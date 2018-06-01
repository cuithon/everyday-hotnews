// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on CRYPTOGAMS code with the following comment:
// # ====================================================================
// # Written by Andy Polyakov <appro@openssl.org> for the OpenSSL
// # project. The module is, however, dual licensed under OpenSSL and
// # CRYPTOGAMS licenses depending on where you obtain it. For further
// # details see http://www.openssl.org/~appro/cryptogams/.
// # ====================================================================

#include "textflag.h"

// SHA512 block routine. See sha512block.go for Go equivalent.
//
// The algorithm is detailed in FIPS 180-4:
//
//  https://csrc.nist.gov/publications/fips/fips180-4/fips-180-4.pdf
//
// Wt = Mt; for 0 <= t <= 15
// Wt = SIGMA1(Wt-2) + SIGMA0(Wt-15) + Wt-16; for 16 <= t <= 79
//
// a = H0
// b = H1
// c = H2
// d = H3
// e = H4
// f = H5
// g = H6
// h = H7
//
// for t = 0 to 79 {
//    T1 = h + BIGSIGMA1(e) + Ch(e,f,g) + Kt + Wt
//    T2 = BIGSIGMA0(a) + Maj(a,b,c)
//    h = g
//    g = f
//    f = e
//    e = d + T1
//    d = c
//    c = b
//    b = a
//    a = T1 + T2
// }
//
// H0 = a + H0
// H1 = b + H1
// H2 = c + H2
// H3 = d + H3
// H4 = e + H4
// H5 = f + H5
// H6 = g + H6
// H7 = h + H7

#define CTX	R3
#define INP	R4
#define END	R5
#define TBL	R6
#define IDX	R7
#define CNT	R8
#define LEN	R9
#define OFFLOAD	R11
#define TEMP	R12

#define HEX00	R0
#define HEX10	R10
#define HEX20	R25
#define HEX30	R26
#define HEX40	R27
#define HEX50	R28
#define HEX60	R29
#define HEX70	R31

// V0-V7 are A-H
// V8-V23 are used for the message schedule
#define KI	V24
#define FUNC	V25
#define S0	V26
#define S1	V27
#define s0	V28
#define s1	V29
#define LEMASK	V31	// Permutation control register for little endian

