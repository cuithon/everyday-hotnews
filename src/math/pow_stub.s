// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build 386 amd64 amd64p32 arm

#include "textflag.h"

TEXT ·Pow(SB),NOSPLIT,$0
    JMP ·pow(SB)

