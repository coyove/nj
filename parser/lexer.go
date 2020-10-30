// Taken from: https://github.com/yuin/gopher-lua/blob/master/parse

/*
The MIT License (MIT)

Copyright (c) 2015 Yusuke Inuzuka

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR Sym PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
	"unsafe"
)

const (
	EOF         = 0xffffffff
	whitespace1 = 1<<'\t' | 1<<' '
	whitespace2 = 1<<'\t' | 1<<'\n' | 1<<'\r' | 1<<' '
)

var numberChars = func() (x [256]bool) {
	for _, r := range "0123456789abcdefABCDEF.xX" {
		x[byte(r)] = true
	}
	return
}()

type Error struct {
	Pos     Position
	Message string
	Token   string
}

func (e *Error) Error() string {
	pos := e.Pos
	if pos.Line == EOF {
		return fmt.Sprintf("%s\n", e.Message)
	} else {
		return fmt.Sprintf("%v at line %d: %s\n", e.Token, pos.Line, e.Message)
	}
}

func writeChar(buf *bytes.Buffer, c uint32) { buf.WriteByte(byte(c)) }

func isIdent(ch uint32, pos int) bool {
	return ch == '_' || 'A' <= ch && ch <= 'Z' || 'a' <= ch && ch <= 'z' || '0' <= ch && ch <= '9' && pos > 0
}

func isDigit(ch uint32) bool {
	return '0' <= ch && ch <= '9' || 'a' <= ch && ch <= 'f' || 'A' <= ch && ch <= 'F'
}

type Scanner struct {
	Pos    Position
	reader *bufio.Reader
	buffer bytes.Buffer
}

func NewScanner(reader io.Reader, source string) *Scanner {
	return &Scanner{
		Pos:    Position{Source: source, Line: 1, Column: 0},
		reader: bufio.NewReaderSize(reader, 4096),
	}
}

func (sc *Scanner) Error(tok string, msg string) *Error { return &Error{sc.Pos, msg, tok} }

func (sc *Scanner) TokenError(tok Token, msg string) *Error { return &Error{tok.Pos, msg, tok.Str} }

func (sc *Scanner) readNext() uint32 {
	ch, err := sc.reader.ReadByte()
	if err == io.EOF {
		return EOF
	}
	return uint32(ch)
}

func (sc *Scanner) Newline(ch uint32) {
	if ch < 0 {
		return
	}
	sc.Pos.Line++
	sc.Pos.Column = 0
	next := sc.Peek()
	if ch == '\n' && next == '\r' || ch == '\r' && next == '\n' {
		sc.reader.ReadByte()
	}
}

func (sc *Scanner) Next() uint32 {
	ch := sc.readNext()
	switch ch {
	case '\n', '\r':
		sc.Newline(ch)
		ch = uint32('\n')
	case EOF:
		sc.Pos.Line = EOF
		sc.Pos.Column = 0
	default:
		sc.Pos.Column++
	}
	return ch
}

func (sc *Scanner) Peek() uint32 {
	ch := sc.readNext()
	if ch != EOF {
		sc.reader.UnreadByte()
	}
	return ch
}

func (sc *Scanner) skipWhiteSpace(whitespace int64) uint32 {
	ch := sc.Next()
	for ; whitespace&(1<<uint(ch)) != 0; ch = sc.Next() {
	}
	return ch
}

func (sc *Scanner) skipComments() {
	for ch := sc.Next(); ; ch = sc.Next() {
		if ch == '\n' || ch == '\r' || ch < 0 || ch == EOF {
			return
		}
	}
}

func (sc *Scanner) skipBlockComments() error {
	for a := sc.Next(); a != EOF; a = sc.Next() {
		if a == ']' && sc.Peek() == ']' {
			sc.Next()
			return nil
		}
	}
	return sc.Error("", "unterminated block comments")
}

func (sc *Scanner) scanIdent(ch uint32, buf *bytes.Buffer) error {
	writeChar(buf, ch)
	for isIdent(sc.Peek(), 1) {
		writeChar(buf, sc.Next())
	}
	return nil
}

func (sc *Scanner) scanNumber(ch uint32, buf *bytes.Buffer) error {
	writeChar(buf, ch)
	for {
		ch := byte(sc.Peek())
		if !numberChars[ch] {
			if ch == '+' || ch == '-' {
				x := buf.Bytes()
				if e := x[len(x)-1]; e == 'e' || e == 'E' {
					if len(x) >= 2 && x[1] != 'x' && x[1] != 'X' {
						goto OK
					}
				}
			}
			return nil
		}
	OK:
		writeChar(buf, sc.Next())
	}
}

func (sc *Scanner) scanString(quote uint32, buf *bytes.Buffer) error {
	ch := sc.Next()
	lastIsSlash := false

	for ch != quote || lastIsSlash {
		if ch == '\n' || ch == '\r' || ch < 0 {
			return sc.Error(buf.String(), "unterminated string")
		}
		lastIsSlash = ch == '\\'
		writeChar(buf, ch)
		ch = sc.Next()
	}

	x := buf.Bytes()
	s := *(*string)(unsafe.Pointer(&x))
	buf.Reset()

	// Hack: escaped string's length is always greater or equal to its unescaped one
	// So we reset the buffer and write unescaped chars back directly because it will never
	// catch up the progress of 'UnquoteChar(escaped_string)'
	var runeTmp [utf8.UTFMax]byte
	for len(s) > 0 {
		c, multibyte, ss, err := strconv.UnquoteChar(s, byte(quote))
		if err != nil {
			return err
		}
		s = ss
		if c < utf8.RuneSelf || !multibyte {
			buf.WriteByte(byte(c))
		} else {
			n := utf8.EncodeRune(runeTmp[:], c)
			buf.Write(runeTmp[:n])
		}
	}
	return nil
}

func (sc *Scanner) scanBlockString(buf *bytes.Buffer) error {
	for {
		ch := sc.Next()
		if ch == EOF {
			return sc.Error(buf.String(), "unexpected end of string block")
		}
		if ch == ']' {
			if sc.Peek() == ']' {
				sc.Next()
				break
			}
		}
		writeChar(buf, ch)
	}
	return nil
}

var reservedWords = map[string]uint32{
	"and":      TAnd,
	"or":       TOr,
	"local":    TLocal,
	"break":    TBreak,
	"else":     TElse,
	"function": TFunc,
	"if":       TIf,
	"elseif":   TElseIf,
	"then":     TThen,
	"end":      TEnd,
	"not":      TNot,
	"return":   TReturn,
	"for":      TFor,
	"while":    TWhile,
	"repeat":   TRepeat,
	"until":    TUntil,
	"do":       TDo,
	"yield":    TYield,
	"goto":     TGoto,
}

func (sc *Scanner) Scan(lexer *Lexer) (Token, error) {
redo:
	var err error
	tok := Token{}

	ch := sc.skipWhiteSpace(whitespace1)
	if ch == '\n' || ch == '\r' {
		ch = sc.skipWhiteSpace(whitespace2)
	}

	buf := &sc.buffer
	buf.Reset()
	tok.Pos = sc.Pos

	switch {
	case isIdent(ch, 0):
		tok.Type = TIdent
		err = sc.scanIdent(ch, buf)
		tok.Str = buf.String()

		if err != nil {
			goto finally
		}
		if typ, ok := reservedWords[tok.Str]; ok {
			crlf := false
			for n := sc.Peek(); (unicode.IsSpace(rune(n)) || n == 'e') && n != EOF; n = sc.Peek() {
				if n == '\n' {
					crlf = true
					break
				}
				if n == 'e' {
					buf, _ := sc.reader.Peek(3)
					crlf = bytes.Equal(buf, []byte("end"))
					break
				}
				sc.Next()
			}

			// return/yield without an arg, but with a CrLf afterward will be considered
			// as return nil/yield nil
			if tok.Str == "return" && crlf {
				tok.Type = TReturnVoid
			} else if tok.Str == "yield" && crlf {
				tok.Type = TYieldVoid
			} else {
				tok.Type = typ
			}
		}
	case ch >= '0' && ch <= '9':
		tok.Type = TNumber
		err = sc.scanNumber(ch, buf)
		tok.Str = buf.String()
	default:
		switch ch {
		case EOF:
			tok.Type = EOF
		case '-':
			switch sc.Peek() {
			case '-':
				sc.Next()
				if sc.Peek() == '[' {
					sc.Next()
					if sc.Peek() == '[' { // --[[ block comment ]]
						sc.Next()
						if err = sc.skipBlockComments(); err != nil {
							goto finally
						}
					}
				}
				sc.skipComments()
				goto redo
			case '=':
				tok.Type = TSubEq
				tok.Str = "-="
				sc.Next()
			default:
				tok.Type = ch
				tok.Str = string(rune(ch))
			}
		case '"', '\'':
			err = sc.scanString(ch, buf)
			tok.Type = TString
			tok.Str = buf.String()
		case '[':
			if sc.Peek() == '[' {
				sc.Next()
				tok.Type = TString
				err = sc.scanBlockString(buf)
				tok.Str = buf.String()
			} else {
				tok.Type = ch
				tok.Str = string(rune(ch))
			}
		case '=':
			if p := sc.Peek(); p == '=' {
				tok.Type = TEqeq
				tok.Str = "=="
				sc.Next()
			} else {
				tok.Type = ch
				tok.Str = string(rune(ch))
			}
		case '!', '~':
			if sc.Peek() == '=' {
				tok.Type = TNeq
				tok.Str = "~="
				sc.Next()
			} else {
				tok.Type = ch
				tok.Str = string(rune(ch))
			}
		case '<':
			if sc.Peek() == '=' {
				tok.Type = TLte
				tok.Str = "<="
				sc.Next()
			} else {
				tok.Type = ch
				tok.Str = string(rune(ch))
			}
		case '>':
			if sc.Peek() == '=' {
				tok.Type = TGte
				tok.Str = ">="
				sc.Next()
			} else {
				tok.Type = ch
				tok.Str = string(rune(ch))
			}
		case '.':
			switch ch2 := sc.Peek(); {
			case ch2 >= '0' && ch2 <= '9':
				tok.Type = TNumber
				err = sc.scanNumber(ch, buf)
				tok.Str = buf.String()
			case ch2 == '.':
				sc.Next()
				if sc.Peek() != '.' {
					tok.Type = TDotDot
					tok.Str = ".."
				} else {
					buf.WriteString("..")
					sc.Next()
					sc.scanIdent('.', buf)
					tok.Type = TIdent
					tok.Str = buf.String()
				}
			default:
				tok.Type = '.'
				tok.Str = buf.String()
			}
		case '(', ')', '{', '}', ']', ';', ',', '#', '^':
			tok.Type = ch
			tok.Str = string(rune(ch))
		case ':':
			if sc.Peek() == ':' {
				tok.Type = TLabel
				tok.Str = "::"
				sc.Next()
			} else {
				tok.Type = ch
				tok.Str = ":"
			}
		case '+', '*', '/', '%':
			switch sc.Peek() {
			case '=':
				tok.Type = [...]uint32{TAddEq, TMulEq, TDivEq, TModEq}[strings.Index("+*/%", string(rune(ch)))]
				tok.Str = string(rune(ch)) + "="
				sc.Next()
			default:
				tok.Type = ch
				tok.Str = string(rune(ch))
			}
		default:
			writeChar(buf, ch)
			err = sc.Error(buf.String(), "invalid token")
			goto finally
		}
	}

