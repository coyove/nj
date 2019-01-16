package potatolang

import (
	"fmt"
	"unsafe"
)

// Env is the environment for a closure in potatolang to run within.
// Env.stack contains arguments used to execute the closure,
// then the local variables will sequentially take the following spaces.
// Env.A stores the result of executing a closure, or a builtin operator;
// Env.R0 ~ Env.R3 store the arguments to call builtin operators (+, -, *, ...).
type Env struct {
	parent *Env
	trace  []stacktrace
	stack  []Value

	A, R0, R1, R2, R3 Value
	Cancel            *uintptr
}

// NewEnv creates the Env for closure to run within
// parent can be nil, which means this is a top Env
func NewEnv(parent *Env, cancel *uintptr) *Env {
	const initCapacity = 16
	return &Env{
		parent: parent,
		stack:  make([]Value, 0, initCapacity),
		A:      Value{},
		Cancel: cancel,
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

func (env *Env) Get(yx uint16) Value {
	if yx == regA {
		return env.A
	}
	y := yx >> 10
REPEAT:
	if y > 0 && env != nil {
		y, env = y-1, env.parent
		goto REPEAT
	}
	index := int(yx & 0x3ff)
	if index >= len(env.stack) {
		return Value{}
	}
	return env.stack[index]
}

func (env *Env) Set(yx uint16, v Value) {
	if yx == regA {
		env.A = v
	} else {
		y := yx >> 10
	REPEAT:
		if y > 0 && env != nil {
			y, env = y-1, env.parent
			goto REPEAT
		}
		index := int(yx & 0x3ff)
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

func (env *Env) reg(i uint16) *Value {
	return (*Value)(unsafe.Pointer(uintptr(unsafe.Pointer(&env.A)) + SizeofValue*uintptr(i)))
}

const (
	CLS_NOENVESCAPE = 1 << iota
	CLS_HASRECEIVER
	CLS_YIELDABLE
	CLS_RECOVERALL
	CLS_PSEUDO_FOREACH
)

// Closure is the closure struct used in potatolang
type Closure struct {
	code      []uint32
	pos       []uint32
	source    string
	consts    []Value
	env       *Env
	caller    Value
	preArgs   []Value
	native    func(env *Env) Value
	argsCount byte
	options   byte
	lastp     uint32
	lastenv   *Env
}

// NewClosure creates a new closure
func NewClosure(code []uint32, consts []Value, env *Env, argsCount byte) *Closure {
	return &Closure{
		code:      code,
		consts:    consts,
		env:       env,
		argsCount: argsCount,
	}
}

// NewNativeValue creates a native function in potatolang
func NewNativeValue(argsCount int, f func(env *Env) Value) Value {
	cls := &Closure{
		argsCount: byte(argsCount),
		native:    f,
	}
	cls.Set(CLS_NOENVESCAPE)
	return NewClosureValue(cls)
}

func (c *Closure) Set(opt byte) { c.options |= opt }

func (c *Closure) Unset(opt byte) { c.options &= ^opt }

func (c *Closure) Isset(opt byte) bool { return (c.options & opt) > 0 }

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

func (c *Closure) PreArgs() []Value { return c.preArgs }

func (c *Closure) SetCode(code []uint32) { c.code = code }

func (c *Closure) Code() []uint32 { return c.code }

func (c *Closure) Pos() []uint32 { return c.pos }

func (c *Closure) Consts() []Value { return c.consts }

func (c *Closure) BytesCode() []byte { return slice64to8(c.code) }

func (c *Closure) SetCaller(cr Value) { c.caller = cr }

func (c *Closure) Caller() Value { return c.caller }

// ArgsCount returns the minimal number of arguments closure accepts
func (c *Closure) ArgsCount() int { return int(c.argsCount) }

// Env returns the env inside closure
func (c *Closure) Env() *Env { return c.env }

// Dup duplicates the closure
func (c *Closure) Dup() *Closure {
	cls := *c
	if c.preArgs != nil {
		cls.preArgs = make([]Value, len(c.preArgs))
		copy(cls.preArgs, c.preArgs)
	}
	return &cls
}

func (c *Closure) String() string {
	if c.native != nil {
		return fmt.Sprintf("<native_%d_%d>", c.argsCount, len(c.preArgs))
	}
	x := fmt.Sprintf("closure_%d_%d", c.argsCount, len(c.preArgs))
	if c.Isset(CLS_YIELDABLE) {
		x += "_yd"
	}
	if c.Isset(CLS_HASRECEIVER) {
		x += "_this"
	}
	if !c.Isset(CLS_NOENVESCAPE) {
		x += "_esc"
	}
	if c.Isset(CLS_RECOVERALL) {
		x += "_safe"
	}
	return "<" + x + ">"
}

func (c *Closure) PrettyString() string {
	if c.native != nil {
		return "[native code]"
	}
	return c.crPrettify(0)
}

// Exec executes the closure with the given env
func (c *Closure) Exec(newEnv *Env) Value {
	if c.native == nil {

		if c.lastenv != nil {
			newEnv = c.lastenv
		} else {
			newEnv.SetParent(c.env)
		}

		v, np, yield := ExecCursor(newEnv, c, c.lastp)
		if yield {
			c.lastp = np
			c.lastenv = newEnv
		} else {
			c.lastp = 0
			c.lastenv = nil
		}
		return v
	}

	// for a native closure, it doesn't have its own env,
	// so newEnv's parent is the env where this native function was called.
	return c.native(newEnv)
}

// MakeCancelable make the closure and all its children cancelable
// store 1 into the returned *uintptr to cancel them
func (c *Closure) MakeCancelable() *uintptr {
	c.lastenv.Cancel = new(uintptr)
	return c.lastenv.Cancel
}
