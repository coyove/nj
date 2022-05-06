package internal

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"unsafe"
)

const UnnamedFunc = "<native>"

var (
	GrowEnvStack func(env unsafe.Pointer, sz int)
	SetEnvStack  func(env unsafe.Pointer, stack unsafe.Pointer)
	SetObjFun    func(obj unsafe.Pointer, fun unsafe.Pointer)
)

type TransparentError struct{}

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
	return os.Getenv("njd") != ""
}

func CatchError(err *error) {
	if r := recover(); r != nil {
		if IsDebug() {
			log.Println(string(debug.Stack()))
		}

		if x, ok := r.(interface{ TransparentError() TransparentError }); ok {
			*err = x.(error)
			return
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

func WriteString(w io.Writer, s string) (int, error) {
	x := struct {
		a string
		b int
	}{s, len(s)}
	return w.Write(*(*[]byte)(unsafe.Pointer(&x)))
}

func IfQuote(v bool, s string) string {
	if v {
		return strconv.Quote(s)
	}
	return s
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
