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
	"unsafe"

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

func reflectLoad(v interface{}, key Value) (Value, bool) {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("reflectLoad %T[%v]: %v", v, key, r))
		}
	}()

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Map {
		v := rv.MapIndex(ToType(key, rv.Type().Key()))
		if v.IsValid() {
			return ValueOf(v.Interface()), true
		}
		return Nil, false
	}

	k := key.AssertType(typ.String, "").Str()
	f := rv.MethodByName(k)
	if !f.IsValid() {
		if rv.Kind() == reflect.Ptr {
			f = rv.Elem().MethodByName(k)
		}
	}
	if !f.IsValid() {
		if strings.HasPrefix(k, "p") {
			if rv.Kind() != reflect.Ptr {
				return Nil, false
			}
			t, ok := rv.Elem().Type().FieldByName(k[1:])
			if !ok {
				return Nil, false
			}
			ptr := (*struct{ a, b uintptr })(unsafe.Pointer(&v)).b + t.Offset
			f = reflect.NewAt(t.Type, unsafe.Pointer(ptr))
		} else {
			f = reflect.Indirect(rv).FieldByName(k)
			if !f.IsValid() {
				return Nil, false
			}
		}
	}
	return ValueOf(f.Interface()), true
}

func reflectStore(subject interface{}, key, value Value) {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("reflectStore %T[%v]: %v", subject, key, r))
		}
	}()

	rv := reflect.ValueOf(subject)
	if rv.Kind() == reflect.Map {
		rv.SetMapIndex(ToType(key, rv.Type().Key()), ToType(value, rv.Type().Elem()))
		return
	}

	rv = reflect.Indirect(rv)
	k := key.AssertType(typ.String, "").Str()
	f := rv.FieldByName(k)
	if !f.IsValid() || !f.CanAddr() {
		internal.Panic("%q not assignable in %T", k, subject)
	}
	f.Set(ToType(value, f.Type()))
}

func reflectString(i interface{}) string {
	if s, ok := i.(fmt.Stringer); ok {
		return s.String()
	}
	if s, ok := i.(error); ok {
		return s.Error()
	}
	return "<" + reflect.TypeOf(i).String() + ">"
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
	case typ.Number, typ.Bool:
		return v.JSONString()
	case typ.String:
		if Len(v) <= 32 {
			return v.JSONString()
		}
		return strconv.Quote(v.Str()[:32] + "...")
	case typ.Object:
		if v.Object().fun != nil { // including named objects
			return v.Object().fun.String()
		}
		return "{" + v.Object().Name() + "}"
	case typ.Native:
		a := v.Native()
		if a.IsInternalArray() {
			return fmt.Sprintf("array(%d)", a.Len())
		}
		return fmt.Sprintf("native(%s)", a.meta.Name)
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
				expecting = typ.Object
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
			if pop := env.Get(popi); IsBytes(pop) {
				fmt.Fprintf(p, tmp.String(), pop.Native().Unwrap())
			} else {
				fmt.Fprintf(p, tmp.String(), pop.String())
			}
		case typ.Number + typ.String:
			if pop := env.Get(popi); pop.Type() == typ.String {
				fmt.Fprintf(p, tmp.String(), pop.Str())
				continue
			} else if IsBytes(pop) {
				fmt.Fprintf(p, tmp.String(), pop.Native().Unwrap())
				continue
			}
			fallthrough
		case typ.Number:
			if pop := env.Num(popi); pop.IsInt64() {
				fmt.Fprintf(p, tmp.String(), pop.Int64())
			} else {
				fmt.Fprintf(p, tmp.String(), pop.Float64())
			}
		case typ.Object:
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
				t.Native().Set(p.i, res)
			}
		}
	}

	var outError error
	if n == 1 {
		if t.Type() == typ.Native {
			t.Native().Foreach(func(k int, v Value) bool { work(fun, &outError, payload{k, v, nil}); return outError == nil })
		} else {
			t.Object().Foreach(func(k Value, v *Value) bool { work(fun, &outError, payload{-1, k, v}); return outError == nil })
		}
	} else {
		var in = make(chan payload, Len(t))
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

		if t.Type() == typ.Native {
			t.Native().Foreach(func(k int, v Value) bool { in <- payload{k, v, nil}; return true })
		} else {
			t.Object().Foreach(func(k Value, v *Value) bool { in <- payload{-1, k, v}; return true })
		}
		close(in)

		wg.Wait()
	}
	internal.PanicErr(outError)
	return t
}
