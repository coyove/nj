package parser

import (
	"fmt"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	r := strings.NewReader("local a = b, c")
	lexer := &Lexer{NewScanner(r, "zzz"), nil, false, Token{Str: ""}, TNil}
	yyParse(lexer)
	chunk := lexer.Stmts
	for _, s := range chunk {
		fmt.Println(s.(*LocalAssignStmt).Names)
	}
}
