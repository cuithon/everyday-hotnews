// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package testdeps provides access to dependencies needed by test execution.
//
// This package is imported by the generated main package, which passes
// TestDeps into testing.Main. This allows tests to use packages at run time
// without making those packages direct dependencies of package testing.
// Direct dependencies of package testing are harder to write tests for.
package testdeps

import (
	"bufio"
	"context"
	"internal/fuzz"
	"internal/testlog"
	"io"
	"os"
	"os/signal"
	"regexp"
	"runtime/pprof"
	"strings"
	"sync"
	"time"
)

// TestDeps is an implementation of the testing.testDeps interface,
// suitable for passing to testing.MainStart.
type TestDeps struct{}

var matchPat string
var matchRe *regexp.Regexp

func (TestDeps) MatchString(pat, str string) (result bool, err error) {
	if matchRe == nil || matchPat != pat {
		matchPat = pat
		matchRe, err = regexp.Compile(matchPat)
		if err != nil {
			return
		}
	}
	return matchRe.MatchString(str), nil
}

func (TestDeps) StartCPUProfile(w io.Writer) error {
	return pprof.StartCPUProfile(w)
}

func (TestDeps) StopCPUProfile() {
	pprof.StopCPUProfile()
}

func (TestDeps) WriteProfileTo(name string, w io.Writer, debug int) error {
	return pprof.Lookup(name).WriteTo(w, debug)
}

// ImportPath is the import path of the testing binary, set by the generated main function.
var ImportPath string

func (TestDeps) ImportPath() string {
	return ImportPath
}

// testLog implements testlog.Interface, logging actions by package os.
type testLog struct {
	mu  sync.Mutex
	w   *bufio.Writer
	set bool
}

func (l *testLog) Getenv(key string) {
	l.add("getenv", key)
}

func (l *testLog) Open(name string) {
	l.add("open", name)
}

func (l *testLog) Stat(name string) {
	l.add("stat", name)
}

func (l *testLog) Chdir(name string) {
	l.add("chdir", name)
}

// add adds the (op, name) pair to the test log.
func (l *testLog) add(op, name string) {
	if strings.Contains(name, "\n") || name == "" {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	if l.w == nil {
		return
	}
	l.w.WriteString(op)
	l.w.WriteByte(' ')
	l.w.WriteString(name)
	l.w.WriteByte('\n')
}

var log testLog

func (TestDeps) StartTestLog(w io.Writer) {
	log.mu.Lock()
	log.w = bufio.NewWriter(w)
	if !log.set {
		// Tests that define TestMain and then run m.Run multiple times
		// will call StartTestLog/StopTestLog multiple times.
		// Checking log.set avoids calling testlog.SetLogger multiple times
		// (which will panic) and also avoids writing the header multiple times.
		log.set = true
		testlog.SetLogger(&log)
		log.w.WriteString("# test log\n") // known to cmd/go/internal/test/test.go
	}
	log.mu.Unlock()
}

func (TestDeps) StopTestLog() error {
	log.mu.Lock()
	defer log.mu.Unlock()
	err := log.w.Flush()
	log.w = nil
	return err
}

// SetPanicOnExit0 tells the os package whether to panic on os.Exit(0).
func (TestDeps) SetPanicOnExit0(v bool) {
	testlog.SetPanicOnExit0(v)
}

func (TestDeps) CoordinateFuzzing(timeout time.Duration, parallel int, seed [][]byte, corpusDir, cacheDir string) error {
	// Fuzzing may be interrupted with a timeout or if the user presses ^C.
	// In either case, we'll stop worker processes gracefully and save
	// crashers and interesting values.
	ctx := context.Background()
	cancel := func() {}
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	}
	interruptC := make(chan os.Signal, 1)
	signal.Notify(interruptC, os.Interrupt)
	go func() {
		<-interruptC
		cancel()
	}()
	defer close(interruptC)

	err := fuzz.CoordinateFuzzing(ctx, parallel, seed, corpusDir, cacheDir)
	if err == ctx.Err() {
		return nil
	}
	return err
}

func (TestDeps) RunFuzzWorker(fn func([]byte) error) error {
	// Worker processes may or may not receive a signal when the user presses ^C
	// On POSIX operating systems, a signal sent to a process group is delivered
	// to all processes in that group. This is not the case on Windows.
	// If the worker is interrupted, return quickly and without error.
	// If only the coordinator process is interrupted, it tells each worker
	// process to stop by closing its "fuzz_in" pipe.
	ctx, cancel := context.WithCancel(context.Background())
	interruptC := make(chan os.Signal, 1)
	signal.Notify(interruptC, os.Interrupt)
	go func() {
		<-interruptC
		cancel()
	}()
	defer close(interruptC)

	err := fuzz.RunFuzzWorker(ctx, fn)
	if err == ctx.Err() {
		return nil
	}
	return nil
}

func (TestDeps) ReadCorpus(dir string) ([][]byte, error) {
	return fuzz.ReadCorpus(dir)
}
