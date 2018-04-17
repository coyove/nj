package compiler

import (
	"fmt"

	"github.com/coyove/eugine/base"
)

func inc(
	stackPtr int16,
	src int32,
	step *token,
	varLookup *base.CMap,
) (code []byte, yx int32, newStackPtr int16, err error) {
	buf := base.NewBytesBuffer()

	switch step.ty {
	case TK_number:
		buf.WriteByte(base.OP_INC_NUM)
		buf.WriteInt32(src)
		buf.WriteDouble(step.v.(float64))
	case TK_string:
		err = fmt.Errorf("can't inc by a string step: %+v", step)
		return
	case TK_compound, TK_atomic:
		code, yx, stackPtr, err = extract(stackPtr, step, varLookup)
		if err != nil {
			return
		}
		buf.Write(code)
		buf.WriteByte(base.OP_INC)
		buf.WriteInt32(src)
		buf.WriteInt32(yx)
		break
	case TK_addr:
		buf.WriteByte(base.OP_INC)
		buf.WriteInt32(src)
		buf.WriteInt32(step.v.(int32))
	}
	return buf.Bytes(), src, stackPtr, nil
}

func compileIncOp(stackPtr int16, atoms []*token, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	if len(atoms) < 3 {
		err = fmt.Errorf("inc must have a src and a step: %+v", atoms[0])
		return
	}

	var src int32
	switch aSrc := atoms[1]; aSrc.ty {
	case TK_number, TK_string:
		err = fmt.Errorf("can't inc an immediate value: %+v", atoms[0])
		return
	case TK_compound:
		if isStoreLoadSugar(aSrc) {
			fatoms := expandStoreLoadSugar(aSrc).v.([]*token)
			code, yx, stackPtr, err = flatWrite(stackPtr, fatoms, varLookup, base.OP_LOAD)
			if err != nil {
				return
			}

			buf := base.NewBytesBuffer()
			buf.Write(code)
			// To prevent inc $a $a
			buf.WriteByte(base.OP_SET)
			yx = int32(stackPtr)
			stackPtr++
			buf.WriteInt32(yx)
			buf.WriteInt32(base.REG_A)

			code, yx, stackPtr, err = inc(stackPtr, yx, atoms[2], varLookup)
			if err != nil {
				return
			}
			buf.Write(code)

			fatoms = append(fatoms, &token{ty: TK_addr, v: int32(base.REG_A)})
			code, yx, stackPtr, err = flatWrite(stackPtr, fatoms, varLookup, base.OP_STORE)
			if err != nil {
				return
			}
			buf.Write(code)

			return buf.Bytes(), base.REG_A, stackPtr, nil
		}
		err = fmt.Errorf("can't inc a compound: %+v", atoms[0])
		return
	case TK_atomic:
		src = varLookup.GetRelPosition(aSrc.v.(string))
		if src == -1 {
			err = fmt.Errorf(ERR_UNDECLARED_VARIABLE, aSrc)
			return
		}
	case TK_addr:
		src = aSrc.v.(int32)
	}

	return inc(stackPtr, src, atoms[2], varLookup)
}
