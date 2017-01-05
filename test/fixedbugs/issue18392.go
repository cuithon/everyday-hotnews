// errorcheck

// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p

type A interface {
	Fn(A.Fn) // ERROR "type A has no method A.Fn"
}
