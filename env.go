package potatolang

// Env is the environment for a closure to run within.
// stack contains arguments used to execute the closure
// then the local variables will take the following spaces sequentially.
// A and V stores the results of the execution (e.g: return a, b, c => env.A = a, env.V = []Value{b, c})
// For native variadic functions, V also stores the incoming varargs.
type Env struct {
	global      []Value
	stackOffset int
	stack       *[]Value
	V           []Value
	A           Value
}

// NewEnv creates the Env for closure to run within
// parent can be nil, which means this is a top Env
func NewEnv() *Env {
	//b := make([]byte, 4096)
	//n := runtime.Stack(b, false)
	//log.Println(string(b[:n]))
	return &Env{stack: new([]Value)}
}

func (env *Env) grow(newSize int) {
	s := *env.stack
	sz := env.stackOffset + newSize
	if sz > cap(s) {
		old := s
		s = make([]Value, sz, sz+newSize)
		copy(s, old)
	}
	*env.stack = s[:sz]
}

// Get gets a value from the current stack
func (env *Env) Get(index int) Value {
	return env._get(uint16(index)&0x3ff, nil)
}

// Set sets a value in the current stack
func (env *Env) Set(index int, value Value) {
	env._set(uint16(index)&0x3ff, value)
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
	return len(*env.stack) - env.stackOffset
}

func (env *Env) _get(yx uint16, cls *Closure) (zzz Value) {
	if yx == regA {
		return env.A
	}
	y := yx >> 10
	index := int(yx & 0x3ff)

	if y == 7 {
		return cls.ConstTable[index]
	}

	s := *env.stack
	if y == 1 {
		s = env.global
		if s == nil {
			panic("nil global")
		}
		return s[index]
	}

	index += env.stackOffset
	if index >= len(s) {
		return Value{}
	}
	// return *(*Value)(unsafe.Pointer(uintptr(unsafe.Pointer(&s[0])) + SizeOfValue*uintptr(index)))
	return s[index]
}

func (env *Env) _set(yx uint16, v Value) {
	if yx == regA {
		env.A = v
	} else {
		index := int(yx & 0x3ff)
		y := yx >> 10
		s := (*env.stack)
		if y == 1 {
			s = env.global
			if s == nil {
				panic("nil global")
			}
			s[index] = v
		} else {
			s[index+env.stackOffset] = v
		}
	}
}

// Stack returns the current stack
func (env *Env) Stack() []Value {
	return (*env.stack)[env.stackOffset:]
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
