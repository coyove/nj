package base

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"strconv"
	"unsafe"
)

const (
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

const (
	minPhysPageSize = 4096 // (0x1000) in mheap.go
	trueValue       = minPhysPageSize + 1
	falseValue      = minPhysPageSize + 2
)

type Value struct {
	ty byte

	// if set true, ptr = ptr + minPhysPageSize
	// to get around of runtime barrier check
	ext bool

	// used in red-black tree
	c color

	p3 byte
	p4 byte
	p5 byte
	p6 byte
	p7 byte

	ptr unsafe.Pointer
}

func NewValue() Value {
	return Value{ty: Tnil}
}

func NewNumberValue(f float64) Value {
	i := *(*uintptr)(unsafe.Pointer(&f))
	if i < minPhysPageSize {
		i += minPhysPageSize
		return Value{ty: Tnumber, ext: true, ptr: unsafe.Pointer(i)}
	}
	return Value{ty: Tnumber, ptr: unsafe.Pointer(i)}
}

func NewStringValue(s string) Value {
	return Value{ty: Tstring, ptr: unsafe.Pointer(&s)}
}

func NewBoolValue(b bool) Value {
	v := Value{ty: Tbool}
	if b {
		v.ptr = unsafe.Pointer(uintptr(trueValue))
	} else {
		v.ptr = unsafe.Pointer(uintptr(falseValue))
	}
	return v
}

func NewListValue(a []Value) Value {
	return Value{ty: Tlist, ptr: unsafe.Pointer(&a)}
}

func NewMapValue(m *Tree) Value {
	return Value{ty: Tmap, ptr: unsafe.Pointer(m)}
}

func NewClosureValue(c *Closure) Value {
	return Value{ty: Tclosure, ptr: unsafe.Pointer(c)}
}

func NewBytesValue(buf []byte) Value {
	return Value{ty: Tbytes, ptr: unsafe.Pointer(&buf)}
}

func NewGenericValue(g interface{}) Value {
	return Value{ty: Tgeneric, ptr: unsafe.Pointer(&g)}
}

func (v Value) Type() byte {
	return v.ty
}

func (v Value) AsBool() bool {
	if v.ty != Tbool {
		log.Panicf("not a boolean: %+v", v)
	}
	return uintptr(v.ptr) == trueValue
}

func (v Value) AsBoolUnsafe() bool {
	return uintptr(v.ptr) == trueValue
}

func (v Value) AsNumber() float64 {
	if v.ty != Tnumber {
		log.Panicf("not a number: %+v", v)
	}
	return v.AsNumberUnsafe()
}

func (v Value) AsUint64() uint64 {
	if v.ty != Tnumber {
		log.Panicf("not a number: %+v", v)
	}
	return math.Float64bits(v.AsNumberUnsafe())
}

func (v Value) AsNumberUnsafe() float64 {
	i := uintptr(v.ptr)
	if v.ext {
		i -= minPhysPageSize
	}
	return *(*float64)(unsafe.Pointer(&i))
}

func (v Value) AsString() string {
	if v.ty != Tstring {
		log.Panicf("not a string: %+v", v)
	}
	return *(*string)(v.ptr)
}

func (v Value) AsStringUnsafe() string {
	return *(*string)(v.ptr)
}

func (v Value) AsList() []Value {
	if v.ty != Tlist {
		log.Panicf("not an array: %+v", v)
	}
	return *(*[]Value)(v.ptr)
}

func (v Value) AsListUnsafe() []Value {
	return *(*[]Value)(v.ptr)
}

func (v Value) AsMap() *Tree {
	if v.ty != Tmap {
		log.Panicf("not a map: %+v", v)
	}
	return (*Tree)(v.ptr)
}

func (v Value) AsMapUnsafe() *Tree {
	return (*Tree)(v.ptr)
}

func (v Value) AsClosure() *Closure {
	if v.ty != Tclosure {
		log.Panicf("not a closure: %+v", v)
	}
	return (*Closure)(v.ptr)
}

func (v Value) AsClosureUnsafe() *Closure {
	return (*Closure)(v.ptr)
}

func (v Value) AsGeneric() interface{} {
	if v.ty != Tgeneric {
		log.Panicf("not a generic: %+v", v)
	}
	return *(*interface{})(v.ptr)
}

func (v Value) AsGenericUnsafe() interface{} {
	return *(*interface{})(v.ptr)
}

func (v Value) AsBytes() []byte {
	if v.ty != Tbytes {
		log.Panicf("not a bytes type: %+v", v)
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
		return "<closure:[" + strconv.Itoa(v.AsClosureUnsafe().argsCount) +
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
