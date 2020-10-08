package potatolang

import (
	"unsafe"
)

// Env is the environment for a closure to run within.
// stack contains arguments used to execute the closure
// then the local variables will take the following spaces sequentially.
// A and V stores the results of the execution (e.g: return a, b, c => env.A = a, env.V = []Value{b, c})
// For native variadic functions, V also stores the incoming varargs.
type Env struct {
	global *Env
	stack  []Value
	V      []Value
	A      Value
}

// NewEnv creates the Env for closure to run within
// parent can be nil, which means this is a top Env
func NewEnv(global *Env) *Env {
	//b := make([]byte, 4096)
	//n := runtime.Stack(b, false)
	//log.Println(string(b[:n]))
	return &Env{global: global}
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
	env.A, env.V = Value{}, nil
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

	if y == 1 {
		env = env.global
		if env == nil {
			panic("nil global")
		}
	}

	s := env.stack
	return *(*Value)(unsafe.Pointer(uintptr(unsafe.Pointer(&s[0])) + SizeOfValue*uintptr(index)))
}

func (env *Env) _set(yx uint16, v Value) {
	if yx == regA {
		env.A = v
	} else {
		y := yx >> 10

		if y == 1 {
			env = env.global
			if env.global == nil {
				panic("nil global")
			}
		}
		index := int(yx & 0x3ff)
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
		panicf("bad argument #%d: expect %q, got %+v", i, typeMappings[expectedType], v)
	}
	return v
}

func (env *Env) Return(a1 Value, an ...Value) {
	env.A, env.V = a1, an
}
