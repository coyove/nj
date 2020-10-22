package potatolang

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/coyove/potatolang/parser"
)

var (
	g   = map[string]Value{}
	now int64
)

func init() {
	go func() {
		for a := range time.Tick(time.Second) {
			now = a.Unix()
		}
	}()

	AddGlobalValue("array", func(env *Env) {
		n := env.In(0, NUM).Int()
		env.Return(Int(n), make([]Value, n)...)
	})
	AddGlobalValue("copyfunction", func(env *Env) {
		f := *env.In(0, FUN).Fun()
		env.A = Fun(&f)
	})
	AddGlobalValue("type", func(env *Env) {
		env.A = Str(typeMappings[env.Get(0).Type()])
	})
	AddGlobalValue("pcall", func(env *Env) {
		defer func() {
			if r := recover(); r != nil {
				env.Return(NumBool(false))
			}
		}()
		a, v := env.In(0, FUN).Fun().Call(env.Stack()[1:]...)
		env.Return(NumBool(true), append([]Value{a}, v...)...)
	})
	AddGlobalValue("select", func(env *Env) {
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
	AddGlobalValue("assert", func(env *Env) {
		if v := env.Get(0); !v.IsFalse() {
			env.A = v
			return
		}
		panic("assertion failed")
	})

	AddGlobalValue("tostring", func(env *Env) {
		v := env.Get(0)
		env.A = Str(v.String())
	})
	AddGlobalValue("tonumber", func(env *Env) {
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
	AddGlobalValue("print", func(env *Env) {
		for _, a := range env.Stack() {
			fmt.Print(a.String())
		}
		fmt.Println()
	})
	AddGlobalValue("println", func(env *Env) {
		for _, a := range env.Stack() {
			fmt.Print(a.String(), " ")
		}
		fmt.Println()
	})
	AddGlobalValue("infinite", Num(math.Inf(1)))
	AddGlobalValue("PI", Num(math.Pi))
	AddGlobalValue("E", Num(math.E))
	AddGlobalValue("randomseed", func(env *Env) {
		rand.Seed(env.In(0, NUM).Int())
	})
	AddGlobalValue("random", func(env *Env) {
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
	AddGlobalValue("sqrt", func(env *Env) { env.A = Num(math.Sqrt(env.In(0, NUM).F64())) })
	AddGlobalValue("floor", func(env *Env) { env.A = Num(math.Floor(env.In(0, NUM).F64())) })
	AddGlobalValue("ceil", func(env *Env) { env.A = Num(math.Ceil(env.In(0, NUM).F64())) })
	AddGlobalValue("fmod", func(env *Env) { env.A = Num(math.Mod(env.In(0, NUM).F64(), env.In(1, NUM).F64())) })
	AddGlobalValue("abs", func(env *Env) { env.A = Num(math.Abs(env.In(0, NUM).F64())) })
	AddGlobalValue("acos", func(env *Env) { env.A = Num(math.Acos(env.In(0, NUM).F64())) })
	AddGlobalValue("asin", func(env *Env) { env.A = Num(math.Asin(env.In(0, NUM).F64())) })
	AddGlobalValue("atan", func(env *Env) { env.A = Num(math.Atan(env.In(0, NUM).F64())) })
	AddGlobalValue("atan2", func(env *Env) { env.A = Num(math.Atan2(env.In(0, NUM).F64(), env.In(1, NUM).F64())) })
	AddGlobalValue("ldexp", func(env *Env) { env.A = Num(math.Ldexp(env.In(0, NUM).F64(), int(env.In(1, NUM).F64()))) })
	AddGlobalValue("modf", func(env *Env) { a, b := math.Modf(env.In(0, NUM).F64()); env.Return(Num(a), Num(float64(b))) })
	AddGlobalValue("min", func(env *Env) {
		if len(env.Stack()) == 0 {
			env.A = Value{}
		} else {
			mathMinMax(env, false)
		}
	})
	AddGlobalValue("max", func(env *Env) {
		if len(env.Stack()) == 0 {
			env.A = Value{}
		} else {
			mathMinMax(env, true)
		}
	})
	AddGlobalValue("int", func(env *Env) {
		env.A = Int(env.In(0, NUM).Int())
	})
	AddGlobalValue("time", func(env *Env) {
		env.A = Num(float64(time.Now().Unix()))
	})
	AddGlobalValue("clock", func(env *Env) {
		x := time.Now()
		s := *(*[2]int64)(unsafe.Pointer(&x))
		env.A = Num(float64(s[1] / 1e9))
	})
	AddGlobalValue("microclock", func(env *Env) {
		x := time.Now()
		s := *(*[2]int64)(unsafe.Pointer(&x))
		env.A = Num(float64(s[1] / 1e3))
	})
	AddGlobalValue("exit", func(env *Env) {
		if v := env.Get(0); !v.IsNil() {
			os.Exit(int(env.In(0, NUM).Int()))
		}
		os.Exit(0)
	})
	AddGlobalValue("char", func(env *Env) {
		env.A = Str(string(rune(env.In(0, NUM).Int())))
	})
	AddGlobalValue("match", func(env *Env) {
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
	AddGlobalValue("mutex", func(env *Env) { env.A = Any(&sync.Mutex{}) })
	AddGlobalValue("error", func(env *Env) { env.A = Any(errors.New(env.InStr(0, ""))) })
	AddGlobalValue("iserror", func(env *Env) { _, ok := env.In(0, ANY).Any().(error); env.A = NumBool(ok) })
	AddGlobalValue("jsonparse", func(env *Env) {
		j := strings.TrimSpace(env.In(0, STR).Str())
		if len(j) == 0 {
			return
		}
		switch j[0] {
		case 'n':
		case 't':
			env.A = NumBool(true)
		case 'f':
			env.A = NumBool(false)
		case '[':
			var a []interface{}
			json.Unmarshal([]byte(j), &a)
			env.A = Any(a)
		case '{':
			a := map[string]interface{}{}
			json.Unmarshal([]byte(j), &a)
			env.A = Any(a)
		default:
			panicf("malformed json string: %q", j)
		}
	})
	AddGlobalValue("json", func(env *Env) {
		var cv func(Value) interface{}
		cv = func(v Value) interface{} {
			if v.Type() == STK {
				x := v.unpackedStack().a
				tmp := make([]interface{}, len(x))
				for i := range x {
					tmp[i] = cv(x[i])
				}
				return tmp
			}
			return v.Any()
		}
		v := env.Get(0)
		if env.Size() > 1 {
			v = unpackedStack(&unpacked{a: env.Stack()})
		}

		i := cv(v)
		if err := reflectCheckCyclicStruct(i); err != nil {
			env.Return(Value{}, Any(err))
			return
		}
		var buf []byte
		var err error
		if ident := env.InStr(1, ""); ident != "" {
			buf, err = json.MarshalIndent(i, "", ident)
		} else {
			buf, err = json.Marshal(i)
		}
		env.Return(StrBytes(buf), Any(err))
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

func ipow(base, exp int64) int64 {
	var result int64 = 1
	for {
		if exp&1 == 1 {
			result *= base
		}
		exp >>= 1
		if exp == 0 {
			break
		}
		base *= base
	}
	return result
}

func AddGlobalValue(k string, v interface{}) {
	switch v := v.(type) {
	case func(*Env):
		g[k] = Fun(&Func{Name: k, native: v})
	default:
		g[k] = Any(v)
	}
}
