package nj

import (
	"bytes"
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
	global *symTable

	codeSeg internal.Packet

	// variable lookup
	sym       bas.Object   // str -> address: uint16
	maskedSym []bas.Object // str -> address: uint16

	forLoops []*breakLabel

	vp uint16

	constMap bas.Object // value -> address: uint16
	funcsMap bas.Object // func name -> address: uint16

	reusableTmps      bas.Object // address: uint16 -> used: bool
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
	x := &bytes.Buffer{}
	node.Dump(x)
	panic(fmt.Sprintf("%s at %s:%d\t", x.String(), table.name, node.GetLine()) + fmt.Sprintf(msg, args...))
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
			if a, ok := n.(parser.Address); ok {
				table.freeAddr(uint16(a))
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
	staticNil   = bas.Str("nil")
	staticTrue  = bas.Str("true")
	staticFalse = bas.Str("false")
	staticThis  = bas.Str("this")
	staticSelf  = bas.Str("self")
	staticA     = bas.Str("a!")
)

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
	panic("loadConst: shouldn't happen")
}

func (table *symTable) writeInst1(op byte, n parser.Node) {
	addr, ok := table.compileStaticNode(n)
	if !ok {
		table.codeSeg.WriteInst(op, table.compileNode(n), 0)
	} else {
		table.codeSeg.WriteInst(op, addr, 0)
	}
}

func (table *symTable) compileAtom(n parser.Node, tmp *[]uint16) uint16 {
	addr, ok := table.compileStaticNode(n)
	if !ok {
		addr := table.borrowAddress()
		table.codeSeg.WriteInst(typ.OpSet, addr, table.compileNode(n))
		*tmp = append(*tmp, addr)
		return addr
	}
	return addr
}

func (table *symTable) writeInst2(op byte, n0, n1 parser.Node) {
	var tmp []uint16
	i := func(n parser.Node) int16 { return int16(bas.Value(n.(parser.Primitive)).Int64()) }
	i64 := func(n parser.Node) int64 { return bas.Value(n.(parser.Primitive)).Int64() }
	switch {
	case op == typ.OpAdd && parser.IsInt16(n1) > 0:
		table.codeSeg.WriteInst3(typ.OpLinear16, table.compileAtom(n0, &tmp), 1, uint16(i(n1)))
	case op == typ.OpAdd && parser.IsInt16(n0) > 0:
		table.codeSeg.WriteInst3(typ.OpLinear16, table.compileAtom(n1, &tmp), 1, uint16(i(n0)))
	case op == typ.OpMul && parser.IsInt16(n1) > 0:
		table.codeSeg.WriteInst3(typ.OpLinear16, table.compileAtom(n0, &tmp), uint16(i(n1)), 0)
	case op == typ.OpMul && parser.IsInt16(n0) > 0:
		table.codeSeg.WriteInst3(typ.OpLinear16, table.compileAtom(n1, &tmp), uint16(i(n0)), 0)
	case op == typ.OpSub && parser.IsInt16(n1) == 2:
		table.codeSeg.WriteInst3(typ.OpLinear16, table.compileAtom(n0, &tmp), 1, uint16(-i(n1)))
	case op == typ.OpSub && parser.IsInt16(n0) > 0:
		table.codeSeg.WriteInst3(typ.OpLinear16, table.compileAtom(n1, &tmp), 65535 /* -1 */, uint16(i(n0)))
	case (op == typ.OpEq || op == typ.OpNeq) && parser.IsInt16(n1) > 0:
		table.codeSeg.WriteInst3(typ.OpEq16, table.compileAtom(n0, &tmp), uint16(i(n1)), uint16(op))
	case (op == typ.OpEq || op == typ.OpNeq) && parser.IsInt16(n0) > 0:
		table.codeSeg.WriteInst3(typ.OpEq16, table.compileAtom(n1, &tmp), uint16(i(n0)), uint16(op))
	case op == typ.OpLess && parser.IsInt16(n1) > 0:
		table.codeSeg.WriteInst3(typ.OpCmp16, table.compileAtom(n0, &tmp), 1, uint16(i(n1)))
	case op == typ.OpLess && parser.IsInt16(n0) == 2:
		table.codeSeg.WriteInst3(typ.OpCmp16, table.compileAtom(n1, &tmp), 65535, uint16(-i(n0)))
	case op == typ.OpLessEq && parser.IsInt16(n1) > 0 && i64(n1)+1 <= math.MaxInt16:
		table.codeSeg.WriteInst3(typ.OpCmp16, table.compileAtom(n0, &tmp), 1, uint16(int16(i64(n1)+1)))
	case op == typ.OpLessEq && parser.IsInt16(n0) > 0 && -i64(n0)+1 >= math.MinInt16:
		table.codeSeg.WriteInst3(typ.OpCmp16, table.compileAtom(n1, &tmp), 65535, uint16(int16(-i64(n0)+1)))
	case op == typ.OpInc && parser.IsInt16(n1) > 0:
		table.codeSeg.WriteInst(typ.OpInc16, table.compileAtom(n0, &tmp), uint16(i(n1)))
	default:
		table.codeSeg.WriteInst(op, table.compileAtom(n0, &tmp), table.compileAtom(n1, &tmp))
	}
	table.freeAddr(tmp)
}

func (table *symTable) writeInst3(op byte, n0, n1, n2 parser.Node) {
	var tmp []uint16
	table.codeSeg.WriteInst3(op, table.compileAtom(n0, &tmp), table.compileAtom(n1, &tmp), table.compileAtom(n2, &tmp))
	table.freeAddr(tmp)
}

func (table *symTable) compileStaticNode(node parser.Node) (uint16, bool) {
	switch v := node.(type) {
	case parser.Address:
		return uint16(v), true
	case parser.Primitive:
		return table.loadConst(bas.Value(v)), true
	case *parser.Symbol:
		idx, ok := table.get(bas.Str(v.Name))
		if !ok {
			if idx := bas.GetGlobalName(bas.Str(v.Name)); idx > 0 {
				c := table.borrowAddress()
				table.codeSeg.WriteInst3(typ.OpLoadGlobal, uint16(idx), typ.RegPhantom, c)
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
			addr := int(typ.RegGlobalFlag | table.borrowAddressNoReuse())
			*v = bas.Int(addr)
			return true
		})
		table.funcsMap = v.Funcs
		table.funcsMap.Foreach(func(k bas.Value, v *bas.Value) bool {
			addr := int(table.borrowAddressNoReuse())
			*v = bas.Int(typ.RegGlobalFlag | addr)
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
	case *parser.Bitwise:
		return compileBitwise(table, v)
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

	coreStack := []bas.Value{bas.Nil}

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

	if opt != nil && opt.Globals != nil {
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
