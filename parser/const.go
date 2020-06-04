package parser

import (
	"unsafe"
)

type Symbol string

var (
	NUM = interfaceType(1.0)
	STR = interfaceType("")
	SYM = interfaceType(Symbol(""))
	CPL = interfaceType([]*Node{})
	ADR = interfaceType(uint16(1))

	chainNode = Nod(ABegin)
	breakNode = Cpl(Nod(ABreak))
	popvNode  = Cpl(Nod(APopV))
	nilNode   = Nod(ANil)
	zeroNode  = Num(0)
	oneNode   = Num(1)
	emptyNode = Cpl()
)

func interfaceType(a interface{}) uintptr {
	return (*(*[2]uintptr)(unsafe.Pointer(&a)))[0]
}

const (
	ANop         Symbol = "nop"
	ADoBlock     Symbol = "do"
	AConcat      Symbol = "con"
	ANil         Symbol = "nil"
	ASet         Symbol = "set"
	AInc         Symbol = "inc"
	AMove        Symbol = "mov"
	AIf          Symbol = "if"
	AFor         Symbol = "for"
	APatchVararg Symbol = "pvag"
	AFunc        Symbol = "fun"
	ABreak       Symbol = "brk"
	AContinue    Symbol = "cont"
	ABegin       Symbol = "prog"
	ALoad        Symbol = "load"
	AStore       Symbol = "stor"
	ACall        Symbol = "call"
	AReturn      Symbol = "ret"
	AYield       Symbol = "yld"
	AHash        Symbol = "hash"
	AHashArray   Symbol = "harr"
	AArray       Symbol = "arr"
	AAdd         Symbol = "add"
	ASub         Symbol = "sub"
	AMul         Symbol = "mul"
	ADiv         Symbol = "div"
	AMod         Symbol = "mod"
	APow         Symbol = "pow"
	AEq          Symbol = "eq"
	ANeq         Symbol = "neq"
	AAnd         Symbol = "and"
	AOr          Symbol = "or"
	ANot         Symbol = "not"
	ALess        Symbol = "lt"
	ALessEq      Symbol = "le"
	ALen         Symbol = "len"
	ARetAddr     Symbol = "reta"
	APopV        Symbol = "popv"
	ALabel       Symbol = "lbl"
	AGoto        Symbol = "goto"
)

func __chain(args ...interface{}) *Node { return Cpl(append([]interface{}{ABegin}, args...)...) }

func __do(args ...interface{}) *Node { return Cpl(append([]interface{}{ADoBlock}, args...)...) }

func __move(dest, src interface{}) *Node { return Cpl(AMove, dest, src) }

func __set(dest, src interface{}) *Node {
	if n, _ := dest.(*Node); n != nil && n.Type() != SYM {
		panic(&Error{Pos: n.Position, Message: "invalid assignments", Token: n.String()})
	}
	return Cpl(ASet, dest, src)
}

func __less(lhs, rhs interface{}) *Node { return Cpl(ALess, lhs, rhs) }

func __lessEq(lhs, rhs interface{}) *Node { return Cpl(ALessEq, lhs, rhs) }

func __inc(subject, step interface{}) *Node { return Cpl(AInc, subject, step) }

func __load(subject, key interface{}) *Node { return Cpl(ALoad, subject, key) }

func __call(cls, args interface{}) *Node { return Cpl(ACall, cls, args) }

func __store(subject, key, value interface{}) *Node { return Cpl(AStore, subject, value, key) }

func __if(cond, truebody, falsebody interface{}) *Node { return Cpl(AIf, cond, truebody, falsebody) }

func __loop(body interface{}) *Node { return Cpl(AFor, body) }

func __func(paramlist, body interface{}) *Node { return Cpl(AFunc, emptyNode, paramlist, body) }
