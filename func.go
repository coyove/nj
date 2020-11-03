package script

import (
	"bytes"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/coyove/script/parser"
)

type Func struct {
	code        packet
	name, doc   string
	numParams   byte
	isVariadic  bool
	stackSize   uint16
	native      func(env *Env)
	loadGlobal  *Program
	params      []string
	debugLocals []string
}

type Program struct {
	Func
	Deadline         int64
	MaxStackSize     int64
	MaxCallStackSize int64
	MaxStringSize    int64
	Stack            *[]Value
	Funcs            []*Func
	Stdout, Stderr   io.Writer
	Stdin            io.Reader
	Logger           *log.Logger
	NilIndex         uint16
	Get              func(string) (Value, error)
	Set              func(string, Value) error
	Survey           struct {
		TotalStringAlloc int64
	}
}

type Arguments map[string]Value

func (a Arguments) GetString(name string, defaultValue string) string {
	if a[name].Type() == VString {
		return a[name].String()
	}
	return defaultValue
}

func (a Arguments) GetInt(name string, defaultValue int64) int64 {
	if a[name].Type() == VNumber {
		return a[name].Int()
	}
	return defaultValue
}

// Native creates a golang-native function
func Native(name string, f func(env *Env), doc ...string) Value {
	return Function(&Func{name: name, native: f, doc: strings.Join(doc, "\n")})
}

func NativeWithParamMap(name string, f func(*Env, Arguments), doc string, params ...string) Value {
	return Function(&Func{
		name:      name,
		params:    params,
		numParams: byte(len(params)),
		doc:       doc,
		native: func(env *Env) {
			stack := env.Stack()
			args := make(map[string]Value, len(stack))
			for i := range stack {
				if i < len(params) {
					args[params[i]] = stack[i]
				}
			}
			f(env, Arguments(args))
		},
	})
}

func (c *Func) Name() string { return c.name }

func (c *Func) IsNative() bool { return c.native != nil }

func (c *Func) Signature() (numParams int, isVariadic bool, stackSize int) {
	return int(c.numParams), c.isVariadic, int(c.stackSize)
}

func (c *Func) String() string {
	p := bytes.Buffer{}
	if c.name != "" {
		p.WriteString(c.name)
	} else if c.native != nil {
		p.WriteString("native")
	} else {
		p.WriteString("function")
	}

	p.WriteString("(")
	for i := 0; i < int(c.numParams); i++ {
		if i < len(c.params) {
			p.WriteString(c.params[i])
		} else {
			p.WriteString("a" + strconv.Itoa(i))
		}
		p.WriteString(",")
	}
	if c.isVariadic {
		p.WriteString("...")
	} else {
		if c.numParams > 0 {
			p.Truncate(p.Len() - 1)
		}
	}
	p.WriteString(")")
	return p.String()
}

func (c *Func) PrettyCode() string {
	if c.native != nil {
		return "[native code]"
	}
	return pkPrettify(c, c.loadGlobal, false, 0)
}

func (c *Func) exec(newEnv Env) (Value, []Value) {
	if c.native != nil {
		c.native(&newEnv)
		return newEnv.A, newEnv.V
	}
	return execCursorLoop(newEnv, c, 0)
}

func (p *Program) Run() (v1 Value, v []Value, err error) {
	return p.Call()
}

func (p *Program) Call() (v1 Value, v []Value, err error) {
	defer parser.CatchError(&err)
	newEnv := Env{
		Global: p,
		stack:  p.Stack,
	}
	v1, v = execCursorLoop(newEnv, &p.Func, 0)
	return
}

func (c *Func) Call(a ...Value) (v1 Value, v []Value, err error) {
	defer parser.CatchError(&err)

	oldLen := len(*c.loadGlobal.Stack)
	newEnv := Env{
		Global:      c.loadGlobal,
		stack:       c.loadGlobal.Stack,
		stackOffset: uint32(oldLen),
	}

	var varg []Value
	for i := range a {
		if i >= int(c.numParams) {
			varg = append(varg, a[i])
		}
		newEnv.Push(a[i])
	}

	if c.native == nil {
		newEnv.growZero(int(c.stackSize))
		if c.isVariadic {
			// newEnv.grow(int(c.numParams) + 1)
			newEnv._set(uint16(c.numParams), _unpackedStack(varg))
		}
	}

	v1, v = c.exec(newEnv)
	*c.loadGlobal.Stack = (*c.loadGlobal.Stack)[:oldLen]
	return
}

func (p *Program) PrettyCode() string { return pkPrettify(&p.Func, p, true, 0) }

func (p *Program) SetTimeout(d time.Duration) { p.Deadline = time.Now().Add(d).UnixNano() }

func (p *Program) SetDeadline(d time.Time) { p.Deadline = d.UnixNano() }

func (p *Program) Print(a ...interface{}) { p.log("", "", a...) }

func (p *Program) Printf(f string, a ...interface{}) { p.log("f", f, a...) }

func (p *Program) Println(a ...interface{}) { p.log("l", "", a...) }

func (p *Program) Fatal(a ...interface{}) { p.log("F", "", a...) }

func (p *Program) Fatalf(f string, a ...interface{}) { p.log("Ff", f, a...) }

func (p *Program) Fatalln(a ...interface{}) { p.log("Fl", "", a...) }

func (p *Program) Panic(a ...interface{}) { p.log("P", "", a...) }

func (p *Program) Panicf(f string, a ...interface{}) { p.log("Pf", f, a...) }

func (p *Program) Panicln(a ...interface{}) { p.log("Pl", "", a...) }

func (p *Program) log(o, f string, a ...interface{}) {
	if p.Logger == nil {
		p.Logger = log.New(p.Stderr, "", log.LstdFlags)
	}
	switch o {
	default:
		p.Logger.Print(a...)
	case "f":
		p.Logger.Printf(f, a...)
	case "l":
		p.Logger.Println(a...)
	case "F":
		p.Logger.Fatal(a...)
	case "Ff":
		p.Logger.Fatalf(f, a...)
	case "Fl":
		p.Logger.Fatalln(a...)
	case "P":
		p.Logger.Panic(a...)
	case "Pf":
		p.Logger.Panicf(f, a...)
	case "Pl":
		p.Logger.Panicln(a...)
	}
}
