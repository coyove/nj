package parser

import (
	"bytes"
	"io"
	"strconv"
)

const (
	NTNumber = iota
	NTString
	NTAtom
	NTCompound
	NTAddr
)

type Node struct {
	Type     byte
	Value    interface{}
	Pos      Position
	Compound []*Node
}

func NewCompoundNode(args ...interface{}) *Node {
	n := &Node{
		Type:     NTCompound,
		Compound: make([]*Node, 0),
	}
	for _, arg := range args {
		switch arg.(type) {
		case string:
			n.Compound = append(n.Compound, &Node{
				Type:  NTAtom,
				Value: arg.(string),
			})
		case *Node:
			if n.Pos.Source == "" {
				n.Pos = arg.(*Node).Pos
			}
			n.Compound = append(n.Compound, arg.(*Node))
		default:
			panic("shouldn't happen")
		}
	}
	return n
}

func NewAtomNode(tok Token) *Node {
	return &Node{
		Type:  NTAtom,
		Value: tok.Str,
		Pos:   tok.Pos,
	}
}

func NewStringNode(arg string) *Node {
	return &Node{
		Type:  NTString,
		Value: arg,
	}
}

func NewNumberNode(arg string) *Node {
	num, err := strconv.ParseFloat(arg, 64)
	if err != nil {
		panic(err)
	}

	return &Node{
		Type:  NTNumber,
		Value: num,
	}
}

func (n *Node) Dump(w io.Writer) {
	switch n.Type {
	case NTNumber:
		io.WriteString(w, "<"+strconv.FormatFloat(n.Value.(float64), 'f', 9, 64)+">")
	case NTString:
		io.WriteString(w, strconv.Quote(n.Value.(string)))
	case NTAtom:
		io.WriteString(w, n.Value.(string))
	case NTCompound:
		io.WriteString(w, "[")
		for _, a := range n.Compound {
			a.Dump(w)
			io.WriteString(w, " ")
		}
		io.WriteString(w, "]")
	}
}

func (n *Node) String() string {
	switch n.Type {
	case NTNumber:
		return "<" + strconv.FormatFloat(n.Value.(float64), 'f', 9, 64) + ">"
	case NTString:
		return strconv.Quote(n.Value.(string))
	case NTAtom:
		return n.Value.(string)
	case NTCompound:
		buf := &bytes.Buffer{}
		buf.WriteString("[")
		for _, a := range n.Compound {
			buf.WriteString(a.String() + " ")
		}
		buf.WriteString("]")
		return buf.String()
	}
	panic(1)
}
