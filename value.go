package nj

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"unicode/utf8"
	"unsafe"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
	"github.com/tidwall/gjson"
)

var (
	baseMarker = func() []byte {
		// Ensures baseMarker is at least 256 bytes long and its memory aligns with 256 bytes
		b := make([]byte, 512)
		for i := range b {
			if byte(uintptr(unsafe.Pointer(&b[i]))) == 0 {
				return b[i:]
			}
		}
		panic("memory")
	}()
	baseStart  = uintptr(unsafe.Pointer(&baseMarker[0]))
	baseLength = uintptr(len(baseMarker))
	baseEnd    = uintptr(unsafe.Pointer(&baseMarker[0])) + baseLength

	int64Marker    = unsafe.Pointer(&baseMarker[int(typ.Number)])
	float64Marker  = unsafe.Pointer(&baseMarker[int(typ.Number)+8])
	trueMarker     = unsafe.Pointer(&baseMarker[int(typ.Bool)])
	falseMarker    = unsafe.Pointer(&baseMarker[int(typ.Bool)+8])
	smallStrMarker = unsafe.Pointer(&baseMarker[int(typ.String)])
	int64Marker2   = uintptr(int64Marker) * 2

	Nil     = Value{}
	Zero    = Int64(0)
	NullStr = Str("")
	False   = Bool(false)
	True    = Bool(true)
)

const (
	ValueSize = unsafe.Sizeof(Value{})

	errNeedNumber           = "operator requires number, got %v"
	errNeedNumbers          = "operator requires numbers, got %v and %v"
	errNeedNumbersOrStrings = "operator requires numbers or strings, got %v and %v"
)

// Value is the basic data type used by the intepreter, an empty Value naturally represent nil
type Value struct {
	v uint64
	p unsafe.Pointer
}

func (v Value) IsValue() {}

// Type returns the type of value
func (v Value) Type() typ.ValueType {
	if uintptr(v.p)^baseStart < baseLength {
		// if uintptr(v.p) >= baseStart && uintptr(v.p) < baseEnd {
		return typ.ValueType(uintptr(v.p) & 7)
	}
	return typ.ValueType(v.v)
}

// IsFalse tests whether value is falsy: nil, false, empty string or 0
func (v Value) IsFalse() bool { return v.v == 0 || v.p == falseMarker }

func (v Value) IsTrue() bool { return !v.IsFalse() }

// IsInt64 tests whether value is an integer number
func (v Value) IsInt64() bool { return v.p == int64Marker }

func (v Value) IsObject() bool { return v.Type() == typ.Object }

func (v Value) IsNil() bool { return v == Nil }

// Bool creates a boolean value
func Bool(v bool) Value {
	if v {
		return Value{uint64(typ.Bool), trueMarker}
	}
	return Value{uint64(typ.Bool), falseMarker}
}

// Float64 creates a number value
func Float64(f float64) Value {
	if float64(int64(f)) == f {
		// if math.Floor(f) == f {
		return Value{v: uint64(int64(f)), p: int64Marker}
	}
	return Value{v: math.Float64bits(f), p: float64Marker}
}

// Int creates a number value
func Int(i int) Value {
	return Int64(int64(i))
}

// Int64 creates a number value
func Int64(i int64) Value {
	return Value{v: uint64(i), p: int64Marker}
}

// Array creates an array consists of given arguments
func Array(m ...Value) Value {
	return (&List{store: m}).Value()
}

// Obj creates a map from `kvs`, which should be laid out as: key1, value1, key2, value2, ...
func Obj(kvs ...Value) Value {
	t := NewObject(len(kvs) / 2)
	for i := 0; i < len(kvs)/2*2; i += 2 {
		t.Set(kvs[i], renameFuncName(kvs[i], kvs[i+1]))
	}
	return Value{v: uint64(typ.Object), p: unsafe.Pointer(t)}
}

// TableMerge merges key-value pairs from `src` into `dst` if both of them are tables
func TableMerge(dst Value, src *Object) Value {
	var t *Object
	switch dst.Type() {
	case typ.Object:
		t = dst.Object()
	case typ.Nil:
		t = NewObject(1)
	default:
		return dst
	}
	t.Merge(src)
	return t.Value()
}

// Proto creates a table whose parent will be set to `p`
func Proto(p *Object, kvs ...Value) Value {
	return Obj(kvs...).Object().SetParent(p).Value()
}

// Str creates a string value
func Str(s string) Value {
	if len(s) <= 8 { // payload 8b
		x := [8]byte{byte(len(s))}
		copy(x[:], s)
		return Value{
			v: binary.BigEndian.Uint64(x[:]),
			p: unsafe.Pointer(uintptr(smallStrMarker) + uintptr(len(s))*8),
		}
	}
	return Value{v: uint64(typ.String), p: unsafe.Pointer(&s)}
}

