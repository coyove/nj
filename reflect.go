package script

import (
	"reflect"
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
		idx := key.MustNum("reflect: load from Go array/slice", 0).Int() - 1
		if idx < int64(rv.Len()) && idx >= 0 {
			return Val(rv.Index(int(idx)).Interface())
		}
	}

	k := (key.MustStr("reflect: load from Go struct", 0))
	f := rv.MethodByName(k)
	if !f.IsValid() {
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		f := rv.FieldByName(k)
		if f.IsValid() {
			return Val(f.Interface())
		}
		// panicf("%q not found in %#v", k, v)
		return Value{}
	}
	return Val(f.Interface())
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
		idx := key.MustNum("reflect: store into Go array/slice", 0).Int()
		if idx >= int64(rv.Len()) || idx < 0 {
			return
		}
		rv.Index(int(idx)).Set(v2.ReflectValue(rv.Type().Elem()))
		return
	}

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	k := (key.MustStr("reflect: store into Go struct", 0))
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
