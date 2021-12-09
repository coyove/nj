package lib

import (
	"sync"

	"github.com/coyove/nj"
)

func init() {
	nj.Globals.SetProp("sync", nj.NamedObject("sync", 0).
		SetMethod("mutex", func(e *nj.Env) { e.A = nj.ValueOf(&sync.Mutex{}) }, "$f() -> *go.sync.Mutex").
		SetMethod("rwmutex", func(e *nj.Env) { e.A = nj.ValueOf(&sync.RWMutex{}) }, "$f() -> *go.sync.RWMutex").
		SetMethod("waitgroup", func(e *nj.Env) { e.A = nj.ValueOf(&sync.WaitGroup{}) }, "$f() -> *go.sync.WaitGroup").
		SetPrototype(nj.StaticObjectProto).
		ToValue())
}