// Byte creates a one-byte string value
func Byte(s byte) Value {
	x := [8]byte{s}
	return Value{v: binary.BigEndian.Uint64(x[:]), p: unsafe.Pointer(uintptr(smallStrMarker) + 8)}
}

// Rune creates a one-rune string value encoded in UTF-8
func Rune(r rune) Value {
	x := [8]byte{}
	n := utf8.EncodeRune(x[:], r)
	return Value{v: binary.BigEndian.Uint64(x[:]), p: unsafe.Pointer(uintptr(smallStrMarker) + uintptr(n)*8)}
}

// Bytes creates a string value from bytes
func Bytes(b []byte) Value { return Str(*(*string)(unsafe.Pointer(&b))) }

// Val creates a `Value` from golang `interface{}`
// `slice`, `array` and `map` will be left as is (except []Value), to convert them recursively, use ValRec instead
func Val(i interface{}) Value {
	switch v := i.(type) {
	case nil:
		return Value{}
	case bool:
		return Bool(v)
	case float64:
		return Float64(v)
	case int:
		return Int64(int64(v))
	case int64:
		return Int64(v)
	case string:
		return Str(v)
	case *Object:
		return v.Value()
	case []Value:
		return Array(v...)
	case Value:
		return v
	case internal.CatchedError:
		return Val(v.Original)
	case reflect.Value:
		return Val(v.Interface())
	case gjson.Result:
		if v.Type == gjson.String {
			return Str(v.Str)
		} else if v.Type == gjson.Number {
			return Float64(v.Float())
		} else if v.Type == gjson.True || v.Type == gjson.False {
			return Bool(v.Bool())
		} else if v.IsArray() {
			x := make([]Value, 0, len(v.Raw)/10)
			v.ForEach(func(k, v gjson.Result) bool { x = append(x, Val(v)); return true })
			return Array(x...)
		} else if v.IsObject() {
			x := NewObject(len(v.Raw) / 10)
			v.ForEach(func(k, v gjson.Result) bool { x.Set(Val(k), Val(v)); return true })
			return x.Value()
		}
		return Nil
	}

	rv := reflect.ValueOf(i)
	if k := rv.Kind(); k >= reflect.Int && k <= reflect.Int64 {
		return Int64(rv.Int())
	} else if k >= reflect.Uint && k <= reflect.Uintptr {
		return Int64(int64(rv.Uint()))
	} else if (k == reflect.Ptr || k == reflect.Interface) && rv.IsNil() {
		return Nil
	} else if k == reflect.Func {
		nf, _ := i.(func(*Env))
		if nf == nil {
			rt := rv.Type()
			nf = func(env *Env) {
				rtNumIn := rt.NumIn()
				ins := make([]reflect.Value, 0, rtNumIn)
				if !rt.IsVariadic() {
					if env.Size() != rtNumIn {
						internal.Panic("call native function, expect %d arguments, got %d", rtNumIn, env.Size())
					}
					for i := 0; i < rtNumIn; i++ {
						ins = append(ins, env.Get(i).ReflectValue(rt.In(i)))
					}
				} else {
					if env.Size() < rtNumIn-1 {
						internal.Panic("call native variadic function, expect at least %d arguments, got %d", rtNumIn-1, env.Size())
					}
					for i := 0; i < rtNumIn-1; i++ {
						ins = append(ins, env.Get(i).ReflectValue(rt.In(i)))
					}
					for i := rtNumIn - 1; i < env.Size(); i++ {
						ins = append(ins, env.Get(i).ReflectValue(rt.In(rtNumIn-1).Elem()))
					}
				}
				if outs := rv.Call(ins); len(outs) == 0 {
				} else if len(outs) == 1 {
					env.A = Val(outs[0].Interface())
				} else {
					env.A = Array(valReflectValues(outs)...)
				}
			}
		}
		return (&Object{callable: &FuncBody{Name: "<" + rv.Type().String() + ">", Native: nf}}).Value()
	}
	return intf(i)
}

func ValRec(v interface{}) Value {
	switch rv := reflect.ValueOf(v); rv.Kind() {
	case reflect.Map:
		m := NewObject(rv.Len() + 1)
		for iter := rv.MapRange(); iter.Next(); {
			m.Set(ValRec(iter.Key()), Val(iter.Value()))
		}
		return m.Value()
	case reflect.Array, reflect.Slice:
		a := make([]Value, rv.Len())
		for i := range a {
			a[i] = ValRec(rv.Index(i))
		}
		return Array(a...)
	}
	return Val(v)
}

func valReflectValues(args []reflect.Value) (a []Value) {
	for i := range args {
		a = append(a, Val(args[i].Interface()))
	}
	return
}

