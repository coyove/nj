package rpn

import (
	"fmt"
	"io"

	"github.com/coyove/potatolang/parser"
)

const ErrUnexpectedInput = "unexpected input token: %v"

type stack []*Token

func (s *stack) push(t *Token) { *s = append(*s, t) }

func (s *stack) empty() bool { return len(*s) == 0 }

func (s *stack) pop() *Token {
	if len(*s) == 0 {
		panic("empty stack")
	}
	t := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return t
}

type funcargs struct {
	m      map[string]int
	parent *funcargs
}

func newfuncargs(parent *funcargs) *funcargs { return &funcargs{make(map[string]int), parent} }

func Parse(r io.Reader) (*parser.Node, error) {
	return parse(NewReader(r), newfuncargs(nil))
}

func parse(r *Reader, fc *funcargs) (*parser.Node, error) {
	node := parser.CNode()
	s := make(stack, 0)

	for {
		tok, err := r.Token()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			return node, nil
		}

		if tok.Is(T_R) {
			return node, nil
		}

		if tok.Is(T_L) {
			n, err := parse(r, newfuncargs(fc))
			if err != nil {
				return nil, err
			}
			node.Cappend(n)
			continue
		}

		if tok.IsKey("fun") {
			// syntax: fun 'funcname? arg0 arg1 ... argn [ body ]
			ntok, err := r.Token()
			if err != nil {
				return nil, fmt.Errorf(ErrUnexpectedInput, tok)
			}

			var body *parser.Node
			params := parser.CNode()
			funcname := (*parser.Node)(nil)

			if ntok.IsIdent() {
				funcname = ntok.ToNode()
			} else if ntok.IsVar() {
				params.Cappend(ntok.ToNode())
			} else {
				return nil, fmt.Errorf(ErrUnexpectedInput, ntok)
			}

			for {
				vtok, err := r.Token()
				if err != nil {
					return nil, fmt.Errorf(ErrUnexpectedInput, ntok)
				}
				if vtok.Is(T_L) {
					body, err = parse(r, newfuncargs(fc))
					if err != nil {
						return nil, err
					}
					break
				}
				if !vtok.IsVar() {
					return nil, fmt.Errorf(ErrUnexpectedInput, vtok)
				}
				params.Cappend(vtok.ToNode())
			}

			if funcname != nil {
				fc.m[funcname.Value.(string)] = params.Cn()
				node.Cappend(parser.CNode("chain",
					parser.CNode("var", funcname, "nil"),
					parser.CNode("set", funcname,
						parser.CNode("lambda", funcname, params, body))))
			} else {
				node.Cappend(parser.CNode("lambda", parser.ANodeS("<a>"), params, body))
			}
			continue
		}

		if tok.IsKey("if") {
			// syntax: if arg0 arg1 .. argn [ true body ] else? [ false body ]?
			condition := parser.CNode()
			body, falsebody := (*parser.Node)(nil), parser.CNode()
			for {
				vtok, err := r.Token()
				if err != nil {
					return nil, unexpectedToken(tok)
				}
				if vtok.Is(T_L) {
					body, err = parse(r, newfuncargs(fc))
					if err != nil {
						return nil, err
					}
					break
				}
				condition.Cappend(vtok.ToNode())
			}
			elsetok, err := r.Token()
			if err == nil {
				if elsetok.IsKey("else") {
					// read false body
					vtok, err := r.Token()
					if err != nil {
						return nil, unexpectedToken(tok)
					}
					if !vtok.Is(T_L) {
						return nil, unexpectedToken(vtok)
					}
					falsebody, err = parse(r, newfuncargs(fc))
					if err != nil {
						return nil, err
					}
				} else {
					r.UnreadToken(elsetok)
				}
			}

			node.Cappend(parser.CNode("if", condition, body, falsebody))
			continue
		}

		if tok.IsIdent() {
			if s.empty() {
				node.Cappend(parser.CNode("var", tok.ToNode(), "nil"))
			} else {
				v := s.pop()
				node.Cappend(parser.CNode("var", tok.ToNode(), v.ToNode()))
			}
			continue
		}

		s.push(tok)
		node.Cappend(tok.ToNode())
	}

	return nil, fmt.Errorf("unexpected end of input")
}
