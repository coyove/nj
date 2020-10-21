package potatolang

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/coyove/potatolang/parser"
)

var g = map[string]Value{}

func init() {
	buildg("ref", func(env *Env) {
		v := env.Get(0)
		env.A = Any(&v)
	})
	buildg("deref", func(env *Env) {
		env.A = *env.In(0, ANY).Any().(*Value)
		if env.A.Type() == STK {
			env.V = env.A.unpackedStack().a
		}
	})
	buildg("array", func(env *Env) {
		env.V = make([]Value, env.In(0, NUM).Int())
		env.A = unpackedStack(&unpacked{a: env.V})
	})
	buildg("copyfunction", func(env *Env) {
		f := *env.In(0, FUN).Fun()
		env.A = Fun(&f)
	})
	buildg("type", func(env *Env) {
		env.A = Str(typeMappings[env.Get(0).Type()])
	})
	buildg("pcall", func(env *Env) {
		defer func() {
			if r := recover(); r != nil {
				env.Return(NumBool(false))
			}
		}()
		a, v := env.In(0, FUN).Fun().Call(env.Stack()[1:]...)
		env.Return(NumBool(true), append([]Value{a}, v...)...)
	})
	buildg("select", func(env *Env) {
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
	})
	buildg("assert", func(env *Env) {
		if v := env.Get(0); !v.IsFalse() {
			env.A = v
			return
		}
		panic("assertion failed")
	})

	buildg("tostring", func(env *Env) {
		v := env.Get(0)
		env.A = Str(v.String())
	})
	buildg("tonumber", func(env *Env) {
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
	})
	buildg("print", func(env *Env) {
		args := make([]interface{}, len(env.Stack()))
		for i := range args {
			args[i] = env.Stack()[i].Any()
		}
		if n, err := fmt.Println(args...); err != nil {
			env.Return(Value{}, Str(err.Error()))
		} else {
			env.Return(Num(float64(n)))
		}
	})
	buildg("infinite", Num(math.Inf(1)))
	buildg("PI", Num(math.Pi))
	buildg("E", Num(math.E))
	buildg("randomseed", func(env *Env) {
		rand.Seed(env.In(0, NUM).Int())
	})
	buildg("random", func(env *Env) {
		switch len(env.Stack()) {
		case 2:
			a, b := int(env.In(0, NUM).Int()), int(env.In(1, NUM).Int())
			env.A = Num(float64(rand.Intn(b-a)+a) + 1)
		case 1:
			env.A = Num(float64(rand.Intn(int(env.In(0, NUM).Int()))) + 1)
		default:
			env.A = Num(rand.Float64())
		}
	})
	buildg("sqrt", func(env *Env) { env.A = Num(math.Sqrt(env.In(0, NUM).F64())) })
	buildg("floor", func(env *Env) { env.A = Num(math.Floor(env.In(0, NUM).F64())) })
	buildg("ceil", func(env *Env) { env.A = Num(math.Ceil(env.In(0, NUM).F64())) })
	buildg("fmod", func(env *Env) { env.A = Num(math.Mod(env.In(0, NUM).F64(), env.In(1, NUM).F64())) })
	buildg("abs", func(env *Env) { env.A = Num(math.Abs(env.In(0, NUM).F64())) })
	buildg("acos", func(env *Env) { env.A = Num(math.Acos(env.In(0, NUM).F64())) })
	buildg("asin", func(env *Env) { env.A = Num(math.Asin(env.In(0, NUM).F64())) })
	buildg("atan", func(env *Env) { env.A = Num(math.Atan(env.In(0, NUM).F64())) })
	buildg("atan2", func(env *Env) { env.A = Num(math.Atan2(env.In(0, NUM).F64(), env.In(1, NUM).F64())) })
	buildg("ldexp", func(env *Env) { env.A = Num(math.Ldexp(env.In(0, NUM).F64(), int(env.In(1, NUM).F64()))) })
	buildg("modf", func(env *Env) { a, b := math.Modf(env.In(0, NUM).F64()); env.Return(Num(a), Num(float64(b))) })
	buildg("min", func(env *Env) {
		if len(env.Stack()) == 0 {
			env.A = Value{}
		} else {
			mathMinMax(env, false)
		}
	})
	buildg("max", func(env *Env) {
		if len(env.Stack()) == 0 {
			env.A = Value{}
		} else {
			mathMinMax(env, true)
		}
	})
	buildg("int", func(env *Env) {
		env.A = Int(env.In(0, NUM).Int())
	})
	buildg("time", func(env *Env) {
		env.A = Num(float64(time.Now().Unix()))
	})
	buildg("clock", func(env *Env) {
		x := time.Now()
		s := *(*[2]int64)(unsafe.Pointer(&x))
		env.A = Num(float64(s[1] / 1e9))
	})
	buildg("microclock", func(env *Env) {
		x := time.Now()
		s := *(*[2]int64)(unsafe.Pointer(&x))
		env.A = Num(float64(s[1] / 1e3))
	})
	buildg("exit", func(env *Env) {
		if v := env.Get(0); !v.IsNil() {
			os.Exit(int(env.In(0, NUM).Int()))
		}
		os.Exit(0)
	})
	buildg("strrep", func(env *Env) {
		env.A = Str(strings.Repeat(env.In(0, STR).Str(), int(env.In(1, NUM).Int())))
	})
	buildg("strchar", func(env *Env) {
		env.A = Str(string(rune(env.In(0, NUM).Int())))
	})
	buildg("match", func(env *Env) {
		m := regexp.MustCompile(env.In(0, STR).Str()).
			FindAllStringSubmatch(env.In(1, STR).Str(), int(env.InNum(2, Int(-1)).Int()))
		var mm []string
		for _, m := range m {
			for _, m := range m {
				mm = append(mm, m)
			}
		}
		if len(mm) > 0 {
			env.A = Str(mm[0])
			for i := 1; i < len(mm); i++ {
				env.V = append(env.V, Str(mm[i]))
			}
		}
	})
	buildg("mutex", func(env *Env) {
		env.A = Any(&sync.Mutex{})
	})
}

