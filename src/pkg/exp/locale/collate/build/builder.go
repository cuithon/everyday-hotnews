// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import (
	"exp/locale/collate"
	"exp/norm"
	"fmt"
	"io"
	"log"
	"sort"
	"strings"
	"unicode/utf8"
)

// TODO: optimizations:
// - expandElem is currently 20K. By putting unique colElems in a separate
//   table and having a byte array of indexes into this table, we can reduce
//   the total size to about 7K. By also factoring out the length bytes, we
//   can reduce this to about 6K.
// - trie valueBlocks are currently 100K. There are a lot of sparse blocks
//   and many consecutive values with the same stride. This can be further
//   compacted.
// - Compress secondary weights into 8 bits.
// - Some LDML specs specify a context element. Currently we simply concatenate
//   those.  Context can be implemented using the contraction trie. If Builder 
//   could analyze and detect when using a context makes sense, there is no
//   need to expose this construct in the API.

// entry is used to keep track of a single entry in the collation element table
// during building. Examples of entries can be found in the Default Unicode
// Collation Element Table.
// See http://www.unicode.org/Public/UCA/6.0.0/allkeys.txt.
type entry struct {
	runes []rune
	elems [][]int // the collation elements for runes
	str   string  // same as string(runes)

	// prev, next, and level are used to keep track of tailorings.
	prev, next *entry
	level      collate.Level // next differs at this level

	decompose bool // can use NFKD decomposition to generate elems
	exclude   bool // do not include in table
	logical   logicalAnchor

	expansionIndex    int // used to store index into expansion table
	contractionHandle ctHandle
	contractionIndex  int // index into contraction elements
}

func (e *entry) String() string {
	return fmt.Sprintf("%X -> %X (ch:%x; ci:%d, ei:%d)",
		e.runes, e.elems, e.contractionHandle, e.contractionIndex, e.expansionIndex)
}

func (e *entry) skip() bool {
	return e.contraction()
}

func (e *entry) expansion() bool {
	return !e.decompose && len(e.elems) > 1
}

func (e *entry) contraction() bool {
	return len(e.runes) > 1
}

func (e *entry) contractionStarter() bool {
	return e.contractionHandle.n != 0
}

// A Builder builds a root collation table.  The user must specify the
// collation elements for each entry.  A common use will be to base the weights
// on those specified in the allkeys* file as provided by the UCA or CLDR.
type Builder struct {
	index  *trieBuilder
	root   ordering
	locale []*Tailoring
	t      *table
	err    error
	built  bool

	minNonVar int // lowest primary recorded for a variable
	varTop    int // highest primary recorded for a non-variable
}

// A Tailoring builds a collation table based on another collation table.
// The table is defined by specifying tailorings to the underlying table.
// See http://unicode.org/reports/tr35/ for an overview of tailoring
// collation tables.  The CLDR contains pre-defined tailorings for a variety
// of languages (See http://www.unicode.org/Public/cldr/2.0.1/core.zip.)
type Tailoring struct {
	id string
	// TODO: implement.
}

// NewBuilder returns a new Builder.
func NewBuilder() *Builder {
	b := &Builder{
		index: newTrieBuilder(),
		root:  makeRootOrdering(),
	}
	return b
}

// Tailoring returns a Tailoring for the given locale.  One should 
// have completed all calls to Add before calling Tailoring.
func (b *Builder) Tailoring(locale string) *Tailoring {
	t := &Tailoring{
		id: locale,
	}
	b.locale = append(b.locale, t)
	return t
}

// Add adds an entry to the collation element table, mapping 
// a slice of runes to a sequence of collation elements.
// A collation element is specified as list of weights: []int{primary, secondary, ...}.
// The entries are typically obtained from a collation element table
// as defined in http://www.unicode.org/reports/tr10/#Data_Table_Format.
// Note that the collation elements specified by colelems are only used
// as a guide.  The actual weights generated by Builder may differ.
// The argument variables is a list of indices into colelems that should contain
// a value for each colelem that is a variable. (See the reference above.)
func (b *Builder) Add(runes []rune, colelems [][]int, variables []int) error {
	str := string(runes)

	elems := make([][]int, len(colelems))
	for i, ce := range colelems {
		elems[i] = append(elems[i], ce...)
		if len(ce) == 0 {
			elems[i] = append(elems[i], []int{0, 0, 0, 0}...)
			break
		}
		if len(ce) == 1 {
			elems[i] = append(elems[i], defaultSecondary)
		}
		if len(ce) <= 2 {
			elems[i] = append(elems[i], defaultTertiary)
		}
		if len(ce) <= 3 {
			elems[i] = append(elems[i], ce[0])
		}
	}
	for i, ce := range elems {
		isvar := false
		for _, j := range variables {
			if i == j {
				isvar = true
			}
		}
		if isvar {
			if ce[0] >= b.minNonVar && b.minNonVar > 0 {
				return fmt.Errorf("primary value %X of variable is larger than the smallest non-variable %X", ce[0], b.minNonVar)
			}
			if ce[0] > b.varTop {
				b.varTop = ce[0]
			}
		} else if ce[0] > 0 {
			if ce[0] <= b.varTop {
				return fmt.Errorf("primary value %X of non-variable is smaller than the highest variable %X", ce[0], b.varTop)
			}
			if b.minNonVar == 0 || ce[0] < b.minNonVar {
				b.minNonVar = ce[0]
			}
		}
	}
	elems, err := convertLargeWeights(elems)
	if err != nil {
		return err
	}
	b.root.newEntry(str, elems)
	return nil
}

