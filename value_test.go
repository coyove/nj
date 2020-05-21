package potatolang

import (
	"math"
	"runtime"
	"testing"
)

func TestFalsyValue(t *testing.T) {
	assert := func(b bool) {
		if !b {
			_, fn, ln, _ := runtime.Caller(1)
			t.Fatal(fn, ln)
		}
	}

	assert(Num(0).IsZero())
	assert(Num(0).IsFalse())
	assert(!Num(1 / math.Inf(-1)).IsFalse())
	assert(!Num(1 / math.Inf(-1)).IsZero())
	assert(!Num(math.NaN()).IsFalse())

	s := Str("")
	assert(s.IsFalse())
	s = Bln(true)
	assert(!s.IsFalse())
	s = Bln(false)
	assert(s.IsFalse())
}

func BenchmarkSmallStringEquality(b *testing.B) {
	a, a0 := Str("true"), Str("true")
	for i := 0; i < b.N; i++ {
		a.Equal(a0)
	}
}

func BenchmarkSmallStringEquality2(b *testing.B) {
	a, a0 := Bln(true), Bln(true)
	for i := 0; i < b.N; i++ {
		a.Equal(a0)
	}
}

func BenchmarkIsZero(b *testing.B) {
	a := Bln(false)
	for i := 0; i < b.N; i++ {
		a.IsZero()
	}
}
