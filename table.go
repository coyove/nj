package nj

import (
	"fmt"
	"io"

	"github.com/coyove/nj/bas"
	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/parser"
	"github.com/coyove/nj/typ"
)

type breakLabel struct {
	continueNode     parser.Node
	continueGoto     int
	breakContinuePos []int
}

// symTable is responsible for recording the state of compilation
type symTable struct {
	name    string
	options *LoadOptions

	// toplevel symtable
	global *symTable

	codeSeg internal.Packet

	// variable lookup
	sym       bas.Object   // str -> address: uint16
	maskedSym []bas.Object // str -> address: uint16

	forLoops []*breakLabel

	vp uint16

	collectConstMode bool
	constMap         bas.Object // value -> address: uint16
	funcsMap         bas.Object // func name -> address: uint16

	reusableTmps      bas.Object // address: uint16 -> used: bool
	reusableTmpsArray []uint16

	forwardGoto map[int]parser.Node // position to goto label node
	labelPos    map[string]int      // label name to position
}

func newSymTable(opt *LoadOptions) *symTable {
	t := &symTable{
		options: opt,
	}
	return t
}

func (table *symTable) panicnode(node parser.Node, msg string, args ...interface{}) {
	panic(fmt.Sprintf("%s at %s:%d\t", node.Value.Str(), table.name, node.SymLine) + fmt.Sprintf(msg, args...))
}

func (table *symTable) symbolsToDebugLocals() []string {
	x := make([]string, table.vp)
	table.sym.Foreach(func(sym bas.Value, addr *bas.Value) bool {
		x[addr.Int64()] = sym.Str()
		return true
	})
	return x
}

func (table *symTable) borrowAddress() uint16 {
	if len(table.reusableTmpsArray) > 0 {
		tmp := bas.Int64(int64(table.reusableTmpsArray[0]))
		table.reusableTmpsArray = table.reusableTmpsArray[1:]
		if table.reusableTmps.Get(tmp).IsFalse() {
			internal.ShouldNotHappen()
		}
		table.reusableTmps.Set(tmp, bas.False)
		return uint16(tmp.Int64())
	}
	if table.vp > typ.RegMaxAddress {
		panic("too many variables in a single scope")
	}
	return table.borrowAddressNoReuse()
}

func (table *symTable) borrowAddressNoReuse() uint16 {
	table.reusableTmps.Set(bas.Int64(int64(table.vp)), bas.False)
	table.vp++
	return table.vp - 1
}

func (table *symTable) freeAddr(a interface{}) {
	switch a := a.(type) {
	case []parser.Node:
		for _, n := range a {
			if n.Type() == parser.ADDR {
				table.freeAddr(uint16(n.Int()))
			}
		}
	case []uint16:
		for _, n := range a {
			table.freeAddr(n)
		}
	case uint16:
		if a == typ.RegA {
			return
		}
		if a > typ.RegLocalMask {
			// We don't free global variables
			return
		}
		if available := table.reusableTmps.Get(bas.Int64(int64(a))); available != bas.Nil && available.IsFalse() {
			table.reusableTmpsArray = append(table.reusableTmpsArray, a)
			table.reusableTmps.Set(bas.Int64(int64(a)), bas.True)
		}

	default:
		internal.ShouldNotHappen()
	}
}

var (
	staticNil    = bas.Str("nil")
	staticTrue   = bas.Str("true")
	staticFalse  = bas.Str("false")
	staticThis   = bas.Str("this")
	staticSelf   = bas.Str("self")
	staticA      = bas.Str("$a")
	nodeCompiler = map[bas.Value]func(*symTable, []parser.Node) uint16{}
)

