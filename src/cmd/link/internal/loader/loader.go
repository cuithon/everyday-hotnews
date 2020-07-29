// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package loader

import (
	"bytes"
	"cmd/internal/bio"
	"cmd/internal/goobj2"
	"cmd/internal/obj"
	"cmd/internal/objabi"
	"cmd/internal/sys"
	"cmd/link/internal/sym"
	"debug/elf"
	"fmt"
	"log"
	"math/bits"
	"os"
	"sort"
	"strings"
)

var _ = fmt.Print

// Sym encapsulates a global symbol index, used to identify a specific
// Go symbol. The 0-valued Sym is corresponds to an invalid symbol.
type Sym int

// Relocs encapsulates the set of relocations on a given symbol; an
// instance of this type is returned by the Loader Relocs() method.
type Relocs struct {
	rs []goobj2.Reloc

	li uint32   // local index of symbol whose relocs we're examining
	r  *oReader // object reader for containing package
	l  *Loader  // loader
}

// ExtReloc contains the payload for an external relocation.
type ExtReloc struct {
	Xsym Sym
	Xadd int64
	Type objabi.RelocType
	Size uint8
}

// Reloc holds a "handle" to access a relocation record from an
// object file.
type Reloc struct {
	*goobj2.Reloc
	r *oReader
	l *Loader

	// External reloc types may not fit into a uint8 which the Go object file uses.
	// Store it here, instead of in the byte of goobj2.Reloc.
	// For Go symbols this will always be zero.
	// goobj2.Reloc.Type() + typ is always the right type, for both Go and external
	// symbols.
	typ objabi.RelocType
}

func (rel Reloc) Type() objabi.RelocType { return objabi.RelocType(rel.Reloc.Type()) + rel.typ }
func (rel Reloc) Sym() Sym               { return rel.l.resolve(rel.r, rel.Reloc.Sym()) }
func (rel Reloc) SetSym(s Sym)           { rel.Reloc.SetSym(goobj2.SymRef{PkgIdx: 0, SymIdx: uint32(s)}) }

func (rel Reloc) SetType(t objabi.RelocType) {
	if t != objabi.RelocType(uint8(t)) {
		panic("SetType: type doesn't fit into Reloc")
	}
	rel.Reloc.SetType(uint8(t))
	if rel.typ != 0 {
		// should use SymbolBuilder.SetRelocType
		panic("wrong method to set reloc type")
	}
}

// Aux holds a "handle" to access an aux symbol record from an
// object file.
type Aux struct {
	*goobj2.Aux
	r *oReader
	l *Loader
}

func (a Aux) Sym() Sym { return a.l.resolve(a.r, a.Aux.Sym()) }

// oReader is a wrapper type of obj.Reader, along with some
// extra information.
type oReader struct {
	*goobj2.Reader
	unit         *sym.CompilationUnit
	version      int    // version of static symbol
	flags        uint32 // read from object file
	pkgprefix    string
	syms         []Sym  // Sym's global index, indexed by local index
	ndef         int    // cache goobj2.Reader.NSym()
	nhashed64def int    // cache goobj2.Reader.NHashed64Def()
	nhasheddef   int    // cache goobj2.Reader.NHashedDef()
	objidx       uint32 // index of this reader in the objs slice
}

// Total number of defined symbols (package symbols, hashed symbols, and
// non-package symbols).
func (r *oReader) NAlldef() int { return r.ndef + r.nhashed64def + r.nhasheddef + r.NNonpkgdef() }

type objIdx struct {
	r *oReader
	i Sym // start index
}

// objSym represents a symbol in an object file. It is a tuple of
// the object and the symbol's local index.
// For external symbols, objidx is the index of l.extReader (extObj),
// s is its index into the payload array.
// {0, 0} represents the nil symbol.
type objSym struct {
	objidx uint32 // index of the object (in l.objs array)
	s      uint32 // local index
}

type nameVer struct {
	name string
	v    int
}

type Bitmap []uint32

// set the i-th bit.
func (bm Bitmap) Set(i Sym) {
	n, r := uint(i)/32, uint(i)%32
	bm[n] |= 1 << r
}

// unset the i-th bit.
func (bm Bitmap) Unset(i Sym) {
	n, r := uint(i)/32, uint(i)%32
	bm[n] &^= (1 << r)
}

// whether the i-th bit is set.
func (bm Bitmap) Has(i Sym) bool {
	n, r := uint(i)/32, uint(i)%32
	return bm[n]&(1<<r) != 0
}

// return current length of bitmap in bits.
func (bm Bitmap) Len() int {
	return len(bm) * 32
}

// return the number of bits set.
func (bm Bitmap) Count() int {
	s := 0
	for _, x := range bm {
		s += bits.OnesCount32(x)
	}
	return s
}

func MakeBitmap(n int) Bitmap {
	return make(Bitmap, (n+31)/32)
}

// growBitmap insures that the specified bitmap has enough capacity,
// reallocating (doubling the size) if needed.
func growBitmap(reqLen int, b Bitmap) Bitmap {
	curLen := b.Len()
	if reqLen > curLen {
		b = append(b, MakeBitmap(reqLen+1-curLen)...)
	}
	return b
}

type symAndSize struct {
	sym  Sym
	size uint32
}

// A Loader loads new object files and resolves indexed symbol references.
//
// Notes on the layout of global symbol index space:
//
// - Go object files are read before host object files; each Go object
//   read adds its defined package symbols to the global index space.
//   Nonpackage symbols are not yet added.
//
// - In loader.LoadNonpkgSyms, add non-package defined symbols and
//   references in all object files to the global index space.
//
// - Host object file loading happens; the host object loader does a
//   name/version lookup for each symbol it finds; this can wind up
//   extending the external symbol index space range. The host object
//   loader stores symbol payloads in loader.payloads using SymbolBuilder.
//
// - Each symbol gets a unique global index. For duplicated and
//   overwriting/overwritten symbols, the second (or later) appearance
//   of the symbol gets the same global index as the first appearance.
type Loader struct {
	start       map[*oReader]Sym // map from object file to its start index
	objs        []objIdx         // sorted by start index (i.e. objIdx.i)
	extStart    Sym              // from this index on, the symbols are externally defined
	builtinSyms []Sym            // global index of builtin symbols

	objSyms []objSym // global index mapping to local index

	hashed64Syms  map[uint64]symAndSize          // short hashed (content-addressable) symbols, keyed by content hash
	hashedSyms    map[goobj2.HashType]symAndSize // hashed (content-addressable) symbols, keyed by content hash
	symsByName    [2]map[string]Sym              // map symbol name to index, two maps are for ABI0 and ABIInternal
	extStaticSyms map[nameVer]Sym                // externally defined static symbols, keyed by name

	extReader    *oReader // a dummy oReader, for external symbols
	payloadBatch []extSymPayload
	payloads     []*extSymPayload // contents of linker-materialized external syms
	values       []int64          // symbol values, indexed by global sym index

	sects    []*sym.Section // sections
	symSects []uint16       // symbol's section, index to sects array

	align []uint8 // symbol 2^N alignment, indexed by global index

	itablink         map[Sym]struct{} // itablink[j] defined if j is go.itablink.*
	deferReturnTramp map[Sym]bool     // whether the symbol is a trampoline of a deferreturn call

	objByPkg map[string]*oReader // map package path to its Go object reader

	anonVersion int // most recently assigned ext static sym pseudo-version

	// Bitmaps and other side structures used to store data used to store
	// symbol flags/attributes; these are to be accessed via the
	// corresponding loader "AttrXXX" and "SetAttrXXX" methods. Please
	// visit the comments on these methods for more details on the
	// semantics / interpretation of the specific flags or attribute.
	attrReachable        Bitmap // reachable symbols, indexed by global index
	attrOnList           Bitmap // "on list" symbols, indexed by global index
	attrLocal            Bitmap // "local" symbols, indexed by global index
	attrNotInSymbolTable Bitmap // "not in symtab" symbols, indexed by global idx
	attrUsedInIface      Bitmap // "used in interface" symbols, indexed by global idx
	attrVisibilityHidden Bitmap // hidden symbols, indexed by ext sym index
	attrDuplicateOK      Bitmap // dupOK symbols, indexed by ext sym index
	attrShared           Bitmap // shared symbols, indexed by ext sym index
	attrExternal         Bitmap // external symbols, indexed by ext sym index

	attrReadOnly         map[Sym]bool     // readonly data for this sym
	attrTopFrame         map[Sym]struct{} // top frame symbols
	attrSpecial          map[Sym]struct{} // "special" frame symbols
	attrCgoExportDynamic map[Sym]struct{} // "cgo_export_dynamic" symbols
	attrCgoExportStatic  map[Sym]struct{} // "cgo_export_static" symbols
	generatedSyms        map[Sym]struct{} // symbols that generate their content

	// Outer and Sub relations for symbols.
	// TODO: figure out whether it's more efficient to just have these
	// as fields on extSymPayload (note that this won't be a viable
	// strategy if somewhere in the linker we set sub/outer for a
	// non-external sym).
	outer map[Sym]Sym
	sub   map[Sym]Sym

	dynimplib   map[Sym]string      // stores Dynimplib symbol attribute
	dynimpvers  map[Sym]string      // stores Dynimpvers symbol attribute
	localentry  map[Sym]uint8       // stores Localentry symbol attribute
	extname     map[Sym]string      // stores Extname symbol attribute
	elfType     map[Sym]elf.SymType // stores elf type symbol property
	elfSym      map[Sym]int32       // stores elf sym symbol property
	localElfSym map[Sym]int32       // stores "local" elf sym symbol property
	symPkg      map[Sym]string      // stores package for symbol, or library for shlib-derived syms
	plt         map[Sym]int32       // stores dynimport for pe objects
	got         map[Sym]int32       // stores got for pe objects
	dynid       map[Sym]int32       // stores Dynid for symbol

	relocVariant map[relocId]sym.RelocVariant // stores variant relocs

	// Used to implement field tracking; created during deadcode if
	// field tracking is enabled. Reachparent[K] contains the index of
	// the symbol that triggered the marking of symbol K as live.
	Reachparent []Sym

	flags uint32

	hasUnknownPkgPath bool // if any Go object has unknown package path

	strictDupMsgs int // number of strict-dup warning/errors, when FlagStrictDups is enabled

	elfsetstring elfsetstringFunc

	errorReporter *ErrorReporter

	npkgsyms int // number of package symbols, for accounting
}

const (
	pkgDef = iota
	hashed64Def
	hashedDef
	nonPkgDef
	nonPkgRef
)

// objidx
const (
	nilObj = iota
	extObj
	goObjStart
)

type elfsetstringFunc func(str string, off int)

// extSymPayload holds the payload (data + relocations) for linker-synthesized
// external symbols (note that symbol value is stored in a separate slice).
type extSymPayload struct {
	name     string // TODO: would this be better as offset into str table?
	size     int64
	ver      int
	kind     sym.SymKind
	objidx   uint32 // index of original object if sym made by cloneToExternal
	relocs   []goobj2.Reloc
	reltypes []objabi.RelocType // relocation types
	data     []byte
	auxs     []goobj2.Aux
}

const (
	// Loader.flags
	FlagStrictDups = 1 << iota
)

