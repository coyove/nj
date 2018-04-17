package main

import (
	"log"
	"runtime"
	"time"

	"github.com/coyove/eugine/base"
	"github.com/coyove/eugine/compiler"
	"github.com/coyove/eugine/vm"

	"net/http"
	_ "net/http/pprof"
)

func main() {

	go func() {
		http.ListenAndServe("0.0.0.0:8080", nil)
	}()

	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	start := time.Now()

	b, err := compiler.LoadFile("tests/fib.txt")
	log.Println(err, b)
	log.Println(base.NewBytesReader(b).Prettify(0))

	vm.Exec(base.NewEnv(nil), b)
	// f, _ := os.Create("1.pbm")
	// f.Write([]byte("P4\n 1600 1600\n"))

	// for i := range ln {
	// 	a := ln[i].Array()
	// 	for j := range a {
	// 		f.Write([]byte{byte(a[j].Number())})
	// 	}
	// }

	// f.Close()
	time.Sleep(time.Minute)
	log.Println(time.Now().Sub(start).Nanoseconds() / 1e6)
}
