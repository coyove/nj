package parser

import "unsafe"

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
	oneNode   = NewNumberNode(1)
	moneNode  = NewNumberNode(-1)
	max32Node = NewNumberNode(0xffffffff)
	emptyNode = CompNode()
)

func interfaceType(a interface{}) uintptr { return (*(*[2]uintptr)(unsafe.Pointer(&a)))[0] }

const (
	ANil   Atom = "nil"
	AChain Atom = "chain"
	ALoad  Atom = "load"
	ASlice Atom = "slice"
	AAdd   Atom = "+"
	ASub   Atom = "-"
	AMul   Atom = "*"
	ALess  Atom = "<"
)
