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
	"unicode/utf8"
	"unsafe"

	"github.com/coyove/script/parser"
	"github.com/coyove/script/typ"
	"github.com/tidwall/gjson"
)

var (
	int64Marker       = unsafe.Pointer(new(int64))
	int64Marker2      = uintptr(int64Marker) * 2
	float64Marker     = unsafe.Pointer(new(int64))
	float64Marker2    = uintptr(float64Marker) * 2
	trueMarker        = unsafe.Pointer(new(int64))
	falseMarker       = unsafe.Pointer(new(int64))
	smallStringMarker = unsafe.Pointer(new([9]int64))

	Nil     = Value{}
	Undef   = Val(new(int))
	Zero    = Int(0)
	NullStr = Str("")
	False   = Bool(false)
	True    = Bool(true)
)

const (
	ValueSize = int64(unsafe.Sizeof(Value{}))

	errNeedNumbers          = "operator requires numbers"
	errNeedNumbersOrStrings = "operator requires number-to-number or string-to-string"
)

// Value is the basic data type used by the intepreter, an empty Value naturally represent nil
type Value struct {
	v uint64
	p unsafe.Pointer
}

// Reverse-reference in 'parser' package
func (v Value) IsValue(parser.Node) {}

// Type returns the type of value, its logic should align IsFalse()
func (v Value) Type() typ.ValueType {
	switch v.p {
	case int64Marker, float64Marker:
		return typ.Number
	case trueMarker, falseMarker:
		return typ.Bool
	case nil:
		return typ.Nil
	}
	if uintptr(v.p) >= uintptr(smallStringMarker) && uintptr(v.p) <= uintptr(smallStringMarker)+8*8 {
		return typ.String
	}
	return typ.ValueType(v.v)
}

// IsFalse tests whether value contains a falsy value: nil, false, empty string or 0
// Note that empty list and nil pointer in golang are considered as 'true'
func (v Value) IsFalse() bool {
	return v == Nil || v == Zero || v == False || v == NullStr
}

// IsInt tests whether value contains an integer (int64)
func (v Value) IsInt() bool { return v.p == int64Marker }

// Bool returns a boolean value
func Bool(v bool) Value {
	if v {
		return Value{uint64(typ.Bool), trueMarker}
	}
	return Value{uint64(typ.Bool), falseMarker}
}

// Float returns a number value
func Float(f float64) Value {
	if float64(int64(f)) == f {
		return Value{v: uint64(int64(f)), p: int64Marker}
	}
	return Value{v: math.Float64bits(f), p: float64Marker}
}

// Int returns a number value like Float does, but it preserves int64 values which overflow float64
func Int(i int64) Value {
	return Value{v: uint64(i), p: int64Marker}
}

// Array returns an array consists of 'm'
func Array(m ...Value) Value {
	x := &HashMap{items: m}
	for _, i := range x.items {
		if i != Nil {
			x.count++
		}
	}
	return x.Value()
}

// Map returns a map, kvs should be laid out as: key1, value1, key2, value2, ...
func Map(kvs ...Value) Value {
	t := NewHashMap(len(kvs) / 2)
	for i := 0; i < len(kvs)/2*2; i += 2 {
		t.Set(kvs[i], kvs[i+1])
	}
	return Value{v: uint64(typ.Map), p: unsafe.Pointer(t)}
}

// MapVal returns a map, kvs should be laid out as: key1, value1, key2, value2, ...
func MapVal(kvs ...interface{}) Value {
	t := NewHashMap(len(kvs) / 2)
	for i := 0; i < len(kvs)/2*2; i += 2 {
		t.Set(Val(kvs[i]), Val(kvs[i+1]))
	}
	return Value{v: uint64(typ.Map), p: unsafe.Pointer(t)}
}

// MapWithParent returns a map whose parent will be set to p
func MapWithParent(p *HashMap, kvs ...Value) Value {
	m := Map(kvs...)
	m.Map().Parent = p
	return m
}

// Str returns a string value
func Str(s string) Value {
	if len(s) <= 8 { // payload 7b
		x := [8]byte{byte(len(s))}
		copy(x[:], s)
		return Value{
			v: binary.BigEndian.Uint64(x[:]),
			p: unsafe.Pointer(uintptr(smallStringMarker) + uintptr(len(s))*8),
		}
	}
	return Value{v: uint64(typ.String), p: unsafe.Pointer(&s)}
}

// Byte returns a one-byte string value
func Byte(s byte) Value {
	x := [8]byte{s}
	return Value{v: binary.BigEndian.Uint64(x[:]), p: unsafe.Pointer(uintptr(smallStringMarker) + 8)}
}

