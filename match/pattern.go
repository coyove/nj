package match

import (
	"fmt"
	"strings"
	"unicode"
	"unsafe"
)

const CAP_UNFINISHED = (-1)
const CAP_POSITION = (-2)
const LUA_MAXCAPTURES = 10
const MAXCCALLS = 200

type MatchState struct {
	src        string
	pat        string
	src_init   uintptr /* init of source string */
	src_end    uintptr /* end ('\0') of source string */
	p_end      uintptr /* end ('\0') of pattern */
	matchdepth int     /* control for recursive depth (to avoid C stack overflow) */
	level      int     /* total number of captures (finished or unfinished) */
	capture    [10]struct {
		init uintptr
		len  int
	}
}

const L_ESC = '%'
const SPECIALS = "^$*+?.([%-"

func panicf(t string, a ...interface{}) { panic(fmt.Errorf(t, a...)) }

func (ms *MatchState) check_capture(l int) int {
	l -= '1'
	if l < 0 || l >= ms.level || ms.capture[l].len == CAP_UNFINISHED {
		panicf("invalid capture index %%%d", l+1)
	}
	return l
}

func (ms *MatchState) capture_to_close() int {
	level := ms.level
	for level--; level >= 0; level-- {
		if ms.capture[level].len == CAP_UNFINISHED {
			return level
		}
	}
	panicf("invalid pattern capture")
	return 0
}

func deref(p uintptr) byte { return *(*byte)(unsafe.Pointer(p)) }

func derefadd(p *uintptr) byte { *p++; return *(*byte)(unsafe.Pointer(*p - 1)) }

func (ms *MatchState) classend(p uintptr) uintptr {
	c := derefadd(&p)
	switch c {
	case L_ESC:
		if p == ms.p_end {
			panicf("malformed pattern (ends with '%%')")
		}
		return p + 1
	case '[':
		if deref(p) == '^' {
			p++
		}
		for { /* look for a ']' */
			if p == ms.p_end {
				panicf("malformed pattern (missing ']')")
			}
			if derefadd(&p) == L_ESC && p < ms.p_end {
				p++ /* skip escapes (e.g. '%]') */
			}
			if deref(p) == ']' {
				break
			}
		}
		return p + 1
	default:
		return p
	}
}

func match_class(c byte, cl byte) bool {
	switch c := rune(c); cl {
	case 'a':
		return unicode.IsLetter(c)
	case 'A':
		return !unicode.IsLetter(c)
	case 'c':
		return unicode.IsControl(c)
	case 'C':
		return !unicode.IsControl(c)
	case 'd':
		return unicode.IsDigit(c)
	case 'D':
		return !unicode.IsDigit(c)
	case 'g':
		return unicode.IsGraphic(c)
	case 'G':
		return !unicode.IsGraphic(c)
	case 'l':
		return unicode.IsLower(c)
	case 'L':
		return !unicode.IsLower(c)
	case 'p':
		return unicode.IsPunct(c)
	case 'P':
		return !unicode.IsPunct(c)
	case 's':
		return unicode.IsSpace(c)
	case 'S':
		return !unicode.IsSpace(c)
	case 'u':
		return unicode.IsUpper(c)
	case 'U':
		return !unicode.IsUpper(c)
	case 'w':
		return unicode.IsLetter(c) || unicode.IsNumber(c)
	case 'W':
		return !(unicode.IsLetter(c) || unicode.IsNumber(c))
	case 'x':
		return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
	case 'X':
		return !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F'))
	case 'z':
		return (c == 0)
	case 'Z':
		return (c != 0)
	default:
		return rune(cl) == c
	}
}

func matchbracketclass(c byte, p uintptr, ec uintptr) bool {
	sig := true
	if deref(p+1) == '^' {
		sig = false
		p++
	}
	for p++; p < ec; p++ {
		if deref(p) == L_ESC {
			if p++; match_class(c, deref(p)) {
				return sig
			}
		} else if (deref(p+1) == '-') && (p+2 < ec) {
			p += 2
			if deref(p-2) <= c && c <= deref(p) {
				return sig
			}
		} else if deref(p) == c {
			return sig
		}
	}
	return !sig
}

