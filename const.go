package potatolang

import "fmt"

const (
	regA   uint16 = 0x1fff // full 13 bits
	regNil uint16 = 0x3ff - 1
)

type _Opcode byte

const (
	_ _Opcode = iota
	OpSet
	OpStore
	OpLoad
	OpAdd
	OpConcat
	OpSub
	OpMul
	OpDiv
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
	OpPushV
	OpPopV
	OpPopVAll
	OpPopVClear
	OpCall
	OpRet
	OpYield
	OpLen
	OpEOB
)

func panicerr(err error) {
	if err != nil {
		panic(err)
	}
}

func panicf(msg string, args ...interface{}) {
	panic(fmt.Errorf(msg, args...))
}
