// Code generated by stringer -i a.out.go -o anames.go -p ppc64; DO NOT EDIT.

package ppc64

import "cmd/internal/obj"

var Anames = []string{
	obj.A_ARCHSPECIFIC: "ADD",
	"ADDCC",
	"ADDIS",
	"ADDV",
	"ADDVCC",
	"ADDC",
	"ADDCCC",
	"ADDCV",
	"ADDCVCC",
	"ADDME",
	"ADDMECC",
	"ADDMEVCC",
	"ADDMEV",
	"ADDE",
	"ADDECC",
	"ADDEVCC",
	"ADDEV",
	"ADDZE",
	"ADDZECC",
	"ADDZEVCC",
	"ADDZEV",
	"ADDEX",
	"AND",
	"ANDCC",
	"ANDN",
	"ANDNCC",
	"ANDISCC",
	"BC",
	"BCL",
	"BEQ",
	"BGE",
	"BGT",
	"BLE",
	"BLT",
	"BNE",
	"BVC",
	"BVS",
	"CMP",
	"CMPU",
	"CMPEQB",
	"CNTLZW",
	"CNTLZWCC",
	"CRAND",
	"CRANDN",
	"CREQV",
	"CRNAND",
	"CRNOR",
	"CROR",
	"CRORN",
	"CRXOR",
	"DIVW",
	"DIVWCC",
	"DIVWVCC",
	"DIVWV",
	"DIVWU",
	"DIVWUCC",
	"DIVWUVCC",
	"DIVWUV",
	"EQV",
	"EQVCC",
	"EXTSB",
	"EXTSBCC",
	"EXTSH",
	"EXTSHCC",
	"FABS",
	"FABSCC",
	"FADD",
	"FADDCC",
	"FADDS",
	"FADDSCC",
	"FCMPO",
	"FCMPU",
	"FCTIW",
	"FCTIWCC",
	"FCTIWZ",
	"FCTIWZCC",
	"FDIV",
	"FDIVCC",
	"FDIVS",
	"FDIVSCC",
	"FMADD",
	"FMADDCC",
	"FMADDS",
	"FMADDSCC",
	"FMOVD",
	"FMOVDCC",
	"FMOVDU",
	"FMOVS",
	"FMOVSU",
	"FMOVSX",
	"FMOVSZ",
	"FMSUB",
	"FMSUBCC",
	"FMSUBS",
	"FMSUBSCC",
	"FMUL",
	"FMULCC",
	"FMULS",
	"FMULSCC",
	"FNABS",
	"FNABSCC",
	"FNEG",
	"FNEGCC",
	"FNMADD",
	"FNMADDCC",
	"FNMADDS",
	"FNMADDSCC",
	"FNMSUB",
	"FNMSUBCC",
	"FNMSUBS",
	"FNMSUBSCC",
	"FRSP",
	"FRSPCC",
	"FSUB",
	"FSUBCC",
	"FSUBS",
	"FSUBSCC",
	"ISEL",
	"MOVMW",
	"LBAR",
	"LHAR",
	"LSW",
	"LWAR",
	"LWSYNC",
	"MOVDBR",
	"MOVWBR",
	"MOVB",
	"MOVBU",
	"MOVBZ",
	"MOVBZU",
	"MOVH",
	"MOVHBR",
	"MOVHU",
	"MOVHZ",
	"MOVHZU",
	"MOVW",
	"MOVWU",
	"MOVFL",
	"MOVCRFS",
	"MTFSB0",
	"MTFSB0CC",
	"MTFSB1",
	"MTFSB1CC",
	"MULHW",
	"MULHWCC",
	"MULHWU",
	"MULHWUCC",
	"MULLW",
	"MULLWCC",
	"MULLWVCC",
	"MULLWV",
	"NAND",
	"NANDCC",
	"NEG",
	"NEGCC",
	"NEGVCC",
	"NEGV",
	"NOR",
	"NORCC",
	"OR",
	"ORCC",
	"ORN",
	"ORNCC",
	"ORIS",
	"REM",
	"REMCC",
	"REMV",
	"REMVCC",
	"REMU",
	"REMUCC",
	"REMUV",
	"REMUVCC",
	"RFI",
	"RLWMI",
	"RLWMICC",
	"RLWNM",
	"RLWNMCC",
	"SLW",
	"SLWCC",
	"SRW",
	"SRAW",
	"SRAWCC",
	"SRWCC",
	"STBCCC",
	"STSW",
	"STWCCC",
	"SUB",
	"SUBCC",
	"SUBVCC",
	"SUBC",
	"SUBCCC",
	"SUBCV",
	"SUBCVCC",
	"SUBME",
	"SUBMECC",
	"SUBMEVCC",
	"SUBMEV",
	"SUBV",
	"SUBE",
	"SUBECC",
	"SUBEV",
	"SUBEVCC",
	"SUBZE",
	"SUBZECC",
	"SUBZEVCC",
	"SUBZEV",
	"SYNC",
	"XOR",
	"XORCC",
	"XORIS",
	"DCBF",
	"DCBI",
	"DCBST",
	"DCBT",
	"DCBTST",
	"DCBZ",
	"ECIWX",
	"ECOWX",
	"EIEIO",
	"ICBI",
	"ISYNC",
	"PTESYNC",
	"TLBIE",
	"TLBIEL",
	"TLBSYNC",
	"TW",
	"SYSCALL",
	"WORD",
	"RFCI",
	"FCPSGN",
	"FCPSGNCC",
	"FRES",
	"FRESCC",
	"FRIM",
	"FRIMCC",
	"FRIP",
	"FRIPCC",
	"FRIZ",
	"FRIZCC",
	"FRIN",
	"FRINCC",
	"FRSQRTE",
	"FRSQRTECC",
	"FSEL",
	"FSELCC",
	"FSQRT",
	"FSQRTCC",
	"FSQRTS",
	"FSQRTSCC",
	"CNTLZD",
	"CNTLZDCC",
	"CMPW",
	"CMPWU",
	"CMPB",
	"FTDIV",
	"FTSQRT",
	"DIVD",
	"DIVDCC",
	"DIVDE",
	"DIVDECC",
	"DIVDEU",
	"DIVDEUCC",
	"DIVDVCC",
	"DIVDV",
	"DIVDU",
	"DIVDUCC",
	"DIVDUVCC",
	"DIVDUV",
	"EXTSW",
	"EXTSWCC",
	"FCFID",
	"FCFIDCC",
	"FCFIDU",
	"FCFIDUCC",
	"FCFIDS",
	"FCFIDSCC",
	"FCTID",
	"FCTIDCC",
	"FCTIDZ",
	"FCTIDZCC",
	"LDAR",
	"MOVD",
	"MOVDU",
	"MOVWZ",
	"MOVWZU",
	"MULHD",
	"MULHDCC",
	"MULHDU",
	"MULHDUCC",
	"MULLD",
	"MULLDCC",
	"MULLDVCC",
	"MULLDV",
	"RFID",
	"RLDMI",
	"RLDMICC",
	"RLDIMI",
	"RLDIMICC",
	"RLDC",
	"RLDCCC",
	"RLDCR",
	"RLDCRCC",
	"RLDICR",
	"RLDICRCC",
	"RLDCL",
	"RLDCLCC",
	"RLDICL",
	"RLDICLCC",
	"ROTL",
	"ROTLW",
	"SLBIA",
	"SLBIE",
	"SLBMFEE",
	"SLBMFEV",
	"SLBMTE",
	"SLD",
	"SLDCC",
	"SRD",
	"SRAD",
	"SRADCC",
	"SRDCC",
	"STDCCC",
	"TD",
	"DWORD",
	"REMD",
	"REMDCC",
	"REMDV",
	"REMDVCC",
	"REMDU",
	"REMDUCC",
	"REMDUV",
	"REMDUVCC",
	"HRFID",
	"POPCNTD",
	"POPCNTW",
	"POPCNTB",
	"COPY",
	"PASTECC",
	"DARN",
	"LDMX",
	"MADDHD",
	"MADDHDU",
	"MADDLD",
	"LV",
	"LVEBX",
	"LVEHX",
	"LVEWX",
	"LVX",
	"LVXL",
	"LVSL",
	"LVSR",
	"STV",
	"STVEBX",
	"STVEHX",
	"STVEWX",
	"STVX",
	"STVXL",
	"VAND",
	"VANDC",
	"VNAND",
	"VOR",
	"VORC",
	"VNOR",
	"VXOR",
	"VEQV",
	"VADDUM",
	"VADDUBM",
	"VADDUHM",
	"VADDUWM",
	"VADDUDM",
	"VADDUQM",
	"VADDCU",
	"VADDCUQ",
	"VADDCUW",
	"VADDUS",
	"VADDUBS",
	"VADDUHS",
	"VADDUWS",
	"VADDSS",
	"VADDSBS",
	"VADDSHS",
	"VADDSWS",
	"VADDE",
	"VADDEUQM",
	"VADDECUQ",
	"VSUBUM",
	"VSUBUBM",
	"VSUBUHM",
	"VSUBUWM",
	"VSUBUDM",
	"VSUBUQM",
	"VSUBCU",
	"VSUBCUQ",
	"VSUBCUW",
	"VSUBUS",
	"VSUBUBS",
	"VSUBUHS",
	"VSUBUWS",
	"VSUBSS",
	"VSUBSBS",
	"VSUBSHS",
	"VSUBSWS",
	"VSUBE",
	"VSUBEUQM",
	"VSUBECUQ",
	"VMULESB",
	"VMULOSB",
	"VMULEUB",
	"VMULOUB",
	"VMULESH",
	"VMULOSH",
	"VMULEUH",
	"VMULOUH",
	"VMULESW",
	"VMULOSW",
	"VMULEUW",
	"VMULOUW",
	"VMULUWM",
	"VPMSUM",
	"VPMSUMB",
	"VPMSUMH",
	"VPMSUMW",
	"VPMSUMD",
	"VMSUMUDM",
	"VR",
	"VRLB",
	"VRLH",
	"VRLW",
	"VRLD",
	"VS",
	"VSLB",
	"VSLH",
	"VSLW",
	"VSL",
	"VSLO",
	"VSRB",
	"VSRH",
	"VSRW",
	"VSR",
	"VSRO",
	"VSLD",
	"VSRD",
	"VSA",
	"VSRAB",
	"VSRAH",
	"VSRAW",
	"VSRAD",
	"VSOI",
	"VSLDOI",
	"VCLZ",
	"VCLZB",
	"VCLZH",
	"VCLZW",
	"VCLZD",
	"VPOPCNT",
	"VPOPCNTB",
	"VPOPCNTH",
	"VPOPCNTW",
	"VPOPCNTD",
	"VCMPEQ",
	"VCMPEQUB",
	"VCMPEQUBCC",
	"VCMPEQUH",
	"VCMPEQUHCC",
	"VCMPEQUW",
	"VCMPEQUWCC",
	"VCMPEQUD",
	"VCMPEQUDCC",
	"VCMPGT",
	"VCMPGTUB",
	"VCMPGTUBCC",
	"VCMPGTUH",
	"VCMPGTUHCC",
	"VCMPGTUW",
	"VCMPGTUWCC",
	"VCMPGTUD",
	"VCMPGTUDCC",
	"VCMPGTSB",
	"VCMPGTSBCC",
	"VCMPGTSH",
	"VCMPGTSHCC",
	"VCMPGTSW",
	"VCMPGTSWCC",
	"VCMPGTSD",
	"VCMPGTSDCC",
	"VCMPNEZB",
	"VCMPNEZBCC",
	"VPERM",
	"VBPERMQ",
	"VBPERMD",
	"VSEL",
	"VSPLT",
	"VSPLTB",
	"VSPLTH",
	"VSPLTW",
	"VSPLTI",
	"VSPLTISB",
	"VSPLTISH",
	"VSPLTISW",
	"VCIPH",
	"VCIPHER",
	"VCIPHERLAST",
	"VNCIPH",
	"VNCIPHER",
	"VNCIPHERLAST",
	"VSBOX",
	"VSHASIGMA",
	"VSHASIGMAW",
	"VSHASIGMAD",
	"LXV",
	"LXVD2X",
	"LXVDSX",
	"LXVW4X",
	"STXV",
	"STXVD2X",
	"STXVW4X",
	"LXS",
	"LXSDX",
	"STXS",
	"STXSDX",
	"LXSI",
	"LXSIWAX",
	"LXSIWZX",
	"STXSI",
	"STXSIWX",
	"MFVSR",
	"MFVSRD",
	"MFFPRD",
	"MFVRD",
	"MFVSRWZ",
	"MFVSRLD",
	"MTVSR",
	"MTVSRD",
	"MTFPRD",
	"MTVRD",
	"MTVSRWA",
	"MTVSRWZ",
	"MTVSRDD",
	"MTVSRWS",
	"XXLAND",
	"XXLANDQ",
	"XXLANDC",
	"XXLEQV",
	"XXLNAND",
	"XXLOR",
	"XXLORC",
	"XXLNOR",
	"XXLORQ",
	"XXLXOR",
	"XXSEL",
	"XXMRG",
	"XXMRGHW",
	"XXMRGLW",
	"XXSPLT",
	"XXSPLTW",
	"XXPERM",
	"XXPERMDI",
	"XXSI",
	"XXSLDWI",
	"XSCV",
	"XSCVDPSP",
	"XSCVSPDP",
	"XSCVDPSPN",
	"XSCVSPDPN",
	"XVCV",
	"XVCVDPSP",
	"XVCVSPDP",
	"XSCVX",
	"XSCVDPSXDS",
	"XSCVDPSXWS",
	"XSCVDPUXDS",
	"XSCVDPUXWS",
	"XSCVXP",
	"XSCVSXDDP",
	"XSCVUXDDP",
	"XSCVSXDSP",
	"XSCVUXDSP",
	"XVCVX",
	"XVCVDPSXDS",
	"XVCVDPSXWS",
	"XVCVDPUXDS",
	"XVCVDPUXWS",
	"XVCVSPSXDS",
	"XVCVSPSXWS",
	"XVCVSPUXDS",
	"XVCVSPUXWS",
	"XVCVXP",
	"XVCVSXDDP",
	"XVCVSXWDP",
	"XVCVUXDDP",
	"XVCVUXWDP",
	"XVCVSXDSP",
	"XVCVSXWSP",
	"XVCVUXDSP",
	"XVCVUXWSP",
	"LAST",
}