// 2 copies of each Kt, to fill both doublewords of a vector register
DATA  ·kcon+0x000(SB)/8, $0x428a2f98d728ae22
DATA  ·kcon+0x008(SB)/8, $0x428a2f98d728ae22
DATA  ·kcon+0x010(SB)/8, $0x7137449123ef65cd
DATA  ·kcon+0x018(SB)/8, $0x7137449123ef65cd
DATA  ·kcon+0x020(SB)/8, $0xb5c0fbcfec4d3b2f
DATA  ·kcon+0x028(SB)/8, $0xb5c0fbcfec4d3b2f
DATA  ·kcon+0x030(SB)/8, $0xe9b5dba58189dbbc
DATA  ·kcon+0x038(SB)/8, $0xe9b5dba58189dbbc
DATA  ·kcon+0x040(SB)/8, $0x3956c25bf348b538
DATA  ·kcon+0x048(SB)/8, $0x3956c25bf348b538
DATA  ·kcon+0x050(SB)/8, $0x59f111f1b605d019
DATA  ·kcon+0x058(SB)/8, $0x59f111f1b605d019
DATA  ·kcon+0x060(SB)/8, $0x923f82a4af194f9b
DATA  ·kcon+0x068(SB)/8, $0x923f82a4af194f9b
DATA  ·kcon+0x070(SB)/8, $0xab1c5ed5da6d8118
DATA  ·kcon+0x078(SB)/8, $0xab1c5ed5da6d8118
DATA  ·kcon+0x080(SB)/8, $0xd807aa98a3030242
DATA  ·kcon+0x088(SB)/8, $0xd807aa98a3030242
DATA  ·kcon+0x090(SB)/8, $0x12835b0145706fbe
DATA  ·kcon+0x098(SB)/8, $0x12835b0145706fbe
DATA  ·kcon+0x0A0(SB)/8, $0x243185be4ee4b28c
DATA  ·kcon+0x0A8(SB)/8, $0x243185be4ee4b28c
DATA  ·kcon+0x0B0(SB)/8, $0x550c7dc3d5ffb4e2
DATA  ·kcon+0x0B8(SB)/8, $0x550c7dc3d5ffb4e2
DATA  ·kcon+0x0C0(SB)/8, $0x72be5d74f27b896f
DATA  ·kcon+0x0C8(SB)/8, $0x72be5d74f27b896f
DATA  ·kcon+0x0D0(SB)/8, $0x80deb1fe3b1696b1
DATA  ·kcon+0x0D8(SB)/8, $0x80deb1fe3b1696b1
DATA  ·kcon+0x0E0(SB)/8, $0x9bdc06a725c71235
DATA  ·kcon+0x0E8(SB)/8, $0x9bdc06a725c71235
DATA  ·kcon+0x0F0(SB)/8, $0xc19bf174cf692694
DATA  ·kcon+0x0F8(SB)/8, $0xc19bf174cf692694
DATA  ·kcon+0x100(SB)/8, $0xe49b69c19ef14ad2
DATA  ·kcon+0x108(SB)/8, $0xe49b69c19ef14ad2
DATA  ·kcon+0x110(SB)/8, $0xefbe4786384f25e3
DATA  ·kcon+0x118(SB)/8, $0xefbe4786384f25e3
DATA  ·kcon+0x120(SB)/8, $0x0fc19dc68b8cd5b5
DATA  ·kcon+0x128(SB)/8, $0x0fc19dc68b8cd5b5
DATA  ·kcon+0x130(SB)/8, $0x240ca1cc77ac9c65
DATA  ·kcon+0x138(SB)/8, $0x240ca1cc77ac9c65
DATA  ·kcon+0x140(SB)/8, $0x2de92c6f592b0275
DATA  ·kcon+0x148(SB)/8, $0x2de92c6f592b0275
DATA  ·kcon+0x150(SB)/8, $0x4a7484aa6ea6e483
DATA  ·kcon+0x158(SB)/8, $0x4a7484aa6ea6e483
DATA  ·kcon+0x160(SB)/8, $0x5cb0a9dcbd41fbd4
DATA  ·kcon+0x168(SB)/8, $0x5cb0a9dcbd41fbd4
DATA  ·kcon+0x170(SB)/8, $0x76f988da831153b5
DATA  ·kcon+0x178(SB)/8, $0x76f988da831153b5
DATA  ·kcon+0x180(SB)/8, $0x983e5152ee66dfab
DATA  ·kcon+0x188(SB)/8, $0x983e5152ee66dfab
DATA  ·kcon+0x190(SB)/8, $0xa831c66d2db43210
DATA  ·kcon+0x198(SB)/8, $0xa831c66d2db43210
DATA  ·kcon+0x1A0(SB)/8, $0xb00327c898fb213f
DATA  ·kcon+0x1A8(SB)/8, $0xb00327c898fb213f
DATA  ·kcon+0x1B0(SB)/8, $0xbf597fc7beef0ee4
DATA  ·kcon+0x1B8(SB)/8, $0xbf597fc7beef0ee4
DATA  ·kcon+0x1C0(SB)/8, $0xc6e00bf33da88fc2
DATA  ·kcon+0x1C8(SB)/8, $0xc6e00bf33da88fc2
DATA  ·kcon+0x1D0(SB)/8, $0xd5a79147930aa725
DATA  ·kcon+0x1D8(SB)/8, $0xd5a79147930aa725
DATA  ·kcon+0x1E0(SB)/8, $0x06ca6351e003826f
DATA  ·kcon+0x1E8(SB)/8, $0x06ca6351e003826f
DATA  ·kcon+0x1F0(SB)/8, $0x142929670a0e6e70
DATA  ·kcon+0x1F8(SB)/8, $0x142929670a0e6e70
DATA  ·kcon+0x200(SB)/8, $0x27b70a8546d22ffc
DATA  ·kcon+0x208(SB)/8, $0x27b70a8546d22ffc
DATA  ·kcon+0x210(SB)/8, $0x2e1b21385c26c926
DATA  ·kcon+0x218(SB)/8, $0x2e1b21385c26c926
DATA  ·kcon+0x220(SB)/8, $0x4d2c6dfc5ac42aed
DATA  ·kcon+0x228(SB)/8, $0x4d2c6dfc5ac42aed
DATA  ·kcon+0x230(SB)/8, $0x53380d139d95b3df
DATA  ·kcon+0x238(SB)/8, $0x53380d139d95b3df
DATA  ·kcon+0x240(SB)/8, $0x650a73548baf63de
DATA  ·kcon+0x248(SB)/8, $0x650a73548baf63de
DATA  ·kcon+0x250(SB)/8, $0x766a0abb3c77b2a8
DATA  ·kcon+0x258(SB)/8, $0x766a0abb3c77b2a8
DATA  ·kcon+0x260(SB)/8, $0x81c2c92e47edaee6
DATA  ·kcon+0x268(SB)/8, $0x81c2c92e47edaee6
DATA  ·kcon+0x270(SB)/8, $0x92722c851482353b
DATA  ·kcon+0x278(SB)/8, $0x92722c851482353b
DATA  ·kcon+0x280(SB)/8, $0xa2bfe8a14cf10364
DATA  ·kcon+0x288(SB)/8, $0xa2bfe8a14cf10364
DATA  ·kcon+0x290(SB)/8, $0xa81a664bbc423001
DATA  ·kcon+0x298(SB)/8, $0xa81a664bbc423001
DATA  ·kcon+0x2A0(SB)/8, $0xc24b8b70d0f89791
DATA  ·kcon+0x2A8(SB)/8, $0xc24b8b70d0f89791
DATA  ·kcon+0x2B0(SB)/8, $0xc76c51a30654be30
DATA  ·kcon+0x2B8(SB)/8, $0xc76c51a30654be30
DATA  ·kcon+0x2C0(SB)/8, $0xd192e819d6ef5218
DATA  ·kcon+0x2C8(SB)/8, $0xd192e819d6ef5218
DATA  ·kcon+0x2D0(SB)/8, $0xd69906245565a910
DATA  ·kcon+0x2D8(SB)/8, $0xd69906245565a910
DATA  ·kcon+0x2E0(SB)/8, $0xf40e35855771202a
DATA  ·kcon+0x2E8(SB)/8, $0xf40e35855771202a
DATA  ·kcon+0x2F0(SB)/8, $0x106aa07032bbd1b8
DATA  ·kcon+0x2F8(SB)/8, $0x106aa07032bbd1b8
DATA  ·kcon+0x300(SB)/8, $0x19a4c116b8d2d0c8
DATA  ·kcon+0x308(SB)/8, $0x19a4c116b8d2d0c8
DATA  ·kcon+0x310(SB)/8, $0x1e376c085141ab53
DATA  ·kcon+0x318(SB)/8, $0x1e376c085141ab53
DATA  ·kcon+0x320(SB)/8, $0x2748774cdf8eeb99
DATA  ·kcon+0x328(SB)/8, $0x2748774cdf8eeb99
DATA  ·kcon+0x330(SB)/8, $0x34b0bcb5e19b48a8
DATA  ·kcon+0x338(SB)/8, $0x34b0bcb5e19b48a8
DATA  ·kcon+0x340(SB)/8, $0x391c0cb3c5c95a63
DATA  ·kcon+0x348(SB)/8, $0x391c0cb3c5c95a63
DATA  ·kcon+0x350(SB)/8, $0x4ed8aa4ae3418acb
DATA  ·kcon+0x358(SB)/8, $0x4ed8aa4ae3418acb
DATA  ·kcon+0x360(SB)/8, $0x5b9cca4f7763e373
DATA  ·kcon+0x368(SB)/8, $0x5b9cca4f7763e373
DATA  ·kcon+0x370(SB)/8, $0x682e6ff3d6b2b8a3
DATA  ·kcon+0x378(SB)/8, $0x682e6ff3d6b2b8a3
DATA  ·kcon+0x380(SB)/8, $0x748f82ee5defb2fc
DATA  ·kcon+0x388(SB)/8, $0x748f82ee5defb2fc
DATA  ·kcon+0x390(SB)/8, $0x78a5636f43172f60
DATA  ·kcon+0x398(SB)/8, $0x78a5636f43172f60
DATA  ·kcon+0x3A0(SB)/8, $0x84c87814a1f0ab72
DATA  ·kcon+0x3A8(SB)/8, $0x84c87814a1f0ab72
DATA  ·kcon+0x3B0(SB)/8, $0x8cc702081a6439ec
DATA  ·kcon+0x3B8(SB)/8, $0x8cc702081a6439ec
DATA  ·kcon+0x3C0(SB)/8, $0x90befffa23631e28
DATA  ·kcon+0x3C8(SB)/8, $0x90befffa23631e28
DATA  ·kcon+0x3D0(SB)/8, $0xa4506cebde82bde9
DATA  ·kcon+0x3D8(SB)/8, $0xa4506cebde82bde9
DATA  ·kcon+0x3E0(SB)/8, $0xbef9a3f7b2c67915
DATA  ·kcon+0x3E8(SB)/8, $0xbef9a3f7b2c67915
DATA  ·kcon+0x3F0(SB)/8, $0xc67178f2e372532b
DATA  ·kcon+0x3F8(SB)/8, $0xc67178f2e372532b
DATA  ·kcon+0x400(SB)/8, $0xca273eceea26619c
DATA  ·kcon+0x408(SB)/8, $0xca273eceea26619c
DATA  ·kcon+0x410(SB)/8, $0xd186b8c721c0c207
DATA  ·kcon+0x418(SB)/8, $0xd186b8c721c0c207
DATA  ·kcon+0x420(SB)/8, $0xeada7dd6cde0eb1e
DATA  ·kcon+0x428(SB)/8, $0xeada7dd6cde0eb1e
DATA  ·kcon+0x430(SB)/8, $0xf57d4f7fee6ed178
DATA  ·kcon+0x438(SB)/8, $0xf57d4f7fee6ed178
DATA  ·kcon+0x440(SB)/8, $0x06f067aa72176fba
DATA  ·kcon+0x448(SB)/8, $0x06f067aa72176fba
DATA  ·kcon+0x450(SB)/8, $0x0a637dc5a2c898a6
DATA  ·kcon+0x458(SB)/8, $0x0a637dc5a2c898a6
DATA  ·kcon+0x460(SB)/8, $0x113f9804bef90dae
DATA  ·kcon+0x468(SB)/8, $0x113f9804bef90dae
DATA  ·kcon+0x470(SB)/8, $0x1b710b35131c471b
DATA  ·kcon+0x478(SB)/8, $0x1b710b35131c471b
DATA  ·kcon+0x480(SB)/8, $0x28db77f523047d84
DATA  ·kcon+0x488(SB)/8, $0x28db77f523047d84
DATA  ·kcon+0x490(SB)/8, $0x32caab7b40c72493
DATA  ·kcon+0x498(SB)/8, $0x32caab7b40c72493
DATA  ·kcon+0x4A0(SB)/8, $0x3c9ebe0a15c9bebc
DATA  ·kcon+0x4A8(SB)/8, $0x3c9ebe0a15c9bebc
DATA  ·kcon+0x4B0(SB)/8, $0x431d67c49c100d4c
DATA  ·kcon+0x4B8(SB)/8, $0x431d67c49c100d4c
DATA  ·kcon+0x4C0(SB)/8, $0x4cc5d4becb3e42b6
DATA  ·kcon+0x4C8(SB)/8, $0x4cc5d4becb3e42b6
DATA  ·kcon+0x4D0(SB)/8, $0x597f299cfc657e2a
DATA  ·kcon+0x4D8(SB)/8, $0x597f299cfc657e2a
DATA  ·kcon+0x4E0(SB)/8, $0x5fcb6fab3ad6faec
DATA  ·kcon+0x4E8(SB)/8, $0x5fcb6fab3ad6faec
DATA  ·kcon+0x4F0(SB)/8, $0x6c44198c4a475817
DATA  ·kcon+0x4F8(SB)/8, $0x6c44198c4a475817
DATA  ·kcon+0x500(SB)/8, $0x0000000000000000
DATA  ·kcon+0x508(SB)/8, $0x0000000000000000
DATA  ·kcon+0x510(SB)/8, $0x1011121314151617
DATA  ·kcon+0x518(SB)/8, $0x0001020304050607
GLOBL ·kcon(SB), RODATA, $1312

