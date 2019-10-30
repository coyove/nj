package potatolang

import (
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
	"unsafe"

	"github.com/buger/jsonparser"
	"github.com/coyove/potatolang/parser"
)

const (
	_ = iota
	PTagUnique
	PTagPhantom
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

func initCoreLibs() {
	lcore := NewMap()
	lcore.Puts("unique", NewNativeValue(0, func(env *Env) Value {
		a := new(int)
		return NewPointerValue(unsafe.Pointer(a), PTagUnique)
	}))
	lcore.Puts("genlist", NewNativeValue(1, func(env *Env) Value {
		return NewMapValue(NewMapSize(int(env.LocalGet(0).MustNumber())))
	}))
	lcore.Puts("apply", NewNativeValue(2, func(env *Env) Value {
		cls := env.LocalGet(0).MustClosure()
		newEnv := NewEnv(cls.Env)
		newEnv.stack = append([]Value{}, cls.PartialArgs...)
		for _, v := range env.LocalGet(1).MustMap().l {
			newEnv.LocalPush(v)
		}
		return cls.Exec(newEnv)
	}))
	lcore.Puts("stacktrace", NewNativeValue(0, func(env *Env) Value {
		panic("not implemented")
		//e := ExecError{stacks: Env.trace}
		//return NewStringValue(e.Error())
	}))
	lcore.Puts("eval", NewNativeValue(1, func(env *Env) Value {
		cls, err := LoadString(env.LocalGet(0).MustString())
		if err != nil {
			return NewStringValue(err.Error())
		}
		return NewClosureValue(cls)
	}))
	lcore.Puts("unicode", NewNativeValue(1, func(env *Env) Value {
		return NewStringValue(string(rune(env.LocalGet(0).MustNumber())))
	}))
	lcore.Puts("char", NewNativeValue(1, func(env *Env) Value {
		r, _ := utf8.DecodeRuneInString(env.LocalGet(0).MustString())
		return NewNumberValue(float64(r))
	}))
	lcore.Puts("index", NewNativeValue(2, func(env *Env) Value {
		switch s := env.LocalGet(0); s.Type() {
		case StringType:
			return NewNumberValue(float64(strings.Index(s.AsString(), env.LocalGet(1).MustString())))
		case MapType:
			m := s.AsMap()
			x := env.LocalGet(1)
			for i, a := range m.l {
				if a.Equal(x) {
					return NewNumberValue(float64(i))
				}
			}
			for k, v := range m.m {
				if v.Equal(x) {
					return NewInterfaceValue(k)
				}
			}
			return Value{}
		default:
			return NewNumberValue(-1)
		}
	}))
	lcore.Puts("sprintf", NewNativeValue(0, stdSprintf))
	lcore.Puts("ftoa", NewNativeValue(1, func(env *Env) Value {
		v := env.LocalGet(0).MustNumber()
		base := byte(env.LocalGet(1).MustNumber())
		digits := int(env.LocalGet(2).MustNumber())
		return NewStringValue(strconv.FormatFloat(v, byte(base), digits, 64))
	}))
	lcore.Puts("sync", NewMapValue(NewMap().
		Puts("run", NewNativeValue(1, func(env *Env) Value {
			cls := env.LocalGet(0).MustClosure()
			newEnv := NewEnv(cls.Env)
			if int(cls.ArgsCount) > env.LocalSize()-1 {
				panic("not enough arguments to start a goroutine")
			}
			newEnv.LocalPushFront(cls.PartialArgs)
			for i := 1; i < env.LocalSize(); i++ {
				newEnv.LocalPush(env.LocalGet(i))
			}
			go cls.Exec(newEnv)
			return Value{}
		})).
		Puts("mutex", NewNativeValue(0, func(env *Env) Value {
			m, mux := NewMap(), &sync.Mutex{}
			m.Puts("lock", NewNativeValue(0, func(env *Env) Value { mux.Lock(); return Value{} }))
			m.Puts("unlock", NewNativeValue(0, func(env *Env) Value { mux.Unlock(); return Value{} }))
			return NewMapValue(m)
		})).
		Puts("waitgroup", NewNativeValue(0, func(env *Env) Value {
			m, wg := NewMap(), &sync.WaitGroup{}
			m.Puts("add", NewNativeValue(1, func(env *Env) Value { wg.Add(int(env.LocalGet(0).MustNumber())); return Value{} }))
			m.Puts("done", NewNativeValue(0, func(env *Env) Value { wg.Done(); return Value{} }))
			m.Puts("wait", NewNativeValue(0, func(env *Env) Value { wg.Wait(); return Value{} }))
			return NewMapValue(m)
		}))))

	lcore.Puts("opcode", NewMapValue(NewMap().
		Puts("closure", NewMapValue(NewMap().
			Puts("empty", NewNativeValue(0, func(env *Env) Value {
				cls := NewClosure(make([]uint32, 0), make([]Value, 0), env.parent, 0)
				return NewClosureValue(cls)
			})).
			Puts("yieldreset", NewNativeValue(1, func(env *Env) Value {
				env.LocalGet(0).MustClosure().lastenv = nil
				return env.LocalGet(0)
			})).
			Puts("get", NewNativeValue(2, func(env *Env) Value {
				cls := env.LocalGet(0).MustClosure()
				switch name := env.LocalGet(1).MustString(); name {
				case "argscount":
					return NewNumberValue(float64(cls.ArgsCount))
				case "yieldable":
					return NewBoolValue(cls.Isset(ClsYieldable))
				case "envescaped":
					return NewBoolValue(!cls.Isset(ClsNoEnvescape))
				case "source":
					return NewStringValue(cls.source)
				}
				return NewClosureValue(cls)
			})))).
		Puts("_", Value{})))

	lcore.Puts("json", NewMapValue(NewMap().
		Puts("parse", NewNativeValue(1, func(env *Env) Value {
			json := []byte(strings.TrimSpace(env.LocalGet(0).MustString()))
			if len(json) == 0 {
				return Value{}
			}
			switch json[0] {
			case '[':
				return walkArray(json)
			case '{':
				return walkObject(json)
			case '"':
				str, err := jsonparser.ParseString(json)
				panicerr(err)
				return NewStringValue(str)
			case 't', 'f':
				b, err := jsonparser.ParseBoolean(json)
				panicerr(err)
				return NewBoolValue(b)
			default:
				num, err := jsonparser.ParseFloat(json)
				panicerr(err)
				return NewNumberValue(num)
			}
		})).
		Puts("stringify", NewNativeValue(1, func(env *Env) Value {
			return NewStringValue(env.LocalGet(0).toString(0, true))
		}))))

	CoreLibs["std"] = NewMapValue(lcore)
	CoreLibs["atoi"] = NewNativeValue(1, func(env *Env) Value {
		v, err := parser.StringToNumber(env.LocalGet(0).MustString())
		if err != nil {
			return Value{}
		}
		return NewNumberValue(v)
	})
	CoreLibs["itoa"] = NewNativeValue(1, func(env *Env) Value {
		v := env.LocalGet(0).MustNumber()
		if float64(int64(v)) == v {
			base := 10
			if env.LocalSize() >= 2 {
				base = int(env.LocalGet(1).MustNumber())
			}
			return NewStringValue(strconv.FormatInt(int64(v), base))
		}
		return NewStringValue(strconv.FormatFloat(v, 'f', -1, 64))
	})

	initIOLib()
	initMathLib()
}

func walkObject(buf []byte) Value {
	m := NewMap()
	jsonparser.ObjectEach(buf, func(key, value []byte, dataType jsonparser.ValueType, offset int) error {
		switch dataType {
		case jsonparser.Unknown:
			panic(value)
		case jsonparser.Null:
			m.Put(NewStringValue(string(key)), Value{})
		case jsonparser.Boolean:
			b, err := jsonparser.ParseBoolean(value)
			panicerr(err)
			m.Put(NewStringValue(string(key)), NewBoolValue(b))
		case jsonparser.Number:
			num, err := jsonparser.ParseFloat(value)
			panicerr(err)
			m.Put(NewStringValue(string(key)), NewNumberValue(num))
		case jsonparser.String:
			str, err := jsonparser.ParseString(value)
			panicerr(err)
			m.Put(NewStringValue(string(key)), NewStringValue(str))
		case jsonparser.Array:
			m.Put(NewStringValue(string(key)), walkArray(value))
		case jsonparser.Object:
			m.Put(NewStringValue(string(key)), walkObject(value))
		}
		return nil
	})
	return NewMapValue(m)
}

func walkArray(buf []byte) Value {
	m := NewMap()
	i := float64(0)
	jsonparser.ArrayEach(buf, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		switch dataType {
		case jsonparser.Unknown:
			panic(value)
		case jsonparser.Null:
			m.Put(NewNumberValue(i), Value{})
		case jsonparser.Boolean:
			b, err := jsonparser.ParseBoolean(value)
			panicerr(err)
			m.Put(NewNumberValue(i), NewBoolValue(b))
		case jsonparser.Number:
			num, err := jsonparser.ParseFloat(value)
			panicerr(err)
			m.Put(NewNumberValue(i), NewNumberValue(num))
		case jsonparser.String:
			str, err := jsonparser.ParseString(value)
			panicerr(err)
			m.Put(NewNumberValue(i), NewStringValue(str))
		case jsonparser.Array:
			m.Put(NewNumberValue(i), walkArray(value))
		case jsonparser.Object:
			m.Put(NewNumberValue(i), walkObject(value))
		}
		i++
	})
	return NewMapValue(m)
}
