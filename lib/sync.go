package lib

import (
	"sync"

	"github.com/coyove/nj/bas"
)

func init() {
	bas.Globals.SetProp("sync", bas.NamedObject("sync", 0).
		SetMethod("mutex", func(e *bas.Env) { e.A = bas.ValueOf(&sync.Mutex{}) }, "$f() -> *go.sync.Mutex").
		SetMethod("rwmutex", func(e *bas.Env) { e.A = bas.ValueOf(&sync.RWMutex{}) }, "$f() -> *go.sync.RWMutex").
		SetMethod("waitgroup", func(e *bas.Env) { e.A = bas.ValueOf(&sync.WaitGroup{}) }, "$f() -> *go.sync.WaitGroup").
		SetPrototype(bas.StaticObjectProto).
		ToValue())
}
