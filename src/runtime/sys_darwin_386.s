// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// System calls and other sys.stuff for 386, Darwin
// See http://fxr.watson.org/fxr/source/bsd/kern/syscalls.c?v=xnu-1228
// or /usr/include/sys/syscall.h (on a Mac) for system call numbers.

#include "go_asm.h"
#include "go_tls.h"
#include "textflag.h"

// Exit the entire program (like C exit)
TEXT runtime·exit(SB),NOSPLIT,$0
	MOVL	$1, AX
	INT	$0x80
	MOVL	$0xf1, 0xf1  // crash
	RET

// Exit this OS thread (like pthread_exit, which eventually
// calls __bsdthread_terminate).
TEXT exit1<>(SB),NOSPLIT,$16-0
	// __bsdthread_terminate takes 4 word-size arguments.
	// Set them all to 0. (None are an exit status.)
	MOVL	$0, 0(SP)
	MOVL	$0, 4(SP)
	MOVL	$0, 8(SP)
	MOVL	$0, 12(SP)
	MOVL	$361, AX
	INT	$0x80
	JAE 2(PC)
	MOVL	$0xf1, 0xf1  // crash
	RET

GLOBL exitStack<>(SB),RODATA,$(4*4)
DATA exitStack<>+0x00(SB)/4, $0
DATA exitStack<>+0x04(SB)/4, $0
DATA exitStack<>+0x08(SB)/4, $0
DATA exitStack<>+0x0c(SB)/4, $0

// func exitThread(wait *uint32)
TEXT runtime·exitThread(SB),NOSPLIT,$0-4
	MOVL	wait+0(FP), AX
	// We're done using the stack.
	MOVL	$0, (AX)
	// __bsdthread_terminate takes 4 arguments, which it expects
	// on the stack. They should all be 0, so switch over to a
	// fake stack of 0s. It won't write to the stack.
	MOVL	$exitStack<>(SB), SP
	MOVL	$361, AX	// __bsdthread_terminate
	INT	$0x80
	MOVL	$0xf1, 0xf1  // crash
	JMP	0(PC)

TEXT runtime·open(SB),NOSPLIT,$0
	MOVL	$5, AX
	INT	$0x80
	JAE	2(PC)
	MOVL	$-1, AX
	MOVL	AX, ret+12(FP)
	RET

TEXT runtime·closefd(SB),NOSPLIT,$0
	MOVL	$6, AX
	INT	$0x80
	JAE	2(PC)
	MOVL	$-1, AX
	MOVL	AX, ret+4(FP)
	RET

TEXT runtime·read(SB),NOSPLIT,$0
	MOVL	$3, AX
	INT	$0x80
	JAE	2(PC)
	MOVL	$-1, AX
	MOVL	AX, ret+12(FP)
	RET

TEXT runtime·write(SB),NOSPLIT,$0
	MOVL	$4, AX
	INT	$0x80
	JAE	2(PC)
	MOVL	$-1, AX
	MOVL	AX, ret+12(FP)
	RET

TEXT runtime·raise(SB),NOSPLIT,$0
	// Ideally we'd send the signal to the current thread,
	// not the whole process, but that's too hard on OS X.
	JMP	runtime·raiseproc(SB)

TEXT runtime·raiseproc(SB),NOSPLIT,$16
	MOVL	$20, AX // getpid
	INT	$0x80
	MOVL	AX, 4(SP)	// pid
	MOVL	sig+0(FP), AX
	MOVL	AX, 8(SP)	// signal
	MOVL	$1, 12(SP)	// posix
	MOVL	$37, AX // kill
	INT	$0x80
	RET

TEXT runtime·mmap(SB),NOSPLIT,$0
	MOVL	$197, AX
	INT	$0x80
	JAE	ok
	MOVL	$0, p+24(FP)
	MOVL	AX, err+28(FP)
	RET
ok:
	MOVL	AX, p+24(FP)
	MOVL	$0, err+28(FP)
	RET

TEXT runtime·madvise(SB),NOSPLIT,$0
	MOVL	$75, AX
	INT	$0x80
	// ignore failure - maybe pages are locked
	RET

TEXT runtime·munmap(SB),NOSPLIT,$0
	MOVL	$73, AX
	INT	$0x80
	JAE	2(PC)
	MOVL	$0xf1, 0xf1  // crash
	RET

TEXT runtime·setitimer(SB),NOSPLIT,$0
	MOVL	$83, AX
	INT	$0x80
	RET

