package typ

type Inst struct {
	Opcode byte
	A      uint16
	B      int32
}

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
	OpLen
	OpNext
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
	OpPushUnpack
	OpCreateArray
	OpCreateObject
	OpCall
	OpTailCall
	OpIsProto
	OpRet
)
