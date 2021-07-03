package script

import (
	"bytes"
	"io"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"github.com/coyove/script/parser"
)

type Func struct {
	Code       packet
	Name       string
	DocString  string
	StackSize  uint16
	Native     func(env *Env)
	loadGlobal *Program
	Params     []string
	Locals     []string
	MethodSrc  *Map
}

type Program struct {
	Func
	Deadline         int64
	MaxCallStackSize int64
	Stack            *[]Value
	Functions        []*Func
	Stdout, Stderr   io.Writer
	Stdin            io.Reader
	Logger           *log.Logger
	NilIndex         uint16
	GLoad            func(string) Value
	GStore           func(string, Value)
	shadowTable      *symtable
	deadsize         int64
}

func (p *Program) SetDeadsize(v int64) { p.deadsize = v }

func (p *Program) GetDeadsize() int64 { return p.deadsize }

func (p *Program) DecrDeadsize(v int64) {
	if p.deadsize == 0 {
		return
	}
	if atomic.AddInt64(&p.deadsize, -v) <= 0 {
		panic("deadsize")
	}
}

// Native creates a golang-Native function
func Native(name string, f func(env *Env), doc ...string) Value {
	return Function(&Func{
		Name:      name,
		Native:    f,
		DocString: fixDocString(strings.Join(doc, "\n"), name, ""),
	})
}

func Native1(name string, f func(*Env, Value) Value, doc ...string) Value {
	return Native(name, func(env *Env) { env.A = f(env, env.Get(0)) }, doc...)
}

func Native2(name string, f func(*Env, Value, Value) Value, doc ...string) Value {
	return Native(name, func(env *Env) { env.A = f(env, env.Get(0), env.Get(1)) }, doc...)
}

func Native3(name string, f func(*Env, Value, Value, Value) Value, doc ...string) Value {
	return Native(name, func(env *Env) { env.A = f(env, env.Get(0), env.Get(1), env.Get(2)) }, doc...)
}

func NativeWithParamMap(name string, f func(*Env), doc string, params ...string) Value {
	return Function(&Func{
		Name:      name,
		Params:    params,
		DocString: fixDocString(doc, name, strings.Join(params, ",")),
		Native: func(env *Env) {
			if env.A.Type() != VMap {
				args := NewSizedMap(env.Size())
				for i := range env.Stack() {
					if i < len(params) {
						args.Set(String(params[i]), env.Stack()[i])
					}
				}
				env.A = args.Value()
			}
			f(env)
		},
	})
}

func (c *Func) NumParams() int {
	return len(c.Params)
}

func (c *Func) IsNative() bool { return c.Native != nil }

func (c *Func) String() string {
	p := bytes.Buffer{}
	if c.Name != "" {
		p.WriteString(c.Name)
	} else if c.Native != nil {
		p.WriteString("native")
	} else {
		p.WriteString("function")
	}

	p.WriteString("(")
	for _, pa := range c.Params {
		p.WriteString(pa)
		p.WriteString(",")
	}
	if c.NumParams() > 0 {
		p.Truncate(p.Len() - 1)
	}
	p.WriteString(")")
	return p.String()
}

func (c *Func) PrettyCode() string {
	if c.Native != nil {
		return "[Native Code]"
	}
	return pkPrettify(c, c.loadGlobal, false)
}

func (c *Func) exec(newEnv Env) Value {
	if c.Native != nil {
		c.Native(&newEnv)
		return newEnv.A
	}
	return InternalExecCursorLoop(newEnv, c, 0)
}

func (p *Program) Run() (v1 Value, err error) {
	return p.Call()
}

func (p *Program) Call() (v1 Value, err error) {
	defer parser.CatchError(&err)
	newEnv := Env{
		Global: p,
		stack:  p.Stack,
	}
	v1 = InternalExecCursorLoop(newEnv, &p.Func, 0)
	return
}

