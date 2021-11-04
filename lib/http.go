package lib

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/coyove/script"
	"github.com/coyove/script/typ"
)

var HostWhitelist = map[string][]string{}

func init() {
	script.AddGlobalValue("url", script.Map(
		script.Str("escape"), script.Native1("escape", func(env *script.Env, a script.Value) script.Value {
			return script.Str(url.QueryEscape(a.MustStr("")))
		}),
		script.Str("unescape"), script.Native1("unescape", func(env *script.Env, a script.Value) script.Value {
			v, err := url.QueryUnescape(a.MustStr(""))
			panicErr(err)
			return script.Str(v)
		}),
	))
	script.AddGlobalValue("http", script.Native("http", func(env *script.Env) {
		args := env.Get(0).Table()

		ctx, cancel, _ := env.Deadline()
		defer func() {
			cancel()
			if r := recover(); r != nil {
				env.A = script.Array(script.Val(r))
			}
		}()

		method := strings.ToUpper(args.Get(script.Str("method")).StringDefault("GET"))

		u, err := url.Parse(args.Get(script.Str("url")).StringDefault("bad://%url%"))
		panicErr(err)

		addKV := func(k string, add func(k, v string)) {
			x := args.Get(script.Str(k))
			if x.Type() == typ.Table {
				p := x.Table()
				for k, v := p.Next(script.Nil); k != script.Nil; k, v = p.Next(k) {
					add(k.String(), v.String())
				}
			}
		}

		if len(HostWhitelist) > 0 {
			ok := false
			for _, allow := range HostWhitelist[u.Host] {
				ok = ok || strings.EqualFold(allow, method)
			}
			if !ok {
				panicErr(fmt.Errorf("%s %v not allowed", method, u))
			}
		}

		additionalQueries := u.Query()
		addKV("query", additionalQueries.Add) // append queries to url
		u.RawQuery = additionalQueries.Encode()

		body := args.Get(script.Str("rawbody")).StringDefault("")
		dataFrom, urlForm, jsonForm := (*multipart.Writer)(nil), false, false
		if body == "" {
			form := url.Values{}
			addKV("form", form.Add) // check "form"
			body = form.Encode()
			urlForm = len(form) > 0
		}
		if body == "" {
			// Check "json"
			if j := args.GetString("json"); j != script.Nil {
				body = j.JSONString()
				jsonForm = true
			}
		}

		var bodyReader io.Reader = strings.NewReader(body)

		if body == "" && method == "POST" {
			// Check form-data
			payload := bytes.Buffer{}
			writer := multipart.NewWriter(&payload)
			addKV("multipart", func(k, v string) {
				if strings.HasPrefix(v, "@") {
					path := v[1:]
					buf := panicErr2(ioutil.ReadFile(path)).([]byte)
					part := panicErr2(writer.CreateFormFile(k, filepath.Base(path))).(io.Writer)
					panicErr2(part.Write(buf))
				} else {
					part := panicErr2(writer.CreateFormField(k)).(io.Writer)
					panicErr2(io.WriteString(part, v))
				}
			})
			panicErr(writer.Close())
			if payload.Len() > 0 {
				bodyReader = &payload
				dataFrom = writer
			}
		}

		req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
		panicErr(err)

		switch {
		case urlForm:
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		case jsonForm:
			req.Header.Add("Content-Type", "application/json")
		case dataFrom != nil:
			req.Header.Add("Content-Type", dataFrom.FormDataContentType())
		}

		addKV("header", req.Header.Add) // append headers

		// Construct HTTP client
		client := &http.Client{}
		if to := args.Get(script.Str("timeout")).IntDefault(0); to > 0 {
			client.Timeout = time.Duration(to) * time.Millisecond
		}
		if v := args.Get(script.Str("jar")); v.Type() == typ.Interface {
			client.Jar, _ = v.Interface().(http.CookieJar)
		}
		if !args.Get(script.Str("no_redirect")).IsFalse() {
			client.CheckRedirect = func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			}
		}
		if p := args.Get(script.Str("proxy")).StringDefault(""); p != "" {
			client.Transport = &http.Transport{
				Proxy: func(r *http.Request) (*url.URL, error) { return url.Parse(p) },
			}
		}

		// Send
		resp, err := client.Do(req)
		panicErr(err)

		var buf script.Value
		if args.GetString("body_reader").IsFalse() && args.GetString("br").IsFalse() {
			resp.Body.Close()
		} else {
			buf = script.TableProto(script.ReadCloser, script.Str("_f"), script.Val(resp.Body))
		}

		hdr := map[string]string{}
		for k := range resp.Header {
			hdr[k] = resp.Header.Get(k)
		}
		env.A = script.Array(script.Int(int64(resp.StatusCode)), script.Val(hdr), buf, script.Val(client.Jar))
	}, "http(options: table) array",
		"\tperforma an HTTP request and return { code, headers, body_reader, cookie_jar }",
		"\t'url' is a mandatory parameter in options, others are optional and pretty self explanatory:",
		"\thttp({url='...'})",
		"\thttp({url='...', no_redirect=true})",
		"\thttp({url='...', body_reader=true})",
		"\thttp({method='POST', url='...'})",
		"\thttp({method='POST', url='...'}, json={...})",
		"\thttp({method='POST', url='...', query={key=value}})",
		"\thttp({method='POST', url='...', header={key=value}, form={key=value}})",
		"\thttp({method='POST', url='...', multipart={file='@path/to/file'}})",
		"\thttp({method='POST', url='...', proxy='http://127.0.0.1:8080'})",
	))
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func panicErr2(v interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return v
}
