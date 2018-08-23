package parser

import (
	"bytes"
	"strings"
	"testing"
)

func TestChainFuns(t *testing.T) {
	c, err := Parse(strings.NewReader(`var a =fun a = fun b =
		 fun c = a + b + c`), "mem")
	if err != nil {
		t.Fatal(err)
	}

	b := &bytes.Buffer{}
	c.Dump(b)

	if b.String() != "[chain [chain [set a [func <a> [a ] [chain [ret [func <a> [b ] [chain [ret [func <a> [c ] [chain [ret [+ [+ a b ] c ] ] ] ] ] ] ] ] ] ] ] ] ]" {
		t.Fatal(b.String())
	}
}
