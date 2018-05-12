package compiler

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/coyove/bracket/base"
	"github.com/coyove/bracket/parser"
)

func parse(n *parser.Node) (code []byte, err error) {
	varLookup := base.NewCMap()
	for i, n := range base.CoreLibNames {
		varLookup.M[n] = int16(i)
	}

	code, _, _, err = compileChainOp(int16(len(varLookup.M)), n, varLookup)
	if err != nil {
		return
	}

	code = append(code, base.OP_EOB)
	return code, nil
}

func LoadFile(path string) ([]byte, error) {
	code, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	n, err := parser.Parse(bytes.NewReader(code), path)
	if err != nil {
		return nil, err
	}

	n.Dump(os.Stderr)
	return parse(n)
}
