package parser

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync/atomic"

	"github.com/coyove/nj/bas"
)

const (
	INVALID = iota
	FLOAT
	INT
	STR
	SYM
	NODES
	ADDR
	JSON
)

var (
	breakNode = Nodes((SBreak))
	zero      = Int(0)
	one       = Int(1)
	emptyNode = Nodes()
)

var (
	Sa        = staticSym("$a")
	SDoBlock  = staticSym("doblock")
	SNil      = staticSym("nil")
	SSet      = staticSym("set")
	SInc      = staticSym("inc")
	SMove     = staticSym("move")
	SIf       = staticSym("if")
	SWhile    = staticSym("loop")
	SFunc     = staticSym("func")
	SBreak    = staticSym("break")
	SContinue = staticSym("continue")
	SBegin    = staticSym("begin")
	SLoad     = staticSym("load")
	SStore    = staticSym("store")
	SSlice    = staticSym("slice")
	SArray    = staticSym("array")
	SObject   = staticSym("object")
	SCall     = staticSym("call")
	STailCall = staticSym("tailcall")
	STryCall  = staticSym("trycall")
	SReturn   = staticSym("return")
	SLen      = staticSym("len")
	SNext     = staticSym("next")
	SAdd      = staticSym("add")
	SSub      = staticSym("sub")
	SMul      = staticSym("mul")
	SDiv      = staticSym("div")
	SIDiv     = staticSym("idiv")
	SMod      = staticSym("mod")
	SBitAnd   = staticSym("bitand")
	SBitOr    = staticSym("bitor")
	SBitXor   = staticSym("bitxor")
	SBitLsh   = staticSym("bitlsh")
	SBitRsh   = staticSym("bitrsh")
	SBitURsh  = staticSym("bitursh")
	SEq       = staticSym("eq")
	SNeq      = staticSym("neq")
	SAnd      = staticSym("and")
	SOr       = staticSym("or")
	SNot      = staticSym("not")
	SLess     = staticSym("less")
	SLessEq   = staticSym("lesseq")
	SFreeAddr = staticSym("freeaddr")
	SLabel    = staticSym("label")
	SGoto     = staticSym("goto")
	SUnpack   = staticSym("unpack")
	SIs       = staticSym("isproto")
)

func __chain(args ...Node) Node { return Nodes(append([]Node{SBegin}, args...)...) }

func __do(args ...Node) Node { return Nodes(append([]Node{SDoBlock}, args...)...) }

func __move(dest, src Node) Node { return Nodes(SMove, dest, src) }

func __set(dest, src Node) Node { return Nodes(SSet, dest, src) }

func __less(lhs, rhs Node) Node { return Nodes(SLess, lhs, rhs) }

func __lessEq(lhs, rhs Node) Node { return Nodes(SLessEq, lhs, rhs) }

func __inc(subject, step Node) Node { return Nodes((SInc), subject, step) }

func __load(subject, key Node) Node { return Nodes((SLoad), subject, key) }

func __store(subject, key, value Node) Node { return Nodes(SStore, subject, key, value) }

func __if(cond, t, f Node) Node { return Nodes((SIf), cond, t, f) }

func __loop(cont Node, body ...Node) Node { return Nodes(SWhile, __chain(body...), cont) }

func __goto(label Node) Node { return Nodes(SGoto, label) }

func __label(name Node) Node { return Nodes(SLabel, name) }

func __ret(v Node) Node { return Nodes(SReturn, v) }

func __func(name Token, paramList Node, stats Node) Node {
	__findTailCall(stats.Nodes())
	funcname := Sym(name)
	return __chain(
		__set(funcname, SNil).At(name),
		__move(funcname, Nodes(SFunc, funcname, paramList, stats).At(name)).At(name),
	)
}

func __method(name Token, paramList Node, stats Node) Node {
	__findTailCall(stats.Nodes())
	return Nodes(SFunc, Sym(name), paramList, stats).At(name)
}

func __lambda(name Token, pp Node, stats Node) Node {
	nodes := stats.Nodes()
	if len(nodes) > 1 && nodes[0].Value == SBegin.Value {
		nodes[len(nodes)-1] = Nodes(SReturn, nodes[len(nodes)-1])
	}
	return __func(name, pp, stats)
}

func __markupFuncName(recv, name Token) Token {
	name.Str = recv.Str + "." + name.Str
	return name
}

var lambdaIndex int64

