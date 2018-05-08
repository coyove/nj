package compiler

import (
	"github.com/coyove/bracket/base"
	"github.com/coyove/bracket/parser"
)

func compileAttachOp(stackPtr int16, atoms []*parser.Node, varLookup *base.CMap) (code []byte, yx int32, newStackPtr int16, err error) {
	var buf *base.BytesReader

}
