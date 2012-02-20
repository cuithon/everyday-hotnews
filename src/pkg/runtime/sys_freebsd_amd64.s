// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// System calls and other sys.stuff for AMD64, FreeBSD
// /usr/src/sys/kern/syscalls.master for syscall numbers.
//

#include "zasm_GOOS_GOARCH.h"
	
TEXT runtime·sys_umtx_op(SB),7,$0
	MOVQ 8(SP), DI
	MOVL 16(SP), SI
	MOVL 20(SP), DX
	MOVQ 24(SP), R10
	MOVQ 32(SP), R8
	MOVL $454, AX
	SYSCALL
	RET

TEXT runtime·thr_new(SB),7,$0
	MOVQ 8(SP), DI
	MOVQ 16(SP), SI
	MOVL $455, AX
	SYSCALL
	RET

TEXT runtime·thr_start(SB),7,$0
	MOVQ	DI, R13	// m

	// set up FS to point at m->tls
	LEAQ	m_tls(R13), DI
	CALL	runtime·settls(SB)	// smashes DI

	// set up m, g
	get_tls(CX)
	MOVQ	R13, m(CX)
	MOVQ	m_g0(R13), DI
	MOVQ	DI, g(CX)

	CALL runtime·stackcheck(SB)
	CALL runtime·mstart(SB)
	MOVQ 0, AX			// crash (not reached)

// Exit the entire program (like C exit)
TEXT runtime·exit(SB),7,$-8
	MOVL	8(SP), DI		// arg 1 exit status
	MOVL	$1, AX
	SYSCALL
	CALL	runtime·notok(SB)
	RET

TEXT runtime·exit1(SB),7,$-8
	MOVQ	8(SP), DI		// arg 1 exit status
	MOVL	$431, AX
	SYSCALL
	CALL	runtime·notok(SB)
	RET

TEXT runtime·write(SB),7,$-8
	MOVL	8(SP), DI		// arg 1 fd
	MOVQ	16(SP), SI		// arg 2 buf
	MOVL	24(SP), DX		// arg 3 count
	MOVL	$4, AX
	SYSCALL
	RET

TEXT runtime·raisesigpipe(SB),7,$16
	// thr_self(&8(SP))
	LEAQ	8(SP), DI	// arg 1 &8(SP)
	MOVL	$432, AX
	SYSCALL
	// thr_kill(self, SIGPIPE)
	MOVQ	8(SP), DI	// arg 1 id
	MOVQ	$13, SI	// arg 2 SIGPIPE
	MOVL	$433, AX
	SYSCALL
	RET

TEXT runtime·setitimer(SB), 7, $-8
	MOVL	8(SP), DI
	MOVQ	16(SP), SI
	MOVQ	24(SP), DX
	MOVL	$83, AX
	SYSCALL
	RET

// func now() (sec int64, nsec int32)
TEXT time·now(SB), 7, $32
	MOVL	$116, AX
	LEAQ	8(SP), DI
	MOVQ	$0, SI
	SYSCALL
	MOVQ	8(SP), AX	// sec
	MOVL	16(SP), DX	// usec

	// sec is in AX, usec in DX
	MOVQ	AX, sec+0(FP)
	IMULQ	$1000, DX
	MOVL	DX, nsec+8(FP)
	RET

TEXT runtime·nanotime(SB), 7, $32
	MOVL	$116, AX
	LEAQ	8(SP), DI
	MOVQ	$0, SI
	SYSCALL
	MOVQ	8(SP), AX	// sec
	MOVL	16(SP), DX	// usec

	// sec is in AX, usec in DX
	// return nsec in AX
	IMULQ	$1000000000, AX
	IMULQ	$1000, DX
	ADDQ	DX, AX
	RET

TEXT runtime·sigaction(SB),7,$-8
	MOVL	8(SP), DI		// arg 1 sig
	MOVQ	16(SP), SI		// arg 2 act
	MOVQ	24(SP), DX		// arg 3 oact
	MOVL	$416, AX
	SYSCALL
	JCC	2(PC)
	CALL	runtime·notok(SB)
	RET