#define SHA512ROUND0(a, b, c, d, e, f, g, h, xi) \
	VSEL		g, f, e, FUNC; \
	VSHASIGMAD	$15, e, $1, S1; \
	VADDUDM		xi, h, h; \
	VSHASIGMAD	$0, a, $1, S0; \
	VADDUDM		FUNC, h, h; \
	VXOR		b, a, FUNC; \
	VADDUDM		S1, h, h; \
	VSEL		b, c, FUNC, FUNC; \
	VADDUDM		KI, g, g; \
	VADDUDM		h, d, d; \
	VADDUDM		FUNC, S0, S0; \
	LVX		(TBL)(IDX), KI; \
	ADD		$16, IDX; \
	VADDUDM		S0, h, h

#define SHA512ROUND1(a, b, c, d, e, f, g, h, xi, xj, xj_1, xj_9, xj_14) \
	VSHASIGMAD	$0, xj_1, $0, s0; \
	VSEL		g, f, e, FUNC; \
	VSHASIGMAD	$15, e, $1, S1; \
	VADDUDM		xi, h, h; \
	VSHASIGMAD	$0, a, $1, S0; \
	VSHASIGMAD	$15, xj_14, $0, s1; \
	VADDUDM		FUNC, h, h; \
	VXOR		b, a, FUNC; \
	VADDUDM		xj_9, xj, xj; \
	VADDUDM		S1, h, h; \
	VSEL		b, c, FUNC, FUNC; \
	VADDUDM		KI, g, g; \
	VADDUDM		h, d, d; \
	VADDUDM		FUNC, S0, S0; \
	VADDUDM		s0, xj, xj; \
	LVX		(TBL)(IDX), KI; \
	ADD		$16, IDX; \
	VADDUDM		S0, h, h; \
	VADDUDM		s1, xj, xj

