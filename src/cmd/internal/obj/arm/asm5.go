// Inferno utils/5l/span.c
// https://bitbucket.org/inferno-os/inferno-os/src/default/utils/5l/span.c
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

package arm

import (
	"cmd/internal/obj"
	"cmd/internal/objabi"
	"fmt"
	"log"
	"math"
	"sort"
)

// ctxt5 holds state while assembling a single function.
// Each function gets a fresh ctxt5.
// This allows for multiple functions to be safely concurrently assembled.
type ctxt5 struct {
	ctxt       *obj.Link
	newprog    obj.ProgAlloc
	cursym     *obj.LSym
	printp     *obj.Prog
	blitrl     *obj.Prog
	elitrl     *obj.Prog
	autosize   int64
	instoffset int64
	pc         int64
	pool       struct {
		start uint32
		size  uint32
		extra uint32
	}
}

type Optab struct {
	as       obj.As
	a1       uint8
	a2       int8
	a3       uint8
	type_    uint8
	size     int8
	param    int16
	flag     int8
	pcrelsiz uint8
}

type Opcross [32][2][32]uint8

const (
	LFROM  = 1 << 0
	LTO    = 1 << 1
	LPOOL  = 1 << 2
	LPCREL = 1 << 3
)

var optab = []Optab{
	/* struct Optab:
	OPCODE,	from, prog->reg, to,		 type,size,param,flag */
	{obj.ATEXT, C_ADDR, C_NONE, C_TEXTSIZE, 0, 0, 0, 0, 0},
	{AADD, C_REG, C_REG, C_REG, 1, 4, 0, 0, 0},
	{AADD, C_REG, C_NONE, C_REG, 1, 4, 0, 0, 0},
	{AMOVW, C_REG, C_NONE, C_REG, 1, 4, 0, 0, 0},
	{AMVN, C_REG, C_NONE, C_REG, 1, 4, 0, 0, 0},
	{ACMP, C_REG, C_REG, C_NONE, 1, 4, 0, 0, 0},
	{AADD, C_RCON, C_REG, C_REG, 2, 4, 0, 0, 0},
	{AADD, C_RCON, C_NONE, C_REG, 2, 4, 0, 0, 0},
	{AMOVW, C_RCON, C_NONE, C_REG, 2, 4, 0, 0, 0},
	{AMVN, C_RCON, C_NONE, C_REG, 2, 4, 0, 0, 0},
	{ACMP, C_RCON, C_REG, C_NONE, 2, 4, 0, 0, 0},
	{AADD, C_SHIFT, C_REG, C_REG, 3, 4, 0, 0, 0},
	{AADD, C_SHIFT, C_NONE, C_REG, 3, 4, 0, 0, 0},
	{AMVN, C_SHIFT, C_NONE, C_REG, 3, 4, 0, 0, 0},
	{ACMP, C_SHIFT, C_REG, C_NONE, 3, 4, 0, 0, 0},
	{AMOVW, C_RACON, C_NONE, C_REG, 4, 4, REGSP, 0, 0},
	{AB, C_NONE, C_NONE, C_SBRA, 5, 4, 0, LPOOL, 0},
	{ABL, C_NONE, C_NONE, C_SBRA, 5, 4, 0, 0, 0},
	{ABX, C_NONE, C_NONE, C_SBRA, 74, 20, 0, 0, 0},
	{ABEQ, C_NONE, C_NONE, C_SBRA, 5, 4, 0, 0, 0},
	{ABEQ, C_RCON, C_NONE, C_SBRA, 5, 4, 0, 0, 0}, // prediction hinted form, hint ignored

	{AB, C_NONE, C_NONE, C_ROREG, 6, 4, 0, LPOOL, 0},
	{ABL, C_NONE, C_NONE, C_ROREG, 7, 4, 0, 0, 0},
	{ABL, C_REG, C_NONE, C_ROREG, 7, 4, 0, 0, 0},
	{ABX, C_NONE, C_NONE, C_ROREG, 75, 12, 0, 0, 0},
	{ABXRET, C_NONE, C_NONE, C_ROREG, 76, 4, 0, 0, 0},
	{ASLL, C_RCON, C_REG, C_REG, 8, 4, 0, 0, 0},
	{ASLL, C_RCON, C_NONE, C_REG, 8, 4, 0, 0, 0},
	{ASLL, C_REG, C_NONE, C_REG, 9, 4, 0, 0, 0},
	{ASLL, C_REG, C_REG, C_REG, 9, 4, 0, 0, 0},
	{ASWI, C_NONE, C_NONE, C_NONE, 10, 4, 0, 0, 0},
	{ASWI, C_NONE, C_NONE, C_LOREG, 10, 4, 0, 0, 0},
	{ASWI, C_NONE, C_NONE, C_LCON, 10, 4, 0, 0, 0},
	{AWORD, C_NONE, C_NONE, C_LCON, 11, 4, 0, 0, 0},
	{AWORD, C_NONE, C_NONE, C_LCONADDR, 11, 4, 0, 0, 0},
	{AWORD, C_NONE, C_NONE, C_ADDR, 11, 4, 0, 0, 0},
	{AWORD, C_NONE, C_NONE, C_TLS_LE, 103, 4, 0, 0, 0},
	{AWORD, C_NONE, C_NONE, C_TLS_IE, 104, 4, 0, 0, 0},
	{AMOVW, C_NCON, C_NONE, C_REG, 12, 4, 0, 0, 0},
	{AMOVW, C_SCON, C_NONE, C_REG, 12, 4, 0, 0, 0},
	{AMOVW, C_LCON, C_NONE, C_REG, 12, 4, 0, LFROM, 0},
	{AMOVW, C_LCONADDR, C_NONE, C_REG, 12, 4, 0, LFROM | LPCREL, 4},
	{AADD, C_NCON, C_REG, C_REG, 13, 8, 0, 0, 0},
	{AADD, C_NCON, C_NONE, C_REG, 13, 8, 0, 0, 0},
	{AMVN, C_NCON, C_NONE, C_REG, 13, 8, 0, 0, 0},
	{ACMP, C_NCON, C_REG, C_NONE, 13, 8, 0, 0, 0},
	{AADD, C_SCON, C_REG, C_REG, 13, 8, 0, 0, 0},
	{AADD, C_SCON, C_NONE, C_REG, 13, 8, 0, 0, 0},
	{AMVN, C_SCON, C_NONE, C_REG, 13, 8, 0, 0, 0},
	{ACMP, C_SCON, C_REG, C_NONE, 13, 8, 0, 0, 0},
	{AADD, C_LCON, C_REG, C_REG, 13, 8, 0, LFROM, 0},
	{AADD, C_LCON, C_NONE, C_REG, 13, 8, 0, LFROM, 0},
	{AMVN, C_LCON, C_NONE, C_REG, 13, 8, 0, LFROM, 0},
	{ACMP, C_LCON, C_REG, C_NONE, 13, 8, 0, LFROM, 0},
	{AMOVB, C_REG, C_NONE, C_REG, 1, 4, 0, 0, 0},
	{AMOVBS, C_REG, C_NONE, C_REG, 14, 8, 0, 0, 0},
	{AMOVBU, C_REG, C_NONE, C_REG, 58, 4, 0, 0, 0},
	{AMOVH, C_REG, C_NONE, C_REG, 1, 4, 0, 0, 0},
	{AMOVHS, C_REG, C_NONE, C_REG, 14, 8, 0, 0, 0},
	{AMOVHU, C_REG, C_NONE, C_REG, 14, 8, 0, 0, 0},
	{AMUL, C_REG, C_REG, C_REG, 15, 4, 0, 0, 0},
	{AMUL, C_REG, C_NONE, C_REG, 15, 4, 0, 0, 0},
	{ADIV, C_REG, C_REG, C_REG, 16, 4, 0, 0, 0},
	{ADIV, C_REG, C_NONE, C_REG, 16, 4, 0, 0, 0},
	{ADIVHW, C_REG, C_REG, C_REG, 105, 4, 0, 0, 0},
	{ADIVHW, C_REG, C_NONE, C_REG, 105, 4, 0, 0, 0},
	{AMULL, C_REG, C_REG, C_REGREG, 17, 4, 0, 0, 0},
	{AMULA, C_REG, C_REG, C_REGREG2, 17, 4, 0, 0, 0},
	{AMOVW, C_REG, C_NONE, C_SAUTO, 20, 4, REGSP, 0, 0},
	{AMOVW, C_REG, C_NONE, C_SOREG, 20, 4, 0, 0, 0},
	{AMOVB, C_REG, C_NONE, C_SAUTO, 20, 4, REGSP, 0, 0},
	{AMOVB, C_REG, C_NONE, C_SOREG, 20, 4, 0, 0, 0},
	{AMOVBS, C_REG, C_NONE, C_SAUTO, 20, 4, REGSP, 0, 0},
	{AMOVBS, C_REG, C_NONE, C_SOREG, 20, 4, 0, 0, 0},
	{AMOVBU, C_REG, C_NONE, C_SAUTO, 20, 4, REGSP, 0, 0},
	{AMOVBU, C_REG, C_NONE, C_SOREG, 20, 4, 0, 0, 0},
	{AMOVW, C_SAUTO, C_NONE, C_REG, 21, 4, REGSP, 0, 0},
	{AMOVW, C_SOREG, C_NONE, C_REG, 21, 4, 0, 0, 0},
	{AMOVBU, C_SAUTO, C_NONE, C_REG, 21, 4, REGSP, 0, 0},
	{AMOVBU, C_SOREG, C_NONE, C_REG, 21, 4, 0, 0, 0},
	{AMOVW, C_REG, C_NONE, C_LAUTO, 30, 8, REGSP, LTO, 0},
	{AMOVW, C_REG, C_NONE, C_LOREG, 30, 8, 0, LTO, 0},
	{AMOVW, C_REG, C_NONE, C_ADDR, 64, 8, 0, LTO | LPCREL, 4},
	{AMOVB, C_REG, C_NONE, C_LAUTO, 30, 8, REGSP, LTO, 0},
	{AMOVB, C_REG, C_NONE, C_LOREG, 30, 8, 0, LTO, 0},
	{AMOVB, C_REG, C_NONE, C_ADDR, 64, 8, 0, LTO | LPCREL, 4},
	{AMOVBS, C_REG, C_NONE, C_LAUTO, 30, 8, REGSP, LTO, 0},
	{AMOVBS, C_REG, C_NONE, C_LOREG, 30, 8, 0, LTO, 0},
	{AMOVBS, C_REG, C_NONE, C_ADDR, 64, 8, 0, LTO | LPCREL, 4},
	{AMOVBU, C_REG, C_NONE, C_LAUTO, 30, 8, REGSP, LTO, 0},
	{AMOVBU, C_REG, C_NONE, C_LOREG, 30, 8, 0, LTO, 0},
	{AMOVBU, C_REG, C_NONE, C_ADDR, 64, 8, 0, LTO | LPCREL, 4},
	{AMOVW, C_TLS_LE, C_NONE, C_REG, 101, 4, 0, LFROM, 0},
	{AMOVW, C_TLS_IE, C_NONE, C_REG, 102, 8, 0, LFROM, 0},
	{AMOVW, C_LAUTO, C_NONE, C_REG, 31, 8, REGSP, LFROM, 0},
	{AMOVW, C_LOREG, C_NONE, C_REG, 31, 8, 0, LFROM, 0},
	{AMOVW, C_ADDR, C_NONE, C_REG, 65, 8, 0, LFROM | LPCREL, 4},
	{AMOVBU, C_LAUTO, C_NONE, C_REG, 31, 8, REGSP, LFROM, 0},
	{AMOVBU, C_LOREG, C_NONE, C_REG, 31, 8, 0, LFROM, 0},
	{AMOVBU, C_ADDR, C_NONE, C_REG, 65, 8, 0, LFROM | LPCREL, 4},
	{AMOVW, C_LACON, C_NONE, C_REG, 34, 8, REGSP, LFROM, 0},
	{AMOVW, C_PSR, C_NONE, C_REG, 35, 4, 0, 0, 0},
	{AMOVW, C_REG, C_NONE, C_PSR, 36, 4, 0, 0, 0},
	{AMOVW, C_RCON, C_NONE, C_PSR, 37, 4, 0, 0, 0},
	{AMOVM, C_REGLIST, C_NONE, C_SOREG, 38, 4, 0, 0, 0},
	{AMOVM, C_SOREG, C_NONE, C_REGLIST, 39, 4, 0, 0, 0},
	{ASWPW, C_SOREG, C_REG, C_REG, 40, 4, 0, 0, 0},
	{ARFE, C_NONE, C_NONE, C_NONE, 41, 4, 0, 0, 0},
	{AMOVF, C_FREG, C_NONE, C_FAUTO, 50, 4, REGSP, 0, 0},
	{AMOVF, C_FREG, C_NONE, C_FOREG, 50, 4, 0, 0, 0},
	{AMOVF, C_FAUTO, C_NONE, C_FREG, 51, 4, REGSP, 0, 0},
	{AMOVF, C_FOREG, C_NONE, C_FREG, 51, 4, 0, 0, 0},
	{AMOVF, C_FREG, C_NONE, C_LAUTO, 52, 12, REGSP, LTO, 0},
	{AMOVF, C_FREG, C_NONE, C_LOREG, 52, 12, 0, LTO, 0},
	{AMOVF, C_LAUTO, C_NONE, C_FREG, 53, 12, REGSP, LFROM, 0},
	{AMOVF, C_LOREG, C_NONE, C_FREG, 53, 12, 0, LFROM, 0},
	{AMOVF, C_FREG, C_NONE, C_ADDR, 68, 8, 0, LTO | LPCREL, 4},
	{AMOVF, C_ADDR, C_NONE, C_FREG, 69, 8, 0, LFROM | LPCREL, 4},
	{AADDF, C_FREG, C_NONE, C_FREG, 54, 4, 0, 0, 0},
	{AADDF, C_FREG, C_REG, C_FREG, 54, 4, 0, 0, 0},
	{AMOVF, C_FREG, C_NONE, C_FREG, 54, 4, 0, 0, 0},
	{AMOVW, C_REG, C_NONE, C_FCR, 56, 4, 0, 0, 0},
	{AMOVW, C_FCR, C_NONE, C_REG, 57, 4, 0, 0, 0},
	{AMOVW, C_SHIFT, C_NONE, C_REG, 59, 4, 0, 0, 0},
	{AMOVBU, C_SHIFT, C_NONE, C_REG, 59, 4, 0, 0, 0},
	{AMOVB, C_SHIFT, C_NONE, C_REG, 60, 4, 0, 0, 0},
	{AMOVBS, C_SHIFT, C_NONE, C_REG, 60, 4, 0, 0, 0},
	{AMOVW, C_REG, C_NONE, C_SHIFT, 61, 4, 0, 0, 0},
	{AMOVB, C_REG, C_NONE, C_SHIFT, 61, 4, 0, 0, 0},
	{AMOVBS, C_REG, C_NONE, C_SHIFT, 61, 4, 0, 0, 0},
	{AMOVBU, C_REG, C_NONE, C_SHIFT, 61, 4, 0, 0, 0},
	{AMOVH, C_REG, C_NONE, C_HAUTO, 70, 4, REGSP, 0, 0},
	{AMOVH, C_REG, C_NONE, C_HOREG, 70, 4, 0, 0, 0},
	{AMOVHS, C_REG, C_NONE, C_HAUTO, 70, 4, REGSP, 0, 0},
	{AMOVHS, C_REG, C_NONE, C_HOREG, 70, 4, 0, 0, 0},
	{AMOVHU, C_REG, C_NONE, C_HAUTO, 70, 4, REGSP, 0, 0},
	{AMOVHU, C_REG, C_NONE, C_HOREG, 70, 4, 0, 0, 0},
	{AMOVB, C_HAUTO, C_NONE, C_REG, 71, 4, REGSP, 0, 0},
	{AMOVB, C_HOREG, C_NONE, C_REG, 71, 4, 0, 0, 0},
	{AMOVBS, C_HAUTO, C_NONE, C_REG, 71, 4, REGSP, 0, 0},
	{AMOVBS, C_HOREG, C_NONE, C_REG, 71, 4, 0, 0, 0},
	{AMOVH, C_HAUTO, C_NONE, C_REG, 71, 4, REGSP, 0, 0},
	{AMOVH, C_HOREG, C_NONE, C_REG, 71, 4, 0, 0, 0},
	{AMOVHS, C_HAUTO, C_NONE, C_REG, 71, 4, REGSP, 0, 0},
	{AMOVHS, C_HOREG, C_NONE, C_REG, 71, 4, 0, 0, 0},
	{AMOVHU, C_HAUTO, C_NONE, C_REG, 71, 4, REGSP, 0, 0},
	{AMOVHU, C_HOREG, C_NONE, C_REG, 71, 4, 0, 0, 0},
	{AMOVH, C_REG, C_NONE, C_LAUTO, 72, 8, REGSP, LTO, 0},
	{AMOVH, C_REG, C_NONE, C_LOREG, 72, 8, 0, LTO, 0},
	{AMOVH, C_REG, C_NONE, C_ADDR, 94, 8, 0, LTO | LPCREL, 4},
	{AMOVHS, C_REG, C_NONE, C_LAUTO, 72, 8, REGSP, LTO, 0},
	{AMOVHS, C_REG, C_NONE, C_LOREG, 72, 8, 0, LTO, 0},
	{AMOVHS, C_REG, C_NONE, C_ADDR, 94, 8, 0, LTO | LPCREL, 4},
	{AMOVHU, C_REG, C_NONE, C_LAUTO, 72, 8, REGSP, LTO, 0},
	{AMOVHU, C_REG, C_NONE, C_LOREG, 72, 8, 0, LTO, 0},
	{AMOVHU, C_REG, C_NONE, C_ADDR, 94, 8, 0, LTO | LPCREL, 4},
	{AMOVB, C_LAUTO, C_NONE, C_REG, 73, 8, REGSP, LFROM, 0},
	{AMOVB, C_LOREG, C_NONE, C_REG, 73, 8, 0, LFROM, 0},
	{AMOVB, C_ADDR, C_NONE, C_REG, 93, 8, 0, LFROM | LPCREL, 4},
	{AMOVBS, C_LAUTO, C_NONE, C_REG, 73, 8, REGSP, LFROM, 0},
	{AMOVBS, C_LOREG, C_NONE, C_REG, 73, 8, 0, LFROM, 0},
	{AMOVBS, C_ADDR, C_NONE, C_REG, 93, 8, 0, LFROM | LPCREL, 4},
	{AMOVH, C_LAUTO, C_NONE, C_REG, 73, 8, REGSP, LFROM, 0},
	{AMOVH, C_LOREG, C_NONE, C_REG, 73, 8, 0, LFROM, 0},
	{AMOVH, C_ADDR, C_NONE, C_REG, 93, 8, 0, LFROM | LPCREL, 4},
	{AMOVHS, C_LAUTO, C_NONE, C_REG, 73, 8, REGSP, LFROM, 0},
	{AMOVHS, C_LOREG, C_NONE, C_REG, 73, 8, 0, LFROM, 0},
	{AMOVHS, C_ADDR, C_NONE, C_REG, 93, 8, 0, LFROM | LPCREL, 4},
	{AMOVHU, C_LAUTO, C_NONE, C_REG, 73, 8, REGSP, LFROM, 0},
	{AMOVHU, C_LOREG, C_NONE, C_REG, 73, 8, 0, LFROM, 0},
	{AMOVHU, C_ADDR, C_NONE, C_REG, 93, 8, 0, LFROM | LPCREL, 4},
	{ALDREX, C_SOREG, C_NONE, C_REG, 77, 4, 0, 0, 0},
	{ASTREX, C_SOREG, C_REG, C_REG, 78, 4, 0, 0, 0},
	{AMOVF, C_ZFCON, C_NONE, C_FREG, 80, 8, 0, 0, 0},
	{AMOVF, C_SFCON, C_NONE, C_FREG, 81, 4, 0, 0, 0},
	{ACMPF, C_FREG, C_REG, C_NONE, 82, 8, 0, 0, 0},
	{ACMPF, C_FREG, C_NONE, C_NONE, 83, 8, 0, 0, 0},
	{AMOVFW, C_FREG, C_NONE, C_FREG, 84, 4, 0, 0, 0},
	{AMOVWF, C_FREG, C_NONE, C_FREG, 85, 4, 0, 0, 0},
	{AMOVFW, C_FREG, C_NONE, C_REG, 86, 8, 0, 0, 0},
	{AMOVWF, C_REG, C_NONE, C_FREG, 87, 8, 0, 0, 0},
	{AMOVW, C_REG, C_NONE, C_FREG, 88, 4, 0, 0, 0},
	{AMOVW, C_FREG, C_NONE, C_REG, 89, 4, 0, 0, 0},
	{ALDREXD, C_SOREG, C_NONE, C_REG, 91, 4, 0, 0, 0},
	{ASTREXD, C_SOREG, C_REG, C_REG, 92, 4, 0, 0, 0},
	{APLD, C_SOREG, C_NONE, C_NONE, 95, 4, 0, 0, 0},
	{obj.AUNDEF, C_NONE, C_NONE, C_NONE, 96, 4, 0, 0, 0},
	{ACLZ, C_REG, C_NONE, C_REG, 97, 4, 0, 0, 0},
	{AMULWT, C_REG, C_REG, C_REG, 98, 4, 0, 0, 0},
	{AMULAWT, C_REG, C_REG, C_REGREG2, 99, 4, 0, 0, 0},
	{obj.APCDATA, C_LCON, C_NONE, C_LCON, 0, 0, 0, 0, 0},
	{obj.AFUNCDATA, C_LCON, C_NONE, C_ADDR, 0, 0, 0, 0, 0},
	{obj.ANOP, C_NONE, C_NONE, C_NONE, 0, 0, 0, 0, 0},
	{obj.ADUFFZERO, C_NONE, C_NONE, C_SBRA, 5, 4, 0, 0, 0}, // same as ABL
	{obj.ADUFFCOPY, C_NONE, C_NONE, C_SBRA, 5, 4, 0, 0, 0}, // same as ABL

	{ADATABUNDLE, C_NONE, C_NONE, C_NONE, 100, 4, 0, 0, 0},
	{ADATABUNDLEEND, C_NONE, C_NONE, C_NONE, 100, 0, 0, 0, 0},
	{obj.AXXX, C_NONE, C_NONE, C_NONE, 0, 4, 0, 0, 0},
}

