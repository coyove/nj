package potatolang

import (
	"bytes"
	"fmt"
	"reflect"
	"sync/atomic"
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
	cursor      uint32
	noenvescape bool
	returnInto  byte
	env         *Env
	cls         *Closure
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
		for i := 0; i < len(r.cls.pos); {
			var op, line uint32
			var opx uint32 = max32
			var col uint16
			i, op, line, col = r.cls.pos.readABC(i)
			if i < len(r.cls.pos)-1 {
				_, opx, _, _ = r.cls.pos.readABC(i)
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

// ExecCursor executes code under the given env from the given start cursor and returns:
// 1. final result 2. yield cursor 3. is yield or not
func ExecCursor(env *Env, K *Closure, cursor uint32) (_v Value, _p uint32, _y bool) {
	var newEnv *Env
	var retStack []stacktrace
	var caddr = kodeaddr(K.code)

	defer func() {
		if r := recover(); r != nil {
			if K.Isset(ClsRecoverable) {
				if rv, ok := r.(Value); ok {
					_v, _p, _y = rv, 0, false
				} else {
					_v, _p, _y = NewStringValue(fmt.Sprintf("%v", r)), 0, false
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
		caddr = kodeaddr(K.code)
		if r.noenvescape {
			newEnv = env
			newEnv.SClear()
		}
		env = r.env
		retStack = retStack[:len(retStack)-1]
	}

	flag := env.Cancel

MAIN:
	for {
		if flag != nil && atomic.LoadUintptr(flag) == 1 {
			panicf("canceled")
		}

		//log.Println(cursor)
		v := *(*uint32)(unsafe.Pointer(uintptr(cursor)*4 + caddr))
		//v := K.code[cursor]
		cursor++
		bop, opa, opb := op(v)

		switch bop {
		case OP_EOB:
			break MAIN
		case OP_NOP:
		case OP_SET:
			env.Set(opa, env.Get(opb, K))
		case OP_INC:
			num := env.Get(opa, K).AsNumber()
			env.A.SetNumberValue(num + env.Get(opb, K).AsNumber())
			env.Set(opa, env.A)
		case OP_ADD:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); testTypes(va, vb) {
			case _NumberNumber:
				env.A.SetNumberValue(va.AsNumber() + vb.AsNumber())
			case _StringString:
				env.A = NewStringValue(va.AsString() + vb.AsString())
			default:
				panicf("can't apply '+' on %+v and %+v", va, vb)
			}
		case OP_SUB:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); testTypes(va, vb) {
			case _NumberNumber:
				env.A.SetNumberValue(env.Get(opa, K).AsNumber() - env.Get(opb, K).AsNumber())
			default:
				panicf("can't apply '-' on %+v and %+v", va, vb)
			}
		case OP_MUL:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); testTypes(va, vb) {
			case _NumberNumber:
				env.A.SetNumberValue(env.Get(opa, K).AsNumber() * env.Get(opb, K).AsNumber())
			default:
				panicf("can't apply '*' on %+v and %+v", va, vb)
			}
		case OP_DIV:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); testTypes(va, vb) {
			case _NumberNumber:
				env.A.SetNumberValue(env.Get(opa, K).AsNumber() / env.Get(opb, K).AsNumber())
			default:
				panicf("can't apply '/' on %+v and %+v", va, vb)
			}
		case OP_MOD:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); testTypes(va, vb) {
			case _NumberNumber:
				env.A.SetNumberValue(float64(int64(env.Get(opa, K).AsNumber()) % int64(env.Get(opb, K).AsNumber())))
			default:
				panicf("can't apply '%%' on %+v and %+v", va, vb)
			}
		case OP_EQ:
			env.A.SetBoolValue(env.Get(opa, K).Equal(env.Get(opb, K)))
		case OP_NEQ:
			env.A.SetBoolValue(!env.Get(opa, K).Equal(env.Get(opb, K)))
		case OP_LESS:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); testTypes(va, vb) {
			case _NumberNumber:
				env.A.SetBoolValue(va.AsNumber() < vb.AsNumber())
			case _StringString:
				env.A.SetBoolValue(va.AsString() < vb.AsString())
			default:
				panicf("can't apply '<' on %+v and %+v", env.Get(opa, K), env.Get(opb, K))
			}
		case OP_LESS_EQ:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); testTypes(va, vb) {
			case _NumberNumber:
				env.A.SetBoolValue(va.AsNumber() <= vb.AsNumber())
			case _StringString:
				env.A.SetBoolValue(va.AsString() <= vb.AsString())
			default:
				panicf("can't apply '<=' on %+v and %+v", env.Get(opa, K), env.Get(opb, K))
			}
		case OP_NOT:
			env.A.SetBoolValue(env.Get(opa, K).IsFalse())
		case OP_BIT_NOT:
			if va := env.Get(opa, K); va.Type() == NumberType {
				env.A.SetNumberValue(float64(^int32(int64(va.AsNumber()) & max32)))
			} else {
				panicf("can't apply 'bit not' on %+v", va)
			}
		case OP_BIT_AND:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); testTypes(va, vb) {
			case _NumberNumber:
				env.A.SetNumberValue(float64(int32(int64(va.AsNumber())&max32) & int32(int64(vb.AsNumber())&max32)))
			default:
				panicf("can't apply '&' on %+v and %+v", env.Get(opa, K), env.Get(opb, K))
			}
		case OP_BIT_OR:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); testTypes(va, vb) {
			case _NumberNumber:
				env.A.SetNumberValue(float64(int32(int64(va.AsNumber())&max32) | int32(int64(vb.AsNumber())&max32)))
			default:
				panicf("can't apply '|' on %+v and %+v", env.Get(opa, K), env.Get(opb, K))
			}
		case OP_BIT_XOR:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); testTypes(va, vb) {
			case _NumberNumber:
				env.A.SetNumberValue(float64(int32(int64(va.AsNumber())&max32) ^ int32(int64(vb.AsNumber())&max32)))
			default:
				panicf("can't apply '^' on %+v and %+v", env.Get(opa, K), env.Get(opb, K))
			}
		case OP_BIT_LSH:
			switch va, vb := env.Get(opa, K), env.Get(opb, K); testTypes(va, vb) {
			case _NumberNumber:
				env.A.SetNumberValue(float64(int32(int64(va.AsNumber())&max32) << uint(vb.AsNumber())))
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
		case OP_BIT_RSH:
			if testTypes(env.Get(opa, K), env.Get(opb, K)) == _NumberNumber {
				env.A.SetNumberValue(float64(int32(int64(env.Get(opa, K).AsNumber())&max32) >> uint(env.Get(opb, K).AsNumber())))
			} else {
				panicf("can't apply '>>' on %+v and %+v", env.Get(opa, K), env.Get(opb, K))
			}
		case OP_BIT_URSH:
			if testTypes(env.Get(opa, K), env.Get(opb, K)) == _NumberNumber {
				env.A.SetNumberValue(float64(uint32(uint64(env.Get(opa, K).AsNumber())&max32) >> uint(env.Get(opb, K).AsNumber())))
			} else {
				panicf("can't apply '>>>' on %+v and %+v", env.Get(opa, K), env.Get(opb, K))
			}
		case OP_ASSERT:
			if a := env.Get(opa, K); a.IsFalse() {
				msg := env.Get(opb, K)
				if msg.Type() == StringType {
					panic(msg.AsString())
				}
				panicf("assertion failed: %+v", a)
			}
			env.A = NewBoolValue(true)
		case OP_LEN:
			switch v := env.Get(opa, K); v.Type() {
			case StringType:
				env.A.SetNumberValue(float64(len(v.AsString())))
			case MapType:
				env.A.SetNumberValue(float64(v.AsMap().Size()))
			default:
				panicf("can't evaluate the length of %+v", v)
			}
		case OP_MAKEMAP:
			if newEnv == nil {
				env.A = NewMapValue(NewMap())
			} else {
				if opa == 1 {
					size := newEnv.SSize()
					m := NewMapSize(size)
					copy(m.l, newEnv.stack)
					newEnv.SClear()
					env.A = NewMapValue(m)
				} else {
					size, m := newEnv.SSize(), NewMap()
					for i := 0; i < size; i += 2 {
						m.Put(newEnv.SGet(i), newEnv.SGet(i+1))
					}
					newEnv.SClear()
					env.A = NewMapValue(m)
				}
			}
		case OP_STORE:
			vidx, v := env.Get(opa, K), env.Get(opb, K)
			switch env.A.Type() {
			case MapType:
				if m := env.A.AsMap(); vidx.Type() == NumberType {
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
				switch testTypes(vidx, v) {
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
			default:
				panicf("can't modify %+v[%+v] to %+v", env.A, vidx, v)
			}
			env.A = v
		case OP_LOAD:
			var v Value
			a := env.Get(opa, K)
			vidx := env.Get(opb, K)
			switch testTypes(a, vidx) {
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
							if cls.argsCount > 0 {
								if len(cls.partialArgs) > 0 {
									panicf("curry function with a receiver")
								}
								cls.argsCount--
								cls.partialArgs = []Value{a}
							}
							v = NewClosureValue(cls)
						}
					}
				} else {
					panicf("can't load %+v[%+v]", a, vidx)
				}
			}
			env.A = v
		case OP_POP:
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
			default:
				env.A = PhantomValue
			}
		case OP_SLICE:
			start, end := int(env.Get(opa, K).Num()), int(env.Get(opb, K).Num())
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
		case OP_PUSH:
			if newEnv == nil {
				newEnv = NewEnv(nil, env.Cancel)
			}
			newEnv.SPush(env.Get(opa, K))
		case OP_RET:
			v := env.Get(opa, K)
			if len(retStack) == 0 {
				return v, 0, false
			}
			returnUpperWorld(v)
		case OP_YIELD:
			return env.Get(opa, K), cursor, true
		case OP_LAMBDA:
			env.A = NewClosureValue(crReadClosure(K.code, &cursor, env, opa, opb))
		case OP_CALL:
			v := env.Get(opa, K)
			if x := v.Type(); x != ClosureType {
				if x == MapType {
					env.A = NewMapValue(v.AsMap().Dup())
					continue
				}
				v.panicType(ClosureType)
			}
			cls := v.AsClosure()
			if cls.lastenv != nil {
				env.A = cls.Exec(nil)
				newEnv = nil
			} else if (newEnv == nil && cls.argsCount > 0) ||
				(newEnv != nil && newEnv.SSize() < int(cls.argsCount)) {
				if newEnv == nil || newEnv.SSize() == 0 {
					env.A = NewClosureValue(cls.Dup())
				} else {
					curry := cls.Dup()
					curry.AppendPreArgs(newEnv.Stack())
					env.A = NewClosureValue(curry)
					newEnv.SClear()
				}
			} else {
				if newEnv == nil {
					newEnv = NewEnv(env, env.Cancel)
				}
				if len(cls.partialArgs) > 0 {
					newEnv.SInsert(0, cls.partialArgs)
				}
				if cls.Isset(ClsYieldable) || cls.native != nil || cls.Isset(ClsRecoverable) {
					newEnv.trace = retStack
					newEnv.parent = env
					env.A = cls.Exec(newEnv)
				} else {
					// log.Println(newEnv.stack)
					last := stacktrace{
						cursor:      cursor,
						env:         env,
						cls:         K,
						noenvescape: cls.Isset(ClsNoEnvescape),
					}

					// switch to the env of cls
					cursor = 0
					K = cls
					caddr = kodeaddr(K.code)
					newEnv.parent = cls.env
					env = newEnv

					retStack = append(retStack, last)
				}
				if cls.native == nil {
					newEnv = nil
				} else {
					newEnv.SClear()
				}
			}

		case OP_JMP:
			cursor = uint32(int32(cursor) + int32(opb) - 1<<12)
		case OP_IFNOT:
			if cond := env.Get(opa, K); cond.IsZero() || cond.IsFalse() {
				cursor = uint32(int32(cursor) + int32(opb) - 1<<12)
			}
		case OP_IF:
			if cond := env.Get(opa, K); !cond.IsZero() || !cond.IsFalse() {
				cursor = uint32(int32(cursor) + int32(opb) - 1<<12)
			}
		case OP_COPY:
			v, r := doCopy(env, env.Get(opa, K).AsNumber(), env.Get(opb, K))
			if r {
				if len(retStack) == 0 {
					return v, 0, false
				}
				returnUpperWorld(v)
			}
		//case OP_CHAR:
		//	switch va := env.Get(opa, K); va.Type() {
		//	case NumberType:
		//		env.A = NewStringValue(string(rune(va.AsNumber())))
		//	case StringType:
		//		env.A = va
		//	default:
		//		panicf("unknown type for char(): %v", va)
		//	}
		case OP_TYPEOF:
			env.A = NewStringValue(TMapping[env.Get(opa, K).Type()])
		}
	}

	if len(retStack) > 0 {
		returnUpperWorld(Value{})
		goto MAIN
	}
	return Value{}, 0, false
}

