package potatolang

import (
	"math/rand"
	"runtime"
	"testing"
)

func TestNewStringValue(t *testing.T) {
	m := make(map[string]Value)
	for i := 0; i < 10000; i++ {
		ln := rand.Intn(10) + 5
		str := ""
		for j := 0; j < ln; j++ {
			str += string(rand.Intn(26) + 'a')
		}
		m[str] = NewStringValue(str)
	}

	runtime.GC()

	for k, v := range m {
		if k != v.AsStringUnsafe() {
			t.Error(k, v)
		}
	}
}
