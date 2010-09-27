// mksyscall.sh syscall_linux.go syscall_linux_amd64.go
// MACHINE GENERATED BY THE COMMAND ABOVE; DO NOT EDIT

package syscall

import "unsafe"

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func open(path string, mode int, perm uint32) (fd int, errno int) {
	r0, _, e1 := Syscall(SYS_OPEN, uintptr(unsafe.Pointer(StringBytePtr(path))), uintptr(mode), uintptr(perm))
	fd = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func openat(dirfd int, path string, flags int, mode uint32) (fd int, errno int) {
	r0, _, e1 := Syscall6(SYS_OPENAT, uintptr(dirfd), uintptr(unsafe.Pointer(StringBytePtr(path))), uintptr(flags), uintptr(mode), 0, 0)
	fd = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func pipe(p *[2]_C_int) (errno int) {
	_, _, e1 := Syscall(SYS_PIPE, uintptr(unsafe.Pointer(p)), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func utimes(path string, times *[2]Timeval) (errno int) {
	_, _, e1 := Syscall(SYS_UTIMES, uintptr(unsafe.Pointer(StringBytePtr(path))), uintptr(unsafe.Pointer(times)), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func futimesat(dirfd int, path *byte, times *[2]Timeval) (errno int) {
	_, _, e1 := Syscall(SYS_FUTIMESAT, uintptr(dirfd), uintptr(unsafe.Pointer(path)), uintptr(unsafe.Pointer(times)))
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Getcwd(buf []byte) (n int, errno int) {
	var _p0 unsafe.Pointer
	if len(buf) > 0 {
		_p0 = unsafe.Pointer(&buf[0])
	}
	r0, _, e1 := Syscall(SYS_GETCWD, uintptr(_p0), uintptr(len(buf)), 0)
	n = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func wait4(pid int, wstatus *_C_int, options int, rusage *Rusage) (wpid int, errno int) {
	r0, _, e1 := Syscall6(SYS_WAIT4, uintptr(pid), uintptr(unsafe.Pointer(wstatus)), uintptr(options), uintptr(unsafe.Pointer(rusage)), 0, 0)
	wpid = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func ptrace(request int, pid int, addr uintptr, data uintptr) (errno int) {
	_, _, e1 := Syscall6(SYS_PTRACE, uintptr(request), uintptr(pid), uintptr(addr), uintptr(data), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Access(path string, mode uint32) (errno int) {
	_, _, e1 := Syscall(SYS_ACCESS, uintptr(unsafe.Pointer(StringBytePtr(path))), uintptr(mode), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Acct(path string) (errno int) {
	_, _, e1 := Syscall(SYS_ACCT, uintptr(unsafe.Pointer(StringBytePtr(path))), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Adjtimex(buf *Timex) (state int, errno int) {
	r0, _, e1 := Syscall(SYS_ADJTIMEX, uintptr(unsafe.Pointer(buf)), 0, 0)
	state = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Chdir(path string) (errno int) {
	_, _, e1 := Syscall(SYS_CHDIR, uintptr(unsafe.Pointer(StringBytePtr(path))), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Chmod(path string, mode uint32) (errno int) {
	_, _, e1 := Syscall(SYS_CHMOD, uintptr(unsafe.Pointer(StringBytePtr(path))), uintptr(mode), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Chroot(path string) (errno int) {
	_, _, e1 := Syscall(SYS_CHROOT, uintptr(unsafe.Pointer(StringBytePtr(path))), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Close(fd int) (errno int) {
	_, _, e1 := Syscall(SYS_CLOSE, uintptr(fd), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Creat(path string, mode uint32) (fd int, errno int) {
	r0, _, e1 := Syscall(SYS_CREAT, uintptr(unsafe.Pointer(StringBytePtr(path))), uintptr(mode), 0)
	fd = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Dup(oldfd int) (fd int, errno int) {
	r0, _, e1 := Syscall(SYS_DUP, uintptr(oldfd), 0, 0)
	fd = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Dup2(oldfd int, newfd int) (fd int, errno int) {
	r0, _, e1 := Syscall(SYS_DUP2, uintptr(oldfd), uintptr(newfd), 0)
	fd = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func EpollCreate(size int) (fd int, errno int) {
	r0, _, e1 := Syscall(SYS_EPOLL_CREATE, uintptr(size), 0, 0)
	fd = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func EpollCtl(epfd int, op int, fd int, event *EpollEvent) (errno int) {
	_, _, e1 := Syscall6(SYS_EPOLL_CTL, uintptr(epfd), uintptr(op), uintptr(fd), uintptr(unsafe.Pointer(event)), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func EpollWait(epfd int, events []EpollEvent, msec int) (n int, errno int) {
	var _p0 unsafe.Pointer
	if len(events) > 0 {
		_p0 = unsafe.Pointer(&events[0])
	}
	r0, _, e1 := Syscall6(SYS_EPOLL_WAIT, uintptr(epfd), uintptr(_p0), uintptr(len(events)), uintptr(msec), 0, 0)
	n = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Exit(code int) {
	Syscall(SYS_EXIT_GROUP, uintptr(code), 0, 0)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Faccessat(dirfd int, path string, mode uint32, flags int) (errno int) {
	_, _, e1 := Syscall6(SYS_FACCESSAT, uintptr(dirfd), uintptr(unsafe.Pointer(StringBytePtr(path))), uintptr(mode), uintptr(flags), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Fallocate(fd int, mode uint32, off int64, len int64) (errno int) {
	_, _, e1 := Syscall6(SYS_FALLOCATE, uintptr(fd), uintptr(mode), uintptr(off), uintptr(len), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Fchdir(fd int) (errno int) {
	_, _, e1 := Syscall(SYS_FCHDIR, uintptr(fd), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Fchmod(fd int, mode uint32) (errno int) {
	_, _, e1 := Syscall(SYS_FCHMOD, uintptr(fd), uintptr(mode), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Fchmodat(dirfd int, path string, mode uint32, flags int) (errno int) {
	_, _, e1 := Syscall6(SYS_FCHMODAT, uintptr(dirfd), uintptr(unsafe.Pointer(StringBytePtr(path))), uintptr(mode), uintptr(flags), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Fchownat(dirfd int, path string, uid int, gid int, flags int) (errno int) {
	_, _, e1 := Syscall6(SYS_FCHOWNAT, uintptr(dirfd), uintptr(unsafe.Pointer(StringBytePtr(path))), uintptr(uid), uintptr(gid), uintptr(flags), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func fcntl(fd int, cmd int, arg int) (val int, errno int) {
	r0, _, e1 := Syscall(SYS_FCNTL, uintptr(fd), uintptr(cmd), uintptr(arg))
	val = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Fdatasync(fd int) (errno int) {
	_, _, e1 := Syscall(SYS_FDATASYNC, uintptr(fd), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Fsync(fd int) (errno int) {
	_, _, e1 := Syscall(SYS_FSYNC, uintptr(fd), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Ftruncate(fd int, length int64) (errno int) {
	_, _, e1 := Syscall(SYS_FTRUNCATE, uintptr(fd), uintptr(length), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Getdents(fd int, buf []byte) (n int, errno int) {
	var _p0 unsafe.Pointer
	if len(buf) > 0 {
		_p0 = unsafe.Pointer(&buf[0])
	}
	r0, _, e1 := Syscall(SYS_GETDENTS64, uintptr(fd), uintptr(_p0), uintptr(len(buf)))
	n = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Getpgid(pid int) (pgid int, errno int) {
	r0, _, e1 := Syscall(SYS_GETPGID, uintptr(pid), 0, 0)
	pgid = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Getpgrp() (pid int) {
	r0, _, _ := Syscall(SYS_GETPGRP, 0, 0, 0)
	pid = int(r0)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Getpid() (pid int) {
	r0, _, _ := Syscall(SYS_GETPID, 0, 0, 0)
	pid = int(r0)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Getppid() (ppid int) {
	r0, _, _ := Syscall(SYS_GETPPID, 0, 0, 0)
	ppid = int(r0)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Getrlimit(resource int, rlim *Rlimit) (errno int) {
	_, _, e1 := Syscall(SYS_GETRLIMIT, uintptr(resource), uintptr(unsafe.Pointer(rlim)), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Getrusage(who int, rusage *Rusage) (errno int) {
	_, _, e1 := Syscall(SYS_GETRUSAGE, uintptr(who), uintptr(unsafe.Pointer(rusage)), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Gettid() (tid int) {
	r0, _, _ := Syscall(SYS_GETTID, 0, 0, 0)
	tid = int(r0)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func InotifyAddWatch(fd int, pathname string, mask uint32) (watchdesc int, errno int) {
	r0, _, e1 := Syscall(SYS_INOTIFY_ADD_WATCH, uintptr(fd), uintptr(unsafe.Pointer(StringBytePtr(pathname))), uintptr(mask))
	watchdesc = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func InotifyInit() (fd int, errno int) {
	r0, _, e1 := Syscall(SYS_INOTIFY_INIT, 0, 0, 0)
	fd = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func InotifyInit1(flags int) (fd int, errno int) {
	r0, _, e1 := Syscall(SYS_INOTIFY_INIT1, uintptr(flags), 0, 0)
	fd = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func InotifyRmWatch(fd int, watchdesc uint32) (success int, errno int) {
	r0, _, e1 := Syscall(SYS_INOTIFY_RM_WATCH, uintptr(fd), uintptr(watchdesc), 0)
	success = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Kill(pid int, sig int) (errno int) {
	_, _, e1 := Syscall(SYS_KILL, uintptr(pid), uintptr(sig), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Klogctl(typ int, buf []byte) (n int, errno int) {
	var _p0 unsafe.Pointer
	if len(buf) > 0 {
		_p0 = unsafe.Pointer(&buf[0])
	}
	r0, _, e1 := Syscall(SYS_SYSLOG, uintptr(typ), uintptr(_p0), uintptr(len(buf)))
	n = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Link(oldpath string, newpath string) (errno int) {
	_, _, e1 := Syscall(SYS_LINK, uintptr(unsafe.Pointer(StringBytePtr(oldpath))), uintptr(unsafe.Pointer(StringBytePtr(newpath))), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Mkdir(path string, mode uint32) (errno int) {
	_, _, e1 := Syscall(SYS_MKDIR, uintptr(unsafe.Pointer(StringBytePtr(path))), uintptr(mode), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Mkdirat(dirfd int, path string, mode uint32) (errno int) {
	_, _, e1 := Syscall(SYS_MKDIRAT, uintptr(dirfd), uintptr(unsafe.Pointer(StringBytePtr(path))), uintptr(mode))
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Mknod(path string, mode uint32, dev int) (errno int) {
	_, _, e1 := Syscall(SYS_MKNOD, uintptr(unsafe.Pointer(StringBytePtr(path))), uintptr(mode), uintptr(dev))
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Mknodat(dirfd int, path string, mode uint32, dev int) (errno int) {
	_, _, e1 := Syscall6(SYS_MKNODAT, uintptr(dirfd), uintptr(unsafe.Pointer(StringBytePtr(path))), uintptr(mode), uintptr(dev), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Nanosleep(time *Timespec, leftover *Timespec) (errno int) {
	_, _, e1 := Syscall(SYS_NANOSLEEP, uintptr(unsafe.Pointer(time)), uintptr(unsafe.Pointer(leftover)), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Pause() (errno int) {
	_, _, e1 := Syscall(SYS_PAUSE, 0, 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func PivotRoot(newroot string, putold string) (errno int) {
	_, _, e1 := Syscall(SYS_PIVOT_ROOT, uintptr(unsafe.Pointer(StringBytePtr(newroot))), uintptr(unsafe.Pointer(StringBytePtr(putold))), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Pread(fd int, p []byte, offset int64) (n int, errno int) {
	var _p0 unsafe.Pointer
	if len(p) > 0 {
		_p0 = unsafe.Pointer(&p[0])
	}
	r0, _, e1 := Syscall6(SYS_PREAD64, uintptr(fd), uintptr(_p0), uintptr(len(p)), uintptr(offset), 0, 0)
	n = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Pwrite(fd int, p []byte, offset int64) (n int, errno int) {
	var _p0 unsafe.Pointer
	if len(p) > 0 {
		_p0 = unsafe.Pointer(&p[0])
	}
	r0, _, e1 := Syscall6(SYS_PWRITE64, uintptr(fd), uintptr(_p0), uintptr(len(p)), uintptr(offset), 0, 0)
	n = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Read(fd int, p []byte) (n int, errno int) {
	var _p0 unsafe.Pointer
	if len(p) > 0 {
		_p0 = unsafe.Pointer(&p[0])
	}
	r0, _, e1 := Syscall(SYS_READ, uintptr(fd), uintptr(_p0), uintptr(len(p)))
	n = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Readlink(path string, buf []byte) (n int, errno int) {
	var _p0 unsafe.Pointer
	if len(buf) > 0 {
		_p0 = unsafe.Pointer(&buf[0])
	}
	r0, _, e1 := Syscall(SYS_READLINK, uintptr(unsafe.Pointer(StringBytePtr(path))), uintptr(_p0), uintptr(len(buf)))
	n = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Rename(oldpath string, newpath string) (errno int) {
	_, _, e1 := Syscall(SYS_RENAME, uintptr(unsafe.Pointer(StringBytePtr(oldpath))), uintptr(unsafe.Pointer(StringBytePtr(newpath))), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Renameat(olddirfd int, oldpath string, newdirfd int, newpath string) (errno int) {
	_, _, e1 := Syscall6(SYS_RENAMEAT, uintptr(olddirfd), uintptr(unsafe.Pointer(StringBytePtr(oldpath))), uintptr(newdirfd), uintptr(unsafe.Pointer(StringBytePtr(newpath))), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Rmdir(path string) (errno int) {
	_, _, e1 := Syscall(SYS_RMDIR, uintptr(unsafe.Pointer(StringBytePtr(path))), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Setdomainname(p []byte) (errno int) {
	var _p0 unsafe.Pointer
	if len(p) > 0 {
		_p0 = unsafe.Pointer(&p[0])
	}
	_, _, e1 := Syscall(SYS_SETDOMAINNAME, uintptr(_p0), uintptr(len(p)), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Sethostname(p []byte) (errno int) {
	var _p0 unsafe.Pointer
	if len(p) > 0 {
		_p0 = unsafe.Pointer(&p[0])
	}
	_, _, e1 := Syscall(SYS_SETHOSTNAME, uintptr(_p0), uintptr(len(p)), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Setpgid(pid int, pgid int) (errno int) {
	_, _, e1 := Syscall(SYS_SETPGID, uintptr(pid), uintptr(pgid), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Setrlimit(resource int, rlim *Rlimit) (errno int) {
	_, _, e1 := Syscall(SYS_SETRLIMIT, uintptr(resource), uintptr(unsafe.Pointer(rlim)), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Setsid() (pid int, errno int) {
	r0, _, e1 := Syscall(SYS_SETSID, 0, 0, 0)
	pid = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Settimeofday(tv *Timeval) (errno int) {
	_, _, e1 := Syscall(SYS_SETTIMEOFDAY, uintptr(unsafe.Pointer(tv)), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Setuid(uid int) (errno int) {
	_, _, e1 := Syscall(SYS_SETUID, uintptr(uid), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Splice(rfd int, roff *int64, wfd int, woff *int64, len int, flags int) (n int64, errno int) {
	r0, _, e1 := Syscall6(SYS_SPLICE, uintptr(rfd), uintptr(unsafe.Pointer(roff)), uintptr(wfd), uintptr(unsafe.Pointer(woff)), uintptr(len), uintptr(flags))
	n = int64(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Symlink(oldpath string, newpath string) (errno int) {
	_, _, e1 := Syscall(SYS_SYMLINK, uintptr(unsafe.Pointer(StringBytePtr(oldpath))), uintptr(unsafe.Pointer(StringBytePtr(newpath))), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Sync() {
	Syscall(SYS_SYNC, 0, 0, 0)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Sysinfo(info *Sysinfo_t) (errno int) {
	_, _, e1 := Syscall(SYS_SYSINFO, uintptr(unsafe.Pointer(info)), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Tee(rfd int, wfd int, len int, flags int) (n int64, errno int) {
	r0, _, e1 := Syscall6(SYS_TEE, uintptr(rfd), uintptr(wfd), uintptr(len), uintptr(flags), 0, 0)
	n = int64(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Tgkill(tgid int, tid int, sig int) (errno int) {
	_, _, e1 := Syscall(SYS_TGKILL, uintptr(tgid), uintptr(tid), uintptr(sig))
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Times(tms *Tms) (ticks uintptr, errno int) {
	r0, _, e1 := Syscall(SYS_TIMES, uintptr(unsafe.Pointer(tms)), 0, 0)
	ticks = uintptr(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Truncate(path string, length int64) (errno int) {
	_, _, e1 := Syscall(SYS_TRUNCATE, uintptr(unsafe.Pointer(StringBytePtr(path))), uintptr(length), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Umask(mask int) (oldmask int) {
	r0, _, _ := Syscall(SYS_UMASK, uintptr(mask), 0, 0)
	oldmask = int(r0)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Uname(buf *Utsname) (errno int) {
	_, _, e1 := Syscall(SYS_UNAME, uintptr(unsafe.Pointer(buf)), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Unlink(path string) (errno int) {
	_, _, e1 := Syscall(SYS_UNLINK, uintptr(unsafe.Pointer(StringBytePtr(path))), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Unlinkat(dirfd int, path string) (errno int) {
	_, _, e1 := Syscall(SYS_UNLINKAT, uintptr(dirfd), uintptr(unsafe.Pointer(StringBytePtr(path))), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Unshare(flags int) (errno int) {
	_, _, e1 := Syscall(SYS_UNSHARE, uintptr(flags), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Ustat(dev int, ubuf *Ustat_t) (errno int) {
	_, _, e1 := Syscall(SYS_USTAT, uintptr(dev), uintptr(unsafe.Pointer(ubuf)), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Utime(path string, buf *Utimbuf) (errno int) {
	_, _, e1 := Syscall(SYS_UTIME, uintptr(unsafe.Pointer(StringBytePtr(path))), uintptr(unsafe.Pointer(buf)), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Write(fd int, p []byte) (n int, errno int) {
	var _p0 unsafe.Pointer
	if len(p) > 0 {
		_p0 = unsafe.Pointer(&p[0])
	}
	r0, _, e1 := Syscall(SYS_WRITE, uintptr(fd), uintptr(_p0), uintptr(len(p)))
	n = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func exitThread(code int) (errno int) {
	_, _, e1 := Syscall(SYS_EXIT, uintptr(code), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func read(fd int, p *byte, np int) (n int, errno int) {
	r0, _, e1 := Syscall(SYS_READ, uintptr(fd), uintptr(unsafe.Pointer(p)), uintptr(np))
	n = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func write(fd int, p *byte, np int) (n int, errno int) {
	r0, _, e1 := Syscall(SYS_WRITE, uintptr(fd), uintptr(unsafe.Pointer(p)), uintptr(np))
	n = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Chown(path string, uid int, gid int) (errno int) {
	_, _, e1 := Syscall(SYS_CHOWN, uintptr(unsafe.Pointer(StringBytePtr(path))), uintptr(uid), uintptr(gid))
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Fchown(fd int, uid int, gid int) (errno int) {
	_, _, e1 := Syscall(SYS_FCHOWN, uintptr(fd), uintptr(uid), uintptr(gid))
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Fstat(fd int, stat *Stat_t) (errno int) {
	_, _, e1 := Syscall(SYS_FSTAT, uintptr(fd), uintptr(unsafe.Pointer(stat)), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Fstatfs(fd int, buf *Statfs_t) (errno int) {
	_, _, e1 := Syscall(SYS_FSTATFS, uintptr(fd), uintptr(unsafe.Pointer(buf)), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Getegid() (egid int) {
	r0, _, _ := Syscall(SYS_GETEGID, 0, 0, 0)
	egid = int(r0)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Geteuid() (euid int) {
	r0, _, _ := Syscall(SYS_GETEUID, 0, 0, 0)
	euid = int(r0)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Getgid() (gid int) {
	r0, _, _ := Syscall(SYS_GETGID, 0, 0, 0)
	gid = int(r0)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Getuid() (uid int) {
	r0, _, _ := Syscall(SYS_GETUID, 0, 0, 0)
	uid = int(r0)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Ioperm(from int, num int, on int) (errno int) {
	_, _, e1 := Syscall(SYS_IOPERM, uintptr(from), uintptr(num), uintptr(on))
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Iopl(level int) (errno int) {
	_, _, e1 := Syscall(SYS_IOPL, uintptr(level), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Lchown(path string, uid int, gid int) (errno int) {
	_, _, e1 := Syscall(SYS_LCHOWN, uintptr(unsafe.Pointer(StringBytePtr(path))), uintptr(uid), uintptr(gid))
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Listen(s int, n int) (errno int) {
	_, _, e1 := Syscall(SYS_LISTEN, uintptr(s), uintptr(n), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Lstat(path string, stat *Stat_t) (errno int) {
	_, _, e1 := Syscall(SYS_LSTAT, uintptr(unsafe.Pointer(StringBytePtr(path))), uintptr(unsafe.Pointer(stat)), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Seek(fd int, offset int64, whence int) (off int64, errno int) {
	r0, _, e1 := Syscall(SYS_LSEEK, uintptr(fd), uintptr(offset), uintptr(whence))
	off = int64(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Select(nfd int, r *FdSet, w *FdSet, e *FdSet, timeout *Timeval) (n int, errno int) {
	r0, _, e1 := Syscall6(SYS_SELECT, uintptr(nfd), uintptr(unsafe.Pointer(r)), uintptr(unsafe.Pointer(w)), uintptr(unsafe.Pointer(e)), uintptr(unsafe.Pointer(timeout)), 0)
	n = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Setfsgid(gid int) (errno int) {
	_, _, e1 := Syscall(SYS_SETFSGID, uintptr(gid), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Setfsuid(uid int) (errno int) {
	_, _, e1 := Syscall(SYS_SETFSUID, uintptr(uid), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Setgid(gid int) (errno int) {
	_, _, e1 := Syscall(SYS_SETGID, uintptr(gid), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Setregid(rgid int, egid int) (errno int) {
	_, _, e1 := Syscall(SYS_SETREGID, uintptr(rgid), uintptr(egid), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Setresgid(rgid int, egid int, sgid int) (errno int) {
	_, _, e1 := Syscall(SYS_SETRESGID, uintptr(rgid), uintptr(egid), uintptr(sgid))
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Setresuid(ruid int, euid int, suid int) (errno int) {
	_, _, e1 := Syscall(SYS_SETRESUID, uintptr(ruid), uintptr(euid), uintptr(suid))
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Setreuid(ruid int, euid int) (errno int) {
	_, _, e1 := Syscall(SYS_SETREUID, uintptr(ruid), uintptr(euid), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Shutdown(fd int, how int) (errno int) {
	_, _, e1 := Syscall(SYS_SHUTDOWN, uintptr(fd), uintptr(how), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Stat(path string, stat *Stat_t) (errno int) {
	_, _, e1 := Syscall(SYS_STAT, uintptr(unsafe.Pointer(StringBytePtr(path))), uintptr(unsafe.Pointer(stat)), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Statfs(path string, buf *Statfs_t) (errno int) {
	_, _, e1 := Syscall(SYS_STATFS, uintptr(unsafe.Pointer(StringBytePtr(path))), uintptr(unsafe.Pointer(buf)), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func SyncFileRange(fd int, off int64, n int64, flags int) (errno int) {
	_, _, e1 := Syscall6(SYS_SYNC_FILE_RANGE, uintptr(fd), uintptr(off), uintptr(n), uintptr(flags), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func accept(s int, rsa *RawSockaddrAny, addrlen *_Socklen) (fd int, errno int) {
	r0, _, e1 := Syscall(SYS_ACCEPT, uintptr(s), uintptr(unsafe.Pointer(rsa)), uintptr(unsafe.Pointer(addrlen)))
	fd = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func bind(s int, addr uintptr, addrlen _Socklen) (errno int) {
	_, _, e1 := Syscall(SYS_BIND, uintptr(s), uintptr(addr), uintptr(addrlen))
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func connect(s int, addr uintptr, addrlen _Socklen) (errno int) {
	_, _, e1 := Syscall(SYS_CONNECT, uintptr(s), uintptr(addr), uintptr(addrlen))
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func getgroups(n int, list *_Gid_t) (nn int, errno int) {
	r0, _, e1 := Syscall(SYS_GETGROUPS, uintptr(n), uintptr(unsafe.Pointer(list)), 0)
	nn = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func setgroups(n int, list *_Gid_t) (errno int) {
	_, _, e1 := Syscall(SYS_SETGROUPS, uintptr(n), uintptr(unsafe.Pointer(list)), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func setsockopt(s int, level int, name int, val uintptr, vallen int) (errno int) {
	_, _, e1 := Syscall6(SYS_SETSOCKOPT, uintptr(s), uintptr(level), uintptr(name), uintptr(val), uintptr(vallen), 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func socket(domain int, typ int, proto int) (fd int, errno int) {
	r0, _, e1 := Syscall(SYS_SOCKET, uintptr(domain), uintptr(typ), uintptr(proto))
	fd = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func socketpair(domain int, typ int, proto int, fd *[2]int) (errno int) {
	_, _, e1 := Syscall6(SYS_SOCKETPAIR, uintptr(domain), uintptr(typ), uintptr(proto), uintptr(unsafe.Pointer(fd)), 0, 0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func getpeername(fd int, rsa *RawSockaddrAny, addrlen *_Socklen) (errno int) {
	_, _, e1 := Syscall(SYS_GETPEERNAME, uintptr(fd), uintptr(unsafe.Pointer(rsa)), uintptr(unsafe.Pointer(addrlen)))
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func getsockname(fd int, rsa *RawSockaddrAny, addrlen *_Socklen) (errno int) {
	_, _, e1 := Syscall(SYS_GETSOCKNAME, uintptr(fd), uintptr(unsafe.Pointer(rsa)), uintptr(unsafe.Pointer(addrlen)))
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func recvfrom(fd int, p []byte, flags int, from *RawSockaddrAny, fromlen *_Socklen) (n int, errno int) {
	var _p0 unsafe.Pointer
	if len(p) > 0 {
		_p0 = unsafe.Pointer(&p[0])
	}
	r0, _, e1 := Syscall6(SYS_RECVFROM, uintptr(fd), uintptr(_p0), uintptr(len(p)), uintptr(flags), uintptr(unsafe.Pointer(from)), uintptr(unsafe.Pointer(fromlen)))
	n = int(r0)
	errno = int(e1)
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func sendto(s int, buf []byte, flags int, to uintptr, addrlen _Socklen) (errno int) {
	var _p0 unsafe.Pointer
	if len(buf) > 0 {
		_p0 = unsafe.Pointer(&buf[0])
	}
	_, _, e1 := Syscall6(SYS_SENDTO, uintptr(s), uintptr(_p0), uintptr(len(buf)), uintptr(flags), uintptr(to), uintptr(addrlen))
	errno = int(e1)
	return
}
