package potatolang

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"unsafe"

	"github.com/coyove/common/rand"
)

// the order can't be changed, for any new type, please also add it in parser.go.y typeof
const (
	// Tnil represents nil type
	Tnil = 0
	// Tnumber represents number type
	Tnumber = 1
	// Tstring represents string type
	Tstring = 2
	// Tmap represents map type
	Tmap = 4
	// Tclosure represents closure type
	Tclosure = 6
	// Tgeneric represents generic type
	Tgeneric = 7
)

const (
	_Tnilnil         = Tnil<<8 | Tnil
	_Tnumbernumber   = Tnumber<<8 | Tnumber
	_Tstringstring   = Tstring<<8 | Tstring
	_Tmapmap         = Tmap<<8 | Tmap
	_Tclosureclosure = Tclosure<<8 | Tclosure
	_Tgenericgeneric = Tgeneric<<8 | Tgeneric
	_Tstringnumber   = Tstring<<8 | Tnumber
	_Tmapnumber      = Tmap<<8 | Tnumber
)

// Value is the basic value used by VM
type Value struct {
	ptr unsafe.Pointer // 8b
	num uint64
}

// Type returns the type of value
func (v Value) Type() byte {
	// return byte(v.num&0xf) & ^(0xe * byte(v.num&1))
	return byte(v.num & 0x7)
}

const clear3LSB = 0xfffffffffffffff8

// NewValue returns a nil value
func NewValue() Value { return Value{num: 0} }

var (
	// TMapping maps type to its string representation
	TMapping = map[byte]string{
		Tnil: "nil", Tnumber: "num", Tstring: "str", Tclosure: "cls", Tgeneric: "gen", Tmap: "map",
	}

	hashkey   [4]uintptr
	hash2Salt = rand.New().Fetch(16)

	// for a numeric Value, its 'ptr' = _Anchor64 + last 3 bits of 'num'
	_anchor64 = new(int64)
	_Anchor64 = uintptr(unsafe.Pointer(_anchor64))

	// PhantomValue is a global readonly value to represent the true "void"
	PhantomValue = NewGenericValue(unsafe.Pointer(_anchor64), GTagPhantom)
)

func init() {
	buf := rand.New().Fetch(32)
	for i := 0; i < 4; i++ {
		hashkey[i] = uintptr(binary.LittleEndian.Uint64(buf[i*8:]))
		hashkey[i] |= 1
	}
	initCoreLibs()
}

// NewNumberValue returns a number value
func NewNumberValue(f float64) Value {
	v := Value{}
	x := math.Float64bits(f)
	v.ptr = unsafe.Pointer(uintptr(x&7) + _Anchor64)
	v.num = x&clear3LSB + Tnumber
	return v
}

// NewBoolValue returns a boolean value
func NewBoolValue(b bool) Value {
	v := Value{}
	x := uint64(*(*byte)(unsafe.Pointer(&b)))
	v.ptr = unsafe.Pointer(_Anchor64)
	v.num = x*0x3ff0000000000000 + Tnumber
	return v
}

// SetNumberValue turns any Value into a numeric Value
func (v *Value) SetNumberValue(f float64) {
	x := math.Float64bits(f)
	v.ptr = unsafe.Pointer(uintptr(x&7) + _Anchor64)
	v.num = x&clear3LSB + Tnumber
}

// SetBoolValue turns any Value into a numeric Value with its value being 0.0 or 1.0
func (v *Value) SetBoolValue(b bool) {
	x := uint64(*(*byte)(unsafe.Pointer(&b)))
	v.ptr = unsafe.Pointer(_Anchor64)
	v.num = x*0x3ff0000000000000 + Tnumber
}

// NewMapValue returns a map value
func NewMapValue(m *Map) Value {
	return Value{num: Tmap, ptr: unsafe.Pointer(m)}
}

// NewClosureValue returns a closure value
func NewClosureValue(c *Closure) Value {
	return Value{num: Tclosure, ptr: unsafe.Pointer(c)}
}

// NewGenericValue returns a generic value
func NewGenericValue(g unsafe.Pointer, tag uint32) Value {
	return Value{num: uint64(tag)<<32 + Tgeneric, ptr: g}
}