finally:
	return tok, err
}

// yacc interface {{{

type Lexer struct {
	scanner *Scanner
	Stmts   Node
	Token   Token
}

func (lx *Lexer) Lex(lval *yySymType) int {
	tok, err := lx.scanner.Scan(lx)
	if err != nil {
		panic(err)
	}
	if tok.Type < 0 {
		return 0
	}
	lval.token = tok
	lx.Token = tok
	t := int32(tok.Type)
	if tok.Type == EOF {
		t = -1
	}
	return int(t)
}

func (lx *Lexer) Error(message string) {
	panic(lx.scanner.Error(lx.Token.Str, message))
}

func (lx *Lexer) TokenError(tok Token, message string) {
	panic(lx.scanner.TokenError(tok, message))
}

func parse(reader io.Reader, name string) (chunk Node, lexer *Lexer, err error) {
	lexer = &Lexer{
		scanner: NewScanner(reader, name),
		Stmts:   Node{},
		Token:   Token{Str: ""},
	}
	defer CatchError(&err)
	yyParse(lexer)
	chunk = lexer.Stmts
	return
}

func Parse(reader io.Reader, name string) (chunk Node, err error) {
	chunk, _, err = parse(reader, name)
	if !chunk.Valid() && err == nil {
		err = fmt.Errorf("invalid chunk")
	}
	return
}

// }}}
