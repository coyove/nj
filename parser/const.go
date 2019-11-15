package parser

import (
	"fmt"
	"unsafe"

	"github.com/coyove/common/lru"
)

type Atom string

var (
	Nnumber   = interfaceType(1.0)
	Nstring   = interfaceType("")
	Natom     = interfaceType(Atom(""))
	Ncompound = interfaceType([]*Node{})
	Naddr     = interfaceType(uint16(1))

	chainNode = NewNode(AChain)
	nilNode   = NewNode(ANil)
	zeroNode  = NewNumberNode(0)
	moneNode  = NewNumberNode(-1)
	oneNode   = NewNumberNode(1)
	max32Node = NewNumberNode(0xffffffff)
	emptyNode = CompNode()
)

func interfaceType(a interface{}) uintptr { return (*(*[2]uintptr)(unsafe.Pointer(&a)))[0] }

const (
	AAssert    Atom = "assert"
	ANil       Atom = "nil"
	ASet       Atom = "set"
	AInc       Atom = "inc"
	AMove      Atom = "move"
	AIf        Atom = "if"
	AFor       Atom = "for"
	AForeach   Atom = "foreach"
	AFunc      Atom = "func"
	ABreak     Atom = "break"
	AContinue  Atom = "continue"
	AChain     Atom = "chain"
	ALoad      Atom = "load"
	AStore     Atom = "store"
	ACall      Atom = "call"
	ASetB      Atom = "setb"
	ASetFromAB Atom = "setab"
	AReturn    Atom = "ret"
	AYield     Atom = "yield"
	ASlice     Atom = "slice"
	AMap       Atom = "map"
	AArray     Atom = "array"
	AAdd       Atom = "+"
	ASub       Atom = "-"
	AMul       Atom = "*"
	ADiv       Atom = "/"
	AMod       Atom = "%"
	ABitAnd    Atom = "&"
	ABitOr     Atom = "|"
	ABitLsh    Atom = "<<"
	ABitRsh    Atom = ">>"
	ABitURsh   Atom = ">>>"
	ABitXor    Atom = "^"
	AEq        Atom = "=="
	ANeq       Atom = "!="
	AAnd       Atom = "and"
	AOr        Atom = "or"
	ANot       Atom = "!"
	ALess      Atom = "<"
	ALessEq    Atom = "<="
	AAddrOf    Atom = "addressof"
	ATypeOf    Atom = "typeof"
	ALen       Atom = "len"
	ADDD       Atom = "..."
)

func __chain(args ...interface{}) *Node { return CompNode(append([]interface{}{AChain}, args...)...) }

func __move(dest, src interface{}) *Node { return CompNode(AMove, dest, src) }

func __set(dest, src interface{}) *Node { return CompNode(ASet, dest, src) }

func __lessEq(lhs, rhs interface{}) *Node { return CompNode(ALessEq, lhs, rhs) }

func __mul(lhs, rhs interface{}) *Node { return CompNode(AMul, lhs, rhs) }

func __sub(lhs, rhs interface{}) *Node { return CompNode(ASub, lhs, rhs) }

func __inc(subject, step interface{}) *Node { return CompNode(AInc, subject, step) }

func __load(subject, key interface{}) *Node { return CompNode(ALoad, subject, key) }

func __call(cls, args interface{}) *Node { return CompNode(ACall, cls, args) }

func __return(value interface{}) *Node { return CompNode(AReturn, value) }

func __store(subject, key, value interface{}) *Node { return CompNode(AStore, subject, key, value) }

func (n *Node) __then(trueBranch *Node) *Node { return n.Cappend(trueBranch) }

func (n *Node) __else(falseBranch *Node) *Node { return n.Cappend(falseBranch) }

func __if(cond interface{}) *Node { return CompNode(AIf, cond) }

func (n *Node) __continue(c *Node) *Node { return n.Cappend(c) }

func (n *Node) __body(body *Node) *Node { return n.Cappend(body) }

func __for(cond interface{}) *Node { return CompNode(AFor, cond) }

func (n *Node) __params(params *Node) *Node { return n.Cappend(params) }

func __func(name interface{}) *Node { return CompNode(AFunc, name) }

func __hash(str string) *Node { return NewNumberNode(HashString(str)) }

var hashDedupCache = lru.NewCache(1024)

func HashString(str string) float64 {
	var hash uint64 = 2166136261
	for _, c := range str {
		hash *= 16777619
		hash ^= uint64(c)
	}
	hash &= 0xffffffffffff
	if t, ok := hashDedupCache.Get(hash); ok && t.(string) != str {
		panic(fmt.Sprint(t, "and", str, "share an identical FNV48 hash:", hash))
	}
	hashDedupCache.Add(hash, str)
	return float64(hash)
}
