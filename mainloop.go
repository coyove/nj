package potatolang

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"sync/atomic"
	"unsafe"

	"github.com/coyove/potatolang/parser"
)

func panicf(msg string, args ...interface{}) {
	panic(fmt.Sprintf(msg, args...))
}

func panicerr(err error) {
	if err != nil {
		panic(err)
	}
}

type stacktrace struct {
	cursor         uint32
	noenvescape    bool
	callReturnInto byte
	env            *Env
	cls            *Closure
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
			var opx uint32 = 0xffffffff
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

func konst(addr uintptr, idx uint16) Value {
	return *(*Value)(unsafe.Pointer(addr + uintptr(idx)*SizeofValue))
}

func kodeaddr(code []uint32) uintptr { return (*reflect.SliceHeader)(unsafe.Pointer(&code)).Data }

// ExecCursor executes code under the given env from the given start cursor and returns:
// 1. final result 2. yield cursor 3. is yield or not
func ExecCursor(env *Env, K *Closure, cursor uint32) (_v Value, _p uint32, _y bool) {
	var newEnv *Env
	var flag *uintptr
	var retStack []stacktrace
	var code = K.code
	var caddr = kodeaddr(code)
	var kaddr = (*reflect.SliceHeader)(unsafe.Pointer(&K.consts)).Data

	defer func() {
		if r := recover(); r != nil {
			if K.Isset(CLS_RECOVERALL) {
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
		code = r.cls.code
		caddr = kodeaddr(code)
		kaddr = (*reflect.SliceHeader)(unsafe.Pointer(&r.cls.consts)).Data
		*r.env.reg(uint16(r.callReturnInto)) = v
		if r.noenvescape {
			newEnv = env
			newEnv.SClear()
		}
		env = r.env
		retStack = retStack[:len(retStack)-1]
	}

	flag = env.Cancel
MAIN:
	for {
		if flag != nil && atomic.LoadUintptr(flag) == 1 {
			panicf("canceled")
		}

		//log.Println(cursor)
		v := *(*uint32)(unsafe.Pointer(uintptr(cursor)*4 + caddr))
		cursor++
		bop, opa, opb := op(v)

		switch bop {
		case OP_EOB:
			break MAIN
		case OP_NOP:
		case OP_SET:
			env.Set(opa, env.Get(opb))
		case OP_SETK:
			env.Set(opa, konst(kaddr, uint16(opb)))
		case OP_INC:
			num := env.Get(opa).AsNumber()
			env.A.SetNumberValue(num + konst(kaddr, uint16(opb)).AsNumber())
			env.Set(opa, env.A)
		case OP_ADD:
			env.reg(opa).SetNumberValue(env.R0.AsNumber() + env.R1.AsNumber())
		case OP_SUB:
			env.reg(opa).SetNumberValue(env.R0.AsNumber() - env.R1.AsNumber())
		case OP_MUL:
			env.reg(opa).SetNumberValue(env.R0.AsNumber() * env.R1.AsNumber())
		case OP_DIV:
			env.reg(opa).SetNumberValue(env.R0.AsNumber() / env.R1.AsNumber())
		case OP_MOD:
			if testTypes(env.R0, env.R1) == _Tnumbernumber {
				env.reg(opa).SetNumberValue(float64(int64(env.R0.AsNumber()) % int64(env.R1.AsNumber())))
			} else {
				panicf("can't apply 'mod' on %+v and %+v", env.R0, env.R1)
			}
		case OP_EQ:
			env.reg(opa).SetBoolValue(env.R0.Equal(env.R1))
		case OP_NEQ:
			env.reg(opa).SetBoolValue(!env.R0.Equal(env.R1))
		case OP_LESS:
			switch testTypes(env.R0, env.R1) {
			case _Tnumbernumber:
				env.reg(opa).SetBoolValue(env.R0.AsNumber() < env.R1.AsNumber())
			case _Tstringstring:
				env.reg(opa).SetBoolValue(env.R0.AsString() < env.R1.AsString())
			default:
				panicf("can't apply 'less' on %+v and %+v", env.R0, env.R1)
			}
		case OP_LESS_EQ:
			switch testTypes(env.R0, env.R1) {
			case _Tnumbernumber:
				env.reg(opa).SetBoolValue(env.R0.AsNumber() <= env.R1.AsNumber())
			case _Tstringstring:
				env.reg(opa).SetBoolValue(env.R0.AsString() <= env.R1.AsString())
			default:
				panicf("can't apply 'less equal' on %+v and %+v", env.R0, env.R1)
			}
		case OP_NOT:
			env.reg(opa).SetBoolValue(env.R0.IsFalse())
		case OP_BIT_NOT:
			if env.R0.Type() == Tnumber {
				env.reg(opa).SetNumberValue(float64(^int32(env.R0.AsNumber())))
			} else {
				panicf("can't apply 'bit not' on %+v", env.R0)
			}
		case OP_BIT_AND:
			switch testTypes(env.R0, env.R1) {
			case _Tnumbernumber:
				env.reg(opa).SetNumberValue(float64(int32(env.R0.AsNumber()) & int32(env.R1.AsNumber())))
			case _Tmapmap:
				tr, m := env.R0.AsMap(), env.R1.AsMap()
				tr.l = append(tr.l, m.l...)
				for _, v := range m.m {
					tr.Put(v[0], v[1])
				}
				*env.reg(opa) = NewMapValue(tr)
			case Tnumber<<8 | Tstring:
				num, err := parser.StringToNumber(env.R1.AsString())
				if err != nil {
					*env.reg(opa) = Value{}
				} else {
					*env.reg(opa) = NewNumberValue(env.R0.AsNumber() + num)
				}
			default:
				if env.R0.Type() == Tstring {
					switch ss := env.R0.AsString(); env.R1.Type() {
					case Tnumber:
						*env.reg(opa) = NewStringValue(ss + strconv.FormatFloat(env.R1.AsNumber(), 'f', -1, 64))
					case Tstring:
						*env.reg(opa) = NewStringValue(ss + env.R1.AsString())
					case Tmap:
						m := env.R1.AsMap()
						buf := make([]byte, len(m.l))
						for i, x := range m.l {
							buf[i] = byte(x.AsNumber())
						}
						*env.reg(opa) = NewStringValue(ss + string(buf))
					default:
						*env.reg(opa) = NewStringValue(ss + env.R1.ToPrintString())
					}
				} else {
					panicf("can't apply 'bit and (concat)' on %+v and %+v", env.R0, env.R1)
				}
			}
		case OP_BIT_OR:
			if testTypes(env.R0, env.R1) == _Tnumbernumber {
				env.reg(opa).SetNumberValue(float64(int32(env.R0.AsNumber()) | int32(env.R1.AsNumber())))
			} else {
				panicf("can't apply 'bit or' on %+v and %+v", env.R0, env.R1)
			}
		case OP_BIT_XOR:
			if testTypes(env.R0, env.R1) == _Tnumbernumber {
				env.reg(opa).SetNumberValue(float64(int32(env.R0.AsNumber()) ^ int32(env.R1.AsNumber())))
			} else {
				panicf("can't apply 'bit xor' on %+v and %+v", env.R0, env.R1)
			}
		case OP_BIT_LSH:
			if testTypes(env.R0, env.R1) == _Tnumbernumber {
				env.reg(opa).SetNumberValue(float64(int32(env.R0.AsNumber()) << uint32(env.R1.AsNumber())))
			} else {
				panicf("can't apply 'bit lsh' on %+v and %+v", env.R0, env.R1)
			}
		case OP_BIT_RSH:
			if testTypes(env.R0, env.R1) == _Tnumbernumber {
				env.reg(opa).SetNumberValue(float64(int32(env.R0.AsNumber()) >> uint32(env.R1.AsNumber())))
			} else {
				panicf("can't apply 'bit rsh' on %+v and %+v", env.R0, env.R1)
			}
		case OP_BIT_URSH:
			if testTypes(env.R0, env.R1) == _Tnumbernumber {
				env.reg(opa).SetNumberValue(float64(uint32(env.R0.AsNumber()) >> uint32(env.R1.AsNumber())))
			} else {
				panicf("can't apply 'bit unsigned rsh' on %+v and %+v", env.R0, env.R1)
			}
		case OP_ASSERT:
			if env.R0.IsFalse() {
				panic(env.R1)
			}
			*env.reg(opa) = NewBoolValue(true)
		case OP_LEN:
			switch v := env.R3; v.Type() {
			case Tstring:
				env.reg(opa).SetNumberValue(float64(len(v.AsString())))
			case Tmap:
				env.reg(opa).SetNumberValue(float64(v.AsMap().Size()))
			case Tgeneric:
				env.reg(opa).SetNumberValue(float64(GLen(v)))
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
			if env.R3.Type() == Tmap {
				if m := env.R3.AsMap(); env.R2.Type() == Tnumber {
					if idx, ln := int(env.R2.AsNumber()), len(m.l); idx < ln {
						m.l[idx] = env.R1
					} else if idx == ln {
						m.l = append(m.l, env.R1)
					} else {
						m.putIntoMap(env.R2, env.R1)
					}
				} else {
					m.putIntoMap(env.R2, env.R1)
				}
			} else if env.R3.Type() == Tgeneric {
				GStore(env.R3, int(env.R2.Num()), env.R1)
			} else {
				panicf("can't store %+v into %+v with key %+v", env.R1, env.R3, env.R2)
			}
			*env.reg(opa) = env.R1
		case OP_LOAD:
			var v Value
			switch testTypes(env.R3, env.R2) {
			case _Tstringnumber:
				v.SetNumberValue(float64(env.R3.AsString()[int(env.R2.AsNumber())]))
			case _Tmapnumber:
				if m, idx := env.R3.AsMap(), int(env.R2.AsNumber()); idx < len(m.l) {
					v = m.l[idx]
					break
				}
				fallthrough
			default:
				if env.R3.Type() == Tmap {
					v, _ = env.R3.AsMap().getFromMap(env.R2)
					if v.Type() == Tclosure {
						v.AsClosure().SetCaller(env.R3)
					}
				} else if env.R3.Type() == Tgeneric {
					v = GLoad(env.R3, int(env.R2.Num()))
				} else {
					panicf("can't load from %+v with key %+v", env.R3, env.R2)
				}
			}
			*env.reg(opa) = v
		case OP_R0:
			env.R0 = env.Get(opa)
		case OP_R0K:
			env.R0 = konst(kaddr, uint16(opa))
		case OP_R1:
			env.R1 = env.Get(opa)
		case OP_R1K:
			env.R1 = konst(kaddr, uint16(opa))
		case OP_R2:
			env.R2 = env.Get(opa)
		case OP_R2K:
			env.R2 = konst(kaddr, uint16(opa))
		case OP_R3:
			env.R3 = env.Get(opa)
		case OP_R3K:
			env.R3 = konst(kaddr, uint16(opa))
		case OP_RX:
			*env.reg(opa + 1) = *env.reg(opb + 1)
		case OP_POP:
			switch env.R3.Type() {
			case Tmap:
				m := env.R3.AsMap()
				l := m.l
				if len(l) == 0 {
					*env.reg(opa) = Value{}
				} else {
					*env.reg(opa) = l[len(l)-1]
					m.l = l[:len(l)-1]
				}
			default:
				*env.reg(opa) = PhantomValue
			}
		case OP_SLICE:
			start, end := int(env.R2.Num()), int(env.R1.Num())
			switch x := env.R3; x.Type() {
			case Tstring:
				if end == -1 {
					*env.reg(opa) = NewStringValue(x.AsString()[start:])
				} else {
					*env.reg(opa) = NewStringValue(x.AsString()[start:end])
				}
			case Tmap:
				m := NewMap()
				if end == -1 {
					m.l = x.AsMap().l[start:]
				} else {
					m.l = x.AsMap().l[start:end]
				}
				*env.reg(opa) = NewMapValue(m)
			case Tgeneric:
				if end == -1 {
					*env.reg(opa) = GSlice(x, start, GLen(x))
				} else {
					*env.reg(opa) = GSlice(x, start, end)
				}
			default:
				panicf("can't slice %+v", x)
			}
		case OP_PUSH:
			if newEnv == nil {
				newEnv = NewEnv(nil, flag)
			}
			newEnv.SPush(env.Get(opa))
		case OP_PUSHK:
			if newEnv == nil {
				newEnv = NewEnv(nil, flag)
			}
			newEnv.SPush(konst(kaddr, uint16(opa)))
		case OP_RET:
			v := env.Get(opa)
			if len(retStack) == 0 {
				return v, 0, false
			}
			returnUpperWorld(v)
		case OP_RETK:
			v := konst(kaddr, uint16(opa))
			if len(retStack) == 0 {
				return v, 0, false
			}
			returnUpperWorld(v)
		case OP_YIELD:
			return env.Get(opa), cursor, true
		case OP_YIELDK:
			return konst(kaddr, uint16(opa)), cursor, true
		case OP_LAMBDA:
			env.A = NewClosureValue(crReadClosure(code, &cursor, env, opa, opb))
		case OP_CALL:
			v := env.Get(opa)
			if v.Type() != Tclosure {
				v.panicType(Tclosure)
			}
			cls := v.AsClosure()
			if cls.lastenv != nil {
				*env.reg(opb) = cls.Exec(nil)
				newEnv = nil
			} else if (newEnv == nil && cls.argsCount > 0) ||
				(newEnv != nil && newEnv.SSize() < int(cls.argsCount)) {
				if newEnv == nil || newEnv.SSize() == 0 {
					*env.reg(opb) = NewClosureValue(cls)
				} else {
					curry := cls.Dup()
					curry.AppendPreArgs(newEnv.Stack())
					*env.reg(opb) = NewClosureValue(curry)
					newEnv.SClear()
				}
			} else {
				if newEnv == nil {
					newEnv = NewEnv(env, flag)
				}
				if len(cls.preArgs) > 0 {
					newEnv.SInsert(0, cls.preArgs)
				}
				if cls.Isset(CLS_HASRECEIVER) {
					newEnv.SPush(cls.caller)
				}
				if cls.Isset(CLS_YIELDABLE) || cls.native != nil || cls.Isset(CLS_RECOVERALL) {
					newEnv.trace = retStack
					newEnv.parent = env
					*env.reg(opb) = cls.Exec(newEnv)
				} else {
					if retStack == nil {
						retStack = make([]stacktrace, 0, 1)
					}
					// log.Println(newEnv.stack)
					last := stacktrace{
						cursor:         cursor,
						env:            env,
						cls:            K,
						noenvescape:    cls.Isset(CLS_NOENVESCAPE),
						callReturnInto: byte(opb),
					}

					// switch to the env of cls
					cursor = 0
					K = cls
					newEnv.parent = cls.env
					env = newEnv
					code = cls.code
					caddr = kodeaddr(code)
					kaddr = (*reflect.SliceHeader)(unsafe.Pointer(&cls.consts)).Data

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
			if cond := env.Get(opa); cond.IsZero() || cond.IsFalse() {
				cursor = uint32(int32(cursor) + int32(opb) - 1<<12)
			}
		case OP_IF:
			if cond := env.Get(opa); !cond.IsZero() || !cond.IsFalse() {
				cursor = uint32(int32(cursor) + int32(opb) - 1<<12)
			}
		case OP_COPY:
			v, r := doCopy(env)
			if r {
				if len(retStack) == 0 {
					return v, 0, false
				}
				returnUpperWorld(v)
			}
		case OP_TYPEOF:
			*env.reg(opa) = NewStringValue(TMapping[env.R0.Type()])
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
func doCopy(env *Env) (_v Value, _b bool) {
	flag := env.R0.AsNumber()
	nopred := env.R2.Type() == Tnil
	alloc := flag == 1

	if flag == 2 {
		ret := NewMapSize(len(env.stack))
		copy(ret.l, env.stack)
		env.A = NewMapValue(ret)
		return
	}

	// immediate value and generic will be returned directly since they can't be truly duplicated
	// however string is an exception
	switch env.R1.Type() {
	case Tnil, Tnumber, Tgeneric:
		env.A = env.R1
		return
	case Tclosure:
		if alloc {
			env.A = NewClosureValue(env.R1.AsClosure().Dup())
		} else {
			env.A = Value{}
		}
		return
	case Tstring:
		if nopred {
			if alloc {
				s := env.R1.AsString()
				m := NewMapSize(len(s))
				for _, x := range s {
					m.l = append(m.l, NewNumberValue(float64(x)))
				}
				env.A = NewMapValue(m)
			}
		} else {
			cls := env.R2.Cls()
			newEnv := NewEnv(cls.Env(), env.Cancel)
			str := env.R1.AsString()
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
					if cls.Isset(CLS_PSEUDO_FOREACH) {
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
		if env.R1.Type() != Tmap {
			env.R1.panicType(Tmap)
		}
		env.A = NewMapValue(env.R1.AsMap().Dup(nil))
		return
	}

	if !alloc && nopred {
		// simple copy(a), but its value is discarded
		// so nothing to do here
		return
	}

	// now R2 should be closure
	cls := env.R2.Cls()
	newEnv := NewEnv(cls.Env(), env.Cancel)
	switch env.R1.Type() {
	case Tmap:
		if alloc {
			env.A = NewMapValue(env.R1.AsMap().Dup(func(k Value, v Value) Value {
				newEnv.SClear()
				newEnv.SPush(k)
				newEnv.SPush(v)
				return cls.Exec(newEnv)
			}))
		} else {
			m := env.R1.AsMap()
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
			for _, v := range m.m {
				newEnv.SClear()
				newEnv.SPush(v[0])
				newEnv.SPush(v[1])
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
