package compiler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/coyove/bracket/base"
	"github.com/coyove/bracket/parser"
)

var opMapping map[string]func(int16, []*parser.Node, *base.CMap) ([]byte, int32, int16, error)

var flatOpMapping map[string]bool

func init() {
	opMapping = make(map[string]func(int16, []*parser.Node, *base.CMap) ([]byte, int32, int16, error))
	opMapping["set"] = compileSetOp
	opMapping["move"] = compileSetOp
	opMapping["ret"] = compileRetOp(base.OP_RET, base.OP_RET_NUM, base.OP_RET_STR)
	opMapping["yield"] = compileRetOp(base.OP_YIELD, base.OP_YIELD_NUM, base.OP_YIELD_STR)
	opMapping["lambda"] = compileLambdaOp
	opMapping["if"] = compileIfOp
	opMapping["while"] = compileWhileOp
	opMapping["continue"] = compileContinueBreakOp
	opMapping["break"] = compileContinueBreakOp
	opMapping["call"] = compileCallOp
	opMapping["list"] = compileListOp
	opMapping["map"] = compileMapOp

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
	flatOpMapping["~"] = true
	flatOpMapping["&"] = true
	flatOpMapping["|"] = true
	flatOpMapping["^"] = true
	flatOpMapping["<<"] = true
	flatOpMapping[">>"] = true
	flatOpMapping["len"] = true
	flatOpMapping["store"] = true
	flatOpMapping["load"] = true
	flatOpMapping["safestore"] = true
	flatOpMapping["safeload"] = true
	flatOpMapping["assert"] = true
	flatOpMapping["nil"] = true
	flatOpMapping["true"] = true
	flatOpMapping["false"] = true
	flatOpMapping["dup"] = true
}

func fill1(buf *base.BytesReader, n *parser.Node, varLookup *base.CMap, ops ...byte) (err error) {
	switch n.Type {
	case parser.NTAtom:
		varIndex := varLookup.GetRelPosition(n.Value.(string))
		if varIndex == -1 {
			err = fmt.Errorf(ERR_UNDECLARED_VARIABLE, n)
			return
		}
		buf.WriteByte(ops[0])
		buf.WriteInt32(varIndex)
	case parser.NTNumber:
		buf.WriteByte(ops[1])
		buf.WriteDouble(n.Value.(float64))
	case parser.NTString:
		buf.WriteByte(ops[2])
		buf.WriteString(n.Value.(string))
	case parser.NTAddr:
		buf.WriteByte(ops[0])
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
		code, yx, stackPtr, err = compile(stackPtr, n.Compound, varLookup)
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
	name, ok := nodes[0].Value.(string)
	if ok {
		if name == "chain" {
			return compileChainOp(stackPtr, &parser.Node{
				Type:     parser.NTCompound,
				Compound: nodes,
			}, varLookup)
		}

		f := opMapping[name]
		if f == nil {
			if flatOpMapping[name] {
				return compileFlatOp(stackPtr, nodes, varLookup)
			}
			panic(name)
		}

		return f(stackPtr, nodes, varLookup)
	}

	panic(nodes[0].Value)
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

func compileNode(n *parser.Node) (code []byte, err error) {
	varLookup := base.NewCMap()
	for i, n := range base.CoreLibNames {
		varLookup.M[n] = int16(i)
	}

	code, _, _, err = compileChainOp(int16(len(varLookup.M)), n, varLookup)
	if err != nil {
		return
	}

	code = append(code, base.OP_EOB)
	return code, nil
}

func LoadFile(path string) ([]byte, error) {
	code, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	n, err := parser.Parse(bytes.NewReader(code), path)
	if err != nil {
		return nil, err
	}

	n.Dump(os.Stderr)
	return compileNode(n)
}
