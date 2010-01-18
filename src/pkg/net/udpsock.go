// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// UDP sockets

package net

import (
	"os"
	"syscall"
)

func sockaddrToUDP(sa syscall.Sockaddr) Addr {
	switch sa := sa.(type) {
	case *syscall.SockaddrInet4:
		return &UDPAddr{&sa.Addr, sa.Port}
	case *syscall.SockaddrInet6:
		return &UDPAddr{&sa.Addr, sa.Port}
	}
	return nil
}

// UDPAddr represents the address of a UDP end point.
type UDPAddr struct {
	IP   IP
	Port int
}

// Network returns the address's network name, "udp".
func (a *UDPAddr) Network() string { return "udp" }

func (a *UDPAddr) String() string { return joinHostPort(a.IP.String(), itoa(a.Port)) }

func (a *UDPAddr) family() int {
	if a == nil || len(a.IP) <= 4 {
		return syscall.AF_INET
	}
	if ip := a.IP.To4(); ip != nil {
		return syscall.AF_INET
	}
	return syscall.AF_INET6
}

func (a *UDPAddr) sockaddr(family int) (syscall.Sockaddr, os.Error) {
	return ipToSockaddr(family, a.IP, a.Port)
}

func (a *UDPAddr) toAddr() sockaddr {
	if a == nil { // nil *UDPAddr
		return nil // nil interface
	}
	return a
}

// ResolveUDPAddr parses addr as a UDP address of the form
// host:port and resolves domain names or port names to
// numeric addresses.  A literal IPv6 host address must be
// enclosed in square brackets, as in "[::]:80".
func ResolveUDPAddr(addr string) (*UDPAddr, os.Error) {
	ip, port, err := hostPortToIP("udp", addr)
	if err != nil {
		return nil, err
	}
	return &UDPAddr{ip, port}, nil
}

// UDPConn is the implementation of the Conn and PacketConn
// interfaces for UDP network connections.
type UDPConn struct {
	fd *netFD
}

func newUDPConn(fd *netFD) *UDPConn { return &UDPConn{fd} }

func (c *UDPConn) ok() bool { return c != nil && c.fd != nil }

// Implementation of the Conn interface - see Conn for documentation.

// Read reads data from a single UDP packet on the connection.
// If the slice b is smaller than the arriving packet,
// the excess packet data may be discarded.
//
// Read can be made to time out and return err == os.EAGAIN
// after a fixed time limit; see SetTimeout and SetReadTimeout.
func (c *UDPConn) Read(b []byte) (n int, err os.Error) {
	if !c.ok() {
		return 0, os.EINVAL
	}
	return c.fd.Read(b)
}

// Write writes data to the connection as a single UDP packet.
//
// Write can be made to time out and return err == os.EAGAIN
// after a fixed time limit; see SetTimeout and SetReadTimeout.
func (c *UDPConn) Write(b []byte) (n int, err os.Error) {
	if !c.ok() {
		return 0, os.EINVAL
	}
	return c.fd.Write(b)
}

// Close closes the UDP connection.
func (c *UDPConn) Close() os.Error {
	if !c.ok() {
		return os.EINVAL
	}
	err := c.fd.Close()
	c.fd = nil
	return err
}

// LocalAddr returns the local network address.
func (c *UDPConn) LocalAddr() Addr {
	if !c.ok() {
		return nil
	}
	return c.fd.laddr
}

// RemoteAddr returns the remote network address, a *UDPAddr.
func (c *UDPConn) RemoteAddr() Addr {
	if !c.ok() {
		return nil
	}
	return c.fd.raddr
}

// SetTimeout sets the read and write deadlines associated
// with the connection.
func (c *UDPConn) SetTimeout(nsec int64) os.Error {
	if !c.ok() {
		return os.EINVAL
	}
	return setTimeout(c.fd, nsec)
}

// SetReadTimeout sets the time (in nanoseconds) that
// Read will wait for data before returning os.EAGAIN.
// Setting nsec == 0 (the default) disables the deadline.
func (c *UDPConn) SetReadTimeout(nsec int64) os.Error {
	if !c.ok() {
		return os.EINVAL
	}
	return setReadTimeout(c.fd, nsec)
}

// SetWriteTimeout sets the time (in nanoseconds) that
// Write will wait to send its data before returning os.EAGAIN.
// Setting nsec == 0 (the default) disables the deadline.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
func (c *UDPConn) SetWriteTimeout(nsec int64) os.Error {
	if !c.ok() {
		return os.EINVAL
	}
	return setWriteTimeout(c.fd, nsec)
}

