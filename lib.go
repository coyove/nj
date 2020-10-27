package script

import (
	"errors"
	"fmt"
	"io"
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
	"github.com/tidwall/gjson"
)

var (
	g   = map[string]Value{}
	now int64
	rg  = struct {
		sync.Mutex
		*rand.Rand
	}{Rand: rand.New(rand.NewSource(1))}
)

func AddGlobalValue(k string, v interface{}) {
	switch v := v.(type) {
	case func(*Env):
		g[k] = Function(&Func{name: k, native: v})
	default:
		tmp := Interface(v)
		if tmp.Type() == VFunction {
			tmp.Function().name = k
		}
		g[k] = tmp
	}
}

func RemoveGlobalValue(k string) {
	delete(g, k)
}

func init() {
	go func() {
		for a := range time.Tick(time.Second) {
			now = a.Unix()
		}
	}()

	AddGlobalValue("True", _interface(true))
	AddGlobalValue("False", _interface(false))
	AddGlobalValue("narray", func(env *Env) {
		n := env.In(0, VNumber).Int()
		env.Return(Int(n), make([]Value, n)...)
	})
	AddGlobalValue("resume", func(env *Env) {
		f := env.In(0, VFunction).Function()
		cursor := env.In(1, VNumber).Int()
		stack := env.Stack()[2:]
		newEnv := *env
		newEnv.stackOffset = uint32(len(*newEnv.stack))
		*newEnv.stack = append(*newEnv.stack, stack...)
		newEnv.grow(int(f.stackSize))
		env.A, env.V = execCursorLoop(newEnv, f, uint32(cursor))
	})
	AddGlobalValue("type", func(env *Env) { env.A = _str(env.Get(0).Type().String()) })
	AddGlobalValue("pcall", func(env *Env) {
		a, v, err := env.In(0, VFunction).Function().Call(env, env.Stack()[1:]...)
		if err == nil {
			env.Return(Bool(true), append([]Value{a}, v...)...)
		} else {
			env.Return(Bool(false), Interface(err))
		}
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
	AddGlobalValue("panic", func(env *Env) { panic(env.InStr(0, "user panic")) })
	AddGlobalValue("assert", func(env *Env) {
		if v := env.Get(0); !v.IsFalse() {
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
			switch v := parser.NewNumberFromString(v._str()); v.Type {
			case parser.Float:
				env.A = Float(v.FloatValue())
			case parser.Int:
				env.A = Int(v.IntValue())
			}
		default:
			env.A = Value{}
		}
	})
	AddGlobalValue("print", func(env *Env) {
		for _, a := range env.Stack() {
			fmt.Fprint(env.Global.Stdout, a.String())
		}
		fmt.Fprintln(env.Global.Stdout)
	})
	AddGlobalValue("println", func(env *Env) {
		for _, a := range env.Stack() {
			fmt.Fprint(env.Global.Stdout, a.String(), " ")
		}
		fmt.Fprintln(env.Global.Stdout)
	})
	AddGlobalValue("scanln", func(env *Env) {
		var results []Value
		nIdx := 0
		if p := env.InStr(0, ""); p != "" {
			fmt.Fprint(env.Global.Stdout, p)
			nIdx++
		}
		var r io.Reader = env.Global.Stdin
		if env.Global.MaxStackSize > 0 {
			r = io.LimitReader(r, (env.Global.MaxStackSize-int64(len(*env.stack)))*16)
		}
		for n := env.InInt(nIdx, 1); n > 0; n-- {
			var s string
			if _, err := fmt.Fscan(r, &s); err != nil {
				break
			}
			results = append(results, env.NewString(s))
		}
		if len(results) > 0 {
			env.Return(results[0], results[1:]...)
		}
	})
	AddGlobalValue("INF", Float(math.Inf(1)))
	AddGlobalValue("PI", Float(math.Pi))
	AddGlobalValue("E", Float(math.E))
	AddGlobalValue("randomseed", func(env *Env) {
		rg.Lock()
		rg.Rand.Seed(env.InInt(0, 1))
		rg.Unlock()
	})
	AddGlobalValue("random", func(env *Env) {
		rg.Lock()
		switch len(env.Stack()) {
		case 2:
			a, b := int(env.In(0, VNumber).Int()), int(env.In(1, VNumber).Int())
			env.A = Float(float64(rg.Intn(b-a)+a) + 1)
		case 1:
			env.A = Float(float64(rg.Intn(int(env.In(0, VNumber).Int()))) + 1)
		default:
			env.A = Float(rg.Float64())
		}
		rg.Unlock()
	})
	AddGlobalValue("bitand", func(env *Env) { env.A = Int(env.In(0, VNumber).Int() & env.In(1, VNumber).Int()) })
	AddGlobalValue("bitor", func(env *Env) { env.A = Int(env.In(0, VNumber).Int() | env.In(1, VNumber).Int()) })
	AddGlobalValue("bitxor", func(env *Env) { env.A = Int(env.In(0, VNumber).Int() ^ env.In(1, VNumber).Int()) })
	AddGlobalValue("bitrsh", func(env *Env) { env.A = Int(env.In(0, VNumber).Int() >> env.In(1, VNumber).Int()) })
	AddGlobalValue("bitlsh", func(env *Env) { env.A = Int(env.In(0, VNumber).Int() << env.In(1, VNumber).Int()) })
	AddGlobalValue("sqrt", func(env *Env) { env.A = Float(math.Sqrt(env.In(0, VNumber).Float())) })
	AddGlobalValue("floor", func(env *Env) { env.A = Float(math.Floor(env.In(0, VNumber).Float())) })
	AddGlobalValue("ceil", func(env *Env) { env.A = Float(math.Ceil(env.In(0, VNumber).Float())) })
	AddGlobalValue("mod", func(env *Env) { env.A = Float(math.Mod(env.In(0, VNumber).Float(), env.In(1, VNumber).Float())) })
	AddGlobalValue("abs", func(env *Env) { env.A = Float(math.Abs(env.In(0, VNumber).Float())) })
	AddGlobalValue("acos", func(env *Env) { env.A = Float(math.Acos(env.In(0, VNumber).Float())) })
	AddGlobalValue("asin", func(env *Env) { env.A = Float(math.Asin(env.In(0, VNumber).Float())) })
	AddGlobalValue("atan", func(env *Env) { env.A = Float(math.Atan(env.In(0, VNumber).Float())) })
	AddGlobalValue("atan2", func(env *Env) { env.A = Float(math.Atan2(env.In(0, VNumber).Float(), env.In(1, VNumber).Float())) })
	AddGlobalValue("ldexp", func(env *Env) { env.A = Float(math.Ldexp(env.In(0, VNumber).Float(), int(env.InInt(1, 0)))) })
	AddGlobalValue("modf", func(env *Env) { a, b := math.Modf(env.In(0, VNumber).Float()); env.Return(Float(a), Float(b)) })
	AddGlobalValue("min", func(env *Env) { mathMinMax(env, false) })
	AddGlobalValue("max", func(env *Env) { mathMinMax(env, true) })
	AddGlobalValue("str", func(env *Env) {
		if v := env.Get(0); v.Type() == VNumber {
			env.A = env.NewString(fmt.Sprintf(env.InStr(1, "%v"), v.Interface()))
		} else {
			env.A = env.NewString(v.String())
		}
	})
	AddGlobalValue("int", func(env *Env) {
		switch v := env.Get(0); v.Type() {
		case VNumber:
			env.A = Int(v.Int())
		default:
			v, err := strconv.ParseInt(v.String(), 0, 64)
			env.Return(Int(v), Interface(err))
		}
	})
	AddGlobalValue("time", func(env *Env) { env.A = Float(float64(time.Now().Unix())) })
	AddGlobalValue("sleep", func(env *Env) { time.Sleep(time.Duration(env.In(0, VNumber).Int()) * time.Millisecond) })
	AddGlobalValue("Go_time", func(env *Env) {
		if env.Size() > 0 {
			loc := time.UTC
			if env.InStr(7, "") == "local" {
				loc = time.Local
			}
			env.A = Interface(time.Date(
				int(env.InInt(0, 1970)), time.Month(env.InInt(1, 1)), int(env.InInt(2, 1)),
				int(env.InInt(3, 0)), int(env.InInt(4, 0)), int(env.InInt(5, 0)),
				int(env.InInt(6, 0)), loc,
			))
		} else {
			env.A = Interface(time.Now())
		}
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
		r, sz := utf8.DecodeRuneInString(env.In(0, VString)._str())
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
				v, _, _ := f.Function().Call(env, env.NewString(in))
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
	AddGlobalValue("json", func(env *Env) {
		cv := func(r gjson.Result) Value {
			switch r.Type {
			case gjson.String:
				return _str(r.Str)
			case gjson.Number:
				return Float(r.Float())
			case gjson.True, gjson.False:
				return Bool(r.Bool())
			}
			return Value{}
		}
		j := strings.TrimSpace(env.In(0, VString)._str())
		result := gjson.Get(j, env.In(1, VString)._str())
		switch result.Type {
		case gjson.String, gjson.Number, gjson.True, gjson.False:
			env.A = cv(result)
		default:
			if result.IsArray() {
				a := result.Array()
				if len(a) > 0 {
					tmp := make([]Value, len(a))
					for i := range a {
						switch a[i].Type {
						case gjson.String, gjson.False, gjson.Number, gjson.True:
							tmp[i] = cv(a[i])
						default:
							tmp[i] = _str(a[i].Raw)
						}
					}
					env.Return(Int(int64(len(a))), tmp...)
				}
			} else if result.IsObject() {
				env.A = _str(result.Raw)
			}
		}
	})
}

func mathMinMax(env *Env, max bool) {
	if len(env.Stack()) <= 0 {
		return
	}
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
