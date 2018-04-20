package parser

import (
	"os"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	r := strings.NewReader(`
		out_println("zzz")
		`)
	lexer := &Lexer{NewScanner(r, "zzz"), nil, false, Token{Str: ""}, TNil}
	yyParse(lexer)
	lexer.Stmts.Dump(os.Stderr)
}
