package potatolang

import (
	"math/rand"
	"testing"
)

func TestNewStack(t *testing.T) {
	st := NewEnv(nil)

	v := NewNumberValue(19930731)
	vi := 0

	for {
		idx := rand.Intn(1000)
		if rand.Intn(100) == 0 {
			st.LocalSet(idx, v)
			vi = idx
			break
		}

		st.LocalSet(idx, Value{})
	}

	if !st.LocalGet(vi).Equal(v) {
		t.Error(v, vi)
	}
}
