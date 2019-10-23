package potatolang

import (
	"strconv"
	"strings"
	"sync"
	"unsafe"

	"github.com/buger/jsonparser"
	"github.com/coyove/potatolang/parser"
)

const (
	GTagError = iota + 1
	GTagUnique
	GTagPhantom
	GTagEnv
	GTagByteArray
	GTagByteClampedArray
	GTagInt8Array
	GTagInt16Array
	GTagUint16Array
	GTagInt32Array
	GTagUint32Array
	GTagStringArray
	GTagBoolArray
	GTagFloat32Array
	GTagFloat64Array
)

type GTagComparator interface {
	Equal(a, b Value) bool
}

var gtagComparators = map[uint64]GTagComparator{}

func RegisterGTagcomparator(a, b uint32, comp GTagComparator) bool {
	x := uint64(a)<<32 + uint64(b)
	y := uint64(b)<<32 + uint64(a)
	if _, ok := gtagComparators[x]; ok {
		return false
	}
	gtagComparators[x] = comp
	gtagComparators[y] = comp
	return true
}

var CoreLibs = map[string]Value{}

// AddCoreValue adds a value to the core libraries
// duplicated name will result in panicking
func AddCoreValue(name string, value Value) {
	if name == "" {
		return
	}
	if CoreLibs[name].Type() != Tnil {
		panicf("core value %s already exists", name)
	}
	CoreLibs[name] = value
}

func char(v float64, ascii bool) string {
	if ascii {
		return string([]byte{byte(v)})
	}
	return string(rune(v))
}

