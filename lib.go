package nj

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
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

var (
	globals = map[string]Value{}

	StrLib    Value
	MathLib   Value
	ObjectLib Value
	FuncLib   Value
	ArrayLib  Value
	OSLib     Value
	IOLib     Value
	ErrorLib  Value
)

func AddGlobalValue(k string, v interface{}, doc ...string) {
	switch v := v.(type) {
	case func(*Env):
		globals[k] = Func(k, v, doc...)
	case Value:
		globals[k] = renameFuncName(Str(k), v)
	default:
		globals[k] = ValueOf(v)
	}
}

func RemoveGlobalValue(k string) {
	delete(globals, k)
}

func init() {
	AddGlobalValue("VERSION", Int64(Version))
	AddGlobalValue("globals", func(e *Env) {
		e.A = e.Global.LocalsObject().ToValue()
	}, "$f() -> object", "\tlist all global symbols and their values")
	AddGlobalValue("new", func(e *Env) {
		m := e.Object(0)
		_ = e.Get(1).IsObject() && e.SetA(e.Object(1).SetProto(m).ToValue()) || e.SetA(NewObject(0).SetProto(m).ToValue())
	})
	AddGlobalValue("prototype", globals["new"])
	AddGlobalValue("len", func(e *Env) { e.A = Int(e.Get(0).Len()) })
	AddGlobalValue("loadfile", func(e *Env) {
		e.A = MustRun(LoadFile(e.Str(0), e.Global.Options))
	}, "$f(path: string) -> value", "\tload and eval file at `path`, globals will be inherited in loaded file")
	AddGlobalValue("eval", func(e *Env) {
		opts := e.Get(1).ToObject()
		if opts.Prop("ast").IsTrue() {
			v, err := parser.Parse(e.Str(0), "")
			internal.PanicErr(err)
			e.A = ValueOf(v)
			return
		}
		p, err := LoadString(e.Str(0), &CompileOptions{Globals: opts.Prop("globals").ToObject()})
		internal.PanicErr(err)
		v, err := p.Run()
		internal.PanicErr(err)
		_ = opts.Prop("returnglobals").IsTrue() && e.SetA(p.LocalsObject().ToValue()) || e.SetA(v)
	}, "$f(code: string, options?: object) -> value", "\tevaluate `code` and return the reuslt")
	AddGlobalValue("closure", func(e *Env) {
		lambda := e.Object(0)
		e.A = Func("<closure-"+lambda.Name()+">", func(e *Env) {
			f := e.Object(-1).Prop("_l").Object()
			stk := append([]Value{e.Object(-1).Prop("_c")}, e.Stack()...)
			e.A = e.Call(f, stk...)
		}).Object().Merge(nil, Str("_l"), e.Get(0), Str("_c"), e.Get(1)).ToValue()
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
		}, "$f() -> array", "\treturn [index1, name1, value1, i2, n2, v2, i3, n3, v3, ...]"),
		Str("globals"), Func("", func(e *Env) {
			var r []Value
			for i, name := range e.Global.Top.Locals {
				r = append(r, Int(i), Str(name), (*e.Global.Stack)[i])
			}
			e.A = Array(r...)
		}, "$f() -> array", "\treturn [index1, name1, value1, i2, n2, v2, i3, n3, v3, ...]"),
		Str("set"), Func("set", func(e *Env) {
			(*e.Global.Stack)[e.Int64(0)] = e.Get(1)
		}, "$f(idx: int, v: value)"),
		Str("trace"), Func("", func(env *Env) {
			stacks := env.Runtime.GetFullStacktrace()
			lines := make([]Value, 0, len(stacks))
			for i := len(stacks) - 1 - env.Get(0).ToInt(0); i >= 0; i-- {
				r := stacks[i]
				src := uint32(0)
				for i := 0; i < len(r.Callable.Code.Pos); {
					var opx uint32 = math.MaxUint32
					ii, op, line := r.Callable.Code.Pos.read(i)
					if ii < len(r.Callable.Code.Pos)-1 {
						_, opx, _ = r.Callable.Code.Pos.read(ii)
					}
					if r.Cursor >= op && r.Cursor < opx {
						src = line
						break
					}
					if r.Cursor < op && i == 0 {
						src = line
						break
					}
					i = ii
				}
				lines = append(lines, Str(r.Callable.Name), Int64(int64(src)), Int64(int64(r.Cursor-1)))
			}
			env.A = Array(lines...)
		}, "$f(skip?: int) -> array", "\treturn [func_name0, line1, cursor1, n2, l2, c2, ...]"),
		Str("disfunc"), Func("", func(e *Env) {
			o := e.Object(0)
			_ = o.IsCallable() && e.SetA(Str(o.Callable.ToCode())) || e.SetA(Nil)
		}),
	))
	AddGlobalValue("type", func(e *Env) {
		e.A = Str(e.Get(0).Type().String())
	}, "$f(v: value) -> string", "\treturn `v`'s type")
	AddGlobalValue("apply", func(e *Env) {
		e.A = CallObject(e.Object(0), e, nil, e.Get(1), e.Stack()[2:]...)
	}, "$f(f: function, this: value, args...: value) -> value")
	AddGlobalValue("panic", func(e *Env) { panic(e.Get(0)) }, "$f(v: value)")
	AddGlobalValue("throw", func(e *Env) { panic(e.Get(0)) }, "$f(v: value)")
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
			e.A = Float64(v.Float64())
		} else {
			e.A = Int64(v.Int64())
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
			_, err := fmt.Fprint(w, a.String())
			e.A = Error(e, err)
		}
	}, "$f(writer: Writer, args...: value)", "\twrite `args` to `writer`")
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
	}, "$f() -> array", "\tread all user inputs and return as [input1, input2, ...]",
		"$f(prompt: string, n?: int) -> array", "\tprint `prompt` then read all (or at most `n`) user inputs")
	AddGlobalValue("time", func(e *Env) {
		e.A = Float64(float64(time.Now().UnixNano()) / 1e9)
	}, "$f() -> float", "\tunix timestamp in seconds")
	AddGlobalValue("sleep", func(e *Env) { time.Sleep(e.Num(0).ToDuration(0)) }, "$f(sec: float)")
	AddGlobalValue("Go_time", func(e *Env) {
		if e.Size() > 0 {
			e.A = ValueOf(time.Date(e.Int(0), time.Month(e.Int(1)), e.Int(2),
				e.Get(3).ToInt(0), e.Get(4).ToInt(0), e.Get(5).ToInt(0), e.Get(6).ToInt(0), time.UTC))
		} else {
			e.A = ValueOf(time.Now())
		}
	},
		"$f() -> go.time.Time",
		"\treturn time.Time of current time",
		"$f(year: int, month: int, day: int, h?: int, m?: int, s?: int, nanoseconds?: int) -> go.time.Time",
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
		e.A = Proto(e.A.Object(), Str("_rx"), ValueOf(regexp.MustCompile(e.Str(0))))
	}, "re(regex: string) -> RegExp", "\tcreate a regular expression object").Object().Merge(nil,
		Str("match"), Func("", func(e *Env) {
			e.A = Bool(e.Object(-1).Prop("_rx").Interface().(*regexp.Regexp).MatchString(e.Str(0)))
		}, "RegExp.$f(text: string) -> bool"),
		Str("find"), Func("", func(e *Env) {
			m := e.Object(-1).Prop("_rx").Interface().(*regexp.Regexp).FindStringSubmatch(e.Str(0))
			e.A = NewSequence(m, stringsSequenceMeta).ToValue()
		}, "RegExp.$f(text: string) -> array"),
		Str("findall"), Func("", func(e *Env) {
			m := e.Object(-1).Prop("_rx").Interface().(*regexp.Regexp).FindAllStringSubmatch(e.Str(0), e.Get(1).ToInt(-1))
			e.A = NewSequence(m, GetGenericSequenceMeta(m)).ToValue()
		}, "RegExp.$f(text: string) -> array"),
		Str("replace"), Func("", func(e *Env) {
			e.A = Str(e.Object(-1).Prop("_rx").Interface().(*regexp.Regexp).ReplaceAllString(e.Str(0), e.Str(1)))
		}, "RegExp.$f(old: string, new: string) -> string"),
	))

	AddGlobalValue("json", Obj(
		Str("stringify"), Func("", func(e *Env) {
			e.A = Str(e.Get(0).JSONString())
		}, "$f(v: value) -> string"),
		Str("parse"), Func("", func(e *Env) {
			e.A = ValueOf(gjson.Parse(strings.TrimSpace(e.Str(0))))
		}, "$f(j: string) -> value"),
		Str("get"), Func("", func(e *Env) {
			result := gjson.Get(e.Str(0), e.Str(1))
			_ = !result.Exists() && e.SetA(e.Get(2)) || e.SetA(ValueOf(result))
		}, "$f(j: string, path: string, default?: value) -> value"),
	))

	AddGlobalValue("sync", Obj(
		Str("mutex"), Func("", func(e *Env) { e.A = ValueOf(&sync.Mutex{}) }, "$f() -> *go.sync.Mutex"),
		Str("rwmutex"), Func("", func(e *Env) { e.A = ValueOf(&sync.RWMutex{}) }, "$f() -> *go.sync.RWMutex"),
		Str("waitgroup"), Func("", func(e *Env) { e.A = ValueOf(&sync.WaitGroup{}) }, "$f() -> *go.sync.WaitGroup"),
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
						res, err := e.Call2(fun, p[0], p[1])
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
			if e.A.Type() == typ.Array {
				t.Array().Foreach(func(k int, v Value) bool { in <- [2]Value{Int(k), v}; return true })
			} else {
				t.Object().Foreach(func(k, v Value) bool { in <- [2]Value{k, v}; return true })
			}
			close(in)
			wg.Wait()
			internal.PanicErr(outError)
		}, "$f(a: object|array, f: function, n: int) -> object|array",
			"\tmap values in `a` into new values using `f(k, v)` concurrently on `n` goroutines (defaults to the number of CPUs)"),
	))

	AddGlobalValue("open", Func("open", func(e *Env) {
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
		e.Object(-1).Set(Zero, ValueOf(f))

		e.A = Func("File", nil).Object().Merge(nil,
			Str("_f"), ValueOf(f),
			Str("path"), Str(f.Name()),
			Str("sync"), Func("", func(e *Env) {
				internal.PanicErr(e.Object(-1).Prop("_f").Interface().(*os.File).Sync())
			}, "File.$f()"),
			Str("stat"), Func("", func(e *Env) {
				fi, err := e.Object(-1).Prop("_f").Interface().(*os.File).Stat()
				internal.PanicErr(err)
				e.A = ValueOf(fi)
			}, "File.$f() -> go.os.FileInfo"),
			Str("truncate"), Func("", func(e *Env) {
				f := e.Object(-1).Prop("_f").Interface().(*os.File)
				internal.PanicErr(f.Truncate(e.Int64(1)))
				t, err := f.Seek(0, 2)
				internal.PanicErr(err)
				e.A = Int64(t)
			}, "File.$f() -> int"),
		).SetProto(ReadWriteSeekCloserProto).ToValue()
	}, "$f(path: string, flag: string, perm: int) -> File").Object().Merge(nil,
		Str("close"), Func("", func(e *Env) {
			f, _ := e.Object(-1).Get(Zero).Interface().(*os.File)
			if f != nil {
				internal.PanicErr(f.Close())
			} else {
				internal.Panic("no opened file yet")
			}
		}, "$f()", "\tclose last opened file"),
	))

	IOLib = Obj(
		Str("reader"), ReaderProto.ToValue(),
		Str("writer"), WriterProto.ToValue(),
		Str("seeker"), SeekerProto.ToValue(),
		Str("closer"), CloserProto.ToValue(),
		Str("readwriter"), ReadWriterProto.ToValue(),
		Str("readcloser"), ReadCloserProto.ToValue(),
		Str("writecloser"), WriteCloserProto.ToValue(),
		Str("readwritecloser"), ReadWriteCloserProto.ToValue(),
		Str("readwriteseekcloser"), ReadWriteSeekCloserProto.ToValue(),
	)
	AddGlobalValue("io", IOLib)

	ObjectLib = Func("object", func(e *Env) {
		switch e.Get(0).Type() {
		case typ.Number:
			e.A = NewObject(e.Int(0)).ToValue()
		case typ.Nil:
			e.A = NewObject(0).ToValue()
		default:
			e.A = e.Object(0).SetFirstProto(e.Object(-1)).ToValue()
		}
	}).Object().Merge(nil,
		Str("get"), Func("", func(e *Env) {
			e.A = e.Object(-1).Get(e.Get(0))
		}, "object.$f(k: value) -> value"),
		Str("set"), Func("", func(e *Env) {
			e.A = e.Object(-1).Set(e.Get(0), e.Get(1))
		}, "object.$f(k: value, v: value) -> value", "\tset value and return previous value"),
		Str("rawget"), Func("", func(e *Env) {
			e.A = e.Object(-1).RawGet(e.Get(0))
		}, "object.$f(k: value) -> value"),
		Str("delete"), Func("", func(e *Env) {
			e.A = e.Object(-1).Delete(e.Get(0))
		}, "object.$f(k: value) -> value", "\tdelete value and return previous value"),
		Str("clear"), Func("", func(e *Env) { e.Object(-1).Clear() }, "object.$f()"),
		Str("copy"), Func("", func(e *Env) {
			e.A = e.Object(-1).Copy().ToValue()
		}, "object.$f() -> object", "\tcopy the object"),
		Str("proto"), Func("", func(e *Env) {
			e.A = e.Object(-1).Proto().ToValue()
		}, "object.$f() -> object", "\treturn object's prototype if any"),
		Str("setproto"), Func("", func(e *Env) {
			e.Object(-1).SetProto(e.Object(0))
		}, "object.$f(p: object)", "\tset object's prototype to `p`"),
		Str("setfirstproto"), Func("", func(e *Env) {
			e.Object(-1).SetFirstProto(e.Object(0))
		}, "object.$f(p: object)", "\tinsert `p` as `t`'s first prototype"),
		Str("size"), Func("", func(e *Env) {
			e.A = Int(e.Object(-1).Size())
		}, "object.$f() -> int", "\treturn the underlay size of object, which always >= object.len()"),
		Str("len"), Func("", func(e *Env) {
			e.A = Int(e.Object(-1).Len())
		}, "object.$f() -> int", "\treturn the count of keys in object"),
		Str("keys"), Func("", func(e *Env) {
			a := make([]Value, 0)
			e.Object(-1).Foreach(func(k, v Value) bool { a = append(a, k); return true })
			e.A = Array(a...)
		}, "object.$f() -> array"),
		Str("values"), Func("", func(e *Env) {
			a := make([]Value, 0)
			e.Object(-1).Foreach(func(k, v Value) bool { a = append(a, v); return true })
			e.A = Array(a...)
		}, "object.$f() -> array"),
		Str("items"), Func("", func(e *Env) {
			a := make([]Value, 0)
			e.Object(-1).Foreach(func(k, v Value) bool { a = append(a, Array(k, v)); return true })
			e.A = Array(a...)
		}, "object.$f() -> [[value, value]]", "return as [[key1, value1], [key2, value2], ...]"),
		Str("foreach"), Func("", func(e *Env) {
			f := e.Object(0)
			e.Object(-1).Foreach(func(k, v Value) bool { return e.Call(f, k, v) == Nil })
		}, "object.$f(f: function)"),
		Str("contains"), Func("", func(e *Env) {
			found, b := false, e.Get(0)
			e.Object(-1).Foreach(func(k, v Value) bool { found = v.Equal(b); return !found })
			e.A = Bool(found)
		}, "object.$f(v: value) -> bool"),
		Str("merge"), Func("", func(e *Env) {
			e.A = e.Object(-1).Merge(e.Object(0)).ToValue()
		}, "object.$f(o: object)", "\tmerge elements from `o` to this"),
		Str("tostring"), Func("", func(e *Env) {
			p := &bytes.Buffer{}
			e.Object(-1).rawPrint(p, 0, typ.MarshalToJSON, true)
			e.A = UnsafeStr(p.Bytes())
		}, "object.$f() -> string", "\tprint raw elements in object"),
		Str("pure"), Func("", func(e *Env) {
			m2 := *e.Object(-1)
			m2.parent = nil
			e.A = m2.ToValue()
		}, "object.$f() -> object", "\treturn a new object who shares all the data from the original, but with no prototype"),
		Str("next"), Func("", func(e *Env) {
			e.A = Array(e.Object(-1).Next(e.Get(0)))
		}, "object.$f(k: value) -> [value, value]", "\tfind next key-value pair after `k` in the object and return as [key, value]"),
		Str("iscallable"), Func("", func(e *Env) {
			e.A = Bool(e.Object(-1).IsCallable())
		}, "object.$f() -> bool"),
	).ToValue()
	FuncLib = Func("function", nil).Object().Merge(nil,
		Str("doc"), Func("", func(e *Env) {
			o := e.Object(-1)
			_ = o.Callable == nil && e.SetA(Nil) || e.SetA(Str(strings.Replace(o.Callable.DocString, "$f", o.Callable.Name, -1)))
		}, "function.$f() -> string", "\treturn function documentation"),
		Str("apply"), Func("", func(e *Env) {
			e.A = CallObject(e.Object(-1), e, nil, e.Get(0), e.Stack()[1:]...)
		}, "function.$f(this: value, args...: value) -> value"),
		Str("call"), Func("", func(e *Env) {
			e.A = e.Call(e.Object(-1), e.Stack()...)
		}, "function.$f(args...: value) -> value"),
		Str("try"), Func("", func(e *Env) {
			a, err := e.Call2(e.Object(-1), e.Stack()...)
			_ = err == nil && e.SetA(a) || e.SetA(Error(e, err))
		}, "function.$f(args...: value) -> value|Error",
			"\trun function, return Error if any panic occurred (if function tends to return n results, these values will all be Error now)"),
		Str("after"), Func("", func(e *Env) {
			f, args, e2 := e.Object(-1), e.CopyStack()[1:], *e
			e2.Stacktrace = append([]Stacktrace{}, e2.Stacktrace...)
			e.A = Func("Timer", nil).Object().Merge(nil,
				Str("t"), ValueOf(time.AfterFunc(e.Num(0).ToDuration(0), func() { e2.Call(f, args...) })),
				Str("stop"), Func("", func(e *Env) {
					e.A = Bool(e.Object(-1).Prop("t").Interface().(*time.Timer).Stop())
				}),
			).ToValue()
		}, "function.$f(secs: float, args...: value) -> Timer"),
		Str("go"), Func("", func(e *Env) {
			f := e.Object(-1)
			args := e.CopyStack()
			w := make(chan Value, 1)
			e2 := *e
			e2.Stacktrace = append([]Stacktrace{}, e2.Stacktrace...)
			go func(f *Object, args []Value) {
				if v, err := e2.Call2(f, args...); err != nil {
					w <- Error(&e2, err)
				} else {
					w <- v
				}
			}(f, args)
			e.A = Func("Goroutine", nil).Object().Merge(nil,
				Str("f"), f.ToValue(),
				Str("w"), intf(w),
				Str("stop"), Func("", func(e *Env) {
					e.Object(-1).Prop("f").Object().Callable.EmergStop()
				}),
				Str("wait"), Func("", func(e *Env) {
					ch := e.Object(-1).Prop("w").Interface().(chan Value)
					select {
					case <-time.After(e.Get(0).ToDuration(1 << 62)):
						panic("timeout")
					case v := <-ch:
						e.A = v
					}
				}),
			).ToValue()
		}, "function.$f(args...: value) -> Goroutine", "\texecute `f` in goroutine"),
	).SetProto(ObjectLib.Object()).ToValue()

	AddGlobalValue("object", ObjectLib)
	AddGlobalValue("func", FuncLib)

	ArrayLib = Func("array", nil).Object().Merge(nil,
		Str("make"), Func("", func(e *Env) {
			e.A = Array(make([]Value, e.Int(0))...)
		}, "array.$f(n: int) -> array", "\tcreate an array of size `n`"),
		Str("len"), Func("", func(e *Env) { e.A = Int(e.Array(-1).Len()) }, "array.$f()"),
		Str("size"), Func("", func(e *Env) { e.A = Int(e.Array(-1).Size()) }, "array.$f()"),
		Str("clear"), Func("", func(e *Env) { e.Array(-1).Clear() }, "array.$f()"),
		Str("append"), Func("", func(e *Env) {
			e.Array(-1).Append(e.Stack()...)
		}, "array.$f(args...: value)", "\tappend values to array"),
		Str("find"), Func("", func(e *Env) {
			a, src, ff := -1, e.Array(-1), e.Get(0)
			for i := 0; i < src.Len(); i++ {
				if src.Get(i).Equal(ff) {
					a = i
					break
				}
			}
			e.A = Int(a)
		}, "array.$f(v: value) -> int", "\tfind the index of first `v` in array"),
		Str("filter"), Func("", func(e *Env) {
			src, ff := e.Array(-1), e.Object(0)
			dest := make([]Value, 0, src.Len())
			src.Foreach(func(k int, v Value) bool {
				if e.Call(ff, v).IsTrue() {
					dest = append(dest, v)
				}
				return true
			})
			e.A = Array(dest...)
		}, "array.$f(f: function) -> array", "\tfilter out all values where f(value) is false"),
		Str("slice"), Func("", func(e *Env) {
			a := e.Array(-1)
			st, en := e.Int(0), e.Get(1).ToInt(a.Len())
			for ; st < 0 && a.Len() > 0; st += a.Len() {
			}
			for ; en < 0 && a.Len() > 0; en += a.Len() {
			}
			e.A = a.Slice(st, en).ToValue()
		}, "array.$f(start: int, end?: int) -> array", "\tslice array, from `start` to `end`"),
		Str("removeat"), Func("", func(e *Env) {
			ma, idx := e.Array(-1), e.Int(0)
			if idx < 0 || idx >= ma.Len() {
				e.A = Nil
			} else {
				old := ma.Get(idx)
				ma.Copy(idx, ma.Len(), ma.Slice(idx+1, ma.Len()))
				ma.SliceInplace(0, ma.Len()-1)
				e.A = old
			}
		}, "array.$f(index: int)", "\tremove value at `index`"),
		Str("copy"), Func("", func(e *Env) {
			e.Array(-1).Copy(e.Int(0), e.Int(1), e.Array(2))
		}, "array.$f(start: int, end: int, src: array) -> array", "\tcopy elements from `src` to `this[start:end]`"),
		Str("concat"), Func("", func(e *Env) {
			e.Array(-1).Concat(e.Array(0))
		}, "array.$f(a: array) -> array", "\tconcat two arrays"),
		Str("istyped"), Func("", func(e *Env) {
			e.A = Bool(e.Array(-1).meta != internalSequenceMeta)
		}, "array.$f() -> bool"),
		Str("untype"), Func("", func(e *Env) {
			e.A = Array(e.Array(-1).Values()...)
		}, "array.$f() -> array"),
	).ToValue()
	AddGlobalValue("array", ArrayLib)

	ErrorLib = Func("Error", func(e *Env) {
		e.A = Error(nil, &ExecError{root: e.Get(0), stacks: e.GetFullStacktrace()})
	}).Object().Merge(nil,
		Str("error"), Func("", func(e *Env) { e.A = ValueOf(e.Array(-1).Unwrap().(*ExecError).root) }),
		Str("getcause"), Func("", func(e *Env) { e.A = intf(e.Array(-1).Unwrap().(*ExecError).root) }),
	).SetProto(ArrayLib.Object()).ToValue()
	AddGlobalValue("error", ErrorLib)

	encDecProto := Func("EncodeDecode", nil).Object().Merge(nil,
		Str("encode"), Func("", func(e *Env) {
			i := e.Object(-1).Prop("_e").Interface()
			e.A = Str(i.(interface{ EncodeToString([]byte) string }).EncodeToString(e.Get(0).ToBytes()))
		}),
		Str("decode"), Func("", func(e *Env) {
			i := e.Object(-1).Prop("_e").Interface()
			v, err := i.(interface{ DecodeString(string) ([]byte, error) }).DecodeString(e.Str(0))
			internal.PanicErr(err)
			e.A = Bytes(v)
		}),
	).SetProto(Func("EncoderDecoder", nil).Object().Merge(nil,
		Str("encoder"), Func("", func(e *Env) {
			enc := Nil
			buf := &bytes.Buffer{}
			switch encoding := e.Object(-1).Prop("_e").Interface().(type) {
			default:
				enc = ValueOf(hex.NewEncoder(buf))
			case *base32.Encoding:
				enc = ValueOf(base32.NewEncoder(encoding, buf))
			case *base64.Encoding:
				enc = ValueOf(base64.NewEncoder(encoding, buf))
			}
			e.A = Func("Encoder", nil).Object().Merge(nil,
				Str("_f"), ValueOf(enc),
				Str("_b"), ValueOf(buf),
				Str("value"), Func("", func(e *Env) {
					e.A = Str(e.Object(-1).Prop("_b").Interface().(*bytes.Buffer).String())
				}),
				Str("bytes"), Func("", func(e *Env) {
					e.A = Bytes(e.Object(-1).Prop("_b").Interface().(*bytes.Buffer).Bytes())
				}),
			).SetProto(WriteCloserProto).ToValue()
		}),
		Str("decoder"), Func("", func(e *Env) {
			src := NewReader(e.Get(0))
			dec := Nil
			switch encoding := e.Object(-1).Prop("_e").Interface().(type) {
			case *base64.Encoding:
				dec = ValueOf(base64.NewDecoder(encoding, src))
			case *base32.Encoding:
				dec = ValueOf(base32.NewDecoder(encoding, src))
			default:
				dec = ValueOf(hex.NewDecoder(src))
			}
			e.A = Func("Decoder", nil).Object().Merge(nil, Str("_f"), ValueOf(dec)).SetProto(ReaderProto).ToValue()
		}),
	))

	StrLib = Func("str", func(e *Env) {
		i, ok := e.Interface(0).([]byte)
		_ = ok && e.SetA(UnsafeStr(i)) || e.SetA(Str(e.Get(0).String()))
	}).Object().Merge(nil,
		Str("from"), Func("", func(e *Env) {
			e.A = Str(fmt.Sprint(e.Interface(0)))
		}, "$f(v: value) -> string", "\tconvert `v` to string"),
		Str("size"), Func("", func(e *Env) {
			e.A = Int(e.StrLen(-1))
		}, "str.$f() -> int", "\tsame as len()"),
		Str("len"), Func("", func(e *Env) {
			e.A = Int(e.StrLen(-1))
		}, "str.$f() -> int", "\tsame as size()"),
		Str("count"), Func("", func(e *Env) {
			e.A = Int(utf8.RuneCountInString(e.Str(-1)))
		}, "str.$f() -> int", "\treturn the count of runes"),
		Str("iequals"), Func("", func(e *Env) {
			e.A = Bool(strings.EqualFold(e.Str(-1), e.Str(0)))
		}, "str.$f(text2: string) -> bool", "\tcase insensitive equal"),
		Str("contains"), Func("", func(e *Env) {
			e.A = Bool(strings.Contains(e.Str(-1), e.Str(0)))
		}, "str.$f(substr: string) -> bool"),
		Str("split"), Func("", func(e *Env) {
			s, d := e.Str(-1), e.Str(0)
			if n := e.Get(1).ToInt(0); n == 0 {
				e.A = NewSequence(strings.Split(s, d), stringsSequenceMeta).ToValue()
			} else {
				e.A = NewSequence(strings.SplitN(s, d, n), stringsSequenceMeta).ToValue()
			}
		}, "str.$f(delim: string, n?: int) -> array"),
		Str("join"), Func("", func(e *Env) {
			d := e.Str(-1)
			buf := &bytes.Buffer{}
			e.Array(0).Foreach(func(k int, v Value) bool {
				buf.WriteString(v.String())
				buf.WriteString(d)
				return true
			})
			if buf.Len() > 0 {
				buf.Truncate(buf.Len() - len(d))
			}
			e.A = UnsafeStr(buf.Bytes())
		}, "str.$f(a: array) -> string"),
		Str("replace"), Func("", func(e *Env) {
			e.A = Str(strings.Replace(e.Str(-1), e.Str(0), e.Str(1), e.Get(2).ToInt(-1)))
		}, "str.$f(old: string, new: string, count?: int) -> string"),
		Str("match"), Func("", func(e *Env) {
			m, err := filepath.Match(e.Str(-1), e.Str(0))
			internal.PanicErr(err)
			e.A = Bool(m)
		}, "str.$f(text: string) -> bool"),
		Str("find"), Func("", func(e *Env) {
			e.A = Int(strings.Index(e.Str(-1), e.Str(0)))
		}, "str.$f(sub: string) -> int", "\tfind the index of first appearence of `sub` in text"),
		Str("findlast"), Func("", func(e *Env) {
			e.A = Int(strings.LastIndex(e.Str(-1), e.Str(0)))
		}, "str.$f(sub: string) -> int", "\tsame as find(), but starting from right to left"),
		Str("sub"), Func("", func(e *Env) {
			s := e.Str(-1)
			st, en := e.Int(0), e.Get(1).ToInt(len(s))
			for ; st < 0 && len(s) > 0; st += len(s) {
			}
			for ; en < 0 && len(s) > 0; en += len(s) {
			}
			e.A = Str(s[st:en])
		}, "str.$f(start: int, end?: int) -> string"),
		Str("trim"), Func("", func(e *Env) {
			_ = e.Get(0).IsNil() && e.SetA(Str(strings.TrimSpace(e.Str(-1)))) || e.SetA(Str(strings.Trim(e.Str(-1), e.Str(0))))
		}, "str.$f(cutset?: string) -> string", "\ttrim spaces (or any chars in `cutset`) at both sides of the text"),
		Str("trimprefix"), Func("", func(e *Env) {
			e.A = Str(strings.TrimPrefix(e.Str(-1), e.Str(0)))
		}, "str.$f(prefix: string) -> string", "\ttrim `prefix` of the text"),
		Str("trimsuffix"), Func("", func(e *Env) {
			e.A = Str(strings.TrimSuffix(e.Str(-1), e.Str(0)))
		}, "str.$f(suffix: string) -> string", "\ttrim `suffix` of the text"),
		Str("trimleft"), Func("", func(e *Env) {
			e.A = Str(strings.TrimLeft(e.Str(-1), e.Str(0)))
		}, "str.$f(cutset: string) -> string", "\ttrim the left side of the text using every char in `cutset`"),
		Str("trimright"), Func("", func(e *Env) {
			e.A = Str(strings.TrimRight(e.Str(-1), e.Str(0)))
		}, "str.$f(cutset: string) -> string", "\ttrim the right side of the text using every char in `cutset`"),
		Str("ord"), Func("", func(e *Env) {
			r, sz := utf8.DecodeRuneInString(e.Str(-1))
			e.A = Array(Int64(int64(r)), Int(sz))
		}, "str.$f() -> [int, int]", "\tdecode first UTF-8 char, return [unicode, bytescount]"),
		Str("startswith"), Func("", func(e *Env) { e.A = Bool(strings.HasPrefix(e.Str(-1), e.Str(0))) }, "str.$f(prefix: string) -> bool"),
		Str("endswith"), Func("", func(e *Env) { e.A = Bool(strings.HasSuffix(e.Str(-1), e.Str(0))) }, "str.$f(suffix: string) -> bool"),
		Str("upper"), Func("", func(e *Env) { e.A = Str(strings.ToUpper(e.Str(-1))) }, "str.$f() -> string"),
		Str("lower"), Func("", func(e *Env) { e.A = Str(strings.ToLower(e.Str(-1))) }, "str.$f() -> string"),
		Str("bytes"), Func("", func(e *Env) {
			_ = e.Get(0).IsInt64() && e.SetA(ValueOf(make([]byte, e.Int(0)))) || e.SetA(ValueOf([]byte(e.Str(0))))
		}, "str.$f() -> bytes", "\tconvert text to byte array",
			"$f(n: int) -> bytes", "\tcreate an n-byte long array"),
		Str("chars"), Func("", func(e *Env) {
			var r []Value
			for s := e.Str(-1); len(s) > 0; {
				_, sz := utf8.DecodeRuneInString(s)
				if sz == 0 {
					break
				}
				r = append(r, Str(s[:sz]))
				s = s[sz:]
			}
			e.A = Array(r...)
		}, "str.$f() -> array", "\tbreak `text` into chars, e.g.: chars('a中c') => ['a', '中', 'c']"),
		Str("format"), Func("", func(e *Env) {
			buf := &bytes.Buffer{}
			sprintf(e, -1, buf)
			e.A = UnsafeStr(buf.Bytes())
		}, "str.$f(args...: value) -> string"),
		Str("buffer"), Func("", func(e *Env) {
			b := &bytes.Buffer{}
			if v := e.Get(0); v != Nil {
				b.WriteString(v.String())
			}
			e.A = Func("Buffer", nil).Object().SetProto(ReadWriterProto).Merge(nil,
				Str("_f"), ValueOf(b),
				Str("reset"), Func("", func(e *Env) {
					e.Object(-1).Prop("_f").Interface().(*bytes.Buffer).Reset()
				}, "Buffer.$f()"),
				Str("value"), Func("", func(e *Env) {
					e.A = UnsafeStr(e.Object(-1).Prop("_f").Interface().(*bytes.Buffer).Bytes())
				}, "Buffer.$f() -> string"),
				Str("bytes"), Func("", func(e *Env) {
					e.A = Bytes(e.Object(-1).Prop("_f").Interface().(*bytes.Buffer).Bytes())
				}, "Buffer.$f() -> bytes"),
			).ToValue()
		}, "$f(v?: string) -> Buffer"),
		Str("hex"), Func("hex", nil).Object().SetProto(encDecProto.Proto()).ToValue(),
		Str("base64"), Func("base64", nil).Object().Merge(nil,
			Str("std"), Proto(encDecProto, Str("_e"), ValueOf(getEncB64(base64.StdEncoding, '='))),
			Str("url"), Proto(encDecProto, Str("_e"), ValueOf(getEncB64(base64.URLEncoding, '='))),
			Str("std2"), Proto(encDecProto, Str("_e"), ValueOf(getEncB64(base64.StdEncoding, -1))),
			Str("url2"), Proto(encDecProto, Str("_e"), ValueOf(getEncB64(base64.URLEncoding, -1))),
		).SetProto(encDecProto).ToValue(),
		Str("base32"), Func("base32", nil).Object().Merge(nil,
			Str("std"), Proto(encDecProto, Str("_e"), ValueOf(getEncB32(base32.StdEncoding, '='))),
			Str("hex"), Proto(encDecProto, Str("_e"), ValueOf(getEncB32(base32.HexEncoding, '='))),
			Str("std2"), Proto(encDecProto, Str("_e"), ValueOf(getEncB32(base32.StdEncoding, -1))),
			Str("hex2"), Proto(encDecProto, Str("_e"), ValueOf(getEncB32(base32.HexEncoding, -1))),
		).SetProto(encDecProto).ToValue(),
	).ToValue()
	AddGlobalValue("str", StrLib)

	MathLib = Obj(
		Str("INF"), Float64(math.Inf(1)),
		Str("NEG_INF"), Float64(math.Inf(-1)),
		Str("PI"), Float64(math.Pi),
		Str("E"), Float64(math.E),
		Str("randomseed"), Func("", func(e *Env) { rand.Seed(e.B(0).ToInt64(1)) }, "$f(seed: int)"),
		Str("random"), Func("", func(e *Env) {
			switch len(e.Stack()) {
			case 2:
				ai, bi := e.Int64(0), e.Int64(1)
				e.A = Int64(rand.Int63n(bi-ai+1) + ai)
			case 1:
				e.A = Int64(rand.Int63n(e.Int64(0)))
			default:
				e.A = Float64(rand.Float64())
			}
		}, "$f() -> float", "\treturn [0, 1)", "$f(n: int) -> int", "\treturn [0, n)", "$f(a: int, b: int) -> int", "\treturn [a, b]"),
		Str("sqrt"), Func("", func(e *Env) { e.A = Float64(math.Sqrt(e.Float64(0))) }, "$f(v: float) -> float"),
		Str("floor"), Func("", func(e *Env) { e.A = Float64(math.Floor(e.Float64(0))) }, "$f(v: float) -> float"),
		Str("ceil"), Func("", func(e *Env) { e.A = Float64(math.Ceil(e.Float64(0))) }, "$f(v: float) -> float"),
		Str("min"), Func("", func(e *Env) { mathMinMax(e, false) }, "$f(a: number, b...: number) -> number"),
		Str("max"), Func("", func(e *Env) { mathMinMax(e, true) }, "$f(a: number, b...: number) -> number"),
		Str("pow"), Func("", func(e *Env) { e.A = Float64(math.Pow(e.Float64(0), e.Float64(1))) }, "$f(a: float, b: float) -> float"),
		Str("abs"), Func("", func(e *Env) {
			if e.A = e.Num(0); e.A.IsInt64() {
				if i := e.A.Int64(); i < 0 {
					e.A = Int64(-i)
				}
			} else {
				e.A = Float64(math.Abs(e.A.Float64()))
			}
		}),
		Str("remainder"), Func("", func(e *Env) { e.A = Float64(math.Remainder(e.Float64(0), e.Float64(1))) }),
		Str("mod"), Func("", func(e *Env) { e.A = Float64(math.Mod(e.Float64(0), e.Float64(1))) }),
		Str("cos"), Func("", func(e *Env) { e.A = Float64(math.Cos(e.Float64(0))) }),
		Str("sin"), Func("", func(e *Env) { e.A = Float64(math.Sin(e.Float64(0))) }),
		Str("tan"), Func("", func(e *Env) { e.A = Float64(math.Tan(e.Float64(0))) }),
		Str("acos"), Func("", func(e *Env) { e.A = Float64(math.Acos(e.Float64(0))) }),
		Str("asin"), Func("", func(e *Env) { e.A = Float64(math.Asin(e.Float64(0))) }),
		Str("atan"), Func("", func(e *Env) { e.A = Float64(math.Atan(e.Float64(0))) }),
		Str("atan2"), Func("", func(e *Env) { e.A = Float64(math.Atan2(e.Float64(0), e.Float64(1))) }),
		Str("ldexp"), Func("", func(e *Env) { e.A = Float64(math.Ldexp(e.Float64(0), e.Int(0))) }),
		Str("modf"), Func("", func(e *Env) { a, b := math.Modf(e.Float64(0)); e.A = Array(Float64(a), Float64(b)) }),
	)
	AddGlobalValue("math", MathLib)

	OSLib = Obj(
		Str("args"), ValueOf(os.Args),
		Str("environ"), Func("", func(e *Env) { e.A = ValueOf(os.Environ()) }),
		Str("shell"), Func("", func(e *Env) {
			p := exec.Command("sh", "-c", e.Str(0))
			opt := e.Get(1)
			opt.ToObject().Prop("env").ToObject().Foreach(func(k, v Value) bool {
				p.Env = append(p.Env, k.String()+"="+v.String())
				return true
			})
			stdout := &bytes.Buffer{}
			p.Stdout, p.Stderr = stdout, stdout
			p.Dir = opt.ToObject().Prop("dir").ToStr("")
			if tmp := opt.ToObject().Prop("stdout"); tmp != Nil {
				p.Stdout = NewWriter(tmp)
			}
			if tmp := opt.ToObject().Prop("stderr"); tmp != Nil {
				p.Stderr = NewWriter(tmp)
			}
			if tmp := opt.ToObject().Prop("stdin"); tmp != Nil {
				p.Stdin = NewReader(tmp)
			}

			out := make(chan error)
			go func() { out <- p.Run() }()
			select {
			case r := <-out:
				internal.PanicErr(r)
			case <-time.After(opt.ToObject().Prop("timeout").ToDuration(time.Duration(1 << 62))):
				p.Process.Kill()
				panic("timeout")
			}
			e.A = Bytes(stdout.Bytes())
		}),
		Str("readdir"), Func("", func(e *Env) {
			fi, err := ioutil.ReadDir(e.Str(0))
			internal.PanicErr(err)
			e.A = ValueOf(fi)
		}),
		Str("remove"), Func("", func(e *Env) {
			path := e.Str(0)
			fi, err := os.Stat(path)
			internal.PanicErr(err)
			if fi.IsDir() {
				internal.PanicErr(os.RemoveAll(path))
			} else {
				internal.PanicErr(os.Remove(path))
			}
		}),
		Str("pstat"), Func("", func(e *Env) {
			fi, err := os.Stat(e.Str(0))
			_ = err == nil && e.SetA(ValueOf(fi)) || e.SetA(Nil)
		}),
	)
	AddGlobalValue("os", OSLib)

}
