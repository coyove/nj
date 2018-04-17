package base

type Env struct {
	parent *Env
	a      Value
	stack  *Stack
	reg    [8]Value
	creg   int
}

func NewEnv(parent *Env) *Env {
	return &Env{
		parent: parent,
		stack:  NewStack(),
		a:      NewValue(),
	}
}

func (e *Env) Reset() {
	e.stack.Clear()
	e.creg = 0
	e.a = NewValue()
}

func (e *Env) Parent() *Env {
	return e.parent
}

func (e *Env) GetA() Value {
	return e.a
}

func (e *Env) SetA(a Value) {
	e.a = a
	e.ClearR()
}

func (e *Env) SetANumber(f float64) {
	e.a = NewNumberValue(f)
	e.ClearR()
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
		return e.GetA()
	}

	env := e.getTop(e, yx)
	return env.stack.Get(int(int16(yx)))
}

func (e *Env) Push(v Value) {
	if v.Type() == TY_phantom {
		return
	}
	e.stack.Add(v)
}

func (e *Env) Size() int {
	return e.stack.Size()
}

func (e *Env) Set(yx int32, v Value) {
	if yx == REG_A {
		e.SetA(v)
	} else {
		env := e.getTop(e, yx)
		env.stack.Set(int(int16(yx)), v)
	}
}

func (e *Env) Stack() *Stack {
	return e.stack
}

func (e *Env) PushR(v Value) { e.reg[e.creg] = v; e.creg++ }
func (e *Env) SizeR() int    { return e.creg }
func (e *Env) ClearR()       { e.creg = 0 }
func (e *Env) R(x int) Value { return e.reg[x] }
func (e *Env) RS() *[8]Value { return &e.reg }
func (e *Env) R0() Value     { return e.reg[0] }
func (e *Env) R1() Value     { return e.reg[1] }
func (e *Env) R2() Value     { return e.reg[2] }
func (e *Env) R3() Value     { return e.reg[3] }
func (e *Env) R4() Value     { return e.reg[4] }
func (e *Env) R5() Value     { return e.reg[5] }
func (e *Env) R6() Value     { return e.reg[6] }
func (e *Env) R7() Value     { return e.reg[7] }
