package potatolang

import (
	"fmt"
	"strings"
)

func fmtPrint(flag byte) func(env *Env) {
	return func(env *Env) {
		args := make([]interface{}, len(env.Vararg))
		for i := range args {
			args[i] = env.Vararg[i]
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
			env.Vararg = []Value{Str(err.Error())}
		} else {
			env.A = Num(float64(n))
		}
	}
}

func fmtSprint(flag byte) func(env *Env) {
	return func(env *Env) {
		args := make([]interface{}, len(env.Vararg))
		for i := range args {
			args[i] = env.Vararg[i].Any()
		}
		var n string
		switch flag {
		case 'l':
			n = fmt.Sprintln(args...)
		case 'f':
			n = fmt.Sprintf(env.Get(0).Expect(STR).Str(), args...)
		default:
			n = fmt.Sprint(args...)
		}
		env.A = Str(n)
	}
}

// func fmtFprint(flag byte) func(env *Env) Value {
// 	return func(env *Env) Value {
// 		args := make([]interface{}, env.Size())
// 		for i := range args {
// 			args[i] = env.Get(i).AsInterface()
// 		}
// 		var n int
// 		var err error
// 		switch flag {
// 		case 'l':
// 			n, err = fmt.Fprintln(args[0].(io.Writer), args[1:]...)
// 		case 'f':
// 			n, err = fmt.Fprintf(args[0].(io.Writer), args[1].(string), args[2:]...)
// 		default:
// 			n, err = fmt.Fprint(args[0].(io.Writer), args[1:]...)
// 		}
//
// 		if err != nil {
// 			env.B = Str(err.Error())
// 		}
// 		return Num(float64(n))
// 	}
// }
//
// func fmtScan(flag string) func(env *Env) Value {
// 	return func(env *Env) Value {
// 		var start int
// 		switch flag {
// 		case "scanf", "sscanln", "sscan", "fscan", "fscanln":
// 			start = 1
// 		case "sscanf", "fscanf":
// 			start = 2
// 		}
//
// 		receivers := make([]interface{}, env.Size())
// 		for i := start; i < len(receivers); i++ {
// 			switch LoadPointerUnsafe(env.Get(i)).Type() {
// 			case STR:
// 				receivers[i] = new(string)
// 			case NUM:
// 				receivers[i] = new(float64)
// 			default:
// 				panicf("Scan: only string and number are supported")
// 			}
// 		}
//
// 		var n int
// 		var err error
//
// 		switch flag {
// 		case "scanln":
// 			n, err = fmt.Scanln(receivers...)
// 		case "scanf":
// 			n, err = fmt.Scanf(string(env.Get(0).MustString()), receivers[1:]...)
// 		case "scan":
// 			n, err = fmt.Scan(receivers...)
// 		case "sscanln":
// 			n, err = fmt.Sscanln(string(env.Get(0).MustString()), receivers[1:]...)
// 		case "sscanf":
// 			n, err = fmt.Sscanf(string(env.Get(0).MustString()), string(env.Get(1).MustString()), receivers[2:]...)
// 		case "sscan":
// 			n, err = fmt.Sscan(string(env.Get(0).MustString()), receivers[1:]...)
// 		case "fscan":
// 			n, err = fmt.Fscan(env.Get(0).AsInterface().(io.Reader), receivers[1:]...)
// 		case "fscanln":
// 			n, err = fmt.Fscanln(env.Get(0).AsInterface().(io.Reader), receivers[1:]...)
// 		case "fscanf":
// 			n, err = fmt.Fscanf(env.Get(0).AsInterface().(io.Reader), string(env.Get(1).MustString()), receivers[1:]...)
// 		}
//
// 		if err == nil {
// 			for i := start; i < len(receivers); i++ {
// 				switch v := receivers[i].(type) {
// 				case *float64:
// 					StorePointerUnsafe(env.Get(i), Num(*v))
// 				case *string:
// 					StorePointerUnsafe(env.Get(i), Str(*v))
// 				}
// 			}
// 		}
//
// 		if err != nil {
// 			env.B = Str(err.Error())
// 		}
// 		return Num(float64(n))
// 	}
// }
//
// func fmtWrite(env *Env) Value {
// 	var n int
// 	var err error
// 	f := env.Get(0).AsInterface().(io.Writer)
//
// 	for i := 1; i < env.Size(); i++ {
// 		switch a := env.Get(i); a.Type() {
// 		case STR:
// 			n, err = f.Write([]byte(env.Get(i).Str()))
// 		default:
// 			panicf("stdWrite can't write: %+v", a)
// 		}
// 	}
//
// 	if err != nil {
// 		env.B = Str(err.Error())
// 	}
// 	return Num(float64(n))
// }
//
// func walkObject(buf []byte) Value {
// 	m := Table{}
// 	jsonparser.ObjectEach(buf, func(key, value []byte, dataType jsonparser.ValueType, offset int) error {
// 		switch dataType {
// 		case jsonparser.Unknown:
// 			panic(value)
// 		case jsonparser.Null:
// 			m.Put(string(key), Value{})
// 		case jsonparser.Boolean:
// 			b, err := jsonparser.ParseBoolean(value)
// 			panicerr(err)
// 			m.Put(string(key), Bln(b))
// 		case jsonparser.Number:
// 			num, err := jsonparser.ParseFloat(value)
// 			panicerr(err)
// 			m.Put(string(key), Num(num))
// 		case jsonparser.String:
// 			str, err := jsonparser.ParseString(value)
// 			panicerr(err)
// 			m.Put(string(key), Str(str))
// 		case jsonparser.Array:
// 			m.Put(string(key), walkArray(value))
// 		case jsonparser.Object:
// 			m.Put(string(key), walkObject(value))
// 		}
// 		return nil
// 	})
// 	return NewStructValue(m)
// }
//
// func walkArray(buf []byte) Value {
// 	m := NewSlice()
// 	i := 0
// 	jsonparser.ArrayEach(buf, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
// 		switch dataType {
// 		case jsonparser.Unknown:
// 			panic(value)
// 		case jsonparser.Null:
// 			m.Put(i, Value{})
// 		case jsonparser.Boolean:
// 			b, err := jsonparser.ParseBoolean(value)
// 			panicerr(err)
// 			m.Put(i, Bln(b))
// 		case jsonparser.Number:
// 			num, err := jsonparser.ParseFloat(value)
// 			panicerr(err)
// 			m.Put(i, Num(num))
// 		case jsonparser.String:
// 			str, err := jsonparser.ParseString(value)
// 			panicerr(err)
// 			m.Put(i, Str(str))
// 		case jsonparser.Array:
// 			m.Put(i, walkArray(value))
// 		case jsonparser.Object:
// 			m.Put(i, walkObject(value))
// 		}
// 		i++
// 	})
// 	return NewSliceValue(m)
// }
//
// func jsonUnmarshal(env *Env) Value {
// 	json := []byte(strings.TrimSpace(env.Get(0).MustString()))
// 	if len(json) == 0 {
// 		return Value{}
// 	}
// 	switch json[0] {
// 	case '[':
// 		return walkArray(json)
// 	case '{':
// 		return walkObject(json)
// 	case '"':
// 		str, err := jsonparser.ParseString(json)
// 		if err != nil {
// 			StorePointerUnsafe(env.Get(1), Str(err.Error()))
// 		}
// 		return Str(str)
// 	case 't', 'f':
// 		b, err := jsonparser.ParseBoolean(json)
// 		if err != nil {
// 			StorePointerUnsafe(env.Get(1), Str(err.Error()))
// 		}
// 		return Bln(b)
// 	default:
// 		num, err := jsonparser.ParseFloat(json)
// 		if err != nil {
// 			StorePointerUnsafe(env.Get(1), Str(err.Error()))
// 		}
// 		return Num(num)
// 	}
// }
//
// func strconvFormatFloat(env *Env) Value {
// 	v := env.Get(0).MustNumber()
// 	base := byte(env.Get(1).MustNumber())
// 	digits := int(env.Get(2).MustNumber())
// 	return Str(strconv.FormatFloat(v, byte(base), digits, 64))
// }
//
// func strconvFormatInt(env *Env) Value {
// 	return Str(strconv.FormatInt(int64(env.Get(0).MustNumber()), int(env.Get(1).MustNumber())))
// }
//
// func strconvParseFloat(env *Env) Value {
// 	v, err := strconv.ParseFloat(string(env.Get(0).MustString()), 64)
// 	if err != nil {
// 		StorePointerUnsafe(env.Get(1), Str(err.Error()))
// 	}
// 	return Num(v)
// }
//
// func stringsIndex(env *Env) Value {
// 	return Num(float64(strings.Index(env.Get(0).MustString(), env.Get(1).MustString())))
// }

