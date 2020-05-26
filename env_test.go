package potatolang

import (
	"math/rand"
	"testing"
)

func TestNewStack(t *testing.T) {
	st := NewEnv(nil)

	v := Num(19930731)
	vi := 0

	for {
		idx := rand.Intn(1000)
		if rand.Intn(100) == 0 {
			st.Set(idx, v)
			vi = idx
			break
		}

		st.Set(idx, Value{})
	}

	if !st.Get(vi).Equal(v) {
		t.Error(v, vi)
	}
}
