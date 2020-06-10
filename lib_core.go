package potatolang

import (
	"math"
	"runtime"
	"strconv"
)

var G = &Table{}

func initCoreLibs() {
	lclosure := &Table{}
	lclosure.Puts("copy", NativeFun(1, func(env *Env) {
		cls := env.In(0, FUN).Fun().Dup()
		env.A = Fun(cls)
	}))
	G.Puts("closure", Tab(lclosure))
	G.Puts("unpack", NativeFun(1, func(env *Env) {
		a := env.In(0, TAB).Tab().a
		start, end := 1, len(a)
		if len(env.V) > 0 {
			start = int(env.V[0].Expect(NUM).Num())
		}
		if len(env.V) > 1 {
			end = int(env.V[1].Expect(NUM).Num())
		}
		env.A = newUnpackedValue(a[start-1 : end])
	}))
	// 	lcore.Put("Eval", NativeFun(1, func(env *Env) Value {
	// 		env.V = Value{}
	// 		cls, err := LoadString(string(env.Get(0).MustString()))
	// 		if err != nil {
	// 			env.V = Str(err.Error())
	// 			return Value{}
	// 		}
	// 		return Fun(cls)
	// 	}))
	// 	lcore.Put("Unicode", NativeFun(1, func(env *Env) Value {
	// 		return Str(string(rune(env.Get(0).MustNumber())))
	// 	}))
	// 	lcore.Put("Char", NativeFun(1, func(env *Env) Value {
	// 		r, _ := utf8.DecodeRuneInString(env.Get(0).MustString())
	// 		return Num(float64(r))
	// 	}))
	// 	lcore.Put("Index", NativeFun(2, func(env *Env) Value {
	// 		x := env.Get(1)
	// 		for i, a := range env.Get(0).MustSlice().l {
	// 			if a.Equal(x) {
	// 				return Num(float64(i))
	// 			}
	// 		}
	// 		return Num(-1)
	// 	}))
	// 	lcore.Put("sync", NewStructValue(NewStruct().
	// 		Put("mutex", NativeFun(0, func(env *Env) Value {
	// 			m, mux := NewStruct(), &sync.Mutex{}
	// 			m.Put("lock", NativeFun(0, func(env *Env) Value { mux.Lock(); return Value{} }))
	// 			m.Put("unlock", NativeFun(0, func(env *Env) Value { mux.Unlock(); return Value{} }))
	// 			return NewStructValue(m)
	// 		})).
	// 		Put("waitgroup", NativeFun(0, func(env *Env) Value {
	// 			m, wg := NewStruct(), &sync.WaitGroup{}
	// 			m.Put("add", NativeFun(1, func(env *Env) Value { wg.Add(int(env.Get(0).MustNumber())); return Value{} }))
	// 			m.Put("done", NativeFun(0, func(env *Env) Value { wg.Done(); return Value{} }))
	// 			m.Put("wait", NativeFun(0, func(env *Env) Value { wg.Wait(); return Value{} }))
	// 			return NewStructValue(m)
	// 		}))))
	// 	G["std"] = NewStructValue(lcore)
	ltable := &Table{}
	ltable.Puts("insert", NativeFun(2, func(env *Env) {
		t := env.In(0, TAB).Tab()
		if len(env.V) > 0 {
			t.Insert(env.In(1, NUM), env.V[0])
			env.A = env.V[0]
		} else {
			t.Insert(Num(float64(t.Len())), env.V[0])
		}
	}))
	ltable.Puts("remove", NativeFun(1, func(env *Env) {
		t := env.In(0, TAB).Tab()
		n := t.Len()
		if len(env.V) > 0 {
			n = int(env.V[0].Expect(NUM).Num())
		}
		t.Remove(n)
	}))
	ltable.Puts("unpack", G.Get(Str("unpack")))
	G.Puts("table", Tab(ltable))
	G.Puts("type", NativeFun(1, func(env *Env) {
		env.A = Str(typeMappings[env.Get(0).Type()])
	}))
	G.Puts("rawset", NativeFun(3, func(env *Env) {
		env.In(0, TAB).Tab().RawPut(env.Get(1), env.Get(2))
	}))
	G.Puts("rawget", NativeFun(2, func(env *Env) {
		env.A = env.In(0, TAB).Tab().RawGet(env.Get(1))
	}))
	G.Puts("rawequal", NativeFun(2, func(env *Env) {
		switch v, r := env.Get(0), env.Get(1); v.Type() + r.Type() {
		case NumNum, BlnBln, NilNil:
			env.A = Bln(v == r)
		case StrStr:
			env.A = Bln(r.Str() == v.Str())
		case AnyAny:
			env.A = Bln(v.Any() == r.Any())
		case TabTab:
			env.A = Bln(v == r)
		case FunFun:
			env.A = Bln(v.Fun() == r.Fun())
		default:
			env.A = Bln(false)
		}
	}))
	G.Puts("next", NativeFun(1, func(env *Env) {
		k, v := env.In(0, TAB).Tab().m.Next(env.Get(1))
		env.Return(k, v)
	}))
	G.Puts("rawlen", NativeFun(1, func(env *Env) {
		switch env.A = env.Get(0); env.A.Type() {
		case TAB:
			env.A = Num(float64(env.A.Tab().Len()))
		case STR:
			env.A = Num(float64(len(env.A.Str())))
		default:
			env.A.ExpectMsg(TAB, "rawlen")
		}
	}))
	G.Puts("call", NativeFun(1, func(env *Env) {
		env.A, env.V = env.In(0, FUN).Fun().Call(env.V...)
	}))
	G.Puts("pcall", NativeFun(1, func(env *Env) {
		defer func() {
			if r := recover(); r != nil {
				env.A, env.V = Bln(false), nil
			}
		}()
		a, v := env.In(0, FUN).Fun().Call(env.V...)
		env.A, env.V = Bln(true), append([]Value{a}, v...)
	}))
	G.Puts("select", NativeFun(2, func(env *Env) {
		switch a := env.Get(0); a.Type() {
		case STR:
			env.A = Num(float64(len(env.In(1, UPK)._Upk())))
		case NUM:
			if u, idx := env.In(1, UPK)._Upk(), int(a.Num())-1; idx < len(u) {
				env.A, env.V = u[idx], u[idx+1:]
			} else {
				env.A, env.V = Value{}, nil
			}
		}
	}))
	G.Puts("pack", NativeFun(1, func(env *Env) {
		t := &Table{a: env.In(0, UPK)._Upk()}
		env.A = Tab(t)
	}))
	G.Puts("setmetatable", NativeFun(2, func(env *Env) {
		if !env.Get(0).GetMetatable().RawGet(M__metatable).IsNil() {
			panicf("cannot change protected metatable")
		}
		if env.Get(1).IsNil() {
			env.Get(0).SetMetatable(nil)
		} else {
			env.Get(0).SetMetatable(env.In(1, TAB).Tab())
		}
		env.A = env.Get(0)
	}))
	G.Puts("getmetatable", NativeFun(1, func(env *Env) {
		t := env.Get(0).GetMetatable()
		if mt := t.RawGet(M__metatable); !mt.IsNil() {
			env.A = mt
		} else {
			env.A = Tab(t)
		}
	}))
	G.Puts("assert", NativeFun(1, func(env *Env) {
		if v := env.Get(0); !v.IsFalse() {
			env.A = v
			return
		}
		panic("assertion failed")
	}))
	G.Puts("pairs", NativeFun(1, func(env *Env) {
		t := env.In(0, TAB).Tab()
		var idx = -1
		var lastk Value
		env.A = NativeFun(0, func(env *Env) {
		AGAIN:
			idx++
			if idx >= len(t.a) {
				k, v := t.m.Next(lastk)
				if k.IsNil() {
					env.A, env.V = Value{}, nil
				} else {
					env.A, env.V = k, []Value{v}
					lastk = k
				}
			} else {
				if t.a[idx].IsNil() {
					goto AGAIN
				}
				env.A, env.V = Num(float64(idx)+1), []Value{t.a[idx]}
			}
		})
	}))
	G.Puts("ipairs", NativeFun(1, func(env *Env) {
		var arr Value
		switch t := env.Get(0); t.Type() {
		case TAB:
			arr = t
		case UPK:
			arr = Tab(&Table{a: t._Upk()})
		default:
			t.ExpectMsg(TAB, "ipairs")
		}
		env.A = NativeFun(2, func(env *Env) {
			idx := int(env.In(1, NUM).Num())
			arr := env.In(0, TAB).Tab().a
		AGAIN:
			idx++
			if idx > len(arr) {
				env.A, env.V = Value{}, nil
			} else {
				if arr[idx-1].IsNil() {
					goto AGAIN
				}
				env.A, env.V = Num(float64(idx)), []Value{arr[idx-1]}
			}
		})
		env.V = []Value{arr, Num(0)}
	}))
	G.Puts("tostring", NativeFun(1, func(env *Env) {
		v := env.Get(0)
		if f := v.GetMetamethod(M__tostring); f.Type() == FUN {
			env.A, env.V = f.Fun().Call(v)
			return
		}
		env.A = Str(v.String())
	}))
	G.Puts("tonumber", NativeFun(1, func(env *Env) {
		v := env.Get(0)
		switch v.Type() {
		case NUM:
			env.A = v
		case STR:
			v, _ := strconv.ParseFloat(v.Str(), 64)
			env.A = Num(v)
		default:
			env.A = Value{}
		}
	}))
	G.Puts("collectgarbage", NativeFun(0, func(env *Env) {
		runtime.GC()
	}))
	// 	G["go"] = NativeFun(1, func(env *Env) Value {
	// 		cls := env.Get(0).MustClosure()
	// 		newEnv := NewEnv(cls.Env)
	// 		newEnv.stack = append([]Value{}, env.stack[1:]...)
	// 		go cls.Exec(newEnv)
	// 		return Value{}
	// 	})
	// 	G["make"] = NativeFun(1, func(env *Env) Value {
	// 		return NewSliceValue(NewSliceSize(int(env.Get(0).MustNumber())))
	// 	})
	//
	// 	// chanDefault := Any(new(int))
	// 	//G["chan"] = NewStructValue(NewStruct().
	// 	//	Put("Default", chanDefault).
	// 	//	Put("Make", NativeFun(1, func(env *Env) Value {
	// 	//		ch := make(chan Value, int(env.Get(0).MustNumber()))
	// 	//		return Any(ch)
	// 	//	})).
	// 	// Put("Send", NativeFun(2, func(env *Env) Value {
	// 	// 	p := env.Get(0).Any().(*chan Value)
	// 	// 	*p <- env.Get(1)
	// 	// 	return env.Get(1)
	// 	// })).
	// 	// Put("Recv", NativeFun(1, func(env *Env) Value {
	// 	// 	p := (*chan Value)(env.Get(0).MustPointer(PTagChan))
	// 	// 	return <-*p
	// 	// })).
	// 	// Put("Close", NativeFun(1, func(env *Env) Value {
	// 	// 	close(*(*chan Value)(env.Get(0).MustPointer(PTagChan)))
	// 	// 	return Value{}
	// 	// })).
	// 	// Put("Select", NativeFun(0, func(env *Env) Value {
	// 	// 	cases := make([]reflect.SelectCase, env.Size())
	// 	// 	chans := make([]chan Value, len(cases))
	// 	// 	for i := range chans {
	// 	// 		if a := env.Get(i); a == chanDefault {
	// 	// 			cases[i] = reflect.SelectCase{Dir: reflect.SelectDefault}
	// 	// 		} else {
	// 	// 			p := (*chan Value)(a.MustPointer(PTagChan))
	// 	// 			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*p)}
	// 	// 			chans[i] = *p
	// 	// 		}
	// 	// 	}
	// 	// 	chosen, value, _ := reflect.Select(cases)
	// 	// 	v, ch := Value{}, Any(unsafe.Pointer(&chans[chosen]), PTagChan)
	// 	// 	if value.IsValid() {
	// 	// 		v, _ = value.Interface().(Value)
	// 	// 	} else {
	// 	// 		ch = Value{}
	// 	// 	}
	// 	// 	env.V = ch
	// 	// 	return v
	// 	// })))
	//
	initLibAux()
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
