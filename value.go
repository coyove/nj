package nj

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
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
	Zero    = Int(0)
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

// IsInt tests whether value is an integer number
func (v Value) IsInt() bool { return v.p == int64Marker }

// Bool creates a boolean value
func Bool(v bool) Value {
	if v {
		return Value{uint64(typ.Bool), trueMarker}
	}
	return Value{uint64(typ.Bool), falseMarker}
}

// Float creates a number value
func Float(f float64) Value {
	if float64(int64(f)) == f {
		// if math.Floor(f) == f {
		return Value{v: uint64(int64(f)), p: int64Marker}
	}
	return Value{v: math.Float64bits(f), p: float64Marker}
}

// Int creates a number value
func Int(i int64) Value {
	return Value{v: uint64(i), p: int64Marker}
}

// Array creates an array consists of given arguments
func Array(m ...Value) Value {
	x := &Table{items: m}
	for _, i := range x.items {
		if i != Nil {
			x.count++
		}
	}
	return x.Value()
}

// Map creates a map from `kvs`, which should be laid out as: key1, value1, key2, value2, ...
func Map(kvs ...Value) Value {
	t := NewTable(len(kvs) / 2)
	for i := 0; i < len(kvs)/2*2; i += 2 {
		k, v := kvs[i], kvs[i+1]
		if v.Type() == typ.Func && v.Func().Name == internal.UnnamedFunc {
			v.Func().Name = k.String()
		}
		t.Set(k, v)
	}
	return Value{v: uint64(typ.Table), p: unsafe.Pointer(t)}
}

// TableMerge merges key-value pairs from `src` into `dst`
func TableMerge(dst Value, src Value) Value {
	var t *Table
	switch dst.Type() {
	case typ.Table:
		t = dst.Table()
	case typ.Nil:
		t = NewTable(1)
	default:
		return dst
	}
	if src.Type() == typ.Table {
		t.Merge(src.Table())
	}
	return t.Value()
}

// TableProto creates a table whose parent will be set to `p`
func TableProto(p *Table, kvs ...Value) Value {
	m := Map(kvs...)
	m.Table().SetParent(p)
	return m
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
		return Float(v)
	case int:
		return Int(int64(v))
	case int64:
		return Int(v)
	case string:
		return Str(v)
	case *Table:
		return v.Value()
	case []Value:
		return Array(v...)
	case *Function:
		return v.Value()
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
			return Float(v.Float())
		} else if v.Type == gjson.True || v.Type == gjson.False {
			return Bool(v.Bool())
		} else if v.IsArray() {
			x := make([]Value, 0, len(v.Raw)/10)
			v.ForEach(func(k, v gjson.Result) bool { x = append(x, Val(v)); return true })
			return Array(x...)
		} else if v.IsObject() {
			x := NewTable(len(v.Raw) / 10)
			v.ForEach(func(k, v gjson.Result) bool { x.Set(Val(k), Val(v)); return true })
			return x.Value()
		}
		return Nil
	}

	rv := reflect.ValueOf(i)
	if k := rv.Kind(); k >= reflect.Int && k <= reflect.Int64 {
		return Int(rv.Int())
	} else if k >= reflect.Uint && k <= reflect.Uintptr {
		return Int(int64(rv.Uint()))
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
		return (&Function{FuncBody: &FuncBody{Name: "<" + rv.Type().String() + ">", Native: nf}}).Value()
	}
	return intf(i)
}

