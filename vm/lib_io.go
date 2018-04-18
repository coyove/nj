package vm

import (
	"bytes"
	"os"
	"strconv"

	"github.com/coyove/bracket/base"
)

func vtoString(v base.Value) string {
	switch v.Type() {
	case base.TY_bool:
		return strconv.FormatBool(v.Bool())
	case base.TY_number:
		n := v.Number()
		if float64(int64(n)) == n {
			return strconv.FormatInt(int64(n), 10)
		}
		return strconv.FormatFloat(n, 'f', 9, 64)
	case base.TY_string:
		return "string (" + v.String() + ")"
	case base.TY_array:
		arr := v.Array()
		buf := &bytes.Buffer{}
		buf.WriteString("list (")
		for _, v := range arr {
			buf.WriteString(vtoString(v))
			buf.WriteString(" ")
		}
		buf.WriteString(")")
		return buf.String()
	case base.TY_map:
		m := v.Map()
		buf := &bytes.Buffer{}
		buf.WriteString("map (")
		for k, v := range m {
			buf.WriteString(k)
			buf.WriteString(":")
			buf.WriteString(vtoString(v))
			buf.WriteString(" ")
		}
		buf.WriteString(")")
		return buf.String()
	case base.TY_bytes:
		arr := v.Bytes()
		buf := &bytes.Buffer{}
		buf.WriteString("bytes (")
		for _, v := range arr {
			buf.WriteString(strconv.Itoa(int(v)))
			buf.WriteString(" ")
		}
		buf.WriteString(")")
		return buf.String()
	case base.TY_closure:
		return v.Closure().String()
	}
	return "nil"
}

func stdPrint(f *os.File, ex bool) func(env *base.Env) base.Value {
	return func(env *base.Env) base.Value {
		if ex {
			s := env.Stack()
			for i := 0; i < s.Size(); i++ {
				f.WriteString(vtoString(s.Get(i)))
			}
		} else {
			for i := 0; i < env.SizeR(); i++ {
				arg := env.R(i)
				f.WriteString(vtoString(arg))
			}
		}

		return base.NewValue()
	}
}

func stdPrintln(f *os.File, ex bool) func(env *base.Env) base.Value {
	return func(env *base.Env) base.Value {
		if ex {
			s := env.Stack()
			for i := 0; i < s.Size(); i++ {
				f.WriteString(vtoString(s.Get(i)) + " ")
			}
		} else {
			for i := 0; i < env.SizeR(); i++ {
				arg := env.R(i)
				f.WriteString(vtoString(arg) + " ")
			}
		}
		f.WriteString("\n")
		return base.NewValue()
	}
}

func stdWrite(f *os.File, ex bool) func(env *base.Env) base.Value {
	return func(env *base.Env) base.Value {
		if ex {
			s := env.Stack()
			for i := 0; i < s.Size(); i++ {
				f.Write(s.Get(i).Bytes())
			}
		} else {
			for i := 0; i < env.SizeR(); i++ {
				arg := env.R(i)
				f.Write(arg.Bytes())
			}
		}
		return base.NewValue()
	}
}

var lib_outprint = LibFunc{name: "stdout/print", args: 0, f: stdPrint(os.Stdout, false), ff: stdPrint(os.Stdout, true)}
var lib_outprintln = LibFunc{name: "stdout/println", args: 0, f: stdPrintln(os.Stdout, false), ff: stdPrintln(os.Stdout, true)}
var lib_outwrite = LibFunc{name: "stdout/write", args: 0, f: stdWrite(os.Stdout, false), ff: stdWrite(os.Stdout, true)}
var lib_errprint = LibFunc{name: "stderr/print", args: 0, f: stdPrint(os.Stderr, false), ff: stdPrint(os.Stderr, true)}
var lib_errprintln = LibFunc{name: "stderr/println", args: 0, f: stdPrintln(os.Stderr, false), ff: stdPrintln(os.Stderr, true)}
var lib_errwrite = LibFunc{name: "stderr/write", args: 0, f: stdWrite(os.Stderr, false), ff: stdWrite(os.Stderr, true)}
