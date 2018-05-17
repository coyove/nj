package base

import (
	"fmt"
	"log"
	"os"
)

func Exec(env *Env, code []byte) Value {
	v, _, _ := ExecCursor(env, code, 0)
	return v
}

func ExecCursor(env *Env, code []byte, cursor uint32) (Value, uint32, bool) {
	var newEnv *Env
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			log.Println(fmt.Sprintf("%x", crHash(code)))
			log.Println("cursor:", cursor)
			os.Exit(1)
		}
	}()

	type ret struct {
		cursor uint32
		env    *Env
		code   []byte
	}

	var retStack []ret

	returnUpperWorld := func(v Value) {
		r := retStack[len(retStack)-1]
		cursor = r.cursor
		code = r.code
		env = r.env
		env.A = v
		retStack = retStack[:len(retStack)-1]
	}

	for {
		bop := crReadByte(code, &cursor)
		if bop == OP_EOB {
			break
		}
		switch bop {
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
			env.Set(crReadInt32(code, &cursor), env.Get(crReadInt32(code, &cursor)))
		case OP_SET_NUM:
			env.Set(crReadInt32(code, &cursor), NewNumberValue(crReadDouble(code, &cursor)))
		case OP_SET_STR:
			env.Set(crReadInt32(code, &cursor), NewStringValue(crReadString(code, &cursor)))
		case OP_ADD:
			switch l := env.R0; l.Type() {
			case Tnumber:
				env.A = NewNumberValue(l.AsNumberUnsafe() + env.R1.AsNumber())
			case Tstring:
				env.A = NewStringValue(l.AsStringUnsafe() + env.R1.AsString())
			case Tlist:
				env.A = NewListValue(append(l.AsListUnsafe(), env.R1))
			case Tbytes:
				env.A = NewBytesValue(append(l.AsBytesUnsafe(), byte(env.R1.AsNumber())))
			default:
				log.Panicf("can't apply 'add' on %+v", l)
			}
		case OP_SUB:
			env.A = NewNumberValue(env.R0.AsNumber() - env.R1.AsNumber())
		case OP_MUL:
			env.A = NewNumberValue(env.R0.AsNumber() * env.R1.AsNumber())
		case OP_DIV:
			env.A = NewNumberValue(env.R0.AsNumber() / env.R1.AsNumber())
		case OP_MOD:
			env.A = NewNumberValue(float64(int64(env.R0.AsNumber()) % int64(env.R1.AsNumber())))
		case OP_EQ:
			env.A = NewBoolValue(env.R0.Equal(env.R1))
		case OP_NEQ:
			env.A = NewBoolValue(!env.R0.Equal(env.R1))
		case OP_LESS:
			env.A = NewBoolValue(env.R0.Less(env.R1))
		case OP_MORE_EQ:
			env.A = NewBoolValue(!env.R0.Less(env.R1))
		case OP_LESS_EQ:
			env.A = NewBoolValue(env.R0.LessEqual(env.R1))
		case OP_MORE:
			env.A = NewBoolValue(!env.R0.LessEqual(env.R1))
		case OP_NOT:
			if env.R0.IsFalse() {
				env.A = NewBoolValue(true)
			} else {
				env.A = NewBoolValue(false)
			}
		case OP_AND:
			env.A = NewBoolValue(!env.R0.IsFalse() && !env.R1.IsFalse())
		case OP_OR:
			if env.R0.IsFalse() {
				if env.R1.IsFalse() {
					env.A = NewBoolValue(false)
				} else {
					env.A = env.R1
				}
			} else {
				env.A = env.R0
			}
		case OP_XOR:
			env.A = NewBoolValue(env.R0.IsFalse() != env.R1.IsFalse())
		case OP_BIT_NOT:
			env.A = NewNumberValue(float64(^int64(env.R0.AsNumber())))
		case OP_BIT_AND:
			switch env.R0.Type() {
			case Tnumber:
				env.A = NewNumberValue(float64(int64(env.R0.AsNumberUnsafe()) & int64(env.R1.AsNumber())))
			case Tlist:
				env.A = NewListValue(append(env.R0.AsListUnsafe(), env.R1.AsList()...))
			case Tbytes:
				env.A = NewBytesValue(append(env.R0.AsBytesUnsafe(), env.R1.AsBytes()...))
			case Tmap:
				tr, m := env.R0.AsMapUnsafe().Dup(), env.R1.AsMap()
				for iter := m.Iterator(); iter.Next(); {
					tr.Put(iter.Key(), iter.Value())
				}
				env.A = NewMapValue(tr)
			default:
				log.Panicf("can't apply bit 'and' on %+v", env.R0)
			}
		case OP_BIT_OR:
			env.A = NewNumberValue(float64(int64(env.R0.AsNumber()) | int64(env.R1.AsNumber())))
		case OP_BIT_XOR:
			env.A = NewNumberValue(float64(int64(env.R0.AsNumber()) ^ int64(env.R1.AsNumber())))
		case OP_BIT_LSH:
			env.A = NewNumberValue(float64(uint64(env.R0.AsNumber()) << uint64(env.R1.AsNumber())))
		case OP_BIT_RSH:
			env.A = NewNumberValue(float64(uint64(env.R0.AsNumber()) >> uint64(env.R1.AsNumber())))
		case OP_ASSERT:
			loc := "assertion failed: " + crReadString(code, &cursor)
			if env.R0.IsFalse() {
				panic(loc)
			}
			env.A = NewBoolValue(true)
		case OP_LEN:
			switch v := env.R0; v.Type() {
			case Tstring:
				env.A = NewNumberValue(float64(len(v.AsStringUnsafe())))
			case Tlist:
				env.A = NewNumberValue(float64(len(v.AsListUnsafe())))
			case Tmap:
				env.A = NewNumberValue(float64(v.AsMapUnsafe().Size()))
			case Tbytes:
				env.A = NewNumberValue(float64(len(v.AsBytesUnsafe())))
			default:
				log.Panicf("can't evaluate the length of %+v", v)
			}
		case OP_LIST:
			if newEnv == nil {
				env.A = NewListValue(make([]Value, 0))
			} else {
				size := newEnv.Stack().Size()
				list := make([]Value, size)
				for i := 0; i < size; i++ {
					list[i] = newEnv.Get(int32(i))
				}
				newEnv.Stack().Clear()
				env.A = NewListValue(list)
			}
		case OP_MAP:
			if newEnv == nil {
				env.A = NewMapValue(new(Tree))
			} else {
				size := newEnv.Stack().Size()
				m := new(Tree)
				for i := 0; i < size; i += 2 {
					m.Put(newEnv.Get(int32(i)).AsString(), newEnv.Get(int32(i+1)))
				}
				newEnv.Stack().Clear()
				env.A = NewMapValue(m)
			}
		case OP_STORE:
			switch env.R0.Type() {
			case Tbytes:
				if b, idx := env.R0.AsBytesUnsafe(), int(env.R1.AsNumber()); idx >= 0 {
					b[idx] = byte(env.R2.AsNumber())
				} else {
					b[len(b)+idx] = byte(env.R2.AsNumber())
				}
			case Tlist:
				if b, idx := env.R0.AsListUnsafe(), int(env.R1.AsNumber()); idx >= 0 {
					b[idx] = env.R2
				} else {
					b[len(b)+idx] = env.R2
				}
			case Tmap:
				env.R0.AsMapUnsafe().Put(env.R1.AsString(), env.R2)
			default:
				log.Panicf("can't store into %+v", env.R0)
			}
			env.A = env.R2
		case OP_LOAD:
			var v Value
			switch env.R0.Type() {
			case Tbytes:
				if b, idx := env.R0.AsBytesUnsafe(), int(env.R1.AsNumber()); idx >= 0 {
					v = NewNumberValue(float64(b[idx]))
				} else {
					v = NewNumberValue(float64(b[len(b)+idx]))
				}
			case Tstring:
				if b, idx := env.R0.AsStringUnsafe(), int(env.R1.AsNumber()); idx >= 0 {
					v = NewNumberValue(float64(b[idx]))
				} else {
					v = NewNumberValue(float64(b[len(b)+idx]))
				}
			case Tlist:
				b, idx := env.R0.AsListUnsafe(), int(env.R1.AsNumber())
				if idx >= 0 {
					v = b[idx]
				} else {
					v = b[len(b)+idx]
				}
				if v.Type() == Tclosure {
					v.AsClosureUnsafe().SetCaller(env.R0)
				}
			case Tmap:
				v, _ = env.R0.AsMapUnsafe().Get(env.R1.AsString())
				if v.Type() == Tclosure {
					v.AsClosureUnsafe().SetCaller(env.R0)
				}
			default:
				log.Panicf("can't load from %+v", env.R0)
			}
			env.A = v
		case OP_SAFE_STORE:
			switch idx := int(env.R1.AsNumber()); env.R0.Type() {
			case Tbytes:
				if bb := env.R0.AsBytesUnsafe(); idx < len(bb) {
					bb[idx] = byte(env.R2.AsNumber())
				}
			case Tlist:
				if bb := env.R0.AsListUnsafe(); idx < len(bb) {
					bb[idx] = env.R2
				}
			case Tmap:
				env.R0.AsMapUnsafe().Put(env.R1.AsString(), env.R2)
			default:
				log.Panicf("can't safe store into %+v", env.R0)
			}
			env.A = env.R2
		case OP_SAFE_LOAD:
			v := NewValue()
			switch idx := int(env.R1.AsNumber()); env.R0.Type() {
			case Tbytes:
				if bb := env.R0.AsBytesUnsafe(); idx < len(bb) {
					v = NewNumberValue(float64(bb[idx]))
				}
			case Tstring:
				if bb := env.R0.AsStringUnsafe(); idx < len(bb) {
					v = NewNumberValue(float64(bb[idx]))
				}
			case Tlist:
				if bb := env.R0.AsListUnsafe(); idx < len(bb) {
					v = bb[idx]
					if v.Type() == Tclosure {
						v.AsClosureUnsafe().SetCaller(env.R0)
					}
				}
			case Tmap:
				v, _ = env.R0.AsMapUnsafe().Get(env.R1.AsString())
				if v.Type() == Tclosure {
					v.AsClosureUnsafe().SetCaller(env.R0)
				}
			default:
				log.Panicf("can't safe load from %+v", env.R0)
			}
			env.A = v
		case OP_R0:
			env.R0 = env.Get(crReadInt32(code, &cursor))
		case OP_R0_NUM:
			env.R0 = NewNumberValue(crReadDouble(code, &cursor))
		case OP_R0_STR:
			env.R0 = NewStringValue(crReadString(code, &cursor))
		case OP_R1:
			env.R1 = env.Get(crReadInt32(code, &cursor))
		case OP_R1_NUM:
			env.R1 = NewNumberValue(crReadDouble(code, &cursor))
		case OP_R1_STR:
			env.R1 = NewStringValue(crReadString(code, &cursor))
		case OP_R2:
			env.R2 = env.Get(crReadInt32(code, &cursor))
		case OP_R2_NUM:
			env.R2 = NewNumberValue(crReadDouble(code, &cursor))
		case OP_R2_STR:
			env.R2 = NewStringValue(crReadString(code, &cursor))
		case OP_R3:
			env.R3 = env.Get(crReadInt32(code, &cursor))
		case OP_R3_NUM:
			env.R3 = NewNumberValue(crReadDouble(code, &cursor))
		case OP_R3_STR:
			env.R3 = NewStringValue(crReadString(code, &cursor))
		case OP_PUSH:
			if newEnv == nil {
				newEnv = NewEnv(nil)
			}
			newEnv.Push(env.Get(crReadInt32(code, &cursor)))
		case OP_PUSH_NUM:
			if newEnv == nil {
				newEnv = NewEnv(nil)
			}
			newEnv.Push(NewNumberValue(crReadDouble(code, &cursor)))
		case OP_PUSH_STR:
			if newEnv == nil {
				newEnv = NewEnv(nil)
			}
			newEnv.Push(NewStringValue(crReadString(code, &cursor)))
		case OP_RET:
			v := env.Get(crReadInt32(code, &cursor))
			if len(retStack) == 0 {
				return v, 0, false
			}
			returnUpperWorld(v)
		case OP_RET_NUM:
			v := NewNumberValue(crReadDouble(code, &cursor))
			if len(retStack) == 0 {
				return v, 0, false
			}
			returnUpperWorld(v)
		case OP_RET_STR:
			v := NewStringValue(crReadString(code, &cursor))
			if len(retStack) == 0 {
				return v, 0, false
			}
			returnUpperWorld(v)
		case OP_YIELD:
			return env.Get(crReadInt32(code, &cursor)), cursor, true
		case OP_YIELD_NUM:
			return NewNumberValue(crReadDouble(code, &cursor)), cursor, true
		case OP_YIELD_STR:
			return NewStringValue(crReadString(code, &cursor)), cursor, true
		case OP_LAMBDA:
			argsCount := int(crReadInt32(code, &cursor))
			yieldable := crReadByte(code, &cursor) == 1
			buf := crReadBytes(code, &cursor, int(crReadInt32(code, &cursor)))
			env.A = NewClosureValue(NewClosure(buf, env, argsCount, yieldable))
		case OP_CALL:
			v := env.Get(crReadInt32(code, &cursor))
			switch v.Type() {
			case Tclosure:
				cls := v.AsClosureUnsafe()
				if newEnv == nil {
					newEnv = NewEnv(nil)
				}

				if newEnv.Size() < cls.ArgsCount() {
					if newEnv.Size() == 0 {
						env.A = (NewClosureValue(cls))
					} else {
						curry := cls.Dup()
						curry.AppendPreArgs(newEnv.Stack().Values())
						env.A = (NewClosureValue(curry))
					}
				} else {
					if cls.PreArgs() != nil && len(cls.PreArgs()) > 0 {
						newEnv.Stack().Insert(0, cls.PreArgs())
					}

					if cls.yieldable || cls.native != nil {
						env.A = cls.Exec(newEnv)
					} else {
						if retStack == nil {
							retStack = make([]ret, 0, 1)
						}

						last := ret{
							cursor: cursor,
							env:    env,
							code:   code,
						}

						// switch to the env of cls
						cursor = 0
						newEnv.parent = cls.env
						env = newEnv
						code = cls.code

						retStack = append(retStack, last)
					}
				}

				newEnv = nil
			case Tlist:
				if bb := v.AsListUnsafe(); newEnv.Size() == 2 {
					env.A = NewListValue(bb[shiftIndex(newEnv.Get(0), len(bb)):shiftIndex(newEnv.Get(1), len(bb))])
				} else if newEnv.Size() == 1 {
					env.A = bb[shiftIndex(newEnv.Get(0), len(bb))]
				} else {
					log.Panicf("too many (or few) arguments to call on list")
				}
				newEnv.Stack().Clear()
			case Tbytes:
				if bb := v.AsBytesUnsafe(); newEnv.Size() == 2 {
					env.A = NewBytesValue(bb[shiftIndex(newEnv.Get(0), len(bb)):shiftIndex(newEnv.Get(1), len(bb))])
				} else if newEnv.Size() == 1 {
					env.A = NewNumberValue(float64(bb[shiftIndex(newEnv.Get(0), len(bb))]))
				} else {
					log.Panicf("too many (or few) arguments to call on bytes")
				}
				newEnv.Stack().Clear()
			case Tstring:
				if bb := v.AsStringUnsafe(); newEnv.Size() == 2 {
					env.A = NewStringValue(bb[shiftIndex(newEnv.Get(0), len(bb)):shiftIndex(newEnv.Get(1), len(bb))])
				} else if newEnv.Size() == 1 {
					env.A = NewNumberValue(float64(bb[shiftIndex(newEnv.Get(0), len(bb))]))
				} else {
					log.Panicf("too many (or few) arguments to call on string")
				}
				newEnv.Stack().Clear()
			default:
				log.Panicf("invalid callee: %+v", v)
			}
		case OP_STACK:
			env.A = NewListValue(env.Stack().data)
		case OP_JMP:
			off := int(crReadInt32(code, &cursor))
			*&cursor += uint32(off)
		case OP_IF:
			cond := env.Get(crReadInt32(code, &cursor))
			off := int(crReadInt32(code, &cursor))
			if cond.IsFalse() {
				*&cursor += uint32(off)
			}
		case OP_DUP:
			switch env.R0.Type() {
			case Tnil, Tnumber, Tstring, Tbool, Tgeneric:
				env.A = env.R0
			case Tclosure:
				env.A = NewClosureValue(env.R0.AsClosureUnsafe().Dup())
			case Tlist:
				list0 := env.R0.AsListUnsafe()
				list1 := make([]Value, len(list0))
				copy(list1, list0)
				env.A = NewListValue(list1)
			case Tmap:
				env.A = NewMapValue(env.R0.AsMapUnsafe().Dup())
			case Tbytes:
				bytes0 := env.R0.AsBytesUnsafe()
				bytes1 := make([]byte, len(bytes0))
				copy(bytes1, bytes0)
				env.A = NewBytesValue(bytes1)
			}
		}
	}

	return NewValue(), 0, false
}

func shiftIndex(index Value, len int) int {
	i := int(index.AsNumberUnsafe())
	if i >= 0 {
		return i
	}
	return i + len
}
