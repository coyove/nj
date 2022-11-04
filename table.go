package nj

import (
	"fmt"
	"io"
	"math"

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
	top, parent *symTable

	codeSeg internal.Packet

	// variable lookup
	sym       bas.Map   // str -> address: uint16
	maskedSym []bas.Map // str -> address: uint16

	forLoops []*breakLabel

	pendingReleases []uint16

	vp uint16

	constMap bas.Map // value -> address: uint16
	funcsMap bas.Map // func name -> address: uint16

	reusableTmps      bas.Map // address: uint16 -> used: bool
	reusableTmpsArray []uint16

	forwardGoto map[int]*parser.GotoLabel // position to goto label node
	labelPos    map[string]int            // label name to position
}

func newSymTable(opt *LoadOptions) *symTable {
	t := &symTable{
		options: opt,
	}
	return t
}

func (table *symTable) panicnode(node parser.GetLine, msg string, args ...interface{}) {
	who, line := node.GetLine()
	panic(fmt.Sprintf("%q at %s:%d\t", who, table.name, line) + fmt.Sprintf(msg, args...))
}

func (table *symTable) symbolsToDebugLocals() []string {
	x := make([]string, table.vp)
	table.sym.Foreach(func(sym bas.Value, addr *bas.Value) bool {
		x[addr.Int64()] = sym.Str()
		return true
	})
	for _, s := range table.maskedSym {
		s.Foreach(func(sym bas.Value, addr *bas.Value) bool {
			x[addr.Int64()] = sym.Str()
			return true
		})
	}
	return x
}

