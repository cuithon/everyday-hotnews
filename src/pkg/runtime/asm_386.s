// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "zasm_GOOS_GOARCH.h"

TEXT _rt0_go(SB),7,$0
	// copy arguments forward on an even stack
	MOVL	argc+0(FP), AX
	MOVL	argv+4(FP), BX
	SUBL	$128, SP		// plenty of scratch
	ANDL	$~15, SP
	MOVL	AX, 120(SP)		// save argc, argv away
	MOVL	BX, 124(SP)

	// set default stack bounds.
	// _cgo_init may update stackguard.
	MOVL	$runtime·g0(SB), BP
	LEAL	(-64*1024+104)(SP), BX
	MOVL	BX, g_stackguard(BP)
	MOVL	BX, g_stackguard0(BP)
	MOVL	SP, g_stackbase(BP)
	
	// find out information about the processor we're on
	MOVL	$0, AX
	CPUID
	CMPL	AX, $0
	JE	nocpuinfo
	MOVL	$1, AX
	CPUID
	MOVL	CX, runtime·cpuid_ecx(SB)
	MOVL	DX, runtime·cpuid_edx(SB)
nocpuinfo:	

	// if there is an _cgo_init, call it to let it
	// initialize and to set up GS.  if not,
	// we set up GS ourselves.
	MOVL	_cgo_init(SB), AX
	TESTL	AX, AX
	JZ	needtls
	MOVL	$setmg_gcc<>(SB), BX
	MOVL	BX, 4(SP)
	MOVL	BP, 0(SP)
	CALL	AX
	// update stackguard after _cgo_init
	MOVL	$runtime·g0(SB), CX
	MOVL	g_stackguard0(CX), AX
	MOVL	AX, g_stackguard(CX)
	// skip runtime·ldt0setup(SB) and tls test after _cgo_init for non-windows
	CMPL runtime·iswindows(SB), $0
	JEQ ok
needtls:
	// skip runtime·ldt0setup(SB) and tls test on Plan 9 in all cases
	CMPL	runtime·isplan9(SB), $1
	JEQ	ok

	// set up %gs
	CALL	runtime·ldt0setup(SB)

	// store through it, to make sure it works
	get_tls(BX)
	MOVL	$0x123, g(BX)
	MOVL	runtime·tls0(SB), AX
	CMPL	AX, $0x123
	JEQ	ok
	MOVL	AX, 0	// abort
ok:
	// set up m and g "registers"
	get_tls(BX)
	LEAL	runtime·g0(SB), CX
	MOVL	CX, g(BX)
	LEAL	runtime·m0(SB), AX
	MOVL	AX, m(BX)

	// save m->g0 = g0
	MOVL	CX, m_g0(AX)

	CALL	runtime·emptyfunc(SB)	// fault if stack check is wrong

	// convention is D is always cleared
	CLD

	CALL	runtime·check(SB)

	// saved argc, argv
	MOVL	120(SP), AX
	MOVL	AX, 0(SP)
	MOVL	124(SP), AX
	MOVL	AX, 4(SP)
	CALL	runtime·args(SB)
	CALL	runtime·osinit(SB)
	CALL	runtime·hashinit(SB)
	CALL	runtime·schedinit(SB)

	// create a new goroutine to start program
	PUSHL	$runtime·main·f(SB)	// entry
	PUSHL	$0	// arg size
	CALL	runtime·newproc(SB)
	POPL	AX
	POPL	AX

	// start this M
	CALL	runtime·mstart(SB)

	INT $3
	RET

DATA	runtime·main·f+0(SB)/4,$runtime·main(SB)
GLOBL	runtime·main·f(SB),8,$4

TEXT runtime·breakpoint(SB),7,$0
	INT $3
	RET

TEXT runtime·asminit(SB),7,$0
	// Linux and MinGW start the FPU in extended double precision.
	// Other operating systems use double precision.
	// Change to double precision to match them,
	// and to match other hardware that only has double.
	PUSHL $0x27F
	FLDCW	0(SP)
	POPL AX
	RET

/*
 *  go-routine
 */

// void gosave(Gobuf*)
// save state in Gobuf; setjmp
TEXT runtime·gosave(SB), 7, $0
	MOVL	4(SP), AX		// gobuf
	LEAL	4(SP), BX		// caller's SP
	MOVL	BX, gobuf_sp(AX)
	MOVL	0(SP), BX		// caller's PC
	MOVL	BX, gobuf_pc(AX)
	MOVL	$0, gobuf_ret(AX)
	MOVL	$0, gobuf_ctxt(AX)
	get_tls(CX)
	MOVL	g(CX), BX
	MOVL	BX, gobuf_g(AX)
	RET

// void gogo(Gobuf*)
// restore state from Gobuf; longjmp
TEXT runtime·gogo(SB), 7, $0
	MOVL	4(SP), BX		// gobuf
	MOVL	gobuf_g(BX), DX
	MOVL	0(DX), CX		// make sure g != nil
	get_tls(CX)
	MOVL	DX, g(CX)
	MOVL	gobuf_sp(BX), SP	// restore SP
	MOVL	gobuf_ret(BX), AX
	MOVL	gobuf_ctxt(BX), DX
	MOVL	$0, gobuf_sp(BX)	// clear to help garbage collector
	MOVL	$0, gobuf_ret(BX)
	MOVL	$0, gobuf_ctxt(BX)
	MOVL	gobuf_pc(BX), BX
	JMP	BX

