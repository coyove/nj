package script

import (
	"fmt"
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
		v := rv.MapIndex(reflect.ValueOf(key.TypedInterface(rv.Type().Elem())))
		if v.IsValid() {
			return Interface(v.Interface())
		}
	case reflect.Slice, reflect.Array:
		idx := key.ExpectMsg(VNumber, "loadarray").Int() - 1
		if idx < int64(rv.Len()) && idx >= 0 {
			return Interface(rv.Index(int(idx)).Interface())
		}
	}

	k := camelKey(key.ExpectMsg(VString, "loadstruct")._str())
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
		start := int(start - 1)
		end := int(end - 1 + 1)
		if start >= 0 && start < rv.Len() && end >= 0 && end < rv.Len() && start <= end {
			return rv.Slice(start, end).Interface()
		}
	}
	return nil
}

func reflectStore(v interface{}, key Value, v2 Value) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map:
		// 	rk := reflect.ValueOf(key.TypedInterface(rv.Type().Elem()))
		// 	v := rv.MapIndex(rk)
		// 	if !v.IsValid() {
		// 		panicf("store: readonly map")
		// 	}
		// 	if v2.IsNil() {
		// 		rv.SetMapIndex(rk, reflect.Value{})
		// 	} else {
		// 		// panicf("store: readonly map")
		// 		rv.SetMapIndex(rk, reflect.ValueOf(v2.TypedInterface(rv.Type().Elem())))
		// 	}
		return
	case reflect.Slice, reflect.Array:
		// 	idx := key.ExpectMsg(VNumber, "storearray").Int() - 1
		// 	if idx >= int64(rv.Len()) || idx < 0 {
		// 		return
		// 	}
		// 	rv.Index(int(idx)).Set(reflect.ValueOf(v2.TypedInterface(rv.Type().Elem())))
		return
	}

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	k := camelKey(key.ExpectMsg(VString, "storestruct")._str())
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
// If there are no cycles, nil is returned.
// If one or more cycles exist, an error will be returned. This error will contain a path to the first cycle found.
var reflectCheckCyclicStruct = func() func(v interface{}) error {

	checkParents := func(value reflect.Value, parents []uintptr, names []string) ([]uintptr, error) {
		kind := value.Kind()
		if kind == reflect.Map || kind == reflect.Ptr || kind == reflect.Slice {
			address := value.Pointer()
			for _, parent := range parents {
				if parent == address {
					return nil, fmt.Errorf("cycle found: %v", names)
				}
			}
			return append(parents, address), nil
		}
		return parents, nil
	}

	var doCheck func(reflect.Value, []uintptr, []string) error
	doCheck = func(value reflect.Value, parents []uintptr, names []string) error {
		kind := value.Kind()

		if kind == reflect.Interface {
			value = value.Elem()
			kind = value.Kind()
		}

		newParents, err := checkParents(value, parents, names)
		if err != nil {
			return err
		}

		if kind == reflect.Map {
			for _, key := range value.MapKeys() {
				err := doCheck(value.MapIndex(key), newParents, append(names, key.String()))
				if err != nil {
					return err
				}
			}
		}

		if kind == reflect.Ptr {
			return doCheck(value.Elem(), newParents, names)
		}

		if kind == reflect.Slice {
			for i := 0; i < value.Len(); i++ {
				err := doCheck(value.Index(i), newParents, append(names, fmt.Sprintf("[%d]", i)))
				if err != nil {
					return err
				}
			}
		}

		if kind == reflect.Struct {
			for i := 0; i < value.NumField(); i++ {
				t := value.Type()
				fieldType := t.Field(i)
				err := doCheck(value.Field(i), newParents, append(names, fieldType.Name))
				if err != nil {
					return err
				}
			}
		}

		return nil
	}

	return func(v interface{}) error {
		return doCheck(reflect.ValueOf(v), nil, nil)
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
