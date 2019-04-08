package rpn

import (
	"fmt"
	"io"

	"github.com/coyove/potatolang/parser"
)

func Parse(r io.Reader) (*parser.Node, error) {
	return parse(NewReader(r))
}

func parse(r *Reader) (*parser.Node, error) {
	node := parser.CNode(parser.ANodeS("chain"))

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
			n, err := parse(r)
			if err != nil {
				return nil, err
			}
			node.Cappend(n)
			continue
		}

		node.Cappend(parser.ANodeS(fmt.Sprintf("%v", tok.Name)).SetPos(parser.Meta{Line: tok.line + 1, Column: tok.col}))
	}

	return nil, fmt.Errorf("unexpected end of input")
}
