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
var playgroundCode = `
-- Author: coyove
_, author = re([[Author: (\S+)]]).find(Program.Source)
println("Author is:", author)

-- Print all global values
local g = debug.globals()

print("version %d, total global values: %d".format(VERSION, #g/3))

function pp(idx, f)
    if f == nil then return end
    if f is callable then
        print(idx, ": function ", f)
    else
        print(idx, ": ", json.stringify(f))
    end		
end

for i=0,#g,3 do
    pp(i//3, g[i + 2])
end`

func PlaygroundHandler(defaultCode string, opt *LoadOptions) func(w http.ResponseWriter, r *http.Request) {
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
			var add2 = func(f, n string, force bool) {
				if force || n[0] >= 'A' && n[0] <= 'Z' {
					add("(" + f + ")." + n)
				}
			}
			var addType func(reflect.Type)
			addType = func(rf reflect.Type) {
				rfs := rf.String()
				if dedup[rfs] {
					return
				}
				dedup[rfs] = true
				rff := rf
				if rf.Kind() == reflect.Ptr {
					rff = rff.Elem()
				}
				if rff.Kind() == reflect.Struct {
					s := rff.String()
					for i := 0; i < rff.NumField(); i++ {
						add2(s, rff.Field(i).Name, false)
						addType(rff.Field(i).Type)
					}
					for i := 0; i < rf.NumMethod(); i++ {
						add2(rfs, rf.Method(i).Name, false)
					}
					if rf != rff {
						for i := 0; i < rff.NumMethod(); i++ {
							add2(s, rff.Method(i).Name, false)
						}
					}
				}
			}
			x := bas.TopSymbols()
			if opt != nil {
				x.Merge(&opt.Globals)
			}
			x.Foreach(func(k bas.Value, v *bas.Value) bool {
				add(k.String())
				switch v.Type() {
				case typ.Object:
					v.Object().Foreach(func(kk bas.Value, vv *bas.Value) bool {
						add2(k.String(), kk.String(), true)
						return true
					})
				case typ.Native:
					addType(reflect.ValueOf(v.Interface()).Type())
				}
				return true
			})

			buf := bytes.Replace(playgroundHTML, []byte("__NAMES__"), []byte(strings.Join(names, ",")), -1)
			if defaultCode != "" {
				buf = bytes.Replace(buf, []byte("__CODE__"), []byte(defaultCode), -1)
			} else {
				buf = bytes.Replace(buf, []byte("__CODE__"), []byte(playgroundCode), -1)
			}
			w.Write(buf)
			return
		}

		start := time.Now()
		bufOut := &internal.LimitedBuffer{Limit: 32 * 1024}

		p, err := LoadString(c, opt)
		if err != nil {
			writeJSON(w, map[string]interface{}{"error": err.Error()})
			return
		}
		p.MaxStackSize = 1000
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
