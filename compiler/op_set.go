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

func isStoreLoadSugar(t *parser.Node) bool {
	ans := false
	if t.Type == parser.NTCompound {
		tokens := t.Compound
		if len(tokens) >= 3 {
			ans = true

			// form: [a : b : c : d ...]
			for i := 1; i < len(tokens); i += 2 {
				if r, ok := tokens[i].Value.(string); ok && r == ":" {
				} else {
					ans = false
					break
				}
			}
		}
	}
	return ans
}

func expandStoreLoadSugar(t *parser.Node) *parser.Node {
	ts := t.Compound
	tokens := make([]*parser.Node, 0, len(ts))
	tokens = append(tokens, nil)
	for i := 0; i < len(ts); i += 2 {
		tokens = append(tokens, ts[i])
	}

	return &parser.Node{
		Type:     parser.NTCompound,
		Pos:      t.Pos,
		Compound: tokens,
	}
}

func compileSetOp(stackPtr int16, atoms []*parser.Node, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	aVar := atoms[1]
	varIndex := int32(0)
	if len(atoms) < 3 {
		err = fmt.Errorf("can't set/declare without value %+v", atoms[0])
		return
	}

	aValue := atoms[2]
	storeSugar := false
	if atoms[0].Value.(string) == "move" && isStoreLoadSugar(aVar) {
		storeSugar = true
	}

	buf := base.NewBytesBuffer()
	if !storeSugar {
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
			varLookup.M[aVar.Value.(string)] = int16(newYX)
		}

		return buf.Bytes(), newYX, stackPtr, nil
	}

	fatoms := expandStoreLoadSugar(aVar).Compound
	fatoms = append(fatoms, aValue)
	return flatWrite(stackPtr, fatoms, varLookup, base.OP_STORE)
}

func compileRetOp(stackPtr int16, atoms []*parser.Node, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	atom := atoms[1]
	buf := base.NewBytesBuffer()

	switch atom.Type {
	case parser.NTAtom, parser.NTNumber, parser.NTString, parser.NTAddr:
		err = fill1(buf, atom, varLookup, base.OP_RET, base.OP_RET_NUM, base.OP_RET_STR)
		if err != nil {
			return
		}
	case parser.NTCompound:
		code, yx, stackPtr, err = extract(stackPtr, atom, varLookup)
		buf.Write(code)
		buf.WriteByte(base.OP_RET)
		buf.WriteInt32(yx)
	}
	return buf.Bytes(), yx, stackPtr, nil
}
