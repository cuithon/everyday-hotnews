// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"os"
	"path/filepath"
)

var Cwd, _ = os.Getwd()

// ShortPath returns an absolute or relative name for path, whatever is shorter.
func ShortPath(path string) string {
	if rel, err := filepath.Rel(Cwd, path); err == nil && len(rel) < len(path) {
		return rel
	}
	return path
}

// RelPaths returns a copy of paths with absolute paths
// made relative to the current directory if they would be shorter.
func RelPaths(paths []string) []string {
	var out []string
	// TODO(rsc): Can this use Cwd from above?
	pwd, _ := os.Getwd()
	for _, p := range paths {
		rel, err := filepath.Rel(pwd, p)
		if err == nil && len(rel) < len(p) {
			p = rel
		}
		out = append(out, p)
	}
	return out
}
