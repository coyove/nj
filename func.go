package script

import (
	"bytes"
	"strconv"
	"time"
)

type Func struct {
	code       packet
	name       string
	constTable []Value
	numParams  byte
	isVariadic bool
	stackSize  uint16
	native     func(env *Env)
	loadGlobal *Program
}

type Program struct {
	Func
	Deadline     int64
	MaxStackSize int64
	Extras       map[string]interface{}
	Stack        *[]Value
	Funcs        []*Func
}

// Native creates a golang-native function
func Native(f func(env *Env)) Value { return Function(&Func{native: f}) }

func (c *Func) Name() string { return c.name }

func (c *Func) IsNative() bool { return c.native != nil }

func (c *Func) Signature() (numParams int, isVariadic bool, stackSize int) {
	return int(c.numParams), c.isVariadic, int(c.stackSize)
}

func (c *Func) String() string {
	if c.native != nil {
		return "<native>"
	}

	p := bytes.Buffer{}
	if c.name != "" {
		p.WriteString(c.name)
	} else {
		p.WriteString("function")
	}

	p.WriteString("$")
	p.WriteString(strconv.Itoa(len(c.constTable)))
	p.WriteString("(")
	for i := 0; i < int(c.numParams); i++ {
		p.WriteString("a")
		p.WriteString(strconv.Itoa(i))
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
	return pkPrettify(c, nil, 0)
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
		global: p,
		stack:  p.Stack,
	}
	v1, v = execCursorLoop(newEnv, &p.Func, 0)
	return
}

func (c *Func) Call(env *Env, a ...Value) (v1 Value, v []Value, err error) {
	defer catchErr(&err)

	var newEnv Env
	var varg []Value
	if env == nil || env.global == nil {
		// panicf("call function without global env")
		x := make([]Value, 0)
		newEnv.stack = &x
		newEnv.global = c.loadGlobal
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
func (c *Func) Terminate() {
	const Stop = uint32(OpEOB) << 26
	for i := range c.code.Code {
		c.code.Code[i] = Stop
	}
}

func (p *Program) PrettyCode() string { return pkPrettify(&p.Func, p, 0) }

func (p *Program) SetTimeout(d time.Duration) { p.Deadline = time.Now().Add(d).Unix() }

func (p *Program) SetDeadline(d time.Time) { p.Deadline = d.Unix() }

func (p *Program) AddValue(k string, v interface{}) *Program {
	if p.Extras == nil {
		p.Extras = map[string]interface{}{}
	}
	p.Extras[k] = v
	return p
}
