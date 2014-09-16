// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync_test

import (
	. "sync"
	"testing"
)

type one int

func (o *one) Increment() {
	*o++
}

func run(t *testing.T, once *Once, o *one, c chan bool) {
	once.Do(func() { o.Increment() })
	if v := *o; v != 1 {
		t.Errorf("once failed inside run: %d is not 1", v)
	}
	c <- true
}

func TestOnce(t *testing.T) {
	o := new(one)
	once := new(Once)
	c := make(chan bool)
	const N = 10
	for i := 0; i < N; i++ {
		go run(t, once, o, c)
	}
	for i := 0; i < N; i++ {
		<-c
	}
	if *o != 1 {
		t.Errorf("once failed outside run: %d is not 1", *o)
	}
}

func TestOncePanic(t *testing.T) {
	once := new(Once)
	for i := 0; i < 2; i++ {
		func() {
			defer func() {
				r := recover()
				if r == nil && i == 0 {
					t.Fatalf("Once.Do() has not panic'ed on first iteration")
				}
				if r != nil && i == 1 {
					t.Fatalf("Once.Do() has panic'ed on second iteration")
				}
			}()
			once.Do(func() {
				panic("failed")
			})
		}()
	}
	once.Do(func() {})
	once.Do(func() {
		t.Fatalf("Once called twice")
	})
}

func BenchmarkOnce(b *testing.B) {
	var once Once
	f := func() {}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			once.Do(f)
		}
	})
}
