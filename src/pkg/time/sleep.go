// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

// Interface to timers implemented in package runtime.
// Must be in sync with ../runtime/runtime.h:/^struct.Timer$
type runtimeTimer struct {
	i      int32
	when   int64
	period int64
	f      func(int64, interface{})
	arg    interface{}
}

func startTimer(*runtimeTimer)
func stopTimer(*runtimeTimer) bool

// The Timer type represents a single event.
// When the Timer expires, the current time will be sent on C,
// unless the Timer was created by AfterFunc.
type Timer struct {
	C <-chan int64
	r runtimeTimer
}

// Stop prevents the Timer from firing.
// It returns true if the call stops the timer, false if the timer has already
// expired or stopped.
func (t *Timer) Stop() (ok bool) {
	return stopTimer(&t.r)
}

// NewTimer creates a new Timer that will send
// the current time on its channel after at least ns nanoseconds.
func NewTimer(ns int64) *Timer {
	c := make(chan int64, 1)
	t := &Timer{
		C: c,
		r: runtimeTimer{
			when: Nanoseconds() + ns,
			f:    sendTime,
			arg:  c,
		},
	}
	startTimer(&t.r)
	return t
}

func sendTime(now int64, c interface{}) {
	// Non-blocking send of time on c.
	// Used in NewTimer, it cannot block anyway (buffer).
	// Used in NewTicker, dropping sends on the floor is
	// the desired behavior when the reader gets behind,
	// because the sends are periodic.
	select {
	case c.(chan int64) <- now:
	default:
	}
}

// After waits at least ns nanoseconds before sending the current time
// on the returned channel.
// It is equivalent to NewTimer(ns).C.
func After(ns int64) <-chan int64 {
	return NewTimer(ns).C
}

// AfterFunc waits at least ns nanoseconds before calling f
// in its own goroutine. It returns a Timer that can
// be used to cancel the call using its Stop method.
func AfterFunc(ns int64, f func()) *Timer {
	t := &Timer{
		r: runtimeTimer{
			when: Nanoseconds() + ns,
			f:    goFunc,
			arg:  f,
		},
	}
	startTimer(&t.r)
	return t
}

func goFunc(now int64, arg interface{}) {
	go arg.(func())()
}
