package nj

import (
	"math"

	"github.com/coyove/nj/bas"
	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/parser"
	"github.com/coyove/nj/typ"
)

var nodeRegA = parser.Address(typ.RegA)

// [prog expr1 expr2 ...]
func compileProgBlock(table *symTable, node *parser.Prog) uint16 {
	if node.DoBlock {
		table.addMaskedSymTable()
	}

	yx := typ.RegA
	for _, a := range node.Stats {
		if a == nil {
			continue
		}
		switch a.(type) {
		case parser.Address, parser.Primitive, *parser.Symbol:
			// e.g.: [prog "a string"] will be transformed into: [prog [set $a "a string"]]
			yx = table.compileNode(a)
			table.codeSeg.WriteInst(typ.OpSet, typ.RegA, yx)
		default:
			yx = table.compileNode(a)
		}
	}

	if node.DoBlock {
		table.removeMaskedSymTable()
	}
	return yx
}

// local a = b
func compileDeclare(table *symTable, node *parser.Declare) uint16 {
	dest := bas.Str(node.Name.Name)
	if bas.GetGlobalName(dest) > 0 || dest == staticTrue || dest == staticFalse || dest == staticThis || dest == staticSelf {
		table.panicnode(node.Name, "can't bound to a global static name")
	}

	destAddr := table.borrowAddress()
	defer table.put(dest, destAddr) // execute in defer in case of: a = 1 do local a = a end
	table.codeSeg.WriteInst(typ.OpSet, destAddr, table.compileNode(node.Value))
	table.codeSeg.WriteLineNum(node.Line)
	return destAddr
}

// a = b
func compileAssign(table *symTable, node *parser.Assign) uint16 {
	dest := bas.Str(node.Name.Name)
	if bas.GetGlobalName(dest) > 0 || dest == staticTrue || dest == staticFalse || dest == staticThis || dest == staticSelf {
		table.panicnode(node.Name, "can't assign to a global static name")
	}
	destAddr, declared := table.get(dest)
	if !declared {
		// a is not declared yet
		destAddr = table.borrowAddress()

		// Do not use t.put() because it may put the symbol into masked tables
		// e.g.: do a = 1 end
		table.sym.Set(dest, bas.Int64(int64(destAddr)))
	} else {
	}
	table.codeSeg.WriteInst(typ.OpSet, destAddr, table.compileNode(node.Value))
	table.codeSeg.WriteLineNum(node.Line)
	return destAddr
}

func compileUnary(table *symTable, node *parser.Unary) uint16 {
	nodes := table.collapse(true, node.A)
	table.writeInst1(node.Op, nodes[0])
	table.freeAddr(nodes)
	table.codeSeg.WriteLineNum(node.Line)
	return typ.RegA
}

func compileBinary(table *symTable, node *parser.Binary) uint16 {
	nodes := table.collapse(true, node.A, node.B)
	table.writeInst2(node.Op, nodes[0], nodes[1])
	table.freeAddr(nodes)
	table.codeSeg.WriteLineNum(node.Line)
	return typ.RegA
}

func compileTenary(table *symTable, node *parser.Tenary) uint16 {
	nodes := table.collapse(true, node.A, node.B, node.C)
	table.writeInst3(node.Op, nodes[0], nodes[1], nodes[2])
	table.freeAddr(nodes)
	table.codeSeg.WriteLineNum(node.Line)
	return typ.RegA
}

func compileBitwise(table *symTable, node *parser.Bitwise) uint16 {
	nodes := table.collapse(true, node.A, node.B)
	// Lookup table:     And,  c, d, e, f, g, h, i, j, k, Lsh,  n, Or,p, q, Rsh,  t, Ursh, w, Xor,  z
	lookup := [26]uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 1, 0, 0, 4, 0, 0, 5, 0, 0, 2, 0, 0}
	table.writeInst3(typ.OpBitOp, nodes[0], nodes[1], parser.Address(lookup[node.Op[0]-'a']))
	table.freeAddr(nodes[1:])
	table.codeSeg.WriteLineNum(node.Line)
	return typ.RegA
}

