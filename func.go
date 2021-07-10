package script

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
	"unsafe"

	"github.com/coyove/script/parser"
)

type Func struct {
	Code       packet
	Name       string
	DocString  string
	StackSize  uint16
	NumParams  uint16
	Variadic   bool
	Native     func(env *Env)
	LoadGlobal *Program
	Locals     []string
	MethodSrc  Value
}

type Program struct {
	Func
	Deadline         int64
	MaxCallStackSize int64
	Stack            *[]Value
	Functions        []*Func
	Stdout           io.Writer
	Stderr           io.Writer
	Stdin            io.Reader
	Logger           *log.Logger
	NilIndex         uint16
	GLoad            func(string) Value
	GStore           func(string, Value)
	shadowTable      *symtable
}

type WrappedFunc struct {
	*Func
}

// Native creates a golang-Native function
func Native(name string, f func(env *Env), doc ...string) Value {
	return (&Func{
		Name:      name,
		Native:    f,
		DocString: fixDocString(strings.Join(doc, "\n"), name, ""),
	}).Value()
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

func Native4(name string, f func(*Env, Value, Value, Value, Value) Value, doc ...string) Value {
	return Native(name, func(env *Env) { env.A = f(env, env.Get(0), env.Get(1), env.Get(2), env.Get(3)) }, doc...)
}

func (c *Func) IsNative() bool { return c.Native != nil }

func (c *Func) Value() Value { return Value{v: uint64(FUNC), p: unsafe.Pointer(c)} }

func (c *Func) WrappedValue() Value { return _interface(&WrappedFunc{c}) }

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
	for i := 0; i < int(c.NumParams); i++ {
		fmt.Fprintf(&p, "a%d,", i)
	}
	if c.Variadic {
		p.Truncate(p.Len() - 1)
		p.WriteString("...")
	} else if c.NumParams > 0 {
		p.Truncate(p.Len() - 1)
	}
	p.WriteString(")")
	return p.String()
}

func (c *Func) PrettyCode() string {
	if c.Native != nil {
		return "[Native Code]"
	}
	return pkPrettify(c, c.LoadGlobal, false)
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

func (c *Func) CallSimple(args ...interface{}) (v1 interface{}, err error) {
	x := make([]Value, len(args))
	for i := range args {
		x[i] = Go(args[i])
	}
	return c.Call(x...)
}

func (c *Func) Call(args ...Value) (v1 Value, err error) {
	defer parser.CatchError(&err)

	newEnv := Env{
		Global: c.LoadGlobal,
		stack:  &args,
	}
	if c.MethodSrc != Nil {
		newEnv.Prepend(c.MethodSrc)
	}

	if c.Native != nil {
		c.Native(&newEnv)
		v1 = newEnv.A
	} else {
		if c.Variadic {
			s := *newEnv.stack
			if len(s) > int(c.NumParams)-1 {
				s[c.NumParams-1] = Array(append([]Value{}, s[c.NumParams-1:]...)...)
			} else {
				newEnv.grow(int(c.NumParams))
				newEnv._set(c.NumParams-1, Array())
			}
		}
		newEnv.growZero(int(c.StackSize))
		v1 = InternalExecCursorLoop(newEnv, c, 0)
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
