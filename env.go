package script

import (
	"bytes"
	"context"
	"time"
	"unsafe"
)

// Env is the environment for a closure to run within.
// stack contains arguments used by the execution and is a Global shared value, local can only use stack[StackOffset:]
// A and V stores the results of the execution (e.g: return a, b, c => env.A = a, env.V = []Value{b, c})
type Env struct {
	Global      *Program
	stack       *[]Value
	A           Value
	StackOffset uint32

	// Debug info for native functions to read
	DebugCursor     uint32
	DebugCaller     *Func
	DebugStacktrace []stacktrace
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
	sz := int(env.StackOffset) + newSize
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
	index += int(env.StackOffset)
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
	*env.stack = (*env.stack)[:env.StackOffset]
	env.A = Value{}
}

// Push pushes a value into the current stack
func (env *Env) Push(v Value) {
	*env.stack = append(*env.stack, v)
}

func (env *Env) Size() int {
	return len(*env.stack) - int(env.StackOffset)
}

func (env *Env) _get(yx uint16) (zzz Value) {
	if yx == regA {
		return env.A
	}

	index := int(yx & 0xfff)
	if yx>>12 == 1 {
		return (*env.Global.Stack)[index]
	}

	s := *env.stack
	index += int(env.StackOffset)
	return s[index]
}

func (env *Env) _set(yx uint16, v Value) {
	if yx == regA {
		env.A = v
	} else {
		index := int(yx & 0xfff)
		s := *env.stack
		if yx>>12 == 1 {
			(*env.Global.Stack)[index] = v
		} else {
			s[index+int(env.StackOffset)] = v
		}
	}
}

func (env *Env) Stack() []Value { return (*env.stack)[env.StackOffset:] }

func (env *Env) CopyStack() []Value { return append([]Value{}, env.Stack()...) }

func (env *Env) StackInterface() []interface{} {
	r := make([]interface{}, env.Size())
	for i := range r {
		r[i] = env.Stack()[i].Interface()
	}
	return r
}

// Some useful helper functions

func (env *Env) Deadline() (context.Context, func(), time.Time) {
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

func (env *Env) NewString(s string) Value {
	if len(s) > 7 {
		env.Global.DecrDeadsize(int64(len(s)))
	}
	return String(s)
}

func (env *Env) NewString2(a, b string) Value {
	if len(a)+len(b) > 7 {
		env.Global.DecrDeadsize(int64(len(a) + len(b)))
	}
	return String(a + b)
}

func (env *Env) NewArray(e ...Value) Value {
	env.Global.DecrDeadsize(int64(len(e)) * ValueSize)
	return Array(e...)
}

func (env *Env) NewMap(kvs ...Value) Value {
	env.Global.DecrDeadsize(int64(len(kvs)) * ValueSize)
	return Dict(kvs...)
}

func (env *Env) NewStringBytes(s []byte) Value {
	return env.NewString(*(*string)(unsafe.Pointer(&s)))
}
