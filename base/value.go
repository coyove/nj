package base

import (
	"bytes"
	"log"
	"unsafe"
)

const (
	TY_nil = iota
	TY_number
	TY_string
	TY_bool
	TY_array
	TY_bytes
	TY_map
	TY_closure
	TY_generic
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

	ptr unsafe.Pointer
}

func NewValue() Value {
	return Value{ty: TY_nil}
}

func NewNumberValue(f float64) Value {
	i := *(*uintptr)(unsafe.Pointer(&f))
	if i < minPhysPageSize {
		i += minPhysPageSize
		return Value{ty: TY_number, ext: true, ptr: unsafe.Pointer(i)}
	}
	return Value{ty: TY_number, ptr: unsafe.Pointer(i)}
}

func NewStringValue(s string) Value {
	return Value{ty: TY_string, ptr: unsafe.Pointer(&s)}
}

func NewBoolValue(b bool) Value {
	v := Value{ty: TY_bool}
	if b {
		v.ptr = unsafe.Pointer(uintptr(trueValue))
	} else {
		v.ptr = unsafe.Pointer(uintptr(falseValue))
	}
	return v
}

func NewArrayValue(a []Value) Value {
	return Value{ty: TY_array, ptr: unsafe.Pointer(&a)}
}

func NewMapValue(m map[string]Value) Value {
	return Value{ty: TY_map, ptr: unsafe.Pointer(&m)}
}

func NewClosureValue(c *Closure) Value {
	return Value{ty: TY_closure, ptr: unsafe.Pointer(c)}
}

func NewBytesValue(buf []byte) Value {
	if buf == nil {
		return Value{ty: TY_bytes, ptr: nil}
	}
	return Value{ty: TY_bytes, ptr: unsafe.Pointer(&buf)}
}

func NewGenericValue(g interface{}) Value {
	return Value{ty: TY_generic, ptr: unsafe.Pointer(&g)}
}

func (v Value) Type() byte {
	return v.ty
}

func (v Value) Bool() bool {
	if v.ty != TY_bool {
		log.Panicf("not a boolean: %d", v.ty)
	}
	return uintptr(v.ptr) == trueValue
}

func (v Value) Number() float64 {
	if v.ty != TY_number {
		log.Panicf("not a number: %d", v.ty)
	}
	i := uintptr(v.ptr)
	if v.ext {
		i -= minPhysPageSize
	}

	return *(*float64)(unsafe.Pointer(&i))
}

func (v Value) String() string {
	if v.ty != TY_string {
		log.Panicf("not a string: %d", v.ty)
	}
	return *(*string)(v.ptr)
}

func (v Value) Array() []Value {
	if v.ty != TY_array {
		log.Panicf("not an array: %d", v.ty)
	}
	return *(*[]Value)(v.ptr)
}

func (v Value) Map() map[string]Value {
	if v.ty != TY_map {
		log.Panicf("not a map: %d", v.ty)
	}
	return *(*map[string]Value)(v.ptr)
}

func (v Value) Closure() *Closure {
	if v.ty != TY_closure {
		log.Panicf("not a closure: %d", v.ty)
	}
	return (*Closure)(v.ptr)
}

func (v Value) Generic() interface{} {
	if v.ty != TY_generic {
		log.Panicf("not a generic: %d", v.ty)
	}
	return *(*interface{})(v.ptr)
}

func (v Value) I() interface{} {
	switch v.Type() {
	case TY_bool:
		return v.Bool()
	case TY_number:
		return v.Number()
	case TY_string:
		return v.String()
	case TY_array:
		return v.Array()
	case TY_map:
		return v.Map()
	case TY_bytes:
		return v.Bytes()
	case TY_closure:
		return v.Closure()
	}
	return nil
}

func (v Value) Bytes() []byte {
	if v.ty != TY_bytes {
		log.Panicf("not a bytes type: %d", v.ty)
	}
	if v.ptr == nil {
		return nil
	}
	return *(*[]byte)(v.ptr)
}

func (v Value) IsFalse() bool {
	if v.ty == TY_nil {
		return true
	}
	switch v.Type() {
	case TY_bool:
		return v.Bool() == false
	case TY_number:
		return v.Number() == 0.0
	case TY_string:
		return v.String() == ""
	case TY_array:
		return len(v.Array()) == 0
	}
	return false
}

func (v Value) Equal(r Value) bool {
	if v.ty == TY_nil || r.ty == TY_nil {
		return v.ty == r.ty
	}

	switch v.ty {
	case TY_number:
		if r.ty == TY_number {
			return r.Number() == v.Number()
		}
	case TY_string:
		if r.ty == TY_string {
			return r.String() == v.String()
		} else if r.ty == TY_bytes {
			return bytes.Equal(r.Bytes(), []byte(v.String()))
		}
	case TY_bool:
		if r.ty == TY_bool {
			return r.Bool() == v.Bool()
		}
	case TY_array:
		if r.ty == TY_array {
			lf, rf := v.Array(), r.Array()

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
	case TY_bytes:
		if r.ty == TY_bytes {
			return bytes.Equal(v.Bytes(), r.Bytes())
		} else if r.ty == TY_string {
			return bytes.Equal(v.Bytes(), []byte(r.String()))
		}
	}

	return false
}

func (v Value) Less(r Value) bool {
	switch v.ty {
	case TY_number:
		if r.ty == TY_number {
			return v.Number() < r.Number()
		}
	case TY_string:
		if r.ty == TY_string {
			return v.String() < r.String()
		}
	}
	log.Panicf("can't compare ty:%d and ty:%d", v.ty, r.ty)
	return false
}

func (v Value) LessEqual(r Value) bool {
	switch v.ty {
	case TY_number:
		if r.ty == TY_number {
			return v.Number() <= r.Number()
		}
	case TY_string:
		if r.ty == TY_string {
			return v.String() <= r.String()
		}
	}
	log.Panicf("can't compare ty:%d and ty:%d", v.ty, r.ty)
	return false
}
