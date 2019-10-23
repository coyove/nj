package potatolang

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"strings"
	"sync"
	"unsafe"

	"github.com/coyove/potatolang/parser"
)

const (
	errUndeclaredVariable    = " %+v: undeclared variable"
	anonyMapIterCallbackFlag = "<anony-map-iter-callback>"
)

var _nop = makeop(OP_NOP, 0, 0)
var _nodeRegA = parser.NewNode(parser.Naddr).SetValue(regA)

func (table *symtable) compileSetOp(atoms []*parser.Node) (code packet, yx uint16, err error) {
	aDest, aSrc := atoms[1], atoms[2]
	buf := newpacket()
	var newYX uint16
	var ok bool
	var noNeedToRecordPos bool
	var varIndex uint16

	if atoms[0].Value.(string) == "set" {
		// compound has its own logic, we won't incr stack here
		if aSrc.Type != parser.Ncompound {
			newYX = table.vp
			table.incrvp()
		}
	} else {
		varIndex, ok = table.get(aDest.Value.(string))
		if !ok {
			err = fmt.Errorf(errUndeclaredVariable, aDest)
			return
		}
		newYX = varIndex
	}

	switch aSrc.Type {
	case parser.Natom:
		valueIndex, ok := table.get(aSrc.Value.(string))
		if !ok {
			err = fmt.Errorf(errUndeclaredVariable, aSrc)
			return
		}
		buf.WriteOP(OP_SET, newYX, valueIndex)
		noNeedToRecordPos = true
	case parser.Nnumber, parser.Nstring:
		buf.WriteOP(OP_SET, newYX, table.loadK(&buf, aSrc.Value))
		noNeedToRecordPos = true
	case parser.Ncompound:
		code, newYX, err = table.compileCompoundInto(aSrc, atoms[0].Value.(string) == "set", varIndex, false)
		if err != nil {
			return
		}
		buf.Write(code)
	}

	if atoms[0].Value.(string) == "set" {
		table.put(aDest.Value.(string), uint16(newYX))
	}
	if !noNeedToRecordPos {
		buf.WritePos(atoms[0].Meta)
	}
	return buf, newYX, nil
}

func (table *symtable) compileRetOp(atoms []*parser.Node) (code packet, yx uint16, err error) {
	isyield := atoms[0].S() == "yield"
	ispseudo := len(table.continueNode) > 0 && table.continueNode[len(table.continueNode)-1].S() == anonyMapIterCallbackFlag

	if isyield && ispseudo {
		err = fmt.Errorf("%+v: yield can't be used inside a pseudo foreach", atoms[0])
		return
	}

	var op byte
	op = OP_RET
	if isyield {
		op = OP_YIELD
		table.y = true
	}
	if ispseudo {
		// in a pseudo foreach, 'yield' is not allowed because we use them to simulate 'return'
		op = OP_YIELD
	}

	buf := newpacket()
	switch atom := atoms[1]; atom.Type {
	case parser.Natom, parser.Nnumber, parser.Nstring, parser.Naddr:
		if err = table.writeOpcode(&buf, op, atom, nil); err != nil {
			return
		}
	case parser.Ncompound:
		if code, yx, err = table.compileNode(atom); err != nil {
			return
		}
		buf.Write(code)
		buf.WriteOP(op, yx, 0)
	}
	buf.WritePos(atoms[0].Meta)
	return buf, yx, nil
}

func (table *symtable) compileMapArrayOp(atoms []*parser.Node) (code packet, yx uint16, err error) {
	var buf packet
	if buf, err = table.decompound(atoms[1].C()); err != nil {
		return
	}
	for _, atom := range atoms[1].C() {
		if err = table.writeOpcode(&buf, OP_PUSH, atom, nil); err != nil {
			return
		}
	}
	if atoms[0].Value.(string) == "map" {
		buf.WriteOP(OP_MAKEMAP, 0, 0)
	} else {
		buf.WriteOP(OP_MAKEMAP, 1, 0)
	}
	buf.WritePos(atoms[0].Meta)
	return buf, regA, nil
}

