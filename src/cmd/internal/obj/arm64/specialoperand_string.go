// Code generated by "stringer -type SpecialOperand -trimprefix SPOP_"; DO NOT EDIT.

package arm64

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SPOP_PLDL1KEEP-0]
	_ = x[SPOP_BEGIN-0]
	_ = x[SPOP_PLDL1STRM-1]
	_ = x[SPOP_PLDL2KEEP-2]
	_ = x[SPOP_PLDL2STRM-3]
	_ = x[SPOP_PLDL3KEEP-4]
	_ = x[SPOP_PLDL3STRM-5]
	_ = x[SPOP_PLIL1KEEP-6]
	_ = x[SPOP_PLIL1STRM-7]
	_ = x[SPOP_PLIL2KEEP-8]
	_ = x[SPOP_PLIL2STRM-9]
	_ = x[SPOP_PLIL3KEEP-10]
	_ = x[SPOP_PLIL3STRM-11]
	_ = x[SPOP_PSTL1KEEP-12]
	_ = x[SPOP_PSTL1STRM-13]
	_ = x[SPOP_PSTL2KEEP-14]
	_ = x[SPOP_PSTL2STRM-15]
	_ = x[SPOP_PSTL3KEEP-16]
	_ = x[SPOP_PSTL3STRM-17]
	_ = x[SPOP_VMALLE1IS-18]
	_ = x[SPOP_VAE1IS-19]
	_ = x[SPOP_ASIDE1IS-20]
	_ = x[SPOP_VAAE1IS-21]
	_ = x[SPOP_VALE1IS-22]
	_ = x[SPOP_VAALE1IS-23]
	_ = x[SPOP_VMALLE1-24]
	_ = x[SPOP_VAE1-25]
	_ = x[SPOP_ASIDE1-26]
	_ = x[SPOP_VAAE1-27]
	_ = x[SPOP_VALE1-28]
	_ = x[SPOP_VAALE1-29]
	_ = x[SPOP_IPAS2E1IS-30]
	_ = x[SPOP_IPAS2LE1IS-31]
	_ = x[SPOP_ALLE2IS-32]
	_ = x[SPOP_VAE2IS-33]
	_ = x[SPOP_ALLE1IS-34]
	_ = x[SPOP_VALE2IS-35]
	_ = x[SPOP_VMALLS12E1IS-36]
	_ = x[SPOP_IPAS2E1-37]
	_ = x[SPOP_IPAS2LE1-38]
	_ = x[SPOP_ALLE2-39]
	_ = x[SPOP_VAE2-40]
	_ = x[SPOP_ALLE1-41]
	_ = x[SPOP_VALE2-42]
	_ = x[SPOP_VMALLS12E1-43]
	_ = x[SPOP_ALLE3IS-44]
	_ = x[SPOP_VAE3IS-45]
	_ = x[SPOP_VALE3IS-46]
	_ = x[SPOP_ALLE3-47]
	_ = x[SPOP_VAE3-48]
	_ = x[SPOP_VALE3-49]
	_ = x[SPOP_VMALLE1OS-50]
	_ = x[SPOP_VAE1OS-51]
	_ = x[SPOP_ASIDE1OS-52]
	_ = x[SPOP_VAAE1OS-53]
	_ = x[SPOP_VALE1OS-54]
	_ = x[SPOP_VAALE1OS-55]
	_ = x[SPOP_RVAE1IS-56]
	_ = x[SPOP_RVAAE1IS-57]
	_ = x[SPOP_RVALE1IS-58]
	_ = x[SPOP_RVAALE1IS-59]
	_ = x[SPOP_RVAE1OS-60]
	_ = x[SPOP_RVAAE1OS-61]
	_ = x[SPOP_RVALE1OS-62]
	_ = x[SPOP_RVAALE1OS-63]
	_ = x[SPOP_RVAE1-64]
	_ = x[SPOP_RVAAE1-65]
	_ = x[SPOP_RVALE1-66]
	_ = x[SPOP_RVAALE1-67]
	_ = x[SPOP_RIPAS2E1IS-68]
	_ = x[SPOP_RIPAS2LE1IS-69]
	_ = x[SPOP_ALLE2OS-70]
	_ = x[SPOP_VAE2OS-71]
	_ = x[SPOP_ALLE1OS-72]
	_ = x[SPOP_VALE2OS-73]
	_ = x[SPOP_VMALLS12E1OS-74]
	_ = x[SPOP_RVAE2IS-75]
	_ = x[SPOP_RVALE2IS-76]
	_ = x[SPOP_IPAS2E1OS-77]
	_ = x[SPOP_RIPAS2E1-78]
	_ = x[SPOP_RIPAS2E1OS-79]
	_ = x[SPOP_IPAS2LE1OS-80]
	_ = x[SPOP_RIPAS2LE1-81]
	_ = x[SPOP_RIPAS2LE1OS-82]
	_ = x[SPOP_RVAE2OS-83]
	_ = x[SPOP_RVALE2OS-84]
	_ = x[SPOP_RVAE2-85]
	_ = x[SPOP_RVALE2-86]
	_ = x[SPOP_ALLE3OS-87]
	_ = x[SPOP_VAE3OS-88]
	_ = x[SPOP_VALE3OS-89]
	_ = x[SPOP_RVAE3IS-90]
	_ = x[SPOP_RVALE3IS-91]
	_ = x[SPOP_RVAE3OS-92]
	_ = x[SPOP_RVALE3OS-93]
	_ = x[SPOP_RVAE3-94]
	_ = x[SPOP_RVALE3-95]
	_ = x[SPOP_DAIFSet-96]
	_ = x[SPOP_DAIFClr-97]
	_ = x[SPOP_EQ-98]
	_ = x[SPOP_NE-99]
	_ = x[SPOP_HS-100]
	_ = x[SPOP_LO-101]
	_ = x[SPOP_MI-102]
	_ = x[SPOP_PL-103]
	_ = x[SPOP_VS-104]
	_ = x[SPOP_VC-105]
	_ = x[SPOP_HI-106]
	_ = x[SPOP_LS-107]
	_ = x[SPOP_GE-108]
	_ = x[SPOP_LT-109]
	_ = x[SPOP_GT-110]
	_ = x[SPOP_LE-111]
	_ = x[SPOP_AL-112]
	_ = x[SPOP_NV-113]
	_ = x[SPOP_END-114]
}

