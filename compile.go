package nj

import (
	"strconv"
	"strings"
	"unsafe"

	"github.com/coyove/nj/bas"
	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/parser"
	"github.com/coyove/nj/typ"
)

var _nodeRegA = parser.Addr(typ.RegA)

// [prog expr1 expr2 ...]
func (table *symTable) compileChain(chain parser.Node) uint16 {
	doblock := chain.Nodes()[0].Sym() == typ.ADoBlock

	if doblock {
		table.addMaskedSymTable()
	}

	yx := typ.RegA
	for i, a := range chain.Nodes() {
		if i == 0 {
			continue
		}
		_, isStatic := table.compileStaticNode(a)
		yx = table.compileNode(a)
		if isStatic {
			// e.g.: [prog "a string"], we will transform it into:
			//       [prog [set $a "a string"]]
			if yx != typ.RegA {
				table.codeSeg.WriteInst(typ.OpSet, typ.RegA, yx)
			}
		}
	}

	if doblock {
		table.removeMaskedSymTable()
	}

	return yx
}

func (table *symTable) compileSetMove(nodes []parser.Node) uint16 {
	dest := nodes[1].Sym()
	destAddr, declared := table.get(dest)
	if nodes[0].Sym() == typ.AMove {
		// a = b
		if !declared {
			// a is not declared yet
			destAddr = table.borrowAddress()

			// Do not use t.put() because it may put the symbol into masked tables
			// e.g.: do a = 1 end
			table.sym[dest] = &typ.Symbol{Address: destAddr}
		}
	} else {
		// local a = b
		destAddr = table.borrowAddress()
		defer table.put(dest, destAddr) // execute in defer in case of: a = 1 do local a = a end
	}

	srcAddr := table.compileNode(nodes[2])
	table.codeSeg.WriteInst(typ.OpSet, destAddr, srcAddr)
	table.codeSeg.WriteLineNum(nodes[0].Line())
	return destAddr
}

// writeInst3 accepts 3 arguments at most, 2 arguments will be encoded into opCode itself, the 3rd one will be in typ.RegA
func (table *symTable) writeInst3(bop byte, atoms []parser.Node) uint16 {
	// first atom: the splitInst Name, tail atoms: the args
	if len(atoms) > 4 {
		panic("DEBUG: too many arguments")
	}

	atoms = append([]parser.Node{}, atoms...) // duplicate

	if bop == typ.OpStore {
		table.collapse(atoms[1:], true)

		// (atoms    1      2    3 )
		// (store subject value key) subject => opa, key => $a, value => opb

		for i := 1; i <= 2; i++ { // subject and value shouldn't use typ.RegA
			if atoms[i].Type() == parser.ADDR && uint16(atoms[i].Int()) == typ.RegA {
				n := parser.Addr(table.borrowAddress())
				table.writeInst(typ.OpSet, n, _nodeRegA)
				atoms[i] = n
			}
		}

		// We would love to see 'key' using typ.RegA, in this case writeInst will just omit it
		table.writeInst(typ.OpSet, _nodeRegA, atoms[3])
		table.writeInst(typ.OpStore, atoms[1], atoms[2])
		table.freeAddr(atoms[1:])
		return typ.RegA
	}

	table.collapse(atoms[1:], true)

	switch bop {
	case typ.OpNot, typ.OpRet, typ.OpBitNot, typ.OpLen:
		// unary splitInst
		table.writeInst(bop, atoms[1], parser.Node{})
	default:
		// binary splitInst
		table.writeInst(bop, atoms[1], atoms[2])
		table.freeAddr(atoms[1:])
	}

	return typ.RegA
}

func (table *symTable) compileOperator(atoms []parser.Node) uint16 {
	head := atoms[0].Sym()
	op, ok := typ.NodeOpcode[head]
	if !ok {
		internal.Panic("DEBUG invalid symbol: %v", atoms[0])
	}
	yx := table.writeInst3(op, atoms)
	if p := atoms[0].Line(); p > 0 {
		table.codeSeg.WriteLineNum(p)
	}
	return yx
}

