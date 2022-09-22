package parser

import (
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/coyove/nj/bas"
	"github.com/coyove/nj/internal"
)

type Token struct {
	Type uint32
	Str  string
	Pos  Position
}

func (t *Token) String() string {
	return t.Str
}

type Position struct {
	Source string
	Line   uint32
	Column uint32
}

func (pos *Position) String() string {
	return fmt.Sprintf("%s:%d:%d", pos.Source, pos.Line, pos.Column)
}

type Node struct {
	bas.Value
	NodeType byte
	SymLine  uint32
}

func staticSym(s string) Node {
	return Node{NodeType: SYM, Value: bas.Str(s)}
}

func jsonValue(v bas.Value) Node {
	return Node{NodeType: JSON, Value: v}
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

func IntBool(b bool) (n Node) {
	if b {
		return Node{NodeType: INT, Value: bas.Int64(1)}
	}
	return Node{NodeType: INT, Value: bas.Int64(0)}
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

func (n Node) numSign() (isNum bool, isNeg bool) {
	switch n.Type() {
	case FLOAT:
		return true, n.UnsafeFloat64() < 0
	case INT:
		return true, n.UnsafeInt64() < 0
	}
	return false, false
}

func (n Node) IsInt16() int {
	if n.Type() == INT {
		a := n.Value.UnsafeInt64()
		if a >= math.MinInt16 && a <= math.MaxInt16 {
			if -a >= math.MinInt16 && -a <= math.MaxInt16 {
				return 2
			}
			return 1
		}
	}
	return 0
}

func (n Node) Nodes() []Node {
	if n.Type() != NODES {
		return nil
	}
	return n.Native().Unwrap().([]Node)
}

func Nodes(args ...Node) Node {
	if len(args) == 3 {
		op, a, b := args[0].Value, args[1], args[2]
		if op == SAdd.Value && a.Type() == STR && b.Type() == STR {
			return Str(a.Str() + b.Str())
		}
		if (a.Type() == INT || a.Type() == FLOAT) && (b.Type() == INT || b.Type() == FLOAT) {
			switch op {
			case SAdd.Value:
				if a.IsInt64() && b.IsInt64() {
					return Int(a.Int64() + b.Int64())
				}
				return Float(a.Float64() + b.Float64())
			case SSub.Value:
				if a.IsInt64() && b.IsInt64() {
					return Int(a.Int64() - b.Int64())
				}
				return Float(a.Float64() - b.Float64())
			case SMul.Value:
				if a.IsInt64() && b.IsInt64() {
					return Int(a.Int64() * b.Int64())
				}
				return Float(a.Float64() * b.Float64())
			case SDiv.Value:
				return Float(a.Float64() / b.Float64())
			case SIDiv.Value:
				return Int(a.Int64() / b.Int64())
			case SMod.Value:
				return Int(a.Int64() % b.Int64())
			case SLessEq.Value:
				if a.Value.Equal(b.Value) {
					return IntBool(true)
				}
				fallthrough
			case SLess.Value:
				if a.IsInt64() && b.IsInt64() {
					return IntBool(a.Int64() < b.Int64())
				}
				return IntBool(a.Float64() < b.Float64())
			case SEq.Value:
				return IntBool(a.Value.Equal(b.Value))
			case SNeq.Value:
				return IntBool(!a.Value.Equal(b.Value))
			}
		}
		if a.Type() == INT && b.Type() == INT {
			switch op {
			case SBitAnd.Value:
				return Int(a.Int64() & b.Int64())
			case SBitOr.Value:
				return Int(a.Int64() | b.Int64())
			case SBitXor.Value:
				return Int(a.Int64() ^ b.Int64())
			case SBitLsh.Value:
				return Int(a.Int64() << b.Int64())
			case SBitRsh.Value:
				return Int(a.Int64() >> b.Int64())
			case SBitURsh.Value:
				return Int(int64(uint64(a.Int64()) >> b.Int64()))
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
			internal.ShouldNotHappen()
		}
		return c[0].Line()
	}
	internal.ShouldNotHappen()
	return 0
}

func (n Node) Dump(w io.Writer, indent string) {
	io.WriteString(w, indent)
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
				a.Dump(w, "  "+indent)
				if a.Type() != NODES {
					io.WriteString(w, "\n")
				}
			}
			io.WriteString(w, indent)
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
		return "<invalid node>"
	}
}

func (n Node) append(n2 Node) Node {
	n.Native().UnwrapFunc(func(i interface{}) interface{} {
		return append(i.([]Node), n2)
	})
	return n
}

func (n Node) moveLoadStore(v Node) Node {
	if len(n.Nodes()) == 3 {
		if s := n.Nodes()[0].Value; s == SLoad.Value {
			return __store(n.Nodes()[1], n.Nodes()[2], v)
		}
	}
	if n.Type() != SYM {
		internal.ShouldNotHappen(n)
	}
	return __move(n, v)
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
