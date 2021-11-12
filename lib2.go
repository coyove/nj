package script

import (
	"bytes"
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

	"github.com/coyove/script/typ"
)

var StringMethods, MathLib, TableLib, OSLib, IOLib Value

func init() {
	IOLib = TableMerge(IOLib, Nil,
		Str("reader"), ReaderProto.Value(),
		Str("writer"), WriterProto.Value(),
		Str("seeker"), SeekerProto.Value(),
		Str("closer"), CloserProto.Value(),
		Str("readwriter"), ReadWriterProto.Value(),
		Str("readcloser"), ReadCloserProto.Value(),
		Str("writecloser"), WriteCloserProto.Value(),
		Str("readwritecloser"), ReadWriteCloserProto.Value(),
		Str("readwriteseekcloser"), ReadWriteSeekCloserProto.Value(),
	)
	AddGlobalValue("io", IOLib)

	TableLib = TableMerge(TableLib, Nil,
		Str("__name"), Str("tablelib"),
		Str("__call"), Func2("__call", func(t, m Value) Value {
			if m == Nil {
				return TableProto(t.MustTable(""))
			}
			m.MustTable("").SetFirstParent(t.MustTable(""))
			return m
		}),
		Str("concurrent"), Func2("concurrent", func(t, m Value) Value {
			x := NewTable(t.MustTable("").Len())
			t.MustTable("").Foreach(func(k, v Value) bool {
				if v.Type() == typ.Func {
					old := v.Func()
					v = Function(v.Func().Name, func(env *Env) {
						mu := env.Get(0).MustTable("").GetString("_mu").Interface().(*sync.Mutex)
						mu.Lock()
						defer mu.Unlock()
						a, err := old.Call(env.Stack()...)
						panicErr(err)
						env.A = a
					}, v.Func().DocString)
				}
				x.Set(k, v)
				return true
			})
			x.SetString("_mu", Val(&sync.Mutex{}))
			if m == Nil {
				return TableProto(x)
			}
			m.MustTable("").SetFirstParent(x)
			return m
		}, "$f(src: table) table", "\tcreate a concurrently accessible table"),
		Str("get"), Func2("get", func(m, k Value) Value {
			return m.MustTable("").Get(k)
		}, "$f({t}: table, k: value) value"),
		Str("set"), Func3("set", func(m, k, v Value) Value {
			m.MustTable("").Set(k, v)
			return m
		}, "$f({t}: table, k: value, v: value) table"),
		Str("make"), Func1("make", func(n Value) Value {
			return NewTable(int(n.MustInt(""))).Value()
		}, "$f(n: int) table", "\treturn a table, preallocate enough hash space for n values"),
		Str("makearray"), Func1("makearray", func(n Value) Value {
			return Array(make([]Value, n.MustInt(""))...)
		}, "$f(n: int) array", "\treturn a table array, preallocate space for n values"),
		Str("clear"), Function("clear", func(env *Env) { env.Get(0).MustTable("").Clear() }, "$f({t}: table)"),
		Str("cleararray"), Function("cleararray", func(env *Env) { env.Get(0).MustTable("").ClearArray() }, "$f({t}: table)"),
		Str("clearmap"), Function("clearmap", func(env *Env) { env.Get(0).MustTable("").ClearMap() }, "$f({t}: table)"),
		Str("slice"), Func3("slice", func(t, s, e Value) Value {
			start, end := int(s.MustInt("")), int(e.MustInt(""))
			return Array(t.MustTable("").items[start:end]...)
		}, "$f({t}: array, start: int, end: int) array", "\tslice array, from start to end"),
		Str("copy"), Func3("copy", func(t, s, e Value) Value {
			if s == Nil && e == Nil {
				return t.MustTable("").Copy().Value()
			}
			start, end := int(s.MustInt("")), int(e.MustInt(""))
			a := t.MustTable("").items
			if start >= 0 && start < len(a) && end >= 0 && end <= len(a) && start <= end {
				return Array(append([]Value{}, a[start:end]...)...)
			}
			m, from := NewTable(0), t.Table()
			for i := start; i < end; i++ {
				m.Set(Int(int64(i-start)), from.Get(Int(int64(i))))
			}
			return m.Value()
		},
			"$f({t}: table) table", "\tcopy the entire table",
			"$f({t}: array, start: int, end: int) array", "\tcopy array, from start to end",
		),
		Str("parent"), Func1("parent", func(v Value) Value {
			return v.MustTable("").Parent().Value()
		}, "$f({t}: table) table", "\treturn table's parent if any"),
		Str("setparent"), Func2("setparent", func(v, p Value) Value {
			v.MustTable("").SetParent(p.MustTable(""))
			return Nil
		}, "$f({t}: table, p: table)", "\tset table's parent"),
		Str("setfirstparent"), Func2("setfirstparent", func(v, p Value) Value {
			v.MustTable("").SetFirstParent(p.MustTable(""))
			return Nil
		}, "$f({t}: table, p: table)", "\tinsert table's first parent"),
		Str("arraylen"), Func1("arraylen", func(v Value) Value {
			return Int(int64(len(v.MustTable("").items)))
		}, "$f({t}: array) int", "\treturn the true size of array (including trailing nils)"),
		Str("maplen"), Func1("maplen", func(v Value) Value {
			return Int(int64(len(v.MustTable("").hashItems)))
		}, "$f({t}: table) int", "\treturn the true size of table map (including empty entries)"),
		Str("keys"), Func1("keys", func(m Value) Value {
			a := make([]Value, 0)
			m.MustTable("").Foreach(func(k, v Value) bool { a = append(a, k); return true })
			return Array(a...)
		}, "$f({t}: table) array"),
		Str("values"), Func1("values", func(m Value) Value {
			a := make([]Value, 0)
			m.MustTable("").Foreach(func(k, v Value) bool { a = append(a, v); return true })
			return Array(a...)
		}, "$f({t}: table) array"),
		Str("items"), Func1("items", func(m Value) Value {
			a := make([]Value, 0)
			m.MustTable("").Foreach(func(k, v Value) bool { a = append(a, Array(k, v)); return true })
			return Array(a...)
		}, "$f({t}: table) array"),
		Str("foreach"), Func2("foreach", func(m, f Value) Value {
			m.MustTable("").Foreach(func(k, v Value) bool {
				v, err := f.MustFunc("").Call(k, v)
				panicErr(err)
				return v == Nil
			})
			return Nil
		}, "$f({t}: table, f: function)"),
		Str("contains"), Func2("contains", func(a, b Value) Value {
			found := false
			a.MustTable("").Foreach(func(k, v Value) bool { found = v.Equal(b); return !found })
			return Bool(found)
		}, "$f({t}: table, v: value)", "\ttest if table contains value"),
		Str("append"), Function("append", func(env *Env) {
			ma := env.Get(0).MustTable("")
			for i := 1; i < env.Size(); i++ {
				ma.Set(Int(int64(len(ma.items))), env.Get(i))
			}
			env.A = ma.Value()
		}, "$f({a}: array, ...args: value)", "\tappend values to array"),
		Str("filter"), Func2("filter", func(a, b Value) Value {
			ma := a.MustTable("")
			a2 := make([]Value, 0, ma.ArrayLen())
			ma.Foreach(func(k, v Value) bool {
				r, err := b.MustFunc("").Call(v)
				panicErr(err)
				if !r.IsFalse() {
					a2 = append(a2, v)
				}
				return true
			})
			return Array(a2...)
		}, "$f({a}: array, f: function)", "\tfilter out all values where f(value) is false"),
		Str("removeat"), Func2("removeat", func(a, b Value) Value {
			ma, idx := a.MustTable(""), b.MustInt("")
			if idx < 0 || idx >= int64(len(ma.items)) {
				return Nil
			}
			old := ma.items[idx]
			ma.items = append(ma.items[:idx], ma.items[idx+1:]...)
			if old != Nil {
				ma.count--
			}
			return old
		}, "$f({a}: array, index: int)", "\tremove value at index from array"),
		Str("concat"), Func2("concat", func(a, b Value) Value {
			ma, mb := a.MustTable(""), b.MustTable("")
			for _, b := range mb.ArrayPart() {
				ma.Set(Int(int64(len(ma.items))), b)
			}
			return ma.Value()
		}, "$f({array1}: array, array2: array)", "\tput elements from array2 to array1's end"),
		Str("merge"), Func2("merge", func(a, b Value) Value {
			ma, mb := a.MustTable(""), b.MustTable("")
			ma.resizeHash(len(mb.hashItems) + len(ma.hashItems))
			mb.Foreach(func(k, v Value) bool {
				ma.Set(k, v)
				return true
			})
			return ma.Value()
		}, "$f({table1}: table, table2: table)", "\tmerge elements from table2 to table1"),
		Str("tostring"), Func1("tostring", func(a Value) Value {
			p := &bytes.Buffer{}
			a.MustTable("").rawPrint(p, 0, false, true)
			return Bytes(p.Bytes())
		}, "$f({t}: table) string", "\tprint raw elements in table, ignore __str"),
		Str("name"), Func1("name", func(a Value) Value {
			return Str(a.MustTable("").Name())
		}, "$f({t}: table) string", "\tprint table's name"),
		Str("pure"), Func1("pure", func(m Value) Value {
			m2 := *m.MustTable("")
			m2.parent = nil
			return m2.Value()
		}, "$f({t}: table) table", "\treturn a new table who shares all the data from t, but with no parent"),
		Str("unwrap"), Func1("unwrap", func(m Value) Value {
			return ValRec(m.Interface())
		}, "unwrap(v: value) value", "\tunwrap Go's array, slice or map into table"),
	)
	AddGlobalValue("table", TableLib)

	StringMethods = TableMerge(StringMethods, Nil,
		Str("__name"), Str("strlib"),
		Str("__call"), Func2("str", func(strObj, src Value) Value {
			switch i := src.Interface().(type) {
			case []byte:
				return Bytes(i)
			default:
				return Str(fmt.Sprint(i))
			}
		}),
		Str("size"), Func1("size", func(src Value) Value {
			return Int(int64(len(src.MustStr(""))))
		}, "size({text}: string) int"),
		Str("len"), Func1("len", func(src Value) Value {
			return Int(int64(len(src.MustStr(""))))
		}, "len({text}: string) int"),
		Str("count"), Func1("count", func(src Value) Value {
			return Int(int64(utf8.RuneCountInString(src.MustStr(""))))
		}, "count({text}: string) int", "\treturn count of runes in text"),
		Str("from"), Func1("from", func(src Value) Value {
			return Str(fmt.Sprint(src.Interface()))
		}, "from(v: value) string", "\tconvert value to string"),
		Str("equals"), Func2("equals", func(src, a Value) Value {
			return Bool(src.MustStr("") == a.MustStr(""))
		}, "$f({text1}: string, text2: string) bool"),
		Str("iequals"), Func2("iequals", func(src, a Value) Value {
			return Bool(strings.EqualFold(src.MustStr(""), a.MustStr("")))
		}, "$f({text1}: string, text2: string) bool"),
		Str("contains"), Func2("contains", func(src, a Value) Value {
			return Bool(strings.Contains(src.MustStr(""), a.MustStr("")))
		}, "$f({text}: string, substr: string) bool"),
		Str("containsany"), Func2("containsany", func(src, a Value) Value {
			return Bool(strings.ContainsAny(src.MustStr(""), a.MustStr("")))
		}, "$f({text}: string, chars: string) bool"),
		Str("split"), Func3("split", func(src, delim, n Value) Value {
			s := src.MustStr("text")
			d := delim.MustStr("delimeter")
			r := []Value{}
			if n := n.MaybeInt(0); n == 0 {
				for _, p := range strings.Split(s, d) {
					r = append(r, Str(p))
				}
			} else {
				for _, p := range strings.SplitN(s, d, int(n)) {
					r = append(r, Str(p))
				}
			}
			return Array(r...)
		}, "split({text}: string, delim: string) array", "split({text}: string, delim: string, n: int) array"),
		Str("replace"), Function("replace", func(env *Env) {
			src := env.Get(0).MustStr("text")
			from := env.Get(1).MustStr("old text")
			to := env.Get(2).MustStr("new text")
			n := env.Get(3).MaybeInt(-1)
			env.A = Str(strings.Replace(src, from, to, int(n)))
		}, ""),
		Str("match"), Func2("match", func(pattern, str Value) Value {
			m, err := filepath.Match(pattern.MustStr("pattern"), str.MustStr("text"))
			panicErr(err)
			return Bool(m)
		}, ""),
		Str("find"), Func2("find", func(src, substr Value) Value {
			return Int(int64(strings.Index(src.MustStr(""), substr.MustStr(""))))
		}, "$f({text}: string, sub: string) int", "\tfind index of first appearence of sub in text"),
		Str("findany"), Func2("findany", func(src, substr Value) Value {
			return Int(int64(strings.IndexAny(src.MustStr(""), substr.MustStr(""))))
		}, "$f({text}: string, charset: string) int", "\tfind index of first appearence of any char from charset in text"),
		Str("rfind"), Func2("rfind", func(src, substr Value) Value {
			return Int(int64(strings.LastIndex(src.MustStr(""), substr.MustStr(""))))
		}, ""),
		Str("rfindany"), Func2("rfindany", func(src, substr Value) Value {
			return Int(int64(strings.LastIndexAny(src.MustStr(""), substr.MustStr(""))))
		}, "$f({text}: string, charset: string) int", "\tsame as findany(), but from right to left"),
		Str("sub"), Func3("sub", func(src, start, end Value) Value {
			s := src.MustStr("")
			st := start.MaybeInt(0)
			en := end.MaybeInt(int64(len(s)))
			for st < 0 && len(s) > 0 {
				st += int64(len(s))
			}
			for en < 0 && len(s) > 0 {
				en += int64(len(s))
			}
			return Str(s[st:en])
		}, "$f({text}: string, start: int, end: int) string"),
		Str("trim"), Func2("trim", func(src, cutset Value) Value {
			if cutset == Nil {
				return Str(strings.TrimSpace(src.MustStr("")))
			}
			return Str(strings.Trim(src.MustStr(""), cutset.MustStr("")))
		},
			"$f{text}: string) string", "\ttrim spaces at left and right side of text",
			"$f{text}: string, cutset: string) string", "\tremove chars both occurred in cutset and left-side/right-side of text"),
		Str("lremove"), Func2("lremove", func(src, cutset Value) Value {
			return Str(strings.TrimPrefix(src.MustStr(""), cutset.MustStr("")))
		}, "$f({text}: string, prefix: string) string", "\tremove prefix in text if any"),
		Str("rremove"), Func2("rremove", func(src, cutset Value) Value {
			return Str(strings.TrimSuffix(src.MustStr(""), cutset.MustStr("")))
		}, "$f({text}: string, suffix: string) string", "\tremove suffix in text if any"),
		Str("ltrim"), Func2("ltrim", func(src, cutset Value) Value {
			return Str(strings.TrimLeft(src.MustStr(""), cutset.MaybeStr(" ")))
		}, "$f({text}: string, cutset: string) string", "\tremove chars both ocurred in cutset and left-side of text"),
		Str("rtrim"), Func2("rtrim", func(src, cutset Value) Value {
			return Str(strings.TrimRight(src.MustStr(""), cutset.MaybeStr(" ")))
		}, "$f({text}: string, cutset: string) string", "\tremove chars both ocurred in cutset and right-side of text"),
		Str("decutf8"), Function("decutf8", func(env *Env) {
			r, sz := utf8.DecodeRuneInString(env.Get(0).MustStr(""))
			env.A = Array(Int(int64(r)), Int(int64(sz)))
		}, "$f({text}: string) array", "\tdecode first char in UTF-8 string, return { char_unicode, width_in_bytes }"),
		Str("startswith"), Func2("startswith", func(t, p Value) Value {
			return Bool(strings.HasPrefix(t.MustStr(""), p.MustStr("")))
		}, "$f(text: string, prefix: string) bool"),
		Str("endswith"), Func2("endswith", func(t, s Value) Value {
			return Bool(strings.HasSuffix(t.MustStr(""), s.MustStr("")))
		}, "$f(text: string, suffix: string) bool"),
		Str("upper"), Func1("upper", func(t Value) Value {
			return Str(strings.ToUpper(t.MustStr("")))
		}, "$f(s: string) string"),
		Str("lower"), Func1("lower", func(t Value) Value {
			return Str(strings.ToLower(t.MustStr("")))
		}, "$f(s: string) string"),
		Str("bytes"), Func1("bytes", func(s Value) Value {
			if s.Type() == typ.Number {
				return Val(make([]byte, s.Int()))
			}
			return Val([]byte(s.MustStr("")))
		},
			"$f(v: string) bytes", "\tcreate a byte array from the given string",
			"$f(n: int) bytes", "\tcreate an n-byte long array",
		),
		Str("chars"), Func2("chars", func(s, n Value) Value {
			var r []Value
			max := n.MaybeInt(0)
			for s := s.MustStr(""); len(s) > 0; {
				_, sz := utf8.DecodeRuneInString(s)
				if sz == 0 {
					break
				}
				r = append(r, Str(s[:sz]))
				if max > 0 && len(r) == int(max) {
					break
				}
				s = s[sz:]
			}
			return Array(r...)
		}, "chars({text}: string) array",
			"\tbreak a string into chars, e.g.: chars('a中c') => { 'a', '中', 'c' }",
			"chars({text}: string, n: int) array",
			"\tbreak a string into at most n chars, e.g.: chars('a中c', 1) => { 'a' }",
		),
		Str("format"), Function("format", func(env *Env) {
			buf := &bytes.Buffer{}
			sprintf(env, buf)
			env.A = Bytes(buf.Bytes())
		}, "format({pattern}: string, ...args: value) string"),
		Str("buffer"), Func1("buffer", func(v Value) Value {
			b := &bytes.Buffer{}
			if v != Nil {
				b.WriteString(v.String())
			}
			return TableProto(ReadWriterProto,
				Str("_f"), Val(b),
				Str("value"), Func1("value", func(a Value) Value {
					return Bytes(a.MustTable("").GetString("_f").Interface().(*bytes.Buffer).Bytes())
				}),
			)
		}),
	)

	AddGlobalValue("str", StringMethods)

	MathLib = TableMerge(MathLib, Nil,
		Str("__name"), Str("mathlib"),
		Str("INF"), Float(math.Inf(1)),
		Str("NEG_INF"), Float(math.Inf(-1)),
		Str("PI"), Float(math.Pi),
		Str("E"), Float(math.E),
		Str("randomseed"), Function("randomseed", func(env *Env) {
			rand.Seed(env.Get(0).MaybeInt(1))
		}, "randomseed(seed: int)"),
		Str("random"), Function("random", func(env *Env) {
			switch len(env.Stack()) {
			case 2:
				ai := env.Get(0).MustInt("")
				bi := env.Get(1).MustInt("")
				env.A = Int(rand.Int63n(bi-ai+1) + ai)
			case 1:
				env.A = Int(rand.Int63n(env.Get(0).MustInt("")))
			default:
				env.A = Float(rand.Float64())
			}
		},
			"$f() float", "\treturn [0, 1)",
			"$f(n: int) int", "\treturn [0, n)",
			"$f(a: int, b: int) int", "\treturn [a, b]"),
		Str("sqrt"), Func1("sqrt", func(v Value) Value { return Float(math.Sqrt(v.MustFloat(""))) }),
		Str("floor"), Func1("floor", func(v Value) Value { return Float(math.Floor(v.MustFloat(""))) }),
		Str("ceil"), Func1("ceil", func(v Value) Value { return Float(math.Ceil(v.MustFloat(""))) }),
		Str("min"), Function("min", func(env *Env) { mathMinMax(env, "#%d arg", false) }, "$f(...a: number) number"),
		Str("max"), Function("max", func(env *Env) { mathMinMax(env, "#%d arg", true) }, "$f(...a: number) number"),
		Str("pow"), Func2("pow", func(a, b Value) Value {
			ai := a.MustNum("base")
			bi := b.MustNum("power")
			if ai.IsInt() {
				return Int(ipow(ai.Int(), bi.Int()))
			}
			return Float(math.Pow(ai.Float(), bi.Float()))
		}, "pow(a: number, b: number) number"),
		Str("abs"), Function("abs", func(env *Env) {
			if i := env.Get(0).MustNum(""); i.IsInt() {
				if i := i.Int(); i < 0 {
					env.A = Int(-i)
				} else {
					env.A = Int(i)
				}
			} else {
				env.A = Float(math.Abs(i.Float()))
			}
		}),
		Str("remainder"), Function("remainder", func(env *Env) { env.A = Float(math.Remainder(env.Get(0).MustFloat(""), env.Get(1).MustFloat(""))) }),
		Str("mod"), Function("mod", func(env *Env) { env.A = Float(math.Mod(env.Get(0).MustFloat(""), env.Get(1).MustFloat(""))) }),
		Str("cos"), Function("cos", func(env *Env) { env.A = Float(math.Cos(env.Get(0).MustFloat(""))) }),
		Str("sin"), Function("sin", func(env *Env) { env.A = Float(math.Sin(env.Get(0).MustFloat(""))) }),
		Str("tan"), Function("tan", func(env *Env) { env.A = Float(math.Tan(env.Get(0).MustFloat(""))) }),
		Str("acos"), Function("acos", func(env *Env) { env.A = Float(math.Acos(env.Get(0).MustFloat(""))) }),
		Str("asin"), Function("asin", func(env *Env) { env.A = Float(math.Asin(env.Get(0).MustFloat(""))) }),
		Str("atan"), Function("atan", func(env *Env) { env.A = Float(math.Atan(env.Get(0).MustFloat(""))) }),
		Str("atan2"), Function("atan2", func(env *Env) { env.A = Float(math.Atan2(env.Get(0).MustFloat(""), env.Get(1).MustFloat(""))) }),
		Str("ldexp"), Function("ldexp", func(env *Env) { env.A = Float(math.Ldexp(env.Get(0).MustFloat(""), int(env.Get(1).MaybeInt(0)))) }),
		Str("modf"), Function("modf", func(env *Env) {
			a, b := math.Modf(env.Get(0).MustFloat(""))
			env.A = Array(Float(a), Float(b))
		}),
	)
	AddGlobalValue("math", MathLib)

	OSLib = TableMerge(OSLib, Nil,
		Str("__name"), Str("oslib"),
		Str("args"), ValRec(os.Args),
		Str("environ"), Function("environ", func(env *Env) { env.A = ValRec(os.Environ()) }),
		Str("shell"), Func2("shell", func(cmd, opt Value) Value {
			timeout := time.Duration(1 << 62) // basically forever
			if tmp := opt.MaybeTableGetString("timeout"); tmp != Nil {
				timeout = time.Duration(tmp.MustFloat("timeout") * float64(time.Second))
			}

			out := make(chan interface{})
			p := exec.Command("sh", "-c", cmd.MustStr(""))

			if tmp := opt.MaybeTableGetString("env"); tmp != Nil {
				tmp.MustTable("env").Foreach(func(k, v Value) bool {
					p.Env = append(p.Env, k.String()+"="+v.String())
					return true
				})
			}
			go func() {
				v, err := p.Output()
				if err != nil {
					out <- err
				} else {
					out <- v
				}
			}()
			select {
			case r := <-out:
				if buf, ok := r.([]byte); ok {
					return Bytes(buf)
				}
				panic(r)
			case <-time.After(timeout):
				p.Process.Kill()
				panic("timeout")
			}
		}),
		Str("readdir"), Func1("readdir", func(path Value) Value {
			p := path.MustStr("")
			fi, err := ioutil.ReadDir(p)
			panicErr(err)
			return ValRec(fi)
		}),
		Str("remove"), Func1("remove", func(path Value) Value {
			p := path.MustStr("")
			fi, err := os.Stat(p)
			panicErr(err)
			if fi.IsDir() {
				err = os.RemoveAll(p)
			} else {
				err = os.Remove(p)
			}
			panicErr(err)
			return Nil
		}),
	)
	AddGlobalValue("os", OSLib)
}

