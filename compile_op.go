package potatolang

import (
	"math"

	"github.com/coyove/potatolang/parser"
)

var _nodeRegA = parser.Node{regA}

func (table *symtable) compileChainOp(chain parser.Node) uint16 {
	doblock := chain.CplIndex(0).Sym().Equals(parser.ADoBlock)

	if doblock {
		table.addMaskedSymTable()
	}

	for _, a := range chain.Cpl() {
		table.compileNode(a)
	}

	if doblock {
		table.removeMaskedSymTable()
	}

	return regA
}

func (table *symtable) compileSetOp(atoms []parser.Node) uint16 {
	aDest := atoms[1].Value.(parser.Symbol)
	newYX := table.get(aDest)
	if atoms[0].Sym().Equals(parser.AMove) {
		// a = b
		if newYX == regNil {
			t, h := table, 0
			for t.parent != nil {
				t, h = t.parent, h+1
				if h > 6 {
					panicf("global variable: too deep")
				}
			}
			yx := t.borrowAddress()

			// Do not use t.put() because it may put the symbol into masked tables
			// e.g.: do a = 1 end
			t.sym[aDest.Text] = &symbol{usage: math.MaxInt32, addr: yx}
			newYX = uint16(h)<<10 | yx
		}
	} else {
		// local a = b
		newYX = table.borrowAddress()
		table.put(aDest.Text, newYX)
	}

	fromYX := table.compileNode(atoms[2])
	table.code.WriteOP(OpSet, newYX, fromYX)
	table.code.WritePos(atoms[0].Pos())
	return newYX
}

func (table *symtable) compileRetOp(atoms []parser.Node) uint16 {
	op := OpRet
	if atoms[0].Sym().Equals(parser.AYield) {
		table.y = true
		op = OpYield
	}

	values := atoms[1].Cpl()
	if len(values) == 0 { // return
		table.code.WriteOP(OpPopV, 0, 0) // clear env.V in case of side effects
		table.code.WriteOP(op, regNil, 0)
		return regA
	}

	// return a1, ..., an
	table.collapse(values, true)

	for i := 1; i < len(values); i++ {
		if i == 1 {
			// First OpPushV will contain the total number of V in opb
			table.writeOpcode(OpPushV, values[i], parser.Node{uint16(len(values) - 1)})
		} else {
			table.writeOpcode(OpPushV, values[i], parser.Node{})
		}
	}
	table.writeOpcode(op, values[0], parser.Node{})

	table.returnAddresses(values)
	return regA
}

func (table *symtable) compileHashArrayOp(atoms []parser.Node) uint16 {
	switch atoms[0].Value.(parser.Symbol).Text {
	case parser.AHash.Text, parser.AArray.Text:
		table.collapse(atoms[1].Cpl(), true)

		args := atoms[1].Cpl()
		for i := 0; i < len(args); i += 2 {
			if i+1 >= len(args) {
				table.writeOpcode(OpPush, args[i], parser.Node{})
			} else {
				table.writeOpcode(OpPush2, args[i], args[i+1])
			}
		}

		table.returnAddresses(args)
	case parser.AHashArray.Text:
		table.collapse(atoms[1].Cpl(), false)
		table.collapse(atoms[2].Cpl(), false)

		arrayElements := atoms[2].Cpl()
		for i := 0; i < len(arrayElements); i += 2 {
			if i+1 >= len(arrayElements) {
				table.writeOpcode(OpPush, arrayElements[i], parser.Node{})
			} else {
				table.writeOpcode(OpPush2, arrayElements[i], arrayElements[i+1])
			}
		}
		table.code.WriteOP(OpMakeTable, 2, 0)

		hashElements := atoms[1].Cpl()
		for i := 0; i < len(hashElements); i += 2 {
			table.writeOpcode(OpPush2, hashElements[i], hashElements[i+1])
		}
		table.code.WriteOP(OpMakeTable, 3, 0)
		table.code.WritePos(atoms[0].Pos())

		table.returnAddresses(arrayElements)
		table.returnAddresses(hashElements)
		return regA
	}

	switch atoms[0].Value.(parser.Symbol).Text {
	case parser.AHash.Text:
		table.code.WriteOP(OpMakeTable, 1, 0)
	case parser.AArray.Text:
		table.code.WriteOP(OpMakeTable, 2, 0)
	}
	table.code.WritePos(atoms[0].Pos())
	return regA
}

