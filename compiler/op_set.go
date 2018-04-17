package compiler

import (
	"fmt"

	"github.com/coyove/eugine/base"
)

const (
	c_set = iota + 1
	c_declare
)

const (
	ERR_UNDECLARED_VARIABLE = "undeclared variable: %+v"
)

func isStoreLoadSugar(t *token) bool {
	ans := false
	if t.ty == TK_compound {
		tokens := t.v.([]*token)
		if len(tokens) >= 3 {
			ans = true

			// form: [a : b : c : d ...]
			for i := 1; i < len(tokens); i += 2 {
				if r, ok := tokens[i].v.(rune); ok && r == ':' {
				} else {
					ans = false
					break
				}
			}
		}
	}
	return ans
}

func expandStoreLoadSugar(t *token) *token {
	ts := t.v.([]*token)
	tokens := make([]*token, 0, len(ts))
	tokens = append(tokens, nil)
	for i := 0; i < len(ts); i += 2 {
		tokens = append(tokens, ts[i])
	}

	return &token{ty: TK_compound, v: tokens}
}

func compileSetOp(stackPtr int16, atoms []*token, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	aVar := atoms[1]
	varIndex := int32(0)
	if len(atoms) < 3 {
		err = fmt.Errorf("can't set/declare without value %+v", atoms[0])
		return
	}

	aValue := atoms[2]
	storeSugar := false
	if atoms[0].v.(string) == "set" && isStoreLoadSugar(aVar) {
		storeSugar = true
	}

	buf := base.NewBytesBuffer()
	if !storeSugar {
		var newYX int32
		if atoms[0].v.(string) == "var" {
			// compound has its own logic, we won't incr stack here
			if aValue.ty != TK_compound {
				newYX = int32(stackPtr)
				stackPtr++
			}
		} else {
			varIndex = varLookup.GetRelPosition(aVar.v.(string))
			if varIndex == -1 {
				err = fmt.Errorf(ERR_UNDECLARED_VARIABLE, aVar)
				return
			}
			newYX = varIndex
		}

		switch aValue.ty {
		case TK_atomic:
			valueIndex := varLookup.GetRelPosition(aValue.v.(string))
			if valueIndex == -1 {
				err = fmt.Errorf(ERR_UNDECLARED_VARIABLE, aValue)
				return
			}

			buf.WriteByte(base.OP_SET)
			buf.WriteInt32(newYX)
			buf.WriteInt32(valueIndex)
		case TK_number:
			buf.WriteByte(base.OP_SET_NUM)
			buf.WriteInt32(newYX)
			buf.WriteDouble(aValue.v.(float64))
		case TK_string:
			buf.WriteByte(base.OP_SET_STR)
			buf.WriteInt32(newYX)
			buf.WriteString(aValue.v.(string))
		case TK_compound:
			code, newYX, stackPtr, err = compileCompoundIntoVariable(stackPtr, aValue, varLookup,
				atoms[0].v.(string) == "var", varIndex)
			if err != nil {
				return
			}
			buf.Write(code)
		}

		if atoms[0].v.(string) == "var" {
			varLookup.M[aVar.v.(string)] = int16(newYX)
		}

		return buf.Bytes(), newYX, stackPtr, nil
	}

	fatoms := expandStoreLoadSugar(aVar).v.([]*token)
	fatoms = append(fatoms, aValue)
	return flatWrite(stackPtr, fatoms, varLookup, base.OP_STORE)
}

func compileRetOp(stackPtr int16, atoms []*token, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	atom := atoms[1]
	buf := base.NewBytesBuffer()

	switch atom.ty {
	case TK_atomic, TK_number, TK_string, TK_addr:
		err = fill1(buf, atom, varLookup, base.OP_RET, base.OP_RET_NUM, base.OP_RET_STR)
		if err != nil {
			return
		}
	case TK_compound:
		code, yx, stackPtr, err = extract(stackPtr, atom, varLookup)
		buf.Write(code)
		buf.WriteByte(base.OP_RET)
		buf.WriteInt32(yx)
	}
	return buf.Bytes(), yx, stackPtr, nil
}
