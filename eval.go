package script

import (
	"bytes"
	"fmt"
	"math"
	"sync/atomic"
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
	msg.WriteString("root panic:\n")
	msg.WriteString(fmt.Sprintf("%v\n", e.r))
	return msg.String()
}

// InternalExecCursorLoop executes 'K' under 'env' from the given start 'cursor'
func InternalExecCursorLoop(env Env, K *Func, cursor uint32) Value {
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
		case OpSet:
			env._set(opa, env._get(opb))
		case OpInc:
			vaf, vai, vaIsInt := env._get(opa).MustNumber("inc sym", 0).Num()
			vbf, vbi, vbIsInt := env._get(opb).MustNumber("inc step", 0).Num()
			if vaIsInt && vbIsInt {
				env.A = Int(vai + vbi)
			} else {
				env.A = Float(vaf + vbf)
			}
			env._set(opa, env.A)
		case OpAdd:
			va, vb := env._get(opa), env._get(opb)
			switch va.Type() + vb.Type() {
			case VNumber + VNumber:
				vaf, vai, vaIsInt := va.Num()
				vbf, vbi, vbIsInt := vb.Num()
				if vaIsInt && vbIsInt {
					env.A = Int(vai + vbi)
				} else {
					env.A = Float(vaf + vbf)
				}
			case VString + VString:
				env.A = String(va.rawStr() + vb.rawStr())
			default:
				env.A = String(va.String() + vb.String())
			}
		case OpSub:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == VNumber+VNumber {
				vaf, vai, vaIsInt := va.Num()
				vbf, vbi, vbIsInt := vb.Num()
				if vaIsInt && vbIsInt {
					env.A = Int(vai - vbi)
				} else {
					env.A = Float(vaf - vbf)
				}
			} else {
				va.MustNumber("sub", 0)
				vb.MustNumber("sub", 0)
			}
		case OpMul:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == VNumber+VNumber {
				vaf, vai, vaIsInt := va.Num()
				vbf, vbi, vbIsInt := vb.Num()
				if vaIsInt && vbIsInt {
					env.A = Int(vai * vbi)
				} else {
					env.A = Float(vaf * vbf)
				}
			} else {
				va.MustNumber("mul", 0)
				vb.MustNumber("mul", 0)
			}
		case OpDiv:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == VNumber+VNumber {
				vaf, vai, vaIsInt := va.Num()
				vbf, vbi, vbIsInt := vb.Num()
				if vaIsInt && vbIsInt && vai%vbi == 0 {
					env.A = Int(vai / vbi)
				} else {
					env.A = Float(vaf / vbf)
				}
			} else {
				va.MustNumber("div", 0)
				vb.MustNumber("div", 0)
			}
		case OpIDiv:
			env.A = Int(env._get(opa).MustNumber("idiv", 0).Int() / env._get(opb).MustNumber("idiv", 0).Int())
		case OpMod:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == VNumber+VNumber {
				vaf, vai, vaIsInt := va.Num()
				vbf, vbi, vbIsInt := vb.Num()
				if vaIsInt && vbIsInt {
					env.A = Int(vai % vbi)
				} else {
					env.A = Float(math.Remainder(vaf, vbf))
				}
			} else {
				va.MustNumber("mod", 0)
				vb.MustNumber("mod", 0)
			}
		case OpEq:
			env.A = Bool(env._get(opa).Equal(env._get(opb)))
		case OpNeq:
			env.A = Bool(!env._get(opa).Equal(env._get(opb)))
		case OpLess:
			switch va, vb := env._get(opa), env._get(opb); va.Type() + vb.Type() {
			case VNumber + VNumber:
				vaf, vai, vaIsInt := va.Num()
				vbf, vbi, vbIsInt := vb.Num()
				if vaIsInt && vbIsInt {
					env.A = Bool(vai < vbi)
				} else {
					env.A = Bool(vaf < vbf)
				}
			case VString + VString:
				env.A = Bool(va.rawStr() < vb.rawStr())
			default:
				va.MustNumber("less", 0)
				vb.MustNumber("less", 0)
			}
		case OpLessEq:
			switch va, vb := env._get(opa), env._get(opb); va.Type() + vb.Type() {
			case VNumber + VNumber:
				vaf, vai, vaIsInt := va.Num()
				vbf, vbi, vbIsInt := vb.Num()
				if vaIsInt && vbIsInt {
					env.A = Bool(vai <= vbi)
				} else {
					env.A = Bool(vaf <= vbf)
				}
			case VString + VString:
				env.A = Bool(va.rawStr() <= vb.rawStr())
			default:
				va.MustNumber("less", 0)
				vb.MustNumber("less", 0)
			}
		case OpNot:
			env.A = Bool(env._get(opa).IsFalse())
		case OpBitAnd:
			env.A = Int(env._get(opa).MustNumber("bitwise and", 0).Int() & env._get(opb).MustNumber("bitwise and", 0).Int())
		case OpBitOr:
			env.A = Int(env._get(opa).MustNumber("bitwise or", 0).Int() | env._get(opb).MustNumber("bitwise or", 0).Int())
		case OpBitXor:
			env.A = Int(env._get(opa).MustNumber("bitwise xor", 0).Int() ^ env._get(opb).MustNumber("bitwise xor", 0).Int())
		case OpBitLsh:
			env.A = Int(env._get(opa).MustNumber("bitwise lsh", 0).Int() << env._get(opb).MustNumber("bitwise lsh", 0).Int())
		case OpBitRsh:
			env.A = Int(env._get(opa).MustNumber("bitwise rsh", 0).Int() >> env._get(opb).MustNumber("bitwise rsh", 0).Int())
		case OpBitURsh:
			a := env._get(opa).MustNumber("bitwise unsigned rsh", 0).Int()
			b := env._get(opb).MustNumber("bitwise unsigned rsh", 0).Int()
			env.A = Int(int64(uint64(a) >> b))
		case OpMapArray:
			env.A = Array(append([]Value{}, stackEnv.Stack()...)...)
			stackEnv.Clear()
		case OpMap:
			env.A = ArrayMap(append([]Value{}, stackEnv.Stack()...)...)
			stackEnv.Clear()
		case OpStore:
			subject, v := env._get(opa), env._get(opb)
			switch subject.Type() {
			case VArray:
				m := subject.Array()
				env.A = m.Set(env.A, v)
			case VInterface:
				reflectStore(subject.Interface(), env.A, v)
				env.A = v
			default:
				panicf("require array or interface to store into")
			}
		case OpLoad:
			switch a := env._get(opa); a.Type() {
			case VArray:
				env.A = a.Array().Get(env._get(opb))
				if env.A.Type() == VFunction && a.Array().Parent != nil {
					f := *env.A.Function()
					f.MethodSrc = a
					env.A = Function(&f)
				}
			case VInterface:
				env.A = reflectLoad(a.Interface(), env._get(opb))
			case VString:
				idx := env._get(opb)
				if idx.Type() == VNumber {
					if s := a.rawStr(); idx.Int() >= 0 && idx.Int() < int64(len(s)) {
						env.A = Int(int64(s[idx.Int()]))
					} else {
						env.A = Nil
					}
					break
				} else if idx.Type() == VString {
					if f := StringMethods.Array().GetString(idx.rawStr()); f != Nil {
						if f.Type() == VFunction {
							f2 := *f.Function()
							f2.MethodSrc = a
							env.A = Function(&f2)
						} else {
							env.A = f
						}
						break
					}
					panicf("string method %q not found", idx.rawStr())
				}
				fallthrough
			default:
				panicf("require array, string or interface to load from")
			}
		case OpPush:
			stackEnv.Push(env._get(opa))
		case OpRet:
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
		case OpLoadFunc:
			env.A = Function(env.Global.Functions[opa])
		case OpCallMap:
			cls := env._get(opa).MustFunc("kwargs invoke operator", 0)
			m := buildCallMap(cls, stackEnv)
			stackEnv.Clear()
			for _, pa := range cls.Params {
				stackEnv.Push(m.Get(String(pa)))
			}
			stackEnv.A = m.Value()
			fallthrough
		case OpCall:
			cls := env._get(opa).MustFunc("invoke operator:", 0)
			if cls.MethodSrc != Nil {
				stackEnv.Prepend(cls.MethodSrc)
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
				stackEnv.Push(watermark) // Used by debug.dumpstk to determine the top of stack
				if bop == OpCallMap {
					stackEnv.Push(stackEnv.A)
				}
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
				env.Global = cls.loadGlobal
				env.A = stackEnv.A

				if opb == 0 {
					retStack = append(retStack, last)
				}

				if env.Global.MaxCallStackSize > 0 && int64(len(retStack)) > env.Global.MaxCallStackSize {
					panicf("call stack overflow, max: %d", env.Global.MaxCallStackSize)
				}

				stackEnv.StackOffset = uint32(len(*env.stack))
			}
		case OpJmp:
			cursor = uint32(int32(cursor) + int32(v&0xffffff) - 1<<23)
		case OpIfNot:
			if env.A.IsFalse() {
				cursor = uint32(int32(cursor) + int32(v&0xffffff) - 1<<23)
			}
		case OpIf:
			if !env.A.IsFalse() {
				cursor = uint32(int32(cursor) + int32(v&0xffffff) - 1<<23)
			}
		}
	}
}

func buildCallMap(f *Func, stackEnv Env) *RHMap {
	m := NewArrayMap(stackEnv.Size() / 2)
	for i := 0; i < stackEnv.Size(); i += 2 {
		a := stackEnv.Stack()[i]
		if a == Nil {
			continue
		}
		var name string
		if a.Type() == VNumber && a.Int() < int64(len(f.Params)) {
			name = f.Params[a.Int()]
		} else {
			name = a.String()
		}
		nv := String(name)
		if m.findHash(nv) >= 0 {
			panicf("call: duplicated parameter: %q", name)
		}
		m.Set(nv, stackEnv.Get(i+1))
	}
	return m
}
