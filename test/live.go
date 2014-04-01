// errorcheck -0 -l -live

// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

func f1() {
	var x *int
	print(&x) // ERROR "live at call to printpointer: x$"
	print(&x) // ERROR "live at call to printpointer: x$"
}

func f2(b bool) {
	if b {
		print(0) // nothing live here
		return
	}
	var x *int
	print(&x) // ERROR "live at call to printpointer: x$"
	print(&x) // ERROR "live at call to printpointer: x$"
}

func f3(b bool) {
	print(0)
	if b == false {
		print(0) // nothing live here
		return
	}

	if b {
		var x *int
		print(&x) // ERROR "live at call to printpointer: x$"
		print(&x) // ERROR "live at call to printpointer: x$"
	} else {
		var y *int
		print(&y) // ERROR "live at call to printpointer: y$"
		print(&y) // ERROR "live at call to printpointer: y$"
	}
	print(0) // ERROR "live at call to printint: x y$" "x \(type \*int\) is ambiguously live" "y \(type \*int\) is ambiguously live"
}

// The old algorithm treated x as live on all code that
// could flow to a return statement, so it included the
// function entry and code above the declaration of x
// but would not include an indirect use of x in an infinite loop.
// Check that these cases are handled correctly.

func f4(b1, b2 bool) { // x not live here
	if b2 {
		print(0) // x not live here
		return
	}
	var z **int
	x := new(int)
	*x = 42
	z = &x
	print(**z) // ERROR "live at call to printint: x z$"
	if b2 {
		print(1) // ERROR "live at call to printint: x$"
		return
	}
	for {
		print(**z) // ERROR "live at call to printint: x z$"
	}
}

func f5(b1 bool) {
	var z **int
	if b1 {
		x := new(int)
		*x = 42
		z = &x
	} else {
		y := new(int)
		*y = 54
		z = &y
	}
	print(**z) // ERROR "live at call to printint: x y$" "x \(type \*int\) is ambiguously live" "y \(type \*int\) is ambiguously live"
}

// confusion about the _ result used to cause spurious "live at entry to f6: _".

func f6() (_, y string) {
	y = "hello"
	return
}

// confusion about addressed results used to cause "live at entry to f7: x".

func f7() (x string) {
	_ = &x
	x = "hello"
	return
}

// ignoring block returns used to cause "live at entry to f8: x, y".

func f8() (x, y string) {
	return g8()
}

func g8() (string, string)

// ignoring block assignments used to cause "live at entry to f9: x"
// issue 7205

var i9 interface{}

func f9() bool {
	g8()
	x := i9
	return x != 99
}

// liveness formerly confused by UNDEF followed by RET,
// leading to "live at entry to f10: ~r1" (unnamed result).

func f10() string {
	panic(1)
}

// liveness formerly confused by select, thinking runtime.selectgo
// can return to next instruction; it always jumps elsewhere.
// note that you have to use at least two cases in the select
// to get a true select; smaller selects compile to optimized helper functions.

var c chan *int
var b bool

// this used to have a spurious "live at entry to f11a: ~r0"
func f11a() *int {
	select { // ERROR "live at call to selectgo: autotmp"
	case <-c: // ERROR "live at call to selectrecv: autotmp"
		return nil
	case <-c: // ERROR "live at call to selectrecv: autotmp"
		return nil
	}
}

func f11b() *int {
	p := new(int)
	if b {
		// At this point p is dead: the code here cannot
		// get to the bottom of the function.
		// This used to have a spurious "live at call to printint: p".
		print(1) // nothing live here!
		select { // ERROR "live at call to selectgo: autotmp"
		case <-c: // ERROR "live at call to selectrecv: autotmp"
			return nil
		case <-c: // ERROR "live at call to selectrecv: autotmp"
			return nil
		}
	}
	println(*p)
	return nil
}

func f11c() *int {
	p := new(int)
	if b {
		// Unlike previous, the cases in this select fall through,
		// so we can get to the println, so p is not dead.
		print(1) // ERROR "live at call to printint: p"
		select { // ERROR "live at call to newselect: p" "live at call to selectgo: autotmp.* p"
		case <-c: // ERROR "live at call to selectrecv: autotmp.* p"
		case <-c: // ERROR "live at call to selectrecv: autotmp.* p"
		}
	}
	println(*p)
	return nil
}

// similarly, select{} does not fall through.
// this used to have a spurious "live at entry to f12: ~r0".

func f12() *int {
	if b {
		select{}
	} else {
		return nil
	}
}

// incorrectly placed VARDEF annotations can cause missing liveness annotations.
// this used to be missing the fact that s is live during the call to g13 (because it is
// needed for the call to h13).

func f13() {
	s := "hello"
	s = h13(s, g13(s)) // ERROR "live at call to g13: s"
}

func g13(string) string
func h13(string, string) string

// more incorrectly placed VARDEF.

func f14() {
	x := g14()
	print(&x) // ERROR "live at call to printpointer: x"
}

