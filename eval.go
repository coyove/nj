package potatolang

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"unsafe"
)

type stacktrace struct {
	cursor uint32
	env    *Env
	cls    *Closure
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
		src := "<unknown>"
		for i := 0; i < len(r.cls.Pos); {
			var op, line uint32
			var opx uint32 = math.MaxUint32
			var col uint16
			i, op, line, col = r.cls.Pos.readABC(i)
			if i < len(r.cls.Pos)-1 {
				_, opx, _, _ = r.cls.Pos.readABC(i)
			}
			if r.cursor >= op && r.cursor < opx {
				src = fmt.Sprintf("%s:%d:%d", r.cls.source, line, col)
				break
			}
		}
		// the recorded cursor was advanced by 1 already
		msg.WriteString(fmt.Sprintf("at %d in %s\n", r.cursor-1, src))
	}
	msg.WriteString("root panic:\n")
	msg.WriteString(fmt.Sprintf("%v\n", e.r))
	return msg.String()
}

func kodeaddr(code []uint32) uintptr { return (*reflect.SliceHeader)(unsafe.Pointer(&code)).Data }

func konstaddr(consts []Value) uintptr { return (*reflect.SliceHeader)(unsafe.Pointer(&consts)).Data }

func returnVararg(a Value, b []Value) (Value, []Value) {
	if len(b) == 0 {
		if a.Type() == UPK {
			u := a._Upk()
			if len(u) == 0 {
				return Value{}, nil
			}
			return u[0], u[1:]
		}
		return a, nil
	}

	count, flag := a._TestUpkLen()
	for _, b := range b {
		c, f := b._TestUpkLen()
		count, f = count+c, flag || f
	}
	if !flag { // both 'a' and 'b' are not (neither containing) unpacked values
		return a, b
	}

	b2 := make([]Value, 0, count)
	b2 = a._AppendTo(b2)
	for _, v := range b {
		b2 = v._AppendTo(b2)
	}
	if len(b2) == 0 {
		return Value{}, nil
	}
	return b2[0], b2[1:]
}

