package nj

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

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

func closeBuffer(p *bytes.Buffer, suffix string) {
	for p.Len() > 0 {
		b := p.Bytes()[p.Len()-1]
		if b == ' ' || b == ',' {
			p.Truncate(p.Len() - 1)
		} else {
			break
		}
	}
	p.WriteString(suffix)
}

func ifstr(v bool, t, f string) string {
	if v {
		return t
	}
	return f
}

func ifquote(v bool, s string) string {
	if v {
		return strconv.Quote(s)
	}
	return s
}

func or(a, b interface{}) interface{} {
	if a != nil {
		return a
	}
	return b
}

func renameFuncName(k, v Value) Value {
	if v.IsObject() {
		if cls := v.Object().Callable; cls != nil && cls.Name == internal.UnnamedFunc {
			cls.Name = k.String()
		}
	}
	return v
}

func setObjectRecv(v, r Value) Value {
	if v.IsObject() {
		v.Object().this = r
	}
	return v
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
	case typ.Array:
		if a := v.Array().any; a != nil {
			if _, ok := a.([]byte); ok {
				return "bytes"
			}
			return reflect.TypeOf(a).String()
		}
		fallthrough
	default:
		return vt.String()
	}
}

func getEncB64(enc *base64.Encoding, padding rune) *base64.Encoding {
	if padding != '=' {
		enc = enc.WithPadding(padding)
	}
	return enc
}

func getEncB32(enc *base32.Encoding, padding rune) *base32.Encoding {
	if padding != '=' {
		enc = enc.WithPadding(padding)
	}
	return enc
}

func mathMinMax(e *Env, max bool) {
	if v := e.Num(0); v.IsInt64() {
		vi := v.Int64()
		for ii := 1; ii < len(e.Stack()); ii++ {
			if x := e.Int64(ii); x >= vi == max {
				vi = x
			}
		}
		e.A = Int64(vi)
	} else {
		vf := v.Float64()
		for i := 1; i < len(e.Stack()); i++ {
			if x := e.Float64(i); x >= vf == max {
				vf = x
			}
		}
		e.A = Float64(vf)
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
