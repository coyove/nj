package compiler

import (
	"fmt"

	"github.com/coyove/bracket/base"
	"github.com/coyove/bracket/parser"
)

const (
	ERR_UNDECLARED_VARIABLE = "undeclared variable: %+v"
)

func compileSetOp(stackPtr int16, atoms []*parser.Node, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	aVar := atoms[1]
	varIndex := int32(0)
	if len(atoms) < 3 {
		err = fmt.Errorf("can't set/move without value %+v", atoms[0])
		return
	}

	aValue := atoms[2]

	buf := base.NewBytesBuffer()
	var newYX int32
	if atoms[0].Value.(string) == "set" {
		// compound has its own logic, we won't incr stack here
		if aValue.Type != parser.NTCompound {
			newYX = int32(stackPtr)
			stackPtr++
		}
	} else {
		varIndex = varLookup.GetRelPosition(aVar.Value.(string))
		if varIndex == -1 {
			err = fmt.Errorf(ERR_UNDECLARED_VARIABLE, aVar)
			return
		}
		newYX = varIndex
	}

	switch aValue.Type {
	case parser.NTAtom:
		valueIndex := varLookup.GetRelPosition(aValue.Value.(string))
		if valueIndex == -1 {
			err = fmt.Errorf(ERR_UNDECLARED_VARIABLE, aValue)
			return
		}

		buf.WriteByte(base.OP_SET)
		buf.WriteInt32(newYX)
		buf.WriteInt32(valueIndex)
	case parser.NTNumber:
		buf.WriteByte(base.OP_SET_NUM)
		buf.WriteInt32(newYX)
		buf.WriteDouble(aValue.Value.(float64))
	case parser.NTString:
		buf.WriteByte(base.OP_SET_STR)
		buf.WriteInt32(newYX)
		buf.WriteString(aValue.Value.(string))
	case parser.NTCompound:
		code, newYX, stackPtr, err = compileCompoundIntoVariable(stackPtr, aValue, varLookup,
			atoms[0].Value.(string) == "set", varIndex)
		if err != nil {
			return
		}
		buf.Write(code)
	}

	if atoms[0].Value.(string) == "set" {
		_, reset := varLookup.M[aVar.Value.(string)]
		if reset {
			err = fmt.Errorf("redeclare: %+v", aVar)
			return
		}

		varLookup.M[aVar.Value.(string)] = int16(newYX)
	}

	varLookup.I = nil
	return buf.Bytes(), newYX, stackPtr, nil
}

func compileRetOp(r, n, s byte) func(stackPtr int16, atoms []*parser.Node, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	return func(stackPtr int16, atoms []*parser.Node, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
		buf := base.NewBytesBuffer()
		if len(atoms) == 1 {
			buf.WriteByte(r)
			buf.WriteInt32(base.REG_A)
			varLookup.I = nil
			return buf.Bytes(), yx, stackPtr, nil
		}

		atom := atoms[1]

		switch atom.Type {
		case parser.NTAtom, parser.NTNumber, parser.NTString, parser.NTAddr:
			err = fill1(buf, atom, varLookup, r, n, s)
			if err != nil {
				return
			}
		case parser.NTCompound:
			code, yx, stackPtr, err = extract(stackPtr, atom, varLookup)
			buf.Write(code)
			buf.WriteByte(r)
			buf.WriteInt32(yx)
		}
		varLookup.I = nil
		return buf.Bytes(), yx, stackPtr, nil
	}
}

func compileListOp(stackPtr int16, atoms []*parser.Node, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	var buf *base.BytesReader
	buf, stackPtr, err = flaten(stackPtr, atoms[1].Compound, varLookup)
	if err != nil {
		return
	}

	for _, atom := range atoms[1].Compound {
		err = fill1(buf, atom, varLookup, base.OP_PUSH, base.OP_PUSH_NUM, base.OP_PUSH_STR)
		if err != nil {
			return
		}
	}

	buf.WriteByte(base.OP_LIST)
	varLookup.I = nil
	return buf.Bytes(), base.REG_A, stackPtr, nil
}

func compileMapOp(stackPtr int16, atoms []*parser.Node, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	if len(atoms[1].Compound)%2 != 0 {
		err = fmt.Errorf("every key in map must have a value: %+v", atoms[1])
		return
	}

	var buf *base.BytesReader
	buf, stackPtr, err = flaten(stackPtr, atoms[1].Compound, varLookup)
	if err != nil {
		return
	}

	for _, atom := range atoms[1].Compound {
		err = fill1(buf, atom, varLookup, base.OP_PUSH, base.OP_PUSH_NUM, base.OP_PUSH_STR)
		if err != nil {
			return
		}
	}

	buf.WriteByte(base.OP_MAP)
	varLookup.I = nil
	return buf.Bytes(), base.REG_A, stackPtr, nil
}
