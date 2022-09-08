package bas

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
	"reflect"
	"strconv"
	"unicode/utf8"
	"unsafe"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

//go:linkname strhash runtime.strhash
func strhash(p unsafe.Pointer, h uintptr) uintptr

//go:linkname memhash64 runtime.memhash64
func memhash64(p unsafe.Pointer, h uintptr) uintptr

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
	False   = Value{0, falseMarker}
	True    = Value{1, trueMarker}
)

const (
	ValueSize = unsafe.Sizeof(Value{})

	errNeedNumber           = "operator requires number, got %v"
	errNeedNumbers          = "operator requires numbers, got %v and %v"
	errNeedNumbersOrStrings = "operator requires numbers or strings, got %v and %v"
)

// Value is the basic data type used by the intepreter, an empty Value naturally represent nil.
type Value struct {
	v uint64
	p unsafe.Pointer
}

func (v Value) IsValue() {}

// Type returns the type of value.
func (v Value) Type() typ.ValueType {
	if uintptr(v.p)&0xffffffffffffff00 == baseStart {
		return typ.ValueType(uintptr(v.p) & 7)
	}
	return typ.ValueType(v.v)
}

func (v Value) pType() typ.ValueType {
	return typ.ValueType(uintptr(v.p) & 7)
}

// IsFalse returns true if value is falsy: nil, false, empty string or 0.
func (v Value) IsFalse() bool { return v.v == 0 }

// IsTrue returns true if value is not falsy.
func (v Value) IsTrue() bool { return v.v != 0 }

// IsInt64 returns true if value is an integer number.
func (v Value) IsInt64() bool { return v.p == int64Marker }

// IsObject returns true if value is an object.
func (v Value) IsObject() bool { return v.Type() == typ.Object }

// IsArray returns true if value is an array.
func (v Value) IsArray() bool {
	return v.Type() == typ.Native && v.Native().meta.Proto.HasPrototype(Proto.Array)
}

// IsNil returns true if value is nil.
func (v Value) IsNil() bool { return v == Nil }

// Bool creates a boolean value.
func Bool(v bool) Value {
	if v {
		return True
	}
	return False
}

// Float64 creates a number value.
func Float64(f float64) Value {
	if float64(int64(f)) == f {
		// if math.Floor(f) == f {
		return Value{v: uint64(int64(f)), p: int64Marker}
	}
	return Value{v: math.Float64bits(f), p: float64Marker}
}

// Int creates a number value.
func Int(i int) Value {
	return Int64(int64(i))
}

// Int64 creates a number value.
func Int64(i int64) Value {
	return Value{v: uint64(i), p: int64Marker}
}

// Str creates a string value.
func Str(s string) Value {
	if len(s) <= 8 { // payload 8b
		var x [8]byte
		switch len(s) {
		case 1:
			a := s[0]
			x = [8]byte{a, a, a, a, a, a, a, a + 1}
		case 2:
			a, b := s[0], s[1]
			x = [8]byte{a, b, a, b, a, b, a, b + 1}
		case 3:
			copy(x[:], s)
			copy(x[3:], s)
			x[6], x[7] = s[0], s[1]+1
		case 4, 5, 6, 7:
			copy(x[:], s)
			copy(x[len(s):], s)
			x[7]++
		case 8:
			if s == "\x00\x00\x00\x00\x00\x00\x00\x00" {
				return Value{v: uint64(typ.String), p: unsafe.Pointer(&s)}
			}
			copy(x[:], s)
		}
		return Value{
			v: binary.BigEndian.Uint64(x[:]),
			p: unsafe.Pointer(uintptr(smallStrMarker) + uintptr(len(s))*8),
		}
	}
	return Value{v: uint64(typ.String), p: unsafe.Pointer(&s)}
}

// UnsafeStr creates a string value from []byte, its content may change if []byte changed.
func UnsafeStr(b []byte) Value {
	var s string
	*(*[2]uintptr)(unsafe.Pointer(&s)) = *(*[2]uintptr)(unsafe.Pointer(&b))
	return Str(s)
}

// Byte creates a one-byte string value.
func Byte(s byte) Value {
	x := [8]byte{s, s, s, s, s, s, s, s + 1}
	return Value{v: binary.BigEndian.Uint64(x[:]), p: unsafe.Pointer(uintptr(smallStrMarker) + 8)}
}

