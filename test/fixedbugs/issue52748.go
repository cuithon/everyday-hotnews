// errorcheck

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package p

import "unsafe"

type S[T any] struct{}

const c = unsafe.Sizeof(S[[c]byte]{}) // ERROR "initialization loop"
