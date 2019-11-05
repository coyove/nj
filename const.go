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
	OpSet
	OpStore
	OpLoad
	OpAdd
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
	OpBitAnd
	OpBitOr
	OpBitXor
	OpBitLsh
	OpBitRsh
	OpBitURsh
	OpIf
	OpIfNot
	OpJmp
	OpMakeStruct
	OpMakeSlice
	OpLambda
	OpPush
	OpPush2
	OpCall
	OpRet
	OpYield
	OpSlice
	OpLen
	OpCopyStack
	OpTypeof
	OpAddressOf
	OpNOP
	OpEOB
)
