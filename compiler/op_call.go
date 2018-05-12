package compiler

import (
	"fmt"

	"github.com/coyove/bracket/base"
	"github.com/coyove/bracket/parser"
)

func compileCallOp(stackPtr int16, nodes []*parser.Node, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	buf := base.NewBytesBuffer()
	callee := nodes[1]

	name, _ := callee.Value.(string)
	switch name {
	case "list":
		return compileListOp(stackPtr, nodes, varLookup)
	case "map":
		return compileMapOp(stackPtr, nodes, varLookup)
	case "who":
		return []byte{base.OP_WHO}, base.REG_A, stackPtr, nil
	case "varargs":
		return []byte{base.OP_VARARGS}, base.REG_A, stackPtr, nil
	}
	if flatOpMapping[name] {
		return compileFlatOp(stackPtr, append(nodes[1:2], nodes[2].Compound...), varLookup)
	}

	atoms := nodes[2].Compound

	for i := 0; i < len(atoms); i++ {
		atom := atoms[i]

		if atom.Type == parser.NTCompound {
			code, yx, stackPtr, err = compileCompoundIntoVariable(stackPtr, atom, varLookup, true, 0)
			if err != nil {
				return
			}
			atoms[i] = &parser.Node{Type: parser.NTAddr, Value: yx}
			buf.Write(code)
		}
	}

	var varIndex int32
	switch callee.Type {
	case parser.NTAtom:
		varIndex = varLookup.GetRelPosition(callee.Value.(string))
		if varIndex == -1 {
			err = fmt.Errorf(ERR_UNDECLARED_VARIABLE, callee)
			return
		}
	case parser.NTCompound:
		code, yx, stackPtr, err = compileCompoundIntoVariable(stackPtr, callee, varLookup, true, 0)
		if err != nil {
			return
		}

		varIndex = yx
		buf.Write(code)
	case parser.NTAddr:
		varIndex = callee.Value.(int32)
	case parser.NTString:
		buf.WriteByte(base.OP_SET_STR)
		buf.WriteInt32(int32(stackPtr))
		buf.WriteString(callee.Value.(string))
		varIndex = int32(stackPtr)
		stackPtr++
	default:
		err = fmt.Errorf("invalid callee: %+v", callee)
		return
	}

	for i := 0; i < len(atoms); i++ {
		err = fill1(buf, atoms[i], varLookup, base.OP_PUSH, base.OP_PUSH_NUM, base.OP_PUSH_STR)
		if err != nil {
			return
		}
	}

	buf.WriteByte(base.OP_CALL)
	buf.WriteInt32(varIndex)

	return buf.Bytes(), base.REG_A, stackPtr, nil
}

func compileLambdaOp(stackPtr int16, atoms []*parser.Node, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	newLookup := base.NewCMap()
	newLookup.Parent = varLookup

	params := atoms[1]
	if params.Type != parser.NTCompound {
		err = fmt.Errorf("invalid lambda parameters: %+v", atoms[0])
		return
	}

	for i, p := range params.Compound {
		newLookup.M[p.Value.(string)] = int16(i)
	}

	ln := len(newLookup.M)
	code, yx, _, err = compileChainOp(int16(ln), atoms[2], newLookup)
	if err != nil {
		return
	}

	code = append(code, base.OP_EOB)
	buf := base.NewBytesBuffer()
	buf.WriteByte(base.OP_LAMBDA)
	buf.WriteInt32(int32(ln))
	buf.WriteInt32(int32(len(code)))
	buf.Write(code)

	return buf.Bytes(), base.REG_A, stackPtr, nil
}
