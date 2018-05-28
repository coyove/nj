package potatolang

import (
	"math/rand"
	"runtime"
	"strconv"
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
		if k != v.AsString() {
			t.Error(k, v)
		}
	}
}

const key = "noenvescape"

func BenchmarkMapValue(b *testing.B) {
	m := make(map[hash128]Value)

	for i := 0; i < 1000; i++ {
		x := NewNumberValue(float64(i))
		m[x.Hash()] = x
	}
	m[NewStringValue(key).Hash()] = NewNumberValue(1)
	// x := b.N % 1000
	for i := 0; i < b.N; i++ {
		if m[NewStringValue(key).Hash()].AsNumber() != 1 {

		}
	}

	// b.Error(hash.Hash())
}

func BenchmarkMapValue2(b *testing.B) {
	m := make(map[string]Value)

	for i := 0; i < 1000; i++ {
		x := NewNumberValue(float64(i))
		m[strconv.Itoa(i)] = x
	}
	m[key] = NewNumberValue(1)
	// x := b.N % 1000
	for i := 0; i < b.N; i++ {
		if m[key].AsNumber() != 1 {
		}
	}
}
