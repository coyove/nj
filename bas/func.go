package bas

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

type funcbody struct {
	name      string
	codeSeg   internal.Packet
	stackSize uint16
	numArgs   byte
	varg      bool
	method    bool
	native    func(env *Env)
	top       *Program
	locals    []string
	caps      []string
}

type Program struct {
	stack     *[]Value
	main      *Object
	symbols   *Map
	functions *Map
	stopped   bool

	File         string
	Source       string
	MaxStackSize int64
	Globals      Map
	Stdout       io.Writer
	Stderr       io.Writer
	Stdin        io.Reader
}

var globals struct {
	sym   Map
	store Map
	stack []Value
}

func GetGlobalName(v Value) int {
	x, ok := globals.sym.Get(v)
	if !ok {
		return 0
	}
	return int(x.UnsafeInt64())

}

func Globals() Map {
	return globals.store.Copy()
}

func AddGlobal(k string, v Value) {
	if len(globals.stack) == 0 {
		globals.stack = append(globals.stack, Nil)
	}
	sk := Str(k)
	idx, ok := globals.sym.Get(sk)
	if ok {
		globals.stack[idx.Int()] = v
	} else {
		idx := len(globals.stack)
		globals.sym.Set(sk, Int(idx))
		globals.stack = append(globals.stack, v)
	}
	globals.store.Set(sk, v)
}

func AddGlobalFunc(k string, f func(*Env)) {
	AddGlobal(k, Func(k, f))
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
	obj.fun = &funcbody{name: name, native: f}
	obj.SetPrototype(&Proto.Func)
	return obj.ToValue()
}

func (p *Program) Run() (v1 Value, err error) {
	p.stopped = false
	if p.MaxStackSize <= 0 {
		p.MaxStackSize = math.MaxInt64
	}

	defer internal.CatchError(&err)
	newEnv := Env{
		top:   p,
		stack: p.stack,
	}
	v1 = internalExecCursorLoop(newEnv, p.main, nil)
	return
}

// Stop stops the program from running unsafely, it can't stop any Go-native functions or goroutines.
func (p *Program) Stop() {
	p.stopped = true
	return
}

func (p *Program) GoString() string {
	x := &bytes.Buffer{}
	p.main.printAll(x)
	p.functions.Foreach(func(f Value, idx *Value) bool {
		x.WriteByte('\n')
		(*p.stack)[idx.Int()&typ.RegLocalMask].Object().printAll(x)
		return true
	})
	return x.String()
}

func (p *Program) Get(k string) (v Value, ok bool) {
	addr, ok := p.symbols.Get(Str(k))
	if !ok {
		return Nil, false
	}
	return (*p.stack)[addr.Int64()], true
}

func (p *Program) Set(k string, v Value) (ok bool) {
	addr, ok := p.symbols.Get(Str(k))
	if !ok {
		return false
	}
	(*p.stack)[addr.Int64()] = v
	return true
}

func (p *Program) LocalsObject() *Object {
	r := NewObject(len(p.main.fun.locals))
	for i, name := range p.main.fun.locals {
		r.Set(Str(name), (*p.stack)[i])
	}
	return r
}

func (m *Object) Apply(e *Env, this Value, args ...Value) Value {
	if e != nil {
		return callobj(m, e.runtime, e.top, nil, this, args...)
	}
	return callobj(m, stacktraces{}, nil, nil, this, args...)
}

func (m *Object) Call(e *Env, args ...Value) (res Value) {
	if e != nil {
		return callobj(m, e.runtime, e.top, nil, m.this, args...)
	}
	return callobj(m, stacktraces{}, nil, nil, m.this, args...)
}

func (m *Object) TryCall(e *Env, args ...Value) (res Value, err error) {
	if e != nil {
		res = callobj(m, e.runtime, e.top, &err, m.this, args...)
	} else {
		res = callobj(m, stacktraces{}, nil, &err, m.this, args...)
	}
	return
}

