package potatolang

import (
	"fmt"
	"os"
	"unsafe"
)

func stdPrint(f *os.File) func(env *Env) Value {
	return func(env *Env) Value {
		for i := 0; i < env.LocalSize(); i++ {
			f.WriteString(env.LocalGet(i).String())
		}

		return Value{}
	}
}

func _sprintf(env *Env) string {
	msg := env.LocalGet(0).MustString()
	for i := range msg {
		if msg[i] == '{' && i < len(msg)-1 && msg[i+1] == '}' {
			msg[i] = '%'
			msg[i+1] = 's'
		}
	}

	args := []interface{}{}
	for i := 1; i < env.LocalSize(); i++ {
		args = append(args, env.LocalGet(i))
	}

	return fmt.Sprintf(*(*string)(unsafe.Pointer(&msg)), args...)
}

func stdPrintf(f *os.File) func(env *Env) Value {
	return func(env *Env) Value {
		f.WriteString(_sprintf(env))
		return Value{}
	}
}

func stdSprintf(env *Env) Value {
	return NewStringValueString(_sprintf(env))
}

func stdPrintln(f *os.File) func(env *Env) Value {
	return func(env *Env) Value {
		for i := 0; i < env.LocalSize(); i++ {
			f.WriteString(env.LocalGet(i).String() + " ")
		}
		f.WriteString("\n")
		return Value{}
	}
}

func stdWrite(f *os.File) func(env *Env) Value {
	return func(env *Env) Value {
		for i := 0; i < env.LocalSize(); i++ {
			switch a := env.LocalGet(i); a.Type() {
			case StringType:
				f.Write(env.LocalGet(i).AsString())
			case SliceType:
				m := a.AsSlice()
				buf := make([]byte, len(m.l))
				for i, b := range m.l {
					buf[i] = byte(b.MustNumber())
				}
				f.Write(buf)
			default:
				panicf("stdWrite can't write: %+v", a)
			}
		}
		return Value{}
	}
}

func initIOLib() {
	lio := NewStruct()
	lio.Put("println", NewNativeValue(0, stdPrintln(os.Stdout)))
	lio.Put("print", NewNativeValue(0, stdPrint(os.Stdout)))
	lio.Put("printf", NewNativeValue(1, stdPrintf(os.Stdout)))
	lio.Put("write", NewNativeValue(0, stdWrite(os.Stdout)))
	lio.Put("err", NewStructValue(NewStruct().
		Put("println", NewNativeValue(0, stdPrintln(os.Stderr))).
		Put("print", NewNativeValue(0, stdPrint(os.Stderr))).
		Put("write", NewNativeValue(0, stdWrite(os.Stderr)))))

	CoreLibs["io"] = NewStructValue(lio)
}
