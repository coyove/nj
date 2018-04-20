package vm

import (
	"sync"

	"github.com/coyove/bracket/base"
)

var lib_go = LibFunc{
	name: "go",
	args: 1,
	ff: func(env *base.Env) base.Value {
		newEnv := base.NewEnv(env.Parent())
		cls := env.Get(0).Closure()

		if cls.ArgsCount() < env.Size()-1 {
			panic("not enough arguments to start")
		}

		for i := 1; i < env.Size(); i++ {
			newEnv.Push(env.Get(int32(i)))
		}

		go Exec(newEnv, cls.Code())
		return base.NewValue()
	},
}

var lib_syncwaitgroupnew = LibFunc{
	name: "wait_group_new",
	args: 0,
	f: func(env *base.Env) base.Value {
		return base.NewGenericValue(&sync.WaitGroup{})
	},
}

var lib_syncwaitgroupadd = LibFunc{
	name: "wait_group_add",
	args: 2,
	f: func(env *base.Env) base.Value {
		wg := env.R0.Generic().(*sync.WaitGroup)
		wg.Add(int(env.R1.Number()))
		return base.NewValue()
	},
}

var lib_syncwaitgroupdone = LibFunc{
	name: "wait_group_done",
	args: 1,
	f: func(env *base.Env) base.Value {
		env.R0.Generic().(*sync.WaitGroup).Done()
		return base.NewValue()
	},
}

var lib_syncwaitgroupwait = LibFunc{
	name: "wait_group_wait",
	args: 1,
	f: func(env *base.Env) base.Value {
		wg := env.R0.Generic().(*sync.WaitGroup)
		wg.Wait()
		return base.NewValue()
	},
}

var lib_syncmutexnew = LibFunc{
	name: "mutex_new",
	args: 0,
	f: func(env *base.Env) base.Value {
		return base.NewGenericValue(&sync.Mutex{})
	},
}

var lib_syncmutexlock = LibFunc{
	name: "mutex_lock",
	args: 1,
	f: func(env *base.Env) base.Value {
		env.R0.Generic().(*sync.Mutex).Lock()
		return base.NewValue()
	},
}

var lib_syncmutexunlock = LibFunc{
	name: "mutex_unlock",
	args: 1,
	f: func(env *base.Env) base.Value {
		env.R0.Generic().(*sync.Mutex).Unlock()
		return base.NewValue()
	},
}
