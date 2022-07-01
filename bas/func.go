package bas

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

var (
	Stdout io.Writer = os.Stdout
	Stderr io.Writer = os.Stderr
	Stdin  io.Reader = os.Stdin
)

type funcbody struct {
	CodeSeg    internal.Packet
	StackSize  uint16
	NumParams  byte
	Variadic   bool
	Method     bool
	EnvgNeeded bool // an execution env (*Env) with Global is required
	Native     func(env *Env)
	Name       string
	LoadGlobal *Program
	Locals     []string
}

type Environment struct {
	MaxStackSize int64
	Globals      *Object
	File         string
	Source       string
	Deadline     time.Time
}

type Program struct {
	top       *Object
	symbols   *Object
	stack     *[]Value
	functions []*Object
	Environment
}

func NewProgram(file, source string, coreStack *Env, top, symbols *Object, funcs []*Object, env *Environment) *Program {
	cls := &Program{top: top}
	cls.stack = coreStack.stack
	cls.symbols = symbols
	cls.functions = funcs
	if env != nil {
		cls.Environment = *env
	}
	cls.File = file
	cls.Source = source

	cls.top.fun.LoadGlobal = cls
	for _, f := range cls.functions {
		f.fun.LoadGlobal = cls
	}
	return cls
}

// Func creates a callable object
func Func(name string, f func(*Env)) Value {
	if name == "" {
		name = internal.UnnamedFunc()
	}
	if f == nil {
		f = func(*Env) {}
	}
	obj := NewObject(0)
	obj.fun = &funcbody{Name: name, Native: f}
	obj.SetPrototype(Proto.Func)
	return obj.ToValue()
}

