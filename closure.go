package potatolang

import (
	"fmt"
)

const (
	ClsNoEnvescape = 1 << iota
	ClsYieldable
	ClsNative
	ClsVararg
)

type Closure struct {
	Code       []uint32
	Pos        posVByte
	source     []byte
	ConstTable []Value
	Env        *Env
	NumParam   byte
	options    byte
	lastp      uint32
	lastenv    *Env
	native     func(env *Env)
}

// NativeFun creates a native function in potatolang
func NativeFun(numParam byte, f func(env *Env)) Value {
	cls := &Closure{
		NumParam: numParam,
		native:   f,
		options:  ClsNative | ClsVararg,
	}
	return Fun(cls)
}

func (c *Closure) setOpt(flag bool, opt byte) {
	if flag {
		c.options |= opt
	}
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
		return fmt.Sprintf("<native-%d>", c.NumParam)
	}
	p := "closure"
	if c.Is(ClsNoEnvescape) {
		p = "function"
	}

	hash := uint32(0)
	for _, v := range c.Code {
		hash = hash*31 + v
	}

	x := fmt.Sprintf("<%s-%d-%04x-%dk", p, c.NumParam, hash/65536, len(c.ConstTable))
	if c.Is(ClsYieldable) {
		x += "-y"
	}
	if c.lastp != 0 {
		x += fmt.Sprintf("-%xy", c.lastp)
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
func (c *Closure) Exec(newEnv *Env) (Value, []Value) {
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
	c.native(newEnv)
	return newEnv.A, newEnv.V
}

func (c *Closure) Call(a ...Value) (Value, []Value) {
	newEnv := NewEnv(c.Env)
	for i := range a {
		if i >= int(c.NumParam) {
			newEnv.V = append(newEnv.V, a[i])
		} else {
			newEnv.Push(a[i])
		}
	}
	if !c.Is(ClsNative) && c.Is(ClsVararg) {
		newEnv._set(uint16(c.NumParam), newUnpackedValue(newEnv.V))
	}
	return c.Exec(newEnv)
}

func (c *Closure) ImmediateStop() {
	const Stop = uint32(OpEOB) << 26
	for i := range c.Code {
		c.Code[i] = Stop
	}
}
