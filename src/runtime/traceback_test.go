// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime_test

import (
	"bytes"
	"fmt"
	"internal/abi"
	"internal/testenv"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"testing"
	_ "unsafe"
)

// Test traceback printing of inlined frames.
func TestTracebackInlined(t *testing.T) {
	testenv.SkipIfOptimizationOff(t) // This test requires inlining
	check := func(t *testing.T, r *ttiResult, funcs ...string) {
		t.Helper()

		// Check the printed traceback.
		frames := parseTraceback1(t, r.printed).frames
		t.Log(r.printed)
		// Find ttiLeaf
		for len(frames) > 0 && frames[0].funcName != "runtime_test.ttiLeaf" {
			frames = frames[1:]
		}
		if len(frames) == 0 {
			t.Errorf("missing runtime_test.ttiLeaf")
			return
		}
		frames = frames[1:]
		// Check the function sequence.
		for i, want := range funcs {
			got := "<end>"
			if i < len(frames) {
				got = frames[i].funcName
				if strings.HasSuffix(want, ")") {
					got += "(" + frames[i].args + ")"
				}
			}
			if got != want {
				t.Errorf("got %s, want %s", got, want)
				return
			}
		}
	}

	t.Run("simple", func(t *testing.T) {
		// Check a simple case of inlining
		r := ttiSimple1()
		check(t, r, "runtime_test.ttiSimple3(...)", "runtime_test.ttiSimple2(...)", "runtime_test.ttiSimple1()")
	})

	t.Run("sigpanic", func(t *testing.T) {
		// Check that sigpanic from an inlined function prints correctly
		r := ttiSigpanic1()
		check(t, r, "runtime_test.ttiSigpanic1.func1()", "panic", "runtime_test.ttiSigpanic3(...)", "runtime_test.ttiSigpanic2(...)", "runtime_test.ttiSigpanic1()")
	})

	t.Run("wrapper", func(t *testing.T) {
		// Check that a method inlined into a wrapper prints correctly
		r := ttiWrapper1()
		check(t, r, "runtime_test.ttiWrapper.m1(...)", "runtime_test.ttiWrapper1()")
	})

	t.Run("excluded", func(t *testing.T) {
		// Check that when F -> G is inlined and F is excluded from stack
		// traces, G still appears.
		r := ttiExcluded1()
		check(t, r, "runtime_test.ttiExcluded3(...)", "runtime_test.ttiExcluded1()")
	})
}

type ttiResult struct {
	printed string
}

//go:noinline
func ttiLeaf() *ttiResult {
	// Get a printed stack trace.
	printed := string(debug.Stack())
	return &ttiResult{printed}
}

//go:noinline
func ttiSimple1() *ttiResult {
	return ttiSimple2()
}
func ttiSimple2() *ttiResult {
	return ttiSimple3()
}
func ttiSimple3() *ttiResult {
	return ttiLeaf()
}

//go:noinline
func ttiSigpanic1() (res *ttiResult) {
	defer func() {
		res = ttiLeaf()
		recover()
	}()
	ttiSigpanic2()
	panic("did not panic")
}
func ttiSigpanic2() {
	ttiSigpanic3()
}
func ttiSigpanic3() {
	var p *int
	*p = 3
}

//go:noinline
func ttiWrapper1() *ttiResult {
	var w ttiWrapper
	m := (*ttiWrapper).m1
	return m(&w)
}

type ttiWrapper struct{}

func (w ttiWrapper) m1() *ttiResult {
	return ttiLeaf()
}

//go:noinline
func ttiExcluded1() *ttiResult {
	return ttiExcluded2()
}

// ttiExcluded2 should be excluded from tracebacks. There are
// various ways this could come up. Linking it to a "runtime." name is
// rather synthetic, but it's easy and reliable. See issue #42754 for
// one way this happened in real code.
//
//go:linkname ttiExcluded2 runtime.ttiExcluded2
//go:noinline
func ttiExcluded2() *ttiResult {
	return ttiExcluded3()
}
func ttiExcluded3() *ttiResult {
	return ttiLeaf()
}

var testTracebackArgsBuf [1000]byte

