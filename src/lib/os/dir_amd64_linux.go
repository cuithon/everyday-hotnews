// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import (
	"syscall";
	"os";
	"unsafe";
)

func clen(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return len(n)
}

// Negative count means read until EOF.
func Readdirnames(fd *FD, count int) (names []string, err *os.Error) {
	// The buffer should be at least a block long.
	// TODO(r): use fstatfs to find fs block size.
	var buf = make([]syscall.Dirent, 8192/unsafe.Sizeof(*new(syscall.Dirent)));
	names = make([]string, 0, 100);	// TODO: could be smarter about size
	for {
		if count == 0 {
			break
		}
		ret, err2 := syscall.Getdents(fd.fd, &buf[0], int64(len(buf) * unsafe.Sizeof(buf[0])));
		if ret < 0 || err2 != 0 {
			return names, os.ErrnoToError(err2)
		}
		if ret == 0 {
			break
		}
		for w, i := uintptr(0),uintptr(0); i < uintptr(ret); i += w {
			if count == 0 {
				break
			}
			dir := unsafe.Pointer((uintptr(unsafe.Pointer(&buf[0])) + i)).(*syscall.Dirent);
			w = uintptr(dir.Reclen);
			if dir.Ino == 0 {
				continue
			}
			count--;
			if len(names) == cap(names) {
				nnames := make([]string, len(names), 2*len(names));
				for i := 0; i < len(names); i++ {
					nnames[i] = names[i]
				}
				names = nnames;
			}
			names = names[0:len(names)+1];
			names[len(names)-1] = string(dir.Name[0:clen(dir.Name)]);
		}
	}
	return names, nil;
}
