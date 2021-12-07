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

type SequenceMeta struct {
	Name         string
	Proto        *Object
	Len          func(*Sequence) int
	Size         func(*Sequence) int
	Clear        func(*Sequence)
	Values       func(*Sequence) []Value
	Get          func(*Sequence, int) Value
	Set          func(*Sequence, int, Value)
	Append       func(*Sequence, ...Value)
	Slice        func(*Sequence, int, int) *Sequence
	SliceInplace func(*Sequence, int, int)
	Copy         func(*Sequence, int, int, *Sequence)
	Concat       func(*Sequence, *Sequence)
	Marshal      func(*Sequence, typ.MarshalType) []byte
}

var (
	internalSequenceMeta     = &SequenceMeta{}
	bytesSequenceMeta        = &SequenceMeta{}
	stringsSequenceMeta      = &SequenceMeta{}
	errorSequenceMeta        = &SequenceMeta{}
	genericSequenceMetaCache sync.Map
)

func init() {
	*internalSequenceMeta = SequenceMeta{
		"internal",
		ArrayLib,
		func(a *Sequence) int { return len(a.internal) },
		func(a *Sequence) int { return cap(a.internal) },
		func(a *Sequence) { a.internal = a.internal[:0] },
		func(a *Sequence) []Value { return a.internal },
		func(a *Sequence, idx int) Value { return a.internal[idx] },
		func(a *Sequence, idx int, v Value) { a.internal[idx] = v },
		func(a *Sequence, v ...Value) { a.internal = append(a.internal, v...) },
		func(a *Sequence, s, e int) *Sequence { return &Sequence{meta: a.meta, internal: a.internal[s:e]} },
		func(a *Sequence, s, e int) { a.internal = a.internal[s:e] },
		func(a *Sequence, s, e int, from *Sequence) {
			if from.meta != a.meta {
				for i := s; i < e; i++ {
					a.internal[i] = from.Get(i - s)
				}
			} else {
				copy(a.internal[s:e], from.internal)
			}
		},
		func(a *Sequence, b *Sequence) {
			if a.meta != b.meta {
				for i := 0; i < b.Len(); i++ {
					a.internal = append(a.internal, b.Get(i))
				}
			} else {
				a.internal = append(a.internal, b.internal...)
			}
		},
		func(a *Sequence, mt typ.MarshalType) []byte {
			p := &bytes.Buffer{}
			p.WriteString("[")
			a.Foreach(func(i int, v Value) bool {
				v.toString(p, 1, mt)
				p.WriteString(",")
				return true
			})
			closeBuffer(p, "]")
			return p.Bytes()
		},
	}
	*bytesSequenceMeta = SequenceMeta{
		"bytes",
		ArrayLib,
		func(a *Sequence) int { return len((a.any).([]byte)) },
		func(a *Sequence) int { return cap((a.any).([]byte)) },
		func(a *Sequence) { a.any = a.any.([]byte)[:0] },
		func(a *Sequence) []Value { a.notSupported("Values"); return nil },
		func(a *Sequence, idx int) Value { return Int64(int64(a.any.([]byte)[idx])) },
		func(a *Sequence, idx int, v Value) {
			a.any.([]byte)[idx] = byte(v.Is(typ.Number, "sequence(bytes).Set").Int())
		},
		func(a *Sequence, v ...Value) {
			p := a.any.([]byte)
			for _, b := range v {
				p = append(p, byte(b.Is(typ.Number, "sequence(bytes).Append").Int()))
			}
			a.any = p
		},
		func(a *Sequence, start, end int) *Sequence {
			return &Sequence{meta: a.meta, any: a.any.([]byte)[start:end]}
		},
		func(a *Sequence, start, end int) {
			a.any = a.any.([]byte)[start:end]
		},
		func(a *Sequence, start, end int, from *Sequence) {
			if from.meta == internalSequenceMeta {
				buf := a.any.([]byte)
				for i := start; i < end; i++ {
					buf[i] = byte(from.Get(i-start).Is(typ.Number, "sequence(bytes).Copy").Int())
				}
			} else {
				copy(a.any.([]byte)[start:end], from.any.([]byte))
			}
		},
		func(a *Sequence, b *Sequence) {
			if b.meta == internalSequenceMeta {
				buf := a.any.([]byte)
				for i := 0; i < b.Len(); i++ {
					buf[i] = byte(b.Get(i).Is(typ.Number, "sequence(bytes).Concat").Int())
				}
				a.any = buf
			} else {
				a.any = append(a.any.([]byte), b.any.([]byte)...)
			}
		},
		func(a *Sequence, mt typ.MarshalType) []byte {
			if mt != typ.MarshalToJSON {
				return sgMarshal(a, mt)
			}
			buf := a.any.([]byte)
			tmp := make([]byte, base64.StdEncoding.EncodedLen(len(buf)))
			base64.StdEncoding.Encode(tmp, buf)
			return tmp
		},
	}
	*stringsSequenceMeta = SequenceMeta{
		"[]string",
		ArrayLib,
		func(a *Sequence) int { return len((a.any).([]string)) },
		func(a *Sequence) int { return cap((a.any).([]string)) },
		func(a *Sequence) { a.any = a.any.([]byte)[:0] },
		func(a *Sequence) []Value {
			res := make([]Value, a.Len())
			for i := 0; i < a.Len(); i++ {
				res[i] = a.Get(i)
			}
			return res
		},
		func(a *Sequence, idx int) Value { return Str(a.any.([]string)[idx]) },
		func(a *Sequence, idx int, v Value) {
			a.any.([]string)[idx] = v.Is(typ.String, "sequence(string).Set").Str()
		},
		func(a *Sequence, v ...Value) {
			p := a.any.([]string)
			for _, b := range v {
				p = append(p, b.Is(typ.String, "sequence(string).Append").Str())
			}
			a.any = p
		},
		func(a *Sequence, start, end int) *Sequence {
			p := a.any.([]string)[start:end]
			return &Sequence{meta: a.meta, any: p}
		},
		func(a *Sequence, start, end int) {
			a.any = a.any.([]string)[start:end]
		},
		func(a *Sequence, start, end int, from *Sequence) {
			if from.meta == internalSequenceMeta {
				buf := a.any.([]string)
				for i := start; i < end; i++ {
					buf[i] = from.Get(i-start).Is(typ.String, "sequence(string).Copy").Str()
				}
			} else {
				copy(a.any.([]byte)[start:end], from.any.([]byte))
			}
		},
		func(a *Sequence, b *Sequence) {
			if b.meta == internalSequenceMeta {
				buf := a.any.([]string)
				for i := 0; i < b.Len(); i++ {
					buf[i] = b.Get(i).Is(typ.String, "sequence(string).Concat").Str()
				}
				a.any = buf
			} else {
				a.any = append(a.any.([]byte), b.any.([]byte)...)
			}
		},
		sgMarshal,
	}
	*errorSequenceMeta = SequenceMeta{
		"error",
		ErrorLib,
		func(a *Sequence) int { return 1 },
		func(a *Sequence) int { return 1 },
		sgClearNotSupported,
		sgValuesNotSupported,
		func(a *Sequence, idx int) Value { return a.ToValue() },
		func(a *Sequence, idx int, v Value) { a.notSupported("Set") },
		sgAppendNotSupported,
		sgSliceNotSupported,
		sgSliceInplaceNotSupported,
		sgCopyNotSupported,
		sgConcatNotSupported,
		func(a *Sequence, mt typ.MarshalType) []byte {
			return []byte(ifquote(mt == typ.MarshalToJSON, a.any.(*ExecError).Error()))
		},
	}
}