// func block(dig *digest, p []byte)
TEXT ·block(SB),0,$128-32
	MOVD	dig+0(FP), CTX
	MOVD	p_base+8(FP), INP
	MOVD	p_len+16(FP), LEN

	SRD	$6, LEN
	SLD	$6, LEN

	ADD	INP, LEN, END

	CMP	INP, END
	BEQ	end

	MOVD	$·kcon(SB), TBL
	MOVD	R1, OFFLOAD

	MOVD	R0, CNT
	MOVWZ	$0x10, HEX10
	MOVWZ	$0x20, HEX20
	MOVWZ	$0x30, HEX30
	MOVWZ	$0x40, HEX40
	MOVWZ	$0x50, HEX50
	MOVWZ	$0x60, HEX60
	MOVWZ	$0x70, HEX70

	MOVWZ	$8, IDX
	LVSL	(IDX)(R0), LEMASK
	VSPLTISB	$0x0F, KI
	VXOR	KI, LEMASK, LEMASK

	LXVD2X	(CTX)(HEX00), VS32	// v0 = vs32
	LXVD2X	(CTX)(HEX10), VS34	// v2 = vs34
	LXVD2X	(CTX)(HEX20), VS36	// v4 = vs36
	// unpack the input values into vector registers
	VSLDOI	$8, V0, V0, V1
	LXVD2X	(CTX)(HEX30), VS38	// v6 = vs38
	VSLDOI	$8, V2, V2, V3
	VSLDOI	$8, V4, V4, V5
	VSLDOI	$8, V6, V6, V7

