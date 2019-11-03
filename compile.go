package potatolang

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/coyove/potatolang/parser"
)

// symtable is responsible for recording the state of compilation
type symtable struct {
	// variable name lookup
	parent *symtable
	sym    map[parser.Atom]uint16

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
		sym:          make(map[parser.Atom]uint16),
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

func (table *symtable) get(varname parser.Atom) (uint16, bool) {
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

func (table *symtable) put(varname parser.Atom, addr uint16) {
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

var flatOpMapping = map[parser.Atom]_Opcode{
	parser.AAdd:     OpAdd,
	parser.ASub:     OpSub,
	parser.AMul:     OpMul,
	parser.ADiv:     OpDiv,
	parser.AMod:     OpMod,
	parser.ALess:    OpLess,
	parser.ALessEq:  OpLessEq,
	parser.AEq:      OpEq,
	parser.ANeq:     OpNeq,
	parser.ANot:     OpNot,
	parser.ABitAnd:  OpBitAnd,
	parser.ABitOr:   OpBitOr,
	parser.ABitXor:  OpBitXor,
	parser.ABitLsh:  OpBitLsh,
	parser.ABitRsh:  OpBitRsh,
	parser.ABitURsh: OpBitURsh,
	parser.AStore:   OpStore,
	parser.ALoad:    OpLoad,
	parser.AAssert:  OpAssert,
	parser.ASlice:   OpSlice,
	parser.ATypeOf:  OpTypeof,
	parser.ALen:     OpLen,
	parser.AForeach: OpForeach,
	parser.AAddrOf:  OpAddressOf,
	parser.AInc:     OpInc,
}

func (table *symtable) writeOpcode(buf *packet, op _Opcode, n0, n1 *parser.Node) (err error) {
	tmp := []uint16{}
	getAddr := func(n *parser.Node) (uint16, error) {
		switch n.Type() {
		case parser.Ncompound:
			code, addr, err := table.compileCompoundInto(n, true, 0)
			if err != nil {
				return 0, err
			}
			buf.Write(code)
			tmp = append(tmp, addr)
			return addr, nil
		case parser.Natom:
			addr, ok := table.get(n.Value.(parser.Atom))
			if !ok {
				return 0, fmt.Errorf(errUndeclaredVariable, n)
			}
			return addr, nil
		case parser.Nnumber, parser.Nstring:
			return table.loadK(buf, n.Value), nil
		case parser.Naddr:
			return n.Value.(uint16), nil
		default:
			panic(fmt.Errorf("unknown type: %d", n.Type()))
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

	switch n.Type() {
	case parser.Natom:
		var ok bool
		varIndex, ok = table.get(n.Value.(parser.Atom))
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
	name, ok := nodes[0].Value.(parser.Atom)
	if !ok {
		nodes[0].Dump(os.Stderr)
		panicf("invalid op: %v", nodes)
	}

	switch name {
	case parser.AChain:
		code, yx, err = table.compileChainOp(compound)
	case parser.ASet, parser.AMove:
		code, yx, err = table.compileSetOp(nodes)
	case parser.AReturn:
		code, yx, err = table.writeOpcode3(OpRet, nodes)
	case parser.AYield:
		table.y = true
		code, yx, err = table.writeOpcode3(OpYield, nodes)
	case parser.AIf:
		code, yx, err = table.compileIfOp(nodes)
	case parser.AFor:
		code, yx, err = table.compileWhileOp(nodes)
	case parser.AContinue, parser.ABreak:
		code, yx, err = table.compileContinueBreakOp(nodes)
	case parser.ACall:
		code, yx, err = table.compileCallOp(nodes)
	case parser.AMap, parser.AArray:
		code, yx, err = table.compileMapArrayOp(nodes)
	case parser.AOr, parser.AAnd:
		code, yx, err = table.compileAndOrOp(nodes)
	case parser.AFunc:
		code, yx, err = table.compileLambdaOp(nodes)
	default:
		if _, ok := flatOpMapping[name]; ok {
			return table.compileFlatOp(nodes)
		}
		log.Println(nodes)
		panic(name)
	}
	return
}

func (table *symtable) compileChainOp(chain *parser.Node) (code packet, yx uint16, err error) {
	buf := newpacket()

	for _, a := range chain.C() {
		if a.Type() != parser.Ncompound {
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
		table.sym[parser.Atom(n)] = uint16(coreStack.LocalSize())
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
			consts[i] = NewStringValueString(k)
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
