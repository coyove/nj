package bas

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

var objEmptyFunc = &funcbody{Name: "object"}

type funcbody struct {
	CodeSeg    internal.Packet
	StackSize  uint16
	NumParams  byte
	Variadic   bool
	Method     bool
	Native     func(env *Env)
	Name       string
	LoadGlobal *Program
	Locals     []string
}

type Program struct {
	top       *Object
	symbols   *Object
	stack     *[]Value
	functions []*Object
	stopped   bool
	file      string
	source    string

	MaxStackSize int64
	Globals      *Object
	Stdout       io.Writer
	Stderr       io.Writer
	Stdin        io.Reader
}

func NewProgram(file, source string, coreStack *Env, top, symbols *Object, funcs []*Object) *Program {
	cls := &Program{top: top}
	cls.stack = coreStack.stack
	cls.symbols = symbols
	cls.functions = funcs
	cls.file = file
	cls.source = source
	cls.Stdout = os.Stdout
	cls.Stdin = os.Stdin
	cls.Stderr = os.Stderr

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

func (p *Program) Run() (v1 Value, err error) {
	p.stopped = false
	defer internal.CatchError(&err)
	newEnv := Env{
		Global: p,
		stack:  p.stack,
	}
	v1 = internalExecCursorLoop(newEnv, p.top, nil)
	return
}

func (p *Program) Stop() {
	p.stopped = true
	return
}

func (p *Program) FileSource() (string, string) {
	return p.file, p.source
}

func (p *Program) GoString() string {
	return p.top.GoString()
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
	e2.runtime.stackN = append([]Stacktrace{}, e2.runtime.stackN...)
	return &e2
}

func Call(m *Object, args ...Value) (res Value) {
	return CallObject(m, Runtime{}, nil, m.this, args...)
}

func Call2(m *Object, args ...Value) (res Value, err error) {
	res = CallObject(m, Runtime{}, &err, m.this, args...)
	return
}

func CallObject(m *Object, r Runtime, outErr *error, this Value, args ...Value) (res Value) {
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
		defer relayPanic(func() []Stacktrace { return newEnv.runtime.Stacktrace() })
		newEnv.runtime = r.push(Stacktrace{
			Callable:        m,
			stackOffsetFlag: internal.FlagNativeCall,
		})
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

	return internalExecCursorLoop(newEnv, m, r.Stacktrace())
}

func (o *Object) funcSig() string {
	c := o.fun
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

func (obj *Object) GoString() string {
	buf := &bytes.Buffer{}
	obj.printAll(buf, true)
	return buf.String()
}

func (obj *Object) printAll(w io.Writer, toplevel bool) {
	defer func() {
		internal.WriteString(w, "keyvalue)\n")
		ki := 0
		obj.Foreach(func(k Value, v *Value) bool {
			fmt.Fprintf(w, "%d)\t%v\t%v\n", ki, k, *v)
			ki++
			return true
		})
		internal.WriteString(w, "endkeyvalue)\n")
	}()

	c, p := obj.fun, obj.fun.LoadGlobal
	if c.Native != nil {
		internal.WriteString(w, "0)\t0\tnative code\n")
		return
	}

	internal.WriteString(w, "start)\t"+obj.funcSig()+"\n")

	if c == p.top.fun {
		internal.WriteString(w, "source)\t"+c.LoadGlobal.file+"\n")
	}

	readAddr := func(a uint16, rValue bool) string {
		if a == typ.RegA {
			return "a"
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
			return fmt.Sprintf("g(%d)", a&typ.RegLocalMask) + suffix
		}
		return fmt.Sprintf("sp(%d)", a&typ.RegLocalMask) + suffix
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
			internal.WriteString(w, fmt.Sprintf("%d)\t%d\t", line, cursor-1))
		} else {
			internal.WriteString(w, fmt.Sprintf("$)\t%d\t", cursor-1))
		}

		switch bop {
		case typ.OpSet:
			internal.WriteString(w, readAddr(a, false)+" = "+readAddr(b, false))
		case typ.OpCreateArray:
			internal.WriteString(w, "createarray")
		case typ.OpCreateObject:
			internal.WriteString(w, "createobject")
		case typ.OpLoadFunc:
			if a == typ.RegA {
				internal.WriteString(w, "loadself")
			} else {
				cls := p.functions[a]
				internal.WriteString(w, fmt.Sprintf("loadfunc(%d)\n", a))
				cls.printAll(w, false)
			}
		case typ.OpTailCall, typ.OpCall, typ.OpTryCall:
			if b != typ.RegPhantom {
				internal.WriteString(w, "push "+readAddr(b, false)+" -> ")
			}
			switch bop {
			case typ.OpTailCall:
				internal.WriteString(w, "tail")
			case typ.OpTryCall:
				internal.WriteString(w, "try")
			}
			internal.WriteString(w, "call "+readAddr(a, true))
		case typ.OpIfNot, typ.OpJmp:
			dest := inst.D()
			pos2 := uint32(int32(cursor) + dest)
			if bop == typ.OpIfNot {
				internal.WriteString(w, "if not $a ")
			}
			internal.WriteString(w, fmt.Sprintf("jmp %d to %d", dest, pos2))
		case typ.OpInc:
			internal.WriteString(w, "inc "+readAddr(a, false)+" "+readAddr(b, false))
			if c != 0 {
				internal.WriteString(w, fmt.Sprintf(" jmp %d to %d", int16(c), int32(cursor)+int32(int16(c))))
			}
		default:
			if bop == typ.OpLoad {
				internal.WriteString(w, "load "+readAddr(a, false)+" "+readAddr(b, false)+" -> "+readAddr(c, false))
			} else if bop == typ.OpStore {
				internal.WriteString(w, "store "+readAddr(a, false)+" "+readAddr(b, false)+" <- "+readAddr(c, false))
			} else if us, ok := typ.UnaryOpcode[bop]; ok {
				internal.WriteString(w, us+" "+readAddr(a, false))
			} else if bs, ok := typ.BinaryOpcode[bop]; ok {
				internal.WriteString(w, bs+" "+readAddr(a, false)+" "+readAddr(b, false))
			} else {
				internal.WriteString(w, fmt.Sprintf("? %02x", bop))
			}
		}

		internal.WriteString(w, "\n")
	}

	internal.WriteString(w, "end)\t"+obj.funcSig()+"\n")
}
