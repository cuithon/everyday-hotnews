// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "zasm_GOOS_GOARCH.h"
#include "../../cmd/ld/textflag.h"
#include "syscall_nacl.h"

#define NACL_SYSCALL(code) \
	MOVW	$(0x10000 + ((code)<<5)), R8; BL (R8)

#define NACL_SYSJMP(code) \
	MOVW	$(0x10000 + ((code)<<5)), R8; B (R8)

TEXT runtime·exit(SB),NOSPLIT,$0
	MOVW	arg1+0(FP), R0
	NACL_SYSJMP(SYS_exit)

TEXT runtime·exit1(SB),NOSPLIT,$0
	MOVW	arg1+0(FP), R0
	NACL_SYSJMP(SYS_thread_exit)

TEXT runtime·open(SB),NOSPLIT,$0
	MOVW	arg1+0(FP), R0
	MOVW	arg2+0(FP), R1
	MOVW	arg3+0(FP), R2
	NACL_SYSJMP(SYS_open)

TEXT runtime·close(SB),NOSPLIT,$0
	MOVW	arg1+0(FP), R0
	NACL_SYSJMP(SYS_close)

TEXT runtime·read(SB),NOSPLIT,$0
	MOVW	arg1+0(FP), R0
	MOVW	arg2+4(FP), R1
	MOVW	arg3+8(FP), R2
	NACL_SYSJMP(SYS_read)

// func naclWrite(fd int, b []byte) int
TEXT syscall·naclWrite(SB),NOSPLIT,$0
	MOVW	arg1+0(FP), R0
	MOVW	arg2+4(FP), R1
	MOVW	arg3+8(FP), R2
	NACL_SYSCALL(SYS_write)
	MOVW	R0, ret+16(FP)
	RET

TEXT runtime·write(SB),NOSPLIT,$0
	MOVW	arg1+0(FP), R0
	MOVW	arg2+4(FP), R1
	MOVW	arg3+8(FP), R2
	NACL_SYSJMP(SYS_write)

TEXT runtime·nacl_exception_stack(SB),NOSPLIT,$0
	MOVW	arg1+0(FP), R0
	MOVW	arg2+4(FP), R1
	NACL_SYSJMP(SYS_exception_stack)

TEXT runtime·nacl_exception_handler(SB),NOSPLIT,$0
	MOVW	arg1+0(FP), R0
	MOVW	arg2+4(FP), R1
	NACL_SYSJMP(SYS_exception_handler)

TEXT runtime·nacl_sem_create(SB),NOSPLIT,$0
	MOVW	arg1+0(FP), R0
	NACL_SYSJMP(SYS_sem_create)

TEXT runtime·nacl_sem_wait(SB),NOSPLIT,$0
	MOVW	arg1+0(FP), R0
	NACL_SYSJMP(SYS_sem_wait)

TEXT runtime·nacl_sem_post(SB),NOSPLIT,$0
	MOVW	arg1+0(FP), R0
	NACL_SYSJMP(SYS_sem_post)

TEXT runtime·nacl_mutex_create(SB),NOSPLIT,$0
	MOVW	arg1+0(FP), R0
	NACL_SYSJMP(SYS_mutex_create)

TEXT runtime·nacl_mutex_lock(SB),NOSPLIT,$0
	MOVW	arg1+0(FP), R0
	NACL_SYSJMP(SYS_mutex_lock)

TEXT runtime·nacl_mutex_trylock(SB),NOSPLIT,$0
	MOVW	arg1+0(FP), R0
	NACL_SYSJMP(SYS_mutex_trylock)

TEXT runtime·nacl_mutex_unlock(SB),NOSPLIT,$0
	MOVW	arg1+0(FP), R0
	NACL_SYSJMP(SYS_mutex_unlock)

TEXT runtime·nacl_cond_create(SB),NOSPLIT,$0
	MOVW	arg1+0(FP), R0
	NACL_SYSJMP(SYS_cond_create)

TEXT runtime·nacl_cond_wait(SB),NOSPLIT,$0
	MOVW	arg1+0(FP), R0
	MOVW	arg2+4(FP), R1
	NACL_SYSJMP(SYS_cond_wait)

TEXT runtime·nacl_cond_signal(SB),NOSPLIT,$0
	MOVW	arg1+0(FP), R0
	NACL_SYSJMP(SYS_cond_signal)

TEXT runtime·nacl_cond_broadcast(SB),NOSPLIT,$0
	MOVW	arg1+0(FP), R0
	NACL_SYSJMP(SYS_cond_broadcast)

TEXT runtime·nacl_cond_timed_wait_abs(SB),NOSPLIT,$0
	MOVW	arg1+0(FP), R0
	MOVW	arg2+4(FP), R1
	MOVW	arg3+8(FP), R2
	NACL_SYSJMP(SYS_cond_timed_wait_abs)

TEXT runtime·nacl_thread_create(SB),NOSPLIT,$0
	MOVW	arg1+0(FP), R0
	MOVW	arg2+4(FP), R1
	MOVW	arg3+8(FP), R2
	MOVW	arg4+12(FP), R3
	NACL_SYSJMP(SYS_thread_create)

TEXT runtime·mstart_nacl(SB),NOSPLIT,$0
	MOVW	0(R9), R0 // TLS
	MOVW	-8(R0), R1 // g
	MOVW	-4(R0), R2 // m
	MOVW	R2, g_m(R1)
	MOVW	R1, g
	B runtime·mstart(SB)