func NewLoader(flags uint32, elfsetstring elfsetstringFunc, reporter *ErrorReporter) *Loader {
	nbuiltin := goobj2.NBuiltin()
	extReader := &oReader{objidx: extObj}
	ldr := &Loader{
		start:                make(map[*oReader]Sym),
		objs:                 []objIdx{{}, {extReader, 0}}, // reserve index 0 for nil symbol, 1 for external symbols
		objSyms:              make([]objSym, 1, 100000),    // reserve index 0 for nil symbol
		extReader:            extReader,
		hashed64Syms:         make(map[uint64]symAndSize, 10000),                                          // TODO: adjust preallocation sizes
		hashedSyms:           make(map[goobj2.HashType]symAndSize, 20000),                                 // TODO: adjust preallocation sizes
		symsByName:           [2]map[string]Sym{make(map[string]Sym, 80000), make(map[string]Sym, 50000)}, // preallocate ~2MB for ABI0 and ~1MB for ABI1 symbols
		objByPkg:             make(map[string]*oReader),
		outer:                make(map[Sym]Sym),
		sub:                  make(map[Sym]Sym),
		dynimplib:            make(map[Sym]string),
		dynimpvers:           make(map[Sym]string),
		localentry:           make(map[Sym]uint8),
		extname:              make(map[Sym]string),
		attrReadOnly:         make(map[Sym]bool),
		elfType:              make(map[Sym]elf.SymType),
		elfSym:               make(map[Sym]int32),
		localElfSym:          make(map[Sym]int32),
		symPkg:               make(map[Sym]string),
		plt:                  make(map[Sym]int32),
		got:                  make(map[Sym]int32),
		dynid:                make(map[Sym]int32),
		attrTopFrame:         make(map[Sym]struct{}),
		attrSpecial:          make(map[Sym]struct{}),
		attrCgoExportDynamic: make(map[Sym]struct{}),
		attrCgoExportStatic:  make(map[Sym]struct{}),
		generatedSyms:        make(map[Sym]struct{}),
		itablink:             make(map[Sym]struct{}),
		deferReturnTramp:     make(map[Sym]bool),
		extStaticSyms:        make(map[nameVer]Sym),
		builtinSyms:          make([]Sym, nbuiltin),
		flags:                flags,
		elfsetstring:         elfsetstring,
		errorReporter:        reporter,
		sects:                []*sym.Section{nil}, // reserve index 0 for nil section
	}
	reporter.ldr = ldr
	return ldr
}

// Add object file r, return the start index.
func (l *Loader) addObj(pkg string, r *oReader) Sym {
	if _, ok := l.start[r]; ok {
		panic("already added")
	}
	pkg = objabi.PathToPrefix(pkg) // the object file contains escaped package path
	if _, ok := l.objByPkg[pkg]; !ok {
		l.objByPkg[pkg] = r
	}
	i := Sym(len(l.objSyms))
	l.start[r] = i
	l.objs = append(l.objs, objIdx{r, i})
	if r.NeedNameExpansion() && !r.FromAssembly() {
		l.hasUnknownPkgPath = true
	}
	return i
}

// Add a symbol from an object file, return the global index.
// If the symbol already exist, it returns the index of that symbol.
func (l *Loader) addSym(name string, ver int, r *oReader, li uint32, kind int, osym *goobj2.Sym) Sym {
	if l.extStart != 0 {
		panic("addSym called after external symbol is created")
	}
	i := Sym(len(l.objSyms))
	addToGlobal := func() {
		l.objSyms = append(l.objSyms, objSym{r.objidx, li})
	}
	if name == "" && kind != hashed64Def && kind != hashedDef {
		addToGlobal()
		return i // unnamed aux symbol
	}
	if ver == r.version {
		// Static symbol. Add its global index but don't
		// add to name lookup table, as it cannot be
		// referenced by name.
		addToGlobal()
		return i
	}
	switch kind {
	case pkgDef:
		// Defined package symbols cannot be dup to each other.
		// We load all the package symbols first, so we don't need
		// to check dup here.
		// We still add it to the lookup table, as it may still be
		// referenced by name (e.g. through linkname).
		l.symsByName[ver][name] = i
		addToGlobal()
		return i
	case hashed64Def, hashedDef:
		// Hashed (content-addressable) symbol. Check the hash
		// but don't add to name lookup table, as they are not
		// referenced by name. Also no need to do overwriting
		// check, as same hash indicates same content.
		var checkHash func() (symAndSize, bool)
		var addToHashMap func(symAndSize)
		var h64 uint64         // only used for hashed64Def
		var h *goobj2.HashType // only used for hashedDef
		if kind == hashed64Def {
			checkHash = func() (symAndSize, bool) {
				h64 = r.Hash64(li - uint32(r.ndef))
				s, existed := l.hashed64Syms[h64]
				return s, existed
			}
			addToHashMap = func(ss symAndSize) { l.hashed64Syms[h64] = ss }
		} else {
			checkHash = func() (symAndSize, bool) {
				h = r.Hash(li - uint32(r.ndef+r.nhashed64def))
				s, existed := l.hashedSyms[*h]
				return s, existed
			}
			addToHashMap = func(ss symAndSize) { l.hashedSyms[*h] = ss }
		}
		siz := osym.Siz()
		if s, existed := checkHash(); existed {
			// The content hash is built from symbol data and relocations. In the
			// object file, the symbol data may not always contain trailing zeros,
			// e.g. for [5]int{1,2,3} and [100]int{1,2,3}, the data is same
			// (although the size is different).
			// Also, for short symbols, the content hash is the identity function of
			// the 8 bytes, and trailing zeros doesn't change the hash value, e.g.
			// hash("A") == hash("A\0\0\0").
			// So when two symbols have the same hash, we need to use the one with
			// larger size.
			if siz > s.size {
				// New symbol has larger size, use the new one. Rewrite the index mapping.
				l.objSyms[s.sym] = objSym{r.objidx, li}
				addToHashMap(symAndSize{s.sym, siz})
			}
			return s.sym
		}
		addToHashMap(symAndSize{i, siz})
		addToGlobal()
		return i
	}

	// Non-package (named) symbol. Check if it already exists.
	oldi, existed := l.symsByName[ver][name]
	if !existed {
		l.symsByName[ver][name] = i
		addToGlobal()
		return i
	}
	// symbol already exists
	if osym.Dupok() {
		if l.flags&FlagStrictDups != 0 {
			l.checkdup(name, r, li, oldi)
		}
		return oldi
	}
	oldr, oldli := l.toLocal(oldi)
	oldsym := oldr.Sym(oldli)
	if oldsym.Dupok() {
		return oldi
	}
	overwrite := r.DataSize(li) != 0
	if overwrite {
		// new symbol overwrites old symbol.
		oldtyp := sym.AbiSymKindToSymKind[objabi.SymKind(oldsym.Type())]
		if !(oldtyp.IsData() && oldr.DataSize(oldli) == 0) {
			log.Fatalf("duplicated definition of symbol " + name)
		}
		l.objSyms[oldi] = objSym{r.objidx, li}
	} else {
		// old symbol overwrites new symbol.
		typ := sym.AbiSymKindToSymKind[objabi.SymKind(oldsym.Type())]
		if !typ.IsData() { // only allow overwriting data symbol
			log.Fatalf("duplicated definition of symbol " + name)
		}
	}
	return oldi
}

// newExtSym creates a new external sym with the specified
// name/version.
func (l *Loader) newExtSym(name string, ver int) Sym {
	i := Sym(len(l.objSyms))
	if l.extStart == 0 {
		l.extStart = i
	}
	l.growValues(int(i) + 1)
	l.growAttrBitmaps(int(i) + 1)
	pi := l.newPayload(name, ver)
	l.objSyms = append(l.objSyms, objSym{l.extReader.objidx, uint32(pi)})
	l.extReader.syms = append(l.extReader.syms, i)
	return i
}

// LookupOrCreateSym looks up the symbol with the specified name/version,
// returning its Sym index if found. If the lookup fails, a new external
// Sym will be created, entered into the lookup tables, and returned.
func (l *Loader) LookupOrCreateSym(name string, ver int) Sym {
	i := l.Lookup(name, ver)
	if i != 0 {
		return i
	}
	i = l.newExtSym(name, ver)
	static := ver >= sym.SymVerStatic || ver < 0
	if static {
		l.extStaticSyms[nameVer{name, ver}] = i
	} else {
		l.symsByName[ver][name] = i
	}
	return i
}

func (l *Loader) IsExternal(i Sym) bool {
	r, _ := l.toLocal(i)
	return l.isExtReader(r)
}

func (l *Loader) isExtReader(r *oReader) bool {
	return r == l.extReader
}

// For external symbol, return its index in the payloads array.
// XXX result is actually not a global index. We (ab)use the Sym type
// so we don't need conversion for accessing bitmaps.
func (l *Loader) extIndex(i Sym) Sym {
	_, li := l.toLocal(i)
	return Sym(li)
}

// Get a new payload for external symbol, return its index in
// the payloads array.
func (l *Loader) newPayload(name string, ver int) int {
	pi := len(l.payloads)
	pp := l.allocPayload()
	pp.name = name
	pp.ver = ver
	l.payloads = append(l.payloads, pp)
	l.growExtAttrBitmaps()
	return pi
}

// getPayload returns a pointer to the extSymPayload struct for an
// external symbol if the symbol has a payload. Will panic if the
// symbol in question is bogus (zero or not an external sym).
func (l *Loader) getPayload(i Sym) *extSymPayload {
	if !l.IsExternal(i) {
		panic(fmt.Sprintf("bogus symbol index %d in getPayload", i))
	}
	pi := l.extIndex(i)
	return l.payloads[pi]
}

// allocPayload allocates a new payload.
func (l *Loader) allocPayload() *extSymPayload {
	batch := l.payloadBatch
	if len(batch) == 0 {
		batch = make([]extSymPayload, 1000)
	}
	p := &batch[0]
	l.payloadBatch = batch[1:]
	return p
}

func (ms *extSymPayload) Grow(siz int64) {
	if int64(int(siz)) != siz {
		log.Fatalf("symgrow size %d too long", siz)
	}
	if int64(len(ms.data)) >= siz {
		return
	}
	if cap(ms.data) < int(siz) {
		cl := len(ms.data)
		ms.data = append(ms.data, make([]byte, int(siz)+1-cl)...)
		ms.data = ms.data[0:cl]
	}
	ms.data = ms.data[:siz]
}

// Convert a local index to a global index.
func (l *Loader) toGlobal(r *oReader, i uint32) Sym {
	return r.syms[i]
}

// Convert a global index to a local index.
func (l *Loader) toLocal(i Sym) (*oReader, uint32) {
	return l.objs[l.objSyms[i].objidx].r, l.objSyms[i].s
}

// Resolve a local symbol reference. Return global index.
func (l *Loader) resolve(r *oReader, s goobj2.SymRef) Sym {
	var rr *oReader
	switch p := s.PkgIdx; p {
	case goobj2.PkgIdxInvalid:
		// {0, X} with non-zero X is never a valid sym reference from a Go object.
		// We steal this space for symbol references from external objects.
		// In this case, X is just the global index.
		if l.isExtReader(r) {
			return Sym(s.SymIdx)
		}
		if s.SymIdx != 0 {
			panic("bad sym ref")
		}
		return 0
	case goobj2.PkgIdxHashed64:
		i := int(s.SymIdx) + r.ndef
		return r.syms[i]
	case goobj2.PkgIdxHashed:
		i := int(s.SymIdx) + r.ndef + r.nhashed64def
		return r.syms[i]
	case goobj2.PkgIdxNone:
		i := int(s.SymIdx) + r.ndef + r.nhashed64def + r.nhasheddef
		return r.syms[i]
	case goobj2.PkgIdxBuiltin:
		return l.builtinSyms[s.SymIdx]
	case goobj2.PkgIdxSelf:
		rr = r
	default:
		pkg := r.Pkg(int(p))
		var ok bool
		rr, ok = l.objByPkg[pkg]
		if !ok {
			log.Fatalf("reference of nonexisted package %s, from %v", pkg, r.unit.Lib)
		}
	}
	return l.toGlobal(rr, s.SymIdx)
}

