// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package Package

import Globals "globals"

export Package
type Package struct {
	ref int;
	ident string;
	path string;
	key string;
	scope *Globals.Scope;
}


export NewPackage;
func NewPackage() *Package {
	pkg := new(Package);
	return pkg;
}
