package bas

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"runtime/debug"
	"testing"
	"time"
)

func isEmpty(o *Object) bool {
	for _, kv := range o.items {
		if kv.key != Nil {
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

func randInt(len, idx int) int {
	buf := make([]byte, 6)
	rand.Read(buf)
	v := rand.Int()
	if Int(v).HashCode()%uint64(len) == uint64(idx) {
		return v
	}
	return randInt(len, idx)
}

func TestObjectForeachDelete(t *testing.T) {
	rand.Seed(time.Now().Unix())
	check := func(o *Object, idx int, k int, dist int32) {
		i := o.items[idx]
		if i.key.Int() != k || i.dist != dist {
			t.Fatal(o.items, string(debug.Stack()))
		}
	}

	old := resizeHash
	resizeHash = func(*Object, int) {}
	o := NewObject(1)
	a := randInt(2, 1)
	b := randInt(2, 1)
	c := randInt(2, 1)
	o.Set(Int(a), Int(a)) // [null, a+0]
	o.Set(Int(b), Int(b)) // [b+1, a+0]
	o.Delete(Int(a))      // [b+1, deleted+0]
	if o.items[0].dist != 1 || !o.items[1].pDeleted {
		t.Fatal(o.items)
	}
	o.Set(Int(a), Int(a)) // [a+1, b+0]
	check(o, 0, a, 1)
	check(o, 1, b, 0)

	o.Delete(Int(a))      // [deleted+1, b+0]
	o.Set(Int(c), Int(c)) // [c+1, b+0]
	check(o, 0, c, 1)
	check(o, 1, b, 0)

	o = NewObject(2)
	a = randInt(4, 1)
	b = randInt(4, 1)
	c = randInt(4, 1)
	d := randInt(4, 1)
	o.Set(Int(a), Int(a))
	o.Set(Int(b), Int(b))
	o.Set(Int(c), Int(c))
	o.Set(Int(d), Int(d)) // [d+3, a+0, b+1, c+2]
	check(o, 0, d, 3)
	check(o, 1, a, 0)
	check(o, 2, b, 1)
	check(o, 3, c, 2)

	o.Delete(Int(b))      // [d+3, a+0, deleted+1, c+2]
	o.Set(Int(b), Int(b)) // [b+3, a+0, c+1, d+2]
	check(o, 0, b, 3)
	check(o, 1, a, 0)
	check(o, 2, c, 1)
	check(o, 3, d, 2)

	o.Delete(Int(a))
	o.Delete(Int(c)) // [b+3, deleted+0, deleted+1, d+2]
	loopCount := 0
	o.Foreach(func(k Value, v *Value) bool { loopCount++; return true })
	if loopCount != 2 {
		t.Fatal(loopCount, o.items)
	}
	for k, _ := o.Next(Nil); k != Nil; k, _ = o.Next(k) {
		loopCount--
	}
	if loopCount != 0 {
		t.Fatal(o.items)
	}
	a = randInt(4, 2)
	o.Set(Int(a), Int(a)) // [a+2, deleted+0, d+1, b+2]
	check(o, 0, a, 2)
	check(o, 2, d, 1)
	check(o, 3, b, 2)

	if o.Get(Int(a)).Int() != a {
		t.Fatal(o.items)
	}

	o = NewObject(4)
	a = randInt(8, 1)
	b = randInt(8, 1)
	c = randInt(8, 1)
	d = randInt(8, 2)
	e := randInt(8, 4)
	o.Set(Int(a), Int(a))
	o.Set(Int(b), Int(b))
	o.Set(Int(c), Int(c))
	o.Set(Int(d), Int(d))
	o.Set(Int(e), Int(e)) // [nil, a+0, b+1, c+2, d+2, e+1, nil, nil]
	check(o, 1, a, 0)
	check(o, 2, b, 1)
	check(o, 3, c, 2)
	check(o, 4, d, 2)
	check(o, 5, e, 1)

	o.Delete(Int(b))
	o.Delete(Int(d)) // [nil, a+0, deleted+1, c+2, deleted+2, e+1, nil, nil]
	if o.Get(Int(b)) != Nil {
		t.Fatal(o.items)
	}
	if o.Get(Int(c)) != Int(c) {
		t.Fatal(o.items)
	}
	o.Set(Int(d), Int(d)) // [nil, a+0, c+1, d+1, e+0, nil, nil, nil]
	check(o, 1, a, 0)
	check(o, 2, c, 1)
	check(o, 3, d, 1)
	check(o, 4, e, 0)

	o.Delete(Int(e))
	o.Set(Int(e), Int(e))
	check(o, 1, a, 0)
	check(o, 2, c, 1)
	check(o, 3, d, 1)
	check(o, 4, e, 0)

	resizeHash = old
}

func BenchmarkRHMap10(b *testing.B)      { benchmarkRHMap(b, 10) }
func BenchmarkGoMap10(b *testing.B)      { benchmarkGoMap(b, 10) }
func BenchmarkRHMap20(b *testing.B)      { benchmarkRHMap(b, 20) }
func BenchmarkGoMap20(b *testing.B)      { benchmarkGoMap(b, 20) }
func BenchmarkRHMap50(b *testing.B)      { benchmarkRHMap(b, 50) }
func BenchmarkGoMap50(b *testing.B)      { benchmarkGoMap(b, 50) }
func BenchmarkRHMap5000(b *testing.B)    { benchmarkRHMap(b, 5000) }
func BenchmarkGoMap5000(b *testing.B)    { benchmarkGoMap(b, 5000) }
func BenchmarkRHMapUnc10(b *testing.B)   { benchmarkRHMapUnconstrainted(b, 10) }
func BenchmarkGoMapUnc10(b *testing.B)   { benchmarkGoMapUnconstrainted(b, 10) }
func BenchmarkRHMapUnc1000(b *testing.B) { benchmarkRHMapUnconstrainted(b, 1000) }
func BenchmarkGoMapUnc1000(b *testing.B) { benchmarkGoMapUnconstrainted(b, 1000) }

func benchmarkRHMap(b *testing.B, n int) {
	rand.Seed(time.Now().Unix())
	m := NewObject(n)
	for i := 0; i < n; i++ {
		m.Set(Int64(int64(i)), Int64(int64(i)))
	}
	for i := 0; i < b.N; i++ {
		idx := rand.Intn(n)
		if m.Find(Int64(int64(idx))) != Int64(int64(idx)) {
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
		if m.Find(Int64(k)).Int64() != v {
			m.Foreach(func(mk Value, mv *Value) bool {
				if mk.Int64() == k {
					t.Log(mk, *mv)
				}
				return true
			})
			t.Fatal(m.Find(Int64(k)), k, v)
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

func TestObjectDistance(t *testing.T) {
	test := func(sz int) {
		o := NewObject(sz)
		for i := 0; i < sz; i++ {
			o.Set(Int(randInt(sz, i)), Int(i))
		}
		for _, i := range o.items {
			if i.key != Nil {
				if i.dist != 0 {
					t.Fatal(o.items)
				}
			}
		}
	}
	for i := 1; i < 16; i++ {
		test(i)
	}
	test = func(sz int) {
		o := NewObject(sz / 2)
		for i := 0; i < sz; i++ {
			o.Set(Int(randInt(sz, i)), Int(i))
		}
		for _, i := range o.items {
			if i.dist != 0 {
				t.Fatal(o.items)
			}
		}
	}
	old := resizeHash
	resizeHash = func(*Object, int) {}
	for i := 2; i <= 16; i += 2 {
		test(i)
	}
	resizeHash = old
}

func TestHashcodeDist(t *testing.T) {
	rand.Seed(time.Now().Unix())
	for _, a := range []string{"a", "b", "c", "z", randString(), randString(), randString(), randString()} {
		fmt.Println(Str(a).HashCode() % 32)
	}

	z := map[uint64]int{}
	rand.Seed(time.Now().Unix())
	for i := 0; i < 1e6; i++ {
		v := Int64(int64(i)).HashCode() % 32
		z[v]++
	}
	fmt.Println(z)

	z = map[uint64]int{}
	for i := 0; i < 1e6; i++ {
		v := Int64(rand.Int63()).HashCode() % 32
		z[v]++
	}
	fmt.Println((z))

	z = map[uint64]int{}
	for i := 0; i < 1e6; i++ {
		v := Str(randString()).HashCode() % 32
		z[v]++
	}
	fmt.Println((z))
}
