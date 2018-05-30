package potatolang

import (
	"math"
)

func initMathLib() {

	lmath := NewMap()

	var _bvalue, _bvalue2 = func(i uint64) Value { return NewNumberValue(math.Float64frombits(i)) }, NewBoolValue
	lmath.Puts("u64", NewMapValue(NewMap().
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

	lmath.Puts("sqrt", NewNativeValue(1, func(env *Env) Value {
		return NewNumberValue(math.Sqrt(env.SGet(0).AsNumber()))
	}))

	CoreLibs["math"] = NewMapValue(lmath)
}
