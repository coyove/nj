package compiler

import (
	"encoding/binary"
	"fmt"

	"github.com/coyove/bracket/parser"
	"github.com/coyove/bracket/vm"

	"github.com/coyove/bracket/base"
)

func flatWrite(
	stackPtr int16,
	atoms []*parser.Node,
	varLookup *base.CMap,
	bop byte,
) (code []byte, yx int32, newStackPtr int16, err error) {

	replacedAtoms := []*parser.Node{}
	buf := base.NewBytesBuffer()

	for i := 1; i < len(atoms); i++ {
		atom := atoms[i]

		if atom.Type == parser.NTCompound {
			code, yx, stackPtr, err = compileCompoundIntoVariable(stackPtr, atom, varLookup, true, 0)
			if err != nil {
				return
			}
			atoms[i] = &parser.Node{Type: parser.NTAddr, Value: yx}
			replacedAtoms = append(replacedAtoms, atoms[i])
			buf.Write(code)
		}
	}

	if len(replacedAtoms) == 1 {
		buf.SetCursor(buf.Len() - 4)
		replacedAtoms[0].Value = buf.ReadInt32()
		buf.Truncate(9)
		stackPtr--
	}

	extflag := false
	if len(atoms) == 3 {
		switch bop {
		case base.OP_SUB, base.OP_MUL, base.OP_DIV, base.OP_MOD,
			base.OP_LESS, base.OP_LESS_EQ, base.OP_MORE, base.OP_MORE_EQ,
			base.OP_BIT_LSH, base.OP_BIT_RSH, base.OP_BIT_AND, base.OP_BIT_OR:
			extflag = true
		}
	}

	if extflag {
		l, r := atoms[1], atoms[2]
		lf := l.Type == parser.NTAtom || l.Type == parser.NTCompound || l.Type == parser.NTAddr
		rf := r.Type == parser.NTAtom || r.Type == parser.NTCompound || r.Type == parser.NTAddr

		extractWrite := func(atom *parser.Node) (int32, error) {
			code, yx, stackPtr, err = extract(stackPtr, atom, varLookup)
			buf.Write(code)
			return yx, err
		}

		var la, ra int32
		if lf && rf {
			la, err = extractWrite(l)
			if err != nil {
				return
			}

			ra, err = extractWrite(r)
			if err != nil {
				return
			}

			buf.WriteByte(base.OP_EXT_F_F)
			buf.WriteByte(bop)
			buf.WriteInt32(la)
			buf.WriteInt32(ra)
		} else if lf && r.Type == parser.NTNumber {
			la, err = extractWrite(l)
			if err != nil {
				return
			}

			buf.WriteByte(base.OP_EXT_F_IMM)
			buf.WriteByte(bop)
			buf.WriteInt32(la)
			buf.WriteDouble(r.Value.(float64))
		} else if l.Type == parser.NTNumber && rf {
			ra, err = extractWrite(r)
			if err != nil {
				return
			}

			buf.WriteByte(base.OP_EXT_IMM_F)
			buf.WriteByte(bop)
			buf.WriteDouble(l.Value.(float64))
			buf.WriteInt32(ra)
		} else {
			extflag = false
		}
	}

	if !extflag {
		if len(atoms) > 9 {
			if bop == base.OP_LIB_CALL {
				for i := 1; i < len(atoms); i++ {
					err = fill1(buf, atoms[i], varLookup, base.OP_PUSH, base.OP_PUSH_NUM, base.OP_PUSH_STR)
					if err != nil {
						return
					}
				}
				buf.WriteByte(base.OP_LIB_CALL_EX)
			} else {
				panic("shouldn't happen")
			}
		} else {
			for i := 1; i < len(atoms); i++ {
				err = fill1(buf, atoms[i], varLookup, base.OP_PUSHF, base.OP_PUSHF_NUM, base.OP_PUSHF_STR)
				if err != nil {
					return
				}
			}
			buf.WriteByte(bop)
		}
	}

	if bop == base.OP_ASSERT {
		buf.WriteString(atoms[0].String())
	}

	return buf.Bytes(), base.REG_A, stackPtr, nil
}

