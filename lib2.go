package nj

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

var (
	StrLib    Value
	MathLib   Value
	ObjectLib Value
	ArrayLib  Value
	OSLib     Value
	IOLib     Value
)

func init() {
	IOLib = Obj(
		Str("reader"), ReaderProto.ToValue(),
		Str("writer"), WriterProto.ToValue(),
		Str("seeker"), SeekerProto.ToValue(),
		Str("closer"), CloserProto.ToValue(),
		Str("readwriter"), ReadWriterProto.ToValue(),
		Str("readcloser"), ReadCloserProto.ToValue(),
		Str("writecloser"), WriteCloserProto.ToValue(),
		Str("readwritecloser"), ReadWriteCloserProto.ToValue(),
		Str("readwriteseekcloser"), ReadWriteSeekCloserProto.ToValue(),
	)
	AddGlobalValue("io", IOLib)

	ObjectLib = Func("object", func(e *Env) {
		_ = e.Get(0).IsNil() && e.SetA(Proto(e.Object(-1))) || e.SetA(e.Object(0).SetFirstParent(e.Object(-1)).ToValue())
	}).Object().Merge(nil,
		Str("make"), Func("", func(e *Env) { e.A = NewObject(e.Get(0).ToInt(0)).ToValue() }, "$f(size?: int) -> object", "\tcreate an object"),
		Str("get"), Func("", func(e *Env) { e.A = e.Object(-1).Get(e.Get(0)) }, "$f(k: value) -> value"),
		Str("set"), Func("", func(e *Env) { e.A = e.Object(-1).Set(e.Get(0), e.Get(1)) }, "$f(k: value, v: value) -> value", "\tset value and return previous value"),
		Str("rawget"), Func("", func(e *Env) { e.A = e.Object(-1).RawGet(e.Get(0)) }, "$f(k: value) -> value"),
		Str("delete"), Func("", func(e *Env) { e.A = e.Object(-1).Delete(e.Get(0)) }, "$f(k: value) -> value", "\tdelete value and return previous value"),
		Str("clear"), Func("", func(e *Env) { e.Object(-1).Clear() }, "$f()"),
		Str("copy"), Func("", func(e *Env) { e.A = e.Object(-1).Copy().ToValue() }, "$f() -> object", "\tcopy the object"),
		Str("proto"), Func("", func(e *Env) { e.A = e.Object(-1).Proto().ToValue() }, "$f() -> object", "\treturn object's parent if any"),
		Str("setproto"), Func("", func(e *Env) { e.Object(-1).SetProto(e.Object(0)) }, "$f(p: object)", "\tset object's prototype to `p`"),
		Str("setfirstproto"), Func("", func(e *Env) { e.Object(-1).SetFirstParent(e.Object(0)) }, "$f(p: object)", "\tinsert `p` as `t`'s first prototype"),
		Str("size"), Func("", func(e *Env) { e.A = Int(e.Object(-1).Size()) }, "$f() -> int", "\treturn the size of object"),
		Str("len"), Func("", func(e *Env) { e.A = Int(e.Object(-1).Len()) }, "$f() -> int", "\treturn the count of keys in object"),
		Str("keys"), Func("", func(e *Env) {
			a := make([]Value, 0)
			e.Object(-1).Foreach(func(k, v Value) bool { a = append(a, k); return true })
			e.A = Array(a...)
		}, "$f() -> array"),
		Str("values"), Func("", func(e *Env) {
			a := make([]Value, 0)
			e.Object(-1).Foreach(func(k, v Value) bool { a = append(a, v); return true })
			e.A = Array(a...)
		}, "$f() -> array"),
		Str("items"), Func("", func(e *Env) {
			a := make([]Value, 0)
			e.Object(-1).Foreach(func(k, v Value) bool { a = append(a, Array(k, v)); return true })
			e.A = Array(a...)
		}, "$f() -> [[value, value]]", "return as [[key1, value1], [key2, value2], ...]"),
		Str("foreach"), Func("", func(e *Env) {
			f := e.Object(0)
			e.Object(-1).Foreach(func(k, v Value) bool { return f.MustCall(k, v) == Nil })
		}, "$f(f: function)"),
		Str("contains"), Func("", func(e *Env) {
			found, b := false, e.Get(0)
			e.Object(-1).Foreach(func(k, v Value) bool { found = v.Equal(b); return !found })
			e.A = Bool(found)
		}, "$f(v: value) -> bool"),
		Str("merge"), Func("", func(e *Env) {
			e.A = e.Object(-1).Merge(e.Object(0)).ToValue()
		}, "$f(o: object)", "\tmerge elements from `o`"),
		Str("tostring"), Func("", func(e *Env) {
			p := &bytes.Buffer{}
			e.Object(-1).rawPrint(p, 0, true, true)
			e.A = UnsafeStr(p.Bytes())
		}, "$f() -> string", "\tprint raw elements in object"),
		Str("pure"), Func("", func(e *Env) {
			m2 := *e.Object(-1)
			m2.parent = nil
			e.A = m2.ToValue()
		}, "$f() -> object", "\treturn a new object who shares all the data from the original, but with no parent"),
		Str("next"), Func("", func(e *Env) {
			e.A = Array(e.Object(-1).Next(e.Get(0)))
		}, "$f(k: value) -> [value, value]", "\tfind next key-value pair after `k` in the object and return as [key, value]"),
		Str("iscallable"), Func("", func(e *Env) { e.A = Bool(e.Object(-1).IsCallable()) }, "$f() -> bool"),
		Str("apply"), Func("", func(e *Env) { e.A = e.Object(-1).MustApply(e.Get(0), e.Stack()[1:]...) }, "$f(this: value, args...: value) -> value"),
		Str("call"), Func("", func(e *Env) { e.A = e.Object(-1).MustCall(e.Stack()...) }, "$f(args...: value) -> value"),
		Str("try"), Func("", func(e *Env) {
			a, err := e.Object(-1).Call(e.Stack()...)
			_ = err == nil && e.SetA(a) || e.SetA(wrapExecError(err))
		}, "$f(args...: value) -> value", "\texecute, catch panic and return as error if any"),
		Str("go"), Func("GoroutineObject", func(e *Env) {
			f := e.Object(-1)
			args := e.CopyStack()
			w := make(chan Value, 1)
			go func(f *Object, args []Value) {
				if v, err := f.Call(args...); err != nil {
					w <- wrapExecError(err)
				} else {
					w <- v
				}
			}(f, args)
			e.A = Obj(
				Str("f"), f.ToValue(), Str("w"), intf(w),
				Str("stop"), Func("", func(e *Env) {
					e.Object(-1).Prop("f").Object().Callable.EmergStop()
				}),
				Str("wait"), Func("", func(e *Env) {
					ch := e.Object(-1).Prop("w").Interface().(chan Value)
					if w := e.Get(0).ToFloat64(0); w > 0 {
						select {
						case <-time.After(time.Duration(w * float64(time.Second))):
							panic("timeout")
						case v := <-ch:
							e.A = v
						}
					} else {
						e.A = <-ch
					}
				}),
			)
		}, "$f(args...: value)", "\texecute `f` in goroutine"),
	).ToValue()
	AddGlobalValue("object", ObjectLib)

	ArrayLib = Func("arary", nil).Object().Merge(nil,
		Str("make"), Func("", func(e *Env) { e.A = Array(make([]Value, e.Int(0))...) }, "$f(n: int) -> array", "\tcreate an array of size `n`"),
		Str("len"), Func("", func(e *Env) { e.A = Int(e.Array(-1).Len()) }, "$f()"),
		Str("size"), Func("", func(e *Env) { e.A = Int(e.Array(-1).Size()) }, "$f()"),
		Str("clear"), Func("", func(e *Env) { e.Array(-1).Clear() }, "$f()"),
		Str("append"), Func("", func(e *Env) { e.Array(-1).Append(e.Stack()...) }, "$f(args...: value)", "\tappend values to array"),
		Str("find"), Func("", func(e *Env) {
			e.A = Int(-1)
			src, ff := e.Array(-1), e.Get(0)
			for i := 0; i < src.Len(); i++ {
				if src.Get(i).Equal(ff) {
					e.A = Int(i)
					break
				}
			}
		}, "$f(v: value) -> int", "\tfind the index of first `v` in array"),
		Str("filter"), Func("", func(e *Env) {
			src, ff := e.Array(-1), e.Object(0)
			dest := make([]Value, 0, src.Len())
			src.Foreach(func(k int, v Value) bool {
				if ff.MustCall(v).IsTrue() {
					dest = append(dest, v)
				}
				return true
			})
			e.A = Array(dest...)
		}, "$f(f: function) -> array", "\tfilter out all values where f(value) is false"),
		Str("slice"), Func("", func(e *Env) {
			e.A = e.Array(-1).Slice(e.Int(0), e.Get(1).ToInt(e.Array(-1).Len())).ToValue()
		}, "$f(start: int, end?: int) -> array", "\tslice array, from `start` to `end`"),
		Str("removeat"), Func("", func(e *Env) {
			ma, idx := e.Array(-1), e.Int(0)
			if idx < 0 || idx >= ma.Len() {
				e.A = Nil
			} else {
				old := ma.Get(idx)
				ma.Slice(0, idx).Concat(ma.Slice(idx+1, ma.Len()))
				e.A = old
			}
		}, "$f(index: int)", "\tremove value at `index`"),
		Str("copy"), Func("", func(e *Env) {
			a := e.Array(-1)
			a.Copy(e.Int(0), e.Int(1), e.Array(2))
			e.A = a.ToValue()
		}, "$f(start: int, end: int, src: array) -> array", "\tcopy elements from `src` to `this[start:end]`"),
		Str("concat"), Func("", func(e *Env) {
			a := e.Array(-1)
			a.Concat(e.Array(0))
			e.A = a.ToValue()
		}, "$f(array2: array) -> array", "\tconcat two arrays"),
	).ToValue()
	AddGlobalValue("array", ArrayLib)

	encDecProto := Func("EncodeDecode", nil).Object().Merge(nil,
		Str("encode"), Func("", func(e *Env) {
			e.A = Str(e.Object(-1).Prop("_e").Interface().(interface {
				EncodeToString([]byte) string
			}).EncodeToString(e.Get(0).ToBytes()))
		}),
		Str("decode"), Func("", func(e *Env) {
			v, err := e.Object(-1).Prop("_e").Interface().(interface {
				DecodeString(string) ([]byte, error)
			}).DecodeString(e.Str(0))
			internal.PanicErr(err)
			e.A = Bytes(v)
		}),
	).SetProto(Func("EncoderDecoder", nil).Object().Merge(nil,
		Str("encoder"), Func("", func(e *Env) {
			enc := Nil
			buf := &bytes.Buffer{}
			switch encoding := e.Object(-1).Prop("_e").Interface().(type) {
			default:
				enc = ValueOf(hex.NewEncoder(buf))
			case *base32.Encoding:
				enc = ValueOf(base32.NewEncoder(encoding, buf))
			case *base64.Encoding:
				enc = ValueOf(base64.NewEncoder(encoding, buf))
			}
			e.A = Proto(WriteCloserProto,
				Str("_f"), ValueOf(enc),
				Str("_b"), ValueOf(buf),
				Str("value"), Func("", func(e *Env) {
					e.A = Str(e.Object(-1).Prop("_b").Interface().(*bytes.Buffer).String())
				}),
				Str("bytes"), Func("", func(e *Env) {
					e.A = Bytes(e.Object(-1).Prop("_b").Interface().(*bytes.Buffer).Bytes())
				}),
			)
		}),
		Str("decoder"), Func("", func(e *Env) {
			src := NewReader(e.Get(0))
			dec := Nil
			switch encoding := e.Object(-1).Prop("_e").Interface().(type) {
			case *base64.Encoding:
				dec = ValueOf(base64.NewDecoder(encoding, src))
			case *base32.Encoding:
				dec = ValueOf(base32.NewDecoder(encoding, src))
			default:
				dec = ValueOf(hex.NewDecoder(src))
			}
			e.A = Proto(ReaderProto, Str("_f"), ValueOf(dec))
		}),
	))

	StrLib = Func("String", func(e *Env) {
		i, ok := e.Interface(0).([]byte)
		_ = ok && e.SetA(UnsafeStr(i)) || e.SetA(Str(e.Get(0).String()))
	}).Object().Merge(nil,
		Str("from"), Func("", func(e *Env) { e.A = Str(fmt.Sprint(e.Interface(0))) }, "$f(v: value) -> string", "\tconvert `v` to string"),
		Str("size"), Func("", func(e *Env) { e.A = Int(e.StrLen(-1)) }, "$f() -> int"),
		Str("len"), Func("", func(e *Env) { e.A = Int(e.StrLen(-1)) }, "$f() -> int"),
		Str("count"), Func("", func(e *Env) { e.A = Int(utf8.RuneCountInString(e.Str(-1))) }, "$f() -> int", "\treturn count of runes in text"),
		Str("iequals"), Func("", func(e *Env) { e.A = Bool(strings.EqualFold(e.Str(-1), e.Str(0))) }, "$f(text2: string) -> bool"),
		Str("contains"), Func("", func(e *Env) { e.A = Bool(strings.Contains(e.Str(-1), e.Str(0))) }, "$f(substr: string) -> bool"),
		Str("split"), Func("", func(e *Env) {
			s, d := e.Str(-1), e.Str(0)
			var r []Value
			if n := e.Get(1).ToInt64(0); n == 0 {
				for _, p := range strings.Split(s, d) {
					r = append(r, Str(p))
				}
			} else {
				for _, p := range strings.SplitN(s, d, int(n)) {
					r = append(r, Str(p))
				}
			}
			e.A = Array(r...)
		}, "$f(delim: string, n?: int) -> array"),
		Str("join"), Func("", func(e *Env) {
			d := e.Str(-1)
			buf := &bytes.Buffer{}
			e.Array(0).Foreach(func(k int, v Value) bool {
				buf.WriteString(v.String())
				buf.WriteString(d)
				return true
			})
			if buf.Len() > 0 {
				buf.Truncate(buf.Len() - len(d))
			}
			e.A = UnsafeStr(buf.Bytes())
		}, "$f(a: array) -> string"),
		Str("replace"), Func("", func(e *Env) {
			e.A = Str(strings.Replace(e.Str(-1), e.Str(0), e.Str(1), e.Get(2).ToInt(-1)))
		}, "$f(old: string, new: string, count?: int) -> string"),
		Str("match"), Func("", func(e *Env) {
			m, err := filepath.Match(e.Str(-1), e.Str(0))
			internal.PanicErr(err)
			e.A = Bool(m)
		}, "$f(text: string) -> bool"),
		Str("find"), Func("", func(e *Env) {
			e.A = Int(strings.Index(e.Str(-1), e.Str(0)))
		}, "$f(sub: string) -> int", "\tfind the index of first appearence of `sub` in text"),
		Str("findlast"), Func("", func(e *Env) {
			e.A = Int(strings.LastIndex(e.Str(-1), e.Str(0)))
		}, "$f(sub: string) -> int", "\tsame as find(), but starting from right to left"),
		Str("sub"), Func("", func(e *Env) {
			s := e.Str(-1)
			st, en := e.Int(0), e.Get(1).ToInt(len(s))
			for ; st < 0 && len(s) > 0; st += len(s) {
			}
			for ; en < 0 && len(s) > 0; en += len(s) {
			}
			e.A = Str(s[st:en])
		}, "$f(start: int, end?: int) -> string"),
		Str("trim"), Func("", func(e *Env) {
			_ = e.Get(0).IsNil() && e.SetA(Str(strings.TrimSpace(e.Str(-1)))) || e.SetA(Str(strings.Trim(e.Str(-1), e.Str(0))))
		}, "$f(cutset?: string) -> string", "\ttrim spaces (or any chars in `cutset`) at both sides of the text"),
		Str("trimprefix"), Func("", func(e *Env) {
			e.A = Str(strings.TrimPrefix(e.Str(-1), e.Str(0)))
		}, "$f(prefix: string) -> string", "\ttrim `prefix` of the text"),
		Str("trimsuffix"), Func("", func(e *Env) {
			e.A = Str(strings.TrimSuffix(e.Str(-1), e.Str(0)))
		}, "$f(suffix: string) -> string", "\ttrim `suffix` of the text"),
		Str("trimleft"), Func("", func(e *Env) {
			e.A = Str(strings.TrimLeft(e.Str(-1), e.Str(0)))
		}, "$f(cutset: string) -> string", "\ttrim the left side of the text using every char in `cutset`"),
		Str("trimright"), Func("", func(e *Env) {
			e.A = Str(strings.TrimRight(e.Str(-1), e.Str(0)))
		}, "$f(cutset: string) -> string", "\ttrim the right side of the text using every char in `cutset`"),
		Str("ord"), Func("", func(e *Env) {
			r, sz := utf8.DecodeRuneInString(e.Str(-1))
			e.A = Array(Int64(int64(r)), Int(sz))
		}, "$f() -> [int, int]", "\tdecode first UTF-8 char, return [unicode, bytescount]"),
		Str("startswith"), Func("", func(e *Env) { e.A = Bool(strings.HasPrefix(e.Str(-1), e.Str(0))) }, "$f(prefix: string) -> bool"),
		Str("endswith"), Func("", func(e *Env) { e.A = Bool(strings.HasSuffix(e.Str(-1), e.Str(0))) }, "$f(suffix: string) -> bool"),
		Str("upper"), Func("", func(e *Env) { e.A = Str(strings.ToUpper(e.Str(-1))) }, "$f() -> string"),
		Str("lower"), Func("", func(e *Env) { e.A = Str(strings.ToLower(e.Str(-1))) }, "$f() -> string"),
		Str("bytes"), Func("", func(e *Env) {
			_ = e.Get(0).IsInt64() && e.SetA(ValueOf(make([]byte, e.Int(0)))) || e.SetA(ValueOf([]byte(e.Str(0))))
		}, "$f() -> bytes", "\tconvert text to byte array",
			"$f(n: int) -> bytes", "\tcreate an n-byte long array"),
		Str("chars"), Func("", func(e *Env) {
			var r []Value
			for s := e.Str(-1); len(s) > 0; {
				_, sz := utf8.DecodeRuneInString(s)
				if sz == 0 {
					break
				}
				r = append(r, Str(s[:sz]))
				s = s[sz:]
			}
			e.A = Array(r...)
		}, "$f() -> array", "\tbreak `text` into chars, e.g.: chars('a中c') => ['a', '中', 'c']"),
		Str("format"), Func("", func(e *Env) {
			buf := &bytes.Buffer{}
			sprintf(e, -1, buf)
			e.A = UnsafeStr(buf.Bytes())
		}, "$f(args...: value) -> string"),
		Str("buffer"), Func("", func(e *Env) {
			b := &bytes.Buffer{}
			if v := e.Get(0); v != Nil {
				b.WriteString(v.String())
			}
			e.A = Func("Buffer", nil).Object().SetProto(ReadWriterProto).Merge(nil,
				Str("_f"), ValueOf(b),
				Str("reset"), Func("", func(e *Env) {
					e.Object(-1).Prop("_f").Interface().(*bytes.Buffer).Reset()
				}),
				Str("value"), Func("", func(e *Env) {
					e.A = UnsafeStr(e.Object(-1).Prop("_f").Interface().(*bytes.Buffer).Bytes())
				}),
				Str("bytes"), Func("", func(e *Env) {
					e.A = Bytes(e.Object(-1).Prop("_f").Interface().(*bytes.Buffer).Bytes())
				}),
			).ToValue()
		}),
		Str("hex"), Func("hex", nil).Object().SetProto(encDecProto.Proto()).ToValue(),
		Str("base64"), Func("base64", nil).Object().Merge(nil,
			Str("std"), Proto(encDecProto, Str("_e"), ValueOf(getEncB64(base64.StdEncoding, '='))),
			Str("url"), Proto(encDecProto, Str("_e"), ValueOf(getEncB64(base64.URLEncoding, '='))),
			Str("std2"), Proto(encDecProto, Str("_e"), ValueOf(getEncB64(base64.StdEncoding, -1))),
			Str("url2"), Proto(encDecProto, Str("_e"), ValueOf(getEncB64(base64.URLEncoding, -1))),
		).SetProto(encDecProto).ToValue(),
		Str("base32"), Func("base32", nil).Object().Merge(nil,
			Str("std"), Proto(encDecProto, Str("_e"), ValueOf(getEncB32(base32.StdEncoding, '='))),
			Str("hex"), Proto(encDecProto, Str("_e"), ValueOf(getEncB32(base32.HexEncoding, '='))),
			Str("std2"), Proto(encDecProto, Str("_e"), ValueOf(getEncB32(base32.StdEncoding, -1))),
			Str("hex2"), Proto(encDecProto, Str("_e"), ValueOf(getEncB32(base32.HexEncoding, -1))),
		).SetProto(encDecProto).ToValue(),
	).ToValue()
	AddGlobalValue("str", StrLib)

	MathLib = Obj(
		Str("INF"), Float64(math.Inf(1)),
		Str("NEG_INF"), Float64(math.Inf(-1)),
		Str("PI"), Float64(math.Pi),
		Str("E"), Float64(math.E),
		Str("randomseed"), Func("", func(e *Env) { rand.Seed(e.B(0).ToInt64(1)) }, "$f(seed: int)"),
		Str("random"), Func("", func(e *Env) {
			switch len(e.Stack()) {
			case 2:
				ai, bi := e.Int64(0), e.Int64(1)
				e.A = Int64(rand.Int63n(bi-ai+1) + ai)
			case 1:
				e.A = Int64(rand.Int63n(e.Int64(0)))
			default:
				e.A = Float64(rand.Float64())
			}
		}, "$f() -> float", "\treturn [0, 1)",
			"$f(n: int) -> int", "\treturn [0, n)",
			"$f(a: int, b: int) -> int", "\treturn [a, b]"),
		Str("sqrt"), Func("", func(e *Env) { e.A = Float64(math.Sqrt(e.Float64(0))) }),
		Str("floor"), Func("", func(e *Env) { e.A = Float64(math.Floor(e.Float64(0))) }),
		Str("ceil"), Func("", func(e *Env) { e.A = Float64(math.Ceil(e.Float64(0))) }),
		Str("min"), Func("", func(e *Env) { mathMinMax(e, false) }, "$f(a: number, b...: number) -> number"),
		Str("max"), Func("", func(e *Env) { mathMinMax(e, true) }, "$f(a: number, b...: number) -> number"),
		Str("pow"), Func("", func(e *Env) {
			e.A = Float64(math.Pow(e.Float64(0), e.Float64(1)))
		}, "$f(a: float, b: float) -> float"),
		Str("abs"), Func("", func(e *Env) {
			if e.A = e.Num(0); e.A.IsInt64() {
				if i := e.A.Int64(); i < 0 {
					e.A = Int64(-i)
				}
			} else {
				e.A = Float64(math.Abs(e.A.Float64()))
			}
		}),
		Str("remainder"), Func("", func(e *Env) { e.A = Float64(math.Remainder(e.Float64(0), e.Float64(1))) }),
		Str("mod"), Func("", func(e *Env) { e.A = Float64(math.Mod(e.Float64(0), e.Float64(1))) }),
		Str("cos"), Func("", func(e *Env) { e.A = Float64(math.Cos(e.Float64(0))) }),
		Str("sin"), Func("", func(e *Env) { e.A = Float64(math.Sin(e.Float64(0))) }),
		Str("tan"), Func("", func(e *Env) { e.A = Float64(math.Tan(e.Float64(0))) }),
		Str("acos"), Func("", func(e *Env) { e.A = Float64(math.Acos(e.Float64(0))) }),
		Str("asin"), Func("", func(e *Env) { e.A = Float64(math.Asin(e.Float64(0))) }),
		Str("atan"), Func("", func(e *Env) { e.A = Float64(math.Atan(e.Float64(0))) }),
		Str("atan2"), Func("", func(e *Env) { e.A = Float64(math.Atan2(e.Float64(0), e.Float64(1))) }),
		Str("ldexp"), Func("", func(e *Env) { e.A = Float64(math.Ldexp(e.Float64(0), e.Int(0))) }),
		Str("modf"), Func("", func(e *Env) { a, b := math.Modf(e.Float64(0)); e.A = Array(Float64(a), Float64(b)) }),
	)
	AddGlobalValue("math", MathLib)

	OSLib = Obj(
		Str("args"), ValueOf(os.Args),
		Str("environ"), Func("", func(e *Env) { e.A = ValueOf(os.Environ()) }),
		Str("shell"), Func("", func(e *Env) {
			p := exec.Command("sh", "-c", e.Str(0))
			opt := e.Get(1)
			timeout := time.Duration(1 << 62) // basically forever
			if tmp := opt.ToObject().Prop("timeout"); tmp != Nil {
				timeout = time.Duration(tmp.Is(typ.Number, "timeout").Float64() * float64(time.Second))
			}
			if tmp := opt.ToObject().Prop("env"); tmp != Nil {
				tmp.Is(typ.Object, "env").Object().Foreach(func(k, v Value) bool {
					p.Env = append(p.Env, k.String()+"="+v.String())
					return true
				})
			}
			stdout := &bytes.Buffer{}
			p.Stdout, p.Stderr = stdout, stdout
			p.Dir = opt.ToObject().Prop("dir").ToStr("")
			if tmp := opt.ToObject().Prop("stdout"); tmp != Nil {
				p.Stdout = NewWriter(tmp)
			}
			if tmp := opt.ToObject().Prop("stderr"); tmp != Nil {
				p.Stderr = NewWriter(tmp)
			}
			if tmp := opt.ToObject().Prop("stdin"); tmp != Nil {
				p.Stdin = NewReader(tmp)
			}

			out := make(chan error)
			go func() { out <- p.Run() }()
			select {
			case r := <-out:
				internal.PanicErr(r)
			case <-time.After(timeout):
				p.Process.Kill()
				panic("timeout")
			}
			e.A = Bytes(stdout.Bytes())
		}),
		Str("readdir"), Func("", func(e *Env) {
			fi, err := ioutil.ReadDir(e.Str(0))
			internal.PanicErr(err)
			e.A = ValueOf(fi)
		}),
		Str("remove"), Func("", func(e *Env) {
			path := e.Str(0)
			fi, err := os.Stat(path)
			internal.PanicErr(err)
			if fi.IsDir() {
				internal.PanicErr(os.RemoveAll(path))
			} else {
				internal.PanicErr(os.Remove(path))
			}
		}),
		Str("pstat"), Func("", func(e *Env) {
			fi, err := os.Stat(e.Str(0))
			_ = err == nil && e.SetA(ValueOf(fi)) || e.SetA(Nil)
		}),
	)
	AddGlobalValue("os", OSLib)
}