// SetAnchor sets the point after which elements passed in subsequent calls to
// Insert will be inserted.  It is equivalent to the reset directive in an LDML
// specification.  See Insert for an example.
// SetAnchor supports the following logical reset positions:
// <first_tertiary_ignorable/>, <last_teriary_ignorable/>, <first_primary_ignorable/>,
// and <last_non_ignorable/>.
func (t *Tailoring) SetAnchor(anchor string) error {
	// TODO: implement.
	return nil
}

// SetAnchorBefore is similar to SetAnchor, except that subsequent calls to
// Insert will insert entries before the anchor.
func (t *Tailoring) SetAnchorBefore(anchor string) error {
	// TODO: implement.
	return nil
}

// Insert sets the ordering of str relative to the entry set by the previous
// call to SetAnchor or Insert.  The argument extend corresponds
// to the extend elements as defined in LDML.  A non-empty value for extend
// will cause the collation elements corresponding to extend to be appended
// to the collation elements generated for the entry added by Insert.
// This has the same net effect as sorting str after the string anchor+extend.
// See http://www.unicode.org/reports/tr10/#Tailoring_Example for details
// on parametric tailoring and http://unicode.org/reports/tr35/#Collation_Elements
// for full details on LDML.
// 
// Examples: create a tailoring for Swedish, where "ä" is ordered after "z"
// at the primary sorting level:
//      t := b.Tailoring("se")
// 		t.SetAnchor("z")
// 		t.Insert(collate.Primary, "ä", "")
// Order "ü" after "ue" at the secondary sorting level:
//		t.SetAnchor("ue")
//		t.Insert(collate.Secondary, "ü","")
// or
//		t.SetAnchor("u")
//		t.Insert(collate.Secondary, "ü", "e")
// Order "q" afer "ab" at the secondary level and "Q" after "q"
// at the tertiary level:
// 		t.SetAnchor("ab")
// 		t.Insert(collate.Secondary, "q", "")
// 		t.Insert(collate.Tertiary, "Q", "")
// Order "b" before "a":
//      t.SetAnchorBefore("a")
//      t.Insert(collate.Primary, "b", "")
// Order "0" after the last primary ignorable:
//      t.SetAnchor("<last_primary_ignorable/>")
//      t.Insert(collate.Primary, "0", "")
func (t *Tailoring) Insert(level collate.Level, str, extend string) error {
	// TODO: implement.
	return nil
}

func (b *Builder) baseColElem(e *entry) uint32 {
	ce := uint32(0)
	var err error
	switch {
	case e.expansion():
		ce, err = makeExpandIndex(e.expansionIndex)
	default:
		if e.decompose {
			log.Fatal("decompose should be handled elsewhere")
		}
		ce, err = makeCE(e.elems[0])
	}
	if err != nil {
		b.error(fmt.Errorf("%s: %X -> %X", err, e.runes, e.elems))
	}
	return ce
}

func (b *Builder) colElem(e *entry) uint32 {
	if e.skip() {
		log.Fatal("cannot build colElem for entry that should be skipped")
	}
	ce := uint32(0)
	var err error
	switch {
	case e.decompose:
		t1 := e.elems[0][2]
		t2 := 0
		if len(e.elems) > 1 {
			t2 = e.elems[1][2]
		}
		ce, err = makeDecompose(t1, t2)
	case e.contractionStarter():
		ce, err = makeContractIndex(e.contractionHandle, e.contractionIndex)
	default:
		if len(e.runes) > 1 {
			log.Fatal("colElem: contractions are handled in contraction trie")
		}
		ce = b.baseColElem(e)
	}
	if err != nil {
		b.error(err)
	}
	return ce
}

func (b *Builder) error(e error) {
	if e != nil {
		b.err = e
	}
}

