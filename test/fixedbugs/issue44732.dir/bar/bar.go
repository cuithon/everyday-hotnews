// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bar

import "issue44732.dir/foo"

type Bar struct {
	Foo *foo.Foo
}
