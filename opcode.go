package script

import (
	"fmt"
)

const (
	regA       uint16 = 0x1fff // full 13 bits
	regPhantom uint16 = 0x1ffe // full 13 bits
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
	OpTailCall
	OpRet
)

func panicf(msg string, args ...interface{}) Value {
	panic(fmt.Errorf(msg, args...))
}
