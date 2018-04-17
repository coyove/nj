package compiler

import (
	"fmt"

	"github.com/coyove/eugine/base"
)

func compileCallOp(stackPtr int16, atoms []*token, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	buf := base.NewBytesBuffer()

	for i := 1; i < len(atoms); i++ {
		atom := atoms[i]

		if atom.ty == TK_compound {
			code, yx, stackPtr, err = compileCompoundIntoVariable(stackPtr, atom, varLookup, true, 0)
			if err != nil {
				return
			}
			atoms[i] = &token{ty: TK_addr, v: yx}
			buf.Write(code)
		}
	}

	for i := 1; i < len(atoms); i++ {
		err = fill1(buf, atoms[i], varLookup, base.OP_PUSH, base.OP_PUSH_NUM, base.OP_PUSH_STR)
		if err != nil {
			return
		}
	}

	callee := atoms[0]
	if callee.ty == TK_atomic {
		varIndex := varLookup.GetRelPosition(callee.v.(string))
		if varIndex == -1 {
			err = fmt.Errorf(ERR_UNDECLARED_VARIABLE, callee)
			return
		}

		buf.WriteByte(base.OP_CALL)
		buf.WriteInt32(varIndex)
	} else if callee.ty == TK_compound {
		code, yx, stackPtr, err = compileCompoundIntoVariable(stackPtr, callee, varLookup, true, 0)
		if err != nil {
			return
		}
		buf.Write(code)
		buf.WriteByte(base.OP_CALL)
		buf.WriteInt32(yx)
	} else if callee.ty == TK_addr {
		buf.WriteByte(base.OP_CALL)
		buf.WriteInt32(callee.v.(int32))
	} else {
		err = fmt.Errorf("invalid callee: %+v", callee)
		return
	}

	return buf.Bytes(), base.REG_A, stackPtr, nil
}

func compileLambdaOp(stackPtr int16, atoms []*token, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	newLookup := base.NewCMap()
	newLookup.Parent = varLookup

	params := atoms[1]
	if params.ty != TK_compound {
		err = fmt.Errorf("invalid lambda parameters: %+v", atoms[0])
		return
	}

	// defer func() {
	// 	if recover() != nil {
	// 		fmt.Println(atoms[0])
	// 	}
	// }()
	for i, p := range params.v.([]*token) {
		newLookup.M[p.v.(string)] = int16(i)
	}

	ln := len(newLookup.M)
	code, yx, _, err = compile(int16(ln), atoms[2:], newLookup)
	if err != nil {
		return
	}

	buf := base.NewBytesBuffer()
	buf.WriteByte(base.OP_LAMBDA)
	buf.WriteInt32(int32(ln))
	buf.WriteInt32(int32(len(code)))
	buf.Write(code)

	return buf.Bytes(), base.REG_A, stackPtr, nil
}