// Rune returns a one-rune string value
func Rune(r rune) Value {
	x := [8]byte{}
	n := utf8.EncodeRune(x[:], r)
	return Value{v: binary.BigEndian.Uint64(x[:]), p: unsafe.Pointer(uintptr(smallStringMarker) + uintptr(n)*8)}
}

// Bytes returns an alterable string value
func Bytes(b []byte) Value {
	return Value{v: uint64(typ.String), p: unsafe.Pointer(&b)}
	// return Str(*(*string)(unsafe.Pointer(&b)))
}

// Val creates a Value from golang interface{}
// []Type, [..]Type and map[Type]Type will be left as is (except []byte and []Value),
// to convert them recursively, use ValRec instead
func Val(i interface{}) Value {
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
		return Str(v)
	case []byte:
		return Bytes(v)
	case *HashMap:
		return v.Value()
	case []Value:
		return Array(v...)
	case *Func:
		return v.Value()
	case Value:
		return v
	case parser.CatchedError:
		return Val(v.Original)
	case reflect.Value:
		return Val(v.Interface())
	case gjson.Result:
		switch v.Type {
		case gjson.String:
			return Str(v.Str)
		case gjson.Number:
			return Float(v.Float())
		case gjson.True, gjson.False:
			return Bool(v.Bool())
		}
		if v.IsArray() {
			a := v.Array()
			x := make([]Value, len(a))
			for i, a := range a {
				x[i] = Val(a)
			}
			return Array(x...)
		}
		if v.IsObject() {
			m := v.Map()
			x := NewHashMap(len(m))
			for k, v := range m {
				x.Set(Str(k), Val(v))
			}
			return x.Value()
		}
		return Nil
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
					return env.Get(i).ReflectValue(t)
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
					env.A = Val(outs[0].Interface())
				default:
					a := make([]Value, len(outs))
					for i := range outs {
						a[i] = Val(outs[i].Interface())
					}
					env.A = Array(a...)
				}
			}
		}
		return (&Func{Name: "<native>", Native: nf}).Value()
	}
	return _interface(i)
}

func ValRec(v interface{}) Value {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map:
		m := NewHashMap(rv.Len() + 1)
		iter := rv.MapRange()
		for iter.Next() {
			m.Set(ValRec(iter.Key()), Val(iter.Value()))
		}
		return m.Value()
	case reflect.Array, reflect.Slice:
		a := make([]Value, rv.Len())
		for i := range a {
			a[i] = ValRec(rv.Index(i))
		}
		return Array(a...)
	default:
		return Val(v)
	}
}

func _interface(i interface{}) Value {
	return Value{v: uint64(typ.Interface), p: unsafe.Pointer(&i)}
}

func (v Value) IsSmallString() bool {
	return uintptr(v.p) >= uintptr(smallStringMarker) && uintptr(v.p) <= uintptr(smallStringMarker)+8*8
}

func (v Value) Str() string {
	if v.IsSmallString() {
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, v.v)
		buf = buf[:(uintptr(v.p)-uintptr(smallStringMarker))/8]
		return *(*string)(unsafe.Pointer(&buf))
	}
	return *(*string)(v.p)
}

func (v Value) UnsafeBytes() []byte {
	if v.IsSmallString() {
		panic("immutable string")
	}
	return *(*[]byte)(v.p)
}

// When isInt == true, use intValue, otherwise floatValue
func (v Value) Num() (floatValue float64, intValue int64, isInt bool) {
	if v.p == int64Marker {
		return float64(int64(v.v)), int64(v.v), true
	}
	x := math.Float64frombits(v.v)
	return x, int64(x), false
}

func (v Value) Int() int64 { _, i, _ := v.Num(); return i }

func (v Value) Float() float64 { f, _, _ := v.Num(); return f }

func (v Value) Bool() bool { return v.p == trueMarker }

func (v Value) Map() *HashMap { return (*HashMap)(v.p) }

// Func cast value to function
func (v Value) Func() *Func { return (*Func)(v.p) }

func (v Value) WrappedFunc() *Func { return v.Interface().(*WrappedFunc).Func }

// Interface returns the interface{} representation of Value
func (v Value) Interface() interface{} {
	switch v.Type() {
	case typ.Bool:
		return v.Bool()
	case typ.Number:
		vf, vi, vIsInt := v.Num()
		if vIsInt {
			return vi
		}
		return vf
	case typ.String:
		return v.Str()
	case typ.Map:
		return v.Map()
	case typ.Func:
		return v.Func()
	case typ.Interface:
		return *(*interface{})(v.p)
	}
	return nil
}

func (v Value) puintptr() uintptr { return uintptr(v.p) }

func (v Value) unsafeint() int64 { return int64(v.v) }

