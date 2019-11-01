package potatolang

import (
	"math"

	"github.com/coyove/common/rand"
)

func initMathLib() {
	r := rand.New()
	lmath := NewStruct()

	lmath.Put("sqrt", NewNativeValue(1, func(env *Env) Value {
		return NewNumberValue(math.Sqrt(env.LocalGet(0).MustNumber()))
	}))
	lmath.Put("rand", NewStructValue(NewStruct().
		Put("intn", NewNativeValue(1, func(env *Env) Value {
			return NewNumberValue(float64(r.Intn(int(env.LocalGet(0).MustNumber()))))
		})).
		Put("bytes", NewNativeValue(1, func(env *Env) Value {
			return NewStringValue(string(r.Fetch(int(env.LocalGet(0).MustNumber()))))
		}))))

	CoreLibs["math"] = NewStructValue(lmath)
}
