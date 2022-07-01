package nj

import (
	"math"
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

func (table *symTable) writeNodes3(bop byte, nodes []parser.Node) uint16 {
	// first atom: the splitInst Name, tail atoms: the args
	if len(nodes) > 4 {
		internal.ShouldNotHappen(nodes)
	}

	nodes = append([]parser.Node{}, nodes...) // duplicate
	table.collapse(nodes[1:], true)

	switch bop {
	case typ.OpStore, typ.OpSlice: // ternary
		table.writeInst3(bop, nodes[1], nodes[2], nodes[3])
	case typ.OpLoad: // special binary
		table.writeInst3(bop, nodes[1], nodes[2], nodeRegA)
	case typ.OpNot, typ.OpRet, typ.OpBitNot, typ.OpLen: // unary
		table.writeInst1(bop, nodes[1])
	default: // binary
		table.writeInst2(bop, nodes[1], nodes[2])
	}
	table.freeAddr(nodes[1:])
	return typ.RegA
}

func makeOPCompiler(op byte) func(table *symTable, nodes []parser.Node) uint16 {
	return func(table *symTable, nodes []parser.Node) uint16 {
		yx := table.writeNodes3(op, nodes)
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
	switch nodes[0].Value {
	case parser.STailCall.Value:
		op = typ.OpTailCall
	case parser.STryCall.Value:
		op = typ.OpTryCall
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
	newtable.codeSeg.Pos = &internal.VByte32{Name: table.name}
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
			table.panicnode(nodes[1], "duplicated parameter %q", n.Value.Str())
		}
		newtable.put(n.Value, uint16(i))
	}

	if ln := newtable.sym.Len(); ln > 255 {
		table.panicnode(nodes[1], "too many parameters (%d > 255)", ln)
	}

	newtable.vp = uint16(newtable.sym.Len())
	newtable.compileNode(nodes[3])
	newtable.patchGoto()

	if a := newtable.sym.Get(staticThis); a != bas.Nil {
		newtable.codeSeg.Code = append([]typ.Inst{
			{Opcode: typ.OpSet, A: uint16(a.Int64()), B: typ.RegA},
		}, newtable.codeSeg.Code...)
	}
	if a := newtable.sym.Get(staticSelf); a != bas.Nil {
		newtable.codeSeg.Code = append([]typ.Inst{
			{Opcode: typ.OpSelf},
			{Opcode: typ.OpSet, A: uint16(a.Int64()), B: typ.RegA},
		}, newtable.codeSeg.Code...)
	}

	code := newtable.codeSeg
	code.WriteInst(typ.OpRet, typ.RegGlobalFlag, 0) // return nil

	cls := &bas.Function{}
	cls.Variadic = varargIdx >= 0
	cls.NumParams = byte(len(params.Nodes()))
	cls.Name = nodes[1].Sym()
	cls.StackSize = newtable.vp
	cls.CodeSeg = code
	cls.Locals = newtable.symbolsToDebugLocals()
	cls.Method = strings.Contains(cls.Name, ".")

	obj := bas.NewObject(0)
	obj.SetPrototype(bas.Proto.Func)
	internal.SetObjFun(unsafe.Pointer(obj), unsafe.Pointer(cls))

	funcIdx := uint16(len(table.getGlobal().funcs))
	table.getGlobal().funcs = append(table.getGlobal().funcs, obj)
	table.codeSeg.WriteInst(typ.OpLoadFunc, funcIdx, 0)
	table.codeSeg.WriteLineNum(nodes[0].Line())
	return typ.RegA
}

// [break|continue]
func compileBreak(table *symTable, atoms []parser.Node) uint16 {
	if len(table.forLoops) == 0 {
		table.panicnode(atoms[0], "outside loop")
	}
	bl := table.forLoops[len(table.forLoops)-1]
	if atoms[0].Value == parser.SContinue.Value {
		table.compileNode(bl.continueNode)
		table.codeSeg.WriteJmpInst(typ.OpJmp, bl.continueGoto-len(table.codeSeg.Code)-1)
	} else {
		bl.breakContinuePos = append(bl.breakContinuePos, table.codeSeg.Len())
		table.codeSeg.WriteJmpInst(typ.OpJmp, 0)
	}
	return typ.RegA
}

// [loop body continue]
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
	for _, idx := range breaks.breakContinuePos {
		table.codeSeg.Code[idx] = typ.JmpInst(typ.OpJmp, table.codeSeg.Len()-idx-1)
	}
	return typ.RegA
}

func compileLabel(table *symTable, nodes []parser.Node) uint16 {
	if table.labelPos == nil {
		table.labelPos = map[string]int{}
	}
	table.labelPos[nodes[1].Value.Str()] = table.codeSeg.Len()
	return typ.RegA
}

func compileGoto(table *symTable, nodes []parser.Node) uint16 {
	label := nodes[1]
	if pos, ok := table.labelPos[label.Value.Str()]; ok {
		table.codeSeg.WriteJmpInst(typ.OpJmp, pos-(table.codeSeg.Len()+1))
	} else {
		table.codeSeg.WriteJmpInst(typ.OpJmp, 0)
		if table.forwardGoto == nil {
			table.forwardGoto = map[int]parser.Node{}
		}
		table.forwardGoto[table.codeSeg.Len()-1] = label
	}
	return typ.RegA
}

func (table *symTable) patchGoto() {
	code := table.codeSeg.Code
	for ipos, node := range table.forwardGoto {
		pos, ok := table.labelPos[node.Value.Str()]
		if !ok {
			table.panicnode(node, "label not found")
		}
		code[ipos] = typ.JmpInst(typ.OpJmp, pos-(ipos+1))
	}
	for i := range code {
		if code[i].Opcode == typ.OpJmp && code[i].D() != 0 {
			// Group continuous jumps into one single jump
			dest := int32(i) + code[i].D() + 1
			for int(dest) < len(code) {
				if c2 := code[dest]; c2.Opcode == typ.OpJmp && c2.D() != 0 {
					dest += c2.D() + 1
					continue
				}
				break
			}
			code[i] = code[i].SetD(dest - int32(i) - 1)
		}
		if code[i].Opcode == typ.OpJmp && i > 0 && code[i-1].Opcode == typ.OpInc {
			// Inc-then-small-jump, see OpInc in eval.go
			if d := code[i].D() + 1; d >= math.MinInt16 && d <= math.MaxInt16 {
				code[i-1].C = uint16(int16(d))
			}
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
		if i := table.codeSeg.LastInst(); i.Opcode == typ.OpSet && i.B == typ.RegA {
			// [set a $a]
			table.codeSeg.TruncLast()
			table.freeAddr(i.A)
			nodes[lastNodeIndex] = parser.Addr(uint16(i.B))
		}
	}
}
