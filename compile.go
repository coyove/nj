package potatolang

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/coyove/potatolang/parser"
)

type symbol struct {
	addr  uint16
	usage int
}

func (s *symbol) String() string { return fmt.Sprintf("symbol:%d", s.addr) }

type gotolabel struct {
	gotoMarker [4]uint32
	labelPos   int
	labelMet   bool
}

type breaklabel struct {
	labelPos []int
}

// symtable is responsible for recording the state of compilation
type symtable struct {
	code packet

	// variable name lookup
	global    *symtable
	sym       map[string]*symbol
	maskedSym []map[string]*symbol

	y      bool // has yield op
	inloop []*breaklabel

	vp uint16

	consts   []interface{}
	constMap map[interface{}]uint16

	reusableTmps map[uint16]bool

	forwardGoto map[int]string
	labelPos    map[string]int
}

func newsymtable() *symtable {
	t := &symtable{
		sym:          make(map[string]*symbol),
		constMap:     make(map[interface{}]uint16),
		reusableTmps: make(map[uint16]bool),
		forwardGoto:  make(map[int]string),
		labelPos:     make(map[string]int),
	}
	return t
}

func (table *symtable) borrowAddress() uint16 {
	for tmp, ok := range table.reusableTmps {
		if ok {
			table.reusableTmps[tmp] = false
			return tmp
		}
	}
	if table.vp > 1000 { //1<<10 {
		panic("Code too complex, may be there are too many variables (1000) in a single scope")
	}
	table.reusableTmps[table.vp] = false
	table.vp++
	return table.vp - 1
}

func (table *symtable) returnAddress(v uint16) {
	if v == regNil || v == regA {
		return
	}
	if v>>10 == 7 {
		return
	}
	//	log.Println("$$", table.reusableTmps, v, table.vp)
	//if v == table.vp-1 {
	//	table.vp--
	//	return
	//}
	if _, existed := table.reusableTmps[v]; existed {
		table.reusableTmps[v] = true
	}
}

func (table *symtable) returnAddresses(a interface{}) {
	switch a := a.(type) {
	case []parser.Node:
		for _, n := range a {
			if n.Type() == parser.ADR {
				table.returnAddress(n.Value.(uint16))
			}
		}
	case []uint16:
		for _, n := range a {
			table.returnAddress(n)
		}
	default:
		panic("returnAddresses: shouldn't happen")
	}
}

func (table *symtable) get(varname parser.Symbol) uint16 {
	depth := uint16(0)

	switch varname.Text {
	case "nil":
		return regNil
	case "true":
		return table.loadK(true)
	case "false":
		return table.loadK(false)
	}

	calc := func(k *symbol) uint16 {
		if depth > 6 || (depth == 6 && k.addr > 1000) {
			panic("too many levels (7) to refer a variable, try simplifing your code")
		}

		addr := (depth << 10) | (uint16(k.addr) & 0x03ff)
		return addr
	}

	for table != nil {
		// Firstly we will iterate the masked symbols
		// Masked symbols are local variables inside do-blocks, like "if then .. end" and "do ... end"
		// The rightmost map of this slice is the innermost do-block
		for i := len(table.maskedSym) - 1; i >= 0; i-- {
			m := table.maskedSym[i]
			if k, ok := m[varname.Text]; ok {
				return calc(k)
			}
		}

		if k, ok := table.sym[varname.Text]; ok {
			return calc(k)
		}

		depth++
		table = table.global
	}

	return regNil
}

func (table *symtable) put(varname string, addr uint16) {
	if addr == regA {
		panic("debug")
	}
	c := math.MaxInt64
	sym := &symbol{
		addr:  addr,
		usage: c,
	}
	if len(table.maskedSym) > 0 {
		table.maskedSym[len(table.maskedSym)-1][varname] = sym
	} else {
		table.sym[varname] = sym
	}
}

func (table *symtable) addMaskedSymTable() {
	table.maskedSym = append(table.maskedSym, map[string]*symbol{})
}

