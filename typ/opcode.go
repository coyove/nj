package typ

import "unsafe"

type Inst struct {
	Opcode    byte
	OpcodeExt byte
	A         uint16
	B         uint16
	C         uint16
}

const InstSize = unsafe.Sizeof(Inst{})

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
	OpExt
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
	OpJmpFalse
	OpJmp
	OpFunction
	OpPush
	OpPushUnpack
	OpCreateArray
	OpCreateObject
	OpCall
	OpTailCall
	OpIsProto
	OpSlice
	OpRet
	OpLoadTop

	_ = iota
	OpExtAdd16
	OpExtRSub16
	OpExtLess16
	OpExtGreat16
	OpExtEq16
	OpExtNeq16
	OpExtInc16
	OpExtLoad16
	OpExtStore16
	OpExtBitAnd
	OpExtBitOr
	OpExtBitXor
	OpExtBitLsh
	OpExtBitRsh
	OpExtBitURsh
	OpExtBitAnd16
	OpExtBitOr16
	OpExtBitXor16
	OpExtBitLsh16
	OpExtBitRsh16
	OpExtBitURsh16
)

func JmpInst(op byte, distance int) Inst {
	if distance < -(1<<30) || distance >= 1<<30 {
		panic("long jump")
	}
	return (Inst{Opcode: op}).SetD(int32(distance))
}

var BinaryOpcode = map[byte]string{
	OpAdd:        "add",
	OpSub:        "sub",
	OpMul:        "mul",
	OpDiv:        "div",
	OpIDiv:       "idiv",
	OpMod:        "mod",
	OpEq:         "eq",
	OpNeq:        "neq",
	OpLess:       "less",
	OpLessEq:     "lesseq",
	OpLoad:       "load",
	OpNext:       "next",
	OpIsProto:    "isproto",
	OpInc:        "inc",
	OpExtBitAnd:  "and",
	OpExtBitOr:   "or",
	OpExtBitXor:  "xor",
	OpExtBitLsh:  "lsh",
	OpExtBitRsh:  "rsh",
	OpExtBitURsh: "ursh",
}

var UnaryOpcode = map[byte]string{
	OpNot:        "not",
	OpRet:        "return",
	OpLen:        "len",
	OpPush:       "push",
	OpPushUnpack: "pushvarg",
}

var TenaryOpcode = map[byte]string{
	OpLoad:  "load",
	OpStore: "store",
	OpSlice: "slice",
}
