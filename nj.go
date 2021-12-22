package nj

import (
	"io/ioutil"
	"os"
	"strings"
	"unsafe"

	"github.com/coyove/nj/bas"
	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/parser"
)

func init() {
	bas.Globals.SetProp("json", bas.NamedObject("json", 0).
		SetMethod("stringify", func(e *bas.Env) {
			e.A = bas.Str(e.Get(0).JSONString())
		}, "$f(v: value) -> string").
		SetMethod("parse", func(e *bas.Env) {
			v, err := parser.ParseJSON(strings.TrimSpace(e.Str(0)))
			internal.PanicErr(err)
			e.A = v
		}, "$f(j: string) -> value").
		ToValue())
	bas.Globals.SetMethod("loadfile", func(e *bas.Env) {
		e.A = MustRun(LoadFile(e.Str(0), &e.Global.Environment))
	}, "$f(path: string) -> value\n\tload and eval file at `path`, globals will be inherited in loaded file")
	bas.Globals.SetMethod("eval", func(e *bas.Env) {
		opts := e.Get(1).Safe().Object()
		if opts.Prop("ast").IsTrue() {
			v, err := parser.Parse(e.Str(0), "")
			internal.PanicErr(err)
			e.A = bas.ValueOf(v)
			return
		}
		p, err := LoadString(e.Str(0), &bas.Environment{Globals: opts.Prop("globals").Safe().Object()})
		internal.PanicErr(err)
		v, err := p.Run()
		internal.PanicErr(err)
		_ = opts.Prop("returnglobals").IsTrue() && e.SetA(p.LocalsObject().ToValue()) || e.SetA(v)
	}, "$f(code: string, options?: object) -> value\n\tevaluate `code` and return the reuslt")
}

func LoadFile(path string, opt *bas.Environment) (*bas.Program, error) {
	code, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return loadCode(*(*string)(unsafe.Pointer(&code)), path, opt)
}

func LoadString(code string, opt *bas.Environment) (*bas.Program, error) {
	return loadCode(code, "<memory>", opt)
}

func loadCode(code, name string, opt *bas.Environment) (*bas.Program, error) {
	n, err := parser.Parse(code, name)
	if err != nil {
		return nil, err
	}
	if internal.IsDebug() {
		n.Dump(os.Stderr, "  ")
	}
	return compileNodeTopLevel(name, code, n, opt)
}

func Run(p *bas.Program, err error) (bas.Value, error) {
	if err != nil {
		return bas.Nil, err
	}
	return p.Run()
}

func MustRun(p *bas.Program, err error) bas.Value {
	internal.PanicErr(err)
	v, err := p.Run()
	internal.PanicErr(err)
	return v
}
