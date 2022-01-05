package nj

import (
	"unsafe"

	"github.com/coyove/nj/bas"
	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/parser"
	"github.com/coyove/nj/typ"
)

type breakLabel struct {
	continueNode parser.Node
	continueGoto int
	labelPos     []int
}

// symTable is responsible for recording the state of compilation
type symTable struct {
	name    string
	options *bas.Environment

	global *symTable
	parent *symTable

	codeSeg internal.Packet

	// toplevel symtable
	funcs []*bas.Object

	// variable lookup
	sym       map[string]*typ.Symbol
	maskedSym []map[string]*typ.Symbol

	forLoops []*breakLabel

	vp uint16

	collectConstMode bool
	constMap         map[interface{}]uint16

	reusableTmps      map[uint16]bool
	reusableTmpsArray []uint16

	forwardGoto map[int]string
	labelPos    map[string]int
}

func newSymTable(opt *bas.Environment) *symTable {
	t := &symTable{
		sym:          make(map[string]*typ.Symbol),
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
		x[info.Address] = sym
	}
	return x
}

func (table *symTable) borrowAddress() uint16 {
	if len(table.reusableTmpsArray) > 0 {
		tmp := table.reusableTmpsArray[0]
		table.reusableTmpsArray = table.reusableTmpsArray[1:]
		if !table.reusableTmps[tmp] {
			internal.ShouldNotHappen()
		}
		table.reusableTmps[tmp] = false
		return tmp
	}
	if table.vp > typ.RegMaxAddress {
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
		if available, existed := table.reusableTmps[a]; existed && !available {
			table.reusableTmpsArray = append(table.reusableTmpsArray, a)
			table.reusableTmps[a] = true
		}

	default:
		internal.ShouldNotHappen()
	}
}

func (table *symTable) get(name string) (uint16, bool) {
	switch name {
	case "nil":
		return table.loadK(nil), true
	case "true":
		return table.loadK(true), true
	case "false":
		return table.loadK(false), true
	case "this":
		if k, ok := table.sym[name]; ok {
			return k.Address, true
		}
		k := &typ.Symbol{Address: table.borrowAddress()}
		table.sym["this"] = k
		return k.Address, true
	case "$a":
		return typ.RegA, true
	}

	calc := func(k *typ.Symbol, depth uint16) (uint16, bool) {
		addr := (depth << 15) | (k.Address & typ.RegLocalMask)
		return addr, true
	}

	// Firstly we will iterate local masked symbols,
	// which are local variables inside do-blocks, like "if then .. end" and "do ... end".
	// The rightmost map of this slice is the innermost do-block
	for i := len(table.maskedSym) - 1; i >= 0; i-- {
		m := table.maskedSym[i]
		if k, ok := m[name]; ok {
			return calc(k, 0)
		}
	}

	// Then local variables
	if k, ok := table.sym[name]; ok {
		return calc(k, 0)
	}

	// Finally global variables
	if table.global != nil {
		if k, ok := table.global.sym[name]; ok {
			return calc(k, 1)
		}
	}

	return table.loadK(nil), false
}

func (table *symTable) put(name string, addr uint16) {
	if addr == typ.RegA {
		panic("DEBUG: put $a?")
	}
	sym := &typ.Symbol{
		Address: addr,
	}
	if len(table.maskedSym) > 0 {
		table.maskedSym[len(table.maskedSym)-1][name] = sym
	} else {
		table.sym[name] = sym
	}
}

func (table *symTable) addMaskedSymTable() {
	table.maskedSym = append(table.maskedSym, map[string]*typ.Symbol{})
}

func (table *symTable) removeMaskedSymTable() {
	last := table.maskedSym[len(table.maskedSym)-1]
	for _, k := range last {
		table.freeAddr(k.Address)
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
		internal.ShouldNotHappen(v)
	}

	idx := typ.RegGlobalFlag | table.borrowAddress()
	table.constMap[v] = idx
	return idx
}

