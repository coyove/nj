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

func JmpInst(op byte, distance int) Inst {
	if distance < -(1<<30) || distance >= 1<<30 {
		panic("long jump")
	}
	return Inst{Opcode: op, B: int32(distance)}
}
