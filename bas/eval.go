package bas

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

type Stacktrace struct {
	Cursor          uint32
	stackOffsetFlag uint32 // LSB: 1=tailcall
	Callable        *Object
}

func (r *Stacktrace) StackOffset() uint32 {
	return r.stackOffsetFlag & internal.MaxStackSize
}

func (r *Stacktrace) IsTailcall() bool {
	return r.stackOffsetFlag&internal.FlagTailCall > 0
}

func (r *Stacktrace) IsNativeCall() bool {
	return r.stackOffsetFlag&internal.FlagNativeCall > 0
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

// ExecError represents the runtime error
type ExecError struct {
	root   interface{}
	stacks []Stacktrace
}

func (e *ExecError) GetCause() interface{} {
	if e == nil {
		return nil
	}
	return e.root
}

func (e *ExecError) Error() string {
	msg := bytes.Buffer{}
	msg.WriteString(fmt.Sprintf("%v\n", e.root))
	msg.WriteString("stacktrace:\n")
	for i := len(e.stacks) - 1; i >= 0; i-- {
		r := e.stacks[i]
		if r.IsNativeCall() {
			msg.WriteString(fmt.Sprintf("%s (native function)\n\t<native code>\n", r.Callable.fun.name))
		} else {
			ln := r.sourceLine()
			msg.WriteString(fmt.Sprintf("%s at %s:%d (i%d)",
				r.Callable.fun.name,
				r.Callable.fun.codeSeg.Pos.Name,
				ln,
				r.Cursor-1, // the recorded cursor was advanced by 1 already
			))
			if r.IsTailcall() {
				msg.WriteString(" (tailcall)")
			}
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
	return msg.String()
}

func relayPanic(onPanic func() []Stacktrace) {
	if r := recover(); r != nil {
		if re, ok := r.(*ExecError); ok {
			panic(re)
		}

		e := &ExecError{}
		e.root = r
		e.stacks = append([]Stacktrace{}, onPanic()...)
		panic(e)
	}
}

func internalExecCursorLoop(env Env, K *Object, retStack []Stacktrace) Value {
	stackEnv := env
	stackEnv.stackOffsetFlag = uint32(len(*env.stack))

	var cursor uint32
	retStackStartSize := len(retStack)

	defer relayPanic(func() []Stacktrace {
		retStack = append(retStack, Stacktrace{
			Cursor:          cursor,
			Callable:        K,
			stackOffsetFlag: env.stackOffsetFlag,
		})
		if stackEnv.runtime.stack0.Callable != nil {
			retStack = append(retStack, stackEnv.runtime.stack0)
		}
		return retStack
	})

	code := K.fun.codeSeg.Code
	for {
		v := code[cursor]
		opa, opb := v.A, v.B
		cursor++

		switch v.Opcode {
		case typ.OpSet:
			env._set(opa, env._get(opb))
		case typ.OpInc:
			va := env._getRef(opa)
			if vb := env._get(opb); va.IsInt64() && vb.IsInt64() {
				*va = Int64(va.UnsafeInt64() + vb.UnsafeInt64())
			} else if va.IsNumber() && vb.IsNumber() {
				*va = Float64(va.Float64() + vb.Float64())
			} else if va.Type() == typ.String && vb.Type() == typ.String {
				*va = Str(va.Str() + vb.Str())
			} else {
				internal.Panic("inc "+errNeedNumbersOrStrings, detail(*va), detail(vb))
			}
			env.A = *va
			cursor = uint32(int32(cursor) + int32(int16(v.C)))
		case typ.OpInc16:
			va := env._getRef(opa)
			if va.IsInt64() {
				*va = Int64(va.UnsafeInt64() + int64(int16(opb)))
			} else if va.IsNumber() {
				*va = Float64(va.Float64() + float64(int16(opb)))
			} else {
				internal.Panic("inc16 "+errNeedNumber, detail(*va))
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
				internal.Panic("can't iterate over %v", detail(va))
			}
		case typ.OpLen:
			env.A = Int(env._get(opa).Len())
		case typ.OpLinear16:
			if va := env._get(opa); va.IsInt64() {
				env.A = Int64(va.UnsafeInt64()*int64(int16(opb)) + int64(int16(v.C)))
			} else if va.IsNumber() {
				env.A = Float64(va.UnsafeFloat64()*float64(int16(opb)) + float64(int16(v.C)))
			} else {
				internal.Panic("arithmetic "+errNeedNumber, detail(va))
			}
		case typ.OpCmp16:
			if va := env._get(opa); va.IsInt64() {
				env.A = Bool(va.UnsafeInt64()*int64(int16(opb)) < int64(int16(v.C)))
			} else if va.IsNumber() {
				env.A = Bool(va.UnsafeFloat64()*float64(int16(opb)) < float64(int16(v.C)))
			} else {
				internal.Panic("comparison "+errNeedNumber, detail(va))
			}
		case typ.OpEq16:
			if va := env._get(opa); va.IsInt64() {
				env.A = Bool((va.UnsafeInt64() == int64(int16(opb))) == (v.C == typ.OpEq))
			} else if va.IsNumber() {
				env.A = Bool((va.UnsafeFloat64() == float64(int16(opb))) == (v.C == typ.OpEq))
			} else {
				internal.Panic("equality "+errNeedNumber, detail(va))
			}
		case typ.OpAdd:
			if va, vb := env._get(opa), env._get(opb); va.IsInt64() && vb.IsInt64() {
				env.A = Int64(va.UnsafeInt64() + vb.UnsafeInt64())
			} else if va.IsNumber() && vb.IsNumber() {
				env.A = Float64(va.Float64() + vb.Float64())
			} else if x := va.Type() + vb.Type(); x == typ.String*2 {
				env.A = Str(va.Str() + vb.Str())
			} else {
				internal.Panic("add "+errNeedNumbersOrStrings, detail(va), detail(vb))
			}
		case typ.OpSub:
			if va, vb := env._get(opa), env._get(opb); va.IsInt64() && vb.IsInt64() {
				env.A = Int64(va.UnsafeInt64() - vb.UnsafeInt64())
			} else if va.IsNumber() && vb.IsNumber() {
				env.A = Float64(va.Float64() - vb.Float64())
			} else if va.IsObject() {
				env.A = va.Object().Delete(vb)
			} else {
				internal.Panic("sub "+errNeedNumbers, detail(va), detail(vb))
			}
		case typ.OpMul:
			if va, vb := env._get(opa), env._get(opb); va.IsInt64() && vb.IsInt64() {
				env.A = Int64(va.UnsafeInt64() * vb.UnsafeInt64())
			} else if va.IsNumber() && vb.IsNumber() {
				env.A = Float64(va.Float64() * vb.Float64())
			} else {
				internal.Panic("mul "+errNeedNumbers, detail(va), detail(vb))
			}
		case typ.OpDiv:
			if va, vb := env._get(opa), env._get(opb); va.IsNumber() && vb.IsNumber() {
				env.A = Float64(va.Float64() / vb.Float64())
			} else {
				internal.Panic("div "+errNeedNumbers, detail(va), detail(vb))
			}
		case typ.OpIDiv:
			if va, vb := env._get(opa), env._get(opb); va.IsNumber() && vb.IsNumber() {
				env.A = Int64(va.Int64() / vb.Int64())
			} else {
				internal.Panic("idiv "+errNeedNumbers, detail(va), detail(vb))
			}
		case typ.OpMod:
			if va, vb := env._get(opa), env._get(opb); va.IsNumber() && vb.IsNumber() {
				env.A = Int64(va.Int64() % vb.Int64())
			} else {
				internal.Panic("mod "+errNeedNumbers, detail(va), detail(vb))
			}
		case typ.OpEq:
			env.A = Bool(env._get(opa).Equal(env._get(opb)))
		case typ.OpNeq:
			env.A = Bool(!env._get(opa).Equal(env._get(opb)))
		case typ.OpLess:
			if va, vb := env._get(opa), env._get(opb); va.IsInt64() && vb.IsInt64() {
				env.A = Bool(va.UnsafeInt64() < vb.UnsafeInt64())
			} else if va.IsNumber() && vb.IsNumber() {
				env.A = Bool(va.Float64() < vb.Float64())
			} else if va.Type() == typ.String && vb.Type() == typ.String {
				env.A = Bool(lessStr(va, vb))
			} else {
				internal.Panic("comparison "+errNeedNumbersOrStrings, detail(va), detail(vb))
			}
		case typ.OpLessEq:
			if va, vb := env._get(opa), env._get(opb); va.IsInt64() && vb.IsInt64() {
				env.A = Bool(va.UnsafeInt64() <= vb.UnsafeInt64())
			} else if va.IsNumber() && vb.IsNumber() {
				env.A = Bool(va.Float64() <= vb.Float64())
			} else if va.Type() == typ.String && vb.Type() == typ.String {
				env.A = Bool(!lessStr(vb, va))
			} else {
				internal.Panic("comparison "+errNeedNumbersOrStrings, detail(va), detail(vb))
			}
		case typ.OpNot:
			env.A = Bool(env._get(opa).IsFalse())
		case typ.OpBitOp:
			va, vb := env._get(opa), env._get(opb)
			if !va.IsInt64() || !vb.IsInt64() {
				internal.Panic("bitwise operation requires integer numbers, got %v and %v", detail(va), detail(vb))
			}
			switch v.C {
			case 0:
				env.A = Int64(va.Int64() & vb.Int64())
			case 1:
				env.A = Int64(va.Int64() | vb.Int64())
			case 2:
				env.A = Int64(va.Int64() ^ vb.Int64())
			case 3:
				env.A = Int64(va.Int64() << vb.Int64())
			case 4:
				env.A = Int64(va.Int64() >> vb.Int64())
			case 5:
				env.A = Int64(int64(uint64(va.Int64()) >> vb.Int64()))
			}
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
			if a, b := env._get(opa), env._get(opb); a.Equal(b) {
				env.A = True
			} else if b.IsString() {
				env.A = Bool(TestShapeFast(a, b.Str()) == nil)
			} else {
				env.A = Bool(a.HasPrototype(b.AssertObject("isprototype")))
			}
		case typ.OpStore:
			subject, k, v := env._get(opa), env._get(opb), env._get(v.C)
			switch subject.Type() {
			case typ.Object:
				subject.Object().Set(k, v)
			case typ.Native:
				if k.IsInt64() {
					if a, idx := subject.Native(), k.Int(); idx == a.Len() {
						a.Append(v)
					} else {
						a.Set(idx, v)
					}
				} else {
					subject.Native().SetKey(k, v)
				}
			default:
				internal.Panic("invalid store: %v, key: %v", detail(subject), detail(k))
			}
			env.A = v
		case typ.OpLoad:
			switch a, idx := env._get(opa), env._get(opb); a.Type() {
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
					env.A = Nil
					if s := a.Str(); idx.UnsafeInt64() >= 0 && idx.UnsafeInt64() < int64(len(s)) {
						env.A = Int64(int64(s[idx.UnsafeInt64()]))
					}
				} else {
					env.A = setObjectRecv(Proto.Str.Get(idx), a)
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
						internal.Panic("slice "+errNeedNumbers, detail(start), detail(end))
					}
					if end := end.Int(); end == -1 {
						env.A = a.Slice(start.Int(), a.Len()).ToValue()
					} else {
						env.A = a.Slice(start.Int(), end).ToValue()
					}
				}
			case typ.String:
				if !start.IsInt64() || !end.IsInt64() {
					internal.Panic("slice "+errNeedNumbers, detail(start), detail(end))
				}
				if end := end.Int(); end == -1 {
					env.A = Str(a.Str()[start.Int():a.Len()])
				} else {
					env.A = Str(a.Str()[start.Int():end])
				}
			case typ.Object:
				env.A = a.Object().GetDefault(start, end)
			default:
				internal.Panic("can't slice %v", detail(a))
			}
		case typ.OpPush:
			stackEnv.push(env._get(opa))
		case typ.OpPushUnpack:
			switch a := env._get(opa); a.Type() {
			case typ.Native:
				*stackEnv.stack = append(*stackEnv.stack, a.Native().Values()...)
			case typ.Nil:
			default:
				internal.Panic("arguments unpacking expects array, got %v", detail(a))
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
			code = K.fun.codeSeg.Code
			env.stackOffsetFlag = r.stackOffsetFlag
			env.A = v
			env.top = K.fun.top
			*env.stack = (*env.stack)[:env.stackOffset()+uint32(r.Callable.fun.stackSize)]
			stackEnv.stackOffsetFlag = uint32(len(*env.stack))
			retStack = retStack[:len(retStack)-1]
		case typ.OpFunction:
			if opa == typ.RegA {
				env.A = K.ToValue()
			} else {
				o := env._get(opa).Object().Copy(true)
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
			a := env._get(opa)
			if a.Type() != typ.Object {
				internal.Panic("can't call %v", detail(a))
			}
			obj := a.Object()
			cls := obj.fun
			if opb != typ.RegPhantom {
				stackEnv.push(env._get(opb))
			}
			stackEnv.A = obj.this
			if cls.varg {
				s, w := stackEnv.Stack(), int(cls.numArgs)-1
				if len(s) > w {
					s[w] = newVarargArray(s[w:]).ToValue()
				} else {
					if len(s) < w {
						internal.PanicNotEnoughArgs(detail(a))
					}
					stackEnv.resize(w + 1)
					stackEnv._set(uint16(w), Nil)
				}
			}
			env.checkStackOverflow()

			last := Stacktrace{
				Callable:        K,
				Cursor:          cursor,
				stackOffsetFlag: env.stackOffsetFlag,
			}

			if cls.native != nil {
				stackEnv.top = env.top
				stackEnv.runtime.stack0 = Stacktrace{Callable: obj, stackOffsetFlag: internal.FlagNativeCall}
				stackEnv.runtime.stack1 = last
				stackEnv.runtime.stackN = retStack
				cls.native(&stackEnv)
				stackEnv.runtime = stacktraces{}
				env.A = stackEnv.A
				stackEnv.clear()
			} else if v.Opcode == typ.OpCall {
				if stackEnv.Size() < int(cls.numArgs) {
					internal.PanicNotEnoughArgs(detail(a))
				}
				// Switch 'env' to 'stackEnv' and move up 'stackEnv'.
				stackEnv.resizeZero(int(cls.stackSize), int(cls.numArgs))
				cursor = 0
				K = obj
				code = obj.fun.codeSeg.Code
				env.stackOffsetFlag = stackEnv.stackOffsetFlag
				env.top = cls.top
				env.A = stackEnv.A

				retStack = append(retStack, last)
				stackEnv.stackOffsetFlag = uint32(len(*env.stack))
			} else {
				if stackEnv.Size() < int(cls.numArgs) {
					internal.PanicNotEnoughArgs(detail(a))
				}
				// Move arguments from 'stackEnv' to 'env'.
				*env.stack = append((*env.stack)[:env.stackOffset()], stackEnv.Stack()...)

				// Resize 'env' to allocate enough space for the next function and move up 'stackEnv'.
				env.resizeZero(int(cls.stackSize), int(cls.numArgs))
				cursor = 0
				K = obj
				code = obj.fun.codeSeg.Code
				env.top = cls.top
				env.stackOffsetFlag |= internal.FlagTailCall
				env.A = stackEnv.A
				stackEnv.stackOffsetFlag = uint32(len(*env.stack))
			}
		case typ.OpJmp:
			cursor = uint32(int32(cursor) + v.D())
		case typ.OpIfNot:
			if env.A.IsFalse() {
				cursor = uint32(int32(cursor) + v.D())
			}
		case typ.OpLoadGlobal:
			if a := globals.stack[opa]; opb != typ.RegPhantom {
				env._set(v.C, a.AssertObject("load global").Get(env._get(opb)))
			} else {
				env._set(v.C, a)
			}
		}
	}
}
