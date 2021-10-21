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
	"unicode/utf8"

	"github.com/coyove/script/typ"
)

var StringMethods, MathLib, TableLib, OSLib Value

func init() {
	TableLib = MapAdd(TableLib,
		Str("makearray"), Native1("makearray", func(env *Env, n Value) Value {
			a := Array(make([]Value, n.MustInt(""))...)
			a.Table().count = 0
			return a
		}, "makearray(n) => { nil, ..., nil }", "\treturn a table array, preallocate space for n values"),
		Str("arraylen"), Native1("arraylen", func(env *Env, v Value) Value { return Int(int64(len(v.MustMap("").items))) }),
		Str("maplen"), Native1("maplen", func(env *Env, v Value) Value { return Int(int64(len(v.MustMap("").hashItems))) }),
		Str("keys"), Native1("keys", func(env *Env, m Value) Value {
			a := make([]Value, 0)
			m.MustMap("").Foreach(func(k, v Value) bool {
				a = append(a, k)
				return true
			})
			return Array(a...)
		}),
		Str("append"), Native2("append", func(env *Env, m, v Value) Value {
			a := m.MustMap("")
			a.Set(Int(int64(len(a.items))), v)
			return m
		}, "append(array, value)", "\tappend value to table array"),
		Str("concat"), Native2("concat", func(env *Env, a, b Value) Value {
			ma, mb := a.MustMap(""), b.MustMap("")
			for _, b := range mb.ArrayPart() {
				ma.Set(Int(int64(len(ma.items))), b)
			}
			return ma.Value()
		}, "concat(array1, array2)", "\tput elements from array2 to array1's end"),
	)
	AddGlobalValue("table", TableLib)

	StringMethods = MapAdd(StringMethods,
		Str("__call"), Native1("str", func(env *Env, src Value) Value {
			return Str(fmt.Sprint(src.Interface()))
		}, ""),
		Str("from"), Native1("from", func(env *Env, src Value) Value {
			return Str(fmt.Sprint(src.Interface()))
		}, ""),
		Str("iequal"), Native2("iequal", func(env *Env, src, a Value) Value {
			s := src.MustStr("")
			return Bool(strings.EqualFold(s, a.MustStr("")))
		}, ""),
		Str("contains"), Native2("contains", func(env *Env, src, a Value) Value {
			s := src.MustStr("")
			return Bool(strings.Contains(s, a.MustStr("")))
		}, ""),
		Str("containsany"), Native2("containsany", func(env *Env, src, a Value) Value {
			s := src.MustStr("")
			return Bool(strings.ContainsAny(s, a.MustStr("")))
		}, ""),
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
		}, "split(text, delim) => {part1, part2, ...}", "split(text, delim, n) => {part1, ..., partN}"),
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
			s := src.MustStr("")
			return Int(int64(strings.Index(s, substr.MustStr(""))))
		}, ""),
		Str("findany"), Native2("findany", func(env *Env, src, substr Value) Value {
			s := src.MustStr("")
			return Int(int64(strings.IndexAny(s, substr.MustStr(""))))
		}, ""),
		Str("rfind"), Native2("rfind", func(env *Env, src, substr Value) Value {
			s := src.MustStr("")
			return Int(int64(strings.LastIndex(s, substr.MustStr(""))))
		}, ""),
		Str("rfindany"), Native2("rfindany", func(env *Env, src, substr Value) Value {
			s := src.MustStr("")
			return Int(int64(strings.LastIndexAny(s, substr.MustStr(""))))
		}, ""),
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
		}, "sub(text, start, end) => text[start:end]"),
		Str("trim"), Native2("trim", func(env *Env, src, cutset Value) Value {
			if cutset == Nil {
				return Str(strings.TrimSpace(src.MustStr("")))
			}
			c := cutset.MustStr("")
			return Str(strings.Trim(src.MustStr(""), c))
		}, "trim(' text ') => 'text'", "trim('text', 'tx') => 'e'"),
		Str("ltrim"), Native2("ltrim", func(env *Env, src, cutset Value) Value {
			c := cutset.MustStr("")
			return Str(strings.TrimLeft(src.MustStr(""), c))
		}, "ltrim('abcdtext', 'abxy') => 'cdtext'"),
		Str("rtrim"), Native2("rtrim", func(env *Env, src, cutset Value) Value {
			c := cutset.MustStr("")
			return Str(strings.TrimRight(src.MustStr(""), c))
		}, "rtrim('textabcd', 'cdxy') => 'textab'"),
		Str("ptrim"), Native2("ptrim", func(env *Env, src, cutset Value) Value {
			c := cutset.MustStr("")
			return Str(strings.TrimPrefix(src.MustStr(""), c))
		}, "ptrim('prefixtext', 'prefix') => 'text'"),
		Str("strim"), Native2("strim", func(env *Env, src, cutset Value) Value {
			c := cutset.MustStr("")
			return Str(strings.TrimSuffix(src.MustStr(""), c))
		}, "strim('textsuffix', 'suffix') => 'text'"),
		Str("decutf8"), Native("decutf8", func(env *Env) {
			r, sz := utf8.DecodeRuneInString(env.Get(0).MustStr(""))
			env.A = Array(Int(int64(r)), Int(int64(sz)))
		}, "$f(string) => { char_unicode, width_in_bytes }"),
		Str("startswith"), Native2("startswith", func(env *Env, t, p Value) Value {
			return Bool(strings.HasPrefix(t.MustStr(""), p.MustStr("")))
		}, "startswith(text, prefix) => bool"),
		Str("endswith"), Native2("endswith", func(env *Env, t, s Value) Value {
			return Bool(strings.HasSuffix(t.MustStr(""), s.MustStr("")))
		}, "endswith(text, suffix) => bool"),
		Str("upper"), Native1("upper", func(env *Env, t Value) Value {
			return Str(strings.ToUpper(t.MustStr("")))
		}, "$f('text') => 'TEXT'"),
		Str("lower"), Native1("lower", func(env *Env, t Value) Value {
			return Str(strings.ToLower(t.MustStr("")))
		}, "$f('TEXT') => 'text'"),
		Str("bytes"), Native2("bytes", func(env *Env, s, n Value) Value {
			sz := s.MustStr("")
			var r []byte
			if max := n.IntDefault(-1); max >= 0 && len(sz) > int(max) {
				r = []byte(sz[:max])
			} else {
				r = []byte(sz)
			}
			return Bytes(r)
		}),
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
		}, "chars(string) => { char1, char2, ... }",
			"chars(string, max) => { char1, char2, ..., char_max }",
			"\tbreak a string into (at most 'max') chars, e.g.:",
			"\tchars('a中c') => { 'a', '中', 'c' }",
			"\tchars('a中c', 1) => { 'a' }",
		),
		Str("format"), Native("format", func(env *Env) {
			f := env.Get(0).MustStr("")
			p, tmp := bytes.Buffer{}, bytes.Buffer{}
			popi := 0
			pop := func() Value { popi++; return env.Get(popi) }
			for len(f) > 0 {
				idx := strings.Index(f, "%")
				if idx == -1 {
					p.WriteString(f)
					break
				} else if idx == len(f)-1 {
					panicf("unexpected '%%' at end")
				}
				p.WriteString(f[:idx])
				if f[idx+1] == '%' {
					p.WriteString("%")
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
					p.WriteString(fmt.Sprintf(tmp.String(), pop().String()))
				case typ.Number:
					f, i, isInt := pop().Num()
					if isInt {
						fmt.Fprintf(&p, tmp.String(), i)
					} else {
						fmt.Fprintf(&p, tmp.String(), f)
					}
				case typ.Number + typ.String:
					v := pop()
					if v.Type() == typ.String {
						fmt.Fprintf(&p, tmp.String(), v.Str())
					} else {
						f, i, isInt := pop().Num()
						if isInt {
							fmt.Fprintf(&p, tmp.String(), i)
						} else {
							fmt.Fprintf(&p, tmp.String(), f)
						}
					}
				case typ.Interface:
					fmt.Fprint(&p, pop())
				}
			}
			env.A = Str(p.String())
		}, "format(pattern, a1, a2, ...)"),
		Str("buffer"), Native1("buffer", func(env *Env, v Value) Value {
			b := &bytes.Buffer{}
			if v != Nil {
				b.WriteString(v.String())
			}
			p := Map(
				Str("_buf"), Val(b),
				Str("value"), Native1("value", func(env *Env, a Value) Value {
					return Bytes(a.MustMap("").GetString("_buf").Interface().(*bytes.Buffer).Bytes())
				}),
				Str("write"), Native2("write", func(env *Env, a, b Value) Value {
					a.MustMap("").GetString("_buf").Interface().(*bytes.Buffer).WriteString(b.String())
					return Nil
				}),
				Str("read"), Native2("read", func(env *Env, a, n Value) Value {
					rd := a.MustMap("").GetString("_buf").Interface().(*bytes.Buffer)
					if n := n.IntDefault(0); n > 0 {
						a := make([]byte, n)
						n, err := rd.Read(a)
						if err != nil && err != io.EOF {
							panic(err)
						}
						return Bytes(a[:n])
					} else {
						return Bytes(rd.Bytes())
					}
				}),
			)
			a := Map()
			a.Table().Parent = p.Table()
			return a
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
		}, "randomseed(int)"),
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
			"$f() => [0, 1)",
			"$f(n) => [0, n)",
			"$f(a, b) => [a, b]"),
		Str("sqrt"), Native1("sqrt", func(env *Env, v Value) Value { return Float(math.Sqrt(v.MustFloat(""))) }),
		Str("floor"), Native1("floor", func(env *Env, v Value) Value { return Float(math.Floor(v.MustFloat(""))) }),
		Str("ceil"), Native1("ceil", func(env *Env, v Value) Value { return Float(math.Ceil(v.MustFloat(""))) }),
		Str("min"), Native("min", func(env *Env) { mathMinMax(env, "#%d arg", false) }, "max(a, b, ...) => largest_number"),
		Str("max"), Native("max", func(env *Env) { mathMinMax(env, "#%d arg", true) }, "min(a, b, ...) => smallest_number"),
		Str("pow"), Native2("pow", func(env *Env, a, b Value) Value {
			af, ai, aIsInt := a.MustNum("base").Num()
			bf, bi, bIsInt := b.MustNum("power").Num()
			if aIsInt && bIsInt {
				return Int(ipow(ai, bi))
			}
			return Float(math.Pow(af, bf))
		}, "pow(a, b) => a to the power of b"),
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
		Str("shell"), Native1("shell", func(env *Env, cmd Value) Value {
			v, err := exec.Command("sh", "-c", cmd.MustStr("")).Output()
			if err != nil {
				panic(err)
			}
			return Bytes(v)
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