// Look up a symbol by name, return global index, or 0 if not found.
// This is more like Syms.ROLookup than Lookup -- it doesn't create
// new symbol.
func (l *Loader) Lookup(name string, ver int) Sym {
	if ver >= sym.SymVerStatic || ver < 0 {
		return l.extStaticSyms[nameVer{name, ver}]
	}
	return l.symsByName[ver][name]
}

// Check that duplicate symbols have same contents.
func (l *Loader) checkdup(name string, r *oReader, li uint32, dup Sym) {
	p := r.Data(li)
	rdup, ldup := l.toLocal(dup)
	pdup := rdup.Data(ldup)
	if bytes.Equal(p, pdup) {
		return
	}
	reason := "same length but different contents"
	if len(p) != len(pdup) {
		reason = fmt.Sprintf("new length %d != old length %d", len(p), len(pdup))
	}
	fmt.Fprintf(os.Stderr, "cmd/link: while reading object for '%v': duplicate symbol '%s', previous def at '%v', with mismatched payload: %s\n", r.unit.Lib, name, rdup.unit.Lib, reason)

	// For the moment, allow DWARF subprogram DIEs for
	// auto-generated wrapper functions. What seems to happen
	// here is that we get different line numbers on formal
	// params; I am guessing that the pos is being inherited
	// from the spot where the wrapper is needed.
	allowed := strings.HasPrefix(name, "go.info.go.interface") ||
		strings.HasPrefix(name, "go.info.go.builtin") ||
		strings.HasPrefix(name, "go.debuglines")
	if !allowed {
		l.strictDupMsgs++
	}
}

func (l *Loader) NStrictDupMsgs() int { return l.strictDupMsgs }

// Number of total symbols.
func (l *Loader) NSym() int {
	return len(l.objSyms)
}

// Number of defined Go symbols.
func (l *Loader) NDef() int {
	return int(l.extStart)
}

// Number of reachable symbols.
func (l *Loader) NReachableSym() int {
	return l.attrReachable.Count()
}

// SymNameLen returns the length of the symbol name, trying hard not to load
// the name.
func (l *Loader) SymNameLen(i Sym) int {
	// Not much we can do about external symbols.
	if l.IsExternal(i) {
		return len(l.SymName(i))
	}
	r, li := l.toLocal(i)
	le := r.Sym(li).NameLen(r.Reader)
	if !r.NeedNameExpansion() {
		return le
	}
	// Just load the symbol name. We don't know how expanded it'll be.
	return len(l.SymName(i))
}

// Returns the raw (unpatched) name of the i-th symbol.
func (l *Loader) RawSymName(i Sym) string {
	if l.IsExternal(i) {
		pp := l.getPayload(i)
		return pp.name
	}
	r, li := l.toLocal(i)
	return r.Sym(li).Name(r.Reader)
}

// Returns the (patched) name of the i-th symbol.
func (l *Loader) SymName(i Sym) string {
	if l.IsExternal(i) {
		pp := l.getPayload(i)
		return pp.name
	}
	r, li := l.toLocal(i)
	name := r.Sym(li).Name(r.Reader)
	if !r.NeedNameExpansion() {
		return name
	}
	return strings.Replace(name, "\"\".", r.pkgprefix, -1)
}

// Returns the version of the i-th symbol.
func (l *Loader) SymVersion(i Sym) int {
	if l.IsExternal(i) {
		pp := l.getPayload(i)
		return pp.ver
	}
	r, li := l.toLocal(i)
	return int(abiToVer(r.Sym(li).ABI(), r.version))
}

func (l *Loader) IsFileLocal(i Sym) bool {
	return l.SymVersion(i) >= sym.SymVerStatic
}

// IsFromAssembly returns true if this symbol is derived from an
// object file generated by the Go assembler.
func (l *Loader) IsFromAssembly(i Sym) bool {
	if l.IsExternal(i) {
		return false
	}
	r, _ := l.toLocal(i)
	return r.FromAssembly()
}

// Returns the type of the i-th symbol.
func (l *Loader) SymType(i Sym) sym.SymKind {
	if l.IsExternal(i) {
		pp := l.getPayload(i)
		if pp != nil {
			return pp.kind
		}
		return 0
	}
	r, li := l.toLocal(i)
	return sym.AbiSymKindToSymKind[objabi.SymKind(r.Sym(li).Type())]
}

// Returns the attributes of the i-th symbol.
func (l *Loader) SymAttr(i Sym) uint8 {
	if l.IsExternal(i) {
		// TODO: do something? External symbols have different representation of attributes.
		// For now, ReflectMethod, NoSplit, GoType, and Typelink are used and they cannot be
		// set by external symbol.
		return 0
	}
	r, li := l.toLocal(i)
	return r.Sym(li).Flag()
}

// Returns the size of the i-th symbol.
func (l *Loader) SymSize(i Sym) int64 {
	if l.IsExternal(i) {
		pp := l.getPayload(i)
		return pp.size
	}
	r, li := l.toLocal(i)
	return int64(r.Sym(li).Siz())
}

// AttrReachable returns true for symbols that are transitively
// referenced from the entry points. Unreachable symbols are not
// written to the output.
func (l *Loader) AttrReachable(i Sym) bool {
	return l.attrReachable.Has(i)
}

// SetAttrReachable sets the reachability property for a symbol (see
// AttrReachable).
func (l *Loader) SetAttrReachable(i Sym, v bool) {
	if v {
		l.attrReachable.Set(i)
	} else {
		l.attrReachable.Unset(i)
	}
}

// AttrOnList returns true for symbols that are on some list (such as
// the list of all text symbols, or one of the lists of data symbols)
// and is consulted to avoid bugs where a symbol is put on a list
// twice.
func (l *Loader) AttrOnList(i Sym) bool {
	return l.attrOnList.Has(i)
}

// SetAttrOnList sets the "on list" property for a symbol (see
// AttrOnList).
func (l *Loader) SetAttrOnList(i Sym, v bool) {
	if v {
		l.attrOnList.Set(i)
	} else {
		l.attrOnList.Unset(i)
	}
}

// AttrLocal returns true for symbols that are only visible within the
// module (executable or shared library) being linked. This attribute
// is applied to thunks and certain other linker-generated symbols.
func (l *Loader) AttrLocal(i Sym) bool {
	return l.attrLocal.Has(i)
}

// SetAttrLocal the "local" property for a symbol (see AttrLocal above).
func (l *Loader) SetAttrLocal(i Sym, v bool) {
	if v {
		l.attrLocal.Set(i)
	} else {
		l.attrLocal.Unset(i)
	}
}

// AttrUsedInIface returns true for a type symbol that is used in
// an interface.
func (l *Loader) AttrUsedInIface(i Sym) bool {
	return l.attrUsedInIface.Has(i)
}

func (l *Loader) SetAttrUsedInIface(i Sym, v bool) {
	if v {
		l.attrUsedInIface.Set(i)
	} else {
		l.attrUsedInIface.Unset(i)
	}
}

// SymAddr checks that a symbol is reachable, and returns its value.
func (l *Loader) SymAddr(i Sym) int64 {
	if !l.AttrReachable(i) {
		panic("unreachable symbol in symaddr")
	}
	return l.values[i]
}

// AttrNotInSymbolTable returns true for symbols that should not be
// added to the symbol table of the final generated load module.
func (l *Loader) AttrNotInSymbolTable(i Sym) bool {
	return l.attrNotInSymbolTable.Has(i)
}

// SetAttrNotInSymbolTable the "not in symtab" property for a symbol
// (see AttrNotInSymbolTable above).
func (l *Loader) SetAttrNotInSymbolTable(i Sym, v bool) {
	if v {
		l.attrNotInSymbolTable.Set(i)
	} else {
		l.attrNotInSymbolTable.Unset(i)
	}
}

// AttrVisibilityHidden symbols returns true for ELF symbols with
// visibility set to STV_HIDDEN. They become local symbols in
// the final executable. Only relevant when internally linking
// on an ELF platform.
func (l *Loader) AttrVisibilityHidden(i Sym) bool {
	if !l.IsExternal(i) {
		return false
	}
	return l.attrVisibilityHidden.Has(l.extIndex(i))
}

// SetAttrVisibilityHidden sets the "hidden visibility" property for a
// symbol (see AttrVisibilityHidden).
func (l *Loader) SetAttrVisibilityHidden(i Sym, v bool) {
	if !l.IsExternal(i) {
		panic("tried to set visibility attr on non-external symbol")
	}
	if v {
		l.attrVisibilityHidden.Set(l.extIndex(i))
	} else {
		l.attrVisibilityHidden.Unset(l.extIndex(i))
	}
}

// AttrDuplicateOK returns true for a symbol that can be present in
// multiple object files.
func (l *Loader) AttrDuplicateOK(i Sym) bool {
	if !l.IsExternal(i) {
		// TODO: if this path winds up being taken frequently, it
		// might make more sense to copy the flag value out of the object
		// into a larger bitmap during preload.
		r, li := l.toLocal(i)
		return r.Sym(li).Dupok()
	}
	return l.attrDuplicateOK.Has(l.extIndex(i))
}

// SetAttrDuplicateOK sets the "duplicate OK" property for an external
// symbol (see AttrDuplicateOK).
func (l *Loader) SetAttrDuplicateOK(i Sym, v bool) {
	if !l.IsExternal(i) {
		panic("tried to set dupok attr on non-external symbol")
	}
	if v {
		l.attrDuplicateOK.Set(l.extIndex(i))
	} else {
		l.attrDuplicateOK.Unset(l.extIndex(i))
	}
}

// AttrShared returns true for symbols compiled with the -shared option.
func (l *Loader) AttrShared(i Sym) bool {
	if !l.IsExternal(i) {
		// TODO: if this path winds up being taken frequently, it
		// might make more sense to copy the flag value out of the
		// object into a larger bitmap during preload.
		r, _ := l.toLocal(i)
		return r.Shared()
	}
	return l.attrShared.Has(l.extIndex(i))
}

// SetAttrShared sets the "shared" property for an external
// symbol (see AttrShared).
func (l *Loader) SetAttrShared(i Sym, v bool) {
	if !l.IsExternal(i) {
		panic(fmt.Sprintf("tried to set shared attr on non-external symbol %d %s", i, l.SymName(i)))
	}
	if v {
		l.attrShared.Set(l.extIndex(i))
	} else {
		l.attrShared.Unset(l.extIndex(i))
	}
}

// AttrExternal returns true for function symbols loaded from host
// object files.
func (l *Loader) AttrExternal(i Sym) bool {
	if !l.IsExternal(i) {
		return false
	}
	return l.attrExternal.Has(l.extIndex(i))
}

// SetAttrExternal sets the "external" property for an host object
// symbol (see AttrExternal).
func (l *Loader) SetAttrExternal(i Sym, v bool) {
	if !l.IsExternal(i) {
		panic(fmt.Sprintf("tried to set external attr on non-external symbol %q", l.RawSymName(i)))
	}
	if v {
		l.attrExternal.Set(l.extIndex(i))
	} else {
		l.attrExternal.Unset(l.extIndex(i))
	}
}

