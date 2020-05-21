package parser

import (
	"unsafe"
)

type Atom string

var (
	Nnumber  = interfaceType(1.0)
	Nstring  = interfaceType("")
	Natom    = interfaceType(Atom(""))
	Ncomplex = interfaceType([]*Node{})
	Naddr    = interfaceType(uint16(1))

	chainNode = NewNode(AChain)
	breakNode = CompNode(NewNode(ABreak))
	nilNode   = NewNode(ANil)
	zeroNode  = NewNumberNode(0)
	moneNode  = NewNumberNode(-1)
	oneNode   = NewNumberNode(1)
	max32Node = NewNumberNode(0xffffffff)
	emptyNode = CompNode()
)

func interfaceType(a interface{}) uintptr {
	return (*(*[2]uintptr)(unsafe.Pointer(&a)))[0]
}

const (
	ADoBlock     Atom = "do"
	AConcat      Atom = "con"
	ANil         Atom = "nil"
	ASet         Atom = "set"
	AInc         Atom = "inc"
	AMove        Atom = "mov"
	AIf          Atom = "if"
	AFor         Atom = "for"
	APatchVararg Atom = "pvag"
	AFunc        Atom = "fun"
	ABreak       Atom = "brk"
	AContinue    Atom = "cont"
	AChain       Atom = "prog"
	ALoad        Atom = "load"
	AStore       Atom = "stor"
	ACall        Atom = "call"
	ASetB        Atom = "setb"
	AGetB        Atom = "getb"
	AReturn      Atom = "ret"
	AYield       Atom = "yld"
	AHash        Atom = "hash"
	AHashArray   Atom = "harr"
	AArray       Atom = "arr"
	AAdd         Atom = "+"
	ASub         Atom = "-"
	AMul         Atom = "*"
	ADiv         Atom = "/"
	AMod         Atom = "%"
	ABitAnd      Atom = "&"
	ABitOr       Atom = "|"
	ABitLsh      Atom = "<<"
	ABitRsh      Atom = ">>"
	ABitURsh     Atom = ">>>"
	ABitXor      Atom = "^"
	AEq          Atom = "=="
	ANeq         Atom = "!="
	AAnd         Atom = "and"
	AOr          Atom = "or"
	ANot         Atom = "not"
	ALess        Atom = "<"
	ALessEq      Atom = "<="
	AAddrOf      Atom = "addr"
	ALen         Atom = "len"
)

func __chain(args ...interface{}) *Node { return CompNode(append([]interface{}{AChain}, args...)...) }

func __do(args ...interface{}) *Node { return CompNode(append([]interface{}{ADoBlock}, args...)...) }

func __move(dest, src interface{}) *Node { return CompNode(AMove, dest, src) }

func __set(dest, src interface{}) *Node {
	if n, _ := dest.(*Node); n != nil && n.Type() != Natom {
		panic(&Error{Pos: n.Position, Message: "invalid assignments", Token: n.String()})
	}
	return CompNode(ASet, dest, src)
}

func __less(lhs, rhs interface{}) *Node { return CompNode(ALess, lhs, rhs) }

func __lessEq(lhs, rhs interface{}) *Node { return CompNode(ALessEq, lhs, rhs) }

func __mul(lhs, rhs interface{}) *Node { return CompNode(AMul, lhs, rhs) }

func __sub(lhs, rhs interface{}) *Node { return CompNode(ASub, lhs, rhs) }

func __inc(subject, step interface{}) *Node { return CompNode(AInc, subject, step) }

func __load(subject, key interface{}) *Node { return CompNode(ALoad, subject, key) }

func __call(cls, args interface{}) *Node { return CompNode(ACall, cls, args) }

func __return(value interface{}) *Node { return CompNode(AReturn, value) }

func __store(subject, key, value interface{}) *Node { return CompNode(AStore, subject, value, key) }

func (n *Node) __then(trueBranch *Node) *Node { return n.Cappend(trueBranch) }

func (n *Node) __else(falseBranch *Node) *Node { return n.Cappend(falseBranch) }

func __if(cond interface{}) *Node { return CompNode(AIf, cond) }

func (n *Node) __continue(c *Node) *Node { return n.Cappend(c) }

func (n *Node) __body(body *Node) *Node { return n.Cappend(body) }

func __for(cond interface{}) *Node { return CompNode(AFor, cond) }

func __func(paramlist interface{}) *Node { return CompNode(AFunc, emptyNode, paramlist) }
