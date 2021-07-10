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
	Undef   = Go(new(int))
	Zero    = Int(0)
	NullStr = Str("")
	False   = Bool(false)
	True    = Bool(true)
)

type ValueType byte

const (
	NIL  ValueType = 0
	BOOL ValueType = 1
	NUM  ValueType = 3
	STR  ValueType = 7
	MAP  ValueType = 15
	FUNC ValueType = 17
	GO   ValueType = 19

	ValueSize = int64(unsafe.Sizeof(Value{}))

	errNeedNumbers          = "operator requires numbers"
	errNeedNumbersOrStrings = "operator requires numbers or strings"
)

func (t ValueType) String() string {
	if t > GO {
		return "?"
	}
	return [...]string{"nil", "bool", "?", "number", "?", "?", "?", "string", "?", "?", "?", "?", "?", "?", "?", "map", "?", "function", "?", "golang"}[t]
}

// Value is the basic data type used by the intepreter, an empty Value naturally represent nil
type Value struct {
	v uint64
	p unsafe.Pointer
}

// Reverse-reference in 'parser' package
func (v Value) IsValue(parser.Node) {}

// Type returns the type of value, its logic should align IsFalse()
func (v Value) Type() ValueType {
	switch v.p {
	case int64Marker, float64Marker:
		return NUM
	case trueMarker, falseMarker:
		return BOOL
	case nil:
		return NIL
	}
	if uintptr(v.p) >= uintptr(smallStringMarker) && uintptr(v.p) <= uintptr(smallStringMarker)+8*8 {
		return STR
	}
	return ValueType(v.v)
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
		return Value{uint64(BOOL), trueMarker}
	}
	return Value{uint64(BOOL), falseMarker}
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
	x := &RHMap{items: m}
	for _, i := range x.items {
		if i != Nil {
			x.count++
		}
	}
	return x.Value()
}

// Map returns a map, kvs should be laid out as: key1, value1, key2, value2, ...
func Map(kvs ...Value) Value {
	t := NewMap(len(kvs) / 2)
	for i := 0; i < len(kvs)/2*2; i += 2 {
		t.Set(kvs[i], kvs[i+1])
	}
	return Value{v: uint64(MAP), p: unsafe.Pointer(t)}
}

// MapWithParent returns a map whose parent will be p
func MapWithParent(p *RHMap, kvs ...Value) Value {
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
	return Value{v: uint64(STR), p: unsafe.Pointer(&s)}
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
	return Str(*(*string)(unsafe.Pointer(&b)))
}

// Go creates Value from golang interface{}
// []Type (except []byte/[]Value), [..]Type and map[Type]Type will be left as is,
// to convert them recursively, use DeepGo instead
func Go(i interface{}) Value {
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
	case *RHMap:
		return v.Value()
	case []Value:
		return Array(v...)
	case *Func:
		return v.Value()
	case Value:
		return v
	case parser.CatchedError:
		return Go(v.Original)
	case reflect.Value:
		return Go(v.Interface())
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
				x[i] = Go(a)
			}
			return Array(x...)
		}
		if v.IsObject() {
			m := v.Map()
			x := NewMap(len(m))
			for k, v := range m {
				x.Set(Str(k), Go(v))
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
					return reflect.ValueOf(env.Get(i).GoType(t))
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
					env.A = Go(outs[0].Interface())
				default:
					a := make([]Value, len(outs))
					for i := range outs {
						a[i] = Go(outs[i].Interface())
					}
					env.A = Array(a...)
				}
			}
		}
		return (&Func{Name: "<native>", Native: nf}).Value()
	}
	return _interface(i)
}

func DeepGo(v interface{}) Value {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map:
		m := NewMap(rv.Len() + 1)
		iter := rv.MapRange()
		for iter.Next() {
			m.Set(DeepGo(iter.Key()), Go(iter.Value()))
		}
		return m.Value()
	case reflect.Array, reflect.Slice:
		a := make([]Value, rv.Len())
		for i := range a {
			a[i] = DeepGo(rv.Index(i))
		}
		return Array(a...)
	default:
		return Go(v)
	}
}

