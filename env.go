package script

import (
	"bytes"
	"context"
	"reflect"
	"time"
	"unsafe"
)

// Env is the environment for a closure to run within.
// stack contains arguments used by the execution and is a Global shared value, local can only use stack[StackOffset:]
// A and V stores the results of the execution (e.g: return a, b, c => env.A = a, env.V = []Value{b, c})
type Env struct {
	// Global
	Global *Program
	stack  *[]Value

	// Local
	StackOffset uint32
	V           []Value
	A           Value

	// Used by Native functions
	NativeSource *Func // points to itself

	// Used by Native debug function
	Debug *debugInfo
}

type debugInfo struct {
	Caller     *Func
	Stacktrace []stacktrace
	Cursor     uint32
}

type DebugState struct {
	Cursor uint32
	Stack  []Value
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
		s = make([]Value, sz, sz+newSize)
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
	env.A, env.V = Value{}, nil
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

func (env *Env) InInterface(i int, allowNil bool, expectedType reflect.Type) interface{} {
	v := env.Get(i)
	if v.Type() == VNil && allowNil {
		return nil
	}
	itf := v.Interface()
	if rt := reflect.TypeOf(itf); rt != expectedType {
		panicf("%s: bad argument #%d, expect %v, got %v",
			env.NativeSource.Name, i+1, expectedType, rt)
	}
	return itf
}

func (env *Env) In(i int, expectedType valueType) Value {
	v := env.Get(i)
	if v.Type() != expectedType {
		panicf("%s: bad argument #%d, expect %v, got %v",
			env.NativeSource.Name, i+1, expectedType, v.Type())
	}
	return v
}

func (env *Env) Return(a ...Value) {
	if len(a) > 0 {
		env.A, env.V = a[0], a[1:]
	}
}

func (env *Env) Return2(a1 Value, an ...Value) {
	env.A, env.V = a1, an
}

func (env *Env) Deadline() (context.Context, func(), time.Time) {
	d := time.Unix(env.Global.Deadline, 0)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	return ctx, cancel, d
}

func (env *Env) String() string {
	buf := bytes.NewBufferString("env(")
	buf.WriteString(env.A.String())
	for _, v := range env.V {
		buf.WriteString(",")
		buf.WriteString(v.String())
	}
	buf.WriteString(")")
	return buf.String()
}

func (env *Env) NewString(s string) Value {
	if env.Global.MaxStringSize > 0 {
		max := env.Global.looseStringSizeLimit()
		if int64(len(s)) > max {
			panicf("string overflow, require %d out of %d", len(s), max)
		}
	}
	env.Global.Survey.StringAlloc += int64(len(s))
	return String(s)
}

func (env *Env) NewStringBytes(s []byte) Value {
	return env.NewString(*(*string)(unsafe.Pointer(&s)))
}

func (env *Env) checkRemainStackSize(sz int) {
	if env.Global.MaxStackSize > 0 && int64(sz+len(*env.stack)) > env.Global.MaxStackSize {
		panicf("stack overflow, require %d out of %d", sz, env.Global.MaxStackSize-int64(len(*env.stack)))
	}
}
