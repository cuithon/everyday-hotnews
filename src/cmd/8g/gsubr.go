// Derived from Inferno utils/8c/txt.c
// http://code.google.com/p/inferno-os/source/browse/utils/8c/txt.c
//
//	Copyright © 1994-1999 Lucent Technologies Inc.  All rights reserved.
//	Portions Copyright © 1995-1997 C H Forsyth (forsyth@terzarima.net)
//	Portions Copyright © 1997-1999 Vita Nuova Limited
//	Portions Copyright © 2000-2007 Vita Nuova Holdings Limited (www.vitanuova.com)
//	Portions Copyright © 2004,2006 Bruce Ellis
//	Portions Copyright © 2005-2007 C H Forsyth (forsyth@terzarima.net)
//	Revisions Copyright © 2000-2007 Lucent Technologies Inc. and others
//	Portions Copyright © 2009 The Go Authors.  All rights reserved.
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

package main

import (
	"cmd/internal/gc"
	"cmd/internal/obj"
	"cmd/internal/obj/x86"
	"fmt"
)

// TODO(rsc): Can make this bigger if we move
// the text segment up higher in 8l for all GOOS.
// At the same time, can raise StackBig in ../../runtime/stack.h.
var unmappedzero uint32 = 4096

/*
 * return Axxx for Oxxx on type t.
 */
