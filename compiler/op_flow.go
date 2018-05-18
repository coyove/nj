package compiler

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/coyove/bracket/base"
	"github.com/coyove/bracket/parser"
)

// [if condition [truechain ...] [falsechain ...]]
func compileIfOp(stackPtr int16, atoms []*parser.Node, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	if len(atoms) < 3 {
		err = fmt.Errorf("if statement should have at least a true branch: %+v", atoms[0])
		return
	}

	condition := atoms[1]
	trueBranch, falseBranch := atoms[2], atoms[3]

	switch condition.Type {
	case parser.NTNumber, parser.NTString:
		err = fmt.Errorf("can't use immediate value as if condition: %+v", atoms[0])
		return
	case parser.NTAtom, parser.NTCompound:
		buf := base.NewBytesWriter()
		code, yx, stackPtr, err = extract(stackPtr, condition, varLookup)
		if err != nil {
			return nil, 0, 0, err
		}

		buf.Write(code)
		buf.WriteByte(base.OP_IF)
		buf.WriteInt32(yx)

		var trueCode, falseCode []byte
		trueCode, yx, stackPtr, err = compileChainOp(stackPtr, trueBranch, varLookup)
		if err != nil {
			return
		}
		falseCode, yx, stackPtr, err = compileChainOp(stackPtr, falseBranch, varLookup)
		if err != nil {
			return
		}
		if len(falseCode) > 0 {
			buf.WriteInt32(int32(len(trueCode)) + 5) // jmp (1b) + offset (4b)
			buf.Write(trueCode)
			buf.WriteByte(base.OP_JMP)
			buf.WriteInt32(int32(len(falseCode)))
			buf.Write(falseCode)
		} else {
			buf.WriteInt32(int32(len(trueCode)))
			buf.Write(trueCode)
		}
		varLookup.I = nil
		return buf.Bytes(), base.REG_A, stackPtr, nil
	}

	err = fmt.Errorf("not a valid condition: %+v", condition)
	return
}

// [call func-name [args ...]]
func compileCallOp(stackPtr int16, nodes []*parser.Node, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	buf := base.NewBytesWriter()
	callee := nodes[1]

	name, _ := callee.Value.(string)
	switch name {
	case "len", "dup", "error":
		atoms := append(nodes[1:2], nodes[2].Compound...)
		if name == "error" {
			return flatWrite(stackPtr, atoms, varLookup, base.OP_ERROR)
		}
		if len(atoms) < 2 {
			err = fmt.Errorf("missing subject to call %s: %v", name, callee)
			return
		}
		switch name {
		case "len":
			return flatWrite(stackPtr, atoms, varLookup, base.OP_LEN)
		case "dup":
			return flatWrite(stackPtr, atoms, varLookup, base.OP_DUP)
		}
	case "who":
		return []byte{base.OP_WHO}, base.REG_A, stackPtr, nil
	case "stack":
		return []byte{base.OP_STACK}, base.REG_A, stackPtr, nil
	}

	atoms := nodes[2].Compound

	for i := 0; i < len(atoms); i++ {
		atom := atoms[i]

		if atom.Type == parser.NTCompound {
			code, yx, stackPtr, err = compileCompoundIntoVariable(stackPtr, atom, varLookup, true, 0)
			if err != nil {
				return
			}
			atoms[i] = &parser.Node{Type: parser.NTAddr, Value: yx}
			buf.Write(code)
		}
	}

	var varIndex int32
	switch callee.Type {
	case parser.NTAtom:
		varIndex = varLookup.GetRelPosition(callee.Value.(string))
		if varIndex == -1 {
			err = fmt.Errorf(ERR_UNDECLARED_VARIABLE, callee)
			return
		}
	case parser.NTCompound:
		code, yx, stackPtr, err = compileCompoundIntoVariable(stackPtr, callee, varLookup, true, 0)
		if err != nil {
			return
		}

		varIndex = yx
		buf.Write(code)
	case parser.NTAddr:
		varIndex = callee.Value.(int32)
	case parser.NTString:
		buf.WriteByte(base.OP_SET_STR)
		buf.WriteInt32(int32(stackPtr))
		buf.WriteString(callee.Value.(string))
		varIndex = int32(stackPtr)
		stackPtr++
	default:
		err = fmt.Errorf("invalid callee: %+v", callee)
		return
	}

	for i := 0; i < len(atoms); i++ {
		err = fill1(buf, atoms[i], varLookup, base.OP_PUSH, base.OP_PUSH_NUM, base.OP_PUSH_STR)
		if err != nil {
			return
		}
	}

	buf.WriteByte(base.OP_CALL)
	buf.WriteInt32(varIndex)

	varLookup.I = nil
	return buf.Bytes(), base.REG_A, stackPtr, nil
}

