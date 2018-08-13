//+build amd64 arm64 ppc64le ppc64 mips64 mips64le s390x

package potatolang

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/coyove/common/rand"
)

var safePointerAddr = unsafe.Pointer(uintptr(0x400000000000ffff))

const (
	minPhysPageSize = 4096 // (0x1000) in mheap.go
	trueValue       = minPhysPageSize + 1
	falseValue      = minPhysPageSize + 2
)

// Value is the basic value used by VM
type Value struct {
	ty byte

	a byte

	// Number (float64) will be stored at &i
	// It can't be stored in ptr because pointer value smaller than minPhysPageSize will violate the heap.
	// For ty == Tstring and len(str) <= 10, i == len(str) + 1
	i byte

	// For ty == Tstring and i > 0, 10 bytes starting from p[0] will be used to store small strings
	p [4]byte

	ptr unsafe.Pointer
}

func (v *Value) SetNumberValue(f float64) {
	*(*[2]uint64)(unsafe.Pointer(v)) = _zeroRaw
	*(*float64)(unsafe.Pointer(&v.ptr)) = f
}

func (v *Value) SetBoolValue(b bool) {
	*(*[2]uint64)(unsafe.Pointer(v)) = _zeroRaw
	*(*float64)(unsafe.Pointer(&v.ptr)) = float64(*(*byte)(unsafe.Pointer(&b)))
}

// NewStringValue returns a string value
func NewStringValue(s string) Value {
	v := Value{ty: Tstring}

	if len(s) < 11 {
		v.ptr = safePointerAddr
		copy((*(*[10]byte)(unsafe.Pointer(&v.p[0])))[:], s)
		v.i = byte(len(s) + 1)
	} else {
		v.ptr = unsafe.Pointer(&s)
	}
	return v
}

var (
	_zeroRaw = *(*[2]uint64)(unsafe.Pointer(&Zero))
)

// IsZero is a fast way to check whether number value equals to +0
func (v Value) IsZero() bool { return *(*[2]uint64)(unsafe.Pointer(&v)) == _zeroRaw }

// AsNumber cast value to float64
func (v Value) AsNumber() float64 { return (*(*[2]float64)(unsafe.Pointer(&v)))[1] }

// AsString cast value to string
func (v Value) AsString() string {
	if v.i > 0 {
		hdr := reflect.StringHeader{}
		hdr.Len = int(v.i - 1)
		hdr.Data = uintptr(unsafe.Pointer(&v.p[0]))
		return *(*string)(unsafe.Pointer(&hdr))
	}
	return *(*string)(v.ptr)
}

// IsFalse tests whether value contains a "false" value
func (v Value) IsFalse() bool {
	if v.ty == Tnumber {
		return *(*[2]uint64)(unsafe.Pointer(&v)) == _zeroRaw
	}
	if v.ty == Tnil {
		return true
	}
	if v.ty == Tstring {
		return v.i == 1
	}
	if v.ty == Tmap {
		m := (*Map)(v.ptr)
		return len(m.l)+len(m.m) == 0
	}
	return false
}

const (
	// Constants for multiplication: four random odd 64-bit numbers.
	m1    = 16877499708836156737
	m2    = 2820277070424839065
	m3    = 9497967016996688599
	m4    = 15839092249703872147
	m5    = 0x9a81d0a6d5a123ed
	iseed = 0x930731
)

var hashkey [4]uintptr

func init() {
	buf := rand.New().Fetch(32)
	for i := 0; i < 4; i++ {
		hashkey[i] = uintptr(binary.LittleEndian.Uint64(buf[i*8:]))
		hashkey[i] |= 1
	}
	initCoreLibs()
}

// The following code is taken from src/runtime/hash64.go

//go:nosplit
func add(p unsafe.Pointer, x uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

// Note: in order to get the compiler to issue rotl instructions, we
// need to constant fold the shift amount by hand.
// TODO: convince the compiler to issue rotl instructions after inlining.
func rotl_31(x uint64) uint64 {
	return (x << 31) | (x >> (64 - 31))
}

func readUnaligned32(p unsafe.Pointer) uint32 {
	return *(*uint32)(p)
}

func readUnaligned64(p unsafe.Pointer) uint64 {
	return *(*uint64)(p)
}

type hashv struct {
	a, b uint64
}

func (v Value) Hash() hashv {
	var a hashv
	switch v.ty {
	case Tnumber, Tnil, Tclosure, Tmap, Tgeneric:
		a = *(*hashv)(unsafe.Pointer(&v))
	case Tstring:
		if v.i > 0 {
			a = *(*hashv)(unsafe.Pointer(&v))
		} else {
			hdr := (*reflect.StringHeader)(v.ptr)
			seed := uintptr(v.ty) ^ iseed
			s := uintptr(hdr.Len)
			p := unsafe.Pointer(hdr.Data)
			h := uint64(seed + s*hashkey[0])
			h0 := uint64(seed<<1 + s*hashkey[0])

		tail:
			switch {
			case s == 0:
			case s < 4:
				h ^= uint64(*(*byte)(p))
				h ^= uint64(*(*byte)(add(p, s>>1))) << 8
				h ^= uint64(*(*byte)(add(p, s-1))) << 16
				h = rotl_31(h*m1) * m2
			case s <= 8:
				h ^= uint64(readUnaligned32(p))
				h ^= uint64(readUnaligned32(add(p, s-4))) << 32
				h = rotl_31(h*m1) * m2
			case s <= 16:
				h ^= readUnaligned64(p)
				h = rotl_31(h*m1) * m2
				h ^= readUnaligned64(add(p, s-8))
				h = rotl_31(h*m1) * m2
				h0 ^= h
			case s <= 32:
				h ^= readUnaligned64(p)
				h = rotl_31(h*m1) * m2
				h ^= readUnaligned64(add(p, 8))
				h = rotl_31(h*m1) * m2
				h ^= readUnaligned64(add(p, s-16))
				h = rotl_31(h*m1) * m2
				h ^= readUnaligned64(add(p, s-8))
				h = rotl_31(h*m1) * m2
				h0 ^= h
			default:
				v1 := h
				v2 := uint64(seed * hashkey[1])
				v3 := uint64(seed * hashkey[2])
				v4 := uint64(seed * hashkey[3])
				for s >= 32 {
					v1 ^= readUnaligned64(p)
					v1 = rotl_31(v1*m1) * m2
					p = add(p, 8)
					v2 ^= readUnaligned64(p)
					v2 = rotl_31(v2*m2) * m3
					p = add(p, 8)
					v3 ^= readUnaligned64(p)
					v3 = rotl_31(v3*m3) * m4
					p = add(p, 8)
					v4 ^= readUnaligned64(p)
					v4 = rotl_31(v4*m4) * m1
					p = add(p, 8)
					s -= 32
				}
				h = v1 ^ v2 ^ v3 ^ v4
				h0 ^= h
				goto tail
			}

			h ^= h >> 29
			h *= m3
			h ^= h >> 32

			h0 ^= h0 >> 29
			h0 *= m5
			h0 ^= h0 >> 32

			a = hashv{h, h0}
		}
	}
	return a
}

func (v Value) hashstr() string {
	h := v.Hash()
	return fmt.Sprintf("%x", *(*[16]byte)(unsafe.Pointer(&h)))
}

var __hash2Salt = rand.New().Fetch(16)

func (v Value) hash2() [2]uint64 {
	h := v.Hash()
	b := *(*[16]byte)(unsafe.Pointer(&h))
	s := sha1.Sum(append(b[:], __hash2Salt...))
	return *(*[2]uint64)(unsafe.Pointer(&s)) // 20 > 16
}
