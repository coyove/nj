package potatolang

import (
	"fmt"
)

const (
	ClsNoEnvescape = 1 << iota
	ClsYieldable
	ClsRecoverable
	ClsNative
	ClsVararg
)

type Closure struct {
	Code        []uint32
	Pos         posVByte
	source      []byte
	ConstTable  []Value
	Env         *Env
	ParamsCount byte
	options     byte
	lastp       uint32
	lastenv     *Env
	native      func(env *Env)
}

// NewClosure creates a new closure
func NewClosure(code []uint32, consts []Value, env *Env, paramsCount byte) *Closure {
	return &Closure{
		Code:        code,
		ConstTable:  consts,
		Env:         env,
		ParamsCount: paramsCount,
	}
}

// NewNativeValue creates a native function in potatolang
func NewNativeValue(paramsCount int, vararg bool, f func(env *Env)) Value {
	cls := &Closure{
		ParamsCount: byte(paramsCount),
		native:      f,
	}
	cls.Set(ClsNative)
	if vararg {
		cls.Set(ClsVararg)
	}
	return Fun(cls)
}

func (c *Closure) Set(opt byte) {
	c.options |= opt
}

func (c *Closure) Is(opt byte) bool {
	return (c.options & opt) > 0
}

func (c *Closure) Source() string {
	return string(c.source)
}

// Dup duplicates the closure
func (c *Closure) Dup() *Closure {
	cls := *c
	return &cls
}

func (c *Closure) String() string {
	if c.native != nil {
		return fmt.Sprintf("<native%d>", c.ParamsCount)
	}
	p := "closure"
	if c.Is(ClsNoEnvescape) {
		p = "func"
	}

	hash := uint32(0)
	for _, v := range c.Code {
		hash = hash*31 + v
	}

	x := fmt.Sprintf("<%s%d_%x_%dk", p, c.ParamsCount, hash/65536, len(c.ConstTable))
	if c.Is(ClsYieldable) {
		x += "_y"
	}
	if c.Is(ClsRecoverable) {
		x += "_safe"
	}
	return x + ">"
}

func (c *Closure) PrettyString() string {
	if c.native != nil {
		return "[native Code]"
	}
	return c.crPrettify(0)
}

// Exec executes the closure with the given Env
func (c *Closure) Exec(newEnv *Env) (Value, Value) {
	if c.native == nil {
		if c.lastenv != nil {
			newEnv = c.lastenv
		} else {
			newEnv.SetParent(c.Env)
		}

		v, vb, np, yield := ExecCursor(newEnv, c, c.lastp)
		if yield {
			c.lastp = np
			c.lastenv = newEnv
		} else {
			c.lastp = 0
			c.lastenv = nil
		}
		return v, vb
	}

	// For a native function, it doesn't have its own Env,
	// so newEnv's parent is the Env where this function was called.

	// Check vararg
	if c.Is(ClsVararg) {
		if len(newEnv.stack) > int(c.ParamsCount) {
			v := newEnv.stack[c.ParamsCount]
			if v.Type() == UPK {
				if len(newEnv.stack) != int(c.ParamsCount)+1 {
					panicf("unpacked values should be the last arguments")
				}
				newEnv.Vararg = v.asUnpacked()
			} else {
				newEnv.Vararg = newEnv.stack[c.ParamsCount:]
				newEnv.stack = newEnv.stack[:c.ParamsCount]
			}
		}
	}

	c.native(newEnv)
	return newEnv.A, newEnv.B
}

func (c *Closure) Call(a ...Value) (Value, Value) {
	if len(a) != int(c.ParamsCount) {
		if !(c.Is(ClsVararg) && len(a) > int(c.ParamsCount)) {
			panicf("expect at least %d arguments (got %d)", c.ParamsCount, len(a))
		}
	}
	newEnv := NewEnv(c.Env)
	for i := range a {
		newEnv.LocalPush(a[i])
	}
	return c.Exec(newEnv)
}

func (c *Closure) ImmediateStop() {
	const Stop = uint32(OpEOB) << 26
	for i := range c.Code {
		c.Code[i] = Stop
	}
}
