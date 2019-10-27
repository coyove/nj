package potatolang

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/coyove/potatolang/parser"
)

const (
	errUndeclaredVariable = " %+v: undeclared variable"
)

var _nop = makeop(OpNOP, 0, 0)
var _nodeRegA = parser.NewNode(parser.Naddr).SetValue(regA)

func (table *symtable) compileSetOp(atoms []*parser.Node) (code packet, yx uint16, err error) {
	aDest := atoms[1].Value.(string)
	newYX, ok := table.get(aDest)
	if atoms[0].S() == "move" {
		if !ok || (strings.HasPrefix(aDest, "$") && (newYX>>10 > 0)) {
			// variable names start with "$" will be a local variable
			newYX = table.borrowAddress()
			table.put(aDest, newYX)
		}
	} else {
		newYX = table.borrowAddress()
		table.put(aDest, newYX)
	}

	buf := newpacket()
	aSrc := atoms[2]
	switch aSrc.Type {
	case parser.Natom:
		valueIndex, ok := table.get(aSrc.Value.(string))
		if !ok {
			err = fmt.Errorf(errUndeclaredVariable, aSrc)
			return
		}
		buf.WriteOP(OpSet, newYX, valueIndex)
		return buf, newYX, nil
	case parser.Nnumber, parser.Nstring:
		buf.WriteOP(OpSet, newYX, table.loadK(&buf, aSrc.Value))
		return buf, newYX, nil
	}

	code, _, err = table.compileCompoundInto(aSrc, false, newYX)
	if err != nil {
		return
	}
	buf.Write(code)
	buf.WritePos(atoms[0].Meta)
	return buf, newYX, nil
}

func (table *symtable) compileRetOp(atoms []*parser.Node) (code packet, yx uint16, err error) {
	if atoms[0].S() == "yield" {
		table.y = true
		return table.writeOpcode3(OpYield, atoms)
	}
	return table.writeOpcode3(OpRet, atoms)
}

func (table *symtable) compileMapArrayOp(atoms []*parser.Node) (code packet, yx uint16, err error) {
	var buf packet
	if buf, err = table.decompound(atoms[1].C()); err != nil {
		return
	}
	for _, atom := range atoms[1].C() {
		if err = table.writeOpcode(&buf, OpPush, atom, nil); err != nil {
			return
		}
	}
	if atoms[0].Value.(string) == "map" {
		buf.WriteOP(OpMakeMap, 0, 0)
	} else {
		buf.WriteOP(OpMakeMap, 1, 0)
	}
	buf.WritePos(atoms[0].Meta)
	return buf, regA, nil
}

// writeOpcode3 accepts 3 arguments at most, 2 arguments will be encoded into opcode itself, the 3rd one will be in regA
func (table *symtable) writeOpcode3(bop _Opcode, atoms []*parser.Node) (buf packet, yx uint16, err error) {
	// first atom: the op name, tail atoms: the args
	if len(atoms) > 4 {
		panic("shouldn't happen: too many arguments")
	}

	atoms = append([]*parser.Node{}, atoms...) // duplicate

	var n0, n1 *parser.Node

	if bop == OpLen || bop == OpLoad || bop == OpPop {
		if buf, err = table.decompound(atoms[1:]); err != nil {
			return
		}

		if len(atoms) >= 2 {
			n0 = atoms[1]
		}

		if len(atoms) == 3 {
			n1 = atoms[2]
		}

		err = table.writeOpcode(&buf, bop, n0, n1)
		for _, a := range atoms {
			if a.Source != "" {
				buf.WritePos(a.Meta)
				break
			}
		}
		return buf, regA, err
	}

	if bop == OpStore || bop == OpSlice {
		if buf, err = table.decompoundWithoutA(atoms[1:]); err != nil {
			return
		}

		if err = table.writeOpcode(&buf, OpSet, _nodeRegA, atoms[1]); err != nil {
			return
		}

		err = table.writeOpcode(&buf, bop, atoms[2], atoms[3])
		return buf, regA, err
	}

	buf, err = table.decompound(atoms[1:])
	if err != nil {
		return
	}

	switch bop {
	case OpTypeof, OpNot, OpBitNot, OpAddressOf, OpRet, OpYield:
		// unary op
		err = table.writeOpcode(&buf, bop, atoms[1], nil)
	case OpAssert:
		if len(atoms) == 3 {
			err = table.writeOpcode(&buf, bop, atoms[1], atoms[2])
		} else {
			err = table.writeOpcode(&buf, bop, atoms[1], nil)
		}
	default:
		// binary op
		err = table.writeOpcode(&buf, bop, atoms[1], atoms[2])
	}

	return buf, regA, err
}

