// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * Input to godefs.
 *
	godefs -f -m64 defs.c >amd64/defs.h
	godefs -f -m32 defs.c >386/defs.h
 */

#include <sys/types.h>
#include <sys/mman.h>
#include <sys/time.h>
#include <sys/unistd.h>
#include <sys/signal.h>
#include <errno.h>
#include <signal.h>

enum {
	$PROT_NONE = PROT_NONE,
	$PROT_READ = PROT_READ,
	$PROT_WRITE = PROT_WRITE,
	$PROT_EXEC = PROT_EXEC,

	$MAP_ANON = MAP_ANON,
	$MAP_PRIVATE = MAP_PRIVATE,
	$MAP_FIXED = MAP_FIXED,

	$SA_SIGINFO = SA_SIGINFO,
	$SA_RESTART = SA_RESTART,
	$SA_ONSTACK = SA_ONSTACK,

	$EINTR = EINTR,
	
	$SIGHUP = SIGHUP,
	$SIGINT = SIGINT,
	$SIGQUIT = SIGQUIT,
	$SIGILL = SIGILL,
	$SIGTRAP = SIGTRAP,
	$SIGABRT = SIGABRT,
	$SIGEMT = SIGEMT,
	$SIGFPE = SIGFPE,
	$SIGKILL = SIGKILL,
	$SIGBUS = SIGBUS,
	$SIGSEGV = SIGSEGV,
	$SIGSYS = SIGSYS,
	$SIGPIPE = SIGPIPE,
	$SIGALRM = SIGALRM,
	$SIGTERM = SIGTERM,
	$SIGURG = SIGURG,
	$SIGSTOP = SIGSTOP,
	$SIGTSTP = SIGTSTP,
	$SIGCONT = SIGCONT,
	$SIGCHLD = SIGCHLD,
	$SIGTTIN = SIGTTIN,
	$SIGTTOU = SIGTTOU,
	$SIGIO = SIGIO,
	$SIGXCPU = SIGXCPU,
	$SIGXFSZ = SIGXFSZ,
	$SIGVTALRM = SIGVTALRM,
	$SIGPROF = SIGPROF,
	$SIGWINCH = SIGWINCH,
	$SIGINFO = SIGINFO,
	$SIGUSR1 = SIGUSR1,
	$SIGUSR2 = SIGUSR2,
	
	$FPE_INTDIV = FPE_INTDIV,
	$FPE_INTOVF = FPE_INTOVF,
	$FPE_FLTDIV = FPE_FLTDIV,
	$FPE_FLTOVF = FPE_FLTOVF,
	$FPE_FLTUND = FPE_FLTUND,
	$FPE_FLTRES = FPE_FLTRES,
	$FPE_FLTINV = FPE_FLTINV,
	$FPE_FLTSUB = FPE_FLTSUB,
	
	$BUS_ADRALN = BUS_ADRALN,
	$BUS_ADRERR = BUS_ADRERR,
	$BUS_OBJERR = BUS_OBJERR,
	
	$SEGV_MAPERR = SEGV_MAPERR,
	$SEGV_ACCERR = SEGV_ACCERR,
	
	$ITIMER_REAL = ITIMER_REAL,
	$ITIMER_VIRTUAL = ITIMER_VIRTUAL,
	$ITIMER_PROF = ITIMER_PROF,
};

typedef struct sigaltstack $Sigaltstack;
typedef sigset_t $Sigset;
typedef siginfo_t $Siginfo;
typedef union sigval $Sigval;

typedef stack_t $StackT;

typedef struct timeval $Timeval;
typedef struct itimerval $Itimerval;

// This is a hack to avoid pulling in machine/fpu.h.
typedef void $sfxsave64;
typedef void $usavefpu;

typedef struct sigcontext $Sigcontext;
