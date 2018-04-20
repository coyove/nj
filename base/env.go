package base

type Env struct {
	parent *Env
	stack  *Stack
	creg   int

	A, R0, R1, R2, R3 Value
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
	e.creg = 0
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
