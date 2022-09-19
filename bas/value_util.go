package bas

import (
	"fmt"
	"io"
	"reflect"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

var (
	ioWriterType = reflect.TypeOf((*io.Writer)(nil)).Elem()
	ioReaderType = reflect.TypeOf((*io.Reader)(nil)).Elem()
	ioCloserType = reflect.TypeOf((*io.Closer)(nil)).Elem()
	errType      = reflect.TypeOf((*error)(nil)).Elem()
	valueType    = reflect.TypeOf(Value{})
)

func ToError(v Value) error {
	if IsError(v) {
		return v.Native().Unwrap().(*ExecError)
	}
	panic("ToError: not error: " + detail(v))
}

func ToErrorRootCause(v Value) interface{} {
	if IsError(v) {
		return v.Native().Unwrap().(*ExecError).root
	}
	panic("ToErrorRootCause: not error: " + detail(v))
}

func Write(w io.Writer, v Value) (int, error) {
	switch v.Type() {
	case typ.Nil:
		return 0, nil
	case typ.Native:
		if v.Native().meta.Proto.HasPrototype(bytesArrayMeta.Proto) {
			return w.Write(v.Native().Unwrap().([]byte))
		}
	case typ.String:
		return internal.WriteString(w, v.Str())
	}
	v.Stringify(w, typ.MarshalToString)
	return 1, nil
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
				flag = DeepEqual(b.Object().Get(k), *v)
				return flag
			})
			return flag
		}
	}
	return false
}

func lessStr(a, b Value) bool {
	if a.isSmallString() && b.isSmallString() {
		al := (a.unsafeAddr() - uintptr(smallStrMarker)) / 8 * 8
		bl := (b.unsafeAddr() - uintptr(smallStrMarker)) / 8 * 8
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
	return a.unsafeAddr() < b.unsafeAddr()
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
		return a.Native().HasPrototype(p)
	}
	return false
}

// ToType converts 'v' to reflect.Value based on reflect.Type.
// The result, even not being Zero, may be illegal to use in certain calls.
func ToType(v Value, t reflect.Type) reflect.Value {
	if t == valueType {
		return reflect.ValueOf(v)
	}
	if t == nil {
		return reflect.ValueOf(v.Interface())
	}

	switch vt := v.Type(); vt {
	case typ.Nil:
		if t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
			return reflect.Zero(t)
		}
	case typ.Object:
		if t.Kind() == reflect.Func {
			return reflect.MakeFunc(t, func(args []reflect.Value) (results []reflect.Value) {
				var a []Value
				for i := range args {
					if i == len(args)-1 && t.IsVariadic() {
						// TODO: performance
						for j := 0; j < args[i].Len(); j++ {
							a = append(a, ValueOf(args[i].Index(j).Interface()))
						}
					} else {
						a = append(a, ValueOf(args[i].Interface()))
					}
				}
				out := v.Object().Call(nil, a...)
				if to := t.NumOut(); to == 1 {
					results = []reflect.Value{ToType(out, t.Out(0))}
				} else if to > 1 {
					if !out.IsArray() {
						internal.Panic("ToType: function should return %d arguments (sig: %v)", to, t)
					}
					results = make([]reflect.Value, t.NumOut())
					for i := range results {
						results[i] = ToType(out.Native().Get(i), t.Out(i))
					}
				}
				return
			})
		}
		if t.Kind() == reflect.Map {
			s := reflect.MakeMap(t)
			kt, vt := t.Key(), t.Elem()
			v.Object().Foreach(func(k Value, v *Value) bool {
				s.SetMapIndex(ToType(k, kt), ToType(*v, vt))
				return true
			})
			return s
		}
		if t.Implements(ioWriterType) || t.Implements(ioReaderType) || t.Implements(ioCloserType) {
			return reflect.ValueOf(valueIO(v))
		}
	case typ.Native:
		a := v.Native().Unwrap()
		if t.Implements(ioWriterType) || t.Implements(ioReaderType) || t.Implements(ioCloserType) {
			return reflect.ValueOf(a)
		}
		if t.Implements(errType) {
			return reflect.ValueOf(ToError(v))
		}
		if t == reflect.TypeOf(a) {
			return reflect.ValueOf(a)
		}
		if v.IsArray() {
			switch a := v.Native(); t.Kind() {
			case reflect.Slice:
				s := reflect.MakeSlice(t, a.Len(), a.Len())
				for i := 0; i < a.Len(); i++ {
					s.Index(i).Set(ToType(a.Get(i), t.Elem()))
				}
				return s
			case reflect.Array:
				s := reflect.New(t).Elem()
				for i := 0; i < a.Len(); i++ {
					s.Index(i).Set(ToType(a.Get(i), t.Elem()))
				}
				return s
			}
		}
	case typ.Number:
		if t.Kind() >= reflect.Int && t.Kind() <= reflect.Float64 {
			return reflect.ValueOf(v.Interface()).Convert(t)
		}
	case typ.Bool:
		if t.Kind() == reflect.Bool {
			return reflect.ValueOf(v.Bool())
		}
	case typ.String:
		if t.Kind() == reflect.String {
			return reflect.ValueOf(v.Str())
		}
	}
	if t.Kind() == reflect.Interface {
		return reflect.ValueOf(v.Interface())
	}
	panic("ToType: failed to convert " + detail(v) + " to " + t.String())
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
	panic("can't measure length of " + detail(v))
}

func setObjectRecv(v, r Value) Value {
	if v.IsObject() {
		v.Object().this = r
	}
	return v
}

func detail(v Value) string {
	switch vt := v.Type(); vt {
	case typ.Object:
		if v.Object().fun != nil {
			return v.Object().funcSig()
		}
		return v.Object().Name() + "{}"
	case typ.Native:
		a := v.Native()
		if a.IsInternalArray() {
			return fmt.Sprintf("array(%d)", a.Len())
		}
		return fmt.Sprintf("native(%s)", a.meta.Name)
	case typ.Number:
		if v.IsInt64() {
			return fmt.Sprintf("int64(%d)", v.Int64())
		}
		return fmt.Sprintf("float64(%f)", v.Float64())
	case typ.Bool:
		return internal.IfStr(v.Bool(), "true", "false")
	case typ.String:
		return fmt.Sprintf("string(%d)", Len(v))
	default:
		return v.Type().String()
	}
}
