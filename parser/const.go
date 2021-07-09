package parser

import (
	"math/rand"
	"strconv"
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
	breakNode = NewComplex(NewSymbol(ABreak))
	zeroNode  = NewNumberFromInt(0)
	oneNode   = NewNumberFromInt(1)
	emptyNode = NewComplex()
)

const (
	ADoBlock  = "do"
	ANil      = "nil"
	ASet      = "set"
	AInc      = "incr"
	AMove     = "move"
	AIf       = "if"
	AFor      = "loop"
	AFunc     = "function"
	ABreak    = "break"
	ABegin    = "prog"
	ALoad     = "load"
	AStore    = "store"
	AArray    = "array"
	AArrayMap = "map"
	ACall     = "call"
	ATailCall = "tailcall"
	AReturn   = "return"
	AAdd      = "add"
	ASub      = "sub"
	AMul      = "mul"
	ADiv      = "div"
	AIDiv     = "idiv"
	AMod      = "mod"
	ABitAnd   = "bitand"
	ABitOr    = "bitor"
	ABitXor   = "bitxor"
	ABitNot   = "bitnot"
	ABitLsh   = "bitlsh"
	ABitRsh   = "bitrsh"
	ABitURsh  = "bitursh"
	AEq       = "eq"
	ANeq      = "neq"
	AAnd      = "and"
	AOr       = "or"
	ANot      = "not"
	ALess     = "lt"
	ALessEq   = "le"
	AFreeAddr = "freeaddr"
	ALabel    = "label"
	AGoto     = "goto"
	AUnpack   = "unpack"
)

func __chain(args ...Node) Node { return NewComplex(append([]Node{NewSymbol(ABegin)}, args...)...) }

func __do(args ...Node) Node { return NewComplex(append([]Node{NewSymbol(ADoBlock)}, args...)...) }

func __move(dest, src Node) Node { return NewComplex(NewSymbol(AMove), dest, src) }

func __set(dest, src Node) Node { return NewComplex(NewSymbol(ASet), dest, src) }

func __less(lhs, rhs Node) Node { return NewComplex(NewSymbol(ALess), lhs, rhs) }

func __lessEq(lhs, rhs Node) Node { return NewComplex(NewSymbol(ALessEq), lhs, rhs) }

func __inc(subject, step Node) Node { return NewComplex(NewSymbol(AInc), subject, step) }

func __load(subject, key Node) Node { return NewComplex(NewSymbol(ALoad), subject, key) }

func __store(subject, key, value Node) Node {
	return NewComplex(NewSymbol(AStore), subject, value, key)
}

func __if(cond, t, f Node) Node { return NewComplex(NewSymbol(AIf), cond, t, f) }

func __loop(body Node) Node { return NewComplex(NewSymbol(AFor), body) }

func __func(name Token, paramList Node, doc string, stats Node) Node {
	__findTailCall(stats.Nodes)
	funcname := NewSymbolFromToken(name)
	p := name.Pos
	return __chain(
		__set(funcname, NewSymbol(ANil)).SetPos(p),
		__move(funcname,
			NewComplex(NewSymbol(AFunc), funcname, paramList, stats, NewString(doc)).SetPos(p)).SetPos(p),
	)
}

func __markupFuncName(recv, name Token) Token {
	name.Str = recv.Str + "." + name.Str
	return name
}

func __call(cls, args Node) Node {
	return NewComplex(NewSymbol(ACall), cls, args)
}

func __callPatch(original, self Node) Node {
	if original.Nodes[0].SymbolValue() == ACall {
		original.Nodes[2].Nodes = append([]Node{self}, original.Nodes[2].Nodes...)
	} else {
		n := original.Nodes[2].Nodes
		n = append([]Node{zeroNode, self}, n...)
		for i := 2; i < len(n); i += 2 {
			if n[i].Type != Int {
				break
			}
			n[i] = NewNumberFromInt(n[i].IntValue() + 1)
		}
		original.Nodes[2].Nodes = n
	}
	return original
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

		}
		return
	}
}

func __local(dest, src []Node, pos Position) Node {
	m, n := len(dest), len(src)
	for i, count := 0, m-n; i < count; i++ {
		src = append(src, NewSymbol(ANil))
	}
	res := __chain()
	for i, v := range dest {
		res = res.append(__set(v, src[i]).SetPos(pos))
	}
	return res
}

func __moveMulti(nodes, src []Node, pos Position) Node {
	m, n := len(nodes), len(src)
	for i, count := 0, m-n; i < count; i++ {
		src = append(src, NewSymbol(ANil))
	}

	res := __chain()
	if head := nodes[0]; len(nodes) == 1 {
		res = head.moveLoadStore(__move, src[0]).SetPos(pos)
	} else {
		// a0, ..., an = b0, ..., bn
		names, retaddr := []Node{}, NewComplex(NewSymbol(AFreeAddr))
		for i := range nodes {
			names = append(names, randomVarname())
			retaddr = retaddr.append(names[i])
			res = res.append(__set(names[i], src[i]).SetPos(pos))
		}
		for i, v := range nodes {
			res = res.append(v.moveLoadStore(__move, names[i]).SetPos(pos))
		}
		res = res.append(retaddr)
	}
	return res
}

func __dotdotdot(expr Node) Node {
	expr.Nodes[len(expr.Nodes)-1] = NewComplex(NewSymbol(AUnpack), expr.Nodes[len(expr.Nodes)-1])
	return expr
}

func randomVarname() Node {
	return NewSymbol("v" + strconv.FormatInt(rand.Int63(), 10)[:6])
}
