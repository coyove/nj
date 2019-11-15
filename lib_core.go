package potatolang

import (
	"bytes"
	"math"
	"reflect"
	"strconv"
	"sync"
	"unicode/utf8"
	"unsafe"

	"github.com/coyove/potatolang/parser"
)

const (
	_ = iota
	PTagUnique
	PTagPhantom
	PTagChan
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
	lcore.Put("unique", NewNativeValue(0, func(env *Env) Value {
		a := new(int)
		return NewPointerValue(unsafe.Pointer(a), PTagUnique)
	}))
	lcore.Put("safe", NewNativeValue(1, func(env *Env) Value {
		cls := env.LocalGet(0).MustClosure()
		cls.Set(ClsRecoverable)
		return NewClosureValue(cls)
	}))
	lcore.Put("eval", NewNativeValue(1, func(env *Env) Value {
		cls, err := LoadString(string(env.LocalGet(0).MustString()))
		if err != nil {
			return NewStringValueString(err.Error())
		}
		return NewClosureValue(cls)
	}))
	lcore.Put("unicode", NewNativeValue(1, func(env *Env) Value {
		return NewStringValueString(string(rune(env.LocalGet(0).MustNumber())))
	}))
	lcore.Put("char", NewNativeValue(1, func(env *Env) Value {
		r, _ := utf8.DecodeRune(env.LocalGet(0).MustString())
		return NewNumberValue(float64(r))
	}))
	lcore.Put("index", NewNativeValue(2, func(env *Env) Value {
		switch s := env.LocalGet(0); s.Type() {
		case StringType:
			return NewNumberValue(float64(bytes.Index(s.AsString(), env.LocalGet(1).MustString())))
		case SliceType:
			m := s.AsSlice()
			x := env.LocalGet(1)
			for i, a := range m.l {
				if a.Equal(x) {
					return NewNumberValue(float64(i))
				}
			}
			return Value{}
		default:
			return NewNumberValue(-1)
		}
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
		v, err := parser.StringToNumber(string(env.LocalGet(0).MustString()))
		if err != nil {
			StorePointerUnsafe(env.LocalGet(1), NewStringValueString(err.Error()))
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
	CoreLibs["chan"] = NewStructValue(NewStruct().
		Put("make", NewNativeValue(1, func(env *Env) Value {
			ch := make(chan Value, int(env.LocalGet(0).MustNumber()))
			return NewPointerValue(unsafe.Pointer(&ch), PTagChan)
		})).
		Put("send", NewNativeValue(2, func(env *Env) Value {
			p := (*chan Value)(env.LocalGet(0).MustPointer(PTagChan))
			*p <- env.LocalGet(1)
			return env.LocalGet(1)
		})).
		Put("recv", NewNativeValue(1, func(env *Env) Value {
			p := (*chan Value)(env.LocalGet(0).MustPointer(PTagChan))
			return <-*p
		})).
		Put("close", NewNativeValue(1, func(env *Env) Value {
			close(*(*chan Value)(env.LocalGet(0).MustPointer(PTagChan)))
			return Value{}
		})).
		Put("select", NewNativeValue(0, func(env *Env) Value {
			cases := make([]reflect.SelectCase, env.LocalSize())
			chans := make([]chan Value, len(cases))
			for i := range chans {
				if a := env.LocalGet(i); a.Type() == StringType && bytes.Equal(a.AsString(), []byte("default")) {
					cases[i] = reflect.SelectCase{Dir: reflect.SelectDefault}
				} else {
					p := (*chan Value)(a.MustPointer(PTagChan))
					cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*p)}
					chans[i] = *p
				}
			}
			chosen, value, ok := reflect.Select(cases)
			v := Value{}
			if value.IsValid() {
				v, _ = value.Interface().(Value)
			}
			return NewStructValue(NewStruct().
				Put("ok", NewBoolValue(ok)).
				Put("value", v).
				Put("chan", NewPointerValue(unsafe.Pointer(&chans[chosen]), PTagChan)))
		})))

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
