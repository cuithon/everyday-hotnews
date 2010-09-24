#!/usr/bin/env bash
# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# The syscall package provides access to the raw system call
# interface of the underlying operating system.  Porting Go to
# a new architecture/operating system combination requires
# some manual effort, though there are tools that automate
# much of the process.  The auto-generated files have names
# beginning with z.
#
# This script runs or (given -n) prints suggested commands to generate z files
# for the current system.  Running those commands is not automatic.
# This script is documentation more than anything else.
#
# * asm_${GOOS}_${GOARCH}.s
#
# This hand-written assembly file implements system call dispatch.
# There are three entry points:
#
# 	func Syscall(trap, a1, a2, a3 uintptr) (r1, r2, err uintptr);
# 	func Syscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, err uintptr);
# 	func RawSyscall(trap, a1, a2, a3 uintptr) (r1, r2, err uintptr);
#
# The first and second are the standard ones; they differ only in
# how many arguments can be passed to the kernel.
# The third is for low-level use by the ForkExec wrapper;
# unlike the first two, it does not call into the scheduler to
# let it know that a system call is running.
#
# * syscall_${GOOS}.go
#
# This hand-written Go file implements system calls that need
# special handling and lists "//sys" comments giving prototypes
# for ones that can be auto-generated.  Mksyscall reads those
# comments to generate the stubs.
#
# * syscall_${GOOS}_${GOARCH}.go
#
# Same as syscall_${GOOS}.go except that it contains code specific
# to ${GOOS} on one particular architecture.
#
# * types_${GOOS}.c
#
# This hand-written C file includes standard C headers and then
# creates typedef or enum names beginning with a dollar sign
# (use of $ in variable names is a gcc extension).  The hardest
# part about preparing this file is figuring out which headers to
# include and which symbols need to be #defined to get the
# actual data structures that pass through to the kernel system calls.
# Some C libraries present alternate versions for binary compatibility
# and translate them on the way in and out of system calls, but
# there is almost always a #define that can get the real ones.
# See types_darwin.c and types_linux.c for examples.
#
# * zerror_${GOOS}_${GOARCH}.go
#
# This machine-generated file defines the system's error numbers,
# error strings, and signal numbers.  The generator is "mkerrors.sh".
# Usually no arguments are needed, but mkerrors.sh will pass its
# arguments on to godefs.
#
# * zsyscall_${GOOS}_${GOARCH}.go
#
# Generated by mksyscall.sh; see syscall_${GOOS}.go above.
#
# * zsysnum_${GOOS}_${GOARCH}.go
#
# Generated by mksysnum_${GOOS}.
#
# * ztypes_${GOOS}_${GOARCH}.go
#
# Generated by godefs; see types_${GOOS}.c above.

GOOSARCH="${GOOS}_${GOARCH}"

# defaults
mksyscall="./mksyscall.sh"
mkerrors="./mkerrors.sh"
run="sh"

case "$1" in
-n)
	run="cat"
	shift
esac

case "$#" in
0)
	;;
*)
	echo 'usage: mkall.sh [-n]' 1>&2
	exit 2
esac

case "$GOOSARCH" in
_* | *_ | _)
	echo 'undefined $GOOS_$GOARCH:' "$GOOSARCH" 1>&2
	exit 1
	;;
freebsd_386)
	mkerrors="$mkerrors -f -m32"
	mksyscall="./mksyscall.sh -l32"
	mksysnum="curl -s 'http://svn.freebsd.org/viewvc/base/head/sys/kern/syscalls.master?view=markup' | ./mksysnum_freebsd.sh"
	mktypes="godefs -gsyscall -f-m32"
	;;
freebsd_amd64)
	mkerrors="$mkerrors -f -m64"
	mksysnum="curl -s 'http://svn.freebsd.org/viewvc/base/head/sys/kern/syscalls.master?view=markup' | ./mksysnum_freebsd.sh"
	mktypes="godefs -gsyscall -f-m64"
	;;
darwin_386)
	mkerrors="$mkerrors -f -m32"
	mksyscall="./mksyscall.sh -l32"
	mksysnum="./mksysnum_darwin.sh /home/rsc/pub/xnu-1228/bsd/kern/syscalls.master"
	mktypes="godefs -gsyscall -f-m32"
	;;
darwin_amd64)
	mkerrors="$mkerrors -f -m64"
	mksysnum="./mksysnum_darwin.sh /home/rsc/pub/xnu-1228/bsd/kern/syscalls.master"
	mktypes="godefs -gsyscall -f-m64"
	mkerrors="./mkerrors.sh"
	;;
linux_386)
	mkerrors="$mkerrors -f -m32"
	mksyscall="./mksyscall.sh -l32"
	mksysnum="./mksysnum_linux.sh /usr/include/asm/unistd_32.h"
	mktypes="godefs -gsyscall -f-m32"
	;;
linux_amd64)
	mkerrors="$mkerrors -f -m64"
	mksysnum="./mksysnum_linux.sh /usr/include/asm/unistd_64.h"
	mktypes="godefs -gsyscall -f-m64"
	;;
nacl_386)
	NACL="/home/rsc/pub/nacl/native_client"
	NACLRUN="$NACL/src/trusted/service_runtime"
	NACLSDK="$NACL/src/third_party/nacl_sdk/linux/sdk/nacl-sdk/nacl"
	mksyscall="./mksyscall.sh -l32 -nacl"
	mksysnum="./mksysnum_nacl.sh $NACLRUN/include/bits/nacl_syscalls.h"
	mktypes="godefs -gsyscall -f-m32 -f-I$NACLSDK/include -f-I$NACL"
	mkerrors="./mkerrors_nacl.sh $NACLRUN/include/sys/errno.h"
	;;
linux_arm)
	mkerrors="$mkerrors"
	mksyscall="./mksyscall.sh -l32"
	mksysnum="./mksysnum_linux.sh /usr/include/asm/unistd.h"
	mktypes="godefs -gsyscall"
	;;
windows_386)
	mksyscall="./mksyscall_windows.sh -l32"
	mksysnum=
	mktypes=
	mkerrors="./mkerrors_windows.sh -f -m32"
	;;
*)
	echo 'unrecognized $GOOS_$GOARCH: ' "$GOOSARCH" 1>&2
	exit 1
	;;
esac

(
	if [ -n "$mkerrors" ]; then echo "$mkerrors |gofmt >zerrors_$GOOSARCH.go"; fi
	syscall_goos="syscall_$GOOS.go"
	case "$GOOS" in
	darwin | freebsd)
		syscall_goos="syscall_bsd.go $syscall_goos"
		;;
	esac
	if [ -n "$mksyscall" ]; then echo "$mksyscall $syscall_goos syscall_$GOOSARCH.go |gofmt >zsyscall_$GOOSARCH.go"; fi
	if [ -n "$mksysnum" ]; then echo "$mksysnum |gofmt >zsysnum_$GOOSARCH.go"; fi
	if [ -n "$mktypes" ]; then echo "$mktypes types_$GOOS.c |gofmt >ztypes_$GOOSARCH.go"; fi
) | $run
