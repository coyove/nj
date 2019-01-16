package potatolang

const regA = 0x1fff

const (
	OP_ASSERT   = iota // 000000
	OP_STORE           // 000001
	OP_LOAD            // 000010
	OP_ADD             // 000011
	OP_SUB             // 000100
	OP_MUL             // 000101
	OP_DIV             // 000110
	OP_MOD             // 000111
	OP_NOT             // 001000
	OP_EQ              // 001001
	OP_NEQ             // 001010
	OP_LESS            // 001011
	OP_LESS_EQ         // 001100
	OP_BIT_NOT         // 001101
	OP_BIT_AND         // 001110
	OP_BIT_OR          // 001111
	OP_BIT_XOR         // 010000
	OP_BIT_LSH         // 010001
	OP_BIT_RSH         // 010010
	OP_BIT_URSH        // 010011
	OP_IF              // 010100
	OP_IFNOT           // 011000
	OP_SET             // 011100
	OP_MAKEMAP         // 000100
	OP_JMP             // 000001
	OP_LAMBDA          // 000001
	OP_CALL            // 000001
	OP_SETK            // 000001
	OP_R0              // 000001
	OP_R0K             // 000001
	OP_R1              // 000001
	OP_R1K             // 000001
	OP_R2              // 000001
	OP_R2K             // 000001
	OP_R3              // 000001
	OP_R3K             // 000001
	OP_RX              // 000001
	OP_PUSH            // 000001
	OP_PUSHK           // 000001
	OP_RET             // 000001
	OP_RETK            // 000001
	OP_YIELD           // 000001
	OP_YIELDK          // 000001
	OP_POP             // 000001
	OP_SLICE           // 000001
	OP_INC             // 000001
	OP_COPY            // 000001
	OP_LEN             // 000001
	OP_TYPEOF          // 000001
	OP_NOP             // 000001
	OP_EOB             // 000001
)
