package potatolang

import (
	"strconv"
	"unsafe"
)

type NativeBytes []byte
type NativeInt64 int64

var (
	nativeBytesMetatable = (&Table{}).
				Put(M__index, NativeFun(func(env *Env) {
			a := env.In(0, ANY).Any().(NativeBytes)
			switch k := env.Get(1); k.Type() {
			case NUM:
				env.A = Num(float64(a[int(k.Num())]))
			case STR:
				switch k.Str() {
				case "append":
					env.A = NativeFun(func(env *Env) {
						a := env.In(0, ANY).Any().(NativeBytes)
						for _, v := range env.stack[1:] {
							a = append(a, byte(v.ExpectMsg(NUM, "append").Num()))
						}
						env.A = Any(a)
					})
					return
				case "tostring":
					env.A = NativeFun(func(env *Env) {
						a := env.In(0, ANY).Any().(NativeBytes)
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
			a := env.In(0, ANY).Any().(NativeBytes)
			a[int(env.In(1, NUM).Num())] = byte(env.In(2, NUM).Num())
		})).
		Put(M__len, NativeFun(func(env *Env) {
			a := env.In(0, ANY).Any().(NativeBytes)
			env.A = Num(float64(len(a)))
		})).
		Put(M__concat, NativeFun(func(env *Env) {
			a := env.In(0, ANY).Any().(NativeBytes)
			b := env.In(1, ANY).Any().(NativeBytes)
			env.A = Any(NativeBytes(append(a, b...)))
		})).
		Put(M__eq, NativeFun(func(env *Env) {
			switch l, r := env.Get(0), env.Get(1); l.Type() + r.Type() {
			case AnyAny:
				a, b := l.Any().(NativeBytes), r.Any().(NativeBytes)
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
				var a NativeBytes
				var b string
				a, _ = l.Any().(NativeBytes)
				if a == nil {
					a, b = r.Any().(NativeBytes), l.Str()
				} else {
					b = r.Str()
				}
				env.A = Bln(*(*string)(unsafe.Pointer(&a)) == b)
			case ANY + TAB:
				var a NativeBytes
				var b *Table
				a, _ = l.Any().(NativeBytes)
				if a == nil {
					a, b = r.Any().(NativeBytes), l.Tab()
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
				Put(M__add, NativeFun(func(env *Env) { env.A = Any(NativeInt64(atoint64(env.Get(0)) + atoint64(env.Get(1)))) })).
				Put(M__sub, NativeFun(func(env *Env) { env.A = Any(NativeInt64(atoint64(env.Get(0)) - atoint64(env.Get(1)))) })).
				Put(M__mul, NativeFun(func(env *Env) { env.A = Any(NativeInt64(atoint64(env.Get(0)) * atoint64(env.Get(1)))) })).
				Put(M__div, NativeFun(func(env *Env) { env.A = Any(NativeInt64(atoint64(env.Get(0)) / atoint64(env.Get(1)))) })).
				Put(M__tostring, NativeFun(func(env *Env) { env.A = Str(strconv.FormatInt(atoint64(env.Get(0)), 10)) }))
)

func atoint64(v Value) int64 {
	switch v.Type() {
	case NUM:
		return int64(v.Num())
	case STR:
		v, _ := strconv.ParseInt(v.Str(), 0, 64)
		return v
	case ANY:
		return int64(v.Any().(NativeInt64))
	}
	panic("not a valid NativeInt64 object")
}

func (a NativeBytes) GetMetatable() *Table   { return nativeBytesMetatable }
func (a NativeBytes) SetMetatable(mt *Table) { nativeBytesMetatable = mt }
func (a NativeInt64) GetMetatable() *Table   { return nativeInt64Metatable }
func (a NativeInt64) SetMetatable(mt *Table) { nativeInt64Metatable = mt }
