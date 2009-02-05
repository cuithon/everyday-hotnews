// Derived from Plan 9 from User Space src/libmach/Linux.c
// http://code.swtch.com/plan9port/src/tip/src/libmach/Linux.c
//
//	Copyright © 1994-1999 Lucent Technologies Inc.
//	Power PC support Copyright © 1995-2004 C H Forsyth (forsyth@terzarima.net).
//	Portions Copyright © 1997-1999 Vita Nuova Limited.
//	Portions Copyright © 2000-2007 Vita Nuova Holdings Limited (www.vitanuova.com).
//	Revisions Copyright © 2000-2004 Lucent Technologies Inc. and others.
//	Portions Copyright © 2001-2007 Russ Cox.
//	Portions Copyright © 2009 The Go Authors.  All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

#include <u.h>
#include <sys/syscall.h>	/* for tkill */
#include <unistd.h>
#include <dirent.h>
#include <sys/ptrace.h>
#include <sys/signal.h>
#include <sys/wait.h>
#include <errno.h>
#include <libc.h>
#include <bio.h>
#include <mach_amd64.h>
#include <ureg_amd64.h>
#undef waitpid

// The old glibc used with crosstool compilers on thresher
// doesn't know these numbers, but the Linux kernel
// had them as far back as 2.6.0.
#ifndef WSTOPPED
#define WSTOPPED 2
#define WCONTINUED 8
#define WIFCONTINUED(x) ((x) == 0xffff)
#endif
#ifndef PTRACE_SETOPTIONS
#define PTRACE_SETOPTIONS 0x4200
#define PTRACE_GETEVENTMSG 0x4201
#define PTRACE_O_TRACEFORK 0x2
#define PTRACE_O_TRACEVFORK 0x4
#define PTRACE_O_TRACECLONE 0x8
#define PTRACE_O_TRACEEXEC 0x10
#define PTRACE_O_TRACEVFORKDONE 0x20
#define PTRACE_O_TRACEEXIT 0x40
#define PTRACE_EVENT_FORK 0x1
#define PTRACE_EVENT_VFORK 0x2
#define PTRACE_EVENT_CLONE 0x3
#define PTRACE_EVENT_EXEC 0x4
#define PTRACE_EVENT_VFORK_DONE 0x5
#define PTRACE_EVENT_EXIT 0x6
#endif

typedef struct Ureg Ureg;

static Maprw ptracesegrw;
static Maprw ptraceregrw;

// /usr/include/asm-x86_64/user.h
struct user_regs_struct {
	unsigned long r15,r14,r13,r12,rbp,rbx,r11,r10;
	unsigned long r9,r8,rax,rcx,rdx,rsi,rdi,orig_rax;
	unsigned long rip,cs,eflags;
	unsigned long rsp,ss;
  	unsigned long fs_base, gs_base;
	unsigned long ds,es,fs,gs;
};

// Linux gets very upset if a debugger forgets the reported state
// of a debugged process, so we keep everything we know about
// a debugged process in the LinuxThread structure.
//
// We can poll for state changes by calling waitpid and interpreting
// the integer status code that comes back.  Wait1 does this.
//
// If the process is already running, it is an error to PTRACE_CONT it.
//
// If the process is already stopped, it is an error to stop it again.
//
// If the process is stopped because of a signal, the debugger must
// relay the signal to the PTRACE_CONT call, or else the signal is
// dropped.
//
// If the process exits, the debugger should detach so that the real
// parent can reap the zombie.
//
// On first attach, the debugger should set a handful of flags in order
// to catch future events like fork, clone, exec, etc.

// One for every attached thread.
typedef struct LinuxThread LinuxThread;
struct LinuxThread
{
	int pid;
	int tid;
	int state;
	int signal;
	int child;
	int exitcode;
};

static int trace = 0;

static LinuxThread **thr;
static int nthr;
static int mthr;

static int realpid(int pid);

enum
{
	Unknown,
	Detached,
	Attached,
	AttachStop,
	Stopped,
	Running,
	Forking,
	Vforking,
	VforkDone,
	Cloning,
	Execing,
	Exiting,
	Exited,
	Killed,

	NSTATE,
};

