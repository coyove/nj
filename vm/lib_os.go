package vm

import (
	"os"
	"os/exec"

	"github.com/coyove/bracket/base"
)

var lib_osargs = LibFunc{name: "os_args", args: 0, f: func(env *base.Env) base.Value {
	list := make([]base.Value, len(os.Args))
	for i, arg := range os.Args {
		list[i] = base.NewStringValue(arg)
	}
	return base.NewListValue(list)
},
}

var lib_startprocess = LibFunc{name: "start_process", args: 1, ff: func(env *base.Env) base.Value {
	v := env.Stack().Values()
	exe := v[0].AsString()
	args := make([]string, len(v)-1)
	for i := 1; i < len(v); i++ {
		args[i-1] = v[i].AsString()
	}

	cmd := exec.Command(exe, args...)
	err := cmd.Run()
	if err == nil {
		return base.NewValue()
	}

	return base.NewStringValue(err.(*exec.ExitError).Error())
},
}

var lib_startprocessbg = LibFunc{name: "start_process_bg", args: 1, ff: func(env *base.Env) base.Value {
	v := env.Stack().Values()
	exe := v[0].AsString()
	args := make([]string, len(v)-1)
	for i := 1; i < len(v); i++ {
		args[i-1] = v[i].AsString()
	}

	cmd := exec.Command(exe, args...)
	go cmd.Run()
	return base.NewValue()
},
}

var lib_createfile = LibFunc{name: "create_file", args: 1, f: func(env *base.Env) base.Value {
	f, err := os.Create(env.R0.AsString())
	if err != nil {
		return base.NewStringValue(err.Error())
	}
	return base.NewGenericValue(f)
},
}

var lib_writefile = LibFunc{name: "write_file", args: 1, f: func(env *base.Env) base.Value {
	f := env.R0.AsGeneric().(*os.File)
	n, _ := f.Write(env.R1.AsBytes())
	return base.NewNumberValue(float64(n))
},
}

var lib_closefile = LibFunc{name: "close_file", args: 1, f: func(env *base.Env) base.Value {
	env.R0.AsGeneric().(*os.File).Close()
	return base.NewValue()
},
}
