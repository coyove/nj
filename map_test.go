package potatolang

import (
	"testing"
)

func BenchmarkSlice(b *testing.B) {
	f := func() {
		v := []Value{}
		for i := 0; i < 10; i++ {
			v = append(v, Num(float64(i)))
		}
	}
	for i := 0; i < b.N; i++ {
		f()
	}
}

func BenchmarkTable(b *testing.B) {
	f := func() {
		v := Table{}
		for i := 0; i < 10; i++ {
			v.Put(Num(float64(i)), Num(float64(i)))
		}
	}
	for i := 0; i < b.N; i++ {
		f()
	}
}

func TestTable(t *testing.T) {
	m := Table{}
	m.Put(Str("hello"), Num(1))
	m.Put(Num(0), Num(0))
	m.Put(Num(1), Num(1))
	m.Put(Num(2), Num(2))
	i := m.Iter()
	for i.Next() {
		t.Log(i.Key())
	}
	m.Put(Num(1), Value{})
	m.Put(Num(2), Value{})
	i = m.Iter()
	for i.Next() {
		t.Log(i.Key(), i.Value())
	}
}