func initCoreLibs() {
	lcore := NewMap()
	lcore.Puts("unique", NewNativeValue(0, func(env *Env) Value {
		a := new(int)
		return NewGenericValue(unsafe.Pointer(a), GTagUnique)
	}))
	lcore.Puts("genlist", NewNativeValue(1, func(env *Env) Value {
		return NewMapValue(NewMapSize(int(env.SGet(0).Num())))
	}))
	lcore.Puts("apply", NewNativeValue(2, func(env *Env) Value {
		x, y := env.SGet(0), env.SGet(1)
		newEnv := NewEnv(x.Cls().env, env.parent.Cancel)
		for _, v := range y.Map().l {
			newEnv.SPush(v)
		}
		return x.Cls().Exec(newEnv)
	}))
	lcore.Puts("storeinto", NewNativeValue(3, func(env *Env) Value {
		e, x, y := env.SGet(0), env.SGet(1), env.SGet(2)
		ep, et := e.Gen()
		if et != GTagEnv {
			panicf("invalid generic tag: %d", et)
		}
		(*Env)(ep).Set(uint16(x.Num()), y)
		return y
	}))
	lcore.Puts("currentenv", NewNativeValue(0, func(env *Env) Value {
		return NewGenericValue(unsafe.Pointer(env.parent), GTagEnv)
	}))
	lcore.Puts("stacktrace", NewNativeValue(0, func(env *Env) Value {
		e := ExecError{stacks: env.trace}
		return NewStringValue(e.Error())
	}))
	lcore.Puts("eval", NewNativeValue(1, func(env *Env) Value {
		x := env.SGet(0).Str()
		cls, err := LoadString(x)
		if err != nil {
			return NewStringValue(err.Error())
		}
		return NewClosureValue(cls)
	}))
	lcore.Puts("remove", NewNativeValue(2, func(env *Env) Value {
		return env.SGet(0).Map().Remove(env.Get(1, nil))
	}))
	lcore.Puts("copy", NewNativeValue(5, func(env *Env) Value {
		dstPos, srcPos := int(env.SGet(1).Num()), int(env.SGet(3).Num())
		length := int(env.SGet(4).Num())

		switch dst, src := env.SGet(0), env.SGet(2); dst.Type() {
		case Tmap:
			return NewNumberValue(float64(copy(dst.Map().l[dstPos:], src.Map().l[srcPos:srcPos+length])))
		case Tgeneric:
			return NewNumberValue(float64(GCopy(dst, src, dstPos, srcPos, srcPos+length)))
		default:
			panicf("can't copy from %+v to %+v", src, dst)
			return Value{}
		}
	}))
	lcore.Puts("char", NewNativeValue(1, func(env *Env) Value {
		switch c := env.SGet(0); c.Type() {
		case Tnumber:
			return NewStringValue(char(c.AsNumber(), true))
		case Tgeneric:
			return NewStringValue(string(*(*[]byte)(c.GenTags(GTagByteArray, GTagByteClampedArray, GTagInt8Array))))
		default:
			panicf("std.char: %+v", c)
			return Value{}
		}
	}))
	lcore.Puts("utf8char", NewNativeValue(1, func(env *Env) Value {
		return NewStringValue(char(env.SGet(0).Num(), false))
	}))
	lcore.Puts("index", NewNativeValue(2, func(env *Env) Value {
		switch s := env.SGet(0); s.Type() {
		case Tstring:
			return NewNumberValue(float64(strings.Index(s.AsString(), env.SGet(1).Str())))
		case Tmap:
			m := s.AsMap()
			x := env.SGet(1)
			for i, a := range m.l {
				if a.Equal(x) {
					return NewNumberValue(float64(i))
				}
			}
			for k, v := range m.m {
				if v.Equal(x) {
					return NewValueFromInterface(k)
				}
			}
			return Value{}
		default:
			return NewNumberValue(-1)
		}
	}))
	lcore.Puts("sprintf", NewNativeValue(0, stdSprintf))
	lcore.Puts("sync", NewMapValue(NewMap().
		Puts("run", NewNativeValue(1, func(env *Env) Value {
			cls := env.SGet(0).Cls()
			newEnv := NewEnv(cls.env, env.parent.Cancel)
			if cls.ArgsCount() > env.SSize()-1 {
				panic("not enough arguments to start a goroutine")
			}
			for i := 1; i < env.SSize(); i++ {
				newEnv.SPush(env.SGet(i))
			}
			if cls.Isset(CLS_HASRECEIVER) {
				newEnv.SPush(cls.caller)
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
			m.Puts("add", NewNativeValue(1, func(env *Env) Value { wg.Add(int(env.SGet(0).Num())); return Value{} }))
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
				env.SGet(0).Cls().lastenv = nil
				return env.SGet(0)
			})).
			Puts("set", NewNativeValue(3, func(env *Env) Value {
				cls := env.SGet(0).Cls()
				switch name := env.SGet(1).Str(); name {
				case "argscount":
					cls.argsCount = byte(env.SGet(2).Num())
				case "yieldable":
					if !env.SGet(2).IsFalse() {
						cls.Set(CLS_YIELDABLE)
					} else {
						cls.Unset(CLS_YIELDABLE)
					}
				case "envescaped":
					if env.SGet(2).IsFalse() {
						cls.Set(CLS_NOENVESCAPE)
					} else {
						cls.Unset(CLS_NOENVESCAPE)
					}
				case "source":
					cls.source = env.SGet(2).Str()
				}
				return NewClosureValue(cls)
			})).
			Puts("get", NewNativeValue(2, func(env *Env) Value {
				cls := env.SGet(0).Cls()
				switch name := env.SGet(1).Str(); name {
				case "argscount":
					return NewNumberValue(float64(cls.argsCount))
				case "yieldable":
					return NewBoolValue(cls.Isset(CLS_YIELDABLE))
				case "envescaped":
					return NewBoolValue(!cls.Isset(CLS_NOENVESCAPE))
				case "source":
					return NewStringValue(cls.source)
				}
				return NewClosureValue(cls)
			})))).
		Puts("_", Value{})))

	_genTyped := func(i interface{}, t uint32) Value {
		return NewGenericValueInterface(i, t)
	}
	lcore.Puts("typed", NewMapValue(NewMap().
		Puts("bytearray", NewNativeValue(1, func(env *Env) Value { return _genTyped(make([]byte, int(env.SGet(0).Num())), GTagByteArray) })).
		Puts("int8array", NewNativeValue(1, func(env *Env) Value { return _genTyped(make([]byte, int(env.SGet(0).Num())), GTagInt8Array) })).
		Puts("uint16array", NewNativeValue(1, func(env *Env) Value { return _genTyped(make([]byte, int(env.SGet(0).Num())), GTagUint16Array) })).
		Puts("int16array", NewNativeValue(1, func(env *Env) Value { return _genTyped(make([]byte, int(env.SGet(0).Num())), GTagInt16Array) })).
		Puts("uint32array", NewNativeValue(1, func(env *Env) Value { return _genTyped(make([]byte, int(env.SGet(0).Num())), GTagUint32Array) })).
		Puts("int32array", NewNativeValue(1, func(env *Env) Value { return _genTyped(make([]byte, int(env.SGet(0).Num())), GTagInt32Array) })).
		Puts("float32array", NewNativeValue(1, func(env *Env) Value { return _genTyped(make([]byte, int(env.SGet(0).Num())), GTagFloat32Array) })).
		Puts("float64array", NewNativeValue(1, func(env *Env) Value { return _genTyped(make([]byte, int(env.SGet(0).Num())), GTagFloat64Array) })).
		Puts("stringarray", NewNativeValue(1, func(env *Env) Value { return _genTyped(make([]byte, int(env.SGet(0).Num())), GTagStringArray) })).
		Puts("boolarray", NewNativeValue(1, func(env *Env) Value { return _genTyped(make([]byte, int(env.SGet(0).Num())), GTagBoolArray) })).
		Puts("_", Value{})))

	lcore.Puts("json", NewMapValue(NewMap().
		Puts("parse", NewNativeValue(1, func(env *Env) Value {
			json := []byte{}
			switch x := env.SGet(0); x.Type() {
			case Tstring:
				json = []byte(x.AsString())
			case Tgeneric:
				json = *(*[]byte)(x.GenTags(GTagByteArray, GTagByteClampedArray, GTagInt8Array))
			}
			for i := 0; i < len(json); i++ {
				switch json[i] {
				case '[':
					return walkArray(json[i:])
				case '{':
					return walkObject(json[i:])
				case '"':
					str, err := jsonparser.ParseString(json[i:])
					panicerr(err)
					return NewStringValue(str)
				case 't', 'f':
					b, err := jsonparser.ParseBoolean(json[i:])
					panicerr(err)
					return NewBoolValue(b)
				case ' ', '\t', '\r', '\n':
					// continue
				default:
					num, err := jsonparser.ParseFloat(json[i:])
					panicerr(err)
					return NewNumberValue(num)
				}
			}
			panic(json)
		})).
		Puts("stringify", NewNativeValue(1, func(env *Env) Value {
			return NewStringValue(env.SGet(0).toString(0, true))
		}))))

	CoreLibs["std"] = NewMapValue(lcore)
	CoreLibs["atoi"] = NewNativeValue(1, func(env *Env) Value {
		v, err := parser.StringToNumber(env.SGet(0).AsString())
		if err != nil {
			return Value{}
		}
		return NewNumberValue(v)
	})
	CoreLibs["itoa"] = NewNativeValue(1, func(env *Env) Value {
		v := env.SGet(0).AsNumber()
		if float64(int64(v)) == v {
			return NewStringValue(strconv.FormatInt(int64(v), 10))
		}
		return NewStringValue(strings.TrimRight(strconv.FormatFloat(v, 'f', 15, 64), "0"))
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
