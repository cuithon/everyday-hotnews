// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

TEXT ·Sqrt(SB),NOSPLIT,$0
	Get SP
	F64Load x+0(FP)
	F64Sqrt
	F64Store ret+8(FP)
	RET