// [and a b] => $a = a if not a then goto out else $a = b end ::out::
func compileAnd(table *symTable, node *parser.And) uint16 {
	table.writeInst2(typ.OpSet, nodeRegA, node.A)

	table.codeSeg.WriteJmpInst(typ.OpIfNot, 0)
	part1 := table.codeSeg.Len()

	table.writeInst2(typ.OpSet, nodeRegA, node.B)
	part2 := table.codeSeg.Len()

	table.codeSeg.Code[part1-1] = typ.JmpInst(typ.OpIfNot, part2-part1)
	return typ.RegA
}

// [or a b]  => $a = a if not a then $a = b end
func compileOr(table *symTable, node *parser.Or) uint16 {
	table.writeInst2(typ.OpSet, nodeRegA, node.A)

	table.codeSeg.WriteJmpInst(typ.OpIfNot, 1)
	table.codeSeg.WriteJmpInst(typ.OpJmp, 0)
	part1 := table.codeSeg.Len()

	table.writeInst2(typ.OpSet, nodeRegA, node.B)
	part2 := table.codeSeg.Len()

	table.codeSeg.Code[part1-1] = typ.JmpInst(typ.OpJmp, part2-part1)
	return typ.RegA
}

// [if condition [true-chain ...] [false-chain ...]]
func compileIf(table *symTable, node *parser.If) uint16 {
	condyx := table.compileNode(node.Cond)

	table.addMaskedSymTable()

	if condyx != typ.RegA {
		table.codeSeg.WriteInst(typ.OpSet, typ.RegA, condyx)
	}

	table.codeSeg.WriteJmpInst(typ.OpIfNot, 0)
	init := table.codeSeg.Len()

	table.compileNode(node.True)
	part1 := table.codeSeg.Len()

	table.codeSeg.WriteJmpInst(typ.OpJmp, 0)

	if node.False != nil {
		table.compileNode(node.False)
		part2 := table.codeSeg.Len()

		table.removeMaskedSymTable()

		table.codeSeg.Code[init-1] = typ.JmpInst(typ.OpIfNot, part1-init+1)
		table.codeSeg.Code[part1] = typ.JmpInst(typ.OpJmp, part2-part1-1)
	} else {
		table.removeMaskedSymTable()

		// The last inst is used to skip the false branch, since we don't have one, we don't need this jmp
		table.codeSeg.TruncLast()
		table.codeSeg.Code[init-1] = typ.JmpInst(typ.OpIfNot, part1-init)
	}
	return typ.RegA
}

// [object [k1, v1, k2, v2, ...]]
func compileObject(table *symTable, node parser.ExprAssignList) uint16 {
	tmp := table.collapse(true, node.ExpandAsExprList()...)
	for i := 0; i < len(tmp); i += 2 {
		table.writeInst1(typ.OpPush, tmp[i])
		table.writeInst1(typ.OpPush, tmp[i+1])
	}
	table.codeSeg.WriteInst(typ.OpCreateObject, 0, 0)
	return typ.RegA
}

// [array [a, b, c, ...]]
func compileArray(table *symTable, node parser.ExprList) uint16 {
	nodes := table.collapse(true, node...)
	for _, x := range nodes {
		table.writeInst1(typ.OpPush, x)
	}
	table.codeSeg.WriteInst(typ.OpCreateArray, 0, 0)
	return typ.RegA
}

// [call callee [args ...]]
func compileCall(table *symTable, node *parser.Call) uint16 {
	tmp := table.collapse(true, append(node.Args, node.Callee)...)
	callee := tmp[len(tmp)-1]
	args := tmp[:len(tmp)-1]
	for i := 0; i < len(args)-1; i++ {
		table.writeInst1(typ.OpPush, args[i])
	}

	if node.Vararg {
		table.writeInst1(typ.OpPushUnpack, args[len(args)-1])
		table.writeInst2(node.Op, callee, parser.Address(typ.RegPhantom))
	} else if len(args) == 0 {
		table.writeInst2(node.Op, callee, parser.Address(typ.RegPhantom))
	} else {
		table.writeInst2(node.Op, callee, args[len(args)-1])
	}

	table.codeSeg.WriteLineNum(node.Line)
	table.freeAddr(tmp)
	return typ.RegA
}