var oprange [ALAST & obj.AMask][]Optab

var xcmp [C_GOK + 1][C_GOK + 1]bool

var (
	deferreturn *obj.LSym
	symdiv      *obj.LSym
	symdivu     *obj.LSym
	symmod      *obj.LSym
	symmodu     *obj.LSym
)

// Note about encoding: Prog.scond holds the condition encoding,
// but XOR'ed with C_SCOND_XOR, so that C_SCOND_NONE == 0.
// The code that shifts the value << 28 has the responsibility
// for XORing with C_SCOND_XOR too.

// asmoutnacl assembles the instruction p. It replaces asmout for NaCl.
// It returns the total number of bytes put in out, and it can change
// p->pc if extra padding is necessary.
// In rare cases, asmoutnacl might split p into two instructions.
// origPC is the PC for this Prog (no padding is taken into account).
func (c *ctxt5) asmoutnacl(origPC int32, p *obj.Prog, o *Optab, out []uint32) int {
	size := int(o.size)

	// instruction specific
	switch p.As {
	default:
		if out != nil {
			c.asmout(p, o, out)
		}

	case ADATABUNDLE, // align to 16-byte boundary
		ADATABUNDLEEND: // zero width instruction, just to align next instruction to 16-byte boundary
		p.Pc = (p.Pc + 15) &^ 15

		if out != nil {
			c.asmout(p, o, out)
		}

	case obj.AUNDEF,
		APLD:
		size = 4
		if out != nil {
			switch p.As {
			case obj.AUNDEF:
				out[0] = 0xe7fedef0 // NACL_INSTR_ARM_ABORT_NOW (UDF #0xEDE0)

			case APLD:
				out[0] = 0xe1a01001 // (MOVW R1, R1)
			}
		}

	case AB, ABL:
		if p.To.Type != obj.TYPE_MEM {
			if out != nil {
				c.asmout(p, o, out)
			}
		} else {
			if p.To.Offset != 0 || size != 4 || p.To.Reg > REG_R15 || p.To.Reg < REG_R0 {
				c.ctxt.Diag("unsupported instruction: %v", p)
			}
			if p.Pc&15 == 12 {
				p.Pc += 4
			}
			if out != nil {
				out[0] = ((uint32(p.Scond)&C_SCOND)^C_SCOND_XOR)<<28 | 0x03c0013f | (uint32(p.To.Reg)&15)<<12 | (uint32(p.To.Reg)&15)<<16 // BIC $0xc000000f, Rx
				if p.As == AB {
					out[1] = ((uint32(p.Scond)&C_SCOND)^C_SCOND_XOR)<<28 | 0x012fff10 | (uint32(p.To.Reg)&15)<<0 // BX Rx
				} else { // ABL
					out[1] = ((uint32(p.Scond)&C_SCOND)^C_SCOND_XOR)<<28 | 0x012fff30 | (uint32(p.To.Reg)&15)<<0 // BLX Rx
				}
			}

			size = 8
		}

		// align the last instruction (the actual BL) to the last instruction in a bundle
		if p.As == ABL {
			if p.To.Sym == deferreturn {
				p.Pc = ((int64(origPC) + 15) &^ 15) + 16 - int64(size)
			} else {
				p.Pc += (16 - ((p.Pc + int64(size)) & 15)) & 15
			}
		}

	case ALDREX,
		ALDREXD,
		AMOVB,
		AMOVBS,
		AMOVBU,
		AMOVD,
		AMOVF,
		AMOVH,
		AMOVHS,
		AMOVHU,
		AMOVM,
		AMOVW,
		ASTREX,
		ASTREXD:
		if p.To.Type == obj.TYPE_REG && p.To.Reg == REG_R15 && p.From.Reg == REG_R13 { // MOVW.W x(R13), PC
			if out != nil {
				c.asmout(p, o, out)
			}
			if size == 4 {
				if out != nil {
					// Note: 5c and 5g reg.c know that DIV/MOD smashes R12
					// so that this return instruction expansion is valid.
					out[0] = out[0] &^ 0x3000                                         // change PC to R12
					out[1] = ((uint32(p.Scond)&C_SCOND)^C_SCOND_XOR)<<28 | 0x03ccc13f // BIC $0xc000000f, R12
					out[2] = ((uint32(p.Scond)&C_SCOND)^C_SCOND_XOR)<<28 | 0x012fff1c // BX R12
				}

				size += 8
				if (p.Pc+int64(size))&15 == 4 {
					p.Pc += 4
				}
				break
			} else {
				// if the instruction used more than 4 bytes, then it must have used a very large
				// offset to update R13, so we need to additionally mask R13.
				if out != nil {
					out[size/4-1] &^= 0x3000                                                 // change PC to R12
					out[size/4] = ((uint32(p.Scond)&C_SCOND)^C_SCOND_XOR)<<28 | 0x03cdd103   // BIC $0xc0000000, R13
					out[size/4+1] = ((uint32(p.Scond)&C_SCOND)^C_SCOND_XOR)<<28 | 0x03ccc13f // BIC $0xc000000f, R12
					out[size/4+2] = ((uint32(p.Scond)&C_SCOND)^C_SCOND_XOR)<<28 | 0x012fff1c // BX R12
				}

				// p->pc+size is only ok at 4 or 12 mod 16.
				if (p.Pc+int64(size))%8 == 0 {
					p.Pc += 4
				}
				size += 12
				break
			}
		}

		if p.To.Type == obj.TYPE_REG && p.To.Reg == REG_R15 {
			c.ctxt.Diag("unsupported instruction (move to another register and use indirect jump instead): %v", p)
		}

		if p.To.Type == obj.TYPE_MEM && p.To.Reg == REG_R13 && (p.Scond&C_WBIT != 0) && size > 4 {
			// function prolog with very large frame size: MOVW.W R14,-100004(R13)
			// split it into two instructions:
			// 	ADD $-100004, R13
			// 	MOVW R14, 0(R13)
			q := c.newprog()

			p.Scond &^= C_WBIT
			*q = *p
			a := &p.To
			var a2 *obj.Addr
			if p.To.Type == obj.TYPE_MEM {
				a2 = &q.To
			} else {
				a2 = &q.From
			}
			nocache(q)
			nocache(p)

			// insert q after p
			q.Link = p.Link

			p.Link = q
			q.Pcond = nil

			// make p into ADD $X, R13
			p.As = AADD

			p.From = *a
			p.From.Reg = 0
			p.From.Type = obj.TYPE_CONST
			p.To = obj.Addr{}
			p.To.Type = obj.TYPE_REG
			p.To.Reg = REG_R13

			// make q into p but load/store from 0(R13)
			q.Spadj = 0

			*a2 = obj.Addr{}
			a2.Type = obj.TYPE_MEM
			a2.Reg = REG_R13
			a2.Sym = nil
			a2.Offset = 0
			size = int(c.oplook(p).size)
			break
		}

		if (p.To.Type == obj.TYPE_MEM && p.To.Reg != REG_R9) || // MOVW Rx, X(Ry), y != 9
			(p.From.Type == obj.TYPE_MEM && p.From.Reg != REG_R9) { // MOVW X(Rx), Ry, x != 9
			var a *obj.Addr
			if p.To.Type == obj.TYPE_MEM {
				a = &p.To
			} else {
				a = &p.From
			}
			reg := int(a.Reg)
			if size == 4 {
				// if addr.reg == 0, then it is probably load from x(FP) with small x, no need to modify.
				if reg == 0 {
					if out != nil {
						c.asmout(p, o, out)
					}
				} else {
					if out != nil {
						out[0] = ((uint32(p.Scond)&C_SCOND)^C_SCOND_XOR)<<28 | 0x03c00103 | (uint32(reg)&15)<<16 | (uint32(reg)&15)<<12 // BIC $0xc0000000, Rx
					}
					if p.Pc&15 == 12 {
						p.Pc += 4
					}
					size += 4
					if out != nil {
						c.asmout(p, o, out[1:])
					}
				}

				break
			} else {
				// if a load/store instruction takes more than 1 word to implement, then
				// we need to separate the instruction into two:
				// 1. explicitly load the address into R11.
				// 2. load/store from R11.
				// This won't handle .W/.P, so we should reject such code.
				if p.Scond&(C_PBIT|C_WBIT) != 0 {
					c.ctxt.Diag("unsupported instruction (.P/.W): %v", p)
				}
				q := c.newprog()
				*q = *p
				var a2 *obj.Addr
				if p.To.Type == obj.TYPE_MEM {
					a2 = &q.To
				} else {
					a2 = &q.From
				}
				nocache(q)
				nocache(p)

				// insert q after p
				q.Link = p.Link

				p.Link = q
				q.Pcond = nil

				// make p into MOVW $X(R), R11
				p.As = AMOVW

				p.From = *a
				p.From.Type = obj.TYPE_ADDR
				p.To = obj.Addr{}
				p.To.Type = obj.TYPE_REG
				p.To.Reg = REG_R11

				// make q into p but load/store from 0(R11)
				*a2 = obj.Addr{}

				a2.Type = obj.TYPE_MEM
				a2.Reg = REG_R11
				a2.Sym = nil
				a2.Offset = 0
				size = int(c.oplook(p).size)
				break
			}
		} else if out != nil {
			c.asmout(p, o, out)
		}
	}

	// destination register specific
	if p.To.Type == obj.TYPE_REG {
		switch p.To.Reg {
		case REG_R9:
			c.ctxt.Diag("invalid instruction, cannot write to R9: %v", p)

		case REG_R13:
			if out != nil {
				out[size/4] = 0xe3cdd103 // BIC $0xc0000000, R13
			}
			if (p.Pc+int64(size))&15 == 0 {
				p.Pc += 4
			}
			size += 4
		}
	}

	return size
}

