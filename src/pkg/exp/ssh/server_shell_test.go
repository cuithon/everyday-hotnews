// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"io"
	"testing"
)

type MockChannel struct {
	toSend       []byte
	bytesPerRead int
	received     []byte
}

func (c *MockChannel) Accept() error {
	return nil
}

func (c *MockChannel) Reject(RejectionReason, string) error {
	return nil
}

func (c *MockChannel) Read(data []byte) (n int, err error) {
	n = len(data)
	if n == 0 {
		return
	}
	if n > len(c.toSend) {
		n = len(c.toSend)
	}
	if n == 0 {
		return 0, io.EOF
	}
	if c.bytesPerRead > 0 && n > c.bytesPerRead {
		n = c.bytesPerRead
	}
	copy(data, c.toSend[:n])
	c.toSend = c.toSend[n:]
	return
}

func (c *MockChannel) Write(data []byte) (n int, err error) {
	c.received = append(c.received, data...)
	return len(data), nil
}

func (c *MockChannel) Close() error {
	return nil
}

func (c *MockChannel) AckRequest(ok bool) error {
	return nil
}

func (c *MockChannel) ChannelType() string {
	return ""
}

func (c *MockChannel) ExtraData() []byte {
	return nil
}

func TestClose(t *testing.T) {
	c := &MockChannel{}
	ss := NewServerShell(c, "> ")
	line, err := ss.ReadLine()
	if line != "" {
		t.Errorf("Expected empty line but got: %s", line)
	}
	if err != io.EOF {
		t.Errorf("Error should have been EOF but got: %s", err)
	}
}

var keyPressTests = []struct {
	in   string
	line string
	err  error
}{
	{
		"",
		"",
		io.EOF,
	},
	{
		"\r",
		"",
		nil,
	},
	{
		"foo\r",
		"foo",
		nil,
	},
	{
		"a\x1b[Cb\r", // right
		"ab",
		nil,
	},
	{
		"a\x1b[Db\r", // left
		"ba",
		nil,
	},
	{
		"a\177b\r", // backspace
		"b",
		nil,
	},
}

func TestKeyPresses(t *testing.T) {
	for i, test := range keyPressTests {
		for j := 0; j < len(test.in); j++ {
			c := &MockChannel{
				toSend:       []byte(test.in),
				bytesPerRead: j,
			}
			ss := NewServerShell(c, "> ")
			line, err := ss.ReadLine()
			if line != test.line {
				t.Errorf("Line resulting from test %d (%d bytes per read) was '%s', expected '%s'", i, j, line, test.line)
				break
			}
			if err != test.err {
				t.Errorf("Error resulting from test %d (%d bytes per read) was '%v', expected '%v'", i, j, err, test.err)
				break
			}
		}
	}
}
