package potatolang

import (
	"bytes"
	"math/rand"
	"time"

	"github.com/coyove/potatolang/parser"
)

var _nodeRegA = parser.Nod(regA)

func (table *symtable) compileChainOp(chain *parser.Node) (code packet, yx uint16) {
	buf := newpacket()
	doblock := chain.CplIndex(0).Sym() == "do"

	if doblock {
		table.addMaskedSymTable()
	}

	for _, a := range chain.Cpl() {
		code, yx = table.compileNode(a)
		buf.Write(code)
	}

	if doblock {
		table.removeMaskedSymTable()
	}

	return buf, yx
}

func (table *symtable) compileSetOp(atoms []*parser.Node) (code packet, yx uint16) {
	calcDest := func() uint16 {
		aDest := atoms[1].Value.(parser.Symbol)
		newYX := table.get(aDest)
		if atoms[0].Sym() == parser.AMove {
			if newYX == regNil {
				newYX = table.borrowAddress()
				table.put(aDest, newYX)
			}
		} else {
			newYX = table.borrowAddress()
			table.put(aDest, newYX)
		}
		return newYX
	}

	buf := newpacket()
	aSrc := atoms[2]
	switch aSrc.Type() {
	case parser.SYM:
		srcName := aSrc.Value.(parser.Symbol)
		valueIndex := table.get(srcName)
		addr := calcDest()
		buf.WriteOP(OpSet, addr, valueIndex)
		return buf, addr
	case parser.NUM, parser.STR:
		addr := calcDest()
		buf.WriteOP(OpSet, addr, table.loadK(aSrc.Value))
		return buf, addr
	}

	code, newYX := table.compileNode(aSrc)
	buf.Write(code)

	addr := calcDest()
	buf.WriteOP(OpSet, addr, newYX)
	buf.WritePos(atoms[0].Position)
	return buf, addr
}

func (table *symtable) compileHashArrayOp(atoms []*parser.Node) (code packet, yx uint16) {
	switch atoms[0].Value.(parser.Symbol) {
	case parser.AHash, parser.AArray:
		code = table.collapse(atoms[1].Cpl(), true)

		args := atoms[1].Cpl()
		for i := 0; i < len(args); i += 2 {
			if i+1 >= len(args) {
				table.writeOpcode(&code, OpPush, args[i], nil)
			} else {
				table.writeOpcode(&code, OpPush2, args[i], args[i+1])
			}
		}

		table.returnAddresses(args)
	case parser.AHashArray:
		hashPart := table.collapse(atoms[1].Cpl(), false)
		code.Write(hashPart)
		arrayPart := table.collapse(atoms[2].Cpl(), false)
		code.Write(arrayPart)

		arrayElements := atoms[2].Cpl()
		for i := 0; i < len(arrayElements); i += 2 {
			if i+1 >= len(arrayElements) {
				table.writeOpcode(&code, OpPush, arrayElements[i], nil)
			} else {
				table.writeOpcode(&code, OpPush2, arrayElements[i], arrayElements[i+1])
			}
		}
		code.WriteOP(OpMakeArray, 0, 0)

		hashElements := atoms[1].Cpl()
		for i := 0; i < len(hashElements); i += 2 {
			table.writeOpcode(&code, OpPush2, hashElements[i], hashElements[i+1])
		}
		code.WriteOP(OpMakeHash, 1, 0)

		table.returnAddresses(arrayElements)
		table.returnAddresses(hashElements)
		return code, regA
	}

	switch atoms[0].Value.(parser.Symbol) {
	case parser.AHash:
		code.WriteOP(OpMakeHash, 0, 0)
	case parser.AArray:
		code.WriteOP(OpMakeArray, 0, 0)
	}
	code.WritePos(atoms[0].Position)
	return code, regA
}

// writeOpcode3 accepts 3 arguments at most, 2 arguments will be encoded into opcode itself, the 3rd one will be in regA
func (table *symtable) writeOpcode3(bop _Opcode, atoms []*parser.Node) (buf packet, yx uint16) {
	// first atom: the op name, tail atoms: the args
	if len(atoms) > 4 {
		panic("shouldn't happen: too many arguments")
	}

	atoms = append([]*parser.Node{}, atoms...) // duplicate

	if bop == OpStore {
		buf = table.collapse(atoms[1:], true)

		// (atoms    1      2    3 )
		// (store subject value key) subject => opa, key => $a, value => opb

		for i := 1; i <= 2; i++ { // subject and value shouldn't use regA
			if atoms[i].Type() == parser.ADR && atoms[i].Value.(uint16) == regA {
				addr := table.borrowAddress()
				table.writeOpcode(&buf, OpSet, parser.Nod(addr), _nodeRegA)
				atoms[i] = parser.Nod(addr)
			}
		}

		// We would love to see key using regA, in this case writeOpcode will just omit it
		table.writeOpcode(&buf, OpSet, _nodeRegA, atoms[3])
		table.writeOpcode(&buf, OpStore, atoms[1], atoms[2])
		table.returnAddresses(atoms[1:])
		return buf, regA
	}

	buf = table.collapse(atoms[1:], true)

	switch bop {
	case OpGetB, OpPatchVararg:
		buf.WriteOP(bop, 0, 0)
	case OpAddressOf:
		yx := table.get(atoms[1].Sym())
		if yx == regNil {
			// For addressof op, atoms[1] is always an identifier, if it is not defined, define it
			yx = table.borrowAddress()
			buf.WriteOP(OpSet, yx, regNil)
			table.put(atoms[1].Sym(), yx)
		}
		buf.WriteOP(bop, yx, 0)
	case OpNot, OpRet, OpYield, OpLen, OpSetB:
		// unary op
		table.writeOpcode(&buf, bop, atoms[1], nil)
	default:
		// binary op
		table.writeOpcode(&buf, bop, atoms[1], atoms[2])
		table.returnAddresses(atoms[1:])
	}

	return buf, regA
}

