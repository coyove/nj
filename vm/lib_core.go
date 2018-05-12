package vm

import (
	"github.com/coyove/bracket/base"
)

var lib_map = []LibFunc{
	LibFunc{name: "iter", args: 1, f: func(env *base.Env) base.Value { return base.NewGenericValue(env.R0.AsMap().Iterator()) }},
	LibFunc{name: "iterend", args: 1, f: func(env *base.Env) base.Value { return base.NewGenericValue(env.R0.AsMap().Iterator().End()) }},
	LibFunc{name: "next", args: 1, f: func(env *base.Env) base.Value { return _bvalue2(env.R0.AsGeneric().(*base.Iterator).Next()) }},
	LibFunc{name: "prev", args: 1, f: func(env *base.Env) base.Value { return _bvalue2(env.R0.AsGeneric().(*base.Iterator).Prev()) }},
	LibFunc{name: "key", args: 1, f: func(env *base.Env) base.Value { return base.NewStringValue(env.R0.AsGeneric().(*base.Iterator).Key()) }},
	LibFunc{name: "val", args: 1, f: func(env *base.Env) base.Value { return env.R0.AsGeneric().(*base.Iterator).Value() }},
}

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

var lib_attach = []LibFunc{
	LibFunc{name: "attach", args: 2, f: func(env *base.Env) base.Value { return env.R0.Attach(0, env.R1) }},
	LibFunc{name: "detach", args: 1, f: func(env *base.Env) base.Value { return env.R0.Detach(0) }},
	LibFunc{name: "attach1", args: 2, f: func(env *base.Env) base.Value { return env.R0.Attach(1, env.R1) }},
	LibFunc{name: "detach1", args: 1, f: func(env *base.Env) base.Value { return env.R0.Detach(1) }},
	LibFunc{name: "attach2", args: 2, f: func(env *base.Env) base.Value { return env.R0.Attach(2, env.R1) }},
	LibFunc{name: "detach2", args: 1, f: func(env *base.Env) base.Value { return env.R0.Detach(2) }},
	LibFunc{name: "attach3", args: 2, f: func(env *base.Env) base.Value { return env.R0.Attach(3, env.R1) }},
	LibFunc{name: "detach3", args: 1, f: func(env *base.Env) base.Value { return env.R0.Detach(3) }},
}
