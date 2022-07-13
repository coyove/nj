package bas

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

func GetGlobals() *Object {
	return globals.store.Copy(true)
}

func GetGlobalsStack() (*Object, []Value) {
	return globals.sym.Copy(true), append([]Value{}, globals.stack...)
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
	internal.GrowEnvStack = func(env unsafe.Pointer, sz int) {
		(*Env)(env).grow(sz)
	}
	internal.SetEnvStack = func(env unsafe.Pointer, stack unsafe.Pointer) {
		(*Env)(env).stack = (*[]Value)(stack)
	}
	internal.CreateRawFunc = func(name string,
		variadic bool, numParams byte, stackSize uint16, locals []string, code internal.Packet) unsafe.Pointer {
		obj := NewObject(0)
		obj.SetPrototype(Proto.Func)
		obj.fun = &funcbody{}
		obj.fun.Variadic = variadic
		obj.fun.NumParams = numParams
		obj.fun.Name = name
		obj.fun.StackSize = stackSize
		obj.fun.CodeSeg = code
		obj.fun.Locals = locals
		obj.fun.Method = strings.Contains(name, ".")
		return unsafe.Pointer(obj)
	}
	internal.NewProgram = func(coreStack, top, symbols unsafe.Pointer, funcs interface{}) unsafe.Pointer {
		cls := &Program{top: (*Object)(top)}
		cls.stack = (*Env)(coreStack).stack
		cls.symbols = (*Object)(symbols)
		cls.functions = funcs.([]*Object)
		cls.Stdout = os.Stdout
		cls.Stdin = os.Stdin
		cls.Stderr = os.Stderr

		cls.top.fun.LoadGlobal = cls
		for _, f := range cls.functions {
			f.fun.LoadGlobal = cls
		}
		return unsafe.Pointer(cls)
	}
	objEmptyFunc.Native = func(e *Env) { e.A = e.A.Object().Copy(true).ToValue() }

	AddGlobal("VERSION", Int64(Version))
	AddGlobalMethod("globals", func(e *Env) {
		e.A = e.MustGlobal().LocalsObject().ToValue()
	})
	AddGlobalMethod("new", func(e *Env) {
		m := e.Object(0)
		_ = e.Get(1).IsObject() && e.SetA(e.Object(1).SetPrototype(m).ToValue()) || e.SetA(NewObject(0).SetPrototype(m).ToValue())
	})
	AddGlobalMethod("createprototype", func(e *Env) {
		e.A = Func(e.Str(0), func(e *Env) {
			o := e.runtime.stack0.Callable
			init := o.Prop("_init").Object()
			n := o.Copy(true).SetPrototype(o)
			CallObject(init, e.runtime, nil, n.ToValue(), e.Stack()...)
			e.A = n.ToValue()
		}).Object().
			Merge(e.Object(2)).
			SetProp("_init", e.Object(1).ToValue()).
			ToValue()
	})

	// Debug libraries
	AddGlobal("debug", NewNamedObject("debug", 0).
		SetProp("self", Func("self", func(e *Env) { e.A = e.Runtime().stack1.Callable.ToValue() })).
		SetProp("locals", Func("locals", func(e *Env) {
			locals := e.Runtime().stack1.Callable.fun.Locals
			start := e.stackOffset() - uint32(e.Runtime().stack1.Callable.fun.StackSize)
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
			for i, name := range e.MustGlobal().top.fun.Locals {
				r = append(r, Int(i), Str(name), (*e.Global.stack)[i])
			}
			e.A = Array(r...)
		})).
		SetProp("set", Func("set", func(e *Env) {
			(*e.MustGlobal().stack)[e.Int64(0)] = e.Get(1)
		})).
		SetProp("trace", Func("trace", func(env *Env) {
			stacks := env.Runtime().Stacktrace(false)
			lines := make([]Value, 0, len(stacks))
			for i := len(stacks) - 1 - env.Get(0).Maybe().Int(0); i >= 0; i-- {
				r := stacks[i]
				lines = append(lines, Str(r.Callable.fun.Name), Int64(int64(r.sourceLine())), Int64(int64(r.Cursor-1)))
			}
			env.A = newArray(lines...).ToValue()
		})).
		SetProp("disfunc", Func("disfunc", func(e *Env) {
			e.A = Str(e.Object(0).GoString())
		})).
		ToValue())

	AddGlobalMethod("type", func(e *Env) { e.A = Str(e.Get(0).Type().String()) })
	AddGlobalMethod("apply", func(e *Env) { e.A = CallObject(e.Object(0), e.runtime, nil, e.Get(1), e.Stack()[2:]...) })
	AddGlobalMethod("panic", func(e *Env) {
		v := e.Get(0)
		if HasPrototype(v, Proto.Error) {
			panic(v.Native().Unwrap().(*ExecError).root)
		}
		panic(v)
	})
	AddGlobalMethod("assert", func(e *Env) {
		if v := e.Get(0); e.Size() <= 1 && v.IsFalse() {
			internal.Panic("assertion failed")
		} else if e.Size() == 2 && !v.Equal(e.Get(1)) {
			internal.Panic("assertion failed: %v and %v", v, e.Get(1))
		} else if e.Size() == 3 && !v.Equal(e.Get(1)) {
			internal.Panic("%s: %v and %v", e.Get(2).String(), v, e.Get(1))
		}
	})

	*Proto.Bool = *Func("bool", func(e *Env) { e.A = Bool(e.Get(0).IsTrue()) }).Object()
	AddGlobal("bool", Proto.Bool.ToValue())

	*Proto.Int = *Func("int", func(e *Env) {
		if v := e.Get(0); v.Type() == typ.Number {
			e.A = Int64(v.Int64())
		} else {
			v, err := strconv.ParseInt(v.String(), e.Get(1).Maybe().Int(0), 64)
			internal.PanicErr(err)
			e.A = Int64(v)
		}
	}).Object()
	AddGlobal("int", Proto.Int.ToValue())

	*Proto.Float = *Func("float", func(e *Env) {
		if v := e.Get(0); v.Type() == typ.Number {
			e.A = v
		} else {
			f, i, isInt, err := internal.ParseNumber(v.String())
			internal.PanicErr(err)
			_ = isInt && e.SetA(Int64(i)) || e.SetA(Float64(f))
		}
	}).Object()
	AddGlobal("float", Proto.Float.ToValue())

	AddGlobal("io", NewNamedObject("io", 0).
		SetProp("reader", Proto.Reader.ToValue()).
		SetProp("writer", Proto.Writer.ToValue()).
		SetProp("seeker", Proto.Seeker.ToValue()).
		SetProp("closer", Proto.Closer.ToValue()).
		SetProp("readwriter", Proto.ReadWriter.ToValue()).
		SetProp("readcloser", Proto.ReadCloser.ToValue()).
		SetProp("writecloser", Proto.WriteCloser.ToValue()).
		SetProp("readwritecloser", Proto.ReadWriteCloser.ToValue()).
		SetProp("readwriteseekcloser", Proto.ReadWriteSeekCloser.ToValue()).
		SetProp("write", Func("write", func(e *Env) {
			w := NewWriter(e.Get(0))
			for _, a := range e.Stack()[1:] {
				w.Write(ToReadonlyBytes(a))
			}
		})).
		ToValue())

	ObjectProto = *NewNamedObject("object", 0).
		SetMethod("new", func(e *Env) { e.A = NewObject(e.Get(0).Maybe().Int(0)).ToValue() }).
		SetMethod("find", func(e *Env) { e.A = e.Object(-1).Find(e.Get(0)) }).
		SetMethod("set", func(e *Env) { e.A = e.Object(-1).Set(e.Get(0), e.Get(1)) }).
		SetMethod("get", func(e *Env) { e.A = e.Object(-1).Get(e.Get(0)) }).
		SetMethod("delete", func(e *Env) { e.A = e.Object(-1).Delete(e.Get(0)) }).
		SetMethod("clear", func(e *Env) { e.Object(-1).Clear() }).
		SetMethod("copy", func(e *Env) { e.A = e.Object(-1).Copy(e.Get(0).Maybe().Bool()).ToValue() }).
		SetMethod("proto", func(e *Env) { e.A = e.Object(-1).Prototype().ToValue() }).
		SetMethod("setproto", func(e *Env) { e.Object(-1).SetPrototype(e.Object(0)) }).
		SetMethod("size", func(e *Env) { e.A = Int(e.Object(-1).Size()) }).
		SetMethod("len", func(e *Env) { e.A = Int(e.Object(-1).Len()) }).
		SetMethod("name", func(e *Env) { e.A = Str(e.Object(-1).Name()) }).
		SetMethod("setname", func(e *Env) { e.Object(-1).setName(e.Str(0)) }).
		SetMethod("keys", func(e *Env) {
			a := make([]Value, 0)
			e.Object(-1).Foreach(func(k Value, v *Value) bool { a = append(a, k); return true })
			e.A = newArray(a...).ToValue()
		}).
		SetMethod("values", func(e *Env) {
			a := make([]Value, 0)
			e.Object(-1).Foreach(func(k Value, v *Value) bool { a = append(a, *v); return true })
			e.A = newArray(a...).ToValue()
		}).
		SetMethod("items", func(e *Env) {
			a := make([]Value, 0)
			e.Object(-1).Foreach(func(k Value, v *Value) bool { a = append(a, newArray(k, *v).ToValue()); return true })
			e.A = newArray(a...).ToValue()
		}).
		SetMethod("foreach", func(e *Env) {
			f := e.Object(0)
			e.Object(-1).Foreach(func(k Value, v *Value) bool { return f.Call(e, k, *v) != False })
		}).
		SetMethod("contains", func(e *Env) { e.A = Bool(e.Object(-1).Contains(e.Get(0), e.Get(1).Maybe().Bool())) }).
		SetMethod("merge", func(e *Env) { e.A = e.Object(-1).Merge(e.Object(0)).ToValue() }).
		SetMethod("tostring", func(e *Env) {
			p := &bytes.Buffer{}
			e.Object(-1).rawPrint(p, typ.MarshalToJSON)
			e.A = UnsafeStr(p.Bytes())
		}).
		SetMethod("printed", func(e *Env) { e.A = Str(e.Object(-1).GoString()) }).
		SetMethod("pure", func(e *Env) { e.A = e.Object(-1).Copy(false).SetPrototype(&ObjectProto).ToValue() }).
		SetMethod("next", func(e *Env) { e.A = newArray(e.Object(-1).NextKeyValue(e.Get(0))).ToValue() })
	ObjectProto.SetPrototype(nil) // object is the topmost 'object', it should not have any prototype

	*Proto.Func = *NewNamedObject("function", 0).
		SetMethod("ismethod", func(e *Env) { e.A = Bool(e.Object(-1).fun.Method) }).
		SetMethod("apply", func(e *Env) { e.A = CallObject(e.Object(-1), e.runtime, nil, e.Get(0), e.Stack()[1:]...) }).
		SetMethod("call", func(e *Env) { e.A = e.Object(-1).Call(e, e.Stack()...) }).
		SetMethod("try", func(e *Env) {
			a, err := e.Object(-1).TryCall(e, e.Stack()...)
			_ = err == nil && e.SetA(a) || e.SetA(Error(e, err))
		}).
		SetMethod("after", func(e *Env) {
			f, args, e2 := e.Object(-1), e.CopyStack()[1:], e.Copy()
			t := time.AfterFunc(time.Duration(e.Float64(0)*1e6)*1e3, func() { f.Call(e2, args...) })
			e.A = NewNamedObject("Timer", 0).
				SetProp("t", ValueOf(t)).
				SetMethod("stop", func(e *Env) { e.A = Bool(e.ThisProp("t").(*time.Timer).Stop()) }).
				ToValue()
		}).
		SetMethod("go", func(e *Env) {
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
			e.A = NewNamedObject("Goroutine", 0).
				SetProp("f", f.ToValue()).
				SetProp("w", NewNative(w).ToValue()).
				SetMethod("wait", func(e *Env) { e.A = <-e.ThisProp("w").(chan Value) }).
				ToValue()
		}).
		SetMethod("map", func(e *Env) {
			if !e.Get(0).IsArray() {
				e.Get(0).AssertType(typ.Object, "map")
			}
			e.A = multiMap(e, e.Object(-1), e.Get(0), e.Get(1).Maybe().Int(1))
		}).
		SetMethod("closure", func(e *Env) {
			lambda := e.Object(-1)
			c := e.CopyStack()
			e.A = Func("<closure-"+lambda.Name()+">", func(e *Env) {
				o := e.runtime.stack0.Callable
				f := o.Prop("_l").Object()
				stk := append(o.Prop("_c").Native().Values(), e.Stack()...)
				e.A = f.Call(e, stk...)
			}).Object().
				SetProp("_l", lambda.ToValue()).
				SetProp("_c", Array(c...)).
				ToValue()
		}).
		SetPrototype(&ObjectProto)

	AddGlobal("object", ObjectProto.ToValue())
	AddGlobal("func", Proto.Func.ToValue())
	AddGlobal("callable", Proto.Func.ToValue())

	*Proto.Native = *NewNamedObject("native", 4).
		SetProp("types", nativeGoObject.ToValue()).
		SetMethod("typename", func(e *Env) {
			e.A = Str(reflect.TypeOf(e.Get(-1).Native().Unwrap()).String())
		}).
		SetMethod("isnil", func(e *Env) {
			switch rv := reflect.ValueOf(e.Native(-1).Unwrap()); rv.Kind() {
			case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
				e.A = Bool(rv.IsNil())
			default:
				e.A = False
			}
		}).
		SetMethod("natptrat", func(e *Env) {
			v := e.Native(-1).Unwrap()
			rv := reflect.ValueOf(v)
			e.A = Nil
			if rv.Kind() == reflect.Ptr && rv.Elem().Kind() == reflect.Struct {
				if t, ok := rv.Elem().Type().FieldByName(e.Str(0)); ok {
					ptr := (*struct{ a, b uintptr })(unsafe.Pointer(&v)).b + t.Offset
					e.A = ValueOf(reflect.NewAt(t.Type, unsafe.Pointer(ptr)).Interface())
				}
			}
		})
	Proto.Native.SetPrototype(nil) // native prototype has no parent
	AddGlobal("native", Proto.Native.ToValue())

	*Proto.Array = *NewNamedObject("array", 0).
		SetMethod("make", func(e *Env) {
			a := make([]Value, e.Int(0))
			if v := e.Get(1); v != Nil {
				for i := range a {
					a[i] = v
				}
			}
			e.A = Array(a...)
		}).
		SetMethod("len", func(e *Env) { e.A = Int(e.Native(-1).Len()) }).
		SetMethod("size", func(e *Env) { e.A = Int(e.Native(-1).Size()) }).
		SetMethod("clear", func(e *Env) { e.Native(-1).Clear() }).
		SetMethod("append", func(e *Env) { e.Native(-1).Append(e.Stack()...) }).
		SetMethod("find", func(e *Env) {
			a, src, ff := -1, e.Native(-1), e.Get(0)
			for i := 0; i < src.Len(); i++ {
				if src.Get(i).Equal(ff) {
					a = i
					break
				}
			}
			e.A = Int(a)
		}).
		SetMethod("filter", func(e *Env) {
			src, ff := e.Native(-1), e.Object(0)
			dest := make([]Value, 0, src.Len())
			for i := 0; i < src.Len(); i++ {
				if v := src.Get(i); ff.Call(e, v).IsTrue() {
					dest = append(dest, v)
				}
			}
			e.A = newArray(dest...).ToValue()
		}).
		SetMethod("removeat", func(e *Env) {
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
		SetMethod("last", func(e *Env) {
			if arr, n := e.Native(-1), e.Int(0); n < 0 {
				e.A = Nil
			} else if n >= arr.Len() {
				e.A = arr.ToValue()
			} else {
				e.A = arr.Slice(arr.Len()-n, arr.Len()).ToValue()
			}
		}).
		SetMethod("copy", func(e *Env) { e.Native(-1).Copy(e.Int(0), e.Int(1), e.Native(2)) }).
		SetMethod("concat", func(e *Env) { e.Native(-1).Concat(e.Native(0)) }).
		SetMethod("sort", func(e *Env) {
			a, rev := e.Native(-1), e.Get(0).Maybe().Bool()
			if kf := e.Get(1).Maybe().Func(nil); kf == nil {
				sort.Slice(a.Unwrap(), func(i, j int) bool {
					return Less(kf.Call(e, a.Get(i)), kf.Call(e, a.Get(j))) != rev
				})
			} else {
				sort.Slice(a.Unwrap(), func(i, j int) bool { return Less(a.Get(i), a.Get(j)) != rev })
			}
		}).
		SetMethod("istyped", func(e *Env) { e.A = Bool(!e.Native(-1).IsInternalArray()) }).
		SetMethod("typename", func(e *Env) { e.A = Str(e.Native(-1).meta.Name) }).
		SetMethod("untype", func(e *Env) { e.A = newArray(e.Native(-1).Values()...).ToValue() }).
		SetMethod("natptrat", func(e *Env) {
			e.A = ValueOf(reflect.ValueOf(e.Native(-1).Unwrap()).Index(e.Int(0)).Addr().Interface())
		}).
		SetPrototype(Proto.Native)
	AddGlobal("array", Proto.Array.ToValue())

	*Proto.Bytes = *Func("bytes", func(e *Env) {
		_ = e.Get(0).IsInt64() && e.SetA(ValueOf(make([]byte, e.Int(0)))) || e.SetA(ValueOf([]byte(e.Str(0))))
	}).Object().SetPrototype(Proto.Array)
	AddGlobal("bytes", Proto.Bytes.ToValue())

	*Proto.Error = *Func("error", func(e *Env) {
		e.A = Error(nil, &ExecError{root: e.Get(0), stacks: e.Runtime().Stacktrace(true)})
	}).Object().
		SetMethod("error", func(e *Env) { e.A = ValueOf(e.Native(-1).Unwrap().(*ExecError).root) }).
		SetMethod("getcause", func(e *Env) { e.A = NewNative(e.Native(-1).Unwrap().(*ExecError).root).ToValue() }).
		SetMethod("trace", func(e *Env) { e.A = ValueOf(e.Native(-1).Unwrap().(*ExecError).stacks) }).
		SetPrototype(Proto.Native)
	AddGlobal("error", Proto.Error.ToValue())

	*Proto.NativeMap = *NewNamedObject("nativemap", 4).
		SetMethod("toobject", func(e *Env) {
			rv := reflect.ValueOf(e.Native(-1).Unwrap())
			o := NewObject(rv.Len())
			for iter := rv.MapRange(); iter.Next(); {
				o.Set(ValueOf(iter.Key().Interface()), ValueOf(iter.Value().Interface()))
			}
			e.A = o.ToValue()
		}).
		SetMethod("delete", func(e *Env) {
			rv := reflect.ValueOf(e.Native(-1).Unwrap())
			rv.SetMapIndex(ToType(e.Get(0), rv.Type().Key()), reflect.Value{})
		}).
		SetMethod("clear", func(e *Env) {
			rv := reflect.ValueOf(e.Native(-1).Unwrap())
			for i := rv.MapRange(); i.Next(); {
				rv.SetMapIndex(i.Key(), reflect.Value{})
			}
		}).
		SetMethod("size", func(e *Env) { e.A = Int(e.Native(-1).Size()) }).
		SetMethod("keys", func(e *Env) {
			e.A = NewNative(reflect.ValueOf(e.Native(-1).Unwrap()).MapKeys()).ToValue()
		}).
		SetMethod("contains", func(e *Env) {
			rv := reflect.ValueOf(e.Native(-1).Unwrap())
			e.A = Bool(rv.MapIndex(ToType(e.Get(0), rv.Type().Key())).IsValid())
		}).
		SetMethod("merge", func(e *Env) {
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
				src.AssertType2(typ.Native, typ.Object, "nativemap.merge")
			}
		}).
		SetPrototype(Proto.Native)
	AddGlobal("nativemap", Proto.NativeMap.ToValue())

	*Proto.NativePtr = *NewNamedObject("nativeptr", 1).
		SetMethod("set", func(e *Env) {
			rv := reflect.ValueOf(e.Native(-1).Unwrap()).Elem()
			rv.Set(ToType(e.Get(0), rv.Type()))
		}).
		SetMethod("deref", func(e *Env) {
			rv := reflect.ValueOf(e.Native(-1).Unwrap())
			_ = rv.IsNil() && e.SetA(Nil) || e.SetA(ValueOf(rv.Elem().Interface()))
		}).
		SetPrototype(Proto.Native)
	AddGlobal("nativeptr", Proto.NativePtr.ToValue())

	*Proto.NativeIntf = *NewNamedObject("nativeintf", 1).
		SetProp("deref", Proto.NativePtr.Prop("deref")).
		SetPrototype(Proto.Native)
	AddGlobal("nativeintf", Proto.NativeIntf.ToValue())

	*Proto.Channel = *Func("channel", func(e *Env) {
		rv := reflect.ValueOf(e.Interface(0))
		_ = rv.Kind() == reflect.Chan && e.SetA(ValueOf(rv.Interface())) || e.SetA(ValueOf(make(chan Value, e.Get(0).Maybe().Int64(0))))
	}).Object().
		SetMethod("len", func(e *Env) { e.A = Int(e.Native(-1).Len()) }).
		SetMethod("size", func(e *Env) { e.A = Int(e.Native(-1).Size()) }).
		SetMethod("close", func(e *Env) { reflect.ValueOf(e.Native(-1).Unwrap()).Close() }).
		SetMethod("send", func(e *Env) {
			rv := reflect.ValueOf(e.Native(-1).Unwrap())
			rv.Send(ToType(e.Get(0), rv.Type().Elem()))
		}).
		SetMethod("recv", func(e *Env) {
			rv := reflect.ValueOf(e.Native(-1).Unwrap())
			v, ok := rv.Recv()
			e.A = newArray(ValueOf(v), Bool(ok)).ToValue()
		}).
		SetMethod("sendmulti", func(e *Env) {
			var cases []reflect.SelectCase
			var channels []Value
			if e.Get(1).Maybe().Bool() {
				cases = append(cases, reflect.SelectCase{Dir: reflect.SelectDefault})
				channels = append(channels, Str("default"))
			}
			e.Object(0).Foreach(func(ch Value, send *Value) bool {
				ch.AssertType(typ.Native, "sendmulti").Native().AssertPrototype(Proto.Channel, "sendmulti")
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
		SetMethod("recvmulti", func(e *Env) {
			var cases []reflect.SelectCase
			var channels []Value
			if e.Get(1).Maybe().Bool() {
				cases = append(cases, reflect.SelectCase{Dir: reflect.SelectDefault})
				channels = append(channels, Str("default"))
			}
			x := e.Native(0).AssertPrototype(Proto.Array, "recvmulti")
			for i := 0; i < x.Len(); i++ {
				ch := x.Get(i).AssertType(typ.Native, "recvmulti").Native().AssertPrototype(Proto.Channel, "recvmulti")
				cases = append(cases, reflect.SelectCase{
					Dir:  reflect.SelectRecv,
					Chan: reflect.ValueOf(ch.Unwrap()),
				})
				channels = append(channels, ch.ToValue())
			}
			chosen, recv, recvOK := reflect.Select(cases)
			e.A = newArray(channels[chosen], ValueOf(recv.Interface()), Bool(recvOK)).ToValue()
		}).
		SetPrototype(Proto.Native)
	AddGlobal("channel", Proto.Channel.ToValue())

	AddGlobal("chr", Func("chr", func(e *Env) { e.A = Rune(rune(e.Int(0))) }))

	*Proto.Str = *Func("str", func(e *Env) {
		i := e.Get(0)
		_ = IsBytes(i) && e.SetA(UnsafeStr(i.Native().Unwrap().([]byte))) || e.SetA(Str(i.String()))
	}).Object().
		SetMethod("from", func(e *Env) { e.A = Str(fmt.Sprint(e.Interface(0))) }).
		SetMethod("size", func(e *Env) { e.A = Int(Len(e.Get(-1))) }).
		SetMethod("len", func(e *Env) { e.A = Int(Len(e.Get(-1))) }).
		SetMethod("count", func(e *Env) { e.A = Int(utf8.RuneCountInString(e.Str(-1))) }).
		SetMethod("iequals", func(e *Env) { e.A = Bool(strings.EqualFold(e.Str(-1), e.Str(0))) }).
		SetMethod("contains", func(e *Env) { e.A = Bool(strings.Contains(e.Str(-1), e.Str(0))) }).
		SetMethod("split", func(e *Env) {
			if n := e.Get(1).Maybe().Int(0); n == 0 {
				e.A = NewNativeWithMeta(strings.Split(e.Str(-1), e.Str(0)), stringsArrayMeta).ToValue()
			} else {
				e.A = NewNativeWithMeta(strings.SplitN(e.Str(-1), e.Str(0), n), stringsArrayMeta).ToValue()
			}
		}).
		SetMethod("join", func(e *Env) {
			d := e.Str(-1)
			buf := &bytes.Buffer{}
			for x, i := e.Native(0).AssertPrototype(Proto.Array, "join"), 0; i < x.Len(); i++ {
				buf.WriteString(x.Get(i).String())
				buf.WriteString(d)
			}
			if buf.Len() > 0 {
				buf.Truncate(buf.Len() - len(d))
			}
			e.A = UnsafeStr(buf.Bytes())
		}).
		SetMethod("replace", func(e *Env) {
			e.A = Str(strings.Replace(e.Str(-1), e.Str(0), e.Str(1), e.Get(2).Maybe().Int(-1)))
		}).
		SetMethod("glob", func(e *Env) {
			m, err := filepath.Match(e.Str(-1), e.Str(0))
			internal.PanicErr(err)
			e.A = Bool(m)
		}).
		SetMethod("find", func(e *Env) {
			start, end := e.Get(1).Maybe().Int(0), e.Get(2).Maybe().Int(Len(e.Get(-1)))
			e.A = Int(strings.Index(e.Str(-1)[start:end], e.Str(0)))
		}).
		SetMethod("findsub", func(e *Env) {
			s := e.Str(-1)
			idx := strings.Index(s, e.Str(0))
			_ = idx > -1 && e.SetA(Str(s[:idx])) || e.SetA(Str(""))
		}).
		SetMethod("findlast", func(e *Env) { e.A = Int(strings.LastIndex(e.Str(-1), e.Str(0))) }).
		SetMethod("sub", func(e *Env) {
			s := e.Str(-1)
			st, en := e.Int(0), e.Get(1).Maybe().Int(len(s))
			for ; st < 0 && len(s) > 0; st += len(s) {
			}
			for ; en < 0 && len(s) > 0; en += len(s) {
			}
			e.A = Str(s[st:en])
		}).
		SetMethod("trim", func(e *Env) {
			cutset := e.Get(0).Maybe().Str("")
			_ = cutset == "" && e.SetA(Str(strings.TrimSpace(e.Str(-1)))) || e.SetA(Str(strings.Trim(e.Str(-1), e.Str(0))))
		}).
		SetMethod("trimprefix", func(e *Env) { e.A = Str(strings.TrimPrefix(e.Str(-1), e.Str(0))) }).
		SetMethod("trimsuffix", func(e *Env) { e.A = Str(strings.TrimSuffix(e.Str(-1), e.Str(0))) }).
		SetMethod("trimleft", func(e *Env) { e.A = Str(strings.TrimLeft(e.Str(-1), e.Str(0))) }).
		SetMethod("trimright", func(e *Env) { e.A = Str(strings.TrimRight(e.Str(-1), e.Str(0))) }).
		SetMethod("ord", func(e *Env) {
			r, sz := utf8.DecodeRuneInString(e.Str(-1))
			e.A = Array(Int64(int64(r)), Int(sz))
		}).
		SetMethod("startswith", func(e *Env) { e.A = Bool(strings.HasPrefix(e.Str(-1), e.Str(0))) }).
		SetMethod("endswith", func(e *Env) { e.A = Bool(strings.HasSuffix(e.Str(-1), e.Str(0))) }).
		SetMethod("upper", func(e *Env) { e.A = Str(strings.ToUpper(e.Str(-1))) }).
		SetMethod("lower", func(e *Env) { e.A = Str(strings.ToLower(e.Str(-1))) }).
		SetMethod("chars", func(e *Env) {
			var r []Value
			for s, code := e.Str(-1), e.Get(0).Maybe().Bool(); len(s) > 0; {
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
		SetMethod("format", func(e *Env) {
			buf := &bytes.Buffer{}
			EnvFprintf(e, -1, buf)
			e.A = UnsafeStr(buf.Bytes())
		})
	AddGlobal("str", Proto.Str.ToValue())
}

func EnvFprintf(env *Env, start int, p io.Writer) {
	args := make([]interface{}, 0)
	for i := start + 1; i < env.Size(); i++ {
		if v := env.Get(i); v.Type() == typ.Number {
			args = append(args, internal.SprintfNumber{I: v.Int64(), F: v.Float64(), IsInt: v.IsInt64()})
		} else {
			args = append(args, v.Interface())
		}
	}
	internal.Fprintf(p, env.Str(start), args...)
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
				defer wg.Done()
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
	internal.PanicErr(outError)
	return t
}
