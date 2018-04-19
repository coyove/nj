package base

import (
	"log"
)

func AddI(env *Env) {
	switch l := env.R0; l.Type() {
	case TY_number:
		env.A = NewNumberValue(l.Number() + env.R1.Number())
	case TY_string:
		env.A = NewStringValue(l.String() + env.R1.String())
	case TY_array:
		v := l.Array()
		v = append(v, env.R1)
		env.A = NewArrayValue(v)
	default:
		log.Panicf("can't add %v", l)
	}
}
