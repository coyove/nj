package nj

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

func reflectLoad(v interface{}, key Value) Value {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map:
		v := rv.MapIndex(key.ReflectValue(rv.Type().Key()))
		if v.IsValid() {
			return Val(v.Interface())
		}
	}

	k := key.MustStr("index key")
	f := rv.MethodByName(k)
	if !f.IsValid() {
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		f = rv.FieldByName(k)
		if !f.IsValid() || !(k[0] >= 'A' && k[0] <= 'Z') {
			return reflectLoadCaseless(v, k)
		}
	}
	return Val(f.Interface())
}

func reflectLoadCaseless(v interface{}, k string) Value {
	rv := reflect.ValueOf(v)
	rt := reflect.TypeOf(v)
	for i := 0; i < rt.NumMethod(); i++ {
		if strings.EqualFold(rt.Method(i).Name, k) {
			return Val(rv.Method(i).Interface())
		}
	}
	if rv.Kind() == reflect.Ptr {
		rv, rt = rv.Elem(), rt.Elem()
	}
	for i := 0; i < rt.NumField(); i++ {
		if strings.EqualFold(rt.Field(i).Name, k) {
			return Val(rv.Field(i).Interface())
		}
	}
	return Nil
}

func reflectStore(subject interface{}, key, value Value) {
	rv := reflect.ValueOf(subject)
	switch rv.Kind() {
	case reflect.Map:
		rk := key.ReflectValue(rv.Type().Key())
		if value == Nil {
			rv.SetMapIndex(rk, reflect.Value{})
		} else {
			rv.SetMapIndex(rk, value.ReflectValue(rv.Type().Elem()))
		}
		return
	}

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	k := key.MustStr("index key")
	f := rv.FieldByName(k)
	if !f.IsValid() || !f.CanAddr() {
		internal.Panic("reflect: %q not assignable in %v", k, subject)
	}
	if f.Type() == reflect.TypeOf(Value{}) {
		f.Set(reflect.ValueOf(value))
	} else {
		f.Set(value.ReflectValue(f.Type()))
	}
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
