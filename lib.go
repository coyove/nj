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
	AddGlobalValue("globals", func(e *Env) {
		r := NewTable(len(e.Global.Top.Locals))
		for i, name := range e.Global.Top.Locals {
			r.Set(Str(name), (*e.Global.Stack)[i])
		}
		e.A = r.Value()
	}, "$f() -> table", "\tlist all global symbols and their values")
	AddGlobalValue("doc", func(e *Env) {
		e.A = Str(e.Func(0).DocString())
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
	AddGlobalValue("len", func(e *Env) {
		switch v := e.Get(0); v.Type() {
		case typ.String:
			e.A = Int(len(v.Str()))
		case typ.Table:
			e.A = Int(v.Table().Len())
		case typ.Func:
			e.A = Int64(int64(v.Func().NumParams))
		case typ.Nil:
			e.A = Zero
		case typ.Number, typ.Bool:
			internal.Panic("can't measure length of %v", v.Type())
		default:
			e.A = Int(reflectLen(v.Interface()))
		}
	})
	AddGlobalValue("sizeof", func(e *Env) {
		e.A = Int64(int64(ValueSize))
		if v := e.Get(0); v.Type() == typ.Table {
			e.A = Int(v.Table().Size())
		}
	})
	AddGlobalValue("eval", func(s, g Value) Value {
		var m map[string]interface{}
		if gt := g.ToTableGets("globals"); gt.Type() == typ.Table {
			m = map[string]interface{}{}
			gt.Table().Foreach(func(k, v Value) bool {
				m[k.String()] = v.Interface()
				return true
			})
		}
		if !g.ToTableGets("compileonly").IsFalse() {
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
		if !g.ToTableGets("returnglobals").IsFalse() {
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
				f := e.Get(0).Recv("lambda").MustFunc("")
				src := e.Get(0).Recv("source")
				e.A = Str("<closure-" + f.Pure().String() + "-" + src.String() + ">")
			}),
			Str("__call"), Func("<closure-"+lambda.Name+">", func(e *Env) {
				f := e.Get(0).Recv("lambda").MustFunc("").Pure()
				stk := append([]Value{e.Get(0).Recv("source")}, e.Stack()[1:]...)
				e.A = MustValue(f.Call(stk...))
			}),
		)
	}, "$f(f: function, v: value) -> function",
		"\tcreate a function out of `f`, when it is called, `v` will be injected into as the first argument:",
		"\t\t closure(f, v)(args...) <=> f(v, args...)")

	// Debug libraries
	AddGlobalValue("debug", Map(
		Str("locals"), Func("", func(e *Env) {
			var r []Value
			start := e.stackOffset - uint32(e.CS.StackSize)
			for i, name := range e.CS.Locals {
				idx := start + uint32(i)
				r = append(r, Int64(int64(idx)), Str(name), (*e.stack)[idx])
			}
			e.A = Array(r...)
		}, "$f() -> array", "\treturn { index1, name1, value1, i2, n2, v2, i3, n3, v3, ... }"),
		Str("globals"), Func("", func(e *Env) {
			var r []Value
			for i, name := range e.Global.Top.Locals {
				r = append(r, Int(i), Str(name), (*e.Global.Stack)[i])
			}
			e.A = Array(r...)
		}, "$f() -> array", "\treturn { index1, name1, value1, i2, n2, v2, i3, n3, v3, ... }"),
		Str("set"), Func("set", func(e *Env) {
			(*e.Global.Stack)[e.Int64(0)] = e.Get(1)
		}, "$f(idx: int, v: value)"),
		Str("trace"), Func("", func(env *Env) {
			stacks := append(env.Stacktrace, stacktrace{cls: env.CS, cursor: env.IP})
			lines := make([]Value, 0, len(stacks))
			for i := len(stacks) - 1 - env.Get(0).ToInt(0); i >= 0; i-- {
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
		Str("disfunc"), Func("", func(e *Env) {
			e.A = Str(e.Func(0).PrettyCode())
		}),
	))
	AddGlobalValue("type", func(e *Env) {
		e.A = Str(e.Get(0).Type().String())
	}, "$f(v: value) -> string", "\treturn value's type")
	AddGlobalValue("apply", func(e *Env) {
		fun := *e.Func(0)
		fun.Receiver = e.Get(1)
		e.A = MustValue(fun.Call(e.Stack()[2:]...))
	}, "$f(f: function, receiver: value, args...: value) -> value")
	AddGlobalValue("pcall", func(e *Env) {
		if a, err := e.Func(0).Call(e.Stack()[1:]...); err == nil {
			e.A = a
		} else {
			e.A = wrapExecError(err)
		}
	}, "$f(f: function, args...: value) -> value", "\texecute `f`, catch panic and return as error if any")
	AddGlobalValue("gcall", Map(
		Str("__name"), Str("gcall"),
		Str("__call"), Func("", func(e *Env) {
			f := e.Func(1).Copy()
			args := e.CopyStack()[2:]
			w := make(chan Value, 1)
			go func(f *Function, args []Value) {
				if v, err := f.Call(args...); err != nil {
					w <- wrapExecError(err)
				} else {
					w <- v
				}
			}(f, args)
			e.A = TableProto(e.Table(0), Str("f"), f.Value(), Str("w"), intf(w))
		}, "$f(f: function, args...: value) -> table^gcall", "\texecute `f` in goroutine"),
		Str("stop"), Func("", func(e *Env) {
			e.Recv("f").MustFunc("").EmergStop()
		}),
		Str("wait"), Func("", func(e *Env) {
			ch := e.Recv("w").Interface().(chan Value)
			if w := e.Get(1).ToFloat64(0); w > 0 {
				select {
				case <-time.After(time.Duration(w * float64(time.Second))):
					panic("timeout")
				case v := <-ch:
					e.A = v
				}
			} else {
				e.A = <-ch
			}
		}),
	))
	AddGlobalValue("panic", func(e *Env) { panic(e.Get(0)) }, "$f(v: value)")
	AddGlobalValue("assert", func(e *Env) {
		if v := e.Get(0); e.Size() <= 1 && v.IsFalse() {
			internal.Panic("assertion failed")
		} else if e.Size() == 2 && !v.Equal(e.Get(1)) {
			internal.Panic("assertion failed: %v and %v", v, e.Get(1))
		} else if e.Size() == 3 && !v.Equal(e.Get(1)) {
			internal.Panic("%s: %v and %v", e.Get(2).String(), v, e.Get(1))
		}
	}, "$f(v: value)", "\tpanic when value is falsy",
		"$f(v1: value, v2: value, msg?: string)", "\tpanic when two values are not equal")
	AddGlobalValue("int", func(e *Env) {
		if v := e.Get(0); v.Type() == typ.Number {
			e.A = Int64(v.Int64())
		} else {
			v, err := strconv.ParseInt(v.String(), e.Get(1).ToInt(0), 64)
			internal.PanicErr(err)
			e.A = Int64(v)
		}
	}, "$f(v: value, base?: int) -> int", "\tconvert `v` to an integer number, panic when failed or overflowed")
	AddGlobalValue("float", func(e *Env) {
		if v := e.Get(0); v.Type() == typ.Number {
			e.A = v
		} else if v := parser.Num(v.String()); v.Type() == parser.FLOAT {
			e.A = Float64(v.Float())
		} else {
			e.A = Int64(v.Int())
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
	AddGlobalValue("write", func(e *Env) {
		w := NewWriter(e.Get(0))
		for _, a := range e.Stack()[1:] {
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
		fmt.Fprint(env.Global.Stdout, prompt.ToStr(""))
		var results []Value
		var r io.Reader = env.Global.Stdin
		for i := n.ToInt64(1); i > 0; i-- {
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
		time.Sleep(time.Duration(e.Get(0).MustFloat64("") * float64(time.Second)))
	}, "$f(sec: float)")
	AddGlobalValue("Go_time", func(e *Env) {
		if e.Size() > 0 {
			e.A = Val(time.Date(
				int(e.Get(0).ToInt64(1970)), time.Month(e.Get(1).ToInt64(1)), int(e.Get(2).ToInt64(1)),
				int(e.Get(3).ToInt64(0)), int(e.Get(4).ToInt64(0)), int(e.Get(5).ToInt64(0)), int(e.Get(6).ToInt64(0)),
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
	AddGlobalValue("exit", func(e *Env) { os.Exit(int(e.Get(0).MustInt64(""))) }, "$f(code: int)")
	AddGlobalValue("chr", func(e *Env) { e.A = Rune(rune(e.Get(0).MustInt64(""))) }, "$f(code: int) -> string")
	AddGlobalValue("byte", func(e *Env) { e.A = Byte(byte(e.Get(0).MustInt64(""))) }, "$f(code: int) -> string")
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
			m := re.Recv("_rx").Interface().(*regexp.Regexp).FindAllStringSubmatch(text.MustStr(""), int(n.ToInt64(-1)))
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
		_, ok := e.Get(0).Interface().(error)
		e.A = Bool(ok)
	}, "$f(v: value) -> bool", "\treturn whether value is an error")

	AddGlobalValue("json", Map(
		Str("stringify"), Func("", func(e *Env) {
			e.A = Str(e.Get(0).JSONString())
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
		Str("map"), Func("", func(e *Env) {
			fun := e.Func(1)
			n, t := e.Get(2).ToInt(runtime.NumCPU()), e.Table(0)
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
			e.A = out.Value()
		}, "$f(t: table, f: function, n: int) -> table",
			"\tmap values in `t` into new values using f(k, v) concurrently on `n` goroutines (defaults to the number of CPUs)"),
	))
	AddGlobalValue("next", func(e *Env) {
		e.A = Array(e.Table(0).Next(e.Get(1)))
	}, "next(t: table, k: value) -> array", "\tfind next key-value pair after `k` in the table and return as { next_key, next_value }")
	AddGlobalValue("parent", func(e *Env) {
		e.A = e.Table(0).Parent().Value()
	}, "parent(t: table) -> table", "\tfind given table's parent, or nil if not existed")

	AddGlobalValue("open", Map(
		Str("__call"), Func("", func(e *Env) {
			path, flag, perm := e.Str(1), e.Get(2).ToStr("r"), e.Get(3).ToInt64(0644)
			var opt int
			for _, f := range flag {
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
			f, err := os.OpenFile(path, opt, fs.FileMode(perm))
			internal.PanicErr(err)
			e.Table(0).Set(Zero, Val(f))

			m := TableProto(ReadWriteSeekCloserProto,
				Str("_f"), Val(f),
				Str("__name"), Str(f.Name()),
				Str("sync"), Func("", func(e *Env) {
					internal.PanicErr(e.Recv("_f").Interface().(*os.File).Sync())
				}),
				Str("stat"), Func("", func(e *Env) {
					fi, err := e.Recv("_f").Interface().(*os.File).Stat()
					internal.PanicErr(err)
					e.A = Val(fi)
				}),
				Str("truncate"), Func("", func(e *Env) {
					f := e.Recv("_f").Interface().(*os.File)
					internal.PanicErr(f.Truncate(e.Int64(1)))
					t, err := f.Seek(0, 2)
					internal.PanicErr(err)
					e.A = Int64(t)
				}),
			)
			e.A = m
		}, "open(path: string, flag: string, perm: int) -> table^"+ReadWriteSeekCloserProto.Name()),
		Str("close"), Func("", func(e *Env) {
			f, _ := e.Table(0).Get(Zero).Interface().(*os.File)
			if f != nil {
				internal.PanicErr(f.Close())
			} else {
				internal.Panic("no opened file yet")
			}
		}, "$f()", "\tclose last opened file"),
	))
}
