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
	case typ.Array:
		if buf, ok := Value(v).Array().Unwrap().([]byte); ok {
			return *(*string)(unsafe.Pointer(&buf))
		}
		fallthrough
	default:
		panic("Str: expects string or nil, got " + simpleString(Value(v)))
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
	if Value(v).Type() == typ.Array && Value(v).Array().meta.Proto.IsPrototype(errorArrayMeta.Proto) {
		return Value(v).Array().Unwrap().(*ExecError)
	}
	panic("ToError: not error: " + simpleString(v))
}

func ToBytes(v Value) []byte {
	if Value(v).Type() == typ.Array && Value(v).Array().meta.Proto.IsPrototype(bytesArrayMeta.Proto) {
		return Value(v).Array().Unwrap().([]byte)
	}
	panic("ToBytes: not []byte: " + simpleString(v))
}

func ToReadonlyBytes(v Value) []byte {
	switch v.Type() {
	case typ.Array:
		if v.Array().meta.Proto.IsPrototype(bytesArrayMeta.Proto) {
			return Value(v).Array().Unwrap().([]byte)
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
	return v.Type() == typ.Array && v.Array().meta.Proto.IsPrototype(bytesArrayMeta.Proto)
}

func IsError(v Value) bool {
	return v.Type() == typ.Array && v.Array().meta.Proto.IsPrototype(errorArrayMeta.Proto)
}

func DeepEqual(a, b Value) bool {
	if a.Equal(b) {
		return true
	}
	if at, bt := a.Type(), b.Type(); at == bt {
		switch at {
		case typ.Array:
			flag := a.Array().Len() == b.Array().Len()
			if !flag {
				return false
			}
			a.Array().Foreach(func(k int, v Value) bool {
				flag = DeepEqual(b.Array().Get(k), v)
				return flag
			})
			return flag
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

func IsPrototype(a Value, p *Object) bool {
	switch a.Type() {
	case typ.Nil:
		return p == nil
	case typ.Object:
		return a.Object().IsPrototype(p)
	case typ.Bool:
		return p == Proto.Bool
	case typ.Number:
		return p == Proto.Float || (a.IsInt64() && p == Proto.Int)
	case typ.String:
		return p == Proto.Str
	case typ.Array:
		return a.Array().meta.Proto.IsPrototype(p)
	}
	return false
}

func IsCallable(a Value) bool {
	return a.Type() == typ.Object && a.Object().IsCallable()
}

// ToType convert Value to reflect.Value based on reflect.Type
func ToType(v Value, t reflect.Type) reflect.Value {
	if t == nil {
		return reflect.ValueOf(v.Interface())
	}
	if t == valueType {
		return reflect.ValueOf(v)
	}

	vt := v.Type()
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
				results = []reflect.Value{ToType(out, t.Out(0))}
			} else if to > 1 {
				out.AssertType(typ.Array, "ToType: expect multiple returned arguments")
				results = make([]reflect.Value, t.NumOut())
				for i := range results {
					results[i] = ToType(out.Array().Get(i), t.Out(i))
				}
			}
			return
		})
	}
	if vt == typ.Number && t.Kind() >= reflect.Int && t.Kind() <= reflect.Float64 {
		return reflect.ValueOf(v.Interface()).Convert(t)
	}
	if vt == typ.Array {
		a := v.Array()
		if t == reflect.TypeOf(a.Unwrap()) {
			return reflect.ValueOf(a.Unwrap())
		}
		switch t.Kind() {
		case reflect.Slice:
			s := reflect.MakeSlice(t, a.Len(), a.Len())
			a.Foreach(func(k int, v Value) bool { s.Index(k).Set(ToType(v, t.Elem())); return true })
			return s
		case reflect.Array:
			s := reflect.New(t).Elem()
			a.Foreach(func(k int, v Value) bool { s.Index(k).Set(ToType(v, t.Elem())); return true })
			return s
		}
	}
	if vt == typ.Object && t.Kind() == reflect.Map {
		s := reflect.MakeMap(t)
		kt, vt := t.Key(), t.Elem()
		v.Object().Foreach(func(k Value, v *Value) bool {
			s.SetMapIndex(ToType(k, kt), ToType(*v, vt))
			return true
		})
		return s
	}
	if vt == typ.Object && t.Kind() == reflect.Struct {
		return objectToStruct(v.Object(), t, false)
	}
	if vt == typ.Object && t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
		return objectToStruct(v.Object(), t.Elem(), true)
	}
	if vt == typ.Bool && t.Kind() == reflect.Bool {
		return reflect.ValueOf(v.Bool())
	}
	if vt == typ.String && t.Kind() == reflect.String {
		return reflect.ValueOf(v.Str())
	}
	if vt == typ.Native {
		if i := v.Interface(); reflect.TypeOf(i) == t {
			return reflect.ValueOf(i)
		}
	}
	panic("ToType: failed to convert " + simpleString(v) + " to " + t.String())
}

func objectToStruct(src *Object, t reflect.Type, ptr bool) reflect.Value {
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
		f.Set(ToType(*v, f.Type()))
		return true
	})
	if ptr {
		src.SetProp("_istruct", intf(s.Interface()))
		src.SetMethod("interop", func(e *Env) {
			this := e.Object(-1)
			v := reflect.ValueOf(e.ThisProp("_istruct"))
			for i := 0; i < v.NumField(); i++ {
				f := v.Field(i)
				if n := v.Type().Field(i).Name; n[0] >= 'A' && n[0] <= 'Z' {
					this.SetProp(n, ValueOf(f.Interface()))
				}
			}
		}, "")
		return vp
	}
	return s
}

func Len(v Value) int {
	switch v.Type() {
	case typ.String:
		if v.isSmallString() {
			return int(uintptr(v.p)-uintptr(smallStrMarker)) / 8
		}
		return len(*(*string)(v.p))
	case typ.Array:
		return v.Array().Len()
	case typ.Object:
		return v.Object().Len()
	case typ.Nil:
		return 0
	case typ.Native:
		return reflect.ValueOf(v.Interface()).Len()
	}
	panic("can't measure length of " + simpleString(v))
}
