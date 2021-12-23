package parser

import (
	"fmt"
	"io"
	"strconv"

	"github.com/coyove/nj/bas"
	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

type Token struct {
	Type uint32
	Str  string
	Pos  typ.Position
}

func (t *Token) String() string {
	return t.Str
}

type Node struct {
	bas.Value
	NodeType byte
	SymLine  uint32
}

type TokenNode struct {
	Token
	Node
}

func staticSym(s string) Node {
	return Node{NodeType: SYM, Value: bas.Str(s)}
}

func Sym(tok Token) Node {
	return Node{NodeType: SYM, SymLine: tok.Pos.Line, Value: bas.Str(tok.Str)}
}

func Addr(yx uint16) Node {
	return Node{NodeType: ADDR, Value: bas.Int(int(yx))}
}

func Str(s string) Node {
	return Node{NodeType: STR, Value: bas.Str(s)}
}

func Num(v string) Node {
	f, i, isInt, err := internal.ParseNumber(v)
	internal.PanicErr(err)
	if isInt {
		return Int(i)
	}
	return Float(f)
}

func Int(v int64) (n Node) {
	return Node{NodeType: INT, Value: bas.Int64(v)}
}

func Float(v float64) (n Node) {
	return Node{NodeType: FLOAT, Value: bas.Float64(v)}
}

func (n Node) Type() byte {
	return n.NodeType
}

func (n Node) Valid() bool {
	return n.Type() != INVALID
}

func (n Node) IsNum() bool {
	t := n.Type()
	return t == FLOAT || t == INT
}

func (n Node) Str() string {
	if n.Type() != STR {
		return ""
	}
	return n.Value.Str()
}

func (n Node) Sym() string {
	if n.Type() != SYM {
		return ""
	}
	return n.Value.Str()
}

func (n Node) Num() (float64, int64, bool) {
	if n.Type() == INT {
		return float64(n.Int64()), n.Int64(), true
	}
	return n.Float64(), int64(n.Float64()), false
}

func (n Node) IsNegativeNumber() bool {
	if n.Type() == FLOAT {
		return n.Float64() < 0
	}
	return n.Int64() < 0
}

func (n Node) Nodes() []Node {
	if n.Type() != NODES {
		return nil
	}
	return n.Array().Unwrap().([]Node)
}

func Nodes(args ...Node) Node {
	if len(args) == 3 {
		op := args[0].Sym()
		a, b := args[1], args[2]
		if op == typ.AAdd && a.Type() == STR && b.Type() == STR {
			return Str(a.Str() + b.Str())
		}
		if a.IsNum() && b.IsNum() {
			switch op {
			case typ.AAdd:
				af, ai, aIsInt := a.Num()
				bf, bi, bIsInt := b.Num()
				if aIsInt && bIsInt {
					return Int(ai + bi)
				}
				return Float(af + bf)
			case typ.ASub:
				af, ai, aIsInt := a.Num()
				bf, bi, bIsInt := b.Num()
				if aIsInt && bIsInt {
					return Int(ai - bi)
				}
				return Float(af - bf)
			case typ.AMul:
				af, ai, aIsInt := a.Num()
				bf, bi, bIsInt := b.Num()
				if aIsInt && bIsInt {
					return Int(ai * bi)
				}
				return Float(af * bf)
			case typ.ADiv:
				af, _, _ := a.Num()
				bf, _, _ := b.Num()
				return Float(af / bf)
			case typ.AIDiv:
				_, ai, _ := a.Num()
				_, bi, _ := b.Num()
				return Int(ai / bi)
			}
		}
	}
	return Node{NodeType: NODES, Value: bas.ValueOf(args)}
}

func (n Node) At(tok Token) Node {
	switch n.Type() {
	case SYM:
		n.SymLine = tok.Pos.Line
	case NODES:
		c := n.Nodes()
		if len(c) > 0 {
			c[0] = c[0].At(tok)
		}
	}
	return n
}

func (n Node) Line() uint32 {
	switch n.Type() {
	case SYM:
		return n.SymLine
	case NODES:
		c := n.Nodes()
		if len(c) == 0 {
			panic("Line()")
		}
		return c[0].Line()
	default:
		panic("Line()")
	}
}

func (n Node) Dump(w io.Writer, ident string) {
	io.WriteString(w, ident)
	switch n.Type() {
	case NODES:
		nocpl := true
		for _, a := range n.Nodes() {
			if a.Type() == NODES {
				nocpl = false
				break
			}
		}

		if !nocpl {
			io.WriteString(w, "[\n")
			for _, a := range n.Nodes() {
				a.Dump(w, "  "+ident)
				if a.Type() != NODES {
					io.WriteString(w, "\n")
				}
			}
			io.WriteString(w, ident)
		} else {
			io.WriteString(w, "[")
			for i, a := range n.Nodes() {
				a.Dump(w, "")
				if i < len(n.Nodes())-1 {
					io.WriteString(w, " ")
				}
			}
		}
		io.WriteString(w, "]\n")
	case SYM:
		io.WriteString(w, fmt.Sprintf("%s/%d", n.Sym(), n.SymLine))
	default:
		io.WriteString(w, n.String())
	}
}

func (n Node) String() string {
	switch n.Type() {
	case INT, FLOAT, STR, JSON:
		return n.Value.JSONString()
	case NODES:
		return n.Value.String()
	case SYM:
		return fmt.Sprintf("%s/%d", n.Sym(), n.SymLine)
	case ADDR:
		return "0x" + strconv.FormatInt(n.Int64(), 16)
	default:
		return fmt.Sprintf("DEBUG: %#v", n)
	}
}

func (n Node) append(n2 ...Node) Node {
	x := append(n.Array().Unwrap().([]Node), n2...)
	return Node{NodeType: NODES, Value: bas.ValueOf(x)}
}

func (n Node) moveLoadStore(sm func(Node, Node) Node, v Node) Node {
	if len(n.Nodes()) == 3 {
		if s := n.Nodes()[0].Sym(); s == typ.ALoad {
			return __store(n.Nodes()[1], n.Nodes()[2], v)
		}
	}
	if n.Type() != SYM {
		panic(fmt.Sprintf("DEBUG: %v invalid assignment", n))
	}
	return sm(n, v)
}

func (n Node) simpleJSON(lex *Lexer) bas.Value {
	switch n.Type() {
	case JSON, STR, INT, FLOAT:
		return n.Value
	case SYM:
		switch n.Value.Str() {
		case "true":
			return bas.True
		case "false":
			return bas.False
		case "null":
			return bas.Nil
		}
	}
	lex.Error("unexpected json symbol: " + n.String())
	return bas.Nil
}
