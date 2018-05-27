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
		panic("potatolang only run under 64bit")
	}
}

type ret struct {
	cursor      uint32
	noenvescape bool
	env         *Env
	code        []uint16
	kaddr       uintptr
	line        string
}

// ExecError represents the runtime error
type ExecError struct {
	r      interface{}
	stacks []ret
}

func (e *ExecError) Error() string {
	msg := ""
	for i := len(e.stacks) - 1; i >= 0; i-- {
		r := e.stacks[i]
		msg += fmt.Sprintf("cursor: %d at <%x>, source: %s\n", r.cursor, crHash(r.code), r.line)
	}
	return msg
}

func konst(addr uintptr, idx uint16) Value { return *(*Value)(unsafe.Pointer(addr + uintptr(idx)*16)) }

func kodeaddr(code []uint16) uintptr { return (*reflect.SliceHeader)(unsafe.Pointer(&code)).Data }

// ExecCursor executes code under the given env from the given start cursor and returns:
// final result, yield cursor, is yield or not
func ExecCursor(env *Env, code []uint16, consts []Value, cursor uint32) (Value, uint32, bool) {
	var newEnv *Env
	var lastCursor uint32
	var lineinfo = "<unknown>"
	var retStack []ret
	var caddr = kodeaddr(code)
	var kaddr = (*reflect.SliceHeader)(unsafe.Pointer(&consts)).Data

	defer func() {
		if r := recover(); r != nil {
			e := &ExecError{r: r}
			e.stacks = make([]ret, len(retStack)+1)
			copy(e.stacks, retStack)
			e.stacks[len(e.stacks)-1] = ret{
				cursor: lastCursor,
				code:   code,
				line:   lineinfo,
			}
			panic(e)
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
		lastCursor = cursor
		// log.Println(cursor)
		bop := cruRead16(caddr, &cursor)
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
			env.Set(cruRead32(caddr, &cursor), env.Get(cruRead32(caddr, &cursor)))
		case OP_SETK:
			env.Set(cruRead32(caddr, &cursor), konst(kaddr, cruRead16(caddr, &cursor)))
		case OP_INC:
			addr := cruRead32(caddr, &cursor)
			num := env.Get(addr).AsNumber()
			env.Set(addr, NewNumberValue(num+konst(kaddr, cruRead16(caddr, &cursor)).AsNumber()))
		case OP_ADD:
			switch testTypes(env.R0, env.R1) {
			case _Tnumbernumber:
				env.A = NewNumberValue(env.R0.AsNumber() + env.R1.AsNumber())
			case _Tstringstring:
				env.A = NewStringValue(env.R0.AsString() + env.R1.AsString())
			case Tbytes<<8 | Tnumber:
				env.A = NewBytesValue(append(env.R0.AsBytes(), byte(env.R1.AsNumber())))
			default:
				if env.R0.ty == Tlist {
					env.A = NewListValue(append(env.R0.AsList(), env.R1))
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
			case _Tlistlist:
				env.A = NewListValue(append(env.R0.AsList(), env.R1.AsList()...))
			case _Tbytesbytes:
				env.A = NewBytesValue(append(env.R0.AsBytes(), env.R1.AsBytes()...))
			case _Tmapmap:
				tr, m := env.R0.AsMap().Dup(nil), env.R1.AsMap()
				if m.t != nil {
					for _, x := range m.t {
						tr.Put(x.k, x.v)
					}
				} else {
					for k, v := range m.m {
						tr.Put(k, v)
					}
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
			case Tlist:
				env.A = NewNumberValue(float64(len(v.AsList())))
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
		case OP_LIST:
			if newEnv == nil {
				env.A = NewListValue(make([]Value, 0))
			} else {
				list := make([]Value, newEnv.SSize())
				copy(list, newEnv.stack)
				newEnv.SClear()
				env.A = NewListValue(list)
			}
		case OP_MAP:
			if newEnv == nil {
				env.A = NewMapValue(NewMap())
			} else {
				size, m := newEnv.SSize(), NewMap()
				for i := 0; i < size; i += 2 {
					if k := newEnv.SGet(i); k.ty == Tstring {
						m.Put(k.AsString(), newEnv.SGet(i+1))
					} else {
						k.panicType(Tstring)
					}
				}
				newEnv.SClear()
				env.A = NewMapValue(m)
			}
		case OP_STORE:
			switch testTypes(env.R0, env.R1) {
			case Tbytes<<8 | Tnumber:
				if env.R2.ty == Tnumber {
					if b, idx := env.R0.AsBytes(), int(env.R1.AsNumber()); idx >= 0 {
						b[idx] = byte(env.R2.AsNumber())
					} else {
						b[len(b)+idx] = byte(env.R2.AsNumber())
					}
				} else {
					log.Panicf("can't store into %+v with key %+v", env.R0, env.R1)
				}
			case Tlist<<8 | Tnumber:
				if b, idx := env.R0.AsList(), int(env.R1.AsNumber()); idx >= 0 {
					b[idx] = env.R2
				} else {
					b[len(b)+idx] = env.R2
				}
			case Tmap<<8 | Tstring:
				env.R0.AsMap().Put(env.R1.AsString(), env.R2)
			default:
				log.Panicf("can't store into %+v with key %+v", env.R0, env.R1)
			}
			env.A = env.R2
		case OP_LOAD:
			var v Value
			switch testTypes(env.R0, env.R1) {
			case Tbytes<<8 | Tnumber:
				if b, idx := env.R0.AsBytes(), int(env.R1.AsNumber()); idx >= 0 {
					v = NewNumberValue(float64(b[idx]))
				} else {
					v = NewNumberValue(float64(b[len(b)+idx]))
				}
			case Tstring<<8 | Tnumber:
				if b, idx := env.R0.AsString(), int(env.R1.AsNumber()); idx >= 0 {
					v = NewNumberValue(float64(b[idx]))
				} else {
					v = NewNumberValue(float64(b[len(b)+idx]))
				}
			case Tlist<<8 | Tnumber:
				b, idx := env.R0.AsList(), int(env.R1.AsNumber())
				if idx >= 0 {
					v = b[idx]
				} else {
					v = b[len(b)+idx]
				}
				if v.Type() == Tclosure {
					v.AsClosure().SetCaller(env.R0)
				}
			case Tmap<<8 | Tstring:
				var found bool
				v, found = env.R0.AsMap().Get(env.R1.AsString())
				if v.Type() == Tclosure {
					v.AsClosure().SetCaller(env.R0)
				}
				if !found {
					env.E = NewBoolValue(true)
				}
			default:
				log.Panicf("can't load from %+v with key %+v", env.R0, env.R1)
			}
			env.A = v
		case OP_R0:
			env.R0 = env.Get(cruRead32(caddr, &cursor))
		case OP_R0K:
			env.R0 = konst(kaddr, cruRead16(caddr, &cursor))
		case OP_R1:
			env.R1 = env.Get(cruRead32(caddr, &cursor))
		case OP_R1K:
			env.R1 = konst(kaddr, cruRead16(caddr, &cursor))
		case OP_R2:
			env.R2 = env.Get(cruRead32(caddr, &cursor))
		case OP_R2K:
			env.R2 = konst(kaddr, cruRead16(caddr, &cursor))
		case OP_R3:
			env.R3 = env.Get(cruRead32(caddr, &cursor))
		case OP_R3K:
			env.R3 = konst(kaddr, cruRead16(caddr, &cursor))
		case OP_PUSH:
			if newEnv == nil {
				newEnv = NewEnv(nil)
			}
			newEnv.SPush(env.Get(cruRead32(caddr, &cursor)))
		case OP_PUSHK:
			if newEnv == nil {
				newEnv = NewEnv(nil)
			}
			newEnv.SPush(konst(kaddr, cruRead16(caddr, &cursor)))
		case OP_RET:
			v := env.Get(cruRead32(caddr, &cursor))
			if len(retStack) == 0 {
				return v, 0, false
			}
			returnUpperWorld(v)
		case OP_RETK:
			v := konst(kaddr, cruRead16(caddr, &cursor))
			if len(retStack) == 0 {
				return v, 0, false
			}
			returnUpperWorld(v)
		case OP_YIELD:
			return env.Get(cruRead32(caddr, &cursor)), cursor, true
		case OP_YIELDK:
			return konst(kaddr, cruRead16(caddr, &cursor)), cursor, true
		case OP_LAMBDA:
			metadata := cruRead32(caddr, &cursor)
			argsCount := byte(metadata >> 24)
			yieldable := byte(metadata>>16) == 1
			errorable := byte(metadata>>8) == 1
			noenvescape := byte(metadata) == 1
			constsLen := cruRead16(caddr, &cursor)
			consts := make([]Value, constsLen)
			for i := uint16(0); i < constsLen; i++ {
				switch cruRead16(caddr, &cursor) {
				case Tnumber:
					consts[i] = NewNumberValue(crReadDouble(code, &cursor))
				case Tstring:
					consts[i] = NewStringValue(crReadString(code, &cursor))
				default:
					panic("shouldn't happen")
				}
			}
			buf := crRead(code, &cursor, int(cruRead32(caddr, &cursor)))
			env.A = NewClosureValue(NewClosure(buf, consts, env, byte(argsCount), yieldable, errorable, noenvescape))
		case OP_CALL:
			switch v := env.Get(cruRead32(caddr, &cursor)); v.ty {
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
			case Tlist:
				if bb := v.AsList(); newEnv.SSize() == 2 {
					env.A = NewListValue(bb[shiftIndex(newEnv.SGet(0), len(bb)):shiftIndex(newEnv.SGet(1), len(bb))])
				} else if newEnv.SSize() == 1 {
					env.A = bb[shiftIndex(newEnv.SGet(0), len(bb))]
				} else {
					log.Panicf("invalid call on %v", v)
				}
				newEnv.SClear()
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
			off := int32(cruRead32(caddr, &cursor))
			cursor = uint32(int32(cursor) + off)
		case OP_IFNOT:
			cond := env.Get(cruRead32(caddr, &cursor))
			off := int32(cruRead32(caddr, &cursor))
			if cond.ty == Tbool && !cond.AsBool() {
				cursor = uint32(int32(cursor) + off)
			} else if cond.IsFalse() {
				cursor = uint32(int32(cursor) + off)
			}
		case OP_IF:
			cond := env.Get(cruRead32(caddr, &cursor))
			off := int32(cruRead32(caddr, &cursor))
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
			stack := env.Stack()
			ret := make([]Value, len(stack))
			copy(ret, stack)
			env.A = NewListValue(ret)
			return
		case 2:
			// return dup()
			env.A = NewListValue(env.Stack())
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
		case Tlist:
			list0 := env.R1.AsList()
			list1 := make([]Value, len(list0))
			copy(list1, list0)
			env.A = NewListValue(list1)
			return
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
	case Tlist:
		var list0 = env.R1.AsList()
		var list1 []Value
		if alloc {
			list1 = make([]Value, 0, len(list0))
		}
		for i, v := range list0 {
			newEnv.SClear()
			newEnv.SPush(NewNumberValue(float64(i)))
			newEnv.SPush(v)
			ret, _, _ := ExecCursor(newEnv, cls.code, cls.consts, 0)
			if newEnv.E.Type() != Tnil {
				break
			}
			if alloc {
				list1 = append(list1, ret)
			}
		}
		if alloc {
			env.A = NewListValue(list1)
		}
	case Tmap:
		if alloc {
			if cls.errorable {
				// the predicator may return error and interrupt the dup, so full copy is not used here
				// however cls.errorable is not 100% accurate because calling error() (to check error) and
				// calling error(...) (to throw error) are different behaviors, but i will left this as a TODO
				m2 := NewMap()
				m := env.R1.AsMap()
				if m.t != nil {
					for _, x := range m.t {
						newEnv.SClear()
						newEnv.SPush(NewStringValue(x.k))
						newEnv.SPush(x.v)
						ret, _, _ := ExecCursor(newEnv, cls.code, cls.consts, 0)
						if newEnv.E.Type() != Tnil {
							break
						}
						m2.Put(x.k, ret)
					}
				} else {
					m2.SwitchToHashmap()
					for k, v := range m.m {
						newEnv.SClear()
						newEnv.SPush(NewStringValue(k))
						newEnv.SPush(v)
						ret, _, _ := ExecCursor(newEnv, cls.code, cls.consts, 0)
						if newEnv.E.Type() != Tnil {
							break
						}
						m2.Put(k, ret)
					}
				}
				env.A = NewMapValue(m2)
			} else {
				// full copy
				env.A = NewMapValue(env.R1.AsMap().Dup(func(k string, v Value) Value {
					newEnv.SClear()
					newEnv.SPush(NewStringValue(k))
					newEnv.SPush(v)
					ret, _, _ := ExecCursor(newEnv, cls.code, cls.consts, 0)
					return ret
				}))
			}
		} else {
			m := env.R1.AsMap()
			for _, x := range m.t {
				newEnv.SClear()
				newEnv.SPush(NewStringValue(x.k))
				newEnv.SPush(x.v)
				ExecCursor(newEnv, cls.code, cls.consts, 0)
				if newEnv.E.Type() != Tnil {
					break
				}
			}
			for k, v := range m.m {
				newEnv.SClear()
				newEnv.SPush(NewStringValue(k))
				newEnv.SPush(v)
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
