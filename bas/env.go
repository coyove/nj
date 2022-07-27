package bas

import (
	"bytes"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

// Env is the environment for a function to run within.
// stack contains arguments used by the execution and is a global shared value, local can only use stack[stackOffset:]
// A stores the result of the execution
type Env struct {
	Global          *Program
	A               Value
	stack           *[]Value
	stackOffsetFlag uint32
	runtime         Runtime
}

type Runtime struct {
	// Stacktrace layout: N, N-1, ..., 2, 1, 0(current)
	stackN []Stacktrace // [N, N-1, ..., 2]
	stack1 Stacktrace   // 1. if null, then Stack0 is the only one in stacktrace
	stack0 Stacktrace   // 0
}

func (r Runtime) Stacktrace(copy bool) []Stacktrace {
	if r.stack0.Callable == nil {
		return nil
	}
	if r.stack1.Callable == nil {
		return []Stacktrace{r.stack0}
	}
	s := append(r.stackN, r.stack1, r.stack0)
	if copy {
		return append([]Stacktrace{}, s...)
	}
	return s
}

func (r Runtime) push(k Stacktrace) Runtime {
	if r.stack1.Callable != nil {
		r.stackN = append(r.stackN, r.stack1)
	}
	r.stack1 = r.stack0
	r.stack0 = k
	return r
}

func (env *Env) Runtime() Runtime {
	return env.runtime
}

func (env *Env) stackOffset() uint32 {
	return env.stackOffsetFlag & internal.MaxStackSize
}

func (env *Env) growZero(newSize, zeroSize int) {
	old := len(*env.stack)
	env.grow(newSize)
	for i := old; i < zeroSize; i++ {
		(*env.stack)[i] = Value{}
	}
}

func (env *Env) grow(newSize int) {
	s := *env.stack
	sz := int(env.stackOffset()) + newSize
	if sz > cap(s) {
		// if env.Global != nil && env.Global.MaxStackSize > 0 && int64(sz) > env.Global.MaxStackSize {
		// 	panic("stack overflow")
		// }
		old := s
		s = make([]Value, sz+newSize)
		copy(s, old)
	}
	*env.stack = s[:sz]
}

// Get gets a value from the current stack
// Get(-1) means env.A
func (env *Env) Get(index int) Value {
	if index == -1 {
		return env.A
	}
	s := *env.stack
	index += int(env.stackOffset())
	if index < len(s) {
		return s[index]
	}
	return Nil
}

// Set sets a value in the current stack
func (env *Env) Set(index int, value Value) {
	env._set(uint16(index)&typ.RegLocalMask, value)
}

func (env *Env) clear() {
	*env.stack = (*env.stack)[:env.stackOffset()]
	env.A = Value{}
}

func (env *Env) push(v Value) {
	*env.stack = append(*env.stack, v)
}

func (env *Env) pushVararg(v []Value) {
	*env.stack = append(*env.stack, v...)
}

func (env *Env) prepend(v Value) {
	*env.stack = append(*env.stack, Nil)
	copy((*env.stack)[env.stackOffset()+1:], (*env.stack)[env.stackOffset():])
	(*env.stack)[env.stackOffset()] = v
}

func (env *Env) Size() int {
	return len(*env.stack) - int(env.stackOffset())
}

func (env *Env) _get(yx uint16) Value {
	if yx == typ.RegA {
		return env.A
	}
	if yx > typ.RegLocalMask {
		return (*env.Global.stack)[yx&typ.RegLocalMask]
	}
	return (*env.stack)[uint32(yx)+(env.stackOffset())]
}

func (env *Env) _set(yx uint16, v Value) {
	if yx == typ.RegA {
		env.A = v
	} else if yx > typ.RegLocalMask {
		(*env.Global.stack)[yx&typ.RegLocalMask] = v
	} else {
		(*env.stack)[uint32(yx)+env.stackOffset()] = v
	}
}

func (env *Env) Stack() []Value { return (*env.stack)[env.stackOffset():] }

func (env *Env) CopyStack() []Value { return append([]Value{}, env.Stack()...) }

func (env *Env) String() string {
	buf := bytes.NewBufferString("env(")
	buf.WriteString(env.A.String())
	buf.WriteString(")")
	return buf.String()
}

func (env *Env) Bool(idx int) bool { return env.mustBe(typ.Bool, idx).Bool() }

func (env *Env) Str(idx int) string { return env.mustBe(typ.String, idx).String() }

func (env *Env) Num(idx int) Value { return env.mustBe(typ.Number, idx) }

func (env *Env) Int64(idx int) int64 { return env.mustBe(typ.Number, idx).Int64() }

func (env *Env) Int(idx int) int { return env.mustBe(typ.Number, idx).Int() }

func (env *Env) Float64(idx int) float64 { return env.mustBe(typ.Number, idx).Float64() }

func (env *Env) Object(idx int) *Object { return env.mustBe(typ.Object, idx).Object() }

func (env *Env) Native(idx int) *Native { return env.mustBe(typ.Native, idx).Native() }

func (env *Env) Interface(idx int) interface{} {
	if idx == -1 {
		return env.A.Interface()
	}
	return env.Get(idx).Interface()
}

func (env *Env) ThisProp(k string) interface{} {
	return env.Object(-1).Prop(k).Interface()
}

func (env *Env) mustBe(t typ.ValueType, idx int) (v Value) {
	if idx == -1 {
		v = env.A
		if v.Type() != t {
			internal.Panic("argument 'this' should be %v, not %v", t, detail(v))
		}
	} else {
		v = env.Get(idx)
		if v.Type() != t {
			internal.Panic("argument %d expects %v, got %v", idx+1, t, detail(v))
		}
	}
	return v
}

func (env *Env) SetA(a Value) bool {
	env.A = a
	return true
}

func (e *Env) MustGlobal() *Program {
	if e.Global != nil {
		return e.Global
	}
	panic("calling out of program")
}

func (e *Env) Copy() *Env {
	stk := e.CopyStack()
	e2 := &Env{}
	e2.A = e.A
	e2.Global = e.Global
	e2.stack = &stk
	e2.stackOffsetFlag = e.stackOffsetFlag - e.stackOffset()
	e2.runtime = e.runtime
	e2.runtime.stackN = append([]Stacktrace{}, e2.runtime.stackN...)
	return e2
}

func (e *Env) checkStackOverflow() {
	if g := e.Global; g != nil {
		if g.MaxStackSize > 0 && int64(len(*g.stack)) > g.MaxStackSize {
			panic("stack overflow")
		}
		if g.stopped {
			panic("program stopped")
		}
	}
}
