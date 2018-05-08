package compiler

import (
	"encoding/binary"
	"fmt"

	"github.com/coyove/bracket/parser"
	"github.com/coyove/bracket/vm"

	"github.com/coyove/bracket/base"
)

func indexToOpR(index int) []byte {
	switch index {
	case 0:
		return []byte{base.OP_R0, base.OP_R0_NUM, base.OP_R0_STR}
	case 1:
		return []byte{base.OP_R1, base.OP_R1_NUM, base.OP_R1_STR}
	case 2:
		return []byte{base.OP_R2, base.OP_R2_NUM, base.OP_R2_STR}
	case 3:
		return []byte{base.OP_R3, base.OP_R3_NUM, base.OP_R3_STR}
	}
	panic("shouldn't happen")
}

func flaten(stackPtr int16, atoms []*parser.Node, varLookup *base.CMap) (buf *base.BytesReader, newStackPtr int16, err error) {

	replacedAtoms := []*parser.Node{}
	buf = base.NewBytesBuffer()

	for i, atom := range atoms {

		var yx int32
		var code []byte

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

	return buf, stackPtr, nil
}

func flatWrite(
	stackPtr int16,
	atoms []*parser.Node,
	varLookup *base.CMap,
	bop byte, bop2 uint32,
) (code []byte, yx int32, newStackPtr int16, err error) {

	var buf *base.BytesReader
	buf, stackPtr, err = flaten(stackPtr, atoms[1:], varLookup)
	if err != nil {
		return
	}

	if bop == base.OP_LIB_CALL {
		ff := vm.Lib[bop2].IsFF()
		for i := 1; i < len(atoms); i++ {
			if ff {
				err = fill1(buf, atoms[i], varLookup, base.OP_PUSH, base.OP_PUSH_NUM, base.OP_PUSH_STR)
			} else {
				err = fill1(buf, atoms[i], varLookup, indexToOpR(i-1)...)
			}
			if err != nil {
				return
			}
		}

		if ff {
			buf.WriteByte(base.OP_LIB_CALL_EX)
		} else {
			buf.WriteByte(base.OP_LIB_CALL)
		}
	} else {
		for i := 1; i < len(atoms); i++ {
			err = fill1(buf, atoms[i], varLookup, indexToOpR(i-1)...)
			if err != nil {
				return
			}
		}
		buf.WriteByte(bop)
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
		code, yx, _, err = flatWrite(bodyStackPtr, atoms, m, bop, bop2)
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

	code, yx, stackPtr, err = flatWrite(stackPtr, atoms, varLookup, bop, bop2)
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
	case "eq":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_EQ, 0, 2)
	case "neq":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_NEQ, 0, 2)
	case "assert":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_ASSERT, 0, 1)
	case "not":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_NOT, 0, 1)
	case "and":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_AND, 0, 2)
	case "or":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_OR, 0, 2)
	case "xor":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_XOR, 0, 2)
	case "~":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_BIT_NOT, 0, 1)
	case "&":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_BIT_AND, 0, 2)
	case "|":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_BIT_OR, 0, 2)
	case "^":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_BIT_XOR, 0, 2)
	case "<<":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_BIT_LSH, 0, 2)
	case ">>":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_BIT_RSH, 0, 2)
	case "len":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_LEN, 0, 1)
	case "nil":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_NIL, 0, 0)
	case "true":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_TRUE, 0, 0)
	case "false":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_FALSE, 0, 0)
	case "bytes":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_BYTES, 0, 1)
	case "store":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_STORE, 0, 3)
	case "load":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_LOAD, 0, 2)
	case "safestore":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_SAFE_STORE, 0, 3)
	case "safeload":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_SAFE_LOAD, 0, 2)
	case "dup":
		return flatCompile(stackPtr, atoms, varLookup, base.OP_DUP, 0, 1)
	}

	if lib, ok := vm.LibLookup[head.Value.(string)]; ok {
		return flatCompile(stackPtr, atoms, varLookup, base.OP_LIB_CALL, uint32(lib), vm.Lib[lib].Args())
	}

	err = fmt.Errorf("invalid flat op %+v", head)
	return
}