// [and a b] => $a = a if not a then goto out else $a = b end ::out::
// [or a b]  => $a = a if not a then $a = b end
func (table *symTable) compileAndOr(atoms []parser.Node) uint16 {
	table.writeInst(typ.OpSet, _nodeRegA, atoms[1])

	if atoms[0].Sym() == (typ.AOr) {
		table.codeSeg.WriteJmpInst(typ.OpIfNot, 1)
		table.codeSeg.WriteJmpInst(typ.OpJmp, 0)
		part1 := table.codeSeg.Len()

		table.writeInst(typ.OpSet, _nodeRegA, atoms[2])
		part2 := table.codeSeg.Len()

		table.codeSeg.Code[part1-1] = typ.JmpInst(typ.OpJmp, part2-part1)
	} else {
		table.codeSeg.WriteJmpInst(typ.OpIfNot, 0)
		part1 := table.codeSeg.Len()

		table.writeInst(typ.OpSet, _nodeRegA, atoms[2])
		part2 := table.codeSeg.Len()

		table.codeSeg.Code[part1-1] = typ.JmpInst(typ.OpIfNot, part2-part1)
	}
	table.codeSeg.WriteLineNum(atoms[0].Line())
	return typ.RegA
}

// [if condition [true-chain ...] [false-chain ...]]
func (table *symTable) compileIf(atoms []parser.Node) uint16 {
	condition := atoms[1]
	trueBranch, falseBranch := atoms[2], atoms[3]

	condyx := table.compileNode(condition)

	table.addMaskedSymTable()

	if condyx != typ.RegA {
		table.codeSeg.WriteInst(typ.OpSet, typ.RegA, condyx)
	}

	table.codeSeg.WriteJmpInst(typ.OpIfNot, 0)
	table.codeSeg.WriteLineNum(atoms[0].Line())
	init := table.codeSeg.Len()

	table.compileNode(trueBranch)
	part1 := table.codeSeg.Len()

	table.codeSeg.WriteJmpInst(typ.OpJmp, 0)

	table.compileNode(falseBranch)
	part2 := table.codeSeg.Len()

	table.removeMaskedSymTable()

	if len(falseBranch.Nodes()) > 0 {
		table.codeSeg.Code[init-1] = typ.JmpInst(typ.OpIfNot, part1-init+1)
		table.codeSeg.Code[part1] = typ.JmpInst(typ.OpJmp, part2-part1-1)
	} else {
		// The last inst is used to skip the false branch, since we don't have one, we don't need this jmp
		table.codeSeg.TruncLast()
		table.codeSeg.Code[init-1] = typ.JmpInst(typ.OpIfNot, part1-init)
	}
	return typ.RegA
}

// [list [a, b, c, ...]]
func (table *symTable) compileList(nodes []parser.Node) uint16 {
	table.collapse(nodes[1].Nodes(), true)
	if nodes[0].Sym() == typ.AArray {
		for _, x := range nodes[1].Nodes() {
			table.writeInst(typ.OpPush, x, parser.Node{})
		}
		table.codeSeg.WriteInst(typ.OpCreateArray, 0, 0)
	} else {
		n := nodes[1].Nodes()
		for i := 0; i < len(n); i += 2 {
			table.writeInst(typ.OpPush, n[i], parser.Node{})
			table.writeInst(typ.OpPush, n[i+1], parser.Node{})
		}
		table.codeSeg.WriteInst(typ.OpCreateObject, 0, 0)
	}
	return typ.RegA
}

