// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// System calls and other sys.stuff for AMD64, FreeBSD
// /usr/src/sys/kern/syscalls.master for syscall numbers.
//

#include "go_asm.h"
#include "go_tls.h"
#include "textflag.h"

TEXT runtime·sys_umtx_op(SB),NOSPLIT,$0
	MOVQ addr+0(FP), DI
	MOVL mode+8(FP), SI
	MOVL val+12(FP), DX
	MOVQ uaddr1+16(FP), R10
	MOVQ ut+24(FP), R8
	MOVL $454, AX
	SYSCALL
	MOVL	AX, ret+32(FP)
	RET

TEXT runtime·thr_new(SB),NOSPLIT,$0
	MOVQ param+0(FP), DI
	MOVL size+8(FP), SI
	MOVL $455, AX
	SYSCALL
	RET

TEXT runtime·thr_start(SB),NOSPLIT,$0
	MOVQ	DI, R13 // m

	// set up FS to point at m->tls
	LEAQ	m_tls(R13), DI
	CALL	runtime·settls(SB)	// smashes DI

	// set up m, g
	get_tls(CX)
	MOVQ	m_g0(R13), DI
	MOVQ	R13, g_m(DI)
	MOVQ	DI, g(CX)

	CALL	runtime·stackcheck(SB)
	CALL	runtime·mstart(SB)

	MOVQ 0, AX			// crash (not reached)

// Exit the entire program (like C exit)
TEXT runtime·exit(SB),NOSPLIT,$-8
	MOVL	code+0(FP), DI		// arg 1 exit status
	MOVL	$1, AX
	SYSCALL
	MOVL	$0xf1, 0xf1  // crash
	RET

// func exitThread(wait *uint32)
TEXT runtime·exitThread(SB),NOSPLIT,$0-8
	MOVQ	wait+0(FP), AX
	// We're done using the stack.
	MOVL	$0, (AX)
	MOVL	$0, DI		// arg 1 long *state
	MOVL	$431, AX	// thr_exit
	SYSCALL
	MOVL	$0xf1, 0xf1  // crash
	JMP	0(PC)

TEXT runtime·open(SB),NOSPLIT,$-8
	MOVQ	name+0(FP), DI		// arg 1 pathname
	MOVL	mode+8(FP), SI		// arg 2 flags
	MOVL	perm+12(FP), DX		// arg 3 mode
	MOVL	$5, AX
	SYSCALL
	JCC	2(PC)
	MOVL	$-1, AX
	MOVL	AX, ret+16(FP)
	RET

TEXT runtime·closefd(SB),NOSPLIT,$-8
	MOVL	fd+0(FP), DI		// arg 1 fd
	MOVL	$6, AX
	SYSCALL
	JCC	2(PC)
	MOVL	$-1, AX
	MOVL	AX, ret+8(FP)
	RET

TEXT runtime·read(SB),NOSPLIT,$-8
	MOVL	fd+0(FP), DI		// arg 1 fd
	MOVQ	p+8(FP), SI		// arg 2 buf
	MOVL	n+16(FP), DX		// arg 3 count
	MOVL	$3, AX
	SYSCALL
	JCC	2(PC)
	MOVL	$-1, AX
	MOVL	AX, ret+24(FP)
	RET

TEXT runtime·write(SB),NOSPLIT,$-8
	MOVQ	fd+0(FP), DI		// arg 1 fd
	MOVQ	p+8(FP), SI		// arg 2 buf
	MOVL	n+16(FP), DX		// arg 3 count
	MOVL	$4, AX
	SYSCALL
	JCC	2(PC)
	MOVL	$-1, AX
	MOVL	AX, ret+24(FP)
	RET

TEXT runtime·getrlimit(SB),NOSPLIT,$-8
	MOVL	kind+0(FP), DI
	MOVQ	limit+8(FP), SI
	MOVL	$194, AX
	SYSCALL
	MOVL	AX, ret+16(FP)
	RET

TEXT runtime·raise(SB),NOSPLIT,$16
	// thr_self(&8(SP))
	LEAQ	8(SP), DI	// arg 1 &8(SP)
	MOVL	$432, AX
	SYSCALL
	// thr_kill(self, SIGPIPE)
	MOVQ	8(SP), DI	// arg 1 id
	MOVL	sig+0(FP), SI	// arg 2
	MOVL	$433, AX
	SYSCALL
	RET

