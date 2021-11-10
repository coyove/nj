package script

import (
	"bytes"
	"fmt"
	"math"
	"unicode/utf8"

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
	native *Func
	stacks []stacktrace
}

func (e *ExecError) GetRootPanic() interface{} { return e.r }

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
		if e.native != nil {
			msg.WriteString(e.native.Name)
			msg.WriteString("() ")
		}
		msg.WriteString(fmt.Sprintf("%v\n", e.r))
	}
	return msg.String()
}

// internalExecCursorLoop executes 'K' under 'env' from the given start 'cursor'
func internalExecCursorLoop(env Env, K *Func, cursor uint32) Value {
	stackEnv := env
	stackEnv.stackOffset = uint32(len(*env.stack))

	var retStack []stacktrace
	var nativeCls *Func

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
				e.native = nativeCls
				e.stacks = make([]stacktrace, len(retStack)+1)
				copy(e.stacks, retStack)
				e.stacks[len(e.stacks)-1] = rr
				panic(e)
			}
		}
	}()

	for {
		v := K.Code.Code[cursor]
		bop, opa, opb := v.op, v.a, uint16(v.b)
		cursor++

		switch bop {
		case typ.OpSet:
			env._set(opa, env._get(opb))
		case typ.OpInc:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				if va.IsInt() && vb.IsInt() {
					env.A = Int(va.unsafeInt() + vb.unsafeInt())
				} else {
					env.A = Float(va.Float() + vb.Float())
				}
				env._set(opa, env.A)
			} else {
				switch va.Type() {
				case typ.Table:
					k, v := va.Table().Next(vb)
					env.A = Array(k, v)
				case typ.String:
					idx := int64(0)
					if vb != Nil {
						idx = vb.MustInt("string iteration")
					}
					r, sz := utf8.DecodeRuneInString(va.Str()[idx:])
					if sz == 0 {
						env.A = Array(Nil, Nil)
					} else {
						env.A = Array(Int(int64(sz)+idx), Rune(r))
					}
				default:
					panicf("inc "+errNeedNumbers, va.Type(), vb.Type())
				}
			}
		case typ.OpAdd:
			va, vb := env._get(opa), env._get(opb)
			switch va.Type() + vb.Type() {
			case typ.Number + typ.Number:
				if sum := va.ptr() + vb.ptr(); sum == int64Marker2 {
					env.A = Int(va.unsafeInt() + vb.unsafeInt())
				} else {
					env.A = Float(va.Float() + vb.Float())
				}
			case typ.String + typ.String:
				env.A = Str(va.Str() + vb.Str())
			case typ.String + typ.Number:
				env.A = Str(va.String() + vb.String())
			default:
				panicf("add "+errNeedNumbersOrStrings, va.Type(), vb.Type())
			}
		case typ.OpSub:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				if sum := va.ptr() + vb.ptr(); sum == int64Marker2 {
					env.A = Int(va.unsafeInt() - vb.unsafeInt())
				} else {
					env.A = Float(va.Float() - vb.Float())
				}
			} else {
				panicf("sub "+errNeedNumbers, va.Type(), vb.Type())
			}
		case typ.OpMul:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				if sum := va.ptr() + vb.ptr(); sum == int64Marker2 {
					env.A = Int(va.unsafeInt() * vb.unsafeInt())
				} else {
					env.A = Float(va.Float() * vb.Float())
				}
			} else {
				panicf("mul "+errNeedNumbers, va.Type(), vb.Type())
			}
		case typ.OpDiv:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Float(va.Float() / vb.Float())
			} else {
				panicf("div "+errNeedNumbers, va.Type(), vb.Type())
			}
		case typ.OpIDiv:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int(va.Int() / vb.Int())
			} else {
				panicf("idiv "+errNeedNumbers, va.Type(), vb.Type())
			}
		case typ.OpMod:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int(va.Int() % vb.Int())
			} else {
				panicf("mod "+errNeedNumbers, va.Type(), vb.Type())
			}
		case typ.OpEq:
			env.A = Bool(env._get(opa).Equal(env._get(opb)))
		case typ.OpNeq:
			env.A = Bool(!env._get(opa).Equal(env._get(opb)))
		case typ.OpLess:
			switch va, vb := env._get(opa), env._get(opb); va.Type() + vb.Type() {
			case typ.Number + typ.Number:
				if sum := va.ptr() + vb.ptr(); sum == int64Marker2 {
					env.A = Bool(va.unsafeInt() < vb.unsafeInt())
				} else {
					env.A = Bool(va.Float() < vb.Float())
				}
			case typ.String + typ.String:
				env.A = Bool(va.Str() < vb.Str())
			default:
				panicf("comparison "+errNeedNumbersOrStrings, va.Type(), vb.Type())
			}
		case typ.OpLessEq:
			switch va, vb := env._get(opa), env._get(opb); va.Type() + vb.Type() {
			case typ.Number + typ.Number:
				if sum := va.ptr() + vb.ptr(); sum == int64Marker2 {
					env.A = Bool(va.unsafeInt() <= vb.unsafeInt())
				} else {
					env.A = Bool(va.Float() <= vb.Float())
				}
			case typ.String + typ.String:
				env.A = Bool(va.Str() <= vb.Str())
			default:
				panicf("comparison "+errNeedNumbersOrStrings, va.Type(), vb.Type())
			}
		case typ.OpNot:
			env.A = Bool(env._get(opa).IsFalse())
		case typ.OpBitAnd:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int(va.Int() & vb.Int())
			} else {
				panicf("bitwise and "+errNeedNumbers, va.Type(), vb.Type())
			}
		case typ.OpBitOr:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int(va.Int() | vb.Int())
			} else {
				panicf("bitwise or "+errNeedNumbers, va.Type(), vb.Type())
			}
		case typ.OpBitXor:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int(va.Int() ^ vb.Int())
			} else {
				panicf("bitwise xor "+errNeedNumbers, va.Type(), vb.Type())
			}
		case typ.OpBitLsh:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int(va.Int() << vb.Int())
			} else {
				panicf("bitwise lsh "+errNeedNumbers, va.Type(), vb.Type())
			}
		case typ.OpBitRsh:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int(va.Int() >> vb.Int())
			} else {
				panicf("bitwise rsh "+errNeedNumbers, va.Type(), vb.Type())
			}
		case typ.OpBitURsh:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int(int64(uint64(va.Int()) >> vb.Int()))
			} else {
				panicf("bitwise ursh "+errNeedNumbers, va.Type(), vb.Type())
			}
		case typ.OpBitNot:
			if a := env._get(opa); a.Type() == typ.Number {
				env.A = Int(^a.Int())
			} else {
				panicf("bitwise not "+errNeedNumber, a.Type())
			}
		case typ.OpArray:
			env.A = Array(append([]Value{}, stackEnv.Stack()...)...)
			stackEnv.Clear()
		case typ.OpMap:
			env.A = Map(stackEnv.Stack()...)
			stackEnv.Clear()
		case typ.OpStore:
			subject, v := env._get(opa), env._get(opb)
			switch subject.Type() {
			case typ.Table:
				m := subject.Table()
				env.A = m.Set(env.A, v)
			case typ.Interface:
				reflectStore(subject.Interface(), env.A, v)
				env.A = v
			default:
				panicf("can't store %v into (%v)[%v]", v.Type(), subject.Type(), env.A.Type())
			}
		case typ.OpLoad:
			switch a, idx := env._get(opa), env._get(opb); a.Type() {
			case typ.Table:
				env.A = a.Table().Get(idx)
			case typ.Interface:
				env.A = reflectLoad(a.Interface(), idx)
			case typ.String:
				if idx.Type() == typ.Number {
					if s := a.Str(); idx.Int() >= 0 && idx.Int() < int64(len(s)) {
						env.A = Int(int64(s[idx.Int()]))
					} else {
						env.A = Nil
					}
					break
				} else if idx.Type() == typ.String {
					if f := StringMethods.Table().GetString(idx.Str()); f != Nil {
						if f.Type() == typ.Func {
							f.Func().Receiver = a
						}
						env.A = f
						break
					}
					panicf("string method %q not found", idx.Str())
				}
				fallthrough
			default:
				panicf("can't load (%v)[%v]", a.Type(), idx.Type())
			}
		case typ.OpPush:
			stackEnv.Push(env._get(opa))
		case typ.OpPushVararg:
			*stackEnv.stack = append(*stackEnv.stack, env._get(opa).mustBe(typ.Table, "unpack arguments", 0).Table().ArrayPart()...)
		case typ.OpRet:
			v := env._get(opa)
			if len(retStack) == 0 {
				return v
			}
			// Return to upper stack
			r := retStack[len(retStack)-1]
			cursor = r.cursor
			K = r.cls
			env.stackOffset = r.stackOffset
			env.A = v
			*env.stack = (*env.stack)[:env.stackOffset+uint32(r.cls.StackSize)]
			stackEnv.stackOffset = uint32(len(*env.stack))
			retStack = retStack[:len(retStack)-1]
		case typ.OpLoadFunc:
			if opb != 0 {
				env.A = env._get(opb)
				if env.A.Type() == typ.Func {
					env.A.Func().Receiver = Nil
				}
			} else {
				env.A = env.Global.Functions[opa].Value()
			}
		case typ.OpCall, typ.OpTailCall:
			a := env._get(opa)
			at := a.Type()
			if at == typ.Table {
				a = a.Table().GetString("__call")
				at = a.Type()
			}
			if at != typ.Func {
				panicf("can't call %v", at)
			}
			cls := a.Func()
			if opb != regPhantom {
				stackEnv.Push(env._get(opb))
			}
			if cls.Receiver != Nil {
				stackEnv.Prepend(cls.Receiver)
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
				stackEnv.CS = K
				stackEnv.IP = cursor
				stackEnv.Stacktrace = retStack
				nativeCls = cls
				cls.Native(&stackEnv)
				nativeCls = nil
				env.A = stackEnv.A
				stackEnv.Clear()
			} else {
				stackEnv.growZero(int(cls.StackSize), int(cls.NumParams))

				last := stacktrace{
					cls:         K,
					cursor:      cursor,
					stackOffset: uint32(env.stackOffset),
				}

				// Switch 'env' to 'stackEnv' and clear 'stackEnv'
				cursor = 0
				K = cls
				env.stackOffset = stackEnv.stackOffset
				env.Global = cls.LoadGlobal
				env.A = stackEnv.A

				if bop == typ.OpCall {
					retStack = append(retStack, last)
				}

				stackEnv.stackOffset = uint32(len(*env.stack))
			}
		case typ.OpJmp:
			cursor = uint32(int32(cursor) + v.b)
		case typ.OpIfNot:
			if env.A.IsFalse() {
				cursor = uint32(int32(cursor) + v.b)
			}
		}
	}
}
