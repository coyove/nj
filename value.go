package potatolang

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"unsafe"
)

const (
	// NilType represents nil type
	NilType = 1

	// NumberType represents number type
	NumberType = 3

	// StringType represents string type
	StringType = 7

	//
	StructType = 15

	// SliceType represents map type
	SliceType = 31

	// ClosureType represents closure type
	ClosureType = 63

	// PointerType represents generic type
	PointerType = 127
)

const (
	_NilNil         = NilType * 2
	_NumberNumber   = NumberType * 2
	_StringString   = StringType * 2
	_SliceSlice     = SliceType * 2
	_StructStruct   = StructType * 2
	_ClosureClosure = ClosureType * 2
	_PointerPointer = PointerType * 2
)

// Value is the basic value used by VM
// It assumes the OS will not map any memory higher than 1 << 48
// Some valid NaN value will not be valid in Value struct
// TODO: 32bit support with padding bytes
type Value struct {
	ptr unsafe.Pointer // 8b
}

const SizeOfValue = unsafe.Sizeof(Value{})

// Type returns the type of value
func (v Value) Type() byte {
	x := uintptr(v.ptr)
	if x >= 1<<48 {
		return NumberType
	}
	if x == 0 {
		return NilType
	}
	return (*Base)(unsafe.Pointer(x)).ptype
}

var (
	typeMappings = map[byte][]byte{
		NilType:     []byte("nil"),
		NumberType:  []byte("number"),
		StringType:  []byte("string"),
		ClosureType: []byte("closure"),
		PointerType: []byte("pointer"),
		SliceType:   []byte("slice"),
		StructType:  []byte("struct"),
	}

	_zero = NewNumberValue(0)
)

func init() {
	initCoreLibs()
}

// NewNumberValue returns a number value
func NewNumberValue(f float64) Value {
	x := *(*uint64)(unsafe.Pointer(&f))
	return Value{unsafe.Pointer(^uintptr(x))}
}

// NewBoolValue returns a boolean value
func NewBoolValue(b bool) Value {
	x := float64(0)
	if b {
		x = 1.0
	}
	return Value{unsafe.Pointer(^uintptr(*(*uint64)(unsafe.Pointer(&x))))}
}

// SetNumberValue turns any Value into a numeric Value
func (v *Value) SetNumberValue(f float64) {
	x := *(*uint64)(unsafe.Pointer(&f))
	//if x>>52 == 0xfff && x<<12 > 0 {
	//	x = math.MaxUint64
	//}
	v.ptr = unsafe.Pointer(^uintptr(x))
}

// SetBoolValue turns any Value into a numeric Value with its value being 0.0 or 1.0
func (v *Value) SetBoolValue(b bool) {
	x := 0.0
	if b {
		x = 1.0
	}
	v.ptr = unsafe.Pointer(^uintptr(*(*uint64)(unsafe.Pointer(&x))))
}

// NewSliceValue returns a map value
func NewSliceValue(m *Slice) Value {
	m.ptype = SliceType
	return Value{ptr: unsafe.Pointer(m)}
}

func NewStructValue(m *Struct) Value {
	m.ptype = StructType
	return Value{ptr: unsafe.Pointer(m)}
}

// NewClosureValue returns a closure value
func NewClosureValue(c *Closure) Value {
	c.ptype = ClosureType
	return Value{unsafe.Pointer(c)}
}

// NewPointerValue returns a generic value
func NewPointerValue(g unsafe.Pointer, tag uint32) Value {
	m := &Pointer{Base: Base{ptype: PointerType, ptag: tag}, ptr: g}
	return Value{unsafe.Pointer(m)}
}

// NewStringValue returns a string value
// Note we use []byte to avoid some unnecessary castings from string to []byte,
// it DOES NOT mean a StringValue is mutable
func NewStringValue(s []byte) Value {
	m := &String{Base: Base{ptype: StringType}, s: s}
	return Value{unsafe.Pointer(m)}
}

func NewStringValueString(s string) Value {
	return NewStringValue([]byte(s))
}

func NewInterfaceValue(i interface{}) Value {
	switch v := i.(type) {
	case float64:
		return NewNumberValue(v)
	case string:
		return NewStringValueString(v)
	case []byte:
		return NewStringValue(v)
	case *Slice:
		return NewSliceValue(v)
	case *Closure:
		return NewClosureValue(v)
	}
	m := &Pointer{Base: Base{ptype: PointerType, ptag: PTagInterface}, ptr: unsafe.Pointer(&i)}
	return Value{unsafe.Pointer(m)}
}

// AsString cast value to string
func (v Value) AsString() []byte {
	return (*String)(v.ptr).s
}

// IsFalse tests whether value contains a "false" value
func (v Value) IsFalse() bool {
	switch v.Type() {
	case NumberType:
		return v.IsZero()
	case NilType:
		return true
	case StringType:
		m := (*String)(v.ptr)
		return len(m.s) == 0
	case SliceType:
		m := (*Slice)(v.ptr)
		return len(m.l) == 0
	}
	return false
}

