package potatolang

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"unsafe"
)

func panicf(msg string, args ...interface{}) {
	panic(fmt.Sprintf(msg, args...))
}

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
						result = NewStringValue(p.Bytes())
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
			stackEnv.LocalClear()
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

		//log.Println(cursor)
		v := *(*uint32)(unsafe.Pointer(uintptr(cursor)*4 + caddr))
		//v := K.Code[cursor]
		cursor++
		bop, opa, opb := op(v)

		switch bop {
		case OpEOB:
			break MAIN
		case OpNOP:
		case OpSet:
			env.Set(opa, env.Get(opb, K))
		case OpSetFromAB:
			env.Set(opa, env.A)
			env.Set(opb, env.B)
		case OpSetB:
			env.B = env.Get(opa, K)
		case OpInc:
			env.A.SetNumberValue(env.Get(opa, K).MustNumber() + env.Get(opb, K).MustNumber())
			env.Set(opa, env.A)
		case OpAdd:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); va.Type() + vb.Type() {
			case _NumberNumber:
				env.A.SetNumberValue(va.AsNumber() + vb.AsNumber())
			case _StringString:
				vab, vbb := va.AsString(), vb.AsString()
				x := make([]byte, len(vab)+len(vbb))
				copy(x[copy(x, vab):], vbb)
				env.A = NewStringValue(x)
			default:
				panicf("can't apply '+' on %#v and %#v", va, vb)
			}
		case OpSub:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); va.Type() + vb.Type() {
			case _NumberNumber:
				env.A.SetNumberValue(va.AsNumber() - vb.AsNumber())
			default:
				panicf("can't apply '-' on %#v and %#v", va, vb)
			}
		case OpMul:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); va.Type() + vb.Type() {
			case _NumberNumber:
				env.A.SetNumberValue(va.AsNumber() * vb.AsNumber())
			default:
				panicf("can't apply '*' on %#v and %#v", va, vb)
			}
		case OpDiv:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); va.Type() + vb.Type() {
			case _NumberNumber:
				env.A.SetNumberValue(va.AsNumber() / vb.AsNumber())
			default:
				panicf("can't apply '/' on %#v and %#v", va, vb)
			}
		case OpMod:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); va.Type() + vb.Type() {
			case _NumberNumber:
				env.A.SetNumberValue(math.Remainder(va.AsNumber(), vb.AsNumber()))
			default:
				panicf("can't apply '%%' on %#v and %#v", va, vb)
			}
		case OpEq:
			env.A.SetBoolValue(env.Get(opa, K).Equal(env.Get(opb, K)))
		case OpNeq:
			env.A.SetBoolValue(!env.Get(opa, K).Equal(env.Get(opb, K)))
		case OpLess:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); va.Type() + vb.Type() {
			case _NumberNumber:
				env.A.SetBoolValue(va.AsNumber() < vb.AsNumber())
			case _StringString:
				env.A.SetBoolValue(bytes.Compare(va.AsString(), vb.AsString()) == -1)
			default:
				panicf("can't apply '<' on %#v and %#v", env.Get(opa, K), env.Get(opb, K))
			}
		case OpLessEq:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); va.Type() + vb.Type() {
			case _NumberNumber:
				env.A.SetBoolValue(va.AsNumber() <= vb.AsNumber())
			case _StringString:
				env.A.SetBoolValue(bytes.Compare(va.AsString(), vb.AsString()) <= 0)
			default:
				panicf("can't apply '<=' on %#v and %#v", env.Get(opa, K), env.Get(opb, K))
			}
		case OpNot:
			env.A.SetBoolValue(env.Get(opa, K).IsFalse())
		case OpBitAnd:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); va.Type() + vb.Type() {
			case _NumberNumber:
				env.A.SetNumberValue(float64(va.AsInt32() & vb.AsInt32()))
			default:
				panicf("can't apply '&' on %#v and %#v", env.Get(opa, K), env.Get(opb, K))
			}
		case OpBitOr:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); va.Type() + vb.Type() {
			case _NumberNumber:
				env.A.SetNumberValue(float64(va.AsInt32() | vb.AsInt32()))
			default:
				panicf("can't apply '|' on %#v and %#v", env.Get(opa, K), env.Get(opb, K))
			}
		case OpBitXor:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); va.Type() + vb.Type() {
			case _NumberNumber:
				env.A.SetNumberValue(float64(va.AsInt32() ^ vb.AsInt32()))
			default:
				panicf("can't apply '^' on %#v and %#v", env.Get(opa, K), env.Get(opb, K))
			}
		case OpBitLsh:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); va.Type() + vb.Type() {
			case _NumberNumber:
				env.A.SetNumberValue(float64(va.AsInt32() << uint(vb.AsNumber())))
			case _SliceSlice:
				vas := va.AsSlice()
				vas.l = append(vas.l, vb.AsSlice().l...)
				env.A = va
			default:
				panicf("can't apply '<<' on %#v and %#v", env.Get(opa, K), env.Get(opb, K))
			}
		case OpBitRsh:
			if va, vb := env.Get(opa, K), env.Get(opb, K); va.Type()+vb.Type() == _NumberNumber {
				env.A.SetNumberValue(float64(va.AsInt32() >> uint(vb.AsNumber())))
			} else {
				panicf("can't apply '>>' on %#v and %#v", env.Get(opa, K), env.Get(opb, K))
			}
		case OpBitURsh:
			if va, vb := env.Get(opa, K), env.Get(opb, K); va.Type()+vb.Type() == _NumberNumber {
				env.A.SetNumberValue(float64(uint32(uint64(va.AsNumber())&math.MaxUint32) >> uint(vb.AsNumber())))
			} else {
				panicf("can't apply '>>>' on %#v and %#v", env.Get(opa, K), env.Get(opb, K))
			}
		case OpLen:
			switch v := env.Get(opa, K); v.Type() {
			case StringType:
				env.A.SetNumberValue(float64(len(v.AsString())))
			case SliceType:
				env.A.SetNumberValue(float64(v.AsSlice().Size()))
			case StructType:
				env.A.SetNumberValue(float64(len(v.AsStruct().l) / 2))
			case NilType:
				env.A.SetNumberValue(0)
			case ClosureType:
				env.A.SetNumberValue(float64(v.AsClosure().ParamsCount))
			default:
				panicf("can't evaluate the length of %#v", v)
			}
		case OpMakeStruct:
			if stackEnv == nil {
				env.A = NewStructValue(NewStruct())
			} else {
				m := NewStruct()
				for i := 0; i < stackEnv.LocalSize(); i += 2 {
					m.l.Add(true, stackEnv.stack[i], stackEnv.stack[i+1])
				}
				stackEnv.LocalClear()
				env.A = NewStructValue(m)
			}
		case OpMakeSlice:
			if stackEnv == nil {
				env.A = NewSliceValue(NewSlice())
			} else {
				m := NewSliceSize(stackEnv.LocalSize())
				copy(m.l, stackEnv.stack)
				stackEnv.LocalClear()
				env.A = NewSliceValue(m)
			}
		case OpStore:
			subject, v := env.Get(opa, K), env.Get(opb, K)
			switch subject.Type() {
			case SliceType:
				subject.AsSlice().Put(int(env.A.MustNumber()), v)
			case StructType:
				if !subject.AsStruct().l.Add(false, env.A, v) {
					panicf("struct attribute %#v not found", env.A)
				}
			case NilType:
				switch env.A.Type() {
				case NumberType:
					x := math.Float64bits(env.A.AsNumber())
					(*Env)(unsafe.Pointer(uintptr(x<<16>>16))).Set(uint16(x>>48), v)
				case NilType:
					// ignore
				default:
					panicf("%#v: address[] = value, not an address", env.A)
				}
			default:
				panicf("can't modify %#v[%#v] to %#v", subject, env.A, v)
			}
			env.A = v
		case OpLoad:
			a := env.Get(opa, K)
			idx := env.Get(opb, K)
			switch uint16(a.Type())<<8 | uint16(idx.Type()) {
			case StringType<<8 | NumberType:
				env.A.SetNumberValue(float64(a.AsString()[int(idx.AsNumber())]))
			case SliceType<<8 | NumberType:
				env.A = a.AsSlice().Get(int(idx.AsNumber()))
			case StructType<<8 | NumberType:
				env.A, _ = a.AsStruct().hashGet(idx)
			default:
				panicf("can't load %#v[%#v]", a, idx)
			}
		case OpSlice:
			start, end := int(env.A.MustNumber()), int(env.Get(opb, K).MustNumber())
			switch x := env.Get(opa, K); x.Type() {
			case StringType:
				if end == -1 {
					env.A = NewStringValue(x.AsString()[start:])
				} else {
					env.A = NewStringValue(x.AsString()[start:end])
				}
			case SliceType:
				m := NewSlice()
				if end == -1 {
					m.l = x.AsSlice().l[start:]
				} else {
					m.l = x.AsSlice().l[start:end]
				}
				env.A = NewSliceValue(m)
			default:
				panicf("can't slice %#v", x)
			}
		case OpPush:
			if stackEnv == nil {
				stackEnv = NewEnv(nil)
			}
			stackEnv.LocalPush(env.Get(opa, K))
		case OpPush2:
			if stackEnv == nil {
				stackEnv = NewEnv(nil)
			}
			stackEnv.LocalPush(env.Get(opa, K))
			stackEnv.LocalPush(env.Get(opb, K))
		case OpPushVararg:
			if stackEnv == nil {
				stackEnv = NewEnv(nil)
			}
			stackEnv.stack = append(stackEnv.stack, env.Get(opa, K).MustSlice().l...)
		case OpRet:
			v := env.Get(opa, K)
			if len(retStack) == 0 {
				return v, env.B, 0, false
			}
			returnUpperWorld(v)
		case OpYield:
			return env.Get(opa, K), env.B, cursor, true
		case OpLambda:
			env.A = NewClosureValue(crReadClosure(K.Code, &cursor, env, opa, opb))
		case OpCall:
			v := env.Get(opa, K)
			if x := v.Type(); x != ClosureType {
				switch x {
				case SliceType:
					env.A = NewSliceValue(v.AsSlice().Dup())
				case StructType:
					env.A = NewStructValue(v.AsStruct().Dup())
				default:
					v.testType(ClosureType)
				}
				if stackEnv != nil {
					stackEnv.LocalClear()
				}
				continue
			}

			cls := v.AsClosure()
			if cls.lastenv != nil {
				env.A, env.B = cls.Exec(nil)
				stackEnv = nil
			} else {
				if stackEnv == nil {
					stackEnv = NewEnv(env)
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
					stackEnv.LocalClear()
				}
			}

		case OpJmp:
			cursor = uint32(int32(cursor) + int32(opb) - 1<<12)
		case OpIfNot:
			if cond := env.Get(opa, K); cond.IsFalse() {
				cursor = uint32(int32(cursor) + int32(opb) - 1<<12)
			}
		case OpIf:
			if cond := env.Get(opa, K); !cond.IsFalse() {
				cursor = uint32(int32(cursor) + int32(opb) - 1<<12)
			}
		case OpCopyStack:
			ret := NewSliceSize(len(env.stack) - int(opa))
			copy(ret.l, env.stack[opa:])
			env.A = NewSliceValue(ret)
		case OpTypeof:
			env.A = NewStringValue(typeMappings[env.Get(opa, K).Type()])
		case OpAddressOf:
			addr := uint64(opa)<<48 | uint64(uintptr(unsafe.Pointer(env)))
			env.A.SetNumberValue(math.Float64frombits(addr))
		}
	}

	if len(retStack) > 0 {
		returnUpperWorld(Value{})
		goto MAIN
	}
	return Value{}, Value{}, 0, false
}
