package parser

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	r := strings.NewReader(`set a = b ;
		(function(b)
		return 1
		return 2
		end)[2]()

		if a > 2 then
		print(1)
		elseif a < 1 then
		a=a+1*3
		print(1)
		end
		`)
	lexer := &Lexer{NewScanner(r, "zzz"), nil, false, Token{Str: ""}, TNil}
	yyParse(lexer)
	buf, _ := json.Marshal(lexer.Stmts)
	b := &bytes.Buffer{}
	json.Indent(b, buf, "", "  ")

	t.Error(b.String())
}