func span5(ctxt *obj.Link, cursym *obj.LSym, newprog obj.ProgAlloc) {
	var p *obj.Prog
	var op *obj.Prog

	p = cursym.Func.Text
	if p == nil || p.Link == nil { // handle external functions and ELF section symbols
		return
	}

	if oprange[AAND&obj.AMask] == nil {
		ctxt.Diag("arm ops not initialized, call arm.buildop first")
	}

	c := ctxt5{ctxt: ctxt, newprog: newprog, cursym: cursym, autosize: p.To.Offset + 4}
	pc := int32(0)

	op = p
	p = p.Link
	var i int
	var m int
	var o *Optab
	for ; p != nil || c.blitrl != nil; op, p = p, p.Link {
		if p == nil {
			if c.checkpool(op, 0) {
				p = op
				continue
			}

			// can't happen: blitrl is not nil, but checkpool didn't flushpool
			ctxt.Diag("internal inconsistency")

			break
		}

		p.Pc = int64(pc)
		o = c.oplook(p)
		if ctxt.Headtype != objabi.Hnacl {
			m = int(o.size)
		} else {
			m = c.asmoutnacl(pc, p, o, nil)
			pc = int32(p.Pc) // asmoutnacl might change pc for alignment
			o = c.oplook(p)  // asmoutnacl might change p in rare cases
		}

		if m%4 != 0 || p.Pc%4 != 0 {
			ctxt.Diag("!pc invalid: %v size=%d", p, m)
		}

		// must check literal pool here in case p generates many instructions
		if c.blitrl != nil {
			i = m
			if c.checkpool(op, i) {
				p = op
				continue
			}
		}

		if m == 0 && (p.As != obj.AFUNCDATA && p.As != obj.APCDATA && p.As != ADATABUNDLEEND && p.As != obj.ANOP) {
			ctxt.Diag("zero-width instruction\n%v", p)
			continue
		}

		switch o.flag & (LFROM | LTO | LPOOL) {
		case LFROM:
			c.addpool(p, &p.From)

		case LTO:
			c.addpool(p, &p.To)

		case LPOOL:
			if p.Scond&C_SCOND == C_SCOND_NONE {
				c.flushpool(p, 0, 0)
			}
		}

		if p.As == AMOVW && p.To.Type == obj.TYPE_REG && p.To.Reg == REGPC && p.Scond&C_SCOND == C_SCOND_NONE {
			c.flushpool(p, 0, 0)
		}
		pc += int32(m)
	}

	c.cursym.Size = int64(pc)

	/*
	 * if any procedure is large enough to
	 * generate a large SBRA branch, then
	 * generate extra passes putting branches
	 * around jmps to fix. this is rare.
	 */
	times := 0

	var bflag int
	var opc int32
	var out [6 + 3]uint32
	for {
		bflag = 0
		pc = 0
		times++
		c.cursym.Func.Text.Pc = 0 // force re-layout the code.
		for p = c.cursym.Func.Text; p != nil; p = p.Link {
			o = c.oplook(p)
			if int64(pc) > p.Pc {
				p.Pc = int64(pc)
			}

			/* very large branches
			if(o->type == 6 && p->pcond) {
				otxt = p->pcond->pc - c;
				if(otxt < 0)
					otxt = -otxt;
				if(otxt >= (1L<<17) - 10) {
					q = emallocz(sizeof(Prog));
					q->link = p->link;
					p->link = q;
					q->as = AB;
					q->to.type = TYPE_BRANCH;
					q->pcond = p->pcond;
					p->pcond = q;
					q = emallocz(sizeof(Prog));
					q->link = p->link;
					p->link = q;
					q->as = AB;
					q->to.type = TYPE_BRANCH;
					q->pcond = q->link->link;
					bflag = 1;
				}
			}
			*/
			opc = int32(p.Pc)

			if ctxt.Headtype != objabi.Hnacl {
				m = int(o.size)
			} else {
				m = c.asmoutnacl(pc, p, o, nil)
			}
			if p.Pc != int64(opc) {
				bflag = 1
			}

			//print("%v pc changed %d to %d in iter. %d\n", p, opc, (int32)p->pc, times);
			pc = int32(p.Pc + int64(m))

			if m%4 != 0 || p.Pc%4 != 0 {
				ctxt.Diag("pc invalid: %v size=%d", p, m)
			}

			if m/4 > len(out) {
				ctxt.Diag("instruction size too large: %d > %d", m/4, len(out))
			}
			if m == 0 && (p.As != obj.AFUNCDATA && p.As != obj.APCDATA && p.As != ADATABUNDLEEND && p.As != obj.ANOP) {
				if p.As == obj.ATEXT {
					c.autosize = p.To.Offset + 4
					continue
				}

				ctxt.Diag("zero-width instruction\n%v", p)
				continue
			}
		}

		c.cursym.Size = int64(pc)
		if bflag == 0 {
			break
		}
	}

	if pc%4 != 0 {
		ctxt.Diag("sym->size=%d, invalid", pc)
	}

	/*
	 * lay out the code.  all the pc-relative code references,
	 * even cross-function, are resolved now;
	 * only data references need to be relocated.
	 * with more work we could leave cross-function
	 * code references to be relocated too, and then
	 * perhaps we'd be able to parallelize the span loop above.
	 */

	p = c.cursym.Func.Text
	c.autosize = p.To.Offset + 4
	c.cursym.Grow(c.cursym.Size)

	bp := c.cursym.P
	pc = int32(p.Pc) // even p->link might need extra padding
	var v int
	for p = p.Link; p != nil; p = p.Link {
		c.pc = p.Pc
		o = c.oplook(p)
		opc = int32(p.Pc)
		if ctxt.Headtype != objabi.Hnacl {
			c.asmout(p, o, out[:])
			m = int(o.size)
		} else {
			m = c.asmoutnacl(pc, p, o, out[:])
			if int64(opc) != p.Pc {
				ctxt.Diag("asmoutnacl broken: pc changed (%d->%d) in last stage: %v", opc, int32(p.Pc), p)
			}
		}

		if m%4 != 0 || p.Pc%4 != 0 {
			ctxt.Diag("final stage: pc invalid: %v size=%d", p, m)
		}

		if int64(pc) > p.Pc {
			ctxt.Diag("PC padding invalid: want %#d, has %#d: %v", p.Pc, pc, p)
		}
		for int64(pc) != p.Pc {
			// emit 0xe1a00000 (MOVW R0, R0)
			bp[0] = 0x00
			bp = bp[1:]

			bp[0] = 0x00
			bp = bp[1:]
			bp[0] = 0xa0
			bp = bp[1:]
			bp[0] = 0xe1
			bp = bp[1:]
			pc += 4
		}

		for i = 0; i < m/4; i++ {
			v = int(out[i])
			bp[0] = byte(v)
			bp = bp[1:]
			bp[0] = byte(v >> 8)
			bp = bp[1:]
			bp[0] = byte(v >> 16)
			bp = bp[1:]
			bp[0] = byte(v >> 24)
			bp = bp[1:]
		}

		pc += int32(m)
	}
}

/*
 * when the first reference to the literal pool threatens
 * to go out of range of a 12-bit PC-relative offset,
 * drop the pool now, and branch round it.
 * this happens only in extended basic blocks that exceed 4k.
 */
func (c *ctxt5) checkpool(p *obj.Prog, sz int) bool {
	if c.pool.size >= 0xff0 || immaddr(int32((p.Pc+int64(sz)+4)+4+int64(12+c.pool.size)-int64(c.pool.start+8))) == 0 {
		return c.flushpool(p, 1, 0)
	} else if p.Link == nil {
		return c.flushpool(p, 2, 0)
	}
	return false
}

func (c *ctxt5) flushpool(p *obj.Prog, skip int, force int) bool {
	if c.blitrl != nil {
		if skip != 0 {
			if false && skip == 1 {
				fmt.Printf("note: flush literal pool at %x: len=%d ref=%x\n", uint64(p.Pc+4), c.pool.size, c.pool.start)
			}
			q := c.newprog()
			q.As = AB
			q.To.Type = obj.TYPE_BRANCH
			q.Pcond = p.Link
			q.Link = c.blitrl
			q.Pos = p.Pos
			c.blitrl = q
		} else if force == 0 && (p.Pc+int64(12+c.pool.size)-int64(c.pool.start) < 2048) { // 12 take into account the maximum nacl literal pool alignment padding size
			return false
		}
		if c.ctxt.Headtype == objabi.Hnacl && c.pool.size%16 != 0 {
			// if pool is not multiple of 16 bytes, add an alignment marker
			q := c.newprog()

			q.As = ADATABUNDLEEND
			c.elitrl.Link = q
			c.elitrl = q
		}

		// The line number for constant pool entries doesn't really matter.
		// We set it to the line number of the preceding instruction so that
		// there are no deltas to encode in the pc-line tables.
		for q := c.blitrl; q != nil; q = q.Link {
			q.Pos = p.Pos
		}

		c.elitrl.Link = p.Link
		p.Link = c.blitrl

		c.blitrl = nil /* BUG: should refer back to values until out-of-range */
		c.elitrl = nil
		c.pool.size = 0
		c.pool.start = 0
		c.pool.extra = 0
		return true
	}

	return false
}

