package base

import (
	"bytes"
	"encoding/binary"
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

	i byte
	p [4]byte

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

func (v *Value) Attachments() int {
	if v.i<<6 > 0 {
		return 4
	}
	if v.i<<4 > 0 {
		return 3
	}
	if v.i<<2 > 0 {
		return 2
	}
	if v.i > 0 {
		return 1
	}
	return 0
}

func (v *Value) Attach(va Value) {
	switch va.ty {
	case Tbool:
		bu := va.AsBoolUnsafe()
		b := *(*byte)(unsafe.Pointer(&bu))
		if v.i == 0 {
			v.p[0] = b
			v.i |= 0x40
		} else if v.i<<2 == 0 {
			v.p[1] = b
			v.i |= 0x10
		} else if v.i<<4 == 0 {
			v.p[2] = b
			v.i |= 0x04
		} else if v.i<<6 == 0 {
			v.p[3] = b
			v.i |= 0x01
		}
	case Tnumber:
		n := va.AsNumberUnsafe()
		if u := uint64(n); float64(u) == n && u < 256 {
			if v.i == 0 {
				v.p[0] = byte(u)
				v.i |= 0x80
			} else if v.i<<2 == 0 {
				v.p[1] = byte(u)
				v.i |= 0x20
			} else if v.i<<4 == 0 {
				v.p[2] = byte(u)
				v.i |= 0x08
			} else if v.i<<6 == 0 {
				v.p[3] = byte(u)
				v.i |= 0x02
			}
		} else {
			if v.i == 0 {
				binary.LittleEndian.PutUint32(v.p[:], math.Float32bits(float32(n)))
				v.i = 0xFF
			}
		}
	}
}

func (v *Value) Detach() Value {
	if v.i == 0xFF {
		v.i = 0
		return NewNumberValue(float64(math.Float32frombits(binary.LittleEndian.Uint32(v.p[:]))))
	}

	ret := func(i, m byte) Value {
		if m == 0x40 {
			return NewBoolValue(*(*bool)(unsafe.Pointer(&v.p[i])))
		}
		return NewNumberValue(float64(v.p[i]))
	}

	if m := v.i << 6; m > 0 {
		v.i &= 0xFC
		return ret(3, m)
	}

	if m := v.i >> 2 << 6; m > 0 {
		v.i &= 0xF3
		return ret(2, m)
	}

	if m := v.i >> 4 << 6; m > 0 {
		v.i &= 0xCF
		return ret(1, m)
	}

	if m := v.i >> 6 << 6; m > 0 {
		v.i &= 0x3F
		return ret(0, m)
	}

	return NewValue()
}
