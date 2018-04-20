package compiler

import (
	"fmt"

	"github.com/coyove/bracket/base"
	"github.com/coyove/bracket/parser"
)

func inc(
	stackPtr int16,
	src int32,
	step *parser.Node,
	varLookup *base.CMap,
) (code []byte, yx int32, newStackPtr int16, err error) {
	buf := base.NewBytesBuffer()

	switch step.Type {
	case parser.NTNumber:
		buf.WriteByte(base.OP_INC_NUM)
		buf.WriteInt32(src)
		buf.WriteDouble(step.Value.(float64))
	case parser.NTString:
		err = fmt.Errorf("can't inc by a string step: %+v", step)
		return
	case parser.NTCompound, parser.NTAtom:
		code, yx, stackPtr, err = extract(stackPtr, step, varLookup)
		if err != nil {
			return
		}
		buf.Write(code)
		buf.WriteByte(base.OP_INC)
		buf.WriteInt32(src)
		buf.WriteInt32(yx)
		break
	case parser.NTAddr:
		buf.WriteByte(base.OP_INC)
		buf.WriteInt32(src)
		buf.WriteInt32(step.Value.(int32))
	}
	return buf.Bytes(), src, stackPtr, nil
}

func compileIncOp(stackPtr int16, atoms []*parser.Node, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	if len(atoms) < 3 {
		err = fmt.Errorf("inc must have a src and a step: %+v", atoms[0])
		return
	}

	var src int32
	switch aSrc := atoms[1]; aSrc.Type {
	case parser.NTNumber, parser.NTString:
		err = fmt.Errorf("can't inc an immediate value: %+v", atoms[0])
		return
	case parser.NTCompound:
		err = fmt.Errorf("can't inc a compound: %+v", atoms[0])
		return
	case parser.NTAtom:
		src = varLookup.GetRelPosition(aSrc.Value.(string))
		if src == -1 {
			err = fmt.Errorf(ERR_UNDECLARED_VARIABLE, aSrc)
			return
		}
	case parser.NTAddr:
		src = aSrc.Value.(int32)
	}

	return inc(stackPtr, src, atoms[2], varLookup)
}