func optoas(op int, t *gc.Type) int {
	if t == nil {
		gc.Fatal("optoas: t is nil")
	}

	a := obj.AXXX
	switch uint32(op)<<16 | uint32(gc.Simtype[t.Etype]) {
	default:
		gc.Fatal("optoas: no entry %v-%v", gc.Oconv(int(op), 0), gc.Tconv(t, 0))

	case gc.OADDR<<16 | gc.TPTR32:
		a = x86.ALEAL

	case gc.OEQ<<16 | gc.TBOOL,
		gc.OEQ<<16 | gc.TINT8,
		gc.OEQ<<16 | gc.TUINT8,
		gc.OEQ<<16 | gc.TINT16,
		gc.OEQ<<16 | gc.TUINT16,
		gc.OEQ<<16 | gc.TINT32,
		gc.OEQ<<16 | gc.TUINT32,
		gc.OEQ<<16 | gc.TINT64,
		gc.OEQ<<16 | gc.TUINT64,
		gc.OEQ<<16 | gc.TPTR32,
		gc.OEQ<<16 | gc.TPTR64,
		gc.OEQ<<16 | gc.TFLOAT32,
		gc.OEQ<<16 | gc.TFLOAT64:
		a = x86.AJEQ

	case gc.ONE<<16 | gc.TBOOL,
		gc.ONE<<16 | gc.TINT8,
		gc.ONE<<16 | gc.TUINT8,
		gc.ONE<<16 | gc.TINT16,
		gc.ONE<<16 | gc.TUINT16,
		gc.ONE<<16 | gc.TINT32,
		gc.ONE<<16 | gc.TUINT32,
		gc.ONE<<16 | gc.TINT64,
		gc.ONE<<16 | gc.TUINT64,
		gc.ONE<<16 | gc.TPTR32,
		gc.ONE<<16 | gc.TPTR64,
		gc.ONE<<16 | gc.TFLOAT32,
		gc.ONE<<16 | gc.TFLOAT64:
		a = x86.AJNE

	case gc.OLT<<16 | gc.TINT8,
		gc.OLT<<16 | gc.TINT16,
		gc.OLT<<16 | gc.TINT32,
		gc.OLT<<16 | gc.TINT64:
		a = x86.AJLT

	case gc.OLT<<16 | gc.TUINT8,
		gc.OLT<<16 | gc.TUINT16,
		gc.OLT<<16 | gc.TUINT32,
		gc.OLT<<16 | gc.TUINT64:
		a = x86.AJCS

	case gc.OLE<<16 | gc.TINT8,
		gc.OLE<<16 | gc.TINT16,
		gc.OLE<<16 | gc.TINT32,
		gc.OLE<<16 | gc.TINT64:
		a = x86.AJLE

	case gc.OLE<<16 | gc.TUINT8,
		gc.OLE<<16 | gc.TUINT16,
		gc.OLE<<16 | gc.TUINT32,
		gc.OLE<<16 | gc.TUINT64:
		a = x86.AJLS

	case gc.OGT<<16 | gc.TINT8,
		gc.OGT<<16 | gc.TINT16,
		gc.OGT<<16 | gc.TINT32,
		gc.OGT<<16 | gc.TINT64:
		a = x86.AJGT

	case gc.OGT<<16 | gc.TUINT8,
		gc.OGT<<16 | gc.TUINT16,
		gc.OGT<<16 | gc.TUINT32,
		gc.OGT<<16 | gc.TUINT64,
		gc.OLT<<16 | gc.TFLOAT32,
		gc.OLT<<16 | gc.TFLOAT64:
		a = x86.AJHI

	case gc.OGE<<16 | gc.TINT8,
		gc.OGE<<16 | gc.TINT16,
		gc.OGE<<16 | gc.TINT32,
		gc.OGE<<16 | gc.TINT64:
		a = x86.AJGE

	case gc.OGE<<16 | gc.TUINT8,
		gc.OGE<<16 | gc.TUINT16,
		gc.OGE<<16 | gc.TUINT32,
		gc.OGE<<16 | gc.TUINT64,
		gc.OLE<<16 | gc.TFLOAT32,
		gc.OLE<<16 | gc.TFLOAT64:
		a = x86.AJCC

	case gc.OCMP<<16 | gc.TBOOL,
		gc.OCMP<<16 | gc.TINT8,
		gc.OCMP<<16 | gc.TUINT8:
		a = x86.ACMPB

	case gc.OCMP<<16 | gc.TINT16,
		gc.OCMP<<16 | gc.TUINT16:
		a = x86.ACMPW

	case gc.OCMP<<16 | gc.TINT32,
		gc.OCMP<<16 | gc.TUINT32,
		gc.OCMP<<16 | gc.TPTR32:
		a = x86.ACMPL

	case gc.OAS<<16 | gc.TBOOL,
		gc.OAS<<16 | gc.TINT8,
		gc.OAS<<16 | gc.TUINT8:
		a = x86.AMOVB

	case gc.OAS<<16 | gc.TINT16,
		gc.OAS<<16 | gc.TUINT16:
		a = x86.AMOVW

	case gc.OAS<<16 | gc.TINT32,
		gc.OAS<<16 | gc.TUINT32,
		gc.OAS<<16 | gc.TPTR32:
		a = x86.AMOVL

	case gc.OAS<<16 | gc.TFLOAT32:
		a = x86.AMOVSS

	case gc.OAS<<16 | gc.TFLOAT64:
		a = x86.AMOVSD

	case gc.OADD<<16 | gc.TINT8,
		gc.OADD<<16 | gc.TUINT8:
		a = x86.AADDB

	case gc.OADD<<16 | gc.TINT16,
		gc.OADD<<16 | gc.TUINT16:
		a = x86.AADDW

	case gc.OADD<<16 | gc.TINT32,
		gc.OADD<<16 | gc.TUINT32,
		gc.OADD<<16 | gc.TPTR32:
		a = x86.AADDL

	case gc.OSUB<<16 | gc.TINT8,
		gc.OSUB<<16 | gc.TUINT8:
		a = x86.ASUBB

	case gc.OSUB<<16 | gc.TINT16,
		gc.OSUB<<16 | gc.TUINT16:
		a = x86.ASUBW

	case gc.OSUB<<16 | gc.TINT32,
		gc.OSUB<<16 | gc.TUINT32,
		gc.OSUB<<16 | gc.TPTR32:
		a = x86.ASUBL

	case gc.OINC<<16 | gc.TINT8,
		gc.OINC<<16 | gc.TUINT8:
		a = x86.AINCB

	case gc.OINC<<16 | gc.TINT16,
		gc.OINC<<16 | gc.TUINT16:
		a = x86.AINCW

	case gc.OINC<<16 | gc.TINT32,
		gc.OINC<<16 | gc.TUINT32,
		gc.OINC<<16 | gc.TPTR32:
		a = x86.AINCL

	case gc.ODEC<<16 | gc.TINT8,
		gc.ODEC<<16 | gc.TUINT8:
		a = x86.ADECB

	case gc.ODEC<<16 | gc.TINT16,
		gc.ODEC<<16 | gc.TUINT16:
		a = x86.ADECW

	case gc.ODEC<<16 | gc.TINT32,
		gc.ODEC<<16 | gc.TUINT32,
		gc.ODEC<<16 | gc.TPTR32:
		a = x86.ADECL

	case gc.OCOM<<16 | gc.TINT8,
		gc.OCOM<<16 | gc.TUINT8:
		a = x86.ANOTB

	case gc.OCOM<<16 | gc.TINT16,
		gc.OCOM<<16 | gc.TUINT16:
		a = x86.ANOTW

	case gc.OCOM<<16 | gc.TINT32,
		gc.OCOM<<16 | gc.TUINT32,
		gc.OCOM<<16 | gc.TPTR32:
		a = x86.ANOTL

	case gc.OMINUS<<16 | gc.TINT8,
		gc.OMINUS<<16 | gc.TUINT8:
		a = x86.ANEGB

	case gc.OMINUS<<16 | gc.TINT16,
		gc.OMINUS<<16 | gc.TUINT16:
		a = x86.ANEGW

	case gc.OMINUS<<16 | gc.TINT32,
		gc.OMINUS<<16 | gc.TUINT32,
		gc.OMINUS<<16 | gc.TPTR32:
		a = x86.ANEGL

	case gc.OAND<<16 | gc.TINT8,
		gc.OAND<<16 | gc.TUINT8:
		a = x86.AANDB

	case gc.OAND<<16 | gc.TINT16,
		gc.OAND<<16 | gc.TUINT16:
		a = x86.AANDW

	case gc.OAND<<16 | gc.TINT32,
		gc.OAND<<16 | gc.TUINT32,
		gc.OAND<<16 | gc.TPTR32:
		a = x86.AANDL

	case gc.OOR<<16 | gc.TINT8,
		gc.OOR<<16 | gc.TUINT8:
		a = x86.AORB

	case gc.OOR<<16 | gc.TINT16,
		gc.OOR<<16 | gc.TUINT16:
		a = x86.AORW

	case gc.OOR<<16 | gc.TINT32,
		gc.OOR<<16 | gc.TUINT32,
		gc.OOR<<16 | gc.TPTR32:
		a = x86.AORL

	case gc.OXOR<<16 | gc.TINT8,
		gc.OXOR<<16 | gc.TUINT8:
		a = x86.AXORB

	case gc.OXOR<<16 | gc.TINT16,
		gc.OXOR<<16 | gc.TUINT16:
		a = x86.AXORW

	case gc.OXOR<<16 | gc.TINT32,
		gc.OXOR<<16 | gc.TUINT32,
		gc.OXOR<<16 | gc.TPTR32:
		a = x86.AXORL

	case gc.OLROT<<16 | gc.TINT8,
		gc.OLROT<<16 | gc.TUINT8:
		a = x86.AROLB

	case gc.OLROT<<16 | gc.TINT16,
		gc.OLROT<<16 | gc.TUINT16:
		a = x86.AROLW

	case gc.OLROT<<16 | gc.TINT32,
		gc.OLROT<<16 | gc.TUINT32,
		gc.OLROT<<16 | gc.TPTR32:
		a = x86.AROLL

	case gc.OLSH<<16 | gc.TINT8,
		gc.OLSH<<16 | gc.TUINT8:
		a = x86.ASHLB

	case gc.OLSH<<16 | gc.TINT16,
		gc.OLSH<<16 | gc.TUINT16:
		a = x86.ASHLW

	case gc.OLSH<<16 | gc.TINT32,
		gc.OLSH<<16 | gc.TUINT32,
		gc.OLSH<<16 | gc.TPTR32:
		a = x86.ASHLL

	case gc.ORSH<<16 | gc.TUINT8:
		a = x86.ASHRB

	case gc.ORSH<<16 | gc.TUINT16:
		a = x86.ASHRW

	case gc.ORSH<<16 | gc.TUINT32,
		gc.ORSH<<16 | gc.TPTR32:
		a = x86.ASHRL

	case gc.ORSH<<16 | gc.TINT8:
		a = x86.ASARB

	case gc.ORSH<<16 | gc.TINT16:
		a = x86.ASARW

	case gc.ORSH<<16 | gc.TINT32:
		a = x86.ASARL

	case gc.OHMUL<<16 | gc.TINT8,
		gc.OMUL<<16 | gc.TINT8,
		gc.OMUL<<16 | gc.TUINT8:
		a = x86.AIMULB

	case gc.OHMUL<<16 | gc.TINT16,
		gc.OMUL<<16 | gc.TINT16,
		gc.OMUL<<16 | gc.TUINT16:
		a = x86.AIMULW

	case gc.OHMUL<<16 | gc.TINT32,
		gc.OMUL<<16 | gc.TINT32,
		gc.OMUL<<16 | gc.TUINT32,
		gc.OMUL<<16 | gc.TPTR32:
		a = x86.AIMULL

	case gc.OHMUL<<16 | gc.TUINT8:
		a = x86.AMULB

	case gc.OHMUL<<16 | gc.TUINT16:
		a = x86.AMULW

	case gc.OHMUL<<16 | gc.TUINT32,
		gc.OHMUL<<16 | gc.TPTR32:
		a = x86.AMULL

	case gc.ODIV<<16 | gc.TINT8,
		gc.OMOD<<16 | gc.TINT8:
		a = x86.AIDIVB

	case gc.ODIV<<16 | gc.TUINT8,
		gc.OMOD<<16 | gc.TUINT8:
		a = x86.ADIVB

	case gc.ODIV<<16 | gc.TINT16,
		gc.OMOD<<16 | gc.TINT16:
		a = x86.AIDIVW

	case gc.ODIV<<16 | gc.TUINT16,
		gc.OMOD<<16 | gc.TUINT16:
		a = x86.ADIVW

	case gc.ODIV<<16 | gc.TINT32,
		gc.OMOD<<16 | gc.TINT32:
		a = x86.AIDIVL

	case gc.ODIV<<16 | gc.TUINT32,
		gc.ODIV<<16 | gc.TPTR32,
		gc.OMOD<<16 | gc.TUINT32,
		gc.OMOD<<16 | gc.TPTR32:
		a = x86.ADIVL

	case gc.OEXTEND<<16 | gc.TINT16:
		a = x86.ACWD

	case gc.OEXTEND<<16 | gc.TINT32:
		a = x86.ACDQ
	}

	return a
}

