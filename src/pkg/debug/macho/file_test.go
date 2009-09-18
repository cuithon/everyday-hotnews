// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package macho

import (
	"reflect";
	"testing";
)

type fileTest struct {
	file string;
	hdr FileHeader;
	segments []*SegmentHeader;
	sections []*SectionHeader;
}

var fileTests = []fileTest {
	fileTest{
		"testdata/gcc-386-darwin-exec",
		FileHeader{0xfeedface, Cpu386, 0x3, 0x2, 0xc, 0x3c0, 0x85},
		[]*SegmentHeader{
			&SegmentHeader{LoadCmdSegment, 0x38, "__PAGEZERO", 0x0, 0x1000, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			&SegmentHeader{LoadCmdSegment, 0xc0, "__TEXT", 0x1000, 0x1000, 0x0, 0x1000, 0x7, 0x5, 0x2, 0x0},
			&SegmentHeader{LoadCmdSegment, 0xc0, "__DATA", 0x2000, 0x1000, 0x1000, 0x1000, 0x7, 0x3, 0x2, 0x0},
			&SegmentHeader{LoadCmdSegment, 0x7c, "__IMPORT", 0x3000, 0x1000, 0x2000, 0x1000, 0x7, 0x7, 0x1, 0x0},
			&SegmentHeader{LoadCmdSegment, 0x38, "__LINKEDIT", 0x4000, 0x1000, 0x3000, 0x12c, 0x7, 0x1, 0x0, 0x0},
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		},
		[]*SectionHeader{
			&SectionHeader{"__text", "__TEXT", 0x1f68, 0x88, 0xf68, 0x2, 0x0, 0x0, 0x80000400},
			&SectionHeader{"__cstring", "__TEXT", 0x1ff0, 0xd, 0xff0, 0x0, 0x0, 0x0, 0x2},
			&SectionHeader{"__data", "__DATA", 0x2000, 0x14, 0x1000, 0x2, 0x0, 0x0, 0x0},
			&SectionHeader{"__dyld", "__DATA", 0x2014, 0x1c, 0x1014, 0x2, 0x0, 0x0, 0x0},
			&SectionHeader{"__jump_table", "__IMPORT", 0x3000, 0xa, 0x2000, 0x6, 0x0, 0x0, 0x4000008},
		},
	},
	fileTest{
		"testdata/gcc-amd64-darwin-exec",
		FileHeader{0xfeedfacf, CpuAmd64, 0x80000003, 0x2, 0xb, 0x568, 0x85},
		[]*SegmentHeader{
			&SegmentHeader{LoadCmdSegment64, 0x48, "__PAGEZERO", 0x0, 0x100000000, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			&SegmentHeader{LoadCmdSegment64, 0x1d8, "__TEXT", 0x100000000, 0x1000, 0x0, 0x1000, 0x7, 0x5, 0x5, 0x0},
			&SegmentHeader{LoadCmdSegment64, 0x138, "__DATA", 0x100001000, 0x1000, 0x1000, 0x1000, 0x7, 0x3, 0x3, 0x0},
			&SegmentHeader{LoadCmdSegment64, 0x48, "__LINKEDIT", 0x100002000, 0x1000, 0x2000, 0x140, 0x7, 0x1, 0x0, 0x0},
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		},
		[]*SectionHeader{
			&SectionHeader{"__text", "__TEXT", 0x100000f14, 0x6d, 0xf14, 0x2, 0x0, 0x0, 0x80000400},
			&SectionHeader{"__symbol_stub1", "__TEXT", 0x100000f81, 0xc, 0xf81, 0x0, 0x0, 0x0, 0x80000408},
			&SectionHeader{"__stub_helper", "__TEXT", 0x100000f90, 0x18, 0xf90, 0x2, 0x0, 0x0, 0x0},
			&SectionHeader{"__cstring", "__TEXT", 0x100000fa8, 0xd, 0xfa8, 0x0, 0x0, 0x0, 0x2},
			&SectionHeader{"__eh_frame", "__TEXT", 0x100000fb8, 0x48, 0xfb8, 0x3, 0x0, 0x0, 0x6000000b},
			&SectionHeader{"__data", "__DATA", 0x100001000, 0x1c, 0x1000, 0x3, 0x0, 0x0, 0x0},
			&SectionHeader{"__dyld", "__DATA", 0x100001020, 0x38, 0x1020, 0x3, 0x0, 0x0, 0x0},
			&SectionHeader{"__la_symbol_ptr", "__DATA", 0x100001058, 0x10, 0x1058, 0x2, 0x0, 0x0, 0x7},
		},
	},
	fileTest{
		"testdata/gcc-amd64-darwin-exec-debug",
		FileHeader{0xfeedfacf, CpuAmd64, 0x80000003, 0xa, 0x4, 0x5a0, 0},
		[]*SegmentHeader{
			nil,
			&SegmentHeader{LoadCmdSegment64, 0x1d8, "__TEXT", 0x100000000, 0x1000, 0x0, 0x0, 0x7, 0x5, 0x5, 0x0},
			&SegmentHeader{LoadCmdSegment64, 0x138, "__DATA", 0x100001000, 0x1000, 0x0, 0x0, 0x7, 0x3, 0x3, 0x0},
			&SegmentHeader{LoadCmdSegment64, 0x278, "__DWARF", 0x100002000, 0x1000, 0x1000, 0x1bc, 0x7, 0x3, 0x7, 0x0},
		},
		[]*SectionHeader{
			&SectionHeader{"__text", "__TEXT", 0x100000f14, 0x0, 0x0, 0x2, 0x0, 0x0, 0x80000400},
			&SectionHeader{"__symbol_stub1", "__TEXT", 0x100000f81, 0x0, 0x0, 0x0, 0x0, 0x0, 0x80000408},
			&SectionHeader{"__stub_helper", "__TEXT", 0x100000f90, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0},
			&SectionHeader{"__cstring", "__TEXT", 0x100000fa8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2},
			&SectionHeader{"__eh_frame", "__TEXT", 0x100000fb8, 0x0, 0x0, 0x3, 0x0, 0x0, 0x6000000b},
			&SectionHeader{"__data", "__DATA", 0x100001000, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0},
			&SectionHeader{"__dyld", "__DATA", 0x100001020, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0},
			&SectionHeader{"__la_symbol_ptr", "__DATA", 0x100001058, 0x0, 0x0, 0x2, 0x0, 0x0, 0x7},
			&SectionHeader{"__debug_abbrev", "__DWARF", 0x100002000, 0x36, 0x1000, 0x0, 0x0, 0x0, 0x0},
			&SectionHeader{"__debug_aranges", "__DWARF", 0x100002036, 0x30, 0x1036, 0x0, 0x0, 0x0, 0x0},
			&SectionHeader{"__debug_frame", "__DWARF", 0x100002066, 0x40, 0x1066, 0x0, 0x0, 0x0, 0x0},
			&SectionHeader{"__debug_info", "__DWARF", 0x1000020a6, 0x54, 0x10a6, 0x0, 0x0, 0x0, 0x0},
			&SectionHeader{"__debug_line", "__DWARF", 0x1000020fa, 0x47, 0x10fa, 0x0, 0x0, 0x0, 0x0},
			&SectionHeader{"__debug_pubnames", "__DWARF", 0x100002141, 0x1b, 0x1141, 0x0, 0x0, 0x0, 0x0},
			&SectionHeader{"__debug_str", "__DWARF", 0x10000215c, 0x60, 0x115c, 0x0, 0x0, 0x0, 0x0},
		},
	},
}

func TestOpen(t *testing.T) {
	for i := range fileTests {
		tt := &fileTests[i];

		f, err := Open(tt.file);
		if err != nil {
			t.Error(err);
			continue;
		}
		if !reflect.DeepEqual(f.FileHeader, tt.hdr) {
			t.Errorf("open %s:\n\thave %#v\n\twant %#v\n", tt.file, f.FileHeader, tt.hdr);
			continue;
		}
		for i, l := range f.Loads {
			if i >= len(tt.segments) {
				break;
			}
			sh := tt.segments[i];
			s, ok := l.(*Segment);
			if sh == nil {
				if ok {
					t.Errorf("open %s, section %d: skipping %#v\n", tt.file, i, &s.SegmentHeader);
				}
				continue;
			}
			if !ok {
				t.Errorf("open %s, section %d: not *Segment\n", tt.file, i);
				continue;
			}
			have := &s.SegmentHeader;
			want := sh;
			if !reflect.DeepEqual(have, want) {
				t.Errorf("open %s, segment %d:\n\thave %#v\n\twant %#v\n", tt.file, i, have, want);
			}
		}
		tn := len(tt.segments);
		fn := len(f.Loads);
		if tn != fn {
			t.Errorf("open %s: len(Loads) = %d, want %d", tt.file, fn, tn);
		}

		for i, sh := range f.Sections {
			if i >= len(tt.sections) {
				break;
			}
			have := &sh.SectionHeader;
			want := tt.sections[i];
			if !reflect.DeepEqual(have, want) {
				t.Errorf("open %s, section %d:\n\thave %#v\n\twant %#v\n", tt.file, i, have, want);
			}
		}
		tn = len(tt.sections);
		fn = len(f.Sections);
		if tn != fn {
			t.Errorf("open %s: len(Sections) = %d, want %d", tt.file, fn, tn);
		}

	}
}