// void mcall(void (*fn)(G*))
// Switch to m->g0's stack, call fn(g).
// Fn must never return.  It should gogo(&g->sched)
// to keep running g.
TEXT runtime·mcall(SB), 7, $0
	MOVL	fn+0(FP), DI
	
	get_tls(CX)
	MOVL	g(CX), AX	// save state in g->sched
	MOVL	0(SP), BX	// caller's PC
	MOVL	BX, (g_sched+gobuf_pc)(AX)
	LEAL	4(SP), BX	// caller's SP
	MOVL	BX, (g_sched+gobuf_sp)(AX)
	MOVL	AX, (g_sched+gobuf_g)(AX)

	// switch to m->g0 & its stack, call fn
	MOVL	m(CX), BX
	MOVL	m_g0(BX), SI
	CMPL	SI, AX	// if g == m->g0 call badmcall
	JNE	2(PC)
	CALL	runtime·badmcall(SB)
	MOVL	SI, g(CX)	// g = m->g0
	MOVL	(g_sched+gobuf_sp)(SI), SP	// sp = m->g0->sched.sp
	PUSHL	AX
	CALL	DI
	POPL	AX
	CALL	runtime·badmcall2(SB)
	RET

/*
 * support for morestack
 */

// Called during function prolog when more stack is needed.
TEXT runtime·morestack(SB),7,$0
	// Cannot grow scheduler stack (m->g0).
	get_tls(CX)
	MOVL	m(CX), BX
	MOVL	m_g0(BX), SI
	CMPL	g(CX), SI
	JNE	2(PC)
	INT	$3
	
	MOVL	DX, m_cret(BX)

	// frame size in DI
	// arg size in AX
	// Save in m.
	MOVL	DI, m_moreframesize(BX)
	MOVL	AX, m_moreargsize(BX)

	// Called from f.
	// Set m->morebuf to f's caller.
	MOVL	4(SP), DI	// f's caller's PC
	MOVL	DI, (m_morebuf+gobuf_pc)(BX)
	LEAL	8(SP), CX	// f's caller's SP
	MOVL	CX, (m_morebuf+gobuf_sp)(BX)
	MOVL	CX, m_moreargp(BX)
	get_tls(CX)
	MOVL	g(CX), SI
	MOVL	SI, (m_morebuf+gobuf_g)(BX)

	// Set m->morepc to f's PC.
	MOVL	0(SP), AX
	MOVL	AX, m_morepc(BX)

	// Call newstack on m->g0's stack.
	MOVL	m_g0(BX), BP
	MOVL	BP, g(CX)
	MOVL	(g_sched+gobuf_sp)(BP), AX
	MOVL	-4(AX), BX	// fault if CALL would, before smashing SP
	MOVL	AX, SP
	CALL	runtime·newstack(SB)
	MOVL	$0, 0x1003	// crash if newstack returns
	RET

// Called from reflection library.  Mimics morestack,
// reuses stack growth code to create a frame
// with the desired args running the desired function.
//
// func call(fn *byte, arg *byte, argsize uint32).
TEXT reflect·call(SB), 7, $0
	get_tls(CX)
	MOVL	m(CX), BX

	// Save our caller's state as the PC and SP to
	// restore when returning from f.
	MOVL	0(SP), AX	// our caller's PC
	MOVL	AX, (m_morebuf+gobuf_pc)(BX)
	LEAL	4(SP), AX	// our caller's SP
	MOVL	AX, (m_morebuf+gobuf_sp)(BX)
	MOVL	g(CX), AX
	MOVL	AX, (m_morebuf+gobuf_g)(BX)

	// Set up morestack arguments to call f on a new stack.
	// We set f's frame size to 1, as a hint to newstack
	// that this is a call from reflect·call.
	// If it turns out that f needs a larger frame than
	// the default stack, f's usual stack growth prolog will
	// allocate a new segment (and recopy the arguments).
	MOVL	4(SP), AX	// fn
	MOVL	8(SP), DX	// arg frame
	MOVL	12(SP), CX	// arg size

	MOVL	AX, m_morepc(BX)	// f's PC
	MOVL	DX, m_moreargp(BX)	// f's argument pointer
	MOVL	CX, m_moreargsize(BX)	// f's argument size
	MOVL	$1, m_moreframesize(BX)	// f's frame size

	// Call newstack on m->g0's stack.
	MOVL	m_g0(BX), BP
	get_tls(CX)
	MOVL	BP, g(CX)
	MOVL	(g_sched+gobuf_sp)(BP), SP
	CALL	runtime·newstack(SB)
	MOVL	$0, 0x1103	// crash if newstack returns
	RET


