// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9 && !windows
// +build !plan9,!windows

package exec

import (
	"io/fs"
	"syscall"
)

func init() {
	skipStdinCopyError = func(err error) bool {
		// Ignore EPIPE errors copying to stdin if the program
		// completed successfully otherwise.
		// See Issue 9173.
		pe, ok := err.(*fs.PathError)
		return ok &&
			pe.Op == "write" && pe.Path == "|1" &&
			pe.Err == syscall.EPIPE
	}
}
