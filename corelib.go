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

func init() {
	lcore := NewMap()
	lcore.Put("genlist", NewNativeValue(1, func(env *Env) Value { return NewListValue(make([]Value, int(env.stack.Get(0).AsNumber()))) }))
	lcore.Put("genbytes", NewNativeValue(1, func(env *Env) Value { return NewBytesValue(make([]byte, int(env.stack.Get(0).AsNumber()))) }))
	lcore.Put("yreset", NewNativeValue(1, func(env *Env) Value { env.stack.Get(0).AsClosure().lastenv = nil; return NewValue() }))
	lcore.Put("tonumber", NewNativeValue(1, func(env *Env) Value {
		switch n := env.stack.Get(0); n.Type() {
		case Tnumber:
			return n
		case Tstring:
			num, err := parser.StringToNumber(n.AsStringUnsafe())
			if err != nil {
				return NewValue()
			}
			return NewNumberValue(num)
		case Tbool:
			if n.AsBoolUnsafe() {
				return NewNumberValue(1)
			}
			return NewNumberValue(0)
		default:
			return NewValue()
		}
	}))
	lcore.Put("remove", NewNativeValue(2, func(env *Env) Value {
		switch s := env.Get(0); s.ty {
		case Tmap:
			return s.AsMapUnsafe().Remove(env.Get(1).AsString())
		case Tlist:
			l := s.AsListUnsafe()
			if env.Size() == 2 {
				idx := int(env.Get(1).AsNumber())
				l = append(l[:idx], l[idx+1:]...)
			} else {
				idx, ln := int(env.Get(1).AsNumber()), int(env.Get(2).AsNumber())
				l = append(l[:idx], l[idx+ln:]...)
			}
			return NewListValue(l)
		case Tbytes:
			l := s.AsBytesUnsafe()
			if env.Size() == 2 {
				idx := int(env.Get(1).AsNumber())
				l = append(l[:idx], l[idx+1:]...)
			} else {
				idx, ln := int(env.Get(1).AsNumber()), int(env.Get(2).AsNumber())
				l = append(l[:idx], l[idx+ln:]...)
			}
			return NewBytesValue(l)
		default:
			log.Panicf("can't call remove on %+v", s)
			return NewValue()
		}
	}))

	lcore.Put("sync", NewMapValue(NewMap().
		Put("run", NewNativeValue(1, func(env *Env) Value {
			cls := env.Get(0).AsClosure()
			newEnv := NewEnv(cls.env)
			if cls.ArgsCount() > env.Size()-1 {
				panic("not enough arguments to start")
			}
			for i := 1; i < env.Size(); i++ {
				newEnv.Push(env.Get(uint32(i)))
			}
			go ExecCursor(newEnv, cls.code, cls.consts, 0)
			return NewValue()
		})).
		Put("mutex", NewNativeValue(0, func(env *Env) Value {
			m, mux := NewMap(), &sync.Mutex{}
			m.Put("lock", NewNativeValue(0, func(env *Env) Value { mux.Lock(); return NewValue() }))
			m.Put("unlock", NewNativeValue(0, func(env *Env) Value { mux.Unlock(); return NewValue() }))
			return NewMapValue(m)
		})).
		Put("waitgroup", NewNativeValue(0, func(env *Env) Value {
			m, wg := NewMap(), &sync.WaitGroup{}
			m.Put("add", NewNativeValue(1, func(env *Env) Value { wg.Add(int(env.stack.Get(0).AsNumber())); return NewValue() }))
			m.Put("done", NewNativeValue(0, func(env *Env) Value { wg.Done(); return NewValue() }))
			m.Put("wait", NewNativeValue(0, func(env *Env) Value { wg.Wait(); return NewValue() }))
			return NewMapValue(m)
		}))))

	var _bvalue, _bvalue2 = func(i uint64) Value { return NewNumberValue(math.Float64frombits(i)) }, NewBoolValue
	lcore.Put("u64", NewMapValue(NewMap().
		Put("inum", NewNativeValue(1, func(env *Env) Value { return NewNumberValue(float64(env.stack.Get(0).u64())) })).
		Put("iint", NewNativeValue(1, func(env *Env) Value { return _bvalue(uint64(env.stack.Get(0).AsNumber())) })).
		Put("iadd", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.stack.Get(0).u64() + env.stack.Get(1).u64()) })).
		Put("isub", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.stack.Get(0).u64() - env.stack.Get(1).u64()) })).
		Put("imul", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.stack.Get(0).u64() * env.stack.Get(1).u64()) })).
		Put("idiv", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.stack.Get(0).u64() / env.stack.Get(1).u64()) })).
		Put("imod", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.stack.Get(0).u64() % env.stack.Get(1).u64()) })).
		Put("inot", NewNativeValue(1, func(env *Env) Value { return _bvalue(^env.stack.Get(0).u64()) })).
		Put("iand", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.stack.Get(0).u64() & env.stack.Get(1).u64()) })).
		Put("ixor", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.stack.Get(0).u64() ^ env.stack.Get(1).u64()) })).
		Put("ilsh", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.stack.Get(0).u64() << env.stack.Get(1).u64()) })).
		Put("irsh", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.stack.Get(0).u64() >> env.stack.Get(1).u64()) })).
		Put("ior", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.stack.Get(0).u64() | env.stack.Get(1).u64()) })).
		Put("ilt", NewNativeValue(2, func(env *Env) Value { return _bvalue2(env.stack.Get(0).u64() < env.stack.Get(1).u64()) })).
		Put("ile", NewNativeValue(2, func(env *Env) Value { return _bvalue2(env.stack.Get(0).u64() <= env.stack.Get(1).u64()) })).
		Put("igt", NewNativeValue(2, func(env *Env) Value { return _bvalue2(env.stack.Get(0).u64() > env.stack.Get(1).u64()) })).
		Put("ige", NewNativeValue(2, func(env *Env) Value { return _bvalue2(env.stack.Get(0).u64() >= env.stack.Get(1).u64()) })).
		Put("ieq", NewNativeValue(2, func(env *Env) Value { return _bvalue2(env.stack.Get(0).u64() == env.stack.Get(1).u64()) })).
		Put("ine", NewNativeValue(2, func(env *Env) Value { return _bvalue2(env.stack.Get(0).u64() != env.stack.Get(1).u64()) }))))

	CoreLibs["std"] = NewMapValue(lcore)

	lio := NewMap()
	lio.Put("println", NewNativeValue(0, stdPrintln(os.Stdout)))
	lio.Put("print", NewNativeValue(0, stdPrint(os.Stdout)))
	lio.Put("write", NewNativeValue(0, stdWrite(os.Stdout)))
	lio.Put("errprintln", NewNativeValue(0, stdPrintln(os.Stderr)))
	lio.Put("errprint", NewNativeValue(0, stdPrint(os.Stderr)))
	lio.Put("errwrite", NewNativeValue(0, stdWrite(os.Stderr)))
	CoreLibs["io"] = NewMapValue(lio)

	lmath := NewMap()
	lmath.Put("sqrt", NewNativeValue(1, func(env *Env) Value { return NewNumberValue(math.Sqrt(env.stack.Get(0).AsNumber())) }))
	CoreLibs["math"] = NewMapValue(lmath)
}

func stdPrint(f *os.File) func(env *Env) Value {
	return func(env *Env) Value {
		s := env.Stack()
		for i := 0; i < s.Size(); i++ {
			f.WriteString(s.Get(i).ToPrintString())
		}

		return NewValue()
	}
}

func stdPrintln(f *os.File) func(env *Env) Value {
	return func(env *Env) Value {
		s := env.Stack()
		for i := 0; i < s.Size(); i++ {
			f.WriteString(s.Get(i).ToPrintString() + " ")
		}
		f.WriteString("\n")
		return NewValue()
	}
}

func stdWrite(f *os.File) func(env *Env) Value {
	return func(env *Env) Value {
		s := env.Stack()
		for i := 0; i < s.Size(); i++ {
			switch a := s.Get(i); a.ty {
			case Tbytes:
				f.Write(s.Get(i).AsBytesUnsafe())
			case Tstring:
				f.Write([]byte(s.Get(i).AsStringUnsafe()))
			default:
				log.Panicf("can't write to output: %+v", a)
			}
		}
		return NewValue()
	}
}
