package script

import (
	"fmt"
)

const (
	regA uint16 = 0x1fff // full 13 bits

	callNormal = 0
	callTail   = 1
)

const (
	_ = iota
	OpSet
	OpStore
	OpLoad
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpIDiv
	OpInc
	OpMod
	OpNot
	OpEq
	OpNeq
	OpLess
	OpLessEq
	OpBitAnd
	OpBitOr
	OpBitXor
	OpBitNot
	OpBitLsh
	OpBitRsh
	OpBitURsh
	OpIfNot
	OpJmp
	OpLoadFunc
	OpPush
	OpPushVararg
	OpArray
	OpMap
	OpCall
	OpRet
)

func panicf(msg string, args ...interface{}) Value {
	panic(fmt.Errorf(msg, args...))
}
