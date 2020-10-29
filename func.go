package script

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"time"
)

type Func struct {
	code       packet
	name, doc  string
	numParams  byte
	isVariadic bool
	stackSize  uint16
	native     func(env *Env)
	loadGlobal *Program
	params     []string
}

type Program struct {
	Func
	Deadline              int64
	MaxStackSize          int64
	MaxCallStackSize      int64
	Extras                map[string]interface{}
	Stack                 *[]Value
	Funcs                 []*Func
	Stdout, Stdin, Stderr io.ReadWriter
	NilIndex              uint16
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
func Native(f func(env *Env), doc ...string) Value {
	return Function(&Func{native: f, doc: strings.Join(doc, "\n")})
}

func NativeWithParamMap(f func(*Env, Arguments), doc string, params ...string) Value {
	return Function(&Func{
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
	defer catchErr(&err)
	newEnv := Env{
		Global: p,
		stack:  p.Stack,
	}
	v1, v = execCursorLoop(newEnv, &p.Func, 0)
	return
}

func (c *Func) Call(env *Env, a ...Value) (v1 Value, v []Value, err error) {
	defer catchErr(&err)

	var newEnv Env
	var varg []Value
	if env == nil || env.Global == nil {
		// panicf("call function without Global env")
		x := make([]Value, 0)
		newEnv.stack = &x
		newEnv.Global = c.loadGlobal
	} else {
		newEnv = *env
		newEnv.stackOffset = uint32(len(*newEnv.stack))
	}

	for i := range a {
		if i >= int(c.numParams) {
			varg = append(varg, a[i])
		}
		newEnv.Push(a[i])
	}

	if c.native == nil {
		newEnv.grow(int(c.stackSize))
		if c.isVariadic {
			// newEnv.grow(int(c.numParams) + 1)
			newEnv._set(uint16(c.numParams), _unpackedStack(&unpacked{a: varg}))
		}
	}

	v1, v = c.exec(newEnv)
	return
}

// Terminate will try to stop the execution, when called the closure (along with duplicates) become invalid immediately
// func (c *Func) Terminate() {
// 	const Stop = uint32(OpRet) << 26
// 	for i := range c.code.Code {
// 		c.code.Code[i] = Stop
// 	}
// }

func (p *Program) PrettyCode() string { return pkPrettify(&p.Func, p, true, 0) }

func (p *Program) SetTimeout(d time.Duration) { p.Deadline = time.Now().Add(d).UnixNano() }

func (p *Program) SetDeadline(d time.Time) { p.Deadline = d.UnixNano() }

func (p *Program) AddValue(k string, v interface{}) *Program {
	if p.Extras == nil {
		p.Extras = map[string]interface{}{}
	}
	p.Extras[k] = v
	return p
}
