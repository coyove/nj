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

var nodeRegA = parser.Addr(typ.RegA)

// [prog expr1 expr2 ...]
func compileProgBlock(table *symTable, nodes []parser.Node) uint16 {
	doblock := nodes[0].Value == parser.SDoBlock.Value
	if doblock {
		table.addMaskedSymTable()
	}

	yx := typ.RegA
	for _, a := range nodes[1:] {
		switch a.Type() {
		case parser.ADDR, parser.STR, parser.FLOAT, parser.INT, parser.SYM:
			// e.g.: [prog "a string"] will be transformed into: [prog [set $a "a string"]]
			yx = table.compileNode(a)
			table.codeSeg.WriteInst(typ.OpSet, typ.RegA, yx)
		default:
			yx = table.compileNode(a)
		}
	}

	if doblock {
		table.removeMaskedSymTable()
	}
	return yx
}

func compileSetMove(table *symTable, nodes []parser.Node) uint16 {
	dest := nodes[1].Value
	destAddr, declared := table.get(dest)
	if nodes[0].Value == parser.SMove.Value {
		// a = b
		if !declared {
			// a is not declared yet
			destAddr = table.borrowAddress()

			// Do not use t.put() because it may put the symbol into masked tables
			// e.g.: do a = 1 end
			table.sym.Set(dest, bas.Int64(int64(destAddr)))
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
func (table *symTable) writeInst3(bop byte, nodes []parser.Node) uint16 {
	// first atom: the splitInst Name, tail atoms: the args
	if len(nodes) > 4 {
		internal.ShouldNotHappen(nodes)
	}

	nodes = append([]parser.Node{}, nodes...) // duplicate

	if bop == typ.OpStore {
		table.collapse(nodes[1:], true)

		// (node     1      2    3 )
		// (store subject value key) subject => opa, key => $a, value => opb

		for i := 1; i <= 2; i++ { // subject and value shouldn't use typ.RegA
			if nodes[i].Type() == parser.ADDR && uint16(nodes[i].Int()) == typ.RegA {
				n := parser.Addr(table.borrowAddress())
				table.writeInst2(typ.OpSet, n, nodeRegA)
				nodes[i] = n
			}
		}

		// We would love to see 'key' using typ.RegA, in this case writeInst will just omit it
		table.writeInst2(typ.OpSet, nodeRegA, nodes[3])
		table.writeInst2(typ.OpStore, nodes[1], nodes[2])
		table.freeAddr(nodes[1:])
		return typ.RegA
	}

	table.collapse(nodes[1:], true)

	switch bop {
	case typ.OpNot, typ.OpRet, typ.OpBitNot, typ.OpLen:
		// unary splitInst
		table.writeInst1(bop, nodes[1])
	default:
		// binary splitInst
		table.writeInst2(bop, nodes[1], nodes[2])
		table.freeAddr(nodes[1:])
	}

	return typ.RegA
}

func makeOPCompiler(op byte) func(table *symTable, nodes []parser.Node) uint16 {
	return func(table *symTable, nodes []parser.Node) uint16 {
		yx := table.writeInst3(op, nodes)
		if p := nodes[0].Line(); p > 0 {
			table.codeSeg.WriteLineNum(p)
		}
		return yx
	}
}

// [and a b] => $a = a if not a then goto out else $a = b end ::out::
// [or a b]  => $a = a if not a then $a = b end
func compileAndOr(table *symTable, nodes []parser.Node) uint16 {
	table.writeInst2(typ.OpSet, nodeRegA, nodes[1])

	if nodes[0].Value == parser.SOr.Value {
		table.codeSeg.WriteJmpInst(typ.OpIfNot, 1)
		table.codeSeg.WriteJmpInst(typ.OpJmp, 0)
		part1 := table.codeSeg.Len()

		table.writeInst2(typ.OpSet, nodeRegA, nodes[2])
		part2 := table.codeSeg.Len()

		table.codeSeg.Code[part1-1] = typ.JmpInst(typ.OpJmp, part2-part1)
	} else {
		table.codeSeg.WriteJmpInst(typ.OpIfNot, 0)
		part1 := table.codeSeg.Len()

		table.writeInst2(typ.OpSet, nodeRegA, nodes[2])
		part2 := table.codeSeg.Len()

		table.codeSeg.Code[part1-1] = typ.JmpInst(typ.OpIfNot, part2-part1)
	}
	table.codeSeg.WriteLineNum(nodes[0].Line())
	return typ.RegA
}

// [if condition [true-chain ...] [false-chain ...]]
func compileIf(table *symTable, atoms []parser.Node) uint16 {
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

// [object [k1, v1, k2, v2, ...]]
func compileObject(table *symTable, nodes []parser.Node) uint16 {
	table.collapse(nodes[1].Nodes(), true)
	n := nodes[1].Nodes()
	for i := 0; i < len(n); i += 2 {
		table.writeInst1(typ.OpPush, n[i])
		table.writeInst1(typ.OpPush, n[i+1])
	}
	table.codeSeg.WriteInst(typ.OpCreateObject, 0, 0)
	return typ.RegA
}

// [array [a, b, c, ...]]
func compileArray(table *symTable, nodes []parser.Node) uint16 {
	table.collapse(nodes[1].Nodes(), true)
	for _, x := range nodes[1].Nodes() {
		table.writeInst1(typ.OpPush, x)
	}
	table.codeSeg.WriteInst(typ.OpCreateArray, 0, 0)
	return typ.RegA
}

// [call callee [args ...]]
func compileCall(table *symTable, nodes []parser.Node) uint16 {
	tmp := append([]parser.Node{nodes[1]}, nodes[2].Nodes()...)
	isVariadic := false
	if last := &tmp[len(tmp)-1]; len(last.Nodes()) == 2 && last.Nodes()[0].Value == parser.SUnpack.Value {
		// [call callee [a b .. [unpack vararg]]]
		*last = last.Nodes()[1]
		table.collapse(tmp, true)
		for i := 1; i < len(tmp)-1; i++ {
			table.writeInst1(typ.OpPush, tmp[i])
		}
		table.writeInst1(typ.OpPushUnpack, tmp[len(tmp)-1])
		isVariadic = true
	} else {
		table.collapse(tmp, true)
		for i := 1; i < len(tmp)-1; i++ {
			table.writeInst1(typ.OpPush, tmp[i])
		}
	}

	op := byte(typ.OpCall)
	if nodes[0].Value == parser.STailCall.Value {
		op = typ.OpTailCall
	}
	if len(tmp) == 1 || isVariadic {
		table.writeInst2(op, tmp[0], parser.Addr(typ.RegPhantom))
	} else {
		table.writeInst2(op, tmp[0], tmp[len(tmp)-1])
	}

	table.codeSeg.WriteLineNum(nodes[0].Line())
	table.freeAddr(tmp)
	return typ.RegA
}

// [function name [paramlist] [chain ...] docstring]
func compileFunction(table *symTable, nodes []parser.Node) uint16 {
	params := nodes[2]
	newtable := newSymTable(table.options)
	newtable.name = table.name
	newtable.codeSeg.Pos.Name = table.name
	newtable.global = table.getGlobal()
	newtable.parent = table

	varargIdx := -1
	for i, p := range params.Nodes() {
		n := p
		if len(p.Nodes()) == 2 && p.Nodes()[0].Value == parser.SUnpack.Value {
			n = p.Nodes()[1]
			varargIdx = i
		}
		if newtable.sym.Contains(n.Value, false) {
			internal.Panic("%v: duplicated parameter: %q", nodes[1], n)
		}
		newtable.put(n.Value, uint16(i))
	}

	if ln := newtable.sym.Len(); ln > 255 {
		internal.Panic("%v: too many parameters, 255 at most", nodes[1])
	}

	newtable.vp = uint16(newtable.sym.Len())
	newtable.compileNode(nodes[3])
	newtable.patchGoto()

	if a := newtable.sym.Prop("this"); a != bas.Nil {
		newtable.codeSeg.Code = append([]typ.Inst{
			{Opcode: typ.OpSet, A: uint16(a.Int64()), B: int32(typ.RegA)},
		}, newtable.codeSeg.Code...)
	}

	code := newtable.codeSeg
	code.WriteInst(typ.OpRet, typ.RegGlobalFlag, 0) // return nil

	cls := &bas.Function{}
	cls.Variadic = varargIdx >= 0
	cls.NumParams = uint16(len(params.Nodes()))
	cls.Name = nodes[1].Sym()
	cls.DocString = nodes[4].Str()
	cls.StackSize = newtable.vp
	cls.CodeSeg = code
	cls.Locals = newtable.symbolsToDebugLocals()

	obj := bas.NewObject(0)
	obj.SetPrototype(bas.Proto.Func)
	internal.SetObjFun(unsafe.Pointer(obj), unsafe.Pointer(cls))

	funcIdx := uint16(len(table.getGlobal().funcs))
	table.getGlobal().funcs = append(table.getGlobal().funcs, obj)
	table.codeSeg.WriteInst(typ.OpLoadFunc, funcIdx, 0)

	if strings.HasPrefix(cls.Name, "<lambda") {
		cls.Name = cls.Name[:len(cls.Name)-1] + "-" + strconv.Itoa(int(funcIdx)) + ">"
	}

	table.codeSeg.WriteLineNum(nodes[0].Line())
	return typ.RegA
}

// [break|continue]
func compileBreak(table *symTable, atoms []parser.Node) uint16 {
	if len(table.forLoops) == 0 {
		internal.Panic("%v: outside loop", atoms[0])
	}
	bl := table.forLoops[len(table.forLoops)-1]
	if atoms[0].Value == parser.SContinue.Value {
		table.compileNode(bl.continueNode)
		table.codeSeg.WriteJmpInst(typ.OpJmp, bl.continueGoto-len(table.codeSeg.Code)-1)
	} else {
		bl.labelPos = append(bl.labelPos, table.codeSeg.Len())
		table.codeSeg.WriteJmpInst(typ.OpJmp, 0)
	}
	return typ.RegA
}

// [loop [chain ...]]
func compileWhile(table *symTable, nodes []parser.Node) uint16 {
	init := table.codeSeg.Len()
	breaks := &breakLabel{
		continueNode: nodes[2],
		continueGoto: init,
	}

	table.forLoops = append(table.forLoops, breaks)
	table.addMaskedSymTable()
	table.compileNode(nodes[1])
	table.removeMaskedSymTable()
	table.forLoops = table.forLoops[:len(table.forLoops)-1]

	table.codeSeg.WriteJmpInst(typ.OpJmp, -(table.codeSeg.Len()-init)-1)
	for _, idx := range breaks.labelPos {
		table.codeSeg.Code[idx] = typ.JmpInst(typ.OpJmp, table.codeSeg.Len()-idx-1)
	}
	return typ.RegA
}

func compileLabel(table *symTable, nodes []parser.Node) uint16 {
	table.labelPos.Set(nodes[1].Value, bas.Int(table.codeSeg.Len()))
	return typ.RegA
}

func compileGoto(table *symTable, nodes []parser.Node) uint16 {
	label := nodes[1].Value
	if pos := table.labelPos.Get(label); pos != bas.Nil {
		table.codeSeg.WriteJmpInst(typ.OpJmp, pos.Int()-(table.codeSeg.Len()+1))
	} else {
		table.codeSeg.WriteJmpInst(typ.OpJmp, 0)
		table.forwardGoto.Set(bas.Int(table.codeSeg.Len()-1), label)
	}
	return typ.RegA
}

func (table *symTable) patchGoto() {
	code := table.codeSeg.Code
	table.forwardGoto.Foreach(func(i bas.Value, l *bas.Value) bool {
		pos := table.labelPos.Get(*l)
		if pos == bas.Nil {
			internal.Panic("label %q not found", l.Str())
		}
		code[i.Int()] = typ.JmpInst(typ.OpJmp, pos.Int()-(i.Int()+1))
		return true
	})
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

func compileFreeAddr(table *symTable, nodes []parser.Node) uint16 {
	for i := 1; i < len(nodes); i++ {
		s := nodes[i].Value
		yx, _ := table.get(s)
		table.freeAddr(yx)
		t := table.sym
		if len(table.maskedSym) > 0 {
			t = table.maskedSym[len(table.maskedSym)-1]
		}
		if !t.Contains(s, false) {
			internal.ShouldNotHappen(nodes)
		}
		t.Delete(s)
	}
	return typ.RegA
}

// collapse will accept a list of expressions, each of them will be collapsed into a temporal variable
// and become an ADR node of this variable. If optLast is true, the last expression won't use one.
func (table *symTable) collapse(nodes []parser.Node, optLast bool) {
	var lastNode parser.Node
	var lastNodeIndex int

	for i, n := range nodes {
		if !n.Valid() {
			break
		}

		if n.Type() == parser.NODES {
			tmp := table.borrowAddress()
			table.codeSeg.WriteInst(typ.OpSet, tmp, table.compileNode(n))
			nodes[i] = parser.Addr(tmp)
			lastNode, lastNodeIndex = n, i
		}
	}

	if optLast && lastNode.Valid() {
		i := table.codeSeg.LastInst()
		// [set a b]
		table.codeSeg.TruncLast()
		table.freeAddr(i.A)
		nodes[lastNodeIndex] = parser.Addr(uint16(i.B))
	}
}
