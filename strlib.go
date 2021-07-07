package script

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

var StringMethods = ArrayMap(
	String("from"), Native1("from", func(env *Env, src Value) Value {
		return String(fmt.Sprint(src.Interface()))
	}, ""),
	String("iequal"), Native2("iequal", func(env *Env, src, a Value) Value {
		s := src.MustString("index", 0)
		return Bool(strings.EqualFold(s, a.MustString("iequal", 0)))
	}, ""),
	String("split"), Native3("split", func(env *Env, src, delim, n Value) Value {
		s := src.MustString("split", 0)
		d := delim.MustString("split delimeter", 0)
		r := []Value{}
		if n := n.IntDefault(0); n == 0 {
			for _, p := range strings.Split(s, d) {
				r = append(r, String(p))
			}
		} else {
			for _, p := range strings.SplitN(s, d, int(n)) {
				r = append(r, String(p))
			}
		}
		return Array(r...)
	}, "split(text, delim) => {part1, part2, ...}", "split(text, delim, n) => {part1, ..., partN}"),
	String("replace"), Native("replace", func(env *Env) {
		src := env.Get(0).MustString("replace", 0)
		from := env.Get(1).MustString("replace from text", 0)
		to := env.Get(2).MustString("replace to text", 0)
		n := env.Get(3).IntDefault(-1)
		env.A = String(strings.Replace(src, from, to, int(n)))
	}, ""),
	String("find"), Native2("index", func(env *Env, src, substr Value) Value {
		s := src.MustString("index", 0)
		return Int(int64(strings.Index(s, substr.MustString("index", 0))))
	}, ""),
	String("findany"), Native2("index_any", func(env *Env, src, substr Value) Value {
		s := src.MustString("index", 0)
		return Int(int64(strings.IndexAny(s, substr.MustString("index_any", 0))))
	}, ""),
	String("rfind"), Native2("last_index", func(env *Env, src, substr Value) Value {
		s := src.MustString("last_index", 0)
		return Int(int64(strings.LastIndex(s, substr.MustString("last_index", 0))))
	}, ""),
	String("rfindany"), Native2("last_index_any", func(env *Env, src, substr Value) Value {
		s := src.MustString("last_index", 0)
		return Int(int64(strings.LastIndexAny(s, substr.MustString("last_index_any", 0))))
	}, ""),
	String("sub"), Native3("sub", func(env *Env, src, start, end Value) Value {
		s := src.MustString("sub", 0)
		st := start.IntDefault(0)
		en := end.IntDefault(int64(len(s)))
		for st < 0 && len(s) > 0 {
			st += int64(len(s))
		}
		for en < 0 && len(s) > 0 {
			en += int64(len(s))
		}
		return String(s[st:en])
	}, ""),
	String("trim"), Native2("trim", func(env *Env, src, cutset Value) Value {
		if cutset == Nil {
			return String(strings.TrimSpace(src.MustString("trim_space", 0)))
		}
		c := cutset.MustString("trim cutset", 0)
		return String(strings.Trim(src.MustString("trim", 0), c))
	}, ""),
	String("ltrim"), Native2("trim_left", func(env *Env, src, cutset Value) Value {
		c := cutset.MustString("trim_left cutset", 0)
		return String(strings.TrimLeft(src.MustString("trim_left", 0), c))
	}, ""),
	String("rtrim"), Native2("trim_right", func(env *Env, src, cutset Value) Value {
		c := cutset.MustString("trim_right cutset", 0)
		return String(strings.TrimRight(src.MustString("trim_right", 0), c))
	}, ""),
	String("ptrim"), Native2("trim_prefix", func(env *Env, src, cutset Value) Value {
		c := cutset.MustString("trim_prefix", 0)
		return String(strings.TrimPrefix(src.MustString("trim_prefix", 0), c))
	}, ""),
	String("strim"), Native2("trim_suffix", func(env *Env, src, cutset Value) Value {
		c := cutset.MustString("trim_suffix", 0)
		return String(strings.TrimSuffix(src.MustString("trim_suffix", 0), c))
	}, ""),
	String("decode_utf8"), Native("decode_utf8", func(env *Env) {
		r, sz := utf8.DecodeRuneInString(env.Get(0).MustString("decode_utf8()", 0))
		env.A = Array(Int(int64(r)), Int(int64(sz)))
	}, "$f(string) => { char_unicode, width_in_bytes }"),
	String("startswith"), Native2("startswith", func(env *Env, t, p Value) Value {
		return Bool(strings.HasPrefix(t.MustString("startswith()", 0), p.MustString("startswith() prefix", 0)))
	}, "startswith(text, prefix) => bool"),
	String("endswith"), Native2("endswith", func(env *Env, t, s Value) Value {
		return Bool(strings.HasSuffix(t.MustString("endswith()", 0), s.MustString("endswith() suffix", 0)))
	}, "endswith(text, suffix) => bool"),
	String("upper"), Native1("upper", func(env *Env, t Value) Value {
		return String(strings.ToUpper(t.MustString("upper()", 0)))
	}, "$f(text) => TEXT"),
	String("lower"), Native1("lower", func(env *Env, t Value) Value {
		return String(strings.ToLower(t.MustString("lower()", 0)))
	}, "$f(TEXT) => text"),
	String("chars"), Native2("chars", func(env *Env, s, n Value) Value {
		var r []Value
		max := n.IntDefault(0)
		for s := s.MustString("chars", 0); len(s) > 0; {
			_, sz := utf8.DecodeRuneInString(s)
			if sz == 0 {
				break
			}
			r = append(r, String(s[:sz]))
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
	String("format"), Native("format", func(env *Env) {
		f := env.Get(0).MustString("format() pattern", 0)
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
			if f[idx+1] == '%' {
				p.WriteString("%")
				f = f[idx+2:]
				continue
			}
			tmp.Reset()
			tmp.WriteByte('%')
			expecting := VNil
			for f = f[idx+1:]; len(f) > 0 && expecting == VNil; {
				switch f[0] {
				case 'b', 'd', 'o', 'O', 'x', 'X', 'c', 'e', 'E', 'f', 'F', 'g', 'G':
					expecting = VNumber
				case 's', 'q', 'U':
					expecting = VString
				case 'v':
					expecting = VInterface
				case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', '-', '+', '#', ' ':
				default:
					panicf("format(): unexpected verb: '%c'", f[0])
				}
				tmp.WriteByte(f[0])
				f = f[1:]
			}
			switch expecting {
			case VString:
				p.WriteString(fmt.Sprintf(tmp.String(), pop().String()))
			case VNumber:
				f, i, isInt := pop().Num()
				if isInt {
					p.Write([]byte(fmt.Sprintf(tmp.String(), i)))
				} else {
					p.Write([]byte(fmt.Sprintf(tmp.String(), f)))
				}
			case VInterface:
				p.Write([]byte(fmt.Sprint(pop())))
			}
		}
		env.A = String(p.String())
	}, "format(pattern, a1, a2, ...)"),
	String("buffer"), Native1("buffer", func(env *Env, v Value) Value {
		b := &bytes.Buffer{}
		if v != Nil {
			b.WriteString(v.String())
		}
		p := ArrayMap(
			String("_buf"), Interface(b),
			String("value"), Native1("value", func(env *Env, a Value) Value {
				return Bytes(a.MustArray("", 0).GetString("_buf").Interface().(*bytes.Buffer).Bytes())
			}),
			String("write"), Native2("write", func(env *Env, a, b Value) Value {
				a.MustArray("", 0).GetString("_buf").Interface().(*bytes.Buffer).WriteString(b.String())
				return Nil
			}),
			String("read"), Native2("read", func(env *Env, a, n Value) Value {
				rd := a.MustArray("", 0).GetString("_buf").Interface().(*bytes.Buffer)
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
		a := ArrayMap()
		a.Array().Parent = p.Array()
		return a
	}),
)
