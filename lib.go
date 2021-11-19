package nj

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
	"unsafe"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/parser"
	"github.com/coyove/nj/typ"
	"github.com/tidwall/gjson"
)

const Version int64 = 327

var g = map[string]Value{}

func AddGlobalValue(k string, v interface{}, doc ...string) {
	switch v := v.(type) {
	case func(*Env):
		g[k] = Func(k, v, doc...)
	case func(Value) Value:
		g[k] = Func1(k, v, doc...)
	case func(Value, Value) Value:
		g[k] = Func2(k, v, doc...)
	case func(Value, Value, Value) Value:
		g[k] = Func3(k, v, doc...)
	default:
		g[k] = Val(v)
	}
}

func RemoveGlobalValue(k string) {
	delete(g, k)
}

func init() {
	AddGlobalValue("VERSION", Int64(Version))
	AddGlobalValue("globals", func(env *Env) {
		r := NewTable(len(env.Global.Top.Locals))
		for i, name := range env.Global.Top.Locals {
			r.Set(Str(name), (*env.Global.Stack)[i])
		}
		env.A = r.Value()
	}, "$f() -> table", "\tlist all global symbols and their values")
	AddGlobalValue("doc", func(f Value) Value {
		return Str(f.MustFunc("").DocString())
	}, "$f(f: function) -> string", "\treturn `f`'s documentation")
	AddGlobalValue("new", func(v, a Value) Value {
		m := v.MustTable("").New()
		if a.Type() != typ.Table {
			return (&Table{parent: m}).Value()
		}
		a.Table().SetParent(m)
		return a
	})
	AddGlobalValue("prototype", g["new"])
	AddGlobalValue("len", func(v Value) Value {
		switch v.Type() {
		case typ.String:
			return Int(len(v.Str()))
		case typ.Table:
			return Int(v.Table().Len())
		case typ.Func:
			return Int64(int64(v.Func().NumParams))
		case typ.Nil:
			return Zero
		case typ.Number, typ.Bool:
			internal.Panic("can't measure length of %v", v.Type())
		}
		return Int(reflectLen(v.Interface()))
	})
	AddGlobalValue("sizeof", func(v Value) Value {
		if v.Type() == typ.Table {
			return Int(v.Table().Size())
		}
		return Int64(int64(ValueSize))
	})
	AddGlobalValue("eval", func(s, g Value) Value {
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
			internal.PanicErr(err)
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
		if !g.MaybeTableGetString("returnglobals").IsFalse() {
			r := NewTable(len(p.Top.Locals))
			for i, name := range p.Top.Locals {
				r.Set(Str(name), (*p.Stack)[i])
			}
			return r.Value()
		}
		return v
	}, "$f(code: string, options?: table) -> value", "\tevaluate `code` and return the reuslt")
	AddGlobalValue("closure", func(f, m Value) Value {
		lambda := f.MustFunc("")
		return Map(
			Str("source"), m,
			Str("lambda"), lambda.Value(),
			Str("__str"), Func("<closure-"+lambda.Name+"__str>", func(e *Env) {
				f := e.B(0).Recv("lambda").MustFunc("")
				src := e.B(0).Recv("source")
				e.A = Str("<closure-" + f.Pure().String() + "-" + src.String() + ">")
			}),
			Str("__call"), Func("<closure-"+lambda.Name+">", func(e *Env) {
				f := e.B(0).Recv("lambda").MustFunc("").Pure()
				stk := append([]Value{e.B(0).Recv("source")}, e.Stack()[1:]...)
				e.A = MustValue(f.Call(stk...))
			}),
		)
	}, "$f(f: function, v: value) -> function",
		"\tcreate a function out of `f`, when it is called, `v` will be injected into as the first argument:",
		"\t\t closure(f, v)(args...) <=> f(v, args...)")

	// Debug libraries
	AddGlobalValue("debug", Map(
		Str("locals"), Func("", func(env *Env) {
			var r []Value
			start := env.stackOffset - uint32(env.CS.StackSize)
			for i, name := range env.CS.Locals {
				idx := start + uint32(i)
				r = append(r, Int64(int64(idx)), Str(name), (*env.stack)[idx])
			}
			env.A = Array(r...)
		}, "$f() -> array", "\treturn { index1, name1, value1, i2, n2, v2, i3, n3, v3, ... }"),
		Str("globals"), Func("", func(env *Env) {
			var r []Value
			for i, name := range env.Global.Top.Locals {
				r = append(r, Int(i), Str(name), (*env.Global.Stack)[i])
			}
			env.A = Array(r...)
		}, "$f() -> array", "\treturn { index1, name1, value1, i2, n2, v2, i3, n3, v3, ... }"),
		Str("set"), Func("set", func(env *Env) {
			(*env.Global.Stack)[env.B(0).MustInt64("")] = env.Get(1)
		}, "$f(idx: int, v: value)"),
		Str("trace"), Func("", func(env *Env) {
			stacks := append(env.Stacktrace, stacktrace{cls: env.CS, cursor: env.IP})
			lines := make([]Value, 0, len(stacks))
			for i := len(stacks) - 1 - env.B(0).MaybeInt(0); i >= 0; i-- {
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
				lines = append(lines, Str(r.cls.Name), Int64(int64(src)), Int64(int64(r.cursor-1)))
			}
			env.A = Array(lines...)
		}, "$f(skip: int) -> array", "\treturn { func_name0, line1, cursor1, n2, l2, c2, ... }"),
		Str("disfunc"), Func1("", func(v Value) Value {
			return Str(v.MustFunc("").PrettyCode())
		}),
	))
	AddGlobalValue("type", func(env *Env) {
		env.A = Str(env.B(0).Type().String())
	}, "$f(v: value) -> string", "\treturn value's type")
	AddGlobalValue("apply", func(e *Env) {
		fun := *e.B(0).MustFunc("")
		fun.Receiver = e.B(1)
		e.A = MustValue(fun.Call(e.Stack()[2:]...))
	}, "$f(f: function, receiver: value, args...: value) -> value")
	AddGlobalValue("pcall", func(env *Env) {
		if a, err := env.B(0).MustFunc("").Call(env.Stack()[1:]...); err == nil {
			env.A = a
		} else {
			env.A = wrapExecError(err)
		}
	}, "$f(f: function, args...: value) -> value", "\texecute `f`, catch panic and return as error if any")
	AddGlobalValue("gcall", Map(
		Str("__name"), Str("gcall"),
		Str("__call"), Func("", func(env *Env) {
			f := env.Get(1).MustFunc("").Copy()
			args := env.CopyStack()[2:]
			w := make(chan Value, 1)
			go func(f *Function, args []Value) {
				if v, err := f.Call(args...); err != nil {
					w <- wrapExecError(err)
				} else {
					w <- v
				}
			}(f, args)
			env.A = TableProto(env.B(0).MustTable(""), Str("f"), f.Value(), Str("w"), intf(w))
		}, "$f(f: function, args...: value) -> table^gcall", "\texecute `f` in goroutine"),
		Str("stop"), Func("", func(env *Env) {
			env.B(0).Recv("f").MustFunc("").EmergStop()
		}),
		Str("wait"), Func2("", func(m, t Value) Value {
			ch := m.Recv("w").Interface().(chan Value)
			if w := t.MaybeFloat(0); w > 0 {
				select {
				case <-time.After(time.Duration(w * float64(time.Second))):
					panic("timeout")
				case v := <-ch:
					return v
				}
			}
			return <-ch
		}),
	))
	AddGlobalValue("panic", func(e *Env) { panic(e.B(0)) }, "$f(v: value)")
	AddGlobalValue("assert", func(e *Env) {
		if v := e.B(0); e.Size() <= 1 && v.IsFalse() {
			internal.Panic("assertion failed")
		} else if e.Size() == 2 && !v.Equal(e.B(1)) {
			internal.Panic("assertion failed: %v and %v", v, e.B(1))
		} else if e.Size() == 3 && !v.Equal(e.B(1)) {
			internal.Panic("%s: %v and %v", e.B(2).String(), v, e.B(1))
		}
	}, "$f(v: value)", "\tpanic when value is falsy",
		"$f(v1: value, v2: value, msg?: string)", "\tpanic when two values are not equal")
	AddGlobalValue("int", func(env *Env) {
		if v := env.B(0); v.Type() == typ.Number {
			env.A = Int64(v.Int64())
		} else {
			v, err := strconv.ParseInt(v.String(), env.Get(1).MaybeInt(0), 64)
			internal.PanicErr(err)
			env.A = Int64(v)
		}
	}, "$f(v: value, base?: int) -> int", "\tconvert `v` to an integer number, panic when failed or overflowed")
	AddGlobalValue("float", func(env *Env) {
		if v := env.B(0); v.Type() == typ.Number {
			env.A = v
		} else if v := parser.Num(v.String()); v.Type() == parser.FLOAT {
			env.A = Float64(v.Float())
		} else {
			env.A = Int64(v.Int())
		}
	}, "$f(v: value) -> number", "\tconvert `v` to a float number, panic when failed")
	AddGlobalValue("print", func(env *Env) {
		for _, a := range env.Stack() {
			fmt.Fprint(env.Global.Stdout, a.String())
		}
		fmt.Fprintln(env.Global.Stdout)
	}, "$f(args...: value)", "\tprint `args` to stdout with no space between them")
	AddGlobalValue("printf", func(env *Env) {
		sprintf(env, env.Global.Stdout)
	}, "$f(format: string, args...: value)")
	AddGlobalValue("write", func(env *Env) {
		w := NewWriter(env.B(0))
		for _, a := range env.Stack()[1:] {
			fmt.Fprint(w, a.String())
		}
	}, "$f(writer: value, args...: value)", "\twrite `args` to `writer`")
	AddGlobalValue("println", func(env *Env) {
		for _, a := range env.Stack() {
			fmt.Fprint(env.Global.Stdout, a.String(), " ")
		}
		fmt.Fprintln(env.Global.Stdout)
	}, "$f(args...: value)", "\tprint values, insert space between each of them")
	AddGlobalValue("scanln", func(env *Env) {
		prompt, n := env.B(0), env.Get(1)
		fmt.Fprint(env.Global.Stdout, prompt.MaybeStr(""))
		var results []Value
		var r io.Reader = env.Global.Stdin
		for i := n.MaybeInt64(1); i > 0; i-- {
			var s string
			if _, err := fmt.Fscan(r, &s); err != nil {
				break
			}
			results = append(results, Str(s))
		}
		env.A = Array(results...)
	}, "$f() -> array", "\tread all user inputs and return as { input1, input2, ... }",
		"$f(prompt: string, n?: int) -> array", "\tprint `prompt` then read all (or at most `n`) user inputs")
	AddGlobalValue("time", func(e *Env) {
		e.A = Float64(float64(time.Now().UnixNano()) / 1e9)
	}, "$f() -> float", "\tunix timestamp in seconds")
	AddGlobalValue("sleep", func(e *Env) {
		time.Sleep(time.Duration(e.B(0).MustFloat("") * float64(time.Second)))
	}, "$f(sec: float)")
	AddGlobalValue("Go_time", func(e *Env) {
		if e.Size() > 0 {
			e.A = Val(time.Date(
				int(e.B(0).MaybeInt64(1970)), time.Month(e.B(1).MaybeInt64(1)), int(e.B(2).MaybeInt64(1)),
				int(e.B(3).MaybeInt64(0)), int(e.B(4).MaybeInt64(0)), int(e.B(5).MaybeInt64(0)), int(e.B(6).MaybeInt64(0)),
				time.UTC))
		} else {
			e.A = Val(time.Now())
		}
	},
		"$f() -> go.time.Time",
		"\treturn time.Time of current time",
		"$f(year: int, month: int, day: int, h: int, m: int, s: int, nanoseconds: int) -> go.time.Time",
		"\treturn time.Time constructed by the given arguments",
	)
	AddGlobalValue("clock", func(e *Env) {
		x := time.Now()
		s := *(*[2]int64)(unsafe.Pointer(&x))
		e.A = Float64(float64(s[1]) / 1e9)
	}, "$f() -> float", "\tseconds since startup (monotonic clock)")
	AddGlobalValue("exit", func(e *Env) { os.Exit(int(e.B(0).MustInt64(""))) }, "$f(code: int)")
	AddGlobalValue("chr", func(e *Env) { e.A = Rune(rune(e.B(0).MustInt64(""))) }, "$f(code: int) -> string")
	AddGlobalValue("byte", func(e *Env) { e.A = Byte(byte(e.B(0).MustInt64(""))) }, "$f(code: int) -> string")
	AddGlobalValue("ord", func(env *Env) {
		r, _ := utf8.DecodeRuneInString(env.B(0).MustStr(""))
		env.A = Int64(int64(r))
	}, "$f(s: string) -> int")

	AddGlobalValue("re", Map(
		Str("__name"), Str("relib"),
		Str("__call"), Func2("", func(re, r Value) Value {
			return TableProto(re.MustTable(""), Str("_rx"), Val(regexp.MustCompile(r.MustStr(""))))
		}, "$f(regex: string) -> table^relib", "\tcreate a regular expression object"),
		Str("match"), Func2("", func(re, text Value) Value {
			return Bool(re.Recv("_rx").Interface().(*regexp.Regexp).MatchString(text.MustStr("")))
		}, "$f({re}: value, text: string) -> bool"),
		Str("find"), Func2("", func(re, text Value) Value {
			m := re.Recv("_rx").Interface().(*regexp.Regexp).FindStringSubmatch(text.MustStr(""))
			mm := []Value{}
			for _, m := range m {
				mm = append(mm, Str(m))
			}
			return Array(mm...)
		}, "$f({re}: value, text: string) -> array"),
		Str("findall"), Func3("", func(re, text, n Value) Value {
			m := re.Recv("_rx").Interface().(*regexp.Regexp).FindAllStringSubmatch(text.MustStr(""), int(n.MaybeInt64(-1)))
			var mm []Value
			for _, m := range m {
				for _, m := range m {
					mm = append(mm, Str(m))
				}
			}
			return Array(mm...)
		}, "$f({re}: value, text: string) -> array"),
		Str("replace"), Func3("", func(re, text, newtext Value) Value {
			m := re.Recv("_rx").Interface().(*regexp.Regexp).ReplaceAllString(text.MustStr(""), newtext.MustStr(""))
			return Str(m)
		}, "$f({re}: value, old: string, new: string) -> string"),
	))

	AddGlobalValue("error", func(msg Value) Value {
		return Val(errors.New(msg.MustStr("")))
	}, "$f(text: string) -> go.error", "\tcreate an error")
	AddGlobalValue("iserror", func(e *Env) {
		_, ok := e.B(0).Interface().(error)
		e.A = Bool(ok)
	}, "$f(v: value) -> bool", "\treturn whether value is an error")

	AddGlobalValue("json", Map(
		Str("stringify"), Func("", func(e *Env) {
			e.A = Str(e.B(0).JSONString())
		}, "$f(v: value) -> string"),
		Str("parse"), Func1("", func(js Value) Value {
			return Val(gjson.Parse(strings.TrimSpace(js.MustStr(""))))
		}, "$f(j: string) -> value"),
		Str("get"), Func3("", func(js, path, et Value) Value {
			result := gjson.Get(js.MustStr("json string"), path.MustStr("selector"))
			if !result.Exists() {
				return et
			}
			return Val(result)
		}, "$f(j: string, selector: string, default?: value) -> value"),
	))

	AddGlobalValue("sync", Map(
		Str("mutex"), Func("", func(e *Env) { e.A = Val(&sync.Mutex{}) }, "$f() -> *go.sync.Mutex"),
		Str("rwmutex"), Func("", func(e *Env) { e.A = Val(&sync.RWMutex{}) }, "$f() -> *go.sync.RWMutex"),
		Str("waitgroup"), Func("", func(e *Env) { e.A = Val(&sync.WaitGroup{}) }, "$f() -> *go.sync.WaitGroup"),
		Str("map"), Func3("", func(list, f, opt Value) Value {
			fun := f.MustFunc("mapping")
			n, t := opt.MaybeInt(runtime.NumCPU()), list.MustTable("")
			if n < 1 || n > runtime.NumCPU()*1e3 {
				internal.Panic("invalid number of goroutines: %v", n)
			}
			var wg = sync.WaitGroup{}
			var in = make(chan [2]Value, t.Len())
			var out, outLock = t.Copy(), sync.Mutex{}
			var outError error
			wg.Add(n)
			for i := 0; i < n; i++ {
				go func() {
					defer wg.Done()
					for p := range in {
						if outError != nil {
							return
						}
						res, err := fun.Call(p[0], p[1])
						if err != nil {
							outError = err
							return
						}
						outLock.Lock()
						out.RawSet(p[0], res)
						outLock.Unlock()
					}
				}()
			}
			t.Foreach(func(k, v Value) bool { in <- [2]Value{k, v}; return true })
			close(in)
			wg.Wait()
			internal.PanicErr(outError)
			return out.Value()
		}, "$f(t: table, f: function, n: int) -> table",
			"\tmap values in `t` into new values using f(k, v) concurrently on `n` goroutines (defaults to the number of CPUs)"),
	))
	AddGlobalValue("next", func(m, k Value) Value {
		return Array(m.MustTable("").Next(k))
	}, "next(t: table, k: value) -> array", "\tfind next key-value pair after `k` in the table and return as { next_key, next_value }")
	AddGlobalValue("parent", func(m Value) Value {
		return m.MustTable("").Parent().Value()
	}, "parent(t: table) -> table", "\tfind given table's parent, or nil if not existed")

	AddGlobalValue("open", Map(
		Str("__call"), Func("", func(env *Env) {
			path, flag, perm := env.Get(1).MustStr("path"), env.Get(2), env.Get(3)
			var opt int
			for _, f := range flag.MaybeStr("r") {
				switch f {
				case 'w':
					opt &^= os.O_RDWR | os.O_RDONLY
					opt |= os.O_WRONLY | os.O_CREATE | os.O_TRUNC
				case 'r':
					opt &^= os.O_RDWR | os.O_WRONLY
					opt |= os.O_RDONLY
				case 'a':
					opt |= os.O_APPEND | os.O_CREATE
				case 'x':
					opt |= os.O_EXCL
				case '+':
					opt &^= os.O_RDONLY | os.O_WRONLY
					opt |= os.O_RDWR | os.O_CREATE
				}
			}
			f, err := os.OpenFile(path, opt, fs.FileMode(perm.MaybeInt64(0644)))
			internal.PanicErr(err)
			env.B(0).MustTable("").Set(Int64(0), Val(f))

			m := TableProto(ReadWriteSeekCloserProto,
				Str("_f"), Val(f),
				Str("__name"), Str(f.Name()),
				Str("sync"), Func1("", func(rx Value) Value {
					internal.PanicErr(rx.Recv("_f").Interface().(*os.File).Sync())
					return Nil
				}),
				Str("stat"), Func1("", func(rx Value) Value {
					fi, err := rx.Recv("_f").Interface().(*os.File).Stat()
					internal.PanicErr(err)
					return Val(fi)
				}),
				Str("truncate"), Func2("", func(rx, n Value) Value {
					f := rx.Recv("_f").Interface().(*os.File)
					internal.PanicErr(f.Truncate(n.MustInt64("")))
					t, err := f.Seek(0, 2)
					internal.PanicErr(err)
					return Int64(t)
				}),
			)
			env.A = m
		}, "open(path: string, flag: string, perm: int) -> table^"+ReadWriteSeekCloserProto.Name()),
		Str("close"), Func("", func(env *Env) {
			f, _ := env.B(0).MustTable("").Get(Int64(0)).Interface().(*os.File)
			if f != nil {
				internal.PanicErr(f.Close())
			} else {
				internal.Panic("no opened file yet")
			}
		}, "$f()", "\tclose last opened file"),
	))
}
