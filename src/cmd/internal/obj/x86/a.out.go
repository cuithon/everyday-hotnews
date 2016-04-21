// Inferno utils/6c/6.out.h
// http://code.google.com/p/inferno-os/source/browse/utils/6c/6.out.h
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

package x86

import "cmd/internal/obj"

//go:generate go run ../stringer.go -i $GOFILE -o anames.go -p x86

const (
	/* mark flags */
	DONE          = 1 << iota
	PRESERVEFLAGS // not allowed to clobber flags
)

/*
 *	amd64
 */
const (
	AAAA = obj.ABaseAMD64 + obj.A_ARCHSPECIFIC + iota
	AAAD
	AAAM
	AAAS
	AADCB
	AADCL
	AADCW
	AADDB
	AADDL
	AADDW
	AADJSP
	AANDB
	AANDL
	AANDW
	AARPL
	ABOUNDL
	ABOUNDW
	ABSFL
	ABSFW
	ABSRL
	ABSRW
	ABTL
	ABTW
	ABTCL
	ABTCW
	ABTRL
	ABTRW
	ABTSL
	ABTSW
	ABYTE
	ACLC
	ACLD
	ACLI
	ACLTS
	ACMC
	ACMPB
	ACMPL
	ACMPW
	ACMPSB
	ACMPSL
	ACMPSW
	ADAA
	ADAS
	ADECB
	ADECL
	ADECQ
	ADECW
	ADIVB
	ADIVL
	ADIVW
	AENTER
	AHADDPD
	AHADDPS
	AHLT
	AHSUBPD
	AHSUBPS
	AIDIVB
	AIDIVL
	AIDIVW
	AIMULB
	AIMULL
	AIMULW
	AINB
	AINL
	AINW
	AINCB
	AINCL
	AINCQ
	AINCW
	AINSB
	AINSL
	AINSW
	AINT
	AINTO
	AIRETL
	AIRETW
	AJCC // >= unsigned
	AJCS // < unsigned
	AJCXZL
	AJEQ // == (zero)
	AJGE // >= signed
	AJGT // > signed
	AJHI // > unsigned
	AJLE // <= signed
	AJLS // <= unsigned
	AJLT // < signed
	AJMI // sign bit set (negative)
	AJNE // != (nonzero)
	AJOC // overflow clear
	AJOS // overflow set
	AJPC // parity clear
	AJPL // sign bit clear (positive)
	AJPS // parity set
	ALAHF
	ALARL
	ALARW
	ALEAL
	ALEAW
	ALEAVEL
	ALEAVEW
	ALOCK
	ALODSB
	ALODSL
	ALODSW
	ALONG
	ALOOP
	ALOOPEQ
	ALOOPNE
	ALSLL
	ALSLW
	AMOVB
	AMOVL
	AMOVW
	AMOVBLSX
	AMOVBLZX
	AMOVBQSX
	AMOVBQZX
	AMOVBWSX
	AMOVBWZX
	AMOVWLSX
	AMOVWLZX
	AMOVWQSX
	AMOVWQZX
	AMOVSB
	AMOVSL
	AMOVSW
	AMULB
	AMULL
	AMULW
	ANEGB
	ANEGL
	ANEGW
	ANOTB
	ANOTL
	ANOTW
	AORB
	AORL
	AORW
	AOUTB
	AOUTL
	AOUTW
	AOUTSB
	AOUTSL
	AOUTSW
	APAUSE
	APOPAL
	APOPAW
	APOPCNTW
	APOPCNTL
	APOPCNTQ
	APOPFL
	APOPFW
	APOPL
	APOPW
	APUSHAL
	APUSHAW
	APUSHFL
	APUSHFW
	APUSHL
	APUSHW
	ARCLB
	ARCLL
	ARCLW
	ARCRB
	ARCRL
	ARCRW
	AREP
	AREPN
	AROLB
	AROLL
	AROLW
	ARORB
	ARORL
	ARORW
	ASAHF
	ASALB
	ASALL
	ASALW
	ASARB
	ASARL
	ASARW
	ASBBB
	ASBBL
	ASBBW
	ASCASB
	ASCASL
	ASCASW
	ASETCC
	ASETCS
	ASETEQ
	ASETGE
	ASETGT
	ASETHI
	ASETLE
	ASETLS
	ASETLT
	ASETMI
	ASETNE
	ASETOC
	ASETOS
	ASETPC
	ASETPL
	ASETPS
	ACDQ
	ACWD
	ASHLB
	ASHLL
	ASHLW
	ASHRB
	ASHRL
	ASHRW
	ASTC
	ASTD
	ASTI
	ASTOSB
	ASTOSL
	ASTOSW
	ASUBB
	ASUBL
	ASUBW
	ASYSCALL
	ATESTB
	ATESTL
	ATESTW
	AVERR
	AVERW
	AWAIT
	AWORD
	AXCHGB
	AXCHGL
	AXCHGW
	AXLAT
	AXORB
	AXORL
	AXORW

	AFMOVB
	AFMOVBP
	AFMOVD
	AFMOVDP
	AFMOVF
	AFMOVFP
	AFMOVL
	AFMOVLP
	AFMOVV
	AFMOVVP
	AFMOVW
	AFMOVWP
	AFMOVX
	AFMOVXP

	AFCOMD
	AFCOMDP
	AFCOMDPP
	AFCOMF
	AFCOMFP
	AFCOML
	AFCOMLP
	AFCOMW
	AFCOMWP
	AFUCOM
	AFUCOMP
	AFUCOMPP

	AFADDDP
	AFADDW
	AFADDL
	AFADDF
	AFADDD

	AFMULDP
	AFMULW
	AFMULL
	AFMULF
	AFMULD

	AFSUBDP
	AFSUBW
	AFSUBL
	AFSUBF
	AFSUBD

	AFSUBRDP
	AFSUBRW
	AFSUBRL
	AFSUBRF
	AFSUBRD

	AFDIVDP
	AFDIVW
	AFDIVL
	AFDIVF
	AFDIVD

	AFDIVRDP
	AFDIVRW
	AFDIVRL
	AFDIVRF
	AFDIVRD

	AFXCHD
	AFFREE

	AFLDCW
	AFLDENV
	AFRSTOR
	AFSAVE
	AFSTCW
	AFSTENV
	AFSTSW

	AF2XM1
	AFABS
	AFCHS
	AFCLEX
	AFCOS
	AFDECSTP
	AFINCSTP
	AFINIT
	AFLD1
	AFLDL2E
	AFLDL2T
	AFLDLG2
	AFLDLN2
	AFLDPI
	AFLDZ
	AFNOP
	AFPATAN
	AFPREM
	AFPREM1
	AFPTAN
	AFRNDINT
	AFSCALE
	AFSIN
	AFSINCOS
	AFSQRT
	AFTST
	AFXAM
	AFXTRACT
	AFYL2X
	AFYL2XP1

	// extra 32-bit operations
	ACMPXCHGB
	ACMPXCHGL
	ACMPXCHGW
	ACMPXCHG8B
	ACPUID
	AINVD
	AINVLPG
	ALFENCE
	AMFENCE
	AMOVNTIL
	ARDMSR
	ARDPMC
	ARDTSC
	ARSM
	ASFENCE
	ASYSRET
	AWBINVD
	AWRMSR
	AXADDB
	AXADDL
	AXADDW

	// conditional move
	ACMOVLCC
	ACMOVLCS
	ACMOVLEQ
	ACMOVLGE
	ACMOVLGT
	ACMOVLHI
	ACMOVLLE
	ACMOVLLS
	ACMOVLLT
	ACMOVLMI
	ACMOVLNE
	ACMOVLOC
	ACMOVLOS
	ACMOVLPC
	ACMOVLPL
	ACMOVLPS
	ACMOVQCC
	ACMOVQCS
	ACMOVQEQ
	ACMOVQGE
	ACMOVQGT
	ACMOVQHI
	ACMOVQLE
	ACMOVQLS
	ACMOVQLT
	ACMOVQMI
	ACMOVQNE
	ACMOVQOC
	ACMOVQOS
	ACMOVQPC
	ACMOVQPL
	ACMOVQPS
	ACMOVWCC
	ACMOVWCS
	ACMOVWEQ
	ACMOVWGE
	ACMOVWGT
	ACMOVWHI
	ACMOVWLE
	ACMOVWLS
	ACMOVWLT
	ACMOVWMI
	ACMOVWNE
	ACMOVWOC
	ACMOVWOS
	ACMOVWPC
	ACMOVWPL
	ACMOVWPS

	// 64-bit
	AADCQ
	AADDQ
	AANDQ
	ABSFQ
	ABSRQ
	ABTCQ
	ABTQ
	ABTRQ
	ABTSQ
	ACMPQ
	ACMPSQ
	ACMPXCHGQ
	ACQO
	ADIVQ
	AIDIVQ
	AIMULQ
	AIRETQ
	AJCXZQ
	ALEAQ
	ALEAVEQ
	ALODSQ
	AMOVQ
	AMOVLQSX
	AMOVLQZX
	AMOVNTIQ
	AMOVSQ
	AMULQ
	ANEGQ
	ANOTQ
	AORQ
	APOPFQ
	APOPQ
	APUSHFQ
	APUSHQ
	ARCLQ
	ARCRQ
	AROLQ
	ARORQ
	AQUAD
	ASALQ
	ASARQ
	ASBBQ
	ASCASQ
	ASHLQ
	ASHRQ
	ASTOSQ
	ASUBQ
	ATESTQ
	AXADDQ
	AXCHGQ
	AXORQ
	AXGETBV

	// media
	AADDPD
	AADDPS
	AADDSD
	AADDSS
	AANDNL
	AANDNQ
	AANDNPD
	AANDNPS
	AANDPD
	AANDPS
	ABEXTRL
	ABEXTRQ
	ABLSIL
	ABLSIQ
	ABLSMSKL
	ABLSMSKQ
	ABLSRL
	ABLSRQ
	ABZHIL
	ABZHIQ
	ACMPPD
	ACMPPS
	ACMPSD
	ACMPSS
	ACOMISD
	ACOMISS
	ACVTPD2PL
	ACVTPD2PS
	ACVTPL2PD
	ACVTPL2PS
	ACVTPS2PD
	ACVTPS2PL
	ACVTSD2SL
	ACVTSD2SQ
	ACVTSD2SS
	ACVTSL2SD
	ACVTSL2SS
	ACVTSQ2SD
	ACVTSQ2SS
	ACVTSS2SD
	ACVTSS2SL
	ACVTSS2SQ
	ACVTTPD2PL
	ACVTTPS2PL
	ACVTTSD2SL
	ACVTTSD2SQ
	ACVTTSS2SL
	ACVTTSS2SQ
	ADIVPD
	ADIVPS
	ADIVSD
	ADIVSS
	AEMMS
	AFXRSTOR
	AFXRSTOR64
	AFXSAVE
	AFXSAVE64
	ALDDQU
	ALDMXCSR
	AMASKMOVOU
	AMASKMOVQ
	AMAXPD
	AMAXPS
	AMAXSD
	AMAXSS
	AMINPD
	AMINPS
	AMINSD
	AMINSS
	AMOVAPD
	AMOVAPS
	AMOVOU
	AMOVHLPS
	AMOVHPD
	AMOVHPS
	AMOVLHPS
	AMOVLPD
	AMOVLPS
	AMOVMSKPD
	AMOVMSKPS
	AMOVNTO
	AMOVNTPD
	AMOVNTPS
	AMOVNTQ
	AMOVO
	AMOVQOZX
	AMOVSD
	AMOVSS
	AMOVUPD
	AMOVUPS
	AMULPD
	AMULPS
	AMULSD
	AMULSS
	AMULXL
	AMULXQ
	AORPD
	AORPS
	APACKSSLW
	APACKSSWB
	APACKUSWB
	APADDB
	APADDL
	APADDQ
	APADDSB
	APADDSW
	APADDUSB
	APADDUSW
	APADDW
	APAND
	APANDN
	APAVGB
	APAVGW
	APCMPEQB
	APCMPEQL
	APCMPEQW
	APCMPGTB
	APCMPGTL
	APCMPGTW
	APDEPL
	APDEPQ
	APEXTL
	APEXTQ
	APEXTRB
	APEXTRD
	APEXTRQ
	APEXTRW
	APHADDD
	APHADDSW
	APHADDW
	APHMINPOSUW
	APHSUBD
	APHSUBSW
	APHSUBW
	APINSRB
	APINSRD
	APINSRQ
	APINSRW
	APMADDWL
	APMAXSW
	APMAXUB
	APMINSW
	APMINUB
	APMOVMSKB
	APMOVSXBD
	APMOVSXBQ
	APMOVSXBW
	APMOVSXDQ
	APMOVSXWD
	APMOVSXWQ
	APMOVZXBD
	APMOVZXBQ
	APMOVZXBW
	APMOVZXDQ
	APMOVZXWD
	APMOVZXWQ
	APMULDQ
	APMULHUW
	APMULHW
	APMULLD
	APMULLW
	APMULULQ
	APOR
	APSADBW
	APSHUFB
	APSHUFHW
	APSHUFL
	APSHUFLW
	APSHUFW
	APSLLL
	APSLLO
	APSLLQ
	APSLLW
	APSRAL
	APSRAW
	APSRLL
	APSRLO
	APSRLQ
	APSRLW
	APSUBB
	APSUBL
	APSUBQ
	APSUBSB
	APSUBSW
	APSUBUSB
	APSUBUSW
	APSUBW
	APUNPCKHBW
	APUNPCKHLQ
	APUNPCKHQDQ
	APUNPCKHWL
	APUNPCKLBW
	APUNPCKLLQ
	APUNPCKLQDQ
	APUNPCKLWL
	APXOR
	ARCPPS
	ARCPSS
	ARSQRTPS
	ARSQRTSS
	ASARXL
	ASARXQ
	ASHLXL
	ASHLXQ
	ASHRXL
	ASHRXQ
	ASHUFPD
	ASHUFPS
	ASQRTPD
	ASQRTPS
	ASQRTSD
	ASQRTSS
	ASTMXCSR
	ASUBPD
	ASUBPS
	ASUBSD
	ASUBSS
	AUCOMISD
	AUCOMISS
	AUNPCKHPD
	AUNPCKHPS
	AUNPCKLPD
	AUNPCKLPS
	AXORPD
	AXORPS
	APCMPESTRI

	ARETFW
	ARETFL
	ARETFQ
	ASWAPGS

	AMODE
	ACRC32B
	ACRC32Q
	AIMUL3Q

	APREFETCHT0
	APREFETCHT1
	APREFETCHT2
	APREFETCHNTA

	AMOVQL
	ABSWAPL
	ABSWAPQ

	AAESENC
	AAESENCLAST
	AAESDEC
	AAESDECLAST
	AAESIMC
	AAESKEYGENASSIST

	AROUNDPS
	AROUNDSS
	AROUNDPD
	AROUNDSD

	APSHUFD
	APCLMULQDQ

	AVZEROUPPER
	AVMOVDQU
	AVMOVNTDQ
	AVMOVDQA
	AVPCMPEQB
	AVPXOR
	AVPMOVMSKB
	AVPAND
	AVPTEST
	AVPBROADCASTB

	// from 386
	AJCXZW
	AFCMOVCC
	AFCMOVCS
	AFCMOVEQ
	AFCMOVHI
	AFCMOVLS
	AFCMOVNE
	AFCMOVNU
	AFCMOVUN
	AFCOMI
	AFCOMIP
	AFUCOMI
	AFUCOMIP

	// TSX
	AXACQUIRE
	AXRELEASE
	AXBEGIN
	AXEND
	AXABORT
	AXTEST

	ALAST
)