func (c *ctxt5) addpool(p *obj.Prog, a *obj.Addr) {
	t := c.newprog()
	t.As = AWORD

	switch c.aclass(a) {
	default:
		t.To.Offset = a.Offset
		t.To.Sym = a.Sym
		t.To.Type = a.Type
		t.To.Name = a.Name

		if c.ctxt.Flag_shared && t.To.Sym != nil {
			t.Rel = p
		}

	case C_SROREG,
		C_LOREG,
		C_ROREG,
		C_FOREG,
		C_SOREG,
		C_HOREG,
		C_FAUTO,
		C_SAUTO,
		C_LAUTO,
		C_LACON:
		t.To.Type = obj.TYPE_CONST
		t.To.Offset = c.instoffset
	}

	if t.Rel == nil {
		for q := c.blitrl; q != nil; q = q.Link { /* could hash on t.t0.offset */
			if q.Rel == nil && q.To == t.To {
				p.Pcond = q
				return
			}
		}
	}

	if c.ctxt.Headtype == objabi.Hnacl && c.pool.size%16 == 0 {
		// start a new data bundle
		q := c.newprog()
		q.As = ADATABUNDLE
		q.Pc = int64(c.pool.size)
		c.pool.size += 4
		if c.blitrl == nil {
			c.blitrl = q
			c.pool.start = uint32(p.Pc)
		} else {
			c.elitrl.Link = q
		}

		c.elitrl = q
	}

	q := c.newprog()
	*q = *t
	q.Pc = int64(c.pool.size)

	if c.blitrl == nil {
		c.blitrl = q
		c.pool.start = uint32(p.Pc)
	} else {
		c.elitrl.Link = q
	}
	c.elitrl = q
	c.pool.size += 4

	p.Pcond = q
}

func (c *ctxt5) regoff(a *obj.Addr) int32 {
	c.instoffset = 0
	c.aclass(a)
	return int32(c.instoffset)
}

func immrot(v uint32) int32 {
	for i := 0; i < 16; i++ {
		if v&^0xff == 0 {
			return int32(uint32(int32(i)<<8) | v | 1<<25)
		}
		v = v<<2 | v>>30
	}

	return 0
}

func immaddr(v int32) int32 {
	if v >= 0 && v <= 0xfff {
		return v&0xfff | 1<<24 | 1<<23 /* pre indexing */ /* pre indexing, up */
	}
	if v >= -0xfff && v < 0 {
		return -v&0xfff | 1<<24 /* pre indexing */
	}
	return 0
}

func immfloat(v int32) bool {
	return v&0xC03 == 0 /* offset will fit in floating-point load/store */
}

func immhalf(v int32) bool {
	if v >= 0 && v <= 0xff {
		return v|1<<24|1<<23 != 0 /* pre indexing */ /* pre indexing, up */
	}
	if v >= -0xff && v < 0 {
		return -v&0xff|1<<24 != 0 /* pre indexing */
	}
	return false
}

func (c *ctxt5) aclass(a *obj.Addr) int {
	switch a.Type {
	case obj.TYPE_NONE:
		return C_NONE

	case obj.TYPE_REG:
		c.instoffset = 0
		if REG_R0 <= a.Reg && a.Reg <= REG_R15 {
			return C_REG
		}
		if REG_F0 <= a.Reg && a.Reg <= REG_F15 {
			return C_FREG
		}
		if a.Reg == REG_FPSR || a.Reg == REG_FPCR {
			return C_FCR
		}
		if a.Reg == REG_CPSR || a.Reg == REG_SPSR {
			return C_PSR
		}
		return C_GOK

	case obj.TYPE_REGREG:
		return C_REGREG

	case obj.TYPE_REGREG2:
		return C_REGREG2

	case obj.TYPE_REGLIST:
		return C_REGLIST

	case obj.TYPE_SHIFT:
		return C_SHIFT

	case obj.TYPE_MEM:
		switch a.Name {
		case obj.NAME_EXTERN,
			obj.NAME_GOTREF,
			obj.NAME_STATIC:
			if a.Sym == nil || a.Sym.Name == "" {
				fmt.Printf("null sym external\n")
				return C_GOK
			}

			c.instoffset = 0 // s.b. unused but just in case
			if a.Sym.Type == objabi.STLSBSS {
				if c.ctxt.Flag_shared {
					return C_TLS_IE
				} else {
					return C_TLS_LE
				}
			}

			return C_ADDR

		case obj.NAME_AUTO:
			c.instoffset = c.autosize + a.Offset
			if t := immaddr(int32(c.instoffset)); t != 0 {
				if immhalf(int32(c.instoffset)) {
					if immfloat(t) {
						return C_HFAUTO
					}
					return C_HAUTO
				}

				if immfloat(t) {
					return C_FAUTO
				}
				return C_SAUTO
			}

			return C_LAUTO

		case obj.NAME_PARAM:
			c.instoffset = c.autosize + a.Offset + 4
			if t := immaddr(int32(c.instoffset)); t != 0 {
				if immhalf(int32(c.instoffset)) {
					if immfloat(t) {
						return C_HFAUTO
					}
					return C_HAUTO
				}

				if immfloat(t) {
					return C_FAUTO
				}
				return C_SAUTO
			}

			return C_LAUTO

		case obj.NAME_NONE:
			c.instoffset = a.Offset
			if t := immaddr(int32(c.instoffset)); t != 0 {
				if immhalf(int32(c.instoffset)) { /* n.b. that it will also satisfy immrot */
					if immfloat(t) {
						return C_HFOREG
					}
					return C_HOREG
				}

				if immfloat(t) {
					return C_FOREG /* n.b. that it will also satisfy immrot */
				}
				if immrot(uint32(c.instoffset)) != 0 {
					return C_SROREG
				}
				if immhalf(int32(c.instoffset)) {
					return C_HOREG
				}
				return C_SOREG
			}

			if immrot(uint32(c.instoffset)) != 0 {
				return C_ROREG
			}
			return C_LOREG
		}

		return C_GOK

	case obj.TYPE_FCONST:
		if c.chipzero5(a.Val.(float64)) >= 0 {
			return C_ZFCON
		}
		if c.chipfloat5(a.Val.(float64)) >= 0 {
			return C_SFCON
		}
		return C_LFCON

	case obj.TYPE_TEXTSIZE:
		return C_TEXTSIZE

	case obj.TYPE_CONST,
		obj.TYPE_ADDR:
		switch a.Name {
		case obj.NAME_NONE:
			c.instoffset = a.Offset
			if a.Reg != 0 {
				return c.aconsize()
			}

			if immrot(uint32(c.instoffset)) != 0 {
				return C_RCON
			}
			if immrot(^uint32(c.instoffset)) != 0 {
				return C_NCON
			}
			if uint32(c.instoffset) <= 0xffff && objabi.GOARM == 7 {
				return C_SCON
			}
			return C_LCON

		case obj.NAME_EXTERN,
			obj.NAME_GOTREF,
			obj.NAME_STATIC:
			s := a.Sym
			if s == nil {
				break
			}
			c.instoffset = 0 // s.b. unused but just in case
			return C_LCONADDR

		case obj.NAME_AUTO:
			c.instoffset = c.autosize + a.Offset
			return c.aconsize()

		case obj.NAME_PARAM:
			c.instoffset = c.autosize + a.Offset + 4
			return c.aconsize()
		}

		return C_GOK

	case obj.TYPE_BRANCH:
		return C_SBRA
	}

	return C_GOK
}

func (c *ctxt5) aconsize() int {
	if immrot(uint32(c.instoffset)) != 0 {
		return C_RACON
	}
	if immrot(uint32(-c.instoffset)) != 0 {
		return C_RACON
	}
	return C_LACON
}

func (c *ctxt5) oplook(p *obj.Prog) *Optab {
	a1 := int(p.Optab)
	if a1 != 0 {
		return &optab[a1-1]
	}
	a1 = int(p.From.Class)
	if a1 == 0 {
		a1 = c.aclass(&p.From) + 1
		p.From.Class = int8(a1)
	}

	a1--
	a3 := int(p.To.Class)
	if a3 == 0 {
		a3 = c.aclass(&p.To) + 1
		p.To.Class = int8(a3)
	}

	a3--
	a2 := C_NONE
	if p.Reg != 0 {
		a2 = C_REG
	}

	if false { /*debug['O']*/
		fmt.Printf("oplook %v %v %v %v\n", p.As, DRconv(a1), DRconv(a2), DRconv(a3))
		fmt.Printf("\t\t%d %d\n", p.From.Type, p.To.Type)
	}

	ops := oprange[p.As&obj.AMask]
	c1 := &xcmp[a1]
	c3 := &xcmp[a3]
	for i := range ops {
		op := &ops[i]
		if int(op.a2) == a2 && c1[op.a1] && c3[op.a3] {
			p.Optab = uint16(cap(optab) - cap(ops) + i + 1)
			return op
		}
	}

	c.ctxt.Diag("illegal combination %v; %v %v %v, %d %d", p, DRconv(a1), DRconv(a2), DRconv(a3), p.From.Type, p.To.Type)
	c.ctxt.Diag("from %d %d to %d %d\n", p.From.Type, p.From.Name, p.To.Type, p.To.Name)
	if ops == nil {
		ops = optab
	}
	return &ops[0]
}

func cmp(a int, b int) bool {
	if a == b {
		return true
	}
	switch a {
	case C_LCON:
		if b == C_RCON || b == C_NCON || b == C_SCON {
			return true
		}

	case C_LACON:
		if b == C_RACON {
			return true
		}

	case C_LFCON:
		if b == C_ZFCON || b == C_SFCON {
			return true
		}

	case C_HFAUTO:
		return b == C_HAUTO || b == C_FAUTO

	case C_FAUTO, C_HAUTO:
		return b == C_HFAUTO

	case C_SAUTO:
		return cmp(C_HFAUTO, b)

	case C_LAUTO:
		return cmp(C_SAUTO, b)

	case C_HFOREG:
		return b == C_HOREG || b == C_FOREG

	case C_FOREG, C_HOREG:
		return b == C_HFOREG

	case C_SROREG:
		return cmp(C_SOREG, b) || cmp(C_ROREG, b)

	case C_SOREG, C_ROREG:
		return b == C_SROREG || cmp(C_HFOREG, b)

	case C_LOREG:
		return cmp(C_SROREG, b)

	case C_LBRA:
		if b == C_SBRA {
			return true
		}

	case C_HREG:
		return cmp(C_SP, b) || cmp(C_PC, b)
	}

	return false
}

type ocmp []Optab

func (x ocmp) Len() int {
	return len(x)
}

func (x ocmp) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x ocmp) Less(i, j int) bool {
	p1 := &x[i]
	p2 := &x[j]
	n := int(p1.as) - int(p2.as)
	if n != 0 {
		return n < 0
	}
	n = int(p1.a1) - int(p2.a1)
	if n != 0 {
		return n < 0
	}
	n = int(p1.a2) - int(p2.a2)
	if n != 0 {
		return n < 0
	}
	n = int(p1.a3) - int(p2.a3)
	if n != 0 {
		return n < 0
	}
	return false
}

func opset(a, b0 obj.As) {
	oprange[a&obj.AMask] = oprange[b0]
}

func buildop(ctxt *obj.Link) {
	if oprange[AAND&obj.AMask] != nil {
		// Already initialized; stop now.
		// This happens in the cmd/asm tests,
		// each of which re-initializes the arch.
		return
	}

	deferreturn = ctxt.Lookup("runtime.deferreturn")

	symdiv = ctxt.Lookup("_div")
	symdivu = ctxt.Lookup("_divu")
	symmod = ctxt.Lookup("_mod")
	symmodu = ctxt.Lookup("_modu")

	var n int

	for i := 0; i < C_GOK; i++ {
		for n = 0; n < C_GOK; n++ {
			if cmp(n, i) {
				xcmp[i][n] = true
			}
		}
	}
	for n = 0; optab[n].as != obj.AXXX; n++ {
		if optab[n].flag&LPCREL != 0 {
			if ctxt.Flag_shared {
				optab[n].size += int8(optab[n].pcrelsiz)
			} else {
				optab[n].flag &^= LPCREL
			}
		}
	}

	sort.Sort(ocmp(optab[:n]))
	for i := 0; i < n; i++ {
		r := optab[i].as
		r0 := r & obj.AMask
		start := i
		for optab[i].as == r {
			i++
		}
		oprange[r0] = optab[start:i]
		i--

		switch r {
		default:
			ctxt.Diag("unknown op in build: %v", r)
			log.Fatalf("bad code")

		case AADD:
			opset(AAND, r0)
			opset(AEOR, r0)
			opset(ASUB, r0)
			opset(ARSB, r0)
			opset(AADC, r0)
			opset(ASBC, r0)
			opset(ARSC, r0)
			opset(AORR, r0)
			opset(ABIC, r0)

		case ACMP:
			opset(ATEQ, r0)
			opset(ACMN, r0)
			opset(ATST, r0)

		case AMVN:
			break

		case ABEQ:
			opset(ABNE, r0)
			opset(ABCS, r0)
			opset(ABHS, r0)
			opset(ABCC, r0)
			opset(ABLO, r0)
			opset(ABMI, r0)
			opset(ABPL, r0)
			opset(ABVS, r0)
			opset(ABVC, r0)
			opset(ABHI, r0)
			opset(ABLS, r0)
			opset(ABGE, r0)
			opset(ABLT, r0)
			opset(ABGT, r0)
			opset(ABLE, r0)

		case ASLL:
			opset(ASRL, r0)
			opset(ASRA, r0)

		case AMUL:
			opset(AMULU, r0)

		case ADIV:
			opset(AMOD, r0)
			opset(AMODU, r0)
			opset(ADIVU, r0)

		case ADIVHW:
			opset(ADIVUHW, r0)

		case AMOVW,
			AMOVB,
			AMOVBS,
			AMOVBU,
			AMOVH,
			AMOVHS,
			AMOVHU:
			break

		case ASWPW:
			opset(ASWPBU, r0)

		case AB,
			ABL,
			ABX,
			ABXRET,
			obj.ADUFFZERO,
			obj.ADUFFCOPY,
			ASWI,
			AWORD,
			AMOVM,
			ARFE,
			obj.ATEXT:
			break

		case AADDF:
			opset(AADDD, r0)
			opset(ASUBF, r0)
			opset(ASUBD, r0)
			opset(AMULF, r0)
			opset(AMULD, r0)
			opset(ADIVF, r0)
			opset(ADIVD, r0)
			opset(ASQRTF, r0)
			opset(ASQRTD, r0)
			opset(AMOVFD, r0)
			opset(AMOVDF, r0)
			opset(AABSF, r0)
			opset(AABSD, r0)
			opset(ANEGF, r0)
			opset(ANEGD, r0)

		case ACMPF:
			opset(ACMPD, r0)

		case AMOVF:
			opset(AMOVD, r0)

		case AMOVFW:
			opset(AMOVDW, r0)

		case AMOVWF:
			opset(AMOVWD, r0)

		case AMULL:
			opset(AMULAL, r0)
			opset(AMULLU, r0)
			opset(AMULALU, r0)

		case AMULWT:
			opset(AMULWB, r0)
			opset(AMULBB, r0)
			opset(AMMUL, r0)

		case AMULAWT:
			opset(AMULAWB, r0)
			opset(AMULABB, r0)
			opset(AMULS, r0)
			opset(AMMULA, r0)
			opset(AMMULS, r0)

		case ACLZ:
			opset(AREV, r0)
			opset(AREV16, r0)
			opset(AREVSH, r0)
			opset(ARBIT, r0)

		case AMULA,
			ALDREX,
			ASTREX,
			ALDREXD,
			ASTREXD,
			ATST,
			APLD,
			obj.AUNDEF,
			obj.AFUNCDATA,
			obj.APCDATA,
			obj.ANOP,
			ADATABUNDLE,
			ADATABUNDLEEND:
			break
		}
	}
}

