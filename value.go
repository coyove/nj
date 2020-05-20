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
// Note that THREE NaN values will not be valid: [0xffffffff`fffffffd, 0xffffffff`ffffffff]
type Value struct {
	v uint64
	p unsafe.Pointer
}

const SizeOfValue = unsafe.Sizeof(Value{})

// Type returns the type of value
func (v Value) Type() byte {
	if v.p == nil {
		if v.v != 0 {
			return NumberType
		}
		return NilType
	}
	return byte(v.v)
}

var (
	typeMappings = map[byte]string{
		NilType:     ("nil"),
		NumberType:  ("number"),
		StringType:  ("string"),
		ClosureType: ("closure"),
		PointerType: ("pointer"),
		SliceType:   ("slice"),
		StructType:  ("struct"),
	}

	_zero = NewNumberValue(0)
)

func init() {
	initCoreLibs()
}

// NewNumberValue returns a number value
func NewNumberValue(f float64) Value {
	x := *(*uint64)(unsafe.Pointer(&f))
	return Value{v: ^x}
}

// NewBoolValue returns a boolean value
func NewBoolValue(b bool) Value {
	x := float64(0)
	if b {
		x = 1.0
	}
	return Value{v: ^(*(*uint64)(unsafe.Pointer(&x)))}
}

// NewSliceValue returns a map value
func NewSliceValue(m *Slice) Value {
	return Value{v: SliceType, p: unsafe.Pointer(m)}
}

func NewStructValue(m *Struct) Value {
	return Value{v: StructType, p: unsafe.Pointer(m)}
}

// NewClosureValue returns a closure value
func NewClosureValue(c *Closure) Value {
	return Value{v: ClosureType, p: unsafe.Pointer(c)}
}

// NewPointerValue returns a generic value
func NewPointerValue(i interface{}) Value {
	return Value{v: PointerType, p: unsafe.Pointer(&i)}
}

// NewStringValue returns a string value
// Note we use []byte to avoid some unnecessary castings from string to []byte,
// it DOES NOT mean a StringValue is mutable
func NewStringValue(s string) Value {
	return Value{v: StringType, p: unsafe.Pointer(&s)}
}

func NewStringValueBytesUnsafe(s []byte) Value {
	return Value{v: StringType, p: unsafe.Pointer(&s)}
}

func NewInterfaceValue(i interface{}) Value {
	switch v := i.(type) {
	case float64:
		return NewNumberValue(v)
	case string:
		return NewStringValue(v)
	case *Slice:
		return NewSliceValue(v)
	case *Closure:
		return NewClosureValue(v)
	}
	return Value{v: PointerType, p: unsafe.Pointer(&i)}
}

// AsString cast value to string
func (v Value) AsString() string {
	return *(*string)(v.p)
}

// IsFalse tests whether value contains a "false" value
func (v Value) IsFalse() bool {
	switch v.Type() {
	case NumberType:
		return v.IsZero()
	case NilType:
		return true
	case StringType:
		m := (*string)(v.p)
		return len(*m) == 0
	case SliceType:
		m := (*Slice)(v.p)
		return len(m.l) == 0
	}
	return false
}

// IsZero is a fast way to check if a numeric Value is +0
func (v Value) IsZero() bool {
	return v == _zero
}

func (v Value) IsNil() bool {
	return v == Value{}
}

// AsNumber cast value to float64
func (v Value) AsNumber() float64 {
	return math.Float64frombits(^v.v)
}

func (v Value) AsInt32() int32 {
	return int32(int64(math.Float64frombits(^v.v)) & 0xffffffff)
}

// AsSlice cast value to map of values
func (v Value) AsSlice() *Slice {
	return (*Slice)(v.p)
}

func (v Value) AsStruct() *Struct {
	return (*Struct)(v.p)
}

// AsClosure cast value to closure
func (v Value) AsClosure() *Closure {
	return (*Closure)(v.p)
}

// AsPointer cast value to unsafe.Pointer
func (v Value) AsPointer() interface{} {
	return *(*interface{})(v.p)
}

// MustSlice safely cast value to map of values
func (v Value) MustSlice() *Slice {
	v.testType(SliceType)
	return (*Slice)(v.p)
}

func (v Value) MustStruct() *Struct {
	v.testType(StructType)
	return (*Struct)(v.p)
}

// MustClosure safely cast value to closure
func (v Value) MustClosure() *Closure {
	v.testType(ClosureType)
	return v.AsClosure()
}

// MustNumber safely cast value to float64
func (v Value) MustNumber() float64 {
	v.testType(NumberType)
	return v.AsNumber()
}

// MustString safely cast value to string
func (v Value) MustString() string {
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
		return v.AsPointer()
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
		return r.AsString() == v.AsString()
	case _SliceSlice:
		return v.AsSlice().Equal(r.AsSlice())
	case _StructStruct:
		return v.AsStruct().Equal(r.AsStruct())
	case _ClosureClosure:
		return v == r
	case _PointerPointer:
		return v.AsPointer() == r.AsPointer()
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
		if json {
			return fmt.Sprintf("\"%v\"", v.AsPointer())
		}
		return fmt.Sprintf("%v", v.AsPointer())
	}
	if json {
		return "null"
	}
	return "nil"
}

func (v Value) Dup() Value {
	switch v.Type() {
	case NilType, NumberType, PointerType, StringType:
		return v
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
