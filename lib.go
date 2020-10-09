package potatolang

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"
	"unsafe"

	"github.com/coyove/potatolang/parser"
)

var G *Table

func init() {
	G = buildtable(
		"copyfunction", func(env *Env) {
			f := *env.In(0, FUN).Fun()
			env.A = Fun(&f)
		},
		"table", buildtable(
			"newhash", func(env *Env) {
				env.A = Tab(&Table{m: *NewMap(int(env.In(0, NUM).Int()))})
			},
			"insert", func(env *Env) {
				t := env.In(0, TAB).Tab()
				if len(env.Stack()) > 2 {
					t.Insert(env.In(1, NUM), env.Get(2))
					env.A = env.Get(2)
				} else {
					t.Insert(Num(float64(t.Len())), env.Get(1))
				}
			},
			"remove", func(env *Env) {
				t := env.In(0, TAB).Tab()
				n := t.Len()
				if len(env.Stack()) > 1 {
					n = int(env.Get(1).Expect(NUM).Int())
				}
				t.Remove(n)
			}),
		"unpack", Native(func(env *Env) {
			a := env.In(0, TAB).Tab().a
			start, end := 1, len(a)
			if len(env.Stack()) > 1 {
				start = int(env.Get(1).Expect(NUM).Int())
			}
			if len(env.Stack()) > 2 {
				end = int(env.Get(2).Expect(NUM).Int())
			}
			env.A = Tab(&Table{a: a[start-1 : end], unpacked: true})
		}),
		"type", Native(func(env *Env) {
			env.A = Str(typeMappings[env.Get(0).Type()])
		}),
		"next", Native(func(env *Env) {
			k, v := env.In(0, TAB).Tab().Next(env.Get(1))
			env.Return(k, v)
		}),
		"pcall", Native(func(env *Env) {
			defer func() {
				if r := recover(); r != nil {
					env.Return(NumBool(false))
				}
			}()
			a, v := env.In(0, FUN).Fun().Call(env.Stack()[1:]...)
			env.Return(NumBool(true), append([]Value{a}, v...)...)
		}),
		"select", Native(func(env *Env) {
			switch a := env.Get(0); a.Type() {
			case STR:
				env.A = Num(float64(len(env.Stack()[1:])))
			case NUM:
				if u, idx := env.Stack()[1:], int(a.Int())-1; idx < len(u) {
					env.Return(u[idx], u[idx+1:]...)
				} else {
					env.Return(Value{})
				}
			}
		}),
		"assert", Native(func(env *Env) {
			if v := env.Get(0); !v.IsFalse() {
				env.A = v
				return
			}
			panic("assertion failed")
		}),
		"pairs", Native(func(env *Env) {
			env.Return(pairsNext, env.In(0, TAB))
		}),
		"ipairs", Native(func(env *Env) {
			env.Return(ipairsNext, env.In(0, TAB), Num(0))
		}),
		"tostring", Native(func(env *Env) {
			v := env.Get(0)
			env.A = Str(v.String())
		}),
		"tonumber", Native(func(env *Env) {
			v := env.Get(0)
			switch v.Type() {
			case NUM:
				env.A = v
			case STR:
				switch v, _ := parser.StringToNumber(v.Str()); v := v.(type) {
				case float64:
					env.A = Num(v)
				case int64:
					env.A = Int(v)
				}
			default:
				env.A = Value{}
			}
		}),
		"collectgarbage", Native(func(env *Env) {
			runtime.GC()
		}),
		"print", Native(func(env *Env) {
			args := make([]interface{}, len(env.Stack()))
			for i := range args {
				args[i] = env.Stack()[i].Any()
			}
			if n, err := fmt.Println(args...); err != nil {
				env.Return(Value{}, Str(err.Error()))
			} else {
				env.Return(Num(float64(n)))
			}
		}),
		"math", buildtable(
			"huge", Num(math.Inf(1)),
			"pi", Num(math.Pi),
			"e", Num(math.E),
			"randomseed", Native(func(env *Env) { rand.Seed(env.In(0, NUM).Int()) }),
			"random", Native(func(env *Env) {
				switch len(env.Stack()) {
				case 2:
					a, b := int(env.In(0, NUM).Int()), int(env.In(1, NUM).Int())
					env.A = Num(float64(rand.Intn(b-a)+a) + 1)
				case 1:
					env.A = Num(float64(rand.Intn(int(env.In(0, NUM).Int()))) + 1)
				default:
					env.A = Num(rand.Float64())
				}
			}),
			"sqrt", func(env *Env) { env.A = Num(math.Sqrt(env.In(0, NUM).F64())) },
			"floor", func(env *Env) { env.A = Num(math.Floor(env.In(0, NUM).F64())) },
			"ceil", func(env *Env) { env.A = Num(math.Ceil(env.In(0, NUM).F64())) },
			"fmod", func(env *Env) { env.A = Num(math.Mod(env.In(0, NUM).F64(), env.In(1, NUM).F64())) },
			"abs", func(env *Env) { env.A = Num(math.Abs(env.In(0, NUM).F64())) },
			"acos", func(env *Env) { env.A = Num(math.Acos(env.In(0, NUM).F64())) },
			"asin", func(env *Env) { env.A = Num(math.Asin(env.In(0, NUM).F64())) },
			"atan", func(env *Env) { env.A = Num(math.Atan(env.In(0, NUM).F64())) },
			"atan2", func(env *Env) { env.A = Num(math.Atan2(env.In(0, NUM).F64(), env.In(1, NUM).F64())) },
			"ldexp", func(env *Env) { env.A = Num(math.Ldexp(env.In(0, NUM).F64(), int(env.In(1, NUM).F64()))) },
			"modf", func(env *Env) { a, b := math.Modf(env.In(0, NUM).F64()); env.Return(Num(a), Num(float64(b))) },
			"min", Native(func(env *Env) {
				if len(env.Stack()) == 0 {
					env.A = Value{}
				} else {
					mathMinMax(env, false)
				}
			}),
			"max", Native(func(env *Env) {
				if len(env.Stack()) == 0 {
					env.A = Value{}
				} else {
					mathMinMax(env, true)
				}
			})),
		"os", buildtable(
			"time", Native(func(env *Env) {
				if v := env.Get(0); !v.IsNil() {
					nz := func(v Value) int {
						if v.Type() == NUM {
							return int(v.Int())
						}
						return 0
					}
					t := env.In(0, TAB).Tab()
					env.A = Num(float64(time.Date(
						nz(t.Get(Str("year"))),
						time.Month(nz(t.Get(Str("month")))),
						nz(t.Get(Str("day"))),
						nz(t.Get(Str("hour"))),
						nz(t.Get(Str("min"))),
						nz(t.Get(Str("sec"))), 0, time.UTC).Unix()))
				} else {
					env.A = Num(float64(time.Now().Unix()))
				}
			}),
			"clock", Native(func(env *Env) {
				x := time.Now()
				s := *(*[2]int64)(unsafe.Pointer(&x))
				env.A = Num(float64(s[1] / 1e9))
			}),
			"microclock", Native(func(env *Env) {
				x := time.Now()
				s := *(*[2]int64)(unsafe.Pointer(&x))
				env.A = Num(float64(s[1] / 1e3))
			}),
			"exit", Native(func(env *Env) {
				if v := env.Get(0); !v.IsNil() {
					os.Exit(int(env.In(0, NUM).Int()))
				}
				os.Exit(0)
			})),
		"string", buildtable(
			"rep", Native(func(env *Env) {
				env.A = Str(strings.Repeat(env.In(0, STR).Str(), int(env.In(1, NUM).Int())))
			}),
			"char", Native(func(env *Env) {
				env.A = Str(string(rune(env.In(0, NUM).Int())))
			})),
	)
}

var (
	mathMinMax = func(env *Env, max bool) {
		f, i, isInt := env.Get(0).Expect(NUM).Num()
		if isInt {
			for ii := 1; ii < len(env.Stack()); ii++ {
				if x := env.Get(ii).Expect(NUM).Int(); x >= i == max {
					i = x
				}
			}
			env.A = Int(i)
		} else {
			for i := 1; i < len(env.Stack()); i++ {
				if x, _, _ := env.Get(i).Expect(NUM).Num(); x >= f == max {
					f = x
				}
			}
			env.A = Num(f)
		}
	}
	ipairsNext = Native(func(env *Env) {
		idx := env.In(1, NUM).Int() + 1
		if v := env.In(0, TAB).Tab().Get(Int(idx)); v.IsNil() {
			env.Return(Value{})
		} else {
			env.Return(Int(idx), v)
		}
	})
	pairsNext = Native(func(env *Env) {
		k := env.Get(1)
		if k, v := env.In(0, TAB).Tab().Next(k); v.IsNil() {
			env.Return(Value{})
		} else {
			env.Return(k, v)
		}
	})
)