func (c *ctxt5) asmout(p *obj.Prog, o *Optab, out []uint32) {
	c.printp = p
	o1 := uint32(0)
	o2 := uint32(0)
	o3 := uint32(0)
	o4 := uint32(0)
	o5 := uint32(0)
	o6 := uint32(0)
	if false { /*debug['P']*/
		fmt.Printf("%x: %v\ttype %d\n", uint32(p.Pc), p, o.type_)
	}
	switch o.type_ {
	default:
		c.ctxt.Diag("%v: unknown asm %d", p, o.type_)

	case 0: /* pseudo ops */
		if false { /*debug['G']*/
			fmt.Printf("%x: %s: arm\n", uint32(p.Pc), p.From.Sym.Name)
		}

	case 1: /* op R,[R],R */
		o1 = c.oprrr(p, p.As, int(p.Scond))

		rf := int(p.From.Reg)
		rt := int(p.To.Reg)
		r := int(p.Reg)
		if p.To.Type == obj.TYPE_NONE {
			rt = 0
		}
		if p.As == AMOVB || p.As == AMOVH || p.As == AMOVW || p.As == AMVN {
			r = 0
		} else if r == 0 {
			r = rt
		}
		o1 |= (uint32(rf)&15)<<0 | (uint32(r)&15)<<16 | (uint32(rt)&15)<<12

	case 2: /* movbu $I,[R],R */
		c.aclass(&p.From)

		o1 = c.oprrr(p, p.As, int(p.Scond))
		o1 |= uint32(immrot(uint32(c.instoffset)))
		rt := int(p.To.Reg)
		r := int(p.Reg)
		if p.To.Type == obj.TYPE_NONE {
			rt = 0
		}
		if p.As == AMOVW || p.As == AMVN {
			r = 0
		} else if r == 0 {
			r = rt
		}
		o1 |= (uint32(r)&15)<<16 | (uint32(rt)&15)<<12

	case 3: /* add R<<[IR],[R],R */
		o1 = c.mov(p)

	case 4: /* MOVW $off(R), R -> add $off,[R],R */
		c.aclass(&p.From)
		if c.instoffset < 0 {
			o1 = c.oprrr(p, ASUB, int(p.Scond))
			o1 |= uint32(immrot(uint32(-c.instoffset)))
		} else {
			o1 = c.oprrr(p, AADD, int(p.Scond))
			o1 |= uint32(immrot(uint32(c.instoffset)))
		}
		r := int(p.From.Reg)
		if r == 0 {
			r = int(o.param)
		}
		o1 |= (uint32(r) & 15) << 16
		o1 |= (uint32(p.To.Reg) & 15) << 12

	case 5: /* bra s */
		o1 = c.opbra(p, p.As, int(p.Scond))

		v := int32(-8)
		if p.To.Sym != nil {
			rel := obj.Addrel(c.cursym)
			rel.Off = int32(c.pc)
			rel.Siz = 4
			rel.Sym = p.To.Sym
			v += int32(p.To.Offset)
			rel.Add = int64(o1) | (int64(v)>>2)&0xffffff
			rel.Type = objabi.R_CALLARM
			break
		}

		if p.Pcond != nil {
			v = int32((p.Pcond.Pc - c.pc) - 8)
		}
		o1 |= (uint32(v) >> 2) & 0xffffff

	case 6: /* b ,O(R) -> add $O,R,PC */
		c.aclass(&p.To)

		o1 = c.oprrr(p, AADD, int(p.Scond))
		o1 |= uint32(immrot(uint32(c.instoffset)))
		o1 |= (uint32(p.To.Reg) & 15) << 16
		o1 |= (REGPC & 15) << 12

	case 7: /* bl (R) -> blx R */
		c.aclass(&p.To)

		if c.instoffset != 0 {
			c.ctxt.Diag("%v: doesn't support BL offset(REG) with non-zero offset %d", p, c.instoffset)
		}
		o1 = c.oprrr(p, ABL, int(p.Scond))
		o1 |= (uint32(p.To.Reg) & 15) << 0
		rel := obj.Addrel(c.cursym)
		rel.Off = int32(c.pc)
		rel.Siz = 0
		rel.Type = objabi.R_CALLIND

	case 8: /* sll $c,[R],R -> mov (R<<$c),R */
		c.aclass(&p.From)

		o1 = c.oprrr(p, p.As, int(p.Scond))
		r := int(p.Reg)
		if r == 0 {
			r = int(p.To.Reg)
		}
		o1 |= (uint32(r) & 15) << 0
		o1 |= uint32((c.instoffset & 31) << 7)
		o1 |= (uint32(p.To.Reg) & 15) << 12

	case 9: /* sll R,[R],R -> mov (R<<R),R */
		o1 = c.oprrr(p, p.As, int(p.Scond))

		r := int(p.Reg)
		if r == 0 {
			r = int(p.To.Reg)
		}
		o1 |= (uint32(r) & 15) << 0
		o1 |= (uint32(p.From.Reg)&15)<<8 | 1<<4
		o1 |= (uint32(p.To.Reg) & 15) << 12

	case 10: /* swi [$con] */
		o1 = c.oprrr(p, p.As, int(p.Scond))

		if p.To.Type != obj.TYPE_NONE {
			c.aclass(&p.To)
			o1 |= uint32(c.instoffset & 0xffffff)
		}

	case 11: /* word */
		c.aclass(&p.To)

		o1 = uint32(c.instoffset)
		if p.To.Sym != nil {
			// This case happens with words generated
			// in the PC stream as part of the literal pool (c.pool).
			rel := obj.Addrel(c.cursym)

			rel.Off = int32(c.pc)
			rel.Siz = 4
			rel.Sym = p.To.Sym
			rel.Add = p.To.Offset

			if c.ctxt.Flag_shared {
				if p.To.Name == obj.NAME_GOTREF {
					rel.Type = objabi.R_GOTPCREL
				} else {
					rel.Type = objabi.R_PCREL
				}
				rel.Add += c.pc - p.Rel.Pc - 8
			} else {
				rel.Type = objabi.R_ADDR
			}
			o1 = 0
		}

	case 12: /* movw $lcon, reg */
		if o.a1 == C_SCON {
			o1 = c.omvs(p, &p.From, int(p.To.Reg))
		} else {
			o1 = c.omvl(p, &p.From, int(p.To.Reg))
		}

		if o.flag&LPCREL != 0 {
			o2 = c.oprrr(p, AADD, int(p.Scond)) | (uint32(p.To.Reg)&15)<<0 | (REGPC&15)<<16 | (uint32(p.To.Reg)&15)<<12
		}

	case 13: /* op $lcon, [R], R */
		if o.a1 == C_SCON {
			o1 = c.omvs(p, &p.From, REGTMP)
		} else {
			o1 = c.omvl(p, &p.From, REGTMP)
		}

		if o1 == 0 {
			break
		}
		o2 = c.oprrr(p, p.As, int(p.Scond))
		o2 |= REGTMP & 15
		r := int(p.Reg)
		if p.As == AMOVW || p.As == AMVN {
			r = 0
		} else if r == 0 {
			r = int(p.To.Reg)
		}
		o2 |= (uint32(r) & 15) << 16
		if p.To.Type != obj.TYPE_NONE {
			o2 |= (uint32(p.To.Reg) & 15) << 12
		}

	case 14: /* movb/movbu/movh/movhu R,R */
		o1 = c.oprrr(p, ASLL, int(p.Scond))

		if p.As == AMOVBU || p.As == AMOVHU {
			o2 = c.oprrr(p, ASRL, int(p.Scond))
		} else {
			o2 = c.oprrr(p, ASRA, int(p.Scond))
		}

		r := int(p.To.Reg)
		o1 |= (uint32(p.From.Reg)&15)<<0 | (uint32(r)&15)<<12
		o2 |= uint32(r)&15 | (uint32(r)&15)<<12
		if p.As == AMOVB || p.As == AMOVBS || p.As == AMOVBU {
			o1 |= 24 << 7
			o2 |= 24 << 7
		} else {
			o1 |= 16 << 7
			o2 |= 16 << 7
		}

	case 15: /* mul r,[r,]r */
		o1 = c.oprrr(p, p.As, int(p.Scond))

		rf := int(p.From.Reg)
		rt := int(p.To.Reg)
		r := int(p.Reg)
		if r == 0 {
			r = rt
		}
		if rt == r {
			r = rf
			rf = rt
		}

		if false {
			if rt == r || rf == REGPC&15 || r == REGPC&15 || rt == REGPC&15 {
				c.ctxt.Diag("%v: bad registers in MUL", p)
			}
		}

		o1 |= (uint32(rf)&15)<<8 | (uint32(r)&15)<<0 | (uint32(rt)&15)<<16

	case 16: /* div r,[r,]r */
		o1 = 0xf << 28

		o2 = 0

	case 17:
		o1 = c.oprrr(p, p.As, int(p.Scond))
		rf := int(p.From.Reg)
		rt := int(p.To.Reg)
		rt2 := int(p.To.Offset)
		r := int(p.Reg)
		o1 |= (uint32(rf)&15)<<8 | (uint32(r)&15)<<0 | (uint32(rt)&15)<<16 | (uint32(rt2)&15)<<12

	case 20: /* mov/movb/movbu R,O(R) */
		c.aclass(&p.To)

		r := int(p.To.Reg)
		if r == 0 {
			r = int(o.param)
		}
		o1 = c.osr(p.As, int(p.From.Reg), int32(c.instoffset), r, int(p.Scond))

	case 21: /* mov/movbu O(R),R -> lr */
		c.aclass(&p.From)

		r := int(p.From.Reg)
		if r == 0 {
			r = int(o.param)
		}
		o1 = c.olr(int32(c.instoffset), r, int(p.To.Reg), int(p.Scond))
		if p.As != AMOVW {
			o1 |= 1 << 22
		}

	case 30: /* mov/movb/movbu R,L(R) */
		o1 = c.omvl(p, &p.To, REGTMP)

		if o1 == 0 {
			break
		}
		r := int(p.To.Reg)
		if r == 0 {
			r = int(o.param)
		}
		o2 = c.osrr(int(p.From.Reg), REGTMP&15, r, int(p.Scond))
		if p.As != AMOVW {
			o2 |= 1 << 22
		}

	case 31: /* mov/movbu L(R),R -> lr[b] */
		o1 = c.omvl(p, &p.From, REGTMP)

		if o1 == 0 {
			break
		}
		r := int(p.From.Reg)
		if r == 0 {
			r = int(o.param)
		}
		o2 = c.olrr(REGTMP&15, r, int(p.To.Reg), int(p.Scond))
		if p.As == AMOVBU || p.As == AMOVBS || p.As == AMOVB {
			o2 |= 1 << 22
		}

	case 34: /* mov $lacon,R */
		o1 = c.omvl(p, &p.From, REGTMP)

		if o1 == 0 {
			break
		}

		o2 = c.oprrr(p, AADD, int(p.Scond))
		o2 |= REGTMP & 15
		r := int(p.From.Reg)
		if r == 0 {
			r = int(o.param)
		}
		o2 |= (uint32(r) & 15) << 16
		if p.To.Type != obj.TYPE_NONE {
			o2 |= (uint32(p.To.Reg) & 15) << 12
		}

	case 35: /* mov PSR,R */
		o1 = 2<<23 | 0xf<<16 | 0<<0

		o1 |= ((uint32(p.Scond) & C_SCOND) ^ C_SCOND_XOR) << 28
		o1 |= (uint32(p.From.Reg) & 1) << 22
		o1 |= (uint32(p.To.Reg) & 15) << 12

	case 36: /* mov R,PSR */
		o1 = 2<<23 | 0x29f<<12 | 0<<4

		if p.Scond&C_FBIT != 0 {
			o1 ^= 0x010 << 12
		}
		o1 |= ((uint32(p.Scond) & C_SCOND) ^ C_SCOND_XOR) << 28
		o1 |= (uint32(p.To.Reg) & 1) << 22
		o1 |= (uint32(p.From.Reg) & 15) << 0

	case 37: /* mov $con,PSR */
		c.aclass(&p.From)

		o1 = 2<<23 | 0x29f<<12 | 0<<4
		if p.Scond&C_FBIT != 0 {
			o1 ^= 0x010 << 12
		}
		o1 |= ((uint32(p.Scond) & C_SCOND) ^ C_SCOND_XOR) << 28
		o1 |= uint32(immrot(uint32(c.instoffset)))
		o1 |= (uint32(p.To.Reg) & 1) << 22
		o1 |= (uint32(p.From.Reg) & 15) << 0

	case 38, 39:
		switch o.type_ {
		case 38: /* movm $con,oreg -> stm */
			o1 = 0x4 << 25

			o1 |= uint32(p.From.Offset & 0xffff)
			o1 |= (uint32(p.To.Reg) & 15) << 16
			c.aclass(&p.To)

		case 39: /* movm oreg,$con -> ldm */
			o1 = 0x4<<25 | 1<<20

			o1 |= uint32(p.To.Offset & 0xffff)
			o1 |= (uint32(p.From.Reg) & 15) << 16
			c.aclass(&p.From)
		}

		if c.instoffset != 0 {
			c.ctxt.Diag("offset must be zero in MOVM; %v", p)
		}
		o1 |= ((uint32(p.Scond) & C_SCOND) ^ C_SCOND_XOR) << 28
		if p.Scond&C_PBIT != 0 {
			o1 |= 1 << 24
		}
		if p.Scond&C_UBIT != 0 {
			o1 |= 1 << 23
		}
		if p.Scond&C_SBIT != 0 {
			o1 |= 1 << 22
		}
		if p.Scond&C_WBIT != 0 {
			o1 |= 1 << 21
		}

	case 40: /* swp oreg,reg,reg */
		c.aclass(&p.From)

		if c.instoffset != 0 {
			c.ctxt.Diag("offset must be zero in SWP")
		}
		o1 = 0x2<<23 | 0x9<<4
		if p.As != ASWPW {
			o1 |= 1 << 22
		}
		o1 |= (uint32(p.From.Reg) & 15) << 16
		o1 |= (uint32(p.Reg) & 15) << 0
		o1 |= (uint32(p.To.Reg) & 15) << 12
		o1 |= ((uint32(p.Scond) & C_SCOND) ^ C_SCOND_XOR) << 28

	case 41: /* rfe -> movm.s.w.u 0(r13),[r15] */
		o1 = 0xe8fd8000

	case 50: /* floating point store */
		v := c.regoff(&p.To)

		r := int(p.To.Reg)
		if r == 0 {
			r = int(o.param)
		}
		o1 = c.ofsr(p.As, int(p.From.Reg), v, r, int(p.Scond), p)

	case 51: /* floating point load */
		v := c.regoff(&p.From)

		r := int(p.From.Reg)
		if r == 0 {
			r = int(o.param)
		}
		o1 = c.ofsr(p.As, int(p.To.Reg), v, r, int(p.Scond), p) | 1<<20

	case 52: /* floating point store, int32 offset UGLY */
		o1 = c.omvl(p, &p.To, REGTMP)

		if o1 == 0 {
			break
		}
		r := int(p.To.Reg)
		if r == 0 {
			r = int(o.param)
		}
		o2 = c.oprrr(p, AADD, int(p.Scond)) | (REGTMP&15)<<12 | (REGTMP&15)<<16 | (uint32(r)&15)<<0
		o3 = c.ofsr(p.As, int(p.From.Reg), 0, REGTMP, int(p.Scond), p)

	case 53: /* floating point load, int32 offset UGLY */
		o1 = c.omvl(p, &p.From, REGTMP)

		if o1 == 0 {
			break
		}
		r := int(p.From.Reg)
		if r == 0 {
			r = int(o.param)
		}
		o2 = c.oprrr(p, AADD, int(p.Scond)) | (REGTMP&15)<<12 | (REGTMP&15)<<16 | (uint32(r)&15)<<0
		o3 = c.ofsr(p.As, int(p.To.Reg), 0, (REGTMP&15), int(p.Scond), p) | 1<<20

	case 54: /* floating point arith */
		o1 = c.oprrr(p, p.As, int(p.Scond))

		rf := int(p.From.Reg)
		rt := int(p.To.Reg)
		r := int(p.Reg)
		if r == 0 {
			r = rt
			if p.As == AMOVF || p.As == AMOVD || p.As == AMOVFD || p.As == AMOVDF || p.As == ASQRTF || p.As == ASQRTD || p.As == AABSF || p.As == AABSD || p.As == ANEGF || p.As == ANEGD {
				r = 0
			}
		}

		o1 |= (uint32(rf)&15)<<0 | (uint32(r)&15)<<16 | (uint32(rt)&15)<<12

	case 56: /* move to FP[CS]R */
		o1 = ((uint32(p.Scond)&C_SCOND)^C_SCOND_XOR)<<28 | 0xe<<24 | 1<<8 | 1<<4

		o1 |= ((uint32(p.To.Reg)&1)+1)<<21 | (uint32(p.From.Reg)&15)<<12

	case 57: /* move from FP[CS]R */
		o1 = ((uint32(p.Scond)&C_SCOND)^C_SCOND_XOR)<<28 | 0xe<<24 | 1<<8 | 1<<4

		o1 |= ((uint32(p.From.Reg)&1)+1)<<21 | (uint32(p.To.Reg)&15)<<12 | 1<<20

	case 58: /* movbu R,R */
		o1 = c.oprrr(p, AAND, int(p.Scond))

		o1 |= uint32(immrot(0xff))
		rt := int(p.To.Reg)
		r := int(p.From.Reg)
		if p.To.Type == obj.TYPE_NONE {
			rt = 0
		}
		if r == 0 {
			r = rt
		}
		o1 |= (uint32(r)&15)<<16 | (uint32(rt)&15)<<12

	case 59: /* movw/bu R<<I(R),R -> ldr indexed */
		if p.From.Reg == 0 {
			if p.As != AMOVW {
				c.ctxt.Diag("byte MOV from shifter operand")
			}
			o1 = c.mov(p)
			break
		}

		if p.From.Offset&(1<<4) != 0 {
			c.ctxt.Diag("bad shift in LDR")
		}
		o1 = c.olrr(int(p.From.Offset), int(p.From.Reg), int(p.To.Reg), int(p.Scond))
		if p.As == AMOVBU {
			o1 |= 1 << 22
		}

	case 60: /* movb R(R),R -> ldrsb indexed */
		if p.From.Reg == 0 {
			c.ctxt.Diag("byte MOV from shifter operand")
			o1 = c.mov(p)
			break
		}

		if p.From.Offset&(^0xf) != 0 {
			c.ctxt.Diag("bad shift in LDRSB")
		}
		o1 = c.olhrr(int(p.From.Offset), int(p.From.Reg), int(p.To.Reg), int(p.Scond))
		o1 ^= 1<<5 | 1<<6

	case 61: /* movw/b/bu R,R<<[IR](R) -> str indexed */
		if p.To.Reg == 0 {
			c.ctxt.Diag("MOV to shifter operand")
		}
		o1 = c.osrr(int(p.From.Reg), int(p.To.Offset), int(p.To.Reg), int(p.Scond))
		if p.As == AMOVB || p.As == AMOVBS || p.As == AMOVBU {
			o1 |= 1 << 22
		}

		/* reloc ops */
	case 64: /* mov/movb/movbu R,addr */
		o1 = c.omvl(p, &p.To, REGTMP)

		if o1 == 0 {
			break
		}
		o2 = c.osr(p.As, int(p.From.Reg), 0, REGTMP, int(p.Scond))
		if o.flag&LPCREL != 0 {
			o3 = o2
			o2 = c.oprrr(p, AADD, int(p.Scond)) | REGTMP&15 | (REGPC&15)<<16 | (REGTMP&15)<<12
		}

	case 65: /* mov/movbu addr,R */
		o1 = c.omvl(p, &p.From, REGTMP)

		if o1 == 0 {
			break
		}
		o2 = c.olr(0, REGTMP, int(p.To.Reg), int(p.Scond))
		if p.As == AMOVBU || p.As == AMOVBS || p.As == AMOVB {
			o2 |= 1 << 22
		}
		if o.flag&LPCREL != 0 {
			o3 = o2
			o2 = c.oprrr(p, AADD, int(p.Scond)) | REGTMP&15 | (REGPC&15)<<16 | (REGTMP&15)<<12
		}

	case 101: /* movw tlsvar,R, local exec*/
		if p.Scond&C_SCOND != C_SCOND_NONE {
			c.ctxt.Diag("conditional tls")
		}
		o1 = c.omvl(p, &p.From, int(p.To.Reg))

	case 102: /* movw tlsvar,R, initial exec*/
		if p.Scond&C_SCOND != C_SCOND_NONE {
			c.ctxt.Diag("conditional tls")
		}
		o1 = c.omvl(p, &p.From, int(p.To.Reg))
		o2 = c.olrr(int(p.To.Reg)&15, (REGPC & 15), int(p.To.Reg), int(p.Scond))

	case 103: /* word tlsvar, local exec */
		if p.To.Sym == nil {
			c.ctxt.Diag("nil sym in tls %v", p)
		}
		if p.To.Offset != 0 {
			c.ctxt.Diag("offset against tls var in %v", p)
		}
		// This case happens with words generated in the PC stream as part of
		// the literal c.pool.
		rel := obj.Addrel(c.cursym)

		rel.Off = int32(c.pc)
		rel.Siz = 4
		rel.Sym = p.To.Sym
		rel.Type = objabi.R_TLS_LE
		o1 = 0

	case 104: /* word tlsvar, initial exec */
		if p.To.Sym == nil {
			c.ctxt.Diag("nil sym in tls %v", p)
		}
		if p.To.Offset != 0 {
			c.ctxt.Diag("offset against tls var in %v", p)
		}
		rel := obj.Addrel(c.cursym)
		rel.Off = int32(c.pc)
		rel.Siz = 4
		rel.Sym = p.To.Sym
		rel.Type = objabi.R_TLS_IE
		rel.Add = c.pc - p.Rel.Pc - 8 - int64(rel.Siz)

	case 68: /* floating point store -> ADDR */
		o1 = c.omvl(p, &p.To, REGTMP)

		if o1 == 0 {
			break
		}
		o2 = c.ofsr(p.As, int(p.From.Reg), 0, REGTMP, int(p.Scond), p)
		if o.flag&LPCREL != 0 {
			o3 = o2
			o2 = c.oprrr(p, AADD, int(p.Scond)) | REGTMP&15 | (REGPC&15)<<16 | (REGTMP&15)<<12
		}

	case 69: /* floating point load <- ADDR */
		o1 = c.omvl(p, &p.From, REGTMP)

		if o1 == 0 {
			break
		}
		o2 = c.ofsr(p.As, int(p.To.Reg), 0, (REGTMP&15), int(p.Scond), p) | 1<<20
		if o.flag&LPCREL != 0 {
			o3 = o2
			o2 = c.oprrr(p, AADD, int(p.Scond)) | REGTMP&15 | (REGPC&15)<<16 | (REGTMP&15)<<12
		}

		/* ArmV4 ops: */
	case 70: /* movh/movhu R,O(R) -> strh */
		c.aclass(&p.To)

		r := int(p.To.Reg)
		if r == 0 {
			r = int(o.param)
		}
		o1 = c.oshr(int(p.From.Reg), int32(c.instoffset), r, int(p.Scond))

	case 71: /* movb/movh/movhu O(R),R -> ldrsb/ldrsh/ldrh */
		c.aclass(&p.From)

		r := int(p.From.Reg)
		if r == 0 {
			r = int(o.param)
		}
		o1 = c.olhr(int32(c.instoffset), r, int(p.To.Reg), int(p.Scond))
		if p.As == AMOVB || p.As == AMOVBS {
			o1 ^= 1<<5 | 1<<6
		} else if p.As == AMOVH || p.As == AMOVHS {
			o1 ^= (1 << 6)
		}

	case 72: /* movh/movhu R,L(R) -> strh */
		o1 = c.omvl(p, &p.To, REGTMP)

		if o1 == 0 {
			break
		}
		r := int(p.To.Reg)
		if r == 0 {
			r = int(o.param)
		}
		o2 = c.oshrr(int(p.From.Reg), REGTMP&15, r, int(p.Scond))

	case 73: /* movb/movh/movhu L(R),R -> ldrsb/ldrsh/ldrh */
		o1 = c.omvl(p, &p.From, REGTMP)

		if o1 == 0 {
			break
		}
		r := int(p.From.Reg)
		if r == 0 {
			r = int(o.param)
		}
		o2 = c.olhrr(REGTMP&15, r, int(p.To.Reg), int(p.Scond))
		if p.As == AMOVB || p.As == AMOVBS {
			o2 ^= 1<<5 | 1<<6
		} else if p.As == AMOVH || p.As == AMOVHS {
			o2 ^= (1 << 6)
		}

	case 74: /* bx $I */
		c.ctxt.Diag("ABX $I")

	case 75: /* bx O(R) */
		c.aclass(&p.To)

		if c.instoffset != 0 {
			c.ctxt.Diag("non-zero offset in ABX")
		}

		/*
			o1 = 	c.oprrr(p, AADD, p->scond) | immrot(0) | ((REGPC&15)<<16) | ((REGLINK&15)<<12);	// mov PC, LR
			o2 = (((p->scond&C_SCOND) ^ C_SCOND_XOR)<<28) | (0x12fff<<8) | (1<<4) | ((p->to.reg&15) << 0);		// BX R
		*/
		// p->to.reg may be REGLINK
		o1 = c.oprrr(p, AADD, int(p.Scond))

		o1 |= uint32(immrot(uint32(c.instoffset)))
		o1 |= (uint32(p.To.Reg) & 15) << 16
		o1 |= (REGTMP & 15) << 12
		o2 = c.oprrr(p, AADD, int(p.Scond)) | uint32(immrot(0)) | (REGPC&15)<<16 | (REGLINK&15)<<12 // mov PC, LR
		o3 = ((uint32(p.Scond)&C_SCOND)^C_SCOND_XOR)<<28 | 0x12fff<<8 | 1<<4 | REGTMP&15            // BX Rtmp

	case 76: /* bx O(R) when returning from fn*/
		c.ctxt.Diag("ABXRET")

	case 77: /* ldrex oreg,reg */
		c.aclass(&p.From)

		if c.instoffset != 0 {
			c.ctxt.Diag("offset must be zero in LDREX")
		}
		o1 = 0x19<<20 | 0xf9f
		o1 |= (uint32(p.From.Reg) & 15) << 16
		o1 |= (uint32(p.To.Reg) & 15) << 12
		o1 |= ((uint32(p.Scond) & C_SCOND) ^ C_SCOND_XOR) << 28

	case 78: /* strex reg,oreg,reg */
		c.aclass(&p.From)

		if c.instoffset != 0 {
			c.ctxt.Diag("offset must be zero in STREX")
		}
		o1 = 0x18<<20 | 0xf90
		o1 |= (uint32(p.From.Reg) & 15) << 16
		o1 |= (uint32(p.Reg) & 15) << 0
		o1 |= (uint32(p.To.Reg) & 15) << 12
		o1 |= ((uint32(p.Scond) & C_SCOND) ^ C_SCOND_XOR) << 28

	case 80: /* fmov zfcon,freg */
		if p.As == AMOVD {
			o1 = 0xeeb00b00 // VMOV imm 64
			o2 = c.oprrr(p, ASUBD, int(p.Scond))
		} else {
			o1 = 0x0eb00a00 // VMOV imm 32
			o2 = c.oprrr(p, ASUBF, int(p.Scond))
		}

		v := int32(0x70) // 1.0
		r := (int(p.To.Reg) & 15) << 0

		// movf $1.0, r
		o1 |= ((uint32(p.Scond) & C_SCOND) ^ C_SCOND_XOR) << 28

		o1 |= (uint32(r) & 15) << 12
		o1 |= (uint32(v) & 0xf) << 0
		o1 |= (uint32(v) & 0xf0) << 12

		// subf r,r,r
		o2 |= (uint32(r)&15)<<0 | (uint32(r)&15)<<16 | (uint32(r)&15)<<12

	case 81: /* fmov sfcon,freg */
		o1 = 0x0eb00a00 // VMOV imm 32
		if p.As == AMOVD {
			o1 = 0xeeb00b00 // VMOV imm 64
		}
		o1 |= ((uint32(p.Scond) & C_SCOND) ^ C_SCOND_XOR) << 28
		o1 |= (uint32(p.To.Reg) & 15) << 12
		v := int32(c.chipfloat5(p.From.Val.(float64)))
		o1 |= (uint32(v) & 0xf) << 0
		o1 |= (uint32(v) & 0xf0) << 12

	case 82: /* fcmp freg,freg, */
		o1 = c.oprrr(p, p.As, int(p.Scond))

		o1 |= (uint32(p.Reg)&15)<<12 | (uint32(p.From.Reg)&15)<<0
		o2 = 0x0ef1fa10 // VMRS R15
		o2 |= ((uint32(p.Scond) & C_SCOND) ^ C_SCOND_XOR) << 28

	case 83: /* fcmp freg,, */
		o1 = c.oprrr(p, p.As, int(p.Scond))

		o1 |= (uint32(p.From.Reg)&15)<<12 | 1<<16
		o2 = 0x0ef1fa10 // VMRS R15
		o2 |= ((uint32(p.Scond) & C_SCOND) ^ C_SCOND_XOR) << 28

	case 84: /* movfw freg,freg - truncate float-to-fix */
		o1 = c.oprrr(p, p.As, int(p.Scond))

		o1 |= (uint32(p.From.Reg) & 15) << 0
		o1 |= (uint32(p.To.Reg) & 15) << 12

	case 85: /* movwf freg,freg - fix-to-float */
		o1 = c.oprrr(p, p.As, int(p.Scond))

		o1 |= (uint32(p.From.Reg) & 15) << 0
		o1 |= (uint32(p.To.Reg) & 15) << 12

		// macro for movfw freg,FTMP; movw FTMP,reg
	case 86: /* movfw freg,reg - truncate float-to-fix */
		o1 = c.oprrr(p, p.As, int(p.Scond))

		o1 |= (uint32(p.From.Reg) & 15) << 0
		o1 |= (FREGTMP & 15) << 12
		o2 = c.oprrr(p, -AMOVFW, int(p.Scond))
		o2 |= (FREGTMP & 15) << 16
		o2 |= (uint32(p.To.Reg) & 15) << 12

		// macro for movw reg,FTMP; movwf FTMP,freg
	case 87: /* movwf reg,freg - fix-to-float */
		o1 = c.oprrr(p, -AMOVWF, int(p.Scond))

		o1 |= (uint32(p.From.Reg) & 15) << 12
		o1 |= (FREGTMP & 15) << 16
		o2 = c.oprrr(p, p.As, int(p.Scond))
		o2 |= (FREGTMP & 15) << 0
		o2 |= (uint32(p.To.Reg) & 15) << 12

	case 88: /* movw reg,freg  */
		o1 = c.oprrr(p, -AMOVWF, int(p.Scond))

		o1 |= (uint32(p.From.Reg) & 15) << 12
		o1 |= (uint32(p.To.Reg) & 15) << 16

	case 89: /* movw freg,reg  */
		o1 = c.oprrr(p, -AMOVFW, int(p.Scond))

		o1 |= (uint32(p.From.Reg) & 15) << 16
		o1 |= (uint32(p.To.Reg) & 15) << 12

	case 91: /* ldrexd oreg,reg */
		c.aclass(&p.From)

		if c.instoffset != 0 {
			c.ctxt.Diag("offset must be zero in LDREX")
		}
		o1 = 0x1b<<20 | 0xf9f
		o1 |= (uint32(p.From.Reg) & 15) << 16
		o1 |= (uint32(p.To.Reg) & 15) << 12
		o1 |= ((uint32(p.Scond) & C_SCOND) ^ C_SCOND_XOR) << 28

	case 92: /* strexd reg,oreg,reg */
		c.aclass(&p.From)

		if c.instoffset != 0 {
			c.ctxt.Diag("offset must be zero in STREX")
		}
		o1 = 0x1a<<20 | 0xf90
		o1 |= (uint32(p.From.Reg) & 15) << 16
		o1 |= (uint32(p.Reg) & 15) << 0
		o1 |= (uint32(p.To.Reg) & 15) << 12
		o1 |= ((uint32(p.Scond) & C_SCOND) ^ C_SCOND_XOR) << 28

	case 93: /* movb/movh/movhu addr,R -> ldrsb/ldrsh/ldrh */
		o1 = c.omvl(p, &p.From, REGTMP)

		if o1 == 0 {
			break
		}
		o2 = c.olhr(0, REGTMP, int(p.To.Reg), int(p.Scond))
		if p.As == AMOVB || p.As == AMOVBS {
			o2 ^= 1<<5 | 1<<6
		} else if p.As == AMOVH || p.As == AMOVHS {
			o2 ^= (1 << 6)
		}
		if o.flag&LPCREL != 0 {
			o3 = o2
			o2 = c.oprrr(p, AADD, int(p.Scond)) | REGTMP&15 | (REGPC&15)<<16 | (REGTMP&15)<<12
		}

	case 94: /* movh/movhu R,addr -> strh */
		o1 = c.omvl(p, &p.To, REGTMP)

		if o1 == 0 {
			break
		}
		o2 = c.oshr(int(p.From.Reg), 0, REGTMP, int(p.Scond))
		if o.flag&LPCREL != 0 {
			o3 = o2
			o2 = c.oprrr(p, AADD, int(p.Scond)) | REGTMP&15 | (REGPC&15)<<16 | (REGTMP&15)<<12
		}

	case 95: /* PLD off(reg) */
		o1 = 0xf5d0f000

		o1 |= (uint32(p.From.Reg) & 15) << 16
		if p.From.Offset < 0 {
			o1 &^= (1 << 23)
			o1 |= uint32((-p.From.Offset) & 0xfff)
		} else {
			o1 |= uint32(p.From.Offset & 0xfff)
		}

	// This is supposed to be something that stops execution.
	// It's not supposed to be reached, ever, but if it is, we'd
	// like to be able to tell how we got there. Assemble as
	// 0xf7fabcfd which is guaranteed to raise undefined instruction
	// exception.
	case 96: /* UNDEF */
		o1 = 0xf7fabcfd

	case 97: /* CLZ Rm, Rd */
		o1 = c.oprrr(p, p.As, int(p.Scond))

		o1 |= (uint32(p.To.Reg) & 15) << 12
		o1 |= (uint32(p.From.Reg) & 15) << 0

	case 98: /* MULW{T,B} Rs, Rm, Rd */
		o1 = c.oprrr(p, p.As, int(p.Scond))

		o1 |= (uint32(p.To.Reg) & 15) << 16
		o1 |= (uint32(p.From.Reg) & 15) << 8
		o1 |= (uint32(p.Reg) & 15) << 0

	case 99: /* MULAW{T,B} Rs, Rm, Rn, Rd */
		o1 = c.oprrr(p, p.As, int(p.Scond))

		o1 |= (uint32(p.To.Reg) & 15) << 12
		o1 |= (uint32(p.From.Reg) & 15) << 8
		o1 |= (uint32(p.Reg) & 15) << 0
		o1 |= uint32((p.To.Offset & 15) << 16)

	// DATABUNDLE: BKPT $0x5be0, signify the start of NaCl data bundle;
	// DATABUNDLEEND: zero width alignment marker
	case 100:
		if p.As == ADATABUNDLE {
			o1 = 0xe125be70
		}

	case 105: /* divhw r,[r,]r */
		o1 = c.oprrr(p, p.As, int(p.Scond))
		rf := int(p.From.Reg)
		rt := int(p.To.Reg)
		r := int(p.Reg)
		if r == 0 {
			r = rt
		}
		o1 |= (uint32(rf)&15)<<8 | (uint32(r)&15)<<0 | (uint32(rt)&15)<<16
	}

	out[0] = o1
	out[1] = o2
	out[2] = o3
	out[3] = o4
	out[4] = o5
	out[5] = o6
	return
}

