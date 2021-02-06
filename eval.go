package script

import (
	"bytes"
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
		for i := 0; i < len(r.cls.Code.Pos); {
			var opx uint32 = math.MaxUint32
			ii, op, line := r.cls.Code.Pos.read(i)
			if ii < len(r.cls.Code.Pos)-1 {
				_, opx, _ = r.cls.Code.Pos.read(ii)
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
		msg.WriteString(fmt.Sprintf("%s at line %d (cursor: %d)\n", r.cls.Name, src, r.cursor-1))
	}
	msg.WriteString("root panic:\n")
	msg.WriteString(fmt.Sprintf("%v\n", e.r))
	return msg.String()
}

// InternalExecCursorLoop executes 'K' under 'env' from the given start 'cursor'
func InternalExecCursorLoop(env Env, K *Func, cursor uint32) Value {
	stackEnv := env
	stackEnv.StackOffset = uint32(len(*env.stack))

	var retStack []stacktrace

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

		v := K.Code.Code[cursor]
		cursor++
		bop, opa, opb := splitInst(v)

		switch bop {
		case OpSet:
			env._set(opa, env._get(opb))
		case OpInc:
			vaf, vai, vaIsInt := env._get(opa).MustBe(VNumber, "inc, read sym", 0).Num()
			vbf, vbi, vbIsInt := env._get(opb).MustBe(VNumber, "inc, read step", 0).Num()
			if vaIsInt && vbIsInt {
				env.A = Int(vai + vbi)
			} else {
				env.A = Float(vaf + vbf)
			}
			env._set(opa, env.A)
		case OpConcat:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == VString+VString {
				env.A = concat(env.Global, va._str(), vb._str())
			} else if va.Type() == VString && vb.Type() == VNumber {
				if vbf, vbi, vbIsInt := vb.Num(); vbIsInt {
					env.A = concat(env.Global, va._str(), strconv.FormatInt(vbi, 10))
				} else {
					env.A = concat(env.Global, va._str(), strconv.FormatFloat(vbf, 'f', 0, 64))
				}
			} else if vb.Type() == VString && va.Type() == VNumber {
				if vaf, vai, vaIsInt := va.Num(); vaIsInt {
					env.A = concat(env.Global, strconv.FormatInt(vai, 10), vb._str())
				} else {
					env.A = concat(env.Global, strconv.FormatFloat(vaf, 'f', 0, 64), vb._str())
				}
			} else {
				va.MustBe(VString, "concat", 0)
				vb.MustBe(VString, "concat", 0)
			}
		case OpAdd:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == VNumber+VNumber {
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
				va.MustBe(VNumber, "add", 0)
				vb.MustBe(VNumber, "add", 0)
			}
		case OpSub:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == VNumber+VNumber {
				vaf, vai, vaIsInt := va.Num()
				vbf, vbi, vbIsInt := vb.Num()
				if vaIsInt && vbIsInt {
					env.A = Int(vai - vbi)
				} else {
					env.A = Float(vaf - vbf)
				}
			} else {
				va.MustBe(VNumber, "sub", 0)
				vb.MustBe(VNumber, "sub", 0)
			}
		case OpMul:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == VNumber+VNumber {
				vaf, vai, vaIsInt := va.Num()
				vbf, vbi, vbIsInt := vb.Num()
				if vaIsInt && vbIsInt {
					env.A = Int(vai * vbi)
				} else {
					env.A = Float(vaf * vbf)
				}
			} else {
				va.MustBe(VNumber, "mul", 0)
				vb.MustBe(VNumber, "mul", 0)
			}
		case OpDiv:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == VNumber+VNumber {
				vaf, vai, vaIsInt := va.Num()
				vbf, vbi, vbIsInt := vb.Num()
				if vaIsInt && vbIsInt && vai%vbi == 0 {
					env.A = Int(vai / vbi)
				} else {
					env.A = Float(vaf / vbf)
				}
			} else {
				va.MustBe(VNumber, "div", 0)
				vb.MustBe(VNumber, "div", 0)
			}
		case OpIDiv:
			env.A = Int(env._get(opa).MustBe(VNumber, "idiv", 0).Int() / env._get(opb).MustBe(VNumber, "idiv", 0).Int())
		case OpMod:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == VNumber+VNumber {
				vaf, vai, vaIsInt := va.Num()
				vbf, vbi, vbIsInt := vb.Num()
				if vaIsInt && vbIsInt {
					env.A = Int(vai % vbi)
				} else {
					env.A = Float(math.Remainder(vaf, vbf))
				}
			} else {
				va.MustBe(VNumber, "mod", 0)
				vb.MustBe(VNumber, "mod", 0)
			}
		case OpEq:
			env.A = Bool(env._get(opa).Equal(env._get(opb)))
		case OpNeq:
			env.A = Bool(!env._get(opa).Equal(env._get(opb)))
		case OpLess:
			switch va, vb := env._get(opa), env._get(opb); va.Type() + vb.Type() {
			case VNumber + VNumber:
				vaf, vai, vaIsInt := va.Num()
				vbf, vbi, vbIsInt := vb.Num()
				if vaIsInt && vbIsInt {
					env.A = Bool(vai < vbi)
				} else {
					env.A = Bool(vaf < vbf)
				}
			case VString + VString:
				env.A = Bool(va._str() < vb._str())
			default:
				va.MustBe(VNumber, "less", 0)
				vb.MustBe(VNumber, "less", 0)
			}
		case OpLessEq:
			switch va, vb := env._get(opa), env._get(opb); va.Type() + vb.Type() {
			case VNumber + VNumber:
				vaf, vai, vaIsInt := va.Num()
				vbf, vbi, vbIsInt := vb.Num()
				if vaIsInt && vbIsInt {
					env.A = Bool(vai <= vbi)
				} else {
					env.A = Bool(vaf <= vbf)
				}
			case VString + VString:
				env.A = Bool(va._str() <= vb._str())
			default:
				va.MustBe(VNumber, "less", 0)
				vb.MustBe(VNumber, "less", 0)
			}
		case OpNot:
			env.A = Bool(env._get(opa).IsFalse())
		case OpPow:
			switch va, vb := env._get(opa), env._get(opb); va.Type() + vb.Type() {
			case VNumber + VNumber:
				vaf, vai, vaIsInt := va.Num()
				vbf, vbi, vbIsInt := vb.Num()
				if vaIsInt && vbIsInt && vbi >= 1 {
					env.A = Int(ipow(vai, vbi))
				} else {
					env.A = Float(math.Pow(vaf, vbf))
				}
			default:
				va.MustBe(VNumber, "pow", 0)
				vb.MustBe(VNumber, "pow", 0)
			}
		case OpLen:
			switch v := env._get(opa); v.Type() {
			case VString:
				env.A = Float(float64(len(v._str())))
			case VArray:
				env.A = Int(int64(v.Array().Len()))
			case VFunction:
				env.A = Float(float64(v.Function().NumParams))
			default:
				env.A = Int(int64(reflectLen(v.Interface())))
			}
		case OpList:
			env.A = env.NewArray(append([]Value{}, stackEnv.Stack()...)...)
			stackEnv.Clear()
		case OpGStore:
			env.A = env._get(opb)
			env.Global.GStore(env._get(opa).MustBe(VString, "gstore", 0)._str(), env.A)
		case OpGLoad:
			env.A = env.Global.GLoad(env._get(opa).MustBe(VString, "gload", 0)._str())
		case OpStore:
			subject, v := env._get(opa), env._get(opb)
			switch subject.Type() {
			case VArray:
				if subject.Array().Put1(env.A.MustBe(VNumber, "store", 0).Int(), v) {
					env.Global.DecrDeadsize(ValueSize)
				}
			case VInterface:
				reflectStore(subject.Interface(), env.A, v)
			default:
				subject = subject.MustBe(VArray, "store", 0)
			}
			env.A = v
		case OpSlice:
			subject := env._get(opa)
			start, end := env.A.MustBe(VNumber, "slice", 0).Int(), env._get(opb).MustBe(VNumber, "slice", 0).Int()
			switch subject.Type() {
			case VArray:
				env.A = Array(subject.Array().Slice1(start, end)...)
			case VString:
				s := subject._str()
				start, end := sliceInRange(start, end, len(s))
				env.A = String(s[start:end])
			case VInterface:
				env.A = Interface(reflectSlice(subject.Interface(), start, end))
			default:
				subject = subject.MustBe(VArray, "slice", 0)
			}
		case OpLoad:
			switch a := env._get(opa); a.Type() {
			case VArray:
				env.A = a.Array().Get1(env._get(opb).MustBe(VNumber, "load", 0).Int())
			case VInterface:
				env.A = reflectLoad(a.Interface(), env._get(opb))
			case VString:
				if idx, s := env._get(opb).MustBe(VNumber, "load", 0).Int(), a._str(); idx >= 1 && idx <= int64(len(s)) {
					env.A = Int(int64(s[idx-1]))
				}
			default:
				a = a.MustBe(VArray, "load", 0)
			}
		case OpPush:
			v := env._get(opa)
			if v.Type() == VArray && v.Array().Unpacked {
				*stackEnv.stack = append(*stackEnv.stack, v.Array().Underlay...)
			} else {
				stackEnv.Push(v)
			}
		case OpRet:
			v := env._get(opa)
			if len(retStack) == 0 {
				return v
			}
			// Return upper stack
			r := retStack[len(retStack)-1]
			cursor = r.cursor
			K = r.cls
			env.StackOffset = r.stackOffset
			env.A = v
			*env.stack = (*env.stack)[:env.StackOffset+uint32(r.cls.StackSize)]
			stackEnv.StackOffset = uint32(len(*env.stack))
			retStack = retStack[:len(retStack)-1]
		case OpLoadFunc:
			env.A = Function(env.Global.Functions[opa])
		case OpCallMap:
			cls := env._get(opa).MustBe(VFunction, "callmap", 0).Function()
			m := make(map[string]Value, stackEnv.Size()/2)
			for i := 0; i < stackEnv.Size(); i += 2 {
				a := stackEnv.Stack()[i]
				if a.IsNil() {
					continue
				}
				var name string
				if a.Type() == VNumber && a.Int() < int64(len(cls.Params)) {
					name = cls.Params[a.Int()]
				} else {
					name = a.String()
				}
				if _, ok := m[name]; ok {
					panicf("call: duplicated parameter: %q", name)
				}
				m[name] = stackEnv.Get(i + 1)
			}
			stackEnv.Clear()
			for i := byte(0); i < cls.NumParams; i++ {
				if int(i) < len(cls.Params) {
					stackEnv.Push(m[cls.Params[i]])
				} else {
					stackEnv.Push(Value{})
				}
			}
			stackEnv.A = Interface(m)
			fallthrough
		case OpCall:
			cls := env._get(opa).MustBe(VFunction, "call", 0).Function()

			if cls.Native != nil {
				stackEnv.Global = env.Global
				stackEnv.NativeSource = cls
				if cls.IsDebug {
					stackEnv.Debug = &debugInfo{
						Caller:     K,
						Cursor:     cursor,
						Stacktrace: append(retStack, stacktrace{cls: K, cursor: cursor}),
					}
					cls.Native(&stackEnv)
					stackEnv.Debug = nil
				} else {
					cls.Native(&stackEnv)
				}
				env.A = stackEnv.A
				stackEnv.NativeSource = nil
				stackEnv.Clear()
			} else {
				stackEnv.growZero(int(cls.StackSize))

				last := stacktrace{
					cls:         K,
					cursor:      cursor,
					stackOffset: uint32(env.StackOffset),
				}

				// Switch 'env' to 'stackEnv' and clear 'stackEnv'
				cursor = 0
				K = cls
				env.StackOffset = stackEnv.StackOffset
				env.Global = cls.loadGlobal
				env.A = stackEnv.A

				if opb == 0 {
					retStack = append(retStack, last)
				}

				if env.Global.MaxCallStackSize > 0 && int64(len(retStack)) > env.Global.MaxCallStackSize {
					panicf("call stack overflow, max: %d", env.Global.MaxCallStackSize)
				}

				stackEnv.StackOffset = uint32(len(*env.stack))
			}
		case OpJmp:
			cursor = uint32(int32(cursor) + int32(v&0xffffff) - 1<<23)
		case OpIfNot:
			if env.A.IsFalse() {
				cursor = uint32(int32(cursor) + int32(v&0xffffff) - 1<<23)
			}
		case OpIf:
			if !env.A.IsFalse() {
				cursor = uint32(int32(cursor) + int32(v&0xffffff) - 1<<23)
			}
		}
	}
}

func concat(p *Program, a, b string) Value {
	p.DecrDeadsize(int64(len(a) + len(b)))
	return String(a + b)
}
