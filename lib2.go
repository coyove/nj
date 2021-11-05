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
	IOLib = MapAdd(IOLib,
		Str("reader"), ReaderProto.Value(),
		Str("writer"), WriterProto.Value(),
		Str("seeker"), SeekerProto.Value(),
		Str("closer"), CloserProto.Value(),
		Str("readwriter"), ReadWriter.Value(),
		Str("readcloser"), ReadCloser.Value(),
		Str("writecloser"), WriteCloser.Value(),
		Str("readwritecloser"), ReadWriteCloser.Value(),
		Str("readwriteseekcloser"), ReadWriteSeekCloser.Value(),
	)
	AddGlobalValue("io", IOLib)

	TableLib = MapAdd(TableLib,
		Str("__call"), Native2("__call", func(env *Env, t, m Value) Value {
			if m == Nil {
				return TableProto(t.MustTable(""))
			}
			m.MustTable("").SetParent(t.MustTable(""))
			return m
		}),
		Str("concurrent"), Native2("concurrent", func(env *Env, t, m Value) Value {
			x := NewTable(t.MustTable("").Len())
			t.MustTable("").Foreach(func(k, v Value) bool {
				if v.Type() == typ.Func {
					old := v.Func()
					v = Native(v.Func().Name, func(env *Env) {
						mu := env.Get(0).MustTable("").GetString("_mu").Interface().(*sync.Mutex)
						mu.Lock()
						defer mu.Unlock()
						a, err := old.Call(env.Stack()...)
						if err != nil {
							panic(err)
						}
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
			m.MustTable("").SetParent(x)
			return m
		}, "$f(src: table) table", "\tcreate a concurrently accessible table"),
		Str("get"), Native2("get", func(env *Env, m, k Value) Value {
			return m.MustTable("").Get(k)
		}, "$f({t}: table, k: value) value"),
		Str("set"), Native3("set", func(env *Env, m, k, v Value) Value {
			m.MustTable("").Set(k, v)
			return m
		}, "$f({t}: table, k: value, v: value) table"),
		Str("make"), Native1("make", func(env *Env, n Value) Value {
			return NewTable(int(n.MustInt(""))).Value()
		}, "$f(n: int) table", "\treturn a table, preallocate enough hash space for n values"),
		Str("makearray"), Native1("makearray", func(env *Env, n Value) Value {
			return Array(make([]Value, n.MustInt(""))...)
		}, "$f(n: int) array", "\treturn a table array, preallocate space for n values"),
		Str("clear"), Native1("clear", func(env *Env, m Value) Value {
			m.MustTable("").Clear()
			return Nil
		}, "$f({t}: table)"),
		Str("slice"), Native3("slice", func(env *Env, t, s, e Value) Value {
			start, end := int(s.MustInt("")), int(e.MustInt(""))
			return Array(t.MustTable("").items[start:end]...)
		}, "$f({t}: array, start: int, end: int) array", "\tslice array, from start to end"),
		Str("copy"), Native3("copy", func(env *Env, t, s, e Value) Value {
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
		Str("parent"), Native1("parent", func(env *Env, v Value) Value {
			return v.MustTable("").Parent().Value()
		}, "$f({t}: table) table", "\treturn table's parent if any"),
		Str("setparent"), Native2("setparent", func(env *Env, v, p Value) Value {
			v.MustTable("").SetParent(p.MustTable(""))
			return Nil
		}, "$f({t}: table, p: table)", "\tset table's parent"),
		Str("arraylen"), Native1("arraylen", func(env *Env, v Value) Value {
			return Int(int64(len(v.MustTable("").items)))
		}, "$f({t}: array) int", "\treturn the true size of array (including trailing nils)"),
		Str("maplen"), Native1("maplen", func(env *Env, v Value) Value {
			return Int(int64(len(v.MustTable("").hashItems)))
		}, "$f({t}: table) int", "\treturn the true size of table map (including empty entries)"),
		Str("keys"), Native1("keys", func(env *Env, m Value) Value {
			a := make([]Value, 0)
			m.MustTable("").Foreach(func(k, v Value) bool { a = append(a, k); return true })
			return Array(a...)
		}, "$f({t}: table) array"),
		Str("values"), Native1("values", func(env *Env, m Value) Value {
			a := make([]Value, 0)
			m.MustTable("").Foreach(func(k, v Value) bool { a = append(a, v); return true })
			return Array(a...)
		}, "$f({t}: table) array"),
		Str("items"), Native1("items", func(env *Env, m Value) Value {
			a := make([]Value, 0)
			m.MustTable("").Foreach(func(k, v Value) bool { a = append(a, Array(k, v)); return true })
			return Array(a...)
		}, "$f({t}: table) array"),
		Str("foreach"), Native2("foreach", func(env *Env, m, f Value) Value {
			m.MustTable("").Foreach(func(k, v Value) bool {
				v, err := f.MustFunc("").Call(k, v)
				if err != nil {
					panic(err)
				}
				return v == Nil
			})
			return Nil
		}, "$f({t}: table, f: function)"),
		Str("contains"), Native2("contains", func(env *Env, a, b Value) Value {
			found := false
			a.MustTable("").Foreach(func(k, v Value) bool { found = v.Equal(b); return !found })
			return Bool(found)
		}, "$f({t}: table, v: value)", "\ttest if table contains value"),
		Str("append"), Native("append", func(env *Env) {
			ma := env.Get(0).MustTable("")
			for i := 1; i < env.Size(); i++ {
				ma.Set(Int(int64(len(ma.items))), env.Get(i))
			}
			env.A = ma.Value()
		}, "$f({a}: array, ...args: value)", "\tappend values to array"),
		Str("filter"), Native2("filter", func(env *Env, a, b Value) Value {
			ma := a.MustTable("")
			a2 := make([]Value, 0, ma.ArrayLen())
			ma.Foreach(func(k, v Value) bool {
				r, err := b.MustFunc("").Call(v)
				if err != nil {
					panic(err)
				}
				if !r.IsFalse() {
					a2 = append(a2, v)
				}
				return true
			})
			return Array(a2...)
		}, "$f({a}: array, f: function)", "\tfilter out all values where f(value) is false"),
		Str("removeat"), Native2("removeat", func(env *Env, a, b Value) Value {
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
		Str("concat"), Native2("concat", func(env *Env, a, b Value) Value {
			ma, mb := a.MustTable(""), b.MustTable("")
			for _, b := range mb.ArrayPart() {
				ma.Set(Int(int64(len(ma.items))), b)
			}
			return ma.Value()
		}, "$f({array1}: array, array2: array)", "\tput elements from array2 to array1's end"),
		Str("merge"), Native2("merge", func(env *Env, a, b Value) Value {
			ma, mb := a.MustTable(""), b.MustTable("")
			ma.resizeHash(len(mb.hashItems) + len(ma.hashItems))
			mb.Foreach(func(k, v Value) bool {
				ma.Set(k, v)
				return true
			})
			return ma.Value()
		}, "$f({table1}: table, table2: table)", "\tmerge elements from table2 to table1"),
		Str("tostring"), Native1("print", func(env *Env, a Value) Value {
			p := &bytes.Buffer{}
			a.MustTable("").rawPrint(p, 0, false)
			return Bytes(p.Bytes())
		}, "$f({t}: table) string", "\tprint raw elements in table, ignore __str"),
	)
	AddGlobalValue("table", TableLib)

	StringMethods = MapAdd(StringMethods,
		Str("__call"), Native2("str", func(env *Env, strObj, src Value) Value {
			switch i := src.Interface().(type) {
			case []byte:
				return Bytes(i)
			default:
				return Str(fmt.Sprint(i))
			}
		}),
		Str("size"), Native1("size", func(env *Env, src Value) Value {
			return Int(int64(len(src.MustStr(""))))
		}, "size({text}: string) int"),
		Str("len"), Native1("len", func(env *Env, src Value) Value {
			return Int(int64(len(src.MustStr(""))))
		}, "len({text}: string) int"),
		Str("count"), Native1("count", func(env *Env, src Value) Value {
			return Int(int64(utf8.RuneCountInString(src.MustStr(""))))
		}, "count({text}: string) int", "\treturn count of runes in text"),
		Str("from"), Native1("from", func(env *Env, src Value) Value {
			return Str(fmt.Sprint(src.Interface()))
		}, "from(v: value) string", "\tconvert value to string"),
		Str("equals"), Native2("equals", func(env *Env, src, a Value) Value {
			return Bool(src.MustStr("") == a.MustStr(""))
		}, "$f({text1}: string, text2: string) bool"),
		Str("iequals"), Native2("iequals", func(env *Env, src, a Value) Value {
			return Bool(strings.EqualFold(src.MustStr(""), a.MustStr("")))
		}, "$f({text1}: string, text2: string) bool"),
		Str("contains"), Native2("contains", func(env *Env, src, a Value) Value {
			return Bool(strings.Contains(src.MustStr(""), a.MustStr("")))
		}, "$f({text}: string, substr: string) bool"),
		Str("containsany"), Native2("containsany", func(env *Env, src, a Value) Value {
			return Bool(strings.ContainsAny(src.MustStr(""), a.MustStr("")))
		}, "$f({text}: string, chars: string) bool"),
		Str("split"), Native3("split", func(env *Env, src, delim, n Value) Value {
			s := src.MustStr("text")
			d := delim.MustStr("delimeter")
			r := []Value{}
			if n := n.IntDefault(0); n == 0 {
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
		Str("replace"), Native("replace", func(env *Env) {
			src := env.Get(0).MustStr("text")
			from := env.Get(1).MustStr("old text")
			to := env.Get(2).MustStr("new text")
			n := env.Get(3).IntDefault(-1)
			env.A = Str(strings.Replace(src, from, to, int(n)))
		}, ""),
		Str("match"), Native2("match", func(env *Env, pattern, str Value) Value {
			m, err := filepath.Match(pattern.MustStr("pattern"), str.MustStr("text"))
			if err != nil {
				panic(err)
			}
			return Bool(m)
		}, ""),
		Str("find"), Native2("find", func(env *Env, src, substr Value) Value {
			return Int(int64(strings.Index(src.MustStr(""), substr.MustStr(""))))
		}, "$f({text}: string, sub: string) int", "\tfind index of first appearence of sub in text"),
		Str("findany"), Native2("findany", func(env *Env, src, substr Value) Value {
			return Int(int64(strings.IndexAny(src.MustStr(""), substr.MustStr(""))))
		}, "$f({text}: string, charset: string) int", "\tfind index of first appearence of any char from charset in text"),
		Str("rfind"), Native2("rfind", func(env *Env, src, substr Value) Value {
			return Int(int64(strings.LastIndex(src.MustStr(""), substr.MustStr(""))))
		}, ""),
		Str("rfindany"), Native2("rfindany", func(env *Env, src, substr Value) Value {
			return Int(int64(strings.LastIndexAny(src.MustStr(""), substr.MustStr(""))))
		}, "$f({text}: string, charset: string) int", "\tsame as findany(), but from right to left"),
		Str("sub"), Native3("sub", func(env *Env, src, start, end Value) Value {
			s := src.MustStr("")
			st := start.IntDefault(0)
			en := end.IntDefault(int64(len(s)))
			for st < 0 && len(s) > 0 {
				st += int64(len(s))
			}
			for en < 0 && len(s) > 0 {
				en += int64(len(s))
			}
			return Str(s[st:en])
		}, "$f({text}: string, start: int, end: int) string"),
		Str("trim"), Native2("trim", func(env *Env, src, cutset Value) Value {
			if cutset == Nil {
				return Str(strings.TrimSpace(src.MustStr("")))
			}
			return Str(strings.Trim(src.MustStr(""), cutset.MustStr("")))
		},
			"$f{text}: string) string", "\ttrim spaces at left and right side of text",
			"$f{text}: string, cutset: string) string", "\tremove chars both occurred in cutset and left-side/right-side of text"),
		Str("lremove"), Native2("lremove", func(env *Env, src, cutset Value) Value {
			return Str(strings.TrimPrefix(src.MustStr(""), cutset.MustStr("")))
		}, "$f({text}: string, prefix: string) string", "\tremove prefix in text if any"),
		Str("rremove"), Native2("rremove", func(env *Env, src, cutset Value) Value {
			return Str(strings.TrimSuffix(src.MustStr(""), cutset.MustStr("")))
		}, "$f({text}: string, suffix: string) string", "\tremove suffix in text if any"),
		Str("ltrim"), Native2("ltrim", func(env *Env, src, cutset Value) Value {
			return Str(strings.TrimLeft(src.MustStr(""), cutset.StringDefault(" ")))
		}, "$f({text}: string, cutset: string) string", "\tremove chars both ocurred in cutset and left-side of text"),
		Str("rtrim"), Native2("rtrim", func(env *Env, src, cutset Value) Value {
			return Str(strings.TrimRight(src.MustStr(""), cutset.StringDefault(" ")))
		}, "$f({text}: string, cutset: string) string", "\tremove chars both ocurred in cutset and right-side of text"),
		Str("decutf8"), Native("decutf8", func(env *Env) {
			r, sz := utf8.DecodeRuneInString(env.Get(0).MustStr(""))
			env.A = Array(Int(int64(r)), Int(int64(sz)))
		}, "$f({text}: string) array", "\tdecode first char in UTF-8 string, return { char_unicode, width_in_bytes }"),
		Str("startswith"), Native2("startswith", func(env *Env, t, p Value) Value {
			return Bool(strings.HasPrefix(t.MustStr(""), p.MustStr("")))
		}, "$f(text: string, prefix: string) bool"),
		Str("endswith"), Native2("endswith", func(env *Env, t, s Value) Value {
			return Bool(strings.HasSuffix(t.MustStr(""), s.MustStr("")))
		}, "$f(text: string, suffix: string) bool"),
		Str("upper"), Native1("upper", func(env *Env, t Value) Value {
			return Str(strings.ToUpper(t.MustStr("")))
		}, "$f(s: string) string"),
		Str("lower"), Native1("lower", func(env *Env, t Value) Value {
			return Str(strings.ToLower(t.MustStr("")))
		}, "$f(s: string) string"),
		Str("bytes"), Native1("bytes", func(env *Env, s Value) Value {
			if s.Type() == typ.Number {
				return Val(make([]byte, s.Int()))
			}
			return Val([]byte(s.MustStr("")))
		},
			"$f(v: string) bytes", "\tcreate a byte array from the given string",
			"$f(n: int) bytes", "\tcreate an n-byte long array",
		),
		Str("chars"), Native2("chars", func(env *Env, s, n Value) Value {
			var r []Value
			max := n.IntDefault(0)
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
		Str("format"), Native("format", func(env *Env) {
			buf := &bytes.Buffer{}
			sprintf(env, buf)
			env.A = Bytes(buf.Bytes())
		}, "format({pattern}: string, ...args: value) string"),
		Str("buffer"), Native1("buffer", func(env *Env, v Value) Value {
			b := &bytes.Buffer{}
			if v != Nil {
				b.WriteString(v.String())
			}
			return TableProto(ReadWriter,
				Str("_f"), Val(b),
				Str("value"), Native1("value", func(env *Env, a Value) Value {
					return Bytes(a.MustTable("").GetString("_f").Interface().(*bytes.Buffer).Bytes())
				}),
			)
		}),
	)

	AddGlobalValue("str", StringMethods)

	var rg = struct {
		sync.Mutex
		*rand.Rand
	}{Rand: rand.New(rand.NewSource(1))}

	MathLib = MapAdd(MathLib,
		Str("INF"), Float(math.Inf(1)),
		Str("NEG_INF"), Float(math.Inf(-1)),
		Str("PI"), Float(math.Pi),
		Str("E"), Float(math.E),
		Str("randomseed"), Native("randomseed", func(env *Env) {
			rg.Lock()
			defer rg.Unlock()
			rg.Rand.Seed(env.Get(0).IntDefault(1))
		}, "randomseed(seed: int)"),
		Str("random"), Native("random", func(env *Env) {
			rg.Lock()
			defer rg.Unlock()
			switch len(env.Stack()) {
			case 2:
				af, ai, aIsInt := env.Get(0).MustNum("").Num()
				bf, bi, bIsInt := env.Get(1).MustNum("").Num()
				if aIsInt && bIsInt {
					env.A = Int(int64(rg.Intn(int(bi-ai+1))) + ai)
				} else {
					env.A = Float(rg.Float64()*(bf-af) + af)
				}
			case 1:
				env.A = Int(int64(rg.Intn(int(env.Get(0).MustNum("").Int()))))
			default:
				env.A = Float(rg.Float64())
			}
		},
			"$f() float", "\treturn [0, 1)",
			"$f(n: int) int", "\treturn [0, n)",
			"$f(a: number, b: number) number", "\treturn [a, b]"),
		Str("sqrt"), Native1("sqrt", func(env *Env, v Value) Value { return Float(math.Sqrt(v.MustFloat(""))) }),
		Str("floor"), Native1("floor", func(env *Env, v Value) Value { return Float(math.Floor(v.MustFloat(""))) }),
		Str("ceil"), Native1("ceil", func(env *Env, v Value) Value { return Float(math.Ceil(v.MustFloat(""))) }),
		Str("min"), Native("min", func(env *Env) { mathMinMax(env, "#%d arg", false) }, "$f(...a: number) number"),
		Str("max"), Native("max", func(env *Env) { mathMinMax(env, "#%d arg", true) }, "$f(...a: number) number"),
		Str("pow"), Native2("pow", func(env *Env, a, b Value) Value {
			af, ai, aIsInt := a.MustNum("base").Num()
			bf, bi, bIsInt := b.MustNum("power").Num()
			if aIsInt && bIsInt {
				return Int(ipow(ai, bi))
			}
			return Float(math.Pow(af, bf))
		}, "pow(a: number, b: number) number"),
		Str("abs"), Native("abs", func(env *Env) {
			switch f, i, isInt := env.Get(0).MustNum("").Num(); {
			case isInt && i < 0:
				env.A = Int(-i)
			case isInt && i >= 0:
				env.A = Int(i)
			default:
				env.A = Float(math.Abs(f))
			}
		}),
		Str("remainder"), Native("remainder", func(env *Env) { env.A = Float(math.Remainder(env.Get(0).MustFloat(""), env.Get(1).MustFloat(""))) }),
		Str("mod"), Native("mod", func(env *Env) { env.A = Float(math.Mod(env.Get(0).MustFloat(""), env.Get(1).MustFloat(""))) }),
		Str("cos"), Native("cos", func(env *Env) { env.A = Float(math.Cos(env.Get(0).MustFloat(""))) }),
		Str("sin"), Native("sin", func(env *Env) { env.A = Float(math.Sin(env.Get(0).MustFloat(""))) }),
		Str("tan"), Native("tan", func(env *Env) { env.A = Float(math.Tan(env.Get(0).MustFloat(""))) }),
		Str("acos"), Native("acos", func(env *Env) { env.A = Float(math.Acos(env.Get(0).MustFloat(""))) }),
		Str("asin"), Native("asin", func(env *Env) { env.A = Float(math.Asin(env.Get(0).MustFloat(""))) }),
		Str("atan"), Native("atan", func(env *Env) { env.A = Float(math.Atan(env.Get(0).MustFloat(""))) }),
		Str("atan2"), Native("atan2", func(env *Env) { env.A = Float(math.Atan2(env.Get(0).MustFloat(""), env.Get(1).MustFloat(""))) }),
		Str("ldexp"), Native("ldexp", func(env *Env) { env.A = Float(math.Ldexp(env.Get(0).MustFloat(""), int(env.Get(1).IntDefault(0)))) }),
		Str("modf"), Native("modf", func(env *Env) {
			a, b := math.Modf(env.Get(0).MustFloat(""))
			env.A = Array(Float(a), Float(b))
		}),
	)
	AddGlobalValue("math", MathLib)

	OSLib = MapAdd(OSLib,
		Str("args"), ValRec(os.Args),
		Str("environ"), Native("environ", func(env *Env) { env.A = ValRec(os.Environ()) }),
		Str("shell"), Native2("shell", func(env *Env, cmd, opt Value) Value {
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
		Str("readdir"), Native1("readdir", func(env *Env, path Value) Value {
			p := path.MustStr("")
			fi, err := ioutil.ReadDir(p)
			if err != nil {
				panic(err)
			}
			return ValRec(fi)
		}),
		Str("remove"), Native1("remove", func(env *Env, path Value) Value {
			p := path.MustStr("")
			fi, err := os.Stat(p)
			if err != nil {
				panic(err)
			}
			if fi.IsDir() {
				err = os.RemoveAll(p)
			} else {
				err = os.Remove(p)
			}
			if err != nil {
				panic(err)
			}
			return Nil
		}),
	)
	AddGlobalValue("os", OSLib)
}

func mathMinMax(env *Env, msg string, max bool) {
	if len(env.Stack()) <= 0 {
		return
	}
	f, i, isInt := env.Get(0).mustBe(typ.Number, msg, 1).Num()
	if isInt {
		for ii := 1; ii < len(env.Stack()); ii++ {
			if x := env.Get(ii).mustBe(typ.Number, msg, ii+1).Int(); x >= i == max {
				i = x
			}
		}
		env.A = Int(i)
	} else {
		for i := 1; i < len(env.Stack()); i++ {
			if x, _, _ := env.Get(i).mustBe(typ.Number, msg, i+1).Num(); x >= f == max {
				f = x
			}
		}
		env.A = Float(f)
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
	pop := func() Value { popi++; return env.Get(popi) }
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
				expecting = typ.Interface
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', '-', '+', '#', ' ':
			default:
				panicf("unexpected verb: '%c'", f[0])
			}
			tmp.WriteByte(f[0])
			f = f[1:]
		}
		switch expecting {
		case typ.String:
			fmt.Fprintf(p, tmp.String(), pop().String())
		case typ.Number:
			f, i, isInt := pop().Num()
			if isInt {
				fmt.Fprintf(p, tmp.String(), i)
			} else {
				fmt.Fprintf(p, tmp.String(), f)
			}
		case typ.Number + typ.String:
			v := pop()
			if v.Type() == typ.String {
				fmt.Fprintf(p, tmp.String(), v.Str())
			} else {
				f, i, isInt := pop().Num()
				if isInt {
					fmt.Fprintf(p, tmp.String(), i)
				} else {
					fmt.Fprintf(p, tmp.String(), f)
				}
			}
		case typ.Interface:
			fmt.Fprint(p, pop())
		}
	}
}
