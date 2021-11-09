package main

import (
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

	"github.com/coyove/script"
	_ "github.com/coyove/script/lib"
)

var (
	output    = flag.String("o", "ret", "separated by comma: (none|compileonly|opcode|bytes|ret|timing)+")
	input     = flag.String("i", "f", "input source, 'f': file, '-': stdin, others: string")
	version   = flag.Bool("v", false, "print version and usage")
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
		http.HandleFunc("/", script.WebREPLHandler(nil, nil))
		log.Println("listen", *apiServer)
		http.ListenAndServe(*apiServer, nil)
		return
	}

	log.SetFlags(0)

	if *version {
		fmt.Println("script virtual machine v" + strconv.FormatInt(script.Version, 10) + " (" + runtime.GOOS + "/" + runtime.GOARCH + ")")
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

	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
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
				b.EmergStop()
				log.Fatalln("timeout")
			}
		}()
	}

	i, err := b.Call()
	finished = true
	if _ret {
		fmt.Print(i)
		fmt.Print(" ", err, "\n")
	}
}
