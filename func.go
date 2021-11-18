package nj

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unsafe"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

type Function struct {
	*FuncBody
	Receiver Value
}

type FuncBody struct {
	Code       packet
	Name       string
	DocString  string
	StackSize  uint16
	NumParams  uint16
	Variadic   bool
	Native     func(env *Env)
	LoadGlobal *Program
	Locals     []string
}

type Program struct {
	Top          *Function
	Symbols      map[string]*symbol
	MaxStackSize int64
	Stack        *[]Value
	Functions    []*Function
	Stdout       io.Writer
	Stderr       io.Writer
	Stdin        io.Reader
}

// Func creates a function
func Func(name string, f func(env *Env), doc ...string) Value {
	if name == "" {
		name = "<native>"
	}
	return (&Function{
		FuncBody: &FuncBody{
			Name:      name,
			Native:    f,
			DocString: fixDocString(strings.Join(doc, "\n"), name, ""),
		},
	}).Value()
}

func Func1(name string, f func(Value) Value, doc ...string) Value {
	return Func(name, func(env *Env) { env.A = f(env.B(0)) }, doc...)
}

func Func2(name string, f func(Value, Value) Value, doc ...string) Value {
	return Func(name, func(env *Env) { env.A = f(env.B(0), env.B(1)) }, doc...)
}

func Func3(name string, f func(Value, Value, Value) Value, doc ...string) Value {
	return Func(name, func(env *Env) { env.A = f(env.B(0), env.B(1), env.B(2)) }, doc...)
}

func (c *Function) Value() Value {
	return Value{v: uint64(typ.Func), p: unsafe.Pointer(c)}
}

func (c *Function) String() string {
	p := bytes.Buffer{}
	if c.Name != "" {
		p.WriteString(c.Name)
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

func (c *Function) PrettyCode() string {
	if c.Native != nil {
		return "[Native Code]"
	}
	return pkPrettify(c, c.LoadGlobal, false)
}

func (p *Program) Run() (v1 Value, err error) {
	defer internal.CatchError(&err)
	newEnv := Env{
		Global: p,
		stack:  p.Stack,
	}
	v1 = internalExecCursorLoop(newEnv, p.Top, 0)
	return
}

// EmergStop terminates the execution of program
// After calling, program will become unavailable for any further operations
func (p *Program) EmergStop() {
	p.Top.EmergStop()
	for _, f := range p.Functions {
		f.EmergStop()
	}
}

// EmergStop terminates the execution of Func
// After calling, Func will become unavailable for any further operations
func (c *Function) EmergStop() {
	for i := range c.Code.Code {
		c.Code.Code[i] = inst(typ.OpRet, regA, 0)
	}
}

func (c *Function) CallVal(args ...interface{}) (v1 interface{}, err error) {
	x := make([]Value, len(args))
	for i := range args {
		x[i] = Val(args[i])
	}
	return c.Call(x...)
}

func (c *Function) Call(args ...Value) (v1 Value, err error) {
	defer internal.CatchErrorFuncCall(&err, c.Name)

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

func (f *Function) Pure() *Function {
	f.Receiver = Nil
	return f
}

func (p *Program) PrettyCode() string {
	return pkPrettify(p.Top, p, true)
}

func (p *Program) Get(k string) (v Value, ok bool) {
	addr := p.Symbols[k]
	if addr == nil {
		return Nil, false
	}
	return (*p.Stack)[addr.addr], true
}

func (p *Program) Set(k string, v Value) (ok bool) {
	addr := p.Symbols[k]
	if addr == nil {
		return false
	}
	(*p.Stack)[addr.addr] = v
	return true
}

func (f *Function) Copy() *Function {
	b2 := *f.FuncBody
	b2.Code.Code = append([]_inst{}, b2.Code.Code...)
	return &Function{
		FuncBody: &b2,
		Receiver: f.Receiver,
	}
}

func fixDocString(in, name, arg string) string {
	in = strings.Replace(in, "$a", arg, -1)
	in = strings.Replace(in, "$f", name, -1)
	return in
}
