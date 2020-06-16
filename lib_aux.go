package potatolang

import (
	"fmt"
	"os"
	"strings"
)

func fmtPrint(flag byte) func(env *Env) {
	return func(env *Env) {
		args := make([]interface{}, len(env.stack))
		for i := range args {
			args[i] = env.stack[i].Any()
		}
		var n int
		var err error

		switch flag {
		case 'l':
			n, err = fmt.Println(args...)
		case 'f':
			n, err = fmt.Printf(args[0].(string), args[1:]...)
		default:
			n, err = fmt.Print(args...)
		}

		if err != nil {
			env.Return(Value{}, Str(err.Error()))
		} else {
			env.Return(Num(float64(n)))
		}
	}
}

func fmtSprint(flag byte) func(env *Env) {
	return func(env *Env) {
		args := make([]interface{}, len(env.stack))
		for i := range args {
			args[i] = env.stack[i].Any()
		}
		var n string
		switch flag {
		case 'l':
			n = fmt.Sprintln(args...)
		case 'f':
			n = fmt.Sprintf(env.Get(0).Expect(STR).Str(), args...)
		default:
			n = fmt.Sprint(args...)
		}
		env.A = Str(n)
	}
}

// func fmtFprint(flag byte) func(env *Env) Value {
// 	return func(env *Env) Value {
// 		args := make([]interface{}, env.Size())
// 		for i := range args {
// 			args[i] = env.Get(i).AsInterface()
// 		}
// 		var n int
// 		var err error
// 		switch flag {
// 		case 'l':
// 			n, err = fmt.Fprintln(args[0].(io.Writer), args[1:]...)
// 		case 'f':
// 			n, err = fmt.Fprintf(args[0].(io.Writer), args[1].(string), args[2:]...)
// 		default:
// 			n, err = fmt.Fprint(args[0].(io.Writer), args[1:]...)
// 		}
//
// 		if err != nil {
// 			env.B = Str(err.Error())
// 		}
// 		return Num(float64(n))
// 	}
// }
//
// func fmtScan(flag string) func(env *Env) Value {
// 	return func(env *Env) Value {
// 		var start int
// 		switch flag {
// 		case "scanf", "sscanln", "sscan", "fscan", "fscanln":
// 			start = 1
// 		case "sscanf", "fscanf":
// 			start = 2
// 		}
//
// 		receivers := make([]interface{}, env.Size())
// 		for i := start; i < len(receivers); i++ {
// 			switch LoadPointerUnsafe(env.Get(i)).Type() {
// 			case STR:
// 				receivers[i] = new(string)
// 			case NUM:
// 				receivers[i] = new(float64)
// 			default:
// 				panicf("Scan: only string and number are supported")
// 			}
// 		}
//
// 		var n int
// 		var err error
//
// 		switch flag {
// 		case "scanln":
// 			n, err = fmt.Scanln(receivers...)
// 		case "scanf":
// 			n, err = fmt.Scanf(string(env.Get(0).MustString()), receivers[1:]...)
// 		case "scan":
// 			n, err = fmt.Scan(receivers...)
// 		case "sscanln":
// 			n, err = fmt.Sscanln(string(env.Get(0).MustString()), receivers[1:]...)
// 		case "sscanf":
// 			n, err = fmt.Sscanf(string(env.Get(0).MustString()), string(env.Get(1).MustString()), receivers[2:]...)
// 		case "sscan":
// 			n, err = fmt.Sscan(string(env.Get(0).MustString()), receivers[1:]...)
// 		case "fscan":
// 			n, err = fmt.Fscan(env.Get(0).AsInterface().(io.Reader), receivers[1:]...)
// 		case "fscanln":
// 			n, err = fmt.Fscanln(env.Get(0).AsInterface().(io.Reader), receivers[1:]...)
// 		case "fscanf":
// 			n, err = fmt.Fscanf(env.Get(0).AsInterface().(io.Reader), string(env.Get(1).MustString()), receivers[1:]...)
// 		}
//
// 		if err == nil {
// 			for i := start; i < len(receivers); i++ {
// 				switch v := receivers[i].(type) {
// 				case *float64:
// 					StorePointerUnsafe(env.Get(i), Num(*v))
// 				case *string:
// 					StorePointerUnsafe(env.Get(i), Str(*v))
// 				}
// 			}
// 		}
//
// 		if err != nil {
// 			env.B = Str(err.Error())
// 		}
// 		return Num(float64(n))
// 	}
// }
//
// func fmtWrite(env *Env) Value {
// 	var n int
// 	var err error
// 	f := env.Get(0).AsInterface().(io.Writer)
//
// 	for i := 1; i < env.Size(); i++ {
// 		switch a := env.Get(i); a.Type() {
// 		case STR:
// 			n, err = f.write([]byte(env.Get(i).Str()))
// 		default:
// 			panicf("stdWrite can't write: %+v", a)
// 		}
// 	}
//
// 	if err != nil {
// 		env.B = Str(err.Error())
// 	}
// 	return Num(float64(n))
// }
//

// func stringsIndex(env *Env) Value {
// 	return Num(float64(strings.Index(env.Get(0).MustString(), env.Get(1).MustString())))
// }

func initLibAux() {
	// 	lfmt.Put("Printf", NativeFun(1, fmtPrint('f')))
	// 	lfmt.Put("Sprintln", NativeFun(0, fmtSprint('l')))
	// 	lfmt.Put("Sprint", NativeFun(0, fmtSprint(0)))
	// 	lfmt.Put("Sprintf", NativeFun(1, fmtSprint('f')))
	// 	lfmt.Put("Fprintln", NativeFun(1, fmtFprint('l')))
	// 	lfmt.Put("Fprint", NativeFun(1, fmtFprint(0)))
	// 	lfmt.Put("Fprintf", NativeFun(2, fmtFprint('f')))
	// 	lfmt.Put("Scanln", NativeFun(1, fmtScan("scanln")))
	// 	lfmt.Put("Scan", NativeFun(1, fmtScan("scan")))
	// 	lfmt.Put("Scanf", NativeFun(1, fmtScan("scanf")))
	// 	lfmt.Put("Sscanln", NativeFun(1, fmtScan("sscanln")))
	// 	lfmt.Put("Sscan", NativeFun(1, fmtScan("sscan")))
	// 	lfmt.Put("Sscanf", NativeFun(2, fmtScan("sscanf")))
	// 	lfmt.Put("Fscanln", NativeFun(1, fmtScan("fscanln")))
	// 	lfmt.Put("Fscan", NativeFun(1, fmtScan("fscan")))
	// 	lfmt.Put("Fscanf", NativeFun(2, fmtScan("fscanf")))
	// 	lfmt.Put("write", NativeFun(0, fmtWrite))
	G.Puts("print", NativeFun(fmtPrint('l')))
	//
	los := &Table{}
	los.Puts("stdout", Any(os.Stdout))
	los.Puts("stdin", Any(os.Stdin))
	los.Puts("stderr", Any(os.Stderr))
	G.Puts("os", Tab(los))

	lstring := &Table{}
	lstring.Puts("format", NativeFun(fmtSprint('f')))
	lstring.Puts("rep", NativeFun(func(env *Env) { env.A = Str(strings.Repeat(env.In(0, STR).Str(), int(env.In(1, NUM).Num()))) }))
	lstring.Puts("char", NativeFun(func(env *Env) { env.A = Str(string(rune(env.In(0, NUM).Num()))) }))
	G.Puts("string", Tab(lstring))
}