// Rune creates a one-rune string value encoded in UTF-8.
func Rune(r rune) Value {
	x := [8]byte{}
	n := utf8.EncodeRune(x[:], r)
	switch n {
	case 1:
		x[1], x[2], x[3], x[4], x[5], x[6], x[7] = x[0], x[0], x[0], x[0], x[0], x[0], x[0]+1
	case 2:
		x[2], x[3], x[4], x[5], x[6], x[7] = x[0], x[1], x[0], x[1], x[0], x[1]+1
	case 3:
		x[3], x[4], x[5], x[6], x[7] = x[0], x[1], x[2], x[0], x[1]+1
	case 4:
		copy(x[4:], x[:4])
		x[7]++
	}
	return Value{v: binary.BigEndian.Uint64(x[:]), p: unsafe.Pointer(uintptr(smallStrMarker) + uintptr(n)*8)}
}

// Bytes creates a bytes array.
func Bytes(b []byte) Value {
	return NewNativeWithMeta(b, bytesArrayMeta).ToValue()
}

// Error creates an error, first argument can be nil, indicating that the returned error has no stacktrace.
func Error(e *Env, err error) Value {
	if err == nil {
		return Nil
	}
	if _, ok := err.(*ExecError); ok {
		return NewNativeWithMeta(err, errorNativeMeta).ToValue()
	}
	ee := &ExecError{root: err}
	if e != nil {
		ee.stacks = e.runtime.Stacktrace(true)
	}
	return NewNativeWithMeta(ee, errorNativeMeta).ToValue()
}

func Array(v ...Value) Value {
	return newArray(v...).ToValue()
}

// ValueOf creates a Value from any golang types.
func ValueOf(i interface{}) Value {
	switch v := i.(type) {
	case nil:
		return Nil
	case bool:
		return Bool(v)
	case float64:
		return Float64(v)
	case int:
		return Int64(int64(v))
	case int8:
		return Int64(int64(v))
	case int16:
		return Int64(int64(v))
	case int32:
		return Int64(int64(v))
	case int64:
		return Int64(v)
	case uint:
		return Int64(int64(uint64(v)))
	case uint8:
		return Int64(int64(uint64(v)))
	case uint16:
		return Int64(int64(uint64(v)))
	case uint32:
		return Int64(int64(uint64(v)))
	case uint64:
		return Int64(int64(v))
	case uintptr:
		return Int64(int64(v))
	case string:
		return Str(v)
	case []byte:
		return Bytes(v)
	case *Object:
		return v.ToValue()
	case []Value:
		return Array(v...)
	case Value:
		return v
	case error:
		return Error(nil, v)
	case reflect.Value:
		return ValueOf(v.Interface())
	case func(*Env):
		return Func(internal.UnnamedFunc(), v)
	}

	if rv := reflect.ValueOf(i); rv.Kind() == reflect.Func {
		rt := rv.Type()
		nf := func(env *Env) {
			var interopFuncs []func()
			rtNumIn := rt.NumIn()
			ins := make([]reflect.Value, 0, rtNumIn)
			if !rt.IsVariadic() {
				if env.Size() != rtNumIn {
					internal.Panic("native function expects %d arguments, got %d", rtNumIn, env.Size())
				}
				for i := 0; i < rtNumIn; i++ {
					ins = append(ins, toTypePtrStruct(env.Get(i), rt.In(i), &interopFuncs))
				}
			} else {
				if env.Size() < rtNumIn-1 {
					internal.Panic("native variadic function expects at least %d arguments, got %d", rtNumIn-1, env.Size())
				}
				for i := 0; i < rtNumIn-1; i++ {
					ins = append(ins, toTypePtrStruct(env.Get(i), rt.In(i), &interopFuncs))
				}
				for i := rtNumIn - 1; i < env.Size(); i++ {
					ins = append(ins, toTypePtrStruct(env.Get(i), rt.In(rtNumIn-1).Elem(), &interopFuncs))
				}
			}
			if outs := rv.Call(ins); len(outs) == 0 {
				env.A = Nil
			} else if len(outs) == 1 {
				env.A = ValueOf(outs[0].Interface())
			} else {
				env.A = NewNative(outs).ToValue()
			}
			for _, f := range interopFuncs {
				f()
			}
		}
		return Func("<"+rt.String()+">", nf)
	}
	return NewNative(i).ToValue()
}