func (b *Builder) build() (*table, error) {
	if !b.built {
		b.built = true
		b.t = &table{
			maxContractLen: utf8.UTFMax,
			variableTop:    uint32(b.varTop),
		}

		b.root.sort()
		b.simplify()
		b.processExpansions()   // requires simplify
		b.processContractions() // requires simplify
		b.buildTrie()           // requires process*
	}
	if b.err != nil {
		return nil, b.err
	}
	return b.t, nil
}

// Build builds the root Collator.
func (b *Builder) Build() (*collate.Collator, error) {
	t, err := b.build()
	if err != nil {
		return nil, err
	}
	return collate.Init(t), nil
}

// Build builds a Collator for Tailoring t.
func (t *Tailoring) Build() (*collate.Collator, error) {
	// TODO: implement.
	return nil, nil
}

// Print prints the tables for b and all its Tailorings as a Go file
// that can be included in the Collate package.
func (b *Builder) Print(w io.Writer) (int, error) {
	t, err := b.build()
	if err != nil {
		return 0, err
	}
	// TODO: support multiple locales
	n, _, err := t.fprint(w, "root")
	return n, err
}

// reproducibleFromNFKD checks whether the given expansion could be generated
// from an NFKD expansion.
func reproducibleFromNFKD(e *entry, exp, nfkd [][]int) bool {
	// Length must be equal.
	if len(exp) != len(nfkd) {
		return false
	}
	for i, ce := range exp {
		// Primary and secondary values should be equal.
		if ce[0] != nfkd[i][0] || ce[1] != nfkd[i][1] {
			return false
		}
		// Tertiary values should be equal to maxTertiary for third element onwards.
		// TODO: there seem to be a lot of cases in CLDR (e.g. ㏭ in zh.xml) that can
		// simply be dropped.  Try this out by dropping the following code.
		if i >= 2 && ce[2] != maxTertiary {
			return false
		}
	}
	return true
}

func equalCE(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < 3; i++ {
		if b[i] != a[i] {
			return false
		}
	}
	return true
}

func equalCEArrays(a, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !equalCE(a[i], b[i]) {
			return false
		}
	}
	return true
}

func (b *Builder) simplify() {
	// Runes that are a starter of a contraction should not be removed.
	// (To date, there is only Kannada character 0CCA.)
	keep := make(map[rune]bool)
	o := b.root

	for e := o.front(); e != nil; e, _ = e.nextIndexed() {
		if len(e.runes) > 1 {
			keep[e.runes[0]] = true
		}
	}
	// Remove entries for which the runes normalize (using NFD) to identical values.
	for e := o.front(); e != nil; e, _ = e.nextIndexed() {
		s := e.str
		nfd := norm.NFD.String(s)
		if len(e.runes) > 1 || keep[e.runes[0]] || nfd == s {
			continue
		}
		if equalCEArrays(o.genColElems(nfd), e.elems) {
			e.remove()
		}
	}

	// Tag entries for which the runes NFKD decompose to identical values.
	for e := o.front(); e != nil; e, _ = e.nextIndexed() {
		s := e.str
		nfkd := norm.NFKD.String(s)
		if len(e.runes) > 1 || keep[e.runes[0]] || nfkd == s {
			continue
		}
		if reproducibleFromNFKD(e, e.elems, o.genColElems(nfkd)) {
			e.decompose = true
		}
	}
}

// convertLargeWeights converts collation elements with large 
// primaries (either double primaries or for illegal runes)
// to our own representation.
// A CJK character C is represented in the DUCET as
//   [.FBxx.0020.0002.C][.BBBB.0000.0000.C]
// We will rewrite these characters to a single CE.
// We assume the CJK values start at 0x8000.
// See http://unicode.org/reports/tr10/#Implicit_Weights
func convertLargeWeights(elems [][]int) (res [][]int, err error) {
	const (
		cjkPrimaryStart   = 0xFB40
		rarePrimaryStart  = 0xFB80
		otherPrimaryStart = 0xFBC0
		illegalPrimary    = 0xFFFE
		highBitsMask      = 0x3F
		lowBitsMask       = 0x7FFF
		lowBitsFlag       = 0x8000
		shiftBits         = 15
	)
	for i := 0; i < len(elems); i++ {
		ce := elems[i]
		p := ce[0]
		if p < cjkPrimaryStart {
			continue
		}
		if p > 0xFFFF {
			return elems, fmt.Errorf("found primary weight %X; should be <= 0xFFFF", p)
		}
		if p >= illegalPrimary {
			ce[0] = illegalOffset + p - illegalPrimary
		} else {
			if i+1 >= len(elems) {
				return elems, fmt.Errorf("second part of double primary weight missing: %v", elems)
			}
			if elems[i+1][0]&lowBitsFlag == 0 {
				return elems, fmt.Errorf("malformed second part of double primary weight: %v", elems)
			}
			np := ((p & highBitsMask) << shiftBits) + elems[i+1][0]&lowBitsMask
			switch {
			case p < rarePrimaryStart:
				np += commonUnifiedOffset
			case p < otherPrimaryStart:
				np += rareUnifiedOffset
			default:
				p += otherOffset
			}
			ce[0] = np
			for j := i + 1; j+1 < len(elems); j++ {
				elems[j] = elems[j+1]
			}
			elems = elems[:len(elems)-1]
		}
	}
	return elems, nil
}

