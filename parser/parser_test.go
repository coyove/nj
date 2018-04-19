package parser

import (
	"os"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	r := strings.NewReader(`set a = b ;
		(function()
		return 1
		return 2
		end)[2]()

		if a == true then
		print(1)
		elseif a < 1 then
		a=a+1*3
		print(1)
		end
		`)
	lexer := &Lexer{NewScanner(r, "zzz"), nil, false, Token{Str: ""}, TNil}
	yyParse(lexer)
	lexer.Stmts.Dump(os.Stderr)
}
