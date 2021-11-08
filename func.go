package script

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unsafe"

	"github.com/coyove/script/parser"
	"github.com/coyove/script/typ"
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
	Receiver   Value
}

type Program struct {
	Func
	MaxCallStackSize int64
	Stack            *[]Value
	Functions        []*Func
	Stdout           io.Writer
	Stderr           io.Writer
	Stdin            io.Reader
	NilIndex         uint16
	shadowTable      *symtable
}

// Native creates a golang-Native function
func Native(name string, f func(env *Env), doc ...string) Value {
	if name == "" {
		name = "<native>"
	}
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

func (c *Func) Value() Value { return Value{v: uint64(typ.Func), p: unsafe.Pointer(c)} }

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
	if c.Receiver != Nil {
		p.WriteString("{" + c.Receiver.Type().String() + "},")
	}
	for i := 0; i < int(c.NumParams); i++ {
		fmt.Fprintf(&p, "a%d,", i)
	}
	if c.Variadic {
		p.Truncate(p.Len() - 1)
		p.WriteString("...")
	} else if p.Bytes()[p.Len()-1] == ',' {
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
	v1 = internalExecCursorLoop(newEnv, &p.Func, 0)
	return
}

func (p *Program) EmergStop() {
	p.Func.EmergStop()
	for _, f := range p.Functions {
		f.EmergStop()
	}
}

// EmergStop terminates the execution of Func
// After calling, Func will become unavailable for any further operations
func (c *Func) EmergStop() {
	for i := range c.Code.Code {
		c.Code.Code[i] = inst(typ.OpRet, regA, 0)
	}
}

func (c *Func) CallSimple(args ...interface{}) (v1 interface{}, err error) {
	x := make([]Value, len(args))
	for i := range args {
		x[i] = Val(args[i])
	}
	return c.Call(x...)
}

func (c *Func) Call(args ...Value) (v1 Value, err error) {
	defer parser.CatchErrorFuncCall(&err, c.Name)

	newEnv := Env{
		Global: c.LoadGlobal,
		stack:  &args,
	}
	if c.Receiver != Nil {
		newEnv.Prepend(c.Receiver)
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
		newEnv.growZero(int(c.StackSize), int(c.NumParams))
		v1 = internalExecCursorLoop(newEnv, c, 0)
	}
	return
}

func (f *Func) Pure() *Func { f.Receiver = Nil; return f }

func (p *Program) PrettyCode() string { return pkPrettify(&p.Func, p, true) }

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