func (table *symtable) removeMaskedSymTable() {
	last := table.maskedSym[len(table.maskedSym)-1]
	for _, k := range last {
		table.returnAddress(k.addr)
	}
	table.maskedSym = table.maskedSym[:len(table.maskedSym)-1]
}

func (table *symtable) loadK(v interface{}) uint16 {
	kaddr := func() uint16 {
		if i, ok := table.constMap[v]; ok {
			return i
		}

		table.consts = append(table.consts, v)
		if len(table.consts) > 1<<10-1 {
			panic("too many constants")
		}

		idx := uint16(len(table.consts) - 1)
		table.constMap[v] = idx
		return idx
	}()

	return 0x7<<10 | kaddr
}

func (table *symtable) constsToValues() []Value {
	consts := make([]Value, len(table.consts))
	for i, k := range table.consts {
		switch k := k.(type) {
		case float64:
			consts[i] = Num(k)
		case int64:
			consts[i] = Int(k)
		case string:
			consts[i] = Str(k)
		case bool:
			consts[i] = NumBool(k)
		}
	}
	return consts
}

var flatOpMapping = map[string]_Opcode{
	parser.AAdd.Text:       OpAdd,
	parser.AConcat.Text:    OpConcat,
	parser.ASub.Text:       OpSub,
	parser.AMul.Text:       OpMul,
	parser.ADiv.Text:       OpDiv,
	parser.AMod.Text:       OpMod,
	parser.ALess.Text:      OpLess,
	parser.ALessEq.Text:    OpLessEq,
	parser.AEq.Text:        OpEq,
	parser.ANeq.Text:       OpNeq,
	parser.ANot.Text:       OpNot,
	parser.APow.Text:       OpPow,
	parser.AStore.Text:     OpStore,
	parser.ALoad.Text:      OpLoad,
	parser.ALen.Text:       OpLen,
	parser.AInc.Text:       OpInc,
	parser.APopV.Text:      OpEOB, // special
	parser.APopVAll.Text:   OpEOB, // special
	parser.APopVAllA.Text:  OpEOB, // special
	parser.APopVClear.Text: OpEOB, // special
}

func (table *symtable) writeOpcode(op _Opcode, n0, n1 parser.Node) {
	var tmp []uint16
	getAddr := func(n parser.Node) uint16 {
		switch n.Type() {
		case parser.CPL:
			addr := table.compileNodeInto(n, true, 0)
			tmp = append(tmp, addr)
			return addr
		case parser.SYM:
			return table.get(n.Value.(parser.Symbol))
		case parser.NUM, parser.STR:
			return table.loadK(n.Value)
		case parser.ADR:
			return n.Value.(uint16)
		default:
			panicf("writeOpcode: shouldn't happend: unknown type: %v", n.TypeName())
			return 0
		}
	}

	defer table.returnAddresses(tmp)

	if !n0.Valid() {
		table.code.writeOP(op, 0, 0)
		return
	}

	n0a := getAddr(n0)
	if !n1.Valid() {
		table.code.writeOP(op, n0a, 0)
		return
	}

	n1a := getAddr(n1)
	if op == OpSet && n0a == n1a {
		return
	}
	table.code.writeOP(op, n0a, n1a)
}

func (table *symtable) compileNodeInto(compound parser.Node, newVar bool, existedVar uint16) uint16 {
	newYX := table.compileNode(compound)

	var yx uint16
	if newVar {
		yx = table.borrowAddress()
	} else {
		yx = existedVar
	}

	table.code.writeOP(OpSet, yx, newYX)
	return yx
}

