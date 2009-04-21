// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exvar

import (
	"exvar";
	"fmt";
	"testing";
)

func TestSimpleCounter(t *testing.T) {
	// Unknown exvar should be zero.
	x := GetInt("requests");
	if x != 0 {
		t.Errorf("Get(nonexistent) = %v, want 0", x)
	}

	IncrementInt("requests", 1);
	IncrementInt("requests", 3);
	x = GetInt("requests");
	if x != 4 {
		t.Errorf("Get('requests') = %v, want 4", x)
	}

	out := String();
	if out != "requests 4\n" {
		t.Errorf("String() = \"%v\", want \"requests 4\n\"",
		         out);
	}
}

func TestMismatchedCounters(t *testing.T) {
	// Make sure some vars exist.
	GetInt("requests");
	GetMap("colours", "red");

	IncrementInt("colours", 1);
	if x := GetInt("x-mismatched-int"); x != 1 {
		t.Errorf("GetInt('x-mismatched-int') = %v, want 1", x)
	}

	IncrementMap("requests", "orange", 1);
	if x := GetMap("x-mismatched-map", "orange"); x != 1 {
		t.Errorf("GetMap('x-mismatched-int', 'orange') = %v, want 1", x)
	}
}

func TestMapCounter(t *testing.T) {
	// Unknown exvar should be zero.
	if x := GetMap("colours", "red"); x != 0 {
		t.Errorf("GetMap(non, existent) = %v, want 0", x)
	}

	IncrementMap("colours", "red", 1);
	IncrementMap("colours", "red", 2);
	IncrementMap("colours", "blue", 4);
	if x := GetMap("colours", "red"); x != 3 {
		t.Errorf("GetMap('colours', 'red') = %v, want 3", x)
	}
	if x := GetMap("colours", "blue"); x != 4 {
		t.Errorf("GetMap('colours', 'blue') = %v, want 4", x)
	}

	// TODO(dsymonds): Test String()
}

func hammer(name string, total int, done chan <- int) {
	for i := 0; i < total; i++ {
		IncrementInt(name, 1)
	}
	done <- 1
}

func TestHammer(t *testing.T) {
	SetInt("hammer-times", 0);
	sync := make(chan int);
	hammer_times := int(1e5);
	go hammer("hammer-times", hammer_times, sync);
	go hammer("hammer-times", hammer_times, sync);
	<-sync;
	<-sync;
	if final := GetInt("hammer-times"); final != 2 * hammer_times {
		t.Errorf("hammer-times = %v, want %v", final, 2 * hammer_times)
	}
}
