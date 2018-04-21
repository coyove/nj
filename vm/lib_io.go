package vm

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/coyove/bracket/base"
)

func vtoString(v base.Value, lv int) string {
	if lv > 32 {
		return "<omit deep nesting>"
	}

	switch v.Type() {
	case base.Tbool:
		return strconv.FormatBool(v.AsBool())
	case base.Tnumber:
		n := v.AsNumber()
		if float64(int64(n)) == n {
			return strconv.FormatInt(int64(n), 10)
		}
		return strconv.FormatFloat(n, 'f', 9, 64)
	case base.Tstring:
		return strconv.Quote(v.AsString())
	case base.Tlist:
		arr := v.AsList()
		buf := &bytes.Buffer{}
		buf.WriteString("[")
		for _, v := range arr {
			buf.WriteString(vtoString(v, lv+1))
			buf.WriteString(",")
		}
		if len(arr) > 0 {
			buf.Truncate(buf.Len() - 1)
		}
		buf.WriteString("]")
		return buf.String()
	case base.Tmap:
		m := v.AsMap()
		buf := &bytes.Buffer{}
		buf.WriteString("{")
		for k, v := range m {
			buf.WriteString(k)
			buf.WriteString(":")
			buf.WriteString(vtoString(v, lv+1))
			buf.WriteString(",")
		}
		if len(m) > 0 {
			buf.Truncate(buf.Len() - 1)
		}
		buf.WriteString("}")
		return buf.String()
	case base.Tbytes:
		arr := v.AsBytes()
		buf := &bytes.Buffer{}
		buf.WriteString("[")
		for _, v := range arr {
			buf.WriteString(fmt.Sprintf("%02x", int(v)))
			buf.WriteString(",")
		}
		if len(arr) > 0 {
			buf.Truncate(buf.Len() - 1)
		}
		buf.WriteString("]")
		return buf.String()
	case base.Tclosure:
		return "<" + v.AsClosure().String() + ">"
	}
	return "nil"
}

func stdPrint(f *os.File, ex bool) func(env *base.Env) base.Value {
	return func(env *base.Env) base.Value {
		s := env.Stack()
		for i := 0; i < s.Size(); i++ {
			f.WriteString(vtoString(s.Get(i), 0))
		}

		return base.NewValue()
	}
}

func stdPrintln(f *os.File, ex bool) func(env *base.Env) base.Value {
	return func(env *base.Env) base.Value {
		s := env.Stack()
		for i := 0; i < s.Size(); i++ {
			f.WriteString(vtoString(s.Get(i), 0) + " ")
		}
		f.WriteString("\n")
		return base.NewValue()
	}
}

func stdWrite(f *os.File, ex bool) func(env *base.Env) base.Value {
	return func(env *base.Env) base.Value {
		s := env.Stack()
		for i := 0; i < s.Size(); i++ {
			f.Write(s.Get(i).AsBytes())
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
