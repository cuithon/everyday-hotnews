// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !netbsd

package stdio

/*
#include <stdio.h>
*/
import "C"

var Stdout = (*File)(C.stdout)
var Stderr = (*File)(C.stderr)
