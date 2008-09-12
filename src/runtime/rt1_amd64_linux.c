// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "runtime.h"
#include "amd64_linux.h"
#include "signals.h"

extern void _rt0_amd64_linux();
byte* startsym = (byte*)_rt0_amd64_linux;

/* From /usr/include/asm-x86_64/sigcontext.h */
struct _fpstate {
  uint16   cwd;
  uint16   swd;
  uint16   twd;    /* Note this is not the same as the 32bit/x87/FSAVE twd */
  uint16   fop;
  uint64   rip;
  uint32   rdp;
  uint32   mxcsr;
  uint32   mxcsr_mask;
  uint32   st_space[32];   /* 8*16 bytes for each FP-reg */
  uint32   xmm_space[64];  /* 16*16 bytes for each XMM-reg  */
  uint32   reserved2[24];
};

struct sigcontext {
  uint64 r8;
  uint64 r9;
  uint64 r10;
  uint64 r11;
  uint64 r12;
  uint64 r13;
  uint64 r14;
  uint64 r15;
  uint64 rdi;
  uint64 rsi;
  uint64 rbp;
  uint64 rbx;
  uint64 rdx;
  uint64 rax;
  uint64 rcx;
  uint64 rsp;
  uint64 rip;
  uint64 eflags;           /* RFLAGS */
  uint16 cs;
  uint16 gs;
  uint16 fs;
  uint16 __pad0;
  uint64 err;
  uint64 trapno;
  uint64 oldmask;
  uint64 cr2;
  struct _fpstate *fpstate;       /* zero when no FPU context */
  uint64 reserved1[8];
};


/* From /usr/include/asm-x86_64/signal.h */
typedef struct sigaltstack {
	void /*__user*/ *ss_sp;
	int32 ss_flags;
	uint64 ss_size;
} stack_t;

typedef uint64 sigset_t;


/* From /usr/include/asm-x86_64/ucontext.h */
struct ucontext {
  uint64            uc_flags;
  struct ucontext  *uc_link;
  stack_t           uc_stack;
  struct sigcontext uc_mcontext;
  sigset_t          uc_sigmask;   /* mask last for extensibility */
};


void
print_sigcontext(struct sigcontext *sc)
{
	prints("\nrax     0x");  sys·printpointer((void*)sc->rax);
	prints("\nrbx     0x");  sys·printpointer((void*)sc->rbx);
	prints("\nrcx     0x");  sys·printpointer((void*)sc->rcx);
	prints("\nrdx     0x");  sys·printpointer((void*)sc->rdx);
	prints("\nrdi     0x");  sys·printpointer((void*)sc->rdi);
	prints("\nrsi     0x");  sys·printpointer((void*)sc->rsi);
	prints("\nrbp     0x");  sys·printpointer((void*)sc->rbp);
	prints("\nrsp     0x");  sys·printpointer((void*)sc->rsp);
	prints("\nr8      0x");  sys·printpointer((void*)sc->r8 );
	prints("\nr9      0x");  sys·printpointer((void*)sc->r9 );
	prints("\nr10     0x");  sys·printpointer((void*)sc->r10);
	prints("\nr11     0x");  sys·printpointer((void*)sc->r11);
	prints("\nr12     0x");  sys·printpointer((void*)sc->r12);
	prints("\nr13     0x");  sys·printpointer((void*)sc->r13);
	prints("\nr14     0x");  sys·printpointer((void*)sc->r14);
	prints("\nr15     0x");  sys·printpointer((void*)sc->r15);
	prints("\nrip     0x");  sys·printpointer((void*)sc->rip);
	prints("\nrflags  0x");  sys·printpointer((void*)sc->eflags);
	prints("\ncs      0x");  sys·printpointer((void*)sc->cs);
	prints("\nfs      0x");  sys·printpointer((void*)sc->fs);
	prints("\ngs      0x");  sys·printpointer((void*)sc->gs);
	prints("\n");
}


/*
 * This assembler routine takes the args from registers, puts them on the stack,
 * and calls sighandler().
 */
extern void sigtramp();

/*
 * Rudimentary reverse-engineered definition of signal interface.
 * You'd think it would be documented.
 */
/* From /usr/include/bits/siginfo.h */
typedef struct siginfo {
	int32	si_signo;		/* signal number */
	int32	si_errno;		/* errno association */
	int32	si_code;		/* signal code */
	int32	si_status;		/* exit value */
	void	*si_addr;		/* faulting address */
	/* more stuff here */
} siginfo;


