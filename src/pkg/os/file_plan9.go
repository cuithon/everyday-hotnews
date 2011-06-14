// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import (
	"runtime"
	"syscall"
)

// Auxiliary information if the File describes a directory
type dirInfo struct {
	buf  [syscall.STATMAX]byte // buffer for directory I/O
	nbuf int                   // length of buf; return value from Read
	bufp int                   // location of next record in buf.
}

func epipecheck(file *File, e syscall.Error) {
}


// DevNull is the name of the operating system's ``null device.''
// On Unix-like systems, it is "/dev/null"; on Windows, "NUL".
const DevNull = "/dev/null"

// OpenFile is the generalized open call; most users will use Open
// or Create instead.  It opens the named file with specified flag
// (O_RDONLY etc.) and perm, (0666 etc.) if applicable.  If successful,
// methods on the returned File can be used for I/O.
// It returns the File and an Error, if any.
func OpenFile(name string, flag int, perm uint32) (file *File, err Error) {
	var (
		fd     int
		e      syscall.Error
		create bool
		excl   bool
		trunc  bool
		append bool
	)

	if flag&O_CREATE == O_CREATE {
		flag = flag & ^O_CREATE
		create = true
	}
	if flag&O_EXCL == O_EXCL {
		excl = true
	}
	if flag&O_TRUNC == O_TRUNC {
		trunc = true
	}
	// O_APPEND is emulated on Plan 9
	if flag&O_APPEND == O_APPEND {
		flag = flag &^ O_APPEND
		append = true
	}

	syscall.ForkLock.RLock()
	if (create && trunc) || excl {
		fd, e = syscall.Create(name, flag, perm)
	} else {
		fd, e = syscall.Open(name, flag)
		if e != nil && create {
			var e1 syscall.Error
			fd, e1 = syscall.Create(name, flag, perm)
			if e1 == nil {
				e = nil
			}
		}
	}
	syscall.ForkLock.RUnlock()

	if e != nil {
		return nil, &PathError{"open", name, e}
	}

	if append {
		if _, e = syscall.Seek(fd, 0, SEEK_END); e != nil {
			return nil, &PathError{"seek", name, e}
		}
	}

	return NewFile(fd, name), nil
}

// Close closes the File, rendering it unusable for I/O.
// It returns an Error, if any.
func (file *File) Close() Error {
	if file == nil || file.fd < 0 {
		return Ebadfd
	}
	var err Error
	syscall.ForkLock.RLock()
	if e := syscall.Close(file.fd); e != nil {
		err = &PathError{"close", file.name, e}
	}
	syscall.ForkLock.RUnlock()
	file.fd = -1 // so it can't be closed again

	// no need for a finalizer anymore
	runtime.SetFinalizer(file, nil)
	return err
}

// Stat returns the FileInfo structure describing file.
// It returns the FileInfo and an error, if any.
func (f *File) Stat() (fi *FileInfo, err Error) {
	d, err := dirstat(f)
	if iserror(err) {
		return nil, err
	}
	return fileInfoFromStat(new(FileInfo), d), err
}

// Truncate changes the size of the file.
// It does not change the I/O offset.
func (f *File) Truncate(size int64) Error {
	var d Dir
	d.Null()

	d.Length = uint64(size)

	if e := syscall.Fwstat(f.fd, pdir(nil, &d)); iserror(e) {
		return &PathError{"truncate", f.name, e}
	}
	return nil
}

// Chmod changes the mode of the file to mode.
func (f *File) Chmod(mode uint32) Error {
	var d Dir
	var mask = ^uint32(0777)

	d.Null()
	odir, e := dirstat(f)
	if iserror(e) {
		return &PathError{"chmod", f.name, e}
	}

	d.Mode = (odir.Mode & mask) | (mode &^ mask)
	if e := syscall.Fwstat(f.fd, pdir(nil, &d)); iserror(e) {
		return &PathError{"chmod", f.name, e}
	}
	return nil
}

// Sync commits the current contents of the file to stable storage.
// Typically, this means flushing the file system's in-memory copy
// of recently written data to disk.
func (f *File) Sync() (err Error) {
	if f == nil {
		return EINVAL
	}

	var d Dir
	d.Null()

	if e := syscall.Fwstat(f.fd, pdir(nil, &d)); iserror(e) {
		return NewSyscallError("fsync", e)
	}
	return nil
}

// read reads up to len(b) bytes from the File.
// It returns the number of bytes read and an error, if any.
func (f *File) read(b []byte) (n int, err syscall.Error) {
	return syscall.Read(f.fd, b)
}

