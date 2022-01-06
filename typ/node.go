package typ

import "fmt"

// Align with S-variables in parser/parser2.go
const (
	ADoBlock  = "do"
	ANil      = "nil"
	ASet      = "set"
	AInc      = "incr"
	AMove     = "move"
	AIf       = "if"
	AFor      = "loop"
	AFunc     = "function"
	ABreak    = "break"
	AContinue = "continue"
	ABegin    = "prog"
	ALoad     = "load"
	AStore    = "store"
	AArray    = "array"
	AObject   = "map"
	ACall     = "call"
	ATailCall = "tailcall"
	AReturn   = "return"
	ALen      = "len"
	ANext     = "next"
	AAdd      = "add"
	ASub      = "sub"
	AMul      = "mul"
	ADiv      = "div"
	AIDiv     = "idiv"
	AMod      = "mod"
	ABitAnd   = "bitand"
	ABitOr    = "bitor"
	ABitXor   = "bitxor"
	ABitNot   = "bitnot"
	ABitLsh   = "bitlsh"
	ABitRsh   = "bitrsh"
	ABitURsh  = "bitursh"
	AEq       = "eq"
	ANeq      = "neq"
	AAnd      = "and"
	AOr       = "or"
	ANot      = "not"
	ALess     = "lt"
	ALessEq   = "le"
	AFreeAddr = "freeaddr"
	ALabel    = "label"
	AGoto     = "goto"
	AUnpack   = "unpack"
	AIs       = "isproto"
)

type Position struct {
	Source string
	Line   uint32
	Column uint32
}

func (pos *Position) String() string {
	return fmt.Sprintf("%s:%d:%d", pos.Source, pos.Line, pos.Column)
}

var BinaryOpcode = map[byte]string{
	OpAdd:     AAdd,
	OpSub:     ASub,
	OpMul:     AMul,
	OpDiv:     ADiv,
	OpIDiv:    AIDiv,
	OpMod:     AMod,
	OpEq:      AEq,
	OpNeq:     ANeq,
	OpLess:    ALess,
	OpLessEq:  ALessEq,
	OpLoad:    ALoad,
	OpStore:   AStore,
	OpBitAnd:  ABitAnd,
	OpBitOr:   ABitOr,
	OpBitXor:  ABitXor,
	OpBitLsh:  ABitLsh,
	OpBitRsh:  ABitRsh,
	OpBitURsh: ABitURsh,
	OpNext:    ANext,
	OpIsProto: AIs,
}

var UnaryOpcode = map[byte]string{
	OpBitNot:     ABitNot,
	OpNot:        ANot,
	OpRet:        AReturn,
	OpLen:        ALen,
	OpPush:       "push",
	OpPushUnpack: "pushvararg",
}