// Return point when leaving stack.
TEXT runtime·lessstack(SB), 7, $0
	// Save return value in m->cret
	get_tls(CX)
	MOVL	m(CX), BX
	MOVL	AX, m_cret(BX)

	// Call oldstack on m->g0's stack.
	MOVL	m_g0(BX), BP
	MOVL	BP, g(CX)
	MOVL	(g_sched+gobuf_sp)(BP), SP
	CALL	runtime·oldstack(SB)
	MOVL	$0, 0x1004	// crash if oldstack returns
	RET


// bool cas(int32 *val, int32 old, int32 new)
// Atomically:
//	if(*val == old){
//		*val = new;
//		return 1;
//	}else
//		return 0;
TEXT runtime·cas(SB), 7, $0
	MOVL	4(SP), BX
	MOVL	8(SP), AX
	MOVL	12(SP), CX
	LOCK
	CMPXCHGL	CX, 0(BX)
	JZ 3(PC)
	MOVL	$0, AX
	RET
	MOVL	$1, AX
	RET

// bool runtime·cas64(uint64 *val, uint64 *old, uint64 new)
// Atomically:
//	if(*val == *old){
//		*val = new;
//		return 1;
//	} else {
//		*old = *val
//		return 0;
//	}
TEXT runtime·cas64(SB), 7, $0
	MOVL	4(SP), BP
	MOVL	8(SP), SI
	MOVL	0(SI), AX
	MOVL	4(SI), DX
	MOVL	12(SP), BX
	MOVL	16(SP), CX
	LOCK
	CMPXCHG8B	0(BP)
	JNZ	cas64_fail
	MOVL	$1, AX
	RET
cas64_fail:
	MOVL	AX, 0(SI)
	MOVL	DX, 4(SI)
	MOVL	$0, AX
	RET

// bool casp(void **p, void *old, void *new)
// Atomically:
//	if(*p == old){
//		*p = new;
//		return 1;
//	}else
//		return 0;
TEXT runtime·casp(SB), 7, $0
	MOVL	4(SP), BX
	MOVL	8(SP), AX
	MOVL	12(SP), CX
	LOCK
	CMPXCHGL	CX, 0(BX)
	JZ 3(PC)
	MOVL	$0, AX
	RET
	MOVL	$1, AX
	RET

// uint32 xadd(uint32 volatile *val, int32 delta)
// Atomically:
//	*val += delta;
//	return *val;
TEXT runtime·xadd(SB), 7, $0
	MOVL	4(SP), BX
	MOVL	8(SP), AX
	MOVL	AX, CX
	LOCK
	XADDL	AX, 0(BX)
	ADDL	CX, AX
	RET

TEXT runtime·xchg(SB), 7, $0
	MOVL	4(SP), BX
	MOVL	8(SP), AX
	XCHGL	AX, 0(BX)
	RET

TEXT runtime·procyield(SB),7,$0
	MOVL	4(SP), AX
again:
	PAUSE
	SUBL	$1, AX
	JNZ	again
	RET

TEXT runtime·atomicstorep(SB), 7, $0
	MOVL	4(SP), BX
	MOVL	8(SP), AX
	XCHGL	AX, 0(BX)
	RET

TEXT runtime·atomicstore(SB), 7, $0
	MOVL	4(SP), BX
	MOVL	8(SP), AX
	XCHGL	AX, 0(BX)
	RET

// uint64 atomicload64(uint64 volatile* addr);
// so actually
// void atomicload64(uint64 *res, uint64 volatile *addr);
TEXT runtime·atomicload64(SB), 7, $0
	MOVL    4(SP), BX
	MOVL	8(SP), AX
	// MOVQ (%EAX), %MM0
	BYTE $0x0f; BYTE $0x6f; BYTE $0x00
	// MOVQ %MM0, 0(%EBX)
	BYTE $0x0f; BYTE $0x7f; BYTE $0x03
	// EMMS
	BYTE $0x0F; BYTE $0x77
	RET

// void runtime·atomicstore64(uint64 volatile* addr, uint64 v);
TEXT runtime·atomicstore64(SB), 7, $0
	MOVL	4(SP), AX
	// MOVQ and EMMS were introduced on the Pentium MMX.
	// MOVQ 0x8(%ESP), %MM0
	BYTE $0x0f; BYTE $0x6f; BYTE $0x44; BYTE $0x24; BYTE $0x08
	// MOVQ %MM0, (%EAX)
	BYTE $0x0f; BYTE $0x7f; BYTE $0x00 
	// EMMS
	BYTE $0x0F; BYTE $0x77
	// This is essentially a no-op, but it provides required memory fencing.
	// It can be replaced with MFENCE, but MFENCE was introduced only on the Pentium4 (SSE2).
	MOVL	$0, AX
	LOCK
	XADDL	AX, (SP)
	RET

// void jmpdefer(fn, sp);
// called from deferreturn.
// 1. pop the caller
// 2. sub 5 bytes from the callers return
// 3. jmp to the argument
TEXT runtime·jmpdefer(SB), 7, $0
	MOVL	4(SP), DX	// fn
	MOVL	8(SP), BX	// caller sp
	LEAL	-4(BX), SP	// caller sp after CALL
	SUBL	$5, (SP)	// return to CALL again
	MOVL	0(DX), BX
	JMP	BX	// but first run the deferred function

