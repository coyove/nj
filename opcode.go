package script

import (
	"fmt"
)

const regA uint16 = 0x1fff // full 13 bits

type opCode byte

const (
	_ opCode = iota
	OpSet
	OpStore
	OpGStore
	OpLoad
	OpGLoad
	OpSlice
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
	OpPow
	OpIf
	OpIfNot
	OpJmp
	OpLoadFunc
	OpPush
	OpMapArray
	OpMap
	OpCall
	OpCallMap
	OpRet
	OpLen
)

func panicf(msg string, args ...interface{}) {
	panic(fmt.Errorf(msg, args...))
}
