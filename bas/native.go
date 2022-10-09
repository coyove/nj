package bas

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sync"
	"unsafe"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

var metaStore struct {
	sync.Mutex
	metaCache   map[reflect.Type]*NativeMeta
	nativeTypes *Object
}

type NativeMeta struct {
	Name         string
	Proto        *Object
	Len          func(*Native) int
	Cap          func(*Native) int
	Clear        func(*Native)
	Values       func(*Native) []Value
	Get          func(*Native, int) Value
	Set          func(*Native, int, Value)
	GetKey       func(*Native, Value) (Value, bool)
	SetKey       func(*Native, Value, Value)
	Append       func(*Native, ...Value)
	Slice        func(*Native, int, int) *Native
	SliceInplace func(*Native, int, int)
	Copy         func(*Native, int, int, *Native)
	Concat       func(*Native, *Native)
	Marshal      func(*Native, io.Writer, typ.MarshalType)
	Next         func(*Native, Value) Value
}

func NewEmptyNativeMeta(name string, proto *Object) *NativeMeta {
	if proto == nil {
		proto = &Proto.Native
	}
	return createNativeMeta(name, proto)
}

func createNativeMeta(name string, proto *Object) *NativeMeta {
	return &NativeMeta{name, proto,
		sgLenNotSupported,
		sgCapNotSupported,
		sgClearNotSupported,
		sgValuesNotSupported,
		func(a *Native, idx int) Value { return a.ToValue() },
		sgSetNotSupported,
		sgGetKeyWithReflection,
		sgSetKeyNotSupported,
		sgAppendNotSupported,
		sgSliceNotSupported,
		sgSliceInplaceNotSupported,
		sgCopyNotSupported,
		sgConcatNotSupported,
		func(a *Native, w io.Writer, mt typ.MarshalType) {
			if mt != typ.MarshalToJSON {
				internal.StringifyTo(w, a.any)
			} else {
				json.NewEncoder(w).Encode(a.any)
			}
		},
		sgNextNotSupported,
	}
}

func getNativeMeta(v interface{}) *NativeMeta {
	metaStore.Lock()
	defer metaStore.Unlock()

	switch v.(type) {
	case []Value:
		return &Proto.ArrayMeta
	case []byte:
		return &Proto.BytesMeta
	case []string:
		return &Proto.StringsMeta
	case *ExecError:
		return &Proto.ErrorMeta
	}
	rt := reflect.TypeOf(v)
	if v, ok := metaStore.metaCache[rt]; ok {
		return v
	}

	var a *NativeMeta
	switch rt.Kind() {
	default:
		a = NewEmptyNativeMeta(reflectTypeName(rt), &Proto.Native)
		a.GetKey = sgGetKeyWithReflection
		a.SetKey = func(a *Native, k, v Value) {
			defer func() {
				if r := recover(); r != nil {
					internal.Panic("%s.SetKey(%v, %v): %v", a.meta.Name, k.simple(), v.simple(), r)
				}
			}()
			rv := reflect.ValueOf(a.any)
			if rv.Kind() == reflect.Map {
				rv.SetMapIndex(k.ToType(rv.Type().Key()), v.ToType(rv.Type().Elem()))
			} else {
				f := reflect.Indirect(rv).FieldByName(k.AssertString("key"))
				f.Set(v.ToType(f.Type()))
			}
		}
		tn := internal.SanitizeName(a.Name)
		pt := Func(string(tn), func(e *Env) { e.A = ValueOf(reflect.New(rt).Elem().Interface()) }).Object()
		switch rt.Kind() {
		case reflect.Map:
			a.Len = sgLen
			a.Cap = sgCap
			a.Marshal = sgMarshal
			a.Next = sgMapNext
			a.Set = func(n *Native, idx int, v Value) {
				rv := reflect.ValueOf(n.any)
				rv.SetMapIndex(Int(idx).ToType(rv.Type().Key()), v.ToType(rv.Type().Elem()))
			}
			a.Get = func(n *Native, idx int) Value {
				v, _ := sgReflectLoadSafe(n.any, Int(idx))
				return v
			}
			a.Proto = pt.SetPrototype(&Proto.NativeMap)
		case reflect.Ptr:
			a.Proto = pt.SetPrototype(&Proto.NativePtr)
		default:
			a.Proto = pt.SetPrototype(&Proto.Native)
		}
		metaStore.nativeTypes.SetProp(pt.Name(), pt.ToValue())
	case reflect.Chan:
		a = NewEmptyNativeMeta(reflectTypeName(rt), &Proto.Channel)
		a.Len = sgLen
		a.Cap = sgCap
		a.Next = func(a *Native, kv Value) Value {
			if kv == Nil {
				kv = Array(Nil, Nil)
			} else {
				if kv.Native().Get(1).IsFalse() {
					return Nil // break because the channel has been closed
				}
			}
			rv, ok := reflect.ValueOf(a.any).Recv()
			if !ok {
				return Nil
			}
			kv.Native().Set(0, ValueOf(rv.Interface()))
			kv.Native().Set(1, Bool(ok))
			return kv
		}
	case reflect.Array, reflect.Slice:
		a = &NativeMeta{reflectTypeName(rt), &Proto.Array,
			sgLen, sgCap, sgClear, sgValues, sgGet, sgSet, sgGetKey, sgSetKeyNotSupported, sgAppend, sgSlice,
			sgSliceInplace, sgCopy, sgConcat, sgMarshal, sgArrayNext}
		if rt.Kind() == reflect.Array {
			a.SliceInplace = sgSliceInplaceNotSupported
			a.Clear = sgClearNotSupported
			a.Append = sgAppendNotSupported
			a.Copy = sgCopyNotSupported
			a.Concat = sgConcatNotSupported
		}
	}
	metaStore.metaCache[rt] = a
	return a
}