const _SpecialOperand_name = "PLDL1KEEPPLDL1STRMPLDL2KEEPPLDL2STRMPLDL3KEEPPLDL3STRMPLIL1KEEPPLIL1STRMPLIL2KEEPPLIL2STRMPLIL3KEEPPLIL3STRMPSTL1KEEPPSTL1STRMPSTL2KEEPPSTL2STRMPSTL3KEEPPSTL3STRMVMALLE1ISVAE1ISASIDE1ISVAAE1ISVALE1ISVAALE1ISVMALLE1VAE1ASIDE1VAAE1VALE1VAALE1IPAS2E1ISIPAS2LE1ISALLE2ISVAE2ISALLE1ISVALE2ISVMALLS12E1ISIPAS2E1IPAS2LE1ALLE2VAE2ALLE1VALE2VMALLS12E1ALLE3ISVAE3ISVALE3ISALLE3VAE3VALE3VMALLE1OSVAE1OSASIDE1OSVAAE1OSVALE1OSVAALE1OSRVAE1ISRVAAE1ISRVALE1ISRVAALE1ISRVAE1OSRVAAE1OSRVALE1OSRVAALE1OSRVAE1RVAAE1RVALE1RVAALE1RIPAS2E1ISRIPAS2LE1ISALLE2OSVAE2OSALLE1OSVALE2OSVMALLS12E1OSRVAE2ISRVALE2ISIPAS2E1OSRIPAS2E1RIPAS2E1OSIPAS2LE1OSRIPAS2LE1RIPAS2LE1OSRVAE2OSRVALE2OSRVAE2RVALE2ALLE3OSVAE3OSVALE3OSRVAE3ISRVALE3ISRVAE3OSRVALE3OSRVAE3RVALE3DAIFSetDAIFClrEQNEHSLOMIPLVSVCHILSGELTGTLEALNVEND"

var _SpecialOperand_index = [...]uint16{0, 9, 18, 27, 36, 45, 54, 63, 72, 81, 90, 99, 108, 117, 126, 135, 144, 153, 162, 171, 177, 185, 192, 199, 207, 214, 218, 224, 229, 234, 240, 249, 259, 266, 272, 279, 286, 298, 305, 313, 318, 322, 327, 332, 342, 349, 355, 362, 367, 371, 376, 385, 391, 399, 406, 413, 421, 428, 436, 444, 453, 460, 468, 476, 485, 490, 496, 502, 509, 519, 530, 537, 543, 550, 557, 569, 576, 584, 593, 601, 611, 621, 630, 641, 648, 656, 661, 667, 674, 680, 687, 694, 702, 709, 717, 722, 728, 735, 742, 744, 746, 748, 750, 752, 754, 756, 758, 760, 762, 764, 766, 768, 770, 772, 774, 777}

func (i SpecialOperand) String() string {
	if i < 0 || i >= SpecialOperand(len(_SpecialOperand_index)-1) {
		return "SpecialOperand(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SpecialOperand_name[_SpecialOperand_index[i]:_SpecialOperand_index[i+1]]
}