func foptoas(op int, t *gc.Type, flg int) int {
	a := obj.AXXX
	et := int(gc.Simtype[t.Etype])

	if !gc.Thearch.Use387 {
		switch uint32(op)<<16 | uint32(et) {
		default:
			gc.Fatal("foptoas-sse: no entry %v-%v", gc.Oconv(int(op), 0), gc.Tconv(t, 0))

		case gc.OCMP<<16 | gc.TFLOAT32:
			a = x86.AUCOMISS

		case gc.OCMP<<16 | gc.TFLOAT64:
			a = x86.AUCOMISD

		case gc.OAS<<16 | gc.TFLOAT32:
			a = x86.AMOVSS

		case gc.OAS<<16 | gc.TFLOAT64:
			a = x86.AMOVSD

		case gc.OADD<<16 | gc.TFLOAT32:
			a = x86.AADDSS

		case gc.OADD<<16 | gc.TFLOAT64:
			a = x86.AADDSD

		case gc.OSUB<<16 | gc.TFLOAT32:
			a = x86.ASUBSS

		case gc.OSUB<<16 | gc.TFLOAT64:
			a = x86.ASUBSD

		case gc.OMUL<<16 | gc.TFLOAT32:
			a = x86.AMULSS

		case gc.OMUL<<16 | gc.TFLOAT64:
			a = x86.AMULSD

		case gc.ODIV<<16 | gc.TFLOAT32:
			a = x86.ADIVSS

		case gc.ODIV<<16 | gc.TFLOAT64:
			a = x86.ADIVSD
		}

		return a
	}

	// If we need Fpop, it means we're working on
	// two different floating-point registers, not memory.
	// There the instruction only has a float64 form.
	if flg&Fpop != 0 {
		et = gc.TFLOAT64
	}

	// clear Frev if unneeded
	switch op {
	case gc.OADD,
		gc.OMUL:
		flg &^= Frev
	}

	switch uint32(op)<<16 | (uint32(et)<<8 | uint32(flg)) {
	case gc.OADD<<16 | (gc.TFLOAT32<<8 | 0):
		return x86.AFADDF

	case gc.OADD<<16 | (gc.TFLOAT64<<8 | 0):
		return x86.AFADDD

	case gc.OADD<<16 | (gc.TFLOAT64<<8 | Fpop):
		return x86.AFADDDP

	case gc.OSUB<<16 | (gc.TFLOAT32<<8 | 0):
		return x86.AFSUBF

	case gc.OSUB<<16 | (gc.TFLOAT32<<8 | Frev):
		return x86.AFSUBRF

	case gc.OSUB<<16 | (gc.TFLOAT64<<8 | 0):
		return x86.AFSUBD

	case gc.OSUB<<16 | (gc.TFLOAT64<<8 | Frev):
		return x86.AFSUBRD

	case gc.OSUB<<16 | (gc.TFLOAT64<<8 | Fpop):
		return x86.AFSUBDP

	case gc.OSUB<<16 | (gc.TFLOAT64<<8 | (Fpop | Frev)):
		return x86.AFSUBRDP

	case gc.OMUL<<16 | (gc.TFLOAT32<<8 | 0):
		return x86.AFMULF

	case gc.OMUL<<16 | (gc.TFLOAT64<<8 | 0):
		return x86.AFMULD

	case gc.OMUL<<16 | (gc.TFLOAT64<<8 | Fpop):
		return x86.AFMULDP

	case gc.ODIV<<16 | (gc.TFLOAT32<<8 | 0):
		return x86.AFDIVF

	case gc.ODIV<<16 | (gc.TFLOAT32<<8 | Frev):
		return x86.AFDIVRF

	case gc.ODIV<<16 | (gc.TFLOAT64<<8 | 0):
		return x86.AFDIVD

	case gc.ODIV<<16 | (gc.TFLOAT64<<8 | Frev):
		return x86.AFDIVRD

	case gc.ODIV<<16 | (gc.TFLOAT64<<8 | Fpop):
		return x86.AFDIVDP

	case gc.ODIV<<16 | (gc.TFLOAT64<<8 | (Fpop | Frev)):
		return x86.AFDIVRDP

	case gc.OCMP<<16 | (gc.TFLOAT32<<8 | 0):
		return x86.AFCOMF

	case gc.OCMP<<16 | (gc.TFLOAT32<<8 | Fpop):
		return x86.AFCOMFP

	case gc.OCMP<<16 | (gc.TFLOAT64<<8 | 0):
		return x86.AFCOMD

	case gc.OCMP<<16 | (gc.TFLOAT64<<8 | Fpop):
		return x86.AFCOMDP

	case gc.OCMP<<16 | (gc.TFLOAT64<<8 | Fpop2):
		return x86.AFCOMDPP

	case gc.OMINUS<<16 | (gc.TFLOAT32<<8 | 0):
		return x86.AFCHS

	case gc.OMINUS<<16 | (gc.TFLOAT64<<8 | 0):
		return x86.AFCHS
	}

	gc.Fatal("foptoas %v %v %#x", gc.Oconv(int(op), 0), gc.Tconv(t, 0), flg)
	return 0
}

var resvd = []int{
	//	REG_DI,	// for movstring
	//	REG_SI,	// for movstring

	x86.REG_AX, // for divide
	x86.REG_CX, // for shift
	x86.REG_DX, // for divide
	x86.REG_SP, // for stack
}

/*
 * generate
 *	as $c, reg
 */
func gconreg(as int, c int64, reg int) {
	var n1 gc.Node
	var n2 gc.Node

	gc.Nodconst(&n1, gc.Types[gc.TINT64], c)
	gc.Nodreg(&n2, gc.Types[gc.TINT64], reg)
	gins(as, &n1, &n2)
}

/*
 * generate
 *	as $c, n
 */
func ginscon(as int, c int64, n2 *gc.Node) {
	var n1 gc.Node
	gc.Nodconst(&n1, gc.Types[gc.TINT32], c)
	gins(as, &n1, n2)
}

/*
 * swap node contents
 */
func nswap(a *gc.Node, b *gc.Node) {
	t := *a
	*a = *b
	*b = t
}

/*
 * return constant i node.
 * overwritten by next call, but useful in calls to gins.
 */

var ncon_n gc.Node

func ncon(i uint32) *gc.Node {
	if ncon_n.Type == nil {
		gc.Nodconst(&ncon_n, gc.Types[gc.TUINT32], 0)
	}
	gc.Mpmovecfix(ncon_n.Val.U.Xval, int64(i))
	return &ncon_n
}

var sclean [10]gc.Node

var nsclean int

/*
 * n is a 64-bit value.  fill in lo and hi to refer to its 32-bit halves.
 */
