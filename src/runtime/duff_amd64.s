// Code generated by mkduff.go; DO NOT EDIT.
// Run go generate from src/runtime to update.
// See mkduff.go for comments.

#include "textflag.h"

TEXT runtime·duffzero(SB), NOSPLIT, $0-0
	MOVUPS	X0,(DI)
	MOVUPS	X0,16(DI)
	MOVUPS	X0,32(DI)
	MOVUPS	X0,48(DI)
	ADDQ	$64,DI

	MOVUPS	X0,(DI)
	MOVUPS	X0,16(DI)
	MOVUPS	X0,32(DI)
	MOVUPS	X0,48(DI)
	ADDQ	$64,DI

	MOVUPS	X0,(DI)
	MOVUPS	X0,16(DI)
	MOVUPS	X0,32(DI)
	MOVUPS	X0,48(DI)
	ADDQ	$64,DI

	MOVUPS	X0,(DI)
	MOVUPS	X0,16(DI)
	MOVUPS	X0,32(DI)
	MOVUPS	X0,48(DI)
	ADDQ	$64,DI

	MOVUPS	X0,(DI)
	MOVUPS	X0,16(DI)
	MOVUPS	X0,32(DI)
	MOVUPS	X0,48(DI)
	ADDQ	$64,DI

	MOVUPS	X0,(DI)
	MOVUPS	X0,16(DI)
	MOVUPS	X0,32(DI)
	MOVUPS	X0,48(DI)
	ADDQ	$64,DI

	MOVUPS	X0,(DI)
	MOVUPS	X0,16(DI)
	MOVUPS	X0,32(DI)
	MOVUPS	X0,48(DI)
	ADDQ	$64,DI

	MOVUPS	X0,(DI)
	MOVUPS	X0,16(DI)
	MOVUPS	X0,32(DI)
	MOVUPS	X0,48(DI)
	ADDQ	$64,DI

	MOVUPS	X0,(DI)
	MOVUPS	X0,16(DI)
	MOVUPS	X0,32(DI)
	MOVUPS	X0,48(DI)
	ADDQ	$64,DI

	MOVUPS	X0,(DI)
	MOVUPS	X0,16(DI)
	MOVUPS	X0,32(DI)
	MOVUPS	X0,48(DI)
	ADDQ	$64,DI

	MOVUPS	X0,(DI)
	MOVUPS	X0,16(DI)
	MOVUPS	X0,32(DI)
	MOVUPS	X0,48(DI)
	ADDQ	$64,DI

	MOVUPS	X0,(DI)
	MOVUPS	X0,16(DI)
	MOVUPS	X0,32(DI)
	MOVUPS	X0,48(DI)
	ADDQ	$64,DI

	MOVUPS	X0,(DI)
	MOVUPS	X0,16(DI)
	MOVUPS	X0,32(DI)
	MOVUPS	X0,48(DI)
	ADDQ	$64,DI

	MOVUPS	X0,(DI)
	MOVUPS	X0,16(DI)
	MOVUPS	X0,32(DI)
	MOVUPS	X0,48(DI)
	ADDQ	$64,DI

	MOVUPS	X0,(DI)
	MOVUPS	X0,16(DI)
	MOVUPS	X0,32(DI)
	MOVUPS	X0,48(DI)
	ADDQ	$64,DI

	MOVUPS	X0,(DI)
	MOVUPS	X0,16(DI)
	MOVUPS	X0,32(DI)
	MOVUPS	X0,48(DI)
	ADDQ	$64,DI

	RET

TEXT runtime·duffcopy(SB), NOSPLIT, $0-0
	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	MOVUPS	(SI), X0
	ADDQ	$16, SI
	MOVUPS	X0, (DI)
	ADDQ	$16, DI

	RET
