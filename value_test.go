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

	assert(NewNumberValue(0).IsZero())
	assert(NewNumberValue(0).IsFalse())
	assert(!NewNumberValue(1 / math.Inf(-1)).IsFalse())
	assert(!NewNumberValue(1 / math.Inf(-1)).IsZero())
	assert(!NewNumberValue(math.NaN()).IsFalse())

	s := NewStringValue("")
	assert(s.IsFalse())
	s = NewBoolValue(true)
	assert(!s.IsFalse())
	s = NewBoolValue(false)
	assert(s.IsFalse())
}

func BenchmarkSmallStringEquality(b *testing.B) {
	a, a0 := NewStringValue("true"), NewStringValue("true")
	for i := 0; i < b.N; i++ {
		a.Equal(a0)
	}
}

func BenchmarkSmallStringEquality2(b *testing.B) {
	a, a0 := NewBoolValue(true), NewBoolValue(true)
	for i := 0; i < b.N; i++ {
		a.Equal(a0)
	}
}

func BenchmarkIsZero(b *testing.B) {
	a := NewBoolValue(false)
	for i := 0; i < b.N; i++ {
		a.IsZero()
	}
}
