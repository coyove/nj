package potatolang

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/coyove/potatolang/parser"
)

// symtable is responsible for recording the state of compilation
type symtable struct {
	// variable name lookup
	parent *symtable
	sym    map[string]uint16

	// has yield op
	y         bool
	envescape bool

	vp uint16

	continueNode []*parser.Node

	consts     []interface{}
	constMap   map[interface{}]uint16
	constStack map[interface{}]uint16

	reusableTmps map[uint16]bool
}

func newsymtable() *symtable {
	t := &symtable{
		sym:          make(map[string]uint16),
		consts:       make([]interface{}, 0),
		constMap:     make(map[interface{}]uint16),
		constStack:   make(map[interface{}]uint16),
		continueNode: make([]*parser.Node, 0),
		reusableTmps: make(map[uint16]bool),
	}
	return t
}

func (table *symtable) incrvp() {
	if table.vp >= 1<<10 {
		panic("too many variables (1024) in a single scope")
	}
	table.vp++
}

func (table *symtable) borrowTmp() uint16 {
	for tmp, ok := range table.reusableTmps {
		if ok {
			table.reusableTmps[tmp] = false
			return tmp
		}
	}
	table.incrvp()
	table.reusableTmps[table.vp-1] = false
	return table.vp - 1
}

func (table *symtable) returnTmp(v uint16) {
	if v == table.vp-1 {
		table.vp--
		return
	}
	if _, existed := table.reusableTmps[v]; existed {
		table.reusableTmps[v] = true
	}
}

func (table *symtable) get(varname string) (uint16, bool) {
	depth := uint16(0)

	if varname == "nil" {
		return table.getnil(), true
	}

	for table != nil {
		k, e := table.sym[varname]
		if e {
			if depth > 7 || (depth == 7 && k == 0x03ff) {
				panic("too many levels (8) to refer a variable, try simplifing your code")
			}
			return (depth << 10) | (uint16(k) & 0x03ff), true
		}

		depth++
		table = table.parent
	}

	return 0, false
}

func (table *symtable) put(varname string, addr uint16) {
	table.sym[varname] = addr
}

func (table *symtable) getnil() uint16 {
	return 0x3ff - 1
}

func (table *symtable) loadK(buf *packet, v interface{}) uint16 {
	addr, ok := table.constStack[v]
	if ok {
		return addr
	}

	addr = table.vp
	kaddr := func() uint16 {
		if i, ok := table.constMap[v]; ok {
			return i
		}

		table.consts = append(table.consts, v)
		if len(table.consts) > 1<<13-1 {
			panic("too many consts")
		}

		idx := uint16(len(table.consts) - 1)
		table.constMap[v] = idx
		return idx
	}()

	buf.WriteOP(OP_SETK, addr, kaddr)
	table.incrvp()
	table.constStack[v] = addr
	return addr
}

var flatOpMapping = map[string]byte{
	"+": OP_ADD, "-": OP_SUB, "*": OP_MUL, "/": OP_DIV, "%": OP_MOD,
	"<": OP_LESS, "<=": OP_LESS_EQ, "==": OP_EQ, "!=": OP_NEQ, "!": OP_NOT,
	"~": OP_BIT_NOT, "&": OP_BIT_AND, "|": OP_BIT_OR, "^": OP_BIT_XOR, "<<": OP_BIT_LSH, ">>": OP_BIT_RSH, ">>>": OP_BIT_URSH,
	"#": OP_POP, "store": OP_STORE, "load": OP_LOAD, "assert": OP_ASSERT, "slice": OP_SLICE, "typeof": OP_TYPEOF, "len": OP_LEN,
}

var flatOpMappingRev = map[byte]string{
	OP_ADD: "+", OP_SUB: "-", OP_MUL: "*", OP_DIV: "/", OP_MOD: "%",
	OP_LESS: "<", OP_LESS_EQ: "<=", OP_EQ: "==", OP_NEQ: "!=", OP_NOT: "!",
	OP_BIT_NOT: "~", OP_BIT_AND: "&", OP_BIT_OR: "|", OP_BIT_XOR: "^", OP_BIT_LSH: "<<", OP_BIT_RSH: ">>", OP_BIT_URSH: ">>>",
	OP_POP: "#", OP_STORE: "store", OP_LOAD: "load", OP_ASSERT: "assert", OP_SLICE: "slice", OP_TYPEOF: "typeof", OP_LEN: "len",
}

