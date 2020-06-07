package potatolang

import (
	"math"
	"unsafe"
)

func initLibMath() {
	//	r := rand.New()
	lmath := &Table{}
	//
	lmath.Puts("sqrt", NativeFun(1, func(env *Env) {
		env.A = Num(math.Sqrt(env.Get(0).Expect(NUM).Num()))
	}))
	lmath.Puts("floor", NativeFun(1, func(env *Env) {
		env.A = Num(math.Floor(env.Get(0).Expect(NUM).Num()))
	}))
	lmath.Puts("max", NativeFun(0, func(env *Env) {
		if len(env.V) == 0 {
			env.A = Value{}
		} else {
			max := env.V[0].Expect(NUM).Num()
			for i := 1; i < len(env.V); i++ {
				if x := env.V[i].Expect(NUM).Num(); x > max {
					max = x
				}
			}
			env.A = Num(max)
		}
	}))
	G.Puts("math", Tab(lmath))

	lnative := &Table{}
	lnative.Puts("bytes", NativeFun(1, func(env *Env) {
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
	}))

	G.Puts("native", Tab(lnative))

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
		Puts("__index", NativeFun(2, func(env *Env) {
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
		Puts("__newindex", NativeFun(3, func(env *Env) {
			a := env.In(0, ANY).Any().(NativeBytes)
			a[int(env.In(1, NUM).Num())] = byte(env.In(2, NUM).Num())
		})).
		Puts("__len", NativeFun(1, func(env *Env) {
			a := env.In(0, ANY).Any().(NativeBytes)
			env.A = Num(float64(len(a)))
		})).
		Puts("__concat", NativeFun(2, func(env *Env) {
			a := env.In(0, ANY).Any().(NativeBytes)
			b := env.In(1, ANY).Any().(NativeBytes)
			env.A = Any(NativeBytes(append(a, b...)))
		})).
		Puts("__eq", NativeFun(2, func(env *Env) {
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
