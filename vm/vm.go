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
		case base.OP_NIL:
			env.A = base.NewValue()
		case base.OP_TRUE:
			env.A = base.NewBoolValue(true)
		case base.OP_FALSE:
			env.A = base.NewBoolValue(false)
		case base.OP_BYTES:
			if n := env.R0; n.Type() == base.TY_string {
				env.A = base.NewBytesValue([]byte(n.String()))
			} else if n.Type() == base.TY_number {
				env.A = base.NewBytesValue(make([]byte, int(n.Number())))
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
			base.AddI(env)
		case base.OP_SUB:
			env.A = base.NewNumberValue(env.R0.Number() - env.R1.Number())
		case base.OP_MUL:
			env.A = base.NewNumberValue(env.R0.Number() * env.R1.Number())
		case base.OP_DIV:
			env.A = base.NewNumberValue(env.R0.Number() / env.R1.Number())
		case base.OP_MOD:
			env.A = base.NewNumberValue(float64(int64(env.R0.Number()) % int64(env.R1.Number())))
		case base.OP_INC:
			id := c.ReadInt32()
			v := env.Get(id).Number() + env.Get(c.ReadInt32()).Number()
			env.Set(id, base.NewNumberValue(v))
			env.A = base.NewNumberValue(v)
		case base.OP_INC_NUM:
			id := c.ReadInt32()
			v := env.Get(id).Number() + c.ReadDouble()
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
		case base.OP_ASSERT:
			loc := "assertion failed: " + c.ReadString()
			if env.R0.IsFalse() {
				panic(loc)
			}
			env.A = base.NewBoolValue(true)
		case base.OP_LEN:
			switch v := env.R0; v.Type() {
			case base.TY_string:
				env.A = base.NewNumberValue(float64(len(v.String())))
			case base.TY_array:
				env.A = base.NewNumberValue(float64(len(v.Array())))
			case base.TY_map:
				env.A = base.NewNumberValue(float64(len(v.Map())))
			case base.TY_bytes:
				env.A = base.NewNumberValue(float64(len(v.Bytes())))
			default:
				log.Panicf("can't evaluate the length of %v", v)
			}
		case base.OP_STORE:
			switch env.R0.Type() {
			case base.TY_bytes:
				env.R0.Bytes()[int(env.R1.Number())] = byte(env.R2.Number())
			case base.TY_array:
				env.R0.Array()[int(env.R1.Number())] = env.R2
			case base.TY_map:
				env.R0.Map()[env.R1.String()] = env.R2
			default:
				log.Panicf("can't store into %v", env.R0)
			}
			env.A = env.R2
		case base.OP_LOAD:
			var v base.Value
			switch env.R0.Type() {
			case base.TY_bytes:
				v = base.NewNumberValue(float64(env.R0.Bytes()[int(env.R1.Number())]))
			case base.TY_string:
				v = base.NewNumberValue(float64(env.R0.String()[int(env.R1.Number())]))
			case base.TY_array:
				v = env.R0.Array()[int(env.R1.Number())]
			case base.TY_map:
				v = env.R0.Map()[env.R2.String()]
			default:
				log.Panicf("can't load from %v", env.R0)
			}
			env.A = v
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
			cls := env.Get(c.ReadInt32()).Closure()
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
				env.A = (Exec(newEnv, cls.Code()))
			}

			newEnv = nil

		case base.OP_JMP:
			off := int(c.ReadInt32())
			c.SetCursor(c.GetCursor() + off)
		case base.OP_IF:
			cond := env.Get(c.ReadInt32())
			off := int(c.ReadInt32())
			if cond.IsFalse() {
				c.SetCursor(c.GetCursor() + off)
			}
		case base.OP_NOT:
			env.A = base.NewBoolValue(!env.R0.Bool())
		case base.OP_AND:
			env.A = base.NewBoolValue(!env.R0.IsFalse() && !env.R1.IsFalse())
		case base.OP_OR:
			env.A = base.NewBoolValue(!env.R0.IsFalse() || !env.R1.IsFalse())
		case base.OP_XOR:
			env.A = base.NewBoolValue(env.R0.IsFalse() != env.R1.IsFalse())
		case base.OP_BIT_NOT:
			env.A = base.NewNumberValue(float64(^int64(env.R0.Number())))
		case base.OP_BIT_AND:
			env.A = base.NewNumberValue(float64(int64(env.R0.Number()) & int64(env.R1.Number())))
		case base.OP_BIT_OR:
			env.A = base.NewNumberValue(float64(int64(env.R0.Number()) | int64(env.R1.Number())))
		case base.OP_BIT_XOR:
			env.A = base.NewNumberValue(float64(int64(env.R0.Number()) ^ int64(env.R1.Number())))
		case base.OP_BIT_LSH:
			env.A = base.NewNumberValue(float64(uint64(env.R0.Number()) << uint64(env.R1.Number())))
		case base.OP_BIT_RSH:
			env.A = base.NewNumberValue(float64(uint64(env.R0.Number()) >> uint64(env.R1.Number())))
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
				panic("shouldn't happen")
			}
			newEnv.SetParent(env)
			env.A = (Lib[libidx].ff(newEnv))
			newEnv.Stack().Clear()
		}
	}

	return base.NewValue()
}