func initLibAux() {
	// 	lfmt.Put("Printf", NativeFun(1, fmtPrint('f')))
	// 	lfmt.Put("Sprintln", NativeFun(0, fmtSprint('l')))
	// 	lfmt.Put("Sprint", NativeFun(0, fmtSprint(0)))
	// 	lfmt.Put("Sprintf", NativeFun(1, fmtSprint('f')))
	// 	lfmt.Put("Fprintln", NativeFun(1, fmtFprint('l')))
	// 	lfmt.Put("Fprint", NativeFun(1, fmtFprint(0)))
	// 	lfmt.Put("Fprintf", NativeFun(2, fmtFprint('f')))
	// 	lfmt.Put("Scanln", NativeFun(1, fmtScan("scanln")))
	// 	lfmt.Put("Scan", NativeFun(1, fmtScan("scan")))
	// 	lfmt.Put("Scanf", NativeFun(1, fmtScan("scanf")))
	// 	lfmt.Put("Sscanln", NativeFun(1, fmtScan("sscanln")))
	// 	lfmt.Put("Sscan", NativeFun(1, fmtScan("sscan")))
	// 	lfmt.Put("Sscanf", NativeFun(2, fmtScan("sscanf")))
	// 	lfmt.Put("Fscanln", NativeFun(1, fmtScan("fscanln")))
	// 	lfmt.Put("Fscan", NativeFun(1, fmtScan("fscan")))
	// 	lfmt.Put("Fscanf", NativeFun(2, fmtScan("fscanf")))
	// 	lfmt.Put("Write", NativeFun(0, fmtWrite))
	G.Puts("print", NativeFun(0, true, fmtPrint('l')), false)
	//
	// 	los := NewStruct()
	// 	los.Put("Stdout", NewInterfaceValue(os.Stdout))
	// 	los.Put("Stdin", NewInterfaceValue(os.Stdin))
	// 	los.Put("Stderr", NewInterfaceValue(os.Stderr))
	// 	G["os"] = NewStructValue(los)
	//
	// 	ljson := NewStruct()
	// 	ljson.Put("Unmarshal", NativeFun(1, jsonUnmarshal))
	// 	ljson.Put("Marshal", NativeFun(1, func(env *Env) Value { return Str(env.Get(0).toString(0, true)) }))
	// 	G["json"] = NewStructValue(ljson)
	//
	// 	lstrconv := NewStruct()
	// 	lstrconv.Put("FormatFloat", NativeFun(3, strconvFormatFloat))
	// 	lstrconv.Put("ParseFloat", NativeFun(1, strconvParseFloat))
	// 	lstrconv.Put("FormatInt", NativeFun(2, strconvFormatInt))
	// 	G["strconv"] = NewStructValue(lstrconv)
	//
	lstring := &Table{}
	lstring.Puts("format", NativeFun(1, true, fmtSprint('f')), false)
	lstring.Puts("rep", NativeFun(2, false, func(env *Env) {
		env.A = Str(strings.Repeat(env.In(0, STR).Str(), int(env.In(1, NUM).Num())))
	}), false)
	lstring.Puts("char", NativeFun(1, false, func(env *Env) {
		env.A = Str(string(rune(env.In(0, NUM).Num())))
	}), false)
	G.Puts("string", Tab(lstring), false)
}
