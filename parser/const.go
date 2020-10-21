package parser

import (
	"math/rand"
	"strconv"
	"strings"
	"unsafe"
)

type Symbol struct {
	Position
	Text string
}

func (s Symbol) Equals(s2 Symbol) bool { return s.Text == s2.Text }

func (s Symbol) SetPos(pos Position) Symbol { s.Position = pos; return s }

func (s Symbol) String() string { return s.Text + "@" + s.Position.String() }

var (
	NUM    = interfaceType(1.0)
	numINT = interfaceType(int64(1))
	STR    = interfaceType("")
	SYM    = interfaceType(Symbol{})
	CPL    = interfaceType([]Node{})
	ADR    = interfaceType(uint16(1))

	breakNode     = Cpl(Node{ABreak})
	popvNode      = Cpl(Node{APopV})
	popvClearNode = Cpl(Node{APopVClear})
	zeroNode      = Num(0)
	oneNode       = Num(1)
	emptyNode     = Cpl()
)

func interfaceType(a interface{}) uintptr {
	return (*(*[2]uintptr)(unsafe.Pointer(&a)))[0]
}

var (
	ANop       = Symbol{Text: "nop"}
	ADoBlock   = Symbol{Text: "do"}
	AConcat    = Symbol{Text: "con"}
	ANil       = Symbol{Text: "nil"}
	ASet       = Symbol{Text: "set"}
	AInc       = Symbol{Text: "inc"}
	AMove      = Symbol{Text: "mov"}
	AIf        = Symbol{Text: "if"}
	AFor       = Symbol{Text: "for"}
	AFunc      = Symbol{Text: "fun"}
	ABreak     = Symbol{Text: "break"}
	AContinue  = Symbol{Text: "cont"}
	ABegin     = Symbol{Text: "prog"}
	ALoad      = Symbol{Text: "load"}
	AStore     = Symbol{Text: "stor"}
	ACall      = Symbol{Text: "call"}
	ATailCall  = Symbol{Text: "tail"}
	AReturn    = Symbol{Text: "ret"}
	AYield     = Symbol{Text: "yield"}
	AAdd       = Symbol{Text: "add"}
	ASub       = Symbol{Text: "sub"}
	AUnm       = Symbol{Text: "unm"}
	AMul       = Symbol{Text: "mul"}
	ADiv       = Symbol{Text: "div"}
	AMod       = Symbol{Text: "mod"}
	APow       = Symbol{Text: "pow"}
	AEq        = Symbol{Text: "eq"}
	ANeq       = Symbol{Text: "neq"}
	AAnd       = Symbol{Text: "and"}
	AOr        = Symbol{Text: "or"}
	ANot       = Symbol{Text: "not"}
	ALess      = Symbol{Text: "lt"}
	ALessEq    = Symbol{Text: "le"}
	ALen       = Symbol{Text: "len"}
	ARetAddr   = Symbol{Text: "reta"}
	APopV      = Symbol{Text: "popv"}
	APopVClear = Symbol{Text: "endv"}
	APopVAll   = Symbol{Text: "allv"}
	APopVAllA  = Symbol{Text: "allva"}
	ALabel     = Symbol{Text: "label"}
	AGoto      = Symbol{Text: "goto"}
	ARef       = Symbol{Text: "ref"}
	ADeref     = Symbol{Text: "deref"}
)

func __chain(args ...Node) Node { return Cpl(append([]Node{Node{ABegin}}, args...)...) }

func __do(args ...Node) Node { return Cpl(append([]Node{Node{ADoBlock}}, args...)...) }

func __removeddd(dest Node) Node {
	sym := dest.Value.(Symbol)
	if sym.Text != "..." {
		sym.Text = strings.TrimLeft(sym.Text, ".")
		dest.Value = sym
	}
	return dest
}

func __move(dest, src Node) Node { return Cpl(Node{AMove}, __removeddd(dest), src) }

func __set(dest, src Node) Node { return Cpl(Node{ASet}, __removeddd(dest), src) }

func __less(lhs, rhs Node) Node { return Cpl(Node{ALess}, lhs, rhs) }

func __lessEq(lhs, rhs Node) Node { return Cpl(Node{ALessEq}, lhs, rhs) }

func __inc(subject, step Node) Node { return Cpl(Node{AInc}, subject, step) }

func __load(subject, key Node) Node { return Cpl(Node{ALoad}, subject, key) }

func __store(subject, key, value Node) Node { return Cpl(Node{AStore}, subject, value, key) }

func __if(cond, truebody, falsebody Node) Node { return Cpl(Node{AIf}, cond, truebody, falsebody) }

func __loop(body Node) Node { return Cpl(Node{AFor}, body) }

func __func(name, paramlist, body Node) Node { return Cpl(Node{AFunc}, name, paramlist, body) }

func __call(cls, args Node) Node { return Cpl(Node{ACall}, cls, args) }

func __popvAll(i int, k Node) Node {
	if i == 0 {
		return __chain(k, Cpl(Node{APopVAllA}))
	}
	return Cpl(Node{APopVAll})
}

func randomVarname() Node {
	return Sym("v" + strconv.FormatInt(rand.Int63(), 10))
}

func forLoop(pos Position, rcv []Node, exprIters []Node, body Node) Node {
	iter := randomVarname()
	subject := randomVarname()
	r := __do(__set(iter, exprIters[0]).SetPos(pos))
	if len(exprIters) > 1 {
		r = r.CplAppend(__set(subject, exprIters[1]).SetPos(pos))
	} else {
		r = r.CplAppend(__set(subject, popvNode).SetPos(pos))
	}
	if len(exprIters) > 2 {
		r = r.CplAppend(__set(rcv[0], exprIters[2]).SetPos(pos))
	} else {
		r = r.CplAppend(__set(rcv[0], __chain(popvNode, popvClearNode)).SetPos(pos))
	}
	rr := __chain()
	for i := 1; i < len(rcv); i++ {
		if i == len(rcv)-1 {
			rr = rr.CplAppend(__set(rcv[i], __chain(popvNode, popvClearNode)).SetPos(pos))
		} else {
			rr = rr.CplAppend(__set(rcv[i], popvNode).SetPos(pos))
		}
	}
	r = r.CplAppend(__loop(
		__chain(
			__move(rcv[0], __call(iter, Cpl(subject, rcv[0])).SetPos(pos)).SetPos(pos),
			rr,
			__if(rcv[0], body, breakNode).SetPos(pos),
		),
	).SetPos(pos))
	return r
}
