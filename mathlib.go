package potatolang

import (
	"math"
)

func initMathLib() {

	lmath := NewMap()
	lmath.Puts("sqrt", NewNativeValue(1, func(env *Env) Value {
		return NewNumberValue(math.Sqrt(env.SGet(0).AsNumber()))
	}))

	CoreLibs["math"] = NewMapValue(lmath)
}
