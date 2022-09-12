// Code generated by "stringer -type=Op -trimprefix=O node.go"; DO NOT EDIT.

package ir

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OXXX-0]
	_ = x[ONAME-1]
	_ = x[ONONAME-2]
	_ = x[OTYPE-3]
	_ = x[OLITERAL-4]
	_ = x[ONIL-5]
	_ = x[OADD-6]
	_ = x[OSUB-7]
	_ = x[OOR-8]
	_ = x[OXOR-9]
	_ = x[OADDSTR-10]
	_ = x[OADDR-11]
	_ = x[OANDAND-12]
	_ = x[OAPPEND-13]
	_ = x[OBYTES2STR-14]
	_ = x[OBYTES2STRTMP-15]
	_ = x[ORUNES2STR-16]
	_ = x[OSTR2BYTES-17]
	_ = x[OSTR2BYTESTMP-18]
	_ = x[OSTR2RUNES-19]
	_ = x[OSLICE2ARR-20]
	_ = x[OSLICE2ARRPTR-21]
	_ = x[OAS-22]
	_ = x[OAS2-23]
	_ = x[OAS2DOTTYPE-24]
	_ = x[OAS2FUNC-25]
	_ = x[OAS2MAPR-26]
	_ = x[OAS2RECV-27]
	_ = x[OASOP-28]
	_ = x[OCALL-29]
	_ = x[OCALLFUNC-30]
	_ = x[OCALLMETH-31]
	_ = x[OCALLINTER-32]
	_ = x[OCAP-33]
	_ = x[OCLOSE-34]
	_ = x[OCLOSURE-35]
	_ = x[OCOMPLIT-36]
	_ = x[OMAPLIT-37]
	_ = x[OSTRUCTLIT-38]
	_ = x[OARRAYLIT-39]
	_ = x[OSLICELIT-40]
	_ = x[OPTRLIT-41]
	_ = x[OCONV-42]
	_ = x[OCONVIFACE-43]
	_ = x[OCONVIDATA-44]
	_ = x[OCONVNOP-45]
	_ = x[OCOPY-46]
	_ = x[ODCL-47]
	_ = x[ODCLFUNC-48]
	_ = x[ODCLCONST-49]
	_ = x[ODCLTYPE-50]
	_ = x[ODELETE-51]
	_ = x[ODOT-52]
	_ = x[ODOTPTR-53]
	_ = x[ODOTMETH-54]
	_ = x[ODOTINTER-55]
	_ = x[OXDOT-56]
	_ = x[ODOTTYPE-57]
	_ = x[ODOTTYPE2-58]
	_ = x[OEQ-59]
	_ = x[ONE-60]
	_ = x[OLT-61]
	_ = x[OLE-62]
	_ = x[OGE-63]
	_ = x[OGT-64]
	_ = x[ODEREF-65]
	_ = x[OINDEX-66]
	_ = x[OINDEXMAP-67]
	_ = x[OKEY-68]
	_ = x[OSTRUCTKEY-69]
	_ = x[OLEN-70]
	_ = x[OMAKE-71]
	_ = x[OMAKECHAN-72]
	_ = x[OMAKEMAP-73]
	_ = x[OMAKESLICE-74]
	_ = x[OMAKESLICECOPY-75]
	_ = x[OMUL-76]
	_ = x[ODIV-77]
	_ = x[OMOD-78]
	_ = x[OLSH-79]
	_ = x[ORSH-80]
	_ = x[OAND-81]
	_ = x[OANDNOT-82]
	_ = x[ONEW-83]
	_ = x[ONOT-84]
	_ = x[OBITNOT-85]
	_ = x[OPLUS-86]
	_ = x[ONEG-87]
	_ = x[OOROR-88]
	_ = x[OPANIC-89]
	_ = x[OPRINT-90]
	_ = x[OPRINTN-91]
	_ = x[OPAREN-92]
	_ = x[OSEND-93]
	_ = x[OSLICE-94]
	_ = x[OSLICEARR-95]
	_ = x[OSLICESTR-96]
	_ = x[OSLICE3-97]
	_ = x[OSLICE3ARR-98]
	_ = x[OSLICEHEADER-99]
	_ = x[OSTRINGHEADER-100]
	_ = x[ORECOVER-101]
	_ = x[ORECOVERFP-102]
	_ = x[ORECV-103]
	_ = x[ORUNESTR-104]
	_ = x[OSELRECV2-105]
	_ = x[OREAL-106]
	_ = x[OIMAG-107]
	_ = x[OCOMPLEX-108]
	_ = x[OALIGNOF-109]
	_ = x[OOFFSETOF-110]
	_ = x[OSIZEOF-111]
	_ = x[OUNSAFEADD-112]
	_ = x[OUNSAFESLICE-113]
	_ = x[OUNSAFESLICEDATA-114]
	_ = x[OUNSAFESTRING-115]
	_ = x[OUNSAFESTRINGDATA-116]
	_ = x[OMETHEXPR-117]
	_ = x[OMETHVALUE-118]
	_ = x[OBLOCK-119]
	_ = x[OBREAK-120]
	_ = x[OCASE-121]
	_ = x[OCONTINUE-122]
	_ = x[ODEFER-123]
	_ = x[OFALL-124]
	_ = x[OFOR-125]
	_ = x[OGOTO-126]
	_ = x[OIF-127]
	_ = x[OLABEL-128]
	_ = x[OGO-129]
	_ = x[ORANGE-130]
	_ = x[ORETURN-131]
	_ = x[OSELECT-132]
	_ = x[OSWITCH-133]
	_ = x[OTYPESW-134]
	_ = x[OFUNCINST-135]
	_ = x[OINLCALL-136]
	_ = x[OEFACE-137]
	_ = x[OITAB-138]
	_ = x[OIDATA-139]
	_ = x[OSPTR-140]
	_ = x[OCFUNC-141]
	_ = x[OCHECKNIL-142]
	_ = x[ORESULT-143]
	_ = x[OINLMARK-144]
	_ = x[OLINKSYMOFFSET-145]
	_ = x[OJUMPTABLE-146]
	_ = x[ODYNAMICDOTTYPE-147]
	_ = x[ODYNAMICDOTTYPE2-148]
	_ = x[ODYNAMICTYPE-149]
	_ = x[OTAILCALL-150]
	_ = x[OGETG-151]
	_ = x[OGETCALLERPC-152]
	_ = x[OGETCALLERSP-153]
	_ = x[OEND-154]
}

