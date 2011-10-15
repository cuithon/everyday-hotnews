// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "runtime.h"
#include "type.h"
#include "defs.h"
#include "os.h"

#pragma dynimport runtime·CloseHandle CloseHandle "kernel32.dll"
#pragma dynimport runtime·CreateEvent CreateEventA "kernel32.dll"
#pragma dynimport runtime·CreateThread CreateThread "kernel32.dll"
#pragma dynimport runtime·CreateWaitableTimer CreateWaitableTimerA "kernel32.dll"
#pragma dynimport runtime·DuplicateHandle DuplicateHandle "kernel32.dll"
#pragma dynimport runtime·ExitProcess ExitProcess "kernel32.dll"
#pragma dynimport runtime·FreeEnvironmentStringsW FreeEnvironmentStringsW "kernel32.dll"
#pragma dynimport runtime·GetEnvironmentStringsW GetEnvironmentStringsW "kernel32.dll"
#pragma dynimport runtime·GetProcAddress GetProcAddress "kernel32.dll"
#pragma dynimport runtime·GetStdHandle GetStdHandle "kernel32.dll"
#pragma dynimport runtime·GetSystemInfo GetSystemInfo "kernel32.dll"
#pragma dynimport runtime·GetThreadContext GetThreadContext "kernel32.dll"
#pragma dynimport runtime·LoadLibrary LoadLibraryW "kernel32.dll"
#pragma dynimport runtime·QueryPerformanceCounter QueryPerformanceCounter "kernel32.dll"
#pragma dynimport runtime·QueryPerformanceFrequency QueryPerformanceFrequency "kernel32.dll"
#pragma dynimport runtime·ResumeThread ResumeThread "kernel32.dll"
#pragma dynimport runtime·SetConsoleCtrlHandler SetConsoleCtrlHandler "kernel32.dll"
#pragma dynimport runtime·SetEvent SetEvent "kernel32.dll"
#pragma dynimport runtime·SetThreadPriority SetThreadPriority "kernel32.dll"
#pragma dynimport runtime·SetWaitableTimer SetWaitableTimer "kernel32.dll"
#pragma dynimport runtime·Sleep Sleep "kernel32.dll"
#pragma dynimport runtime·SuspendThread SuspendThread "kernel32.dll"
#pragma dynimport runtime·timeBeginPeriod timeBeginPeriod "winmm.dll"
#pragma dynimport runtime·WaitForSingleObject WaitForSingleObject "kernel32.dll"
#pragma dynimport runtime·WriteFile WriteFile "kernel32.dll"

extern void *runtime·CloseHandle;
extern void *runtime·CreateEvent;
extern void *runtime·CreateThread;
extern void *runtime·CreateWaitableTimer;
extern void *runtime·DuplicateHandle;
extern void *runtime·ExitProcess;
extern void *runtime·FreeEnvironmentStringsW;
extern void *runtime·GetEnvironmentStringsW;
extern void *runtime·GetProcAddress;
extern void *runtime·GetStdHandle;
extern void *runtime·GetSystemInfo;
extern void *runtime·GetThreadContext;
extern void *runtime·LoadLibrary;
extern void *runtime·QueryPerformanceCounter;
extern void *runtime·QueryPerformanceFrequency;
extern void *runtime·ResumeThread;
extern void *runtime·SetConsoleCtrlHandler;
extern void *runtime·SetEvent;
extern void *runtime·SetThreadPriority;
extern void *runtime·SetWaitableTimer;
extern void *runtime·Sleep;
extern void *runtime·SuspendThread;
extern void *runtime·timeBeginPeriod;
extern void *runtime·WaitForSingleObject;
extern void *runtime·WriteFile;

static int64 timerfreq;

static int32
getproccount(void)
{
	SystemInfo info;

	runtime·stdcall(runtime·GetSystemInfo, 1, &info);
	return info.dwNumberOfProcessors;
}

void
runtime·osinit(void)
{
	// -1 = current process, -2 = current thread
	runtime·stdcall(runtime·DuplicateHandle, 7,
		(uintptr)-1, (uintptr)-2, (uintptr)-1, &m->thread,
		(uintptr)0, (uintptr)0, (uintptr)DUPLICATE_SAME_ACCESS);
	runtime·stdcall(runtime·QueryPerformanceFrequency, 1, &timerfreq);
	runtime·stdcall(runtime·SetConsoleCtrlHandler, 2, runtime·ctrlhandler, (uintptr)1);
	runtime·stdcall(runtime·timeBeginPeriod, 1, (uintptr)1);
	runtime·ncpu = getproccount();
}

void
runtime·goenvs(void)
{
	extern Slice os·Envs;

	uint16 *env;
	String *s;
	int32 i, n;
	uint16 *p;

	env = runtime·stdcall(runtime·GetEnvironmentStringsW, 0);

	n = 0;
	for(p=env; *p; n++)
		p += runtime·findnullw(p)+1;

	s = runtime·malloc(n*sizeof s[0]);

	p = env;
	for(i=0; i<n; i++) {
		s[i] = runtime·gostringw(p);
		p += runtime·findnullw(p)+1;
	}
	os·Envs.array = (byte*)s;
	os·Envs.len = n;
	os·Envs.cap = n;

	runtime·stdcall(runtime·FreeEnvironmentStringsW, 1, env);
}

