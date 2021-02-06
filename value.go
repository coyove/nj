package script

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/coyove/script/parser"
)

var (
	int64Marker   = unsafe.Pointer(new(int64))
	trueMarker    = unsafe.Pointer(new(int64))
	falseMarker   = unsafe.Pointer(new(int64))
	sStringMarker = unsafe.Pointer(new(int64))
	falseValue    = Bool(false)
	zeroValue     = Int(0)
)

type ValueType byte

const (
	VNil       ValueType = 0
	VBool                = 1
	VNumber              = 3
	VString              = 7
	VArray               = 15
	VFunction            = 17
	VInterface           = 19

	ValueSize = int64(unsafe.Sizeof(Value{}))
)

func (t ValueType) String() string {
	if t > VInterface {
		return "?"
	}
	return [...]string{"nil", "bool", "?", "number", "?", "?", "?", "string", "?", "?", "?", "?", "?", "?", "?", "array", "?", "function", "?", "interface"}[t]
}

// Value is the basic data type used by the intepreter
// For float numbers there is one NaN which is not representable: 0xffffffff_ffffffff
// An empty Value naturally represent nil
type Value struct {
	v uint64
	p unsafe.Pointer
}

// Reverse-reference in 'parser' package
func (v Value) IsValue(parser.Node) {}

// Type returns the type of value, its logic should align IsFalse()
func (v Value) Type() ValueType {
	switch v.p {
	case int64Marker:
		return VNumber
	case trueMarker, falseMarker:
		return VBool
	case nil:
		if v.v == 0 {
			return VNil
		}
		return VNumber
	}
	return ValueType(v.v)
}

// IsFalse tests whether value contains a falsy value: nil or 0
func (v Value) IsFalse() bool {
	return v == Value{} || v == zeroValue || v == falseValue
}

func (v Value) IsNil() bool {
	return v == Value{}
}

// Bool returns a boolean value
func Bool(v bool) Value {
	if v {
		return Value{VBool, trueMarker}
	}
	return Value{VBool, falseMarker}
}

// Float returns a number value
func Float(f float64) Value {
	if float64(int64(f)) == f {
		return Value{v: uint64(int64(f)), p: int64Marker}
	}
	return Value{v: ^math.Float64bits(f)}
}

// Int returns a number value like Float does, but it preserves int64 values which overflow float64
func Int(i int64) Value {
	return Value{v: uint64(i), p: int64Marker}
}

// Array returns an array consists of 'm'
func Array(m ...Value) Value {
	return Value{v: VArray, p: unsafe.Pointer(&Values{Underlay: m})}
}

// Function returns a closure value
func Function(c *Func) Value {
	if c.Name == "" {
		c.Name = "function"
	}
	return Value{v: VFunction, p: unsafe.Pointer(c)}
}

func String(s string) Value {
	if len(s) <= 6 { // length 1b + payload 6b + VString type 1b = uint64
		x := [8]byte{}
		copy(x[1:], s)
		u64 := uint64(len(s))<<56 | *(*uint64)(unsafe.Pointer(&x)) | uint64(VString)
		return Value{v: u64, p: sStringMarker}
	}
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
		return Array(v...)
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
					env.A = Interface(outs[0].Interface())
				default:
					a := make([]Value, len(outs))
					for i := range outs {
						a[i] = Interface(outs[i].Interface())
					}
					env.A = Array(a...)
				}
			}
		}
		return Function(&Func{Name: "<native>", Native: nf})
	}
	return _interface(i)
}

func _interface(i interface{}) Value {
	return Value{v: VInterface, p: unsafe.Pointer(&i)}
}

// _str cast value to string
func (v Value) _str() string {
	if v.p == sStringMarker {
		tmp := *(*[8]byte)(unsafe.Pointer(&v.v))
		return string(tmp[1 : 1+v.v>>56])
	}
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

func (v Value) Bool() bool { return v.p == trueMarker }

func (v Value) Array() *Values { return (*Values)(v.p) }

// Function cast value to function
func (v Value) Function() *Func { return (*Func)(v.p) }

// Interface returns the interface{}
func (v Value) Interface() interface{} {
	switch v.Type() {
	case VBool:
		return v.Bool()
	case VNumber:
		vf, vi, vIsInt := v.Num()
		if vIsInt {
			return vi
		}
		return vf
	case VString:
		return v._str()
	case VArray:
		a := v.Array().Underlay
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
	case VArray:
		a := v.Array().Underlay
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

func (v Value) MustBe(t ValueType, msg string, msgArg int) Value {
	if v.Type() != t {
		if msgArg > 0 {
			panicf("%s %d: expect %v, got %v", msg, msgArg, t, v.Type())
		}
		panicf("%s: expect %v, got %v", msg, t, v.Type())
	}
	return v
}

func (v Value) Equal(r Value) bool {
	if v == r {
		return true
	}
	switch v.Type() + r.Type() {
	case VString * 2:
		return r._str() == v._str()
	case VInterface * 2:
		return v.Interface() == r.Interface()
	}
	return false
}

func (v Value) Less(r Value) bool {
	switch v.Type() + r.Type() {
	case VNumber * 2:
		vf, vi, vIsInt := v.Num()
		rf, ri, rIsInt := r.Num()
		if vIsInt && rIsInt {
			return vi < ri
		}
		return vf < rf
	case VString * 2:
		return r._str() < v._str()
	case VString + VNumber:
		if v.Type() == VNumber {
			return true
		}
		return false
	}
	return false
}

func (v Value) HashCode() [2]uint64 {
	if v.Type() != VString {
		return *(*[2]uint64)(unsafe.Pointer(&v))
	}
	code := [2]uint64{1 << 63, 5381}
	for _, r := range v._str() {
		old := code[1]
		code[1] = code[1]*33 + uint64(r)
		if code[1] < old {
			code[0]++
		}
	}
	return code
}

func (v Value) String() string { return v.toString(0, false) }

func (v Value) JSONString() string { return v.toString(0, true) }

func (v Value) MarshalJSON() ([]byte, error) { return []byte(v.toString(0, true)), nil }

func (v Value) toString(lv int, j bool) string {
	if lv > 32 {
		return "<omit deep nesting>"
	}
	switch v.Type() {
	case VBool:
		return strconv.FormatBool(v.Bool())
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
		p := bytes.NewBufferString("[")
		for _, a := range v.Array().Underlay {
			p.WriteString(a.toString(lv+1, j))
			p.WriteString(",")
		}
		return strings.TrimRight(p.String(), ", ") + "]"
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
	Unpacked bool
	Underlay []Value
}

func (t *Values) Len() int {
	return len(t.Underlay)
}

func (t *Values) Slice1(start int64, end int64) []Value {
	start2, end2 := sliceInRange(start, end, len(t.Underlay))
	return t.Underlay[start2:end2]
}

func (t *Values) Put1(idx int64, v Value) (appended bool) {
	idx--
	if idx < int64(len(t.Underlay)) && idx >= 0 {
		t.Underlay[idx] = v
		return false
	}
	if idx == int64(len(t.Underlay)) {
		t.Underlay = append(t.Underlay, v)
		return true
	}
	return false
}

func (t *Values) Get1(idx int64) (v Value) {
	idx--
	if idx < int64(len(t.Underlay)) && idx >= 0 {
		return t.Underlay[idx]
	}
	return Value{}
}
