#!/usr/bin/env bash
# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Generate Go code listing errors and other #defined constant
# values (ENAMETOOLONG etc.), by asking the preprocessor
# about the definitions.

unset LANG
export LC_ALL=C
export LC_CTYPE=C

case "$GOARCH" in
arm)
	GCC=arm-gcc
	;;
*)
	GCC=gcc
	;;
esac

uname=$(uname)

includes_Linux='
#define _LARGEFILE_SOURCE
#define _LARGEFILE64_SOURCE
#define _FILE_OFFSET_BITS 64
#define _GNU_SOURCE

#include <sys/types.h>
#include <sys/epoll.h>
#include <linux/ptrace.h>
#include <linux/wait.h>
'

includes_Darwin='
#define __DARWIN_UNIX03 0
#define KERNEL
#define _DARWIN_USE_64_BIT_INODE
#include <sys/wait.h>
#include <sys/event.h>
'

includes_FreeBSD='
#include <sys/wait.h>
#include <sys/event.h>
'

includes='
#include <sys/types.h>
#include <fcntl.h>
#include <dirent.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <netinet/ip.h>
#include <netinet/ip6.h>
#include <netinet/tcp.h>
#include <errno.h>
#include <sys/signal.h>
#include <signal.h>
'

ccflags=""
next=false
for i
do
	if $next; then
		ccflags="$ccflags $i"
		next=false
	elif [ "$i" = "-f" ]; then
		next=true
	fi
done

# These are go errors that will be mapped directly to windows errors
goerrors='
ENOENT:ERROR_FILE_NOT_FOUND
ENOTDIR:ERROR_PATH_NOT_FOUND
'

# Pull out just the error names for later.
i=$(
	for j in "$goerrors"
	do
		echo "$j"
	done |
	awk -F: '
		{ if (NR > 1) printf("|") }
		{ printf("%s", $1) }
	'
)
errors=$(
	echo '#include <errno.h>' | $GCC -x c - -E -dM $ccflags |
	awk '
		$1 != "#define" || $2 ~ /\(/ {next}
		$2 ~ /^('$i')$/ {next}
		$2 ~ /^E[A-Z0-9_]+$/ { print $2 }
		{next}
	' | sort
)

echo '// mkerrors_windows.sh' "$@"
echo '// MACHINE GENERATED BY THE COMMAND ABOVE; DO NOT EDIT'
echo
echo 'package syscall'

# Run C program to print error strings.
(
	/bin/echo "
#include <stdio.h>
#include <errno.h>
#include <ctype.h>
#include <string.h>

#define nelem(x) (sizeof(x)/sizeof((x)[0]))

enum { A = 'A', Z = 'Z', a = 'a', z = 'z' }; // avoid need for single quotes below

struct {
	char *goname;
	char *winname;
} goerrors[] = {
"
	for i in $goerrors
	do
		j=`echo $i | cut -d: -f1`
		k=`echo $i | cut -d: -f2`
		echo '	{"'$j'", "'$k'"},'
	done

	# Use /bin/echo to avoid builtin echo,
	# which interprets \n itself
	/bin/echo '
};

struct {
	char *name;
	int value;
} errors[] = {
'
	for i in $errors
	do
		echo '	{"'$i'",' $i'},'
	done

	# Use /bin/echo to avoid builtin echo,
	# which interprets \n itself
	/bin/echo '
};

int
main(void)
{
	int i, e, iota = 1;
	char buf[1024];

	printf("\n// Go names for Windows errors.\n");
	printf("const (\n");
	for(i=0; i<nelem(goerrors); i++) {
		printf("\t%s = %s\n", goerrors[i].goname, goerrors[i].winname);
			
	}
	printf(")\n");

	printf("\n// Windows reserves errors >= 1<<29 for application use.\n");
	printf("const APPLICATION_ERROR = 1 << 29\n");

	printf("\n// Invented values to support what package os and others expects.\n");
	printf("const (\n");
	for(i=0; i<nelem(errors); i++) {
		printf("\t%s", errors[i].name);
		if(iota) {
			printf(" = APPLICATION_ERROR + iota");
			iota = !iota;
		}
		printf("\n");
			
	}
	printf("\tEWINDOWS\n");
	printf(")\n");

	printf("\n// Error strings for invented errors\n");
	printf("var errors = [...]string {\n");
	for(i=0; i<nelem(errors); i++) {
		e = errors[i].value;
		strcpy(buf, strerror(e));
		// lowercase first letter: Bad -> bad, but STREAM -> STREAM.
		if(A <= buf[0] && buf[0] <= Z && a <= buf[1] && buf[1] <= z)
			buf[0] += a - A;
		printf("\t%s - APPLICATION_ERROR: \"%s\",\n", errors[i].name, buf);
	next:;
	}
	printf("\tEWINDOWS - APPLICATION_ERROR: \"not supported by windows\",\n");
	printf("}\n\n");
	return 0;
}

'
) >_errors.c

$GCC $ccflags -static -o _errors _errors.c && $GORUN ./_errors && rm -f _errors.c _errors
