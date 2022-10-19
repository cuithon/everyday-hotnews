// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !unix

package runtime

const canCreateFile = false

func create(name *byte, perm int32) int32 {
	throw("unimplemented")
	return -1
}
