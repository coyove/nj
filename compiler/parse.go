package compiler

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/coyove/eugine/base"
)

const (
	TK_lparen = iota
	TK_rparen
	TK_atomic
	TK_string
	TK_number
	TK_compound
	TK_addr
	TK_error
)

type token struct {
	ty          byte
	v           interface{}
	line, index int
	source      string
}

func (t *token) String() string {
	return fmt.Sprintf("(%+v,%d,%s:%d:%d)", t.v, t.ty, t.source, t.line, t.index)
}

type tokenReader struct {
	f           io.Reader
	last        rune
	line, index int
	source      string
}

func (f *tokenReader) rollbackRune(r rune) {
	f.last = r
}

func (f *tokenReader) nextRune() (r rune, err error) {
	if l := f.last; l != 0 {
		f.last = 0
		return l, nil
	}
	buf := make([]byte, 3)
	n, err := f.f.Read(buf[:1])
	if n != 1 {
		return 0, err
	}

	b := buf[0]
	if b < 128 {
		return rune(b), nil
	}
	if (b >> 5) == 6 {
		n, err := f.f.Read(buf[1:2])
		if n != 1 {
			return 0, err
		}
		r, _ := utf8.DecodeRune(buf[:2])
		return r, nil
	}
	if (b >> 5) == 7 {
		n, err := f.f.Read(buf[1:3])
		if n != 2 {
			return 0, err
		}
		r, _ := utf8.DecodeRune(buf)
		return r, nil
	}
	panic("ignore plane X")
}

func (f *tokenReader) nextHexByte() (s byte, err error) {
	buf := make([]byte, 2)
	n, err := f.f.Read(buf)
	if n != 2 {
		return 0, err
	}
	i, err := strconv.ParseInt(string(buf), 16, 64)
	return byte(i), err
}

func (f *tokenReader) nextToken() (tk *token) {
	defer func() {
		if tk != nil {
			tk.line, tk.index, tk.source = f.line, f.index, f.source
		}
	}()

	r, err := f.nextRune()
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return &token{ty: TK_error, v: err.Error()}
	}

	f.index++
	switch r {
	case '\n':
		f.line++
		f.index = 0
		fallthrough
	case '\r', ' ', '\t':
		return f.nextToken()
	case '[', '(', '{':
		return &token{ty: TK_lparen, v: r}
	case ']', ')', '}':
		return &token{ty: TK_rparen, v: r}
	case ':':
		return &token{ty: TK_atomic, v: r}
	case '#':
		for r, err := f.nextRune(); err == nil; r, err = f.nextRune() {
			if r == '\n' {
				f.line++
				f.index = 0
				return f.nextToken()
			}
			f.index++
		}
		return nil
	case '"':
		buf, escaped := &bytes.Buffer{}, false
		for r, err := f.nextRune(); err == nil; r, err = f.nextRune() {
			f.index++

			if r == '\\' && !escaped {
				escaped = true
				continue
			}

			if escaped {
				escaped = false
				switch r {
				case '\\':
					buf.WriteRune('\\')
				case 'n':
					buf.WriteRune('\n')
				case 't':
					buf.WriteRune('\t')
				case '"':
					buf.WriteRune('"')
				case 'x':
					b, err := f.nextHexByte()
					if err != nil {
						return &token{ty: TK_error, v: "unexpected hex escape"}
					}
					f.index += 2
					buf.WriteByte(b)
				case 'u':
					b1, err1 := f.nextHexByte()
					b2, err2 := f.nextHexByte()
					if err1 != nil || err2 != nil {
						return &token{ty: TK_error, v: "unexpected unicode escape"}
					}
					f.index += 4
					buf.WriteRune(rune(b1)*256 + rune(b2))
				default:
					return &token{ty: TK_error, v: "unexpected escape char: " + string(r)}
				}
			} else if r == '"' {
				return &token{ty: TK_string, v: buf.String()}
			} else {
				buf.WriteRune(r)
			}
		}
		return &token{ty: TK_error, v: "unexpected end of string"}
	case '+', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
		buf := &bytes.Buffer{}
		buf.WriteRune(r)
		for r, err := f.nextRune(); err == nil; r, err = f.nextRune() {
			if !strings.ContainsRune("0123456789.eE-+", r) {
				f.rollbackRune(r)
				break
			}
			buf.WriteRune(r)
			f.index++
		}
		num, err := strconv.ParseFloat(buf.String(), 64)
		if err != nil {
			return &token{ty: TK_atomic, v: buf.String()}
		}
		return &token{ty: TK_number, v: num}
	default:
		buf := &bytes.Buffer{}
		buf.WriteRune(r)
		for r, err := f.nextRune(); err == nil; r, err = f.nextRune() {
			if strings.ContainsRune("[](){}:#\"\r\n\t ", r) {
				f.rollbackRune(r)
				break
			}
			buf.WriteRune(r)
			f.index++
		}
		return &token{ty: TK_atomic, v: buf.String()}
	}
}

func (f *tokenReader) Close() {
	if f.f.(*os.File) != nil {
		f.f.(*os.File).Close()
	}
}

func newTokenReader(path string, content bool) (*tokenReader, error) {
	frr := &tokenReader{}
	if content {
		frr.f = strings.NewReader(path)
		frr.source = "vm"
	} else {
		src, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		frr.f = src
		frr.source = path
	}
	return frr, nil
}

func (f *tokenReader) parseNext(t *token) (*token, error) {
	if t == nil {
		return nil, nil
	}

	switch t.ty {
	case TK_rparen:
		return nil, fmt.Errorf("unexpected right bracket: " + t.String())
	case TK_lparen:
		tokens := make([]*token, 0, 8)
		for t = f.nextToken(); t != nil && t.ty != TK_rparen; t = f.nextToken() {
			if t.ty == TK_error {
				return nil, fmt.Errorf(t.v.(string))
			}
			ts, err := f.parseNext(t)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, ts)
		}
		if t == nil || t.ty != TK_rparen {
			return nil, fmt.Errorf("unexepcted end of code")
		}
		return &token{ty: TK_compound, v: tokens}, nil
	default:
		return t, nil
	}
}

func (f *tokenReader) parse() (code []byte, err error) {
	// tokens := make([]*token, 0, 16)
	buf := base.NewBytesBuffer()
	varLookup := base.NewCMap()

	var stackPtr int16
	var a *token

	for t := f.nextToken(); t != nil; t = f.nextToken() {
		a, err = f.parseNext(t)
		if err != nil {
			return nil, err
		}
		// tokens = append(tokens, ts)

		if a.ty != TK_compound {
			err = fmt.Errorf("every atom in the chain must be a compound: %+v", a)
			return
		}

		code, _, stackPtr, err = compileImpl(stackPtr, a.v.([]*token), varLookup)
		if err != nil {
			return
		}

		buf.Write(code)
	}

	buf.WriteByte(base.OP_EOB)
	return buf.Bytes(), nil
}

func equalI(l, r interface{}) bool {
	switch l.(type) {
	case string:
		if rf, ok := r.(string); ok {
			return l.(string) == rf
		}
	}
	return false
}

func LoadFile(path string) ([]byte, error) {
	tr, err := newTokenReader(path, false)
	if err != nil {
		return nil, err
	}

	defer tr.Close()
	return tr.parse()
}
