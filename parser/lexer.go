// Modified upon: yuin/gopher-lua
package parser

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
	"unsafe"
)

const EOF = 0xffffffff

var numberChars = func() (x [256]bool) {
	for _, r := range "0123456789abcdefABCDEF.xX_" {
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
		return fmt.Sprintf("%q at line %d: %s\n", e.Token, pos.Line, e.Message)
	}
}

func isIdent(ch uint32, pos int) bool {
	return ch == '_' || 'A' <= ch && ch <= 'Z' || 'a' <= ch && ch <= 'z' || '0' <= ch && ch <= '9' && pos > 0
}

type Scanner struct {
	Pos    Position
	buffer bytes.Buffer
	offset int64
	text   string
}

func NewScanner(text string, source string) *Scanner {
	return &Scanner{
		Pos:  Position{Source: source, Line: 1, Column: 0},
		text: text,
	}
}

func (sc *Scanner) Error(tok string, msg string) *Error {
	return &Error{sc.Pos, msg, tok}
}

func (sc *Scanner) TokenError(tok Token, msg string) *Error {
	return &Error{tok.Pos, msg, tok.Str}
}

func (sc *Scanner) Peek() uint32 {
	if sc.offset >= int64(len(sc.text)) {
		return EOF
	}
	return uint32(sc.text[sc.offset])
}

func (sc *Scanner) Next() uint32 {
	ch := sc.Peek()
	sc.offset++
	switch ch {
	case '\r':
		return sc.Next()
	case '\n':
		sc.Pos.Line++
		sc.Pos.Column = 0
	case EOF:
		sc.Pos.Line = EOF
		sc.Pos.Column = 0
	default:
		sc.Pos.Column++
	}
	return ch
}

