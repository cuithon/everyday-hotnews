// Code generated by "stringer -type=RelocType"; DO NOT EDIT.

package objabi

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[R_ADDR-1]
	_ = x[R_ADDRPOWER-2]
	_ = x[R_ADDRARM64-3]
	_ = x[R_ADDRMIPS-4]
	_ = x[R_ADDROFF-5]
	_ = x[R_SIZE-6]
	_ = x[R_CALL-7]
	_ = x[R_CALLARM-8]
	_ = x[R_CALLARM64-9]
	_ = x[R_CALLIND-10]
	_ = x[R_CALLPOWER-11]
	_ = x[R_CALLMIPS-12]
	_ = x[R_CONST-13]
	_ = x[R_PCREL-14]
	_ = x[R_TLS_LE-15]
	_ = x[R_TLS_IE-16]
	_ = x[R_GOTOFF-17]
	_ = x[R_PLT0-18]
	_ = x[R_PLT1-19]
	_ = x[R_PLT2-20]
	_ = x[R_USEFIELD-21]
	_ = x[R_USETYPE-22]
	_ = x[R_USEIFACE-23]
	_ = x[R_USEIFACEMETHOD-24]
	_ = x[R_USENAMEDMETHOD-25]
	_ = x[R_METHODOFF-26]
	_ = x[R_KEEP-27]
	_ = x[R_POWER_TOC-28]
	_ = x[R_GOTPCREL-29]
	_ = x[R_JMPMIPS-30]
	_ = x[R_DWARFSECREF-31]
	_ = x[R_DWARFFILEREF-32]
	_ = x[R_ARM64_TLS_LE-33]
	_ = x[R_ARM64_TLS_IE-34]
	_ = x[R_ARM64_GOTPCREL-35]
	_ = x[R_ARM64_GOT-36]
	_ = x[R_ARM64_PCREL-37]
	_ = x[R_ARM64_PCREL_LDST8-38]
	_ = x[R_ARM64_PCREL_LDST16-39]
	_ = x[R_ARM64_PCREL_LDST32-40]
	_ = x[R_ARM64_PCREL_LDST64-41]
	_ = x[R_ARM64_LDST8-42]
	_ = x[R_ARM64_LDST16-43]
	_ = x[R_ARM64_LDST32-44]
	_ = x[R_ARM64_LDST64-45]
	_ = x[R_ARM64_LDST128-46]
	_ = x[R_POWER_TLS_LE-47]
	_ = x[R_POWER_TLS_IE-48]
	_ = x[R_POWER_TLS-49]
	_ = x[R_POWER_TLS_IE_PCREL34-50]
	_ = x[R_POWER_TLS_LE_TPREL34-51]
	_ = x[R_ADDRPOWER_DS-52]
	_ = x[R_ADDRPOWER_GOT-53]
	_ = x[R_ADDRPOWER_GOT_PCREL34-54]
	_ = x[R_ADDRPOWER_PCREL-55]
	_ = x[R_ADDRPOWER_TOCREL-56]
	_ = x[R_ADDRPOWER_TOCREL_DS-57]
	_ = x[R_ADDRPOWER_D34-58]
	_ = x[R_ADDRPOWER_PCREL34-59]
	_ = x[R_RISCV_JAL-60]
	_ = x[R_RISCV_JAL_TRAMP-61]
	_ = x[R_RISCV_CALL-62]
	_ = x[R_RISCV_PCREL_ITYPE-63]
	_ = x[R_RISCV_PCREL_STYPE-64]
	_ = x[R_RISCV_TLS_IE-65]
	_ = x[R_RISCV_TLS_LE-66]
	_ = x[R_RISCV_GOT_HI20-67]
	_ = x[R_RISCV_PCREL_HI20-68]
	_ = x[R_RISCV_PCREL_LO12_I-69]
	_ = x[R_RISCV_PCREL_LO12_S-70]
	_ = x[R_RISCV_BRANCH-71]
	_ = x[R_RISCV_RVC_BRANCH-72]
	_ = x[R_RISCV_RVC_JUMP-73]
	_ = x[R_PCRELDBL-74]
	_ = x[R_ADDRLOONG64-75]
	_ = x[R_ADDRLOONG64U-76]
	_ = x[R_ADDRLOONG64TLS-77]
	_ = x[R_ADDRLOONG64TLSU-78]
	_ = x[R_CALLLOONG64-79]
	_ = x[R_LOONG64_TLS_IE_PCREL_HI-80]
	_ = x[R_LOONG64_TLS_IE_LO-81]
	_ = x[R_JMPLOONG64-82]
	_ = x[R_ADDRMIPSU-83]
	_ = x[R_ADDRMIPSTLS-84]
	_ = x[R_ADDRCUOFF-85]
	_ = x[R_WASMIMPORT-86]
	_ = x[R_XCOFFREF-87]
	_ = x[R_PEIMAGEOFF-88]
	_ = x[R_INITORDER-89]
}

