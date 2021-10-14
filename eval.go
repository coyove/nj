package script

import (
	"bytes"
	"fmt"
	"math"
	"sync/atomic"

	"github.com/coyove/script/typ"
)

type stacktrace struct {
	cursor      uint32
	stackOffset uint32
	cls         *Func
}

// ExecError represents the runtime error
type ExecError struct {
	r      interface{}
	stacks []stacktrace
}

func (e *ExecError) Error() string {
	msg := bytes.Buffer{}
	msg.WriteString("stacktrace:\n")
	for i := len(e.stacks) - 1; i >= 0; i-- {
		r := e.stacks[i]
		src := uint32(0)
		for i := 0; i < len(r.cls.Code.Pos); {
			var opx uint32 = math.MaxUint32
			ii, op, line := r.cls.Code.Pos.read(i)
			if ii < len(r.cls.Code.Pos)-1 {
				_, opx, _ = r.cls.Code.Pos.read(ii)
			}
			if r.cursor >= op && r.cursor < opx {
				src = line
				break
			}
			if r.cursor < op && i == 0 {
				src = line
				break
			}
			i = ii
		}
		// the recorded cursor was advanced by 1 already
		msg.WriteString(fmt.Sprintf("%s at line %d (cursor: %d)\n", r.cls.Name, src, r.cursor-1))
	}
	if e.r != nil {
		msg.WriteString("root panic:\n")
		msg.WriteString(fmt.Sprintf("%v\n", e.r))
	}
	return msg.String()
}