/* From /usr/include/bits/sigaction.h */
/* (gri) Is this correct? See e.g. /usr/include/asm-x86_64/signal.h */
typedef struct sigaction {
 	union {
		void (*sa_handler)(int32);
		void (*sa_sigaction)(int32, siginfo *, void *);
	} u;				/* signal handler */
	uint8 sa_mask[128];		/* signal mask to apply. 128? are they MORONS? */
	int32 sa_flags;			/* see signal options below */
	void (*sa_restorer) (void);	/* unused here; needed to return from trap? */
} sigaction;


void
sighandler(int32 sig, siginfo* info, void** context)
{
	if(panicking)	// traceback already printed
		sys·exit(2);

        struct sigcontext *sc = &(((struct ucontext *)context)->uc_mcontext);

	if(!inlinetrap(sig, (byte *)sc->rip)) {
		if(sig < 0 || sig >= NSIG){
			prints("Signal ");
			sys·printint(sig);
		}else{
			prints(sigtab[sig].name);
		}
	}

        prints("\nFaulting address: 0x");  sys·printpointer(info->si_addr);
        prints("\npc: 0x");  sys·printpointer((void *)sc->rip);
        prints("\n\n");

	traceback((void *)sc->rip, (void *)sc->rsp, (void *)sc->r15);
	tracebackothers((void*)sc->r15);
	print_sigcontext(sc);

	sys·breakpoint();
	sys·exit(2);
}


static sigaction a;

void
initsig(void)
{
	int32 i;
	a.u.sa_sigaction = (void*)sigtramp;
	a.sa_flags = 0x04;  /* SA_SIGINFO */
	for(i=0; i<sizeof(a.sa_mask); i++)
		a.sa_mask[i] = 0xFF;

	for(i = 0; i<NSIG; i++)
		if(sigtab[i].catch){
			sys·rt_sigaction(i, &a, (void*)0, 8);
		}
}

// Linux futex.
//
//	futexsleep(uint32 *addr, uint32 val)
//	futexwakeup(uint32 *addr)
//
// Futexsleep atomically checks if *addr == val and if so, sleeps on addr.
// Futexwakeup wakes up one thread sleeping on addr.
// Futexsleep is allowed to wake up spuriously.

enum
{
	FUTEX_WAIT = 0,
	FUTEX_WAKE = 1,

	EINTR = 4,
	EAGAIN = 11,
};

// TODO(rsc) I tried using 1<<40 here but futex woke up (-ETIMEDOUT).
// I wonder if the timespec that gets to the kernel
// actually has two 32-bit numbers in it, so tha
// a 64-bit 1<<40 ends up being 0 seconds,
// 1<<8 nanoseconds.
static struct timespec longtime =
{
	1<<30,	// 34 years
	0
};

// Atomically,
//	if(*addr == val) sleep
// Might be woken up spuriously; that's allowed.
static void
futexsleep(uint32 *addr, uint32 val)
{
	int64 ret;

	ret = futex(addr, FUTEX_WAIT, val, &longtime, nil, 0);
	if(ret >= 0 || ret == -EAGAIN || ret == -EINTR)
		return;

	prints("futexsleep addr=");
	sys·printpointer(addr);
	prints(" val=");
	sys·printint(val);
	prints(" returned ");
	sys·printint(ret);
	prints("\n");
	*(int32*)0x1005 = 0x1005;
}

// If any procs are sleeping on addr, wake up at least one.
static void
futexwakeup(uint32 *addr)
{
	int64 ret;

	ret = futex(addr, FUTEX_WAKE, 1, nil, nil, 0);

	if(ret >= 0)
		return;

	// I don't know that futex wakeup can return
	// EAGAIN or EINTR, but if it does, it would be
	// safe to loop and call futex again.

	prints("futexwakeup addr=");
	sys·printpointer(addr);
	prints(" returned ");
	sys·printint(ret);
	prints("\n");
	*(int32*)0x1006 = 0x1006;
}


// Lock and unlock.
//
// The lock state is a single 32-bit word that holds
// a 31-bit count of threads waiting for the lock
// and a single bit (the low bit) saying whether the lock is held.
// The uncontended case runs entirely in user space.
// When contention is detected, we defer to the kernel (futex).
//
// A reminder: compare-and-swap cas(addr, old, new) does
//	if(*addr == old) { *addr = new; return 1; }
// 	else return 0;
// but atomically.

