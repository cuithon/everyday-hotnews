// asmcheck

// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codegen

func issue63332(c chan int) {
	x := 0
	// amd64:-`MOVQ`
	x += 2
	c <- x
}