func (table *symtable) compileFlatOp(atoms []*parser.Node) (code packet, yx uint16, err error) {
	head := atoms[0].Value.(string)
	switch head {
	case "nop":
		return newpacket(), regA, nil
	}
	op, ok := flatOpMapping[head]
	if !ok {
		err = fmt.Errorf("%+v: invalid op", atoms[0])
		return
	}
	code, yx, err = table.writeOpcode3(op, atoms)
	for _, a := range atoms {
		if a.Meta.Source != "" {
			code.WritePos(a.Meta)
			break
		}
	}
	return
}

func (table *symtable) compileIncOp(atoms []*parser.Node) (code packet, yx uint16, err error) {
	subject, ok := table.get(atoms[1].Value.(string))
	buf := newpacket()
	if !ok {
		return newpacket(), 0, fmt.Errorf(errUndeclaredVariable, atoms[1])
	}
	buf.WriteOP(OpInc, subject, table.loadK(&buf, atoms[2].Value))
	code.WritePos(atoms[1].Meta)
	return buf, regA, nil
}

// [and a b] => $a = a if not a then return else $a = b end
// [or a b]  => $a = a if a then do nothing else $a = b end
func (table *symtable) compileAndOrOp(atoms []*parser.Node) (code packet, yx uint16, err error) {
	bop := OpIfNot
	if atoms[0].Value.(string) == "or" {
		bop = OpIf
	}

	buf := newpacket()

	if err = table.writeOpcode(&buf, OpSet, _nodeRegA, atoms[1]); err != nil {
		return
	}
	buf.WriteOP(bop, regA, 0)
	c2 := buf.Len()

	if err = table.writeOpcode(&buf, OpSet, _nodeRegA, atoms[2]); err != nil {
		return
	}

	_, yx, _ = op(buf.data[c2-1])
	buf.data[c2-1] = makejmpop(bop, yx, buf.Len()-c2)
	buf.WritePos(atoms[0].Meta)
	return buf, regA, nil
}

// [if condition [true-chain ...] [false-chain ...]]
func (table *symtable) compileIfOp(atoms []*parser.Node) (code packet, yx uint16, err error) {
	condition := atoms[1]
	trueBranch, falseBranch := atoms[2], atoms[3]
	buf := newpacket()

	code, yx, err = table.compileNode(condition)
	if err != nil {
		return newpacket(), 0, err
	}
	buf.Write(code)
	condyx := yx

	var trueCode, falseCode packet
	trueCode, yx, err = table.compileChainOp(trueBranch)
	if err != nil {
		return
	}
	falseCode, yx, err = table.compileChainOp(falseBranch)
	if err != nil {
		return
	}
	if len(falseCode.data) > 0 {
		buf.WriteJmpOP(OpIfNot, condyx, len(trueCode.data)+1)
		buf.WritePos(atoms[0].Meta)
		buf.Write(trueCode)
		buf.WriteJmpOP(OpJmp, 0, len(falseCode.data))
		buf.Write(falseCode)
	} else {
		buf.WriteJmpOP(OpIfNot, condyx, len(trueCode.data))
		buf.WritePos(atoms[0].Meta)
		buf.Write(trueCode)
	}
	return buf, regA, nil
}

