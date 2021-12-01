package nj

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

type FuncBody struct {
	Code       packet
	StackSize  uint16
	NumParams  uint16
	Variadic   bool
	Native     func(env *Env)
	Name       string
	DocString  string
	LoadGlobal *Program
	Locals     []string
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
func Func(name string, f func(*Env), doc ...string) Value {
	if name == "" {
		name = internal.UnnamedFunc
	}
	if f == nil {
		f = dummyFunc
	}
	return (&Object{
		Callable: &FuncBody{
			Name:      name,
			Native:    f,
			DocString: strings.Join(doc, "\n"),
		},
	}).ToValue()
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
		A:      this,
		Global: c.LoadGlobal,
		stack:  &args,
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
		v1 = internalExecCursorLoop(newEnv, c, r.Stacktrace)
	}
	return
}