func GetGenericSequenceMeta(v interface{}) *SequenceMeta {
	switch v.(type) {
	case []Value:
		return internalSequenceMeta
	case []byte:
		return bytesSequenceMeta
	case []string:
		return stringsSequenceMeta
	}
	rt := reflect.TypeOf(v)
	if rt.Kind() != reflect.Slice && rt.Kind() != reflect.Array {
		internal.Panic("not generic sequence: %v", rt.String())
	}
	if v, ok := genericSequenceMetaCache.Load(rt); ok {
		return v.(*SequenceMeta)
	}
	a := &SequenceMeta{rt.String(), ArrayLib, sgLen, sgSize, sgClear, sgValues, sgGet, sgSet, sgAppend, sgSlice, sgSliceInplace, sgCopy, sgConcat, sgMarshal}
	if rt.Kind() == reflect.Array {
		a.SliceInplace = sgSliceInplaceNotSupported
		a.Clear = sgClearNotSupported
		a.Append = sgAppendNotSupported
		a.Copy = sgCopyNotSupported
		a.Concat = sgConcatNotSupported
	}
	genericSequenceMetaCache.Store(rt, a)
	return a
}

type Sequence struct {
	meta     *SequenceMeta
	internal []Value
	any      interface{}
}

// Error creates a builtin error, env can be nil
func Error(e *Env, err error) Value {
	if err == nil {
		return Nil
	} else if _, ok := err.(*ExecError); ok {
		return NewSequence(err, errorSequenceMeta).ToValue()
	}
	ee := &ExecError{root: err}
	if e != nil {
		ee.stacks = e.GetFullStacktrace()
	}
	return NewSequence(ee, errorSequenceMeta).ToValue()
}

func NewSequence(any interface{}, meta *SequenceMeta) *Sequence {
	return &Sequence{meta: meta, any: any}
}

func (a *Sequence) ToValue() Value {
	return Value{v: uint64(typ.Array), p: unsafe.Pointer(a)}
}

func (a *Sequence) Unwrap() interface{} {
	if a.meta == internalSequenceMeta {
		return a.internal
	}
	return a.any
}

func (a *Sequence) Len() int {
	if a.meta == internalSequenceMeta {
		return len(a.internal)
	}
	return a.meta.Len(a)
}

func (a *Sequence) Size() int {
	if a.meta == internalSequenceMeta {
		return cap(a.internal)
	}
	return a.meta.Size(a)
}

func (a *Sequence) Values() []Value {
	if a.meta == internalSequenceMeta {
		return a.internal
	}
	return a.meta.Values(a)
}