func intf(i interface{}) Value {
	return Value{v: uint64(typ.Native), p: unsafe.Pointer(&i)}
}

func showType(v Value) string {
	switch vt := v.Type(); vt {
	case typ.Number, typ.Bool, typ.Native:
		return v.JSONString()
	case typ.String:
		if v.StrLen() <= 32 {
			return v.JSONString()
		}
		return strconv.Quote(v.Str()[:32] + "...")
	case typ.Object:
		return "{" + v.Object().Name() + "}"
	default:
		return vt.String()
	}
}

func (v Value) isSmallString() bool {
	return uintptr(v.p) >= uintptr(smallStrMarker) && uintptr(v.p) <= uintptr(smallStrMarker)+8*8
}

// Str returns value as a string without checking Type()
func (v Value) Str() string {
	if v.isSmallString() {
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, v.v)
		buf = buf[:(uintptr(v.p)-uintptr(smallStrMarker))/8]
		return *(*string)(unsafe.Pointer(&buf))
	}
	return *(*string)(v.p)
}

// StrLen returns the length of string without checking Type()
func (v Value) StrLen() int {
	if v.isSmallString() {
		return int(uintptr(v.p)-uintptr(smallStrMarker)) / 8
	}
	return len(*(*string)(v.p))
}

// Int returns value as an integer without checking Type()
func (v Value) Int() int { return int(v.Int64()) }

// Int64 returns value as an integer without checking Type()
func (v Value) Int64() int64 {
	if v.p == int64Marker {
		return int64(v.v)
	}
	return int64(math.Float64frombits(v.v))
}

// Float64 returns value as a float without checking Type()
func (v Value) Float64() float64 {
	if v.p == int64Marker {
		return float64(int64(v.v))
	}
	return math.Float64frombits(v.v)
}

// Bool returns value as a boolean without checking Type()
func (v Value) Bool() bool { return v.p == trueMarker }

// Object returns value as a table without checking Type()
func (v Value) Object() *Object { return (*Object)(v.p) }

func (v Value) Array() *List { return (*List)(v.p) }

// Interface returns value as an interface{}
func (v Value) Interface() interface{} {
	switch v.Type() {
	case typ.Bool:
		return v.Bool()
	case typ.Number:
		if v.IsInt64() {
			return v.Int64()
		}
		return v.Float64()
	case typ.String:
		return v.Str()
	case typ.Object:
		return v.Object()
	case typ.Array:
		return v.Array()
	case typ.Native:
		return *(*interface{})(v.p)
	}
	return nil
}

func (v Value) ptr() uintptr { return uintptr(v.p) }

func (v Value) unsafeInt() int64 { return int64(v.v) }

