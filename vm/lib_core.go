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
			for k, v := range env.R0.AsMapUnsafe() {
				newEnv.Stack().Clear()
				newEnv.Push(base.NewStringValue(k))
				newEnv.Push(v)
				Exec(newEnv, cls.Code())
			}
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

var lib_dup = LibFunc{
	name: "dup",
	args: 1,
	f: func(env *base.Env) base.Value {
		switch env.R0.Type() {
		case base.Tnil, base.Tnumber, base.Tstring, base.Tbool, base.Tclosure, base.Tgeneric:
			env.A = env.R0
		case base.Tlist:
			list0 := env.R0.AsListUnsafe()
			list1 := make([]base.Value, len(list0))
			copy(list1, list0)
			env.A = base.NewListValue(list1)
		case base.Tmap:
			map0 := env.R0.AsMapUnsafe()
			map1 := make(map[string]base.Value)
			for k, v := range map0 {
				map1[k] = v
			}
			env.A = base.NewMapValue(map1)
		case base.Tbytes:
			bytes0 := env.R0.AsBytesUnsafe()
			bytes1 := make([]byte, len(bytes0))
			copy(bytes1, bytes0)
			env.A = base.NewBytesValue(bytes1)
		}
		panic("shouldn't happen")
	},
}
