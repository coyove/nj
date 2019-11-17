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
	Column uint32
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

type Node struct {
	Meta
	Value interface{}
}

func NewNode(v interface{}) *Node {
	return &Node{Value: v}
}

func ANode(tok Token) *Node {
	n := NewNode(Atom(tok.Str))
	n.SetPos(tok.Pos)
	return n
}

func ANodeS(s string) *Node {
	n := NewNode(Atom(s))
	return n
}

func NewNumberNode(arg interface{}) *Node {
	n := NewNode(nil)
	switch x := arg.(type) {
	case string:
		num, err := StringToNumber(x)
		if err != nil {
			panic(err)
		}
		n.Value = num
	case float64:
		n.Value = x
	case int:
		n.Value = float64(x)
	default:
		panic("shouldn't happen")
	}
	return n
}

func (n *Node) Type() uintptr {
	return interfaceType(n.Value)
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
		panic(fmt.Sprintf("SetPos: shouldn't happen: %v", p))
	}
	n.Meta = m
	return n
}

func (n *Node) C() []*Node { return n.Value.([]*Node) }

func (n *Node) Cappend(na ...*Node) *Node {
	n.Value = append(n.C(), na...)
	if n.Meta.Source == "" {
		for _, na := range na {
			if na.Meta.Source != "" {
				n.Meta = na.Meta
				break
			}
		}
	}
	return n
}

func (n *Node) Cprepend(n2 *Node) *Node {
	n.Value = append([]*Node{n2}, n.C()...)
	return n
}

func (n *Node) Cx(i int) *Node { return n.Value.([]*Node)[i] }

func (n *Node) Cn() int { a, _ := n.Value.([]*Node); return len(a) }

func (n *Node) Cy() bool { _, ok := n.Value.([]*Node); return ok }

func (n *Node) S() string { a, _ := n.Value.(string); return a }

func (n *Node) A() Atom { a, _ := n.Value.(Atom); return a }

func (n *Node) N() float64 { a, _ := n.Value.(float64); return a }

func CompNode(args ...interface{}) *Node {
	const max32 = 0xffffffff

	if len(args) >= 2 {
		op, _ := args[0].(string)
		if op == "" {
			x, _ := args[0].(Atom)
			op = string(x)
		}
		a, _ := args[1].(*Node)
		if len(args) == 3 {
			b, _ := args[2].(*Node)
			if a != nil && b != nil {
				if op == "+" && a.Type() == Nstring && b.Type() == Nstring {
					return NewNode(a.Value.(string) + b.Value.(string))
				}
				if a.Type() == Nnumber && b.Type() == Nnumber {
					v1, v2 := a.Value.(float64), b.Value.(float64)
					switch op {
					case "+":
						return NewNumberNode(v1 + v2)
					case "-":
						return NewNumberNode(v1 - v2)
					case "*":
						return NewNumberNode(v1 * v2)
					case "/":
						return NewNumberNode(v1 / v2)
					case "%":
						return NewNumberNode(math.Mod(v1, v2))
					case "==":
						if v1 == v2 {
							return NewNumberNode(1)
						}
						return NewNumberNode(0)
					case "!=":
						if v1 != v2 {
							return NewNumberNode(1)
						}
						return NewNumberNode(0)
					case "<":
						if v1 < v2 {
							return NewNumberNode(1)
						}
						return NewNumberNode(0)
					case "<=":
						if v1 <= v2 {
							return NewNumberNode(1)
						}
						return NewNumberNode(0)
					case "&":
						return NewNumberNode(float64(int32(int64(v1)&max32) & int32(int64(v2)&max32)))
					case "|":
						return NewNumberNode(float64(int32(int64(v1)&max32) | int32(int64(v2)&max32)))
					case "^":
						return NewNumberNode(float64(int32(int64(v1)&max32) ^ int32(int64(v2)&max32)))
					case "<<":
						return NewNumberNode(float64(int32(int64(v1)&max32) << uint(v2)))
					case ">>":
						return NewNumberNode(float64(int32(int64(v1)&max32) >> uint(v2)))
					case ">>>":
						return NewNumberNode(float64(uint32(uint64(v1)&max32) >> uint(v2)))
					}
				}
			}
		} else if len(args) == 2 {
			if a != nil && a.Type() == Nnumber {
				v1 := a.Value.(float64)
				switch op {
				case "not":
					if v1 == 0 {
						return NewNumberNode(1)
					}
					return NewNumberNode(0)
				case "~":
					return NewNumberNode(float64(^int32(int64(v1) & max32)))
				}
			}
		}
	}

	arr := make([]*Node, 0, len(args))
	n := NewNode(arr)
	for _, arg := range args {
		switch x := arg.(type) {
		case string:
			if x == string(AChain) {
				arr = append(arr, chainNode)
			} else {
				arr = append(arr, NewNode(Atom(x)))
			}
		case Atom:
			if x == AChain {
				arr = append(arr, chainNode)
			} else {
				arr = append(arr, NewNode(x))
			}
		case *Node:
			if n.Source == "" {
				n.SetPos(x.Meta)
			}
			arr = append(arr, x)
		case Token:
			arr = append(arr, ANode(x))
		default:
			panic(fmt.Sprintf("CompNode: shouldn't happen: %v", x))
		}
	}
	n.Value = arr
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

func (n *Node) pos0(p interface{}) *Node {
	if n.Type() == Nnumber || n.Type() == Nstring {
		return n
	}
	n.Cx(0).SetPos(p)
	return n
}

func (n *Node) setPos(p interface{}) *Node { n.SetPos(p); return n }

func (n *Node) SetPos0(p interface{}) *Node { return n.pos0(p) }

func (n *Node) Dump(w io.Writer) {
	switch n.Type() {
	case Nnumber:
		io.WriteString(w, "<"+strconv.FormatFloat(n.Value.(float64), 'f', 9, 64)+">")
	case Nstring:
		io.WriteString(w, strconv.Quote(n.Value.(string)))
	case Natom:
		io.WriteString(w, string(n.Value.(Atom)))
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
	pos := ""
	if n.Source != "" {
		pos = "@" + n.Meta.String()
	}
	switch n.Type() {
	case Nnumber:
		return strconv.FormatFloat(n.Value.(float64), 'f', -1, 64) + pos
	case Nstring:
		return strconv.Quote(n.Value.(string)) + pos
	case Natom:
		return string(n.Value.(Atom)) + pos
	case Ncompound:
		buf := make([]string, n.Cn())
		for i, a := range n.C() {
			buf[i] = a.String()
		}
		return "[" + strings.Join(buf, " ") + "]" + pos
	case Naddr:
		return "0x" + strconv.FormatInt(int64(n.Value.(uint16)), 16) + pos
	}
	panic("shouldn't happen")
}

func (n *Node) isSimpleAddSub() (a Atom, v float64) {
	if n.Type() != Ncompound || n.Cn() < 3 {
		return
	}

	// a = a + v
	if n.Cx(1).Type() == Natom && n.Cx(0).A() == "+" && n.Cx(2).Type() == Nnumber {
		a, v = n.Cx(1).A(), n.Cx(2).Value.(float64)
	}

	// a = v + a
	if n.Cx(2).Type() == Natom && n.Cx(0).A() == "+" && n.Cx(1).Type() == Nnumber {
		a, v = n.Cx(2).A(), n.Cx(1).Value.(float64)
	}

	// a = a - v
	if n.Cx(1).Type() == Natom && n.Cx(0).A() == "-" && n.Cx(2).Type() == Nnumber {
		a, v = n.Cx(1).A(), -n.Cx(2).Value.(float64)
	}

	return
}
