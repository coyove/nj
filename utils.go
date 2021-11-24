package nj

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
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
			return Val(v.Interface())
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
	return Val(f.Interface())
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
		if cls := v.Object().callable; cls != nil && cls.Name == internal.UnnamedFunc {
			cls.Name = k.String()
		}
	}
	return v
}

func setObjectRecv(v, r Value) Value {
	if v.IsObject() {
		v.Object().receiver = r
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
	default:
		return vt.String()
	}
}
