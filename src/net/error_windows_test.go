// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import "syscall"

var (
	errTimedout       = syscall.ETIMEDOUT
	errOpNotSupported = syscall.EOPNOTSUPP

	abortedConnRequestErrors = []error{syscall.ERROR_NETNAME_DELETED, syscall.WSAECONNRESET} // see accept in fd_windows.go
)

func isPlatformError(err error) bool {
	_, ok := err.(syscall.Errno)
	return ok
}