// NewGenericValueInterface returns a generic value from an interface{}
func NewGenericValueInterface(i interface{}, tag uint32) Value {
	g := (*(*[2]unsafe.Pointer)(unsafe.Pointer(&i)))[1]
	return Value{num: uint64(tag)<<32 + Tgeneric, ptr: g}
}

// NewStringValue returns a string value
func NewStringValue(s string) Value {
	v := Value{}
	if len(s) < 7 {
		// for a string containing only 0~6 bytes, it will be stored into Value directly
		copy((*(*[8]byte)(unsafe.Pointer(&v.num)))[1:7], s)
		v.num |= uint64(Tstring + (len(s)+1)<<4)
	} else {
		v.num = Tstring
		v.ptr = unsafe.Pointer(&s)
	}
	return v
}

// AsString cast value to string
func (v Value) AsString() string {
	if ln := byte(v.num) >> 4; ln > 0 {
		buf := make([]byte, ln-1)
		copy(buf, (*(*[8]byte)(unsafe.Pointer(&v.num)))[1:ln])
		return string(buf)
	}
	return *(*string)(v.ptr)
}

// IsFalse tests whether value contains a "false" value
func (v Value) IsFalse() bool {
	if v.Type() == Tnumber {
		return v.IsZero()
	}
	if v.Type() == Tnil {
		return true
	}
	if v.Type() == Tstring {
		return byte(v.num)>>4 == 1
	}
	if v.Type() == Tmap {
		m := (*Map)(v.ptr)
		return len(m.l)+len(m.m) == 0
	}
	return false
}

var _zero = NewNumberValue(0)

// IsZero is a fast way to check if a numeric Value is +0
func (v Value) IsZero() bool {
	return v == _zero
}

// AsNumber cast value to float64
func (v Value) AsNumber() float64 {
	return math.Float64frombits(v.num&clear3LSB + uint64(uintptr(v.ptr)-_Anchor64))
}

// AsMap cast value to map of values
func (v Value) AsMap() *Map { return (*Map)(v.ptr) }

// AsClosure cast value to closure
func (v Value) AsClosure() *Closure { return (*Closure)(v.ptr) }

// AsGeneric cast value to unsafe.Pointer
func (v Value) AsGeneric() (unsafe.Pointer, uint32) { return v.ptr, uint32(v.num >> 32) }

// Map safely cast value to map of values
func (v Value) Map() *Map { v.testType(Tmap); return (*Map)(v.ptr) }

// Cls safely cast value to closure
func (v Value) Cls() *Closure { v.testType(Tclosure); return (*Closure)(v.ptr) }

// Gen safely cast value to unsafe.Pointer
func (v Value) Gen() (unsafe.Pointer, uint32) { v.testType(Tgeneric); return v.AsGeneric() }

func (v Value) GenTags(tags ...uint32) unsafe.Pointer {
	v.testType(Tgeneric)
	vp, vt := v.AsGeneric()
	for _, tag := range tags {
		if vt == tag {
			return vp
		}
	}
	panicf("expecting tags: %v, got %d", tags, vt)
	return vp
}

func (v Value) u64() uint64 { return math.Float64bits(v.Num()) }

// Num safely cast value to float64
func (v Value) Num() float64 { v.testType(Tnumber); return v.AsNumber() }

// Str safely cast value to string
func (v Value) Str() string { v.testType(Tstring); return v.AsString() }