func (c *ctxt5) mov(p *obj.Prog) uint32 {
	c.aclass(&p.From)
	o1 := c.oprrr(p, p.As, int(p.Scond))
	o1 |= uint32(p.From.Offset)
	rt := int(p.To.Reg)
	if p.To.Type == obj.TYPE_NONE {
		rt = 0
	}
	r := int(p.Reg)
	if p.As == AMOVW || p.As == AMVN {
		r = 0
	} else if r == 0 {
		r = rt
	}
	o1 |= (uint32(r)&15)<<16 | (uint32(rt)&15)<<12
	return o1
}

func (c *ctxt5) oprrr(p *obj.Prog, a obj.As, sc int) uint32 {
	o := ((uint32(sc) & C_SCOND) ^ C_SCOND_XOR) << 28
	if sc&C_SBIT != 0 {
		o |= 1 << 20
	}
	if sc&(C_PBIT|C_WBIT) != 0 {
		c.ctxt.Diag(".nil/.W on dp instruction")
	}
	switch a {
	case ADIVHW:
		return o | 0x71<<20 | 0xf<<12 | 0x1<<4
	case ADIVUHW:
		return o | 0x73<<20 | 0xf<<12 | 0x1<<4
	case AMMUL:
		return o | 0x75<<20 | 0xf<<12 | 0x1<<4
	case AMULS:
		return o | 0x6<<20 | 0x9<<4
	case AMMULA:
		return o | 0x75<<20 | 0x1<<4
	case AMMULS:
		return o | 0x75<<20 | 0xd<<4
	case AMULU, AMUL:
		return o | 0x0<<21 | 0x9<<4
	case AMULA:
		return o | 0x1<<21 | 0x9<<4
	case AMULLU:
		return o | 0x4<<21 | 0x9<<4
	case AMULL:
		return o | 0x6<<21 | 0x9<<4
	case AMULALU:
		return o | 0x5<<21 | 0x9<<4
	case AMULAL:
		return o | 0x7<<21 | 0x9<<4
	case AAND:
		return o | 0x0<<21
	case AEOR:
		return o | 0x1<<21
	case ASUB:
		return o | 0x2<<21
	case ARSB:
		return o | 0x3<<21
	case AADD:
		return o | 0x4<<21
	case AADC:
		return o | 0x5<<21
	case ASBC:
		return o | 0x6<<21
	case ARSC:
		return o | 0x7<<21
	case ATST:
		return o | 0x8<<21 | 1<<20
	case ATEQ:
		return o | 0x9<<21 | 1<<20
	case ACMP:
		return o | 0xa<<21 | 1<<20
	case ACMN:
		return o | 0xb<<21 | 1<<20
	case AORR:
		return o | 0xc<<21

	case AMOVB, AMOVH, AMOVW:
		return o | 0xd<<21
	case ABIC:
		return o | 0xe<<21
	case AMVN:
		return o | 0xf<<21
	case ASLL:
		return o | 0xd<<21 | 0<<5
	case ASRL:
		return o | 0xd<<21 | 1<<5
	case ASRA:
		return o | 0xd<<21 | 2<<5
	case ASWI:
		return o | 0xf<<24

	case AADDD:
		return o | 0xe<<24 | 0x3<<20 | 0xb<<8 | 0<<4
	case AADDF:
		return o | 0xe<<24 | 0x3<<20 | 0xa<<8 | 0<<4
	case ASUBD:
		return o | 0xe<<24 | 0x3<<20 | 0xb<<8 | 4<<4
	case ASUBF:
		return o | 0xe<<24 | 0x3<<20 | 0xa<<8 | 4<<4
	case AMULD:
		return o | 0xe<<24 | 0x2<<20 | 0xb<<8 | 0<<4
	case AMULF:
		return o | 0xe<<24 | 0x2<<20 | 0xa<<8 | 0<<4
	case ADIVD:
		return o | 0xe<<24 | 0x8<<20 | 0xb<<8 | 0<<4
	case ADIVF:
		return o | 0xe<<24 | 0x8<<20 | 0xa<<8 | 0<<4
	case ASQRTD:
		return o | 0xe<<24 | 0xb<<20 | 1<<16 | 0xb<<8 | 0xc<<4
	case ASQRTF:
		return o | 0xe<<24 | 0xb<<20 | 1<<16 | 0xa<<8 | 0xc<<4
	case AABSD:
		return o | 0xe<<24 | 0xb<<20 | 0<<16 | 0xb<<8 | 0xc<<4
	case AABSF:
		return o | 0xe<<24 | 0xb<<20 | 0<<16 | 0xa<<8 | 0xc<<4
	case ANEGD:
		return o | 0xe<<24 | 0xb<<20 | 1<<16 | 0xb<<8 | 0x4<<4
	case ANEGF:
		return o | 0xe<<24 | 0xb<<20 | 1<<16 | 0xa<<8 | 0x4<<4
	case ACMPD:
		return o | 0xe<<24 | 0xb<<20 | 4<<16 | 0xb<<8 | 0xc<<4
	case ACMPF:
		return o | 0xe<<24 | 0xb<<20 | 4<<16 | 0xa<<8 | 0xc<<4

	case AMOVF:
		return o | 0xe<<24 | 0xb<<20 | 0<<16 | 0xa<<8 | 4<<4
	case AMOVD:
		return o | 0xe<<24 | 0xb<<20 | 0<<16 | 0xb<<8 | 4<<4

	case AMOVDF:
		return o | 0xe<<24 | 0xb<<20 | 7<<16 | 0xa<<8 | 0xc<<4 | 1<<8 // dtof
	case AMOVFD:
		return o | 0xe<<24 | 0xb<<20 | 7<<16 | 0xa<<8 | 0xc<<4 | 0<<8 // dtof

	case AMOVWF:
		if sc&C_UBIT == 0 {
			o |= 1 << 7 /* signed */
		}
		return o | 0xe<<24 | 0xb<<20 | 8<<16 | 0xa<<8 | 4<<4 | 0<<18 | 0<<8 // toint, double

	case AMOVWD:
		if sc&C_UBIT == 0 {
			o |= 1 << 7 /* signed */
		}
		return o | 0xe<<24 | 0xb<<20 | 8<<16 | 0xa<<8 | 4<<4 | 0<<18 | 1<<8 // toint, double

	case AMOVFW:
		if sc&C_UBIT == 0 {
			o |= 1 << 16 /* signed */
		}
		return o | 0xe<<24 | 0xb<<20 | 8<<16 | 0xa<<8 | 4<<4 | 1<<18 | 0<<8 | 1<<7 // toint, double, trunc

	case AMOVDW:
		if sc&C_UBIT == 0 {
			o |= 1 << 16 /* signed */
		}
		return o | 0xe<<24 | 0xb<<20 | 8<<16 | 0xa<<8 | 4<<4 | 1<<18 | 1<<8 | 1<<7 // toint, double, trunc

	case -AMOVWF: // copy WtoF
		return o | 0xe<<24 | 0x0<<20 | 0xb<<8 | 1<<4

	case -AMOVFW: // copy FtoW
		return o | 0xe<<24 | 0x1<<20 | 0xb<<8 | 1<<4

	case -ACMP: // cmp imm
		return o | 0x3<<24 | 0x5<<20

		// CLZ doesn't support .nil
	case ACLZ:
		return o&(0xf<<28) | 0x16f<<16 | 0xf1<<4

	case AREV:
		return o&(0xf<<28) | 0x6bf<<16 | 0xf3<<4

	case AREV16:
		return o&(0xf<<28) | 0x6bf<<16 | 0xfb<<4

	case AREVSH:
		return o&(0xf<<28) | 0x6ff<<16 | 0xfb<<4

	case ARBIT:
		return o&(0xf<<28) | 0x6ff<<16 | 0xf3<<4

	case AMULWT:
		return o&(0xf<<28) | 0x12<<20 | 0xe<<4

	case AMULWB:
		return o&(0xf<<28) | 0x12<<20 | 0xa<<4

	case AMULBB:
		return o&(0xf<<28) | 0x16<<20 | 0xf<<12 | 0x8<<4

	case AMULAWT:
		return o&(0xf<<28) | 0x12<<20 | 0xc<<4

	case AMULAWB:
		return o&(0xf<<28) | 0x12<<20 | 0x8<<4

	case AMULABB:
		return o&(0xf<<28) | 0x10<<20 | 0x8<<4

	case ABL: // BLX REG
		return o&(0xf<<28) | 0x12fff3<<4
	}

	c.ctxt.Diag("%v: bad rrr %d", p, a)
	return 0
}

