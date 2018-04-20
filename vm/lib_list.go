package vm

import (
	"github.com/coyove/bracket/base"
)

var lib_listmakelen = LibFunc{
	name: "list/make-len",
	args: 1,
	f: func(env *base.Env) base.Value {
		list := make([]base.Value, int(env.R0.Number()))
		return base.NewArrayValue(list)
	},
}