func (table *symtable) compileFlatOp(atoms []*parser.Node) (code packet, yx uint16) {
	head := atoms[0].Value.(parser.Symbol)
	switch head {
	case "nop":
		return newpacket(), regA
	}
	op, ok := flatOpMapping[head]
	if !ok {
		panicf("compileFlatOp: shouldn't happen: invalid op: %#v", atoms[0])
	}
	code, yx = table.writeOpcode3(op, atoms)
	for _, a := range atoms {
		if a.Position.Source != "" {
			code.WritePos(a.Position)
			break
		}
	}
	return
}

// [and a b] => $a = a if not a then return else $a = b end
// [or a b]  => $a = a if a then do nothing else $a = b end
func (table *symtable) compileAndOrOp(atoms []*parser.Node) (code packet, yx uint16) {
	bop := OpIfNot
	if atoms[0].Value.(parser.Symbol) == parser.AOr {
		bop = OpIf
	}

	buf := newpacket()

	table.writeOpcode(&buf, OpSet, _nodeRegA, atoms[1])
	buf.WriteOP(bop, regA, 0)
	c2 := buf.Len()

	table.writeOpcode(&buf, OpSet, _nodeRegA, atoms[2])

	_, yx, _ = op(buf.data[c2-1])
	buf.data[c2-1] = makejmpop(bop, yx, buf.Len()-c2)
	buf.WritePos(atoms[0].Position)
	return buf, regA
}

// [if condition [true-chain ...] [false-chain ...]]
func (table *symtable) compileIfOp(atoms []*parser.Node) (code packet, yx uint16) {
	condition := atoms[1]
	trueBranch, falseBranch := atoms[2], atoms[3]
	buf := newpacket()

	code, yx = table.compileNode(condition)
	buf.Write(code)
	condyx := yx

	table.addMaskedSymTable()
	trueCode, _ := table.compileNode(trueBranch)
	falseCode, _ := table.compileNode(falseBranch)
	table.removeMaskedSymTable()
	if len(falseCode.data) > 0 {
		buf.WriteJmpOP(OpIfNot, condyx, len(trueCode.data)+1)
		buf.WritePos(atoms[0].Position)
		buf.Write(trueCode)
		buf.WriteJmpOP(OpJmp, 0, len(falseCode.data))
		buf.Write(falseCode)
	} else {
		buf.WriteJmpOP(OpIfNot, condyx, len(trueCode.data))
		buf.WritePos(atoms[0].Position)
		buf.Write(trueCode)
	}
	return buf, regA
}

// [call callee [args ...]]
func (table *symtable) compileCallOp(nodes []*parser.Node) (code packet, yx uint16) {
	tmp := append([]*parser.Node{nodes[1]}, nodes[2].Cpl()...)
	code = table.collapse(tmp, true)

	for i := 1; i < len(tmp); i += 2 {
		if i+1 >= len(tmp) {
			table.writeOpcode(&code, OpPush, tmp[i], nil)
		} else {
			table.writeOpcode(&code, OpPush2, tmp[i], tmp[i+1])
		}
	}

	table.writeOpcode(&code, OpCall, tmp[0], nil)
	code.WritePos(nodes[0].Position)

	table.returnAddresses(tmp)
	return code, regA
}

