// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build mips64 || mips64le
// +build mips64 mips64le

package unix

const (
	getrandomTrap     uintptr = 5313
	copyFileRangeTrap uintptr = 5320
)
