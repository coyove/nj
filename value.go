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
	Tnil = iota
	// Tnumber represents number type
	Tnumber
	// Tstring represents string type
	Tstring
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
	_Tmapmap         = Tmap<<8 | Tmap
	_Tclosureclosure = Tclosure<<8 | Tclosure
	_Tgenericgeneric = Tgeneric<<8 | Tgeneric
	_Tstringnumber   = Tstring<<8 | Tnumber
	_Tmapnumber      = Tmap<<8 | Tnumber
)

// Type returns the type of value
func (v Value) Type() byte {
	return v.ty
}

// NewValue returns a nil value
func NewValue() Value {
	return Value{ty: Tnil}
}

var (
	// Zero is +0
	Zero = NewNumberValue(0)

	// One is 1
	One = NewNumberValue(1)

	// TMapping maps type to its string representation
	TMapping = map[byte]string{
		Tnil: "nil", Tnumber: "number", Tstring: "string", Tclosure: "closure", Tgeneric: "generic", Tmap: "map",
	}

	hashkey [4]uintptr

	hash2Salt = rand.New().Fetch(16)
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
	v := Value{ty: Tnumber}
	*(*float64)(unsafe.Pointer(&v.ptr)) = f
	return v
}

// NewBoolValue returns a boolean value
func NewBoolValue(b bool) Value {
	if b {
		return One
	}
	return Zero
}

// NewMapValue returns a map value
func NewMapValue(m *Map) Value {
	return Value{ty: Tmap, ptr: unsafe.Pointer(m)}
}

// NewClosureValue returns a closure value
func NewClosureValue(c *Closure) Value {
	return Value{ty: Tclosure, ptr: unsafe.Pointer(c)}
}

// NewGenericValue returns a generic value
func NewGenericValue(g unsafe.Pointer) Value {
	return Value{ty: Tgeneric, ptr: g}
}

// AsMap cast value to map of values
func (v Value) AsMap() *Map { return (*Map)(v.ptr) }

// AsClosure cast value to closure
func (v Value) AsClosure() *Closure { return (*Closure)(v.ptr) }

// AsGeneric cast value to unsafe.Pointer
func (v Value) AsGeneric() unsafe.Pointer { return v.ptr }

// Map safely cast value to map of values
func (v Value) Map() *Map { v.testType(Tmap); return (*Map)(v.ptr) }

// Cls safely cast value to closure
func (v Value) Cls() *Closure { v.testType(Tclosure); return (*Closure)(v.ptr) }

// Gen safely cast value to unsafe.Pointer
func (v Value) Gen() unsafe.Pointer { v.testType(Tgeneric); return v.ptr }

func (v Value) u64() uint64 { return math.Float64bits(v.Num()) }

// Num safely cast value to float64
func (v Value) Num() float64 { v.testType(Tnumber); return v.AsNumber() }

// Str safely cast value to string
func (v Value) Str() string { v.testType(Tstring); return v.AsString() }

// I returns the golang interface representation of value
// it is not the same as AsGeneric()
func (v Value) I() interface{} {
	switch v.ty {
	case Tnumber:
		return v.AsNumber()
	case Tstring:
		return v.AsString()
	case Tmap:
		return v.AsMap()
	case Tclosure:
		return v.AsClosure()
	case Tgeneric:
		return v.AsGeneric()
	}
	return nil
}

func (v Value) String() string {
	switch v.ty {
	case Tstring:
		return strconv.Quote(v.AsString())
	default:
		return v.ToPrintString()
	}
}

// Equal tests whether value is equal to another value
// This is a strict test
func (v Value) Equal(r Value) bool {
	if v.ty == Tnil || r.ty == Tnil {
		return v.ty == r.ty
	}
	switch testTypes(v, r) {
	case _Tnumbernumber:
		return r.AsNumber() == v.AsNumber()
	case _Tstringstring:
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
			bytes.Equal(slice64to8(c0.code), slice64to8(c1.code)) &&
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
		return v.AsGeneric() == r.AsGeneric()
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
	case Tnumber:
		x := v.AsNumber()
		if float64(uint64(x)) == x {
			return strconv.FormatUint(uint64(x), 10)
		}
		return strconv.FormatFloat(x, 'f', -1, 64)
	case Tstring:
		return v.AsString()
	case Tmap:
		m, buf := v.AsMap(), &bytes.Buffer{}
		buf.WriteString("{")
		for _, v := range m.l {
			buf.WriteString(v.toString(lv + 1))
			buf.WriteString(",")
		}
		for _, v := range m.m {
			buf.WriteString(v[0].String())
			buf.WriteString(":")
			buf.WriteString(v[1].toString(lv + 1))
			buf.WriteString(",")
		}
		if m.Size() > 0 {
			buf.Truncate(buf.Len() - 1)
		}
		buf.WriteString("}")
		return buf.String()
	case Tclosure:
		return v.AsClosure().String()
	case Tgeneric:
		return fmt.Sprintf("%v", v.AsGeneric())
	}
	return "nil"
}

func (v Value) panicType(expected byte) {
	panicf("expecting %s, got %+v", TMapping[expected], v)
}

func (v Value) testType(expected byte) Value {
	if v.ty != expected {
		panicf("expecting %s, got %+v", TMapping[expected], v)
	}
	return v
}

func testTypes(v1, v2 Value) uint16 {
	return uint16(v1.ty)<<8 + uint16(v2.ty)
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
