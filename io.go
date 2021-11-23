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
	ReaderProto = Proto(ObjectLib.Object(), Str("__name"), Str("reader"),
		Str("read"), Func("", func(e *Env) {
			f := e.Object(-1).Gets("_f").Interface().(io.Reader)
			switch n := e.Get(0); n.Type() {
			case typ.Number:
				p := make([]byte, n.ToInt64(0))
				rn, err := f.Read(p)
				if err == nil || rn > 0 {
					e.A = Bytes(p[:rn])
				} else if err == io.EOF {
					e.A = Nil
				} else {
					panic(err)
				}
			default:
				buf, err := ioutil.ReadAll(f)
				internal.PanicErr(err)
				e.A = Bytes(buf)
			}
		}, "$f() -> string", "\tread all bytes, return nil if EOF reached",
			"$f(n: int) -> string", "\tread `n` bytes"),
		Str("readbuf"), Func("", func(e *Env) {
			rn, err := e.Object(-1).Gets("_f").Interface().(io.Reader).Read(e.Interface(0).([]byte))
			e.A = Array(Int(rn), Val(err)) // return in Go style
		}, "$f(buf: bytes) array", "\tread into `buf` and return [int, go.error] in Go style"),
		Str("readlines"), Func("", func(e *Env) {
			f := e.Object(-1).Gets("_f").Interface().(io.Reader)
			delim := e.Object(-1).Gets("delim").ToStr("\n")
			if e.Get(0) == Nil {
				buf, err := ioutil.ReadAll(f)
				internal.PanicErr(err)
				parts := bytes.Split(buf, []byte(delim))
				var res []Value
				for i, line := range parts {
					if i < len(parts)-1 {
						line = append(line, delim...)
					}
					res = append(res, Bytes(line))
				}
				e.A = Array(res...)
				return
			}
			for cb, rd := e.Object(0), bufio.NewReader(f); ; {
				line, err := rd.ReadString(delim[0])
				if len(line) > 0 {
					if v := cb.MustCall(Str(line)); v != Nil {
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
			"readlines() array", "\tread the whole file and return lines as a table array",
			"readlines(f: function)", "\tfor every line read, f(line) will be called", "\tto exit the reading, return anything other than nil in f",
		),
	).Object()

	WriterProto = Proto(ObjectLib.Object(), Str("__name"), Str("writer"),
		Str("write"), Func("", func(e *Env) {
			f := e.Object(-1).Gets("_f").Interface().(io.Writer)
			wn, err := f.Write([]byte(e.Str(0)))
			internal.PanicErr(err)
			e.A = Int(wn)
		}, "$f({w}: value, buf: string) int", "\twrite buf to w"),
		Str("pipe"), Func("pipe", func(e *Env) {
			var wn int64
			var err error
			if n := e.Get(1).ToInt64(0); n > 0 {
				wn, err = io.CopyN(NewWriter(e.Get(-1)), NewReader(e.Get(0)), n)
			} else {
				wn, err = io.Copy(NewWriter(e.Get(-1)), NewReader(e.Get(0)))
			}
			internal.PanicErr(err)
			e.A = Int64(wn)
		}, "$f({w}: value, r: value) int", "\tcopy bytes from r to w, return number of bytes copied",
			"$f({w}: value, r: value, n: int) int", "\tcopy at most n bytes from r to w"),
	).Object()

	SeekerProto = Proto(ObjectLib.Object(), Str("__name"), Str("seeker"),
		Str("seek"), Func3("seek", func(rx, off, where Value) Value {
			f := rx.Object().Gets("_f").Interface().(io.Seeker)
			wn, err := f.Seek(off.MustInt64("offset"), int(where.MustInt64("where")))
			internal.PanicErr(err)
			return Int64(int64(wn))
		}, "")).Object()

	CloserProto = Proto(ObjectLib.Object(), Str("__name"), Str("closer"),
		Str("close"), Func("", func(e *Env) {
			internal.PanicErr(e.Object(-1).Gets("_f").Interface().(io.Closer).Close())
		}, "")).Object()

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
	case typ.Object:
		if rb := Value(m).Object().Gets("readbuf"); rb.IsObject() {
			v, err := rb.Func().Call(Val(p))
			if err != nil {
				return 0, err
			}
			t := v.Is(typ.Array, "ValueIO.Read: use readbuf()").Array()
			n := t.Get(Int64(0)).Is(typ.Number, "ValueIO.Read: (int, error)").Int()
			err, _ = t.Get(Int64(1)).Interface().(error)
			return int(n), err
		}
		if rb := Value(m).Object().Gets("read"); rb.IsObject() {
			v, err := rb.Func().Call(Int(len(p)))
			if err != nil {
				return 0, err
			}
			if v == Nil {
				return 0, io.EOF
			}
			return copy(p, []byte(v.Is(typ.String, "ValueIO.Read: use read()").Str())), nil
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
		if rb := Value(m).Object().Gets("write"); rb.IsObject() {
			v, err := rb.Func().Call(Bytes(p))
			if err != nil {
				return 0, err
			}
			return int(v.MustInt64("ValueIO.Write: (int, error)")), nil
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
		if rb := Value(m).Object().Gets("close"); rb.IsObject() {
			v, err := rb.Func().Call()
			if err != nil {
				return err
			}
			return v.Interface().(error)
		}
	}
	return fmt.Errorf("closer not implemented")
}
