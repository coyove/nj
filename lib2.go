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
	"sync"
	"time"
	"unicode/utf8"
	"unsafe"

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
	IOLib = TableMerge(IOLib, Obj(
		Str("reader"), ReaderProto.Value(),
		Str("writer"), WriterProto.Value(),
		Str("seeker"), SeekerProto.Value(),
		Str("closer"), CloserProto.Value(),
		Str("readwriter"), ReadWriterProto.Value(),
		Str("readcloser"), ReadCloserProto.Value(),
		Str("writecloser"), WriteCloserProto.Value(),
		Str("readwritecloser"), ReadWriteCloserProto.Value(),
		Str("readwriteseekcloser"), ReadWriteSeekCloserProto.Value(),
	).Object())
	AddGlobalValue("io", IOLib)

	ObjectLib = TableMerge(ObjectLib, Func("object", func(e *Env) {
		if e.Get(0) == Nil {
			e.A = Proto(e.Object(-1))
		} else {
			e.A = e.Object(0).SetFirstParent(e.Object(-1)).Value()
		}
	}).Object().Merge(nil,
		Str("concurrent"), Func("", func(e *Env) {
			x := NewObject(e.Object(-1).Len())
			ObjectLib.Object().Foreach(func(k, v Value) bool {
				if v.IsObject() {
					if old := v.Object(); old.IsCallable() {
						v = Func(old.callable.Name, func(e *Env) {
							mu := e.Object(-1).Gets("_mu").Interface().(*sync.Mutex)
							mu.Lock()
							defer mu.Unlock()
							e.A = old.MustCall(e.Stack()...)
						}, old.callable.DocString)
					}
				}
				x.Set(k, v)
				return true
			})
			x.Sets("_mu", Val(&sync.Mutex{}))
			_ = e.Get(1).IsNil() && e.SetA(Proto(x)) || e.SetA(e.Object(1).SetFirstParent(x).Value())
		}, "$f() -> object", "\tcreate a concurrently accessible object"),
		Str("make"), Func("", func(e *Env) { e.A = NewObject(e.Get(0).ToInt(0)).Value() }, "$f(n: int) -> object", "\tcreate an object"),
		Str("get"), Func("", func(e *Env) { e.A = e.Object(-1).Get(e.Get(0)) }, "$f(k: value) -> value"),
		Str("set"), Func("", func(e *Env) { e.A = e.Object(-1).Set(e.Get(0), e.Get(1)) }, "$f(k: value, v: value) -> value", "\tset value and return previous value"),
		Str("rawget"), Func("", func(e *Env) { e.A = e.Object(-1).RawGet(e.Get(0)) }, "$f(k: value) -> value"),
		Str("clear"), Func("", func(e *Env) { e.Object(-1).Clear() }, "$f()"),
		Str("copy"), Func("", func(e *Env) { e.A = e.Object(-1).Copy().Value() }, "$f() -> object", "\tcopy the object"),
		Str("parent"), Func("", func(e *Env) { e.A = e.Object(-1).Parent().Value() }, "$f() -> object", "\treturn object's parent if any"),
		Str("setparent"), Func("", func(e *Env) { e.Object(-1).SetParent(e.Object(0)) }, "$f(p: table)", "\tset object's parent"),
		Str("setfirstparent"), Func("", func(e *Env) { e.Object(-1).SetFirstParent(e.Object(0)) }, "$f(p: table)", "\tinsert `p` as `t`'s first parent"),
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
		}, "$f() -> array", "return as [[key1, value1], [key2, value2], ...]"),
		Str("foreach"), Func("", func(e *Env) {
			f := e.Object(0)
			e.Object(-1).Foreach(func(k, v Value) bool { return f.MustCall(k, v) == Nil })
		}, "$f(f: function)"),
		Str("contains"), Func("", func(e *Env) {
			found, b := false, e.Get(0)
			e.Object(-1).Foreach(func(k, v Value) bool { found = v.Equal(b); return !found })
			e.A = Bool(found)
		}, "$f(v: value) -> bool"),
		Str("merge"), Func2("", func(a, b Value) Value {
			return a.MustTable("").Merge(b.MustTable("")).Value()
		}, "$f({table1}: table, table2: table)", "\tmerge elements from table2 to table1"),
		Str("tostring"), Func("", func(e *Env) {
			p := &bytes.Buffer{}
			e.Object(-1).rawPrint(p, 0, true, true)
			e.A = Bytes(p.Bytes())
		}, "$f() -> string", "\tprint raw elements in table"),
		Str("pure"), Func1("", func(m Value) Value {
			m2 := *m.MustTable("")
			m2.parent = nil
			return m2.Value()
		}, "$f({t}: table) -> object", "\treturn a new table who shares all the data from t, but with no parent"),
		Str("unwrap"), Func("", func(e *Env) {
			v := e.Get(0)
			_ = v.Type() == typ.Native && e.SetA(ValRec(v.Interface())) || e.SetA(v)
		}, "$f(v: value) -> object", "\tunwrap Go's array, slice or map into object"),
		Str("next"), Func("", func(e *Env) {
			e.A = Array(e.Object(-1).Next(e.Get(0)))
		}, "$f(k: value) -> array", "\tfind next key-value pair after `k` in the object and return as [key, value]"),
	))
	AddGlobalValue("table", ObjectLib)

	ArrayLib = TableMerge(ArrayLib, Obj(
		Str("make"), Func("", func(e *Env) { e.A = Array(make([]Value, e.Int(0))...) }, "$f(n: int) -> array", "\tcreate an array of size `n`"),
		Str("copy"), Func("", func(e *Env) { e.A = Array(append([]Value{}, e.Array(-1).store...)...) }, "$f() -> array", "\t copy the array"),
		Str("len"), Func("", func(e *Env) { e.A = Int(e.Array(-1).Len()) }, "$f()"),
		Str("size"), Func("", func(e *Env) { e.A = Int(e.Array(-1).Size()) }, "$f()"),
		Str("clear"), Func("", func(e *Env) { e.Array(-1).Clear() }, "$f()"),
		Str("append"), Func("", func(e *Env) {
			ma := e.Array(-1)
			ma.store = append(ma.store, e.Stack()...)
		}, "$f(args...: value)", "\tappend values to array"),
		Str("find"), Func("", func(e *Env) {
			e.A = Int(-1)
			src, ff := e.Array(-1), e.Get(0)
			for i, v := range src.store {
				if v.Equal(ff) {
					e.A = Int(i)
					break
				}
			}
		}, "$f(v: value) -> int", "\tfind the index of first `v` in array"),
		Str("filter"), Func("", func(e *Env) {
			src, ff := e.Array(-1), e.Object(0)
			dest := make([]Value, 0, src.Len())
			src.Foreach(func(k, v Value) bool {
				if MustValue(ff.Call(v)).IsTrue() {
					dest = append(dest, v)
				}
				return true
			})
			e.A = Array(dest...)
		}, "$f(f: function) -> array", "\tfilter out all values where f(value) is false"),
		Str("slice"), Func("", func(e *Env) {
			e.A = Array(e.Array(-1).store[e.Int(0):e.Get(1).ToInt(e.Array(-1).Len())]...)
		}, "$f(start: int, end?: int) -> array", "\tslice array, from start to end"),
		Str("removeat"), Func("", func(e *Env) {
			ma, idx := e.Array(-1), e.Int(0)
			if idx < 0 || idx >= ma.Len() {
				e.A = Nil
			} else {
				old := ma.store[idx]
				ma.store = append(ma.store[:idx], ma.store[idx+1:]...)
				e.A = old
			}
		}, "$f(index: int)", "\tremove value at `index`"),
		Str("concat"), Func("", func(e *Env) {
			ma := e.Array(-1)
			ma.store = append(ma.store, e.Array(0).store...)
			e.A = ma.Value()
		}, "$f(array2: array)", "\tconcat two arrays"),
	).Object())
	AddGlobalValue("array", ArrayLib)

	encDecProto := Proto(Obj(
		Str("encoder"), Func("", func(e *Env) {
			enc := Nil
			buf := &bytes.Buffer{}
			switch encoding := e.Object(-1).Gets("_e").Interface().(type) {
			default:
				enc = Val(hex.NewEncoder(buf))
			case *base32.Encoding:
				enc = Val(base32.NewEncoder(encoding, buf))
			case *base64.Encoding:
				enc = Val(base64.NewEncoder(encoding, buf))
			}
			e.A = Proto(WriteCloserProto,
				Str("_f"), Val(enc),
				Str("_b"), Val(buf),
				Str("value"), Func("", func(e *Env) {
					e.A = Bytes(e.Object(-1).Gets("_b").Interface().(*bytes.Buffer).Bytes())
				}),
			)
		}),
		Str("decoder"), Func("", func(e *Env) {
			src := NewReader(e.Get(0))
			dec := Nil
			switch encoding := e.Object(-1).Gets("_e").Interface().(type) {
			case *base64.Encoding:
				dec = Val(base64.NewDecoder(encoding, src))
			case *base32.Encoding:
				dec = Val(base32.NewDecoder(encoding, src))
			default:
				dec = Val(hex.NewDecoder(src))
			}
			e.A = Proto(ReaderProto, Str("_f"), Val(dec))
		}),
	).Object(),
		Str("__name"), Str("encdecfast"),
		Str("encode"), Func("", func(e *Env) {
			x := struct {
				a string
				c int
			}{a: e.Str(0), c: e.StrLen(0)}
			e.A = Str(e.Object(-1).Gets("_e").Interface().(interface {
				EncodeToString([]byte) string
			}).EncodeToString(*(*[]byte)(unsafe.Pointer(&x))))
		}),
		Str("decode"), Func("", func(e *Env) {
			v, err := e.Object(-1).Gets("_e").Interface().(interface {
				DecodeString(string) ([]byte, error)
			}).DecodeString(e.Str(0))
			internal.PanicErr(err)
			e.A = Bytes(v)
		}),
	)

	StrLib = TableMerge(StrLib, Func("String", func(e *Env) {
		i, ok := e.Interface(1).([]byte)
		_ = ok && e.SetA(Bytes(i)) || e.SetA(Str(e.Get(1).String()))
	}).Object().Merge(nil,
		Str("size"), Func("", func(e *Env) { e.A = Int(e.StrLen(-1)) }, "$f() -> int"),
		Str("len"), Func("", func(e *Env) { e.A = Int(e.StrLen(-1)) }, "$f() -> int"),
		Str("count"), Func("", func(e *Env) { e.A = Int(utf8.RuneCountInString(e.Str(-1))) }, "$f() -> int", "\treturn count of runes in text"),
		Str("from"), Func("", func(e *Env) { e.A = Str(fmt.Sprint(e.Interface(0))) }, "$f(v: value) -> string", "\tconvert value to string"),
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
			for _, a := range e.Array(0).store {
				buf.WriteString(a.String())
				buf.WriteString(d)
			}
			if buf.Len() > 0 {
				buf.Truncate(buf.Len() - len(d))
			}
			e.A = Bytes(buf.Bytes())
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
			_ = e.Get(0).IsInt64() && e.SetA(Val(make([]byte, e.Int(0)))) || e.SetA(Val([]byte(e.Str(0))))
		}, "$f() -> go.[]byte", "\tconvert text to []byte",
			"$f(n: int) -> go.[]byte", "\tcreate an n-byte long array"),
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
			e.A = Bytes(buf.Bytes())
		}, "$f(args...: value) -> string"),
		Str("buffer"), Func("", func(e *Env) {
			b := &bytes.Buffer{}
			if v := e.Get(0); v != Nil {
				b.WriteString(v.String())
			}
			e.A = Func("Buffer", nil).Object().SetParent(ReadWriterProto).Merge(nil,
				Str("_f"), Val(b),
				Str("reset"), Func("", func(e *Env) {
					e.Object(-1).Gets("_f").Interface().(*bytes.Buffer).Reset()
				}),
				Str("value"), Func("", func(e *Env) {
					e.A = Bytes(e.Object(-1).Gets("_f").Interface().(*bytes.Buffer).Bytes())
				}),
			).Value()
		}),
		Str("hex"), Proto(encDecProto.Object().Parent(), Str("__name"), Str("hex")),
		Str("base64"), Obj(Str("__name"), Str("base64"),
			Str("std"), Proto(encDecProto.Object(), Str("_e"), Val(getEncB64(base64.StdEncoding, '='))),
			Str("url"), Proto(encDecProto.Object(), Str("_e"), Val(getEncB64(base64.URLEncoding, '='))),
			Str("stdp"), Proto(encDecProto.Object(), Str("_e"), Val(getEncB64(base64.StdEncoding, -1))),
			Str("urlp"), Proto(encDecProto.Object(), Str("_e"), Val(getEncB64(base64.URLEncoding, -1))),
		),
		Str("base32"), Obj(Str("__name"), Str("base32"),
			Str("std"), Proto(encDecProto.Object(), Str("_e"), Val(getEncB32(base32.StdEncoding, '='))),
			Str("hex"), Proto(encDecProto.Object(), Str("_e"), Val(getEncB32(base32.HexEncoding, '='))),
			Str("stdp"), Proto(encDecProto.Object(), Str("_e"), Val(getEncB32(base32.StdEncoding, -1))),
			Str("hexp"), Proto(encDecProto.Object(), Str("_e"), Val(getEncB32(base32.HexEncoding, -1))),
		),
	))
	AddGlobalValue("str", StrLib)

	MathLib = TableMerge(MathLib, Obj(
		Str("__name"), Str("mathlib"),
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
		Str("modf"), Func("", func(e *Env) {
			a, b := math.Modf(e.Float64(0))
			e.A = Array(Float64(a), Float64(b))
		}),
	).Object())
	AddGlobalValue("math", MathLib)

	OSLib = TableMerge(OSLib, Obj(
		Str("args"), ValRec(os.Args),
		Str("environ"), Func("", func(e *Env) { e.A = ValRec(os.Environ()) }),
		Str("shell"), Func("", func(e *Env) {
			p := exec.Command("sh", "-c", e.Str(0))
			opt := e.Get(1)
			timeout := time.Duration(1 << 62) // basically forever
			if tmp := opt.ToTableGets("timeout"); tmp != Nil {
				timeout = time.Duration(tmp.MustFloat64("timeout") * float64(time.Second))
			}
			if tmp := opt.ToTableGets("env"); tmp != Nil {
				tmp.MustTable("env").Foreach(func(k, v Value) bool {
					p.Env = append(p.Env, k.String()+"="+v.String())
					return true
				})
			}
			stdout := &bytes.Buffer{}
			p.Stdout, p.Stderr = stdout, stdout
			p.Dir = opt.ToTableGets("dir").ToStr("")
			if tmp := opt.ToTableGets("stdout"); tmp != Nil {
				p.Stdout = NewWriter(tmp)
			}
			if tmp := opt.ToTableGets("stderr"); tmp != Nil {
				p.Stderr = NewWriter(tmp)
			}
			if tmp := opt.ToTableGets("stdin"); tmp != Nil {
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
			e.A = ValRec(fi)
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
			_ = err == nil && e.SetA(Val(fi))
		}),
	).Object())
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
			fmt.Fprintf(p, tmp.String(), env.Str(popi))
		case typ.Number + typ.String:
			if pop := env.Get(popi); pop.Type() == typ.String {
				fmt.Fprintf(p, tmp.String(), pop.Str())
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