// Save state of caller into g->sched.
TEXT gosave<>(SB),7,$0
	PUSHL	AX
	PUSHL	BX
	get_tls(BX)
	MOVL	g(BX), BX
	LEAL	arg+0(FP), AX
	MOVL	AX, (g_sched+gobuf_sp)(BX)
	MOVL	-4(AX), AX
	MOVL	AX, (g_sched+gobuf_pc)(BX)
	MOVL	$0, (g_sched+gobuf_ret)(BX)
	MOVL	$0, (g_sched+gobuf_ctxt)(BX)
	POPL	BX
	POPL	AX
	RET

// asmcgocall(void(*fn)(void*), void *arg)
// Call fn(arg) on the scheduler stack,
// aligned appropriately for the gcc ABI.
// See cgocall.c for more details.
TEXT runtime·asmcgocall(SB),7,$0
	MOVL	fn+0(FP), AX
	MOVL	arg+4(FP), BX
	MOVL	SP, DX

	// Figure out if we need to switch to m->g0 stack.
	// We get called to create new OS threads too, and those
	// come in on the m->g0 stack already.
	get_tls(CX)
	MOVL	m(CX), BP
	MOVL	m_g0(BP), SI
	MOVL	g(CX), DI
	CMPL	SI, DI
	JEQ	4(PC)
	CALL	gosave<>(SB)
	MOVL	SI, g(CX)
	MOVL	(g_sched+gobuf_sp)(SI), SP

	// Now on a scheduling stack (a pthread-created stack).
	SUBL	$32, SP
	ANDL	$~15, SP	// alignment, perhaps unnecessary
	MOVL	DI, 8(SP)	// save g
	MOVL	DX, 4(SP)	// save SP
	MOVL	BX, 0(SP)	// first argument in x86-32 ABI
	CALL	AX

	// Restore registers, g, stack pointer.
	get_tls(CX)
	MOVL	8(SP), DI
	MOVL	DI, g(CX)
	MOVL	4(SP), SP
	RET

// cgocallback(void (*fn)(void*), void *frame, uintptr framesize)
// Turn the fn into a Go func (by taking its address) and call
// cgocallback_gofunc.
TEXT runtime·cgocallback(SB),7,$12
	LEAL	fn+0(FP), AX
	MOVL	AX, 0(SP)
	MOVL	frame+4(FP), AX
	MOVL	AX, 4(SP)
	MOVL	framesize+8(FP), AX
	MOVL	AX, 8(SP)
	MOVL	$runtime·cgocallback_gofunc(SB), AX
	CALL	AX
	RET

// cgocallback_gofunc(FuncVal*, void *frame, uintptr framesize)
// See cgocall.c for more details.
TEXT runtime·cgocallback_gofunc(SB),7,$12
	// If m is nil, Go did not create the current thread.
	// Call needm to obtain one for temporary use.
	// In this case, we're running on the thread stack, so there's
	// lots of space, but the linker doesn't know. Hide the call from
	// the linker analysis by using an indirect call through AX.
	get_tls(CX)
#ifdef GOOS_windows
	CMPL	CX, $0
	JNE	3(PC)
	PUSHL	$0
	JMP needm
#endif
	MOVL	m(CX), BP
	PUSHL	BP
	CMPL	BP, $0
	JNE	havem
needm:
	MOVL	$runtime·needm(SB), AX
	CALL	AX
	get_tls(CX)
	MOVL	m(CX), BP

havem:
	// Now there's a valid m, and we're running on its m->g0.
	// Save current m->g0->sched.sp on stack and then set it to SP.
	// Save current sp in m->g0->sched.sp in preparation for
	// switch back to m->curg stack.
	MOVL	m_g0(BP), SI
	PUSHL	(g_sched+gobuf_sp)(SI)
	MOVL	SP, (g_sched+gobuf_sp)(SI)

	// Switch to m->curg stack and call runtime.cgocallbackg
	// with the three arguments.  Because we are taking over
	// the execution of m->curg but *not* resuming what had
	// been running, we need to save that information (m->curg->sched)
	// so that we can restore it when we're done. 
	// We can restore m->curg->sched.sp easily, because calling
	// runtime.cgocallbackg leaves SP unchanged upon return.
	// To save m->curg->sched.pc, we push it onto the stack.
	// This has the added benefit that it looks to the traceback
	// routine like cgocallbackg is going to return to that
	// PC (because we defined cgocallbackg to have
	// a frame size of 12, the same amount that we use below),
	// so that the traceback will seamlessly trace back into
	// the earlier calls.
	MOVL	fn+0(FP), AX
	MOVL	frame+4(FP), BX
	MOVL	framesize+8(FP), DX

	MOVL	m_curg(BP), SI
	MOVL	SI, g(CX)
	MOVL	(g_sched+gobuf_sp)(SI), DI  // prepare stack as DI

	// Push gobuf.pc
	MOVL	(g_sched+gobuf_pc)(SI), BP
	SUBL	$4, DI
	MOVL	BP, 0(DI)

	// Push arguments to cgocallbackg.
	// Frame size here must match the frame size above
	// to trick traceback routines into doing the right thing.
	SUBL	$12, DI
	MOVL	AX, 0(DI)
	MOVL	BX, 4(DI)
	MOVL	DX, 8(DI)
	
	// Switch stack and make the call.
	MOVL	DI, SP
	CALL	runtime·cgocallbackg(SB)

	// Restore g->sched (== m->curg->sched) from saved values.
	get_tls(CX)
	MOVL	g(CX), SI
	MOVL	12(SP), BP
	MOVL	BP, (g_sched+gobuf_pc)(SI)
	LEAL	(12+4)(SP), DI
	MOVL	DI, (g_sched+gobuf_sp)(SI)

	// Switch back to m->g0's stack and restore m->g0->sched.sp.
	// (Unlike m->curg, the g0 goroutine never uses sched.pc,
	// so we do not have to restore it.)
	MOVL	m(CX), BP
	MOVL	m_g0(BP), SI
	MOVL	SI, g(CX)
	MOVL	(g_sched+gobuf_sp)(SI), SP
	POPL	(g_sched+gobuf_sp)(SI)
	
	// If the m on entry was nil, we called needm above to borrow an m
	// for the duration of the call. Since the call is over, return it with dropm.
	POPL	BP
	CMPL	BP, $0
	JNE 3(PC)
	MOVL	$runtime·dropm(SB), AX
	CALL	AX

	// Done!
	RET

