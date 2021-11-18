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
	ReaderProto = Map(Str("__name"), Str("reader"),
		Str("read"), Func2("read", func(rx, n Value) Value {
			f := rx.Table().GetString("_f").Interface().(io.Reader)
			switch n.Type() {
			case typ.Number:
				p := make([]byte, n.MaybeInt(0))
				rn, err := f.Read(p)
				if err == nil || rn > 0 {
					return Bytes(p[:rn])
				}
				if err == io.EOF {
					return Nil
				}
				panic(err)
			default:
				buf, err := ioutil.ReadAll(f)
				internal.PanicErr(err)
				return Bytes(buf)
			}
		}, "read() string", "\tread all bytes, return nil if EOF reached", "read(n: int) string", "\tread n bytes"),
		Str("readbuf"), Func2("readbuf", func(rx, n Value) Value {
			rn, err := rx.Table().GetString("_f").Interface().(io.Reader).Read(n.Interface().([]byte))
			return Array(Int(int64(rn)), Val(err)) // return in Go style
		}, "$f(buf: bytes) array", "\tread into buf and return { bytes_read, error } in Go style"),
		Str("readlines"), Func2("readlines", func(rx, cb Value) Value {
			f := rx.Table().GetString("_f").Interface().(io.Reader)
			delim := rx.Table().GetString("delim").MaybeStr("\n")
			if cb == Nil {
				buf, err := ioutil.ReadAll(f)
				if err != nil {
					panic(err)
				}
				parts := bytes.Split(buf, []byte(delim))
				var res []Value
				for i, line := range parts {
					if i < len(parts)-1 {
						line = append(line, delim...)
					}
					res = append(res, Bytes(line))
				}
				return Array(res...)
			}
			for rd := bufio.NewReader(f); ; {
				line, err := rd.ReadString(delim[0])
				if len(line) > 0 {
					if v, err := cb.MustFunc("callback").Call(Str(line)); err != nil {
						panic(err)
					} else if v != Nil {
						return v
					}
				}
				if err != nil {
					if err != io.EOF {
						panic(err)
					}
					break
				}
			}
			return Nil
		},
			"readlines() array", "\tread the whole file and return lines as a table array",
			"readlines(f: function)", "\tfor every line read, f(line) will be called", "\tto exit the reading, return anything other than nil in f",
		),
	).Table()

	WriterProto = Map(Str("__name"), Str("writer"),
		Str("write"), Func2("write", func(rx, buf Value) Value {
			f := rx.Table().GetString("_f").Interface().(io.Writer)
			wn, err := f.Write([]byte(buf.MustStr("")))
			internal.PanicErr(err)
			return Int(int64(wn))
		}, "$f({w}: value, buf: string) int", "\twrite buf to w"),
		Str("pipe"), Func3("pipe", func(dest, src, n Value) Value {
			var wn int64
			var err error
			if n := n.MaybeInt(0); n > 0 {
				wn, err = io.CopyN(NewWriter(dest), NewReader(src), n)
			} else {
				wn, err = io.Copy(NewWriter(dest), NewReader(src))
			}
			internal.PanicErr(err)
			return Int(wn)
		}, "$f({w}: value, r: value) int", "\tcopy bytes from r to w, return number of bytes copied",
			"$f({w}: value, r: value, n: int) int", "\tcopy at most n bytes from r to w"),
	).Table()

	SeekerProto = Map(Str("__name"), Str("seeker"),
		Str("seek"), Func3("seek", func(rx, off, where Value) Value {
			f := rx.Table().GetString("_f").Interface().(io.Seeker)
			wn, err := f.Seek(off.MustInt("offset"), int(where.MustInt("where")))
			internal.PanicErr(err)
			return Int(int64(wn))
		}, "")).Table()

	CloserProto = Map(Str("__name"), Str("closer"),
		Str("close"), Func1("close", func(rx Value) Value {
			internal.PanicErr(rx.Table().GetString("_f").Interface().(io.Closer).Close())
			return Nil
		}, "")).Table()

	ReadWriterProto = ReaderProto.Copy().Merge(WriterProto, Str("__name"), Str("readwriter"))

	ReadCloserProto = ReaderProto.Copy().Merge(CloserProto, Str("__name"), Str("readcloser"))

	WriteCloserProto = WriterProto.Copy().Merge(CloserProto, Str("__name"), Str("writecloser"))

	ReadWriteCloserProto = ReadWriterProto.Copy().Merge(CloserProto, Str("__name"), Str("readwritecloser"))

	ReadWriteSeekCloserProto = ReadWriteCloserProto.Copy().Merge(SeekerProto, Str("__name"), Str("readwriteseekcloser"))
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
	case typ.Table:
		if rb := Value(m).Table().GetString("readbuf"); rb.Type() == typ.Func {
			v, err := rb.Func().Call(Val(p))
			if err != nil {
				return 0, err
			}
			t := v.MustTable("TableIO.Read: readbuf()")
			n := t.Get(Int(0)).MustInt("TableIO.Read: (int, error)")
			err, _ = t.Get(Int(1)).Interface().(error)
			return int(n), err
		}
		if rb := Value(m).Table().GetString("read"); rb.Type() == typ.Func {
			v, err := rb.Func().Call(Int(int64(len(p))))
			if err != nil {
				return 0, err
			}
			if v == Nil {
				return 0, io.EOF
			}
			return copy(p, v.MustStr("TableIO.Read: read()")), nil
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
	case typ.Table:
		if rb := Value(m).Table().GetString("write"); rb.Type() == typ.Func {
			v, err := rb.Func().Call(Bytes(p))
			if err != nil {
				return 0, err
			}
			return int(v.MustInt("TableIO.Write: (int, error)")), nil
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
	case typ.Table:
		if rb := Value(m).Table().GetString("close"); rb.Type() == typ.Func {
			v, err := rb.Func().Call()
			if err != nil {
				return err
			}
			return v.Interface().(error)
		}
	}
	return fmt.Errorf("closer not implemented")
}
