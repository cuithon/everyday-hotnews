// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build cgo

package cgotest

// Issue 43639: No runtime test needed, make sure package misc/cgo/test/issue43639 compiles well.

import _ "misc/cgo/test/issue43639"
