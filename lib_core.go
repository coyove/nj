package potatolang

import (
	"math"
	"math/rand"
	"runtime"
	"strconv"
)

var G = &Table{}

func initCoreLibs() {
	lclosure := &Table{}
	lclosure.Puts("copy", NativeFun(func(env *Env) {
		cls := env.In(0, FUN).Fun().Dup()
		env.A = Fun(cls)
	}))
	G.Puts("closure", Tab(lclosure))
	G.Puts("unpack", NativeFun(func(env *Env) {
		a := env.In(0, TAB).Tab().a
		start, end := 1, len(a)
		if len(env.Stack()) > 1 {
			start = int(env.Get(1).Expect(NUM).Num())
		}
		if len(env.Stack()) > 2 {
			end = int(env.Get(2).Expect(NUM).Num())
		}
		env.A = newUnpackedValue(a[start-1 : end])
	}))
	ltable := &Table{}
	ltable.Puts("insert", NativeFun(func(env *Env) {
		t := env.In(0, TAB).Tab()
		if len(env.Stack()) > 2 {
			t.Insert(env.In(1, NUM), env.Get(2))
			env.A = env.Get(2)
		} else {
			t.Insert(Num(float64(t.Len())), env.Get(1))
		}
	}))
	ltable.Puts("remove", NativeFun(func(env *Env) {
		t := env.In(0, TAB).Tab()
		n := t.Len()
		if len(env.Stack()) > 1 {
			n = int(env.Get(1).Expect(NUM).Num())
		}
		t.Remove(n)
	}))
	ltable.Puts("unpack", G.Get(Str("unpack")))
	G.Puts("table", Tab(ltable))
	G.Puts("type", NativeFun(func(env *Env) {
		env.A = Str(typeMappings[env.Get(0).Type()])
	}))
	G.Puts("rawset", NativeFun(func(env *Env) {
		env.In(0, TAB).Tab().RawPut(env.Get(1), env.Get(2))
	}))
	G.Puts("rawget", NativeFun(func(env *Env) {
		env.A = env.In(0, TAB).Tab().RawGet(env.Get(1))
	}))
	G.Puts("rawequal", NativeFun(func(env *Env) {
		switch v, r := env.Get(0), env.Get(1); v.Type() + r.Type() {
		case NumNum, BlnBln, NilNil:
			env.A = Bln(v == r)
		case StrStr:
			env.A = Bln(r.Str() == v.Str())
		case AnyAny:
			env.A = Bln(v.Any() == r.Any())
		case TabTab:
			env.A = Bln(v == r)
		case FunFun:
			env.A = Bln(v.Fun() == r.Fun())
		default:
			env.A = Bln(false)
		}
	}))
	G.Puts("next", NativeFun(func(env *Env) {
		k, v := env.In(0, TAB).Tab().Next(env.Get(1))
		env.Return(k, v)
	}))
	G.Puts("rawlen", NativeFun(func(env *Env) {
		switch env.A = env.Get(0); env.A.Type() {
		case TAB:
			env.A = Num(float64(env.A.Tab().Len()))
		case STR:
			env.A = Num(float64(len(env.A.Str())))
		default:
			env.A.ExpectMsg(TAB, "rawlen")
		}
	}))
	G.Puts("pcall", NativeFun(func(env *Env) {
		defer func() {
			if r := recover(); r != nil {
				env.Return(Bln(false))
			}
		}()
		a, v := env.In(0, FUN).Fun().Call(env.Stack()[1:]...)
		env.Return(Bln(true), append([]Value{a}, v...)...)
	}))
	G.Puts("select", NativeFun(func(env *Env) {
		switch a := env.Get(0); a.Type() {
		case STR:
			env.A = Num(float64(len(env.Stack()[1:])))
		case NUM:
			if u, idx := env.Stack()[1:], int(a.Num())-1; idx < len(u) {
				env.Return(u[idx], u[idx+1:]...)
			} else {
				env.Return(Value{})
			}
		}
	}))
	G.Puts("pack", NativeFun(func(env *Env) {
		t := &Table{a: env.In(0, UPK)._Upk()}
		env.A = Tab(t)
	}))
	G.Puts("setmetatable", NativeFun(func(env *Env) {
		if !env.Get(0).GetMetatable().RawGet(M__metatable).IsNil() {
			panicf("cannot change protected metatable")
		}
		if env.Get(1).IsNil() {
			env.Get(0).SetMetatable(nil)
		} else {
			env.Get(0).SetMetatable(env.In(1, TAB).Tab())
		}
		env.A = env.Get(0)
	}))
	G.Puts("getmetatable", NativeFun(func(env *Env) {
		t := env.Get(0).GetMetatable()
		if mt := t.RawGet(M__metatable); !mt.IsNil() {
			env.A = mt
		} else {
			env.A = Tab(t)
		}
	}))
	G.Puts("assert", NativeFun(func(env *Env) {
		if v := env.Get(0); !v.IsFalse() {
			env.A = v
			return
		}
		panic("assertion failed")
	}))
	G.Puts("pairs", NativeFun(func(env *Env) {
		env.Return(pairsNext, env.In(0, TAB))
	}))
	G.Puts("ipairs", NativeFun(func(env *Env) {
		if f := env.Get(0).GetMetamethod(M__ipairs); !f.IsNil() {
			env.A, env.V = f.ExpectMsg(FUN, "metamethod: ipairs").Fun().Call(env.Get(0))
			return
		}
		env.Return(ipairsNext, env.In(0, TAB), Num(0))
	}))
	G.Puts("tostring", NativeFun(func(env *Env) {
		v := env.Get(0)
		if f := v.GetMetamethod(M__tostring); f.Type() == FUN {
			env.A, _ = f.Fun().Call(v)
			return
		}
		env.A = Str(v.String())
	}))
	G.Puts("tonumber", NativeFun(func(env *Env) {
		v := env.Get(0)
		switch v.Type() {
		case NUM:
			env.A = v
		case STR:
			v, _ := strconv.ParseFloat(v.Str(), 64)
			env.A = Num(v)
		default:
			env.A = Value{}
		}
	}))
	G.Puts("collectgarbage", NativeFun(func(env *Env) {
		runtime.GC()
	}))
	//
	initLibAux()
	//	r := rand.New()
	lmath := &Table{}
	lmath.Puts("huge", Num(math.Inf(1)))
	lmath.Puts("pi", Num(math.Pi))
	lmath.Puts("e", Num(math.E))
	lmath.Puts("randomseed", NativeFun(func(env *Env) { rand.Seed(int64(env.In(0, NUM).Num())) }))
	lmath.Puts("random", NativeFun(func(env *Env) {
		switch len(env.Stack()) {
		case 2:
			a, b := int(env.In(0, NUM).Num()), int(env.In(1, NUM).Num())
			env.A = Num(float64(rand.Intn(b-a)+a) + 1)
		case 1:
			env.A = Num(float64(rand.Intn(int(env.In(0, NUM).Num()))) + 1)
		default:
			env.A = Num(rand.Float64())
		}
	}))
	lmath.Puts("sqrt", NativeFun(func(env *Env) { env.A = Num(math.Sqrt(env.In(0, NUM).Num())) }))
	lmath.Puts("floor", NativeFun(func(env *Env) { env.A = Num(math.Floor(env.In(0, NUM).Num())) }))
	lmath.Puts("ceil", NativeFun(func(env *Env) { env.A = Num(math.Ceil(env.In(0, NUM).Num())) }))
	lmath.Puts("fmod", NativeFun(func(env *Env) { env.A = Num(math.Mod(env.In(0, NUM).Num(), env.In(1, NUM).Num())) }))
	lmath.Puts("abs", NativeFun(func(env *Env) { env.A = Num(math.Abs(env.In(0, NUM).Num())) }))
	lmath.Puts("acos", NativeFun(func(env *Env) { env.A = Num(math.Acos(env.In(0, NUM).Num())) }))
	lmath.Puts("asin", NativeFun(func(env *Env) { env.A = Num(math.Asin(env.In(0, NUM).Num())) }))
	lmath.Puts("atan", NativeFun(func(env *Env) { env.A = Num(math.Atan(env.In(0, NUM).Num())) }))
	lmath.Puts("atan2", NativeFun(func(env *Env) { env.A = Num(math.Atan2(env.In(0, NUM).Num(), env.In(1, NUM).Num())) }))
	lmath.Puts("ldexp", NativeFun(func(env *Env) { env.A = Num(math.Ldexp(env.In(0, NUM).Num(), int(env.In(1, NUM).Num()))) }))
	lmath.Puts("modf", NativeFun(func(env *Env) { a, b := math.Modf(env.In(0, NUM).Num()); env.Return(Num(a), Num(float64(b))) }))
	lmath.Puts("min", NativeFun(func(env *Env) {
		if len(env.Stack()) == 0 {
			env.A = Value{}
		} else {
			min := env.Get(0).Expect(NUM).Num()
			for i := 1; i < len(env.Stack()); i++ {
				if x := env.Get(i).Expect(NUM).Num(); x < min {
					min = x
				}
			}
			env.A = Num(min)
		}
	}))
	lmath.Puts("max", NativeFun(func(env *Env) {
		if len(env.Stack()) == 0 {
			env.A = Value{}
		} else {
			max := env.Get(0).Expect(NUM).Num()
			for i := 1; i < len(env.Stack()); i++ {
				if x := env.Get(i).Expect(NUM).Num(); x > max {
					max = x
				}
			}
			env.A = Num(max)
		}
	}))
	G.Puts("math", Tab(lmath))

	lnative := &Table{}
	lnative.Puts("int64", NativeFun(func(env *Env) {
		if len(env.Stack()) == 2 {
			v := int64(uint32(env.In(0, NUM).Num()))<<32 | int64(uint32(env.In(1, NUM).Num()))
			env.A = Any(Int64(v))
		} else {
			env.A = Any(Int64(atoint64(env.Get(0))))
		}
	}))
	lnative.Puts("uint64", NativeFun(func(env *Env) {
		if len(env.Stack()) == 2 {
			v := uint64(uint32(env.In(0, NUM).Num()))<<32 | uint64(uint32(env.In(1, NUM).Num()))
			env.A = Any(UInt64(v))
		} else {
			env.A = Any(UInt64(atouint64(env.Get(0))))
		}
	}))
	lnative.Puts("bytes", NativeFun(func(env *Env) {
		switch v := env.Get(0); v.Type() {
		case NUM:
			env.A = Any(make(Bytes, int(v.Num())))
		case STR:
			env.A = Any(Bytes(v.Str()))
		case ANY:
			env.A = Any(Bytes(append([]byte{}, v.Any().(Bytes)...)))
		default:
			env.A = Value{}
		}
	}))
	G.Puts("native", Tab(lnative))
}

var (
	ipairsNext = NativeFun(func(env *Env) {
		idx := env.In(1, NUM).Num() + 1
		if v := env.In(0, TAB).Tab().Get(Num(idx)); v.IsNil() {
			env.Return(Value{})
		} else {
			env.Return(Num(idx), v)
		}
	})
	pairsNext = NativeFun(func(env *Env) {
		k := env.Get(1)
		if k, v := env.In(0, TAB).Tab().Next(k); v.IsNil() {
			env.Return(Value{})
		} else {
			env.Return(k, v)
		}
	})
)