func mathMinMax(e *Env, max bool) {
	if v := e.Num(0); v.IsInt64() {
		vi := v.Int64()
		for ii := 1; ii < len(e.Stack()); ii++ {
			if x := e.Int64(ii); x >= vi == max {
				vi = x
			}
		}
		e.A = Int64(vi)
	} else {
		vf := v.Float64()
		for i := 1; i < len(e.Stack()); i++ {
			if x := e.Float64(i); x >= vf == max {
				vf = x
			}
		}
		e.A = Float64(vf)
	}
}

func sprintf(env *Env, start int, p io.Writer) {
	f := env.Str(start)
	tmp := bytes.Buffer{}
	popi := start
	for len(f) > 0 {
		idx := strings.Index(f, "%")
		if idx == -1 {
			fmt.Fprint(p, f)
			break
		} else if idx == len(f)-1 {
			internal.Panic("unexpected '%%' at end")
		}
		fmt.Fprint(p, f[:idx])
		if f[idx+1] == '%' {
			p.Write([]byte("%"))
			f = f[idx+2:]
			continue
		}
		tmp.Reset()
		tmp.WriteByte('%')
		expecting := typ.Nil
		for f = f[idx+1:]; len(f) > 0 && expecting == typ.Nil; {
			switch f[0] {
			case 'b', 'd', 'o', 'O', 'c', 'e', 'E', 'f', 'F', 'g', 'G':
				expecting = typ.Number
			case 's', 'q', 'U':
				expecting = typ.String
			case 'x', 'X':
				expecting = typ.String + typ.Number
			case 'v':
				expecting = typ.Native
			case 't':
				expecting = typ.Bool
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', '-', '+', '#', ' ':
			default:
				internal.Panic("unexpected verb: '%c'", f[0])
			}
			tmp.WriteByte(f[0])
			f = f[1:]
		}

		popi++
		switch expecting {
		case typ.Bool:
			fmt.Fprint(p, env.Bool(popi))
		case typ.String:
			if pop := env.Get(popi); pop.IsBytes() {
				fmt.Fprintf(p, tmp.String(), pop.Array().Unwrap())
			} else {
				fmt.Fprintf(p, tmp.String(), pop.String())
			}
		case typ.Number + typ.String:
			if pop := env.Get(popi); pop.Type() == typ.String {
				fmt.Fprintf(p, tmp.String(), pop.Str())
				continue
			} else if pop.IsBytes() {
				fmt.Fprintf(p, tmp.String(), pop.Array().Unwrap())
				continue
			}
			fallthrough
		case typ.Number:
			if pop := env.Num(popi); pop.IsInt64() {
				fmt.Fprintf(p, tmp.String(), pop.Int64())
			} else {
				fmt.Fprintf(p, tmp.String(), pop.Float64())
			}
		case typ.Native:
			fmt.Fprint(p, env.Interface(popi))
		}
	}
}

func getEncB64(enc *base64.Encoding, padding rune) *base64.Encoding {
	if padding != '=' {
		enc = enc.WithPadding(padding)
	}
	return enc
}

func getEncB32(enc *base32.Encoding, padding rune) *base32.Encoding {
	if padding != '=' {
		enc = enc.WithPadding(padding)
	}
	return enc
}