// appendExpansion converts the given collation sequence to
// collation elements and adds them to the expansion table.
// It returns an index to the expansion table.
func (b *Builder) appendExpansion(e *entry) int {
	t := b.t
	i := len(t.expandElem)
	ce := uint32(len(e.elems))
	t.expandElem = append(t.expandElem, ce)
	for _, w := range e.elems {
		ce, err := makeCE(w)
		if err != nil {
			b.error(err)
			return -1
		}
		t.expandElem = append(t.expandElem, ce)
	}
	return i
}

// processExpansions extracts data necessary to generate
// the extraction tables.
func (b *Builder) processExpansions() {
	eidx := make(map[string]int)
	o := b.root
	for e := o.front(); e != nil; e, _ = e.nextIndexed() {
		if !e.expansion() {
			continue
		}
		key := fmt.Sprintf("%v", e.elems)
		i, ok := eidx[key]
		if !ok {
			i = b.appendExpansion(e)
			eidx[key] = i
		}
		e.expansionIndex = i
	}
}

func (b *Builder) processContractions() {
	// Collate contractions per starter rune.
	starters := []rune{}
	o := b.root
	cm := make(map[rune][]*entry)
	for e := o.front(); e != nil; e, _ = e.nextIndexed() {
		if e.contraction() {
			if len(e.str) > b.t.maxContractLen {
				b.t.maxContractLen = len(e.str)
			}
			r := e.runes[0]
			if _, ok := cm[r]; !ok {
				starters = append(starters, r)
			}
			cm[r] = append(cm[r], e)
		}
	}
	// Add entries of single runes that are at a start of a contraction.
	for e := o.front(); e != nil; e, _ = e.nextIndexed() {
		if !e.contraction() {
			r := e.runes[0]
			if _, ok := cm[r]; ok {
				cm[r] = append(cm[r], e)
			}
		}
	}
	// Build the tries for the contractions.
	t := b.t
	handlemap := make(map[string]ctHandle)
	for _, r := range starters {
		l := cm[r]
		// Compute suffix strings. There are 31 different contraction suffix
		// sets for 715 contractions and 82 contraction starter runes as of
		// version 6.0.0.
		sufx := []string{}
		hasSingle := false
		for _, e := range l {
			if len(e.runes) > 1 {
				sufx = append(sufx, string(e.runes[1:]))
			} else {
				hasSingle = true
			}
		}
		if !hasSingle {
			b.error(fmt.Errorf("no single entry for starter rune %U found", r))
			continue
		}
		// Unique the suffix set.
		sort.Strings(sufx)
		key := strings.Join(sufx, "\n")
		handle, ok := handlemap[key]
		if !ok {
			var err error
			handle, err = t.contractTries.appendTrie(sufx)
			if err != nil {
				b.error(err)
			}
			handlemap[key] = handle
		}
		// Bucket sort entries in index order.
		es := make([]*entry, len(l))
		for _, e := range l {
			var o, sn int
			if len(e.runes) > 1 {
				str := []byte(string(e.runes[1:]))
				o, sn = t.contractTries.lookup(handle, str)
				if sn != len(str) {
					log.Fatalf("processContractions: unexpected length for '%X'; len=%d; want %d", e.runes, sn, len(str))
				}
			}
			if es[o] != nil {
				log.Fatalf("Multiple contractions for position %d for rune %U", o, e.runes[0])
			}
			es[o] = e
		}
		// Store info in entry for starter rune.
		es[0].contractionIndex = len(t.contractElem)
		es[0].contractionHandle = handle
		// Add collation elements for contractions.
		for _, e := range es {
			t.contractElem = append(t.contractElem, b.baseColElem(e))
		}
	}
}

func (b *Builder) buildTrie() {
	t := newNode()
	o := b.root
	for e := o.front(); e != nil; e, _ = e.nextIndexed() {
		if !e.skip() {
			ce := b.colElem(e)
			t.insert(e.runes[0], ce)
		}
	}
	b.t.root = b.index.addTrie(t)
	i, err := b.index.generate()
	b.t.index = *i
	b.error(err)
}
