#include "textflag.h"

// This is the assembly version of getting a value from *Env
// However without inlining it won't be any faster than the native Go one

// func envGet(env *Env, yx uint16, K *Closure) Value
TEXT ·envGet(SB),NOSPLIT,$0-32
	MOVQ env+0(FP), R9
	MOVQ yx+8(FP), R10
	CMPW R10, $0x1fff
	JE return_reg_a
	MOVQ R10, R12
	SHRW $0xa, R10							
	ANDL $0x3ff, R12
	CMPW R10, $0x07
	JE return_const
find_parent:
	TESTW R10, R10							
	JBE 4(PC) // goto got_parent
	DECL R10							
	MOVQ 0(R9), R9 // env = env.parent
	JMP find_parent
got_parnet:
	MOVQ 0x10(R9), CX						
	CMPQ R12, CX
	JGE return_nil
	MOVQ 0x8(R9), CX						
	MOVQ 0(CX)(R12*8), R9
	MOVQ R9, ret+24(FP) // retuen env.stack[index]
	RET
return_nil:
	XORQ R9, R9
	MOVQ R9, ret+24(FP) // Value{}
	RET
return_const:
	MOVQ K+16(FP), DX
	MOVQ 0x48(DX), DX						
	MOVQ 0(DX)(R12*8), R9
	MOVQ R9, ret+24(FP) // K.consts[index]
	RET
return_reg_a:
	MOVQ 0x20(R9), R9
	MOVQ R9, ret+24(FP) // 0x20: env.A
	RET