type Native struct {
	meta     *NativeMeta
	internal []Value
	any      interface{}
}

func NewNative(any interface{}) *Native {
	return NewNativeWithMeta(any, getNativeMeta(any))
}

func newArray(m ...Value) *Native {
	return &Native{meta: &Proto.ArrayMeta, internal: m}
}

func newVarargArray(m []Value) *Native {
	return &Native{meta: &Proto.VarargMeta, internal: append([]Value{}, m...)}
}

func NewNativeWithMeta(any interface{}, meta *NativeMeta) *Native {
	return &Native{meta: meta, any: any}
}

func (a *Native) ToValue() Value {
	return Value{v: uint64(typ.Native), p: unsafe.Pointer(a)}
}

func (a *Native) Unwrap() interface{} {
	if a.meta == &Proto.ArrayMeta {
		return a.internal
	}
	return a.any
}

func (a *Native) UnwrapFunc(f func(interface{}) interface{}) {
	if a.meta == &Proto.ArrayMeta {
		a.internal = f(a.Unwrap()).([]Value)
	} else {
		a.any = f(a.Unwrap())
	}
}

func (a *Native) Len() int {
	if a.meta == &Proto.ArrayMeta {
		return len(a.internal)
	}
	return a.meta.Len(a)
}

func (a *Native) Cap() int {
	if a.meta == &Proto.ArrayMeta {
		return cap(a.internal)
	}
	return a.meta.Cap(a)
}

func (a *Native) Values() []Value {
	if a.meta == &Proto.ArrayMeta {
		return a.internal
	}
	return a.meta.Values(a)
}

func (a *Native) Get(v int) Value {
	if a.meta == &Proto.ArrayMeta {
		return a.internal[v]
	}
	return a.meta.Get(a, v)
}

func (a *Native) Set(idx int, v Value) {
	if a.meta == &Proto.ArrayMeta {
		a.internal[idx] = v
	} else {
		a.meta.Set(a, idx, v)
	}
}

func (a *Native) GetKey(k Value) (Value, bool) {
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
	if a.meta == &Proto.ArrayMeta {
		a.internal = a.internal[:0]
	} else {
		a.meta.Clear(a)
	}
}

func (a *Native) Copy(start, end int, from *Native) {
	if a.meta == &Proto.ArrayMeta || from.meta == &Proto.ArrayMeta {
	} else if a.meta != from.meta {
		internal.Panic("copy array with different types: from %q to %q", from.meta.Name, a.meta.Name)
	}
	a.meta.Copy(a, start, end, from)
}

func (a *Native) Concat(b *Native) {
	if b == nil {
		return
	}
	if a.meta == &Proto.ArrayMeta || b.meta == &Proto.ArrayMeta {
	} else if a.meta != b.meta {
		internal.Panic("concat array with different types: from %q to %q", b.meta.Name, a.meta.Name)
	}
	a.meta.Concat(a, b)
}

func (a *Native) Marshal(w io.Writer, mt typ.MarshalType) {
	a.meta.Marshal(a, w, mt)
}

func (a *Native) internalNext(k Value) Value {
	if a.meta == &Proto.ArrayMeta {
		return sgArrayNext(a, k)
	}
	return a.meta.Next(a, k)
}

func (a *Native) IsUntypedArray() bool {
	return a.meta == &Proto.ArrayMeta
}

func (a *Native) IsTypedArray() bool {
	return a.meta.Proto.HasPrototype(&Proto.Array) && a.meta != &Proto.ArrayMeta
}

func (a *Native) Prototype() *Object {
	return a.meta.Proto
}