TEXT runtime·nacl_nanosleep(SB),NOSPLIT,$0
	MOVW	arg1+0(FP), R0
	MOVW	arg2+4(FP), R1
	NACL_SYSJMP(SYS_nanosleep)

TEXT runtime·osyield(SB),NOSPLIT,$0
	NACL_SYSJMP(SYS_sched_yield)

TEXT runtime·mmap(SB),NOSPLIT,$8
	MOVW	arg1+0(FP), R0
	MOVW	arg2+4(FP), R1
	MOVW	arg3+8(FP), R2
	MOVW	arg4+12(FP), R3
	MOVW	arg5+16(FP), R4
	// arg6:offset should be passed as a pointer (to int64)
	MOVW	arg6+20(FP), R5
	MOVW	R5, 4(R13)
	MOVW	$0, R6
	MOVW	R6, 8(R13)
	MOVW	$4(R13), R5
	MOVM.DB.W [R4,R5], (R13) // arg5 and arg6 are passed on stack
	NACL_SYSCALL(SYS_mmap)
	MOVM.IA.W (R13), [R4, R5]
	CMP	$-4095, R0
	RSB.HI	$0, R0
	RET

TEXT time·now(SB),NOSPLIT,$16
	MOVW	$0, R0 // real time clock
	MOVW	$4(R13), R1
	NACL_SYSCALL(SYS_clock_gettime)
	MOVW	4(R13), R0 // low 32-bit sec
	MOVW	8(R13), R1 // high 32-bit sec
	MOVW	12(R13), R2 // nsec
	MOVW	R0, sec+0(FP)
	MOVW	R1, sec+4(FP)
	MOVW	R2, sec+8(FP)
	RET

TEXT syscall·now(SB),NOSPLIT,$0
	B time·now(SB)

TEXT runtime·nacl_clock_gettime(SB),NOSPLIT,$0
	MOVW	arg1+0(FP), R0
	MOVW	arg2+4(FP), R1
	NACL_SYSJMP(SYS_clock_gettime)

// int64 nanotime(void) so really
// void nanotime(int64 *nsec)
TEXT runtime·nanotime(SB),NOSPLIT,$16
	MOVW	$0, R0 // real time clock
	MOVW	$4(R13), R1
	NACL_SYSCALL(SYS_clock_gettime)
	MOVW	4(R13), R0 // low 32-bit sec
	MOVW	8(R13), R1 // high 32-bit sec (ignored for now)
	MOVW	12(R13), R2 // nsec
	MOVW	$1000000000, R3
	MULLU	R0, R3, (R1, R0)
	MOVW	$0, R4
	ADD.S	R2, R0
	ADC	R4, R1
	MOVW	0(FP), R2
	MOVW	R0, 0(R2)
	MOVW	R1, 4(R2)
	RET

TEXT runtime·sigtramp(SB),NOSPLIT,$80
	// load g from thread context
	MOVW	$ctxt+-4(FP), R0
	MOVW	(16*4+10*4)(R0), g

	// check that g exists
	CMP	$0, g
	BNE 	4(PC)
	MOVW  	$runtime·badsignal2(SB), R11
	BL	(R11)
	RET

	// save g
	MOVW	g, R3
	MOVW	g, 20(R13)

	// g = m->gsignal
	MOVW	g_m(g), R8
	MOVW	m_gsignal(R8), g

	// copy arguments for call to sighandler
	MOVW	$11, R0
	MOVW	R0, 4(R13) // signal
	MOVW	$0, R0
	MOVW	R0, 8(R13) // siginfo
	MOVW	$ctxt+-4(FP), R0
	MOVW	R0, 12(R13) // context
	MOVW	R3, 16(R13) // g

	BL	runtime·sighandler(SB)

	// restore g
	MOVW	20(R13), g

sigtramp_ret:
	// Enable exceptions again.
	NACL_SYSCALL(SYS_exception_clear_flag)

	// Restore registers as best we can. Impossible to do perfectly.
	// See comment in sys_nacl_386.s for extended rationale.
	MOVW	$ctxt+-4(FP), R1
	ADD	$64, R1
	MOVW	(0*4)(R1), R0
	MOVW	(2*4)(R1), R2
	MOVW	(3*4)(R1), R3
	MOVW	(4*4)(R1), R4
	MOVW	(5*4)(R1), R5
	MOVW	(6*4)(R1), R6
	MOVW	(7*4)(R1), R7
	MOVW	(8*4)(R1), R8
	// cannot write to R9
	MOVW	(10*4)(R1), g
	MOVW	(11*4)(R1), R11
	MOVW	(12*4)(R1), R12
	MOVW	(13*4)(R1), R13
	MOVW	(14*4)(R1), R14
	MOVW	(15*4)(R1), R1
	B	(R1)

nog:
	MOVW	$0, R0
	RET

TEXT runtime·nacl_sysinfo(SB),NOSPLIT,$16
	RET

TEXT runtime·casp(SB),NOSPLIT,$0
	B	runtime·cas(SB)

// This is only valid for ARMv6+, however, NaCl/ARM is only defined
// for ARMv7A anyway.
// bool armcas(int32 *val, int32 old, int32 new)
// AtomiBLy:
//	if(*val == old){
//		*val = new;
//		return 1;
//	}else
//		return 0;
TEXT runtime·cas(SB),NOSPLIT,$0
	B runtime·armcas(SB)

TEXT runtime·read_tls_fallback(SB),NOSPLIT,$-4
	WORD $0xe7fedef0 // NACL_INSTR_ARM_ABORT_NOW (UDF #0xEDE0)
