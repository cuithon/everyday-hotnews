// Derived from Inferno utils/6c/reg.c
// https://bitbucket.org/inferno-os/inferno-os/src/default/utils/6c/reg.c
//
//	Copyright © 1994-1999 Lucent Technologies Inc.  All rights reserved.
//	Portions Copyright © 1995-1997 C H Forsyth (forsyth@terzarima.net)
//	Portions Copyright © 1997-1999 Vita Nuova Limited
//	Portions Copyright © 2000-2007 Vita Nuova Holdings Limited (www.vitanuova.com)
//	Portions Copyright © 2004,2006 Bruce Ellis
//	Portions Copyright © 2005-2007 C H Forsyth (forsyth@terzarima.net)
//	Revisions Copyright © 2000-2007 Lucent Technologies Inc. and others
//	Portions Copyright © 2009 The Go Authors. All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package amd64

import "cmd/internal/obj/x86"

// For ProgInfo.
const (
	AX  = 1 << (x86.REG_AX - x86.REG_AX)
	BX  = 1 << (x86.REG_BX - x86.REG_AX)
	CX  = 1 << (x86.REG_CX - x86.REG_AX)
	DX  = 1 << (x86.REG_DX - x86.REG_AX)
	DI  = 1 << (x86.REG_DI - x86.REG_AX)
	SI  = 1 << (x86.REG_SI - x86.REG_AX)
	R15 = 1 << (x86.REG_R15 - x86.REG_AX)
	X0  = 1 << 16
)

func RtoB(r int) uint64 {
	if r < x86.REG_AX || r > x86.REG_R15 {
		return 0
	}
	return 1 << uint(r-x86.REG_AX)
}
