// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unicode_test

import (
	"testing";
	. "unicode";
)

var testDigit = []int {
	0x0030,
	0x0039,
	0x0661,
	0x06F1,
	0x07C9,
	0x0966,
	0x09EF,
	0x0A66,
	0x0AEF,
	0x0B66,
	0x0B6F,
	0x0BE6,
	0x0BEF,
	0x0C66,
	0x0CEF,
	0x0D66,
	0x0D6F,
	0x0E50,
	0x0E59,
	0x0ED0,
	0x0ED9,
	0x0F20,
	0x0F29,
	0x1040,
	0x1049,
	0x1090,
	0x1091,
	0x1099,
	0x17E0,
	0x17E9,
	0x1810,
	0x1819,
	0x1946,
	0x194F,
	0x19D0,
	0x19D9,
	0x1B50,
	0x1B59,
	0x1BB0,
	0x1BB9,
	0x1C40,
	0x1C49,
	0x1C50,
	0x1C59,
	0xA620,
	0xA629,
	0xA8D0,
	0xA8D9,
	0xA900,
	0xA909,
	0xAA50,
	0xAA59,
	0xFF10,
	0xFF19,
	0x104A1,
	0x1D7CE,
}

var testLetter = []int {
	0x0041,
	0x0061,
	0x00AA,
	0x00BA,
	0x00C8,
	0x00DB,
	0x00F9,
	0x02EC,
	0x0535,
	0x06E6,
	0x093D,
	0x0A15,
	0x0B99,
	0x0DC0,
	0x0EDD,
	0x1000,
	0x1200,
	0x1312,
	0x1401,
	0x1885,
	0x2C00,
	0xA800,
	0xF900,
	0xFA30,
	0xFFDA,
	0xFFDC,
	0x10000,
	0x10300,
	0x10400,
	0x20000,
	0x2F800,
	0x2FA1D,
}

func TestDigit(t *testing.T) {
	for i, r := range testDigit {
		if !IsDigit(r) {
			t.Errorf("IsDigit(U+%04X) = false, want true\n", r);
		}
	}
	for i, r := range testLetter {
		if IsDigit(r) {
			t.Errorf("IsDigit(U+%04X) = true, want false\n", r);
		}
	}
}

// Test that the special case in IsDigit agrees with the table
func TestDigitOptimization(t *testing.T) {
	for i := 0; i < 0x100; i++ {
		if Is(Digit, i) != IsDigit(i) {
			t.Errorf("IsDigit(U+%04X) disagrees with Is(Digit)", i)
		}
	}
}
