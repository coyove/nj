package parser

import (
	"fmt"
	"io"
	"math"
	"reflect"
	"strconv"
	"strings"
)

type Position struct {
	Source string
	Line   uint32
	Column uint32
}

func (pos *Position) String() string {
	if pos.Source == "" {
		return ""
	}
	return fmt.Sprintf("%s:%d:%d", pos.Source, pos.Line, pos.Column)
}

type Token struct {
	Type uint32
	Str  string
	Pos  Position
}

func (self *Token) String() string {
	return self.Str
}

type Node struct {
	Position
	Value interface{}
}

func Nod(v interface{}) *Node {
	return &Node{Value: v}
}

func SymTok(tok Token) *Node {
	n := Nod(Symbol(tok.Str))
	n.SetPos(tok.Pos)
	return n
}

func Sym(s string) *Node {
	n := Nod(Symbol(s))
	return n
}

func Num(arg interface{}) *Node {
	n := Nod(nil)
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

func (n *Node) TypeName() string {
	return reflect.TypeOf(n.Value).String()
}

func (n *Node) SetPos(p interface{}) *Node {
	var m Position
	switch x := p.(type) {
	case *Node:
		m = x.Position
	case Token:
		m = x.Pos
	case Position:
		m = x
	default:
		panic(fmt.Sprintf("SetPos: shouldn't happen: %v", p))
	}
	n.Position = m
	return n
}

func (n *Node) Cpl() []*Node { a, _ := n.Value.([]*Node); return a }

func (n *Node) CplAppend(na ...*Node) *Node {
	n.Value = append(n.Cpl(), na...)
	if n.Position.Source == "" {
		for _, na := range na {
			if na.Position.Source != "" {
				n.Position = na.Position
				break
			}
		}
	}
	return n
}

func (n *Node) CplPrepend(n2 *Node) *Node {
	n.Value = append([]*Node{n2}, n.Cpl()...)
	return n
}

func (n *Node) CplIndex(i int) *Node { return n.Value.([]*Node)[i] }

func (n *Node) Str() string { a, _ := n.Value.(string); return a }

func (n *Node) Sym() Symbol { a, _ := n.Value.(Symbol); return a }

func (n *Node) Num() float64 { a, _ := n.Value.(float64); return a }

func Cpl(args ...interface{}) *Node {
	if len(args) == 3 {
		op, _ := args[0].(Symbol)
		a, _ := args[1].(*Node)
		b, _ := args[2].(*Node)
		if a != nil && b != nil {
			if op == AConcat && a.Type() == STR && b.Type() == STR {
				return Nod(a.Str() + b.Str())
			}
			if a.Type() == NUM && b.Type() == NUM {
				switch v1, v2 := a.Num(), b.Num(); op {
				case AAdd:
					return Num(v1 + v2)
				case ASub:
					return Num(v1 - v2)
				case AMul:
					return Num(v1 * v2)
				case ADiv:
					return Num(v1 / v2)
				case AMod:
					return Num(math.Mod(v1, v2))
				case APow:
					return Num(math.Pow(v1, v2))
				}
			}
		}
	}

	arr := make([]*Node, 0, len(args))
	n := Nod(arr)
	for _, arg := range args {
		switch x := arg.(type) {
		case string:
			if x == string(ABegin) {
				arr = append(arr, chainNode)
			} else {
				arr = append(arr, Nod(Symbol(x)))
			}
		case Symbol:
			if x == ABegin {
				arr = append(arr, chainNode)
			} else {
				arr = append(arr, Nod(x))
			}
		case *Node:
			if n.Source == "" {
				n.SetPos(x.Position)
			}
			arr = append(arr, x)
		case Token:
			arr = append(arr, SymTok(x))
		default:
			panic(fmt.Sprintf("Cpl: shouldn't happen: %v", x))
		}
	}
	n.Value = arr
	return n
}

func StringToNumber(arg string) (float64, error) {
	i, err := strconv.ParseInt(arg, 0, 64)
	if err == nil {
		return float64(i), nil
	}
	return strconv.ParseFloat(arg, 64)
}

func (n *Node) pos0(p interface{}) *Node {
	if n.Type() != CPL {
		return n
	}
	n.CplIndex(0).SetPos(p)
	n.SetPos(p)
	return n
}

func (n *Node) SetPos0(p interface{}) *Node { return n.pos0(p) }

func (n *Node) Dump(w io.Writer, ident string) {
	io.WriteString(w, ident)
	switch n.Type() {
	case NUM:
		io.WriteString(w, strconv.FormatFloat(n.Value.(float64), 'f', -1, 64))
	case STR:
		io.WriteString(w, strconv.Quote(n.Value.(string)))
	case SYM:
		io.WriteString(w, n.String())
	case CPL:
		nocpl := true
		for _, a := range n.Cpl() {
			if a.Type() == CPL {
				nocpl = false
				break
			}
		}

		if !nocpl {
			io.WriteString(w, "[\n")
			for _, a := range n.Cpl() {
				a.Dump(w, "  "+ident)
				if a.Type() != CPL {
					io.WriteString(w, "\n")
				}
			}
			io.WriteString(w, ident)
		} else {
			io.WriteString(w, "[")
			for i, a := range n.Cpl() {
				a.Dump(w, "")
				if i < len(n.Cpl())-1 {
					io.WriteString(w, " ")
				}
			}
		}
		io.WriteString(w, "]\n")
	}
}

func (n *Node) String() string {
	switch n.Type() {
	case NUM:
		return strconv.FormatFloat(n.Value.(float64), 'f', -1, 64)
	case STR:
		return strconv.Quote(n.Value.(string))
	case SYM:
		return string(n.Value.(Symbol)) + "@" + n.Position.String()
	case CPL:
		buf := make([]string, len(n.Cpl()))
		for i, a := range n.Cpl() {
			buf[i] = a.String()
		}
		return "[" + strings.Join(buf, " ") + "]"
	case ADR:
		return "0x" + strconv.FormatInt(int64(n.Value.(uint16)), 16)
	}
	panic("shouldn't happen")
}

func (n *Node) isSimpleAddSub() (a Symbol, v float64) {
	if n.Type() != CPL || len(n.Cpl()) < 3 {
		return
	}
	// a = a + v
	if n.CplIndex(1).Type() == SYM && n.CplIndex(0).Sym() == "+" && n.CplIndex(2).Type() == NUM {
		a, v = n.CplIndex(1).Sym(), n.CplIndex(2).Value.(float64)
	}
	// a = v + a
	if n.CplIndex(2).Type() == SYM && n.CplIndex(0).Sym() == "+" && n.CplIndex(1).Type() == NUM {
		a, v = n.CplIndex(2).Sym(), n.CplIndex(1).Value.(float64)
	}
	// a = a - v
	if n.CplIndex(1).Type() == SYM && n.CplIndex(0).Sym() == "-" && n.CplIndex(2).Type() == NUM {
		a, v = n.CplIndex(1).Sym(), -n.CplIndex(2).Value.(float64)
	}
	return
}

func (n *Node) moveLoadStore(sm func(interface{}, interface{}) *Node, v *Node) *Node {
	if len(n.Cpl()) == 3 && n.CplIndex(0).Sym() == ALoad {
		return __store(n.CplIndex(1), n.CplIndex(2), v)
	}
	if n.Type() != SYM {
		panic(fmt.Sprintf("%#v: invalid assignment"))
	}
	return sm(n, v)
}
