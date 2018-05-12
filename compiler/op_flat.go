package compiler

import (
	"encoding/binary"
	"fmt"

	"github.com/coyove/bracket/parser"

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
		cursor := buf.Len() - 4
		replacedAtoms[0].Value = int32(binary.LittleEndian.Uint32(buf.Bytes()[:cursor]))
		buf.Truncate(9)
		stackPtr--
	}

	return buf, stackPtr, nil
}

func flatWrite(stackPtr int16, atoms []*parser.Node, varLookup *base.CMap, bop byte) (code []byte, yx int32, newStackPtr int16, err error) {

	var buf *base.BytesReader
	buf, stackPtr, err = flaten(stackPtr, atoms[1:], varLookup)
	if err != nil {
		return
	}

	for i := 1; i < len(atoms); i++ {
		err = fill1(buf, atoms[i], varLookup, indexToOpR(i-1)...)
		if err != nil {
			return
		}
	}
	buf.WriteByte(bop)

	if bop == base.OP_ASSERT {
		buf.WriteString(atoms[0].String())
	}

	return buf.Bytes(), base.REG_A, stackPtr, nil
}

func compileFlatOp(stackPtr int16, atoms []*parser.Node, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {

	head := atoms[0]
	switch head.Value.(string) {
	case "+":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_ADD)
	case "-":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_SUB)
	case "*":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_MUL)
	case "/":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_DIV)
	case "%":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_MOD)
	case "<":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_LESS)
	case "<=":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_LESS_EQ)
	case ">":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_MORE)
	case ">=":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_MORE_EQ)
	case "eq":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_EQ)
	case "neq":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_NEQ)
	case "assert":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_ASSERT)
	case "not":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_NOT)
	case "and":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_AND)
	case "or":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_OR)
	case "xor":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_XOR)
	case "~":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_BIT_NOT)
	case "&":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_BIT_AND)
	case "|":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_BIT_OR)
	case "^":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_BIT_XOR)
	case "<<":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_BIT_LSH)
	case ">>":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_BIT_RSH)
	case "nil":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_NIL)
	case "true":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_TRUE)
	case "false":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_FALSE)
	case "store":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_STORE)
	case "load":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_LOAD)
	case "safestore":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_SAFE_STORE)
	case "safeload":
		return flatWrite(stackPtr, atoms, varLookup, base.OP_SAFE_LOAD)
	}

	err = fmt.Errorf("invalid flat op %+v", head)
	return
}
