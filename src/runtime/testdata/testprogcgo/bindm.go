// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9 && !windows

// Test that callbacks from C to Go in the same C-thread always get the same m.
// Make sure the extra M bind to the C-thread.

package main

/*
#include <stdint.h>
#include <pthread.h>
#include <unistd.h>

extern void GoCheckBindM();

#define CTHREADS 2
#define CHECKCALLS 100

static void* checkBindMThread(void* thread) {
	int i;
	for (i = 0; i < CHECKCALLS; i++) {
		GoCheckBindM((uintptr_t)thread);
		usleep(1);
	}
	return NULL;
}

static void CheckBindM() {
	int i;
	pthread_t s[CTHREADS];

	for (i = 0; i < CTHREADS; i++) {
		pthread_create(&s[i], NULL, checkBindMThread, &s[i]);
	}
	for (i = 0; i < CTHREADS; i++) {
		pthread_join(s[i], NULL);
	}
}
*/
import "C"

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
)

var (
	mutex      = sync.Mutex{}
	cThreadToM = map[uintptr]uintptr{}
	started    = atomic.Uint32{}
)

// same as CTHREADS in C, make sure all the C threads are actually started.
const cThreadNum = 2

func init() {
	register("EnsureBindM", EnsureBindM)
}

//export GoCheckBindM
func GoCheckBindM(thread uintptr) {
	// Wait all threads start
	if started.Load() != cThreadNum {
		// Only once for each thread, since it will wait all threads start.
		started.Add(1)
		for started.Load() < cThreadNum {
			runtime.Gosched()
		}
	}
	m := runtime_getm_for_test()
	mutex.Lock()
	defer mutex.Unlock()
	if savedM, ok := cThreadToM[thread]; ok && savedM != m {
		fmt.Printf("m == %x want %x\n", m, savedM)
		os.Exit(1)
	}
	cThreadToM[thread] = m
}

func EnsureBindM() {
	C.CheckBindM()
	fmt.Println("OK")
}
