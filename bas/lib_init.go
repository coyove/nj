package bas

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
	"unsafe"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

const Version int64 = 379

var Globals = NewObject(0)

func init() {
	internal.GrowEnvStack = func(env unsafe.Pointer, sz int) {
		(*Env)(env).grow(sz)
	}
	internal.SetObjFun = func(obj, fun unsafe.Pointer) {
		(*Object)(obj).fun = (*Function)(fun)
		(*Function)(fun).obj = (*Object)(obj)
	}

	Globals.SetProp("VERSION", Int64(Version))
	Globals.SetMethod("globals", func(e *Env) {
		e.A = e.Global.LocalsObject().ToValue()
	}, "$f() -> object\n\tlist all global symbols and their values")
	Globals.SetMethod("new", func(e *Env) {
		m := e.Object(0)
		_ = e.Get(1).IsObject() && e.SetA(e.Object(1).SetPrototype(m).ToValue()) || e.SetA(NewObject(0).SetPrototype(m).ToValue())
	}, "$f(p: object, o?: object) -> object")

	// Debug libraries
	Globals.SetProp("debug", NamedObject("debug", 0).
		SetMethod("self", func(e *Env) {
			e.A = e.Runtime().Stack1.Callable.obj.ToValue()
		}, "$f() -> function\n\treturn caller").
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
				e.A = NewArray(r...).ToValue()
			}
		}, "$f() -> array\n\treturn [index1, name1, value1, i2, n2, v2, i3, n3, v3, ...]").
		SetMethod("globals", func(e *Env) {
			var r []Value
			for i, name := range e.Global.top.Locals {
				r = append(r, Int(i), Str(name), (*e.Global.stack)[i])
			}
			e.A = NewArray(r...).ToValue()
		}, "$f() -> array\n\treturn [index1, name1, value1, i2, n2, v2, i3, n3, v3, ...]").
		SetMethod("set", func(e *Env) {
			(*e.Global.stack)[e.Int64(0)] = e.Get(1)
		}, "$f(idx: int, v: value)").
		SetMethod("trace", func(env *Env) {
			stacks := env.Runtime().Stacktrace()
			lines := make([]Value, 0, len(stacks))
			for i := len(stacks) - 1 - env.Get(0).Safe().Int(0); i >= 0; i-- {
				r := stacks[i]
				lines = append(lines, Str(r.Callable.Name), Int64(int64(r.sourceLine())), Int64(int64(r.Cursor-1)))
			}
			env.A = NewArray(lines...).ToValue()
		}, "$f(skip?: int) -> array\n\treturn [func_name0, line1, cursor1, n2, l2, c2, ...]").
		SetMethod("disfunc", func(e *Env) {
			o := e.Object(0)
			_ = o.IsCallable() && e.SetA(Str(o.fun.GoString())) || e.SetA(Nil)
		}, "").
		SetPrototype(Proto.StaticObject).
		ToValue())

	Globals.SetMethod("type", func(e *Env) {
		e.A = Str(e.Get(0).Type().String())
	}, "$f(v: value) -> string\n\treturn `v`'s type")
	Globals.SetMethod("apply", func(e *Env) {
		e.A = CallObject(e.Object(0), e, nil, e.Get(1), e.Stack()[2:]...)
	}, "$f(f: function, this: value, args...: value) -> value")
	Globals.SetMethod("panic", func(e *Env) {
		v := e.Get(0)
		if IsPrototype(v, Proto.Error) {
			panic(v.Array().Unwrap().(*ExecError).root)
		}
		panic(v)
	}, "$f(v: value)")
	Globals.SetProp("throw", Globals.Prop("panic"))
	Globals.SetMethod("assert", func(e *Env) {
		if v := e.Get(0); e.Size() <= 1 && v.IsFalse() {
			internal.Panic("assertion failed")
		} else if e.Size() == 2 && !v.Equal(e.Get(1)) {
			internal.Panic("assertion failed: %v and %v", v, e.Get(1))
		} else if e.Size() == 3 && !v.Equal(e.Get(1)) {
			internal.Panic("%s: %v and %v", e.Get(2).String(), v, e.Get(1))
		}
	}, "$f(v: value)\n\tpanic when value is falsy\n"+
		"$f(v1: value, v2: value, msg?: string)\n\tpanic when two values are not equal")
	*Proto.Bool = *Func("bool", func(e *Env) { e.A = Bool(e.Get(0).IsTrue()) }, "$f(v: value) -> bool").Object()
	Globals.SetProp("bool", Proto.Bool.ToValue())
	*Proto.Int = *Func("int", func(e *Env) {
		if v := e.Get(0); v.Type() == typ.Number {
			e.A = Int64(v.Int64())
		} else {
			v, err := strconv.ParseInt(v.String(), e.Get(1).Safe().Int(0), 64)
			internal.PanicErr(err)
			e.A = Int64(v)
		}
	}, "$f(v: value, base?: int) -> int\n\tconvert `v` to an integer number, panic when failed or overflowed").Object()
	Globals.SetProp("int", Proto.Int.ToValue())
	*Proto.Float = *Func("float", func(e *Env) {
		if v := e.Get(0); v.Type() == typ.Number {
			e.A = v
		} else {
			f, i, isInt, err := internal.ParseNumber(v.String())
			internal.PanicErr(err)
			_ = isInt && e.SetA(Int64(i)) || e.SetA(Float64(f))
		}
	}, "$f(v: value) -> number\n\tconvert `v` to a float number, panic when failed").Object()
	Globals.SetProp("float", Proto.Float.ToValue())
	*Proto.Bytes = *Func("bytes", func(e *Env) {
		_ = e.Get(0).IsInt64() && e.SetA(ValueOf(make([]byte, e.Int(0)))) || e.SetA(ValueOf([]byte(e.Str(0))))
	}, "$f(s: str) -> bytes\n\tcreate bytes from string\n$f(n: int) -> bytes\n\tcreate an n-byte long array").
		Object().
		SetPrototype(Proto.Array)
	Globals.SetProp("bytes", Proto.Bool.ToValue())

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
		SetPrototype(Proto.StaticObject).
		ToValue())

	ObjectProto = *NamedObject("object", 0).
		SetMethod("new", func(e *Env) {
			e.A = NewObject(e.Get(0).Safe().Int(0)).ToValue()
		}, "object.$f(sz?: int) -> object").
		SetMethod("newstatic", func(e *Env) {
			e.A = NewObject(e.Get(0).Safe().Int(0)).SetPrototype(Proto.StaticObject).ToValue()
		}, "object.$f(sz?: int) -> staticobject").
		SetMethod("get", func(e *Env) {
			e.A = e.Object(-1).Get(e.Get(0))
		}, "object.$f(k: value) -> value").
		SetMethod("set", func(e *Env) {
			e.A = e.Object(-1).Set(e.Get(0), e.Get(1))
		}, "object.$f(k: value, v: value) -> value\n\tset value and return previous value").
		SetMethod("rawget", func(e *Env) {
			e.A = e.Object(-1).RawGet(e.Get(0))
		}, "object.$f(k: value) -> value").
		SetMethod("delete", func(e *Env) {
			e.A = e.Object(-1).Delete(e.Get(0))
		}, "object.$f(k: value) -> value\n\tdelete value and return previous value").
		SetMethod("clear", func(e *Env) { e.Object(-1).Clear() }, "object.$f()").
		SetMethod("copy", func(e *Env) {
			e.A = e.Object(-1).Copy(e.Get(0).IsTrue()).ToValue()
		}, "object.$f(copydata?: bool) -> object\n\tcopy the object (and its key-value data if flag is set)").
		SetMethod("proto", func(e *Env) {
			e.A = e.Object(-1).Prototype().ToValue()
		}, "object.$f() -> object\n\treturn object's prototype if any").
		SetMethod("setproto", func(e *Env) {
			e.Object(-1).SetPrototype(e.Object(0))
		}, "object.$f(p: object)\n\tset object's prototype to `p`").
		SetMethod("size", func(e *Env) {
			e.A = Int(e.Object(-1).Size())
		}, "object.$f() -> int\n\treturn the underlay size of object, which always >= object.len()").
		SetMethod("len", func(e *Env) {
			e.A = Int(e.Object(-1).Len())
		}, "object.$f() -> int\n\treturn the count of keys in object").
		SetMethod("name", func(e *Env) {
			e.A = Str(e.Object(-1).Name())
		}, "object.$f() -> string\n\treturn object's name").
		SetMethod("keys", func(e *Env) {
			a := make([]Value, 0)
			e.Object(-1).Foreach(func(k Value, v *Value) bool { a = append(a, k); return true })
			e.A = NewArray(a...).ToValue()
		}, "object.$f() -> array").
		SetMethod("values", func(e *Env) {
			a := make([]Value, 0)
			e.Object(-1).Foreach(func(k Value, v *Value) bool { a = append(a, *v); return true })
			e.A = NewArray(a...).ToValue()
		}, "object.$f() -> array").
		SetMethod("items", func(e *Env) {
			a := make([]Value, 0)
			e.Object(-1).Foreach(func(k Value, v *Value) bool { a = append(a, NewArray(k, *v).ToValue()); return true })
			e.A = NewArray(a...).ToValue()
		}, "object.$f() -> [[value, value]]\n\treturn as [[key1, value1], [key2, value2], ...]").
		SetMethod("foreach", func(e *Env) {
			f := e.Object(0)
			e.Object(-1).Foreach(func(k Value, v *Value) bool { return e.Call(f, k, *v) != False })
		}, "object.$f(f: function)").
		SetMethod("find", func(e *Env) {
			found, b := false, e.Get(0)
			e.Object(-1).Foreach(func(k Value, v *Value) bool { found = v.Equal(b); return !found })
			e.A = Bool(found)
		}, "object.$f(val: value) -> bool").
		SetMethod("contains", func(e *Env) {
			e.A = Bool(e.Object(-1).Contains(e.Get(0)))
		}, "object.$f(key: value) -> bool").
		SetMethod("merge", func(e *Env) {
			e.A = e.Object(-1).Merge(e.Object(0)).ToValue()
		}, "object.$f(o: object)\n\tmerge elements from `o`").
		SetMethod("tostring", func(e *Env) {
			p := &bytes.Buffer{}
			e.Object(-1).rawPrint(p, 0, typ.MarshalToJSON, true)
			e.A = UnsafeStr(p.Bytes())
		}, "object.$f() -> string\n\tprint raw elements in object").
		SetMethod("pure", func(e *Env) {
			m2 := e.Object(-1).Copy(false)
			e.A = m2.SetPrototype(&ObjectProto).ToValue()
		}, "object.$f() -> object\n\treturn a new object which shares all data from the original, but with no prototype").
		SetMethod("next", func(e *Env) {
			e.A = NewArray(e.Object(-1).Next(e.Get(0))).ToValue()
		}, "object.$f(k: value) -> [value, value]\n\tfind next key-value pair after `k` in the object and return as [key, value]")
	ObjectProto.SetPrototype(nil) // objectlib is the topmost object, it should not have any prototype

	*Proto.Func = *NamedObject("function", 0).
		SetMethod("doc", func(e *Env) {
			o := e.Object(-1)
			e.A = Str(strings.Replace(o.fun.DocString, "$f", o.fun.Name, -1))
		}, "function.$f() -> string\n\treturn function documentation").
		SetMethod("apply", func(e *Env) {
			e.A = CallObject(e.Object(-1), e, nil, e.Get(0), e.Stack()[1:]...)
		}, "function.$f(this: value, args...: value) -> value").
		SetMethod("call", func(e *Env) {
			e.A = e.Call(e.Object(-1), e.Stack()...)
		}, "function.$f(args...: value) -> value").
		SetMethod("try", func(e *Env) {
			a, err := e.Call2(e.Object(-1), e.Stack()...)
			_ = err == nil && e.SetA(a) || e.SetA(Error(e, err))
		}, "function.$f(args...: value) -> value|Error\n"+
			"\trun function, return Error if any panic occurred (if function tends to return n results, these values will all be Error by then)").
		SetMethod("after", func(e *Env) {
			f, args, e2 := e.Object(-1), e.CopyStack()[1:], EnvForAsyncCall(e)
			t := time.AfterFunc(e.Num(0).Safe().Duration(0), func() { e2.Call(f, args...) })
			e.A = NamedObject("Timer", 0).
				SetProp("t", ValueOf(t)).
				SetMethod("stop", func(e *Env) {
					e.A = Bool(e.This("t").(*time.Timer).Stop())
				}, "").
				ToValue()
		}, "function.$f(secs: float, args...: value) -> Timer").
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
				SetProp("w", intf(w)).
				SetMethod("wait", func(e *Env) {
					select {
					case <-time.After(e.Get(0).Safe().Duration(1 << 62)):
						panic("timeout")
					case v := <-e.This("w").(chan Value):
						e.A = v
					}
				}, "").
				ToValue()
		}, "function.$f(args...: value) -> Goroutine\n\texecute `f` in goroutine").
		SetMethod("map", func(e *Env) {
			e.A = multiMap(e, e.Object(-1), e.Get(0), e.Get(1).Safe().Int(1))
		}, "function.$f(a: object|array, n?: int) -> object|array\n"+
			"\tmap values in `a` into new values using `f(k, v)`,\n"+
			"\tsetting `n` higher than 1 will spread the load to n workers").
		SetMethod("closure", func(e *Env) {
			lambda := e.Object(-1)
			c := e.CopyStack()
			e.A = Func("<closure-"+lambda.Name()+">", func(e *Env) {
				o := e.runtime.Callable0.obj
				f := o.Prop("_l").Object()
				stk := append(o.Prop("_c").Array().Values(), e.Stack()...)
				e.A = e.Call(f, stk...)
			}, "").Object().
				SetProp("_l", lambda.ToValue()).
				SetProp("_c", NewArray(c...).ToValue()).
				ToValue()
		}, "function.$f(args...: value) -> function\n"+
			"\tcreate a closure, when it is called, `args` will be prepended to the argument list:\n"+
			"\t\tf.closure(args...)(args2...) <=> f(args..., args2...)").
		SetPrototype(&ObjectProto)

	Globals.SetProp("object", ObjectProto.ToValue())
	Globals.SetProp("staticobject", Proto.StaticObject.ToValue())
	Globals.SetProp("func", Proto.Func.ToValue())
	Globals.SetProp("callable", Proto.Func.ToValue())

	*Proto.Array = *NamedObject("array", 0).
		SetMethod("make", func(e *Env) {
			e.A = NewArray(make([]Value, e.Int(0))...).ToValue()
		}, "array.$f(n: int) -> array\n\tcreate an array of size `n`").
		SetMethod("len", func(e *Env) { e.A = Int(e.Array(-1).Len()) }, "array.$f()").
		SetMethod("size", func(e *Env) { e.A = Int(e.Array(-1).Size()) }, "array.$f()").
		SetMethod("clear", func(e *Env) { e.Array(-1).Clear() }, "array.$f()").
		SetMethod("append", func(e *Env) {
			e.Array(-1).Append(e.Stack()...)
		}, "array.$f(args...: value)\n\tappend values to array").
		SetMethod("find", func(e *Env) {
			a, src, ff := -1, e.Array(-1), e.Get(0)
			for i := 0; i < src.Len(); i++ {
				if src.Get(i).Equal(ff) {
					a = i
					break
				}
			}
			e.A = Int(a)
		}, "array.$f(v: value) -> int\n\tfind the index of first `v` in array").
		SetMethod("filter", func(e *Env) {
			src, ff := e.Array(-1), e.Object(0)
			dest := make([]Value, 0, src.Len())
			src.ForeachIndex(func(k int, v Value) bool {
				if e.Call(ff, v).IsTrue() {
					dest = append(dest, v)
				}
				return true
			})
			e.A = NewArray(dest...).ToValue()
		}, "array.$f(f: function) -> array\n\tfilter out all values where f(value) is false").
		SetMethod("slice", func(e *Env) {
			a := e.Array(-1)
			st, en := e.Int(0), e.Get(1).Safe().Int(a.Len())
			for ; st < 0 && a.Len() > 0; st += a.Len() {
			}
			for ; en < 0 && a.Len() > 0; en += a.Len() {
			}
			e.A = a.Slice(st, en).ToValue()
		}, "array.$f(start: int, end?: int) -> array\n\tslice from `start` to `end`").
		SetMethod("removeat", func(e *Env) {
			ma, idx := e.Array(-1), e.Int(0)
			if idx < 0 || idx >= ma.Len() {
				e.A = Nil
			} else {
				old := ma.Get(idx)
				ma.Copy(idx, ma.Len(), ma.Slice(idx+1, ma.Len()))
				ma.SliceInplace(0, ma.Len()-1)
				e.A = old
			}
		}, "array.$f(index: int)\n\tremove value at `index`").
		SetMethod("copy", func(e *Env) {
			e.Array(-1).Copy(e.Int(0), e.Int(1), e.Array(2))
		}, "array.$f(start: int, end: int, src: array) -> array\n\tcopy elements from `src` to `this[start:end]`").
		SetMethod("concat", func(e *Env) {
			e.Array(-1).Concat(e.Array(0))
		}, "array.$f(a: array) -> array\n\tconcat two arrays").
		SetMethod("sort", func(e *Env) {
			a, rev := e.Array(-1), e.Get(0).IsTrue()
			if kf := e.Get(1); kf.IsCallable() {
				sort.Slice(a.Unwrap(), func(i, j int) bool {
					return Less(e.Call(kf.Object(), a.Get(i)), e.Call(kf.Object(), a.Get(j))) != rev
				})
			} else {
				sort.Slice(a.Unwrap(), func(i, j int) bool { return Less(a.Get(i), a.Get(j)) != rev })
			}
		}, "array.$f(reverse?: bool, key?: function)\n\tsort array elements").
		SetMethod("istyped", func(e *Env) {
			e.A = Bool(e.Array(-1).meta != internalArrayMeta)
		}, "array.$f() -> bool").
		SetMethod("typename", func(e *Env) {
			e.A = Str(e.Array(-1).meta.Name)
		}, "array.$f() -> string").
		SetMethod("untype", func(e *Env) {
			e.A = NewArray(e.Array(-1).Values()...).ToValue()
		}, "array.$f() -> array")
	Globals.SetProp("array", Proto.Array.ToValue())

	*Proto.Error = *Func("error", func(e *Env) {
		e.A = Error(nil, &ExecError{root: e.Get(0), stacks: e.Runtime().Stacktrace()})
	}, "").Object().
		SetMethod("error", func(e *Env) { e.A = ValueOf(e.Array(-1).Unwrap().(*ExecError).root) }, "").
		SetMethod("getcause", func(e *Env) { e.A = intf(e.Array(-1).Unwrap().(*ExecError).root) }, "").
		SetMethod("trace", func(e *Env) { e.A = ValueOf(e.Array(-1).Unwrap().(*ExecError).stacks) }, "")
	Globals.SetProp("error", Proto.Error.ToValue())

	*Proto.Str = *Func("str", func(e *Env) {
		i, ok := e.Interface(0).([]byte)
		_ = ok && e.SetA(UnsafeStr(i)) || e.SetA(Str(e.Get(0).String()))
	}, "").Object().
		SetMethod("from", func(e *Env) {
			e.A = Str(fmt.Sprint(e.Interface(0)))
		}, "$f(v: value) -> string\n\tconvert `v` to string").
		SetMethod("size", func(e *Env) {
			e.A = Int(e.StrLen(-1))
		}, "str.$f() -> int\n\tsame as len()").
		SetMethod("len", func(e *Env) {
			e.A = Int(e.StrLen(-1))
		}, "str.$f() -> int\n\tsame as size()").
		SetMethod("count", func(e *Env) {
			e.A = Int(utf8.RuneCountInString(e.Str(-1)))
		}, "str.$f() -> int\n\treturn the count of runes").
		SetMethod("iequals", func(e *Env) {
			e.A = Bool(strings.EqualFold(e.Str(-1), e.Str(0)))
		}, "str.$f(text: string) -> bool\n\tcase insensitive equal").
		SetMethod("contains", func(e *Env) {
			e.A = Bool(strings.Contains(e.Str(-1), e.Str(0)))
		}, "str.$f(substr: string) -> bool").
		SetMethod("split", func(e *Env) {
			if n := e.Get(1).Safe().Int(0); n == 0 {
				e.A = NewTypedArray(strings.Split(e.Str(-1), e.Str(0)), stringsArrayMeta).ToValue()
			} else {
				e.A = NewTypedArray(strings.SplitN(e.Str(-1), e.Str(0), n), stringsArrayMeta).ToValue()
			}
		}, "str.$f(delim: string, n?: int) -> array").
		SetMethod("join", func(e *Env) {
			d := e.Str(-1)
			buf := &bytes.Buffer{}
			e.Array(0).ForeachIndex(func(k int, v Value) bool {
				buf.WriteString(v.String())
				buf.WriteString(d)
				return true
			})
			if buf.Len() > 0 {
				buf.Truncate(buf.Len() - len(d))
			}
			e.A = UnsafeStr(buf.Bytes())
		}, "str.$f(a: array) -> string").
		SetMethod("replace", func(e *Env) {
			e.A = Str(strings.Replace(e.Str(-1), e.Str(0), e.Str(1), e.Get(2).Safe().Int(-1)))
		}, "str.$f(old: string, new: string, count?: int) -> string").
		SetMethod("glob", func(e *Env) {
			m, err := filepath.Match(e.Str(-1), e.Str(0))
			internal.PanicErr(err)
			e.A = Bool(m)
		}, "str.$f(text: string) -> bool").
		SetMethod("find", func(e *Env) {
			e.A = Int(strings.Index(e.Str(-1), e.Str(0)))
		}, "str.$f(sub: string) -> int\n\tfind the index of first appearence of `sub` in text").
		SetMethod("findlast", func(e *Env) {
			e.A = Int(strings.LastIndex(e.Str(-1), e.Str(0)))
		}, "str.$f(sub: string) -> int\n\tsame as find(), but starting from right to left").
		SetMethod("sub", func(e *Env) {
			s := e.Str(-1)
			st, en := e.Int(0), e.Get(1).Safe().Int(len(s))
			for ; st < 0 && len(s) > 0; st += len(s) {
			}
			for ; en < 0 && len(s) > 0; en += len(s) {
			}
			e.A = Str(s[st:en])
		}, "str.$f(start: int, end?: int) -> string").
		SetMethod("trim", func(e *Env) {
			_ = e.Get(0).IsNil() && e.SetA(Str(strings.TrimSpace(e.Str(-1)))) || e.SetA(Str(strings.Trim(e.Str(-1), e.Str(0))))
		}, "str.$f(cutset?: string) -> string\n\ttrim spaces (or any chars in `cutset`) at both sides of the text").
		SetMethod("trimprefix", func(e *Env) {
			e.A = Str(strings.TrimPrefix(e.Str(-1), e.Str(0)))
		}, "str.$f(prefix: string) -> string\n\ttrim `prefix` of the text").
		SetMethod("trimsuffix", func(e *Env) {
			e.A = Str(strings.TrimSuffix(e.Str(-1), e.Str(0)))
		}, "str.$f(suffix: string) -> string\n\ttrim `suffix` of the text").
		SetMethod("trimleft", func(e *Env) {
			e.A = Str(strings.TrimLeft(e.Str(-1), e.Str(0)))
		}, "str.$f(cutset: string) -> string\n\ttrim the left side of the text using every char in `cutset`").
		SetMethod("trimright", func(e *Env) {
			e.A = Str(strings.TrimRight(e.Str(-1), e.Str(0)))
		}, "str.$f(cutset: string) -> string\n\ttrim the right side of the text using every char in `cutset`").
		SetMethod("ord", func(e *Env) {
			r, sz := utf8.DecodeRuneInString(e.Str(-1))
			e.A = NewArray(Int64(int64(r)), Int(sz)).ToValue()
		}, "str.$f() -> [int, int]\n\tdecode first UTF-8 char, return [unicode, width]").
		SetMethod("startswith", func(e *Env) { e.A = Bool(strings.HasPrefix(e.Str(-1), e.Str(0))) }, "str.$f(prefix: string) -> bool").
		SetMethod("endswith", func(e *Env) { e.A = Bool(strings.HasSuffix(e.Str(-1), e.Str(0))) }, "str.$f(suffix: string) -> bool").
		SetMethod("upper", func(e *Env) { e.A = Str(strings.ToUpper(e.Str(-1))) }, "str.$f() -> string").
		SetMethod("lower", func(e *Env) { e.A = Str(strings.ToLower(e.Str(-1))) }, "str.$f() -> string").
		SetMethod("chars", func(e *Env) {
			var r []Value
			for s, code := e.Str(-1), e.Get(0).IsTrue(); len(s) > 0; {
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
			e.A = NewArray(r...).ToValue()
		}, "str.$f(code?: bool) -> array\n\tbreak string into chars or unicodes, e.g.:\n"+
			"\t\t('a中c').chars() => ['a', '中', 'c']\n\t\t('a中c').chars(true) => [97, 20013, 99]").
		SetMethod("format", func(e *Env) {
			buf := &bytes.Buffer{}
			sprintf(e, -1, buf)
			e.A = UnsafeStr(buf.Bytes())
		}, "str.$f(args...: value) -> string").
		SetMethod("buffer", func(e *Env) {
			b := &bytes.Buffer{}
			if v := e.Get(0); v != Nil {
				b.WriteString(v.String())
			}
			e.A = NamedObject("Buffer", 0).
				SetPrototype(Proto.ReadWriter).
				SetProp("_f", ValueOf(b)).
				SetMethod("reset", func(e *Env) {
					e.This("_f").(*bytes.Buffer).Reset()
				}, "Buffer.$f()").
				SetMethod("value", func(e *Env) {
					e.A = UnsafeStr(e.This("_f").(*bytes.Buffer).Bytes())
				}, "Buffer.$f() -> string").
				SetMethod("bytes", func(e *Env) {
					e.A = Bytes(e.This("_f").(*bytes.Buffer).Bytes())
				}, "Buffer.$f() -> bytes").
				ToValue()
		}, "$f(v?: string) -> Buffer")
	Globals.SetProp("str", Proto.Str.ToValue())

	Globals.SetMethod("printf", func(e *Env) {
		sprintf(e, 0, e.Global.Stdout)
	}, "$f(format: string, args...: value)")
}
