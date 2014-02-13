// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package macho

import (
	"reflect"
	"testing"
)

type fileTest struct {
	file     string
	hdr      FileHeader
	segments []*SegmentHeader
	sections []*SectionHeader
}

var fileTests = []fileTest{
	{
		"testdata/gcc-386-darwin-exec",
		FileHeader{0xfeedface, Cpu386, 0x3, 0x2, 0xc, 0x3c0, 0x85},
		[]*SegmentHeader{
			{LoadCmdSegment, 0x38, "__PAGEZERO", 0x0, 0x1000, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			{LoadCmdSegment, 0xc0, "__TEXT", 0x1000, 0x1000, 0x0, 0x1000, 0x7, 0x5, 0x2, 0x0},
			{LoadCmdSegment, 0xc0, "__DATA", 0x2000, 0x1000, 0x1000, 0x1000, 0x7, 0x3, 0x2, 0x0},
			{LoadCmdSegment, 0x7c, "__IMPORT", 0x3000, 0x1000, 0x2000, 0x1000, 0x7, 0x7, 0x1, 0x0},
			{LoadCmdSegment, 0x38, "__LINKEDIT", 0x4000, 0x1000, 0x3000, 0x12c, 0x7, 0x1, 0x0, 0x0},
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		},
		[]*SectionHeader{
			{"__text", "__TEXT", 0x1f68, 0x88, 0xf68, 0x2, 0x0, 0x0, 0x80000400},
			{"__cstring", "__TEXT", 0x1ff0, 0xd, 0xff0, 0x0, 0x0, 0x0, 0x2},
			{"__data", "__DATA", 0x2000, 0x14, 0x1000, 0x2, 0x0, 0x0, 0x0},
			{"__dyld", "__DATA", 0x2014, 0x1c, 0x1014, 0x2, 0x0, 0x0, 0x0},
			{"__jump_table", "__IMPORT", 0x3000, 0xa, 0x2000, 0x6, 0x0, 0x0, 0x4000008},
		},
	},
	{
		"testdata/gcc-amd64-darwin-exec",
		FileHeader{0xfeedfacf, CpuAmd64, 0x80000003, 0x2, 0xb, 0x568, 0x85},
		[]*SegmentHeader{
			{LoadCmdSegment64, 0x48, "__PAGEZERO", 0x0, 0x100000000, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			{LoadCmdSegment64, 0x1d8, "__TEXT", 0x100000000, 0x1000, 0x0, 0x1000, 0x7, 0x5, 0x5, 0x0},
			{LoadCmdSegment64, 0x138, "__DATA", 0x100001000, 0x1000, 0x1000, 0x1000, 0x7, 0x3, 0x3, 0x0},
			{LoadCmdSegment64, 0x48, "__LINKEDIT", 0x100002000, 0x1000, 0x2000, 0x140, 0x7, 0x1, 0x0, 0x0},
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		},
		[]*SectionHeader{
			{"__text", "__TEXT", 0x100000f14, 0x6d, 0xf14, 0x2, 0x0, 0x0, 0x80000400},
			{"__symbol_stub1", "__TEXT", 0x100000f81, 0xc, 0xf81, 0x0, 0x0, 0x0, 0x80000408},
			{"__stub_helper", "__TEXT", 0x100000f90, 0x18, 0xf90, 0x2, 0x0, 0x0, 0x0},
			{"__cstring", "__TEXT", 0x100000fa8, 0xd, 0xfa8, 0x0, 0x0, 0x0, 0x2},
			{"__eh_frame", "__TEXT", 0x100000fb8, 0x48, 0xfb8, 0x3, 0x0, 0x0, 0x6000000b},
			{"__data", "__DATA", 0x100001000, 0x1c, 0x1000, 0x3, 0x0, 0x0, 0x0},
			{"__dyld", "__DATA", 0x100001020, 0x38, 0x1020, 0x3, 0x0, 0x0, 0x0},
			{"__la_symbol_ptr", "__DATA", 0x100001058, 0x10, 0x1058, 0x2, 0x0, 0x0, 0x7},
		},
	},
	{
		"testdata/gcc-amd64-darwin-exec-debug",
		FileHeader{0xfeedfacf, CpuAmd64, 0x80000003, 0xa, 0x4, 0x5a0, 0},
		[]*SegmentHeader{
			nil,
			{LoadCmdSegment64, 0x1d8, "__TEXT", 0x100000000, 0x1000, 0x0, 0x0, 0x7, 0x5, 0x5, 0x0},
			{LoadCmdSegment64, 0x138, "__DATA", 0x100001000, 0x1000, 0x0, 0x0, 0x7, 0x3, 0x3, 0x0},
			{LoadCmdSegment64, 0x278, "__DWARF", 0x100002000, 0x1000, 0x1000, 0x1bc, 0x7, 0x3, 0x7, 0x0},
		},
		[]*SectionHeader{
			{"__text", "__TEXT", 0x100000f14, 0x0, 0x0, 0x2, 0x0, 0x0, 0x80000400},
			{"__symbol_stub1", "__TEXT", 0x100000f81, 0x0, 0x0, 0x0, 0x0, 0x0, 0x80000408},
			{"__stub_helper", "__TEXT", 0x100000f90, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0},
			{"__cstring", "__TEXT", 0x100000fa8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2},
			{"__eh_frame", "__TEXT", 0x100000fb8, 0x0, 0x0, 0x3, 0x0, 0x0, 0x6000000b},
			{"__data", "__DATA", 0x100001000, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0},
			{"__dyld", "__DATA", 0x100001020, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0},
			{"__la_symbol_ptr", "__DATA", 0x100001058, 0x0, 0x0, 0x2, 0x0, 0x0, 0x7},
			{"__debug_abbrev", "__DWARF", 0x100002000, 0x36, 0x1000, 0x0, 0x0, 0x0, 0x0},
			{"__debug_aranges", "__DWARF", 0x100002036, 0x30, 0x1036, 0x0, 0x0, 0x0, 0x0},
			{"__debug_frame", "__DWARF", 0x100002066, 0x40, 0x1066, 0x0, 0x0, 0x0, 0x0},
			{"__debug_info", "__DWARF", 0x1000020a6, 0x54, 0x10a6, 0x0, 0x0, 0x0, 0x0},
			{"__debug_line", "__DWARF", 0x1000020fa, 0x47, 0x10fa, 0x0, 0x0, 0x0, 0x0},
			{"__debug_pubnames", "__DWARF", 0x100002141, 0x1b, 0x1141, 0x0, 0x0, 0x0, 0x0},
			{"__debug_str", "__DWARF", 0x10000215c, 0x60, 0x115c, 0x0, 0x0, 0x0, 0x0},
		},
	},
}

func TestOpen(t *testing.T) {
	for i := range fileTests {
		tt := &fileTests[i]

		f, err := Open(tt.file)
		if err != nil {
			t.Error(err)
			continue
		}
		if !reflect.DeepEqual(f.FileHeader, tt.hdr) {
			t.Errorf("open %s:\n\thave %#v\n\twant %#v\n", tt.file, f.FileHeader, tt.hdr)
			continue
		}
		for i, l := range f.Loads {
			if i >= len(tt.segments) {
				break
			}
			sh := tt.segments[i]
			s, ok := l.(*Segment)
			if sh == nil {
				if ok {
					t.Errorf("open %s, section %d: skipping %#v\n", tt.file, i, &s.SegmentHeader)
				}
				continue
			}
			if !ok {
				t.Errorf("open %s, section %d: not *Segment\n", tt.file, i)
				continue
			}
			have := &s.SegmentHeader
			want := sh
			if !reflect.DeepEqual(have, want) {
				t.Errorf("open %s, segment %d:\n\thave %#v\n\twant %#v\n", tt.file, i, have, want)
			}
		}
		tn := len(tt.segments)
		fn := len(f.Loads)
		if tn != fn {
			t.Errorf("open %s: len(Loads) = %d, want %d", tt.file, fn, tn)
		}

		for i, sh := range f.Sections {
			if i >= len(tt.sections) {
				break
			}
			have := &sh.SectionHeader
			want := tt.sections[i]
			if !reflect.DeepEqual(have, want) {
				t.Errorf("open %s, section %d:\n\thave %#v\n\twant %#v\n", tt.file, i, have, want)
			}
		}
		tn = len(tt.sections)
		fn = len(f.Sections)
		if tn != fn {
			t.Errorf("open %s: len(Sections) = %d, want %d", tt.file, fn, tn)
		}

	}
}

func TestOpenFailure(t *testing.T) {
	filename := "file.go"    // not a Mach-O file
	_, err := Open(filename) // don't crash
	if err == nil {
		t.Errorf("open %s: succeeded unexpectedly", filename)
	}
}

func TestOpenFat(t *testing.T) {
	ff, err := OpenFat("testdata/fat-gcc-386-amd64-darwin-exec")
	if err != nil {
		t.Fatal(err)
	}

	if ff.Magic != MagicFat {
		t.Errorf("OpenFat: got magic number %#x, want %#x", ff.Magic, MagicFat)
	}
	if len(ff.Arches) != 2 {
		t.Errorf("OpenFat: got %d architectures, want 2", len(ff.Arches))
	}

	for i := range ff.Arches {
		arch := &ff.Arches[i]
		ftArch := &fileTests[i]

		if arch.Cpu != ftArch.hdr.Cpu || arch.SubCpu != ftArch.hdr.SubCpu {
			t.Error("OpenFat: architecture #%d got cpu=%#x subtype=%#x, expected cpu=%#x, subtype=%#x", i, arch.Cpu, arch.SubCpu, ftArch.hdr.Cpu, ftArch.hdr.SubCpu)
		}

		if !reflect.DeepEqual(arch.FileHeader, ftArch.hdr) {
			t.Errorf("OpenFat header:\n\tgot %#v\n\twant %#v\n", arch.FileHeader, ftArch.hdr)
		}
	}
}

func TestOpenFatFailure(t *testing.T) {
	filename := "file.go" // not a Mach-O file
	if _, err := OpenFat(filename); err == nil {
		t.Errorf("OpenFat %s: succeeded unexpectedly", filename)
	}

	filename = "testdata/gcc-386-darwin-exec" // not a fat Mach-O
	ff, err := OpenFat(filename)
	if err != ErrNotFat {
		t.Errorf("OpenFat %s: got %v, want ErrNotFat", err)
	}
	if ff != nil {
		t.Errorf("OpenFat %s: got %v, want nil", ff)
	}
}
