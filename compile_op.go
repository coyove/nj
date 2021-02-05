package script

import (
	"github.com/coyove/script/parser"
)

var _nodeRegA = parser.NewAddress(regA)

func (table *symtable) compileChain(chain parser.Node) uint16 {
	doblock := chain.Nodes[0].SymbolValue() == (parser.ADoBlock)

	if doblock {
		table.addMaskedSymTable()
	}

	yx := regA
	for _, a := range chain.Nodes {
		yx = table.compileNode(a)
	}

	if doblock {
		table.removeMaskedSymTable()
	}

	return yx
}

func (table *symtable) compileSetMove(atoms []parser.Node) uint16 {
	aDest := atoms[1].SymbolValue()
	newYX := table.get(aDest)
	if atoms[0].SymbolValue() == (parser.AMove) {
		// a = b
		if newYX == table.loadK(nil) {
			t, h := table, 0
			if t.global != nil {
				t, h = t.global, 1
			}
			yx := t.borrowAddress()

			// Do not use t.put() because it may put the symbol into masked tables
			// e.g.: do a = 1 end
			t.sym[aDest] = &symbol{addr: yx}
			newYX = uint16(h)<<12 | yx
		}
	} else {
		// local a = b
		newYX = table.borrowAddress()
		defer table.put(aDest, newYX) // execute in defer in case of: a = 1 do local a = a end
	}

	fromYX := table.compileNode(atoms[2])
	table.code.writeInst(OpSet, newYX, fromYX)
	table.code.writePos(atoms[0].Pos())
	return newYX
}

func (table *symtable) compileReturn(atoms []parser.Node) uint16 {
	table.writeInst(OpRet, atoms[1], parser.Node{})
	return regA
}

// writeInst3 accepts 3 arguments at most, 2 arguments will be encoded into opCode itself, the 3rd one will be in regA
func (table *symtable) writeInst3(bop opCode, atoms []parser.Node) uint16 {
	// first atom: the splitInst Name, tail atoms: the args
	if len(atoms) > 4 {
		panic("DEBUG: too many arguments")
	}

	atoms = append([]parser.Node{}, atoms...) // duplicate

	if bop == OpStore {
		table.collapse(atoms[1:], true)

		// (atoms    1      2    3 )
		// (store subject value key) subject => opa, key => $a, value => opb

		for i := 1; i <= 2; i++ { // subject and value shouldn't use regA
			if atoms[i].Type == parser.Address && atoms[i].Addr == regA {
				n := parser.NewAddress(table.borrowAddress())
				table.writeInst(OpSet, n, _nodeRegA)
				atoms[i] = n
			}
		}

		// We would love to see 'key' using regA, in this case writeInst will just omit it
		table.writeInst(OpSet, _nodeRegA, atoms[3])
		table.writeInst(OpStore, atoms[1], atoms[2])
		table.returnAddresses(atoms[1:])
		return regA
	}

	if bop == OpSlice {
		table.collapse(atoms[1:], true)

		// (atoms    1      2    3 )
		// (slice subject start end) subject => opa, start => $a, end => opb

		for i := 1; i <= 3; i += 2 { // subject and end shouldn't use regA
			if atoms[i].Type == parser.Address && atoms[i].Addr == regA {
				n := parser.NewAddress(table.borrowAddress())
				table.writeInst(OpSet, n, _nodeRegA)
				atoms[i] = n
			}
		}

		// We would love to see 'start' using regA, in this case writeInst will just omit it
		table.writeInst(OpSet, _nodeRegA, atoms[2])
		table.writeInst(OpSlice, atoms[1], atoms[3])
		table.returnAddresses(atoms[1:])
		return regA
	}

	table.collapse(atoms[1:], true)

	switch bop {
	case OpNot, OpRet, OpLen, OpGLoad:
		// unary splitInst
		table.writeInst(bop, atoms[1], parser.Node{})
	default:
		// binary splitInst
		table.writeInst(bop, atoms[1], atoms[2])
		table.returnAddresses(atoms[1:])
	}

	return regA
}

