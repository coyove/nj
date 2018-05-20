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

const (
	// the order can't be changed, for any new type, please also add it in parser.go.y typeof
	Tnil = iota
	Tnumber
	Tstring
	Tbool
	Tlist
	Tbytes
	Tmap
	Tclosure
	Tgeneric
)

var TMapping = map[byte]string{
	Tnil: "nil", Tnumber: "number", Tstring: "string", Tbool: "bool",
	Tclosure: "closure", Tgeneric: "generic", Tlist: "list", Tmap: "map", Tbytes: "bytes",
}

const (
	minPhysPageSize = 4096 // (0x1000) in mheap.go
	trueValue       = minPhysPageSize + 1
	falseValue      = minPhysPageSize + 2
)

// Value is the basic value used by VM
type Value struct {
	ty byte

	// used in red-black tree
	c color

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
	// i := *(*uintptr)(unsafe.Pointer(&f))
	// if i < minPhysPageSize {
	// 	i += minPhysPageSize
	// 	return Value{ty: Tnumber, ext: true, ptr: unsafe.Pointer(i)}
	// }
	// return Value{ty: Tnumber, ptr: unsafe.Pointer(i)}
	v := Value{ty: Tnumber}
	*(*float64)(unsafe.Pointer(&v.i)) = f
	return v
}

