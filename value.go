package script

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"unsafe"
)

var int64Marker = unsafe.Pointer(new(int64))

// Value is the basic value used by the intepreter
// For float numbers there is one NaN which is not representable: 0xffffffff_ffffffff
// An empty Value naturally represent nil
type Value struct {
	v uint64
	p unsafe.Pointer
}

// Type returns the type of value, its logic should align IsFalse()
func (v Value) Type() valueType {
	if v.p == int64Marker {
		return VNumber
	}
	if v.p == nil {
		if v.v == 0 {
			return VNil
		}
		return VNumber
	}
	if v.v&0xffff_ffff_ffff > 4096 {
		return VInterface
	}
	return valueType(v.v)
}

// IsFalse tests whether value contains a falsy value: nil or 0
func (v Value) IsFalse() bool {
	x := uintptr(v.p) + uintptr(v.v)
	return x == 0 || x == uintptr(int64Marker)
}

func (v Value) IsNil() bool {
	return v == Value{}
}

// Bool returns a NUMBER value: true -> 1 and false -> 0
func Bool(v bool) Value {
	if v {
		return Float(1)
	}
	return Float(0)
}

// Float returns a number value
func Float(f float64) Value {
	if float64(int64(f)) == f {
		return Value{v: uint64(int64(f)), p: int64Marker}
	}
	return Value{v: ^math.Float64bits(f)}
}

// Int returns a number value like Float does, but it preserves int64 values which may overflow float64
func Int(i int64) Value {
	return Value{v: uint64(i), p: int64Marker}
}

func _unpackedStack(m *unpacked) Value {
	if m == nil {
		return Value{}
	}
	return Value{v: VStack, p: unsafe.Pointer(m)}
}

// Function returns a closure value
func Function(c *Func) Value {
	return Value{v: VFunction, p: unsafe.Pointer(c)}
}

func _str(s string) Value {
	return Value{v: VString, p: unsafe.Pointer(&s)}
}

func Interface(i interface{}) Value {
	switch v := i.(type) {
	case nil:
		return Value{}
	case bool:
		return Bool(v)
	case float64:
		return Float(v)
	case float32:
		return Float(float64(v))
	case int64:
		return Int(v)
	case string:
		return _str(v)
	case *unpacked:
		return _unpackedStack(v)
	case *Func:
		return Function(v)
	case Value:
		return v
	}
	rv := reflect.ValueOf(i)
	if k := rv.Kind(); k >= reflect.Int && k <= reflect.Int64 {
		return Int(rv.Int())
	} else if k >= reflect.Uint && k <= reflect.Uintptr {
		return Int(int64(rv.Uint()))
	}

	x := *(*[2]uintptr)(unsafe.Pointer(&i))
	return Value{v: uint64(x[0]), p: unsafe.Pointer(x[1])}
	// return Value{v: VInterface, p: unsafe.Pointer(&i)}
}

// _str cast value to string
func (v Value) _str() string {
	return *(*string)(v.p)
}

func (v Value) _unsafeBytes() []byte {
	var ss []byte
	b := (*[3]uintptr)(unsafe.Pointer(&ss))
	if l := v.v >> 56; l > 0 {
		(*b)[0] = uintptr(v.p)
		(*b)[1], (*b)[2] = uintptr(l-1), uintptr(l-1)
	} else {
		vpp := *(*[2]uintptr)(v.p)
		(*b)[0] = vpp[0]
		(*b)[1], (*b)[2] = vpp[1], vpp[1]
	}
	return ss
}

func (v Value) Num() (floatValue float64, intValue int64, isInt bool) {
	if v.p == int64Marker {
		return float64(int64(v.v)), int64(v.v), true
	}
	x := math.Float64frombits(^v.v)
	return x, int64(x), false
}

func (v Value) Int() int64 { _, i, _ := v.Num(); return i }

func (v Value) Float() float64 { f, _, _ := v.Num(); return f }

