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
	unnamedCounter int64
	debugMode      = os.Getenv("njd") != ""
)

func UnnamedFunc() string {
	return fmt.Sprintf("<native-%d>", atomic.AddInt64(&unnamedCounter, 1))
}

func UnnamedLoadString() string {
	return fmt.Sprintf("<memory-%d>", atomic.AddInt64(&unnamedCounter, 1))
}

func Unnamed() string {
	return "tmp." + strconv.FormatInt(atomic.AddInt64(&unnamedCounter, 1), 10)
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

func PanicNotEnoughArgs(a string) {
	panic("not enough arguments to call " + a)
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
	a := struct {
		b string
		c int
	}{s, len(s)}
	var x []byte
	*(*[3]int)(unsafe.Pointer(&x)) = *(*[3]int)(unsafe.Pointer(&a))
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
