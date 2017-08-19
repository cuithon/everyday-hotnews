// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

var sigtable = [...]sigTabT{
	/* 0 */ {0, "SIGNONE: no trap"},
	/* 1 */ {_SigNotify + _SigKill, "SIGHUP: hangup"},
	/* 2 */ {_SigNotify + _SigKill, "SIGINT: interrupt (rubout)"},
	/* 3 */ {_SigNotify + _SigThrow, "SIGQUIT: quit (ASCII FS)"},
	/* 4 */ {_SigThrow + _SigUnblock, "SIGILL: illegal instruction (not reset when caught)"},
	/* 5 */ {_SigThrow + _SigUnblock, "SIGTRAP: trace trap (not reset when caught)"},
	/* 6 */ {_SigNotify + _SigThrow, "SIGABRT: used by abort, replace SIGIOT in the future"},
	/* 7 */ {_SigThrow, "SIGEMT: EMT instruction"},
	/* 8 */ {_SigPanic + _SigUnblock, "SIGFPE: floating point exception"},
	/* 9 */ {0, "SIGKILL: kill (cannot be caught or ignored)"},
	/* 10 */ {_SigPanic + _SigUnblock, "SIGBUS: bus error"},
	/* 11 */ {_SigPanic + _SigUnblock, "SIGSEGV: segmentation violation"},
	/* 12 */ {_SigThrow, "SIGSYS: bad argument to system call"},
	/* 13 */ {_SigNotify, "SIGPIPE: write on a pipe with no one to read it"},
	/* 14 */ {_SigNotify, "SIGALRM: alarm clock"},
	/* 15 */ {_SigNotify + _SigKill, "SIGTERM: software termination signal from kill"},
	/* 16 */ {_SigNotify, "SIGUSR1: user defined signal 1"},
	/* 17 */ {_SigNotify, "SIGUSR2: user defined signal 2"},
	/* 18 */ {_SigNotify + _SigUnblock + _SigIgn, "SIGCHLD: child status change alias (POSIX)"},
	/* 19 */ {_SigNotify, "SIGPWR: power-fail restart"},
	/* 20 */ {_SigNotify + _SigIgn, "SIGWINCH: window size change"},
	/* 21 */ {_SigNotify + _SigIgn, "SIGURG: urgent socket condition"},
	/* 22 */ {_SigNotify, "SIGPOLL: pollable event occurred"},
	/* 23 */ {0, "SIGSTOP: stop (cannot be caught or ignored)"},
	/* 24 */ {_SigNotify + _SigDefault + _SigIgn, "SIGTSTP: user stop requested from tty"},
	/* 25 */ {_SigNotify + _SigDefault + _SigIgn, "SIGCONT: stopped process has been continued"},
	/* 26 */ {_SigNotify + _SigDefault + _SigIgn, "SIGTTIN: background tty read attempted"},
	/* 27 */ {_SigNotify + _SigDefault + _SigIgn, "SIGTTOU: background tty write attempted"},
	/* 28 */ {_SigNotify, "SIGVTALRM: virtual timer expired"},
	/* 29 */ {_SigNotify + _SigUnblock, "SIGPROF: profiling timer expired"},
	/* 30 */ {_SigNotify, "SIGXCPU: exceeded cpu limit"},
	/* 31 */ {_SigNotify, "SIGXFSZ: exceeded file size limit"},
	/* 32 */ {_SigNotify, "SIGWAITING: reserved signal no longer used by"},
	/* 33 */ {_SigNotify, "SIGLWP: reserved signal no longer used by"},
	/* 34 */ {_SigNotify, "SIGFREEZE: special signal used by CPR"},
	/* 35 */ {_SigNotify, "SIGTHAW: special signal used by CPR"},
	/* 36 */ {_SigSetStack + _SigUnblock, "SIGCANCEL: reserved signal for thread cancellation"}, // Oracle's spelling of cancelation.
	/* 37 */ {_SigNotify, "SIGLOST: resource lost (eg, record-lock lost)"},
	/* 38 */ {_SigNotify, "SIGXRES: resource control exceeded"},
	/* 39 */ {_SigNotify, "SIGJVM1: reserved signal for Java Virtual Machine"},
	/* 40 */ {_SigNotify, "SIGJVM2: reserved signal for Java Virtual Machine"},

	/* TODO(aram): what should be do about these signals? _SigDefault or _SigNotify? is this set static? */
	/* 41 */ {_SigNotify, "real time signal"},
	/* 42 */ {_SigNotify, "real time signal"},
	/* 43 */ {_SigNotify, "real time signal"},
	/* 44 */ {_SigNotify, "real time signal"},
	/* 45 */ {_SigNotify, "real time signal"},
	/* 46 */ {_SigNotify, "real time signal"},
	/* 47 */ {_SigNotify, "real time signal"},
	/* 48 */ {_SigNotify, "real time signal"},
	/* 49 */ {_SigNotify, "real time signal"},
	/* 50 */ {_SigNotify, "real time signal"},
	/* 51 */ {_SigNotify, "real time signal"},
	/* 52 */ {_SigNotify, "real time signal"},
	/* 53 */ {_SigNotify, "real time signal"},
	/* 54 */ {_SigNotify, "real time signal"},
	/* 55 */ {_SigNotify, "real time signal"},
	/* 56 */ {_SigNotify, "real time signal"},
	/* 57 */ {_SigNotify, "real time signal"},
	/* 58 */ {_SigNotify, "real time signal"},
	/* 59 */ {_SigNotify, "real time signal"},
	/* 60 */ {_SigNotify, "real time signal"},
	/* 61 */ {_SigNotify, "real time signal"},
	/* 62 */ {_SigNotify, "real time signal"},
	/* 63 */ {_SigNotify, "real time signal"},
	/* 64 */ {_SigNotify, "real time signal"},
	/* 65 */ {_SigNotify, "real time signal"},
	/* 66 */ {_SigNotify, "real time signal"},
	/* 67 */ {_SigNotify, "real time signal"},
	/* 68 */ {_SigNotify, "real time signal"},
	/* 69 */ {_SigNotify, "real time signal"},
	/* 70 */ {_SigNotify, "real time signal"},
	/* 71 */ {_SigNotify, "real time signal"},
	/* 72 */ {_SigNotify, "real time signal"},
}
