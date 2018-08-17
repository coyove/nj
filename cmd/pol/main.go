package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/coyove/potatolang"
)

var goroutinePerCPU = flag.Int("goroutine", 2, "goroutines per CPU")
var output = flag.String("o", "none", "output [none, opcode, bytes, ret, timing]+")
var input = flag.String("i", "f", "input source, 'f': file, '-': stdin, others: string")
var version = flag.Bool("v", false, "print pol version")

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
		fmt.Println("\"pol\": potatolang virtual machine (" + runtime.GOOS + "/" + runtime.GOARCH + ")")
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

	var _opcode, _timing, _bytes, _ret bool

ARG:
	for _, a := range strings.Split(*output, ",") {
		switch a {
		case "n", "no", "none":
			_opcode, _ret, _bytes, _timing = false, false, false, false
			break ARG
		case "o", "opcode", "op":
			_opcode = true
		case "r", "ret", "return":
			_ret = true
		case "b", "byte", "bytes":
			_bytes = true
		case "t", "time", "timing":
			_timing = true
		}
	}

	runtime.GOMAXPROCS(runtime.NumCPU() * *goroutinePerCPU)
	start := time.Now()
	defer func() {
		if _timing {
			fmt.Println("Total execution time:", time.Now().Sub(start).Nanoseconds()/1e6, "ms")
		}
	}()

	var b *potatolang.Closure
	var err error

	if *input == "f" {
		b, err = potatolang.LoadFile(source)
	} else {
		b, err = potatolang.LoadString(source)
	}
	if err != nil {
		log.Fatalln(err)
	}

	i := b.Exec(nil)

	if _opcode {
		fmt.Println(b.PrettyString())
	}
	if _bytes {
		os.Stdout.Write(b.BytesCode())
	}
	if _ret {
		fmt.Println(i.I())
	}
}
