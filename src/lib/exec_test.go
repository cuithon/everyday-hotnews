// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exec

import (
	"exec";
	"io";
	"testing";
)

func TestRunCat(t *testing.T) {
	cmd, err := exec.Run("/bin/cat", []string{"cat"}, nil,
		exec.Pipe, exec.Pipe, exec.DevNull);
	if err != nil {
		t.Fatalf("opencmd /bin/cat: %v", err);
	}
	io.WriteString(cmd.Stdin, "hello, world\n");
	cmd.Stdin.Close();
	var buf [64]byte;
	n, err1 := io.FullRead(cmd.Stdout, &buf);
	if err1 != nil && err1 != io.ErrEOF {
		t.Fatalf("reading from /bin/cat: %v", err1);
	}
	if string(buf[0:n]) != "hello, world\n" {
		t.Fatalf("reading from /bin/cat: got %q", buf[0:n]);
	}
	if err1 = cmd.Close(); err1 != nil {
		t.Fatalf("closing /bin/cat: %v", err1);
	}
}

func TestRunEcho(t *testing.T) {
	cmd, err := Run("/bin/echo", []string{"echo", "hello", "world"}, nil,
		exec.DevNull, exec.Pipe, exec.DevNull);
	if err != nil {
		t.Fatalf("opencmd /bin/echo: %v", err);
	}
	var buf [64]byte;
	n, err1 := io.FullRead(cmd.Stdout, &buf);
	if err1 != nil && err1 != io.ErrEOF {
		t.Fatalf("reading from /bin/echo: %v", err1);
	}
	if string(buf[0:n]) != "hello world\n" {
		t.Fatalf("reading from /bin/echo: got %q", buf[0:n]);
	}
	if err1 = cmd.Close(); err1 != nil {
		t.Fatalf("closing /bin/echo: %v", err1);
	}
}
