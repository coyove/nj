package bas

import (
	"fmt"
	"io"
	"reflect"
	"runtime"
	"strconv"
	"sync"

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
		if IsBytes(v) {
			return w.Write(v.Native().Unwrap().([]byte))
		}
	case typ.String:
		return internal.WriteString(w, v.Str())
	}
	v.Stringify(w, typ.MarshalToString)
	return 1, nil
}

func IsBytes(v Value) bool {
	return v.Type() == typ.Native && v.Native().meta.Proto.HasPrototype(&Proto.Bytes)
}

func IsError(v Value) bool {
	return v.Type() == typ.Native && v.Native().meta.Proto.HasPrototype(&Proto.Error)
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

func (a Value) Less(b Value) bool {
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

func (a Value) HasPrototype(p *Object) bool {
	switch a.Type() {
	case typ.Nil:
		return p == nil
	case typ.Object:
		return a.Object().HasPrototype(p)
	case typ.Bool:
		return p == &Proto.Bool
	case typ.Number:
		return p == &Proto.Float || (a.IsInt64() && p == &Proto.Int)
	case typ.String:
		return p == &Proto.Str
	case typ.Native:
		return a.Native().HasPrototype(p)
	}
	return false
}

// ToType converts value to reflect.Value based on reflect.Type.
// The result, even not being Zero, may be illegal to use in certain calls.
func (v Value) ToType(t reflect.Type) reflect.Value {
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
					results = []reflect.Value{out.ToType(t.Out(0))}
				} else if to > 1 {
					if !out.IsArray() {
						internal.Panic("ToType: function should return %d arguments (sig: %v)", to, t)
					}
					results = make([]reflect.Value, t.NumOut())
					for i := range results {
						results[i] = out.Native().Get(i).ToType(t.Out(i))
					}
				}
				return
			})
		}
		if t.Kind() == reflect.Map {
			s := reflect.MakeMap(t)
			kt, vt := t.Key(), t.Elem()
			v.Object().Foreach(func(k Value, v *Value) bool {
				s.SetMapIndex(k.ToType(kt), v.ToType(vt))
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
					s.Index(i).Set(a.Get(i).ToType(t.Elem()))
				}
				return s
			case reflect.Array:
				s := reflect.New(t).Elem()
				for i := 0; i < a.Len(); i++ {
					s.Index(i).Set(a.Get(i).ToType(t.Elem()))
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

func (v Value) Len() int {
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
		if f := v.Object().fun; f != nil && f != objDefaultFun {
			return v.Object().funcSig()
		}
		return v.Object().Name() + "{}"
	case typ.Native:
		a := v.Native()
		if a.IsUntypedArray() {
			return fmt.Sprintf("array.%d", a.Len())
		}
		return a.meta.Name
	case typ.Number, typ.Bool:
		return v.String()
	case typ.String:
		ln := (v).Len()
		if ln < 16 {
			return strconv.Quote(v.Str())
		}
		return fmt.Sprintf("string.%d", ln)
	default:
		return v.Type().String()
	}
}

func Fprintf(w io.Writer, f string, values ...Value) {
	args := make([]interface{}, 0, len(values))
	for _, v := range values {
		if v.Type() == typ.Number {
			args = append(args, internal.SprintfNumber{Int: v.Int64(), Float: v.Float64(), IsInt: v.IsInt64()})
		} else {
			args = append(args, v.Interface())
		}
	}
	internal.Fprintf(w, f, args...)
}

func Fprint(w io.Writer, values ...Value) {
	for _, v := range values {
		v.Stringify(w, typ.MarshalToString)
	}
}

func multiMap(e *Env, fun *Object, t Value, n int) Value {
	if n < 1 || n > runtime.NumCPU()*1e3 {
		internal.Panic("invalid number of goroutines: %v", n)
	}

	type payload struct {
		i int
		k Value
		v *Value
	}

	work := func(e *Env, fun *Object, outError *error, p payload) {
		if p.i == -1 {
			res, err := fun.TryCall(e, p.k, *p.v)
			if err != nil {
				*outError = err
			} else {
				*p.v = res
			}
		} else {
			res, err := fun.TryCall(e, Int(p.i), p.k)
			if err != nil {
				*outError = err
			} else {
				t.Native().Set(p.i, res)
			}
		}
	}

	var outError error
	if n == 1 {
		if t.IsArray() {
			for i := 0; outError == nil && i < t.Native().Len(); i++ {
				work(e, fun, &outError, payload{i, t.Native().Get(i), nil})
			}
		} else {
			t.Object().Foreach(func(k Value, v *Value) bool {
				work(e, fun, &outError, payload{-1, k, v})
				return outError == nil
			})
		}
	} else {
		var in = make(chan payload, t.Len())
		var wg sync.WaitGroup
		wg.Add(n)
		for i := 0; i < n; i++ {
			go func(e *Env) {
				defer func() {
					wg.Done()
					if r := recover(); r != nil {
						outError = fmt.Errorf("map fatal error: %v", r)
					}
				}()
				for p := range in {
					if outError != nil {
						return
					}
					work(e, fun, &outError, p)
				}
			}(e.Copy())
		}

		if t.IsArray() {
			for i := 0; i < t.Native().Len(); i++ {
				in <- payload{i, t.Native().Get(i), nil}
			}
		} else {
			t.Object().Foreach(func(k Value, v *Value) bool {
				in <- payload{-1, k, v}
				return true
			})
		}
		close(in)

		wg.Wait()
	}
	if outError != nil {
		return Error(e, outError)
	}
	return t
}
