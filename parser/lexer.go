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

	"github.com/coyove/nj/bas"
	"github.com/coyove/nj/internal"
)

const EOF = 0xffffffff

var numberChars [256]bool

func init() {
	for _, r := range "0123456789abcdefABCDEF.xX_" {
		numberChars[byte(r)] = true
	}
	update := func(idx int, name string) { yyToknames[idx-yyPrivate+1] = name }
	for name, idx := range reservedWords {
		update(int(idx), name)
	}
	update(TInv, "'-'")
	update(TLParen, "'('")
	update(TLBracket, "'['")
	update(TLsh, "'<<'")
	update(TRsh, "'>>'")
	update(TURsh, "'>>>'")
	update(TEqeq, "'=='")
	update(TNeq, "'!='")
	update(TLte, "'<='")
	update(TGte, "'>='")
	update(TIDiv, "'//'")
	update(TDotDotDot, "'...'")
	update(TLabel, "goto label")
	update(TReturnVoid, "return")
	update(TIdent, "identifier")
	update(TNumber, "number")
	update(TString, "string")
	update(TAddEq, "+=")
	update(TSubEq, "-=")
	update(TMulEq, "*=")
	update(TDivEq, "/=")
	update(TIDivEq, "//=")
	update(TModEq, "%=")
	update(TBitAndEq, "&=")
	update(TBitOrEq, "|=")
	update(TBitXorEq, "^=")
	update(TBitLshEq, "<<=")
	update(TBitRshEq, ">>=")
	update(TBitURshEq, ">>>=")
	yyToknames[0] = "<EOF>"
}

type Error struct {
	Pos     Position
	Message string
	Token   string
}

func (e *Error) Error() string {
	pos := e.Pos
	if pos.Line == EOF {
		return e.Message
	} else {
		msg := fmt.Sprintf("%q at %s:%d\t%s", e.Token, pos.Source, pos.Line, e.Message)
		// if e.Message == "syntax error: unexpected '('" {
		// 	msg += ", is there any space(' ') or newline('\\n') before it?"
		// }
		return (msg)
	}
}

func isIdent(ch uint32, pos int) bool {
	return ch == '_' || 'A' <= ch && ch <= 'Z' || 'a' <= ch && ch <= 'z' || '0' <= ch && ch <= '9' && pos > 0
}

type Scanner struct {
	Pos       Position
	buffer    bytes.Buffer
	offset    int64
	text      string
	lastToken Token
	constants bas.Map
	functions bas.Map
}

