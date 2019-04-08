package rpn

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Type byte

const (
	T_Common Type = 1 << iota
	T_String
	T_Number
	T_L
	T_R
	T_Isolated
)

type Token struct {
	Name interface{}
	ty   Type
	col  uint16 // save some space
	line uint32
}

func (t *Token) String() string {
	ts := ""
	x := t.Name
	switch t.ty {
	case T_Common:
		ts = "----"
	case T_String:
		ts = "str"
		x = strconv.Quote(t.Name.(string))
	case T_Number:
		ts = "num"
	case T_L:
		ts = "["
	case T_R:
		ts = "]"
	case T_L | T_Isolated:
		ts = "i["
	default:
		panic("shouldn't happen")
	}

	// line and col are all zero-based, so line + 1, however col is always advanced by 1, so col + 1 - 1 -> col
	return fmt.Sprintf("|% 4s %4d:%3d | %v", ts, t.line+1, t.col, x)
}

func (t *Token) Line() int { return int(t.line + 1) }

func (t *Token) Column() int { return int(t.col) }

func (t *Token) Is(ty Type) bool {
	return t.ty&ty > 0
}

type Reader struct {
	r         *bufio.Reader
	dotNames  []byte
	line, col int
}

func NewReader(r io.Reader) *Reader {
	return &Reader{r: bufio.NewReader(r)}
}

func NewReaderFromString(str string) *Reader {
	return NewReader(strings.NewReader(str))
}

func (r *Reader) Token() (tok *Token, err error) {
	defer func() {
		if tok != nil {
			tok.line = uint32(r.line)
			tok.col = uint16(r.col)
		}
	}()

	if len(r.dotNames) > 0 {
		if r.dotNames[0] == '.' {
			idx := bytes.IndexByte(r.dotNames[1:], '.')
			if idx == -1 {
				idx = len(r.dotNames) - 1
			}
			name := string(r.dotNames[1 : idx+1])
			r.dotNames = r.dotNames[1:]
			return &Token{Name: name, ty: T_String}, nil
		}

		idx := bytes.IndexByte(r.dotNames, '.')
		if idx == -1 {
			r.dotNames = r.dotNames[:0]
		} else {
			r.dotNames = r.dotNames[idx:]
		}
		return &Token{Name: "<load>", ty: T_Common}, nil
	}

	metDim := false
AGAIN:
	lead, err := r.r.ReadByte()
	if err != nil {
		return nil, err
	}
	r.col++
	if lead == ' ' || lead == '\t' || lead == '\r' {
		metDim = true
		goto AGAIN
	}
	if lead == '\n' {
		r.col = 0
		r.line++
		metDim = true
		goto AGAIN
	}
	if lead == '[' {
		if metDim {
			return &Token{ty: T_L | T_Isolated}, nil
		}
		return &Token{ty: T_L}, nil
	}
	if lead == ']' {
		return &Token{ty: T_R}, nil
	}
	if isNum(lead, false) {
		return r.readNum(lead)
	}
	if lead == '"' {
		return r.readString()
	}
	if lead == '!' {
		return &Token{Name: "!", ty: T_Common}, nil
	}
	if !isSep(lead, false) {
		return r.readCommon(lead)
	}

	return nil, fmt.Errorf("unexpected char: %s", string(lead))
}

func (r *Reader) readNum(lead byte) (*Token, error) {
	var buf bytes.Buffer
	buf.WriteByte(lead)
	for {
		x, err := r.r.ReadByte()
		if err != nil {
			break
		}
		if !isNum(x, true) {
			r.r.UnreadByte()
			break
		}
		buf.WriteByte(x)
		r.col++
	}

	x := buf.String()
	if x == "-" || x == "+" || x == "--" || x == "++" {
		return &Token{Name: x, ty: T_Common}, nil
	}

	v, err := strconv.ParseFloat(x, 64)
	if err != nil {
		return nil, err
	}

	return &Token{Name: v, ty: T_Number}, nil
}

func (r *Reader) readCommon(lead byte) (*Token, error) {
	var buf bytes.Buffer
	buf.WriteByte(lead)
	dot := false
	for {
		x, err := r.r.ReadByte()
		if err != nil {
			break
		}
		if isSep(x, true) {
			r.r.UnreadByte()
			break
		}
		if x == '.' {
			dot = true
		}
		buf.WriteByte(x)
		r.col++
	}

	if !dot {
		return &Token{Name: buf.String(), ty: T_Common}, nil
	}

	x := buf.Bytes()
	idx := bytes.IndexByte(x, '.')
	r.dotNames = x[idx:]
	return &Token{Name: string(x[:idx]), ty: T_Common}, nil
}

func (r *Reader) readString() (*Token, error) {
	escape := false
	buf := bytes.Buffer{}
	for {
		x, n, err := r.r.ReadRune()
		if err != nil {
			return nil, err
		}

		r.col += n
		if x == '\n' {
			r.col = 0
			r.line++
		}

		if x == '\\' {
			if escape {
				buf.WriteRune(x)
				goto ESCAPE
			}
			escape = true
			continue
		}
		if x == 'n' && escape {
			buf.WriteRune('\n')
			goto ESCAPE
		}
		if x == 'r' && escape {
			buf.WriteRune('\r')
			goto ESCAPE
		}
		if x == 't' && escape {
			buf.WriteRune('\t')
			goto ESCAPE
		}
		if x == '"' {
			if escape {
				buf.WriteRune('"')
				goto ESCAPE
			}
			break
		}

		buf.WriteRune(x)
		continue

	ESCAPE:
		escape = false
	}

	return &Token{Name: buf.String(), ty: T_String}, nil

}

func isSep(r byte, sq bool) bool {
	x := r == ' ' || r == '\t' || r == '\r' || r == '\n' || r == '[' || r == ']' || r == '"' || r == '!'
	if sq {
		x = x || r == '\''
	}
	return x
}

func isAlpha(r byte, num bool) bool {
	if num {
		return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_'
	}
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
}

func isNum(r byte, e bool) bool {
	x := (r >= '0' && r <= '9') || r == '.' || r == '+' || r == '-'
	if e {
		x = x || r == 'e' || r == 'E'
	}
	return x
}
