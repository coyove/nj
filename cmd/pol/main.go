package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/coyove/potatolang"
)

var recordTime = flag.Bool("timing", false, "record the execution time")
var goroutinePerCPU = flag.Int("goroutine", 2, "goroutines per CPU")
var output = flag.String("o", "none", "output [none, opcode, bytes, ret]+")

func main() {
	source := ""
	for i, arg := range os.Args {
		if _, err := os.Stat(arg); err == nil && i > 0 {
			source = arg
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
		}
	}
	if source == "" {
		log.Fatalln("Please specify the input file")
	}

	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU() * *goroutinePerCPU)
	start := time.Now()
	defer func() {
		if *recordTime {
			log.Println("Total execution time:", time.Now().Sub(start).Nanoseconds()/1e6, "ms")
		}
	}()

	b, err := potatolang.LoadFile(source)
	if err != nil {
		log.Fatalln(err)
	}

	_opcode := false
	_bytes := false
	_ret := false

ARG:
	for _, a := range strings.Split(*output, ",") {
		switch a {
		case "n", "no", "none":
			_opcode, _ret, _bytes = false, false, false
			break ARG
		case "o", "opcode", "op":
			_opcode = true
		case "r", "ret", "return":
			_ret = true
		case "b", "byte", "bytes":
			_bytes = true
		}
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
