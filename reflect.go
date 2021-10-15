package script

import (
	"reflect"
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
		idx := key.MustNum("index key").Int() - 1
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
		panicf("reflect: %q not assignable in %#v", k, v)
	}
	if f.Type() == reflect.TypeOf(Value{}) {
		f.Set(reflect.ValueOf(v2))
	} else {
		f.Set(v2.ReflectValue(f.Type()))
	}
}

// https://github.com/theothertomelliott/acyclic
// reflectCheckCyclicStruct performs a DFS traversal over an interface to determine if it contains any cycles.
var reflectCheckCyclicStruct = func() func(v interface{}) bool {

	checkParents := func(value reflect.Value, parents []uintptr) ([]uintptr, bool) {
		kind := value.Kind()
		if kind == reflect.Map || kind == reflect.Ptr || kind == reflect.Slice {
			address := value.Pointer()
			for _, parent := range parents {
				if parent == address {
					return nil, false
				}
			}
			return append(parents, address), true
		}
		return parents, true
	}

	var doCheck func(reflect.Value, []uintptr) bool
	doCheck = func(value reflect.Value, parents []uintptr) bool {
		if value.IsZero() || !value.IsValid() {
			return true
		}

		if value.Type() == reflect.TypeOf(Value{}) {
			v := value.Interface().(Value)
			return doCheck(reflect.ValueOf(v.Interface()), parents)
		}

		kind := value.Kind()

		if kind == reflect.Interface {
			value = value.Elem()
			kind = value.Kind()
		}

		newParents, ok := checkParents(value, parents)
		if !ok {
			return ok
		}

		if kind == reflect.Map {
			for _, key := range value.MapKeys() {
				if !doCheck(value.MapIndex(key), newParents) {
					return false
				}
			}
		}

		if kind == reflect.Ptr {
			return doCheck(value.Elem(), newParents)
		}

		if kind == reflect.Slice {
			for i := 0; i < value.Len(); i++ {
				if !doCheck(value.Index(i), newParents) {
					return false
				}
			}
		}

		if kind == reflect.Struct {
			for i := 0; i < value.NumField(); i++ {
				if !doCheck(value.Field(i), newParents) {
					return false
				}
			}
		}

		return true
	}

	return func(v interface{}) bool {
		return doCheck(reflect.ValueOf(v), nil)
	}
}()
