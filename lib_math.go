package potatolang

import "unsafe"

type NativeBytes []byte

var (
	nativeBytesMetatable = (&Table{}).
		Put(M__index, NativeFun(2, func(env *Env) {
			a := env.In(0, ANY).Any().(NativeBytes)
			switch k := env.Get(1); k.Type() {
			case NUM:
				env.A = Num(float64(a[int(k.Num())]))
			case STR:
				switch k.Str() {
				case "append":
					env.A = NativeFun(1, func(env *Env) {
						a := env.In(0, ANY).Any().(NativeBytes)
						for _, v := range env.V {
							a = append(a, byte(v.ExpectMsg(NUM, "append").Num()))
						}
						env.A = Any(a)
					})
					return
				case "tostring":
					env.A = NativeFun(1, func(env *Env) {
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
		Put(M__newindex, NativeFun(3, func(env *Env) {
			a := env.In(0, ANY).Any().(NativeBytes)
			a[int(env.In(1, NUM).Num())] = byte(env.In(2, NUM).Num())
		})).
		Put(M__len, NativeFun(1, func(env *Env) {
			a := env.In(0, ANY).Any().(NativeBytes)
			env.A = Num(float64(len(a)))
		})).
		Put(M__concat, NativeFun(2, func(env *Env) {
			a := env.In(0, ANY).Any().(NativeBytes)
			b := env.In(1, ANY).Any().(NativeBytes)
			env.A = Any(NativeBytes(append(a, b...)))
		})).
		Put(M__eq, NativeFun(2, func(env *Env) {
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
)

func (a NativeBytes) GetMetatable() *Table   { return nativeBytesMetatable }
func (a NativeBytes) SetMetatable(mt *Table) { nativeBytesMetatable = mt }
