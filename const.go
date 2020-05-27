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
	OpGetB
	OpSetB
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
	OpMakeHash
	OpMakeArray
	OpLambda
	OpPush
	OpPush2
	OpCall
	OpRet
	OpYield
	OpLen
	OpPatchVararg
	OpAddressOf
	OpNOP
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