// AttrTopFrame returns true for a function symbol that is an entry
// point, meaning that unwinders should stop when they hit this
// function.
func (l *Loader) AttrTopFrame(i Sym) bool {
	_, ok := l.attrTopFrame[i]
	return ok
}

// SetAttrTopFrame sets the "top frame" property for a symbol (see
// AttrTopFrame).
func (l *Loader) SetAttrTopFrame(i Sym, v bool) {
	if v {
		l.attrTopFrame[i] = struct{}{}
	} else {
		delete(l.attrTopFrame, i)
	}
}

// AttrSpecial returns true for a symbols that do not have their
// address (i.e. Value) computed by the usual mechanism of
// data.go:dodata() & data.go:address().
func (l *Loader) AttrSpecial(i Sym) bool {
	_, ok := l.attrSpecial[i]
	return ok
}

// SetAttrSpecial sets the "special" property for a symbol (see
// AttrSpecial).
func (l *Loader) SetAttrSpecial(i Sym, v bool) {
	if v {
		l.attrSpecial[i] = struct{}{}
	} else {
		delete(l.attrSpecial, i)
	}
}

// AttrCgoExportDynamic returns true for a symbol that has been
// specially marked via the "cgo_export_dynamic" compiler directive
// written by cgo (in response to //export directives in the source).
func (l *Loader) AttrCgoExportDynamic(i Sym) bool {
	_, ok := l.attrCgoExportDynamic[i]
	return ok
}

// SetAttrCgoExportDynamic sets the "cgo_export_dynamic" for a symbol
// (see AttrCgoExportDynamic).
func (l *Loader) SetAttrCgoExportDynamic(i Sym, v bool) {
	if v {
		l.attrCgoExportDynamic[i] = struct{}{}
	} else {
		delete(l.attrCgoExportDynamic, i)
	}
}

// AttrCgoExportStatic returns true for a symbol that has been
// specially marked via the "cgo_export_static" directive
// written by cgo.
func (l *Loader) AttrCgoExportStatic(i Sym) bool {
	_, ok := l.attrCgoExportStatic[i]
	return ok
}

// SetAttrCgoExportStatic sets the "cgo_export_static" for a symbol
// (see AttrCgoExportStatic).
func (l *Loader) SetAttrCgoExportStatic(i Sym, v bool) {
	if v {
		l.attrCgoExportStatic[i] = struct{}{}
	} else {
		delete(l.attrCgoExportStatic, i)
	}
}

// IsGeneratedSym returns true if a symbol's been previously marked as a
// generator symbol through the SetIsGeneratedSym. The functions for generator
// symbols are kept in the Link context.
func (l *Loader) IsGeneratedSym(i Sym) bool {
	_, ok := l.generatedSyms[i]
	return ok
}

// SetIsGeneratedSym marks symbols as generated symbols. Data shouldn't be
// stored in generated symbols, and a function is registered and called for
// each of these symbols.
func (l *Loader) SetIsGeneratedSym(i Sym, v bool) {
	if !l.IsExternal(i) {
		panic("only external symbols can be generated")
	}
	if v {
		l.generatedSyms[i] = struct{}{}
	} else {
		delete(l.generatedSyms, i)
	}
}

func (l *Loader) AttrCgoExport(i Sym) bool {
	return l.AttrCgoExportDynamic(i) || l.AttrCgoExportStatic(i)
}

// AttrReadOnly returns true for a symbol whose underlying data
// is stored via a read-only mmap.
func (l *Loader) AttrReadOnly(i Sym) bool {
	if v, ok := l.attrReadOnly[i]; ok {
		return v
	}
	if l.IsExternal(i) {
		pp := l.getPayload(i)
		if pp.objidx != 0 {
			return l.objs[pp.objidx].r.ReadOnly()
		}
		return false
	}
	r, _ := l.toLocal(i)
	return r.ReadOnly()
}

// SetAttrReadOnly sets the "data is read only" property for a symbol
// (see AttrReadOnly).
func (l *Loader) SetAttrReadOnly(i Sym, v bool) {
	l.attrReadOnly[i] = v
}

// AttrSubSymbol returns true for symbols that are listed as a
// sub-symbol of some other outer symbol. The sub/outer mechanism is
// used when loading host objects (sections from the host object
// become regular linker symbols and symbols go on the Sub list of
// their section) and for constructing the global offset table when
// internally linking a dynamic executable.
//
// Note that in later stages of the linker, we set Outer(S) to some
// container symbol C, but don't set Sub(C). Thus we have two
// distinct scenarios:
//
// - Outer symbol covers the address ranges of its sub-symbols.
//   Outer.Sub is set in this case.
// - Outer symbol doesn't conver the address ranges. It is zero-sized
//   and doesn't have sub-symbols. In the case, the inner symbol is
//   not actually a "SubSymbol". (Tricky!)
//
// This method returns TRUE only for sub-symbols in the first scenario.
//
// FIXME: would be better to do away with this and have a better way
// to represent container symbols.

func (l *Loader) AttrSubSymbol(i Sym) bool {
	// we don't explicitly store this attribute any more -- return
	// a value based on the sub-symbol setting.
	o := l.OuterSym(i)
	if o == 0 {
		return false
	}
	return l.SubSym(o) != 0
}

// Note that we don't have a 'SetAttrSubSymbol' method in the loader;
// clients should instead use the AddInteriorSym method to establish
// containment relationships for host object symbols.

// Returns whether the i-th symbol has ReflectMethod attribute set.
func (l *Loader) IsReflectMethod(i Sym) bool {
	return l.SymAttr(i)&goobj2.SymFlagReflectMethod != 0
}

// Returns whether the i-th symbol is nosplit.
func (l *Loader) IsNoSplit(i Sym) bool {
	return l.SymAttr(i)&goobj2.SymFlagNoSplit != 0
}

// Returns whether this is a Go type symbol.
func (l *Loader) IsGoType(i Sym) bool {
	return l.SymAttr(i)&goobj2.SymFlagGoType != 0
}

// Returns whether this symbol should be included in typelink.
func (l *Loader) IsTypelink(i Sym) bool {
	return l.SymAttr(i)&goobj2.SymFlagTypelink != 0
}

// Returns whether this is a "go.itablink.*" symbol.
func (l *Loader) IsItabLink(i Sym) bool {
	if _, ok := l.itablink[i]; ok {
		return true
	}
	return false
}

// Return whether this is a trampoline of a deferreturn call.
func (l *Loader) IsDeferReturnTramp(i Sym) bool {
	return l.deferReturnTramp[i]
}

// Set that i is a trampoline of a deferreturn call.
func (l *Loader) SetIsDeferReturnTramp(i Sym, v bool) {
	l.deferReturnTramp[i] = v
}

// growValues grows the slice used to store symbol values.
func (l *Loader) growValues(reqLen int) {
	curLen := len(l.values)
	if reqLen > curLen {
		l.values = append(l.values, make([]int64, reqLen+1-curLen)...)
	}
}

// SymValue returns the value of the i-th symbol. i is global index.
func (l *Loader) SymValue(i Sym) int64 {
	return l.values[i]
}

// SetSymValue sets the value of the i-th symbol. i is global index.
func (l *Loader) SetSymValue(i Sym, val int64) {
	l.values[i] = val
}

// AddToSymValue adds to the value of the i-th symbol. i is the global index.
func (l *Loader) AddToSymValue(i Sym, val int64) {
	l.values[i] += val
}

// Returns the symbol content of the i-th symbol. i is global index.
func (l *Loader) Data(i Sym) []byte {
	if l.IsExternal(i) {
		pp := l.getPayload(i)
		if pp != nil {
			return pp.data
		}
		return nil
	}
	r, li := l.toLocal(i)
	return r.Data(li)
}

// FreeData clears the symbol data of an external symbol, allowing the memory
// to be freed earlier. No-op for non-external symbols.
// i is global index.
func (l *Loader) FreeData(i Sym) {
	if l.IsExternal(i) {
		pp := l.getPayload(i)
		if pp != nil {
			pp.data = nil
		}
	}
}

// SymAlign returns the alignment for a symbol.
func (l *Loader) SymAlign(i Sym) int32 {
	if int(i) >= len(l.align) {
		// align is extended lazily -- it the sym in question is
		// outside the range of the existing slice, then we assume its
		// alignment has not yet been set.
		return 0
	}
	// TODO: would it make sense to return an arch-specific
	// alignment depending on section type? E.g. STEXT => 32,
	// SDATA => 1, etc?
	abits := l.align[i]
	if abits == 0 {
		return 0
	}
	return int32(1 << (abits - 1))
}

// SetSymAlign sets the alignment for a symbol.
func (l *Loader) SetSymAlign(i Sym, align int32) {
	// Reject nonsense alignments.
	if align < 0 || align&(align-1) != 0 {
		panic("bad alignment value")
	}
	if int(i) >= len(l.align) {
		l.align = append(l.align, make([]uint8, l.NSym()-len(l.align))...)
	}
	if align == 0 {
		l.align[i] = 0
	}
	l.align[i] = uint8(bits.Len32(uint32(align)))
}

// SymValue returns the section of the i-th symbol. i is global index.
func (l *Loader) SymSect(i Sym) *sym.Section {
	if int(i) >= len(l.symSects) {
		// symSects is extended lazily -- it the sym in question is
		// outside the range of the existing slice, then we assume its
		// section has not yet been set.
		return nil
	}
	return l.sects[l.symSects[i]]
}

// SetSymSect sets the section of the i-th symbol. i is global index.
func (l *Loader) SetSymSect(i Sym, sect *sym.Section) {
	if int(i) >= len(l.symSects) {
		l.symSects = append(l.symSects, make([]uint16, l.NSym()-len(l.symSects))...)
	}
	l.symSects[i] = sect.Index
}

// growSects grows the slice used to store symbol sections.
func (l *Loader) growSects(reqLen int) {
	curLen := len(l.symSects)
	if reqLen > curLen {
		l.symSects = append(l.symSects, make([]uint16, reqLen+1-curLen)...)
	}
}

// NewSection creates a new (output) section.
func (l *Loader) NewSection() *sym.Section {
	sect := new(sym.Section)
	idx := len(l.sects)
	if idx != int(uint16(idx)) {
		panic("too many sections created")
	}
	sect.Index = uint16(idx)
	l.sects = append(l.sects, sect)
	return sect
}

// SymDynImplib returns the "dynimplib" attribute for the specified
// symbol, making up a portion of the info for a symbol specified
// on a "cgo_import_dynamic" compiler directive.
func (l *Loader) SymDynimplib(i Sym) string {
	return l.dynimplib[i]
}

// SetSymDynimplib sets the "dynimplib" attribute for a symbol.
func (l *Loader) SetSymDynimplib(i Sym, value string) {
	// reject bad symbols
	if i >= Sym(len(l.objSyms)) || i == 0 {
		panic("bad symbol index in SetDynimplib")
	}
	if value == "" {
		delete(l.dynimplib, i)
	} else {
		l.dynimplib[i] = value
	}
}

// SymDynimpvers returns the "dynimpvers" attribute for the specified
// symbol, making up a portion of the info for a symbol specified
// on a "cgo_import_dynamic" compiler directive.
func (l *Loader) SymDynimpvers(i Sym) string {
	return l.dynimpvers[i]
}

