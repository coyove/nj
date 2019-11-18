package potatolang

import (
	"math"
	"reflect"
	"strconv"
	"sync"
	"unicode/utf8"
	"unsafe"

	"github.com/coyove/potatolang/hash50"
	"github.com/coyove/potatolang/parser"
)

const (
	_ = iota
	PTagUnique
	PTagPhantom
	PTagChan
	PTagMap
	PTagInterface
)

var CoreLibs = map[string]Value{}

// AddCoreValue adds a value to the core libraries
// duplicated name will result in panicking
func AddCoreValue(name string, value Value) {
	if name == "" {
		return
	}
	if CoreLibs[name].Type() != NilType {
		panicf("core value %s already exists", name)
	}
	CoreLibs[name] = value
}

func panicerr(err error) {
	if err != nil {
		panic(err)
	}
}

func initCoreLibs() {
	lcore := NewStruct()
	lcore.Put("Unique", NewNativeValue(0, func(env *Env) Value {
		a := new(int)
		return NewPointerValue(unsafe.Pointer(a), PTagUnique)
	}))
	lcore.Put("Safe", NewNativeValue(1, func(env *Env) Value {
		cls := env.LocalGet(0).MustClosure()
		cls.Set(ClsRecoverable)
		return NewClosureValue(cls)
	}))
	lcore.Put("Eval", NewNativeValue(1, func(env *Env) Value {
		env.B = Value{}
		cls, err := LoadString(string(env.LocalGet(0).MustString()))
		if err != nil {
			env.B = NewStringValueString(err.Error())
			return Value{}
		}
		return NewClosureValue(cls)
	}))
	lcore.Put("Unicode", NewNativeValue(1, func(env *Env) Value {
		return NewStringValueString(string(rune(env.LocalGet(0).MustNumber())))
	}))
	lcore.Put("Char", NewNativeValue(1, func(env *Env) Value {
		r, _ := utf8.DecodeRune(env.LocalGet(0).MustString())
		return NewNumberValue(float64(r))
	}))
	lcore.Put("Index", NewNativeValue(2, func(env *Env) Value {
		x := env.LocalGet(1)
		for i, a := range env.LocalGet(0).MustSlice().l {
			if a.Equal(x) {
				return NewNumberValue(float64(i))
			}
		}
		return NewNumberValue(-1)
	}))
	lcore.Put("PopBack", NewNativeValue(2, func(env *Env) Value {
		s := env.LocalGet(0).MustSlice()
		if len(s.l) == 0 {
			env.B = Value{}
			return Value{}
		}
		res := s.l[len(s.l)-1]
		s.l = s.l[:len(s.l)-1]
		env.B = NewSliceValue(s)
		return res
	}))
	lcore.Put("sync", NewStructValue(NewStruct().
		Put("mutex", NewNativeValue(0, func(env *Env) Value {
			m, mux := NewStruct(), &sync.Mutex{}
			m.Put("lock", NewNativeValue(0, func(env *Env) Value { mux.Lock(); return Value{} }))
			m.Put("unlock", NewNativeValue(0, func(env *Env) Value { mux.Unlock(); return Value{} }))
			return NewStructValue(m)
		})).
		Put("waitgroup", NewNativeValue(0, func(env *Env) Value {
			m, wg := NewStruct(), &sync.WaitGroup{}
			m.Put("add", NewNativeValue(1, func(env *Env) Value { wg.Add(int(env.LocalGet(0).MustNumber())); return Value{} }))
			m.Put("done", NewNativeValue(0, func(env *Env) Value { wg.Done(); return Value{} }))
			m.Put("wait", NewNativeValue(0, func(env *Env) Value { wg.Wait(); return Value{} }))
			return NewStructValue(m)
		}))))
	CoreLibs["std"] = NewStructValue(lcore)
	CoreLibs["atoi"] = NewNativeValue(1, func(env *Env) Value {
		env.B = Value{}
		v, err := parser.StringToNumber(string(env.LocalGet(0).MustString()))
		if err != nil {
			env.B = NewStringValueString(err.Error())
			return Value{}
		}
		return NewNumberValue(v)
	})
	CoreLibs["itoa"] = NewNativeValue(1, func(env *Env) Value {
		v := env.LocalGet(0).MustNumber()
		base := 10
		if env.LocalSize() >= 2 {
			base = int(env.LocalGet(1).MustNumber())
		}
		return NewStringValueString(strconv.FormatInt(int64(v), base))
	})
	CoreLibs["assert"] = NewNativeValue(1, func(env *Env) Value {
		if v := env.LocalGet(0); !v.IsFalse() {
			return v
		}
		if env.LocalSize() > 1 {
			panic(env.LocalGet(1))
		}
		panic("assertion failed")
	})
	CoreLibs["append"] = NewNativeValue(1, func(env *Env) Value {
		switch v := env.LocalGet(0); v.Type() {
		case SliceType:
			s := NewSlice()
			s.l = append(v.AsSlice().l, env.stack[1:]...)
			return NewSliceValue(s)
		case StringType:
			x := append([]byte{}, v.AsString()...)
			for _, a := range env.stack[1:] {
				switch a.Type() {
				case NumberType:
					x = append(x, byte(a.AsNumber()))
				case StringType:
					x = append(x, a.AsString()...)
				default:
					v.testType(SliceType)
				}
			}
			return NewStringValue(x)
		default:
			v.testType(SliceType)
			return Value{}
		}
	})
	CoreLibs["copy"] = NewNativeValue(2, func(env *Env) Value {
		if env.LocalSize() == 2 {
			switch v := env.LocalGet(1); v.Type() {
			case StringType:
				arr := env.LocalGet(0).MustSlice().l
				str := v.AsString()
				n := 0
				for i := range arr {
					if n >= len(str) {
						break
					}
					arr[i] = NewNumberValue(float64(str[n]))
					n++
				}
				return NewNumberValue(float64(len(arr)))
			default:
				return NewNumberValue(float64(copy(env.LocalGet(0).MustSlice().l, v.MustSlice().l)))
			}
		}
		return env.LocalGet(0).Dup()
	})
	CoreLibs["go"] = NewNativeValue(1, func(env *Env) Value {
		cls := env.LocalGet(0).MustClosure()
		newEnv := NewEnv(cls.Env)
		newEnv.stack = append([]Value{}, env.stack[1:]...)
		go cls.Exec(newEnv)
		return Value{}
	})
	CoreLibs["make"] = NewNativeValue(1, func(env *Env) Value {
		return NewSliceValue(NewSliceSize(int(env.LocalGet(0).MustNumber())))
	})

	chanDefault := NewPointerValue(unsafe.Pointer(new(int)), PTagUnique)
	CoreLibs["chan"] = NewStructValue(NewStruct().
		Put("Default", chanDefault).
		Put("Make", NewNativeValue(1, func(env *Env) Value {
			ch := make(chan Value, int(env.LocalGet(0).MustNumber()))
			return NewPointerValue(unsafe.Pointer(&ch), PTagChan)
		})).
		Put("Send", NewNativeValue(2, func(env *Env) Value {
			p := (*chan Value)(env.LocalGet(0).MustPointer(PTagChan))
			*p <- env.LocalGet(1)
			return env.LocalGet(1)
		})).
		Put("Recv", NewNativeValue(1, func(env *Env) Value {
			p := (*chan Value)(env.LocalGet(0).MustPointer(PTagChan))
			return <-*p
		})).
		Put("Close", NewNativeValue(1, func(env *Env) Value {
			close(*(*chan Value)(env.LocalGet(0).MustPointer(PTagChan)))
			return Value{}
		})).
		Put("Select", NewNativeValue(0, func(env *Env) Value {
			cases := make([]reflect.SelectCase, env.LocalSize())
			chans := make([]chan Value, len(cases))
			for i := range chans {
				if a := env.LocalGet(i); a == chanDefault {
					cases[i] = reflect.SelectCase{Dir: reflect.SelectDefault}
				} else {
					p := (*chan Value)(a.MustPointer(PTagChan))
					cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*p)}
					chans[i] = *p
				}
			}
			chosen, value, _ := reflect.Select(cases)
			v, ch := Value{}, NewPointerValue(unsafe.Pointer(&chans[chosen]), PTagChan)
			if value.IsValid() {
				v, _ = value.Interface().(Value)
			} else {
				ch = Value{}
			}
			env.B = ch
			return v
		})))

	CoreLibs["map"] = NewNativeValue(0, func(env *Env) Value {
		var m map[string]Value
		if env.LocalSize() == 1 {
			switch a := env.LocalGet(0); a.Type() {
			case NumberType:
				m = make(map[string]Value, int(a.AsNumber()))
			case StructType:
				s := a.AsStruct().l
				m = make(map[string]Value, len(s)/2)
				for i, v := range s[:len(s)/2] {
					m[string(hash50.FindStringHash(v.AsNumber()))] = s[i+len(s)/2]
				}
			default:
				a.testType(NumberType)
			}
			// 	} else if env.LocalSize() == 2 {
			// 		a, b := env.LocalGet(0), env.LocalGet(1)
			// 		if a.Type()+b.Type() == SliceType*2 {

			// 		} else if a.Type() == SliceType && b.Type() == ClosureType {

			// 		} else {
			// 			a.testType(NumberType)
			// 		}
		} else {
			m = make(map[string]Value)
		}
		return NewStructValue(NewStruct().
			Put("Get", NewNativeValue(1, func(env *Env) Value {
				buf := env.LocalGet(0).MustString()
				v, ok := m[*(*string)(unsafe.Pointer(&buf))]
				env.B = NewBoolValue(ok)
				return v
			})).
			Put("Put", NewNativeValue(2, func(env *Env) Value {
				buf := env.LocalGet(0).MustString()
				v := env.LocalGet(1)
				m[*(*string)(unsafe.Pointer(&buf))] = v
				return v
			})).
			Put("Len", NewNativeValue(1, func(env *Env) Value {
				return NewNumberValue(float64(len(m)))
			})).
			Put("Delete", NewNativeValue(1, func(env *Env) Value {
				buf := env.LocalGet(0).MustString()
				v, ok := m[*(*string)(unsafe.Pointer(&buf))]
				env.B = NewBoolValue(ok)
				delete(m, *(*string)(unsafe.Pointer(&buf)))
				return v
			})).
			Put("Range", NewNativeValue(1, func(env *Env) Value {
				cls := env.LocalGet(0).MustClosure()
				newEnv := NewEnv(env)
				for k, v := range m {
					newEnv.LocalClear()
					newEnv.LocalPush(NewStringValueString(k))
					newEnv.LocalPush(v)
					ok, _ := cls.Exec(newEnv)
					if ok.IsZero() {
						break
					}
				}
				return Value{}
			})))
	})

	initLibAux()
	initLibMath()
}

func StorePointerUnsafe(ptr Value, v Value) {
	if ptr.IsFalse() {
		return
	}
	x := math.Float64bits(ptr.MustNumber())
	(*Env)(unsafe.Pointer(uintptr(x<<16>>16))).Set(uint16(x>>48), v)
}

func LoadPointerUnsafe(ptr Value) Value {
	if ptr.IsFalse() {
		return Value{}
	}
	x := math.Float64bits(ptr.MustNumber())
	return (*Env)(unsafe.Pointer(uintptr(x<<16>>16))).Get(uint16(x>>48), nil)
}