func g14() string

func f15() {
	var x string
	_ = &x
	x = g15() // ERROR "live at call to g15: x"
	print(x) // ERROR "live at call to printstring: x"
}

func g15() string

// Checking that various temporaries do not persist or cause
// ambiguously live values that must be zeroed.
// The exact temporary names are inconsequential but we are
// trying to check that there is only one at any given site,
// and also that none show up in "ambiguously live" messages.

var m map[string]int

func f16() {
	if b {
		delete(m, "hi") // ERROR "live at call to mapdelete: autotmp_[0-9]+$"
	}
	delete(m, "hi") // ERROR "live at call to mapdelete: autotmp_[0-9]+$"
	delete(m, "hi") // ERROR "live at call to mapdelete: autotmp_[0-9]+$"
}

var m2s map[string]*byte
var m2 map[[2]string]*byte
var x2 [2]string
var bp *byte

func f17a() {
	// value temporary only
	if b {
		m2[x2] = nil // ERROR "live at call to mapassign1: autotmp_[0-9]+$"
	}
	m2[x2] = nil // ERROR "live at call to mapassign1: autotmp_[0-9]+$"
	m2[x2] = nil // ERROR "live at call to mapassign1: autotmp_[0-9]+$"
}

func f17b() {
	// key temporary only
	if b {
		m2s["x"] = bp // ERROR "live at call to mapassign1: autotmp_[0-9]+$"
	}
	m2s["x"] = bp // ERROR "live at call to mapassign1: autotmp_[0-9]+$"
	m2s["x"] = bp // ERROR "live at call to mapassign1: autotmp_[0-9]+$"
}

func f17c() {
	// key and value temporaries
	if b {
		m2s["x"] = nil // ERROR "live at call to mapassign1: autotmp_[0-9]+ autotmp_[0-9]+$"
	}
	m2s["x"] = nil // ERROR "live at call to mapassign1: autotmp_[0-9]+ autotmp_[0-9]+$"
	m2s["x"] = nil // ERROR "live at call to mapassign1: autotmp_[0-9]+ autotmp_[0-9]+$"
}

func g18() [2]string

func f18() {
	// key temporary for mapaccess.
	// temporary introduced by orderexpr.
	var z *byte
	if b {
		z = m2[g18()] // ERROR "live at call to mapaccess1: autotmp_[0-9]+$"
	}
	z = m2[g18()] // ERROR "live at call to mapaccess1: autotmp_[0-9]+$"
	z = m2[g18()] // ERROR "live at call to mapaccess1: autotmp_[0-9]+$"
	print(z)
}

var ch chan *byte

func f19() {
	// dest temporary for channel receive.
	var z *byte
	
	if b {
		z = <-ch // ERROR "live at call to chanrecv1: autotmp_[0-9]+$"
	}
	z = <-ch // ERROR "live at call to chanrecv1: autotmp_[0-9]+$"
	z = <-ch // ERROR "live at call to chanrecv1: autotmp_[0-9]+$"
	print(z)
}

func f20() {
	// src temporary for channel send
	if b {
		ch <- nil // ERROR "live at call to chansend1: autotmp_[0-9]+$"
	}
	ch <- nil // ERROR "live at call to chansend1: autotmp_[0-9]+$"
	ch <- nil // ERROR "live at call to chansend1: autotmp_[0-9]+$"
}

func f21() {
	// key temporary for mapaccess using array literal key.
	var z *byte
	if b {
		z = m2[[2]string{"x", "y"}] // ERROR "live at call to mapaccess1: autotmp_[0-9]+$"
	}
	z = m2[[2]string{"x", "y"}] // ERROR "live at call to mapaccess1: autotmp_[0-9]+$"
	z = m2[[2]string{"x", "y"}] // ERROR "live at call to mapaccess1: autotmp_[0-9]+$"
	print(z)
}

func f23() {
	// key temporary for two-result map access using array literal key.
	var z *byte
	var ok bool
	if b {
		z, ok = m2[[2]string{"x", "y"}] // ERROR "live at call to mapaccess2: autotmp_[0-9]+$"
	}
	z, ok = m2[[2]string{"x", "y"}] // ERROR "live at call to mapaccess2: autotmp_[0-9]+$"
	z, ok = m2[[2]string{"x", "y"}] // ERROR "live at call to mapaccess2: autotmp_[0-9]+$"
	print(z, ok)
}

func f24() {
	// key temporary for map access using array literal key.
	// value temporary too.
	if b {
		m2[[2]string{"x", "y"}] = nil // ERROR "live at call to mapassign1: autotmp_[0-9]+ autotmp_[0-9]+$"
	}
	m2[[2]string{"x", "y"}] = nil // ERROR "live at call to mapassign1: autotmp_[0-9]+ autotmp_[0-9]+$"
	m2[[2]string{"x", "y"}] = nil // ERROR "live at call to mapassign1: autotmp_[0-9]+ autotmp_[0-9]+$"
}
