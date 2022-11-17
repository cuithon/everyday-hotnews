// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package a

type A struct {
	New func() any
}

func NewA(i int) *A {
	return &A{
		New: func() any {
			_ = i
			return nil
		},
	}
}
