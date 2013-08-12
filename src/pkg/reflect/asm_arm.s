// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "../../cmd/ld/textflag.h"

// makeFuncStub is jumped to by the code generated by MakeFunc.
// See the comment on the declaration of makeFuncStub in makefunc.go
// for more details.
// No argsize here, gc generates argsize info at call site.
TEXT ·makeFuncStub(SB),NOSPLIT,$8
	MOVW	R7, 4(R13)
	MOVW	$argframe+0(FP), R1
	MOVW	R1, 8(R13)
	BL	·callReflect(SB)
	RET

// methodValueCall is the code half of the function returned by makeMethodValue.
// See the comment on the declaration of methodValueCall in makefunc.go
// for more details.
// No argsize here, gc generates argsize info at call site.
TEXT ·methodValueCall(SB),NOSPLIT,$8
	MOVW	R7, 4(R13)
	MOVW	$argframe+0(FP), R1
	MOVW	R1, 8(R13)
	BL	·callMethod(SB)
	RET
