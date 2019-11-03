package potatolang

import (
	"encoding/base64"
	"fmt"
	"hash/crc32"
	"unsafe"
)

// Env is the environment for a closure in potatolang to run within.
// Env.stack contains arguments used to execute the closure,
// then the local variables will sequentially take the following spaces.
// Env.A stores the result of an operation
type Env struct {
	parent *Env
	stack  []Value
	A      Value
}

// NewEnv creates the Env for closure to run within
// parent can be nil, which means this is a top Env
func NewEnv(parent *Env) *Env {
	//b := make([]byte, 4096)
	//n := runtime.Stack(b, false)
	//log.Println(string(b[:n]))
	return &Env{
		parent: parent,
	}
}

func (env *Env) grow(newSize int) {
	s := env.stack
	if newSize > cap(s) {
		old := s
		s = make([]Value, newSize, newSize*2)
		copy(s, old)
	}
	env.stack = s[:newSize]
}

// LocalGet gets a value from the current stack
func (env *Env) LocalGet(index int) Value {
	if index >= len(env.stack) {
		return Value{}
	}
	return env.stack[index]
}

// LocalSet sets a value in the current stack
func (env *Env) LocalSet(index int, value Value) {
	if index >= len(env.stack) {
		env.grow(index + 1)
	}
	env.stack[index] = value
}

// LocalClear clears the current stack
func (env *Env) LocalClear() {
	env.stack = env.stack[:0]
	env.A = Value{}
}

// LocalPushFront inserts another stack into the current stack at front
func (env *Env) LocalPushFront(data []Value) {
	ln := len(env.stack)
	env.grow(ln + len(data))
	copy(env.stack[len(env.stack)-ln:], env.stack)
	copy(env.stack, data)
}

// LocalPush pushes a value into the current stack
func (env *Env) LocalPush(v Value) {
	// e.stack.Add(v)
	ln := len(env.stack)
	env.grow(ln + 1)
	env.stack[ln] = v
}

func (env *Env) LocalSize() int {
	//if env == nil {
	//	return 0
	//}
	return len(env.stack)
}

func (env *Env) Parent() *Env {
	return env.parent
}

func (env *Env) SetParent(parent *Env) {
	env.parent = parent
}

func (env *Env) Get(yx uint16, cls *Closure) Value {
	if yx == regA {
		return env.A
	}
	y := yx >> 10
	index := int(yx & 0x3ff)

	if y == 7 {
		return cls.ConstTable[index]
	}

REPEAT:
	if y > 0 && env != nil {
		y, env = y-1, env.parent
		goto REPEAT
	}
	if s := env.stack; index < len(s) {
		return *(*Value)(unsafe.Pointer(uintptr(unsafe.Pointer(&s[0])) + SizeOfValue*uintptr(index)))
		// return s[index]
	}
	return Value{}
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

const (
	ClsNoEnvescape = 1 << iota
	ClsYieldable
	ClsRecoverable
	ClsNative
)

// NewClosure creates a new closure
func NewClosure(code []uint32, consts []Value, env *Env, argsCount byte) *Closure {
	return &Closure{
		Code:       code,
		ConstTable: consts,
		Env:        env,
		ArgsCount:  argsCount,
	}
}

// NewNativeValue creates a native function in potatolang
func NewNativeValue(argsCount int, f func(env *Env) Value) Value {
	cls := &Closure{
		ArgsCount: byte(argsCount),
		native:    f,
	}
	cls.Set(ClsNative)
	return NewClosureValue(cls)
}

func (c *Closure) Set(opt byte) { c.options |= opt }

func (c *Closure) Unset(opt byte) { c.options &= ^opt }

func (c *Closure) Isset(opt byte) bool { return (c.options & opt) > 0 }

func (c *Closure) AppendPartialArgs(preArgs []Value) {
	c.PartialArgs = append(c.PartialArgs, preArgs...)
	if c.ArgsCount < byte(len(preArgs)) {
		panic("negative args count")
	}
	c.ArgsCount -= byte(len(preArgs))
}

func (c *Closure) BytesCode() []byte { return u32Bytes(c.Code) }

// Dup duplicates the closure
func (c *Closure) Dup() *Closure {
	cls := *c
	if len(c.PartialArgs) > 0 {
		cls.PartialArgs = append([]Value{}, c.PartialArgs...)
	}
	return &cls
}

func (c *Closure) String() string {
	if c.native != nil {
		return fmt.Sprintf("<native_%da%dc>", c.ArgsCount, len(c.PartialArgs))
	}
	p := "closure"
	if c.Isset(ClsNoEnvescape) {
		p = "pfun"
	}
	h := crc32.New(crc32.IEEETable)
	h.Write(u32Bytes(c.Code))
	hash := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzzz0123456789").EncodeToString(h.Sum(nil)[:3])

	x := fmt.Sprintf("<%s_%s_%da%dc%dk", p, hash, c.ArgsCount, len(c.PartialArgs), len(c.ConstTable))
	if c.Isset(ClsYieldable) {
		x += "_y"
	}
	if c.Isset(ClsRecoverable) {
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
func (c *Closure) Exec(newEnv *Env) Value {
	if c.native == nil {
		if c.lastenv != nil {
			newEnv = c.lastenv
		} else {
			newEnv.SetParent(c.Env)
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

	// for a native closure, it doesn't have its own Env,
	// so newEnv's parent is the Env where this native function was called.
	return c.native(newEnv)
}

func (c *Closure) ImmediateStop() {
	const Stop = uint32(OpEOB) << 26
	for i := range c.Code {
		c.Code[i] = Stop
	}
}
