// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

TEXT _rt0_amd64_openbsd(SB),7,$-8
	MOVQ	$_rt0_amd64(SB), DX
	MOVQ	SP, DI
	JMP	DX
