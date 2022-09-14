package nj

import (
	"io"
	"io/ioutil"
	"os"
	"unsafe"

	"github.com/coyove/nj/bas"
	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/parser"
)

type LoadOptions struct {
	Globals      *bas.Object
	MaxStackSize int64
	Stdout       io.Writer
	Stderr       io.Writer
	Stdin        io.Reader
}

func LoadFile(path string, opt *LoadOptions) (*bas.Program, error) {
	code, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return loadCode(*(*string)(unsafe.Pointer(&code)), path, opt)
}

func LoadString(code string, opt *LoadOptions) (*bas.Program, error) {
	return loadCode(code, internal.UnnamedLoadString(), opt)
}

func loadCode(code, name string, opt *LoadOptions) (*bas.Program, error) {
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
