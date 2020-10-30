package parser

import (
	"strings"
	"testing"
)

func TestHashString(t *testing.T) {
	t.Log(Parse(strings.NewReader(`
a = 1e3
-- local b, ... = a()
`), ""))
}
