package script

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

//go:embed typ/index.html
var indexBytes []byte

func WebREPLHandler(opt *CompileOptions, cb func(*Program)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() { recover() }()

		c := getCode(r)
		if c == "" {
			w.Header().Add("Content-Type", "text/html")
			w.Write(indexBytes)
			return
		}

		start := time.Now()
		bufOut := &limitedWriter{limit: 32 * 1024}

		p, err := LoadString(c, opt)
		if err != nil {
			writeJSON(w, map[string]interface{}{"error": err.Error()})
			return
		}
		finished := false
		go func() {
			time.Sleep(time.Second * 2)
			if !finished {
				p.EmergStop()
			}
		}()
		p.MaxCallStackSize = 100
		p.Stdout = bufOut
		p.Stderr = bufOut
		if cb != nil {
			cb(p)
		}
		code := p.PrettyCode()
		v, err := p.Run()
		finished = true
		if err != nil {
			writeJSON(w, map[string]interface{}{
				"error":   err.Error(),
				"elapsed": time.Since(start).Seconds(),
				"stdout":  bufOut.String(),
				"opcode":  code,
			})
			return
		}
		writeJSON(w, map[string]interface{}{
			"elapsed": time.Since(start).Seconds(),
			"result":  fmt.Sprint(v),
			"stdout":  bufOut.String(),
			"opcode":  code,
		})
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

type limitedWriter struct {
	limit int
	bytes.Buffer
}

func (w *limitedWriter) Write(b []byte) (int, error) {
	if w.Len()+len(b) > w.limit {
		b = b[:w.limit-w.Len()]
	}
	return w.Buffer.Write(b)
}
