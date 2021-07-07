package lib

import (
	"math"

	"github.com/coyove/script"
)

func init() {
	script.AddGlobalValue("mod", func(env *script.Env) {
		env.A = script.Float(math.Mod(env.Get(0).MustNum("mod #", 1).Float(), env.Get(1).MustNum("mod #", 2).Float()))
	})
	script.AddGlobalValue("cos", func(env *script.Env) {
		env.A = script.Float(math.Cos(env.Get(0).MustNum("cos", 0).Float()))
	})
	script.AddGlobalValue("sin", func(env *script.Env) {
		env.A = script.Float(math.Sin(env.Get(0).MustNum("sin", 0).Float()))
	})
	script.AddGlobalValue("tan", func(env *script.Env) {
		env.A = script.Float(math.Tan(env.Get(0).MustNum("tan", 0).Float()))
	})
	script.AddGlobalValue("acos", func(env *script.Env) {
		env.A = script.Float(math.Acos(env.Get(0).MustNum("acos", 0).Float()))
	})
	script.AddGlobalValue("asin", func(env *script.Env) {
		env.A = script.Float(math.Asin(env.Get(0).MustNum("asin", 0).Float()))
	})
	script.AddGlobalValue("atan", func(env *script.Env) {
		env.A = script.Float(math.Atan(env.Get(0).MustNum("atan", 0).Float()))
	})
	script.AddGlobalValue("atan2", func(env *script.Env) {
		env.A = script.Float(math.Atan2(env.Get(0).MustNum("atan2 #", 1).Float(), env.Get(1).MustNum("atan #", 2).Float()))
	})
	script.AddGlobalValue("ldexp", func(env *script.Env) {
		env.A = script.Float(math.Ldexp(env.Get(0).MustNum("ldexp", 0).Float(), int(env.Get(1).IntDefault(0))))
	})
	script.AddGlobalValue("modf", func(env *script.Env) {
		a, b := math.Modf(env.Get(0).MustNum("modf", 0).Float())
		env.A = script.Array(script.Float(a), script.Float(b))
	})
}