// pread reads len(b) bytes from the File starting at byte offset off.
// It returns the number of bytes read and the error, if any.
// EOF is signaled by a zero count with err set to nil.
func (f *File) pread(b []byte, off int64) (n int, err syscall.Error) {
	return syscall.Pread(f.fd, b, off)
}

// write writes len(b) bytes to the File.
// It returns the number of bytes written and an error, if any.
func (f *File) write(b []byte) (n int, err syscall.Error) {
	return syscall.Write(f.fd, b)
}

// pwrite writes len(b) bytes to the File starting at byte offset off.
// It returns the number of bytes written and an error, if any.
func (f *File) pwrite(b []byte, off int64) (n int, err syscall.Error) {
	return syscall.Pwrite(f.fd, b, off)
}

// seek sets the offset for the next Read or Write on file to offset, interpreted
// according to whence: 0 means relative to the origin of the file, 1 means
// relative to the current offset, and 2 means relative to the end.
// It returns the new offset and an error, if any.
func (f *File) seek(offset int64, whence int) (ret int64, err syscall.Error) {
	return syscall.Seek(f.fd, offset, whence)
}

// Truncate changes the size of the named file.
// If the file is a symbolic link, it changes the size of the link's target.
func Truncate(name string, size int64) Error {
	var d Dir
	d.Null()

	d.Length = uint64(size)

	if e := syscall.Wstat(name, pdir(nil, &d)); iserror(e) {
		return &PathError{"truncate", name, e}
	}
	return nil
}

// Remove removes the named file or directory.
func Remove(name string) Error {
	if e := syscall.Remove(name); iserror(e) {
		return &PathError{"remove", name, e}
	}
	return nil
}

// Rename renames a file.
func Rename(oldname, newname string) Error {
	var d Dir
	d.Null()

	d.Name = newname

	if e := syscall.Wstat(oldname, pdir(nil, &d)); iserror(e) {
		return &PathError{"rename", oldname, e}
	}
	return nil
}

// Chmod changes the mode of the named file to mode.
func Chmod(name string, mode uint32) Error {
	var d Dir
	var mask = ^uint32(0777)

	d.Null()
	odir, e := dirstat(name)
	if iserror(e) {
		return &PathError{"chmod", name, e}
	}

	d.Mode = (odir.Mode & mask) | (mode &^ mask)
	if e := syscall.Wstat(name, pdir(nil, &d)); iserror(e) {
		return &PathError{"chmod", name, e}
	}
	return nil
}

// ChownPlan9 changes the uid and gid strings of the named file.
func ChownPlan9(name, uid, gid string) Error {
	var d Dir
	d.Null()

	d.Uid = uid
	d.Gid = gid

	if e := syscall.Wstat(name, pdir(nil, &d)); iserror(e) {
		return &PathError{"chown_plan9", name, e}
	}
	return nil
}

// Chtimes changes the access and modification times of the named
// file, similar to the Unix utime() or utimes() functions.
//
// The argument times are in nanoseconds, although the underlying
// filesystem may truncate or round the values to a more
// coarse time unit.
func Chtimes(name string, atimeNs int64, mtimeNs int64) Error {
	var d Dir
	d.Null()

	d.Atime = uint32(atimeNs / 1e9)
	d.Mtime = uint32(mtimeNs / 1e9)

	if e := syscall.Wstat(name, pdir(nil, &d)); iserror(e) {
		return &PathError{"chtimes", name, e}
	}
	return nil
}

func Pipe() (r *File, w *File, err Error) {
	var p [2]int

	syscall.ForkLock.RLock()
	if e := syscall.Pipe(p[0:]); iserror(e) {
		syscall.ForkLock.RUnlock()
		return nil, nil, NewSyscallError("pipe", e)
	}
	syscall.ForkLock.RUnlock()

	return NewFile(p[0], "|0"), NewFile(p[1], "|1"), nil
}


// not supported on Plan 9

// Link creates a hard link.
func Link(oldname, newname string) Error {
	return EPLAN9
}

func Symlink(oldname, newname string) Error {
	return EPLAN9
}

func Readlink(name string) (string, Error) {
	return "", EPLAN9
}

func Chown(name string, uid, gid int) Error {
	return EPLAN9
}

func Lchown(name string, uid, gid int) Error {
	return EPLAN9
}

func (f *File) Chown(uid, gid int) Error {
	return EPLAN9
}
