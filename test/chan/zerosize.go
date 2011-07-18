// $G $D/$F.go && $L $F.$A && ./$A.out

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Making channels of a zero-sized type should not panic.

package main

func main() {
	_ = make(chan [0]byte)
	_ = make(chan [0]byte, 1)
	_ = make(chan struct{})
	_ = make(chan struct{}, 1)
}
