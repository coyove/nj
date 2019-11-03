package potatolang

const (
	regA   uint16 = 0x1fff // full 13 bits
	regNil uint16 = 0x3ff - 1
)

type (
	_Opcode byte
)

const (
	_ _Opcode = iota
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
	OpBitAnd
	OpBitOr
	OpBitXor
	OpBitLsh
	OpBitRsh
	OpBitURsh
	OpIf
	OpIfNot
	OpSet
	OpMakeStruct
	OpMakeSlice
	OpJmp
	OpLambda
	OpCall
	OpPush
	OpPush2
	OpRet
	OpYield
	OpSlice
	OpInc
	OpCopyStack
	OpLen
	OpTypeof
	OpAddressOf
	OpNOP
	OpEOB
)