// EnvFunc creates a callable object which strictly requires a valid execution env
func EnvFunc(name string, f func(*Env)) Value {
	o := Func(name, f).Object()
	o.fun.EnvgNeeded = true
	return o.ToValue()
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

func (p *Program) Stop() {
	return
}

func (p *Program) GoString() string {
	return pkPrettify(p.top.fun, p, true)
}

func (p *Program) Get(k string) (v Value, ok bool) {
	addr := p.symbols.Prop(k)
	if addr == Nil {
		return Nil, false
	}
	return (*p.stack)[addr.Int64()], true
}

func (p *Program) Set(k string, v Value) (ok bool) {
	addr := p.symbols.Prop(k)
	if addr == Nil {
		return false
	}
	(*p.stack)[addr.Int64()] = v
	return true
}

func (p *Program) LocalsObject() *Object {
	r := NewObject(len(p.top.fun.Locals))
	for i, name := range p.top.fun.Locals {
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

func CallObject(m *Object, e *Env, outErr *error, this Value, args ...Value) (res Value) {
	if !m.IsCallable() {
		if outErr == nil {
			internal.Panic("%v not callable", detail(m.ToValue()))
		} else {
			*outErr = fmt.Errorf("%v not callable", detail(m.ToValue()))
		}
		return
	}

	c := m.fun
	newEnv := Env{
		A:      this,
		Global: c.LoadGlobal,
		stack:  &args,
	}

	if outErr != nil {
		defer internal.CatchError(outErr)
	}

	if c.Native != nil {
		st := Stacktrace{Callable: m, stackOffsetFlag: internal.FlagNativeCall}
		defer relayPanic(func() []Stacktrace { return newEnv.runtime.Stacktrace() })
		if e == nil {
			newEnv.runtime.Stack0 = st
		} else {
			newEnv.runtime = e.runtime.Push(st)
		}
		if newEnv.Global == nil && c.EnvgNeeded {
			internal.Panic("native function %q requires global env", c.Name)
		}
		c.Native(&newEnv)
		return newEnv.A
	}

	if c.Variadic {
		s := *newEnv.stack
		if len(s) > int(c.NumParams)-1 {
			s[c.NumParams-1] = newArray(append([]Value{}, s[c.NumParams-1:]...)...).ToValue()
		} else {
			newEnv.grow(int(c.NumParams))
			newEnv._set(uint16(c.NumParams)-1, newArray().ToValue())
		}
	}
	newEnv.growZero(int(c.StackSize), int(c.NumParams))

	var stk []Stacktrace
	if e != nil {
		stk = e.runtime.Stacktrace()
	}
	return internalExecCursorLoop(newEnv, m, stk)
}

func (c *funcbody) String() string {
	p := bytes.NewBufferString(c.Name)
	p.WriteString(internal.IfStr(c.Method, "([this],", "("))
	if c.Native != nil {
		p.WriteString("...")
	} else {
		for i := 0; i < int(c.NumParams); i++ {
			fmt.Fprintf(p, "a%d,", i)
		}
	}
	internal.CloseBuffer(p, internal.IfStr(c.Variadic, "...)", ")"))
	return p.String()
}

func (c *funcbody) GoString() string {
	return pkPrettify(c, c.LoadGlobal, false)
}

func pkPrettify(c *funcbody, p *Program, toplevel bool) string {
	sb := &bytes.Buffer{}
	sb.WriteString("+ START " + c.String() + "\n")
	if c.Native != nil {
		sb.WriteString("| NATIVE CODE\n")
	}

	readAddr := func(a uint16, rValue bool) string {
		if a == typ.RegA {
			return "$a"
		}

		suffix := ""
		if rValue {
			if a > typ.RegLocalMask || toplevel {
				x := (*p.stack)[a&typ.RegLocalMask]
				if x != Nil {
					suffix = ":" + detail(x)
				}
			}
		}

		if a > typ.RegLocalMask {
			return fmt.Sprintf("g$%d", a&typ.RegLocalMask) + suffix
		}
		return fmt.Sprintf("$%d", a&typ.RegLocalMask) + suffix
	}

	oldpos := c.CodeSeg.Pos

	for i, inst := range c.CodeSeg.Code {
		cursor := uint32(i) + 1
		bop, a, b, c := inst.Opcode, inst.A, inst.B, inst.C

		if oldpos.Len() > 0 {
			_, op, line := oldpos.Read(0)
			// log.Println(cursor, splitInst, unsafe.Pointer(&Pos))
			for uint32(cursor) > op && oldpos.Len() > 0 {
				op, line = oldpos.Pop()
			}
			x := strconv.Itoa(int(line))
			sb.WriteString(fmt.Sprintf("|%-5s % 4d| ", x+"L", cursor-1))
		} else {
			sb.WriteString(fmt.Sprintf("|$     % 4d| ", cursor-1))
		}

		switch bop {
		case typ.OpSet:
			sb.WriteString(readAddr(a, false) + " = " + readAddr(b, true))
		case typ.OpCreateArray:
			sb.WriteString("createarray")
		case typ.OpCreateObject:
			sb.WriteString("createobject")
		case typ.OpLoadFunc:
			cls := p.functions[a]
			sb.WriteString("loadfunc " + cls.fun.Name + "\n")
			sb.WriteString(pkPrettify(cls.fun, cls.fun.LoadGlobal, false))
		case typ.OpTailCall, typ.OpCall, typ.OpTryCall:
			if b != typ.RegPhantom {
				sb.WriteString("push " + readAddr(b, true) + " -> ")
			}
			switch bop {
			case typ.OpTailCall:
				sb.WriteString("tail")
			case typ.OpTryCall:
				sb.WriteString("try")
			}
			sb.WriteString("call " + readAddr(a, true))
		case typ.OpIfNot, typ.OpJmp:
			dest := inst.D()
			pos2 := uint32(int32(cursor) + dest)
			if bop == typ.OpIfNot {
				sb.WriteString("if not $a ")
			}
			sb.WriteString(fmt.Sprintf("jmp %d to %d", dest, pos2))
		case typ.OpInc:
			sb.WriteString("inc " + readAddr(a, false) + " " + readAddr(b, true))
			if c != 0 {
				sb.WriteString(fmt.Sprintf(" jmp %d to %d", int16(c), int32(cursor)+int32(int16(c))))
			}
		default:
			if bop == typ.OpLoad {
				sb.WriteString("load " + readAddr(a, true) + " " + readAddr(b, true) + " -> " + readAddr(c, false))
			} else if bop == typ.OpStore {
				sb.WriteString("store " + readAddr(a, true) + " " + readAddr(b, true) + " " + readAddr(c, true))
			} else if us, ok := typ.UnaryOpcode[bop]; ok {
				sb.WriteString(us + " " + readAddr(a, true))
			} else if bs, ok := typ.BinaryOpcode[bop]; ok {
				sb.WriteString(bs + " " + readAddr(a, true) + " " + readAddr(b, true))
			} else {
				sb.WriteString(fmt.Sprintf("? %02x", bop))
			}
		}

		sb.WriteString("\n")
	}

	sb.WriteString("+ END " + c.String())
	return sb.String()
}