func split64(n *gc.Node, lo *gc.Node, hi *gc.Node) {
	if !gc.Is64(n.Type) {
		gc.Fatal("split64 %v", gc.Tconv(n.Type, 0))
	}

	if nsclean >= len(sclean) {
		gc.Fatal("split64 clean")
	}
	sclean[nsclean].Op = gc.OEMPTY
	nsclean++
	switch n.Op {
	default:
		switch n.Op {
		default:
			var n1 gc.Node
			if !dotaddable(n, &n1) {
				gc.Igen(n, &n1, nil)
				sclean[nsclean-1] = n1
			}

			n = &n1

		case gc.ONAME:
			if n.Class == gc.PPARAMREF {
				var n1 gc.Node
				gc.Cgen(n.Heapaddr, &n1)
				sclean[nsclean-1] = n1
				n = &n1
			}

			// nothing
		case gc.OINDREG:
			break
		}

		*lo = *n
		*hi = *n
		lo.Type = gc.Types[gc.TUINT32]
		if n.Type.Etype == gc.TINT64 {
			hi.Type = gc.Types[gc.TINT32]
		} else {
			hi.Type = gc.Types[gc.TUINT32]
		}
		hi.Xoffset += 4

	case gc.OLITERAL:
		var n1 gc.Node
		gc.Convconst(&n1, n.Type, &n.Val)
		i := gc.Mpgetfix(n1.Val.U.Xval)
		gc.Nodconst(lo, gc.Types[gc.TUINT32], int64(uint32(i)))
		i >>= 32
		if n.Type.Etype == gc.TINT64 {
			gc.Nodconst(hi, gc.Types[gc.TINT32], int64(int32(i)))
		} else {
			gc.Nodconst(hi, gc.Types[gc.TUINT32], int64(uint32(i)))
		}
	}
}

func splitclean() {
	if nsclean <= 0 {
		gc.Fatal("splitclean")
	}
	nsclean--
	if sclean[nsclean].Op != gc.OEMPTY {
		gc.Regfree(&sclean[nsclean])
	}
}

/*
 * set up nodes representing fp constants
 */
var zerof gc.Node

var two64f gc.Node

var two63f gc.Node

var bignodes_did int

func bignodes() {
	if bignodes_did != 0 {
		return
	}
	bignodes_did = 1

	two64f = *ncon(0)
	two64f.Type = gc.Types[gc.TFLOAT64]
	two64f.Val.Ctype = gc.CTFLT
	two64f.Val.U.Fval = new(gc.Mpflt)
	gc.Mpmovecflt(two64f.Val.U.Fval, 18446744073709551616.)

	two63f = two64f
	two63f.Val.U.Fval = new(gc.Mpflt)
	gc.Mpmovecflt(two63f.Val.U.Fval, 9223372036854775808.)

	zerof = two64f
	zerof.Val.U.Fval = new(gc.Mpflt)
	gc.Mpmovecflt(zerof.Val.U.Fval, 0)
}

func memname(n *gc.Node, t *gc.Type) {
	gc.Tempname(n, t)
	n.Sym = gc.Lookup("." + n.Sym.Name[1:]) // keep optimizer from registerizing
	n.Orig.Sym = n.Sym
}

func gmove(f *gc.Node, t *gc.Node) {
	if gc.Debug['M'] != 0 {
		fmt.Printf("gmove %v -> %v\n", gc.Nconv(f, 0), gc.Nconv(t, 0))
	}

	ft := gc.Simsimtype(f.Type)
	tt := gc.Simsimtype(t.Type)
	cvt := t.Type

	if gc.Iscomplex[ft] || gc.Iscomplex[tt] {
		gc.Complexmove(f, t)
		return
	}

	if gc.Isfloat[ft] || gc.Isfloat[tt] {
		floatmove(f, t)
		return
	}

	// cannot have two integer memory operands;
	// except 64-bit, which always copies via registers anyway.
	var r1 gc.Node
	var a int
	if gc.Isint[ft] && gc.Isint[tt] && !gc.Is64(f.Type) && !gc.Is64(t.Type) && gc.Ismem(f) && gc.Ismem(t) {
		goto hard
	}

	// convert constant to desired type
	if f.Op == gc.OLITERAL {
		var con gc.Node
		gc.Convconst(&con, t.Type, &f.Val)
		f = &con
		ft = gc.Simsimtype(con.Type)
	}

	// value -> value copy, only one memory operand.
	// figure out the instruction to use.
	// break out of switch for one-instruction gins.
	// goto rdst for "destination must be register".
	// goto hard for "convert to cvt type first".
	// otherwise handle and return.

	switch uint32(ft)<<16 | uint32(tt) {
	default:
		// should not happen
		gc.Fatal("gmove %v -> %v", gc.Nconv(f, 0), gc.Nconv(t, 0))
		return

		/*
		 * integer copy and truncate
		 */
	case gc.TINT8<<16 | gc.TINT8, // same size
		gc.TINT8<<16 | gc.TUINT8,
		gc.TUINT8<<16 | gc.TINT8,
		gc.TUINT8<<16 | gc.TUINT8:
		a = x86.AMOVB

	case gc.TINT16<<16 | gc.TINT8, // truncate
		gc.TUINT16<<16 | gc.TINT8,
		gc.TINT32<<16 | gc.TINT8,
		gc.TUINT32<<16 | gc.TINT8,
		gc.TINT16<<16 | gc.TUINT8,
		gc.TUINT16<<16 | gc.TUINT8,
		gc.TINT32<<16 | gc.TUINT8,
		gc.TUINT32<<16 | gc.TUINT8:
		a = x86.AMOVB

		goto rsrc

	case gc.TINT64<<16 | gc.TINT8, // truncate low word
		gc.TUINT64<<16 | gc.TINT8,
		gc.TINT64<<16 | gc.TUINT8,
		gc.TUINT64<<16 | gc.TUINT8:
		var flo gc.Node
		var fhi gc.Node
		split64(f, &flo, &fhi)

		var r1 gc.Node
		gc.Nodreg(&r1, t.Type, x86.REG_AX)
		gmove(&flo, &r1)
		gins(x86.AMOVB, &r1, t)
		splitclean()
		return

	case gc.TINT16<<16 | gc.TINT16, // same size
		gc.TINT16<<16 | gc.TUINT16,
		gc.TUINT16<<16 | gc.TINT16,
		gc.TUINT16<<16 | gc.TUINT16:
		a = x86.AMOVW

	case gc.TINT32<<16 | gc.TINT16, // truncate
		gc.TUINT32<<16 | gc.TINT16,
		gc.TINT32<<16 | gc.TUINT16,
		gc.TUINT32<<16 | gc.TUINT16:
		a = x86.AMOVW

		goto rsrc

	case gc.TINT64<<16 | gc.TINT16, // truncate low word
		gc.TUINT64<<16 | gc.TINT16,
		gc.TINT64<<16 | gc.TUINT16,
		gc.TUINT64<<16 | gc.TUINT16:
		var flo gc.Node
		var fhi gc.Node
		split64(f, &flo, &fhi)

		var r1 gc.Node
		gc.Nodreg(&r1, t.Type, x86.REG_AX)
		gmove(&flo, &r1)
		gins(x86.AMOVW, &r1, t)
		splitclean()
		return

	case gc.TINT32<<16 | gc.TINT32, // same size
		gc.TINT32<<16 | gc.TUINT32,
		gc.TUINT32<<16 | gc.TINT32,
		gc.TUINT32<<16 | gc.TUINT32:
		a = x86.AMOVL

	case gc.TINT64<<16 | gc.TINT32, // truncate
		gc.TUINT64<<16 | gc.TINT32,
		gc.TINT64<<16 | gc.TUINT32,
		gc.TUINT64<<16 | gc.TUINT32:
		var fhi gc.Node
		var flo gc.Node
		split64(f, &flo, &fhi)

		var r1 gc.Node
		gc.Nodreg(&r1, t.Type, x86.REG_AX)
		gmove(&flo, &r1)
		gins(x86.AMOVL, &r1, t)
		splitclean()
		return

	case gc.TINT64<<16 | gc.TINT64, // same size
		gc.TINT64<<16 | gc.TUINT64,
		gc.TUINT64<<16 | gc.TINT64,
		gc.TUINT64<<16 | gc.TUINT64:
		var fhi gc.Node
		var flo gc.Node
		split64(f, &flo, &fhi)

		var tlo gc.Node
		var thi gc.Node
		split64(t, &tlo, &thi)
		if f.Op == gc.OLITERAL {
			gins(x86.AMOVL, &flo, &tlo)
			gins(x86.AMOVL, &fhi, &thi)
		} else {
			var r1 gc.Node
			gc.Nodreg(&r1, gc.Types[gc.TUINT32], x86.REG_AX)
			var r2 gc.Node
			gc.Nodreg(&r2, gc.Types[gc.TUINT32], x86.REG_DX)
			gins(x86.AMOVL, &flo, &r1)
			gins(x86.AMOVL, &fhi, &r2)
			gins(x86.AMOVL, &r1, &tlo)
			gins(x86.AMOVL, &r2, &thi)
		}

		splitclean()
		splitclean()
		return

		/*
		 * integer up-conversions
		 */
	case gc.TINT8<<16 | gc.TINT16, // sign extend int8
		gc.TINT8<<16 | gc.TUINT16:
		a = x86.AMOVBWSX

		goto rdst

	case gc.TINT8<<16 | gc.TINT32,
		gc.TINT8<<16 | gc.TUINT32:
		a = x86.AMOVBLSX
		goto rdst

	case gc.TINT8<<16 | gc.TINT64, // convert via int32
		gc.TINT8<<16 | gc.TUINT64:
		cvt = gc.Types[gc.TINT32]

		goto hard

	case gc.TUINT8<<16 | gc.TINT16, // zero extend uint8
		gc.TUINT8<<16 | gc.TUINT16:
		a = x86.AMOVBWZX

		goto rdst

	case gc.TUINT8<<16 | gc.TINT32,
		gc.TUINT8<<16 | gc.TUINT32:
		a = x86.AMOVBLZX
		goto rdst

	case gc.TUINT8<<16 | gc.TINT64, // convert via uint32
		gc.TUINT8<<16 | gc.TUINT64:
		cvt = gc.Types[gc.TUINT32]

		goto hard

	case gc.TINT16<<16 | gc.TINT32, // sign extend int16
		gc.TINT16<<16 | gc.TUINT32:
		a = x86.AMOVWLSX

		goto rdst

	case gc.TINT16<<16 | gc.TINT64, // convert via int32
		gc.TINT16<<16 | gc.TUINT64:
		cvt = gc.Types[gc.TINT32]

		goto hard

	case gc.TUINT16<<16 | gc.TINT32, // zero extend uint16
		gc.TUINT16<<16 | gc.TUINT32:
		a = x86.AMOVWLZX

		goto rdst

	case gc.TUINT16<<16 | gc.TINT64, // convert via uint32
		gc.TUINT16<<16 | gc.TUINT64:
		cvt = gc.Types[gc.TUINT32]

		goto hard

	case gc.TINT32<<16 | gc.TINT64, // sign extend int32
		gc.TINT32<<16 | gc.TUINT64:
		var thi gc.Node
		var tlo gc.Node
		split64(t, &tlo, &thi)

		var flo gc.Node
		gc.Nodreg(&flo, tlo.Type, x86.REG_AX)
		var fhi gc.Node
		gc.Nodreg(&fhi, thi.Type, x86.REG_DX)
		gmove(f, &flo)
		gins(x86.ACDQ, nil, nil)
		gins(x86.AMOVL, &flo, &tlo)
		gins(x86.AMOVL, &fhi, &thi)
		splitclean()
		return

	case gc.TUINT32<<16 | gc.TINT64, // zero extend uint32
		gc.TUINT32<<16 | gc.TUINT64:
		var tlo gc.Node
		var thi gc.Node
		split64(t, &tlo, &thi)

		gmove(f, &tlo)
		gins(x86.AMOVL, ncon(0), &thi)
		splitclean()
		return
	}

	gins(a, f, t)
	return

	// requires register source
rsrc:
	gc.Regalloc(&r1, f.Type, t)

	gmove(f, &r1)
	gins(a, &r1, t)
	gc.Regfree(&r1)
	return

	// requires register destination
rdst:
	{
		gc.Regalloc(&r1, t.Type, t)

		gins(a, f, &r1)
		gmove(&r1, t)
		gc.Regfree(&r1)
		return
	}

	// requires register intermediate
hard:
	gc.Regalloc(&r1, cvt, t)

	gmove(f, &r1)
	gmove(&r1, t)
	gc.Regfree(&r1)
	return
}