// writeOpcode3 accepts 3 arguments at most, 2 arguments will be encoded into opcode itself, the 3rd one will be in regA
func (table *symtable) writeOpcode3(atoms []*parser.Node, bop byte) (buf packet, yx uint16, err error) {
	// first atom: the op name, tail atoms: the args
	if len(atoms) > 4 {
		panic("shouldn't happen: too many arguments")
	}

	var n0, n1 *parser.Node

	if bop == OP_COPY {
		if err = table.writeOpcode(&buf, OP_SET, _nodeRegA, atoms[2]); err != nil {
			return
		}
		err = table.writeOpcode(&buf, bop, atoms[1], atoms[3])
		return buf, regA, err
	}

	if bop == OP_STORE || bop == OP_SLICE || bop == OP_LEN || bop == OP_LOAD || bop == OP_POP {
		if err = table.writeOpcode(&buf, OP_SET, _nodeRegA, atoms[1]); err != nil {
			return
		}

		if len(atoms) >= 3 {
			n0 = atoms[2]
		}

		if len(atoms) == 4 {
			n1 = atoms[3]
		}

		err = table.writeOpcode(&buf, bop, n0, n1)
		return buf, regA, err
	}

	buf, err = table.decompound(atoms[1:])
	if err != nil {
		return
	}

	switch bop {
	case OP_TYPEOF, OP_NOT, OP_ASSERT:
		err = table.writeOpcode(&buf, bop, atoms[1], nil)
	default:
		err = table.writeOpcode(&buf, bop, atoms[1], atoms[2])
	}

	return buf, regA, err
}

func (table *symtable) compileFlatOp(atoms []*parser.Node) (code packet, yx uint16, err error) {
	head := atoms[0].Value.(string)
	if head == "nop" {
		return newpacket(), regA, nil
	}
	op, ok := flatOpMapping[head]
	if !ok {
		err = fmt.Errorf("%+v: invalid op", atoms[0])
		return
	}
	code, yx, err = table.writeOpcode3(atoms, op)
	code.WritePos(atoms[0].Meta)
	return
}

func (table *symtable) compileIncOp(atoms []*parser.Node) (code packet, yx uint16, err error) {
	panic("TODO")
	//subject, ok := table.get(atoms[1].Value.(string))
	//buf := newpacket()
	//if !ok {
	//	return newpacket(), 0, fmt.Errorf(errUndeclaredVariable, atoms[1])
	//}
	//buf.WriteOP(OP_INC, subject, table.addConst(atoms[2].Value))
	//code.WritePos(atoms[1].Meta)
	//return buf, regA, nil
}

func checkjmpdist(jmp int) int {
	if jmp < -(1<<12) || jmp >= 1<<12 {
		panic("too long jump")
	}
	return jmp
}

// [and a b] => $a = a if not a then return else $a = b end
// [or a b]  => $a = a if a then do nothing else $a = b end
func (table *symtable) compileAndOrOp(atoms []*parser.Node) (code packet, yx uint16, err error) {
	bop := byte(OP_IFNOT)
	if atoms[0].Value.(string) == "or" {
		bop = OP_IF
	}

	buf := newpacket()

	if err = table.writeOpcode(&buf, OP_SET, _nodeRegA, atoms[1]); err != nil {
		return
	}
	buf.WriteOP(bop, regA, 0)
	c2 := buf.Len()

	if err = table.writeOpcode(&buf, OP_SET, _nodeRegA, atoms[2]); err != nil {
		return
	}
	jmp := checkjmpdist(buf.Len() - c2)

	_, yx, _ = op(buf.data[c2-1])
	buf.data[c2-1] = makeop(bop, yx, uint16(jmp+1<<12))
	buf.WritePos(atoms[0].Meta)
	return buf, regA, nil
}

// [if condition [true-chain ...] [false-chain ...]]
func (table *symtable) compileIfOp(atoms []*parser.Node) (code packet, yx uint16, err error) {
	condition := atoms[1]
	trueBranch, falseBranch := atoms[2], atoms[3]
	buf := newpacket()

	switch condition.Type {
	case parser.Nnumber, parser.Nstring:
		buf.WriteOP(OP_SET, regA, table.loadK(&buf, condition.Value))
		yx = regA
	case parser.Natom, parser.Ncompound:
		code, yx, err = table.compileNode(condition)
		if err != nil {
			return newpacket(), 0, err
		}
		buf.Write(code)
	}
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
		checkjmpdist(len(trueCode.data) + 1)
		buf.WriteOP(OP_IFNOT, condyx, uint16(len(trueCode.data)+1+1<<12))
		buf.WritePos(condition.Meta)
		buf.Write(trueCode)
		checkjmpdist(len(falseCode.data))
		buf.WriteOP(OP_JMP, 0, uint16(len(falseCode.data)+1<<12))
		buf.Write(falseCode)
	} else {
		checkjmpdist(len(trueCode.data))
		buf.WriteOP(OP_IFNOT, condyx, uint16(len(trueCode.data)+1<<12))
		buf.WritePos(condition.Meta)
		buf.Write(trueCode)
	}
	return buf, regA, nil
}

