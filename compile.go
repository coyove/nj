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

	consts   []interface{}
	constMap map[interface{}]uint16

	reusableTmps map[uint16]bool
}

func newsymtable() *symtable {
	t := &symtable{
		sym:          make(map[string]uint16),
		constMap:     make(map[interface{}]uint16),
		reusableTmps: make(map[uint16]bool),
	}
	return t
}

func (table *symtable) borrowAddress() uint16 {
	for tmp, ok := range table.reusableTmps {
		if ok {
			table.reusableTmps[tmp] = false
			return tmp
		}
	}
	if table.vp > 1000 { //1<<10 {
		panic("Code too complex, may be there are too many variables (1000) in a single scope")
	}
	table.reusableTmps[table.vp] = false
	table.vp++
	return table.vp - 1
}

func (table *symtable) returnAddress(v uint16) {
	//	log.Println("$$", table.reusableTmps, v, table.vp)
	//if v == table.vp-1 {
	//	table.vp--
	//	return
	//}
	if _, existed := table.reusableTmps[v]; existed {
		table.reusableTmps[v] = true
	}
}

func (table *symtable) get(varname string) (uint16, bool) {
	depth := uint16(0)

	switch varname {
	case "nil":
		return regNil, true
	case "true":
		return table.loadK(nil, 1.0), true
	case "false":
		return table.loadK(nil, 0.0), true
	}

	for table != nil {
		k, e := table.sym[varname]
		if e {
			if depth > 6 || (depth == 6 && k > 1000) {
				panic("too many levels (7) to refer a variable, try simplifing your Code")
			}
			return (depth << 10) | (uint16(k) & 0x03ff), true
		}

		depth++
		table = table.parent
	}

	return 0, false
}

func (table *symtable) put(varname string, addr uint16) {
	if addr == regA {
		panic("debug")
	}
	table.sym[varname] = addr
}

func (table *symtable) loadK(buf *packet, v interface{}) uint16 {
	kaddr := func() uint16 {
		if i, ok := table.constMap[v]; ok {
			return i
		}

		table.consts = append(table.consts, v)
		if len(table.consts) > 1<<10-1 {
			panic("too many ConstTable")
		}

		idx := uint16(len(table.consts) - 1)
		table.constMap[v] = idx
		return idx
	}()

	return 0x7<<10 | kaddr
}

var flatOpMapping = map[string]_Opcode{
	"+": OpAdd, "-": OpSub, "*": OpMul, "/": OpDiv, "%": OpMod,
	"<": OpLess, "<=": OpLessEq, "==": OpEq, "!=": OpNeq, "!": OpNot,
	"~": OpBitNot, "&": OpBitAnd, "|": OpBitOr, "^": OpBitXor, "<<": OpBitLsh, ">>": OpBitRsh, ">>>": OpBitURsh, "#": OpPop,
	"store": OpStore, "load": OpLoad, "assert": OpAssert, "slice": OpSlice, "typeof": OpTypeof, "len": OpLen, "foreach": OpForeach,
	"addressof": OpAddressOf,
}

func (table *symtable) writeOpcode(buf *packet, op _Opcode, n0, n1 *parser.Node) (err error) {
	tmp := []uint16{}
	getAddr := func(n *parser.Node) (uint16, error) {
		switch n.Type {
		case parser.Ncompound:
			code, addr, err := table.compileCompoundInto(n, true, 0)
			if err != nil {
				return 0, err
			}
			buf.Write(code)
			tmp = append(tmp, addr)
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

	defer func() {
		for _, tmp := range tmp {
			table.returnAddress(tmp)
		}
	}()

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

	if op == OpSet && n0a == n1a {
		return nil
	}

	buf.WriteOP(op, n0a, n1a)
	return nil
}

func (table *symtable) compileCompoundInto(compound *parser.Node, newVar bool, existedVar uint16) (code packet, yx uint16, err error) {
	buf := newpacket()

	var newYX uint16
	code, newYX, err = table.compileCompound(compound)
	if err != nil {
		return
	}

	buf.Write(code)
	if newVar {
		yx = table.borrowAddress()
	} else {
		yx = existedVar
	}

	buf.WriteOP(OpSet, yx, newYX)
	return buf, yx, nil
}

func (table *symtable) compileNode(n *parser.Node) (code packet, yx uint16, err error) {
	var varIndex uint16

	switch n.Type {
	case parser.Natom:
		var ok bool
		varIndex, ok = table.get(n.Value.(string))
		if !ok {
			err = fmt.Errorf(errUndeclaredVariable, n)
			return
		}
	case parser.Naddr:
		varIndex = n.Value.(uint16)
	case parser.Nstring, parser.Nnumber:
		varIndex = table.loadK(nil, n.Value)
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

	//log.Println(table.vp)

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

	coreStack := NewEnv(nil)
	for n, v := range CoreLibs {
		table.sym[n] = uint16(coreStack.LocalSize())
		coreStack.LocalPush(v)
	}

	table.vp = uint16(len(table.sym))
	code, _, err := table.compileChainOp(n)
	if err != nil {
		return nil, err
	}

	code.WriteOP(OpEOB, 0, 0)
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
	cls.lastenv = NewEnv(nil)
	cls.Pos = code.pos
	cls.source = "root" + cls.String() + "@" + code.source
	cls.lastenv.stack = coreStack.stack
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
	//panic(10)
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
