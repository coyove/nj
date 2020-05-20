package potatolang

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
)

func fmtPrint(flag byte) func(env *Env) Value {
	return func(env *Env) Value {
		args := make([]interface{}, env.LocalSize())
		for i := range args {
			args[i] = env.LocalGet(i).AsInterface()
		}
		var n int
		var err error

		switch flag {
		case 'l':
			n, err = fmt.Println(args...)
		case 'f':
			n, err = fmt.Printf(args[0].(string), args[1:]...)
		default:
			n, err = fmt.Print(args...)
		}

		if err != nil {
			env.B = NewStringValue(err.Error())
		}
		return NewNumberValue(float64(n))
	}
}

func fmtSprint(flag byte) func(env *Env) Value {
	return func(env *Env) Value {
		args := make([]interface{}, env.LocalSize())
		for i := range args {
			args[i] = env.LocalGet(i).AsInterface()
		}
		var n string
		switch flag {
		case 'l':
			n = fmt.Sprintln(args...)
		case 'f':
			n = fmt.Sprintf(args[0].(string), args[1:]...)
		default:
			n = fmt.Sprint(args...)
		}
		return NewStringValue(n)
	}
}

func fmtFprint(flag byte) func(env *Env) Value {
	return func(env *Env) Value {
		args := make([]interface{}, env.LocalSize())
		for i := range args {
			args[i] = env.LocalGet(i).AsInterface()
		}
		var n int
		var err error
		switch flag {
		case 'l':
			n, err = fmt.Fprintln(args[0].(io.Writer), args[1:]...)
		case 'f':
			n, err = fmt.Fprintf(args[0].(io.Writer), args[1].(string), args[2:]...)
		default:
			n, err = fmt.Fprint(args[0].(io.Writer), args[1:]...)
		}

		if err != nil {
			env.B = NewStringValue(err.Error())
		}
		return NewNumberValue(float64(n))
	}
}

func fmtScan(flag string) func(env *Env) Value {
	return func(env *Env) Value {
		var start int
		switch flag {
		case "scanf", "sscanln", "sscan", "fscan", "fscanln":
			start = 1
		case "sscanf", "fscanf":
			start = 2
		}

		receivers := make([]interface{}, env.LocalSize())
		for i := start; i < len(receivers); i++ {
			switch LoadPointerUnsafe(env.LocalGet(i)).Type() {
			case StringType:
				receivers[i] = new(string)
			case NumberType:
				receivers[i] = new(float64)
			default:
				panicf("Scan: only string and number are supported")
			}
		}

		var n int
		var err error

		switch flag {
		case "scanln":
			n, err = fmt.Scanln(receivers...)
		case "scanf":
			n, err = fmt.Scanf(string(env.LocalGet(0).MustString()), receivers[1:]...)
		case "scan":
			n, err = fmt.Scan(receivers...)
		case "sscanln":
			n, err = fmt.Sscanln(string(env.LocalGet(0).MustString()), receivers[1:]...)
		case "sscanf":
			n, err = fmt.Sscanf(string(env.LocalGet(0).MustString()), string(env.LocalGet(1).MustString()), receivers[2:]...)
		case "sscan":
			n, err = fmt.Sscan(string(env.LocalGet(0).MustString()), receivers[1:]...)
		case "fscan":
			n, err = fmt.Fscan(env.LocalGet(0).AsInterface().(io.Reader), receivers[1:]...)
		case "fscanln":
			n, err = fmt.Fscanln(env.LocalGet(0).AsInterface().(io.Reader), receivers[1:]...)
		case "fscanf":
			n, err = fmt.Fscanf(env.LocalGet(0).AsInterface().(io.Reader), string(env.LocalGet(1).MustString()), receivers[1:]...)
		}

		if err == nil {
			for i := start; i < len(receivers); i++ {
				switch v := receivers[i].(type) {
				case *float64:
					StorePointerUnsafe(env.LocalGet(i), NewNumberValue(*v))
				case *string:
					StorePointerUnsafe(env.LocalGet(i), NewStringValue(*v))
				}
			}
		}

		if err != nil {
			env.B = NewStringValue(err.Error())
		}
		return NewNumberValue(float64(n))
	}
}