// [call callee [args ...]]
func (table *symTable) compileCall(nodes []parser.Node) uint16 {
	tmp := append([]parser.Node{nodes[1]}, nodes[2].Nodes()...)
	isVariadic := false
	if last := &tmp[len(tmp)-1]; len(last.Nodes()) == 2 && last.Nodes()[0].Sym() == typ.AUnpack {
		// [call callee [a b .. [unpack vararg]]]
		*last = last.Nodes()[1]
		table.collapse(tmp, true)
		for i := 1; i < len(tmp)-1; i++ {
			table.writeInst(typ.OpPush, tmp[i], parser.Addr(0))
		}
		table.writeInst(typ.OpPushUnpack, tmp[len(tmp)-1], parser.Addr(0))
		isVariadic = true
	} else {
		table.collapse(tmp, true)
		for i := 1; i < len(tmp)-1; i++ {
			table.writeInst(typ.OpPush, tmp[i], parser.Addr(0))
		}
	}

	op := byte(typ.OpCall)
	if nodes[0].Sym() == typ.ATailCall {
		op = typ.OpTailCall
	}
	if len(tmp) == 1 || isVariadic {
		table.writeInst(op, tmp[0], parser.Addr(typ.RegPhantom))
	} else {
		table.writeInst(op, tmp[0], tmp[len(tmp)-1])
	}

	table.codeSeg.WriteLineNum(nodes[0].Line())
	table.freeAddr(tmp)
	return typ.RegA
}

// [function name [paramlist] [chain ...] docstring]
func (table *symTable) compileFunction(atoms []parser.Node) uint16 {
	params := atoms[2]
	newtable := newSymTable(table.options)
	newtable.name = table.name
	newtable.codeSeg.Pos.Name = table.name
	if table.global == nil {
		newtable.global = table
	} else {
		newtable.global = table.global
	}
	newtable.parent = table

	varargIdx := -1
	for i, p := range params.Nodes() {
		n := p.Sym()
		if len(p.Nodes()) == 2 && p.Nodes()[0].Sym() == typ.AUnpack {
			n = p.Nodes()[1].Sym()
			varargIdx = i
		}
		if _, ok := newtable.sym[n]; ok {
			internal.Panic("%v: duplicated parameter: %q", atoms[1], n)
		}
		newtable.put(n, uint16(i))
	}

	ln := len(newtable.sym)
	if ln > 255 {
		internal.Panic("%v: too many parameters, 255 at most", atoms[1])
	}

	newtable.vp = uint16(len(newtable.sym))
	newtable.compileNode(atoms[3])
	newtable.patchGoto()

	if a := newtable.sym["this"]; a != nil {
		newtable.codeSeg.Code = append([]typ.Inst{
			{Opcode: typ.OpSet, A: a.Address, B: int32(typ.RegA)},
		}, newtable.codeSeg.Code...)
	}

	code := newtable.codeSeg
	code.WriteInst(typ.OpRet, table.loadK(nil), 0)
	// code.writeInst(typ.OpRet, typ.RegA, 0)

	cls := &bas.Function{}
	cls.Variadic = varargIdx >= 0
	cls.NumParams = uint16(len(params.Nodes()))
	cls.Name = atoms[1].Sym()
	cls.DocString = atoms[4].Str()
	cls.StackSize = newtable.vp
	cls.CodeSeg = code
	cls.Locals = newtable.symbolsToDebugLocals()

	var loadFuncIndex uint16
	obj := bas.NewObject(0)
	obj.SetPrototype(bas.FuncProto)
	internal.SetObjFun(unsafe.Pointer(obj), unsafe.Pointer(cls))
	if table.global != nil {
		x := table.global
		loadFuncIndex = uint16(len(x.funcs))
		x.funcs = append(x.funcs, obj)
	} else {
		loadFuncIndex = uint16(len(table.funcs))
		table.funcs = append(table.funcs, obj)
	}
	table.codeSeg.WriteInst(typ.OpLoadFunc, loadFuncIndex, 0)
	if strings.HasPrefix(cls.Name, "<lambda") {
		cls.Name = cls.Name[:len(cls.Name)-1] + "-" + strconv.Itoa(int(loadFuncIndex)) + ">"
	}
	table.codeSeg.WriteLineNum(atoms[0].Line())
	return typ.RegA
}

