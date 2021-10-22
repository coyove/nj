package parser

import (
	"math/rand"
	"strconv"
	"unsafe"
)

const (
	FLOAT = iota + 1
	INT
	STR
	SYM
	NODES
	ADDR
	INVALID
)

var (
	intNode   = unsafe.Pointer(new(int))
	floatNode = unsafe.Pointer(new(int))
	breakNode = Nodes((SBreak))
	zero      = Int(0)
	one       = Int(1)
	emptyNode = Nodes()
)

var (
	ADoBlock, SDoBlock   = "do", staticSym("do")
	ANil, SNil           = "nil", staticSym("nil")
	ASet, SSet           = "set", staticSym("set")
	AInc, SInc           = "incr", staticSym("incr")
	AMove, SMove         = "move", staticSym("move")
	AIf, SIf             = "if", staticSym("if")
	AFor, SFor           = "loop", staticSym("loop")
	AFunc, SFunc         = "function", staticSym("function")
	ABreak, SBreak       = "break", staticSym("break")
	ABegin, SBegin       = "prog", staticSym("prog")
	ALoad, SLoad         = "load", staticSym("load")
	AStore, SStore       = "store", staticSym("store")
	AArray, SArray       = "array", staticSym("array")
	AArrayMap, SArrayMap = "map", staticSym("map")
	ACall, SCall         = "call", staticSym("call")
	ATailCall, STailCall = "tailcall", staticSym("tailcall")
	AReturn, SReturn     = "return", staticSym("return")
	AAdd, SAdd           = "add", staticSym("add")
	ASub, SSub           = "sub", staticSym("sub")
	AMul, SMul           = "mul", staticSym("mul")
	ADiv, SDiv           = "div", staticSym("div")
	AIDiv, SIDiv         = "idiv", staticSym("idiv")
	AMod, SMod           = "mod", staticSym("mod")
	ABitAnd, SBitAnd     = "bitand", staticSym("bitand")
	ABitOr, SBitOr       = "bitor", staticSym("bitor")
	ABitXor, SBitXor     = "bitxor", staticSym("bitxor")
	ABitNot, SBitNot     = "bitnot", staticSym("bitnot")
	ABitLsh, SBitLsh     = "bitlsh", staticSym("bitlsh")
	ABitRsh, SBitRsh     = "bitrsh", staticSym("bitrsh")
	ABitURsh, SBitURsh   = "bitursh", staticSym("bitursh")
	AEq, SEq             = "eq", staticSym("eq")
	ANeq, SNeq           = "neq", staticSym("neq")
	AAnd, SAnd           = "and", staticSym("and")
	AOr, SOr             = "or", staticSym("or")
	ANot, SNot           = "not", staticSym("not")
	ALess, SLess         = "lt", staticSym("lt")
	ALessEq, SLessEq     = "le", staticSym("le")
	AFreeAddr, SFreeAddr = "freeaddr", staticSym("freeaddr")
	ALabel, SLabel       = "label", staticSym("label")
	AGoto, SGoto         = "goto", staticSym("goto")
	AUnpack, SUnpack     = "unpack", staticSym("unpack")
)

func __chain(args ...Node) Node { return Nodes(append([]Node{SBegin}, args...)...) }

func __do(args ...Node) Node { return Nodes(append([]Node{SDoBlock}, args...)...) }

func __move(dest, src Node) Node { return Nodes(SMove, dest, src) }

func __set(dest, src Node) Node { return Nodes(SSet, dest, src) }

func __less(lhs, rhs Node) Node { return Nodes(SLess, lhs, rhs) }

func __lessEq(lhs, rhs Node) Node { return Nodes(SLessEq, lhs, rhs) }

func __inc(subject, step Node) Node { return Nodes((SInc), subject, step) }

func __load(subject, key Node) Node { return Nodes((SLoad), subject, key) }

func __store(subject, key, value Node) Node { return Nodes((SStore), subject, value, key) }

func __if(cond, t, f Node) Node { return Nodes((SIf), cond, t, f) }