// [lambda name? [namelist] [chain ...]]
func (table *symtable) compileLambdaOp(atoms []*parser.Node) (code packet, yx uint16) {
	table.envescape = true
	vararg := false
	params := atoms[2]
	newtable := newsymtable()
	newtable.parent = table

	for i, p := range params.Cpl() {
		argname := p.Value.(parser.Symbol)
		if argname == "..." {
			if i != len(params.Cpl())-1 {
				panicf("%#v: vararg must be the last parameter", atoms[0])
			}
			atoms[3] = parser.Cpl(
				parser.ABegin,
				parser.Cpl(
					parser.ASet,
					parser.Symbol("arg"),
					parser.Cpl(parser.APatchVararg).SetPos0(atoms[0]),
				).SetPos0(atoms[0]),
				atoms[3],
			)
			vararg = true
			break
		}
		if _, ok := newtable.sym[argname]; ok {
			panicf("%#v: duplicated parameter: %s", atoms[0], argname)
		}
		newtable.put(argname, uint16(i))
	}

	ln := len(newtable.sym)
	if ln > 255 {
		panicf("%#v: too many parameters", atoms[0])
	}

	newtable.vp = uint16(ln)

	code, yx = newtable.compileNode(atoms[3])

	code.WriteOP(OpEOB, 0, 0)
	buf := newpacket()
	cls := Closure{}
	cls.NumParam = byte(ln)
	if newtable.y {
		cls.Set(ClsYieldable)
	}
	if !newtable.envescape {
		cls.Set(ClsNoEnvescape)
	}
	if vararg {
		cls.Set(ClsVararg)
	}

	// (ln: 8bit) + (cls.options: 8bit) + (len(consts): 10bit)
	opaopb := uint32(ln)<<18 | uint32(cls.options)<<10 | uint32(len(newtable.consts))
	buf.WriteOP(OpLambda, uint16(opaopb>>13), uint16(opaopb&0x1fff))
	buf.WriteConsts(newtable.consts)

	cls.Code = code.data
	src := cls.String() + "@" + code.source
	buf.WriteString(src)

	buf.Write32(uint32(len(code.pos)))
	buf.WriteRaw(u32FromBytes(code.pos))

	buf.Write32(uint32(len(code.data)))

	// Note buf.source will be set to Code.source in buf.Write
	// but buf.source shouldn't be changed, so set code.source to buf.source
	code.source = buf.source

	buf.Write(code)
	buf.WritePos(atoms[0].Position)
	return buf, regA
}

var staticWhileHack [8]uint32

func init() {
	rand.Seed(time.Now().Unix())
	for i := range staticWhileHack {
		staticWhileHack[i] = makeop(OpNOP, uint16(rand.Uint32()), uint16(rand.Uint32()))
	}
}

// [continue | break]
func (table *symtable) compileContinueBreakOp(atoms []*parser.Node) (code packet, yx uint16) {
	buf := newpacket()
	if table.inloop == 0 {
		panicf("%#v: invalid statement outside loop", atoms[0])
	}

	if atoms[0].Value.(parser.Symbol) == parser.AContinue {
		buf.WriteRaw(staticWhileHack[:4]) // write a 'continue' placeholder
	} else {
		buf.WriteRaw(staticWhileHack[4:]) // write a 'break' placeholder
	}
	return buf, regA
}

// [loop [chain ...]]
func (table *symtable) compileWhileOp(atoms []*parser.Node) (code packet, yx uint16) {
	buf := newpacket()

	table.inloop++
	table.addMaskedSymTable()
	code, yx = table.compileNode(atoms[1])
	table.removeMaskedSymTable()
	table.inloop--

	buf.Write(code)
	buf.WriteJmpOP(OpJmp, 0, -code.Len()-1)

	code = buf
	code2 := u32Bytes(code.data)

	// Search for special 'continue' placeholder and replace it with a OP_JMP to the
	// beginning of the Code
	continueflag := u32Bytes(staticWhileHack[:4])
	for i := 0; i < len(code2); {
		x := bytes.Index(code2[i:], continueflag)
		if x == -1 {
			break
		}
		idx := (i + x) / 4
		code.data[idx] = makejmpop(OpJmp, 0, -idx-1)
		i += x + 4
	}

	// Search for special 'break' placeholder and replace it with a OP_JMP to the
	// end of the Code
	breakflag := u32Bytes(staticWhileHack[4:])
	for i := 0; i < len(code2); {
		x := bytes.Index(code2[i:], breakflag)
		if x == -1 {
			break
		}
		idx := (i + x) / 4
		code.data[idx] = makejmpop(OpJmp, 0, len(code.data)-idx-1)
		i += x + 4
	}
	return buf, regA
}

// collapse will accept a list of nodes, for every expression inside,
// it will be collapsed into a temp variable and be replaced with a ADR node,
// For the last expression, it will be collapsed but not use a temp variable unless optLast == false
func (table *symtable) collapse(atoms []*parser.Node, optLast bool) (buf packet) {
	buf = newpacket()

	var lastCompound struct {
		n *parser.Node
		i int
	}

	for i, atom := range atoms {
		if atom == nil {
			break
		}

		var yx uint16
		var code packet

		if atom.Type() == parser.CPL {
			code, yx = table.compileNodeInto(atom, true, 0)

			atoms[i] = parser.Nod(yx)
			buf.Write(code)

			lastCompound.n = atom
			lastCompound.i = i
		}
	}

	if lastCompound.n != nil {
		if optLast {
			_, old, opb := op(buf.data[len(buf.data)-1])
			buf.TruncateLast(1)
			table.returnAddress(old)
			atoms[lastCompound.i].Value = opb
		}
	}

	return
}
