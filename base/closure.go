package base

import (
	"fmt"
)

type Closure struct {
	code      []byte
	env       *Env
	argsCount int
	preArgs   []Value
}

func NewClosure(code []byte, env *Env, argsCount int) *Closure {
	return &Closure{
		code:      code,
		env:       env,
		argsCount: argsCount,
	}
}

func (c *Closure) AppendPreArgs(preArgs []Value) {
	if c.preArgs == nil {
		c.preArgs = make([]Value, 0, 4)
	}

	c.preArgs = append(c.preArgs, preArgs...)
	c.argsCount -= len(preArgs)
	if c.argsCount < 0 {
		panic("negative args count")
	}
}

func (c *Closure) PreArgs() []Value {
	return c.preArgs
}

func (c *Closure) SetCode(code []byte) {
	c.code = code
}

func (c *Closure) Code() []byte {
	return c.code
}

func (c *Closure) ArgsCount() int {
	return c.argsCount
}

func (c *Closure) Env() *Env {
	return c.env
}

func (c *Closure) Dup() *Closure {
	cls := NewClosure(c.code, c.env, c.argsCount)
	//     cls.preArgs = preArgs == null ? null : preArgs.clone();
	if c.preArgs != nil {
		cls.preArgs = make([]Value, len(c.preArgs))
		copy(cls.preArgs, c.preArgs)
	}
	return cls
}

func (c *Closure) String() string {
	return fmt.Sprintf("lambda %d (\n", c.argsCount) + NewBytesReader(c.code).Prettify(4) + ")"
}
