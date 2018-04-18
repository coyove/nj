package compiler

import (
	"fmt"

	"github.com/coyove/bracket/base"
)

func compileTypeofOp(stackPtr int16, atoms []*token, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	if len(atoms) < 3 {
		err = fmt.Errorf("typeof needs 2 arguments: %+v", atoms[0])
		return
	}

	atom := atoms[1]
	buf := base.NewBytesBuffer()
	switch atom.ty {
	case TK_string, TK_number:
		err = fmt.Errorf("no need to assert the type of an immediate value: %+v", atom)
		return
	case TK_addr:
		buf.WriteByte(base.OP_TYPEOF)
		buf.WriteInt32(atom.v.(int32))
		break
	case TK_compound, TK_atomic:
		code, yx, stackPtr, err = extract(stackPtr, atom, varLookup)
		if err != nil {
			return nil, 0, 0, err
		}
		buf.Write(code)
		buf.WriteByte(base.OP_TYPEOF)
		buf.WriteInt32(yx)
	}

	t := atoms[2]
	if t.ty != TK_atomic {
		err = fmt.Errorf("typeof needs an atom to test: %+v", t)
		return
	}

	switch t.v.(string) {
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