func compileFunction(table *symTable, node *parser.Function) uint16 {
	newtable := newSymTable(table.options)
	newtable.name = table.name
	newtable.codeSeg.Pos.Name = table.name
	newtable.global = table.getGlobal()

	for i, p := range node.Args {
		name := p.(*parser.Symbol).Name
		if newtable.sym.Contains(bas.Str(name)) {
			table.panicnode(node, "duplicated parameter %q", name)
		}
		newtable.put(bas.Str(name), uint16(i))
	}

	if ln := newtable.sym.Len(); ln > 255 {
		table.panicnode(node, "too many parameters (%d > 255)", ln)
	}

	newtable.vp = uint16(newtable.sym.Len())
	newtable.compileNode(node.Body)
	newtable.patchGoto()

	if a := newtable.sym.Get(staticSelf); a != bas.Nil {
		newtable.codeSeg.Code = append([]typ.Inst{
			{Opcode: typ.OpFunction, A: typ.RegA},
			{Opcode: typ.OpSet, A: uint16(a.Int64()), B: typ.RegA},
		}, newtable.codeSeg.Code...)
		newtable.codeSeg.Pos.Offset += 2
	}

	if a := newtable.sym.Get(staticThis); a != bas.Nil {
		newtable.codeSeg.Code = append([]typ.Inst{
			{Opcode: typ.OpSet, A: uint16(a.Int64()), B: typ.RegA},
		}, newtable.codeSeg.Code...)
		newtable.codeSeg.Pos.Offset += 1
	}

	code := newtable.codeSeg
	code.WriteInst(typ.OpRet, typ.RegGlobalFlag, 0) // return nil

	localDeclare := table.borrowAddress()
	table.put(bas.Str(node.Name), localDeclare)

	var captureList []string
	if table.global != nil {
		captureList = table.symbolsToDebugLocals()
	}

	obj := bas.NewBareFunc(
		node.Name,
		node.Vararg,
		byte(len(node.Args)),
		newtable.vp,
		newtable.symbolsToDebugLocals(),
		captureList,
		code,
	)

	fm := &table.getGlobal().funcsMap
	fidx := fm.Get(bas.Str(node.Name))
	table.getGlobal().constMap.Set(obj.ToValue(), fidx)

	table.codeSeg.WriteInst3(typ.OpFunction, uint16(fidx.Int()),
		uint16(internal.IfInt(table.global == nil, 0, 1)),
		typ.RegA,
	)
	table.codeSeg.WriteInst(typ.OpSet, localDeclare, typ.RegA)
	table.codeSeg.WriteLineNum(node.Line)
	return typ.RegA
}

// [break|continue]
func compileBreakContinue(table *symTable, node *parser.BreakContinue) uint16 {
	if len(table.forLoops) == 0 {
		table.panicnode(node, "outside loop")
	}
	bl := table.forLoops[len(table.forLoops)-1]
	if !node.Break {
		table.compileNode(bl.continueNode)
		table.codeSeg.WriteJmpInst(typ.OpJmp, bl.continueGoto-len(table.codeSeg.Code)-1)
	} else {
		bl.breakContinuePos = append(bl.breakContinuePos, table.codeSeg.Len())
		table.codeSeg.WriteJmpInst(typ.OpJmp, 0)
	}
	return typ.RegA
}

