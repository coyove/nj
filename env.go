package potatolang

import (
	"unsafe"
)

// Env is the environment for a closure to run within.
// stack contains arguments used to execute the closure (variadic arguments are inside Vararg)
// then the local variables will take the following spaces sequentially.
// A and B store the results of the execution
type Env struct {
	parent *Env
	stack  []Value
	Vararg []Value
	A, B   Value
}

// NewEnv creates the Env for closure to run within
// parent can be nil, which means this is a top Env
func NewEnv(parent *Env) *Env {
	//b := make([]byte, 4096)
	//n := runtime.Stack(b, false)
	//log.Println(string(b[:n]))
	return &Env{parent: parent}
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

// Get gets a value from the current stack
func (env *Env) Get(index int) Value {
	if index >= len(env.stack) {
		return Value{}
	}
	return env.stack[index]
}

// Set sets a value in the current stack
func (env *Env) Set(index int, value Value) {
	if index >= len(env.stack) {
		env.grow(index + 1)
	}
	env.stack[index] = value
}

// Clear clears the current stack
func (env *Env) Clear() {
	env.stack = env.stack[:0]
	env.A, env.B = Value{}, Value{}
}

// Push pushes a value into the current stack
func (env *Env) Push(v Value) {
	// e.stack.Add(v)
	ln := len(env.stack)
	env.grow(ln + 1)
	env.stack[ln] = v
}

func (env *Env) Size() int {
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

// go:noescape
// func envGet(env *Env, yx uint16, K *Closure) Value

func (env *Env) _get(yx uint16, cls *Closure) (zzz Value) {
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

func (env *Env) _set(yx uint16, v Value) {
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

func (env *Env) In(i int, expectedType byte) Value {
	v := env.Get(i)
	if v.Type() != expectedType {
		panicf("argument #%d: expect %q, got %+v", i, typeMappings[expectedType], v)
	}
	return v
}
