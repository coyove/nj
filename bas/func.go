package bas

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

type Function struct {
	CodeSeg    Packet
	StackSize  uint16
	NumParams  uint16
	Variadic   bool
	Dummy      bool
	Native     func(env *Env)
	Name       string
	DocString  string
	LoadGlobal *Program
	Locals     []string
	obj        *Object
}
type Environment struct {
	MaxStackSize int64
	Globals      *Object
	Stdout       io.Writer
	Stderr       io.Writer
	Stdin        io.Reader
}

type Program struct {
	top       *Function
	symbols   map[string]*typ.Symbol
	stack     *[]Value
	functions []*Object
	Environment
}

func NewProgram(coreStack *Env, top *Function, symbols map[string]*typ.Symbol, funcs []*Object, env *Environment) *Program {
	cls := &Program{top: top}
	cls.stack = coreStack.stack
	cls.symbols = symbols
	cls.functions = funcs
	if env != nil {
		cls.Environment = *env
	}
	cls.Stdout = or(cls.Stdout, os.Stdout).(io.Writer)
	cls.Stdin = or(cls.Stdin, os.Stdin).(io.Reader)
	cls.Stderr = or(cls.Stderr, os.Stderr).(io.Writer)

	cls.top.LoadGlobal = cls
	for _, f := range cls.functions {
		f.fun.LoadGlobal = cls
	}
	return cls
}

// Func creates a callable object
func Func(name string, f func(*Env), doc string) Value {
	if name == "" {
		name = internal.UnnamedFunc
	}
	obj := NewObject(0)
	obj.fun = &Function{
		Name:      name,
		Native:    f,
		DocString: doc,
		obj:       obj,
	}
	if f == nil {
		obj.fun.Native = func(*Env) {}
		obj.fun.Dummy = true
	}
	obj.SetPrototype(Proto.Func)
	return obj.ToValue()
}

func (p *Program) Run() (v1 Value, err error) {
	defer internal.CatchError(&err)
	newEnv := Env{
		Global: p,
		stack:  p.stack,
	}
	v1 = internalExecCursorLoop(newEnv, p.top, nil)
	return
}

// Stop terminates the execution of program
// After calling, program will become unavailable for any further operations
// There is no way to terminate goroutines and blocking I/Os
func (p *Program) Stop() {
	stop := func(c *Function) {
		for i := range c.CodeSeg.Code {
			c.CodeSeg.Code[i] = typ.Inst{Opcode: typ.OpRet, A: typ.RegA}
		}
	}
	stop(p.top)
	for _, f := range p.functions {
		stop(f.fun)
	}
}

func (p *Program) GoString() string {
	return pkPrettify(p.top, p, true)
}

func (p *Program) Get(k string) (v Value, ok bool) {
	addr, ok := p.symbols[k]
	if !ok {
		return Nil, false
	}
	return (*p.stack)[addr.Address], true
}

func (p *Program) Set(k string, v Value) (ok bool) {
	addr, ok := p.symbols[k]
	if !ok {
		return false
	}
	(*p.stack)[addr.Address] = v
	return true
}

func (p *Program) LocalsObject() *Object {
	r := NewObject(len(p.top.Locals))
	for i, name := range p.top.Locals {
		r.Set(Str(name), (*p.stack)[i])
	}
	return r
}

func EnvForAsyncCall(e *Env) *Env {
	e2 := *e
	e2.runtime.StackN = append([]Stacktrace{}, e2.runtime.StackN...)
	return &e2
}

func Call(m *Object, args ...Value) (res Value) {
	return CallObject(m, nil, nil, m.this, args...)
}

func Call2(m *Object, args ...Value) (res Value, err error) {
	res = CallObject(m, nil, &err, m.this, args...)
	return
}

func CallObject(m *Object, e *Env, err *error, this Value, args ...Value) (res Value) {
	if !m.IsCallable() {
		if err == nil {
			internal.Panic("%v not callable", simpleString(m.ToValue()))
		} else {
			*err = fmt.Errorf("%v not callable", simpleString(m.ToValue()))
		}
		return
	}
	if err != nil {
		defer internal.CatchErrorFuncCall(err, m.fun.Name)
	}

	c := m.fun
	newEnv := Env{
		A:      this,
		Global: c.LoadGlobal,
		stack:  &args,
	}

	if c.Native != nil {
		if e == nil {
			newEnv.runtime.Callable0 = c
		} else {
			newEnv.runtime = e.runtime.Push(c)
		}
		c.Native(&newEnv)
		return newEnv.A
	}

	if c.Variadic {
		s := *newEnv.stack
		if len(s) > int(c.NumParams)-1 {
			s[c.NumParams-1] = NewArray(append([]Value{}, s[c.NumParams-1:]...)...).ToValue()
		} else {
			newEnv.grow(int(c.NumParams))
			newEnv._set(c.NumParams-1, NewArray().ToValue())
		}
	}
	newEnv.growZero(int(c.StackSize), int(c.NumParams))

	var stk []Stacktrace
	if e != nil {
		stk = e.runtime.Stacktrace()
	}
	return internalExecCursorLoop(newEnv, c, stk)
}

func (c *Function) String() string {
	p := bytes.Buffer{}
	if c.Name != "" {
		p.WriteString(c.Name)
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
	} else if p.Bytes()[p.Len()-1] == ',' {
		p.Truncate(p.Len() - 1)
	}
	p.WriteString(")")
	return p.String()
}

func (c *Function) GoString() string {
	if c.Native != nil {
		return "[Native Code]"
	}
	return pkPrettify(c, c.LoadGlobal, false)
}
