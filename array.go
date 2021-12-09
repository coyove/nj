package nj

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"unsafe"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

type ArrayMeta struct {
	Name         string
	Proto        *Object
	Len          func(*Array) int
	Size         func(*Array) int
	Clear        func(*Array)
	Values       func(*Array) []Value
	Get          func(*Array, int) Value
	Set          func(*Array, int, Value)
	Append       func(*Array, ...Value)
	Slice        func(*Array, int, int) *Array
	SliceInplace func(*Array, int, int)
	Copy         func(*Array, int, int, *Array)
	Concat       func(*Array, *Array)
	Marshal      func(*Array, typ.MarshalType) []byte
}

var (
	internalArrayMeta     = &ArrayMeta{}
	bytesArrayMeta        = &ArrayMeta{}
	stringsArrayMeta      = &ArrayMeta{}
	errorArrayMeta        = &ArrayMeta{}
	genericArrayMetaCache sync.Map
)

func init() {
	*internalArrayMeta = ArrayMeta{
		"internal",
		ArrayProto,
		func(a *Array) int { return len(a.internal) },
		func(a *Array) int { return cap(a.internal) },
		func(a *Array) { a.internal = a.internal[:0] },
		func(a *Array) []Value { return a.internal },
		func(a *Array, idx int) Value { return a.internal[idx] },
		func(a *Array, idx int, v Value) { a.internal[idx] = v },
		func(a *Array, v ...Value) { a.internal = append(a.internal, v...) },
		func(a *Array, s, e int) *Array { return &Array{meta: a.meta, internal: a.internal[s:e]} },
		func(a *Array, s, e int) { a.internal = a.internal[s:e] },
		func(a *Array, s, e int, from *Array) {
			if from.meta != a.meta {
				for i := s; i < e; i++ {
					a.internal[i] = from.Get(i - s)
				}
			} else {
				copy(a.internal[s:e], from.internal)
			}
		},
		func(a *Array, b *Array) {
			if a.meta != b.meta {
				for i := 0; i < b.Len(); i++ {
					a.internal = append(a.internal, b.Get(i))
				}
			} else {
				a.internal = append(a.internal, b.internal...)
			}
		},
		func(a *Array, mt typ.MarshalType) []byte {
			p := &bytes.Buffer{}
			p.WriteString("[")
			a.ForeachIndex(func(i int, v Value) bool {
				v.toString(p, 1, mt)
				p.WriteString(",")
				return true
			})
			closeBuffer(p, "]")
			return p.Bytes()
		},
	}
	*bytesArrayMeta = ArrayMeta{
		"bytes",
		ArrayProto,
		func(a *Array) int { return len((a.any).([]byte)) },
		func(a *Array) int { return cap((a.any).([]byte)) },
		func(a *Array) { a.any = a.any.([]byte)[:0] },
		func(a *Array) []Value { a.notSupported("Values"); return nil },
		func(a *Array, idx int) Value { return Int64(int64(a.any.([]byte)[idx])) },
		func(a *Array, idx int, v Value) {
			a.any.([]byte)[idx] = byte(v.Is(typ.Number, "bytes.Set").Int())
		},
		func(a *Array, v ...Value) {
			p := a.any.([]byte)
			for _, b := range v {
				p = append(p, byte(b.Is(typ.Number, "bytes.Append").Int()))
			}
			a.any = p
		},
		func(a *Array, start, end int) *Array {
			return &Array{meta: a.meta, any: a.any.([]byte)[start:end]}
		},
		func(a *Array, start, end int) {
			a.any = a.any.([]byte)[start:end]
		},
		func(a *Array, start, end int, from *Array) {
			if from.meta == internalArrayMeta {
				buf := a.any.([]byte)
				for i := start; i < end; i++ {
					buf[i] = byte(from.Get(i-start).Is(typ.Number, "bytes.Copy").Int())
				}
			} else {
				copy(a.any.([]byte)[start:end], from.any.([]byte))
			}
		},
		func(a *Array, b *Array) {
			if b.meta == internalArrayMeta {
				buf := a.any.([]byte)
				for i := 0; i < b.Len(); i++ {
					buf[i] = byte(b.Get(i).Is(typ.Number, "bytes.Concat").Int())
				}
				a.any = buf
			} else {
				a.any = append(a.any.([]byte), b.any.([]byte)...)
			}
		},
		func(a *Array, mt typ.MarshalType) []byte {
			if mt != typ.MarshalToJSON {
				return sgMarshal(a, mt)
			}
			buf := a.any.([]byte)
			tmp := make([]byte, base64.StdEncoding.EncodedLen(len(buf)))
			base64.StdEncoding.Encode(tmp, buf)
			return tmp
		},
	}
	*stringsArrayMeta = ArrayMeta{
		"[]string",
		ArrayProto,
		func(a *Array) int { return len((a.any).([]string)) },
		func(a *Array) int { return cap((a.any).([]string)) },
		func(a *Array) { a.any = a.any.([]byte)[:0] },
		func(a *Array) []Value {
			res := make([]Value, a.Len())
			for i := 0; i < a.Len(); i++ {
				res[i] = a.Get(i)
			}
			return res
		},
		func(a *Array, idx int) Value { return Str(a.any.([]string)[idx]) },
		func(a *Array, idx int, v Value) {
			a.any.([]string)[idx] = v.Is(typ.String, "[]string.Set").Str()
		},
		func(a *Array, v ...Value) {
			p := a.any.([]string)
			for _, b := range v {
				p = append(p, b.Is(typ.String, "[]string.Append").Str())
			}
			a.any = p
		},
		func(a *Array, start, end int) *Array {
			p := a.any.([]string)[start:end]
			return &Array{meta: a.meta, any: p}
		},
		func(a *Array, start, end int) {
			a.any = a.any.([]string)[start:end]
		},
		func(a *Array, start, end int, from *Array) {
			if from.meta == internalArrayMeta {
				buf := a.any.([]string)
				for i := start; i < end; i++ {
					buf[i] = from.Get(i-start).Is(typ.String, "[]string.Copy").Str()
				}
			} else {
				copy(a.any.([]byte)[start:end], from.any.([]byte))
			}
		},
		func(a *Array, b *Array) {
			if b.meta == internalArrayMeta {
				buf := a.any.([]string)
				for i := 0; i < b.Len(); i++ {
					buf[i] = b.Get(i).Is(typ.String, "[]string.Concat").Str()
				}
				a.any = buf
			} else {
				a.any = append(a.any.([]byte), b.any.([]byte)...)
			}
		},
		sgMarshal,
	}
	*errorArrayMeta = ArrayMeta{
		"error",
		ErrorProto,
		func(a *Array) int { return 1 },
		func(a *Array) int { return 1 },
		sgClearNotSupported,
		sgValuesNotSupported,
		func(a *Array, idx int) Value { return a.ToValue() },
		func(a *Array, idx int, v Value) { a.notSupported("Set") },
		sgAppendNotSupported,
		sgSliceNotSupported,
		sgSliceInplaceNotSupported,
		sgCopyNotSupported,
		sgConcatNotSupported,
		func(a *Array, mt typ.MarshalType) []byte {
			return []byte(ifquote(mt == typ.MarshalToJSON, a.any.(*ExecError).Error()))
		},
	}
}

