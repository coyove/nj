package script

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

var StringMethods = Map(
	Str("from"), Native1("from", func(env *Env, src Value) Value {
		return Str(fmt.Sprint(src.Any()))
	}, ""),
	Str("iequal"), Native2("iequal", func(env *Env, src, a Value) Value {
		s := src.MustStr("index", 0)
		return Bool(strings.EqualFold(s, a.MustStr("iequal", 0)))
	}, ""),
	Str("split"), Native3("split", func(env *Env, src, delim, n Value) Value {
		s := src.MustStr("split", 0)
		d := delim.MustStr("split delimeter", 0)
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
		src := env.Get(0).MustStr("replace", 0)
		from := env.Get(1).MustStr("replace from text", 0)
		to := env.Get(2).MustStr("replace to text", 0)
		n := env.Get(3).IntDefault(-1)
		env.A = Str(strings.Replace(src, from, to, int(n)))
	}, ""),
	Str("find"), Native2("index", func(env *Env, src, substr Value) Value {
		s := src.MustStr("index", 0)
		return Int(int64(strings.Index(s, substr.MustStr("index", 0))))
	}, ""),
	Str("findany"), Native2("index_any", func(env *Env, src, substr Value) Value {
		s := src.MustStr("index", 0)
		return Int(int64(strings.IndexAny(s, substr.MustStr("index_any", 0))))
	}, ""),
	Str("rfind"), Native2("last_index", func(env *Env, src, substr Value) Value {
		s := src.MustStr("last_index", 0)
		return Int(int64(strings.LastIndex(s, substr.MustStr("last_index", 0))))
	}, ""),
	Str("rfindany"), Native2("last_index_any", func(env *Env, src, substr Value) Value {
		s := src.MustStr("last_index", 0)
		return Int(int64(strings.LastIndexAny(s, substr.MustStr("last_index_any", 0))))
	}, ""),
	Str("sub"), Native3("sub", func(env *Env, src, start, end Value) Value {
		s := src.MustStr("sub", 0)
		st := start.IntDefault(0)
		en := end.IntDefault(int64(len(s)))
		for st < 0 && len(s) > 0 {
			st += int64(len(s))
		}
		for en < 0 && len(s) > 0 {
			en += int64(len(s))
		}
		return Str(s[st:en])
	}, ""),
	Str("trim"), Native2("trim", func(env *Env, src, cutset Value) Value {
		if cutset == Nil {
			return Str(strings.TrimSpace(src.MustStr("trim_space", 0)))
		}
		c := cutset.MustStr("trim cutset", 0)
		return Str(strings.Trim(src.MustStr("trim", 0), c))
	}, ""),
	Str("ltrim"), Native2("trim_left", func(env *Env, src, cutset Value) Value {
		c := cutset.MustStr("trim_left cutset", 0)
		return Str(strings.TrimLeft(src.MustStr("trim_left", 0), c))
	}, ""),
	Str("rtrim"), Native2("trim_right", func(env *Env, src, cutset Value) Value {
		c := cutset.MustStr("trim_right cutset", 0)
		return Str(strings.TrimRight(src.MustStr("trim_right", 0), c))
	}, ""),
	Str("ptrim"), Native2("trim_prefix", func(env *Env, src, cutset Value) Value {
		c := cutset.MustStr("trim_prefix", 0)
		return Str(strings.TrimPrefix(src.MustStr("trim_prefix", 0), c))
	}, ""),
	Str("strim"), Native2("trim_suffix", func(env *Env, src, cutset Value) Value {
		c := cutset.MustStr("trim_suffix", 0)
		return Str(strings.TrimSuffix(src.MustStr("trim_suffix", 0), c))
	}, ""),
	Str("decode_utf8"), Native("decode_utf8", func(env *Env) {
		r, sz := utf8.DecodeRuneInString(env.Get(0).MustStr("decode_utf8()", 0))
		env.A = Array(Int(int64(r)), Int(int64(sz)))
	}, "$f(string) => { char_unicode, width_in_bytes }"),
	Str("startswith"), Native2("startswith", func(env *Env, t, p Value) Value {
		return Bool(strings.HasPrefix(t.MustStr("startswith()", 0), p.MustStr("startswith() prefix", 0)))
	}, "startswith(text, prefix) => bool"),
	Str("endswith"), Native2("endswith", func(env *Env, t, s Value) Value {
		return Bool(strings.HasSuffix(t.MustStr("endswith()", 0), s.MustStr("endswith() suffix", 0)))
	}, "endswith(text, suffix) => bool"),
	Str("upper"), Native1("upper", func(env *Env, t Value) Value {
		return Str(strings.ToUpper(t.MustStr("upper()", 0)))
	}, "$f(text) => TEXT"),
	Str("lower"), Native1("lower", func(env *Env, t Value) Value {
		return Str(strings.ToLower(t.MustStr("lower()", 0)))
	}, "$f(TEXT) => text"),
	Str("chars"), Native2("chars", func(env *Env, s, n Value) Value {
		var r []Value
		max := n.IntDefault(0)
		for s := s.MustStr("chars", 0); len(s) > 0; {
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
	}, "chars(string) => { char1, char2, ... }", "chars(string, max) => { char1, char2, ..., char_max }",
		"\tbreak a string into (at most 'max') chars, e.g.:",
		"\tchars('a中c') => { 'a', '中', 'c' }",
		"\tchars('a中c', 1) => { 'a' }",
	),
	Str("format"), Native("format", func(env *Env) {
		f := env.Get(0).MustStr("format() pattern", 0)
		p, tmp := bytes.Buffer{}, bytes.Buffer{}
		popi := 0
		pop := func() Value { popi++; return env.Get(popi) }
		for len(f) > 0 {
			idx := strings.Index(f, "%")
			if idx == -1 {
				p.WriteString(f)
				break
			} else if idx == len(f)-1 {
				panicf("format(): unexpected '%%' at end")
			}
			p.WriteString(f[:idx])
			if f[idx+1] == '%' {
				p.WriteString("%")
				f = f[idx+2:]
				continue
			}
			tmp.Reset()
			tmp.WriteByte('%')
			expecting := NIL
			for f = f[idx+1:]; len(f) > 0 && expecting == NIL; {
				switch f[0] {
				case 'b', 'd', 'o', 'O', 'c', 'e', 'E', 'f', 'F', 'g', 'G':
					expecting = NUM
				case 's', 'q', 'U':
					expecting = STR
				case 'x', 'X':
					expecting = STR + NUM
				case 'v':
					expecting = ANY
				case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', '-', '+', '#', ' ':
				default:
					panicf("format(): unexpected verb: '%c'", f[0])
				}
				tmp.WriteByte(f[0])
				f = f[1:]
			}
			switch expecting {
			case STR:
				p.WriteString(fmt.Sprintf(tmp.String(), pop().String()))
			case NUM:
				f, i, isInt := pop().Num()
				if isInt {
					fmt.Fprintf(&p, tmp.String(), i)
				} else {
					fmt.Fprintf(&p, tmp.String(), f)
				}
			case NUM + STR:
				v := pop()
				if v.Type() == STR {
					fmt.Fprintf(&p, tmp.String(), v.Str())
				} else {
					f, i, isInt := pop().Num()
					if isInt {
						fmt.Fprintf(&p, tmp.String(), i)
					} else {
						fmt.Fprintf(&p, tmp.String(), f)
					}
				}
			case ANY:
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
			Str("_buf"), Any(b),
			Str("value"), Native1("value", func(env *Env, a Value) Value {
				return Bytes(a.MustMap("", 0).GetString("_buf").Any().(*bytes.Buffer).Bytes())
			}),
			Str("write"), Native2("write", func(env *Env, a, b Value) Value {
				a.MustMap("", 0).GetString("_buf").Any().(*bytes.Buffer).WriteString(b.String())
				return Nil
			}),
			Str("read"), Native2("read", func(env *Env, a, n Value) Value {
				rd := a.MustMap("", 0).GetString("_buf").Any().(*bytes.Buffer)
				if n := n.IntDefault(0); n > 0 {
					a := make([]byte, n)
					n, err := rd.Read(a)
					if err != nil && err != io.EOF {
						panicf("read(): %v", err)
					}
					return Bytes(a[:n])
				} else {
					return Bytes(rd.Bytes())
				}
			}),
		)
		a := Map()
		a.Map().Parent = p.Map()
		return a
	}),
)
