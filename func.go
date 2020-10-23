package script

import (
	"bytes"
	"fmt"
	"strconv"
)

const (
	FuncYield = 1 << iota
	FuncVararg
	FuncLocked
)

type Func struct {
	packet
	Name       string
	ConstTable []Value
	NumParam   byte
	options    byte
	stackSize  uint16
	yCursor    uint32
	yEnv       Env
	native     func(env *Env)
	loadGlobal *Global
}

// Native creates a golang-native function
func Native(f func(env *Env)) Value {
	return Function(&Func{native: f})
}

func (c *Func) setOpt(flag bool, opt byte) {
	if flag {
		c.options |= opt
	} else {
		c.options &^= opt
	}
}

func (c *Func) Is(opt byte) bool { return (c.options & opt) > 0 }

func (c *Func) String() string {
	if c.native != nil {
		return "<native>"
	}

	p := bytes.Buffer{}
	if c.Name != "" {
		p.WriteString(c.Name)
	} else {
		p.WriteString("function")
	}

	p.WriteString("$")
	p.WriteString(strconv.Itoa(len(c.ConstTable)))
	p.WriteString("(")
	for i := 0; i < int(c.NumParam); i++ {
		p.WriteString("a")
		p.WriteString(strconv.Itoa(i))
		p.WriteString(",")
	}

	if c.Is(FuncVararg) {
		p.WriteString("...")
	} else {
		if c.NumParam > 0 {
			p.Truncate(p.Len() - 1)
		}
	}
	p.WriteString(")")

	if c.yEnv.stack != nil {
		p.WriteString(fmt.Sprintf("@%x", c.yCursor))
	}

	return p.String()
}

func (c *Func) PrettyString() string {
	if c.native != nil {
		return "[native Code]"
	}
	return pkPrettify(c, 0)
}

// exec executes the closure with the given Env
func (c *Func) exec(newEnv Env) (Value, []Value) {
	if c.native != nil {
		c.native(&newEnv)
		return newEnv.A, newEnv.V
	}

	if c.Is(FuncLocked) {
		panicf("reenter yielded function")
	}
	if c.yEnv.stack != nil {
		newEnv = c.yEnv
	}
	if c.Is(FuncYield) {
		c.setOpt(true, FuncLocked)
	}

	v, vb, np, yield := execCursorLoop(newEnv, c, c.yCursor)
	c.setOpt(false, FuncLocked)

	if yield {
		c.yCursor, c.yEnv = np, newEnv
	} else {
		c.yCursor, c.yEnv = 0, Env{}
	}
	return v, vb
}

func (c *Func) Call(a ...Value) (Value, []Value) {
	return c.CallEnv(nil, a...)
}

func (c *Func) PCall(a ...Value) (v1 Value, v []Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err, _ = r.(error)
			if err == nil {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	v1, v = c.CallEnv(nil, a...)
	return
}

func (c *Func) CallEnv(env *Env, a ...Value) (Value, []Value) {
	if c.yEnv.stack != nil {
		return c.exec(c.yEnv)
	}

	var newEnv Env
	var varg []Value
	if env == nil || env.global == nil {
		// panicf("call function without global env")
		x := make([]Value, 0)
		newEnv.stack = &x
		newEnv.global = c.loadGlobal
	} else {
		newEnv = *env
		newEnv.stackOffset = len(*newEnv.stack)
	}

	for i := range a {
		if i >= int(c.NumParam) {
			varg = append(varg, a[i])
		}
		newEnv.Push(a[i])
	}

	if c.native == nil {
		if c.Is(FuncVararg) {
			newEnv.grow(int(c.NumParam) + 1)
			newEnv._set(uint16(c.NumParam), _unpackedStack(&unpacked{a: varg}))
		}
		if c.Is(FuncYield) {
			x := append([]Value{}, newEnv.Stack()...)
			newEnv.stack = &x
			newEnv.stackOffset = 0
		}
		newEnv.grow(int(c.stackSize))
	}
	return c.exec(newEnv)
}

// Terminate will try to stop the execution, when called the closure (along with duplicates) become invalid immediately
func (c *Func) Terminate() {
	const Stop = uint32(OpEOB) << 26
	for i := range c.Code {
		c.Code[i] = Stop
	}
}
