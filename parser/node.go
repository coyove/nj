package parser

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

type Meta struct {
	Source string
	Line   uint32
	Column uint16
	Type   byte
}

func (pos *Meta) String() string {
	if pos.Source == "" {
		return fmt.Sprintf("0:%d:%d", pos.Line, pos.Column)
	}
	return fmt.Sprintf("%s:%d:%d", pos.Source, pos.Line, pos.Column)
}

type Token struct {
	Type uint32
	Str  string
	Pos  Meta
}

func (self *Token) String() string {
	return self.Str
}

const (
	Nnumber = iota
	Nstring
	Natom
	Ncompound
	Naddr
)

type Node struct {
	Meta
	Value interface{}
}

func NewNode(t byte) *Node {
	n := &Node{}
	n.Type = t
	return n
}

func (n *Node) SetPos(p interface{}) *Node {
	var m Meta
	switch x := p.(type) {
	case *Node:
		m = x.Meta
	case Token:
		m = x.Pos
	case Meta:
		m = x
	default:
		panic("shouldn't happen")
	}
	n.Meta.Column = m.Column
	n.Meta.Line = m.Line
	n.Meta.Source = m.Source
	return n
}

func (n *Node) SetValue(v interface{}) *Node { n.Value = v; return n }

func (n *Node) C() []*Node { return n.Value.([]*Node) }

func (n *Node) Cappend(na *Node) *Node { n.Value = append(n.C(), na); return n }

func (n *Node) Cx(i int) *Node { return n.Value.([]*Node)[i] }

func (n *Node) Cn() int { a, _ := n.Value.([]*Node); return len(a) }

func (n *Node) Cy() bool { _, ok := n.Value.([]*Node); return ok }

func (n *Node) S() string { a, _ := n.Value.(string); return a }

func (n *Node) N() float64 { a, _ := n.Value.(float64); return a }

func CNode(args ...interface{}) *Node {
	n := NewNode(Ncompound)
	arr := make([]*Node, 0, len(args))
	for _, arg := range args {
		switch arg.(type) {
		case string:
			arr = append(arr, NewNode(Natom).SetValue(arg))
		case *Node:
			if n.Source == "" {
				n.SetPos(arg.(*Node).Meta)
			}
			arr = append(arr, arg.(*Node))
		default:
			panic("shouldn't happen")
		}
	}
	n.Value = arr
	return n
}

func ANode(tok Token) *Node {
	n := NewNode(Natom)
	n.Value = tok.Str
	n.SetPos(tok.Pos)
	return n
}

func NilNode() *Node {
	n := NewNode(Natom)
	n.Value = "nil"
	return n
}

func SNode(arg string) *Node {
	n := NewNode(Nstring)
	n.Value = arg
	return n
}

func NNode(arg interface{}) *Node {
	n := NewNode(Nnumber)
	switch x := arg.(type) {
	case string:
		num, err := StringToNumber(x)
		if err != nil {
			panic(err)
		}
		n.Value = num
	case float64:
		n.Value = x
	default:
		panic("shouldn't happen")
	}
	return n
}

func StringToNumber(arg string) (float64, error) {
	if arg[0] == '0' && len(arg) > 1 {
		var num uint64
		var err error
		switch arg[1] {
		case 'x', 'X':
			num, err = strconv.ParseUint(arg[2:], 16, 64)
		case 'b', 'B':
			num, err = strconv.ParseUint(arg[2:], 2, 64)
		case 'i', 'I':
			num, err = strconv.ParseUint(arg[2:], 16, 64)
			if err == nil {
				return math.Float64frombits(uint64(num)), nil
			}
		default:
			num, err = strconv.ParseUint(arg[1:], 8, 64)
		}
		if err == nil {
			return float64(num), nil
		}
	}

	num, err := strconv.ParseFloat(arg, 64)
	if err != nil {
		return 0, err
	}

	return num, nil
}

func (n *Node) setPos0(p interface{}) *Node { n.Cx(0).SetPos(p); return n }

func (n *Node) setPos(p interface{}) *Node { n.SetPos(p); return n }

func (n *Node) Dump(w io.Writer) {
	switch n.Type {
	case Nnumber:
		io.WriteString(w, "<"+strconv.FormatFloat(n.Value.(float64), 'f', 9, 64)+">")
	case Nstring:
		io.WriteString(w, strconv.Quote(n.Value.(string)))
	case Natom:
		io.WriteString(w, n.Value.(string))
	case Ncompound:
		io.WriteString(w, "[")
		for _, a := range n.C() {
			a.Dump(w)
			io.WriteString(w, " ")
		}
		io.WriteString(w, "]")
	}
}

func (n *Node) String() string {
	pos := fmt.Sprintf("@%s:%d:%d", n.Source, n.Line, n.Column)
	switch n.Type {
	case Nnumber:
		return strconv.FormatFloat(n.Value.(float64), 'f', 9, 64) + pos
	case Nstring:
		return strconv.Quote(n.Value.(string)) + pos
	case Natom:
		return n.Value.(string) + pos
	case Ncompound:
		buf := make([]string, n.Cn())
		for i, a := range n.C() {
			buf[i] = a.String()
		}
		return "[" + strings.Join(buf, " ") + "]" + pos
	}
	panic("shouldn't happen")
}

func (n *Node) isIsolatedDupCall() bool {
	if n.Cn() < 3 || n.Cx(0).S() != "call" || n.Cx(1).S() != "copy" {
		return false
	}
	// [call copy [a0, a1, a2]]
	return n.Cx(2).Cn() == 3
}

func (n *Node) isSimpleAddSub() (a string, b string, s float64) {
	if n.Type != Ncompound || n.Cn() < 3 {
		return
	}
	s = 1
	if c := n.Cx(0).S(); c != "+" && c != "-" {
		return
	} else if c == "-" {
		s = -1
	}
	if c := n.Cx(1).S(); n.Cx(1).Type == Natom {
		a = c
		if n.Cx(2).Type != Nnumber {
			a = ""
		}
	}
	if c := n.Cx(2).S(); n.Cx(2).Type == Natom {
		b = c
		if n.Cx(1).Type != Nnumber {
			b = ""
		}
	}
	return
}
