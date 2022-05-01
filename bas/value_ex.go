package bas

import (
	"reflect"
	"unsafe"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

type MaybeValue Value

func (v MaybeValue) Str(defaultValue string) string {
	switch t := Value(v).Type(); t {
	case typ.String:
		return Value(v).Str()
	case typ.Nil:
		return defaultValue
	case typ.Native:
		if buf, ok := Value(v).Native().Unwrap().([]byte); ok {
			return *(*string)(unsafe.Pointer(&buf))
		}
		fallthrough
	default:
		panic("Str: expects string, bytes or nil, got " + simpleString(Value(v)))
	}
}

func (v MaybeValue) Bool() bool {
	switch t := Value(v).Type(); t {
	case typ.Number, typ.Bool:
		return Value(v).IsTrue()
	case typ.Nil:
		return false
	default:
		panic("Bool: expects boolean or nil, got " + simpleString(Value(v)))
	}
}

func (v MaybeValue) Int64(defaultValue int64) int64 {
	switch t := Value(v).Type(); t {
	case typ.Number:
		return Value(v).Int64()
	case typ.Nil:
		return defaultValue
	default:
		panic("Int64: expects integer number or nil, got " + simpleString(Value(v)))
	}
}

func (v MaybeValue) Int(defaultValue int) int {
	return int(v.Int64(int64(defaultValue)))
}

func (v MaybeValue) Float64(defaultValue float64) float64 {
	switch t := Value(v).Type(); t {
	case typ.Number:
		return Value(v).Float64()
	case typ.Nil:
		return defaultValue
	default:
		panic("Float64: expects float number or nil, got " + simpleString(Value(v)))
	}
}

func (v MaybeValue) Object(defaultValue *Object) *Object {
	switch t := Value(v).Type(); t {
	case typ.Object:
		return Value(v).Object()
	case typ.Nil:
		return defaultValue
	default:
		panic("Object: expects object or nil, got " + simpleString(Value(v)))
	}
}

func (v MaybeValue) Func(defaultValue *Object) *Object {
	o := v.Object(defaultValue)
	if o == defaultValue {
		return o
	}
	if o != nil {
		if o.IsCallable() {
			return o
		}
		panic("Func: expects function or nil, got " + simpleString(Value(v)))
	}
	return o
}

func ToError(v Value) error {
	if Value(v).Type() == typ.Native && Value(v).Native().meta.Proto.HasPrototype(errorNativeMeta.Proto) {
		return Value(v).Native().Unwrap().(*ExecError)
	}
	panic("ToError: not error: " + simpleString(v))
}

func ToErrorRootCause(v Value) interface{} {
	if Value(v).Type() == typ.Native && Value(v).Native().meta.Proto.HasPrototype(errorNativeMeta.Proto) {
		return Value(v).Native().Unwrap().(*ExecError).root
	}
	panic("ToErrorRootCause: not error: " + simpleString(v))
}

func ToBytes(v Value) []byte {
	if Value(v).Type() == typ.Native && Value(v).Native().meta.Proto.HasPrototype(bytesArrayMeta.Proto) {
		return Value(v).Native().Unwrap().([]byte)
	}
	panic("ToBytes: not []byte: " + simpleString(v))
}

func ToReadonlyBytes(v Value) []byte {
	switch v.Type() {
	case typ.Nil:
		return nil
	case typ.Native:
		if v.Native().meta.Proto.HasPrototype(bytesArrayMeta.Proto) {
			return Value(v).Native().Unwrap().([]byte)
		}
	case typ.String:
		var s struct {
			a string
			i int
		}
		s.a = v.Str()
		s.i = len(s.a)
		return *(*[]byte)(unsafe.Pointer(&s))
	}
	panic("ToReadonlyBytes: not []byte or string: " + simpleString(v))
}

func IsBytes(v Value) bool {
	return v.Type() == typ.Native && v.Native().meta.Proto.HasPrototype(bytesArrayMeta.Proto)
}

func IsError(v Value) bool {
	return v.Type() == typ.Native && v.Native().meta.Proto.HasPrototype(errorNativeMeta.Proto)
}

func DeepEqual(a, b Value) bool {
	if a.Equal(b) {
		return true
	}
	if at, bt := a.Type(), b.Type(); at == bt {
		switch at {
		case typ.Native:
			if a.IsArray() && b.IsArray() {
				flag := a.Native().Len() == b.Native().Len()
				if !flag {
					return false
				}
				for i := 0; flag && i < a.Native().Len(); i++ {
					flag = DeepEqual(b.Native().Get(i), a.Native().Get(i))
				}
				return flag
			}
		case typ.Object:
			flag := a.Object().Len() == b.Object().Len()
			if !flag {
				return false
			}
			a.Object().Foreach(func(k Value, v *Value) bool {
				flag = DeepEqual(b.Object().Find(k), *v)
				return flag
			})
			return flag
		}
	}
	return false
}

func lessStr(a, b Value) bool {
	if a.isSmallString() && b.isSmallString() {
		al := (a.UnsafeAddr() - uintptr(smallStrMarker)) / 8 * 8
		bl := (b.UnsafeAddr() - uintptr(smallStrMarker)) / 8 * 8
		av := a.v >> (64 - al)
		bv := b.v >> (64 - bl)
		return av < bv
	}
	return a.Str() < b.Str()
}

func Less(a, b Value) bool {
	at, bt := a.Type(), b.Type()
	if at != bt {
		return at < bt
	}
	switch at {
	case typ.Number:
		if a.IsInt64() && b.IsInt64() {
			return a.UnsafeInt64() < b.UnsafeInt64()
		}
		return a.Float64() < b.Float64()
	case typ.String:
		return lessStr(a, b)
	}
	return a.UnsafeAddr() < b.UnsafeAddr()
}

func HasPrototype(a Value, p *Object) bool {
	switch a.Type() {
	case typ.Nil:
		return p == nil
	case typ.Object:
		return a.Object().HasPrototype(p)
	case typ.Bool:
		return p == Proto.Bool
	case typ.Number:
		return p == Proto.Float || (a.IsInt64() && p == Proto.Int)
	case typ.String:
		return p == Proto.Str
	case typ.Native:
		return a.Native().meta.Proto.HasPrototype(p)
	}
	return false
}

func IsCallable(a Value) bool {
	return a.Type() == typ.Object && a.Object().IsCallable()
}

// ToType convert Value to reflect.Value based on reflect.Type
func ToType(v Value, t reflect.Type) reflect.Value {
	return toTypePtrStruct(v, t, nil)
}

func toTypePtrStruct(v Value, t reflect.Type, interopFuncs *[]func()) reflect.Value {
	if t == nil {
		return reflect.ValueOf(v.Interface())
	}
	if t == valueType {
		return reflect.ValueOf(v)
	}

	vt := v.Type()
	if interopFuncs != nil && t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct && vt == typ.Object {
		this := v.Object()
		vp := objectToStruct(this, t.Elem(), interopFuncs)
		*interopFuncs = append(*interopFuncs, func() {
			for i := 0; i < vp.NumField(); i++ {
				f := vp.Field(i)
				if n := vp.Type().Field(i).Name; n[0] >= 'A' && n[0] <= 'Z' {
					this.SetProp(n, ValueOf(f.Interface()))
				}
			}
		})
		return vp.Addr()
	}

	if vt == typ.Nil && (t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface) {
		return reflect.Zero(t)
	}
	if t.Implements(ioWriterType) || t.Implements(ioReaderType) || t.Implements(ioCloserType) {
		return reflect.ValueOf(ValueIO(v))
	}
	if t.Implements(errType) {
		return reflect.ValueOf(ToError(v))
	}
	if IsCallable(v) && t.Kind() == reflect.Func {
		return reflect.MakeFunc(t, func(args []reflect.Value) (results []reflect.Value) {
			var a []Value
			for i := range args {
				a = append(a, ValueOf(args[i].Interface()))
			}
			out := Call(v.Object(), a...)
			if to := t.NumOut(); to == 1 {
				results = []reflect.Value{toTypePtrStruct(out, t.Out(0), interopFuncs)}
			} else if to > 1 {
				out.AssertType(typ.Native, "ToType: requires function to return multiple arguments")
				results = make([]reflect.Value, t.NumOut())
				for i := range results {
					results[i] = toTypePtrStruct(out.Native().Get(i), t.Out(i), interopFuncs)
				}
			}
			return
		})
	}
	if vt == typ.Number && t.Kind() >= reflect.Int && t.Kind() <= reflect.Float64 {
		return reflect.ValueOf(v.Interface()).Convert(t)
	}
	if vt == typ.Native {
		a := v.Native()
		if t == reflect.TypeOf(a.Unwrap()) {
			return reflect.ValueOf(a.Unwrap())
		}
		switch t.Kind() {
		case reflect.Slice:
			a.AssertPrototype(Proto.Array, "ToType")
			s := reflect.MakeSlice(t, a.Len(), a.Len())
			for i := 0; i < a.Len(); i++ {
				s.Index(i).Set(toTypePtrStruct(a.Get(i), t.Elem(), interopFuncs))
			}
			return s
		case reflect.Array:
			a.AssertPrototype(Proto.Array, "ToType")
			s := reflect.New(t).Elem()
			for i := 0; i < a.Len(); i++ {
				s.Index(i).Set(toTypePtrStruct(a.Get(i), t.Elem(), interopFuncs))
			}
			return s
		}
	}
	if vt == typ.Object && t.Kind() == reflect.Map {
		s := reflect.MakeMap(t)
		kt, vt := t.Key(), t.Elem()
		v.Object().Foreach(func(k Value, v *Value) bool {
			s.SetMapIndex(toTypePtrStruct(k, kt, interopFuncs), toTypePtrStruct(*v, vt, interopFuncs))
			return true
		})
		return s
	}
	if vt == typ.Object && t.Kind() == reflect.Struct {
		return objectToStruct(v.Object(), t, interopFuncs)
	}
	if vt == typ.Bool && t.Kind() == reflect.Bool {
		return reflect.ValueOf(v.Bool())
	}
	if vt == typ.String && t.Kind() == reflect.String {
		return reflect.ValueOf(v.Str())
	}

	panic("ToType: failed to convert " + simpleString(v) + " to " + t.String())
}

func objectToStruct(src *Object, t reflect.Type, interopFuncs *[]func()) reflect.Value {
	vp := reflect.New(t)
	s := vp.Elem()
	src.Foreach(func(k Value, v *Value) bool {
		field := k.AssertType(typ.String, "ToStruct: field name").Str()
		if field == "" || field[0] < 'A' || field[0] > 'Z' {
			return true
		}
		f := s.FieldByName(field)
		if !f.IsValid() {
			internal.Panic("ToStruct: field %q not found", field)
		}
		f.Set(toTypePtrStruct(*v, f.Type(), interopFuncs))
		return true
	})
	return s
}

func Len(v Value) int {
	switch v.Type() {
	case typ.String:
		if v.isSmallString() {
			return int(uintptr(v.p)-uintptr(smallStrMarker)) / 8
		}
		return len(*(*string)(v.p))
	case typ.Native:
		return v.Native().Len()
	case typ.Object:
		return v.Object().Len()
	case typ.Nil:
		return 0
	}
	panic("can't measure length of " + simpleString(v))
}
