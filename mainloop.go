package potatolang

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"unsafe"
)

const max32 = 0xffffffff

func panicf(msg string, args ...interface{}) {
	panic(fmt.Sprintf(msg, args...))
}

func panicerr(err error) {
	if err != nil {
		panic(err)
	}
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
			var opx uint32 = max32
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
func ExecCursor(env *Env, K *Closure, cursor uint32) (result Value, nextCursor uint32, yielded bool) {
	var newEnv *Env
	var retStack []stacktrace
	var caddr = kodeaddr(K.Code)

	defer func() {
		if r := recover(); r != nil {
			if K.Isset(ClsRecoverable) {
				nextCursor, yielded = 0, false
				if rv, ok := r.(Value); ok {
					result = rv
				} else {
					result = NewStringValue(fmt.Sprintf("%v", r))
				}
				return
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
		caddr = kodeaddr(K.Code)
		if r.cls.Isset(ClsNoEnvescape) {
			newEnv = env
			newEnv.LocalClear()
		}
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
		case OpInc:
			env.A.SetNumberValue(env.Get(opa, K).MustNumber() + env.Get(opb, K).MustNumber())
			env.Set(opa, env.A)
		case OpAdd:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); combineTypes(va, vb) {
			case _NumberNumber:
				env.A.SetNumberValue(va.AsNumber() + vb.AsNumber())
			case _StringString:
				env.A = NewStringValue(va.AsString() + vb.AsString())
			default:
				panicf("can't apply '+' on %+v and %+v", va, vb)
			}
		case OpSub:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); combineTypes(va, vb) {
			case _NumberNumber:
				env.A.SetNumberValue(va.AsNumber() - vb.AsNumber())
			default:
				panicf("can't apply '-' on %+v and %+v", va, vb)
			}
		case OpMul:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); combineTypes(va, vb) {
			case _NumberNumber:
				env.A.SetNumberValue(va.AsNumber() * vb.AsNumber())
			default:
				panicf("can't apply '*' on %+v and %+v", va, vb)
			}
		case OpDiv:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); combineTypes(va, vb) {
			case _NumberNumber:
				env.A.SetNumberValue(va.AsNumber() / vb.AsNumber())
			default:
				panicf("can't apply '/' on %+v and %+v", va, vb)
			}
		case OpMod:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); combineTypes(va, vb) {
			case _NumberNumber:
				env.A.SetNumberValue(math.Remainder(va.AsNumber(), vb.AsNumber()))
			default:
				panicf("can't apply '%%' on %+v and %+v", va, vb)
			}
		case OpEq:
			env.A.SetBoolValue(env.Get(opa, K).Equal(env.Get(opb, K)))
		case OpNeq:
			env.A.SetBoolValue(!env.Get(opa, K).Equal(env.Get(opb, K)))
		case OpLess:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); combineTypes(va, vb) {
			case _NumberNumber:
				env.A.SetBoolValue(va.AsNumber() < vb.AsNumber())
			case _StringString:
				env.A.SetBoolValue(va.AsString() < vb.AsString())
			default:
				panicf("can't apply '<' on %+v and %+v", env.Get(opa, K), env.Get(opb, K))
			}
		case OpLessEq:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); combineTypes(va, vb) {
			case _NumberNumber:
				env.A.SetBoolValue(va.AsNumber() <= vb.AsNumber())
			case _StringString:
				env.A.SetBoolValue(va.AsString() <= vb.AsString())
			default:
				panicf("can't apply '<=' on %+v and %+v", env.Get(opa, K), env.Get(opb, K))
			}
		case OpNot:
			env.A.SetBoolValue(env.Get(opa, K).IsFalse())
		case OpBitAnd:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); combineTypes(va, vb) {
			case _NumberNumber:
				env.A.SetNumberValue(float64(va.AsInt32() & vb.AsInt32()))
			default:
				panicf("can't apply '&' on %+v and %+v", env.Get(opa, K), env.Get(opb, K))
			}
		case OpBitOr:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); combineTypes(va, vb) {
			case _NumberNumber:
				env.A.SetNumberValue(float64(va.AsInt32() | vb.AsInt32()))
			default:
				panicf("can't apply '|' on %+v and %+v", env.Get(opa, K), env.Get(opb, K))
			}
		case OpBitXor:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); combineTypes(va, vb) {
			case _NumberNumber:
				env.A.SetNumberValue(float64(va.AsInt32() ^ vb.AsInt32()))
			default:
				panicf("can't apply '^' on %+v and %+v", env.Get(opa, K), env.Get(opb, K))
			}
		case OpBitLsh:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); combineTypes(va, vb) {
			case _NumberNumber:
				env.A.SetNumberValue(float64(va.AsInt32() << uint(vb.AsNumber())))
			case _MapMap:
				{
					va, vb := va.AsMap(), vb.AsMap()
					va.l = append(va.l, vb.l...)
					if va.m == nil && vb.m != nil {
						va.m = make(map[interface{}]Value, len(vb.m))
						for k, v := range vb.m {
							va.m[k] = v
						}
					}
				}
				env.A = va
			default:
				panicf("can't apply '<<' on %+v and %+v", env.Get(opa, K), env.Get(opb, K))
			}
		case OpBitRsh:
			if va, vb := env.Get(opa, K), env.Get(opb, K); combineTypes(va, vb) == _NumberNumber {
				env.A.SetNumberValue(float64(va.AsInt32() >> uint(vb.AsNumber())))
			} else {
				panicf("can't apply '>>' on %+v and %+v", env.Get(opa, K), env.Get(opb, K))
			}
		case OpBitURsh:
			if va, vb := env.Get(opa, K), env.Get(opb, K); combineTypes(va, vb) == _NumberNumber {
				env.A.SetNumberValue(float64(uint32(uint64(va.AsNumber())&max32) >> uint(vb.AsNumber())))
			} else {
				panicf("can't apply '>>>' on %+v and %+v", env.Get(opa, K), env.Get(opb, K))
			}
		case OpAssert:
			if a := env.Get(opa, K); a.IsFalse() {
				msg := env.Get(opb, K)
				if msg.Type() == StringType {
					panic(msg.AsString())
				}
				panicf("assertion failed: %+v", a)
			}
			env.A = NewBoolValue(true)
		case OpLen:
			switch v := env.Get(opa, K); v.Type() {
			case StringType:
				env.A.SetNumberValue(float64(len(v.AsString())))
			case MapType:
				env.A.SetNumberValue(float64(v.AsMap().Size()))
			default:
				panicf("can't evaluate the length of %+v", v)
			}
		case OpMakeMap:
			if newEnv == nil {
				env.A = NewMapValue(NewMap())
			} else {
				size, m := newEnv.LocalSize(), NewMap()
				for i := 0; i < size; i += 2 {
					m.Put(newEnv.LocalGet(i), newEnv.LocalGet(i+1))
				}
				newEnv.LocalClear()
				env.A = NewMapValue(m)
			}
		case OpMakeArray:
			if newEnv == nil {
				env.A = NewMapValue(NewMap())
			} else {
				m := NewMapSize(newEnv.LocalSize())
				copy(m.l, newEnv.stack)
				newEnv.LocalClear()
				env.A = NewMapValue(m)
			}
		case OpStore:
			vidx, v := env.Get(opa, K), env.Get(opb, K)
			switch env.A.Type() {
			case MapType:
				m := env.A.AsMap()
				if v == Phantom {
					m.Remove(vidx)
				} else if vidx.Type() == NumberType {
					idx, ln := int(vidx.AsNumber()), len(m.l)
					if idx < ln {
						m.l[idx] = v
					} else if idx == ln {
						m.l = append(m.l, v)
					} else {
						m.putIntoMap(vidx, v)
					}
				} else {
					m.putIntoMap(vidx, v)
				}
			case StringType:
				var p []byte
				switch combineTypes(vidx, v) {
				case _NumberNumber:
					p = []byte(env.A.AsString())
					p[int(vidx.AsNumber())] = byte(v.AsNumber())
				case NumberType<<8 | StringType:
					idx, as, vs := int(vidx.AsNumber()), env.A.AsString(), v.AsString()
					if len(vs) == 1 {
						p = []byte(as)
						p[idx] = vs[0]
					} else {
						p = make([]byte, len(as)+len(vs)-1)
						copy(p, as[:idx])
						copy(p[idx:], vs)
						copy(p[idx+len(vs):], as[idx+1:])
					}
				default:
					panicf("can't modify string %+v[%+v] to %+v", env.A, vidx, v)
				}
				(*Map)(env.A.ptr).ptr = unsafe.Pointer(&p) // unsafely cast p to string
			case NilType:
				switch va := env.Get(opa, K); va.Type() {
				case NumberType:
					x := math.Float64bits(va.AsNumber())
					(*Env)(unsafe.Pointer(uintptr(x<<16>>16))).Set(uint16(x>>48), v)
				case NilType:
					// ignore
				default:
					panicf("%+v: move(address, value), not an address", va)
				}
			default:
				panicf("can't modify %+v[%+v] to %+v", env.A, vidx, v)
			}
			env.A = v
		case OpLoad:
			var v Value
			a := env.Get(opa, K)
			vidx := env.Get(opb, K)
			switch combineTypes(a, vidx) {
			case _StringNumber:
				v.SetNumberValue(float64(a.AsString()[int(vidx.AsNumber())]))
			case _MapNumber:
				if m, idx := a.AsMap(), int(vidx.AsNumber()); idx < len(m.l) {
					v = m.l[idx]
					break
				}
				fallthrough
			default:
				if a.Type() == MapType {
					v, _ = a.AsMap().getFromMap(vidx)
					if v.Type() == ClosureType {
						if cls := v.AsClosure(); cls.Isset(ClsHasReceiver) {
							cls = cls.Dup()
							if cls.ArgsCount > 0 {
								if len(cls.PartialArgs) > 0 {
									panicf("curry function with a receiver")
								}
								cls.ArgsCount--
								cls.PartialArgs = []Value{a}
							}
							v = NewClosureValue(cls)
						}
					}
				} else {
					panicf("can't load %+v[%+v]", a, vidx)
				}
			}
			env.A = v
		case OpPop:
			a := env.Get(opa, K)
			switch a.Type() {
			case MapType:
				m := a.AsMap()
				l := m.l
				if len(l) == 0 {
					env.A = Value{}
				} else {
					env.A = l[len(l)-1]
					m.l = l[:len(l)-1]
				}
			case NilType:
				env.A = Phantom
			default:
				panicf("can't pop %+v", a)
			}
		case OpSlice:
			start, end := int(env.Get(opa, K).MustNumber()), int(env.Get(opb, K).MustNumber())
			switch x := env.A; x.Type() {
			case StringType:
				if end == -1 {
					env.A = NewStringValue(x.AsString()[start:])
				} else {
					env.A = NewStringValue(x.AsString()[start:end])
				}
			case MapType:
				m := NewMap()
				if end == -1 {
					m.l = x.AsMap().l[start:]
				} else {
					m.l = x.AsMap().l[start:end]
				}
				env.A = NewMapValue(m)
			default:
				panicf("can't slice %+v", x)
			}
		case OpPush:
			if newEnv == nil {
				newEnv = NewEnv(nil)
			}
			newEnv.LocalPush(env.Get(opa, K))
		case OpPush2:
			if newEnv == nil {
				newEnv = NewEnv(nil)
			}
			newEnv.LocalPush(env.Get(opa, K))
			newEnv.LocalPush(env.Get(opb, K))
		case OpRet:
			v := env.Get(opa, K)
			if len(retStack) == 0 {
				return v, 0, false
			}
			returnUpperWorld(v)
		case OpYield:
			return env.Get(opa, K), cursor, true
		case OpLambda:
			env.A = NewClosureValue(crReadClosure(K.Code, &cursor, env, opa, opb))
		case OpCall:
			v := env.Get(opa, K)
			if x := v.Type(); x != ClosureType {
				if x == MapType {
					if newEnv != nil && newEnv.LocalSize() > 0 {
						dest := newEnv.LocalGet(0).MustMap()
						n := copy(dest.l, v.AsMap().l)
						m := NewMap()
						m.l = dest.l[:n]
						env.A = NewMapValue(m)
						newEnv.LocalClear()
					} else {
						env.A = NewMapValue(v.AsMap().Dup())
					}
					continue
				}
				v.panicType(ClosureType)
			}
			cls := v.AsClosure()
			if cls.lastenv != nil {
				env.A = cls.Exec(nil)
				newEnv = nil
			} else {
				if newEnv == nil {
					newEnv = NewEnv(env)
				}

				if newEnv.LocalSize() >= int(cls.ArgsCount) {
					if len(cls.PartialArgs) > 0 {
						newEnv.LocalPushFront(cls.PartialArgs)
					}
					if cls.Isset(ClsYieldable | ClsRecoverable | ClsNative) {
						newEnv.parent = env
						env.A = cls.Exec(newEnv)
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
						newEnv.parent = cls.Env
						env = newEnv

						retStack = append(retStack, last)
					}
					if cls.native == nil {
						newEnv = nil
					} else {
						newEnv.LocalClear()
					}
				} else if newEnv.LocalSize() == 0 {
					env.A = NewClosureValue(cls.Dup())
				} else {
					curry := cls.Dup()
					curry.AppendPartialArgs(newEnv.Stack())
					env.A = NewClosureValue(curry)
					newEnv.LocalClear()
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
		case OpForeach:
			x := env.Get(opa, K)
			if x.Type() == NilType {
				ret := NewMapSize(len(env.stack))
				copy(ret.l, env.stack)
				env.A = NewMapValue(ret)
				continue
			}
			m := x.MustMap()
			cls := env.Get(opb, K).MustClosure()
			forEnv := NewEnv(cls.Env)
			for i := len(m.l) - 1; i >= 0; i-- {
				forEnv.LocalClear()
				forEnv.LocalPush(NewNumberValue(float64(i)))
				forEnv.LocalPush(m.l[i])
				if res := cls.Exec(forEnv); res.IsZero() {
					continue MAIN
				}
			}
			for k, v := range m.m {
				forEnv.LocalClear()
				forEnv.LocalPush(NewInterfaceValue(k))
				forEnv.LocalPush(v)
				if res := cls.Exec(forEnv); res.IsZero() {
					continue MAIN
				}
			}
			env.A = Value{}
		case OpTypeof:
			v := env.Get(opa, K)
			if v == Phantom {
				env.A = NewStringValue("#nil")
			} else {
				env.A = NewStringValue(TMapping[v.Type()])
			}
		case OpAddressOf:
			addr := uint64(opa)<<48 | uint64(uintptr(unsafe.Pointer(env)))
			env.A.SetNumberValue(math.Float64frombits(addr))
		}
	}

	if len(retStack) > 0 {
		returnUpperWorld(Value{})
		goto MAIN
	}
	return Value{}, 0, false
}
