// compile -d=libfuzzer

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p

func f(x, y int) {
	_ = x > y
	_ = y > x
}