func mathMinMax(env *Env, max bool) {
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

func camelKey(k string) string {
	if k == "" {
		return k
	}
	if k[0] >= 'a' && k[0] <= 'z' {
		return string(k[0]-'a'+'A') + k[1:]
	}
	return k
}

func reflectLoad(v interface{}, k string) Value {
	rv := reflect.ValueOf(v)
	f := rv.MethodByName(k)
	if !f.IsValid() {
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		f := rv.FieldByName(k)
		if f.IsValid() {
			return Any(f.Interface())
		}
		panicf("%q not found in %#v", k, v)
	}
	return Fun(&Func{
		Name: k,
		native: func(env *Env) {
			rt := f.Type()
			rtNumIn := rt.NumIn()
			ins := make([]reflect.Value, 0, rtNumIn)
			getter := func(i int, t reflect.Type) reflect.Value {
				return reflect.ValueOf(env.Get(i).AnyTyped(t))
			}

			if !rt.IsVariadic() {
				for i := 0; i < rtNumIn; i++ {
					ins = append(ins, getter(i, rt.In(i)))
				}
			} else {
				for i := 0; i < rtNumIn-1; i++ {
					ins = append(ins, getter(i, rt.In(i)))
				}
				for i := rtNumIn - 1; i < env.Size(); i++ {
					ins = append(ins, getter(i, rt.In(rtNumIn-1)))
				}
			}

			outs := f.Call(ins)
			if rt.NumOut() == 0 {
				return
			}

			a := make([]Value, len(outs))
			for i := range outs {
				a[i] = Any(outs[i])
			}
			env.Return(a[0], a[1:]...)
		},
	})
}

func reflectStore(v interface{}, k string, v2 Value) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	f := rv.FieldByName(k)
	if !f.IsValid() || !f.CanAddr() {
		panicf("%q not assignable in %#v", k, v)
	}
	if f.Type() == reflect.TypeOf(Value{}) {
		f.Set(reflect.ValueOf(v2))
	} else {
		f.Set(reflect.ValueOf(v2.AnyTyped(f.Type())))
	}
}

func buildg(k string, v interface{}) {
	switch v := v.(type) {
	case func(*Env):
		g[k] = Fun(&Func{Name: k, native: v})
	case Value:
		g[k] = v
	}
}
