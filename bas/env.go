package bas

import (
	"bytes"
	"unsafe"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

// Env is the environment for a function to run within.
// 'stack' represents the global stack, a running function use 'stack[stackOffset:]' as its local stack.
// 'A' stores the result of the execution. 'global' is the topmost function scope, a.k.a. Program.
type Env struct {
	stack           *[]Value
	top             *Program
	A               Value
	stackOffsetFlag uint32
	runtime         stacktraces
}

type stacktraces struct {
	// Stacktrace layout: N, N-1, ..., 2, 1, 0(current)
	stackN []Stacktrace // [N, N-1, ..., 2]
	stack1 Stacktrace   // 1. If nil, then 'stack0' is the only one in stacktrace
	stack0 Stacktrace   // 0
}

func (r stacktraces) Stacktrace(copy bool) []Stacktrace {
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

func (r stacktraces) push(k Stacktrace) stacktraces {
	if r.stack1.Callable != nil {
		r.stackN = append(r.stackN, r.stack1)
	}
	r.stack1 = r.stack0
	r.stack0 = k
	return r
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
		old := s
		s = make([]Value, sz+newSize)
		copy(s, old)
	}
	*env.stack = s[:sz]
}

// Get gets value at 'index' in current stack, Get(-1) means env.A.
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

// Set sets 'value' at 'index' in current stack.
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

func (env *Env) Size() int {
	return len(*env.stack) - int(env.stackOffset())
}

func (env *Env) _get(yx uint16) Value {
	if yx == typ.RegA {
		return env.A
	}
	if yx > typ.RegLocalMask {
		offset := uintptr(yx&typ.RegLocalMask) * ValueSize
		return *(*Value)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(env.top.stack)) + offset))
		// return (*env.global.stack)[yx&typ.RegLocalMask]
	}
	offset := uintptr(uint32(yx)+env.stackOffset()) * ValueSize
	return *(*Value)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(env.stack)) + offset))
	// return (*env.stack)[uint32(yx)+env.stackOffset()]
}

func (env *Env) _set(yx uint16, v Value) {
	if yx == typ.RegA {
		env.A = v
	} else if yx > typ.RegLocalMask {
		offset := uintptr(yx&typ.RegLocalMask) * ValueSize
		*(*Value)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(env.top.stack)) + offset)) = v
		// (*env.global.stack)[yx&typ.RegLocalMask] = v
	} else {
		offset := uintptr(uint32(yx)+env.stackOffset()) * ValueSize
		*(*Value)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(env.stack)) + offset)) = v
		//(*env.stack)[uint32(yx)+env.stackOffset()] = v
	}
}

// Stack returns current stack as a reference.
func (env *Env) Stack() []Value { return (*env.stack)[env.stackOffset():] }

// CopyStack returns a copy of current stack.
func (env *Env) CopyStack() []Value { return append([]Value{}, env.Stack()...) }

func (env *Env) String() string {
	buf := bytes.NewBufferString("env(")
	buf.WriteString(env.A.String())
	buf.WriteString(")")
	return buf.String()
}

// Bool returns value at 'idx' in current stack and asserts its Type() to be a boolean.
func (env *Env) Bool(idx int) bool { return env.mustBe(typ.Bool, idx).Bool() }

// Str returns value at 'idx' in current stack and asserts its Type() to be a string.
func (env *Env) Str(idx int) string { return env.mustBe(typ.String, idx).String() }

// Num returns value at 'idx' in current stack and asserts its Type() to be a number.
func (env *Env) Num(idx int) Value { return env.mustBe(typ.Number, idx) }

// Int64 returns value at 'idx' in current stack and asserts its Type() to be a number.
func (env *Env) Int64(idx int) int64 { return env.mustBe(typ.Number, idx).Int64() }

// Int returns value at 'idx' in current stack and asserts its Type() to be a number.
func (env *Env) Int(idx int) int { return env.mustBe(typ.Number, idx).Int() }

// Float64 returns value at 'idx' in current stack and asserts its Type() to be a number.
func (env *Env) Float64(idx int) float64 { return env.mustBe(typ.Number, idx).Float64() }

// Object returns value at 'idx' in current stack and asserts its Type() to be an Object.
func (env *Env) Object(idx int) *Object { return env.mustBe(typ.Object, idx).Object() }

// Native returns value at 'idx' in current stack and asserts its Type() to be a Native.
func (env *Env) Native(idx int) *Native { return env.mustBe(typ.Native, idx).Native() }

// Interface returns value at 'idx' in current stack as interface{}
func (env *Env) Interface(idx int) interface{} {
	if idx == -1 {
		return env.A.Interface()
	}
	return env.Get(idx).Interface()
}

// ThisProp returns value, by property 'k' of 'this', as interface{}.
func (env *Env) ThisProp(k string) interface{} {
	return env.Object(-1).Prop(k).Interface()
}

func (env *Env) mustBe(t typ.ValueType, idx int) (v Value) {
	if idx == -1 {
		v = env.A
		if v.Type() != t {
			internal.Panic("expects 'this' to be %v, got %v", t, detail(v))
		}
	} else {
		v = env.Get(idx)
		if v.Type() != t {
			internal.Panic("expects argument #%d to be %v, got %v", idx+1, t, detail(v))
		}
	}
	return v
}

func (env *Env) SetA(a Value) bool {
	env.A = a
	return true
}

func (e *Env) MustProgram() *Program {
	if e.top != nil {
		return e.top
	}
	panic("out of program")
}

func (e *Env) Copy() *Env {
	stk := e.CopyStack()
	e2 := &Env{}
	e2.A = e.A
	e2.top = e.top
	e2.stack = &stk
	e2.stackOffsetFlag = e.stackOffsetFlag - e.stackOffset()
	e2.runtime = e.runtime
	e2.runtime.stackN = append([]Stacktrace{}, e2.runtime.stackN...)
	return e2
}

func (e *Env) checkStackOverflow() {
	if g := e.top; g != nil {
		if int64(len(*g.stack)) > g.MaxStackSize {
			panic("stack overflow")
		}
		if g.stopped {
			panic("program stopped")
		}
	}
}
