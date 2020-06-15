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
func NativeFun(f func(env *Env)) Value {
	return Fun(&Closure{native: f, options: ClsNative})
}

func (c *Closure) setOpt(flag bool, opt byte) {
	if flag {
		c.options |= opt
	}
}

func (c *Closure) Is(opt byte) bool { return (c.options & opt) > 0 }

func (c *Closure) Source() string { return string(c.source) }

func (c Closure) Dup() *Closure { return &c }

func (c *Closure) String() string {
	if c.native != nil {
		return "<native>"
	}
	p := "closure"
	if c.Is(ClsNoEnvescape) {
		p = "func"
	}
	if c.Is(ClsVararg) {
		p = "varg-" + p
	}
	if c.Is(ClsYieldable) {
		p = "yield-" + p
	}

	hash := uint32(0)
	for _, v := range c.Code {
		hash = hash*31 + v
	}

	x := fmt.Sprintf("<%s-%d-%04x-%dk", p, c.NumParam, hash/65536, len(c.ConstTable))
	if c.lastenv != nil {
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

// exec executes the closure with the given Env
func (c *Closure) exec(newEnv *Env) (Value, []Value) {
	if c.native == nil {
		if c.lastenv != nil {
			newEnv = c.lastenv
		} else {
			newEnv.SetParent(c.Env)
		}

		v, vb, np, yield := execCursorLoop(newEnv, c, c.lastp)
		if yield {
			c.lastp = np
			c.lastenv = newEnv
		} else {
			c.lastp = 0
			c.lastenv = nil
		}
		return v, vb
	}

	// Native function doesn't have its own Env,
	// so newEnv's parent is the Env where this function was called.
	c.native(newEnv)
	return newEnv.A, newEnv.V
}

func (c *Closure) Call(a ...Value) (Value, []Value) {
	var newEnv *Env
	var varg []Value
	if c.lastenv == nil {
		newEnv = NewEnv(c.Env)
		for i := range a {
			if i >= int(c.NumParam) {
				varg = append(varg, a[i])
			}
			newEnv.Push(a[i])
		}
		if !c.Is(ClsNative) && c.Is(ClsVararg) {
			newEnv._set(uint16(c.NumParam), newUnpackedValue(varg))
		}
	}
	return c.exec(newEnv)
}

// ImmediateStop will try to stop the execution, when called the closure (along with duplicates) become invalid immediately
func (c *Closure) ImmediateStop() {
	const Stop = uint32(OpEOB) << 26
	for i := range c.Code {
		c.Code[i] = Stop
	}
}
