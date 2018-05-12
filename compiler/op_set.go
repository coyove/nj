package compiler

import (
	"fmt"

	"github.com/coyove/bracket/base"
	"github.com/coyove/bracket/parser"
)

const (
	c_set = iota + 1
	c_declare
)

const (
	ERR_UNDECLARED_VARIABLE = "undeclared variable: %+v"
)

func compileSetOp(stackPtr int16, atoms []*parser.Node, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	aVar := atoms[1]
	varIndex := int32(0)
	if len(atoms) < 3 {
		err = fmt.Errorf("can't set/declare without value %+v", atoms[0])
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

	return buf.Bytes(), newYX, stackPtr, nil
}

func compileRetOp(r, n, s byte) func(stackPtr int16, atoms []*parser.Node, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	return func(stackPtr int16, atoms []*parser.Node, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
		buf := base.NewBytesBuffer()
		if len(atoms) == 1 {
			buf.WriteByte(r)
			buf.WriteInt32(base.REG_A)
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
		return buf.Bytes(), yx, stackPtr, nil
	}
}