// IsZero is a fast way to check if a numeric Value is +0
func (v Value) IsZero() bool {
	return v == _zero
}

// AsNumber cast value to float64
func (v Value) AsNumber() float64 {
	return math.Float64frombits(^uint64(uintptr(v.ptr)))
}

func (v Value) AsInt32() int32 {
	return int32(int64(math.Float64frombits(^uint64(uintptr(v.ptr)))) & 0xffffffff)
}

// AsSlice cast value to map of values
func (v Value) AsSlice() *Slice {
	return (*Slice)(v.ptr)
}

func (v Value) AsStruct() *Struct {
	return (*Struct)(v.ptr)
}

// AsClosure cast value to closure
func (v Value) AsClosure() *Closure {
	return (*Closure)(v.ptr)
}

// AsPointer cast value to unsafe.Pointer
func (v Value) AsPointer() (unsafe.Pointer, uint32) {
	return (*Pointer)(v.ptr).ptr, (*Pointer)(v.ptr).ptag
}

// MustSlice safely cast value to map of values
func (v Value) MustSlice() *Slice {
	v.testType(SliceType)
	return (*Slice)(v.ptr)
}

func (v Value) MustStruct() *Struct {
	v.testType(StructType)
	return (*Struct)(v.ptr)
}

// MustClosure safely cast value to closure
func (v Value) MustClosure() *Closure {
	v.testType(ClosureType)
	return v.AsClosure()
}

// MustPointer safely cast value to unsafe.Pointer
func (v Value) MustPointer(tag uint32) unsafe.Pointer {
	v.testType(PointerType)
	p, t := v.AsPointer()
	if t != tag {
		panicf("expecting %x, got %x", tag, t)
	}
	return p
}

// MustNumber safely cast value to float64
func (v Value) MustNumber() float64 {
	v.testType(NumberType)
	return v.AsNumber()
}

// MustString safely cast value to string
func (v Value) MustString() []byte {
	v.testType(StringType)
	return v.AsString()
}

// AsInterface returns the golang interface representation of value
func (v Value) AsInterface() interface{} {
	switch v.Type() {
	case NumberType:
		return v.AsNumber()
	case StringType:
		return string(v.AsString())
	case SliceType:
		return v.AsSlice()
	case ClosureType:
		return v.AsClosure()
	case PointerType:
		ptr, ptag := v.AsPointer()
		if ptag == PTagInterface {
			return *(*interface{})(ptr)
		}
	}
	return nil
}

func (v Value) String() string {
	return v.toString(0, false)
}

func (v Value) GoString() string {
	return v.toString(0, true)
}

// Equal tests whether value is equal to another value
func (v Value) Equal(r Value) bool {
	switch v.Type() + r.Type() {
	case _NilNil:
		return true
	case _NumberNumber:
		return v == r
	case _StringString:
		return bytes.Equal(r.AsString(), v.AsString())
	case _SliceSlice:
		return v.AsSlice().Equal(r.AsSlice())
	case _StructStruct:
		return v.AsStruct().Equal(r.AsStruct())
	case _ClosureClosure:
		return v == r
	case _PointerPointer:
		vp, vt := v.AsPointer()
		rp, rt := r.AsPointer()
		return vp == rp && vt == rt
	}
	return false
}

func (v Value) toString(lv int, json bool) string {
	if lv > 32 {
		if json {
			return "\"<omit deep nesting>\""
		}
		return "<omit deep nesting>"
	}

	switch v.Type() {
	case NumberType:
		return strconv.FormatFloat(v.AsNumber(), 'f', -1, 64)
	case StringType:
		if json {
			return strconv.Quote(string(v.AsString()))
		}
		return string(v.AsString())
	case SliceType:
		buf := bytes.Buffer{}
		buf.WriteString("[")
		for _, v := range v.AsSlice().l {
			buf.WriteString(v.toString(lv+1, json))
			buf.WriteString(",")
		}
		if len(v.AsSlice().l) > 0 {
			buf.Truncate(buf.Len() - 1)
		}
		buf.WriteString("]")
		return buf.String()
	case StructType:
		if json {
			return "{}"
		}
		return "<struct>"
	case ClosureType:
		if json {
			return "\"" + v.AsClosure().String() + "\""
		}
		return v.AsClosure().String()
	case PointerType:
		vp, vt := v.AsPointer()
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

func (v Value) Dup() Value {
	switch v.Type() {
	case NilType, NumberType, PointerType:
		return v
	case StringType:
		return NewStringValue(append([]byte{}, v.AsString()...))
	case ClosureType:
		return NewClosureValue(v.AsClosure().Dup())
	case SliceType:
		return NewSliceValue(v.AsSlice().Dup())
	case StructType:
		return NewStructValue(v.AsStruct().Dup())
	default:
		panic("unreachable Code")
	}
}

func (v Value) testType(expected byte) Value {
	if v.Type() != expected {
		panicf("expecting %q, got %+v", typeMappings[expected], v)
	}
	return v
}
