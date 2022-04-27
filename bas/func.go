package bas

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

type Function struct {
	CodeSeg    internal.Packet
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
	Deadline     time.Time
}

type Program struct {
	top       *Function
	symbols   *Object
	stack     *[]Value
	functions []*Object
	Environment
}

func NewProgram(coreStack *Env, top *Function, symbols *Object, funcs []*Object, env *Environment) *Program {
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

func (p *Program) Run() (Value, error) {
	if p.Deadline.IsZero() {
		return p.runTop()
	}
	if time.Now().After(p.Deadline) {
		return Nil, fmt.Errorf("timeout")
	}

	dummy := make(map[*Function]*internal.Packet)
	dummy[p.top] = p.top.CodeSeg.Copy()
	for _, f := range p.functions {
		dummy[f.fun] = f.fun.CodeSeg.Copy()
	}

	var mu sync.Mutex
	var reverted, timedout bool
	finished := make(chan struct{}, 1)
	go func(start time.Time) {
		select {
		case <-finished:
		case <-time.After(time.Until(p.Deadline)):
			mu.Lock()
			if !reverted { // if code has been reverted to original, don't Stop again
				p.Stop()
				timedout = true
			}
			mu.Unlock()
		}
	}(time.Now())
	v, err := p.runTop()
	finished <- struct{}{}

	// Revert code anyway
	mu.Lock()
	p.top.CodeSeg = *dummy[p.top]
	for _, f := range p.functions {
		f.fun.CodeSeg = *dummy[f.fun]
	}
	reverted = true
	mu.Unlock()

	if timedout {
		return Nil, fmt.Errorf("timeout")
	}
	return v, err
}

func (p *Program) runTop() (v1 Value, err error) {
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

func (c *Function) Object() *Object {
	return c.obj
}

func pkPrettify(c *Function, p *Program, toplevel bool) string {
	sb := &bytes.Buffer{}
	sb.WriteString("+ START " + c.String() + "\n")

	readAddr := func(a uint16, rValue bool) string {
		if a == typ.RegA {
			return "$a"
		}

		suffix := ""
		if rValue {
			if a > typ.RegLocalMask || toplevel {
				suffix = ":" + simpleString((*p.stack)[a&typ.RegLocalMask])
			}
		}

		if a > typ.RegLocalMask {
			return fmt.Sprintf("g$%d", a&typ.RegLocalMask) + suffix
		}
		return fmt.Sprintf("$%d", a&typ.RegLocalMask) + suffix
	}

	oldpos := c.CodeSeg.Pos
	lastLine := uint32(0)

	for i, inst := range c.CodeSeg.Code {
		cursor := uint32(i) + 1
		bop, a, b := inst.Opcode, inst.A, uint16(inst.B)

		if c.CodeSeg.Pos.Len() > 0 {
			op, line := c.CodeSeg.Pos.Pop()
			// log.Println(cursor, splitInst, unsafe.Pointer(&Pos))
			for uint32(cursor) > op && c.CodeSeg.Pos.Len() > 0 {
				if op, line = c.CodeSeg.Pos.Pop(); uint32(cursor) <= op {
					break
				}
			}

			if op == uint32(cursor) {
				x := "."
				if line != lastLine {
					x = strconv.Itoa(int(line))
					lastLine = line
				}
				sb.WriteString(fmt.Sprintf("|%-4s % 4d| ", x, cursor-1))
			} else {
				sb.WriteString(fmt.Sprintf("|     % 4d| ", cursor-1))
			}
		} else {
			sb.WriteString(fmt.Sprintf("|$    % 4d| ", cursor-1))
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
		case typ.OpTailCall, typ.OpCall:
			if b != typ.RegPhantom {
				sb.WriteString("push " + readAddr(b, true) + " -> ")
			}
			if bop == typ.OpTailCall {
				sb.WriteString("tail")
			}
			sb.WriteString("call " + readAddr(a, true))
		case typ.OpIfNot, typ.OpJmp:
			pos := inst.B
			pos2 := uint32(int32(cursor) + pos)
			if bop == typ.OpIfNot {
				sb.WriteString("if not $a ")
			}
			sb.WriteString(fmt.Sprintf("jmp %d to %d", pos, pos2))
		case typ.OpInc:
			sb.WriteString("inc " + readAddr(a, false) + " " + readAddr(b, true))
		default:
			if us, ok := typ.UnaryOpcode[bop]; ok {
				sb.WriteString(us + " " + readAddr(a, true))
			} else if bs, ok := typ.BinaryOpcode[bop]; ok {
				sb.WriteString(bs + " " + readAddr(a, true) + " " + readAddr(b, true))
			} else {
				sb.WriteString(fmt.Sprintf("? %02x", bop))
			}
		}

		sb.WriteString("\n")
	}

	c.CodeSeg.Pos = oldpos

	sb.WriteString("+ END " + c.String())
	return sb.String()
}