static char* statestr[NSTATE] = {
	"Unknown",
	"Detached",
	"Attached",
	"AttachStop",
	"Stopped",
	"Running",
	"Forking",
	"Vforking",
	"VforkDone",
	"Cloning",
	"Execing",
	"Exiting",
	"Exited",
	"Killed"
};

static LinuxThread*
attachthread(int pid, int tid, int *new, int newstate)
{
	int i, n, status;
	LinuxThread **p, *t;
	uintptr flags;

	if(new)
		*new = 0;

	for(i=0; i<nthr; i++)
		if((pid == 0 || thr[i]->pid == pid) && thr[i]->tid == tid) {
			t = thr[i];
			goto fixup;
		}

	if(!new)
		return nil;

	if(nthr >= mthr) {
		n = mthr;
		if(n == 0)
			n = 64;
		else
			n *= 2;
		p = realloc(thr, n*sizeof thr[0]);
		if(p == nil)
			return nil;
		thr = p;
		mthr = n;
	}

	t = malloc(sizeof *t);
	if(t == nil)
		return nil;

	thr[nthr++] = t;
	t->pid = pid;
	t->tid = tid;
	t->state = newstate;
	if(trace)
		fprint(2, "new thread %d %d\n", t->pid, t->tid);
	if(new)
		*new = 1;

fixup:
	if(t->state == Detached) {
		if(ptrace(PTRACE_ATTACH, tid, 0, 0) < 0) {
			fprint(2, "ptrace ATTACH %d: %r\n", tid);
			return nil;
		}
		t->state = Attached;
	}

	if(t->state == Attached) {
		// wait for stop, so we can set options
		if(waitpid(tid, &status, __WALL|WUNTRACED|WSTOPPED) < 0)
			return nil;
		if(!WIFSTOPPED(status)) {
			fprint(2, "waitpid %d: status=%#x not stopped\n", tid);
			return nil;
		}
		t->state = AttachStop;
	}

	if(t->state == AttachStop) {
		// set options so we'll find out about new threads
		flags = PTRACE_O_TRACEFORK |
			PTRACE_O_TRACEVFORK |
			PTRACE_O_TRACECLONE |
			PTRACE_O_TRACEEXEC |
			PTRACE_O_TRACEVFORKDONE |
			PTRACE_O_TRACEEXIT;
		if(ptrace(PTRACE_SETOPTIONS, tid, 0, (void*)flags) < 0)	{
			fprint(2, "ptrace PTRACE_SETOPTIONS %d: %r\n", tid);
			return nil;
		}
		t->state = Stopped;
	}

	return t;
}

static LinuxThread*
findthread(int tid)
{
	return attachthread(0, tid, nil, 0);
}

int
procthreadpids(int pid, int *p, int np)
{
	int i, n;
	LinuxThread *t;

	n = 0;
	for(i=0; i<nthr; i++) {
		t = thr[i];
		if(t->pid == pid) {
			switch(t->state) {
			case Exited:
			case Detached:
			case Killed:
				break;

			default:
				if(n < np)
					p[n] = t->tid;
				n++;
				break;
			}
		}
	}
	return n;
}

