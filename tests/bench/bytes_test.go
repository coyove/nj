package bench

import (
	"bytes"
	"math/rand"
	"strconv"
	"testing"
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

type IntMap struct {
	values []uintptr
}

func (m *IntMap) Get(k uintptr) (v uintptr) {
	offset := len(m.values) / 2
	i, j := 0, offset

	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h

		k2 := m.values[h]
		if k2 == k {
			return m.values[h+offset]
		}
		// i ≤ h < j
		if k2 < k {
			i = h + 1 // preserves f(i-1) == false
		} else {
			j = h // preserves f(j) == true
		}
	}

	return 0
}

func (m *IntMap) Set(k uintptr, v uintptr) {
	offset := len(m.values) / 2
	i, j := 0, offset

	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		// i ≤ h < j
		if m.values[h] == k {
			m.values[h+offset] = v
			return
		}

		if m.values[h] < k {
			i = h + 1 // preserves f(i-1) == false
		} else {
			j = h // preserves f(j) == true
		}
	}

	m.values = append(m.values, 0, 0)
	copy(m.values[i+offset+2:], m.values[i+offset:])
	m.values[i+offset+1] = v
	copy(m.values[i+1:], m.values[i:i+offset])
	m.values[i] = k
}

var N = 100

func BenchmarkNativeMapIndex(b *testing.B) {
	b.StopTimer()
	m := map[uintptr]uintptr{}
	for i := 0; i < N; i++ {
		x := uintptr(rand.Int())
		m[x] = x
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		if m[uintptr(i)] == 999 {
			b.Fatal(m)
		}
	}
}

func TestMapIndex(t *testing.T) {
	m := IntMap{}
	m2 := map[uintptr]bool{}
	for i := 0; i < N; i++ {
		x := uintptr(rand.Int())
		m.Set(x, x)
		m2[x] = true
	}

	for k := range m2 {
		if m.Get(k) != k {
			t.Fatal(m)
		}
	}

	//b.Log(m)
}

func BenchmarkMapIndex(b *testing.B) {
	b.StopTimer()
	m := IntMap{}
	for i := 0; i < N; i++ {
		x := uintptr(rand.Int())
		m.Set(x, x)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		if m.Get(uintptr(i)) == 999 {
			b.Fatal(m)
		}
	}

	//b.Log(m)
}

//func BenchmarkArrayIndex(b *testing.B) {
//	m := make([]uint32, 4)
//	for i := 0; i < b.N; i++ {
//		if m[uint32(i%4)] == 1 {
//			b.Fatal(m)
//		}
//	}
//}