void
runtime·exit(int32 code)
{
	runtime·stdcall(runtime·ExitProcess, 1, (uintptr)code);
}

int32
runtime·write(int32 fd, void *buf, int32 n)
{
	void *handle;
	uint32 written;

	written = 0;
	switch(fd) {
	case 1:
		handle = runtime·stdcall(runtime·GetStdHandle, 1, (uintptr)-11);
		break;
	case 2:
		handle = runtime·stdcall(runtime·GetStdHandle, 1, (uintptr)-12);
		break;
	default:
		return -1;
	}
	runtime·stdcall(runtime·WriteFile, 5, handle, buf, (uintptr)n, &written, (uintptr)0);
	return written;
}

void
runtime·osyield(void)
{
	runtime·stdcall(runtime·Sleep, 1, (uintptr)0);
}

void
runtime·usleep(uint32 us)
{
	us /= 1000;
	if(us == 0)
		us = 1;
	runtime·stdcall(runtime·Sleep, 1, (uintptr)us);
}

// Thread-safe allocation of an event.
static void
initevent(void **pevent)
{
	void *event;

	event = runtime·stdcall(runtime·CreateEvent, 4, (uintptr)0, (uintptr)0, (uintptr)0, (uintptr)0);
	if(!runtime·casp(pevent, 0, event)) {
		// Someone else filled it in.  Use theirs.
		runtime·stdcall(runtime·CloseHandle, 1, event);
	}
}

#define LOCK_HELD ((M*)-1)

static void
eventlock(Lock *l)
{
	// Allocate event if needed.
	if(m->event == nil)
		initevent(&m->event);

	for(;;) {
		m->nextwaitm = runtime·atomicloadp(&l->waitm);
		if(m->nextwaitm == nil) {
			if(runtime·casp(&l->waitm, nil, LOCK_HELD))
				return;
		// Someone else has it.
		// l->waitm points to a linked list of M's waiting
		// for this lock, chained through m->nextwaitm.
		// Queue this M.
		} else if(runtime·casp(&l->waitm, m->nextwaitm, m))
			break;
	}

	// Wait.
	runtime·stdcall(runtime·WaitForSingleObject, 2, m->event, (uintptr)-1);
}

static void
eventunlock(Lock *l)
{
	M *mp;

	for(;;) {
		mp = runtime·atomicloadp(&l->waitm);
		if(mp == LOCK_HELD) {
			if(runtime·casp(&l->waitm, LOCK_HELD, nil))
				return;
		// Other M's are waiting for the lock.
		// Dequeue a M.
		} else if(runtime·casp(&l->waitm, mp, mp->nextwaitm))
			break;
	}

	// Wake that M.
	runtime·stdcall(runtime·SetEvent, 1, mp->event);
}

void
runtime·lock(Lock *l)
{
	if(m->locks < 0)
		runtime·throw("lock count");
	m->locks++;
	eventlock(l);
}

void
runtime·unlock(Lock *l)
{
	m->locks--;
	if(m->locks < 0)
		runtime·throw("lock count");
	eventunlock(l);
}

void
runtime·noteclear(Note *n)
{
	n->lock.waitm = nil;
	eventlock(&n->lock);
}

void
runtime·notewakeup(Note *n)
{
	eventunlock(&n->lock);
}

void
runtime·notesleep(Note *n)
{
	eventlock(&n->lock);
	eventunlock(&n->lock);	// Let other sleepers find out too.
}

void
runtime·newosproc(M *m, G *g, void *stk, void (*fn)(void))
{
	void *thandle;

	USED(stk);
	USED(g);	// assuming g = m->g0
	USED(fn);	// assuming fn = mstart

	thandle = runtime·stdcall(runtime·CreateThread, 6,
		nil, nil, runtime·tstart_stdcall, m, nil, nil);
	if(thandle == nil) {
		runtime·printf("runtime: failed to create new OS thread (have %d already; errno=%d)\n", runtime·mcount(), runtime·getlasterror());
		runtime·throw("runtime.newosproc");
	}
	runtime·atomicstorep(&m->thread, thandle);
}

// Called to initialize a new m (including the bootstrap m).
void
runtime·minit(void)
{
}

void
runtime·gettime(int64 *sec, int32 *usec)
{
	int64 count;

	runtime·stdcall(runtime·QueryPerformanceCounter, 1, &count);
	*sec = count / timerfreq;
	count %= timerfreq;
	*usec = count*1000000 / timerfreq;
}

// Calling stdcall on os stack.
#pragma textflag 7
void *
runtime·stdcall(void *fn, int32 count, ...)
{
	WinCall c;

	c.fn = fn;
	c.n = count;
	c.args = (uintptr*)&count + 1;
	runtime·asmcgocall(runtime·asmstdcall, &c);
	return (void*)c.r1;
}