void
lock(Lock *l)
{
	uint32 v;

again:
	v = l->key;
	if((v&1) == 0){
		if(cas(&l->key, v, v|1)){
			// Lock wasn't held; we grabbed it.
			return;
		}
		goto again;
	}

	// Lock was held; try to add ourselves to the waiter count.
	if(!cas(&l->key, v, v+2))
		goto again;

	// We're accounted for, now sleep in the kernel.
	//
	// We avoid the obvious lock/unlock race because
	// the kernel won't put us to sleep if l->key has
	// changed underfoot and is no longer v+2.
	//
	// We only really care that (v&1) == 1 (the lock is held),
	// and in fact there is a futex variant that could
	// accomodate that check, but let's not get carried away.)
	futexsleep(&l->key, v+2);

	// We're awake: remove ourselves from the count.
	for(;;){
		v = l->key;
		if(v < 2)
			throw("bad lock key");
		if(cas(&l->key, v, v-2))
			break;
	}

	// Try for the lock again.
	goto again;
}

void
unlock(Lock *l)
{
	uint32 v;

	// Atomically get value and clear lock bit.
again:
	v = l->key;
	if((v&1) == 0)
		throw("unlock of unlocked lock");
	if(!cas(&l->key, v, v&~1))
		goto again;

	// If there were waiters, wake one.
	if(v & ~1)
		futexwakeup(&l->key);
}


// One-time notifications.
//
// Since the lock/unlock implementation already
// takes care of sleeping in the kernel, we just reuse it.
// (But it's a weird use, so it gets its own interface.)
//
// We use a lock to represent the event:
// unlocked == event has happened.
// Thus the lock starts out locked, and to wait for the
// event you try to lock the lock.  To signal the event,
// you unlock the lock.

void
noteclear(Note *n)
{
	n->lock.key = 0;	// memset(n, 0, sizeof *n)
	lock(&n->lock);
}

void
notewakeup(Note *n)
{
	unlock(&n->lock);
}

void
notesleep(Note *n)
{
	lock(&n->lock);
	unlock(&n->lock);	// Let other sleepers find out too.
}


// Clone, the Linux rfork.
enum
{
	CLONE_VM = 0x100,
	CLONE_FS = 0x200,
	CLONE_FILES = 0x400,
	CLONE_SIGHAND = 0x800,
	CLONE_PTRACE = 0x2000,
	CLONE_VFORK = 0x4000,
	CLONE_PARENT = 0x8000,
	CLONE_THREAD = 0x10000,
	CLONE_NEWNS = 0x20000,
	CLONE_SYSVSEM = 0x40000,
	CLONE_SETTLS = 0x80000,
	CLONE_PARENT_SETTID = 0x100000,
	CLONE_CHILD_CLEARTID = 0x200000,
	CLONE_UNTRACED = 0x800000,
	CLONE_CHILD_SETTID = 0x1000000,
	CLONE_STOPPED = 0x2000000,
	CLONE_NEWUTS = 0x4000000,
	CLONE_NEWIPC = 0x8000000,
};

void
newosproc(M *m, G *g, void *stk, void (*fn)(void))
{
	int64 ret;
	int32 flags;

	flags = CLONE_PARENT	/* getppid doesn't change in child */
		| CLONE_VM	/* share memory */
		| CLONE_FS	/* share cwd, etc */
		| CLONE_FILES	/* share fd table */
		| CLONE_SIGHAND	/* share sig handler table */
		| CLONE_PTRACE	/* revisit - okay for now */
		| CLONE_THREAD	/* revisit - okay for now */
		;

	if(0){
		prints("newosproc stk=");
		sys·printpointer(stk);
		prints(" m=");
		sys·printpointer(m);
		prints(" g=");
		sys·printpointer(g);
		prints(" fn=");
		sys·printpointer(fn);
		prints(" clone=");
		sys·printpointer(clone);
		prints("\n");
	}

	ret = clone(flags, stk, m, g, fn);
	if(ret < 0)
		*(int32*)123 = 123;
}

void
sys·sleep(int64 ms)
{
	struct timeval tv;

	tv.tv_sec = ms/1000;
	tv.tv_usec = ms%1000 * 1000;
	select(0, nil, nil, nil, &tv);
}

void
osinit(void)
{
}
