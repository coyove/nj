package script

import (
	"strings"

	"github.com/coyove/script/parser"
)

var _nodeRegA = parser.Node{Addr: regA, Type: parser.Address}

func (table *symtable) compileChainOp(chain parser.Node) uint16 {
	doblock := chain.Nodes[0].SymbolValue() == (parser.ADoBlock)

	if doblock {
		table.addMaskedSymTable()
	}

	for _, a := range chain.Nodes {
		table.compileNode(a)
	}

	if doblock {
		table.removeMaskedSymTable()
	}

	return regA
}

func (table *symtable) compileSetOp(atoms []parser.Node) uint16 {
	aDest := atoms[1].SymbolValue()
	newYX := table.get(aDest)
	if atoms[0].SymbolValue() == (parser.AMove) {
		// a = b
		if newYX == regNil {
			t, h := table, 0
			if t.global != nil {
				t, h = t.global, 1
			}
			yx := t.borrowAddress()

			// Do not use t.put() because it may put the symbol into masked tables
			// e.g.: do a = 1 end
			t.sym[aDest] = &symbol{addr: yx}
			newYX = uint16(h)<<11 | yx
		}
	} else {
		// local a = b
		newYX = table.borrowAddress()
		table.put(aDest, newYX)
	}

	fromYX := table.compileNode(atoms[2])
	table.code.writeOP(OpSet, newYX, fromYX)
	table.code.writePos(atoms[0].Pos())
	return newYX
}

func (table *symtable) compileRetOp(atoms []parser.Node) uint16 {
	op := OpRet
	if atoms[0].SymbolValue() == (parser.AYield) {
		op = OpYield
	}

	values := atoms[1].Nodes
	if len(values) == 0 { // return
		table.code.writeOP(op, regNil, 0)
		return regA
	}

	// return a1, ..., an
	table.collapse(values, true)

	for i := 1; i < len(values); i++ {
		if i == 1 {
			// First OpPushV will contain the total number of V in opb
			table.writeOpcode(OpPushV, values[i], parser.Node{Type: parser.Address, Addr: uint16(len(values) - 1)})
		} else {
			table.writeOpcode(OpPushV, values[i], parser.Node{})
		}
	}
	table.writeOpcode(op, values[0], parser.Node{})

	table.returnAddresses(values)
	return regA
}

// writeOpcode3 accepts 3 arguments at most, 2 arguments will be encoded into opCode itself, the 3rd one will be in regA
func (table *symtable) writeOpcode3(bop opCode, atoms []parser.Node) uint16 {
	// first atom: the op name, tail atoms: the args
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
				table.writeOpcode(OpSet, n, _nodeRegA)
				atoms[i] = n
			}
		}

		// We would love to see 'key' using regA, in this case writeOpcode will just omit it
		table.writeOpcode(OpSet, _nodeRegA, atoms[3])
		table.writeOpcode(OpStore, atoms[1], atoms[2])
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
				table.writeOpcode(OpSet, n, _nodeRegA)
				atoms[i] = n
			}
		}

		// We would love to see 'start' using regA, in this case writeOpcode will just omit it
		table.writeOpcode(OpSet, _nodeRegA, atoms[2])
		table.writeOpcode(OpSlice, atoms[1], atoms[3])
		table.returnAddresses(atoms[1:])
		return regA
	}

	table.collapse(atoms[1:], true)

	switch bop {
	case OpNot, OpRet, OpYield, OpLen:
		// unary op
		table.writeOpcode(bop, atoms[1], parser.Node{})
	default:
		// binary op
		table.writeOpcode(bop, atoms[1], atoms[2])
		table.returnAddresses(atoms[1:])
	}

	return regA
}

func (table *symtable) compileFlatOp(atoms []parser.Node) uint16 {
	head := atoms[0].SymbolValue()
	switch head {
	case parser.APopV:
		table.code.writeOP(OpPopV, 0, 0)
		return regA
	case parser.APopVClear:
		table.code.writeOP(OpPopVClear, 0, 0)
		return regA
	case parser.APopVAll:
		table.code.writeOP(OpPopVAll, 0, 0)
		return regA
	case parser.APopVAllA:
		table.code.writeOP(OpPopVAll, 1, 0)
		return regA
	}

	op, ok := flatOpMapping[head]
	if !ok {
		panicf("DEBUG compileFlatOp invalid op: %v", atoms[0])
	}
	yx := table.writeOpcode3(op, atoms)
	if p := atoms[0].Pos(); p.Source != "" {
		table.code.writePos(p)
	}
	return yx
}

// [and a b] => $a = a if not a then return else $a = b end
// [or a b]  => $a = a if a then do nothing else $a = b end
func (table *symtable) compileAndOrOp(atoms []parser.Node) uint16 {
	bop := OpIfNot
	if atoms[0].SymbolValue() == (parser.AOr) {
		bop = OpIf
	}

	table.writeOpcode(OpSet, _nodeRegA, atoms[1])
	table.code.writeOP(bop, regA, 0)
	part1 := table.code.Len()

	table.writeOpcode(OpSet, _nodeRegA, atoms[2])
	part2 := table.code.Len()

	table.code.Code[part1-1] = makejmpop(bop, regA, part2-part1)
	table.code.writePos(atoms[0].Pos())
	return regA
}

