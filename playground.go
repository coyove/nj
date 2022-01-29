package nj

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/coyove/nj/bas"
	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

//go:embed playground.html
var playgroundHTML []byte

func PlaygroundHandler(opt *bas.Environment) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() { recover() }()

		c := getCode(r)
		if c == "" {
			w.Header().Add("Content-Type", "text/html")
			var names []string
			var dedup = map[string]bool{}
			var add = func(n string) {
				if !dedup[n] {
					dedup[n], names = true, append(names, strconv.Quote(n))
				}
			}
			var add2 = func(f, n string) { add("(" + f + ")." + n) }
			x := bas.Globals.Copy(true)
			if opt != nil {
				x.Merge(opt.Globals)
			}
			x.Foreach(func(k bas.Value, v *bas.Value) bool {
				add(k.String())
				switch v.Type() {
				case typ.Object:
					v.Object().Foreach(func(kk bas.Value, vv *bas.Value) bool {
						add2(k.String(), kk.String())
						return true
					})
				case typ.Native:
					rv := reflect.ValueOf(v.Interface())
					rf := reflect.Indirect(rv).Type()
					if rf.Kind() == reflect.Struct {
						for i := 0; i < rf.NumField(); i++ {
							add2(rf.String(), rf.Field(i).Name)
						}
						rf := rv.Type()
						for i := 0; i < rf.NumMethod(); i++ {
							add2(rf.String(), rf.Method(i).Name)
						}
						if rv.Kind() == reflect.Ptr {
							rf := rv.Elem().Type()
							for i := 0; i < rf.NumMethod(); i++ {
								add2(rf.String(), rf.Method(i).Name)
							}
						}
					}
				}
				return true
			})
			w.Write(bytes.Replace(playgroundHTML, []byte("__NAMES__"), []byte(strings.Join(names, ",")), -1))
			return
		}

		start := time.Now()
		bufOut := &internal.LimitedBuffer{Limit: 32 * 1024}

		p, err := LoadString(c, opt)
		if err != nil {
			writeJSON(w, map[string]interface{}{"error": err.Error()})
			return
		}
		p.MaxStackSize = 100
		p.Deadline = start.Add(time.Second * 2)
		p.Stdout = bufOut
		p.Stderr = bufOut
		code := p.GoString()
		v, err := p.Run()
		if err != nil {
			writeJSON(w, map[string]interface{}{
				"error":   err.Error(),
				"elapsed": time.Since(start).Seconds(),
				"stdout":  bufOut.String(),
				"opcode":  code,
			})
			return
		}
		switch outf := r.URL.Query().Get("output"); outf {
		case "stdout":
			writeJSON(w, map[string]interface{}{"stdout": bufOut.String()})
		case "result":
			writeJSON(w, map[string]interface{}{"result": fmt.Sprint(v)})
		default:
			writeJSON(w, map[string]interface{}{
				"elapsed": time.Since(start).Seconds(),
				"result":  fmt.Sprint(v),
				"stdout":  bufOut.String(),
				"opcode":  code,
			})
		}
	}
}

func writeJSON(w http.ResponseWriter, m map[string]interface{}) {
	// w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	buf, _ := json.Marshal(m)
	w.Write(buf)
}

func getCode(r *http.Request) string {
	c := strings.TrimSpace(r.FormValue("code"))
	if c == "" {
		c = strings.TrimSpace(r.URL.Query().Get("code"))
	}
	if len(c) > 16*1024 {
		c = c[:16*1024]
	}
	return c
}