// ReflectValue returns value as a reflect.Value based on reflect.Type
func (v Value) ReflectValue(t reflect.Type) reflect.Value {
	if t == nil {
		return reflect.ValueOf(v.Interface())
	} else if t == reflect.TypeOf(Value{}) {
		return reflect.ValueOf(v)
	} else if t.Implements(ioWriterType) || t.Implements(ioReaderType) || t.Implements(ioCloserType) {
		return reflect.ValueOf(ValueIO(v))
	} else if vt := v.Type(); vt == typ.Nil && (t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface) {
		return reflect.Zero(t)
	} else if v.IsObject() && t.Kind() == reflect.Func {
		return reflect.MakeFunc(t, func(args []reflect.Value) (results []reflect.Value) {
			out := v.Object().MustCall(valReflectValues(args)...)
			if to := t.NumOut(); to == 1 {
				results = []reflect.Value{out.ReflectValue(t.Out(0))}
			} else if to > 1 {
				out.Is(typ.Array, "ReflectValue: expect multiple returned arguments")
				results = make([]reflect.Value, t.NumOut())
				for i := range results {
					results[i] = out.Array().Get(Int(i)).ReflectValue(t.Out(i))
				}
			}
			return
		})
	} else if vt == typ.Number && t.Kind() >= reflect.Int && t.Kind() <= reflect.Float64 {
		return reflect.ValueOf(v.Interface()).Convert(t)
	} else if vt == typ.Array {
		switch a := v.Array(); t.Kind() {
		case reflect.Slice:
			s := reflect.MakeSlice(t, a.Len(), a.Len())
			a.Foreach(func(k, v Value) bool { s.Index(k.Int()).Set(v.ReflectValue(t.Elem())); return true })
			return s
		case reflect.Array:
			s := reflect.New(t).Elem()
			a.Foreach(func(k, v Value) bool { s.Index(k.Int()).Set(v.ReflectValue(t.Elem())); return true })
			return s
		}
	} else if vt == typ.Object {
		switch a := v.Object(); t.Kind() {
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

func (v Value) MustBool(msg string) bool { return v.Is(typ.Bool, msg).Bool() }

func (v Value) MustStr(msg string) string { return v.Is(typ.String, msg).String() }

func (v Value) MustStrLen(msg string) int { return v.Is(typ.String, msg).StrLen() }

func (v Value) MustNum(msg string) Value { return v.Is(typ.Number, msg) }

func (v Value) MustInt64(msg string) int64 { return v.Is(typ.Number, msg).Int64() }

func (v Value) MustInt(msg string) int { return v.Is(typ.Number, msg).Int() }

func (v Value) MustFloat64(msg string) float64 { return v.Is(typ.Number, msg).Float64() }

func (v Value) MustTable(msg string) *Object { return v.Is(typ.Object, msg).Object() }

func (v Value) Is(t typ.ValueType, msg string) Value {
	if v.Type() != t {
		if msg != "" {
			internal.Panic("%s: expect %v, got %v", msg, t, showType(v))
		}
		internal.Panic("expect %v, got %v", t, showType(v))
	}
	return v
}

// Equal tests whether two values are equal
func (v Value) Equal(r Value) bool {
	if v == r {
		return true
	}
	return v.v == uint64(typ.String) && v.v == r.v && *(*string)(v.p) == *(*string)(r.p)
}

func (v Value) HashCode() uint64 {
	if typ.ValueType(v.v) == typ.String {
		code := uint64(5381)
		for _, r := range v.Str() {
			code = code*33 + uint64(r)
		}
		return code
	}
	return v.v * uint64(uintptr(v.p))
}

func (v Value) String() string {
	return v.toString(&bytes.Buffer{}, 0, false).String()
}

func (v Value) JSONString() string {
	return v.toString(&bytes.Buffer{}, 0, true).String()
}

func (v Value) MarshalJSON() ([]byte, error) {
	return v.toString(&bytes.Buffer{}, 0, true).Bytes(), nil
}

func (v Value) toString(p *bytes.Buffer, lv int, j bool) *bytes.Buffer {
	if lv > 10 {
		p.WriteString(ifstr(j, "{}", "..."))
		return p
	}
	switch v.Type() {
	case typ.Bool:
		p.WriteString(strconv.FormatBool(v.Bool()))
	case typ.Number:
		if v.IsInt64() {
			p.WriteString(strconv.FormatInt(v.Int64(), 10))
		} else {
			p.WriteString(strconv.FormatFloat(v.Float64(), 'f', -1, 64))
		}
	case typ.String:
		p.WriteString(ifquote(j, v.Str()))
	case typ.Object:
		v.Object().rawPrint(p, lv, j, false)
	case typ.Array:
		p.WriteString("[")
		v.Array().Foreach(func(i, v Value) bool {
			v.toString(p, lv+1, j)
			p.WriteString(",")
			return true
		})
		closeBuffer(p, "]")
	case typ.Native:
		i := v.Interface()
		if s, ok := i.(fmt.Stringer); ok {
			p.WriteString(ifquote(j, s.String()))
		} else if s, ok := i.(error); ok {
			p.WriteString(ifquote(j, s.Error()))
		} else {
			p.WriteString(ifquote(j, "<"+reflect.TypeOf(i).String()+">"))
		}
	default:
		p.WriteString(ifstr(j, "null", "nil"))
	}
	return p
}

func (v Value) ToStr(d string) string {
	if v.Type() == typ.String {
		return v.Str()
	}
	return d
}

func (v Value) ToInt(d int) int {
	return int(v.ToInt64(int64(d)))
}

func (v Value) ToInt64(d int64) int64 {
	if v.Type() == typ.Number {
		return v.Int64()
	}
	return d
}

func (v Value) ToFloat64(d float64) float64 {
	if v.Type() == typ.Number {
		return v.Float64()
	}
	return d
}

func (v Value) ToObject() *Object {
	if v.Type() != typ.Object {
		return nil
	}
	return v.Object()
}

func (v Value) ToTableGets(key string) Value {
	if v.Type() != typ.Object {
		return Nil
	}
	return v.Object().Gets(key)
}

func (v Value) ForEach(f func(k, v Value) bool) {
	switch v.Type() {
	case typ.Object:
		v.Object().Foreach(f)
	case typ.Array:
		v.Array().Foreach(f)
	default:
		internal.Panic("can't iterate %v", v.Type())
	}
}

func (v Value) Len() int {
	switch v.Type() {
	case typ.String:
		return v.StrLen()
	case typ.Array:
		return v.Array().Len()
	case typ.Object:
		return v.Object().Len()
	case typ.Nil:
		return 0
	case typ.Number, typ.Bool:
		internal.Panic("can't measure length of %v", v.Type())
	}
	return reflectLen(v.Interface())
}