// OS X comm page time offsets
// http://www.opensource.apple.com/source/xnu/xnu-1699.26.8/osfmk/i386/cpu_capabilities.h
#define	cpu_capabilities	0x20
#define	nt_tsc_base	0x50
#define	nt_scale	0x58
#define	nt_shift	0x5c
#define	nt_ns_base	0x60
#define	nt_generation	0x68
#define	gtod_generation	0x6c
#define	gtod_ns_base	0x70
#define	gtod_sec_base	0x78

// called from assembly
// 64-bit unix nanoseconds returned in DX:AX.
// I'd much rather write this in C but we need
// assembly for the 96-bit multiply and RDTSC.
//
// Note that we could arrange to return monotonic time here
// as well, but we don't bother, for two reasons:
// 1. macOS only supports 64-bit systems, so no one should
// be using the 32-bit code in production.
// This code is only maintained to make it easier for developers
// using Macs to test the 32-bit compiler.
// 2. On some (probably now unsupported) CPUs,
// the code falls back to the system call always,
// so it can't even use the comm page at all. 
TEXT runtime·now(SB),NOSPLIT,$40
	MOVL	$0xffff0000, BP /* comm page base */
	
	// Test for slow CPU. If so, the math is completely
	// different, and unimplemented here, so use the
	// system call.
	MOVL	cpu_capabilities(BP), AX
	TESTL	$0x4000, AX
	JNZ	systime

	// Loop trying to take a consistent snapshot
	// of the time parameters.
timeloop:
	MOVL	gtod_generation(BP), BX
	TESTL	BX, BX
	JZ	systime
	MOVL	nt_generation(BP), CX
	TESTL	CX, CX
	JZ	timeloop
	RDTSC
	MOVL	nt_tsc_base(BP), SI
	MOVL	(nt_tsc_base+4)(BP), DI
	MOVL	SI, 0(SP)
	MOVL	DI, 4(SP)
	MOVL	nt_scale(BP), SI
	MOVL	SI, 8(SP)
	MOVL	nt_ns_base(BP), SI
	MOVL	(nt_ns_base+4)(BP), DI
	MOVL	SI, 12(SP)
	MOVL	DI, 16(SP)
	CMPL	nt_generation(BP), CX
	JNE	timeloop
	MOVL	gtod_ns_base(BP), SI
	MOVL	(gtod_ns_base+4)(BP), DI
	MOVL	SI, 20(SP)
	MOVL	DI, 24(SP)
	MOVL	gtod_sec_base(BP), SI
	MOVL	(gtod_sec_base+4)(BP), DI
	MOVL	SI, 28(SP)
	MOVL	DI, 32(SP)
	CMPL	gtod_generation(BP), BX
	JNE	timeloop

	// Gathered all the data we need. Compute time.
	//	((tsc - nt_tsc_base) * nt_scale) >> 32 + nt_ns_base - gtod_ns_base + gtod_sec_base*1e9
	// The multiply and shift extracts the top 64 bits of the 96-bit product.
	SUBL	0(SP), AX // DX:AX = (tsc - nt_tsc_base)
	SBBL	4(SP), DX

	// We have x = tsc - nt_tsc_base - DX:AX to be
	// multiplied by y = nt_scale = 8(SP), keeping the top 64 bits of the 96-bit product.
	// x*y = (x&0xffffffff)*y + (x&0xffffffff00000000)*y
	// (x*y)>>32 = ((x&0xffffffff)*y)>>32 + (x>>32)*y
	MOVL	DX, CX // SI = (x&0xffffffff)*y >> 32
	MOVL	$0, DX
	MULL	8(SP)
	MOVL	DX, SI

	MOVL	CX, AX // DX:AX = (x>>32)*y
	MOVL	$0, DX
	MULL	8(SP)

	ADDL	SI, AX	// DX:AX += (x&0xffffffff)*y >> 32
	ADCL	$0, DX
	
	// DX:AX is now ((tsc - nt_tsc_base) * nt_scale) >> 32.
	ADDL	12(SP), AX	// DX:AX += nt_ns_base
	ADCL	16(SP), DX
	SUBL	20(SP), AX	// DX:AX -= gtod_ns_base
	SBBL	24(SP), DX
	MOVL	AX, SI	// DI:SI = DX:AX
	MOVL	DX, DI
	MOVL	28(SP), AX	// DX:AX = gtod_sec_base*1e9
	MOVL	32(SP), DX
	MOVL	$1000000000, CX
	MULL	CX
	ADDL	SI, AX	// DX:AX += DI:SI
	ADCL	DI, DX
	RET

