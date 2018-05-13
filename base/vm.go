package base

import (
	"log"
	"math"
	"os"
	"strconv"
	"sync"
)

var CoreLibNames = []string{
	"core", "io", "os", "math",
}

var CoreLibs = map[string]Value{}

func init() {
	lcore := new(Tree)
	lcore.Put("listn", NewNativeClosureValue(1, func(env *Env) Value { return NewListValue(make([]Value, int(env.stack.Get(0).AsNumber()))) }))
	lcore.Put("typeof", NewNativeClosureValue(1, func(env *Env) Value { return NewStringValue(TMapping[env.stack.Get(0).Type()]) }))
	lcore.Put("yreset", NewNativeClosureValue(1, func(env *Env) Value { env.stack.Get(0).AsClosure().lastenv = nil; return NewValue() }))
	lcore.Put("bytes", NewNativeClosureValue(1, func(env *Env) Value {
		if n := env.stack.Get(0); n.Type() == Tstring {
			return NewBytesValue([]byte(n.AsStringUnsafe()))
		} else if n.Type() == Tnumber {
			return NewBytesValue(make([]byte, int(n.AsNumberUnsafe())))
		} else {
			panic("can't generate the bytes")
		}
	}))
	lcore.Put("tostring", NewNativeClosureValue(1, func(env *Env) Value {
		switch n := env.stack.Get(0); n.Type() {
		case Tnumber:
			num := n.AsNumberUnsafe()
			if float64(int64(num)) == num {
				return NewStringValue(strconv.FormatInt(int64(num), 10))
			}
			return NewStringValue(strconv.FormatFloat(num, 'f', 9, 64))
		case Tstring:
			return n
		default:
			return NewStringValue(n.ToPrintString())
		}
	}))
	lcore.Put("tonumber", NewNativeClosureValue(1, func(env *Env) Value {
		switch n := env.stack.Get(0); n.Type() {
		case Tnumber:
			return n
		case Tstring:
			num, err := strconv.ParseFloat(n.AsStringUnsafe(), 64)
			if err != nil {
				return NewValue()
			}
			return NewNumberValue(num)
		default:
			return NewValue()
		}
	}))
	lcore.Put("go", NewNativeClosureValue(1, func(env *Env) Value {
		cls := env.Get(0).AsClosure()
		newEnv := NewEnv(cls.env)
		if cls.ArgsCount() > env.Size()-1 {
			panic("not enough arguments to start")
		}
		for i := 1; i < env.Size(); i++ {
			newEnv.Push(env.Get(int32(i)))
		}
		go Exec(newEnv, cls.Code())
		return NewValue()
	}))
	lcore.Put("foreach", NewNativeClosureValue(2, func(env *Env) Value {
		cls := env.stack.Get(1).AsClosure()
		newEnv := NewEnv(cls.Env())
		switch env.stack.Get(0).Type() {
		case Tlist:
			for i, v := range env.stack.Get(0).AsListUnsafe() {
				newEnv.Stack().Clear()
				newEnv.Push(NewNumberValue(float64(i)))
				newEnv.Push(v)
				Exec(newEnv, cls.Code())
			}
		case Tmap:
			for iter := env.stack.Get(0).AsMapUnsafe().Iterator(); iter.Next(); {
				newEnv.Stack().Clear()
				newEnv.Push(NewStringValue(iter.Key()))
				newEnv.Push(iter.Value())
				Exec(newEnv, cls.Code())
			}
		case Tbytes:
			for i, v := range env.stack.Get(0).AsBytesUnsafe() {
				newEnv.Stack().Clear()
				newEnv.Push(NewNumberValue(float64(i)))
				newEnv.Push(NewNumberValue(float64(v)))
				Exec(newEnv, cls.Code())
			}
		}
		return NewValue()
	}))
	lcore.Put("mutex", NewNativeClosureValue(0, func(env *Env) Value {
		m, mux := new(Tree), &sync.Mutex{}
		m.Put("lock", NewNativeClosureValue(0, func(env *Env) Value { mux.Lock(); return NewValue() }))
		m.Put("unlock", NewNativeClosureValue(0, func(env *Env) Value { mux.Unlock(); return NewValue() }))
		return NewMapValue(m)
	}))
	lcore.Put("waitgroup", NewNativeClosureValue(0, func(env *Env) Value {
		m, wg := new(Tree), &sync.WaitGroup{}
		m.Put("add", NewNativeClosureValue(1, func(env *Env) Value { wg.Add(int(env.stack.Get(0).AsNumber())); return NewValue() }))
		m.Put("done", NewNativeClosureValue(0, func(env *Env) Value { wg.Done(); return NewValue() }))
		m.Put("wait", NewNativeClosureValue(0, func(env *Env) Value { wg.Wait(); return NewValue() }))
		return NewMapValue(m)
	}))
	CoreLibs["core"] = NewMapValue(lcore)

	lio := new(Tree)
	lio.Put("println", NewNativeClosureValue(0, stdPrintln(os.Stdout)))
	lio.Put("print", NewNativeClosureValue(0, stdPrint(os.Stdout)))
	lio.Put("write", NewNativeClosureValue(0, stdWrite(os.Stdout)))
	lio.Put("errprintln", NewNativeClosureValue(0, stdPrintln(os.Stderr)))
	lio.Put("errprint", NewNativeClosureValue(0, stdPrint(os.Stderr)))
	lio.Put("errwrite", NewNativeClosureValue(0, stdWrite(os.Stderr)))
	CoreLibs["io"] = NewMapValue(lio)

	los := new(Tree)
	CoreLibs["os"] = NewMapValue(los)

	lmath := new(Tree)
	lmath.Put("sqrt", NewNativeClosureValue(1, func(env *Env) Value { return NewNumberValue(math.Sqrt(env.stack.Get(0).AsNumber())) }))
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
