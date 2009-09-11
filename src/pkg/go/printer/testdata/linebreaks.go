// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linebreaks

import (
	"bytes";
	"fmt";
	"io";
	"os";
	"reflect";
	"strings";
	"testing";
)

type untarTest struct {
	file string;
	headers []*Header;
}

var untarTests = []*untarTest{
	&untarTest{
		file: "testdata/gnu.tar",
		headers: []*Header{
			&Header{
				Name: "small.txt",
				Mode: 0640,
				Uid: 73025,
				Gid: 5000,
				Size: 5,
				Mtime: 1244428340,
				Typeflag: '0',
				Uname: "dsymonds",
				Gname: "eng",
			},
			&Header{
				Name: "small2.txt",
				Mode: 0640,
				Uid: 73025,
				Gid: 5000,
				Size: 11,
				Mtime: 1244436044,
				Typeflag: '0',
				Uname: "dsymonds",
				Gname: "eng",
			},
		},
	},
	&untarTest{
		file: "testdata/star.tar",
		headers: []*Header{
			&Header{
				Name: "small.txt",
				Mode: 0640,
				Uid: 73025,
				Gid: 5000,
				Size: 5,
				Mtime: 1244592783,
				Typeflag: '0',
				Uname: "dsymonds",
				Gname: "eng",
				Atime: 1244592783,
				Ctime: 1244592783,
			},
			&Header{
				Name: "small2.txt",
				Mode: 0640,
				Uid: 73025,
				Gid: 5000,
				Size: 11,
				Mtime: 1244592783,
				Typeflag: '0',
				Uname: "dsymonds",
				Gname: "eng",
				Atime: 1244592783,
				Ctime: 1244592783,
			},
		},
	},
	&untarTest{
		file: "testdata/v7.tar",
		headers: []*Header{
			&Header{
				Name: "small.txt",
				Mode: 0444,
				Uid: 73025,
				Gid: 5000,
				Size: 5,
				Mtime: 1244593104,
				Typeflag: '\x00',
			},
			&Header{
				Name: "small2.txt",
				Mode: 0444,
				Uid: 73025,
				Gid: 5000,
				Size: 11,
				Mtime: 1244593104,
				Typeflag: '\x00',
			},
		},
	},
}

var facts = map[int] string {
	0: "1",
	1: "1",
	2: "2",
	10: "3628800",
	20: "2432902008176640000",
	100: "933262154439441526816992388562667004907159682643816214685929"
		"638952175999932299156089414639761565182862536979208272237582"
		"51185210916864000000000000000000000000",
}

func TestReader(t *testing.T) {
testLoop:
	for i, test := range untarTests {
		f, err := os.Open(test.file, os.O_RDONLY, 0444);
		if err != nil {
			t.Errorf("test %d: Unexpected error: %v", i, err);
			continue
		}
		tr := NewReader(f);
		for j, header := range test.headers {
			hdr, err := tr.Next();
			if err != nil || hdr == nil {
				t.Errorf("test %d, entry %d: Didn't get entry: %v", i, j, err);
				f.Close();
				continue testLoop
			}
			if !reflect.DeepEqual(hdr, header) {
				t.Errorf("test %d, entry %d: Incorrect header:\nhave %+v\nwant %+v",
					 i, j, *hdr, *header);
			}
		}
		hdr, err := tr.Next();
		if hdr != nil || err != nil {
			t.Errorf("test %d: Unexpected entry or error: hdr=%v err=%v", i, err);
		}
		f.Close();
	}
}
