package potatolang

import (
	"sync"
	"unsafe"
)

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
	if CoreLibs[name].ty != Tnil {
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
		return NewGenericValue(unsafe.Pointer(a))
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
	lcore.Puts("storeinto", NewNativeValue(3, func(env *Env) Value {
		e, x, y := env.SGet(0), env.SGet(1), env.SGet(2)
		(*Env)(e.Gen()).Set(uint32(x.Num()), y)
		return y
	}))
	lcore.Puts("currentenv", NewNativeValue(0, func(env *Env) Value {
		return NewGenericValue(unsafe.Pointer(env.parent))
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
		dst, src := env.SGet(0).Map(), env.SGet(2).Map()
		dstPos, srcPos := int(env.SGet(1).Num()), int(env.SGet(3).Num())
		length := int(env.SGet(4).Num())
		return NewNumberValue(float64(copy(dst.l[dstPos:], src.l[srcPos:srcPos+length])))
	}))
	lcore.Puts("char", NewNativeValue(1, func(env *Env) Value {
		return NewStringValue(char(env.SGet(0).Num(), true))
	}))
	lcore.Puts("utf8char", NewNativeValue(1, func(env *Env) Value {
		return NewStringValue(char(env.SGet(0).Num(), false))
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

	CoreLibs["std"] = NewMapValue(lcore)

	initIOLib()
	initMathLib()
}