func init() {
	nodeCompiler[parser.SDoBlock.Value] = compileProgBlock
	nodeCompiler[parser.SBegin.Value] = compileProgBlock
	nodeCompiler[parser.SSet.Value] = compileSetMove
	nodeCompiler[parser.SMove.Value] = compileSetMove
	nodeCompiler[parser.SIf.Value] = compileIf
	nodeCompiler[parser.SWhile.Value] = compileWhile
	nodeCompiler[parser.SBreak.Value] = compileBreak
	nodeCompiler[parser.SContinue.Value] = compileBreak
	nodeCompiler[parser.SCall.Value] = compileCall
	nodeCompiler[parser.STailCall.Value] = compileCall
	nodeCompiler[parser.STryCall.Value] = compileCall
	nodeCompiler[parser.SArray.Value] = compileArray
	nodeCompiler[parser.SObject.Value] = compileObject
	nodeCompiler[parser.SOr.Value] = compileAndOr
	nodeCompiler[parser.SAnd.Value] = compileAndOr
	nodeCompiler[parser.SFunc.Value] = compileFunction
	nodeCompiler[parser.SFreeAddr.Value] = compileFreeAddr
	nodeCompiler[parser.SGoto.Value] = compileGoto
	nodeCompiler[parser.SLabel.Value] = compileLabel
	nodeCompiler[parser.SAdd.Value] = makeOPCompiler(typ.OpAdd)
	nodeCompiler[parser.SSub.Value] = makeOPCompiler(typ.OpSub)
	nodeCompiler[parser.SMul.Value] = makeOPCompiler(typ.OpMul)
	nodeCompiler[parser.SDiv.Value] = makeOPCompiler(typ.OpDiv)
	nodeCompiler[parser.SIDiv.Value] = makeOPCompiler(typ.OpIDiv)
	nodeCompiler[parser.SMod.Value] = makeOPCompiler(typ.OpMod)
	nodeCompiler[parser.SLess.Value] = makeOPCompiler(typ.OpLess)
	nodeCompiler[parser.SLessEq.Value] = makeOPCompiler(typ.OpLessEq)
	nodeCompiler[parser.SEq.Value] = makeOPCompiler(typ.OpEq)
	nodeCompiler[parser.SNeq.Value] = makeOPCompiler(typ.OpNeq)
	nodeCompiler[parser.SNot.Value] = makeOPCompiler(typ.OpNot)
	nodeCompiler[parser.SBitAnd.Value] = makeOPCompiler(typ.OpBitAnd)
	nodeCompiler[parser.SBitOr.Value] = makeOPCompiler(typ.OpBitOr)
	nodeCompiler[parser.SBitXor.Value] = makeOPCompiler(typ.OpBitXor)
	nodeCompiler[parser.SBitNot.Value] = makeOPCompiler(typ.OpBitNot)
	nodeCompiler[parser.SBitLsh.Value] = makeOPCompiler(typ.OpBitLsh)
	nodeCompiler[parser.SBitRsh.Value] = makeOPCompiler(typ.OpBitRsh)
	nodeCompiler[parser.SBitURsh.Value] = makeOPCompiler(typ.OpBitURsh)
	nodeCompiler[parser.SStore.Value] = makeOPCompiler(typ.OpStore)
	nodeCompiler[parser.SLoad.Value] = makeOPCompiler(typ.OpLoad)
	nodeCompiler[parser.SSlice.Value] = makeOPCompiler(typ.OpSlice)
	nodeCompiler[parser.SInc.Value] = makeOPCompiler(typ.OpInc)
	nodeCompiler[parser.SNext.Value] = makeOPCompiler(typ.OpNext)
	nodeCompiler[parser.SLen.Value] = makeOPCompiler(typ.OpLen)
	nodeCompiler[parser.SIs.Value] = makeOPCompiler(typ.OpIsProto)
	nodeCompiler[parser.SReturn.Value] = makeOPCompiler(typ.OpRet)
}

