package potatolang

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/coyove/common/rand"
)

// the order can't be changed, for any new type, please also add it in parser.go.y typeof
const (
	// Tnil represents nil type
	Tnil = iota
	// Tnumber represents number type
	Tnumber
	// Tstring represents string type
	Tstring
	// Tbool represents bool type
	Tbool
	// Tlist represents list type
	Tlist
	// Tbytes represents bytes list type
	Tbytes
	// Tmap represents map type
	Tmap
	// Tclosure represents closure type
	Tclosure
	// Tgeneric represents generic type
	Tgeneric
)

const (
	_Tnilnil         = Tnil<<8 | Tnil
	_Tnumbernumber   = Tnumber<<8 | Tnumber
	_Tstringstring   = Tstring<<8 | Tstring
	_Tboolbool       = Tbool<<8 | Tbool
	_Tlistlist       = Tlist<<8 | Tlist
	_Tbytesbytes     = Tbytes<<8 | Tbytes
	_Tmapmap         = Tmap<<8 | Tmap
	_Tclosureclosure = Tclosure<<8 | Tclosure
	_Tgenericgeneric = Tgeneric<<8 | Tgeneric
	_Tbytesnumber    = Tbytes<<8 | Tnumber
	_Tlistnumber     = Tlist<<8 | Tnumber
	_Tstringnumber   = Tstring<<8 | Tnumber
	_Tmapnumber      = Tmap<<8 | Tnumber
)

// TMapping maps type to its string representation
var TMapping = map[byte]string{
	Tnil: "nil", Tnumber: "number", Tstring: "string", Tbool: "bool",
	Tclosure: "closure", Tgeneric: "generic", Tlist: "list", Tmap: "map", Tbytes: "bytes",
}

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

// NewValue returns a nil value
func NewValue() Value {
	return Value{ty: Tnil}
}

// NewNumberValue returns a number value
func NewNumberValue(f float64) Value {
	v := Value{ty: Tnumber}
	v.ptr = safePointerAddr
	*(*float64)(unsafe.Pointer(&v.i)) = f
	return v
}

// NewStringValue returns a string value
func NewStringValue(s string) Value {
	v := NewValue()
	v.ty = Tstring

	if len(s) < 11 {
		v.ptr = safePointerAddr
		copy((*(*[10]byte)(unsafe.Pointer(&v.p[0])))[:], s)
		v.i = byte(len(s) + 1)
	} else {
		v.ptr = unsafe.Pointer(&s)
	}
	return v
}

// NewBoolValue returns a boolean value
func NewBoolValue(b bool) Value {
	v := Value{ty: Tbool}
	if b {
		v.ptr = unsafe.Pointer(uintptr(trueValue))
	} else {
		v.ptr = unsafe.Pointer(uintptr(falseValue))
	}
	return v
}

// NewListValue returns a list value
func NewListValue(a []Value) Value {
	return Value{ty: Tlist, ptr: unsafe.Pointer(&a)}
}

// NewMapValue returns a map value
func NewMapValue(m *Map) Value {
	return Value{ty: Tmap, ptr: unsafe.Pointer(m)}
}

// NewClosureValue returns a closure value
func NewClosureValue(c *Closure) Value {
	return Value{ty: Tclosure, ptr: unsafe.Pointer(c)}
}

// NewBytesValue returns a bytes value
func NewBytesValue(buf []byte) Value {
	return Value{ty: Tbytes, ptr: unsafe.Pointer(&buf)}
}

// NewGenericValue returns a generic value
func NewGenericValue(g interface{}) Value {
	return Value{ty: Tgeneric, ptr: unsafe.Pointer(&g)}
}

// Type returns the type of value
func (v Value) Type() byte {
	return v.ty
}

// AsBool cast value to bool
func (v Value) AsBool() bool { return uintptr(v.ptr) == trueValue }

func (v Value) u64() uint64 {
	if v.ty != Tnumber {
		log.Panicf("expecting number, got %+v", v)
	}
	return math.Float64bits(v.AsNumber())
}

// AsNumber cast value to float64
func (v Value) AsNumber() float64 {
	return *(*float64)(unsafe.Pointer(&v.i))
}

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

// AsMap cast value to map of values
func (v Value) AsMap() *Map { return (*Map)(v.ptr) }

// AsClosure cast value to closure
func (v Value) AsClosure() *Closure { return (*Closure)(v.ptr) }

// AsGeneric cast value to interface{}
func (v Value) AsGeneric() interface{} { return *(*interface{})(v.ptr) }

// AsList cast value to slice of values
func (v Value) AsList() []Value {
	if v.ptr == nil {
		return nil
	}
	return *(*[]Value)(v.ptr)
}

// AsBytes cast value to []byte
func (v Value) AsBytes() []byte {
	if v.ptr == nil {
		return nil
	}
	return *(*[]byte)(v.ptr)
}

