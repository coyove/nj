package potatolang

import (
	"fmt"
)

const (
	// INIT_CAPACITY defines the inital capacity of the stack
	INIT_CAPACITY = 16
)

// Stack is a special structure which will automatically grow when index overflows
type Stack struct {
	data []Value
}

// NewStack creates a new stack
func NewStack() *Stack {
	return &Stack{
		data: make([]Value, 0, INIT_CAPACITY),
	}
}

func (s *Stack) grow(newSize int) {
	if newSize > cap(s.data) {
		old := s.data
		s.data = make([]Value, newSize, newSize*3/2)
		copy(s.data, old)
	}
	s.data = s.data[:newSize]
}

func (s *Stack) Size() int {
	return len(s.data)
}

func (s *Stack) Get(index int) Value {
	if index >= len(s.data) {
		return NewValue()
	}
	return s.data[index]
}

func (s *Stack) Set(index int, value Value) {
	if index >= len(s.data) {
		s.grow(index + 1)
	}
	s.data[index] = value
}

func (s *Stack) Add(value Value) {
	s.Set(len(s.data), value)
}

func (s *Stack) Clear() {
	s.data = s.data[:0]
}

func (s *Stack) InsertStack(index int, s2 *Stack) {
	s.Insert(index, s2.data)
}

func (s *Stack) Insert(index int, data []Value) {
	if index <= len(s.data) {
		ln := len(s.data)
		s.grow(ln + len(data))
		copy(s.data[len(s.data)-(ln-index):], s.data[index:])
	} else {
		s.grow(index + len(data))
	}
	copy(s.data[index:], data)
}

func (s *Stack) Values() []Value {
	return s.data
}

type Env struct {
	parent *Env
	stack  *Stack

	A, C, E, R0, R1, R2, R3 Value
}

func NewTopEnv() *Env {
	e := NewEnv(nil)
	for _, name := range CoreLibNames {
		e.Push(CoreLibs[name])
	}
	return e
}

func NewEnv(parent *Env) *Env {
	return &Env{
		parent: parent,
		stack:  NewStack(),
		A:      NewValue(),
	}
}

func (e *Env) Reset() {
	e.stack.Clear()
	e.A = NewValue()
}

func (e *Env) Parent() *Env {
	return e.parent
}

func (e *Env) SetParent(parent *Env) {
	e.parent = parent
}

func (e *Env) getTop(start *Env, yx int32) *Env {
	env := start
	y := yx >> 16
	for y > 0 && env != nil {
		env = env.parent
		y--
	}

	if env == nil {
		panic("get: null parent")
	}

	return env
}

func (e *Env) Get(yx int32) Value {
	if yx == REG_A {
		return e.A
	}

	env := e.getTop(e, yx)
	return env.stack.Get(int(int16(yx)))
}

func (e *Env) Push(v Value) {
	e.stack.Add(v)
}

func (e *Env) Size() int {
	return e.stack.Size()
}

func (e *Env) Set(yx int32, v Value) {
	if yx == REG_A {
		e.A = v
	} else {
		env := e.getTop(e, yx)
		env.stack.Set(int(int16(yx)), v)
	}
}

func (e *Env) Stack() *Stack {
	return e.stack
}

// Closure is the closure struct used in potatolang
type Closure struct {
	code      []byte
	env       *Env
	caller    Value
	preArgs   []Value
	native    func(env *Env) Value
	argsCount byte
	status    byte
	yieldable bool
	errorable bool
	lastp     uint32
	lastenv   *Env
}

// NewClosure creates a new closure
func NewClosure(code []byte, env *Env, argsCount byte, yieldable, errorable bool) *Closure {
	return &Closure{
		code:      code,
		env:       env,
		argsCount: argsCount,
		yieldable: yieldable,
		errorable: errorable,
	}
}

// NewNativeValue creates a native function in potatolang
func NewNativeValue(argsCount int, f func(env *Env) Value) Value {
	return NewClosureValue(&Closure{
		argsCount: byte(argsCount),
		native:    f,
	})
}

func (c *Closure) AppendPreArgs(preArgs []Value) {
	if c.preArgs == nil {
		c.preArgs = make([]Value, 0, 4)
	}

	c.preArgs = append(c.preArgs, preArgs...)
	c.argsCount -= byte(len(preArgs))
	if c.argsCount < 0 {
		panic("negative args count")
	}
}

func (c *Closure) PreArgs() []Value {
	return c.preArgs
}

func (c *Closure) SetCode(code []byte) {
	c.code = code
}

func (c *Closure) Code() []byte {
	return c.code
}

func (c *Closure) SetCaller(cr Value) {
	c.caller = cr
}

func (c *Closure) Caller() Value {
	return c.caller
}

// ArgsCount returns the minimal number of arguments closure accepts
func (c *Closure) ArgsCount() int {
	return int(c.argsCount)
}

// Env returns the env inside closure
func (c *Closure) Env() *Env {
	return c.env
}

// Dup duplicates the closure
func (c *Closure) Dup() *Closure {
	cls := NewClosure(c.code, c.env, c.argsCount, c.yieldable, c.errorable)
	cls.caller = c.caller
	cls.lastp = c.lastp
	cls.native = c.native
	if c.preArgs != nil {
		cls.preArgs = make([]Value, len(c.preArgs))
		copy(cls.preArgs, c.preArgs)
	}
	return cls
}

func (c *Closure) String() string {
	if c.native == nil {
		return "closure (\n" + crPrettifyLambda(int(c.argsCount), len(c.preArgs), c.yieldable, c.errorable, c.code, 4) + ")"
	}
	return fmt.Sprintf("closure (\n    <args: %d>\n    <curry: %d>\n    [...] native code\n)", c.argsCount, len(c.preArgs))
}

// Exec executes the closure with the given env
func (c *Closure) Exec(newEnv *Env) Value {

	if c.lastenv != nil {
		newEnv = c.lastenv
	} else {
		newEnv.SetParent(c.env)
		newEnv.C = c.caller
	}

	if c.native == nil {
		v, np, yield := ExecCursor(newEnv, c.code, c.lastp)
		if yield {
			c.lastp = np
			c.lastenv = newEnv
		} else {
			c.lastp = 0
			c.lastenv = nil
		}
		return v
	}
	return c.native(newEnv)
}

// symtable is responsible for recording extra states of compilation
type symtable struct {
	// variable name lookup
	Parent *symtable
	M      map[string]int16

	// flat op immediate value
	I  *float64
	Is *string

	// has yield op
	Y bool

	// has error op
	E bool

	// record line info at chain
	LineInfo bool
}

func (c *symtable) GetRelPosition(key string) int32 {
	m := c
	depth := int32(0)

	for m != nil {
		k, e := m.M[key]
		if e {
			return (depth << 16) | int32(k)
		}

		depth++
		m = m.Parent
	}

	return -1
}