// writeOpcode3 accepts 3 arguments at most, 2 arguments will be encoded into opcode itself, the 3rd one will be in regA
func (table *symtable) writeOpcode3(bop _Opcode, atoms []parser.Node) uint16 {
	// first atom: the op name, tail atoms: the args
	if len(atoms) > 4 {
		panic("shouldn't happen: too many arguments")
	}

	atoms = append([]parser.Node{}, atoms...) // duplicate

	if bop == OpStore {
		table.collapse(atoms[1:], true)

		// (atoms    1      2    3 )
		// (store subject value key) subject => opa, key => $a, value => opb

		for i := 1; i <= 2; i++ { // subject and value shouldn't use regA
			if atoms[i].Type() == parser.ADR && atoms[i].Value.(uint16) == regA {
				addr := table.borrowAddress()
				table.writeOpcode(OpSet, parser.Node{addr}, _nodeRegA)
				atoms[i] = parser.Node{addr}
			}
		}

		// We would love to see key using regA, in this case writeOpcode will just omit it
		table.writeOpcode(OpSet, _nodeRegA, atoms[3])
		table.writeOpcode(OpStore, atoms[1], atoms[2])
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
	head := atoms[0].Value.(parser.Symbol)
	switch head.Text {
	case parser.APopVAll.Text:
		table.code.WriteOP(OpPopV, 2, 0)
		return regA
	case parser.APopV.Text:
		table.code.WriteOP(OpPopV, 1, 0)
		return regA
	case parser.APopVEnd.Text:
		table.code.WriteOP(OpPopV, 0, 0)
		return regA
	}

	op, ok := flatOpMapping[head.Text]
	if !ok {
		panicf("compileFlatOp: shouldn't happen: invalid op: %#v", atoms[0])
	}
	yx := table.writeOpcode3(op, atoms)
	if p := atoms[0].Pos(); p.Source != "" {
		table.code.WritePos(p)
	}
	return yx
}

// [and a b] => $a = a if not a then return else $a = b end
// [or a b]  => $a = a if a then do nothing else $a = b end
func (table *symtable) compileAndOrOp(atoms []parser.Node) uint16 {
	bop := OpIfNot
	if atoms[0].Value.(parser.Symbol).Equals(parser.AOr) {
		bop = OpIf
	}

	table.writeOpcode(OpSet, _nodeRegA, atoms[1])
	table.code.WriteOP(bop, regA, 0)
	part1 := table.code.Len()

	table.writeOpcode(OpSet, _nodeRegA, atoms[2])
	part2 := table.code.Len()

	table.code.data[part1-1] = makejmpop(bop, regA, part2-part1)
	table.code.WritePos(atoms[0].Pos())
	return regA
}

// [if condition [true-chain ...] [false-chain ...]]
func (table *symtable) compileIfOp(atoms []parser.Node) uint16 {
	condition := atoms[1]
	trueBranch, falseBranch := atoms[2], atoms[3]

	condyx := table.compileNode(condition)

	table.addMaskedSymTable()

	table.code.WriteJmpOP(OpIfNot, condyx, 0)
	table.code.WritePos(atoms[0].Pos())
	init := table.code.Len()

	table.compileNode(trueBranch)
	part1 := table.code.Len()

	table.code.WriteJmpOP(OpJmp, 0, 0)

	table.compileNode(falseBranch)
	part2 := table.code.Len()

	table.removeMaskedSymTable()

	if len(falseBranch.Cpl()) > 0 {
		table.code.data[init-1] = makejmpop(OpIfNot, condyx, part1-init+1)
		table.code.data[part1] = makejmpop(OpJmp, 0, part2-part1-1)
	} else {
		table.code.TruncateLast(1) // the last op is used to skip the false branch, since we don't have one, we don't need this jmp
		table.code.data[init-1] = makejmpop(OpIfNot, condyx, part1-init)
	}
	return regA
}

// [call callee [args ...]]
func (table *symtable) compileCallOp(nodes []parser.Node) uint16 {
	tmp := append([]parser.Node{nodes[1]}, nodes[2].Cpl()...)
	table.collapse(tmp, true)

	for i := 1; i < len(tmp); i += 2 {
		if i+1 >= len(tmp) {
			table.writeOpcode(OpPush, tmp[i], parser.Node{})
		} else {
			table.writeOpcode(OpPush2, tmp[i], tmp[i+1])
		}
	}

	var opb uint16
	if nodes[0].Sym().Equals(parser.ATailCall) {
		opb = 1
	}

	table.writeOpcode(OpCall, tmp[0], parser.Node{opb})
	table.code.WritePos(nodes[0].Pos())

	table.returnAddresses(tmp)
	return regA
}

// [lambda name? [namelist] [chain ...]]
func (table *symtable) compileLambdaOp(atoms []parser.Node) uint16 {
	table.envescape = true
	vararg := false
	params := atoms[2]
	newtable := newsymtable()
	newtable.parent = table

	for i, p := range params.Cpl() {
		argname := p.Value.(parser.Symbol)
		if argname.Text == "..." {
			if i != len(params.Cpl())-1 {
				panicf("%#v: vararg must be the last parameter", atoms[0])
			}
			vararg = true
		}
		if _, ok := newtable.sym[argname.Text]; ok {
			panicf("%#v: duplicated parameter: %s", atoms[0], argname)
		}
		newtable.put(argname.Text, uint16(i))
	}

	ln := len(newtable.sym)
	if vararg {
		ln--
	}
	if ln > 255 {
		panicf("%#v: too many parameters", atoms[0])
	}

	newtable.vp = uint16(len(newtable.sym))

	newtable.compileNode(atoms[3])
	newtable.patchGoto()

	code := newtable.code
	code.source = atoms[0].Pos().Source
	code.WriteOP(OpEOB, 0, 0)

	cls := Closure{}
	cls.NumParam = byte(ln)
	cls.setOpt(newtable.y, ClsYieldable)
	cls.setOpt(!newtable.envescape, ClsNoEnvescape)
	cls.setOpt(vararg, ClsVararg)

	// (ln: 8bit) + (cls.options: 8bit) + (len(consts): 10bit)
	opaopb := uint32(ln)<<18 | uint32(cls.options)<<10 | uint32(len(newtable.consts))
	table.code.WriteOP(OpLambda, uint16(opaopb>>13), uint16(opaopb&0x1fff))
	table.code.WriteConsts(newtable.consts)

	cls.Code = code.data
	src := cls.String() + " " + code.source
	table.code.WriteString(src)

	table.code.Write32(uint32(len(code.pos)))
	table.code.WriteRaw(u32FromBytes(code.pos))

	table.code.Write32(uint32(len(code.data)))

	table.code.Write(code)
	table.code.WritePos(atoms[0].Pos())
	return regA
}

// [continue | break]
func (table *symtable) compileContinueBreakOp(atoms []parser.Node) uint16 {
	if len(table.inloop) == 0 {
		panicf("break outside loop")
	}
	table.inloop[len(table.inloop)-1].labelPos = append(table.inloop[len(table.inloop)-1].labelPos, table.code.Len())
	table.code.WriteJmpOP(OpJmp, 0, 0)
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

	table.code.WriteJmpOP(OpJmp, 0, -(table.code.Len()-init)-1)
	for _, idx := range breaks.labelPos {
		table.code.data[idx] = makejmpop(OpJmp, 0, table.code.Len()-idx-1)
	}
	return regA
}

func (table *symtable) compileGotoOp(atoms []parser.Node) uint16 {
	label := atoms[1].Sym()
	if atoms[0].Sym().Equals(parser.ALabel) { // :: label ::
		table.labelPos[label.Text] = table.code.Len()
	} else { // goto label
		if pos, ok := table.labelPos[label.Text]; ok {
			table.code.WriteJmpOP(OpJmp, 0, pos-(table.code.Len()+1))
		} else {
			table.code.WriteJmpOP(OpJmp, 0, 0)
			table.forwardGoto[table.code.Len()-1] = label.Text
		}
	}
	return regA
}

func (table *symtable) patchGoto() {
	code := table.code.data
	for i, l := range table.forwardGoto {
		pos, ok := table.labelPos[l]
		if !ok {
			panicf("label %s not found", l)
		}
		code[i] = makejmpop(OpJmp, 0, pos-(i+1))
	}
}

func (table *symtable) compileRetAddrOp(atoms []parser.Node) uint16 {
	for i := 1; i < len(atoms); i++ {
		s := atoms[i].Sym()
		yx := table.get(s)
		table.returnAddress(yx)
		if len(table.maskedSym) > 0 {
			delete(table.maskedSym[len(table.maskedSym)-1], s.Text)
		} else {
			delete(table.sym, s.Text)
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

		if atom.Type() == parser.CPL {
			yx := table.compileNodeInto(atom, true, 0)
			atoms[i] = parser.Node{yx}

			lastCompound.n = atom
			lastCompound.i = i
		}
	}

	if lastCompound.n.Valid() {
		if optLast {
			_, old, opb := op(table.code.data[len(table.code.data)-1])
			table.code.TruncateLast(1)
			table.returnAddress(old)
			atoms[lastCompound.i].Value = opb
		}
	}

	return
}
