package nj

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

type ValueIO Value

var (
	ioWriterType = reflect.TypeOf((*io.Writer)(nil)).Elem()
	ioReaderType = reflect.TypeOf((*io.Reader)(nil)).Elem()
	ioCloserType = reflect.TypeOf((*io.Closer)(nil)).Elem()
)

var (
	ReaderProto = Func("Reader", nil).Object().Merge(nil,
		Str("read"), Func("", func(e *Env) {
			buf := ioRead(e)
			_ = buf == nil && e.SetA(Nil) || e.SetA(UnsafeStr(buf))
		}, "$f() -> string", "\tread all bytes as string, return nil if EOF reached",
			"$f(n: int) -> string", "\tread `n` bytes as string"),
		Str("readbytes"), Func("", func(e *Env) {
			buf := ioRead(e)
			_ = buf == nil && e.SetA(Nil) || e.SetA(Bytes(buf))
		}, "$f() -> bytes", "\tread all bytes, return nil if EOF reached",
			"$f(n: int) -> bytes", "\tread `n` bytes"),
		Str("readbuf"), Func("", func(e *Env) {
			rn, err := e.Object(-1).Prop("_f").Interface().(io.Reader).Read(e.Array(0).Unwrap().([]byte))
			e.A = Array(Int(rn), ValueOf(err)) // return in Go style
		}, "$f(buf: bytes) -> [int, go.error]", "\tread into `buf` and return in Go style"),
		Str("readlines"), Func("", func(e *Env) {
			f := e.Object(-1).Prop("_f").Interface().(io.Reader)
			delim := e.Object(-1).Prop("delim").ToStr("\n")
			if e.Get(0) == Nil {
				buf, err := ioutil.ReadAll(f)
				internal.PanicErr(err)
				parts := bytes.Split(buf, []byte(delim))
				var res []Value
				for i, line := range parts {
					if i < len(parts)-1 {
						line = append(line, delim...)
					}
					res = append(res, UnsafeStr(line))
				}
				e.A = Array(res...)
				return
			}
			for cb, rd := e.Object(0), bufio.NewReader(f); ; {
				line, err := rd.ReadString(delim[0])
				if len(line) > 0 {
					if v := cb.MustCall(Str(line)); v == False {
						e.A = v
						return
					}
				}
				if err != nil {
					if err != io.EOF {
						panic(err)
					}
					break
				}
			}
			e.A = Nil
		},
			"$f() -> array", "\tread the whole file and return lines as an array",
			"$f(f: function)", "\tfor every line read, `f(line)` will be called", "\tto exit the reading, return `false` in `f`",
		),
	)

	WriterProto = Func("Writer", nil).Object().Merge(nil,
		Str("write"), Func("", func(e *Env) {
			wn, err := e.Object(-1).Prop("_f").Interface().(io.Writer).Write(e.Get(0).ToBytes())
			internal.PanicErr(err)
			e.A = Int(wn)
		}, "$f(buf: string|bytes) int", "\twrite `buf` to writer"),
		Str("pipe"), Func("", func(e *Env) {
			var wn int64
			var err error
			if n := e.Get(1).ToInt64(0); n > 0 {
				wn, err = io.CopyN(NewWriter(e.Get(-1)), NewReader(e.Get(0)), n)
			} else {
				wn, err = io.Copy(NewWriter(e.Get(-1)), NewReader(e.Get(0)))
			}
			internal.PanicErr(err)
			e.A = Int64(wn)
		}, "$f(r: Reader, n?: int) -> int", "\tcopy (at most `n`) bytes from `r` to writer, return number of bytes copied"),
	)

	SeekerProto = Func("Seeker", nil).Object().Merge(nil,
		Str("seek"), Func("", func(e *Env) {
			f := e.Object(-1).Prop("_f").Interface().(io.Seeker)
			wn, err := f.Seek(e.Int64(0), e.Int(1))
			internal.PanicErr(err)
			e.A = Int64(wn)
		}, "$f(offset: int, whence: int) -> int"))

	CloserProto = Func("Closer", nil).Object().Merge(nil,
		Str("close"), Func("", func(e *Env) {
			internal.PanicErr(e.Object(-1).Prop("_f").Interface().(io.Closer).Close())
		}, "$f()"))

	ReadWriterProto = Func("ReadWriter", nil).Object().Merge(ReaderProto).Merge(WriterProto)

	ReadCloserProto = Func("ReadCloser", nil).Object().Merge(ReaderProto).Merge(CloserProto)

	WriteCloserProto = Func("WriteCloser", nil).Object().Merge(WriterProto).Merge(CloserProto)

	ReadWriteCloserProto = Func("ReadWriteCloser", nil).Object().Merge(ReadWriterProto).Merge(CloserProto)

	ReadWriteSeekCloserProto = Func("ReadWriteSeekCloserProto", nil).Object().Merge(ReadWriteCloserProto).Merge(SeekerProto)
)

