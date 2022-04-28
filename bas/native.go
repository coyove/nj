package bas

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sync"
	"unsafe"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

type NativeMeta struct {
	Name         string
	Proto        *Object
	Len          func(*Native) int
	Size         func(*Native) int
	Clear        func(*Native)
	Values       func(*Native) []Value
	Get          func(*Native, int) Value
	Set          func(*Native, int, Value)
	GetKey       func(*Native, Value) Value
	SetKey       func(*Native, Value, Value)
	Append       func(*Native, ...Value)
	Slice        func(*Native, int, int) *Native
	SliceInplace func(*Native, int, int)
	Copy         func(*Native, int, int, *Native)
	Concat       func(*Native, *Native)
	Marshal      func(*Native, io.Writer, typ.MarshalType)
}

var (
	internalArrayMeta = &NativeMeta{}
	bytesArrayMeta    = &NativeMeta{}
	stringsArrayMeta  = &NativeMeta{}
	errorArrayMeta    = &NativeMeta{}
	genericMetaCache  sync.Map
)

func init() {
	*internalArrayMeta = NativeMeta{
		"internal",
		Proto.Array,
		func(a *Native) int { return len(a.internal) },
		func(a *Native) int { return cap(a.internal) },
		func(a *Native) { a.internal = a.internal[:0] },
		func(a *Native) []Value { return a.internal },
		func(a *Native, idx int) Value { return a.internal[idx] },
		func(a *Native, idx int, v Value) { a.internal[idx] = v },
		sgGetKey,
		sgSetKeyNotSupported,
		func(a *Native, v ...Value) { a.internal = append(a.internal, v...) },
		func(a *Native, s, e int) *Native { return &Native{meta: a.meta, internal: a.internal[s:e]} },
		func(a *Native, s, e int) { a.internal = a.internal[s:e] },
		func(a *Native, s, e int, from *Native) {
			if from.meta != a.meta {
				for i := s; i < e; i++ {
					a.internal[i] = from.Get(i - s)
				}
			} else {
				copy(a.internal[s:e], from.internal)
			}
		},
		func(a *Native, b *Native) {
			if a.meta != b.meta {
				for i := 0; i < b.Len(); i++ {
					a.internal = append(a.internal, b.Get(i))
				}
			} else {
				a.internal = append(a.internal, b.internal...)
			}
		},
		func(a *Native, w io.Writer, mt typ.MarshalType) {
			w.Write([]byte("["))
			a.Foreach(func(i int, v Value) bool {
				w.Write([]byte(internal.IfStr(i == 0, "", ",")))
				v.Stringify(w, mt)
				return true
			})
			w.Write([]byte("]"))
		},
	}
	*bytesArrayMeta = NativeMeta{
		"bytes",
		Proto.Bytes,
		func(a *Native) int { return len((a.any).([]byte)) },
		func(a *Native) int { return cap((a.any).([]byte)) },
		func(a *Native) { a.any = a.any.([]byte)[:0] },
		func(a *Native) []Value { a.notSupported("Values"); return nil },
		func(a *Native, idx int) Value { return Int64(int64(a.any.([]byte)[idx])) },
		func(a *Native, idx int, v Value) {
			a.any.([]byte)[idx] = byte(v.AssertType(typ.Number, "bytes.Set").Int())
		},
		sgGetKey,
		sgSetKeyNotSupported,
		func(a *Native, v ...Value) {
			p := a.any.([]byte)
			for _, b := range v {
				p = append(p, byte(b.AssertType(typ.Number, "bytes.Append").Int()))
			}
			a.any = p
		},
		func(a *Native, start, end int) *Native {
			return &Native{meta: a.meta, any: a.any.([]byte)[start:end]}
		},
		func(a *Native, start, end int) {
			a.any = a.any.([]byte)[start:end]
		},
		func(a *Native, start, end int, from *Native) {
			if from.meta == internalArrayMeta {
				buf := a.any.([]byte)
				for i := start; i < end; i++ {
					buf[i] = byte(from.Get(i-start).AssertType(typ.Number, "bytes.Copy").Int())
				}
			} else {
				copy(a.any.([]byte)[start:end], from.any.([]byte))
			}
		},
		func(a *Native, b *Native) {
			if b.meta == internalArrayMeta {
				buf := a.any.([]byte)
				for i := 0; i < b.Len(); i++ {
					buf[i] = byte(b.Get(i).AssertType(typ.Number, "bytes.Concat").Int())
				}
				a.any = buf
			} else {
				a.any = append(a.any.([]byte), b.any.([]byte)...)
			}
		},
		func(a *Native, w io.Writer, mt typ.MarshalType) {
			if mt != typ.MarshalToJSON {
				sgMarshal(a, w, mt)
			} else {
				enc := base64.NewEncoder(base64.StdEncoding, w)
				enc.Write(a.any.([]byte))
				enc.Close()
			}
		},
	}
	*stringsArrayMeta = NativeMeta{
		"[]string",
		Proto.Array,
		func(a *Native) int { return len((a.any).([]string)) },
		func(a *Native) int { return cap((a.any).([]string)) },
		func(a *Native) { a.any = a.any.([]byte)[:0] },
		func(a *Native) []Value {
			res := make([]Value, a.Len())
			for i := 0; i < a.Len(); i++ {
				res[i] = a.Get(i)
			}
			return res
		},
		func(a *Native, idx int) Value { return Str(a.any.([]string)[idx]) },
		func(a *Native, idx int, v Value) {
			a.any.([]string)[idx] = v.AssertType(typ.String, "[]string.Set").Str()
		},
		sgGetKey,
		sgSetKeyNotSupported,
		func(a *Native, v ...Value) {
			p := a.any.([]string)
			for _, b := range v {
				p = append(p, b.AssertType(typ.String, "[]string.Append").Str())
			}
			a.any = p
		},
		func(a *Native, start, end int) *Native {
			p := a.any.([]string)[start:end]
			return &Native{meta: a.meta, any: p}
		},
		func(a *Native, start, end int) {
			a.any = a.any.([]string)[start:end]
		},
		func(a *Native, start, end int, from *Native) {
			if from.meta == internalArrayMeta {
				buf := a.any.([]string)
				for i := start; i < end; i++ {
					buf[i] = from.Get(i-start).AssertType(typ.String, "[]string.Copy").Str()
				}
			} else {
				copy(a.any.([]byte)[start:end], from.any.([]byte))
			}
		},
		func(a *Native, b *Native) {
			if b.meta == internalArrayMeta {
				buf := a.any.([]string)
				for i := 0; i < b.Len(); i++ {
					buf[i] = b.Get(i).AssertType(typ.String, "[]string.Concat").Str()
				}
				a.any = buf
			} else {
				a.any = append(a.any.([]byte), b.any.([]byte)...)
			}
		},
		sgMarshal,
	}
	*errorArrayMeta = NativeMeta{
		"error",
		Proto.Error,
		func(a *Native) int { return 1 },
		func(a *Native) int { return 1 },
		sgClearNotSupported,
		sgValuesNotSupported,
		func(a *Native, idx int) Value { return a.ToValue() },
		sgSetNotSupported,
		sgGetKey,
		sgSetKeyNotSupported,
		sgAppendNotSupported,
		sgSliceNotSupported,
		sgSliceInplaceNotSupported,
		sgCopyNotSupported,
		sgConcatNotSupported,
		func(a *Native, w io.Writer, mt typ.MarshalType) {
			w.Write([]byte(internal.IfQuote(mt == typ.MarshalToJSON, a.any.(*ExecError).Error())))
		},
	}
}

