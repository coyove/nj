package nj

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
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
	bas.Globals.SetProp("json", bas.NamedObject("json", 0).
		SetMethod("stringify", func(e *bas.Env) { e.A = bas.Str(e.Get(0).JSONString()) }).
		SetMethod("dump", func(e *bas.Env) { e.Get(1).Stringify(bas.NewWriter(e.Get(0)), typ.MarshalToJSON) }).
		SetMethod("parse", func(e *bas.Env) {
			s := strings.TrimSpace(e.Str(0))
			f := parser.ParseJSON
			if e.Get(1).IsTrue() {
				f = ParseStrictJSON
			}
			v, err := f(s)
			internal.PanicErr(err)
			e.A = v
		}).
		ToValue())
	bas.Globals.SetMethod("loadfile", func(e *bas.Env) {
		path := e.Str(0)
		if e.Get(1).Maybe().Bool() && e.Global.File != "" {
			path = filepath.Join(filepath.Dir(e.Global.File), path)
		}
		e.A = MustRun(LoadFile(path, &e.Global.Environment))
	})
	bas.Globals.SetMethod("eval", func(e *bas.Env) {
		opts := e.Get(1).Maybe().Object(nil)
		if opts.Prop("ast").IsTrue() {
			v, err := parser.Parse(e.Str(0), "")
			internal.PanicErr(err)
			e.A = bas.ValueOf(v)
			return
		}
		p, err := LoadString(e.Str(0), &bas.Environment{Globals: opts.Prop("globals").Maybe().Object(nil)})
		internal.PanicErr(err)
		v, err := p.Run()
		internal.PanicErr(err)
		_ = opts.Prop("returnglobals").IsTrue() && e.SetA(p.LocalsObject().ToValue()) || e.SetA(v)
	})

	bas.Globals.SetProp("stdout", bas.ValueOf(os.Stdout))
	bas.Globals.SetProp("stdin", bas.ValueOf(os.Stdin))
	bas.Globals.SetProp("stderr", bas.ValueOf(os.Stderr))
	bas.Globals.SetMethod("scanln", func(env *bas.Env) {
		prompt, n := env.Get(0), env.Get(1)
		fmt.Fprint(env.Global.Stdout, prompt.Maybe().Str(""))
		var results []bas.Value
		var r io.Reader = env.Global.Stdin
		for i := n.Maybe().Int64(1); i > 0; i-- {
			var s string
			if _, err := fmt.Fscan(r, &s); err != nil {
				break
			}
			results = append(results, bas.Str(s))
		}
		env.A = bas.Array(results...)
	})
	bas.Globals.SetMethod("sleep", func(e *bas.Env) { time.Sleep(time.Duration(e.Float64(0)*1e6) * 1e3) })
	bas.Globals.SetMethod("Go_time", func(e *bas.Env) {
		if e.Size() > 0 {
			e.A = bas.ValueOf(time.Date(e.Int(0), time.Month(e.Int(1)), e.Int(2),
				e.Get(3).Maybe().Int(0), e.Get(4).Maybe().Int(0), e.Get(5).Maybe().Int(0), e.Get(6).Maybe().Int(0), time.UTC))
		} else {
			e.A = bas.ValueOf(time.Now())
		}
	})
	bas.Globals.SetMethod("clock", func(e *bas.Env) {
		x := time.Now()
		s := *(*[2]int64)(unsafe.Pointer(&x))
		e.A = bas.Float64(float64(s[1]) / 1e9)
	})
	bas.Globals.SetMethod("exit", func(e *bas.Env) { os.Exit(e.Int(0)) })
	bas.Globals.SetMethod("chr", func(e *bas.Env) { e.A = bas.Rune(rune(e.Int(0))) })
	bas.Globals.SetMethod("byte", func(e *bas.Env) { e.A = bas.Byte(byte(e.Int(0))) })
	bas.Globals.SetMethod("ord", func(e *bas.Env) { r, _ := utf8.DecodeRuneInString(e.Str(0)); e.A = bas.Int64(int64(r)) })

	bas.Globals.SetProp("re", bas.Func("RegExp", func(e *bas.Env) {
		rx := regexp.MustCompile(e.Str(0))
		e.A = bas.NewObject(1).SetPrototype(e.A.Object()).SetProp("_rx", bas.ValueOf(rx)).ToValue()
	}).Object().
		SetMethod("match", func(e *bas.Env) {
			e.A = bas.Bool(e.ThisProp("_rx").(*regexp.Regexp).MatchString(e.Str(0)))
		}).
		SetMethod("find", func(e *bas.Env) {
			e.A = bas.NewNative(e.ThisProp("_rx").(*regexp.Regexp).FindStringSubmatch(e.Str(0))).ToValue()
		}).
		SetMethod("findall", func(e *bas.Env) {
			m := e.ThisProp("_rx").(*regexp.Regexp).FindAllStringSubmatch(e.Str(0), e.Get(1).Maybe().Int(-1))
			e.A = bas.NewNative(m).ToValue()
		}).
		SetMethod("replace", func(e *bas.Env) {
			e.A = bas.Str(e.ThisProp("_rx").(*regexp.Regexp).ReplaceAllString(e.Str(0), e.Str(1)))
		}).
		ToValue())

	bas.Globals.SetProp("open", bas.Func("open", func(e *bas.Env) {
		path, flag, perm := e.Str(0), e.Get(1).Maybe().Str("r"), e.Get(2).Maybe().Int64(0644)
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
		e.Object(-1).Set(bas.Zero, bas.ValueOf(f))

		e.A = bas.NamedObject("File", 0).
			SetProp("_f", bas.ValueOf(f)).
			SetProp("path", bas.Str(f.Name())).
			SetMethod("sync", func(e *bas.Env) { internal.PanicErr(e.ThisProp("_f").(*os.File).Sync()) }).
			SetMethod("stat", func(e *bas.Env) { e.A = valueOrPanic(e.ThisProp("_f").(*os.File).Stat()) }).
			SetMethod("truncate", func(e *bas.Env) {
				f := e.ThisProp("_f").(*os.File)
				internal.PanicErr(f.Truncate(e.Int64(1)))
				t, err := f.Seek(0, 2)
				internal.PanicErr(err)
				e.A = bas.Int64(t)
			}).
			SetPrototype(bas.Proto.ReadWriteSeekCloser).
			ToValue()
	}).Object().
		SetMethod("close", func(e *bas.Env) {
			if f, _ := e.Object(-1).Find(bas.Zero).Interface().(*os.File); f != nil {
				internal.PanicErr(f.Close())
			} else {
				internal.Panic("no opened file yet")
			}
		}).ToValue(),
	)

	bas.Globals.SetProp("math", bas.NamedObject("math", 0).
		SetProp("INF", bas.Float64(math.Inf(1))).
		SetProp("NEG_INF", bas.Float64(math.Inf(-1))).
		SetProp("PI", bas.Float64(math.Pi)).
		SetProp("E", bas.Float64(math.E)).
		SetMethod("randomseed", func(e *bas.Env) { rand.Seed(e.Get(0).Maybe().Int64(1)) }).
		SetMethod("random", func(e *bas.Env) {
			switch len(e.Stack()) {
			case 2:
				ai, bi := e.Int64(0), e.Int64(1)
				e.A = bas.Int64(rand.Int63n(bi-ai+1) + ai)
			case 1:
				e.A = bas.Int64(rand.Int63n(e.Int64(0)))
			default:
				e.A = bas.Float64(rand.Float64())
			}
		}).
		SetMethod("sqrt", func(e *bas.Env) { e.A = bas.Float64(math.Sqrt(e.Float64(0))) }).
		SetMethod("floor", func(e *bas.Env) { e.A = bas.Float64(math.Floor(e.Float64(0))) }).
		SetMethod("ceil", func(e *bas.Env) { e.A = bas.Float64(math.Ceil(e.Float64(0))) }).
		SetMethod("min", func(e *bas.Env) { mathMinMax(e, false) }).
		SetMethod("max", func(e *bas.Env) { mathMinMax(e, true) }).
		SetMethod("pow", func(e *bas.Env) { e.A = bas.Float64(math.Pow(e.Float64(0), e.Float64(1))) }).
		SetMethod("abs", func(e *bas.Env) {
			if e.A = e.Num(0); e.A.IsInt64() {
				if i := e.A.Int64(); i < 0 {
					e.A = bas.Int64(-i)
				}
			} else {
				e.A = bas.Float64(math.Abs(e.A.Float64()))
			}
		}).
		SetMethod("remainder", func(e *bas.Env) { e.A = bas.Float64(math.Remainder(e.Float64(0), e.Float64(1))) }).
		SetMethod("mod", func(e *bas.Env) { e.A = bas.Float64(math.Mod(e.Float64(0), e.Float64(1))) }).
		SetMethod("cos", func(e *bas.Env) { e.A = bas.Float64(math.Cos(e.Float64(0))) }).
		SetMethod("sin", func(e *bas.Env) { e.A = bas.Float64(math.Sin(e.Float64(0))) }).
		SetMethod("tan", func(e *bas.Env) { e.A = bas.Float64(math.Tan(e.Float64(0))) }).
		SetMethod("acos", func(e *bas.Env) { e.A = bas.Float64(math.Acos(e.Float64(0))) }).
		SetMethod("asin", func(e *bas.Env) { e.A = bas.Float64(math.Asin(e.Float64(0))) }).
		SetMethod("atan", func(e *bas.Env) { e.A = bas.Float64(math.Atan(e.Float64(0))) }).
		SetMethod("atan2", func(e *bas.Env) { e.A = bas.Float64(math.Atan2(e.Float64(0), e.Float64(1))) }).
		SetMethod("ldexp", func(e *bas.Env) { e.A = bas.Float64(math.Ldexp(e.Float64(0), e.Int(0))) }).
		SetMethod("modf", func(e *bas.Env) {
			a, b := math.Modf(e.Float64(0))
			e.A = bas.Array(bas.Float64(a), bas.Float64(b))
		}).
		SetPrototype(bas.Proto.StaticObject).
		ToValue())

	bas.Globals.SetProp("os", bas.NamedObject("os", 0).
		SetProp("pid", bas.Int(os.Getpid())).
		SetProp("numcpus", bas.Int(runtime.NumCPU())).
		SetProp("args", bas.ValueOf(os.Args)).
		SetMethod("environ", func(e *bas.Env) {
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
		}).
		SetMethod("shell", func(e *bas.Env) {
			win := runtime.GOOS == "windows"
			p := exec.Command(internal.IfStr(win, "cmd", "sh"), internal.IfStr(win, "/c", "-c"), e.Str(0))
			opt := e.Get(1).Maybe().Object(nil)
			opt.Prop("env").Maybe().Object(nil).Foreach(func(k bas.Value, v *bas.Value) bool {
				p.Env = append(p.Env, k.String()+"="+v.String())
				return true
			})
			stdout := &bytes.Buffer{}
			p.Stdout, p.Stderr = stdout, stdout
			p.Dir = opt.Prop("dir").Maybe().Str("")
			if tmp := opt.Prop("stdout"); tmp != bas.Nil {
				p.Stdout = bas.NewWriter(tmp)
			}
			if tmp := opt.Prop("stderr"); tmp != bas.Nil {
				p.Stderr = bas.NewWriter(tmp)
			}
			if tmp := opt.Prop("stdin"); tmp != bas.Nil {
				p.Stdin = bas.NewReader(tmp)
			}

			out := make(chan error)
			go func() { out <- p.Run() }()
			select {
			case r := <-out:
				internal.PanicErr(r)
			case <-time.After(time.Duration(opt.Prop("timeout").Maybe().Float64(1<<52)*1e6) * 1e3):
				p.Process.Kill()
				panic("timeout")
			}
			e.A = bas.Bytes(stdout.Bytes())
		}).
		SetMethod("readdir", func(e *bas.Env) { e.A = valueOrPanic(ioutil.ReadDir(e.Str(0))) }).
		SetMethod("remove", func(e *bas.Env) {
			path := e.Str(0)
			fi, err := os.Stat(path)
			internal.PanicErr(err)
			if fi.IsDir() {
				internal.PanicErr(os.RemoveAll(path))
			} else {
				internal.PanicErr(os.Remove(path))
			}
		}).
		SetMethod("pstat", func(e *bas.Env) {
			fi, err := os.Stat(e.Str(0))
			_ = err == nil && e.SetA(bas.ValueOf(fi)) || e.SetA(bas.Nil)
		}).
		SetPrototype(bas.Proto.StaticObject).
		ToValue())

	bas.Globals.SetProp("sync", bas.NamedObject("sync", 0).
		SetMethod("mutex", func(e *bas.Env) { e.A = bas.ValueOf(&sync.Mutex{}) }).
		SetMethod("rwmutex", func(e *bas.Env) { e.A = bas.ValueOf(&sync.RWMutex{}) }).
		SetMethod("waitgroup", func(e *bas.Env) { e.A = bas.ValueOf(&sync.WaitGroup{}) }).
		SetPrototype(bas.Proto.StaticObject).
		ToValue())

	encDecProto := bas.NamedObject("EncodeDecode", 0).
		SetMethod("encode", func(e *bas.Env) {
			i := e.ThisProp("_e")
			e.A = bas.Str(i.(interface{ EncodeToString([]byte) string }).EncodeToString(bas.ToReadonlyBytes(e.Get(0))))
		}).
		SetMethod("decode", func(e *bas.Env) {
			i := e.ThisProp("_e")
			v, err := i.(interface{ DecodeString(string) ([]byte, error) }).DecodeString(e.Str(0))
			internal.PanicErr(err)
			e.A = bas.Bytes(v)
		}).
		SetPrototype(bas.NamedObject("EncoderDecoder", 0).
			SetMethod("encoder", func(e *bas.Env) {
				enc := bas.Nil
				buf := &bytes.Buffer{}
				switch encoding := e.ThisProp("_e").(type) {
				default:
					enc = bas.ValueOf(hex.NewEncoder(buf))
				case *base32.Encoding:
					enc = bas.ValueOf(base32.NewEncoder(encoding, buf))
				case *base64.Encoding:
					enc = bas.ValueOf(base64.NewEncoder(encoding, buf))
				}
				e.A = bas.NamedObject("Encoder", 0).
					SetProp("_f", bas.ValueOf(enc)).
					SetProp("_b", bas.ValueOf(buf)).
					SetMethod("value", func(e *bas.Env) { e.A = bas.Str(e.ThisProp("_b").(*bytes.Buffer).String()) }).
					SetMethod("bytes", func(e *bas.Env) { e.A = bas.Bytes(e.ThisProp("_b").(*bytes.Buffer).Bytes()) }).
					SetPrototype(bas.Proto.WriteCloser).
					ToValue()
			}).
			SetMethod("decoder", func(e *bas.Env) {
				src := bas.NewReader(e.Get(0))
				dec := bas.Nil
				switch encoding := e.ThisProp("_e").(type) {
				case *base64.Encoding:
					dec = bas.ValueOf(base64.NewDecoder(encoding, src))
				case *base32.Encoding:
					dec = bas.ValueOf(base32.NewDecoder(encoding, src))
				default:
					dec = bas.ValueOf(hex.NewDecoder(src))
				}
				e.A = bas.NamedObject("Decoder", 0).
					SetProp("_f", bas.ValueOf(dec)).
					SetPrototype(bas.Proto.Reader).
					ToValue()
			}))

	bas.Globals.SetProp("hex", bas.NamedObject("hex", 0).SetPrototype(encDecProto.Prototype()).ToValue())
	bas.Globals.SetProp("base64", bas.NamedObject("base64", 0).
		SetProp("std", bas.NewObject(1).SetPrototype(encDecProto).SetProp("_e", bas.ValueOf(base64.StdEncoding)).ToValue()).
		SetProp("url", bas.NewObject(1).SetPrototype(encDecProto).SetProp("_e", bas.ValueOf(base64.URLEncoding)).ToValue()).
		SetProp("std2", bas.NewObject(1).SetPrototype(encDecProto).SetProp("_e", bas.ValueOf(base64.StdEncoding.WithPadding(-1))).ToValue()).
		SetProp("url2", bas.NewObject(1).SetPrototype(encDecProto).SetProp("_e", bas.ValueOf(base64.URLEncoding.WithPadding(-1))).ToValue()).
		SetPrototype(encDecProto).
		ToValue())
	bas.Globals.SetProp("base32", bas.NamedObject("base32", 0).
		SetProp("std", bas.NewObject(1).SetPrototype(encDecProto).SetProp("_e", bas.ValueOf(base32.StdEncoding)).ToValue()).
		SetProp("hex", bas.NewObject(1).SetPrototype(encDecProto).SetProp("_e", bas.ValueOf(base32.HexEncoding)).ToValue()).
		SetProp("std2", bas.NewObject(1).SetPrototype(encDecProto).SetProp("_e", bas.ValueOf(base32.StdEncoding.WithPadding(-1))).ToValue()).
		SetProp("hex2", bas.NewObject(1).SetPrototype(encDecProto).SetProp("_e", bas.ValueOf(base32.HexEncoding.WithPadding(-1))).ToValue()).
		SetPrototype(encDecProto).
		ToValue())

	bas.Globals.SetProp("time", bas.Func("time", func(e *bas.Env) {
		e.A = bas.Float64(float64(time.Now().UnixNano()) / 1e9)
	}).Object().
		SetMethod("now", func(e *bas.Env) { e.A = bas.ValueOf(time.Now()) }).
		SetMethod("after", func(e *bas.Env) { e.A = bas.ValueOf(time.After(time.Duration(e.Float64(0)*1e6) * 1e3)) }).
		SetMethod("parse", func(e *bas.Env) { e.A = valueOrPanic(time.Parse(getTimeFormat(e.Str(0)), e.Str(1))) }).
		SetMethod("format", func(e *bas.Env) {
			tt, ok := e.Get(1).Interface().(time.Time)
			if !ok {
				if t := e.Get(1); t.Type() == typ.Number {
					tt = time.Unix(0, int64(t.Float64()*1e9))
				} else {
					tt = time.Now()
				}
			}
			e.A = bas.Str(tt.Format(getTimeFormat(e.Get(0).Maybe().Str(""))))
		}).
		ToValue())

	bas.Globals.SetProp("url", bas.Func("url", nil).Object().
		SetMethod("escape", func(e *bas.Env) { e.A = bas.Str(url.QueryEscape(e.Str(0))) }).
		SetMethod("unescape", func(e *bas.Env) { e.A = valueOrPanic(url.QueryUnescape(e.Str(0))) }).
		ToValue())

	httpLib := bas.Func("http", func(e *bas.Env) {
		args := e.Object(0)
		to := args.Prop("timeout").Maybe().Float64(1 << 30)
		method := strings.ToUpper(args.Find(bas.Str("method")).Maybe().Str("GET"))

		u, err := url.Parse(args.Find(bas.Str("url")).Maybe().Str("bad://%url%"))
		internal.PanicErr(err)

		addKV := func(k string, add func(k, v string)) {
			x := args.Find(bas.Str(k))
			x.Maybe().Object(nil).Foreach(func(k bas.Value, v *bas.Value) bool { add(k.String(), v.String()); return true })
		}

		additionalQueries := u.Query()
		addKV("query", additionalQueries.Add) // append queries to url
		u.RawQuery = additionalQueries.Encode()

		var bodyReader io.Reader
		dataForm, urlForm, jsonForm := (*multipart.Writer)(nil), false, false

		if j := args.Prop("json"); j != bas.Nil {
			bodyReader = strings.NewReader(j.JSONString())
			jsonForm = true
		} else {
			var form url.Values
			if args.Contains(bas.Str("form"), false) {
				form = url.Values{}
				addKV("form", form.Add) // check "form"
			}
			urlForm = len(form) > 0
			if urlForm {
				bodyReader = strings.NewReader(form.Encode())
			} else if rd := args.Prop("data"); rd != bas.Nil {
				bodyReader = bas.NewReader(rd)
			}
		}

		if bodyReader == nil && (method == "POST" || method == "PUT") {
			// Check form-data
			payload := bytes.Buffer{}
			writer := multipart.NewWriter(&payload)
			if x := args.Prop("multipart"); x.Type() == typ.Object {
				x.Object().Foreach(func(k bas.Value, v *bas.Value) bool {
					key, rd := k.String(), *v
					if rd.Type() == typ.Native && bas.Len(rd) == 2 { // [filename, reader]
						part, err := writer.CreateFormFile(key, rd.Native().Get(0).Maybe().Str(""))
						internal.PanicErr(err)
						_, err = io.Copy(part, bas.NewReader(rd.Native().Get(1)))
						internal.PanicErr(err)
					} else {
						part, err := writer.CreateFormField(key)
						internal.PanicErr(err)
						_, err = io.Copy(part, bas.NewReader(rd))
						internal.PanicErr(err)
					}
					return true
				})
			}
			internal.PanicErr(writer.Close())
			if payload.Len() > 0 {
				bodyReader = &payload
				dataForm = writer
			}
		}

		req, err := http.NewRequest(method, u.String(), bodyReader)
		internal.PanicErr(err)

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
		client.Timeout = time.Duration(to * float64(time.Second))
		if v := args.Find(bas.Str("jar")); v.Type() == typ.Native {
			client.Jar, _ = v.Interface().(http.CookieJar)
		}
		if !args.Find(bas.Str("noredirect")).IsFalse() {
			client.CheckRedirect = func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			}
		}
		if p := args.Find(bas.Str("proxy")).Maybe().Str(""); p != "" {
			client.Transport = &http.Transport{
				Proxy: func(r *http.Request) (*url.URL, error) { return url.Parse(p) },
			}
		}
		send := func(e *bas.Env, panic bool) (code, headers, buf, jar bas.Value) {
			resp, err := client.Do(req)
			if panic {
				internal.PanicErr(err)
			} else if err != nil {
				return bas.Error(e, err), bas.Error(e, err), bas.Error(e, err), bas.Error(e, err)
			}

			if args.Prop("bodyreader").IsFalse() && args.Prop("br").IsFalse() {
				resp.Body.Close()
			} else {
				buf = bas.NewObject(1).SetProp("_f", bas.ValueOf(resp.Body)).SetPrototype(bas.Proto.ReadCloser).ToValue()
			}

			hdr := map[string]string{}
			for k := range resp.Header {
				hdr[k] = resp.Header.Get(k)
			}
			return bas.Int(resp.StatusCode), bas.ValueOf(hdr), buf, bas.ValueOf(client.Jar)
		}
		if f := args.Prop("async"); bas.IsCallable(f) {
			go func(e *bas.Env) {
				code, hdr, buf, jar := send(e, false)
				e.Call(f.Object(), code, hdr, buf, jar)
			}(bas.EnvForAsyncCall(e))
			return
		}
		e.A = bas.Array(send(e, true))
	}).Object()
	for _, m := range []string{"get", "post", "put", "delete", "head", "patch"} {
		httpLib = httpLib.SetMethod(m, func(e *bas.Env) {
			ex := e.Get(1).Maybe().Object(nil)
			e.A = e.Call(e.Object(-1), bas.NewObject(0).SetProp("method", bas.Str(m)).SetProp("url", e.Get(0)).Merge(ex).ToValue())
		})
	}
	bas.Globals.SetProp("http", httpLib.ToValue())

	bas.Globals.SetMethod("buffer", func(e *bas.Env) {
		b := &internal.LimitedBuffer{Limit: e.Get(1).Maybe().Int(0)}
		b.Write(bas.ToReadonlyBytes(e.Get(0)))
		e.A = bas.NamedObject("Buffer", 0).
			SetPrototype(bas.Proto.ReadWriter).
			SetProp("_f", bas.ValueOf(b)).
			SetMethod("reset", func(e *bas.Env) { e.ThisProp("_f").(*internal.LimitedBuffer).Reset() }).
			SetMethod("value", func(e *bas.Env) { e.A = bas.UnsafeStr(e.ThisProp("_f").(*internal.LimitedBuffer).Bytes()) }).
			SetMethod("bytes", func(e *bas.Env) { e.A = bas.Bytes(e.ThisProp("_f").(*internal.LimitedBuffer).Bytes()) }).
			ToValue()
	})
}