// internalExecCursorLoop executes 'K' under 'env' from the given start 'cursor'
func internalExecCursorLoop(env Env, K *Func, cursor uint32) Value {
	stackEnv := env
	stackEnv.StackOffset = uint32(len(*env.stack))

	var retStack []stacktrace

	defer func() {
		if r := recover(); r != nil {
			rr := stacktrace{
				cursor: cursor,
				cls:    K,
			}

			if re, ok := r.(*ExecError); ok {
				retStack = append(retStack, rr)
				re.stacks = append(retStack, re.stacks...)
				panic(re)
			} else {
				e := &ExecError{}
				e.r = r // root panic
				e.stacks = make([]stacktrace, len(retStack)+1)
				copy(e.stacks, retStack)
				e.stacks[len(e.stacks)-1] = rr
				panic(e)
			}
		}
	}()

	for {
		if env.Global.Deadline != 0 {
			if atomic.LoadInt64(&now) > env.Global.Deadline {
				panicf("timeout")
			}
		}

		v := K.Code.Code[cursor]
		cursor++
		bop, opa, opb := splitInst(v)

		switch bop {
		case typ.OpSet:
			env._set(opa, env._get(opb))
		case typ.OpInc:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				vaf, vai, vaIsInt := va.Num()
				vbf, vbi, vbIsInt := vb.Num()
				if vaIsInt && vbIsInt {
					env.A = Int(vai + vbi)
				} else {
					env.A = Float(vaf + vbf)
				}
				env._set(opa, env.A)
			} else if va.Type() == typ.Map {
				k, v := va.Map().Next(vb)
				env.A = Array(k, v)
			} else {
				panicf("inc " + errNeedNumbers)
			}
		case typ.OpAdd:
			va, vb := env._get(opa), env._get(opb)
			switch va.Type() + vb.Type() {
			case typ.Number + typ.Number:
				if sum := va.puintptr() + vb.puintptr(); sum == int64Marker2 {
					env.A = Int(va.unsafeint() + vb.unsafeint())
				} else if sum == float64Marker2 {
					env.A = Float(va.unsafefloat() + vb.unsafefloat())
				} else {
					env.A = Float(va.Float() + vb.Float())
				}
			case typ.String + typ.String:
				env.A = Str(va.Str() + vb.Str())
			case typ.String + typ.Number:
				env.A = Str(va.String() + vb.String())
			default:
				panicf("add " + errNeedNumbersOrStrings)
			}
		case typ.OpSub:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				if sum := va.puintptr() + vb.puintptr(); sum == int64Marker2 {
					env.A = Int(va.unsafeint() - vb.unsafeint())
				} else if sum == float64Marker2 {
					env.A = Float(va.unsafefloat() - vb.unsafefloat())
				} else {
					env.A = Float(va.Float() - vb.Float())
				}
			} else {
				panicf("sub " + errNeedNumbers)
			}
		case typ.OpMul:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				if sum := va.puintptr() + vb.puintptr(); sum == int64Marker2 {
					env.A = Int(va.unsafeint() * vb.unsafeint())
				} else if sum == float64Marker2 {
					env.A = Float(va.unsafefloat() * vb.unsafefloat())
				} else {
					env.A = Float(va.Float() * vb.Float())
				}
			} else {
				panicf("mul " + errNeedNumbers)
			}
		case typ.OpDiv:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Float(va.Float() / vb.Float())
			} else {
				panicf("div " + errNeedNumbers)
			}
		case typ.OpIDiv:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int(va.Int() / vb.Int())
			} else {
				panicf("idiv " + errNeedNumbers)
			}
		case typ.OpMod:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				if sum := va.puintptr() + vb.puintptr(); sum == int64Marker2 {
					env.A = Int(va.unsafeint() % vb.unsafeint())
				} else {
					panicf("modulating floating numbers is ambiguous, use math.mod/math.remainder instead")
				}
			} else {
				panicf("mod " + errNeedNumbers)
			}
		case typ.OpEq:
			env.A = Bool(env._get(opa).Equal(env._get(opb)))
		case typ.OpNeq:
			env.A = Bool(!env._get(opa).Equal(env._get(opb)))
		case typ.OpLess:
			switch va, vb := env._get(opa), env._get(opb); va.Type() + vb.Type() {
			case typ.Number + typ.Number:
				if sum := va.puintptr() + vb.puintptr(); sum == int64Marker2 {
					env.A = Bool(va.unsafeint() < vb.unsafeint())
				} else if sum == float64Marker2 {
					env.A = Bool(va.unsafefloat() < vb.unsafefloat())
				} else {
					env.A = Bool(va.Float() < vb.Float())
				}
			case typ.String + typ.String:
				env.A = Bool(va.Str() < vb.Str())
			default:
				panicf("comparison " + errNeedNumbersOrStrings)
			}
		case typ.OpLessEq:
			switch va, vb := env._get(opa), env._get(opb); va.Type() + vb.Type() {
			case typ.Number + typ.Number:
				if sum := va.puintptr() + vb.puintptr(); sum == int64Marker2 {
					env.A = Bool(va.unsafeint() <= vb.unsafeint())
				} else if sum == float64Marker2 {
					env.A = Bool(va.unsafefloat() <= vb.unsafefloat())
				} else {
					env.A = Bool(va.Float() <= vb.Float())
				}
			case typ.String + typ.String:
				env.A = Bool(va.Str() <= vb.Str())
			default:
				panicf("comparison " + errNeedNumbersOrStrings)
			}
		case typ.OpNot:
			env.A = Bool(env._get(opa).IsFalse())
		case typ.OpBitAnd:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int(va.Int() & vb.Int())
			} else {
				panicf("bitwise and " + errNeedNumbers)
			}
		case typ.OpBitOr:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int(va.Int() | vb.Int())
			} else {
				panicf("bitwise or " + errNeedNumbers)
			}
		case typ.OpBitXor:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int(va.Int() ^ vb.Int())
			} else {
				panicf("bitwise xor " + errNeedNumbers)
			}
		case typ.OpBitLsh:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int(va.Int() << vb.Int())
			} else {
				panicf("bitwise lsh " + errNeedNumbers)
			}
		case typ.OpBitRsh:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int(va.Int() >> vb.Int())
			} else {
				panicf("bitwise rsh " + errNeedNumbers)
			}
		case typ.OpBitURsh:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int(int64(uint64(va.Int()) >> vb.Int()))
			} else {
				panicf("bitwise ursh " + errNeedNumbers)
			}
		case typ.OpBitNot:
			env.A = Int(^env._get(opa).MustNum("bitwise not "+errNeedNumbers, 0).Int())
		case typ.OpArray:
			env.A = Array(append([]Value{}, stackEnv.Stack()...)...)
			stackEnv.Clear()
		case typ.OpMap:
			env.A = Map(append([]Value{}, stackEnv.Stack()...)...)
			stackEnv.Clear()
		case typ.OpStore:
			subject, v := env._get(opa), env._get(opb)
			switch subject.Type() {
			case typ.Map:
				m := subject.Map()
				env.A = m.Set(env.A, v)
			case typ.Interface:
				reflectStore(subject.Interface(), env.A, v)
				env.A = v
			default:
				panicf("cannot store %v into (%v)[%v]", v.Type(), subject.Type(), env.A.Type())
			}
		case typ.OpLoad:
			switch a := env._get(opa); a.Type() {
			case typ.Map:
				env.A = a.Map().Get(env._get(opb))
			case typ.Interface:
				env.A = reflectLoad(a.Interface(), env._get(opb))
			case typ.String:
				idx := env._get(opb)
				if idx.Type() == typ.Number {
					if s := a.Str(); idx.Int() >= 0 && idx.Int() < int64(len(s)) {
						env.A = Int(int64(s[idx.Int()]))
					} else {
						env.A = Nil
					}
					break
				} else if idx.Type() == typ.String {
					if f := StringMethods.Map().GetString(idx.Str()); f != Nil {
						if f.Type() == typ.Func {
							f2 := *f.Func()
							f2.MethodSrc = a
							env.A = f2.Value()
						} else {
							env.A = f
						}
						break
					}
					panicf("string method %q not found", idx.Str())
				}
				fallthrough
			default:
				panicf("cannot load (%v)[%v]", a.Type(), env._get(opb).Type())
			}
		case typ.OpPush:
			stackEnv.Push(env._get(opa))
		case typ.OpPushVararg:
			*stackEnv.stack = append(*stackEnv.stack, env._get(opa).MustMap("unpack arguments", 0).Array()...)
		case typ.OpRet:
			v := env._get(opa)
			if len(retStack) == 0 {
				return v
			}
			// Return upper stack
			r := retStack[len(retStack)-1]
			cursor = r.cursor
			K = r.cls
			env.StackOffset = r.stackOffset
			env.A = v
			*env.stack = (*env.stack)[:env.StackOffset+uint32(r.cls.StackSize)]
			stackEnv.StackOffset = uint32(len(*env.stack))
			retStack = retStack[:len(retStack)-1]
		case typ.OpLoadFunc:
			env.A = env.Global.Functions[opa].Value()
		case typ.OpCall, typ.OpTailCall:
			cls := env._get(opa).MustFunc("invoke function", 0)
			if opb != regPhantom {
				stackEnv.Push(env._get(opb))
			}
			if cls.MethodSrc != Nil {
				stackEnv.Prepend(cls.MethodSrc)
			}
			if cls.Variadic {
				s, w := stackEnv.Stack(), int(cls.NumParams)-1
				if len(s) > w {
					s[w] = Array(append([]Value{}, s[w:]...)...)
				} else {
					stackEnv.grow(w + 1)
					stackEnv._set(uint16(w), Array())
				}
			}
			if cls.Native != nil {
				stackEnv.Global = env.Global
				stackEnv.DebugCaller = K
				stackEnv.DebugCursor = cursor
				stackEnv.DebugStacktrace = retStack
				cls.Native(&stackEnv)
				env.A = stackEnv.A
				stackEnv.Clear()
			} else {
				stackEnv.growZero(int(cls.StackSize))

				last := stacktrace{
					cls:         K,
					cursor:      cursor,
					stackOffset: uint32(env.StackOffset),
				}

				// Switch 'env' to 'stackEnv' and clear 'stackEnv'
				cursor = 0
				K = cls
				env.StackOffset = stackEnv.StackOffset
				env.Global = cls.LoadGlobal
				env.A = stackEnv.A

				if bop == typ.OpCall {
					retStack = append(retStack, last)
				}

				if env.Global.MaxCallStackSize > 0 && int64(len(retStack)) > env.Global.MaxCallStackSize {
					panicf("call stack overflow, max: %d", env.Global.MaxCallStackSize)
				}

				stackEnv.StackOffset = uint32(len(*env.stack))
			}
		case typ.OpJmp:
			cursor = uint32(int32(cursor) + int32(v&0xffffff) - 1<<23)
		case typ.OpIfNot:
			if env.A.IsFalse() {
				cursor = uint32(int32(cursor) + int32(v&0xffffff) - 1<<23)
			}
		}
	}
}
