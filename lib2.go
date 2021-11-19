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
		Str("__call"), Func2("", func(t, m Value) Value {
			if m == Nil {
				return TableProto(t.MustTable(""))
			}
			m.MustTable("").SetFirstParent(t.MustTable(""))
			return m
		}),
		Str("concurrent"), Func2("", func(t, m Value) Value {
			x := NewTable(t.MustTable("").Len())
			t.MustTable("").Foreach(func(k, v Value) bool {
				if v.IsFunc() {
					old := v.Func()
					v = Func(v.Func().Name, func(e *Env) {
						mu := e.B(0).Recv("_mu").Interface().(*sync.Mutex)
						mu.Lock()
						defer mu.Unlock()
						e.A = MustValue(old.Call(e.Stack()...))
					}, v.Func().docString)
				}
				x.Set(k, v)
				return true
			})
			x.Sets("_mu", Val(&sync.Mutex{}))
			if m == Nil {
				return TableProto(x)
			}
			m.MustTable("").SetFirstParent(x)
			return m
		}, "$f(src: table) -> table", "\tcreate a concurrently accessible table"),
		Str("get"), Func2("", func(m, k Value) Value {
			return m.MustTable("").Get(k)
		}, "$f({t}: table, k: value) -> value"),
		Str("set"), Func3("", func(m, k, v Value) Value {
			return m.MustTable("").Set(k, v)
		}, "$f({t}: table, k: value, v: value) -> value", "\tset value by key, return previous value"),
		Str("rawget"), Func2("", func(m, k Value) Value {
			return m.MustTable("").RawGet(k)
		}, "$f({t}: table, k: value) -> value"),
		Str("rawset"), Func3("", func(m, k, v Value) Value {
			return m.MustTable("").RawSet(k, v)
		}, "$f({t}: table, k: value, v: value) -> value", "\tset value by key, return previous value"),
		Str("make"), Func1("", func(n Value) Value {
			return NewTable(n.MustInt("")).Value()
		}, "$f(n: int) -> table", "\treturn a table, preallocate enough hash space for n values"),
		Str("makearray"), Func1("", func(n Value) Value {
			return Array(make([]Value, n.MustInt64(""))...)
		}, "$f(n: int) -> array", "\treturn a table array, preallocate space for n values"),
		Str("clear"), Func("", func(env *Env) {
			env.B(0).MustTable("").Clear()
		}, "$f({t}: table)"),
		Str("cleararray"), Func("", func(env *Env) {
			env.B(0).MustTable("").ClearArray()
		}, "$f({t}: table)"),
		Str("clearmap"), Func("", func(env *Env) {
			env.B(0).MustTable("").ClearMap()
		}, "$f({t}: table)"),
		Str("slice"), Func3("", func(t, s, e Value) Value {
			return Array(t.MustTable("").items[s.MustInt(""):e.MustInt("")]...)
		}, "$f({t}: array, start: int, end: int) -> array", "\tslice array, from start to end"),
		Str("copy"), Func3("", func(t, s, e Value) Value {
			if s == Nil && e == Nil {
				return t.MustTable("").Copy().Value()
			}
			start, end := s.MustInt(""), e.MustInt("")
			a := t.MustTable("").items
			if start >= 0 && start < len(a) && end >= 0 && end <= len(a) && start <= end {
				return Array(append([]Value{}, a[start:end]...)...)
			}
			m, from := NewTable(0), t.Table()
			for i := start; i < end; i++ {
				m.Set(Int(i-start), from.Get(Int(i)))
			}
			return m.Value()
		},
			"$f({t}: table) -> table", "\tcopy the entire table",
			"$f({t}: array, start: int, end: int) -> array", "\tcopy array, from start to end",
		),
		Str("parent"), Func1("", func(v Value) Value {
			return v.MustTable("").Parent().Value()
		}, "$f({t}: table) -> table", "\treturn table's parent if any"),
		Str("setparent"), Func("", func(e *Env) {
			e.B(0).MustTable("").SetParent(e.B(1).MustTable(""))
		}, "$f({t}: table, p: table)", "\tset table's parent"),
		Str("setfirstparent"), Func("", func(e *Env) {
			e.B(0).MustTable("").SetFirstParent(e.B(1).MustTable(""))
		}, "$f({t}: table, p: table)", "\tinsert `p` as `t`'s first parent"),
		Str("arraylen"), Func1("", func(v Value) Value {
			return Int(v.MustTable("").ArrayLen())
		}, "$f({t}: array) -> int", "\treturn the length of array"),
		Str("maplen"), Func1("", func(v Value) Value {
			return Int(v.MustTable("").MapLen())
		}, "$f({t}: table) -> int", "\treturn the size of table map"),
		Str("arraysize"), Func1("", func(v Value) Value {
			return Int(len(v.MustTable("").items))
		}, "$f({t}: array) -> int", "\treturn the true size of array (including nils)"),
		Str("mapsize"), Func1("", func(v Value) Value {
			return Int(len(v.MustTable("").hashItems))
		}, "$f({t}: table) -> int", "\treturn the true size of table map (including empty nil entries)"),
		Str("keys"), Func1("", func(m Value) Value {
			a := make([]Value, 0)
			m.MustTable("").Foreach(func(k, v Value) bool { a = append(a, k); return true })
			return Array(a...)
		}, "$f({t}: table) -> array"),
		Str("values"), Func1("", func(m Value) Value {
			a := make([]Value, 0)
			m.MustTable("").Foreach(func(k, v Value) bool { a = append(a, v); return true })
			return Array(a...)
		}, "$f({t}: table) -> array"),
		Str("items"), Func1("", func(m Value) Value {
			a := make([]Value, 0)
			m.MustTable("").Foreach(func(k, v Value) bool { a = append(a, Array(k, v)); return true })
			return Array(a...)
		}, "$f({t}: table) -> array"),
		Str("foreach"), Func("", func(e *Env) {
			f := e.B(1).MustFunc("")
			e.B(0).MustTable("").Foreach(func(k, v Value) bool { return MustValue(f.Call(k, v)) == Nil })
		}, "$f({t}: table, f: function)"),
		Str("contains"), Func2("", func(a, b Value) Value {
			found := false
			a.MustTable("").Foreach(func(k, v Value) bool { found = v.Equal(b); return !found })
			return Bool(found)
		}, "$f({t}: table, v: value) -> bool", "\ttest if table contains value"),
		Str("append"), Func("", func(env *Env) {
			ma := env.B(0).MustTable("")
			ma.items = append(ma.items, env.Stack()[1:]...)
			env.A = ma.Value()
		}, "$f({a}: array, args...: value)", "\tappend values to array"),
		Str("filter"), Func2("", func(a, b Value) Value {
			ma := a.MustTable("")
			a2 := make([]Value, 0, ma.ArrayLen())
			ma.Foreach(func(k, v Value) bool {
				if MustValue(b.MustFunc("").Call(v)).IsTrue() {
					a2 = append(a2, v)
				}
				return true
			})
			return Array(a2...)
		}, "$f({a}: array, f: function)", "\tfilter out all values where f(value) is false"),
		Str("removeat"), Func2("", func(a, b Value) Value {
			ma, idx := a.MustTable(""), b.MustInt("")
			if idx < 0 || idx >= len(ma.items) {
				return Nil
			}
			old := ma.items[idx]
			ma.items = append(ma.items[:idx], ma.items[idx+1:]...)
			if old != Nil {
				ma.count--
			}
			return old
		}, "$f({a}: array, index: int)", "\tremove value at `index`"),
		Str("concat"), Func2("", func(a, b Value) Value {
			ma, mb := a.MustTable(""), b.MustTable("")
			for _, b := range mb.ArrayPart() {
				ma.Set(Int(len(ma.items)), b)
			}
			return ma.Value()
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
		Str("unwrap"), Func1("unwrap", func(m Value) Value {
			if m.Type() == typ.Native {
				return ValRec(m.Interface())
			}
			return m
		}, "unwrap(v: value) -> table", "\tunwrap Go's array, slice or map into table"),
	))
	AddGlobalValue("table", TableLib)

	encDecProto := TableProto(Map(
		Str("__name"), Str("encdec"),
		Str("encoder"), Func1("", func(m Value) Value {
			enc := Nil
			buf := &bytes.Buffer{}
			switch encoding := m.Recv("_e").Interface().(type) {
			default:
				enc = Val(hex.NewEncoder(buf))
			case *base32.Encoding:
				enc = Val(base32.NewEncoder(encoding, buf))
			case *base64.Encoding:
				enc = Val(base64.NewEncoder(encoding, buf))
			}
			return TableProto(WriteCloserProto,
				Str("_f"), Val(enc),
				Str("_b"), Val(buf),
				Str("value"), Func1("", func(p Value) Value {
					return Bytes(p.Recv("_b").Interface().(*bytes.Buffer).Bytes())
				}),
			)
		}),
		Str("decoder"), Func("", func(e *Env) {
			src := NewReader(e.B(1))
			dec := Nil
			switch encoding := e.B(0).Recv("_e").Interface().(type) {
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
		Str("__call"), Func2("", func(strObj, src Value) Value {
			if i, ok := src.Interface().([]byte); ok {
				return Bytes(i)
			}
			return Str(src.String())
		}),
		Str("size"), Func1("", func(src Value) Value {
			return Int(src.MustStrLen(""))
		}, "$f({text}: string) -> int"),
		Str("len"), Func1("", func(src Value) Value {
			return Int(src.MustStrLen(""))
		}, "$f({text}: string) -> int"),
		Str("count"), Func1("", func(src Value) Value {
			return Int(utf8.RuneCountInString(src.MustStr("")))
		}, "$f({text}: string) -> int", "\treturn count of runes in text"),
		Str("from"), Func1("", func(src Value) Value {
			return Str(fmt.Sprint(src.Interface()))
		}, "$f(v: value) -> string", "\tconvert value to string"),
		Str("iequals"), Func2("", func(src, a Value) Value {
			return Bool(strings.EqualFold(src.MustStr(""), a.MustStr("")))
		}, "$f({text1}: string, text2: string) -> bool"),
		Str("contains"), Func2("", func(src, a Value) Value {
			return Bool(strings.Contains(src.MustStr(""), a.MustStr("")))
		}, "$f({text}: string, substr: string) bool"),
		Str("containsany"), Func2("", func(src, a Value) Value {
			return Bool(strings.ContainsAny(src.MustStr(""), a.MustStr("")))
		}, "$f({text}: string, chars: string) bool"),
		Str("split"), Func3("", func(src, delim, n Value) Value {
			s := src.MustStr("text")
			d := delim.MustStr("delimeter")
			r := []Value{}
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
			src := e.B(0).MustStr("text")
			from, to := e.B(1).MustStr("old"), e.B(2).MustStr("new")
			e.A = Str(strings.Replace(src, from, to, e.B(3).ToInt(-1)))
		}, "$f({text}: string, old: string, new: string) -> string"),
		Str("match"), Func2("", func(pattern, str Value) Value {
			m, err := filepath.Match(pattern.MustStr("pattern"), str.MustStr("text"))
			internal.PanicErr(err)
			return Bool(m)
		}, "$f({pattern}: string, text: string) -> bool"),
		Str("find"), Func2("", func(src, substr Value) Value {
			return Int(strings.Index(src.MustStr(""), substr.MustStr("")))
		}, "$f({text}: string, sub: string) -> int", "\tfind index of first appearence of `sub` in `text`"),
		Str("findany"), Func2("", func(src, substr Value) Value {
			return Int(strings.IndexAny(src.MustStr(""), substr.MustStr("")))
		}, "$f({text}: string, charset: string) -> int", "\tfind index of first appearence of any char from `charset` in `text`"),
		Str("rfind"), Func2("", func(src, substr Value) Value {
			return Int(strings.LastIndex(src.MustStr(""), substr.MustStr("")))
		}, "$f({text}: string, sub: string) -> int", "\tsame as find(), but from right to left"),
		Str("rfindany"), Func2("", func(src, substr Value) Value {
			return Int(strings.LastIndexAny(src.MustStr(""), substr.MustStr("")))
		}, "$f({text}: string, charset: string) -> int", "\tsame as findany(), but from right to left"),
		Str("sub"), Func3("", func(src, start, end Value) Value {
			s := src.MustStr("")
			st, en := start.ToInt(0), end.ToInt(len(s))
			for ; st < 0 && len(s) > 0; st += len(s) {
			}
			for ; en < 0 && len(s) > 0; en += len(s) {
			}
			return Str(s[st:en])
		}, "$f({text}: string, start: int, end: int) -> string"),
		Str("trim"), Func2("", func(src, cutset Value) Value {
			if cutset == Nil {
				return Str(strings.TrimSpace(src.MustStr("")))
			}
			return Str(strings.Trim(src.MustStr(""), cutset.MustStr("")))
		}, "$f({text}: string, cutset?: string) -> string", "\ttrim spaces (or any chars in `cutset`) at left and right side of `text`"),
		Str("lremove"), Func2("", func(src, cutset Value) Value {
			return Str(strings.TrimPrefix(src.MustStr(""), cutset.MustStr("")))
		}, "$f({text}: string, prefix: string) -> string", "\tremove `prefix` in `text` if any"),
		Str("rremove"), Func2("", func(src, cutset Value) Value {
			return Str(strings.TrimSuffix(src.MustStr(""), cutset.MustStr("")))
		}, "$f({text}: string, suffix: string) -> string", "\tremove `suffix` in `text` if any"),
		Str("ltrim"), Func2("", func(src, cutset Value) Value {
			return Str(strings.TrimLeft(src.MustStr(""), cutset.ToStr(" ")))
		}, "$f({text}: string, cutset: string) -> string", "\tremove chars both ocurred in `cutset` and left-side of `text`"),
		Str("rtrim"), Func2("", func(src, cutset Value) Value {
			return Str(strings.TrimRight(src.MustStr(""), cutset.ToStr(" ")))
		}, "$f({text}: string, cutset: string) -> string", "\tremove chars both ocurred in `cutset` and right-side of `text`"),
		Str("ord"), Func("", func(env *Env) {
			r, sz := utf8.DecodeRuneInString(env.B(0).MustStr(""))
			env.A = Array(Int64(int64(r)), Int(sz))
		}, "$f({text}: string) -> array", "\tdecode first UTF-8 char in `text`, return { unicode, width_in_bytes }"),
		Str("startswith"), Func2("", func(t, p Value) Value {
			return Bool(strings.HasPrefix(t.MustStr(""), p.MustStr("")))
		}, "$f({text}: string, prefix: string) -> bool"),
		Str("endswith"), Func2("", func(t, s Value) Value {
			return Bool(strings.HasSuffix(t.MustStr(""), s.MustStr("")))
		}, "$f({text}: string, suffix: string) -> bool"),
		Str("upper"), Func1("", func(t Value) Value {
			return Str(strings.ToUpper(t.MustStr("")))
		}, "$f({text}: string) -> string"),
		Str("lower"), Func1("", func(t Value) Value {
			return Str(strings.ToLower(t.MustStr("")))
		}, "$f({text}: string) -> string"),
		Str("bytes"), Func1("", func(s Value) Value {
			if s.Type() == typ.Number {
				return Val(make([]byte, s.Int64()))
			}
			return Val([]byte(s.MustStr("")))
		}, "$f({text}: string) -> go.[]byte", "\tcreate a byte array from `text`",
			"$f(n: int) -> go.[]byte", "\tcreate an n-byte long array"),
		Str("chars"), Func1("", func(s Value) Value {
			var r []Value
			for s := s.MustStr(""); len(s) > 0; {
				_, sz := utf8.DecodeRuneInString(s)
				if sz == 0 {
					break
				}
				r = append(r, Str(s[:sz]))
				s = s[sz:]
			}
			return Array(r...)
		}, "$f({text}: string) -> array", "\tbreak `text` into chars, e.g.: chars('a中c') => { 'a', '中', 'c' }"),
		Str("format"), Func("", func(e *Env) {
			buf := &bytes.Buffer{}
			sprintf(e, buf)
			e.A = Bytes(buf.Bytes())
		}, "$f({pattern}: string, args...: value) -> string"),
		Str("buffer"), Func1("", func(v Value) Value {
			b := &bytes.Buffer{}
			if v != Nil {
				b.WriteString(v.String())
			}
			return TableProto(ReadWriterProto,
				Str("_f"), Val(b),
				Str("reset"), Func("", func(e *Env) {
					e.B(0).Recv("_f").Interface().(*bytes.Buffer).Reset()
				}),
				Str("value"), Func1("", func(a Value) Value {
					return Bytes(a.Recv("_f").Interface().(*bytes.Buffer).Bytes())
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
		Str("ldexp"), Func("", func(e *Env) { e.A = Float64(math.Ldexp(e.Float64(0), e.B(1).ToInt(0))) }),
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
			if fi, err := os.Stat(e.Str(0)); err == nil {
				e.A = Val(fi)
			}
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
			fmt.Fprint(p, env.Get(popi).Interface())
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
