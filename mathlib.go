package potatolang

import (
	"math"

	"github.com/coyove/common/rand"
)

func initMathLib() {
	r := rand.New()
	lmath := NewMap()

	var _bvalue, _bvalue2 = func(i uint64) Value { return NewNumberValue(math.Float64frombits(i)) }, NewBoolValue
	lmath.Puts("u64", NewMapValue(NewMap().
		Puts("num", NewNativeValue(1, func(env *Env) Value { return NewNumberValue(float64(env.SGet(0).u64())) })).
		Puts("int", NewNativeValue(1, func(env *Env) Value { return _bvalue(uint64(env.SGet(0).Num())) })).
		Puts("add", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.SGet(0).u64() + env.SGet(1).u64()) })).
		Puts("sub", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.SGet(0).u64() - env.SGet(1).u64()) })).
		Puts("mul", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.SGet(0).u64() * env.SGet(1).u64()) })).
		Puts("div", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.SGet(0).u64() / env.SGet(1).u64()) })).
		Puts("mod", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.SGet(0).u64() % env.SGet(1).u64()) })).
		Puts("not", NewNativeValue(1, func(env *Env) Value { return _bvalue(^env.SGet(0).u64()) })).
		Puts("and", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.SGet(0).u64() & env.SGet(1).u64()) })).
		Puts("xor", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.SGet(0).u64() ^ env.SGet(1).u64()) })).
		Puts("lsh", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.SGet(0).u64() << env.SGet(1).u64()) })).
		Puts("rsh", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.SGet(0).u64() >> env.SGet(1).u64()) })).
		Puts("or", NewNativeValue(2, func(env *Env) Value { return _bvalue(env.SGet(0).u64() | env.SGet(1).u64()) })).
		Puts("lt", NewNativeValue(2, func(env *Env) Value { return _bvalue2(env.SGet(0).u64() < env.SGet(1).u64()) })).
		Puts("le", NewNativeValue(2, func(env *Env) Value { return _bvalue2(env.SGet(0).u64() <= env.SGet(1).u64()) })).
		Puts("gt", NewNativeValue(2, func(env *Env) Value { return _bvalue2(env.SGet(0).u64() > env.SGet(1).u64()) })).
		Puts("ge", NewNativeValue(2, func(env *Env) Value { return _bvalue2(env.SGet(0).u64() >= env.SGet(1).u64()) })).
		Puts("eq", NewNativeValue(2, func(env *Env) Value { return _bvalue2(env.SGet(0).u64() == env.SGet(1).u64()) })).
		Puts("ne", NewNativeValue(2, func(env *Env) Value { return _bvalue2(env.SGet(0).u64() != env.SGet(1).u64()) }))))

	var _bvalue32 = func(i uint32) Value { return NewNumberValue(float64(i)) }
	lmath.Puts("u32", NewMapValue(NewMap().
		Puts("add", NewNativeValue(2, func(env *Env) Value { return _bvalue32(uint32(env.SGet(0).Num()) + uint32(env.SGet(1).Num())) })).
		Puts("sub", NewNativeValue(2, func(env *Env) Value { return _bvalue32(uint32(env.SGet(0).Num()) - uint32(env.SGet(1).Num())) })).
		Puts("mul", NewNativeValue(2, func(env *Env) Value { return _bvalue32(uint32(env.SGet(0).Num()) * uint32(env.SGet(1).Num())) })).
		Puts("div", NewNativeValue(2, func(env *Env) Value { return _bvalue32(uint32(env.SGet(0).Num()) / uint32(env.SGet(1).Num())) })).
		Puts("mod", NewNativeValue(2, func(env *Env) Value { return _bvalue32(uint32(env.SGet(0).Num()) % uint32(env.SGet(1).Num())) })).
		Puts("lsh", NewNativeValue(2, func(env *Env) Value { return _bvalue32(uint32(env.SGet(0).Num()) << uint32(env.SGet(1).Num())) })).
		Puts("rsh", NewNativeValue(2, func(env *Env) Value { return _bvalue32(uint32(env.SGet(0).Num()) >> uint32(env.SGet(1).Num())) }))))

	lmath.Puts("sqrt", NewNativeValue(1, func(env *Env) Value {
		return NewNumberValue(math.Sqrt(env.SGet(0).Num()))
	}))
	lmath.Puts("rand", NewMapValue(NewMap().
		Puts("intn", NewNativeValue(1, func(env *Env) Value {
			return NewNumberValue(float64(r.Intn(int(env.SGet(0).Num()))))
		})).
		Puts("bytes", NewNativeValue(1, func(env *Env) Value {
			return NewStringValue(string(r.Fetch(int(env.SGet(0).Num()))))
		}))))

	CoreLibs["math"] = NewMapValue(lmath)
}
