package typ

import "unsafe"

type Inst struct {
	Opcode byte
	A      uint16
	B      uint16
	C      uint16
}

func (i Inst) D() int32 {
	// return int32(uint32(i.B)<<16 | uint32(i.C))
	return *(*int32)(unsafe.Pointer(&i.B))
}

func (i Inst) SetD(d int32) Inst {
	// i.B = uint16(uint32(d) >> 16)
	// i.C = uint16(uint32(d))
	*(*int32)(unsafe.Pointer(&i.B)) = d
	return i
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
	OpSlice
	OpRet
)

func JmpInst(op byte, distance int) Inst {
	if distance < -(1<<30) || distance >= 1<<30 {
		panic("long jump")
	}
	return (Inst{Opcode: op}).SetD(int32(distance))
}