func (sc *Scanner) skipComments() {
	for ch := sc.Next(); ; ch = sc.Next() {
		if ch == '\n' || ch < 0 || ch == EOF {
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

func (sc *Scanner) scanIdent(offsetOffset int64) string {
	start := sc.offset - 1 - offsetOffset
	for isIdent(sc.Peek(), 1) {
		sc.Next()
	}
	return sc.text[start:sc.offset]
}

func (sc *Scanner) scanNumber() string {
	start := sc.offset - 1
	for {
		ch := byte(sc.Peek())
		if !numberChars[ch] {
			if ch == '+' || ch == '-' {
				before := sc.text[sc.offset-1]
				if before == 'e' || before == 'E' {
					if x := sc.text[start+1]; x != 'x' && x != 'X' {
						// Not a hexdecimal string, so it is a float64 value (maybe)
						goto OK
					}
				}
			}
			dxx := sc.text[start:sc.offset]
			return dxx
		}
	OK:
		sc.Next()
	}
}

func (sc *Scanner) scanString(quote uint32) (string, error) {
	lastIsSlash := false
	buf := &sc.buffer
	buf.Reset()

	ch := sc.Next()
	for ch != quote || lastIsSlash {
		if ch == '\n' || ch < 0 {
			return "", sc.Error(buf.String(), "unterminated string")
		}
		lastIsSlash = ch == '\\'
		buf.WriteByte(byte(ch))
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
			return "", err
		}
		s = ss
		if c < utf8.RuneSelf || !multibyte {
			buf.WriteByte(byte(c))
		} else {
			n := utf8.EncodeRune(runeTmp[:], c)
			buf.Write(runeTmp[:n])
		}
	}
	return buf.String(), nil
}

func (sc *Scanner) scanBlockString() (string, error) {
	start := sc.offset
	for {
		ch := sc.Next()
		if ch == EOF {
			return "", sc.Error("", "unexpected end of string block")
		}
		if ch == ']' && sc.Peek() == ']' {
			sc.Next()
			break
		}
	}
	return sc.text[start : sc.offset-2], nil // -2: exclude ']' and ']' at end
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
	"in":       TIn,
	"goto":     TGoto,
}

func (sc *Scanner) Scan(lexer *Lexer) (Token, error) {
redo:
	var err error
	var tok Token

skipspaces:
	ch := sc.Next()
	if unicode.IsSpace(rune(ch)) {
		goto skipspaces
	}

	tok.Pos = sc.Pos

	switch {
	case isIdent(ch, 0):
		tok.Type = TIdent
		tok.Str = sc.scanIdent(0)

		if typ, ok := reservedWords[tok.Str]; ok {
			crlf := false
			for n := sc.Peek(); unicode.IsSpace(rune(n)) || n == 'e' || n == EOF; n = sc.Peek() {
				if n == '\n' || n == EOF {
					crlf = true
					break
				}
				if n == 'e' {
					crlf = strings.HasPrefix(sc.text[sc.offset:], "end")
					break
				}
				sc.Next()
			}

			// 'return' without an arg, but with a CrLf afterward will be considered
			// as 'return nil'. This rule implies the following syntax:
			//   1. return end
			//   2. return \n end
			if tok.Str == "return" && crlf {
				tok.Type = TReturnVoid
			} else {
				tok.Type = typ
			}
		}
	case ch >= '0' && ch <= '9':
		tok.Type = TNumber
		tok.Str = sc.scanNumber()
	default:
		switch ch {
		case EOF:
			tok.Type = EOF
		case '-':
			if sc.Peek() == '-' {
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
			}
			tok.Type = ch
			tok.Str = "-"
		case '"', '\'':
			tok.Type = TString
			tok.Str, err = sc.scanString(ch)
		case '[':
			if sc.Peek() == '[' {
				sc.Next()
				tok.Type = TString
				tok.Str, err = sc.scanBlockString()
			} else {
				tok.Type = ch
				tok.Str = "["
			}
		case '=', '!', '~', '<', '>':
			idx := strings.IndexByte("=!~<>", byte(ch))
			if p := sc.Peek(); p == '=' {
				tok.Type = [...]uint32{TEqeq, TNeq, TNeq, TLte, TGte}[idx]
				tok.Str = [...]string{"==", "!=", "~=", "<=", ">="}[idx]
				sc.Next()
			} else if p == ch && ch == '<' {
				tok.Type, tok.Str = TLsh, "<<"
				sc.Next()
			} else if p == ch && ch == '>' {
				sc.Next()
				if sc.Peek() == '>' {
					tok.Type, tok.Str = TURsh, ">>>"
					sc.Next()
				} else {
					tok.Type, tok.Str = TRsh, ">>"
				}
			} else {
				tok.Type = ch
				tok.Str = [...]string{"=", "!", "~", "<", ">"}[idx]
			}
		case '.':
			switch ch2 := sc.Peek(); {
			case ch2 >= '0' && ch2 <= '9':
				tok.Type = TNumber
				tok.Str = sc.scanNumber()
			default:
				tok.Type = '.'
				tok.Str = "."
			}
		case '(', ')', '{', '}', ']', ';', ',', '#', '^', '|', '@', '&':
			const pat = "(){}];,#^|@&"
			idx := strings.IndexByte(pat, byte(ch))
			tok.Type = ch
			tok.Str = pat[idx : idx+1]
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
			ii := strings.IndexByte("+*/%", byte(ch))
			if sc.Peek() == '/' && ch == '/' {
				tok.Type = TIDiv
				tok.Str = "//"
				sc.Next()
			} else {
				tok.Type = ch
				tok.Str = [...]string{"+", "*", "/", "%"}[ii]
			}
		default:
			err = sc.Error(string(rune(ch)), "invalid token")
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

func parse(reader string, name string) (chunk Node, lexer *Lexer, err error) {
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

func Parse(text, name string) (chunk Node, err error) {
	yyErrorVerbose = true
	yyDebug = 1
	chunk, _, err = parse(text, name)
	if !chunk.Valid() && err == nil {
		err = fmt.Errorf("invalid chunk")
	}
	return
}

// }}}