// SetReadBuffer sets the size of the operating system's
// receive buffer associated with the connection.
func (c *UDPConn) SetReadBuffer(bytes int) os.Error {
	if !c.ok() {
		return os.EINVAL
	}
	return setReadBuffer(c.fd, bytes)
}

// SetWriteBuffer sets the size of the operating system's
// transmit buffer associated with the connection.
func (c *UDPConn) SetWriteBuffer(bytes int) os.Error {
	if !c.ok() {
		return os.EINVAL
	}
	return setWriteBuffer(c.fd, bytes)
}

// UDP-specific methods.

// ReadFromUDP reads a UDP packet from c, copying the payload into b.
// It returns the number of bytes copied into b and the return address
// that was on the packet.
//
// ReadFromUDP can be made to time out and return err == os.EAGAIN
// after a fixed time limit; see SetTimeout and SetReadTimeout.
func (c *UDPConn) ReadFromUDP(b []byte) (n int, addr *UDPAddr, err os.Error) {
	if !c.ok() {
		return 0, nil, os.EINVAL
	}
	n, sa, err := c.fd.ReadFrom(b)
	switch sa := sa.(type) {
	case *syscall.SockaddrInet4:
		addr = &UDPAddr{&sa.Addr, sa.Port}
	case *syscall.SockaddrInet6:
		addr = &UDPAddr{&sa.Addr, sa.Port}
	}
	return
}

// ReadFrom reads a UDP packet from c, copying the payload into b.
// It returns the number of bytes copied into b and the return address
// that was on the packet.
//
// ReadFrom can be made to time out and return err == os.EAGAIN
// after a fixed time limit; see SetTimeout and SetReadTimeout.
func (c *UDPConn) ReadFrom(b []byte) (n int, addr Addr, err os.Error) {
	if !c.ok() {
		return 0, nil, os.EINVAL
	}
	n, uaddr, err := c.ReadFromUDP(b)
	return n, uaddr.toAddr(), err
}

// WriteToUDP writes a UDP packet to addr via c, copying the payload from b.
//
// WriteToUDP can be made to time out and return err == os.EAGAIN
// after a fixed time limit; see SetTimeout and SetWriteTimeout.
// On packet-oriented connections such as UDP, write timeouts are rare.
func (c *UDPConn) WriteToUDP(b []byte, addr *UDPAddr) (n int, err os.Error) {
	if !c.ok() {
		return 0, os.EINVAL
	}
	sa, err := addr.sockaddr(c.fd.family)
	if err != nil {
		return 0, err
	}
	return c.fd.WriteTo(b, sa)
}

// WriteTo writes a UDP packet with payload b to addr via c.
//
// WriteTo can be made to time out and return err == os.EAGAIN
// after a fixed time limit; see SetTimeout and SetWriteTimeout.
// On packet-oriented connections such as UDP, write timeouts are rare.
func (c *UDPConn) WriteTo(b []byte, addr Addr) (n int, err os.Error) {
	if !c.ok() {
		return 0, os.EINVAL
	}
	a, ok := addr.(*UDPAddr)
	if !ok {
		return 0, &OpError{"writeto", "udp", addr, os.EINVAL}
	}
	return c.WriteToUDP(b, a)
}

// DialUDP connects to the remote address raddr on the network net,
// which must be "udp", "udp4", or "udp6".  If laddr is not nil, it is used
// as the local address for the connection.
func DialUDP(net string, laddr, raddr *UDPAddr) (c *UDPConn, err os.Error) {
	switch net {
	case "udp", "udp4", "udp6":
	default:
		return nil, UnknownNetworkError(net)
	}
	if raddr == nil {
		return nil, &OpError{"dial", "udp", nil, errMissingAddress}
	}
	fd, e := internetSocket(net, laddr.toAddr(), raddr.toAddr(), syscall.SOCK_DGRAM, "dial", sockaddrToUDP)
	if e != nil {
		return nil, e
	}
	return newUDPConn(fd), nil
}

// ListenUDP listens for incoming UDP packets addressed to the
// local address laddr.  The returned connection c's ReadFrom
// and WriteTo methods can be used to receive and send UDP
// packets with per-packet addressing.
func ListenUDP(net string, laddr *UDPAddr) (c *UDPConn, err os.Error) {
	switch net {
	case "udp", "udp4", "udp6":
	default:
		return nil, UnknownNetworkError(net)
	}
	if laddr == nil {
		return nil, &OpError{"listen", "udp", nil, errMissingAddress}
	}
	fd, e := internetSocket(net, laddr.toAddr(), nil, syscall.SOCK_DGRAM, "dial", sockaddrToUDP)
	if e != nil {
		return nil, e
	}
	return newUDPConn(fd), nil
}