func (c *ctxt5) opbra(p *obj.Prog, a obj.As, sc int) uint32 {
	if sc&(C_SBIT|C_PBIT|C_WBIT) != 0 {
		c.ctxt.Diag("%v: .nil/.nil/.W on bra instruction", p)
	}
	sc &= C_SCOND
	sc ^= C_SCOND_XOR
	if a == ABL || a == obj.ADUFFZERO || a == obj.ADUFFCOPY {
		return uint32(sc)<<28 | 0x5<<25 | 0x1<<24
	}
	if sc != 0xe {
		c.ctxt.Diag("%v: .COND on bcond instruction", p)
	}
	switch a {
	case ABEQ:
		return 0x0<<28 | 0x5<<25
	case ABNE:
		return 0x1<<28 | 0x5<<25
	case ABCS:
		return 0x2<<28 | 0x5<<25
	case ABHS:
		return 0x2<<28 | 0x5<<25
	case ABCC:
		return 0x3<<28 | 0x5<<25
	case ABLO:
		return 0x3<<28 | 0x5<<25
	case ABMI:
		return 0x4<<28 | 0x5<<25
	case ABPL:
		return 0x5<<28 | 0x5<<25
	case ABVS:
		return 0x6<<28 | 0x5<<25
	case ABVC:
		return 0x7<<28 | 0x5<<25
	case ABHI:
		return 0x8<<28 | 0x5<<25
	case ABLS:
		return 0x9<<28 | 0x5<<25
	case ABGE:
		return 0xa<<28 | 0x5<<25
	case ABLT:
		return 0xb<<28 | 0x5<<25
	case ABGT:
		return 0xc<<28 | 0x5<<25
	case ABLE:
		return 0xd<<28 | 0x5<<25
	case AB:
		return 0xe<<28 | 0x5<<25
	}

	c.ctxt.Diag("%v: bad bra %v", p, a)
	return 0
}

