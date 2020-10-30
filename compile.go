package script

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime/debug"

	"github.com/coyove/script/parser"
)

type symbol struct {
	addr uint16
}

func (s *symbol) String() string { return fmt.Sprintf("symbol:%d", s.addr) }

type breaklabel struct {
	labelPos []int
}

// symtable is responsible for recording the state of compilation
type symtable struct {
	global *symtable

	code packet

	// toplevel symtable
	funcs []*Func

	// variable name lookup
	sym       map[string]*symbol
	maskedSym []map[string]*symbol

	inloop []*breaklabel

	vp uint16

	insideJSONGenerator bool

	collectConstsMode bool
	constMap          map[interface{}]uint16

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
	if table.vp > 4000 { // 12 bits {
		panic("too many variables (4000) in a single scope")
	}
	table.reusableTmps[table.vp] = false
	table.vp++
	return table.vp - 1
}

func (table *symtable) returnAddress(v uint16) {
	if v == regA {
		return
	}
	if v>>12 == 1 {
		panic("DEBUG")
	}
	if _, existed := table.reusableTmps[v]; existed {
		table.reusableTmps[v] = true
	}
}

func (table *symtable) returnAddresses(a interface{}) {
	switch a := a.(type) {
	case []parser.Node:
		for _, n := range a {
			if n.Type == parser.Address {
				table.returnAddress(n.Addr)
			}
		}
	case []uint16:
		for _, n := range a {
			table.returnAddress(n)
		}
	default:
		panic("DEBUG returnAddresses")
	}
}

func (table *symtable) get(varname string) uint16 {
	depth := uint16(0)
	regNil := table.loadK(nil)

	switch varname {
	case "nil":
		return regNil
	case "true":
		return table.loadK(int64(1))
	case "false":
		return table.loadK(int64(0))
	}

	calc := func(k *symbol) uint16 {
		addr := (depth << 12) | (uint16(k.addr) & 0xfff)
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
		table = table.global
	}

	return regNil
}

