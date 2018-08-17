package potatolang

import (
	"strings"
	"sync"
	"unsafe"

	"github.com/buger/jsonparser"
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

var CoreLibNames = []string{
	"std", "io", "math",
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
	CoreLibNames = append(CoreLibNames, name)
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
		newEnv := NewEnv(x.Cls().env)
		for _, v := range y.Map().l {
			newEnv.SPush(v)
		}
		return x.Cls().Exec(newEnv)
	}))
	lcore.Puts("id", NewNativeValue(1, func(env *Env) Value {
		return NewStringValue(env.SGet(0).hashstr())
	}))
	lcore.Puts("r0", NewNativeValue(1, func(env *Env) Value {
		env.parent.R0 = env.SGet(0)
		return Value{}
	}))
	lcore.Puts("r1", NewNativeValue(1, func(env *Env) Value {
		env.parent.R1 = env.SGet(0)
		return Value{}
	}))
	lcore.Puts("r2", NewNativeValue(1, func(env *Env) Value {
		env.parent.R2 = env.SGet(0)
		return Value{}
	}))
	lcore.Puts("r3", NewNativeValue(1, func(env *Env) Value {
		env.parent.R3 = env.SGet(0)
		return Value{}
	}))
	lcore.Puts("storeinto", NewNativeValue(3, func(env *Env) Value {
		e, x, y := env.SGet(0), env.SGet(1), env.SGet(2)
		ep, et := e.Gen()
		if et != GTagEnv {
			panicf("invalid generic tag: %d", et)
		}
		(*Env)(ep).Set(uint32(x.Num()), y)
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
		return env.SGet(0).Map().Remove(env.Get(1))
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
			return NewValue()
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
			return NewValue()
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
			for _, v := range m.m {
				if v[1].Equal(x) {
					return v[0]
				}
			}
			return NewValue()
		default:
			return NewNumberValue(-1)
		}
	}))
	lcore.Puts("sync", NewMapValue(NewMap().
		Puts("run", NewNativeValue(1, func(env *Env) Value {
			cls := env.SGet(0).Cls()
			newEnv := NewEnv(cls.env)
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
			return NewValue()
		})).
		Puts("mutex", NewNativeValue(0, func(env *Env) Value {
			m, mux := NewMap(), &sync.Mutex{}
			m.Puts("lock", NewNativeValue(0, func(env *Env) Value { mux.Lock(); return NewValue() }))
			m.Puts("unlock", NewNativeValue(0, func(env *Env) Value { mux.Unlock(); return NewValue() }))
			return NewMapValue(m)
		})).
		Puts("waitgroup", NewNativeValue(0, func(env *Env) Value {
			m, wg := NewMap(), &sync.WaitGroup{}
			m.Puts("add", NewNativeValue(1, func(env *Env) Value { wg.Add(int(env.SGet(0).Num())); return NewValue() }))
			m.Puts("done", NewNativeValue(0, func(env *Env) Value { wg.Done(); return NewValue() }))
			m.Puts("wait", NewNativeValue(0, func(env *Env) Value { wg.Wait(); return NewValue() }))
			return NewMapValue(m)
		}))))

	lcore.Puts("opcode", NewMapValue(NewMap().
		Puts("A", NewNumberValue(regA)).Puts("EOB", NewNumberValue(OP_EOB)).
		Puts("LOAD", NewNumberValue(OP_LOAD)).Puts("STORE", NewNumberValue(OP_STORE)).
		Puts("ADD", NewNumberValue(OP_ADD)).Puts("SUB", NewNumberValue(OP_SUB)).
		Puts("MUL", NewNumberValue(OP_MUL)).Puts("DIV", NewNumberValue(OP_DIV)).
		Puts("LESS", NewNumberValue(OP_LESS)).Puts("LESSEQ", NewNumberValue(OP_LESS_EQ)).
		Puts("IFNOT", NewNumberValue(OP_IFNOT)).Puts("IF", NewNumberValue(OP_IF)).
		Puts("CALL", NewNumberValue(OP_CALL)).Puts("JMP", NewNumberValue(OP_JMP)).
		Puts("PUSH", NewNumberValue(OP_PUSH)).Puts("PUSHK", NewNumberValue(OP_PUSHK)).
		Puts("RET", NewNumberValue(OP_RET)).Puts("RETK", NewNumberValue(OP_RETK)).
		Puts("YIELD", NewNumberValue(OP_YIELD)).Puts("YIELDK", NewNumberValue(OP_YIELDK)).
		Puts("R0", NewNumberValue(OP_R0)).Puts("R0K", NewNumberValue(OP_R0K)).
		Puts("R1", NewNumberValue(OP_R1)).Puts("R1K", NewNumberValue(OP_R1K)).
		Puts("R2", NewNumberValue(OP_R2)).Puts("R2K", NewNumberValue(OP_R2K)).
		Puts("R3", NewNumberValue(OP_R3)).Puts("R3K", NewNumberValue(OP_R3K)).
		Puts("R0R2", NewNumberValue(OP_R0R2)).Puts("R1R2", NewNumberValue(OP_R1R2)).
		Puts("SET", NewNumberValue(OP_SET)).Puts("SETK", NewNumberValue(OP_SETK)).
		Puts("closure", NewMapValue(NewMap().
			Puts("empty", NewNativeValue(0, func(env *Env) Value {
				cls := NewClosure(make([]uint64, 0), make([]Value, 0), env.parent, 0)
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
			})).
			Puts("write", NewNativeValue(4, func(env *Env) Value {
				cls := env.SGet(0).Cls()
				cls.code = append(cls.code, makeop(
					byte(env.SGet(1).Num()),
					uint32(env.SGet(2).Num()),
					uint32(env.SGet(3).Num()),
				))
				return Value{}
			})).
			Puts("writeconst", NewNativeValue(2, func(env *Env) Value {
				cls := env.SGet(0).Cls()
				cls.consts = append(cls.consts, env.SGet(1))
				return Value{}
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
		Puts("_", NewValue())))

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
			m.Put(NewStringValue(string(key)), NewValue())
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
			m.Put(NewNumberValue(i), NewValue())
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