func (c *ctxt5) olr(v int32, b int, r int, sc int) uint32 {
	if sc&C_SBIT != 0 {
		c.ctxt.Diag(".nil on LDR/STR instruction")
	}
	o := ((uint32(sc) & C_SCOND) ^ C_SCOND_XOR) << 28
	if sc&C_PBIT == 0 {
		o |= 1 << 24
	}
	if sc&C_UBIT == 0 {
		o |= 1 << 23
	}
	if sc&C_WBIT != 0 {
		o |= 1 << 21
	}
	o |= 1<<26 | 1<<20
	if v < 0 {
		if sc&C_UBIT != 0 {
			c.ctxt.Diag(".U on neg offset")
		}
		v = -v
		o ^= 1 << 23
	}

	if v >= 1<<12 || v < 0 {
		c.ctxt.Diag("literal span too large: %d (R%d)\n%v", v, b, c.printp)
	}
	o |= uint32(v)
	o |= (uint32(b) & 15) << 16
	o |= (uint32(r) & 15) << 12
	return o
}

func (c *ctxt5) olhr(v int32, b int, r int, sc int) uint32 {
	if sc&C_SBIT != 0 {
		c.ctxt.Diag(".nil on LDRH/STRH instruction")
	}
	o := ((uint32(sc) & C_SCOND) ^ C_SCOND_XOR) << 28
	if sc&C_PBIT == 0 {
		o |= 1 << 24
	}
	if sc&C_WBIT != 0 {
		o |= 1 << 21
	}
	o |= 1<<23 | 1<<20 | 0xb<<4
	if v < 0 {
		v = -v
		o ^= 1 << 23
	}

	if v >= 1<<8 || v < 0 {
		c.ctxt.Diag("literal span too large: %d (R%d)\n%v", v, b, c.printp)
	}
	o |= uint32(v)&0xf | (uint32(v)>>4)<<8 | 1<<22
	o |= (uint32(b) & 15) << 16
	o |= (uint32(r) & 15) << 12
	return o
}

func (c *ctxt5) osr(a obj.As, r int, v int32, b int, sc int) uint32 {
	o := c.olr(v, b, r, sc) ^ (1 << 20)
	if a != AMOVW {
		o |= 1 << 22
	}
	return o
}

func (c *ctxt5) oshr(r int, v int32, b int, sc int) uint32 {
	o := c.olhr(v, b, r, sc) ^ (1 << 20)
	return o
}

func (c *ctxt5) osrr(r int, i int, b int, sc int) uint32 {
	return c.olr(int32(i), b, r, sc) ^ (1<<25 | 1<<20)
}

func (c *ctxt5) oshrr(r int, i int, b int, sc int) uint32 {
	return c.olhr(int32(i), b, r, sc) ^ (1<<22 | 1<<20)
}

func (c *ctxt5) olrr(i int, b int, r int, sc int) uint32 {
	return c.olr(int32(i), b, r, sc) ^ (1 << 25)
}

func (c *ctxt5) olhrr(i int, b int, r int, sc int) uint32 {
	return c.olhr(int32(i), b, r, sc) ^ (1 << 22)
}

func (c *ctxt5) ofsr(a obj.As, r int, v int32, b int, sc int, p *obj.Prog) uint32 {
	if sc&C_SBIT != 0 {
		c.ctxt.Diag(".nil on FLDR/FSTR instruction: %v", p)
	}
	o := ((uint32(sc) & C_SCOND) ^ C_SCOND_XOR) << 28
	if sc&C_PBIT == 0 {
		o |= 1 << 24
	}
	if sc&C_WBIT != 0 {
		o |= 1 << 21
	}
	o |= 6<<25 | 1<<24 | 1<<23 | 10<<8
	if v < 0 {
		v = -v
		o ^= 1 << 23
	}

	if v&3 != 0 {
		c.ctxt.Diag("odd offset for floating point op: %d\n%v", v, p)
	} else if v >= 1<<10 || v < 0 {
		c.ctxt.Diag("literal span too large: %d\n%v", v, p)
	}
	o |= (uint32(v) >> 2) & 0xFF
	o |= (uint32(b) & 15) << 16
	o |= (uint32(r) & 15) << 12

	switch a {
	default:
		c.ctxt.Diag("bad fst %v", a)
		fallthrough

	case AMOVD:
		o |= 1 << 8
		fallthrough

	case AMOVF:
		break
	}

	return o
}

// MOVW $"lower 16-bit", Reg
func (c *ctxt5) omvs(p *obj.Prog, a *obj.Addr, dr int) uint32 {
	var o1 uint32
	o1 = ((uint32(p.Scond) & C_SCOND) ^ C_SCOND_XOR) << 28
	o1 |= 0x30 << 20
	o1 |= (uint32(dr) & 15) << 12
	o1 |= uint32(a.Offset) & 0x0fff
	o1 |= (uint32(a.Offset) & 0xf000) << 4
	return o1
}

func (c *ctxt5) omvl(p *obj.Prog, a *obj.Addr, dr int) uint32 {
	var o1 uint32
	if p.Pcond == nil {
		c.aclass(a)
		v := immrot(^uint32(c.instoffset))
		if v == 0 {
			c.ctxt.Diag("%v: missing literal", p)
			return 0
		}

		o1 = c.oprrr(p, AMVN, int(p.Scond)&C_SCOND)
		o1 |= uint32(v)
		o1 |= (uint32(dr) & 15) << 12
	} else {
		v := int32(p.Pcond.Pc - p.Pc - 8)
		o1 = c.olr(v, REGPC, dr, int(p.Scond)&C_SCOND)
	}

	return o1
}

func (c *ctxt5) chipzero5(e float64) int {
	// We use GOARM=7 to gate the use of VFPv3 vmov (imm) instructions.
	if objabi.GOARM < 7 || e != 0 {
		return -1
	}
	return 0
}

func (c *ctxt5) chipfloat5(e float64) int {
	// We use GOARM=7 to gate the use of VFPv3 vmov (imm) instructions.
	if objabi.GOARM < 7 {
		return -1
	}

	ei := math.Float64bits(e)
	l := uint32(ei)
	h := uint32(ei >> 32)

	if l != 0 || h&0xffff != 0 {
		return -1
	}
	h1 := h & 0x7fc00000
	if h1 != 0x40000000 && h1 != 0x3fc00000 {
		return -1
	}
	n := 0

	// sign bit (a)
	if h&0x80000000 != 0 {
		n |= 1 << 7
	}

	// exp sign bit (b)
	if h1 == 0x3fc00000 {
		n |= 1 << 6
	}

	// rest of exp and mantissa (cd-efgh)
	n |= int((h >> 16) & 0x3f)

	//print("match %.8lux %.8lux %d\n", l, h, n);
	return n
}

func nocache(p *obj.Prog) {
	p.Optab = 0
	p.From.Class = 0
	if p.From3 != nil {
		p.From3.Class = 0
	}
	p.To.Class = 0
}