func (table *symtable) put(varname string, addr uint16) {
	if addr == regA {
		panic("DEBUG: put $a?")
	}
	sym := &symbol{
		addr: addr,
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
	if table.global != nil {
		return 1<<12 | table.global.loadK(v)
	}

	if i, ok := table.constMap[v]; ok {
		return i
	}

	if !table.collectConstsMode {
		panicf("DEBUG: collect consts %#v", v)
	}

	idx := table.borrowAddress()
	table.constMap[v] = idx
	return idx
}

var flatOpMapping = map[string]opCode{
	parser.AAdd:       OpAdd,
	parser.AConcat:    OpConcat,
	parser.ASub:       OpSub,
	parser.AMul:       OpMul,
	parser.ADiv:       OpDiv,
	parser.AMod:       OpMod,
	parser.ALess:      OpLess,
	parser.ALessEq:    OpLessEq,
	parser.AEq:        OpEq,
	parser.ANeq:       OpNeq,
	parser.ANot:       OpNot,
	parser.APow:       OpPow,
	parser.AStore:     OpStore,
	parser.ALoad:      OpLoad,
	parser.ASlice:     OpSlice,
	parser.ALen:       OpLen,
	parser.AInc:       OpInc,
	parser.APopV:      OpRet, // special
	parser.APopVAll:   OpRet, // special
	parser.APopVAllA:  OpRet, // special
	parser.APopVClear: OpRet, // special
}

func (table *symtable) writeInst(op opCode, n0, n1 parser.Node) {
	var tmp []uint16
	getAddr := func(n parser.Node) uint16 {
		switch n.Type {
		case parser.Complex:
			addr := table.compileNodeInto(n, true, 0)
			tmp = append(tmp, addr)
			return addr
		case parser.Symbol:
			return table.get(n.SymbolValue())
		case parser.String:
			return table.loadK(n.StringValue())
		case parser.Float:
			return table.loadK(n.FloatValue())
		case parser.Int:
			return table.loadK(n.IntValue())
		case parser.Address:
			return n.Addr
		default:
			panicf("DEBUG writeInst unknown type: %#v", n)
			return 0
		}
	}

	defer table.returnAddresses(tmp)

	if !n0.Valid() {
		table.code.writeInst(op, 0, 0)
		return
	}

	n0a := getAddr(n0)
	if !n1.Valid() {
		table.code.writeInst(op, n0a, 0)
		return
	}

	n1a := getAddr(n1)
	if op == OpSet && n0a == n1a {
		return
	}
	table.code.writeInst(op, n0a, n1a)
}

func (table *symtable) compileNodeInto(compound parser.Node, newVar bool, existedVar uint16) uint16 {
	newYX := table.compileNode(compound)

	var yx uint16
	if newVar {
		yx = table.borrowAddress()
	} else {
		yx = existedVar
	}

	table.code.writeInst(OpSet, yx, newYX)
	return yx
}

func (table *symtable) compileNode(node parser.Node) uint16 {
	switch node.Type {
	case parser.Address:
		return node.Addr
	case parser.String:
		return table.loadK(node.StringValue())
	case parser.Float:
		return table.loadK(node.FloatValue())
	case parser.Int:
		return table.loadK(node.IntValue())
	case parser.Symbol:
		return table.get(node.SymbolValue())
	}

	nodes := node.Nodes
	if len(nodes) == 0 {
		return regA
	}

	name := nodes[0].SymbolValue()
	var yx uint16
	switch name {
	case parser.ADoBlock, parser.ABegin:
		yx = table.compileChainOp(node)
	case parser.ASet, parser.AMove:
		yx = table.compileSetOp(nodes)
	case parser.AReturn, parser.AYield:
		yx = table.compileRetOp(nodes)
	case parser.AIf:
		yx = table.compileIfOp(nodes)
	case parser.AFor:
		yx = table.compileWhileOp(nodes)
	case parser.ABreak:
		yx = table.compileBreakOp(nodes)
	case parser.ACall, parser.ATailCall, parser.ACallMap:
		yx = table.compileCallOp(nodes)
	case parser.AOr, parser.AAnd:
		yx = table.compileAndOrOp(nodes)
	case parser.AFunc:
		yx = table.compileLambdaOp(nodes)
	case parser.ARetAddr:
		yx = table.compileRetAddrOp(nodes)
	case parser.AGoto, parser.ALabel:
		yx = table.compileGotoOp(nodes)
	case parser.AJSON:
		yx = table.compileJSONOp(nodes)
	default:
		if _, ok := flatOpMapping[name]; ok {
			return table.compileFlatOp(nodes)
		}
		panicf("DEBUG: compileNode unknown symbol: %q", name)
	}
	return yx
}

func (table *symtable) collectConsts(node parser.Node) {
	switch node.Type {
	case parser.String:
		table.loadK(node.StringValue())
	case parser.Float:
		table.loadK(node.FloatValue())
	case parser.Int:
		table.loadK(node.IntValue())
	case parser.Complex:
		for _, n := range node.Nodes {
			table.collectConsts(n)
		}
	}
}

func compileNodeTopLevel(n parser.Node, globalKeyValues ...interface{}) (cls *Program, err error) {
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

	if len(globalKeyValues)%2 != 0 {
		globalKeyValues = append(globalKeyValues, nil)
	}
	for i := 0; i < len(globalKeyValues); i += 2 {
		k, ok := globalKeyValues[i].(string)
		if ok {
			table.put(k, uint16(coreStack.Size()))
			coreStack.Push(Interface(globalKeyValues[i+1]))
		}
	}

	table.vp = uint16(coreStack.Size())

	// Find and fill consts
	table.collectConstsMode = true
	table.loadK(nil)
	table.loadK(int64(1))
	table.loadK(int64(0))
	table.collectConsts(n)
	table.collectConstsMode = false

	table.compileNode(n)
	table.code.writeInst(OpRet, table.loadK(nil), 0)
	table.patchGoto()

	coreStack.grow(int(table.vp))
	for k, stackPos := range table.constMap {
		switch k := k.(type) {
		case float64:
			coreStack.Set(int(stackPos), Float(k))
		case int64:
			coreStack.Set(int(stackPos), Int(k))
		case string:
			coreStack.Set(int(stackPos), _str(k))
		case nil:
			coreStack.Set(int(stackPos), Value{})
		default:
			panic("DEBUG")
		}
	}

	cls = &Program{}
	cls.name = "main"
	cls.code = table.code
	cls.stackSize = table.vp
	cls.Stack = coreStack.stack
	cls.Funcs = table.funcs
	cls.Stdout = os.Stdout
	cls.Stdin = os.Stdin
	cls.Stderr = os.Stderr
	for _, f := range cls.Funcs {
		f.loadGlobal = cls
	}
	cls.loadGlobal = cls
	cls.NilIndex = table.loadK(nil)
	return cls, err
}

func LoadFile(path string, globalKeyValues ...interface{}) (*Program, error) {
	code, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	n, err := parser.Parse(bytes.NewReader(code), path)
	if err != nil {
		return nil, err
	}
	// n.Dump(os.Stderr, "  ")
	return compileNodeTopLevel(n, globalKeyValues...)
}

func LoadString(code string, globalKeyValues ...interface{}) (*Program, error) {
	n, err := parser.Parse(bytes.NewReader([]byte(code)), "")
	if err != nil {
		return nil, err
	}
	// n.Dump(os.Stderr, "  ")
	return compileNodeTopLevel(n, globalKeyValues...)
}

func MustRun(p *Program, err error) (Value, []Value) {
	if err != nil {
		panic(err)
	}
	v, v1, err := p.Run()
	if err != nil {
		panic(err)
	}
	return v, v1
}
