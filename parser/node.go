package parser

import (
	"fmt"
	"io"
	"math/big"
	"reflect"
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
	Value interface{}
}

func SymTok(tok Token) Node { return Node{Symbol{tok.Pos, tok.Str}} }

func Sym(s string) Node { return Node{Symbol{Text: s}} }

func Num(arg interface{}) Node {
	n := Node{}
	switch x := arg.(type) {
	case *big.Float:
		v, acc := x.Int64()
		if acc == big.Exact {
			n.Value = v
		} else {
			n.Value, _ = x.Float64()
		}
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

func (n Node) Valid() bool {
	t := n.Type()
	return t == SYM || t == NUM || t == STR || t == CPL || t == ADR
}

func (n Node) Type() uintptr {
	t := interfaceType(n.Value)
	if t == numINT {
		t = NUM
	}
	return t
}

func (n Node) TypeName() string { return reflect.TypeOf(n.Value).String() }

func (n Node) Cpl() []Node { a, _ := n.Value.([]Node); return a }

func (n Node) CplAppend(na ...Node) Node {
	n.Value = append(n.Cpl(), na...)
	return n
}

func (n Node) CplPrepend(n2 Node) Node {
	n.Value = append([]Node{n2}, n.Cpl()...)
	return n
}

func (n Node) CplIndex(i int) Node { return n.Value.([]Node)[i] }

func (n Node) Str() string { a, _ := n.Value.(string); return a }

func (n Node) Sym() Symbol { a, _ := n.Value.(Symbol); return a }

func (n Node) Num() (float64, int64) {
	a, ok := n.Value.(float64)
	if ok {
		return a, 0
	}
	return 0, n.Value.(int64)
}

func (n Node) Num1() *big.Float {
	a, ok := n.Value.(float64)
	if ok {
		return (&big.Float{}).SetFloat64(a)
	}
	return (&big.Float{}).SetInt64(n.Value.(int64))
}

func Cpl(args ...Node) Node {
	if len(args) == 3 {
		op, _ := args[0].Value.(Symbol)
		a, b := args[1], args[2]
		if op.Equals(AConcat) && a.Type() == STR && b.Type() == STR {
			return Node{a.Str() + b.Str()}
		}
		if a.Type() == NUM && b.Type() == NUM {
			switch v1, v2 := a.Num1(), b.Num1(); op.Text {
			case AAdd.Text:
				return Num(v1.Add(v1, v2))
			case ASub.Text:
				return Num(v1.Sub(v1, v2))
			case AMul.Text:
				return Num(v1.Mul(v1, v2))
			case ADiv.Text:
				return Num(v1.Quo(v1, v2))
			}
		}
	}
	if len(args) == 2 {
		if args[0].Sym().Equals(AUnm) && args[1].Type() == NUM {
			f, i := args[1].Num()
			if i == 0 {
				return Node{-f}
			}
			return Node{-i}
		}
	}
	return Node{args}
}

func (n Node) SetPos(pos Position) Node {
	switch n.Type() {
	case SYM:
		s := n.Value.(Symbol)
		if s.Position.Source == "" {
			s.Position = pos
		}
		return Node{s}
	case CPL:
		c := n.Cpl()
		if len(c) == 0 {
			return n
		}
		c[0] = c[0].SetPos(pos)
		return Node{c}
	default:
		return n
	}
}

func (n Node) Pos() Position {
	switch n.Type() {
	case SYM:
		return n.Value.(Symbol).Position
	case CPL:
		c := n.Cpl()
		if len(c) == 0 {
			panic("Pos()")
		}
		return c[0].Pos()
	default:
		panic("Pos()")
	}
}

func StringToNumber(arg string) (interface{}, error) {
	i, err := strconv.ParseInt(arg, 0, 64)
	if err == nil {
		if int64(float64(i)) == i {
			return float64(i), nil
		}
		return i, nil
	}
	return strconv.ParseFloat(arg, 64)
}

func (n Node) Dump(w io.Writer, ident string) {
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

func (n Node) String() string {
	switch n.Type() {
	case NUM:
		return strconv.FormatFloat(n.Value.(float64), 'f', -1, 64)
	case STR:
		return strconv.Quote(n.Value.(string))
	case SYM:
		return n.Value.(Symbol).String()
	case CPL:
		buf := make([]string, len(n.Cpl()))
		for i, a := range n.Cpl() {
			buf[i] = a.String()
		}
		return "[" + strings.Join(buf, " ") + "]"
	case ADR:
		return "0x" + strconv.FormatInt(int64(n.Value.(uint16)), 16)
	default:
		return fmt.Sprintf("shouldn't happen: %v", n.Value)
	}
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
	return len(n.Cpl()) > 0 && n.CplIndex(0).Sym().Equals(ACall)
}

func (n Node) moveLoadStore(sm func(Node, Node) Node, v Node) Node {
	if len(n.Cpl()) == 3 && n.CplIndex(0).Sym().Equals(ALoad) {
		return __store(n.CplIndex(1), n.CplIndex(2), v)
	}
	if n.Type() != SYM {
		panic(fmt.Sprintf("%v: invalid assignment", n))
	}
	return sm(n, v)
}