func GetTypedArrayMeta(v interface{}) *ArrayMeta {
	switch v.(type) {
	case []Value:
		return internalArrayMeta
	case []byte:
		return bytesArrayMeta
	case []string:
		return stringsArrayMeta
	}
	rt := reflect.TypeOf(v)
	if rt.Kind() != reflect.Slice && rt.Kind() != reflect.Array {
		internal.Panic("not array or slice: %v", rt.String())
	}
	if v, ok := genericArrayMetaCache.Load(rt); ok {
		return v.(*ArrayMeta)
	}
	a := &ArrayMeta{rt.String(), ArrayProto, sgLen, sgSize, sgClear, sgValues, sgGet, sgSet, sgAppend, sgSlice, sgSliceInplace, sgCopy, sgConcat, sgMarshal}
	if rt.Kind() == reflect.Array {
		a.SliceInplace = sgSliceInplaceNotSupported
		a.Clear = sgClearNotSupported
		a.Append = sgAppendNotSupported
		a.Copy = sgCopyNotSupported
		a.Concat = sgConcatNotSupported
	}
	genericArrayMetaCache.Store(rt, a)
	return a
}

type Array struct {
	meta     *ArrayMeta
	internal []Value
	any      interface{}
}

// NewArray creates an array consists of given arguments
func NewArray(m ...Value) *Array {
	return &Array{meta: internalArrayMeta, internal: m}
}

