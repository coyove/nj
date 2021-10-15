package script

import (
	"errors"
	"fmt"
	"io"
	"math"
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
	"github.com/coyove/script/typ"
	"github.com/tidwall/gjson"
)

const Version int64 = 301

var (
	g   = map[string]Value{}
	now int64
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
		g[k] = Val(v)
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
	AddGlobalValue("_", func(env *Env) { env.A = env.Get(0) }, "nop(a) => a")
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
	AddGlobalValue("doc", func(env *Env, f, doc Value) Value {
		if doc == Nil {
			return Str(f.MustFunc("").DocString)
		}
		f.MustFunc("").DocString = doc.String()
		return doc
	}, "doc(function) => string", "\treturn function's documentation",
		"doc(function, docstring)", "\tupdate function's documentation")
	AddGlobalValue("new", func(env *Env, v, a Value) Value {
		m := *v.MustMap("")
		m.hashItems = append([]hashItem{}, m.hashItems...)
		m.items = append([]Value{}, m.items...)
		if a.Type() != typ.Table {
			return (&Table{Parent: &m}).Value()
		}
		a.Table().Parent = &m
		return a
	})
	AddGlobalValue("prototype", g["new"])
	AddGlobalValue("len", func(env *Env, v Value) Value {
		switch v.Type() {
		case typ.String:
			return Int(int64(len(v.Str())))
		case typ.Table:
			return Int(int64(v.Table().Len()))
		case typ.Func:
			return Int(int64(v.Func().NumParams))
		case typ.Number, typ.Nil, typ.Bool:
			return panicf("can't measure length of %v", v.Type())
		default:
			return Int(int64(reflectLen(v.Interface())))
		}
	})
	AddGlobalValue("eval", func(env *Env, s, g Value) Value {
		var m map[string]interface{}
		if g.Type() == typ.Table {
			m = map[string]interface{}{}
			g.Table().Foreach(func(k, v Value) bool {
				m[k.String()] = v.Interface()
				return true
			})
		}
		wrap := func(err error) error { return fmt.Errorf("panic inside: %v", err) }
		f, err := LoadString(s.MustStr(""), &CompileOptions{GlobalKeyValues: m})
		if err != nil {
			panic(wrap(err))
		}
		v, err := f.Run()
		if err != nil {
			panic(wrap(err))
		}
		return v
	}, "eval(string, globals) => evaluate the string")
	AddGlobalValue("apply", func(env *Env, f Value) Value {
		res, err := f.MustFunc("").Call(env.Stack()[1:]...)
		if err != nil {
			panic(err)
		}
		return res
	}, "apply(function, array) => call function using arguments in array")

	AddGlobalValue("go", Map(
		Str("new"), Native1("new", func(env *Env, f Value) Value {
			wg := &sync.WaitGroup{}
			b := Map(
				Str("_f"), f.MustFunc("").WrappedValue(),
				Str("_r"), Undef,
				Str("_wg"), Val(wg),
				Str("start"), Native1("start", func(env *Env, t Value) Value {
					m := t.Table()
					wg := m.GetString("_wg").Interface().(*sync.WaitGroup)
					wg.Add(1)
					args := append([]Value{}, env.Stack()[1:]...)
					go func() {
						defer wg.Done()
						res, err := m.GetString("_f").WrappedFunc().Call(args...)
						if err != nil {
							panic(err)
						}
						m.SetString("_r", res)
					}()
					return Nil
				}),
				Str("wait"), Native1("wait", func(env *Env, t Value) Value {
					t.Table().GetString("_wg").Interface().(*sync.WaitGroup).Wait()
					return t.Table().GetString("_r")
				}),
			)
			return TableProto(b.Table())
		}),
	))

	// Debug libraries
	AddGlobalValue("debug", Map(
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
			(*env.Global.Stack)[idx.MustNum("").Int()] = value
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
	AddGlobalValue("table", Map(
		Str("makearray"), Native1("makearray", func(env *Env, n Value) Value {
			a := Array(make([]Value, n.MustNum("").Int())...)
			a.Table().count = 0
			return a
		}, "makearray(n) => { nil, ..., nil }", "\treturn a table array, preallocate space for n values"),
		Str("arraylen"), Native1("arraylen", func(env *Env, v Value) Value { return Int(int64(len(v.MustMap("").items))) }),
		Str("maplen"), Native1("maplen", func(env *Env, v Value) Value { return Int(int64(len(v.MustMap("").hashItems))) }),
		Str("keys"), Native1("keys", func(env *Env, m Value) Value {
			a := make([]Value, 0)
			m.MustMap("").Foreach(func(k, v Value) bool {
				a = append(a, k)
				return true
			})
			return Array(a...)
		}),
		Str("append"), Native2("append", func(env *Env, m, v Value) Value {
			a := m.MustMap("")
			a.Set(Int(int64(len(a.items))), v)
			return m
		}, "append(array, value) => append value to array"),
		Str("concat"), Native2("concat", func(env *Env, a, b Value) Value {
			ma, mb := a.MustMap(""), b.MustMap("")
			for _, b := range mb.ArrayPart() {
				ma.Set(Int(int64(len(ma.items))), b)
			}
			return ma.Value()
		}, "concat(array1, array2) => put elements from array2 to array1's end"),
	))
	AddGlobalValue("type", func(env *Env) { env.A = Str(env.Get(0).Type().String()) }, "type(value) => string", "\treturn value's type")
	AddGlobalValue("pcall", func(env *Env, f Value) Value {
		a, err := f.MustFunc("").Call(env.Stack()[1:]...)
		if err == nil {
			return a
		}
		if err, ok := err.(*ExecError); ok {
			return Val(err.r)
		}
		return Val(err)
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
		case typ.Number:
			env.A = v
		case typ.String:
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
	AddGlobalValue("stderr", func(env *Env) { env.A = _interface(env.Global.Stderr) }, "$f() => fd", "\treturn stderr fd")
	AddGlobalValue("stdin", func(env *Env) { env.A = _interface(env.Global.Stdin) }, "$f() => fd", "\treturn stdin fd")
	AddGlobalValue("print", func(env *Env) {
		for _, a := range env.Stack() {
			fmt.Fprint(env.Global.Stdout, a.String())
		}
		fmt.Fprintln(env.Global.Stdout)
	}, "print(a, b, c, ...)", "\tprint values, no space between them")
	AddGlobalValue("write", func(env *Env) {
		w := env.Get(0).Interface().(io.Writer)
		for _, a := range env.Stack()[1:] {
			fmt.Fprint(w, a.String())
		}
	}, "write(writer, a, b, c, ...)", "\twrite values to writer")
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
		"$f(prompt='', n=1) => { s1, s2, ..., sN }", "\tprint prompt and read N user inputs",
	)
	AddGlobalValue("math", MathLib)
	AddGlobalValue("int", func(env *Env) {
		switch v := env.Get(0); v.Type() {
		case typ.Number:
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
	AddGlobalValue("sleep", func(env *Env, milli Value) Value {
		time.Sleep(time.Duration(milli.IntDefault(0)) * time.Millisecond)
		return Nil
	}, "sleep(milliseconds)")
	AddGlobalValue("Go_time", func(env *Env) {
		if env.Size() > 0 {
			loc := time.UTC
			if env.Get(7).StringDefault("") == "local" {
				loc = time.Local
			}
			env.A = Val(time.Date(
				int(env.Get(0).IntDefault(1970)), time.Month(env.Get(1).IntDefault(1)), int(env.Get(2).IntDefault(1)),
				int(env.Get(3).IntDefault(0)), int(env.Get(4).IntDefault(0)), int(env.Get(5).IntDefault(0)),
				int(env.Get(6).IntDefault(0)), loc,
			))
		} else {
			env.A = Val(time.Now())
		}
	},
		"Go_time() => time.Time",
		"\treturn time.Time of current time",
		"Go_time(year, month, day, h, m, s, nanoseconds, 'local'|'utc') => time.Time",
		"\treturn time.Time constructed by the given arguments",
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
		os.Exit(int(code.MustNum("").Int()))
		return Nil
	}, "exit(code)")
	AddGlobalValue("chr", func(env *Env) { env.A = Rune(rune(env.Get(0).MustNum("").Int())) }, "chr(unicode) => string")
	AddGlobalValue("byte", func(env *Env, a Value) Value { return Byte(byte(a.MustNum("").Int())) }, "byte(int) => one byte string")
	AddGlobalValue("ord", func(env *Env) {
		r, _ := utf8.DecodeRuneInString(env.Get(0).MustStr(""))
		env.A = Int(int64(r))
	}, "$f(string) => unicode")
	AddGlobalValue("re", func(env *Env, r Value) Value {
		rx, err := regexp.Compile(r.MustStr(""))
		if err != nil {
			panic(err)
		}
		a := Map(
			Str("_rx"), Val(rx),
			Str("match"), Native2("match", func(e *Env, rx, text Value) Value {
				return Bool(rx.Table().GetString("_rx").Interface().(*regexp.Regexp).MatchString(text.MustStr("")))
			}, ""),
			Str("find"), Native2("find", func(e *Env, rx, text Value) Value {
				m := rx.Table().GetString("_rx").Interface().(*regexp.Regexp).FindStringSubmatch(text.MustStr(""))
				mm := []Value{}
				for _, m := range m {
					mm = append(mm, Str(m))
				}
				return Array(mm...)
			}, ""),
			Str("findall"), Native3("findall", func(e *Env, rx, text, n Value) Value {
				m := rx.Table().GetString("_rx").Interface().(*regexp.Regexp).FindAllStringSubmatch(text.MustStr(""), int(n.IntDefault(-1)))
				mm := []Value{}
				for _, m := range m {
					for _, m := range m {
						mm = append(mm, Str(m))
					}
				}
				return Array(mm...)
			}, ""),
			Str("replace"), Native3("replace", func(e *Env, rx, text, newtext Value) Value {
				m := rx.Table().GetString("_rx").Interface().(*regexp.Regexp).ReplaceAllString(text.MustStr(""), newtext.MustStr(""))
				return Str(m)
			}, ""),
		)
		b := Map()
		b.Table().Parent = a.Table()
		return b
	}, "re(string) => create a regular expression object")
	AddGlobalValue("error", func(env *Env, msg Value) Value {
		return Val(errors.New(msg.MustStr("")))
	}, "error(text)", "\tcreate an error")
	AddGlobalValue("iserror", func(env *Env) {
		_, ok := env.Get(0).Interface().(error)
		env.A = Bool(ok)
	}, "iserror(value) => bool", "\ttest whether value is an error")

	AddGlobalValue("json", Map(
		Str("stringify"), Native("stringify", func(env *Env) {
			env.A = Str(env.Get(0).JSONString())
		}, "$f(value) => json_string"),
		Str("parse"), Native1("parse", func(env *Env, js Value) Value {
			j := strings.TrimSpace(js.MustStr(""))
			return Val(gjson.Parse(j))
		}, "$f(json_string) => object"),
		Str("get"), Native3("get", func(env *Env, js, path, et Value) Value {
			j := strings.TrimSpace(js.MustStr("json string"))
			result := gjson.Get(j, path.MustStr("selector"))
			if !result.Exists() {
				return et
			}
			return Val(result)
		}, "$f(json_string, selector, default?) => object"),
	))

	AddGlobalValue("sync", Map(
		Str("mutex"), Native("mutex", func(env *Env) { env.A = Val(&sync.Mutex{}) }, "$f() => mutex"),
		Str("rwmutex"), Native("rwmutex", func(env *Env) { env.A = Val(&sync.RWMutex{}) }, "$f() => read-write mutex"),
		Str("waitgroup"), Native("waitgroup", func(env *Env) { env.A = Val(&sync.WaitGroup{}) }, "$f() => wait group"),
	))
	AddGlobalValue("next", func(env *Env, m, k Value) Value {
		nk, nv := m.MustMap("").Next(k)
		return Array(nk, nv)
	}, "next(map, start_key) => { next_key, next_value }", "\tfind next key-value pair in the map")
	AddGlobalValue("parent", func(env *Env, m Value) Value {
		return m.MustMap("").Parent.Value()
	}, "parent(table) => table", "\tfind given table's parent, or nil if not existed")
}
