package vm

import (
	"github.com/coyove/eugine/base"
)

var lib_foreach = LibFunc{
	name: "for-each",
	args: 2,
	f: func(env *base.Env) base.Value {
		newEnv := base.NewEnv(nil)
		cls := env.R1().Closure()

		switch env.R0().Type() {
		case base.TY_array:
			for i, v := range env.R0().Array() {
				newEnv.Stack().Clear()
				newEnv.Push(base.NewNumberValue(float64(i)))
				newEnv.Push(v)
				Exec(newEnv, cls.Code())
			}
		case base.TY_map:
			for k, v := range env.R0().Map() {
				newEnv.Stack().Clear()
				newEnv.Push(base.NewStringValue(k))
				newEnv.Push(v)
				Exec(newEnv, cls.Code())
			}
		case base.TY_bytes:
			for i, v := range env.R0().Bytes() {
				newEnv.Stack().Clear()
				newEnv.Push(base.NewNumberValue(float64(i)))
				newEnv.Push(base.NewNumberValue(float64(v)))
				Exec(newEnv, cls.Code())
			}
		}
		return base.NewValue()
	},
}
