// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types2

// This file will eventually define an Alias type.
// For now it declares asNamed. Once Alias types
// exist, asNamed will need to indirect through
// them as needed.

// asNamed returns t as *Named if that is t's
// actual type. It returns nil otherwise.
func asNamed(t Type) *Named {
	n, _ := t.(*Named)
	return n
}
