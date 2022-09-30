package parser

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"unsafe"

	"github.com/coyove/nj/bas"
	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

type Token struct {
	Type uint32
	Str  string
	Pos  Position
}

func (t *Token) Line() uint32 {
	return t.Pos.Line
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

type Node interface {
	Dump(io.Writer)
}

type GetLine interface {
	Node
	GetLine() (string, int)
}

type Symbol struct {
	Name string
	Line uint32
}

func (s *Symbol) Dump(w io.Writer) {
	if s.Line == 0 {
		internal.WriteString(w, s.Name)
	} else {
		internal.WriteString(w, fmt.Sprintf("%s/%d", s.Name, s.Line))
	}
}

func (s *Symbol) GetLine() (string, int) {
	return s.Name, int(s.Line)
}

type Address uint16

func (s Address) Dump(w io.Writer) {
	fmt.Fprintf(w, "@%08x", s)
}

type Prog struct {
	DoBlock bool
	Stats   []Node
}

func (p *Prog) Append(stat Node) *Prog {
	if p2, ok := stat.(*Prog); ok && !p2.DoBlock {
		p.Stats = append(p.Stats, p2.Stats...)
	} else {
		p.Stats = append(p.Stats, stat)
	}
	return p
}

func (p *Prog) Dump(w io.Writer) {
	if p == nil {
		internal.WriteString(w, "(prog)")
		return
	}
	internal.WriteString(w, internal.IfStr(p.DoBlock, "(do", "(prog"))
	for _, stat := range p.Stats {
		if stat != nil {
			internal.WriteString(w, " ")
			stat.Dump(w)
		}
	}
	internal.WriteString(w, ")")
}

func (p *Prog) String() string {
	buf := &bytes.Buffer{}
	p.Dump(buf)
	return buf.String()
}

type LoadConst struct {
	Table bas.Map
	Funcs bas.Map
}

func (p *LoadConst) Dump(w io.Writer) {
	internal.WriteString(w, "(loadconst")
	p.Table.Foreach(func(k bas.Value, v *bas.Value) bool {
		internal.WriteString(w, " ")
		k.Stringify(w, typ.MarshalToJSON)
		return true
	})
	p.Funcs.Foreach(func(k bas.Value, v *bas.Value) bool {
		internal.WriteString(w, " f:")
		k.Stringify(w, typ.MarshalToString)
		return true
	})
	internal.WriteString(w, ")")
}

type Primitive bas.Value

func (p Primitive) Dump(w io.Writer) {
	bas.Value(p).Stringify(w, typ.MarshalToJSON)
}

func (p Primitive) Value() bas.Value {
	return bas.Value(p)
}

type JValue bas.Value

func (p JValue) Dump(w io.Writer) {
	bas.Value(p).Stringify(w, typ.MarshalToJSON)
}

type If struct {
	Cond        Node
	True, False Node
}

func (p *If) Dump(w io.Writer) {
	internal.WriteString(w, "(if ")
	p.Cond.Dump(w)
	internal.WriteString(w, " ")
	p.True.Dump(w)
	if p.False != nil {
		internal.WriteString(w, " ")
		p.False.Dump(w)
	}
	internal.WriteString(w, ")")
}

type Unary struct {
	Op   byte
	A    Node
	Line uint32
}

func (b *Unary) Dump(w io.Writer) {
	fmt.Fprintf(w, "(%s/%d ", typ.UnaryOpcode[b.Op], b.Line)
	b.A.Dump(w)
	internal.WriteString(w, ")")
}

func (b *Unary) GetLine() (string, int) {
	return typ.UnaryOpcode[b.Op], int(b.Line)
}

type Binary struct {
	A, B Node
	Op   byte
	Line uint32
}

func (b *Binary) Dump(w io.Writer) {
	fmt.Fprintf(w, "(%s/%d ", typ.BinaryOpcode[b.Op], b.Line)
	b.A.Dump(w)
	internal.WriteString(w, " ")
	b.B.Dump(w)
	internal.WriteString(w, ")")
}

func (b *Binary) GetLine() (string, int) {
	return typ.BinaryOpcode[b.Op], int(b.Line)
}

type Bitwise struct {
	A, B Node
	Op   string
	Line uint32
}

func (b *Bitwise) Dump(w io.Writer) {
	fmt.Fprintf(w, "(bit%s/%d ", b.Op, b.Line)
	b.A.Dump(w)
	internal.WriteString(w, " ")
	b.B.Dump(w)
	internal.WriteString(w, ")")
}

func (b *Bitwise) GetLine() (string, int) {
	return "bit" + b.Op, int(b.Line)
}

type Declare struct {
	Name  *Symbol
	Value Node
	Line  uint32
}

func (p *Declare) Dump(w io.Writer) {
	fmt.Fprintf(w, "(declare/%d %s ", p.Line, p.Name.Name)
	p.Value.Dump(w)
	internal.WriteString(w, ")")
}

func (b *Declare) GetLine() (string, int) {
	return b.Name.Name, int(b.Line)
}

type Assign Declare

func (p *Assign) Dump(w io.Writer) {
	fmt.Fprintf(w, "(assign/%d %s ", p.Line, p.Name.Name)
	p.Value.Dump(w)
	internal.WriteString(w, ")")
}

func (b *Assign) GetLine() (string, int) {
	return b.Name.Name, int(b.Line)
}

type Release []*Symbol

func (p Release) Dump(w io.Writer) {
	internal.WriteString(w, "(release")
	for _, name := range p {
		internal.WriteString(w, " ")
		name.Dump(w)
	}
	internal.WriteString(w, ")")
}

type Tenary struct {
	Op      byte
	A, B, C Node
	Line    uint32
}

func (p *Tenary) Dump(w io.Writer) {
	fmt.Fprintf(w, "(%s/%d ", typ.TenaryOpcode[p.Op], p.Line)
	p.A.Dump(w)
	internal.WriteString(w, " ")
	p.B.Dump(w)
	internal.WriteString(w, " ")
	p.C.Dump(w)
	internal.WriteString(w, ")")
}

func (b *Tenary) GetLine() (string, int) {
	return typ.TenaryOpcode[b.Op], int(b.Line)
}

type IdentList []Node

func (p IdentList) Dump(w io.Writer) {
	internal.WriteString(w, "(identlist ")
	for _, n := range p {
		internal.WriteString(w, " ")
		n.Dump(w)
	}
	internal.WriteString(w, ")")
}

type IdentVarargList struct {
	IdentList
}

type ExprList []Node

func (p ExprList) Dump(w io.Writer) {
	internal.WriteString(w, "(list")
	for _, n := range p {
		internal.WriteString(w, " ")
		n.Dump(w)
	}
	internal.WriteString(w, ")")
}

type DeclList []Node

func (p DeclList) Dump(w io.Writer) {}

type ExprAssignList [][2]Node

func (p *ExprAssignList) ExpandAsExprList() (tmp ExprList) {
	*(*[3]int)(unsafe.Pointer(&tmp)) = *(*[3]int)(unsafe.Pointer(p))
	(*(*[3]int)(unsafe.Pointer(&tmp)))[1] = len(*p) * 2
	(*(*[3]int)(unsafe.Pointer(&tmp)))[2] = cap(*p) * 2
	return
}

func (p ExprAssignList) Dump(w io.Writer) {
	internal.WriteString(w, "(assignlist")
	for _, n2 := range p {
		internal.WriteString(w, " (")
		n2[0].Dump(w)
		internal.WriteString(w, " ")
		n2[1].Dump(w)
		internal.WriteString(w, ")")
	}
	internal.WriteString(w, ")")
}

type And struct {
	A, B Node
}

func (p And) Dump(w io.Writer) {
	internal.WriteString(w, "(and ")
	p.A.Dump(w)
	internal.WriteString(w, " ")
	p.B.Dump(w)
	internal.WriteString(w, ")")
}

type Or And

func (p Or) Dump(w io.Writer) {
	internal.WriteString(w, "(or ")
	p.A.Dump(w)
	internal.WriteString(w, " ")
	p.B.Dump(w)
	internal.WriteString(w, ")")
}

type Loop struct {
	Continue Node
	Body     Node
}

func (p *Loop) Dump(w io.Writer) {
	internal.WriteString(w, "(loop ")
	p.Continue.Dump(w)
	internal.WriteString(w, " ")
	p.Body.Dump(w)
	internal.WriteString(w, ")")
}

type Call struct {
	Op     byte
	Callee Node
	Args   ExprList
	Vararg bool
	Line   uint32
}

func (p *Call) Dump(w io.Writer) {
	internal.WriteString(w, "(")
	internal.WriteString(w, internal.IfStr(p.Op == typ.OpTailCall, "tail", ""))
	internal.WriteString(w, "call")
	internal.WriteString(w, internal.IfStr(p.Vararg, "varg", ""))
	fmt.Fprintf(w, "/%d ", p.Line)
	p.Callee.Dump(w)
	internal.WriteString(w, " ")
	p.Args.Dump(w)
	internal.WriteString(w, ")")
}

func (b *Call) GetLine() (string, int) {
	buf := &bytes.Buffer{}
	b.Callee.Dump(buf)
	return buf.String(), int(b.Line)
}

type Function struct {
	Name   string
	Args   IdentList
	Body   Node
	Vararg bool
	Line   uint32
}

func (p *Function) Dump(w io.Writer) {
	fmt.Fprintf(w, "(function/%d %s ", p.Line, p.Name)
	p.Args.Dump(w)
	internal.WriteString(w, " ")
	p.Body.Dump(w)
	internal.WriteString(w, ")")
}

func (b *Function) GetLine() (string, int) {
	return b.Name, int(b.Line)
}

type GotoLabel struct {
	Label string
	Goto  bool
	Line  uint32
}

func (p *GotoLabel) Dump(w io.Writer) {
	fmt.Fprintf(w, "(%s/%d %s)", internal.IfStr(p.Goto, "goto", "label"), p.Line, p.Label)
}

func (p *GotoLabel) GetLine() (string, int) {
	return p.Label, int(p.Line)
}

type BreakContinue struct {
	Break bool
	Line  uint32
}

func (p *BreakContinue) Dump(w io.Writer) {
	fmt.Fprintf(w, "(%s/%d)", internal.IfStr(p.Break, "break", "continue"), p.Line)
}

func (p *BreakContinue) GetLine() (string, int) {
	return internal.IfStr(p.Break, "break", "continue"), int(p.Line)
}

func Sym(tok Token) *Symbol {
	return &Symbol{tok.Str, tok.Line()}
}

func (lex *Lexer) Str(s string) Primitive {
	x := Primitive(bas.Str(s))
	if !lex.scanner.jsonMode {
		lex.scanner.constants.Set(bas.Value(x), bas.Nil)
	}
	return x
}

func (lex *Lexer) Num(v string) Primitive {
	f, i, isInt, err := internal.ParseNumber(v)
	internal.PanicErr(err)
	if isInt {
		return lex.Int(i)
	}
	return lex.Float(f)
}

func (lex *Lexer) Int(v int64) Primitive {
	x := Primitive(bas.Int64(v))
	if v != 0 && v != 1 && !lex.scanner.jsonMode {
		lex.scanner.constants.Set(bas.Value(x), bas.Nil)
	}
	return x
}

func (lex *Lexer) IntBool(b bool) (n Node) {
	if b {
		return lex.Int(1)
	}
	return lex.Int(0)
}

func (lex *Lexer) Float(v float64) Primitive {
	x := Primitive(bas.Float64(v))
	if !lex.scanner.jsonMode {
		lex.scanner.constants.Set(bas.Value(x), bas.Nil)
	}
	return x
}

func IsInt16(n Node) int {
	if isInt64(n) {
		a := bas.Value(n.(Primitive)).UnsafeInt64()
		if a >= math.MinInt16 && a <= math.MaxInt16 {
			if -a >= math.MinInt16 && -a <= math.MaxInt16 {
				return 2
			}
			return 1
		}
	}
	return 0
}

func pUnary(op byte, a Node, pos Token) Node {
	return &Unary{Op: op, A: a, Line: pos.Line()}
}

func (lex *Lexer) pProg(do bool, n ...Node) Node {
	return &Prog{do, n}
}

func (lex *Lexer) pBinary(op byte, a, b Node, pos Token) Node {
	if op == typ.OpAdd && isString(a) && isString(b) {
		as := bas.Value(a.(Primitive)).Str()
		bs := bas.Value(b.(Primitive)).Str()
		return Primitive(lex.Str(as + bs))
	}
	if isNumber(a) && isNumber(b) {
		a := bas.Value(a.(Primitive))
		b := bas.Value(b.(Primitive))
		switch op {
		case typ.OpAdd:
			if a.IsInt64() && b.IsInt64() {
				return lex.Int(a.Int64() + b.Int64())
			}
			return lex.Float(a.Float64() + b.Float64())
		case typ.OpSub:
			if a.IsInt64() && b.IsInt64() {
				return lex.Int(a.Int64() - b.Int64())
			}
			return lex.Float(a.Float64() - b.Float64())
		case typ.OpMul:
			if a.IsInt64() && b.IsInt64() {
				return lex.Int(a.Int64() * b.Int64())
			}
			return lex.Float(a.Float64() * b.Float64())
		case typ.OpDiv:
			return lex.Float(a.Float64() / b.Float64())
		case typ.OpIDiv:
			return lex.Int(a.Int64() / b.Int64())
		case typ.OpMod:
			return lex.Int(a.Int64() % b.Int64())
		case typ.OpLessEq:
			if a.Equal(b) {
				return lex.IntBool(true)
			}
			fallthrough
		case typ.OpLess:
			if a.IsInt64() && b.IsInt64() {
				return lex.IntBool(a.Int64() < b.Int64())
			}
			return lex.IntBool(a.Float64() < b.Float64())
		case typ.OpEq:
			return lex.IntBool(a.Equal(b))
		case typ.OpNeq:
			return lex.IntBool(!a.Equal(b))
		}
	}
	return &Binary{Op: op, A: a, B: b, Line: pos.Line()}
}

func (lex *Lexer) pBitwise(op string, a, b Node, pos Token) Node {
	if isInt64(a) && isInt64(b) {
		a := bas.Value(a.(Primitive))
		b := bas.Value(b.(Primitive))
		switch op {
		case "and":
			return lex.Int(a.Int64() & b.Int64())
		case "or":
			return lex.Int(a.Int64() | b.Int64())
		case "xor":
			return lex.Int(a.Int64() ^ b.Int64())
		case "lsh":
			return lex.Int(a.Int64() << b.Int64())
		case "rsh":
			return lex.Int(a.Int64() >> b.Int64())
		case "ursh":
			return lex.Int(int64(uint64(a.Int64()) >> b.Int64()))
		}
	}
	return &Bitwise{Op: op, A: a, B: b, Line: pos.Line()}
}

func assignLoadStore(n, v Node, pos Token) Node {
	if load, ok := n.(*Tenary); ok && load.Op == typ.OpLoad {
		return &Tenary{typ.OpStore, load.A, load.B, v, pos.Line()}
	}
	return &Assign{n.(*Symbol), v, pos.Line()}
}

func (lex *Lexer) pSimpleJSON(n Node) bas.Value {
	switch v := n.(type) {
	case JValue:
		return bas.Value(v)
	case Primitive:
		return bas.Value(v)
	case *Symbol:
		switch v.Name {
		case "true":
			return bas.True
		case "false":
			return bas.False
		case "null":
			return bas.Nil
		}
	}
	buf := bytes.NewBufferString("unexpected json symbol: ")
	n.Dump(buf)
	lex.Error(buf.String())
	return bas.Nil
}

func isString(n Node) bool {
	n2, ok := n.(Primitive)
	return ok && bas.Value(n2).IsString()
}

func isNumber(n Node) bool {
	n2, ok := n.(Primitive)
	return ok && bas.Value(n2).IsNumber()
}

func isInt64(n Node) bool {
	n2, ok := n.(Primitive)
	return ok && bas.Value(n2).IsInt64()
}
