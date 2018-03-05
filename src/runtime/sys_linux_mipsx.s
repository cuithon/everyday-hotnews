// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux
// +build mips mipsle

//
// System calls and other sys.stuff for mips, Linux
//

#include "go_asm.h"
#include "go_tls.h"
#include "textflag.h"

#define SYS_exit		4001
#define SYS_read		4003
#define SYS_write		4004
#define SYS_open		4005
#define SYS_close		4006
#define SYS_getpid		4020
#define SYS_kill		4037
#define SYS_brk			4045
#define SYS_fcntl		4055
#define SYS_mmap		4090
#define SYS_munmap		4091
#define SYS_setitimer		4104
#define SYS_clone		4120
#define SYS_newselect		4142
#define SYS_sched_yield		4162
#define SYS_rt_sigreturn	4193
#define SYS_rt_sigaction	4194
#define SYS_rt_sigprocmask	4195
#define SYS_sigaltstack		4206
#define SYS_madvise		4218
#define SYS_mincore		4217
#define SYS_gettid		4222
#define SYS_tkill		4236
#define SYS_futex		4238
#define SYS_sched_getaffinity	4240
#define SYS_exit_group		4246
#define SYS_epoll_create	4248
#define SYS_epoll_ctl		4249
#define SYS_epoll_wait		4250
#define SYS_clock_gettime	4263
#define SYS_epoll_create1	4326

TEXT runtime·exit(SB),NOSPLIT,$0-4
	MOVW	code+0(FP), R4
	MOVW	$SYS_exit_group, R2
	SYSCALL
	UNDEF
	RET

// func exitThread(wait *uint32)
TEXT runtime·exitThread(SB),NOSPLIT,$0-4
	MOVW	wait+0(FP), R1
	// We're done using the stack.
	MOVW	$0, R2
	SYNC
	MOVW	R2, (R1)
	SYNC
	MOVW	$0, R4	// exit code
	MOVW	$SYS_exit, R2
	SYSCALL
	UNDEF
	JMP	0(PC)

TEXT runtime·open(SB),NOSPLIT,$0-16
	MOVW	name+0(FP), R4
	MOVW	mode+4(FP), R5
	MOVW	perm+8(FP), R6
	MOVW	$SYS_open, R2
	SYSCALL
	BEQ	R7, 2(PC)
	MOVW	$-1, R2
	MOVW	R2, ret+12(FP)
	RET

TEXT runtime·closefd(SB),NOSPLIT,$0-8
	MOVW	fd+0(FP), R4
	MOVW	$SYS_close, R2
	SYSCALL
	BEQ	R7, 2(PC)
	MOVW	$-1, R2
	MOVW	R2, ret+4(FP)
	RET

TEXT runtime·write(SB),NOSPLIT,$0-16
	MOVW	fd+0(FP), R4
	MOVW	p+4(FP), R5
	MOVW	n+8(FP), R6
	MOVW	$SYS_write, R2
	SYSCALL
	BEQ	R7, 2(PC)
	MOVW	$-1, R2
	MOVW	R2, ret+12(FP)
	RET

TEXT runtime·read(SB),NOSPLIT,$0-16
	MOVW	fd+0(FP), R4
	MOVW	p+4(FP), R5
	MOVW	n+8(FP), R6
	MOVW	$SYS_read, R2
	SYSCALL
	BEQ	R7, 2(PC)
	MOVW	$-1, R2
	MOVW	R2, ret+12(FP)
	RET

TEXT runtime·usleep(SB),NOSPLIT,$28-4
	MOVW	usec+0(FP), R3
	MOVW	R3, R5
	MOVW	$1000000, R4
	DIVU	R4, R3
	MOVW	LO, R3
	MOVW	R3, 24(R29)
	MULU	R3, R4
	MOVW	LO, R4
	SUBU	R4, R5
	MOVW	R5, 28(R29)

	// select(0, 0, 0, 0, &tv)
	MOVW	$0, R4
	MOVW	$0, R5
	MOVW	$0, R6
	MOVW	$0, R7
	ADDU	$24, R29, R8
	MOVW	R8, 16(R29)
	MOVW	$SYS_newselect, R2
	SYSCALL
	RET

