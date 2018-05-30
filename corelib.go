package potatolang

import (
	"log"
	"sync"
)

var CoreLibNames = []string{
	"std", "io", "math",
}

var CoreLibs = map[string]Value{}

func char(v float64, ascii bool) string {
	if ascii {
		return string([]byte{byte(v)})
	}
	return string(rune(v))
}

func initCoreLibs() {
	lcore := NewMap()
	lcore.Puts("genlist", NewNativeValue(1, func(env *Env) Value {
		v := env.SGet(0)
		if v.ty != Tnumber {
			v.panicType(Tnumber)
		}
		return NewMapValue(NewMapSize(int(v.AsNumber())))
	}))
	lcore.Puts("noenvescape", NewNativeValue(1, func(env *Env) Value {
		if env.SGet(0).ty != Tclosure {
			env.SGet(0).panicType(Tclosure)
		}
		env.SGet(0).AsClosure().noenvescape = true
		return env.SGet(0)
	}))
	lcore.Puts("stacktrace", NewNativeValue(0, func(env *Env) Value {
		e := ExecError{stacks: env.trace}
		return NewStringValue(e.Error())
	}))
	lcore.Puts("yreset", NewNativeValue(1, func(env *Env) Value {
		env.SGet(0).AsClosure().lastenv = nil
		return NewValue()
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

	CoreLibs["std"] = NewMapValue(lcore)

	initIOLib()
	initMathLib()
}
