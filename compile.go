package potatolang

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/coyove/potatolang/parser"
)

type kinfo struct {
	ty    byte
	value interface{}
}

// symtable is responsible for recording extra states of compilation
type symtable struct {
	// variable name lookup
	parent *symtable
	sym    map[string]uint16

	// flat op immediate value
	im  *float64
	ims *string

	// has yield op
	y bool

	envescape bool

	noredecl bool

	// record line info at chain
	lineInfo bool

	regs [4]struct {
		addr  uint32
		kaddr uint16
		k     bool
	}

	continueNode []*parser.Node

	consts         []kinfo
	constStringMap map[string]uint16
	constFloatMap  map[float64]uint16
}

func newsymtable() *symtable {
	t := &symtable{
		sym:            make(map[string]uint16),
		consts:         make([]kinfo, 0),
		constStringMap: make(map[string]uint16),
		constFloatMap:  make(map[float64]uint16),
		continueNode:   make([]*parser.Node, 0),
	}
	for i := range t.regs {
		t.regs[i].addr = regA
		t.regs[i].k = false
	}
	return t
}

func (m *symtable) get(varname string) (uint32, bool) {
	depth := uint32(0)

	for m != nil {
		k, e := m.sym[varname]
		if e {
			return (depth << 16) | uint32(k), true
		}

		depth++
		m = m.parent
	}

	return 0, false
}

func (m *symtable) put(varname string, addr uint16) {
	m.sym[varname] = addr
}

func (m *symtable) clearRegRecord(addr uint32) {
	for i, x := range m.regs {
		if !x.k && x.addr == addr {
			m.regs[i].addr = regA
		}
	}
}

func (m *symtable) clearAllRegRecords() {
	for i := range m.regs {
		m.regs[i].k = false
		m.regs[i].addr = regA
	}
}

func (m *symtable) addConst(v interface{}) uint16 {
	var k kinfo
	k.value = v

	switch v.(type) {
	case float64:
		k.ty = Tnumber
		if i, ok := m.constFloatMap[v.(float64)]; ok {
			return i
		}
	case string:
		k.ty = Tstring
		if i, ok := m.constStringMap[v.(string)]; ok {
			return i
		}
	default:
		panic("shouldn't happen")
	}

	m.consts = append(m.consts, k)
	idx := uint16(len(m.consts))

	switch v.(type) {
	case float64:
		m.constFloatMap[v.(float64)] = idx
	case string:
		m.constStringMap[v.(string)] = idx
	}

	return idx
}

type compileFunc func(uint16, []*parser.Node, *symtable) ([]uint64, uint32, uint16, error)

var opMapping map[string]compileFunc

var flatOpMapping map[string]byte

