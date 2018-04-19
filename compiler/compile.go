package compiler

import (
	"fmt"

	"github.com/coyove/bracket/base"
	"github.com/coyove/bracket/parser"
)

var opMapping map[string]func(int16, []*parser.Node, *base.CMap) ([]byte, int32, int16, error)

var flatOpMapping map[string]bool

func init() {
	opMapping = make(map[string]func(int16, []*parser.Node, *base.CMap) ([]byte, int32, int16, error))
	opMapping["set"] = compileSetOp
	opMapping["move"] = compileSetOp
	opMapping["ret"] = compileRetOp
	opMapping["inc"] = compileIncOp
	opMapping["lambda"] = compileLambdaOp
	opMapping["if"] = compileIfOp
	opMapping["while"] = compileWhileOp
	opMapping["continue"] = compileContinueBreakOp
	opMapping["break"] = compileContinueBreakOp
	opMapping["typeof"] = compileTypeofOp
	opMapping["call"] = compileCallOp

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

func fill1(buf *base.BytesReader, n *parser.Node, varLookup *base.CMap, a, i, s byte) (err error) {
	switch n.Type {
	case parser.NTAtom:
		varIndex := varLookup.GetRelPosition(n.Value.(string))
		if varIndex == -1 {
			err = fmt.Errorf(ERR_UNDECLARED_VARIABLE, n)
			return
		}
		buf.WriteByte(a)
		buf.WriteInt32(varIndex)
	case parser.NTNumber:
		buf.WriteByte(i)
		buf.WriteDouble(n.Value.(float64))
	case parser.NTString:
		buf.WriteByte(s)
		buf.WriteString(n.Value.(string))
	case parser.NTAddr:
		buf.WriteByte(a)
		buf.WriteInt32(n.Value.(int32))
	default:
		return fmt.Errorf("fill1 unknown type: %d", n.Type)
	}
	return nil
}

func compileCompoundIntoVariable(
	stackPtr int16,
	compound *parser.Node,
	varLookup *base.CMap,
	intoNewVar bool,
	intoExistedVar int32,
) (code []byte, yx int32, newStackPtr int16, err error) {
	buf := base.NewBytesBuffer()
	if isStoreLoadSugar(compound) {
		code, yx, stackPtr, err = flatWrite(stackPtr, expandStoreLoadSugar(compound).Compound, varLookup, base.OP_LOAD)
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
	code, newYX, stackPtr, err = compile(stackPtr, compound.Compound, varLookup)
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

func extract(stackPtr int16, n *parser.Node, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	var varIndex int32

	switch n.Type {
	case parser.NTAtom:
		varIndex = varLookup.GetRelPosition(n.Value.(string))
		if varIndex == -1 {
			err = fmt.Errorf(ERR_UNDECLARED_VARIABLE, n)
			return
		}
	case parser.NTAddr:
		varIndex = n.Value.(int32)
	default:
		if isStoreLoadSugar(n) {
			code, yx, stackPtr, err = flatWrite(stackPtr, expandStoreLoadSugar(n).Compound, varLookup, base.OP_LOAD)
		} else {
			code, yx, stackPtr, err = compile(stackPtr, n.Compound, varLookup)
		}
		if err != nil {
			return
		}
		varIndex = yx
	}
	return code, varIndex, stackPtr, nil
}

func compile(stackPtr int16, nodes []*parser.Node, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	if len(nodes) == 0 {
		return nil, base.REG_A, stackPtr, nil
	}

	if name, ok := nodes[0].Value.(string); ok {
		f := opMapping[name]
		if f == nil {
			if flatOpMapping[name] {
				return compileFlatOp(stackPtr, nodes, varLookup)
			}
			panic(name)
		}

		return f(stackPtr, nodes, varLookup)
	}

	panic(1)
}

func compileChainOp(stackPtr int16, chain *parser.Node, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	buf := base.NewBytesBuffer()

	for _, a := range chain.Compound {
		if a.Type != parser.NTCompound {
			continue
		}

		code, yx, stackPtr, err = compile(stackPtr, a.Compound, varLookup)
		if err != nil {
			return
		}

		buf.Write(code)
	}

	return buf.Bytes(), yx, stackPtr, err
}
