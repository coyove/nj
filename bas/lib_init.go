package bas

import (
	"bytes"
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
	"unsafe"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

const Version int64 = 412

var Globals = NewObject(0)

func init() {
	internal.GrowEnvStack = func(env unsafe.Pointer, sz int) {
		(*Env)(env).grow(sz)
	}
	internal.SetObjFun = func(obj, fun unsafe.Pointer) {
		(*Object)(obj).fun = (*Function)(fun)
		(*Function)(fun).obj = (*Object)(obj)
	}
	internal.SetEnvStack = func(env unsafe.Pointer, stack unsafe.Pointer) {
		(*Env)(env).stack = (*[]Value)(stack)
	}

	Globals.SetProp("VERSION", Int64(Version))
	Globals.SetProp("globals", EnvFunc("globals", func(e *Env) {
		e.A = e.Global.LocalsObject().ToValue()
	}))
	Globals.SetMethod("new", func(e *Env) {
		m := e.Object(0)
		_ = e.Get(1).IsObject() && e.SetA(e.Object(1).SetPrototype(m).ToValue()) || e.SetA(NewObject(0).SetPrototype(m).ToValue())
	})

	// Debug libraries
	Globals.SetProp("debug", NamedObject("debug", 0).
		SetMethod("self", func(e *Env) { e.A = e.Runtime().Stack1.Callable.obj.ToValue() }).
		SetMethod("locals", func(e *Env) {
			locals := e.Runtime().Stack1.Callable.Locals
			start := e.stackOffset - uint32(e.Runtime().Stack1.Callable.StackSize)
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
		}).
		SetProp("globals", EnvFunc("globals", func(e *Env) {
			var r []Value
			for i, name := range e.Global.top.Locals {
				r = append(r, Int(i), Str(name), (*e.Global.stack)[i])
			}
			e.A = newArray(r...).ToValue()
		})).
		SetProp("set", EnvFunc("set", func(e *Env) {
			(*e.Global.stack)[e.Int64(0)] = e.Get(1)
		})).
		SetMethod("trace", func(env *Env) {
			stacks := env.Runtime().Stacktrace()
			lines := make([]Value, 0, len(stacks))
			for i := len(stacks) - 1 - env.Get(0).Maybe().Int(0); i >= 0; i-- {
				r := stacks[i]
				lines = append(lines, Str(r.Callable.Name), Int64(int64(r.sourceLine())), Int64(int64(r.Cursor-1)))
			}
			env.A = newArray(lines...).ToValue()
		}).
		SetMethod("disfunc", func(e *Env) {
			o := e.Object(0)
			_ = o.IsCallable() && e.SetA(Str(o.fun.GoString())) || e.SetA(Nil)
		}).
		SetPrototype(Proto.StaticObject).
		ToValue())

	Globals.
		SetMethod("type", func(e *Env) { e.A = Str(e.Get(0).Type().String()) }).
		SetMethod("apply", func(e *Env) { e.A = CallObject(e.Object(0), e, nil, e.Get(1), e.Stack()[2:]...) }).
		SetMethod("panic", func(e *Env) {
			v := e.Get(0)
			if HasPrototype(v, Proto.Error) {
				panic(v.Native().Unwrap().(*ExecError).root)
			}
			panic(v)
		}).
		SetProp("throw", Globals.Prop("panic")).
		SetMethod("assert", func(e *Env) {
			if v := e.Get(0); e.Size() <= 1 && v.IsFalse() {
				internal.Panic("assertion failed")
			} else if e.Size() == 2 && !v.Equal(e.Get(1)) {
				internal.Panic("assertion failed: %v and %v", v, e.Get(1))
			} else if e.Size() == 3 && !v.Equal(e.Get(1)) {
				internal.Panic("%s: %v and %v", e.Get(2).String(), v, e.Get(1))
			}
		})

	*Proto.Bool = *Func("bool", func(e *Env) { e.A = Bool(e.Get(0).IsTrue()) }).Object()
	Globals.SetProp("bool", Proto.Bool.ToValue())

	*Proto.Int = *Func("int", func(e *Env) {
		if v := e.Get(0); v.Type() == typ.Number {
			e.A = Int64(v.Int64())
		} else {
			v, err := strconv.ParseInt(v.String(), e.Get(1).Maybe().Int(0), 64)
			internal.PanicErr(err)
			e.A = Int64(v)
		}
	}).Object()
	Globals.SetProp("int", Proto.Int.ToValue())

	*Proto.Float = *Func("float", func(e *Env) {
		if v := e.Get(0); v.Type() == typ.Number {
			e.A = v
		} else {
			f, i, isInt, err := internal.ParseNumber(v.String())
			internal.PanicErr(err)
			_ = isInt && e.SetA(Int64(i)) || e.SetA(Float64(f))
		}
	}).Object()
	Globals.SetProp("float", Proto.Float.ToValue())

	Globals.SetProp("io", NamedObject("io", 0).
		SetProp("reader", Proto.Reader.ToValue()).
		SetProp("writer", Proto.Writer.ToValue()).
		SetProp("seeker", Proto.Seeker.ToValue()).
		SetProp("closer", Proto.Closer.ToValue()).
		SetProp("readwriter", Proto.ReadWriter.ToValue()).
		SetProp("readcloser", Proto.ReadCloser.ToValue()).
		SetProp("writecloser", Proto.WriteCloser.ToValue()).
		SetProp("readwritecloser", Proto.ReadWriteCloser.ToValue()).
		SetProp("readwriteseekcloser", Proto.ReadWriteSeekCloser.ToValue()).
		SetMethod("write", func(e *Env) {
			w := NewWriter(e.Get(0))
			for _, a := range e.Stack()[1:] {
				w.Write(ToReadonlyBytes(a))
			}
		}).
		SetPrototype(Proto.StaticObject).
		ToValue())

	ObjectProto = *NamedObject("object", 0).
		SetMethod("new", func(e *Env) { e.A = NewObject(e.Get(0).Maybe().Int(0)).ToValue() }).
		SetMethod("newstatic", func(e *Env) { e.A = NewObject(e.Get(0).Maybe().Int(0)).SetPrototype(Proto.StaticObject).ToValue() }).
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
			e.Object(-1).Foreach(func(k Value, v *Value) bool { return e.Call(f, k, *v) != False })
		}).
		SetMethod("contains", func(e *Env) { e.A = Bool(e.Object(-1).Contains(e.Get(0), e.Get(1).Maybe().Bool())) }).
		SetMethod("merge", func(e *Env) { e.A = e.Object(-1).Merge(e.Object(0)).ToValue() }).
		SetMethod("tostring", func(e *Env) {
			p := &bytes.Buffer{}
			e.Object(-1).rawPrint(p, typ.MarshalToJSON, true)
			e.A = UnsafeStr(p.Bytes())
		}).
		SetMethod("pure", func(e *Env) { e.A = e.Object(-1).Copy(false).SetPrototype(&ObjectProto).ToValue() }).
		SetMethod("next", func(e *Env) { e.A = newArray(e.Object(-1).NextKeyValue(e.Get(0))).ToValue() })
	ObjectProto.SetPrototype(nil) // object is the topmost 'object', it should not have any prototype

	*Proto.Func = *NamedObject("function", 0).
		SetMethod("apply", func(e *Env) { e.A = CallObject(e.Object(-1), e, nil, e.Get(0), e.Stack()[1:]...) }).
		SetMethod("call", func(e *Env) { e.A = e.Call(e.Object(-1), e.Stack()...) }).
		SetMethod("try", func(e *Env) {
			a, err := e.Call2(e.Object(-1), e.Stack()...)
			_ = err == nil && e.SetA(a) || e.SetA(Error(e, err))
		}).
		SetMethod("after", func(e *Env) {
			f, args, e2 := e.Object(-1), e.CopyStack()[1:], EnvForAsyncCall(e)
			t := time.AfterFunc(time.Duration(e.Float64(0)*1e6)*1e3, func() { e2.Call(f, args...) })
			e.A = NamedObject("Timer", 0).
				SetProp("t", ValueOf(t)).
				SetMethod("stop", func(e *Env) { e.A = Bool(e.ThisProp("t").(*time.Timer).Stop()) }).
				ToValue()
		}).
		SetMethod("go", func(e *Env) {
			f := e.Object(-1)
			args := e.CopyStack()
			w := make(chan Value, 1)
			e2 := EnvForAsyncCall(e)
			go func(f *Object, args []Value) {
				if v, err := e2.Call2(f, args...); err != nil {
					w <- Error(e2, err)
				} else {
					w <- v
				}
			}(f, args)
			e.A = NamedObject("Goroutine", 0).
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
				o := e.runtime.Callable0.obj
				f := o.Prop("_l").Object()
				stk := append(o.Prop("_c").Native().Values(), e.Stack()...)
				e.A = e.Call(f, stk...)
			}).Object().
				SetProp("_l", lambda.ToValue()).
				SetProp("_c", Array(c...)).
				ToValue()
		}).
		SetPrototype(&ObjectProto)

	Globals.SetProp("object", ObjectProto.ToValue())
	Globals.SetProp("staticobject", Proto.StaticObject.ToValue())
	Globals.SetProp("func", Proto.Func.ToValue())
	Globals.SetProp("callable", Proto.Func.ToValue())

	*Proto.Native = *NamedObject("native", 4).
		SetProp("types", nativeGoObject.ToValue()).
		SetMethod("typename", func(e *Env) {
			e.A = Str(reflect.TypeOf(e.Get(-1).Native().Unwrap()).String())
		})
	Proto.Native.SetPrototype(nil) // native prototype has no parent
	Globals.SetProp("native", Proto.Native.ToValue())

	*Proto.Array = *NamedObject("array", 0).
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
				if v := src.Get(i); e.Call(ff, v).IsTrue() {
					dest = append(dest, v)
				}
			}
			e.A = newArray(dest...).ToValue()
		}).
		SetMethod("slice", func(e *Env) {
			a := e.Native(-1)
			st, en := e.Int(0), e.Get(1).Maybe().Int(a.Len())
			for ; st < 0 && a.Len() > 0; st += a.Len() {
			}
			for ; en < 0 && a.Len() > 0; en += a.Len() {
			}
			e.A = a.Slice(st, en).ToValue()
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
		SetMethod("copy", func(e *Env) { e.Native(-1).Copy(e.Int(0), e.Int(1), e.Native(2)) }).
		SetMethod("concat", func(e *Env) { e.Native(-1).Concat(e.Native(0)) }).
		SetMethod("sort", func(e *Env) {
			a, rev := e.Native(-1), e.Get(0).Maybe().Bool()
			if kf := e.Get(1).Maybe().Func(nil); kf == nil {
				sort.Slice(a.Unwrap(), func(i, j int) bool {
					return Less(e.Call(kf, a.Get(i)), e.Call(kf, a.Get(j))) != rev
				})
			} else {
				sort.Slice(a.Unwrap(), func(i, j int) bool { return Less(a.Get(i), a.Get(j)) != rev })
			}
		}).
		SetMethod("istyped", func(e *Env) { e.A = Bool(!e.Native(-1).IsInternalArray()) }).
		SetMethod("typename", func(e *Env) { e.A = Str(e.Native(-1).meta.Name) }).
		SetMethod("untype", func(e *Env) { e.A = newArray(e.Native(-1).Values()...).ToValue() }).
		SetPrototype(Proto.Native)
	Globals.SetProp("array", Proto.Array.ToValue())

	*Proto.Bytes = *Func("bytes", func(e *Env) {
		_ = e.Get(0).IsInt64() && e.SetA(ValueOf(make([]byte, e.Int(0)))) || e.SetA(ValueOf([]byte(e.Str(0))))
	}).Object().SetPrototype(Proto.Array)
	Globals.SetProp("bytes", Proto.Bytes.ToValue())

	*Proto.Error = *Func("error", func(e *Env) {
		e.A = Error(nil, &ExecError{root: e.Get(0), stacks: e.Runtime().Stacktrace()})
	}).Object().
		SetMethod("error", func(e *Env) { e.A = ValueOf(e.Native(-1).Unwrap().(*ExecError).root) }).
		SetMethod("getcause", func(e *Env) { e.A = NewNative(e.Native(-1).Unwrap().(*ExecError).root).ToValue() }).
		SetMethod("trace", func(e *Env) { e.A = ValueOf(e.Native(-1).Unwrap().(*ExecError).stacks) }).
		SetPrototype(Proto.Native)
	Globals.SetProp("error", Proto.Error.ToValue())

	*Proto.NativeMap = *NamedObject("nativemap", 4).
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
	Globals.SetProp("nativemap", Proto.NativeMap.ToValue())

	*Proto.NativePtr = *NamedObject("nativeptr", 1).
		SetMethod("deref", func(e *Env) { e.A = ValueOf(reflect.ValueOf(e.Native(-1).Unwrap()).Elem().Interface()) }).
		SetPrototype(Proto.Native)
	Globals.SetProp("nativeptr", Proto.NativePtr.ToValue())

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
	Globals.SetProp("channel", Proto.Channel.ToValue())

	*Proto.Str = *Func("str", func(e *Env) {
		i, ok := e.Interface(0).([]byte)
		_ = ok && e.SetA(UnsafeStr(i)) || e.SetA(Str(e.Get(0).String()))
	}).Object().
		SetMethod("from", func(e *Env) { e.A = Str(fmt.Sprint(e.Interface(0))) }).
		SetMethod("size", func(e *Env) { e.A = Int(Len(e.Get(-1))) }).
		SetMethod("len", func(e *Env) { e.A = Int(Len(e.Get(-1))) }).
		SetMethod("count", func(e *Env) { e.A = Int(utf8.RuneCountInString(e.Str(-1))) }).
		SetMethod("iequals", func(e *Env) { e.A = Bool(strings.EqualFold(e.Str(-1), e.Str(0))) }).
		SetMethod("contains", func(e *Env) { e.A = Bool(strings.Contains(e.Str(-1), e.Str(0))) }).
		SetMethod("split", func(e *Env) {
			if n := e.Get(1).Maybe().Int(0); n == 0 {
				e.A = newNativeWithType(strings.Split(e.Str(-1), e.Str(0)), stringsArrayMeta).ToValue()
			} else {
				e.A = newNativeWithType(strings.SplitN(e.Str(-1), e.Str(0), n), stringsArrayMeta).ToValue()
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
			sprintf(e, -1, buf)
			e.A = UnsafeStr(buf.Bytes())
		})
	Globals.SetProp("str", Proto.Str.ToValue())

	Globals.SetProp("printf", EnvFunc("printf", func(e *Env) {
		sprintf(e, 0, e.Global.Stdout)
	}))
	Globals.SetProp("println", EnvFunc("println", func(e *Env) {
		for _, a := range e.Stack() {
			fmt.Fprint(e.Global.Stdout, a.String(), " ")
		}
		fmt.Fprintln(e.Global.Stdout)
	}))
	Globals.SetProp("print", EnvFunc("print", func(e *Env) {
		for _, a := range e.Stack() {
			fmt.Fprint(e.Global.Stdout, a.String())
		}
		fmt.Fprintln(e.Global.Stdout)
	}))
}
