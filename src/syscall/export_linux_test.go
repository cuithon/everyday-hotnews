// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syscall

var (
	RawSyscallNoError = rawSyscallNoError
	ForceClone3       = &forceClone3
)

const (
	Sys_GETEUID           = sys_GETEUID
	Sys_pidfd_send_signal = _SYS_pidfd_send_signal
)
