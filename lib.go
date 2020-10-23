package script

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
	"unsafe"

	"github.com/coyove/script/parser"
)

var (
	g   = map[string]Value{}
	now int64
)

func AddGlobalValue(k string, v interface{}) {
	switch v := v.(type) {
	case func(*Env):
		g[k] = Function(&Func{Name: k, native: v})
	default:
		g[k] = Interface(v)
	}
}

func init() {
	go func() {
		for a := range time.Tick(time.Second) {
			now = a.Unix()
		}
	}()

	AddGlobalValue("array", func(env *Env) {
		n := env.In(0, VNumber).Int()
		env.Return(Int(n), make([]Value, n)...)
	})
	AddGlobalValue("copyfunction", func(env *Env) {
		f := *env.In(0, VFunction).Function()
		env.A = Function(&f)
	})
	AddGlobalValue("type", func(env *Env) {
		env.A = _str(typeMappings[env.Get(0).Type()])
	})
	AddGlobalValue("pcall", func(env *Env) {
		defer func() {
			if r := recover(); r != nil {
				env.Return(Bool(false), Interface(errors.New(fmt.Sprint(r))))
			}
		}()
		a, v := env.In(0, VFunction).Function().CallEnv(env, env.Stack()[1:]...)
		env.Return(Bool(true), append([]Value{a}, v...)...)
	})
	AddGlobalValue("select", func(env *Env) {
		switch a := env.Get(0); a.Type() {
		case VString:
			env.A = Float(float64(len(env.Stack()[1:])))
		case VNumber:
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
	AddGlobalValue("num", func(env *Env) {
		v := env.Get(0)
		switch v.Type() {
		case VNumber:
			env.A = v
		case VString:
			switch v, _ := parser.StringToNumber(v._str()); v := v.(type) {
			case float64:
				env.A = Float(v)
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
	AddGlobalValue("infinite", Float(math.Inf(1)))
	AddGlobalValue("PI", Float(math.Pi))
	AddGlobalValue("E", Float(math.E))
	AddGlobalValue("randomseed", func(env *Env) { rand.Seed(env.In(0, VNumber).Int()) })
	AddGlobalValue("random", func(env *Env) {
		switch len(env.Stack()) {
		case 2:
			a, b := int(env.In(0, VNumber).Int()), int(env.In(1, VNumber).Int())
			env.A = Float(float64(rand.Intn(b-a)+a) + 1)
		case 1:
			env.A = Float(float64(rand.Intn(int(env.In(0, VNumber).Int()))) + 1)
		default:
			env.A = Float(rand.Float64())
		}
	})
	AddGlobalValue("sqrt", func(env *Env) { env.A = Float(math.Sqrt(env.In(0, VNumber).Float())) })
	AddGlobalValue("floor", func(env *Env) { env.A = Float(math.Floor(env.In(0, VNumber).Float())) })
	AddGlobalValue("ceil", func(env *Env) { env.A = Float(math.Ceil(env.In(0, VNumber).Float())) })
	AddGlobalValue("fmod", func(env *Env) { env.A = Float(math.Mod(env.In(0, VNumber).Float(), env.In(1, VNumber).Float())) })
	AddGlobalValue("abs", func(env *Env) { env.A = Float(math.Abs(env.In(0, VNumber).Float())) })
	AddGlobalValue("acos", func(env *Env) { env.A = Float(math.Acos(env.In(0, VNumber).Float())) })
	AddGlobalValue("asin", func(env *Env) { env.A = Float(math.Asin(env.In(0, VNumber).Float())) })
	AddGlobalValue("atan", func(env *Env) { env.A = Float(math.Atan(env.In(0, VNumber).Float())) })
	AddGlobalValue("atan2", func(env *Env) { env.A = Float(math.Atan2(env.In(0, VNumber).Float(), env.In(1, VNumber).Float())) })
	AddGlobalValue("ldexp", func(env *Env) { env.A = Float(math.Ldexp(env.In(0, VNumber).Float(), int(env.In(1, VNumber).Float()))) })
	AddGlobalValue("modf", func(env *Env) { a, b := math.Modf(env.In(0, VNumber).Float()); env.Return(Float(a), Float(float64(b))) })
	AddGlobalValue("min", func(env *Env) {
		if len(env.Stack()) > 0 {
			mathMinMax(env, false)
		}
	})
	AddGlobalValue("max", func(env *Env) {
		if len(env.Stack()) > 0 {
			mathMinMax(env, true)
		}
	})
	AddGlobalValue("str", func(env *Env) {
		env.A = env.NewString(env.Get(0).String())
	})
	AddGlobalValue("int", func(env *Env) {
		switch v := env.Get(0); v.Type() {
		case VNumber:
			env.A = Int(v.Int())
		default:
			v, err := strconv.ParseInt(v.String(), int(env.InInt(1, 10)), 64)
			env.Return(Int(v), Interface(err))
		}
	})
	AddGlobalValue("time", func(env *Env) {
		env.A = Float(float64(time.Now().Unix()))
	})
	AddGlobalValue("clock", func(env *Env) {
		x := time.Now()
		s := *(*[2]int64)(unsafe.Pointer(&x))
		switch env.InStr(0, "") {
		case "nano":
			env.A = Int(s[1])
		case "micro":
			env.A = Int(s[1] / 1e3)
		case "milli":
			env.A = Int(s[1] / 1e6)
		default:
			env.A = Int(s[1] / 1e9)
		}
	})
	AddGlobalValue("exit", func(env *Env) { os.Exit(int(env.InInt(0, 0))) })
	AddGlobalValue("char", func(env *Env) {
		env.A = env.NewString(string(rune(env.In(0, VNumber).Int())))
	})
	AddGlobalValue("rune", func(env *Env) {
		r, sz := utf8.DecodeRuneInString(env.In(0, VString).Str())
		env.Return(Int(int64(r)), Int(int64(sz)))
	})
	AddGlobalValue("match", func(env *Env) {
		rx, err := regexp.Compile(env.In(1, VString)._str())
		if err != nil {
			env.A = Interface(err)
			return
		}
		m := rx.FindAllStringSubmatch(env.Get(0).String(), int(env.InInt(2, -1)))
		var mm []string
		for _, m := range m {
			for _, m := range m {
				mm = append(mm, m)
			}
		}
		if len(mm) > 0 {
			env.A = _str(mm[0])
			for i := 1; i < len(mm); i++ {
				env.V = append(env.V, _str(mm[i]))
			}
		}
	})
	AddGlobalValue("startswith", func(env *Env) {
		env.A = Bool(strings.HasPrefix(env.In(0, VString)._str(), env.In(1, VString)._str()))
	})
	AddGlobalValue("endswith", func(env *Env) {
		env.A = Bool(strings.HasSuffix(env.In(0, VString)._str(), env.In(1, VString)._str()))
	})
	AddGlobalValue("trim", func(env *Env) {
		switch a, cutset := env.In(0, VString)._str(), env.InStr(1, " "); env.InStr(2, "") {
		case "left", "l":
			env.A = _str(strings.TrimLeft(a, cutset))
		case "right", "r":
			env.A = _str(strings.TrimRight(a, cutset))
		case "prefix", "start":
			env.A = _str(strings.TrimPrefix(a, cutset))
		case "suffix", "end":
			env.A = _str(strings.TrimSuffix(a, cutset))
		default:
			env.A = _str(strings.Trim(a, cutset))
		}
	})
	AddGlobalValue("replace", func(env *Env) {
		a := env.In(0, VString)._str()
		rx, err := regexp.Compile(env.In(1, VString)._str())
		if err != nil {
			env.A = Interface(err)
			return
		}
		switch f := env.Get(2); f.Type() {
		case VString:
			env.A = env.NewString(rx.ReplaceAllString(a, f._str()))
		case VFunction:
			env.A = env.NewString(rx.ReplaceAllStringFunc(a, func(in string) string {
				v, _ := f.Function().CallEnv(env, env.NewString(in))
				return v.String()
			}))
		}
	})
	AddGlobalValue("split", func(env *Env) {
		x := strings.Split(env.In(0, VString)._str(), env.In(1, VString)._str())
		v := make([]Value, len(x))
		for i := range x {
			v[i] = _str(x[i])
		}
		env.Return(v[0], v[1:]...)
	})
	AddGlobalValue("substr", func(env *Env) {
		s, a := env.In(0, VString)._str(), env.InInt(1, 1)
		b := env.InInt(2, int64(len(s)))
		env.A = _str(s[a-1 : b-1+1])
	})
	AddGlobalValue("strpos", func(env *Env) {
		a, b := env.In(0, VString)._str(), env.In(1, VString)._str()
		if env.InStr(1, "") == "last" {
			env.A = Int(int64(strings.LastIndex(a, b)) + 1)
		} else {
			env.A = Int(int64(strings.Index(a, b)) + 1)
		}
	})
	AddGlobalValue("mutex", func(env *Env) { env.A = Interface(&sync.Mutex{}) })
	AddGlobalValue("error", func(env *Env) { env.A = Interface(errors.New(env.InStr(0, ""))) })
	AddGlobalValue("iserror", func(env *Env) { _, ok := env.In(0, VInterface).Interface().(error); env.A = Bool(ok) })
	AddGlobalValue("jsonparse", func(env *Env) {
		j := strings.TrimSpace(env.In(0, VString)._str())
		switch j {
		case "", "null":
			return
		case "true", "false":
			env.A = Bool(j == "true")
			return
		}
		switch j[0] {
		case '[':
			var a []interface{}
			err := json.Unmarshal([]byte(j), &a)
			env.Return(Interface(a), Interface(err))
		case '{':
			a := map[string]interface{}{}
			err := json.Unmarshal([]byte(j), &a)
			env.Return(Interface(a), Interface(err))
		default:
			env.Return(Value{}, Interface(fmt.Errorf("malformed json string: %q", j)))
		}
	})
	AddGlobalValue("json", func(env *Env) {
		v := env.Get(0)
		if env.Size() > 1 {
			v = _unpackedStack(&unpacked{a: env.Stack()})
		}
		i := v.Interface()
		if err := reflectCheckCyclicStruct(i); err != nil {
			env.Return(Value{}, Interface(err))
			return
		}
		var buf []byte
		var err error
		if ident := env.InStr(1, ""); ident != "" {
			buf, err = json.MarshalIndent(i, "", ident)
		} else {
			buf, err = json.Marshal(i)
		}
		env.Return(env.NewStringBytes(buf), Interface(err))
	})
}

func mathMinMax(env *Env, max bool) {
	f, i, isInt := env.Get(0).Expect(VNumber).Num()
	if isInt {
		for ii := 1; ii < len(env.Stack()); ii++ {
			if x := env.Get(ii).Expect(VNumber).Int(); x >= i == max {
				i = x
			}
		}
		env.A = Int(i)
	} else {
		for i := 1; i < len(env.Stack()); i++ {
			if x, _, _ := env.Get(i).Expect(VNumber).Num(); x >= f == max {
				f = x
			}
		}
		env.A = Float(f)
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
