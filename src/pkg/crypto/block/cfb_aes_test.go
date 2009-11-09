// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// CFB AES test vectors.

// See U.S. National Institute of Standards and Technology (NIST)
// Special Publication 800-38A, ``Recommendation for Block Cipher
// Modes of Operation,'' 2001 Edition, pp. 29-52.

package block

import (
	"bytes";
	"crypto/aes";
	"io";
	"testing";
)

type cfbTest struct {
	name	string;
	s	int;
	key	[]byte;
	iv	[]byte;
	in	[]byte;
	out	[]byte;
}

var cfbAESTests = []cfbTest{
	cfbTest{
		"CFB1-AES128",
		1,
		commonKey128,
		commonIV,
		[]byte{
			0<<7 | 1<<6 | 1<<5 | 0<<4 | 1<<3 | 0<<2 | 1<<1,
			1<<7 | 1<<6 | 0<<5 | 0<<4 | 0<<3 | 0<<2 | 0<<1,
		},
		[]byte{
			0<<7 | 1<<6 | 1<<5 | 0<<4 | 1<<3 | 0<<2 | 0<<1,
			1<<7 | 0<<6 | 1<<5 | 1<<4 | 0<<3 | 0<<2 | 1<<1,
		},
	},
	cfbTest{
		"CFB1-AES192",
		1,
		commonKey192,
		commonIV,
		[]byte{
			0<<7 | 1<<6 | 1<<5 | 0<<4 | 1<<3 | 0<<2 | 1<<1,
			1<<7 | 1<<6 | 0<<5 | 0<<4 | 0<<3 | 0<<2 | 0<<1,
		},
		[]byte{
			1<<7 | 0<<6 | 0<<5 | 1<<4 | 0<<3 | 0<<2 | 1<<1,
			0<<7 | 1<<6 | 0<<5 | 1<<4 | 1<<3 | 0<<2 | 0<<1,
		},
	},
	cfbTest{
		"CFB1-AES256",
		1,
		commonKey256,
		commonIV,
		[]byte{
			0<<7 | 1<<6 | 1<<5 | 0<<4 | 1<<3 | 0<<2 | 1<<1,
			1<<7 | 1<<6 | 0<<5 | 0<<4 | 0<<3 | 0<<2 | 0<<1,
		},
		[]byte{
			1<<7 | 0<<6 | 0<<5 | 1<<4 | 0<<3 | 0<<2 | 0<<1,
			0<<7 | 0<<6 | 1<<5 | 0<<4 | 1<<3 | 0<<2 | 0<<1,
		},
	},

	cfbTest{
		"CFB8-AES128",
		8,
		commonKey128,
		commonIV,
		[]byte{
			0x6b,
			0xc1,
			0xbe,
			0xe2,
			0x2e,
			0x40,
			0x9f,
			0x96,
			0xe9,
			0x3d,
			0x7e,
			0x11,
			0x73,
			0x93,
			0x17,
			0x2a,
			0xae,
			0x2d,
		},
		[]byte{
			0x3b,
			0x79,
			0x42,
			0x4c,
			0x9c,
			0x0d,
			0xd4,
			0x36,
			0xba,
			0xce,
			0x9e,
			0x0e,
			0xd4,
			0x58,
			0x6a,
			0x4f,
			0x32,
			0xb9,
		},
	},

	cfbTest{
		"CFB8-AES192",
		8,
		commonKey192,
		commonIV,
		[]byte{
			0x6b,
			0xc1,
			0xbe,
			0xe2,
			0x2e,
			0x40,
			0x9f,
			0x96,
			0xe9,
			0x3d,
			0x7e,
			0x11,
			0x73,
			0x93,
			0x17,
			0x2a,
			0xae,
			0x2d,
		},
		[]byte{
			0xcd,
			0xa2,
			0x52,
			0x1e,
			0xf0,
			0xa9,
			0x05,
			0xca,
			0x44,
			0xcd,
			0x05,
			0x7c,
			0xbf,
			0x0d,
			0x47,
			0xa0,
			0x67,
			0x8a,
		},
	},

	cfbTest{
		"CFB8-AES256",
		8,
		commonKey256,
		commonIV,
		[]byte{
			0x6b,
			0xc1,
			0xbe,
			0xe2,
			0x2e,
			0x40,
			0x9f,
			0x96,
			0xe9,
			0x3d,
			0x7e,
			0x11,
			0x73,
			0x93,
			0x17,
			0x2a,
			0xae,
			0x2d,
		},
		[]byte{
			0xdc,
			0x1f,
			0x1a,
			0x85,
			0x20,
			0xa6,
			0x4d,
			0xb5,
			0x5f,
			0xcc,
			0x8a,
			0xc5,
			0x54,
			0x84,
			0x4e,
			0x88,
			0x97,
			0x00,
		},
	},

	cfbTest{
		"CFB128-AES128",
		128,
		commonKey128,
		commonIV,
		[]byte{
			0x6b, 0xc1, 0xbe, 0xe2, 0x2e, 0x40, 0x9f, 0x96, 0xe9, 0x3d, 0x7e, 0x11, 0x73, 0x93, 0x17, 0x2a,
			0xae, 0x2d, 0x8a, 0x57, 0x1e, 0x03, 0xac, 0x9c, 0x9e, 0xb7, 0x6f, 0xac, 0x45, 0xaf, 0x8e, 0x51,
			0x30, 0xc8, 0x1c, 0x46, 0xa3, 0x5c, 0xe4, 0x11, 0xe5, 0xfb, 0xc1, 0x19, 0x1a, 0x0a, 0x52, 0xef,
			0xf6, 0x9f, 0x24, 0x45, 0xdf, 0x4f, 0x9b, 0x17, 0xad, 0x2b, 0x41, 0x7b, 0xe6, 0x6c, 0x37, 0x10,
		},
		[]byte{
			0x3b, 0x3f, 0xd9, 0x2e, 0xb7, 0x2d, 0xad, 0x20, 0x33, 0x34, 0x49, 0xf8, 0xe8, 0x3c, 0xfb, 0x4a,
			0xc8, 0xa6, 0x45, 0x37, 0xa0, 0xb3, 0xa9, 0x3f, 0xcd, 0xe3, 0xcd, 0xad, 0x9f, 0x1c, 0xe5, 0x8b,
			0x26, 0x75, 0x1f, 0x67, 0xa3, 0xcb, 0xb1, 0x40, 0xb1, 0x80, 0x8c, 0xf1, 0x87, 0xa4, 0xf4, 0xdf,
			0xc0, 0x4b, 0x05, 0x35, 0x7c, 0x5d, 0x1c, 0x0e, 0xea, 0xc4, 0xc6, 0x6f, 0x9f, 0xf7, 0xf2, 0xe6,
		},
	},

	cfbTest{
		"CFB128-AES192",
		128,
		commonKey192,
		commonIV,
		[]byte{
			0x6b, 0xc1, 0xbe, 0xe2, 0x2e, 0x40, 0x9f, 0x96, 0xe9, 0x3d, 0x7e, 0x11, 0x73, 0x93, 0x17, 0x2a,
			0xae, 0x2d, 0x8a, 0x57, 0x1e, 0x03, 0xac, 0x9c, 0x9e, 0xb7, 0x6f, 0xac, 0x45, 0xaf, 0x8e, 0x51,
			0x30, 0xc8, 0x1c, 0x46, 0xa3, 0x5c, 0xe4, 0x11, 0xe5, 0xfb, 0xc1, 0x19, 0x1a, 0x0a, 0x52, 0xef,
			0xf6, 0x9f, 0x24, 0x45, 0xdf, 0x4f, 0x9b, 0x17, 0xad, 0x2b, 0x41, 0x7b, 0xe6, 0x6c, 0x37, 0x10,
		},
		[]byte{
			0xcd, 0xc8, 0x0d, 0x6f, 0xdd, 0xf1, 0x8c, 0xab, 0x34, 0xc2, 0x59, 0x09, 0xc9, 0x9a, 0x41, 0x74,
			0x67, 0xce, 0x7f, 0x7f, 0x81, 0x17, 0x36, 0x21, 0x96, 0x1a, 0x2b, 0x70, 0x17, 0x1d, 0x3d, 0x7a,
			0x2e, 0x1e, 0x8a, 0x1d, 0xd5, 0x9b, 0x88, 0xb1, 0xc8, 0xe6, 0x0f, 0xed, 0x1e, 0xfa, 0xc4, 0xc9,
			0xc0, 0x5f, 0x9f, 0x9c, 0xa9, 0x83, 0x4f, 0xa0, 0x42, 0xae, 0x8f, 0xba, 0x58, 0x4b, 0x09, 0xff,
		},
	},

	cfbTest{
		"CFB128-AES256",
		128,
		commonKey256,
		commonIV,
		[]byte{
			0x6b, 0xc1, 0xbe, 0xe2, 0x2e, 0x40, 0x9f, 0x96, 0xe9, 0x3d, 0x7e, 0x11, 0x73, 0x93, 0x17, 0x2a,
			0xae, 0x2d, 0x8a, 0x57, 0x1e, 0x03, 0xac, 0x9c, 0x9e, 0xb7, 0x6f, 0xac, 0x45, 0xaf, 0x8e, 0x51,
			0x30, 0xc8, 0x1c, 0x46, 0xa3, 0x5c, 0xe4, 0x11, 0xe5, 0xfb, 0xc1, 0x19, 0x1a, 0x0a, 0x52, 0xef,
			0xf6, 0x9f, 0x24, 0x45, 0xdf, 0x4f, 0x9b, 0x17, 0xad, 0x2b, 0x41, 0x7b, 0xe6, 0x6c, 0x37, 0x10,
		},
		[]byte{
			0xdc, 0x7e, 0x84, 0xbf, 0xda, 0x79, 0x16, 0x4b, 0x7e, 0xcd, 0x84, 0x86, 0x98, 0x5d, 0x38, 0x60,
			0x39, 0xff, 0xed, 0x14, 0x3b, 0x28, 0xb1, 0xc8, 0x32, 0x11, 0x3c, 0x63, 0x31, 0xe5, 0x40, 0x7b,
			0xdf, 0x10, 0x13, 0x24, 0x15, 0xe5, 0x4b, 0x92, 0xa1, 0x3e, 0xd0, 0xa8, 0x26, 0x7a, 0xe2, 0xf9,
			0x75, 0xa3, 0x85, 0x74, 0x1a, 0xb9, 0xce, 0xf8, 0x20, 0x31, 0x62, 0x3d, 0x55, 0xb1, 0xe4, 0x71,
		},
	},
}

