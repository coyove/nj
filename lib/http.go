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
)

var HostWhitelist = map[string][]string{}

func init() {
	script.AddGlobalValue("http", script.NativeWithParamMap("http", func(env *script.Env) {
		args := env.A.Array()

		ctx, cancel, _ := env.Deadline()
		defer func() {
			cancel()
			if r := recover(); r != nil {
				env.A = script.Array(script.Interface(r))
			}
		}()

		method := strings.ToUpper(args.Get(script.String("method")).StringDefault("GET"))

		u, err := url.Parse(args.Get(script.String("url")).StringDefault("bad://%url%"))
		panicErr(err)

		addKV := func(k string, add func(k, v string)) {
			x := args.Get(script.String(k))
			if x.Type() == script.VArray {
				p := x.Array()
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

		body := args.Get(script.String("rawbody")).StringDefault("")
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
		if to := args.Get(script.String("timeout")).IntDefault(0); to > 0 {
			client.Timeout = time.Duration(to) * time.Millisecond
		}
		if v := args.Get(script.String("jar")); v.Type() == script.VInterface {
			client.Jar, _ = v.Interface().(http.CookieJar)
		}
		if !args.Get(script.String("no_redirect")).IsFalse() {
			client.CheckRedirect = func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			}
		}
		if p := args.Get(script.String("proxy")).StringDefault(""); p != "" {
			client.Transport = &http.Transport{
				Proxy: func(r *http.Request) (*url.URL, error) { return url.Parse(p) },
			}
		}

		// Send
		resp, err := client.Do(req)
		panicErr(err)

		defer resp.Body.Close()

		var r io.Reader = resp.Body
		buf := panicErr2(ioutil.ReadAll(r)).([]byte)

		hdr := map[string]string{}
		for k := range resp.Header {
			hdr[k] = resp.Header.Get(k)
		}
		env.A = script.Array(
			script.Int(int64(resp.StatusCode)),
			script.Interface(hdr),
			script.Bytes(buf),
			script.Interface(client.Jar),
		)
	}, `http($a...a$) => code, body, headers, cookie_jar
    'url' is a mandatory parameter, others are optional and pretty self explanatory:
	http(url="...") -- GET req
	http(url="...", no_redirect=true)
	http("POST", "...")
	http("POST", "...", form={key=value})
    http("POST", "...", multipart={file='@path/to/file'})`,

		"method", "url", "rawbody", "timeout", "proxy",
		"header", "query", "form", "multipart",
		"json", "jar", "no_redirect"))
}

func splitKV(line string) (k string, v string, ok bool) {
	if idx := strings.Index(line, ":"); idx > -1 {
		k, v, ok = strings.TrimSpace(line[:idx]), strings.TrimSpace(line[idx+1:]), true
	}
	if idx := strings.Index(line, "="); idx > -1 {
		k, v, ok = strings.TrimSpace(line[:idx]), strings.TrimSpace(line[idx+1:]), true
	}
	if idx := strings.Index(line, " "); idx > -1 {
		k, v, ok = strings.TrimSpace(line[:idx]), strings.TrimSpace(line[idx+1:]), true
	}
	if idx := strings.Index(line, ","); idx > -1 {
		k, v, ok = strings.TrimSpace(line[:idx]), strings.TrimSpace(line[idx+1:]), true
	}
	tmp, err := url.QueryUnescape(v)
	if err == nil {
		v = tmp
	}
	return
}

func iterStrings(v script.Value, f func(string)) {
	switch v.Type() {
	case script.VString:
		f(v.String())
	case script.VArray:
		for _, line := range v.Array().Array() {
			f(line.String())
		}
	}
}

func iterStringPairs(v1, v2 script.Value, f func(string, string)) {
	switch v1.Type() + v2.Type() {
	case script.VString * 2:
		f(v1.String(), v2.String())
	case script.VArray * 2:
		for i, line := range v1.Array().Array() {
			if i < len(v2.Array().Array()) {
				f(line.String(), v2.Array().Array()[i].String())
			}
		}
	}
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
