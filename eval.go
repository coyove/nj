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

// ExecCursor executes 'K' under 'Env' from the given start 'cursor'
func ExecCursor(env *Env, K *Closure, cursor uint32) (result, resultB Value, nextCursor uint32, yielded bool) {
	var stackEnv *Env
	var retStack []stacktrace
	var recycledStacks []*Env
	var caddr = kodeaddr(K.Code)

	defer func() {
		if r := recover(); r != nil {
			stk := append(retStack, stacktrace{cls: K})
			for i := len(stk) - 1; i >= 0; i-- {
				if stk[i].cls.Is(ClsRecoverable) {
					nextCursor, yielded = 0, false
					if rv, ok := r.(Value); ok {
						result = rv
						resultB = env.B
					} else {
						p := bytes.Buffer{}
						fmt.Fprint(&p, r)
						result = Str(p.String())
					}
					return
				}
			}

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
		r.env.A = v
		r.env.B = env.B
		caddr = kodeaddr(K.Code)
		if r.cls.Is(ClsNoEnvescape) {
			if stackEnv != nil {
				for i := range stackEnv.stack {
					stackEnv.stack[i] = Value{}
				}
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
		//v := K.Code[cursor]
		cursor++
		bop, opa, opb := op(v)

		switch bop {
		case OpEOB:
			break MAIN
		case OpNOP:
		case OpSet:
			env._set(opa, env._get(opb, K))
		case OpGetB:
			env.A = env.B
		case OpSetB:
			env.B = env._get(opa, K)
		case OpInc:
			env.A = Num(env._get(opa, K).Expect(NUM).Num() + env._get(opb, K).Expect(NUM).Num())
			env._set(opa, env.A)
		case OpConcat:
			switch va, vb := env._get(opa, K), env._get(opb, K); va.Type() + vb.Type() {
			case StrStr:
				env.A = Str(va.Str() + vb.Str())
			default:
				env.A, env.B = va.DummyTable().mm("__concat").ExpectMsg(FUN, "operator ..").Fun().Call(va, vb)
			}
		case OpAdd:
			switch va, vb := env._get(opa, K), env._get(opb, K); va.Type() + vb.Type() {
			case NumNum:
				env.A = Num(va.Num() + vb.Num())
			default:
				env.A, env.B = va.DummyTable().mm("__add").ExpectMsg(FUN, "operator +").Fun().Call(va, vb)
			}
		case OpSub:
			switch va, vb := env._get(opa, K), env._get(opb, K); va.Type() + vb.Type() {
			case NumNum:
				env.A = Num(va.Num() - vb.Num())
			default:
				env.A, env.B = va.DummyTable().mm("__sub").ExpectMsg(FUN, "operator -").Fun().Call(va, vb)
			}
		case OpMul:
			switch va, vb := env._get(opa, K), env._get(opb, K); va.Type() + vb.Type() {
			case NumNum:
				env.A = Num(va.Num() * vb.Num())
			default:
				env.A, env.B = va.DummyTable().mm("__mul").ExpectMsg(FUN, "operator *").Fun().Call(va, vb)
			}
		case OpDiv:
			switch va, vb := env._get(opa, K), env._get(opb, K); va.Type() + vb.Type() {
			case NumNum:
				env.A = Num(va.Num() / vb.Num())
			default:
				env.A, env.B = va.DummyTable().mm("__div").ExpectMsg(FUN, "operator /").Fun().Call(va, vb)
			}
		case OpMod:
			switch va, vb := env._get(opa, K), env._get(opb, K); va.Type() + vb.Type() {
			case NumNum:
				env.A = Num(math.Remainder(va.Num(), vb.Num()))
			default:
				env.A, env.B = va.DummyTable().mm("__mod").ExpectMsg(FUN, "operator %").Fun().Call(va, vb)
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
				env.A, env.B = va.DummyTable().mm("__lt").ExpectMsg(FUN, "operator <").Fun().Call(va, vb)
			}
		case OpLessEq:
			switch va, vb := env._get(opa, K), env._get(opb, K); va.Type() + vb.Type() {
			case NumNum:
				env.A = Bln(va.Num() <= vb.Num())
			case StrStr:
				env.A = Bln(va.Str() <= vb.Str())
			default:
				env.A, env.B = va.DummyTable().mm("__le").ExpectMsg(FUN, "operator <=").Fun().Call(va, vb)
			}
		case OpNot:
			env.A = Bln(env._get(opa, K).IsFalse())
		case OpPow:
			switch va, vb := env._get(opa, K), env._get(opb, K); va.Type() + vb.Type() {
			case NumNum:
				a, b := va.Num(), vb.Num()
				if a == 2 && b > 0 && float64(int(b)) == b {
					env.A = Num(float64(uint(2) << uint(b-1)))
				} else {
					env.A = Num(math.Pow(a, b))
				}
			default:
				env.A, env.B = va.DummyTable().mm("__pow").ExpectMsg(FUN, "operator ^").Fun().Call(va, vb)
			}
		case OpLen:
			switch v := env._get(opa, K); v.Type() {
			case STR:
				env.A = Num(float64(len(v.Str())))
			case TAB:
				t := v.Tab()
				if t.mm("__len").Type() == FUN {
					env.A, env.B = t.mm("__len").Fun().Call(v)
				} else {
					env.A = Num(float64(t.Len()))
				}
			case FUN:
				env.A = Num(float64(v.Fun().NumParam))
			default:
				env.A, env.B = v.DummyTable().mm("__len").ExpectMsg(FUN, "operator #").Fun().Call(v)
			}
		case OpMakeHash:
			if stackEnv == nil {
				env.A = Tab(&Table{})
			} else {
				var m *Table
				if opa == 1 {
					m = env.A.Tab()
				} else {
					m = &Table{}
				}
				for i := 0; i < stackEnv.Size(); i += 2 {
					m.Put(stackEnv.stack[i], stackEnv.stack[i+1], true)
				}
				stackEnv.Clear()
				env.A = Tab(m)
			}
		case OpMakeArray:
			if stackEnv == nil {
				env.A = Tab(&Table{})
			} else {
				m := &Table{a: make([]Value, 0, len(stackEnv.stack))}
				for _, v := range stackEnv.stack {
					if v.Type() == UPK {
						m.a = append(m.a, v.asUnpacked()...)
					} else {
						m.a = append(m.a, v)
					}
				}
				stackEnv.Clear()
				env.A = Tab(m)
			}
		case OpStore:
			subject, v := env._get(opa, K), env._get(opb, K)
			switch subject.Type() {
			case TAB:
				subject.Tab().Put(env.A, v, false)
			case NIL:
				switch env.A.Type() {
				case NUM:
					x := math.Float64bits(env.A.Num())
					(*Env)(unsafe.Pointer(uintptr(x<<16>>16)))._set(uint16(x>>48), v)
				case NIL:
					// ignore
				default:
					panicf("%#v: address[] = value, not an address", env.A)
				}
			default:
				env.A, env.B = subject.DummyTable().mm("__newindex").ExpectMsg(FUN, "store").Fun().Call(subject, env.A, v)
			}
			env.A = v
		case OpLoad:
			switch a, idx := env._get(opa, K), env._get(opb, K); a.Type() {
			case TAB:
				env.A = a.Tab().Get(idx, false)
			default:
				env.A, env.B = a.DummyTable().mm("__index").ExpectMsg(FUN, "load").Fun().Call(a, idx)
			}
		case OpPush:
			if stackEnv == nil {
				stackEnv = NewEnv(nil)
			}
			stackEnv.Push(env._get(opa, K))
		case OpPush2:
			if stackEnv == nil {
				stackEnv = NewEnv(nil)
			}
			stackEnv.Push(env._get(opa, K))
			stackEnv.Push(env._get(opb, K))
		case OpRet:
			v := env._get(opa, K)
			if len(retStack) == 0 {
				return v, env.B, 0, false
			}
			returnUpperWorld(v)
		case OpYield:
			return env._get(opa, K), env.B, cursor, true
		case OpLambda:
			env.A = Fun(crReadClosure(K.Code, &cursor, env, opa, opb))
		case OpCall:
			var cls *Closure
			switch a := env._get(opa, K); a.Type() {
			case FUN:
				cls = a.Fun()
			default:
				cls = a.DummyTable().mm("__call").ExpectMsg(FUN, "invoke").Fun()
				stackEnv.stack = append([]Value{a}, stackEnv.stack...)
			}
			if cls.lastenv != nil {
				env.A, env.B = cls.Exec(nil)
				stackEnv = nil
			} else {
				if stackEnv == nil {
					stackEnv = NewEnv(env)
				}

				if stackEnv.Size() != int(cls.NumParam) {
					if !(cls.Is(ClsVararg) && stackEnv.Size() > int(cls.NumParam)) {
						panicf("expect %d arguments (got %d)", cls.NumParam, stackEnv.Size())
					}
				}

				if cls.Is(ClsYieldable | ClsRecoverable | ClsNative) {
					stackEnv.parent = env
					env.A, env.B = cls.Exec(stackEnv)
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

					retStack = append(retStack, last)
				}

				if cls.native == nil {
					if len(recycledStacks) == 0 {
						stackEnv = nil
					} else {
						stackEnv = recycledStacks[len(recycledStacks)-1]
						recycledStacks = recycledStacks[:len(recycledStacks)-1]
					}
				} else {
					stackEnv.Clear()
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
		case OpPatchVararg:
			ret := &Table{}
			for i, v := range env.stack {
				if v.Type() == UPK {
					if i != len(env.stack)-1 {
						panicf("misuse of unpack(...): it should be the last argument")
					}
					ret.a = v.asUnpacked()
					env.stack = env.stack[:len(env.stack)-1]
					break
				}
				if i >= int(K.NumParam) {
					ret.a = append(ret.a, v)
				}
			}
			env.A = Tab(ret)
		case OpAddressOf:
			addr := uint64(opa)<<48 | uint64(uintptr(unsafe.Pointer(env)))
			env.A = Num(math.Float64frombits(addr))
		}
	}

	if len(retStack) > 0 {
		returnUpperWorld(Value{})
		goto MAIN
	}
	return Value{}, Value{}, 0, false
}