func (a *Native) HasPrototype(p *Object) bool {
	return a.meta.Proto.HasPrototype(p)
}

func (a *Native) notSupported(method string) {
	panic(a.meta.Name + "." + method + " not supported")
}

func sgLen(a *Native) int {
	return reflect.ValueOf(a.any).Len()
}

func sgCap(a *Native) int {
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
	rv.Index(idx).Set(v.ToType(rv.Type().Elem()))
}

func sgAppend(a *Native, v ...Value) {
	rv := reflect.ValueOf(a.any)
	rt := rv.Type().Elem()
	for _, b := range v {
		rv = reflect.Append(rv, b.ToType(rt))
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
	if from.meta == &Proto.ArrayMeta {
		rv := reflect.ValueOf(a.any)
		rt := rv.Type().Elem()
		for i := start; i < end; i++ {
			rv.Index(i).Set(from.Get(i - start).ToType(rt))
		}
	} else {
		reflect.Copy(reflect.ValueOf(a.any).Slice(start, end), reflect.ValueOf(from.any))
	}
}

func sgConcat(a *Native, b *Native) {
	if b.meta == &Proto.ArrayMeta {
		rv := reflect.ValueOf(a.any)
		rt := rv.Type().Elem()
		for i := 0; i < b.Len(); i++ {
			rv = reflect.Append(rv, b.Get(i).ToType(rt))
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

func sgCapNotSupported(a *Native) int {
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

func sgNextNotSupported(a *Native, b Value) Value {
	a.notSupported("Next")
	return Nil
}

func sgArrayNext(a *Native, kv Value) Value {
	al := a.Len()
	if al == 0 {
		return Nil
	}
	if kv == Nil {
		return Array(Int(0), a.Get(0))
	}
	idx := kv.Native().Get(0).AssertNumber("array iteration").Int() + 1
	if idx >= al {
		return Nil
	}
	kv.Native().Set(0, Int(idx))
	kv.Native().Set(1, a.Get(idx))
	return kv
}

func sgMapNext(a *Native, kv Value) Value {
	if a.Len() == 0 {
		return Nil
	}
	if kv == Nil {
		iter := reflect.ValueOf(a.any).MapRange()
		if iter.Next() {
			return Array(ValueOf(iter.Key()), ValueOf(iter.Value()), NewNative(iter).ToValue())
		}
		return Nil
	}
	iter := kv.Native().Get(2).Interface().(*reflect.MapIter)
	if !iter.Next() {
		return Nil
	}
	kv.Native().Set(0, ValueOf(iter.Key()))
	kv.Native().Set(1, ValueOf(iter.Value()))
	return kv
}

func sgGetKey(a *Native, k Value) (Value, bool) {
	f, ok := a.meta.Proto.Find(k)
	if f != Nil {
		f = setObjectRecv(f, a.ToValue())
	}
	return f, ok
}

func sgGetKeyWithReflection(a *Native, k Value) (Value, bool) {
	if v, ok := sgReflectLoadSafe(a.any, k); ok {
		return v, true
	}
	return sgGetKey(a, k)
}

func sgReflectLoadSafe(v interface{}, key Value) (value Value, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			value, ok = Nil, false
		}
	}()

	if reflect.TypeOf(v).Kind() == reflect.Map {
		rv := reflect.ValueOf(v)
		if v := rv.MapIndex(key.ToType(rv.Type().Key())); v.IsValid() {
			return ValueOf(v.Interface()), true
		}
		return Nil, false
	}
	if key.Type() != typ.String {
		return Nil, false
	}
	k := key.Str()
	if len(k) > 0 && k[0] >= 'A' && k[0] <= 'Z' {
		rv := reflect.ValueOf(v)
		f := rv.MethodByName(k)
		if !f.IsValid() {
			f = reflect.Indirect(rv).MethodByName(k)
		}
		if !f.IsValid() {
			f = reflect.Indirect(rv).FieldByName(k)
		}
		if !f.IsValid() {
			return Nil, false
		}
		return ValueOf(f.Interface()), true
	}
	return Nil, false
}

func reflectTypeName(t reflect.Type) string {
	res := []byte(t.String())
	for i := 0; i < len(res); i++ {
		switch res[i] {
		case '<', '(', '[', '{', ':', ')', '}', ']', '>', ' ', ',':
			res[i] = '-'
		}
	}
	return string(res)
}

func (a *Native) wrapCall(k string, args ...Value) Value {
	o := a.meta.Proto.Get(Str(k)).AssertShape("No", "wrapper."+k).Object()
	if o == nil {
		internal.Panic("wrapper.%s not implemented", k)
	}
	o = o.Copy()
	o.this = a.any.(*Object).ToValue()
	return o.Call(nil, args...)
}