// SetSymDynimpvers sets the "dynimpvers" attribute for a symbol.
func (l *Loader) SetSymDynimpvers(i Sym, value string) {
	// reject bad symbols
	if i >= Sym(len(l.objSyms)) || i == 0 {
		panic("bad symbol index in SetDynimpvers")
	}
	if value == "" {
		delete(l.dynimpvers, i)
	} else {
		l.dynimpvers[i] = value
	}
}

// SymExtname returns the "extname" value for the specified
// symbol.
func (l *Loader) SymExtname(i Sym) string {
	if s, ok := l.extname[i]; ok {
		return s
	}
	return l.SymName(i)
}

// SetSymExtname sets the  "extname" attribute for a symbol.
func (l *Loader) SetSymExtname(i Sym, value string) {
	// reject bad symbols
	if i >= Sym(len(l.objSyms)) || i == 0 {
		panic("bad symbol index in SetExtname")
	}
	if value == "" {
		delete(l.extname, i)
	} else {
		l.extname[i] = value
	}
}

// SymElfType returns the previously recorded ELF type for a symbol
// (used only for symbols read from shared libraries by ldshlibsyms).
// It is not set for symbols defined by the packages being linked or
// by symbols read by ldelf (and so is left as elf.STT_NOTYPE).
func (l *Loader) SymElfType(i Sym) elf.SymType {
	if et, ok := l.elfType[i]; ok {
		return et
	}
	return elf.STT_NOTYPE
}

// SetSymElfType sets the elf type attribute for a symbol.
func (l *Loader) SetSymElfType(i Sym, et elf.SymType) {
	// reject bad symbols
	if i >= Sym(len(l.objSyms)) || i == 0 {
		panic("bad symbol index in SetSymElfType")
	}
	if et == elf.STT_NOTYPE {
		delete(l.elfType, i)
	} else {
		l.elfType[i] = et
	}
}

// SymElfSym returns the ELF symbol index for a given loader
// symbol, assigned during ELF symtab generation.
func (l *Loader) SymElfSym(i Sym) int32 {
	return l.elfSym[i]
}

// SetSymElfSym sets the elf symbol index for a symbol.
func (l *Loader) SetSymElfSym(i Sym, es int32) {
	if i == 0 {
		panic("bad sym index")
	}
	if es == 0 {
		delete(l.elfSym, i)
	} else {
		l.elfSym[i] = es
	}
}

// SymLocalElfSym returns the "local" ELF symbol index for a given loader
// symbol, assigned during ELF symtab generation.
func (l *Loader) SymLocalElfSym(i Sym) int32 {
	return l.localElfSym[i]
}

// SetSymLocalElfSym sets the "local" elf symbol index for a symbol.
func (l *Loader) SetSymLocalElfSym(i Sym, es int32) {
	if i == 0 {
		panic("bad sym index")
	}
	if es == 0 {
		delete(l.localElfSym, i)
	} else {
		l.localElfSym[i] = es
	}
}

// SymPlt returns the plt value for pe symbols.
func (l *Loader) SymPlt(s Sym) int32 {
	if v, ok := l.plt[s]; ok {
		return v
	}
	return -1
}

// SetPlt sets the plt value for pe symbols.
func (l *Loader) SetPlt(i Sym, v int32) {
	if i >= Sym(len(l.objSyms)) || i == 0 {
		panic("bad symbol for SetPlt")
	}
	if v == -1 {
		delete(l.plt, i)
	} else {
		l.plt[i] = v
	}
}

// SymGot returns the got value for pe symbols.
func (l *Loader) SymGot(s Sym) int32 {
	if v, ok := l.got[s]; ok {
		return v
	}
	return -1
}

// SetGot sets the got value for pe symbols.
func (l *Loader) SetGot(i Sym, v int32) {
	if i >= Sym(len(l.objSyms)) || i == 0 {
		panic("bad symbol for SetGot")
	}
	if v == -1 {
		delete(l.got, i)
	} else {
		l.got[i] = v
	}
}

// SymDynid returns the "dynid" property for the specified symbol.
func (l *Loader) SymDynid(i Sym) int32 {
	if s, ok := l.dynid[i]; ok {
		return s
	}
	return -1
}

// SetSymDynid sets the "dynid" property for a symbol.
func (l *Loader) SetSymDynid(i Sym, val int32) {
	// reject bad symbols
	if i >= Sym(len(l.objSyms)) || i == 0 {
		panic("bad symbol index in SetSymDynid")
	}
	if val == -1 {
		delete(l.dynid, i)
	} else {
		l.dynid[i] = val
	}
}

// DynIdSyms returns the set of symbols for which dynID is set to an
// interesting (non-default) value. This is expected to be a fairly
// small set.
func (l *Loader) DynidSyms() []Sym {
	sl := make([]Sym, 0, len(l.dynid))
	for s := range l.dynid {
		sl = append(sl, s)
	}
	sort.Slice(sl, func(i, j int) bool { return sl[i] < sl[j] })
	return sl
}

// SymGoType returns the 'Gotype' property for a given symbol (set by
// the Go compiler for variable symbols). This version relies on
// reading aux symbols for the target sym -- it could be that a faster
// approach would be to check for gotype during preload and copy the
// results in to a map (might want to try this at some point and see
// if it helps speed things up).
func (l *Loader) SymGoType(i Sym) Sym {
	var r *oReader
	var auxs []goobj2.Aux
	if l.IsExternal(i) {
		pp := l.getPayload(i)
		r = l.objs[pp.objidx].r
		auxs = pp.auxs
	} else {
		var li uint32
		r, li = l.toLocal(i)
		auxs = r.Auxs(li)
	}
	for j := range auxs {
		a := &auxs[j]
		switch a.Type() {
		case goobj2.AuxGotype:
			return l.resolve(r, a.Sym())
		}
	}
	return 0
}

// SymUnit returns the compilation unit for a given symbol (which will
// typically be nil for external or linker-manufactured symbols).
func (l *Loader) SymUnit(i Sym) *sym.CompilationUnit {
	if l.IsExternal(i) {
		pp := l.getPayload(i)
		if pp.objidx != 0 {
			r := l.objs[pp.objidx].r
			return r.unit
		}
		return nil
	}
	r, _ := l.toLocal(i)
	return r.unit
}

// SymPkg returns the package where the symbol came from (for
// regular compiler-generated Go symbols), but in the case of
// building with "-linkshared" (when a symbol is read from a
// shared library), will hold the library name.
// NOTE: this correspondes to sym.Symbol.File field.
func (l *Loader) SymPkg(i Sym) string {
	if f, ok := l.symPkg[i]; ok {
		return f
	}
	if l.IsExternal(i) {
		pp := l.getPayload(i)
		if pp.objidx != 0 {
			r := l.objs[pp.objidx].r
			return r.unit.Lib.Pkg
		}
		return ""
	}
	r, _ := l.toLocal(i)
	return r.unit.Lib.Pkg
}

// SetSymPkg sets the package/library for a symbol. This is
// needed mainly for external symbols, specifically those imported
// from shared libraries.
func (l *Loader) SetSymPkg(i Sym, pkg string) {
	// reject bad symbols
	if i >= Sym(len(l.objSyms)) || i == 0 {
		panic("bad symbol index in SetSymPkg")
	}
	l.symPkg[i] = pkg
}

// SymLocalentry returns the "local entry" value for the specified
// symbol.
func (l *Loader) SymLocalentry(i Sym) uint8 {
	return l.localentry[i]
}

// SetSymLocalentry sets the "local entry" attribute for a symbol.
func (l *Loader) SetSymLocalentry(i Sym, value uint8) {
	// reject bad symbols
	if i >= Sym(len(l.objSyms)) || i == 0 {
		panic("bad symbol index in SetSymLocalentry")
	}
	if value == 0 {
		delete(l.localentry, i)
	} else {
		l.localentry[i] = value
	}
}

// Returns the number of aux symbols given a global index.
func (l *Loader) NAux(i Sym) int {
	if l.IsExternal(i) {
		return 0
	}
	r, li := l.toLocal(i)
	return r.NAux(li)
}

// Returns the "handle" to the j-th aux symbol of the i-th symbol.
func (l *Loader) Aux(i Sym, j int) Aux {
	if l.IsExternal(i) {
		return Aux{}
	}
	r, li := l.toLocal(i)
	if j >= r.NAux(li) {
		return Aux{}
	}
	return Aux{r.Aux(li, j), r, l}
}

// GetFuncDwarfAuxSyms collects and returns the auxiliary DWARF
// symbols associated with a given function symbol.  Prior to the
// introduction of the loader, this was done purely using name
// lookups, e.f. for function with name XYZ we would then look up
// go.info.XYZ, etc.
func (l *Loader) GetFuncDwarfAuxSyms(fnSymIdx Sym) (auxDwarfInfo, auxDwarfLoc, auxDwarfRanges, auxDwarfLines Sym) {
	if l.SymType(fnSymIdx) != sym.STEXT {
		log.Fatalf("error: non-function sym %d/%s t=%s passed to GetFuncDwarfAuxSyms", fnSymIdx, l.SymName(fnSymIdx), l.SymType(fnSymIdx).String())
	}
	if l.IsExternal(fnSymIdx) {
		// Current expectation is that any external function will
		// not have auxsyms.
		return
	}
	r, li := l.toLocal(fnSymIdx)
	auxs := r.Auxs(li)
	for i := range auxs {
		a := &auxs[i]
		switch a.Type() {
		case goobj2.AuxDwarfInfo:
			auxDwarfInfo = l.resolve(r, a.Sym())
			if l.SymType(auxDwarfInfo) != sym.SDWARFFCN {
				panic("aux dwarf info sym with wrong type")
			}
		case goobj2.AuxDwarfLoc:
			auxDwarfLoc = l.resolve(r, a.Sym())
			if l.SymType(auxDwarfLoc) != sym.SDWARFLOC {
				panic("aux dwarf loc sym with wrong type")
			}
		case goobj2.AuxDwarfRanges:
			auxDwarfRanges = l.resolve(r, a.Sym())
			if l.SymType(auxDwarfRanges) != sym.SDWARFRANGE {
				panic("aux dwarf ranges sym with wrong type")
			}
		case goobj2.AuxDwarfLines:
			auxDwarfLines = l.resolve(r, a.Sym())
			if l.SymType(auxDwarfLines) != sym.SDWARFLINES {
				panic("aux dwarf lines sym with wrong type")
			}
		}
	}
	return
}

// AddInteriorSym sets up 'interior' as an interior symbol of
// container/payload symbol 'container'. An interior symbol does not
// itself have data, but gives a name to a subrange of the data in its
// container symbol. The container itself may or may not have a name.
// This method is intended primarily for use in the host object
// loaders, to capture the semantics of symbols and sections in an
// object file. When reading a host object file, we'll typically
// encounter a static section symbol (ex: ".text") containing content
// for a collection of functions, then a series of ELF (or macho, etc)
// symbol table entries each of which points into a sub-section
// (offset and length) of its corresponding container symbol. Within
// the go linker we create a loader.Sym for the container (which is
// expected to have the actual content/payload) and then a set of
// interior loader.Sym's that point into a portion of the container.
func (l *Loader) AddInteriorSym(container Sym, interior Sym) {
	// Container symbols are expected to have content/data.
	// NB: this restriction may turn out to be too strict (it's possible
	// to imagine a zero-sized container with an interior symbol pointing
	// into it); it's ok to relax or remove it if we counter an
	// oddball host object that triggers this.
	if l.SymSize(container) == 0 && len(l.Data(container)) == 0 {
		panic("unexpected empty container symbol")
	}
	// The interior symbols for a container are not expected to have
	// content/data or relocations.
	if len(l.Data(interior)) != 0 {
		panic("unexpected non-empty interior symbol")
	}
	// Interior symbol is expected to be in the symbol table.
	if l.AttrNotInSymbolTable(interior) {
		panic("interior symbol must be in symtab")
	}
	// Only a single level of containment is allowed.
	if l.OuterSym(container) != 0 {
		panic("outer has outer itself")
	}
	// Interior sym should not already have a sibling.
	if l.SubSym(interior) != 0 {
		panic("sub set for subsym")
	}
	// Interior sym should not already point at a container.
	if l.OuterSym(interior) != 0 {
		panic("outer already set for subsym")
	}
	l.sub[interior] = l.sub[container]
	l.sub[container] = interior
	l.outer[interior] = container
}