func ValRec(v interface{}) Value {
	switch rv := reflect.ValueOf(v); rv.Kind() {
	case reflect.Map:
		m := NewTable(rv.Len() + 1)
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

func stringType(v Value) string {
	switch vt := v.Type(); vt {
	case typ.Number, typ.Bool, typ.Native:
		return v.JSONString()
	case typ.String:
		if v.StrLen() <= 32 {
			return v.JSONString()
		}
		return strconv.Quote(v.Str()[:32] + "...")
	case typ.Table:
		return "{" + v.Table().Name() + "}"
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

// Int returns value as an int without checking Type()
func (v Value) Int() int64 {
	if v.p == int64Marker {
		return int64(v.v)
	}
	return int64(math.Float64frombits(v.v))
}

// Float returns value as a float without checking Type()
func (v Value) Float() float64 {
	if v.p == int64Marker {
		return float64(int64(v.v))
	}
	return math.Float64frombits(v.v)
}

// Bool returns value as a boolean without checking Type()
func (v Value) Bool() bool { return v.p == trueMarker }

// Table returns value as a table without checking Type()
func (v Value) Table() *Table { return (*Table)(v.p) }

// Func returns value as a function without checking Type()
func (v Value) Func() *Function { return (*Function)(v.p) }

// Interface returns value as an interface{}
func (v Value) Interface() interface{} {
	switch v.Type() {
	case typ.Bool:
		return v.Bool()
	case typ.Number:
		if v.IsInt() {
			return v.Int()
		}
		return v.Float()
	case typ.String:
		return v.Str()
	case typ.Table:
		return v.Table()
	case typ.Func:
		return v.Func()
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
	} else if vt == typ.Func && t.Kind() == reflect.Func {
		return reflect.MakeFunc(t, func(args []reflect.Value) (results []reflect.Value) {
			out, err := v.Func().Call(valReflectValues(args)...)
			internal.PanicErr(err)
			if to := t.NumOut(); to == 1 {
				results = []reflect.Value{out.ReflectValue(t.Out(0))}
			} else if to > 1 {
				out.mustBe(typ.Table, "ReflectValue: expect multiple returned arguments", 0)
				results = make([]reflect.Value, t.NumOut())
				for i := range results {
					results[i] = out.Table().Get(Int(int64(i))).ReflectValue(t.Out(i))
				}
			}
			return
		})
	} else if vt == typ.Number && t.Kind() >= reflect.Int && t.Kind() <= reflect.Float64 {
		return reflect.ValueOf(v.Interface()).Convert(t)
	} else if vt == typ.Table {
		switch a := v.Table(); t.Kind() {
		case reflect.Slice:
			s := reflect.MakeSlice(t, len(a.ArrayPart()), len(a.ArrayPart()))
			for i, a := range a.ArrayPart() {
				s.Index(i).Set(a.ReflectValue(t.Elem()))
			}
			return s
		case reflect.Array:
			s := reflect.New(t).Elem()
			for i, a := range a.ArrayPart() {
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

func (v Value) MustBool(msg string) bool { return v.mustBe(typ.Bool, msg, 0).Bool() }

func (v Value) MustStr(msg string) string { return v.mustBe(typ.String, msg, 0).String() }

func (v Value) MustStrLen(msg string) int { return v.mustBe(typ.String, msg, 0).StrLen() }

func (v Value) MustNum(msg string) Value { return v.mustBe(typ.Number, msg, 0) }

func (v Value) MustInt(msg string) int64 { return v.mustBe(typ.Number, msg, 0).Int() }

func (v Value) MustFloat(msg string) float64 { return v.mustBe(typ.Number, msg, 0).Float() }

func (v Value) MustTable(msg string) *Table { return v.mustBe(typ.Table, msg, 0).Table() }

func (v Value) MustFunc(msg string) *Function {
	if vt := v.Type(); vt == typ.Table {
		return v.Table().GetString("__call").MustFunc(msg)
	} else if vt == typ.Func {
	} else {
		internal.Panic(msg+ifstr(msg != "", ": ", "")+"expect function or callable table, got %v", stringType(v))
	}
	return v.Func()
}

func (v Value) mustBe(t typ.ValueType, msg string, msgArg int) Value {
	if v.Type() != t {
		if strings.Contains(msg, "%d") {
			msg = fmt.Sprintf(msg, msgArg)
		}
		if msg != "" {
			internal.Panic("%s: expect %v, got %v", msg, t, stringType(v))
		}
		internal.Panic("expect %v, got %v", t, stringType(v))
	}
	return v
}

func (v Value) Recv(k string) Value {
	if v.Type() != typ.Table {
		internal.Panic("method expects receiver, got %v, did you misuse 'table.key' and 'table:key'?", stringType(v))
	}
	return v.Table().GetString(k)
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
	return v.v ^ uint64(uintptr(v.p))
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
		if v.IsInt() {
			p.WriteString(strconv.FormatInt(v.Int(), 10))
		} else {
			p.WriteString(strconv.FormatFloat(v.Float(), 'f', -1, 64))
		}
	case typ.String:
		p.WriteString(ifquote(j, v.Str()))
	case typ.Table:
		m := v.Table()
		if sf := m.GetString("__str"); sf.Type() == typ.Func {
			if v, err := sf.Func().Call(); err != nil {
				p.WriteString(fmt.Sprintf("<table.__str: %v>", err))
			} else {
				v.toString(p, lv+1, j)
			}
			return p
		}
		m.rawPrint(p, lv, j, false)
	case typ.Func:
		p.WriteString(ifquote(j, v.Func().String()))
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

func (v Value) MaybeStr(d string) string {
	if v.Type() == typ.String {
		return v.Str()
	}
	return d
}

func (v Value) MaybeInt(d int64) int64 {
	if v.Type() == typ.Number {
		return v.Int()
	}
	return d
}

func (v Value) MaybeFloat(d float64) float64 {
	if v.Type() == typ.Number {
		return v.Float()
	}
	return d
}

func (v Value) MaybeTableGetString(key string) Value {
	if v.Type() != typ.Table {
		return Nil
	}
	return v.Table().GetString(key)
}
