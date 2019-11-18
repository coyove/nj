package hash50

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"unsafe"
)

var hashMap = struct {
	sync.RWMutex
	rev map[uint64][]byte
}{
	rev: map[uint64][]byte{},
}

var lookup = []byte{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 0, 0, 0, 0, 0, 0,
	0, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25,
	26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 0, 0, 0, 0, 37,
	0, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52,
	53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 0, 0, 0, 0, 0,
}

var revLookup = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz"

func genTable() {
	m := [64]byte{}

	c := 1
	for i := 0; i < 128; i++ {
		if strings.Contains(revLookup, string(rune(i))) {
			fmt.Println(c, ",")
			m[c] = byte(i)
			c++
		} else {
			fmt.Println(0, ",")
		}
	}
	fmt.Println("====")
	for _, v := range m {
		fmt.Println(v, ",")
	}
}

//go:nosplit
func HashString(str string) float64 {
	s := (*reflect.StringHeader)(unsafe.Pointer(&str))
	p := &reflect.SliceHeader{
		Data: s.Data,
		Len:  s.Len,
		Cap:  s.Len,
	}
	return Hash(*(*[]byte)(unsafe.Pointer(p)))
}

func Hash(strbuf []byte) float64 {
	var hash uint64 = 2166136261

	if len(strbuf) < 8 {
		hash = 0
		for i, v := range strbuf {
			hash |= (uint64(lookup[v]) & 0x3f) << ((7 - i) * 6)
		}
		return float64(hash)
	}

	for _, c := range strbuf {
		hash *= 16777619
		hash ^= uint64(c)
	}
	hash &= 0x3ffffffffffff
	hash |= 1 << 49

	hashMap.Lock()
	if t, ok := hashMap.rev[hash]; ok && !bytes.Equal(t, strbuf) {
		panic(fmt.Sprint(string(t), "and", string(strbuf), "share an identical hash:", hash))
	}
	hashMap.rev[hash] = strbuf
	hashMap.Unlock()

	return float64(hash)
}

func FindStringHash(h float64) []byte {
	x := uint64(h)
	if x>>49 == 0 {
		i, buf := 0, [8]byte{}
		for ; i < 8; i++ {
			v := (x >> ((7 - i) * 6)) & 0x3f
			if v == 0 {
				break
			}
			buf[i] = revLookup[v-1]
		}
		return buf[:i]
	}

	hashMap.RLock()
	v := hashMap.rev[x]
	hashMap.RUnlock()
	return v
}
