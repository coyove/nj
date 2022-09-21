package bas

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
	"unsafe"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

const Version int64 = 422

var globals struct {
	sym   Object
	store Object
	stack []Value
}

func GetGlobalName(v Value) int {
	return int(globals.sym.Get(v).UnsafeInt64())
}

func Globals() *Object {
	return globals.store.Copy(true)
}

func AddGlobal(k string, v Value) {
	if len(globals.stack) == 0 {
		globals.stack = append(globals.stack, Nil)
	}
	sk := Str(k)
	idx := globals.sym.Get(sk)
	if idx != Nil {
		globals.stack[idx.Int()] = v
	} else {
		idx := len(globals.stack)
		globals.sym.Set(sk, Int(idx))
		globals.stack = append(globals.stack, v)
	}
	globals.store.Set(sk, v)
}

func AddGlobalMethod(k string, f func(*Env)) {
	AddGlobal(k, Func(k, f))
}

func init() {
	internal.NewFunc = func(f string, varg bool, np byte, ss uint16, locals, caps []string, code internal.Packet) interface{} {
		obj := NewObject(0)
		obj.SetPrototype(Proto.Func)
		obj.fun = &funcbody{}
		obj.fun.varg = varg
		obj.fun.numParams = np
		obj.fun.name = f
		obj.fun.stackSize = ss
		obj.fun.codeSeg = code
		obj.fun.locals = locals
		obj.fun.method = strings.Contains(f, ".")
		obj.fun.caps = caps
		return obj
	}
	internal.NewProgram = func(coreStack, top, symbols, funcs interface{}) interface{} {
		cls := &Program{}
		cls.main = top.(*Object)
		cls.stack = coreStack.(*[]Value)
		cls.symbols = symbols.(*Object)
		cls.functions = funcs.(*Object)
		cls.Stdout = os.Stdout
		cls.Stdin = os.Stdin
		cls.Stderr = os.Stderr

		cls.main.fun.top = cls
		cls.functions.Foreach(func(_ Value, f *Value) bool {
			f.Object().fun.top = cls
			return true
		})
		return cls
	}
	objEmptyFunc.native = func(e *Env) { e.A = e.A.Object().Copy(true).ToValue() }

	AddGlobal("VERSION", Int64(Version))
	AddGlobalMethod("globals", func(e *Env) {
		e.A = globals.store.Copy(true).ToValue()
	})
	AddGlobalMethod("new", func(e *Env) {
		m := e.Object(0)
		_ = e.Get(1).IsObject() && e.SetA(e.Object(1).SetPrototype(m).ToValue()) || e.SetA(NewObject(0).SetPrototype(m).ToValue())
	})
	AddGlobalMethod("createprototype", func(e *Env) {
		e.A = Func(e.Str(0), func(e *Env) {
			o := e.Self()
			init := o.Get(Str("_init")).Object()
			n := o.Copy(true).SetPrototype(o)
			callobj(init, e.runtime, e.top, nil, n.ToValue(), e.Stack()...)
			e.A = n.ToValue()
		}).Object().
			Merge(e.Shape(2, "No").Object()).
			SetProp("_init", e.Object(1).ToValue()).
			ToValue()
	})

	// Debug libraries
	AddGlobal("debug", NewNamedObject("debug", 0).
		SetProp("self", Func("self", func(e *Env) { e.A = e.Caller().ToValue() })).
		SetProp("locals", Func("locals", func(e *Env) {
			locals := e.Caller().fun.locals
			start := e.stackOffset() - uint32(e.Caller().fun.stackSize)
			if e.Get(0).IsTrue() {
				r := NewObject(0)
				for i, name := range locals {
					r.SetProp(name, (*e.stack)[start+uint32(i)])
				}
				e.A = r.ToValue()
			} else {
				var r []Value
				for i, name := range locals {
					idx := start + uint32(i)
					r = append(r, Int64(int64(idx)), Str(name), (*e.stack)[idx])
				}
				e.A = newArray(r...).ToValue()
			}
		})).
		SetProp("globals", Func("globals", func(e *Env) {
			var r []Value
			for i, name := range e.MustProgram().main.fun.locals {
				r = append(r, Int(i), Str(name), (*e.top.stack)[i])
			}
			e.A = Array(r...)
		})).
		SetProp("set", Func("set", func(e *Env) {
			(*e.MustProgram().stack)[e.Int64(0)] = e.Get(1)
		})).
		SetProp("trace", Func("trace", func(env *Env) {
			stacks := env.runtime.Stacktrace(false)
			lines := make([]Value, 0, len(stacks))
			for i := len(stacks) - 1 - env.IntDefault(0, 0); i >= 0; i-- {
				r := stacks[i]
				lines = append(lines, Str(r.Callable.fun.name), Int64(int64(r.sourceLine())), Int64(int64(r.Cursor-1)))
			}
			env.A = newArray(lines...).ToValue()
		})).
		SetProp("disfunc", Func("disfunc", func(e *Env) {
			e.A = Str(e.Object(0).GoString())
		})).
		ToValue())

	AddGlobalMethod("type", func(e *Env) { e.A = Str(e.Get(0).Type().String()) })
	AddGlobalMethod("apply", func(e *Env) {
		e.A = callobj(e.Object(0), e.runtime, e.top, nil, e.Get(1), e.Stack()[2:]...)
	})
	AddGlobalMethod("panic", func(e *Env) {
		v := e.Get(0)
		if HasPrototype(v, Proto.Error) {
			panic(v.Native().Unwrap().(*ExecError).root)
		}
		panic(v)
	})
	AddGlobal("assert", Func("assert", func(e *Env) {
		if v := e.Get(0); e.Size() <= 1 && v.IsFalse() {
			internal.Panic("assertion failed")
		} else if e.Size() == 2 && !v.Equal(e.Get(1)) {
			internal.Panic("assertion failed: %v and %v", v, e.Get(1))
		} else if e.Size() == 3 && !v.Equal(e.Get(1)) {
			internal.Panic("%s: %v and %v", e.Get(2).String(), v, e.Get(1))
		}
	}).Object().
		SetProp("shape", Func("shape", func(e *Env) { e.Get(0).AssertShape(e.Str(1), "assert.shape") }).
			Object().
			SetProp("describe", Func("describe", func(e *Env) {
				e.A = Str(buildShape(e.Str(0)).String())
			})).
			ToValue()).
		ToValue())

	*Proto.Bool = *Func("bool", func(e *Env) { e.A = Bool(e.Get(0).IsTrue()) }).Object()
	AddGlobal("bool", Proto.Bool.ToValue())

	*Proto.Int = *Func("int", func(e *Env) {
		if v := e.Get(0); v.Type() == typ.Number {
			e.A = Int64(v.Int64())
		} else {
			v, err := strconv.ParseInt(v.String(), e.IntDefault(1, 0), 64)
			_ = err == nil && e.SetA(Int64(v)) || e.SetA(Error(e, err))
		}
	}).Object()
	AddGlobal("int", Proto.Int.ToValue())

	*Proto.Float = *Func("float", func(e *Env) {
		if v := e.Get(0); v.Type() == typ.Number {
			e.A = v
		} else {
			if f, i, isInt, err := internal.ParseNumber(v.String()); err != nil {
				e.A = Error(e, err)
			} else {
				_ = isInt && e.SetA(Int64(i)) || e.SetA(Float64(f))
			}
		}
	}).Object()
	AddGlobal("float", Proto.Float.ToValue())

	AddGlobal("io", NewNamedObject("io", 0).
		SetProp("write", Func("write", func(e *Env) {
			w := e.Get(0).Writer()
			for _, a := range e.Stack()[1:] {
				Write(w, a)
			}
		})).
		SetProp("eof", Error(nil, io.EOF)).
		ToValue())

	ObjectProto = *NewNamedObject("object", 0)
	ObjectProto.
		AddMethod("new", func(e *Env) { e.A = NewObject(e.IntDefault(0, 0)).ToValue() }).
		AddMethod("set", func(e *Env) { e.A = e.Object(-1).Set(e.Get(0), e.Get(1)) }).
		AddMethod("get", func(e *Env) { e.A = e.Object(-1).Get(e.Get(0)) }).
		AddMethod("delete", func(e *Env) { e.A = e.Object(-1).Delete(e.Get(0)) }).
		AddMethod("clear", func(e *Env) { e.Object(-1).Clear() }).
		AddMethod("copy", func(e *Env) { e.A = e.Object(-1).Copy(e.Shape(0, "Nb").IsTrue()).ToValue() }).
		AddMethod("proto", func(e *Env) { e.A = e.Object(-1).Prototype().ToValue() }).
		AddMethod("setproto", func(e *Env) { e.Object(-1).SetPrototype(e.Object(0)) }).
		AddMethod("size", func(e *Env) { e.A = Int(e.Object(-1).Size()) }).
		AddMethod("len", func(e *Env) { e.A = Int(e.Object(-1).Len()) }).
		AddMethod("name", func(e *Env) { e.A = Str(e.Object(-1).Name()) }).
		AddMethod("setname", func(e *Env) { e.Object(-1).setName(e.Str(0)) }).
		AddMethod("keys", func(e *Env) {
			a := make([]Value, 0)
			e.Object(-1).Foreach(func(k Value, v *Value) bool { a = append(a, k); return true })
			e.A = newArray(a...).ToValue()
		}).
		AddMethod("values", func(e *Env) {
			a := make([]Value, 0)
			e.Object(-1).Foreach(func(k Value, v *Value) bool { a = append(a, *v); return true })
			e.A = newArray(a...).ToValue()
		}).
		AddMethod("items", func(e *Env) {
			a := make([]Value, 0)
			e.Object(-1).Foreach(func(k Value, v *Value) bool { a = append(a, newArray(k, *v).ToValue()); return true })
			e.A = newArray(a...).ToValue()
		}).
		AddMethod("foreach", func(e *Env) {
			f := e.Object(0)
			e.Object(-1).Foreach(func(k Value, v *Value) bool { return f.Call(e, k, *v) != False })
		}).
		AddMethod("contains", func(e *Env) { e.A = Bool(e.Object(-1).Contains(e.Get(0))) }).
		AddMethod("hasownproperty", func(e *Env) { e.A = Bool(e.Object(-1).HasOwnProperty(e.Get(0))) }).
		AddMethod("merge", func(e *Env) { e.A = e.Object(-1).Merge(e.Object(0)).ToValue() }).
		AddMethod("tostring", func(e *Env) {
			p := &bytes.Buffer{}
			e.Object(-1).rawPrint(p, typ.MarshalToJSON)
			e.A = UnsafeStr(p.Bytes())
		}).
		AddMethod("printed", func(e *Env) { e.A = Str(e.Object(-1).GoString()) }).
		AddMethod("debugprinted", func(e *Env) { e.A = Str(e.Object(-1).DebugString()) }).
		AddMethod("pure", func(e *Env) { e.A = e.Object(-1).Copy(false).SetPrototype(&ObjectProto).ToValue() }).
		AddMethod("next", func(e *Env) { e.A = newArray(e.Object(-1).FindNext(e.Get(0))).ToValue() })
	ObjectProto.SetPrototype(nil) // object is the topmost 'object', it should not have any prototype

	*Proto.Func = *NewNamedObject("function", 0).
		AddMethod("ismethod", func(e *Env) { e.A = Bool(e.Object(-1).fun.method) }).
		AddMethod("isvarg", func(e *Env) { e.A = Bool(e.Object(-1).fun.varg) }).
		AddMethod("argcount", func(e *Env) { e.A = Int(int(e.Object(-1).fun.numParams)) }).
		AddMethod("caplist", func(e *Env) { e.A = ValueOf(e.Object(-1).fun.caps) }).
		AddMethod("apply", func(e *Env) {
			e.A = callobj(e.Object(-1), e.runtime, e.top, nil, e.Get(0), e.Stack()[1:]...)
		}).
		AddMethod("call", func(e *Env) { e.A = e.Object(-1).Call(e, e.Stack()...) }).
		AddMethod("try", func(e *Env) {
			a, err := e.Object(-1).TryCall(e, e.Stack()...)
			_ = err == nil && e.SetA(a) || e.SetA(Error(e, err))
		}).
		AddMethod("after", func(e *Env) {
			f, args, e2 := e.Object(-1), e.CopyStack()[1:], e.Copy()
			t := time.AfterFunc(time.Duration(e.Float64(0)*1e6)*1e3, func() { f.Call(e2, args...) })
			e.A = NewNative(t).ToValue()
		}).
		AddMethod("go", func(e *Env) {
			f := e.Object(-1)
			args := e.CopyStack()
			w := make(chan Value, 1)
			e2 := e.Copy()
			go func(f *Object, args []Value) {
				if v, err := f.TryCall(e2, args...); err != nil {
					w <- Error(e2, err)
				} else {
					w <- v
				}
			}(f, args)
			e.A = NewNative(w).ToValue()
		}).
		AddMethod("map", func(e *Env) {
			e.A = multiMap(e, e.Object(-1), e.Shape(0, "<@array,{}>"), e.IntDefault(1, 1))
		}).
		// SetMethod("closure", func(e *Env) {
		// 	scope := e.runtime.stack1.Callable
		// 	lambda := e.Object(-1).Merge(scope).Merge(e.Shape(0, "No").Object())
		// 	start := e.stackOffset() - uint32(scope.fun.stackSize)
		// 	for addr, name := range lambda.fun.caps {
		// 		if name == "" {
		// 			continue
		// 		}
		// 		lambda.Set(Str(name), (*e.stack)[start+uint32(addr)])
		// 	}
		// 	e.A = lambda.ToValue()
		// }).
		SetPrototype(&ObjectProto)

	AddGlobal("object", ObjectProto.ToValue())
	AddGlobal("func", Proto.Func.ToValue())
	AddGlobal("callable", Proto.Func.ToValue())

	*Proto.Native = *NewNamedObject("native", 4).
		SetProp("types", nativeGoObject.ToValue()).
		AddMethod("name", func(e *Env) {
			e.A = Str(e.Get(-1).Native().meta.Name)
		}).
		AddMethod("typename", func(e *Env) {
			e.A = Str(reflect.TypeOf(e.Get(-1).Native().Unwrap()).String())
		}).
		AddMethod("isnil", func(e *Env) {
			switch rv := reflect.ValueOf(e.Native(-1).Unwrap()); rv.Kind() {
			case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
				e.A = Bool(rv.IsNil())
			default:
				e.A = False
			}
		}).
		AddMethod("natptrat", func(e *Env) {
			v := e.Native(-1).Unwrap()
			rv := reflect.ValueOf(v)
			e.A = Nil
			if rv.Kind() == reflect.Ptr && rv.Elem().Kind() == reflect.Struct {
				if t, ok := rv.Elem().Type().FieldByName(e.Str(0)); ok {
					ptr := (*struct{ a, b uintptr })(unsafe.Pointer(&v)).b + t.Offset
					e.A = ValueOf(reflect.NewAt(t.Type, unsafe.Pointer(ptr)).Interface())
				}
			}
		}).
		AddMethod("toreader", func(e *Env) { e.Native(-1).meta = NativeMetaProto.Reader }).
		AddMethod("towriter", func(e *Env) { e.Native(-1).meta = NativeMetaProto.Writer }).
		AddMethod("tocloser", func(e *Env) { e.Native(-1).meta = NativeMetaProto.Closer }).
		AddMethod("toreadwriter", func(e *Env) { e.Native(-1).meta = NativeMetaProto.ReadWriter }).
		AddMethod("toreadcloser", func(e *Env) { e.Native(-1).meta = NativeMetaProto.ReadCloser }).
		AddMethod("towritecloser", func(e *Env) { e.Native(-1).meta = NativeMetaProto.WriteCloser }).
		AddMethod("toreadwritecloser", func(e *Env) { e.Native(-1).meta = NativeMetaProto.ReadWriteCloser })

	Proto.Native.SetPrototype(nil) // native prototype has no parent
	AddGlobal("native", Proto.Native.ToValue())

	*Proto.Array = *NewNamedObject("array", 0).
		AddMethod("make", func(e *Env) {
			a := make([]Value, e.Int(0))
			if v := e.Get(1); v != Nil {
				for i := range a {
					a[i] = v
				}
			}
			e.A = Array(a...)
		}).
		AddMethod("len", func(e *Env) { e.A = Int(e.Native(-1).Len()) }).
		AddMethod("size", func(e *Env) { e.A = Int(e.Native(-1).Size()) }).
		AddMethod("clear", func(e *Env) { e.Native(-1).Clear() }).
		AddMethod("append", func(e *Env) { e.Native(-1).Append(e.Stack()...) }).
		AddMethod("find", func(e *Env) {
			a, src, ff := -1, e.Native(-1), e.Get(0)
			for i := 0; i < src.Len(); i++ {
				if src.Get(i).Equal(ff) {
					a = i
					break
				}
			}
			e.A = Int(a)
		}).
		AddMethod("filter", func(e *Env) {
			src, ff := e.Native(-1), e.Object(0)
			dest := make([]Value, 0, src.Len())
			for i := 0; i < src.Len(); i++ {
				if v := src.Get(i); ff.Call(e, v).IsTrue() {
					dest = append(dest, v)
				}
			}
			e.A = newArray(dest...).ToValue()
		}).
		AddMethod("removeat", func(e *Env) {
			ma, idx := e.Native(-1), e.Int(0)
			if idx < 0 || idx >= ma.Len() {
				e.A = Nil
			} else {
				old := ma.Get(idx)
				ma.Copy(idx, ma.Len(), ma.Slice(idx+1, ma.Len()))
				ma.SliceInplace(0, ma.Len()-1)
				e.A = old
			}
		}).
		AddMethod("last", func(e *Env) {
			if arr, n := e.Native(-1), e.Int(0); n < 0 {
				e.A = Nil
			} else if n >= arr.Len() {
				e.A = arr.ToValue()
			} else {
				e.A = arr.Slice(arr.Len()-n, arr.Len()).ToValue()
			}
		}).
		AddMethod("copy", func(e *Env) { e.Native(-1).Copy(e.Int(0), e.Int(1), e.Native(2)) }).
		AddMethod("concat", func(e *Env) { e.Native(-1).Concat(e.Shape(0, "<N,@native>").Native()) }).
		AddMethod("sort", func(e *Env) {
			a, rev := e.Native(-1), e.Shape(0, "Nb").IsTrue()
			if kf := e.Shape(1, "No").Object(); kf != nil {
				sort.Slice(a.Unwrap(), func(i, j int) bool {
					return Less(kf.Call(e, a.Get(i)), kf.Call(e, a.Get(j))) != rev
				})
			} else {
				sort.Slice(a.Unwrap(), func(i, j int) bool { return Less(a.Get(i), a.Get(j)) != rev })
			}
		}).
		AddMethod("istyped", func(e *Env) { e.A = Bool(!e.Native(-1).IsInternalArray()) }).
		AddMethod("typename", func(e *Env) { e.A = Str(e.Native(-1).meta.Name) }).
		AddMethod("untype", func(e *Env) { e.A = newArray(e.Native(-1).Values()...).ToValue() }).
		AddMethod("natptrat", func(e *Env) {
			e.A = ValueOf(reflect.ValueOf(e.Native(-1).Unwrap()).Index(e.Int(0)).Addr().Interface())
		}).
		SetPrototype(Proto.Native)
	AddGlobal("array", Proto.Array.ToValue())

	*Proto.Bytes = *Func("bytes", func(e *Env) {
		switch v := e.Get(0); v.Type() {
		case typ.Number:
			e.A = Bytes(make([]byte, v.Int()))
		case typ.String:
			e.A = Bytes([]byte(v.Str()))
		case typ.Native:
			buf := make([]byte, v.Native().Len())
			for i := 0; i < v.Native().Len(); i++ {
				buf[i] = byte(v.Native().Get(i).AssertNumber("bytes").Int())
			}
			e.A = Bytes(buf)
		default:
			v.AssertShape("<n,s,@array>", "bytes")
		}
	}).Object().
		AddMethod("unsafestr", func(e *Env) { e.A = UnsafeStr(e.A.Native().Unwrap().([]byte)) }).
		SetPrototype(Proto.Array)
	AddGlobal("bytes", Proto.Bytes.ToValue())

	*Proto.Error = *Func("error", func(e *Env) {
		e.A = Error(nil, &ExecError{root: e.Get(0), stacks: e.runtime.Stacktrace(true)})
	}).Object().
		AddMethod("error", func(e *Env) { e.A = ValueOf(e.Native(-1).Unwrap().(*ExecError).root) }).
		AddMethod("getcause", func(e *Env) { e.A = NewNative(e.Native(-1).Unwrap().(*ExecError).root).ToValue() }).
		AddMethod("trace", func(e *Env) { e.A = ValueOf(e.Native(-1).Unwrap().(*ExecError).stacks) }).
		SetPrototype(Proto.Native)
	AddGlobal("error", Proto.Error.ToValue())

	*Proto.NativeMap = *NewNamedObject("nativemap", 4).
		AddMethod("toobject", func(e *Env) {
			rv := reflect.ValueOf(e.Native(-1).Unwrap())
			o := NewObject(rv.Len())
			for iter := rv.MapRange(); iter.Next(); {
				o.Set(ValueOf(iter.Key().Interface()), ValueOf(iter.Value().Interface()))
			}
			e.A = o.ToValue()
		}).
		AddMethod("delete", func(e *Env) {
			rv := reflect.ValueOf(e.Native(-1).Unwrap())
			rv.SetMapIndex(ToType(e.Get(0), rv.Type().Key()), reflect.Value{})
		}).
		AddMethod("clear", func(e *Env) {
			rv := reflect.ValueOf(e.Native(-1).Unwrap())
			for i := rv.MapRange(); i.Next(); {
				rv.SetMapIndex(i.Key(), reflect.Value{})
			}
		}).
		AddMethod("size", func(e *Env) { e.A = Int(e.Native(-1).Size()) }).
		AddMethod("keys", func(e *Env) {
			e.A = NewNative(reflect.ValueOf(e.Native(-1).Unwrap()).MapKeys()).ToValue()
		}).
		AddMethod("contains", func(e *Env) {
			rv := reflect.ValueOf(e.Native(-1).Unwrap())
			e.A = Bool(rv.MapIndex(ToType(e.Get(0), rv.Type().Key())).IsValid())
		}).
		AddMethod("merge", func(e *Env) {
			rv, src := reflect.ValueOf(e.Native(-1).Unwrap()), e.Get(0)
			if src.Type() == typ.Object {
				rtk, rtv := rv.Type().Key(), rv.Type().Elem()
				src.Object().Foreach(func(k Value, v *Value) bool {
					rv.SetMapIndex(ToType(k, rtk), ToType(*v, rtv))
					return true
				})
			} else if src.Type() == typ.Native && src.Native().Prototype() == e.Native(-1).Prototype() {
				for i := reflect.ValueOf(src.Native().Unwrap()).MapRange(); i.Next(); {
					rv.SetMapIndex(i.Key(), i.Value())
				}
			} else {
				src.AssertShape("<@nativemap,@object>", "nativemap.merge")
			}
		}).
		SetPrototype(Proto.Native)
	AddGlobal("nativemap", Proto.NativeMap.ToValue())

	*Proto.NativePtr = *NewNamedObject("nativeptr", 1).
		AddMethod("set", func(e *Env) {
			rv := reflect.ValueOf(e.Native(-1).Unwrap()).Elem()
			rv.Set(ToType(e.Get(0), rv.Type()))
		}).
		AddMethod("deref", func(e *Env) {
			rv := reflect.ValueOf(e.Native(-1).Unwrap())
			_ = rv.IsNil() && e.SetA(Nil) || e.SetA(ValueOf(rv.Elem().Interface()))
		}).
		SetPrototype(Proto.Native)
	AddGlobal("nativeptr", Proto.NativePtr.ToValue())

	*Proto.NativeIntf = *NewNamedObject("nativeintf", 1).
		SetProp("deref", Proto.NativePtr.Get(Str("deref"))).
		SetPrototype(Proto.Native)
	AddGlobal("nativeintf", Proto.NativeIntf.ToValue())

	*Proto.Channel = *Func("channel", func(e *Env) {
		rv := reflect.ValueOf(e.Interface(0))
		_ = rv.Kind() == reflect.Chan && e.SetA(ValueOf(rv.Interface())) || e.SetA(ValueOf(make(chan Value, e.IntDefault(0, 0))))
	}).Object().
		AddMethod("len", func(e *Env) { e.A = Int(e.Native(-1).Len()) }).
		AddMethod("size", func(e *Env) { e.A = Int(e.Native(-1).Size()) }).
		AddMethod("close", func(e *Env) { reflect.ValueOf(e.Native(-1).Unwrap()).Close() }).
		AddMethod("send", func(e *Env) {
			rv := reflect.ValueOf(e.Native(-1).Unwrap())
			rv.Send(ToType(e.Get(0), rv.Type().Elem()))
		}).
		AddMethod("recv", func(e *Env) {
			rv := reflect.ValueOf(e.Native(-1).Unwrap())
			v, ok := rv.Recv()
			e.A = newArray(ValueOf(v), Bool(ok)).ToValue()
		}).
		AddMethod("sendmulti", func(e *Env) {
			var cases []reflect.SelectCase
			var channels []Value
			if e.Shape(1, "Nb").IsTrue() {
				cases = append(cases, reflect.SelectCase{Dir: reflect.SelectDefault})
				channels = append(channels, Str("default"))
			}
			e.Shape(0, "{@channel:v}").Object().Foreach(func(ch Value, send *Value) bool {
				chr := reflect.ValueOf(ch.Native().Unwrap())
				cases = append(cases, reflect.SelectCase{
					Dir:  reflect.SelectSend,
					Chan: chr,
					Send: ToType(*send, chr.Type().Elem()),
				})
				channels = append(channels, ch)
				return true
			})
			chosen, _, _ := reflect.Select(cases)
			e.A = channels[chosen]
		}).
		AddMethod("recvmulti", func(e *Env) {
			var cases []reflect.SelectCase
			var channels []Value
			if e.Shape(1, "Nb").IsTrue() {
				cases = append(cases, reflect.SelectCase{Dir: reflect.SelectDefault})
				channels = append(channels, Str("default"))
			}
			x := e.Shape(0, "[@channel]").Native()
			for i := 0; i < x.Len(); i++ {
				ch := x.Get(i).Native()
				cases = append(cases, reflect.SelectCase{
					Dir:  reflect.SelectRecv,
					Chan: reflect.ValueOf(ch.Unwrap()),
				})
				channels = append(channels, ch.ToValue())
			}
			chosen, recv, recvOK := reflect.Select(cases)
			e.A = Array(channels[chosen], ValueOf(recv.Interface()), Bool(recvOK))
		}).
		SetPrototype(Proto.Native)
	AddGlobal("channel", Proto.Channel.ToValue())

	AddGlobalMethod("chr", func(e *Env) { e.A = Rune(rune(e.Int(0))) })
	AddGlobalMethod("byte", func(e *Env) { e.A = Byte(byte(e.Int(0))) })
	AddGlobalMethod("ord", func(e *Env) { r, _ := utf8.DecodeRuneInString(e.Str(0)); e.A = Int64(int64(r)) })

	*Proto.Str = *Func("str", func(e *Env) {
		i := e.Get(0)
		_ = IsBytes(i) && e.SetA(UnsafeStr(i.Native().Unwrap().([]byte))) || e.SetA(Str(i.String()))
	}).Object().
		AddMethod("from", func(e *Env) { e.A = Str(fmt.Sprint(e.Interface(0))) }).
		AddMethod("size", func(e *Env) { e.A = Int(Len(e.Get(-1))) }).
		AddMethod("len", func(e *Env) { e.A = Int(Len(e.Get(-1))) }).
		AddMethod("count", func(e *Env) { e.A = Int(utf8.RuneCountInString(e.Str(-1))) }).
		AddMethod("iequals", func(e *Env) { e.A = Bool(strings.EqualFold(e.Str(-1), e.Str(0))) }).
		AddMethod("contains", func(e *Env) { e.A = Bool(strings.Contains(e.Str(-1), e.Str(0))) }).
		AddMethod("split", func(e *Env) {
			if n := e.IntDefault(1, 0); n == 0 {
				e.A = NewNativeWithMeta(strings.Split(e.Str(-1), e.Str(0)), stringsArrayMeta).ToValue()
			} else {
				e.A = NewNativeWithMeta(strings.SplitN(e.Str(-1), e.Str(0), n), stringsArrayMeta).ToValue()
			}
		}).
		AddMethod("join", func(e *Env) {
			d := e.Str(-1)
			buf := &bytes.Buffer{}
			for x, i := e.Shape(0, "@array").Native(), 0; i < x.Len(); i++ {
				buf.WriteString(x.Get(i).String())
				buf.WriteString(d)
			}
			if buf.Len() > 0 {
				buf.Truncate(buf.Len() - len(d))
			}
			e.A = UnsafeStr(buf.Bytes())
		}).
		AddMethod("replace", func(e *Env) {
			e.A = Str(strings.Replace(e.Str(-1), e.Str(0), e.Str(1), e.IntDefault(2, -1)))
		}).
		AddMethod("find", func(e *Env) {
			start, end := e.IntDefault(1, 0), e.IntDefault(2, Len(e.A))
			e.A = Int(strings.Index(e.Str(-1)[start:end], e.Str(0)))
		}).
		AddMethod("findsub", func(e *Env) {
			s := e.Str(-1)
			idx := strings.Index(s, e.Str(0))
			_ = idx > -1 && e.SetA(Str(s[:idx])) || e.SetA(Str(""))
		}).
		AddMethod("findlast", func(e *Env) { e.A = Int(strings.LastIndex(e.Str(-1), e.Str(0))) }).
		AddMethod("sub", func(e *Env) {
			s := e.Str(-1)
			st, en := e.Int(0), e.IntDefault(1, len(s))
			for ; st < 0 && len(s) > 0; st += len(s) {
			}
			for ; en < 0 && len(s) > 0; en += len(s) {
			}
			e.A = Str(s[st:en])
		}).
		AddMethod("trim", func(e *Env) {
			cutset := e.StrDefault(0, "", 0)
			_ = cutset == "" && e.SetA(Str(strings.TrimSpace(e.Str(-1)))) || e.SetA(Str(strings.Trim(e.Str(-1), e.Str(0))))
		}).
		AddMethod("trimprefix", func(e *Env) { e.A = Str(strings.TrimPrefix(e.Str(-1), e.Str(0))) }).
		AddMethod("trimsuffix", func(e *Env) { e.A = Str(strings.TrimSuffix(e.Str(-1), e.Str(0))) }).
		AddMethod("trimleft", func(e *Env) { e.A = Str(strings.TrimLeft(e.Str(-1), e.Str(0))) }).
		AddMethod("trimright", func(e *Env) { e.A = Str(strings.TrimRight(e.Str(-1), e.Str(0))) }).
		AddMethod("ord", func(e *Env) {
			r, sz := utf8.DecodeRuneInString(e.Str(-1))
			e.A = Array(Int64(int64(r)), Int(sz))
		}).
		AddMethod("startswith", func(e *Env) { e.A = Bool(strings.HasPrefix(e.Str(-1), e.Str(0))) }).
		AddMethod("endswith", func(e *Env) { e.A = Bool(strings.HasSuffix(e.Str(-1), e.Str(0))) }).
		AddMethod("upper", func(e *Env) { e.A = Str(strings.ToUpper(e.Str(-1))) }).
		AddMethod("lower", func(e *Env) { e.A = Str(strings.ToLower(e.Str(-1))) }).
		AddMethod("chars", func(e *Env) {
			var r []Value
			for s, code := e.Str(-1), e.Shape(0, "Nb").IsTrue(); len(s) > 0; {
				rn, sz := utf8.DecodeRuneInString(s)
				if sz == 0 {
					break
				}
				if code {
					r = append(r, Int64(int64(rn)))
				} else {
					r = append(r, Str(s[:sz]))
				}
				s = s[sz:]
			}
			e.A = newArray(r...).ToValue()
		}).
		AddMethod("format", func(e *Env) {
			buf := &bytes.Buffer{}
			Fprintf(buf, e.Str(-1), e.Stack()...)
			e.A = UnsafeStr(buf.Bytes())
		})
	AddGlobal("str", Proto.Str.ToValue())
}

