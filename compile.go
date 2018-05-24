package potatolang

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/coyove/potatolang/parser"
)

type compileFunc func(uint16, []*parser.Node, *symtable) ([]uint16, uint32, uint16, error)

var opMapping map[string]compileFunc

var flatOpMapping map[string]bool

func init() {
	clearI := func(f compileFunc) compileFunc {
		return func(s uint16, n []*parser.Node, v *symtable) ([]uint16, uint32, uint16, error) {
			a, b, c, d := f(s, n, v)
			v.im = nil
			v.ims = nil
			return a, b, c, d
		}
	}

	opMapping = make(map[string]compileFunc)
	opMapping["set"] = clearI(compileSetOp)
	opMapping["move"] = clearI(compileSetOp)
	opMapping["ret"] = clearI(compileRetOp(OP_RET, OP_RETK))
	opMapping["yield"] = clearI(compileRetOp(OP_YIELD, OP_YIELDK))
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
	opMapping["inc"] = clearI(compileIncOp)

	flatOpMapping = map[string]bool{
		"+": true, "-": true, "*": true, "/": true, "%": true,
		"<": true, "<=": true, "eq": true, "neq": true, "not": true,
		"~": true, "&": true, "|": true, "^": true, "<<": true, ">>": true,
		"store": true, "load": true, "safestore": true, "safeload": true,
		"assert": true, "nil": true, "true": true, "false": true,
	}
}

func fill(buf *BytesWriter, n *parser.Node, table *symtable, op, opk uint16) (err error) {
	switch n.Type {
	case parser.NTAtom:
		addr, ok := table.get(n.Value.(string))
		if !ok {
			return fmt.Errorf(ERR_UNDECLARED_VARIABLE, n)
		}
		buf.Write16(op)
		buf.Write32(addr)
	case parser.NTNumber, parser.NTString:
		buf.Write16(opk)
		buf.Write16(table.addConst(n.Value))
	case parser.NTAddr:
		buf.Write16(op)
		buf.Write32(n.Value.(uint32))
	default:
		return fmt.Errorf("unknown type: %d", n.Type)
	}
	return nil
}

func compileCompoundIntoVariable(
	sp uint16,
	compound *parser.Node,
	table *symtable,
	intoNewVar bool,
	intoExistedVar uint32,
) (code []uint16, yx uint32, newsp uint16, err error) {
	buf := NewBytesWriter()

	var newYX uint32
	code, newYX, sp, err = compile(sp, compound.Compound, table)
	if err != nil {
		return nil, 0, 0, err
	}

	buf.Write(code)
	buf.Write16(OP_SET)
	if intoNewVar {
		yx = uint32(sp)
		sp++
	} else {
		yx = intoExistedVar
	}
	buf.Write32(yx)
	buf.Write32(newYX)
	return buf.data, yx, sp, nil
}

func extract(sp uint16, n *parser.Node, table *symtable) (code []uint16, yx uint32, newsp uint16, err error) {
	var varIndex uint32

	switch n.Type {
	case parser.NTAtom:
		var ok bool
		varIndex, ok = table.get(n.Value.(string))
		if !ok {
			err = fmt.Errorf(ERR_UNDECLARED_VARIABLE, n)
			return
		}
	case parser.NTAddr:
		varIndex = n.Value.(uint32)
	default:
		code, yx, sp, err = compile(sp, n.Compound, table)
		if err != nil {
			return
		}
		varIndex = yx
	}
	return code, varIndex, sp, nil
}

func compile(sp uint16, nodes []*parser.Node, table *symtable) (code []uint16, yx uint32, newsp uint16, err error) {
	if len(nodes) == 0 {
		return nil, regA, sp, nil
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

func compileChainOp(sp uint16, chain *parser.Node, table *symtable) (code []uint16, yx uint32, newsp uint16, err error) {
	buf := NewBytesWriter()
	table.im = nil

	for _, a := range chain.Compound {
		if a.Type != parser.NTCompound {
			continue
		}
		if table.lineInfo {
			for _, n := range a.Compound {
				if n.Pos.Source != "" {
					buf.Write16(OP_LINE)
					buf.WriteString(n.Pos.String())
					break
				}
			}
		}
		code, yx, sp, err = compile(sp, a.Compound, table)
		if err != nil {
			return
		}
		buf.Write(code)
	}

	return buf.data, yx, sp, err
}

func compileNode(n *parser.Node, lineinfo bool) (cls *Closure, err error) {
	table := newsymtable()
	table.lineInfo = lineinfo
	for i, n := range CoreLibNames {
		table.sym[n] = uint16(i)
	}

	code, _, _, err := compileChainOp(uint16(len(table.sym)), n, table)
	if err != nil {
		return nil, err
	}

	code = append(code, OP_EOB)
	consts := make([]Value, len(table.consts))
	for i, k := range table.consts {
		switch k.ty {
		case Tnumber:
			consts[i] = NewNumberValue(k.value.(float64))
		case Tstring:
			consts[i] = NewStringValue(k.value.(string))
		}
	}
	return NewClosure(code, consts, nil, 0, false, false), err
}

func LoadFile(path string, lineinfo bool) (*Closure, error) {
	code, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	n, err := parser.Parse(bytes.NewReader(code), path)
	if err != nil {
		return nil, err
	}

	return compileNode(n, lineinfo)
}

func LoadString(code string, lineinfo bool) (*Closure, error) {
	n, err := parser.Parse(bytes.NewReader([]byte(code)), "mem")
	if err != nil {
		return nil, err
	}
	return compileNode(n, lineinfo)
}
