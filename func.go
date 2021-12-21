package nj

import (
	"bytes"
	"fmt"
	"io"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

type function struct {
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
	top       *function
	symbols   map[string]*symbol
	stack     *[]Value
	functions []*Object
	Environment
}

// Func creates a callable object
func Func(name string, f func(*Env), doc string) Value {
	if name == "" {
		name = internal.UnnamedFunc
	}
	obj := NewObject(0)
	obj.fun = &function{
		Name:      name,
		Native:    f,
		DocString: doc,
		obj:       obj,
	}
	if f == nil {
		obj.fun.Native = func(*Env) {}
		obj.fun.Dummy = true
	}
	obj.SetPrototype(FuncProto)
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
	stop := func(c *function) {
		for i := range c.CodeSeg.Code {
			c.CodeSeg.Code[i] = inst(typ.OpRet, regA, 0)
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
	addr := p.symbols[k]
	if addr == nil {
		return Nil, false
	}
	return (*p.stack)[addr.addr], true
}

func (p *Program) Set(k string, v Value) (ok bool) {
	addr := p.symbols[k]
	if addr == nil {
		return false
	}
	(*p.stack)[addr.addr] = v
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
			internal.Panic("%v not callable", showType(m.ToValue()))
		} else {
			*err = fmt.Errorf("%v not callable", showType(m.ToValue()))
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

func (c *function) String() string {
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

func (c *function) GoString() string {
	if c.Native != nil {
		return "[Native Code]"
	}
	return pkPrettify(c, c.LoadGlobal, false)
}