func GetNativeMeta(v interface{}) *NativeMeta {
	switch v.(type) {
	case []Value:
		return internalArrayMeta
	case []byte:
		return bytesArrayMeta
	case []string:
		return stringsArrayMeta
	}
	rt := reflect.TypeOf(v)
	if v, ok := genericMetaCache.Load(rt); ok {
		return v.(*NativeMeta)
	}
	var a *NativeMeta
	switch rt.Kind() {
	default:
		a = &NativeMeta{rt.String(), Proto.Native,
			sgLenNotSupported,
			sgSizeNotSupported,
			sgClearNotSupported,
			sgValuesNotSupported,
			func(a *Native, idx int) Value { return a.ToValue() },
			sgSetNotSupported,
			func(a *Native, k Value) Value {
				if v, ok := reflectLoad(a.any, k); ok {
					return v
				}
				return sgGetKey(a, k)
			},
			func(a *Native, k, v Value) { reflectStore(a.any, k, v) },
			sgAppendNotSupported,
			sgSliceNotSupported,
			sgSliceInplaceNotSupported,
			sgCopyNotSupported,
			sgConcatNotSupported,
			sgMarshalTypeOnly,
		}
		if rt.Kind() == reflect.Map {
			a.Proto = Proto.NativeMap
			a.Len = sgLen
			a.Size = sgSize
			a.Marshal = sgMarshal
		}
	case reflect.Chan:
		a = &NativeMeta{rt.String(), Proto.Channel,
			sgLen,
			sgSize,
			sgClearNotSupported,
			sgValuesNotSupported,
			func(a *Native, idx int) Value { return a.ToValue() },
			sgSetNotSupported,
			sgGetKey,
			sgSetKeyNotSupported,
			sgAppendNotSupported,
			sgSliceNotSupported,
			sgSliceInplaceNotSupported,
			sgCopyNotSupported,
			sgConcatNotSupported,
			sgMarshalTypeOnly,
		}
	case reflect.Array, reflect.Slice:
		a = &NativeMeta{rt.String(), Proto.Array,
			sgLen,
			sgSize,
			sgClear,
			sgValues,
			sgGet,
			sgSet,
			sgGetKey,
			sgSetKeyNotSupported,
			sgAppend,
			sgSlice,
			sgSliceInplace,
			sgCopy,
			sgConcat,
			sgMarshal,
		}
		if rt.Kind() == reflect.Array {
			a.SliceInplace = sgSliceInplaceNotSupported
			a.Clear = sgClearNotSupported
			a.Append = sgAppendNotSupported
			a.Copy = sgCopyNotSupported
			a.Concat = sgConcatNotSupported
		}
	}
	genericMetaCache.Store(rt, a)
	return a
}