func (table *symTable) borrowAddress() uint16 {
	if len(table.reusableTmpsArray) > 0 {
		tmp := bas.Int64(int64(table.reusableTmpsArray[0]))
		table.reusableTmpsArray = table.reusableTmpsArray[1:]
		if v, _ := table.reusableTmps.Get(tmp); v.IsFalse() {
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

func (table *symTable) releaseAddr(a interface{}) {
	switch a := a.(type) {
	case []parser.Node:
		for _, n := range a {
			if a, ok := n.(parser.Address); ok {
				table.releaseAddr(uint16(a))
			}
		}
	case []uint16:
		for _, n := range a {
			table.releaseAddr(n)
		}
	case uint16:
		if a == typ.RegA {
			return
		}
		if a > typ.RegLocalMask {
			// We don't free global variables
			return
		}
		if available, ok := table.reusableTmps.Get(bas.Int64(int64(a))); ok && available.IsFalse() {
			table.reusableTmpsArray = append(table.reusableTmpsArray, a)
			table.reusableTmps.Set(bas.Int64(int64(a)), bas.True)
		}
	default:
		internal.ShouldNotHappen()
	}
}

var (
	staticNil   = parser.SNil.Name
	staticTrue  = bas.Str("true")
	staticFalse = bas.Str("false")
	staticThis  = bas.Str("this")
	staticSelf  = bas.Str("self")
	staticA     = parser.Sa.Name
)

func (table *symTable) get(name bas.Value) (uint16, bool) {
	stubLoad := func(name bas.Value) uint16 {
		k, _ := table.sym.Get(name)
		if k.Type() == typ.Number {
			return uint16(k.Int64())
		}
		k = bas.Int64(int64(table.borrowAddressNoReuse()))
		table.sym.Set(name, k)
		return uint16(k.Int64())
	}

	switch name {
	case staticNil:
		return typ.RegNil, true
	case staticTrue:
		return table.loadConst(bas.True), true
	case staticFalse:
		return table.loadConst(bas.False), true
	case staticThis, staticSelf:
		return stubLoad(name), true
	case staticA:
		return typ.RegA, true
	}

	calc := func(k uint16, depth uint16) (uint16, bool) {
		addr := (depth << 15) | (k & typ.RegLocalMask)
		return addr, true
	}

	// Firstly we will iterate local masked symbols,
	// which are local variables inside do-blocks, like "if then .. end" and "do ... end".
	// The rightmost map of this slice is the innermost do-block.
	for i := len(table.maskedSym) - 1; i >= 0; i-- {
		if k, ok := table.maskedSym[i].Get(name); ok {
			return calc(uint16(k.Int64()), 0)
		}
	}

	// Then local variables.
	if k, ok := table.sym.Get(name); ok {
		return calc(uint16(k.Int64()), 0)
	}

	// // If parent exists and parent != top, it means we are inside a closure.
	// for p := table.parent; p != nil && p != table.top; p = p.parent {
	// 	if _, ok := p.sym.Get(name); ok {
	// 		// self := stubLoad(staticSelf)
	// 		table.codeSeg.WriteInst3(typ.OpLoad, self)
	// 	}
	// }

	// Finally top variables.
	if table.top != nil {
		if k, ok := table.top.sym.Get(name); ok {
			return calc(uint16(k.Int64()), 1)
		}
	}

	return typ.RegNil, false
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
	table.maskedSym = append(table.maskedSym, bas.Map{})
}

func (table *symTable) removeMaskedSymTable() {
	table.maskedSym[len(table.maskedSym)-1].Foreach(func(sym bas.Value, addr *bas.Value) bool {
		table.releaseAddr(uint16(addr.Int64()))
		return true
	})
	table.maskedSym = table.maskedSym[:len(table.maskedSym)-1]
}

func (table *symTable) loadConst(v bas.Value) uint16 {
	if table.top != nil {
		return table.top.loadConst(v)
	}
	if i, ok := table.constMap.Get(v); ok {
		return uint16(i.Int64())
	}
	panic("loadConst: shouldn't happen")
}

func (table *symTable) compileOpcode1Node(op byte, n parser.Node) {
	addr, ok := table.compileStaticNode(n)
	if !ok {
		table.codeSeg.WriteInst(op, table.compileNode(n), 0)
	} else {
		table.codeSeg.WriteInst(op, addr, 0)
	}
}

func (table *symTable) compileAtom(n parser.Node, releases *[]uint16) uint16 {
	addr, ok := table.compileStaticNode(n)
	if !ok {
		addr := table.borrowAddress()
		table.codeSeg.WriteInst(typ.OpSet, addr, table.compileNode(n))
		*releases = append(*releases, addr)
		return addr
	}
	return addr
}

func (table *symTable) compileOpcode2Node(op byte, n0, n1 parser.Node) {
	var r []uint16
	var __n0__16, __n0__ = toInt16(n0)
	var __n1__16, __n1__ = toInt16(n1)
	switch {
	case op == typ.OpAdd && __n1__:
		table.codeSeg.WriteInst2Ext(typ.OpExtAdd16, table.compileAtom(n0, &r), __n1__16)
	case op == typ.OpAdd && __n0__:
		table.codeSeg.WriteInst2Ext(typ.OpExtAdd16, table.compileAtom(n1, &r), __n0__16)
	case op == typ.OpSub && __n0__:
		table.codeSeg.WriteInst2Ext(typ.OpExtRSub16, table.compileAtom(n1, &r), __n0__16)
	case op == typ.OpSub && __n1__:
		table.codeSeg.WriteInst2Ext(typ.OpExtAdd16, table.compileAtom(n0, &r), uint16(-int16(__n1__16)))
	case op == typ.OpEq && __n1__:
		table.codeSeg.WriteInst2Ext(typ.OpExtEq16, table.compileAtom(n0, &r), __n1__16)
	case op == typ.OpEq && __n0__:
		table.codeSeg.WriteInst2Ext(typ.OpExtEq16, table.compileAtom(n1, &r), __n0__16)
	case op == typ.OpNeq && __n1__:
		table.codeSeg.WriteInst2Ext(typ.OpExtNeq16, table.compileAtom(n0, &r), __n1__16)
	case op == typ.OpNeq && __n0__:
		table.codeSeg.WriteInst2Ext(typ.OpExtNeq16, table.compileAtom(n1, &r), __n0__16)
	case op == typ.OpLess && __n1__:
		table.codeSeg.WriteInst2Ext(typ.OpExtLess16, table.compileAtom(n0, &r), __n1__16)
	case op == typ.OpLess && __n0__:
		table.codeSeg.WriteInst2Ext(typ.OpExtGreat16, table.compileAtom(n1, &r), __n0__16)
	case op == typ.OpLessEq && __n1__ && int16(__n1__16) <= math.MaxInt16-1:
		table.codeSeg.WriteInst2Ext(typ.OpExtLess16, table.compileAtom(n0, &r), uint16(int16(__n1__16+1)))
	case op == typ.OpLessEq && __n0__ && int16(__n0__16) >= math.MinInt16+1:
		table.codeSeg.WriteInst2Ext(typ.OpExtGreat16, table.compileAtom(n1, &r), uint16(int16(__n0__16-1)))
	case op == typ.OpInc && __n1__:
		table.codeSeg.WriteInst2Ext(typ.OpExtInc16, table.compileAtom(n0, &r), __n1__16)
	default:
		table.codeSeg.WriteInst(op, table.compileAtom(n0, &r), table.compileAtom(n1, &r))
	}
	table.releaseAddr(r)
}

func (table *symTable) compileOpcode3Node(op byte, n0, n1, n2 parser.Node) {
	var r []uint16
	var __n1__16, __n1__ = toInt16(n1)
	switch {
	case op == typ.OpLoad && __n1__:
		table.codeSeg.WriteInst3Ext(typ.OpExtLoad16, table.compileAtom(n0, &r), __n1__16, table.compileAtom(n2, &r))
	case op == typ.OpLoad && isStrNode(n1):
		name := bas.Value(n1.(parser.Primitive)).Str()
		a := table.compileAtom(n0, &r)
		table.codeSeg.WriteInst3Ext(typ.OpExtLoadString, a, uint16(len(name)), table.compileAtom(n2, &r))
		table.codeSeg.Code = append(table.codeSeg.Code, internal.CreateRawBytesInst(name)...)
		table.codeSeg.Code = append(table.codeSeg.Code, typ.Inst{Opcode: typ.OpExt, OpcodeExt: typ.OpExtLoadString})
	case op == typ.OpStore && __n1__:
		table.codeSeg.WriteInst3Ext(typ.OpExtStore16, table.compileAtom(n0, &r), __n1__16, table.compileAtom(n2, &r))
	default:
		table.codeSeg.WriteInst3(op, table.compileAtom(n0, &r), table.compileAtom(n1, &r), table.compileAtom(n2, &r))
	}
	table.releaseAddr(r)
}

func (table *symTable) compileStaticNode(node parser.Node) (uint16, bool) {
	switch v := node.(type) {
	case parser.Address:
		return uint16(v), true
	case parser.Primitive:
		return table.loadConst(bas.Value(v)), true
	case *parser.Symbol:
		idx, ok := table.get(v.Name)
		if !ok {
			if idx := bas.GetTopIndex(v.Name); idx > 0 {
				c := table.borrowAddress()
				table.codeSeg.WriteInst3(typ.OpLoadTop, uint16(idx), typ.RegPhantom, c)
				table.pendingReleases = append(table.pendingReleases, c)
				return c, true
			}
			table.panicnode(v, "symbol not defined")
		}
		return idx, true
	}
	return 0, false
}

func (table *symTable) compileNode(node parser.Node) uint16 {
	if addr, ok := table.compileStaticNode(node); ok {
		return addr
	}

	switch v := node.(type) {
	case *parser.LoadConst:
		table.constMap = v.Table
		table.constMap.Foreach(func(k bas.Value, v *bas.Value) bool {
			addr := int(typ.RegA | table.borrowAddressNoReuse())
			*v = bas.Int(addr)
			return true
		})
		table.funcsMap = v.Funcs
		table.funcsMap.Foreach(func(k bas.Value, v *bas.Value) bool {
			addr := int(table.borrowAddressNoReuse())
			*v = bas.Int(typ.RegA | addr)
			table.sym.Set(k, bas.Int(addr))
			return true
		})
		return typ.RegA
	case *parser.Prog:
		return compileProgBlock(table, v)
	case *parser.Declare:
		return compileDeclare(table, v)
	case *parser.Assign:
		return compileAssign(table, v)
	case parser.Release:
		return compileRelease(table, v)
	case *parser.Unary:
		return compileUnary(table, v)
	case *parser.Binary:
		return compileBinary(table, v)
	case *parser.Tenary:
		return compileTenary(table, v)
	case *parser.And:
		return compileAnd(table, v)
	case *parser.Or:
		return compileOr(table, v)
	case parser.ExprList:
		return compileArray(table, v)
	case parser.ExprAssignList:
		return compileObject(table, v)
	case *parser.GotoLabel:
		return compileGotoLabel(table, v)
	case *parser.Call:
		return compileCall(table, v)
	case *parser.If:
		return compileIf(table, v)
	case *parser.Loop:
		return compileLoop(table, v)
	case *parser.BreakContinue:
		return compileBreakContinue(table, v)
	case *parser.Function:
		return compileFunction(table, v)
	}

	panic("compileNode: shouldn't happen")
}

func (table *symTable) getTopTable() *symTable {
	if table.top != nil {
		return table.top
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

	coreStack := []bas.Value{bas.Nil, bas.Nil}

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

	if opt != nil {
		opt.Globals.Foreach(func(k bas.Value, v *bas.Value) bool { push(k, *v); return true })
	}

	gi := push(bas.Str("Program"), bas.Nil)

	table.vp = uint16(len(coreStack))

	table.compileNode(n)
	table.codeSeg.WriteInst(typ.OpRet, typ.RegA, 0)
	table.patchGoto()

	coreStack = append(coreStack, make([]bas.Value, int(table.vp)-len(coreStack))...)
	table.constMap.Foreach(func(konst bas.Value, addr *bas.Value) bool {
		coreStack[addr.Int64()&typ.RegLocalMask] = konst
		return true
	})

	cls = bas.NewBareProgram(
		coreStack,
		bas.NewBareFunc(
			"main",
			false,
			0,
			table.vp,
			table.symbolsToDebugLocals(),
			nil,
			table.labelPos,
			table.codeSeg,
		),
		&table.sym,
		&table.funcsMap)
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

func toInt16(n parser.Node) (uint16, bool) {
	if a, ok := n.(parser.Primitive); ok && bas.Value(a).IsInt64() {
		a := bas.Value(a).UnsafeInt64()
		if a >= math.MinInt16+1 && a <= math.MaxInt16 {
			return uint16(int16(a)), true // don't take -1<<15 into consideration because we may negate n.
		}
	}
	return 0, false
}

func isStrNode(n parser.Node) bool {
	a, ok := n.(parser.Primitive)
	return ok && bas.Value(a).IsString()
}
