package potatolang

import (
	"strconv"
	"unsafe"
)

type (
	Bytes  []byte
	Int64  int64
	UInt64 uint64
)

var (
	nativeBytesMetatable = (&Table{}).
				Put(M__index, NativeFun(func(env *Env) {
			a := env.In(0, ANY).Any().(Bytes)
			switch k := env.Get(1); k.Type() {
			case NUM:
				env.A = Num(float64(a[int(k.Num())]))
			case STR:
				switch k.Str() {
				case "append":
					env.A = NativeFun(func(env *Env) {
						a := env.In(0, ANY).Any().(Bytes)
						for _, v := range env.Stack()[1:] {
							a = append(a, byte(v.ExpectMsg(NUM, "append").Num()))
						}
						env.A = Any(a)
					})
					return
				case "tostring":
					env.A = NativeFun(func(env *Env) {
						a := env.In(0, ANY).Any().(Bytes)
						env.A = Str(string(a))
					})
					return
				}
				fallthrough
			default:
				panicf("invalid index: %#v", k)
			}
		})).
		Put(M__newindex, NativeFun(func(env *Env) {
			a := env.In(0, ANY).Any().(Bytes)
			a[int(env.In(1, NUM).Num())] = byte(env.In(2, NUM).Num())
		})).
		Put(M__len, NativeFun(func(env *Env) {
			a := env.In(0, ANY).Any().(Bytes)
			env.A = Num(float64(len(a)))
		})).
		Put(M__concat, NativeFun(func(env *Env) {
			a := env.In(0, ANY).Any().(Bytes)
			b := env.In(1, ANY).Any().(Bytes)
			env.A = Any(Bytes(append(a, b...)))
		})).
		Put(M__eq, NativeFun(func(env *Env) {
			switch l, r := env.Get(0), env.Get(1); l.Type() + r.Type() {
			case AnyAny:
				a, b := l.Any().(Bytes), r.Any().(Bytes)
				env.A = Bln(false)
				if len(a) != len(b) {
					return
				}
				for i := range a {
					if a[i] != b[i] {
						return
					}
				}
				env.A = Bln(true)
			case ANY + STR:
				var a Bytes
				var b string
				a, _ = l.Any().(Bytes)
				if a == nil {
					a, b = r.Any().(Bytes), l.Str()
				} else {
					b = r.Str()
				}
				env.A = Bln(*(*string)(unsafe.Pointer(&a)) == b)
			case ANY + TAB:
				var a Bytes
				var b *Table
				a, _ = l.Any().(Bytes)
				if a == nil {
					a, b = r.Any().(Bytes), l.Tab()
				} else {
					b = r.Tab()
				}
				env.A = Bln(false)
				if len(a) != len(b.a) {
					return
				}
				for i := range a {
					bv := b.a[i]
					if bv.Type() != NUM {
						return
					}
					if float64(a[i]) != bv.Num() {
						return
					}
				}
				env.A = Bln(true)
			}
		}))
	nativeInt64Metatable = (&Table{}).
				Put(M__add, NativeFun(func(env *Env) { env.A = Any(Int64(atoint64(env.Get(0)) + atoint64(env.Get(1)))) })).
				Put(M__sub, NativeFun(func(env *Env) { env.A = Any(Int64(atoint64(env.Get(0)) - atoint64(env.Get(1)))) })).
				Put(M__mul, NativeFun(func(env *Env) { env.A = Any(Int64(atoint64(env.Get(0)) * atoint64(env.Get(1)))) })).
				Put(M__div, NativeFun(func(env *Env) { env.A = Any(Int64(atoint64(env.Get(0)) / atoint64(env.Get(1)))) })).
				Put(M__unm, NativeFun(func(env *Env) { env.A = Any(Int64(-atoint64(env.Get(0)))) })).
				Put(M__eq, NativeFun(func(env *Env) { env.A = Bln(atoint64(env.Get(0)) == atoint64(env.Get(1))) })).
				Put(M__lt, NativeFun(func(env *Env) { env.A = Bln(atoint64(env.Get(0)) < atoint64(env.Get(1))) })).
				Put(M__le, NativeFun(func(env *Env) { env.A = Bln(atoint64(env.Get(0)) <= atoint64(env.Get(1))) })).
				Put(M__tostring, NativeFun(func(env *Env) { env.A = Str(strconv.FormatInt(atoint64(env.Get(0)), 10)) }))
	nativeUInt64Metatable = (&Table{}).
				Put(M__add, NativeFun(func(env *Env) { env.A = Any(UInt64(atouint64(env.Get(0)) + atouint64(env.Get(1)))) })).
				Put(M__sub, NativeFun(func(env *Env) { env.A = Any(UInt64(atouint64(env.Get(0)) - atouint64(env.Get(1)))) })).
				Put(M__mul, NativeFun(func(env *Env) { env.A = Any(UInt64(atouint64(env.Get(0)) * atouint64(env.Get(1)))) })).
				Put(M__div, NativeFun(func(env *Env) { env.A = Any(UInt64(atouint64(env.Get(0)) / atouint64(env.Get(1)))) })).
				Put(M__unm, NativeFun(func(env *Env) { env.A = Any(UInt64(-atouint64(env.Get(0)))) })).
				Put(M__eq, NativeFun(func(env *Env) { env.A = Bln(atouint64(env.Get(0)) == atouint64(env.Get(1))) })).
				Put(M__lt, NativeFun(func(env *Env) { env.A = Bln(atouint64(env.Get(0)) < atouint64(env.Get(1))) })).
				Put(M__le, NativeFun(func(env *Env) { env.A = Bln(atouint64(env.Get(0)) <= atouint64(env.Get(1))) })).
				Put(M__tostring, NativeFun(func(env *Env) { env.A = Str(strconv.FormatUint(atouint64(env.Get(0)), 10)) }))
)

func atoint64(v Value) int64 {
	switch v.Type() {
	case NUM:
		return int64(v.Num())
	case STR:
		v, _ := strconv.ParseInt(v.Str(), 0, 64)
		return v
	case ANY:
		switch v := v.Any().(type) {
		case Int64:
			return int64(v)
		case UInt64:
			return int64(v)
		}
	}
	v.ExpectMsg(NUM, "Int64")
	return 0
}

func atouint64(v Value) uint64 {
	switch v.Type() {
	case NUM:
		return uint64(v.Num())
	case STR:
		v, _ := strconv.ParseUint(v.Str(), 0, 64)
		return v
	case ANY:
		switch v := v.Any().(type) {
		case Int64:
			return uint64(v)
		case UInt64:
			return uint64(v)
		}
	}
	v.ExpectMsg(NUM, "UInt64")
	return 0
}

func (a Bytes) GetMetatable() *Table    { return nativeBytesMetatable }
func (a Bytes) SetMetatable(mt *Table)  { nativeBytesMetatable = mt }
func (a Int64) GetMetatable() *Table    { return nativeInt64Metatable }
func (a Int64) SetMetatable(mt *Table)  { nativeInt64Metatable = mt }
func (a UInt64) GetMetatable() *Table   { return nativeUInt64Metatable }
func (a UInt64) SetMetatable(mt *Table) { nativeUInt64Metatable = mt }
