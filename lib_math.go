package potatolang

import "math"

func initLibMath() {
	//	r := rand.New()
	lmath := &Table{}
	//
	lmath.Puts("sqrt", NewNativeValue(1, false, func(env *Env) {
		env.A = Num(math.Sqrt(env.Get(0).Expect(NUM).Num()))
	}), false)
	lmath.Puts("max", NewNativeValue(0, true, func(env *Env) {
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
	//	lmath.Put("rand", NewStructValue(NewStruct().
	//		Put("intn", NewNativeValue(1, func(env *Env) Value {
	//			return Num(float64(r.Intn(int(env.Get(0).MustNumber()))))
	//		})).
	//		Put("bytes", NewNativeValue(1, func(env *Env) Value {
	//			return Str(string(r.Fetch(int(env.Get(0).MustNumber()))))
	//		}))))
	//
	G.Puts("math", Tab(lmath), false)
}
