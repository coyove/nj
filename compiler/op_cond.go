package compiler

import (
	"fmt"

	"github.com/coyove/bracket/base"
)

func truncLastByte(b []byte) []byte {
	if len(b) > 0 {
		return b[:len(b)-1]
	}
	return b
}

func compileIfOp(stackPtr int16, atoms []*token, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	/* [if [condition]
			[true branch code]
			[true branch code]
			...
		else
			[false branch code]
			[false branch code]
			...
	   ] */
	if len(atoms) < 3 {
		err = fmt.Errorf("if statement should have at least a true branch: %+v", atoms[0])
		return
	}

	condition := atoms[1]
	trueBranch, falseBranch := make([]*token, 0, 16), make([]*token, 0, 16)
	currentBranch := &trueBranch

	for i := 2; i < len(atoms); i++ {
		l := atoms[i]
		if equalI("else", l.v) {
			if currentBranch == &falseBranch {
				err = fmt.Errorf("if already has an else: %+v", atoms[0])
				return
			}
			currentBranch = &falseBranch
		} else {
			*currentBranch = append(*currentBranch, l)
		}
	}

	switch condition.ty {
	case TK_number, TK_string:
		fflag := false
		if condition.ty == TK_number && condition.v.(float64) == 0 {
			fflag = true
		} else if condition.ty == TK_string && condition.v.(string) == "" {
			fflag = true
		}

		if fflag {
			if len(falseBranch) > 0 {
				code, yx, stackPtr, err = compile(stackPtr, falseBranch, varLookup)
			}
		} else {
			code, yx, stackPtr, err = compile(stackPtr, trueBranch, varLookup)
		}

		code = truncLastByte(code)
		yx = base.REG_A
		return
	case TK_atomic, TK_compound:
		buf := base.NewBytesBuffer()
		code, yx, stackPtr, err = extract(stackPtr, condition, varLookup)
		if err != nil {
			return nil, 0, 0, err
		}

		buf.Write(code)
		buf.WriteByte(base.OP_IF)
		buf.WriteInt32(yx)

		var trueCode, falseCode []byte
		trueCode, yx, stackPtr, err = compile(stackPtr, trueBranch, varLookup)
		trueCode = truncLastByte(trueCode)

		falseCode, yx, stackPtr, err = compile(stackPtr, falseBranch, varLookup)
		falseCode = truncLastByte(falseCode)

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
		return buf.Bytes(), base.REG_A, stackPtr, nil
	}

	err = fmt.Errorf("not a valid condition: %+v", condition)
	return
}
