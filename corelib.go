package potatolang

import (
	"log"
	"sync"
	"unsafe"
)

var CoreLibNames = []string{
	"std", "io", "math",
}

var CoreLibs = map[string]Value{}

var cancelled = 1

func char(v float64, ascii bool) string {
	if ascii {
		return string([]byte{byte(v)})
	}
	return string(rune(v))
}

func initCoreLibs() {
	lcore := NewMap()
	lcore.Puts("cancelled", NewGenericValue(unsafe.Pointer(&cancelled)))
	lcore.Puts("genlist", NewNativeValue(1, func(env *Env) Value {
		v := env.SGet(0)
		if v.ty != Tnumber {
			v.panicType(Tnumber)
		}
		return NewMapValue(NewMapSize(int(v.AsNumber())))
	}))
	lcore.Puts("apply", NewNativeValue(2, func(env *Env) Value {
		x, y := env.SGet(0), env.SGet(1)
		if x.ty != Tclosure {
			x.panicType(Tclosure)
		}
		if y.ty != Tmap {
			y.panicType(Tmap)
		}
		newEnv := NewEnv(x.AsClosure().env)
		for _, v := range y.AsMap().l {
			newEnv.SPush(v)
		}
		return x.AsClosure().Exec(newEnv)
	}))
	lcore.Puts("storeinto", NewNativeValue(3, func(env *Env) Value {
		e, x, y := env.SGet(0), env.SGet(1), env.SGet(2)
		if x.ty != Tnumber {
			x.panicType(Tnumber)
		}
		if e.ty != Tgeneric {
			e.panicType(Tgeneric)
		}
		(*Env)(e.AsGeneric()).Set(uint32(x.AsNumber()), y)
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
		x := env.SGet(0)
		if x.ty != Tstring {
			x.panicType(Tstring)
		}
		cls, err := LoadString(x.AsString())
		if err != nil {
			return NewStringValue(err.Error())
		}
		return NewClosureValue(cls)
	}))
	lcore.Puts("remove", NewNativeValue(2, func(env *Env) Value {
		s := env.SGet(0)
		if s.ty != Tmap {
			s.panicType(Tmap)
		}
		return s.AsMap().Remove(env.Get(1))
	}))
	lcore.Puts("copy", NewNativeValue(5, func(env *Env) Value {
		dst, src := env.SGet(0).testType(Tmap).AsMap(), env.SGet(2).testType(Tmap).AsMap()
		dstPos, srcPos := int(env.SGet(1).testType(Tnumber).AsNumber()), int(env.SGet(3).testType(Tnumber).AsNumber())
		length := int(env.SGet(4).testType(Tnumber).AsNumber())
		return NewNumberValue(float64(copy(dst.l[dstPos:], src.l[srcPos:srcPos+length])))
	}))
	lcore.Puts("sub", NewNativeValue(2, func(env *Env) Value {
		src := env.SGet(0)
		start, end := int(env.SGet(1).testType(Tnumber).AsNumber()), -1
		if env.SSize() > 2 {
			end = int(env.SGet(2).testType(Tnumber).AsNumber())
		}
		switch src.ty {
		case Tmap:
			m, m2 := NewMap(), src.AsMap()
			if end == -1 {
				end = len(m2.l)
			}
			m.l = make([]Value, end-start)
			copy(m.l, m2.l[start:end])
			return NewMapValue(m)
		case Tstring:
			buf2 := src.AsString()
			if end == -1 {
				end = len(buf2)
			}
			buf := make([]byte, end-start)
			copy(buf, buf2[start:end])
			return NewStringValue(string(buf))
		default:
			log.Panicf("can't call sub on %v", src)
		}
		return Value{}
	}))
	lcore.Puts("char", NewNativeValue(1, func(env *Env) Value {
		return NewStringValue(char(env.SGet(0).AsNumber(), true))
	}))
	lcore.Puts("utf8char", NewNativeValue(1, func(env *Env) Value {
		return NewStringValue(char(env.SGet(0).AsNumber(), false))
	}))
	lcore.Puts("append", NewNativeValue(2, func(env *Env) Value {
		src, v := env.SGet(0), env.SGet(1)
		if src.ty != Tmap {
			src.panicType(Tmap)
		}
		m := src.AsMap()
		m.l = append(m.l, v)
		return src
	}))

	lcore.Puts("sync", NewMapValue(NewMap().
		Puts("run", NewNativeValue(1, func(env *Env) Value {
			if env.SGet(0).ty != Tclosure {
				env.SGet(0).panicType(Tclosure)
			}
			cls := env.SGet(0).AsClosure()
			newEnv := NewEnv(cls.env)
			if cls.ArgsCount() > env.SSize()-1 {
				panic("not enough arguments to start")
			}
			for i := 1; i < env.SSize(); i++ {
				newEnv.SPush(env.SGet(i))
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
			m.Puts("add", NewNativeValue(1, func(env *Env) Value { wg.Add(int(env.SGet(0).AsNumber())); return NewValue() }))
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
				env.SGet(0).AsClosure().lastenv = nil
				return env.SGet(0)
			})).
			Puts("set", NewNativeValue(3, func(env *Env) Value {
				cls := env.SGet(0).AsClosure()
				switch name := env.SGet(1).AsString(); name {
				case "argscount":
					cls.argsCount = byte(env.SGet(2).AsNumber())
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
					cls.source = env.SGet(2).AsString()
				}
				return NewClosureValue(cls)
			})).
			Puts("write", NewNativeValue(4, func(env *Env) Value {
				cls := env.SGet(0).AsClosure()
				cls.code = append(cls.code, makeop(
					byte(env.SGet(1).AsNumber()),
					uint32(env.SGet(2).AsNumber()),
					uint32(env.SGet(3).AsNumber()),
				))
				return Value{}
			})).
			Puts("writeconst", NewNativeValue(2, func(env *Env) Value {
				cls := env.SGet(0).AsClosure()
				cls.consts = append(cls.consts, env.SGet(1))
				return Value{}
			})))).
		Puts("_", Value{})))

	CoreLibs["std"] = NewMapValue(lcore)

	initIOLib()
	initMathLib()
}