func flatCompile(stackPtr int16,
	atoms []*parser.Node,
	varLookup *base.CMap,
	bop byte, bop2 uint32,
	expectedArgsCount int,
) (code []byte, yx int32, newStackPtr int16, err error) {
	argsCount := len(atoms) - 1
	if argsCount < expectedArgsCount {
		//          0       1
		// (lambda (x1) (ret (op x0)))
		missingArgsCount := expectedArgsCount - argsCount
		bodyStackPtr := int16(missingArgsCount)
		for i := 0; i < missingArgsCount; i++ {
			atoms = append(atoms, &parser.Node{Type: parser.NTAddr, Value: int32(i)})
		}

		m := base.NewCMap()
		m.Parent = varLookup
		code, yx, _, err = flatWrite(bodyStackPtr, atoms, m, bop)
		if err != nil {
			return
		}

		body, buf := base.NewBytesBuffer(), base.NewBytesBuffer()
		body.Write(code)

		if bop == base.OP_LIB_CALL {
			body.WriteInt32(int32(bop2))
		}

		body.WriteByte(base.OP_RET)
		body.WriteInt32(base.REG_A)
		body.WriteByte(base.OP_EOB)

		buf.WriteByte(base.OP_LAMBDA)
		buf.WriteInt32(int32(missingArgsCount))
		buf.WriteInt32(int32(body.Len()))
		buf.Write(body.Bytes())

		return buf.Bytes(), base.REG_A, stackPtr, nil
	}

	code, yx, stackPtr, err = flatWrite(stackPtr, atoms, varLookup, bop)
	if err == nil && bop == base.OP_LIB_CALL {
		code = append(code, 0, 0, 0, 0)
		binary.LittleEndian.PutUint32(code[len(code)-4:], bop2)
	}
	return code, yx, stackPtr, err
}

// @Override
func compileFlatOp(stackPtr int16, atoms []*parser.Node, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {

	head := atoms[0]
	switch head.Value.(string) {
	case "+":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_ADD, 0, 2)
	case "-":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_SUB, 0, 1)
	case "*":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_MUL, 0, 2)
	case "/":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_DIV, 0, 2)
	case "inc":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_INC, 0, 2)
	case "%":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_MOD, 0, 2)
	case "<":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_LESS, 0, 2)
	case "<=":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_LESS_EQ, 0, 2)
	case ">":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_MORE, 0, 2)
	case ">=":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_MORE_EQ, 0, 2)
	case "eq", "==":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_EQ, 0, 2)
	case "neq":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_NEQ, 0, 2)
	case "assert":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_ASSERT, 0, 1)
	case "list":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_LIST, 0, 0)
	case "map":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_MAP, 0, 0)
	case "expand":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_EXPAND, 0, 1)
	case "not":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_NOT, 0, 1)
	case "and":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_AND, 0, 2)
	case "or":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_OR, 0, 2)
	case "xor":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_XOR, 0, 2)
	case "b/not":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_BIT_NOT, 0, 1)
	case "b/and":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_BIT_AND, 0, 2)
	case "b/or":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_BIT_OR, 0, 2)
	case "b/xor":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_BIT_XOR, 0, 2)
	case "b/lsh":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_BIT_LSH, 0, 2)
	case "b/rsh":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_BIT_RSH, 0, 2)
	case "store":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_STORE, 0, 3)
	case "load":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_LOAD, 0, 2)
	case "len":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_LEN, 0, 1)
	case "nil":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_NIL, 0, 0)
	case "true":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_TRUE, 0, 0)
	case "false":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_FALSE, 0, 0)
	case "bytes":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_BYTES, 0, 0)
	}

	if lib, ok := vm.LibLookup[head.Value.(string)]; ok {
		return flatCompile(stackPtr, atoms, varLookup, base.OP_LIB_CALL, uint32(lib), vm.Lib[lib].Args())
	}

	err = fmt.Errorf("invalid flat op %+v", head)
	return
}
