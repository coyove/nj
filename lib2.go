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
	StrLib   Value
	MathLib  Value
	TableLib Value
	OSLib    Value
	IOLib    Value
)

func init() {
	IOLib = TableMerge(IOLib, Map(
		Str("reader"), ReaderProto.Value(),
		Str("writer"), WriterProto.Value(),
		Str("seeker"), SeekerProto.Value(),
		Str("closer"), CloserProto.Value(),
		Str("readwriter"), ReadWriterProto.Value(),
		Str("readcloser"), ReadCloserProto.Value(),
		Str("writecloser"), WriteCloserProto.Value(),
		Str("readwritecloser"), ReadWriteCloserProto.Value(),
		Str("readwriteseekcloser"), ReadWriteSeekCloserProto.Value(),
	))
	AddGlobalValue("io", IOLib)

	TableLib = TableMerge(TableLib, Map(
		Str("__name"), Str("tablelib"),
		Str("__call"), Func("", func(e *Env) {
			if e.Get(1) == Nil {
				e.A = TableProto(e.Table(0))
			} else {
				e.A = e.Table(1).SetFirstParent(e.Table(0)).Value()
			}
		}),
		Str("concurrent"), Func("", func(e *Env) {
			x := NewTable(e.Table(0).Len())
			e.Table(0).Foreach(func(k, v Value) bool {
				if v.IsFunc() {
					old := v.Func()
					v = Func(v.Func().Name, func(e *Env) {
						mu := e.Recv("_mu").Interface().(*sync.Mutex)
						mu.Lock()
						defer mu.Unlock()
						e.A = MustValue(old.Call(e.Stack()...))
					}, v.Func().docString)
				}
				x.Set(k, v)
				return true
			})
			x.Sets("_mu", Val(&sync.Mutex{}))
			_ = e.Get(1).IsNil() && e.SetA(TableProto(x)) || e.SetA(e.Table(1).SetFirstParent(x).Value())
		}, "$f(src: table) -> table", "\tcreate a concurrently accessible table"),
		Str("get"), Func("", func(e *Env) {
			e.A = e.Table(0).Get(e.Get(1))
		}, "$f({t}: table, k: value) -> value"),
		Str("set"), Func("", func(e *Env) {
			e.A = e.Table(0).Set(e.Get(1), e.Get(2))
		}, "$f({t}: table, k: value, v: value) -> value", "\tset value by key, return previous value"),
		Str("rawget"), Func("", func(e *Env) {
			e.A = e.Table(0).RawGet(e.Get(1))
		}, "$f({t}: table, k: value) -> value"),
		Str("rawset"), Func("", func(e *Env) {
			e.A = e.Table(0).RawSet(e.Get(1), e.Get(2))
		}, "$f({t}: table, k: value, v: value) -> value", "\tset value by key, return previous value"),
		Str("make"), Func("", func(e *Env) {
			e.A = NewTable(e.Int(0)).Value()
		}, "$f(n: int) -> table", "\treturn a table, preallocate enough hash space for n values"),
		Str("makearray"), Func("", func(e *Env) {
			e.A = Array(make([]Value, e.Int(0))...)
		}, "$f(n: int) -> array", "\treturn a table array, preallocate space for n values"),
		Str("clear"), Func("", func(e *Env) {
			e.Table(0).Clear()
		}, "$f({t}: table)"),
		Str("cleararray"), Func("", func(e *Env) {
			e.Table(0).ClearArray()
		}, "$f({t}: table)"),
		Str("clearmap"), Func("", func(e *Env) {
			e.Table(0).ClearMap()
		}, "$f({t}: table)"),
		Str("slice"), Func("", func(e *Env) {
			e.A = Array(e.Table(0).items[e.Int(1):e.Int(2)]...)
		}, "$f({t}: array, start: int, end: int) -> array", "\tslice array, from start to end"),
		Str("copy"), Func("", func(e *Env) {
			if e.Size() == 1 {
				e.A = e.Table(0).Copy().Value()
			} else {
				a, start, end := e.Table(0).items, e.Int(1), e.Int(2)
				if start >= 0 && start < len(a) && end >= 0 && end <= len(a) && start <= end {
					e.A = Array(append([]Value{}, a[start:end]...)...)
				} else {
					m, from := NewTable(0), e.Table(0)
					for i := start; i < end; i++ {
						m.Set(Int(i-start), from.Get(Int(i)))
					}
					e.A = m.Value()
				}
			}
		}, "$f({t}: table) -> table", "\tcopy the entire table",
			"$f({t}: array, start: int, end: int) -> array", "\tcopy array, from start to end"),
		Str("parent"), Func("", func(e *Env) {
			e.A = e.Table(0).Parent().Value()
		}, "$f({t}: table) -> table", "\treturn table's parent if any"),
		Str("setparent"), Func("", func(e *Env) {
			e.Table(0).SetParent(e.Table(1))
		}, "$f({t}: table, p: table)", "\tset table's parent"),
		Str("setfirstparent"), Func("", func(e *Env) {
			e.Table(0).SetFirstParent(e.Table(1))
		}, "$f({t}: table, p: table)", "\tinsert `p` as `t`'s first parent"),
		Str("arraylen"), Func("", func(e *Env) {
			e.A = Int(e.Table(0).ArrayLen())
		}, "$f({t}: array) -> int", "\treturn the length of array"),
		Str("maplen"), Func("", func(e *Env) {
			e.A = Int(e.Table(0).MapLen())
		}, "$f({t}: table) -> int", "\treturn the size of table map"),
		Str("arraysize"), Func("", func(e *Env) {
			e.A = Int(len(e.Table(0).items))
		}, "$f({t}: array) -> int", "\treturn the true size of array (including nils)"),
		Str("mapsize"), Func("", func(e *Env) {
			e.A = Int(len(e.Table(0).hashItems))
		}, "$f({t}: table) -> int", "\treturn the true size of table map (including empty nil entries)"),
		Str("keys"), Func("", func(e *Env) {
			a := make([]Value, 0)
			e.Table(0).Foreach(func(k, v Value) bool { a = append(a, k); return true })
			e.A = Array(a...)
		}, "$f({t}: table) -> array"),
		Str("values"), Func("", func(e *Env) {
			a := make([]Value, 0)
			e.Table(0).Foreach(func(k, v Value) bool { a = append(a, v); return true })
			e.A = Array(a...)
		}, "$f({t}: table) -> array"),
		Str("items"), Func("", func(e *Env) {
			a := make([]Value, 0)
			e.Table(0).Foreach(func(k, v Value) bool { a = append(a, Array(k, v)); return true })
			e.A = Array(a...)
		}, "$f({t}: table) -> array"),
		Str("foreach"), Func("", func(e *Env) {
			f := e.Func(1)
			e.Table(0).Foreach(func(k, v Value) bool { return MustValue(f.Call(k, v)) == Nil })
		}, "$f({t}: table, f: function)"),
		Str("contains"), Func("", func(e *Env) {
			found, b := false, e.Get(1)
			e.Table(0).Foreach(func(k, v Value) bool { found = v.Equal(b); return !found })
			e.A = Bool(found)
		}, "$f({t}: table, v: value) -> bool", "\ttest if table contains value"),
		Str("append"), Func("", func(e *Env) {
			ma := e.Array(0)
			ma.items = append(ma.items, e.Stack()[1:]...)
			e.A = ma.Value()
		}, "$f({a}: array, args...: value)", "\tappend values to array"),
		Str("filter"), Func("", func(e *Env) {
			ma, ff := e.Array(0), e.Func(1)
			a2 := make([]Value, 0, ma.ArrayLen())
			ma.Foreach(func(k, v Value) bool {
				if MustValue(ff.Call(v)).IsTrue() {
					a2 = append(a2, v)
				}
				return true
			})
			e.A = Array(a2...)
		}, "$f({a}: array, f: function)", "\tfilter out all values where f(value) is false"),
		Str("removeat"), Func("", func(e *Env) {
			ma, idx := e.Array(0), e.Int(1)
			if idx < 0 || idx >= len(ma.items) {
				e.A = Nil
				return
			}
			old := ma.items[idx]
			ma.items = append(ma.items[:idx], ma.items[idx+1:]...)
			if old != Nil {
				ma.count--
			}
			e.A = old
		}, "$f({a}: array, index: int)", "\tremove value at `index`"),
		Str("concat"), Func("", func(e *Env) {
			ma, mb := e.Table(0), e.Table(1)
			for _, b := range mb.ArrayPart() {
				ma.Set(Int(len(ma.items)), b)
			}
			e.A = ma.Value()
		}, "$f({array1}: array, array2: array)", "\tput elements from array2 to array1's end"),
		Str("merge"), Func2("", func(a, b Value) Value {
			return a.MustTable("").Merge(b.MustTable("")).Value()
		}, "$f({table1}: table, table2: table)", "\tmerge elements from table2 to table1"),
		Str("tostring"), Func1("", func(a Value) Value {
			p := &bytes.Buffer{}
			a.MustTable("").rawPrint(p, 0, true, true)
			return Bytes(p.Bytes())
		}, "$f({t}: table) -> string", "\tprint raw elements in table, ignore __str"),
		Str("name"), Func1("", func(a Value) Value {
			return Str(a.MustTable("").Name())
		}, "$f({t}: table) -> string", "\tprint table's name"),
		Str("pure"), Func1("", func(m Value) Value {
			m2 := *m.MustTable("")
			m2.parent = nil
			return m2.Value()
		}, "$f({t}: table) -> table", "\treturn a new table who shares all the data from t, but with no parent"),
		Str("unwrap"), Func("", func(e *Env) {
			v := e.Get(0)
			_ = v.Type() == typ.Native && e.SetA(ValRec(v.Interface())) || e.SetA(v)
		}, "unwrap(v: value) -> table", "\tunwrap Go's array, slice or map into table"),
	))
	AddGlobalValue("table", TableLib)

	encDecProto := TableProto(Map(
		Str("__name"), Str("encdec"),
		Str("encoder"), Func("", func(e *Env) {
			enc := Nil
			buf := &bytes.Buffer{}
			switch encoding := e.Recv("_e").Interface().(type) {
			default:
				enc = Val(hex.NewEncoder(buf))
			case *base32.Encoding:
				enc = Val(base32.NewEncoder(encoding, buf))
			case *base64.Encoding:
				enc = Val(base64.NewEncoder(encoding, buf))
			}
			e.A = TableProto(WriteCloserProto,
				Str("_f"), Val(enc),
				Str("_b"), Val(buf),
				Str("value"), Func("", func(e *Env) {
					e.A = Bytes(e.Recv("_b").Interface().(*bytes.Buffer).Bytes())
				}),
			)
		}),
		Str("decoder"), Func("", func(e *Env) {
			src := NewReader(e.Get(1))
			dec := Nil
			switch encoding := e.Recv("_e").Interface().(type) {
			case *base64.Encoding:
				dec = Val(base64.NewDecoder(encoding, src))
			case *base32.Encoding:
				dec = Val(base32.NewDecoder(encoding, src))
			default:
				dec = Val(hex.NewDecoder(src))
			}
			e.A = TableProto(ReaderProto, Str("_f"), Val(dec))
		}),
	).Table(),
		Str("__name"), Str("encdecfast"),
		Str("encode"), Func2("", func(m, t Value) Value {
			x := struct {
				a string
				c int
			}{a: t.MustStr(""), c: t.StrLen()}
			return Str(m.Recv("_e").Interface().(interface {
				EncodeToString([]byte) string
			}).EncodeToString(*(*[]byte)(unsafe.Pointer(&x))))
		}),
		Str("decode"), Func2("", func(m, t Value) Value {
			v, err := m.Recv("_e").Interface().(interface {
				DecodeString(string) ([]byte, error)
			}).DecodeString(t.MustStr(""))
			internal.PanicErr(err)
			return Bytes(v)
		}),
	)

	StrLib = TableMerge(StrLib, Map(
		Str("__name"), Str("strlib"),
		Str("__call"), Func("", func(e *Env) {
			i, ok := e.Interface(1).([]byte)
			_ = ok && e.SetA(Bytes(i)) || e.SetA(Str(e.Get(1).String()))
		}),
		Str("size"), Func("", func(e *Env) {
			e.A = Int(e.StrLen(0))
		}, "$f({text}: string) -> int"),
		Str("len"), Func("", func(e *Env) {
			e.A = Int(e.StrLen(0))
		}, "$f({text}: string) -> int"),
		Str("count"), Func("", func(e *Env) {
			e.A = Int(utf8.RuneCountInString(e.Str(0)))
		}, "$f({text}: string) -> int", "\treturn count of runes in text"),
		Str("from"), Func("", func(e *Env) {
			e.A = Str(fmt.Sprint(e.Interface(0)))
		}, "$f(v: value) -> string", "\tconvert value to string"),
		Str("iequals"), Func("", func(e *Env) {
			e.A = Bool(strings.EqualFold(e.Str(0), e.Str(1)))
		}, "$f({text1}: string, text2: string) -> bool"),
		Str("contains"), Func("", func(e *Env) {
			e.A = Bool(strings.Contains(e.Str(0), e.Str(1)))
		}, "$f({text}: string, substr: string) bool"),
		Str("containsany"), Func("", func(e *Env) {
			e.A = Bool(strings.ContainsAny(e.Str(0), e.Str(1)))
		}, "$f({text}: string, chars: string) bool"),
		Str("split"), Func3("", func(src, delim, n Value) Value {
			s := src.MustStr("text")
			d := delim.MustStr("delimeter")
			var r []Value
			if n := n.ToInt64(0); n == 0 {
				for _, p := range strings.Split(s, d) {
					r = append(r, Str(p))
				}
			} else {
				for _, p := range strings.SplitN(s, d, int(n)) {
					r = append(r, Str(p))
				}
			}
			return Array(r...)
		}, "split({text}: string, delim: string, n?: int) array"),
		Str("join"), Func2("", func(delim, array Value) Value {
			d := delim.MustStr("delimeter")
			buf := &bytes.Buffer{}
			for _, a := range array.MustTable("").ArrayPart() {
				buf.WriteString(a.String())
				buf.WriteString(d)
			}
			if buf.Len() > 0 {
				buf.Truncate(buf.Len() - len(d))
			}
			return Bytes(buf.Bytes())
		}, "$f({text}: string, a: array) -> string"),
		Str("replace"), Func("", func(e *Env) {
			e.A = Str(strings.Replace(e.Str(0), e.Str(1), e.Str(2), e.Get(3).ToInt(-1)))
		}, "$f({text}: string, old: string, new: string) -> string"),
		Str("match"), Func2("", func(pattern, str Value) Value {
			m, err := filepath.Match(pattern.MustStr("pattern"), str.MustStr("text"))
			internal.PanicErr(err)
			return Bool(m)
		}, "$f({pattern}: string, text: string) -> bool"),
		Str("find"), Func("", func(e *Env) {
			e.A = Int(strings.Index(e.Str(0), e.Str(1)))
		}, "$f({text}: string, sub: string) -> int", "\tfind index of first appearence of `sub` in `text`"),
		Str("findany"), Func("", func(e *Env) {
			e.A = Int(strings.IndexAny(e.Str(0), e.Str(1)))
		}, "$f({text}: string, charset: string) -> int", "\tfind index of first appearence of any char from `charset` in `text`"),
		Str("rfind"), Func("", func(e *Env) {
			e.A = Int(strings.LastIndex(e.Str(0), e.Str(1)))
		}, "$f({text}: string, sub: string) -> int", "\tsame as find(), but from right to left"),
		Str("rfindany"), Func("", func(e *Env) {
			e.A = Int(strings.LastIndexAny(e.Str(0), e.Str(1)))
		}, "$f({text}: string, charset: string) -> int", "\tsame as findany(), but from right to left"),
		Str("sub"), Func("", func(e *Env) {
			s := e.Str(0)
			st, en := e.Int(1), e.Get(2).ToInt(len(s))
			for ; st < 0 && len(s) > 0; st += len(s) {
			}
			for ; en < 0 && len(s) > 0; en += len(s) {
			}
			e.A = Str(s[st:en])
		}, "$f({text}: string, start: int, end?: int) -> string"),
		Str("trim"), Func("", func(e *Env) {
			_ = e.Get(1).IsNil() && e.SetA(Str(strings.TrimSpace(e.Str(0)))) || e.SetA(Str(strings.Trim(e.Str(0), e.Str(1))))
		}, "$f({text}: string, cutset?: string) -> string", "\ttrim spaces (or any chars in `cutset`) at left and right side of `text`"),
		Str("lremove"), Func("", func(e *Env) {
			e.A = Str(strings.TrimPrefix(e.Str(0), e.Str(1)))
		}, "$f({text}: string, prefix: string) -> string", "\tremove `prefix` of `text`"),
		Str("rremove"), Func("", func(e *Env) {
			e.A = Str(strings.TrimSuffix(e.Str(0), e.Str(1)))
		}, "$f({text}: string, suffix: string) -> string", "\tremove `suffix` of `text`"),
		Str("ltrim"), Func("", func(e *Env) {
			e.A = Str(strings.TrimLeft(e.Str(0), e.Str(1)))
		}, "$f({text}: string, cutset: string) -> string", "\tremove chars both ocurred in `cutset` and left-side of `text`"),
		Str("rtrim"), Func("", func(e *Env) {
			e.A = Str(strings.TrimRight(e.Str(0), e.Str(1)))
		}, "$f({text}: string, cutset: string) -> string", "\tremove chars both ocurred in `cutset` and right-side of `text`"),
		Str("ord"), Func("", func(e *Env) {
			r, sz := utf8.DecodeRuneInString(e.Str(0))
			e.A = Array(Int64(int64(r)), Int(sz))
		}, "$f({text}: string) -> array", "\tdecode first UTF-8 char in `text`, return { unicode, width_in_bytes }"),
		Str("startswith"), Func("", func(e *Env) {
			e.A = Bool(strings.HasPrefix(e.Str(0), e.Str(1)))
		}, "$f({text}: string, prefix: string) -> bool"),
		Str("endswith"), Func("", func(e *Env) {
			e.A = Bool(strings.HasSuffix(e.Str(0), e.Str(1)))
		}, "$f({text}: string, suffix: string) -> bool"),
		Str("upper"), Func("", func(e *Env) {
			e.A = Str(strings.ToUpper(e.Str(0)))
		}, "$f({text}: string) -> string"),
		Str("lower"), Func("", func(e *Env) {
			e.A = Str(strings.ToLower(e.Str(0)))
		}, "$f({text}: string) -> string"),
		Str("bytes"), Func("", func(e *Env) {
			_ = e.Get(0).IsInt64() && e.SetA(Val(make([]byte, e.Int(0)))) || e.SetA(Val([]byte(e.Str(0))))
		}, "$f({text}: string) -> go.[]byte", "\tcreate a byte array from `text`",
			"$f(n: int) -> go.[]byte", "\tcreate an n-byte long array"),
		Str("chars"), Func("", func(e *Env) {
			var r []Value
			for s := e.Str(0); len(s) > 0; {
				_, sz := utf8.DecodeRuneInString(s)
				if sz == 0 {
					break
				}
				r = append(r, Str(s[:sz]))
				s = s[sz:]
			}
			e.A = Array(r...)
		}, "$f({text}: string) -> array", "\tbreak `text` into chars, e.g.: chars('a中c') => { 'a', '中', 'c' }"),
		Str("format"), Func("", func(e *Env) {
			buf := &bytes.Buffer{}
			sprintf(e, buf)
			e.A = Bytes(buf.Bytes())
		}, "$f({pattern}: string, args...: value) -> string"),
		Str("buffer"), Func("", func(e *Env) {
			b := &bytes.Buffer{}
			if v := e.Get(0); v != Nil {
				b.WriteString(v.String())
			}
			e.A = TableProto(ReadWriterProto,
				Str("_f"), Val(b),
				Str("reset"), Func("", func(e *Env) {
					e.Recv("_f").Interface().(*bytes.Buffer).Reset()
				}),
				Str("value"), Func("", func(e *Env) {
					e.A = Bytes(e.Recv("_f").Interface().(*bytes.Buffer).Bytes())
				}),
			)
		}),
		Str("hex"), TableProto(encDecProto.Table().Parent(), Str("__name"), Str("hex")),
		Str("base64"), Map(Str("__name"), Str("base64"),
			Str("std"), TableProto(encDecProto.Table(), Str("_e"), Val(getEncB64(base64.StdEncoding, '='))),
			Str("url"), TableProto(encDecProto.Table(), Str("_e"), Val(getEncB64(base64.URLEncoding, '='))),
			Str("stdp"), TableProto(encDecProto.Table(), Str("_e"), Val(getEncB64(base64.StdEncoding, -1))),
			Str("urlp"), TableProto(encDecProto.Table(), Str("_e"), Val(getEncB64(base64.URLEncoding, -1))),
		),
		Str("base32"), Map(Str("__name"), Str("base32"),
			Str("std"), TableProto(encDecProto.Table(), Str("_e"), Val(getEncB32(base32.StdEncoding, '='))),
			Str("hex"), TableProto(encDecProto.Table(), Str("_e"), Val(getEncB32(base32.HexEncoding, '='))),
			Str("stdp"), TableProto(encDecProto.Table(), Str("_e"), Val(getEncB32(base32.StdEncoding, -1))),
			Str("hexp"), TableProto(encDecProto.Table(), Str("_e"), Val(getEncB32(base32.HexEncoding, -1))),
		),
	))
	AddGlobalValue("str", StrLib)

	MathLib = TableMerge(MathLib, Map(
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
	))
	AddGlobalValue("math", MathLib)

	OSLib = TableMerge(OSLib, Map(
		Str("__name"), Str("oslib"),
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
	))
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

func sprintf(env *Env, p io.Writer) {
	f := env.Str(0)
	tmp := bytes.Buffer{}
	popi := 0
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
