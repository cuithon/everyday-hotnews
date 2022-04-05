// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zlib

import (
	"bytes"
	"io"
	"testing"
)

type zlibTest struct {
	desc       string
	raw        string
	compressed []byte
	dict       []byte
	err        error
}

// Compare-to-golden test data was generated by the ZLIB example program at
// https://www.zlib.net/zpipe.c

var zlibTests = []zlibTest{
	{
		"truncated empty",
		"",
		[]byte{},
		nil,
		io.ErrUnexpectedEOF,
	},
	{
		"truncated dict",
		"",
		[]byte{0x78, 0xbb},
		[]byte{0x00},
		io.ErrUnexpectedEOF,
	},
	{
		"truncated checksum",
		"",
		[]byte{0x78, 0xbb, 0x00, 0x01, 0x00, 0x01, 0xca, 0x48,
			0xcd, 0xc9, 0xc9, 0xd7, 0x51, 0x28, 0xcf, 0x2f,
			0xca, 0x49, 0x01, 0x04, 0x00, 0x00, 0xff, 0xff,
		},
		[]byte{0x00},
		io.ErrUnexpectedEOF,
	},
	{
		"empty",
		"",
		[]byte{0x78, 0x9c, 0x03, 0x00, 0x00, 0x00, 0x00, 0x01},
		nil,
		nil,
	},
	{
		"goodbye",
		"goodbye, world",
		[]byte{
			0x78, 0x9c, 0x4b, 0xcf, 0xcf, 0x4f, 0x49, 0xaa,
			0x4c, 0xd5, 0x51, 0x28, 0xcf, 0x2f, 0xca, 0x49,
			0x01, 0x00, 0x28, 0xa5, 0x05, 0x5e,
		},
		nil,
		nil,
	},
	{
		"bad header (CINFO)",
		"",
		[]byte{0x88, 0x98, 0x03, 0x00, 0x00, 0x00, 0x00, 0x01},
		nil,
		ErrHeader,
	},
	{
		"bad header (FCHECK)",
		"",
		[]byte{0x78, 0x9f, 0x03, 0x00, 0x00, 0x00, 0x00, 0x01},
		nil,
		ErrHeader,
	},
	{
		"bad checksum",
		"",
		[]byte{0x78, 0x9c, 0x03, 0x00, 0x00, 0x00, 0x00, 0xff},
		nil,
		ErrChecksum,
	},
	{
		"not enough data",
		"",
		[]byte{0x78, 0x9c, 0x03, 0x00, 0x00, 0x00},
		nil,
		io.ErrUnexpectedEOF,
	},
	{
		"excess data is silently ignored",
		"",
		[]byte{
			0x78, 0x9c, 0x03, 0x00, 0x00, 0x00, 0x00, 0x01,
			0x78, 0x9c, 0xff,
		},
		nil,
		nil,
	},
	{
		"dictionary",
		"Hello, World!\n",
		[]byte{
			0x78, 0xbb, 0x1c, 0x32, 0x04, 0x27, 0xf3, 0x00,
			0xb1, 0x75, 0x20, 0x1c, 0x45, 0x2e, 0x00, 0x24,
			0x12, 0x04, 0x74,
		},
		[]byte{
			0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x0a,
		},
		nil,
	},
	{
		"wrong dictionary",
		"",
		[]byte{
			0x78, 0xbb, 0x1c, 0x32, 0x04, 0x27, 0xf3, 0x00,
			0xb1, 0x75, 0x20, 0x1c, 0x45, 0x2e, 0x00, 0x24,
			0x12, 0x04, 0x74,
		},
		[]byte{
			0x48, 0x65, 0x6c, 0x6c,
		},
		ErrDictionary,
	},
	{
		"truncated zlib stream amid raw-block",
		"hello",
		[]byte{
			0x78, 0x9c, 0x00, 0x0c, 0x00, 0xf3, 0xff, 0x68, 0x65, 0x6c, 0x6c, 0x6f,
		},
		nil,
		io.ErrUnexpectedEOF,
	},
	{
		"truncated zlib stream amid fixed-block",
		"He",
		[]byte{
			0x78, 0x9c, 0xf2, 0x48, 0xcd,
		},
		nil,
		io.ErrUnexpectedEOF,
	},
}

func TestDecompressor(t *testing.T) {
	b := new(bytes.Buffer)
	for _, tt := range zlibTests {
		in := bytes.NewReader(tt.compressed)
		zr, err := NewReaderDict(in, tt.dict)
		if err != nil {
			if err != tt.err {
				t.Errorf("%s: NewReader: %s", tt.desc, err)
			}
			continue
		}
		defer zr.Close()

		// Read and verify correctness of data.
		b.Reset()
		n, err := io.Copy(b, zr)
		if err != nil {
			if err != tt.err {
				t.Errorf("%s: io.Copy: %v want %v", tt.desc, err, tt.err)
			}
			continue
		}
		s := b.String()
		if s != tt.raw {
			t.Errorf("%s: got %d-byte %q want %d-byte %q", tt.desc, n, s, len(tt.raw), tt.raw)
		}

		// Check for sticky errors.
		if n, err := zr.Read([]byte{0}); n != 0 || err != io.EOF {
			t.Errorf("%s: Read() = (%d, %v), want (0, io.EOF)", tt.desc, n, err)
		}
		if err := zr.Close(); err != nil {
			t.Errorf("%s: Close() = %v, want nil", tt.desc, err)
		}
	}
}
