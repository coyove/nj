package parser

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func (lex *Lexer) loadFile(path string) *Node {
	if strings.Contains(lex.loop, path) {
		lex.Error(fmt.Sprintf("%s and %s are importing each other", lex.scanner.Pos.Source, path))
	}

	if n, ok := lex.cache[path]; ok {
		return cloneNode(n)
	}

	code, err := ioutil.ReadFile(path)
	if err != nil {
		lex.Error(err.Error())
	}
	n, _, err := parse(bytes.NewReader(code), path, lex.cache, lex.loop+"?"+path)
	if err != nil {
		lex.Error(err.Error())
	}

	// Now the required code is loaded, we will wrap them into a closure
	pos := Position{Line: 1, Column: 1, Source: path}
	cls := __func(Cpl(), n).pos0(pos)
	node := __call(cls, Cpl()).pos0(pos)
	lex.cache[path] = node
	return cloneNode(node)
}

func joinSourcePath(path1 string, filename2 string) string {
	return filepath.Join(filepath.Dir(path1), filename2)
}

func moduleNameFromPath(path string) string {
	fn := filepath.Base(path)
	fn = fn[:len(fn)-len(filepath.Ext(fn))]
	return fn
}

func cloneNode(n *Node) *Node {
	if n.Type() == CPL {
		n2 := make([]*Node, 0, len(n.Cpl()))
		for _, n := range n.Cpl() {
			n2 = append(n2, cloneNode(n))
		}
		tmp := *n
		tmp.Value = n2
		return &tmp
	}
	return n
}