// OuterSym gets the outer symbol for host object loaded symbols.
func (l *Loader) OuterSym(i Sym) Sym {
	// FIXME: add check for isExternal?
	return l.outer[i]
}

// SubSym gets the subsymbol for host object loaded symbols.
func (l *Loader) SubSym(i Sym) Sym {
	// NB: note -- no check for l.isExternal(), since I am pretty sure
	// that later phases in the linker set subsym for "type." syms
	return l.sub[i]
}

// SetCarrierSym declares that 'c' is the carrier or container symbol
// for 's'. Carrier symbols are used in the linker to as a container
// for a collection of sub-symbols where the content of the
// sub-symbols is effectively concatenated to form the content of the
// carrier. The carrier is given a name in the output symbol table
// while the sub-symbol names are not. For example, the Go compiler
// emits named string symbols (type SGOSTRING) when compiling a
// package; after being deduplicated, these symbols are collected into
// a single unit by assigning them a new carrier symbol named
// "go.string.*" (which appears in the final symbol table for the
// output load module).
func (l *Loader) SetCarrierSym(s Sym, c Sym) {
	if c == 0 {
		panic("invalid carrier in SetCarrierSym")
	}
	if s == 0 {
		panic("invalid sub-symbol in SetCarrierSym")
	}
	// Carrier symbols are not expected to have content/data. It is
	// ok for them to have non-zero size (to allow for use of generator
	// symbols).
	if len(l.Data(c)) != 0 {
		panic("unexpected non-empty carrier symbol")
	}
	l.outer[s] = c
	// relocsym's foldSubSymbolOffset requires that we only
	// have a single level of containment-- enforce here.
	if l.outer[c] != 0 {
		panic("invalid nested carrier sym")
	}
}

// Initialize Reachable bitmap and its siblings for running deadcode pass.
func (l *Loader) InitReachable() {
	l.growAttrBitmaps(l.NSym() + 1)
}

type symWithVal struct {
	s Sym
	v int64
}
type bySymValue []symWithVal

func (s bySymValue) Len() int           { return len(s) }
func (s bySymValue) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s bySymValue) Less(i, j int) bool { return s[i].v < s[j].v }

// SortSub walks through the sub-symbols for 's' and sorts them
// in place by increasing value. Return value is the new
// sub symbol for the specified outer symbol.
func (l *Loader) SortSub(s Sym) Sym {

	if s == 0 || l.sub[s] == 0 {
		return s
	}

	// Sort symbols using a slice first. Use a stable sort on the off
	// chance that there's more than once symbol with the same value,
	// so as to preserve reproducible builds.
	sl := []symWithVal{}
	for ss := l.sub[s]; ss != 0; ss = l.sub[ss] {
		sl = append(sl, symWithVal{s: ss, v: l.SymValue(ss)})
	}
	sort.Stable(bySymValue(sl))

	// Then apply any changes needed to the sub map.
	ns := Sym(0)
	for i := len(sl) - 1; i >= 0; i-- {
		s := sl[i].s
		l.sub[s] = ns
		ns = s
	}

	// Update sub for outer symbol, then return
	l.sub[s] = sl[0].s
	return sl[0].s
}

// Insure that reachable bitmap and its siblings have enough size.
func (l *Loader) growAttrBitmaps(reqLen int) {
	if reqLen > l.attrReachable.Len() {
		// These are indexed by global symbol
		l.attrReachable = growBitmap(reqLen, l.attrReachable)
		l.attrOnList = growBitmap(reqLen, l.attrOnList)
		l.attrLocal = growBitmap(reqLen, l.attrLocal)
		l.attrNotInSymbolTable = growBitmap(reqLen, l.attrNotInSymbolTable)
		l.attrUsedInIface = growBitmap(reqLen, l.attrUsedInIface)
	}
	l.growExtAttrBitmaps()
}

func (l *Loader) growExtAttrBitmaps() {
	// These are indexed by external symbol index (e.g. l.extIndex(i))
	extReqLen := len(l.payloads)
	if extReqLen > l.attrVisibilityHidden.Len() {
		l.attrVisibilityHidden = growBitmap(extReqLen, l.attrVisibilityHidden)
		l.attrDuplicateOK = growBitmap(extReqLen, l.attrDuplicateOK)
		l.attrShared = growBitmap(extReqLen, l.attrShared)
		l.attrExternal = growBitmap(extReqLen, l.attrExternal)
	}
}

func (relocs *Relocs) Count() int { return len(relocs.rs) }

// At returns the j-th reloc for a global symbol.
func (relocs *Relocs) At(j int) Reloc {
	if relocs.l.isExtReader(relocs.r) {
		pp := relocs.l.payloads[relocs.li]
		return Reloc{&relocs.rs[j], relocs.r, relocs.l, pp.reltypes[j]}
	}
	return Reloc{&relocs.rs[j], relocs.r, relocs.l, 0}
}

// Relocs returns a Relocs object for the given global sym.
func (l *Loader) Relocs(i Sym) Relocs {
	r, li := l.toLocal(i)
	if r == nil {
		panic(fmt.Sprintf("trying to get oreader for invalid sym %d\n\n", i))
	}
	return l.relocs(r, li)
}

// Relocs returns a Relocs object given a local sym index and reader.
func (l *Loader) relocs(r *oReader, li uint32) Relocs {
	var rs []goobj2.Reloc
	if l.isExtReader(r) {
		pp := l.payloads[li]
		rs = pp.relocs
	} else {
		rs = r.Relocs(li)
	}
	return Relocs{
		rs: rs,
		li: li,
		r:  r,
		l:  l,
	}
}

// FuncInfo provides hooks to access goobj2.FuncInfo in the objects.
type FuncInfo struct {
	l       *Loader
	r       *oReader
	data    []byte
	auxs    []goobj2.Aux
	lengths goobj2.FuncInfoLengths
}

func (fi *FuncInfo) Valid() bool { return fi.r != nil }

func (fi *FuncInfo) Args() int {
	return int((*goobj2.FuncInfo)(nil).ReadArgs(fi.data))
}

func (fi *FuncInfo) Locals() int {
	return int((*goobj2.FuncInfo)(nil).ReadLocals(fi.data))
}

func (fi *FuncInfo) FuncID() objabi.FuncID {
	return objabi.FuncID((*goobj2.FuncInfo)(nil).ReadFuncID(fi.data))
}

func (fi *FuncInfo) Pcsp() []byte {
	pcsp, end := (*goobj2.FuncInfo)(nil).ReadPcsp(fi.data)
	return fi.r.BytesAt(fi.r.PcdataBase()+pcsp, int(end-pcsp))
}

func (fi *FuncInfo) Pcfile() []byte {
	pcf, end := (*goobj2.FuncInfo)(nil).ReadPcfile(fi.data)
	return fi.r.BytesAt(fi.r.PcdataBase()+pcf, int(end-pcf))
}

func (fi *FuncInfo) Pcline() []byte {
	pcln, end := (*goobj2.FuncInfo)(nil).ReadPcline(fi.data)
	return fi.r.BytesAt(fi.r.PcdataBase()+pcln, int(end-pcln))
}

// Preload has to be called prior to invoking the various methods
// below related to pcdata, funcdataoff, files, and inltree nodes.
func (fi *FuncInfo) Preload() {
	fi.lengths = (*goobj2.FuncInfo)(nil).ReadFuncInfoLengths(fi.data)
}

func (fi *FuncInfo) Pcinline() []byte {
	if !fi.lengths.Initialized {
		panic("need to call Preload first")
	}
	pcinl, end := (*goobj2.FuncInfo)(nil).ReadPcinline(fi.data, fi.lengths.PcdataOff)
	return fi.r.BytesAt(fi.r.PcdataBase()+pcinl, int(end-pcinl))
}

func (fi *FuncInfo) NumPcdata() uint32 {
	if !fi.lengths.Initialized {
		panic("need to call Preload first")
	}
	return fi.lengths.NumPcdata
}

func (fi *FuncInfo) Pcdata(k int) []byte {
	if !fi.lengths.Initialized {
		panic("need to call Preload first")
	}
	pcdat, end := (*goobj2.FuncInfo)(nil).ReadPcdata(fi.data, fi.lengths.PcdataOff, uint32(k))
	return fi.r.BytesAt(fi.r.PcdataBase()+pcdat, int(end-pcdat))
}

func (fi *FuncInfo) NumFuncdataoff() uint32 {
	if !fi.lengths.Initialized {
		panic("need to call Preload first")
	}
	return fi.lengths.NumFuncdataoff
}

func (fi *FuncInfo) Funcdataoff(k int) int64 {
	if !fi.lengths.Initialized {
		panic("need to call Preload first")
	}
	return (*goobj2.FuncInfo)(nil).ReadFuncdataoff(fi.data, fi.lengths.FuncdataoffOff, uint32(k))
}

func (fi *FuncInfo) Funcdata(syms []Sym) []Sym {
	if !fi.lengths.Initialized {
		panic("need to call Preload first")
	}
	if int(fi.lengths.NumFuncdataoff) > cap(syms) {
		syms = make([]Sym, 0, fi.lengths.NumFuncdataoff)
	} else {
		syms = syms[:0]
	}
	for j := range fi.auxs {
		a := &fi.auxs[j]
		if a.Type() == goobj2.AuxFuncdata {
			syms = append(syms, fi.l.resolve(fi.r, a.Sym()))
		}
	}
	return syms
}

func (fi *FuncInfo) NumFile() uint32 {
	if !fi.lengths.Initialized {
		panic("need to call Preload first")
	}
	return fi.lengths.NumFile
}

func (fi *FuncInfo) File(k int) Sym {
	if !fi.lengths.Initialized {
		panic("need to call Preload first")
	}
	sr := (*goobj2.FuncInfo)(nil).ReadFile(fi.data, fi.lengths.FileOff, uint32(k))
	return fi.l.resolve(fi.r, sr)
}

type InlTreeNode struct {
	Parent   int32
	File     Sym
	Line     int32
	Func     Sym
	ParentPC int32
}

func (fi *FuncInfo) NumInlTree() uint32 {
	if !fi.lengths.Initialized {
		panic("need to call Preload first")
	}
	return fi.lengths.NumInlTree
}

func (fi *FuncInfo) InlTree(k int) InlTreeNode {
	if !fi.lengths.Initialized {
		panic("need to call Preload first")
	}
	node := (*goobj2.FuncInfo)(nil).ReadInlTree(fi.data, fi.lengths.InlTreeOff, uint32(k))
	return InlTreeNode{
		Parent:   node.Parent,
		File:     fi.l.resolve(fi.r, node.File),
		Line:     node.Line,
		Func:     fi.l.resolve(fi.r, node.Func),
		ParentPC: node.ParentPC,
	}
}

