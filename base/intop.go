package base

import (
	"bytes"
	"log"
)

func AddI(env *Env) {
	switch l := env.R0(); l.Type() {
	case TY_number:
		v := l.Number()
		for i := 1; i < env.SizeR(); i++ {
			v += env.R(i).Number()
		}
		env.SetANumber(v)
	case TY_string:
		buf := &bytes.Buffer{}
		for i := 0; i < env.SizeR(); i++ {
			buf.WriteString(env.R(i).String())
		}
		env.SetA(NewStringValue(buf.String()))
	case TY_array:
		v := make([]Value, 0, len(l.Array())*env.SizeR())
		for i := 0; i < env.SizeR(); i++ {
			v = append(v, env.R(i).Array()...)
		}
		env.SetA(NewArrayValue(v))
	default:
		log.Panicf("can't add %v", l)
	}
}

func SubI(env *Env) {
	v := env.R(0).Number()
	for i := 1; i < env.SizeR(); i++ {
		v -= env.R(i).Number()
	}
	env.SetANumber(v)
}

func MulI(env *Env) {
	v := env.R(0).Number()
	for i := 1; i < env.SizeR(); i++ {
		v *= env.R(i).Number()
	}
	env.SetANumber(v)
}

func DivI(env *Env) {
	v := env.R(0).Number()
	for i := 1; i < env.SizeR(); i++ {
		v /= env.R(i).Number()
	}
	env.SetANumber(v)
}

func ModI(env *Env) {
	v := int64(env.R(0).Number())
	for i := 1; i < env.SizeR(); i++ {
		v %= int64(env.R(i).Number())
	}
	env.SetA(NewNumberValue(float64(v)))
}

func NotI(env *Env) {
	env.SetA(NewBoolValue(!env.R(0).Bool()))
}

func AndI(env *Env) {
	for i := 0; i < env.SizeR(); i++ {
		if env.R(i).IsFalse() {
			env.SetA(NewBoolValue(false))
			return
		}
	}
	env.SetA(NewBoolValue(true))
}

func OrI(env *Env) {
	for i := 0; i < env.SizeR(); i++ {
		if !env.R(i).IsFalse() {
			env.SetA(NewBoolValue(true))
			return
		}
	}
	env.SetA(NewBoolValue(false))
}

func XorI(env *Env) {
	first := !env.R(0).IsFalse()
	for i := 1; i < env.SizeR(); i++ {
		if !env.R(i).IsFalse() != first {
			env.SetA(NewBoolValue(true))
			return
		}
	}
	env.SetA(NewBoolValue(false))
}

func BitNotI(env *Env) {
	l := env.R(0)
	env.SetA(NewNumberValue(float64(^int64(l.Number()))))
}

func BitAndI(env *Env) {
	d := int64(env.R(0).Number())
	for i := 1; i < env.SizeR(); i++ {
		dv := int64(env.R(i).Number())
		d &= dv
	}
	env.SetA(NewNumberValue(float64(d)))
}

func BitOrI(env *Env) {
	d := int64(env.R(0).Number())
	for i := 1; i < env.SizeR(); i++ {
		dv := int64(env.R(i).Number())
		d |= dv
	}
	env.SetA(NewNumberValue(float64(d)))
}

func BitXorI(env *Env) {
	d := int64(env.R(0).Number())
	for i := 1; i < env.SizeR(); i++ {
		dv := int64(env.R(i).Number())
		d ^= dv
	}
	env.SetA(NewNumberValue(float64(d)))
}

func BitLshI(env *Env) {
	d := uint64(env.R(0).Number()) << uint64(env.R(1).Number())
	env.SetA(NewNumberValue(float64(d)))
}

func BitRshI(env *Env) {
	d := uint64(env.R(0).Number()) >> uint64(env.R(1).Number())
	env.SetA(NewNumberValue(float64(d)))
}

func LogicCompare(env *Env, comp func(r Value) bool, undirsedResult bool) {
	for i := 1; i < env.SizeR(); i++ {
		if comp(env.R(i)) == undirsedResult {
			env.SetA(NewBoolValue(false))
			return
		}
	}
	env.SetA(NewBoolValue(true))
}
