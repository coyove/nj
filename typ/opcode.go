package typ

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
	OpRet
)

type ValueType byte

const (
	Nil    ValueType = 0
	Bool   ValueType = 1
	Number ValueType = 3
	String ValueType = 7
	Object ValueType = 15
	Array  ValueType = 17
	Native ValueType = 19
)

func (t ValueType) String() string {
	if t > Native {
		return "?"
	}
	return [...]string{"nil", "bool", "?", "number", "?", "?", "?", "string", "?", "?", "?", "?", "?", "?", "?", "object", "?", "array", "?", "native"}[t]
}