// Execute a single wait and update the corresponding thread.
static int
wait1(int nohang)
{
	int tid, new, status, event;
	ulong data;
	LinuxThread *t;
	enum
	{
		NormalStop = 0x137f
	};

	if(nohang != 0)
		nohang = WNOHANG;

	tid = waitpid(-1, &status, __WALL|WUNTRACED|WSTOPPED|WCONTINUED|nohang);
	if(tid < 0)
		return -1;
	if(tid == 0)
		return 0;

	if(trace > 0 && status != NormalStop)
		fprint(2, "TID %d: %#x\n", tid, status);

	// If we've not heard of this tid, something is wrong.
	t = findthread(tid);
	if(t == nil) {
		fprint(2, "ptrace waitpid: unexpected new tid %d status %#x\n", tid, status);
		return -1;
	}

	if(WIFSTOPPED(status)) {
		t->state = Stopped;
		t->signal = WSTOPSIG(status);
		if(trace)
			fprint(2, "tid %d: stopped %#x%s\n", tid, status,
				status != NormalStop ? " ***" : "");
		if(t->signal == SIGTRAP && (event = status>>16) != 0) {	// ptrace event
			switch(event) {
			case PTRACE_EVENT_FORK:
				t->state = Forking;
				goto child;

			case PTRACE_EVENT_VFORK:
				t->state = Vforking;
				goto child;

			case PTRACE_EVENT_CLONE:
				t->state = Cloning;
				goto child;

			child:
				if(ptrace(PTRACE_GETEVENTMSG, t->tid, 0, &data) < 0) {
					fprint(2, "ptrace GETEVENTMSG tid %d: %r\n", tid);
					break;
				}
				t->child = data;
				attachthread(t->pid, t->child, &new, Running);
				if(!new)
					fprint(2, "ptrace child: not new\n");
				break;

			case PTRACE_EVENT_EXEC:
				t->state = Execing;
				break;

			case PTRACE_EVENT_VFORK_DONE:
				t->state = VforkDone;
				break;

			case PTRACE_EVENT_EXIT:
				if(trace)
					fprint(2, "tid %d: exiting %#x\n", tid, status);
				t->state = Exiting;
				if(ptrace(PTRACE_GETEVENTMSG, t->tid, 0, &data) < 0) {
					fprint(2, "ptrace GETEVENTMSG tid %d: %r\n", tid);
					break;
				}
				t->exitcode = data;
				break;
			}
		}
	}
	if(WIFCONTINUED(status)) {
		if(trace)
			fprint(2, "tid %d: continued %#x\n", tid, status);
		t->state = Running;
	}
	if(WIFEXITED(status)) {
		if(trace)
			fprint(2, "tid %d: exited %#x\n", tid, status);
		t->state = Exited;
		t->exitcode = WEXITSTATUS(status);
		t->signal = -1;
		ptrace(PTRACE_DETACH, t->tid, 0, 0);
		if(trace)
			fprint(2, "tid %d: detach exited\n", tid);
	}
	if(WIFSIGNALED(status)) {
		if(trace)
			fprint(2, "tid %d: signaled %#x\n", tid, status);
		t->state = Exited;
		t->signal = WTERMSIG(status);
		t->exitcode = -1;
		ptrace(PTRACE_DETACH, t->tid, 0, 0);
		if(trace)
			fprint(2, "tid %d: detach signaled\n", tid);
	}
	return 1;
}

static int
waitstop(LinuxThread *t)
{
	while(t->state == Running)
		if(wait1(0) < 0)
			return -1;
	return 0;
}

// Attach to and stop all threads in process pid.
// Must stop everyone in order to make sure we set
// the "tell me about new threads" option in every
// task.
int
attachallthreads(int pid)
{
	int tid, foundnew, new;
	char buf[100];
	DIR *d;
	struct dirent *de;
	LinuxThread *t;

	if(pid == 0) {
		fprint(2, "attachallthreads(0)\n");
		return -1;
	}

	pid = realpid(pid);

	snprint(buf, sizeof buf, "/proc/%d/task", pid);
	if((d = opendir(buf)) == nil) {
		fprint(2, "opendir %s: %r\n", buf);
		return -1;
	}

	// Loop in case new threads are being created right now.
	// We stop every thread as we find it, so eventually
	// this has to stop (or the system runs out of procs).
	do {
		foundnew = 0;
		while((de = readdir(d)) != nil) {
			tid = atoi(de->d_name);
			if(tid == 0)
				continue;
			t = attachthread(pid, tid, &new, Detached);
			foundnew |= new;
			if(t)
				waitstop(t);
		}
		rewinddir(d);
	} while(foundnew);
	closedir(d);

	return 0;
}

Map*
attachproc(int pid, Fhdr *fp)
{
	Map *map;

	if(pid == 0) {
		fprint(2, "attachproc(0)\n");
		return nil;
	}

	if(findthread(pid) == nil && attachallthreads(pid) < 0)
		return nil;

	map = newmap(0, 4);
	if (!map)
		return 0;
	map->pid = pid;
	if(mach->regsize)
		setmap(map, -1, 0, mach->regsize, 0, "regs", ptraceregrw);
//	if(mach->fpregsize)
//		setmap(map, -1, mach->regsize, mach->regsize+mach->fpregsize, 0, "fpregs", ptraceregrw);
	setmap(map, -1, fp->txtaddr, fp->txtaddr+fp->txtsz, fp->txtaddr, "*text", ptracesegrw);
	setmap(map, -1, fp->dataddr, mach->utop, fp->dataddr, "*data", ptracesegrw);
	return map;
}

