package potatolang

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type File struct {
	*os.File
}

func (s File) GetMetatable() *Table { return fileMetatable }

var (
	DefaultInput  = File{os.Stdin}
	DefaultOutput = File{os.Stdout}
	fileWriteImpl = NativeFun(func(env *Env) {
		out := env.In(0, ANY).Any().(File)
		for _, v := range env.stack[1:] {
			switch v.Type() {
			case STR:
				out.Write(v._StrBytes())
			default:
				out.Write([]byte(v.String()))
			}
		}
	})
	fileReadImpl = NativeFun(func(env *Env) {
		DefaultInput := env.In(0, ANY).Any().(File)
		ret := func(i int, v Value) {
			if i == 0 {
				env.A = v
			} else {
				env.V = append(env.V, v)
			}
		}
		for i, a := range env.stack[1:] {
			switch a.Type() {
			case NUM:
				b := make([]byte, int(a.Num()))
				for i := range b {
					if n, _ := DefaultInput.Read(b[i : i+1]); n != 1 {
						b = b[:i]
						break
					}
				}
				if len(b) == 0 {
					ret(i, Value{})
				} else {
					ret(i, _StrBytes(b))
				}
			case STR:
				switch s := a.Str(); {
				case strings.HasPrefix(s, "*a"):
					buf, _ := ioutil.ReadAll(DefaultInput)
					ret(i, _StrBytes(buf))
				case strings.HasPrefix(s, "*l"):
					b := make([]byte, 0, 16)
					for {
						b = append(b, 0)
						if n, _ := DefaultInput.Read(b[len(b)-1:]); n != 1 || b[len(b)-1] == '\n' {
							b = b[:len(b)-1]
							break
						}
					}
					ret(i, _StrBytes(b))
				case strings.HasPrefix(s, "*n"):
					var n float64
					if _, err := fmt.Scanf("%f", &n); err != nil {
						ret(i, Value{})
					} else {
						ret(i, Num(n))
					}
				default:
					panicf("bad argument")
				}
			default:
				a.ExpectMsg(STR, "io.read")
			}
		}
	})
	fileSeekImpl = NativeFun(func(env *Env) {
		var off int64 = 0
		var when = os.SEEK_CUR
		if len(env.stack) >= 2 {
			switch env.In(1, STR).Str() {
			case "set":
				when = os.SEEK_SET
			case "cur":
				when = os.SEEK_CUR
			case "end":
				when = os.SEEK_END
			default:
				panic("bad argument")
			}
		}
		if len(env.stack) == 3 {
			off = int64(env.In(2, NUM).Num())
		}
		if n, err := env.In(0, ANY).Any().(File).Seek(off, when); err == nil {
			env.A = Num(float64(n))
		} else {
			env.Return(Value{}, Str(err.Error()))
		}
	})
	fileCloseImpl = NativeFun(func(env *Env) {
		if err := env.In(0, ANY).Any().(File).Close(); err == nil {
			env.A = Bln(true)
		} else {
			env.Return(Value{}, Str(err.Error()))
		}
	})
	fileFlushImpl = NativeFun(func(env *Env) {
		if err := env.In(0, ANY).Any().(File).Sync(); err == nil {
			env.A = Bln(true)
		} else {
			env.Return(Value{}, Str(err.Error()))
		}
	})
	fileMetatable = (&Table{}).
			Put(M__index, NativeFun(func(env *Env) {
			switch env.In(1, STR).Str() {
			case "read":
				env.A = fileReadImpl
			case "write":
				env.A = fileWriteImpl
			case "close":
				env.A = fileCloseImpl
			case "seek":
				env.A = fileSeekImpl
			case "flush":
				env.A = fileFlushImpl
			}
		}))
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
	G.Puts("print", NativeFun(fmtPrint('l')))

	lio := &Table{}
	lio.Puts("open", NativeFun(func(env *Env) {
		perm := os.FileMode(0777)
		flag := os.O_RDONLY
		if len(env.stack) == 2 {
			switch strings.Replace(env.In(1, STR).Str(), "b", "", 1) {
			case "w":
				flag = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
			case "a":
				flag = os.O_WRONLY | os.O_CREATE | os.O_APPEND
			case "r+":
				flag = os.O_RDWR
			case "w+":
				flag = os.O_RDWR | os.O_CREATE | os.O_TRUNC
			case "a+":
				flag = os.O_RDWR | os.O_CREATE | os.O_APPEND
			default:
				panic("bad argument")
			}
		}
		f, err := os.OpenFile(env.In(0, STR).Str(), flag, perm)
		if err == nil {
			env.A = Any(File{f})
		} else {
			env.Return(Value{}, Str(err.Error()))
		}
	}))
	lio.Puts("read", NativeFun(func(env *Env) {
		env.A, env.V = fileReadImpl.Fun().Call(append([]Value{Any(DefaultInput)}, env.stack...)...)
	}))
	lio.Puts("write", NativeFun(func(env *Env) {
		env.A, env.V = fileWriteImpl.Fun().Call(append([]Value{Any(DefaultInput)}, env.stack...)...)
	}))
	G.Puts("io", Tab(lio))
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