func (ms *MatchState) singlematch(s uintptr, p, ep uintptr) bool {
	if s >= ms.src_end {
		return false
	}

	switch c := deref(s); deref(p) {
	case '.':
		return true /* matches any char */
	case L_ESC:
		return match_class(c, deref(p+1))
	case '[':
		return matchbracketclass(c, p, ep-1)
	default:
		return deref(p) == c
	}
}

func (ms *MatchState) matchbalance(s uintptr, p uintptr) uintptr {
	if p >= ms.p_end-1 {
		panicf("malformed pattern (missing arguments to '%%b')")
	}
	if deref(s) != deref(p) {
		return 0
	} else {
		b := deref(p)
		e := deref(p + 1)
		cont := 1
		for s++; s < ms.src_end; s++ {
			if deref(s) == e {
				if cont--; cont == 0 {
					return s + 1
				}
			} else if deref(s) == b {
				cont++
			}
		}
	}
	return 0 /* string ends out of balance */
}

func (ms *MatchState) max_expand(s, p uintptr, ep uintptr) uintptr {
	i := 0 /* counts maximum expand for item */
	for ms.singlematch(s+uintptr(i), p, ep) {
		i++
	}
	/* keeps trying to match with the maximum repetitions */
	for i >= 0 {
		res := ms.match((s + uintptr(i)), ep+1)
		if res != 0 {
			return res
		}
		i-- /* else didn't match; reduce 1 repetition to try again */
	}
	return 0
}

func (ms *MatchState) min_expand(s, p, ep uintptr) uintptr {
	for {
		res := ms.match(s, ep+1)
		if res != 0 {
			return res
		} else if ms.singlematch(s, p, ep) {
			s++ /* try with one more repetition */
		} else {
			return 0
		}
	}
}

func (ms *MatchState) start_capture(s, p uintptr, what int) uintptr {
	var res uintptr
	level := ms.level
	if level >= LUA_MAXCAPTURES {
		panicf("too many captures")
	}
	ms.capture[level].init = s
	ms.capture[level].len = (what)
	ms.level = level + 1
	if res = ms.match(s, p); res == 0 { /* match failed? */
		ms.level-- /* undo capture */
	}
	return res
}

func (ms *MatchState) end_capture(s, p uintptr) uintptr {
	l := ms.capture_to_close()
	var res uintptr
	ms.capture[l].len = int(s - ms.capture[l].init) /* close capture */
	if res = ms.match(s, p); res == 0 {             /* match failed? */
		ms.capture[l].len = CAP_UNFINISHED /* undo capture */
	}
	return res
}

func memcmp(a, b uintptr, len int) bool {
	for i := 0; i < len; i++ {
		if deref(a+uintptr(i)) != deref(b+uintptr(i)) {
			return false
		}
	}
	return true
}

func (ms *MatchState) match_capture(s uintptr, l int) uintptr {
	l = ms.check_capture(l)
	len := ms.capture[l].len
	if int(ms.src_end-s) >= len && memcmp(ms.capture[l].init, s, len) {
		return s + uintptr(len)
	}
	return 0
}

