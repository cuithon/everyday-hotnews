// Inferno utils/6l/span.c
// https://bitbucket.org/inferno-os/inferno-os/src/default/utils/6l/span.c
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

package x86

import (
	"cmd/internal/obj"
	"cmd/internal/objabi"
	"cmd/internal/sys"
	"encoding/binary"
	"fmt"
	"log"
	"strings"
)

var (
	plan9privates *obj.LSym
	deferreturn   *obj.LSym
)

// Instruction layout.

const (
	// Loop alignment constants:
	// want to align loop entry to LoopAlign-byte boundary,
	// and willing to insert at most MaxLoopPad bytes of NOP to do so.
	// We define a loop entry as the target of a backward jump.
	//
	// gcc uses MaxLoopPad = 10 for its 'generic x86-64' config,
	// and it aligns all jump targets, not just backward jump targets.
	//
	// As of 6/1/2012, the effect of setting MaxLoopPad = 10 here
	// is very slight but negative, so the alignment is disabled by
	// setting MaxLoopPad = 0. The code is here for reference and
	// for future experiments.
	//
	LoopAlign  = 16
	MaxLoopPad = 0
)

type Optab struct {
	as     obj.As
	ytab   []ytab
	prefix uint8
	op     [23]uint8
}

type Movtab struct {
	as   obj.As
	ft   uint8
	f3t  uint8
	tt   uint8
	code uint8
	op   [4]uint8
}

const (
	Yxxx = iota
	Ynone
	Yi0 // $0
	Yi1 // $1
	Yu2 // $x, x fits in uint2
	Yi8 // $x, x fits in int8
	Yu8 // $x, x fits in uint8
	Yu7 // $x, x in 0..127 (fits in both int8 and uint8)
	Ys32
	Yi32
	Yi64
	Yiauto
	Yal
	Ycl
	Yax
	Ycx
	Yrb
	Yrl
	Yrl32 // Yrl on 32-bit system
	Yrf
	Yf0
	Yrx
	Ymb
	Yml
	Ym
	Ybr
	Ycs
	Yss
	Yds
	Yes
	Yfs
	Ygs
	Ygdtr
	Yidtr
	Yldtr
	Ymsw
	Ytask
	Ycr0
	Ycr1
	Ycr2
	Ycr3
	Ycr4
	Ycr5
	Ycr6
	Ycr7
	Ycr8
	Ydr0
	Ydr1
	Ydr2
	Ydr3
	Ydr4
	Ydr5
	Ydr6
	Ydr7
	Ytr0
	Ytr1
	Ytr2
	Ytr3
	Ytr4
	Ytr5
	Ytr6
	Ytr7
	Ymr
	Ymm
	Yxr
	Yxm
	Yxvm // VSIB vector array; vm32x/vm64x
	Yyr
	Yym
	Yyvm // VSIB vector array; vm32y/vm64y
	Ytls
	Ytextsize
	Yindir
	Ymax
)

const (
	Zxxx = iota
	Zlit
	Zlitm_r
	Z_rp
	Zbr
	Zcall
	Zcallcon
	Zcallduff
	Zcallind
	Zcallindreg
	Zib_
	Zib_rp
	Zibo_m
	Zibo_m_xm
	Zil_
	Zil_rp
	Ziq_rp
	Zilo_m
	Zjmp
	Zjmpcon
	Zloop
	Zo_iw
	Zm_o
	Zm_r
	Zm2_r
	Zm_r_xm
	Zm_r_i_xm
	Zm_r_xm_nr
	Zr_m_xm_nr
	Zibm_r /* mmx1,mmx2/mem64,imm8 */
	Zibr_m
	Zmb_r
	Zaut_r
	Zo_m
	Zo_m64
	Zpseudo
	Zr_m
	Zr_m_xm
	Zrp_
	Z_ib
	Z_il
	Zm_ibo
	Zm_ilo
	Zib_rr
	Zil_rr
	Zbyte
	Zvex_rm_v_r
	Zvex_rm_v_ro
	Zvex_r_v_rm
	Zvex_v_rm_r
	Zvex_i_rm_r
	Zvex_i_r_v
	Zvex_i_rm_v_r
	Zvex
	Zvex_rm_r_vo
	Zvex_i_r_rm
	Zvex_hr_rm_v_r

	Zmax
)

const (
	Px   = 0
	Px1  = 1    // symbolic; exact value doesn't matter
	P32  = 0x32 /* 32-bit only */
	Pe   = 0x66 /* operand escape */
	Pm   = 0x0f /* 2byte opcode escape */
	Pq   = 0xff /* both escapes: 66 0f */
	Pb   = 0xfe /* byte operands */
	Pf2  = 0xf2 /* xmm escape 1: f2 0f */
	Pf3  = 0xf3 /* xmm escape 2: f3 0f */
	Pef3 = 0xf5 /* xmm escape 2 with 16-bit prefix: 66 f3 0f */
	Pq3  = 0x67 /* xmm escape 3: 66 48 0f */
	Pq4  = 0x68 /* xmm escape 4: 66 0F 38 */
	Pq4w = 0x69 /* Pq4 with Rex.w 66 0F 38 */
	Pq5  = 0x6a /* xmm escape 5: F3 0F 38 */
	Pq5w = 0x6b /* Pq5 with Rex.w F3 0F 38 */
	Pfw  = 0xf4 /* Pf3 with Rex.w: f3 48 0f */
	Pw   = 0x48 /* Rex.w */
	Pw8  = 0x90 // symbolic; exact value doesn't matter
	Py   = 0x80 /* defaults to 64-bit mode */
	Py1  = 0x81 // symbolic; exact value doesn't matter
	Py3  = 0x83 // symbolic; exact value doesn't matter
	Pvex = 0x84 // symbolic: exact value doesn't matter

	Rxw = 1 << 3 /* =1, 64-bit operand size */
	Rxr = 1 << 2 /* extend modrm reg */
	Rxx = 1 << 1 /* extend sib index */
	Rxb = 1 << 0 /* extend modrm r/m, sib base, or opcode reg */
)

const (
	// Encoding for VEX prefix in tables.
	// The P, L, and W fields are chosen to match
	// their eventual locations in the VEX prefix bytes.

	// V field - 4 bits; ignored by encoder
	vexNOVSR = 0 // No VEX-SPECIFIED-REGISTER
	vexNDS   = 0
	vexNDD   = 0
	vexDDS   = 0
	// P field - 2 bits
	vex66 = 1 << 0
	vexF3 = 2 << 0
	vexF2 = 3 << 0
	// L field - 1 bit
	vexLZ  = 0 << 2
	vexLIG = 0 << 2
	vex128 = 0 << 2
	vex256 = 1 << 2
	// W field - 1 bit
	vexWIG = 0 << 7
	vexW0  = 0 << 7
	vexW1  = 1 << 7
	// M field - 5 bits, but mostly reserved; we can store up to 4
	vex0F   = 1 << 3
	vex0F38 = 2 << 3
	vex0F3A = 3 << 3

	// Combinations used in the manual.
	VEX_DDS_LIG_66_0F38_W1    = vexDDS | vexLIG | vex66 | vex0F38 | vexW1
	VEX_NDD_128_66_0F_WIG     = vexNDD | vex128 | vex66 | vex0F | vexWIG
	VEX_NDD_256_66_0F_WIG     = vexNDD | vex256 | vex66 | vex0F | vexWIG
	VEX_NDD_LZ_F2_0F38_W0     = vexNDD | vexLZ | vexF2 | vex0F38 | vexW0
	VEX_NDD_LZ_F2_0F38_W1     = vexNDD | vexLZ | vexF2 | vex0F38 | vexW1
	VEX_NDS_128_66_0F_WIG     = vexNDS | vex128 | vex66 | vex0F | vexWIG
	VEX_NDS_128_66_0F38_WIG   = vexNDS | vex128 | vex66 | vex0F38 | vexWIG
	VEX_NDS_128_F2_0F_WIG     = vexNDS | vex128 | vexF2 | vex0F | vexWIG
	VEX_NDS_256_66_0F_WIG     = vexNDS | vex256 | vex66 | vex0F | vexWIG
	VEX_NDS_256_66_0F38_WIG   = vexNDS | vex256 | vex66 | vex0F38 | vexWIG
	VEX_NDS_256_66_0F3A_W0    = vexNDS | vex256 | vex66 | vex0F3A | vexW0
	VEX_NDS_256_66_0F3A_WIG   = vexNDS | vex256 | vex66 | vex0F3A | vexWIG
	VEX_NDS_LZ_0F38_W0        = vexNDS | vexLZ | vex0F38 | vexW0
	VEX_NDS_LZ_0F38_W1        = vexNDS | vexLZ | vex0F38 | vexW1
	VEX_NDS_LZ_66_0F38_W0     = vexNDS | vexLZ | vex66 | vex0F38 | vexW0
	VEX_NDS_LZ_66_0F38_W1     = vexNDS | vexLZ | vex66 | vex0F38 | vexW1
	VEX_NDS_LZ_F2_0F38_W0     = vexNDS | vexLZ | vexF2 | vex0F38 | vexW0
	VEX_NDS_LZ_F2_0F38_W1     = vexNDS | vexLZ | vexF2 | vex0F38 | vexW1
	VEX_NDS_LZ_F3_0F38_W0     = vexNDS | vexLZ | vexF3 | vex0F38 | vexW0
	VEX_NDS_LZ_F3_0F38_W1     = vexNDS | vexLZ | vexF3 | vex0F38 | vexW1
	VEX_NOVSR_128_66_0F_WIG   = vexNOVSR | vex128 | vex66 | vex0F | vexWIG
	VEX_NOVSR_128_66_0F38_W0  = vexNOVSR | vex128 | vex66 | vex0F38 | vexW0
	VEX_NOVSR_128_66_0F38_WIG = vexNOVSR | vex128 | vex66 | vex0F38 | vexWIG
	VEX_NOVSR_128_F2_0F_WIG   = vexNOVSR | vex128 | vexF2 | vex0F | vexWIG
	VEX_NOVSR_128_F3_0F_WIG   = vexNOVSR | vex128 | vexF3 | vex0F | vexWIG
	VEX_NOVSR_256_66_0F_WIG   = vexNOVSR | vex256 | vex66 | vex0F | vexWIG
	VEX_NOVSR_256_66_0F38_W0  = vexNOVSR | vex256 | vex66 | vex0F38 | vexW0
	VEX_NOVSR_256_66_0F38_WIG = vexNOVSR | vex256 | vex66 | vex0F38 | vexWIG
	VEX_NOVSR_256_F2_0F_WIG   = vexNOVSR | vex256 | vexF2 | vex0F | vexWIG
	VEX_NOVSR_256_F3_0F_WIG   = vexNOVSR | vex256 | vexF3 | vex0F | vexWIG
	VEX_NOVSR_LZ_F2_0F3A_W0   = vexNOVSR | vexLZ | vexF2 | vex0F3A | vexW0
	VEX_NOVSR_LZ_F2_0F3A_W1   = vexNOVSR | vexLZ | vexF2 | vex0F3A | vexW1
)

var ycover [Ymax * Ymax]uint8

var reg [MAXREG]int

var regrex [MAXREG + 1]int

var ynone = []ytab{
	{Zlit, 1, argList{}},
}

var ytext = []ytab{
	{Zpseudo, 0, argList{Ymb, Ytextsize}},
	{Zpseudo, 1, argList{Ymb, Yi32, Ytextsize}},
}

var ynop = []ytab{
	{Zpseudo, 0, argList{}},
	{Zpseudo, 0, argList{Yiauto}},
	{Zpseudo, 0, argList{Yml}},
	{Zpseudo, 0, argList{Yrf}},
	{Zpseudo, 0, argList{Yxr}},
	{Zpseudo, 0, argList{Yiauto}},
	{Zpseudo, 0, argList{Yml}},
	{Zpseudo, 0, argList{Yrf}},
	{Zpseudo, 1, argList{Yxr}},
}

var yfuncdata = []ytab{
	{Zpseudo, 0, argList{Yi32, Ym}},
}

var ypcdata = []ytab{
	{Zpseudo, 0, argList{Yi32, Yi32}},
}

var yxorb = []ytab{
	{Zib_, 1, argList{Yi32, Yal}},
	{Zibo_m, 2, argList{Yi32, Ymb}},
	{Zr_m, 1, argList{Yrb, Ymb}},
	{Zm_r, 1, argList{Ymb, Yrb}},
}

var yaddl = []ytab{
	{Zibo_m, 2, argList{Yi8, Yml}},
	{Zil_, 1, argList{Yi32, Yax}},
	{Zilo_m, 2, argList{Yi32, Yml}},
	{Zr_m, 1, argList{Yrl, Yml}},
	{Zm_r, 1, argList{Yml, Yrl}},
}

var yincl = []ytab{
	{Z_rp, 1, argList{Yrl}},
	{Zo_m, 2, argList{Yml}},
}

var yincq = []ytab{
	{Zo_m, 2, argList{Yml}},
}

var ycmpb = []ytab{
	{Z_ib, 1, argList{Yal, Yi32}},
	{Zm_ibo, 2, argList{Ymb, Yi32}},
	{Zm_r, 1, argList{Ymb, Yrb}},
	{Zr_m, 1, argList{Yrb, Ymb}},
}

var ycmpl = []ytab{
	{Zm_ibo, 2, argList{Yml, Yi8}},
	{Z_il, 1, argList{Yax, Yi32}},
	{Zm_ilo, 2, argList{Yml, Yi32}},
	{Zm_r, 1, argList{Yml, Yrl}},
	{Zr_m, 1, argList{Yrl, Yml}},
}

var yshb = []ytab{
	{Zo_m, 2, argList{Yi1, Ymb}},
	{Zibo_m, 2, argList{Yu8, Ymb}},
	{Zo_m, 2, argList{Ycx, Ymb}},
}

var yshl = []ytab{
	{Zo_m, 2, argList{Yi1, Yml}},
	{Zibo_m, 2, argList{Yu8, Yml}},
	{Zo_m, 2, argList{Ycl, Yml}},
	{Zo_m, 2, argList{Ycx, Yml}},
}

var ytestl = []ytab{
	{Zil_, 1, argList{Yi32, Yax}},
	{Zilo_m, 2, argList{Yi32, Yml}},
	{Zr_m, 1, argList{Yrl, Yml}},
	{Zm_r, 1, argList{Yml, Yrl}},
}

var ymovb = []ytab{
	{Zr_m, 1, argList{Yrb, Ymb}},
	{Zm_r, 1, argList{Ymb, Yrb}},
	{Zib_rp, 1, argList{Yi32, Yrb}},
	{Zibo_m, 2, argList{Yi32, Ymb}},
}

var ybtl = []ytab{
	{Zibo_m, 2, argList{Yi8, Yml}},
	{Zr_m, 1, argList{Yrl, Yml}},
}

var ymovw = []ytab{
	{Zr_m, 1, argList{Yrl, Yml}},
	{Zm_r, 1, argList{Yml, Yrl}},
	{Zil_rp, 1, argList{Yi32, Yrl}},
	{Zilo_m, 2, argList{Yi32, Yml}},
	{Zaut_r, 2, argList{Yiauto, Yrl}},
}

var ymovl = []ytab{
	{Zr_m, 1, argList{Yrl, Yml}},
	{Zm_r, 1, argList{Yml, Yrl}},
	{Zil_rp, 1, argList{Yi32, Yrl}},
	{Zilo_m, 2, argList{Yi32, Yml}},
	{Zm_r_xm, 1, argList{Yml, Ymr}}, // MMX MOVD
	{Zr_m_xm, 1, argList{Ymr, Yml}}, // MMX MOVD
	{Zm_r_xm, 2, argList{Yml, Yxr}}, // XMM MOVD (32 bit)
	{Zr_m_xm, 2, argList{Yxr, Yml}}, // XMM MOVD (32 bit)
	{Zaut_r, 2, argList{Yiauto, Yrl}},
}

var yret = []ytab{
	{Zo_iw, 1, argList{}},
	{Zo_iw, 1, argList{Yi32}},
}

var ymovq = []ytab{
	// valid in 32-bit mode
	{Zm_r_xm_nr, 1, argList{Ym, Ymr}},  // 0x6f MMX MOVQ (shorter encoding)
	{Zr_m_xm_nr, 1, argList{Ymr, Ym}},  // 0x7f MMX MOVQ
	{Zm_r_xm_nr, 2, argList{Yxr, Ymr}}, // Pf2, 0xd6 MOVDQ2Q
	{Zm_r_xm_nr, 2, argList{Yxm, Yxr}}, // Pf3, 0x7e MOVQ xmm1/m64 -> xmm2
	{Zr_m_xm_nr, 2, argList{Yxr, Yxm}}, // Pe, 0xd6 MOVQ xmm1 -> xmm2/m64

	// valid only in 64-bit mode, usually with 64-bit prefix
	{Zr_m, 1, argList{Yrl, Yml}},      // 0x89
	{Zm_r, 1, argList{Yml, Yrl}},      // 0x8b
	{Zilo_m, 2, argList{Ys32, Yrl}},   // 32 bit signed 0xc7,(0)
	{Ziq_rp, 1, argList{Yi64, Yrl}},   // 0xb8 -- 32/64 bit immediate
	{Zilo_m, 2, argList{Yi32, Yml}},   // 0xc7,(0)
	{Zm_r_xm, 1, argList{Ymm, Ymr}},   // 0x6e MMX MOVD
	{Zr_m_xm, 1, argList{Ymr, Ymm}},   // 0x7e MMX MOVD
	{Zm_r_xm, 2, argList{Yml, Yxr}},   // Pe, 0x6e MOVD xmm load
	{Zr_m_xm, 2, argList{Yxr, Yml}},   // Pe, 0x7e MOVD xmm store
	{Zaut_r, 1, argList{Yiauto, Yrl}}, // 0 built-in LEAQ
}

var ym_rl = []ytab{
	{Zm_r, 1, argList{Ym, Yrl}},
}

var yrl_m = []ytab{
	{Zr_m, 1, argList{Yrl, Ym}},
}

var ymb_rl = []ytab{
	{Zmb_r, 1, argList{Ymb, Yrl}},
}

var yml_rl = []ytab{
	{Zm_r, 1, argList{Yml, Yrl}},
}

var yrl_ml = []ytab{
	{Zr_m, 1, argList{Yrl, Yml}},
}

var yml_mb = []ytab{
	{Zr_m, 1, argList{Yrb, Ymb}},
	{Zm_r, 1, argList{Ymb, Yrb}},
}

var yrb_mb = []ytab{
	{Zr_m, 1, argList{Yrb, Ymb}},
}

var yxchg = []ytab{
	{Z_rp, 1, argList{Yax, Yrl}},
	{Zrp_, 1, argList{Yrl, Yax}},
	{Zr_m, 1, argList{Yrl, Yml}},
	{Zm_r, 1, argList{Yml, Yrl}},
}

var ydivl = []ytab{
	{Zm_o, 2, argList{Yml}},
}

var ydivb = []ytab{
	{Zm_o, 2, argList{Ymb}},
}

var yimul = []ytab{
	{Zm_o, 2, argList{Yml}},
	{Zib_rr, 1, argList{Yi8, Yrl}},
	{Zil_rr, 1, argList{Yi32, Yrl}},
	{Zm_r, 2, argList{Yml, Yrl}},
}

var yimul3 = []ytab{
	{Zibm_r, 2, argList{Yi8, Yml, Yrl}},
}

var ybyte = []ytab{
	{Zbyte, 1, argList{Yi64}},
}

var yin = []ytab{
	{Zib_, 1, argList{Yi32}},
	{Zlit, 1, argList{}},
}

var yint = []ytab{
	{Zib_, 1, argList{Yi32}},
}

var ypushl = []ytab{
	{Zrp_, 1, argList{Yrl}},
	{Zm_o, 2, argList{Ym}},
	{Zib_, 1, argList{Yi8}},
	{Zil_, 1, argList{Yi32}},
}

var ypopl = []ytab{
	{Z_rp, 1, argList{Yrl}},
	{Zo_m, 2, argList{Ym}},
}

var yclflush = []ytab{
	{Zo_m, 2, argList{Ym}},
}

var ybswap = []ytab{
	{Z_rp, 2, argList{Yrl}},
}

var yscond = []ytab{
	{Zo_m, 2, argList{Ymb}},
}

var yjcond = []ytab{
	{Zbr, 0, argList{Ybr}},
	{Zbr, 0, argList{Yi0, Ybr}},
	{Zbr, 1, argList{Yi1, Ybr}},
}

var yloop = []ytab{
	{Zloop, 1, argList{Ybr}},
}

var ycall = []ytab{
	{Zcallindreg, 0, argList{Yml}},
	{Zcallindreg, 2, argList{Yrx, Yrx}},
	{Zcallind, 2, argList{Yindir}},
	{Zcall, 0, argList{Ybr}},
	{Zcallcon, 1, argList{Yi32}},
}

var yduff = []ytab{
	{Zcallduff, 1, argList{Yi32}},
}

var yjmp = []ytab{
	{Zo_m64, 2, argList{Yml}},
	{Zjmp, 0, argList{Ybr}},
	{Zjmpcon, 1, argList{Yi32}},
}

var yfmvd = []ytab{
	{Zm_o, 2, argList{Ym, Yf0}},
	{Zo_m, 2, argList{Yf0, Ym}},
	{Zm_o, 2, argList{Yrf, Yf0}},
	{Zo_m, 2, argList{Yf0, Yrf}},
}

var yfmvdp = []ytab{
	{Zo_m, 2, argList{Yf0, Ym}},
	{Zo_m, 2, argList{Yf0, Yrf}},
}

var yfmvf = []ytab{
	{Zm_o, 2, argList{Ym, Yf0}},
	{Zo_m, 2, argList{Yf0, Ym}},
}

var yfmvx = []ytab{
	{Zm_o, 2, argList{Ym, Yf0}},
}

var yfmvp = []ytab{
	{Zo_m, 2, argList{Yf0, Ym}},
}

var yfcmv = []ytab{
	{Zm_o, 2, argList{Yrf, Yf0}},
}

var yfadd = []ytab{
	{Zm_o, 2, argList{Ym, Yf0}},
	{Zm_o, 2, argList{Yrf, Yf0}},
	{Zo_m, 2, argList{Yf0, Yrf}},
}

var yfxch = []ytab{
	{Zo_m, 2, argList{Yf0, Yrf}},
	{Zm_o, 2, argList{Yrf, Yf0}},
}

var ycompp = []ytab{
	{Zo_m, 2, argList{Yf0, Yrf}}, /* botch is really f0,f1 */
}

var ystsw = []ytab{
	{Zo_m, 2, argList{Ym}},
	{Zlit, 1, argList{Yax}},
}

var ysvrs_mo = []ytab{
	{Zm_o, 2, argList{Ym}},
}

// unaryDst version of "ysvrs_mo".
var ysvrs_om = []ytab{
	{Zo_m, 2, argList{Ym}},
}

var ymm = []ytab{
	{Zm_r_xm, 1, argList{Ymm, Ymr}},
	{Zm_r_xm, 2, argList{Yxm, Yxr}},
}

var yxm = []ytab{
	{Zm_r_xm, 1, argList{Yxm, Yxr}},
}

var yxm_q4 = []ytab{
	{Zm_r, 1, argList{Yxm, Yxr}},
}

var yxcvm1 = []ytab{
	{Zm_r_xm, 2, argList{Yxm, Yxr}},
	{Zm_r_xm, 2, argList{Yxm, Ymr}},
}

var yxcvm2 = []ytab{
	{Zm_r_xm, 2, argList{Yxm, Yxr}},
	{Zm_r_xm, 2, argList{Ymm, Yxr}},
}

var yxr = []ytab{
	{Zm_r_xm, 1, argList{Yxr, Yxr}},
}

var yxr_ml = []ytab{
	{Zr_m_xm, 1, argList{Yxr, Yml}},
}

var ymr = []ytab{
	{Zm_r, 1, argList{Ymr, Ymr}},
}

var ymr_ml = []ytab{
	{Zr_m_xm, 1, argList{Ymr, Yml}},
}

var yxcmpi = []ytab{
	{Zm_r_i_xm, 2, argList{Yxm, Yxr, Yi8}},
}

var yxmov = []ytab{
	{Zm_r_xm, 1, argList{Yxm, Yxr}},
	{Zr_m_xm, 1, argList{Yxr, Yxm}},
}

var yxcvfl = []ytab{
	{Zm_r_xm, 1, argList{Yxm, Yrl}},
}

var yxcvlf = []ytab{
	{Zm_r_xm, 1, argList{Yml, Yxr}},
}

var yxcvfq = []ytab{
	{Zm_r_xm, 2, argList{Yxm, Yrl}},
}

var yxcvqf = []ytab{
	{Zm_r_xm, 2, argList{Yml, Yxr}},
}

var yps = []ytab{
	{Zm_r_xm, 1, argList{Ymm, Ymr}},
	{Zibo_m_xm, 2, argList{Yi8, Ymr}},
	{Zm_r_xm, 2, argList{Yxm, Yxr}},
	{Zibo_m_xm, 3, argList{Yi8, Yxr}},
}

var yxrrl = []ytab{
	{Zm_r, 1, argList{Yxr, Yrl}},
}

var ymrxr = []ytab{
	{Zm_r, 1, argList{Ymr, Yxr}},
	{Zm_r_xm, 1, argList{Yxm, Yxr}},
}

var ymshuf = []ytab{
	{Zibm_r, 2, argList{Yi8, Ymm, Ymr}},
}

var ymshufb = []ytab{
	{Zm2_r, 2, argList{Yxm, Yxr}},
}

// It should never have more than 1 entry,
// because some optab entries you opcode secuences that
// are longer than 2 bytes (zoffset=2 here),
// ROUNDPD and ROUNDPS and recently added BLENDPD,
// to name a few.
var yxshuf = []ytab{
	{Zibm_r, 2, argList{Yu8, Yxm, Yxr}},
}

var yextrw = []ytab{
	{Zibm_r, 2, argList{Yu8, Yxr, Yrl}},
}

var yextr = []ytab{
	{Zibr_m, 3, argList{Yu8, Yxr, Ymm}},
}

var yinsrw = []ytab{
	{Zibm_r, 2, argList{Yu8, Yml, Yxr}},
}

var yinsr = []ytab{
	{Zibm_r, 3, argList{Yu8, Ymm, Yxr}},
}

var ypsdq = []ytab{
	{Zibo_m, 2, argList{Yi8, Yxr}},
}

var ymskb = []ytab{
	{Zm_r_xm, 2, argList{Yxr, Yrl}},
	{Zm_r_xm, 1, argList{Ymr, Yrl}},
}

var ycrc32l = []ytab{
	{Zlitm_r, 0, argList{Yml, Yrl}},
}

var yprefetch = []ytab{
	{Zm_o, 2, argList{Ym}},
}

var yaes = []ytab{
	{Zlitm_r, 2, argList{Yxm, Yxr}},
}

var yxbegin = []ytab{
	{Zjmp, 1, argList{Ybr}},
}

var yxabort = []ytab{
	{Zib_, 1, argList{Yu8}},
}

var ylddqu = []ytab{
	{Zm_r, 1, argList{Ym, Yxr}},
}

var ypalignr = []ytab{
	{Zibm_r, 2, argList{Yu8, Yxm, Yxr}},
}

// VEX instructions that come in two forms:
//	VTHING xmm2/m128, xmmV, xmm1
//	VTHING ymm2/m256, ymmV, ymm1
// The opcode array in the corresponding Optab entry
// should contain the (VEX prefixes, opcode byte) pair
// for each of the two forms.
// For example, the entries for VPXOR are:
//
//	VPXOR xmm2/m128, xmmV, xmm1
//	VEX.NDS.128.66.0F.WIG EF /r
//
//	VPXOR ymm2/m256, ymmV, ymm1
//	VEX.NDS.256.66.0F.WIG EF /r
//
// Produce this Optab entry:
//
//	{AVPXOR, yvex_xy3, Pvex, [23]uint8{VEX_NDS_128_66_0F_WIG, 0xEF, VEX_NDS_256_66_0F_WIG, 0xEF}}
//
var yvex_xy3 = []ytab{
	{Zvex_rm_v_r, 2, argList{Yxm, Yxr, Yxr}},
	{Zvex_rm_v_r, 2, argList{Yym, Yyr, Yyr}},
}

