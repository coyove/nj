package script

import (
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"reflect"
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

const Version int64 = 241

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
	AddGlobalValue("True", _interface(true))
	AddGlobalValue("False", _interface(false))
	AddGlobalValue("globals", func(env *Env) {
		env.A = Int(int64(len(g)))
		keys := make([]string, 0, len(g))
		for k := range g {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			env.V = append(env.V, g[k])
		}
	}, "globals() => n, g1, g2, ...", "\tlist all global values")
	AddGlobalValue("doc", func(env *Env) {
		env.A = String(env.In(0, VFunction).Function().DocString)
	}, "doc(function) => string", "\treturn function's documentation")
	AddGlobalValue("debug_state", func(env *Env) {
		code := env.Debug.Caller.Code.Code
		c := env.Debug.Cursor
		for ; int(c) < len(code); c++ {
			op, _, _ := splitInst(code[c])
			if op == OpRet {
				c++
				break
			}
		}
		start := env.StackOffset - uint32(env.Debug.Caller.StackSize)
		copystack := append([]Value{}, (*env.stack)[start:env.StackOffset]...)
		env.A = Interface(&DebugState{c, copystack})
	}, "debug_state() => state", "\tget current function's executing state")
	AddGlobalValue("debug_resume", func(env *Env) {
		var (
			f      = env.In(0, VFunction).Function()
			cursor uint32
			newEnv = *env
			stack  []Value
		)
		if state, ok := env.Get(1).Interface().(*DebugState); ok {
			stack, cursor = state.Stack, state.Cursor
		}
		newEnv.StackOffset = uint32(len(*newEnv.stack))
		*newEnv.stack = append(*newEnv.stack, stack...)
		newEnv.grow(int(f.StackSize))
		if cursor >= uint32(f.Code.Len()) {
			panicf("cursor overflowed")
		}
		env.A, env.V = InternalExecCursorLoop(newEnv, f, cursor)
	}, "$f(function, state) => a1, ..., an, new_state", "\texecute the function using the given state and return the new state")
	AddGlobalValue("debug_locals", func(env *Env) {
		var r []Value
		start := env.StackOffset - uint32(env.Debug.Caller.StackSize)
		for i, name := range env.Debug.Caller.Locals {
			idx := start + uint32(i)
			r = append(r, Int(int64(idx)), String(name), (*env.stack)[idx])
		}
		env.Return(r...)
	}, "debug_locals() => index1, name1, value1, i2, n2, v2, i3, n3, v3, ...")
	AddGlobalValue("debug_globals", func(env *Env) {
		var r []Value
		for i, name := range env.Global.Func.Locals {
			r = append(r, Int(int64(i)), String(name), (*env.Global.Stack)[i])
		}
		env.Return(r...)
	}, "debug_globals() => index1, name1, value1, i2, n2, v2, i3, n3, v3, ...")
	AddGlobalValue("debug_set", func(env *Env) {
		(*env.Global.Stack)[env.In(0, VNumber).Int()] = env.Get(1)
	}, "debug_set(idx, value)")
	AddGlobalValue("debug_stacktrace", func(env *Env) {
		stacks := env.Debug.Stacktrace
		lines := make([]Value, 0, len(stacks))
		for i := len(stacks) - 1 - int(env.InInt(0, 0)); i >= 0; i-- {
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
		env.Return(lines...)
	}, "debug_stacktrace(skip) => func_name1, line1, cursor1, n2, l2, c2, ...")
	AddGlobalValue("narray", func(env *Env) {
		n := env.In(0, VNumber).Int()
		env.Return2(Int(n), make([]Value, n)...)
	}, "narray(n) => n, ...array", "\treturn an array size of n, filled with nil")
	AddGlobalValue("array", func(env *Env) {
		env.Return2(Int(int64(env.Size())), append([]Value{}, env.Stack()...)...)
	}, "array(a, b, c, ...) => n, a, b, c, ...", "\treturn the number of input arguments, followed by the arugments themselves")
	AddGlobalValue("type", func(env *Env) {
		env.A = String(env.Get(0).Type().String())
	}, "type(value) => string", "\treturn value's type")
	AddGlobalValue("pcall", func(env *Env) {
		a, v, err := env.In(0, VFunction).Function().Call(env.Stack()[1:]...)
		if err == nil {
			env.Return2(Bool(true), append([]Value{a}, v...)...)
		} else {
			switch err := err.(type) {
			case *ExecError:
				env.Return2(Bool(false), Interface(err.r))
			default:
				env.Return2(Bool(false), Interface(err))
			}
		}
	}, "pcall(function, arg1, arg2, ...)",
		"\texecute the function, catch panic and return: 0, panic_as_error",
		"\tif everything went well, return what function returned: 1, ...")
	AddGlobalValue("select", func(env *Env) {
		env.A = Value{}
		switch a := env.Get(0); a.Type() {
		case VString:
			env.A = Int(int64(env.Size() - 1))
		case VNumber:
			if u, idx := env.CopyStack()[1:], int(a.Int())-1; idx < len(u) {
				env.Return(u[idx:]...)
			}
		}
	}, "select(n, ...)", "\tlua-style select function")
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
			switch v := parser.NewNumberFromString(v._str()); v.Type {
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
		w, ok := env.In(0, VInterface).Interface().(io.Writer)
		if !ok {
			panicf("invalid stdout")
		}
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
		env.Return(results...)
	},
		"$f(prompt) => string", "\tprint prompt and read one user input",
		"$f(n) => s1, s2, ..., sn", "\tread N user input",
		"$f(prompt, n) => s1, s2, ..., sn", "\tprint prompt and read N user input",
	)
	AddGlobalValue("INF", Float(math.Inf(1)))
	AddGlobalValue("PI", Float(math.Pi))
	AddGlobalValue("E", Float(math.E))
	AddGlobalValue("randomseed", func(env *Env) {
		rg.Lock()
		rg.Rand.Seed(env.InInt(0, 1))
		rg.Unlock()
	}, "randomseed(int)")
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
	},
		"$f() => [0,1]",
		"$f(n) => [1, n]",
		"$f(a, b) => [a, b]")
	AddGlobalValue("bitand", func(env *Env) { env.A = Int(env.In(0, VNumber).Int() & env.In(1, VNumber).Int()) })
	AddGlobalValue("bitor", func(env *Env) { env.A = Int(env.In(0, VNumber).Int() | env.In(1, VNumber).Int()) })
	AddGlobalValue("bitxor", func(env *Env) { env.A = Int(env.In(0, VNumber).Int() ^ env.In(1, VNumber).Int()) })
	AddGlobalValue("bitrsh", func(env *Env) { env.A = Int(env.In(0, VNumber).Int() >> env.In(1, VNumber).Int()) })
	AddGlobalValue("bitlsh", func(env *Env) { env.A = Int(env.In(0, VNumber).Int() << env.In(1, VNumber).Int()) })
	AddGlobalValue("sqrt", func(env *Env) { env.A = Float(math.Sqrt(env.In(0, VNumber).Float())) })
	AddGlobalValue("floor", func(env *Env) { env.A = Float(math.Floor(env.In(0, VNumber).Float())) })
	AddGlobalValue("ceil", func(env *Env) { env.A = Float(math.Ceil(env.In(0, VNumber).Float())) })
	AddGlobalValue("min", func(env *Env) { mathMinMax(env, false) }, "max(a, b, ...) => number")
	AddGlobalValue("max", func(env *Env) { mathMinMax(env, true) }, "min(a, b, ...) => number")
	AddGlobalValue("abs", func(env *Env) {
		switch f, i, isInt := env.In(0, VNumber).Num(); {
		case isInt && i < 0:
			env.A = Int(-i)
		case isInt && i >= 0:
			env.A = Int(i)
		default:
			env.A = Float(math.Abs(f))
		}
	})
	AddGlobalValue("str", func(env *Env) {
		switch v := env.Get(0); v.Type() {
		case VString:
			env.A = env.NewString(fmt.Sprintf(env.InStr(1, "%s"), v.String()))
		default:
			env.A = env.NewString(fmt.Sprintf(env.InStr(1, "%v"), v.Interface()))
		}
	},
		"str(value) => string", "\tconvert value to string",
		"str(value, format) => string", "\tconvert value to string using format",
	)
	AddGlobalValue("int", func(env *Env) {
		switch v := env.Get(0); v.Type() {
		case VNumber:
			env.A = Int(v.Int())
		default:
			v, err := strconv.ParseInt(v.String(), 0, 64)
			env.Return2(Int(v), Interface(err))
		}
	}, "int(value) => integer", "\tconvert value to integer number (int64)")
	AddGlobalValue("time", func(env *Env) {
		switch env.InStr(0, "") {
		case "nano":
			env.A = Int(time.Now().UnixNano())
		case "micro":
			env.A = Int(time.Now().UnixNano() / 1e3)
		case "milli":
			env.A = Int(time.Now().UnixNano() / 1e6)
		default:
			env.A = Int(time.Now().Unix())
		}
	},
		"time('nano'|'micro'|'milli') => int", "\tunix timestamp in (nano|micro|milli)seconds",
		"time() => int", "\tunix timestamp in seconds",
	)
	AddGlobalValue("sleep", func(env *Env) {
		if env.Get(0).Type() == VString {
			d, _ := time.ParseDuration(env.InStr(0, ""))
			time.Sleep(d)
		} else {
			time.Sleep(time.Duration(env.In(0, VNumber).Int()) * time.Millisecond)
		}
	}, "sleep(milliseconds|duration_string)")
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
	},
		"Go_time() => time.Time",
		"Go_time(year, month, day, h, m, s, nanoseconds, 'local'|'utc') => time.Time")
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
	},
		"clock('nano'|'micro'|'milli') => int", "\t(nano|micro|milli)seconds since startup",
		"clock() => int", "\tseconds since startup",
	)
	AddGlobalValue("exit", func(env *Env) { os.Exit(int(env.InInt(0, 0))) }, "exit(Code)")
	AddGlobalValue("char", func(env *Env) {
		env.A = env.NewString(string(rune(env.In(0, VNumber).Int())))
	}, "char(number) => string")
	AddGlobalValue("unicode", func(env *Env) {
		r, sz := utf8.DecodeRuneInString(env.In(0, VString)._str())
		env.Return2(Int(int64(r)), Int(int64(sz)))
	}, "unicode(one_char_string) => char_unicode")
	AddGlobalValue("chars", func(env *Env) {
		var r []Value
		var max = env.InInt(1, 0)
		for s := env.In(0, VString)._str(); len(s) > 0; {
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
		env.Return(r...)
	}, "chars(string) => char1, char2, ...",
		"chars(string, max) => char1, char2, ..., char_max",
		"\tbreak a string into (at most 'max') chars, e.g.:",
		"\tchars('a中c') => 'a', '中', 'c'",
		"\tchars('a中c', 1) => 'a'",
	)
	AddGlobalValue("match", func(env *Env) {
		rx, err := regexp.Compile(env.In(1, VString)._str())
		if err != nil {
			env.A = Interface(err)
			return
		}
		m := rx.FindAllStringSubmatch(env.Get(0).String(), int(env.InInt(2, -1)))
		mm := []Value{Int(int64(len(m)))}
		for _, m := range m {
			for _, m := range m {
				mm = append(mm, String(m))
			}
		}
		env.Return(mm...)
	}, "match(string, regex) => n_or_err, match1, match2, ..., matchn")
	AddGlobalValue("startswith", func(env *Env) {
		env.A = Bool(strings.HasPrefix(env.In(0, VString)._str(), env.In(1, VString)._str()))
	}, "startswith(text, prefix) => bool")
	AddGlobalValue("endswith", func(env *Env) {
		env.A = Bool(strings.HasSuffix(env.In(0, VString)._str(), env.In(1, VString)._str()))
	}, "endswith(text, suffix) => bool")
	AddGlobalValue("stricmp", func(env *Env) {
		env.A = Bool(strings.EqualFold(env.In(0, VString)._str(), env.In(1, VString)._str()))
	}, "stricmp(text1, text2) => bool", "\tcompare two strings case insensitively")
	AddGlobalValue("trim", func(env *Env) {
		switch a, cutset := env.In(0, VString)._str(), env.InStr(1, " \n\t\r"); env.InStr(2, "") {
		case "left", "l":
			env.A = String(strings.TrimLeft(a, cutset))
		case "right", "r":
			env.A = String(strings.TrimRight(a, cutset))
		case "prefix", "start":
			env.A = String(strings.TrimPrefix(a, cutset))
		case "suffix", "end":
			env.A = String(strings.TrimSuffix(a, cutset))
		default:
			env.A = String(strings.Trim(a, cutset))
		}
	},
		"$f(text) => string", "\ttrim spaces",
		"$f(text, cutset) => string", "\ttrim chars inside 'cutset'",
		"$f(text, cutset, 'left'|'right') => string", "\ttrim right/left chars inside 'cutset'",
		"$f(text, pattern, 'suffix'|'prefix') => string", "\ttrim prefix/suffix",
	)
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
				v, _, err := f.Function().Call(env.NewString(in))
				if err != nil {
					panic(err)
				}
				return v.String()
			}))
		}
	},
		"replace(text, regex, newtext) => string",
		"replace(text, regex, callback) => string",
		"\tcallback will be called in such way: new_string = f(captured_string)",
	)
	AddGlobalValue("split", func(env *Env) {
		x := strings.Split(env.In(0, VString)._str(), env.In(1, VString)._str())
		v := make([]Value, len(x))
		for i := range x {
			v[i] = String(x[i])
		}
		env.Return(v...)
	}, "split(text, sep) => part1, part2, ...")
	AddGlobalValue("strpos", func(env *Env) {
		a, b := env.In(0, VString)._str(), env.In(1, VString)._str()
		if env.InStr(2, "") == "last" {
			env.A = Int(int64(strings.LastIndex(a, b)) + 1)
		} else {
			env.A = Int(int64(strings.Index(a, b)) + 1)
		}
	},
		"strpos(text, needle) => int", "\tfirst occurrence of needle in text, or 0 if not found",
		"strpos(text, needle, 'last') => int", "\tlast occurrence of needle in text",
	)
	AddGlobalValue("format", func(env *Env) {
		f := strings.Replace(env.In(0, VString)._str(), "%", "%%", -1)
		f = strings.Replace(f, "{}", "%v", -1)
		env.A = env.NewString(fmt.Sprintf(f, env.StackInterface()[1:]...))
	}, "format(pattern, a1, a2, ...)", "\t'{}' is the placeholder, no need to escape '%'")
	AddGlobalValue("error", func(env *Env) {
		env.A = Interface(errors.New(env.InStr(0, "")))
	}, "error(text)", "\tcreate an error")
	AddGlobalValue("iserror", func(env *Env) {
		_, ok := env.Get(0).Interface().(error)
		env.A = Bool(ok)
	}, "iserror(value)", "\ttest whether value is an error")
	AddGlobalValue("json", func(env *Env) {
		env.A = env.NewString(env.Get(0).JSONString())
	}, "json(v) => json_string")
	AddGlobalValue("json_get", func(env *Env) {
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
		j := strings.TrimSpace(env.In(0, VString)._str())
		result := gjson.Get(j, env.In(1, VString)._str())

		if expectedType := env.InStr(2, ""); expectedType != "" {
			if !strings.EqualFold(result.Type.String(), expectedType) {
				return
			}
		}

		env.A = cv(result)
		if result.IsArray() {
			a := result.Array()
			tmp := make([]Value, len(a))
			for i := range a {
				tmp[i] = cv(a[i])
			}
			env.Return2(Int(int64(len(a))), tmp...)
		}
	},
		"$f(json_string, selector) => true|false|number|string",
		"$f(json_string, selector) => n, ...array",
		"$f(json_string, selector) => object_string",
		"$f(json_string, selector, expected_type) => value",
	)
	AddGlobalValue("dict", func(env *Env) {
		if env.Size() > 0 {
			env.checkRemainStackSize(env.Size())
			env.A = _interface(env.CopyStack())
		} else {
			env.A = env.A
		}
	}, "dict(a1, a2, ..., an) => []Value{ a1, ..., an }",
		"dict(k1 = v1, ..., kn = vn) => map[string]Value{ k1: v1, ..., kn: vn }")
	AddGlobalValue("iter", func(env *Env) {
		iter := reflect.ValueOf(env.Get(0).Interface()).MapRange()
		env.A = Interface(iter)
	}, "iter(map[string]Value) => map_iterator",
		"\texample:",
		"\twhile map_iter.next() do",
		"\t\tprintln(map_iter.key(), map_iter.value())",
		"\tend")
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