TEXT runtime·gettid(SB),NOSPLIT,$0-4
	MOVW	$SYS_gettid, R2
	SYSCALL
	MOVW	R2, ret+0(FP)
	RET

TEXT runtime·raise(SB),NOSPLIT,$0-4
	MOVW	$SYS_gettid, R2
	SYSCALL
	MOVW	R2, R4	// arg 1 tid
	MOVW	sig+0(FP), R5	// arg 2
	MOVW	$SYS_tkill, R2
	SYSCALL
	RET

TEXT runtime·raiseproc(SB),NOSPLIT,$0
	MOVW	$SYS_getpid, R2
	SYSCALL
	MOVW	R2, R4	// arg 1 pid
	MOVW	sig+0(FP), R5	// arg 2
	MOVW	$SYS_kill, R2
	SYSCALL
	RET

TEXT runtime·setitimer(SB),NOSPLIT,$0-12
	MOVW	mode+0(FP), R4
	MOVW	new+4(FP), R5
	MOVW	old+8(FP), R6
	MOVW	$SYS_setitimer, R2
	SYSCALL
	RET

TEXT runtime·mincore(SB),NOSPLIT,$0-16
	MOVW	addr+0(FP), R4
	MOVW	n+4(FP), R5
	MOVW	dst+8(FP), R6
	MOVW	$SYS_mincore, R2
	SYSCALL
	SUBU	R2, R0, R2	// caller expects negative errno
	MOVW	R2, ret+12(FP)
	RET

// func walltime() (sec int64, nsec int32)
TEXT runtime·walltime(SB),NOSPLIT,$8-12
	MOVW	$0, R4	// CLOCK_REALTIME
	MOVW	$4(R29), R5
	MOVW	$SYS_clock_gettime, R2
	SYSCALL
	MOVW	4(R29), R3	// sec
	MOVW	8(R29), R5	// nsec
#ifdef GOARCH_mips
	MOVW	R3, sec_lo+4(FP)
	MOVW	R0, sec_hi+0(FP)
#else
	MOVW	R3, sec_lo+0(FP)
	MOVW	R0, sec_hi+4(FP)
#endif
	MOVW	R5, nsec+8(FP)
	RET

TEXT runtime·nanotime(SB),NOSPLIT,$8-8
	MOVW	$1, R4	// CLOCK_MONOTONIC
	MOVW	$4(R29), R5
	MOVW	$SYS_clock_gettime, R2
	SYSCALL
	MOVW	4(R29), R3	// sec
	MOVW	8(R29), R5	// nsec
	// sec is in R3, nsec in R5
	// return nsec in R3
	MOVW	$1000000000, R4
	MULU	R4, R3
	MOVW	LO, R3
	ADDU	R5, R3
	SGTU	R5, R3, R4
#ifdef GOARCH_mips
	MOVW	R3, ret_lo+4(FP)
#else
	MOVW	R3, ret_lo+0(FP)
#endif
	MOVW	HI, R3
	ADDU	R4, R3
#ifdef GOARCH_mips
	MOVW	R3, ret_hi+0(FP)
#else
	MOVW	R3, ret_hi+4(FP)
#endif
	RET

TEXT runtime·rtsigprocmask(SB),NOSPLIT,$0-16
	MOVW	how+0(FP), R4
	MOVW	new+4(FP), R5
	MOVW	old+8(FP), R6
	MOVW	size+12(FP), R7
	MOVW	$SYS_rt_sigprocmask, R2
	SYSCALL
	BEQ	R7, 2(PC)
	UNDEF	// crash
	RET

TEXT runtime·rt_sigaction(SB),NOSPLIT,$0-20
	MOVW	sig+0(FP), R4
	MOVW	new+4(FP), R5
	MOVW	old+8(FP), R6
	MOVW	size+12(FP), R7
	MOVW	$SYS_rt_sigaction, R2
	SYSCALL
	MOVW	R2, ret+16(FP)
	RET

