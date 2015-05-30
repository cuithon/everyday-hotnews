// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package arm

import (
	"cmd/compile/internal/gc"
	"cmd/internal/obj"
	"cmd/internal/obj/arm"
)

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
var progtable = [arm.ALAST]obj.ProgInfo{
	obj.ATYPE:     {Flags: gc.Pseudo | gc.Skip},
	obj.ATEXT:     {Flags: gc.Pseudo},
	obj.AFUNCDATA: {Flags: gc.Pseudo},
	obj.APCDATA:   {Flags: gc.Pseudo},
	obj.AUNDEF:    {Flags: gc.Break},
	obj.AUSEFIELD: {Flags: gc.OK},
	obj.ACHECKNIL: {Flags: gc.LeftRead},
	obj.AVARDEF:   {Flags: gc.Pseudo | gc.RightWrite},
	obj.AVARKILL:  {Flags: gc.Pseudo | gc.RightWrite},

	// NOP is an internal no-op that also stands
	// for USED and SET annotations, not the Intel opcode.
	obj.ANOP: {Flags: gc.LeftRead | gc.RightWrite},

	// Integer.
	arm.AADC:    {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm.AADD:    {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm.AAND:    {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm.ABIC:    {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm.ACMN:    {Flags: gc.SizeL | gc.LeftRead | gc.RightRead},
	arm.ACMP:    {Flags: gc.SizeL | gc.LeftRead | gc.RightRead},
	arm.ADIVU:   {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm.ADIV:    {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm.AEOR:    {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm.AMODU:   {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm.AMOD:    {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm.AMULALU: {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | RightRdwr},
	arm.AMULAL:  {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | RightRdwr},
	arm.AMULA:   {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | RightRdwr},
	arm.AMULU:   {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm.AMUL:    {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm.AMULL:   {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm.AMULLU:  {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm.AMVN:    {Flags: gc.SizeL | gc.LeftRead | gc.RightWrite},
	arm.AORR:    {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm.ARSB:    {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm.ARSC:    {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm.ASBC:    {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm.ASLL:    {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm.ASRA:    {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm.ASRL:    {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm.ASUB:    {Flags: gc.SizeL | gc.LeftRead | gc.RegRead | gc.RightWrite},
	arm.ATEQ:    {Flags: gc.SizeL | gc.LeftRead | gc.RightRead},
	arm.ATST:    {Flags: gc.SizeL | gc.LeftRead | gc.RightRead},

	// Floating point.
	arm.AADDD:  {Flags: gc.SizeD | gc.LeftRead | RightRdwr},
	arm.AADDF:  {Flags: gc.SizeF | gc.LeftRead | RightRdwr},
	arm.ACMPD:  {Flags: gc.SizeD | gc.LeftRead | gc.RightRead},
	arm.ACMPF:  {Flags: gc.SizeF | gc.LeftRead | gc.RightRead},
	arm.ADIVD:  {Flags: gc.SizeD | gc.LeftRead | RightRdwr},
	arm.ADIVF:  {Flags: gc.SizeF | gc.LeftRead | RightRdwr},
	arm.AMULD:  {Flags: gc.SizeD | gc.LeftRead | RightRdwr},
	arm.AMULF:  {Flags: gc.SizeF | gc.LeftRead | RightRdwr},
	arm.ASUBD:  {Flags: gc.SizeD | gc.LeftRead | RightRdwr},
	arm.ASUBF:  {Flags: gc.SizeF | gc.LeftRead | RightRdwr},
	arm.ASQRTD: {Flags: gc.SizeD | gc.LeftRead | RightRdwr},

	// Conversions.
	arm.AMOVWD: {Flags: gc.SizeD | gc.LeftRead | gc.RightWrite | gc.Conv},
	arm.AMOVWF: {Flags: gc.SizeF | gc.LeftRead | gc.RightWrite | gc.Conv},
	arm.AMOVDF: {Flags: gc.SizeF | gc.LeftRead | gc.RightWrite | gc.Conv},
	arm.AMOVDW: {Flags: gc.SizeL | gc.LeftRead | gc.RightWrite | gc.Conv},
	arm.AMOVFD: {Flags: gc.SizeD | gc.LeftRead | gc.RightWrite | gc.Conv},
	arm.AMOVFW: {Flags: gc.SizeL | gc.LeftRead | gc.RightWrite | gc.Conv},

	// Moves.
	arm.AMOVB: {Flags: gc.SizeB | gc.LeftRead | gc.RightWrite | gc.Move},
	arm.AMOVD: {Flags: gc.SizeD | gc.LeftRead | gc.RightWrite | gc.Move},
	arm.AMOVF: {Flags: gc.SizeF | gc.LeftRead | gc.RightWrite | gc.Move},
	arm.AMOVH: {Flags: gc.SizeW | gc.LeftRead | gc.RightWrite | gc.Move},
	arm.AMOVW: {Flags: gc.SizeL | gc.LeftRead | gc.RightWrite | gc.Move},

	// In addtion, duffzero reads R0,R1 and writes R1.  This fact is
	// encoded in peep.c
	obj.ADUFFZERO: {Flags: gc.Call},

	// In addtion, duffcopy reads R1,R2 and writes R0,R1,R2.  This fact is
	// encoded in peep.c
	obj.ADUFFCOPY: {Flags: gc.Call},

	// These should be split into the two different conversions instead
	// of overloading the one.
	arm.AMOVBS: {Flags: gc.SizeB | gc.LeftRead | gc.RightWrite | gc.Conv},
	arm.AMOVBU: {Flags: gc.SizeB | gc.LeftRead | gc.RightWrite | gc.Conv},
	arm.AMOVHS: {Flags: gc.SizeW | gc.LeftRead | gc.RightWrite | gc.Conv},
	arm.AMOVHU: {Flags: gc.SizeW | gc.LeftRead | gc.RightWrite | gc.Conv},

	// Jumps.
	arm.AB:   {Flags: gc.Jump | gc.Break},
	arm.ABL:  {Flags: gc.Call},
	arm.ABEQ: {Flags: gc.Cjmp},
	arm.ABNE: {Flags: gc.Cjmp},
	arm.ABCS: {Flags: gc.Cjmp},
	arm.ABHS: {Flags: gc.Cjmp},
	arm.ABCC: {Flags: gc.Cjmp},
	arm.ABLO: {Flags: gc.Cjmp},
	arm.ABMI: {Flags: gc.Cjmp},
	arm.ABPL: {Flags: gc.Cjmp},
	arm.ABVS: {Flags: gc.Cjmp},
	arm.ABVC: {Flags: gc.Cjmp},
	arm.ABHI: {Flags: gc.Cjmp},
	arm.ABLS: {Flags: gc.Cjmp},
	arm.ABGE: {Flags: gc.Cjmp},
	arm.ABLT: {Flags: gc.Cjmp},
	arm.ABGT: {Flags: gc.Cjmp},
	arm.ABLE: {Flags: gc.Cjmp},
	obj.ARET: {Flags: gc.Break},
}

func proginfo(p *obj.Prog) {
	info := &p.Info
	*info = progtable[p.As]
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
}
