package parser

import (
	"fmt"
	"io"
	"math"
	"math/big"
	"strconv"
	"strings"
)

type Position struct {
	Source string
	Line   uint32
	Column uint32
	// a      []int
}

func (pos *Position) String() string {
	// if pos.Source == "" {
	// 	return ""
	// }
	// return fmt.Sprintf("%s:%d:%d", pos.Source, pos.Line, pos.Column)
	return fmt.Sprintf("line %d", pos.Line)
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
	Type    byte   // Node type
	Addr    uint16 // raw address value
	symLine uint32 // symbol position
	num     uint64 // float64 or int64 value
	strSym  string // string or symbol value
	Nodes   []Node // Nodes value
}

func NewSymbolFromToken(tok Token) Node {
	return Node{Type: Symbol, symLine: tok.Pos.Line, strSym: tok.Str}
}

func NewAddress(yx uint16) Node { return Node{Type: Address, Addr: yx} }

func NewSymbol(s string) Node { return Node{Type: Symbol, strSym: s} }

func NewString(s string) Node { return Node{Type: String, strSym: s} }

func NewNumberFromBig(x *big.Float) Node {
	n := Node{}
	v, acc := x.Int64()
	if acc == big.Exact {
		n.Type = Int
		n.num = uint64(v)
	} else {
		n.Type = Float
		f, _ := x.Float64()
		n.num = math.Float64bits(f)
	}
	return n
}

func NewNumberFromString(v string) Node {
	i, err := strconv.ParseInt(v, 0, 64)
	if err == nil {
		if int64(float64(i)) == i {
			return Node{Type: Float, num: math.Float64bits(float64(i))}
		}
		return Node{Type: Int, num: uint64(i)}
	}
	f, _ := strconv.ParseFloat(v, 64)
	return Node{Type: Float, num: math.Float64bits(f)}
}

func NewNumberFromInt(v int64) Node {
	if int64(float64(v)) == v {
		return Node{Type: Float, num: math.Float64bits(float64(v))}
	}
	return Node{num: uint64(v), Type: Int}
}

func (n Node) Valid() bool {
	t := n.Type
	return t == Symbol || t == Float || t == String || t == Complex || t == Address || t == Int
}

func (n Node) IsNumber() bool {
	t := n.Type
	return t == Float || t == Int
}

func (n Node) StringValue() string { return n.strSym }

func (n Node) SymbolValue() string { return n.strSym }

func (n Node) IsSymbolDotDotDot() bool { return strings.HasPrefix(n.strSym, "...") }

func (n Node) IntValue() int64 { return int64(n.num) }

func (n Node) FloatValue() float64 { return math.Float64frombits(n.num) }

func (n Node) NumberValue() (float64, int64, bool) {
	if n.Type == Int {
		return (&big.Float{}).SetFloat64(n.FloatValue())
	}
	x := (&big.Float{}).SetInt64(n.IntValue())
	return x
}

func (n Node) IsNegativeNumber() bool {
	if n.Type == Float {
		return n.FloatValue() < 0
	}
	return n.IntValue() < 0
}

func NewComplex(args ...Node) Node {
	if len(args) == 3 {
		op := args[0].SymbolValue()
		a, b := args[1], args[2]
		if op == AConcat && a.Type == String && b.Type == String {
			return NewString(a.StringValue() + b.StringValue())
		}
		if a.IsNumber() && b.IsNumber() {
			// switch v1, v2 := a.BigValue(), b.BigValue(); op {
			// case AAdd:
			// 	return NewNumberFromBig(v1.Add(v1, v2))
			// case ASub:
			// 	return NewNumberFromBig(v1.Sub(v1, v2))
			// case AMul:
			// 	return NewNumberFromBig(v1.Mul(v1, v2))
			// case ADiv:
			// 	return NewNumberFromBig(v1.Quo(v1, v2))
			// }
		}
	}
	return Node{Nodes: args, Type: Complex}
}

