package script

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"sync/atomic"
)

type stacktrace struct {
	cursor      uint32
	stackOffset uint32
	cls         *Func
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
		src := uint32(0)
		for i := 0; i < len(r.cls.code.Pos); {
			var opx uint32 = math.MaxUint32
			ii, op, line := r.cls.code.Pos.read(i)
			if ii < len(r.cls.code.Pos)-1 {
				_, opx, _ = r.cls.code.Pos.read(ii)
			}
			if r.cursor >= op && r.cursor < opx {
				src = line
				break
			}
			if r.cursor < op && i == 0 {
				src = line
				break
			}
			i = ii
		}
		// the recorded cursor was advanced by 1 already
		msg.WriteString(fmt.Sprintf("%s at line %d (cursor: %d)\n", r.cls.name, src, r.cursor-1))
	}
	msg.WriteString("root panic:\n")
	msg.WriteString(fmt.Sprintf("%v\n", e.r))
	return msg.String()
}

func returnVararg(env *Env, a Value, b []Value) (Value, []Value) {
	flag := a.Type() == VStack
	if len(b) == 0 {
		if flag {
			u := a._unpackedStack().a
			if len(u) == 0 {
				return Value{}, nil
			}
			return u[0], u[1:]
		}
		return a, nil
	}

	for _, b := range b {
		flag = flag || b.Type() == VStack
	}

	if !flag {
		// both 'a' and 'b' are not (neither containing) unpacked values
		env.checkRemainStackSize(len(b))
		return a, b
	}

	var b2 []Value
	if a.Type() == VStack {
		b2 = append(b2, a._unpackedStack().a...)
	} else {
		b2 = append(b2, a)
	}
	for _, b := range b {
		if b.Type() == VStack {
			b2 = append(b2, b._unpackedStack().a...)
		} else {
			b2 = append(b2, b)
		}
		env.checkRemainStackSize(len(b2))
	}
	if len(b2) == 0 {
		return Value{}, nil
	}
	env.Global.Survey.AdjustedReturns += int64(len(b2))
	return b2[0], b2[1:]
}

