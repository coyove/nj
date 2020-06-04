package potatolang

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// type str struct {
// 	p unsafe.Pointer
// 	a int
// }
//
// func mstr(s string) str {
// 	bs := [8]byte{}
// 	copy(bs[:], s)
// 	// ss := *(*[2]uintptr)(unsafe.Pointer(&s))
// 	// return str{
// 	// 	p: unsafe.Pointer(ss[0]),
// 	// 	a: int(ss[1]),
// 	// }
// 	return str{p: unsafe.Pointer(&bs), a: len(s)}
// }
//
// func (s str) String() string {
// 	var ss string
// 	v := (*[2]uintptr)(unsafe.Pointer(&ss))
// 	(*v)[0] = uintptr(s.p)
// 	(*v)[1] = uintptr(s.a)
// 	return ss
// }
//
// func TestStr(t *testing.T) {
// 	// a := "b"
// 	var s []str
// 	for i := 0; i < 3; i++ {
// 		s = append(s, mstr("a"+strconv.Itoa(i)))
// 	}
// 	runtime.GC()
// 	t.Log(s)
// }

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

const benchMapSize = 1e3
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

}

func TestMapNext(t *testing.T) {
	rand.Seed(time.Now().Unix())
	m := &Map{}
	v := rand.Float64()
	m.Put(Num(v), Num(v))

	count := 0
	for k, _ := m.Next(Value{}); !k.IsNil(); k, _ = m.Next(k) {
		if k.Num() != v {
			t.Fatal(m.Len(), k, v)
		}
		count++
	}

	if count != 1 {
		t.Fatal(count)
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
	for v, i := range rand.Perm(benchMapSize * 2) {
		c[i] = v
		m.Put(Num(float64(i)), Num(float64(v)))
	}

	for k := range c {
		m.Put(Num(float64(k)), Value{})
		delete(c, k)
		if len(c) == benchMapSize {
			break
		}
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		count := 0
		for k, v := m.Next(Value{}); !k.IsNil(); k, v = m.Next(k) {
			v2 := c[int(k.Num())]
			if int(v.Num()) != v2 {
				b.Fatal(m)
			}
			count++
		}

		if count != len(c) {
			b.Fatal(count)
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
			if i+v == -99 {
				b.Fatal(i)
			}
		}
	}
}

//
// func BenchmarkTreeMapGetN(b *testing.B) {
// 	b.StopTimer()
//
// 	rand.Seed(time.Now().Unix())
// 	m := &Map{}
// 	keys := rand.Perm(benchMapSize2)
// 	for _, k := range keys {
// 		m.Put(Num(float64(k)), Num(float64(k)))
// 	}
//
// 	b.StartTimer()
//
// 	k := Num(float64(keys[rand.Intn(len(keys))]))
// 	for i := 0; i < b.N; i++ {
// 		m.Get(k)
// 	}
// }
//
// func BenchmarkMapGetN(b *testing.B) {
// 	b.StopTimer()
//
// 	rand.Seed(time.Now().Unix())
// 	m := map[tablekey]Value{}
// 	keys := rand.Perm(benchMapSize2)
// 	for _, k := range keys {
// 		m[tablekey{g: Num(float64(k))}] = Num(float64(k))
// 	}
//
// 	b.StartTimer()
//
// 	k := tablekey{g: Num(float64(keys[rand.Intn(len(keys))]))}
// 	for i := 0; i < b.N; i++ {
// 		_ = m[k]
// 	}
// }
