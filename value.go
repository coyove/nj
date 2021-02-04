package script

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/coyove/script/parser"
)

var int64Marker = unsafe.Pointer(new(int64))

// Value is the basic value used by the intepreter
// For float numbers there is one NaN which is not representable: 0xffffffff_ffffffff
// An empty Value naturally represent nil
type Value struct {
	v uint64
	p unsafe.Pointer
}

// Reverse-reference in 'parser' package
func (v Value) IsValue(parser.Node) {}

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
	// if v.v&0xffff_ffff_ffff > 4096 {
	// 	return VInterface
	// }
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

// Array returns an array consists of 'm'
func Array(m []Value) Value {
	return Value{v: VArray, p: unsafe.Pointer(&Values{a: m})}
}

// Function returns a closure value
func Function(c *Func) Value {
	if c.Name == "" {
		c.Name = "function"
	}
	return Value{v: VFunction, p: unsafe.Pointer(c)}
}

func String(s string) Value {
	return Value{v: VString, p: unsafe.Pointer(&s)}
}

func Bytes(b []byte) Value {
	return Value{v: VString, p: unsafe.Pointer(&b)}
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
		return String(v)
	case []byte:
		return Bytes(v)
	case []Value:
		return Array(v)
	case *Func:
		return Function(v)
	case Value:
		return v
	case parser.CatchedError:
		return Interface(v.Original)
	case reflect.Value:
		return Interface(v.Interface())
	}

	rv := reflect.ValueOf(i)
	if k := rv.Kind(); k >= reflect.Int && k <= reflect.Int64 {
		return Int(rv.Int())
	} else if k >= reflect.Uint && k <= reflect.Uintptr {
		return Int(int64(rv.Uint()))
	} else if k == reflect.Func {
		nf, _ := i.(func(*Env))
		if nf == nil {
			rt := rv.Type()
			rtNumIn := rt.NumIn()
			nf = func(env *Env) {
				getter := func(i int, t reflect.Type) reflect.Value {
					return reflect.ValueOf(env.Get(i).TypedInterface(t))
				}
				ins := make([]reflect.Value, 0, rtNumIn)
				if !rt.IsVariadic() {
					for i := 0; i < rtNumIn; i++ {
						ins = append(ins, getter(i, rt.In(i)))
					}
				} else {
					for i := 0; i < rtNumIn-1; i++ {
						ins = append(ins, getter(i, rt.In(i)))
					}
					for i := rtNumIn - 1; i < env.Size(); i++ {
						ins = append(ins, getter(i, rt.In(rtNumIn-1).Elem()))
					}
				}
				switch outs := rv.Call(ins); rt.NumOut() {
				case 1:
					env.A, env.V = Interface2(outs[0].Interface())
				default:
					a := make([]Value, len(outs))
					for i := range outs {
						a[i] = Interface(outs[i].Interface())
					}
					env.Return(a...)
				}
			}
		}
		return Function(&Func{Name: "<native>", Native: nf})
	}
	return _interface(i)
}

func Interface2(i interface{}) (Value, []Value) {
	if s, ok := i.([]Value); ok {
		return Int(int64(len(s))), s
	}
	return Interface(i), nil
}

func _interface(i interface{}) Value {
	return Value{v: VInterface, p: unsafe.Pointer(&i)}
}

// _str cast value to string
func (v Value) _str() string {
	return *(*string)(v.p)
}

func (v Value) _unsafeBytes() []byte {
	var ss []byte
	b := (*[3]uintptr)(unsafe.Pointer(&ss))
	vpp := *(*[2]uintptr)(v.p)
	(*b)[0] = vpp[0]
	(*b)[1], (*b)[2] = vpp[1], vpp[1]
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

func (v Value) Array() *Values { return (*Values)(v.p) }

func (v Value) Stack() []Value { return v.Array().a }

func (v Value) _append(v2 Value) Value {
	if v.Type() == VArray {
		s := v.Array()
		s.a = append(s.a, v2)
		return Array(s.a)
	}
	return Array([]Value{v, v2})
}

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
	case VArray:
		a := v.Array().a
		x := make([]interface{}, len(a))
		for i := range a {
			x[i] = a[i].Interface()
		}
		return x
	case VFunction:
		return v.Function()
	case VInterface:
		return *(*interface{})(v.p)
	}
	return nil
}

func (v Value) TypedInterface(t reflect.Type) interface{} {
	if t == reflect.TypeOf(Value{}) {
		return v
	}

	switch v.Type() {
	case VNumber:
		if t.Kind() >= reflect.Int && t.Kind() <= reflect.Float64 {
			rv := reflect.ValueOf(v.Interface())
			rv = rv.Convert(t)
			return rv.Interface()
		}
		if t.Kind() == reflect.Bool {
			return !v.IsFalse()
		}
	case VArray:
		a := v.Array().a
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
		panicf("expect %v, got %v", t, v.Type())
	}
	return v
}

func (v Value) ExpectMsg(t valueType, msg string) Value {
	if v.Type() != t {
		panicf("%s: expect %v, got %v", msg, t, v.Type())
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

func (v Value) Less(r Value) bool {
	switch v.Type() + r.Type() {
	case _NumNum:
		vf, vi, vIsInt := v.Num()
		rf, ri, rIsInt := r.Num()
		if vIsInt && rIsInt {
			return vi < ri
		}
		return vf < rf
	case _StrStr:
		return r._str() < v._str()
	case VString + VNumber:
		if v.Type() == VNumber {
			return true
		}
		return false
	}
	return false
}

func (v Value) String() string { return v.toString(0, false) }

func (v Value) JSONString() string { return v.toString(0, true) }

func (v Value) MarshalJSON() ([]byte, error) { return []byte(v.toString(0, true)), nil }

func (v Value) toString(lv int, j bool) string {
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
		if j {
			return strconv.Quote(v._str())
		}
		return v._str()
	case VArray:
		t := v.Array()
		p := bytes.NewBufferString("[")
		for _, a := range t.a {
			p.WriteString(a.toString(lv+1, j))
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
		if !reflectCheckCyclicStruct(i) {
			i = fmt.Sprintf("<interface: omit deep nesting>")
		}
		if j {
			buf, _ := json.Marshal(i)
			return string(buf)
		}
		return fmt.Sprintf("%v", i)
	}
	return "nil"
}

type Values struct {
	a []Value
}

func (t *Values) Slice(start int64, end int64) []Value {
	start2, end2 := sliceInRange(start, end, len(t.a))
	return t.a[start2:end2]
}

func (t *Values) Put(idx int64, v Value) {
	idx--
	if idx < int64(len(t.a)) && idx >= 0 {
		t.a[idx] = v
	}
}

func (t *Values) Get(idx int64) (v Value) {
	idx--
	if idx < int64(len(t.a)) && idx >= 0 {
		return t.a[idx]
	}
	return Value{}
}