func (table *symtable) compileFlat(atoms []parser.Node) uint16 {
	head := atoms[0].SymbolValue()
	op, ok := flatOpMapping[head]
	if !ok {
		panicf("DEBUG compileFlat invalid symbol: %v", atoms[0])
	}
	yx := table.writeInst3(op, atoms)
	if p := atoms[0].Pos(); p.Source != "" {
		table.code.writePos(p)
	}
	return yx
}

// [and a b] => $a = a if not a then return else $a = b end
// [or a b]  => $a = a if a then do nothing else $a = b end
func (table *symtable) compileAndOr(atoms []parser.Node) uint16 {
	bop := OpIfNot
	if atoms[0].SymbolValue() == (parser.AOr) {
		bop = OpIf
	}

	table.writeInst(OpSet, _nodeRegA, atoms[1])
	table.code.writeInst(bop, regA, 0)
	part1 := table.code.Len()

	table.writeInst(OpSet, _nodeRegA, atoms[2])
	part2 := table.code.Len()

	table.code.Code[part1-1] = jmpInst(bop, part2-part1)
	table.code.writePos(atoms[0].Pos())
	return regA
}

// [if condition [true-chain ...] [false-chain ...]]
func (table *symtable) compileIf(atoms []parser.Node) uint16 {
	condition := atoms[1]
	trueBranch, falseBranch := atoms[2], atoms[3]

	condyx := table.compileNode(condition)

	table.addMaskedSymTable()

	if condyx != regA {
		table.code.writeInst(OpSet, regA, condyx)
	}

	table.code.writeJmpInst(OpIfNot, 0)
	table.code.writePos(atoms[0].Pos())
	init := table.code.Len()

	table.compileNode(trueBranch)
	part1 := table.code.Len()

	table.code.writeJmpInst(OpJmp, 0)

	table.compileNode(falseBranch)
	part2 := table.code.Len()

	table.removeMaskedSymTable()

	if len(falseBranch.Nodes) > 0 {
		table.code.Code[init-1] = jmpInst(OpIfNot, part1-init+1)
		table.code.Code[part1] = jmpInst(OpJmp, part2-part1-1)
	} else {
		// The last inst is used to skip the false branch, since we don't have one, we don't need this jmp
		table.code.truncateLast()
		table.code.Code[init-1] = jmpInst(OpIfNot, part1-init)
	}
	return regA
}

func (table *symtable) compileList(nodes []parser.Node) uint16 {
	// [list [a, b, c, ...]]
	table.collapse(nodes[1].Nodes, true)
	for _, x := range nodes[1].Nodes {
		table.writeInst(OpPush, x, parser.Node{})
	}
	table.code.writeInst(OpList, 0, 0)
	return regA
}

// [call callee [args ...]]
func (table *symtable) compileCall(nodes []parser.Node) uint16 {
	tmp := append([]parser.Node{nodes[1]}, nodes[2].Nodes...)
	table.collapse(tmp, true)

	for i := 1; i < len(tmp); i++ {
		table.writeInst(OpPush, tmp[i], parser.NewAddress(uint16(i-1)))
	}

	switch nodes[0].SymbolValue() {
	case parser.ACallMap:
		table.writeInst(OpCallMap, tmp[0], parser.Node{})
	case parser.ACall:
		table.writeInst(OpCall, tmp[0], parser.NewAddress(0))
	case parser.ATailCall:
		table.writeInst(OpCall, tmp[0], parser.NewAddress(1))
	}

	if _, opa, _ := splitInst(table.code.Code[table.code.Len()-1]); opa == table.loadK(nil) {
		panicf("function %s not found", nodes[1])
	}

	table.code.writePos(nodes[0].Pos())
	table.returnAddresses(tmp)
	return regA
}

