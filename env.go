package script

import (
	"bytes"
	"context"
	"time"
)

// Env is the environment for a function to run within.
// stack contains arguments used by the execution and is a global shared value, local can only use stack[stackOffset:]
// A stores the result of the execution
type Env struct {
	Global      *Program
	A           Value
	stack       *[]Value
	stackOffset uint32

	// Debug info for native functions to read
	IP         uint32
	CS         *Func
	Stacktrace []stacktrace
}

func (env *Env) growZero(newSize int) {
	old := len(*env.stack)
	env.grow(newSize)
	for i := old; i < len(*env.stack); i++ {
		(*env.stack)[i] = Value{}
	}
}

func (env *Env) grow(newSize int) {
	s := *env.stack
	sz := int(env.stackOffset) + newSize
	if sz > cap(s) {
		old := s
		s = make([]Value, sz+newSize)
		copy(s, old)
	}
	*env.stack = s[:sz]
}

// Get gets a value from the current stack
func (env *Env) Get(index int) Value {
	s := *env.stack
	index += int(env.stackOffset)
	if index < len(*env.stack) {
		return s[index]
	}
	return Value{}
}

// Set sets a value in the current stack
func (env *Env) Set(index int, value Value) {
	env._set(uint16(index)&0xfff, value)
}

// Clear clears the current stack
func (env *Env) Clear() {
	*env.stack = (*env.stack)[:env.stackOffset]
	env.A = Value{}
}

// Push pushes a value into the current stack
func (env *Env) Push(v Value) {
	*env.stack = append(*env.stack, v)
}

func (env *Env) PushVararg(v []Value) {
	*env.stack = append(*env.stack, v...)
}

func (env *Env) Prepend(v Value) {
	*env.stack = append(*env.stack, Nil)
	copy((*env.stack)[env.stackOffset+1:], (*env.stack)[env.stackOffset:])
	(*env.stack)[env.stackOffset] = v
}

func (env *Env) Size() int {
	return len(*env.stack) - int(env.stackOffset)
}

func (env *Env) _get(yx uint16) Value {
	if yx == regA {
		return env.A
	}

	// if yx>>12 == 1 {
	if yx >= 1<<12 {
		return (*env.Global.Stack)[yx&0xfff]
	}
	return (*env.stack)[yx+uint16(env.stackOffset)]
}

func (env *Env) _set(yx uint16, v Value) {
	if yx == regA {
		env.A = v
	} else if yx >= 1<<12 {
		(*env.Global.Stack)[yx&0xfff] = v
	} else {
		(*env.stack)[yx+uint16(env.stackOffset)] = v
	}
}

func (env *Env) Stack() []Value { return (*env.stack)[env.stackOffset:] }

func (env *Env) CopyStack() []Value { return append([]Value{}, env.Stack()...) }

func (env *Env) Deadline() (context.Context, func(), time.Time) {
	if env.Global.Deadline == 0 {
		return context.TODO(), func() {}, time.Time{}
	}
	d := time.Unix(0, env.Global.Deadline)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	return ctx, cancel, d
}

func (env *Env) String() string {
	buf := bytes.NewBufferString("env(")
	buf.WriteString(env.A.String())
	buf.WriteString(")")
	return buf.String()
}
