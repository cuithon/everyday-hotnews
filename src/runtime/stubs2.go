// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !solaris

package runtime

import "unsafe"

func read(fd int32, p unsafe.Pointer, n int32) int32
func close(fd int32) int32

func exit(code int32)
func nanotime() int64
func usleep(usec uint32)

func mmap(addr unsafe.Pointer, n uintptr, prot, flags, fd int32, off uint32) unsafe.Pointer
func munmap(addr unsafe.Pointer, n uintptr)

//go:noescape
func write(fd uintptr, p unsafe.Pointer, n int32) int32

//go:noescape
func open(name *byte, mode, perm int32) int32

func madvise(addr unsafe.Pointer, n uintptr, flags int32)