// [function Name [paramlist] [chain ...] docstring]
func (table *symtable) compileFunction(atoms []parser.Node) uint16 {
	if table.options.DefineFuncFilter != nil && !table.options.DefineFuncFilter(atoms[1].SymbolValue()) {
		panicf("%v: function definition is not allowed by options", atoms[1])
	}

	params := atoms[2]
	newtable := newsymtable(table.options)
	paramsString := []string{}

	if table.global == nil {
		newtable.global = table
	} else {
		newtable.global = table.global
	}

	for i, p := range params.Nodes {
		n := p.SymbolValue()
		if _, ok := newtable.sym[n]; ok {
			panicf("%v: duplicated parameter: %q", atoms[1], n)
		}
		newtable.put(n, uint16(i))
		paramsString = append(paramsString, n)
	}

	ln := len(newtable.sym)
	if ln > 255 {
		panicf("%v: too many parameters, 255 at most", atoms[1])
	}

	newtable.vp = uint16(len(newtable.sym))
	newtable.compileNode(atoms[3])
	newtable.patchGoto()

	code := newtable.code
	code.writeInst(OpRet, table.loadK(nil), 0)

	cls := &Func{}
	cls.Name = atoms[1].SymbolValue()
	cls.DocString = atoms[4].StringValue()
	cls.NumParams = byte(ln)
	cls.StackSize = newtable.vp
	cls.Code = code
	cls.Params = paramsString
	cls.Locals = newtable.symbolsToDebugLocals()

	table.funcs = append(table.funcs, cls)
	table.code.writeInst(OpLoadFunc, uint16(len(table.funcs))-1, 0)
	table.code.writePos(atoms[0].Pos())
	return regA
}

// [break]
func (table *symtable) compileBreak(atoms []parser.Node) uint16 {
	if len(table.forLoops) == 0 {
		panicf("break outside loop")
	}
	table.forLoops[len(table.forLoops)-1].labelPos = append(table.forLoops[len(table.forLoops)-1].labelPos, table.code.Len())
	table.code.writeJmpInst(OpJmp, 0)
	return regA
}

// [loop [chain ...]]
func (table *symtable) compileWhile(atoms []parser.Node) uint16 {
	init := table.code.Len()
	breaks := &breaklabel{}

	table.forLoops = append(table.forLoops, breaks)
	table.addMaskedSymTable()
	table.compileNode(atoms[1])
	table.removeMaskedSymTable()
	table.forLoops = table.forLoops[:len(table.forLoops)-1]

	table.code.writeJmpInst(OpJmp, -(table.code.Len()-init)-1)
	for _, idx := range breaks.labelPos {
		table.code.Code[idx] = jmpInst(OpJmp, table.code.Len()-idx-1)
	}
	return regA
}

func (table *symtable) compileGoto(atoms []parser.Node) uint16 {
	label := atoms[1].SymbolValue()
	if atoms[0].SymbolValue() == parser.ALabel { // :: label ::
		table.labelPos[label] = table.code.Len()
	} else { // goto label
		if pos, ok := table.labelPos[label]; ok {
			table.code.writeJmpInst(OpJmp, pos-(table.code.Len()+1))
		} else {
			table.code.writeJmpInst(OpJmp, 0)
			table.forwardGoto[table.code.Len()-1] = label
		}
	}
	return regA
}

func (table *symtable) patchGoto() {
	code := table.code.Code
	for i, l := range table.forwardGoto {
		pos, ok := table.labelPos[l]
		if !ok {
			panicf("label %q not found", l)
		}
		code[i] = jmpInst(OpJmp, pos-(i+1))
	}
}

func (table *symtable) compileRetAddr(atoms []parser.Node) uint16 {
	for i := 1; i < len(atoms); i++ {
		s := atoms[i].SymbolValue()
		yx := table.get(s)
		table.returnAddress(yx)
		if len(table.maskedSym) > 0 {
			delete(table.maskedSym[len(table.maskedSym)-1], s)
		} else {
			delete(table.sym, s)
		}
	}
	return regA
}

// collapse will accept a list of nodes, for every expression inside,
// it will be collapsed into a temp variable and be replaced with a ADR node,
// For the last expression, it will be collapsed but not use a temp variable unless optLast == false
func (table *symtable) collapse(atoms []parser.Node, optLast bool) {
	var lastCompound struct {
		n parser.Node
		i int
	}

	for i, atom := range atoms {
		if !atom.Valid() {
			break
		}

		if atom.Type == parser.Complex {
			yx := table.compileNodeInto(atom, true, 0)
			atoms[i] = parser.NewAddress(yx)

			lastCompound.n = atom
			lastCompound.i = i
		}
	}

	if lastCompound.n.Valid() {
		if optLast {
			_, old, opb := splitInst(table.code.Code[len(table.code.Code)-1])
			table.code.truncateLast()
			table.returnAddress(old)
			atoms[lastCompound.i] = parser.NewAddress(opb)
		}
	}

	return
}
