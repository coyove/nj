package vm

import (
	"log"
	"os"

	"github.com/coyove/bracket/base"
)

func Exec(env *base.Env, code []byte) base.Value {
	var newEnv *base.Env
	c := base.NewBytesReader(code)

	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			log.Println("cursor:", c.GetCursor())
			os.Exit(1)
		}
	}()

	for {
		bop := c.ReadByte()
		if bop == base.OP_EOB {
			break
		}
		// log.Println(c.GetCursor())
		switch bop {
		case base.OP_NOP:
		case base.OP_WHO:
			env.A = env.C
		case base.OP_NIL:
			env.A = base.NewValue()
		case base.OP_TRUE:
			env.A = base.NewBoolValue(true)
		case base.OP_FALSE:
			env.A = base.NewBoolValue(false)
		case base.OP_BYTES:
			if n := env.R0; n.Type() == base.Tstring {
				env.A = base.NewBytesValue([]byte(n.AsStringUnsafe()))
			} else if n.Type() == base.Tnumber {
				env.A = base.NewBytesValue(make([]byte, int(n.AsNumberUnsafe())))
			} else {
				panic("can't generate the bytes")
			}
		case base.OP_SET:
			env.Set(c.ReadInt32(), env.Get(c.ReadInt32()))
		case base.OP_SET_NUM:
			env.Set(c.ReadInt32(), base.NewNumberValue(c.ReadDouble()))
		case base.OP_SET_STR:
			env.Set(c.ReadInt32(), base.NewStringValue(c.ReadString()))
		case base.OP_ADD:
			switch l := env.R0; l.Type() {
			case base.Tnumber:
				env.A = base.NewNumberValue(l.AsNumberUnsafe() + env.R1.AsNumber())
			case base.Tstring:
				env.A = base.NewStringValue(l.AsStringUnsafe() + env.R1.AsString())
			case base.Tlist:
				env.A = base.NewListValue(append(l.AsListUnsafe(), env.R1))
			default:
				log.Panicf("can't add %+v", l)
			}
		case base.OP_SUB:
			env.A = base.NewNumberValue(env.R0.AsNumber() - env.R1.AsNumber())
		case base.OP_MUL:
			env.A = base.NewNumberValue(env.R0.AsNumber() * env.R1.AsNumber())
		case base.OP_DIV:
			env.A = base.NewNumberValue(env.R0.AsNumber() / env.R1.AsNumber())
		case base.OP_MOD:
			env.A = base.NewNumberValue(float64(int64(env.R0.AsNumber()) % int64(env.R1.AsNumber())))
		case base.OP_INC:
			id := c.ReadInt32()
			v := env.Get(id).AsNumber() + env.Get(c.ReadInt32()).AsNumber()
			env.Set(id, base.NewNumberValue(v))
			env.A = base.NewNumberValue(v)
		case base.OP_INC_NUM:
			id := c.ReadInt32()
			v := env.Get(id).AsNumber() + c.ReadDouble()
			env.Set(id, base.NewNumberValue(v))
			env.A = base.NewNumberValue(v)
		case base.OP_EQ:
			env.A = base.NewBoolValue(env.R0.Equal(env.R1))
		case base.OP_NEQ:
			env.A = base.NewBoolValue(!env.R0.Equal(env.R1))
		case base.OP_LESS:
			env.A = base.NewBoolValue(env.R0.Less(env.R1))
		case base.OP_MORE_EQ:
			env.A = base.NewBoolValue(!env.R0.Less(env.R1))
		case base.OP_LESS_EQ:
			env.A = base.NewBoolValue(env.R0.LessEqual(env.R1))
		case base.OP_MORE:
			env.A = base.NewBoolValue(!env.R0.LessEqual(env.R1))
		case base.OP_NOT:
			env.A = base.NewBoolValue(!env.R0.AsBool())
		case base.OP_AND:
			env.A = base.NewBoolValue(!env.R0.IsFalse() && !env.R1.IsFalse())
		case base.OP_OR:
			if env.R0.IsFalse() {
				if env.R1.IsFalse() {
					env.A = base.NewBoolValue(false)
				} else {
					env.A = env.R1
				}
			} else {
				env.A = env.R0
			}
		case base.OP_XOR:
			env.A = base.NewBoolValue(env.R0.IsFalse() != env.R1.IsFalse())
		case base.OP_BIT_NOT:
			env.A = base.NewNumberValue(float64(^int64(env.R0.AsNumber())))
		case base.OP_BIT_AND:
			env.A = base.NewNumberValue(float64(int64(env.R0.AsNumber()) & int64(env.R1.AsNumber())))
		case base.OP_BIT_OR:
			env.A = base.NewNumberValue(float64(int64(env.R0.AsNumber()) | int64(env.R1.AsNumber())))
		case base.OP_BIT_XOR:
			env.A = base.NewNumberValue(float64(int64(env.R0.AsNumber()) ^ int64(env.R1.AsNumber())))
		case base.OP_BIT_LSH:
			env.A = base.NewNumberValue(float64(uint64(env.R0.AsNumber()) << uint64(env.R1.AsNumber())))
		case base.OP_BIT_RSH:
			env.A = base.NewNumberValue(float64(uint64(env.R0.AsNumber()) >> uint64(env.R1.AsNumber())))
		case base.OP_ASSERT:
			loc := "assertion failed: " + c.ReadString()
			if env.R0.IsFalse() {
				panic(loc)
			}
			env.A = base.NewBoolValue(true)
		case base.OP_LEN:
			switch v := env.R0; v.Type() {
			case base.Tstring:
				env.A = base.NewNumberValue(float64(len(v.AsStringUnsafe())))
			case base.Tlist:
				env.A = base.NewNumberValue(float64(len(v.AsListUnsafe())))
			case base.Tmap:
				env.A = base.NewNumberValue(float64(len(v.AsMapUnsafe())))
			case base.Tbytes:
				env.A = base.NewNumberValue(float64(len(v.AsBytesUnsafe())))
			default:
				log.Panicf("can't evaluate the length of %+v", v)
			}
		case base.OP_LIST:
			size := newEnv.Stack().Size()
			list := make([]base.Value, size)
			for i := 0; i < size; i++ {
				list[i] = newEnv.Get(int32(i))
			}
			newEnv.Stack().Clear()
			env.A = base.NewListValue(list)
		case base.OP_MAP:
			size := newEnv.Stack().Size()
			m := make(map[string]base.Value)
			for i := 0; i < size; i += 2 {
				m[newEnv.Get(int32(i)).AsString()] = newEnv.Get(int32(i + 1))
			}
			newEnv.Stack().Clear()
			env.A = base.NewMapValue(m)
		case base.OP_STORE:
			switch env.R0.Type() {
			case base.Tbytes:
				if b, idx := env.R0.AsBytesUnsafe(), int(env.R1.AsNumber()); idx >= 0 {
					b[idx] = byte(env.R2.AsNumber())
				} else {
					b[len(b)+idx] = byte(env.R2.AsNumber())
				}
			case base.Tlist:
				if b, idx := env.R0.AsListUnsafe(), int(env.R1.AsNumber()); idx >= 0 {
					b[idx] = env.R2
				} else {
					b[len(b)+idx] = env.R2
				}
			case base.Tmap:
				env.R0.AsMapUnsafe()[env.R1.AsString()] = env.R2
			default:
				log.Panicf("can't store into %+v", env.R0)
			}
			env.A = env.R2
		case base.OP_LOAD:
			var v base.Value
			switch env.R0.Type() {
			case base.Tbytes:
				if b, idx := env.R0.AsBytesUnsafe(), int(env.R1.AsNumber()); idx >= 0 {
					v = base.NewNumberValue(float64(b[idx]))
				} else {
					v = base.NewNumberValue(float64(b[len(b)+idx]))
				}
			case base.Tstring:
				if b, idx := env.R0.AsStringUnsafe(), int(env.R1.AsNumber()); idx >= 0 {
					v = base.NewNumberValue(float64(b[idx]))
				} else {
					v = base.NewNumberValue(float64(b[len(b)+idx]))
				}
			case base.Tlist:
				b, idx := env.R0.AsListUnsafe(), int(env.R1.AsNumber())
				if idx >= 0 {
					v = b[idx]
				} else {
					v = b[len(b)+idx]
				}
				if v.Type() == base.Tclosure {
					v.AsClosureUnsafe().SetCaller(env.R0)
				}
			case base.Tmap:
				v = env.R0.AsMapUnsafe()[env.R1.AsString()]
				if v.Type() == base.Tclosure {
					v.AsClosureUnsafe().SetCaller(env.R0)
				}
			default:
				log.Panicf("can't load from %+v", env.R0)
			}
			env.A = v
		case base.OP_SAFE_STORE:
			switch idx := int(env.R1.AsNumber()); env.R0.Type() {
			case base.Tbytes:
				bb := env.R0.AsBytesUnsafe()
				bb[(idx+len(bb))%len(bb)] = byte(env.R2.AsNumber())
			case base.Tlist:
				bb := env.R0.AsListUnsafe()
				bb[(idx+len(bb))%len(bb)] = env.R2
			default:
				log.Panicf("can't safe store into %+v", env.R0)
			}
			env.A = env.R2
		case base.OP_SAFE_LOAD:
			v := base.NewValue()
			switch idx := int(env.R1.AsNumber()); env.R0.Type() {
			case base.Tbytes:
				bb := env.R0.AsBytesUnsafe()
				v = base.NewNumberValue(float64(bb[(idx+len(bb))%len(bb)]))
			case base.Tstring:
				bb := env.R0.AsStringUnsafe()
				v = base.NewNumberValue(float64(bb[(idx+len(bb))%len(bb)]))
			case base.Tlist:
				bb := env.R0.AsListUnsafe()
				v = bb[(idx+len(bb))%len(bb)]
				if v.Type() == base.Tclosure {
					v.AsClosureUnsafe().SetCaller(env.R0)
				}
			case base.Tmap:
				v = env.R0.AsMapUnsafe()[env.R1.AsString()]
				if v.Type() == base.Tclosure {
					v.AsClosureUnsafe().SetCaller(env.R0)
				}
			default:
				log.Panicf("can't safe load from %+v", env.R0)
			}
			env.A = v
		case base.OP_R0:
			env.R0 = env.Get(c.ReadInt32())
		case base.OP_R0_NUM:
			env.R0 = base.NewNumberValue(c.ReadDouble())
		case base.OP_R0_STR:
			env.R0 = base.NewStringValue(c.ReadString())
		case base.OP_R1:
			env.R1 = env.Get(c.ReadInt32())
		case base.OP_R1_NUM:
			env.R1 = base.NewNumberValue(c.ReadDouble())
		case base.OP_R1_STR:
			env.R1 = base.NewStringValue(c.ReadString())
		case base.OP_R2:
			env.R2 = env.Get(c.ReadInt32())
		case base.OP_R2_NUM:
			env.R2 = base.NewNumberValue(c.ReadDouble())
		case base.OP_R2_STR:
			env.R2 = base.NewStringValue(c.ReadString())
		case base.OP_R3:
			env.R3 = env.Get(c.ReadInt32())
		case base.OP_R3_NUM:
			env.R3 = base.NewNumberValue(c.ReadDouble())
		case base.OP_R3_STR:
			env.R3 = base.NewStringValue(c.ReadString())
		case base.OP_PUSH:
			if newEnv == nil {
				newEnv = base.NewEnv(nil)
			}
			newEnv.Push(env.Get(c.ReadInt32()))
		case base.OP_PUSH_NUM:
			if newEnv == nil {
				newEnv = base.NewEnv(nil)
			}
			newEnv.Push(base.NewNumberValue(c.ReadDouble()))
		case base.OP_PUSH_STR:
			if newEnv == nil {
				newEnv = base.NewEnv(nil)
			}
			newEnv.Push(base.NewStringValue(c.ReadString()))
		case base.OP_RET:
			return env.Get(c.ReadInt32())
		case base.OP_RET_NUM:
			return base.NewNumberValue(c.ReadDouble())
		case base.OP_RET_STR:
			return base.NewStringValue(c.ReadString())
		case base.OP_LAMBDA:
			argsCount := int(c.ReadInt32())
			buf := c.Read(int(c.ReadInt32()))
			env.A = base.NewClosureValue(base.NewClosure(buf, env, argsCount))
		case base.OP_CALL:
			v := env.Get(c.ReadInt32())
			if v.Type() == base.Tclosure {
				cls := v.AsClosureUnsafe()
				if newEnv == nil {
					newEnv = base.NewEnv(nil)
				}

				if newEnv.Size() < cls.ArgsCount() {
					if newEnv.Size() == 0 {
						env.A = (base.NewClosureValue(cls))
					} else {
						curry := cls.Dup()
						curry.AppendPreArgs(newEnv.Stack().Values())
						env.A = (base.NewClosureValue(curry))
					}
				} else {
					if cls.PreArgs() != nil && len(cls.PreArgs()) > 0 {
						newEnv.Stack().Insert(0, cls.PreArgs())
					}

					newEnv.SetParent(cls.Env())
					newEnv.C = cls.Caller()
					env.A = Exec(newEnv, cls.Code())
				}

				newEnv = nil
			} else if v.Type() == base.Tlist {
				if bb := v.AsListUnsafe(); newEnv.Size() == 2 {
					env.A = base.NewListValue(bb[shiftIndex(newEnv.Get(0), len(bb)):shiftIndex(newEnv.Get(1), len(bb))])
				} else {
					log.Panicf("too many (or few) arguments to call on list")
				}
				newEnv.Stack().Clear()
			} else if v.Type() == base.Tbytes {
				if bb := v.AsBytesUnsafe(); newEnv.Size() == 2 {
					env.A = base.NewBytesValue(bb[shiftIndex(newEnv.Get(0), len(bb)):shiftIndex(newEnv.Get(1), len(bb))])
				} else {
					log.Panicf("too many (or few) arguments to call on bytes")
				}
				newEnv.Stack().Clear()
			} else if v.Type() == base.Tstring {
				if bb := v.AsStringUnsafe(); newEnv.Size() == 2 {
					env.A = base.NewStringValue(bb[shiftIndex(newEnv.Get(0), len(bb)):shiftIndex(newEnv.Get(1), len(bb))])
				} else {
					log.Panicf("too many (or few) arguments to call on string")
				}
				newEnv.Stack().Clear()
			} else {
				log.Panicf("invalid callee: %+v", v)
			}
		case base.OP_VARARGS:
			list0 := env.Stack().Values()
			list1 := make([]base.Value, len(list0))
			copy(list1, list0)
			env.A = base.NewListValue(list1)
		case base.OP_JMP:
			off := int(c.ReadInt32())
			c.SetCursor(c.GetCursor() + off)
		case base.OP_IF:
			cond := env.Get(c.ReadInt32())
			off := int(c.ReadInt32())
			if cond.IsFalse() {
				c.SetCursor(c.GetCursor() + off)
			}
		case base.OP_TYPEOF:
			v, t := env.Get(c.ReadInt32()), c.ReadInt32()
			env.A = base.NewBoolValue(v.Type() == byte(t))
		case base.OP_LIB_CALL:
			libidx := int(uint32(c.ReadInt32()))
			if libidx >= len(Lib) {
				panic("lib call index overflows")
			}
			env.A = (Lib[libidx].f(env))
		case base.OP_LIB_CALL_EX:
			libidx := int(uint32(c.ReadInt32()))
			if libidx >= len(Lib) {
				panic("lib call index overflows")
			}
			if newEnv == nil {
				newEnv = base.NewEnv(nil)
			}
			newEnv.SetParent(env)
			env.A = (Lib[libidx].ff(newEnv))
			newEnv.Stack().Clear()
		case base.OP_DUP:
			switch env.R0.Type() {
			case base.Tnil, base.Tnumber, base.Tstring, base.Tbool, base.Tclosure, base.Tgeneric:
				env.A = env.R0
			case base.Tlist:
				list0 := env.R0.AsListUnsafe()
				list1 := make([]base.Value, len(list0))
				copy(list1, list0)
				env.A = base.NewListValue(list1)
			case base.Tmap:
				map0 := env.R0.AsMapUnsafe()
				map1 := make(map[string]base.Value)
				for k, v := range map0 {
					map1[k] = v
				}
				env.A = base.NewMapValue(map1)
			case base.Tbytes:
				bytes0 := env.R0.AsBytesUnsafe()
				bytes1 := make([]byte, len(bytes0))
				copy(bytes1, bytes0)
				env.A = base.NewBytesValue(bytes1)
			}
		}
	}

	return base.NewValue()
}

func shiftIndex(index base.Value, len int) int {
	i := int(index.AsNumberUnsafe())
	if i >= 0 {
		return i
	}
	return i + len
}