func TestTracebackArgs(t *testing.T) {
	if *flagQuick {
		t.Skip("-quick")
	}
	optimized := !testenv.OptimizationOff()
	abiSel := func(x, y string) string {
		// select expected output based on ABI
		// In noopt build we always spill arguments so the output is the same as stack ABI.
		if optimized && abi.IntArgRegs > 0 {
			return x
		}
		return y
	}

	tests := []struct {
		fn     func() int
		expect string
	}{
		// simple ints
		{
			func() int { return testTracebackArgs1(1, 2, 3, 4, 5) },
			"testTracebackArgs1(0x1, 0x2, 0x3, 0x4, 0x5)",
		},
		// some aggregates
		{
			func() int {
				return testTracebackArgs2(false, struct {
					a, b, c int
					x       [2]int
				}{1, 2, 3, [2]int{4, 5}}, [0]int{}, [3]byte{6, 7, 8})
			},
			"testTracebackArgs2(0x0, {0x1, 0x2, 0x3, {0x4, 0x5}}, {}, {0x6, 0x7, 0x8})",
		},
		{
			func() int { return testTracebackArgs3([3]byte{1, 2, 3}, 4, 5, 6, [3]byte{7, 8, 9}) },
			"testTracebackArgs3({0x1, 0x2, 0x3}, 0x4, 0x5, 0x6, {0x7, 0x8, 0x9})",
		},
		// too deeply nested type
		{
			func() int { return testTracebackArgs4(false, [1][1][1][1][1][1][1][1][1][1]int{}) },
			"testTracebackArgs4(0x0, {{{{{...}}}}})",
		},
		// a lot of zero-sized type
		{
			func() int {
				z := [0]int{}
				return testTracebackArgs5(false, struct {
					x int
					y [0]int
					z [2][0]int
				}{1, z, [2][0]int{}}, z, z, z, z, z, z, z, z, z, z, z, z)
			},
			"testTracebackArgs5(0x0, {0x1, {}, {{}, {}}}, {}, {}, {}, {}, {}, ...)",
		},

		// edge cases for ...
		// no ... for 10 args
		{
			func() int { return testTracebackArgs6a(1, 2, 3, 4, 5, 6, 7, 8, 9, 10) },
			"testTracebackArgs6a(0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa)",
		},
		// has ... for 11 args
		{
			func() int { return testTracebackArgs6b(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11) },
			"testTracebackArgs6b(0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, ...)",
		},
		// no ... for aggregates with 10 words
		{
			func() int { return testTracebackArgs7a([10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) },
			"testTracebackArgs7a({0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa})",
		},
		// has ... for aggregates with 11 words
		{
			func() int { return testTracebackArgs7b([11]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}) },
			"testTracebackArgs7b({0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, ...})",
		},
		// no ... for aggregates, but with more args
		{
			func() int { return testTracebackArgs7c([10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 11) },
			"testTracebackArgs7c({0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa}, ...)",
		},
		// has ... for aggregates and also for more args
		{
			func() int { return testTracebackArgs7d([11]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, 12) },
			"testTracebackArgs7d({0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, ...}, ...)",
		},
		// nested aggregates, no ...
		{
			func() int { return testTracebackArgs8a(testArgsType8a{1, 2, 3, 4, 5, 6, 7, 8, [2]int{9, 10}}) },
			"testTracebackArgs8a({0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, {0x9, 0xa}})",
		},
		// nested aggregates, ... in inner but not outer
		{
			func() int { return testTracebackArgs8b(testArgsType8b{1, 2, 3, 4, 5, 6, 7, 8, [3]int{9, 10, 11}}) },
			"testTracebackArgs8b({0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, {0x9, 0xa, ...}})",
		},
		// nested aggregates, ... in outer but not inner
		{
			func() int { return testTracebackArgs8c(testArgsType8c{1, 2, 3, 4, 5, 6, 7, 8, [2]int{9, 10}, 11}) },
			"testTracebackArgs8c({0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, {0x9, 0xa}, ...})",
		},
		// nested aggregates, ... in both inner and outer
		{
			func() int { return testTracebackArgs8d(testArgsType8d{1, 2, 3, 4, 5, 6, 7, 8, [3]int{9, 10, 11}, 12}) },
			"testTracebackArgs8d({0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, {0x9, 0xa, ...}, ...})",
		},

		// Register argument liveness.
		// 1, 3 are used and live, 2, 4 are dead (in register ABI).
		// Address-taken (7) and stack ({5, 6}) args are always live.
		{
			func() int {
				poisonStack() // poison arg area to make output deterministic
				return testTracebackArgs9(1, 2, 3, 4, [2]int{5, 6}, 7)
			},
			abiSel(
				"testTracebackArgs9(0x1, 0xffffffff?, 0x3, 0xff?, {0x5, 0x6}, 0x7)",
				"testTracebackArgs9(0x1, 0x2, 0x3, 0x4, {0x5, 0x6}, 0x7)"),
		},
		// No live.
		// (Note: this assume at least 5 int registers if register ABI is used.)
		{
			func() int {
				poisonStack() // poison arg area to make output deterministic
				return testTracebackArgs10(1, 2, 3, 4, 5)
			},
			abiSel(
				"testTracebackArgs10(0xffffffff?, 0xffffffff?, 0xffffffff?, 0xffffffff?, 0xffffffff?)",
				"testTracebackArgs10(0x1, 0x2, 0x3, 0x4, 0x5)"),
		},
		// Conditional spills.
		// Spill in conditional, not executed.
		{
			func() int {
				poisonStack() // poison arg area to make output deterministic
				return testTracebackArgs11a(1, 2, 3)
			},
			abiSel(
				"testTracebackArgs11a(0xffffffff?, 0xffffffff?, 0xffffffff?)",
				"testTracebackArgs11a(0x1, 0x2, 0x3)"),
		},
		// 2 spills in conditional, not executed; 3 spills in conditional, executed, but not statically known.
		// So print 0x3?.
		{
			func() int {
				poisonStack() // poison arg area to make output deterministic
				return testTracebackArgs11b(1, 2, 3, 4)
			},
			abiSel(
				"testTracebackArgs11b(0xffffffff?, 0xffffffff?, 0x3?, 0x4)",
				"testTracebackArgs11b(0x1, 0x2, 0x3, 0x4)"),
		},
	}
	for _, test := range tests {
		n := test.fn()
		got := testTracebackArgsBuf[:n]
		if !bytes.Contains(got, []byte(test.expect)) {
			t.Errorf("traceback does not contain expected string: want %q, got\n%s", test.expect, got)
		}
	}
}

//go:noinline
func testTracebackArgs1(a, b, c, d, e int) int {
	n := runtime.Stack(testTracebackArgsBuf[:], false)
	if a < 0 {
		// use in-reg args to keep them alive
		return a + b + c + d + e
	}
	return n
}

//go:noinline
func testTracebackArgs2(a bool, b struct {
	a, b, c int
	x       [2]int
}, _ [0]int, d [3]byte) int {
	n := runtime.Stack(testTracebackArgsBuf[:], false)
	if a {
		// use in-reg args to keep them alive
		return b.a + b.b + b.c + b.x[0] + b.x[1] + int(d[0]) + int(d[1]) + int(d[2])
	}
	return n

}

//go:noinline
//go:registerparams
func testTracebackArgs3(x [3]byte, a, b, c int, y [3]byte) int {
	n := runtime.Stack(testTracebackArgsBuf[:], false)
	if a < 0 {
		// use in-reg args to keep them alive
		return int(x[0]) + int(x[1]) + int(x[2]) + a + b + c + int(y[0]) + int(y[1]) + int(y[2])
	}
	return n
}

//go:noinline
func testTracebackArgs4(a bool, x [1][1][1][1][1][1][1][1][1][1]int) int {
	n := runtime.Stack(testTracebackArgsBuf[:], false)
	if a {
		panic(x) // use args to keep them alive
	}
	return n
}

//go:noinline
func testTracebackArgs5(a bool, x struct {
	x int
	y [0]int
	z [2][0]int
}, _, _, _, _, _, _, _, _, _, _, _, _ [0]int) int {
	n := runtime.Stack(testTracebackArgsBuf[:], false)
	if a {
		panic(x) // use args to keep them alive
	}
	return n
}

//go:noinline
func testTracebackArgs6a(a, b, c, d, e, f, g, h, i, j int) int {
	n := runtime.Stack(testTracebackArgsBuf[:], false)
	if a < 0 {
		// use in-reg args to keep them alive
		return a + b + c + d + e + f + g + h + i + j
	}
	return n
}

//go:noinline
func testTracebackArgs6b(a, b, c, d, e, f, g, h, i, j, k int) int {
	n := runtime.Stack(testTracebackArgsBuf[:], false)
	if a < 0 {
		// use in-reg args to keep them alive
		return a + b + c + d + e + f + g + h + i + j + k
	}
	return n
}

//go:noinline
func testTracebackArgs7a(a [10]int) int {
	n := runtime.Stack(testTracebackArgsBuf[:], false)
	if a[0] < 0 {
		// use in-reg args to keep them alive
		return a[1] + a[2] + a[3] + a[4] + a[5] + a[6] + a[7] + a[8] + a[9]
	}
	return n
}

//go:noinline
func testTracebackArgs7b(a [11]int) int {
	n := runtime.Stack(testTracebackArgsBuf[:], false)
	if a[0] < 0 {
		// use in-reg args to keep them alive
		return a[1] + a[2] + a[3] + a[4] + a[5] + a[6] + a[7] + a[8] + a[9] + a[10]
	}
	return n
}

//go:noinline
func testTracebackArgs7c(a [10]int, b int) int {
	n := runtime.Stack(testTracebackArgsBuf[:], false)
	if a[0] < 0 {
		// use in-reg args to keep them alive
		return a[1] + a[2] + a[3] + a[4] + a[5] + a[6] + a[7] + a[8] + a[9] + b
	}
	return n
}

//go:noinline
func testTracebackArgs7d(a [11]int, b int) int {
	n := runtime.Stack(testTracebackArgsBuf[:], false)
	if a[0] < 0 {
		// use in-reg args to keep them alive
		return a[1] + a[2] + a[3] + a[4] + a[5] + a[6] + a[7] + a[8] + a[9] + a[10] + b
	}
	return n
}

type testArgsType8a struct {
	a, b, c, d, e, f, g, h int
	i                      [2]int
}
type testArgsType8b struct {
	a, b, c, d, e, f, g, h int
	i                      [3]int
}
type testArgsType8c struct {
	a, b, c, d, e, f, g, h int
	i                      [2]int
	j                      int
}
type testArgsType8d struct {
	a, b, c, d, e, f, g, h int
	i                      [3]int
	j                      int
}

//go:noinline
func testTracebackArgs8a(a testArgsType8a) int {
	n := runtime.Stack(testTracebackArgsBuf[:], false)
	if a.a < 0 {
		// use in-reg args to keep them alive
		return a.b + a.c + a.d + a.e + a.f + a.g + a.h + a.i[0] + a.i[1]
	}
	return n
}

//go:noinline
func testTracebackArgs8b(a testArgsType8b) int {
	n := runtime.Stack(testTracebackArgsBuf[:], false)
	if a.a < 0 {
		// use in-reg args to keep them alive
		return a.b + a.c + a.d + a.e + a.f + a.g + a.h + a.i[0] + a.i[1] + a.i[2]
	}
	return n
}

//go:noinline
func testTracebackArgs8c(a testArgsType8c) int {
	n := runtime.Stack(testTracebackArgsBuf[:], false)
	if a.a < 0 {
		// use in-reg args to keep them alive
		return a.b + a.c + a.d + a.e + a.f + a.g + a.h + a.i[0] + a.i[1] + a.j
	}
	return n
}

//go:noinline
func testTracebackArgs8d(a testArgsType8d) int {
	n := runtime.Stack(testTracebackArgsBuf[:], false)
	if a.a < 0 {
		// use in-reg args to keep them alive
		return a.b + a.c + a.d + a.e + a.f + a.g + a.h + a.i[0] + a.i[1] + a.i[2] + a.j
	}
	return n
}

// nosplit to avoid preemption or morestack spilling registers.
//
//go:nosplit
//go:noinline
func testTracebackArgs9(a int64, b int32, c int16, d int8, x [2]int, y int) int {
	if a < 0 {
		println(&y) // take address, make y live, even if no longer used at traceback
	}
	n := runtime.Stack(testTracebackArgsBuf[:], false)
	if a < 0 {
		// use half of in-reg args to keep them alive, the other half are dead
		return int(a) + int(c)
	}
	return n
}

// nosplit to avoid preemption or morestack spilling registers.
//
//go:nosplit
//go:noinline
func testTracebackArgs10(a, b, c, d, e int32) int {
	// no use of any args
	return runtime.Stack(testTracebackArgsBuf[:], false)
}

// norace to avoid race instrumentation changing spill locations.
// nosplit to avoid preemption or morestack spilling registers.
//
//go:norace
//go:nosplit
//go:noinline
func testTracebackArgs11a(a, b, c int32) int {
	if a < 0 {
		println(a, b, c) // spill in a conditional, may not execute
	}
	if b < 0 {
		return int(a + b + c)
	}
	return runtime.Stack(testTracebackArgsBuf[:], false)
}

// norace to avoid race instrumentation changing spill locations.
// nosplit to avoid preemption or morestack spilling registers.
//
//go:norace
//go:nosplit
//go:noinline
func testTracebackArgs11b(a, b, c, d int32) int {
	var x int32
	if a < 0 {
		print() // spill b in a conditional
		x = b
	} else {
		print() // spill c in a conditional
		x = c
	}
	if d < 0 { // d is always needed
		return int(x + d)
	}
	return runtime.Stack(testTracebackArgsBuf[:], false)
}

// Poison the arg area with deterministic values.
//
//go:noinline
func poisonStack() [20]int {
	return [20]int{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}
}

func TestTracebackParentChildGoroutines(t *testing.T) {
	parent := fmt.Sprintf("goroutine %d", runtime.Goid())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		buf := make([]byte, 1<<10)
		// We collect the stack only for this goroutine (by passing
		// false to runtime.Stack). We expect to see the current
		// goroutine ID, and the parent goroutine ID in a message like
		// "created by ... in goroutine N".
		stack := string(buf[:runtime.Stack(buf, false)])
		child := fmt.Sprintf("goroutine %d", runtime.Goid())
		if !strings.Contains(stack, parent) || !strings.Contains(stack, child) {
			t.Errorf("did not see parent (%s) and child (%s) IDs in stack, got %s", parent, child, stack)
		}
	}()
	wg.Wait()
}