func floatmove(f *gc.Node, t *gc.Node) {
	var r1 gc.Node

	ft := gc.Simsimtype(f.Type)
	tt := gc.Simsimtype(t.Type)
	cvt := t.Type

	// cannot have two floating point memory operands.
	if gc.Isfloat[ft] && gc.Isfloat[tt] && gc.Ismem(f) && gc.Ismem(t) {
		goto hard
	}

	// convert constant to desired type
	if f.Op == gc.OLITERAL {
		var con gc.Node
		gc.Convconst(&con, t.Type, &f.Val)
		f = &con
		ft = gc.Simsimtype(con.Type)

		// some constants can't move directly to memory.
		if gc.Ismem(t) {
			// float constants come from memory.
			if gc.Isfloat[tt] {
				goto hard
			}
		}
	}

	// value -> value copy, only one memory operand.
	// figure out the instruction to use.
	// break out of switch for one-instruction gins.
	// goto rdst for "destination must be register".
	// goto hard for "convert to cvt type first".
	// otherwise handle and return.

	switch uint32(ft)<<16 | uint32(tt) {
	default:
		if gc.Thearch.Use387 {
			floatmove_387(f, t)
		} else {
			floatmove_sse(f, t)
		}
		return

		// float to very long integer.
	case gc.TFLOAT32<<16 | gc.TINT64,
		gc.TFLOAT64<<16 | gc.TINT64:
		if f.Op == gc.OREGISTER {
			cvt = f.Type
			goto hardmem
		}

		var r1 gc.Node
		gc.Nodreg(&r1, gc.Types[ft], x86.REG_F0)
		if ft == gc.TFLOAT32 {
			gins(x86.AFMOVF, f, &r1)
		} else {
			gins(x86.AFMOVD, f, &r1)
		}

		// set round to zero mode during conversion
		var t1 gc.Node
		memname(&t1, gc.Types[gc.TUINT16])

		var t2 gc.Node
		memname(&t2, gc.Types[gc.TUINT16])
		gins(x86.AFSTCW, nil, &t1)
		gins(x86.AMOVW, ncon(0xf7f), &t2)
		gins(x86.AFLDCW, &t2, nil)
		if tt == gc.TINT16 {
			gins(x86.AFMOVWP, &r1, t)
		} else if tt == gc.TINT32 {
			gins(x86.AFMOVLP, &r1, t)
		} else {
			gins(x86.AFMOVVP, &r1, t)
		}
		gins(x86.AFLDCW, &t1, nil)
		return

	case gc.TFLOAT32<<16 | gc.TUINT64,
		gc.TFLOAT64<<16 | gc.TUINT64:
		if !gc.Ismem(f) {
			cvt = f.Type
			goto hardmem
		}

		bignodes()
		var f0 gc.Node
		gc.Nodreg(&f0, gc.Types[ft], x86.REG_F0)
		var f1 gc.Node
		gc.Nodreg(&f1, gc.Types[ft], x86.REG_F0+1)
		var ax gc.Node
		gc.Nodreg(&ax, gc.Types[gc.TUINT16], x86.REG_AX)

		if ft == gc.TFLOAT32 {
			gins(x86.AFMOVF, f, &f0)
		} else {
			gins(x86.AFMOVD, f, &f0)
		}

		// if 0 > v { answer = 0 }
		gins(x86.AFMOVD, &zerof, &f0)

		gins(x86.AFUCOMIP, &f0, &f1)
		p1 := gc.Gbranch(optoas(gc.OGT, gc.Types[tt]), nil, 0)

		// if 1<<64 <= v { answer = 0 too }
		gins(x86.AFMOVD, &two64f, &f0)

		gins(x86.AFUCOMIP, &f0, &f1)
		p2 := gc.Gbranch(optoas(gc.OGT, gc.Types[tt]), nil, 0)
		gc.Patch(p1, gc.Pc)
		gins(x86.AFMOVVP, &f0, t) // don't care about t, but will pop the stack
		var thi gc.Node
		var tlo gc.Node
		split64(t, &tlo, &thi)
		gins(x86.AMOVL, ncon(0), &tlo)
		gins(x86.AMOVL, ncon(0), &thi)
		splitclean()
		p1 = gc.Gbranch(obj.AJMP, nil, 0)
		gc.Patch(p2, gc.Pc)

		// in range; algorithm is:
		//	if small enough, use native float64 -> int64 conversion.
		//	otherwise, subtract 2^63, convert, and add it back.

		// set round to zero mode during conversion
		var t1 gc.Node
		memname(&t1, gc.Types[gc.TUINT16])

		var t2 gc.Node
		memname(&t2, gc.Types[gc.TUINT16])
		gins(x86.AFSTCW, nil, &t1)
		gins(x86.AMOVW, ncon(0xf7f), &t2)
		gins(x86.AFLDCW, &t2, nil)

		// actual work
		gins(x86.AFMOVD, &two63f, &f0)

		gins(x86.AFUCOMIP, &f0, &f1)
		p2 = gc.Gbranch(optoas(gc.OLE, gc.Types[tt]), nil, 0)
		gins(x86.AFMOVVP, &f0, t)
		p3 := gc.Gbranch(obj.AJMP, nil, 0)
		gc.Patch(p2, gc.Pc)
		gins(x86.AFMOVD, &two63f, &f0)
		gins(x86.AFSUBDP, &f0, &f1)
		gins(x86.AFMOVVP, &f0, t)
		split64(t, &tlo, &thi)
		gins(x86.AXORL, ncon(0x80000000), &thi) // + 2^63
		gc.Patch(p3, gc.Pc)
		splitclean()

		// restore rounding mode
		gins(x86.AFLDCW, &t1, nil)

		gc.Patch(p1, gc.Pc)
		return

		/*
		 * integer to float
		 */
	case gc.TINT64<<16 | gc.TFLOAT32,
		gc.TINT64<<16 | gc.TFLOAT64:
		if t.Op == gc.OREGISTER {
			goto hardmem
		}
		var f0 gc.Node
		gc.Nodreg(&f0, t.Type, x86.REG_F0)
		gins(x86.AFMOVV, f, &f0)
		if tt == gc.TFLOAT32 {
			gins(x86.AFMOVFP, &f0, t)
		} else {
			gins(x86.AFMOVDP, &f0, t)
		}
		return

		// algorithm is:
	//	if small enough, use native int64 -> float64 conversion.
	//	otherwise, halve (rounding to odd?), convert, and double.
	case gc.TUINT64<<16 | gc.TFLOAT32,
		gc.TUINT64<<16 | gc.TFLOAT64:
		var ax gc.Node
		gc.Nodreg(&ax, gc.Types[gc.TUINT32], x86.REG_AX)

		var dx gc.Node
		gc.Nodreg(&dx, gc.Types[gc.TUINT32], x86.REG_DX)
		var cx gc.Node
		gc.Nodreg(&cx, gc.Types[gc.TUINT32], x86.REG_CX)
		var t1 gc.Node
		gc.Tempname(&t1, f.Type)
		var tlo gc.Node
		var thi gc.Node
		split64(&t1, &tlo, &thi)
		gmove(f, &t1)
		gins(x86.ACMPL, &thi, ncon(0))
		p1 := gc.Gbranch(x86.AJLT, nil, 0)

		// native
		var r1 gc.Node
		gc.Nodreg(&r1, gc.Types[tt], x86.REG_F0)

		gins(x86.AFMOVV, &t1, &r1)
		if tt == gc.TFLOAT32 {
			gins(x86.AFMOVFP, &r1, t)
		} else {
			gins(x86.AFMOVDP, &r1, t)
		}
		p2 := gc.Gbranch(obj.AJMP, nil, 0)

		// simulated
		gc.Patch(p1, gc.Pc)

		gmove(&tlo, &ax)
		gmove(&thi, &dx)
		p1 = gins(x86.ASHRL, ncon(1), &ax)
		p1.From.Index = x86.REG_DX // double-width shift DX -> AX
		p1.From.Scale = 0
		gins(x86.AMOVL, ncon(0), &cx)
		gins(x86.ASETCC, nil, &cx)
		gins(x86.AORL, &cx, &ax)
		gins(x86.ASHRL, ncon(1), &dx)
		gmove(&dx, &thi)
		gmove(&ax, &tlo)
		gc.Nodreg(&r1, gc.Types[tt], x86.REG_F0)
		var r2 gc.Node
		gc.Nodreg(&r2, gc.Types[tt], x86.REG_F0+1)
		gins(x86.AFMOVV, &t1, &r1)
		gins(x86.AFMOVD, &r1, &r1)
		gins(x86.AFADDDP, &r1, &r2)
		if tt == gc.TFLOAT32 {
			gins(x86.AFMOVFP, &r1, t)
		} else {
			gins(x86.AFMOVDP, &r1, t)
		}
		gc.Patch(p2, gc.Pc)
		splitclean()
		return
	}

	// requires register intermediate
hard:
	gc.Regalloc(&r1, cvt, t)

	gmove(f, &r1)
	gmove(&r1, t)
	gc.Regfree(&r1)
	return

	// requires memory intermediate
hardmem:
	gc.Tempname(&r1, cvt)

	gmove(f, &r1)
	gmove(&r1, t)
	return
}