func callobj(m *Object, r stacktraces, g *Program, outErr *error, this Value, args ...Value) (res Value) {
	c := m.fun
	newEnv := Env{
		A:     this,
		top:   c.top,
		stack: &args,
	}

	if c.top == nil {
		newEnv.top = g
	}

	if outErr != nil {
		defer internal.CatchError(outErr)
	}

	if c.native != nil {
		defer relayPanic(func() []Stacktrace { return newEnv.runtime.Stacktrace(false) })
		newEnv.runtime = r.push(Stacktrace{
			Callable: m,
		})
		c.native(&newEnv)
		return newEnv.A
	}

	if c.varg {
		s := *newEnv.stack
		if len(s) > int(c.numArgs)-1 {
			s[c.numArgs-1] = newVarargArray(s[c.numArgs-1:]).ToValue()
		} else {
			if newEnv.Size() < int(c.numArgs)-1 {
				internal.PanicNotEnoughArgs(m.ToValue().simple())
			}
			newEnv.resize(int(c.numArgs))
			newEnv._set(uint16(c.numArgs)-1, Nil)
		}
	} else {
		if newEnv.Size() < int(c.numArgs) {
			internal.PanicNotEnoughArgs(m.ToValue().simple())
		}
	}
	newEnv.resizeZero(int(c.stackSize), int(c.numArgs))

	return internalExecCursorLoop(newEnv, m, r.Stacktrace(false))
}

func (o *Object) funcSig() string {
	c := o.fun
	p := bytes.NewBufferString(c.name)
	p.WriteString(internal.IfStr(c.method, "({this},", "("))
	if c.native != nil {
		p.WriteString("...")
	} else {
		for i := 0; i < int(c.numArgs); i++ {
			fmt.Fprintf(p, "a%d,", i)
		}
	}
	internal.CloseBuffer(p, internal.IfStr(c.varg, "...)", ")"))
	return p.String()
}

func (obj *Object) GoString() string {
	buf := &bytes.Buffer{}
	obj.printAll(buf)
	return buf.String()
}