func NewTypedArray(any interface{}, meta *ArrayMeta) *Array {
	return &Array{meta: meta, any: any}
}

// Error creates a builtin error, env can be nil
func Error(e *Env, err error) Value {
	if err == nil {
		return Nil
	} else if _, ok := err.(*ExecError); ok {
		return NewTypedArray(err, errorArrayMeta).ToValue()
	}
	ee := &ExecError{root: err}
	if e != nil {
		ee.stacks = e.GetFullStacktrace()
	}
	return NewTypedArray(ee, errorArrayMeta).ToValue()
}

func (a *Array) ToValue() Value {
	return Value{v: uint64(typ.Array), p: unsafe.Pointer(a)}
}

func (a *Array) Unwrap() interface{} {
	if a.meta == internalArrayMeta {
		return a.internal
	}
	return a.any
}

func (a *Array) Len() int {
	if a.meta == internalArrayMeta {
		return len(a.internal)
	}
	return a.meta.Len(a)
}

func (a *Array) Size() int {
	if a.meta == internalArrayMeta {
		return cap(a.internal)
	}
	return a.meta.Size(a)
}

func (a *Array) Values() []Value {
	if a.meta == internalArrayMeta {
		return a.internal
	}
	return a.meta.Values(a)
}

func (a *Array) Get(v int) Value {
	if a.meta == internalArrayMeta {
		return a.internal[v]
	}
	return a.meta.Get(a, v)
}

func (a *Array) Set(idx int, v Value) {
	if a.meta == internalArrayMeta {
		a.internal[idx] = v
	} else {
		a.meta.Set(a, idx, v)
	}
}

func (a *Array) Append(v ...Value) {
	a.meta.Append(a, v...)
}

func (a *Array) Slice(start, end int) *Array {
	return a.meta.Slice(a, start, end)
}

func (a *Array) SliceInplace(start, end int) {
	a.meta.SliceInplace(a, start, end)
}

func (a *Array) Clear() {
	if a.meta == internalArrayMeta {
		a.internal = a.internal[:0]
	} else {
		a.meta.Clear(a)
	}
}

func (a *Array) Copy(start, end int, from *Array) {
	if a.meta == internalArrayMeta || from.meta == internalArrayMeta {
	} else if a.meta != from.meta {
		internal.Panic("copy array with different types: from %q to %q", from.meta.Name, a.meta.Name)
	}
	a.meta.Copy(a, start, end, from)
}

func (a *Array) Concat(b *Array) {
	if a.meta == internalArrayMeta || b.meta == internalArrayMeta {
	} else if a.meta != b.meta {
		internal.Panic("concat array with different types: from %q to %q", b.meta.Name, a.meta.Name)
	}
	a.meta.Concat(a, b)
}

func (a *Array) Marshal(mt typ.MarshalType) []byte {
	return a.meta.Marshal(a, mt)
}

func (a *Array) ForeachIndex(f func(k int, v Value) bool) {
	for i := 0; i < a.Len(); i++ {
		if !f(i, a.Get(i)) {
			break
		}
	}
}