void
detachproc(Map *m)
{
	LinuxThread *t;

	t = findthread(m->pid);
	if(t != nil) {
		ptrace(PTRACE_DETACH, t->tid, 0, 0);
		t->state = Detached;
		if(trace)
			fprint(2, "tid %d: detachproc\n", t->tid);
		// TODO(rsc): Reclaim thread structs somehow?
	}
	free(m);
}

/* /proc/pid/stat contains
	pid
	command in parens
	0. state
	1. ppid
	2. pgrp
	3. session
	4. tty_nr
	5. tpgid
	6. flags (math=4, traced=10)
	7. minflt
	8. cminflt
	9. majflt
	10. cmajflt
	11. utime
	12. stime
	13. cutime
	14. cstime
	15. priority
	16. nice
	17. 0
	18. itrealvalue
	19. starttime
	20. vsize
	21. rss
	22. rlim
	23. startcode
	24. endcode
	25. startstack
	26. kstkesp
	27. kstkeip
	28. pending signal bitmap
	29. blocked signal bitmap
	30. ignored signal bitmap
	31. caught signal bitmap
	32. wchan
	33. nswap
	34. cnswap
	35. exit_signal
	36. processor
*/

static int
readstat(int pid, char *buf, int nbuf, char **f, int nf)
{
	int fd, n;
	char *p;

	snprint(buf, nbuf, "/proc/%d/stat", pid);
	if((fd = open(buf, OREAD)) < 0){
		fprint(2, "open %s: %r\n", buf);
		return -1;
	}
	n = read(fd, buf, nbuf-1);
	close(fd);
	if(n <= 0){
		fprint(2, "read %s: %r\n", buf);
		return -1;
	}
	buf[n] = 0;

	/* command name is in parens, no parens afterward */
	p = strrchr(buf, ')');
	if(p == nil || *++p != ' '){
		fprint(2, "bad format in /proc/%d/stat\n", pid);
		return -1;
	}
	++p;

	nf = tokenize(p, f, nf);
	if(0) print("code 0x%lux-0x%lux stack 0x%lux kstk 0x%lux keip 0x%lux pending 0x%lux\n",
		strtoul(f[23], 0, 0), strtoul(f[24], 0, 0), strtoul(f[25], 0, 0),
		strtoul(f[26], 0, 0), strtoul(f[27], 0, 0), strtoul(f[28], 0, 0));

	return nf;
}

static char*
readstatus(int pid, char *buf, int nbuf, char *key)
{
	int fd, n;
	char *p;

	snprint(buf, nbuf, "/proc/%d/status", pid);
	if((fd = open(buf, OREAD)) < 0){
		fprint(2, "open %s: %r\n", buf);
		return nil;
	}
	n = read(fd, buf, nbuf-1);
	close(fd);
	if(n <= 0){
		fprint(2, "read %s: %r\n", buf);
		return nil;
	}
	buf[n] = 0;
	p = strstr(buf, key);
	if(p)
		return p+strlen(key);
	return nil;
}

int
procnotes(int pid, char ***pnotes)
{
	char buf[1024], *f[40];
	int i, n, nf;
	char *s, **notes;
	ulong sigs;
	extern char *_p9sigstr(int, char*);

	*pnotes = nil;
	nf = readstat(pid, buf, sizeof buf, f, nelem(f));
	if(nf <= 28)
		return -1;

	sigs = strtoul(f[28], 0, 0) & ~(1<<SIGCONT);
	if(sigs == 0){
		*pnotes = nil;
		return 0;
	}

	notes = mallocz(32*sizeof(char*), 0);
	if(notes == nil)
		return -1;
	n = 0;
	for(i=0; i<32; i++){
		if((sigs&(1<<i)) == 0)
			continue;
		if((s = _p9sigstr(i, nil)) == nil)
			continue;
		notes[n++] = s;
	}
	*pnotes = notes;
	return n;
}

static int
realpid(int pid)
{
	char buf[1024], *p;

	p = readstatus(pid, buf, sizeof buf, "\nTgid:");
	if(p == nil)
		return pid;
	return atoi(p);
}

