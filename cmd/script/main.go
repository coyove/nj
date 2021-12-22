package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"

	"github.com/coyove/nj"
	"github.com/coyove/nj/bas"
	"github.com/coyove/nj/lib"
)

var (
	output    = flag.String("o", "ret", "separated by comma: (none|compileonly|opcode|bytes|ret|timing)+")
	input     = flag.String("i", "f", "input source, 'f': file, '-': stdin, others: string")
	version   = flag.Bool("v", false, "print version and usage")
	repl      = flag.Bool("repl", false, "repl mode")
	timeout   = flag.Int("t", 0, "max execution time in ms")
	stackSize = flag.Int("ss", 1e6, "max stack size (counted by 16 bytes)")
	apiServer = flag.String("serve", "", "start as language playground")
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
		http.HandleFunc("/", lib.PlaygroundHandler(nil))
		log.Println("listen", *apiServer)
		http.ListenAndServe(*apiServer, nil)
		return
	}

	if *repl {
		runRepl()
	}

	log.SetFlags(0)

	if *version {
		fmt.Println("nj virtual machine v" + strconv.FormatInt(bas.Version, 10) + " (" + runtime.GOOS + "/" + runtime.GOARCH + ")")
		flag.Usage()
		return
	}

	if p, _ := os.Getwd(); !strings.Contains(p, "/cmd/nj") {
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
			log.Fatalln("Please specify the input file: ./nj <filename>")
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

	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	start := time.Now()

	var b *bas.Program
	var err error

	defer func() {
		if _opcode {
			log.Println(b.GoString())
		}
		if _timing {
			log.Printf("Time elapsed: %v\n", time.Since(start))
		}
	}()

	if *input == "f" {
		b, err = nj.LoadFile(source, nil)
	} else {
		b, err = nj.LoadString(source, nil)
	}
	if err != nil {
		log.Fatalln(err)
	}
	b.MaxStackSize = int64(*stackSize)

	if _compileonly {
		return
	}

	var finished bool
	if *timeout > 0 {
		// b.SetTimeout(time.Second * time.Duration(*timeout))
		go func() {
			time.Sleep(time.Second * time.Duration(*timeout))
			if !finished {
				b.Stop()
				log.Fatalln("timeout")
			}
		}()
	}

	i, err := b.Run()
	finished = true
	if _ret {
		fmt.Print(i)
		fmt.Print(" ", err, "\n")
	}
}

func runRepl() {
	var code []string
	var globals *bas.Object
	rd := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		s, err := rd.ReadString('\n')
		if err != nil {
			fmt.Println("Exit")
			break
		}
		s = strings.TrimSuffix(s, "\n")
		code = append(code, s)
		if s == "" || strings.HasSuffix(s, ";") {
			text := strings.TrimSuffix(strings.Join(code, "\n"), ";")
			p, err := nj.LoadString(text, &bas.Environment{Globals: globals})
			if err != nil {
				fmt.Println("x", err)
			} else {
				res, err := p.Run()
				if err != nil {
					fmt.Println("x", err)
				} else {
					fmt.Println("=>", res)
				}
				globals = p.LocalsObject()
			}
			code = code[:0]
		}
	}
}
