// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

TEXT errors(SB),$0
	MOVW	(F0), R1           // ERROR "illegal base register"
	MOVB	(F0), R1           // ERROR "illegal base register"
	MOVH	(F0), R1           // ERROR "illegal base register"
	MOVF	(F0), F1           // ERROR "illegal base register"
	MOVD	(F0), F1           // ERROR "illegal base register"
	MOVW	R1, (F0)           // ERROR "illegal base register"
	MOVB	R2, (F0)           // ERROR "illegal base register"
	MOVH	R3, (F0)           // ERROR "illegal base register"
	MOVF	F4, (F0)           // ERROR "illegal base register"
	MOVD	F5, (F0)           // ERROR "illegal base register"
	MOVM.IA	(F1), [R0-R4]      // ERROR "illegal base register"
	MOVM.DA	(F1), [R0-R4]      // ERROR "illegal base register"
	MOVM.IB	(F1), [R0-R4]      // ERROR "illegal base register"
	MOVM.DB	(F1), [R0-R4]      // ERROR "illegal base register"
	MOVM.IA	[R0-R4], (F1)      // ERROR "illegal base register"
	MOVM.DA	[R0-R4], (F1)      // ERROR "illegal base register"
	MOVM.IB	[R0-R4], (F1)      // ERROR "illegal base register"
	MOVM.DB	[R0-R4], (F1)      // ERROR "illegal base register"
	MOVW	R0<<0(F1), R1      // ERROR "illegal base register"
	MOVB	R0<<0(F1), R1      // ERROR "illegal base register"
	MOVW	R1, R0<<0(F1)      // ERROR "illegal base register"
	MOVB	R2, R0<<0(F1)      // ERROR "illegal base register"
	MOVF	0x00ffffff(F2), F1 // ERROR "illegal base register"
	MOVD	0x00ffffff(F2), F1 // ERROR "illegal base register"
	MOVF	F2, 0x00ffffff(F2) // ERROR "illegal base register"
	MOVD	F2, 0x00ffffff(F2) // ERROR "illegal base register"
	MULS.S	R1, R2, R3, R4     // ERROR "invalid .S suffix"
	ADD.P	R1, R2, R3         // ERROR "invalid .P suffix"
	SUB.W	R2, R3             // ERROR "invalid .W suffix"
	BL	4(R4)              // ERROR "non-zero offset"
	ADDF	F0, R1, F2         // ERROR "illegal combination"
	SWI	(R0)               // ERROR "illegal combination"
	MULAD	F0, F1             // ERROR "illegal combination"
	MULAF	F0, F1             // ERROR "illegal combination"
	MULSD	F0, F1             // ERROR "illegal combination"
	MULSF	F0, F1             // ERROR "illegal combination"
	NMULAD	F0, F1             // ERROR "illegal combination"
	NMULAF	F0, F1             // ERROR "illegal combination"
	NMULSD	F0, F1             // ERROR "illegal combination"
	NMULSF	F0, F1             // ERROR "illegal combination"
	FMULAD	F0, F1             // ERROR "illegal combination"
	FMULAF	F0, F1             // ERROR "illegal combination"
	FMULSD	F0, F1             // ERROR "illegal combination"
	FMULSF	F0, F1             // ERROR "illegal combination"
	FNMULAD	F0, F1             // ERROR "illegal combination"
	FNMULAF	F0, F1             // ERROR "illegal combination"
	FNMULSD	F0, F1             // ERROR "illegal combination"
	FNMULSF	F0, F1             // ERROR "illegal combination"
	NEGF	F0, F1, F2         // ERROR "illegal combination"
	NEGD	F0, F1, F2         // ERROR "illegal combination"
	ABSF	F0, F1, F2         // ERROR "illegal combination"
	ABSD	F0, F1, F2         // ERROR "illegal combination"
	SQRTF	F0, F1, F2         // ERROR "illegal combination"
	SQRTD	F0, F1, F2         // ERROR "illegal combination"
	MOVF	F0, F1, F2         // ERROR "illegal combination"
	MOVD	F0, F1, F2         // ERROR "illegal combination"
	MOVDF	F0, F1, F2         // ERROR "illegal combination"
	MOVFD	F0, F1, F2         // ERROR "illegal combination"
	MOVM.IA	4(R1), [R0-R4]     // ERROR "offset must be zero"
	MOVM.DA	4(R1), [R0-R4]     // ERROR "offset must be zero"
	MOVM.IB	4(R1), [R0-R4]     // ERROR "offset must be zero"
	MOVM.DB	4(R1), [R0-R4]     // ERROR "offset must be zero"
	MOVM.IA	[R0-R4], 4(R1)     // ERROR "offset must be zero"
	MOVM.DA	[R0-R4], 4(R1)     // ERROR "offset must be zero"
	MOVM.IB	[R0-R4], 4(R1)     // ERROR "offset must be zero"
	MOVM.DB	[R0-R4], 4(R1)     // ERROR "offset must be zero"
	MOVW	CPSR, FPSR         // ERROR "illegal combination"
	MOVW	FPSR, CPSR         // ERROR "illegal combination"
	MOVW	CPSR, errors(SB)   // ERROR "illegal combination"
	MOVW	errors(SB), CPSR   // ERROR "illegal combination"
	MOVW	FPSR, errors(SB)   // ERROR "illegal combination"
	MOVW	errors(SB), FPSR   // ERROR "illegal combination"
	MOVW	F0, errors(SB)     // ERROR "illegal combination"
	MOVW	errors(SB), F0     // ERROR "illegal combination"
	MOVW	$20, errors(SB)    // ERROR "illegal combination"
	MOVW	errors(SB), $20    // ERROR "illegal combination"
	MOVB	$245, R1           // ERROR "illegal combination"
	MOVH	$245, R1           // ERROR "illegal combination"
	MOVB	$0xff000000, R1    // ERROR "illegal combination"
	MOVH	$0xff000000, R1    // ERROR "illegal combination"
	MOVB	$0x00ffffff, R1    // ERROR "illegal combination"
	MOVH	$0x00ffffff, R1    // ERROR "illegal combination"
	MOVB	FPSR, g            // ERROR "illegal combination"
	MOVH	FPSR, g            // ERROR "illegal combination"
	MOVB	g, FPSR            // ERROR "illegal combination"
	MOVH	g, FPSR            // ERROR "illegal combination"
	MOVB	CPSR, g            // ERROR "illegal combination"
	MOVH	CPSR, g            // ERROR "illegal combination"
	MOVB	g, CPSR            // ERROR "illegal combination"
	MOVH	g, CPSR            // ERROR "illegal combination"
	MOVB	$0xff000000, CPSR  // ERROR "illegal combination"
	MOVH	$0xff000000, CPSR  // ERROR "illegal combination"
	MOVB	$0xff000000, FPSR  // ERROR "illegal combination"
	MOVH	$0xff000000, FPSR  // ERROR "illegal combination"
	MOVB	$0xffffff00, CPSR  // ERROR "illegal combination"
	MOVH	$0xffffff00, CPSR  // ERROR "illegal combination"
	MOVB	$0xfffffff0, FPSR  // ERROR "illegal combination"
	MOVH	$0xfffffff0, FPSR  // ERROR "illegal combination"
	MOVB.IA	4(R1), [R0-R4]     // ERROR "illegal combination"
	MOVB.DA	4(R1), [R0-R4]     // ERROR "illegal combination"
	MOVH.IA	4(R1), [R0-R4]     // ERROR "illegal combination"
	MOVH.DA	4(R1), [R0-R4]     // ERROR "illegal combination"
	MOVB	$0xff(R0), R1      // ERROR "illegal combination"
	MOVH	$0xff(R0), R1      // ERROR "illegal combination"
	MOVB	$errors(SB), R2    // ERROR "illegal combination"
	MOVH	$errors(SB), R2    // ERROR "illegal combination"
	MOVB	F0, R0             // ERROR "illegal combination"
	MOVH	F0, R0             // ERROR "illegal combination"
	MOVB	R0, F0             // ERROR "illegal combination"
	MOVH	R0, F0             // ERROR "illegal combination"
	MOVB	R0>>0(R1), R2      // ERROR "bad shift"
	MOVB	R0->0(R1), R2      // ERROR "bad shift"
	MOVB	R0@>0(R1), R2      // ERROR "bad shift"
	MOVBS	R0>>0(R1), R2      // ERROR "bad shift"
	MOVBS	R0->0(R1), R2      // ERROR "bad shift"
	MOVBS	R0@>0(R1), R2      // ERROR "bad shift"
	MOVF	CPSR, F1           // ERROR "illegal combination"
	MOVD	R1, CPSR           // ERROR "illegal combination"
	MOVW	F1, F2             // ERROR "illegal combination"
	MOVB	F1, F2             // ERROR "illegal combination"
	MOVH	F1, F2             // ERROR "illegal combination"
	MOVF	R1, F2             // ERROR "illegal combination"
	MOVD	R1, F2             // ERROR "illegal combination"
	MOVF	R1, R1             // ERROR "illegal combination"
	MOVD	R1, R2             // ERROR "illegal combination"
	MOVFW	R1, R2             // ERROR "illegal combination"
	MOVDW	R1, R2             // ERROR "illegal combination"
	MOVWF	R1, R2             // ERROR "illegal combination"
	MOVWD	R1, R2             // ERROR "illegal combination"
	MOVWD	CPSR, R2           // ERROR "illegal combination"
	MOVWF	CPSR, R2           // ERROR "illegal combination"
	MOVWD	R1, CPSR           // ERROR "illegal combination"
	MOVWF	R1, CPSR           // ERROR "illegal combination"
	MOVDW	CPSR, R2           // ERROR "illegal combination"
	MOVFW	CPSR, R2           // ERROR "illegal combination"
	MOVDW	R1, CPSR           // ERROR "illegal combination"
	MOVFW	R1, CPSR           // ERROR "illegal combination"
	BFX	$12, $41, R2, R3   // ERROR "wrong width or LSB"
	BFX	$12, $-2, R2       // ERROR "wrong width or LSB"
	BFXU	$40, $4, R2, R3    // ERROR "wrong width or LSB"
	BFXU	$-40, $4, R2       // ERROR "wrong width or LSB"
	BFX	$-2, $4, R2, R3    // ERROR "wrong width or LSB"
	BFXU	$4, R2, R5, R2     // ERROR "missing or wrong LSB"
	BFXU	$4, R2, R5         // ERROR "missing or wrong LSB"

	END
