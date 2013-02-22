// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build race

TEXT	runtime·racefuncenter(SB),7,$0
	PUSHQ	DX // save function entry context (for closures)
	CALL	runtime·racefuncenter1(SB)
	POPQ	DX
	RET
