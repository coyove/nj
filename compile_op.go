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
		table.clearRegRecord(varIndex)
	}

	switch aSrc.Type {
	case parser.Natom:
		if aSrc.Value.(string) == "nil" {
			buf.WriteOP(OP_SETK, newYX, 0)
		} else {
			valueIndex, ok := table.get(aSrc.Value.(string))
			if !ok {
				err = fmt.Errorf(errUndeclaredVariable, aSrc)
				return
			}
			buf.WriteOP(OP_SET, newYX, valueIndex)
		}
		noNeedToRecordPos = true
	case parser.Nnumber, parser.Nstring:
		buf.WriteOP(OP_SETK, newYX, table.addConst(aSrc.Value))
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

	var op, opk byte
	op, opk = OP_RET, OP_RETK
	if isyield {
		op, opk = OP_YIELD, OP_YIELDK
		table.y = true
	}
	if ispseudo {
		// in a pseudo foreach, 'yield' is not allowed because we use them to simulate 'return'
		op, opk = OP_YIELD, OP_YIELDK
	}

	buf := newpacket()
	switch atom := atoms[1]; atom.Type {
	case parser.Natom, parser.Nnumber, parser.Nstring, parser.Naddr:
		if err = table.fill(&buf, atom, op, opk); err != nil {
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
	if buf, err = table.decompound(atoms[1].C(), nil, true); err != nil {
		return
	}
	for _, atom := range atoms[1].C() {
		if err = table.fill(&buf, atom, OP_PUSH, OP_PUSHK); err != nil {
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

func (table *symtable) flatWrite(atoms []*parser.Node, bop byte) (code packet, yx uint16, err error) {
	var op = []uint16{OP_R0<<8 | OP_R0K, OP_R1<<8 | OP_R1K, OP_R2<<8 | OP_R2K, OP_R3<<8 | OP_R3K}
	var canWeUseR2 = true

	switch bop {
	case OP_STORE, OP_LOAD, OP_SLICE:
		canWeUseR2 = false
		fallthrough
	case OP_LEN, OP_POP:
		op[0], op[3] = op[3], op[0]
		op[1], op[2] = op[2], op[1]
	}

	var buf packet
	buf, err = table.decompound(atoms[1:], op, canWeUseR2)
	if err != nil {
		return
	}

	immediateRet := func(n float64) {
		buf.Clear()
		buf.WriteOP(OP_SETK, regA, table.addConst(n))
		table.im = &n
	}

	switch bop {
	case OP_ADD:
		if atoms[1].Type == parser.Nstring && atoms[2].Type == parser.Nstring {
			str := atoms[1].Value.(string) + atoms[2].Value.(string)
			buf.Clear()
			buf.WriteOP(OP_SETK, regA, table.addConst(str))
			table.ims = &str
			return buf, regA, nil
		}
		if atoms[1].Type == parser.Nnumber && atoms[2].Type == parser.Nnumber {
			immediateRet(atoms[1].Value.(float64) + atoms[2].Value.(float64))
			return buf, regA, nil
		}
	}

	if len(atoms) == 2 && atoms[1].Type == parser.Nnumber {
		v1 := atoms[1].Value.(float64)
		switch bop {
		case OP_NOT:
			if v1 == 0 {
				immediateRet(1)
			} else {
				immediateRet(0)
			}
			return buf, regA, nil
		case OP_BIT_NOT:
			immediateRet(float64(^int32(v1)))
			return buf, regA, nil
		}
	}

	if len(atoms) > 2 && atoms[1].Type == parser.Nnumber && atoms[2].Type == parser.Nnumber {
		v1, v2 := atoms[1].Value.(float64), atoms[2].Value.(float64)
		switch bop {
		case OP_SUB:
			immediateRet(v1 - v2)
		case OP_MUL:
			immediateRet(v1 * v2)
		case OP_DIV:
			immediateRet(v1 / v2)
		case OP_MOD:
			immediateRet(float64(int64(v1) % int64(v2)))
		case OP_EQ:
			if v1 == v2 {
				immediateRet(1)
			} else {
				immediateRet(0)
			}
		case OP_NEQ:
			if v1 != v2 {
				immediateRet(1)
			} else {
				immediateRet(0)
			}
		case OP_LESS:
			if v1 < v2 {
				immediateRet(1)
			} else {
				immediateRet(0)
			}
		case OP_LESS_EQ:
			if v1 <= v2 {
				immediateRet(1)
			} else {
				immediateRet(0)
			}
		case OP_BIT_AND:
			immediateRet(float64(int32(v1) & int32(v2)))
		case OP_BIT_OR:
			immediateRet(float64(int32(v1) | int32(v2)))
		case OP_BIT_XOR:
			immediateRet(float64(int32(v1) ^ int32(v2)))
		case OP_BIT_LSH:
			immediateRet(float64(int32(v1) << uint32(v2)))
		case OP_BIT_RSH:
			immediateRet(float64(int32(v1) >> uint32(v2)))
		case OP_BIT_URSH:
			immediateRet(float64(uint32(v1) >> uint32(v2)))
		default:
			goto IM_PASS
		}
		return buf, regA, nil
	}

IM_PASS:

	count := buf.Len()
	for i := 1; i < len(atoms); i++ {
		if op[i-1] == OP_NOP {
			// ignore me ...
			continue
		}
		if err = table.fill(&buf, atoms[i], byte(op[i-1]>>8), byte(op[i-1])); err != nil {
			return
		}
	}
	if buf.Len() == count {
		// TODO:
		// No argument opcode was written into buf, which means either this op accepts 0 argument,
		// or all arguments have been inside the registers already.
		// So we should evaluate the case that A still holds the result and skip the current op, e.g.:

		// set a = (i + j) / (i + j + 1) =>
		// R0 = i, R1 = j, A = i + j
	}
	buf.WriteOP(bop, 0, 0)

	return buf, regA, nil
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
	code, yx, err = table.flatWrite(atoms, op)
	code.WritePos(atoms[0].Meta)
	return
}

func (table *symtable) compileIncOp(atoms []*parser.Node) (code packet, yx uint16, err error) {
	subject, ok := table.get(atoms[1].Value.(string))
	buf := newpacket()
	if !ok {
		return newpacket(), 0, fmt.Errorf(errUndeclaredVariable, atoms[1])
	}
	table.clearRegRecord(subject)
	buf.WriteOP(OP_INC, subject, table.addConst(atoms[2].Value))
	code.WritePos(atoms[1].Meta)
	return buf, regA, nil
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
	a1, a2 := atoms[1], atoms[2]
	buf := newpacket()
	switch a1.Type {
	case parser.Nnumber, parser.Nstring:
		buf.WriteOP(OP_SETK, regA, table.addConst(a1.Value))
		buf.WriteOP(bop, regA, 0)
	case parser.Natom, parser.Ncompound:
		code, yx, err = table.compileNode(a1)
		if err != nil {
			return newpacket(), 0, err
		}
		buf.Write(code)
		buf.WriteOP(OP_SET, regA, yx)
		buf.WriteOP(bop, regA, 0)
	}
	c2 := buf.Len()

	switch a2.Type {
	case parser.Nnumber, parser.Nstring:
		buf.WriteOP(OP_SETK, regA, table.addConst(a2.Value))
	case parser.Natom, parser.Ncompound:
		code, yx, err = table.compileNode(a2)
		if err != nil {
			return newpacket(), 0, err
		}
		buf.Write(code)
		if yx != regA {
			buf.WriteOP(OP_SET, regA, yx)
		}
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
		buf.WriteOP(OP_SETK, regA, table.addConst(condition.Value))
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
	table.clearAllRegRecords()
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
		buf.WriteOP(OP_SETK, regA, table.addConst(float64(address)))
		return buf, regA, nil
	case "copy":
		x := append([]*parser.Node{nodes[1]}, nodes[2].C()...)
		code, yx, err = table.flatWrite(x, OP_COPY)
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
		err = table.fill(&buf, atoms[i], OP_PUSH, OP_PUSHK)
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
			buf.WriteOP(OP_RETK, 0, 0)
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
		buf.WriteOP(OP_R3K, 0, 0)
		buf.WriteOP(OP_POP, 0, 0)
		buf.WriteOP(OP_RET, regA, 0)
		return buf, regA, nil
	}
	buf.WriteRaw(staticWhileHack.breakFlag)
	return buf, regA, nil
}

// [for condition incr [chain ...]]
func (table *symtable) compileWhileOp(atoms []*parser.Node) (code packet, yx uint16, err error) {
	table.clearAllRegRecords()
	condition := atoms[1]
	buf := newpacket()
	var varIndex uint16

	switch condition.Type {
	case parser.Naddr:
		varIndex = condition.Value.(uint16)
	case parser.Nnumber, parser.Nstring:
		varIndex = table.vp
		table.incrvp()
		buf.WriteOP(OP_SETK, varIndex, table.addConst(condition.Value))
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
	table.clearAllRegRecords()
	return buf, regA, nil
}