TEXT runtime·raiseproc(SB),NOSPLIT,$0
	// getpid
	MOVL	$20, AX
	SYSCALL
	// kill(self, sig)
	MOVQ	AX, DI		// arg 1 pid
	MOVL	sig+0(FP), SI	// arg 2 sig
	MOVL	$37, AX
	SYSCALL
	RET

TEXT runtime·setitimer(SB), NOSPLIT, $-8
	MOVL	mode+0(FP), DI
	MOVQ	new+8(FP), SI
	MOVQ	old+16(FP), DX
	MOVL	$83, AX
	SYSCALL
	RET

// func walltime() (sec int64, nsec int32)
TEXT runtime·walltime(SB), NOSPLIT, $32
	MOVL	$232, AX // clock_gettime
	MOVQ	$0, DI		// CLOCK_REALTIME
	LEAQ	8(SP), SI
	SYSCALL
	MOVQ	8(SP), AX	// sec
	MOVQ	16(SP), DX	// nsec

	// sec is in AX, nsec in DX
	MOVQ	AX, sec+0(FP)
	MOVL	DX, nsec+8(FP)
	RET

TEXT runtime·nanotime(SB), NOSPLIT, $32
	MOVL	$232, AX
	// We can use CLOCK_MONOTONIC_FAST here when we drop
	// support for FreeBSD 8-STABLE.
	MOVQ	$4, DI		// CLOCK_MONOTONIC
	LEAQ	8(SP), SI
	SYSCALL
	MOVQ	8(SP), AX	// sec
	MOVQ	16(SP), DX	// nsec

	// sec is in AX, nsec in DX
	// return nsec in AX
	IMULQ	$1000000000, AX
	ADDQ	DX, AX
	MOVQ	AX, ret+0(FP)
	RET

TEXT runtime·sigaction(SB),NOSPLIT,$-8
	MOVL	sig+0(FP), DI		// arg 1 sig
	MOVQ	new+8(FP), SI		// arg 2 act
	MOVQ	old+16(FP), DX		// arg 3 oact
	MOVL	$416, AX
	SYSCALL
	JCC	2(PC)
	MOVL	$0xf1, 0xf1  // crash
	RET

TEXT runtime·sigfwd(SB),NOSPLIT,$0-32
	MOVQ	fn+0(FP),    AX
	MOVL	sig+8(FP),   DI
	MOVQ	info+16(FP), SI
	MOVQ	ctx+24(FP),  DX
	PUSHQ	BP
	MOVQ	SP, BP
	ANDQ	$~15, SP     // alignment for x86_64 ABI
	CALL	AX
	MOVQ	BP, SP
	POPQ	BP
	RET

TEXT runtime·sigtramp(SB),NOSPLIT,$72
	// Save callee-saved C registers, since the caller may be a C signal handler.
	MOVQ	BX, bx-8(SP)
	MOVQ	BP, bp-16(SP)  // save in case GOEXPERIMENT=noframepointer is set
	MOVQ	R12, r12-24(SP)
	MOVQ	R13, r13-32(SP)
	MOVQ	R14, r14-40(SP)
	MOVQ	R15, r15-48(SP)
	// We don't save mxcsr or the x87 control word because sigtrampgo doesn't
	// modify them.

	MOVQ	DX, ctx-56(SP)
	MOVQ	SI, info-64(SP)
	MOVQ	DI, signum-72(SP)
	CALL	runtime·sigtrampgo(SB)

	MOVQ	r15-48(SP), R15
	MOVQ	r14-40(SP), R14
	MOVQ	r13-32(SP), R13
	MOVQ	r12-24(SP), R12
	MOVQ	bp-16(SP),  BP
	MOVQ	bx-8(SP),   BX
	RET

TEXT runtime·mmap(SB),NOSPLIT,$0
	MOVQ	addr+0(FP), DI		// arg 1 addr
	MOVQ	n+8(FP), SI		// arg 2 len
	MOVL	prot+16(FP), DX		// arg 3 prot
	MOVL	flags+20(FP), R10		// arg 4 flags
	MOVL	fd+24(FP), R8		// arg 5 fid
	MOVL	off+28(FP), R9		// arg 6 offset
	MOVL	$477, AX
	SYSCALL
	JCC	ok
	MOVQ	$0, p+32(FP)
	MOVQ	AX, err+40(FP)
	RET
ok:
	MOVQ	AX, p+32(FP)
	MOVQ	$0, err+40(FP)
	RET

TEXT runtime·munmap(SB),NOSPLIT,$0
	MOVQ	addr+0(FP), DI		// arg 1 addr
	MOVQ	n+8(FP), SI		// arg 2 len
	MOVL	$73, AX
	SYSCALL
	JCC	2(PC)
	MOVL	$0xf1, 0xf1  // crash
	RET

