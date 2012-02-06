#!/bin/bash
# Copyright 2011 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

set -e

source utils.bash

if ! test -f ../../src/env.bash; then
	echo "package.bash must be run from $GOROOT/misc/osx" 1>&2
fi

BUILD=/tmp/go.build.tmp
ROOT=`hg root`
export GOROOT=$BUILD/root/usr/local/go
export GOROOT_FINAL=/usr/local/go

echo "Removing old images"
rm -f *.pkg *.dmg

echo "Preparing temporary directory"
rm -rf $BUILD
mkdir -p $BUILD

echo "Copying go source distribution"
mkdir -p $BUILD/root/usr/local
cp -r $ROOT $GOROOT
cp -r etc $BUILD/root/etc

echo "Building go"
pushd $GOROOT > /dev/null
src/version.bash -save
rm -rf .hg .hgignore .hgtags
cd src
./all.bash | sed "s/^/  /"
popd > /dev/null

echo "Building package"
# $PM came from utils.bahs
$PM -v -r $BUILD/root -o "Go `hg id`.pkg" \
	--scripts scripts \
	--id com.googlecode.go \
	--title Go \
	--version "0.1" \
	--target "10.5"

echo "Removing temporary directory"
rm -rf $BUILD
