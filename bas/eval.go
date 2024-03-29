package bas

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"unicode/utf8"
	"unsafe"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

type Stacktrace struct {
	Cursor      uint32
	StackOffset uint32
	Callable    *Object
}

func (r *Stacktrace) sourceLine() (src uint32) {
	posv := r.Callable.fun.codeSeg.Pos
	lastLine := uint32(math.MaxUint32)
	cursor := r.Cursor - 1
	if posv.Len() > 0 {
		_, op, line := posv.Read(0)
		for cursor > op && posv.Len() > 0 {
			op, line = posv.Pop()
			// fmt.Println(r.Callable.fun.name, cursor, op, line)
		}
		if cursor <= op {
			return line
		}
		lastLine = line
	}
	return lastLine
}

// ExecError represents the runtime error/panic.
type ExecError struct {
	fatal  bool
	root   interface{}
	stacks []Stacktrace
}

func (e *ExecError) IsPanic() bool {
	return e.fatal
}

func (e *ExecError) GetCause() interface{} {
	if e == nil {
		return nil
	}
	return e.root
}

func (e *ExecError) Error() string {
	msg := bytes.Buffer{}
	e.printError(&msg, 0)
	return msg.String()
}

func (e *ExecError) printError(msg *bytes.Buffer, idx int) {
	if ee, ok := e.root.(*ExecError); ok {
		ee.printError(msg, idx+1)
	} else if ee, ok := e.root.(Value); ok && ee.IsError() {
		ee.Error().printError(msg, idx+1)
	} else {
		fmt.Fprint(msg, e.root)
	}
	msg.WriteByte('\n')

	if e.fatal {
		msg.WriteString("fatal ")
	}

	if idx == 0 {
		msg.WriteString("stacktrace:\n")
	} else {
		msg.WriteString(fmt.Sprintf("stacktrace %d:\n", idx))
	}

	if len(e.stacks) == 0 {
		msg.WriteString("none\n")
	} else {
		for i := len(e.stacks) - 1; i >= 0; i-- {
			r := e.stacks[i]
			if r.Callable.fun.native != nil {
				msg.WriteString(fmt.Sprintf("%s (native function)\n\t<native code>\n", r.Callable.fun.name))
			} else {
				ln := r.sourceLine()
				msg.WriteString(fmt.Sprintf("%s at %s:%d (i%d)",
					r.Callable.fun.name,
					r.Callable.fun.codeSeg.Pos.Name,
					ln,
					r.Cursor-1, // the recorded cursor was advanced by 1 already
				))
				msg.WriteString("\n\t")
				line, ok := internal.LineOf(r.Callable.fun.top.Source, int(ln))
				if ok {
					msg.WriteString(strings.TrimSpace(line))
				} else {
					msg.WriteString("<unknown source>")
				}
				msg.WriteString("\n")
			}
		}
	}
	msg.Truncate(msg.Len() - 1)
}

func relayPanic(onPanic func() []Stacktrace) {
	if r := recover(); r != nil {
		if re, ok := r.(*ExecError); ok {
			panic(re)
		}

		e := &ExecError{}
		e.fatal = true
		e.root = r
		e.stacks = append([]Stacktrace{}, onPanic()...)
		panic(e)
	}
}