func __loop(body ...Node) Node { return Nodes(SFor, __chain(body...)) }

func __goto(label Node) Node { return Nodes(SGoto, label) }

func __label(name Node) Node { return Nodes(SLabel, name) }

func __ret(v Node) Node { return Nodes(SReturn, v) }

func __func(name Token, paramList Node, doc string, stats Node) Node {
	__findTailCall(stats.Nodes())
	funcname := Sym(name)
	p := name
	return __chain(
		__set(funcname, SNil).At(p),
		__move(funcname,
			Nodes((SFunc), funcname, paramList, stats, Str(doc)).At(p)).At(p),
	)
}

func __markupFuncName(recv, name Token) Token {
	name.Str = recv.Str + "." + name.Str
	return name
}

func __markupLambdaName(lambda Token) Token {
	lambda.Str = "<lambda" + strconv.Itoa(int(lambda.Pos.Line)) + ">"
	return lambda
}

func __call(cls, args Node) Node {
	return Nodes((SCall), cls, args)
}

func __findTailCall(stats []Node) {
	if len(stats) > 0 {
		x := stats[len(stats)-1]
		c := x.Nodes()
		if len(c) == 3 && c[0].Sym() == ACall {
			old := c[0].symLine
			c[0] = STailCall
			c[0].symLine = old
			return
		}

		if len(c) > 0 {
			if c[0].Sym() == (ABegin) {
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
		res = head.moveLoadStore(__move, src[0]).At(pos)
	} else {
		// a0, ..., an = b0, ..., bn
		names, retaddr := []Node{}, Nodes((SFreeAddr))
		for i := range nodes {
			names = append(names, randomVarname())
			retaddr = retaddr.append(names[i])
			res = res.append(__set(names[i], src[i]).At(pos))
		}
		for i, v := range nodes {
			res = res.append(v.moveLoadStore(__move, names[i]).At(pos))
		}
		res = res.append(retaddr)
	}
	return res
}

func __dotdotdot(expr Node) Node {
	expr.Nodes()[len(expr.Nodes())-1] = Nodes((SUnpack), expr.Nodes()[len(expr.Nodes())-1])
	return expr
}

func __forIn(key, value Token, expr, skip, body Node, pos Token) Node {
	k, v, e, tmp := Sym(key), Sym(value), randomVarname(), randomVarname()
	next := Nodes((SInc), e, k).At(pos)
	moveNext := __chain(
		__move(tmp, next).At(pos),
		__move(k, __load(tmp, zero).At(pos)).At(pos),
		__move(v, __load(tmp, one).At(pos)).At(pos),
	)
	skipExpr := __chain()
	if skip.Int() != 1 {
		skipVar, skipEnd, repeatLabel, exitLabel := randomVarname(), randomVarname(), randomVarname(), randomVarname()
		skipExpr = skipExpr.append(__set(skipEnd, skip).At(pos))
		moveNext = __chain(
			__set(skipVar, skipEnd).At(pos),
			__chain(
				__label(repeatLabel),
				__if(__lessEq(skipVar, zero).At(pos),
					__goto(exitLabel),
					__chain(
						moveNext,
						__if(
							Nodes((SEq), k, SNil).At(pos),
							breakNode,
							__inc(skipVar, Int(-1)).At(pos), // repeat
						).At(pos),
					),
				).At(pos),
				__goto(repeatLabel),
				__label(exitLabel),
			),
		)
	}
	return __do(
		__set(e, expr).At(pos),
		skipExpr,
		__set(k, SNil).At(pos),
		__set(tmp, next).At(pos), // init, move to the first key
		__move(k, __load(tmp, zero).At(pos)).At(pos),
		__set(v, __load(tmp, one).At(pos)).At(pos),
		__loop(
			__if(
				Nodes((SEq), k, SNil).At(pos),
				breakNode,
				__chain(
					body,
					moveNext,
				),
			).At(pos),
		),
	)
}

func randomVarname() Node {
	return staticSym("v" + strconv.FormatInt(rand.Int63(), 10)[:6])
}
