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
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
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
)

const EOF = 0xffffffff
const whitespace1 = 1<<'\t' | 1<<' '
const whitespace2 = 1<<'\t' | 1<<'\n' | 1<<'\r' | 1<<' '

type Error struct {
	Pos     Meta
	Message string
	Token   string
}

func (e *Error) Error() string {
	pos := e.Pos
	if pos.Line == EOF {
		return fmt.Sprintf("%v:eof: %s\n", pos.Source, e.Message)
	} else {
		return fmt.Sprintf("%v@%v:%d:%d: %s\n", e.Token, pos.Source, pos.Line, pos.Column, e.Message)
	}
}

func writeChar(buf *bytes.Buffer, c uint32) { buf.WriteByte(byte(c)) }

func isDecimal(ch uint32) bool { return '0' <= ch && ch <= '9' }

func isIdent(ch uint32, pos int) bool {
	return ch == '_' || 'A' <= ch && ch <= 'Z' || 'a' <= ch && ch <= 'z' || isDecimal(ch) && pos > 0
}

func isDigit(ch uint32) bool {
	return '0' <= ch && ch <= '9' || 'a' <= ch && ch <= 'f' || 'A' <= ch && ch <= 'F'
}

type Scanner struct {
	Pos    Meta
	reader *bufio.Reader
}

