package script

import (
	"reflect"
)

func reflectLen(v interface{}) int {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array:
		return rv.Len()
	}
	return -1
}

func reflectLoad(v interface{}, key Value) Value {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map:
		v := rv.MapIndex(reflect.ValueOf(key.TypedInterface(rv.Type().Key())))
		if v.IsValid() {
			return Interface(v.Interface())
		}
	case reflect.Slice, reflect.Array:
		idx := key.MustBe(VNumber, "load array", 0).Int() - 1
		if idx < int64(rv.Len()) && idx >= 0 {
			return Interface(rv.Index(int(idx)).Interface())
		}
	}

	k := camelKey(key.MustBe(VString, "load struct", 0)._str())
	f := rv.MethodByName(k)
	if !f.IsValid() {
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		f := rv.FieldByName(k)
		if f.IsValid() {
			return Interface(f.Interface())
		}
		// panicf("%q not found in %#v", k, v)
		return Value{}
	}
	return Interface(f.Interface())
}

func reflectSlice(v interface{}, start, end int64) interface{} {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array, reflect.String:
		return rv.Slice(sliceInRange(start, end, rv.Len())).Interface()
	}
	return nil
}

func reflectStore(v interface{}, key Value, v2 Value) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map:
		rk := reflect.ValueOf(key.TypedInterface(rv.Type().Key()))
		v := rv.MapIndex(rk)
		if !v.IsValid() {
			panicf("store: readonly map")
		}
		if v2.IsNil() {
			rv.SetMapIndex(rk, reflect.Value{})
		} else {
			// panicf("store: readonly map")
			rv.SetMapIndex(rk, reflect.ValueOf(v2.TypedInterface(rv.Type().Elem())))
		}
		return
	case reflect.Slice, reflect.Array:
		panicf("store: readonly slice")
		idx := key.MustBe(VNumber, "store array", 0).Int() - 1
		if idx >= int64(rv.Len()) || idx < 0 {
			return
		}
		rv.Index(int(idx)).Set(reflect.ValueOf(v2.TypedInterface(rv.Type().Elem())))
		return
	}

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	k := camelKey(key.MustBe(VString, "store struct", 0)._str())
	f := rv.FieldByName(k)
	if !f.IsValid() || !f.CanAddr() {
		panicf("%q not assignable in %#v", k, v)
	}
	if f.Type() == reflect.TypeOf(Value{}) {
		f.Set(reflect.ValueOf(v2))
	} else {
		f.Set(reflect.ValueOf(v2.TypedInterface(f.Type())))
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

func camelKey(k string) string {
	if k == "" {
		return k
	}
	if k[0] >= 'a' && k[0] <= 'z' {
		return string(k[0]-'a'+'A') + k[1:]
	}
	return k
}

func sliceInRange(start, end int64, length int) (int, int) {
	{
		start := int(start - 1)
		end := int(end - 1 + 1)
		if start >= 0 && start <= length && end >= 0 && end <= length && start <= end {
			return start, int(end)
		}
	}
	panicf("slice [%d,%d] overflows [1,%d]", start, end, length)
	return 0, 0
}
