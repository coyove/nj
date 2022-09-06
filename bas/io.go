package bas

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/typ"
)

type valueIO Value

func init() {
	ioReadlinesIter = NewEmptyNativeMeta("readlines", Proto.Native)
	ioReadlinesIter.Next = func(a *Native, k Value) Value {
		if k.IsNil() {
			k = Array(Int(-1), Nil)
		}
		var er error
		if s := a.any.(*ioReadlinesStruct); s.bytes {
			line, err := s.rd.ReadBytes(s.delim)
			if len(line) > 0 {
				k.Native().Set(0, Int(k.Native().Get(0).Int()+1))
				k.Native().Set(1, Bytes(line))
				return k
			}
			er = err
		} else {
			line, err := s.rd.ReadString(s.delim)
			if len(line) > 0 {
				k.Native().Set(0, Int(k.Native().Get(0).Int()+1))
				k.Native().Set(1, Str(line))
				return k
			}
			er = err
		}
		if er == io.EOF {
			return Nil
		}
		panic(er)
	}

	Proto.Reader.Proto.
		SetMethod("read", func(e *Env) {
			buf := ioRead(e)
			_ = buf == nil && e.SetA(Nil) || e.SetA(UnsafeStr(buf))
		}).
		SetMethod("readbytes", func(e *Env) {
			buf := ioRead(e)
			_ = buf == nil && e.SetA(Nil) || e.SetA(Bytes(buf))
		}).
		SetMethod("readbuf", func(e *Env) {
			rn, err := e.A.Reader().Read(e.Native(0).Unwrap().([]byte))
			e.A = Array(Int(rn), Error(e, err)) // return in Go style
		}).
		SetMethod("readlines", func(e *Env) {
			e.A = NewNativeWithMeta(&ioReadlinesStruct{
				rd:    bufio.NewReader(e.A.Reader()),
				delim: e.Get(0).NilStr("\n")[0],
				bytes: e.Get(1).NilBool(),
			}, ioReadlinesIter).ToValue()
		})

	Proto.Writer.Proto.
		SetMethod("write", func(e *Env) {
			wn, err := e.A.Writer().Write(ToReadonlyBytes(e.Get(0)))
			internal.PanicErr(err)
			e.A = Int(wn)
		}).
		SetMethod("writebytes", func(e *Env) {
			wn, err := e.A.Writer().Write(ToReadonlyBytes(e.Get(0)))
			internal.PanicErr(err)
			e.A = Int(wn)
		}).
		SetMethod("writebuf", func(e *Env) {
			wn, err := e.A.Writer().Write(ToReadonlyBytes(e.Get(0)))
			e.A = Array(Int(wn), Error(e, err))
		}).
		SetMethod("pipe", func(e *Env) {
			var wn int64
			var err error
			if n := e.Get(1).NilInt64(0); n > 0 {
				wn, err = io.CopyN(e.Get(-1).Writer(), e.Get(0).Reader(), n)
			} else {
				wn, err = io.Copy(e.Get(-1).Writer(), e.Get(0).Reader())
			}
			internal.PanicErr(err)
			e.A = Int64(wn)
		})

	Proto.Closer.Proto.
		SetMethod("close", func(e *Env) {
			internal.PanicErr(e.A.Closer().Close())
		})

	Proto.ReadWriter.Proto.Merge(Proto.Reader.Proto).Merge(Proto.Writer.Proto)

	Proto.ReadCloser.Proto.Merge(Proto.Reader.Proto).Merge(Proto.Closer.Proto)

	Proto.WriteCloser.Proto.Merge(Proto.Writer.Proto).Merge(Proto.Closer.Proto)

	Proto.ReadWriteCloser.Proto.Merge(Proto.ReadWriter.Proto).Merge(Proto.Closer.Proto)
}

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
		if rb := Value(m).Object().Prop("readbuf"); rb.IsObject() {
			t := rb.Object().Call(nil, Bytes(p)).AssertShape("(i, Ev)", "Reader.readbuf").Native()
			if IsError(t.Get(1)) {
				return t.Get(0).Int(), ToErrorRootCause(t.Get(1)).(error)
			}
			return t.Get(0).Int(), nil
		}
		if rb := Value(m).Object().Prop("readbytes"); rb.IsObject() {
			v := rb.Object().Call(nil, Int(len(p)))
			if v == Nil {
				return 0, io.EOF
			}
			return copy(p, v.AssertPrototype(Proto.Bytes, "Reader.readbytes").Native().Unwrap().([]byte)), nil
		}
		if rb := Value(m).Object().Prop("read"); rb.IsObject() {
			v := rb.Object().Call(nil, Int(len(p)))
			if v == Nil {
				return 0, io.EOF
			}
			return copy(p, v.AssertType(typ.String, "Reader.read").Str()), nil
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
		if rb := Value(m).Object().Prop("write"); rb.IsObject() {
			v := rb.Object().Call(nil, UnsafeStr(p))
			if IsError(v) {
				return 0, ToError(v)
			}
			return v.AssertType(typ.Number, "Writer.write").Int(), nil
		}
		if rb := Value(m).Object().Prop("writebytes"); rb.IsObject() {
			v := rb.Object().Call(nil, Bytes(p))
			if IsError(v) {
				return 0, ToError(v)
			}
			return v.AssertType(typ.Number, "Writer.writebytes").Int(), nil
		}
		if rb := Value(m).Object().Prop("writebuf"); rb.IsObject() {
			t := rb.Object().Call(nil, Bytes(p)).AssertShape("(i, Ev)", "Writer.writebuf").Native()
			if IsError(t.Get(1)) {
				return t.Get(0).Int(), ToErrorRootCause(t.Get(1)).(error)
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
		if rb := Value(m).Object().Prop("close"); rb.IsObject() {
			if v := rb.Object().Call(nil); IsError(v) {
				return ToError(v)
			}
			return nil
		}
	}
	return fmt.Errorf("closer not implemented")
}

func ioRead(e *Env) []byte {
	f := e.A.Reader()
	if n := e.Get(0); n.Type() == typ.Number {
		p := make([]byte, n.Int())
		rn, err := f.Read(p)
		if err == nil || rn > 0 {
			return p[:rn]
		} else if err == io.EOF {
			return nil
		}
		panic(err)
	}
	buf, err := ioutil.ReadAll(f)
	internal.PanicErr(err)
	return buf
}

var ioReadlinesIter *NativeMeta

type ioReadlinesStruct struct {
	rd    *bufio.Reader
	delim byte
	bytes bool
}
