package potatolang

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/coyove/potatolang/parser"
)

type kinfo struct {
	ty    byte
	value interface{}
}

// symtable is responsible for recording the state of compilation
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

	sp uint16

	regs [4]struct {
		addr  uint32
		kaddr uint16
		k     bool
	}

	continueNode []*parser.Node

	consts         []kinfo
	constStringMap map[string]uint16
	constFloatMap  map[float64]uint16

	gotoTable map[string][2]uint64
}

func newsymtable() *symtable {
	t := &symtable{
		sym:            make(map[string]uint16),
		consts:         make([]kinfo, 0),
		constStringMap: make(map[string]uint16),
		constFloatMap:  make(map[float64]uint16),
		continueNode:   make([]*parser.Node, 0),
		gotoTable:      make(map[string][2]uint64),
	}
	for i := range t.regs {
		t.regs[i].addr = regA
		t.regs[i].k = false
	}
	return t
}

func (table *symtable) get(varname string) (uint32, bool) {
	depth := uint32(0)

	for table != nil {
		k, e := table.sym[varname]
		if e {
			return (depth << 16) | uint32(k), true
		}

		depth++
		table = table.parent
	}

	return 0, false
}

func (table *symtable) put(varname string, addr uint16) {
	table.sym[varname] = addr
}

func (table *symtable) clearRegRecord(addr uint32) {
	for i, x := range table.regs {
		if !x.k && x.addr == addr {
			table.regs[i].addr = regA
		}
	}
}

func (table *symtable) clearAllRegRecords() {
	for i := range table.regs {
		table.regs[i].k = false
		table.regs[i].addr = regA
	}
}

func (table *symtable) addConst(v interface{}) uint16 {
	var k kinfo
	k.value = v

	switch v.(type) {
	case float64:
		k.ty = Tnumber
		if i, ok := table.constFloatMap[v.(float64)]; ok {
			return i
		}
	case string:
		k.ty = Tstring
		if i, ok := table.constStringMap[v.(string)]; ok {
			return i
		}
	default:
		panic("shouldn't happen")
	}

	table.consts = append(table.consts, k)
	idx := uint16(len(table.consts))

	switch v.(type) {
	case float64:
		table.constFloatMap[v.(float64)] = idx
	case string:
		table.constStringMap[v.(string)] = idx
	}

	return idx
}

var flatOpMapping = map[string]byte{
	"+": OP_ADD, "-": OP_SUB, "*": OP_MUL, "/": OP_DIV, "%": OP_MOD,
	"<": OP_LESS, "<=": OP_LESS_EQ, "==": OP_EQ, "!=": OP_NEQ, "!": OP_NOT,
	"~": OP_BIT_NOT, "&": OP_BIT_AND, "|": OP_BIT_OR, "^": OP_BIT_XOR, "<<": OP_BIT_LSH, ">>": OP_BIT_RSH,
	"#": OP_POP, "store": OP_STORE, "load": OP_LOAD, "assert": OP_ASSERT, "slice": OP_SLICE,
}

var registerOpMappings = map[byte]int{OP_R0: 0, OP_R1: 1, OP_R2: 2, OP_R3: 3}

func (table *symtable) fill(buf *packet, n *parser.Node, op, opk byte) (err error) {
	idx, isreg := registerOpMappings[op]
	if isreg {
		if rs := table.regs[idx]; rs.k {
			if n.Type == parser.Nnumber || n.Type == parser.Nstring {
				if kidx := table.addConst(n.Value); kidx == rs.kaddr {
					// the register contains what we want already,
					return
				}
			}
		} else {
			var addr uint32 = regA
			switch n.Type {
			case parser.Natom:
				addr, _ = table.get(n.Value.(string))
			case parser.Naddr:
				addr = n.Value.(uint32)
			}
			if addr != regA && rs.addr == addr {
				// the register contains what we want already,
				return
			}
		}
	}

	switch n.Type {
	case parser.Natom:
		if n.Value.(string) == "nil" {
			buf.WriteOP(opk, 0, 0)
			if isreg {
				table.regs[idx].k, table.regs[idx].kaddr = true, 0
			}
		} else {
			addr, ok := table.get(n.Value.(string))
			if !ok {
				return fmt.Errorf(errUndeclaredVariable, n)
			}
			buf.WriteOP(op, addr, 0)
			if isreg {
				table.regs[idx].k, table.regs[idx].addr = false, addr
			}
		}
	case parser.Nnumber, parser.Nstring:
		kidx := table.addConst(n.Value)
		buf.WriteOP(opk, uint32(kidx), 0)
		if isreg {
			table.regs[idx].k, table.regs[idx].kaddr = true, kidx
		}
	case parser.Naddr:
		buf.WriteOP(op, n.Value.(uint32), 0)
		if isreg {
			table.regs[idx].k, table.regs[idx].addr = false, n.Value.(uint32)
		}
	default:
		return fmt.Errorf("unknown type: %d", n.Type)
	}
	return nil
}

