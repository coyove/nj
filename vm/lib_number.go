package vm

import (
	"math"
	"strconv"
	"unsafe"

	"github.com/coyove/bracket/base"
)

var lib_numlongbits = LibFunc{
	name: "bits_long",
	args: 1,
	f: func(env *base.Env) base.Value {
		num := env.R0.AsNumber()
		return base.NewNumberValue(float64(*(*int64)(unsafe.Pointer(&num))))
	},
}

var lib_mathsqrt = LibFunc{
	name: "math_sqrt",
	args: 1,
	f: func(env *base.Env) base.Value {
		num := env.R0.AsNumber()
		return base.NewNumberValue(math.Sqrt(num))
	},
}

var lib_numtostring = LibFunc{
	name: "to_string",
	args: 1,
	f: func(env *base.Env) base.Value {
		switch env.R0.Type() {
		case base.Tnumber:
			num := env.R0.AsNumberUnsafe()
			if float64(int64(num)) == num {
				return base.NewStringValue(strconv.FormatInt(int64(num), 10))
			}
			return base.NewStringValue(strconv.FormatFloat(num, 'f', 9, 64))
		case base.Tstring:
			return env.R0
		default:
			return base.NewValue()
		}
	},
}

var lib_stringtonum = LibFunc{
	name: "to_number",
	args: 1,
	f: func(env *base.Env) base.Value {
		switch env.R0.Type() {
		case base.Tstring:
			str := env.R0.AsStringUnsafe()
			num, err := strconv.ParseFloat(str, 64)
			if err != nil {
				return base.NewValue()
			}
			return base.NewNumberValue(num)
		case base.Tnumber:
			return env.R0
		default:
			return base.NewValue()
		}
	},
}
