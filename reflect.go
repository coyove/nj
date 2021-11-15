package script

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func reflectLen(v interface{}) int {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array:
		return rv.Len()
	default:
		panicf("reflect: can't measure length of %T", v)
	}
	return -1
}

func reflectLoad(v interface{}, key Value) Value {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map:
		v := rv.MapIndex(key.ReflectValue(rv.Type().Key()))
		if v.IsValid() {
			return Val(v.Interface())
		}
	case reflect.Slice, reflect.Array:
		idx := key.MustNum("index key").Int()
		if idx < int64(rv.Len()) && idx >= 0 {
			return Val(rv.Index(int(idx)).Interface())
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

func reflectStore(v interface{}, key Value, v2 Value) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map:
		rk := key.ReflectValue(rv.Type().Key())
		if v2 == Nil {
			rv.SetMapIndex(rk, reflect.Value{})
		} else {
			rv.SetMapIndex(rk, v2.ReflectValue(rv.Type().Elem()))
		}
		return
	case reflect.Slice, reflect.Array:
		idx := key.MustNum("index key").Int()
		if idx >= int64(rv.Len()) || idx < 0 {
			return
		}
		rv.Index(int(idx)).Set(v2.ReflectValue(rv.Type().Elem()))
		return
	}

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	k := key.MustStr("index key")
	f := rv.FieldByName(k)
	if !f.IsValid() || !f.CanAddr() {
		panicf("reflect: %q not assignable in %v", k, v)
	}
	if f.Type() == reflect.TypeOf(Value{}) {
		f.Set(reflect.ValueOf(v2))
	} else {
		f.Set(v2.ReflectValue(f.Type()))
	}
}

func Stringify(v interface{}) string {
	if v == nil {
		return "nil"
	}
	p := bytes.NewBufferString("")
	reflectStringify(p, reflect.ValueOf(v), 0, false)
	return p.String()
}

func reflectStringify(p *bytes.Buffer, rv reflect.Value, depth int, showType bool) {
	const maxDepth = 10
	if rv.Type().Implements(reflect.TypeOf((*fmt.Stringer)(nil)).Elem()) {
		if !rv.CanInterface() {
			p.WriteString("...")
		} else {
			p.WriteString(rv.Interface().(fmt.Stringer).String())
		}
		return
	}
	if rv.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		if !rv.CanInterface() {
			p.WriteString("...")
		} else {
			p.WriteString(rv.Interface().(error).Error())
		}
		return
	}

	k := rv.Kind()

	if showType {
		if k == reflect.Ptr {
			rv = rv.Elem()
			p.WriteString("(*")
		} else if k == reflect.Interface {
			rv = rv.Elem()
			p.WriteString("(^")
		}
		if !rv.IsValid() {
			p.WriteString("nil)")
			return
		}
		p.WriteString(rv.Type().String())
		p.WriteString(")(")
	} else {
		if k == reflect.Ptr || k == reflect.Interface {
			rv = rv.Elem()
		}
		if !rv.IsValid() {
			p.WriteString("nil)")
			return
		}
	}

	switch rv.Kind() {
	case reflect.Bool:
		p.WriteString(strconv.FormatBool(rv.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p.WriteString(strconv.FormatInt(rv.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p.WriteString(strconv.FormatUint(rv.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		p.WriteString(strconv.FormatFloat(rv.Float(), 'f', -1, 64))
	case reflect.Complex64, reflect.Complex128:
		c := rv.Complex()
		p.WriteString(strconv.FormatFloat(real(c), 'f', -1, 64))
		p.WriteString("+")
		p.WriteString(strconv.FormatFloat(imag(c), 'f', -1, 64))
		p.WriteString("i")
	case reflect.String:
		p.WriteString(rv.String())
	case reflect.Struct:
		count := 0
		p.WriteString(rv.Type().String())
		p.WriteString("{ ")
		if depth > maxDepth {
			p.WriteString("...")
		} else {
			rt := rv.Type()
			for i := 0; i < rt.NumField(); i++ {
				p.WriteString(rt.Field(i).Name)
				p.WriteString(": ")
				reflectStringify(p, rv.Field(i), depth+1, rt.Field(i).Type.Kind() == reflect.Interface)
				p.WriteString(", ")
				count++
			}
		}
		closeBuffer(p, " }")
	case reflect.Array, reflect.Slice:
		p.WriteString("[")
		if rv.Kind() == reflect.Array {
			p.WriteString(strconv.Itoa(int(rv.Len())))
		}
		p.WriteString("]")
		p.WriteString(rv.Type().Elem().String())
		p.WriteString("{ ")
		if depth > maxDepth {
			p.WriteString("...")
		} else {
			showType := rv.Type().Elem().Kind() == reflect.Interface
			for i := 0; i < rv.Len(); i++ {
				reflectStringify(p, rv.Index(i), depth+1, showType)
				p.WriteString(", ")
			}
		}
		closeBuffer(p, " }")
	case reflect.Map:
		p.WriteString("map[")
		p.WriteString(rv.Type().Key().String())
		p.WriteString("]")
		p.WriteString(rv.Type().Elem().String())
		p.WriteString("{ ")
		if depth > maxDepth {
			p.WriteString("...")
		} else {
			iter := rv.MapRange()
			showType1 := rv.Type().Key().Kind() == reflect.Interface
			showType2 := rv.Type().Elem().Kind() == reflect.Interface
			for iter.Next() {
				reflectStringify(p, iter.Key(), depth+1, showType1)
				p.WriteString(": ")
				reflectStringify(p, iter.Value(), depth+1, showType2)
				p.WriteString(", ")
			}
		}
		closeBuffer(p, " }")
	default:
		p.WriteString(rv.Type().String())
		p.WriteString("{}")
	}

	if showType {
		p.WriteString(")")
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