func NewScanner(reader io.Reader, source string) *Scanner {
	return &Scanner{
		Pos:    Meta{source, 1, 0, 0},
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

func (sc *Scanner) skipComments(ch uint32) error {
	for {
		if ch == '\n' || ch == '\r' || ch < 0 || ch == EOF {
			break
		}
		ch = sc.Next()
	}
	return nil
}

func (sc *Scanner) skipBlockComments() error {
	for {
		a := sc.Next()
		if a == EOF {
			return sc.Error("", "unterminated block comments")
		}
		b := sc.Peek()
		if a == '*' && b == '/' {
			sc.Next()
			return nil
		}
	}
}

func (sc *Scanner) scanIdent(ch uint32, buf *bytes.Buffer) error {
	writeChar(buf, ch)
	for isIdent(sc.Peek(), 1) {
		writeChar(buf, sc.Next())
	}
	return nil
}

func (sc *Scanner) scanDecimal(ch uint32, buf *bytes.Buffer) error {
	writeChar(buf, ch)
	for isDecimal(sc.Peek()) {
		writeChar(buf, sc.Next())
	}
	return nil
}

func (sc *Scanner) scanNumber(ch uint32, buf *bytes.Buffer) error {
	if ch == '0' { // octal
		switch sc.Peek() {
		case 'x', 'X', 'b', 'B', 'i', 'I':
			writeChar(buf, ch)
			writeChar(buf, sc.Next())
			hasvalue := false
			for isDigit(sc.Peek()) {
				writeChar(buf, sc.Next())
				hasvalue = true
			}
			if !hasvalue {
				return sc.Error(buf.String(), "illegal number")
			}
			return nil
		}
	}
	sc.scanDecimal(ch, buf)
	if sc.Peek() == '.' {
		sc.scanDecimal(sc.Next(), buf)
	}
	if ch = sc.Peek(); ch == 'e' || ch == 'E' {
		writeChar(buf, sc.Next())
		if ch = sc.Peek(); ch == '-' || ch == '+' {
			writeChar(buf, sc.Next())
		}
		sc.scanDecimal(sc.Next(), buf)
	}

	return nil
}

func (sc *Scanner) scanString(quote uint32, buf *bytes.Buffer) error {
	ch := sc.Next()
	for ch != quote {
		if ch == '\n' || ch == '\r' || ch < 0 {
			return sc.Error(buf.String(), "unterminated string")
		}
		if ch == '\\' {
			if err := sc.scanEscape(ch, buf); err != nil {
				return err
			}
		} else {
			writeChar(buf, ch)
		}
		ch = sc.Next()
	}
	return nil
}

func (sc *Scanner) scanEscape(ch uint32, buf *bytes.Buffer) error {
	ch = sc.Next()
	switch ch {
	case 'a':
		buf.WriteByte('\a')
	case 'b':
		buf.WriteByte('\b')
	case 'f':
		buf.WriteByte('\f')
	case 'n':
		buf.WriteByte('\n')
	case 'r':
		buf.WriteByte('\r')
	case 't':
		buf.WriteByte('\t')
	case 'v':
		buf.WriteByte('\v')
	case '\\':
		buf.WriteByte('\\')
	case '"':
		buf.WriteByte('"')
	case '\'':
		buf.WriteByte('\'')
	case '\n':
		buf.WriteByte('\n')
	case '\r':
		buf.WriteByte('\n')
		sc.Newline('\r')
	case 'u':
		ubuf := make([]byte, 4)
		for i := 0; i < 4; i++ {
			if isDigit(sc.Peek()) {
				ubuf[i] = byte(sc.Next())
			} else {
				return sc.Error(buf.String(), "invalid unicode escape sequence")
			}
		}
		val, _ := strconv.ParseInt(string(ubuf), 16, 32)
		buf.WriteRune(rune(val))
	case 'x':
		bbuf := make([]byte, 2)
		for i := 0; i < 2; i++ {
			if isDigit(sc.Peek()) {
				bbuf[i] = byte(sc.Next())
			} else {
				return sc.Error(buf.String(), "invalid hex escape sequence")
			}
		}
		val, _ := strconv.ParseInt(string(bbuf), 16, 32)
		buf.WriteByte(byte(val))
	default:
		if '0' <= ch && ch <= '9' {
			bytes := []byte{byte(ch)}
			for i := 0; i < 2 && isDecimal(sc.Peek()); i++ {
				bytes = append(bytes, byte(sc.Next()))
			}
			val, _ := strconv.ParseInt(string(bytes), 10, 32)
			writeChar(buf, uint32(val))
		} else {
			buf.WriteByte('\\')
			writeChar(buf, ch)
			return sc.Error(buf.String(), "invalid escape sequence")
		}
	}
	return nil
}

func (sc *Scanner) countSep(ch uint32) (int, uint32) {
	count := 0
	for ; ch == '='; count = count + 1 {
		ch = sc.Next()
	}
	return count, ch
}

func (sc *Scanner) scanBlockString(buf *bytes.Buffer) error {
	for {
		ch := sc.Next()
		if ch == EOF {
			return sc.Error(buf.String(), "unexpected end of string block")
		}
		if ch == '`' {
			break
		}
		writeChar(buf, ch)
	}
	return nil
}

var reservedWords = map[string]uint32{
	"and":       TAnd,
	"or":        TOr,
	"assert":    TAssert,
	"break":     TBreak,
	"case":      TCase,
	"continue":  TContinue,
	"else":      TElse,
	"var":       TVar,
	"fun":       TFunc,
	"if":        TIf,
	"len":       TLen,
	"not":       TNot,
	"return":    TReturn,
	"use":       TUse,
	"switch":    TSwitch,
	"addressof": TAddressof,
	"typeof":    TTypeof,
	"for":       TFor,
	"while":     TWhile,
	"foreach":   TForeach,
	"yield":     TYield,
}

func (sc *Scanner) Scan(lexer *Lexer) (Token, error) {
redo:
	var err error
	tok := Token{}

	ch := sc.skipWhiteSpace(whitespace1)
	if ch == '\n' || ch == '\r' {
		ch = sc.skipWhiteSpace(whitespace2)
	}

	var _buf bytes.Buffer
	buf := &_buf
	tok.Pos = sc.Pos

	switch {
	case ch == '$', isIdent(ch, 0):
		tok.Type = TIdent
		err = sc.scanIdent(ch, buf)
		tok.Str = buf.String()

		if err != nil {
			goto finally
		}
		if typ, ok := reservedWords[tok.Str]; ok {
			tok.Type = typ
		}
	case isDecimal(ch):
		tok.Type = TNumber
		err = sc.scanNumber(ch, buf)
		tok.Str = buf.String()
	default:
		switch ch {
		case EOF:
			tok.Type = EOF
		case '/':
			switch sc.Peek() {
			case '/':
				err = sc.skipComments(sc.Next())
				if err != nil {
					goto finally
				}
				goto redo
			case '*':
				sc.Next()
				err = sc.skipBlockComments()
				if err != nil {
					goto finally
				}
				goto redo
			case '=':
				tok.Type = TDivEq
				tok.Str = "/="
				sc.Next()
			default:
				tok.Type = ch
				tok.Str = string(ch)
			}
		case '"', '\'':
			tok.Type = TString
			err = sc.scanString(ch, buf)
			tok.Str = buf.String()
		case '`':
			tok.Type = TString
			err = sc.scanBlockString(buf)
			tok.Str = buf.String()
		case '[':
			tok.Type = ch
			tok.Str = string(ch)
		case '=', '+', '-':
			p := sc.Peek()
			if p == ch {
				switch ch {
				case '=':
					tok.Type = TEqeq
				case '+':
					tok.Type = TAddAdd
				case '-':
					tok.Type = TSubSub
				}
				tok.Str = string(ch) + string(ch)
				sc.Next()
			} else if p == '=' && (ch == '+' || ch == '-') {
				if ch == '+' {
					tok.Str = "+="
					tok.Type = TAddEq
				} else {
					tok.Str = "-="
					tok.Type = TSubEq
				}
				sc.Next()
			} else {
				tok.Type = ch
				tok.Str = string(ch)
			}
		case '!':
			if sc.Peek() == '=' {
				tok.Type = TNeq
				tok.Str = "!="
				sc.Next()
			} else {
				tok.Type = ch
				tok.Str = string(ch)
			}
		case '<':
			if sc.Peek() == '=' {
				tok.Type = TLte
				tok.Str = "<="
				sc.Next()
			} else if sc.Peek() == '<' {
				tok.Type = TLsh
				tok.Str = "<<"
				sc.Next()
				if sc.Peek() == '=' {
					tok.Type = TLshEq
					tok.Str = "<<="
					sc.Next()
				}
			} else {
				tok.Type = ch
				tok.Str = string(ch)
			}
		case '>':
			if sc.Peek() == '=' {
				tok.Type = TGte
				tok.Str = ">="
				sc.Next()
			} else if sc.Peek() == '>' {
				tok.Type = TRsh
				tok.Str = ">>"
				sc.Next()
				switch sc.Peek() {
				case '>':
					tok.Type = TURsh
					tok.Str = ">>>"
					sc.Next()
					if sc.Peek() == '=' {
						tok.Type = TURshEq
						tok.Str = ">>>="
						sc.Next()
					}
				case '=':
					tok.Type = TRshEq
					tok.Str = ">>="
					sc.Next()
				}
			} else {
				tok.Type = ch
				tok.Str = string(ch)
			}
		case '.':
			ch2 := sc.Peek()
			switch {
			case isDecimal(ch2):
				tok.Type = TNumber
				err = sc.scanNumber(ch, buf)
				tok.Str = buf.String()
			default:
				tok.Type = '.'
			}
			tok.Str = buf.String()
		case '(', ')', '{', '}', ']', ';', ':', ',', '~', '#':
			tok.Type = ch
			tok.Str = string(ch)
		case '*', '%', '&', '|', '^':
			switch sc.Peek() {
			case '=':
				tok.Type = [5]uint32{TMulEq, TModEq, TAndEq, TOrEq, TXorEq}[strings.Index("*%&|^", string(ch))]
				tok.Str = string(ch) + "="
				sc.Next()
			default:
				tok.Type = ch
				tok.Str = string(ch)
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
	loop    string
	cache   map[string]*Node
	Stmts   *Node
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

func parse(reader io.Reader, name string, cache map[string]*Node, loop string) (chunk *Node, lexer *Lexer, err error) {
	lexer = &Lexer{
		scanner: NewScanner(reader, name),
		loop:    loop,
		Stmts:   nil,
		Token:   Token{Str: ""},
		cache:   cache,
	}
	defer func() {
		if e := recover(); e != nil {
			err, _ = e.(error)
		}
	}()
	yyParse(lexer)
	chunk = lexer.Stmts
	return
}

func Parse(reader io.Reader, name string) (chunk *Node, err error) {
	chunk, _, err = parse(reader, name, make(map[string]*Node), name)
	return
}

// }}}
