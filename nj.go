package nj

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync/atomic"
	"unsafe"

	"github.com/coyove/nj/bas"
	"github.com/coyove/nj/internal"
	"github.com/coyove/nj/parser"
)

var loadIndex int64

func LoadFile(path string, opt *bas.Environment) (*bas.Program, error) {
	code, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return loadCode(*(*string)(unsafe.Pointer(&code)), path, opt)
}

func LoadString(code string, opt *bas.Environment) (*bas.Program, error) {
	return loadCode(code, fmt.Sprintf("<memory-%d>", atomic.AddInt64(&loadIndex, 1)), opt)
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