TEXT runtime·sigfwd(SB),NOSPLIT,$0-16
	MOVW	sig+4(FP), R4
	MOVW	info+8(FP), R5
	MOVW	ctx+12(FP), R6
	MOVW	fn+0(FP), R25
	MOVW	R29, R22
	SUBU	$16, R29
	AND	$~7, R29	// shadow space for 4 args aligned to 8 bytes as per O32 ABI
	JAL	(R25)
	MOVW	R22, R29
	RET

TEXT runtime·sigtramp(SB),NOSPLIT,$12
	// this might be called in external code context,
	// where g is not set.
	MOVB	runtime·iscgo(SB), R1
	BEQ	R1, 2(PC)
	JAL	runtime·load_g(SB)

	MOVW	R4, 4(R29)
	MOVW	R5, 8(R29)
	MOVW	R6, 12(R29)
	MOVW	$runtime·sigtrampgo(SB), R1
	JAL	(R1)
	RET

TEXT runtime·cgoSigtramp(SB),NOSPLIT,$0
	JMP	runtime·sigtramp(SB)

TEXT runtime·mmap(SB),NOSPLIT,$20-32
	MOVW	addr+0(FP), R4
	MOVW	n+4(FP), R5
	MOVW	prot+8(FP), R6
	MOVW	flags+12(FP), R7
	MOVW	fd+16(FP), R8
	MOVW	off+20(FP), R9
	MOVW	R8, 16(R29)
	MOVW	R9, 20(R29)

	MOVW	$SYS_mmap, R2
	SYSCALL
	BEQ	R7, ok
	MOVW	$0, p+24(FP)
	MOVW	R2, err+28(FP)
	RET
ok:
	MOVW	R2, p+24(FP)
	MOVW	$0, err+28(FP)
	RET

TEXT runtime·munmap(SB),NOSPLIT,$0-8
	MOVW	addr+0(FP), R4
	MOVW	n+4(FP), R5
	MOVW	$SYS_munmap, R2
	SYSCALL
	BEQ	R7, 2(PC)
	UNDEF	// crash
	RET

TEXT runtime·madvise(SB),NOSPLIT,$0-12
	MOVW	addr+0(FP), R4
	MOVW	n+4(FP), R5
	MOVW	flags+8(FP), R6
	MOVW	$SYS_madvise, R2
	SYSCALL
	// ignore failure - maybe pages are locked
	RET

// int32 futex(int32 *uaddr, int32 op, int32 val, struct timespec *timeout, int32 *uaddr2, int32 val2);
TEXT runtime·futex(SB),NOSPLIT,$20-28
	MOVW	addr+0(FP), R4
	MOVW	op+4(FP), R5
	MOVW	val+8(FP), R6
	MOVW	ts+12(FP), R7

	MOVW	addr2+16(FP), R8
	MOVW	val3+20(FP), R9

	MOVW	R8, 16(R29)
	MOVW	R9, 20(R29)

	MOVW	$SYS_futex, R2
	SYSCALL
	MOVW	R2, ret+24(FP)
	RET


// int32 clone(int32 flags, void *stk, M *mp, G *gp, void (*fn)(void));
TEXT runtime·clone(SB),NOSPLIT|NOFRAME,$0-24
	MOVW	flags+0(FP), R4
	MOVW	stk+4(FP), R5
	MOVW	R0, R6	// ptid
	MOVW	R0, R7	// tls

	// O32 syscall handler unconditionally copies arguments 5-8 from stack,
	// even for syscalls with less than 8 arguments. Reserve 32 bytes of new
	// stack so that any syscall invoked immediately in the new thread won't fail.
	ADD	$-32, R5

	// Copy mp, gp, fn off parent stack for use by child.
	MOVW	mp+8(FP), R16
	MOVW	gp+12(FP), R17
	MOVW	fn+16(FP), R18

	MOVW	$1234, R1

	MOVW	R16, 0(R5)
	MOVW	R17, 4(R5)
	MOVW	R18, 8(R5)

	MOVW	R1, 12(R5)

	MOVW	$SYS_clone, R2
	SYSCALL

	// In parent, return.
	BEQ	R2, 5(PC)
	SUBU	R2, R0, R3
	CMOVN	R7, R3, R2
	MOVW	R2, ret+20(FP)
	RET

	// In child, on new stack.
	// Check that SP is as we expect
	MOVW	12(R29), R16
	MOVW	$1234, R1
	BEQ	R16, R1, 2(PC)
	MOVW	(R0), R0

	// Initialize m->procid to Linux tid
	MOVW	$SYS_gettid, R2
	SYSCALL

	MOVW	0(R29), R16	// m
	MOVW	4(R29), R17	// g
	MOVW	8(R29), R18	// fn

	BEQ	R16, nog
	BEQ	R17, nog

	MOVW	R2, m_procid(R16)

	// In child, set up new stack
	MOVW	R16, g_m(R17)
	MOVW	R17, g