func fmtWrite(env *Env) Value {
	var n int
	var err error
	f := env.LocalGet(0).AsInterface().(io.Writer)

	for i := 1; i < env.LocalSize(); i++ {
		switch a := env.LocalGet(i); a.Type() {
		case StringType:
			n, err = f.Write([]byte(env.LocalGet(i).AsString()))
		case SliceType:
			m := a.AsSlice()
			buf := make([]byte, len(m.l))
			for i, b := range m.l {
				buf[i] = byte(b.MustNumber())
			}
			n, err = f.Write(buf)
		default:
			panicf("stdWrite can't write: %+v", a)
		}
	}

	if err != nil {
		env.B = NewStringValue(err.Error())
	}
	return NewNumberValue(float64(n))
}

func walkObject(buf []byte) Value {
	m := NewStruct()
	jsonparser.ObjectEach(buf, func(key, value []byte, dataType jsonparser.ValueType, offset int) error {
		switch dataType {
		case jsonparser.Unknown:
			panic(value)
		case jsonparser.Null:
			m.Put(string(key), Value{})
		case jsonparser.Boolean:
			b, err := jsonparser.ParseBoolean(value)
			panicerr(err)
			m.Put(string(key), NewBoolValue(b))
		case jsonparser.Number:
			num, err := jsonparser.ParseFloat(value)
			panicerr(err)
			m.Put(string(key), NewNumberValue(num))
		case jsonparser.String:
			str, err := jsonparser.ParseString(value)
			panicerr(err)
			m.Put(string(key), NewStringValue(str))
		case jsonparser.Array:
			m.Put(string(key), walkArray(value))
		case jsonparser.Object:
			m.Put(string(key), walkObject(value))
		}
		return nil
	})
	return NewStructValue(m)
}

func walkArray(buf []byte) Value {
	m := NewSlice()
	i := 0
	jsonparser.ArrayEach(buf, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		switch dataType {
		case jsonparser.Unknown:
			panic(value)
		case jsonparser.Null:
			m.Put(i, Value{})
		case jsonparser.Boolean:
			b, err := jsonparser.ParseBoolean(value)
			panicerr(err)
			m.Put(i, NewBoolValue(b))
		case jsonparser.Number:
			num, err := jsonparser.ParseFloat(value)
			panicerr(err)
			m.Put(i, NewNumberValue(num))
		case jsonparser.String:
			str, err := jsonparser.ParseString(value)
			panicerr(err)
			m.Put(i, NewStringValue(str))
		case jsonparser.Array:
			m.Put(i, walkArray(value))
		case jsonparser.Object:
			m.Put(i, walkObject(value))
		}
		i++
	})
	return NewSliceValue(m)
}

func jsonUnmarshal(env *Env) Value {
	json := []byte(strings.TrimSpace(env.LocalGet(0).MustString()))
	if len(json) == 0 {
		return Value{}
	}
	switch json[0] {
	case '[':
		return walkArray(json)
	case '{':
		return walkObject(json)
	case '"':
		str, err := jsonparser.ParseString(json)
		if err != nil {
			StorePointerUnsafe(env.LocalGet(1), NewStringValue(err.Error()))
		}
		return NewStringValue(str)
	case 't', 'f':
		b, err := jsonparser.ParseBoolean(json)
		if err != nil {
			StorePointerUnsafe(env.LocalGet(1), NewStringValue(err.Error()))
		}
		return NewBoolValue(b)
	default:
		num, err := jsonparser.ParseFloat(json)
		if err != nil {
			StorePointerUnsafe(env.LocalGet(1), NewStringValue(err.Error()))
		}
		return NewNumberValue(num)
	}
}

