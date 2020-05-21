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
	"strconv"
	"strings"

	"github.com/coyove/potatolang/parser"
)

type symbol struct {
	addr  uint16
	usage int
}

// symtable is responsible for recording the state of compilation
type symtable struct {
	// variable name lookup
	parent    *symtable
	sym       map[parser.Atom]*symbol
	maskedSym []map[parser.Atom]*symbol

	y         bool // has yield op
	envescape bool
	inloop    bool

	vp uint16

	consts   []interface{}
	constMap map[interface{}]uint16

	reusableTmps map[uint16]bool
}

func newsymtable() *symtable {
	t := &symtable{
		sym:          make(map[parser.Atom]*symbol),
		constMap:     make(map[interface{}]uint16),
		reusableTmps: make(map[uint16]bool),
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
	case []*parser.Node:
		for _, n := range a {
			if n.Type() == parser.Naddr {
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

func (table *symtable) get(varname parser.Atom) uint16 {
	depth := uint16(0)

	switch varname {
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

		if k.usage--; k.usage == 0 {
			table.returnAddress(k.addr)
			delete(table.sym, varname)
		}

		return addr
	}

	for table != nil {
		// Firstly we will iterate the masked symbols
		// Masked symbols are local variables inside do-blocks, like "if then .. end" and "do ... end"
		// The rightmost map of this slice is the innermost do-block
		for i := len(table.maskedSym) - 1; i >= 0; i-- {
			m := table.maskedSym[i]
			if k, ok := m[varname]; ok {
				return calc(k)
			}
		}

		if k, ok := table.sym[varname]; ok {
			return calc(k)
		}

		depth++
		table = table.parent
	}

	return regNil
}

func (table *symtable) put(varname parser.Atom, addr uint16) {
	if addr == regA {
		panic("debug")
	}
	c := math.MaxInt64
	if strings.HasPrefix(string(varname), "(") {
		c, _ = strconv.Atoi(string(varname[1:strings.Index(string(varname), ")")]))
	}
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
	table.maskedSym = append(table.maskedSym, map[parser.Atom]*symbol{})
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

var flatOpMapping = map[parser.Atom]_Opcode{
	parser.AAdd:         OpAdd,
	parser.AConcat:      OpConcat,
	parser.ASub:         OpSub,
	parser.AMul:         OpMul,
	parser.ADiv:         OpDiv,
	parser.AMod:         OpMod,
	parser.ALess:        OpLess,
	parser.ALessEq:      OpLessEq,
	parser.AEq:          OpEq,
	parser.ANeq:         OpNeq,
	parser.ANot:         OpNot,
	parser.ABitAnd:      OpBitAnd,
	parser.ABitOr:       OpBitOr,
	parser.ABitXor:      OpBitXor,
	parser.ABitLsh:      OpBitLsh,
	parser.ABitRsh:      OpBitRsh,
	parser.ABitURsh:     OpBitURsh,
	parser.AStore:       OpStore,
	parser.ALoad:        OpLoad,
	parser.ALen:         OpLen,
	parser.APatchVararg: OpPatchVararg,
	parser.AAddrOf:      OpAddressOf,
	parser.AInc:         OpInc,
	parser.ASetB:        OpSetB,
	parser.AGetB:        OpGetB,
}

func (table *symtable) writeOpcode(buf *packet, op _Opcode, n0, n1 *parser.Node) {
	var tmp []uint16
	getAddr := func(n *parser.Node) uint16 {
		switch n.Type() {
		case parser.Ncomplex:
			code, addr := table.compileNodeInto(n, true, 0)
			buf.Write(code)
			tmp = append(tmp, addr)
			return addr
		case parser.Natom:
			return table.get(n.Value.(parser.Atom))
		case parser.Nnumber, parser.Nstring:
			return table.loadK(n.Value)
		case parser.Naddr:
			return n.Value.(uint16)
		default:
			panicf("writeOpcode: shouldn't happend: unknown type: %v", n.TypeName())
			return 0
		}
	}

	defer table.returnAddresses(tmp)

	if n0 == nil {
		buf.WriteOP(op, 0, 0)
		return
	}

	n0a := getAddr(n0)
	if n1 == nil {
		buf.WriteOP(op, n0a, 0)
		return
	}

	n1a := getAddr(n1)
	if op == OpSet && n0a == n1a {
		return
	}
	buf.WriteOP(op, n0a, n1a)
}

func (table *symtable) compileNodeInto(compound *parser.Node, newVar bool, existedVar uint16) (code packet, yx uint16) {
	buf := newpacket()

	var newYX uint16
	code, newYX = table.compileNode(compound)

	buf.Write(code)
	if newVar {
		yx = table.borrowAddress()
	} else {
		yx = existedVar
	}

	buf.WriteOP(OpSet, yx, newYX)
	return buf, yx
}

func (table *symtable) compileNode(node *parser.Node) (code packet, yx uint16) {
	switch node.Type() {
	case parser.Naddr:
		return code, node.Value.(uint16)
	case parser.Nstring, parser.Nnumber:
		return code, table.loadK(node.Value)
	case parser.Natom:
		return code, table.get(node.A())
	}

	nodes := node.C()
	if len(nodes) == 0 {
		return newpacket(), regA
	}
	name, ok := nodes[0].Value.(parser.Atom)
	if !ok {
		nodes[0].Dump(os.Stderr)
		panicf("compileNode: shouldn't happend: invalid op: %v", nodes)
	}

	switch name {
	case parser.ADoBlock, parser.AChain:
		code, yx = table.compileChainOp(node)
	case parser.ASet, parser.AMove:
		code, yx = table.compileSetOp(nodes)
	case parser.AReturn:
		code, yx = table.writeOpcode3(OpRet, nodes)
	case parser.AYield:
		table.y = true
		code, yx = table.writeOpcode3(OpYield, nodes)
	case parser.AIf:
		code, yx = table.compileIfOp(nodes)
	case parser.AFor:
		code, yx = table.compileWhileOp(nodes)
	case parser.AContinue, parser.ABreak:
		code, yx = table.compileContinueBreakOp(nodes)
	case parser.ACall:
		code, yx = table.compileCallOp(nodes)
	case parser.AHash, parser.AHashArray, parser.AArray:
		code, yx = table.compileHashArrayOp(nodes)
	case parser.AOr, parser.AAnd:
		code, yx = table.compileAndOrOp(nodes)
	case parser.AFunc:
		code, yx = table.compileLambdaOp(nodes)
	default:
		if _, ok := flatOpMapping[name]; ok {
			return table.compileFlatOp(nodes)
		}
		panicf("compileNode: shouldn't happen: unknown symbol: %s", name)
	}
	return
}

func compileNodeTopLevel(n *parser.Node) (cls *Closure, err error) {
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
			log.Println(n)
		}
	}()

	table := newsymtable()

	coreStack := NewEnv(nil)
	for n, v := range CoreLibs {
		table.put(parser.Atom(n), uint16(coreStack.LocalSize()))
		coreStack.LocalPush(v)
	}

	table.vp = uint16(len(table.sym))
	code, _ := table.compileNode(n)
	code.WriteOP(OpEOB, 0, 0)
	consts := make([]Value, len(table.consts))
	for i, k := range table.consts {
		switch k := k.(type) {
		case float64:
			consts[i] = Num(k)
		case string:
			consts[i] = Str(k)
		case bool:
			consts[i] = Bln(k)
		}
	}
	cls = NewClosure(code.data, consts, nil, 0)
	cls.lastenv = NewEnv(nil)
	cls.Pos = code.pos
	cls.source = []byte("root" + cls.String() + "@" + code.source)
	cls.lastenv.stack = coreStack.stack
	return cls, err
}

func LoadFile(path string) (*Closure, error) {
	code, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	n, err := parser.Parse(bytes.NewReader(code), path)
	if err != nil {
		return nil, err
	}
	// n.Dump(os.Stderr)
	return compileNodeTopLevel(n)
}

func LoadString(code string) (*Closure, error) {
	_, fn, _, _ := runtime.Caller(1)
	return loadStringName(code, fn)
}

func loadStringName(code, name string) (*Closure, error) {
	n, err := parser.Parse(bytes.NewReader([]byte(code)), name)
	if err != nil {
		return nil, err
	}
	// n.Dump(os.Stderr)
	return compileNodeTopLevel(n)
}