// TODO(mips32): doesn't have runtime·stackcheck(SB)

nog:
	// Call fn
	ADDU	$32, R29
	JAL	(R18)

	// It shouldn't return.	 If it does, exit that thread.
	ADDU	$-32, R29
	MOVW	$0xf4, R4
	MOVW	$SYS_exit, R2
	SYSCALL
	UNDEF

TEXT runtime·sigaltstack(SB),NOSPLIT,$0
	MOVW	new+0(FP), R4
	MOVW	old+4(FP), R5
	MOVW	$SYS_sigaltstack, R2
	SYSCALL
	BEQ	R7, 2(PC)
	UNDEF	// crash
	RET

TEXT runtime·osyield(SB),NOSPLIT,$0
	MOVW	$SYS_sched_yield, R2
	SYSCALL
	RET

TEXT runtime·sched_getaffinity(SB),NOSPLIT,$0-16
	MOVW	pid+0(FP), R4
	MOVW	len+4(FP), R5
	MOVW	buf+8(FP), R6
	MOVW	$SYS_sched_getaffinity, R2
	SYSCALL
	MOVW	R2, ret+12(FP)
	RET

// int32 runtime·epollcreate(int32 size);
TEXT runtime·epollcreate(SB),NOSPLIT,$0-8
	MOVW	size+0(FP), R4
	MOVW	$SYS_epoll_create, R2
	SYSCALL
	MOVW	R2, ret+4(FP)
	RET

// int32 runtime·epollcreate1(int32 flags);
TEXT runtime·epollcreate1(SB),NOSPLIT,$0-8
	MOVW	flags+0(FP), R4
	MOVW	$SYS_epoll_create1, R2
	SYSCALL
	MOVW	R2, ret+4(FP)
	RET

// func epollctl(epfd, op, fd int32, ev *epollEvent) int
TEXT runtime·epollctl(SB),NOSPLIT,$0-20
	MOVW	epfd+0(FP), R4
	MOVW	op+4(FP), R5
	MOVW	fd+8(FP), R6
	MOVW	ev+12(FP), R7
	MOVW	$SYS_epoll_ctl, R2
	SYSCALL
	MOVW	R2, ret+16(FP)
	RET

// int32 runtime·epollwait(int32 epfd, EpollEvent *ev, int32 nev, int32 timeout);
TEXT runtime·epollwait(SB),NOSPLIT,$0-20
	MOVW	epfd+0(FP), R4
	MOVW	ev+4(FP), R5
	MOVW	nev+8(FP), R6
	MOVW	timeout+12(FP), R7
	MOVW	$SYS_epoll_wait, R2
	SYSCALL
	MOVW	R2, ret+16(FP)
	RET

// void runtime·closeonexec(int32 fd);
TEXT runtime·closeonexec(SB),NOSPLIT,$0-4
	MOVW	fd+0(FP), R4	// fd
	MOVW	$2, R5	// F_SETFD
	MOVW	$1, R6	// FD_CLOEXEC
	MOVW	$SYS_fcntl, R2
	SYSCALL
	RET

// func sbrk0() uintptr
TEXT runtime·sbrk0(SB),NOSPLIT,$0-4
	// Implemented as brk(NULL).
	MOVW	$0, R4
	MOVW	$SYS_brk, R2
	SYSCALL
	MOVW	R2, ret+0(FP)
	RET
