package vm

import (
	"github.com/coyove/bracket/base"
)

var lib_foreach = LibFunc{
	name: "foreach",
	args: 2,
	f: func(env *base.Env) base.Value {
		cls := env.R1.Closure()
		newEnv := base.NewEnv(cls.Env())

		switch env.R0.Type() {
		case base.TY_array:
			for i, v := range env.R0.Array() {
				newEnv.Stack().Clear()
				newEnv.Push(base.NewNumberValue(float64(i)))
				newEnv.Push(v)
				Exec(newEnv, cls.Code())
			}
		case base.TY_map:
			for k, v := range env.R0.Map() {
				newEnv.Stack().Clear()
				newEnv.Push(base.NewStringValue(k))
				newEnv.Push(v)
				Exec(newEnv, cls.Code())
			}
		case base.TY_bytes:
			for i, v := range env.R0.Bytes() {
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
		case base.TY_nil:
			env.A = base.NewStringValue("nil")
		case base.TY_number:
			env.A = base.NewStringValue("number")
		case base.TY_string:
			env.A = base.NewStringValue("string")
		case base.TY_bool:
			env.A = base.NewStringValue("bool")
		case base.TY_closure:
			env.A = base.NewStringValue("closure")
		case base.TY_generic:
			env.A = base.NewStringValue("generic")
		case base.TY_array:
			env.A = base.NewStringValue("list")
		case base.TY_map:
			env.A = base.NewStringValue("map")
		case base.TY_bytes:
			env.A = base.NewStringValue("bytes")
		}
		return base.NewValue()
	},
}