func (a *Sequence) Get(v int) Value {
	if a.meta == internalSequenceMeta {
		return a.internal[v]
	}
	return a.meta.Get(a, v)
}

func (a *Sequence) Set(idx int, v Value) {
	if a.meta == internalSequenceMeta {
		a.internal[idx] = v
	} else {
		a.meta.Set(a, idx, v)
	}
}

func (a *Sequence) Append(v ...Value) {
	a.meta.Append(a, v...)
}

func (a *Sequence) Slice(start, end int) *Sequence {
	return a.meta.Slice(a, start, end)
}

func (a *Sequence) SliceInplace(start, end int) {
	a.meta.SliceInplace(a, start, end)
}

func (a *Sequence) Clear() {
	if a.meta == internalSequenceMeta {
		a.internal = a.internal[:0]
	} else {
		a.meta.Clear(a)
	}
}

func (a *Sequence) Copy(start, end int, from *Sequence) {
	if a.meta == internalSequenceMeta || from.meta == internalSequenceMeta {
	} else if a.meta != from.meta {
		internal.Panic("copy sequences with different types: from %q to %q", from.meta.Name, a.meta.Name)
	}
	a.meta.Copy(a, start, end, from)
}

func (a *Sequence) Concat(b *Sequence) {
	if a.meta == internalSequenceMeta || b.meta == internalSequenceMeta {
	} else if a.meta != b.meta {
		internal.Panic("concat sequences with different types: from %q to %q", b.meta.Name, a.meta.Name)
	}
	a.meta.Concat(a, b)
}

func (a *Sequence) Marshal(mt typ.MarshalType) []byte {
	return a.meta.Marshal(a, mt)
}

func (a *Sequence) Foreach(f func(k int, v Value) bool) {
	for i := 0; i < a.Len(); i++ {
		if !f(i, a.Get(i)) {
			break
		}
	}
}

func (a *Sequence) Typed() bool {
	return a.meta != internalSequenceMeta
}

func (a *Sequence) notSupported(method string) {
	panic("sequence(" + a.meta.Name + ")." + method + " not allowed")
}

func sgLen(a *Sequence) int {
	return reflect.ValueOf(a.any).Len()
}

func sgSize(a *Sequence) int {
	return reflect.ValueOf(a.any).Cap()
}

func sgClear(a *Sequence) {
	a.any = reflect.ValueOf(a.any).Slice(0, 0).Interface()
}

func sgValues(a *Sequence) []Value {
	res := make([]Value, a.Len())
	for i := 0; i < len(res); i++ {
		res[i] = a.Get(i)
	}
	return res
}

func sgGet(a *Sequence, idx int) Value {
	return ValueOf(reflect.ValueOf(a.any).Index(idx).Interface())
}

func sgSet(a *Sequence, idx int, v Value) {
	rv := reflect.ValueOf(a.any)
	rv.Index(idx).Set(v.ReflectValue(rv.Type().Elem()))
}

func sgAppend(a *Sequence, v ...Value) {
	rv := reflect.ValueOf(a.any)
	rt := rv.Type().Elem()
	for _, b := range v {
		reflect.Append(rv, b.ReflectValue(rt))
	}
	a.any = rv.Interface()
}

func sgSlice(a *Sequence, start, end int) *Sequence {
	return &Sequence{meta: a.meta, any: reflect.ValueOf(a.any).Slice(start, end).Interface()}
}

func sgSliceInplace(a *Sequence, start, end int) {
	a.any = reflect.ValueOf(a.any).Slice(start, end).Interface()
}

func sgCopy(a *Sequence, start, end int, from *Sequence) {
	if from.meta == internalSequenceMeta {
		rv := reflect.ValueOf(a.any)
		rt := rv.Type().Elem()
		for i := start; i < end; i++ {
			rv.Index(i).Set(from.Get(i - start).ReflectValue(rt))
		}
	} else {
		reflect.Copy(reflect.ValueOf(a.any).Slice(start, end), reflect.ValueOf(from.any))
	}
}

func sgConcat(a *Sequence, b *Sequence) {
	if b.meta == internalSequenceMeta {
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

func sgMarshal(a *Sequence, mt typ.MarshalType) []byte {
	if mt != typ.MarshalToJSON {
		return []byte(fmt.Sprint(a.any))
	}
	buf, _ := json.Marshal(a.any)
	return buf
}

func sgSliceNotSupported(a *Sequence, start int, end int) *Sequence {
	a.notSupported("Slice")
	return nil
}

func sgSliceInplaceNotSupported(a *Sequence, start int, end int) {
	a.notSupported("SliceInplace")
}

func sgClearNotSupported(a *Sequence) {
	a.notSupported("Clear")
}

func sgValuesNotSupported(a *Sequence) []Value {
	a.notSupported("Values")
	return nil
}

func sgAppendNotSupported(a *Sequence, v ...Value) {
	a.notSupported("Append")
}

func sgCopyNotSupported(a *Sequence, start, end int, from *Sequence) {
	a.notSupported("Copy")
}

func sgConcatNotSupported(a *Sequence, b *Sequence) {
	a.notSupported("Concat")
}
