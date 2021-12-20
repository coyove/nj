package nj

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"unsafe"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/parser"
	"github.com/coyove/nj/typ"
)

const (
	regA          uint16 = 0xffff
	regPhantom    uint16 = 0xfffe
	regLocalMask         = 0x7fff
	regGlobalFlag        = 0x8000
	maxAddress           = 0x7f00
)

type symbol struct {
	addr uint16
}

func (s *symbol) String() string { return fmt.Sprintf("symbol:%d", s.addr) }

type breakLabel struct {
	continueNode parser.Node
	continueGoto int
	labelPos     []int
}

type CompileOptions struct {
	MaxStackSize int64
	Globals      *Object
	Stdout       io.Writer
	Stderr       io.Writer
	Stdin        io.Reader
}

// symTable is responsible for recording the state of compilation
type symTable struct {
	name    string
	options *CompileOptions

	global *symTable
	parent *symTable

	codeSeg Packet

	// toplevel symtable
	funcs []*Object

	// variable lookup
	sym       map[string]*symbol
	maskedSym []map[string]*symbol

	forLoops []*breakLabel

	vp uint16

	collectConstMode bool
	constMap         map[interface{}]uint16

	reusableTmps      map[uint16]bool
	reusableTmpsArray []uint16

	forwardGoto map[int]string
	labelPos    map[string]int
}

func newSymTable(opt *CompileOptions) *symTable {
	t := &symTable{
		sym:          make(map[string]*symbol),
		constMap:     make(map[interface{}]uint16),
		reusableTmps: make(map[uint16]bool),
		forwardGoto:  make(map[int]string),
		labelPos:     make(map[string]int),
		options:      opt,
	}
	return t
}

func (table *symTable) symbolsToDebugLocals() []string {
	x := make([]string, table.vp)
	for sym, info := range table.sym {
		x[info.addr] = sym
	}
	return x
}

func (table *symTable) borrowAddress() uint16 {
	if len(table.reusableTmpsArray) > 0 {
		tmp := table.reusableTmpsArray[0]
		table.reusableTmpsArray = table.reusableTmpsArray[1:]
		if !table.reusableTmps[tmp] {
			panic("DEBUG: corrupted reusable map")
		}
		table.reusableTmps[tmp] = false
		return tmp
	}
	if table.vp > maxAddress {
		panic("too many variables in a single scope")
	}
	table.reusableTmps[table.vp] = false
	table.vp++
	return table.vp - 1
}

func (table *symTable) freeAddr(a interface{}) {
	switch a := a.(type) {
	case []parser.Node:
		for _, n := range a {
			if n.Type() == parser.ADDR {
				table.freeAddr(n.Addr)
			}
		}
	case []uint16:
		for _, n := range a {
			table.freeAddr(n)
		}
	case uint16:
		if a == regA {
			return
		}
		if a > regLocalMask {
			// We don't free global variables
			return
		}
		if available, existed := table.reusableTmps[a]; existed && !available {
			table.reusableTmpsArray = append(table.reusableTmpsArray, a)
			table.reusableTmps[a] = true
		}

	default:
		panic("DEBUG freeAddr")
	}
}

