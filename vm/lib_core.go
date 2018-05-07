package vm

import (
	"github.com/coyove/bracket/base"
)

var lib_foreach = LibFunc{
	name: "foreach",
	args: 2,
	f: func(env *base.Env) base.Value {
		cls := env.R1.AsClosure()
		newEnv := base.NewEnv(cls.Env())

		switch env.R0.Type() {
		case base.Tlist:
			for i, v := range env.R0.AsListUnsafe() {
				newEnv.Stack().Clear()
				newEnv.Push(base.NewNumberValue(float64(i)))
				newEnv.Push(v)
				Exec(newEnv, cls.Code())
			}
		case base.Tmap:
			// for k, v := range env.R0.AsMapUnsafe() {
			// 	newEnv.Stack().Clear()
			// 	newEnv.Push(base.NewStringValue(k))
			// 	newEnv.Push(v)
			// 	Exec(newEnv, cls.Code())
			// }
		case base.Tbytes:
			for i, v := range env.R0.AsBytesUnsafe() {
				newEnv.Stack().Clear()
				newEnv.Push(base.NewNumberValue(float64(i)))
				newEnv.Push(base.NewNumberValue(float64(v)))
				Exec(newEnv, cls.Code())
			}
		}
		return base.NewValue()
	},
}

var lib_typeof = LibFunc{
	name: "typeof",
	args: 1,
	f: func(env *base.Env) base.Value {
		switch env.R0.Type() {
		case base.Tnil:
			env.A = base.NewStringValue("nil")
		case base.Tnumber:
			env.A = base.NewStringValue("number")
		case base.Tstring:
			env.A = base.NewStringValue("string")
		case base.Tbool:
			env.A = base.NewStringValue("bool")
		case base.Tclosure:
			env.A = base.NewStringValue("closure")
		case base.Tgeneric:
			env.A = base.NewStringValue("generic")
		case base.Tlist:
			env.A = base.NewStringValue("list")
		case base.Tmap:
			env.A = base.NewStringValue("map")
		case base.Tbytes:
			env.A = base.NewStringValue("bytes")
		}
		panic("shouldn't happen")
	},
}
