package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"github.com/coyove/potatolang"
)

const VERSION = "0.1.0"

var goroutinePerCPU = flag.Int("goroutine", 2, "goroutines per CPU")
var output = flag.String("o", "none", "separated by comma: [none, compileonly, compiledsize, opcode, bytes, ret, timing]+")
var input = flag.String("i", "f", "input source, 'f': file, '-': stdin, others: string")
var version = flag.Bool("v", false, "print version and usage")
var quiet = flag.Bool("quieterr", false, "suppress the error output (if any)")
var timeout = flag.Int("t", 0, "max execution time in ms")

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
	log.SetFlags(0)

	if *version {
		fmt.Println("\"pol\": potatolang virtual machine v" + VERSION + " (" + runtime.GOOS + "/" + runtime.GOARCH + ")")
		flag.Usage()
		return
	}

	switch *input {
	case "f":
		if source == "" {
			log.Fatalln("Please specify the input file: ./pol <filename>")
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

	var _opcode, _timing, _bytes, _ret, _compileonly, _roughsize bool

ARG:
	for _, a := range strings.Split(*output, ",") {
		switch a {
		case "n", "no", "none":
			_opcode, _ret, _bytes, _timing, _compileonly, _roughsize =
				false, false, false, false, false, false
			break ARG
		case "o", "opcode", "op":
			_opcode = true
		case "r", "ret", "return":
			_ret = true
		case "b", "byte", "bytes":
			_bytes = true
		case "t", "time", "timing":
			_timing = true
		case "co", "compile", "compileonly":
			_compileonly = true
		case "cs", "compiledsize", "rs", "roughsize":
			_roughsize = true
		}
	}

	runtime.GOMAXPROCS(runtime.NumCPU() * *goroutinePerCPU)
	start := time.Now()

	var b *potatolang.Closure
	var err error

	defer func() {
		if *quiet {
			recover()
		}

		if _opcode {
			log.Println(b.PrettyString())
		}
		if _bytes {
			os.Stderr.Write(b.BytesCode())
		}
		if _roughsize {
			ln := len(b.BytesCode())
			ln += len(b.Consts()) * 16
			ln += len(b.Pos()) * 8

			// 1.1: a factor
			log.Printf("Compiled size: ~%.1fK with %d opcode\n", float64(ln)/1024*1.1, len(b.Code()))
		}
		if _timing {
			e := float64(time.Now().Sub(start).Nanoseconds()) / 1e6
			if e < 1000 {
				log.Printf("Time elapsed: %.1fms\n", e)
			} else {
				log.Printf("Time elapsed: %.3fs\n", e/1e3)
			}
		}
	}()

	if *input == "f" {
		b, err = potatolang.LoadFile(source)
	} else {
		b, err = potatolang.LoadString(source)
	}
	if err != nil {
		log.Fatalln(err)
	}

	if _compileonly {
		return
	}

	var exit *uintptr
	var ok = make(chan bool, 1)
	if *timeout > 0 {
		exit = b.MakeCancelable()
		go func() {
			select {
			case <-time.After(time.Duration(*timeout) * time.Millisecond):
				atomic.StoreUintptr(exit, 1)
				log.Println("Timeout:", *timeout, "ms")
			case <-ok:
				// peacefully exit
			}
		}()
	}

	i := b.Exec(nil)
	ok <- true
	if _ret {
		fmt.Println(i.ToPrintString())
	}
}
