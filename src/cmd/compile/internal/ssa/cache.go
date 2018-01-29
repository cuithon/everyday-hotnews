// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssa

import (
	"cmd/internal/obj"
	"sort"
)

// A Cache holds reusable compiler state.
// It is intended to be re-used for multiple Func compilations.
type Cache struct {
	// Storage for low-numbered values and blocks.
	values [2000]Value
	blocks [200]Block
	locs   [2000]Location

	// Reusable stackAllocState.
	// See stackalloc.go's {new,put}StackAllocState.
	stackAllocState *stackAllocState

	domblockstore []ID         // scratch space for computing dominators
	scrSparse     []*sparseSet // scratch sparse sets to be re-used.

	ValueToProgAfter []*obj.Prog
	blockDebug       []BlockDebug
	valueNames       [][]SlotID
	slotLocs         []VarLoc
	regContents      [][]SlotID
	pendingEntries   []pendingEntry
	pendingSlotLocs  []VarLoc

	liveSlotSliceBegin int
	liveSlots          []liveSlot
}

func (c *Cache) Reset() {
	nv := sort.Search(len(c.values), func(i int) bool { return c.values[i].ID == 0 })
	xv := c.values[:nv]
	for i := range xv {
		xv[i] = Value{}
	}
	nb := sort.Search(len(c.blocks), func(i int) bool { return c.blocks[i].ID == 0 })
	xb := c.blocks[:nb]
	for i := range xb {
		xb[i] = Block{}
	}
	nl := sort.Search(len(c.locs), func(i int) bool { return c.locs[i] == nil })
	xl := c.locs[:nl]
	for i := range xl {
		xl[i] = nil
	}

	c.liveSlots = c.liveSlots[:0]
	c.liveSlotSliceBegin = 0
}

func (c *Cache) AppendLiveSlot(ls liveSlot) {
	c.liveSlots = append(c.liveSlots, ls)
}

func (c *Cache) GetLiveSlotSlice() []liveSlot {
	s := c.liveSlots[c.liveSlotSliceBegin:]
	c.liveSlotSliceBegin = len(c.liveSlots)
	return s
}