// [if condition [true-chain ...] [false-chain ...]]
func (table *symtable) compileIfOp(atoms []parser.Node) uint16 {
	condition := atoms[1]
	trueBranch, falseBranch := atoms[2], atoms[3]

	condyx := table.compileNode(condition)

	table.addMaskedSymTable()

	table.code.writeJmpOP(OpIfNot, condyx, 0)
	table.code.writePos(atoms[0].Pos())
	init := table.code.Len()

	table.compileNode(trueBranch)
	part1 := table.code.Len()

	table.code.writeJmpOP(OpJmp, 0, 0)

	table.compileNode(falseBranch)
	part2 := table.code.Len()

	table.removeMaskedSymTable()

	if len(falseBranch.Nodes) > 0 {
		table.code.Code[init-1] = makejmpop(OpIfNot, condyx, part1-init+1)
		table.code.Code[part1] = makejmpop(OpJmp, 0, part2-part1-1)
	} else {
		table.code.truncateLast() // the last op is used to skip the false branch, since we don't have one, we don't need this jmp
		table.code.Code[init-1] = makejmpop(OpIfNot, condyx, part1-init)
	}
	return regA
}

// [call callee [args ...]]
func (table *symtable) compileCallOp(nodes []parser.Node) uint16 {
	tmp := append([]parser.Node{nodes[1]}, nodes[2].Nodes...)
	table.collapse(tmp, true)

	for i := 1; i < len(tmp); i++ {
		table.writeOpcode(OpPush, tmp[i], parser.NewAddress(uint16(i-1)))
	}

	var opb uint16
	if nodes[0].SymbolValue() == (parser.ATailCall) {
		opb = 1
	}

	table.writeOpcode(OpCall, tmp[0], parser.NewAddress(opb))
	table.code.writePos(nodes[0].Pos())

	table.returnAddresses(tmp)
	return regA
}

// [lambda name? [namelist] [chain ...]]
func (table *symtable) compileLambdaOp(atoms []parser.Node) uint16 {
	vararg := false
	params := atoms[2]
	newtable := newsymtable()

	if table.global == nil {
		newtable.global = table
	} else {
		newtable.global = table.global
	}

	for i, p := range params.Nodes {
		n := p.SymbolValue()
		if p.IsSymbolDotDotDot() {
			if i != len(params.Nodes)-1 {
				panicf("%v: vararg must be the last parameter", atoms[0])
			}
			vararg = true
			if n != "..." {
				n = strings.TrimLeft(n, ".")
			}
		}
		if _, ok := newtable.sym[n]; ok {
			panicf("%v: duplicated parameter: %q", atoms[0], n)
		}
		newtable.put(n, uint16(i))
	}

	ln := len(newtable.sym)
	if vararg {
		ln--
	}
	if ln > 255 {
		panicf("%v: too many parameters", atoms[0])
	}

	newtable.vp = uint16(len(newtable.sym))
	newtable.compileNode(atoms[3])
	newtable.patchGoto()

	code := newtable.code
	code.writeOP(OpEOB, 0, 0)

	cls := &Func{}
	cls.name = atoms[1].SymbolValue()
	cls.numParams = byte(ln)
	cls.stackSize = newtable.vp
	cls.isVariadic = vararg
	cls.code = code
	cls.constTable = newtable.constsToValues()

	table.funcs = append(table.funcs, cls)
	table.code.writeOP(OpLoadFunc, uint16(len(table.funcs))-1, 0)
	table.code.writePos(atoms[0].Pos())
	return regA
}

// [break]
func (table *symtable) compileBreakOp(atoms []parser.Node) uint16 {
	if len(table.inloop) == 0 {
		panicf("break outside loop")
	}
	table.inloop[len(table.inloop)-1].labelPos = append(table.inloop[len(table.inloop)-1].labelPos, table.code.Len())
	table.code.writeJmpOP(OpJmp, 0, 0)
	return regA
}

// [loop [chain ...]]
func (table *symtable) compileWhileOp(atoms []parser.Node) uint16 {
	init := table.code.Len()
	breaks := &breaklabel{}

	table.inloop = append(table.inloop, breaks)
	table.addMaskedSymTable()
	table.compileNode(atoms[1])
	table.removeMaskedSymTable()
	table.inloop = table.inloop[:len(table.inloop)-1]

	table.code.writeJmpOP(OpJmp, 0, -(table.code.Len()-init)-1)
	for _, idx := range breaks.labelPos {
		table.code.Code[idx] = makejmpop(OpJmp, 0, table.code.Len()-idx-1)
	}
	return regA
}

func (table *symtable) compileGotoOp(atoms []parser.Node) uint16 {
	label := atoms[1].SymbolValue()
	if atoms[0].SymbolValue() == (parser.ALabel) { // :: label ::
		table.labelPos[label] = table.code.Len()
	} else { // goto label
		if pos, ok := table.labelPos[label]; ok {
			table.code.writeJmpOP(OpJmp, 0, pos-(table.code.Len()+1))
		} else {
			table.code.writeJmpOP(OpJmp, 0, 0)
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
		code[i] = makejmpop(OpJmp, 0, pos-(i+1))
	}
}

func (table *symtable) compileRetAddrOp(atoms []parser.Node) uint16 {
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
			_, old, opb := op(table.code.Code[len(table.code.Code)-1])
			table.code.truncateLast()
			table.returnAddress(old)
			atoms[lastCompound.i] = parser.NewAddress(opb)
		}
	}

	return
}