const _Op_name = "XXXNAMENONAMETYPELITERALNILADDSUBORXORADDSTRADDRANDANDAPPENDBYTES2STRBYTES2STRTMPRUNES2STRSTR2BYTESSTR2BYTESTMPSTR2RUNESSLICE2ARRSLICE2ARRPTRASAS2AS2DOTTYPEAS2FUNCAS2MAPRAS2RECVASOPCALLCALLFUNCCALLMETHCALLINTERCAPCLOSECLOSURECOMPLITMAPLITSTRUCTLITARRAYLITSLICELITPTRLITCONVCONVIFACECONVIDATACONVNOPCOPYDCLDCLFUNCDCLCONSTDCLTYPEDELETEDOTDOTPTRDOTMETHDOTINTERXDOTDOTTYPEDOTTYPE2EQNELTLEGEGTDEREFINDEXINDEXMAPKEYSTRUCTKEYLENMAKEMAKECHANMAKEMAPMAKESLICEMAKESLICECOPYMULDIVMODLSHRSHANDANDNOTNEWNOTBITNOTPLUSNEGORORPANICPRINTPRINTNPARENSENDSLICESLICEARRSLICESTRSLICE3SLICE3ARRSLICEHEADERSTRINGHEADERRECOVERRECOVERFPRECVRUNESTRSELRECV2REALIMAGCOMPLEXALIGNOFOFFSETOFSIZEOFUNSAFEADDUNSAFESLICEUNSAFESLICEDATAUNSAFESTRINGUNSAFESTRINGDATAMETHEXPRMETHVALUEBLOCKBREAKCASECONTINUEDEFERFALLFORGOTOIFLABELGORANGERETURNSELECTSWITCHTYPESWFUNCINSTINLCALLEFACEITABIDATASPTRCFUNCCHECKNILRESULTINLMARKLINKSYMOFFSETJUMPTABLEDYNAMICDOTTYPEDYNAMICDOTTYPE2DYNAMICTYPETAILCALLGETGGETCALLERPCGETCALLERSPEND"

var _Op_index = [...]uint16{0, 3, 7, 13, 17, 24, 27, 30, 33, 35, 38, 44, 48, 54, 60, 69, 81, 90, 99, 111, 120, 129, 141, 143, 146, 156, 163, 170, 177, 181, 185, 193, 201, 210, 213, 218, 225, 232, 238, 247, 255, 263, 269, 273, 282, 291, 298, 302, 305, 312, 320, 327, 333, 336, 342, 349, 357, 361, 368, 376, 378, 380, 382, 384, 386, 388, 393, 398, 406, 409, 418, 421, 425, 433, 440, 449, 462, 465, 468, 471, 474, 477, 480, 486, 489, 492, 498, 502, 505, 509, 514, 519, 525, 530, 534, 539, 547, 555, 561, 570, 581, 593, 600, 609, 613, 620, 628, 632, 636, 643, 650, 658, 664, 673, 684, 699, 711, 727, 735, 744, 749, 754, 758, 766, 771, 775, 778, 782, 784, 789, 791, 796, 802, 808, 814, 820, 828, 835, 840, 844, 849, 853, 858, 866, 872, 879, 892, 901, 915, 930, 941, 949, 953, 964, 975, 978}

func (i Op) String() string {
	if i >= Op(len(_Op_index)-1) {
		return "Op(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Op_name[_Op_index[i]:_Op_index[i+1]]
}
