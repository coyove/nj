package internal

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"sync/atomic"
	"unicode/utf8"
	"unsafe"
)

var (
	NewFunc    func(name string, variadic bool, np byte, ss uint16, locals, caps []string, code Packet) interface{}
	NewProgram func(coreStack, top, symbols, funcs interface{}) interface{}

	unnamedCounter int64
	debugMode      = os.Getenv("njd") != ""
)

func UnnamedFunc() string {
	return fmt.Sprintf("<native-%d>", atomic.AddInt64(&unnamedCounter, 1))
}

func UnnamedLoadString() string {
	return fmt.Sprintf("<memory-%d>", atomic.AddInt64(&unnamedCounter, 1))
}

func ShouldNotHappen(args ...interface{}) {
	if len(args) > 0 {
		panic(fmt.Errorf("fatal: should not happen, bad values: %v", args...))
	}
	panic(fmt.Errorf("fatal: should not happen"))
}

func Panic(msg string, args ...interface{}) {
	panic(fmt.Errorf(msg, args...))
}

func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func IsDebug() bool {
	return debugMode
}

func CatchError(err *error) {
	if r := recover(); r != nil {
		if IsDebug() {
			log.Println(string(debug.Stack()))
		}

		*err, _ = r.(error)
		if *err == nil {
			*err = fmt.Errorf("%v", r)
		}
	}
}

func CloseBuffer(p *bytes.Buffer, suffix string) {
	for p.Len() > 0 {
		b := p.Bytes()[p.Len()-1]
		if b == ' ' || b == ',' {
			p.Truncate(p.Len() - 1)
		} else {
			break
		}
	}
	p.WriteString(suffix)
}

func IfStr(v bool, t, f string) string {
	if v {
		return t
	}
	return f
}

func IfInt(v bool, t, f int) int {
	if v {
		return t
	}
	return f
}

func WriteString(w io.Writer, s string) (int, error) {
	type a struct {
		b string
		c int
	}
	var x []byte
	*(*a)(unsafe.Pointer(&x)) = a{s, len(s)}
	return w.Write(x)
}

func IfQuote(v bool, s string) string {
	if v {
		return strconv.Quote(s)
	}
	return s
}

func Or(a, b interface{}) interface{} {
	if a != nil {
		return a
	}
	return b
}

func ParseNumber(v string) (vf float64, vi int64, isInt bool, err error) {
	i, err := strconv.ParseInt(v, 0, 64)
	if err == nil {
		return 0, i, true, nil
	}
	if err.(*strconv.NumError).Err == strconv.ErrRange {
		i, err := strconv.ParseUint(v, 0, 64)
		if err == nil {
			return 0, (int64(i)), true, nil
		}
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, 0, false, fmt.Errorf("invalid number format: %q", v)
	}
	return f, 0, false, nil
}

func LineOf(text string, line int) (string, bool) {
	for line > 0 {
		idx := strings.IndexByte(text, '\n')
		line--
		if line == 0 {
			if idx >= 0 {
				text = text[:idx]
			}
			if trunc := 256; len(text) > trunc {
				for i := trunc - 1; i >= 0; i-- {
					if r, _ := utf8.DecodeLastRuneInString(text[:i]); r != utf8.RuneError {
						text = text[:i] + " ... truncated code"
						break
					}
				}
			}
			return text, true
		}
		if idx == -1 {
			break
		}
		text = text[idx+1:]
	}
	return "", false
}

func SanitizeName(s string) string {
	tn := []byte(s)
	for i := range tn {
		switch tn[i] {
		case '*':
			tn[i] = 'p'
		case '.', '[', ']':
			tn[i] = '_'
		}
	}
	return *(*string)(unsafe.Pointer(&tn))
}

func StringifyTo(w io.Writer, i interface{}) {
	switch s := i.(type) {
	case fmt.Stringer:
		WriteString(w, s.String())
	case error:
		WriteString(w, s.Error())
	default:
		WriteString(w, "<")
		WriteString(w, reflect.TypeOf(i).String())
		WriteString(w, ">")
	}
}

func Fprintf(w io.Writer, format string, args ...interface{}) {
	tmp := bytes.Buffer{}
	ai := 0
NEXT:
	for len(format) > 0 {
		idx := strings.Index(format, "%")
		if idx == -1 {
			WriteString(w, format)
			break
		}
		if idx == len(format)-1 {
			WriteString(w, "%?(NOVERB)")
			break
		}
		WriteString(w, format[:idx])
		if format[idx+1] == '%' {
			WriteString(w, "%")
			format = format[idx+2:]
			continue
		}

		tmp.Reset()
		tmp.WriteByte('%')
		format = format[idx+1:]

		preferNumber := ' '
		for found := false; len(format) > 0 && !found; {
			head := format[0]
			tmp.WriteByte(head)
			format = format[1:]
			switch head {
			case 'b', 'd', 'o', 'O', 'c', 'U':
				preferNumber = 'i'
				found = true
			case 'f', 'F', 'g', 'G', 'e', 'E':
				preferNumber = 'f'
				found = true
			case 's', 'q', 'x', 'X', 'v', 't', 'p':
				found = true
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', '-', '+', '#', ' ':
			default:
				WriteString(w, tmp.String()+"(BAD)")
				goto NEXT
			}
		}
		if ai >= len(args) {
			WriteString(w, tmp.String()+"(MISSING)")
		} else {
			v := args[ai]
			if sn, ok := v.(SprintfNumber); ok {
				if preferNumber == 'i' {
					v = sn.Int
				} else if preferNumber == 'f' {
					v = sn.Float
				} else if sn.IsInt {
					v = sn.Int
				} else {
					v = sn.Float
				}
			}
			fmt.Fprintf(w, tmp.String(), v)
		}
		ai++
	}
}

type SprintfNumber struct {
	Int   int64
	Float float64
	IsInt bool
}