func (table *symtable) compileCompoundInto(compound *parser.Node, newVar bool, existedVar uint32) (code packet, yx uint32, err error) {
	buf := newpacket()

	var newYX uint32
	code, newYX, err = table.compileCompound(compound)
	if err != nil {
		return
	}

	buf.Write(code)
	if newVar {
		yx = uint32(table.sp)
		table.sp++
	} else {
		yx = existedVar
		table.clearRegRecord(yx)
	}
	buf.WriteOP(OP_SET, yx, newYX)
	return buf, yx, nil
}

func (table *symtable) compileNode(n *parser.Node) (code packet, yx uint32, err error) {
	var varIndex uint32

	switch n.Type {
	case parser.Natom:
		if n.Value.(string) == "nil" {
			buf := newpacket()
			yx = uint32(table.sp)
			buf.WriteOP(OP_SETK, yx, 0)
			table.sp++
			return buf, yx, nil
		}

		var ok bool
		varIndex, ok = table.get(n.Value.(string))
		if !ok {
			err = fmt.Errorf(errUndeclaredVariable, n)
			return
		}
	case parser.Naddr:
		varIndex = n.Value.(uint32)
	default:
		code, yx, err = table.compileCompound(n)
		if err != nil {
			return
		}
		varIndex = yx
	}
	return code, varIndex, nil
}

func (table *symtable) compileCompound(compound *parser.Node) (code packet, yx uint32, err error) {
	nodes := compound.C()
	if len(nodes) == 0 {
		return newpacket(), regA, nil
	}
	name, ok := nodes[0].Value.(string)
	if !ok {
		nodes[0].Dump(os.Stderr)
		panicf("invalid op: %v", nodes)
	}

	switch name {
	case "chain":
		code, yx, err = table.compileChainOp(compound)
	case "set", "move":
		code, yx, err = table.compileSetOp(nodes)
	case "ret", "yield":
		code, yx, err = table.compileRetOp(nodes)
	case "func", "safefunc":
		code, yx, err = table.compileLambdaOp(nodes)
	case "if":
		code, yx, err = table.compileIfOp(nodes)
	case "for":
		code, yx, err = table.compileWhileOp(nodes)
	case "continue", "break":
		code, yx, err = table.compileContinueBreakOp(nodes)
	case "call":
		code, yx, err = table.compileCallOp(nodes)
	case "map", "array":
		code, yx, err = table.compileMapArrayOp(nodes)
	case "or", "and":
		code, yx, err = table.compileAndOrOp(nodes)
	case "inc":
		code, yx, err = table.compileIncOp(nodes)
	default:
		if flatOpMapping[name] != 0 {
			return table.compileFlatOp(nodes)
		}
		panic(name)
	}
	table.im, table.ims = nil, nil
	return
}

func (table *symtable) compileChainOp(chain *parser.Node) (code packet, yx uint32, err error) {
	buf := newpacket()
	table.im = nil

	for _, a := range chain.C() {
		if a.Type != parser.Ncompound {
			continue
		}
		code, yx, err = table.compileCompound(a)
		if err != nil {
			return
		}
		buf.Write(code)
	}

	return buf, yx, err
}

func compileNode(n *parser.Node) (cls *Closure, err error) {
	table := newsymtable()
	for i, n := range CoreLibNames {
		table.sym[n] = uint16(i)
	}

	table.sp = uint16(len(table.sym))
	code, _, err := table.compileChainOp(n)
	if err != nil {
		return nil, err
	}

	code.WriteOP(OP_EOB, 0, 0)
	consts := make([]Value, len(table.consts)+1)
	for i, k := range table.consts {
		switch k.ty {
		case Tnumber:
			consts[i+1] = NewNumberValue(k.value.(float64))
		case Tstring:
			consts[i+1] = NewStringValue(k.value.(string))
		}
	}
	cls = NewClosure(code.data, consts, nil, 0)
	cls.lastenv = NewEnv(nil)
	cls.pos = code.pos
	cls.source = "<root>@" + code.source
	for _, name := range CoreLibNames {
		cls.lastenv.SPush(CoreLibs[name])
	}
	return cls, err
}

func LoadFile(path string) (*Closure, error) {
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
	return compileNode(n)
}

func LoadString(code string) (*Closure, error) {
	_, fn, _, _ := runtime.Caller(1)
	return loadStringName(code, fn)
}

func loadStringName(code, name string) (*Closure, error) {
	n, err := parser.Parse(bytes.NewReader([]byte(code)), name)
	if err != nil {
		return nil, err
	}
	return compileNode(n)
}
