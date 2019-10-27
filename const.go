package potatolang

const (
	regA   uint16 = 0x1fff // full 13 bits
	regNil uint16 = 0x3ff - 1
)

type (
	_Opcode byte
)

const (
	OpAssert _Opcode = iota
	OpStore
	OpLoad
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpMod
	OpNot
	OpEq
	OpNeq
	OpLess
	OpLessEq
	OpBitNot
	OpBitAnd
	OpBitOr
	OpBitXor
	OpBitLsh
	OpBitRsh
	OpBitURsh
	OpIf
	OpIfNot
	OpSet
	OpMakeMap
	OpJmp
	OpLambda
	OpCall
	OpPush
	OpRet
	OpYield
	OpPop
	OpSlice
	OpInc
	OpForeach
	OpLen
	OpTypeof
	OpAddressOf
	OpNOP
	OpEOB
)
