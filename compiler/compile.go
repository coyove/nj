package compiler

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/coyove/bracket/base"
	"github.com/coyove/bracket/parser"
)

type compileFunc func(int16, []*parser.Node, *base.CMap) ([]byte, int32, int16, error)

var opMapping map[string]compileFunc

var flatOpMapping map[string]bool

func init() {
	clearI := func(f compileFunc) compileFunc {
		return func(s int16, n []*parser.Node, v *base.CMap) ([]byte, int32, int16, error) {
			a, b, c, d := f(s, n, v)
			v.I = nil
			v.Is = nil
			return a, b, c, d
		}
	}

	opMapping = make(map[string]compileFunc)
	opMapping["set"] = clearI(compileSetOp)
	opMapping["move"] = clearI(compileSetOp)
	opMapping["ret"] = clearI(compileRetOp(base.OP_RET, base.OP_RET_NUM, base.OP_RET_STR))
	opMapping["yield"] = clearI(compileRetOp(base.OP_YIELD, base.OP_YIELD_NUM, base.OP_YIELD_STR))
	opMapping["lambda"] = clearI(compileLambdaOp)
	opMapping["if"] = clearI(compileIfOp)
	opMapping["while"] = clearI(compileWhileOp)
	opMapping["continue"] = clearI(compileContinueBreakOp)
	opMapping["break"] = clearI(compileContinueBreakOp)
	opMapping["call"] = clearI(compileCallOp)
	opMapping["list"] = clearI(compileListOp)
	opMapping["map"] = clearI(compileMapOp)

	flatOpMapping = map[string]bool{
		"+": true, "-": true, "*": true, "/": true, "%": true,
		"<": true, "<=": true, ">": true, ">=": true, "eq": true, "neq": true, "not": true, "and": true, "or": true, "xor": true,
		"~": true, "&": true, "|": true, "^": true, "<<": true, ">>": true,
		"store": true, "load": true, "safestore": true, "safeload": true,
		"assert": true, "nil": true, "true": true, "false": true,
	}
}

func fill1(buf *base.BytesWriter, n *parser.Node, varLookup *base.CMap, ops ...byte) (err error) {
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
	buf := base.NewBytesWriter()

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
	buf := base.NewBytesWriter()
	varLookup.I = nil

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

	// n.Dump(os.Stderr)
	return compileNode(n)
}
