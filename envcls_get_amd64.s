// envcls.go:95		0x11049d2		664181faff1f		CMPW $0x1fff, R10						
// envcls.go:95		0x11049d8		743f				JE 0x1104a19							
// ---> envcls.go:96		0x1104a19		488b4e20			MOVQ 0x20(SI), CX // 0x20: env.A
// 
// envcls.go:98		0x11049da		6641c1ea0a			SHRW $0xa, R10							
// envcls.go:99		0x11049df		4181e4ff030000	    ANDL $0x3ff, R12						
// envcls.go:101		0x11049e6		664183fa07			CMPW $0x7, R10							
// envcls.go:101		0x11049eb		740f				JE 0x11049fc							
// ---> envcls.go:102		0x11049fc		488b4a50			MOVQ 0x50(DX), CX						
//      envcls.go:102		0x1104a00		488b5248			MOVQ 0x48(DX), DX						
//      envcls.go:102		0x1104a04		4939cc				CMPQ CX, R12							
//      envcls.go:102		0x1104a07		0f8374720000		JAE 0x110bc81 // index overflow check
//      envcls.go:102		0x1104a0d		4a8b0ce2			MOVQ 0(DX)(R12*8), CX
// 
// mainloop.go:137	    0x11049ed		4889f0				MOVQ SI, AX							
// envcls.go:106		0x1104901		664585d2			TESTW R10, R10							
// envcls.go:106		0x1104905		7605				JBE 0x110490c							
// envcls.go:106		0x1104907		4885f6				TESTQ SI, SI							
// envcls.go:106		0x110490a		75ef				JNE 0x11048fb							
// --->   envcls.go:107		0x11048fb		41ffca				DECL R10							
//   	   envcls.go:107		0x11048fe		488b36				MOVQ 0(SI), SI
// 	   JMP 0x11049
// 
// envcls.go:110		0x110490c		488b5608			MOVQ 0x8(SI), DX						
// envcls.go:110		0x1104910		488b4e10			MOVQ 0x10(SI), CX						
// envcls.go:110		0x1104914		4939cc				CMPQ CX, R12							
// envcls.go:110		0x1104917		0f8dd8000000		JGE 0x11049f5							
// ---> envcls.go:106		0x11049f5		31c9				XORL CX, CX	// return Value{}
// 
// envcls.go:111		0x110491d		4885c9				TESTQ CX, CX							
// envcls.go:111		0x1104920		0f8654730000	 	JBE 0x110bc7a							
// envcls.go:111		0x1104926		4a8b0ce2			MOVQ 0(DX)(R12*8), CX // retuen env.stack[i]

#include "textflag.h"

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