func Fprintf(w io.Writer, f string, values ...Value) {
	args := make([]interface{}, 0, len(values))
	for _, v := range values {
		if v.Type() == typ.Number {
			args = append(args, internal.SprintfNumber{Int: v.Int64(), Float: v.Float64(), IsInt: v.IsInt64()})
		} else {
			args = append(args, v.Interface())
		}
	}
	internal.Fprintf(w, f, args...)
}

func Fprint(w io.Writer, values ...Value) {
	for _, v := range values {
		v.Stringify(w, typ.MarshalToString)
	}
}

func multiMap(e *Env, fun *Object, t Value, n int) Value {
	if n < 1 || n > runtime.NumCPU()*1e3 {
		internal.Panic("invalid number of goroutines: %v", n)
	}

	type payload struct {
		i int
		k Value
		v *Value
	}

	work := func(e *Env, fun *Object, outError *error, p payload) {
		if p.i == -1 {
			res, err := fun.TryCall(e, p.k, *p.v)
			if err != nil {
				*outError = err
			} else {
				*p.v = res
			}
		} else {
			res, err := fun.TryCall(e, Int(p.i), p.k)
			if err != nil {
				*outError = err
			} else {
				t.Native().Set(p.i, res)
			}
		}
	}

	var outError error
	if n == 1 {
		if t.IsArray() {
			for i := 0; outError == nil && i < t.Native().Len(); i++ {
				work(e, fun, &outError, payload{i, t.Native().Get(i), nil})
			}
		} else {
			t.Object().Foreach(func(k Value, v *Value) bool {
				work(e, fun, &outError, payload{-1, k, v})
				return outError == nil
			})
		}
	} else {
		var in = make(chan payload, Len(t))
		var wg sync.WaitGroup
		wg.Add(n)
		for i := 0; i < n; i++ {
			go func(e *Env) {
				defer func() {
					wg.Done()
					if r := recover(); r != nil {
						outError = fmt.Errorf("map fatal error: %v", r)
					}
				}()
				for p := range in {
					if outError != nil {
						return
					}
					work(e, fun, &outError, p)
				}
			}(e.Copy())
		}

		if t.IsArray() {
			for i := 0; i < t.Native().Len(); i++ {
				in <- payload{i, t.Native().Get(i), nil}
			}
		} else {
			t.Object().Foreach(func(k Value, v *Value) bool {
				in <- payload{-1, k, v}
				return true
			})
		}
		close(in)

		wg.Wait()
	}
	if outError != nil {
		return Error(e, outError)
	}
	return t
}
