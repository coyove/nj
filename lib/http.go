package lib

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/coyove/script"
	"github.com/coyove/script/typ"
)

var HostWhitelist = map[string][]string{}

func init() {
	script.AddGlobalValue("url", script.Map(
		script.Str("escape"), script.Func1("escape", func(a script.Value) script.Value {
			return script.Str(url.QueryEscape(a.MustStr("")))
		}),
		script.Str("unescape"), script.Func1("unescape", func(a script.Value) script.Value {
			v, err := url.QueryUnescape(a.MustStr(""))
			panicErr(err)
			return script.Str(v)
		}),
	))
	script.AddGlobalValue("http", script.Function("http", func(env *script.Env) {
		args := env.Get(0).Table()
		to := args.GetString("timeout").MaybeFloat(1 << 30)

		method := strings.ToUpper(args.Get(script.Str("method")).MaybeStr("GET"))

		u, err := url.Parse(args.Get(script.Str("url")).MaybeStr("bad://%url%"))
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

		var bodyReader io.Reader
		dataFrom, urlForm, jsonForm := (*multipart.Writer)(nil), false, false

		if j := args.GetString("json"); j != script.Nil {
			bodyReader = strings.NewReader(j.JSONString())
			jsonForm = true
		} else {
			var form url.Values
			addKV("form", form.Add) // check "form"
			urlForm = len(form) > 0
			if urlForm {
				bodyReader = strings.NewReader(form.Encode())
			} else if rd := args.GetString("data"); rd != script.Nil {
				bodyReader = script.NewReader(rd)
			}
		}

		if bodyReader == nil && method == "POST" {
			// Check form-data
			payload := bytes.Buffer{}
			writer := multipart.NewWriter(&payload)
			if x := args.GetString("multipart"); x.Type() == typ.Table {
				x.Table().Foreach(func(k, v script.Value) bool {
					key := k.MustStr("multipart key")
					filename := ""
					if strings.Contains(key, "/") {
						filename = key[strings.Index(key, "/")+1:]
						key = key[:strings.Index(key, "/")]
					}
					if filename != "" {
						part := panicErr2(writer.CreateFormFile(key, filename)).(io.Writer)
						panicErr2(io.Copy(part, script.NewReader(v)))
					} else {
						part := panicErr2(writer.CreateFormField(key)).(io.Writer)
						panicErr2(io.Copy(part, script.NewReader(v)))
					}
					return true
				})
			}
			panicErr(writer.Close())
			if payload.Len() > 0 {
				bodyReader = &payload
				dataFrom = writer
			}
		}

		req, err := http.NewRequest(method, u.String(), bodyReader)
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
		client.Timeout = time.Duration(to * float64(time.Second))
		if v := args.Get(script.Str("jar")); v.Type() == typ.Native {
			client.Jar, _ = v.Interface().(http.CookieJar)
		}
		if !args.Get(script.Str("noredirect")).IsFalse() {
			client.CheckRedirect = func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			}
		}
		if p := args.Get(script.Str("proxy")).MaybeStr(""); p != "" {
			client.Transport = &http.Transport{
				Proxy: func(r *http.Request) (*url.URL, error) { return url.Parse(p) },
			}
		}

		// Send
		resp, err := client.Do(req)
		panicErr(err)

		var buf script.Value
		if args.GetString("bodyreader").IsFalse() && args.GetString("br").IsFalse() {
			resp.Body.Close()
		} else {
			buf = script.TableProto(script.ReadCloserProto, script.Str("_f"), script.Val(resp.Body))
		}

		hdr := map[string]string{}
		for k := range resp.Header {
			hdr[k] = resp.Header.Get(k)
		}
		env.A = script.Array(script.Int(int64(resp.StatusCode)), script.Val(hdr), buf, script.Val(client.Jar))
	}, "http(options: table) array",
		"\tperform an HTTP request and return { code, headers, body_reader, cookie_jar }",
		"\t'url' is a mandatory parameter in options, others are optional and pretty self explanatory:",
		"\thttp({url='...'})",
		"\thttp({url='...', noredirect=true})",
		"\thttp({url='...', bodyreader=true})",
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