func (table *symTable) get(name bas.Value) (uint16, bool) {
	switch name {
	case staticNil:
		return typ.RegGlobalFlag, true
	case staticTrue:
		return table.loadConst(bas.True), true
	case staticFalse:
		return table.loadConst(bas.False), true
	case staticThis, staticSelf:
		k := table.sym.Get(name)
		if k.Type() == typ.Number {
			return uint16(k.Int64()), true
		}
		k = bas.Int64(int64(table.borrowAddressNoReuse()))
		table.sym.Set(name, k)
		return uint16(k.Int64()), true
	case staticA:
		return typ.RegA, true
	}

	calc := func(k uint16, depth uint16) (uint16, bool) {
		addr := (depth << 15) | (k & typ.RegLocalMask)
		return addr, true
	}

	// Firstly we will iterate local masked symbols,
	// which are local variables inside do-blocks, like "if then .. end" and "do ... end".
	// The rightmost map of this slice is the innermost do-block
	for i := len(table.maskedSym) - 1; i >= 0; i-- {
		if k := table.maskedSym[i].Get(name); k != bas.Nil {
			return calc(uint16(k.Int64()), 0)
		}
	}

	// Then local variables
	if k := table.sym.Get(name); k != bas.Nil {
		return calc(uint16(k.Int64()), 0)
	}

	// Finally global variables
	if table.global != nil {
		if k := table.global.sym.Get(name); k != bas.Nil {
			return calc(uint16(k.Int64()), 1)
		}
	}

	return typ.RegGlobalFlag, false
}

func (table *symTable) put(name bas.Value, addr uint16) {
	if addr == typ.RegA {
		internal.ShouldNotHappen()
	}
	sym := bas.Int64(int64(addr))
	if len(table.maskedSym) > 0 {
		table.maskedSym[len(table.maskedSym)-1].Set(name, sym)
	} else {
		table.sym.Set(name, sym)
	}
}

func (table *symTable) addMaskedSymTable() {
	table.maskedSym = append(table.maskedSym, bas.Object{})
}

func (table *symTable) removeMaskedSymTable() {
	table.maskedSym[len(table.maskedSym)-1].Foreach(func(sym bas.Value, addr *bas.Value) bool {
		table.freeAddr(uint16(addr.Int64()))
		return true
	})
	table.maskedSym = table.maskedSym[:len(table.maskedSym)-1]
}

func (table *symTable) loadConst(v bas.Value) uint16 {
	if table.global != nil {
		return table.global.loadConst(v)
	}

	if i := table.constMap.Get(v); i != bas.Nil {
		return uint16(i.Int64())
	}

	if !table.collectConstMode {
		internal.ShouldNotHappen(v)
	}

	idx := bas.Int64(int64(typ.RegGlobalFlag | table.borrowAddress()))
	table.constMap.Set(v, idx)
	return uint16(idx.Int64())
}

func (table *symTable) writeInst1(op byte, n parser.Node) {
	if !n.Valid() {
		internal.ShouldNotHappen(n)
	}
	if n.Type() == parser.NODES {
		table.codeSeg.WriteInst(op, table.compileNode(n), 0)
	} else {
		addr, ok := table.compileStaticNode(n)
		if !ok {
			internal.ShouldNotHappen(n)
		}
		table.codeSeg.WriteInst(op, addr, 0)
	}
}

func (table *symTable) compileAtom(n parser.Node, tmp *[]uint16) uint16 {
	switch n.Type() {
	case parser.NODES:
		addr := table.borrowAddress()
		table.codeSeg.WriteInst(typ.OpSet, addr, table.compileNode(n))
		*tmp = append(*tmp, addr)
		return addr
	default:
		addr, ok := table.compileStaticNode(n)
		if !ok {
			internal.ShouldNotHappen(n)
		}
		return addr
	}
}

func (table *symTable) writeInst2(op byte, n0, n1 parser.Node) {
	if !n0.Valid() || !n1.Valid() {
		internal.ShouldNotHappen(n0, n1)
	}
	var tmp []uint16
	table.codeSeg.WriteInst(op, table.compileAtom(n0, &tmp), table.compileAtom(n1, &tmp))
	table.freeAddr(tmp)
}

func (table *symTable) writeInst3(op byte, n0, n1, n2 parser.Node) {
	if !n0.Valid() || !n1.Valid() || !n2.Valid() {
		internal.ShouldNotHappen(n0, n1)
	}
	var tmp []uint16
	table.codeSeg.WriteInst3(op, table.compileAtom(n0, &tmp), table.compileAtom(n1, &tmp), table.compileAtom(n2, &tmp))
	table.freeAddr(tmp)
}