loop:
	LVX	(TBL)(HEX00), KI
	MOVWZ	$16, IDX

	LXVD2X	(INP)(R0), VS40	// load v8 (=vs40) in advance
	ADD	$16, INP

	STVX	V0, (OFFLOAD+HEX00)
	STVX	V1, (OFFLOAD+HEX10)
	STVX	V2, (OFFLOAD+HEX20)
	STVX	V3, (OFFLOAD+HEX30)
	STVX	V4, (OFFLOAD+HEX40)
	STVX	V5, (OFFLOAD+HEX50)
	STVX	V6, (OFFLOAD+HEX60)
	STVX	V7, (OFFLOAD+HEX70)

	VADDUDM	KI, V7, V7	// h+K[i]
	LVX	(TBL)(IDX), KI
	ADD	$16, IDX

	VPERM	V8, V8, LEMASK, V8
	SHA512ROUND0(V0, V1, V2, V3, V4, V5, V6, V7, V8)
	LXVD2X	(INP)(R0), VS42	// load v10 (=vs42) in advance
	ADD	$16, INP, INP
	VSLDOI	$8, V8, V8, V9
	SHA512ROUND0(V7, V0, V1, V2, V3, V4, V5, V6, V9)
	VPERM	V10, V10, LEMASK, V10
	SHA512ROUND0(V6, V7, V0, V1, V2, V3, V4, V5, V10)
	LXVD2X	(INP)(R0), VS44	// load v12 (=vs44) in advance
	ADD	$16, INP, INP
	VSLDOI	$8, V10, V10, V11
	SHA512ROUND0(V5, V6, V7, V0, V1, V2, V3, V4, V11)
	VPERM	V12, V12, LEMASK, V12
	SHA512ROUND0(V4, V5, V6, V7, V0, V1, V2, V3, V12)
	LXVD2X	(INP)(R0), VS46	// load v14 (=vs46) in advance
	ADD	$16, INP, INP
	VSLDOI	$8, V12, V12, V13
	SHA512ROUND0(V3, V4, V5, V6, V7, V0, V1, V2, V13)
	VPERM	V14, V14, LEMASK, V14
	SHA512ROUND0(V2, V3, V4, V5, V6, V7, V0, V1, V14)
	LXVD2X	(INP)(R0), VS48	// load v16 (=vs48) in advance
	ADD	$16, INP, INP
	VSLDOI	$8, V14, V14, V15
	SHA512ROUND0(V1, V2, V3, V4, V5, V6, V7, V0, V15)
	VPERM	V16, V16, LEMASK, V16
	SHA512ROUND0(V0, V1, V2, V3, V4, V5, V6, V7, V16)
	LXVD2X	(INP)(R0), VS50	// load v18 (=vs50) in advance
	ADD	$16, INP, INP
	VSLDOI	$8, V16, V16, V17
	SHA512ROUND0(V7, V0, V1, V2, V3, V4, V5, V6, V17)
	VPERM	V18, V18, LEMASK, V18
	SHA512ROUND0(V6, V7, V0, V1, V2, V3, V4, V5, V18)
	LXVD2X	(INP)(R0), VS52	// load v20 (=vs52) in advance
	ADD	$16, INP, INP
	VSLDOI	$8, V18, V18, V19
	SHA512ROUND0(V5, V6, V7, V0, V1, V2, V3, V4, V19)
	VPERM	V20, V20, LEMASK, V20
	SHA512ROUND0(V4, V5, V6, V7, V0, V1, V2, V3, V20)
	LXVD2X	(INP)(R0), VS54	// load v22 (=vs54) in advance
	ADD	$16, INP, INP
	VSLDOI	$8, V20, V20, V21
	SHA512ROUND0(V3, V4, V5, V6, V7, V0, V1, V2, V21)
	VPERM	V22, V22, LEMASK, V22
	SHA512ROUND0(V2, V3, V4, V5, V6, V7, V0, V1, V22)
	VSLDOI	$8, V22, V22, V23
	SHA512ROUND1(V1, V2, V3, V4, V5, V6, V7, V0, V23, V8, V9, V17, V22)

	MOVWZ	$4, TEMP
	MOVWZ	TEMP, CTR