func (table *symtable) compileNode(node parser.Node) uint16 {
	switch node.Type() {
	case parser.ADR:
		return node.Value.(uint16)
	case parser.STR, parser.NUM:
		return table.loadK(node.Value)
	case parser.SYM:
		return table.get(node.Sym())
	}

	nodes := node.Cpl()
	if len(nodes) == 0 {
		return regA
	}
	name, ok := nodes[0].Value.(parser.Symbol)
	if !ok {
		panicf("compileNode: shouldn't happend: invalid op: %v", nodes)
	}

	var yx uint16
	switch name.Text {
	case parser.ADoBlock.Text, parser.ABegin.Text:
		yx = table.compileChainOp(node)
	case parser.ASet.Text, parser.AMove.Text:
		yx = table.compileSetOp(nodes)
	case parser.AReturn.Text, parser.AYield.Text:
		yx = table.compileRetOp(nodes)
	case parser.AIf.Text:
		yx = table.compileIfOp(nodes)
	case parser.AFor.Text:
		yx = table.compileWhileOp(nodes)
	case parser.AContinue.Text, parser.ABreak.Text:
		yx = table.compileBreakOp(nodes)
	case parser.ACall.Text, parser.ATailCall.Text:
		yx = table.compileCallOp(nodes)
	case parser.AOr.Text, parser.AAnd.Text:
		yx = table.compileAndOrOp(nodes)
	case parser.AFunc.Text:
		yx = table.compileLambdaOp(nodes)
	case parser.ARetAddr.Text:
		yx = table.compileRetAddrOp(nodes)
	case parser.AGoto.Text, parser.ALabel.Text:
		yx = table.compileGotoOp(nodes)
	default:
		if _, ok := flatOpMapping[name.Text]; ok {
			return table.compileFlatOp(nodes)
		}
		panicf("compileNode: shouldn't happen: unknown symbol: %s", name)
	}
	return yx
}

func compileNodeTopLevel(n parser.Node) (cls *Func, err error) {
	defer func() {
		if r := recover(); r != nil {
			cls = nil
			if err, _ = r.(error); err == nil {
				err = fmt.Errorf("recovered panic: %v", r)
			}
			if os.Getenv("PL_STACK") != "" {
				log.Println(string(debug.Stack()))
			}
		}
		if os.Getenv("PL_STACK") != "" {
			n.Dump(os.Stderr, "")
		}
	}()

	table := newsymtable()

	coreStack := &Env{stack: new([]Value)}
	for k, v := range g {
		table.put(k, uint16(coreStack.Size()))
		coreStack.Push(v)
	}

	table.vp = uint16(len(table.sym))
	table.compileNode(n)
	table.code.writeOP(OpEOB, 0, 0)
	table.patchGoto()
	cls = &Func{}
	cls.Name = "main"
	cls.packet = table.code
	cls.ConstTable = table.constsToValues()
	cls.yEnv.stack = coreStack.stack
	cls.yEnv.grow(int(table.vp))
	cls.yEnv.global = &Global{Stack: cls.yEnv.stack} // yEnv itself is the global stack
	cls.stackSize = table.vp
	return cls, err
}

func LoadFile(path string) (*Func, error) {
	code, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	n, err := parser.Parse(bytes.NewReader(code), path)
	if err != nil {
		return nil, err
	}
	// n.Dump(os.Stderr, "  ")
	return compileNodeTopLevel(n)
}

func LoadString(code string) (*Func, error) {
	_, fn, _, _ := runtime.Caller(1)
	return loadStringName(code, fn)
}

func loadStringName(code, name string) (*Func, error) {
	n, err := parser.Parse(bytes.NewReader([]byte(code)), name)
	if err != nil {
		return nil, err
	}
	// n.Dump(os.Stderr, "  ")
	return compileNodeTopLevel(n)
}

func WithTimeout(f *Func, d time.Duration) { f.yEnv.global.Deadline = time.Now().Add(d).Unix() }

func WithDeadline(f *Func, d time.Time) { f.yEnv.global.Deadline = d.Unix() }

func WithMaxStackSize(f *Func, sz int64) { f.yEnv.global.MaxStackSize = sz }

func WithMaxStringSize(f *Func, sz int64) { f.yEnv.global.MaxStringSize = sz }

func WithValue(f *Func, k string, v interface{}) {
	if f.yEnv.global.Extras == nil {
		f.yEnv.global.Extras = map[string]interface{}{}
	}
	f.yEnv.global.Extras[k] = v
}
