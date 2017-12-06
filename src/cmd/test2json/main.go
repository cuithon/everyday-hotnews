// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test2json converts go test output to a machine-readable JSON stream.
//
// Usage:
//
//	go tool test2json [-p pkg] [-t] [./pkg.test -test.v]
//
// Test2json runs the given test command and converts its output to JSON;
// with no command specified, test2json expects test output on standard input.
// It writes a corresponding stream of JSON events to standard output.
// There is no unnecessary input or output buffering, so that
// the JSON stream can be read for “live updates” of test status.
//
// The -p flag sets the package reported in each test event.
//
// The -t flag requests that time stamps be added to each test event.
//
// Note that test2json is only intended for converting a single test
// binary's output. To convert the output of a "go test" command,
// use "go test -json" instead of invoking test2json directly.
//
// Output Format
//
// The JSON stream is a newline-separated sequence of TestEvent objects
// corresponding to the Go struct:
//
//	type TestEvent struct {
//		Time    time.Time // encodes as an RFC3339-format string
//		Action  string
//		Package string
//		Test    string
//		Elapsed float64 // seconds
//		Output  string
//	}
//
// The Time field holds the time the event happened.
// It is conventionally omitted for cached test results.
//
// The Action field is one of a fixed set of action descriptions:
//
//	run    - the test has started running
//	pause  - the test has been paused
//	cont   - the test has continued running
//	pass   - the test passed
//	fail   - the test failed
//	output - the test printed output
//
// The Package field, if present, specifies the package being tested.
// When the go command runs parallel tests in -json mode, events from
// different tests are interlaced; the Package field allows readers to
// separate them.
//
// The Test field, if present, specifies the test or example, or benchmark
// function that caused the event. Events for the overall package test
// do not set Test.
//
// The Elapsed field is set for "pass" and "fail" events. It gives the time
// elapsed for the specific test or the overall package test that passed or failed.
//
// The Output field is set for Action == "output" and is a portion of the test's output
// (standard output and standard error merged together). The output is
// unmodified except that invalid UTF-8 output from a test is coerced
// into valid UTF-8 by use of replacement characters. With that one exception,
// the concatenation of the Output fields of all output events is the exact
// output of the test execution.
//
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"

	"cmd/internal/test2json"
)

var (
	flagP = flag.String("p", "", "report `pkg` as the package being tested in each event")
	flagT = flag.Bool("t", false, "include timestamps in events")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: go tool test2json [-p pkg] [-t] [./pkg.test -test.v]\n")
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	var mode test2json.Mode
	if *flagT {
		mode |= test2json.Timestamp
	}
	c := test2json.NewConverter(os.Stdout, *flagP, mode)
	defer c.Close()

	if flag.NArg() == 0 {
		io.Copy(c, os.Stdin)
	} else {
		args := flag.Args()
		cmd := exec.Command(args[0], args[1:]...)
		w := &countWriter{0, c}
		cmd.Stdout = w
		cmd.Stderr = w
		if err := cmd.Run(); err != nil {
			if w.n > 0 {
				// Assume command printed why it failed.
			} else {
				fmt.Fprintf(c, "test2json: %v\n", err)
			}
			c.Close()
			os.Exit(1)
		}
	}
}

type countWriter struct {
	n int64
	w io.Writer
}

func (w *countWriter) Write(b []byte) (int, error) {
	w.n += int64(len(b))
	return w.w.Write(b)
}
