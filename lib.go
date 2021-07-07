package script

import (
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
	"unsafe"

	"github.com/coyove/script/parser"
	"github.com/tidwall/gjson"
)

const Version int64 = 301

var (
	g   = map[string]Value{}
	now int64
	rg  = struct {
		sync.Mutex
		*rand.Rand
	}{Rand: rand.New(rand.NewSource(1))}
)

func AddGlobalValue(k string, v interface{}, doc ...string) {
	switch v := v.(type) {
	case func(*Env):
		g[k] = Native(k, v, doc...)
	case func(*Env, Value) Value:
		g[k] = Native1(k, v, doc...)
	case func(*Env, Value, Value) Value:
		g[k] = Native2(k, v, doc...)
	case func(*Env, Value, Value, Value) Value:
		g[k] = Native3(k, v, doc...)
	default:
		g[k] = Any(v)
	}
}

func RemoveGlobalValue(k string) {
	delete(g, k)
}

func init() {
	go func() {
		for a := range time.Tick(time.Second / 2) {
			now = a.UnixNano()
		}
	}()

	AddGlobalValue("VERSION", Int(Version))
	AddGlobalValue("undefined", Undef)
	AddGlobalValue("_", func(env *Env) { env.A = env.Get(0) }, "nop()")
	AddGlobalValue("globals", func(env *Env) {
		keys := make([]string, 0, len(g))
		for k := range g {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		globals := make([]Value, len(keys))
		for i, k := range keys {
			globals[i] = g[k]
		}
		env.A = Array(globals...)
	}, "globals() => { g1, g2, ... }", "\tlist all global values")
	AddGlobalValue("doc", func(env *Env, f Value) Value {
		return Str(f.MustFunc("doc", 0).DocString)
	}, "doc(function) => string", "\treturn function's documentation")
	AddGlobalValue("new", func(env *Env, v Value) Value {
		m := *v.MustMap("new", 0)
		m.hashItems = append([]hashItem{}, m.hashItems...)
		m.items = append([]Value{}, m.items...)
		a := &RHMap{Parent: &m}
		return a.Value()
	})
	AddGlobalValue("len", func(env *Env, v Value) Value {
		switch v.Type() {
		case STR:
			return Float(float64(len(v.Str())))
		case MAP:
			return Int(int64(v.Map().Len()))
		case FUNC:
			return Float(float64(v.Func().NumParams()))
		case NUM, NIL, BOOL:
			return panicf("can't measure length of %v", v.Type())
		default:
			return Int(int64(reflectLen(v.Any())))
		}
	})
	AddGlobalValue("apply", func(env *Env, f, a Value) Value {
		res, err := f.MustFunc("apply()", 0).Call(a)
		if err != nil {
			panic(err)
		}
		return res
	}, "apply(function, array) => call function using arguments in array")
	AddGlobalValue("go", func(env *Env, f, a Value) (res Value) {
		return a
	}, "$f(function, array) => call function in new goroutine using arguments in array")

	// Debug libraries
	AddGlobalValue("debug", Map(
		Str("dumpstk"), Native("dumpstk", func(env *Env) {
			start := env.Size()
			for _, s := range env.DebugStacktrace {
				start += int(s.cls.StackSize)
			}
			var r []Value
			stack := (*env.stack)[start:]
			for _, el := range stack[:cap(stack)] {
				if el == watermark {
					break
				}
				r = append(r, el)
			}
			env.A = Array(r...)
		}, "$f() => { v1, v2, v3, ... }"),
		Str("kwargs"), Native("kwargs", func(env *Env) {
			start := env.Size()
			for _, s := range env.DebugStacktrace {
				start += int(s.cls.StackSize)
			}
			stack := (*env.stack)[start:]
			stack = stack[:cap(stack)]
			for i, el := range stack {
				if el == watermark {
					env.A = stack[i+1]
					break
				}
			}
		}, "$f() => { key1: value, key2: value2, ... }"),
		Str("locals"), Native("locals", func(env *Env) {
			var r []Value
			start := env.StackOffset - uint32(env.DebugCaller.StackSize)
			for i, name := range env.DebugCaller.Locals {
				idx := start + uint32(i)
				r = append(r, Int(int64(idx)), Str(name), (*env.stack)[idx])
			}
			env.A = Array(r...)
		}, "$f() => { index1, name1, value1, i2, n2, v2, i3, n3, v3, ... }"),
		Str("globals"), Native("globals", func(env *Env) {
			var r []Value
			for i, name := range env.Global.Func.Locals {
				r = append(r, Int(int64(i)), Str(name), (*env.Global.Stack)[i])
			}
			env.A = Array(r...)
		}, "$f() => { index1, name1, value1, i2, n2, v2, i3, n3, v3, ... }"),
		Str("set"), Native2("set", func(env *Env, idx, value Value) Value {
			(*env.Global.Stack)[idx.MustNum("set", 1).Int()] = value
			return Nil
		}, "$f(idx, value)"),
		Str("trace"), Native1("trace", func(env *Env, skip Value) Value {
			stacks := append(env.DebugStacktrace, stacktrace{cls: env.DebugCaller, cursor: env.DebugCursor})
			lines := make([]Value, 0, len(stacks))
			for i := len(stacks) - 1 - int(skip.IntDefault(0)); i >= 0; i-- {
				r := stacks[i]
				src := uint32(0)
				for i := 0; i < len(r.cls.Code.Pos); {
					var opx uint32 = math.MaxUint32
					ii, op, line := r.cls.Code.Pos.read(i)
					if ii < len(r.cls.Code.Pos)-1 {
						_, opx, _ = r.cls.Code.Pos.read(ii)
					}
					if r.cursor >= op && r.cursor < opx {
						src = line
						break
					}
					if r.cursor < op && i == 0 {
						src = line
						break
					}
					i = ii
				}
				lines = append(lines, Str(r.cls.Name), Int(int64(src)), Int(int64(r.cursor-1)))
			}
			return Array(lines...)
		}, "$f(skip) => { func_name1, line1, cursor1, n2, l2, c2, ... }"),
	))
	AddGlobalValue("narray", func(env *Env, n Value) Value {
		return Array(make([]Value, n.MustNum("narray()", 0).Int())...)
	}, "narray(n) => { nil, ..., nil }", "\treturn an array size of n, filled with nil")
	AddGlobalValue("type", func(env *Env) {
		env.A = Str(env.Get(0).Type().String())
	}, "type(value) => string", "\treturn value's type")
	AddGlobalValue("pcall", func(env *Env, f, a Value) Value {
		a, err := f.MustFunc("pcall", 0).Call(a)
		if err == nil {
			return a
		}
		if err, ok := err.(*ExecError); ok {
			return Any(err.r)
		}
		return Any(err)
	}, "pcall(function, arg1, arg2, ...) => result_of_function",
		"\texecute the function, catch panic and return as error")
	AddGlobalValue("panic", func(env *Env) { panic(env.Get(0)) }, "panic(value)")
	AddGlobalValue("assert", func(env *Env) {
		v := env.Get(0)
		if env.Size() <= 1 && v.IsFalse() {
			panicf("assertion failed")
		}
		if env.Size() == 2 && !v.Equal(env.Get(1)) {
			panicf("assertion failed: %v and %v", v, env.Get(1))
		}
		if env.Size() == 3 && !v.Equal(env.Get(1)) {
			panicf("%s: %v and %v", env.Get(2).String(), v, env.Get(1))
		}
	}, "assert(value)", "\tpanic when value is falsy",
		"assert(value1, value2)", "\tpanic when two values are not equal",
		"assert(value1, value2, msg)", "\tpanic message when two values are not equal",
	)
	AddGlobalValue("float", func(env *Env) {
		v := env.Get(0)
		switch v.Type() {
		case NUM:
			env.A = v
		case STR:
			switch v := parser.NewNumberFromString(v.Str()); v.Type {
			case parser.Float:
				env.A = Float(v.FloatValue())
			case parser.Int:
				env.A = Int(v.IntValue())
			}
		default:
			env.A = Value{}
		}
	}, "$f(value) => number", "\tconvert string to number")
	AddGlobalValue("stdout", func(env *Env) { env.A = _interface(env.Global.Stdout) }, "$f() => fd", "\treturn stdout fd")
	AddGlobalValue("stdin", func(env *Env) { env.A = _interface(env.Global.Stdin) }, "$f() => fd", "\treturn stdin fd")
	AddGlobalValue("print", func(env *Env) {
		length := 0
		for _, a := range env.Stack() {
			s := a.String()
			length += len(s)
			fmt.Fprint(env.Global.Stdout, s)
		}
		fmt.Fprintln(env.Global.Stdout)
		env.A = Int(int64(length))
	}, "print(a, b, c, ...)", "\tprint values, no space between them")
	AddGlobalValue("write", func(env *Env) {
		w := env.Get(0).Any().(io.Writer)
		for _, a := range env.Stack()[1:] {
			fmt.Fprint(w, a.String())
		}
	}, "write(stdout, a, b, c, ...)", "\twrite raw values to stdout")
	AddGlobalValue("println", func(env *Env) {
		for _, a := range env.Stack() {
			fmt.Fprint(env.Global.Stdout, a.String(), " ")
		}
		fmt.Fprintln(env.Global.Stdout)
	}, "println(a, b, c, ...)", "\tprint values, insert space between each of them")
	AddGlobalValue("scanln", func(env *Env, prompt, n Value) Value {
		fmt.Fprint(env.Global.Stdout, prompt.StringDefault(""))
		var results []Value
		var r io.Reader = env.Global.Stdin
		for i := n.IntDefault(1); i > 0; i-- {
			var s string
			if _, err := fmt.Fscan(r, &s); err != nil {
				break
			}
			results = append(results, Str(s))
		}
		return Array(results...)
	},
		"$f(prompt='', n=1) => { s1, s2, ..., sn }", "\tprint prompt and read N user inputs",
	)
	AddGlobalValue("math", Map(
		Str("INF"), Float(math.Inf(1)),
		Str("NEG_INF"), Float(math.Inf(-1)),
		Str("PI"), Float(math.Pi),
		Str("E"), Float(math.E),
		Str("randomseed"), Native("randomseed", func(env *Env) {
			rg.Lock()
			defer rg.Unlock()
			rg.Rand.Seed(env.Get(0).IntDefault(1))
		}, "randomseed(int)"),
		Str("random"), Native("random", func(env *Env) {
			rg.Lock()
			defer rg.Unlock()
			switch len(env.Stack()) {
			case 2:
				af, ai, aIsInt := env.Get(0).MustNum("random #", 1).Num()
				bf, bi, bIsInt := env.Get(1).MustNum("random #", 2).Num()
				if aIsInt && bIsInt {
					env.A = Int(int64(rg.Intn(int(bi-ai+1))) + ai)
				} else {
					env.A = Float(rg.Float64()*(bf-af) + af)
				}
			case 1:
				env.A = Int(int64(rg.Intn(int(env.Get(0).MustNum("random", 0).Int()))))
			default:
				env.A = Float(rg.Float64())
			}
		},
			"$f() => [0, 1)",
			"$f(n) => [0, n)",
			"$f(a, b) => [a, b]"),
		Str("sqrt"), Native1("sqrt", func(env *Env, v Value) Value { return Float(math.Sqrt(v.MustNum("sqrt", 0).Float())) }),
		Str("floor"), Native1("floor", func(env *Env, v Value) Value { return Float(math.Floor(v.MustNum("floor", 0).Float())) }),
		Str("ceil"), Native1("ceil", func(env *Env, v Value) Value { return Float(math.Ceil(v.MustNum("ceil", 0).Float())) }),
		Str("min"), Native("min", func(env *Env) { mathMinMax(env, "min #", false) }, "max(a, b, ...) => largest_number"),
		Str("max"), Native("max", func(env *Env) { mathMinMax(env, "max #", true) }, "min(a, b, ...) => smallest_number"),
		Str("pow"), Native2("pow", func(env *Env, a, b Value) Value {
			af, ai, aIsInt := a.MustNum("pow", 1).Num()
			bf, bi, bIsInt := b.MustNum("pow", 2).Num()
			if aIsInt && bIsInt {
				return Int(ipow(ai, bi))
			}
			return Float(math.Pow(af, bf))
		}, "pow(a, b) => a to the power of b"),
		Str("abs"), Native("abs", func(env *Env) {
			switch f, i, isInt := env.Get(0).MustNum("abs", 0).Num(); {
			case isInt && i < 0:
				env.A = Int(-i)
			case isInt && i >= 0:
				env.A = Int(i)
			default:
				env.A = Float(math.Abs(f))
			}
		}),
	))
	AddGlobalValue("str", StringMethods)
	AddGlobalValue("int", func(env *Env) {
		switch v := env.Get(0); v.Type() {
		case NUM:
			env.A = Int(v.Int())
		default:
			v, _ := strconv.ParseInt(v.String(), 0, 64)
			env.A = Int(v)
		}
	}, "int(value) => integer", "\tconvert value to integer number (int64)")
	AddGlobalValue("time", func(env *Env, prefix Value) Value {
		switch prefix.StringDefault("") {
		case "nano":
			return Int(time.Now().UnixNano())
		case "micro":
			return Int(time.Now().UnixNano() / 1e3)
		case "milli":
			return Int(time.Now().UnixNano() / 1e6)
		}
		return Int(time.Now().Unix())
	}, "time(nil|'nano'|'micro'|'milli') => int", "\tunix timestamp in (nano|micro|milli)seconds")
	AddGlobalValue("sleep", func(env *Env, milli Value) {
		time.Sleep(time.Duration(milli.IntDefault(0)) * time.Millisecond)
	}, "sleep(milliseconds)")
	AddGlobalValue("Go_time", func(env *Env) {
		if env.Size() > 0 {
			loc := time.UTC
			if env.Get(7).StringDefault("") == "local" {
				loc = time.Local
			}
			env.A = Any(time.Date(
				int(env.Get(0).IntDefault(1970)), time.Month(env.Get(1).IntDefault(1)), int(env.Get(2).IntDefault(1)),
				int(env.Get(3).IntDefault(0)), int(env.Get(4).IntDefault(0)), int(env.Get(5).IntDefault(0)),
				int(env.Get(6).IntDefault(0)), loc,
			))
		} else {
			env.A = Any(time.Now())
		}
	},
		"Go_time() => time.Time",
		"\treturns time.Time struct of current time",
		"Go_time(year, month, day, h, m, s, nanoseconds, 'local'|'utc') => time.Time",
		"\treturns time.Time struct constructed by the given arguments",
	)
	AddGlobalValue("clock", func(env *Env, prefix Value) Value {
		x := time.Now()
		s := *(*[2]int64)(unsafe.Pointer(&x))
		switch prefix.StringDefault("") {
		case "nano":
			return Int(s[1])
		case "micro":
			return Int(s[1] / 1e3)
		case "milli":
			return Int(s[1] / 1e6)
		}
		return Int(s[1] / 1e9)
	}, "clock(nil|'nano'|'micro'|'milli') => int", "\t(nano|micro|milli)seconds since startup")
	AddGlobalValue("exit", func(env *Env, code Value) Value {
		os.Exit(int(code.MustNum("exit", 0).Int()))
		return Nil
	}, "exit(code)")
	AddGlobalValue("chr", func(env *Env) { env.A = Str(string(rune(env.Get(0).MustNum("chr()", 0).Int()))) }, "chr(unicode) => string")
	AddGlobalValue("ord", func(env *Env) {
		r, _ := utf8.DecodeRuneInString(env.Get(0).MustStr("ord()", 0))
		env.A = Int(int64(r))
	}, "$f(string) => unicode")
	AddGlobalValue("re", func(env *Env, r Value) Value {
		rx, err := regexp.Compile(r.MustStr("build regexp", 0))
		if err != nil {
			panic(err)
		}
		a := Map(
			Str("_rx"), Any(rx),
			Str("match"), Native2("match", func(e *Env, rx, text Value) Value {
				return Bool(rx.Map().GetString("_rx").Any().(*regexp.Regexp).MatchString(text.MustStr("match text", 0)))
			}, ""),
			Str("find"), Native2("find", func(e *Env, rx, text Value) Value {
				m := rx.Map().GetString("_rx").Any().(*regexp.Regexp).FindStringSubmatch(text.MustStr("find text", 0))
				mm := []Value{}
				for _, m := range m {
					mm = append(mm, Str(m))
				}
				return Array(mm...)
			}, ""),
			Str("findall"), Native3("findall", func(e *Env, rx, text, n Value) Value {
				m := rx.Map().GetString("_rx").Any().(*regexp.Regexp).FindAllStringSubmatch(text.MustStr("findall text", 0), int(n.IntDefault(-1)))
				mm := []Value{}
				for _, m := range m {
					for _, m := range m {
						mm = append(mm, Str(m))
					}
				}
				return Array(mm...)
			}, ""),
			Str("replace"), Native3("replace", func(e *Env, rx, text, newtext Value) Value {
				m := rx.Map().GetString("_rx").Any().(*regexp.Regexp).ReplaceAllString(text.MustStr("replace text", 0), newtext.MustStr("replace text new", 0))
				return Str(m)
			}, ""),
		)
		b := Map()
		b.Map().Parent = a.Map()
		return b
	}, "re(string) => creates a regular expression object")
	AddGlobalValue("error", func(env *Env, msg Value) Value {
		return Any(errors.New(msg.MustStr("error() message", 0)))
	}, "error(text)", "\tcreate an error")
	AddGlobalValue("iserror", func(env *Env) {
		_, ok := env.Get(0).Any().(error)
		env.A = Bool(ok)
	}, "iserror(value)", "\ttest whether value is an error")

	AddGlobalValue("json", Map(
		Str("stringify"), Native("stringify", func(env *Env) {
			env.A = Str(env.Get(0).JSONString())
		}, "$f(value) => json_string"),
		Str("parse"), Native1("parse", func(env *Env, js Value) Value {
			j := strings.TrimSpace(js.MustStr("json.parse() json string", 0))
			return gjsonConvert(gjson.Parse(j))
		}, "$f(json_string) => array"),
		Str("get"), Native3("get", func(env *Env, js, path, et Value) Value {
			j := strings.TrimSpace(js.MustStr("json.get() json string", 0))
			result := gjson.Get(j, path.MustStr("json.get() path", 0))
			if !result.Exists() {
				return et
			}
			return gjsonConvert(result)
		}, "$f(json_string, selector, default?) => bool|number|string|array"),
	))

	AddGlobalValue("sync", Map(
		Str("mutex"), Native("mutex", func(env *Env) { env.A = Any(&sync.Mutex{}) }, "$f() => creates a mutex"),
		Str("rwmutex"), Native("rwmutex", func(env *Env) { env.A = Any(&sync.RWMutex{}) }, "$f() => creates a read-write mutex"),
		Str("waitgroup"), Native("rwmutex", func(env *Env) { env.A = Any(&sync.WaitGroup{}) }, "$f() => creates a wait group"),
	))

	// Array related functions
	AddGlobalValue("append", func(env *Env, m, v Value) Value {
		a := m.MustMap("append()", 0)
		a.Set(Int(int64(a.Len())), v)
		return m
	}, "append(array, value) => append value to array")
	AddGlobalValue("concat", func(env *Env, a, b Value) Value {
		ma, mb := a.MustMap("concat()", 1), b.MustMap("concat()", 2)
		for _, b := range mb.Array() {
			ma.Set(Int(int64(ma.Len())), b)
		}
		return ma.Value()
	}, "append(array, value) => append value to array")
	AddGlobalValue("next", func(env *Env, m, k Value) Value {
		nk, nv := m.MustMap("next()", 0).Next(k)
		return Array(nk, nv)
	})
	AddGlobalValue("parent", func(env *Env, m Value) Value {
		return m.MustMap("parent()", 0).Parent.Value()
	})
	AddGlobalValue("keys", func(env *Env, m Value) Value {
		a := make([]Value, 0)
		m.MustMap("keys()", 0).Foreach(func(k, v Value) bool {
			a = append(a, k)
			return true
		})
		return Array(a...)
	})
	AddGlobalValue("iter", func(env *Env, m Value) Value {
		a := Map(
			Str("key"), Undef,
			Str("value"), Undef,
			Str("_src"), m.MustMap("iter()", 0).Value(),
			Str("next"), Native1("next", func(env *Env, self Value) Value {
				m := self.Map()
				var k, v Value
				if pk := m.GetString("key"); pk == Undef {
					k, v = m.GetString("_src").Map().Next(Nil)
				} else {
					k, v = m.GetString("_src").Map().Next(pk)
				}
				if k == Nil {
					return Bool(false)
				}
				m.Set(Str("key"), k)
				m.Set(Str("value"), v)
				return Bool(true)
			}),
		)
		b := Map()
		b.Map().Parent = a.Map()
		return b
	})
}

func mathMinMax(env *Env, msg string, max bool) {
	if len(env.Stack()) <= 0 {
		return
	}
	f, i, isInt := env.Get(0).MustNum(msg, 1).Num()
	if isInt {
		for ii := 1; ii < len(env.Stack()); ii++ {
			if x := env.Get(ii).MustNum(msg, ii+1).Int(); x >= i == max {
				i = x
			}
		}
		env.A = Int(i)
	} else {
		for i := 1; i < len(env.Stack()); i++ {
			if x, _, _ := env.Get(i).MustNum(msg, i+1).Num(); x >= f == max {
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

func gjsonConvert(r gjson.Result) Value {
	switch r.Type {
	case gjson.String:
		return Str(r.Str)
	case gjson.Number:
		return Float(r.Float())
	case gjson.True, gjson.False:
		return Bool(r.Bool())
	}
	if r.IsArray() {
		a := r.Array()
		x := make([]Value, len(a))
		for i, a := range a {
			x[i] = gjsonConvert(a)
		}
		return Array(x...)
	}
	if r.IsObject() {
		m := r.Map()
		x := NewMap(len(m))
		for k, v := range m {
			x.Set(Str(k), gjsonConvert(v))
		}
		return x.Value()
	}
	return Nil
}