func NewScanner(text string, source string) *Scanner {
	o := bas.Map{}
	o.Set(bas.True, bas.Nil)
	o.Set(bas.False, bas.Nil)
	o.Set(bas.Int(0), bas.Nil)
	o.Set(bas.Int(1), bas.Nil)
	return &Scanner{
		constants: o,
		Pos:       Position{Source: source, Line: 1, Column: 0},
		text:      text,
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

func (sc *Scanner) isLastTokenSymbolClosed() bool {
	last := sc.lastToken
	// foo(...)
	// foo?(...)
	// foo()(...)
	// foo[bar](...)
	// {...}(...)
	// if(cond, true, false)
	// function() end(...)
	return last.Type == TEnd || last.Type == TIf || last.Type == TIdent ||
		last.Type == ')' || last.Type == ']' || last.Type == '}' || last.Type == '?'
}

func (sc *Scanner) isLastTokenSymbolOrNumberClosed() bool {
	return sc.isLastTokenSymbolClosed() || sc.lastToken.Type == TNumber
}

func (sc *Scanner) skipComments() {
	for ch := sc.Next(); ; ch = sc.Next() {
		if ch == '\n' || ch == EOF {
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
	needUnescape := false
	buf := &sc.buffer
	buf.Reset()

	ch := sc.Next()
	for ch != quote || lastIsSlash {
		if ch == '\n' || ch == EOF {
			return "", sc.Error(buf.String(), "unterminated string")
		}
		lastIsSlash = ch == '\\'
		needUnescape = needUnescape || lastIsSlash
		buf.WriteByte(byte(ch))
		ch = sc.Next()
	}
	if !needUnescape {
		return buf.String(), nil
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
			return "", sc.Error(buf.String(), "escape: "+err.Error())
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
	startPos := sc.Pos
	for {
		ch := sc.Next()
		if ch == EOF {
			sc.Pos = startPos
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
	"continue": TContinue,
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
	"is":       TIs,
}

var opAssignMap = map[uint32]uint32{
	'+': TAddEq, '-': TSubEq, '*': TMulEq, '/': TDivEq, TIDiv: TIDivEq, '%': TModEq,
	'&': TBitAndEq, '|': TBitOrEq, '^': TBitXorEq, TLsh: TBitLshEq, TRsh: TBitRshEq, TURsh: TBitURshEq,
}

func (sc *Scanner) Scan(lexer *Lexer) (Token, error) {
	var metSpaces bool

redo:
	var err error
	var tok Token

skipspaces:
	ch := sc.Next()
	if unicode.IsSpace(rune(ch)) {
		metSpaces = true
		goto skipspaces
	}

	tok.Pos = sc.Pos

	switch {
	case isIdent(ch, 0):
		tok.Type = TIdent
		tok.Str = sc.scanIdent(0)
		if len(tok.Str) > 1024 {
			err = sc.Error("", "identifier too long")
			goto finally
		}

		if typ, ok := reservedWords[tok.Str]; ok {
			tok.Type = typ
			if typ == TReturn {
				tail := strings.TrimLeft(sc.text[sc.offset:], " \t\r")
				if strings.HasPrefix(tail, "\n") { // return \n
					tok.Type = TReturnVoid
				} else if tail = strings.TrimLeft(tail, "\n"); tail == "" { // return <EOF>
					tok.Type = TReturnVoid
				} else if strings.HasPrefix(tail, ";") { // return ;
					tok.Type = TReturnVoid
				} else if strings.HasPrefix(tail, "--") { // return --comments
					tok.Type = TReturnVoid
				} else if strings.HasPrefix(tail, "end") { // return end
					tmp := tail[3:]
					r, _ := utf8.DecodeRuneInString(tmp)
					if tmp == "" || unicode.IsSpace(r) {
						// return <spaces> <keyword> (<spaces>|<EOF>)
						tok.Type = TReturnVoid
					}
				} else {
					// for k := range reservedWords {
					// 	if k == "function" || k == "if" {
					// 		// return function() end / return if(a, b, c)
					// 		continue
					// 	}
					// 	if strings.HasPrefix(tail, k) {
					// 		tmp := tail[len(k):]
					// 		r, _ := utf8.DecodeRuneInString(tmp)
					// 		if tmp == "" || unicode.IsSpace(r) {
					// 			// return <spaces> <keyword> (<spaces>|<EOF>)
					// 			tok.Type = TReturnVoid
					// 			break
					// 		}
					// 	}
					// }
				}
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
			if p := sc.Peek(); p == '-' {
				sc.Next()
				if strings.HasPrefix(sc.text[sc.offset:], "[[") {
					// --[[ block comment ]]
					if err = sc.skipBlockComments(); err != nil {
						goto finally
					}
				} else {
					sc.skipComments()
				}
				goto redo
			} else if p == '=' {
				sc.Next()
				tok.Type = TSubEq
				tok.Str = "-="
				goto finally
			} else if metSpaces || !sc.isLastTokenSymbolOrNumberClosed() {
				if p == '.' || (p >= '0' && p <= '9') {
					tok.Type = TNumber
					tok.Str = sc.scanNumber()
					goto finally
				} else if !unicode.IsSpace(rune(p)) {
					tok.Type = TInv
					tok.Str = "-"
					goto finally
				}
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
				if (sc.isLastTokenSymbolClosed() ||
					sc.lastToken.Type == TString) && // TString in prefix_expr
					!metSpaces {
					tok.Type = TLBracket
				} else {
					tok.Type = ch
				}
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
				tok = sc.opAssign(tok)
			} else if p == ch && ch == '>' {
				sc.Next()
				if sc.Peek() == '>' {
					tok.Type, tok.Str = TURsh, ">>>"
					sc.Next()
				} else {
					tok.Type, tok.Str = TRsh, ">>"
				}
				tok = sc.opAssign(tok)
			} else {
				tok.Type = ch
				tok.Str = [...]string{"=", "!", "~", "<", ">"}[idx]
			}
		case '.':
			switch ch2 := sc.Peek(); {
			case ch2 >= '0' && ch2 <= '9':
				tok.Type = TNumber
				tok.Str = sc.scanNumber()
			case ch2 == '.':
				sc.Next()
				if sc.Peek() == '.' {
					sc.Next()
					tok.Type = TDotDotDot
					tok.Str = "..."
				} else {
					err = sc.Error(string(rune(ch)), "unexpected dots")
					goto finally
				}
			default:
				tok.Type = '.'
				tok.Str = "."
			}
		case '(', ')', '{', '}', ']', ';', ',', '#', '^', '|', '&', '?':
			if ch == '(' && sc.isLastTokenSymbolClosed() && !metSpaces {
				tok.Type = TLParen
				tok.Str = "("
			} else {
				const pat = "(){}];,#^|&?"
				idx := strings.IndexByte(pat, byte(ch))
				tok.Type = ch
				tok.Str = pat[idx : idx+1]
				tok = sc.opAssign(tok)
			}
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
			tok = sc.opAssign(tok)
		default:
			err = sc.Error(string(rune(ch)), "invalid token")
			goto finally
		}
	}

finally:
	sc.lastToken = tok
	return tok, err
}

func (sc *Scanner) opAssign(tok Token) Token {
	if sc.Peek() == '=' {
		if m, ok := opAssignMap[tok.Type]; ok {
			tok.Type = m
		} else {
			return tok
		}
		sc.Next()
		tok.Str += "="
	}
	return tok
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

func parse(reader, name string) (chunk Node, lexer *Lexer, err error) {
	lexer = &Lexer{
		scanner: NewScanner(reader, name),
		Stmts:   nil,
		Token:   Token{Str: ""},
	}
	defer internal.CatchError(&err)
	yyParse(lexer)
	chunk = lexer.Stmts
	return
}

func Parse(text, name string) (chunk Node, err error) {
	yyErrorVerbose = true
	yyDebug = 1
	var lexer *Lexer
	chunk, lexer, err = parse(text, name)
	if chunk == nil && err == nil {
		err = fmt.Errorf("invalid chunk")
	}
	if chunk != nil {
		lc := &LoadConst{Table: lexer.scanner.constants, Funcs: lexer.scanner.functions}
		chunk.(*Prog).Stats = append([]Node{lc}, chunk.(*Prog).Stats...)
	}
	return
}

// }}}
