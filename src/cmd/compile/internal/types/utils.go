// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"cmd/internal/obj"
	"fmt"
)

const BADWIDTH = -1000000000

// Initialized by frontend. Exists only here.
var Tptr EType // either TPTR32 or TPTR64

// The following variables must be initialized early by the frontend.
// They are here to break import cycles.
// TODO(gri) eliminate these dependencies.
var (
	Widthptr    int
	Dowidth     func(*Type)
	Fatalf      func(string, ...interface{})
	Sconv       func(*Sym, int, int) string       // orig: func sconv(s *Sym, flag FmtFlag, mode fmtMode) string
	Tconv       func(*Type, int, int, int) string // orig: func tconv(t *Type, flag FmtFlag, mode fmtMode, depth int) string
	FormatSym   func(*Sym, fmt.State, rune, int)  // orig: func symFormat(sym *Sym, s fmt.State, verb rune, mode fmtMode)
	FormatType  func(*Type, fmt.State, rune, int) // orig: func typeFormat(t *Type, s fmt.State, verb rune, mode fmtMode)
	TypeLinkSym func(*Type) *obj.LSym
	Ctxt        *obj.Link

	FmtLeft     int
	FmtUnsigned int
	FErr        int
)

func (s *Sym) String() string {
	return Sconv(s, 0, FErr)
}

func (sym *Sym) Format(s fmt.State, verb rune) {
	FormatSym(sym, s, verb, FErr)
}

func (t *Type) String() string {
	// This is an external entry point, so we pass depth 0 to tconv.
	// The implementation of tconv (including typefmt and fldconv)
	// must take care not to use a type in a formatting string
	// to avoid resetting the recursion counter.
	return Tconv(t, 0, FErr, 0)
}

// ShortString generates a short description of t.
// It is used in autogenerated method names, reflection,
// and itab names.
func (t *Type) ShortString() string {
	return Tconv(t, FmtLeft, FErr, 0)
}

// LongString generates a complete description of t.
// It is useful for reflection,
// or when a unique fingerprint or hash of a type is required.
func (t *Type) LongString() string {
	return Tconv(t, FmtLeft|FmtUnsigned, FErr, 0)
}

func (t *Type) Format(s fmt.State, verb rune) {
	FormatType(t, s, verb, FErr)
}

type bitset8 uint8

func (f *bitset8) set(mask uint8, b bool) {
	if b {
		*(*uint8)(f) |= mask
	} else {
		*(*uint8)(f) &^= mask
	}
}

var etnames = []string{
	Txxx:        "Txxx",
	TINT:        "INT",
	TUINT:       "UINT",
	TINT8:       "INT8",
	TUINT8:      "UINT8",
	TINT16:      "INT16",
	TUINT16:     "UINT16",
	TINT32:      "INT32",
	TUINT32:     "UINT32",
	TINT64:      "INT64",
	TUINT64:     "UINT64",
	TUINTPTR:    "UINTPTR",
	TFLOAT32:    "FLOAT32",
	TFLOAT64:    "FLOAT64",
	TCOMPLEX64:  "COMPLEX64",
	TCOMPLEX128: "COMPLEX128",
	TBOOL:       "BOOL",
	TPTR32:      "PTR32",
	TPTR64:      "PTR64",
	TFUNC:       "FUNC",
	TARRAY:      "ARRAY",
	TSLICE:      "SLICE",
	TSTRUCT:     "STRUCT",
	TCHAN:       "CHAN",
	TMAP:        "MAP",
	TINTER:      "INTER",
	TFORW:       "FORW",
	TSTRING:     "STRING",
	TUNSAFEPTR:  "TUNSAFEPTR",
	TANY:        "ANY",
	TIDEAL:      "TIDEAL",
	TNIL:        "TNIL",
	TBLANK:      "TBLANK",
	TFUNCARGS:   "TFUNCARGS",
	TCHANARGS:   "TCHANARGS",
	TDDDFIELD:   "TDDDFIELD",
	TSSA:        "TSSA",
	TTUPLE:      "TTUPLE",
}

func (et EType) String() string {
	if int(et) < len(etnames) && etnames[et] != "" {
		return etnames[et]
	}
	return fmt.Sprintf("E-%d", et)
}