// I returns the golang interface representation of value
// Tgeneric will not be returned, use Gen() instead
func (v Value) I() interface{} {
	switch v.Type() {
	case Tnumber:
		return v.AsNumber()
	case Tstring:
		return v.AsString()
	case Tmap:
		return v.AsMap()
	case Tclosure:
		return v.AsClosure()
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

// Equal tests whether value is equal to another value
// This is a strict test
func (v Value) Equal(r Value) bool {
	if v.Type() == Tnil || r.Type() == Tnil {
		return v.Type() == r.Type()
	}
	switch testTypes(v, r) {
	case _Tnumbernumber:
		return v == r
	case _Tstringstring:
		if ln := byte(v.num) >> 4; ln > 0 {
			return v == r
		}
		return r.AsString() == v.AsString()
	case _Tmapmap:
		return v.AsMap().Equal(r.AsMap())
	case _Tclosureclosure:
		c0, c1 := v.AsClosure(), r.AsClosure()
		e := c0.argsCount == c1.argsCount &&
			c0.options == c1.options &&
			c0.env == c1.env &&
			c0.lastenv == c1.lastenv &&
			c0.lastp == c1.lastp &&
			bytes.Equal(u32Bytes(c0.code), u32Bytes(c1.code)) &&
			c0.caller.Equal(c1.caller) &&
			len(c0.preArgs) == len(c1.preArgs)
		if !e {
			return false
		}
		for i, arg := range c0.preArgs {
			if !arg.Equal(c1.preArgs[i]) {
				return false
			}
		}
		return true
	case _Tgenericgeneric:
		vp, vt := v.AsGeneric()
		rp, rt := r.AsGeneric()
		eq := gtagComparators[uint64(vt)<<32+uint64(rt)]
		if eq != nil {
			return eq.Equal(v, r)
		}
		return vp == rp && vt == rt
	}
	return false
}

// ToPrintString returns the printable string of value
// it won't wrap a string with double quotes, String() will
func (v Value) ToPrintString() string {
	return v.toString(0, false)
}

func (v Value) toString(lv int, json bool) string {
	if lv > 32 {
		if json {
			return "\"<omit deep nesting>\""
		}
		return "<omit deep nesting>"
	}

	switch v.Type() {
	case Tnumber:
		x := v.AsNumber()
		if float64(uint64(x)) == x {
			return strconv.FormatUint(uint64(x), 10)
		}
		return strconv.FormatFloat(x, 'f', -1, 64)
	case Tstring:
		if json {
			return "\"" + strconv.Quote(v.AsString()) + "\""
		}
		return v.AsString()
	case Tmap:
		m, buf := v.AsMap(), &bytes.Buffer{}
		if json {
			if len(m.m) == 0 {
				// treat it as an array
				buf.WriteString("[")
				for _, v := range m.l {
					buf.WriteString(v.toString(lv+1, json))
					buf.WriteString(",")
				}
				if len(m.l) > 0 {
					buf.Truncate(buf.Len() - 1)
				}
				buf.WriteString("]")
			} else {
				// treat it as an object
				buf.WriteString("{")
				for i, v := range m.l {
					buf.WriteString("\"" + strconv.Itoa(i) + "\":")
					buf.WriteString(v.toString(lv+1, json))
					buf.WriteString(",")
				}
				for _, v := range m.m {
					if v[0].Type() != Tstring {
						panicf("non-string key is not allowed")
					}
					buf.WriteString(v[0].String())
					buf.WriteString(":")
					buf.WriteString(v[1].toString(lv+1, json))
					buf.WriteString(",")
				}
				if m.Size() > 0 {
					buf.Truncate(buf.Len() - 1)
				}
				buf.WriteString("}")
			}
		} else {
			buf.WriteString("{")
			for _, v := range m.l {
				buf.WriteString(v.toString(lv+1, json))
				buf.WriteString(",")
			}
			for _, v := range m.m {
				buf.WriteString(v[0].String())
				buf.WriteString(":")
				buf.WriteString(v[1].toString(lv+1, json))
				buf.WriteString(",")
			}
			if m.Size() > 0 {
				buf.Truncate(buf.Len() - 1)
			}
			buf.WriteString("}")
		}
		return buf.String()
	case Tclosure:
		if json {
			return "\"" + v.AsClosure().String() + "\""
		}
		return v.AsClosure().String()
	case Tgeneric:
		vp, vt := v.AsGeneric()
		if json {
			return fmt.Sprintf("\"<tag%x:%v>\"", vt, vp)
		}
		return fmt.Sprintf("<tag%x:%v>", vt, vp)
	}
	if json {
		return "null"
	}
	return "nil"
}

func (v Value) panicType(expected byte) {
	panicf("expecting %s, got %+v", TMapping[expected], v)
}

func (v Value) testType(expected byte) Value {
	if v.Type() != expected {
		panicf("expecting %s, got %+v", TMapping[expected], v)
	}
	return v
}

func testTypes(v1, v2 Value) uint16 {
	return uint16(v1.Type())<<8 + uint16(v2.Type())
}

//go:nosplit
func add(p unsafe.Pointer, x uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

func readUnaligned32(p unsafe.Pointer) uint32 {
	return *(*uint32)(p)
}

func readUnaligned64(p unsafe.Pointer) uint64 {
	return *(*uint64)(p)
}
