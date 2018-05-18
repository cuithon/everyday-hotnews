// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package issue12839 is a go/doc test to test association of a function
// that returns multiple types.
// See golang.org/issue/12839.
package issue12839

import "p"

type T1 struct{}

type T2 struct{}

func (t T1) hello() string {
	return "hello"
}

// F1 should not be associated with T1
func F1() (*T1, *T2) {
	return &T1{}, &T2{}
}

// F2 should be associated with T1
func F2() (a, b, c T1) {
	return T1{}, T1{}, T1{}
}

// F3 should be associated with T1 because b.T3 is from a different package
func F3() (a T1, b p.T3) {
	return T1{}, p.T3{}
}

// F4 should not be associated with a type (same as F1)
func F4() (a T1, b T2) {
	return T1{}, T2{}
}