func __markupLambdaName(lambda Token) Token {
	lambda.Str = fmt.Sprintf("<lambda-%d-%d>", lambda.Pos.Line, atomic.AddInt64(&lambdaIndex, 1))
	return lambda
}

func __call(cls, args Node) Node { return Nodes(SCall, cls, args) }

func __tryCall(cls, args Node) Node { return Nodes(STryCall, cls, args) }

func __findTailCall(stats []Node) {
	if len(stats) > 0 {
		x := stats[len(stats)-1]
		c := x.Nodes()
		if len(c) == 3 && c[0].Value == SCall.Value {
			old := c[0].SymLine
			c[0] = STailCall
			c[0].SymLine = old
			return
		}

		if len(c) > 0 {
			if c[0].Value == SBegin.Value {
				__findTailCall(c)
				return
			}

		}
	}
}

func __local(dest, src []Node, pos Token) Node {
	m, n := len(dest), len(src)
	for i, count := 0, m-n; i < count; i++ {
		src = append(src, SNil)
	}
	res := __chain()
	for i, v := range dest {
		res = res.append(__set(v, src[i]).At(pos))
	}
	return res
}

func __moveMulti(nodes, src []Node, pos Token) Node {
	m, n := len(nodes), len(src)
	for i, count := 0, m-n; i < count; i++ {
		src = append(src, SNil)
	}

	res := __chain()
	if head := nodes[0]; len(nodes) == 1 {
		res = head.moveLoadStore(src[0]).At(pos)
	} else {
		// a0, ..., an = b0, ..., bn
		names, retaddr := []Node{}, Nodes((SFreeAddr))
		for i := range nodes {
			names = append(names, randomVarname())
			retaddr = retaddr.append(names[i])
			res = res.append(__set(names[i], src[i]).At(pos))
		}
		for i, v := range nodes {
			res = res.append(v.moveLoadStore(names[i]).At(pos))
		}
		res = res.append(retaddr)
	}
	return res
}

func __dotdotdot(expr Node) Node {
	expr.Nodes()[len(expr.Nodes())-1] = Nodes((SUnpack), expr.Nodes()[len(expr.Nodes())-1])
	return expr
}

func __forIn(key, value Token, expr, body Node, pos Token) Node {
	k, v, subject, kv := Sym(key), Sym(value), randomVarname(), randomVarname()
	moveNext := __chain(
		__move(kv, Nodes(SNext, subject, kv).At(pos)).At(pos),
		__move(k, __load(kv, zero).At(pos)).At(pos),
		__move(v, __load(kv, one).At(pos)).At(pos),
	)
	return __do(
		__set(subject, expr).At(pos),
		__set(k, SNil).At(pos),
		__set(v, SNil).At(pos),
		__set(kv, SNil).At(pos),
		__loop(
			one,
			moveNext,
			__if(Nodes(SEq, kv, SNil).At(pos), breakNode, body).At(pos),
		),
	)
}

func (lex *Lexer) __arrayBuild(list, arg Node) Node {
	if lex.scanner.jsonMode {
		if list.Valid() {
			list.simpleJSON(lex).Native().Append(arg.simpleJSON(lex))
			return list
		}
		return jsonValue(bas.Array(arg.simpleJSON(lex)))
	}
	if list.Valid() {
		return list.append(arg)
	}
	return Nodes(arg)
}

func (lex *Lexer) __objectBuild(list, k, v Node) Node {
	if lex.scanner.jsonMode {
		if list.Valid() {
			list.simpleJSON(lex).Object().Set(k.simpleJSON(lex), v.simpleJSON(lex))
			return list
		}
		o := bas.NewObject(0)
		o.Set(k.simpleJSON(lex), v.simpleJSON(lex))
		return jsonValue(o.ToValue())
	}
	if list.Valid() {
		return list.append(k).append(v)
	}
	return Nodes(k, v)
}

func (lex *Lexer) __array(tok Token, args Node) Node {
	if lex.scanner.jsonMode {
		if args == emptyNode {
			return Node{NodeType: JSON, Value: bas.Array()}
		}
		return args
	}
	return Nodes(SArray, args).At(tok)
}

func (lex *Lexer) __object(tok Token, args Node) Node {
	if lex.scanner.jsonMode {
		if args == emptyNode {
			return jsonValue(bas.NewObject(0).ToValue())
		}
		return args
	}
	return Nodes(SObject, args).At(tok)
}

func randomVarname() Node {
	return staticSym("v" + strconv.FormatInt(rand.Int63(), 10)[:6])
}