systime:
	// Fall back to system call (usually first call in this thread)
	LEAL	16(SP), AX	// must be non-nil, unused
	MOVL	AX, 4(SP)
	MOVL	$0, 8(SP)	// time zone pointer
	MOVL	$0, 12(SP)	// required as of Sierra; Issue 16570
	MOVL	$116, AX // SYS_GETTIMEOFDAY
	INT	$0x80
	CMPL	AX, $0
	JNE	inreg
	MOVL	16(SP), AX
	MOVL	20(SP), DX
inreg:
	// sec is in AX, usec in DX
	// convert to DX:AX nsec
	MOVL	DX, BX
	MOVL	$1000000000, CX
	MULL	CX
	IMULL	$1000, BX
	ADDL	BX, AX
	ADCL	$0, DX
	RET

// func now() (sec int64, nsec int32, mono uint64)
TEXT time·now(SB),NOSPLIT,$0-20
	CALL	runtime·now(SB)
	MOVL	AX, BX
	MOVL	DX, BP
	SUBL	runtime·startNano(SB), BX
	SBBL	runtime·startNano+4(SB), BP
	MOVL	BX, mono+12(FP)
	MOVL	BP, mono+16(FP)
	MOVL	$1000000000, CX
	DIVL	CX
	MOVL	AX, sec+0(FP)
	MOVL	$0, sec+4(FP)
	MOVL	DX, nsec+8(FP)
	RET

// func nanotime() int64
TEXT runtime·nanotime(SB),NOSPLIT,$0
	CALL	runtime·now(SB)
	SUBL	runtime·startNano(SB), AX
	SBBL	runtime·startNano+4(SB), DX
	MOVL	AX, ret_lo+0(FP)
	MOVL	DX, ret_hi+4(FP)
	RET

TEXT runtime·sigprocmask(SB),NOSPLIT,$0
	MOVL	$329, AX  // pthread_sigmask (on OS X, sigprocmask==entire process)
	INT	$0x80
	JAE	2(PC)
	MOVL	$0xf1, 0xf1  // crash
	RET

TEXT runtime·sigaction(SB),NOSPLIT,$0
	MOVL	$46, AX
	INT	$0x80
	JAE	2(PC)
	MOVL	$0xf1, 0xf1  // crash
	RET

TEXT runtime·sigfwd(SB),NOSPLIT,$0-16
	MOVL	fn+0(FP), AX
	MOVL	sig+4(FP), BX
	MOVL	info+8(FP), CX
	MOVL	ctx+12(FP), DX
	MOVL	SP, SI
	SUBL	$32, SP
	ANDL	$~15, SP	// align stack: handler might be a C function
	MOVL	BX, 0(SP)
	MOVL	CX, 4(SP)
	MOVL	DX, 8(SP)
	MOVL	SI, 12(SP)	// save SI: handler might be a Go function
	CALL	AX
	MOVL	12(SP), AX
	MOVL	AX, SP
	RET

// Sigtramp's job is to call the actual signal handler.
// It is called with the following arguments on the stack:
//	0(SP)	"return address" - ignored
//	4(SP)	actual handler
//	8(SP)	siginfo style
//	12(SP)	signal number
//	16(SP)	siginfo
//	20(SP)	context
TEXT runtime·sigtramp(SB),NOSPLIT,$20
	MOVL	sig+8(FP), BX
	MOVL	BX, 0(SP)
	MOVL	info+12(FP), BX
	MOVL	BX, 4(SP)
	MOVL	ctx+16(FP), BX
	MOVL	BX, 8(SP)
	CALL	runtime·sigtrampgo(SB)

	// call sigreturn
	MOVL	ctx+16(FP), CX
	MOVL	infostyle+4(FP), BX
	MOVL	$0, 0(SP)	// "caller PC" - ignored
	MOVL	CX, 4(SP)
	MOVL	BX, 8(SP)
	MOVL	$184, AX	// sigreturn(ucontext, infostyle)
	INT	$0x80
	MOVL	$0xf1, 0xf1  // crash
	RET

TEXT runtime·sigaltstack(SB),NOSPLIT,$0
	MOVL	$53, AX
	INT	$0x80
	JAE	2(PC)
	MOVL	$0xf1, 0xf1  // crash
	RET

