package rpn

import (
	"testing"
)

func TestBasic(t *testing.T) {
	r := NewReaderFromString(`[
	fun 'fib n [
		if n 2 < [n return]
		n ++ -- 1 <<
		n 1 - fib!
		n 2 - fib! +
	]]
	`)

	if true {
		t.Log(parse(r, newfuncargs(nil)))
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
