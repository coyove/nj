package lib

import (
	"math"

	"github.com/coyove/script"
)

func init() {
	script.AddGlobalValue("mod", func(env *script.Env) {
		env.A = script.Float(math.Mod(env.In(0, script.VNumber).Float(),
			env.In(1, script.VNumber).Float()))
	})
	script.AddGlobalValue("cos", func(env *script.Env) {
		env.A = script.Float(math.Cos(env.In(0, script.VNumber).Float()))
	})
	script.AddGlobalValue("sin", func(env *script.Env) {
		env.A = script.Float(math.Sin(env.In(0, script.VNumber).Float()))
	})
	script.AddGlobalValue("tan", func(env *script.Env) {
		env.A = script.Float(math.Tan(env.In(0, script.VNumber).Float()))
	})
	script.AddGlobalValue("acos", func(env *script.Env) {
		env.A = script.Float(math.Acos(env.In(0, script.VNumber).Float()))
	})
	script.AddGlobalValue("asin", func(env *script.Env) {
		env.A = script.Float(math.Asin(env.In(0, script.VNumber).Float()))
	})
	script.AddGlobalValue("atan", func(env *script.Env) {
		env.A = script.Float(math.Atan(env.In(0, script.VNumber).Float()))
	})
	script.AddGlobalValue("atan2", func(env *script.Env) {
		env.A = script.Float(math.Atan2(env.In(0, script.VNumber).Float(), env.In(1, script.VNumber).Float()))
	})
	script.AddGlobalValue("ldexp", func(env *script.Env) {
		env.A = script.Float(math.Ldexp(env.In(0, script.VNumber).Float(), int(env.InInt(1, 0))))
	})
	script.AddGlobalValue("modf", func(env *script.Env) {
		a, b := math.Modf(env.In(0, script.VNumber).Float())
		env.Return2(script.Float(a), script.Float(b))
	})
}
