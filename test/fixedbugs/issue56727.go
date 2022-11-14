// compile

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p

type I interface {
	M()
}

type S struct{}

func (*S) M() {}

type slice []I

func f() {
	ss := struct {
		i I
	}{
		i: &S{},
	}

	_ = [...]struct {
		s slice
	}{
		{
			s: slice{ss.i},
		},
		{
			s: slice{ss.i},
		},
		{
			s: slice{ss.i},
		},
		{
			s: slice{ss.i},
		},
		{
			s: slice{ss.i},
		},
	}
}