// [lambda [namelist] [chain ...]]
func compileLambdaOp(stackPtr int16, atoms []*parser.Node, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	newLookup := base.NewCMap()
	newLookup.Parent = varLookup

	params := atoms[1]
	if params.Type != parser.NTCompound {
		err = fmt.Errorf("invalid lambda parameters: %+v", atoms[0])
		return
	}

	for i, p := range params.Compound {
		newLookup.M[p.Value.(string)] = int16(i)
	}

	ln := len(newLookup.M)
	code, yx, _, err = compileChainOp(int16(ln), atoms[2], newLookup)
	if err != nil {
		return
	}

	code = append(code, base.OP_EOB)
	buf := base.NewBytesWriter()
	buf.WriteByte(base.OP_LAMBDA)
	buf.WriteInt32(int32(ln))
	if newLookup.Y {
		buf.WriteByte(1)
	} else {
		buf.WriteByte(0)
	}
	buf.WriteInt32(int32(len(code)))
	buf.Write(code)

	return buf.Bytes(), base.REG_A, stackPtr, nil
}

var staticWhileHack struct {
	sync.Mutex
	continueFlag []byte
	breakFlag    []byte
}

// [continue | break]
func compileContinueBreakOp(stackPtr int16, atoms []*parser.Node, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	staticWhileHack.Lock()
	defer staticWhileHack.Unlock()
	if atoms[0].Value.(string) == "continue" {
		if staticWhileHack.continueFlag == nil {
			staticWhileHack.continueFlag = make([]byte, 9)
			if _, err := rand.Read(staticWhileHack.continueFlag); err != nil {
				panic(err)
			}
		}
		return staticWhileHack.continueFlag, base.REG_A, stackPtr, nil
	}

	if staticWhileHack.breakFlag == nil {
		staticWhileHack.breakFlag = make([]byte, 9)
		if _, err := rand.Read(staticWhileHack.breakFlag); err != nil {
			panic(err)
		}
	}
	return staticWhileHack.breakFlag, base.REG_A, stackPtr, nil
}

// [while condition [chain ...]]
func compileWhileOp(stackPtr int16, atoms []*parser.Node, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	if len(atoms) < 3 {
		err = fmt.Errorf("while statement should have condition and body: %+v", atoms[0])
		return
	}
	condition := atoms[1]
	buf := base.NewBytesWriter()
	var varIndex int32

	switch condition.Type {
	case parser.NTAddr:
		varIndex = condition.Value.(int32)
	case parser.NTNumber:
		buf.WriteByte(base.OP_SET_NUM)
		varIndex = int32(stackPtr)
		stackPtr++
		buf.WriteInt32(varIndex)
		buf.WriteDouble(condition.Value.(float64))
	case parser.NTString:
		buf.WriteByte(base.OP_SET_STR)
		varIndex = int32(stackPtr)
		stackPtr++
		buf.WriteInt32(varIndex)
		buf.WriteString(condition.Value.(string))
	case parser.NTCompound, parser.NTAtom:
		code, yx, stackPtr, err = extract(stackPtr, condition, varLookup)
		if err != nil {
			return
		}
		buf.Write(code)
		varIndex = yx
	}

	code, yx, stackPtr, err = compileChainOp(stackPtr, atoms[2], varLookup)
	if err != nil {
		return
	}

	buf.WriteByte(base.OP_IF)
	buf.WriteInt32(varIndex)
	buf.WriteInt32(int32(len(code)) + 5)
	buf.Write(code)
	buf.WriteByte(base.OP_JMP)
	buf.WriteInt32(-int32(buf.Len()) - 4)

	code = buf.Bytes()
	i := 0
	for i < len(code) && staticWhileHack.continueFlag != nil {
		x := bytes.Index(code[i:], staticWhileHack.continueFlag)
		if x == -1 {
			break
		}
		idx := i + x
		code[idx] = base.OP_JMP
		binary.LittleEndian.PutUint32(code[idx+1:], uint32(-(idx + 5)))
		copy(code[idx+5:], []byte{base.OP_NOP, base.OP_NOP, base.OP_NOP, base.OP_NOP})
		i = idx + 9
		// 				buf.WriteInt(i + 1, bop == Op.JMP_BREAK ?
		// 						(f.code.size() - i) - 5 :
		// 						-(i + 5));
	}

	i = 0
	for i < len(code) && staticWhileHack.breakFlag != nil {
		x := bytes.Index(code[i:], staticWhileHack.breakFlag)
		if x == -1 {
			break
		}
		idx := i + x
		code[idx] = base.OP_JMP
		binary.LittleEndian.PutUint32(code[idx+1:], uint32((len(code)-idx)-5))
		copy(code[idx+5:], []byte{base.OP_NOP, base.OP_NOP, base.OP_NOP, base.OP_NOP})
		i = idx + 9
	}

	return code, base.REG_A, stackPtr, nil
}
