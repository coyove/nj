package compiler

import (
	"fmt"

	"github.com/coyove/bracket/base"
	"github.com/coyove/bracket/parser"
)

func compileListOp(stackPtr int16, atoms []*parser.Node, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	var buf *base.BytesReader
	buf, stackPtr, err = flaten(stackPtr, atoms[2:], varLookup)
	if err != nil {
		return
	}

	for i := 2; i < len(atoms); i++ {
		err = fill1(buf, atoms[i], varLookup, base.OP_PUSH, base.OP_PUSH_NUM, base.OP_PUSH_STR)
		if err != nil {
			return
		}
	}

	buf.WriteByte(base.OP_LIST)
	return buf.Bytes(), base.REG_A, stackPtr, nil
}

func compileMapOp(stackPtr int16, atoms []*parser.Node, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	if len(atoms)%2 != 0 {
		err = fmt.Errorf("every key in map must have a value: %+v", atoms[1])
	}

	var buf *base.BytesReader
	buf, stackPtr, err = flaten(stackPtr, atoms[2:], varLookup)
	if err != nil {
		return
	}

	for i := 2; i < len(atoms); i++ {
		err = fill1(buf, atoms[i], varLookup, base.OP_PUSH, base.OP_PUSH_NUM, base.OP_PUSH_STR)
		if err != nil {
			return
		}
	}

	buf.WriteByte(base.OP_MAP)
	return buf.Bytes(), base.REG_A, stackPtr, nil
}