type Native struct {
	meta     *NativeMeta
	internal []Value
	any      interface{}
}

// NewArray creates an array consists of given arguments
func NewArray(m ...Value) *Native {
	return &Native{meta: internalArrayMeta, internal: m}
}

func NewNative(any interface{}) *Native {
	return newTypedArray(any, GetNativeMeta(any))
}

func newTypedArray(any interface{}, meta *NativeMeta) *Native {
	return &Native{meta: meta, any: any}
}

// Error creates a builtin error, env can be nil
func Error(e *Env, err error) Value {
	if err == nil {
		return Nil
	} else if _, ok := err.(*ExecError); ok {
		return newTypedArray(err, errorArrayMeta).ToValue()
	}
	ee := &ExecError{root: err}
	if e != nil {
		ee.stacks = e.Runtime().Stacktrace()
	}
	return newTypedArray(ee, errorArrayMeta).ToValue()
}

func (a *Native) ToValue() Value {
	return Value{v: uint64(typ.Native), p: unsafe.Pointer(a)}
}

func (a *Native) Unwrap() interface{} {
	if a.meta == internalArrayMeta {
		return a.internal
	}
	return a.any
}

func (a *Native) Len() int {
	if a.meta == internalArrayMeta {
		return len(a.internal)
	}
	return a.meta.Len(a)
}

func (a *Native) Size() int {
	if a.meta == internalArrayMeta {
		return cap(a.internal)
	}
	return a.meta.Size(a)
}

func (a *Native) Values() []Value {
	if a.meta == internalArrayMeta {
		return a.internal
	}
	return a.meta.Values(a)
}

func (a *Native) Get(v int) Value {
	if a.meta == internalArrayMeta {
		return a.internal[v]
	}
	return a.meta.Get(a, v)
}

func (a *Native) Set(idx int, v Value) {
	if a.meta == internalArrayMeta {
		a.internal[idx] = v
	} else {
		a.meta.Set(a, idx, v)
	}
}

func (a *Native) GetKey(k Value) Value {
	return a.meta.GetKey(a, k)
}

func (a *Native) SetKey(k, v Value) {
	a.meta.SetKey(a, k, v)
}

func (a *Native) Append(v ...Value) {
	a.meta.Append(a, v...)
}

func (a *Native) Slice(start, end int) *Native {
	return a.meta.Slice(a, start, end)
}

func (a *Native) SliceInplace(start, end int) {
	a.meta.SliceInplace(a, start, end)
}

func (a *Native) Clear() {
	if a.meta == internalArrayMeta {
		a.internal = a.internal[:0]
	} else {
		a.meta.Clear(a)
	}
}

func (a *Native) Copy(start, end int, from *Native) {
	if a.meta == internalArrayMeta || from.meta == internalArrayMeta {
	} else if a.meta != from.meta {
		internal.Panic("copy array with different types: from %q to %q", from.meta.Name, a.meta.Name)
	}
	a.meta.Copy(a, start, end, from)
}

func (a *Native) Concat(b *Native) {
	if a.meta == internalArrayMeta || b.meta == internalArrayMeta {
	} else if a.meta != b.meta {
		internal.Panic("concat array with different types: from %q to %q", b.meta.Name, a.meta.Name)
	}
	a.meta.Concat(a, b)
}

func (a *Native) Marshal(w io.Writer, mt typ.MarshalType) {
	a.meta.Marshal(a, w, mt)
}

