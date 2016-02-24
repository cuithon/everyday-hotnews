// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build cgo

package runtime_test

import (
	"internal/testenv"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestCgoCrashHandler(t *testing.T) {
	testCrashHandler(t, true)
}

func TestCgoSignalDeadlock(t *testing.T) {
	if testing.Short() && runtime.GOOS == "windows" {
		t.Skip("Skipping in short mode") // takes up to 64 seconds
	}
	got := runTestProg(t, "testprogcgo", "CgoSignalDeadlock")
	want := "OK\n"
	if got != want {
		t.Fatalf("expected %q, but got:\n%s", want, got)
	}
}

func TestCgoTraceback(t *testing.T) {
	got := runTestProg(t, "testprogcgo", "CgoTraceback")
	want := "OK\n"
	if got != want {
		t.Fatalf("expected %q, but got:\n%s", want, got)
	}
}

func TestCgoCallbackGC(t *testing.T) {
	if runtime.GOOS == "plan9" || runtime.GOOS == "windows" {
		t.Skipf("no pthreads on %s", runtime.GOOS)
	}
	if testing.Short() {
		switch {
		case runtime.GOOS == "dragonfly":
			t.Skip("see golang.org/issue/11990")
		case runtime.GOOS == "linux" && runtime.GOARCH == "arm":
			t.Skip("too slow for arm builders")
		}
	}
	got := runTestProg(t, "testprogcgo", "CgoCallbackGC")
	want := "OK\n"
	if got != want {
		t.Fatalf("expected %q, but got:\n%s", want, got)
	}
}

func TestCgoExternalThreadPanic(t *testing.T) {
	if runtime.GOOS == "plan9" {
		t.Skipf("no pthreads on %s", runtime.GOOS)
	}
	got := runTestProg(t, "testprogcgo", "CgoExternalThreadPanic")
	want := "panic: BOOM"
	if !strings.Contains(got, want) {
		t.Fatalf("want failure containing %q. output:\n%s\n", want, got)
	}
}

func TestCgoExternalThreadSIGPROF(t *testing.T) {
	// issue 9456.
	switch runtime.GOOS {
	case "plan9", "windows":
		t.Skipf("no pthreads on %s", runtime.GOOS)
	case "darwin":
		if runtime.GOARCH != "arm" && runtime.GOARCH != "arm64" {
			// static constructor needs external linking, but we don't support
			// external linking on OS X 10.6.
			out, err := exec.Command("uname", "-r").Output()
			if err != nil {
				t.Fatalf("uname -r failed: %v", err)
			}
			// OS X 10.6 == Darwin 10.x
			if strings.HasPrefix(string(out), "10.") {
				t.Skipf("no external linking on OS X 10.6")
			}
		}
	}
	if runtime.GOARCH == "ppc64" {
		// TODO(austin) External linking not implemented on
		// ppc64 (issue #8912)
		t.Skipf("no external linking on ppc64")
	}
	got := runTestProg(t, "testprogcgo", "CgoExternalThreadSIGPROF")
	want := "OK\n"
	if got != want {
		t.Fatalf("expected %q, but got:\n%s", want, got)
	}
}

func TestCgoExternalThreadSignal(t *testing.T) {
	// issue 10139
	switch runtime.GOOS {
	case "plan9", "windows":
		t.Skipf("no pthreads on %s", runtime.GOOS)
	}
	got := runTestProg(t, "testprogcgo", "CgoExternalThreadSignal")
	want := "OK\n"
	if got != want {
		t.Fatalf("expected %q, but got:\n%s", want, got)
	}
}

func TestCgoDLLImports(t *testing.T) {
	// test issue 9356
	if runtime.GOOS != "windows" {
		t.Skip("skipping windows specific test")
	}
	got := runTestProg(t, "testprogcgo", "CgoDLLImportsMain")
	want := "OK\n"
	if got != want {
		t.Fatalf("expected %q, but got %v", want, got)
	}
}

func TestCgoExecSignalMask(t *testing.T) {
	// Test issue 13164.
	switch runtime.GOOS {
	case "windows", "plan9":
		t.Skipf("skipping signal mask test on %s", runtime.GOOS)
	}
	got := runTestProg(t, "testprogcgo", "CgoExecSignalMask")
	want := "OK\n"
	if got != want {
		t.Errorf("expected %q, got %v", want, got)
	}
}

func TestEnsureDropM(t *testing.T) {
	// Test for issue 13881.
	switch runtime.GOOS {
	case "windows", "plan9":
		t.Skipf("skipping dropm test on %s", runtime.GOOS)
	}
	got := runTestProg(t, "testprogcgo", "EnsureDropM")
	want := "OK\n"
	if got != want {
		t.Errorf("expected %q, got %v", want, got)
	}
}

// Test for issue 14387.
// Test that the program that doesn't need any cgo pointer checking
// takes about the same amount of time with it as without it.
func TestCgoCheckBytes(t *testing.T) {
	// Make sure we don't count the build time as part of the run time.
	testenv.MustHaveGoBuild(t)
	exe, err := buildTestProg(t, "testprogcgo")
	if err != nil {
		t.Fatal(err)
	}

	cmd := testEnv(exec.Command(exe, "CgoCheckBytes"))
	cmd.Env = append(cmd.Env, "GODEBUG=cgocheck=0")

	start := time.Now()
	cmd.Run()
	d1 := time.Since(start)

	cmd = testEnv(exec.Command(exe, "CgoCheckBytes"))

	start = time.Now()
	cmd.Run()
	d2 := time.Since(start)

	if d1*20 < d2 {
		t.Errorf("cgo check too slow: got %v, expected at most %v", d1, d2*10)
	}
}

func TestCgoPanicDeadlock(t *testing.T) {
	// test issue 14432
	got := runTestProg(t, "testprogcgo", "CgoPanicDeadlock")
	want := "panic: cgo error\n\n"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("output does not start with %q:\n%s", want, got)
	}
}
