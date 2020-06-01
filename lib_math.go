package potatolang

import (
	"math"
	"unsafe"
)

func initLibMath() {
	//	r := rand.New()
	lmath := &Table{}
	//
	lmath.Puts("sqrt", NativeFun(1, false, func(env *Env) {
		env.A = Num(math.Sqrt(env.Get(0).Expect(NUM).Num()))
	}), false)
	lmath.Puts("floor", NativeFun(1, false, func(env *Env) {
		env.A = Num(math.Floor(env.Get(0).Expect(NUM).Num()))
	}), false)
	lmath.Puts("max", NativeFun(0, true, func(env *Env) {
		if len(env.Vararg) == 0 {
			env.A = Value{}
		} else {
			max := env.Vararg[0].Expect(NUM).Num()
			for i := 1; i < len(env.Vararg); i++ {
				if x := env.Vararg[i].Expect(NUM).Num(); x > max {
					max = x
				}
			}
			env.A = Num(max)
		}
	}), false)
	G.Puts("math", Tab(lmath), false)

	lnative := &Table{}
	lnative.Puts("bytes", NativeFun(1, false, func(env *Env) {
		switch v := env.Get(0); v.Type() {
		case NUM:
			env.A = Any(make(NativeBytes, int(v.Num())))
		case STR:
			env.A = Any(NativeBytes(v.Str()))
		case ANY:
			env.A = Any(NativeBytes(append([]byte{}, v.Any().(NativeBytes)...)))
		default:
			env.A = Value{}
		}
	}), false)

	G.Puts("native", Tab(lnative), false)

	//	lmath.Put("rand", NewStructValue(NewStruct().
	//		Put("intn", NativeFun(1, func(env *Env) Value {
	//			return Num(float64(r.Intn(int(env.Get(0).MustNumber()))))
	//		})).
	//		Put("bytes", NativeFun(1, func(env *Env) Value {
	//			return Str(string(r.Fetch(int(env.Get(0).MustNumber()))))
	//		}))))
	//
}

type NativeBytes []byte

var (
	nativeBytesMetatable = (&Table{}).
		Puts("__index", NativeFun(2, false, func(env *Env) {
			a := env.In(0, ANY).Any().(NativeBytes)
			switch k := env.Get(1); k.Type() {
			case NUM:
				env.A = Num(float64(a[int(k.Num())]))
			case STR:
				switch k.Str() {
				case "append":
					env.A = NativeFun(1, true, func(env *Env) {
						a := env.In(0, ANY).Any().(NativeBytes)
						for _, v := range env.Vararg {
							a = append(a, byte(v.ExpectMsg(NUM, "append").Num()))
						}
						env.A = Any(a)
					})
					return
				case "tostring":
					env.A = NativeFun(1, false, func(env *Env) {
						a := env.In(0, ANY).Any().(NativeBytes)
						env.A = Str(string(a))
					})
					return
				}
				fallthrough
			default:
				panicf("invalid index: %#v", k)
			}
		}), false).
		Puts("__newindex", NativeFun(3, false, func(env *Env) {
			a := env.In(0, ANY).Any().(NativeBytes)
			a[int(env.In(1, NUM).Num())] = byte(env.In(2, NUM).Num())
		}), false).
		Puts("__len", NativeFun(1, false, func(env *Env) {
			a := env.In(0, ANY).Any().(NativeBytes)
			env.A = Num(float64(len(a)))
		}), false).
		Puts("__concat", NativeFun(2, false, func(env *Env) {
			a := env.In(0, ANY).Any().(NativeBytes)
			b := env.In(1, ANY).Any().(NativeBytes)
			env.A = Any(NativeBytes(append(a, b...)))
		}), false).
		Puts("__eq", NativeFun(2, false, func(env *Env) {
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
		}), false)
)

func (a NativeBytes) GetMetatable() *Table   { return nativeBytesMetatable }
func (a NativeBytes) SetMetatable(mt *Table) { nativeBytesMetatable = mt }
