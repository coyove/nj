package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/coyove/script"
	"github.com/coyove/script/lib"
)

const VERSION = "0.2.0"

var (
	goroutinePerCPU = flag.Int("goroutine", 2, "goroutines per CPU")
	output          = flag.String("o", "none", "separated by comma: (none|compileonly|opcode|bytes|ret|timing)+")
	input           = flag.String("i", "f", "input source, 'f': file, '-': stdin, others: string")
	version         = flag.Bool("v", false, "print version and usage")
	timeout         = flag.Int("t", 0, "max execution time in ms")
	apiServer       = flag.String("serve", "", "start as language playground")
)

func main() {
	source := ""
	for i, arg := range os.Args {
		if _, err := os.Stat(arg); err == nil && i > 0 {
			source = arg
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
			break
		}
	}

	flag.Parse()

	if *apiServer != "" {
		script.RemoveGlobalValue("sleep")
		script.RemoveGlobalValue("scanln")
		lib.HostWhitelist["httpbin.org"] = []string{"DELETE", "GET", "PATCH", "POST", "PUT"}
		lib.HostWhitelist["example.com"] = []string{"DELETE", "GET", "PATCH", "POST", "PUT"}
		lib.HostWhitelist["bokete.jp"] = []string{"DELETE", "GET", "PATCH", "POST", "PUT"}
		lib.HostWhitelist["cn.bing.com"] = []string{"DELETE", "GET", "PATCH", "POST", "PUT"}

		apiServerStatic := func() string {
			buf, _ := exec.Command("go", "env", "GOPATH").Output()
			return filepath.Join(strings.TrimSpace(string(buf)), "/src/github.com/coyove/script/cmd/script/docs")
		}()
		http.Handle("/", http.FileServer(http.Dir(apiServerStatic)))
		http.HandleFunc("/share", func(w http.ResponseWriter, r *http.Request) {
			defer func() { recover() }()
			read := func(resp *http.Response, err error) ([]byte, error) {
				if err != nil {
					return nil, err
				}
				buf, _ := ioutil.ReadAll(resp.Body)
				resp.Body.Close()
				return bytes.TrimSpace(buf), nil
			}
			var buf []byte
			var err error
			if src := r.URL.Query().Get("get"); src != "" {
				buf, err = read(http.Get("http://sprunge.us/" + src))
			} else {
				buf, err = read(http.Post("http://sprunge.us", "application/x-www-form-urlencoded",
					strings.NewReader("sprunge="+url.QueryEscape(getCode(r)))))
			}
			if err != nil {
				writeJSON(w, map[string]interface{}{"error": err.Error()})
			} else {
				writeJSON(w, map[string]interface{}{"data": string(buf)})
			}
		})
		http.HandleFunc("/eval", func(w http.ResponseWriter, r *http.Request) {
			defer func() { recover() }()

			start := time.Now()
			c := getCode(r)

			p, err := script.LoadString(c, nil)
			if err != nil {
				writeJSON(w, map[string]interface{}{"error": err.Error()})
				return
			}
			bufOut := &limitedWriter{limit: 16 * 1024}
			p.SetTimeout(time.Second * 2)
			p.MaxCallStackSize = 100
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
				"result":  v.Any(),
				"stdout":  bufOut.String(),
				"opcode":  code,
			})
		})
		log.Println("listen", *apiServer)
		http.ListenAndServe(*apiServer, nil)
		return
	}

	log.SetFlags(0)

	if *version {
		fmt.Println("\"script\": script virtual machine v" + VERSION + " (" + runtime.GOOS + "/" + runtime.GOARCH + ")")
		flag.Usage()
		return
	}

	if p, _ := os.Getwd(); !strings.Contains(p, "/cmd/script") {
		f, err := os.Create("cpuprofile")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	switch *input {
	case "f":
		if source == "" {
			log.Fatalln("Please specify the input file: ./script <filename>")
		}
	case "-":
		buf, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalln(err)
		}
		source = string(buf)
	default:
		if _, err := os.Stat(*input); err == nil {
			source = *input
			*input = "f"
		} else {
			source = *input
		}
	}

	var _opcode, _timing, _ret, _compileonly bool
	for _, a := range strings.Split(*output, ",") {
		switch a {
		case "o", "opcode", "op":
			_opcode = true
		case "r", "ret", "return":
			_ret = true
		case "t", "time", "timing":
			_timing = true
		case "co", "compile", "compileonly":
			_compileonly = true
		}
	}

	runtime.GOMAXPROCS(runtime.NumCPU() * *goroutinePerCPU)
	start := time.Now()

	var b *script.Program
	var err error

	defer func() {
		if _opcode {
			log.Println(b.PrettyCode())
		}
		if _timing {
			log.Printf("Time elapsed: %v\n", time.Since(start))
		}
	}()

	if *input == "f" {
		b, err = script.LoadFile(source, nil)
	} else {
		b, err = script.LoadString(source, nil)
	}
	if err != nil {
		log.Fatalln(err)
	}

	if _compileonly {
		return
	}
	if *timeout > 0 {
		b.SetTimeout(time.Second * time.Duration(*timeout))
	}

	i, err := b.Call()
	if _ret {
		fmt.Print(i)
		fmt.Print(" ", err, "\n")
	}
}

func writeJSON(w http.ResponseWriter, m map[string]interface{}) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
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