type traceback struct {
	frames    []*tbFrame
	createdBy *tbFrame // no args
}

type tbFrame struct {
	funcName string
	args     string
	inlined  bool
}

// parseTraceback parses a printed traceback to make it easier for tests to
// check the result.
func parseTraceback(t *testing.T, tb string) []*traceback {
	lines := strings.Split(tb, "\n")
	nLines := len(lines)
	fatal := func(f string, args ...any) {
		lineNo := nLines - len(lines) + 1
		msg := fmt.Sprintf(f, args...)
		t.Fatalf("%s (line %d):\n%s", msg, lineNo, tb)
	}
	parseFrame := func(funcName, args string) *tbFrame {
		// Consume file/line/etc
		if len(lines) == 0 || !strings.HasPrefix(lines[0], "\t") {
			fatal("missing source line")
		}
		lines = lines[1:]
		inlined := args == "..."
		return &tbFrame{funcName, args, inlined}
	}
	var tbs []*traceback
	var cur *traceback
	for len(lines) > 0 {
		line := lines[0]
		lines = lines[1:]
		switch {
		case strings.HasPrefix(line, "goroutine "):
			cur = &traceback{}
			tbs = append(tbs, cur)
		case line == "":
			cur = nil
		case line[0] == '\t':
			fatal("unexpected indent")
		case strings.HasPrefix(line, "created by "):
			funcName := line[len("created by "):]
			cur.createdBy = parseFrame(funcName, "")
		case strings.HasSuffix(line, ")"):
			line = line[:len(line)-1] // Trim trailing ")"
			funcName, args, found := strings.Cut(line, "(")
			if !found {
				fatal("missing (")
			}
			frame := parseFrame(funcName, args)
			cur.frames = append(cur.frames, frame)
		}
	}
	return tbs
}

// parseTraceback1 is like parseTraceback, but expects tb to contain exactly one
// goroutine.
func parseTraceback1(t *testing.T, tb string) *traceback {
	tbs := parseTraceback(t, tb)
	if len(tbs) != 1 {
		t.Fatalf("want 1 goroutine, got %d:\n%s", len(tbs), tb)
	}
	return tbs[0]
}
