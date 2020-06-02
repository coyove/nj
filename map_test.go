package potatolang

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
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
			v.Put(Num(float64(i)), Num(float64(i)), false)
		}
	}
	for i := 0; i < b.N; i++ {
		f()
	}
}

const benchMapSize = 1e5
const benchMapSize2 = 100

func TestTreeMap(t *testing.T) {
	rand.Seed(time.Now().Unix())
	m := &Map{}
	c := map[int]int{}
	for v, i := range rand.Perm(benchMapSize) {
		c[i] = v
		m.Put(Num(float64(i)), Num(float64(v)))
	}

	// fmt.Println(m.GoString())

	for i, v := range c {
		v2 := int(m.Get(Num(float64(i))).Num())
		if v != v2 {
			t.Fatal(m)
		}
	}

	mid := benchMapSize / 2
	t.Log(mid)
	for k, _ := Next(m, Num(float64(mid))); !k.IsNil(); k, _ = Next(m, k) {
		if k.Num() != float64(mid)-1 {
			t.Fatal(m.GoString())
		}
		mid--
		// 	fmt.Println(m.GoString())
	}
	if mid != 0 {
		t.Fatal(mid)
	}

	// fmt.Println(m.GoString())
}

func TestTreeMapDelete(t *testing.T) {
	rand.Seed(time.Now().Unix())
	m := &Map{}
	c := map[int]int{}
	for v, i := range rand.Perm(benchMapSize) {
		c[i] = v
		m.Put(Num(float64(i)), Num(float64(v)))
		if rand.Intn(2) == 1 {
			delete(c, i)
			m.Put(Num(float64(i)), Value{})
		}
	}

	fmt.Println(m.Len(), len(c))

	for i, v := range c {
		v2 := int(m.Get(Num(float64(i))).Num())
		if v != v2 {
			t.Fatal(m)
		}
	}

	// fmt.Println(m.GoString())
}

func BenchmarkTreeMap(b *testing.B) {
	b.StopTimer()

	rand.Seed(time.Now().Unix())
	m := &Map{}
	c := map[int]int{}
	for v, i := range rand.Perm(benchMapSize) {
		c[i] = v
		m.Put(Num(float64(i)), Num(float64(v)))
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for i, v := range c {
			v2 := int(m.Get(Num(float64(i))).Num())
			if v != v2 {
				b.Fatal(m)
			}
		}
	}
}

func BenchmarkTreeMapRev(b *testing.B) {
	b.StopTimer()

	rand.Seed(time.Now().Unix())
	m := &Map{}
	c := map[int]int{}
	for v, i := range rand.Perm(benchMapSize) {
		c[i] = v
		m.Put(Num(float64(i)), Num(float64(v)))
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for k, v := Next(m, Value{}); !k.IsNil(); k, v = Next(m, k) {
			v2 := c[int(k.Num())]
			if int(v.Num()) != v2 {
				b.Fatal(m)
			}
		}
	}
}

func BenchmarkMap(b *testing.B) {
	b.StopTimer()

	rand.Seed(time.Now().Unix())
	c := map[int]int{}
	for v, i := range rand.Perm(benchMapSize) {
		c[i] = v
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for i, v := range c {
			if i+v == 0 {
				b.Fatal(i)
			}
		}
	}
}

func BenchmarkTreeMapGetN(b *testing.B) {
	b.StopTimer()

	rand.Seed(time.Now().Unix())
	m := &Map{}
	keys := rand.Perm(benchMapSize2)
	for _, k := range keys {
		m.Put(Num(float64(k)), Num(float64(k)))
	}

	b.StartTimer()

	k := Num(float64(keys[rand.Intn(len(keys))]))
	for i := 0; i < b.N; i++ {
		m.Get(k)
	}
}

func BenchmarkMapGetN(b *testing.B) {
	b.StopTimer()

	rand.Seed(time.Now().Unix())
	m := map[tablekey]Value{}
	keys := rand.Perm(benchMapSize2)
	for _, k := range keys {
		m[tablekey{g: Num(float64(k))}] = Num(float64(k))
	}

	b.StartTimer()

	k := tablekey{g: Num(float64(keys[rand.Intn(len(keys))]))}
	for i := 0; i < b.N; i++ {
		_ = m[k]
	}
}