func mathMinMax(e *bas.Env, max bool) {
	if v := e.Num(0); v.IsInt64() {
		vi := v.Int64()
		for ii := 1; ii < len(e.Stack()); ii++ {
			if x := e.Int64(ii); x >= vi == max {
				vi = x
			}
		}
		e.A = bas.Int64(vi)
	} else {
		vf := v.Float64()
		for i := 1; i < len(e.Stack()); i++ {
			if x := e.Float64(i); x >= vf == max {
				vf = x
			}
		}
		e.A = bas.Float64(vf)
	}
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

func ParseStrictJSON(s string) (bas.Value, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return bas.Nil, fmt.Errorf("empty value")
	}
	switch s[0] {
	case 'n':
		return bas.Nil, nil
	case 't', 'f':
		v, err := strconv.ParseBool(s)
		return bas.Bool(v), err
	case '[':
		a := []interface{}{}
		err := json.Unmarshal([]byte(s), &a)
		return parseJSON(a), err
	case '{':
		a := map[string]interface{}{}
		err := json.Unmarshal([]byte(s), &a)
		return parseJSON(a), err
	default:
		return bas.Nil, fmt.Errorf("invalid value")
	}
}

func parseJSON(v interface{}) bas.Value {
	switch v := v.(type) {
	case []interface{}:
		a := make([]bas.Value, len(v))
		for i := range a {
			a[i] = parseJSON(v[i])
		}
		return bas.Array(a...)
	case map[string]interface{}:
		a := bas.NewObject(len(v) / 2)
		for k, v := range v {
			a.SetProp(k, parseJSON(v))
		}
		return a.ToValue()
	}
	return bas.ValueOf(v)
}

func valueOrPanic(v interface{}, err error) bas.Value {
	internal.PanicErr(err)
	return bas.ValueOf(v)
}