// NewStringValue returns a string value
func NewStringValue(s string) Value {
	v := NewValue()
	v.ty = Tstring

	if len(s) < 11 {
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
func NewMapValue(m *Tree) Value {
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

func (v Value) AsBool() bool {
	if v.ty != Tbool {
		log.Panicf("expecting boolean, got %+v", v)
	}
	return uintptr(v.ptr) == trueValue
}

func (v Value) AsBoolUnsafe() bool {
	return uintptr(v.ptr) == trueValue
}

func (v Value) AsNumber() float64 {
	if v.ty != Tnumber {
		log.Panicf("expecting number, got %+v", v)
	}
	return v.AsNumberUnsafe()
}

func (v Value) u64() uint64 {
	if v.ty != Tnumber {
		log.Panicf("expecting number, got %+v", v)
	}
	return math.Float64bits(v.AsNumberUnsafe())
}

func (v Value) AsNumberUnsafe() float64 {
	// i := uintptr(v.ptr)
	// if v.ext {
	// 	i -= minPhysPageSize
	// }
	// return *(*float64)(unsafe.Pointer(&i))
	return *(*float64)(unsafe.Pointer(&v.i))
}

func (v Value) AsString() string {
	if v.ty != Tstring {
		log.Panicf("expecting string, got %+v", v)
	}
	return v.AsStringUnsafe()
}

func (v Value) AsStringUnsafe() string {
	if v.i > 0 {
		hdr := reflect.StringHeader{}
		hdr.Len = int(v.i - 1)
		hdr.Data = uintptr(unsafe.Pointer(&v.p[0]))
		return *(*string)(unsafe.Pointer(&hdr))
	}
	return *(*string)(v.ptr)
}

func (v Value) AsList() []Value {
	if v.ty != Tlist {
		log.Panicf("expecting array, got %+v", v)
	}
	return *(*[]Value)(v.ptr)
}

func (v Value) AsListUnsafe() []Value {
	return *(*[]Value)(v.ptr)
}

func (v Value) AsMap() *Tree {
	if v.ty != Tmap {
		log.Panicf("expecting map, got %+v", v)
	}
	return (*Tree)(v.ptr)
}

func (v Value) AsMapUnsafe() *Tree {
	return (*Tree)(v.ptr)
}

func (v Value) AsClosure() *Closure {
	if v.ty != Tclosure {
		log.Panicf("expecting closure, got %+v", v)
	}
	return (*Closure)(v.ptr)
}

func (v Value) AsClosureUnsafe() *Closure {
	return (*Closure)(v.ptr)
}

func (v Value) AsGeneric() interface{} {
	if v.ty != Tgeneric {
		log.Panicf("expecting generic, got %+v", v)
	}
	return *(*interface{})(v.ptr)
}

func (v Value) AsGenericUnsafe() interface{} {
	return *(*interface{})(v.ptr)
}

func (v Value) AsBytes() []byte {
	if v.ty != Tbytes {
		log.Panicf("expecting bytes, got %+v", v)
	}
	if v.ptr == nil {
		return nil
	}
	return *(*[]byte)(v.ptr)
}

func (v Value) AsBytesUnsafe() []byte {
	return *(*[]byte)(v.ptr)
}

func (v Value) I() interface{} {
	switch v.Type() {
	case Tbool:
		return v.AsBoolUnsafe()
	case Tnumber:
		return v.AsNumberUnsafe()
	case Tstring:
		return v.AsStringUnsafe()
	case Tlist:
		return v.AsListUnsafe()
	case Tmap:
		return v.AsMapUnsafe()
	case Tbytes:
		return v.AsBytesUnsafe()
	case Tclosure:
		return v.AsClosureUnsafe()
	case Tgeneric:
		return v.AsGenericUnsafe()
	}
	return nil
}

func (v Value) String() string {
	switch v.Type() {
	case Tbool:
		return "<bool:" + strconv.FormatBool(v.AsBoolUnsafe()) + ">"
	case Tnumber:
		return "<number:" + strconv.FormatFloat(v.AsNumberUnsafe(), 'f', 9, 64) + ">"
	case Tstring:
		return "<string:" + strconv.Quote(v.AsStringUnsafe()) + ">"
	case Tlist:
		return "<list:[" + strconv.Itoa(len(v.AsListUnsafe())) + "]>"
	case Tmap:
		return "<map:[" + strconv.Itoa(v.AsMapUnsafe().Size()) + "]>"
	case Tbytes:
		return "<bytes:[" + strconv.Itoa(len(v.AsBytesUnsafe())) + "]>"
	case Tclosure:
		return "<closure:[" + strconv.Itoa(v.AsClosureUnsafe().ArgsCount()) +
			"/" + strconv.Itoa(len(v.AsClosureUnsafe().preArgs)) + "]>"
	case Tgeneric:
		return fmt.Sprintf("<generic:%+v>", v.AsGenericUnsafe())
	}
	return "<nil>"
}

func (v Value) IsFalse() bool {
	switch v.Type() {
	case Tnil:
		return true
	case Tbool:
		return v.AsBoolUnsafe() == false
	case Tnumber:
		return v.AsNumberUnsafe() == 0.0
	case Tstring:
		return v.AsStringUnsafe() == ""
	case Tlist:
		return len(v.AsListUnsafe()) == 0
	case Tbytes:
		return len(v.AsBytesUnsafe()) == 0
	case Tmap:
		return v.AsMapUnsafe().Size() == 0
	}
	return false
}

func (v Value) Equal(r Value) bool {
	if v.ty == Tnil || r.ty == Tnil {
		return v.ty == r.ty
	}

	switch v.ty {
	case Tnumber:
		if r.ty == Tnumber {
			return r.AsNumberUnsafe() == v.AsNumberUnsafe()
		}
	case Tstring:
		if r.ty == Tstring {
			return r.AsStringUnsafe() == v.AsStringUnsafe()
		} else if r.ty == Tbytes {
			return bytes.Equal(r.AsBytesUnsafe(), []byte(v.AsStringUnsafe()))
		}
	case Tbool:
		if r.ty == Tbool {
			return r.AsBoolUnsafe() == v.AsBoolUnsafe()
		}
	case Tlist:
		if r.ty == Tlist {
			lf, rf := v.AsListUnsafe(), r.AsListUnsafe()

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
			return bytes.Equal(v.AsBytesUnsafe(), r.AsBytesUnsafe())
		} else if r.ty == Tstring {
			return bytes.Equal(v.AsBytesUnsafe(), []byte(r.AsStringUnsafe()))
		}
	}

	return false
}

func (v Value) Less(r Value) bool {
	switch v.ty {
	case Tnumber:
		if r.ty == Tnumber {
			return v.AsNumberUnsafe() < r.AsNumberUnsafe()
		}
	case Tstring:
		if r.ty == Tstring {
			return v.AsStringUnsafe() < r.AsStringUnsafe()
		}
	}
	log.Panicf("can't compare %+v and %+v", v, r)
	return false
}

func (v Value) LessEqual(r Value) bool {
	switch v.ty {
	case Tnumber:
		if r.ty == Tnumber {
			return v.AsNumberUnsafe() <= r.AsNumberUnsafe()
		}
	case Tstring:
		if r.ty == Tstring {
			return v.AsStringUnsafe() <= r.AsStringUnsafe()
		}
	}
	log.Panicf("can't compare %+v and %+v", v, r)
	return false
}

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
		n := v.AsNumber()
		if float64(int64(n)) == n {
			return strconv.FormatInt(int64(n), 10)
		}
		return strconv.FormatFloat(n, 'f', 9, 64)
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
		m := v.AsMap()
		buf, iter := &bytes.Buffer{}, m.Iterator()
		buf.WriteString("{")
		for iter.Next() {
			buf.WriteString(iter.Key())
			buf.WriteString(":")
			buf.WriteString(iter.Value().toString(lv + 1))
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
		return "<" + v.AsClosure().String() + ">"
	case Tgeneric:
		return fmt.Sprintf("%v", v.AsGenericUnsafe())
	}
	return "nil"
}
