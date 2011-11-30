// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import (
	"syscall"
	"time"
)

func isSymlink(stat *syscall.Stat_t) bool {
	return stat.Mode&syscall.S_IFMT == syscall.S_IFLNK
}

func fileInfoFromStat(name string, fi *FileInfo, lstat, stat *syscall.Stat_t) *FileInfo {
	fi.Dev = uint64(stat.Dev)
	fi.Ino = stat.Ino
	fi.Nlink = uint64(stat.Nlink)
	fi.Mode = uint32(stat.Mode)
	fi.Uid = int(stat.Uid)
	fi.Gid = int(stat.Gid)
	fi.Rdev = uint64(stat.Rdev)
	fi.Size = stat.Size
	fi.Blksize = int64(stat.Blksize)
	fi.Blocks = stat.Blocks
	fi.AccessTime = timespecToTime(stat.Atimespec)
	fi.ModTime = timespecToTime(stat.Mtimespec)
	fi.ChangeTime = timespecToTime(stat.Ctimespec)
	fi.Name = basename(name)
	if isSymlink(lstat) && !isSymlink(stat) {
		fi.FollowedSymlink = true
	}
	return fi
}

func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}