// [break|continue]
func (table *symTable) compileBreak(atoms []parser.Node) uint16 {
	if len(table.forLoops) == 0 {
		internal.Panic("%v: outside loop", atoms[0])
	}
	bl := table.forLoops[len(table.forLoops)-1]
	if atoms[0].Sym() == typ.AContinue {
		table.compileNode(bl.continueNode)
		table.codeSeg.WriteJmpInst(typ.OpJmp, bl.continueGoto-len(table.codeSeg.Code)-1)
	} else {
		bl.labelPos = append(bl.labelPos, table.codeSeg.Len())
		table.codeSeg.WriteJmpInst(typ.OpJmp, 0)
	}
	return typ.RegA
}

// [loop [chain ...]]
func (table *symTable) compileWhile(atoms []parser.Node) uint16 {
	init := table.codeSeg.Len()
	breaks := &breakLabel{
		continueNode: atoms[2],
		continueGoto: init,
	}

	table.forLoops = append(table.forLoops, breaks)
	table.addMaskedSymTable()
	table.compileNode(atoms[1])
	table.removeMaskedSymTable()
	table.forLoops = table.forLoops[:len(table.forLoops)-1]

	table.codeSeg.WriteJmpInst(typ.OpJmp, -(table.codeSeg.Len()-init)-1)
	for _, idx := range breaks.labelPos {
		table.codeSeg.Code[idx] = typ.JmpInst(typ.OpJmp, table.codeSeg.Len()-idx-1)
	}
	return typ.RegA
}

func (table *symTable) compileGoto(atoms []parser.Node) uint16 {
	label := atoms[1].Sym()
	if atoms[0].Sym() == typ.ALabel { // :: label ::
		table.labelPos[label] = table.codeSeg.Len()
	} else { // goto label
		if pos, ok := table.labelPos[label]; ok {
			table.codeSeg.WriteJmpInst(typ.OpJmp, pos-(table.codeSeg.Len()+1))
		} else {
			table.codeSeg.WriteJmpInst(typ.OpJmp, 0)
			table.forwardGoto[table.codeSeg.Len()-1] = label
		}
	}
	return typ.RegA
}

func (table *symTable) patchGoto() {
	code := table.codeSeg.Code
	for i, l := range table.forwardGoto {
		pos, ok := table.labelPos[l]
		if !ok {
			internal.Panic("label %q not found", l)
		}
		code[i] = typ.JmpInst(typ.OpJmp, pos-(i+1))
	}
	for i, c := range code {
		if c.Opcode == typ.OpJmp && c.B != 0 {
			dest := int32(i) + c.B + 1
			for int(dest) < len(code) {
				if c2 := code[dest]; c2.Opcode == typ.OpJmp && c2.B != 0 {
					dest += c2.B + 1
					continue
				}
				break
			}
			code[i].B = dest - int32(i) - 1
		}
	}
}

func (table *symTable) compileFreeAddr(atoms []parser.Node) uint16 {
	for i := 1; i < len(atoms); i++ {
		s := atoms[i].Sym()
		yx, _ := table.get(s)
		table.freeAddr(yx)
		if len(table.maskedSym) > 0 {
			delete(table.maskedSym[len(table.maskedSym)-1], s)
		} else {
			delete(table.sym, s)
		}
	}
	return typ.RegA
}

// collapse will accept a list of expressions, for each of them,
// it will be collapsed into a temp variable and be replaced with a ADR node,
// the last expression will be collapsed and not using a temp variable if optLast is true.
func (table *symTable) collapse(nodes []parser.Node, optLast bool) {
	var lastCompound struct {
		n parser.Node
		i int
	}

	for i, atom := range nodes {
		if !atom.Valid() {
			break
		}

		if atom.Type() == parser.NODES {
			yx := table.compileNodeInto(atom, true, 0)
			nodes[i] = parser.Addr(yx)

			lastCompound.n = atom
			lastCompound.i = i
		}
	}

	if lastCompound.n.Valid() {
		if optLast {
			i := table.codeSeg.LastInst()
			if i.Opcode == typ.OpSet {
				table.codeSeg.TruncLast()
				table.freeAddr(i.A)
				nodes[lastCompound.i] = parser.Addr(uint16(i.B))
			}
		}
	}
}
