// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package net provides a portable interface for network I/O, including
TCP/IP, UDP, domain name resolution, and Unix domain sockets.

Although the package provides access to low-level networking
primitives, most clients will need only the basic interface provided
by the Dial, Listen, and Accept functions and the associated
Conn and Listener interfaces. The crypto/tls package uses
the same interfaces and similar Dial and Listen functions.

The Dial function connects to a server:

	conn, err := net.Dial("tcp", "google.com:80")
	if err != nil {
		// handle error
	}
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	status, err := bufio.NewReader(conn).ReadString('\n')
	// ...

The Listen function creates servers:

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
		go handleConnection(conn)
	}
*/
package net

// TODO(rsc):
//	support for raw ethernet sockets

import (
	"errors"
	"time"
)

// Addr represents a network end point address.
type Addr interface {
	Network() string // name of the network
	String() string  // string form of address
}

// Conn is a generic stream-oriented network connection.
type Conn interface {
	// Read reads data from the connection.
	// Read can be made to time out and return a Error with Timeout() == true
	// after a fixed time limit; see SetDeadline and SetReadDeadline.
	Read(b []byte) (n int, err error)

	// Write writes data to the connection.
	// Write can be made to time out and return a Error with Timeout() == true
	// after a fixed time limit; see SetDeadline and SetWriteDeadline.
	Write(b []byte) (n int, err error)

	// Close closes the connection.
	Close() error

	// LocalAddr returns the local network address.
	LocalAddr() Addr

	// RemoteAddr returns the remote network address.
	RemoteAddr() Addr

	// SetDeadline sets the read and write deadlines associated
	// with the connection. It is equivalent to calling both
	// SetReadDeadline and SetWriteDeadline.
	//
	// A deadline is an absolute time after which I/O operations
	// fail with a timeout (see type Error) instead of
	// blocking. The deadline applies to all future I/O, not just
	// the immediately following call to Read or Write.
	//
	// An idle timeout can be implemented by repeatedly extending
	// the deadline after successful Read or Write calls.
	//
	// A zero value for t means I/O operations will not time out.
	SetDeadline(t time.Time) error

	// SetReadDeadline sets the deadline for Read calls.
	// A zero value for t means Read will not time out.
	SetReadDeadline(t time.Time) error

	// SetWriteDeadline sets the deadline for Write calls.
	// Even if write times out, it may return n > 0, indicating that
	// some of the data was successfully written.
	// A zero value for t means Write will not time out.
	SetWriteDeadline(t time.Time) error
}

// An Error represents a network error.
type Error interface {
	error
	Timeout() bool   // Is the error a timeout?
	Temporary() bool // Is the error temporary?
}

// PacketConn is a generic packet-oriented network connection.
type PacketConn interface {
	// ReadFrom reads a packet from the connection,
	// copying the payload into b.  It returns the number of
	// bytes copied into b and the return address that
	// was on the packet.
	// ReadFrom can be made to time out and return
	// an error with Timeout() == true after a fixed time limit;
	// see SetDeadline and SetReadDeadline.
	ReadFrom(b []byte) (n int, addr Addr, err error)

	// WriteTo writes a packet with payload b to addr.
	// WriteTo can be made to time out and return
	// an error with Timeout() == true after a fixed time limit;
	// see SetDeadline and SetWriteDeadline.
	// On packet-oriented connections, write timeouts are rare.
	WriteTo(b []byte, addr Addr) (n int, err error)

	// Close closes the connection.
	Close() error

	// LocalAddr returns the local network address.
	LocalAddr() Addr

	// SetDeadline sets the read and write deadlines associated
	// with the connection.
	SetDeadline(t time.Time) error

	// SetReadDeadline sets the deadline for all Read calls to return.
	// If the deadline is reached, Read will fail with a timeout
	// (see type Error) instead of blocking.
	// A zero value for t means Read will not time out.
	SetReadDeadline(t time.Time) error

	// SetWriteDeadline sets the deadline for all Write calls to return.
	// If the deadline is reached, Write will fail with a timeout
	// (see type Error) instead of blocking.
	// A zero value for t means Write will not time out.
	// Even if write times out, it may return n > 0, indicating that
	// some of the data was successfully written.
	SetWriteDeadline(t time.Time) error
}

// A Listener is a generic network listener for stream-oriented protocols.
type Listener interface {
	// Accept waits for and returns the next connection to the listener.
	Accept() (c Conn, err error)

	// Close closes the listener.
	Close() error

	// Addr returns the listener's network address.
	Addr() Addr
}

var errMissingAddress = errors.New("missing address")

type OpError struct {
	Op   string
	Net  string
	Addr Addr
	Err  error
}

func (e *OpError) Error() string {
	if e == nil {
		return "<nil>"
	}
	s := e.Op
	if e.Net != "" {
		s += " " + e.Net
	}
	if e.Addr != nil {
		s += " " + e.Addr.String()
	}
	s += ": " + e.Err.Error()
	return s
}

type temporary interface {
	Temporary() bool
}

func (e *OpError) Temporary() bool {
	t, ok := e.Err.(temporary)
	return ok && t.Temporary()
}

type timeout interface {
	Timeout() bool
}

func (e *OpError) Timeout() bool {
	t, ok := e.Err.(timeout)
	return ok && t.Timeout()
}

type timeoutError struct{}

func (e *timeoutError) Error() string   { return "i/o timeout" }
func (e *timeoutError) Timeout() bool   { return true }
func (e *timeoutError) Temporary() bool { return true }

var errTimeout error = &timeoutError{}

type AddrError struct {
	Err  string
	Addr string
}

func (e *AddrError) Error() string {
	if e == nil {
		return "<nil>"
	}
	s := e.Err
	if e.Addr != "" {
		s += " " + e.Addr
	}
	return s
}

func (e *AddrError) Temporary() bool {
	return false
}

func (e *AddrError) Timeout() bool {
	return false
}

type UnknownNetworkError string

func (e UnknownNetworkError) Error() string   { return "unknown network " + string(e) }
func (e UnknownNetworkError) Temporary() bool { return false }
func (e UnknownNetworkError) Timeout() bool   { return false }

// DNSConfigError represents an error reading the machine's DNS configuration.
type DNSConfigError struct {
	Err error
}

func (e *DNSConfigError) Error() string {
	return "error reading DNS config: " + e.Err.Error()
}

func (e *DNSConfigError) Timeout() bool   { return false }
func (e *DNSConfigError) Temporary() bool { return false }