func (obj *Object) printAll(w io.Writer) {
	cls, p := obj.fun, obj.fun.top
	internal.WriteString(w, "start)\t"+obj.funcSig()+"\n")
	if obj.parent != nil {
		internal.WriteString(w, "proto)\t"+obj.parent.Name()+"\n")
	}
	if cls.native != nil {
		internal.WriteString(w, "0)\t0\tnative code\n")
	} else {
		if cls == p.main.fun {
			internal.WriteString(w, "source)\t"+cls.top.File+"\n")
		}

		readAddr := func(a uint16, rValue bool) string {
			if a == typ.RegA {
				return "a"
			}

			suffix := ""
			if addr := a & typ.RegLocalMask; a != addr && rValue && int(addr) < len(*p.stack) {
				x := (*p.stack)[addr]
				if x != Nil {
					suffix = ":" + x.simple()
				}
			}

			if a > typ.RegLocalMask {
				return fmt.Sprintf("g%d", a&typ.RegLocalMask) + suffix
			}
			return fmt.Sprintf("sp+%d", a&typ.RegLocalMask) + suffix
		}

		oldpos := cls.codeSeg.Pos

		for i, inst := range cls.codeSeg.Code {
			cursor := uint32(i) + 1
			bop, a, b, c := inst.Opcode, inst.A, inst.B, inst.C

			if oldpos.Len() > 0 {
				c1 := cursor - 1
				_, op, line := oldpos.Read(0)
				for uint32(c1) > op && oldpos.Len() > 0 {
					op, line = oldpos.Pop()
				}
				internal.WriteString(w, fmt.Sprintf("%d)\t%d\t", line, c1))
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
			case typ.OpFunction:
				if a == typ.RegA {
					internal.WriteString(w, "moveself")
				} else if b == 0 {
					internal.WriteString(w, "copyfunction "+readAddr(a, true))
					internal.WriteString(w, " -> "+readAddr(c, false))
				} else if b == 1 {
					internal.WriteString(w, "copyclosure "+readAddr(a, true))
					internal.WriteString(w, " -> "+readAddr(c, false))
				}
			case typ.OpTailCall, typ.OpCall:
				if b != typ.RegPhantom {
					internal.WriteString(w, "push "+readAddr(b, false)+" -> ")
				}
				internal.WriteString(w, internal.IfStr(bop == typ.OpTailCall, "tailcall ", "call "))
				internal.WriteString(w, readAddr(a, true))
			case typ.OpIfNot, typ.OpJmp:
				dest := inst.D()
				pos2 := uint32(int32(cursor) + dest)
				if bop == typ.OpIfNot {
					internal.WriteString(w, "if not a ")
				}
				internal.WriteString(w, fmt.Sprintf("jmp %d to %d", dest, pos2))
			case typ.OpInc, typ.OpInc16:
				if bop == typ.OpInc16 {
					internal.WriteString(w, "inc16 "+readAddr(a, false)+" "+strconv.Itoa(int(int16(b))))
				} else {
					internal.WriteString(w, "inc "+readAddr(a, false)+" "+readAddr(b, false))
				}
				if c != 0 {
					internal.WriteString(w, fmt.Sprintf(" jmp %d to %d", int16(c), int32(cursor)+int32(int16(c))))
				}
			case typ.OpBitOp:
				internal.WriteString(w, "bit")
				internal.WriteString(w, [...]string{"and ", "or ", "xor ", "lsh ", "rsh ", "ursh "}[c])
				internal.WriteString(w, readAddr(a, false)+" "+readAddr(b, false))
			case typ.OpLoad:
				internal.WriteString(w, "load "+readAddr(a, false)+" "+readAddr(b, false)+" -> "+readAddr(c, false))
			case typ.OpStore:
				internal.WriteString(w, "store "+readAddr(a, false)+" "+readAddr(b, false)+" <- "+readAddr(c, false))
			case typ.OpSlice:
				internal.WriteString(w, "sliceload "+readAddr(a, false)+" "+readAddr(b, false)+" : "+readAddr(c, false))
			case typ.OpLoadGlobal:
				internal.WriteString(w, "loadglobal "+globals.stack[a].simple())
				if b != typ.RegPhantom {
					internal.WriteString(w, " "+readAddr(b, true))
				}
				internal.WriteString(w, " -> "+readAddr(c, false))
			case typ.OpLinear16:
				if b == 1 {
					internal.WriteString(w, "a = "+readAddr(a, false)+fmt.Sprintf(" + %d", int16(c)))
				} else if b == 65535 {
					internal.WriteString(w, "a = -"+readAddr(a, false)+fmt.Sprintf(" + %d", int16(c)))
				} else if c == 0 {
					internal.WriteString(w, "a = "+readAddr(a, false)+fmt.Sprintf(" * %d", int16(b)))
				} else {
					internal.WriteString(w, "a = "+readAddr(a, false)+fmt.Sprintf(" * %d + %d", int16(b), int16(c)))
				}
			case typ.OpCmp16:
				internal.WriteString(w, "a = "+readAddr(a, false)+fmt.Sprintf(" * %d < %d", int16(b), int16(c)))
			case typ.OpEq16:
				internal.WriteString(w, "a = "+readAddr(a, false)+internal.IfStr(c == typ.OpEq, " == ", " != ")+strconv.Itoa(int(int16(b))))
			default:
				if us, ok := typ.UnaryOpcode[bop]; ok {
					internal.WriteString(w, us+" "+readAddr(a, false))
				} else if bs, ok := typ.BinaryOpcode[bop]; ok {
					internal.WriteString(w, bs+" "+readAddr(a, false)+" "+readAddr(b, false))
				} else {
					internal.WriteString(w, fmt.Sprintf("? %02x", bop))
				}
			}

			internal.WriteString(w, "\n")
		}
	}
	ki := 0
	obj.Foreach(func(k Value, v *Value) bool {
		fmt.Fprintf(w, "prop%d)\t%v\t%v\n", ki, k, *v)
		ki++
		return true
	})
	internal.WriteString(w, "end)\t"+obj.funcSig())
}

func NewBareFunc(f string, varg bool, np byte, ss uint16, locals, caps []string, code internal.Packet) *Object {
	obj := NewObject(0)
	obj.SetPrototype(&Proto.Func)
	obj.fun = &funcbody{}
	obj.fun.varg = varg
	obj.fun.numArgs = np
	obj.fun.name = f
	obj.fun.stackSize = ss
	obj.fun.codeSeg = code
	obj.fun.locals = locals
	obj.fun.method = strings.Contains(f, ".")
	obj.fun.caps = caps
	return obj
}

func NewBareProgram(coreStack []Value, top *Object, symbols, funcs *Map) *Program {
	cls := &Program{}
	cls.main = top
	cls.stack = &coreStack
	cls.symbols = symbols
	cls.functions = funcs
	cls.Stdout = os.Stdout
	cls.Stdin = os.Stdin
	cls.Stderr = os.Stderr

	cls.main.fun.top = cls
	cls.functions.Foreach(func(f Value, idx *Value) bool {
		(*cls.stack)[idx.Int()&typ.RegLocalMask].Object().fun.top = cls
		return true
	})
	return cls
}
