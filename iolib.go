package potatolang

import (
	"bytes"
	"os"
	"strconv"
)

func stdPrint(f *os.File) func(env *Env) Value {
	return func(env *Env) Value {
		for i := 0; i < env.SSize(); i++ {
			f.WriteString(env.SGet(i).ToPrintString())
		}

		return NewValue()
	}
}

func _sprintf(env *Env) string {
	msg := []rune(env.SGet(0).Str())
	buf, numbuf := bytes.Buffer{}, bytes.Buffer{}
	i := 0
	for i < len(msg) {
		if msg[i] != '~' {
			buf.WriteRune(msg[i])
			i++
			continue
		}
		if i+1 >= len(msg) {
			break
		}
		if msg[i+1] == '~' {
			buf.WriteRune('~')
			i += 2
			continue
		}
		numbuf.Reset()
		j := i + 1
		for ; j < len(msg); j++ {
			if msg[j] >= '0' && msg[j] <= '9' {
				numbuf.WriteRune(msg[j])
			} else {
				break
			}
		}
		if j == i+1 {
			i++
			continue
		}
		i = j
		num, _ := strconv.Atoi(numbuf.String())
		buf.WriteString(env.SGet(num).ToPrintString())
	}

	return buf.String()
}

func stdPrintf(f *os.File) func(env *Env) Value {
	return func(env *Env) Value {
		f.WriteString(_sprintf(env))
		return Value{}
	}
}

func stdSprintf(env *Env) Value {
	return NewStringValue(_sprintf(env))
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
					buf[0] = byte(b.Num())
					f.Write(buf)
				}
			default:
				panicf("stdWrite can't write: %+v", a)
			}
		}
		return NewValue()
	}
}

func initIOLib() {
	lio := NewMap()
	lio.Puts("println", NewNativeValue(0, stdPrintln(os.Stdout)))
	lio.Puts("print", NewNativeValue(0, stdPrint(os.Stdout)))
	lio.Puts("printf", NewNativeValue(1, stdPrintf(os.Stdout)))
	lio.Puts("sprintf", NewNativeValue(1, stdSprintf))
	lio.Puts("write", NewNativeValue(0, stdWrite(os.Stdout)))
	lio.Puts("err", NewMapValue(NewMap().
		Puts("println", NewNativeValue(0, stdPrintln(os.Stderr))).
		Puts("print", NewNativeValue(0, stdPrint(os.Stderr))).
		Puts("write", NewNativeValue(0, stdWrite(os.Stderr)))))

	CoreLibs["io"] = NewMapValue(lio)
}
