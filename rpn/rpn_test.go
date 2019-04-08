package rpn

import (
	"testing"
)

func TestBasic(t *testing.T) {
	r := NewReaderFromString(`[
	fun 'fib n [
		if n 2 < [n return]
		n ++ --
		n 1 - fib!
		n 2 - fib! +
	]]
	`)

	if true {
		t.Log(parse(r))
		return
	}

	for {
		tok, err := r.Token()
		if err != nil {
			t.Log(err)
			break
		}
		t.Log(tok)
	}
}