// [call callee [args ...]]
func (table *symtable) compileCallOp(nodes []*parser.Node) (code packet, yx uint16, err error) {
	tmp := append([]*parser.Node{nodes[1]}, nodes[2].C()...)
	code, err = table.decompound(tmp)
	if err != nil {
		return
	}

	var varIndex uint16
	var ok bool
	switch callee := tmp[0]; callee.Type {
	case parser.Natom:
		varIndex, ok = table.get(callee.Value.(string))
		if !ok {
			err = fmt.Errorf(errUndeclaredVariable, callee)
			return
		}
	case parser.Naddr:
		varIndex = callee.Value.(uint16)
	default:
		err = fmt.Errorf("%+v: invalid callee", callee)
		return
	}

	for i := 1; i < len(tmp); i++ {
		err = table.writeOpcode(&code, OpPush, tmp[i], nil)
		if err != nil {
			return
		}
		if tmp[i].Type == parser.Naddr {
			table.returnAddress(tmp[i].Value.(uint16))
		}
	}

	code.WriteOP(OpCall, varIndex, 0)
	code.WritePos(nodes[0].Meta)
	return code, regA, nil
}

// [lambda name? [namelist] [chain ...]]
func (table *symtable) compileLambdaOp(atoms []*parser.Node) (code packet, yx uint16, err error) {
	table.envescape = true
	isSafe, isVar := false, false
	for _, s := range strings.Split(atoms[0].Value.(string), ",") {
		switch s {
		case "safe":
			isSafe = true
		case "var":
			isVar = true
		}
	}

	name := atoms[1].Value.(string)
	newtable := newsymtable()
	newtable.parent = table

	params := atoms[2]
	if params.Type != parser.Ncompound {
		err = fmt.Errorf("%+v: invalid arguments list", atoms[2])
		return
	}

	for i, p := range params.C() {
		argname := p.Value.(string)
		if _, ok := newtable.sym[argname]; ok {
			return newpacket(), 0, fmt.Errorf("duplicated parameter: %s", argname)
		}
		if argname == "this" && i != 0 {
			return newpacket(), 0, fmt.Errorf("%+v: 'this' must be the first parameter inside a lambda", atoms[2])
		}
		newtable.put(argname, uint16(i))
	}

	ln := len(newtable.sym)
	if ln > 255 {
		return newpacket(), 0, fmt.Errorf("%+v: do you really need more than 255 arguments?", atoms[2])
	}

	newtable.vp = uint16(ln)

	if isVar {
		comps := append(atoms[3].C(), nil)
		copy(comps[2:], comps[1:])
		comps[1] = parser.CNode("set", "arguments", parser.CNode(
			"foreach", parser.ANodeS("nil"), parser.ANodeS("nil"),
		).SetPos0(atoms[0].Meta),
		).SetPos0(atoms[0].Meta)
		atoms[3].SetValue(comps)
	}

	code, yx, err = newtable.compileChainOp(atoms[3])
	if err != nil {
		return
	}

	code.WriteOP(OpEOB, 0, 0)
	buf := newpacket()
	cls := Closure{}
	cls.argsCount = byte(ln)
	if newtable.y || isSafe {
		cls.Set(ClsYieldable)
	}
	if _, ok := newtable.sym["this"]; ok {
		cls.Set(ClsHasReceiver)
	}
	if !newtable.envescape {
		cls.Set(ClsNoEnvescape)
	}
	if isSafe {
		cls.Set(ClsRecoverable)
	}

	buf.WriteOP(OpLambda, uint16(ln), uint16(cls.options))
	buf.Write32(uint32(len(newtable.consts)))
	buf.WriteConsts(newtable.consts)

	cls.code = code.data
	src := name + cls.String() + "@" + code.source
	if len(src) > 4095 {
		src = src[:4095]
	}

	// 26bit code pos length, 26bit code data length, 12bit code source length
	buf.Write64(uint64(uint32(len(code.pos))&0x03ffffff)<<38 +
		uint64(uint32(len(code.data))&0x03ffffff)<<12 +
		uint64(uint16(len(src))&0x0fff))
	buf.WriteRaw(u32FromBytes([]byte(src)))
	buf.WriteRaw(u32FromBytes(code.pos))

	// Note buf.source will be set to code.source in buf.Write
	// but buf.source shouldn't be changed
	code.source = buf.source
	buf.Write(code)
	buf.WritePos(atoms[0].Meta)
	return buf, regA, nil
}