func _interface(i interface{}) Value {
	return Value{v: uint64(GO), p: unsafe.Pointer(&i)}
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

func (v Value) Map() *RHMap { return (*RHMap)(v.p) }

// Func cast value to function
func (v Value) Func() *Func { return (*Func)(v.p) }

func (v Value) WrappedFunc() *Func { return v.Go().(*WrappedFunc).Func }

// Go returns the interface{} representation of Value
func (v Value) Go() interface{} {
	switch v.Type() {
	case BOOL:
		return v.Bool()
	case NUM:
		vf, vi, vIsInt := v.Num()
		if vIsInt {
			return vi
		}
		return vf
	case STR:
		return v.Str()
	case MAP:
		return v.Map()
	case FUNC:
		return v.Func()
	case GO:
		return *(*interface{})(v.p)
	}
	return nil
}

func (v Value) puintptr() uintptr { return uintptr(v.p) }

func (v Value) unsafeint() int64 { return int64(v.v) }

func (v Value) unsafefloat() float64 { return math.Float64frombits(v.v) }

// GoType returns the interface{} representation of Value which will be converted to t if needed
func (v Value) GoType(t reflect.Type) interface{} {
	if t == nil {
		return v.Go()
	}
	if t == reflect.TypeOf(Value{}) {
		return v
	}

	switch v.Type() {
	case NUM:
		if t.Kind() >= reflect.Int && t.Kind() <= reflect.Float64 {
			rv := reflect.ValueOf(v.Go())
			rv = rv.Convert(t)
			return rv.Interface()
		}
	case MAP:
		a := v.Map()
		switch t.Kind() {
		case reflect.Slice:
			s := reflect.MakeSlice(t, len(a.Array()), len(a.Array()))
			for i, a := range a.Array() {
				s.Index(i).Set(reflect.ValueOf(a.GoType(t.Elem())))
			}
			return s.Interface()
		case reflect.Array:
			s := reflect.New(t).Elem()
			for i, a := range a.Array() {
				s.Index(i).Set(reflect.ValueOf(a.GoType(t.Elem())))
			}
			return s.Interface()
		case reflect.Map:
			s := reflect.MakeMap(t)
			kt, vt := t.Key(), t.Elem()
			a.Foreach(func(k, v Value) bool {
				s.SetMapIndex(reflect.ValueOf(k.GoType(kt)), reflect.ValueOf(v.GoType(vt)))
				return true
			})
			return s.Interface()
		}
	}
	return v.Go()
}

func (v Value) MustBool(msg string, a int) bool { return v.mustBe(BOOL, msg, a).Bool() }

func (v Value) MustStr(msg string, a int) string { return v.mustBe(STR, msg, a).String() }

func (v Value) MustNum(msg string, a int) Value { return v.mustBe(NUM, msg, a) }

func (v Value) MustMap(msg string, a int) *RHMap { return v.mustBe(MAP, msg, a).Map() }

func (v Value) MustFunc(msg string, a int) *Func { return v.mustBe(FUNC, msg, a).Func() }

func (v Value) mustBe(t ValueType, msg string, msgArg int) Value {
	if v.Type() != t {
		if strings.Contains(msg, "%d") {
			msg = fmt.Sprintf(msg, msgArg)
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
	case STR * 2:
		return r.Str() == v.Str()
	case GO * 2:
		return v.Go() == r.Go()
	}
	return false
}

func (v Value) HashCode() uint64 {
	code := uint64(5381)
	if v.Type() != STR || v.IsSmallString() {
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
	case BOOL:
		return strconv.FormatBool(v.Bool())
	case NUM:
		vf, vi, vIsInt := v.Num()
		if vIsInt {
			return strconv.FormatInt(vi, 10)
		}
		return strconv.FormatFloat(vf, 'f', -1, 64)
	case STR:
		if j {
			return strconv.Quote(v.Str())
		}
		return v.Str()
	case MAP:
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
	case FUNC:
		return v.Func().String()
	case GO:
		i := v.Go()
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
	if v.Type() == STR {
		return v.Str()
	}
	return d
}

func (v Value) IntDefault(d int64) int64 {
	if v.Type() == NUM {
		return v.Int()
	}
	return d
}
