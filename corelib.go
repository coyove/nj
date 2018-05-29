package potatolang

import (
	"log"
	"math"
	"os"
	"sync"

	"github.com/coyove/potatolang/parser"
)

var CoreLibNames = []string{
	"std", "io", "math",
}

var CoreLibs = map[string]Value{}

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
	lcore.Puts("genbytes", NewNativeValue(1, func(env *Env) Value {
		return NewBytesValue(make([]byte, int(env.SGet(0).AsNumber())))
	}))
	lcore.Puts("noenvescape", NewNativeValue(1, func(env *Env) Value {
		if env.SGet(0).ty != Tclosure {
			env.SGet(0).panicType(Tclosure)
		}
		env.SGet(0).AsClosure().noenvescape = true
		return env.SGet(0)
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
		switch s := env.Get(0); testTypes(s, env.SGet(1)) {
		case Tmap<<8 | Tstring:
			return s.AsMap().Remove(env.Get(1))
		case Tbytes<<8 | Tnumber:
			l := s.AsBytes()
			if env.SSize() == 2 {
				idx := int(env.Get(1).AsNumber())
				l = append(l[:idx], l[idx+1:]...)
			} else if env.SGet(2).ty == Tnumber {
				idx, ln := int(env.Get(1).AsNumber()), int(env.Get(2).AsNumber())
				l = append(l[:idx], l[idx+ln:]...)
			} else {
				log.Panicf("can't call remove on %+v with index %+v", s, env.SGet(2))
			}
			return NewBytesValue(l)
		default:
			log.Panicf("can't call remove on %+v", s)
			return NewValue()
		}
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

	lio := NewMap()
	lio.Puts("println", NewNativeValue(0, stdPrintln(os.Stdout)))
	lio.Puts("print", NewNativeValue(0, stdPrint(os.Stdout)))
	lio.Puts("write", NewNativeValue(0, stdWrite(os.Stdout)))
	lio.Puts("errprintln", NewNativeValue(0, stdPrintln(os.Stderr)))
	lio.Puts("errprint", NewNativeValue(0, stdPrint(os.Stderr)))
	lio.Puts("errwrite", NewNativeValue(0, stdWrite(os.Stderr)))
	CoreLibs["io"] = NewMapValue(lio)

	lmath := NewMap()
	lmath.Puts("sqrt", NewNativeValue(1, func(env *Env) Value { return NewNumberValue(math.Sqrt(env.SGet(0).AsNumber())) }))
	CoreLibs["math"] = NewMapValue(lmath)
}

func stdPrint(f *os.File) func(env *Env) Value {
	return func(env *Env) Value {
		for i := 0; i < env.SSize(); i++ {
			f.WriteString(env.SGet(i).ToPrintString())
		}

		return NewValue()
	}
}

func stdPrintln(f *os.File) func(env *Env) Value {
	return func(env *Env) Value {
		for i := 0; i < env.SSize(); i++ {
			f.WriteString(env.SGet(i).ToPrintString() + " ")
		}
		f.WriteString("\n")
		return NewValue()
	}
}

func stdWrite(f *os.File) func(env *Env) Value {
	return func(env *Env) Value {
		for i := 0; i < env.SSize(); i++ {
			switch a := env.SGet(i); a.ty {
			case Tbytes:
				f.Write(env.SGet(i).AsBytes())
			case Tstring:
				f.Write([]byte(env.SGet(i).AsString()))
			default:
				log.Panicf("can't write to output: %+v", a)
			}
		}
		return NewValue()
	}
}