// [call func-name [args ...]]
func (table *symtable) compileCallOp(nodes []*parser.Node) (code packet, yx uint16, err error) {
	buf := newpacket()
	callee := nodes[1]
	name, _ := callee.Value.(string)

	if strings.HasPrefix(name, "$") {
		return table.compileRawOp(nodes)
	}

	switch name {
	case "addressof":
		varname := nodes[2].Cx(0).Value.(string)
		address, ok := table.get(varname)
		if !ok {
			err = fmt.Errorf(errUndeclaredVariable, callee)
			return
		}
		buf.WriteOP(OP_SET, regA, table.loadK(&buf, float64(address)))
		return buf, regA, nil
	case "copy":
		x := append([]*parser.Node{nodes[1]}, nodes[2].C()...)
		code, yx, err = table.writeOpcode3(x, OP_COPY)
		code.WritePos(nodes[0].Meta)
		return
	}

	atoms, replacedAtoms := nodes[2].C(), []*parser.Node{}
	for i := 0; i < len(atoms); i++ {
		atom := atoms[i]

		if atom.Type == parser.Ncompound {
			code, yx, err = table.compileCompoundInto(atom, true, 0, true)
			if err != nil {
				return
			}
			atoms[i] = parser.NewNode(parser.Naddr).SetValue(yx)
			replacedAtoms = append(replacedAtoms, atoms[i])
			buf.Write(code)
		}
	}

	// note: [call [..] [..]] is different, which will be handled in more belowed code
	if callee.Type != parser.Ncompound {
		if len(replacedAtoms) >= 1 {
			var tmp uint16
			_, tmp, replacedAtoms[len(replacedAtoms)-1].Value = op(buf.data[len(buf.data)-1])
			buf.TruncateLast(1)
			//table.vp--
			table.returnTmp(tmp)
		}
	}

	var varIndex uint16
	var ok bool
	switch callee.Type {
	case parser.Natom:
		varIndex, ok = table.get(callee.Value.(string))
		if !ok {
			err = fmt.Errorf(errUndeclaredVariable, callee)
			return
		}
	case parser.Ncompound:
		code, yx, err = table.compileCompoundInto(callee, true, 0, false)
		if err != nil {
			return
		}
		varIndex = yx
		if len(replacedAtoms) == 0 {
			_, _, varIndex = op(code.data[len(code.data)-1])
			code.data = code.data[:len(code.data)-1]
			table.vp--
		}
		buf.Write(code)
	case parser.Naddr:
		varIndex = callee.Value.(uint16)
	default:
		err = fmt.Errorf("%+v: invalid callee", callee)
		return
	}

	for i := 0; i < len(atoms); i++ {
		err = table.writeOpcode(&buf, OP_PUSH, atoms[i], nil)
		if err != nil {
			return
		}
		if atoms[i].Type == parser.Naddr {
			table.returnTmp(atoms[i].Value.(uint16))
		}
	}

	buf.WriteOP(OP_CALL, varIndex, 0)
	buf.WritePos(nodes[0].Meta)
	return buf, regA, nil
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
		err = fmt.Errorf("%+v: invalid arguments list", atoms[0])
		return
	}

	var this bool
	i := 0
	for _, p := range params.C() {
		argname := p.Value.(string)
		if argname == "this" {
			this = true
			continue
		}
		newtable.put(argname, uint16(i))
		i++
	}

	ln := len(newtable.sym)
	if this {
		newtable.put("this", uint16(ln))
	}

	if ln > 255 {
		return newpacket(), 0, fmt.Errorf("do you really need more than 255 arguments?")
	}

	// this is a special function, inside it any 'continue' will be converted to 'return nil'
	// and any 'break' will be converted to 'return #nil'
	if name == anonyMapIterCallbackFlag {
		newtable.continueNode = append(newtable.continueNode, parser.NewNode(parser.Natom).SetValue(name))
	}

	newtable.vp = uint16(ln)

	if isVar {
		comps := append(atoms[3].C(), nil)
		copy(comps[2:], comps[1:])
		comps[1] = parser.CNode("set", "arguments", parser.CNode(
			"call", "copy", parser.CNode(parser.NNode(2), parser.NilNode(), parser.NilNode()),
		))
		atoms[3].SetValue(comps)
	}

	code, yx, err = newtable.compileChainOp(atoms[3])
	if err != nil {
		return
	}

	if name == anonyMapIterCallbackFlag {
		newtable.continueNode = newtable.continueNode[:len(newtable.continueNode)-1]
	}

	code.WriteOP(OP_EOB, 0, 0)
	buf := newpacket()
	cls := Closure{}
	cls.argsCount = byte(ln)
	if newtable.y || isSafe {
		cls.Set(CLS_YIELDABLE)
	}
	if this {
		cls.Set(CLS_HASRECEIVER)
	}
	if !newtable.envescape {
		cls.Set(CLS_NOENVESCAPE)
	}
	if isSafe {
		cls.Set(CLS_RECOVERALL)
	}
	if name == anonyMapIterCallbackFlag {
		cls.Set(CLS_PSEUDO_FOREACH)
	}

	buf.WriteOP(OP_LAMBDA, uint16(ln), uint16(cls.options))
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

