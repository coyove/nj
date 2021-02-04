package script

import (
	"bytes"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/coyove/script/parser"
)

type Func struct {
	Code       packet
	Name       string
	DocString  string
	NumParams  byte
	IsVariadic bool
	IsDebug    bool
	StackSize  uint16
	Native     func(env *Env)
	loadGlobal *Program
	Params     []string
	Locals     []string
}

type Program struct {
	Func
	Deadline         int64
	MaxStackSize     int64
	MaxCallStackSize int64
	MaxStringSize    int64
	Stack            *[]Value
	Functions        []*Func
	Stdout, Stderr   io.Writer
	Stdin            io.Reader
	Logger           *log.Logger
	NilIndex         uint16
	Survey           struct {
		Elapsed         float64
		StringAlloc     int64
		AdjustedReturns int64
	}
	shadowTable *symtable
}

type Arguments map[string]Value

func (a Arguments) GetStringOrJSON(name string, defaultValue string) string {
	if a[name].Type() == VString {
		return a[name].String()
	}
	if !a[name].IsNil() {
		return a[name].JSONString()
	}
	return defaultValue
}

func (a Arguments) GetString(name string, defaultValue string) string {
	if a[name].Type() == VString {
		return a[name].String()
	}
	return defaultValue
}

func (a Arguments) GetInt(name string, defaultValue int64) int64 {
	if a[name].Type() == VNumber {
		return a[name].Int()
	}
	return defaultValue
}

// Native creates a golang-Native function
func Native(name string, f func(env *Env), doc ...string) Value {
	return Function(&Func{
		Name:      name,
		Native:    f,
		DocString: fixDocString(strings.Join(doc, "\n"), name, ""),
		IsDebug:   strings.HasPrefix(name, "debug_"),
	})
}

func NativeWithParamMap(name string, f func(*Env, Arguments), doc string, params ...string) Value {
	return Function(&Func{
		Name:      name,
		Params:    params,
		NumParams: byte(len(params)),
		DocString: fixDocString(doc, name, strings.Join(params, ",")),
		IsDebug:   strings.HasPrefix(name, "debug_"),
		Native: func(env *Env) {
			stack := env.Stack()
			args := make(map[string]Value, len(stack))
			for i := range stack {
				if i < len(params) {
					args[params[i]] = stack[i]
				}
			}
			f(env, Arguments(args))
		},
	})
}

func (c *Func) IsNative() bool { return c.Native != nil }

func (c *Func) String() string {
	p := bytes.Buffer{}
	if c.Name != "" {
		p.WriteString(c.Name)
	} else if c.Native != nil {
		p.WriteString("native")
	} else {
		p.WriteString("function")
	}

	p.WriteString("(")
	for i := 0; i < int(c.NumParams); i++ {
		if i < len(c.Params) {
			p.WriteString(c.Params[i])
		} else {
			p.WriteString("a" + strconv.Itoa(i))
		}
		p.WriteString(",")
	}
	if c.IsVariadic {
		p.WriteString("...")
	} else {
		if c.NumParams > 0 {
			p.Truncate(p.Len() - 1)
		}
	}
	p.WriteString(")")
	return p.String()
}

func (c *Func) PrettyCode() string {
	if c.Native != nil {
		return "[Native Code]"
	}
	return pkPrettify(c, c.loadGlobal, false, 0)
}

func (c *Func) exec(newEnv Env) (Value, []Value) {
	if c.Native != nil {
		c.Native(&newEnv)
		return newEnv.A, newEnv.V
	}
	return InternalExecCursorLoop(newEnv, c, 0)
}

func (p *Program) Run() (v1 Value, v []Value, err error) {
	return p.Call()
}

func (p *Program) Call() (v1 Value, v []Value, err error) {
	defer parser.CatchError(&err)
	start := time.Now()
	newEnv := Env{
		Global: p,
		stack:  p.Stack,
	}
	v1, v = InternalExecCursorLoop(newEnv, &p.Func, 0)
	p.Survey.Elapsed = time.Since(start).Seconds()
	return
}

func (c *Func) Call(a ...Value) (v1 Value, v []Value, err error) {
	defer parser.CatchError(&err)

	oldLen := len(*c.loadGlobal.Stack)
	newEnv := Env{
		Global:      c.loadGlobal,
		stack:       c.loadGlobal.Stack,
		StackOffset: uint32(oldLen),
	}

	var varg []Value
	for i := range a {
		if i >= int(c.NumParams) {
			varg = append(varg, a[i])
		}
		newEnv.Push(a[i])
	}

	if c.Native == nil {
		newEnv.growZero(int(c.StackSize))
		if c.IsVariadic {
			// newEnv.grow(int(c.NumParams) + 1)
			newEnv._set(uint16(c.NumParams), Array(varg))
		}
	}

	v1, v = c.exec(newEnv)
	*c.loadGlobal.Stack = (*c.loadGlobal.Stack)[:oldLen]
	return
}

func (p *Program) PrettyCode() string { return pkPrettify(&p.Func, p, true, 0) }

func (p *Program) SetTimeout(d time.Duration) { p.Deadline = time.Now().Add(d).UnixNano() }

func (p *Program) SetDeadline(d time.Time) { p.Deadline = d.UnixNano() }

func (p *Program) Print(a ...interface{}) { p.log("", "", a...) }

func (p *Program) Printf(f string, a ...interface{}) { p.log("f", f, a...) }

func (p *Program) Println(a ...interface{}) { p.log("l", "", a...) }

func (p *Program) Fatal(a ...interface{}) { p.log("F", "", a...) }

func (p *Program) Fatalf(f string, a ...interface{}) { p.log("Ff", f, a...) }

func (p *Program) Fatalln(a ...interface{}) { p.log("Fl", "", a...) }

func (p *Program) Panic(a ...interface{}) { p.log("P", "", a...) }

func (p *Program) Panicf(f string, a ...interface{}) { p.log("Pf", f, a...) }

func (p *Program) Panicln(a ...interface{}) { p.log("Pl", "", a...) }

func (p *Program) log(o, f string, a ...interface{}) {
	if p.Logger == nil {
		p.Logger = log.New(p.Stderr, "", log.LstdFlags)
	}
	switch o {
	default:
		p.Logger.Print(a...)
	case "f":
		p.Logger.Printf(f, a...)
	case "l":
		p.Logger.Println(a...)
	case "F":
		p.Logger.Fatal(a...)
	case "Ff":
		p.Logger.Fatalf(f, a...)
	case "Fl":
		p.Logger.Fatalln(a...)
	case "P":
		p.Logger.Panic(a...)
	case "Pf":
		p.Logger.Panicf(f, a...)
	case "Pl":
		p.Logger.Panicln(a...)
	}
}

func (p *Program) Get(k string) (v Value, err error) {
	defer parser.CatchError(&err)
	return (*p.Stack)[int(p.shadowTable.mustGetSymbol(k))], nil
}

func (p *Program) Set(k string, v Value) (err error) {
	defer parser.CatchError(&err)
	(*p.Stack)[int(p.shadowTable.mustGetSymbol(k))] = v
	return nil
}

func (p *Program) ResetSurvey() {
	p2 := Program{}
	p.Survey = p2.Survey
}

func (p *Program) looseStringSizeLimit() int64 {
	if p.MaxStackSize == 0 {
		return p.MaxStringSize
	}
	return p.MaxStringSize - p.MaxStringSize*int64(len(*p.Stack))/p.MaxStackSize
}

func fixDocString(in, name, arg string) string {
	in = strings.Replace(in, "$a", arg, -1)
	in = strings.Replace(in, "$f", name, -1)
	return in
}