func mathMinMax(env *Env, msg string, max bool) {
	if len(env.Stack()) <= 0 {
		return
	}
	v := env.Get(0).mustBe(typ.Number, msg, 1)
	if v.IsInt() {
		vi := v.Int()
		for ii := 1; ii < len(env.Stack()); ii++ {
			if x := env.Get(ii).mustBe(typ.Number, msg, ii+1).Int(); x >= vi == max {
				vi = x
			}
		}
		env.A = Int(vi)
	} else {
		vf := v.Float()
		for i := 1; i < len(env.Stack()); i++ {
			if x := env.Get(i).mustBe(typ.Number, msg, i+1).Float(); x >= vf == max {
				vf = x
			}
		}
		env.A = Float(vf)
	}
}

func ipow(base, exp int64) int64 {
	var result int64 = 1
	for {
		if exp&1 == 1 {
			result *= base
		}
		exp >>= 1
		if exp == 0 {
			break
		}
		base *= base
	}
	return result
}

func sprintf(env *Env, p io.Writer) {
	f := env.Get(0).MustStr("")
	tmp := bytes.Buffer{}
	popi := 0
	for len(f) > 0 {
		idx := strings.Index(f, "%")
		if idx == -1 {
			fmt.Fprint(p, f)
			break
		} else if idx == len(f)-1 {
			panicf("unexpected '%%' at end")
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
				panicf("unexpected verb: '%c'", f[0])
			}
			tmp.WriteByte(f[0])
			f = f[1:]
		}

		popi++
		pop := env.Get(popi)
		switch expecting {
		case typ.Bool:
			fmt.Fprint(p, !pop.IsFalse())
		case typ.String:
			fmt.Fprintf(p, tmp.String(), pop.String())
		case typ.Number:
			if pop.mustBe(typ.Number, "arg #%d", popi-1).IsInt() {
				fmt.Fprintf(p, tmp.String(), pop.Int())
			} else {
				fmt.Fprintf(p, tmp.String(), pop.Float())
			}
		case typ.Number + typ.String:
			if pop.Type() == typ.String {
				fmt.Fprintf(p, tmp.String(), pop.Str())
			} else if pop.mustBe(typ.Number, "arg #%d", popi-1).IsInt() {
				fmt.Fprintf(p, tmp.String(), pop.Int())
			} else {
				fmt.Fprintf(p, tmp.String(), pop.Float())
			}
		case typ.Native:
			fmt.Fprint(p, pop.Interface())
		}
	}
}