var staticWhileHack struct {
	sync.Mutex
	continueFlag []uint32
	breakFlag    []uint32
}

func gen128bit() ([4]uint32, []uint32) {
	var p [4]uint32
	buf := [16]byte{}
	if _, err := rand.Read(buf[:]); err != nil {
		panic(err)
	}
	copy(p[:], (*(*[4]uint32)(unsafe.Pointer(&buf)))[:])
	return p, p[:]
}

// [continue | break]
func (table *symtable) compileContinueBreakOp(atoms []*parser.Node) (code packet, yx uint16, err error) {
	staticWhileHack.Lock()
	defer staticWhileHack.Unlock()

	buf := newpacket()
	if atoms[0].Value.(string) == "continue" {
		if staticWhileHack.continueFlag == nil {
			_, staticWhileHack.continueFlag = gen128bit()
		}
		if len(table.continueNode) == 0 {
			err = fmt.Errorf("%+v: invalid continue statement", atoms[0])
			return
		}
		cn := table.continueNode[len(table.continueNode)-1]
		if cn.S() == anonyMapIterCallbackFlag {
			buf.WriteOP(OP_RET, 0, 0)
			return buf, regA, nil
		}
		code, yx, err = table.compileChainOp(cn)
		if err != nil {
			return
		}
		buf.Write(code)
		buf.WriteRaw(staticWhileHack.continueFlag)
		return buf, regA, nil
	}
	if staticWhileHack.breakFlag == nil {
		_, staticWhileHack.breakFlag = gen128bit()
	}
	if len(table.continueNode) == 0 {
		err = fmt.Errorf("%+v: invalid break statement", atoms[0])
		return
	}
	if table.continueNode[len(table.continueNode)-1].S() == anonyMapIterCallbackFlag {
		buf.WriteOP(OP_SET, regA, table.getnil())
		buf.WriteOP(OP_POP, 0, 0)
		buf.WriteOP(OP_RET, regA, 0) // PhantomValue
		return buf, regA, nil
	}
	buf.WriteRaw(staticWhileHack.breakFlag)
	return buf, regA, nil
}

// [for condition incr [chain ...]]
func (table *symtable) compileWhileOp(atoms []*parser.Node) (code packet, yx uint16, err error) {
	condition := atoms[1]
	buf := newpacket()
	var varIndex uint16

	switch condition.Type {
	case parser.Naddr:
		varIndex = condition.Value.(uint16)
	case parser.Nnumber, parser.Nstring:
		varIndex = table.vp
		table.incrvp()
		buf.WriteOP(OP_SET, varIndex, table.loadK(&buf, condition.Value))
	case parser.Ncompound, parser.Natom:
		code, yx, err = table.compileNode(condition)
		if err != nil {
			return
		}
		buf.Write(code)
		varIndex = yx
	}

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
	checkjmpdist(len(code.data) + 1)
	buf.WriteOP(OP_IFNOT, varIndex, uint16(len(code.data)+1+1<<12))
	buf.Write(code)
	checkjmpdist(-buf.Len() - 1)
	buf.WriteOP(OP_JMP, 0, uint16(-buf.Len()-1+1<<12))

	code = buf
	code2 := u32Bytes(code.data)
	if staticWhileHack.continueFlag != nil {
		flag := u32Bytes(staticWhileHack.continueFlag)
		for i := 0; i < len(code2); {
			x := bytes.Index(code2[i:], flag)
			if x == -1 {
				break
			}
			idx := (i + x) / 4
			checkjmpdist(-idx - 1)
			code.data[idx] = makeop(OP_JMP, 0, uint16(-idx-1+1<<12))
			code.data[idx+1], code.data[idx+2], code.data[idx+3] = _nop, _nop, _nop
			i = idx*4 + 4
		}
	}

	if staticWhileHack.breakFlag != nil {
		flag := u32Bytes(staticWhileHack.breakFlag)
		for i := 0; i < len(code2); {
			x := bytes.Index(code2[i:], flag)
			if x == -1 {
				break
			}
			idx := (i + x) / 4
			checkjmpdist(len(code.data) - idx - 1)
			code.data[idx] = makeop(OP_JMP, 0, uint16(len(code.data)-idx-1+1<<12))
			code.data[idx+1], code.data[idx+2], code.data[idx+3] = _nop, _nop, _nop
			i = idx*4 + 4
		}
	}
	return buf, regA, nil
}