func strconvFormatFloat(env *Env) Value {
	v := env.LocalGet(0).MustNumber()
	base := byte(env.LocalGet(1).MustNumber())
	digits := int(env.LocalGet(2).MustNumber())
	return NewStringValue(strconv.FormatFloat(v, byte(base), digits, 64))
}

func strconvFormatInt(env *Env) Value {
	return NewStringValue(strconv.FormatInt(int64(env.LocalGet(0).MustNumber()), int(env.LocalGet(1).MustNumber())))
}

func strconvParseFloat(env *Env) Value {
	v, err := strconv.ParseFloat(string(env.LocalGet(0).MustString()), 64)
	if err != nil {
		StorePointerUnsafe(env.LocalGet(1), NewStringValue(err.Error()))
	}
	return NewNumberValue(v)
}

func stringsIndex(env *Env) Value {
	return NewNumberValue(float64(strings.Index(env.LocalGet(0).MustString(), env.LocalGet(1).MustString())))
}

func initLibAux() {
	lfmt := NewStruct()
	lfmt.Put("Println", NewNativeValue(0, fmtPrint('l')))
	lfmt.Put("Print", NewNativeValue(0, fmtPrint(0)))
	lfmt.Put("Printf", NewNativeValue(1, fmtPrint('f')))
	lfmt.Put("Sprintln", NewNativeValue(0, fmtSprint('l')))
	lfmt.Put("Sprint", NewNativeValue(0, fmtSprint(0)))
	lfmt.Put("Sprintf", NewNativeValue(1, fmtSprint('f')))
	lfmt.Put("Fprintln", NewNativeValue(1, fmtFprint('l')))
	lfmt.Put("Fprint", NewNativeValue(1, fmtFprint(0)))
	lfmt.Put("Fprintf", NewNativeValue(2, fmtFprint('f')))
	lfmt.Put("Scanln", NewNativeValue(1, fmtScan("scanln")))
	lfmt.Put("Scan", NewNativeValue(1, fmtScan("scan")))
	lfmt.Put("Scanf", NewNativeValue(1, fmtScan("scanf")))
	lfmt.Put("Sscanln", NewNativeValue(1, fmtScan("sscanln")))
	lfmt.Put("Sscan", NewNativeValue(1, fmtScan("sscan")))
	lfmt.Put("Sscanf", NewNativeValue(2, fmtScan("sscanf")))
	lfmt.Put("Fscanln", NewNativeValue(1, fmtScan("fscanln")))
	lfmt.Put("Fscan", NewNativeValue(1, fmtScan("fscan")))
	lfmt.Put("Fscanf", NewNativeValue(2, fmtScan("fscanf")))
	lfmt.Put("Write", NewNativeValue(0, fmtWrite))
	CoreLibs["fmt"] = NewStructValue(lfmt)

	los := NewStruct()
	los.Put("Stdout", NewInterfaceValue(os.Stdout))
	los.Put("Stdin", NewInterfaceValue(os.Stdin))
	los.Put("Stderr", NewInterfaceValue(os.Stderr))
	CoreLibs["os"] = NewStructValue(los)

	ljson := NewStruct()
	ljson.Put("Unmarshal", NewNativeValue(1, jsonUnmarshal))
	ljson.Put("Marshal", NewNativeValue(1, func(env *Env) Value { return NewStringValue(env.LocalGet(0).toString(0, true)) }))
	CoreLibs["json"] = NewStructValue(ljson)

	lstrconv := NewStruct()
	lstrconv.Put("FormatFloat", NewNativeValue(3, strconvFormatFloat))
	lstrconv.Put("ParseFloat", NewNativeValue(1, strconvParseFloat))
	lstrconv.Put("FormatInt", NewNativeValue(2, strconvFormatInt))
	CoreLibs["strconv"] = NewStructValue(lstrconv)

	lstrings := NewStruct()
	lstrings.Put("Index", NewNativeValue(2, stringsIndex))
	CoreLibs["strings"] = NewStructValue(lstrconv)
}
