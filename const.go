package potatolang

const regA = 0xffffffff

const (
	// basic flat op
	OP_ASSERT = iota
	OP_NIL
	OP_TRUE
	OP_FALSE
	OP_STORE
	OP_LOAD
	OP_ADD
	OP_SUB
	OP_MUL
	OP_DIV
	OP_MOD
	OP_NOT
	OP_EQ
	OP_NEQ
	OP_LESS
	OP_LESS_EQ
	OP_BIT_NOT
	OP_BIT_AND
	OP_BIT_OR
	OP_BIT_XOR
	OP_BIT_LSH
	OP_BIT_RSH

	// complex structure op
	OP_LIST
	OP_MAP

	// flow control op
	OP_IF
	OP_IFNOT
	OP_JMP
	OP_LAMBDA
	OP_CALL
	OP_SET
	OP_SETK
	OP_R0
	OP_R0K
	OP_R1
	OP_R1K
	OP_R2
	OP_R2K
	OP_R3
	OP_R3K
	OP_PUSH
	OP_PUSHK
	OP_RET
	OP_RETK
	OP_YIELD
	OP_YIELDK

	// inc op
	OP_INC

	// special builtin op
	OP_WHO
	OP_DUP
	OP_ERROR
	OP_LEN
	OP_TYPEOF

	OP_LINE = 0xFD
	OP_NOP  = 0xFE
	OP_EOB  = 0xFF
)
