// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"cmd/internal/obj"
	"cmd/internal/obj/arm"
)
import "cmd/internal/gc"

const (
	RightRdwr = gc.RightRead | gc.RightWrite
)

// This table gives the basic information about instruction
// generated by the compiler and processed in the optimizer.
// See opt.h for bit definitions.
//
// Instructions not generated need not be listed.
// As an exception to that rule, we typically write down all the
// size variants of an operation even if we just use a subset.
//
// The table is formatted for 8-space tabs.
var progtable = [arm.ALAST]gc.ProgInfo{
	obj.ATYPE:     gc.ProgInfo{gc.Pseudo | gc.Skip, 0, 0, 0},
	obj.ATEXT:     gc.ProgInfo{gc.Pseudo, 0, 0, 0},
	obj.AFUNCDATA: gc.ProgInfo{gc.Pseudo, 0, 0, 0},
	obj.APCDATA:   gc.ProgInfo{gc.Pseudo, 0, 0, 0},
	obj.AUNDEF:    gc.ProgInfo{gc.Break, 0, 0, 0},
	obj.AUSEFIELD: gc.ProgInfo{gc.OK, 0, 0, 0},
	obj.ACHECKNIL: gc.ProgInfo{gc.LeftRead, 0, 0, 0},
	obj.AVARDEF:   gc.ProgInfo{gc.Pseudo | gc.RightWrite, 0, 0, 0},
	obj.AVARKILL:  gc.ProgInfo{gc.Pseudo | gc.RightWrite, 0, 0, 0},

	// NOP is an internal no-op that also stands
	// for USED and SET annotations, not the Intel opcode.
	obj.ANOP: gc.ProgInfo{gc.LeftRead | gc.RightWrite, 0, 0, 0},

	// Integer.
	arm.AADC:    gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite, 0, 0, 0},
	arm.AADD:    gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite, 0, 0, 0},
	arm.AAND:    gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite, 0, 0, 0},
	arm.ABIC:    gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite, 0, 0, 0},
	arm.ACMN:    gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RightRead, 0, 0, 0},
	arm.ACMP:    gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RightRead, 0, 0, 0},
	arm.ADIVU:   gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite, 0, 0, 0},
	arm.ADIV:    gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite, 0, 0, 0},
	arm.AEOR:    gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite, 0, 0, 0},
	arm.AMODU:   gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite, 0, 0, 0},
	arm.AMOD:    gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite, 0, 0, 0},
	arm.AMULALU: gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RegRead | RightRdwr, 0, 0, 0},
	arm.AMULAL:  gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RegRead | RightRdwr, 0, 0, 0},
	arm.AMULA:   gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RegRead | RightRdwr, 0, 0, 0},
	arm.AMULU:   gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite, 0, 0, 0},
	arm.AMUL:    gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite, 0, 0, 0},
	arm.AMULL:   gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite, 0, 0, 0},
	arm.AMULLU:  gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite, 0, 0, 0},
	arm.AMVN:    gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RightWrite, 0, 0, 0},
	arm.AORR:    gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite, 0, 0, 0},
	arm.ARSB:    gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite, 0, 0, 0},
	arm.ARSC:    gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite, 0, 0, 0},
	arm.ASBC:    gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite, 0, 0, 0},
	arm.ASLL:    gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite, 0, 0, 0},
	arm.ASRA:    gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite, 0, 0, 0},
	arm.ASRL:    gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite, 0, 0, 0},
	arm.ASUB:    gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite, 0, 0, 0},
	arm.ATEQ:    gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RightRead, 0, 0, 0},
	arm.ATST:    gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RightRead, 0, 0, 0},

	// Floating point.
	arm.AADDD: gc.ProgInfo{gc.SizeD | gc.LeftRead | RightRdwr, 0, 0, 0},
	arm.AADDF: gc.ProgInfo{gc.SizeF | gc.LeftRead | RightRdwr, 0, 0, 0},
	arm.ACMPD: gc.ProgInfo{gc.SizeD | gc.LeftRead | gc.RightRead, 0, 0, 0},
	arm.ACMPF: gc.ProgInfo{gc.SizeF | gc.LeftRead | gc.RightRead, 0, 0, 0},
	arm.ADIVD: gc.ProgInfo{gc.SizeD | gc.LeftRead | RightRdwr, 0, 0, 0},
	arm.ADIVF: gc.ProgInfo{gc.SizeF | gc.LeftRead | RightRdwr, 0, 0, 0},
	arm.AMULD: gc.ProgInfo{gc.SizeD | gc.LeftRead | RightRdwr, 0, 0, 0},
	arm.AMULF: gc.ProgInfo{gc.SizeF | gc.LeftRead | RightRdwr, 0, 0, 0},
	arm.ASUBD: gc.ProgInfo{gc.SizeD | gc.LeftRead | RightRdwr, 0, 0, 0},
	arm.ASUBF: gc.ProgInfo{gc.SizeF | gc.LeftRead | RightRdwr, 0, 0, 0},

	// Conversions.
	arm.AMOVWD: gc.ProgInfo{gc.SizeD | gc.LeftRead | gc.RightWrite | gc.Conv, 0, 0, 0},
	arm.AMOVWF: gc.ProgInfo{gc.SizeF | gc.LeftRead | gc.RightWrite | gc.Conv, 0, 0, 0},
	arm.AMOVDF: gc.ProgInfo{gc.SizeF | gc.LeftRead | gc.RightWrite | gc.Conv, 0, 0, 0},
	arm.AMOVDW: gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RightWrite | gc.Conv, 0, 0, 0},
	arm.AMOVFD: gc.ProgInfo{gc.SizeD | gc.LeftRead | gc.RightWrite | gc.Conv, 0, 0, 0},
	arm.AMOVFW: gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RightWrite | gc.Conv, 0, 0, 0},

	// Moves.
	arm.AMOVB: gc.ProgInfo{gc.SizeB | gc.LeftRead | gc.RightWrite | gc.Move, 0, 0, 0},
	arm.AMOVD: gc.ProgInfo{gc.SizeD | gc.LeftRead | gc.RightWrite | gc.Move, 0, 0, 0},
	arm.AMOVF: gc.ProgInfo{gc.SizeF | gc.LeftRead | gc.RightWrite | gc.Move, 0, 0, 0},
	arm.AMOVH: gc.ProgInfo{gc.SizeW | gc.LeftRead | gc.RightWrite | gc.Move, 0, 0, 0},
	arm.AMOVW: gc.ProgInfo{gc.SizeL | gc.LeftRead | gc.RightWrite | gc.Move, 0, 0, 0},

	// In addtion, duffzero reads R0,R1 and writes R1.  This fact is
	// encoded in peep.c
	obj.ADUFFZERO: gc.ProgInfo{gc.Call, 0, 0, 0},

	// In addtion, duffcopy reads R1,R2 and writes R0,R1,R2.  This fact is
	// encoded in peep.c
	obj.ADUFFCOPY: gc.ProgInfo{gc.Call, 0, 0, 0},

	// These should be split into the two different conversions instead
	// of overloading the one.
	arm.AMOVBS: gc.ProgInfo{gc.SizeB | gc.LeftRead | gc.RightWrite | gc.Conv, 0, 0, 0},
	arm.AMOVBU: gc.ProgInfo{gc.SizeB | gc.LeftRead | gc.RightWrite | gc.Conv, 0, 0, 0},
	arm.AMOVHS: gc.ProgInfo{gc.SizeW | gc.LeftRead | gc.RightWrite | gc.Conv, 0, 0, 0},
	arm.AMOVHU: gc.ProgInfo{gc.SizeW | gc.LeftRead | gc.RightWrite | gc.Conv, 0, 0, 0},

	// Jumps.
	arm.AB:   gc.ProgInfo{gc.Jump | gc.Break, 0, 0, 0},
	arm.ABL:  gc.ProgInfo{gc.Call, 0, 0, 0},
	arm.ABEQ: gc.ProgInfo{gc.Cjmp, 0, 0, 0},
	arm.ABNE: gc.ProgInfo{gc.Cjmp, 0, 0, 0},
	arm.ABCS: gc.ProgInfo{gc.Cjmp, 0, 0, 0},
	arm.ABHS: gc.ProgInfo{gc.Cjmp, 0, 0, 0},
	arm.ABCC: gc.ProgInfo{gc.Cjmp, 0, 0, 0},
	arm.ABLO: gc.ProgInfo{gc.Cjmp, 0, 0, 0},
	arm.ABMI: gc.ProgInfo{gc.Cjmp, 0, 0, 0},
	arm.ABPL: gc.ProgInfo{gc.Cjmp, 0, 0, 0},
	arm.ABVS: gc.ProgInfo{gc.Cjmp, 0, 0, 0},
	arm.ABVC: gc.ProgInfo{gc.Cjmp, 0, 0, 0},
	arm.ABHI: gc.ProgInfo{gc.Cjmp, 0, 0, 0},
	arm.ABLS: gc.ProgInfo{gc.Cjmp, 0, 0, 0},
	arm.ABGE: gc.ProgInfo{gc.Cjmp, 0, 0, 0},
	arm.ABLT: gc.ProgInfo{gc.Cjmp, 0, 0, 0},
	arm.ABGT: gc.ProgInfo{gc.Cjmp, 0, 0, 0},
	arm.ABLE: gc.ProgInfo{gc.Cjmp, 0, 0, 0},
	obj.ARET: gc.ProgInfo{gc.Break, 0, 0, 0},
}

func proginfo(p *obj.Prog) (info gc.ProgInfo) {
	info = progtable[p.As]
	if info.Flags == 0 {
		gc.Fatal("unknown instruction %v", p)
	}

	if p.From.Type == obj.TYPE_ADDR && p.From.Sym != nil && (info.Flags&gc.LeftRead != 0) {
		info.Flags &^= gc.LeftRead
		info.Flags |= gc.LeftAddr
	}

	if (info.Flags&gc.RegRead != 0) && p.Reg == 0 {
		info.Flags &^= gc.RegRead
		info.Flags |= gc.CanRegRead | gc.RightRead
	}

	if (p.Scond&arm.C_SCOND != arm.C_SCOND_NONE) && (info.Flags&gc.RightWrite != 0) {
		info.Flags |= gc.RightRead
	}

	switch p.As {
	case arm.ADIV,
		arm.ADIVU,
		arm.AMOD,
		arm.AMODU:
		info.Regset |= RtoB(arm.REG_R12)
	}

	return
}