// void setmg(M*, G*); set m and g. for use by needm.
TEXT runtime·setmg(SB), 7, $0
#ifdef GOOS_windows
	MOVL	mm+0(FP), AX
	CMPL	AX, $0
	JNE	settls
	MOVL	$0, 0x14(FS)
	RET
settls:
	LEAL	m_tls(AX), AX
	MOVL	AX, 0x14(FS)
#endif
	MOVL	mm+0(FP), AX
	get_tls(CX)
	MOVL	mm+0(FP), AX
	MOVL	AX, m(CX)
	MOVL	gg+4(FP), BX
	MOVL	BX, g(CX)
	RET

// void setmg_gcc(M*, G*); set m and g. for use by gcc
TEXT setmg_gcc<>(SB), 7, $0	
	get_tls(AX)
	MOVL	mm+0(FP), DX
	MOVL	DX, m(AX)
	MOVL	gg+4(FP), DX
	MOVL	DX,g (AX)
	RET

// check that SP is in range [g->stackbase, g->stackguard)
TEXT runtime·stackcheck(SB), 7, $0
	get_tls(CX)
	MOVL	g(CX), AX
	CMPL	g_stackbase(AX), SP
	JHI	2(PC)
	INT	$3
	CMPL	SP, g_stackguard(AX)
	JHI	2(PC)
	INT	$3
	RET

TEXT runtime·memclr(SB),7,$0
	MOVL	4(SP), DI		// arg 1 addr
	MOVL	8(SP), CX		// arg 2 count
	MOVL	CX, BX
	ANDL	$3, BX
	SHRL	$2, CX
	MOVL	$0, AX
	CLD
	REP
	STOSL
	MOVL	BX, CX
	REP
	STOSB
	RET

TEXT runtime·getcallerpc(SB),7,$0
	MOVL	x+0(FP),AX		// addr of first arg
	MOVL	-4(AX),AX		// get calling pc
	RET

TEXT runtime·setcallerpc(SB),7,$0
	MOVL	x+0(FP),AX		// addr of first arg
	MOVL	x+4(FP), BX
	MOVL	BX, -4(AX)		// set calling pc
	RET

TEXT runtime·getcallersp(SB), 7, $0
	MOVL	sp+0(FP), AX
	RET

// int64 runtime·cputicks(void), so really
// void runtime·cputicks(int64 *ticks)
TEXT runtime·cputicks(SB),7,$0
	RDTSC
	MOVL	ret+0(FP), DI
	MOVL	AX, 0(DI)
	MOVL	DX, 4(DI)
	RET

TEXT runtime·ldt0setup(SB),7,$16
	// set up ldt 7 to point at tls0
	// ldt 1 would be fine on Linux, but on OS X, 7 is as low as we can go.
	// the entry number is just a hint.  setldt will set up GS with what it used.
	MOVL	$7, 0(SP)
	LEAL	runtime·tls0(SB), AX
	MOVL	AX, 4(SP)
	MOVL	$32, 8(SP)	// sizeof(tls array)
	CALL	runtime·setldt(SB)
	RET

TEXT runtime·emptyfunc(SB),0,$0
	RET

TEXT runtime·abort(SB),7,$0
	INT $0x3

TEXT runtime·stackguard(SB),7,$0
	MOVL	SP, DX
	MOVL	DX, sp+0(FP)
	get_tls(CX)
	MOVL	g(CX), BX
	MOVL	g_stackguard(BX), DX
	MOVL	DX, limit+4(FP)
	RET

GLOBL runtime·tls0(SB), $32

// hash function using AES hardware instructions
TEXT runtime·aeshash(SB),7,$0
	MOVL	4(SP), DX	// ptr to hash value
	MOVL	8(SP), CX	// size
	MOVL	12(SP), AX	// ptr to data
	JMP	runtime·aeshashbody(SB)

