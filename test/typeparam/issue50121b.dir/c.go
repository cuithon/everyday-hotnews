// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package c

import (
	"./b"
)

func BuildInt() int {
	return b.IntBuilder.New()
}
