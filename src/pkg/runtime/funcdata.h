// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file defines the IDs for PCDATA and FUNCDATA instructions
// in Go binaries. It is included by both C and assembly, so it must
// be written using #defines. It is included by the runtime package
// as well as the compilers.

#define PCDATA_ArgSize 0

// To be used in assembly.
#define ARGSIZE(n) PCDATA $PCDATA_ArgSize, $n