TEXT runtime·aeshashstr(SB),7,$0
	MOVL	4(SP), DX	// ptr to hash value
	MOVL	12(SP), AX	// ptr to string struct
	MOVL	4(AX), CX	// length of string
	MOVL	(AX), AX	// string data
	JMP	runtime·aeshashbody(SB)

// AX: data
// CX: length
// DX: ptr to seed input / hash output
TEXT runtime·aeshashbody(SB),7,$0
	MOVL	(DX), X0	// seed to low 32 bits of xmm0
	PINSRD	$1, CX, X0	// size to next 32 bits of xmm0
	MOVO	runtime·aeskeysched+0(SB), X2
	MOVO	runtime·aeskeysched+16(SB), X3
	CMPL	CX, $16
	JB	aessmall
aesloop:
	CMPL	CX, $16
	JBE	aesloopend
	MOVOU	(AX), X1
	AESENC	X2, X0
	AESENC	X1, X0
	SUBL	$16, CX
	ADDL	$16, AX
	JMP	aesloop
// 1-16 bytes remaining
aesloopend:
	// This load may overlap with the previous load above.
	// We'll hash some bytes twice, but that's ok.
	MOVOU	-16(AX)(CX*1), X1
	JMP	partial
// 0-15 bytes
aessmall:
	TESTL	CX, CX
	JE	finalize	// 0 bytes

	CMPB	AX, $0xf0
	JA	highpartial

	// 16 bytes loaded at this address won't cross
	// a page boundary, so we can load it directly.
	MOVOU	(AX), X1
	ADDL	CX, CX
	PAND	masks(SB)(CX*8), X1
	JMP	partial
highpartial:
	// address ends in 1111xxxx.  Might be up against
	// a page boundary, so load ending at last byte.
	// Then shift bytes down using pshufb.
	MOVOU	-16(AX)(CX*1), X1
	ADDL	CX, CX
	PSHUFB	shifts(SB)(CX*8), X1
partial:
	// incorporate partial block into hash
	AESENC	X3, X0
	AESENC	X1, X0
finalize:	
	// finalize hash
	AESENC	X2, X0
	AESENC	X3, X0
	AESENC	X2, X0
	MOVL	X0, (DX)
	RET

TEXT runtime·aeshash32(SB),7,$0
	MOVL	4(SP), DX	// ptr to hash value
	MOVL	12(SP), AX	// ptr to data
	MOVL	(DX), X0	// seed
	PINSRD	$1, (AX), X0	// data
	AESENC	runtime·aeskeysched+0(SB), X0
	AESENC	runtime·aeskeysched+16(SB), X0
	AESENC	runtime·aeskeysched+0(SB), X0
	MOVL	X0, (DX)
	RET

TEXT runtime·aeshash64(SB),7,$0
	MOVL	4(SP), DX	// ptr to hash value
	MOVL	12(SP), AX	// ptr to data
	MOVQ	(AX), X0	// data
	PINSRD	$2, (DX), X0	// seed
	AESENC	runtime·aeskeysched+0(SB), X0
	AESENC	runtime·aeskeysched+16(SB), X0
	AESENC	runtime·aeskeysched+0(SB), X0
	MOVL	X0, (DX)
	RET


// simple mask to get rid of data in the high part of the register.
TEXT masks(SB),7,$0
	LONG $0x00000000
	LONG $0x00000000
	LONG $0x00000000
	LONG $0x00000000
	
	LONG $0x000000ff
	LONG $0x00000000
	LONG $0x00000000
	LONG $0x00000000
	
	LONG $0x0000ffff
	LONG $0x00000000
	LONG $0x00000000
	LONG $0x00000000
	
	LONG $0x00ffffff
	LONG $0x00000000
	LONG $0x00000000
	LONG $0x00000000
	
	LONG $0xffffffff
	LONG $0x00000000
	LONG $0x00000000
	LONG $0x00000000
	
	LONG $0xffffffff
	LONG $0x000000ff
	LONG $0x00000000
	LONG $0x00000000
	
	LONG $0xffffffff
	LONG $0x0000ffff
	LONG $0x00000000
	LONG $0x00000000
	
	LONG $0xffffffff
	LONG $0x00ffffff
	LONG $0x00000000
	LONG $0x00000000
	
	LONG $0xffffffff
	LONG $0xffffffff
	LONG $0x00000000
	LONG $0x00000000
	
	LONG $0xffffffff
	LONG $0xffffffff
	LONG $0x000000ff
	LONG $0x00000000
	
	LONG $0xffffffff
	LONG $0xffffffff
	LONG $0x0000ffff
	LONG $0x00000000
	
	LONG $0xffffffff
	LONG $0xffffffff
	LONG $0x00ffffff
	LONG $0x00000000
	
	LONG $0xffffffff
	LONG $0xffffffff
	LONG $0xffffffff
	LONG $0x00000000
	
	LONG $0xffffffff
	LONG $0xffffffff
	LONG $0xffffffff
	LONG $0x000000ff
	
	LONG $0xffffffff
	LONG $0xffffffff
	LONG $0xffffffff
	LONG $0x0000ffff
	
	LONG $0xffffffff
	LONG $0xffffffff
	LONG $0xffffffff
	LONG $0x00ffffff

	// these are arguments to pshufb.  They move data down from
	// the high bytes of the register to the low bytes of the register.
	// index is how many bytes to move.
