package nj

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/coyove/nj/bas"
	"github.com/coyove/nj/internal"
)

//go:embed playground.html
var playgroundHTML []byte

func PlaygroundHandler(opt *bas.Environment) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() { recover() }()

		c := getCode(r)
		if c == "" {
			w.Header().Add("Content-Type", "text/html")
			w.Write(playgroundHTML)
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
