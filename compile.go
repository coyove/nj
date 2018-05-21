package potatolang

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/coyove/potatolang/parser"
)

type compileFunc func(int16, []*parser.Node, *symtable) ([]byte, int32, int16, error)

var opMapping map[string]compileFunc

var flatOpMapping map[string]bool

func init() {
	clearI := func(f compileFunc) compileFunc {
		return func(s int16, n []*parser.Node, v *symtable) ([]byte, int32, int16, error) {
			a, b, c, d := f(s, n, v)
			v.I = nil
			v.Is = nil
			return a, b, c, d
		}
	}

	opMapping = make(map[string]compileFunc)
	opMapping["set"] = clearI(compileSetOp)
	opMapping["move"] = clearI(compileSetOp)
	opMapping["ret"] = clearI(compileRetOp(OP_RET, OP_RET_NUM, OP_RET_STR))
	opMapping["yield"] = clearI(compileRetOp(OP_YIELD, OP_YIELD_NUM, OP_YIELD_STR))
	opMapping["lambda"] = clearI(compileLambdaOp)
	opMapping["if"] = clearI(compileIfOp)
	opMapping["while"] = clearI(compileWhileOp)
	opMapping["continue"] = clearI(compileContinueBreakOp)
	opMapping["break"] = clearI(compileContinueBreakOp)
	opMapping["call"] = clearI(compileCallOp)
	opMapping["list"] = clearI(compileListOp)
	opMapping["map"] = clearI(compileMapOp)
	opMapping["or"] = clearI(compileAndOrOp(OP_IF))
	opMapping["and"] = clearI(compileAndOrOp(OP_IFNOT))

	flatOpMapping = map[string]bool{
		"+": true, "-": true, "*": true, "/": true, "%": true,
		"<": true, "<=": true, "eq": true, "neq": true, "not": true,
		"~": true, "&": true, "|": true, "^": true, "<<": true, ">>": true,
		"store": true, "load": true, "safestore": true, "safeload": true,
		"assert": true, "nil": true, "true": true, "false": true,
	}
}

func fill1(buf *BytesWriter, n *parser.Node, table *symtable, ops ...byte) (err error) {
	switch n.Type {
	case parser.NTAtom:
		varIndex := table.GetRelPosition(n.Value.(string))
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
	sp int16,
	compound *parser.Node,
	table *symtable,
	intoNewVar bool,
	intoExistedVar int32,
) (code []byte, yx int32, newsp int16, err error) {
	buf := NewBytesWriter()

	var newYX int32
	code, newYX, sp, err = compile(sp, compound.Compound, table)
	if err != nil {
		return nil, 0, 0, err
	}

	buf.Write(code)
	buf.WriteByte(OP_SET)
	if intoNewVar {
		yx = int32(sp)
		sp++
	} else {
		yx = intoExistedVar
	}
	buf.WriteInt32(yx)
	buf.WriteInt32(newYX)
	return buf.Bytes(), yx, sp, nil
}

func extract(sp int16, n *parser.Node, table *symtable) (code []byte, yx int32, newsp int16, err error) {
	var varIndex int32

	switch n.Type {
	case parser.NTAtom:
		varIndex = table.GetRelPosition(n.Value.(string))
		if varIndex == -1 {
			err = fmt.Errorf(ERR_UNDECLARED_VARIABLE, n)
			return
		}
	case parser.NTAddr:
		varIndex = n.Value.(int32)
	default:
		code, yx, sp, err = compile(sp, n.Compound, table)
		if err != nil {
			return
		}
		varIndex = yx
	}
	return code, varIndex, sp, nil
}

func compile(sp int16, nodes []*parser.Node, table *symtable) (code []byte, yx int32, newsp int16, err error) {
	if len(nodes) == 0 {
		return nil, REG_A, sp, nil
	}
	name, ok := nodes[0].Value.(string)
	if ok {
		if name == "chain" {
			return compileChainOp(sp, &parser.Node{
				Type:     parser.NTCompound,
				Compound: nodes,
			}, table)
		}

		f := opMapping[name]
		if f == nil {
			if flatOpMapping[name] {
				return compileFlatOp(sp, nodes, table)
			}
			panic(name)
		}

		return f(sp, nodes, table)
	}

	panic(nodes[0].Value)
}

func compileChainOp(sp int16, chain *parser.Node, table *symtable) (code []byte, yx int32, newsp int16, err error) {
	buf := NewBytesWriter()
	table.I = nil

	for _, a := range chain.Compound {
		if a.Type != parser.NTCompound {
			continue
		}
		if table.LineInfo && len(a.Compound) > 0 && a.Compound[0].Pos.Source != "" {
			buf.WriteByte(OP_LINE)
			buf.WriteString(a.Compound[0].Pos.String())
		}
		code, yx, sp, err = compile(sp, a.Compound, table)
		if err != nil {
			return
		}
		buf.Write(code)
	}

	return buf.Bytes(), yx, sp, err
}

func compileNode(n *parser.Node, lineinfo bool) (code []byte, err error) {
	table := &symtable{M: make(map[string]int16), LineInfo: lineinfo}
	for i, n := range CoreLibNames {
		table.M[n] = int16(i)
	}

	code, _, _, err = compileChainOp(int16(len(table.M)), n, table)
	if err != nil {
		return
	}

	code = append(code, OP_EOB)
	return code, nil
}

func LoadFile(path string, lineinfo bool) ([]byte, error) {
	code, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	n, err := parser.Parse(bytes.NewReader(code), path)
	if err != nil {
		return nil, err
	}
	n.Dump(os.Stderr)
	return compileNode(n, lineinfo)
}

func LoadString(code string, lineinfo bool) ([]byte, error) {
	n, err := parser.Parse(bytes.NewReader([]byte(code)), "mem")
	if err != nil {
		return nil, err
	}
	return compileNode(n, lineinfo)
}
