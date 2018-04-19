package compiler

import (
	"fmt"

	"github.com/coyove/bracket/base"
	"github.com/coyove/bracket/parser"
)

func compileTypeofOp(stackPtr int16, atoms []*parser.Node, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	if len(atoms) < 3 {
		err = fmt.Errorf("typeof needs 2 arguments: %+v", atoms[0])
		return
	}

	atom := atoms[1]
	buf := base.NewBytesBuffer()
	switch atom.Type {
	case parser.NTString, parser.NTNumber:
		err = fmt.Errorf("no need to assert the type of an immediate value: %+v", atom)
		return
	case parser.NTAddr:
		buf.WriteByte(base.OP_TYPEOF)
		buf.WriteInt32(atom.Value.(int32))
		break
	case parser.NTCompound, parser.NTAtom:
		code, yx, stackPtr, err = extract(stackPtr, atom, varLookup)
		if err != nil {
			return nil, 0, 0, err
		}
		buf.Write(code)
		buf.WriteByte(base.OP_TYPEOF)
		buf.WriteInt32(yx)
	}

	t := atoms[2]
	if t.Type != parser.NTAtom {
		err = fmt.Errorf("typeof needs an atom to test: %+v", t)
		return
	}

	switch t.Value.(string) {
	case "number":
		buf.WriteInt32(base.TY_number)
	case "string":
		buf.WriteInt32(base.TY_string)
	case "list":
		buf.WriteInt32(base.TY_array)
	case "closure":
		buf.WriteInt32(base.TY_closure)
	case "nil":
		buf.WriteInt32(base.TY_nil)
	case "map":
		buf.WriteInt32(base.TY_map)
	case "bool":
		buf.WriteInt32(base.TY_bool)
	default:
		err = fmt.Errorf("invalid type to test: %+v", t)
		return
	}

	return buf.Bytes(), base.REG_A, stackPtr, nil
}