TEXT runtime·usleep(SB),NOSPLIT,$32
	MOVL	$0, DX
	MOVL	usec+0(FP), AX
	MOVL	$1000000, CX
	DIVL	CX
	MOVL	AX, 24(SP)  // sec
	MOVL	DX, 28(SP)  // usec

	// select(0, 0, 0, 0, &tv)
	MOVL	$0, 0(SP)  // "return PC" - ignored
	MOVL	$0, 4(SP)
	MOVL	$0, 8(SP)
	MOVL	$0, 12(SP)
	MOVL	$0, 16(SP)
	LEAL	24(SP), AX
	MOVL	AX, 20(SP)
	MOVL	$93, AX
	INT	$0x80
	RET

// Invoke Mach system call.
// Assumes system call number in AX,
// caller PC on stack, caller's caller PC next,
// and then the system call arguments.
//
// Can be used for BSD too, but we don't,
// because if you use this interface the BSD
// system call numbers need an extra field
// in the high 16 bits that seems to be the
// argument count in bytes but is not always.
// INT $0x80 works fine for those.
TEXT runtime·sysenter(SB),NOSPLIT,$0
	POPL	DX
	MOVL	SP, CX
	SYSENTER
	// returns to DX with SP set to CX

TEXT runtime·mach_msg_trap(SB),NOSPLIT,$0
	MOVL	$-31, AX
	CALL	runtime·sysenter(SB)
	MOVL	AX, ret+28(FP)
	RET

TEXT runtime·mach_reply_port(SB),NOSPLIT,$0
	MOVL	$-26, AX
	CALL	runtime·sysenter(SB)
	MOVL	AX, ret+0(FP)
	RET

TEXT runtime·mach_task_self(SB),NOSPLIT,$0
	MOVL	$-28, AX
	CALL	runtime·sysenter(SB)
	MOVL	AX, ret+0(FP)
	RET

// Mach provides trap versions of the semaphore ops,
// instead of requiring the use of RPC.

// func mach_semaphore_wait(sema uint32) int32
TEXT runtime·mach_semaphore_wait(SB),NOSPLIT,$0
	MOVL	$-36, AX
	CALL	runtime·sysenter(SB)
	MOVL	AX, ret+4(FP)
	RET

// func mach_semaphore_timedwait(sema, sec, nsec uint32) int32
TEXT runtime·mach_semaphore_timedwait(SB),NOSPLIT,$0
	MOVL	$-38, AX
	CALL	runtime·sysenter(SB)
	MOVL	AX, ret+12(FP)
	RET

// func mach_semaphore_signal(sema uint32) int32
TEXT runtime·mach_semaphore_signal(SB),NOSPLIT,$0
	MOVL	$-33, AX
	CALL	runtime·sysenter(SB)
	MOVL	AX, ret+4(FP)
	RET

// func mach_semaphore_signal_all(sema uint32) int32
TEXT runtime·mach_semaphore_signal_all(SB),NOSPLIT,$0
	MOVL	$-34, AX
	CALL	runtime·sysenter(SB)
	MOVL	AX, ret+4(FP)
	RET

// func setldt(entry int, address int, limit int)
// entry and limit are ignored.
TEXT runtime·setldt(SB),NOSPLIT,$32
	MOVL	address+4(FP), BX	// aka base

	/*
	 * When linking against the system libraries,
	 * we use its pthread_create and let it set up %gs
	 * for us.  When we do that, the private storage
	 * we get is not at 0(GS) but at 0x18(GS).
	 * The linker rewrites 0(TLS) into 0x18(GS) for us.
	 * To accommodate that rewrite, we translate the
	 * address here so that 0x18(GS) maps to 0(address).
	 *
	 * Constant must match the one in cmd/link/internal/ld/sym.go.
	 */
	SUBL	$0x18, BX

	/*
	 * Must set up as USER_CTHREAD segment because
	 * Darwin forces that value into %gs for signal handlers,
	 * and if we don't set one up, we'll get a recursive
	 * fault trying to get into the signal handler.
	 * Since we have to set one up anyway, it might as
	 * well be the value we want.  So don't bother with
	 * i386_set_ldt.
	 */
	MOVL	BX, 4(SP)
	MOVL	$3, AX	// thread_fast_set_cthread_self - machdep call #3
	INT	$0x82	// sic: 0x82, not 0x80, for machdep call

	XORL	AX, AX
	MOVW	GS, AX
	RET

TEXT runtime·sysctl(SB),NOSPLIT,$0
	MOVL	$202, AX
	INT	$0x80
	JAE	4(PC)
	NEGL	AX
	MOVL	AX, ret+24(FP)
	RET
	MOVL	$0, AX
	MOVL	AX, ret+24(FP)
	RET

