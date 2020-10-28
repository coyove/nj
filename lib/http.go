package lib

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/coyove/script"
)

var HostWhitelist = map[string][]string{}

func init() {
	script.AddGlobalValue("http", script.NativeWithParamMap(func(env *script.Env, args script.Arguments) {
		ctx, cancel, _ := env.Deadline()
		defer cancel()

		method := args.GetString("method", "GET")

		u, err := url.Parse(args.GetString("url", "bad://%url%"))
		if err != nil {
			env.Return(script.Interface(err))
			return
		}

		if len(HostWhitelist) > 0 {
			ok := false
			for _, allow := range HostWhitelist[u.Host] {
				if allow == method {
					ok = true
				}
			}
			if !ok {
				env.Return(script.Interface(fmt.Errorf("%s %v not allowed", method, u)))
				return
			}
		}

		{ // Append queries to url
			q := u.Query()
			iterStrings(args["query"], func(line string) {
				if k, v, ok := splitKV(line); ok {
					q.Add(k, v)
				}
			})
			u.RawQuery = q.Encode()
		}

		body := args.GetString("rawbody", "")
		urlForm, jsonForm := false, false
		if body == "" {
			// Check "form"
			form := url.Values{}
			iterStrings(args["form"], func(line string) {
				if k, v, ok := splitKV(line); ok {
					form.Add(k, v)
				}
			})
			body = form.Encode()
			urlForm = len(form) > 0
		}
		if body == "" {
			// Check "json"
			body = args.GetString("json", "")
			jsonForm = len(body) > 0
		}

		req, err := http.NewRequestWithContext(ctx, method, u.String(), strings.NewReader(body))
		if err != nil {
			env.Return(script.Interface(err))
			return
		}

		switch {
		case urlForm:
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		case jsonForm:
			req.Header.Add("Content-Type", "application/json")
		}

		iterStrings(args["header"], func(line string) {
			if k, v, ok := splitKV(line); ok {
				req.Header.Add(k, v)
			}
		})

		client := &http.Client{}
		if to := args.GetInt("timeout", 0); to > 0 {
			client.Timeout = time.Duration(to) * time.Millisecond
		}

		resp, err := client.Do(req)
		if err != nil {
			env.Return(script.Interface(err))
			return
		}

		defer resp.Body.Close()

		var r io.Reader = resp.Body
		if env.Global.MaxStackSize > 0 {
			r = io.LimitReader(r, env.Global.MaxStackSize*16)
		}
		buf, _ := ioutil.ReadAll(r)

		hdr := map[string]string{}
		for k := range resp.Header {
			hdr[k] = resp.Header.Get(k)
		}
		env.Return(
			script.Int(int64(resp.StatusCode)),
			script.Interface(hdr),
			env.NewStringBytes(buf),
		)
	}, "method", "url", "rawbody", "header", "query", "timeout", "form", "json"))
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
	case script.VStack:
		for _, line := range v.Stack() {
			f(line.String())
		}
	}
}