func (v Value) _unpackedStack() *unpacked { return (*unpacked)(v.p) }

// Function cast value to function
func (v Value) Function() *Func { return (*Func)(v.p) }

// Interface returns the interface{}
func (v Value) Interface() interface{} {
	switch v.Type() {
	case VNumber:
		vf, vi, vIsInt := v.Num()
		if vIsInt {
			return vi
		}
		return vf
	case VString:
		return v._str()
	case VStack:
		a := v._unpackedStack().a
		x := make([]interface{}, len(a))
		for i := range a {
			x[i] = a[i].Interface()
		}
		return x
	case VFunction:
		return v.Function()
	case VInterface:
		// return *(*interface{})(v.p)
		var i interface{}
		x := (*[2]uintptr)(unsafe.Pointer(&i))
		(*x)[0] = uintptr(v.v)
		(*x)[1] = uintptr(v.p)
		return i
	}
	return nil
}

func (v Value) TypedInterface(t reflect.Type) interface{} {
	switch v.Type() {
	case VNumber:
		if t.Kind() >= reflect.Int && t.Kind() <= reflect.Float64 {
			rv := reflect.ValueOf(v.Interface())
			rv = rv.Convert(t)
			return rv.Interface()
		}
	case VStack:
		a := v._unpackedStack().a
		if t.Kind() == reflect.Slice {
			e := t.Elem()
			s := reflect.MakeSlice(t, len(a), len(a))
			for i := range a {
				s.Index(i).Set(reflect.ValueOf(a[i].TypedInterface(e)))
			}
			return s.Interface()
		}
	}
	return v.Interface()
}

func (v Value) Expect(t valueType) Value {
	if v.Type() != t {
		panicf("expect %s, got %s", typeMappings[t], typeMappings[v.Type()])
	}
	return v
}

func (v Value) ExpectMsg(t valueType, msg string) Value {
	if v.Type() != t {
		panicf("%s: expect %s, got %s", msg, typeMappings[t], typeMappings[v.Type()])
	}
	return v
}

// Equal tests whether value is equal to another value
func (v Value) Equal(r Value) bool {
	switch v.Type() + r.Type() {
	case _NumNum, VNil * 2:
		return v == r
	case _StrStr:
		return r._str() == v._str()
	case VFunction * 2:
		return v.Function() == r.Function()
	case VInterface * 2:
		return v.Interface() == r.Interface()
	}
	return false
}

func (v Value) String() string { return v.toString(0) }

func (v Value) toString(lv int) string {
	if lv > 32 {
		return "<omit deep nesting>"
	}
	switch v.Type() {
	case VNumber:
		vf, vi, vIsInt := v.Num()
		if vIsInt {
			return strconv.FormatInt(vi, 10)
		}
		return strconv.FormatFloat(vf, 'f', -1, 64)
	case VString:
		return v._str()
	case VStack:
		t := v._unpackedStack()
		p := bytes.NewBufferString("[")
		for _, a := range t.a {
			p.WriteString(a.toString(lv + 1))
			p.WriteString(",")
		}
		if len(t.a) > 0 {
			p.Truncate(p.Len() - 1)
		}
		p.WriteString("]")
		return p.String()
	case VFunction:
		return v.Function().String()
	case VInterface:
		i := v.Interface()
		if err := reflectCheckCyclicStruct(i); err != nil {
			return fmt.Sprintf("<any: omit deep nesting>")
		}
		return fmt.Sprintf("%v", i)
	}
	return "nil"
}

type unpacked struct{ a []Value }

func (t *unpacked) Put(idx int64, v Value) {
	idx--
	if idx < int64(len(t.a)) && idx >= 0 {
		t.a[idx] = v
	}
}

func (t *unpacked) Get(idx int64) (v Value) {
	idx--
	if idx < int64(len(t.a)) && idx >= 0 {
		return t.a[idx]
	}
	return Value{}
}