func init() {
	clearI := func(f compileFunc) compileFunc {
		return func(s uint16, n []*parser.Node, v *symtable) ([]uint64, uint32, uint16, error) {
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
	opMapping["for"] = clearI(compileWhileOp)
	opMapping["continue"] = clearI(compileContinueBreakOp)
	opMapping["break"] = clearI(compileContinueBreakOp)
	opMapping["call"] = clearI(compileCallOp)
	opMapping["map"] = clearI(compileMapOp)
	opMapping["or"] = clearI(compileAndOrOp(OP_IF))
	opMapping["and"] = clearI(compileAndOrOp(OP_IFNOT))
	opMapping["inc"] = clearI(compileIncOp)

	flatOpMapping = map[string]byte{
		"+": OP_ADD, "-": OP_SUB, "*": OP_MUL, "/": OP_DIV, "%": OP_MOD,
		"<": OP_LESS, "<=": OP_LESS_EQ, "==": OP_EQ, "!=": OP_NEQ, "!": OP_NOT,
		"~": OP_BIT_NOT, "&": OP_BIT_AND, "|": OP_BIT_OR, "^": OP_BIT_XOR, "<<": OP_BIT_LSH, ">>": OP_BIT_RSH,
		"#": OP_POP, "store": OP_STORE, "load": OP_LOAD, "assert": OP_ASSERT,
	}
}

var registerOpMappings = map[byte]int{OP_R0: 0, OP_R1: 1, OP_R2: 2, OP_R3: 3}

func fill(buf *opwriter, n *parser.Node, table *symtable, op, opk byte) (err error) {
	idx, isreg := registerOpMappings[op]
	if isreg {
		if rs := table.regs[idx]; rs.k {
			if n.Type == parser.NTNumber || n.Type == parser.NTString {
				if kidx := table.addConst(n.Value); kidx == rs.kaddr {
					// the register contains what we want already,
					return
				}
			}
		} else {
			var addr uint32 = regA
			switch n.Type {
			case parser.NTAtom:
				addr, _ = table.get(n.Value.(string))
			case parser.NTAddr:
				addr = n.Value.(uint32)
			}
			if addr != regA && rs.addr == addr {
				// the register contains what we want already,
				return
			}
		}
	}

	switch n.Type {
	case parser.NTAtom:
		if n.Value.(string) == "nil" {
			buf.WriteOP(opk, 0, 0)
			if isreg {
				table.regs[idx].k, table.regs[idx].kaddr = true, 0
			}
		} else {
			addr, ok := table.get(n.Value.(string))
			if !ok {
				return fmt.Errorf(ERR_UNDECLARED_VARIABLE, n)
			}
			buf.WriteOP(op, addr, 0)
			if isreg {
				table.regs[idx].k, table.regs[idx].addr = false, addr
			}
		}
	case parser.NTNumber, parser.NTString:
		kidx := table.addConst(n.Value)
		buf.WriteOP(opk, uint32(kidx), 0)
		if isreg {
			table.regs[idx].k, table.regs[idx].kaddr = true, kidx
		}
	case parser.NTAddr:
		buf.WriteOP(op, n.Value.(uint32), 0)
		if isreg {
			table.regs[idx].k, table.regs[idx].addr = false, n.Value.(uint32)
		}
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
) (code []uint64, yx uint32, newsp uint16, err error) {
	buf := newopwriter()

	var newYX uint32
	code, newYX, sp, err = compile(sp, compound.Compound, table)
	if err != nil {
		return nil, 0, 0, err
	}

	buf.Write(code)
	if intoNewVar {
		yx = uint32(sp)
		sp++
	} else {
		yx = intoExistedVar
		table.clearRegRecord(yx)
	}
	buf.WriteOP(OP_SET, yx, newYX)
	return buf.data, yx, sp, nil
}

func extract(sp uint16, n *parser.Node, table *symtable) (code []uint64, yx uint32, newsp uint16, err error) {
	var varIndex uint32

	switch n.Type {
	case parser.NTAtom:
		if n.Value.(string) == "nil" {
			buf := newopwriter()
			buf.WriteOP(OP_SETK, uint32(sp), 0)
			return buf.data, uint32(sp), sp + 1, nil
		} else {
			var ok bool
			varIndex, ok = table.get(n.Value.(string))
			if !ok {
				err = fmt.Errorf(ERR_UNDECLARED_VARIABLE, n)
				return
			}
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

func compile(sp uint16, nodes []*parser.Node, table *symtable) (code []uint64, yx uint32, newsp uint16, err error) {
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
			if flatOpMapping[name] != 0 {
				return compileFlatOp(sp, nodes, table)
			}
			panic(name)
		}

		return f(sp, nodes, table)
	}

	nodes[0].Dump(os.Stderr)
	log.Panicf("invalid op: %v", nodes)
	return
}

func compileChainOp(sp uint16, chain *parser.Node, table *symtable) (code []uint64, yx uint32, newsp uint16, err error) {
	buf := newopwriter()
	table.im = nil

	for _, a := range chain.Compound {
		if a.Type != parser.NTCompound {
			continue
		}
		if table.lineInfo {
			for _, n := range a.Compound {
				if n.Pos.Source != "" {
					buf.WriteOP(OP_LINE, 0, 0)
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

	code = append(code, makeop(OP_EOB, 0, 0))
	consts := make([]Value, len(table.consts)+1)
	for i, k := range table.consts {
		switch k.ty {
		case Tnumber:
			consts[i+1] = NewNumberValue(k.value.(float64))
		case Tstring:
			consts[i+1] = NewStringValue(k.value.(string))
		}
	}
	cls = NewClosure(code, consts, nil, 0, false, false, false, false)
	cls.lastenv = NewEnv(nil)
	for _, name := range CoreLibNames {
		cls.lastenv.SPush(CoreLibs[name])
	}
	return cls, err
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
	// n.Dump(os.Stderr)
	// panic(10)
	return compileNode(n, lineinfo)
}

func LoadString(code string, lineinfo bool) (*Closure, error) {
	n, err := parser.Parse(bytes.NewReader([]byte(code)), "mem")
	if err != nil {
		return nil, err
	}
	return compileNode(n, lineinfo)
}
