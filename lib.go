package nj

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"math"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
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

	"github.com/coyove/nj/bas"
	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/parser"
	"github.com/coyove/nj/typ"
)

func init() {
	bas.AddTopFunc("chr", func(e *bas.Env) { e.A = bas.Rune(rune(e.Int(0))) })
	bas.AddTopFunc("byte", func(e *bas.Env) { e.A = bas.Byte(byte(e.Int(0))) })
	bas.AddTopFunc("ord", func(e *bas.Env) {
		r, _ := utf8.DecodeRuneInString(e.Str(0))
		e.A = bas.Int64(int64(r))
	})

	bas.AddTopValue("json", bas.NewNamedObject("json", 0).
		SetProp("stringify", bas.Func("stringify", func(e *bas.Env) { e.A = bas.Str(e.Get(0).JSONString()) })).
		SetProp("dump", bas.Func("dump", func(e *bas.Env) { e.Get(1).Stringify(e.Get(0).Writer(), typ.MarshalToJSON) })).
		SetProp("parse", bas.Func("parse", func(e *bas.Env) {
			s := strings.TrimSpace(e.Str(0))
			if s == "" {
				e.SetError(fmt.Errorf("empty value"))
				return
			}
			switch s[0] {
			case 'n':
				_ = s == "null" && e.SetA(bas.Nil) || e.SetError(fmt.Errorf("invalid value"))
			case 't', 'f':
				e.A = (*env)(e).valueOrError(strconv.ParseBool(s))
			case '[':
				a := []interface{}{}
				err := json.Unmarshal([]byte(s), &a)
				e.A = (*env)(e).valueOrError(bas.ValueOf(a), err)
			case '{':
				a := map[string]interface{}{}
				err := json.Unmarshal([]byte(s), &a)
				e.A = (*env)(e).valueOrError(bas.ValueOf(a), err)
			default:
				e.SetError(fmt.Errorf("invalid value"))
			}
		})).
		ToValue())
	bas.AddTopFunc("loadfile", func(e *bas.Env) {
		path := e.Str(0)
		if e.Shape(1, "Nb").IsTrue() && e.MustProgram().File != "" {
			path = filepath.Join(filepath.Dir(e.MustProgram().File), path)
		}
		e.A = MustRun(LoadFile(path, &LoadOptions{
			MaxStackSize: e.MustProgram().MaxStackSize,
			Globals:      e.MustProgram().Globals,
		}))
	})
	bas.AddTopValue("eval", bas.Func("eval", func(e *bas.Env) {
		p, err := LoadString(e.Str(0), &LoadOptions{
			Globals: e.Shape(1, "No").Object().ToMap(),
		})
		if err != nil {
			e.A = bas.Error(e, err)
		} else {
			e.A = (*env)(e).valueOrError(p.Run())
		}
	}).Object().
		SetProp("parse", bas.Func("parse", func(e *bas.Env) {
			e.A = (*env)(e).valueOrError(parser.Parse(e.Str(0), "eval.parse"))
		})).
		SetProp("op", bas.NewObject(0).
			SetProp("set", bas.Int(typ.OpSet)).
			SetProp("store", bas.Int(typ.OpStore)).
			SetProp("load", bas.Int(typ.OpLoad)).
			SetProp("add", bas.Int(typ.OpAdd)).
			SetProp("sub", bas.Int(typ.OpSub)).
			SetProp("mul", bas.Int(typ.OpMul)).
			SetProp("div", bas.Int(typ.OpDiv)).
			SetProp("idiv", bas.Int(typ.OpIDiv)).
			SetProp("inc", bas.Int(typ.OpInc)).
			SetProp("mod", bas.Int(typ.OpMod)).
			SetProp("len", bas.Int(typ.OpLen)).
			SetProp("next", bas.Int(typ.OpNext)).
			SetProp("not", bas.Int(typ.OpNot)).
			SetProp("eq", bas.Int(typ.OpEq)).
			SetProp("neq", bas.Int(typ.OpNeq)).
			SetProp("less", bas.Int(typ.OpLess)).
			SetProp("lesseq", bas.Int(typ.OpLessEq)).
			SetProp("ifnot", bas.Int(typ.OpJmpFalse)).
			SetProp("jmp", bas.Int(typ.OpJmp)).
			SetProp("function", bas.Int(typ.OpFunction)).
			SetProp("push", bas.Int(typ.OpPush)).
			SetProp("pushunpack", bas.Int(typ.OpPushUnpack)).
			SetProp("createarray", bas.Int(typ.OpCreateArray)).
			SetProp("createobject", bas.Int(typ.OpCreateObject)).
			SetProp("call", bas.Int(typ.OpCall)).
			SetProp("tailcall", bas.Int(typ.OpTailCall)).
			SetProp("isproto", bas.Int(typ.OpIsProto)).
			SetProp("slice", bas.Int(typ.OpSlice)).
			SetProp("ret", bas.Int(typ.OpRet)).
			SetProp("loadtop", bas.Int(typ.OpLoadTop)).
			SetProp("ext", bas.Int(typ.OpExt)).
			ToValue(),
		).
		ToValue())

	bas.AddTopValue("printf", bas.Func("printf", func(e *bas.Env) {
		bas.Fprintf(e.MustProgram().Stdout, e.Str(0), e.Stack()[1:]...)
	}))
	bas.AddTopValue("println", bas.Func("println", func(e *bas.Env) {
		for _, a := range e.Stack() {
			fmt.Fprint(e.MustProgram().Stdout, a.String(), " ")
		}
		fmt.Fprintln(e.MustProgram().Stdout)
	}))
	bas.AddTopValue("print", bas.Func("print", func(e *bas.Env) {
		for _, a := range e.Stack() {
			fmt.Fprint(e.MustProgram().Stdout, a.String())
		}
		fmt.Fprintln(e.MustProgram().Stdout)
	}))
	bas.AddTopValue("input", bas.NewObject(0).
		SetProp("int", bas.Func("int", func(e *bas.Env) {
			fmt.Fprint(e.MustProgram().Stdout, e.StrDefault(0, "", 0))
			var v int64
			fmt.Fscanf(e.MustProgram().Stdin, "%d", &v)
			e.A = bas.Int64(v)
		})).
		SetProp("num", bas.Func("num", func(e *bas.Env) {
			fmt.Fprint(e.MustProgram().Stdout, e.StrDefault(0, "", 0))
			var v float64
			fmt.Fscanf(e.MustProgram().Stdin, "%f", &v)
			e.A = bas.Float64(v)
		})).
		SetProp("str", bas.Func("str", func(e *bas.Env) {
			fmt.Fprint(e.MustProgram().Stdout, e.StrDefault(0, "", 0))
			var v string
			fmt.Fscanln(e.MustProgram().Stdin, &v)
			e.A = bas.Str(v)
		})).
		ToValue())

	rxMeta := bas.NewEmptyNativeMeta("RegExp", bas.NewObject(0).
		AddMethod("match", func(e *bas.Env) {
			rx := e.This().Interface().(*regexp.Regexp)
			e.A = bas.Bool(rx.MatchString(e.Str(0)))
		}).
		AddMethod("find", func(e *bas.Env) {
			rx := e.This().Interface().(*regexp.Regexp)
			e.A = bas.NewNative(rx.FindStringSubmatch(e.Str(0))).ToValue()
		}).
		AddMethod("findall", func(e *bas.Env) {
			rx := e.This().Interface().(*regexp.Regexp)
			e.A = bas.NewNative(rx.FindAllStringSubmatch(e.Str(0), e.IntDefault(1, -1))).ToValue()
		}).
		AddMethod("replace", func(e *bas.Env) {
			rx := e.This().Interface().(*regexp.Regexp)
			e.A = bas.Str(rx.ReplaceAllString(e.Str(0), e.Str(1)))
		}).
		SetPrototype(&bas.Proto.Native))

	bas.AddTopFunc("re", func(e *bas.Env) {
		if rx, err := regexp.Compile(e.Str(0)); err != nil {
			e.SetError(err)
		} else {
			e.A = bas.NewNativeWithMeta(rx, rxMeta).ToValue()
		}
	})

	fileMeta := bas.NewEmptyNativeMeta("File", bas.NewObject(0).
		AddMethod("name", func(e *bas.Env) {
			e.A = bas.Str(e.A.Interface().(*os.File).Name())
		}).
		AddMethod("seek", func(e *bas.Env) {
			e.A = (*env)(e).valueOrError(e.A.Interface().(*os.File).Seek(e.Int64(0), e.Int(1)))
		}).
		AddMethod("sync", func(e *bas.Env) {
			e.A = bas.Error(e, e.A.Interface().(*os.File).Sync())
		}).
		AddMethod("stat", func(e *bas.Env) {
			e.A = (*env)(e).valueOrError(e.A.Interface().(*os.File).Stat())
		}).
		AddMethod("truncate", func(e *bas.Env) {
			f := e.A.Interface().(*os.File)
			if err := f.Truncate(e.Int64(1)); err != nil {
				e.A = bas.Error(e, err)
			} else {
				e.A = (*env)(e).valueOrError(f.Seek(0, 2))
			}
		}).
		SetPrototype(bas.Proto.ReadWriteCloser.Proto))

	bas.AddTopValue("open", bas.Func("open", func(e *bas.Env) {
		path, flag, perm := e.Str(0), e.StrDefault(1, "r", 1), e.IntDefault(2, 0644)
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
		if err != nil {
			e.A = bas.Error(e, err)
			return
		}
		e.Object(-1).Set(bas.Zero, bas.ValueOf(f))

		e.A = bas.NewNativeWithMeta(f, fileMeta).ToValue()
	}).Object().
		AddMethod("close", func(e *bas.Env) {
			if f, _ := e.Object(-1).Get(bas.Zero).Interface().(*os.File); f != nil {
				e.A = bas.Error(e, f.Close())
			} else {
				panic("no opened file yet")
			}
		}).ToValue(),
	)

	bas.AddTopValue("math", bas.NewNamedObject("math", 0).
		SetProp("INF", bas.Float64(math.Inf(1))).
		SetProp("NEG_INF", bas.Float64(math.Inf(-1))).
		SetProp("PI", bas.Float64(math.Pi)).
		SetProp("E", bas.Float64(math.E)).
		SetProp("randomseed", bas.Func("randomseed", func(e *bas.Env) { rand.Seed(int64(e.IntDefault(0, 1))) })).
		SetProp("random", bas.Func("random", func(e *bas.Env) {
			switch len(e.Stack()) {
			case 2:
				ai, bi := e.Int64(0), e.Int64(1)
				e.A = bas.Int64(rand.Int63n(bi-ai+1) + ai)
			case 1:
				e.A = bas.Int64(rand.Int63n(e.Int64(0)))
			default:
				e.A = bas.Float64(rand.Float64())
			}
		})).
		SetProp("sqrt", bas.Func("sqrt", func(e *bas.Env) { e.A = bas.Float64(math.Sqrt(e.Float64(0))) })).
		SetProp("floor", bas.Func("floor", func(e *bas.Env) { e.A = bas.Float64(math.Floor(e.Float64(0))) })).
		SetProp("ceil", bas.Func("ceil", func(e *bas.Env) { e.A = bas.Float64(math.Ceil(e.Float64(0))) })).
		SetProp("min", bas.Func("min", func(e *bas.Env) { minMax(e, false) })).
		SetProp("max", bas.Func("max", func(e *bas.Env) { minMax(e, true) })).
		SetProp("pow", bas.Func("pow", func(e *bas.Env) { e.A = bas.Float64(math.Pow(e.Float64(0), e.Float64(1))) })).
		SetProp("abs", bas.Func("abs", func(e *bas.Env) {
			if e.A = e.Num(0); e.A.IsInt64() {
				if i := e.A.Int64(); i < 0 {
					e.A = bas.Int64(-i)
				}
			} else {
				e.A = bas.Float64(math.Abs(e.A.Float64()))
			}
		})).
		SetProp("remainder", bas.Func("remainder", func(e *bas.Env) { e.A = bas.Float64(math.Remainder(e.Float64(0), e.Float64(1))) })).
		SetProp("mod", bas.Func("mod", func(e *bas.Env) { e.A = bas.Float64(math.Mod(e.Float64(0), e.Float64(1))) })).
		SetProp("cos", bas.Func("cos", func(e *bas.Env) { e.A = bas.Float64(math.Cos(e.Float64(0))) })).
		SetProp("sin", bas.Func("sin", func(e *bas.Env) { e.A = bas.Float64(math.Sin(e.Float64(0))) })).
		SetProp("tan", bas.Func("tan", func(e *bas.Env) { e.A = bas.Float64(math.Tan(e.Float64(0))) })).
		SetProp("acos", bas.Func("acos", func(e *bas.Env) { e.A = bas.Float64(math.Acos(e.Float64(0))) })).
		SetProp("asin", bas.Func("asin", func(e *bas.Env) { e.A = bas.Float64(math.Asin(e.Float64(0))) })).
		SetProp("atan", bas.Func("atan", func(e *bas.Env) { e.A = bas.Float64(math.Atan(e.Float64(0))) })).
		SetProp("atan2", bas.Func("atan2", func(e *bas.Env) { e.A = bas.Float64(math.Atan2(e.Float64(0), e.Float64(1))) })).
		SetProp("ldexp", bas.Func("ldexp", func(e *bas.Env) { e.A = bas.Float64(math.Ldexp(e.Float64(0), e.Int(0))) })).
		SetProp("modf", bas.Func("modf", func(e *bas.Env) {
			a, b := math.Modf(e.Float64(0))
			e.A = bas.Array(bas.Float64(a), bas.Float64(b))
		})).
		ToValue())

	bas.AddTopValue("os", bas.NewNamedObject("os", 0).
		SetProp("stdout", bas.ValueOf(os.Stdout)).
		SetProp("stdin", bas.ValueOf(os.Stdin)).
		SetProp("stderr", bas.ValueOf(os.Stderr)).
		SetProp("pid", bas.Int(os.Getpid())).
		SetProp("numcpus", bas.Int(runtime.NumCPU())).
		SetProp("args", bas.ValueOf(os.Args)).
		SetProp("exit", bas.Func("exit", func(e *bas.Env) { os.Exit(e.IntDefault(0, 0)) })).
		SetProp("environ", bas.Func("environ", func(e *bas.Env) {
			if env := os.Environ(); e.Get(0).IsTrue() {
				obj := bas.NewObject(len(env))
				for _, e := range env {
					if idx := strings.Index(e, "="); idx > -1 {
						obj.SetProp(e[:idx], bas.Str(e[idx+1:]))
					}
				}
				e.A = obj.ToValue()
			} else {
				e.A = bas.ValueOf(env)
			}
		})).
		SetProp("shell", bas.Func("shell", func(e *bas.Env) {
			win := runtime.GOOS == "windows"
			p := exec.Command(internal.IfStr(win, "cmd", "sh"), internal.IfStr(win, "/c", "-c"), e.Str(0))
			opt := e.Shape(1, "No").Object()
			opt.Get(bas.Str("env")).AssertShape("No", "environment").Object().Foreach(func(k bas.Value, v *bas.Value) bool {
				p.Env = append(p.Env, k.String()+"="+v.String())
				return true
			})
			stdout := &bytes.Buffer{}
			p.Stdout, p.Stderr = stdout, stdout
			p.Dir = opt.GetDefault(bas.Str("dir"), bas.NullStr).AssertShape("s", "directory").Str()
			if tmp := opt.Get(bas.Str("stdout")); tmp != bas.Nil {
				p.Stdout = tmp.Writer()
			}
			if tmp := opt.Get(bas.Str("stderr")); tmp != bas.Nil {
				p.Stderr = tmp.Writer()
			}
			if tmp := opt.Get(bas.Str("stdin")); tmp != bas.Nil {
				p.Stdin = tmp.Reader()
			}

			to := opt.GetDefault(bas.Str("timeout"), bas.Float64(1<<30)).AssertNumber("timeout seconds").Duration()
			out := make(chan error)
			go func() { out <- p.Run() }()
			select {
			case r := <-out:
				if r != nil {
					e.A = bas.Error(e, r)
					return
				}
			case <-time.After(to):
				p.Process.Kill()
				e.A = bas.Error(e, fmt.Errorf("os.shell timeout: %v", e.Get(0)))
				return
			}
			e.A = bas.Bytes(stdout.Bytes())
		})).
		SetProp("readdir", bas.Func("readdir", func(e *bas.Env) {
			e.A = (*env)(e).valueOrError(ioutil.ReadDir(e.Str(0)))
		})).
		SetProp("remove", bas.Func("remove", func(e *bas.Env) {
			path := e.Str(0)
			fi, err := os.Stat(path)
			if err != nil {
				e.A = bas.Error(e, err)
			} else if fi.IsDir() {
				e.A = bas.Error(e, os.RemoveAll(path))
			} else {
				e.A = bas.Error(e, os.Remove(path))
			}
		})).
		SetProp("pstat", bas.Func("pstat", func(e *bas.Env) {
			fi, err := os.Stat(e.Str(0))
			_ = err == nil && e.SetA(bas.ValueOf(fi)) || e.SetA(bas.Nil)
		})).
		ToValue())

	bas.AddTopValue("sync", bas.NewNamedObject("sync", 0).
		SetProp("mutex", bas.Func("mutex", func(e *bas.Env) { e.A = bas.ValueOf(&sync.Mutex{}) })).
		SetProp("rwmutex", bas.Func("rwmutex", func(e *bas.Env) { e.A = bas.ValueOf(&sync.RWMutex{}) })).
		SetProp("waitgroup", bas.Func("waitgroup", func(e *bas.Env) { e.A = bas.ValueOf(&sync.WaitGroup{}) })).
		ToValue())

	bas.AddTopValue("base64", bas.NewObject(0).
		SetProp("std", bas.ValueOf(base64.StdEncoding)).SetProp("rawstd", bas.ValueOf(base64.RawStdEncoding)).
		SetProp("url", bas.ValueOf(base64.URLEncoding)).SetProp("rawurl", bas.ValueOf(base64.RawURLEncoding)).
		ToValue())

	bas.AddTopValue("base32", bas.NewObject(0).
		SetProp("std", bas.ValueOf(base32.StdEncoding)).SetProp("rawstd", bas.ValueOf(base32.StdEncoding.WithPadding(-1))).
		SetProp("hex", bas.ValueOf(base32.HexEncoding)).SetProp("rawhex", bas.ValueOf(base32.HexEncoding.WithPadding(-1))).
		ToValue())

	bas.AddTopValue("time", bas.Func("time", func(e *bas.Env) {
		e.A = bas.Float64(float64(time.Now().UnixNano()) / 1e9)
	}).Object().
		SetProp("sleep", bas.Func("sleep", func(e *bas.Env) { time.Sleep(time.Duration(e.Float64(0)*1e6) * 1e3) })).
		SetProp("ymd", bas.Func("ymd", func(e *bas.Env) {
			e.A = bas.ValueOf(time.Date(e.Int(0), time.Month(e.Int(1)), e.Int(2),
				e.IntDefault(3, 0), e.IntDefault(4, 0), e.IntDefault(5, 0), e.IntDefault(6, 0), time.UTC))
		})).
		SetProp("clock", bas.Func("clock", func(e *bas.Env) {
			x := time.Now()
			s := *(*[2]int64)(unsafe.Pointer(&x))
			e.A = bas.Float64(float64(s[1]) / 1e9)
		})).
		SetProp("now", bas.Func("now", func(e *bas.Env) { e.A = bas.ValueOf(time.Now()) })).
		SetProp("after", bas.Func("after", func(e *bas.Env) { e.A = bas.ValueOf(time.After(time.Duration(e.Float64(0)*1e6) * 1e3)) })).
		SetProp("parse", bas.Func("parse", func(e *bas.Env) {
			e.A = (*env)(e).valueOrError(time.Parse(getTimeFormat(e.Str(0)), e.Str(1)))
		})).
		SetProp("format", bas.Func("format", func(e *bas.Env) {
			tt, ok := e.Get(1).Interface().(time.Time)
			if !ok {
				if t := e.Get(1); t.Type() == typ.Number {
					tt = time.Unix(0, int64(t.Float64()*1e9))
				} else {
					tt = time.Now()
				}
			}
			e.A = bas.Str(tt.Format(getTimeFormat(e.StrDefault(0, "", 0))))
		})).
		ToValue())

	httpLib := bas.Func("http", func(e *bas.Env) {
		args := e.Object(0)
		method := strings.ToUpper(args.GetDefault(bas.Str("method"), bas.Str("GET")).Str())

		u, err := url.Parse(args.Get(bas.Str("url")).AssertString("http URL"))
		if err != nil {
			e.A = bas.Error(e, err)
			return
		}

		addKV := func(k string, add func(k, v string)) {
			x := args.Get(bas.Str(k)).AssertShape("No", k)
			x.Object().Foreach(func(k bas.Value, v *bas.Value) bool { add(k.String(), v.String()); return true })
		}

		additionalQueries := u.Query()
		addKV("query", additionalQueries.Add) // append queries to url
		u.RawQuery = additionalQueries.Encode()

		var bodyReader io.Reader
		dataForm, urlForm, jsonForm := (*multipart.Writer)(nil), false, false

		if j := args.Get(bas.Str("json")); j != bas.Nil {
			bodyReader = strings.NewReader(j.JSONString())
			jsonForm = true
		} else {
			var form url.Values
			if args.Contains(bas.Str("form")) {
				form = url.Values{}
				addKV("form", form.Add) // check "form"
			}
			urlForm = len(form) > 0
			if urlForm {
				bodyReader = strings.NewReader(form.Encode())
			} else if rd := args.Get(bas.Str("data")); rd != bas.Nil {
				bodyReader = rd.Reader()
			}
		}

		if bodyReader == nil && (method == "POST" || method == "PUT") {
			// Check form-data
			payload := bytes.Buffer{}
			writer := multipart.NewWriter(&payload)
			outError := (error)(nil)
			args.Get(bas.Str("multipart")).AssertShape("No", "multipart").Object().
				Foreach(func(k bas.Value, v *bas.Value) bool {
					key, rd := k.String(), *v
					rd.AssertShape("<s,(s,R)>", "http multipart form data format")
					if rd.Type() == typ.Native && rd.Len() == 2 { // [filename, reader]
						fn := rd.Native().Get(0).Str()
						if part, err := writer.CreateFormFile(key, fn); err != nil {
							outError = fmt.Errorf("%s: %v", fn, err)
							return false
						} else if _, err = io.Copy(part, rd.Native().Get(1).Reader()); err != nil {
							outError = fmt.Errorf("%s: %v", fn, err)
							return false
						}
					} else {
						if part, err := writer.CreateFormField(key); err != nil {
							outError = fmt.Errorf("%s: %v", key, err)
							return false
						} else if _, err = io.Copy(part, rd.Reader()); err != nil {
							outError = fmt.Errorf("%s: %v", key, err)
							return false
						}
					}
					return true
				})
			if outError != nil {
				e.A = bas.Error(e, err)
				return
			}
			if err := writer.Close(); err != nil {
				e.A = bas.Error(e, err)
				return
			}
			if payload.Len() > 0 {
				bodyReader = &payload
				dataForm = writer
			}
		}

		req, err := http.NewRequest(method, u.String(), bodyReader)
		if err != nil {
			e.A = bas.Error(e, err)
			return
		}

		switch {
		case urlForm:
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		case jsonForm:
			req.Header.Add("Content-Type", "application/json")
		case dataForm != nil:
			req.Header.Add("Content-Type", dataForm.FormDataContentType())
		}

		addKV("header", req.Header.Add) // append headers

		// Construct HTTP client
		client := &http.Client{}
		client.Timeout = args.GetDefault(bas.Str("timeout"), bas.Float64(1<<30)).AssertNumber("timeout seconds").Duration()
		if v := args.Get(bas.Str("jar")); v.Type() == typ.Native {
			client.Jar, _ = v.Interface().(http.CookieJar)
		}
		if !args.Get(bas.Str("noredirect")).IsFalse() {
			client.CheckRedirect = func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			}
		}
		if p := args.GetDefault(bas.Str("proxy"), bas.NullStr).Str(); p != "" {
			client.Transport = &http.Transport{
				Proxy: func(r *http.Request) (*url.URL, error) { return url.Parse(p) },
			}
		}
		send := func(e *bas.Env) (code, headers, buf, jar bas.Value) {
			resp, err := client.Do(req)
			if err != nil {
				err := bas.Error(e, err)
				return err, err, err, err
			}

			if args.Get(bas.Str("br")).IsFalse() {
				resp.Body.Close()
			} else {
				buf = bas.NewNativeWithMeta(resp.Body, &bas.Proto.ReadCloser).ToValue()
			}
			return bas.Int(resp.StatusCode), bas.ValueOf(resp.Header), buf, bas.ValueOf(client.Jar)
		}
		if f := args.Get(bas.Str("async")); f.IsObject() {
			go func(e *bas.Env) {
				code, hdr, buf, jar := send(e)
				f.Object().Call(e, code, hdr, buf, jar)
			}(e.Copy())
			return
		}
		e.A = bas.Array(send(e))
	}).Object().
		SetProp("urlescape", bas.Func("urlescape", func(e *bas.Env) { e.A = bas.Str(url.QueryEscape(e.Str(0))) })).
		SetProp("urlunescape", bas.Func("urlunescape", func(e *bas.Env) {
			e.A = (*env)(e).valueOrError(url.QueryUnescape(e.Str(0)))
		})).
		SetProp("pathescape", bas.Func("pathescape", func(e *bas.Env) { e.A = bas.Str(url.PathEscape(e.Str(0))) })).
		SetProp("pathunescape", bas.Func("pathunescape", func(e *bas.Env) {
			e.A = (*env)(e).valueOrError(url.PathUnescape(e.Str(0)))
		}))
	for _, m := range []string{"get", "post", "put", "delete", "head", "patch"} {
		m := m
		httpLib = httpLib.AddMethod(m, func(e *bas.Env) {
			ex := e.Shape(1, "No").Object()
			e.A = e.Object(-1).Call(e, bas.NewObject(0).SetProp("method", bas.Str(m)).SetProp("url", e.Get(0)).Merge(ex).ToValue())
		})
	}
	bas.AddTopValue("http", httpLib.ToValue())

	bufferMeta := bas.NewEmptyNativeMeta("Buffer", bas.NewObject(0).
		SetPrototype(bas.Proto.ReadWriter.Proto).
		AddMethod("reset", func(e *bas.Env) { e.A.Interface().(*internal.LimitedBuffer).Reset() }).
		AddMethod("value", func(e *bas.Env) { e.A = bas.UnsafeStr(e.A.Interface().(*internal.LimitedBuffer).Bytes()) }).
		AddMethod("bytes", func(e *bas.Env) { e.A = bas.Bytes(e.A.Interface().(*internal.LimitedBuffer).Bytes()) }))

	bas.AddTopFunc("buffer", func(e *bas.Env) {
		b := &internal.LimitedBuffer{Limit: e.IntDefault(1, 0)}
		bas.Write(b, e.Get(0))
		e.A = bas.NewNativeWithMeta(b, bufferMeta).ToValue()
	})
}

