// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package poll

import "syscall"

const (
	// spliceNonblock makes calls to splice(2) non-blocking.
	spliceNonblock = 0x2

	// maxSpliceSize is the maximum amount of data Splice asks
	// the kernel to move in a single call to splice(2).
	maxSpliceSize = 4 << 20
)

// Splice transfers at most remain bytes of data from src to dst, using the
// splice system call to minimize copies of data from and to userspace.
//
// Splice creates a temporary pipe, to serve as a buffer for the data transfer.
// src and dst must both be stream-oriented sockets.
//
// If err != nil, sc is the system call which caused the error.
func Splice(dst, src *FD, remain int64) (written int64, handled bool, sc string, err error) {
	prfd, pwfd, sc, err := newTempPipe()
	if err != nil {
		return 0, false, sc, err
	}
	defer destroyTempPipe(prfd, pwfd)
	// From here on, the operation should be considered handled,
	// even if Splice doesn't transfer any data.
	if err := src.readLock(); err != nil {
		return 0, true, "splice", err
	}
	defer src.readUnlock()
	if err := dst.writeLock(); err != nil {
		return 0, true, "splice", err
	}
	defer dst.writeUnlock()
	if err := src.pd.prepareRead(src.isFile); err != nil {
		return 0, true, "splice", err
	}
	if err := dst.pd.prepareWrite(dst.isFile); err != nil {
		return 0, true, "splice", err
	}
	var inPipe, n int
	for err == nil && remain > 0 {
		max := maxSpliceSize
		if int64(max) > remain {
			max = int(remain)
		}
		inPipe, err = spliceDrain(pwfd, src, max)
		// spliceDrain should never return EAGAIN, so if err != nil,
		// Splice cannot continue. If inPipe == 0 && err == nil,
		// src is at EOF, and the transfer is complete.
		if err != nil || (inPipe == 0 && err == nil) {
			break
		}
		n, err = splicePump(dst, prfd, inPipe)
		if n > 0 {
			written += int64(n)
			remain -= int64(n)
		}
	}
	if err != nil {
		return written, true, "splice", err
	}
	return written, true, "", nil
}

// spliceDrain moves data from a socket to a pipe.
//
// Invariant: when entering spliceDrain, the pipe is empty. It is either in its
// initial state, or splicePump has emptied it previously.
//
// Given this, spliceDrain can reasonably assume that the pipe is ready for
// writing, so if splice returns EAGAIN, it must be because the socket is not
// ready for reading.
//
// If spliceDrain returns (0, nil), src is at EOF.
func spliceDrain(pipefd int, sock *FD, max int) (int, error) {
	for {
		n, err := splice(pipefd, sock.Sysfd, max, spliceNonblock)
		if err != syscall.EAGAIN {
			return n, err
		}
		if err := sock.pd.waitRead(sock.isFile); err != nil {
			return n, err
		}
	}
}

// splicePump moves all the buffered data from a pipe to a socket.
//
// Invariant: when entering splicePump, there are exactly inPipe
// bytes of data in the pipe, from a previous call to spliceDrain.
//
// By analogy to the condition from spliceDrain, splicePump
// only needs to poll the socket for readiness, if splice returns
// EAGAIN.
//
// If splicePump cannot move all the data in a single call to
// splice(2), it loops over the buffered data until it has written
// all of it to the socket. This behavior is similar to the Write
// step of an io.Copy in userspace.
func splicePump(sock *FD, pipefd int, inPipe int) (int, error) {
	written := 0
	for inPipe > 0 {
		n, err := splice(sock.Sysfd, pipefd, inPipe, spliceNonblock)
		// Here, the condition n == 0 && err == nil should never be
		// observed, since Splice controls the write side of the pipe.
		if n > 0 {
			inPipe -= n
			written += n
			continue
		}
		if err != syscall.EAGAIN {
			return written, err
		}
		if err := sock.pd.waitWrite(sock.isFile); err != nil {
			return written, err
		}
	}
	return written, nil
}

// splice wraps the splice system call. Since the current implementation
// only uses splice on sockets and pipes, the offset arguments are unused.
// splice returns int instead of int64, because callers never ask it to
// move more data in a single call than can fit in an int32.
func splice(out int, in int, max int, flags int) (int, error) {
	n, err := syscall.Splice(in, nil, out, nil, max, flags)
	return int(n), err
}

// newTempPipe sets up a temporary pipe for a splice operation.
func newTempPipe() (prfd, pwfd int, sc string, err error) {
	var fds [2]int
	const flags = syscall.O_CLOEXEC | syscall.O_NONBLOCK
	if err := syscall.Pipe2(fds[:], flags); err != nil {
		// pipe2 was added in 2.6.27 and our minimum requirement
		// is 2.6.23, so it might not be implemented.
		if err == syscall.ENOSYS {
			return newTempPipeFallback(fds[:])
		}
		return -1, -1, "pipe2", err
	}
	return fds[0], fds[1], "", nil
}

// newTempPipeFallback is a fallback for newTempPipe, for systems
// which do not support pipe2.
func newTempPipeFallback(fds []int) (prfd, pwfd int, sc string, err error) {
	syscall.ForkLock.RLock()
	defer syscall.ForkLock.RUnlock()
	if err := syscall.Pipe(fds); err != nil {
		return -1, -1, "pipe", err
	}
	prfd, pwfd = fds[0], fds[1]
	syscall.CloseOnExec(prfd)
	syscall.CloseOnExec(pwfd)
	if err := syscall.SetNonblock(prfd, true); err != nil {
		CloseFunc(prfd)
		CloseFunc(pwfd)
		return -1, -1, "setnonblock", err
	}
	if err := syscall.SetNonblock(pwfd, true); err != nil {
		CloseFunc(prfd)
		CloseFunc(pwfd)
		return -1, -1, "setnonblock", err
	}
	return prfd, pwfd, "", nil
}

// destroyTempPipe destroys a temporary pipe.
func destroyTempPipe(prfd, pwfd int) error {
	err := CloseFunc(prfd)
	err1 := CloseFunc(pwfd)
	if err == nil {
		return err1
	}
	return err
}
