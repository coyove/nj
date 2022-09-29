package bas

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
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

const Version int64 = 481

var objDefaultFun = &funcbody{name: "object"}

var Proto struct {
	Object, Bool, Str, Int, Float, Func, Native, Bytes, Array, Error             Object
	NativeMap, NativePtr, NativeIntf, Channel                                    Object
	Reader, Writer, Closer, ReadWriter, ReadCloser, WriteCloser, ReadWriteCloser NativeMeta
	ReadlinesMeta, ArrayMeta, BytesMeta, StringsMeta, ErrorMeta                  NativeMeta
}

func init() {
	objDefaultFun.native = func(e *Env) { e.A = e.A.Object().Copy(true).ToValue() }
	metaStore.metaCache = map[reflect.Type]*NativeMeta{}
	metaStore.nativeTypes = NewObject(0)

	Proto.ArrayMeta = NativeMeta{
		"internal",
		&Proto.Array,
		func(a *Native) int { return len(a.internal) },
		func(a *Native) int { return cap(a.internal) },
		func(a *Native) { a.internal = a.internal[:0] },
		func(a *Native) []Value { return a.internal },
		func(a *Native, idx int) Value { return a.internal[idx] },
		func(a *Native, idx int, v Value) { a.internal[idx] = v },
		sgGetKey,
		sgSetKeyNotSupported,
		func(a *Native, v ...Value) { a.internal = append(a.internal, v...) },
		func(a *Native, s, e int) *Native { return &Native{meta: a.meta, internal: a.internal[s:e]} },
		func(a *Native, s, e int) { a.internal = a.internal[s:e] },
		func(a *Native, s, e int, from *Native) {
			if from.meta != a.meta {
				for i := s; i < e; i++ {
					a.internal[i] = from.Get(i - s)
				}
			} else {
				copy(a.internal[s:e], from.internal)
			}
		},
		func(a *Native, b *Native) {
			if a.meta != b.meta {
				for i := 0; i < b.Len(); i++ {
					a.internal = append(a.internal, b.Get(i))
				}
			} else {
				a.internal = append(a.internal, b.internal...)
			}
		},
		func(a *Native, w io.Writer, mt typ.MarshalType) {
			internal.WriteString(w, "[")
			for i, v := range a.internal {
				internal.WriteString(w, internal.IfStr(i == 0, "", ","))
				v.Stringify(w, mt.NoRec())
			}
			internal.WriteString(w, "]")
		},
		sgArrayNext,
	}

	Proto.BytesMeta = NativeMeta{
		"bytes",
		&Proto.Bytes,
		func(a *Native) int { return len((a.any).([]byte)) },
		func(a *Native) int { return cap((a.any).([]byte)) },
		func(a *Native) { a.any = a.any.([]byte)[:0] },
		sgValuesNotSupported,
		func(a *Native, idx int) Value { return Int64(int64(a.any.([]byte)[idx])) },
		func(a *Native, idx int, v Value) {
			a.any.([]byte)[idx] = byte(v.AssertNumber("bytes.Set").Int())
		},
		sgGetKey,
		sgSetKeyNotSupported,
		func(a *Native, v ...Value) {
			p := a.any.([]byte)
			for _, b := range v {
				p = append(p, byte(b.AssertNumber("bytes.Append").Int()))
			}
			a.any = p
		},
		func(a *Native, start, end int) *Native {
			return &Native{meta: a.meta, any: a.any.([]byte)[start:end]}
		},
		func(a *Native, start, end int) {
			a.any = a.any.([]byte)[start:end]
		},
		func(a *Native, start, end int, from *Native) {
			if from.meta == &Proto.ArrayMeta {
				buf := a.any.([]byte)
				for i := start; i < end; i++ {
					buf[i] = byte(from.Get(i - start).AssertNumber("bytes.Copy").Int())
				}
			} else {
				copy(a.any.([]byte)[start:end], from.any.([]byte))
			}
		},
		func(a *Native, b *Native) {
			if b.meta == &Proto.ArrayMeta {
				buf := a.any.([]byte)
				for i := 0; i < b.Len(); i++ {
					buf[i] = byte(b.Get(i).AssertNumber("bytes.Concat").Int())
				}
				a.any = buf
			} else {
				a.any = append(a.any.([]byte), b.any.([]byte)...)
			}
		},
		sgMarshal,
		sgArrayNext,
	}

	Proto.StringsMeta = NativeMeta{
		"strings",
		&Proto.Array,
		func(a *Native) int { return len((a.any).([]string)) },
		func(a *Native) int { return cap((a.any).([]string)) },
		func(a *Native) { a.any = a.any.([]byte)[:0] },
		func(a *Native) []Value {
			res := make([]Value, a.Len())
			for i := 0; i < a.Len(); i++ {
				res[i] = a.Get(i)
			}
			return res
		},
		func(a *Native, idx int) Value { return Str(a.any.([]string)[idx]) },
		func(a *Native, idx int, v Value) {
			a.any.([]string)[idx] = v.AssertString("strings.Set")
		},
		sgGetKey,
		sgSetKeyNotSupported,
		func(a *Native, v ...Value) {
			p := a.any.([]string)
			for _, b := range v {
				p = append(p, b.AssertString("strings.Append"))
			}
			a.any = p
		},
		func(a *Native, start, end int) *Native {
			return &Native{meta: a.meta, any: a.any.([]string)[start:end]}
		},
		func(a *Native, start, end int) { a.any = a.any.([]string)[start:end] },
		func(a *Native, start, end int, from *Native) {
			if from.meta == &Proto.ArrayMeta {
				buf := a.any.([]string)
				for i := start; i < end; i++ {
					buf[i] = from.Get(i - start).AssertString("strings.Copy")
				}
			} else {
				copy(a.any.([]byte)[start:end], from.any.([]byte))
			}
		},
		func(a *Native, b *Native) {
			if b.meta == &Proto.ArrayMeta {
				buf := a.any.([]string)
				for i := 0; i < b.Len(); i++ {
					buf[i] = b.Get(i).AssertString("strings.Concat")
				}
				a.any = buf
			} else {
				a.any = append(a.any.([]byte), b.any.([]byte)...)
			}
		},
		sgMarshal,
		sgArrayNext,
	}

	Proto.ErrorMeta = *NewEmptyNativeMeta("error", &Proto.Error)
	Proto.ErrorMeta.Marshal = func(a *Native, w io.Writer, mt typ.MarshalType) {
		internal.WriteString(w, internal.IfQuote(mt == typ.MarshalToJSON, a.any.(*ExecError).Error()))
	}

	Proto.ReadlinesMeta = *NewEmptyNativeMeta("readlines", &Proto.Native)
	Proto.ReadlinesMeta.Next = func(a *Native, k Value) Value {
		if k.IsNil() {
			k = Array(Int(-1), Nil)
		}
		var er error
		if s := a.any.(*ioReadlinesStruct); s.bytes {
			line, err := s.rd.ReadBytes(s.delim)
			if len(line) > 0 {
				k.Native().Set(0, Int(k.Native().Get(0).Int()+1))
				k.Native().Set(1, Bytes(line))
				return k
			}
			er = err
		} else {
			line, err := s.rd.ReadString(s.delim)
			if len(line) > 0 {
				k.Native().Set(0, Int(k.Native().Get(0).Int()+1))
				k.Native().Set(1, Str(line))
				return k
			}
			er = err
		}
		if er == io.EOF {
			return Nil
		}
		return Error(nil, er)
	}

	Proto.Reader = *createNativeMeta("Reader", NewNamedObject("Reader", 0).
		AddMethod("read", func(e *Env) {
			buf, err := func(e *Env) ([]byte, error) {
				f := e.A.Reader()
				if n := e.Get(0); n.Type() == typ.Number {
					p := make([]byte, n.Int())
					rn, err := f.Read(p)
					if err == nil || rn > 0 {
						return p[:rn], nil
					} else if err == io.EOF {
						return nil, nil
					}
					return nil, err
				}
				return ioutil.ReadAll(f)
			}(e)
			_ = err != nil && e.SetA(Error(e, err)) || e.SetA(Bytes(buf))
		}).
		AddMethod("read2", func(e *Env) {
			rn, err := e.A.Reader().Read(e.Shape(0, "B").Native().Unwrap().([]byte))
			e.A = Array(Int(rn), Error(e, err)) // return in Go style
		}).
		AddMethod("readlines", func(e *Env) {
			e.A = NewNativeWithMeta(&ioReadlinesStruct{
				rd:    bufio.NewReader(e.A.Reader()),
				delim: e.StrDefault(0, "\n", 1)[0],
				bytes: e.Shape(1, "Nb").IsTrue(),
			}, &Proto.ReadlinesMeta).ToValue()
		}).
		SetPrototype(&Proto.Native))

	Proto.Writer = *createNativeMeta("Writer", NewNamedObject("Writer", 0).
		AddMethod("write", func(e *Env) {
			wn, err := Write(e.A.Writer(), e.Get(0))
			_ = err == nil && e.SetA(Int(wn)) || e.SetA(Error(e, err))
		}).
		AddMethod("write2", func(e *Env) {
			wn, err := Write(e.A.Writer(), e.Get(0))
			e.A = Array(Int(wn), Error(e, err))
		}).
		AddMethod("pipe", func(e *Env) {
			var wn int64
			var err error
			if n := e.IntDefault(1, 0); n > 0 {
				wn, err = io.CopyN(e.Get(-1).Writer(), e.Get(0).Reader(), int64(n))
			} else {
				wn, err = io.Copy(e.Get(-1).Writer(), e.Get(0).Reader())
			}
			_ = err == nil && e.SetA(Int64(wn)) || e.SetA(Error(e, err))
		}).
		SetPrototype(&Proto.Native))

	Proto.Closer = *createNativeMeta("Closer", NewNamedObject("Closer", 0).
		AddMethod("close", func(e *Env) {
			e.A = Error(e, e.A.Closer().Close())
		}).
		SetPrototype(&Proto.Native))

	Proto.ReadWriter = *createNativeMeta("ReadWriter", NewNamedObject("ReadWriter", 0).
		Merge(Proto.Reader.Proto).
		Merge(Proto.Writer.Proto).SetPrototype(&Proto.Native))

	Proto.ReadCloser = *createNativeMeta("ReadCloser", NewNamedObject("ReadCloser", 0).
		Merge(Proto.Reader.Proto).
		Merge(Proto.Closer.Proto).SetPrototype(&Proto.Native))

	Proto.WriteCloser = *createNativeMeta("WriteCloser", NewNamedObject("WriteCloser", 0).
		Merge(Proto.Writer.Proto).
		Merge(Proto.Closer.Proto).SetPrototype(&Proto.Native))

	Proto.ReadWriteCloser = *createNativeMeta("ReadWriteCloser", NewNamedObject("ReadWriteCloser", 0).
		Merge(Proto.ReadWriter.Proto).
		Merge(Proto.Closer.Proto).SetPrototype(&Proto.Native))

	AddGlobal("Version", Int64(Version))
	AddGlobalMethod("globals", func(e *Env) { e.A = globals.store.Copy(true).ToValue() })
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
		SetProp("deepequal", Func("deepequal", func(e *Env) {
			e.A = Bool(DeepEqual(e.Get(0), e.Get(1)))
		})).
		ToValue())

	AddGlobalMethod("type", func(e *Env) { e.A = Str(e.Get(0).Type().String()) })
	AddGlobalMethod("apply", func(e *Env) {
		e.A = callobj(e.Object(0), e.runtime, e.top, nil, e.Get(1), e.Stack()[2:]...)
	})
	AddGlobalMethod("panic", func(e *Env) {
		v := e.Get(0)
		if v.HasPrototype(&Proto.Error) {
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

	Proto.Bool = *Func("bool", func(e *Env) { e.A = Bool(e.Get(0).IsTrue()) }).Object()
	AddGlobal("bool", Proto.Bool.ToValue())

	Proto.Int = *Func("int", func(e *Env) {
		if v := e.Get(0); v.Type() == typ.Number {
			e.A = Int64(v.Int64())
		} else {
			v, err := strconv.ParseInt(v.String(), e.IntDefault(1, 0), 64)
			_ = err == nil && e.SetA(Int64(v)) || e.SetA(Error(e, err))
		}
	}).Object()
	AddGlobal("int", Proto.Int.ToValue())

	Proto.Float = *Func("float", func(e *Env) {
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

	Proto.Object = *NewNamedObject("object", 0)
	Proto.Object.
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
		AddMethod("merge", func(e *Env) { e.A = e.Object(-1).Merge(e.Shape(0, "No").Object()).ToValue() }).
		AddMethod("tostring", func(e *Env) {
			p := &bytes.Buffer{}
			e.Object(-1).rawPrint(p, typ.MarshalToJSON)
			e.A = UnsafeStr(p.Bytes())
		}).
		AddMethod("printed", func(e *Env) { e.A = Str(e.Object(-1).GoString()) }).
		AddMethod("debugprinted", func(e *Env) { e.A = Str(e.Object(-1).DebugString()) }).
		AddMethod("pure", func(e *Env) { e.A = e.Object(-1).Copy(false).SetPrototype(&Proto.Object).ToValue() }).
		AddMethod("next", func(e *Env) { e.A = newArray(e.Object(-1).FindNext(e.Get(0))).ToValue() })
	Proto.Object.SetPrototype(nil) // object is the topmost 'object', it should not have any prototype

	Proto.Func = *NewNamedObject("function", 0).
		AddMethod("ismethod", func(e *Env) { e.A = Bool(e.Object(-1).fun.method) }).
		AddMethod("isvarg", func(e *Env) { e.A = Bool(e.Object(-1).fun.varg) }).
		AddMethod("isnative", func(e *Env) { e.A = Bool(e.Object(-1).fun.native != nil) }).
		AddMethod("argcount", func(e *Env) { e.A = Int(int(e.Object(-1).fun.numArgs)) }).
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
			go func(f *Object, args []Value) { w <- f.Call(e2, args...) }(f, args)
			e.A = NewNative(w).ToValue()
		}).
		AddMethod("map", func(e *Env) {
			e.A = multiMap(e, e.Object(-1), e.Shape(0, "<@array,{}>"), e.IntDefault(1, 1))
		}).
		SetPrototype(&Proto.Object)

	AddGlobal("object", Proto.Object.ToValue())
	AddGlobal("func", Proto.Func.ToValue())
	AddGlobal("callable", Proto.Func.ToValue())

	Proto.Native = *NewNamedObject("native", 0).
		SetProp("types", metaStore.nativeTypes.ToValue()).
		AddMethod("typename", func(e *Env) { e.A = Str(e.Get(-1).Native().meta.Name) }).
		AddMethod("nativename", func(e *Env) { e.A = Str(reflect.TypeOf(e.Get(-1).Native().Unwrap()).String()) }).
		AddMethod("isnil", func(e *Env) {
			switch rv := reflect.ValueOf(e.Native(-1).Unwrap()); rv.Kind() {
			case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
				e.A = Bool(rv.IsNil())
			default:
				e.A = False
			}
		}).
		AddMethod("fields", func(e *Env) {
			var list []string
			rt := reflect.Indirect(reflect.ValueOf(e.Native(-1).Unwrap())).Type()
			for i := 0; i < rt.NumField(); i++ {
				list = append(list, rt.Field(i).Name)
			}
			e.A = NewNative(list).ToValue()
		}).
		AddMethod("methods", func(e *Env) {
			var list []string
			rv := reflect.ValueOf(e.Native(-1).Unwrap())
			rt := rv.Type()
			for i := 0; i < rt.NumMethod(); i++ {
				list = append(list, rt.Method(i).Name)
			}
			if rt.Kind() == reflect.Ptr && rv.CanAddr() {
				for rt, i := rv.Addr().Type(), 0; i < rt.NumMethod(); i++ {
					list = append(list, rt.Method(i).Name)
				}
			}
			e.A = NewNative(list).ToValue()
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
		AddMethod("toreader", func(e *Env) { e.Native(-1).meta = &Proto.Reader }).
		AddMethod("towriter", func(e *Env) { e.Native(-1).meta = &Proto.Writer }).
		AddMethod("tocloser", func(e *Env) { e.Native(-1).meta = &Proto.Closer }).
		AddMethod("toreadwriter", func(e *Env) { e.Native(-1).meta = &Proto.ReadWriter }).
		AddMethod("toreadcloser", func(e *Env) { e.Native(-1).meta = &Proto.ReadCloser }).
		AddMethod("towritecloser", func(e *Env) { e.Native(-1).meta = &Proto.WriteCloser }).
		AddMethod("toreadwritecloser", func(e *Env) { e.Native(-1).meta = &Proto.ReadWriteCloser })

	Proto.Native.SetPrototype(nil) // native prototype has no parent
	AddGlobal("native", Proto.Native.ToValue())

	Proto.Array = *NewNamedObject("array", 0).
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
					return kf.Call(e, a.Get(i)).Less(kf.Call(e, a.Get(j))) != rev
				})
			} else {
				sort.Slice(a.Unwrap(), func(i, j int) bool { return a.Get(i).Less(a.Get(j)) != rev })
			}
		}).
		AddMethod("istyped", func(e *Env) { e.A = Bool(e.Native(-1).IsTypedArray()) }).
		AddMethod("untype", func(e *Env) { e.A = Array(e.Native(-1).Values()...) }).
		AddMethod("natptrat", func(e *Env) {
			e.A = ValueOf(reflect.ValueOf(e.Native(-1).Unwrap()).Index(e.Int(0)).Addr().Interface())
		}).
		SetPrototype(&Proto.Native)
	AddGlobal("array", Proto.Array.ToValue())

	Proto.Bytes = *Func("bytes", func(e *Env) {
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
		SetPrototype(&Proto.Array)
	AddGlobal("bytes", Proto.Bytes.ToValue())

	Proto.Error = *Func("error", func(e *Env) {
		e.A = Error(nil, &ExecError{root: e.Get(0), stacks: e.runtime.Stacktrace(true)})
	}).Object().
		AddMethod("equals", func(e *Env) { e.A = Bool(ToErrorRootCause(e.A) == ToErrorRootCause(e.Shape(0, "E"))) }).
		AddMethod("error", func(e *Env) { e.A = ValueOf(e.Native(-1).Unwrap().(*ExecError).root) }).
		AddMethod("getcause", func(e *Env) { e.A = NewNative(e.Native(-1).Unwrap().(*ExecError).root).ToValue() }).
		AddMethod("trace", func(e *Env) { e.A = ValueOf(e.Native(-1).Unwrap().(*ExecError).stacks) }).
		SetPrototype(&Proto.Native)
	AddGlobal("error", Proto.Error.ToValue())

	Proto.NativeMap = *Func("nativemap", func(e *Env) {
		m := make(map[Value]Value, e.Int(0))
		e.A = NewNative(m).ToValue()
	}).Object().
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
			rv.SetMapIndex(e.Get(0).ToType(rv.Type().Key()), reflect.Value{})
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
			e.A = Bool(rv.MapIndex(e.Get(0).ToType(rv.Type().Key())).IsValid())
		}).
		AddMethod("merge", func(e *Env) {
			rv, src := reflect.ValueOf(e.Native(-1).Unwrap()), e.Get(0)
			if src.Type() == typ.Object {
				rtk, rtv := rv.Type().Key(), rv.Type().Elem()
				src.Object().Foreach(func(k Value, v *Value) bool {
					rv.SetMapIndex(k.ToType(rtk), v.ToType(rtv))
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
		SetPrototype(&Proto.Native)
	AddGlobal("nativemap", Proto.NativeMap.ToValue())

	Proto.NativePtr = *NewNamedObject("nativeptr", 1).
		AddMethod("set", func(e *Env) {
			rv := reflect.ValueOf(e.Native(-1).Unwrap()).Elem()
			rv.Set(e.Get(0).ToType(rv.Type()))
		}).
		AddMethod("deref", func(e *Env) {
			rv := reflect.ValueOf(e.Native(-1).Unwrap())
			_ = rv.IsNil() && e.SetA(Nil) || e.SetA(ValueOf(rv.Elem().Interface()))
		}).
		SetPrototype(&Proto.Native)
	AddGlobal("nativeptr", Proto.NativePtr.ToValue())

	Proto.NativeIntf = *NewNamedObject("nativeintf", 1).
		AddMethod("deref", func(e *Env) {
			rv := reflect.ValueOf(e.Native(-1).Unwrap())
			_ = rv.IsNil() && e.SetA(Nil) || e.SetA(ValueOf(rv.Elem().Interface()))
		}).
		SetPrototype(&Proto.Native)
	AddGlobal("nativeintf", Proto.NativeIntf.ToValue())

	Proto.Channel = *Func("channel", func(e *Env) {
		e.A = ValueOf(make(chan Value, e.IntDefault(0, 0)))
	}).Object().
		AddMethod("len", func(e *Env) { e.A = Int(e.Native(-1).Len()) }).
		AddMethod("size", func(e *Env) { e.A = Int(e.Native(-1).Size()) }).
		AddMethod("close", func(e *Env) { reflect.ValueOf(e.Native(-1).Unwrap()).Close() }).
		AddMethod("send", func(e *Env) {
			rv := reflect.ValueOf(e.Native(-1).Unwrap())
			rv.Send(e.Get(0).ToType(rv.Type().Elem()))
		}).
		AddMethod("recv", func(e *Env) {
			rv := reflect.ValueOf(e.Native(-1).Unwrap())
			v, ok := rv.Recv()
			e.A = newArray(ValueOf(v), Bool(ok)).ToValue()
		}).
		SetProp("sendmulti", Func("sendmulti", func(e *Env) {
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
					Send: send.ToType(chr.Type().Elem()),
				})
				channels = append(channels, ch)
				return true
			})
			chosen, _, _ := reflect.Select(cases)
			e.A = channels[chosen]
		})).
		SetProp("recvmulti", Func("recvmulti", func(e *Env) {
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
		})).
		SetPrototype(&Proto.Native)
	AddGlobal("channel", Proto.Channel.ToValue())

	AddGlobalMethod("chr", func(e *Env) { e.A = Rune(rune(e.Int(0))) })
	AddGlobalMethod("byte", func(e *Env) { e.A = Byte(byte(e.Int(0))) })
	AddGlobalMethod("ord", func(e *Env) { r, _ := utf8.DecodeRuneInString(e.Str(0)); e.A = Int64(int64(r)) })

	Proto.Str = *Func("str", func(e *Env) {
		i := e.Get(0)
		_ = IsBytes(i) && e.SetA(Str(string(i.Native().Unwrap().([]byte)))) || e.SetA(Str(i.String()))
	}).Object().
		AddMethod("size", func(e *Env) { e.A = Int(e.Get(-1).Len()) }).
		AddMethod("len", func(e *Env) { e.A = Int(e.Get(-1).Len()) }).
		AddMethod("count", func(e *Env) { e.A = Int(utf8.RuneCountInString(e.Str(-1))) }).
		AddMethod("iequals", func(e *Env) { e.A = Bool(strings.EqualFold(e.Str(-1), e.Str(0))) }).
		AddMethod("contains", func(e *Env) { e.A = Bool(strings.Contains(e.Str(-1), e.Str(0))) }).
		AddMethod("split", func(e *Env) {
			if n := e.IntDefault(1, 0); n == 0 {
				e.A = NewNativeWithMeta(strings.Split(e.Str(-1), e.Str(0)), &Proto.StringsMeta).ToValue()
			} else {
				e.A = NewNativeWithMeta(strings.SplitN(e.Str(-1), e.Str(0), n), &Proto.StringsMeta).ToValue()
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
			start, end := e.IntDefault(1, 0), e.IntDefault(2, e.A.Len())
			e.A = Int(strings.Index(e.Str(-1)[start:end], e.Str(0)))
		}).
		AddMethod("findsub", func(e *Env) {
			s := e.Str(-1)
			idx := strings.Index(s, e.Str(0))
			_ = idx > -1 && e.SetA(Str(s[:idx])) || e.SetA(Str(""))
		}).
		AddMethod("findlast", func(e *Env) { e.A = Int(strings.LastIndex(e.Str(-1), e.Str(0))) }).
		AddMethod("trim", func(e *Env) {
			cutset := e.StrDefault(0, "", 0)
			_ = cutset == "" && e.SetA(Str(strings.TrimSpace(e.Str(-1)))) || e.SetA(Str(strings.Trim(e.Str(-1), e.Str(0))))
		}).
		AddMethod("trimprefix", func(e *Env) { e.A = Str(strings.TrimPrefix(e.Str(-1), e.Str(0))) }).
		AddMethod("trimsuffix", func(e *Env) { e.A = Str(strings.TrimSuffix(e.Str(-1), e.Str(0))) }).
		AddMethod("trimleft", func(e *Env) { e.A = Str(strings.TrimLeft(e.Str(-1), e.Str(0))) }).
		AddMethod("trimright", func(e *Env) { e.A = Str(strings.TrimRight(e.Str(-1), e.Str(0))) }).
		AddMethod("decodeutf8", func(e *Env) {
			r, sz := utf8.DecodeRuneInString(e.Str(-1))
			e.A = Array(Int64(int64(r)), Int(sz))
		}).
		AddMethod("startswith", func(e *Env) { e.A = Bool(strings.HasPrefix(e.Str(-1), e.Str(0))) }).
		AddMethod("endswith", func(e *Env) { e.A = Bool(strings.HasSuffix(e.Str(-1), e.Str(0))) }).
		AddMethod("upper", func(e *Env) { e.A = Str(strings.ToUpper(e.Str(-1))) }).
		AddMethod("lower", func(e *Env) { e.A = Str(strings.ToLower(e.Str(-1))) }).
		AddMethod("repeat", func(e *Env) { e.A = Str(strings.Repeat(e.Str(-1), e.Int(0))) }).
		AddMethod("format", func(e *Env) {
			buf := &bytes.Buffer{}
			Fprintf(buf, e.Str(-1), e.Stack()...)
			e.A = UnsafeStr(buf.Bytes())
		})

	AddGlobal("str", Proto.Str.ToValue())
}
