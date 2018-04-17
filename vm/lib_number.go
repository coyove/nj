package vm

import (
	"math"
	"strconv"
	"unsafe"

	"github.com/coyove/eugine/base"
)

var lib_numlongbits = LibFunc{
	name: "bits/long",
	args: 1,
	f: func(env *base.Env) base.Value {
		num := env.R0().Number()
		return base.NewNumberValue(float64(*(*int64)(unsafe.Pointer(&num))))
	},
}

var lib_mathsqrt = LibFunc{
	name: "math/sqrt",
	args: 1,
	f: func(env *base.Env) base.Value {
		num := env.R0().Number()
		return base.NewNumberValue(math.Sqrt(num))
	},
}

var lib_numtostring = LibFunc{
	name: "num/to-string",
	args: 1,
	f: func(env *base.Env) base.Value {
		num := env.R0().Number()
		if float64(int64(num)) == num {
			return base.NewStringValue(strconv.FormatInt(int64(num), 10))
		}
		if env.SizeR() == 1 {
			return base.NewStringValue(strconv.FormatFloat(num, 'f', 9, 64))
		}
		return base.NewStringValue(strconv.FormatFloat(num, 'f', int(env.R1().Number()), 64))
	},
}
