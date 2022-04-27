package bas

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"strconv"
	"unicode/utf8"
	"unsafe"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
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
	// baseEnd    = uintptr(unsafe.Pointer(&baseMarker[0])) + baseLength

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

// IsTrue returns the same way as !IsFalse()
func (v Value) IsTrue() bool { return !v.IsFalse() }

// IsInt64 tests whether value is an integer number
func (v Value) IsInt64() bool { return v.p == int64Marker }

// IsObject tests whether value is an object
func (v Value) IsObject() bool { return v.Type() == typ.Object }

// IsNil tests whether value is nil
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

// Str creates a string value
func Str(s string) Value {
	if len(s) <= 8 { // payload 8b
		var x [8]byte
		switch len(s) {
		case 1:
			a := s[0]
			x = [8]byte{a, a, a, a, a, a, a, a}
		case 2:
			a, b := s[0], s[1]
			x = [8]byte{a, b, a, b, a, b, a, b}
		case 3:
			copy(x[:], s)
			copy(x[3:], s)
			copy(x[6:], s)
		case 4, 5, 6, 7:
			copy(x[:], s)
			copy(x[len(s):], s)
		case 8:
			copy(x[:], s)
		}
		return Value{
			v: binary.BigEndian.Uint64(x[:]),
			p: unsafe.Pointer(uintptr(smallStrMarker) + uintptr(len(s))*8),
		}
	}
	return Value{v: uint64(typ.String), p: unsafe.Pointer(&s)}
}

// UnsafeStr creates a string value from []byte, its content may change if []byte changed
func UnsafeStr(s []byte) Value {
	return Str(*(*string)(unsafe.Pointer(&s)))
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

// Bytes creates a bytes array
func Bytes(b []byte) Value {
	return NewTypedArray(b, bytesArrayMeta).ToValue()
}

// ValueOf creates a `Value` from golang `interface{}`
func ValueOf(i interface{}) Value {
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
	case []byte:
		return Bytes(v)
	case *Object:
		return v.ToValue()
	case []Value:
		return NewArray(v...).ToValue()
	case Value:
		return v
	case error:
		return Error(nil, v)
	case reflect.Value:
		return ValueOf(v.Interface())
	case os.FileInfo:
		return fileInfo(v).ToValue()
	}

	rv := reflect.ValueOf(i)
	if k := rv.Kind(); k >= reflect.Int && k <= reflect.Int64 {
		return Int64(rv.Int())
	} else if k >= reflect.Uint && k <= reflect.Uintptr {
		return Int64(int64(rv.Uint()))
	} else if (k == reflect.Ptr || k == reflect.Interface) && rv.IsNil() {
		return Nil
	} else if k == reflect.Array || k == reflect.Slice {
		return NewTypedArray(i, GetTypedArrayMeta(i)).ToValue()
	} else if k == reflect.Chan {
		return NewObject(0).SetProp("_ch", intf(rv.Interface())).SetPrototype(Proto.Channel).ToValue()
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
						ins = append(ins, ToType(env.Get(i), rt.In(i)))
					}
				} else {
					if env.Size() < rtNumIn-1 {
						internal.Panic("call native variadic function, expect at least %d arguments, got %d", rtNumIn-1, env.Size())
					}
					for i := 0; i < rtNumIn-1; i++ {
						ins = append(ins, ToType(env.Get(i), rt.In(i)))
					}
					for i := rtNumIn - 1; i < env.Size(); i++ {
						ins = append(ins, ToType(env.Get(i), rt.In(rtNumIn-1).Elem()))
					}
				}
				if outs := rv.Call(ins); len(outs) == 0 {
					env.A = Nil
				} else if len(outs) == 1 {
					env.A = ValueOf(outs[0].Interface())
				} else {
					env.A = NewTypedArray(outs, GetTypedArrayMeta(outs)).ToValue()
				}
			}
		}
		return Func("<"+rv.Type().String()+">", nf, "")
	} else if k == reflect.Map {
		o := NewObject(rv.Len())
		for iter := rv.MapRange(); iter.Next(); {
			o.Set(ValueOf(iter.Key().Interface()), ValueOf(iter.Value().Interface()))
		}
		return o.ToValue()
	}
	return intf(i)
}

func intf(i interface{}) Value {
	return Value{v: uint64(typ.Native), p: unsafe.Pointer(&i)}
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

// Object returns value as an object without checking Type()
func (v Value) Object() *Object { return (*Object)(v.p) }

// Array returns value as a sequence without checking Type()
func (v Value) Array() *Array { return (*Array)(v.p) }

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
		return v.Array().Unwrap()
	case typ.Native:
		return *(*interface{})(v.p)
	}
	return nil
}

func (v Value) UnsafeAddr() uintptr { return uintptr(v.p) }

func (v Value) UnsafeInt64() int64 { return int64(v.v) }

func (v Value) AssertType(t typ.ValueType, msg string) Value {
	if v.Type() != t {
		if msg != "" {
			internal.Panic("%s: expect %v, got %v", msg, t, simpleString(v))
		}
		internal.Panic("expect %v, got %v", t, simpleString(v))
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
	return (v.v)*uint64(uintptr(v.p)) ^ (v.v >> 33)
	// v.v ^= uint64(uintptr(v.p))
	// v.v ^= v.v >> 33
	// v.v *= 0xff51afd7ed558ccd
	// v.v ^= v.v >> 33
	// v.v *= 0xc4ceb9fe1a85ec53
	// v.v ^= v.v >> 33
	// return v.v
}

func (v Value) String() string {
	p := &bytes.Buffer{}
	v.Stringify(p, typ.MarshalToString)
	return p.String()
}

func (v Value) JSONString() string {
	p := &bytes.Buffer{}
	v.Stringify(p, typ.MarshalToJSON)
	return p.String()
}

func (v Value) MarshalJSON() ([]byte, error) {
	p := &bytes.Buffer{}
	v.Stringify(p, typ.MarshalToJSON)
	return p.Bytes(), nil
}

func (v Value) Stringify(p io.Writer, j typ.MarshalType) {
	switch v.Type() {
	case typ.Bool:
		internal.WriteString(p, strconv.FormatBool(v.Bool()))
	case typ.Number:
		if v.IsInt64() {
			internal.WriteString(p, strconv.FormatInt(v.Int64(), 10))
		} else {
			internal.WriteString(p, strconv.FormatFloat(v.Float64(), 'f', -1, 64))
		}
	case typ.String:
		internal.WriteString(p, internal.IfQuote(j == typ.MarshalToJSON, v.Str()))
	case typ.Object:
		v.Object().rawPrint(p, j, false)
	case typ.Array:
		v.Array().Marshal(p, j)
	case typ.Native:
		i := v.Interface()
		if s, ok := i.(fmt.Stringer); ok {
			internal.WriteString(p, internal.IfQuote(j == typ.MarshalToJSON, s.String()))
		} else if s, ok := i.(error); ok {
			internal.WriteString(p, internal.IfQuote(j == typ.MarshalToJSON, s.Error()))
		} else {
			internal.WriteString(p, internal.IfQuote(j == typ.MarshalToJSON, "<"+reflect.TypeOf(i).String()+">"))
		}
	default:
		internal.WriteString(p, internal.IfStr(j == typ.MarshalToJSON, "null", "nil"))
	}
}

func (v Value) Maybe() MaybeValue { return MaybeValue(v) }