func (table *symTable) compileStaticNode(node parser.Node) (uint16, bool) {
	switch node.Type() {
	case parser.ADDR:
		return uint16(node.Int()), true
	case parser.STR, parser.FLOAT, parser.INT:
		return table.loadConst(node.Value), true
	case parser.SYM:
		idx, ok := table.get(node.Value)
		if !ok {
			table.panicnode(node, "symbol not defined")
		}
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
		return typ.RegA
	}

	return nodeCompiler[nodes[0].Value](table, nodes)
}

func (table *symTable) collectConsts(node parser.Node) {
	switch node.Type() {
	case parser.STR, parser.FLOAT, parser.INT:
		table.loadConst(node.Value)
	case parser.NODES:
		nodes := node.Nodes()
		if len(nodes) == 4 && nodes[0].Value == parser.SFunc.Value {
			addr := int(typ.RegGlobalFlag | table.borrowAddressNoReuse())
			if table.getGlobal().funcsMap.Set(nodes[1].Value, bas.Int(addr)) != bas.Nil {
				table.panicnode(nodes[1], "duplicated function names")
			}
		}
		for _, n := range nodes {
			table.collectConsts(n)
		}
	}
}

func (table *symTable) getGlobal() *symTable {
	if table.global != nil {
		return table.global
	}
	return table
}

func compileNodeTopLevel(name, source string, n parser.Node, opt *LoadOptions) (cls *bas.Program, err error) {
	defer internal.CatchError(&err)

	table := newSymTable(opt)
	table.name = name
	table.codeSeg.Pos.Name = name
	// Load nil first to ensure its address == 0
	table.borrowAddress()

	obj, coreStack := bas.GetGlobalsStack()

	push := func(k, v bas.Value) uint16 {
		idx, ok := table.get(k)
		if ok {
			coreStack[idx] = v
		} else {
			idx = uint16(len(coreStack))
			table.put(k, idx)
			coreStack = append(coreStack, v)
		}
		return idx
	}

	table.sym = *obj

	if opt != nil && opt.Globals != nil {
		opt.Globals.Foreach(func(k bas.Value, v *bas.Value) bool { push(k, *v); return true })
	}

	gi := push(bas.Str("PROGRAM"), bas.Nil)
	push(bas.Str("SOURCE_CODE"), bas.Str(source))

	table.vp = uint16(len(coreStack))

	// Search constants and functions
	table.collectConstMode = true
	table.loadConst(bas.True)
	table.loadConst(bas.False)
	table.collectConsts(n)
	table.collectConstMode = false

	table.compileNode(n)
	table.codeSeg.WriteInst(typ.OpRet, typ.RegA, 0)
	table.patchGoto()

	coreStack = append(coreStack, make([]bas.Value, int(table.vp)-len(coreStack))...)
	table.constMap.Foreach(func(konst bas.Value, addr *bas.Value) bool {
		coreStack[addr.Int64()&typ.RegLocalMask] = konst
		return true
	})

	cls = internal.NewProgram(
		&coreStack,
		internal.NewFunc(
			"main",
			false,
			0,
			table.vp,
			table.symbolsToDebugLocals(),
			table.codeSeg,
		),
		&table.sym,
		&table.funcsMap).(*bas.Program)
	cls.File = name
	cls.Source = source
	if opt != nil {
		cls.MaxStackSize = opt.MaxStackSize
		cls.Globals = opt.Globals
		cls.Stdout = internal.Or(opt.Stdout, cls.Stdout).(io.Writer)
		cls.Stderr = internal.Or(opt.Stderr, cls.Stderr).(io.Writer)
		cls.Stdin = internal.Or(opt.Stdin, cls.Stdin).(io.Reader)
	}

	coreStack[gi] = bas.ValueOf(cls)
	return cls, err
}
