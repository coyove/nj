package main

import (
	"log"
	"runtime"
	"time"

	"github.com/coyove/bracket/base"
	"github.com/coyove/bracket/compiler"
	"github.com/coyove/bracket/vm"

	"net/http"
	_ "net/http/pprof"
)

func main() {

	go func() {
		http.ListenAndServe("0.0.0.0:8080", nil)
	}()

	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	start := time.Now()

	b, err := compiler.LoadFile("tests/test.txt")
	log.Println(err, b)
	log.Println(base.NewBytesReader(b).Prettify(0))

	log.Println(vm.Exec(base.NewEnv(nil), b).I())
	log.Println(time.Now().Sub(start).Nanoseconds() / 1e6)
}