func (table *symTable) get(varname string) (uint16, bool) {
	depth := uint16(0)
	regNil := table.loadK(nil)

	switch varname {
	case "nil":
		return regNil, true
	case "true":
		return table.loadK(true), true
	case "false":
		return table.loadK(false), true
	case "this":
		if k, ok := table.sym[varname]; ok {
			return k.addr, true
		}
		table.sym["this"] = &symbol{table.borrowAddress()}
	case "$a":
		return regA, true
	}

	calc := func(k *symbol) (uint16, bool) {
		addr := (depth << 15) | (uint16(k.addr) & regLocalMask)
		return addr, true
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

	return regNil, false
}

func (table *symTable) put(varname string, addr uint16) {
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

func (table *symTable) addMaskedSymTable() {
	table.maskedSym = append(table.maskedSym, map[string]*symbol{})
}

func (table *symTable) removeMaskedSymTable() {
	last := table.maskedSym[len(table.maskedSym)-1]
	for _, k := range last {
		table.freeAddr(k.addr)
	}
	table.maskedSym = table.maskedSym[:len(table.maskedSym)-1]
}

func (table *symTable) loadK(v interface{}) uint16 {
	if table.global != nil {
		return table.global.loadK(v)
	}

	if i, ok := table.constMap[v]; ok {
		return i
	}

	if !table.collectConstMode {
		internal.Panic("DEBUG: collect consts %#v", v)
	}

	idx := regGlobalFlag | table.borrowAddress()
	table.constMap[v] = idx
	return idx
}

var operatorMapping = map[string]byte{
	parser.AAdd:     typ.OpAdd,
	parser.ASub:     typ.OpSub,
	parser.AMul:     typ.OpMul,
	parser.ADiv:     typ.OpDiv,
	parser.AIDiv:    typ.OpIDiv,
	parser.AMod:     typ.OpMod,
	parser.ALess:    typ.OpLess,
	parser.ALessEq:  typ.OpLessEq,
	parser.AEq:      typ.OpEq,
	parser.ANeq:     typ.OpNeq,
	parser.ANot:     typ.OpNot,
	parser.ABitAnd:  typ.OpBitAnd,
	parser.ABitOr:   typ.OpBitOr,
	parser.ABitXor:  typ.OpBitXor,
	parser.ABitNot:  typ.OpBitNot,
	parser.ABitLsh:  typ.OpBitLsh,
	parser.ABitRsh:  typ.OpBitRsh,
	parser.ABitURsh: typ.OpBitURsh,
	parser.AStore:   typ.OpStore,
	parser.ALoad:    typ.OpLoad,
	parser.AInc:     typ.OpInc,
	parser.ANext:    typ.OpNext,
	parser.ALen:     typ.OpLen,
	parser.AIs:      typ.OpIsProto,
	parser.AReturn:  typ.OpRet,
}

func (table *symTable) writeInst(op byte, n0, n1 parser.Node) {
	var tmp []uint16
	getAddr := func(n parser.Node, intoNewAddr bool) uint16 {
		switch n.Type() {
		case parser.NODES:
			addr := table.compileNodeInto(n, intoNewAddr, regA)
			tmp = append(tmp, addr)
			return addr
		default:
			addr, ok := table.compileStaticNode(n)
			if !ok {
				internal.Panic("DEBUG writeInst unknown type: %#v", n)
			}
			return addr
		}
	}

	if !n0.Valid() {
		table.codeSeg.writeInst(op, 0, 0)
		return
	}

	n0a := getAddr(n0, n1.Valid())
	if !n1.Valid() {
		table.codeSeg.writeInst(op, n0a, 0)
		table.freeAddr(tmp)
		return
	}

	n1a := getAddr(n1, true)
	table.codeSeg.writeInst(op, n0a, n1a)
	table.freeAddr(tmp)
}

func (table *symTable) compileNodeInto(compound parser.Node, newVar bool, existedVar uint16) uint16 {
	newYX := table.compileNode(compound)

	var yx uint16
	if newVar {
		yx = table.borrowAddress()
	} else {
		yx = existedVar
	}

	table.codeSeg.writeInst(typ.OpSet, yx, newYX)
	return yx
}

func (table *symTable) compileStaticNode(node parser.Node) (uint16, bool) {
	switch node.Type() {
	case parser.ADDR:
		return node.Addr, true
	case parser.STR:
		return table.loadK(node.Str()), true
	case parser.FLOAT:
		return table.loadK(node.Float64()), true
	case parser.INT:
		return table.loadK(node.Int64()), true
	case parser.SYM:
		idx, _ := table.get(node.Sym())
		return idx, true
	}
	return 0, false
}

func (table *symTable) compileNode(node parser.Node) uint16 {
	if addr, ok := table.compileStaticNode(node); ok {
		return addr
	}

	nodes := node.Nodes()
	if len(nodes) == 0 {
		return regA
	}

	name := nodes[0].Sym()
	var yx uint16
	switch name {
	case parser.ADoBlock, parser.ABegin:
		yx = table.compileChain(node)
	case parser.ASet, parser.AMove:
		yx = table.compileSetMove(nodes)
	case parser.AIf:
		yx = table.compileIf(nodes)
	case parser.AFor:
		yx = table.compileWhile(nodes)
	case parser.ABreak, parser.AContinue:
		yx = table.compileBreak(nodes)
	case parser.ACall, parser.ATailCall:
		yx = table.compileCall(nodes)
	case parser.AArray, parser.AObject:
		yx = table.compileList(nodes)
	case parser.AOr, parser.AAnd:
		yx = table.compileAndOr(nodes)
	case parser.AFunc:
		yx = table.compileFunction(nodes)
	case parser.AFreeAddr:
		yx = table.compileFreeAddr(nodes)
	case parser.AGoto, parser.ALabel:
		yx = table.compileGoto(nodes)
	default:
		yx = table.compileOperator(nodes)
	}
	return yx
}

func (table *symTable) collectConsts(node parser.Node) {
	switch node.Type() {
	case parser.STR:
		table.loadK(node.Str())
	case parser.FLOAT:
		table.loadK(node.Float64())
	case parser.INT:
		table.loadK(node.Int64())
	case parser.NODES:
		for _, n := range node.Nodes() {
			table.collectConsts(n)
		}
	}
}

func compileNodeTopLevel(name, source string, n parser.Node, opt *CompileOptions) (cls *Program, err error) {
	defer internal.CatchError(&err)

	table := newSymTable(opt)
	table.collectConstMode = true
	table.name = name
	table.codeSeg.Pos.Name = name
	coreStack := &Env{stack: new([]Value)}

	// Load nil first so it will be at the top
	table.loadK(nil)
	coreStack.Push(Nil)

	push := func(k string, v Value) uint16 {
		idx, ok := table.get(k)
		if ok {
			coreStack.Set(int(idx), v)
		} else {
			idx = uint16(coreStack.Size())
			table.put(k, idx)
			coreStack.Push(v)
		}
		return idx
	}

	Globals.Foreach(func(k Value, v *Value) bool { push(k.String(), *v); return true })

	if opt != nil && opt.Globals != nil {
		opt.Globals.Foreach(func(k Value, v *Value) bool { push(k.String(), *v); return true })
	}

	gi := push("PROGRAM", Nil)
	push("COMPILE_OPTIONS", ValueOf(opt))
	push("SOURCE_CODE", Str(source))

	table.vp = uint16(coreStack.Size())

	// Find and fill consts
	table.loadK(true)
	table.loadK(false)
	table.collectConsts(n)
	table.collectConstMode = false

	table.compileNode(n)
	table.codeSeg.writeInst(typ.OpRet, regA, 0)
	table.patchGoto()

	coreStack.grow(int(table.vp))
	for k, stackPos := range table.constMap {
		switch k := k.(type) {
		case float64:
			coreStack.Set(int(stackPos), Float64(k))
		case int64:
			coreStack.Set(int(stackPos), Int64(k))
		case string:
			coreStack.Set(int(stackPos), Str(k))
		case bool:
			coreStack.Set(int(stackPos), Bool(k))
		case nil:
			coreStack.Set(int(stackPos), Nil)
		default:
			panic("DEBUG")
		}
	}

	cls = &Program{top: &function{}}
	cls.top.Name = "main"
	cls.top.CodeSeg = table.codeSeg
	cls.top.StackSize = table.vp
	cls.top.Locals = table.symbolsToDebugLocals()
	cls.top.LoadGlobal = cls
	cls.stack = coreStack.stack
	cls.symbols = table.sym
	cls.Functions = table.funcs
	if opt != nil {
		cls.Options = *opt
	}
	cls.Options.Stdout = or(cls.Options.Stdout, os.Stdout).(io.Writer)
	cls.Options.Stdin = or(cls.Options.Stdin, os.Stdin).(io.Reader)
	cls.Options.Stderr = or(cls.Options.Stderr, os.Stderr).(io.Writer)
	for _, f := range cls.Functions {
		f.fun.LoadGlobal = cls
	}
	(*cls.stack)[gi] = intf(cls)
	return cls, err
}

func LoadFile(path string, opt *CompileOptions) (*Program, error) {
	code, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return loadCode(*(*string)(unsafe.Pointer(&code)), path, opt)
}

func LoadString(code string, opt *CompileOptions) (*Program, error) {
	return loadCode(code, "<memory>", opt)
}

func loadCode(code, name string, opt *CompileOptions) (*Program, error) {
	n, err := parser.Parse(code, name)
	if err != nil {
		return nil, err
	}
	if internal.IsDebug() {
		n.Dump(os.Stderr, "  ")
	}
	return compileNodeTopLevel(name, code, n, opt)
}

func Run(p *Program, err error) (Value, error) {
	if err != nil {
		return Nil, err
	}
	return p.Run()
}

func MustRun(p *Program, err error) Value {
	internal.PanicErr(err)
	v, err := p.Run()
	internal.PanicErr(err)
	return v
}
