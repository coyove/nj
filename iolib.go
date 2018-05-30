package potatolang

import (
	"log"
	"os"
)

func stdPrint(f *os.File) func(env *Env) Value {
	return func(env *Env) Value {
		for i := 0; i < env.SSize(); i++ {
			f.WriteString(env.SGet(i).ToPrintString())
		}

		return NewValue()
	}
}

func stdPrintln(f *os.File) func(env *Env) Value {
	return func(env *Env) Value {
		for i := 0; i < env.SSize(); i++ {
			f.WriteString(env.SGet(i).ToPrintString() + " ")
		}
		f.WriteString("\n")
		return NewValue()
	}
}

func stdWrite(f *os.File) func(env *Env) Value {
	return func(env *Env) Value {
		for i := 0; i < env.SSize(); i++ {
			switch a := env.SGet(i); a.ty {
			case Tstring:
				f.WriteString(env.SGet(i).AsString())
			case Tmap:
				buf := make([]byte, 1)
				for _, b := range a.AsMap().l {
					buf[0] = byte(b.AsNumber())
					f.Write(buf)
				}
			default:
				log.Panicf("stdWrite can't write: %+v", a)
			}
		}
		return NewValue()
	}
}

func initIOLib() {
	lio := NewMap()
	lio.Puts("println", NewNativeValue(0, stdPrintln(os.Stdout)))
	lio.Puts("print", NewNativeValue(0, stdPrint(os.Stdout)))
	lio.Puts("write", NewNativeValue(0, stdWrite(os.Stdout)))
	lio.Puts("err", NewMapValue(NewMap().
		Puts("println", NewNativeValue(0, stdPrintln(os.Stderr))).
		Puts("print", NewNativeValue(0, stdPrint(os.Stderr))).
		Puts("write", NewNativeValue(0, stdWrite(os.Stderr)))))

	CoreLibs["io"] = NewMapValue(lio)
}
