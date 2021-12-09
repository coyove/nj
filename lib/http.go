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

	"github.com/coyove/nj"
	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

var HostWhitelist = map[string][]string{}

func init() {
	nj.Globals.SetProp("url", nj.Func("url", nil, "").Object().
		SetMethod("escape", func(e *nj.Env) {
			e.A = nj.Str(url.QueryEscape(e.Str(0)))
		}, "").
		SetMethod("unescape", func(e *nj.Env) {
			v, err := url.QueryUnescape(e.Str(0))
			internal.PanicErr(err)
			e.A = nj.Str(v)
		}, "").
		ToValue())
	nj.Globals.SetProp("http", nj.Func("http", func(env *nj.Env) {
		args := env.Get(0).Object()
		to := args.Prop("timeout").Safe().Float64(1 << 30)

		method := strings.ToUpper(args.Get(nj.Str("method")).Safe().Str("GET"))

		u, err := url.Parse(args.Get(nj.Str("url")).Safe().Str("bad://%url%"))
		internal.PanicErr(err)

		addKV := func(k string, add func(k, v string)) {
			x := args.Get(nj.Str(k))
			if x.Type() == typ.Object {
				x.Object().Foreach(func(k nj.Value, v *nj.Value) bool { add(k.String(), v.String()); return true })
			}
		}

		if len(HostWhitelist) > 0 {
			ok := false
			for _, allow := range HostWhitelist[u.Host] {
				ok = ok || strings.EqualFold(allow, method)
			}
			if !ok {
				internal.PanicErr(fmt.Errorf("%s %v not allowed", method, u))
			}
		}

		additionalQueries := u.Query()
		addKV("query", additionalQueries.Add) // append queries to url
		u.RawQuery = additionalQueries.Encode()

		var bodyReader io.Reader
		dataFrom, urlForm, jsonForm := (*multipart.Writer)(nil), false, false

		if j := args.Prop("json"); j != nj.Nil {
			bodyReader = strings.NewReader(j.JSONString())
			jsonForm = true
		} else {
			var form url.Values
			addKV("form", form.Add) // check "form"
			urlForm = len(form) > 0
			if urlForm {
				bodyReader = strings.NewReader(form.Encode())
			} else if rd := args.Prop("data"); rd != nj.Nil {
				bodyReader = nj.NewReader(rd)
			}
		}

		if bodyReader == nil && method == "POST" {
			// Check form-data
			payload := bytes.Buffer{}
			writer := multipart.NewWriter(&payload)
			if x := args.Prop("multipart"); x.Type() == typ.Object {
				x.Object().Foreach(func(k nj.Value, v *nj.Value) bool {
					key := k.String()
					filename := ""
					if strings.Contains(key, "/") {
						filename = key[strings.Index(key, "/")+1:]
						key = key[:strings.Index(key, "/")]
					}
					if filename != "" {
						part, err := writer.CreateFormFile(key, filename)
						internal.PanicErr(err)
						_, err = io.Copy(part, nj.NewReader(*v))
						internal.PanicErr(err)
					} else {
						part, err := writer.CreateFormField(key)
						internal.PanicErr(err)
						_, err = io.Copy(part, nj.NewReader(*v))
						internal.PanicErr(err)
					}
					return true
				})
			}
			internal.PanicErr(writer.Close())
			if payload.Len() > 0 {
				bodyReader = &payload
				dataFrom = writer
			}
		}

		req, err := http.NewRequest(method, u.String(), bodyReader)
		internal.PanicErr(err)

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
		if v := args.Get(nj.Str("jar")); v.Type() == typ.Native {
			client.Jar, _ = v.Interface().(http.CookieJar)
		}
		if !args.Get(nj.Str("noredirect")).IsFalse() {
			client.CheckRedirect = func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			}
		}
		if p := args.Get(nj.Str("proxy")).Safe().Str(""); p != "" {
			client.Transport = &http.Transport{
				Proxy: func(r *http.Request) (*url.URL, error) { return url.Parse(p) },
			}
		}

		// Send
		resp, err := client.Do(req)
		internal.PanicErr(err)

		var buf nj.Value
		if args.Prop("bodyreader").IsFalse() && args.Prop("br").IsFalse() {
			resp.Body.Close()
		} else {
			buf = nj.NewObject(1).SetProp("_f", nj.ValueOf(resp.Body)).SetPrototype(nj.ReadCloserProto).ToValue()
		}

		hdr := map[string]string{}
		for k := range resp.Header {
			hdr[k] = resp.Header.Get(k)
		}
		env.A = nj.NewArray(nj.Int(resp.StatusCode), nj.ValueOf(hdr), buf, nj.ValueOf(client.Jar)).ToValue()
	}, "$f(options: object) -> array\n"+
		"\tperform an HTTP request and return [code, headers, body_reader, cookie_jar]\n"+
		"\t'url' is a mandatory parameter in `options`, others are optional and pretty self explanatory:\n"+
		"\thttp({url='...'})\n"+
		"\thttp({url='...', noredirect=true})\n"+
		"\thttp({url='...', bodyreader=true})\n"+
		"\thttp({method='POST', url='...'})\n"+
		"\thttp({method='POST', url='...'}, json={...})\n"+
		"\thttp({method='POST', url='...', query={key=value}})\n"+
		"\thttp({method='POST', url='...', header={key=value}, form={key=value}})\n"+
		"\thttp({method='POST', url='...', multipart={file={reader}}})\n"+
		"\thttp({method='POST', url='...', proxy='http://127.0.0.1:8080'})",
	))
}