func (v Value) unsafefloat() float64 { return math.Float64frombits(v.v) }

// ReflectValue returns reflect.Value based on reflect.Type
func (v Value) ReflectValue(t reflect.Type) reflect.Value {
	if t == nil {
		return reflect.ValueOf(v.Interface())
	}
	if t == reflect.TypeOf(Value{}) {
		return reflect.ValueOf(v)
	}
	if v.Type() == typ.Nil && (t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface) {
		return reflect.Zero(t)
	}
	if v.Type() == typ.String && t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Uint8 {
		return reflect.ValueOf(v.UnsafeBytes())
	}
	if v.Type() == typ.Func && t.Kind() == reflect.Func {
		return reflect.MakeFunc(t, func(args []reflect.Value) (results []reflect.Value) {
			in := make([]Value, len(args))
			for i := range in {
				in[i] = Val(args[i].Interface())
			}
			out, err := v.Func().Call(in...)
			if err != nil {
				panic(err)
			}
			switch t.NumOut() {
			case 0:
			case 1:
				results = []reflect.Value{out.ReflectValue(t.Out(0))}
			default:
				out.MustMap("expect multiple returned arguments", 0)
				results = make([]reflect.Value, t.NumOut())
				for i := range results {
					results[i] = out.Map().Get(Int(int64(i))).ReflectValue(t.Out(i))
				}
			}
			return
		})
	}

	switch v.Type() {
	case typ.Number:
		if t.Kind() >= reflect.Int && t.Kind() <= reflect.Float64 {
			rv := reflect.ValueOf(v.Interface())
			return rv.Convert(t)
		}
	case typ.Map:
		a := v.Map()
		switch t.Kind() {
		case reflect.Slice:
			s := reflect.MakeSlice(t, len(a.Array()), len(a.Array()))
			for i, a := range a.Array() {
				s.Index(i).Set(a.ReflectValue(t.Elem()))
			}
			return s
		case reflect.Array:
			s := reflect.New(t).Elem()
			for i, a := range a.Array() {
				s.Index(i).Set(a.ReflectValue(t.Elem()))
			}
			return s
		case reflect.Map:
			s := reflect.MakeMap(t)
			kt, vt := t.Key(), t.Elem()
			a.Foreach(func(k, v Value) bool {
				s.SetMapIndex(k.ReflectValue(kt), v.ReflectValue(vt))
				return true
			})
			return s
		}
	}
	return reflect.ValueOf(v.Interface())
}

func (v Value) MustBool(msg string, a int) bool { return v.mustBe(typ.Bool, msg, a).Bool() }

func (v Value) MustStr(msg string, a int) string { return v.mustBe(typ.String, msg, a).String() }

func (v Value) MustNum(msg string, a int) Value { return v.mustBe(typ.Number, msg, a) }

func (v Value) MustMap(msg string, a int) *HashMap { return v.mustBe(typ.Map, msg, a).Map() }

func (v Value) MustFunc(msg string, a int) *Func { return v.mustBe(typ.Func, msg, a).Func() }

func (v Value) mustBe(t typ.ValueType, msg string, msgArg int) Value {
	if v.Type() != t {
		if strings.Contains(msg, "%d") {
			msg = fmt.Sprintf(msg, msgArg)
		}
		if msg != "" {
			panicf("%s: expect %v, got %v", msg, t, v.Type())
		}
		panicf("expect %v, got %v", t, v.Type())
	}
	return v
}

func (v Value) Equal(r Value) bool {
	if v == r {
		return true
	}
	switch v.Type() + r.Type() {
	case typ.String * 2:
		return r.Str() == v.Str()
	case typ.Interface * 2:
		return v.Interface() == r.Interface()
	}
	return false
}

func (v Value) HashCode() uint64 {
	code := uint64(5381)
	if v.Type() != typ.String || v.IsSmallString() {
		for _, r := range *(*[ValueSize]byte)(unsafe.Pointer(&v)) {
			code = code*33 + uint64(r)
		}
	} else {
		for _, r := range v.Str() {
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
	case typ.Bool:
		return strconv.FormatBool(v.Bool())
	case typ.Number:
		vf, vi, vIsInt := v.Num()
		if vIsInt {
			return strconv.FormatInt(vi, 10)
		}
		return strconv.FormatFloat(vf, 'f', -1, 64)
	case typ.String:
		if j {
			return strconv.Quote(v.Str())
		}
		return v.Str()
	case typ.Map:
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
	case typ.Func:
		return v.Func().String()
	case typ.Interface:
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
	if v.Type() == typ.String {
		return v.Str()
	}
	return d
}

func (v Value) IntDefault(d int64) int64 {
	if v.Type() == typ.Number {
		return v.Int()
	}
	return d
}
