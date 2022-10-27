package bas

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/coyove/nj/typ"
)

type valueIO Value

// Reader creates an io.Reader from value, Read() may fail if value doesn't support reading.
func (v Value) Reader() io.Reader {
	switch rd := v.Interface().(type) {
	case io.Reader:
		return rd
	case []byte:
		return bytes.NewReader(rd)
	case string:
		return strings.NewReader(v.Str())
	}
	return valueIO(v)
}

// Writer creates an io.Writer from value, Write() may fail if value doesn't support writing.
func (v Value) Writer() io.Writer {
	switch rd := v.Interface().(type) {
	case io.Writer:
		return rd
	case []byte:
		w := bytes.NewBuffer(rd)
		w.Reset()
		return w
	}
	return valueIO(v)
}

// Closer creates an io.Closer from value, Close() may fail if value doesn't support closing.
func (v Value) Closer() io.Closer {
	if rd, ok := v.Interface().(io.Closer); ok {
		return rd
	}
	return valueIO(v)
}

func (m valueIO) Read(p []byte) (int, error) {
	switch Value(m).Type() {
	case typ.Native:
		if rd, _ := Value(m).Interface().(io.Reader); rd != nil {
			return rd.Read(p)
		}
	case typ.Object:
		if rb := Value(m).Object().Get(Str("read2")); rb.IsObject() {
			t := rb.Object().Call(nil, Bytes(p)).AssertShape("(i, EN)", "Reader.read2").Native()
			if t.Get(1).IsError() {
				return t.Get(0).Int(), t.Get(1).Error()
			}
			return t.Get(0).Int(), nil
		}
		if rb := Value(m).Object().Get(Str("read")); rb.IsObject() {
			switch v := rb.Object().Call(nil, Int(len(p))); v.Type() {
			case typ.Nil:
				return 0, io.EOF
			case typ.String:
				return copy(p, v.Str()), nil
			case typ.Native:
				return copy(p, v.AssertShape("sB", "Reader.read").Native().Unwrap().([]byte)), nil
			default:
				v.AssertShape("sB", "Reader.read")
			}
		}
	}
	return 0, fmt.Errorf("reader not implemented")
}

func (m valueIO) WriteString(p string) (int, error) {
	return m.Write([]byte(p))
}

func (m valueIO) Write(p []byte) (int, error) {
	switch Value(m).Type() {
	case typ.Native:
		if rd, _ := Value(m).Interface().(io.Writer); rd != nil {
			return rd.Write(p)
		}
	case typ.Object:
		if rb := Value(m).Object().Get(Str("write")); rb.IsObject() {
			v := rb.Object().Call(nil, Bytes(p))
			if v.IsError() {
				return 0, v.Error()
			}
			return v.AssertNumber("Writer.write").Int(), nil
		}
		if rb := Value(m).Object().Get(Str("write2")); rb.IsObject() {
			t := rb.Object().Call(nil, Bytes(p)).AssertShape("(i, EN)", "Writer.write2").Native()
			if t.Get(1).IsError() {
				return t.Get(0).Int(), t.Get(1).Error()
			}
			return t.Get(0).Int(), nil
		}
	}
	return 0, fmt.Errorf("writer not implemented")
}

func (m valueIO) Close() error {
	switch Value(m).Type() {
	case typ.Native:
		if rd, _ := Value(m).Interface().(io.Closer); rd != nil {
			return rd.Close()
		}
	case typ.Object:
		if rb := Value(m).Object().Get(Str("close")); rb.IsObject() {
			if v := rb.Object().Call(nil); v.IsError() {
				return v.Error()
			}
			return nil
		}
	}
	return fmt.Errorf("closer not implemented")
}

type ioReadlinesStruct struct {
	rd    *bufio.Reader
	delim byte
	bytes bool
}
