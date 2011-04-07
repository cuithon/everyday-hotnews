// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Linux system calls.
// This file is compiled as ordinary Go code,
// but it is also input to mksyscall,
// which parses the //sys lines and generates system call stubs.
// Note that sometimes we use a lowercase //sys name and
// wrap it in our own nicer implementation.

package syscall

import "unsafe"

const OS = "linux"

/*
 * Wrapped
 */

//sys	open(path string, mode int, perm uint32) (fd int, errno int)
func Open(path string, mode int, perm uint32) (fd int, errno int) {
	return open(path, mode|O_LARGEFILE, perm)
}

//sys	openat(dirfd int, path string, flags int, mode uint32) (fd int, errno int)
func Openat(dirfd int, path string, flags int, mode uint32) (fd int, errno int) {
	return openat(dirfd, path, flags|O_LARGEFILE, mode)
}

//sysnb	pipe(p *[2]_C_int) (errno int)
func Pipe(p []int) (errno int) {
	if len(p) != 2 {
		return EINVAL
	}
	var pp [2]_C_int
	errno = pipe(&pp)
	p[0] = int(pp[0])
	p[1] = int(pp[1])
	return
}

//sys	utimes(path string, times *[2]Timeval) (errno int)
func Utimes(path string, tv []Timeval) (errno int) {
	if len(tv) != 2 {
		return EINVAL
	}
	return utimes(path, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
}

//sys	futimesat(dirfd int, path *byte, times *[2]Timeval) (errno int)
func Futimesat(dirfd int, path string, tv []Timeval) (errno int) {
	if len(tv) != 2 {
		return EINVAL
	}
	return futimesat(dirfd, StringBytePtr(path), (*[2]Timeval)(unsafe.Pointer(&tv[0])))
}

func Futimes(fd int, tv []Timeval) (errno int) {
	// Believe it or not, this is the best we can do on Linux
	// (and is what glibc does).
	return Utimes("/proc/self/fd/"+itoa(fd), tv)
}

const ImplementsGetwd = true

//sys	Getcwd(buf []byte) (n int, errno int)
func Getwd() (wd string, errno int) {
	var buf [PathMax]byte
	n, err := Getcwd(buf[0:])
	if err != 0 {
		return "", err
	}
	// Getcwd returns the number of bytes written to buf, including the NUL.
	if n < 1 || n > len(buf) || buf[n-1] != 0 {
		return "", EINVAL
	}
	return string(buf[0 : n-1]), 0
}

func Getgroups() (gids []int, errno int) {
	n, err := getgroups(0, nil)
	if err != 0 {
		return nil, errno
	}
	if n == 0 {
		return nil, 0
	}

	// Sanity check group count.  Max is 1<<16 on Linux.
	if n < 0 || n > 1<<20 {
		return nil, EINVAL
	}

	a := make([]_Gid_t, n)
	n, err = getgroups(n, &a[0])
	if err != 0 {
		return nil, errno
	}
	gids = make([]int, n)
	for i, v := range a[0:n] {
		gids[i] = int(v)
	}
	return
}

func Setgroups(gids []int) (errno int) {
	if len(gids) == 0 {
		return setgroups(0, nil)
	}

	a := make([]_Gid_t, len(gids))
	for i, v := range gids {
		a[i] = _Gid_t(v)
	}
	return setgroups(len(a), &a[0])
}

type WaitStatus uint32

// Wait status is 7 bits at bottom, either 0 (exited),
// 0x7F (stopped), or a signal number that caused an exit.
// The 0x80 bit is whether there was a core dump.
// An extra number (exit code, signal causing a stop)
// is in the high bits.  At least that's the idea.
// There are various irregularities.  For example, the
// "continued" status is 0xFFFF, distinguishing itself
// from stopped via the core dump bit.

const (
	mask    = 0x7F
	core    = 0x80
	exited  = 0x00
	stopped = 0x7F
	shift   = 8
)

func (w WaitStatus) Exited() bool { return w&mask == exited }

func (w WaitStatus) Signaled() bool { return w&mask != stopped && w&mask != exited }

func (w WaitStatus) Stopped() bool { return w&0xFF == stopped }

func (w WaitStatus) Continued() bool { return w == 0xFFFF }

func (w WaitStatus) CoreDump() bool { return w.Signaled() && w&core != 0 }

func (w WaitStatus) ExitStatus() int {
	if !w.Exited() {
		return -1
	}
	return int(w>>shift) & 0xFF
}

func (w WaitStatus) Signal() int {
	if !w.Signaled() {
		return -1
	}
	return int(w & mask)
}

func (w WaitStatus) StopSignal() int {
	if !w.Stopped() {
		return -1
	}
	return int(w>>shift) & 0xFF
}

func (w WaitStatus) TrapCause() int {
	if w.StopSignal() != SIGTRAP {
		return -1
	}
	return int(w>>shift) >> 8
}

//sys	wait4(pid int, wstatus *_C_int, options int, rusage *Rusage) (wpid int, errno int)
func Wait4(pid int, wstatus *WaitStatus, options int, rusage *Rusage) (wpid int, errno int) {
	var status _C_int
	wpid, errno = wait4(pid, &status, options, rusage)
	if wstatus != nil {
		*wstatus = WaitStatus(status)
	}
	return
}

func Sleep(nsec int64) (errno int) {
	tv := NsecToTimeval(nsec)
	_, err := Select(0, nil, nil, nil, &tv)
	return err
}

// For testing: clients can set this flag to force
// creation of IPv6 sockets to return EAFNOSUPPORT.
var SocketDisableIPv6 bool

type Sockaddr interface {
	sockaddr() (ptr uintptr, len _Socklen, errno int) // lowercase; only we can define Sockaddrs
}

type SockaddrInet4 struct {
	Port int
	Addr [4]byte
	raw  RawSockaddrInet4
}

func (sa *SockaddrInet4) sockaddr() (uintptr, _Socklen, int) {
	if sa.Port < 0 || sa.Port > 0xFFFF {
		return 0, 0, EINVAL
	}
	sa.raw.Family = AF_INET
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	for i := 0; i < len(sa.Addr); i++ {
		sa.raw.Addr[i] = sa.Addr[i]
	}
	return uintptr(unsafe.Pointer(&sa.raw)), SizeofSockaddrInet4, 0
}

type SockaddrInet6 struct {
	Port int
	Addr [16]byte
	raw  RawSockaddrInet6
}

func (sa *SockaddrInet6) sockaddr() (uintptr, _Socklen, int) {
	if sa.Port < 0 || sa.Port > 0xFFFF {
		return 0, 0, EINVAL
	}
	sa.raw.Family = AF_INET6
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	for i := 0; i < len(sa.Addr); i++ {
		sa.raw.Addr[i] = sa.Addr[i]
	}
	return uintptr(unsafe.Pointer(&sa.raw)), SizeofSockaddrInet6, 0
}

type SockaddrUnix struct {
	Name string
	raw  RawSockaddrUnix
}

func (sa *SockaddrUnix) sockaddr() (uintptr, _Socklen, int) {
	name := sa.Name
	n := len(name)
	if n >= len(sa.raw.Path) || n == 0 {
		return 0, 0, EINVAL
	}
	sa.raw.Family = AF_UNIX
	for i := 0; i < n; i++ {
		sa.raw.Path[i] = int8(name[i])
	}
	// length is family (uint16), name, NUL.
	sl := 2 + _Socklen(n) + 1
	if sa.raw.Path[0] == '@' {
		sa.raw.Path[0] = 0
		// Don't count trailing NUL for abstract address.
		sl--
	}

	return uintptr(unsafe.Pointer(&sa.raw)), sl, 0
}

type SockaddrLinklayer struct {
	Protocol uint16
	Ifindex  int
	Hatype   uint16
	Pkttype  uint8
	Halen    uint8
	Addr     [8]byte
	raw      RawSockaddrLinklayer
}

func (sa *SockaddrLinklayer) sockaddr() (uintptr, _Socklen, int) {
	if sa.Ifindex < 0 || sa.Ifindex > 0x7fffffff {
		return 0, 0, EINVAL
	}
	sa.raw.Family = AF_PACKET
	sa.raw.Protocol = sa.Protocol
	sa.raw.Ifindex = int32(sa.Ifindex)
	sa.raw.Hatype = sa.Hatype
	sa.raw.Pkttype = sa.Pkttype
	sa.raw.Halen = sa.Halen
	for i := 0; i < len(sa.Addr); i++ {
		sa.raw.Addr[i] = sa.Addr[i]
	}
	return uintptr(unsafe.Pointer(&sa.raw)), SizeofSockaddrLinklayer, 0
}

func anyToSockaddr(rsa *RawSockaddrAny) (Sockaddr, int) {
	switch rsa.Addr.Family {
	case AF_PACKET:
		pp := (*RawSockaddrLinklayer)(unsafe.Pointer(rsa))
		sa := new(SockaddrLinklayer)
		sa.Protocol = pp.Protocol
		sa.Ifindex = int(pp.Ifindex)
		sa.Hatype = pp.Hatype
		sa.Pkttype = pp.Pkttype
		sa.Halen = pp.Halen
		for i := 0; i < len(sa.Addr); i++ {
			sa.Addr[i] = pp.Addr[i]
		}
		return sa, 0

	case AF_UNIX:
		pp := (*RawSockaddrUnix)(unsafe.Pointer(rsa))
		sa := new(SockaddrUnix)
		if pp.Path[0] == 0 {
			// "Abstract" Unix domain socket.
			// Rewrite leading NUL as @ for textual display.
			// (This is the standard convention.)
			// Not friendly to overwrite in place,
			// but the callers below don't care.
			pp.Path[0] = '@'
		}

		// Assume path ends at NUL.
		// This is not technically the Linux semantics for
		// abstract Unix domain sockets--they are supposed
		// to be uninterpreted fixed-size binary blobs--but
		// everyone uses this convention.
		n := 0
		for n < len(pp.Path) && pp.Path[n] != 0 {
			n++
		}
		bytes := (*[10000]byte)(unsafe.Pointer(&pp.Path[0]))[0:n]
		sa.Name = string(bytes)
		return sa, 0

	case AF_INET:
		pp := (*RawSockaddrInet4)(unsafe.Pointer(rsa))
		sa := new(SockaddrInet4)
		p := (*[2]byte)(unsafe.Pointer(&pp.Port))
		sa.Port = int(p[0])<<8 + int(p[1])
		for i := 0; i < len(sa.Addr); i++ {
			sa.Addr[i] = pp.Addr[i]
		}
		return sa, 0

	case AF_INET6:
		pp := (*RawSockaddrInet6)(unsafe.Pointer(rsa))
		sa := new(SockaddrInet6)
		p := (*[2]byte)(unsafe.Pointer(&pp.Port))
		sa.Port = int(p[0])<<8 + int(p[1])
		for i := 0; i < len(sa.Addr); i++ {
			sa.Addr[i] = pp.Addr[i]
		}
		return sa, 0
	}
	return nil, EAFNOSUPPORT
}

func Accept(fd int) (nfd int, sa Sockaddr, errno int) {
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	nfd, errno = accept(fd, &rsa, &len)
	if errno != 0 {
		return
	}
	sa, errno = anyToSockaddr(&rsa)
	if errno != 0 {
		Close(nfd)
		nfd = 0
	}
	return
}

func Getsockname(fd int) (sa Sockaddr, errno int) {
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	if errno = getsockname(fd, &rsa, &len); errno != 0 {
		return
	}
	return anyToSockaddr(&rsa)
}

func Getpeername(fd int) (sa Sockaddr, errno int) {
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	if errno = getpeername(fd, &rsa, &len); errno != 0 {
		return
	}
	return anyToSockaddr(&rsa)
}

func Bind(fd int, sa Sockaddr) (errno int) {
	ptr, n, err := sa.sockaddr()
	if err != 0 {
		return err
	}
	return bind(fd, ptr, n)
}

func Connect(fd int, sa Sockaddr) (errno int) {
	ptr, n, err := sa.sockaddr()
	if err != 0 {
		return err
	}
	return connect(fd, ptr, n)
}

func Socket(domain, typ, proto int) (fd, errno int) {
	if domain == AF_INET6 && SocketDisableIPv6 {
		return -1, EAFNOSUPPORT
	}
	fd, errno = socket(domain, typ, proto)
	return
}

func Socketpair(domain, typ, proto int) (fd [2]int, errno int) {
	errno = socketpair(domain, typ, proto, &fd)
	return
}

func GetsockoptInt(fd, level, opt int) (value, errno int) {
	var n int32
	vallen := _Socklen(4)
	errno = getsockopt(fd, level, opt, uintptr(unsafe.Pointer(&n)), &vallen)
	return int(n), errno
}

func SetsockoptInt(fd, level, opt int, value int) (errno int) {
	var n = int32(value)
	return setsockopt(fd, level, opt, uintptr(unsafe.Pointer(&n)), 4)
}

func SetsockoptTimeval(fd, level, opt int, tv *Timeval) (errno int) {
	return setsockopt(fd, level, opt, uintptr(unsafe.Pointer(tv)), unsafe.Sizeof(*tv))
}

func SetsockoptLinger(fd, level, opt int, l *Linger) (errno int) {
	return setsockopt(fd, level, opt, uintptr(unsafe.Pointer(l)), unsafe.Sizeof(*l))
}

func SetsockoptIpMreq(fd, level, opt int, mreq *IpMreq) (errno int) {
	return setsockopt(fd, level, opt, uintptr(unsafe.Pointer(mreq)), unsafe.Sizeof(*mreq))
}

func SetsockoptString(fd, level, opt int, s string) (errno int) {
	return setsockopt(fd, level, opt, uintptr(unsafe.Pointer(&[]byte(s)[0])), len(s))
}

func Recvfrom(fd int, p []byte, flags int) (n int, from Sockaddr, errno int) {
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	if n, errno = recvfrom(fd, p, flags, &rsa, &len); errno != 0 {
		return
	}
	from, errno = anyToSockaddr(&rsa)
	return
}

func Sendto(fd int, p []byte, flags int, to Sockaddr) (errno int) {
	ptr, n, err := to.sockaddr()
	if err != 0 {
		return err
	}
	return sendto(fd, p, flags, ptr, n)
}

func Recvmsg(fd int, p, oob []byte, flags int) (n, oobn int, recvflags int, from Sockaddr, errno int) {
	var msg Msghdr
	var rsa RawSockaddrAny
	msg.Name = (*byte)(unsafe.Pointer(&rsa))
	msg.Namelen = uint32(SizeofSockaddrAny)
	var iov Iovec
	if len(p) > 0 {
		iov.Base = (*byte)(unsafe.Pointer(&p[0]))
		iov.SetLen(len(p))
	}
	var dummy byte
	if len(oob) > 0 {
		// receive at least one normal byte
		if len(p) == 0 {
			iov.Base = &dummy
			iov.SetLen(1)
		}
		msg.Control = (*byte)(unsafe.Pointer(&oob[0]))
		msg.SetControllen(len(oob))
	}
	msg.Iov = &iov
	msg.Iovlen = 1
	if n, errno = recvmsg(fd, &msg, flags); errno != 0 {
		return
	}
	oobn = int(msg.Controllen)
	recvflags = int(msg.Flags)
	// source address is only specified if the socket is unconnected
	if rsa.Addr.Family != 0 {
		from, errno = anyToSockaddr(&rsa)
	}
	return
}

func Sendmsg(fd int, p, oob []byte, to Sockaddr, flags int) (errno int) {
	var ptr uintptr
	var nsock _Socklen
	if to != nil {
		var err int
		ptr, nsock, err = to.sockaddr()
		if err != 0 {
			return err
		}
	}
	var msg Msghdr
	msg.Name = (*byte)(unsafe.Pointer(ptr))
	msg.Namelen = uint32(nsock)
	var iov Iovec
	if len(p) > 0 {
		iov.Base = (*byte)(unsafe.Pointer(&p[0]))
		iov.SetLen(len(p))
	}
	var dummy byte
	if len(oob) > 0 {
		// send at least one normal byte
		if len(p) == 0 {
			iov.Base = &dummy
			iov.SetLen(1)
		}
		msg.Control = (*byte)(unsafe.Pointer(&oob[0]))
		msg.SetControllen(len(oob))
	}
	msg.Iov = &iov
	msg.Iovlen = 1
	if errno = sendmsg(fd, &msg, flags); errno != 0 {
		return
	}
	return
}

// BindToDevice binds the socket associated with fd to device.
func BindToDevice(fd int, device string) (errno int) {
	return SetsockoptString(fd, SOL_SOCKET, SO_BINDTODEVICE, device)
}

//sys	ptrace(request int, pid int, addr uintptr, data uintptr) (errno int)

func ptracePeek(req int, pid int, addr uintptr, out []byte) (count int, errno int) {
	// The peek requests are machine-size oriented, so we wrap it
	// to retrieve arbitrary-length data.

	// The ptrace syscall differs from glibc's ptrace.
	// Peeks returns the word in *data, not as the return value.

	var buf [sizeofPtr]byte

	// Leading edge.  PEEKTEXT/PEEKDATA don't require aligned
	// access (PEEKUSER warns that it might), but if we don't
	// align our reads, we might straddle an unmapped page
	// boundary and not get the bytes leading up to the page
	// boundary.
	n := 0
	if addr%sizeofPtr != 0 {
		errno = ptrace(req, pid, addr-addr%sizeofPtr, uintptr(unsafe.Pointer(&buf[0])))
		if errno != 0 {
			return 0, errno
		}
		n += copy(out, buf[addr%sizeofPtr:])
		out = out[n:]
	}

	// Remainder.
	for len(out) > 0 {
		// We use an internal buffer to gaurantee alignment.
		// It's not documented if this is necessary, but we're paranoid.
		errno = ptrace(req, pid, addr+uintptr(n), uintptr(unsafe.Pointer(&buf[0])))
		if errno != 0 {
			return n, errno
		}
		copied := copy(out, buf[0:])
		n += copied
		out = out[copied:]
	}

	return n, 0
}

func PtracePeekText(pid int, addr uintptr, out []byte) (count int, errno int) {
	return ptracePeek(PTRACE_PEEKTEXT, pid, addr, out)
}

func PtracePeekData(pid int, addr uintptr, out []byte) (count int, errno int) {
	return ptracePeek(PTRACE_PEEKDATA, pid, addr, out)
}

func ptracePoke(pokeReq int, peekReq int, pid int, addr uintptr, data []byte) (count int, errno int) {
	// As for ptracePeek, we need to align our accesses to deal
	// with the possibility of straddling an invalid page.

	// Leading edge.
	n := 0
	if addr%sizeofPtr != 0 {
		var buf [sizeofPtr]byte
		errno = ptrace(peekReq, pid, addr-addr%sizeofPtr, uintptr(unsafe.Pointer(&buf[0])))
		if errno != 0 {
			return 0, errno
		}
		n += copy(buf[addr%sizeofPtr:], data)
		word := *((*uintptr)(unsafe.Pointer(&buf[0])))
		errno = ptrace(pokeReq, pid, addr-addr%sizeofPtr, word)
		if errno != 0 {
			return 0, errno
		}
		data = data[n:]
	}

	// Interior.
	for len(data) > sizeofPtr {
		word := *((*uintptr)(unsafe.Pointer(&data[0])))
		errno = ptrace(pokeReq, pid, addr+uintptr(n), word)
		if errno != 0 {
			return n, errno
		}
		n += sizeofPtr
		data = data[sizeofPtr:]
	}

	// Trailing edge.
	if len(data) > 0 {
		var buf [sizeofPtr]byte
		errno = ptrace(peekReq, pid, addr+uintptr(n), uintptr(unsafe.Pointer(&buf[0])))
		if errno != 0 {
			return n, errno
		}
		copy(buf[0:], data)
		word := *((*uintptr)(unsafe.Pointer(&buf[0])))
		errno = ptrace(pokeReq, pid, addr+uintptr(n), word)
		if errno != 0 {
			return n, errno
		}
		n += len(data)
	}

	return n, 0
}

func PtracePokeText(pid int, addr uintptr, data []byte) (count int, errno int) {
	return ptracePoke(PTRACE_POKETEXT, PTRACE_PEEKTEXT, pid, addr, data)
}

func PtracePokeData(pid int, addr uintptr, data []byte) (count int, errno int) {
	return ptracePoke(PTRACE_POKEDATA, PTRACE_PEEKDATA, pid, addr, data)
}

func PtraceGetRegs(pid int, regsout *PtraceRegs) (errno int) {
	return ptrace(PTRACE_GETREGS, pid, 0, uintptr(unsafe.Pointer(regsout)))
}

func PtraceSetRegs(pid int, regs *PtraceRegs) (errno int) {
	return ptrace(PTRACE_SETREGS, pid, 0, uintptr(unsafe.Pointer(regs)))
}

func PtraceSetOptions(pid int, options int) (errno int) {
	return ptrace(PTRACE_SETOPTIONS, pid, 0, uintptr(options))
}

func PtraceGetEventMsg(pid int) (msg uint, errno int) {
	var data _C_long
	errno = ptrace(PTRACE_GETEVENTMSG, pid, 0, uintptr(unsafe.Pointer(&data)))
	msg = uint(data)
	return
}

func PtraceCont(pid int, signal int) (errno int) {
	return ptrace(PTRACE_CONT, pid, 0, uintptr(signal))
}

func PtraceSingleStep(pid int) (errno int) { return ptrace(PTRACE_SINGLESTEP, pid, 0, 0) }

func PtraceAttach(pid int) (errno int) { return ptrace(PTRACE_ATTACH, pid, 0, 0) }

func PtraceDetach(pid int) (errno int) { return ptrace(PTRACE_DETACH, pid, 0, 0) }

//sys	reboot(magic1 uint, magic2 uint, cmd int, arg string) (errno int)
func Reboot(cmd int) (errno int) {
	return reboot(LINUX_REBOOT_MAGIC1, LINUX_REBOOT_MAGIC2, cmd, "")
}

func clen(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return len(n)
}

func ReadDirent(fd int, buf []byte) (n int, errno int) {
	return Getdents(fd, buf)
}

func ParseDirent(buf []byte, max int, names []string) (consumed int, count int, newnames []string) {
	origlen := len(buf)
	count = 0
	for max != 0 && len(buf) > 0 {
		dirent := (*Dirent)(unsafe.Pointer(&buf[0]))
		buf = buf[dirent.Reclen:]
		if dirent.Ino == 0 { // File absent in directory.
			continue
		}
		bytes := (*[10000]byte)(unsafe.Pointer(&dirent.Name[0]))
		var name = string(bytes[0:clen(bytes[:])])
		if name == "." || name == ".." { // Useless names
			continue
		}
		max--
		count++
		names = append(names, name)
	}
	return origlen - len(buf), count, names
}

// Sendto
// Recvfrom
// Socketpair

/*
 * Direct access
 */
//sys	Access(path string, mode uint32) (errno int)
//sys	Acct(path string) (errno int)
//sys	Adjtimex(buf *Timex) (state int, errno int)
//sys	Chdir(path string) (errno int)
//sys	Chmod(path string, mode uint32) (errno int)
//sys	Chroot(path string) (errno int)
//sys	Close(fd int) (errno int)
//sys	Creat(path string, mode uint32) (fd int, errno int)
//sysnb	Dup(oldfd int) (fd int, errno int)
//sysnb	Dup2(oldfd int, newfd int) (fd int, errno int)
//sysnb	EpollCreate(size int) (fd int, errno int)
//sysnb	EpollCtl(epfd int, op int, fd int, event *EpollEvent) (errno int)
//sys	EpollWait(epfd int, events []EpollEvent, msec int) (n int, errno int)
//sys	Exit(code int) = SYS_EXIT_GROUP
//sys	Faccessat(dirfd int, path string, mode uint32, flags int) (errno int)
//sys	Fallocate(fd int, mode uint32, off int64, len int64) (errno int)
//sys	Fchdir(fd int) (errno int)
//sys	Fchmod(fd int, mode uint32) (errno int)
//sys	Fchmodat(dirfd int, path string, mode uint32, flags int) (errno int)
//sys	Fchownat(dirfd int, path string, uid int, gid int, flags int) (errno int)
//sys	fcntl(fd int, cmd int, arg int) (val int, errno int)
//sys	Fdatasync(fd int) (errno int)
//sys	Fsync(fd int) (errno int)
//sys	Getdents(fd int, buf []byte) (n int, errno int) = SYS_GETDENTS64
//sysnb	Getpgid(pid int) (pgid int, errno int)
//sysnb	Getpgrp() (pid int)
//sysnb	Getpid() (pid int)
//sysnb	Getppid() (ppid int)
//sysnb	Getrlimit(resource int, rlim *Rlimit) (errno int)
//sysnb	Getrusage(who int, rusage *Rusage) (errno int)
//sysnb	Gettid() (tid int)
//sys	InotifyAddWatch(fd int, pathname string, mask uint32) (watchdesc int, errno int)
//sysnb	InotifyInit() (fd int, errno int)
//sysnb	InotifyInit1(flags int) (fd int, errno int)
//sysnb	InotifyRmWatch(fd int, watchdesc uint32) (success int, errno int)
//sysnb	Kill(pid int, sig int) (errno int)
//sys	Klogctl(typ int, buf []byte) (n int, errno int) = SYS_SYSLOG
//sys	Link(oldpath string, newpath string) (errno int)
//sys	Mkdir(path string, mode uint32) (errno int)
//sys	Mkdirat(dirfd int, path string, mode uint32) (errno int)
//sys	Mknod(path string, mode uint32, dev int) (errno int)
//sys	Mknodat(dirfd int, path string, mode uint32, dev int) (errno int)
//sys	Mount(source string, target string, fstype string, flags int, data string) (errno int)
//sys	Nanosleep(time *Timespec, leftover *Timespec) (errno int)
//sys	Pause() (errno int)
//sys	PivotRoot(newroot string, putold string) (errno int) = SYS_PIVOT_ROOT
//sys	Read(fd int, p []byte) (n int, errno int)
//sys	Readlink(path string, buf []byte) (n int, errno int)
//sys	Rename(oldpath string, newpath string) (errno int)
//sys	Renameat(olddirfd int, oldpath string, newdirfd int, newpath string) (errno int)
//sys	Rmdir(path string) (errno int)
//sys	Setdomainname(p []byte) (errno int)
//sys	Sethostname(p []byte) (errno int)
//sysnb	Setpgid(pid int, pgid int) (errno int)
//sysnb	Setrlimit(resource int, rlim *Rlimit) (errno int)
//sysnb	Setsid() (pid int, errno int)
//sysnb	Settimeofday(tv *Timeval) (errno int)
//sysnb	Setuid(uid int) (errno int)
//sys	Symlink(oldpath string, newpath string) (errno int)
//sys	Sync()
//sysnb	Sysinfo(info *Sysinfo_t) (errno int)
//sys	Tee(rfd int, wfd int, len int, flags int) (n int64, errno int)
//sysnb	Tgkill(tgid int, tid int, sig int) (errno int)
//sysnb	Times(tms *Tms) (ticks uintptr, errno int)
//sysnb	Umask(mask int) (oldmask int)
//sysnb	Uname(buf *Utsname) (errno int)
//sys	Unlink(path string) (errno int)
//sys	Unlinkat(dirfd int, path string) (errno int)
//sys	Unmount(target string, flags int) (errno int) = SYS_UMOUNT2
//sys	Unshare(flags int) (errno int)
//sys	Ustat(dev int, ubuf *Ustat_t) (errno int)
//sys	Utime(path string, buf *Utimbuf) (errno int)
//sys	Write(fd int, p []byte) (n int, errno int)
//sys	exitThread(code int) (errno int) = SYS_EXIT
//sys	read(fd int, p *byte, np int) (n int, errno int)
//sys	write(fd int, p *byte, np int) (n int, errno int)

// mmap varies by architecture; see syscall_linux_*.go.
//sys	munmap(addr uintptr, length uintptr) (errno int)

var mapper = &mmapper{
	active: make(map[*byte][]byte),
	mmap:   mmap,
	munmap: munmap,
}

func Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, errno int) {
	return mapper.Mmap(fd, offset, length, prot, flags)
}

