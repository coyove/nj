package bench

import (
	"bytes"
	"math/rand"
	"sort"
	"strconv"
	"testing"
	"time"
	"unsafe"
)

func concatBytes(a, b []byte) []byte {
	c := make([]byte, len(a)+len(b))
	copy(c, a)
	copy(c[len(a):], b)
	return c
}

func bytesString(a []byte) string {
	return *(*string)(unsafe.Pointer(&a))
}

func bytesString2(a []byte) string {
	return string(a)
}

func TestBytesInterface(t *testing.T) {
	a := []byte("zzz")
	m := map[interface{}]int{}
	m[bytesString(a)] = 99
	a[1] = 'a'
	t.Log(m)
}

var dummies = []string{strconv.Itoa(rand.Int()), strconv.Itoa(rand.Int()), strconv.Itoa(rand.Int()), strconv.Itoa(rand.Int())}

func BenchmarkBytes(b *testing.B) {
	v1, v2 := make([]byte, 18), make([]byte, 18)
	rand.Read(v1)
	rand.Read(v2)

	for i := 0; i < b.N; i++ {
		v3 := concatBytes(v1, v2)
		if bytes.Equal(v3, nil) {
			b.Fatal(1)
		}
	}
}

func BenchmarkString(b *testing.B) {
	v1, v2 := dummies[rand.Intn(len(dummies))], dummies[rand.Intn(len(dummies))]
	for i := 0; i < b.N; i++ {
		v3 := v1 + v2
		if v3 == "" {
			b.Fatal(1)
		}
	}
}

func BenchmarkBytesReplace(b *testing.B) {
	v1 := dummies[rand.Intn(len(dummies))]

	for i := 0; i < b.N; i++ {
		x := []byte(v1)
		x[1] = 'z'
		v3 := string(x)
		if v3 == "" {
			b.Fatal(1)
		}
	}
}

func BenchmarkStringReplace(b *testing.B) {
	v1 := dummies[rand.Intn(len(dummies))]
	for i := 0; i < b.N; i++ {
		v3 := v1[:1] + "z" + v1[2:]
		if v3 == "" {
			b.Fatal(1)
		}
	}
}

type IntMap []uint64

func (m IntMap) Get(k uint64) (v uint64, ok bool) {
	offset := len(m) / 2
	j := offset

	if j > 0 && m[0] == k {
		return m[offset], true
	}

	if j > 1 && m[1] == k {
		return m[1+offset], true
	}
	//	return 0, false

	for i := 2; i < j; {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h

		k2 := &m[h]
		if *k2 == k {
			return *(*uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(k2)) + uintptr(offset)*unsafe.Sizeof(uint64(0)))), true
			// return m.values[h+offset], true
		}

		if *k2 < k {
			i = h + 1
		} else {
			j = h
		}
	}

	return 0, false
}

func (m *IntMap) Add(k uint64, v uint64) {
	offset := len(*m) / 2
	i, j := 0, offset

	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		// i ≤ h < j
		if (*m)[h] == k {
			(*m)[h+offset] = v
			return
		}

		if (*m)[h] < k {
			i = h + 1 // preserves f(i-1) == false
		} else {
			j = h // preserves f(j) == true
		}
	}

	*m = append(*m, 0, 0)
	copy((*m)[i+offset+2:], (*m)[i+offset:])
	(*m)[i+offset+1] = v
	copy((*m)[i+1:], (*m)[i:i+offset])
	(*m)[i] = k
}

type kvSwapper []uint64

func (s kvSwapper) Len() int {
	return len(s) / 2
}

func (s kvSwapper) Less(i, j int) bool {
	return (s)[i] < (s)[j]
}

func (s kvSwapper) Swap(i, j int) {
	(s)[i], (s)[j] = (s)[j], (s)[i]
	i, j = i+len(s)/2, j+len(s)/2
	(s)[i], (s)[j] = (s)[j], (s)[i]
}

func (m *IntMap) BatchSet(kv []uint64) {
	*m = kv
	sort.Sort(kvSwapper(*m))
}

var N = 16

func TestMapBatch(t *testing.T) {
	rand.Seed(time.Now().Unix())

	m := IntMap{}
	m2 := map[uint64]bool{}

	args := make([]uint64, N*2)
	for i := 0; i < N; i++ {
		x := uint64(rand.Int())
		m2[x] = true
		args[i] = x
		args[i+N] = x
	}

	m.BatchSet(args)

	for k := range m2 {
		if v, _ := m.Get(k); v != k {
			t.Fatal(m)
		}
	}
}

func TestMapAdd(t *testing.T) {
	m := IntMap{}
	m2 := map[uint64]bool{}
	for i := 0; i < N; i++ {
		x := uint64(rand.Int())
		m.Add(x, x)
		m2[x] = true
	}

	for k := range m2 {
		if v, _ := m.Get(k); v != k {
			t.Fatal(m)
		}
	}
}

func BenchmarkNativeMapIndex(b *testing.B) {
	b.StopTimer()
	m := map[uint64]uint64{}
	for i := 0; i < N; i++ {
		x := uint64(rand.Int())
		m[x] = x
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		if m[uint64(i)] == 999 {
			b.Fatal(m)
		}
	}
}

func BenchmarkMapIndex(b *testing.B) {
	b.StopTimer()
	m := IntMap{}
	for i := 0; i < N; i++ {
		x := uint64(rand.Int())
		m.Add(x, x)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		if v, _ := m.Get(uint64(i)); v == 999 {
			b.Fatal(m)
		}
	}

	//b.Log(m)
}

func BenchmarkNativeMapAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := map[uint64]uint64{}
		for i := 0; i < N; i++ {
			x := uint64(rand.Int())
			m[x] = x
		}
	}
}

func BenchmarkMapAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := IntMap{}
		for i := 0; i < N; i++ {
			x := uint64(rand.Int())
			m.Add(x, x)
		}
	}
}

func BenchmarkMapBatch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := IntMap{}
		kv := make([]uint64, N*2)
		for i := 0; i < N; i++ {
			x := uint64(rand.Int())
			kv[i] = x
			kv[i+N] = x
		}
		m.BatchSet(kv)
	}
}

//func BenchmarkArrayIndex(b *testing.B) {
//	m := make([]uint32, 4)
//	for i := 0; i < b.N; i++ {
//		if m[uint32(i%4)] == 1 {
//			b.Fatal(m)
//		}
//	}
//}
