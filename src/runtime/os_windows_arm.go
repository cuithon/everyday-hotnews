// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

//go:nosplit
func cputicks() int64 {
	return nanotime()
}

func checkgoarm() {
	if goarm < 7 {
		print("Need atomic synchronization instructions, coprocessor ",
			"access instructions. Recompile using GOARM=7.\n")
		exit(1)
	}
}