// I returns the golang interface representation of value
// it is not the same as AsGeneric()
func (v Value) I() interface{} {
	switch v.Type() {
	case Tbool:
		return v.AsBool()
	case Tnumber:
		return v.AsNumber()
	case Tstring:
		return v.AsString()
	case Tlist:
		return v.AsList()
	case Tmap:
		return v.AsMap()
	case Tbytes:
		return v.AsBytes()
	case Tclosure:
		return v.AsClosure()
	case Tgeneric:
		return v.AsGeneric()
	}
	return nil
}

func (v Value) String() string {
	switch v.Type() {
	case Tstring:
		return strconv.Quote(v.AsString())
	default:
		return v.ToPrintString()
	}
}

// IsFalse tests whether value contains a "false" value
func (v Value) IsFalse() bool {
	switch v.Type() {
	case Tnil:
		return true
	case Tbool:
		return v.AsBool() == false
	case Tnumber:
		return v.AsNumber() == 0.0
	case Tstring:
		return v.AsString() == ""
	case Tlist:
		return len(v.AsList()) == 0
	case Tbytes:
		return len(v.AsBytes()) == 0
	case Tmap:
		return v.AsMap().Size() == 0
	}
	return false
}

// Equal tests whether value is equal to another value
func (v Value) Equal(r Value) bool {
	if v.ty == Tnil || r.ty == Tnil {
		return v.ty == r.ty
	}

	switch v.ty {
	case Tnumber:
		if r.ty == Tnumber {
			return r.AsNumber() == v.AsNumber()
		}
	case Tstring:
		if r.ty == Tstring {
			return r.AsString() == v.AsString()
		} else if r.ty == Tbytes {
			return bytes.Equal(r.AsBytes(), []byte(v.AsString()))
		}
	case Tbool:
		if r.ty == Tbool {
			return r.AsBool() == v.AsBool()
		}
	case Tlist:
		if r.ty == Tlist {
			lf, rf := v.AsList(), r.AsList()

			if len(lf) != len(rf) {
				return false
			}

			for i := 0; i < len(lf); i++ {
				if !lf[i].Equal(rf[i]) {
					return false
				}
			}

			return true
		}
	case Tbytes:
		if r.ty == Tbytes {
			return bytes.Equal(v.AsBytes(), r.AsBytes())
		} else if r.ty == Tstring {
			return bytes.Equal(v.AsBytes(), []byte(r.AsString()))
		}
	case Tmap:
		if r.ty == Tmap {
			return v.AsMap().Equal(r.AsMap())
		}
	}

	return false
}

// ToPrintString returns the printable string of value
func (v Value) ToPrintString() string {
	return v.toString(0)
}

func (v Value) toString(lv int) string {
	if lv > 32 {
		return "<omit deep nesting>"
	}

	switch v.Type() {
	case Tbool:
		return strconv.FormatBool(v.AsBool())
	case Tnumber:
		return strconv.FormatFloat(v.AsNumber(), 'f', -1, 64)
	case Tstring:
		return v.AsString()
	case Tlist:
		arr := v.AsList()
		buf := &bytes.Buffer{}
		buf.WriteString("[")
		for _, v := range arr {
			buf.WriteString(v.toString(lv + 1))
			buf.WriteString(",")
		}
		if len(arr) > 0 {
			buf.Truncate(buf.Len() - 1)
		}
		buf.WriteString("]")
		return buf.String()
	case Tmap:
		m, buf := v.AsMap(), &bytes.Buffer{}
		buf.WriteString("{")
		for i, v := range m.l {
			buf.WriteString(strconv.Itoa(i))
			buf.WriteString("=")
			buf.WriteString(v.toString(lv + 1))
			buf.WriteString(",")
		}
		for _, v := range m.m {
			buf.WriteString(v[0].String())
			buf.WriteString("=")
			buf.WriteString(v[1].toString(lv + 1))
			buf.WriteString(",")
		}
		if m.Size() > 0 {
			buf.Truncate(buf.Len() - 1)
		}
		buf.WriteString("}")
		return buf.String()
	case Tbytes:
		arr := v.AsBytes()
		buf := &bytes.Buffer{}
		buf.WriteString("[")
		for _, v := range arr {
			buf.WriteString(fmt.Sprintf("%02x", int(v)))
			buf.WriteString(",")
		}
		if len(arr) > 0 {
			buf.Truncate(buf.Len() - 1)
		}
		buf.WriteString("]")
		return buf.String()
	case Tclosure:
		return v.AsClosure().String()
	case Tgeneric:
		return fmt.Sprintf("%v", v.AsGeneric())
	}
	return "nil"
}

func (v Value) panicType(expected byte) {
	log.Panicf("expecting %s, got %+v", TMapping[expected], v)
}

func testTypes(v1, v2 Value) uint16 {
	return uint16(v1.ty)<<8 + uint16(v2.ty)
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

func (v Value) Hash() hash128 {
	var a hash128
	switch v.ty {
	case Tnumber, Tbool, Tnil, Tclosure, Tlist, Tmap, Tgeneric:
		a = *(*hash128)(unsafe.Pointer(&v))
	case Tstring, Tbytes:
		if v.i > 0 {
			a = *(*hash128)(unsafe.Pointer(&v))
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

			a = hash128{h, h0}
		}
	}
	return a
}