// OP_COPY takes 3 arguments:
//   1. number:
//		0: do a generic copy but the result will be discarded
//		1: do a generic copy, the result will be stored into somewhere
//		2: the result will be the current stack (copied)
//      all other numbers panic
//   2. any: the subject to be duplicated
//   3. predicator: a closure or nil
func doCopy(env *Env, flag float64, pred Value) (_v Value, _b bool) {
	nopred := pred.Type() == NilType
	alloc := flag == 1

	if flag == 2 {
		ret := NewMapSize(len(env.stack))
		copy(ret.l, env.stack)
		env.A = NewMapValue(ret)
		return
	}

	// immediate value and generic will be returned directly since they can't be truly duplicated
	// however string is an exception
	switch env.A.Type() {
	case NilType, NumberType, PointerType:
		return
	case ClosureType:
		if alloc {
			env.A = NewClosureValue(env.A.AsClosure().Dup())
		} else {
			env.A = Value{}
		}
		return
	case StringType:
		if nopred {
			if alloc {
				s := env.A.AsString()
				m := NewMapSize(len(s))
				for _, x := range s {
					m.l = append(m.l, NewNumberValue(float64(x)))
				}
				env.A = NewMapValue(m)
			}
		} else {
			cls := pred.Cls()
			newEnv := NewEnv(cls.Env(), env.Cancel)
			str := env.A.AsString()
			var newstr []Value
			if alloc {
				newstr = make([]Value, 0, len(str))
			}
			for i, v := range str {
				newEnv.SClear()
				newEnv.SPush(NewNumberValue(float64(i)))
				newEnv.SPush(NewNumberValue(float64(v)))
				if alloc {
					newstr = append(newstr, cls.Exec(newEnv))
				} else {
					if cls.Isset(_ClsPseudoForeach) {
						if res := cls.Exec(newEnv); res == PhantomValue {
							break
						} else if cls.lastp > 0 {
							_b = true
							_v = res
							return
						}
					}
				}
			}
			if alloc {
				m := NewMap()
				m.l = newstr
				env.A = NewMapValue(m)
			}
		}
		return
	}

	// now R1 can only be a map
	if alloc && nopred {
		// simple dup of map
		if env.A.Type() != MapType {
			env.A.panicType(MapType)
		}
		env.A = NewMapValue(env.A.AsMap().Dup())
		return
	}

	if !alloc && nopred {
		// simple copy(a), but its value is discarded
		// so nothing to do here
		return
	}

	// now R2 should be closure
	cls := pred.Cls()
	newEnv := NewEnv(cls.Env(), env.Cancel)
	switch env.A.Type() {
	case MapType:
		if alloc {
			env.A = NewMapValue(env.A.AsMap().Dup())
		} else {
			m := env.A.AsMap()
			for i, v := range m.l {
				newEnv.SClear()
				newEnv.SPush(NewNumberValue(float64(i)))
				newEnv.SPush(v)
				if res := cls.Exec(newEnv); res == PhantomValue {
					goto BREAK_ALL
				} else if cls.lastp > 0 {
					_b = true
					_v = res
					return
				}
			}
			for k, v := range m.m {
				newEnv.SClear()
				newEnv.SPush(NewValueFromInterface(k))
				newEnv.SPush(v)
				if res := cls.Exec(newEnv); res == PhantomValue {
					goto BREAK_ALL
				} else if cls.lastp > 0 {
					_b = true
					_v = res
					return
				}
			}
		BREAK_ALL:
		}
	}

	return
}
