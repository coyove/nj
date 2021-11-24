package nj

import (
	"reflect"
	"sync"
	"unsafe"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

type SequenceMeta struct {
	Name   string
	Len    func(*Sequence) int
	Size   func(*Sequence) int
	Clear  func(*Sequence)
	Values func(*Sequence) []Value
	Get    func(*Sequence, int) Value
	Set    func(*Sequence, int, Value)
	Append func(*Sequence, ...Value)
	Slice  func(*Sequence, int, int) *Sequence
	Copy   func(*Sequence, int, int, *Sequence)
	Concat func(*Sequence, *Sequence)
}

var internalSequenceMeta = &SequenceMeta{
	"internal",
	func(a *Sequence) int { return len(a.internal) },
	func(a *Sequence) int { return cap(a.internal) },
	func(a *Sequence) { a.internal = a.internal[:0] },
	func(a *Sequence) []Value { return a.internal },
	func(a *Sequence, idx int) Value { return a.internal[idx] },
	func(a *Sequence, idx int, v Value) { a.internal[idx] = v },
	func(a *Sequence, v ...Value) { a.internal = append(a.internal, v...) },
	func(a *Sequence, s, e int) *Sequence { return &Sequence{meta: a.meta, internal: a.internal[s:e]} },
	func(a *Sequence, s, e int, from *Sequence) { copy(a.internal[s:e], from.internal) },
	func(a *Sequence, b *Sequence) { a.internal = append(a.internal, b.internal...) },
}

var bytesSequenceMeta = &SequenceMeta{
	"bytes",
	func(a *Sequence) int { return len((a.any).([]byte)) },
	func(a *Sequence) int { return cap((a.any).([]byte)) },
	func(a *Sequence) { a.any = a.any.([]byte)[:0] },
	func(a *Sequence) []Value { panic("sequence(bytes).values: can't convert bytes") },
	func(a *Sequence, idx int) Value { return Int64(int64(a.any.([]byte)[idx])) },
	func(a *Sequence, idx int, v Value) {
		a.any.([]byte)[idx] = byte(v.Is(typ.Number, "sequence(bytes).set").Int())
	},
	func(a *Sequence, v ...Value) {
		p := a.any.([]byte)
		for _, b := range v {
			p = append(p, byte(b.Is(typ.Number, "sequence(bytes).append").Int()))
		}
		a.any = p
	},
	func(a *Sequence, start, end int) *Sequence {
		p := a.any.([]byte)[start:end]
		return &Sequence{meta: a.meta, any: p}
	},
	func(a *Sequence, start, end int, from *Sequence) {
		copy(a.any.([]byte)[start:end], from.any.([]byte))
	},
	func(a *Sequence, b *Sequence) {
		a.any = append(a.any.([]byte), b.any.([]byte)...)
	},
}

var genericSequenceMetaCache sync.Map

func GetGenericSequenceMeta(v interface{}) *SequenceMeta {
	rt := reflect.TypeOf(v)
	n := rt.String()
	if v, ok := genericSequenceMetaCache.Load(n); ok {
		return v.(*SequenceMeta)
	}
	a := &SequenceMeta{n, sgLen, sgSize, sgClear, sgValues, sgGet, sgSet, sgAppend, sgSlice, sgCopy, sgConcat}
	if rt.Kind() == reflect.Array {
		a.Clear = func(a *Sequence) { panic("sequence(" + a.meta.Name + ").clear") }
		a.Append = func(a *Sequence, v ...Value) { panic("sequence(" + a.meta.Name + ").append") }
		a.Copy = func(a *Sequence, start, end int, from *Sequence) { panic("sequence(" + a.meta.Name + ").copy") }
		a.Concat = func(a *Sequence, b *Sequence) { panic("sequence(" + a.meta.Name + ").concat") }
	}
	genericSequenceMetaCache.Store(n, a)
	return a
}

type Sequence struct {
	meta     *SequenceMeta
	internal []Value
	any      interface{}
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

func (a *Sequence) Clear() {
	if a.meta == internalSequenceMeta {
		a.internal = a.internal[:0]
	} else {
		a.meta.Clear(a)
	}
}

func (a *Sequence) Copy(start, end int, from *Sequence) {
	if a.meta != from.meta {
		internal.Panic("copy sequences with different types: from %q to %q", from.meta.Name, a.meta.Name)
	}
	a.meta.Copy(a, start, end, from)
}

func (a *Sequence) Concat(b *Sequence) {
	if a.meta != b.meta {
		internal.Panic("concat sequences with different types: from %q to %q", b.meta.Name, a.meta.Name)
	}
	a.meta.Concat(a, b)
}

func (a *Sequence) Foreach(f func(k, v Value) bool) {
	for i := 0; i < a.Len(); i++ {
		if !f(Int(i), a.Get(i)) {
			break
		}
	}
}

func sgLen(a *Sequence) int {
	return reflect.ValueOf(a.any).Len()
}

func sgSize(a *Sequence) int {
	return reflect.ValueOf(a.any).Len()
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
	return Val(reflect.ValueOf(a.any).Index(idx).Interface())
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

func sgCopy(a *Sequence, start, end int, from *Sequence) {
	reflect.Copy(reflect.ValueOf(a.any).Slice(start, end), reflect.ValueOf(from.any))
}

func sgConcat(a *Sequence, b *Sequence) {
	a.any = reflect.AppendSlice(reflect.ValueOf(a.any), reflect.ValueOf(b.any)).Interface()
}