// [loop body continue]
func compileLoop(table *symTable, node *parser.Loop) uint16 {
	init := table.codeSeg.Len()
	breaks := &breakLabel{
		continueNode: node.Continue,
		continueGoto: init,
	}

	table.forLoops = append(table.forLoops, breaks)
	table.addMaskedSymTable()
	table.compileNode(node.Body)
	table.removeMaskedSymTable()
	table.forLoops = table.forLoops[:len(table.forLoops)-1]

	table.codeSeg.WriteJmpInst(typ.OpJmp, -(table.codeSeg.Len()-init)-1)
	for _, idx := range breaks.breakContinuePos {
		table.codeSeg.Code[idx] = typ.JmpInst(typ.OpJmp, table.codeSeg.Len()-idx-1)
	}
	return typ.RegA
}

func compileGotoLabel(table *symTable, node *parser.GotoLabel) uint16 {
	if !node.Goto {
		if table.labelPos == nil {
			table.labelPos = map[string]int{}
		}
		table.labelPos[node.Label] = table.codeSeg.Len()
		return typ.RegA
	}

	if pos, ok := table.labelPos[node.Label]; ok {
		table.codeSeg.WriteJmpInst(typ.OpJmp, pos-(table.codeSeg.Len()+1))
	} else {
		table.codeSeg.WriteJmpInst(typ.OpJmp, 0)
		if table.forwardGoto == nil {
			table.forwardGoto = map[int]*parser.GotoLabel{}
		}
		table.forwardGoto[table.codeSeg.Len()-1] = node
	}
	return typ.RegA
}

func (table *symTable) patchGoto() {
	code := table.codeSeg.Code
	for ipos, node := range table.forwardGoto {
		pos, ok := table.labelPos[node.Label]
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
		if code[i].Opcode == typ.OpJmp && i > 0 && (code[i-1].Opcode == typ.OpInc || code[i-1].Opcode == typ.OpInc16) {
			// Inc-then-small-jump, see OpInc in eval.go
			if d := code[i].D() + 1; d >= math.MinInt16 && d <= math.MaxInt16 {
				code[i-1].C = uint16(int16(d))
			}
		}
	}
}

func compileRelease(table *symTable, node parser.Release) uint16 {
	for _, s := range node {
		s := bas.Str(s.Name)
		yx, _ := table.get(s)
		table.freeAddr(yx)
		t := table.sym
		if len(table.maskedSym) > 0 {
			t = table.maskedSym[len(table.maskedSym)-1]
		}
		if !t.Contains(s) {
			internal.ShouldNotHappen(node)
		}
		t.Delete(s)
	}
	return typ.RegA
}

// collapse will accept a list of expressions, each of them will be collapsed into a temporal variable
// and become an Address node. If optLast is true, the last expression will be directly using regA.
func (table *symTable) collapse(optLast bool, nodes ...parser.Node) []parser.Node {
	var lastNode parser.Node
	var lastNodeIndex int

	for i, n := range nodes {
		switch n.(type) {
		case parser.Address, parser.Primitive, *parser.Symbol:
			// No need to collapse
		case *parser.If:
			// 'if' is special because it can be used as an expresison, we can't optimize just one branch.
			// e.g.: if(cond, a[0], a[1])
			tmp := table.borrowAddress()
			res := compileIf(table, n.(*parser.If))
			table.codeSeg.Code = append(table.codeSeg.Code, typ.Inst{Opcode: typ.OpSet, A: tmp, B: res})
			nodes[i] = parser.Address(tmp)
			lastNode, lastNodeIndex = n, i
		default:
			tmp := table.borrowAddress()
			table.codeSeg.WriteInst(typ.OpSet, tmp, table.compileNode(n))
			nodes[i] = parser.Address(tmp)
			lastNode, lastNodeIndex = n, i
		}
	}

	if optLast && lastNode != nil {
		if i := table.codeSeg.LastInst(); i.Opcode == typ.OpSet && i.B == typ.RegA {
			// [set something $a]
			table.codeSeg.TruncLast()
			table.freeAddr(i.A)
			nodes[lastNodeIndex] = parser.Address(uint16(i.B))
		}
	}
	return nodes
}