var staticWhileHack [8]uint32

func init() {
	rand.Seed(time.Now().Unix())
	for i := range staticWhileHack {
		staticWhileHack[i] = rand.Uint32()
	}
}

// [continue | break]
func (table *symtable) compileContinueBreakOp(atoms []*parser.Node) (code packet, yx uint16, err error) {
	buf := newpacket()
	if len(table.continueNode) == 0 {
		err = fmt.Errorf("%+v: invalid statement outside a for loop", atoms[0])
		return
	}

	if atoms[0].Value.(string) == "continue" {
		cn := table.continueNode[len(table.continueNode)-1]
		code, yx, err = table.compileChainOp(cn)
		if err != nil {
			return
		}
		buf.Write(code)
		buf.WriteRaw(staticWhileHack[:4]) // write a 'continue' placeholder
		return buf, regA, nil
	}

	buf.WriteRaw(staticWhileHack[4:]) // write a 'continue' placeholder
	return buf, regA, nil
}

// [for condition incr [chain ...]]
func (table *symtable) compileWhileOp(atoms []*parser.Node) (code packet, yx uint16, err error) {
	condition := atoms[1]
	buf := newpacket()

	code, yx, err = table.compileNode(condition)
	if err != nil {
		return
	}
	buf.Write(code)
	cond := yx

	table.continueNode = append(table.continueNode, atoms[2])
	code, yx, err = table.compileChainOp(atoms[3])
	if err != nil {
		return
	}
	var icode packet
	icode, yx, err = table.compileChainOp(atoms[2])
	if err != nil {
		return
	}
	table.continueNode = table.continueNode[:len(table.continueNode)-1]

	code.Write(icode)
	buf.WriteJmpOP(OpIfNot, cond, len(code.data)+1) // +---> test cond  ---+
	buf.Write(code)                                 // |     code ...      | (break)
	buf.WriteJmpOP(OpJmp, 0, -buf.Len()-1)          // +---- continue      |
	//                                                   following code <--+

	code = buf
	code2 := u32Bytes(code.data)

	// search for special 'continue' placeholder and replace it with a OP_JMP to the
	// beginning of the code
	flag := u32Bytes(staticWhileHack[:4])
	for i := 0; i < len(code2); {
		x := bytes.Index(code2[i:], flag)
		if x == -1 {
			break
		}
		idx := (i + x) / 4
		code.data[idx] = makejmpop(OpJmp, 0, -idx-1)
		code.data[idx+1], code.data[idx+2], code.data[idx+3] = _nop, _nop, _nop
		i = idx*4 + 4
	}

	// search for special 'break' placeholder and replace it with a OP_JMP to the
	// end of the code
	flag = u32Bytes(staticWhileHack[4:])
	for i := 0; i < len(code2); {
		x := bytes.Index(code2[i:], flag)
		if x == -1 {
			break
		}
		idx := (i + x) / 4
		code.data[idx] = makejmpop(OpJmp, 0, len(code.data)-idx-1)
		code.data[idx+1], code.data[idx+2], code.data[idx+3] = _nop, _nop, _nop
		i = idx*4 + 4
	}
	return buf, regA, nil
}