TEXT runtime·madvise(SB),NOSPLIT,$0
	MOVQ	addr+0(FP), DI
	MOVQ	n+8(FP), SI
	MOVL	flags+16(FP), DX
	MOVQ	$75, AX	// madvise
	SYSCALL
	// ignore failure - maybe pages are locked
	RET
	
TEXT runtime·sigaltstack(SB),NOSPLIT,$-8
	MOVQ	new+0(FP), DI
	MOVQ	old+8(FP), SI
	MOVQ	$53, AX
	SYSCALL
	JCC	2(PC)
	MOVL	$0xf1, 0xf1  // crash
	RET

TEXT runtime·usleep(SB),NOSPLIT,$16
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
	RET

// set tls base to DI
TEXT runtime·settls(SB),NOSPLIT,$8
	ADDQ	$8, DI	// adjust for ELF: wants to use -8(FS) for g and m
	MOVQ	DI, 0(SP)
	MOVQ	SP, SI
	MOVQ	$129, DI	// AMD64_SET_FSBASE
	MOVQ	$165, AX	// sysarch
	SYSCALL
	JCC	2(PC)
	MOVL	$0xf1, 0xf1  // crash
	RET

TEXT runtime·sysctl(SB),NOSPLIT,$0
	MOVQ	mib+0(FP), DI		// arg 1 - name
	MOVL	miblen+8(FP), SI		// arg 2 - namelen
	MOVQ	out+16(FP), DX		// arg 3 - oldp
	MOVQ	size+24(FP), R10		// arg 4 - oldlenp
	MOVQ	dst+32(FP), R8		// arg 5 - newp
	MOVQ	ndst+40(FP), R9		// arg 6 - newlen
	MOVQ	$202, AX		// sys___sysctl
	SYSCALL
	JCC 4(PC)
	NEGQ	AX
	MOVL	AX, ret+48(FP)
	RET
	MOVL	$0, AX
	MOVL	AX, ret+48(FP)
	RET

TEXT runtime·osyield(SB),NOSPLIT,$-4
	MOVL	$331, AX		// sys_sched_yield
	SYSCALL
	RET

TEXT runtime·sigprocmask(SB),NOSPLIT,$0
	MOVL	how+0(FP), DI		// arg 1 - how
	MOVQ	new+8(FP), SI		// arg 2 - set
	MOVQ	old+16(FP), DX		// arg 3 - oset
	MOVL	$340, AX		// sys_sigprocmask
	SYSCALL
	JAE	2(PC)
	MOVL	$0xf1, 0xf1  // crash
	RET

// int32 runtime·kqueue(void);
TEXT runtime·kqueue(SB),NOSPLIT,$0
	MOVQ	$0, DI
	MOVQ	$0, SI
	MOVQ	$0, DX
	MOVL	$362, AX
	SYSCALL
	JCC	2(PC)
	NEGQ	AX
	MOVL	AX, ret+0(FP)
	RET

// int32 runtime·kevent(int kq, Kevent *changelist, int nchanges, Kevent *eventlist, int nevents, Timespec *timeout);
TEXT runtime·kevent(SB),NOSPLIT,$0
	MOVL	kq+0(FP), DI
	MOVQ	ch+8(FP), SI
	MOVL	nch+16(FP), DX
	MOVQ	ev+24(FP), R10
	MOVL	nev+32(FP), R8
	MOVQ	ts+40(FP), R9
	MOVL	$363, AX
	SYSCALL
	JCC	2(PC)
	NEGQ	AX
	MOVL	AX, ret+48(FP)
	RET

// void runtime·closeonexec(int32 fd);
TEXT runtime·closeonexec(SB),NOSPLIT,$0
	MOVL	fd+0(FP), DI	// fd
	MOVQ	$2, SI		// F_SETFD
	MOVQ	$1, DX		// FD_CLOEXEC
	MOVL	$92, AX		// fcntl
	SYSCALL
	RET

// func cpuset_getaffinity(level int, which int, id int64, size int, mask *byte) int32
TEXT runtime·cpuset_getaffinity(SB), NOSPLIT, $0-44
	MOVQ	level+0(FP), DI
	MOVQ	which+8(FP), SI
	MOVQ	id+16(FP), DX
	MOVQ	size+24(FP), R10
	MOVQ	mask+32(FP), R8
	MOVL	$487, AX
	SYSCALL
	JCC	2(PC)
	NEGQ	AX
	MOVL	AX, ret+40(FP)
	RET
