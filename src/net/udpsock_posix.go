// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build aix || darwin || dragonfly || freebsd || (js && wasm) || linux || netbsd || openbsd || solaris || windows

package net

import (
	"context"
	"net/netip"
	"syscall"
)

func sockaddrToUDP(sa syscall.Sockaddr) Addr {
	switch sa := sa.(type) {
	case *syscall.SockaddrInet4:
		return &UDPAddr{IP: sa.Addr[0:], Port: sa.Port}
	case *syscall.SockaddrInet6:
		return &UDPAddr{IP: sa.Addr[0:], Port: sa.Port, Zone: zoneCache.name(int(sa.ZoneId))}
	}
	return nil
}

func (a *UDPAddr) family() int {
	if a == nil || len(a.IP) <= IPv4len {
		return syscall.AF_INET
	}
	if a.IP.To4() != nil {
		return syscall.AF_INET
	}
	return syscall.AF_INET6
}

func (a *UDPAddr) sockaddr(family int) (syscall.Sockaddr, error) {
	if a == nil {
		return nil, nil
	}
	return ipToSockaddr(family, a.IP, a.Port, a.Zone)
}

func (a *UDPAddr) toLocal(net string) sockaddr {
	return &UDPAddr{loopbackIP(net), a.Port, a.Zone}
}

func (c *UDPConn) readFrom(b []byte, addr *UDPAddr) (int, *UDPAddr, error) {
	var n int
	var err error
	switch c.fd.family {
	case syscall.AF_INET:
		var from syscall.SockaddrInet4
		n, err = c.fd.readFromInet4(b, &from)
		if err == nil {
			ip := from.Addr // copy from.Addr; ip escapes, so this line allocates 4 bytes
			*addr = UDPAddr{IP: ip[:], Port: from.Port}
		}
	case syscall.AF_INET6:
		var from syscall.SockaddrInet6
		n, err = c.fd.readFromInet6(b, &from)
		if err == nil {
			ip := from.Addr // copy from.Addr; ip escapes, so this line allocates 16 bytes
			*addr = UDPAddr{IP: ip[:], Port: from.Port, Zone: zoneCache.name(int(from.ZoneId))}
		}
	}
	if err != nil {
		// No sockaddr, so don't return UDPAddr.
		addr = nil
	}
	return n, addr, err
}

func (c *UDPConn) readMsg(b, oob []byte) (n, oobn, flags int, addr netip.AddrPort, err error) {
	var sa syscall.Sockaddr
	n, oobn, flags, sa, err = c.fd.readMsg(b, oob, 0)
	switch sa := sa.(type) {
	case *syscall.SockaddrInet4:
		ip := netip.AddrFrom4(sa.Addr)
		addr = netip.AddrPortFrom(ip, uint16(sa.Port))
	case *syscall.SockaddrInet6:
		ip := netip.AddrFrom16(sa.Addr).WithZone(zoneCache.name(int(sa.ZoneId)))
		addr = netip.AddrPortFrom(ip, uint16(sa.Port))
	}
	return
}

func (c *UDPConn) writeTo(b []byte, addr *UDPAddr) (int, error) {
	if c.fd.isConnected {
		return 0, ErrWriteToConnected
	}
	if addr == nil {
		return 0, errMissingAddress
	}

	switch c.fd.family {
	case syscall.AF_INET:
		sa, err := ipToSockaddrInet4(addr.IP, addr.Port)
		if err != nil {
			return 0, err
		}
		return c.fd.writeToInet4(b, sa)
	case syscall.AF_INET6:
		sa, err := ipToSockaddrInet6(addr.IP, addr.Port, addr.Zone)
		if err != nil {
			return 0, err
		}
		return c.fd.writeToInet6(b, sa)
	default:
		return 0, &AddrError{Err: "invalid address family", Addr: addr.IP.String()}
	}
}

func (c *UDPConn) writeMsg(b, oob []byte, addr *UDPAddr) (n, oobn int, err error) {
	if c.fd.isConnected && addr != nil {
		return 0, 0, ErrWriteToConnected
	}
	if !c.fd.isConnected && addr == nil {
		return 0, 0, errMissingAddress
	}
	sa, err := addr.sockaddr(c.fd.family)
	if err != nil {
		return 0, 0, err
	}
	return c.fd.writeMsg(b, oob, sa)
}

func (sd *sysDialer) dialUDP(ctx context.Context, laddr, raddr *UDPAddr) (*UDPConn, error) {
	fd, err := internetSocket(ctx, sd.network, laddr, raddr, syscall.SOCK_DGRAM, 0, "dial", sd.Dialer.Control)
	if err != nil {
		return nil, err
	}
	return newUDPConn(fd), nil
}

func (sl *sysListener) listenUDP(ctx context.Context, laddr *UDPAddr) (*UDPConn, error) {
	fd, err := internetSocket(ctx, sl.network, laddr, nil, syscall.SOCK_DGRAM, 0, "listen", sl.ListenConfig.Control)
	if err != nil {
		return nil, err
	}
	return newUDPConn(fd), nil
}

func (sl *sysListener) listenMulticastUDP(ctx context.Context, ifi *Interface, gaddr *UDPAddr) (*UDPConn, error) {
	fd, err := internetSocket(ctx, sl.network, gaddr, nil, syscall.SOCK_DGRAM, 0, "listen", sl.ListenConfig.Control)
	if err != nil {
		return nil, err
	}
	c := newUDPConn(fd)
	if ip4 := gaddr.IP.To4(); ip4 != nil {
		if err := listenIPv4MulticastUDP(c, ifi, ip4); err != nil {
			c.Close()
			return nil, err
		}
	} else {
		if err := listenIPv6MulticastUDP(c, ifi, gaddr.IP); err != nil {
			c.Close()
			return nil, err
		}
	}
	return c, nil
}

func listenIPv4MulticastUDP(c *UDPConn, ifi *Interface, ip IP) error {
	if ifi != nil {
		if err := setIPv4MulticastInterface(c.fd, ifi); err != nil {
			return err
		}
	}
	if err := setIPv4MulticastLoopback(c.fd, false); err != nil {
		return err
	}
	if err := joinIPv4Group(c.fd, ifi, ip); err != nil {
		return err
	}
	return nil
}

func listenIPv6MulticastUDP(c *UDPConn, ifi *Interface, ip IP) error {
	if ifi != nil {
		if err := setIPv6MulticastInterface(c.fd, ifi); err != nil {
			return err
		}
	}
	if err := setIPv6MulticastLoopback(c.fd, false); err != nil {
		return err
	}
	if err := joinIPv6Group(c.fd, ifi, ip); err != nil {
		return err
	}
	return nil
}
