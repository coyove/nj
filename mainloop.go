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
		r.env.A = v
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
			env.A.SetNumberValue(env.Get(opa).AsNumber() + env.Get(opb).AsNumber())
		case OP_SUB:
			env.A.SetNumberValue(env.Get(opa).AsNumber() - env.Get(opb).AsNumber())
		case OP_MUL:
			env.A.SetNumberValue(env.Get(opa).AsNumber() * env.Get(opb).AsNumber())
		case OP_DIV:
			env.A.SetNumberValue(env.Get(opa).AsNumber() / env.Get(opb).AsNumber())
		case OP_MOD:
			env.A.SetNumberValue(float64(int64(env.Get(opa).AsNumber()) % int64(env.Get(opb).AsNumber())))
		case OP_EQ:
			env.A.SetBoolValue(env.Get(opa).Equal(env.Get(opb)))
		case OP_NEQ:
			env.A.SetBoolValue(!env.Get(opa).Equal(env.Get(opb)))
		case OP_LESS:
			switch testTypes(env.Get(opa), env.Get(opb)) {
			case _Tnumbernumber:
				env.A.SetBoolValue(env.Get(opa).AsNumber() < env.Get(opb).AsNumber())
			case _Tstringstring:
				env.A.SetBoolValue(env.Get(opa).AsString() < env.Get(opb).AsString())
			default:
				panicf("can't apply 'less' on %+v and %+v", env.Get(opa), env.Get(opb))
			}
		case OP_LESS_EQ:
			switch testTypes(env.Get(opa), env.Get(opb)) {
			case _Tnumbernumber:
				env.A.SetBoolValue(env.Get(opa).AsNumber() <= env.Get(opb).AsNumber())
			case _Tstringstring:
				env.A.SetBoolValue(env.Get(opa).AsString() <= env.Get(opb).AsString())
			default:
				panicf("can't apply 'less equal' on %+v and %+v", env.Get(opa), env.Get(opb))
			}
		case OP_NOT:
			env.A.SetBoolValue(env.Get(opa).IsFalse())
		case OP_BIT_NOT:
			if env.Get(opa).Type() == Tnumber {
				env.A.SetNumberValue(float64(^int32(env.Get(opa).AsNumber())))
			} else {
				panicf("can't apply 'bit not' on %+v", env.Get(opa))
			}
		case OP_BIT_AND:
			switch testTypes(env.Get(opa), env.Get(opb)) {
			case _Tnumbernumber:
				env.A.SetNumberValue(float64(int32(env.Get(opa).AsNumber()) & int32(env.Get(opb).AsNumber())))
			case _Tmapmap:
				tr, m := env.Get(opa).AsMap(), env.Get(opb).AsMap()
				tr.l = append(tr.l, m.l...)
				for _, v := range m.m {
					tr.Put(v[0], v[1])
				}
				env.A = NewMapValue(tr)
			case Tnumber<<8 | Tstring:
				num, err := parser.StringToNumber(env.Get(opb).AsString())
				if err != nil {
					env.A = Value{}
				} else {
					env.A = NewNumberValue(env.Get(opa).AsNumber() + num)
				}
			default:
				if env.Get(opa).Type() == Tstring {
					switch ss := env.Get(opa).AsString(); env.Get(opb).Type() {
					case Tnumber:
						env.A = NewStringValue(ss + strconv.FormatFloat(env.Get(opb).AsNumber(), 'f', -1, 64))
					case Tstring:
						env.A = NewStringValue(ss + env.Get(opb).AsString())
					case Tmap:
						m := env.Get(opb).AsMap()
						buf := make([]byte, len(m.l))
						for i, x := range m.l {
							buf[i] = byte(x.AsNumber())
						}
						env.A = NewStringValue(ss + string(buf))
					default:
						env.A = NewStringValue(ss + env.Get(opb).ToPrintString())
					}
				} else {
					panicf("can't apply 'bit and (concat)' on %+v and %+v", env.Get(opa), env.Get(opb))
				}
			}
		case OP_BIT_OR:
			if testTypes(env.Get(opa), env.Get(opb)) == _Tnumbernumber {
				env.A.SetNumberValue(float64(int32(env.Get(opa).AsNumber()) | int32(env.Get(opb).AsNumber())))
			} else {
				panicf("can't apply 'bit or' on %+v and %+v", env.Get(opa), env.Get(opb))
			}
		case OP_BIT_XOR:
			if testTypes(env.Get(opa), env.Get(opb)) == _Tnumbernumber {
				env.A.SetNumberValue(float64(int32(env.Get(opa).AsNumber()) ^ int32(env.Get(opb).AsNumber())))
			} else {
				panicf("can't apply 'bit xor' on %+v and %+v", env.Get(opa), env.Get(opb))
			}
		case OP_BIT_LSH:
			if testTypes(env.Get(opa), env.Get(opb)) == _Tnumbernumber {
				env.A.SetNumberValue(float64(int32(env.Get(opa).AsNumber()) << uint32(env.Get(opb).AsNumber())))
			} else {
				panicf("can't apply 'bit lsh' on %+v and %+v", env.Get(opa), env.Get(opb))
			}
		case OP_BIT_RSH:
			if testTypes(env.Get(opa), env.Get(opb)) == _Tnumbernumber {
				env.A.SetNumberValue(float64(int32(env.Get(opa).AsNumber()) >> uint32(env.Get(opb).AsNumber())))
			} else {
				panicf("can't apply 'bit rsh' on %+v and %+v", env.Get(opa), env.Get(opb))
			}
		case OP_BIT_URSH:
			if testTypes(env.Get(opa), env.Get(opb)) == _Tnumbernumber {
				env.A.SetNumberValue(float64(uint32(env.Get(opa).AsNumber()) >> uint32(env.Get(opb).AsNumber())))
			} else {
				panicf("can't apply 'bit unsigned rsh' on %+v and %+v", env.Get(opa), env.Get(opb))
			}
		case OP_ASSERT:
			if env.Get(opa).IsFalse() {
				panic(env.Get(opb))
			}
			env.A = NewBoolValue(true)
		case OP_LEN:
			switch v := env.A; v.Type() {
			case Tstring:
				env.A.SetNumberValue(float64(len(v.AsString())))
			case Tmap:
				env.A.SetNumberValue(float64(v.AsMap().Size()))
			case Tgeneric:
				env.A.SetNumberValue(float64(GLen(v)))
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
			if env.A.Type() == Tmap {
				vidx := env.Get(opa)
				if m := env.A.AsMap(); vidx.Type() == Tnumber {
					if idx, ln := int(vidx.AsNumber()), len(m.l); idx < ln {
						m.l[idx] = env.Get(opb)
					} else if idx == ln {
						m.l = append(m.l, env.Get(opb))
					} else {
						m.putIntoMap(vidx, env.Get(opb))
					}
				} else {
					m.putIntoMap(vidx, env.Get(opb))
				}
			} else {
				panicf("can't store %+v into %+v with key %+v", env.Get(opb), env.A, env.Get(opa))
			}
			env.A = env.Get(opb)
		case OP_LOAD:
			var v Value
			vidx := env.Get(opa)
			switch testTypes(env.A, vidx) {
			case _Tstringnumber:
				v.SetNumberValue(float64(env.A.AsString()[int(vidx.AsNumber())]))
			case _Tmapnumber:
				if m, idx := env.A.AsMap(), int(vidx.AsNumber()); idx < len(m.l) {
					v = m.l[idx]
					break
				}
				fallthrough
			default:
				if env.A.Type() == Tmap {
					v, _ = env.A.AsMap().getFromMap(vidx)
					if v.Type() == Tclosure {
						v.AsClosure().SetCaller(env.A)
					}
				} else {
					panicf("can't load from %+v with key %+v", env.A, vidx)
				}
			}
			env.A = v
		case OP_POP:
			switch env.A.Type() {
			case Tmap:
				m := env.A.AsMap()
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
			start, end := int(env.Get(opa).Num()), int(env.Get(opb).Num())
			switch x := env.A; x.Type() {
			case Tstring:
				if end == -1 {
					env.A = NewStringValue(x.AsString()[start:])
				} else {
					env.A = NewStringValue(x.AsString()[start:end])
				}
			case Tmap:
				m := NewMap()
				if end == -1 {
					m.l = x.AsMap().l[start:]
				} else {
					m.l = x.AsMap().l[start:end]
				}
				env.A = NewMapValue(m)
			case Tgeneric:
				if end == -1 {
					env.A = GSlice(x, start, GLen(x))
				} else {
					env.A = GSlice(x, start, end)
				}
			default:
				panicf("can't slice %+v", x)
			}
		case OP_PUSH:
			if newEnv == nil {
				newEnv = NewEnv(nil, flag)
			}
			newEnv.SPush(env.Get(opa))
		case OP_RET:
			v := env.Get(opa)
			if len(retStack) == 0 {
				return v, 0, false
			}
			returnUpperWorld(v)
		case OP_YIELD:
			return env.Get(opa), cursor, true
		case OP_LAMBDA:
			env.A = NewClosureValue(crReadClosure(code, &cursor, env, opa, opb))
		case OP_CALL:
			v := env.Get(opa)
			if v.Type() != Tclosure {
				v.panicType(Tclosure)
			}
			cls := v.AsClosure()
			if cls.lastenv != nil {
				env.A = cls.Exec(nil)
				newEnv = nil
			} else if (newEnv == nil && cls.argsCount > 0) ||
				(newEnv != nil && newEnv.SSize() < int(cls.argsCount)) {
				if newEnv == nil || newEnv.SSize() == 0 {
					env.A = NewClosureValue(cls)
				} else {
					curry := cls.Dup()
					curry.AppendPreArgs(newEnv.Stack())
					env.A = NewClosureValue(curry)
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
					env.A = cls.Exec(newEnv)
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
			v, r := doCopy(env, env.Get(opa).AsNumber(), env.Get(opb))
			if r {
				if len(retStack) == 0 {
					return v, 0, false
				}
				returnUpperWorld(v)
			}
		case OP_TYPEOF:
			env.A = NewStringValue(TMapping[env.Get(opa).Type()])
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
	nopred := pred.Type() == Tnil
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
	case Tnil, Tnumber, Tgeneric:
		return
	case Tclosure:
		if alloc {
			env.A = NewClosureValue(env.A.AsClosure().Dup())
		} else {
			env.A = Value{}
		}
		return
	case Tstring:
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
		if env.A.Type() != Tmap {
			env.A.panicType(Tmap)
		}
		env.A = NewMapValue(env.A.AsMap().Dup(nil))
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
	case Tmap:
		if alloc {
			env.A = NewMapValue(env.A.AsMap().Dup(func(k Value, v Value) Value {
				newEnv.SClear()
				newEnv.SPush(k)
				newEnv.SPush(v)
				return cls.Exec(newEnv)
			}))
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