func (l *Loader) FuncInfo(i Sym) FuncInfo {
	var r *oReader
	var auxs []goobj2.Aux
	if l.IsExternal(i) {
		pp := l.getPayload(i)
		if pp.objidx == 0 {
			return FuncInfo{}
		}
		r = l.objs[pp.objidx].r
		auxs = pp.auxs
	} else {
		var li uint32
		r, li = l.toLocal(i)
		auxs = r.Auxs(li)
	}
	for j := range auxs {
		a := &auxs[j]
		if a.Type() == goobj2.AuxFuncInfo {
			b := r.Data(a.Sym().SymIdx)
			return FuncInfo{l, r, b, auxs, goobj2.FuncInfoLengths{}}
		}
	}
	return FuncInfo{}
}

// Preload a package: add autolibs, add defined package symbols to the symbol table.
// Does not add non-package symbols yet, which will be done in LoadNonpkgSyms.
// Does not read symbol data.
// Returns the fingerprint of the object.
func (l *Loader) Preload(localSymVersion int, f *bio.Reader, lib *sym.Library, unit *sym.CompilationUnit, length int64) goobj2.FingerprintType {
	roObject, readonly, err := f.Slice(uint64(length)) // TODO: no need to map blocks that are for tools only (e.g. RefName)
	if err != nil {
		log.Fatal("cannot read object file:", err)
	}
	r := goobj2.NewReaderFromBytes(roObject, readonly)
	if r == nil {
		if len(roObject) >= 8 && bytes.Equal(roObject[:8], []byte("\x00go114ld")) {
			log.Fatalf("found object file %s in old format", f.File().Name())
		}
		panic("cannot read object file")
	}
	pkgprefix := objabi.PathToPrefix(lib.Pkg) + "."
	ndef := r.NSym()
	nhashed64def := r.NHashed64def()
	nhasheddef := r.NHasheddef()
	or := &oReader{
		Reader:       r,
		unit:         unit,
		version:      localSymVersion,
		flags:        r.Flags(),
		pkgprefix:    pkgprefix,
		syms:         make([]Sym, ndef+nhashed64def+nhasheddef+r.NNonpkgdef()+r.NNonpkgref()),
		ndef:         ndef,
		nhasheddef:   nhasheddef,
		nhashed64def: nhashed64def,
		objidx:       uint32(len(l.objs)),
	}

	// Autolib
	lib.Autolib = append(lib.Autolib, r.Autolib()...)

	// DWARF file table
	nfile := r.NDwarfFile()
	unit.DWARFFileTable = make([]string, nfile)
	for i := range unit.DWARFFileTable {
		unit.DWARFFileTable[i] = r.DwarfFile(i)
	}

	l.addObj(lib.Pkg, or)
	l.preloadSyms(or, pkgDef)

	// The caller expects us consuming all the data
	f.MustSeek(length, os.SEEK_CUR)

	return r.Fingerprint()
}

// Preload symbols of given kind from an object.
func (l *Loader) preloadSyms(r *oReader, kind int) {
	var start, end uint32
	switch kind {
	case pkgDef:
		start = 0
		end = uint32(r.ndef)
	case hashed64Def:
		start = uint32(r.ndef)
		end = uint32(r.ndef + r.nhashed64def)
	case hashedDef:
		start = uint32(r.ndef + r.nhashed64def)
		end = uint32(r.ndef + r.nhashed64def + r.nhasheddef)
		if l.hasUnknownPkgPath {
			// The content hash depends on symbol name expansion. If any package is
			// built without fully expanded names, the content hash is unreliable.
			// Treat them as named symbols.
			// This is rare.
			// (We don't need to do this for hashed64Def case, as there the hash
			// function is simply the identity function, which doesn't depend on
			// name expansion.)
			kind = nonPkgDef
		}
	case nonPkgDef:
		start = uint32(r.ndef + r.nhashed64def + r.nhasheddef)
		end = uint32(r.ndef + r.nhashed64def + r.nhasheddef + r.NNonpkgdef())
	default:
		panic("preloadSyms: bad kind")
	}
	l.growAttrBitmaps(len(l.objSyms) + int(end-start))
	needNameExpansion := r.NeedNameExpansion()
	loadingRuntimePkg := r.unit.Lib.Pkg == "runtime"
	for i := start; i < end; i++ {
		osym := r.Sym(i)
		var name string
		var v int
		if kind != hashed64Def && kind != hashedDef { // we don't need the name, etc. for hashed symbols
			name = osym.Name(r.Reader)
			if needNameExpansion {
				name = strings.Replace(name, "\"\".", r.pkgprefix, -1)
			}
			v = abiToVer(osym.ABI(), r.version)
		}
		gi := l.addSym(name, v, r, i, kind, osym)
		r.syms[i] = gi
		if osym.TopFrame() {
			l.SetAttrTopFrame(gi, true)
		}
		if osym.Local() {
			l.SetAttrLocal(gi, true)
		}
		if osym.UsedInIface() {
			l.SetAttrUsedInIface(gi, true)
		}
		if strings.HasPrefix(name, "go.itablink.") {
			l.itablink[gi] = struct{}{}
		}
		if strings.HasPrefix(name, "runtime.") ||
			(loadingRuntimePkg && strings.HasPrefix(name, "type.")) {
			if bi := goobj2.BuiltinIdx(name, v); bi != -1 {
				// This is a definition of a builtin symbol. Record where it is.
				l.builtinSyms[bi] = gi
			}
		}
		if a := int32(osym.Align()); a != 0 && a > l.SymAlign(gi) {
			l.SetSymAlign(gi, a)
		}
	}
}

// Add hashed (content-addressable) symbols, non-package symbols, and
// references to external symbols (which are always named).
func (l *Loader) LoadNonpkgSyms(arch *sys.Arch) {
	l.npkgsyms = l.NSym()
	for _, o := range l.objs[goObjStart:] {
		l.preloadSyms(o.r, hashed64Def)
		l.preloadSyms(o.r, hashedDef)
		l.preloadSyms(o.r, nonPkgDef)
	}
	for _, o := range l.objs[goObjStart:] {
		loadObjRefs(l, o.r, arch)
	}
	l.values = make([]int64, l.NSym(), l.NSym()+1000) // +1000 make some room for external symbols
}

func loadObjRefs(l *Loader, r *oReader, arch *sys.Arch) {
	// load non-package refs
	ndef := uint32(r.NAlldef())
	needNameExpansion := r.NeedNameExpansion()
	for i, n := uint32(0), uint32(r.NNonpkgref()); i < n; i++ {
		osym := r.Sym(ndef + i)
		name := osym.Name(r.Reader)
		if needNameExpansion {
			name = strings.Replace(name, "\"\".", r.pkgprefix, -1)
		}
		v := abiToVer(osym.ABI(), r.version)
		r.syms[ndef+i] = l.LookupOrCreateSym(name, v)
		gi := r.syms[ndef+i]
		if osym.Local() {
			l.SetAttrLocal(gi, true)
		}
		if osym.UsedInIface() {
			l.SetAttrUsedInIface(gi, true)
		}
	}

	// load flags of package refs
	for i, n := 0, r.NRefFlags(); i < n; i++ {
		rf := r.RefFlags(i)
		gi := l.resolve(r, rf.Sym())
		if rf.Flag2()&goobj2.SymFlagUsedInIface != 0 {
			l.SetAttrUsedInIface(gi, true)
		}
	}
}

func abiToVer(abi uint16, localSymVersion int) int {
	var v int
	if abi == goobj2.SymABIstatic {
		// Static
		v = localSymVersion
	} else if abiver := sym.ABIToVersion(obj.ABI(abi)); abiver != -1 {
		// Note that data symbols are "ABI0", which maps to version 0.
		v = abiver
	} else {
		log.Fatalf("invalid symbol ABI: %d", abi)
	}
	return v
}

// ResolveABIAlias given a symbol returns the ABI alias target of that
// symbol. If the sym in question is not an alias, the sym itself is
// returned.
func (l *Loader) ResolveABIAlias(s Sym) Sym {
	if s == 0 {
		return 0
	}
	if l.SymType(s) != sym.SABIALIAS {
		return s
	}
	relocs := l.Relocs(s)
	target := relocs.At(0).Sym()
	if l.SymType(target) == sym.SABIALIAS {
		panic(fmt.Sprintf("ABI alias %s references another ABI alias %s", l.SymName(s), l.SymName(target)))
	}
	return target
}

// TopLevelSym tests a symbol (by name and kind) to determine whether
// the symbol first class sym (participating in the link) or is an
// anonymous aux or sub-symbol containing some sub-part or payload of
// another symbol.
func (l *Loader) TopLevelSym(s Sym) bool {
	return topLevelSym(l.RawSymName(s), l.SymType(s))
}

// topLevelSym tests a symbol name and kind to determine whether
// the symbol first class sym (participating in the link) or is an
// anonymous aux or sub-symbol containing some sub-part or payload of
// another symbol.
func topLevelSym(sname string, skind sym.SymKind) bool {
	if sname != "" {
		return true
	}
	switch skind {
	case sym.SDWARFFCN, sym.SDWARFABSFCN, sym.SDWARFTYPE, sym.SDWARFCONST, sym.SDWARFCUINFO, sym.SDWARFRANGE, sym.SDWARFLOC, sym.SDWARFLINES, sym.SGOFUNC:
		return true
	default:
		return false
	}
}

// cloneToExternal takes the existing object file symbol (symIdx)
// and creates a new external symbol payload that is a clone with
// respect to name, version, type, relocations, etc. The idea here
// is that if the linker decides it wants to update the contents of
// a symbol originally discovered as part of an object file, it's
// easier to do this if we make the updates to an external symbol
// payload.
func (l *Loader) cloneToExternal(symIdx Sym) {
	if l.IsExternal(symIdx) {
		panic("sym is already external, no need for clone")
	}

	// Read the particulars from object.
	r, li := l.toLocal(symIdx)
	osym := r.Sym(li)
	sname := osym.Name(r.Reader)
	if r.NeedNameExpansion() {
		sname = strings.Replace(sname, "\"\".", r.pkgprefix, -1)
	}
	sver := abiToVer(osym.ABI(), r.version)
	skind := sym.AbiSymKindToSymKind[objabi.SymKind(osym.Type())]

	// Create new symbol, update version and kind.
	pi := l.newPayload(sname, sver)
	pp := l.payloads[pi]
	pp.kind = skind
	pp.ver = sver
	pp.size = int64(osym.Siz())
	pp.objidx = r.objidx

	// If this is a def, then copy the guts. We expect this case
	// to be very rare (one case it may come up is with -X).
	if li < uint32(r.NAlldef()) {

		// Copy relocations
		relocs := l.Relocs(symIdx)
		pp.relocs = make([]goobj2.Reloc, relocs.Count())
		pp.reltypes = make([]objabi.RelocType, relocs.Count())
		for i := range pp.relocs {
			// Copy the relocs slice.
			// Convert local reference to global reference.
			rel := relocs.At(i)
			pp.relocs[i].Set(rel.Off(), rel.Siz(), 0, rel.Add(), goobj2.SymRef{PkgIdx: 0, SymIdx: uint32(rel.Sym())})
			pp.reltypes[i] = rel.Type()
		}

		// Copy data
		pp.data = r.Data(li)
	}

	// If we're overriding a data symbol, collect the associated
	// Gotype, so as to propagate it to the new symbol.
	auxs := r.Auxs(li)
	pp.auxs = auxs

	// Install new payload to global index space.
	// (This needs to happen at the end, as the accessors above
	// need to access the old symbol content.)
	l.objSyms[symIdx] = objSym{l.extReader.objidx, uint32(pi)}
	l.extReader.syms = append(l.extReader.syms, symIdx)
}

