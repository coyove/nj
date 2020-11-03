package parser

import (
	"testing"
)

func TestHashString(t *testing.T) {
	t.Log(Parse(`
a = 1e3
b = a + 1
_=[[
s]]
-- local b, ... = a()
`, ""))
}