const (
	REG_NONE = 0
)

const (
	REG_AL = obj.RBaseAMD64 + iota
	REG_CL
	REG_DL
	REG_BL
	REG_SPB
	REG_BPB
	REG_SIB
	REG_DIB
	REG_R8B
	REG_R9B
	REG_R10B
	REG_R11B
	REG_R12B
	REG_R13B
	REG_R14B
	REG_R15B

	REG_AX
	REG_CX
	REG_DX
	REG_BX
	REG_SP
	REG_BP
	REG_SI
	REG_DI
	REG_R8
	REG_R9
	REG_R10
	REG_R11
	REG_R12
	REG_R13
	REG_R14
	REG_R15

	REG_AH
	REG_CH
	REG_DH
	REG_BH

	REG_F0
	REG_F1
	REG_F2
	REG_F3
	REG_F4
	REG_F5
	REG_F6
	REG_F7

	REG_M0
	REG_M1
	REG_M2
	REG_M3
	REG_M4
	REG_M5
	REG_M6
	REG_M7

	REG_X0
	REG_X1
	REG_X2
	REG_X3
	REG_X4
	REG_X5
	REG_X6
	REG_X7
	REG_X8
	REG_X9
	REG_X10
	REG_X11
	REG_X12
	REG_X13
	REG_X14
	REG_X15

	REG_Y0
	REG_Y1
	REG_Y2
	REG_Y3
	REG_Y4
	REG_Y5
	REG_Y6
	REG_Y7
	REG_Y8
	REG_Y9
	REG_Y10
	REG_Y11
	REG_Y12
	REG_Y13
	REG_Y14
	REG_Y15

	REG_CS
	REG_SS
	REG_DS
	REG_ES
	REG_FS
	REG_GS

	REG_GDTR /* global descriptor table register */
	REG_IDTR /* interrupt descriptor table register */
	REG_LDTR /* local descriptor table register */
	REG_MSW  /* machine status word */
	REG_TASK /* task register */

	REG_CR0
	REG_CR1
	REG_CR2
	REG_CR3
	REG_CR4
	REG_CR5
	REG_CR6
	REG_CR7
	REG_CR8
	REG_CR9
	REG_CR10
	REG_CR11
	REG_CR12
	REG_CR13
	REG_CR14
	REG_CR15

	REG_DR0
	REG_DR1
	REG_DR2
	REG_DR3
	REG_DR4
	REG_DR5
	REG_DR6
	REG_DR7

	REG_TR0
	REG_TR1
	REG_TR2
	REG_TR3
	REG_TR4
	REG_TR5
	REG_TR6
	REG_TR7

	REG_TLS

	MAXREG

	REG_CR = REG_CR0
	REG_DR = REG_DR0
	REG_TR = REG_TR0

	REGARG   = -1
	REGRET   = REG_AX
	FREGRET  = REG_X0
	REGSP    = REG_SP
	REGCTXT  = REG_DX
	REGEXT   = REG_R15     /* compiler allocates external registers R15 down */
	FREGMIN  = REG_X0 + 5  /* first register variable */
	FREGEXT  = REG_X0 + 15 /* first external register */
	T_TYPE   = 1 << 0
	T_INDEX  = 1 << 1
	T_OFFSET = 1 << 2
	T_FCONST = 1 << 3
	T_SYM    = 1 << 4
	T_SCONST = 1 << 5
	T_64     = 1 << 6
	T_GOTYPE = 1 << 7
)