func floatmove_387(f *gc.Node, t *gc.Node) {
	var r1 gc.Node
	var a int

	ft := gc.Simsimtype(f.Type)
	tt := gc.Simsimtype(t.Type)
	cvt := t.Type

	switch uint32(ft)<<16 | uint32(tt) {
	default:
		goto fatal

		/*
		* float to integer
		 */
	case gc.TFLOAT32<<16 | gc.TINT16,
		gc.TFLOAT32<<16 | gc.TINT32,
		gc.TFLOAT32<<16 | gc.TINT64,
		gc.TFLOAT64<<16 | gc.TINT16,
		gc.TFLOAT64<<16 | gc.TINT32,
		gc.TFLOAT64<<16 | gc.TINT64:
		if t.Op == gc.OREGISTER {
			goto hardmem
		}
		var r1 gc.Node
		gc.Nodreg(&r1, gc.Types[ft], x86.REG_F0)
		if f.Op != gc.OREGISTER {
			if ft == gc.TFLOAT32 {
				gins(x86.AFMOVF, f, &r1)
			} else {
				gins(x86.AFMOVD, f, &r1)
			}
		}

		// set round to zero mode during conversion
		var t1 gc.Node
		memname(&t1, gc.Types[gc.TUINT16])

		var t2 gc.Node
		memname(&t2, gc.Types[gc.TUINT16])
		gins(x86.AFSTCW, nil, &t1)
		gins(x86.AMOVW, ncon(0xf7f), &t2)
		gins(x86.AFLDCW, &t2, nil)
		if tt == gc.TINT16 {
			gins(x86.AFMOVWP, &r1, t)
		} else if tt == gc.TINT32 {
			gins(x86.AFMOVLP, &r1, t)
		} else {
			gins(x86.AFMOVVP, &r1, t)
		}
		gins(x86.AFLDCW, &t1, nil)
		return

		// convert via int32.
	case gc.TFLOAT32<<16 | gc.TINT8,
		gc.TFLOAT32<<16 | gc.TUINT16,
		gc.TFLOAT32<<16 | gc.TUINT8,
		gc.TFLOAT64<<16 | gc.TINT8,
		gc.TFLOAT64<<16 | gc.TUINT16,
		gc.TFLOAT64<<16 | gc.TUINT8:
		var t1 gc.Node
		gc.Tempname(&t1, gc.Types[gc.TINT32])

		gmove(f, &t1)
		switch tt {
		default:
			gc.Fatal("gmove %v", gc.Nconv(t, 0))

		case gc.TINT8:
			gins(x86.ACMPL, &t1, ncon(-0x80&(1<<32-1)))
			p1 := gc.Gbranch(optoas(gc.OLT, gc.Types[gc.TINT32]), nil, -1)
			gins(x86.ACMPL, &t1, ncon(0x7f))
			p2 := gc.Gbranch(optoas(gc.OGT, gc.Types[gc.TINT32]), nil, -1)
			p3 := gc.Gbranch(obj.AJMP, nil, 0)
			gc.Patch(p1, gc.Pc)
			gc.Patch(p2, gc.Pc)
			gmove(ncon(-0x80&(1<<32-1)), &t1)
			gc.Patch(p3, gc.Pc)
			gmove(&t1, t)

		case gc.TUINT8:
			gins(x86.ATESTL, ncon(0xffffff00), &t1)
			p1 := gc.Gbranch(x86.AJEQ, nil, +1)
			gins(x86.AMOVL, ncon(0), &t1)
			gc.Patch(p1, gc.Pc)
			gmove(&t1, t)

		case gc.TUINT16:
			gins(x86.ATESTL, ncon(0xffff0000), &t1)
			p1 := gc.Gbranch(x86.AJEQ, nil, +1)
			gins(x86.AMOVL, ncon(0), &t1)
			gc.Patch(p1, gc.Pc)
			gmove(&t1, t)
		}

		return

		// convert via int64.
	case gc.TFLOAT32<<16 | gc.TUINT32,
		gc.TFLOAT64<<16 | gc.TUINT32:
		cvt = gc.Types[gc.TINT64]

		goto hardmem

		/*
		 * integer to float
		 */
	case gc.TINT16<<16 | gc.TFLOAT32,
		gc.TINT16<<16 | gc.TFLOAT64,
		gc.TINT32<<16 | gc.TFLOAT32,
		gc.TINT32<<16 | gc.TFLOAT64,
		gc.TINT64<<16 | gc.TFLOAT32,
		gc.TINT64<<16 | gc.TFLOAT64:
		if t.Op != gc.OREGISTER {
			goto hard
		}
		if f.Op == gc.OREGISTER {
			cvt = f.Type
			goto hardmem
		}

		switch ft {
		case gc.TINT16:
			a = x86.AFMOVW

		case gc.TINT32:
			a = x86.AFMOVL

		default:
			a = x86.AFMOVV
		}

		// convert via int32 memory
	case gc.TINT8<<16 | gc.TFLOAT32,
		gc.TINT8<<16 | gc.TFLOAT64,
		gc.TUINT16<<16 | gc.TFLOAT32,
		gc.TUINT16<<16 | gc.TFLOAT64,
		gc.TUINT8<<16 | gc.TFLOAT32,
		gc.TUINT8<<16 | gc.TFLOAT64:
		cvt = gc.Types[gc.TINT32]

		goto hardmem

		// convert via int64 memory
	case gc.TUINT32<<16 | gc.TFLOAT32,
		gc.TUINT32<<16 | gc.TFLOAT64:
		cvt = gc.Types[gc.TINT64]

		goto hardmem

		// The way the code generator uses floating-point
	// registers, a move from F0 to F0 is intended as a no-op.
	// On the x86, it's not: it pushes a second copy of F0
	// on the floating point stack.  So toss it away here.
	// Also, F0 is the *only* register we ever evaluate
	// into, so we should only see register/register as F0/F0.
	/*
	 * float to float
	 */
	case gc.TFLOAT32<<16 | gc.TFLOAT32,
		gc.TFLOAT64<<16 | gc.TFLOAT64:
		if gc.Ismem(f) && gc.Ismem(t) {
			goto hard
		}
		if f.Op == gc.OREGISTER && t.Op == gc.OREGISTER {
			if f.Val.U.Reg != x86.REG_F0 || t.Val.U.Reg != x86.REG_F0 {
				goto fatal
			}
			return
		}

		a = x86.AFMOVF
		if ft == gc.TFLOAT64 {
			a = x86.AFMOVD
		}
		if gc.Ismem(t) {
			if f.Op != gc.OREGISTER || f.Val.U.Reg != x86.REG_F0 {
				gc.Fatal("gmove %v", gc.Nconv(f, 0))
			}
			a = x86.AFMOVFP
			if ft == gc.TFLOAT64 {
				a = x86.AFMOVDP
			}
		}

	case gc.TFLOAT32<<16 | gc.TFLOAT64:
		if gc.Ismem(f) && gc.Ismem(t) {
			goto hard
		}
		if f.Op == gc.OREGISTER && t.Op == gc.OREGISTER {
			if f.Val.U.Reg != x86.REG_F0 || t.Val.U.Reg != x86.REG_F0 {
				goto fatal
			}
			return
		}

		if f.Op == gc.OREGISTER {
			gins(x86.AFMOVDP, f, t)
		} else {
			gins(x86.AFMOVF, f, t)
		}
		return

	case gc.TFLOAT64<<16 | gc.TFLOAT32:
		if gc.Ismem(f) && gc.Ismem(t) {
			goto hard
		}
		if f.Op == gc.OREGISTER && t.Op == gc.OREGISTER {
			var r1 gc.Node
			gc.Tempname(&r1, gc.Types[gc.TFLOAT32])
			gins(x86.AFMOVFP, f, &r1)
			gins(x86.AFMOVF, &r1, t)
			return
		}

		if f.Op == gc.OREGISTER {
			gins(x86.AFMOVFP, f, t)
		} else {
			gins(x86.AFMOVD, f, t)
		}
		return
	}

	gins(a, f, t)
	return

	// requires register intermediate
hard:
	gc.Regalloc(&r1, cvt, t)

	gmove(f, &r1)
	gmove(&r1, t)
	gc.Regfree(&r1)
	return

	// requires memory intermediate
hardmem:
	gc.Tempname(&r1, cvt)

	gmove(f, &r1)
	gmove(&r1, t)
	return

	// should not happen
fatal:
	gc.Fatal("gmove %v -> %v", gc.Nconv(f, obj.FmtLong), gc.Nconv(t, obj.FmtLong))

	return
}

