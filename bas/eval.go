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
	Callable        *Function
}

func (r *Stacktrace) StackOffset() uint32 {
	return r.stackOffsetFlag & 0x7fffffff
}

func (r *Stacktrace) IsTailcall() bool {
	return r.stackOffsetFlag>>31 == 1
}

func (r *Stacktrace) sourceLine() (src uint32) {
	posv := r.Callable.CodeSeg.Pos
	if posv.Len() > 0 {
		_, op, line := posv.Read(0)
		for r.Cursor > op && posv.Len() > 0 {
			op, line = posv.Pop()
		}
		if r.Cursor <= op {
			return line
		}
	}
	return math.MaxUint32
}

// ExecError represents the runtime error
type ExecError struct {
	root   interface{}
	native *Function
	stacks []Stacktrace
}

func (e *ExecError) TransparentError() internal.TransparentError {
	panic(nil)
}

func (e *ExecError) GetCause() error {
	if e == nil {
		return nil
	}
	if err, ok := e.root.(error); ok {
		return err
	}
	return e
}

func (e *ExecError) Error() string {
	msg := bytes.Buffer{}
	if e.root != nil {
		if e.native != nil {
			msg.WriteString(e.native.Name)
			msg.WriteString("(): ")
		}
		msg.WriteString(fmt.Sprintf("%v\n", e.root))
	}
	msg.WriteString("stacktrace:\n")
	for i := len(e.stacks) - 1; i >= 0; i-- {
		r := e.stacks[i]
		if r.Cursor == internal.NativeCallCursor {
			msg.WriteString(fmt.Sprintf("%s (native)\n", r.Callable.Name))
		} else {
			ln := r.sourceLine()
			msg.WriteString(fmt.Sprintf("%s at %s:%d (i%d)",
				r.Callable.Name,
				r.Callable.CodeSeg.Pos.Name,
				ln,
				r.Cursor-1, // the recorded cursor was advanced by 1 already
			))
			if r.IsTailcall() {
				msg.WriteString(" (tailcall)")
			}
			msg.WriteString("\n\t")
			line, ok := lineOf(r.Callable.LoadGlobal.Source, int(ln))
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

func internalExecCursorLoop(env Env, K *Function, retStack []Stacktrace) Value {
	stackEnv := env
	stackEnv.stackOffset = uint32(len(*env.stack))

	var cursor uint32
	retStackStartSize := len(retStack)

	defer func() {
		if r := recover(); r != nil {
			rr := Stacktrace{
				Cursor:   cursor,
				Callable: K,
			}

			if re, ok := r.(*ExecError); ok {
				retStack = append(retStack, rr)
				re.stacks = append(retStack, re.stacks...)[retStackStartSize:]
				panic(re)
			} else {
				e := &ExecError{}
				e.root = r // root panic
				e.native = stackEnv.runtime.Callable0
				e.stacks = make([]Stacktrace, len(retStack)-retStackStartSize+1)
				copy(e.stacks, retStack[retStackStartSize:])
				e.stacks[len(e.stacks)-1] = rr
				panic(e)
			}
		}
	}()

	for {
		v := K.CodeSeg.Code[cursor]
		bop, opa, opb, opc := v.Opcode, v.A, v.B, v.C
		cursor++

		switch bop {
		case typ.OpSet:
			env._set(opa, env._get(opb))
		case typ.OpInc:
			va, vb := env._get(opa), env._get(opb)
			switch va.Type() + vb.Type() {
			case typ.Number + typ.Number:
				if va.IsInt64() && vb.IsInt64() {
					env.A = Int64(va.UnsafeInt64() + vb.UnsafeInt64())
				} else {
					env.A = Float64(va.Float64() + vb.Float64())
				}
			case typ.String + typ.String:
				env.A = Str(va.Str() + vb.Str())
			default:
				internal.Panic("inc "+errNeedNumbersOrStrings, simpleString(va), simpleString(vb))
			}
			env._set(opa, env.A)
			cursor = uint32(int32(cursor) + int32(int16(opc)))
		case typ.OpNext:
			va, vb := env._get(opa), env._get(opb)
			switch va.Type() {
			case typ.Nil:
				env.A = Array(Nil, Nil)
			case typ.Native:
				env.A = va.Native().Next(vb)
			case typ.Object:
				env.A = va.Object().Next(vb)
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
					vb.Native().Set(1, Rune(r))
				}
				env.A = vb
			default:
				internal.Panic("can't iterate over %v", simpleString(va))
			}
		case typ.OpLen:
			env.A = Int(Len(env._get(opa)))
		case typ.OpAdd:
			va, vb := env._get(opa), env._get(opb)
			switch va.Type() + vb.Type() {
			case typ.Number + typ.Number:
				if sum := va.UnsafeAddr() + vb.UnsafeAddr(); sum == int64Marker2 {
					env.A = Int64(va.UnsafeInt64() + vb.UnsafeInt64())
				} else {
					env.A = Float64(va.Float64() + vb.Float64())
				}
			case typ.String + typ.String:
				env.A = Str(va.Str() + vb.Str())
			case typ.String + typ.Number:
				env.A = Str(va.String() + vb.String())
			default:
				internal.Panic("add "+errNeedNumbersOrStrings, simpleString(va), simpleString(vb))
			}
		case typ.OpSub:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				if sum := va.UnsafeAddr() + vb.UnsafeAddr(); sum == int64Marker2 {
					env.A = Int64(va.UnsafeInt64() - vb.UnsafeInt64())
				} else {
					env.A = Float64(va.Float64() - vb.Float64())
				}
			} else {
				internal.Panic("sub "+errNeedNumbers, simpleString(va), simpleString(vb))
			}
		case typ.OpMul:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				if sum := va.UnsafeAddr() + vb.UnsafeAddr(); sum == int64Marker2 {
					env.A = Int64(va.UnsafeInt64() * vb.UnsafeInt64())
				} else {
					env.A = Float64(va.Float64() * vb.Float64())
				}
			} else {
				internal.Panic("mul "+errNeedNumbers, simpleString(va), simpleString(vb))
			}
		case typ.OpDiv:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Float64(va.Float64() / vb.Float64())
			} else {
				internal.Panic("div "+errNeedNumbers, simpleString(va), simpleString(vb))
			}
		case typ.OpIDiv:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int64(va.Int64() / vb.Int64())
			} else {
				internal.Panic("idiv "+errNeedNumbers, simpleString(va), simpleString(vb))
			}
		case typ.OpMod:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int64(va.Int64() % vb.Int64())
			} else {
				internal.Panic("mod "+errNeedNumbers, simpleString(va), simpleString(vb))
			}
		case typ.OpEq:
			env.A = Bool(env._get(opa).Equal(env._get(opb)))
		case typ.OpNeq:
			env.A = Bool(!env._get(opa).Equal(env._get(opb)))
		case typ.OpLess:
			switch va, vb := env._get(opa), env._get(opb); va.Type() + vb.Type() {
			case typ.Number + typ.Number:
				if sum := va.UnsafeAddr() + vb.UnsafeAddr(); sum == int64Marker2 {
					env.A = Bool(va.UnsafeInt64() < vb.UnsafeInt64())
				} else {
					env.A = Bool(va.Float64() < vb.Float64())
				}
			case typ.String + typ.String:
				env.A = Bool(lessStr(va, vb))
			default:
				internal.Panic("comparison "+errNeedNumbersOrStrings, simpleString(va), simpleString(vb))
			}
		case typ.OpLessEq:
			switch va, vb := env._get(opa), env._get(opb); va.Type() + vb.Type() {
			case typ.Number + typ.Number:
				if sum := va.UnsafeAddr() + vb.UnsafeAddr(); sum == int64Marker2 {
					env.A = Bool(va.UnsafeInt64() <= vb.UnsafeInt64())
				} else {
					env.A = Bool(va.Float64() <= vb.Float64())
				}
			case typ.String + typ.String:
				env.A = Bool(!lessStr(vb, va))
			default:
				internal.Panic("comparison "+errNeedNumbersOrStrings, simpleString(va), simpleString(vb))
			}
		case typ.OpNot:
			env.A = Bool(env._get(opa).IsFalse())
		case typ.OpBitAnd:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int64(va.Int64() & vb.Int64())
			} else {
				internal.Panic("bitwise and "+errNeedNumbers, simpleString(va), simpleString(vb))
			}
		case typ.OpBitOr:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int64(va.Int64() | vb.Int64())
			} else {
				internal.Panic("bitwise or "+errNeedNumbers, simpleString(va), simpleString(vb))
			}
		case typ.OpBitXor:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int64(va.Int64() ^ vb.Int64())
			} else {
				internal.Panic("bitwise xor "+errNeedNumbers, simpleString(va), simpleString(vb))
			}
		case typ.OpBitLsh:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int64(va.Int64() << vb.Int64())
			} else {
				internal.Panic("bitwise lsh "+errNeedNumbers, simpleString(va), simpleString(vb))
			}
		case typ.OpBitRsh:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int64(va.Int64() >> vb.Int64())
			} else {
				internal.Panic("bitwise rsh "+errNeedNumbers, simpleString(va), simpleString(vb))
			}
		case typ.OpBitURsh:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int64(int64(uint64(va.Int64()) >> vb.Int64()))
			} else {
				internal.Panic("bitwise ursh "+errNeedNumbers, simpleString(va), simpleString(vb))
			}
		case typ.OpBitNot:
			if a := env._get(opa); a.Type() == typ.Number {
				env.A = Int64(^a.Int64())
			} else {
				internal.Panic("bitwise not "+errNeedNumber, simpleString(a))
			}
		case typ.OpCreateArray:
			env.A = newArray(append([]Value{}, stackEnv.Stack()...)...).ToValue()
			stackEnv.Clear()
		case typ.OpCreateObject:
			stk := stackEnv.Stack()
			o := NewObject(len(stk))
			for i := 0; i < len(stk); i += 2 {
				o.Set(stk[i], stk[i+1])
			}
			env.A = o.ToValue()
			stackEnv.Clear()
		case typ.OpIsProto:
			if a, b := env._get(opa), env._get(opb); a.Equal(b) {
				env.A = True
			} else {
				env.A = Bool(HasPrototype(a, b.AssertType(typ.Object, "isprototype").Object()))
			}
		case typ.OpStore:
			subject, k, v := env._get(opa), env._get(opb), env._get(opc)
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
				internal.Panic("invalid store: %v, key: %v", simpleString(subject), simpleString(k))
			}
			env.A = v
		case typ.OpLoad:
			switch a, idx := env._get(opa), env._get(opb); a.Type() {
			case typ.Nil, typ.Number, typ.Bool:
				env.A = Nil
			case typ.Object:
				env.A = a.Object().Find(idx)
			case typ.Native:
				if idx.IsInt64() {
					env.A = a.Native().Get(idx.Int())
				} else {
					env.A = a.Native().GetKey(idx)
				}
			case typ.String:
				if idx.IsInt64() {
					env.A = Nil
					if s := a.Str(); idx.UnsafeInt64() >= 0 && idx.UnsafeInt64() < int64(len(s)) {
						env.A = Int64(int64(s[idx.UnsafeInt64()]))
					}
				} else {
					env.A = setObjectRecv(Proto.Str.Find(idx), a)
				}
			default:
				internal.Panic("invalid load: %v, key: %v", simpleString(a), simpleString(idx))
			}
			env._set(opc, env.A)
		case typ.OpSlice:
			a, start, end := env._get(opa), env._get(opb), env._get(opc)
			if start.Type()+end.Type() != typ.Number+typ.Number {
				internal.Panic("slice "+errNeedNumbers, simpleString(start), simpleString(end))
			}
			switch a.Type() {
			case typ.Native:
				if end := end.Int(); end == -1 {
					env.A = a.Native().Slice(start.Int(), a.Native().Len()).ToValue()
				} else {
					env.A = a.Native().Slice(start.Int(), end).ToValue()
				}
			case typ.String:
				if end := end.Int(); end == -1 {
					env.A = Str(a.Str()[start.Int():Len(a)])
				} else {
					env.A = Str(a.Str()[start.Int():end])
				}
			default:
				internal.Panic("can't slice %v", simpleString(a))
			}
		case typ.OpPush:
			stackEnv.Push(env._get(opa))
		case typ.OpPushUnpack:
			switch a := env._get(opa); a.Type() {
			case typ.Native:
				*stackEnv.stack = append(*stackEnv.stack, a.Native().Values()...)
			case typ.Nil:
			default:
				internal.Panic("arguments unpacking expects array, got %v", simpleString(a))
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
			env.stackOffset = r.StackOffset()
			env.A = v
			*env.stack = (*env.stack)[:env.stackOffset+uint32(r.Callable.StackSize)]
			stackEnv.stackOffset = uint32(len(*env.stack))
			retStack = retStack[:len(retStack)-1]
		case typ.OpLoadFunc:
			env.A = env.Global.functions[opa].ToValue()
		case typ.OpCall, typ.OpTailCall:
			a := env._get(opa)
			if a.Type() != typ.Object {
				internal.Panic("can't call %v", simpleString(a))
			}
			cls := a.Object().fun
			if cls == nil {
				internal.Panic("%v not callable", simpleString(a))
			}
			if opb != typ.RegPhantom {
				stackEnv.Push(env._get(opb))
			}
			stackEnv.A = a.Object().this
			if cls.Variadic {
				s, w := stackEnv.Stack(), int(cls.NumParams)-1
				if len(s) > w {
					s[w] = newArray(append([]Value{}, s[w:]...)...).ToValue()
				} else {
					stackEnv.grow(w + 1)
					stackEnv._set(uint16(w), newArray().ToValue())
				}
			}
			if cls.Native != nil {
				stackEnv.Global = env.Global
				stackEnv.runtime.Callable0 = cls
				stackEnv.runtime.Stack1 = Stacktrace{Callable: K, Cursor: cursor}
				stackEnv.runtime.StackN = retStack
				cls.Native(&stackEnv)
				stackEnv.runtime = Runtime{}
				env.A = stackEnv.A
				stackEnv.Clear()
			} else {
				stackEnv.growZero(int(cls.StackSize), int(cls.NumParams))

				last := Stacktrace{
					Callable:        K,
					Cursor:          cursor,
					stackOffsetFlag: uint32(env.stackOffset),
				}

				// Switch 'env' to 'stackEnv' and clear 'stackEnv'
				cursor = 0
				K = cls
				env.stackOffset = stackEnv.stackOffset
				env.Global = cls.LoadGlobal
				env.A = stackEnv.A

				if bop == typ.OpCall {
					retStack = append(retStack, last)
				} else if len(retStack) > 0 {
					retStack[len(retStack)-1].stackOffsetFlag |= 0x80000000
				}

				stackEnv.stackOffset = uint32(len(*env.stack))
			}
		case typ.OpJmp:
			cursor = uint32(int32(cursor) + v.D())
		case typ.OpIfNot:
			if env.A.IsFalse() {
				cursor = uint32(int32(cursor) + v.D())
			}
		}
	}
}
