// $G $D/$F.go && $L $F.$A && ./$A.out

// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tests verifying the semantics of the select statement
// for basic empty/non-empty cases.

package main

import "time"

const always = "function did not"
const never = "function did"


func unreachable() {
	panic("control flow shouldn't reach here")
}


// Calls f and verifies that f always/never panics depending on signal.
func testPanic(signal string, f func()) {
	defer func() {
		s := never
		if recover() != nil {
			s = always // f panicked
		}
		if s != signal {
			panic(signal + " panic")
		}
	}()
	f()
}


// Calls f and empirically verifies that f always/never blocks depending on signal.
func testBlock(signal string, f func()) {
	c := make(chan string)
	go func() {
		f()
		c <- never // f didn't block
	}()
	go func() {
		time.Sleep(1e8) // 0.1s seems plenty long
		c <- always     // f blocked always
	}()
	if <-c != signal {
		panic(signal + " block")
	}
}


func main() {
	const async = 1 // asynchronous channels
	var nilch chan int
	closedch := make(chan int)
	close(closedch)

	// sending/receiving from a nil channel outside a select panics
	testPanic(always, func() {
		nilch <- 7
	})
	testPanic(always, func() {
		<-nilch
	})

	// sending/receiving from a nil channel inside a select never panics
	testPanic(never, func() {
		select {
		case nilch <- 7:
			unreachable()
		default:
		}
	})
	testPanic(never, func() {
		select {
		case <-nilch:
			unreachable()
		default:
		}
	})

	// sending to an async channel with free buffer space never blocks
	testBlock(never, func() {
		ch := make(chan int, async)
		ch <- 7
	})

	// receiving (a small number of times) from a closed channel never blocks
	testBlock(never, func() {
		for i := 0; i < 10; i++ {
			if <-closedch != 0 {
				panic("expected zero value when reading from closed channel")
			}
		}
	})

	// sending (a small number of times) to a closed channel is not specified
	// but the current implementation doesn't block: test that different
	// implementations behave the same
	testBlock(never, func() {
		for i := 0; i < 10; i++ {
			closedch <- 7
		}
	})

	// receiving from a non-ready channel always blocks
	testBlock(always, func() {
		ch := make(chan int)
		<-ch
	})

	// empty selects always block
	testBlock(always, func() {
		select {
		}
	})

	// selects with only nil channels always block
	testBlock(always, func() {
		select {
		case <-nilch:
			unreachable()
		}
	})
	testBlock(always, func() {
		select {
		case nilch <- 7:
			unreachable()
		}
	})
	testBlock(always, func() {
		select {
		case <-nilch:
			unreachable()
		case nilch <- 7:
			unreachable()
		}
	})

	// selects with non-ready non-nil channels always block
	testBlock(always, func() {
		ch := make(chan int)
		select {
		case <-ch:
			unreachable()
		}
	})

	// selects with default cases don't block
	testBlock(never, func() {
		select {
		default:
		}
	})
	testBlock(never, func() {
		select {
		case <-nilch:
			unreachable()
		default:
		}
	})
	testBlock(never, func() {
		select {
		case nilch <- 7:
			unreachable()
		default:
		}
	})

	// selects with ready channels don't block
	testBlock(never, func() {
		ch := make(chan int, async)
		select {
		case ch <- 7:
		default:
			unreachable()
		}
	})
	testBlock(never, func() {
		ch := make(chan int, async)
		ch <- 7
		select {
		case <-ch:
		default:
			unreachable()
		}
	})

	// selects with closed channels don't block
	testBlock(never, func() {
		select {
		case <-closedch:
		}
	})
	testBlock(never, func() {
		select {
		case closedch <- 7:
		}
	})
}
