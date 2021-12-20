package nj

import (
	"bytes"
	"fmt"
	"io"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

type FuncBody struct {
	Code       packet
	StackSize  uint16
	NumParams  uint16
	Variadic   bool
	Dummy      bool
	Native     func(env *Env)
	Name       string
	DocString  string
	LoadGlobal *Program
	Locals     []string
	Object     *Object
}

type Program struct {
	Top          *FuncBody
	Symbols      map[string]*symbol
	MaxStackSize int64
	Options      *CompileOptions
	Stack        *[]Value
	Functions    []*Object
	Stdout       io.Writer
	Stderr       io.Writer
	Stdin        io.Reader
}

func dummyFunc(*Env) {}

// Func creates a callable object
func Func(name string, f func(*Env), doc string) Value {
	if name == "" {
		name = internal.UnnamedFunc
	}
	obj := NewObject(0)
	obj.Callable = &FuncBody{
		Name:      name,
		Native:    f,
		DocString: doc,
		Object:    obj,
	}
	if f == nil {
		obj.Callable.Native = dummyFunc
		obj.Callable.Dummy = true
	}
	obj.SetPrototype(FuncProto)
	return obj.ToValue()
}

func (p *Program) Run() (v1 Value, err error) {
	defer internal.CatchError(&err)
	newEnv := Env{
		Global: p,
		stack:  p.Stack,
	}
	v1 = internalExecCursorLoop(newEnv, p.Top, nil)
	return
}

// EmergStop terminates the execution of program
// After calling, program will become unavailable for any further operations
func (p *Program) EmergStop() {
	p.Top.EmergStop()
	for _, f := range p.Functions {
		f.Callable.EmergStop()
	}
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

func (p *Program) LocalsObject() *Object {
	r := NewObject(len(p.Top.Locals))
	for i, name := range p.Top.Locals {
		r.Set(Str(name), (*p.Stack)[i])
	}
	return r
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
			*err = fmt.Errorf("not callable")
		}
		return
	}
	if err != nil {
		defer internal.CatchErrorFuncCall(err, m.Callable.Name)
	}
	r := Runtime{Stacktrace: e.Runtime().StacktraceWithCurrent()}
	return m.Callable.execute(r, this, args...)
}

func (c *FuncBody) String() string {
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

func (c *FuncBody) ToCode() string {
	if c.Native != nil {
		return "[Native Code]"
	}
	return pkPrettify(c, c.LoadGlobal, false)
}

func (f *FuncBody) Copy() *FuncBody {
	f2 := *f
	f2.Code.Code = append([]_inst{}, f2.Code.Code...)
	return &f2
}

// EmergStop terminates the execution of Func
// After calling, FuncBody will become unavailable for any further operations
func (c *FuncBody) EmergStop() {
	for i := range c.Code.Code {
		c.Code.Code[i] = inst(typ.OpRet, regA, 0)
	}
}

func (c *FuncBody) execute(r Runtime, this Value, args ...Value) (v1 Value) {
	newEnv := Env{
		A:       this,
		Global:  c.LoadGlobal,
		stack:   &args,
		runtime: r,
	}

	if c.Native != nil {
		newEnv.NativeSelf = c
		c.Native(&newEnv)
		v1 = newEnv.A
	} else {
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
		v1 = internalExecCursorLoop(newEnv, c, r.Stacktrace)
	}
	return
}
