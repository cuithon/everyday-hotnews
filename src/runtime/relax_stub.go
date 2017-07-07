// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows

package runtime

// osRelaxDelay is the number of nanoseconds of idleness to tolerate
// before performing an osRelax. Since osRelax may reduce the
// precision of timers, this should be enough larger than the relaxed
// timer precision to keep the timer error acceptable.
const osRelaxDelay = 0

// osRelax is called by the scheduler when transitioning to and from
// all Ps being idle.
func osRelax(relax bool) {}
