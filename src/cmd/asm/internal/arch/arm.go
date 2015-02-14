// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file encapsulates some of the odd characteristics of the ARM
// instruction set, to minimize its interaction with the core of the
// assembler.

package arch

import (
	"strings"

	"cmd/internal/obj"
	"cmd/internal/obj/arm"
)

var armLS = map[string]uint8{
	"U":  arm.C_UBIT,
	"S":  arm.C_SBIT,
	"W":  arm.C_WBIT,
	"P":  arm.C_PBIT,
	"PW": arm.C_WBIT | arm.C_PBIT,
	"WP": arm.C_WBIT | arm.C_PBIT,
}

var armSCOND = map[string]uint8{
	"EQ":  arm.C_SCOND_EQ,
	"NE":  arm.C_SCOND_NE,
	"CS":  arm.C_SCOND_HS,
	"HS":  arm.C_SCOND_HS,
	"CC":  arm.C_SCOND_LO,
	"LO":  arm.C_SCOND_LO,
	"MI":  arm.C_SCOND_MI,
	"PL":  arm.C_SCOND_PL,
	"VS":  arm.C_SCOND_VS,
	"VC":  arm.C_SCOND_VC,
	"HI":  arm.C_SCOND_HI,
	"LS":  arm.C_SCOND_LS,
	"GE":  arm.C_SCOND_GE,
	"LT":  arm.C_SCOND_LT,
	"GT":  arm.C_SCOND_GT,
	"LE":  arm.C_SCOND_LE,
	"AL":  arm.C_SCOND_NONE,
	"U":   arm.C_UBIT,
	"S":   arm.C_SBIT,
	"W":   arm.C_WBIT,
	"P":   arm.C_PBIT,
	"PW":  arm.C_WBIT | arm.C_PBIT,
	"WP":  arm.C_WBIT | arm.C_PBIT,
	"F":   arm.C_FBIT,
	"IBW": arm.C_WBIT | arm.C_PBIT | arm.C_UBIT,
	"IAW": arm.C_WBIT | arm.C_UBIT,
	"DBW": arm.C_WBIT | arm.C_PBIT,
	"DAW": arm.C_WBIT,
	"IB":  arm.C_PBIT | arm.C_UBIT,
	"IA":  arm.C_UBIT,
	"DB":  arm.C_PBIT,
	"DA":  0,
}

var armJump = map[string]bool{
	"B":    true,
	"BL":   true,
	"BEQ":  true,
	"BNE":  true,
	"BCS":  true,
	"BHS":  true,
	"BCC":  true,
	"BLO":  true,
	"BMI":  true,
	"BPL":  true,
	"BVS":  true,
	"BVC":  true,
	"BHI":  true,
	"BLS":  true,
	"BGE":  true,
	"BLT":  true,
	"BGT":  true,
	"BLE":  true,
	"CALL": true,
}

func jumpArm(word string) bool {
	return armJump[word]
}

// IsARMCMP reports whether the op (as defined by an arm.A* constant) is
// one of the comparison instructions that require special handling.
func IsARMCMP(op int) bool {
	switch op {
	case arm.ACMN, arm.ACMP, arm.ATEQ, arm.ATST:
		return true
	}
	return false
}

// IsARMSTREX reports whether the op (as defined by an arm.A* constant) is
// one of the STREX-like instructions that require special handling.
func IsARMSTREX(op int) bool {
	switch op {
	case arm.ASTREX, arm.ASTREXD, arm.ASWPW, arm.ASWPBU:
		return true
	}
	return false
}

// IsARMMRC reports whether the op (as defined by an arm.A* constant) is
// MRC or MCR
func IsARMMRC(op int) bool {
	switch op {
	case arm.AMRC /*, arm.AMCR*/ :
		return true
	}
	return false
}

// IsARMMULA reports whether the op (as defined by an arm.A* constant) is
// MULA, MULAWT or MULAWB, the 4-operand instructions.
func IsARMMULA(op int) bool {
	switch op {
	case arm.AMULA, arm.AMULAWB, arm.AMULAWT:
		return true
	}
	return false
}

var bcode = []int{
	arm.ABEQ,
	arm.ABNE,
	arm.ABCS,
	arm.ABCC,
	arm.ABMI,
	arm.ABPL,
	arm.ABVS,
	arm.ABVC,
	arm.ABHI,
	arm.ABLS,
	arm.ABGE,
	arm.ABLT,
	arm.ABGT,
	arm.ABLE,
	arm.AB,
	obj.ANOP,
}

// ARMConditionCodes handles the special condition code situation for the ARM.
// It returns a boolean to indicate success; failure means cond was unrecognized.
func ARMConditionCodes(prog *obj.Prog, cond string) bool {
	if cond == "" {
		return true
	}
	bits, ok := parseARMCondition(cond)
	if !ok {
		return false
	}
	/* hack to make B.NE etc. work: turn it into the corresponding conditional */
	if prog.As == arm.AB {
		prog.As = int16(bcode[(bits^arm.C_SCOND_XOR)&0xf])
		bits = (bits &^ 0xf) | arm.C_SCOND_NONE
	}
	prog.Scond = bits
	return true
}

// parseARMCondition parses the conditions attached to an ARM instruction.
// The input is a single string consisting of period-separated condition
// codes, such as ".P.W". An initial period is ignored.
func parseARMCondition(cond string) (uint8, bool) {
	if strings.HasPrefix(cond, ".") {
		cond = cond[1:]
	}
	if cond == "" {
		return arm.C_SCOND_NONE, true
	}
	names := strings.Split(cond, ".")
	bits := uint8(0)
	for _, name := range names {
		if b, present := armLS[name]; present {
			bits |= b
			continue
		}
		if b, present := armSCOND[name]; present {
			bits = (bits &^ arm.C_SCOND) | b
			continue
		}
		return 0, false
	}
	return bits, true
}
