package script

import (
	"bytes"
	"encoding/binary"
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
	int64Marker       = unsafe.Pointer(new(int64))
	trueMarker        = unsafe.Pointer(new(int64))
	falseMarker       = unsafe.Pointer(new(int64))
	smallStringMarker = unsafe.Pointer(new(int64))
	falseValue        = Bool(false)
	zeroValue         = Int(0)
	watermark         = Interface(new(int))

	Nil = Value{}
)

type ValueType byte

const (
	VNil       ValueType = 0
	VBool      ValueType = 1
	VNumber    ValueType = 3
	VString    ValueType = 7
	VMap       ValueType = 15
	VFunction  ValueType = 17
	VInterface ValueType = 19

	ValueSize = int64(unsafe.Sizeof(Value{}))
)

func (t ValueType) String() string {
	if t > VInterface {
		return "?"
	}
	return [...]string{"nil", "bool", "?", "number", "?", "?", "?", "string", "?", "?", "?", "?", "?", "?", "?", "array", "?", "function", "?", "native"}[t]
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
	case smallStringMarker:
		return VString
	case nil:
		if v.v == 0 {
			return VNil
		}
		return VNumber
	}
	return ValueType(v.v)
}

// IsFalse tests whether value contains a falsy value: nil, false or 0
func (v Value) IsFalse() bool { return v == Nil || v == zeroValue || v == falseValue }

func (v Value) IsInt() bool { return v.p == int64Marker }

// Bool returns a boolean value
func Bool(v bool) Value {
	if v {
		return Value{uint64(VBool), trueMarker}
	}
	return Value{uint64(VBool), falseMarker}
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
	return (&Map{items: m, count: uint32(len(m))}).Value()
}

func ArrayMap(kvs ...Value) Value {
	t := &Map{}
	for i := 0; i < len(kvs)/2*2; i += 2 {
		t.Set(kvs[i], kvs[i+1])
	}
	return Value{v: uint64(VMap), p: unsafe.Pointer(t)}
}

// Function returns a closure value
func Function(c *Func) Value {
	if c.Name == "" {
		c.Name = "function"
	}
	return Value{v: uint64(VFunction), p: unsafe.Pointer(c)}
}

func String(s string) Value {
	if len(s) <= 7 { // length 1b + payload 7b
		x := [8]byte{byte(len(s))}
		copy(x[1:], s)
		return Value{v: binary.BigEndian.Uint64(x[:]), p: smallStringMarker}
	}
	return Value{v: uint64(VString), p: unsafe.Pointer(&s)}
}

func Bytes(b []byte) Value {
	return String(*(*string)(unsafe.Pointer(&b)))
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
	case *Map:
		return v.Value()
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
				case 0:
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
	return Value{v: uint64(VInterface), p: unsafe.Pointer(&i)}
}

// rawStr cast value to string
func (v Value) rawStr() string {
	if v.p == smallStringMarker {
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, v.v)
		buf = buf[1 : 1+buf[0]]
		return *(*string)(unsafe.Pointer(&buf))
	}
	return *(*string)(v.p)
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

func (v Value) Map() *Map { return (*Map)(v.p) }

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
		return v.rawStr()
	case VMap:
		return v.Map()
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
	case VMap:
		a := v.Map()
		if t.Kind() == reflect.Slice {
			e := t.Elem()
			s := reflect.MakeSlice(t, len(a.Array()), len(a.Array()))
			for i, a := range a.Array() {
				s.Index(i).Set(reflect.ValueOf(a.TypedInterface(e)))
			}
			return s.Interface()
		}
	}
	return v.Interface()
}

func (v Value) MustBool(msg string, a int) bool { return v.mustBe(VBool, msg, a).Bool() }

func (v Value) MustString(msg string, a int) string { return v.mustBe(VString, msg, a).String() }

func (v Value) MustNumber(msg string, a int) Value { return v.mustBe(VNumber, msg, a) }

func (v Value) MustMap(msg string, a int) *Map { return v.mustBe(VMap, msg, a).Map() }

func (v Value) MustFunc(msg string, a int) *Func { return v.mustBe(VFunction, msg, a).Function() }

func (v Value) mustBe(t ValueType, msg string, msgArg int) Value {
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
		return r.rawStr() == v.rawStr()
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
		return v.rawStr() < r.rawStr()
	case VString + VNumber:
		if v.Type() == VNumber {
			return true
		}
		return false
	}
	return false
}

func (v Value) HashCode() uint64 {
	code := uint64(5381)
	if v.Type() != VString || v.p == smallStringMarker {
		for _, r := range *(*[ValueSize]byte)(unsafe.Pointer(&v)) {
			code = code*33 + uint64(r)
		}
	} else {
		for _, r := range v.rawStr() {
			code = code*33 + uint64(r)
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
			return strconv.Quote(v.rawStr())
		}
		return v.rawStr()
	case VMap:
		m := v.Map()
		if len(m.hashItems) == 0 {
			p := bytes.NewBufferString("[")
			for _, a := range m.Array() {
				p.WriteString(a.toString(lv+1, j))
				p.WriteString(",")
			}
			return strings.TrimRight(p.String(), ", ") + "]"
		}
		p := bytes.NewBufferString("{")
		for k, v := m.Next(Nil); k != Nil; k, v = m.Next(k) {
			p.WriteString(k.toString(lv+1, j))
			p.WriteString(":")
			p.WriteString(v.toString(lv+1, j))
			p.WriteString(",")
		}
		return strings.TrimRight(p.String(), ", ") + "}"
	case VFunction:
		return v.Function().String()
	case VInterface:
		i := v.Interface()
		if !reflectCheckCyclicStruct(i) {
			i = "<interface: omit deep nesting>"
		}
		if j {
			buf, _ := json.Marshal(i)
			return string(buf)
		}
		return fmt.Sprintf("%v", i)
	}
	return "nil"
}

func (v Value) StringDefault(d string) string {
	if v.Type() == VString {
		return v.rawStr()
	}
	return d
}

func (v Value) IntDefault(d int64) int64 {
	if v.Type() == VNumber {
		_, i, _ := v.Num()
		return i
	}
	return d
}

type Values struct {
	Underlay []Value
}

func (t *Values) Len() int {
	return len(t.Underlay)
}

func (t *Values) Slice(start int64, end int64) []Value {
	start, end = sliceInRange(start, end, int64(len(t.Underlay)))
	return t.Underlay[start:end]
}

func (t *Values) Put(idx int64, v Value) (appended bool) {
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

func (t *Values) Get(idx int64) (v Value) {
	if idx < int64(len(t.Underlay)) && idx >= 0 {
		return t.Underlay[idx]
	}
	return Value{}
}