int
ctlproc(int pid, char *msg)
{
	int new;
	LinuxThread *t;
	uintptr data;

	while(wait1(1) > 0)
		;

	if(strcmp(msg, "attached") == 0){
		t = attachthread(pid, pid, &new, Attached);
		if(t == nil)
			return -1;
		return 0;
	}

	if(strcmp(msg, "hang") == 0){
		if(pid == getpid())
			return ptrace(PTRACE_TRACEME, 0, 0, 0);
		werrstr("can only hang self");
		return -1;
	}

	t = findthread(pid);
	if(t == nil) {
		werrstr("not attached to pid %d", pid);
		return -1;
	}
	if(t->state == Exited) {
		werrstr("pid %d has exited", pid);
		return -1;
	}
	if(t->state == Killed) {
		werrstr("pid %d has been killed", pid);
		return -1;
	}

	if(strcmp(msg, "kill") == 0) {
		if(ptrace(PTRACE_KILL, pid, 0, 0) < 0)
			return -1;
		t->state = Killed;
		return 0;
	}
	if(strcmp(msg, "startstop") == 0){
		if(ctlproc(pid, "start") < 0)
			return -1;
		return waitstop(t);
	}
	if(strcmp(msg, "sysstop") == 0){
		if(ptrace(PTRACE_SYSCALL, pid, 0, 0) < 0)
			return -1;
		t->state = Running;
		return waitstop(t);
	}
	if(strcmp(msg, "stop") == 0){
		if(trace > 1)
			fprint(2, "tid %d: tkill stop\n", pid);
		if(t->state == Stopped)
			return 0;
		if(syscall(__NR_tkill, pid, SIGSTOP) < 0)
			return -1;
		return waitstop(t);
	}
	if(strcmp(msg, "step") == 0){
		if(t->state == Running) {
			werrstr("cannot single-step unstopped %d", pid);
			return -1;
		}
		if(ptrace(PTRACE_SINGLESTEP, pid, 0, 0) < 0)
			return -1;
		return waitstop(t);
	}
	if(strcmp(msg, "start") == 0) {
		if(t->state == Running)
			return 0;
		data = 0;
		if(t->state == Stopped && t->signal != SIGSTOP)
			data = t->signal;
		if(trace && data)
			fprint(2, "tid %d: continue %lud\n", pid, (ulong)data);
		if(ptrace(PTRACE_CONT, pid, 0, (void*)data) < 0)
			return -1;
		t->state = Running;
		return 0;
	}
	if(strcmp(msg, "waitstop") == 0) {
		return waitstop(t);
	}
	werrstr("unknown control message '%s'", msg);
	return -1;
}

char*
proctextfile(int pid)
{
	static char buf[1024], pbuf[128];

	snprint(pbuf, sizeof pbuf, "/proc/%d/exe", pid);
	if(readlink(pbuf, buf, sizeof buf) >= 0)
		return strdup(buf);
	if(access(pbuf, AEXIST) >= 0)
		return strdup(pbuf);
	return nil;
}


static int
ptracerw(int type, int xtype, int isr, int pid, uvlong addr, void *v, uint n)
{
	int i;
	uintptr u;
	uchar buf[sizeof(uintptr)];

	for(i=0; i<n; i+=sizeof(uintptr)){
		if(isr){
			errno = 0;
			u = ptrace(type, pid, addr+i, 0);
			if(errno)
				goto ptraceerr;
			if(n-i >= sizeof(uintptr))
				*(uintptr*)((char*)v+i) = u;
			else{
				*(uintptr*)buf = u;
				memmove((char*)v+i, buf, n-i);
			}
		}else{
			if(n-i >= sizeof(uintptr))
				u = *(uintptr*)((char*)v+i);
			else{
				errno = 0;
				u = ptrace(xtype, pid, addr+i, 0);
				if(errno)
					return -1;
				*(uintptr*)buf = u;
				memmove(buf, (char*)v+i, n-i);
				u = *(uintptr*)buf;
			}
			if(ptrace(type, pid, addr+i, u) < 0)
				goto ptraceerr;
		}
	}
	return 0;

ptraceerr:
	werrstr("ptrace %s addr=%#llux pid=%d: %r", isr ? "read" : "write", addr, pid);
	return -1;
}

