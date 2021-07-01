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

const Version int64 = 243

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
		g[k] = Interface(v)
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
		env.A = env.NewArray(globals...)
	}, "globals() => { g1, g2, ... }", "\tlist all global values")
	AddGlobalValue("doc", func(env *Env, f Value) Value {
		return String(f.MustFunc("doc", 0).DocString)
	}, "doc(function) => string", "\treturn function's documentation")
	AddGlobalValue("len", func(env *Env, v Value) Value {
		switch v.Type() {
		case VString:
			return Float(float64(len(v.rawStr())))
		case VMap:
			return Int(int64(v.Map().Len()))
		case VFunction:
			return Float(float64(v.Function().NumParams()))
		default:
			return Int(int64(reflectLen(v.Interface())))
		}
	})
	AddGlobalValue("acall", func(env *Env, f, a Value) Value {
		return Nil
		// res, err := f.MustFunc("acall", 0).Call(a.MustMap("acall", 1).Underlay...)
		// if err != nil {
		// 	panic(err)
		// }
		// return res
	})
	AddGlobalValue("mcall", func(env *Env, f, a Value) Value {
		return Nil
		// fn := f.MustFunc("mcall", 0)
		// m := buildCallMap(fn, Env{stack: &a.MustMap("mcall", 1).Underlay})
		// res, err := fn.CallMap(m)
		// if err != nil {
		// 	panic(err)
		// }
		// return res
	})
	AddGlobalValue("dbg_dumpstack", func(env *Env) {
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
	}, "dbg_dumpstack() => { v1, v2, v3, ... }")
	AddGlobalValue("dbg_kwargs", func(env *Env) {
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
	}, "dbg_kwargs() => { key1: value, key2: value2, ... }")
	AddGlobalValue("debug_locals", func(env *Env) {
		var r []Value
		start := env.StackOffset - uint32(env.DebugCaller.StackSize)
		for i, name := range env.DebugCaller.Locals {
			idx := start + uint32(i)
			r = append(r, Int(int64(idx)), String(name), (*env.stack)[idx])
		}
		env.A = env.NewArray(r...)
	}, "debug_locals() => { index1, name1, value1, i2, n2, v2, i3, n3, v3, ... }")
	AddGlobalValue("debug_globals", func(env *Env) {
		var r []Value
		for i, name := range env.Global.Func.Locals {
			r = append(r, Int(int64(i)), String(name), (*env.Global.Stack)[i])
		}
		env.A = env.NewArray(r...)
	}, "debug_globals() => { index1, name1, value1, i2, n2, v2, i3, n3, v3, ... }")
	AddGlobalValue("debug_set", func(env *Env, idx, value Value) Value {
		(*env.Global.Stack)[idx.MustNumber("debug_set", 1).Int()] = value
		return Nil
	}, "debug_set(idx, value)")
	AddGlobalValue("debug_stacktrace", func(env *Env, skip Value) Value {
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
			lines = append(lines, String(r.cls.Name), Int(int64(src)), Int(int64(r.cursor-1)))
		}
		return env.NewArray(lines...)
	}, "debug_stacktrace(skip) => { func_name1, line1, cursor1, n2, l2, c2, ... }")
	AddGlobalValue("narray", func(env *Env, n Value) Value {
		return env.NewArray(make([]Value, n.MustNumber("narray", 0).Int())...)
	}, "narray(n) => { nil, ..., nil }", "\treturn an array size of n, filled with nil")
	AddGlobalValue("type", func(env *Env) {
		env.A = String(env.Get(0).Type().String())
	}, "type(value) => string", "\treturn value's type")
	AddGlobalValue("pcall", func(env *Env, f Value) Value {
		a, err := f.MustFunc("pcall", 0).Call(env.Stack()[1:]...)
		if err == nil {
			return a
		}
		if err, ok := err.(*ExecError); ok {
			return Interface(err.r)
		}
		return Interface(err)
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
	AddGlobalValue("num", func(env *Env) {
		v := env.Get(0)
		switch v.Type() {
		case VNumber:
			env.A = v
		case VString:
			switch v := parser.NewNumberFromString(v.rawStr()); v.Type {
			case parser.Float:
				env.A = Float(v.FloatValue())
			case parser.Int:
				env.A = Int(v.IntValue())
			}
		default:
			env.A = Value{}
		}
	}, "num(value) => number", "\tconvert string to number")
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
		w := env.Get(0).Interface().(io.Writer)
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
		if env.Global.GetDeadsize() > 0 {
			r = io.LimitReader(r, env.Global.GetDeadsize())
		}
		for i := n.IntDefault(1); i > 0; i-- {
			var s string
			if _, err := fmt.Fscan(r, &s); err != nil {
				break
			}
			results = append(results, env.NewString(s))
		}
		return env.NewArray(results...)
	},
		"$f(prompt='', n=1) => { s1, s2, ..., sn }", "\tprint prompt and read N user inputs",
	)
	AddGlobalValue("INF", Float(math.Inf(1)))
	AddGlobalValue("PI", Float(math.Pi))
	AddGlobalValue("E", Float(math.E))
	AddGlobalValue("randomseed", func(env *Env) {
		rg.Lock()
		defer rg.Unlock()
		rg.Rand.Seed(env.Get(0).IntDefault(1))
	}, "randomseed(int)")
	AddGlobalValue("random", func(env *Env) {
		rg.Lock()
		defer rg.Unlock()
		switch len(env.Stack()) {
		case 2:
			af, ai, aIsInt := env.Get(0).MustNumber("random #", 1).Num()
			bf, bi, bIsInt := env.Get(1).MustNumber("random #", 2).Num()
			if aIsInt && bIsInt {
				env.A = Int(int64(rg.Intn(int(bi-ai+1))) + ai)
			} else {
				env.A = Float(rg.Float64()*(bf-af) + af)
			}
		case 1:
			env.A = Int(int64(rg.Intn(int(env.Get(0).MustNumber("random", 0).Int()))) + 1)
		default:
			env.A = Float(rg.Float64())
		}
	},
		"$f() => [0,1]",
		"$f(n) => [1, n]",
		"$f(a, b) => [a, b]")
	AddGlobalValue("bitand", func(env *Env, a, b Value) Value {
		return Int(a.MustNumber("bitand", 1).Int() & b.MustNumber("bitand", 2).Int())
	})
	AddGlobalValue("bitor", func(env *Env, a, b Value) Value {
		return Int(a.MustNumber("bitor", 1).Int() | b.MustNumber("bitor", 2).Int())
	})
	AddGlobalValue("bitxor", func(env *Env, a, b Value) Value {
		return Int(a.MustNumber("bitxir", 1).Int() ^ b.MustNumber("bitxir", 2).Int())
	})
	AddGlobalValue("bitrsh", func(env *Env, a, b Value) Value {
		return Int(a.MustNumber("bitrsh", 1).Int() >> b.MustNumber("bitrsh", 2).Int())
	})
	AddGlobalValue("bitlsh", func(env *Env, a, b Value) Value {
		return Int(a.MustNumber("bitlsh", 1).Int() << b.MustNumber("bitlsh", 2).Int())
	})
	AddGlobalValue("sqrt", func(env *Env, v Value) Value {
		return Float(math.Sqrt(v.MustNumber("sqrt", 0).Float()))
	})
	AddGlobalValue("floor", func(env *Env, v Value) Value {
		return Float(math.Floor(v.MustNumber("floor", 0).Float()))
	})
	AddGlobalValue("ceil", func(env *Env, v Value) Value {
		return Float(math.Ceil(v.MustNumber("ceil", 0).Float()))
	})
	AddGlobalValue("min", func(env *Env) { mathMinMax(env, "min #", false) }, "max(a, b, ...) => largest_number")
	AddGlobalValue("max", func(env *Env) { mathMinMax(env, "max #", true) }, "min(a, b, ...) => smallest_number")
	AddGlobalValue("pow", func(env *Env, a, b Value) Value {
		af, ai, aIsInt := a.MustNumber("pow", 1).Num()
		bf, bi, bIsInt := b.MustNumber("pow", 2).Num()
		if aIsInt && bIsInt {
			return Int(ipow(ai, bi))
		}
		return Float(math.Pow(af, bf))
	}, "min(a, b, ...) => smallest_number")
	AddGlobalValue("abs", func(env *Env) {
		switch f, i, isInt := env.Get(0).MustNumber("abs", 0).Num(); {
		case isInt && i < 0:
			env.A = Int(-i)
		case isInt && i >= 0:
			env.A = Int(i)
		default:
			env.A = Float(math.Abs(f))
		}
	})
	AddGlobalValue("str", func(env *Env, v, format Value) Value {
		return env.NewString(fmt.Sprintf(format.StringDefault("%v"), v.Interface()))
	},
		"str(value, format='%v') => string", "\tconvert value to string using format",
	)
	AddGlobalValue("int", func(env *Env) {
		switch v := env.Get(0); v.Type() {
		case VNumber:
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
			env.A = Interface(time.Date(
				int(env.Get(0).IntDefault(1970)), time.Month(env.Get(1).IntDefault(1)), int(env.Get(2).IntDefault(1)),
				int(env.Get(3).IntDefault(0)), int(env.Get(4).IntDefault(0)), int(env.Get(5).IntDefault(0)),
				int(env.Get(6).IntDefault(0)), loc,
			))
		} else {
			env.A = Interface(time.Now())
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
		os.Exit(int(code.MustNumber("exit", 0).Int()))
		return Nil
	}, "exit(code)")
	AddGlobalValue("char", func(env *Env) {
		env.A = env.NewString(string(rune(env.Get(0).MustNumber("char", 0).Int())))
	}, "char(number) => string")
	AddGlobalValue("unicode", func(env *Env) {
		r, sz := utf8.DecodeRuneInString(env.Get(0).MustString("unicode", 0))
		env.A = env.NewArray(Int(int64(r)), Int(int64(sz)))
	}, "unicode(one_char_string) => { char_unicode, width_in_bytes }")
	AddGlobalValue("substr", func(env *Env, s, i, j Value) Value {
		ss := s.MustString("substr", 0)
		return String(ss[i.MustNumber("substr start", 0).Int():j.IntDefault(int64(len(ss)))])
	})
	AddGlobalValue("chars", func(env *Env, s, n Value) Value {
		var r []Value
		max := n.IntDefault(0)
		for s := s.MustString("chars", 0); len(s) > 0; {
			_, sz := utf8.DecodeRuneInString(s)
			if sz == 0 {
				break
			}
			r = append(r, String(s[:sz]))
			if max > 0 && len(r) == int(max) {
				break
			}
			s = s[sz:]
		}
		return env.NewArray(r...)
	}, "chars(string) => { char1, char2, ... }", "chars(string, max) => { char1, char2, ..., char_max }",
		"\tbreak a string into (at most 'max') chars, e.g.:",
		"\tchars('a中c') => { 'a', '中', 'c' }",
		"\tchars('a中c', 1) => { 'a' }",
	)
	AddGlobalValue("match", func(env *Env, in, r, n Value) Value {
		rx, err := regexp.Compile(r.MustString("match #", 2))
		if err != nil {
			return env.NewArray(Interface(err))
		}
		m := rx.FindAllStringSubmatch(in.MustString("match #", 1), int(n.IntDefault(-1)))
		mm := []Value{}
		for _, m := range m {
			for _, m := range m {
				mm = append(mm, String(m))
			}
		}
		return env.NewArray(mm...)
	}, "match(string, regex, n=-1) => { match1, match2, ..., matchn }")
	AddGlobalValue("startswith", func(env *Env, t, p Value) Value {
		return Bool(strings.HasPrefix(t.MustString("startswith", 0), p.MustString("startswith prefix", 0)))
	}, "startswith(text, prefix) => bool")
	AddGlobalValue("endswith", func(env *Env, t, s Value) Value {
		return Bool(strings.HasSuffix(t.MustString("endswith", 0), s.MustString("endswith suffix", 0)))
	}, "endswith(text, suffix) => bool")
	AddGlobalValue("stricmp", func(env *Env, a, b Value) Value {
		return Bool(strings.EqualFold(a.MustString("stricmp #", 1), b.MustString("stricmp #", 2)))
	}, "stricmp(text1, text2) => bool", "\tcompare two strings case insensitively")
	AddGlobalValue("trimspace", func(env *Env, txt Value) Value {
		return String(strings.TrimSpace(txt.MustString("trimspace", 0)))
	})
	AddGlobalValue("trim", func(env *Env, txt, p, t Value) Value {
		switch a, cutset := txt.MustString("trim", 0), p.StringDefault(" \n\t\r"); t.StringDefault("") {
		case "left", "l":
			return String(strings.TrimLeft(a, cutset))
		case "right", "r":
			return String(strings.TrimRight(a, cutset))
		case "prefix", "start":
			return String(strings.TrimPrefix(a, cutset))
		case "suffix", "end":
			return String(strings.TrimSuffix(a, cutset))
		default:
			return String(strings.Trim(a, cutset))
		}
	},
		"$f(text) => string", "\ttrim spaces",
		"$f(text, cutset) => string", "\ttrim chars inside 'cutset'",
		"$f(text, cutset, 'left'|'right') => string", "\ttrim right/left chars inside 'cutset'",
		"$f(text, pattern, 'suffix'|'prefix') => string", "\ttrim prefix/suffix",
	)
	AddGlobalValue("replace", func(env *Env, text, src, dst Value) Value {
		a := text.MustString("replace text", 0)
		rx, err := regexp.Compile(src.MustString("replace source", 0))
		if err != nil {
			return Interface(err)
		}
		switch f := env.Get(2); f.Type() {
		case VString:
			return env.NewString(rx.ReplaceAllString(a, f.rawStr()))
		case VFunction:
			return env.NewString(rx.ReplaceAllStringFunc(a, func(in string) string {
				v, err := f.Function().Call(String(in))
				if err != nil {
					panic(err)
				}
				return v.String()
			}))
		}
		return Nil
	},
		"replace(text, regex, newtext|callback) => string",
		"\tcallback will be called in such way: new_string = f(captured_string)",
	)
	AddGlobalValue("split", func(env *Env, text, sep Value) Value {
		x := strings.Split(text.MustString("split", 0), sep.MustString("split sep", 0))
		v := make([]Value, len(x))
		for i := range x {
			v[i] = String(x[i])
		}
		return env.NewArray(v...)
	}, "split(text, sep) => { part1, part2, ... }")
	AddGlobalValue("strpos", func(env *Env, txt, n, t Value) Value {
		a, b := txt.MustString("strpos #", 1), n.MustString("strpos #", 2)
		if t.StringDefault("") == "last" {
			return Int(int64(strings.LastIndex(a, b)) + 1)
		}
		return Int(int64(strings.Index(a, b)) + 1)
	},
		"strpos(text, needle) => int", "\tfirst occurrence of needle in text, or 0 if not found",
		"strpos(text, needle, 'last') => int", "\tlast occurrence of needle in text",
	)
	AddGlobalValue("format", func(env *Env, p Value) Value {
		f := strings.Replace(p.MustString("format text", 0), "%", "%%", -1)
		f = strings.Replace(f, "{}", "%v", -1)
		return env.NewString(fmt.Sprintf(f, env.StackInterface()[1:]...))
	}, "format(pattern, a1, a2, ...)", "\t'{}' is the placeholder, no need to escape '%'")
	AddGlobalValue("error", func(env *Env, msg Value) Value {
		return Interface(errors.New(msg.MustString("error", 0)))
	}, "error(text)", "\tcreate an error")
	AddGlobalValue("iserror", func(env *Env) {
		_, ok := env.Get(0).Interface().(error)
		env.A = Bool(ok)
	}, "iserror(value)", "\ttest whether value is an error")
	AddGlobalValue("json", func(env *Env) {
		env.A = env.NewString(env.Get(0).JSONString())
	}, "json(v) => json_string")
	AddGlobalValue("json_get", func(env *Env, js, path, et Value) Value {
		cv := func(r gjson.Result) Value {
			switch r.Type {
			case gjson.String:
				return String(r.Str)
			case gjson.Number:
				return Float(r.Float())
			case gjson.True, gjson.False:
				return Bool(r.Bool())
			}
			return String(r.Raw)
		}
		j := strings.TrimSpace(js.MustString("json_get #", 1))
		result := gjson.Get(j, path.MustString("json_get #", 2))

		if expectedType := et.StringDefault(""); expectedType != "" {
			if !strings.EqualFold(result.Type.String(), expectedType) {
				return Nil
			}
		}

		if result.IsArray() {
			a := result.Array()
			tmp := make([]Value, len(a))
			for i := range a {
				tmp[i] = cv(a[i])
			}
			return env.NewArray(tmp...)
		}
		return cv(result)
	}, "$f(json_string, selector, nil|expected_type) => true|false|number|string|array|object_string")
	AddGlobalValue("next", func(env *Env, m, k Value) Value {
		nk, nv := m.MustMap("next", 0).Next(k)
		return Array(nk, nv)
	})
}

func mathMinMax(env *Env, msg string, max bool) {
	if len(env.Stack()) <= 0 {
		return
	}
	f, i, isInt := env.Get(0).MustNumber(msg, 1).Num()
	if isInt {
		for ii := 1; ii < len(env.Stack()); ii++ {
			if x := env.Get(ii).MustNumber(msg, ii+1).Int(); x >= i == max {
				i = x
			}
		}
		env.A = Int(i)
	} else {
		for i := 1; i < len(env.Stack()); i++ {
			if x, _, _ := env.Get(i).MustNumber(msg, i+1).Num(); x >= f == max {
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
