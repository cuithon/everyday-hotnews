// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ld

import (
	"cmd/link/internal/loader"
	"sync"
)

// Assembling the binary is broken into two steps:
//  - writing out the code/data/dwarf Segments
//  - writing out the architecture specific pieces.
// This function handles the first part.
func asmb(ctxt *Link, ldr *loader.Loader) {
	// TODO(jfaller): delete me.
	if thearch.Asmb != nil {
		thearch.Asmb(ctxt, ldr)
		return
	}

	if ctxt.IsELF {
		Asmbelfsetup()
	}

	var wg sync.WaitGroup
	sect := Segtext.Sections[0]
	offset := sect.Vaddr - Segtext.Vaddr + Segtext.Fileoff
	f := func(ctxt *Link, out *OutBuf, start, length int64) {
		pad := thearch.CodePad
		if pad == nil {
			pad = zeros[:]
		}
		CodeblkPad(ctxt, out, start, length, pad)
	}

	if !thearch.WriteTextBlocks {
		writeParallel(&wg, f, ctxt, offset, sect.Vaddr, sect.Length)
		for _, sect := range Segtext.Sections[1:] {
			offset := sect.Vaddr - Segtext.Vaddr + Segtext.Fileoff
			writeParallel(&wg, Datblk, ctxt, offset, sect.Vaddr, sect.Length)
		}
	} else {
		// TODO why can't we handle all sections this way?
		for _, sect := range Segtext.Sections {
			offset := sect.Vaddr - Segtext.Vaddr + Segtext.Fileoff
			// Handle additional text sections with Codeblk
			if sect.Name == ".text" {
				writeParallel(&wg, f, ctxt, offset, sect.Vaddr, sect.Length)
			} else {
				writeParallel(&wg, Datblk, ctxt, offset, sect.Vaddr, sect.Length)
			}
		}
	}

	if Segrodata.Filelen > 0 {
		writeParallel(&wg, Datblk, ctxt, Segrodata.Fileoff, Segrodata.Vaddr, Segrodata.Filelen)
	}

	if Segrelrodata.Filelen > 0 {
		writeParallel(&wg, Datblk, ctxt, Segrelrodata.Fileoff, Segrelrodata.Vaddr, Segrelrodata.Filelen)
	}

	writeParallel(&wg, Datblk, ctxt, Segdata.Fileoff, Segdata.Vaddr, Segdata.Filelen)

	writeParallel(&wg, dwarfblk, ctxt, Segdwarf.Fileoff, Segdwarf.Vaddr, Segdwarf.Filelen)

	wg.Wait()
}
