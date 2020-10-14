// Code generated by mkpreempt.go; DO NOT EDIT.

#include "go_asm.h"
#include "textflag.h"

// Note: asyncPreempt doesn't use the internal ABI, but we must be able to inject calls to it from the signal handler, so Go code has to see the PC of this function literally.
TEXT ·asyncPreempt<ABIInternal>(SB),NOSPLIT|NOFRAME,$0-0
	MOVW.W R14, -188(R13)
	MOVW R0, 4(R13)
	MOVW R1, 8(R13)
	MOVW R2, 12(R13)
	MOVW R3, 16(R13)
	MOVW R4, 20(R13)
	MOVW R5, 24(R13)
	MOVW R6, 28(R13)
	MOVW R7, 32(R13)
	MOVW R8, 36(R13)
	MOVW R9, 40(R13)
	MOVW R11, 44(R13)
	MOVW R12, 48(R13)
	MOVW CPSR, R0
	MOVW R0, 52(R13)
	MOVB ·goarm(SB), R0
	CMP $6, R0
	BLT nofp
	MOVW FPCR, R0
	MOVW R0, 56(R13)
	MOVD F0, 60(R13)
	MOVD F1, 68(R13)
	MOVD F2, 76(R13)
	MOVD F3, 84(R13)
	MOVD F4, 92(R13)
	MOVD F5, 100(R13)
	MOVD F6, 108(R13)
	MOVD F7, 116(R13)
	MOVD F8, 124(R13)
	MOVD F9, 132(R13)
	MOVD F10, 140(R13)
	MOVD F11, 148(R13)
	MOVD F12, 156(R13)
	MOVD F13, 164(R13)
	MOVD F14, 172(R13)
	MOVD F15, 180(R13)
nofp:
	CALL ·asyncPreempt2(SB)
	MOVB ·goarm(SB), R0
	CMP $6, R0
	BLT nofp2
	MOVD 180(R13), F15
	MOVD 172(R13), F14
	MOVD 164(R13), F13
	MOVD 156(R13), F12
	MOVD 148(R13), F11
	MOVD 140(R13), F10
	MOVD 132(R13), F9
	MOVD 124(R13), F8
	MOVD 116(R13), F7
	MOVD 108(R13), F6
	MOVD 100(R13), F5
	MOVD 92(R13), F4
	MOVD 84(R13), F3
	MOVD 76(R13), F2
	MOVD 68(R13), F1
	MOVD 60(R13), F0
	MOVW 56(R13), R0
	MOVW R0, FPCR
nofp2:
	MOVW 52(R13), R0
	MOVW R0, CPSR
	MOVW 48(R13), R12
	MOVW 44(R13), R11
	MOVW 40(R13), R9
	MOVW 36(R13), R8
	MOVW 32(R13), R7
	MOVW 28(R13), R6
	MOVW 24(R13), R5
	MOVW 20(R13), R4
	MOVW 16(R13), R3
	MOVW 12(R13), R2
	MOVW 8(R13), R1
	MOVW 4(R13), R0
	MOVW 188(R13), R14
	MOVW.P 192(R13), R15
	UNDEF
