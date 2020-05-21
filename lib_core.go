package potatolang

import (
	"runtime"
	"strconv"
)

var CoreLibs = map[string]Value{}

// AddCoreValue adds a value to the core libraries
// duplicated name will result in panicking
func AddCoreValue(name string, value Value) {
	if name == "" {
		return
	}
	if CoreLibs[name].Type() != NIL {
		panicf("core value %s already exists", name)
	}
	CoreLibs[name] = value
}

func initCoreLibs() {
	lclosure := &Table{}
	lclosure.Puts("copy", NewNativeValue(1, false, func(env *Env) {
		cls := env.In(0, FUN).Fun().Dup()
		env.A = Fun(cls)
	}), false)
	// 	lcore.Put("Safe", NewNativeValue(1, func(env *Env) Value {
	// 		cls := env.LocalGet(0).MustClosure()
	// 		cls.Set(ClsRecoverable)
	// 		return Fun(cls)
	// 	}))
	// 	lcore.Put("Eval", NewNativeValue(1, func(env *Env) Value {
	// 		env.B = Value{}
	// 		cls, err := LoadString(string(env.LocalGet(0).MustString()))
	// 		if err != nil {
	// 			env.B = Str(err.Error())
	// 			return Value{}
	// 		}
	// 		return Fun(cls)
	// 	}))
	// 	lcore.Put("Unicode", NewNativeValue(1, func(env *Env) Value {
	// 		return Str(string(rune(env.LocalGet(0).MustNumber())))
	// 	}))
	// 	lcore.Put("Char", NewNativeValue(1, func(env *Env) Value {
	// 		r, _ := utf8.DecodeRuneInString(env.LocalGet(0).MustString())
	// 		return Num(float64(r))
	// 	}))
	// 	lcore.Put("Index", NewNativeValue(2, func(env *Env) Value {
	// 		x := env.LocalGet(1)
	// 		for i, a := range env.LocalGet(0).MustSlice().l {
	// 			if a.Equal(x) {
	// 				return Num(float64(i))
	// 			}
	// 		}
	// 		return Num(-1)
	// 	}))
	// 	lcore.Put("PopBack", NewNativeValue(2, func(env *Env) Value {
	// 		s := env.LocalGet(0).MustSlice()
	// 		if len(s.l) == 0 {
	// 			env.B = Value{}
	// 			return Value{}
	// 		}
	// 		res := s.l[len(s.l)-1]
	// 		s.l = s.l[:len(s.l)-1]
	// 		env.B = NewSliceValue(s)
	// 		return res
	// 	}))
	// 	lcore.Put("sync", NewStructValue(NewStruct().
	// 		Put("mutex", NewNativeValue(0, func(env *Env) Value {
	// 			m, mux := NewStruct(), &sync.Mutex{}
	// 			m.Put("lock", NewNativeValue(0, func(env *Env) Value { mux.Lock(); return Value{} }))
	// 			m.Put("unlock", NewNativeValue(0, func(env *Env) Value { mux.Unlock(); return Value{} }))
	// 			return NewStructValue(m)
	// 		})).
	// 		Put("waitgroup", NewNativeValue(0, func(env *Env) Value {
	// 			m, wg := NewStruct(), &sync.WaitGroup{}
	// 			m.Put("add", NewNativeValue(1, func(env *Env) Value { wg.Add(int(env.LocalGet(0).MustNumber())); return Value{} }))
	// 			m.Put("done", NewNativeValue(0, func(env *Env) Value { wg.Done(); return Value{} }))
	// 			m.Put("wait", NewNativeValue(0, func(env *Env) Value { wg.Wait(); return Value{} }))
	// 			return NewStructValue(m)
	// 		}))))
	// 	CoreLibs["std"] = NewStructValue(lcore)
	CoreLibs["type"] = NewNativeValue(1, false, func(env *Env) {
		env.A = Str(typeMappings[env.LocalGet(0).Type()])
	})
	CoreLibs["rawset"] = NewNativeValue(3, false, func(env *Env) {
		env.In(0, TAB).Tab().Put(env.LocalGet(1), env.LocalGet(2), true)
	})
	CoreLibs["rawget"] = NewNativeValue(2, false, func(env *Env) {
		env.A = env.In(0, TAB).Tab().Get(env.LocalGet(1), true)
	})
	CoreLibs["pcall"] = NewNativeValue(1, true, func(env *Env) {
		defer func() {
			if r := recover(); r != nil {
				env.A, env.B = Value{}, Value{}
			}
		}()
		env.A, env.B = env.In(0, FUN).Fun().Call(env.Vararg...)
	})
	CoreLibs["closure"] = Tab(lclosure)
	CoreLibs["unpack"] = NewNativeValue(1, true, func(env *Env) {
		a := env.In(0, TAB).Tab().a
		start, end := 0, len(a)
		if len(env.Vararg) > 0 {
			start = int(env.Vararg[0].Expect(NUM).Num())
		}
		if len(env.Vararg) > 1 {
			end = int(env.Vararg[1].Expect(NUM).Num())
		}
		env.A = newUnpackedValue(a[start:end])
	})
	CoreLibs["setmetatable"] = NewNativeValue(2, false, func(env *Env) {
		t := env.In(0, TAB).Tab()
		if t.mt != nil && !t.mt.Gets("__metatable", false).IsNil() {
			panicf("cannot change protected metatable")
		}
		if env.LocalGet(1).IsNil() {
			t.mt = nil
		} else {
			t.mt = env.In(1, TAB).Tab()
		}
		env.A = env.LocalGet(0)
	})
	CoreLibs["getmetatable"] = NewNativeValue(0, true, func(env *Env) {
		if len(env.Vararg) == 0 {
			env.A = Value{}
			return
		}
		if t := env.Vararg[0]; t.Type() == TAB {
			if t := t.Tab().mt; t != nil {
				if _mt := t.Gets("__metatable", false); !_mt.IsNil() {
					env.A = _mt
				} else {
					env.A = Tab(t)
				}
			}
		} else {
			env.A = Value{}
		}
	})
	CoreLibs["assert"] = NewNativeValue(1, false, func(env *Env) {
		if v := env.LocalGet(0); !v.IsFalse() {
			return
		}
		panic("assertion failed")
	})
	CoreLibs["pairs"] = NewNativeValue(1, false, func(env *Env) {
		iter := env.In(0, TAB).Tab().Iter()
		env.A = NewNativeValue(0, false, func(env *Env) {
			if !iter.Next() {
				env.A = Value{}
			} else {
				env.A, env.B = iter.Key(), iter.Value()
			}
		})
	})
	CoreLibs["ipairs"] = NewNativeValue(1, false, func(env *Env) {
		t := env.In(0, TAB).Tab()
		idx := -1
		env.A = NewNativeValue(0, false, func(env *Env) {
		AGAIN:
			idx++
			if idx >= len(t.a) {
				env.A, env.B = Value{}, Value{}
			} else {
				if t.a[idx].IsNil() {
					goto AGAIN
				}
				env.A, env.B = Num(float64(idx)), t.a[idx]
			}
		})
	})
	CoreLibs["tostring"] = NewNativeValue(1, false, func(env *Env) {
		v := env.LocalGet(0)
		if v.Type() == TAB && v.Tab().mt != nil {
			_tostring := v.Tab().mt.Gets("__tostring", false)
			if _tostring.Type() == FUN {
				env.A, _ = _tostring.Fun().Call(v)
				return
			}
		}
		env.A = Str(v.String())
	})
	CoreLibs["tonumber"] = NewNativeValue(1, false, func(env *Env) {
		v := env.LocalGet(0)
		switch v.Type() {
		case NUM:
			env.A = v
		case STR:
			v, _ := strconv.ParseFloat(v.Str(), 64)
			env.A = Num(v)
		default:
			env.A = Value{}
		}
	})
	CoreLibs["collectgarbage"] = NewNativeValue(0, false, func(env *Env) {
		runtime.GC()
	})
	// 	CoreLibs["copy"] = NewNativeValue(2, func(env *Env) Value {
	// 		if env.LocalSize() == 2 {
	// 			switch v := env.LocalGet(1); v.Type() {
	// 			case STR:
	// 				arr := env.LocalGet(0).MustSlice().l
	// 				str := v.Str()
	// 				n := 0
	// 				for i := range arr {
	// 					if n >= len(str) {
	// 						break
	// 					}
	// 					arr[i] = Num(float64(str[n]))
	// 					n++
	// 				}
	// 				return Num(float64(len(arr)))
	// 			default:
	// 				return Num(float64(copy(env.LocalGet(0).MustSlice().l, v.MustSlice().l)))
	// 			}
	// 		}
	// 		return env.LocalGet(0).Dup()
	// 	})
	// 	CoreLibs["go"] = NewNativeValue(1, func(env *Env) Value {
	// 		cls := env.LocalGet(0).MustClosure()
	// 		newEnv := NewEnv(cls.Env)
	// 		newEnv.stack = append([]Value{}, env.stack[1:]...)
	// 		go cls.Exec(newEnv)
	// 		return Value{}
	// 	})
	// 	CoreLibs["make"] = NewNativeValue(1, func(env *Env) Value {
	// 		return NewSliceValue(NewSliceSize(int(env.LocalGet(0).MustNumber())))
	// 	})
	//
	// 	// chanDefault := NewPointerValue(new(int))
	// 	//CoreLibs["chan"] = NewStructValue(NewStruct().
	// 	//	Put("Default", chanDefault).
	// 	//	Put("Make", NewNativeValue(1, func(env *Env) Value {
	// 	//		ch := make(chan Value, int(env.LocalGet(0).MustNumber()))
	// 	//		return NewPointerValue(ch)
	// 	//	})).
	// 	// Put("Send", NewNativeValue(2, func(env *Env) Value {
	// 	// 	p := env.LocalGet(0).Any().(*chan Value)
	// 	// 	*p <- env.LocalGet(1)
	// 	// 	return env.LocalGet(1)
	// 	// })).
	// 	// Put("Recv", NewNativeValue(1, func(env *Env) Value {
	// 	// 	p := (*chan Value)(env.LocalGet(0).MustPointer(PTagChan))
	// 	// 	return <-*p
	// 	// })).
	// 	// Put("Close", NewNativeValue(1, func(env *Env) Value {
	// 	// 	close(*(*chan Value)(env.LocalGet(0).MustPointer(PTagChan)))
	// 	// 	return Value{}
	// 	// })).
	// 	// Put("Select", NewNativeValue(0, func(env *Env) Value {
	// 	// 	cases := make([]reflect.SelectCase, env.LocalSize())
	// 	// 	chans := make([]chan Value, len(cases))
	// 	// 	for i := range chans {
	// 	// 		if a := env.LocalGet(i); a == chanDefault {
	// 	// 			cases[i] = reflect.SelectCase{Dir: reflect.SelectDefault}
	// 	// 		} else {
	// 	// 			p := (*chan Value)(a.MustPointer(PTagChan))
	// 	// 			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*p)}
	// 	// 			chans[i] = *p
	// 	// 		}
	// 	// 	}
	// 	// 	chosen, value, _ := reflect.Select(cases)
	// 	// 	v, ch := Value{}, NewPointerValue(unsafe.Pointer(&chans[chosen]), PTagChan)
	// 	// 	if value.IsValid() {
	// 	// 		v, _ = value.Interface().(Value)
	// 	// 	} else {
	// 	// 		ch = Value{}
	// 	// 	}
	// 	// 	env.B = ch
	// 	// 	return v
	// 	// })))
	//
	// 	CoreLibs["map"] = NewNativeValue(0, func(env *Env) Value {
	// 		var m map[string]Value
	// 		if env.LocalSize() == 1 {
	// 			switch a := env.LocalGet(0); a.Type() {
	// 			case NUM:
	// 				m = make(map[string]Value, int(a.Num()))
	// 			case StructType:
	// 				s := a.AsStruct().l
	// 				m = make(map[string]Value, len(s)/2)
	// 				for i, v := range s[:len(s)/2] {
	// 					m[string(hash50.FindStringHash(v.Num()))] = s[i+len(s)/2]
	// 				}
	// 			default:
	// 				a.testType(NUM)
	// 			}
	// 			// 	} else if env.LocalSize() == 2 {
	// 			// 		a, b := env.LocalGet(0), env.LocalGet(1)
	// 			// 		if a.Type()+b.Type() == SliceType*2 {
	//
	// 			// 		} else if a.Type() == SliceType && b.Type() == FUN {
	//
	// 			// 		} else {
	// 			// 			a.testType(NUM)
	// 			// 		}
	// 		} else {
	// 			m = make(map[string]Value)
	// 		}
	// 		return NewStructValue(NewStruct().
	// 			Put("Get", NewNativeValue(1, func(env *Env) Value {
	// 				buf := env.LocalGet(0).MustString()
	// 				v, ok := m[*(*string)(unsafe.Pointer(&buf))]
	// 				env.B = Bln(ok)
	// 				return v
	// 			})).
	// 			Put("Put", NewNativeValue(2, func(env *Env) Value {
	// 				buf := env.LocalGet(0).MustString()
	// 				v := env.LocalGet(1)
	// 				m[*(*string)(unsafe.Pointer(&buf))] = v
	// 				return v
	// 			})).
	// 			Put("Len", NewNativeValue(1, func(env *Env) Value {
	// 				return Num(float64(len(m)))
	// 			})).
	// 			Put("Delete", NewNativeValue(1, func(env *Env) Value {
	// 				buf := env.LocalGet(0).MustString()
	// 				v, ok := m[*(*string)(unsafe.Pointer(&buf))]
	// 				env.B = Bln(ok)
	// 				delete(m, *(*string)(unsafe.Pointer(&buf)))
	// 				return v
	// 			})).
	// 			Put("Range", NewNativeValue(1, func(env *Env) Value {
	// 				cls := env.LocalGet(0).MustClosure()
	// 				newEnv := NewEnv(env)
	// 				for k, v := range m {
	// 					newEnv.LocalClear()
	// 					newEnv.LocalPush(Str(k))
	// 					newEnv.LocalPush(v)
	// 					ok, _ := cls.Exec(newEnv)
	// 					if ok.IsZero() {
	// 						break
	// 					}
	// 				}
	// 				return Value{}
	// 			})))
	// 	})
	//
	initLibAux()
	// 	initLibMath()
}
