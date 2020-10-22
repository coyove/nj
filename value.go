package potatolang

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"unsafe"
)

const (
	NIL    = 0  // nil
	NUM    = 3  // number
	STR    = 7  // string
	STK    = 15 // stack
	FUN    = 31 // function
	ANY    = 63 // generic
	NumNum = NUM * 2
	StrStr = STR * 2
)

// Value is the basic value used by the intepreter
// For float numbers there is one NaN which is not representable: 0xffffffff_ffffffff
// An empty Value naturally represent nil
type Value struct {
	v uint64
	p unsafe.Pointer
}

// Type returns the type of value, its logic should align IsFalse()
func (v Value) Type() byte {
	if v.p == nil || v.p == int64Marker {
		if v.v == 0 {
			return NIL
		}
		return NUM
	}
	if v.v&0xffff_ffff_ffff > 4096 {
		return ANY
	}
	return byte(v.v)
}

// IsFalse tests whether value contains a falsy value: nil or 0
func (v Value) IsFalse() bool {
	x := uintptr(v.p) + uintptr(v.v)
	return x == 0 || x == 0xffffffff_ffffffff
}

var (
	typeMappings = map[byte]string{
		NIL: "nil", NUM: "number", STR: "string", FUN: "function", ANY: "any", STK: "stack",
	}
	int64Marker = unsafe.Pointer(new(int64))
)

func NumBool(v bool) Value {
	if v {
		return Num(1)
	}
	return Num(0)
}

// Num returns a number value
func Num(f float64) Value {
	return Value{v: ^math.Float64bits(f)}
}

// Int also returns a number value as Num does, but it preserves int64 values which may be truncated in float64
func Int(i int64) Value {
	if int64(float64(i)) == i {
		return Value{v: ^math.Float64bits(float64(i))}
	}
	return Value{v: uint64(i), p: int64Marker}
}

// unpackedStack returns a table value
func unpackedStack(m *unpacked) Value {
	if m == nil {
		return Value{}
	}
	return Value{v: STK, p: unsafe.Pointer(m)}
}

// Fun returns a closure value
func Fun(c *Func) Value {
	return Value{v: FUN, p: unsafe.Pointer(c)}
}

// Str returns a string value
func Str(s string) Value {
	if len(s) <= 16 {
		b := [16]byte{}
		copy(b[:], s)
		return Value{v: uint64(len(s)+1)<<56 | STR, p: unsafe.Pointer(&b)}
	}
	return Value{v: STR, p: unsafe.Pointer(&s)}
}

// StrBytes returns a string value
func StrBytes(s []byte) Value {
	return Str(*(*string)(unsafe.Pointer(&s)))
}

func Any(i interface{}) Value {
	switch v := i.(type) {
	case nil:
		return Value{}
	case bool:
		return NumBool(v)
	case float64:
		return Num(v)
	case float32:
		return Num(float64(v))
	case int64:
		return Int(v)
	case string:
		return Str(v)
	case *unpacked:
		return unpackedStack(v)
	case *Func:
		return Fun(v)
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
	// return Value{v: ANY, p: unsafe.Pointer(&i)}
}

// Str cast value to string
func (v Value) Str() string {
	if l := v.v >> 56; l > 0 {
		var ss string
		b := (*[2]uintptr)(unsafe.Pointer(&ss))
		(*b)[0] = uintptr(v.p)
		(*b)[1] = uintptr(l - 1)
		return ss
	}
	return *(*string)(v.p)
}

func (v Value) _StrBytes() []byte {
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

func (v Value) IsNil() bool { return v == Value{} }

func (v Value) Num() (float64, int64, bool) {
	if v.p == int64Marker {
		return float64(int64(v.v)), int64(v.v), true
	}
	x := math.Float64frombits(^v.v)
	return x, int64(x), false
}

func (v Value) Int() int64 {
	_, i, _ := v.Num()
	return i
}

func (v Value) F64() float64 {
	f, _, _ := v.Num()
	return f
}

// unpackedStack cast value to map of values
func (v Value) unpackedStack() *unpacked { return (*unpacked)(v.p) }

// Fun cast value to closure
func (v Value) Fun() *Func { return (*Func)(v.p) }

// Any returns the interface{}
func (v Value) Any() interface{} {
	switch v.Type() {
	case NUM:
		vf, vi, vIsInt := v.Num()
		if vIsInt {
			return vi
		}
		return vf
	case STR:
		return v.Str()
	case STK:
		return v.unpackedStack()
	case FUN:
		return v.Fun()
	case ANY:
		// return *(*interface{})(v.p)
		var i interface{}
		x := (*[2]uintptr)(unsafe.Pointer(&i))
		(*x)[0] = uintptr(v.v)
		(*x)[1] = uintptr(v.p)
		return i
	}
	return nil
}

func (v Value) AnyTyped(t reflect.Type) interface{} {
	if v.Type() == NUM && t.Kind() >= reflect.Int && t.Kind() <= reflect.Float64 {
		rv := reflect.ValueOf(v.Any())
		rv = rv.Convert(t)
		return rv.Interface()
	}
	return v.Any()
}

func (v Value) Expect(t byte) Value {
	if v.Type() != t {
		panicf("expect %s, got %s", typeMappings[t], typeMappings[v.Type()])
	}
	return v
}

func (v Value) ExpectMsg(t byte, msg string) Value {
	if v.Type() != t {
		panicf("%s: expect %s, got %s", msg, typeMappings[t], typeMappings[v.Type()])
	}
	return v
}

// Equal tests whether value is equal to another value
func (v Value) Equal(r Value) bool {
	switch v.Type() + r.Type() {
	case NumNum, NIL * 2:
		return v == r
	case StrStr:
		return r.Str() == v.Str()
	case FUN * 2:
		return v.Fun() == r.Fun()
	case ANY * 2:
		return v.Any() == r.Any()
	}
	return false
}

func (v Value) String() string { return v.toString(0) }

func (v Value) toString(lv int) string {
	if lv > 32 {
		return "<omit deep nesting>"
	}
	switch v.Type() {
	case NUM:
		vf, vi, vIsInt := v.Num()
		if vIsInt {
			return strconv.FormatInt(vi, 10)
		}
		return strconv.FormatFloat(vf, 'f', -1, 64)
	case STR:
		return v.Str()
	case STK:
		return v.unpackedStack().toString(lv + 1)
	case FUN:
		return v.Fun().String()
	case ANY:
		i := v.Any()
		if err := reflectCheckCyclicStruct(i); err != nil {
			return fmt.Sprintf("<any: omit deep nesting>")
		}
		return fmt.Sprintf("%#v", i)
	}
	return "nil"
}
