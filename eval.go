package nj

import (
	"bytes"
	"fmt"
	"math"
	"unicode/utf8"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

type Stacktrace struct {
	Cursor      uint32
	StackOffset uint32
	Callable    *FuncBody
}

// ExecError represents the runtime error
type ExecError struct {
	root   interface{}
	native *FuncBody
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
	msg.WriteString("stacktrace:\n")
	for i := len(e.stacks) - 1; i >= 0; i-- {
		r := e.stacks[i]
		src := uint32(0)
		for i := 0; i < len(r.Callable.Code.Pos); {
			var opx uint32 = math.MaxUint32
			ii, op, line := r.Callable.Code.Pos.read(i)
			if ii < len(r.Callable.Code.Pos)-1 {
				_, opx, _ = r.Callable.Code.Pos.read(ii)
			}
			if r.Cursor >= op && r.Cursor < opx {
				src = line
				break
			}
			if r.Cursor < op && i == 0 {
				src = line
				break
			}
			i = ii
		}
		// the recorded cursor was advanced by 1 already
		msg.WriteString(fmt.Sprintf("%s at line %d (cursor: %d)\n", r.Callable.Name, src, r.Cursor-1))
	}
	if e.root != nil {
		msg.WriteString("root panic:\n")
		if e.native != nil {
			msg.WriteString(e.native.Name)
			msg.WriteString("() ")
		}
		msg.WriteString(fmt.Sprintf("%v\n", e.root))
	}
	return msg.String()
}

func internalExecCursorLoop(env Env, K *FuncBody, retStack []Stacktrace) Value {
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
				e.native = stackEnv.Native
				e.stacks = make([]Stacktrace, len(retStack)-retStackStartSize+1)
				copy(e.stacks, retStack[retStackStartSize:])
				e.stacks[len(e.stacks)-1] = rr
				panic(e)
			}
		}
	}()

	for {
		v := K.Code.Code[cursor]
		bop, opa, opb := v.op, v.a, uint16(v.b)
		cursor++

		switch bop {
		case typ.OpSet:
			env._set(opa, env._get(opb))
		case typ.OpInc:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				if va.IsInt64() && vb.IsInt64() {
					env.A = Int64(va.unsafeInt() + vb.unsafeInt())
				} else {
					env.A = Float64(va.Float64() + vb.Float64())
				}
				env._set(opa, env.A)
			} else {
				internal.Panic("inc "+errNeedNumbers, showType(va), showType(vb))
			}
		case typ.OpNext:
			va, vb := env._get(opa), env._get(opb)
			switch va.Type() {
			case typ.Nil:
				env.A = NewArray(Nil, Nil).ToValue()
			case typ.Array:
				idx := 0
				if vb != Nil {
					idx = vb.Is(typ.Number, "array iteration").Int() + 1
				}
				a := va.Array()
				_ = idx >= a.Len() && env.SetA(NewArray(Nil, Nil).ToValue()) || env.SetA(NewArray(Int(idx), a.Get(idx)).ToValue())
			case typ.Object:
				env.A = NewArray(va.Object().Next(vb)).ToValue()
			case typ.String:
				idx := int64(0)
				if vb != Nil {
					idx = vb.Is(typ.Number, "string iteration").Int64()
				}
				if r, sz := utf8.DecodeRuneInString(va.Str()[idx:]); sz == 0 {
					env.A = NewArray(Nil, Nil).ToValue()
				} else {
					env.A = NewArray(Int64(int64(sz)+idx), Rune(r)).ToValue()
				}
			default:
				internal.Panic("can't iterate %v using %v", showType(va), showType(vb))
			}
		case typ.OpLen:
			env.A = Int(env._get(opa).Len())
		case typ.OpAdd:
			va, vb := env._get(opa), env._get(opb)
			switch va.Type() + vb.Type() {
			case typ.Number + typ.Number:
				if sum := va.ptr() + vb.ptr(); sum == int64Marker2 {
					env.A = Int64(va.unsafeInt() + vb.unsafeInt())
				} else {
					env.A = Float64(va.Float64() + vb.Float64())
				}
			case typ.String + typ.String:
				env.A = Str(va.Str() + vb.Str())
			case typ.String + typ.Number:
				env.A = Str(va.String() + vb.String())
			default:
				internal.Panic("add "+errNeedNumbersOrStrings, showType(va), showType(vb))
			}
		case typ.OpSub:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				if sum := va.ptr() + vb.ptr(); sum == int64Marker2 {
					env.A = Int64(va.unsafeInt() - vb.unsafeInt())
				} else {
					env.A = Float64(va.Float64() - vb.Float64())
				}
			} else {
				internal.Panic("sub "+errNeedNumbers, showType(va), showType(vb))
			}
		case typ.OpMul:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				if sum := va.ptr() + vb.ptr(); sum == int64Marker2 {
					env.A = Int64(va.unsafeInt() * vb.unsafeInt())
				} else {
					env.A = Float64(va.Float64() * vb.Float64())
				}
			} else {
				internal.Panic("mul "+errNeedNumbers, showType(va), showType(vb))
			}
		case typ.OpDiv:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Float64(va.Float64() / vb.Float64())
			} else {
				internal.Panic("div "+errNeedNumbers, showType(va), showType(vb))
			}
		case typ.OpIDiv:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int64(va.Int64() / vb.Int64())
			} else {
				internal.Panic("idiv "+errNeedNumbers, showType(va), showType(vb))
			}
		case typ.OpMod:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int64(va.Int64() % vb.Int64())
			} else {
				internal.Panic("mod "+errNeedNumbers, showType(va), showType(vb))
			}
		case typ.OpEq:
			env.A = Bool(env._get(opa).Equal(env._get(opb)))
		case typ.OpNeq:
			env.A = Bool(!env._get(opa).Equal(env._get(opb)))
		case typ.OpLess:
			switch va, vb := env._get(opa), env._get(opb); va.Type() + vb.Type() {
			case typ.Number + typ.Number:
				if sum := va.ptr() + vb.ptr(); sum == int64Marker2 {
					env.A = Bool(va.unsafeInt() < vb.unsafeInt())
				} else {
					env.A = Bool(va.Float64() < vb.Float64())
				}
			case typ.String + typ.String:
				env.A = Bool(lessStr(va, vb))
			default:
				internal.Panic("comparison "+errNeedNumbersOrStrings, showType(va), showType(vb))
			}
		case typ.OpLessEq:
			switch va, vb := env._get(opa), env._get(opb); va.Type() + vb.Type() {
			case typ.Number + typ.Number:
				if sum := va.ptr() + vb.ptr(); sum == int64Marker2 {
					env.A = Bool(va.unsafeInt() <= vb.unsafeInt())
				} else {
					env.A = Bool(va.Float64() <= vb.Float64())
				}
			case typ.String + typ.String:
				env.A = Bool(!lessStr(vb, va))
			default:
				internal.Panic("comparison "+errNeedNumbersOrStrings, showType(va), showType(vb))
			}
		case typ.OpNot:
			env.A = Bool(env._get(opa).IsFalse())
		case typ.OpBitAnd:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int64(va.Int64() & vb.Int64())
			} else {
				internal.Panic("bitwise and "+errNeedNumbers, showType(va), showType(vb))
			}
		case typ.OpBitOr:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int64(va.Int64() | vb.Int64())
			} else {
				internal.Panic("bitwise or "+errNeedNumbers, showType(va), showType(vb))
			}
		case typ.OpBitXor:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int64(va.Int64() ^ vb.Int64())
			} else {
				internal.Panic("bitwise xor "+errNeedNumbers, showType(va), showType(vb))
			}
		case typ.OpBitLsh:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int64(va.Int64() << vb.Int64())
			} else {
				internal.Panic("bitwise lsh "+errNeedNumbers, showType(va), showType(vb))
			}
		case typ.OpBitRsh:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int64(va.Int64() >> vb.Int64())
			} else {
				internal.Panic("bitwise rsh "+errNeedNumbers, showType(va), showType(vb))
			}
		case typ.OpBitURsh:
			if va, vb := env._get(opa), env._get(opb); va.Type()+vb.Type() == typ.Number+typ.Number {
				env.A = Int64(int64(uint64(va.Int64()) >> vb.Int64()))
			} else {
				internal.Panic("bitwise ursh "+errNeedNumbers, showType(va), showType(vb))
			}
		case typ.OpBitNot:
			if a := env._get(opa); a.Type() == typ.Number {
				env.A = Int64(^a.Int64())
			} else {
				internal.Panic("bitwise not "+errNeedNumber, showType(a))
			}
		case typ.OpCreateArray:
			env.A = NewArray(append([]Value{}, stackEnv.Stack()...)...).ToValue()
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
				env.A = False
				switch p := b.Safe().Object(); a.Type() {
				case typ.Nil:
					env.A = Bool(p == nil)
				case typ.Object:
					env.A = Bool(a.Object().IsPrototype(p))
				case typ.Bool:
					env.A = Bool(p == BoolProto)
				case typ.Number:
					env.A = Bool(p == FloatProto || (a.IsInt64() && p == IntProto))
				case typ.String:
					env.A = Bool(p == StrProto)
				case typ.Array:
					env.A = Bool(a.Array().meta.Proto.IsPrototype(p))
				}
			}
		case typ.OpStore:
			subject, v := env._get(opa), env._get(opb)
			switch subject.Type() {
			case typ.Object:
				subject.Object().Set(env.A, v)
			case typ.Array:
				if env.A.IsInt64() {
					if a, idx := subject.Array(), env.A.Int(); idx == a.Len() {
						a.Append(v)
					} else {
						a.Set(idx, v)
					}
				} else {
					internal.Panic("can't store %v into array[%v]", showType(v), showType(env.A))
				}
			case typ.Native:
				reflectStore(subject.Interface(), env.A, v)
			default:
				internal.Panic("can't store %v into (%v)[%v]", showType(v), showType(subject), showType(env.A))
			}
			env.A = v
		case typ.OpLoad:
			switch a, idx := env._get(opa), env._get(opb); a.Type() {
			case typ.Nil, typ.Number, typ.Bool:
				env.A = Nil
			case typ.Object:
				env.A = a.Object().Get(idx)
			case typ.Array:
				if idx.IsInt64() {
					env.A = a.Array().Get(idx.Int())
				} else if idx.Type() == typ.String {
					if f := a.Array().meta.Proto.Get(idx); f != Nil {
						env.A = setObjectRecv(f, a)
						break
					}
					internal.Panic("array method %q not found", idx.Str())
				} else {
					internal.Panic("can't load array[%v]", showType(idx))
				}
			case typ.Native:
				env.A = reflectLoad(a.Interface(), idx)
			case typ.String:
				if idx.Type() == typ.Number {
					if s := a.Str(); idx.Int64() >= 0 && idx.Int64() < int64(len(s)) {
						env.A = Int64(int64(s[idx.Int64()]))
					} else {
						env.A = Nil
					}
					break
				} else if idx.Type() == typ.String {
					if f := StrProto.Get(idx); f != Nil {
						env.A = setObjectRecv(f, a)
						break
					}
					internal.Panic("string method %q not found", idx.Str())
				}
				fallthrough
			default:
				internal.Panic("can't load (%v)[%v]", showType(a), showType(idx))
			}
		case typ.OpPush:
			stackEnv.Push(env._get(opa))
		case typ.OpPushUnpack:
			switch a := env._get(opa); a.Type() {
			case typ.Array:
				*stackEnv.stack = append(*stackEnv.stack, a.Array().Values()...)
			case typ.Nil:
			default:
				internal.Panic("arguments unpacking expects array, got %v", showType(a))
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
			env.stackOffset = r.StackOffset
			env.A = v
			*env.stack = (*env.stack)[:env.stackOffset+uint32(r.Callable.StackSize)]
			stackEnv.stackOffset = uint32(len(*env.stack))
			retStack = retStack[:len(retStack)-1]
		case typ.OpLoadFunc:
			env.A = env.Global.Functions[opa].ToValue()
		case typ.OpCall, typ.OpTailCall:
			a := env._get(opa)
			if a.Type() != typ.Object {
				internal.Panic("can't call %v", showType(a))
			}
			cls := a.Object().Callable
			if cls == nil || cls.Dummy {
				internal.Panic("%v not callable", showType(a))
			}
			if opb != regPhantom {
				stackEnv.Push(env._get(opb))
			}
			stackEnv.A = a.Object().this
			if cls.Variadic {
				s, w := stackEnv.Stack(), int(cls.NumParams)-1
				if len(s) > w {
					s[w] = NewArray(append([]Value{}, s[w:]...)...).ToValue()
				} else {
					stackEnv.grow(w + 1)
					stackEnv._set(uint16(w), NewArray().ToValue())
				}
			}
			if cls.Native != nil {
				stackEnv.Global = env.Global
				stackEnv.CS = K
				stackEnv.IP = cursor
				stackEnv.Stacktrace = retStack
				stackEnv.Native = cls
				cls.Native(&stackEnv)
				stackEnv.Runtime = Runtime{}
				env.A = stackEnv.A
				stackEnv.Clear()
			} else {
				stackEnv.growZero(int(cls.StackSize), int(cls.NumParams))

				last := Stacktrace{
					Callable:    K,
					Cursor:      cursor,
					StackOffset: uint32(env.stackOffset),
				}

				// Switch 'env' to 'stackEnv' and clear 'stackEnv'
				cursor = 0
				K = cls
				env.stackOffset = stackEnv.stackOffset
				env.Global = cls.LoadGlobal
				env.A = stackEnv.A

				if bop == typ.OpCall {
					retStack = append(retStack, last)
				}

				stackEnv.stackOffset = uint32(len(*env.stack))
			}
		case typ.OpJmp:
			cursor = uint32(int32(cursor) + v.b)
		case typ.OpIfNot:
			if env.A.IsFalse() {
				cursor = uint32(int32(cursor) + v.b)
			}
		}
	}
}