TEXT runtime·sigtramp(SB),7,$64
	get_tls(BX)
	
	// save g
	MOVQ	g(BX), R10
	MOVQ	R10, 40(SP)
	
	// g = m->signal
	MOVQ	m(BX), BP
	MOVQ	m_gsignal(BP), BP
	MOVQ	BP, g(BX)
	
	MOVQ	DI, 0(SP)
	MOVQ	SI, 8(SP)
	MOVQ	DX, 16(SP)
	MOVQ	R10, 24(SP)
	
	CALL	runtime·sighandler(SB)

	// restore g
	get_tls(BX)
	MOVQ	40(SP), R10
	MOVQ	R10, g(BX)
	RET

TEXT runtime·mmap(SB),7,$0
	MOVQ	8(SP), DI		// arg 1 addr
	MOVQ	16(SP), SI		// arg 2 len
	MOVL	24(SP), DX		// arg 3 prot
	MOVL	28(SP), R10		// arg 4 flags
	MOVL	32(SP), R8		// arg 5 fid
	MOVL	36(SP), R9		// arg 6 offset
	MOVL	$477, AX
	SYSCALL
	RET

TEXT runtime·munmap(SB),7,$0
	MOVQ	8(SP), DI		// arg 1 addr
	MOVQ	16(SP), SI		// arg 2 len
	MOVL	$73, AX
	SYSCALL
	JCC	2(PC)
	CALL	runtime·notok(SB)
	RET

TEXT runtime·notok(SB),7,$-8
	MOVL	$0xf1, BP
	MOVQ	BP, (BP)
	RET

TEXT runtime·sigaltstack(SB),7,$-8
	MOVQ	new+8(SP), DI
	MOVQ	old+16(SP), SI
	MOVQ	$53, AX
	SYSCALL
	JCC	2(PC)
	CALL	runtime·notok(SB)
	RET

TEXT runtime·usleep(SB),7,$16
	MOVL	$0, DX
	MOVL	usec+0(FP), AX
	MOVL	$1000000, CX
	DIVL	CX
	MOVQ	AX, 0(SP)		// tv_sec
	MOVL	$1000, AX
	MULL	DX
	MOVQ	AX, 8(SP)		// tv_nsec

	MOVQ	SP, DI			// arg 1 - rqtp
	MOVQ	$0, SI			// arg 2 - rmtp
	MOVL	$240, AX		// sys_nanosleep
	SYSCALL
	JCC	2(PC)
	CALL	runtime·notok(SB)
	RET

// set tls base to DI
TEXT runtime·settls(SB),7,$8
	ADDQ	$16, DI	// adjust for ELF: wants to use -16(FS) and -8(FS) for g and m
	MOVQ	DI, 0(SP)
	MOVQ	SP, SI
	MOVQ	$129, DI	// AMD64_SET_FSBASE
	MOVQ	$165, AX	// sysarch
	SYSCALL
	JCC	2(PC)
	CALL	runtime·notok(SB)
	RET

TEXT runtime·sysctl(SB),7,$0
	MOVQ	8(SP), DI		// arg 1 - name
	MOVL	16(SP), SI		// arg 2 - namelen
	MOVQ	24(SP), DX		// arg 3 - oldp
	MOVQ	32(SP), R10		// arg 4 - oldlenp
	MOVQ	40(SP), R8		// arg 5 - newp
	MOVQ	48(SP), R9		// arg 6 - newlen
	MOVQ	$202, AX		// sys___sysctl
	SYSCALL
	JCC 3(PC)
	NEGL	AX
	RET
	MOVL	$0, AX
	RET

TEXT runtime·osyield(SB),7,$-4
	MOVL	$331, AX		// sys_sched_yield
	INT	$0x80
	RET
