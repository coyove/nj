package potatolang

import (
	"fmt"
)

// Env is the environment for a closure in potatolang to run within.
// The stack contains arguments used to execute the closure,
// all the local variables will sequentially take the following spaces.
// A stores the result of executing a closure, or a builtin operator;
// C stores the caller;
// E stores the error;
// R0 ~ R3 store the arguments to call builtin operators (+, -, *, ...).
// After each function calling in potato, E will be propagated to the upper Env.
// Explicitly calling error() will get E and clear E.
type Env struct {
	parent *Env
	stack  []Value

	A, C, E, R0, R1, R2, R3 Value
}

// NewEnv creates the Env for closure to run within
// parent can be nil, which means this is a top Env
func NewEnv(parent *Env) *Env {
	const initCapacity = 16
	return &Env{
		parent: parent,
		stack:  make([]Value, 0, initCapacity),
		A:      Value{},
	}
}

func (env *Env) grow(newSize int) {
	if newSize > cap(env.stack) {
		old := env.stack
		env.stack = make([]Value, newSize, newSize*3/2)
		copy(env.stack, old)
	}
	env.stack = env.stack[:newSize]
}

// SGet gets a value from the current stack
func (env *Env) SGet(index int) Value {
	if index >= len(env.stack) {
		return Value{}
	}
	return env.stack[index]
}

// SSet sets a value in the current stack
func (env *Env) SSet(index int, value Value) {
	if index >= len(env.stack) {
		env.grow(index + 1)
	}
	env.stack[index] = value
}

// SClear clears the current stack
func (env *Env) SClear() {
	env.stack = env.stack[:0]
	env.A = Value{}
}

// SInsert inserts another stack into the current stack
func (env *Env) SInsert(index int, data []Value) {
	if index <= len(env.stack) {
		ln := len(env.stack)
		env.grow(ln + len(data))
		copy(env.stack[len(env.stack)-(ln-index):], env.stack[index:])
	} else {
		env.grow(index + len(data))
	}
	copy(env.stack[index:], data)
}

// SPush pushes a value into the current stack
func (env *Env) SPush(v Value) {
	// e.stack.Add(v)
	ln := len(env.stack)
	env.grow(ln + 1)
	env.stack[ln] = v
}

func (env *Env) SSize() int {
	return len(env.stack)
}

func (e *Env) Parent() *Env {
	return e.parent
}

func (e *Env) SetParent(parent *Env) {
	e.parent = parent
}

func (env *Env) Get(yx uint32) Value {
	if yx == regA {
		return env.A
	}
	y := yx >> 16
REPEAT:
	if y > 0 && env != nil {
		y, env = y-1, env.parent
		goto REPEAT
	}
	index := int(uint16(yx))
	if index >= len(env.stack) {
		return Value{}
	}
	return env.stack[index]
}

func (env *Env) Set(yx uint32, v Value) {
	if yx == regA {
		env.A = v
	} else {
		y := yx >> 16
	REPEAT:
		if y > 0 && env != nil {
			y, env = y-1, env.parent
			goto REPEAT
		}
		index := int(uint16(yx))
		if index >= len(env.stack) {
			env.grow(index + 1)
		}
		env.stack[index] = v
	}
}

// Stack returns the current stack
func (env *Env) Stack() []Value {
	return env.stack
}

// Closure is the closure struct used in potatolang
type Closure struct {
	code        []uint16
	consts      []Value
	env         *Env
	caller      Value
	preArgs     []Value
	native      func(env *Env) Value
	argsCount   byte
	noenvescape bool
	yieldable   bool
	errorable   bool
	lastp       uint32
	lastenv     *Env
}

// NewClosure creates a new closure
func NewClosure(code []uint16, consts []Value, env *Env, argsCount byte, yieldable, errorable, noenvescape bool) *Closure {
	return &Closure{
		code:        code,
		consts:      consts,
		env:         env,
		argsCount:   argsCount,
		yieldable:   yieldable,
		errorable:   errorable,
		noenvescape: noenvescape,
	}
}

// NewNativeValue creates a native function in potatolang
func NewNativeValue(argsCount int, f func(env *Env) Value) Value {
	return NewClosureValue(&Closure{
		argsCount:   byte(argsCount),
		native:      f,
		noenvescape: true,
	})
}

func (c *Closure) AppendPreArgs(preArgs []Value) {
	if c.preArgs == nil {
		c.preArgs = make([]Value, 0, 4)
	}

	c.preArgs = append(c.preArgs, preArgs...)
	c.argsCount -= byte(len(preArgs))
	if c.argsCount < 0 {
		panic("negative args count")
	}
}

func (c *Closure) PreArgs() []Value {
	return c.preArgs
}

func (c *Closure) SetCode(code []uint16) {
	c.code = code
}

func (c *Closure) Code() []uint16 {
	return c.code
}

func (c *Closure) SetCaller(cr Value) {
	c.caller = cr
}

func (c *Closure) Caller() Value {
	return c.caller
}

// ArgsCount returns the minimal number of arguments closure accepts
func (c *Closure) ArgsCount() int {
	return int(c.argsCount)
}

// Env returns the env inside closure
func (c *Closure) Env() *Env {
	return c.env
}

// Dup duplicates the closure
func (c *Closure) Dup() *Closure {
	cls := NewClosure(c.code, c.consts, c.env, c.argsCount, c.yieldable, c.errorable, c.noenvescape)
	cls.caller = c.caller
	cls.lastp = c.lastp
	cls.native = c.native
	if c.preArgs != nil {
		cls.preArgs = make([]Value, len(c.preArgs))
		copy(cls.preArgs, c.preArgs)
	}
	return cls
}

func (c *Closure) String() string {
	if c.native == nil {
		return "closure (\n" +
			crPrettifyLambda(int(c.argsCount), len(c.preArgs),
				c.yieldable, c.errorable, !c.noenvescape, c.code, c.consts, 4) + ")"
	}
	return fmt.Sprintf("closure (\n    <args: %d>\n    <curry: %d>\n    [...] native code\n)", c.argsCount, len(c.preArgs))
}

// Exec executes the closure with the given env
func (c *Closure) Exec(newEnv *Env) Value {

	if c.lastenv != nil {
		newEnv = c.lastenv
	} else {
		newEnv.SetParent(c.env)
		newEnv.C = c.caller
	}

	if c.native == nil {
		v, np, yield := ExecCursor(newEnv, c.code, c.consts, c.lastp)
		if yield {
			c.lastp = np
			c.lastenv = newEnv
		} else {
			c.lastp = 0
			c.lastenv = nil
		}
		return v
	}
	return c.native(newEnv)
}
