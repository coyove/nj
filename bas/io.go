package bas

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

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
		return Error(nil, er)
	}

	NativeMetaProto.Reader.Proto.
		AddMethod("read", func(e *Env) {
			buf, err := func(e *Env) ([]byte, error) {
				f := e.A.Reader()
				if n := e.Get(0); n.Type() == typ.Number {
					p := make([]byte, n.Int())
					rn, err := f.Read(p)
					if err == nil || rn > 0 {
						return p[:rn], nil
					} else if err == io.EOF {
						return nil, nil
					}
					return nil, err
				}
				return ioutil.ReadAll(f)
			}(e)
			_ = err != nil && e.SetA(Error(e, err)) || e.SetA(Bytes(buf))
		}).
		AddMethod("read2", func(e *Env) {
			rn, err := e.A.Reader().Read(e.Shape(0, "B").Native().Unwrap().([]byte))
			e.A = Array(Int(rn), Error(e, err)) // return in Go style
		}).
		AddMethod("readlines", func(e *Env) {
			e.A = NewNativeWithMeta(&ioReadlinesStruct{
				rd:    bufio.NewReader(e.A.Reader()),
				delim: e.StrDefault(0, "\n", 1)[0],
				bytes: e.Shape(1, "Nb").IsTrue(),
			}, ioReadlinesIter).ToValue()
		}).
		SetPrototype(Proto.Native)

	NativeMetaProto.Writer.Proto.
		AddMethod("write", func(e *Env) {
			wn, err := Write(e.A.Writer(), e.Get(0))
			_ = err == nil && e.SetA(Int(wn)) || e.SetA(Error(e, err))
		}).
		AddMethod("write2", func(e *Env) {
			wn, err := Write(e.A.Writer(), e.Get(0))
			e.A = Array(Int(wn), Error(e, err))
		}).
		AddMethod("pipe", func(e *Env) {
			var wn int64
			var err error
			if n := e.IntDefault(1, 0); n > 0 {
				wn, err = io.CopyN(e.Get(-1).Writer(), e.Get(0).Reader(), int64(n))
			} else {
				wn, err = io.Copy(e.Get(-1).Writer(), e.Get(0).Reader())
			}
			_ = err == nil && e.SetA(Int64(wn)) || e.SetA(Error(e, err))
		}).
		SetPrototype(Proto.Native)

	NativeMetaProto.Closer.Proto.
		AddMethod("close", func(e *Env) {
			e.A = Error(e, e.A.Closer().Close())
		}).
		SetPrototype(Proto.Native)

	NativeMetaProto.ReadWriter.Proto.
		Merge(NativeMetaProto.Reader.Proto).
		Merge(NativeMetaProto.Writer.Proto).SetPrototype(Proto.Native)

	NativeMetaProto.ReadCloser.Proto.
		Merge(NativeMetaProto.Reader.Proto).
		Merge(NativeMetaProto.Closer.Proto).SetPrototype(Proto.Native)

	NativeMetaProto.WriteCloser.Proto.
		Merge(NativeMetaProto.Writer.Proto).
		Merge(NativeMetaProto.Closer.Proto).SetPrototype(Proto.Native)

	NativeMetaProto.ReadWriteCloser.Proto.
		Merge(NativeMetaProto.ReadWriter.Proto).
		Merge(NativeMetaProto.Closer.Proto).SetPrototype(Proto.Native)
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
		if rb := Value(m).Object().Get(Str("read2")); rb.IsObject() {
			t := rb.Object().Call(nil, Bytes(p)).AssertShape("(i, EN)", "Reader.read2").Native()
			if IsError(t.Get(1)) {
				return t.Get(0).Int(), ToErrorRootCause(t.Get(1)).(error)
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
			if IsError(v) {
				return 0, ToError(v)
			}
			return v.AssertNumber("Writer.write").Int(), nil
		}
		if rb := Value(m).Object().Get(Str("write2")); rb.IsObject() {
			t := rb.Object().Call(nil, Bytes(p)).AssertShape("(i, EN)", "Writer.write2").Native()
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
		if rb := Value(m).Object().Get(Str("close")); rb.IsObject() {
			if v := rb.Object().Call(nil); IsError(v) {
				return ToError(v)
			}
			return nil
		}
	}
	return fmt.Errorf("closer not implemented")
}

var ioReadlinesIter *NativeMeta

type ioReadlinesStruct struct {
	rd    *bufio.Reader
	delim byte
	bytes bool
}
