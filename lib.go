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
		r := NewObject(len(e.Global.Top.Locals))
		for i, name := range e.Global.Top.Locals {
			r.Set(Str(name), (*e.Global.Stack)[i])
		}
		e.A = r.ToValue()
	}, "$f() -> object", "\tlist all global symbols and their values")
	AddGlobalValue("doc", func(e *Env) {
		o := e.Object(0)
		_ = o.callable == nil && e.SetA(Nil) || e.SetA(Str(strings.Replace(o.callable.DocString, "$f", o.callable.Name, -1)))
	}, "$f(f: function) -> string", "\treturn `f`'s documentation")
	AddGlobalValue("new", func(e *Env) {
		m := e.Object(0)
		_ = e.Get(1).IsObject() && e.SetA(e.Object(1).SetParent(m).ToValue()) || e.SetA((&Object{parent: m}).ToValue())
	})
	AddGlobalValue("prototype", g["new"])
	AddGlobalValue("len", func(e *Env) { e.A = Int(e.Get(0).Len()) })
	AddGlobalValue("eval", func(e *Env) {
		var m map[string]interface{}
		if gt := e.Get(1).ToObject().Prop("globals"); gt.Type() == typ.Object {
			m = map[string]interface{}{}
			gt.Object().Foreach(func(k, v Value) bool {
				m[k.String()] = v.Interface()
				return true
			})
		}
		if e.Get(1).ToObject().Prop("compileonly").IsTrue() {
			v, err := parser.Parse(e.Str(0), "")
			internal.PanicErr(err)
			e.A = Val(v)
			return
		}
		wrap := func(err error) error { return fmt.Errorf("panic inside: %v", err) }
		p, err := LoadString(e.Str(0), &CompileOptions{GlobalKeyValues: m})
		if err != nil {
			panic(wrap(err))
		}
		v, err := p.Run()
		if err != nil {
			panic(wrap(err))
		}
		if !e.Get(1).ToObject().Prop("returnglobals").IsFalse() {
			r := NewObject(len(p.Top.Locals))
			for i, name := range p.Top.Locals {
				r.Set(Str(name), (*p.Stack)[i])
			}
			e.A = r.ToValue()
		} else {
			e.A = v
		}
	}, "$f(code: string, options?: object) -> value", "\tevaluate `code` and return the reuslt")
	AddGlobalValue("closure", func(e *Env) {
		lambda := e.Object(0)
		e.A = Func("<closure-"+lambda.Name()+">", func(e *Env) {
			f := e.Object(-1).Prop("_l").Object()
			stk := append([]Value{e.Object(-1).Prop("_c")}, e.Stack()...)
			e.A = f.MustCall(stk...)
		}).Object().Merge(nil,
			Str("_l"), e.Get(0),
			Str("_c"), e.Get(1),
		).ToValue()
	}, "$f(f: function, v: value) -> function",
		"\tcreate a function out of `f`, when it is called, `v` will be injected into as the first argument:",
		"\t\t closure(f, v)(args...) <=> f(v, args...)")

	// Debug libraries
	AddGlobalValue("debug", Obj(
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
			o := e.Object(0)
			_ = o.IsCallable() && e.SetA(Str(o.callable.ToCode())) || e.SetA(Nil)
		}),
	))
	AddGlobalValue("type", func(e *Env) {
		e.A = Str(e.Get(0).Type().String())
	}, "$f(v: value) -> string", "\treturn `v`'s type")
	AddGlobalValue("apply", func(e *Env) {
		e.A = e.Object(0).MustApply(e.Get(1), e.Stack()[2:]...)
	}, "$f(f: function, this: value, args...: value) -> value")
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
	AddGlobalValue("printf", func(e *Env) {
		sprintf(e, 0, e.Global.Stdout)
	}, "$f(format: string, args...: value)")
	AddGlobalValue("write", func(e *Env) {
		w := NewWriter(e.Get(0))
		for _, a := range e.Stack()[1:] {
			fmt.Fprint(w, a.String())
		}
	}, "$f(writer: value, args...: value)", "\twrite `args` to `writer`")
	AddGlobalValue("println", func(e *Env) {
		for _, a := range e.Stack() {
			fmt.Fprint(e.Global.Stdout, a.String(), " ")
		}
		fmt.Fprintln(e.Global.Stdout)
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
		time.Sleep(time.Duration(e.Float64(0) * float64(time.Second)))
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
	AddGlobalValue("exit", func(e *Env) { os.Exit(e.Int(0)) }, "$f(code: int)")
	AddGlobalValue("chr", func(e *Env) { e.A = Rune(rune(e.Int(0))) }, "$f(code: int) -> string")
	AddGlobalValue("byte", func(e *Env) { e.A = Byte(byte(e.Int(0))) }, "$f(code: int) -> string")
	AddGlobalValue("ord", func(e *Env) { r, _ := utf8.DecodeRuneInString(e.Str(0)); e.A = Int64(int64(r)) }, "$f(s: string) -> int")

	AddGlobalValue("re", Func("RegExp", func(e *Env) {
		e.A = Proto(e.A.Object(), Str("_rx"), Val(regexp.MustCompile(e.Str(0))))
	}, "re(regex: string) -> RegExp", "\tcreate a regular expression object").Object().Merge(nil,
		Str("match"), Func("", func(e *Env) {
			e.A = Bool(e.A.Object().Prop("_rx").Interface().(*regexp.Regexp).MatchString(e.Str(0)))
		}, "$f(text: string) -> bool"),
		Str("find"), Func("", func(e *Env) {
			m := e.A.Object().Prop("_rx").Interface().(*regexp.Regexp).FindStringSubmatch(e.Str(0))
			var mm []Value
			for _, m := range m {
				mm = append(mm, Str(m))
			}
			e.A = Array(mm...)
		}, "$f(text: string) -> array"),
		Str("findall"), Func("", func(e *Env) {
			m := e.A.Object().Prop("_rx").Interface().(*regexp.Regexp).FindAllStringSubmatch(e.Str(0), e.Get(1).ToInt(-1))
			var mm []Value
			for _, m := range m {
				for _, m := range m {
					mm = append(mm, Str(m))
				}
			}
			e.A = Array(mm...)
		}, "$f(text: string) -> array"),
		Str("replace"), Func("", func(e *Env) {
			e.A = Str(e.A.Object().Prop("_rx").Interface().(*regexp.Regexp).ReplaceAllString(e.Str(0), e.Str(1)))
		}, "$f(old: string, new: string) -> string"),
	))

	AddGlobalValue("error", func(e *Env) {
		e.A = Val(errors.New(e.Str(0)))
	}, "$f(text: string) -> go.error", "\tcreate an error")
	AddGlobalValue("iserror", func(e *Env) {
		_, ok := e.Interface(0).(error)
		e.A = Bool(ok)
	}, "$f(v: value) -> bool", "\treturn whether value is an error")

	AddGlobalValue("json", Obj(
		Str("stringify"), Func("", func(e *Env) {
			e.A = Str(e.Get(0).JSONString())
		}, "$f(v: value) -> string"),
		Str("parse"), Func("", func(e *Env) {
			e.A = Val(gjson.Parse(strings.TrimSpace(e.Str(0))))
		}, "$f(j: string) -> value"),
		Str("get"), Func("", func(e *Env) {
			result := gjson.Get(e.Str(0), e.Str(1))
			_ = !result.Exists() && e.SetA(e.Get(2)) || e.SetA(Val(result))
		}, "$f(j: string, selector: string, default?: value) -> value"),
	))

	AddGlobalValue("sync", Obj(
		Str("mutex"), Func("", func(e *Env) { e.A = Val(&sync.Mutex{}) }, "$f() -> *go.sync.Mutex"),
		Str("rwmutex"), Func("", func(e *Env) { e.A = Val(&sync.RWMutex{}) }, "$f() -> *go.sync.RWMutex"),
		Str("waitgroup"), Func("", func(e *Env) { e.A = Val(&sync.WaitGroup{}) }, "$f() -> *go.sync.WaitGroup"),
		Str("map"), Func("", func(e *Env) {
			fun := e.Object(1)
			n, t := e.Get(2).ToInt(runtime.NumCPU()), e.Get(0)
			if n < 1 || n > runtime.NumCPU()*1e3 {
				internal.Panic("invalid number of goroutines: %v", n)
			}
			var wg = sync.WaitGroup{}
			var in = make(chan [2]Value, t.Len())
			var outLock = sync.Mutex{}
			var outError error
			_ = t.Type() == typ.Array && e.SetA(Array(make([]Value, t.Len())...)) || e.SetA(NewObject(t.Len()).ToValue())
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
						if e.A.Type() == typ.Array {
							e.A.Array().Set(p[0].Int(), res)
						} else {
							e.A.Object().Set(p[0], res)
						}
						outLock.Unlock()
					}
				}()
			}
			t.ForEach(func(k, v Value) bool { in <- [2]Value{k, v}; return true })
			close(in)
			wg.Wait()
			internal.PanicErr(outError)
		}, "$f(a: object, f: function, n: int) -> object",
			"\tmap values in `a` into new values using f(k, v) concurrently on `n` goroutines (defaults to the number of CPUs)"),
	))

	AddGlobalValue("open", Func("", func(e *Env) {
		path, flag, perm := e.Str(0), e.Get(1).ToStr("r"), e.Get(2).ToInt64(0644)
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
		e.Object(-1).Set(Zero, Val(f))

		e.A = Func("FileObject", nil).Object().Merge(nil,
			Str("_f"), Val(f),
			Str("path"), Str(f.Name()),
			Str("sync"), Func("", func(e *Env) {
				internal.PanicErr(e.Object(-1).Prop("_f").Interface().(*os.File).Sync())
			}),
			Str("stat"), Func("", func(e *Env) {
				fi, err := e.Object(-1).Prop("_f").Interface().(*os.File).Stat()
				internal.PanicErr(err)
				e.A = Val(fi)
			}),
			Str("truncate"), Func("", func(e *Env) {
				f := e.Object(-1).Prop("_f").Interface().(*os.File)
				internal.PanicErr(f.Truncate(e.Int64(1)))
				t, err := f.Seek(0, 2)
				internal.PanicErr(err)
				e.A = Int64(t)
			}),
		).SetParent(ReadWriteSeekCloserProto).ToValue()
	}, "$f(path: string, flag: string, perm: int) -> FileObject").Object().Merge(nil,
		Str("close"), Func("", func(e *Env) {
			f, _ := e.Object(-1).Get(Zero).Interface().(*os.File)
			if f != nil {
				internal.PanicErr(f.Close())
			} else {
				internal.Panic("no opened file yet")
			}
		}, "$f()", "\tclose last opened file"),
	))
}