static int
ptracesegrw(Map *map, Seg *seg, uvlong addr, void *v, uint n, int isr)
{
	return ptracerw(isr ? PTRACE_PEEKDATA : PTRACE_POKEDATA, PTRACE_PEEKDATA,
		isr, map->pid, addr, v, n);
}

static int
go2linux(uvlong addr)
{
	switch(addr){
	case offsetof(Ureg, ax):
		return offsetof(struct user_regs_struct, rax);
	case offsetof(Ureg, bx):
		return offsetof(struct user_regs_struct, rbx);
	case offsetof(Ureg, cx):
		return offsetof(struct user_regs_struct, rcx);
	case offsetof(Ureg, dx):
		return offsetof(struct user_regs_struct, rdx);
	case offsetof(Ureg, si):
		return offsetof(struct user_regs_struct, rsi);
	case offsetof(Ureg, di):
		return offsetof(struct user_regs_struct, rdi);
	case offsetof(Ureg, bp):
		return offsetof(struct user_regs_struct, rbp);
	case offsetof(Ureg, r8):
		return offsetof(struct user_regs_struct, r8);
	case offsetof(Ureg, r9):
		return offsetof(struct user_regs_struct, r9);
	case offsetof(Ureg, r10):
		return offsetof(struct user_regs_struct, r10);
	case offsetof(Ureg, r11):
		return offsetof(struct user_regs_struct, r11);
	case offsetof(Ureg, r12):
		return offsetof(struct user_regs_struct, r12);
	case offsetof(Ureg, r13):
		return offsetof(struct user_regs_struct, r13);
	case offsetof(Ureg, r14):
		return offsetof(struct user_regs_struct, r14);
	case offsetof(Ureg, r15):
		return offsetof(struct user_regs_struct, r15);
	case offsetof(Ureg, ds):
		return offsetof(struct user_regs_struct, ds);
	case offsetof(Ureg, es):
		return offsetof(struct user_regs_struct, es);
	case offsetof(Ureg, fs):
		return offsetof(struct user_regs_struct, fs);
	case offsetof(Ureg, gs):
		return offsetof(struct user_regs_struct, gs);
	case offsetof(Ureg, ip):
		return offsetof(struct user_regs_struct, rip);
	case offsetof(Ureg, cs):
		return offsetof(struct user_regs_struct, cs);
	case offsetof(Ureg, flags):
		return offsetof(struct user_regs_struct, eflags);
	case offsetof(Ureg, sp):
		return offsetof(struct user_regs_struct, rsp);
	case offsetof(Ureg, ss):
		return offsetof(struct user_regs_struct, ss);
	}
	return -1;
}

static int
ptraceregrw(Map *map, Seg *seg, uvlong addr, void *v, uint n, int isr)
{
	int laddr;
	uvlong u;

	if((laddr = go2linux(addr)) < 0){
		if(isr){
			memset(v, 0, n);
			return 0;
		}
		werrstr("register %llud not available", addr);
		return -1;
	}

	if(isr){
		errno = 0;
		u = ptrace(PTRACE_PEEKUSER, map->pid, laddr, 0);
		if(errno)
			goto ptraceerr;
		switch(n){
		case 1:
			*(uint8*)v = u;
			break;
		case 2:
			*(uint16*)v = u;
			break;
		case 4:
			*(uint32*)v = u;
			break;
		case 8:
			*(uint64*)v = u;
			break;
		default:
			werrstr("bad register size");
			return -1;
		}
	}else{
		switch(n){
		case 1:
			u = *(uint8*)v;
			break;
		case 2:
			u = *(uint16*)v;
			break;
		case 4:
			u = *(uint32*)v;
			break;
		case 8:
			u = *(uint64*)v;
			break;
		default:
			werrstr("bad register size");
			return -1;
		}
		if(ptrace(PTRACE_POKEUSER, map->pid, laddr, (void*)(uintptr)u) < 0)
			goto ptraceerr;
	}
	return 0;

ptraceerr:
	werrstr("ptrace %s register laddr=%d pid=%d: %r", isr ? "read" : "write", laddr, map->pid);
	return -1;
}

char*
procstatus(int pid)
{
	LinuxThread *t;

	t = findthread(pid);
	if(t == nil)
		return "???";

	return statestr[t->state];
}
