// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import (
	"syscall"
	"unsafe"
)

const _P_PID = 1 // not 0 as on FreeBSD and Dragonfly!

func wait6(idtype, id, options int) (status int, errno syscall.Errno) {
	var status32 int32 // C.int
	_, _, errno = syscall.Syscall6(syscall.SYS_WAIT6, uintptr(idtype), uintptr(id), uintptr(unsafe.Pointer(&status32)), uintptr(options), 0, 0)
	return int(status32), errno
}
