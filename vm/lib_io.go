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
		s := env.Stack()
		for i := 0; i < s.Size(); i++ {
			f.WriteString(vtoString(s.Get(i)))
		}

		return base.NewValue()
	}
}

func stdPrintln(f *os.File, ex bool) func(env *base.Env) base.Value {
	return func(env *base.Env) base.Value {
		s := env.Stack()
		for i := 0; i < s.Size(); i++ {
			f.WriteString(vtoString(s.Get(i)) + " ")
		}
		f.WriteString("\n")
		return base.NewValue()
	}
}

func stdWrite(f *os.File, ex bool) func(env *base.Env) base.Value {
	return func(env *base.Env) base.Value {
		s := env.Stack()
		for i := 0; i < s.Size(); i++ {
			f.Write(s.Get(i).Bytes())
		}
		return base.NewValue()
	}
}

var lib_outprint = LibFunc{name: "out_print", args: 0, ff: stdPrint(os.Stdout, true)}
var lib_outprintln = LibFunc{name: "out_println", args: 0, ff: stdPrintln(os.Stdout, true)}
var lib_outwrite = LibFunc{name: "out_write", args: 0, ff: stdWrite(os.Stdout, true)}
var lib_errprint = LibFunc{name: "err_print", args: 0, ff: stdPrint(os.Stderr, true)}
var lib_errprintln = LibFunc{name: "err_println", args: 0, ff: stdPrintln(os.Stderr, true)}
var lib_errwrite = LibFunc{name: "err_write", args: 0, ff: stdWrite(os.Stderr, true)}
