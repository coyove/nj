package bas

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
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

func reflectLoad(v interface{}, key Value) Value {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("reflectLoad %T.%s: %v", v, key, r))
		}
	}()

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map:
		v := rv.MapIndex(key.ReflectValue(rv.Type().Key()))
		if v.IsValid() {
			return ValueOf(v.Interface())
		}
	}

	k := key.Is(typ.String, "").Str()
	f := rv.MethodByName(k)
	if !f.IsValid() {
		if rv.Kind() == reflect.Ptr {
			f = rv.Elem().MethodByName(k)
		} else if rv.Kind() == reflect.Struct && rv.CanAddr() {
			f = rv.Addr().MethodByName(k)
		}
	}
	if !f.IsValid() {
		f = reflect.Indirect(rv).FieldByName(k)
		if !f.IsValid() {
			return Nil
		}
	}
	return ValueOf(f.Interface())
}

func reflectStore(subject interface{}, key, value Value) {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("reflectStore %T.%s: %v", subject, key, r))
		}
	}()

	rv := reflect.ValueOf(subject)
	if rv.Kind() == reflect.Map {
		rv.SetMapIndex(key.ReflectValue(rv.Type().Key()), value.ReflectValue(rv.Type().Elem()))
		return
	}

	rv = reflect.Indirect(rv)
	k := key.Is(typ.String, "").Str()
	f := rv.FieldByName(k)
	if !f.IsValid() || !f.CanAddr() {
		internal.Panic("reflect: %q not assignable in %v", k, subject)
	}
	f.Set(value.ReflectValue(f.Type()))
}

func or(a, b interface{}) interface{} {
	if a != nil {
		return a
	}
	return b
}

func setObjectRecv(v, r Value) Value {
	if v.IsObject() {
		v.Object().this = r
	}
	return v
}

func simpleString(v Value) string {
	switch vt := v.Type(); vt {
	case typ.Number, typ.Bool, typ.Native:
		return v.JSONString()
	case typ.String:
		if v.StrLen() <= 32 {
			return v.JSONString()
		}
		return strconv.Quote(v.Str()[:32] + "...")
	case typ.Object:
		if v.Object().fun != nil { // including named objects
			return v.Object().fun.String()
		}
		return "{" + v.Object().Name() + "}"
	case typ.Array:
		a := v.Array()
		if a.Typed() {
			return fmt.Sprintf("array(%s)", a.meta.Name)
		}
		return fmt.Sprintf("array(%d)", a.Len())
	default:
		return vt.String()
	}
}

func sprintf(env *Env, start int, p io.Writer) {
	f := env.Str(start)
	tmp := bytes.Buffer{}
	popi := start
	for len(f) > 0 {
		idx := strings.Index(f, "%")
		if idx == -1 {
			fmt.Fprint(p, f)
			break
		} else if idx == len(f)-1 {
			internal.Panic("unexpected '%%' at end")
		}
		fmt.Fprint(p, f[:idx])
		if f[idx+1] == '%' {
			p.Write([]byte("%"))
			f = f[idx+2:]
			continue
		}
		tmp.Reset()
		tmp.WriteByte('%')
		expecting := typ.Nil
		for f = f[idx+1:]; len(f) > 0 && expecting == typ.Nil; {
			switch f[0] {
			case 'b', 'd', 'o', 'O', 'c', 'e', 'E', 'f', 'F', 'g', 'G':
				expecting = typ.Number
			case 's', 'q', 'U':
				expecting = typ.String
			case 'x', 'X':
				expecting = typ.String + typ.Number
			case 'v':
				expecting = typ.Native
			case 't':
				expecting = typ.Bool
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', '-', '+', '#', ' ':
			default:
				internal.Panic("unexpected verb: '%c'", f[0])
			}
			tmp.WriteByte(f[0])
			f = f[1:]
		}

		popi++
		switch expecting {
		case typ.Bool:
			fmt.Fprint(p, env.Bool(popi))
		case typ.String:
			if pop := env.Get(popi); pop.IsBytes() {
				fmt.Fprintf(p, tmp.String(), pop.Array().Unwrap())
			} else {
				fmt.Fprintf(p, tmp.String(), pop.String())
			}
		case typ.Number + typ.String:
			if pop := env.Get(popi); pop.Type() == typ.String {
				fmt.Fprintf(p, tmp.String(), pop.Str())
				continue
			} else if pop.IsBytes() {
				fmt.Fprintf(p, tmp.String(), pop.Array().Unwrap())
				continue
			}
			fallthrough
		case typ.Number:
			if pop := env.Num(popi); pop.IsInt64() {
				fmt.Fprintf(p, tmp.String(), pop.Int64())
			} else {
				fmt.Fprintf(p, tmp.String(), pop.Float64())
			}
		case typ.Native:
			fmt.Fprint(p, env.Interface(popi))
		}
	}
}

func fileInfo(fi os.FileInfo) *Object {
	return NewObject(0).
		SetProp("filename", Str(fi.Name())).
		SetProp("size", Int64(fi.Size())).
		SetProp("mode", Int64(int64(fi.Mode()))).
		SetProp("modestr", Str(fi.Mode().String())).
		SetProp("modtime", ValueOf(fi.ModTime())).
		SetProp("isdir", Bool(fi.IsDir()))
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

	work := func(fun *Object, outError *error, p payload) {
		if p.i == -1 {
			res, err := e.Call2(fun, p.k, *p.v)
			if err != nil {
				*outError = err
			} else {
				*p.v = res
			}
		} else {
			res, err := e.Call2(fun, Int(p.i), p.k)
			if err != nil {
				*outError = err
			} else {
				t.Array().Set(p.i, res)
			}
		}
	}

	var outError error
	if n == 1 {
		if t.Type() == typ.Array {
			t.Array().Foreach(func(k int, v Value) bool { work(fun, &outError, payload{k, v, nil}); return outError == nil })
		} else {
			t.Object().Foreach(func(k Value, v *Value) bool { work(fun, &outError, payload{-1, k, v}); return outError == nil })
		}
	} else {
		var in = make(chan payload, t.Len())
		var wg sync.WaitGroup
		wg.Add(n)
		for i := 0; i < n; i++ {
			go func() {
				defer wg.Done()
				for p := range in {
					if outError != nil {
						return
					}
					work(fun, &outError, p)
				}
			}()
		}

		if t.Type() == typ.Array {
			t.Array().Foreach(func(k int, v Value) bool { in <- payload{k, v, nil}; return true })
		} else {
			t.Object().Foreach(func(k Value, v *Value) bool { in <- payload{-1, k, v}; return true })
		}
		close(in)

		wg.Wait()
	}
	internal.PanicErr(outError)
	return t
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