func internalExecCursorLoop(env Env, K *Object, retStack []Stacktrace) Value {
	stackEnv := env
	stackEnv.stackOffset = uint32(len(*env.stack))

	var cursor uint32
	retStackStartSize := len(retStack)

	defer relayPanic(func() []Stacktrace {
		retStack = append(retStack, Stacktrace{
			Cursor:      cursor,
			Callable:    K,
			StackOffset: env.stackOffset,
		})
		if stackEnv.runtime.stack0.Callable != nil {
			retStack = append(retStack, stackEnv.runtime.stack0)
		}
		return retStack
	})

	code := K.fun.codeSeg.CodePtr()
	for {
		v := (*typ.Inst)(unsafe.Pointer(code + typ.InstSize*uintptr(cursor)))
		opa, opb := v.A, v.B
		cursor++

		switch v.Opcode {
		case typ.OpSet:
			env._set(opa, env._get(opb))
		case typ.OpInc:
			va := env._ref(opa)
			if vb := env._get(opb); va.IsInt64() && vb.IsInt64() {
				*va = Int64(va.UnsafeInt64() + vb.UnsafeInt64())
			} else if va.IsNumber() && vb.IsNumber() {
				*va = Float64(va.Float64() + vb.Float64())
			} else if vat := va.Type(); vat == typ.String && vb.Type() == typ.String {
				*va = Str(va.Str() + vb.Str())
			} else if vat == typ.Native {
				va.Native().Append(vb)
			} else {
				internal.Panic("inc "+errNeedNumbersOrStrings, va.simple(), vb.simple())
			}
			env.A = *va
			cursor = uint32(int32(cursor) + int32(int16(v.C)))
		case typ.OpNext:
			va, vb := env._get(opa), env._get(opb)
			switch va.Type() {
			case typ.Nil:
				env.A = Nil
			case typ.Native:
				env.A = va.Native().internalNext(vb)
			case typ.Object:
				env.A = va.Object().internalNext(vb)
			case typ.String:
				idx := 0
				if vb != Nil {
					idx = vb.Native().Get(0).Int()
				} else {
					vb = Array(Nil, Nil)
				}
				if r, sz := utf8.DecodeRuneInString(va.Str()[idx:]); sz == 0 {
					vb.Native().Set(0, Nil)
					vb.Native().Set(1, Nil)
				} else {
					vb.Native().Set(0, Int(idx+sz))
					vb.Native().Set(1, Int(int(r)))
				}
				env.A = vb
			default:
				internal.Panic("can't iterate over %v", va.simple())
			}
		case typ.OpLen:
			env.A = Int(env._get(opa).Len())
		case typ.OpExt:
			switch v.OpcodeExt {
			case typ.OpExtInc16:
				va := env._ref(opa)
				if va.IsInt64() {
					*va = Int64(va.UnsafeInt64() + int64(int16(opb)))
				} else if va.IsNumber() {
					*va = Float64(va.Float64() + float64(int16(opb)))
				} else if va.Type() == typ.Native {
					va.Native().Append(Int(int(int16(opb))))
				} else {
					internal.Panic("inc "+errNeedNumbers, va.simple(), int16(opb))
				}
				env.A = *va
				cursor = uint32(int32(cursor) + int32(int16(v.C)))
			case typ.OpExtLess16:
				if va := env._ref(opa); va.IsInt64() {
					env.A = Bool(va.UnsafeInt64() < int64(int16(opb)))
				} else if va.IsNumber() {
					env.A = Bool(va.UnsafeFloat64() < float64(int16(opb)))
				} else {
					internal.Panic("comparison "+errNeedNumbers, va.simple(), int16(opb))
				}
			case typ.OpExtAdd16:
				if va := env._ref(opa); va.IsInt64() {
					env.A = Int64(va.UnsafeInt64() + int64(int16(opb)))
				} else if va.IsNumber() {
					env.A = Float64(va.UnsafeFloat64() + float64(int16(opb)))
				} else if va.Type() == typ.Native {
					va.Native().Append(Int(int(int16(opb))))
					env.A = *va
				} else {
					internal.Panic("add "+errNeedNumbers, va.simple(), int16(opb))
				}
			case typ.OpExtRSub16:
				if va := env._ref(opa); va.IsInt64() {
					env.A = Int64(int64(int16(opb)) - va.UnsafeInt64())
				} else if va.IsNumber() {
					env.A = Float64(float64(int16(opb)) - va.UnsafeFloat64())
				} else {
					internal.Panic("sub "+errNeedNumbers, va.simple(), int16(opb))
				}
			case typ.OpExtGreat16:
				if va := env._ref(opa); va.IsInt64() {
					env.A = Bool(va.UnsafeInt64() > int64(int16(opb)))
				} else if va.IsNumber() {
					env.A = Bool(va.UnsafeFloat64() > float64(int16(opb)))
				} else {
					internal.Panic("comparison "+errNeedNumbers, va.simple(), int16(opb))
				}
			case typ.OpExtNeq16:
				env.A = Bool(env._get(opa) != Int64(int64(int16(opb))))
			case typ.OpExtEq16:
				env.A = Bool(env._get(opa) == Int64(int64(int16(opb))))
			case typ.OpExtStore16:
				subject, v := env._ref(opa), env._get(v.C)
				if st := subject.Type(); st == typ.Native {
					subject.Native().Set(int(int16(opb)), v)
				} else if st == typ.Object {
					subject.Object().Set(Int(int(int16(opb))), v)
				} else {
					internal.Panic("can't index %v using %v", subject.simple(), int16(opb))
				}
				env.A = v
			case typ.OpExtLoad16:
				a := env._ref(opa)
				if at := a.Type(); at == typ.Native {
					env.A = a.Native().Get(int(int16(opb)))
				} else if at == typ.String {
					env.A = Int64(int64(a.Str()[int16(opb)]))
				} else if at == typ.Object {
					env.A = a.Object().Get(Int(int(int16(opb))))
				} else {
					env.A = Nil
				}
				if v.C != typ.RegA {
					env._set(v.C, env.A)
				}
			case typ.OpExtBitAnd:
				va, vb := env._ref(opa), env._ref(opb)
				if va.Type() == typ.Native && vb.Type() == typ.Native {
					va.Native().Concat(vb.Native())
					env.A = *va
				} else {
					bassertTwoInts("and", va, vb)
					env.A = Int64(va.Int64() & vb.Int64())
				}
			case typ.OpExtBitOr:
				va, vb := env._ref(opa), env._ref(opb)
				bassertTwoInts("or", va, vb)
				env.A = Int64(va.Int64() | vb.Int64())
			case typ.OpExtBitXor:
				va, vb := env._ref(opa), env._ref(opb)
				bassertTwoInts("xor", va, vb)
				env.A = Int64(va.Int64() ^ vb.Int64())
			case typ.OpExtBitLsh:
				va, vb := env._ref(opa), env._ref(opb)
				bassertTwoInts("lsh", va, vb)
				env.A = Int64(va.Int64() << vb.Int64())
			case typ.OpExtBitRsh:
				va, vb := env._ref(opa), env._ref(opb)
				bassertTwoInts("rsh", va, vb)
				env.A = Int64(va.Int64() >> vb.Int64())
			case typ.OpExtBitURsh:
				va, vb := env._ref(opa), env._ref(opb)
				bassertTwoInts("ursh", va, vb)
				env.A = Int64(int64(uint64(va.Int64()) >> vb.Int64()))
			case typ.OpExtBitAnd16:
				va := env._get(opa)
				bassertOneInt("and", va)
				env.A = Int64(va.Int64() & int64(int16(opb)))
			case typ.OpExtBitOr16:
				va := env._get(opa)
				bassertOneInt("or", va)
				env.A = Int64(va.Int64() | int64(int16(opb)))
			case typ.OpExtBitXor16:
				va := env._get(opa)
				bassertOneInt("xor", va)
				env.A = Int64(va.Int64() ^ int64(int16(opb)))
			case typ.OpExtBitLsh16:
				va := env._get(opa)
				bassertOneInt("lsh", va)
				env.A = Int64(va.Int64() << int64(int16(opb)))
			case typ.OpExtBitRsh16:
				va := env._get(opa)
				bassertOneInt("rsh", va)
				env.A = Int64(va.Int64() >> int64(int16(opb)))
			case typ.OpExtBitURsh16:
				va := env._get(opa)
				bassertOneInt("uRsh", va)
				env.A = Int64(int64(uint64(va.Int64()) >> int64(int16(opb))))
			}
		case typ.OpAdd:
			if va, vb := env._ref(opa), env._ref(opb); va.IsInt64() && vb.IsInt64() {
				env.A = Int64(va.UnsafeInt64() + vb.UnsafeInt64())
			} else if va.IsNumber() && vb.IsNumber() {
				env.A = Float64(va.Float64() + vb.Float64())
			} else if vat, vbt := va.Type(), vb.Type(); vat == typ.String && vbt == typ.String {
				env.A = Str(va.Str() + vb.Str())
			} else if vat == typ.Native {
				va.Native().Append(*vb)
				env.A = *va
			} else {
				internal.Panic("add "+errNeedNumbersOrStrings, va.simple(), vb.simple())
			}
		case typ.OpSub:
			if va, vb := env._ref(opa), env._ref(opb); va.IsInt64() && vb.IsInt64() {
				env.A = Int64(va.UnsafeInt64() - vb.UnsafeInt64())
			} else if va.IsNumber() && vb.IsNumber() {
				env.A = Float64(va.Float64() - vb.Float64())
			} else if va.IsObject() {
				env.A = va.Object().Delete(*vb)
			} else {
				internal.Panic("sub "+errNeedNumbers, va.simple(), vb.simple())
			}
		case typ.OpMul:
			if va, vb := env._ref(opa), env._ref(opb); va.IsInt64() && vb.IsInt64() {
				env.A = Int64(va.UnsafeInt64() * vb.UnsafeInt64())
			} else if va.IsNumber() && vb.IsNumber() {
				env.A = Float64(va.Float64() * vb.Float64())
			} else {
				internal.Panic("mul "+errNeedNumbers, va.simple(), vb.simple())
			}
		case typ.OpDiv:
			if va, vb := env._ref(opa), env._ref(opb); va.IsNumber() && vb.IsNumber() {
				env.A = Float64(va.Float64() / vb.Float64())
			} else {
				internal.Panic("div "+errNeedNumbers, va.simple(), vb.simple())
			}
		case typ.OpIDiv:
			if va, vb := env._ref(opa), env._ref(opb); va.IsNumber() && vb.IsNumber() {
				env.A = Int64(va.Int64() / vb.Int64())
			} else {
				internal.Panic("idiv "+errNeedNumbers, va.simple(), vb.simple())
			}
		case typ.OpMod:
			if va, vb := env._get(opa), env._get(opb); va.IsNumber() && vb.IsNumber() {
				env.A = Int64(va.Int64() % vb.Int64())
			} else {
				internal.Panic("mod "+errNeedNumbers, va.simple(), vb.simple())
			}
		case typ.OpEq:
			env.A = Bool(env._ref(opa).Equal(env._get(opb)))
		case typ.OpNeq:
			env.A = Bool(!env._ref(opa).Equal(env._get(opb)))
		case typ.OpLess:
			if va, vb := env._ref(opa), env._ref(opb); va.IsInt64() && vb.IsInt64() {
				env.A = Bool(va.UnsafeInt64() < vb.UnsafeInt64())
			} else if va.IsNumber() && vb.IsNumber() {
				env.A = Bool(va.Float64() < vb.Float64())
			} else if va.Type() == typ.String && vb.Type() == typ.String {
				env.A = Bool(lessStr(*va, *vb))
			} else {
				internal.Panic("comparison "+errNeedNumbersOrStrings, va.simple(), vb.simple())
			}
		case typ.OpLessEq:
			if va, vb := env._ref(opa), env._ref(opb); va.IsInt64() && vb.IsInt64() {
				env.A = Bool(va.UnsafeInt64() <= vb.UnsafeInt64())
			} else if va.IsNumber() && vb.IsNumber() {
				env.A = Bool(va.Float64() <= vb.Float64())
			} else if va.Type() == typ.String && vb.Type() == typ.String {
				env.A = Bool(!lessStr(*vb, *va))
			} else {
				internal.Panic("comparison "+errNeedNumbersOrStrings, va.simple(), vb.simple())
			}
		case typ.OpNot:
			env.A = Bool(env._get(opa).IsFalse())
		case typ.OpCreateArray:
			env.A = newArray(append([]Value{}, stackEnv.Stack()...)...).ToValue()
			stackEnv.clear()
		case typ.OpCreateObject:
			stk := stackEnv.Stack()
			o := NewObject(len(stk) / 2)
			for i := 0; i < len(stk); i += 2 {
				o.Set(stk[i], stk[i+1])
			}
			env.A = o.ToValue()
			stackEnv.clear()
		case typ.OpIsProto:
			if a, b := env._get(opa), env._get(opb); b.IsString() {
				env.A = Bool(TestShapeFast(a, b.Str()) == nil)
			} else if b.IsObject() {
				env.A = Bool(a.HasPrototype(b.Object()))
			} else {
				env.A = Bool(a.Equal(b))
			}
		case typ.OpStore:
			subject, k, v := env._ref(opa), env._get(opb), env._get(v.C)
			switch subject.Type() {
			case typ.Object:
				subject.Object().Set(k, v)
			case typ.Native:
				if k.IsInt64() {
					subject.Native().Set(k.Int(), v)
				} else {
					subject.Native().SetKey(k, v)
				}
			default:
				internal.Panic("can't index %v using %v", subject.simple(), k.simple())
			}
			env.A = v
		case typ.OpLoad:
			switch a, idx := env._ref(opa), env._get(opb); a.Type() {
			case typ.Object:
				env.A = a.Object().Get(idx)
			case typ.Native:
				if idx.IsInt64() {
					env.A = a.Native().Get(idx.Int())
				} else {
					env.A, _ = a.Native().GetKey(idx)
				}
			case typ.String:
				if idx.IsInt64() {
					env.A = Int64(int64(a.Str()[idx.UnsafeInt64()]))
				} else {
					env.A = setObjectRecv(Proto.Str.Get(idx), *a)
				}
			default:
				env.A = Nil
			}
			env._set(v.C, env.A)
		case typ.OpSlice:
			a, start, end := env._get(opa), env._get(opb), env._get(v.C)
			switch a.Type() {
			case typ.Native:
				if a := a.Native(); a.HasPrototype(&Proto.NativeMap) {
					if v, ok := a.GetKey(start); ok {
						env.A = v
					} else {
						env.A = end
					}
				} else {
					if !start.IsInt64() || !end.IsInt64() {
						internal.Panic("can't slice %v using %v and %v", a.ToValue().simple(), start.simple(), end.simple())
					}
					if end := end.Int(); end == -1 {
						env.A = a.Slice(start.Int(), a.Len()).ToValue()
					} else {
						env.A = a.Slice(start.Int(), end).ToValue()
					}
				}
			case typ.String:
				if !start.IsInt64() || !end.IsInt64() {
					internal.Panic("can't slice %v using %v and %v", a.simple(), start.simple(), end.simple())
				}
				if end := end.Int(); end == -1 {
					env.A = Str(a.Str()[start.Int():a.Len()])
				} else {
					env.A = Str(a.Str()[start.Int():end])
				}
			case typ.Object:
				env.A = a.Object().GetDefault(start, end)
			default:
				internal.Panic("can't slice %v", a.simple())
			}
		case typ.OpPush:
			stackEnv.push(env._get(opa))
		case typ.OpPushUnpack:
			switch a := env._get(opa); a.Type() {
			case typ.Native:
				*stackEnv.stack = append(*stackEnv.stack, a.Native().Values()...)
			case typ.Nil:
			default:
				internal.Panic("arguments unpacking expects array, got %v", a.simple())
			}
		case typ.OpRet:
			v := env._get(opa)
			if len(retStack) == retStackStartSize {
				return v
			}
			// Return to upper stack
			r := retStack[len(retStack)-1]
			cursor = r.Cursor
			K = r.Callable
			code = K.fun.codeSeg.CodePtr()
			env.stackOffset = r.StackOffset
			env.A = v
			env.top = K.fun.top
			*env.stack = (*env.stack)[:env.stackOffset+uint32(r.Callable.fun.stackSize)]
			stackEnv.stackOffset = uint32(len(*env.stack))
			retStack = retStack[:len(retStack)-1]
		case typ.OpFunction:
			if opa == typ.RegA {
				env.A = K.ToValue()
			} else {
				o := env._get(opa).Object().Copy()
				if opb == 1 {
					o.Merge(K)
					for addr, name := range o.fun.caps {
						if name == "" {
							continue
						}
						if uint16(addr) == v.C {
							// Recursive closure, e.g.:
							// function foo()
							//   function bar()
							//     self.bar()
							//   end
							//   return bar
							// end
							o.Set(Str(name), o.ToValue())
						} else {
							o.Set(Str(name), env._get(uint16(addr)))
						}
					}
				}
				env._set(v.C, o.ToValue())
			}
		case typ.OpCall, typ.OpTailCall:
			a := env._refgp(opa)
			if a.Type() != typ.Object {
				internal.Panic("can't call %v", a.simple())
			}
			obj := a.Object()
			cls := obj.fun
			if v.OpcodeExt == 1 {
				stackEnv.push(env._get(opb))
			} else if v.OpcodeExt == 2 {
				stackEnv.push(env._get(opb))
				stackEnv.push(env._get(v.C))
			}
			stackEnv.A = obj.this
			if cls.varg {
				s, w := stackEnv.Stack(), int(cls.numArgs)-1
				if len(s) > w {
					s[w] = newVarargArray(s[w:]).ToValue()
				} else {
					if len(s) < w {
						panicNotEnoughArgs(obj)
					}
					stackEnv.resize(w + 1)
					stackEnv._set(uint16(w), Nil)
				}
			}
			env.checkStackOverflow()

			last := Stacktrace{
				Callable:    K,
				Cursor:      cursor,
				StackOffset: env.stackOffset,
			}

			if cls.native != nil {
				stackEnv.top = env.top
				stackEnv.runtime.stack0 = Stacktrace{Callable: obj}
				stackEnv.runtime.stack1 = last
				stackEnv.runtime.stackN = retStack
				cls.native(&stackEnv)
				cursor = stackEnv.runtime.stack1.Cursor
				env.A = stackEnv.A
				stackEnv.runtime = stacktraces{}
				stackEnv.clear()
			} else if v.Opcode == typ.OpCall {
				if stackEnv.Size() < int(cls.numArgs) {
					panicNotEnoughArgs(obj)
				}
				// Switch 'env' to 'stackEnv' and move up 'stackEnv'.
				stackEnv.resizeZero(int(cls.stackSize), int(cls.numArgs))
				cursor = 0
				K = obj
				code = obj.fun.codeSeg.CodePtr()
				env.stackOffset = stackEnv.stackOffset
				env.top = cls.top
				env.A = stackEnv.A

				retStack = append(retStack, last)
				stackEnv.stackOffset = uint32(len(*env.stack))
			} else {
				if stackEnv.Size() < int(cls.numArgs) {
					panicNotEnoughArgs(obj)
				}
				// Move arguments from 'stackEnv' to 'env'.
				*env.stack = append((*env.stack)[:env.stackOffset], stackEnv.Stack()...)

				// Resize 'env' to allocate enough space for the next function and move up 'stackEnv'.
				env.resizeZero(int(cls.stackSize), int(cls.numArgs))
				cursor = 0
				K = obj
				code = obj.fun.codeSeg.CodePtr()
				env.top = cls.top
				env.A = stackEnv.A
				stackEnv.stackOffset = uint32(len(*env.stack))
			}
		case typ.OpJmp:
			cursor = uint32(int32(cursor) + v.D())
		case typ.OpJmpFalse:
			if env.A.IsFalse() {
				cursor = uint32(int32(cursor) + v.D())
			}
		case typ.OpLoadTop:
			if a := topSymbols.stack[opa]; opb != typ.RegPhantom {
				env._set(v.C, a.AssertObject("load top").Get(env._get(opb)))
			} else {
				env._set(v.C, a)
			}
		}
	}
}