uint32
runtime·issigpanic(uint32 code)
{
	switch(code) {
	case EXCEPTION_ACCESS_VIOLATION:
	case EXCEPTION_INT_DIVIDE_BY_ZERO:
	case EXCEPTION_INT_OVERFLOW:
	case EXCEPTION_FLT_DENORMAL_OPERAND:
	case EXCEPTION_FLT_DIVIDE_BY_ZERO:
	case EXCEPTION_FLT_INEXACT_RESULT:
	case EXCEPTION_FLT_OVERFLOW:
	case EXCEPTION_FLT_UNDERFLOW:
		return 1;
	}
	return 0;
}

void
runtime·sigpanic(void)
{
	switch(g->sig) {
	case EXCEPTION_ACCESS_VIOLATION:
		if(g->sigcode1 < 0x1000)
			runtime·panicstring("invalid memory address or nil pointer dereference");
		runtime·printf("unexpected fault address %p\n", g->sigcode1);
		runtime·throw("fault");
	case EXCEPTION_INT_DIVIDE_BY_ZERO:
		runtime·panicstring("integer divide by zero");
	case EXCEPTION_INT_OVERFLOW:
		runtime·panicstring("integer overflow");
	case EXCEPTION_FLT_DENORMAL_OPERAND:
	case EXCEPTION_FLT_DIVIDE_BY_ZERO:
	case EXCEPTION_FLT_INEXACT_RESULT:
	case EXCEPTION_FLT_OVERFLOW:
	case EXCEPTION_FLT_UNDERFLOW:
		runtime·panicstring("floating point error");
	}
	runtime·throw("fault");
}

String
runtime·signame(int32 sig)
{
	int8 *s;

	switch(sig) {
	case SIGINT:
		s = "SIGINT: interrupt";
		break;
	default:
		return runtime·emptystring;
	}
	return runtime·gostringnocopy((byte*)s);
}

uint32
runtime·ctrlhandler1(uint32 type)
{
	int32 s;

	switch(type) {
	case CTRL_C_EVENT:
	case CTRL_BREAK_EVENT:
		s = SIGINT;
		break;
	default:
		return 0;
	}

	if(runtime·sigsend(s))
		return 1;
	runtime·exit(2);	// SIGINT, SIGTERM, etc
	return 0;
}

extern void runtime·dosigprof(Context *r, G *gp);
extern void runtime·profileloop(void);
static void *profiletimer;

static void
profilem(M *mp)
{
	extern M runtime·m0;
	extern uint32 runtime·tls0[];
	byte rbuf[sizeof(Context)+15];
	Context *r;
	void *tls;
	G *gp;

	tls = mp->tls;
	if(mp == &runtime·m0)
		tls = runtime·tls0;
	gp = *(G**)tls;

	if(gp != nil && gp != mp->g0 && gp->status != Gsyscall) {
		// align Context to 16 bytes
		r = (Context*)((uintptr)(&rbuf[15]) & ~15);
		r->ContextFlags = CONTEXT_CONTROL;
		runtime·stdcall(runtime·GetThreadContext, 2, mp->thread, r);
		runtime·dosigprof(r, gp);
	}
}

void
runtime·profileloop1(void)
{
	M *mp, *allm;
	void *thread;

	runtime·stdcall(runtime·SetThreadPriority, 2,
		(uintptr)-2, (uintptr)THREAD_PRIORITY_HIGHEST);

	for(;;) {
		runtime·stdcall(runtime·WaitForSingleObject, 2, profiletimer, (uintptr)-1);
		allm = runtime·atomicloadp(&runtime·allm);
		for(mp = allm; mp != nil; mp = mp->alllink) {
			thread = runtime·atomicloadp(&mp->thread);
			if(thread == nil)
				continue;
			runtime·stdcall(runtime·SuspendThread, 1, thread);
			if(mp->profilehz != 0)
				profilem(mp);
			runtime·stdcall(runtime·ResumeThread, 1, thread);
		}
	}
}

void
runtime·resetcpuprofiler(int32 hz)
{
	static Lock lock;
	void *timer, *thread;
	int32 ms;
	int64 due;

	runtime·lock(&lock);
	if(profiletimer == nil) {
		timer = runtime·stdcall(runtime·CreateWaitableTimer, 3, nil, nil, nil);
		runtime·atomicstorep(&profiletimer, timer);
		thread = runtime·stdcall(runtime·CreateThread, 6,
			nil, nil, runtime·profileloop, nil, nil, nil);
		runtime·stdcall(runtime·CloseHandle, 1, thread);
	}
	runtime·unlock(&lock);

	ms = 0;
	due = 1LL<<63;
	if(hz > 0) {
		ms = 1000 / hz;
		if(ms == 0)
			ms = 1;
		due = ms * -10000;
	}
	runtime·stdcall(runtime·SetWaitableTimer, 6,
		profiletimer, &due, (uintptr)ms, nil, nil, nil);
	runtime·atomicstore((uint32*)&m->profilehz, hz);
}

void
os·sigpipe(void)
{
	runtime·throw("too many writes on closed pipe");
}
