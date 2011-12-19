#!/bin/sh
# Copyright 2011 The Go Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# This script generates the various derived files involved in
# building package runtime. 
#
#	autogen.sh rebuilds everything
#	autogen.sh -clean deletes the generated files

GOARCHES="
	386
	amd64
	arm
"

GOOSES="
	darwin
	freebsd
	linux
	netbsd
	openbsd
	plan9
	windows
"

GOOSARCHES="
	darwin_386
	darwin_amd64
	freebsd_386
	freebsd_amd64
	linux_386
	linux_amd64
	linux_arm
	netbsd_386
	netbsd_amd64
	openbsd_386
	openbsd_amd64
	plan9_386
	windows_386
	windows_amd64
"

HELPERS="goc2c mkversion"

rm -f $HELPERS z*

if [ "$1" = "-clean" ]; then
	exit 0
fi

set -e

if [ "$GOROOT" = "" ]; then
	echo "$0"': $GOROOT must be set' >&2
	exit 2
fi

# Use goc2c to translate .goc files into arch-specific .c files.
quietgcc -o goc2c -I "$GOROOT/include" goc2c.c "$GOROOT/lib/lib9.a"
for file in *.goc
do
	for arch in $GOARCHES
	do
		base=$(echo $file | sed 's/\.goc$//')
		GOARCH=$arch ./goc2c $file >z.tmp
		mv -f z.tmp z${base}_$arch.c
	done
done

# Version constants.
quietgcc -o mkversion -I "$GOROOT/include" mkversion.c "$GOROOT/lib/lib9.a"
GOROOT="$GOROOT_FINAL" ./mkversion >z.tmp
mv z.tmp zversion.go

for arch in $GOARCHES
do
	(
		echo '// AUTO-GENERATED by autogen.sh; DO NOT EDIT'
		echo
		echo 'package runtime'
		echo
		echo 'const theGoarch = "'$arch'"'
	) >zgoarch_$arch.go
done

for os in $GOOSES
do
	(
		echo '// AUTO-GENERATED by autogen.sh; DO NOT EDIT'
		echo
		echo 'package runtime'
		echo
		echo 'const theGoos = "'$os'"'
	) >zgoos_$os.go
done

# Definitions of runtime structs, translated from C to Go.
for osarch in $GOOSARCHES
do
	./mkgodefs.sh $osarch proc.c iface.c hashmap.c chan.c >z.tmp
	mv -f z.tmp zruntime_defs_$osarch.go
done

# Struct field offsets, for use by assembly files.
for osarch in $GOOSARCHES
do
	./mkasmh.sh $osarch proc.c defs.h >z.tmp
	mv -f z.tmp zasm_$osarch.h
done

rm -f $HELPERS
