// errorcheck

// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Verify that we get a single non-confusing error
// message for embedded fields/interfaces that use
// a qualified identifier with non-existing package.

package p

type _ struct {
	F sync.Mutex // ERROR "undefined: sync"
}

type _ struct {
	sync.Mutex // ERROR "undefined: sync"
}

type _ interface {
	sync.Mutex // ERROR "undefined: sync"
}
