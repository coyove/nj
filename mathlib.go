package potatolang

import (
	"math"

	"github.com/coyove/common/rand"
)

func initMathLib() {
	r := rand.New()
	lmath := NewMap()

	lmath.Puts("sqrt", NewNativeValue(1, func(env *Env) Value {
		return NewNumberValue(math.Sqrt(env.LocalGet(0).MustNumber()))
	}))
	lmath.Puts("rand", NewMapValue(NewMap().
		Puts("intn", NewNativeValue(1, func(env *Env) Value {
			return NewNumberValue(float64(r.Intn(int(env.LocalGet(0).MustNumber()))))
		})).
		Puts("bytes", NewNativeValue(1, func(env *Env) Value {
			return NewStringValue(string(r.Fetch(int(env.LocalGet(0).MustNumber()))))
		}))))

	CoreLibs["math"] = NewMapValue(lmath)
}