TEXT shifts(SB),7,$0
	LONG $0x00000000
	LONG $0x00000000
	LONG $0x00000000
	LONG $0x00000000
	
	LONG $0xffffff0f
	LONG $0xffffffff
	LONG $0xffffffff
	LONG $0xffffffff
	
	LONG $0xffff0f0e
	LONG $0xffffffff
	LONG $0xffffffff
	LONG $0xffffffff
	
	LONG $0xff0f0e0d
	LONG $0xffffffff
	LONG $0xffffffff
	LONG $0xffffffff
	
	LONG $0x0f0e0d0c
	LONG $0xffffffff
	LONG $0xffffffff
	LONG $0xffffffff
	
	LONG $0x0e0d0c0b
	LONG $0xffffff0f
	LONG $0xffffffff
	LONG $0xffffffff
	
	LONG $0x0d0c0b0a
	LONG $0xffff0f0e
	LONG $0xffffffff
	LONG $0xffffffff
	
	LONG $0x0c0b0a09
	LONG $0xff0f0e0d
	LONG $0xffffffff
	LONG $0xffffffff
	
	LONG $0x0b0a0908
	LONG $0x0f0e0d0c
	LONG $0xffffffff
	LONG $0xffffffff
	
	LONG $0x0a090807
	LONG $0x0e0d0c0b
	LONG $0xffffff0f
	LONG $0xffffffff
	
	LONG $0x09080706
	LONG $0x0d0c0b0a
	LONG $0xffff0f0e
	LONG $0xffffffff
	
	LONG $0x08070605
	LONG $0x0c0b0a09
	LONG $0xff0f0e0d
	LONG $0xffffffff
	
	LONG $0x07060504
	LONG $0x0b0a0908
	LONG $0x0f0e0d0c
	LONG $0xffffffff
	
	LONG $0x06050403
	LONG $0x0a090807
	LONG $0x0e0d0c0b
	LONG $0xffffff0f
	
	LONG $0x05040302
	LONG $0x09080706
	LONG $0x0d0c0b0a
	LONG $0xffff0f0e
	
	LONG $0x04030201
	LONG $0x08070605
	LONG $0x0c0b0a09
	LONG $0xff0f0e0d

TEXT runtime·memeq(SB),7,$0
	MOVL	a+0(FP), SI
	MOVL	b+4(FP), DI
	MOVL	count+8(FP), BX
	JMP	runtime·memeqbody(SB)


TEXT bytes·Equal(SB),7,$0
	MOVL	a_len+4(FP), BX
	MOVL	b_len+16(FP), CX
	XORL	AX, AX
	CMPL	BX, CX
	JNE	eqret
	MOVL	a+0(FP), SI
	MOVL	b+12(FP), DI
	CALL	runtime·memeqbody(SB)
eqret:
	MOVB	AX, ret+24(FP)
	RET

// a in SI
// b in DI
// count in BX
TEXT runtime·memeqbody(SB),7,$0
	XORL	AX, AX

	CMPL	BX, $4
	JB	small

	// 64 bytes at a time using xmm registers
hugeloop:
	CMPL	BX, $64
	JB	bigloop
	TESTL	$0x4000000, runtime·cpuid_edx(SB) // check for sse2
	JE	bigloop
	MOVOU	(SI), X0
	MOVOU	(DI), X1
	MOVOU	16(SI), X2
	MOVOU	16(DI), X3
	MOVOU	32(SI), X4
	MOVOU	32(DI), X5
	MOVOU	48(SI), X6
	MOVOU	48(DI), X7
	PCMPEQB	X1, X0
	PCMPEQB	X3, X2
	PCMPEQB	X5, X4
	PCMPEQB	X7, X6
	PAND	X2, X0
	PAND	X6, X4
	PAND	X4, X0
	PMOVMSKB X0, DX
	ADDL	$64, SI
	ADDL	$64, DI
	SUBL	$64, BX
	CMPL	DX, $0xffff
	JEQ	hugeloop
	RET

	// 4 bytes at a time using 32-bit register
bigloop:
	CMPL	BX, $4
	JBE	leftover
	MOVL	(SI), CX
	MOVL	(DI), DX
	ADDL	$4, SI
	ADDL	$4, DI
	SUBL	$4, BX
	CMPL	CX, DX
	JEQ	bigloop
	RET

	// remaining 0-4 bytes
leftover:
	MOVL	-4(SI)(BX*1), CX
	MOVL	-4(DI)(BX*1), DX
	CMPL	CX, DX
	SETEQ	AX
	RET

small:
	CMPL	BX, $0
	JEQ	equal

	LEAL	0(BX*8), CX
	NEGL	CX

	MOVL	SI, DX
	CMPB	DX, $0xfc
	JA	si_high

	// load at SI won't cross a page boundary.
	MOVL	(SI), SI
	JMP	si_finish
