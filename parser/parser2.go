package parser

import (
	"fmt"
	"sync/atomic"

	"github.com/coyove/nj/bas"
	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

var (
	emptyBreak = &BreakContinue{Break: true}
	emptyProg  = &Prog{}
	zero       = Primitive(bas.Int(0))
	one        = Primitive(bas.Int(1))
)

var (
	Sa   = &Symbol{Name: "a!"}
	SNil = &Symbol{Name: "nil"}
)

func (lex *Lexer) pFunc(method bool, name Token, args Node, stats Node, pos Token) Node {
	namev := bas.Str(name.Str)
	if lex.scanner.functions.Contains(namev) {
		lex.Error(fmt.Sprintf("function %s already existed", name.Str))
	}
	lex.scanner.functions.Set(namev, bas.Nil)

	// funcname := Sym(name)
	lex.pFindTailCall(stats)

	f := &Function{Name: name.Str, Body: stats, Line: pos.Line()}
	switch vargs := args.(type) {
	case IdentVarargExpandList:
		f.Args, f.Vararg, f.VargExpand = append(vargs.IdentList, randomVarname()), true, vargs.Expand
		for i := range vargs.Expand {
			lex.scanner.constants.Set(bas.Int(i), bas.Nil)
		}
	case IdentVarargList:
		f.Args, f.Vararg = vargs.IdentList, true
	default:
		f.Args = args.(IdentList)
	}
	if method {
		return f
	}
	return f //  &Assign{funcname, f, pos.Line()}
}

func __markupFuncName(recv, name Token) Token {
	name.Str = recv.Str + "." + name.Str
	return name
}

var lambdaIndex int64

func __markupLambdaName(lambda Token) Token {
	lambda.Str = fmt.Sprintf("lambda.%d.%d", lambda.Pos.Line, atomic.AddInt64(&lambdaIndex, 1))
	return lambda
}

func (lex *Lexer) pFindTailCall(stat Node) {
	switch v := stat.(type) {
	case *Call:
		v.Op = typ.OpTailCall
	case *Prog:
		for i, stat := range v.Stats {
			if i == len(v.Stats)-1 {
				lex.pFindTailCall(stat)
			} else if b, ok := v.Stats[i+1].(*Unary); ok && b.Op == typ.OpRet && b.A == SNil {
				lex.pFindTailCall(stat)
			}
		}
	}
}

func (lex *Lexer) pLoop(body ...Node) Node {
	return &Loop{emptyProg, lex.pProg(false, body...)}
}

func (lex *Lexer) pForRange(v Token, start, end, step, body Node, pos Token) Node {
	forVar := Sym(v)
	if v, ok := step.(Primitive); ok && bas.Value(v).IsNumber() { // step is a static number, easy case
		var cmp Node
		if bas.Value(v).Less(bas.Int(0)) {
			cmp = lex.pBinary(typ.OpLess, end, forVar, pos)
		} else {
			cmp = lex.pBinary(typ.OpLess, forVar, end, pos)
		}
		return lex.pProg(true,
			&Declare{forVar, start, pos.Line()},
			&Loop{
				lex.pBinary(typ.OpInc, forVar, step, pos),
				&If{cmp,
					lex.pProg(false, body, lex.pBinary(typ.OpInc, forVar, step, pos)),
					emptyBreak,
				},
			},
		)
	}
	return lex.pProg(true,
		&Declare{forVar, start, pos.Line()},
		&Loop{
			lex.pBinary(typ.OpInc, forVar, step, pos),
			&If{
				lex.pBinary(typ.OpLess, lex.Int(0), step, pos),
				&If{
					lex.pBinary(typ.OpLess, forVar, end, pos), // +step
					lex.pProg(false, body, lex.pBinary(typ.OpInc, forVar, step, pos)),
					emptyBreak,
				},
				&If{
					lex.pBinary(typ.OpLess, end, forVar, pos), // -step
					lex.pProg(false, body, lex.pBinary(typ.OpInc, forVar, step, pos)),
					emptyBreak,
				},
			},
		},
	)
}

func (lex *Lexer) pForIn(key, value Token, expr, body Node, pos Token) Node {
	k, v, subject, kv := Sym(key), Sym(value), randomVarname(), randomVarname()
	return lex.pProg(true,
		&Declare{subject, expr, pos.Line()},
		&Declare{k, SNil, pos.Line()},
		&Declare{v, SNil, pos.Line()},
		&Declare{kv, SNil, pos.Line()},
		&Loop{
			emptyProg,
			lex.pProg(false,
				&Assign{kv, lex.pBinary(typ.OpNext, subject, kv, pos), pos.Line()},
				&Tenary{typ.OpLoad, kv, lex.Int(0), k, pos.Line()},
				&Tenary{typ.OpLoad, kv, lex.Int(1), v, pos.Line()},
				&If{
					lex.pBinary(typ.OpEq, k, SNil, pos),
					emptyBreak,
					body,
				},
			),
		},
	)
}

func (lex *Lexer) pDeclareAssign(dest []Node, src ExprList, assign bool, pos Token) Node {
	if len(src) == 1 && len(dest) > 1 {
		tmp := randomVarname()
		p := (&Prog{}).Append(&Declare{tmp, src[0], pos.Line()})
		for i, ident := range dest {
			op := &Tenary{typ.OpLoad, tmp, lex.Int(int64(i)), Address(typ.RegA), pos.Line()}
			if assign {
				p.Append(assignLoadStore(ident, op, pos))
			} else {
				p.Append(&Declare{ident.(*Symbol), op, pos.Line()})
			}
		}
		return p
	}
	if len(dest) == 1 && len(src) == 1 {
		if assign {
			return assignLoadStore(dest[0], src[0], pos)
		}
		return &Declare{dest[0].(*Symbol), src[0], pos.Line()}
	}
	res := &Prog{}
	if !assign {
		for i, v := range dest {
			if i >= len(src) {
				res.Append(&Declare{v.(*Symbol), SNil, pos.Line()})
			} else {
				res.Append(&Declare{v.(*Symbol), src[i], pos.Line()})
			}
		}
	} else {
		if len(dest) != len(src) {
			lex.Error(fmt.Sprintf("unmatched number of assignments"))
		}
		// a0, ..., an = b0, ..., bn
		var tmp Release
		for i := range dest {
			x := randomVarname()
			tmp = append(tmp, x)
			res.Append(&Assign{x, src[i], pos.Line()})
		}
		for i, v := range dest {
			res.Append(assignLoadStore(v, tmp[i], pos))
		}
		res.Append(tmp)
	}
	return res
}

func (lex *Lexer) pArray(list, arg Node) Node {
	if lex.scanner.jsonMode {
		if list != nil {
			lex.pSimpleJSON(list).Native().Append(lex.pSimpleJSON(arg))
			return list
		}
		return JValue(bas.Array(lex.pSimpleJSON(arg)))
	}
	if list != nil {
		return append(list.(ExprList), arg)
	}
	return ExprList([]Node{arg})
}

func (lex *Lexer) pObject(list, k, v Node) Node {
	if lex.scanner.jsonMode {
		if list != nil {
			lex.pSimpleJSON(list).Object().Set(lex.pSimpleJSON(k), lex.pSimpleJSON(v))
			return list
		}
		o := bas.NewObject(0)
		o.Set(lex.pSimpleJSON(k), lex.pSimpleJSON(v))
		return JValue(o.ToValue())
	}
	if list != nil {
		return append(list.(ExprAssignList), [2]Node{k, v})
	}
	return append((ExprAssignList)(nil), [2]Node{k, v})
}

func (lex *Lexer) pEmptyArray() Node {
	if lex.scanner.jsonMode {
		return JValue(bas.Array())
	}
	return ExprList(nil)
}

func (lex *Lexer) pEmptyObject() Node {
	if lex.scanner.jsonMode {
		return JValue(bas.NewObject(0).ToValue())
	}
	return ExprAssignList(nil)
}

func randomVarname() *Symbol {
	return &Symbol{Name: internal.Unnamed()}
}