func (table *symTable) writeInst(op byte, n0, n1 parser.Node) {
	var tmp []uint16
	getAddr := func(n parser.Node, intoNewAddr bool) uint16 {
		switch n.Type() {
		case parser.NODES:
			addr := table.compileNodeInto(n, intoNewAddr, typ.RegA)
			tmp = append(tmp, addr)
			return addr
		default:
			addr, ok := table.compileStaticNode(n)
			if !ok {
				internal.ShouldNotHappen(n)
			}
			return addr
		}
	}

	if !n0.Valid() {
		table.codeSeg.WriteInst(op, 0, 0)
		return
	}

	n0a := getAddr(n0, n1.Valid())
	if !n1.Valid() {
		table.codeSeg.WriteInst(op, n0a, 0)
		table.freeAddr(tmp)
		return
	}

	n1a := getAddr(n1, true)
	table.codeSeg.WriteInst(op, n0a, n1a)
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

	table.codeSeg.WriteInst(typ.OpSet, yx, newYX)
	return yx
}

func (table *symTable) compileStaticNode(node parser.Node) (uint16, bool) {
	switch node.Type() {
	case parser.ADDR:
		return uint16(node.Int()), true
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
		return typ.RegA
	}

	name := nodes[0].Sym()
	var yx uint16
	switch name {
	case typ.ADoBlock, typ.ABegin:
		yx = table.compileChain(node)
	case typ.ASet, typ.AMove:
		yx = table.compileSetMove(nodes)
	case typ.AIf:
		yx = table.compileIf(nodes)
	case typ.AFor:
		yx = table.compileWhile(nodes)
	case typ.ABreak, typ.AContinue:
		yx = table.compileBreak(nodes)
	case typ.ACall, typ.ATailCall:
		yx = table.compileCall(nodes)
	case typ.AArray, typ.AObject:
		yx = table.compileList(nodes)
	case typ.AOr, typ.AAnd:
		yx = table.compileAndOr(nodes)
	case typ.AFunc:
		yx = table.compileFunction(nodes)
	case typ.AFreeAddr:
		yx = table.compileFreeAddr(nodes)
	case typ.AGoto, typ.ALabel:
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

func compileNodeTopLevel(name, source string, n parser.Node, env *bas.Environment) (cls *bas.Program, err error) {
	defer internal.CatchError(&err)

	table := newSymTable(env)
	table.collectConstMode = true
	table.name = name
	table.codeSeg.Pos.Name = name
	coreStack := bas.NewEnv()

	// Load nil first to ensure its address == 0
	table.loadK(nil)
	coreStack.Push(bas.Nil)

	push := func(k string, v bas.Value) uint16 {
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

	bas.Globals.Foreach(func(k bas.Value, v *bas.Value) bool { push(k.String(), *v); return true })

	if env != nil && env.Globals != nil {
		env.Globals.Foreach(func(k bas.Value, v *bas.Value) bool { push(k.String(), *v); return true })
	}

	gi := push("PROGRAM", bas.Nil)
	push("COMPILE_OPTIONS", bas.ValueOf(env))
	push("SOURCE_CODE", bas.Str(source))

	table.vp = uint16(coreStack.Size())

	// Find and fill consts
	table.loadK(true)
	table.loadK(false)
	table.collectConsts(n)
	table.collectConstMode = false

	table.compileNode(n)
	table.codeSeg.WriteInst(typ.OpRet, typ.RegA, 0)
	table.patchGoto()

	internal.GrowEnvStack(unsafe.Pointer(coreStack), int(table.vp))
	for k, stackPos := range table.constMap {
		switch k := k.(type) {
		case float64:
			coreStack.Set(int(stackPos), bas.Float64(k))
		case int64:
			coreStack.Set(int(stackPos), bas.Int64(k))
		case string:
			coreStack.Set(int(stackPos), bas.Str(k))
		case bool:
			coreStack.Set(int(stackPos), bas.Bool(k))
		case nil:
			coreStack.Set(int(stackPos), bas.Nil)
		default:
			internal.ShouldNotHappen()
		}
	}

	cls = bas.NewProgram(coreStack, &bas.Function{
		Name:      "main",
		CodeSeg:   table.codeSeg,
		StackSize: table.vp,
		Locals:    table.symbolsToDebugLocals(),
	}, table.sym, table.funcs, env)
	coreStack.Set(int(gi), bas.ValueOf(cls))
	return cls, err
}
