package bas

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/coyove/nj/typ"
)

func isEmpty(o *Object) bool {
	for _, kv := range o.items {
		if kv.Key != Nil {
			return false
		}
	}
	return true
}

func randString() string {
	buf := make([]byte, 6)
	rand.Read(buf)
	return base64.StdEncoding.EncodeToString(buf)
}

func TestObjectForeachDelete(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for i := 1; i <= 16; i++ {
		o := NewObject(i)
		o.SetProp("a", Zero)
		o.Foreach(func(k Value, v *Value) int { return typ.ForeachDeleteContinue })
		if !isEmpty(o) {
			t.Fatal(o.items)
		}
	}

	for i := 1; i <= 16; i++ {
		o := NewObject(i)
		o.SetProp("a", Int(1))
		o.SetProp("b", Int(2))
		o.Foreach(func(k Value, v *Value) int { return typ.ForeachDeleteContinue })
		if !isEmpty(o) {
			t.Fatal(o.items)
		}
	}

	for i := 1; i <= 16; i++ {
		o := NewObject(i)
		m := map[string]Value{}
		o.SetProp("a", Int(1))
		o.SetProp("b", Int(2))
		m["a"] = Int(1)
		m["b"] = Int(2)
		o.Foreach(func(k Value, v *Value) int { delete(m, k.Str()); return typ.ForeachDeleteBreak })
		for k, v := range m {
			if o.Prop(k).Int() != v.Int() {
				t.Fatal(o.items)
			}
		}
	}

	for i := 1; i <= 16; i++ {
		o := NewObject(i)
		m := map[string]int64{}
		for i := 0; i < 1e5; i++ {
			k, v := randString(), rand.Int63()
			o.SetProp(k, Int64(v))
			m[k] = v
		}
		// fmt.Println(o.items)
		// fmt.Println()
		o.Foreach(func(k Value, v *Value) int {
			if rand.Intn(2) == 0 {
				delete(m, k.Str())
				return typ.ForeachDeleteContinue
			}
			return typ.ForeachContinue
		})
		for k, v := range m {
			if ov := o.Prop(k); ov.Int64() != v {
				// t.Fatal(o.Len(), len(m), ov, v)
				t.Fatal(k, Str(k).HashCode()%uint64(len(o.items)), o.items, m)
			}
		}
	}
}

func BenchmarkRHMap10(b *testing.B)    { benchmarkRHMap(b, 10) }
func BenchmarkGoMap10(b *testing.B)    { benchmarkGoMap(b, 10) }
func BenchmarkRHMap20(b *testing.B)    { benchmarkRHMap(b, 20) }
func BenchmarkGoMap20(b *testing.B)    { benchmarkGoMap(b, 20) }
func BenchmarkRHMap50(b *testing.B)    { benchmarkRHMap(b, 50) }
func BenchmarkGoMap50(b *testing.B)    { benchmarkGoMap(b, 50) }
func BenchmarkRHMapUnc10(b *testing.B) { benchmarkRHMapUnconstrainted(b, 10) }
func BenchmarkGoMapUnc10(b *testing.B) { benchmarkGoMapUnconstrainted(b, 10) }

func benchmarkRHMap(b *testing.B, n int) {
	rand.Seed(time.Now().Unix())
	m := NewObject(n)
	for i := 0; i < n; i++ {
		m.Set(Int64(int64(i)), Int64(int64(i)))
	}
	for i := 0; i < b.N; i++ {
		idx := rand.Intn(n)
		if m.Get(Int64(int64(idx))) != Int64(int64(idx)) {
			b.Fatal(idx, m)
		}
	}
}

func benchmarkRHMapUnconstrainted(b *testing.B, n int) {
	rand.Seed(time.Now().Unix())
	m := NewObject(1)
	for i := 0; i < b.N; i++ {
		for i := 0; i < n; i++ {
			x := rand.Intn(n)
			m.Set(Int64(int64(x)), Int64(int64(i)))
		}
	}
}

func benchmarkGoMap(b *testing.B, n int) {
	rand.Seed(time.Now().Unix())
	m := map[int]int{}
	for i := 0; i < n; i++ {
		m[i] = i
	}
	for i := 0; i < b.N; i++ {
		idx := rand.Intn(n)
		if m[idx] == -1 {
			b.Fatal(idx, m)
		}
	}
}

func benchmarkGoMapUnconstrainted(b *testing.B, n int) {
	rand.Seed(time.Now().Unix())
	m := map[int]int{}
	for i := 0; i < b.N; i++ {
		for i := 0; i < n; i++ {
			idx := rand.Intn(n)
			m[idx] = i
		}
	}
}

func TestRHMap(t *testing.T) {
	rand.Seed(time.Now().Unix())
	m := NewObject(0)
	m2 := map[int64]int64{}
	counter := int64(0)
	for i := 0; i < 1e6; i++ {
		x := rand.Int63()
		if x%2 == 0 {
			x = counter
			counter++
		}
		m.Set(Int64(x), Int64(x))
		m2[x] = x
	}
	for k := range m2 {
		delete(m2, k)
		m.Delete(Int64(k))
		if rand.Intn(10000) == 0 {
			break
		}
	}

	fmt.Println(m.Len(), m.Size(), len(m2))

	for k, v := range m2 {
		if m.Get(Int64(k)).Int64() != v {
			m.Foreach(func(mk Value, mv *Value) int {
				if mk.Int64() == k {
					t.Log(mk, *mv)
				}
				return typ.ForeachContinue
			})
			t.Fatal(m.Get(Int64(k)), k, v)
		}
	}

	if m.Len() != len(m2) {
		t.Fatal(m.Len(), len(m2))
	}

	for k, v := m.Next(Nil); k != Nil; k, v = m.Next(k) {
		if _, ok := m2[k.Int64()]; !ok {
			t.Fatal(k, v, len(m2))
		}
		delete(m2, k.Int64())
	}
	if len(m2) != 0 {
		t.Fatal(len(m2))
	}

	m.Clear()
	m.Set(Int64(0), Int64(0))
	m.Set(Int64(1), Int64(1))
	m.Set(Int64(2), Int64(2))

	for i := 4; i < 9; i++ {
		m.Set(Int64(int64(i*i)), Int64(0))
	}

	for k, v := m.Next(Nil); k != Nil; k, v = m.Next(k) {
		fmt.Println(k, v)
	}
}