func Munmap(b []byte) (errno int) {
	return mapper.Munmap(b)
}

/*
 * Unimplemented
 */
// AddKey
// AfsSyscall
// Alarm
// ArchPrctl
// Brk
// Capget
// Capset
// ClockGetres
// ClockGettime
// ClockNanosleep
// ClockSettime
// Clone
// CreateModule
// DeleteModule
// EpollCtlOld
// EpollPwait
// EpollWaitOld
// Eventfd
// Execve
// Fadvise64
// Fgetxattr
// Flistxattr
// Flock
// Fork
// Fremovexattr
// Fsetxattr
// Futex
// GetKernelSyms
// GetMempolicy
// GetRobustList
// GetThreadArea
// Getitimer
// Getpmsg
// Getpriority
// Getxattr
// IoCancel
// IoDestroy
// IoGetevents
// IoSetup
// IoSubmit
// Ioctl
// IoprioGet
// IoprioSet
// KexecLoad
// Keyctl
// Lgetxattr
// Listxattr
// Llistxattr
// LookupDcookie
// Lremovexattr
// Lsetxattr
// Madvise
// Mbind
// MigratePages
// Mincore
// Mlock
// Mmap
// ModifyLdt
// Mount
// MovePages
// Mprotect
// MqGetsetattr
// MqNotify
// MqOpen
// MqTimedreceive
// MqTimedsend
// MqUnlink
// Mremap
// Msgctl
// Msgget
// Msgrcv
// Msgsnd
// Msync
// Munlock
// Munlockall
// Munmap
// Newfstatat
// Nfsservctl
// Personality
// Poll
// Ppoll
// Prctl
// Pselect6
// Ptrace
// Putpmsg
// QueryModule
// Quotactl
// Readahead
// Readv
// RemapFilePages
// Removexattr
// RequestKey
// RestartSyscall
// RtSigaction
// RtSigpending
// RtSigprocmask
// RtSigqueueinfo
// RtSigreturn
// RtSigsuspend
// RtSigtimedwait
// SchedGetPriorityMax
// SchedGetPriorityMin
// SchedGetaffinity
// SchedGetparam
// SchedGetscheduler
// SchedRrGetInterval
// SchedSetaffinity
// SchedSetparam
// SchedYield
// Security
// Semctl
// Semget
// Semop
// Semtimedop
// Sendfile
// SetMempolicy
// SetRobustList
// SetThreadArea
// SetTidAddress
// Setpriority
// Setxattr
// Shmat
// Shmctl
// Shmdt
// Shmget
// Sigaltstack
// Signalfd
// Swapoff
// Swapon
// Sysfs
// TimerCreate
// TimerDelete
// TimerGetoverrun
// TimerGettime
// TimerSettime
// Timerfd
// Tkill (obsolete)
// Tuxcall
// Umount2
// Uselib
// Utimensat
// Vfork
// Vhangup
// Vmsplice
// Vserver
// Waitid
// Writev
// _Sysctl