var yvex_x3 = []ytab{
	{Zvex_rm_v_r, 2, argList{Yxm, Yxr, Yxr}},
}

var yvex_ri3 = []ytab{
	{Zvex_i_rm_r, 2, argList{Yi8, Ymb, Yrl}},
}

var yvex_xyi3 = []ytab{
	{Zvex_i_rm_r, 2, argList{Yu8, Yxm, Yxr}},
	{Zvex_i_rm_r, 2, argList{Yu8, Yym, Yyr}},
	{Zvex_i_rm_r, 2, argList{Yi8, Yxm, Yxr}},
	{Zvex_i_rm_r, 2, argList{Yi8, Yym, Yyr}},
}

var yvex_yyi4 = []ytab{
	{Zvex_i_rm_v_r, 2, argList{Yu8, Yym, Yyr, Yyr}},
}

var yvex_xyi4 = []ytab{
	{Zvex_i_rm_v_r, 2, argList{Yu8, Yxm, Yyr, Yyr}},
}

var yvex_shift = []ytab{
	{Zvex_i_r_v, 3, argList{Yi8, Yxr, Yxr}},
	{Zvex_i_r_v, 3, argList{Yi8, Yyr, Yyr}},
	{Zvex_rm_v_r, 2, argList{Yxm, Yxr, Yxr}},
	{Zvex_rm_v_r, 2, argList{Yxm, Yyr, Yyr}},
}

var yvex_shift_dq = []ytab{
	{Zvex_i_r_v, 3, argList{Yi8, Yxr, Yxr}},
	{Zvex_i_r_v, 3, argList{Yi8, Yyr, Yyr}},
}

var yvex_r3 = []ytab{
	{Zvex_rm_v_r, 2, argList{Yml, Yrl, Yrl}},
}

var yvex_vmr3 = []ytab{
	{Zvex_v_rm_r, 2, argList{Yrl, Yml, Yrl}},
}

var yvex_xy2 = []ytab{
	{Zvex_rm_v_r, 2, argList{Yxm, Yxr}},
	{Zvex_rm_v_r, 2, argList{Yym, Yyr}},
}

var yvex_xyr2 = []ytab{
	{Zvex_rm_v_r, 2, argList{Yxr, Yrl}},
	{Zvex_rm_v_r, 2, argList{Yyr, Yrl}},
}

var yvex_vmovdqa = []ytab{
	{Zvex_rm_v_r, 2, argList{Yxm, Yxr}},
	{Zvex_r_v_rm, 2, argList{Yxr, Yxm}},
	{Zvex_rm_v_r, 2, argList{Yym, Yyr}},
	{Zvex_r_v_rm, 2, argList{Yyr, Yym}},
}

var yvex_vmovntdq = []ytab{
	{Zvex_r_v_rm, 2, argList{Yxr, Ym}},
	{Zvex_r_v_rm, 2, argList{Yyr, Ym}},
}

var yvex_vpbroadcast = []ytab{
	{Zvex_rm_v_r, 2, argList{Yxm, Yxr}},
	{Zvex_rm_v_r, 2, argList{Yxm, Yyr}},
}

var yvex_vpbroadcast_sd = []ytab{
	{Zvex_rm_v_r, 2, argList{Yxm, Yyr}},
}

var yvex_vpextrw = []ytab{
	{Zvex_i_rm_r, 2, argList{Yi8, Yxr, Yrl}},
	{Zvex_i_r_rm, 2, argList{Yi8, Yxr, Yml}},
}

var yvex_m = []ytab{
	{Zvex_rm_v_ro, 3, argList{Ym}},
}

var yvex_xx3 = []ytab{
	{Zvex_rm_v_r, 2, argList{Yxr, Yxr, Yxr}},
}

var yvex_yi3 = []ytab{
	{Zvex_i_r_rm, 2, argList{Yi8, Yyr, Yxm}},
}

var yvex_mxy = []ytab{
	{Zvex_rm_v_r, 2, argList{Ym, Yxr}},
	{Zvex_rm_v_r, 2, argList{Ym, Yyr}},
}

var yvex_yy3 = []ytab{
	{Zvex_rm_v_r, 2, argList{Yym, Yyr, Yyr}},
}

var yvex_xi3 = []ytab{
	{Zvex_i_rm_r, 2, argList{Yi8, Yxm, Yxr}},
}

var yvex_vpermpd = []ytab{
	{Zvex_i_rm_r, 2, argList{Yi8, Yym, Yyr}},
}

var yvex_vpermilp = []ytab{
	{Zvex_i_rm_r, 2, argList{Yi8, Yxm, Yxr}},
	{Zvex_rm_v_r, 2, argList{Yxm, Yxr, Yxr}},
	{Zvex_i_rm_r, 2, argList{Yi8, Yym, Yyr}},
	{Zvex_rm_v_r, 2, argList{Yym, Yyr, Yyr}},
}

var yvex_vcvtps2ph = []ytab{
	{Zvex_i_r_rm, 2, argList{Yi8, Yyr, Yxm}},
	{Zvex_i_r_rm, 2, argList{Yi8, Yxr, Yxm}},
}

var yvex_vbroadcastf = []ytab{
	{Zvex_rm_v_r, 2, argList{Ym, Yyr}},
}

var yvex_vmovd = []ytab{
	{Zvex_r_v_rm, 2, argList{Yxr, Yml}},
	{Zvex_rm_v_r, 2, argList{Yml, Yxr}},
}

var yvex_x2 = []ytab{
	{Zvex_rm_v_r, 2, argList{Yxm, Yxr}},
}

var yvex_y2 = []ytab{
	{Zvex_rm_v_r, 2, argList{Yym, Yxr}},
}

var yvex = []ytab{
	{Zvex, 2, argList{}},
}

var yvex_xx2 = []ytab{
	{Zvex_rm_v_r, 2, argList{Yxr, Yxr}},
}

var yvex_vpalignr = []ytab{
	{Zvex_i_rm_v_r, 2, argList{Yu8, Yxm, Yxr, Yxr}},
	{Zvex_i_rm_v_r, 2, argList{Yu8, Yym, Yyr, Yyr}},
}

var yvex_rxi4 = []ytab{
	{Zvex_i_rm_v_r, 2, argList{Yu8, Yml, Yxr, Yxr}},
}

var yvex_xxi4 = []ytab{
	{Zvex_i_rm_v_r, 2, argList{Yu8, Yxm, Yxr, Yxr}},
}

var yvex_xy4 = []ytab{
	{Zvex_hr_rm_v_r, 2, argList{Yxr, Yxm, Yxr, Yxr}},
	{Zvex_hr_rm_v_r, 2, argList{Yyr, Yym, Yyr, Yyr}},
}

var yvex_vpbroadcast_ss = []ytab{
	{Zvex_rm_v_r, 2, argList{Ym, Yxr}},
	{Zvex_rm_v_r, 2, argList{Yxr, Yxr}},
	{Zvex_rm_v_r, 2, argList{Ym, Yyr}},
	{Zvex_rm_v_r, 2, argList{Yxr, Yyr}},
}

var yvex_vblendvpd = []ytab{
	{Zvex_r_v_rm, 2, argList{Yxr, Yxr, Yml}},
	{Zvex_r_v_rm, 2, argList{Yyr, Yyr, Yml}},
	{Zvex_rm_v_r, 2, argList{Ym, Yxr, Yxr}},
	{Zvex_rm_v_r, 2, argList{Ym, Yyr, Yyr}},
}

var yvex_vmov = []ytab{
	{Zvex_r_v_rm, 2, argList{Yxr, Ym}},
	{Zvex_rm_v_r, 2, argList{Ym, Yxr}},
	{Zvex_rm_v_r, 2, argList{Yxr, Yxr, Yxr}},
	{Zvex_rm_v_r, 2, argList{Yxr, Yxr, Yxr}},
}

var yvex_vps = []ytab{
	{Zvex_rm_v_r, 2, argList{Yxm, Yxr, Yxr}},
	{Zvex_i_r_v, 3, argList{Yi8, Yxr, Yxr}},
	{Zvex_rm_v_r, 2, argList{Yxm, Yyr, Yyr}},
	{Zvex_i_r_v, 3, argList{Yi8, Yyr, Yyr}},
}

var yvex_r2 = []ytab{
	{Zvex_rm_r_vo, 3, argList{Yml, Yrl}},
}

var yvex_vpextr = []ytab{
	{Zvex_i_r_rm, 2, argList{Yi8, Yxr, Yml}},
}

var yvex_rx3 = []ytab{
	{Zvex_rm_v_r, 2, argList{Yml, Yxr, Yxr}},
}

var yvex_vcvtsd2si = []ytab{
	{Zvex_rm_v_r, 2, argList{Yxm, Yrl}},
}

var yvex_vmovhpd = []ytab{
	{Zvex_r_v_rm, 2, argList{Yxr, Ym}},
	{Zvex_rm_v_r, 2, argList{Ym, Yxr, Yxr}},
}

var yvex_vmovq = []ytab{
	{Zvex_r_v_rm, 2, argList{Yxr, Yml}},
	{Zvex_rm_v_r, 2, argList{Ym, Yxr}},
	{Zvex_rm_v_r, 2, argList{Yml, Yxr}},
	{Zvex_rm_v_r, 2, argList{Yxr, Yxr}},
	{Zvex_r_v_rm, 2, argList{Yxr, Yxm}},
}

var yvpgatherdq = []ytab{
	{Zvex_v_rm_r, 2, argList{Yxr, Yxvm, Yxr}},
	{Zvex_v_rm_r, 2, argList{Yyr, Yxvm, Yyr}},
}

var yvpgatherqq = []ytab{
	{Zvex_v_rm_r, 2, argList{Yxr, Yxvm, Yxr}},
	{Zvex_v_rm_r, 2, argList{Yyr, Yyvm, Yyr}},
}

var yvgatherqps = []ytab{
	{Zvex_v_rm_r, 2, argList{Yxr, Yxvm, Yxr}},
	{Zvex_v_rm_r, 2, argList{Yxr, Yyvm, Yxr}},
}

var ymmxmm0f38 = []ytab{
	{Zlitm_r, 3, argList{Ymm, Ymr}},
	{Zlitm_r, 5, argList{Yxm, Yxr}},
}

var yextractps = []ytab{
	{Zibr_m, 2, argList{Yu2, Yxr, Yml}},
}

/*
 * You are doasm, holding in your hand a *obj.Prog with p.As set to, say,
 * ACRC32, and p.From and p.To as operands (obj.Addr).  The linker scans optab
 * to find the entry with the given p.As and then looks through the ytable for
 * that instruction (the second field in the optab struct) for a line whose
 * first two values match the Ytypes of the p.From and p.To operands.  The
 * function oclass computes the specific Ytype of an operand and then the set
 * of more general Ytypes that it satisfies is implied by the ycover table, set
 * up in instinit.  For example, oclass distinguishes the constants 0 and 1
 * from the more general 8-bit constants, but instinit says
 *
 *        ycover[Yi0*Ymax+Ys32] = 1
 *        ycover[Yi1*Ymax+Ys32] = 1
 *        ycover[Yi8*Ymax+Ys32] = 1
 *
 * which means that Yi0, Yi1, and Yi8 all count as Ys32 (signed 32)
 * if that's what an instruction can handle.
 *
 * In parallel with the scan through the ytable for the appropriate line, there
 * is a z pointer that starts out pointing at the strange magic byte list in
 * the Optab struct.  With each step past a non-matching ytable line, z
 * advances by the 4th entry in the line.  When a matching line is found, that
 * z pointer has the extra data to use in laying down the instruction bytes.
 * The actual bytes laid down are a function of the 3rd entry in the line (that
 * is, the Ztype) and the z bytes.
 *
 * For example, let's look at AADDL.  The optab line says:
 *        {AADDL, yaddl, Px, [23]uint8{0x83, 00, 0x05, 0x81, 00, 0x01, 0x03}},
 *
 * and yaddl says
 *        var yaddl = []ytab{
 *                {Yi8, Ynone, Yml, Zibo_m, 2},
 *                {Yi32, Ynone, Yax, Zil_, 1},
 *                {Yi32, Ynone, Yml, Zilo_m, 2},
 *                {Yrl, Ynone, Yml, Zr_m, 1},
 *                {Yml, Ynone, Yrl, Zm_r, 1},
 *        }
 *
 * so there are 5 possible types of ADDL instruction that can be laid down, and
 * possible states used to lay them down (Ztype and z pointer, assuming z
 * points at [23]uint8{0x83, 00, 0x05,0x81, 00, 0x01, 0x03}) are:
 *
 *        Yi8, Yml -> Zibo_m, z (0x83, 00)
 *        Yi32, Yax -> Zil_, z+2 (0x05)
 *        Yi32, Yml -> Zilo_m, z+2+1 (0x81, 0x00)
 *        Yrl, Yml -> Zr_m, z+2+1+2 (0x01)
 *        Yml, Yrl -> Zm_r, z+2+1+2+1 (0x03)
 *
 * The Pconstant in the optab line controls the prefix bytes to emit.  That's
 * relatively straightforward as this program goes.
 *
 * The switch on yt.zcase in doasm implements the various Z cases.  Zibo_m, for
 * example, is an opcode byte (z[0]) then an asmando (which is some kind of
 * encoded addressing mode for the Yml arg), and then a single immediate byte.
 * Zilo_m is the same but a long (32-bit) immediate.
 */