// Copy the payload of symbol src to dst. Both src and dst must be external
// symbols.
// The intended use case is that when building/linking against a shared library,
// where we do symbol name mangling, the Go object file may have reference to
// the original symbol name whereas the shared library provides a symbol with
// the mangled name. When we do mangling, we copy payload of mangled to original.
func (l *Loader) CopySym(src, dst Sym) {
	if !l.IsExternal(dst) {
		panic("dst is not external") //l.newExtSym(l.SymName(dst), l.SymVersion(dst))
	}
	if !l.IsExternal(src) {
		panic("src is not external") //l.cloneToExternal(src)
	}
	l.payloads[l.extIndex(dst)] = l.payloads[l.extIndex(src)]
	l.SetSymPkg(dst, l.SymPkg(src))
	// TODO: other attributes?
}

// CopyAttributes copies over all of the attributes of symbol 'src' to
// symbol 'dst'.
func (l *Loader) CopyAttributes(src Sym, dst Sym) {
	l.SetAttrReachable(dst, l.AttrReachable(src))
	l.SetAttrOnList(dst, l.AttrOnList(src))
	l.SetAttrLocal(dst, l.AttrLocal(src))
	l.SetAttrNotInSymbolTable(dst, l.AttrNotInSymbolTable(src))
	if l.IsExternal(dst) {
		l.SetAttrVisibilityHidden(dst, l.AttrVisibilityHidden(src))
		l.SetAttrDuplicateOK(dst, l.AttrDuplicateOK(src))
		l.SetAttrShared(dst, l.AttrShared(src))
		l.SetAttrExternal(dst, l.AttrExternal(src))
	} else {
		// Some attributes are modifiable only for external symbols.
		// In such cases, don't try to transfer over the attribute
		// from the source even if there is a clash. This comes up
		// when copying attributes from a dupOK ABI wrapper symbol to
		// the real target symbol (which may not be marked dupOK).
	}
	l.SetAttrTopFrame(dst, l.AttrTopFrame(src))
	l.SetAttrSpecial(dst, l.AttrSpecial(src))
	l.SetAttrCgoExportDynamic(dst, l.AttrCgoExportDynamic(src))
	l.SetAttrCgoExportStatic(dst, l.AttrCgoExportStatic(src))
	l.SetAttrReadOnly(dst, l.AttrReadOnly(src))
}

// CreateExtSym creates a new external symbol with the specified name
// without adding it to any lookup tables, returning a Sym index for it.
func (l *Loader) CreateExtSym(name string, ver int) Sym {
	return l.newExtSym(name, ver)
}

// CreateStaticSym creates a new static symbol with the specified name
// without adding it to any lookup tables, returning a Sym index for it.
func (l *Loader) CreateStaticSym(name string) Sym {
	// Assign a new unique negative version -- this is to mark the
	// symbol so that it is not included in the name lookup table.
	l.anonVersion--
	return l.newExtSym(name, l.anonVersion)
}

func (l *Loader) FreeSym(i Sym) {
	if l.IsExternal(i) {
		pp := l.getPayload(i)
		*pp = extSymPayload{}
	}
}

// relocId is essentially a <S,R> tuple identifying the Rth
// relocation of symbol S.
type relocId struct {
	sym  Sym
	ridx int
}

// SetRelocVariant sets the 'variant' property of a relocation on
// some specific symbol.
func (l *Loader) SetRelocVariant(s Sym, ri int, v sym.RelocVariant) {
	// sanity check
	if relocs := l.Relocs(s); ri >= relocs.Count() {
		panic("invalid relocation ID")
	}
	if l.relocVariant == nil {
		l.relocVariant = make(map[relocId]sym.RelocVariant)
	}
	if v != 0 {
		l.relocVariant[relocId{s, ri}] = v
	} else {
		delete(l.relocVariant, relocId{s, ri})
	}
}

// RelocVariant returns the 'variant' property of a relocation on
// some specific symbol.
func (l *Loader) RelocVariant(s Sym, ri int) sym.RelocVariant {
	return l.relocVariant[relocId{s, ri}]
}

// UndefinedRelocTargets iterates through the global symbol index
// space, looking for symbols with relocations targeting undefined
// references. The linker's loadlib method uses this to determine if
// there are unresolved references to functions in system libraries
// (for example, libgcc.a), presumably due to CGO code. Return
// value is a list of loader.Sym's corresponding to the undefined
// cross-refs. The "limit" param controls the maximum number of
// results returned; if "limit" is -1, then all undefs are returned.
func (l *Loader) UndefinedRelocTargets(limit int) []Sym {
	result := []Sym{}
	for si := Sym(1); si < Sym(len(l.objSyms)); si++ {
		relocs := l.Relocs(si)
		for ri := 0; ri < relocs.Count(); ri++ {
			r := relocs.At(ri)
			rs := r.Sym()
			if rs != 0 && l.SymType(rs) == sym.SXREF && l.RawSymName(rs) != ".got" {
				result = append(result, rs)
				if limit != -1 && len(result) >= limit {
					break
				}
			}
		}
	}
	return result
}

// AssignTextSymbolOrder populates the Textp slices within each
// library and compilation unit, insuring that packages are laid down
// in dependency order (internal first, then everything else). Return value
// is a slice of all text syms.
func (l *Loader) AssignTextSymbolOrder(libs []*sym.Library, intlibs []bool, extsyms []Sym) []Sym {

	// Library Textp lists should be empty at this point.
	for _, lib := range libs {
		if len(lib.Textp) != 0 {
			panic("expected empty Textp slice for library")
		}
		if len(lib.DupTextSyms) != 0 {
			panic("expected empty DupTextSyms slice for library")
		}
	}

	// Used to record which dupok symbol we've assigned to a unit.
	// Can't use the onlist attribute here because it will need to
	// clear for the later assignment of the sym.Symbol to a unit.
	// NB: we can convert to using onList once we no longer have to
	// call the regular addToTextp.
	assignedToUnit := MakeBitmap(l.NSym() + 1)

	// Start off textp with reachable external syms.
	textp := []Sym{}
	for _, sym := range extsyms {
		if !l.attrReachable.Has(sym) {
			continue
		}
		textp = append(textp, sym)
	}

	// Walk through all text symbols from Go object files and append
	// them to their corresponding library's textp list.
	for _, o := range l.objs[goObjStart:] {
		r := o.r
		lib := r.unit.Lib
		for i, n := uint32(0), uint32(r.NAlldef()); i < n; i++ {
			gi := l.toGlobal(r, i)
			if !l.attrReachable.Has(gi) {
				continue
			}
			osym := r.Sym(i)
			st := sym.AbiSymKindToSymKind[objabi.SymKind(osym.Type())]
			if st != sym.STEXT {
				continue
			}
			dupok := osym.Dupok()
			if r2, i2 := l.toLocal(gi); r2 != r || i2 != i {
				// A dupok text symbol is resolved to another package.
				// We still need to record its presence in the current
				// package, as the trampoline pass expects packages
				// are laid out in dependency order.
				lib.DupTextSyms = append(lib.DupTextSyms, sym.LoaderSym(gi))
				continue // symbol in different object
			}
			if dupok {
				lib.DupTextSyms = append(lib.DupTextSyms, sym.LoaderSym(gi))
				continue
			}

			lib.Textp = append(lib.Textp, sym.LoaderSym(gi))
		}
	}

	// Now assemble global textp, and assign text symbols to units.
	for _, doInternal := range [2]bool{true, false} {
		for idx, lib := range libs {
			if intlibs[idx] != doInternal {
				continue
			}
			lists := [2][]sym.LoaderSym{lib.Textp, lib.DupTextSyms}
			for i, list := range lists {
				for _, s := range list {
					sym := Sym(s)
					if l.attrReachable.Has(sym) && !assignedToUnit.Has(sym) {
						textp = append(textp, sym)
						unit := l.SymUnit(sym)
						if unit != nil {
							unit.Textp = append(unit.Textp, s)
							assignedToUnit.Set(sym)
						}
						// Dupok symbols may be defined in multiple packages; the
						// associated package for a dupok sym is chosen sort of
						// arbitrarily (the first containing package that the linker
						// loads). Canonicalizes its Pkg to the package with which
						// it will be laid down in text.
						if i == 1 /* DupTextSyms2 */ && l.SymPkg(sym) != lib.Pkg {
							l.SetSymPkg(sym, lib.Pkg)
						}
					}
				}
			}
			lib.Textp = nil
			lib.DupTextSyms = nil
		}
	}

	return textp
}

// ErrorReporter is a helper class for reporting errors.
type ErrorReporter struct {
	ldr              *Loader
	AfterErrorAction func()
}

// Errorf method logs an error message.
//
// After each error, the error actions function will be invoked; this
// will either terminate the link immediately (if -h option given)
// or it will keep a count and exit if more than 20 errors have been printed.
//
// Logging an error means that on exit cmd/link will delete any
// output file and return a non-zero error code.
//
func (reporter *ErrorReporter) Errorf(s Sym, format string, args ...interface{}) {
	if s != 0 && reporter.ldr.SymName(s) != "" {
		format = reporter.ldr.SymName(s) + ": " + format
	} else {
		format = fmt.Sprintf("sym %d: %s", s, format)
	}
	format += "\n"
	fmt.Fprintf(os.Stderr, format, args...)
	reporter.AfterErrorAction()
}

// GetErrorReporter returns the loader's associated error reporter.
func (l *Loader) GetErrorReporter() *ErrorReporter {
	return l.errorReporter
}

// Errorf method logs an error message. See ErrorReporter.Errorf for details.
func (l *Loader) Errorf(s Sym, format string, args ...interface{}) {
	l.errorReporter.Errorf(s, format, args...)
}

// Symbol statistics.
func (l *Loader) Stat() string {
	s := fmt.Sprintf("%d symbols, %d reachable\n", l.NSym(), l.NReachableSym())
	s += fmt.Sprintf("\t%d package symbols, %d hashed symbols, %d non-package symbols, %d external symbols\n",
		l.npkgsyms, len(l.hashed64Syms)+len(l.hashedSyms), int(l.extStart)-l.npkgsyms-len(l.hashed64Syms)-len(l.hashedSyms), l.NSym()-int(l.extStart))
	return s
}

// For debugging.
func (l *Loader) Dump() {
	fmt.Println("objs")
	for _, obj := range l.objs[goObjStart:] {
		if obj.r != nil {
			fmt.Println(obj.i, obj.r.unit.Lib)
		}
	}
	fmt.Println("extStart:", l.extStart)
	fmt.Println("Nsyms:", len(l.objSyms))
	fmt.Println("syms")
	for i := Sym(1); i < Sym(len(l.objSyms)); i++ {
		pi := interface{}("")
		if l.IsExternal(i) {
			pi = fmt.Sprintf("<ext %d>", l.extIndex(i))
		}
		fmt.Println(i, l.SymName(i), l.SymType(i), pi)
	}
	fmt.Println("symsByName")
	for name, i := range l.symsByName[0] {
		fmt.Println(i, name, 0)
	}
	for name, i := range l.symsByName[1] {
		fmt.Println(i, name, 1)
	}
	fmt.Println("payloads:")
	for i := range l.payloads {
		pp := l.payloads[i]
		fmt.Println(i, pp.name, pp.ver, pp.kind)
	}
}
