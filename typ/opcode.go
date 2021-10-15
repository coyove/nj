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
	OpPushVararg
	OpArray
	OpMap
	OpCall
	OpTailCall
	OpRet
)

type ValueType byte

const (
	Nil       ValueType = 0
	Bool      ValueType = 1
	Number    ValueType = 3
	String    ValueType = 7
	Table     ValueType = 15
	Func      ValueType = 17
	Interface ValueType = 19
)

func (t ValueType) String() string {
	if t > Interface {
		return "?"
	}
	return [...]string{"nil", "bool", "?", "number", "?", "?", "?", "string", "?", "?", "?", "?", "?", "?", "?", "table", "?", "function", "?", "interface"}[t]
}
