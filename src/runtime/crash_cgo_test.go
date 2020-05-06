// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build cgo

package runtime_test

import (
	"bytes"
	"fmt"
	"internal/testenv"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestCgoCrashHandler(t *testing.T) {
	t.Parallel()
	testCrashHandler(t, true)
}

func TestCgoSignalDeadlock(t *testing.T) {
	// Don't call t.Parallel, since too much work going on at the
	// same time can cause the testprogcgo code to overrun its
	// timeouts (issue #18598).

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
	t.Parallel()
	got := runTestProg(t, "testprogcgo", "CgoTraceback")
	want := "OK\n"
	if got != want {
		t.Fatalf("expected %q, but got:\n%s", want, got)
	}
}

func TestCgoCallbackGC(t *testing.T) {
	t.Parallel()
	switch runtime.GOOS {
	case "plan9", "windows":
		t.Skipf("no pthreads on %s", runtime.GOOS)
	}
	if testing.Short() {
		switch {
		case runtime.GOOS == "dragonfly":
			t.Skip("see golang.org/issue/11990")
		case runtime.GOOS == "linux" && runtime.GOARCH == "arm":
			t.Skip("too slow for arm builders")
		case runtime.GOOS == "linux" && (runtime.GOARCH == "mips64" || runtime.GOARCH == "mips64le"):
			t.Skip("too slow for mips64x builders")
		}
	}
	got := runTestProg(t, "testprogcgo", "CgoCallbackGC")
	want := "OK\n"
	if got != want {
		t.Fatalf("expected %q, but got:\n%s", want, got)
	}
}

func TestCgoExternalThreadPanic(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	// issue 9456.
	switch runtime.GOOS {
	case "plan9", "windows":
		t.Skipf("no pthreads on %s", runtime.GOOS)
	}
	if runtime.GOARCH == "ppc64" && runtime.GOOS == "linux" {
		// TODO(austin) External linking not implemented on
		// linux/ppc64 (issue #8912)
		t.Skipf("no external linking on ppc64")
	}

	exe, err := buildTestProg(t, "testprogcgo", "-tags=threadprof")
	if err != nil {
		t.Fatal(err)
	}

	got, err := testenv.CleanCmdEnv(exec.Command(exe, "CgoExternalThreadSIGPROF")).CombinedOutput()
	if err != nil {
		t.Fatalf("exit status: %v\n%s", err, got)
	}

	if want := "OK\n"; string(got) != want {
		t.Fatalf("expected %q, but got:\n%s", want, got)
	}
}

func TestCgoExternalThreadSignal(t *testing.T) {
	t.Parallel()
	// issue 10139
	switch runtime.GOOS {
	case "plan9", "windows":
		t.Skipf("no pthreads on %s", runtime.GOOS)
	}

	exe, err := buildTestProg(t, "testprogcgo", "-tags=threadprof")
	if err != nil {
		t.Fatal(err)
	}

	got, err := testenv.CleanCmdEnv(exec.Command(exe, "CgoExternalThreadSIGPROF")).CombinedOutput()
	if err != nil {
		t.Fatalf("exit status: %v\n%s", err, got)
	}

	want := []byte("OK\n")
	if !bytes.Equal(got, want) {
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
	// Make sure we don't count the build time as part of the run time.
	testenv.MustHaveGoBuild(t)
	exe, err := buildTestProg(t, "testprogcgo")
	if err != nil {
		t.Fatal(err)
	}

	// Try it 10 times to avoid flakiness.
	const tries = 10
	var tot1, tot2 time.Duration
	for i := 0; i < tries; i++ {
		cmd := testenv.CleanCmdEnv(exec.Command(exe, "CgoCheckBytes"))
		cmd.Env = append(cmd.Env, "GODEBUG=cgocheck=0", fmt.Sprintf("GO_CGOCHECKBYTES_TRY=%d", i))

		start := time.Now()
		cmd.Run()
		d1 := time.Since(start)

		cmd = testenv.CleanCmdEnv(exec.Command(exe, "CgoCheckBytes"))
		cmd.Env = append(cmd.Env, fmt.Sprintf("GO_CGOCHECKBYTES_TRY=%d", i))

		start = time.Now()
		cmd.Run()
		d2 := time.Since(start)

		if d1*20 > d2 {
			// The slow version (d2) was less than 20 times
			// slower than the fast version (d1), so OK.
			return
		}

		tot1 += d1
		tot2 += d2
	}

	t.Errorf("cgo check too slow: got %v, expected at most %v", tot2/tries, (tot1/tries)*20)
}

func TestCgoPanicDeadlock(t *testing.T) {
	t.Parallel()
	// test issue 14432
	got := runTestProg(t, "testprogcgo", "CgoPanicDeadlock")
	want := "panic: cgo error\n\n"
	if !strings.HasPrefix(got, want) {
		t.Fatalf("output does not start with %q:\n%s", want, got)
	}
}

func TestCgoCCodeSIGPROF(t *testing.T) {
	t.Parallel()
	got := runTestProg(t, "testprogcgo", "CgoCCodeSIGPROF")
	want := "OK\n"
	if got != want {
		t.Errorf("expected %q got %v", want, got)
	}
}

func TestCgoCrashTraceback(t *testing.T) {
	t.Parallel()
	switch platform := runtime.GOOS + "/" + runtime.GOARCH; platform {
	case "darwin/amd64":
	case "linux/amd64":
	case "linux/ppc64le":
	default:
		t.Skipf("not yet supported on %s", platform)
	}
	got := runTestProg(t, "testprogcgo", "CrashTraceback")
	for i := 1; i <= 3; i++ {
		if !strings.Contains(got, fmt.Sprintf("cgo symbolizer:%d", i)) {
			t.Errorf("missing cgo symbolizer:%d", i)
		}
	}
}

func TestCgoTracebackContext(t *testing.T) {
	t.Parallel()
	got := runTestProg(t, "testprogcgo", "TracebackContext")
	want := "OK\n"
	if got != want {
		t.Errorf("expected %q got %v", want, got)
	}
}

func testCgoPprof(t *testing.T, buildArg, runArg, top, bottom string) {
	t.Parallel()
	if runtime.GOOS != "linux" || (runtime.GOARCH != "amd64" && runtime.GOARCH != "ppc64le") {
		t.Skipf("not yet supported on %s/%s", runtime.GOOS, runtime.GOARCH)
	}
	testenv.MustHaveGoRun(t)

	exe, err := buildTestProg(t, "testprogcgo", buildArg)
	if err != nil {
		t.Fatal(err)
	}

	// pprofCgoTraceback is called whenever CGO code is executing and a signal
	// is received. Disable signal preemption to increase the likelihood at
	// least one SIGPROF signal fired to capture a sample. See issue #37201.
	cmd := testenv.CleanCmdEnv(exec.Command(exe, runArg))
	cmd.Env = append(cmd.Env, "GODEBUG=asyncpreemptoff=1")

	got, err := cmd.CombinedOutput()
	if err != nil {
		if testenv.Builder() == "linux-amd64-alpine" {
			// See Issue 18243 and Issue 19938.
			t.Skipf("Skipping failing test on Alpine (golang.org/issue/18243). Ignoring error: %v", err)
		}
		t.Fatalf("%s\n\n%v", got, err)
	}
	fn := strings.TrimSpace(string(got))
	defer os.Remove(fn)

	for try := 0; try < 2; try++ {
		cmd := testenv.CleanCmdEnv(exec.Command(testenv.GoToolPath(t), "tool", "pprof", "-traces"))
		// Check that pprof works both with and without explicit executable on command line.
		if try == 0 {
			cmd.Args = append(cmd.Args, exe, fn)
		} else {
			cmd.Args = append(cmd.Args, fn)
		}

		found := false
		for i, e := range cmd.Env {
			if strings.HasPrefix(e, "PPROF_TMPDIR=") {
				cmd.Env[i] = "PPROF_TMPDIR=" + os.TempDir()
				found = true
				break
			}
		}
		if !found {
			cmd.Env = append(cmd.Env, "PPROF_TMPDIR="+os.TempDir())
		}

		out, err := cmd.CombinedOutput()
		t.Logf("%s:\n%s", cmd.Args, out)
		if err != nil {
			t.Error(err)
			continue
		}

		trace := findTrace(string(out), top)
		if len(trace) == 0 {
			t.Errorf("%s traceback missing.", top)
			continue
		}
		if trace[len(trace)-1] != bottom {
			t.Errorf("invalid traceback origin: got=%v; want=[%s ... %s]", trace, top, bottom)
		}
	}
}

func TestCgoPprof(t *testing.T) {
	testCgoPprof(t, "", "CgoPprof", "cpuHog", "runtime.main")
}

func TestCgoPprofPIE(t *testing.T) {
	testCgoPprof(t, "-buildmode=pie", "CgoPprof", "cpuHog", "runtime.main")
}

func TestCgoPprofThread(t *testing.T) {
	testCgoPprof(t, "", "CgoPprofThread", "cpuHogThread", "cpuHogThread2")
}

func TestCgoPprofThreadNoTraceback(t *testing.T) {
	testCgoPprof(t, "", "CgoPprofThreadNoTraceback", "cpuHogThread", "runtime._ExternalCode")
}

func TestRaceProf(t *testing.T) {
	if (runtime.GOOS != "linux" && runtime.GOOS != "freebsd") || runtime.GOARCH != "amd64" {
		t.Skipf("not yet supported on %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	testenv.MustHaveGoRun(t)

	// This test requires building various packages with -race, so
	// it's somewhat slow.
	if testing.Short() {
		t.Skip("skipping test in -short mode")
	}

	exe, err := buildTestProg(t, "testprogcgo", "-race")
	if err != nil {
		t.Fatal(err)
	}

	got, err := testenv.CleanCmdEnv(exec.Command(exe, "CgoRaceprof")).CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	want := "OK\n"
	if string(got) != want {
		t.Errorf("expected %q got %s", want, got)
	}
}

func TestRaceSignal(t *testing.T) {
	t.Parallel()
	if (runtime.GOOS != "linux" && runtime.GOOS != "freebsd") || runtime.GOARCH != "amd64" {
		t.Skipf("not yet supported on %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	testenv.MustHaveGoRun(t)

	// This test requires building various packages with -race, so
	// it's somewhat slow.
	if testing.Short() {
		t.Skip("skipping test in -short mode")
	}

	exe, err := buildTestProg(t, "testprogcgo", "-race")
	if err != nil {
		t.Fatal(err)
	}

	got, err := testenv.CleanCmdEnv(exec.Command(exe, "CgoRaceSignal")).CombinedOutput()
	if err != nil {
		t.Logf("%s\n", got)
		t.Fatal(err)
	}
	want := "OK\n"
	if string(got) != want {
		t.Errorf("expected %q got %s", want, got)
	}
}

func TestCgoNumGoroutine(t *testing.T) {
	switch runtime.GOOS {
	case "windows", "plan9":
		t.Skipf("skipping numgoroutine test on %s", runtime.GOOS)
	}
	t.Parallel()
	got := runTestProg(t, "testprogcgo", "NumGoroutine")
	want := "OK\n"
	if got != want {
		t.Errorf("expected %q got %v", want, got)
	}
}

func TestCatchPanic(t *testing.T) {
	t.Parallel()
	switch runtime.GOOS {
	case "plan9", "windows":
		t.Skipf("no signals on %s", runtime.GOOS)
	case "darwin":
		if runtime.GOARCH == "amd64" {
			t.Skipf("crash() on darwin/amd64 doesn't raise SIGABRT")
		}
	}

	testenv.MustHaveGoRun(t)

	exe, err := buildTestProg(t, "testprogcgo")
	if err != nil {
		t.Fatal(err)
	}

	for _, early := range []bool{true, false} {
		cmd := testenv.CleanCmdEnv(exec.Command(exe, "CgoCatchPanic"))
		// Make sure a panic results in a crash.
		cmd.Env = append(cmd.Env, "GOTRACEBACK=crash")
		if early {
			// Tell testprogcgo to install an early signal handler for SIGABRT
			cmd.Env = append(cmd.Env, "CGOCATCHPANIC_EARLY_HANDLER=1")
		}
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Errorf("testprogcgo CgoCatchPanic failed: %v\n%s", err, out)
		}
	}
}

func TestCgoLockOSThreadExit(t *testing.T) {
	switch runtime.GOOS {
	case "plan9", "windows":
		t.Skipf("no pthreads on %s", runtime.GOOS)
	}
	t.Parallel()
	testLockOSThreadExit(t, "testprogcgo")
}

func TestWindowsStackMemoryCgo(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("skipping windows specific test")
	}
	testenv.SkipFlaky(t, 22575)
	o := runTestProg(t, "testprogcgo", "StackMemory")
	stackUsage, err := strconv.Atoi(o)
	if err != nil {
		t.Fatalf("Failed to read stack usage: %v", err)
	}
	if expected, got := 100<<10, stackUsage; got > expected {
		t.Fatalf("expected < %d bytes of memory per thread, got %d", expected, got)
	}
}

func TestSigStackSwapping(t *testing.T) {
	switch runtime.GOOS {
	case "plan9", "windows":
		t.Skipf("no sigaltstack on %s", runtime.GOOS)
	}
	t.Parallel()
	got := runTestProg(t, "testprogcgo", "SigStack")
	want := "OK\n"
	if got != want {
		t.Errorf("expected %q got %v", want, got)
	}
}

func TestCgoTracebackSigpanic(t *testing.T) {
	// Test unwinding over a sigpanic in C code without a C
	// symbolizer. See issue #23576.
	if runtime.GOOS == "windows" {
		// On Windows if we get an exception in C code, we let
		// the Windows exception handler unwind it, rather
		// than injecting a sigpanic.
		t.Skip("no sigpanic in C on windows")
	}
	t.Parallel()
	got := runTestProg(t, "testprogcgo", "TracebackSigpanic")
	want := "runtime.sigpanic"
	if !strings.Contains(got, want) {
		t.Fatalf("want failure containing %q. output:\n%s\n", want, got)
	}
	nowant := "unexpected return pc"
	if strings.Contains(got, nowant) {
		t.Fatalf("failure incorrectly contains %q. output:\n%s\n", nowant, got)
	}
}

// Test that C code called via cgo can use large Windows thread stacks
// and call back in to Go without crashing. See issue #20975.
//
// See also TestBigStackCallbackSyscall.
func TestBigStackCallbackCgo(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("skipping windows specific test")
	}
	t.Parallel()
	got := runTestProg(t, "testprogcgo", "BigStack")
	want := "OK\n"
	if got != want {
		t.Errorf("expected %q got %v", want, got)
	}
}

func nextTrace(lines []string) ([]string, []string) {
	var trace []string
	for n, line := range lines {
		if strings.HasPrefix(line, "---") {
			return trace, lines[n+1:]
		}
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) == 0 {
			continue
		}
		// Last field contains the function name.
		trace = append(trace, fields[len(fields)-1])
	}
	return nil, nil
}

func findTrace(text, top string) []string {
	lines := strings.Split(text, "\n")
	_, lines = nextTrace(lines) // Skip the header.
	for len(lines) > 0 {
		var t []string
		t, lines = nextTrace(lines)
		if len(t) == 0 {
			continue
		}
		if t[0] == top {
			return t
		}
	}
	return nil
}

func TestSegv(t *testing.T) {
	switch runtime.GOOS {
	case "plan9", "windows":
		t.Skipf("no signals on %s", runtime.GOOS)
	}

	for _, test := range []string{"Segv", "SegvInCgo"} {
		t.Run(test, func(t *testing.T) {
			t.Parallel()
			got := runTestProg(t, "testprogcgo", test)
			t.Log(got)
			if !strings.Contains(got, "SIGSEGV") {
				t.Errorf("expected crash from signal")
			}
		})
	}
}
