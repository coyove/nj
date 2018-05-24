package potatolang

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"
	"unsafe"
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
)

// TMapping maps type to its string representation
var TMapping = map[byte]string{
	Tnil: "nil", Tnumber: "number", Tstring: "string", Tbool: "bool",
	Tclosure: "closure", Tgeneric: "generic", Tlist: "list", Tmap: "map", Tbytes: "bytes",
}

var safePointerAddr = unsafe.Pointer(uintptr(0x4000000000000000))

const (
	minPhysPageSize = 4096 // (0x1000) in mheap.go
	trueValue       = minPhysPageSize + 1
	falseValue      = minPhysPageSize + 2
)

// Value is the basic value used by VM
type Value struct {
	ty byte

	// float64 will be stored at &i
	// it can't be stored in ptr because pointer value smaller than minPhysPageSize will violate the heap
	// for ty == Tstring and len(str) <= 10, i == len(str) + 1
	i byte

	// for ty == Tstring and i > 0, 10 bytes starting from p[0] will be used to store small strings
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
func (v Value) AsNumber() float64 { return *(*float64)(unsafe.Pointer(&v.i)) }

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
		for _, x := range m.t {
			buf.WriteString(x.k)
			buf.WriteString(":")
			buf.WriteString(x.v.toString(lv + 1))
			buf.WriteString(",")
		}
		for k, v := range m.m {
			buf.WriteString(k)
			buf.WriteString(":")
			buf.WriteString(v.toString(lv + 1))
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
