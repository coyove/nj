package main

import (
	"log"
	"net/http"
	// _ "net/http/pprof"
	"runtime"
	"time"

	"github.com/coyove/bracket/base"
	"github.com/coyove/bracket/compiler"
)

func main() {

	// go func() {
	// 	http.ListenAndServe("0.0.0.0:8080", nil)
	// }()

	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	start := time.Now()

	base.CoreLibNames = append(base.CoreLibNames, "http")
	m := new(base.Tree)
	m.Put("handle", base.NewNativeClosureValue(2, func(env *base.Env) base.Value {
		pattern := env.Get(0).AsString()
		http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			q := r.FormValue("q")
			newEnv := base.NewEnv(nil)
			newEnv.Push(base.NewStringValue(q))
			w.Write(env.Get(1).AsClosure().Exec(newEnv).AsBytes())
		})
		return base.NewValue()
	}))
	m.Put("listen", base.NewNativeClosureValue(1, func(env *base.Env) base.Value {
		http.ListenAndServe(env.Get(0).AsString(), nil)
		return base.NewValue()
	}))
	base.CoreLibs["http"] = (base.NewMapValue(m))

	b, err := compiler.LoadFile("tests/builtin.txt")
	log.Println(err)
	log.Println(base.Prettify(b))

	e := base.NewEnv(nil)
	for _, name := range base.CoreLibNames {
		e.Push(base.CoreLibs[name])
	}

	i := base.Exec(e, b)
	log.Println(i.I())
	log.Println(time.Now().Sub(start).Nanoseconds() / 1e6)
}
