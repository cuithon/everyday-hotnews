// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// makeFuncStub is jumped to by the code generated by MakeFunc.
// See the comment on the declaration of makeFuncStub in makefunc.go
// for more details.
TEXT ·makeFuncStub(SB),7,$8
	MOVW	R7, 4(R13)
	MOVW	$argframe+0(FP), R1
	MOVW	R1, 8(R13)
	BL	·callReflect(SB)
	RET

// methodValueCall is the code half of the function returned by makeMethodValue.
// See the comment on the declaration of methodValueCall in makefunc.go
// for more details.
TEXT ·methodValueCall(SB),7,$8
	MOVW	R7, 4(R13)
	MOVW	$argframe+0(FP), R1
	MOVW	R1, 8(R13)
	BL	·callMethod(SB)
	RET