func (n Node) SetPos(pos Position) Node {
	switch n.Type {
	case Symbol:
		n.symLine = pos.Line
	case Complex:
		c := n.Nodes
		if len(c) > 0 {
			c[0] = c[0].SetPos(pos)
		}
	}
	return n
}

func (n Node) Pos() Position {
	switch n.Type {
	case Symbol:
		return Position{Line: n.symLine}
	case Complex:
		c := n.Nodes
		if len(c) == 0 {
			panic("Pos()")
		}
		return c[0].Pos()
	default:
		panic("Pos()")
	}
}

func (n Node) Dump(w io.Writer, ident string) {
	io.WriteString(w, ident)
	switch n.Type {
	case Complex:
		nocpl := true
		for _, a := range n.Nodes {
			if a.Type == Complex {
				nocpl = false
				break
			}
		}

		if !nocpl {
			io.WriteString(w, "[\n")
			for _, a := range n.Nodes {
				a.Dump(w, "  "+ident)
				if a.Type != Complex {
					io.WriteString(w, "\n")
				}
			}
			io.WriteString(w, ident)
		} else {
			io.WriteString(w, "[")
			for i, a := range n.Nodes {
				a.Dump(w, "")
				if i < len(n.Nodes)-1 {
					io.WriteString(w, " ")
				}
			}
		}
		io.WriteString(w, "]\n")
	case Symbol:
		io.WriteString(w, fmt.Sprintf("%s.%d", n.strSym, n.symLine))
	default:
		io.WriteString(w, n.String())
	}
}

func (n Node) String() string {
	switch n.Type {
	case Int:
		return strconv.FormatInt(n.IntValue(), 10)
	case Float:
		return strconv.FormatFloat(n.FloatValue(), 'f', -1, 64)
	case String:
		return strconv.Quote(n.strSym)
	case Symbol:
		return fmt.Sprintf("%s at line %d", n.strSym, n.symLine)
	case Complex:
		buf := make([]string, len(n.Nodes))
		for i, a := range n.Nodes {
			buf[i] = a.String()
		}
		return "[" + strings.Join(buf, " ") + "]"
	case Address:
		return "0x" + strconv.FormatInt(int64(n.Addr), 16)
	default:
		return fmt.Sprintf("FIXME: %#v", n)
	}
}

func (n Node) append(n2 Node) Node {
	n.Nodes = append(n.Nodes, n2)
	return n
}

// func (n Node) isSimpleAddSub() (a Symbol, v float64, ok bool) {
// 	if n.Type() != CPL || len(n.Cpl()) < 3 {
// 		return
// 	}
// 	// a = a + v
// 	if n.CplIndex(1).Type() == SYM && n.CplIndex(0).Sym().Equals(AAdd) && n.CplIndex(2).Type() == NUM {
// 		a, v, ok = n.CplIndex(1).Sym(), n.CplIndex(2).Value.(float64), true
// 	}
// 	// a = v + a
// 	if n.CplIndex(2).Type() == SYM && n.CplIndex(0).Sym().Equals(AAdd) && n.CplIndex(1).Type() == NUM {
// 		a, v, ok = n.CplIndex(2).Sym(), n.CplIndex(1).Value.(float64), true
// 	}
// 	// a = a - v
// 	if n.CplIndex(1).Type() == SYM && n.CplIndex(0).Sym().Equals(ASub) && n.CplIndex(2).Type() == NUM {
// 		a, v, ok = n.CplIndex(1).Sym(), -n.CplIndex(2).Value.(float64), true
// 	}
// 	return
// }

func (n Node) isCallStat() bool {
	return len(n.Nodes) > 0 && n.Nodes[0].SymbolValue() == ACall
}

func (n Node) moveLoadStore(sm func(Node, Node) Node, v Node) Node {
	if len(n.Nodes) == 3 && n.Nodes[0].SymbolValue() == ALoad {
		return __store(n.Nodes[1], n.Nodes[2], v)
	}
	if n.Type != Symbol {
		panic(fmt.Sprintf("%v: invalid assignment", n))
	}
	return sm(n, v)
}
