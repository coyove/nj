package potatolang

import (
	"log"
	"math"
	"sync"

	"github.com/coyove/potatolang/parser"
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
		m := NewMap()
		m.l = make([]Value, int(v.AsNumber()))
		return NewMapValue(m)
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
	lcore.Puts("tonumber", NewNativeValue(1, func(env *Env) Value {
		switch n := env.SGet(0); n.Type() {
		case Tnumber:
			return n
		case Tstring:
			num, err := parser.StringToNumber(n.AsString())
			if err != nil {
				return NewValue()
			}
			return NewNumberValue(num)
		case Tbool:
			if n.AsBool() {
				return NewNumberValue(1)
			}
			return NewNumberValue(0)
		default:
			return NewValue()
		}
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
		switch src.ty {
		case Tmap:
			m := src.AsMap()
			m.l = append(m.l, v)
			return src
		case Tstring:
			if v.ty == Tstring {
				return NewStringValue(src.AsString() + v.AsString())
			}
			fallthrough
		default:
			log.Panicf("can't call append on %v", src)
		}
		return Value{}
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
			go ExecCursor(newEnv, cls.code, cls.consts, 0)
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

	var _bvalue, _bvalue2 = func(i uint64) Value { return NewNumberValue(math.Float64frombits(i)) }, NewBoolValue
	lcore.Puts("u64", NewMapValue(NewMap().
		Puts("inum", NewNativeValue(1, func(env *Env) Value { return NewNumberValue(float64(env.SGet(0).u64())) })).
		Puts("iint", NewNativeValue(1, func(env *Env) Value { return _bvalue(uint64(env.SGet(0).AsNumber())) })).
		Puts("iadd", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.SGet(0).u64() + env.SGet(1).u64()) })).
		Puts("isub", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.SGet(0).u64() - env.SGet(1).u64()) })).
		Puts("imul", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.SGet(0).u64() * env.SGet(1).u64()) })).
		Puts("idiv", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.SGet(0).u64() / env.SGet(1).u64()) })).
		Puts("imod", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.SGet(0).u64() % env.SGet(1).u64()) })).
		Puts("inot", NewNativeValue(1, func(env *Env) Value { return _bvalue(^env.SGet(0).u64()) })).
		Puts("iand", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.SGet(0).u64() & env.SGet(1).u64()) })).
		Puts("ixor", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.SGet(0).u64() ^ env.SGet(1).u64()) })).
		Puts("ilsh", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.SGet(0).u64() << env.SGet(1).u64()) })).
		Puts("irsh", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.SGet(0).u64() >> env.SGet(1).u64()) })).
		Puts("ior", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.SGet(0).u64() | env.SGet(1).u64()) })).
		Puts("ilt", NewNativeValue(2, func(env *Env) Value { return _bvalue2(env.SGet(0).u64() < env.SGet(1).u64()) })).
		Puts("ile", NewNativeValue(2, func(env *Env) Value { return _bvalue2(env.SGet(0).u64() <= env.SGet(1).u64()) })).
		Puts("igt", NewNativeValue(2, func(env *Env) Value { return _bvalue2(env.SGet(0).u64() > env.SGet(1).u64()) })).
		Puts("ige", NewNativeValue(2, func(env *Env) Value { return _bvalue2(env.SGet(0).u64() >= env.SGet(1).u64()) })).
		Puts("ieq", NewNativeValue(2, func(env *Env) Value { return _bvalue2(env.SGet(0).u64() == env.SGet(1).u64()) })).
		Puts("ine", NewNativeValue(2, func(env *Env) Value { return _bvalue2(env.SGet(0).u64() != env.SGet(1).u64()) }))))

	CoreLibs["std"] = NewMapValue(lcore)

	initIOLib()
	initMathLib()
}
