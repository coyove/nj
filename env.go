package script

import (
	"context"
	"time"
	"unsafe"
)

// Env is the environment for a closure to run within.
// stack contains arguments used by the execution and is a Global shared value, local can only use stack[stackOffset:]
// A and V stores the results of the execution (e.g: return a, b, c => env.A = a, env.V = []Value{b, c})
type Env struct {
	// Global
	Global *Program
	stack  *[]Value

	// Local
	stackOffset uint32
	V           []Value
	A           Value

	nativeSource *Func
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
		s = make([]Value, sz, sz+newSize)
		copy(s, old)
	}
	*env.stack = s[:sz]
}

// Get gets a value from the current stack
func (env *Env) Get(index int) Value {
	return env._get(uint16(index) & 0xfff)
}

// Set sets a value in the current stack
func (env *Env) Set(index int, value Value) {
	env._set(uint16(index)&0xfff, value)
}

// Clear clears the current stack
func (env *Env) Clear() {
	*env.stack = (*env.stack)[:env.stackOffset]
	env.A, env.V = Value{}, nil
}

// Push pushes a value into the current stack
func (env *Env) Push(v Value) {
	*env.stack = append(*env.stack, v)
}

func (env *Env) Size() int {
	return len(*env.stack) - int(env.stackOffset)
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
	index += int(env.stackOffset)
	if index >= len(s) {
		return Value{}
	}
	return s[index]
}

func (env *Env) _set(yx uint16, v Value) {
	if yx == regA {
		env.A = v
	} else {
		index := int(yx & 0xfff)
		s := (*env.stack)
		if yx>>12 == 1 {
			(*env.Global.Stack)[index] = v
		} else {
			s[index+int(env.stackOffset)] = v
		}
	}
}

func (env *Env) Stack() []Value {
	return (*env.stack)[env.stackOffset:]
}

func (env *Env) StackInterface() []interface{} {
	r := make([]interface{}, env.Size())
	for i := range r {
		r[i] = env.Stack()[i].Interface()
	}
	return r
}

// Some useful helper functions

func (env *Env) InInt(i int, defaultValue int64) int64 {
	v := env.Get(i)
	if v.Type() != VNumber {
		return defaultValue
	}
	return v.Int()
}

func (env *Env) InStr(i int, defaultValue string) string {
	v := env.Get(i)
	if v.Type() != VString {
		return defaultValue
	}
	return v._str()
}

func (env *Env) In(i int, expectedType valueType) Value {
	v := env.Get(i)
	if v.Type() != expectedType {
		panicf("%s: bad argument #%d, expect %v, got %v",
			env.nativeSource.name, i+1, expectedType, v.Type())
	}
	return v
}

func (env *Env) Return(a1 Value, an ...Value) {
	env.A, env.V = a1, an
}

func (env *Env) Deadline() (context.Context, func(), time.Time) {
	d := time.Unix(env.Global.Deadline, 0)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	return ctx, cancel, d
}

func (env *Env) GetGlobalCustomValue(key string) interface{} {
	return env.Global.Extras[key]
}

func (e *Env) NewString(s string) Value {
	if e.Global.MaxStackSize > 0 {
		// Loosely control the string size
		remain := (e.Global.MaxStackSize - int64(len(*e.stack))) * 16
		if int64(len(s)) > remain {
			panicf("string overflow, max: %d", remain)
		}
	}
	return _str(s)
}

func (e *Env) NewUnlimitedString(s string) Value {
	return _str(s)
}

func (e *Env) NewStringBytes(s []byte) Value {
	return e.NewString(*(*string)(unsafe.Pointer(&s)))
}

func (e *Env) NewUnlimitedStringBytes(s []byte) Value {
	return e.NewUnlimitedString(*(*string)(unsafe.Pointer(&s)))
}

func (e *Env) checkRemainStackSize(sz int) {
	if e.Global.MaxStackSize > 0 && int64(sz+len(*e.stack)) > e.Global.MaxStackSize {
		panicf("stack overflow, max: %d", e.Global.MaxStackSize)
	}
}
