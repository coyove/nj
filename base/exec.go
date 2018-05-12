package base

import (
	"fmt"
	"log"
	"os"
)

func Exec(env *Env, code []byte) Value {
	var newEnv *Env
	c := NewBytesReader(code)

	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			log.Println(fmt.Sprintf("%x", NewBytesReader(code).Hash()))
			log.Println("cursor:", c.GetCursor())
			os.Exit(1)
		}
	}()

	for {
		bop := c.ReadByte()
		if bop == OP_EOB {
			break
		}
		// log.Println(c.GetCursor())
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
			env.Set(c.ReadInt32(), env.Get(c.ReadInt32()))
		case OP_SET_NUM:
			env.Set(c.ReadInt32(), NewNumberValue(c.ReadDouble()))
		case OP_SET_STR:
			env.Set(c.ReadInt32(), NewStringValue(c.ReadString()))
		case OP_ADD:
			switch l := env.R0; l.Type() {
			case Tnumber:
				env.A = NewNumberValue(l.AsNumberUnsafe() + env.R1.AsNumber())
			case Tstring:
				env.A = NewStringValue(l.AsStringUnsafe() + env.R1.AsString())
			case Tlist:
				env.A = NewListValue(append(l.AsListUnsafe(), env.R1))
			default:
				log.Panicf("can't add %+v", l)
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
			env.A = NewNumberValue(float64(int64(env.R0.AsNumber()) & int64(env.R1.AsNumber())))
		case OP_BIT_OR:
			env.A = NewNumberValue(float64(int64(env.R0.AsNumber()) | int64(env.R1.AsNumber())))
		case OP_BIT_XOR:
			env.A = NewNumberValue(float64(int64(env.R0.AsNumber()) ^ int64(env.R1.AsNumber())))
		case OP_BIT_LSH:
			env.A = NewNumberValue(float64(uint64(env.R0.AsNumber()) << uint64(env.R1.AsNumber())))
		case OP_BIT_RSH:
			env.A = NewNumberValue(float64(uint64(env.R0.AsNumber()) >> uint64(env.R1.AsNumber())))
		case OP_ASSERT:
			loc := "assertion failed: " + c.ReadString()
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
				bb := env.R0.AsBytesUnsafe()
				bb[(idx+len(bb))%len(bb)] = byte(env.R2.AsNumber())
			case Tlist:
				bb := env.R0.AsListUnsafe()
				bb[(idx+len(bb))%len(bb)] = env.R2
			default:
				log.Panicf("can't safe store into %+v", env.R0)
			}
			env.A = env.R2
		case OP_SAFE_LOAD:
			v := NewValue()
			switch idx := int(env.R1.AsNumber()); env.R0.Type() {
			case Tbytes:
				bb := env.R0.AsBytesUnsafe()
				v = NewNumberValue(float64(bb[(idx+len(bb))%len(bb)]))
			case Tstring:
				bb := env.R0.AsStringUnsafe()
				v = NewNumberValue(float64(bb[(idx+len(bb))%len(bb)]))
			case Tlist:
				bb := env.R0.AsListUnsafe()
				v = bb[(idx+len(bb))%len(bb)]
				if v.Type() == Tclosure {
					v.AsClosureUnsafe().SetCaller(env.R0)
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
			env.R0 = env.Get(c.ReadInt32())
		case OP_R0_NUM:
			env.R0 = NewNumberValue(c.ReadDouble())
		case OP_R0_STR:
			env.R0 = NewStringValue(c.ReadString())
		case OP_R1:
			env.R1 = env.Get(c.ReadInt32())
		case OP_R1_NUM:
			env.R1 = NewNumberValue(c.ReadDouble())
		case OP_R1_STR:
			env.R1 = NewStringValue(c.ReadString())
		case OP_R2:
			env.R2 = env.Get(c.ReadInt32())
		case OP_R2_NUM:
			env.R2 = NewNumberValue(c.ReadDouble())
		case OP_R2_STR:
			env.R2 = NewStringValue(c.ReadString())
		case OP_R3:
			env.R3 = env.Get(c.ReadInt32())
		case OP_R3_NUM:
			env.R3 = NewNumberValue(c.ReadDouble())
		case OP_R3_STR:
			env.R3 = NewStringValue(c.ReadString())
		case OP_PUSH:
			if newEnv == nil {
				newEnv = NewEnv(nil)
			}
			newEnv.Push(env.Get(c.ReadInt32()))
		case OP_PUSH_NUM:
			if newEnv == nil {
				newEnv = NewEnv(nil)
			}
			newEnv.Push(NewNumberValue(c.ReadDouble()))
		case OP_PUSH_STR:
			if newEnv == nil {
				newEnv = NewEnv(nil)
			}
			newEnv.Push(NewStringValue(c.ReadString()))
		case OP_RET:
			return env.Get(c.ReadInt32())
		case OP_RET_NUM:
			return NewNumberValue(c.ReadDouble())
		case OP_RET_STR:
			return NewStringValue(c.ReadString())
		case OP_LAMBDA:
			argsCount := int(c.ReadInt32())
			buf := c.Read(int(c.ReadInt32()))
			env.A = NewClosureValue(NewClosure(buf, env, argsCount))
		case OP_CALL:
			v := env.Get(c.ReadInt32())
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
					env.A = cls.Exec(newEnv)
				}

				newEnv = nil
			case Tlist:
				if bb := v.AsListUnsafe(); newEnv.Size() == 2 {
					env.A = NewListValue(bb[shiftIndex(newEnv.Get(0), len(bb)):shiftIndex(newEnv.Get(1), len(bb))])
				} else {
					log.Panicf("too many (or few) arguments to call on list")
				}
				newEnv.Stack().Clear()
			case Tbytes:
				if bb := v.AsBytesUnsafe(); newEnv.Size() == 2 {
					env.A = NewBytesValue(bb[shiftIndex(newEnv.Get(0), len(bb)):shiftIndex(newEnv.Get(1), len(bb))])
				} else {
					log.Panicf("too many (or few) arguments to call on bytes")
				}
				newEnv.Stack().Clear()
			case Tstring:
				if bb := v.AsStringUnsafe(); newEnv.Size() == 2 {
					env.A = NewStringValue(bb[shiftIndex(newEnv.Get(0), len(bb)):shiftIndex(newEnv.Get(1), len(bb))])
				} else {
					log.Panicf("too many (or few) arguments to call on string")
				}
				newEnv.Stack().Clear()
			default:
				log.Panicf("invalid callee: %+v", v)
			}
		case OP_VARARGS:
			list0 := env.Stack().Values()
			list1 := make([]Value, len(list0))
			copy(list1, list0)
			env.A = NewListValue(list1)
		case OP_JMP:
			off := int(c.ReadInt32())
			c.SetCursor(c.GetCursor() + off)
		case OP_IF:
			cond := env.Get(c.ReadInt32())
			off := int(c.ReadInt32())
			if cond.IsFalse() {
				c.SetCursor(c.GetCursor() + off)
			}
		case OP_DUP:
			switch env.R0.Type() {
			case Tnil, Tnumber, Tstring, Tbool, Tclosure, Tgeneric:
				env.A = env.R0
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

	return NewValue()
}

func shiftIndex(index Value, len int) int {
	i := int(index.AsNumberUnsafe())
	if i >= 0 {
		return i
	}
	return i + len
}