func (a *Array) Foreach(f func(k Value, v *Value) bool) {
	if a.meta != internalArrayMeta {
		internal.Panic("can't iterate typed array using untyped foreach")
	}
	for i := 0; i < a.Len(); i++ {
		if !f(Int(i), &a.internal[i]) {
			break
		}
	}
}

func (a *Array) Typed() bool {
	return a.meta != internalArrayMeta
}

func (a *Array) notSupported(method string) {
	panic(a.meta.Name + "." + method + " not allowed")
}

func sgLen(a *Array) int {
	return reflect.ValueOf(a.any).Len()
}

func sgSize(a *Array) int {
	return reflect.ValueOf(a.any).Cap()
}

func sgClear(a *Array) {
	a.any = reflect.ValueOf(a.any).Slice(0, 0).Interface()
}

func sgValues(a *Array) []Value {
	res := make([]Value, a.Len())
	for i := 0; i < len(res); i++ {
		res[i] = a.Get(i)
	}
	return res
}

func sgGet(a *Array, idx int) Value {
	return ValueOf(reflect.ValueOf(a.any).Index(idx).Interface())
}

func sgSet(a *Array, idx int, v Value) {
	rv := reflect.ValueOf(a.any)
	rv.Index(idx).Set(v.ReflectValue(rv.Type().Elem()))
}

func sgAppend(a *Array, v ...Value) {
	rv := reflect.ValueOf(a.any)
	rt := rv.Type().Elem()
	for _, b := range v {
		rv = reflect.Append(rv, b.ReflectValue(rt))
	}
	a.any = rv.Interface()
}

func sgSlice(a *Array, start, end int) *Array {
	return &Array{meta: a.meta, any: reflect.ValueOf(a.any).Slice(start, end).Interface()}
}

func sgSliceInplace(a *Array, start, end int) {
	a.any = reflect.ValueOf(a.any).Slice(start, end).Interface()
}

func sgCopy(a *Array, start, end int, from *Array) {
	if from.meta == internalArrayMeta {
		rv := reflect.ValueOf(a.any)
		rt := rv.Type().Elem()
		for i := start; i < end; i++ {
			rv.Index(i).Set(from.Get(i - start).ReflectValue(rt))
		}
	} else {
		reflect.Copy(reflect.ValueOf(a.any).Slice(start, end), reflect.ValueOf(from.any))
	}
}

func sgConcat(a *Array, b *Array) {
	if b.meta == internalArrayMeta {
		rv := reflect.ValueOf(a.any)
		rt := rv.Type().Elem()
		for i := 0; i < b.Len(); i++ {
			rv = reflect.Append(rv, b.Get(i).ReflectValue(rt))
		}
		a.any = rv.Interface()
	} else {
		a.any = reflect.AppendSlice(reflect.ValueOf(a.any), reflect.ValueOf(b.any)).Interface()
	}
}

func sgMarshal(a *Array, mt typ.MarshalType) []byte {
	if mt != typ.MarshalToJSON {
		return []byte(fmt.Sprint(a.any))
	}
	buf, _ := json.Marshal(a.any)
	return buf
}

func sgSliceNotSupported(a *Array, start int, end int) *Array {
	a.notSupported("Slice")
	return nil
}

func sgSliceInplaceNotSupported(a *Array, start int, end int) {
	a.notSupported("SliceInplace")
}

func sgClearNotSupported(a *Array) {
	a.notSupported("Clear")
}

func sgValuesNotSupported(a *Array) []Value {
	a.notSupported("Values")
	return nil
}

func sgAppendNotSupported(a *Array, v ...Value) {
	a.notSupported("Append")
}

func sgCopyNotSupported(a *Array, start, end int, from *Array) {
	a.notSupported("Copy")
}

func sgConcatNotSupported(a *Array, b *Array) {
	a.notSupported("Concat")
}
