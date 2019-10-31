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