L16_xx:
	SHA512ROUND1(V0, V1, V2, V3, V4, V5, V6, V7, V8, V9, V10, V18, V23)
	SHA512ROUND1(V7, V0, V1, V2, V3, V4, V5, V6, V9, V10, V11, V19, V8)
	SHA512ROUND1(V6, V7, V0, V1, V2, V3, V4, V5, V10, V11, V12, V20, V9)
	SHA512ROUND1(V5, V6, V7, V0, V1, V2, V3, V4, V11, V12, V13, V21, V10)
	SHA512ROUND1(V4, V5, V6, V7, V0, V1, V2, V3, V12, V13, V14, V22, V11)
	SHA512ROUND1(V3, V4, V5, V6, V7, V0, V1, V2, V13, V14, V15, V23, V12)
	SHA512ROUND1(V2, V3, V4, V5, V6, V7, V0, V1, V14, V15, V16, V8, V13)
	SHA512ROUND1(V1, V2, V3, V4, V5, V6, V7, V0, V15, V16, V17, V9, V14)
	SHA512ROUND1(V0, V1, V2, V3, V4, V5, V6, V7, V16, V17, V18, V10, V15)
	SHA512ROUND1(V7, V0, V1, V2, V3, V4, V5, V6, V17, V18, V19, V11, V16)
	SHA512ROUND1(V6, V7, V0, V1, V2, V3, V4, V5, V18, V19, V20, V12, V17)
	SHA512ROUND1(V5, V6, V7, V0, V1, V2, V3, V4, V19, V20, V21, V13, V18)
	SHA512ROUND1(V4, V5, V6, V7, V0, V1, V2, V3, V20, V21, V22, V14, V19)
	SHA512ROUND1(V3, V4, V5, V6, V7, V0, V1, V2, V21, V22, V23, V15, V20)
	SHA512ROUND1(V2, V3, V4, V5, V6, V7, V0, V1, V22, V23, V8, V16, V21)
	SHA512ROUND1(V1, V2, V3, V4, V5, V6, V7, V0, V23, V8, V9, V17, V22)

	BC	0x10, 0, L16_xx		// bdnz

	LVX	(OFFLOAD)(HEX00), V10

	LVX	(OFFLOAD)(HEX10), V11
	VADDUDM	V10, V0, V0
	LVX	(OFFLOAD)(HEX20), V12
	VADDUDM	V11, V1, V1
	LVX	(OFFLOAD)(HEX30), V13
	VADDUDM	V12, V2, V2
	LVX	(OFFLOAD)(HEX40), V14
	VADDUDM	V13, V3, V3
	LVX	(OFFLOAD)(HEX50), V15
	VADDUDM	V14, V4, V4
	LVX	(OFFLOAD)(HEX60), V16
	VADDUDM	V15, V5, V5
	LVX	(OFFLOAD)(HEX70), V17
	VADDUDM	V16, V6, V6
	VADDUDM	V17, V7, V7

	CMPU	INP, END
	BLT	loop

	VPERM	V0, V1, KI, V0
	VPERM	V2, V3, KI, V2
	VPERM	V4, V5, KI, V4
	VPERM	V6, V7, KI, V6
	STXVD2X	VS32, (CTX+HEX00)	// v0 = vs32
	STXVD2X	VS34, (CTX+HEX10)	// v2 = vs34
	STXVD2X	VS36, (CTX+HEX20)	// v4 = vs36
	STXVD2X	VS38, (CTX+HEX30)	// v6 = vs38

end:
	RET