func (c *Func) Call(a ...Value) (v1 Value, err error) {
	defer parser.CatchError(&err)

	if c.Native != nil {
		newEnv := Env{
			stack: &a,
		}
		newEnv.growZero(int(c.StackSize))
		v1 = c.exec(newEnv)
	} else {
		oldLen := len(*c.loadGlobal.Stack)
		newEnv := Env{
			Global:      c.loadGlobal,
			stack:       c.loadGlobal.Stack,
			StackOffset: uint32(oldLen),
		}

		for _, a := range a {
			newEnv.Push(a)
		}
		newEnv.growZero(int(c.StackSize))

		v1 = c.exec(newEnv)
		*c.loadGlobal.Stack = (*c.loadGlobal.Stack)[:oldLen]
	}
	return
}

func (c *Func) CallMap(a *Map) (v1 Value, err error) {
	defer parser.CatchError(&err)

	if c.Native != nil {
		s := []Value{}
		newEnv := Env{
			stack: &s,
			A:     Interface(a),
		}
		for _, pa := range c.Params {
			newEnv.Push(a.Get(String(pa)))
		}
		newEnv.growZero(int(c.StackSize))
		v1 = c.exec(newEnv)
	} else {
		oldLen := len(*c.loadGlobal.Stack)
		newEnv := Env{
			Global:      c.loadGlobal,
			stack:       c.loadGlobal.Stack,
			StackOffset: uint32(oldLen),
			A:           Interface(a),
		}

		for _, pa := range c.Params {
			newEnv.Push(a.Get(String(pa)))
		}
		newEnv.growZero(int(c.StackSize))

		v1 = c.exec(newEnv)
		*c.loadGlobal.Stack = (*c.loadGlobal.Stack)[:oldLen]
	}
	return
}

func (p *Program) PrettyCode() string { return pkPrettify(&p.Func, p, true) }

func (p *Program) SetTimeout(d time.Duration) { p.Deadline = time.Now().Add(d).UnixNano() }

func (p *Program) SetDeadline(d time.Time) { p.Deadline = d.UnixNano() }

func (p *Program) Print(a ...interface{}) { p.log("", "", a...) }

func (p *Program) Printf(f string, a ...interface{}) { p.log("f", f, a...) }

func (p *Program) Println(a ...interface{}) { p.log("l", "", a...) }

func (p *Program) Fatal(a ...interface{}) { p.log("F", "", a...) }

func (p *Program) Fatalf(f string, a ...interface{}) { p.log("Ff", f, a...) }

func (p *Program) Fatalln(a ...interface{}) { p.log("Fl", "", a...) }

func (p *Program) Panic(a ...interface{}) { p.log("P", "", a...) }

func (p *Program) Panicf(f string, a ...interface{}) { p.log("Pf", f, a...) }

func (p *Program) Panicln(a ...interface{}) { p.log("Pl", "", a...) }

func (p *Program) log(o, f string, a ...interface{}) {
	if p.Logger == nil {
		p.Logger = log.New(p.Stderr, "", log.LstdFlags)
	}
	switch o {
	default:
		p.Logger.Print(a...)
	case "f":
		p.Logger.Printf(f, a...)
	case "l":
		p.Logger.Println(a...)
	case "F":
		p.Logger.Fatal(a...)
	case "Ff":
		p.Logger.Fatalf(f, a...)
	case "Fl":
		p.Logger.Fatalln(a...)
	case "P":
		p.Logger.Panic(a...)
	case "Pf":
		p.Logger.Panicf(f, a...)
	case "Pl":
		p.Logger.Panicln(a...)
	}
}

func (p *Program) Get(k string) (v Value, err error) {
	defer parser.CatchError(&err)
	return (*p.Stack)[int(p.shadowTable.mustGetSymbol(k))], nil
}

func (p *Program) Set(k string, v Value) (err error) {
	defer parser.CatchError(&err)
	(*p.Stack)[int(p.shadowTable.mustGetSymbol(k))] = v
	return nil
}

func fixDocString(in, name, arg string) string {
	in = strings.Replace(in, "$a", arg, -1)
	in = strings.Replace(in, "$f", name, -1)
	return in
}