func (ms *MatchState) match(s, p uintptr) uintptr {
	if ms.matchdepth--; ms.matchdepth == 0 {
		panicf("pattern too complex")
	}
init: /* using goto's to optimize tail recursion */
	if p != ms.p_end { /* end of pattern? */
		switch deref(p) {
		case '(': /* start capture */
			if deref(p+1) == ')' { /* position capture? */
				s = ms.start_capture(s, p+2, CAP_POSITION)
			} else {
				s = ms.start_capture(s, p+1, CAP_UNFINISHED)
			}
		case ')': /* end capture */
			s = ms.end_capture(s, p+1)
		case '$':
			if (p + 1) != ms.p_end { /* is the '$' the last char in pattern? */
				goto dflt /* no; go to default */
			}
			if s != ms.src_end { /* check end of string */
				s = 0
			}
		case L_ESC: /* escaped sequences not in the format class[*+?-]? */
			switch deref(p + 1) {
			case 'b': /* balanced string? */
				s = ms.matchbalance(s, p+2)
				if s != 0 {
					p += 4
					goto init /* return match(ms, s, p + 4); */
				} /* else fail (s == NULL) */
			case 'f': /* frontier? */
				var previous byte
				p += 2
				if deref(p) != '[' {
					panicf("missing '[' after '%%f' in pattern")
				}
				ep := ms.classend(p) /* points to what is next */
				if s == ms.src_init {
					previous = '\x00'
				} else {
					previous = deref(s - 1)
				}
				if !matchbracketclass(previous, p, ep-1) && matchbracketclass(deref(s), p, ep-1) {
					p = ep
					goto init /* return match(ms, s, ep); */
				}
				s = 0 /* match failed */
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9': /* capture results (%0-%9)? */
				s = ms.match_capture(s, int(deref(p+1)))
				if s != 0 {
					p += 2
					goto init /* return match(ms, s, p + 2) */
				}
			default:
				goto dflt
			}
		default:
			goto dflt
		}
		ms.matchdepth++
		return s
		//default:
	dflt: /* pattern class plus optional suffix */
		ep := ms.classend(p) /* points to optional suffix */
		/* does not match at least once? */
		if !ms.singlematch(s, p, ep) {
			if deref(ep) == '*' || deref(ep) == '?' || deref(ep) == '-' { /* accept empty? */
				p = ep + 1
				goto init /* return match(ms, s, ep + 1); */
			} else { /* '+' or no suffix */
				s = 0 /* fail */
			}
		} else { /* matched once */
			switch deref(ep) { /* handle optional suffix */
			case '?': /* optional */
				res := ms.match(s+1, ep+1)
				if res != 0 {
					s = res
				} else {
					p = ep + 1
					goto init /* else return match(ms, s, ep + 1); */
				}
			case '+': /* 1 or more repetitions */
				s += 1 /* 1 match already done */
				fallthrough
			case '*': /* 0 or more repetitions */
				s = ms.max_expand(s, p, ep)
			case '-': /* 0 or more repetitions (minimum) */
				s = ms.min_expand(s, p, ep)
			default: /* no suffix */
				s++
				p = ep
				goto init /* return match(ms, s + 1, ep); */
			}
		}
	}
	ms.matchdepth++
	return s
}

func ptr(a string) uintptr {
	return *(*uintptr)(unsafe.Pointer(&a))
}

func Find(s string, p string, init int, plain bool) (int, int, []string) {
	if plain {
		idx := strings.Index(s, p)
		if idx == -1 {
			return -1, -1, nil
		}
		return idx + 1, idx + len(p), nil
	}

	s1 := uintptr(int(ptr(s)) + init)
	ms := MatchState{
		src:        s,
		pat:        p,
		src_init:   ptr(s),
		src_end:    ptr(s) + uintptr(len(s)),
		p_end:      ptr(p) + uintptr(len(p)),
		matchdepth: MAXCCALLS,
	}
	for {
		res := ms.match(s1, ptr(p))
		if res != 0 {
			var caps []string
			for _, c := range ms.capture {
				if c.init > 0 && c.len >= 0 {
					caps = append(caps, s[c.init-ptr(s):int(c.init-ptr(s))+c.len])
				}
			}
			return int(s1-ptr(s)) + 1, int(res - ptr(s)), caps
		}
		if s1++; s1 >= ms.src_end {
			break
		}
	}
	return -1, -1, nil
}

func Match(s string, p string, init int) []string {
	start, end, m := Find(s, p, init, false)
	if len(m) == 0 {
		if start == -1 {
			return nil
		}
		return []string{s[start-1 : end]}
	}
	return m
}
