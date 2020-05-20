package potatolang

import (
	"testing"
)

func BenchmarkSlice(b *testing.B) {
	f := func() {
		v := []Value{}
		for i := 0; i < 10; i++ {
			v = append(v, NewNumberValue(float64(i)))
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
			v.Put(NewNumberValue(float64(i)), NewNumberValue(float64(i)))
		}
	}
	for i := 0; i < b.N; i++ {
		f()
	}
}

func TestTable(t *testing.T) {
	m := Table{}
	m.Put(NewStringValue("hello"), NewNumberValue(1))
	m.Put(NewNumberValue(0), NewNumberValue(0))
	m.Put(NewNumberValue(1), NewNumberValue(1))
	m.Put(NewNumberValue(2), NewNumberValue(2))
	i := m.Iter()
	for i.Next() {
		t.Log(i.Key())
	}
	m.Put(NewNumberValue(1), Value{})
	m.Put(NewNumberValue(2), Value{})
	i = m.Iter()
	for i.Next() {
		t.Log(i.Key(), i.Value())
	}
}
