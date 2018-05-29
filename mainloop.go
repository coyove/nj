package potatolang

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"unicode/utf8"
	"unsafe"
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

type ret struct {
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
	stacks []ret
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
	var retStack []ret
	var caddr = kodeaddr(code)
	var kaddr = (*reflect.SliceHeader)(unsafe.Pointer(&consts)).Data

	defer func() {
		if r := recover(); r != nil {
			rr := ret{
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
				e.stacks = make([]ret, len(retStack)+1)
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
		r.env.A, r.env.E = v, env.E
		if r.noenvescape {
			newEnv = env
			newEnv.SClear()
		}
		env = r.env
		retStack = retStack[:len(retStack)-1]
	}
MAIN:
	for {
		// lastCursor = cursor
		// log.Println(cruop(caddr, &cursor))
		// log.Println(op(cruRead64(caddr, &lastCursor)))
		// os.Exit(1)
		// log.Println(cursor)
		bop, opa, opb := cruop(caddr, &cursor)
		switch bop {
		case OP_LINE:
			lineinfo = crReadString(code, &cursor)
		case OP_EOB:
			break MAIN
		case OP_NOP:
		case OP_WHO:
			env.A = env.C
		case OP_NIL:
			env.A = NewValue()
		case OP_TRUE:
			env.A = NewBoolValue(true)
		case OP_FALSE:
			env.A = NewBoolValue(false)
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
				env.A = NewNumberValue(env.R0.AsNumber() + env.R1.AsNumber())
			case _Tstringstring:
				env.A = NewStringValue(env.R0.AsString() + env.R1.AsString())
			case _Tbytesnumber:
				buf := env.R0.AsBytes()
				xbuf := make([]byte, len(buf)+1)
				copy(xbuf, buf)
				xbuf[len(xbuf)-1] = byte(env.R1.AsNumber())
				env.A = NewBytesValue(xbuf)
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
				env.A = NewNumberValue(env.R0.AsNumber() - env.R1.AsNumber())
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
				env.A = NewNumberValue(env.R0.AsNumber() / env.R1.AsNumber())
			} else {
				log.Panicf("can't apply 'div' on %+v and %+v", env.R0, env.R1)
			}
		case OP_MOD:
			if testTypes(env.R0, env.R1) == _Tnumbernumber {
				env.A = NewNumberValue(float64(int64(env.R0.AsNumber()) % int64(env.R1.AsNumber())))
			} else {
				log.Panicf("can't apply 'mod' on %+v and %+v", env.R0, env.R1)
			}
		case OP_EQ:
			env.A = NewBoolValue(env.R0.Equal(env.R1))
		case OP_NEQ:
			env.A = NewBoolValue(!env.R0.Equal(env.R1))
		case OP_LESS:
			switch testTypes(env.R0, env.R1) {
			case _Tnumbernumber:
				env.A = NewBoolValue(env.R0.AsNumber() < env.R1.AsNumber())
			case _Tstringstring:
				env.A = NewBoolValue(env.R0.AsString() < env.R1.AsString())
			default:
				log.Panicf("can't apply 'less' on %+v and %+v", env.R0, env.R1)
			}
		case OP_LESS_EQ:
			switch testTypes(env.R0, env.R1) {
			case _Tnumbernumber:
				env.A = NewBoolValue(env.R0.AsNumber() <= env.R1.AsNumber())
			case _Tstringstring:
				env.A = NewBoolValue(env.R0.AsString() <= env.R1.AsString())
			default:
				log.Panicf("can't apply 'less equal' on %+v and %+v", env.R0, env.R1)
			}
		case OP_NOT:
			if env.R0.ty == Tbool {
				env.A = NewBoolValue(!env.R0.AsBool())
			} else if env.R0.IsFalse() {
				env.A = NewBoolValue(true)
			} else {
				env.A = NewBoolValue(false)
			}
		case OP_BIT_NOT:
			if env.R0.ty == Tnumber {
				env.A = NewNumberValue(float64(^int64(env.R0.AsNumber())))
			} else {
				log.Panicf("can't apply 'bit not' on %+v", env.R0)
			}
		case OP_BIT_AND:
			switch testTypes(env.R0, env.R1) {
			case _Tnumbernumber:
				env.A = NewNumberValue(float64(int64(env.R0.AsNumber()) & int64(env.R1.AsNumber())))
			case _Tbytesbytes:
				buf, buf2 := env.R0.AsBytes(), env.R1.AsBytes()
				xbuf := make([]byte, 0, len(buf)+len(buf2))
				xbuf = append(append(xbuf, buf...), buf2...)
				env.A = NewBytesValue(xbuf)
			case _Tmapmap:
				tr, m := env.R0.AsMap().Dup(nil), env.R1.AsMap()
				for _, v := range m.l {
					tr.l = append(tr.l, v)
				}
				for _, v := range m.m {
					tr.Put(v[0], v[1])
				}
				env.A = NewMapValue(tr)
			default:
				if env.R0.ty == Tstring {
					switch ss := env.R0.AsString(); env.R1.ty {
					case Tnumber:
						num := env.R1.AsNumber()
						if float64(int64(num)) == num {
							env.A = NewStringValue(ss + strconv.FormatInt(int64(num), 10))
						} else {
							env.A = NewStringValue(ss + strconv.FormatFloat(num, 'f', -1, 64))
						}
					case Tbool:
						env.A = NewStringValue(ss + strconv.FormatBool(env.R1.AsBool()))
					case Tstring:
						env.A = NewStringValue(ss + env.R1.AsString())
					case Tbytes:
						env.A = NewStringValue(ss + string(env.R1.AsBytes()))
					default:
						env.A = NewStringValue(ss + env.R1.ToPrintString())
					}
				} else {
					log.Panicf("can't apply 'bit and' on %+v and %+v", env.R0, env.R1)
				}
			}
		case OP_BIT_OR:
			if testTypes(env.R0, env.R1) == _Tnumbernumber {
				env.A = NewNumberValue(float64(int64(env.R0.AsNumber()) | int64(env.R1.AsNumber())))
			} else {
				log.Panicf("can't apply 'bit or' on %+v and %+v", env.R0, env.R1)
			}
		case OP_BIT_XOR:
			if testTypes(env.R0, env.R1) == _Tnumbernumber {
				env.A = NewNumberValue(float64(int64(env.R0.AsNumber()) ^ int64(env.R1.AsNumber())))
			} else {
				log.Panicf("can't apply 'bit xor' on %+v and %+v", env.R0, env.R1)
			}
		case OP_BIT_LSH:
			if testTypes(env.R0, env.R1) == _Tnumbernumber {
				env.A = NewNumberValue(float64(uint64(env.R0.AsNumber()) << uint64(env.R1.AsNumber())))
			} else {
				log.Panicf("can't apply 'bit lsh' on %+v and %+v", env.R0, env.R1)
			}
		case OP_BIT_RSH:
			if testTypes(env.R0, env.R1) == _Tnumbernumber {
				env.A = NewNumberValue(float64(uint64(env.R0.AsNumber()) >> uint64(env.R1.AsNumber())))
			} else {
				log.Panicf("can't apply 'bit rsh' on %+v and %+v", env.R0, env.R1)
			}
		case OP_ASSERT:
			loc := "assertion failed: " + crReadString(code, &cursor)
			if env.R0.IsFalse() {
				panic(loc)
			}
			env.A = NewBoolValue(true)
		case OP_LEN:
			switch v := env.R0; v.Type() {
			case Tstring:
				env.A = NewNumberValue(float64(len(v.AsString())))
			case Tmap:
				env.A = NewNumberValue(float64(v.AsMap().Size()))
			case Tbytes:
				env.A = NewNumberValue(float64(len(v.AsBytes())))
			default:
				log.Panicf("can't evaluate the length of %+v", v)
			}
		case OP_ERROR:
			if env.R0.Type() != Tnil {
				env.E = env.R0
			} else {
				env.A = env.E
				env.E = NewValue()
			}
		case OP_MAP:
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
			switch testTypes(env.R0, env.R1) {
			case _Tbytesnumber:
				if env.R2.ty == Tnumber {
					if b, idx := env.R0.AsBytes(), int(env.R1.AsNumber()); idx >= 0 {
						b[idx] = byte(env.R2.AsNumber())
					} else {
						b[len(b)+idx] = byte(env.R2.AsNumber())
					}
				} else {
					log.Panicf("can't store into %+v with key %+v", env.R0, env.R1)
				}
			case _Tmapnumber:
				m := env.R0.AsMap()
				if idx, ln := int(env.R1.AsNumber()), len(m.l); idx < ln {
					m.l[idx] = env.R2
					break
				} else if idx == ln {
					m.l = append(m.l, env.R2)
					break
				}
				fallthrough
			default:
				if env.R0.ty == Tmap {
					env.R0.AsMap().putIntoMap(env.R1, env.R2)
				} else {
					log.Panicf("can't store into %+v with key %+v", env.R0, env.R1)
				}
			}
			env.A = env.R2
		case OP_LOAD:
			var v Value
			switch testTypes(env.R0, env.R1) {
			case _Tbytesnumber:
				if b, idx := env.R0.AsBytes(), int(env.R1.AsNumber()); idx >= 0 {
					v = NewNumberValue(float64(b[idx]))
				} else {
					v = NewNumberValue(float64(b[len(b)+idx]))
				}
			case _Tstringnumber:
				if b, idx := env.R0.AsString(), int(env.R1.AsNumber()); idx >= 0 {
					v = NewNumberValue(float64(b[idx]))
				} else {
					v = NewNumberValue(float64(b[len(b)+idx]))
				}
			case _Tmapnumber:
				if m, idx := env.R0.AsMap(), int(env.R1.AsNumber()); idx < len(m.l) {
					v = m.l[idx]
					break
				}
				fallthrough
			default:
				if env.R0.ty == Tmap {
					v, _ = env.R0.AsMap().getFromMap(env.R1)
					if v.Type() == Tclosure {
						v.AsClosure().SetCaller(env.R0)
					}
				} else {
					log.Panicf("can't load from %+v with key %+v", env.R0, env.R1)
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
			yieldable := byte(metadata>>16) == 1
			errorable := byte(metadata>>8) == 1
			noenvescape := byte(metadata) == 1
			constsLen := opa
			consts := make([]Value, constsLen)
			for i := uint32(0); i < constsLen; i++ {
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
			env.A = NewClosureValue(NewClosure(buf, consts, env, byte(argsCount), yieldable, errorable, noenvescape))
		case OP_CALL:
			switch v := env.Get(opa); v.ty {
			case Tclosure:
				cls := v.AsClosure()
				if cls.lastenv != nil {
					env.A = cls.Exec(newEnv)
					if newEnv != nil {
						env.E = newEnv.E
					}
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
						newEnv = NewEnv(nil)
					}
					if len(cls.preArgs) > 0 {
						newEnv.SInsert(0, cls.preArgs)
					}

					if cls.yieldable || cls.native != nil {
						env.A = cls.Exec(newEnv)
						env.E = newEnv.E
					} else {
						if retStack == nil {
							retStack = make([]ret, 0, 1)
						}
						//  log.Println(newEnv.stack)
						last := ret{
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
						newEnv.C = cls.caller
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
			case Tbytes:
				if bb := v.AsBytes(); newEnv.SSize() == 2 {
					env.A = NewBytesValue(bb[shiftIndex(newEnv.SGet(0), len(bb)):shiftIndex(newEnv.SGet(1), len(bb))])
				} else if newEnv.SSize() == 1 {
					env.A = NewNumberValue(float64(bb[shiftIndex(newEnv.SGet(0), len(bb))]))
				} else {
					log.Panicf("invalid call on %v", v)
				}
				newEnv.SClear()
			case Tstring:
				if bb := v.AsString(); newEnv.SSize() == 2 {
					env.A = NewStringValue(bb[shiftIndex(newEnv.SGet(0), len(bb)):shiftIndex(newEnv.SGet(1), len(bb))])
				} else if newEnv.SSize() == 1 {
					env.A = NewNumberValue(float64(bb[shiftIndex(newEnv.SGet(0), len(bb))]))
				} else {
					log.Panicf("invalid call on %v", v)
				}
				newEnv.SClear()
			default:
				log.Panicf("invalid callee: %+v", v)
			}
		case OP_JMP:
			off := int32(opb)
			cursor = uint32(int32(cursor) + off)
		case OP_IFNOT:
			cond := env.Get(opa)
			off := int32(opb)
			if cond.ty == Tbool && !cond.AsBool() {
				cursor = uint32(int32(cursor) + off)
			} else if cond.IsFalse() {
				cursor = uint32(int32(cursor) + off)
			}
		case OP_IF:
			cond := env.Get(opa)
			off := int32(opb)
			if cond.ty == Tbool && cond.AsBool() {
				cursor = uint32(int32(cursor) + off)
			} else if !cond.IsFalse() {
				cursor = uint32(int32(cursor) + off)
			}
		case OP_DUP:
			doDup(env)
		case OP_TYPEOF:
			if env.R1.ty == Tnumber {
				if n := byte(env.R1.AsNumber()); n == 255 {
					env.A = NewStringValue(TMapping[env.R0.ty])
				} else {
					env.A = NewBoolValue(env.R0.ty == n)
				}
			} else {
				env.A = NewBoolValue(TMapping[env.R0.ty] == env.R1.AsString())
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
			ret := NewMap()
			ret.l = make([]Value, len(env.stack))
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
	case Tnil, Tnumber, Tbool, Tgeneric:
		env.A = env.R1
		return
	case Tclosure:
		env.A = NewClosureValue(env.R1.AsClosure().Dup())
		return
	case Tstring:
		if nopred {
			env.A = NewBytesValue([]byte(env.R1.AsString()))
		} else {
			if env.R2.ty != Tclosure {
				env.R2.panicType(Tclosure)
			}
			cls := env.R2.AsClosure()
			newEnv := NewEnv(cls.Env())
			str := env.R1.AsString()
			var newstr []byte
			if alloc {
				newstr = make([]byte, 0, len(str))
			}
			for i, v := range str {
				newEnv.SClear()
				newEnv.SPush(NewNumberValue(float64(i)))
				newEnv.SPush(NewNumberValue(float64(v)))
				newEnv.SPush(NewNumberValue(float64(len(newstr))))
				if newEnv.E.Type() != Tnil {
					break
				}
				if alloc {
					v, _, _ := ExecCursor(newEnv, cls.code, cls.consts, 0)
					if v.ty != Tnumber {
						v.panicType(Tnumber)
					}
					r := rune(v.AsNumber())
					idx := len(newstr)
					newstr = append(newstr, 0, 0, 0, 0)
					newstr = newstr[:idx+utf8.EncodeRune(newstr[idx:], r)]
				}
			}
			if alloc {
				env.A = NewBytesValue(newstr)
			}
		}
		return
	}

	if alloc && nopred {
		// simple dup of list, map and bytes
		switch env.R1.Type() {
		case Tmap:
			env.A = NewMapValue(env.R1.AsMap().Dup(nil))
			return
		case Tbytes:
			bytes0 := env.R1.AsBytes()
			bytes1 := make([]byte, len(bytes0))
			copy(bytes1, bytes0)
			env.A = NewBytesValue(bytes1)
			return
		}
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
			if cls.errorable {
				// the predicator may return error and interrupt the dup, so full copy is not used here
				// however cls.errorable is not 100% accurate because calling error() (to check error) and
				// calling error(...) (to throw error) are different behaviors, but i will left this as a TODO
				m2 := NewMap()
				m := env.R1.AsMap()
				for i, v := range m.l {
					idx := NewNumberValue(float64(i))
					newEnv.SClear()
					newEnv.SPush(idx)
					newEnv.SPush(v)
					ret, _, _ := ExecCursor(newEnv, cls.code, cls.consts, 0)
					if newEnv.E.Type() != Tnil {
						break
					}
					m2.Put(idx, ret)
				}
				for _, v := range m.m {
					newEnv.SClear()
					newEnv.SPush(v[0])
					newEnv.SPush(v[1])
					ret, _, _ := ExecCursor(newEnv, cls.code, cls.consts, 0)
					if newEnv.E.Type() != Tnil {
						break
					}
					m2.Put(v[0], ret)
				}
				env.A = NewMapValue(m2)
			} else {
				// full copy
				env.A = NewMapValue(env.R1.AsMap().Dup(func(k Value, v Value) Value {
					newEnv.SClear()
					newEnv.SPush(k)
					newEnv.SPush(v)
					ret, _, _ := ExecCursor(newEnv, cls.code, cls.consts, 0)
					return ret
				}))
			}
		} else {
			m := env.R1.AsMap()
			for i, v := range m.l {
				newEnv.SClear()
				newEnv.SPush(NewNumberValue(float64(i)))
				newEnv.SPush(v)
				ExecCursor(newEnv, cls.code, cls.consts, 0)
				if newEnv.E.Type() != Tnil {
					break
				}
			}
			for _, v := range m.m {
				newEnv.SClear()
				newEnv.SPush(v[0])
				newEnv.SPush(v[1])
				ExecCursor(newEnv, cls.code, cls.consts, 0)
				if newEnv.E.Type() != Tnil {
					break
				}
			}
		}
	case Tbytes:
		var list0 = env.R1.AsBytes()
		var list1 []byte
		if alloc {
			list1 = make([]byte, 0, len(list0))
		}
		for i, v := range list0 {
			newEnv.SClear()
			newEnv.SPush(NewNumberValue(float64(i)))
			newEnv.SPush(NewNumberValue(float64(v)))
			ret, _, _ := ExecCursor(newEnv, cls.code, cls.consts, 0)
			if newEnv.E.Type() != Tnil {
				break
			}
			if alloc {
				if ret.ty != Tnumber {
					ret.panicType(Tnumber)
				}
				list1 = append(list1, byte(ret.AsNumber()))
			}
		}
		if alloc {
			env.A = NewBytesValue(list1)
		}
	}
}
