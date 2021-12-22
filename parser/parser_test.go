package parser

import (
	"testing"
)

func TestHashString(t *testing.T) {
	t.Log(ParseJSON(`{A=
		a}`))
}