func floatmove_sse(f *gc.Node, t *gc.Node) {
	var r1 gc.Node
	var cvt *gc.Type
	var a int

	ft := gc.Simsimtype(f.Type)
	tt := gc.Simsimtype(t.Type)

	switch uint32(ft)<<16 | uint32(tt) {
	// should not happen
	default:
		gc.Fatal("gmove %v -> %v", gc.Nconv(f, 0), gc.Nconv(t, 0))

		return

		// convert via int32.
	/*
	* float to integer
	 */
	case gc.TFLOAT32<<16 | gc.TINT16,
		gc.TFLOAT32<<16 | gc.TINT8,
		gc.TFLOAT32<<16 | gc.TUINT16,
		gc.TFLOAT32<<16 | gc.TUINT8,
		gc.TFLOAT64<<16 | gc.TINT16,
		gc.TFLOAT64<<16 | gc.TINT8,
		gc.TFLOAT64<<16 | gc.TUINT16,
		gc.TFLOAT64<<16 | gc.TUINT8:
		cvt = gc.Types[gc.TINT32]

		goto hard

		// convert via int64.
	case gc.TFLOAT32<<16 | gc.TUINT32,
		gc.TFLOAT64<<16 | gc.TUINT32:
		cvt = gc.Types[gc.TINT64]

		goto hardmem

	case gc.TFLOAT32<<16 | gc.TINT32:
		a = x86.ACVTTSS2SL
		goto rdst

	case gc.TFLOAT64<<16 | gc.TINT32:
		a = x86.ACVTTSD2SL
		goto rdst

		// convert via int32 memory
	/*
	 * integer to float
	 */
	case gc.TINT8<<16 | gc.TFLOAT32,
		gc.TINT8<<16 | gc.TFLOAT64,
		gc.TINT16<<16 | gc.TFLOAT32,
		gc.TINT16<<16 | gc.TFLOAT64,
		gc.TUINT16<<16 | gc.TFLOAT32,
		gc.TUINT16<<16 | gc.TFLOAT64,
		gc.TUINT8<<16 | gc.TFLOAT32,
		gc.TUINT8<<16 | gc.TFLOAT64:
		cvt = gc.Types[gc.TINT32]

		goto hard

		// convert via int64 memory
	case gc.TUINT32<<16 | gc.TFLOAT32,
		gc.TUINT32<<16 | gc.TFLOAT64:
		cvt = gc.Types[gc.TINT64]

		goto hardmem

	case gc.TINT32<<16 | gc.TFLOAT32:
		a = x86.ACVTSL2SS
		goto rdst

	case gc.TINT32<<16 | gc.TFLOAT64:
		a = x86.ACVTSL2SD
		goto rdst

		/*
		 * float to float
		 */
	case gc.TFLOAT32<<16 | gc.TFLOAT32:
		a = x86.AMOVSS

	case gc.TFLOAT64<<16 | gc.TFLOAT64:
		a = x86.AMOVSD

	case gc.TFLOAT32<<16 | gc.TFLOAT64:
		a = x86.ACVTSS2SD
		goto rdst

	case gc.TFLOAT64<<16 | gc.TFLOAT32:
		a = x86.ACVTSD2SS
		goto rdst
	}

	gins(a, f, t)
	return

	// requires register intermediate
hard:
	gc.Regalloc(&r1, cvt, t)

	gmove(f, &r1)
	gmove(&r1, t)
	gc.Regfree(&r1)
	return

	// requires memory intermediate
hardmem:
	gc.Tempname(&r1, cvt)

	gmove(f, &r1)
	gmove(&r1, t)
	return

	// requires register destination
rdst:
	gc.Regalloc(&r1, t.Type, t)

	gins(a, f, &r1)
	gmove(&r1, t)
	gc.Regfree(&r1)
	return
}

