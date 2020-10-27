package parser

import (
	"math/rand"
	"strconv"
	"strings"
)

const (
	Float = iota + 1
	Int
	String
	Symbol
	Complex
	Address
)

var (
	breakNode     = NewComplex(NewSymbol(ABreak))
	popvNode      = NewComplex(NewSymbol(APopV))
	popvClearNode = NewComplex(NewSymbol(APopVClear))
	zeroNode      = NewNumberFromInt(0)
	oneNode       = NewNumberFromInt(1)
	emptyNode     = NewComplex()
)

const (
	ANop       = "nop"
	ADoBlock   = "do"
	AConcat    = "concat"
	ANil       = "nil"
	ASet       = "set"
	AInc       = "incr"
	AMove      = "move"
	AIf        = "if"
	AFor       = "loop"
	AFunc      = "function"
	ABreak     = "break"
	ABegin     = "prog"
	ALoad      = "load"
	AStore     = "store"
	ASlice     = "slice"
	ACall      = "call"
	ACallMap   = "callmap"
	ATailCall  = "tailcall"
	AReturn    = "return"
	AYield     = "yield"
	AAdd       = "add"
	ASub       = "sub"
	AMul       = "mul"
	ADiv       = "div"
	AMod       = "mod"
	APow       = "pow"
	AEq        = "eq"
	ANeq       = "neq"
	AAnd       = "and"
	AOr        = "or"
	ANot       = "not"
	ALess      = "lt"
	ALessEq    = "le"
	ALen       = "len"
	ARetAddr   = "retaddr"
	APopV      = "popv"
	APopVClear = "clearv"
	APopVAll   = "popallv"
	APopVAllA  = "popallva"
	ALabel     = "label"
	AGoto      = "goto"
	AJSON      = "map"
)

func __chain(args ...Node) Node {
	return NewComplex(append([]Node{NewSymbol(ABegin)}, args...)...)
}

func __do(args ...Node) Node {
	return NewComplex(append([]Node{NewSymbol(ADoBlock)}, args...)...)
}

func RemoveDDD(dest Node) Node {
	sym := dest.strSym
	if sym != "..." {
		sym = strings.TrimLeft(sym, ".")
		dest.strSym = sym
	}
	return dest
}

func __move(dest, src Node) Node { return NewComplex(NewSymbol(AMove), RemoveDDD(dest), src) }

func __set(dest, src Node) Node { return NewComplex(NewSymbol(ASet), RemoveDDD(dest), src) }

func __less(lhs, rhs Node) Node { return NewComplex(NewSymbol(ALess), lhs, rhs) }

func __lessEq(lhs, rhs Node) Node { return NewComplex(NewSymbol(ALessEq), lhs, rhs) }

func __inc(subject, step Node) Node { return NewComplex(NewSymbol(AInc), subject, step) }

func __load(subject, key Node) Node { return NewComplex(NewSymbol(ALoad), subject, key) }

func __store(subject, key, value Node) Node { return NewComplex(NewSymbol(AStore), subject, value, key) }

func __if(cond, truebody, falsebody Node) Node {
	return NewComplex(NewSymbol(AIf), cond, truebody, falsebody)
}

func __loop(body Node) Node { return NewComplex(NewSymbol(AFor), body) }

func __func(name, paramlist, body Node) Node {
	return NewComplex(NewSymbol(AFunc), name, paramlist, body)
}

func __call(cls, args Node) Node { return NewComplex(NewSymbol(ACall), cls, args) }

func __popvAll(i int, k Node) Node {
	if i == 0 {
		return __chain(k, NewComplex(NewSymbol(APopVAllA)))
	}
	return NewComplex(NewSymbol(APopVAll))
}

func __findTailCall(stats []Node) {
	for len(stats) > 0 {
		x := stats[len(stats)-1]
		c := x.Nodes
		if len(c) == 3 && c[0].SymbolValue() == ACall {
			c[0].strSym = ATailCall
			return
		}

		if len(c) > 0 {
			if c[0].SymbolValue() == (ABegin) {
				__findTailCall(c)
				return
			}

			switch c[0].SymbolValue() {
			case APopV, APopVClear, APopVAll, APopVAllA:
				stats = stats[:len(stats)-1]
				continue
			}
		}
		return
	}
}

func randomVarname() Node {
	return NewSymbol("v" + strconv.FormatInt(rand.Int63(), 10))
}
