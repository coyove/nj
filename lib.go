package script

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"math"
	"os"
	"regexp"
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

const Version int64 = 290

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
	AddGlobalValue("globals", func(env *Env) {
		r := NewTable(len(env.Global.Func.Locals))
		for i, name := range env.Global.Func.Locals {
			r.Set(Str(name), (*env.Global.Stack)[i])
		}
		env.A = r.Value()
	}, "globals() table", "\tlist all global values as key-value pairs")
	AddGlobalValue("doc", func(env *Env, f, doc Value) Value {
		if doc == Nil {
			return Str(f.MustFunc("").DocString)
		}
		f.MustFunc("").DocString = doc.String()
		return doc
	}, "doc(f: function) string", "\treturn function's documentation",
		"doc(f: function, docstring: string)", "\tupdate function's documentation")
	AddGlobalValue("new", func(env *Env, v, a Value) Value {
		m := v.MustTable("").Copy()
		if a.Type() != typ.Table {
			return (&Table{Parent: m}).Value()
		}
		a.Table().Parent = m
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
		if gt := g.MaybeTableGetString("globals"); gt.Type() == typ.Table {
			m = map[string]interface{}{}
			gt.Table().Foreach(func(k, v Value) bool {
				m[k.String()] = v.Interface()
				return true
			})
		}
		if !g.MaybeTableGetString("compileonly").IsFalse() {
			v, err := parser.Parse(s.MustStr(""), "")
			if err != nil {
				panic(err)
			}
			return Val(v)
		}
		wrap := func(err error) error { return fmt.Errorf("panic inside: %v", err) }
		p, err := LoadString(s.MustStr(""), &CompileOptions{GlobalKeyValues: m})
		if err != nil {
			panic(wrap(err))
		}
		v, err := p.Run()
		if err != nil {
			panic(wrap(err))
		}
		return v
	}, "eval(code: string, globals: table) value", "\tevaluate the string and return the executed reuslt")
	AddGlobalValue("apply", func(env *Env, f Value) Value {
		res, err := f.MustFunc("").Call(env.Stack()[1:]...)
		if err != nil {
			panic(err)
		}
		return res
	}, "apply(f: function, args: array) value", "\tcall function using an array of arguments")
	AddGlobalValue("closure", func(env *Env, f, m Value) Value {
		lambda := f.MustFunc("")
		return Map(
			Str("__source"), m,
			Str("__lambda"), lambda.Value(),
			Str("__call"), Native("<closure-"+lambda.Name+">", func(env *Env) {
				recv := env.Get(0).MustTable("")
				f := recv.GetString("__lambda").MustFunc("")
				f.MethodSrc = Nil
				stk := append([]Value{recv.GetString("__source")}, env.Stack()[1:]...)
				res, err := f.Call(stk...)
				if err != nil {
					panic(err)
				}
				env.A = res
			}),
		)
	}, "closure(v: value, lambda: function) value", "\tbind v to lambda, when lambda is called, v will be passed in as the first argument")

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
		}, "$f() array", "\treturn { index1, name1, value1, i2, n2, v2, i3, n3, v3, ... }"),
		Str("globals"), Native("globals", func(env *Env) {
			var r []Value
			for i, name := range env.Global.Func.Locals {
				r = append(r, Int(int64(i)), Str(name), (*env.Global.Stack)[i])
			}
			env.A = Array(r...)
		}, "$f() array", "\treturn { index1, name1, value1, i2, n2, v2, i3, n3, v3, ... }"),
		Str("set"), Native2("set", func(env *Env, idx, value Value) Value {
			(*env.Global.Stack)[idx.MustInt("")] = value
			return Nil
		}, "$f(idx: int, v: value)"),
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
		}, "$f(skip: int) array", "\treturn { func_name0, line1, cursor1, n2, l2, c2, ... }"),
	))
	AddGlobalValue("type", func(env *Env) { env.A = Str(env.Get(0).Type().String()) }, "type(v value) string", "\treturn value's type")
	AddGlobalValue("pcall", func(env *Env, f Value) Value {
		a, err := f.MustFunc("").Call(env.Stack()[1:]...)
		if err == nil {
			return a
		}
		if err, ok := err.(*ExecError); ok {
			return Val(err.r)
		}
		return Val(err)
	}, "pcall(f: function, ...arg: value) value", "\texecute the function, catch panic and return as error if any")
	AddGlobalValue("panic", func(env *Env) { panic(env.Get(0)) }, "panic(v: value)")
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
	}, "assert(v: value)", "\tpanic when value is falsy",
		"assert(v1: value, v2: value)", "\tpanic when two values are not equal",
		"assert(v1: value, v2: value, msg: string)", "\tpanic message when two values are not equal",
	)
	AddGlobalValue("int", func(env *Env) {
		env.A = Nil
		switch v := env.Get(0); v.Type() {
		case typ.Number:
			env.A = Int(v.Int())
		default:
			if v, err := strconv.ParseInt(v.String(), int(env.Get(1).IntDefault(0)), 64); err == nil {
				env.A = Int(v)
			}
		}
	}, "int(v: value) int", "\tconvert value to integer number (int64)")
	AddGlobalValue("float", func(env *Env) {
		v := env.Get(0)
		switch v.Type() {
		case typ.Number:
			env.A = v
		case typ.String:
			switch v := parser.Num(v.Str()); v.Type() {
			case parser.FLOAT:
				env.A = Float(v.Float())
			case parser.INT:
				env.A = Int(v.Int())
			}
		default:
			env.A = Value{}
		}
	}, "$f(v: value) float", "\tconvert string to number")
	AddGlobalValue("stdout", func(env *Env) { env.A = _interface(env.Global.Stdout) }, "$f() value", "\treturn stdout")
	AddGlobalValue("stderr", func(env *Env) { env.A = _interface(env.Global.Stderr) }, "$f() value", "\treturn stderr")
	AddGlobalValue("stdin", func(env *Env) { env.A = _interface(env.Global.Stdin) }, "$f() value", "\treturn stdin")
	AddGlobalValue("print", func(env *Env) {
		for _, a := range env.Stack() {
			fmt.Fprint(env.Global.Stdout, a.String())
		}
		fmt.Fprintln(env.Global.Stdout)
	}, "print(...args: value)", "\tprint values, no space between them")
	AddGlobalValue("write", func(env *Env) {
		w := env.Get(0).Interface().(io.Writer)
		for _, a := range env.Stack()[1:] {
			fmt.Fprint(w, a.String())
		}
	}, "write(writer: value, ...args: value)", "\twrite values to writer")
	AddGlobalValue("println", func(env *Env) {
		for _, a := range env.Stack() {
			fmt.Fprint(env.Global.Stdout, a.String(), " ")
		}
		fmt.Fprintln(env.Global.Stdout)
	}, "println(...args: value)", "\tprint values, insert space between each of them")
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
		"$f() array", "\tread all user inputs and return { input1, input2, ... }",
		"$f(prompt: string) array", "\tprint prompt and read all user inputs",
		"$f(prompt: string, n: int) array", "\tprint prompt and read at most n user inputs",
	)
	AddGlobalValue("time", func(env *Env) { env.A = Float(float64(time.Now().UnixNano()) / 1e9) }, "time() float", "\tunix timestamp in seconds")
	AddGlobalValue("sleep", func(env *Env) { time.Sleep(time.Duration(env.Get(0).MustFloat("") * float64(time.Second))) }, "sleep(sec: float)")
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
		"Go_time() value",
		"\treturn time.Time of current time",
		"Go_time(year: int, month: int, day: int, h: int, m: int, s: int, nanoseconds: int, loc: string) value",
		"\treturn time.Time constructed by the given arguments, loc defaults to 'local'",
	)
	AddGlobalValue("clock", func(env *Env, prefix Value) Value {
		x := time.Now()
		s := *(*[2]int64)(unsafe.Pointer(&x))
		return Float(float64(s[1]) / 1e9)
	}, "clock() float", "\tseconds since startup (monotonic clock)")
	AddGlobalValue("exit", func(env *Env) { os.Exit(int(env.Get(0).MustInt(""))) }, "exit(code: int)")
	AddGlobalValue("chr", func(env *Env) { env.A = Rune(rune(env.Get(0).MustInt(""))) }, "chr(code: int) string")
	AddGlobalValue("byte", func(env *Env, a Value) Value { return Byte(byte(a.MustInt(""))) }, "byte(code: int) string")
	AddGlobalValue("ord", func(env *Env) {
		r, _ := utf8.DecodeRuneInString(env.Get(0).MustStr(""))
		env.A = Int(int64(r))
	}, "$f(s: string) int")

	AddGlobalValue("re", Map(
		Str("__call"), Native2("", func(env *Env, re, r Value) Value {
			rx, err := regexp.Compile(r.MustStr(""))
			if err != nil {
				panic(err)
			}
			return TableProto(re.MustTable(""), Str("_rx"), Val(rx))
		}, "$f(regex: string) table", "\tcreate a regular expression object"),
		Str("match"), Native2("match", func(e *Env, rx, text Value) Value {
			return Bool(rx.Table().GetString("_rx").Interface().(*regexp.Regexp).MatchString(text.MustStr("")))
		}, "$f({re}: value, text: string) bool"),
		Str("find"), Native2("find", func(e *Env, rx, text Value) Value {
			m := rx.Table().GetString("_rx").Interface().(*regexp.Regexp).FindStringSubmatch(text.MustStr(""))
			mm := []Value{}
			for _, m := range m {
				mm = append(mm, Str(m))
			}
			return Array(mm...)
		}, "$f({re}: value, text: string) array"),
		Str("findall"), Native3("findall", func(e *Env, rx, text, n Value) Value {
			m := rx.Table().GetString("_rx").Interface().(*regexp.Regexp).FindAllStringSubmatch(text.MustStr(""), int(n.IntDefault(-1)))
			mm := []Value{}
			for _, m := range m {
				for _, m := range m {
					mm = append(mm, Str(m))
				}
			}
			return Array(mm...)
		}, "$f({re}: value, text: string) array"),
		Str("replace"), Native3("replace", func(e *Env, rx, text, newtext Value) Value {
			m := rx.Table().GetString("_rx").Interface().(*regexp.Regexp).ReplaceAllString(text.MustStr(""), newtext.MustStr(""))
			return Str(m)
		}, "$f({re}: value, old: string, new: string) string"),
	))

	AddGlobalValue("error", func(env *Env, msg Value) Value {
		return Val(errors.New(msg.MustStr("")))
	}, "error(text: string) value", "\tcreate an error")
	AddGlobalValue("iserror", func(env *Env) {
		_, ok := env.Get(0).Interface().(error)
		env.A = Bool(ok)
	}, "iserror(v: value) bool", "\ttest whether value is an error")

	AddGlobalValue("json", Map(
		Str("stringify"), Native("stringify", func(env *Env) {
			env.A = Str(env.Get(0).JSONString())
		}, "$f(v: value) string"),
		Str("parse"), Native1("parse", func(env *Env, js Value) Value {
			j := strings.TrimSpace(js.MustStr(""))
			return Val(gjson.Parse(j))
		}, "$f(json: string) value"),
		Str("get"), Native3("get", func(env *Env, js, path, et Value) Value {
			j := strings.TrimSpace(js.MustStr("json string"))
			result := gjson.Get(j, path.MustStr("selector"))
			if !result.Exists() {
				return et
			}
			return Val(result)
		}, "$f(json: string, selector: string) value", "$f(json: string, selector: string, default: value) value"),
	))

	AddGlobalValue("sync", Map(
		Str("mutex"), Native("mutex", func(env *Env) { env.A = Val(&sync.Mutex{}) }, "$f() value", "\tcreate a mutex"),
		Str("rwmutex"), Native("rwmutex", func(env *Env) { env.A = Val(&sync.RWMutex{}) }, "$f() value", "\tcreate a read-write mutex"),
		Str("waitgroup"), Native("waitgroup", func(env *Env) { env.A = Val(&sync.WaitGroup{}) }, "$f() value", "\tcreate a wait group"),
	))
	AddGlobalValue("next", func(env *Env, m, k Value) Value {
		nk, nv := m.MustTable("").Next(k)
		return Array(nk, nv)
	}, "next(t: table, k: value) array", "\tfind next key-value pair after k in the table and return as { next_key, next_value }")
	AddGlobalValue("parent", func(env *Env, m Value) Value {
		return m.MustTable("").Parent.Value()
	}, "parent(t: table) table", "\tfind given table's parent, or nil if not existed")
	AddGlobalValue("unwrap", func(env *Env, m Value) Value {
		return ValRec(m.Interface())
	}, "unwrap(v: value) value", "\tunwrap Go's array, slice or map into table")
	AddGlobalValue("open", func(env *Env, path, flag, perm Value) Value {
		var opt int
		for _, f := range flag.StringDefault("") {
			switch f {
			case 'w':
				opt |= os.O_WRONLY | os.O_CREATE | os.O_TRUNC
			case 'r':
				opt |= os.O_RDONLY
			case 'a':
				opt |= os.O_APPEND | os.O_CREATE
			case 'x':
				opt |= os.O_EXCL
			case '+':
				opt |= os.O_RDWR | os.O_CREATE
			}
		}
		f, err := os.OpenFile(path.MustStr("path"), opt, fs.FileMode(perm.IntDefault(0)))
		if err != nil {
			panic(err)
		}
		a := Map(
			Str("_f"), Val(f),
			Str("close"), Native1("close", func(e *Env, rx Value) Value {
				return Val(rx.Table().GetString("_f").Interface().(*os.File).Close())
			}, ""),
			Str("read"), Native2("read", func(e *Env, rx, n Value) Value {
				f := rx.Table().GetString("_f").Interface().(*os.File)
				if c := n.IntDefault(0); c == 0 {
					buf, err := ioutil.ReadAll(f)
					if err != nil {
						panic(err)
					}
					return Bytes(buf)
				}
				p := make([]byte, n.IntDefault(0))
				rn, err := f.Read(p)
				if rn > 0 {
					return Bytes(p[:rn])
				}
				if err == io.EOF {
					return Nil
				}
				panic(err)
			}, ""),
			Str("write"), Native2("write", func(e *Env, rx, buf Value) Value {
				f := rx.Table().GetString("_f").Interface().(*os.File)
				wn, err := f.WriteString(buf.MustStr(""))
				if err != nil {
					panic(err)
				}
				return Int(int64(wn))
			}, ""),
			Str("seek"), Native3("seek", func(e *Env, rx, off, where Value) Value {
				f := rx.Table().GetString("_f").Interface().(*os.File)
				wn, err := f.Seek(off.MustInt("offset"), int(where.MustInt("where")))
				if err != nil {
					panic(err)
				}
				return Int(int64(wn))
			}, ""),
			Str("readlines"), Native2("readlines", func(e *Env, rx, cb Value) Value {
				f := rx.Table().GetString("_f").Interface().(*os.File)
				if _, err := f.Seek(0, 0); err != nil {
					panic(err)
				}
				if cb == Nil {
					buf, err := ioutil.ReadAll(f)
					if err != nil {
						panic(err)
					}
					res := []Value{}
					for _, line := range bytes.Split(buf, []byte("\n")) {
						res = append(res, Bytes(line))
					}
					return Array(res...)
				}
				for rd := bufio.NewReader(f); ; {
					line, err := rd.ReadString('\n')
					if len(line) > 0 {
						if v, err := cb.MustFunc("callback").Call(Str(line)); err != nil {
							panic(err)
						} else if v != Nil {
							return v
						}
					}
					if err != nil {
						if err != io.EOF {
							panic(err)
						}
						break
					}
				}
				return Nil
			},
				"readlines() array", "\tread the whole file and return lines as a table array",
				"readlines(f: function)", "\tfor every line read, f(line) will be called", "\tto exit the reading, return anything other than nil in callback",
			),
		)
		b := Map()
		b.Table().Parent = a.Table()
		return b
	}, "open(path: string, flag: string, perm: int) value")
}
