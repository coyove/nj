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
	NilType = 0

	// NumberType represents number type
	NumberType = 1

	// StringType represents string type
	StringType = 2

	//
	StructType = 3

	// SliceType represents map type
	SliceType = 4

	// ClosureType represents closure type
	ClosureType = 6

	// PointerType represents generic type
	PointerType = 7
)

const (
	_NilNil         = NilType<<8 | NilType
	_NumberNumber   = NumberType<<8 | NumberType
	_StringString   = StringType<<8 | StringType
	_SliceSlice     = SliceType<<8 | SliceType
	_StructStruct   = StructType<<8 | StructType
	_ClosureClosure = ClosureType<<8 | ClosureType
	_PointerPointer = PointerType<<8 | PointerType
	_StringNumber   = StringType<<8 | NumberType
	_SliceNumber    = SliceType<<8 | NumberType
	_StructNumber   = StructType<<8 | NumberType
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
	if x > 0xffffffffffff {
		return NumberType
	}

	if x == 0 {
		return NilType
	}

	m := (*baseSlice)(unsafe.Pointer(x))
	return m.ptype
}

var (
	// TMapping maps type to its string representation
	TMapping = map[byte]string{
		NilType: "nil", NumberType: "num", StringType: "str", ClosureType: "cls", PointerType: "ptr", SliceType: "map",
	}

	Phantom = NewPointerValue(unsafe.Pointer(new(int)), PTagPhantom)

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

// NewMapValue returns a map value
func NewMapValue(m *baseSlice) Value {
	m.ptype = SliceType
	return Value{ptr: unsafe.Pointer(m)}
}

func NewStructValue(m *baseStruct) Value {
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
	m := &basePointer{ptype: PointerType, ptr: g, ptag: tag}
	return Value{unsafe.Pointer(m)}
}

// NewStringValue returns a string value
func NewStringValue(s string) Value {
	m := &baseString{ptype: StringType, s: s}
	return Value{unsafe.Pointer(m)}
}

func NewInterfaceValue(i interface{}) Value {
	switch v := i.(type) {
	case float64:
		return NewNumberValue(v)
	case string:
		return NewStringValue(v)
	case *baseSlice:
		return NewMapValue(v)
	case *Closure:
		return NewClosureValue(v)
	}
	return Value{}
}

// AsString cast value to string
func (v Value) AsString() string {
	return (*baseString)(v.ptr).s
}

// IsFalse tests whether value contains a "false" value
func (v Value) IsFalse() bool {
	switch v.Type() {
	case NumberType:
		return v.IsZero()
	case NilType:
		return true
	case StringType:
		m := (*baseString)(v.ptr)
		return len(m.s) == 0
	case SliceType:
		m := (*baseSlice)(v.ptr)
		return len(m.l) == 0
	}
	return v == Phantom
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
func (v Value) AsSlice() *baseSlice {
	return (*baseSlice)(v.ptr)
}

func (v Value) AsStruct() *baseStruct {
	return (*baseStruct)(v.ptr)
}

// AsClosure cast value to closure
func (v Value) AsClosure() *Closure {
	return (*Closure)(v.ptr)
}

// AsPointer cast value to unsafe.Pointer
func (v Value) AsPointer() (unsafe.Pointer, uint32) {
	return (*basePointer)(v.ptr).ptr, (*basePointer)(v.ptr).ptag
}

// MustMap safely cast value to map of values
func (v Value) MustMap() *baseSlice { v.testType(SliceType); return (*baseSlice)(v.ptr) }

// MustClosure safely cast value to closure
func (v Value) MustClosure() *Closure { v.testType(ClosureType); return v.AsClosure() }

// MustPointer safely cast value to unsafe.Pointer
func (v Value) MustPointer() (unsafe.Pointer, uint32) { v.testType(PointerType); return v.AsPointer() }

func (v Value) u64() uint64 { return math.Float64bits(v.MustNumber()) }

// MustNumber safely cast value to float64
func (v Value) MustNumber() float64 { v.testType(NumberType); return v.AsNumber() }

// MustString safely cast value to string
func (v Value) MustString() string { v.testType(StringType); return v.AsString() }

// AsInterface returns the golang interface representation of value
func (v Value) AsInterface() interface{} {
	switch v.Type() {
	case NumberType:
		return v.AsNumber()
	case StringType:
		return v.AsString()
	case SliceType:
		return v.AsSlice()
	case ClosureType:
		return v.AsClosure()
	}
	return nil
}

func (v Value) String() string {
	return v.toString(0, false)
}

// Equal tests whether value is equal to another value
// This is a strict test
func (v Value) Equal(r Value) bool {
	switch combineTypes(v, r) {
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
		c0, c1 := v.AsClosure(), r.AsClosure()
		e := c0.ArgsCount == c1.ArgsCount &&
			c0.options == c1.options &&
			c0.Env == c1.Env &&
			c0.lastenv == c1.lastenv &&
			c0.lastp == c1.lastp &&
			bytes.Equal(u32Bytes(c0.Code), u32Bytes(c1.Code)) &&
			len(c0.PartialArgs) == len(c1.PartialArgs)
		if !e {
			return false
		}
		for i, arg := range c0.PartialArgs {
			if !arg.Equal(c1.PartialArgs[i]) {
				return false
			}
		}
		return true
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
		x := v.AsNumber()
		if float64(uint64(x)) == x {
			return strconv.FormatUint(uint64(x), 10)
		}
		return strconv.FormatFloat(x, 'f', -1, 64)
	case StringType:
		if json {
			return strconv.Quote(v.AsString())
		}
		return v.AsString()
	case SliceType:
		m, buf := v.AsSlice(), &bytes.Buffer{}
		if json {
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
			buf.WriteString("{")
			for _, v := range m.l {
				buf.WriteString(v.String())
				buf.WriteString(",")
			}
			if m.Size() > 0 {
				buf.Truncate(buf.Len() - 1)
			}
			buf.WriteString("}")
		}
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
	case NilType, NumberType, StringType, PointerType:
		return v
	case ClosureType:
		return NewClosureValue(v.AsClosure().Dup())
	case SliceType:
		return NewMapValue(v.AsSlice().Dup())
	default:
		panic("unreachable Code")
	}
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

func combineTypes(v1, v2 Value) uint16 {
	return uint16(v1.Type())<<8 + uint16(v2.Type())
}