func minMax(e *bas.Env, max bool) {
	v := e.Get(0)
	for i := 1; i < e.Size(); i++ {
		if x := e.Get(i); v.Less(x) == max {
			v = x
		}
	}
	e.A = v
}

var timeFormatMapping = map[interface{}]string{
	"ansic": time.ANSIC, "ANSIC": time.ANSIC,
	"unixdate": time.UnixDate, "UnixDate": time.UnixDate,
	"rubydate": time.RubyDate, "RubyDate": time.RubyDate,
	"rfc822": time.RFC822, "RFC822": time.RFC822,
	"rfc822z": time.RFC822Z, "RFC822Z": time.RFC822Z,
	"rfc850": time.RFC850, "RFC850": time.RFC850,
	"rfc1123": time.RFC1123, "RFC1123": time.RFC1123,
	"rfc1123z": time.RFC1123Z, "RFC1123Z": time.RFC1123Z,
	"rfc3339": time.RFC3339, "RFC3339": time.RFC3339,
	"rfc3339nano": time.RFC3339Nano, "RFC3339Nano": time.RFC3339Nano,
	"kitchen": time.Kitchen, "Kitchen": time.Kitchen,
	"stamp": time.Stamp, "Stamp": time.Stamp,
	"stampmilli": time.StampMilli, "StampMilli": time.StampMilli,
	"stampmicro": time.StampMicro, "StampMicro": time.StampMicro,
	"stampnano": time.StampNano, "StampNano": time.StampNano,
	'd': "02", 'D': "Mon", 'j': "2", 'l': "Monday", 'F': "January", 'z': "002", 'm': "01",
	'M': "Jan", 'n': "1", 'Y': "2006", 'y': "06", 'a': "pm", 'A': "PM", 'g': "3", 'G': "15",
	'h': "03", 'H': "15", 'i': "04", 's': "05", 'u': "05.000000", 'v': "05.000", 'O': "+0700",
	'P': "-07:00", 'T': "MST",
	'c': "2006-01-02T15:04:05-07:00",       //	ISO 860,
	'r': "Mon, 02 Jan 2006 15:04:05 -0700", //	RFC 282,
}

func getTimeFormat(f string) string {
	if tf, ok := timeFormatMapping[f]; ok {
		return tf
	}
	buf := bytes.Buffer{}
	for len(f) > 0 {
		r, sz := utf8.DecodeRuneInString(f)
		if sz == 0 {
			break
		}
		if tf, ok := timeFormatMapping[r]; ok {
			buf.WriteString(tf)
		} else {
			buf.WriteRune(r)
		}
		f = f[sz:]
	}
	return buf.String()
}

type env bas.Env

func (e *env) valueOrError(v interface{}, err error) bas.Value {
	if err != nil {
		return bas.Error((*bas.Env)(e), err)
	}
	return bas.ValueOf(v)
}