func samaddr(f *gc.Node, t *gc.Node) bool {
	if f.Op != t.Op {
		return false
	}

	switch f.Op {
	case gc.OREGISTER:
		if f.Val.U.Reg != t.Val.U.Reg {
			break
		}
		return true
	}

	return false
}

/*
 * generate one instruction:
 *	as f, t
 */
func gins(as int, f *gc.Node, t *gc.Node) *obj.Prog {
	if as == x86.AFMOVF && f != nil && f.Op == gc.OREGISTER && t != nil && t.Op == gc.OREGISTER {
		gc.Fatal("gins MOVF reg, reg")
	}
	if as == x86.ACVTSD2SS && f != nil && f.Op == gc.OLITERAL {
		gc.Fatal("gins CVTSD2SS const")
	}
	if as == x86.AMOVSD && t != nil && t.Op == gc.OREGISTER && t.Val.U.Reg == x86.REG_F0 {
		gc.Fatal("gins MOVSD into F0")
	}

	if as == x86.AMOVL && f != nil && f.Op == gc.OADDR && f.Left.Op == gc.ONAME && f.Left.Class != gc.PEXTERN && f.Left.Class != gc.PFUNC {
		// Turn MOVL $xxx(FP/SP) into LEAL xxx.
		// These should be equivalent but most of the backend
		// only expects to see LEAL, because that's what we had
		// historically generated. Various hidden assumptions are baked in by now.
		as = x86.ALEAL
		f = f.Left
	}

	switch as {
	case x86.AMOVB,
		x86.AMOVW,
		x86.AMOVL:
		if f != nil && t != nil && samaddr(f, t) {
			return nil
		}

	case x86.ALEAL:
		if f != nil && gc.Isconst(f, gc.CTNIL) {
			gc.Fatal("gins LEAL nil %v", gc.Tconv(f.Type, 0))
		}
	}

	p := gc.Prog(as)
	gc.Naddr(&p.From, f)
	gc.Naddr(&p.To, t)

	if gc.Debug['g'] != 0 {
		fmt.Printf("%v\n", p)
	}

	w := 0
	switch as {
	case x86.AMOVB:
		w = 1

	case x86.AMOVW:
		w = 2

	case x86.AMOVL:
		w = 4
	}

	if true && w != 0 && f != nil && (p.From.Width > int64(w) || p.To.Width > int64(w)) {
		gc.Dump("bad width from:", f)
		gc.Dump("bad width to:", t)
		gc.Fatal("bad width: %v (%d, %d)\n", p, p.From.Width, p.To.Width)
	}

	if p.To.Type == obj.TYPE_ADDR && w > 0 {
		gc.Fatal("bad use of addr: %v", p)
	}

	return p
}

func ginsnop() {
	var reg gc.Node
	gc.Nodreg(&reg, gc.Types[gc.TINT], x86.REG_AX)
	gins(x86.AXCHGL, &reg, &reg)
}

func dotaddable(n *gc.Node, n1 *gc.Node) bool {
	if n.Op != gc.ODOT {
		return false
	}

	var oary [10]int64
	var nn *gc.Node
	o := gc.Dotoffset(n, oary[:], &nn)
	if nn != nil && nn.Addable && o == 1 && oary[0] >= 0 {
		*n1 = *nn
		n1.Type = n.Type
		n1.Xoffset += oary[0]
		return true
	}

	return false
}

func sudoclean() {
}

func sudoaddable(as int, n *gc.Node, a *obj.Addr) bool {
	*a = obj.Addr{}
	return false
}