func TestCFB_AES(t *testing.T) {
	for _, tt := range cfbAESTests {
		test := tt.name;

		if tt.s == 1 {
			// 1-bit CFB not implemented
			continue
		}

		c, err := aes.NewCipher(tt.key);
		if err != nil {
			t.Errorf("%s: NewCipher(%d bytes) = %s", test, len(tt.key), err);
			continue;
		}

		var crypt bytes.Buffer;
		w := NewCFBEncrypter(c, tt.s, tt.iv, &crypt);
		var r io.Reader = bytes.NewBuffer(tt.in);
		n, err := io.Copy(w, r);
		if n != int64(len(tt.in)) || err != nil {
			t.Errorf("%s: CFBEncrypter io.Copy = %d, %v want %d, nil", test, n, err, len(tt.in))
		} else if d := crypt.Bytes(); !same(tt.out, d) {
			t.Errorf("%s: CFBEncrypter\nhave %x\nwant %x", test, d, tt.out)
		}

		var plain bytes.Buffer;
		r = NewCFBDecrypter(c, tt.s, tt.iv, bytes.NewBuffer(tt.out));
		w = &plain;
		n, err = io.Copy(w, r);
		if n != int64(len(tt.out)) || err != nil {
			t.Errorf("%s: CFBDecrypter io.Copy = %d, %v want %d, nil", test, n, err, len(tt.out))
		} else if d := plain.Bytes(); !same(tt.in, d) {
			t.Errorf("%s: CFBDecrypter\nhave %x\nwant %x", test, d, tt.in)
		}

		if t.Failed() {
			break
		}
	}
}
