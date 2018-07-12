package potatolang

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"sync"
	"unsafe"

	"github.com/coyove/potatolang/parser"
)

const (
	errUndeclaredVariable = " %+v: undeclared variable"
)

func (table *symtable) compileSetOp(atoms []*parser.Node) (code packet, yx uint32, err error) {
	aDest, aSrc := atoms[1], atoms[2]
	varIndex := uint32(0)
	buf := newpacket()
	var newYX uint32
	var ok bool

	if atoms[0].Value.(string) == "set" {
		// compound has its own logic, we won't incr stack here
		if aSrc.Type != parser.NTCompound {
			newYX = uint32(table.sp)
			table.sp++
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
	case parser.NTAtom:
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
	case parser.NTNumber, parser.NTString:
		buf.WriteOP(OP_SETK, newYX, uint32(table.addConst(aSrc.Value)))
	case parser.NTCompound:
		code, newYX, err = table.compileCompoundInto(aSrc, atoms[0].Value.(string) == "set", varIndex)
		if err != nil {
			return
		}
		buf.Write(code)
	}

	if atoms[0].Value.(string) == "set" {
		_, redecl := table.get(aDest.Value.(string))
		if redecl && table.noredecl {
			err = fmt.Errorf("redeclare: %+v", aDest)
			return
		}
		table.put(aDest.Value.(string), uint16(newYX))
	}
	buf.WritePos(atoms[0].Pos)
	return buf, newYX, nil
}

func (table *symtable) compileRetOp(atoms []*parser.Node) (code packet, yx uint32, err error) {
	var op, opk byte
	op, opk = OP_RET, OP_RETK
	if atoms[0].Value.(string) == "yield" {
		op, opk = OP_YIELD, OP_YIELDK
		table.y = true
	}

	buf := newpacket()
	if len(atoms) == 1 {
		buf.WriteOP(op, regA, 0)
		return buf, yx, nil
	}

	switch atom := atoms[1]; atom.Type {
	case parser.NTAtom, parser.NTNumber, parser.NTString, parser.NTAddr:
		if err = table.fill(&buf, atom, op, opk); err != nil {
			return
		}
	case parser.NTCompound:
		if code, yx, err = table.compileNode(atom); err != nil {
			return
		}
		buf.Write(code)
		buf.WriteOP(op, yx, 0)
	}
	buf.WritePos(atoms[0].Pos)
	return buf, yx, nil
}

func (table *symtable) compileMapOp(atoms []*parser.Node) (code packet, yx uint32, err error) {
	if len(atoms[1].Compound)%2 != 0 {
		err = fmt.Errorf("%+v: every key in map must have a value", atoms[1])
		return
	}
	var buf packet
	if buf, err = table.flaten(atoms[1].Compound); err != nil {
		return
	}
	for _, atom := range atoms[1].Compound {
		if err = table.fill(&buf, atom, OP_PUSH, OP_PUSHK); err != nil {
			return
		}
	}
	buf.WriteOP(OP_MAKEMAP, 0, 0)
	buf.WritePos(atoms[0].Pos)
	return buf, regA, nil
}

func (table *symtable) flaten(atoms []*parser.Node) (buf packet, err error) {
	replacedAtoms := []*parser.Node{}
	buf = newpacket()

	for i, atom := range atoms {
		var yx uint32
		var code packet

		if atom.Type == parser.NTCompound {
			if code, yx, err = table.compileCompoundInto(atom, true, 0); err != nil {
				return
			}
			if table.im != nil {
				atoms[i] = &parser.Node{Type: parser.NTNumber, Value: *table.im}
				table.im = nil
			} else if table.ims != nil {
				atoms[i] = &parser.Node{Type: parser.NTString, Value: *table.ims}
				table.ims = nil
			} else {
				atoms[i] = &parser.Node{Type: parser.NTAddr, Value: yx}
				replacedAtoms = append(replacedAtoms, atoms[i])
				buf.Write(code)
			}
		}
	}

	if len(replacedAtoms) > 0 {
		_, _, replacedAtoms[len(replacedAtoms)-1].Value = op(buf.data[len(buf.data)-1])
		buf.TruncateLast(1)
		table.sp--
	}

	return buf, nil
}

func (table *symtable) flatWrite(atoms []*parser.Node, bop byte) (code packet, yx uint32, err error) {
	var buf packet
	buf, err = table.flaten(atoms[1:])
	if err != nil {
		return
	}

	immediateRet := func(n float64) {
		buf.Clear()
		buf.WriteOP(OP_SETK, regA, uint32(table.addConst(n)))
		table.im = &n
	}

	v1 := func() float64 { return atoms[1].Value.(float64) }
	v2 := func() float64 { return atoms[2].Value.(float64) }

	switch bop {
	case OP_ADD:
		if atoms[1].Type == parser.NTString && atoms[2].Type == parser.NTString {
			str := atoms[1].Value.(string) + atoms[2].Value.(string)
			buf.Clear()
			buf.WriteOP(OP_SETK, regA, uint32(table.addConst(str)))
			table.ims = &str
			return buf, regA, nil
		}
		if atoms[1].Type == parser.NTNumber && atoms[2].Type == parser.NTNumber {
			immediateRet(v1() + v2())
			return buf, regA, nil
		}
	case OP_SUB:
		if atoms[1].Type == parser.NTNumber && atoms[2].Type == parser.NTNumber {
			immediateRet(v1() - v2())
			return buf, regA, nil
		}
	case OP_MUL:
		if atoms[1].Type == parser.NTNumber && atoms[2].Type == parser.NTNumber {
			immediateRet(v1() * v2())
			return buf, regA, nil
		}
	case OP_DIV:
		if atoms[1].Type == parser.NTNumber && atoms[2].Type == parser.NTNumber {
			immediateRet(v1() / v2())
			return buf, regA, nil
		}
	}

	var op = [4]uint16{OP_R0<<8 | OP_R0K, OP_R1<<8 | OP_R1K, OP_R2<<8 | OP_R2K, OP_R3<<8 | OP_R3K}
	switch bop {
	case OP_LEN, OP_STORE, OP_LOAD, OP_SLICE, OP_POP:
		op[0], op[3] = op[3], op[0]
		op[1], op[2] = op[2], op[1]
	}

	count := buf.Len()
	for i := 1; i < len(atoms); i++ {
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

func (table *symtable) compileFlatOp(atoms []*parser.Node) (code packet, yx uint32, err error) {
	head := atoms[0].Value.(string)
	op := flatOpMapping[head]
	if op == 0 {
		err = fmt.Errorf("%+v: invalid op", atoms[0])
		return
	}
	code, yx, err = table.flatWrite(atoms, op)
	code.WritePos(atoms[0].Pos)
	return
}

func (table *symtable) compileIncOp(atoms []*parser.Node) (code packet, yx uint32, err error) {
	subject, ok := table.get(atoms[1].Value.(string))
	buf := newpacket()
	if !ok {
		return newpacket(), 0, fmt.Errorf(errUndeclaredVariable, atoms[1])
	}
	table.clearRegRecord(subject)
	buf.WriteOP(OP_INC, subject, uint32(table.addConst(atoms[2].Value)))
	code.WritePos(atoms[1].Pos)
	return buf, regA, nil
}

// [and a b] => $a = a if not a then return else $a = b end
// [or a b]  => $a = a if a then do nothing else $a = b end
func (table *symtable) compileAndOrOp(atoms []*parser.Node) (code packet, yx uint32, err error) {
	bop := byte(OP_IFNOT)
	if atoms[0].Value.(string) == "or" {
		bop = OP_IF
	}
	a1, a2 := atoms[1], atoms[2]
	buf := newpacket()
	switch a1.Type {
	case parser.NTNumber, parser.NTString:
		buf.WriteOP(OP_SETK, regA, uint32(table.addConst(a1.Value)))
		buf.WriteOP(bop, regA, 0)
	case parser.NTAtom, parser.NTCompound:
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
	case parser.NTNumber, parser.NTString:
		buf.WriteOP(OP_SETK, regA, uint32(table.addConst(a2.Value)))
	case parser.NTAtom, parser.NTCompound:
		code, yx, err = table.compileNode(a2)
		if err != nil {
			return newpacket(), 0, err
		}
		buf.Write(code)
		if yx != regA {
			buf.WriteOP(OP_SET, regA, yx)
		}
	}
	jmp := buf.Len() - c2

	_, yx, _ = op(buf.data[c2-1])
	buf.data[c2-1] = makeop(bop, yx, uint32(int32(jmp)))
	code.WritePos(atoms[0].Pos)
	return buf, regA, nil
}

// [if condition [true-chain ...] [false-chain ...]]
func (table *symtable) compileIfOp(atoms []*parser.Node) (code packet, yx uint32, err error) {
	condition := atoms[1]
	trueBranch, falseBranch := atoms[2], atoms[3]
	buf := newpacket()

	switch condition.Type {
	case parser.NTNumber, parser.NTString:
		buf.WriteOP(OP_SETK, regA, uint32(table.addConst(condition.Value)))
		yx = regA
	case parser.NTAtom, parser.NTCompound:
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
		buf.WriteOP(OP_IFNOT, condyx, uint32(len(trueCode.data))+1)
		buf.WritePos(condition.Pos)
		buf.Write(trueCode)
		buf.WriteOP(OP_JMP, 0, uint32(len(falseCode.data)))
		buf.Write(falseCode)
	} else {
		buf.WriteOP(OP_IFNOT, condyx, uint32(len(trueCode.data)))
		buf.WritePos(condition.Pos)
		buf.Write(trueCode)
	}
	table.clearAllRegRecords()
	return buf, regA, nil
}

// [call func-name [args ...]]
func (table *symtable) compileCallOp(nodes []*parser.Node) (code packet, yx uint32, err error) {
	buf := newpacket()
	callee := nodes[1]
	name, _ := callee.Value.(string)
	switch name {
	case "addressof":
		varname := nodes[2].Compound[0].Value.(string)
		address, ok := table.get(varname)
		if !ok {
			err = fmt.Errorf(errUndeclaredVariable, callee)
			return
		}
		buf.WriteOP(OP_SETK, regA, uint32(table.addConst(float64(address))))
		return buf, regA, nil
	case "len":
		code, yx, err = table.flatWrite(append(nodes[1:2], nodes[2].Compound...), OP_LEN)
		code.WritePos(nodes[0].Pos)
		return
	case "copy":
		x := append(nodes[1:2], nodes[2].Compound...)
		if y, ok := x[3].Value.(float64); ok && y == 2 {
			// return stack, env is escaped
			table.envescape = true
		}
		code, yx, err = table.flatWrite(x, OP_COPY)
		code.WritePos(nodes[0].Pos)
		return
	case "typeof":
		return table.flatWrite(append(nodes[1:2], nodes[2].Compound...), OP_TYPEOF)
	}

	atoms, replacedAtoms := nodes[2].Compound, []*parser.Node{}
	for i := 0; i < len(atoms); i++ {
		atom := atoms[i]

		if atom.Type == parser.NTCompound {
			code, yx, err = table.compileCompoundInto(atom, true, 0)
			if err != nil {
				return
			}
			atoms[i] = &parser.Node{Type: parser.NTAddr, Value: yx}
			replacedAtoms = append(replacedAtoms, atoms[i])
			buf.Write(code)
		}
	}

	// note: [call [..] [..]] is different
	if len(replacedAtoms) == 1 && callee.Type != parser.NTCompound {
		_, _, replacedAtoms[0].Value = op(buf.data[len(buf.data)-1])
		buf.TruncateLast(1)
		table.sp--
	}

	var varIndex uint32
	var ok bool
	switch callee.Type {
	case parser.NTAtom:
		varIndex, ok = table.get(callee.Value.(string))
		if !ok {
			err = fmt.Errorf(errUndeclaredVariable, callee)
			return
		}
	case parser.NTCompound:
		code, yx, err = table.compileCompoundInto(callee, true, 0)
		if err != nil {
			return
		}
		varIndex = yx
		if len(replacedAtoms) == 0 {
			_, _, varIndex = op(code.data[len(code.data)-1])
			code.data = code.data[:len(code.data)-1]
			table.sp--
		}
		buf.Write(code)
	case parser.NTAddr:
		varIndex = callee.Value.(uint32)
	default:
		err = fmt.Errorf("invalid callee: %+v", callee)
		return
	}

	for i := 0; i < len(atoms); i++ {
		err = table.fill(&buf, atoms[i], OP_PUSH, OP_PUSHK)
		if err != nil {
			return
		}
	}

	buf.WriteOP(OP_CALL, varIndex, 0)
	buf.WritePos(nodes[0].Pos)
	return buf, regA, nil
}

// [lambda [namelist] [chain ...]]
func (table *symtable) compileLambdaOp(atoms []*parser.Node) (code packet, yx uint32, err error) {
	table.envescape = true
	isSafe := atoms[0].Value.(string) == "safefunc"
	newtable := newsymtable()
	newtable.parent = table
	newtable.noredecl = table.noredecl

	params := atoms[1]
	if params.Type != parser.NTCompound {
		err = fmt.Errorf("%+v: invalid arguments list", atoms[0])
		return
	}

	var this bool
	i := 0
	for _, p := range params.Compound {
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

	newtable.sp = uint16(ln)
	code, yx, err = newtable.compileChainOp(atoms[2])
	if err != nil {
		return
	}
	if len(code.source) > 4096 {
		return newpacket(), 0, fmt.Errorf("does your path really contain more than 4096 chars?")
	}

	code.WriteOP(OP_EOB, 0, 0)
	buf := newpacket()
	cls := Closure{}
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

	buf.WriteOP(OP_LAMBDA, uint32(len(newtable.consts)), uint32(byte(ln))<<24+uint32(cls.options))
	buf.WriteConsts(newtable.consts)
	code.data = patchGotoCode(code.data, newtable.gotoTable)
	buf.WriteCode(code)
	buf.WritePos(atoms[0].Pos)
	return buf, regA, nil
}

var staticWhileHack struct {
	sync.Mutex
	continueFlag []uint64
	breakFlag    []uint64
}

func gen128bit() ([2]uint64, []uint64) {
	var p [2]uint64
	buf := [16]byte{}
	if _, err := rand.Read(buf[:]); err != nil {
		panic(err)
	}
	copy(p[:], (*(*[2]uint64)(unsafe.Pointer(&buf)))[:])
	return p, p[:]
}

// [continue | break]
func (table *symtable) compileContinueBreakOp(atoms []*parser.Node) (code packet, yx uint32, err error) {
	staticWhileHack.Lock()
	defer staticWhileHack.Unlock()

	buf := newpacket()
	if atoms[0].Value.(string) == "continue" {
		if staticWhileHack.continueFlag == nil {
			_, staticWhileHack.continueFlag = gen128bit()
		}
		code, yx, err = table.compileChainOp(table.continueNode[len(table.continueNode)-1])
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
	buf.WriteRaw(staticWhileHack.breakFlag)
	return buf, regA, nil
}

func (table *symtable) compileLabelGotoOp(atoms []*parser.Node) (code packet, yx uint32, err error) {
	buf := newpacket()
	label := atoms[1].Value.(string)

	if atoms[0].Value.(string) == "label" {
		if _, exist := table.gotoTable[label]; exist {
			return buf, 0, fmt.Errorf("label '%s' already exists", label)
		}
		x, _ := gen128bit() // placeholder
		table.gotoTable[label] = x
		buf.WriteRaw(x[:])
		return buf, regA, nil
	}

	x, exist := table.gotoTable[label]
	if !exist {
		return buf, 0, fmt.Errorf("label '%s' doesn't exist", label)
	}
	buf.WriteOP(OP_JMP, 0xffffffff, 0)
	buf.WriteRaw(x[:])
	return buf, regA, nil
}

// [for condition incr [chain ...]]
func (table *symtable) compileWhileOp(atoms []*parser.Node) (code packet, yx uint32, err error) {
	table.clearAllRegRecords()
	condition := atoms[1]
	buf := newpacket()
	var varIndex uint32

	switch condition.Type {
	case parser.NTAddr:
		varIndex = condition.Value.(uint32)
	case parser.NTNumber, parser.NTString:
		varIndex = uint32(table.sp)
		table.sp++
		buf.WriteOP(OP_SETK, varIndex, uint32(table.addConst(condition.Value)))
	case parser.NTCompound, parser.NTAtom:
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
	buf.WriteOP(OP_IFNOT, varIndex, uint32(len(code.data))+1)
	buf.Write(code)
	buf.WriteOP(OP_JMP, 0, -uint32(buf.Len())-1)

	code = buf
	code2 := slice64to8(code.data)
	if staticWhileHack.continueFlag != nil {
		flag := slice64to8(staticWhileHack.continueFlag)
		for i := 0; i < len(code2); {
			x := bytes.Index(code2[i:], flag)
			if x == -1 {
				break
			}
			idx := (i + x) / 8
			code.data[idx] = makeop(OP_JMP, 0, uint32(int32(-idx-1)))
			code.data[idx+1] = makeop(OP_NOP, 0, 0)
			i = idx*8 + 2
		}
	}

	if staticWhileHack.breakFlag != nil {
		flag := slice64to8(staticWhileHack.breakFlag)
		for i := 0; i < len(code2); {
			x := bytes.Index(code2[i:], flag)
			if x == -1 {
				break
			}
			idx := (i + x) / 8
			code.data[idx] = makeop(OP_JMP, 0, uint32(int32(len(code.data)-idx-1)))
			code.data[idx+1] = makeop(OP_NOP, 0, 0)
			i = idx*8 + 2
		}
	}
	table.clearAllRegRecords()
	return buf, regA, nil
}