// NewReader creates an io.Reader from value if possible
func NewReader(v Value) io.Reader {
	switch v.Type() {
	case typ.Native:
		switch rd := v.Interface().(type) {
		case io.Reader:
			return rd
		case []byte:
			return bytes.NewReader(rd)
		}
	case typ.String:
		return strings.NewReader(v.Str())
	}
	return ValueIO(v)
}

// NewWriter creates an io.Writer from value if possible
func NewWriter(v Value) io.Writer {
	switch v.Type() {
	case typ.Native:
		switch rd := v.Interface().(type) {
		case io.Writer:
			return rd
		case []byte:
			w := bytes.NewBuffer(rd)
			w.Reset()
			return w
		}
	}
	return ValueIO(v)
}

// NewCloser creates an io.Closer from value if possible
func NewCloser(v Value) io.Closer {
	if v.Type() == typ.Native {
		if rd, ok := v.Interface().(io.Closer); ok {
			return rd
		}
	}
	return ValueIO(v)
}

func (m ValueIO) Read(p []byte) (int, error) {
	switch Value(m).Type() {
	case typ.Native:
		if rd, _ := Value(m).Interface().(io.Reader); rd != nil {
			return rd.Read(p)
		}
	case typ.Object:
		if rb := Value(m).Object().Prop("readbuf"); rb.IsObject() {
			v, err := rb.Object().Call(Bytes(p))
			if err != nil {
				return 0, err
			}
			t := v.Is(typ.Array, "ValueIO.Read: use readbuf()").Array()
			n := t.Get(0).Is(typ.Number, "ValueIO.Read: (int, error)").Int()
			err, _ = t.Get(1).Interface().(error)
			return int(n), err
		}
		if rb := Value(m).Object().Prop("readbytes"); rb.IsObject() {
			v, err := rb.Object().Call(Int(len(p)))
			if err != nil {
				return 0, err
			} else if v == Nil {
				return 0, io.EOF
			}
			return copy(p, v.ToBytes()), nil
		}
		if rb := Value(m).Object().Prop("read"); rb.IsObject() {
			v, err := rb.Object().Call(Int(len(p)))
			if err != nil {
				return 0, err
			} else if v == Nil {
				return 0, io.EOF
			}
			return copy(p, v.ToStr("")), nil
		}
	}
	return 0, fmt.Errorf("reader not implemented")
}

func (m ValueIO) Write(p []byte) (int, error) {
	switch Value(m).Type() {
	case typ.Native:
		if rd, _ := Value(m).Interface().(io.Writer); rd != nil {
			return rd.Write(p)
		}
	case typ.Object:
		if rb := Value(m).Object().Prop("write"); rb.IsObject() {
			v, err := rb.Object().Call(Bytes(p))
			if err != nil {
				return 0, err
			}
			return v.Is(typ.Number, "ValueIO.Write: (int, error)").Int(), nil
		}
	}
	return 0, fmt.Errorf("writer not implemented")
}

func (m ValueIO) Close() error {
	switch Value(m).Type() {
	case typ.Native:
		if rd, _ := Value(m).Interface().(io.Closer); rd != nil {
			return rd.Close()
		}
	case typ.Object:
		if rb := Value(m).Object().Prop("close"); rb.IsObject() {
			if _, err := rb.Object().Call(); err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("closer not implemented")
}

func ioRead(e *Env) []byte {
	f := e.Object(-1).Prop("_f").Interface().(io.Reader)
	if n := e.Get(0); n.Type() == typ.Number {
		p := make([]byte, n.ToInt64(0))
		rn, err := f.Read(p)
		if err == nil || rn > 0 {
			return p[:rn]
		} else if err == io.EOF {
			return nil
		} else {
			panic(err)
		}
	}
	buf, err := ioutil.ReadAll(f)
	internal.PanicErr(err)
	return buf
}
