// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package a

type I interface {
	I2
}
type I2 interface {
	M()
}
type S struct{}

func (*S) M() {}

func New() I {
	return &S{}
}
