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

func WebREPLHandler(loader func(string) (*Program, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() { recover() }()

		c := getCode(r)
		if c == "" {
			w.Header().Add("Content-Type", "text/html")
			w.Write(indexBytes)
			return
		}

		start := time.Now()
		bufOut := &limitedWriter{limit: 16 * 1024}

		var p *Program
		var err error
		if loader == nil {
			p, err = LoadString(c, nil)
			if err != nil {
				writeJSON(w, map[string]interface{}{"error": err.Error()})
				return
			}
			p.SetTimeout(time.Second * 2)
			p.MaxCallStackSize = 100
		} else {
			p, err = loader(c)
			if err != nil {
				writeJSON(w, map[string]interface{}{"error": err.Error()})
				return
			}
		}
		p.Stdout = bufOut
		p.Stderr = bufOut
		code := p.PrettyCode()
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
		writeJSON(w, map[string]interface{}{
			"elapsed": time.Since(start).Seconds(),
			"result":  v.Interface(),
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
	if w.Len() > w.limit {
		return 0, fmt.Errorf("overflow")
	}
	return w.Buffer.Write(b)
}