// execCursorLoop executes 'K' under 'Env' from the given start 'cursor'
func execCursorLoop(env Env, K *Func, cursor uint32) (result Value, resultV []Value) {
	var stackEnv = env
	var retStack []stacktrace

	stackEnv.stackOffset = uint32(len(*env.stack))

	defer func() {
		if r := recover(); r != nil {
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

	for {
		if env.Global.Deadline != 0 {
			if atomic.LoadInt64(&now) > env.Global.Deadline {
				panicf("timeout")
			}
		}

		v := K.code.Code[cursor]
		cursor++
		bop, opa, opb := splitInst(v)

		switch bop {
		case OpSet:
			env._set(opa, env._get(opb))
		case OpPushV:
			// if opb != 0 {
			// 	env.V = make([]Value, 0, opb)
			// }
			env.V = append(env.V, env._get(opa))
		case OpPopVClear:
			env.V = nil
		case OpPopVAll:
			if opa == 1 { // popv-all-with-a, e.g.: local ... = foo()
				env.A = _unpackedStack(append([]Value{env.A}, env.V...))
			} else { // popv-all, e.g.: local a, ... = foo()
				env.A = _unpackedStack(env.V)
			}
			env.V = nil
		case OpPopV:
			if len(env.V) == 0 {
				env.A = Value{}
			} else {
				env.A, env.V = env.V[0], env.V[1:]
			}
		case OpInc:
			vaf, vai, vaIsInt := env._get(opa).Expect(VNumber).Num()
			vbf, vbi, vbIsInt := env._get(opb).Expect(VNumber).Num()
			if vaIsInt && vbIsInt {
				env.A = Int(vai + vbi)
			} else {
				env.A = Float(vaf + vbf)
			}
			env._set(opa, env.A)
		case OpJSON:
			var buf []byte
			switch opa {
			case 0:
				buf, _ = json.Marshal(stackEnv.StackInterface())
			case 1:
				a := make(map[string]interface{}, stackEnv.Size()/2)
				for i := 0; i < stackEnv.Size(); i += 2 {
					if k := stackEnv.Stack()[i]; !k.IsNil() {
						a[k.String()] = stackEnv.Get(i + 1).Interface()
					}
				}
				buf, _ = json.Marshal(a)
			}
			if opb == 0 {
				env.A = Interface(jsonQuotedString(buf))
			} else {
				env.A = env.NewStringBytes(buf)
			}
			stackEnv.Clear()
		case OpConcat:
			var x string
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == _StrStr {
				x = va._str() + vb._str()
			} else if va.Type() == VString && vb.Type() == VNumber {
				if vbf, vbi, vbIsInt := vb.Num(); vbIsInt {
					x = va._str() + strconv.FormatInt(vbi, 10)
				} else {
					x = va._str() + strconv.FormatFloat(vbf, 'f', 0, 64)
				}
			} else if vb.Type() == VString && va.Type() == VNumber {
				if vaf, vai, vaIsInt := va.Num(); vaIsInt {
					x = strconv.FormatInt(vai, 10) + vb._str()
				} else {
					x = strconv.FormatFloat(vaf, 'f', 0, 64) + vb._str()
				}
			} else {
				va, vb = va.Expect(VString), vb.Expect(VString)
			}
			env.A = env.NewString(x)
		case OpAdd:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == _NumNum {
				vaf, vai, vaIsInt := va.Num()
				vbf, vbi, vbIsInt := vb.Num()
				if vaIsInt && vbIsInt {
					env.A = Int(vai + vbi)
				} else {
					env.A = Float(vaf + vbf)
				}
			} else if va.Type() == VNumber && vb.Type() == VString {
				vaf, vai, vaIsInt := va.Num()
				if vaIsInt {
					vbi, _ := strconv.ParseInt(vb._str(), 0, 64)
					env.A = Int(vai + vbi)
				} else {
					vbf, _ := strconv.ParseFloat(vb._str(), 64)
					env.A = Float(vaf + vbf)
				}
			} else {
				va, vb = va.Expect(VNumber), vb.Expect(VNumber)
			}
		case OpSub:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == _NumNum {
				vaf, vai, vaIsInt := va.Num()
				vbf, vbi, vbIsInt := vb.Num()
				if vaIsInt && vbIsInt {
					env.A = Int(vai - vbi)
				} else {
					env.A = Float(vaf - vbf)
				}
			} else {
				va, vb = va.Expect(VNumber), vb.Expect(VNumber)
			}
		case OpMul:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == _NumNum {
				vaf, vai, vaIsInt := va.Num()
				vbf, vbi, vbIsInt := vb.Num()
				if vaIsInt && vbIsInt {
					env.A = Int(vai * vbi)
				} else {
					env.A = Float(vaf * vbf)
				}
			} else {
				va, vb = va.Expect(VNumber), vb.Expect(VNumber)
			}
		case OpDiv:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == _NumNum {
				vaf, vai, vaIsInt := va.Num()
				vbf, vbi, vbIsInt := vb.Num()
				if vaIsInt && vbIsInt && vai%vbi == 0 {
					env.A = Int(vai / vbi)
				} else {
					env.A = Float(vaf / vbf)
				}
			} else {
				va, vb = va.Expect(VNumber), vb.Expect(VNumber)
			}
		case OpMod:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == _NumNum {
				vaf, vai, vaIsInt := va.Num()
				vbf, vbi, vbIsInt := vb.Num()
				if vaIsInt && vbIsInt {
					env.A = Int(vai % vbi)
				} else {
					env.A = Float(math.Remainder(vaf, vbf))
				}
			} else {
				va, vb = va.Expect(VNumber), vb.Expect(VNumber)
			}
		case OpEq:
			env.A = Bool(env._get(opa).Equal(env._get(opb)))
		case OpNeq:
			env.A = Bool(!env._get(opa).Equal(env._get(opb)))
		case OpLess:
			switch va, vb := env._get(opa), env._get(opb); va.Type() + vb.Type() {
			case _NumNum:
				vaf, vai, vaIsInt := va.Num()
				vbf, vbi, vbIsInt := vb.Num()
				if vaIsInt && vbIsInt {
					env.A = Bool(vai < vbi)
				} else {
					env.A = Bool(vaf < vbf)
				}
			case _StrStr:
				env.A = Bool(va._str() < vb._str())
			default:
				va, vb = va.Expect(VNumber), vb.Expect(VNumber)
			}
		case OpLessEq:
			switch va, vb := env._get(opa), env._get(opb); va.Type() + vb.Type() {
			case _NumNum:
				vaf, vai, vaIsInt := va.Num()
				vbf, vbi, vbIsInt := vb.Num()
				if vaIsInt && vbIsInt {
					env.A = Bool(vai <= vbi)
				} else {
					env.A = Bool(vaf <= vbf)
				}
			case _StrStr:
				env.A = Bool(va._str() <= vb._str())
			default:
				va, vb = va.Expect(VNumber), vb.Expect(VNumber)
			}
		case OpNot:
			env.A = Bool(env._get(opa).IsFalse())
		case OpPow:
			switch va, vb := env._get(opa), env._get(opb); va.Type() + vb.Type() {
			case _NumNum:
				vaf, vai, vaIsInt := va.Num()
				vbf, vbi, vbIsInt := vb.Num()
				if vaIsInt && vbIsInt && vbi >= 1 {
					env.A = Int(ipow(vai, vbi))
				} else {
					env.A = Float(math.Pow(vaf, vbf))
				}
			default:
				va, vb = va.Expect(VNumber), vb.Expect(VNumber)
			}
		case OpLen:
			switch v := env._get(opa); v.Type() {
			case VString:
				env.A = Float(float64(len(v._str())))
			case VStack:
				env.A = Int(int64(len(v._unpackedStack().a)))
			case VFunction:
				env.A = Float(float64(v.Function().numParams))
			default:
				env.A = Int(int64(reflectLen(v.Interface())))
			}
		case OpStore:
			subject, v := env._get(opa), env._get(opb)
			switch subject.Type() {
			case VStack:
				subject._unpackedStack().Put(env.A.ExpectMsg(VNumber, "store").Int(), v)
			case VInterface:
				reflectStore(subject.Interface(), env.A, v)
			default:
				subject = subject.Expect(VStack)
			}
			env.A = v
		case OpSlice:
			subject := env._get(opa)
			start, end := env.A.ExpectMsg(VNumber, "slice").Int(), env._get(opb).ExpectMsg(VNumber, "slice").Int()
			switch subject.Type() {
			case VStack:
				env.A = _unpackedStack(subject._unpackedStack().Slice(start, end))
			case VString:
				s := subject._str()
				start, end := sliceInRange(start, end, len(s))
				env.A = _str(s[start:end])
			case VInterface:
				env.A = Interface(reflectSlice(subject.Interface(), start, end))
			default:
				subject = subject.Expect(VStack)
			}
		case OpLoad:
			switch a := env._get(opa); a.Type() {
			case VStack:
				env.A = a._unpackedStack().Get(env._get(opb).ExpectMsg(VNumber, "load").Int())
			case VInterface:
				env.A = reflectLoad(a.Interface(), env._get(opb))
			case VString:
				if idx, s := env._get(opb).ExpectMsg(VNumber, "load").Int(), a._str(); idx >= 1 && idx <= int64(len(s)) {
					env.A = Int(int64(s[idx-1]))
				}
			default:
				a = a.Expect(VStack)
			}
		case OpPush:
			if v := env._get(opa); v.Type() == VStack {
				*stackEnv.stack = append(*stackEnv.stack, v._unpackedStack().a...)
			} else {
				stackEnv.Push(v)
			}
			if opa == regA && len(env.V) > 0 {
				*stackEnv.stack = append(*stackEnv.stack, env.V...)
				env.V = env.V[:0]
			}
			if env.Global.MaxStackSize > 0 && int64(len(*stackEnv.stack)) > env.Global.MaxStackSize {
				panicf("stack overflow, max: %d", env.Global.MaxStackSize)
			}
		case OpYield:
			env.V = append(env.V, Int(int64(cursor)))
			env.V = append(env.V, _interface(&Env{V: append([]Value{}, env.Stack()...)}))
			env.Global.Survey.YieldSize += int64(env.Size())
			fallthrough
		case OpRet:
			v := env._get(opa)
			if len(retStack) == 0 {
				v, env.V = returnVararg(&env, v, env.V)
				return v, env.V
			}
			// Return2 upper stack
			r := retStack[len(retStack)-1]
			cursor = r.cursor
			K = r.cls
			env.stackOffset = r.stackOffset
			env.A, env.V = returnVararg(&env, v, env.V)
			*env.stack = (*env.stack)[:env.stackOffset+uint32(r.cls.stackSize)]
			stackEnv.stackOffset = uint32(len(*env.stack))
			retStack = retStack[:len(retStack)-1]
		case OpLoadFunc:
			env.A = Function(env.Global.Funcs[opa])
		case OpCallMap:
			cls := env._get(opa).ExpectMsg(VFunction, "callmap").Function()
			m := make(map[string]Value, stackEnv.Size()/2)
			for i := 0; i < stackEnv.Size(); i += 2 {
				var name string
				if a := stackEnv.Stack()[i]; a.Type() == VNumber && a.Int() < int64(len(cls.params)) {
					name = cls.params[a.Int()]
				} else {
					name = a.String()
				}
				if v, ok := m[name]; ok {
					m[name] = v._append(stackEnv.Get(i + 1))
				} else {
					m[name] = stackEnv.Get(i + 1)
				}
			}
			stackEnv.Clear() // TODO: support variadic function?
			for i := byte(0); i < cls.numParams; i++ {
				if int(i) < len(cls.params) {
					stackEnv.Push(m[cls.params[i]])
				} else {
					stackEnv.Push(Value{})
				}
			}
			fallthrough
		case OpCall:
			cls := env._get(opa).ExpectMsg(VFunction, "call").Function()

			if cls.native != nil {
				stackEnv.Global = env.Global
				stackEnv.nativeSource = cls
				if cls.isDebug {
					stackEnv.debug = &debugInfo{
						Caller:     K,
						Cursor:     cursor,
						Stacktrace: append(retStack, stacktrace{cls: K, cursor: cursor}),
					}
					cls.native(&stackEnv)
					stackEnv.debug = nil
				} else {
					cls.native(&stackEnv)
				}
				env.A, env.V = returnVararg(&env, stackEnv.A, stackEnv.V)
				stackEnv.nativeSource = nil
				stackEnv.Clear()
			} else {
				if cls.isVariadic {
					var varg []Value
					if stackEnv.Size() > int(cls.numParams) {
						varg = append([]Value{}, stackEnv.Stack()[cls.numParams:]...)
					}
					stackEnv.growZero(int(cls.stackSize))
					stackEnv._set(uint16(cls.numParams), _unpackedStack(varg))
				} else {
					stackEnv.growZero(int(cls.stackSize))
				}

				last := stacktrace{
					cls:         K,
					cursor:      cursor,
					stackOffset: uint32(env.stackOffset),
				}

				// Switch 'env' to 'stackEnv' and clear 'stackEnv'
				cursor = 0
				K = cls
				env.stackOffset = stackEnv.stackOffset
				env.Global = cls.loadGlobal
				env.V = env.V[:0]

				if opb == 0 {
					retStack = append(retStack, last)
				}

				if env.Global.MaxCallStackSize > 0 && int64(len(retStack)) > env.Global.MaxCallStackSize {
					panicf("call stack overflow, max: %d", env.Global.MaxCallStackSize)
				}

				stackEnv.stackOffset = uint32(len(*env.stack))
			}
		case OpJmp:
			cursor = uint32(int32(cursor) + int32(opb) - 1<<12)
		case OpIfNot:
			if cond := env._get(opa); cond.IsFalse() {
				cursor = uint32(int32(cursor) + int32(opb) - 1<<12)
			}
		case OpIf:
			if cond := env._get(opa); !cond.IsFalse() {
				cursor = uint32(int32(cursor) + int32(opb) - 1<<12)
			}
		}
	}
}