const _RelocType_name = "R_ADDRR_ADDRPOWERR_ADDRARM64R_ADDRMIPSR_ADDROFFR_SIZER_CALLR_CALLARMR_CALLARM64R_CALLINDR_CALLPOWERR_CALLMIPSR_CONSTR_PCRELR_TLS_LER_TLS_IER_GOTOFFR_PLT0R_PLT1R_PLT2R_USEFIELDR_USETYPER_USEIFACER_USEIFACEMETHODR_USENAMEDMETHODR_METHODOFFR_KEEPR_POWER_TOCR_GOTPCRELR_JMPMIPSR_DWARFSECREFR_DWARFFILEREFR_ARM64_TLS_LER_ARM64_TLS_IER_ARM64_GOTPCRELR_ARM64_GOTR_ARM64_PCRELR_ARM64_PCREL_LDST8R_ARM64_PCREL_LDST16R_ARM64_PCREL_LDST32R_ARM64_PCREL_LDST64R_ARM64_LDST8R_ARM64_LDST16R_ARM64_LDST32R_ARM64_LDST64R_ARM64_LDST128R_POWER_TLS_LER_POWER_TLS_IER_POWER_TLSR_POWER_TLS_IE_PCREL34R_POWER_TLS_LE_TPREL34R_ADDRPOWER_DSR_ADDRPOWER_GOTR_ADDRPOWER_GOT_PCREL34R_ADDRPOWER_PCRELR_ADDRPOWER_TOCRELR_ADDRPOWER_TOCREL_DSR_ADDRPOWER_D34R_ADDRPOWER_PCREL34R_RISCV_JALR_RISCV_JAL_TRAMPR_RISCV_CALLR_RISCV_PCREL_ITYPER_RISCV_PCREL_STYPER_RISCV_TLS_IER_RISCV_TLS_LER_RISCV_GOT_HI20R_RISCV_PCREL_HI20R_RISCV_PCREL_LO12_IR_RISCV_PCREL_LO12_SR_RISCV_BRANCHR_RISCV_RVC_BRANCHR_RISCV_RVC_JUMPR_PCRELDBLR_ADDRLOONG64R_ADDRLOONG64UR_ADDRLOONG64TLSR_ADDRLOONG64TLSUR_CALLLOONG64R_LOONG64_TLS_IE_PCREL_HIR_LOONG64_TLS_IE_LOR_JMPLOONG64R_ADDRMIPSUR_ADDRMIPSTLSR_ADDRCUOFFR_WASMIMPORTR_XCOFFREFR_PEIMAGEOFFR_INITORDER"

var _RelocType_index = [...]uint16{0, 6, 17, 28, 38, 47, 53, 59, 68, 79, 88, 99, 109, 116, 123, 131, 139, 147, 153, 159, 165, 175, 184, 194, 210, 226, 237, 243, 254, 264, 273, 286, 300, 314, 328, 344, 355, 368, 387, 407, 427, 447, 460, 474, 488, 502, 517, 531, 545, 556, 578, 600, 614, 629, 652, 669, 687, 708, 723, 742, 753, 770, 782, 801, 820, 834, 848, 864, 882, 902, 922, 936, 954, 970, 980, 993, 1007, 1023, 1040, 1053, 1078, 1097, 1109, 1120, 1133, 1144, 1156, 1166, 1178, 1189}

func (i RelocType) String() string {
	i -= 1
	if i < 0 || i >= RelocType(len(_RelocType_index)-1) {
		return "RelocType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _RelocType_name[_RelocType_index[i]:_RelocType_index[i+1]]
}
