// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"cmd/compile/internal/base"
	"cmd/internal/src"
)

var basicTypes = [...]struct {
	name  string
	etype Kind
}{
	{"int8", TINT8},
	{"int16", TINT16},
	{"int32", TINT32},
	{"int64", TINT64},
	{"uint8", TUINT8},
	{"uint16", TUINT16},
	{"uint32", TUINT32},
	{"uint64", TUINT64},
	{"float32", TFLOAT32},
	{"float64", TFLOAT64},
	{"complex64", TCOMPLEX64},
	{"complex128", TCOMPLEX128},
	{"bool", TBOOL},
	{"string", TSTRING},
}

var typedefs = [...]struct {
	name     string
	etype    Kind
	sameas32 Kind
	sameas64 Kind
}{
	{"int", TINT, TINT32, TINT64},
	{"uint", TUINT, TUINT32, TUINT64},
	{"uintptr", TUINTPTR, TUINT32, TUINT64},
}

func InitTypes(defTypeName func(sym *Sym, typ *Type) Object) {
	if PtrSize == 0 {
		base.Fatalf("typeinit before betypeinit")
	}

	SlicePtrOffset = 0
	SliceLenOffset = Rnd(SlicePtrOffset+int64(PtrSize), int64(PtrSize))
	SliceCapOffset = Rnd(SliceLenOffset+int64(PtrSize), int64(PtrSize))
	SliceSize = Rnd(SliceCapOffset+int64(PtrSize), int64(PtrSize))

	// string is same as slice wo the cap
	StringSize = Rnd(SliceLenOffset+int64(PtrSize), int64(PtrSize))

	for et := Kind(0); et < NTYPE; et++ {
		SimType[et] = et
	}

	Types[TANY] = newType(TANY)
	Types[TINTER] = NewInterface(LocalPkg, nil, false)

	defBasic := func(kind Kind, pkg *Pkg, name string) *Type {
		typ := newType(kind)
		obj := defTypeName(pkg.Lookup(name), typ)
		typ.sym = obj.Sym()
		typ.nod = obj
		if kind != TANY {
			CheckSize(typ)
		}
		return typ
	}

	for _, s := range &basicTypes {
		Types[s.etype] = defBasic(s.etype, BuiltinPkg, s.name)
	}

	for _, s := range &typedefs {
		sameas := s.sameas32
		if PtrSize == 8 {
			sameas = s.sameas64
		}
		SimType[s.etype] = sameas

		Types[s.etype] = defBasic(s.etype, BuiltinPkg, s.name)
	}

	// We create separate byte and rune types for better error messages
	// rather than just creating type alias *Sym's for the uint8 and
	// int32  Hence, (bytetype|runtype).Sym.isAlias() is false.
	// TODO(gri) Should we get rid of this special case (at the cost
	// of less informative error messages involving bytes and runes)?
	// (Alternatively, we could introduce an OTALIAS node representing
	// type aliases, albeit at the cost of having to deal with it everywhere).
	ByteType = defBasic(TUINT8, BuiltinPkg, "byte")
	RuneType = defBasic(TINT32, BuiltinPkg, "rune")

	// error type
	DeferCheckSize()
	ErrorType = defBasic(TFORW, BuiltinPkg, "error")
	ErrorType.SetUnderlying(makeErrorInterface())
	ResumeCheckSize()

	// comparable type (interface)
	DeferCheckSize()
	ComparableType = defBasic(TFORW, BuiltinPkg, "comparable")
	ComparableType.SetUnderlying(makeComparableInterface())
	ResumeCheckSize()

	// any type (interface)
	if base.Flag.G > 0 {
		DeferCheckSize()
		AnyType = defBasic(TFORW, BuiltinPkg, "any")
		AnyType.SetUnderlying(NewInterface(NoPkg, []*Field{}, false))
		ResumeCheckSize()
	}

	Types[TUNSAFEPTR] = defBasic(TUNSAFEPTR, UnsafePkg, "Pointer")

	Types[TBLANK] = newType(TBLANK)
	Types[TNIL] = newType(TNIL)

	// simple aliases
	SimType[TMAP] = TPTR
	SimType[TCHAN] = TPTR
	SimType[TFUNC] = TPTR
	SimType[TUNSAFEPTR] = TPTR

	for et := TINT8; et <= TUINT64; et++ {
		IsInt[et] = true
	}
	IsInt[TINT] = true
	IsInt[TUINT] = true
	IsInt[TUINTPTR] = true

	IsFloat[TFLOAT32] = true
	IsFloat[TFLOAT64] = true

	IsComplex[TCOMPLEX64] = true
	IsComplex[TCOMPLEX128] = true
}

func makeErrorInterface() *Type {
	sig := NewSignature(NoPkg, FakeRecv(), nil, nil, []*Field{
		NewField(src.NoXPos, nil, Types[TSTRING]),
	})
	method := NewField(src.NoXPos, LocalPkg.Lookup("Error"), sig)
	return NewInterface(NoPkg, []*Field{method}, false)
}

func makeComparableInterface() *Type {
	sig := NewSignature(NoPkg, FakeRecv(), nil, nil, nil)
	method := NewField(src.NoXPos, LocalPkg.Lookup("=="), sig)
	return NewInterface(NoPkg, []*Field{method}, false)
}
