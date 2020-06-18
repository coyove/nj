package potatolang

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unsafe"
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
		in := env.In(0, ANY).Any().(File)
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
					if n, _ := in.Read(b[i : i+1]); n != 1 {
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
					buf, _ := ioutil.ReadAll(in)
					ret(i, _StrBytes(buf))
				case strings.HasPrefix(s, "*l"):
					b := make([]byte, 0, 16)
					for {
						b = append(b, 0)
						if n, _ := in.Read(b[len(b)-1:]); n != 1 || b[len(b)-1] == '\n' {
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
		var when = io.SeekCurrent
		if len(env.stack) >= 2 {
			switch env.In(1, STR).Str() {
			case "set":
				when = io.SeekStart
			case "cur":
				when = io.SeekCurrent
			case "end":
				when = io.SeekEnd
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
	fileLinesImpl = NativeFun(func(env *Env) {
		f := bufio.NewReader(env.In(0, ANY).Any().(File))
		env.A = NativeFun(func(env *Env) {
			buf, err := f.ReadBytes('\n')
			if err != nil {
				env.Return(Value{}, Str(err.Error()))
			} else {
				env.A = Str(string(bytes.TrimRight(buf, "\n")))
			}
		})
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
			case "lines":
				env.A = fileLinesImpl
			case "flush":
				env.A = fileFlushImpl
			}
		}))
)

func initLibAux() {
	G.Puts("print", NativeFun(func(env *Env) {
		args := make([]interface{}, len(env.stack))
		for i := range args {
			args[i] = env.stack[i].Any()
		}
		if n, err := fmt.Println(args...); err != nil {
			env.Return(Value{}, Str(err.Error()))
		} else {
			env.Return(Num(float64(n)))
		}
	}))

	lio := &Table{}
	lio.Puts("open", NativeFun(func(env *Env) {
		perm := os.FileMode(0666)
		flag := os.O_RDONLY
		if len(env.stack) == 2 {
			switch strings.Replace(env.In(1, STR).Str(), "b", "", 1) {
			case "r":
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
		env.A, env.V = fileWriteImpl.Fun().Call(append([]Value{Any(DefaultOutput)}, env.stack...)...)
	}))
	lio.Puts("type", NativeFun(func(env *Env) {
		env.A = Value{}
		if v := env.Get(0); v.Type() == ANY {
			if f, ok := v.Any().(File); ok {
				buf := [0]byte{}
				if _, err := f.Read(buf[:]); err == nil {
					env.A = Str("file")
				} else if err.(*os.PathError).Err == os.ErrClosed {
					env.A = Str("closed file")
				} else {
					env.Return(Str("unknown file"), Str(err.Error()))
				}
			}
		}
	}))
	lio.Puts("lines", NativeFun(func(env *Env) {
		i := DefaultInput
		if len(env.stack) == 1 {
			f, err := os.Open(env.In(0, STR).Str())
			if err != nil {
				env.Return(Value{}, Str(err.Error()))
				return
			}
			i = File{f}
		}
		env.A, env.V = fileLinesImpl.Fun().Call(append([]Value{Any(i)}, env.stack...)...)
	}))
	lio.Puts("tmpfile", NativeFun(func(env *Env) {
		p := filepath.Join(os.TempDir(), fmt.Sprintf("pol%d%d", time.Now().Unix(), rand.Int()))
		f, err := os.Create(p)
		if err != nil {
			env.Return(Value{}, Str(err.Error()))
			return
		}
		env.Return(Any(File{f}), Str(p))
	}))
	G.Puts("io", Tab(lio))
	//
	los := &Table{}
	los.Puts("time", NativeFun(func(env *Env) {
		if v := env.Get(0); !v.IsNil() {
			nz := func(v Value) int {
				if v.Type() == NUM {
					return int(v.Num())
				}
				return 0
			}
			t := env.In(0, TAB).Tab()
			env.A = Num(float64(time.Date(
				nz(t.Get(Str("year"))),
				time.Month(nz(t.Get(Str("month")))),
				nz(t.Get(Str("day"))),
				nz(t.Get(Str("hour"))),
				nz(t.Get(Str("min"))),
				nz(t.Get(Str("sec"))), 0, time.UTC).Unix()))
		} else {
			env.A = Num(float64(time.Now().Unix()))
		}
	}))
	los.Puts("clock", NativeFun(func(env *Env) {
		x := time.Now()
		s := *(*[2]int64)(unsafe.Pointer(&x))
		env.A = Num(float64(s[1] / 1e9))
	}))
	los.Puts("microclock", NativeFun(func(env *Env) {
		x := time.Now()
		s := *(*[2]int64)(unsafe.Pointer(&x))
		env.A = Num(float64(s[1] / 1e3))
	}))
	los.Puts("exit", NativeFun(func(env *Env) {
		if v := env.Get(0); !v.IsNil() {
			os.Exit(int(env.In(0, NUM).Num()))
		}
		os.Exit(0)
	}))
	G.Puts("os", Tab(los))

	lstring := &Table{}
	lstring.Puts("format", NativeFun(func(env *Env) {
		f := env.In(0, STR).Str()
		args := make([]interface{}, 0, len(env.stack))
		for x, i, f := byte(0), 1, f; ; {
			x, f = findNextFormat(f)
			if x == 0 {
				break
			}
			switch x {
			case 'c', 'd', 'i':
				args = append(args, atoint64(env.Get(i)))
			case 'o', 'u', 'X', 'x':
				args = append(args, atouint64(env.Get(i)))
			default:
				args = append(args, env.Get(i).Any())
			}
		}
		env.A = Str(fmt.Sprintf(f, args...))
	}))
	lstring.Puts("rep", NativeFun(func(env *Env) { env.A = Str(strings.Repeat(env.In(0, STR).Str(), int(env.In(1, NUM).Num()))) }))
	lstring.Puts("char", NativeFun(func(env *Env) { env.A = Str(string(rune(env.In(0, NUM).Num()))) }))
	G.Puts("string", Tab(lstring))
}

func findNextFormat(f string) (byte, string) {
	i := strings.Index(f, "%")
	if i == -1 || i == len(f)-1 {
		return 0, "" // no more format strings
	}
	i++
	if f[i] == '%' { // %%
		return findNextFormat(f[i+1:])
	}
	for i < len(f) {
		switch f[i] {
		case '#', '+', '-', ' ', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			i++
		default:
			return f[i], f[i+1:]
		}
	}
	return 0, ""
}
