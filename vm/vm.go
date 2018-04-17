package vm

import (
	"log"
	"os"

	"github.com/coyove/eugine/base"
)

func extOp(c *base.BytesReader, ext byte, env *base.Env) {
	bop := c.ReadByte()
	first, second := 0.0, 0.0
	switch ext {
	case base.OP_EXT_F_F:
		first, second = env.Get(c.ReadInt32()).Number(), env.Get(c.ReadInt32()).Number()
	case base.OP_EXT_F_IMM:
		first, second = env.Get(c.ReadInt32()).Number(), c.ReadDouble()
	case base.OP_EXT_IMM_F:
		first, second = c.ReadDouble(), env.Get(c.ReadInt32()).Number()
	}

	switch bop {
	case base.OP_SUB:
		env.SetANumber(first - second)
	case base.OP_MUL:
		env.SetANumber(first * second)
	case base.OP_DIV:
		env.SetANumber(first / second)
	case base.OP_MOD:
		env.SetANumber(float64(int64(first) % int64(second)))
	case base.OP_BIT_LSH:
		env.SetANumber(float64(uint64(first) << uint64(second)))
	case base.OP_BIT_RSH:
		env.SetANumber(float64(uint64(first) >> uint64(second)))
	case base.OP_BIT_AND:
		env.SetANumber(float64(int64(first) & int64(second)))
	case base.OP_BIT_OR:
		env.SetANumber(float64(int64(first) | int64(second)))
	case base.OP_LESS:
		env.SetA(base.NewBoolValue(first < second))
	case base.OP_LESS_EQ:
		env.SetA(base.NewBoolValue(first <= second))
	case base.OP_MORE:
		env.SetA(base.NewBoolValue(first > second))
	case base.OP_MORE_EQ:
		env.SetA(base.NewBoolValue(first >= second))
	}
}

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
			env.SetA(base.NewValue())
		case base.OP_TRUE:
			env.SetA(base.NewBoolValue(true))
		case base.OP_FALSE:
			env.SetA(base.NewBoolValue(false))
		case base.OP_BYTES:
			if env.SizeR() == 0 {
				env.SetA(base.NewBytesValue(nil))
			} else if n := env.R0(); n.Type() == base.TY_string {
				env.SetA(base.NewBytesValue([]byte(n.String())))
			} else if n.Type() == base.TY_number {
				env.SetA(base.NewBytesValue(make([]byte, int(n.Number()))))
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
			base.SubI(env)
		case base.OP_MUL:
			base.MulI(env)
		case base.OP_DIV:
			base.DivI(env)
		case base.OP_MOD:
			base.ModI(env)
		case base.OP_INC:
			id := c.ReadInt32()
			v := env.Get(id).Number() + env.Get(c.ReadInt32()).Number()
			env.Set(id, base.NewNumberValue(v))
			env.SetANumber(v)
		case base.OP_INC_NUM:
			id := c.ReadInt32()
			v := env.Get(id).Number() + c.ReadDouble()
			env.Set(id, base.NewNumberValue(v))
			env.SetANumber(v)
		case base.OP_EQ:
			base.LogicCompare(env, env.R(0).Equal, false)
		case base.OP_NEQ:
			base.LogicCompare(env, env.R(0).Equal, true)
		case base.OP_LESS:
			base.LogicCompare(env, env.R(0).Less, false)
		case base.OP_MORE_EQ:
			base.LogicCompare(env, env.R(0).Less, true)
		case base.OP_LESS_EQ:
			base.LogicCompare(env, env.R(0).LessEqual, false)
		case base.OP_MORE:
			base.LogicCompare(env, env.R(0).LessEqual, true)
		case base.OP_ASSERT:
			loc := "assertion failed: " + c.ReadString()
			first := env.R(0)
			if env.SizeR() == 1 {
				if first.IsFalse() {
					panic(loc)
				}
			} else {
				for i := 1; i < env.SizeR(); i++ {
					if !first.Equal(env.R(i)) {
						panic(loc)
					}
				}
			}
			env.SetA(base.NewBoolValue(true))
		case base.OP_LIST:
			list := make([]base.Value, env.SizeR())
			for i := 0; i < env.SizeR(); i++ {
				list[i] = env.R(i)
			}
			env.SetA(base.NewArrayValue(list))
		case base.OP_MAP:
			m := make(map[string]base.Value)
			for i := 0; i < env.SizeR(); i += 2 {
				m[env.R(i).String()] = env.R(i + 1)
			}
			env.SetA(base.NewMapValue(m))
		case base.OP_LEN:
			switch v := env.R(0); v.Type() {
			case base.TY_string:
				env.SetANumber(float64(len(v.String())))
			case base.TY_array:
				env.SetANumber(float64(len(v.Array())))
			case base.TY_map:
				env.SetANumber(float64(len(v.Map())))
			case base.TY_bytes:
				env.SetANumber(float64(len(v.Bytes())))
			default:
				log.Panicf("can't evaluate the length of %v", v)
			}
		case base.OP_STORE:
			obj, v := env.R(0), env.R(env.SizeR()-1)
			for i, sz := 1, env.SizeR(); i < sz-1; i++ {
				switch h := env.R(i); obj.Type() {
				case base.TY_bytes:
					obj.Bytes()[int(h.Number())] = byte(v.Number())
				case base.TY_array:
					idx := int(h.Number())
					if i == sz-2 {
						obj.Array()[idx] = v
					} else {
						obj = (obj.Array())[idx]
					}
				case base.TY_map:
					key := h.String()
					if i == sz-2 {
						obj.Map()[key] = v
					} else {
						obj = obj.Map()[key]
					}
				default:
					log.Panicf("can't store into %v", obj)
				}
			}
			env.SetA(v)
		case base.OP_LOAD:
			obj := env.R0()
			for i, sz := 1, env.SizeR(); i < sz; i++ {
				switch h := env.R(i); obj.Type() {
				case base.TY_bytes:
					obj = base.NewNumberValue(float64(obj.Bytes()[int(h.Number())]))
				case base.TY_string:
					obj = base.NewNumberValue(float64(obj.String()[int(h.Number())]))
				case base.TY_array:
					obj = obj.Array()[int(h.Number())]
				case base.TY_map:
					obj = obj.Map()[h.String()]
				default:
					log.Panicf("can't load from %v", obj)
				}
			}
			env.SetA(obj)
		case base.OP_PUSHF:
			env.PushR(env.Get(c.ReadInt32()))
		case base.OP_PUSHF_NUM:
			env.PushR(base.NewNumberValue(c.ReadDouble()))
		case base.OP_PUSHF_STR:
			env.PushR(base.NewStringValue(c.ReadString()))
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
			env.SetA(base.NewClosureValue(base.NewClosure(buf, env, argsCount)))
		case base.OP_CALL:
			cls := env.Get(c.ReadInt32()).Closure()
			if newEnv == nil {
				newEnv = base.NewEnv(nil)
			}

			if newEnv.Size() < cls.ArgsCount() {
				if newEnv.Size() == 0 {
					env.SetA(base.NewClosureValue(cls))
				} else {
					curry := cls.Dup()
					curry.AppendPreArgs(newEnv.Stack().Values())
					env.SetA(base.NewClosureValue(curry))
				}
			} else {
				if cls.PreArgs() != nil && len(cls.PreArgs()) > 0 {
					newEnv.Stack().Insert(0, cls.PreArgs())
				}

				newEnv.SetParent(cls.Env())
				env.SetA(Exec(newEnv, cls.Code()))
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
			base.NotI(env)
		case base.OP_AND:
			base.AndI(env)
		case base.OP_OR:
			base.OrI(env)
		case base.OP_XOR:
			base.XorI(env)
		case base.OP_BIT_NOT:
			base.BitNotI(env)
		case base.OP_BIT_AND:
			base.BitAndI(env)
		case base.OP_BIT_OR:
			base.BitOrI(env)
		case base.OP_BIT_XOR:
			base.BitXorI(env)
		case base.OP_BIT_LSH:
			base.BitLshI(env)
		case base.OP_BIT_RSH:
			base.BitRshI(env)
		case base.OP_TYPEOF:
			v, t := env.Get(c.ReadInt32()), c.ReadInt32()
			env.SetA(base.NewBoolValue(v.Type() == byte(t)))
		case base.OP_EXT_F_F, base.OP_EXT_F_IMM, base.OP_EXT_IMM_F:
			extOp(c, bop, env)
		case base.OP_LIB_CALL:
			libidx := int(uint32(c.ReadInt32()))
			if libidx >= len(Lib) {
				panic("lib call index overflows")
			}
			env.SetA(Lib[libidx].f(env))
		case base.OP_LIB_CALL_EX:
			libidx := int(uint32(c.ReadInt32()))
			if libidx >= len(Lib) {
				panic("lib call index overflows")
			}
			if newEnv == nil {
				panic("shouldn't happen")
			}
			newEnv.SetParent(env)
			env.SetA(Lib[libidx].ff(newEnv))
			newEnv.Stack().Clear()
		}
	}

	return base.NewValue()
}