// func kqueue() int32
TEXT runtime·kqueue(SB),NOSPLIT,$0
	MOVL	$362, AX
	INT	$0x80
	JAE	2(PC)
	NEGL	AX
	MOVL	AX, ret+0(FP)
	RET

// func kevent(kq int32, ch *keventt, nch int32, ev *keventt, nev int32, ts *timespec) int32
TEXT runtime·kevent(SB),NOSPLIT,$0
	MOVL	$363, AX
	INT	$0x80
	JAE	2(PC)
	NEGL	AX
	MOVL	AX, ret+24(FP)
	RET

// func closeonexec(fd int32)
TEXT runtime·closeonexec(SB),NOSPLIT,$32
	MOVL	$92, AX  // fcntl
	// 0(SP) is where the caller PC would be; kernel skips it
	MOVL	fd+0(FP), BX
	MOVL	BX, 4(SP)  // fd
	MOVL	$2, 8(SP)  // F_SETFD
	MOVL	$1, 12(SP)  // FD_CLOEXEC
	INT	$0x80
	JAE	2(PC)
	NEGL	AX
	RET

// mstart_stub is the first function executed on a new thread started by pthread_create.
// It just does some low-level setup and then calls mstart.
// Note: called with the C calling convention.
TEXT runtime·mstart_stub(SB),NOSPLIT,$0
	// The value at SP+4 points to the m.
	// We are already on m's g0 stack.

	MOVL	SP, AX       // hide argument read from vet (vet thinks this function is using the Go calling convention)
	MOVL	4(AX), DI    // m
	MOVL	m_g0(DI), DX // g

	// Initialize TLS entry.
	// See cmd/link/internal/ld/sym.go:computeTLSOffset.
	MOVL	DX, 0x18(GS)

	// Someday the convention will be D is always cleared.
	CLD

	CALL	runtime·stackcheck(SB) // just in case
	CALL	runtime·mstart(SB)

	// mstart shouldn't ever return, and if it does, we shouldn't ever join to this thread
	// to get its return status. But tell pthread everything is ok, just in case.
	XORL	AX, AX
	RET

TEXT runtime·pthread_attr_init_trampoline(SB),NOSPLIT,$0-8
	// move args into registers
	MOVL	attr+0(FP), AX

	// save SP, BP
	PUSHL	BP
	MOVL	SP, BP

	// allocate space for args
	SUBL	$4, SP

	// align stack to 16 bytes
	ANDL	$~15, SP

	// call libc function
	MOVL	AX, 0(SP)
	CALL	libc_pthread_attr_init(SB)

	// restore BP, SP
	MOVL	BP, SP
	POPL	BP

	// save result.
	MOVL	AX, ret+4(FP)
	RET

TEXT runtime·pthread_attr_setstack_trampoline(SB),NOSPLIT,$0-16
	MOVL	attr+0(FP), AX
	MOVL	addr+4(FP), CX
	MOVL	size+8(FP), DX

	PUSHL	BP
	MOVL	SP, BP

	SUBL	$12, SP
	ANDL	$~15, SP

	MOVL	AX, 0(SP)
	MOVL	CX, 4(SP)
	MOVL	DX, 8(SP)
	CALL	libc_pthread_attr_setstack(SB)

	MOVL	BP, SP
	POPL	BP

	MOVL	AX, ret+12(FP)
	RET

TEXT runtime·pthread_attr_setdetachstate_trampoline(SB),NOSPLIT,$0-12
	MOVL	attr+0(FP), AX
	MOVL	state+4(FP), CX

	PUSHL	BP
	MOVL	SP, BP

	SUBL	$8, SP
	ANDL	$~15, SP

	MOVL	AX, 0(SP)
	MOVL	CX, 4(SP)
	CALL	libc_pthread_attr_setdetachstate(SB)

	MOVL	BP, SP
	POPL	BP

	MOVL	AX, ret+8(FP)
	RET

TEXT runtime·pthread_create_trampoline(SB),NOSPLIT,$0-20
	MOVL	t+0(FP), AX
	MOVL	attr+4(FP), CX
	MOVL	start+8(FP), DX
	MOVL	arg+12(FP), BX

	PUSHL	BP
	MOVL	SP, BP

	SUBL	$16, SP
	ANDL	$~15, SP

	MOVL	AX, 0(SP)
	MOVL	CX, 4(SP)
	MOVL	DX, 8(SP)
	MOVL	BX, 12(SP)
	CALL	libc_pthread_create(SB)

	MOVL	BP, SP
	POPL	BP

	MOVL	AX, ret+16(FP)
	RET