// execCursorLoop executes 'K' under 'Env' from the given start 'cursor'
func execCursorLoop(env *Env, K *Closure, cursor uint32) (result Value, resultV []Value, nextCursor uint32, yielded bool) {
	var stackEnv *Env
	var retStack []stacktrace
	var recycledStacks []*Env
	var caddr = kodeaddr(K.Code)

	defer func() {
		if r := recover(); r != nil {
			// stk := append(retStack, stacktrace{cls: K})
			// for i := len(stk) - 1; i >= 0; i-- {
			// 	if stk[i].cls.Is(ClsRecoverable) {
			// 		nextCursor, yielded = 0, false
			// 		if rv, ok := r.(Value); ok {
			// 			result = rv
			// 			resultV = env.V
			// 		} else {
			// 			p := bytes.Buffer{}
			// 			fmt.Fprint(&p, r)
			// 			result = Str(p.String())
			// 		}
			// 		return
			// 	}
			// }

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

	returnUpperWorld := func(v Value) {
		r := retStack[len(retStack)-1]
		cursor = r.cursor
		K = r.cls
		r.env.A, r.env.V = returnVararg(v, env.V)
		caddr = kodeaddr(K.Code)
		if r.cls.Is(ClsNoEnvEscape) {
			if stackEnv != nil {
				for i := range stackEnv.stack {
					stackEnv.stack[i] = Value{}
				}
				stackEnv.stack = stackEnv.stack[:0]
				recycledStacks = append(recycledStacks, stackEnv)
			}
			stackEnv = env
			stackEnv.Clear()
		}
		// log.Println(unsafe.Pointer(&stackEnv.stack))
		env = r.env
		retStack = retStack[:len(retStack)-1]
	}

MAIN:
	for {
		//	if flag != nil && atomic.LoadUintptr(flag) == 1 {
		//		panicf("canceled")
		//	}
		v := *(*uint32)(unsafe.Pointer(uintptr(cursor)*4 + caddr))
		cursor++
		bop, opa, opb := op(v)

		switch bop {
		case OpEOB:
			break MAIN
		case OpSet:
			env._set(opa, env._get(opb, K))
		case OpPushV:
			if opb != 0 {
				env.V = make([]Value, 0, opb)
			}
			env.V = append(env.V, env._get(opa, K))
		case OpPopV:
			switch opa {
			case 4: // popv-all-with-a
				if env.A.IsNil() {
					env.A, env.V = newUnpackedValue(env.V), nil
				} else {
					env.A, env.V = newUnpackedValue(append([]Value{env.A}, env.V...)), nil
				}
			case 3: // popv-clear
				env.V = nil
			case 2: // popv-all-and-clear
				if len(env.V) == 0 {
					env.A = newUnpackedValue(nil)
				} else {
					env.A, env.V = newUnpackedValue(env.V), nil
				}
			case 1: // popv
				if len(env.V) == 0 {
					env.A = Value{}
				} else {
					env.A, env.V = env.V[0], env.V[1:]
				}
			case 0: // popv-last-and-clear
				if len(env.V) == 0 {
					env.A = Value{}
				} else {
					env.A, env.V = env.V[0], env.V[1:]
				}
				env.V = nil
			}
		case OpInc:
			env.A = Num(env._get(opa, K).Expect(NUM).Num() + env._get(opb, K).Expect(NUM).Num())
			env._set(opa, env.A)
		case OpConcat:
			switch va, vb := env._get(opa, K), env._get(opb, K); va.Type() + vb.Type() {
			case StrStr:
				env.A = Str(va.Str() + vb.Str())
			default:
				env.A, _ = findmm(va, vb, M__concat).ExpectMsg(FUN, "metamethod operator ..").Fun().Call(va, vb)
			}
		case OpAdd:
			switch va, vb := env._get(opa, K), env._get(opb, K); va.Type() + vb.Type() {
			case NumNum:
				env.A = Num(va.Num() + vb.Num())
			default:
				env.A, _ = findmm(va, vb, M__add).ExpectMsg(FUN, "metamethod operator +").Fun().Call(va, vb)
			}
		case OpSub:
			switch va, vb := env._get(opa, K), env._get(opb, K); va.Type() + vb.Type() {
			case NumNum:
				env.A = Num(va.Num() - vb.Num())
			default:
				env.A, _ = findmm(va, vb, M__sub).ExpectMsg(FUN, "metamethod operator -").Fun().Call(va, vb)
			}
		case OpUnm:
			switch va := env._get(opa, K); va.Type() {
			case NUM:
				env.A = Num(-va.Num())
			default:
				env.A, _ = va.GetMetamethod(M__unm).ExpectMsg(FUN, "metamethod operator unary -").Fun().Call(va)
			}
		case OpMul:
			switch va, vb := env._get(opa, K), env._get(opb, K); va.Type() + vb.Type() {
			case NumNum:
				env.A = Num(va.Num() * vb.Num())
			default:
				env.A, _ = findmm(va, vb, M__mul).ExpectMsg(FUN, "metamethod operator *").Fun().Call(va, vb)
			}
		case OpDiv:
			switch va, vb := env._get(opa, K), env._get(opb, K); va.Type() + vb.Type() {
			case NumNum:
				env.A = Num(va.Num() / vb.Num())
			default:
				env.A, _ = findmm(va, vb, M__div).ExpectMsg(FUN, "metamethod operator /").Fun().Call(va, vb)
			}
		case OpMod:
			switch va, vb := env._get(opa, K), env._get(opb, K); va.Type() + vb.Type() {
			case NumNum:
				env.A = Num(math.Remainder(va.Num(), vb.Num()))
			default:
				env.A, _ = findmm(va, vb, M__mod).ExpectMsg(FUN, "metamethod operator %").Fun().Call(va, vb)
			}
		case OpEq:
			env.A = Bln(env._get(opa, K).Equal(env._get(opb, K)))
		case OpNeq:
			env.A = Bln(!env._get(opa, K).Equal(env._get(opb, K)))
		case OpLess:
			switch va, vb := env._get(opa, K), env._get(opb, K); va.Type() + vb.Type() {
			case NumNum:
				env.A = Bln(va.Num() < vb.Num())
			case StrStr:
				env.A = Bln(va.Str() < vb.Str())
			default:
				env.A, _ = findmm(va, vb, M__lt).ExpectMsg(FUN, "metamethod operator <").Fun().Call(va, vb)
			}
		case OpLessEq:
			switch va, vb := env._get(opa, K), env._get(opb, K); va.Type() + vb.Type() {
			case NumNum:
				env.A = Bln(va.Num() <= vb.Num())
			case StrStr:
				env.A = Bln(va.Str() <= vb.Str())
			default:
				env.A, _ = findmm(va, vb, M__le).ExpectMsg(FUN, "metamethod operator <=").Fun().Call(va, vb)
			}
		case OpNot:
			env.A = Bln(env._get(opa, K).IsFalse())
		case OpPow:
			switch va, vb := env._get(opa, K), env._get(opb, K); va.Type() + vb.Type() {
			case NumNum:
				env.A = Num(math.Pow(va.Num(), vb.Num()))
			default:
				env.A, _ = findmm(va, vb, M__pow).ExpectMsg(FUN, "metamethod operator ^").Fun().Call(va, vb)
			}
		case OpLen:
			switch v := env._get(opa, K); v.Type() {
			case STR:
				if f := v.GetMetamethod(M__len); f.Type() == FUN {
					env.A, _ = f.Fun().Call(v)
				} else {
					env.A = Num(float64(len(v.Str())))
				}
			case TAB:
				t := v.Tab()
				if l := t.mt.RawGet(M__len); l.Type() == FUN {
					env.A, _ = l.Fun().Call(v)
				} else {
					env.A = Num(float64(t.Len()))
				}
			case FUN:
				if f := v.GetMetamethod(M__len); f.Type() == FUN {
					env.A, _ = f.Fun().Call(v)
				} else {
					env.A = Num(float64(v.Fun().NumParam))
				}
			case UPK:
				env.A = Num(float64(len(v._Upk())))
			default:
				env.A, _ = v.GetMetamethod(M__len).ExpectMsg(FUN, "metamethod operator #").Fun().Call(v)
			}
		case OpMakeTable:
			if stackEnv == nil {
				env.A = Tab(&Table{})
			} else {
				switch opa {
				case 1, 3: // 1: make hash; 3: make hash part of the table stored in $a
					var m *Table
					if opa == 3 {
						m = env.A.Tab()
					} else {
						m = &Table{}
					}
					for i := 0; i < stackEnv.Size(); i += 2 {
						m.RawPut(stackEnv.stack[i], stackEnv.stack[i+1])
					}
					stackEnv.Clear()
					env.A = Tab(m)
				case 2: // 2: make array
					m := &Table{a: append([]Value{}, stackEnv.stack...)}
					stackEnv.Clear()
					env.A = Tab(m)
				}
			}
		case OpStore:
			subject, v := env._get(opa, K), env._get(opb, K)
			switch subject.Type() {
			case TAB:
				subject.Tab().Put(env.A, v)
			case UPK:
				subject._Upk()[int(env.A.ExpectMsg(NUM, "unpacked store").Num())-1] = v
			default:
				env.A, _ = subject.GetMetamethod(M__newindex).ExpectMsg(FUN, "metamethod newindex").Fun().Call(subject, env.A, v)
			}
			env.A = v
		case OpLoad:
			switch a, idx := env._get(opa, K), env._get(opb, K); a.Type() {
			case TAB:
				env.A = a.Tab().Get(idx)
			case UPK:
				env.A = a._Upk()[int(idx.ExpectMsg(NUM, "unpacked load").Num())-1]
			default:
				env.A, _ = a.GetMetamethod(M__index).ExpectMsg(FUN, "metamethod index").Fun().Call(a, idx)
			}
		case OpPush:
			if stackEnv == nil {
				stackEnv = NewEnv(nil)
			}
			// if opb == 0 {
			// 	if len(stackEnv.stack) > 0 && opb == 0 {
			// 		log.Println(stackEnv.stack)
			// 	}
			// 	stackEnv.stack = stackEnv.stack[:0]
			// }
			if v := env._get(opa, K); v.Type() == UPK {
				stackEnv.stack = append(stackEnv.stack, v._Upk()...)
			} else {
				stackEnv.Push(v)
			}
			if opa == regA && len(env.V) > 0 {
				stackEnv.stack = append(stackEnv.stack, env.V...)
			}
		case OpRet:
			v := env._get(opa, K)
			if len(retStack) == 0 {
				v, env.V = returnVararg(v, env.V)
				return v, env.V, 0, false
			}
			returnUpperWorld(v)
		case OpYield:
			v := env._get(opa, K)
			v, env.V = returnVararg(v, env.V)
			return v, env.V, cursor, true
		case OpLambda:
			env.A = Fun(pkReadClosure(K.Code, &cursor, env, opa, opb))
		case OpCall:
			var cls *Closure
			switch a := env._get(opa, K); a.Type() {
			case FUN:
				cls = a.Fun()
			default:
				cls = a.GetMetamethod(M__call).ExpectMsg(FUN, "metamethod call").Fun()
				stackEnv.stack = append([]Value{a}, stackEnv.stack...)
			}
			if cls.lastEnv != nil { // resume yielded coroutine
				env.A, env.V = cls.exec(nil)
				if stackEnv != nil {
					stackEnv.stack = stackEnv.stack[:0]
				}
			} else {
				if stackEnv == nil {
					stackEnv = NewEnv(env)
				}

				if cls.Is(ClsVararg) && !cls.Is(ClsNative) {
					var varg []Value
					if stackEnv.Size() > int(cls.NumParam) {
						varg = append([]Value{}, stackEnv.stack[cls.NumParam:]...)
					}
					stackEnv._set(uint16(cls.NumParam), newUnpackedValue(varg))
				}

				if cls.Is(ClsYieldable | ClsNative) {
					stackEnv.parent = env
					env.A, env.V = cls.exec(stackEnv)

					if cls.Is(ClsYieldable) {
						stackEnv = nil
					} else {
						stackEnv.Clear()
					}
				} else {
					last := stacktrace{
						cls:    K,
						cursor: cursor,
						env:    env,
					}

					// switch to the Env of cls
					cursor = 0
					K = cls
					caddr = kodeaddr(K.Code)
					stackEnv.parent = cls.Env
					env = stackEnv

					if opb == 0 {
						retStack = append(retStack, last)
					}

					if len(recycledStacks) == 0 {
						stackEnv = nil
					} else {
						stackEnv = recycledStacks[len(recycledStacks)-1]
						recycledStacks = recycledStacks[:len(recycledStacks)-1]
					}
				}
			}

		case OpJmp:
			cursor = uint32(int32(cursor) + int32(opb) - 1<<12)
		case OpIfNot:
			if cond := env._get(opa, K); cond.IsFalse() {
				cursor = uint32(int32(cursor) + int32(opb) - 1<<12)
			}
		case OpIf:
			if cond := env._get(opa, K); !cond.IsFalse() {
				cursor = uint32(int32(cursor) + int32(opb) - 1<<12)
			}
		}
	}

	if len(retStack) > 0 {
		returnUpperWorld(Value{})
		goto MAIN
	}
	return Value{}, nil, 0, false
}