func (v Value) isSmallString() bool {
	return uintptr(v.p) >= uintptr(smallStrMarker) && uintptr(v.p) <= uintptr(smallStrMarker)+8*8
}

// Str returns value as a string, Type() should be checked beforehand.
func (v Value) Str() string {
	if v.isSmallString() {
		sz := (uintptr(v.p) - uintptr(smallStrMarker)) / 8
		if sz == 0 {
			return ""
		}
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, v.v)
		buf = buf[:sz]
		return *(*string)(unsafe.Pointer(&buf))
	}
	return *(*string)(v.p)
}

// Int returns value as an int.
func (v Value) Int() int { return int(v.Int64()) }

// Int64 returns value as an int64 (floats will be truncated to integers), Type() should be checked beforehand.
func (v Value) Int64() int64 {
	if v.p == int64Marker {
		return int64(v.v)
	}
	return int64(math.Float64frombits(v.v))
}

// Float64 returns value as a float (integers will be promoted to floats), Type() should be checked beforehand.
func (v Value) Float64() float64 {
	if v.p == int64Marker {
		return float64(int64(v.v))
	}
	return math.Float64frombits(v.v)
}

// Bool returns value as a boolean, Type() should be checked beforehand.
func (v Value) Bool() bool { return v.p == trueMarker }

// Object returns value as an Object, Type() should be checked beforehand.
func (v Value) Object() *Object { return (*Object)(v.p) }

// Native returns value as a Native, Type() should be checked beforehand.
func (v Value) Native() *Native { return (*Native)(v.p) }

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
	case typ.Native:
		return v.Native().Unwrap()
	}
	return nil
}

func (v Value) unsafeAddr() uintptr { return uintptr(v.p) }

// UnsafeInt64 returns value as an int64 without any type-checkings.
func (v Value) UnsafeInt64() int64 { return int64(v.v) }

// AssertType asserts the type of value. 'msg' will be panicked out if failed.
func (v Value) AssertType(t typ.ValueType, msg string) Value {
	if v.Type() != t {
		if msg != "" {
			internal.Panic("%s: expects %v, got %v", msg, t, detail(v))
		}
		internal.Panic("expects %v, got %v", t, detail(v))
	}
	return v
}

// AssertPrototype asserts the prototype of value. 'msg' will be panicked out if failed.
func (v Value) AssertPrototype(p *Object, msg string) Value {
	if !HasPrototype(v, p) {
		if msg != "" {
			internal.Panic("%s: expects prototype %v, got %v", msg, detail(p.ToValue()), detail(v))
		}
		internal.Panic("expects prototype %v, got %v", detail(p.ToValue()), detail(v))
	}
	return v
}

// Equal returns true if two values are equal.
func (v Value) Equal(r Value) bool {
	if v == r {
		return true
	}
	return v.v == uint64(typ.String) && v.v == r.v && *(*string)(v.p) == *(*string)(r.p)
}

// HashCode returns the hash of value.
func (v Value) HashCode() uint32 {
	if typ.ValueType(v.v) == typ.String {
		return uint32(strhash(v.p, 0))
	}

	x := uint64(uintptr(v.p))
	xp := uintptr(unsafe.Pointer(&x))
	h := memhash64(unsafe.Pointer(xp^0), uintptr(v.v))
	return uint32(h)

	// x := v.v ^ uint64(uintptr(v.p))
	// h := uint32(x) ^ uint32(x>>32)
	// h ^= h >> 16
	// h *= 0x85ebca6b
	// h ^= h >> 13
	// // h *= 0xc2b2ae35
	// // h ^= h >> 16
	// return h
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
		if j == typ.MarshalToJSON || j == typ.MarshalToString {
			v.Object().rawPrint(p, j)
		} else {
			internal.WriteString(p, "@"+v.Object().Name())
		}
	case typ.Native:
		if j == typ.MarshalToJSON || j == typ.MarshalToString {
			v.Native().Marshal(p, j)
		} else {
			internal.WriteString(p, "@"+v.Native().meta.Name)
		}
	default:
		internal.WriteString(p, internal.IfStr(j == typ.MarshalToJSON, "null", "nil"))
	}
}
