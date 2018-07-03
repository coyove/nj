package potatolang

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/coyove/potatolang/parser"
)

func init() {
	if strconv.IntSize != 64 {
		panic("potatolang can only run under 64bit")
	}
	one := uint32(1)
	if *(*byte)(unsafe.Pointer(&one)) != 1 {
		panic("potatolang only support little endian arch now")
	}
}

type stacktrace struct {
	cursor      uint32
	noenvescape bool
	env         *Env
	code        []uint64
	kaddr       uintptr
	line        string
}

// ExecError represents the runtime error
type ExecError struct {
	r      interface{}
	stacks []stacktrace
}

func (e *ExecError) Error() string {
	msg := "stacktrace:\n"
	for i := len(e.stacks) - 1; i >= 0; i-- {
		r := e.stacks[i]
		// the recorded cursor was advanced by 1 already
		msg += fmt.Sprintf("cursor: %d at <%08x>, source: %s\n", r.cursor-1, crHash(r.code), r.line)
	}
	return msg
}

func konst(addr uintptr, idx uint16) Value { return *(*Value)(unsafe.Pointer(addr + uintptr(idx)*16)) }

func kodeaddr(code []uint64) uintptr { return (*reflect.SliceHeader)(unsafe.Pointer(&code)).Data }

// ExecCursor executes code under the given env from the given start cursor and returns:
// 1. final result 2. yield cursor 3. is yield or not
func ExecCursor(env *Env, code []uint64, consts []Value, cursor uint32) (Value, uint32, bool) {
	var newEnv *Env
	// var lastCursor uint32
	var lineinfo = "<unknown>"
	var retStack []stacktrace
	var caddr = kodeaddr(code)
	var kaddr = (*reflect.SliceHeader)(unsafe.Pointer(&consts)).Data

	defer func() {
		if r := recover(); r != nil {
			rr := stacktrace{
				cursor: cursor,
				code:   code,
				line:   lineinfo,
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
		code = r.code
		caddr = kodeaddr(code)
		kaddr = r.kaddr
		r.env.A = v
		if r.noenvescape {
			newEnv = env
			newEnv.SClear()
		}
		env = r.env
		retStack = retStack[:len(retStack)-1]
	}
MAIN:
	for {
		// log.Println(cursor)
		bop, opa, opb := cruop(caddr, &cursor)
		switch bop {
		case OP_LINE:
			lineinfo = crReadString(code, &cursor)
		case OP_EOB:
			break MAIN
		case OP_NOP:
		case OP_SET:
			env.Set(opa, env.Get(opb))
		case OP_SETK:
			env.Set(opa, konst(kaddr, uint16(opb)))
		case OP_INC:
			num := env.Get(opa).AsNumber()
			env.Set(opa, NewNumberValue(num+konst(kaddr, uint16(opb)).AsNumber()))
		case OP_ADD:
			switch testTypes(env.R0, env.R1) {
			case _Tnumbernumber:
				env.A.SetNumberValue(env.R0.AsNumber() + env.R1.AsNumber())
			case _Tstringstring:
				env.A = NewStringValue(env.R0.AsString() + env.R1.AsString())
			default:
				if env.R0.ty == Tmap {
					m := env.R0.AsMap().Dup(nil)
					m.l = append(m.l, env.R1)
					env.A = NewMapValue(m)
				} else {
					log.Panicf("can't apply 'add' on %+v and %+v", env.R0, env.R1)
				}
			}
		case OP_SUB:
			if testTypes(env.R0, env.R1) == _Tnumbernumber {
				env.A.SetNumberValue(env.R0.AsNumber() - env.R1.AsNumber())
			} else {
				log.Panicf("can't apply 'sub' on %+v and %+v", env.R0, env.R1)
			}
		case OP_MUL:
			if testTypes(env.R0, env.R1) == _Tnumbernumber {
				env.A = NewNumberValue(env.R0.AsNumber() * env.R1.AsNumber())
			} else {
				log.Panicf("can't apply 'mul' on %+v and %+v", env.R0, env.R1)
			}
		case OP_DIV:
			if testTypes(env.R0, env.R1) == _Tnumbernumber {
				env.A.SetNumberValue(env.R0.AsNumber() / env.R1.AsNumber())
			} else {
				log.Panicf("can't apply 'div' on %+v and %+v", env.R0, env.R1)
			}
		case OP_MOD:
			if testTypes(env.R0, env.R1) == _Tnumbernumber {
				env.A.SetNumberValue(float64(int64(env.R0.AsNumber()) % int64(env.R1.AsNumber())))
			} else {
				log.Panicf("can't apply 'mod' on %+v and %+v", env.R0, env.R1)
			}
		case OP_EQ:
			env.A.SetBoolValue(env.R0.Equal(env.R1))
		case OP_NEQ:
			env.A.SetBoolValue(!env.R0.Equal(env.R1))
		case OP_LESS:
			switch testTypes(env.R0, env.R1) {
			case _Tnumbernumber:
				env.A.SetBoolValue(env.R0.AsNumber() < env.R1.AsNumber())
			case _Tstringstring:
				env.A.SetBoolValue(env.R0.AsString() < env.R1.AsString())
			default:
				log.Panicf("can't apply 'less' on %+v and %+v", env.R0, env.R1)
			}
		case OP_LESS_EQ:
			switch testTypes(env.R0, env.R1) {
			case _Tnumbernumber:
				env.A.SetBoolValue(env.R0.AsNumber() <= env.R1.AsNumber())
			case _Tstringstring:
				env.A.SetBoolValue(env.R0.AsString() <= env.R1.AsString())
			default:
				log.Panicf("can't apply 'less equal' on %+v and %+v", env.R0, env.R1)
			}
		case OP_NOT:
			env.A.SetBoolValue(env.R0.IsFalse())
		case OP_BIT_NOT:
			if env.R0.ty == Tnumber {
				env.A.SetNumberValue(float64(^int32(env.R0.AsNumber())))
			} else {
				log.Panicf("can't apply 'bit not' on %+v", env.R0)
			}
		case OP_BIT_AND:
			switch testTypes(env.R0, env.R1) {
			case _Tnumbernumber:
				env.A.SetNumberValue(float64(int32(env.R0.AsNumber()) & int32(env.R1.AsNumber())))
			case _Tmapmap:
				tr, m := env.R0.AsMap().Dup(nil), env.R1.AsMap()
				for _, v := range m.l {
					tr.l = append(tr.l, v)
				}
				for _, v := range m.m {
					tr.Put(v[0], v[1])
				}
				env.A = NewMapValue(tr)
			case Tnumber<<8 | Tstring:
				num, err := parser.StringToNumber(env.R1.AsString())
				if err != nil {
					env.A = Value{}
				} else {
					env.A = NewNumberValue(env.R0.AsNumber() + num)
				}
			default:
				if env.R0.ty == Tstring {
					switch ss := env.R0.AsString(); env.R1.ty {
					case Tnumber:
						env.A = NewStringValue(ss + strconv.FormatFloat(env.R1.AsNumber(), 'f', -1, 64))
					case Tstring:
						env.A = NewStringValue(ss + env.R1.AsString())
					case Tmap:
						m := env.R1.AsMap()
						buf := make([]byte, len(m.l))
						for i, x := range m.l {
							buf[i] = byte(x.AsNumber())
						}
						env.A = NewStringValue(ss + string(buf))
					default:
						env.A = NewStringValue(ss + env.R1.ToPrintString())
					}
				} else {
					log.Panicf("can't apply 'bit and (concat)' on %+v and %+v", env.R0, env.R1)
				}
			}
		case OP_BIT_OR:
			if testTypes(env.R0, env.R1) == _Tnumbernumber {
				env.A.SetNumberValue(float64(int32(env.R0.AsNumber()) | int32(env.R1.AsNumber())))
			} else {
				log.Panicf("can't apply 'bit or' on %+v and %+v", env.R0, env.R1)
			}
		case OP_BIT_XOR:
			if testTypes(env.R0, env.R1) == _Tnumbernumber {
				env.A.SetNumberValue(float64(int32(env.R0.AsNumber()) ^ int32(env.R1.AsNumber())))
			} else {
				log.Panicf("can't apply 'bit xor' on %+v and %+v", env.R0, env.R1)
			}
		case OP_BIT_LSH:
			if testTypes(env.R0, env.R1) == _Tnumbernumber {
				env.A.SetNumberValue(float64(int32(env.R0.AsNumber()) << uint32(env.R1.AsNumber())))
			} else {
				log.Panicf("can't apply 'bit lsh' on %+v and %+v", env.R0, env.R1)
			}
		case OP_BIT_RSH:
			if testTypes(env.R0, env.R1) == _Tnumbernumber {
				env.A.SetNumberValue(float64(int32(env.R0.AsNumber()) >> uint32(env.R1.AsNumber())))
			} else {
				log.Panicf("can't apply 'bit rsh' on %+v and %+v", env.R0, env.R1)
			}
		case OP_ASSERT:
			loc := crReadString(code, &cursor)
			if env.R0.IsFalse() {
				panic("assertion failed: " + loc)
			}
			env.A = NewBoolValue(true)
		case OP_LEN:
			switch v := env.R3; v.Type() {
			case Tstring:
				env.A.SetNumberValue(float64(len(v.AsString())))
			case Tmap:
				env.A.SetNumberValue(float64(v.AsMap().Size()))
			default:
				log.Panicf("can't evaluate the length of %+v", v)
			}
		case OP_MAKEMAP:
			if newEnv == nil {
				env.A = NewMapValue(NewMap())
			} else {
				size, m := newEnv.SSize(), NewMap()
				for i := 0; i < size; i += 2 {
					m.Put(newEnv.SGet(i), newEnv.SGet(i+1))
				}
				newEnv.SClear()
				env.A = NewMapValue(m)
			}
		case OP_STORE:
			if env.R3.ty == Tmap {
				if m := env.R3.AsMap(); env.R2.ty == Tnumber {
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
			} else {
				log.Panicf("can't store %+v into %+v with key %+v", env.R1, env.R3, env.R2)
			}
			env.A = env.R2
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
				if env.R3.ty == Tmap {
					v, _ = env.R3.AsMap().getFromMap(env.R2)
					if v.Type() == Tclosure {
						v.AsClosure().SetCaller(env.R3)
					}
				} else {
					log.Panicf("can't load from %+v with key %+v", env.R3, env.R2)
				}
			}
			env.A = v
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
		case OP_R0R2:
			env.R0 = env.R2
		case OP_R1R2:
			env.R1 = env.R2
		case OP_POP:
			if x := env.R0; x.ty != Tmap {
				x.panicType(Tmap)
			} else {
				m := x.AsMap()
				l := m.l
				if len(l) == 0 {
					env.A = Value{}
				} else {
					env.A = l[len(l)-1]
					m.l = l[:len(l)-1]
				}
			}
		case OP_PUSH:
			if newEnv == nil {
				newEnv = NewEnv(nil)
			}
			newEnv.SPush(env.Get(opa))
		case OP_PUSHK:
			if newEnv == nil {
				newEnv = NewEnv(nil)
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
			metadata := opb
			argsCount := byte(metadata >> 24)
			yieldable := byte(metadata<<8>>28) == 1
			errorable := byte(metadata<<12>>28) == 1
			noenvescape := byte(metadata<<16>>28) == 1
			receiver := byte(metadata<<20>>28) == 1
			constsLen := opa
			consts := make([]Value, constsLen+1)
			for i := uint32(1); i <= constsLen; i++ {
				switch cruRead64(caddr, &cursor) {
				case Tnumber:
					consts[i] = NewNumberValue(crReadDouble(code, &cursor))
				case Tstring:
					consts[i] = NewStringValue(crReadString(code, &cursor))
				default:
					panic("shouldn't happen")
				}
			}
			buf := crRead(code, &cursor, int(cruRead64(caddr, &cursor)))
			env.A = NewClosureValue(NewClosure(buf, consts, env, byte(argsCount), yieldable, errorable, receiver, noenvescape))
		case OP_CALL:
			v := env.Get(opa)
			if v.ty != Tclosure {
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
					newEnv = NewEnv(env)
				}
				if len(cls.preArgs) > 0 {
					newEnv.SInsert(0, cls.preArgs)
				}
				if cls.receiver {
					newEnv.SPush(cls.caller)
				}
				if cls.Yieldable() || cls.native != nil {
					newEnv.trace = retStack
					newEnv.parent = env
					env.A = cls.Exec(newEnv)
				} else {
					if retStack == nil {
						retStack = make([]stacktrace, 0, 1)
					}
					//  log.Println(newEnv.stack)
					last := stacktrace{
						cursor:      cursor,
						env:         env,
						code:        code,
						kaddr:       kaddr,
						line:        lineinfo,
						noenvescape: cls.noenvescape,
					}

					// switch to the env of cls
					cursor = 0
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
			cursor = uint32(int32(cursor) + int32(opb))
		case OP_IFNOT:
			if cond := env.Get(opa); cond.IsZero() || cond.IsFalse() {
				cursor = uint32(int32(cursor) + int32(opb))
			}
		case OP_IF:
			if cond := env.Get(opa); !cond.IsZero() || !cond.IsFalse() {
				cursor = uint32(int32(cursor) + int32(opb))
			}
		case OP_DUP:
			doDup(env)
		case OP_TYPEOF:
			if env.R1.ty == Tnumber {
				if n := byte(env.R1.AsNumber()); n == 255 {
					env.A = NewStringValue(TMapping[env.R0.ty])
				} else {
					env.A.SetBoolValue(env.R0.ty == n)
				}
			} else {
				env.A.SetBoolValue(TMapping[env.R0.ty] == env.R1.AsString())
			}
		}
	}

	if len(retStack) > 0 {
		returnUpperWorld(NewValue())
		goto MAIN
	}
	return NewValue(), 0, false
}

func shiftIndex(index Value, len int) int {
	if index.ty != Tnumber {
		index.panicType(Tnumber)
	}
	i := int(index.AsNumber())
	if i >= 0 {
		return i
	}
	return i + len
}

// OP_DUP takes 3 arguments:
//   1. number: 0 means the dup result will be discarded, 1 means the result will be stored into somewhere
//   2. any: the subject to be duplicated
//   3. number/closure: 0 means no predicator, 1 means dup stack, 2 means return stack, otherwise the closure will be used
func doDup(env *Env) {
	alloc := env.R0.AsNumber() == 1
	nopred := false
	if env.R2.ty == Tnumber {
		switch env.R2.AsNumber() {
		case 0:
			// dup(a)
			nopred = true
		case 1:
			// dup()
			ret := NewMapSize(len(env.stack))
			copy(ret.l, env.stack)
			env.A = NewMapValue(ret)
			return
		case 2:
			ret := NewMap()
			ret.l = env.stack
			env.A = NewMapValue(ret)
			return
		default:
			panic("serious error")
		}
	}

	// immediate value and generic will be returned directly since they can't be truly duplicated
	// however string is an exception
	switch env.R1.Type() {
	case Tnil, Tnumber, Tgeneric:
		env.A = env.R1
		return
	case Tclosure:
		env.A = NewClosureValue(env.R1.AsClosure().Dup())
		return
	case Tstring:
		if nopred {
			s := env.R1.AsString()
			m := NewMapSize(len(s))
			for _, x := range s {
				m.l = append(m.l, NewNumberValue(float64(x)))
			}
			env.A = NewMapValue(m)
		} else {
			if env.R2.ty != Tclosure {
				env.R2.panicType(Tclosure)
			}
			cls := env.R2.AsClosure()
			newEnv := NewEnv(cls.Env())
			str := env.R1.AsString()
			var newstr []Value
			if alloc {
				newstr = make([]Value, 0, len(str))
			}
			for i, v := range str {
				newEnv.SClear()
				newEnv.SPush(NewNumberValue(float64(i)))
				newEnv.SPush(NewNumberValue(float64(v)))
				newEnv.SPush(NewNumberValue(float64(len(newstr))))
				if alloc {
					v, _, _ := ExecCursor(newEnv, cls.code, cls.consts, 0)
					newstr = append(newstr, v)
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

	if alloc && nopred {
		// simple dup of map
		if env.R1.ty != Tmap {
			env.R1.panicType(Tmap)
		}
		env.A = NewMapValue(env.R1.AsMap().Dup(nil))
		return
	}

	if !alloc && nopred {
		// simple dup(a), but its value is discarded
		// so nothing to do here
		return
	}

	// now R2 should be closure
	if env.R2.ty != Tclosure {
		env.R2.panicType(Tclosure)
	}
	cls := env.R2.AsClosure()
	newEnv := NewEnv(cls.Env())
	switch env.R1.Type() {
	case Tmap:
		if alloc {
			env.A = NewMapValue(env.R1.AsMap().Dup(func(k Value, v Value) Value {
				newEnv.SClear()
				newEnv.SPush(k)
				newEnv.SPush(v)
				ret, _, _ := ExecCursor(newEnv, cls.code, cls.consts, 0)
				return ret
			}))
		} else {
			m := env.R1.AsMap()
			for i := len(m.l) - 1; i >= 0; i-- {
				newEnv.SClear()
				newEnv.SPush(NewNumberValue(float64(i)))
				newEnv.SPush(m.l[i])
				ExecCursor(newEnv, cls.code, cls.consts, 0)
			}
			for _, v := range m.m {
				newEnv.SClear()
				newEnv.SPush(v[0])
				newEnv.SPush(v[1])
				ExecCursor(newEnv, cls.code, cls.consts, 0)
			}
		}
	}
}