package vm

import (
	"math"
	"strconv"

	"github.com/coyove/bracket/base"
)

// _helper functions
var _fbits, _bvalue, _bvalue2 = math.Float64frombits, base.NewNumberValue, base.NewBoolValue
var _bvalues = func(a, b float64) base.Value { return base.NewListValue([]base.Value{_bvalue(a), _bvalue(b)}) }
var _bvalues2 = func(a float64, b int) base.Value { return _bvalues(a, float64(b)) }

var lib_math = []LibFunc{
	LibFunc{name: "inum", args: 1, f: func(env *base.Env) base.Value { return _bvalue(float64(env.R0.AsUint64())) }},
	LibFunc{name: "iint", args: 1, f: func(env *base.Env) base.Value { return _bvalue(_fbits(uint64(env.R0.AsNumber()))) }},
	LibFunc{name: "iadd", args: 2, f: func(env *base.Env) base.Value { return _bvalue(_fbits(env.R0.AsUint64() + env.R1.AsUint64())) }},
	LibFunc{name: "isub", args: 2, f: func(env *base.Env) base.Value { return _bvalue(_fbits(env.R0.AsUint64() - env.R1.AsUint64())) }},
	LibFunc{name: "imul", args: 2, f: func(env *base.Env) base.Value { return _bvalue(_fbits(env.R0.AsUint64() * env.R1.AsUint64())) }},
	LibFunc{name: "idiv", args: 2, f: func(env *base.Env) base.Value { return _bvalue(_fbits(env.R0.AsUint64() / env.R1.AsUint64())) }},
	LibFunc{name: "imod", args: 2, f: func(env *base.Env) base.Value { return _bvalue(_fbits(env.R0.AsUint64() % env.R1.AsUint64())) }},
	LibFunc{name: "inot", args: 1, f: func(env *base.Env) base.Value { return _bvalue(_fbits(^env.R0.AsUint64())) }},
	LibFunc{name: "iand", args: 2, f: func(env *base.Env) base.Value { return _bvalue(_fbits(env.R0.AsUint64() & env.R1.AsUint64())) }},
	LibFunc{name: "ixor", args: 2, f: func(env *base.Env) base.Value { return _bvalue(_fbits(env.R0.AsUint64() ^ env.R1.AsUint64())) }},
	LibFunc{name: "ilsh", args: 2, f: func(env *base.Env) base.Value { return _bvalue(_fbits(env.R0.AsUint64() << env.R1.AsUint64())) }},
	LibFunc{name: "irsh", args: 2, f: func(env *base.Env) base.Value { return _bvalue(_fbits(env.R0.AsUint64() >> env.R1.AsUint64())) }},
	LibFunc{name: "ior", args: 2, f: func(env *base.Env) base.Value { return _bvalue(_fbits(env.R0.AsUint64() | env.R1.AsUint64())) }},
	LibFunc{name: "ilt", args: 2, f: func(env *base.Env) base.Value { return _bvalue2(env.R0.AsUint64() < env.R1.AsUint64()) }},
	LibFunc{name: "ile", args: 2, f: func(env *base.Env) base.Value { return _bvalue2(env.R0.AsUint64() <= env.R1.AsUint64()) }},
	LibFunc{name: "igt", args: 2, f: func(env *base.Env) base.Value { return _bvalue2(env.R0.AsUint64() > env.R1.AsUint64()) }},
	LibFunc{name: "ige", args: 2, f: func(env *base.Env) base.Value { return _bvalue2(env.R0.AsUint64() >= env.R1.AsUint64()) }},
	LibFunc{name: "ieq", args: 2, f: func(env *base.Env) base.Value { return _bvalue2(env.R0.AsUint64() == env.R1.AsUint64()) }},
	LibFunc{name: "ine", args: 2, f: func(env *base.Env) base.Value { return _bvalue2(env.R0.AsUint64() != env.R1.AsUint64()) }},
	LibFunc{name: "mabs", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Abs(env.R0.AsNumber())) }},
	LibFunc{name: "macos", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Acos(env.R0.AsNumber())) }},
	LibFunc{name: "macosh", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Acosh(env.R0.AsNumber())) }},
	LibFunc{name: "masin", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Asin(env.R0.AsNumber())) }},
	LibFunc{name: "masinh", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Asinh(env.R0.AsNumber())) }},
	LibFunc{name: "matan", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Atan(env.R0.AsNumber())) }},
	LibFunc{name: "matan2", args: 2, f: func(env *base.Env) base.Value { return _bvalue(math.Atan2(env.R0.AsNumber(), env.R1.AsNumber())) }},
	LibFunc{name: "matanh", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Atanh(env.R0.AsNumber())) }},
	LibFunc{name: "mcbrt", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Cbrt(env.R0.AsNumber())) }},
	LibFunc{name: "mceil", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Ceil(env.R0.AsNumber())) }},
	LibFunc{name: "mcopysign", args: 2, f: func(env *base.Env) base.Value { return _bvalue(math.Copysign(env.R0.AsNumber(), env.R1.AsNumber())) }},
	LibFunc{name: "mcos", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Cos(env.R0.AsNumber())) }},
	LibFunc{name: "mcosh", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Cosh(env.R0.AsNumber())) }},
	LibFunc{name: "mdim", args: 2, f: func(env *base.Env) base.Value { return _bvalue(math.Dim(env.R0.AsNumber(), env.R1.AsNumber())) }},
	LibFunc{name: "merf", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Erf(env.R0.AsNumber())) }},
	LibFunc{name: "merfc", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Erfc(env.R0.AsNumber())) }},
	LibFunc{name: "mexp", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Exp(env.R0.AsNumber())) }},
	LibFunc{name: "mexp2", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Exp2(env.R0.AsNumber())) }},
	LibFunc{name: "mexpm1", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Expm1(env.R0.AsNumber())) }},
	LibFunc{name: "mfloor", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Floor(env.R0.AsNumber())) }},
	LibFunc{name: "mgamma", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Gamma(env.R0.AsNumber())) }},
	LibFunc{name: "mhypot", args: 2, f: func(env *base.Env) base.Value { return _bvalue(math.Hypot(env.R0.AsNumber(), env.R1.AsNumber())) }},
	LibFunc{name: "minf", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Inf(int(env.R0.AsNumber()))) }},
	LibFunc{name: "misinf", args: 2, f: func(env *base.Env) base.Value { return _bvalue2(math.IsInf(env.R0.AsNumber(), int(env.R1.AsNumber()))) }},
	LibFunc{name: "misnan", args: 1, f: func(env *base.Env) base.Value { return _bvalue2(math.IsNaN(env.R0.AsNumber())) }},
	LibFunc{name: "mj0", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.J0(env.R0.AsNumber())) }},
	LibFunc{name: "mj1", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.J1(env.R0.AsNumber())) }},
	LibFunc{name: "mjn", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Jn(int(env.R0.AsNumber()), env.R1.AsNumber())) }},
	LibFunc{name: "mlog", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Log(env.R0.AsNumber())) }},
	LibFunc{name: "mlog10", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Log10(env.R0.AsNumber())) }},
	LibFunc{name: "mlog1p", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Log1p(env.R0.AsNumber())) }},
	LibFunc{name: "mlog2", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Log2(env.R0.AsNumber())) }},
	LibFunc{name: "mlogb", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Logb(env.R0.AsNumber())) }},
	LibFunc{name: "mmax", args: 2, f: func(env *base.Env) base.Value { return _bvalue(math.Max(env.R0.AsNumber(), env.R1.AsNumber())) }},
	LibFunc{name: "mmin", args: 2, f: func(env *base.Env) base.Value { return _bvalue(math.Min(env.R0.AsNumber(), env.R1.AsNumber())) }},
	LibFunc{name: "mmod", args: 2, f: func(env *base.Env) base.Value { return _bvalue(math.Mod(env.R0.AsNumber(), env.R1.AsNumber())) }},
	LibFunc{name: "mnan", args: 0, f: func(env *base.Env) base.Value { return _bvalue(math.NaN()) }},
	LibFunc{name: "mnextafter", args: 2, f: func(env *base.Env) base.Value { return _bvalue(math.Nextafter(env.R0.AsNumber(), env.R1.AsNumber())) }},
	LibFunc{name: "mpow", args: 2, f: func(env *base.Env) base.Value { return _bvalue(math.Pow(env.R0.AsNumber(), env.R1.AsNumber())) }},
	LibFunc{name: "mpow10", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Pow10(int(env.R0.AsNumber()))) }},
	LibFunc{name: "mremainder", args: 2, f: func(env *base.Env) base.Value { return _bvalue(math.Remainder(env.R0.AsNumber(), env.R1.AsNumber())) }},
	LibFunc{name: "msignbit", args: 1, f: func(env *base.Env) base.Value { return _bvalue2(math.Signbit(env.R0.AsNumber())) }},
	LibFunc{name: "msin", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Sin(env.R0.AsNumber())) }},
	LibFunc{name: "msinh", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Sinh(env.R0.AsNumber())) }},
	LibFunc{name: "msqrt", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Sqrt(env.R0.AsNumber())) }},
	LibFunc{name: "mtan", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Tan(env.R0.AsNumber())) }},
	LibFunc{name: "mtanh", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Tanh(env.R0.AsNumber())) }},
	LibFunc{name: "mtrunc", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Trunc(env.R0.AsNumber())) }},
	LibFunc{name: "my0", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Y0(env.R0.AsNumber())) }},
	LibFunc{name: "my1", args: 1, f: func(env *base.Env) base.Value { return _bvalue(math.Y1(env.R0.AsNumber())) }},
	LibFunc{name: "myn", args: 2, f: func(env *base.Env) base.Value { return _bvalue(math.Yn(int(env.R0.AsNumber()), env.R1.AsNumber())) }},
	LibFunc{name: "msincos", args: 1, f: func(env *base.Env) base.Value { return _bvalues(math.Sincos(env.R0.AsNumber())) }},
	LibFunc{name: "mmodf", args: 1, f: func(env *base.Env) base.Value { return _bvalues(math.Modf(env.R0.AsNumber())) }},
	LibFunc{name: "mldexp", args: 2, f: func(env *base.Env) base.Value { return _bvalue(math.Ldexp(env.R0.AsNumber(), int(env.R1.AsNumber()))) }},
	LibFunc{name: "mlgamma", args: 1, f: func(env *base.Env) base.Value { return _bvalues2(math.Lgamma(env.R0.AsNumber())) }},
	LibFunc{name: "mfrexp", args: 1, f: func(env *base.Env) base.Value { return _bvalues2(math.Frexp(env.R0.AsNumber())) }},
}

var lib_numtostring = LibFunc{
	name: "tostring",
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
	name: "tonumber",
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