func (a *Native) Foreach(f func(k int, v Value) bool) {
	for i := 0; i < a.Len(); i++ {
		if !f(i, a.Get(i)) {
			break
		}
	}
}

func (a *Native) Typed() bool {
	return a.meta != internalArrayMeta
}

func (a *Native) notSupported(method string) {
	panic(a.meta.Name + "." + method + " not allowed")
}

func sgLen(a *Native) int {
	return reflect.ValueOf(a.any).Len()
}

func sgSize(a *Native) int {
	return reflect.ValueOf(a.any).Cap()
}

func sgClear(a *Native) {
	a.any = reflect.ValueOf(a.any).Slice(0, 0).Interface()
}

func sgValues(a *Native) []Value {
	res := make([]Value, a.Len())
	for i := 0; i < len(res); i++ {
		res[i] = a.Get(i)
	}
	return res
}

func sgGet(a *Native, idx int) Value {
	return ValueOf(reflect.ValueOf(a.any).Index(idx).Interface())
}

func sgSet(a *Native, idx int, v Value) {
	rv := reflect.ValueOf(a.any)
	rv.Index(idx).Set(ToType(v, rv.Type().Elem()))
}

func sgAppend(a *Native, v ...Value) {
	rv := reflect.ValueOf(a.any)
	rt := rv.Type().Elem()
	for _, b := range v {
		rv = reflect.Append(rv, ToType(b, rt))
	}
	a.any = rv.Interface()
}

func sgSlice(a *Native, start, end int) *Native {
	return &Native{meta: a.meta, any: reflect.ValueOf(a.any).Slice(start, end).Interface()}
}

func sgSliceInplace(a *Native, start, end int) {
	a.any = reflect.ValueOf(a.any).Slice(start, end).Interface()
}

func sgCopy(a *Native, start, end int, from *Native) {
	if from.meta == internalArrayMeta {
		rv := reflect.ValueOf(a.any)
		rt := rv.Type().Elem()
		for i := start; i < end; i++ {
			rv.Index(i).Set(ToType(from.Get(i-start), rt))
		}
	} else {
		reflect.Copy(reflect.ValueOf(a.any).Slice(start, end), reflect.ValueOf(from.any))
	}
}

func sgConcat(a *Native, b *Native) {
	if b.meta == internalArrayMeta {
		rv := reflect.ValueOf(a.any)
		rt := rv.Type().Elem()
		for i := 0; i < b.Len(); i++ {
			rv = reflect.Append(rv, ToType(b.Get(i), rt))
		}
		a.any = rv.Interface()
	} else {
		a.any = reflect.AppendSlice(reflect.ValueOf(a.any), reflect.ValueOf(b.any)).Interface()
	}
}

func sgMarshal(a *Native, w io.Writer, mt typ.MarshalType) {
	if mt != typ.MarshalToJSON {
		fmt.Fprint(w, a.any)
	} else {
		json.NewEncoder(w).Encode(a.any)
	}
}

func sgLenNotSupported(a *Native) int {
	a.notSupported("Len")
	return 0
}

func sgSizeNotSupported(a *Native) int {
	a.notSupported("Size")
	return 0
}

func sgSliceNotSupported(a *Native, start int, end int) *Native {
	a.notSupported("Slice")
	return nil
}

func sgSliceInplaceNotSupported(a *Native, start int, end int) {
	a.notSupported("SliceInplace")
}

func sgClearNotSupported(a *Native) {
	a.notSupported("Clear")
}

func sgValuesNotSupported(a *Native) []Value {
	a.notSupported("Values")
	return nil
}

func sgAppendNotSupported(a *Native, v ...Value) {
	a.notSupported("Append")
}

func sgCopyNotSupported(a *Native, start, end int, from *Native) {
	a.notSupported("Copy")
}

func sgConcatNotSupported(a *Native, b *Native) {
	a.notSupported("Concat")
}

func sgSetNotSupported(a *Native, b int, c Value) {
	a.notSupported("Set")
}

func sgSetKeyNotSupported(a *Native, b, c Value) {
	a.notSupported("SetKey")
}

func sgGetKey(a *Native, k Value) Value {
	f := a.meta.Proto.Find(k)
	if f != Nil {
		f = setObjectRecv(f, a.ToValue())
	}
	return f
}

func sgMarshalTypeOnly(a *Native, w io.Writer, mt typ.MarshalType) {
	if mt != typ.MarshalToJSON {
		fmt.Fprint(w, reflectString(a.any))
	} else {
		json.NewEncoder(w).Encode(reflectString(a.any))
	}
}