var optab =
/*	as, ytab, andproto, opcode */
[]Optab{
	{obj.AXXX, nil, 0, [23]uint8{}},
	{AAAA, ynone, P32, [23]uint8{0x37}},
	{AAAD, ynone, P32, [23]uint8{0xd5, 0x0a}},
	{AAAM, ynone, P32, [23]uint8{0xd4, 0x0a}},
	{AAAS, ynone, P32, [23]uint8{0x3f}},
	{AADCB, yxorb, Pb, [23]uint8{0x14, 0x80, 02, 0x10, 0x12}},
	{AADCL, yaddl, Px, [23]uint8{0x83, 02, 0x15, 0x81, 02, 0x11, 0x13}},
	{AADCQ, yaddl, Pw, [23]uint8{0x83, 02, 0x15, 0x81, 02, 0x11, 0x13}},
	{AADCW, yaddl, Pe, [23]uint8{0x83, 02, 0x15, 0x81, 02, 0x11, 0x13}},
	{AADCXL, yml_rl, Pq4, [23]uint8{0xf6}},
	{AADCXQ, yml_rl, Pq4w, [23]uint8{0xf6}},
	{AADDB, yxorb, Pb, [23]uint8{0x04, 0x80, 00, 0x00, 0x02}},
	{AADDL, yaddl, Px, [23]uint8{0x83, 00, 0x05, 0x81, 00, 0x01, 0x03}},
	{AADDPD, yxm, Pq, [23]uint8{0x58}},
	{AADDPS, yxm, Pm, [23]uint8{0x58}},
	{AADDQ, yaddl, Pw, [23]uint8{0x83, 00, 0x05, 0x81, 00, 0x01, 0x03}},
	{AADDSD, yxm, Pf2, [23]uint8{0x58}},
	{AADDSS, yxm, Pf3, [23]uint8{0x58}},
	{AADDSUBPD, yxm, Pq, [23]uint8{0xd0}},
	{AADDSUBPS, yxm, Pf2, [23]uint8{0xd0}},
	{AADDW, yaddl, Pe, [23]uint8{0x83, 00, 0x05, 0x81, 00, 0x01, 0x03}},
	{AADOXL, yml_rl, Pq5, [23]uint8{0xf6}},
	{AADOXQ, yml_rl, Pq5w, [23]uint8{0xf6}},
	{AADJSP, nil, 0, [23]uint8{}},
	{AANDB, yxorb, Pb, [23]uint8{0x24, 0x80, 04, 0x20, 0x22}},
	{AANDL, yaddl, Px, [23]uint8{0x83, 04, 0x25, 0x81, 04, 0x21, 0x23}},
	{AANDNPD, yxm, Pq, [23]uint8{0x55}},
	{AANDNPS, yxm, Pm, [23]uint8{0x55}},
	{AANDPD, yxm, Pq, [23]uint8{0x54}},
	{AANDPS, yxm, Pm, [23]uint8{0x54}},
	{AANDQ, yaddl, Pw, [23]uint8{0x83, 04, 0x25, 0x81, 04, 0x21, 0x23}},
	{AANDW, yaddl, Pe, [23]uint8{0x83, 04, 0x25, 0x81, 04, 0x21, 0x23}},
	{AARPL, yrl_ml, P32, [23]uint8{0x63}},
	{ABOUNDL, yrl_m, P32, [23]uint8{0x62}},
	{ABOUNDW, yrl_m, Pe, [23]uint8{0x62}},
	{ABSFL, yml_rl, Pm, [23]uint8{0xbc}},
	{ABSFQ, yml_rl, Pw, [23]uint8{0x0f, 0xbc}},
	{ABSFW, yml_rl, Pq, [23]uint8{0xbc}},
	{ABSRL, yml_rl, Pm, [23]uint8{0xbd}},
	{ABSRQ, yml_rl, Pw, [23]uint8{0x0f, 0xbd}},
	{ABSRW, yml_rl, Pq, [23]uint8{0xbd}},
	{ABSWAPL, ybswap, Px, [23]uint8{0x0f, 0xc8}},
	{ABSWAPQ, ybswap, Pw, [23]uint8{0x0f, 0xc8}},
	{ABTCL, ybtl, Pm, [23]uint8{0xba, 07, 0xbb}},
	{ABTCQ, ybtl, Pw, [23]uint8{0x0f, 0xba, 07, 0x0f, 0xbb}},
	{ABTCW, ybtl, Pq, [23]uint8{0xba, 07, 0xbb}},
	{ABTL, ybtl, Pm, [23]uint8{0xba, 04, 0xa3}},
	{ABTQ, ybtl, Pw, [23]uint8{0x0f, 0xba, 04, 0x0f, 0xa3}},
	{ABTRL, ybtl, Pm, [23]uint8{0xba, 06, 0xb3}},
	{ABTRQ, ybtl, Pw, [23]uint8{0x0f, 0xba, 06, 0x0f, 0xb3}},
	{ABTRW, ybtl, Pq, [23]uint8{0xba, 06, 0xb3}},
	{ABTSL, ybtl, Pm, [23]uint8{0xba, 05, 0xab}},
	{ABTSQ, ybtl, Pw, [23]uint8{0x0f, 0xba, 05, 0x0f, 0xab}},
	{ABTSW, ybtl, Pq, [23]uint8{0xba, 05, 0xab}},
	{ABTW, ybtl, Pq, [23]uint8{0xba, 04, 0xa3}},
	{ABYTE, ybyte, Px, [23]uint8{1}},
	{obj.ACALL, ycall, Px, [23]uint8{0xff, 02, 0xff, 0x15, 0xe8}},
	{ACDQ, ynone, Px, [23]uint8{0x99}},
	{ACLC, ynone, Px, [23]uint8{0xf8}},
	{ACLD, ynone, Px, [23]uint8{0xfc}},
	{ACLFLUSH, yclflush, Pm, [23]uint8{0xae, 07}},
	{ACLI, ynone, Px, [23]uint8{0xfa}},
	{ACLTS, ynone, Pm, [23]uint8{0x06}},
	{ACMC, ynone, Px, [23]uint8{0xf5}},
	{ACMOVLCC, yml_rl, Pm, [23]uint8{0x43}},
	{ACMOVLCS, yml_rl, Pm, [23]uint8{0x42}},
	{ACMOVLEQ, yml_rl, Pm, [23]uint8{0x44}},
	{ACMOVLGE, yml_rl, Pm, [23]uint8{0x4d}},
	{ACMOVLGT, yml_rl, Pm, [23]uint8{0x4f}},
	{ACMOVLHI, yml_rl, Pm, [23]uint8{0x47}},
	{ACMOVLLE, yml_rl, Pm, [23]uint8{0x4e}},
	{ACMOVLLS, yml_rl, Pm, [23]uint8{0x46}},
	{ACMOVLLT, yml_rl, Pm, [23]uint8{0x4c}},
	{ACMOVLMI, yml_rl, Pm, [23]uint8{0x48}},
	{ACMOVLNE, yml_rl, Pm, [23]uint8{0x45}},
	{ACMOVLOC, yml_rl, Pm, [23]uint8{0x41}},
	{ACMOVLOS, yml_rl, Pm, [23]uint8{0x40}},
	{ACMOVLPC, yml_rl, Pm, [23]uint8{0x4b}},
	{ACMOVLPL, yml_rl, Pm, [23]uint8{0x49}},
	{ACMOVLPS, yml_rl, Pm, [23]uint8{0x4a}},
	{ACMOVQCC, yml_rl, Pw, [23]uint8{0x0f, 0x43}},
	{ACMOVQCS, yml_rl, Pw, [23]uint8{0x0f, 0x42}},
	{ACMOVQEQ, yml_rl, Pw, [23]uint8{0x0f, 0x44}},
	{ACMOVQGE, yml_rl, Pw, [23]uint8{0x0f, 0x4d}},
	{ACMOVQGT, yml_rl, Pw, [23]uint8{0x0f, 0x4f}},
	{ACMOVQHI, yml_rl, Pw, [23]uint8{0x0f, 0x47}},
	{ACMOVQLE, yml_rl, Pw, [23]uint8{0x0f, 0x4e}},
	{ACMOVQLS, yml_rl, Pw, [23]uint8{0x0f, 0x46}},
	{ACMOVQLT, yml_rl, Pw, [23]uint8{0x0f, 0x4c}},
	{ACMOVQMI, yml_rl, Pw, [23]uint8{0x0f, 0x48}},
	{ACMOVQNE, yml_rl, Pw, [23]uint8{0x0f, 0x45}},
	{ACMOVQOC, yml_rl, Pw, [23]uint8{0x0f, 0x41}},
	{ACMOVQOS, yml_rl, Pw, [23]uint8{0x0f, 0x40}},
	{ACMOVQPC, yml_rl, Pw, [23]uint8{0x0f, 0x4b}},
	{ACMOVQPL, yml_rl, Pw, [23]uint8{0x0f, 0x49}},
	{ACMOVQPS, yml_rl, Pw, [23]uint8{0x0f, 0x4a}},
	{ACMOVWCC, yml_rl, Pq, [23]uint8{0x43}},
	{ACMOVWCS, yml_rl, Pq, [23]uint8{0x42}},
	{ACMOVWEQ, yml_rl, Pq, [23]uint8{0x44}},
	{ACMOVWGE, yml_rl, Pq, [23]uint8{0x4d}},
	{ACMOVWGT, yml_rl, Pq, [23]uint8{0x4f}},
	{ACMOVWHI, yml_rl, Pq, [23]uint8{0x47}},
	{ACMOVWLE, yml_rl, Pq, [23]uint8{0x4e}},
	{ACMOVWLS, yml_rl, Pq, [23]uint8{0x46}},
	{ACMOVWLT, yml_rl, Pq, [23]uint8{0x4c}},
	{ACMOVWMI, yml_rl, Pq, [23]uint8{0x48}},
	{ACMOVWNE, yml_rl, Pq, [23]uint8{0x45}},
	{ACMOVWOC, yml_rl, Pq, [23]uint8{0x41}},
	{ACMOVWOS, yml_rl, Pq, [23]uint8{0x40}},
	{ACMOVWPC, yml_rl, Pq, [23]uint8{0x4b}},
	{ACMOVWPL, yml_rl, Pq, [23]uint8{0x49}},
	{ACMOVWPS, yml_rl, Pq, [23]uint8{0x4a}},
	{ACMPB, ycmpb, Pb, [23]uint8{0x3c, 0x80, 07, 0x38, 0x3a}},
	{ACMPL, ycmpl, Px, [23]uint8{0x83, 07, 0x3d, 0x81, 07, 0x39, 0x3b}},
	{ACMPPD, yxcmpi, Px, [23]uint8{Pe, 0xc2}},
	{ACMPPS, yxcmpi, Pm, [23]uint8{0xc2, 0}},
	{ACMPQ, ycmpl, Pw, [23]uint8{0x83, 07, 0x3d, 0x81, 07, 0x39, 0x3b}},
	{ACMPSB, ynone, Pb, [23]uint8{0xa6}},
	{ACMPSD, yxcmpi, Px, [23]uint8{Pf2, 0xc2}},
	{ACMPSL, ynone, Px, [23]uint8{0xa7}},
	{ACMPSQ, ynone, Pw, [23]uint8{0xa7}},
	{ACMPSS, yxcmpi, Px, [23]uint8{Pf3, 0xc2}},
	{ACMPSW, ynone, Pe, [23]uint8{0xa7}},
	{ACMPW, ycmpl, Pe, [23]uint8{0x83, 07, 0x3d, 0x81, 07, 0x39, 0x3b}},
	{ACOMISD, yxm, Pe, [23]uint8{0x2f}},
	{ACOMISS, yxm, Pm, [23]uint8{0x2f}},
	{ACPUID, ynone, Pm, [23]uint8{0xa2}},
	{ACVTPL2PD, yxcvm2, Px, [23]uint8{Pf3, 0xe6, Pe, 0x2a}},
	{ACVTPL2PS, yxcvm2, Pm, [23]uint8{0x5b, 0, 0x2a, 0}},
	{ACVTPD2PL, yxcvm1, Px, [23]uint8{Pf2, 0xe6, Pe, 0x2d}},
	{ACVTPD2PS, yxm, Pe, [23]uint8{0x5a}},
	{ACVTPS2PL, yxcvm1, Px, [23]uint8{Pe, 0x5b, Pm, 0x2d}},
	{ACVTPS2PD, yxm, Pm, [23]uint8{0x5a}},
	{ACVTSD2SL, yxcvfl, Pf2, [23]uint8{0x2d}},
	{ACVTSD2SQ, yxcvfq, Pw, [23]uint8{Pf2, 0x2d}},
	{ACVTSD2SS, yxm, Pf2, [23]uint8{0x5a}},
	{ACVTSL2SD, yxcvlf, Pf2, [23]uint8{0x2a}},
	{ACVTSQ2SD, yxcvqf, Pw, [23]uint8{Pf2, 0x2a}},
	{ACVTSL2SS, yxcvlf, Pf3, [23]uint8{0x2a}},
	{ACVTSQ2SS, yxcvqf, Pw, [23]uint8{Pf3, 0x2a}},
	{ACVTSS2SD, yxm, Pf3, [23]uint8{0x5a}},
	{ACVTSS2SL, yxcvfl, Pf3, [23]uint8{0x2d}},
	{ACVTSS2SQ, yxcvfq, Pw, [23]uint8{Pf3, 0x2d}},
	{ACVTTPD2PL, yxcvm1, Px, [23]uint8{Pe, 0xe6, Pe, 0x2c}},
	{ACVTTPS2PL, yxcvm1, Px, [23]uint8{Pf3, 0x5b, Pm, 0x2c}},
	{ACVTTSD2SL, yxcvfl, Pf2, [23]uint8{0x2c}},
	{ACVTTSD2SQ, yxcvfq, Pw, [23]uint8{Pf2, 0x2c}},
	{ACVTTSS2SL, yxcvfl, Pf3, [23]uint8{0x2c}},
	{ACVTTSS2SQ, yxcvfq, Pw, [23]uint8{Pf3, 0x2c}},
	{ACWD, ynone, Pe, [23]uint8{0x99}},
	{ACQO, ynone, Pw, [23]uint8{0x99}},
	{ADAA, ynone, P32, [23]uint8{0x27}},
	{ADAS, ynone, P32, [23]uint8{0x2f}},
	{ADECB, yscond, Pb, [23]uint8{0xfe, 01}},
	{ADECL, yincl, Px1, [23]uint8{0x48, 0xff, 01}},
	{ADECQ, yincq, Pw, [23]uint8{0xff, 01}},
	{ADECW, yincq, Pe, [23]uint8{0xff, 01}},
	{ADIVB, ydivb, Pb, [23]uint8{0xf6, 06}},
	{ADIVL, ydivl, Px, [23]uint8{0xf7, 06}},
	{ADIVPD, yxm, Pe, [23]uint8{0x5e}},
	{ADIVPS, yxm, Pm, [23]uint8{0x5e}},
	{ADIVQ, ydivl, Pw, [23]uint8{0xf7, 06}},
	{ADIVSD, yxm, Pf2, [23]uint8{0x5e}},
	{ADIVSS, yxm, Pf3, [23]uint8{0x5e}},
	{ADIVW, ydivl, Pe, [23]uint8{0xf7, 06}},
	{ADPPD, yxshuf, Pq, [23]uint8{0x3a, 0x41, 0}},
	{ADPPS, yxshuf, Pq, [23]uint8{0x3a, 0x40, 0}},
	{AEMMS, ynone, Pm, [23]uint8{0x77}},
	{AEXTRACTPS, yextractps, Pq, [23]uint8{0x3a, 0x17, 0}},
	{AENTER, nil, 0, [23]uint8{}}, /* botch */
	{AFXRSTOR, ysvrs_mo, Pm, [23]uint8{0xae, 01, 0xae, 01}},
	{AFXSAVE, ysvrs_om, Pm, [23]uint8{0xae, 00, 0xae, 00}},
	{AFXRSTOR64, ysvrs_mo, Pw, [23]uint8{0x0f, 0xae, 01, 0x0f, 0xae, 01}},
	{AFXSAVE64, ysvrs_om, Pw, [23]uint8{0x0f, 0xae, 00, 0x0f, 0xae, 00}},
	{AHLT, ynone, Px, [23]uint8{0xf4}},
	{AIDIVB, ydivb, Pb, [23]uint8{0xf6, 07}},
	{AIDIVL, ydivl, Px, [23]uint8{0xf7, 07}},
	{AIDIVQ, ydivl, Pw, [23]uint8{0xf7, 07}},
	{AIDIVW, ydivl, Pe, [23]uint8{0xf7, 07}},
	{AIMULB, ydivb, Pb, [23]uint8{0xf6, 05}},
	{AIMULL, yimul, Px, [23]uint8{0xf7, 05, 0x6b, 0x69, Pm, 0xaf}},
	{AIMULQ, yimul, Pw, [23]uint8{0xf7, 05, 0x6b, 0x69, Pm, 0xaf}},
	{AIMULW, yimul, Pe, [23]uint8{0xf7, 05, 0x6b, 0x69, Pm, 0xaf}},
	{AIMUL3Q, yimul3, Pw, [23]uint8{0x6b, 00}},
	{AINB, yin, Pb, [23]uint8{0xe4, 0xec}},
	{AINCB, yscond, Pb, [23]uint8{0xfe, 00}},
	{AINCL, yincl, Px1, [23]uint8{0x40, 0xff, 00}},
	{AINCQ, yincq, Pw, [23]uint8{0xff, 00}},
	{AINCW, yincq, Pe, [23]uint8{0xff, 00}},
	{AINL, yin, Px, [23]uint8{0xe5, 0xed}},
	{AINSB, ynone, Pb, [23]uint8{0x6c}},
	{AINSL, ynone, Px, [23]uint8{0x6d}},
	{AINSERTPS, yxshuf, Pq, [23]uint8{0x3a, 0x21, 0}},
	{AINSW, ynone, Pe, [23]uint8{0x6d}},
	{AINT, yint, Px, [23]uint8{0xcd}},
	{AINTO, ynone, P32, [23]uint8{0xce}},
	{AINW, yin, Pe, [23]uint8{0xe5, 0xed}},
	{AIRETL, ynone, Px, [23]uint8{0xcf}},
	{AIRETQ, ynone, Pw, [23]uint8{0xcf}},
	{AIRETW, ynone, Pe, [23]uint8{0xcf}},
	{AJCC, yjcond, Px, [23]uint8{0x73, 0x83, 00}},
	{AJCS, yjcond, Px, [23]uint8{0x72, 0x82}},
	{AJCXZL, yloop, Px, [23]uint8{0xe3}},
	{AJCXZW, yloop, Px, [23]uint8{0xe3}},
	{AJCXZQ, yloop, Px, [23]uint8{0xe3}},
	{AJEQ, yjcond, Px, [23]uint8{0x74, 0x84}},
	{AJGE, yjcond, Px, [23]uint8{0x7d, 0x8d}},
	{AJGT, yjcond, Px, [23]uint8{0x7f, 0x8f}},
	{AJHI, yjcond, Px, [23]uint8{0x77, 0x87}},
	{AJLE, yjcond, Px, [23]uint8{0x7e, 0x8e}},
	{AJLS, yjcond, Px, [23]uint8{0x76, 0x86}},
	{AJLT, yjcond, Px, [23]uint8{0x7c, 0x8c}},
	{AJMI, yjcond, Px, [23]uint8{0x78, 0x88}},
	{obj.AJMP, yjmp, Px, [23]uint8{0xff, 04, 0xeb, 0xe9}},
	{AJNE, yjcond, Px, [23]uint8{0x75, 0x85}},
	{AJOC, yjcond, Px, [23]uint8{0x71, 0x81, 00}},
	{AJOS, yjcond, Px, [23]uint8{0x70, 0x80, 00}},
	{AJPC, yjcond, Px, [23]uint8{0x7b, 0x8b}},
	{AJPL, yjcond, Px, [23]uint8{0x79, 0x89}},
	{AJPS, yjcond, Px, [23]uint8{0x7a, 0x8a}},
	{AHADDPD, yxm, Pq, [23]uint8{0x7c}},
	{AHADDPS, yxm, Pf2, [23]uint8{0x7c}},
	{AHSUBPD, yxm, Pq, [23]uint8{0x7d}},
	{AHSUBPS, yxm, Pf2, [23]uint8{0x7d}},
	{ALAHF, ynone, Px, [23]uint8{0x9f}},
	{ALARL, yml_rl, Pm, [23]uint8{0x02}},
	{ALARW, yml_rl, Pq, [23]uint8{0x02}},
	{ALDDQU, ylddqu, Pf2, [23]uint8{0xf0}},
	{ALDMXCSR, ysvrs_mo, Pm, [23]uint8{0xae, 02, 0xae, 02}},
	{ALEAL, ym_rl, Px, [23]uint8{0x8d}},
	{ALEAQ, ym_rl, Pw, [23]uint8{0x8d}},
	{ALEAVEL, ynone, P32, [23]uint8{0xc9}},
	{ALEAVEQ, ynone, Py, [23]uint8{0xc9}},
	{ALEAVEW, ynone, Pe, [23]uint8{0xc9}},
	{ALEAW, ym_rl, Pe, [23]uint8{0x8d}},
	{ALOCK, ynone, Px, [23]uint8{0xf0}},
	{ALODSB, ynone, Pb, [23]uint8{0xac}},
	{ALODSL, ynone, Px, [23]uint8{0xad}},
	{ALODSQ, ynone, Pw, [23]uint8{0xad}},
	{ALODSW, ynone, Pe, [23]uint8{0xad}},
	{ALONG, ybyte, Px, [23]uint8{4}},
	{ALOOP, yloop, Px, [23]uint8{0xe2}},
	{ALOOPEQ, yloop, Px, [23]uint8{0xe1}},
	{ALOOPNE, yloop, Px, [23]uint8{0xe0}},
	{ALSLL, yml_rl, Pm, [23]uint8{0x03}},
	{ALSLW, yml_rl, Pq, [23]uint8{0x03}},
	{AMASKMOVOU, yxr, Pe, [23]uint8{0xf7}},
	{AMASKMOVQ, ymr, Pm, [23]uint8{0xf7}},
	{AMAXPD, yxm, Pe, [23]uint8{0x5f}},
	{AMAXPS, yxm, Pm, [23]uint8{0x5f}},
	{AMAXSD, yxm, Pf2, [23]uint8{0x5f}},
	{AMAXSS, yxm, Pf3, [23]uint8{0x5f}},
	{AMINPD, yxm, Pe, [23]uint8{0x5d}},
	{AMINPS, yxm, Pm, [23]uint8{0x5d}},
	{AMINSD, yxm, Pf2, [23]uint8{0x5d}},
	{AMINSS, yxm, Pf3, [23]uint8{0x5d}},
	{AMOVAPD, yxmov, Pe, [23]uint8{0x28, 0x29}},
	{AMOVAPS, yxmov, Pm, [23]uint8{0x28, 0x29}},
	{AMOVB, ymovb, Pb, [23]uint8{0x88, 0x8a, 0xb0, 0xc6, 00}},
	{AMOVBLSX, ymb_rl, Pm, [23]uint8{0xbe}},
	{AMOVBLZX, ymb_rl, Pm, [23]uint8{0xb6}},
	{AMOVBQSX, ymb_rl, Pw, [23]uint8{0x0f, 0xbe}},
	{AMOVBQZX, ymb_rl, Pm, [23]uint8{0xb6}},
	{AMOVBWSX, ymb_rl, Pq, [23]uint8{0xbe}},
	{AMOVBWZX, ymb_rl, Pq, [23]uint8{0xb6}},
	{AMOVO, yxmov, Pe, [23]uint8{0x6f, 0x7f}},
	{AMOVOU, yxmov, Pf3, [23]uint8{0x6f, 0x7f}},
	{AMOVHLPS, yxr, Pm, [23]uint8{0x12}},
	{AMOVHPD, yxmov, Pe, [23]uint8{0x16, 0x17}},
	{AMOVHPS, yxmov, Pm, [23]uint8{0x16, 0x17}},
	{AMOVL, ymovl, Px, [23]uint8{0x89, 0x8b, 0xb8, 0xc7, 00, 0x6e, 0x7e, Pe, 0x6e, Pe, 0x7e, 0}},
	{AMOVLHPS, yxr, Pm, [23]uint8{0x16}},
	{AMOVLPD, yxmov, Pe, [23]uint8{0x12, 0x13}},
	{AMOVLPS, yxmov, Pm, [23]uint8{0x12, 0x13}},
	{AMOVLQSX, yml_rl, Pw, [23]uint8{0x63}},
	{AMOVLQZX, yml_rl, Px, [23]uint8{0x8b}},
	{AMOVMSKPD, yxrrl, Pq, [23]uint8{0x50}},
	{AMOVMSKPS, yxrrl, Pm, [23]uint8{0x50}},
	{AMOVNTO, yxr_ml, Pe, [23]uint8{0xe7}},
	{AMOVNTDQA, ylddqu, Pq4, [23]uint8{0x2a}},
	{AMOVNTPD, yxr_ml, Pe, [23]uint8{0x2b}},
	{AMOVNTPS, yxr_ml, Pm, [23]uint8{0x2b}},
	{AMOVNTQ, ymr_ml, Pm, [23]uint8{0xe7}},
	{AMOVQ, ymovq, Pw8, [23]uint8{0x6f, 0x7f, Pf2, 0xd6, Pf3, 0x7e, Pe, 0xd6, 0x89, 0x8b, 0xc7, 00, 0xb8, 0xc7, 00, 0x6e, 0x7e, Pe, 0x6e, Pe, 0x7e, 0}},
	{AMOVQOZX, ymrxr, Pf3, [23]uint8{0xd6, 0x7e}},
	{AMOVSB, ynone, Pb, [23]uint8{0xa4}},
	{AMOVSD, yxmov, Pf2, [23]uint8{0x10, 0x11}},
	{AMOVSL, ynone, Px, [23]uint8{0xa5}},
	{AMOVSQ, ynone, Pw, [23]uint8{0xa5}},
	{AMOVSS, yxmov, Pf3, [23]uint8{0x10, 0x11}},
	{AMOVSW, ynone, Pe, [23]uint8{0xa5}},
	{AMOVUPD, yxmov, Pe, [23]uint8{0x10, 0x11}},
	{AMOVUPS, yxmov, Pm, [23]uint8{0x10, 0x11}},
	{AMOVW, ymovw, Pe, [23]uint8{0x89, 0x8b, 0xb8, 0xc7, 00, 0}},
	{AMOVWLSX, yml_rl, Pm, [23]uint8{0xbf}},
	{AMOVWLZX, yml_rl, Pm, [23]uint8{0xb7}},
	{AMOVWQSX, yml_rl, Pw, [23]uint8{0x0f, 0xbf}},
	{AMOVWQZX, yml_rl, Pw, [23]uint8{0x0f, 0xb7}},
	{AMPSADBW, yxshuf, Pq, [23]uint8{0x3a, 0x42, 0}},
	{AMULB, ydivb, Pb, [23]uint8{0xf6, 04}},
	{AMULL, ydivl, Px, [23]uint8{0xf7, 04}},
	{AMULPD, yxm, Pe, [23]uint8{0x59}},
	{AMULPS, yxm, Ym, [23]uint8{0x59}},
	{AMULQ, ydivl, Pw, [23]uint8{0xf7, 04}},
	{AMULSD, yxm, Pf2, [23]uint8{0x59}},
	{AMULSS, yxm, Pf3, [23]uint8{0x59}},
	{AMULW, ydivl, Pe, [23]uint8{0xf7, 04}},
	{ANEGB, yscond, Pb, [23]uint8{0xf6, 03}},
	{ANEGL, yscond, Px, [23]uint8{0xf7, 03}},
	{ANEGQ, yscond, Pw, [23]uint8{0xf7, 03}},
	{ANEGW, yscond, Pe, [23]uint8{0xf7, 03}},
	{obj.ANOP, ynop, Px, [23]uint8{0, 0}},
	{ANOTB, yscond, Pb, [23]uint8{0xf6, 02}},
	{ANOTL, yscond, Px, [23]uint8{0xf7, 02}}, // TODO(rsc): yscond is wrong here.
	{ANOTQ, yscond, Pw, [23]uint8{0xf7, 02}},
	{ANOTW, yscond, Pe, [23]uint8{0xf7, 02}},
	{AORB, yxorb, Pb, [23]uint8{0x0c, 0x80, 01, 0x08, 0x0a}},
	{AORL, yaddl, Px, [23]uint8{0x83, 01, 0x0d, 0x81, 01, 0x09, 0x0b}},
	{AORPD, yxm, Pq, [23]uint8{0x56}},
	{AORPS, yxm, Pm, [23]uint8{0x56}},
	{AORQ, yaddl, Pw, [23]uint8{0x83, 01, 0x0d, 0x81, 01, 0x09, 0x0b}},
	{AORW, yaddl, Pe, [23]uint8{0x83, 01, 0x0d, 0x81, 01, 0x09, 0x0b}},
	{AOUTB, yin, Pb, [23]uint8{0xe6, 0xee}},
	{AOUTL, yin, Px, [23]uint8{0xe7, 0xef}},
	{AOUTSB, ynone, Pb, [23]uint8{0x6e}},
	{AOUTSL, ynone, Px, [23]uint8{0x6f}},
	{AOUTSW, ynone, Pe, [23]uint8{0x6f}},
	{AOUTW, yin, Pe, [23]uint8{0xe7, 0xef}},
	{APABSB, yxm_q4, Pq4, [23]uint8{0x1c}},
	{APABSD, yxm_q4, Pq4, [23]uint8{0x1e}},
	{APABSW, yxm_q4, Pq4, [23]uint8{0x1d}},
	{APACKSSLW, ymm, Py1, [23]uint8{0x6b, Pe, 0x6b}},
	{APACKSSWB, ymm, Py1, [23]uint8{0x63, Pe, 0x63}},
	{APACKUSDW, yxm_q4, Pq4, [23]uint8{0x2b}},
	{APACKUSWB, ymm, Py1, [23]uint8{0x67, Pe, 0x67}},
	{APADDB, ymm, Py1, [23]uint8{0xfc, Pe, 0xfc}},
	{APADDL, ymm, Py1, [23]uint8{0xfe, Pe, 0xfe}},
	{APADDQ, yxm, Pe, [23]uint8{0xd4}},
	{APADDSB, ymm, Py1, [23]uint8{0xec, Pe, 0xec}},
	{APADDSW, ymm, Py1, [23]uint8{0xed, Pe, 0xed}},
	{APADDUSB, ymm, Py1, [23]uint8{0xdc, Pe, 0xdc}},
	{APADDUSW, ymm, Py1, [23]uint8{0xdd, Pe, 0xdd}},
	{APADDW, ymm, Py1, [23]uint8{0xfd, Pe, 0xfd}},
	{APALIGNR, ypalignr, Pq, [23]uint8{0x3a, 0x0f}},
	{APAND, ymm, Py1, [23]uint8{0xdb, Pe, 0xdb}},
	{APANDN, ymm, Py1, [23]uint8{0xdf, Pe, 0xdf}},
	{APAUSE, ynone, Px, [23]uint8{0xf3, 0x90}},
	{APAVGB, ymm, Py1, [23]uint8{0xe0, Pe, 0xe0}},
	{APAVGW, ymm, Py1, [23]uint8{0xe3, Pe, 0xe3}},
	{APBLENDW, yxshuf, Pq, [23]uint8{0x3a, 0x0e, 0}},
	{APCMPEQB, ymm, Py1, [23]uint8{0x74, Pe, 0x74}},
	{APCMPEQL, ymm, Py1, [23]uint8{0x76, Pe, 0x76}},
	{APCMPEQQ, yxm_q4, Pq4, [23]uint8{0x29}},
	{APCMPEQW, ymm, Py1, [23]uint8{0x75, Pe, 0x75}},
	{APCMPGTB, ymm, Py1, [23]uint8{0x64, Pe, 0x64}},
	{APCMPGTL, ymm, Py1, [23]uint8{0x66, Pe, 0x66}},
	{APCMPGTQ, yxm_q4, Pq4, [23]uint8{0x37}},
	{APCMPGTW, ymm, Py1, [23]uint8{0x65, Pe, 0x65}},
	{APCMPISTRI, yxshuf, Pq, [23]uint8{0x3a, 0x63, 0}},
	{APCMPISTRM, yxshuf, Pq, [23]uint8{0x3a, 0x62, 0}},
	{APEXTRW, yextrw, Pq, [23]uint8{0xc5, 00}},
	{APEXTRB, yextr, Pq, [23]uint8{0x3a, 0x14, 00}},
	{APEXTRD, yextr, Pq, [23]uint8{0x3a, 0x16, 00}},
	{APEXTRQ, yextr, Pq3, [23]uint8{0x3a, 0x16, 00}},
	{APHADDD, ymmxmm0f38, Px, [23]uint8{0x0F, 0x38, 0x02, 0, 0x66, 0x0F, 0x38, 0x02, 0}},
	{APHADDSW, yxm_q4, Pq4, [23]uint8{0x03}},
	{APHADDW, yxm_q4, Pq4, [23]uint8{0x01}},
	{APHMINPOSUW, yxm_q4, Pq4, [23]uint8{0x41}},
	{APHSUBD, yxm_q4, Pq4, [23]uint8{0x06}},
	{APHSUBSW, yxm_q4, Pq4, [23]uint8{0x07}},
	{APHSUBW, yxm_q4, Pq4, [23]uint8{0x05}},
	{APINSRW, yinsrw, Pq, [23]uint8{0xc4, 00}},
	{APINSRB, yinsr, Pq, [23]uint8{0x3a, 0x20, 00}},
	{APINSRD, yinsr, Pq, [23]uint8{0x3a, 0x22, 00}},
	{APINSRQ, yinsr, Pq3, [23]uint8{0x3a, 0x22, 00}},
	{APMADDUBSW, yxm_q4, Pq4, [23]uint8{0x04}},
	{APMADDWL, ymm, Py1, [23]uint8{0xf5, Pe, 0xf5}},
	{APMAXSB, yxm_q4, Pq4, [23]uint8{0x3c}},
	{APMAXSD, yxm_q4, Pq4, [23]uint8{0x3d}},
	{APMAXSW, yxm, Pe, [23]uint8{0xee}},
	{APMAXUB, yxm, Pe, [23]uint8{0xde}},
	{APMAXUD, yxm_q4, Pq4, [23]uint8{0x3f}},
	{APMAXUW, yxm_q4, Pq4, [23]uint8{0x3e}},
	{APMINSB, yxm_q4, Pq4, [23]uint8{0x38}},
	{APMINSD, yxm_q4, Pq4, [23]uint8{0x39}},
	{APMINSW, yxm, Pe, [23]uint8{0xea}},
	{APMINUB, yxm, Pe, [23]uint8{0xda}},
	{APMINUD, yxm_q4, Pq4, [23]uint8{0x3b}},
	{APMINUW, yxm_q4, Pq4, [23]uint8{0x3a}},
	{APMOVMSKB, ymskb, Px, [23]uint8{Pe, 0xd7, 0xd7}},
	{APMOVSXBD, yxm_q4, Pq4, [23]uint8{0x21}},
	{APMOVSXBQ, yxm_q4, Pq4, [23]uint8{0x22}},
	{APMOVSXBW, yxm_q4, Pq4, [23]uint8{0x20}},
	{APMOVSXDQ, yxm_q4, Pq4, [23]uint8{0x25}},
	{APMOVSXWD, yxm_q4, Pq4, [23]uint8{0x23}},
	{APMOVSXWQ, yxm_q4, Pq4, [23]uint8{0x24}},
	{APMOVZXBD, yxm_q4, Pq4, [23]uint8{0x31}},
	{APMOVZXBQ, yxm_q4, Pq4, [23]uint8{0x32}},
	{APMOVZXBW, yxm_q4, Pq4, [23]uint8{0x30}},
	{APMOVZXDQ, yxm_q4, Pq4, [23]uint8{0x35}},
	{APMOVZXWD, yxm_q4, Pq4, [23]uint8{0x33}},
	{APMOVZXWQ, yxm_q4, Pq4, [23]uint8{0x34}},
	{APMULDQ, yxm_q4, Pq4, [23]uint8{0x28}},
	{APMULHRSW, yxm_q4, Pq4, [23]uint8{0x0b}},
	{APMULHUW, ymm, Py1, [23]uint8{0xe4, Pe, 0xe4}},
	{APMULHW, ymm, Py1, [23]uint8{0xe5, Pe, 0xe5}},
	{APMULLD, yxm_q4, Pq4, [23]uint8{0x40}},
	{APMULLW, ymm, Py1, [23]uint8{0xd5, Pe, 0xd5}},
	{APMULULQ, ymm, Py1, [23]uint8{0xf4, Pe, 0xf4}},
	{APOPAL, ynone, P32, [23]uint8{0x61}},
	{APOPAW, ynone, Pe, [23]uint8{0x61}},
	{APOPCNTW, yml_rl, Pef3, [23]uint8{0xb8}},
	{APOPCNTL, yml_rl, Pf3, [23]uint8{0xb8}},
	{APOPCNTQ, yml_rl, Pfw, [23]uint8{0xb8}},
	{APOPFL, ynone, P32, [23]uint8{0x9d}},
	{APOPFQ, ynone, Py, [23]uint8{0x9d}},
	{APOPFW, ynone, Pe, [23]uint8{0x9d}},
	{APOPL, ypopl, P32, [23]uint8{0x58, 0x8f, 00}},
	{APOPQ, ypopl, Py, [23]uint8{0x58, 0x8f, 00}},
	{APOPW, ypopl, Pe, [23]uint8{0x58, 0x8f, 00}},
	{APOR, ymm, Py1, [23]uint8{0xeb, Pe, 0xeb}},
	{APSADBW, yxm, Pq, [23]uint8{0xf6}},
	{APSHUFHW, yxshuf, Pf3, [23]uint8{0x70, 00}},
	{APSHUFL, yxshuf, Pq, [23]uint8{0x70, 00}},
	{APSHUFLW, yxshuf, Pf2, [23]uint8{0x70, 00}},
	{APSHUFW, ymshuf, Pm, [23]uint8{0x70, 00}},
	{APSHUFB, ymshufb, Pq, [23]uint8{0x38, 0x00}},
	{APSIGNB, yxm_q4, Pq4, [23]uint8{0x08}},
	{APSIGND, yxm_q4, Pq4, [23]uint8{0x0a}},
	{APSIGNW, yxm_q4, Pq4, [23]uint8{0x09}},
	{APSLLO, ypsdq, Pq, [23]uint8{0x73, 07}},
	{APSLLL, yps, Py3, [23]uint8{0xf2, 0x72, 06, Pe, 0xf2, Pe, 0x72, 06}},
	{APSLLQ, yps, Py3, [23]uint8{0xf3, 0x73, 06, Pe, 0xf3, Pe, 0x73, 06}},
	{APSLLW, yps, Py3, [23]uint8{0xf1, 0x71, 06, Pe, 0xf1, Pe, 0x71, 06}},
	{APSRAL, yps, Py3, [23]uint8{0xe2, 0x72, 04, Pe, 0xe2, Pe, 0x72, 04}},
	{APSRAW, yps, Py3, [23]uint8{0xe1, 0x71, 04, Pe, 0xe1, Pe, 0x71, 04}},
	{APSRLO, ypsdq, Pq, [23]uint8{0x73, 03}},
	{APSRLL, yps, Py3, [23]uint8{0xd2, 0x72, 02, Pe, 0xd2, Pe, 0x72, 02}},
	{APSRLQ, yps, Py3, [23]uint8{0xd3, 0x73, 02, Pe, 0xd3, Pe, 0x73, 02}},
	{APSRLW, yps, Py3, [23]uint8{0xd1, 0x71, 02, Pe, 0xd1, Pe, 0x71, 02}},
	{APSUBB, yxm, Pe, [23]uint8{0xf8}},
	{APSUBL, yxm, Pe, [23]uint8{0xfa}},
	{APSUBQ, yxm, Pe, [23]uint8{0xfb}},
	{APSUBSB, yxm, Pe, [23]uint8{0xe8}},
	{APSUBSW, yxm, Pe, [23]uint8{0xe9}},
	{APSUBUSB, yxm, Pe, [23]uint8{0xd8}},
	{APSUBUSW, yxm, Pe, [23]uint8{0xd9}},
	{APSUBW, yxm, Pe, [23]uint8{0xf9}},
	{APTEST, yxm_q4, Pq4, [23]uint8{0x17}},
	{APUNPCKHBW, ymm, Py1, [23]uint8{0x68, Pe, 0x68}},
	{APUNPCKHLQ, ymm, Py1, [23]uint8{0x6a, Pe, 0x6a}},
	{APUNPCKHQDQ, yxm, Pe, [23]uint8{0x6d}},
	{APUNPCKHWL, ymm, Py1, [23]uint8{0x69, Pe, 0x69}},
	{APUNPCKLBW, ymm, Py1, [23]uint8{0x60, Pe, 0x60}},
	{APUNPCKLLQ, ymm, Py1, [23]uint8{0x62, Pe, 0x62}},
	{APUNPCKLQDQ, yxm, Pe, [23]uint8{0x6c}},
	{APUNPCKLWL, ymm, Py1, [23]uint8{0x61, Pe, 0x61}},
	{APUSHAL, ynone, P32, [23]uint8{0x60}},
	{APUSHAW, ynone, Pe, [23]uint8{0x60}},
	{APUSHFL, ynone, P32, [23]uint8{0x9c}},
	{APUSHFQ, ynone, Py, [23]uint8{0x9c}},
	{APUSHFW, ynone, Pe, [23]uint8{0x9c}},
	{APUSHL, ypushl, P32, [23]uint8{0x50, 0xff, 06, 0x6a, 0x68}},
	{APUSHQ, ypushl, Py, [23]uint8{0x50, 0xff, 06, 0x6a, 0x68}},
	{APUSHW, ypushl, Pe, [23]uint8{0x50, 0xff, 06, 0x6a, 0x68}},
	{APXOR, ymm, Py1, [23]uint8{0xef, Pe, 0xef}},
	{AQUAD, ybyte, Px, [23]uint8{8}},
	{ARCLB, yshb, Pb, [23]uint8{0xd0, 02, 0xc0, 02, 0xd2, 02}},
	{ARCLL, yshl, Px, [23]uint8{0xd1, 02, 0xc1, 02, 0xd3, 02, 0xd3, 02}},
	{ARCLQ, yshl, Pw, [23]uint8{0xd1, 02, 0xc1, 02, 0xd3, 02, 0xd3, 02}},
	{ARCLW, yshl, Pe, [23]uint8{0xd1, 02, 0xc1, 02, 0xd3, 02, 0xd3, 02}},
	{ARCPPS, yxm, Pm, [23]uint8{0x53}},
	{ARCPSS, yxm, Pf3, [23]uint8{0x53}},
	{ARCRB, yshb, Pb, [23]uint8{0xd0, 03, 0xc0, 03, 0xd2, 03}},
	{ARCRL, yshl, Px, [23]uint8{0xd1, 03, 0xc1, 03, 0xd3, 03, 0xd3, 03}},
	{ARCRQ, yshl, Pw, [23]uint8{0xd1, 03, 0xc1, 03, 0xd3, 03, 0xd3, 03}},
	{ARCRW, yshl, Pe, [23]uint8{0xd1, 03, 0xc1, 03, 0xd3, 03, 0xd3, 03}},
	{AREP, ynone, Px, [23]uint8{0xf3}},
	{AREPN, ynone, Px, [23]uint8{0xf2}},
	{obj.ARET, ynone, Px, [23]uint8{0xc3}},
	{ARETFW, yret, Pe, [23]uint8{0xcb, 0xca}},
	{ARETFL, yret, Px, [23]uint8{0xcb, 0xca}},
	{ARETFQ, yret, Pw, [23]uint8{0xcb, 0xca}},
	{AROLB, yshb, Pb, [23]uint8{0xd0, 00, 0xc0, 00, 0xd2, 00}},
	{AROLL, yshl, Px, [23]uint8{0xd1, 00, 0xc1, 00, 0xd3, 00, 0xd3, 00}},
	{AROLQ, yshl, Pw, [23]uint8{0xd1, 00, 0xc1, 00, 0xd3, 00, 0xd3, 00}},
	{AROLW, yshl, Pe, [23]uint8{0xd1, 00, 0xc1, 00, 0xd3, 00, 0xd3, 00}},
	{ARORB, yshb, Pb, [23]uint8{0xd0, 01, 0xc0, 01, 0xd2, 01}},
	{ARORL, yshl, Px, [23]uint8{0xd1, 01, 0xc1, 01, 0xd3, 01, 0xd3, 01}},
	{ARORQ, yshl, Pw, [23]uint8{0xd1, 01, 0xc1, 01, 0xd3, 01, 0xd3, 01}},
	{ARORW, yshl, Pe, [23]uint8{0xd1, 01, 0xc1, 01, 0xd3, 01, 0xd3, 01}},
	{ARSQRTPS, yxm, Pm, [23]uint8{0x52}},
	{ARSQRTSS, yxm, Pf3, [23]uint8{0x52}},
	{ASAHF, ynone, Px1, [23]uint8{0x9e, 00, 0x86, 0xe0, 0x50, 0x9d}}, /* XCHGB AH,AL; PUSH AX; POPFL */
	{ASALB, yshb, Pb, [23]uint8{0xd0, 04, 0xc0, 04, 0xd2, 04}},
	{ASALL, yshl, Px, [23]uint8{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASALQ, yshl, Pw, [23]uint8{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASALW, yshl, Pe, [23]uint8{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASARB, yshb, Pb, [23]uint8{0xd0, 07, 0xc0, 07, 0xd2, 07}},
	{ASARL, yshl, Px, [23]uint8{0xd1, 07, 0xc1, 07, 0xd3, 07, 0xd3, 07}},
	{ASARQ, yshl, Pw, [23]uint8{0xd1, 07, 0xc1, 07, 0xd3, 07, 0xd3, 07}},
	{ASARW, yshl, Pe, [23]uint8{0xd1, 07, 0xc1, 07, 0xd3, 07, 0xd3, 07}},
	{ASBBB, yxorb, Pb, [23]uint8{0x1c, 0x80, 03, 0x18, 0x1a}},
	{ASBBL, yaddl, Px, [23]uint8{0x83, 03, 0x1d, 0x81, 03, 0x19, 0x1b}},
	{ASBBQ, yaddl, Pw, [23]uint8{0x83, 03, 0x1d, 0x81, 03, 0x19, 0x1b}},
	{ASBBW, yaddl, Pe, [23]uint8{0x83, 03, 0x1d, 0x81, 03, 0x19, 0x1b}},
	{ASCASB, ynone, Pb, [23]uint8{0xae}},
	{ASCASL, ynone, Px, [23]uint8{0xaf}},
	{ASCASQ, ynone, Pw, [23]uint8{0xaf}},
	{ASCASW, ynone, Pe, [23]uint8{0xaf}},
	{ASETCC, yscond, Pb, [23]uint8{0x0f, 0x93, 00}},
	{ASETCS, yscond, Pb, [23]uint8{0x0f, 0x92, 00}},
	{ASETEQ, yscond, Pb, [23]uint8{0x0f, 0x94, 00}},
	{ASETGE, yscond, Pb, [23]uint8{0x0f, 0x9d, 00}},
	{ASETGT, yscond, Pb, [23]uint8{0x0f, 0x9f, 00}},
	{ASETHI, yscond, Pb, [23]uint8{0x0f, 0x97, 00}},
	{ASETLE, yscond, Pb, [23]uint8{0x0f, 0x9e, 00}},
	{ASETLS, yscond, Pb, [23]uint8{0x0f, 0x96, 00}},
	{ASETLT, yscond, Pb, [23]uint8{0x0f, 0x9c, 00}},
	{ASETMI, yscond, Pb, [23]uint8{0x0f, 0x98, 00}},
	{ASETNE, yscond, Pb, [23]uint8{0x0f, 0x95, 00}},
	{ASETOC, yscond, Pb, [23]uint8{0x0f, 0x91, 00}},
	{ASETOS, yscond, Pb, [23]uint8{0x0f, 0x90, 00}},
	{ASETPC, yscond, Pb, [23]uint8{0x0f, 0x9b, 00}},
	{ASETPL, yscond, Pb, [23]uint8{0x0f, 0x99, 00}},
	{ASETPS, yscond, Pb, [23]uint8{0x0f, 0x9a, 00}},
	{ASHLB, yshb, Pb, [23]uint8{0xd0, 04, 0xc0, 04, 0xd2, 04}},
	{ASHLL, yshl, Px, [23]uint8{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASHLQ, yshl, Pw, [23]uint8{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASHLW, yshl, Pe, [23]uint8{0xd1, 04, 0xc1, 04, 0xd3, 04, 0xd3, 04}},
	{ASHRB, yshb, Pb, [23]uint8{0xd0, 05, 0xc0, 05, 0xd2, 05}},
	{ASHRL, yshl, Px, [23]uint8{0xd1, 05, 0xc1, 05, 0xd3, 05, 0xd3, 05}},
	{ASHRQ, yshl, Pw, [23]uint8{0xd1, 05, 0xc1, 05, 0xd3, 05, 0xd3, 05}},
	{ASHRW, yshl, Pe, [23]uint8{0xd1, 05, 0xc1, 05, 0xd3, 05, 0xd3, 05}},
	{ASHUFPD, yxshuf, Pq, [23]uint8{0xc6, 00}},
	{ASHUFPS, yxshuf, Pm, [23]uint8{0xc6, 00}},
	{ASQRTPD, yxm, Pe, [23]uint8{0x51}},
	{ASQRTPS, yxm, Pm, [23]uint8{0x51}},
	{ASQRTSD, yxm, Pf2, [23]uint8{0x51}},
	{ASQRTSS, yxm, Pf3, [23]uint8{0x51}},
	{ASTC, ynone, Px, [23]uint8{0xf9}},
	{ASTD, ynone, Px, [23]uint8{0xfd}},
	{ASTI, ynone, Px, [23]uint8{0xfb}},
	{ASTMXCSR, ysvrs_om, Pm, [23]uint8{0xae, 03, 0xae, 03}},
	{ASTOSB, ynone, Pb, [23]uint8{0xaa}},
	{ASTOSL, ynone, Px, [23]uint8{0xab}},
	{ASTOSQ, ynone, Pw, [23]uint8{0xab}},
	{ASTOSW, ynone, Pe, [23]uint8{0xab}},
	{ASUBB, yxorb, Pb, [23]uint8{0x2c, 0x80, 05, 0x28, 0x2a}},
	{ASUBL, yaddl, Px, [23]uint8{0x83, 05, 0x2d, 0x81, 05, 0x29, 0x2b}},
	{ASUBPD, yxm, Pe, [23]uint8{0x5c}},
	{ASUBPS, yxm, Pm, [23]uint8{0x5c}},
	{ASUBQ, yaddl, Pw, [23]uint8{0x83, 05, 0x2d, 0x81, 05, 0x29, 0x2b}},
	{ASUBSD, yxm, Pf2, [23]uint8{0x5c}},
	{ASUBSS, yxm, Pf3, [23]uint8{0x5c}},
	{ASUBW, yaddl, Pe, [23]uint8{0x83, 05, 0x2d, 0x81, 05, 0x29, 0x2b}},
	{ASWAPGS, ynone, Pm, [23]uint8{0x01, 0xf8}},
	{ASYSCALL, ynone, Px, [23]uint8{0x0f, 0x05}}, /* fast syscall */
	{ATESTB, yxorb, Pb, [23]uint8{0xa8, 0xf6, 00, 0x84, 0x84}},
	{ATESTL, ytestl, Px, [23]uint8{0xa9, 0xf7, 00, 0x85, 0x85}},
	{ATESTQ, ytestl, Pw, [23]uint8{0xa9, 0xf7, 00, 0x85, 0x85}},
	{ATESTW, ytestl, Pe, [23]uint8{0xa9, 0xf7, 00, 0x85, 0x85}},
	{obj.ATEXT, ytext, Px, [23]uint8{}},
	{AUCOMISD, yxm, Pe, [23]uint8{0x2e}},
	{AUCOMISS, yxm, Pm, [23]uint8{0x2e}},
	{AUNPCKHPD, yxm, Pe, [23]uint8{0x15}},
	{AUNPCKHPS, yxm, Pm, [23]uint8{0x15}},
	{AUNPCKLPD, yxm, Pe, [23]uint8{0x14}},
	{AUNPCKLPS, yxm, Pm, [23]uint8{0x14}},
	{AVERR, ydivl, Pm, [23]uint8{0x00, 04}},
	{AVERW, ydivl, Pm, [23]uint8{0x00, 05}},
	{AWAIT, ynone, Px, [23]uint8{0x9b}},
	{AWORD, ybyte, Px, [23]uint8{2}},
	{AXCHGB, yml_mb, Pb, [23]uint8{0x86, 0x86}},
	{AXCHGL, yxchg, Px, [23]uint8{0x90, 0x90, 0x87, 0x87}},
	{AXCHGQ, yxchg, Pw, [23]uint8{0x90, 0x90, 0x87, 0x87}},
	{AXCHGW, yxchg, Pe, [23]uint8{0x90, 0x90, 0x87, 0x87}},
	{AXLAT, ynone, Px, [23]uint8{0xd7}},
	{AXORB, yxorb, Pb, [23]uint8{0x34, 0x80, 06, 0x30, 0x32}},
	{AXORL, yaddl, Px, [23]uint8{0x83, 06, 0x35, 0x81, 06, 0x31, 0x33}},
	{AXORPD, yxm, Pe, [23]uint8{0x57}},
	{AXORPS, yxm, Pm, [23]uint8{0x57}},
	{AXORQ, yaddl, Pw, [23]uint8{0x83, 06, 0x35, 0x81, 06, 0x31, 0x33}},
	{AXORW, yaddl, Pe, [23]uint8{0x83, 06, 0x35, 0x81, 06, 0x31, 0x33}},
	{AFMOVB, yfmvx, Px, [23]uint8{0xdf, 04}},
	{AFMOVBP, yfmvp, Px, [23]uint8{0xdf, 06}},
	{AFMOVD, yfmvd, Px, [23]uint8{0xdd, 00, 0xdd, 02, 0xd9, 00, 0xdd, 02}},
	{AFMOVDP, yfmvdp, Px, [23]uint8{0xdd, 03, 0xdd, 03}},
	{AFMOVF, yfmvf, Px, [23]uint8{0xd9, 00, 0xd9, 02}},
	{AFMOVFP, yfmvp, Px, [23]uint8{0xd9, 03}},
	{AFMOVL, yfmvf, Px, [23]uint8{0xdb, 00, 0xdb, 02}},
	{AFMOVLP, yfmvp, Px, [23]uint8{0xdb, 03}},
	{AFMOVV, yfmvx, Px, [23]uint8{0xdf, 05}},
	{AFMOVVP, yfmvp, Px, [23]uint8{0xdf, 07}},
	{AFMOVW, yfmvf, Px, [23]uint8{0xdf, 00, 0xdf, 02}},
	{AFMOVWP, yfmvp, Px, [23]uint8{0xdf, 03}},
	{AFMOVX, yfmvx, Px, [23]uint8{0xdb, 05}},
	{AFMOVXP, yfmvp, Px, [23]uint8{0xdb, 07}},
	{AFCMOVCC, yfcmv, Px, [23]uint8{0xdb, 00}},
	{AFCMOVCS, yfcmv, Px, [23]uint8{0xda, 00}},
	{AFCMOVEQ, yfcmv, Px, [23]uint8{0xda, 01}},
	{AFCMOVHI, yfcmv, Px, [23]uint8{0xdb, 02}},
	{AFCMOVLS, yfcmv, Px, [23]uint8{0xda, 02}},
	{AFCMOVNE, yfcmv, Px, [23]uint8{0xdb, 01}},
	{AFCMOVNU, yfcmv, Px, [23]uint8{0xdb, 03}},
	{AFCMOVUN, yfcmv, Px, [23]uint8{0xda, 03}},
	{AFCOMD, yfadd, Px, [23]uint8{0xdc, 02, 0xd8, 02, 0xdc, 02}},  /* botch */
	{AFCOMDP, yfadd, Px, [23]uint8{0xdc, 03, 0xd8, 03, 0xdc, 03}}, /* botch */
	{AFCOMDPP, ycompp, Px, [23]uint8{0xde, 03}},
	{AFCOMF, yfmvx, Px, [23]uint8{0xd8, 02}},
	{AFCOMFP, yfmvx, Px, [23]uint8{0xd8, 03}},
	{AFCOMI, yfmvx, Px, [23]uint8{0xdb, 06}},
	{AFCOMIP, yfmvx, Px, [23]uint8{0xdf, 06}},
	{AFCOML, yfmvx, Px, [23]uint8{0xda, 02}},
	{AFCOMLP, yfmvx, Px, [23]uint8{0xda, 03}},
	{AFCOMW, yfmvx, Px, [23]uint8{0xde, 02}},
	{AFCOMWP, yfmvx, Px, [23]uint8{0xde, 03}},
	{AFUCOM, ycompp, Px, [23]uint8{0xdd, 04}},
	{AFUCOMI, ycompp, Px, [23]uint8{0xdb, 05}},
	{AFUCOMIP, ycompp, Px, [23]uint8{0xdf, 05}},
	{AFUCOMP, ycompp, Px, [23]uint8{0xdd, 05}},
	{AFUCOMPP, ycompp, Px, [23]uint8{0xda, 13}},
	{AFADDDP, ycompp, Px, [23]uint8{0xde, 00}},
	{AFADDW, yfmvx, Px, [23]uint8{0xde, 00}},
	{AFADDL, yfmvx, Px, [23]uint8{0xda, 00}},
	{AFADDF, yfmvx, Px, [23]uint8{0xd8, 00}},
	{AFADDD, yfadd, Px, [23]uint8{0xdc, 00, 0xd8, 00, 0xdc, 00}},
	{AFMULDP, ycompp, Px, [23]uint8{0xde, 01}},
	{AFMULW, yfmvx, Px, [23]uint8{0xde, 01}},
	{AFMULL, yfmvx, Px, [23]uint8{0xda, 01}},
	{AFMULF, yfmvx, Px, [23]uint8{0xd8, 01}},
	{AFMULD, yfadd, Px, [23]uint8{0xdc, 01, 0xd8, 01, 0xdc, 01}},
	{AFSUBDP, ycompp, Px, [23]uint8{0xde, 05}},
	{AFSUBW, yfmvx, Px, [23]uint8{0xde, 04}},
	{AFSUBL, yfmvx, Px, [23]uint8{0xda, 04}},
	{AFSUBF, yfmvx, Px, [23]uint8{0xd8, 04}},
	{AFSUBD, yfadd, Px, [23]uint8{0xdc, 04, 0xd8, 04, 0xdc, 05}},
	{AFSUBRDP, ycompp, Px, [23]uint8{0xde, 04}},
	{AFSUBRW, yfmvx, Px, [23]uint8{0xde, 05}},
	{AFSUBRL, yfmvx, Px, [23]uint8{0xda, 05}},
	{AFSUBRF, yfmvx, Px, [23]uint8{0xd8, 05}},
	{AFSUBRD, yfadd, Px, [23]uint8{0xdc, 05, 0xd8, 05, 0xdc, 04}},
	{AFDIVDP, ycompp, Px, [23]uint8{0xde, 07}},
	{AFDIVW, yfmvx, Px, [23]uint8{0xde, 06}},
	{AFDIVL, yfmvx, Px, [23]uint8{0xda, 06}},
	{AFDIVF, yfmvx, Px, [23]uint8{0xd8, 06}},
	{AFDIVD, yfadd, Px, [23]uint8{0xdc, 06, 0xd8, 06, 0xdc, 07}},
	{AFDIVRDP, ycompp, Px, [23]uint8{0xde, 06}},
	{AFDIVRW, yfmvx, Px, [23]uint8{0xde, 07}},
	{AFDIVRL, yfmvx, Px, [23]uint8{0xda, 07}},
	{AFDIVRF, yfmvx, Px, [23]uint8{0xd8, 07}},
	{AFDIVRD, yfadd, Px, [23]uint8{0xdc, 07, 0xd8, 07, 0xdc, 06}},
	{AFXCHD, yfxch, Px, [23]uint8{0xd9, 01, 0xd9, 01}},
	{AFFREE, nil, 0, [23]uint8{}},
	{AFLDCW, ysvrs_mo, Px, [23]uint8{0xd9, 05, 0xd9, 05}},
	{AFLDENV, ysvrs_mo, Px, [23]uint8{0xd9, 04, 0xd9, 04}},
	{AFRSTOR, ysvrs_mo, Px, [23]uint8{0xdd, 04, 0xdd, 04}},
	{AFSAVE, ysvrs_om, Px, [23]uint8{0xdd, 06, 0xdd, 06}},
	{AFSTCW, ysvrs_om, Px, [23]uint8{0xd9, 07, 0xd9, 07}},
	{AFSTENV, ysvrs_om, Px, [23]uint8{0xd9, 06, 0xd9, 06}},
	{AFSTSW, ystsw, Px, [23]uint8{0xdd, 07, 0xdf, 0xe0}},
	{AF2XM1, ynone, Px, [23]uint8{0xd9, 0xf0}},
	{AFABS, ynone, Px, [23]uint8{0xd9, 0xe1}},
	{AFCHS, ynone, Px, [23]uint8{0xd9, 0xe0}},
	{AFCLEX, ynone, Px, [23]uint8{0xdb, 0xe2}},
	{AFCOS, ynone, Px, [23]uint8{0xd9, 0xff}},
	{AFDECSTP, ynone, Px, [23]uint8{0xd9, 0xf6}},
	{AFINCSTP, ynone, Px, [23]uint8{0xd9, 0xf7}},
	{AFINIT, ynone, Px, [23]uint8{0xdb, 0xe3}},
	{AFLD1, ynone, Px, [23]uint8{0xd9, 0xe8}},
	{AFLDL2E, ynone, Px, [23]uint8{0xd9, 0xea}},
	{AFLDL2T, ynone, Px, [23]uint8{0xd9, 0xe9}},
	{AFLDLG2, ynone, Px, [23]uint8{0xd9, 0xec}},
	{AFLDLN2, ynone, Px, [23]uint8{0xd9, 0xed}},
	{AFLDPI, ynone, Px, [23]uint8{0xd9, 0xeb}},
	{AFLDZ, ynone, Px, [23]uint8{0xd9, 0xee}},
	{AFNOP, ynone, Px, [23]uint8{0xd9, 0xd0}},
	{AFPATAN, ynone, Px, [23]uint8{0xd9, 0xf3}},
	{AFPREM, ynone, Px, [23]uint8{0xd9, 0xf8}},
	{AFPREM1, ynone, Px, [23]uint8{0xd9, 0xf5}},
	{AFPTAN, ynone, Px, [23]uint8{0xd9, 0xf2}},
	{AFRNDINT, ynone, Px, [23]uint8{0xd9, 0xfc}},
	{AFSCALE, ynone, Px, [23]uint8{0xd9, 0xfd}},
	{AFSIN, ynone, Px, [23]uint8{0xd9, 0xfe}},
	{AFSINCOS, ynone, Px, [23]uint8{0xd9, 0xfb}},
	{AFSQRT, ynone, Px, [23]uint8{0xd9, 0xfa}},
	{AFTST, ynone, Px, [23]uint8{0xd9, 0xe4}},
	{AFXAM, ynone, Px, [23]uint8{0xd9, 0xe5}},
	{AFXTRACT, ynone, Px, [23]uint8{0xd9, 0xf4}},
	{AFYL2X, ynone, Px, [23]uint8{0xd9, 0xf1}},
	{AFYL2XP1, ynone, Px, [23]uint8{0xd9, 0xf9}},
	{ACMPXCHGB, yrb_mb, Pb, [23]uint8{0x0f, 0xb0}},
	{ACMPXCHGL, yrl_ml, Px, [23]uint8{0x0f, 0xb1}},
	{ACMPXCHGW, yrl_ml, Pe, [23]uint8{0x0f, 0xb1}},
	{ACMPXCHGQ, yrl_ml, Pw, [23]uint8{0x0f, 0xb1}},
	{ACMPXCHG8B, yscond, Pm, [23]uint8{0xc7, 01}},
	{AINVD, ynone, Pm, [23]uint8{0x08}},
	{AINVLPG, ydivb, Pm, [23]uint8{0x01, 07}},
	{ALFENCE, ynone, Pm, [23]uint8{0xae, 0xe8}},
	{AMFENCE, ynone, Pm, [23]uint8{0xae, 0xf0}},
	{AMOVNTIL, yrl_ml, Pm, [23]uint8{0xc3}},
	{AMOVNTIQ, yrl_ml, Pw, [23]uint8{0x0f, 0xc3}},
	{ARDMSR, ynone, Pm, [23]uint8{0x32}},
	{ARDPMC, ynone, Pm, [23]uint8{0x33}},
	{ARDTSC, ynone, Pm, [23]uint8{0x31}},
	{ARSM, ynone, Pm, [23]uint8{0xaa}},
	{ASFENCE, ynone, Pm, [23]uint8{0xae, 0xf8}},
	{ASYSRET, ynone, Pm, [23]uint8{0x07}},
	{AWBINVD, ynone, Pm, [23]uint8{0x09}},
	{AWRMSR, ynone, Pm, [23]uint8{0x30}},
	{AXADDB, yrb_mb, Pb, [23]uint8{0x0f, 0xc0}},
	{AXADDL, yrl_ml, Px, [23]uint8{0x0f, 0xc1}},
	{AXADDQ, yrl_ml, Pw, [23]uint8{0x0f, 0xc1}},
	{AXADDW, yrl_ml, Pe, [23]uint8{0x0f, 0xc1}},
	{ACRC32B, ycrc32l, Px, [23]uint8{0xf2, 0x0f, 0x38, 0xf0, 0}},
	{ACRC32Q, ycrc32l, Pw, [23]uint8{0xf2, 0x0f, 0x38, 0xf1, 0}},
	{APREFETCHT0, yprefetch, Pm, [23]uint8{0x18, 01}},
	{APREFETCHT1, yprefetch, Pm, [23]uint8{0x18, 02}},
	{APREFETCHT2, yprefetch, Pm, [23]uint8{0x18, 03}},
	{APREFETCHNTA, yprefetch, Pm, [23]uint8{0x18, 00}},
	{AMOVQL, yrl_ml, Px, [23]uint8{0x89}},
	{obj.AUNDEF, ynone, Px, [23]uint8{0x0f, 0x0b}},
	{AAESENC, yaes, Pq, [23]uint8{0x38, 0xdc, 0}},
	{AAESENCLAST, yaes, Pq, [23]uint8{0x38, 0xdd, 0}},
	{AAESDEC, yaes, Pq, [23]uint8{0x38, 0xde, 0}},
	{AAESDECLAST, yaes, Pq, [23]uint8{0x38, 0xdf, 0}},
	{AAESIMC, yaes, Pq, [23]uint8{0x38, 0xdb, 0}},
	{AAESKEYGENASSIST, yxshuf, Pq, [23]uint8{0x3a, 0xdf, 0}},
	{AROUNDPD, yxshuf, Pq, [23]uint8{0x3a, 0x09, 0}},
	{AROUNDPS, yxshuf, Pq, [23]uint8{0x3a, 0x08, 0}},
	{AROUNDSD, yxshuf, Pq, [23]uint8{0x3a, 0x0b, 0}},
	{AROUNDSS, yxshuf, Pq, [23]uint8{0x3a, 0x0a, 0}},
	{APSHUFD, yxshuf, Pq, [23]uint8{0x70, 0}},
	{APCLMULQDQ, yxshuf, Pq, [23]uint8{0x3a, 0x44, 0}},
	{APCMPESTRI, yxshuf, Pq, [23]uint8{0x3a, 0x61, 0}},
	{APCMPESTRM, yxshuf, Pq, [23]uint8{0x3a, 0x60, 0}},
	{AMOVDDUP, yxm, Pf2, [23]uint8{0x12}},
	{AMOVSHDUP, yxm, Pf3, [23]uint8{0x16}},
	{AMOVSLDUP, yxm, Pf3, [23]uint8{0x12}},

	{ABLENDPD, yxshuf, Pq, [23]uint8{0x3a, 0x0d, 0}},
	{ABLENDPS, yxshuf, Pq, [23]uint8{0x3a, 0x0c, 0}},
	{AXACQUIRE, ynone, Px, [23]uint8{0xf2}},
	{AXRELEASE, ynone, Px, [23]uint8{0xf3}},
	{AXBEGIN, yxbegin, Px, [23]uint8{0xc7, 0xf8}},
	{AXABORT, yxabort, Px, [23]uint8{0xc6, 0xf8}},
	{AXEND, ynone, Px, [23]uint8{0x0f, 01, 0xd5}},
	{AXTEST, ynone, Px, [23]uint8{0x0f, 01, 0xd6}},
	{AXGETBV, ynone, Pm, [23]uint8{01, 0xd0}},
	{obj.AFUNCDATA, yfuncdata, Px, [23]uint8{0, 0}},
	{obj.APCDATA, ypcdata, Px, [23]uint8{0, 0}},
	{obj.ADUFFCOPY, yduff, Px, [23]uint8{0xe8}},
	{obj.ADUFFZERO, yduff, Px, [23]uint8{0xe8}},

	// AVX2 gather instructions.
	// Added as a part of VSIB support implementation,
	// when x86avxgen will output these, they will be moved to
	// vex_optabs.go where they belong.
	{AVGATHERDPD, yvpgatherdq, Pvex, [23]uint8{
		vexDDS | vex128 | vex66 | vex0F38 | vexW1, 0x92,
		vexDDS | vex256 | vex66 | vex0F38 | vexW1, 0x92,
	}},
	{AVGATHERQPD, yvpgatherqq, Pvex, [23]uint8{
		vexDDS | vex128 | vex66 | vex0F38 | vexW1, 0x93,
		vexDDS | vex256 | vex66 | vex0F38 | vexW1, 0x93,
	}},
	{AVGATHERDPS, yvpgatherqq, Pvex, [23]uint8{
		vexDDS | vex128 | vex66 | vex0F38 | vexW0, 0x92,
		vexDDS | vex256 | vex66 | vex0F38 | vexW0, 0x92,
	}},
	{AVGATHERQPS, yvgatherqps, Pvex, [23]uint8{
		vexDDS | vex128 | vex66 | vex0F38 | vexW0, 0x93,
		vexDDS | vex256 | vex66 | vex0F38 | vexW0, 0x93,
	}},
	{AVPGATHERDD, yvpgatherqq, Pvex, [23]uint8{
		vexDDS | vex128 | vex66 | vex0F38 | vexW0, 0x90,
		vexDDS | vex256 | vex66 | vex0F38 | vexW0, 0x90,
	}},
	{AVPGATHERQD, yvgatherqps, Pvex, [23]uint8{
		vexDDS | vex128 | vex66 | vex0F38 | vexW0, 0x91,
		vexDDS | vex256 | vex66 | vex0F38 | vexW0, 0x91,
	}},
	{AVPGATHERDQ, yvpgatherdq, Pvex, [23]uint8{
		vexDDS | vex128 | vex66 | vex0F38 | vexW1, 0x90,
		vexDDS | vex256 | vex66 | vex0F38 | vexW1, 0x90,
	}},
	{AVPGATHERQQ, yvpgatherqq, Pvex, [23]uint8{
		vexDDS | vex128 | vex66 | vex0F38 | vexW1, 0x91,
		vexDDS | vex256 | vex66 | vex0F38 | vexW1, 0x91,
	}},

	{obj.AEND, nil, 0, [23]uint8{}},
	{0, nil, 0, [23]uint8{}},
}

var opindex [(ALAST + 1) & obj.AMask]*Optab

// useAbs reports whether s describes a symbol that must avoid pc-relative addressing.
// This happens on systems like Solaris that call .so functions instead of system calls.
// It does not seem to be necessary for any other systems. This is probably working
// around a Solaris-specific bug that should be fixed differently, but we don't know
// what that bug is. And this does fix it.
func useAbs(ctxt *obj.Link, s *obj.LSym) bool {
	if ctxt.Headtype == objabi.Hsolaris {
		// All the Solaris dynamic imports from libc.so begin with "libc_".
		return strings.HasPrefix(s.Name, "libc_")
	}
	return ctxt.Arch.Family == sys.I386 && !ctxt.Flag_shared
}

// single-instruction no-ops of various lengths.
// constructed by hand and disassembled with gdb to verify.
// see http://www.agner.org/optimize/optimizing_assembly.pdf for discussion.
var nop = [][16]uint8{
	{0x90},
	{0x66, 0x90},
	{0x0F, 0x1F, 0x00},
	{0x0F, 0x1F, 0x40, 0x00},
	{0x0F, 0x1F, 0x44, 0x00, 0x00},
	{0x66, 0x0F, 0x1F, 0x44, 0x00, 0x00},
	{0x0F, 0x1F, 0x80, 0x00, 0x00, 0x00, 0x00},
	{0x0F, 0x1F, 0x84, 0x00, 0x00, 0x00, 0x00, 0x00},
	{0x66, 0x0F, 0x1F, 0x84, 0x00, 0x00, 0x00, 0x00, 0x00},
}

// Native Client rejects the repeated 0x66 prefix.
// {0x66, 0x66, 0x0F, 0x1F, 0x84, 0x00, 0x00, 0x00, 0x00, 0x00},
func fillnop(p []byte, n int) {
	var m int

	for n > 0 {
		m = n
		if m > len(nop) {
			m = len(nop)
		}
		copy(p[:m], nop[m-1][:m])
		p = p[m:]
		n -= m
	}
}

func naclpad(ctxt *obj.Link, s *obj.LSym, c int32, pad int32) int32 {
	s.Grow(int64(c) + int64(pad))
	fillnop(s.P[c:], int(pad))
	return c + pad
}

func spadjop(ctxt *obj.Link, p *obj.Prog, l, q obj.As) obj.As {
	if ctxt.Arch.Family != sys.AMD64 || ctxt.Arch.PtrSize == 4 {
		return l
	}
	return q
}

func span6(ctxt *obj.Link, s *obj.LSym, newprog obj.ProgAlloc) {
	if s.P != nil {
		return
	}

	if ycover[0] == 0 {
		ctxt.Diag("x86 tables not initialized, call x86.instinit first")
	}

	var asmbuf AsmBuf

	for p := s.Func.Text; p != nil; p = p.Link {
		if p.To.Type == obj.TYPE_BRANCH {
			if p.Pcond == nil {
				p.Pcond = p
			}
		}
		if p.As == AADJSP {
			p.To.Type = obj.TYPE_REG
			p.To.Reg = REG_SP
			v := int32(-p.From.Offset)
			p.From.Offset = int64(v)
			p.As = spadjop(ctxt, p, AADDL, AADDQ)
			if v < 0 {
				p.As = spadjop(ctxt, p, ASUBL, ASUBQ)
				v = -v
				p.From.Offset = int64(v)
			}

			if v == 0 {
				p.As = obj.ANOP
			}
		}
	}

	var q *obj.Prog
	var count int64 // rough count of number of instructions
	for p := s.Func.Text; p != nil; p = p.Link {
		count++
		p.Back = 2 // use short branches first time through
		q = p.Pcond
		if q != nil && (q.Back&2 != 0) {
			p.Back |= 1 // backward jump
			q.Back |= 4 // loop head
		}

		if p.As == AADJSP {
			p.To.Type = obj.TYPE_REG
			p.To.Reg = REG_SP
			v := int32(-p.From.Offset)
			p.From.Offset = int64(v)
			p.As = spadjop(ctxt, p, AADDL, AADDQ)
			if v < 0 {
				p.As = spadjop(ctxt, p, ASUBL, ASUBQ)
				v = -v
				p.From.Offset = int64(v)
			}

			if v == 0 {
				p.As = obj.ANOP
			}
		}
	}
	s.GrowCap(count * 5) // preallocate roughly 5 bytes per instruction

	n := 0
	var c int32
	errors := ctxt.Errors
	for {
		loop := int32(0)
		for i := range s.R {
			s.R[i] = obj.Reloc{}
		}
		s.R = s.R[:0]
		s.P = s.P[:0]
		c = 0
		for p := s.Func.Text; p != nil; p = p.Link {
			if ctxt.Headtype == objabi.Hnacl && p.Isize > 0 {

				// pad everything to avoid crossing 32-byte boundary
				if c>>5 != (c+int32(p.Isize)-1)>>5 {
					c = naclpad(ctxt, s, c, -c&31)
				}

				// pad call deferreturn to start at 32-byte boundary
				// so that subtracting 5 in jmpdefer will jump back
				// to that boundary and rerun the call.
				if p.As == obj.ACALL && p.To.Sym == deferreturn {
					c = naclpad(ctxt, s, c, -c&31)
				}

				// pad call to end at 32-byte boundary
				if p.As == obj.ACALL {
					c = naclpad(ctxt, s, c, -(c+int32(p.Isize))&31)
				}

				// the linker treats REP and STOSQ as different instructions
				// but in fact the REP is a prefix on the STOSQ.
				// make sure REP has room for 2 more bytes, so that
				// padding will not be inserted before the next instruction.
				if (p.As == AREP || p.As == AREPN) && c>>5 != (c+3-1)>>5 {
					c = naclpad(ctxt, s, c, -c&31)
				}

				// same for LOCK.
				// various instructions follow; the longest is 4 bytes.
				// give ourselves 8 bytes so as to avoid surprises.
				if p.As == ALOCK && c>>5 != (c+8-1)>>5 {
					c = naclpad(ctxt, s, c, -c&31)
				}
			}

			if (p.Back&4 != 0) && c&(LoopAlign-1) != 0 {
				// pad with NOPs
				v := -c & (LoopAlign - 1)

				if v <= MaxLoopPad {
					s.Grow(int64(c) + int64(v))
					fillnop(s.P[c:], int(v))
					c += v
				}
			}

			p.Pc = int64(c)

			// process forward jumps to p
			for q = p.Rel; q != nil; q = q.Forwd {
				v := int32(p.Pc - (q.Pc + int64(q.Isize)))
				if q.Back&2 != 0 { // short
					if v > 127 {
						loop++
						q.Back ^= 2
					}

					if q.As == AJCXZL || q.As == AXBEGIN {
						s.P[q.Pc+2] = byte(v)
					} else {
						s.P[q.Pc+1] = byte(v)
					}
				} else {
					binary.LittleEndian.PutUint32(s.P[q.Pc+int64(q.Isize)-4:], uint32(v))
				}
			}

			p.Rel = nil

			p.Pc = int64(c)
			asmbuf.asmins(ctxt, s, p)
			m := asmbuf.Len()
			if int(p.Isize) != m {
				p.Isize = uint8(m)
				loop++
			}

			s.Grow(p.Pc + int64(m))
			copy(s.P[p.Pc:], asmbuf.Bytes())
			c += int32(m)
		}

		n++
		if n > 20 {
			ctxt.Diag("span must be looping")
			log.Fatalf("loop")
		}
		if loop == 0 {
			break
		}
		if ctxt.Errors > errors {
			return
		}
	}

	if ctxt.Headtype == objabi.Hnacl {
		c = naclpad(ctxt, s, c, -c&31)
	}

	s.Size = int64(c)

	if false { /* debug['a'] > 1 */
		fmt.Printf("span1 %s %d (%d tries)\n %.6x", s.Name, s.Size, n, 0)
		var i int
		for i = 0; i < len(s.P); i++ {
			fmt.Printf(" %.2x", s.P[i])
			if i%16 == 15 {
				fmt.Printf("\n  %.6x", uint(i+1))
			}
		}

		if i%16 != 0 {
			fmt.Printf("\n")
		}

		for i := 0; i < len(s.R); i++ {
			r := &s.R[i]
			fmt.Printf(" rel %#.4x/%d %s%+d\n", uint32(r.Off), r.Siz, r.Sym.Name, r.Add)
		}
	}
}

func instinit(ctxt *obj.Link) {
	if ycover[0] != 0 {
		// Already initialized; stop now.
		// This happens in the cmd/asm tests,
		// each of which re-initializes the arch.
		return
	}

	switch ctxt.Headtype {
	case objabi.Hplan9:
		plan9privates = ctxt.Lookup("_privates")
	case objabi.Hnacl:
		deferreturn = ctxt.Lookup("runtime.deferreturn")
	}

	for i := range vexOptab {
		c := vexOptab[i].as
		if opindex[c&obj.AMask] != nil {
			ctxt.Diag("phase error in vexOptab: %d (%v)", i, c)
		}
		opindex[c&obj.AMask] = &vexOptab[i]
	}
	for i := 1; optab[i].as != 0; i++ {
		c := optab[i].as
		if opindex[c&obj.AMask] != nil {
			ctxt.Diag("phase error in optab: %d (%v)", i, c)
		}
		opindex[c&obj.AMask] = &optab[i]
	}

	for i := 0; i < Ymax; i++ {
		ycover[i*Ymax+i] = 1
	}

	ycover[Yi0*Ymax+Yu2] = 1
	ycover[Yi1*Ymax+Yu2] = 1

	ycover[Yi0*Ymax+Yi8] = 1
	ycover[Yi1*Ymax+Yi8] = 1
	ycover[Yu2*Ymax+Yi8] = 1
	ycover[Yu7*Ymax+Yi8] = 1

	ycover[Yi0*Ymax+Yu7] = 1
	ycover[Yi1*Ymax+Yu7] = 1
	ycover[Yu2*Ymax+Yu7] = 1

	ycover[Yi0*Ymax+Yu8] = 1
	ycover[Yi1*Ymax+Yu8] = 1
	ycover[Yu2*Ymax+Yu8] = 1
	ycover[Yu7*Ymax+Yu8] = 1

	ycover[Yi0*Ymax+Ys32] = 1
	ycover[Yi1*Ymax+Ys32] = 1
	ycover[Yu2*Ymax+Ys32] = 1
	ycover[Yu7*Ymax+Ys32] = 1
	ycover[Yu8*Ymax+Ys32] = 1
	ycover[Yi8*Ymax+Ys32] = 1

	ycover[Yi0*Ymax+Yi32] = 1
	ycover[Yi1*Ymax+Yi32] = 1
	ycover[Yu2*Ymax+Yi32] = 1
	ycover[Yu7*Ymax+Yi32] = 1
	ycover[Yu8*Ymax+Yi32] = 1
	ycover[Yi8*Ymax+Yi32] = 1
	ycover[Ys32*Ymax+Yi32] = 1

	ycover[Yi0*Ymax+Yi64] = 1
	ycover[Yi1*Ymax+Yi64] = 1
	ycover[Yu7*Ymax+Yi64] = 1
	ycover[Yu2*Ymax+Yi64] = 1
	ycover[Yu8*Ymax+Yi64] = 1
	ycover[Yi8*Ymax+Yi64] = 1
	ycover[Ys32*Ymax+Yi64] = 1
	ycover[Yi32*Ymax+Yi64] = 1

	ycover[Yal*Ymax+Yrb] = 1
	ycover[Ycl*Ymax+Yrb] = 1
	ycover[Yax*Ymax+Yrb] = 1
	ycover[Ycx*Ymax+Yrb] = 1
	ycover[Yrx*Ymax+Yrb] = 1
	ycover[Yrl*Ymax+Yrb] = 1 // but not Yrl32

	ycover[Ycl*Ymax+Ycx] = 1

	ycover[Yax*Ymax+Yrx] = 1
	ycover[Ycx*Ymax+Yrx] = 1

	ycover[Yax*Ymax+Yrl] = 1
	ycover[Ycx*Ymax+Yrl] = 1
	ycover[Yrx*Ymax+Yrl] = 1
	ycover[Yrl32*Ymax+Yrl] = 1

	ycover[Yf0*Ymax+Yrf] = 1

	ycover[Yal*Ymax+Ymb] = 1
	ycover[Ycl*Ymax+Ymb] = 1
	ycover[Yax*Ymax+Ymb] = 1
	ycover[Ycx*Ymax+Ymb] = 1
	ycover[Yrx*Ymax+Ymb] = 1
	ycover[Yrb*Ymax+Ymb] = 1
	ycover[Yrl*Ymax+Ymb] = 1 // but not Yrl32
	ycover[Ym*Ymax+Ymb] = 1

	ycover[Yax*Ymax+Yml] = 1
	ycover[Ycx*Ymax+Yml] = 1
	ycover[Yrx*Ymax+Yml] = 1
	ycover[Yrl*Ymax+Yml] = 1
	ycover[Yrl32*Ymax+Yml] = 1
	ycover[Ym*Ymax+Yml] = 1

	ycover[Yax*Ymax+Ymm] = 1
	ycover[Ycx*Ymax+Ymm] = 1
	ycover[Yrx*Ymax+Ymm] = 1
	ycover[Yrl*Ymax+Ymm] = 1
	ycover[Yrl32*Ymax+Ymm] = 1
	ycover[Ym*Ymax+Ymm] = 1
	ycover[Ymr*Ymax+Ymm] = 1

	ycover[Ym*Ymax+Yxm] = 1
	ycover[Yxr*Ymax+Yxm] = 1

	ycover[Ym*Ymax+Yym] = 1
	ycover[Yyr*Ymax+Yym] = 1

	for i := 0; i < MAXREG; i++ {
		reg[i] = -1
		if i >= REG_AL && i <= REG_R15B {
			reg[i] = (i - REG_AL) & 7
			if i >= REG_SPB && i <= REG_DIB {
				regrex[i] = 0x40
			}
			if i >= REG_R8B && i <= REG_R15B {
				regrex[i] = Rxr | Rxx | Rxb
			}
		}

		if i >= REG_AH && i <= REG_BH {
			reg[i] = 4 + ((i - REG_AH) & 7)
		}
		if i >= REG_AX && i <= REG_R15 {
			reg[i] = (i - REG_AX) & 7
			if i >= REG_R8 {
				regrex[i] = Rxr | Rxx | Rxb
			}
		}

		if i >= REG_F0 && i <= REG_F0+7 {
			reg[i] = (i - REG_F0) & 7
		}
		if i >= REG_M0 && i <= REG_M0+7 {
			reg[i] = (i - REG_M0) & 7
		}
		if i >= REG_X0 && i <= REG_X0+15 {
			reg[i] = (i - REG_X0) & 7
			if i >= REG_X0+8 {
				regrex[i] = Rxr | Rxx | Rxb
			}
		}
		if i >= REG_Y0 && i <= REG_Y0+15 {
			reg[i] = (i - REG_Y0) & 7
			if i >= REG_Y0+8 {
				regrex[i] = Rxr | Rxx | Rxb
			}
		}

		if i >= REG_CR+8 && i <= REG_CR+15 {
			regrex[i] = Rxr
		}
	}
}

var isAndroid = (objabi.GOOS == "android")

func prefixof(ctxt *obj.Link, p *obj.Prog, a *obj.Addr) int {
	if a.Reg < REG_CS && a.Index < REG_CS { // fast path
		return 0
	}
	if a.Type == obj.TYPE_MEM && a.Name == obj.NAME_NONE {
		switch a.Reg {
		case REG_CS:
			return 0x2e

		case REG_DS:
			return 0x3e

		case REG_ES:
			return 0x26

		case REG_FS:
			return 0x64

		case REG_GS:
			return 0x65

		case REG_TLS:
			// NOTE: Systems listed here should be only systems that
			// support direct TLS references like 8(TLS) implemented as
			// direct references from FS or GS. Systems that require
			// the initial-exec model, where you load the TLS base into
			// a register and then index from that register, do not reach
			// this code and should not be listed.
			if ctxt.Arch.Family == sys.I386 {
				switch ctxt.Headtype {
				default:
					if isAndroid {
						return 0x65 // GS
					}
					log.Fatalf("unknown TLS base register for %v", ctxt.Headtype)

				case objabi.Hdarwin,
					objabi.Hdragonfly,
					objabi.Hfreebsd,
					objabi.Hnetbsd,
					objabi.Hopenbsd:
					return 0x65 // GS
				}
			}

			switch ctxt.Headtype {
			default:
				log.Fatalf("unknown TLS base register for %v", ctxt.Headtype)

			case objabi.Hlinux:
				if isAndroid {
					return 0x64 // FS
				}

				if ctxt.Flag_shared {
					log.Fatalf("unknown TLS base register for linux with -shared")
				} else {
					return 0x64 // FS
				}

			case objabi.Hdragonfly,
				objabi.Hfreebsd,
				objabi.Hnetbsd,
				objabi.Hopenbsd,
				objabi.Hsolaris:
				return 0x64 // FS

			case objabi.Hdarwin:
				return 0x65 // GS
			}
		}
	}

	if ctxt.Arch.Family == sys.I386 {
		if a.Index == REG_TLS && ctxt.Flag_shared {
			// When building for inclusion into a shared library, an instruction of the form
			//     MOVL 0(CX)(TLS*1), AX
			// becomes
			//     mov %gs:(%ecx), %eax
			// which assumes that the correct TLS offset has been loaded into %ecx (today
			// there is only one TLS variable -- g -- so this is OK). When not building for
			// a shared library the instruction it becomes
			//     mov 0x0(%ecx), $eax
			// and a R_TLS_LE relocation, and so does not require a prefix.
			if a.Offset != 0 {
				ctxt.Diag("cannot handle non-0 offsets to TLS")
			}
			return 0x65 // GS
		}
		return 0
	}

	switch a.Index {
	case REG_CS:
		return 0x2e

	case REG_DS:
		return 0x3e

	case REG_ES:
		return 0x26

	case REG_TLS:
		if ctxt.Flag_shared && ctxt.Headtype != objabi.Hwindows {
			// When building for inclusion into a shared library, an instruction of the form
			//     MOV 0(CX)(TLS*1), AX
			// becomes
			//     mov %fs:(%rcx), %rax
			// which assumes that the correct TLS offset has been loaded into %rcx (today
			// there is only one TLS variable -- g -- so this is OK). When not building for
			// a shared library the instruction does not require a prefix.
			if a.Offset != 0 {
				log.Fatalf("cannot handle non-0 offsets to TLS")
			}
			return 0x64
		}

	case REG_FS:
		return 0x64

	case REG_GS:
		return 0x65
	}

	return 0
}

func oclass(ctxt *obj.Link, p *obj.Prog, a *obj.Addr) int {
	switch a.Type {
	case obj.TYPE_NONE:
		return Ynone

	case obj.TYPE_BRANCH:
		return Ybr

	case obj.TYPE_INDIR:
		if a.Name != obj.NAME_NONE && a.Reg == REG_NONE && a.Index == REG_NONE && a.Scale == 0 {
			return Yindir
		}
		return Yxxx

	case obj.TYPE_MEM:
		if a.Index == REG_SP {
			// Can't use SP as the index register
			return Yxxx
		}
		if a.Index >= REG_X0 && a.Index <= REG_X15 {
			if ctxt.Arch.Family == sys.I386 && a.Index > REG_X7 {
				return Yxxx
			}
			return Yxvm
		}
		if a.Index >= REG_Y0 && a.Index <= REG_Y15 {
			if ctxt.Arch.Family == sys.I386 && a.Index > REG_Y7 {
				return Yxxx
			}
			return Yyvm
		}
		if ctxt.Arch.Family == sys.AMD64 {
			// Offset must fit in a 32-bit signed field (or fit in a 32-bit unsigned field
			// where the sign extension doesn't matter).
			// Note: The latter happens only in assembly, for example crypto/sha1/sha1block_amd64.s.
			if !(a.Offset == int64(int32(a.Offset)) ||
				a.Offset == int64(uint32(a.Offset)) && p.As == ALEAL) {
				return Yxxx
			}
			switch a.Name {
			case obj.NAME_EXTERN, obj.NAME_STATIC, obj.NAME_GOTREF:
				// Global variables can't use index registers and their
				// base register is %rip (%rip is encoded as REG_NONE).
				if a.Reg != REG_NONE || a.Index != REG_NONE || a.Scale != 0 {
					return Yxxx
				}
			case obj.NAME_AUTO, obj.NAME_PARAM:
				// These names must have a base of SP.  The old compiler
				// uses 0 for the base register. SSA uses REG_SP.
				if a.Reg != REG_SP && a.Reg != 0 {
					return Yxxx
				}
			case obj.NAME_NONE:
				// everything is ok
			default:
				// unknown name
				return Yxxx
			}
		}
		return Ym

	case obj.TYPE_ADDR:
		switch a.Name {
		case obj.NAME_GOTREF:
			ctxt.Diag("unexpected TYPE_ADDR with NAME_GOTREF")
			return Yxxx

		case obj.NAME_EXTERN,
			obj.NAME_STATIC:
			if a.Sym != nil && useAbs(ctxt, a.Sym) {
				return Yi32
			}
			return Yiauto // use pc-relative addressing

		case obj.NAME_AUTO,
			obj.NAME_PARAM:
			return Yiauto
		}

		// TODO(rsc): DUFFZERO/DUFFCOPY encoding forgot to set a->index
		// and got Yi32 in an earlier version of this code.
		// Keep doing that until we fix yduff etc.
		if a.Sym != nil && strings.HasPrefix(a.Sym.Name, "runtime.duff") {
			return Yi32
		}

		if a.Sym != nil || a.Name != obj.NAME_NONE {
			ctxt.Diag("unexpected addr: %v", obj.Dconv(p, a))
		}
		fallthrough

		// fall through

	case obj.TYPE_CONST:
		if a.Sym != nil {
			ctxt.Diag("TYPE_CONST with symbol: %v", obj.Dconv(p, a))
		}

		v := a.Offset
		if ctxt.Arch.Family == sys.I386 {
			v = int64(int32(v))
		}
		if v == 0 {
			return Yi0
		}
		if v == 1 {
			return Yi1
		}
		if v >= 0 && v <= 3 {
			return Yu2
		}
		if v >= 0 && v <= 127 {
			return Yu7
		}
		if v >= 0 && v <= 255 {
			return Yu8
		}
		if v >= -128 && v <= 127 {
			return Yi8
		}
		if ctxt.Arch.Family == sys.I386 {
			return Yi32
		}
		l := int32(v)
		if int64(l) == v {
			return Ys32 /* can sign extend */
		}
		if v>>32 == 0 {
			return Yi32 /* unsigned */
		}
		return Yi64

	case obj.TYPE_TEXTSIZE:
		return Ytextsize
	}

	if a.Type != obj.TYPE_REG {
		ctxt.Diag("unexpected addr1: type=%d %v", a.Type, obj.Dconv(p, a))
		return Yxxx
	}

	switch a.Reg {
	case REG_AL:
		return Yal

	case REG_AX:
		return Yax

		/*
			case REG_SPB:
		*/
	case REG_BPB,
		REG_SIB,
		REG_DIB,
		REG_R8B,
		REG_R9B,
		REG_R10B,
		REG_R11B,
		REG_R12B,
		REG_R13B,
		REG_R14B,
		REG_R15B:
		if ctxt.Arch.Family == sys.I386 {
			return Yxxx
		}
		fallthrough

	case REG_DL,
		REG_BL,
		REG_AH,
		REG_CH,
		REG_DH,
		REG_BH:
		return Yrb

	case REG_CL:
		return Ycl

	case REG_CX:
		return Ycx

	case REG_DX, REG_BX:
		return Yrx

	case REG_R8, /* not really Yrl */
		REG_R9,
		REG_R10,
		REG_R11,
		REG_R12,
		REG_R13,
		REG_R14,
		REG_R15:
		if ctxt.Arch.Family == sys.I386 {
			return Yxxx
		}
		fallthrough

	case REG_SP, REG_BP, REG_SI, REG_DI:
		if ctxt.Arch.Family == sys.I386 {
			return Yrl32
		}
		return Yrl

	case REG_F0 + 0:
		return Yf0

	case REG_F0 + 1,
		REG_F0 + 2,
		REG_F0 + 3,
		REG_F0 + 4,
		REG_F0 + 5,
		REG_F0 + 6,
		REG_F0 + 7:
		return Yrf

	case REG_M0 + 0,
		REG_M0 + 1,
		REG_M0 + 2,
		REG_M0 + 3,
		REG_M0 + 4,
		REG_M0 + 5,
		REG_M0 + 6,
		REG_M0 + 7:
		return Ymr

	case REG_X0 + 0,
		REG_X0 + 1,
		REG_X0 + 2,
		REG_X0 + 3,
		REG_X0 + 4,
		REG_X0 + 5,
		REG_X0 + 6,
		REG_X0 + 7,
		REG_X0 + 8,
		REG_X0 + 9,
		REG_X0 + 10,
		REG_X0 + 11,
		REG_X0 + 12,
		REG_X0 + 13,
		REG_X0 + 14,
		REG_X0 + 15:
		return Yxr

	case REG_Y0 + 0,
		REG_Y0 + 1,
		REG_Y0 + 2,
		REG_Y0 + 3,
		REG_Y0 + 4,
		REG_Y0 + 5,
		REG_Y0 + 6,
		REG_Y0 + 7,
		REG_Y0 + 8,
		REG_Y0 + 9,
		REG_Y0 + 10,
		REG_Y0 + 11,
		REG_Y0 + 12,
		REG_Y0 + 13,
		REG_Y0 + 14,
		REG_Y0 + 15:
		return Yyr

	case REG_CS:
		return Ycs
	case REG_SS:
		return Yss
	case REG_DS:
		return Yds
	case REG_ES:
		return Yes
	case REG_FS:
		return Yfs
	case REG_GS:
		return Ygs
	case REG_TLS:
		return Ytls

	case REG_GDTR:
		return Ygdtr
	case REG_IDTR:
		return Yidtr
	case REG_LDTR:
		return Yldtr
	case REG_MSW:
		return Ymsw
	case REG_TASK:
		return Ytask

	case REG_CR + 0:
		return Ycr0
	case REG_CR + 1:
		return Ycr1
	case REG_CR + 2:
		return Ycr2
	case REG_CR + 3:
		return Ycr3
	case REG_CR + 4:
		return Ycr4
	case REG_CR + 5:
		return Ycr5
	case REG_CR + 6:
		return Ycr6
	case REG_CR + 7:
		return Ycr7
	case REG_CR + 8:
		return Ycr8

	case REG_DR + 0:
		return Ydr0
	case REG_DR + 1:
		return Ydr1
	case REG_DR + 2:
		return Ydr2
	case REG_DR + 3:
		return Ydr3
	case REG_DR + 4:
		return Ydr4
	case REG_DR + 5:
		return Ydr5
	case REG_DR + 6:
		return Ydr6
	case REG_DR + 7:
		return Ydr7

	case REG_TR + 0:
		return Ytr0
	case REG_TR + 1:
		return Ytr1
	case REG_TR + 2:
		return Ytr2
	case REG_TR + 3:
		return Ytr3
	case REG_TR + 4:
		return Ytr4
	case REG_TR + 5:
		return Ytr5
	case REG_TR + 6:
		return Ytr6
	case REG_TR + 7:
		return Ytr7
	}

	return Yxxx
}

// AsmBuf is a simple buffer to assemble variable-length x86 instructions into
// and hold assembly state.
type AsmBuf struct {
	buf     [100]byte
	off     int
	rexflag int
	vexflag int
	rep     int
	repn    int
	lock    bool
}

// Put1 appends one byte to the end of the buffer.
func (a *AsmBuf) Put1(x byte) {
	a.buf[a.off] = x
	a.off++
}

// Put2 appends two bytes to the end of the buffer.
func (a *AsmBuf) Put2(x, y byte) {
	a.buf[a.off+0] = x
	a.buf[a.off+1] = y
	a.off += 2
}

// Put3 appends three bytes to the end of the buffer.
func (a *AsmBuf) Put3(x, y, z byte) {
	a.buf[a.off+0] = x
	a.buf[a.off+1] = y
	a.buf[a.off+2] = z
	a.off += 3
}

// Put4 appends four bytes to the end of the buffer.
func (a *AsmBuf) Put4(x, y, z, w byte) {
	a.buf[a.off+0] = x
	a.buf[a.off+1] = y
	a.buf[a.off+2] = z
	a.buf[a.off+3] = w
	a.off += 4
}

// PutInt16 writes v into the buffer using little-endian encoding.
func (a *AsmBuf) PutInt16(v int16) {
	a.buf[a.off+0] = byte(v)
	a.buf[a.off+1] = byte(v >> 8)
	a.off += 2
}

// PutInt32 writes v into the buffer using little-endian encoding.
func (a *AsmBuf) PutInt32(v int32) {
	a.buf[a.off+0] = byte(v)
	a.buf[a.off+1] = byte(v >> 8)
	a.buf[a.off+2] = byte(v >> 16)
	a.buf[a.off+3] = byte(v >> 24)
	a.off += 4
}

// PutInt64 writes v into the buffer using little-endian encoding.
func (a *AsmBuf) PutInt64(v int64) {
	a.buf[a.off+0] = byte(v)
	a.buf[a.off+1] = byte(v >> 8)
	a.buf[a.off+2] = byte(v >> 16)
	a.buf[a.off+3] = byte(v >> 24)
	a.buf[a.off+4] = byte(v >> 32)
	a.buf[a.off+5] = byte(v >> 40)
	a.buf[a.off+6] = byte(v >> 48)
	a.buf[a.off+7] = byte(v >> 56)
	a.off += 8
}

// Put copies b into the buffer.
func (a *AsmBuf) Put(b []byte) {
	copy(a.buf[a.off:], b)
	a.off += len(b)
}

// Insert inserts b at offset i.
func (a *AsmBuf) Insert(i int, b byte) {
	a.off++
	copy(a.buf[i+1:a.off], a.buf[i:a.off-1])
	a.buf[i] = b
}

// Last returns the byte at the end of the buffer.
func (a *AsmBuf) Last() byte { return a.buf[a.off-1] }

// Len returns the length of the buffer.
func (a *AsmBuf) Len() int { return a.off }

// Bytes returns the contents of the buffer.
func (a *AsmBuf) Bytes() []byte { return a.buf[:a.off] }

// Reset empties the buffer.
func (a *AsmBuf) Reset() { a.off = 0 }

// At returns the byte at offset i.
func (a *AsmBuf) At(i int) byte { return a.buf[i] }

// asmidx emits SIB byte.
func (asmbuf *AsmBuf) asmidx(ctxt *obj.Link, scale int, index int, base int) {
	var i int

	// X/Y index register is used in VSIB.
	switch index {
	default:
		goto bad

	case REG_NONE:
		i = 4 << 3
		goto bas

	case REG_R8,
		REG_R9,
		REG_R10,
		REG_R11,
		REG_R12,
		REG_R13,
		REG_R14,
		REG_R15,
		REG_X8,
		REG_X9,
		REG_X10,
		REG_X11,
		REG_X12,
		REG_X13,
		REG_X14,
		REG_X15,
		REG_Y8,
		REG_Y9,
		REG_Y10,
		REG_Y11,
		REG_Y12,
		REG_Y13,
		REG_Y14,
		REG_Y15:
		if ctxt.Arch.Family == sys.I386 {
			goto bad
		}
		fallthrough

	case REG_AX,
		REG_CX,
		REG_DX,
		REG_BX,
		REG_BP,
		REG_SI,
		REG_DI,
		REG_X0,
		REG_X1,
		REG_X2,
		REG_X3,
		REG_X4,
		REG_X5,
		REG_X6,
		REG_X7,
		REG_Y0,
		REG_Y1,
		REG_Y2,
		REG_Y3,
		REG_Y4,
		REG_Y5,
		REG_Y6,
		REG_Y7:
		i = reg[index] << 3
	}

	switch scale {
	default:
		goto bad

	case 1:
		break

	case 2:
		i |= 1 << 6

	case 4:
		i |= 2 << 6

	case 8:
		i |= 3 << 6
	}

bas:
	switch base {
	default:
		goto bad

	case REG_NONE: /* must be mod=00 */
		i |= 5

	case REG_R8,
		REG_R9,
		REG_R10,
		REG_R11,
		REG_R12,
		REG_R13,
		REG_R14,
		REG_R15:
		if ctxt.Arch.Family == sys.I386 {
			goto bad
		}
		fallthrough

	case REG_AX,
		REG_CX,
		REG_DX,
		REG_BX,
		REG_SP,
		REG_BP,
		REG_SI,
		REG_DI:
		i |= reg[base]
	}

	asmbuf.Put1(byte(i))
	return

bad:
	ctxt.Diag("asmidx: bad address %d/%d/%d", scale, index, base)
	asmbuf.Put1(0)
	return
}

func (asmbuf *AsmBuf) relput4(ctxt *obj.Link, cursym *obj.LSym, p *obj.Prog, a *obj.Addr) {
	var rel obj.Reloc

	v := vaddr(ctxt, p, a, &rel)
	if rel.Siz != 0 {
		if rel.Siz != 4 {
			ctxt.Diag("bad reloc")
		}
		r := obj.Addrel(cursym)
		*r = rel
		r.Off = int32(p.Pc + int64(asmbuf.Len()))
	}

	asmbuf.PutInt32(int32(v))
}

func vaddr(ctxt *obj.Link, p *obj.Prog, a *obj.Addr, r *obj.Reloc) int64 {
	if r != nil {
		*r = obj.Reloc{}
	}

	switch a.Name {
	case obj.NAME_STATIC,
		obj.NAME_GOTREF,
		obj.NAME_EXTERN:
		s := a.Sym
		if r == nil {
			ctxt.Diag("need reloc for %v", obj.Dconv(p, a))
			log.Fatalf("reloc")
		}

		if a.Name == obj.NAME_GOTREF {
			r.Siz = 4
			r.Type = objabi.R_GOTPCREL
		} else if useAbs(ctxt, s) {
			r.Siz = 4
			r.Type = objabi.R_ADDR
		} else {
			r.Siz = 4
			r.Type = objabi.R_PCREL
		}

		r.Off = -1 // caller must fill in
		r.Sym = s
		r.Add = a.Offset

		return 0
	}

	if (a.Type == obj.TYPE_MEM || a.Type == obj.TYPE_ADDR) && a.Reg == REG_TLS {
		if r == nil {
			ctxt.Diag("need reloc for %v", obj.Dconv(p, a))
			log.Fatalf("reloc")
		}

		if !ctxt.Flag_shared || isAndroid || ctxt.Headtype == objabi.Hdarwin {
			r.Type = objabi.R_TLS_LE
			r.Siz = 4
			r.Off = -1 // caller must fill in
			r.Add = a.Offset
		}
		return 0
	}

	return a.Offset
}

func (asmbuf *AsmBuf) asmandsz(ctxt *obj.Link, cursym *obj.LSym, p *obj.Prog, a *obj.Addr, r int, rex int, m64 int) {
	var base int
	var rel obj.Reloc

	rex &= 0x40 | Rxr
	switch {
	case int64(int32(a.Offset)) == a.Offset:
		// Offset fits in sign-extended 32 bits.
	case int64(uint32(a.Offset)) == a.Offset && asmbuf.rexflag&Rxw == 0:
		// Offset fits in zero-extended 32 bits in a 32-bit instruction.
		// This is allowed for assembly that wants to use 32-bit hex
		// constants, e.g. LEAL 0x99999999(AX), AX.
	default:
		ctxt.Diag("offset too large in %s", p)
	}
	v := int32(a.Offset)
	rel.Siz = 0

	switch a.Type {
	case obj.TYPE_ADDR:
		if a.Name == obj.NAME_NONE {
			ctxt.Diag("unexpected TYPE_ADDR with NAME_NONE")
		}
		if a.Index == REG_TLS {
			ctxt.Diag("unexpected TYPE_ADDR with index==REG_TLS")
		}
		goto bad

	case obj.TYPE_REG:
		if a.Reg < REG_AL || REG_Y0+15 < a.Reg {
			goto bad
		}
		if v != 0 {
			goto bad
		}
		asmbuf.Put1(byte(3<<6 | reg[a.Reg]<<0 | r<<3))
		asmbuf.rexflag |= regrex[a.Reg]&(0x40|Rxb) | rex
		return
	}

	if a.Type != obj.TYPE_MEM {
		goto bad
	}

	if a.Index != REG_NONE && a.Index != REG_TLS {
		base := int(a.Reg)
		switch a.Name {
		case obj.NAME_EXTERN,
			obj.NAME_GOTREF,
			obj.NAME_STATIC:
			if !useAbs(ctxt, a.Sym) && ctxt.Arch.Family == sys.AMD64 {
				goto bad
			}
			if ctxt.Arch.Family == sys.I386 && ctxt.Flag_shared {
				// The base register has already been set. It holds the PC
				// of this instruction returned by a PC-reading thunk.
				// See obj6.go:rewriteToPcrel.
			} else {
				base = REG_NONE
			}
			v = int32(vaddr(ctxt, p, a, &rel))

		case obj.NAME_AUTO,
			obj.NAME_PARAM:
			base = REG_SP
		}

		asmbuf.rexflag |= regrex[int(a.Index)]&Rxx | regrex[base]&Rxb | rex
		if base == REG_NONE {
			asmbuf.Put1(byte(0<<6 | 4<<0 | r<<3))
			asmbuf.asmidx(ctxt, int(a.Scale), int(a.Index), base)
			goto putrelv
		}

		if v == 0 && rel.Siz == 0 && base != REG_BP && base != REG_R13 {
			asmbuf.Put1(byte(0<<6 | 4<<0 | r<<3))
			asmbuf.asmidx(ctxt, int(a.Scale), int(a.Index), base)
			return
		}

		if v >= -128 && v < 128 && rel.Siz == 0 {
			asmbuf.Put1(byte(1<<6 | 4<<0 | r<<3))
			asmbuf.asmidx(ctxt, int(a.Scale), int(a.Index), base)
			asmbuf.Put1(byte(v))
			return
		}

		asmbuf.Put1(byte(2<<6 | 4<<0 | r<<3))
		asmbuf.asmidx(ctxt, int(a.Scale), int(a.Index), base)
		goto putrelv
	}

	base = int(a.Reg)
	switch a.Name {
	case obj.NAME_STATIC,
		obj.NAME_GOTREF,
		obj.NAME_EXTERN:
		if a.Sym == nil {
			ctxt.Diag("bad addr: %v", p)
		}
		if ctxt.Arch.Family == sys.I386 && ctxt.Flag_shared {
			// The base register has already been set. It holds the PC
			// of this instruction returned by a PC-reading thunk.
			// See obj6.go:rewriteToPcrel.
		} else {
			base = REG_NONE
		}
		v = int32(vaddr(ctxt, p, a, &rel))

	case obj.NAME_AUTO,
		obj.NAME_PARAM:
		base = REG_SP
	}

	if base == REG_TLS {
		v = int32(vaddr(ctxt, p, a, &rel))
	}

	asmbuf.rexflag |= regrex[base]&Rxb | rex
	if base == REG_NONE || (REG_CS <= base && base <= REG_GS) || base == REG_TLS {
		if (a.Sym == nil || !useAbs(ctxt, a.Sym)) && base == REG_NONE && (a.Name == obj.NAME_STATIC || a.Name == obj.NAME_EXTERN || a.Name == obj.NAME_GOTREF) || ctxt.Arch.Family != sys.AMD64 {
			if a.Name == obj.NAME_GOTREF && (a.Offset != 0 || a.Index != 0 || a.Scale != 0) {
				ctxt.Diag("%v has offset against gotref", p)
			}
			asmbuf.Put1(byte(0<<6 | 5<<0 | r<<3))
			goto putrelv
		}

		// temporary
		asmbuf.Put2(
			byte(0<<6|4<<0|r<<3), // sib present
			0<<6|4<<3|5<<0,       // DS:d32
		)
		goto putrelv
	}

	if base == REG_SP || base == REG_R12 {
		if v == 0 {
			asmbuf.Put1(byte(0<<6 | reg[base]<<0 | r<<3))
			asmbuf.asmidx(ctxt, int(a.Scale), REG_NONE, base)
			return
		}

		if v >= -128 && v < 128 {
			asmbuf.Put1(byte(1<<6 | reg[base]<<0 | r<<3))
			asmbuf.asmidx(ctxt, int(a.Scale), REG_NONE, base)
			asmbuf.Put1(byte(v))
			return
		}

		asmbuf.Put1(byte(2<<6 | reg[base]<<0 | r<<3))
		asmbuf.asmidx(ctxt, int(a.Scale), REG_NONE, base)
		goto putrelv
	}

	if REG_AX <= base && base <= REG_R15 {
		if a.Index == REG_TLS && !ctxt.Flag_shared {
			rel = obj.Reloc{}
			rel.Type = objabi.R_TLS_LE
			rel.Siz = 4
			rel.Sym = nil
			rel.Add = int64(v)
			v = 0
		}

		if v == 0 && rel.Siz == 0 && base != REG_BP && base != REG_R13 {
			asmbuf.Put1(byte(0<<6 | reg[base]<<0 | r<<3))
			return
		}

		if v >= -128 && v < 128 && rel.Siz == 0 {
			asmbuf.Put2(byte(1<<6|reg[base]<<0|r<<3), byte(v))
			return
		}

		asmbuf.Put1(byte(2<<6 | reg[base]<<0 | r<<3))
		goto putrelv
	}

	goto bad

putrelv:
	if rel.Siz != 0 {
		if rel.Siz != 4 {
			ctxt.Diag("bad rel")
			goto bad
		}

		r := obj.Addrel(cursym)
		*r = rel
		r.Off = int32(p.Pc + int64(asmbuf.Len()))
	}

	asmbuf.PutInt32(v)
	return

bad:
	ctxt.Diag("asmand: bad address %v", obj.Dconv(p, a))
	return
}

func (asmbuf *AsmBuf) asmand(ctxt *obj.Link, cursym *obj.LSym, p *obj.Prog, a *obj.Addr, ra *obj.Addr) {
	asmbuf.asmandsz(ctxt, cursym, p, a, reg[ra.Reg], regrex[ra.Reg], 0)
}

func (asmbuf *AsmBuf) asmando(ctxt *obj.Link, cursym *obj.LSym, p *obj.Prog, a *obj.Addr, o int) {
	asmbuf.asmandsz(ctxt, cursym, p, a, o, 0, 0)
}

func bytereg(a *obj.Addr, t *uint8) {
	if a.Type == obj.TYPE_REG && a.Index == REG_NONE && (REG_AX <= a.Reg && a.Reg <= REG_R15) {
		a.Reg += REG_AL - REG_AX
		*t = 0
	}
}

func unbytereg(a *obj.Addr, t *uint8) {
	if a.Type == obj.TYPE_REG && a.Index == REG_NONE && (REG_AL <= a.Reg && a.Reg <= REG_R15B) {
		a.Reg += REG_AX - REG_AL
		*t = 0
	}
}

const (
	E = 0xff
)

var ymovtab = []Movtab{
	/* push */
	{APUSHL, Ycs, Ynone, Ynone, 0, [4]uint8{0x0e, E, 0, 0}},
	{APUSHL, Yss, Ynone, Ynone, 0, [4]uint8{0x16, E, 0, 0}},
	{APUSHL, Yds, Ynone, Ynone, 0, [4]uint8{0x1e, E, 0, 0}},
	{APUSHL, Yes, Ynone, Ynone, 0, [4]uint8{0x06, E, 0, 0}},
	{APUSHL, Yfs, Ynone, Ynone, 0, [4]uint8{0x0f, 0xa0, E, 0}},
	{APUSHL, Ygs, Ynone, Ynone, 0, [4]uint8{0x0f, 0xa8, E, 0}},
	{APUSHQ, Yfs, Ynone, Ynone, 0, [4]uint8{0x0f, 0xa0, E, 0}},
	{APUSHQ, Ygs, Ynone, Ynone, 0, [4]uint8{0x0f, 0xa8, E, 0}},
	{APUSHW, Ycs, Ynone, Ynone, 0, [4]uint8{Pe, 0x0e, E, 0}},
	{APUSHW, Yss, Ynone, Ynone, 0, [4]uint8{Pe, 0x16, E, 0}},
	{APUSHW, Yds, Ynone, Ynone, 0, [4]uint8{Pe, 0x1e, E, 0}},
	{APUSHW, Yes, Ynone, Ynone, 0, [4]uint8{Pe, 0x06, E, 0}},
	{APUSHW, Yfs, Ynone, Ynone, 0, [4]uint8{Pe, 0x0f, 0xa0, E}},
	{APUSHW, Ygs, Ynone, Ynone, 0, [4]uint8{Pe, 0x0f, 0xa8, E}},

	/* pop */
	{APOPL, Ynone, Ynone, Yds, 0, [4]uint8{0x1f, E, 0, 0}},
	{APOPL, Ynone, Ynone, Yes, 0, [4]uint8{0x07, E, 0, 0}},
	{APOPL, Ynone, Ynone, Yss, 0, [4]uint8{0x17, E, 0, 0}},
	{APOPL, Ynone, Ynone, Yfs, 0, [4]uint8{0x0f, 0xa1, E, 0}},
	{APOPL, Ynone, Ynone, Ygs, 0, [4]uint8{0x0f, 0xa9, E, 0}},
	{APOPQ, Ynone, Ynone, Yfs, 0, [4]uint8{0x0f, 0xa1, E, 0}},
	{APOPQ, Ynone, Ynone, Ygs, 0, [4]uint8{0x0f, 0xa9, E, 0}},
	{APOPW, Ynone, Ynone, Yds, 0, [4]uint8{Pe, 0x1f, E, 0}},
	{APOPW, Ynone, Ynone, Yes, 0, [4]uint8{Pe, 0x07, E, 0}},
	{APOPW, Ynone, Ynone, Yss, 0, [4]uint8{Pe, 0x17, E, 0}},
	{APOPW, Ynone, Ynone, Yfs, 0, [4]uint8{Pe, 0x0f, 0xa1, E}},
	{APOPW, Ynone, Ynone, Ygs, 0, [4]uint8{Pe, 0x0f, 0xa9, E}},

	/* mov seg */
	{AMOVW, Yes, Ynone, Yml, 1, [4]uint8{0x8c, 0, 0, 0}},
	{AMOVW, Ycs, Ynone, Yml, 1, [4]uint8{0x8c, 1, 0, 0}},
	{AMOVW, Yss, Ynone, Yml, 1, [4]uint8{0x8c, 2, 0, 0}},
	{AMOVW, Yds, Ynone, Yml, 1, [4]uint8{0x8c, 3, 0, 0}},
	{AMOVW, Yfs, Ynone, Yml, 1, [4]uint8{0x8c, 4, 0, 0}},
	{AMOVW, Ygs, Ynone, Yml, 1, [4]uint8{0x8c, 5, 0, 0}},
	{AMOVW, Yml, Ynone, Yes, 2, [4]uint8{0x8e, 0, 0, 0}},
	{AMOVW, Yml, Ynone, Ycs, 2, [4]uint8{0x8e, 1, 0, 0}},
	{AMOVW, Yml, Ynone, Yss, 2, [4]uint8{0x8e, 2, 0, 0}},
	{AMOVW, Yml, Ynone, Yds, 2, [4]uint8{0x8e, 3, 0, 0}},
	{AMOVW, Yml, Ynone, Yfs, 2, [4]uint8{0x8e, 4, 0, 0}},
	{AMOVW, Yml, Ynone, Ygs, 2, [4]uint8{0x8e, 5, 0, 0}},

	/* mov cr */
	{AMOVL, Ycr0, Ynone, Yml, 3, [4]uint8{0x0f, 0x20, 0, 0}},
	{AMOVL, Ycr2, Ynone, Yml, 3, [4]uint8{0x0f, 0x20, 2, 0}},
	{AMOVL, Ycr3, Ynone, Yml, 3, [4]uint8{0x0f, 0x20, 3, 0}},
	{AMOVL, Ycr4, Ynone, Yml, 3, [4]uint8{0x0f, 0x20, 4, 0}},
	{AMOVL, Ycr8, Ynone, Yml, 3, [4]uint8{0x0f, 0x20, 8, 0}},
	{AMOVQ, Ycr0, Ynone, Yml, 3, [4]uint8{0x0f, 0x20, 0, 0}},
	{AMOVQ, Ycr2, Ynone, Yml, 3, [4]uint8{0x0f, 0x20, 2, 0}},
	{AMOVQ, Ycr3, Ynone, Yml, 3, [4]uint8{0x0f, 0x20, 3, 0}},
	{AMOVQ, Ycr4, Ynone, Yml, 3, [4]uint8{0x0f, 0x20, 4, 0}},
	{AMOVQ, Ycr8, Ynone, Yml, 3, [4]uint8{0x0f, 0x20, 8, 0}},
	{AMOVL, Yml, Ynone, Ycr0, 4, [4]uint8{0x0f, 0x22, 0, 0}},
	{AMOVL, Yml, Ynone, Ycr2, 4, [4]uint8{0x0f, 0x22, 2, 0}},
	{AMOVL, Yml, Ynone, Ycr3, 4, [4]uint8{0x0f, 0x22, 3, 0}},
	{AMOVL, Yml, Ynone, Ycr4, 4, [4]uint8{0x0f, 0x22, 4, 0}},
	{AMOVL, Yml, Ynone, Ycr8, 4, [4]uint8{0x0f, 0x22, 8, 0}},
	{AMOVQ, Yml, Ynone, Ycr0, 4, [4]uint8{0x0f, 0x22, 0, 0}},
	{AMOVQ, Yml, Ynone, Ycr2, 4, [4]uint8{0x0f, 0x22, 2, 0}},
	{AMOVQ, Yml, Ynone, Ycr3, 4, [4]uint8{0x0f, 0x22, 3, 0}},
	{AMOVQ, Yml, Ynone, Ycr4, 4, [4]uint8{0x0f, 0x22, 4, 0}},
	{AMOVQ, Yml, Ynone, Ycr8, 4, [4]uint8{0x0f, 0x22, 8, 0}},

	/* mov dr */
	{AMOVL, Ydr0, Ynone, Yml, 3, [4]uint8{0x0f, 0x21, 0, 0}},
	{AMOVL, Ydr6, Ynone, Yml, 3, [4]uint8{0x0f, 0x21, 6, 0}},
	{AMOVL, Ydr7, Ynone, Yml, 3, [4]uint8{0x0f, 0x21, 7, 0}},
	{AMOVQ, Ydr0, Ynone, Yml, 3, [4]uint8{0x0f, 0x21, 0, 0}},
	{AMOVQ, Ydr6, Ynone, Yml, 3, [4]uint8{0x0f, 0x21, 6, 0}},
	{AMOVQ, Ydr7, Ynone, Yml, 3, [4]uint8{0x0f, 0x21, 7, 0}},
	{AMOVL, Yml, Ynone, Ydr0, 4, [4]uint8{0x0f, 0x23, 0, 0}},
	{AMOVL, Yml, Ynone, Ydr6, 4, [4]uint8{0x0f, 0x23, 6, 0}},
	{AMOVL, Yml, Ynone, Ydr7, 4, [4]uint8{0x0f, 0x23, 7, 0}},
	{AMOVQ, Yml, Ynone, Ydr0, 4, [4]uint8{0x0f, 0x23, 0, 0}},
	{AMOVQ, Yml, Ynone, Ydr6, 4, [4]uint8{0x0f, 0x23, 6, 0}},
	{AMOVQ, Yml, Ynone, Ydr7, 4, [4]uint8{0x0f, 0x23, 7, 0}},

	/* mov tr */
	{AMOVL, Ytr6, Ynone, Yml, 3, [4]uint8{0x0f, 0x24, 6, 0}},
	{AMOVL, Ytr7, Ynone, Yml, 3, [4]uint8{0x0f, 0x24, 7, 0}},
	{AMOVL, Yml, Ynone, Ytr6, 4, [4]uint8{0x0f, 0x26, 6, E}},
	{AMOVL, Yml, Ynone, Ytr7, 4, [4]uint8{0x0f, 0x26, 7, E}},

	/* lgdt, sgdt, lidt, sidt */
	{AMOVL, Ym, Ynone, Ygdtr, 4, [4]uint8{0x0f, 0x01, 2, 0}},
	{AMOVL, Ygdtr, Ynone, Ym, 3, [4]uint8{0x0f, 0x01, 0, 0}},
	{AMOVL, Ym, Ynone, Yidtr, 4, [4]uint8{0x0f, 0x01, 3, 0}},
	{AMOVL, Yidtr, Ynone, Ym, 3, [4]uint8{0x0f, 0x01, 1, 0}},
	{AMOVQ, Ym, Ynone, Ygdtr, 4, [4]uint8{0x0f, 0x01, 2, 0}},
	{AMOVQ, Ygdtr, Ynone, Ym, 3, [4]uint8{0x0f, 0x01, 0, 0}},
	{AMOVQ, Ym, Ynone, Yidtr, 4, [4]uint8{0x0f, 0x01, 3, 0}},
	{AMOVQ, Yidtr, Ynone, Ym, 3, [4]uint8{0x0f, 0x01, 1, 0}},

	/* lldt, sldt */
	{AMOVW, Yml, Ynone, Yldtr, 4, [4]uint8{0x0f, 0x00, 2, 0}},
	{AMOVW, Yldtr, Ynone, Yml, 3, [4]uint8{0x0f, 0x00, 0, 0}},

	/* lmsw, smsw */
	{AMOVW, Yml, Ynone, Ymsw, 4, [4]uint8{0x0f, 0x01, 6, 0}},
	{AMOVW, Ymsw, Ynone, Yml, 3, [4]uint8{0x0f, 0x01, 4, 0}},

	/* ltr, str */
	{AMOVW, Yml, Ynone, Ytask, 4, [4]uint8{0x0f, 0x00, 3, 0}},
	{AMOVW, Ytask, Ynone, Yml, 3, [4]uint8{0x0f, 0x00, 1, 0}},

	/* load full pointer - unsupported
	Movtab{AMOVL, Yml, Ycol, 5, [4]uint8{0, 0, 0, 0}},
	Movtab{AMOVW, Yml, Ycol, 5, [4]uint8{Pe, 0, 0, 0}},
	*/

	/* double shift */
	{ASHLL, Yi8, Yrl, Yml, 6, [4]uint8{0xa4, 0xa5, 0, 0}},
	{ASHLL, Ycl, Yrl, Yml, 6, [4]uint8{0xa4, 0xa5, 0, 0}},
	{ASHLL, Ycx, Yrl, Yml, 6, [4]uint8{0xa4, 0xa5, 0, 0}},
	{ASHRL, Yi8, Yrl, Yml, 6, [4]uint8{0xac, 0xad, 0, 0}},
	{ASHRL, Ycl, Yrl, Yml, 6, [4]uint8{0xac, 0xad, 0, 0}},
	{ASHRL, Ycx, Yrl, Yml, 6, [4]uint8{0xac, 0xad, 0, 0}},
	{ASHLQ, Yi8, Yrl, Yml, 6, [4]uint8{Pw, 0xa4, 0xa5, 0}},
	{ASHLQ, Ycl, Yrl, Yml, 6, [4]uint8{Pw, 0xa4, 0xa5, 0}},
	{ASHLQ, Ycx, Yrl, Yml, 6, [4]uint8{Pw, 0xa4, 0xa5, 0}},
	{ASHRQ, Yi8, Yrl, Yml, 6, [4]uint8{Pw, 0xac, 0xad, 0}},
	{ASHRQ, Ycl, Yrl, Yml, 6, [4]uint8{Pw, 0xac, 0xad, 0}},
	{ASHRQ, Ycx, Yrl, Yml, 6, [4]uint8{Pw, 0xac, 0xad, 0}},
	{ASHLW, Yi8, Yrl, Yml, 6, [4]uint8{Pe, 0xa4, 0xa5, 0}},
	{ASHLW, Ycl, Yrl, Yml, 6, [4]uint8{Pe, 0xa4, 0xa5, 0}},
	{ASHLW, Ycx, Yrl, Yml, 6, [4]uint8{Pe, 0xa4, 0xa5, 0}},
	{ASHRW, Yi8, Yrl, Yml, 6, [4]uint8{Pe, 0xac, 0xad, 0}},
	{ASHRW, Ycl, Yrl, Yml, 6, [4]uint8{Pe, 0xac, 0xad, 0}},
	{ASHRW, Ycx, Yrl, Yml, 6, [4]uint8{Pe, 0xac, 0xad, 0}},

	/* load TLS base */
	{AMOVL, Ytls, Ynone, Yrl, 7, [4]uint8{0, 0, 0, 0}},
	{AMOVQ, Ytls, Ynone, Yrl, 7, [4]uint8{0, 0, 0, 0}},
	{0, 0, 0, 0, 0, [4]uint8{}},
}

func isax(a *obj.Addr) bool {
	switch a.Reg {
	case REG_AX, REG_AL, REG_AH:
		return true
	}

	if a.Index == REG_AX {
		return true
	}
	return false
}

func subreg(p *obj.Prog, from int, to int) {
	if false { /* debug['Q'] */
		fmt.Printf("\n%v\ts/%v/%v/\n", p, rconv(from), rconv(to))
	}

	if int(p.From.Reg) == from {
		p.From.Reg = int16(to)
		p.Ft = 0
	}

	if int(p.To.Reg) == from {
		p.To.Reg = int16(to)
		p.Tt = 0
	}

	if int(p.From.Index) == from {
		p.From.Index = int16(to)
		p.Ft = 0
	}

	if int(p.To.Index) == from {
		p.To.Index = int16(to)
		p.Tt = 0
	}

	if false { /* debug['Q'] */
		fmt.Printf("%v\n", p)
	}
}

func (asmbuf *AsmBuf) mediaop(ctxt *obj.Link, o *Optab, op int, osize int, z int) int {
	switch op {
	case Pm, Pe, Pf2, Pf3:
		if osize != 1 {
			if op != Pm {
				asmbuf.Put1(byte(op))
			}
			asmbuf.Put1(Pm)
			z++
			op = int(o.op[z])
			break
		}
		fallthrough

	default:
		if asmbuf.Len() == 0 || asmbuf.Last() != Pm {
			asmbuf.Put1(Pm)
		}
	}

	asmbuf.Put1(byte(op))
	return z
}

var bpduff1 = []byte{
	0x48, 0x89, 0x6c, 0x24, 0xf0, // MOVQ BP, -16(SP)
	0x48, 0x8d, 0x6c, 0x24, 0xf0, // LEAQ -16(SP), BP
}

var bpduff2 = []byte{
	0x48, 0x8b, 0x6d, 0x00, // MOVQ 0(BP), BP
}

// Emit VEX prefix and opcode byte.
// The three addresses are the r/m, vvvv, and reg fields.
// The reg and rm arguments appear in the same order as the
// arguments to asmand, which typically follows the call to asmvex.
// The final two arguments are the VEX prefix (see encoding above)
// and the opcode byte.
// For details about vex prefix see:
// https://en.wikipedia.org/wiki/VEX_prefix#Technical_description
func (asmbuf *AsmBuf) asmvex(ctxt *obj.Link, rm, v, r *obj.Addr, vex, opcode uint8) {
	asmbuf.vexflag = 1
	rexR := 0
	if r != nil {
		rexR = regrex[r.Reg] & Rxr
	}
	rexB := 0
	rexX := 0
	if rm != nil {
		rexB = regrex[rm.Reg] & Rxb
		rexX = regrex[rm.Index] & Rxx
	}
	vexM := (vex >> 3) & 0xF
	vexWLP := vex & 0x87
	vexV := byte(0)
	if v != nil {
		vexV = byte(reg[v.Reg]|(regrex[v.Reg]&Rxr)<<1) & 0xF
	}
	vexV ^= 0xF
	if vexM == 1 && (rexX|rexB) == 0 && vex&vexW1 == 0 {
		// Can use 2-byte encoding.
		asmbuf.Put2(0xc5, byte(rexR<<5)^0x80|vexV<<3|vexWLP)
	} else {
		// Must use 3-byte encoding.
		asmbuf.Put3(0xc4,
			(byte(rexR|rexX|rexB)<<5)^0xE0|vexM,
			vexV<<3|vexWLP,
		)
	}
	asmbuf.Put1(opcode)
}

// regIndex returns register index that fits in 4 bits.
//
// Examples:
//   REG_X15 => 15
//   REG_R9  => 9
//   REG_AX  => 0
//
func regIndex(r int16) int {
	lower3bits := reg[r]
	high4bit := regrex[r] & Rxr << 1
	return lower3bits | high4bit
}

// avx2gatherValid returns true if p satisfies AVX2 gather constraints.
// Reports errors via ctxt.
func avx2gatherValid(ctxt *obj.Link, p *obj.Prog) bool {
	// If any pair of the index, mask, or destination registers
	// are the same, this instruction results a #UD fault.
	index := regIndex(p.GetFrom3().Index)
	mask := regIndex(p.From.Reg)
	dest := regIndex(p.To.Reg)
	if dest == mask || dest == index || mask == index {
		ctxt.Diag("mask, index, and destination registers should be distinct: %v", p)
		return false
	}

	return true
}

func (asmbuf *AsmBuf) doasm(ctxt *obj.Link, cursym *obj.LSym, p *obj.Prog) {
	o := opindex[p.As&obj.AMask]

	if o == nil {
		ctxt.Diag("asmins: missing op %v", p)
		return
	}

	pre := prefixof(ctxt, p, &p.From)
	if pre != 0 {
		asmbuf.Put1(byte(pre))
	}
	pre = prefixof(ctxt, p, &p.To)
	if pre != 0 {
		asmbuf.Put1(byte(pre))
	}

	// TODO(rsc): This special case is for SHRQ $3, AX:DX,
	// which encodes as SHRQ $32(DX*0), AX.
	// Similarly SHRQ CX, AX:DX is really SHRQ CX(DX*0), AX.
	// Change encoding generated by assemblers and compilers and remove.
	if (p.From.Type == obj.TYPE_CONST || p.From.Type == obj.TYPE_REG) && p.From.Index != REG_NONE && p.From.Scale == 0 {
		p.SetFrom3(obj.Addr{
			Type: obj.TYPE_REG,
			Reg:  p.From.Index,
		})
		p.From.Index = 0
	}

	// TODO(rsc): This special case is for PINSRQ etc, CMPSD etc.
	// Change encoding generated by assemblers and compilers (if any) and remove.
	switch p.As {
	case AIMUL3Q, APEXTRW, APINSRW, APINSRD, APINSRQ, APSHUFHW, APSHUFL, APSHUFW, ASHUFPD, ASHUFPS, AAESKEYGENASSIST, APSHUFD, APCLMULQDQ:
		if p.From3Type() == obj.TYPE_NONE {
			p.SetFrom3(p.From)
			p.From = obj.Addr{}
			p.From.Type = obj.TYPE_CONST
			p.From.Offset = p.To.Offset
			p.To.Offset = 0
		}
	case ACMPSD, ACMPSS, ACMPPS, ACMPPD:
		if p.From3Type() == obj.TYPE_NONE {
			p.SetFrom3(p.To)
			p.To = obj.Addr{}
			p.To.Type = obj.TYPE_CONST
			p.To.Offset = p.GetFrom3().Offset
			p.GetFrom3().Offset = 0
		}

	case AVGATHERDPD,
		AVGATHERQPD,
		AVGATHERDPS,
		AVGATHERQPS,
		AVPGATHERDD,
		AVPGATHERQD,
		AVPGATHERDQ,
		AVPGATHERQQ:
		if !avx2gatherValid(ctxt, p) {
			return
		}
	}

	if p.Ft == 0 {
		p.Ft = uint8(oclass(ctxt, p, &p.From))
	}
	if p.Tt == 0 {
		p.Tt = uint8(oclass(ctxt, p, &p.To))
	}

	ft := int(p.Ft) * Ymax
	var f3t int
	tt := int(p.Tt) * Ymax

	xo := obj.Bool2int(o.op[0] == 0x0f)
	z := 0
	var a *obj.Addr
	var l int
	var op int
	var q *obj.Prog
	var r *obj.Reloc
	var rel obj.Reloc
	var v int64

	args := make([]int, 0, 6)
	if ft != Ynone*Ymax {
		args = append(args, ft)
	}
	for i := range p.RestArgs {
		args = append(args, oclass(ctxt, p, &p.RestArgs[i])*Ymax)
	}
	if tt != Ynone*Ymax {
		args = append(args, tt)
	}

	for _, yt := range o.ytab {
		if !yt.match(args) {
			z += int(yt.zoffset) + xo
		} else {
			switch o.prefix {
			case Px1: /* first option valid only in 32-bit mode */
				if ctxt.Arch.Family == sys.AMD64 && z == 0 {
					z += int(yt.zoffset) + xo
					continue
				}
			case Pq: /* 16 bit escape and opcode escape */
				asmbuf.Put2(Pe, Pm)

			case Pq3: /* 16 bit escape and opcode escape + REX.W */
				asmbuf.rexflag |= Pw
				asmbuf.Put2(Pe, Pm)

			case Pq4: /*  66 0F 38 */
				asmbuf.Put3(0x66, 0x0F, 0x38)

			case Pq4w: /*  66 0F 38 + REX.W */
				asmbuf.rexflag |= Pw
				asmbuf.Put3(0x66, 0x0F, 0x38)

			case Pq5: /*  F3 0F 38 */
				asmbuf.Put3(0xF3, 0x0F, 0x38)

			case Pq5w: /*  F3 0F 38 + REX.W */
				asmbuf.rexflag |= Pw
				asmbuf.Put3(0xF3, 0x0F, 0x38)

			case Pf2, /* xmm opcode escape */
				Pf3:
				asmbuf.Put2(o.prefix, Pm)

			case Pef3:
				asmbuf.Put3(Pe, Pf3, Pm)

			case Pfw: /* xmm opcode escape + REX.W */
				asmbuf.rexflag |= Pw
				asmbuf.Put2(Pf3, Pm)

			case Pm: /* opcode escape */
				asmbuf.Put1(Pm)

			case Pe: /* 16 bit escape */
				asmbuf.Put1(Pe)

			case Pw: /* 64-bit escape */
				if ctxt.Arch.Family != sys.AMD64 {
					ctxt.Diag("asmins: illegal 64: %v", p)
				}
				asmbuf.rexflag |= Pw

			case Pw8: /* 64-bit escape if z >= 8 */
				if z >= 8 {
					if ctxt.Arch.Family != sys.AMD64 {
						ctxt.Diag("asmins: illegal 64: %v", p)
					}
					asmbuf.rexflag |= Pw
				}

			case Pb: /* botch */
				if ctxt.Arch.Family != sys.AMD64 && (isbadbyte(&p.From) || isbadbyte(&p.To)) {
					goto bad
				}
				// NOTE(rsc): This is probably safe to do always,
				// but when enabled it chooses different encodings
				// than the old cmd/internal/obj/i386 code did,
				// which breaks our "same bits out" checks.
				// In particular, CMPB AX, $0 encodes as 80 f8 00
				// in the original obj/i386, and it would encode
				// (using a valid, shorter form) as 3c 00 if we enabled
				// the call to bytereg here.
				if ctxt.Arch.Family == sys.AMD64 {
					bytereg(&p.From, &p.Ft)
					bytereg(&p.To, &p.Tt)
				}

			case P32: /* 32 bit but illegal if 64-bit mode */
				if ctxt.Arch.Family == sys.AMD64 {
					ctxt.Diag("asmins: illegal in 64-bit mode: %v", p)
				}

			case Py: /* 64-bit only, no prefix */
				if ctxt.Arch.Family != sys.AMD64 {
					ctxt.Diag("asmins: illegal in %d-bit mode: %v", ctxt.Arch.RegSize*8, p)
				}

			case Py1: /* 64-bit only if z < 1, no prefix */
				if z < 1 && ctxt.Arch.Family != sys.AMD64 {
					ctxt.Diag("asmins: illegal in %d-bit mode: %v", ctxt.Arch.RegSize*8, p)
				}

			case Py3: /* 64-bit only if z < 3, no prefix */
				if z < 3 && ctxt.Arch.Family != sys.AMD64 {
					ctxt.Diag("asmins: illegal in %d-bit mode: %v", ctxt.Arch.RegSize*8, p)
				}
			}

			if z >= len(o.op) {
				log.Fatalf("asmins bad table %v", p)
			}
			op = int(o.op[z])
			// In vex case 0x0f is actually VEX_256_F2_0F_WIG
			if op == 0x0f && o.prefix != Pvex {
				asmbuf.Put1(byte(op))
				z++
				op = int(o.op[z])
			}

			switch yt.zcase {
			default:
				ctxt.Diag("asmins: unknown z %d %v", yt.zcase, p)
				return

			case Zpseudo:
				break

			case Zlit:
				for ; ; z++ {
					op = int(o.op[z])
					if op == 0 {
						break
					}
					asmbuf.Put1(byte(op))
				}

			case Zlitm_r:
				for ; ; z++ {
					op = int(o.op[z])
					if op == 0 {
						break
					}
					asmbuf.Put1(byte(op))
				}
				asmbuf.asmand(ctxt, cursym, p, &p.From, &p.To)

			case Zmb_r:
				bytereg(&p.From, &p.Ft)
				fallthrough

			case Zm_r:
				asmbuf.Put1(byte(op))
				asmbuf.asmand(ctxt, cursym, p, &p.From, &p.To)

			case Zm2_r:
				asmbuf.Put2(byte(op), o.op[z+1])
				asmbuf.asmand(ctxt, cursym, p, &p.From, &p.To)

			case Zm_r_xm:
				asmbuf.mediaop(ctxt, o, op, int(yt.zoffset), z)
				asmbuf.asmand(ctxt, cursym, p, &p.From, &p.To)

			case Zm_r_xm_nr:
				asmbuf.rexflag = 0
				asmbuf.mediaop(ctxt, o, op, int(yt.zoffset), z)
				asmbuf.asmand(ctxt, cursym, p, &p.From, &p.To)

			case Zm_r_i_xm:
				asmbuf.mediaop(ctxt, o, op, int(yt.zoffset), z)
				asmbuf.asmand(ctxt, cursym, p, &p.From, p.GetFrom3())
				asmbuf.Put1(byte(p.To.Offset))

			case Zibm_r, Zibr_m:
				for {
					tmp1 := z
					z++
					op = int(o.op[tmp1])
					if op == 0 {
						break
					}
					asmbuf.Put1(byte(op))
				}
				if yt.zcase == Zibr_m {
					asmbuf.asmand(ctxt, cursym, p, &p.To, p.GetFrom3())
				} else {
					asmbuf.asmand(ctxt, cursym, p, p.GetFrom3(), &p.To)
				}
				asmbuf.Put1(byte(p.From.Offset))

			case Zaut_r:
				asmbuf.Put1(0x8d) // leal
				if p.From.Type != obj.TYPE_ADDR {
					ctxt.Diag("asmins: Zaut sb type ADDR")
				}
				p.From.Type = obj.TYPE_MEM
				asmbuf.asmand(ctxt, cursym, p, &p.From, &p.To)
				p.From.Type = obj.TYPE_ADDR

			case Zm_o:
				asmbuf.Put1(byte(op))
				asmbuf.asmando(ctxt, cursym, p, &p.From, int(o.op[z+1]))

			case Zr_m:
				asmbuf.Put1(byte(op))
				asmbuf.asmand(ctxt, cursym, p, &p.To, &p.From)

			case Zvex:
				asmbuf.asmvex(ctxt, &p.From, p.GetFrom3(), &p.To, o.op[z], o.op[z+1])

			case Zvex_rm_v_r:
				asmbuf.asmvex(ctxt, &p.From, p.GetFrom3(), &p.To, o.op[z], o.op[z+1])
				asmbuf.asmand(ctxt, cursym, p, &p.From, &p.To)

			case Zvex_rm_v_ro:
				asmbuf.asmvex(ctxt, &p.From, p.GetFrom3(), &p.To, o.op[z], o.op[z+1])
				asmbuf.asmando(ctxt, cursym, p, &p.From, int(o.op[z+2]))

			case Zvex_i_r_v:
				asmbuf.asmvex(ctxt, p.GetFrom3(), &p.To, nil, o.op[z], o.op[z+1])
				regnum := byte(0x7)
				if p.GetFrom3().Reg >= REG_X0 && p.GetFrom3().Reg <= REG_X15 {
					regnum &= byte(p.GetFrom3().Reg - REG_X0)
				} else {
					regnum &= byte(p.GetFrom3().Reg - REG_Y0)
				}
				asmbuf.Put1(byte(o.op[z+2]) | regnum)
				asmbuf.Put1(byte(p.From.Offset))

			case Zvex_i_rm_v_r:
				imm, from, from3, to := unpackOps4(p)
				asmbuf.asmvex(ctxt, from, from3, to, o.op[z], o.op[z+1])
				asmbuf.asmand(ctxt, cursym, p, from, to)
				asmbuf.Put1(byte(imm.Offset))

			case Zvex_i_rm_r:
				asmbuf.asmvex(ctxt, p.GetFrom3(), nil, &p.To, o.op[z], o.op[z+1])
				asmbuf.asmand(ctxt, cursym, p, p.GetFrom3(), &p.To)
				asmbuf.Put1(byte(p.From.Offset))

			case Zvex_v_rm_r:
				asmbuf.asmvex(ctxt, p.GetFrom3(), &p.From, &p.To, o.op[z], o.op[z+1])
				asmbuf.asmand(ctxt, cursym, p, p.GetFrom3(), &p.To)

			case Zvex_r_v_rm:
				asmbuf.asmvex(ctxt, &p.To, p.GetFrom3(), &p.From, o.op[z], o.op[z+1])
				asmbuf.asmand(ctxt, cursym, p, &p.To, &p.From)

			case Zvex_rm_r_vo:
				asmbuf.asmvex(ctxt, &p.From, &p.To, p.GetFrom3(), o.op[z], o.op[z+1])
				asmbuf.asmando(ctxt, cursym, p, &p.From, int(o.op[z+2]))

			case Zvex_i_r_rm:
				asmbuf.asmvex(ctxt, &p.To, nil, p.GetFrom3(), o.op[z], o.op[z+1])
				asmbuf.asmand(ctxt, cursym, p, &p.To, p.GetFrom3())
				asmbuf.Put1(byte(p.From.Offset))

			case Zvex_hr_rm_v_r:
				hr, from, from3, to := unpackOps4(p)
				asmbuf.asmvex(ctxt, from, from3, to, o.op[z], o.op[z+1])
				asmbuf.asmand(ctxt, cursym, p, from, to)
				asmbuf.Put1(byte(regIndex(hr.Reg) << 4))

			case Zr_m_xm:
				asmbuf.mediaop(ctxt, o, op, int(yt.zoffset), z)
				asmbuf.asmand(ctxt, cursym, p, &p.To, &p.From)

			case Zr_m_xm_nr:
				asmbuf.rexflag = 0
				asmbuf.mediaop(ctxt, o, op, int(yt.zoffset), z)
				asmbuf.asmand(ctxt, cursym, p, &p.To, &p.From)

			case Zo_m:
				asmbuf.Put1(byte(op))
				asmbuf.asmando(ctxt, cursym, p, &p.To, int(o.op[z+1]))

			case Zcallindreg:
				r = obj.Addrel(cursym)
				r.Off = int32(p.Pc)
				r.Type = objabi.R_CALLIND
				r.Siz = 0
				fallthrough

			case Zo_m64:
				asmbuf.Put1(byte(op))
				asmbuf.asmandsz(ctxt, cursym, p, &p.To, int(o.op[z+1]), 0, 1)

			case Zm_ibo:
				asmbuf.Put1(byte(op))
				asmbuf.asmando(ctxt, cursym, p, &p.From, int(o.op[z+1]))
				asmbuf.Put1(byte(vaddr(ctxt, p, &p.To, nil)))

			case Zibo_m:
				asmbuf.Put1(byte(op))
				asmbuf.asmando(ctxt, cursym, p, &p.To, int(o.op[z+1]))
				asmbuf.Put1(byte(vaddr(ctxt, p, &p.From, nil)))

			case Zibo_m_xm:
				z = asmbuf.mediaop(ctxt, o, op, int(yt.zoffset), z)
				asmbuf.asmando(ctxt, cursym, p, &p.To, int(o.op[z+1]))
				asmbuf.Put1(byte(vaddr(ctxt, p, &p.From, nil)))

			case Z_ib, Zib_:
				if yt.zcase == Zib_ {
					a = &p.From
				} else {
					a = &p.To
				}
				asmbuf.Put1(byte(op))
				if p.As == AXABORT {
					asmbuf.Put1(o.op[z+1])
				}
				asmbuf.Put1(byte(vaddr(ctxt, p, a, nil)))

			case Zib_rp:
				asmbuf.rexflag |= regrex[p.To.Reg] & (Rxb | 0x40)
				asmbuf.Put2(byte(op+reg[p.To.Reg]), byte(vaddr(ctxt, p, &p.From, nil)))

			case Zil_rp:
				asmbuf.rexflag |= regrex[p.To.Reg] & Rxb
				asmbuf.Put1(byte(op + reg[p.To.Reg]))
				if o.prefix == Pe {
					v = vaddr(ctxt, p, &p.From, nil)
					asmbuf.PutInt16(int16(v))
				} else {
					asmbuf.relput4(ctxt, cursym, p, &p.From)
				}

			case Zo_iw:
				asmbuf.Put1(byte(op))
				if p.From.Type != obj.TYPE_NONE {
					v = vaddr(ctxt, p, &p.From, nil)
					asmbuf.PutInt16(int16(v))
				}

			case Ziq_rp:
				v = vaddr(ctxt, p, &p.From, &rel)
				l = int(v >> 32)
				if l == 0 && rel.Siz != 8 {
					//p->mark |= 0100;
					//print("zero: %llux %v\n", v, p);
					asmbuf.rexflag &^= (0x40 | Rxw)

					asmbuf.rexflag |= regrex[p.To.Reg] & Rxb
					asmbuf.Put1(byte(0xb8 + reg[p.To.Reg]))
					if rel.Type != 0 {
						r = obj.Addrel(cursym)
						*r = rel
						r.Off = int32(p.Pc + int64(asmbuf.Len()))
					}

					asmbuf.PutInt32(int32(v))
				} else if l == -1 && uint64(v)&(uint64(1)<<31) != 0 { /* sign extend */

					//p->mark |= 0100;
					//print("sign: %llux %v\n", v, p);
					asmbuf.Put1(0xc7)
					asmbuf.asmando(ctxt, cursym, p, &p.To, 0)

					asmbuf.PutInt32(int32(v)) // need all 8
				} else {
					//print("all: %llux %v\n", v, p);
					asmbuf.rexflag |= regrex[p.To.Reg] & Rxb
					asmbuf.Put1(byte(op + reg[p.To.Reg]))
					if rel.Type != 0 {
						r = obj.Addrel(cursym)
						*r = rel
						r.Off = int32(p.Pc + int64(asmbuf.Len()))
					}

					asmbuf.PutInt64(v)
				}

			case Zib_rr:
				asmbuf.Put1(byte(op))
				asmbuf.asmand(ctxt, cursym, p, &p.To, &p.To)
				asmbuf.Put1(byte(vaddr(ctxt, p, &p.From, nil)))

			case Z_il, Zil_:
				if yt.zcase == Zil_ {
					a = &p.From
				} else {
					a = &p.To
				}
				asmbuf.Put1(byte(op))
				if o.prefix == Pe {
					v = vaddr(ctxt, p, a, nil)
					asmbuf.PutInt16(int16(v))
				} else {
					asmbuf.relput4(ctxt, cursym, p, a)
				}

			case Zm_ilo, Zilo_m:
				asmbuf.Put1(byte(op))
				if yt.zcase == Zilo_m {
					a = &p.From
					asmbuf.asmando(ctxt, cursym, p, &p.To, int(o.op[z+1]))
				} else {
					a = &p.To
					asmbuf.asmando(ctxt, cursym, p, &p.From, int(o.op[z+1]))
				}

				if o.prefix == Pe {
					v = vaddr(ctxt, p, a, nil)
					asmbuf.PutInt16(int16(v))
				} else {
					asmbuf.relput4(ctxt, cursym, p, a)
				}

			case Zil_rr:
				asmbuf.Put1(byte(op))
				asmbuf.asmand(ctxt, cursym, p, &p.To, &p.To)
				if o.prefix == Pe {
					v = vaddr(ctxt, p, &p.From, nil)
					asmbuf.PutInt16(int16(v))
				} else {
					asmbuf.relput4(ctxt, cursym, p, &p.From)
				}

			case Z_rp:
				asmbuf.rexflag |= regrex[p.To.Reg] & (Rxb | 0x40)
				asmbuf.Put1(byte(op + reg[p.To.Reg]))

			case Zrp_:
				asmbuf.rexflag |= regrex[p.From.Reg] & (Rxb | 0x40)
				asmbuf.Put1(byte(op + reg[p.From.Reg]))

			case Zcallcon, Zjmpcon:
				if yt.zcase == Zcallcon {
					asmbuf.Put1(byte(op))
				} else {
					asmbuf.Put1(o.op[z+1])
				}
				r = obj.Addrel(cursym)
				r.Off = int32(p.Pc + int64(asmbuf.Len()))
				r.Type = objabi.R_PCREL
				r.Siz = 4
				r.Add = p.To.Offset
				asmbuf.PutInt32(0)

			case Zcallind:
				asmbuf.Put2(byte(op), o.op[z+1])
				r = obj.Addrel(cursym)
				r.Off = int32(p.Pc + int64(asmbuf.Len()))
				if ctxt.Arch.Family == sys.AMD64 {
					r.Type = objabi.R_PCREL
				} else {
					r.Type = objabi.R_ADDR
				}
				r.Siz = 4
				r.Add = p.To.Offset
				r.Sym = p.To.Sym
				asmbuf.PutInt32(0)

			case Zcall, Zcallduff:
				if p.To.Sym == nil {
					ctxt.Diag("call without target")
					ctxt.DiagFlush()
					log.Fatalf("bad code")
				}

				if yt.zcase == Zcallduff && ctxt.Flag_dynlink {
					ctxt.Diag("directly calling duff when dynamically linking Go")
				}

				if ctxt.Framepointer_enabled && yt.zcase == Zcallduff && ctxt.Arch.Family == sys.AMD64 {
					// Maintain BP around call, since duffcopy/duffzero can't do it
					// (the call jumps into the middle of the function).
					// This makes it possible to see call sites for duffcopy/duffzero in
					// BP-based profiling tools like Linux perf (which is the
					// whole point of obj.Framepointer_enabled).
					// MOVQ BP, -16(SP)
					// LEAQ -16(SP), BP
					asmbuf.Put(bpduff1)
				}
				asmbuf.Put1(byte(op))
				r = obj.Addrel(cursym)
				r.Off = int32(p.Pc + int64(asmbuf.Len()))
				r.Sym = p.To.Sym
				r.Add = p.To.Offset
				r.Type = objabi.R_CALL
				r.Siz = 4
				asmbuf.PutInt32(0)

				if ctxt.Framepointer_enabled && yt.zcase == Zcallduff && ctxt.Arch.Family == sys.AMD64 {
					// Pop BP pushed above.
					// MOVQ 0(BP), BP
					asmbuf.Put(bpduff2)
				}

			// TODO: jump across functions needs reloc
			case Zbr, Zjmp, Zloop:
				if p.As == AXBEGIN {
					asmbuf.Put1(byte(op))
				}
				if p.To.Sym != nil {
					if yt.zcase != Zjmp {
						ctxt.Diag("branch to ATEXT")
						ctxt.DiagFlush()
						log.Fatalf("bad code")
					}

					asmbuf.Put1(o.op[z+1])
					r = obj.Addrel(cursym)
					r.Off = int32(p.Pc + int64(asmbuf.Len()))
					r.Sym = p.To.Sym
					r.Type = objabi.R_PCREL
					r.Siz = 4
					asmbuf.PutInt32(0)
					break
				}

				// Assumes q is in this function.
				// TODO: Check in input, preserve in brchain.

				// Fill in backward jump now.
				q = p.Pcond

				if q == nil {
					ctxt.Diag("jmp/branch/loop without target")
					ctxt.DiagFlush()
					log.Fatalf("bad code")
				}

				if p.Back&1 != 0 {
					v = q.Pc - (p.Pc + 2)
					if v >= -128 && p.As != AXBEGIN {
						if p.As == AJCXZL {
							asmbuf.Put1(0x67)
						}
						asmbuf.Put2(byte(op), byte(v))
					} else if yt.zcase == Zloop {
						ctxt.Diag("loop too far: %v", p)
					} else {
						v -= 5 - 2
						if p.As == AXBEGIN {
							v--
						}
						if yt.zcase == Zbr {
							asmbuf.Put1(0x0f)
							v--
						}

						asmbuf.Put1(o.op[z+1])
						asmbuf.PutInt32(int32(v))
					}

					break
				}

				// Annotate target; will fill in later.
				p.Forwd = q.Rel

				q.Rel = p
				if p.Back&2 != 0 && p.As != AXBEGIN { // short
					if p.As == AJCXZL {
						asmbuf.Put1(0x67)
					}
					asmbuf.Put2(byte(op), 0)
				} else if yt.zcase == Zloop {
					ctxt.Diag("loop too far: %v", p)
				} else {
					if yt.zcase == Zbr {
						asmbuf.Put1(0x0f)
					}
					asmbuf.Put1(o.op[z+1])
					asmbuf.PutInt32(0)
				}

				break

			/*
				v = q->pc - p->pc - 2;
				if((v >= -128 && v <= 127) || p->pc == -1 || q->pc == -1) {
					*ctxt->andptr++ = op;
					*ctxt->andptr++ = v;
				} else {
					v -= 5-2;
					if(yt.zcase == Zbr) {
						*ctxt->andptr++ = 0x0f;
						v--;
					}
					*ctxt->andptr++ = o->op[z+1];
					*ctxt->andptr++ = v;
					*ctxt->andptr++ = v>>8;
					*ctxt->andptr++ = v>>16;
					*ctxt->andptr++ = v>>24;
				}
			*/

			case Zbyte:
				v = vaddr(ctxt, p, &p.From, &rel)
				if rel.Siz != 0 {
					rel.Siz = uint8(op)
					r = obj.Addrel(cursym)
					*r = rel
					r.Off = int32(p.Pc + int64(asmbuf.Len()))
				}

				asmbuf.Put1(byte(v))
				if op > 1 {
					asmbuf.Put1(byte(v >> 8))
					if op > 2 {
						asmbuf.PutInt16(int16(v >> 16))
						if op > 4 {
							asmbuf.PutInt32(int32(v >> 32))
						}
					}
				}
			}

			return
		}
	}
	f3t = Ynone * Ymax
	if p.GetFrom3() != nil {
		f3t = oclass(ctxt, p, p.GetFrom3()) * Ymax
	}
	for mo := ymovtab; mo[0].as != 0; mo = mo[1:] {
		var pp obj.Prog
		var t []byte
		if p.As == mo[0].as {
			if ycover[ft+int(mo[0].ft)] != 0 && ycover[f3t+int(mo[0].f3t)] != 0 && ycover[tt+int(mo[0].tt)] != 0 {
				t = mo[0].op[:]
				switch mo[0].code {
				default:
					ctxt.Diag("asmins: unknown mov %d %v", mo[0].code, p)

				case 0: /* lit */
					for z = 0; t[z] != E; z++ {
						asmbuf.Put1(t[z])
					}

				case 1: /* r,m */
					asmbuf.Put1(t[0])
					asmbuf.asmando(ctxt, cursym, p, &p.To, int(t[1]))

				case 2: /* m,r */
					asmbuf.Put1(t[0])
					asmbuf.asmando(ctxt, cursym, p, &p.From, int(t[1]))

				case 3: /* r,m - 2op */
					asmbuf.Put2(t[0], t[1])
					asmbuf.asmando(ctxt, cursym, p, &p.To, int(t[2]))
					asmbuf.rexflag |= regrex[p.From.Reg] & (Rxr | 0x40)

				case 4: /* m,r - 2op */
					asmbuf.Put2(t[0], t[1])
					asmbuf.asmando(ctxt, cursym, p, &p.From, int(t[2]))
					asmbuf.rexflag |= regrex[p.To.Reg] & (Rxr | 0x40)

				case 5: /* load full pointer, trash heap */
					if t[0] != 0 {
						asmbuf.Put1(t[0])
					}
					switch p.To.Index {
					default:
						goto bad

					case REG_DS:
						asmbuf.Put1(0xc5)

					case REG_SS:
						asmbuf.Put2(0x0f, 0xb2)

					case REG_ES:
						asmbuf.Put1(0xc4)

					case REG_FS:
						asmbuf.Put2(0x0f, 0xb4)

					case REG_GS:
						asmbuf.Put2(0x0f, 0xb5)
					}

					asmbuf.asmand(ctxt, cursym, p, &p.From, &p.To)

				case 6: /* double shift */
					if t[0] == Pw {
						if ctxt.Arch.Family != sys.AMD64 {
							ctxt.Diag("asmins: illegal 64: %v", p)
						}
						asmbuf.rexflag |= Pw
						t = t[1:]
					} else if t[0] == Pe {
						asmbuf.Put1(Pe)
						t = t[1:]
					}

					switch p.From.Type {
					default:
						goto bad

					case obj.TYPE_CONST:
						asmbuf.Put2(0x0f, t[0])
						asmbuf.asmandsz(ctxt, cursym, p, &p.To, reg[p.GetFrom3().Reg], regrex[p.GetFrom3().Reg], 0)
						asmbuf.Put1(byte(p.From.Offset))

					case obj.TYPE_REG:
						switch p.From.Reg {
						default:
							goto bad

						case REG_CL, REG_CX:
							asmbuf.Put2(0x0f, t[1])
							asmbuf.asmandsz(ctxt, cursym, p, &p.To, reg[p.GetFrom3().Reg], regrex[p.GetFrom3().Reg], 0)
						}
					}

				// NOTE: The systems listed here are the ones that use the "TLS initial exec" model,
				// where you load the TLS base register into a register and then index off that
				// register to access the actual TLS variables. Systems that allow direct TLS access
				// are handled in prefixof above and should not be listed here.
				case 7: /* mov tls, r */
					if ctxt.Arch.Family == sys.AMD64 && p.As != AMOVQ || ctxt.Arch.Family == sys.I386 && p.As != AMOVL {
						ctxt.Diag("invalid load of TLS: %v", p)
					}

					if ctxt.Arch.Family == sys.I386 {
						// NOTE: The systems listed here are the ones that use the "TLS initial exec" model,
						// where you load the TLS base register into a register and then index off that
						// register to access the actual TLS variables. Systems that allow direct TLS access
						// are handled in prefixof above and should not be listed here.
						switch ctxt.Headtype {
						default:
							log.Fatalf("unknown TLS base location for %v", ctxt.Headtype)

						case objabi.Hlinux,
							objabi.Hnacl:
							if ctxt.Flag_shared {
								// Note that this is not generating the same insns as the other cases.
								//     MOV TLS, dst
								// becomes
								//     call __x86.get_pc_thunk.dst
								//     movl (gotpc + g@gotntpoff)(dst), dst
								// which is encoded as
								//     call __x86.get_pc_thunk.dst
								//     movq 0(dst), dst
								// and R_CALL & R_TLS_IE relocs. This all assumes the only tls variable we access
								// is g, which we can't check here, but will when we assemble the second
								// instruction.
								dst := p.To.Reg
								asmbuf.Put1(0xe8)
								r = obj.Addrel(cursym)
								r.Off = int32(p.Pc + int64(asmbuf.Len()))
								r.Type = objabi.R_CALL
								r.Siz = 4
								r.Sym = ctxt.Lookup("__x86.get_pc_thunk." + strings.ToLower(rconv(int(dst))))
								asmbuf.PutInt32(0)

								asmbuf.Put2(0x8B, byte(2<<6|reg[dst]|(reg[dst]<<3)))
								r = obj.Addrel(cursym)
								r.Off = int32(p.Pc + int64(asmbuf.Len()))
								r.Type = objabi.R_TLS_IE
								r.Siz = 4
								r.Add = 2
								asmbuf.PutInt32(0)
							} else {
								// ELF TLS base is 0(GS).
								pp.From = p.From

								pp.From.Type = obj.TYPE_MEM
								pp.From.Reg = REG_GS
								pp.From.Offset = 0
								pp.From.Index = REG_NONE
								pp.From.Scale = 0
								asmbuf.Put2(0x65, // GS
									0x8B)
								asmbuf.asmand(ctxt, cursym, p, &pp.From, &p.To)
							}
						case objabi.Hplan9:
							pp.From = obj.Addr{}
							pp.From.Type = obj.TYPE_MEM
							pp.From.Name = obj.NAME_EXTERN
							pp.From.Sym = plan9privates
							pp.From.Offset = 0
							pp.From.Index = REG_NONE
							asmbuf.Put1(0x8B)
							asmbuf.asmand(ctxt, cursym, p, &pp.From, &p.To)

						case objabi.Hwindows:
							// Windows TLS base is always 0x14(FS).
							pp.From = p.From

							pp.From.Type = obj.TYPE_MEM
							pp.From.Reg = REG_FS
							pp.From.Offset = 0x14
							pp.From.Index = REG_NONE
							pp.From.Scale = 0
							asmbuf.Put2(0x64, // FS
								0x8B)
							asmbuf.asmand(ctxt, cursym, p, &pp.From, &p.To)
						}
						break
					}

					switch ctxt.Headtype {
					default:
						log.Fatalf("unknown TLS base location for %v", ctxt.Headtype)

					case objabi.Hlinux:
						if !ctxt.Flag_shared {
							log.Fatalf("unknown TLS base location for linux without -shared")
						}
						// Note that this is not generating the same insn as the other cases.
						//     MOV TLS, R_to
						// becomes
						//     movq g@gottpoff(%rip), R_to
						// which is encoded as
						//     movq 0(%rip), R_to
						// and a R_TLS_IE reloc. This all assumes the only tls variable we access
						// is g, which we can't check here, but will when we assemble the second
						// instruction.
						asmbuf.rexflag = Pw | (regrex[p.To.Reg] & Rxr)

						asmbuf.Put2(0x8B, byte(0x05|(reg[p.To.Reg]<<3)))
						r = obj.Addrel(cursym)
						r.Off = int32(p.Pc + int64(asmbuf.Len()))
						r.Type = objabi.R_TLS_IE
						r.Siz = 4
						r.Add = -4
						asmbuf.PutInt32(0)

					case objabi.Hplan9:
						pp.From = obj.Addr{}
						pp.From.Type = obj.TYPE_MEM
						pp.From.Name = obj.NAME_EXTERN
						pp.From.Sym = plan9privates
						pp.From.Offset = 0
						pp.From.Index = REG_NONE
						asmbuf.rexflag |= Pw
						asmbuf.Put1(0x8B)
						asmbuf.asmand(ctxt, cursym, p, &pp.From, &p.To)

					case objabi.Hsolaris: // TODO(rsc): Delete Hsolaris from list. Should not use this code. See progedit in obj6.c.
						// TLS base is 0(FS).
						pp.From = p.From

						pp.From.Type = obj.TYPE_MEM
						pp.From.Name = obj.NAME_NONE
						pp.From.Reg = REG_NONE
						pp.From.Offset = 0
						pp.From.Index = REG_NONE
						pp.From.Scale = 0
						asmbuf.rexflag |= Pw
						asmbuf.Put2(0x64, // FS
							0x8B)
						asmbuf.asmand(ctxt, cursym, p, &pp.From, &p.To)

					case objabi.Hwindows:
						// Windows TLS base is always 0x28(GS).
						pp.From = p.From

						pp.From.Type = obj.TYPE_MEM
						pp.From.Name = obj.NAME_NONE
						pp.From.Reg = REG_GS
						pp.From.Offset = 0x28
						pp.From.Index = REG_NONE
						pp.From.Scale = 0
						asmbuf.rexflag |= Pw
						asmbuf.Put2(0x65, // GS
							0x8B)
						asmbuf.asmand(ctxt, cursym, p, &pp.From, &p.To)
					}
				}
				return
			}
		}
	}
	goto bad

bad:
	if ctxt.Arch.Family != sys.AMD64 {
		/*
		 * here, the assembly has failed.
		 * if its a byte instruction that has
		 * unaddressable registers, try to
		 * exchange registers and reissue the
		 * instruction with the operands renamed.
		 */
		pp := *p

		unbytereg(&pp.From, &pp.Ft)
		unbytereg(&pp.To, &pp.Tt)

		z := int(p.From.Reg)
		if p.From.Type == obj.TYPE_REG && z >= REG_BP && z <= REG_DI {
			// TODO(rsc): Use this code for x86-64 too. It has bug fixes not present in the amd64 code base.
			// For now, different to keep bit-for-bit compatibility.
			if ctxt.Arch.Family == sys.I386 {
				breg := byteswapreg(ctxt, &p.To)
				if breg != REG_AX {
					asmbuf.Put1(0x87) // xchg lhs,bx
					asmbuf.asmando(ctxt, cursym, p, &p.From, reg[breg])
					subreg(&pp, z, breg)
					asmbuf.doasm(ctxt, cursym, &pp)
					asmbuf.Put1(0x87) // xchg lhs,bx
					asmbuf.asmando(ctxt, cursym, p, &p.From, reg[breg])
				} else {
					asmbuf.Put1(byte(0x90 + reg[z])) // xchg lsh,ax
					subreg(&pp, z, REG_AX)
					asmbuf.doasm(ctxt, cursym, &pp)
					asmbuf.Put1(byte(0x90 + reg[z])) // xchg lsh,ax
				}
				return
			}

			if isax(&p.To) || p.To.Type == obj.TYPE_NONE {
				// We certainly don't want to exchange
				// with AX if the op is MUL or DIV.
				asmbuf.Put1(0x87) // xchg lhs,bx
				asmbuf.asmando(ctxt, cursym, p, &p.From, reg[REG_BX])
				subreg(&pp, z, REG_BX)
				asmbuf.doasm(ctxt, cursym, &pp)
				asmbuf.Put1(0x87) // xchg lhs,bx
				asmbuf.asmando(ctxt, cursym, p, &p.From, reg[REG_BX])
			} else {
				asmbuf.Put1(byte(0x90 + reg[z])) // xchg lsh,ax
				subreg(&pp, z, REG_AX)
				asmbuf.doasm(ctxt, cursym, &pp)
				asmbuf.Put1(byte(0x90 + reg[z])) // xchg lsh,ax
			}
			return
		}

		z = int(p.To.Reg)
		if p.To.Type == obj.TYPE_REG && z >= REG_BP && z <= REG_DI {
			// TODO(rsc): Use this code for x86-64 too. It has bug fixes not present in the amd64 code base.
			// For now, different to keep bit-for-bit compatibility.
			if ctxt.Arch.Family == sys.I386 {
				breg := byteswapreg(ctxt, &p.From)
				if breg != REG_AX {
					asmbuf.Put1(0x87) //xchg rhs,bx
					asmbuf.asmando(ctxt, cursym, p, &p.To, reg[breg])
					subreg(&pp, z, breg)
					asmbuf.doasm(ctxt, cursym, &pp)
					asmbuf.Put1(0x87) // xchg rhs,bx
					asmbuf.asmando(ctxt, cursym, p, &p.To, reg[breg])
				} else {
					asmbuf.Put1(byte(0x90 + reg[z])) // xchg rsh,ax
					subreg(&pp, z, REG_AX)
					asmbuf.doasm(ctxt, cursym, &pp)
					asmbuf.Put1(byte(0x90 + reg[z])) // xchg rsh,ax
				}
				return
			}

			if isax(&p.From) {
				asmbuf.Put1(0x87) // xchg rhs,bx
				asmbuf.asmando(ctxt, cursym, p, &p.To, reg[REG_BX])
				subreg(&pp, z, REG_BX)
				asmbuf.doasm(ctxt, cursym, &pp)
				asmbuf.Put1(0x87) // xchg rhs,bx
				asmbuf.asmando(ctxt, cursym, p, &p.To, reg[REG_BX])
			} else {
				asmbuf.Put1(byte(0x90 + reg[z])) // xchg rsh,ax
				subreg(&pp, z, REG_AX)
				asmbuf.doasm(ctxt, cursym, &pp)
				asmbuf.Put1(byte(0x90 + reg[z])) // xchg rsh,ax
			}
			return
		}
	}

	ctxt.Diag("invalid instruction: %v", p)
	//	ctxt.Diag("doasm: notfound ft=%d tt=%d %v %d %d", p.Ft, p.Tt, p, oclass(ctxt, p, &p.From), oclass(ctxt, p, &p.To))
	return
}

// byteswapreg returns a byte-addressable register (AX, BX, CX, DX)
// which is not referenced in a.
// If a is empty, it returns BX to account for MULB-like instructions
// that might use DX and AX.
func byteswapreg(ctxt *obj.Link, a *obj.Addr) int {
	cana, canb, canc, cand := true, true, true, true
	if a.Type == obj.TYPE_NONE {
		cana, cand = false, false
	}

	if a.Type == obj.TYPE_REG || ((a.Type == obj.TYPE_MEM || a.Type == obj.TYPE_ADDR) && a.Name == obj.NAME_NONE) {
		switch a.Reg {
		case REG_NONE:
			cana, cand = false, false
		case REG_AX, REG_AL, REG_AH:
			cana = false
		case REG_BX, REG_BL, REG_BH:
			canb = false
		case REG_CX, REG_CL, REG_CH:
			canc = false
		case REG_DX, REG_DL, REG_DH:
			cand = false
		}
	}

	if a.Type == obj.TYPE_MEM || a.Type == obj.TYPE_ADDR {
		switch a.Index {
		case REG_AX:
			cana = false
		case REG_BX:
			canb = false
		case REG_CX:
			canc = false
		case REG_DX:
			cand = false
		}
	}

	switch {
	case cana:
		return REG_AX
	case canb:
		return REG_BX
	case canc:
		return REG_CX
	case cand:
		return REG_DX
	default:
		ctxt.Diag("impossible byte register")
		ctxt.DiagFlush()
		log.Fatalf("bad code")
		return 0
	}
}

func isbadbyte(a *obj.Addr) bool {
	return a.Type == obj.TYPE_REG && (REG_BP <= a.Reg && a.Reg <= REG_DI || REG_BPB <= a.Reg && a.Reg <= REG_DIB)
}

var naclret = []uint8{
	0x5e, // POPL SI
	// 0x8b, 0x7d, 0x00, // MOVL (BP), DI - catch return to invalid address, for debugging
	0x83,
	0xe6,
	0xe0, // ANDL $~31, SI
	0x4c,
	0x01,
	0xfe, // ADDQ R15, SI
	0xff,
	0xe6, // JMP SI
}

var naclret8 = []uint8{
	0x5d, // POPL BP
	// 0x8b, 0x7d, 0x00, // MOVL (BP), DI - catch return to invalid address, for debugging
	0x83,
	0xe5,
	0xe0, // ANDL $~31, BP
	0xff,
	0xe5, // JMP BP
}

var naclspfix = []uint8{0x4c, 0x01, 0xfc} // ADDQ R15, SP

var naclbpfix = []uint8{0x4c, 0x01, 0xfd} // ADDQ R15, BP

var naclmovs = []uint8{
	0x89,
	0xf6, // MOVL SI, SI
	0x49,
	0x8d,
	0x34,
	0x37, // LEAQ (R15)(SI*1), SI
	0x89,
	0xff, // MOVL DI, DI
	0x49,
	0x8d,
	0x3c,
	0x3f, // LEAQ (R15)(DI*1), DI
}

var naclstos = []uint8{
	0x89,
	0xff, // MOVL DI, DI
	0x49,
	0x8d,
	0x3c,
	0x3f, // LEAQ (R15)(DI*1), DI
}

func (asmbuf *AsmBuf) nacltrunc(ctxt *obj.Link, reg int) {
	if reg >= REG_R8 {
		asmbuf.Put1(0x45)
	}
	reg = (reg - REG_AX) & 7
	asmbuf.Put2(0x89, byte(3<<6|reg<<3|reg))
}

func (asmbuf *AsmBuf) asmins(ctxt *obj.Link, cursym *obj.LSym, p *obj.Prog) {
	asmbuf.Reset()

	if ctxt.Headtype == objabi.Hnacl && ctxt.Arch.Family == sys.I386 {
		switch p.As {
		case obj.ARET:
			asmbuf.Put(naclret8)
			return

		case obj.ACALL,
			obj.AJMP:
			if p.To.Type == obj.TYPE_REG && REG_AX <= p.To.Reg && p.To.Reg <= REG_DI {
				asmbuf.Put3(0x83, byte(0xe0|(p.To.Reg-REG_AX)), 0xe0)
			}

		case AINT:
			asmbuf.Put1(0xf4)
			return
		}
	}

	if ctxt.Headtype == objabi.Hnacl && ctxt.Arch.Family == sys.AMD64 {
		if p.As == AREP {
			asmbuf.rep++
			return
		}

		if p.As == AREPN {
			asmbuf.repn++
			return
		}

		if p.As == ALOCK {
			asmbuf.lock = true
			return
		}

		if p.As != ALEAQ && p.As != ALEAL {
			if p.From.Index != REG_NONE && p.From.Scale > 0 {
				asmbuf.nacltrunc(ctxt, int(p.From.Index))
			}
			if p.To.Index != REG_NONE && p.To.Scale > 0 {
				asmbuf.nacltrunc(ctxt, int(p.To.Index))
			}
		}

		switch p.As {
		case obj.ARET:
			asmbuf.Put(naclret)
			return

		case obj.ACALL,
			obj.AJMP:
			if p.To.Type == obj.TYPE_REG && REG_AX <= p.To.Reg && p.To.Reg <= REG_DI {
				// ANDL $~31, reg
				asmbuf.Put3(0x83, byte(0xe0|(p.To.Reg-REG_AX)), 0xe0)
				// ADDQ R15, reg
				asmbuf.Put3(0x4c, 0x01, byte(0xf8|(p.To.Reg-REG_AX)))
			}

			if p.To.Type == obj.TYPE_REG && REG_R8 <= p.To.Reg && p.To.Reg <= REG_R15 {
				// ANDL $~31, reg
				asmbuf.Put4(0x41, 0x83, byte(0xe0|(p.To.Reg-REG_R8)), 0xe0)
				// ADDQ R15, reg
				asmbuf.Put3(0x4d, 0x01, byte(0xf8|(p.To.Reg-REG_R8)))
			}

		case AINT:
			asmbuf.Put1(0xf4)
			return

		case ASCASB,
			ASCASW,
			ASCASL,
			ASCASQ,
			ASTOSB,
			ASTOSW,
			ASTOSL,
			ASTOSQ:
			asmbuf.Put(naclstos)

		case AMOVSB, AMOVSW, AMOVSL, AMOVSQ:
			asmbuf.Put(naclmovs)
		}

		if asmbuf.rep != 0 {
			asmbuf.Put1(0xf3)
			asmbuf.rep = 0
		}

		if asmbuf.repn != 0 {
			asmbuf.Put1(0xf2)
			asmbuf.repn = 0
		}

		if asmbuf.lock {
			asmbuf.Put1(0xf0)
			asmbuf.lock = false
		}
	}

	asmbuf.rexflag = 0
	asmbuf.vexflag = 0
	mark := asmbuf.Len()
	asmbuf.doasm(ctxt, cursym, p)
	if asmbuf.rexflag != 0 && asmbuf.vexflag == 0 {
		/*
		 * as befits the whole approach of the architecture,
		 * the rex prefix must appear before the first opcode byte
		 * (and thus after any 66/67/f2/f3/26/2e/3e prefix bytes, but
		 * before the 0f opcode escape!), or it might be ignored.
		 * note that the handbook often misleadingly shows 66/f2/f3 in `opcode'.
		 */
		if ctxt.Arch.Family != sys.AMD64 {
			ctxt.Diag("asmins: illegal in mode %d: %v (%d %d)", ctxt.Arch.RegSize*8, p, p.Ft, p.Tt)
		}
		n := asmbuf.Len()
		var np int
		for np = mark; np < n; np++ {
			c := asmbuf.At(np)
			if c != 0xf2 && c != 0xf3 && (c < 0x64 || c > 0x67) && c != 0x2e && c != 0x3e && c != 0x26 {
				break
			}
		}
		asmbuf.Insert(np, byte(0x40|asmbuf.rexflag))
	}

	n := asmbuf.Len()
	for i := len(cursym.R) - 1; i >= 0; i-- {
		r := &cursym.R[i]
		if int64(r.Off) < p.Pc {
			break
		}
		if asmbuf.rexflag != 0 && asmbuf.vexflag == 0 {
			r.Off++
		}
		if r.Type == objabi.R_PCREL {
			if ctxt.Arch.Family == sys.AMD64 || p.As == obj.AJMP || p.As == obj.ACALL {
				// PC-relative addressing is relative to the end of the instruction,
				// but the relocations applied by the linker are relative to the end
				// of the relocation. Because immediate instruction
				// arguments can follow the PC-relative memory reference in the
				// instruction encoding, the two may not coincide. In this case,
				// adjust addend so that linker can keep relocating relative to the
				// end of the relocation.
				r.Add -= p.Pc + int64(n) - (int64(r.Off) + int64(r.Siz))
			} else if ctxt.Arch.Family == sys.I386 {
				// On 386 PC-relative addressing (for non-call/jmp instructions)
				// assumes that the previous instruction loaded the PC of the end
				// of that instruction into CX, so the adjustment is relative to
				// that.
				r.Add += int64(r.Off) - p.Pc + int64(r.Siz)
			}
		}
		if r.Type == objabi.R_GOTPCREL && ctxt.Arch.Family == sys.I386 {
			// On 386, R_GOTPCREL makes the same assumptions as R_PCREL.
			r.Add += int64(r.Off) - p.Pc + int64(r.Siz)
		}

	}

	if ctxt.Arch.Family == sys.AMD64 && ctxt.Headtype == objabi.Hnacl && p.As != ACMPL && p.As != ACMPQ && p.To.Type == obj.TYPE_REG {
		switch p.To.Reg {
		case REG_SP:
			asmbuf.Put(naclspfix)
		case REG_BP:
			asmbuf.Put(naclbpfix)
		}
	}
}

// Extract 4 operands from p.
func unpackOps4(p *obj.Prog) (*obj.Addr, *obj.Addr, *obj.Addr, *obj.Addr) {
	return &p.From, &p.RestArgs[0], &p.RestArgs[1], &p.To
}
