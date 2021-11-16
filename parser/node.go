package parser

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unsafe"
)

type Position struct {
	Source string
	Line   uint32
	Column uint32
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

func (t *Token) String() string {
	return t.Str
}

type Node struct {
	typ     byte
	pad     byte
	Addr    uint16
	symLine uint32
	ptr     unsafe.Pointer
}

func Sym(tok Token) Node {
	return Node{typ: SYM, symLine: tok.Pos.Line, ptr: unsafe.Pointer(&tok.Str)}
}

func Addr(yx uint16) Node { return Node{typ: ADDR, Addr: yx} }

func staticSym(s string) Node {
	return Node{typ: SYM, ptr: unsafe.Pointer(&s)}
}

func Str(s string) Node { return Node{typ: STR, ptr: unsafe.Pointer(&s)} }

func Num(v string) Node {
	i, err := strconv.ParseInt(v, 0, 64)
	if err == nil {
		return Int(i)
	}
	if err.(*strconv.NumError).Err == strconv.ErrRange {
		i, err := strconv.ParseUint(v, 0, 64)
		if err == nil {
			return Int(int64(i))
		}
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		panic("invalid number format: \"" + v + "\"")
	}
	return Float(f)
}

func Int(v int64) (n Node) {
	*(*int64)(unsafe.Pointer(&n)) = v
	n.ptr = intNode
	return
}

func Float(v float64) (n Node) {
	if float64(int64(v)) == v {
		return Int(int64(v))
	}
	*(*float64)(unsafe.Pointer(&n)) = v
	n.ptr = floatNode
	return
}

func (n Node) Type() byte {
	switch n.ptr {
	case intNode:
		return INT
	case floatNode:
		return FLOAT
	default:
		switch n.typ {
		case STR, SYM, ADDR, NODES:
			return n.typ
		}
	}
	return INVALID
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
	return *(*string)(n.ptr)
}

func (n Node) Sym() string {
	if n.Type() != SYM {
		return ""
	}
	return *(*string)(n.ptr)
}

func (n Node) Int() int64 { return *(*int64)(unsafe.Pointer(&n)) }

func (n Node) Float() float64 { return *(*float64)(unsafe.Pointer(&n)) }

func (n Node) Num() (float64, int64, bool) {
	if n.Type() == INT {
		return float64(n.Int()), n.Int(), true
	}
	return n.Float(), int64(n.Float()), false
}

func (n Node) IsNegativeNumber() bool {
	if n.Type() == FLOAT {
		return n.Float() < 0
	}
	return n.Int() < 0
}

func (n Node) Nodes() []Node {
	if n.Type() != NODES {
		return nil
	}
	return *(*[]Node)(n.ptr)
}

func Nodes(args ...Node) Node {
	if len(args) == 3 {
		op := args[0].Sym()
		a, b := args[1], args[2]
		if op == AAdd && a.Type() == STR && b.Type() == STR {
			return Str(a.Str() + b.Str())
		}
		if a.IsNum() && b.IsNum() {
			switch op {
			case AAdd:
				af, ai, aIsInt := a.Num()
				bf, bi, bIsInt := b.Num()
				if aIsInt && bIsInt {
					return Int(ai + bi)
				} else {
					return Float(af + bf)
				}
			case ASub:
				af, ai, aIsInt := a.Num()
				bf, bi, bIsInt := b.Num()
				if aIsInt && bIsInt {
					return Int(ai - bi)
				} else {
					return Float(af - bf)

				}
			case AMul:
				af, ai, aIsInt := a.Num()
				bf, bi, bIsInt := b.Num()
				if aIsInt && bIsInt {
					return Int(ai * bi)
				} else {
					return Float(af * bf)
				}
			case ADiv:
				af, ai, aIsInt := a.Num()
				bf, bi, bIsInt := b.Num()
				if aIsInt && bIsInt && ai%bi == 0 {
					return Int(ai / bi)
				} else {
					return Float(af / bf)
				}
			}
		}
	}
	return Node{ptr: unsafe.Pointer(&args), typ: NODES}
}

func (n Node) At(tok Token) Node {
	switch n.Type() {
	case SYM:
		n.symLine = tok.Pos.Line
	case NODES:
		c := n.Nodes()
		if len(c) > 0 {
			c[0] = c[0].At(tok)
		}
	}
	return n
}

func (n Node) Pos() Position {
	switch n.Type() {
	case SYM:
		return Position{Line: n.symLine}
	case NODES:
		c := n.Nodes()
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
		io.WriteString(w, fmt.Sprintf("%s/%d", n.Sym(), n.symLine))
	default:
		io.WriteString(w, n.String())
	}
}

func (n Node) String() string {
	switch n.Type() {
	case INT:
		return strconv.FormatInt(n.Int(), 10)
	case FLOAT:
		return strconv.FormatFloat(n.Float(), 'f', -1, 64)
	case STR:
		return strconv.Quote(n.Str())
	case SYM:
		return fmt.Sprintf("%s/%d", n.Sym(), n.symLine)
	case NODES:
		buf := make([]string, len(n.Nodes()))
		for i, a := range n.Nodes() {
			buf[i] = a.String()
		}
		return "[" + strings.Join(buf, " ") + "]"
	case ADDR:
		return "0x" + strconv.FormatInt(int64(n.Addr), 16)
	default:
		return fmt.Sprintf("DEBUG: %#v", n)
	}
}

func (n Node) append(n2 ...Node) Node {
	*(*[]Node)(n.ptr) = append(n.Nodes(), n2...)
	return n
}

func (n Node) moveLoadStore(sm func(Node, Node) Node, v Node) Node {
	if len(n.Nodes()) == 3 {
		if s := n.Nodes()[0].Sym(); s == ALoad || s == ALoadStatic {
			return __store(n.Nodes()[1], n.Nodes()[2], v)
		}
	}
	if n.Type() != SYM {
		panic(fmt.Sprintf("DEBUG: %v invalid assignment", n))
	}
	return sm(n, v)
}