func (table *symtable) writeOpcode(buf *packet, op byte, n0, n1 *parser.Node) (err error) {
	getAddr := func(n *parser.Node) (uint16, error) {
		switch n.Type {
		case parser.Ncompound:
			code, addr, err := table.compileCompoundInto(n, true, 0, false)
			if err != nil {
				return 0, err
			}
			buf.Write(code)
			return addr, nil
		case parser.Natom:
			addr, ok := table.get(n.Value.(string))
			if !ok {
				return 0, fmt.Errorf(errUndeclaredVariable, n)
			}
			return addr, nil
		case parser.Nnumber, parser.Nstring:
			return table.loadK(buf, n.Value), nil
		case parser.Naddr:
			return n.Value.(uint16), nil
		default:
			panic(fmt.Errorf("unknown type: %d", n.Type))
		}
	}

	if n0 == nil {
		buf.WriteOP(op, 0, 0)
		return nil
	}

	n0a, err := getAddr(n0)
	if err != nil {
		return err
	}

	if n1 == nil {
		buf.WriteOP(op, n0a, 0)
		return nil
	}

	n1a, err := getAddr(n1)
	if err != nil {
		return err
	}

	if op == OP_SET && n0a == n1a {
		return nil
	}
	buf.WriteOP(op, n0a, n1a)
	return nil
}

func (table *symtable) compileCompoundInto(compound *parser.Node, newVar bool, existedVar uint16, tmp bool) (code packet, yx uint16, err error) {
	buf := newpacket()

	var newYX uint16
	code, newYX, err = table.compileCompound(compound)
	if err != nil {
		return
	}

	buf.Write(code)
	if newVar {
		if tmp {
			yx = table.borrowTmp()
		} else {
			yx = table.vp
			table.incrvp()
		}
	} else {
		yx = existedVar
	}

	buf.WriteOP(OP_SET, yx, newYX)
	return buf, yx, nil
}

func (table *symtable) compileNode(n *parser.Node) (code packet, yx uint16, err error) {
	var varIndex uint16

	switch n.Type {
	case parser.Natom:
		if n.Value.(string) == "nil" {
			buf := newpacket()
			yx = table.vp
			buf.WriteOP(OP_SET, yx, 0)
			table.incrvp()
			return buf, yx, nil
		}

		var ok bool
		varIndex, ok = table.get(n.Value.(string))
		if !ok {
			err = fmt.Errorf(errUndeclaredVariable, n)
			return
		}
	case parser.Naddr:
		varIndex = n.Value.(uint16)
	default:
		code, yx, err = table.compileCompound(n)
		if err != nil {
			return
		}
		varIndex = yx
	}
	return code, varIndex, nil
}

func (table *symtable) compileCompound(compound *parser.Node) (code packet, yx uint16, err error) {
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
		if strings.Contains(name, "func") {
			code, yx, err = table.compileLambdaOp(nodes)
		} else {
			if _, ok := flatOpMapping[name]; ok {
				return table.compileFlatOp(nodes)
			}
			panic(name)
		}
	}
	return
}

func (table *symtable) compileChainOp(chain *parser.Node) (code packet, yx uint16, err error) {
	buf := newpacket()

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
	defer func() {
		if r := recover(); r != nil {
			cls = nil
			err = fmt.Errorf("recovered panic: %v, from: %s", r, debug.Stack())
		}
	}()

	table := newsymtable()
	for i, n := range CoreLibNames {
		table.sym[n] = uint16(i)
	}

	table.vp = uint16(len(table.sym))
	code, _, err := table.compileChainOp(n)
	if err != nil {
		return nil, err
	}

	code.WriteOP(OP_EOB, 0, 0)
	consts := make([]Value, len(table.consts))
	for i, k := range table.consts {
		switch k := k.(type) {
		case float64:
			consts[i] = NewNumberValue(k)
		case string:
			consts[i] = NewStringValue(k)
		}
	}
	cls = NewClosure(code.data, consts, nil, 0)
	cls.lastenv = NewEnv(nil, nil)
	cls.pos = code.pos
	cls.source = "root" + cls.String() + "@" + code.source
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
	//n.Dump(os.Stderr)
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
	// n.Dump(os.Stderr)
	return compileNode(n)
}
