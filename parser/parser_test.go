package parser

import (
	"os"
	"strings"
	"testing"
)

func TestTokenName(t *testing.T) {
	c, err := Parse(strings.NewReader(`if (0 == 1 ) {
		assert 0;
	} else if (2 == 2) {
		assert 1;
	} else {
		assert 0;
	}`), "mem")
	t.Error(err)
	c.Dump(os.Stderr)
}
