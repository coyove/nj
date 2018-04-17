package compiler

import (
	"fmt"

	"github.com/coyove/eugine/base"
	"github.com/coyove/eugine/vm"
)

var opMapping map[string]func(int16, []*token, *base.CMap) ([]byte, int32, int16, error)

var flatOpMapping map[string]bool

func init() {
	opMapping = make(map[string]func(int16, []*token, *base.CMap) ([]byte, int32, int16, error))
	opMapping["var"] = compileSetOp
	opMapping["set"] = compileSetOp
	opMapping["ret"] = compileRetOp
	opMapping["inc"] = compileIncOp
	opMapping["lambda"] = compileLambdaOp
	opMapping["if"] = compileIfOp
	opMapping["while"] = compileWhileOp
	opMapping["continue"] = compileContinueBreakOp
	opMapping["break"] = compileContinueBreakOp
	opMapping["typeof"] = compileTypeofOp

	flatOpMapping = make(map[string]bool)
	flatOpMapping["+"] = true
	flatOpMapping["-"] = true
	flatOpMapping["*"] = true
	flatOpMapping["/"] = true
	flatOpMapping["%"] = true
	flatOpMapping["<"] = true
	flatOpMapping["<="] = true
	flatOpMapping[">"] = true
	flatOpMapping[">="] = true
	flatOpMapping["eq"] = true
	flatOpMapping["neq"] = true
	flatOpMapping["not"] = true
	flatOpMapping["and"] = true
	flatOpMapping["or"] = true
	flatOpMapping["xor"] = true
	flatOpMapping["b/not"] = true
	flatOpMapping["b/and"] = true
	flatOpMapping["b/or"] = true
	flatOpMapping["b/xor"] = true
	flatOpMapping["b/lsh"] = true
	flatOpMapping["b/rsh"] = true
	flatOpMapping["list"] = true
	flatOpMapping["map"] = true
	flatOpMapping["len"] = true
	flatOpMapping["store"] = true
	flatOpMapping["load"] = true
	flatOpMapping["assert"] = true
	flatOpMapping["expand"] = true
	flatOpMapping["nil"] = true
	flatOpMapping["bytes"] = true
	flatOpMapping["true"] = true
	flatOpMapping["false"] = true
}

func fill1(buf *base.BytesReader, atom *token, varLookup *base.CMap, a, i, s byte) (err error) {
	switch atom.ty {
	case TK_atomic:
		varIndex := varLookup.GetRelPosition(atom.v.(string))
		if varIndex == -1 {
			err = fmt.Errorf(ERR_UNDECLARED_VARIABLE, atom)
			return
		}
		buf.WriteByte(a)
		buf.WriteInt32(varIndex)
	case TK_number:
		buf.WriteByte(i)
		buf.WriteDouble(atom.v.(float64))
	case TK_string:
		buf.WriteByte(s)
		buf.WriteString(atom.v.(string))
	case TK_addr:
		buf.WriteByte(a)
		buf.WriteInt32(atom.v.(int32))
	default:
		return fmt.Errorf("fill1 unknown type")
	}
	return nil
}

func compileCompoundIntoVariable(
	stackPtr int16,
	compound *token,
	varLookup *base.CMap,
	intoNewVar bool,
	intoExistedVar int32,
) (code []byte, yx int32, newStackPtr int16, err error) {
	buf := base.NewBytesBuffer()
	if isStoreLoadSugar(compound) {
		code, yx, stackPtr, err = flatWrite(stackPtr, expandStoreLoadSugar(compound).v.([]*token), varLookup, base.OP_LOAD)
		buf.Write(code)
		buf.WriteByte(base.OP_SET)
		if intoNewVar {
			yx = int32(stackPtr)
			stackPtr++
		} else {
			yx = intoExistedVar
		}
		buf.WriteInt32(yx)
		buf.WriteInt32(base.REG_A)
		return buf.Bytes(), yx, stackPtr, err
	}

	var newYX int32
	code, newYX, stackPtr, err = compileImpl(stackPtr, compound.v.([]*token), varLookup)
	if err != nil {
		return nil, 0, 0, err
	}

	buf.Write(code)
	buf.WriteByte(base.OP_SET)
	if intoNewVar {
		yx = int32(stackPtr)
		stackPtr++
	} else {
		yx = intoExistedVar
	}
	buf.WriteInt32(yx)
	buf.WriteInt32(newYX)
	return buf.Bytes(), yx, stackPtr, nil
}

func extract(stackPtr int16, atom *token, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	var varIndex int32

	switch atom.ty {
	case TK_atomic:
		varIndex = varLookup.GetRelPosition(atom.v.(string))
		if varIndex == -1 {
			err = fmt.Errorf(ERR_UNDECLARED_VARIABLE, atom)
			return
		}
	case TK_addr:
		varIndex = atom.v.(int32)
	default:
		if isStoreLoadSugar(atom) {
			code, yx, stackPtr, err = flatWrite(stackPtr, expandStoreLoadSugar(atom).v.([]*token), varLookup, base.OP_LOAD)
		} else {
			code, yx, stackPtr, err = compileImpl(stackPtr, atom.v.([]*token), varLookup)
		}
		if err != nil {
			return
		}
		varIndex = yx
	}
	return code, varIndex, stackPtr, nil
}

func compileImpl(stackPtr int16, atoms []*token, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	if len(atoms) == 0 {
		return nil, base.REG_A, stackPtr, nil
	}

	if name, ok := atoms[0].v.(string); ok {
		f := opMapping[name]
		if f == nil {
			if flatOpMapping[name] {
				return compileFlatOp(stackPtr, atoms, varLookup)
			}

			if _, ok := vm.LibLookup[name]; ok {
				return compileFlatOp(stackPtr, atoms, varLookup)
			}

			return compileCallOp(stackPtr, atoms, varLookup)
		}

		return f(stackPtr, atoms, varLookup)
	}
	return compileCallOp(stackPtr, atoms, varLookup)
}

func compile(stackPtr int16, atoms []*token, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	buf := base.NewBytesBuffer()

	for i := 0; i < len(atoms); i++ {
		a := atoms[i]
		if a.ty != TK_compound {
			err = fmt.Errorf("every atom in the chain must be a compound: %+v", a)
			return
		}

		code, yx, stackPtr, err = compileImpl(stackPtr, a.v.([]*token), varLookup)
		if err != nil {
			return
		}

		buf.Write(code)
	}

	buf.WriteByte(base.OP_EOB)
	return buf.Bytes(), yx, stackPtr, err
}