si_high:
	// address ends in 111111xx.  Load up to bytes we want, move to correct position.
	MOVL	-4(SI)(BX*1), SI
	SHRL	CX, SI
si_finish:

	// same for DI.
	MOVL	DI, DX
	CMPB	DX, $0xfc
	JA	di_high
	MOVL	(DI), DI
	JMP	di_finish
di_high:
	MOVL	-4(DI)(BX*1), DI
	SHRL	CX, DI
di_finish:

	SUBL	SI, DI
	SHLL	CX, DI
equal:
	SETEQ	AX
	RET

TEXT runtime·cmpstring(SB),7,$0
	MOVL	s1+0(FP), SI
	MOVL	s1+4(FP), BX
	MOVL	s2+8(FP), DI
	MOVL	s2+12(FP), DX
	CALL	runtime·cmpbody(SB)
	MOVL	AX, res+16(FP)
	RET

TEXT bytes·Compare(SB),7,$0
	MOVL	s1+0(FP), SI
	MOVL	s1+4(FP), BX
	MOVL	s2+12(FP), DI
	MOVL	s2+16(FP), DX
	CALL	runtime·cmpbody(SB)
	MOVL	AX, res+24(FP)
	RET

// input:
//   SI = a
//   DI = b
//   BX = alen
//   DX = blen
// output:
//   AX = 1/0/-1
TEXT runtime·cmpbody(SB),7,$0
	CMPL	SI, DI
	JEQ	cmp_allsame
	CMPL	BX, DX
	MOVL	DX, BP
	CMOVLLT	BX, BP // BP = min(alen, blen)
	CMPL	BP, $4
	JB	cmp_small
	TESTL	$0x4000000, runtime·cpuid_edx(SB) // check for sse2
	JE	cmp_mediumloop
cmp_largeloop:
	CMPL	BP, $16
	JB	cmp_mediumloop
	MOVOU	(SI), X0
	MOVOU	(DI), X1
	PCMPEQB X0, X1
	PMOVMSKB X1, AX
	XORL	$0xffff, AX	// convert EQ to NE
	JNE	cmp_diff16	// branch if at least one byte is not equal
	ADDL	$16, SI
	ADDL	$16, DI
	SUBL	$16, BP
	JMP	cmp_largeloop

cmp_diff16:
	BSFL	AX, BX	// index of first byte that differs
	XORL	AX, AX
	MOVB	(SI)(BX*1), CX
	CMPB	CX, (DI)(BX*1)
	SETHI	AX
	LEAL	-1(AX*2), AX	// convert 1/0 to +1/-1
	RET

cmp_mediumloop:
	CMPL	BP, $4
	JBE	cmp_0through4
	MOVL	(SI), AX
	MOVL	(DI), CX
	CMPL	AX, CX
	JNE	cmp_diff4
	ADDL	$4, SI
	ADDL	$4, DI
	SUBL	$4, BP
	JMP	cmp_mediumloop

cmp_0through4:
	MOVL	-4(SI)(BP*1), AX
	MOVL	-4(DI)(BP*1), CX
	CMPL	AX, CX
	JEQ	cmp_allsame

cmp_diff4:
	BSWAPL	AX	// reverse order of bytes
	BSWAPL	CX
	XORL	AX, CX	// find bit differences
	BSRL	CX, CX	// index of highest bit difference
	SHRL	CX, AX	// move a's bit to bottom
	ANDL	$1, AX	// mask bit
	LEAL	-1(AX*2), AX // 1/0 => +1/-1
	RET

	// 0-3 bytes in common
cmp_small:
	LEAL	(BP*8), CX
	NEGL	CX
	JEQ	cmp_allsame

	// load si
	CMPB	SI, $0xfc
	JA	cmp_si_high
	MOVL	(SI), SI
	JMP	cmp_si_finish
cmp_si_high:
	MOVL	-4(SI)(BP*1), SI
	SHRL	CX, SI
cmp_si_finish:
	SHLL	CX, SI

	// same for di
	CMPB	DI, $0xfc
	JA	cmp_di_high
	MOVL	(DI), DI
	JMP	cmp_di_finish
cmp_di_high:
	MOVL	-4(DI)(BP*1), DI
	SHRL	CX, DI
cmp_di_finish:
	SHLL	CX, DI

	BSWAPL	SI	// reverse order of bytes
	BSWAPL	DI
	XORL	SI, DI	// find bit differences
	JEQ	cmp_allsame
	BSRL	DI, CX	// index of highest bit difference
	SHRL	CX, SI	// move a's bit to bottom
	ANDL	$1, SI	// mask bit
	LEAL	-1(SI*2), AX // 1/0 => +1/-1
	RET

	// all the bytes in common are the same, so we just need
	// to compare the lengths.
cmp_allsame:
	XORL	AX, AX
	XORL	CX, CX
	CMPL	BX, DX
	SETGT	AX	// 1 if alen > blen
	SETEQ	CX	// 1 if alen == blen
	LEAL	-1(CX)(AX*2), AX	// 1,0,-1 result
	RET
