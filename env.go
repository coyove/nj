package potatolang

// Env is the environment for a closure to run within.
// stack contains arguments used by the execution and is a global shared value, local can only use stack[stackOffset:]
// A and V stores the results of the execution (e.g: return a, b, c => env.A = a, env.V = []Value{b, c})
type Env struct {
	// Global
	global *Global
	stack  *[]Value

	// Local
	stackOffset int
	V           []Value
	A           Value
}

type Global struct {
	Stack *[]Value
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

func (env *Env) _get(yx uint16, cls *Func) (zzz Value) {
	if yx == regA {
		return env.A
	}
	y := yx >> 10
	index := int(yx & 0x3ff)

	if y == 7 {
		return cls.ConstTable[index]
	}

	if y == 1 {
		if env.global == nil {
			panic("nil global")
		}
		return (*env.global.Stack)[index]
	}

	s := *env.stack
	index += env.stackOffset
	if index >= len(s) {
		return Value{}
	}
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
			if env.global == nil {
				panic("nil global")
			}
			(*env.global.Stack)[index] = v
		} else {
			s[index+env.stackOffset] = v
		}
	}
}

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
