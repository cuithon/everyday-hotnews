// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
)

/*
 * Helpers for building runtime.
 */

// mkzversion writes zversion.go:
//
//	package runtime
//	const defaultGoroot = <goroot>
//	const theVersion = <version>
//
func mkzversion(dir, file string) {
	out := fmt.Sprintf(
		"// auto generated by go tool dist\n"+
			"\n"+
			"package runtime\n"+
			"\n"+
			"const defaultGoroot = `%s`\n"+
			"const theVersion = `%s`\n"+
			"\n"+
			"var buildVersion = theVersion\n", goroot_final, goversion)

	writefile(out, file, 0)
}

// mkzexperiment writes zaexperiment.h (sic):
//
//	#define GOEXPERIMENT "experiment string"
//
func mkzexperiment(dir, file string) {
	out := fmt.Sprintf(
		"// auto generated by go tool dist\n"+
			"\n"+
			"#define GOEXPERIMENT \"%s\"\n", os.Getenv("GOEXPERIMENT"))

	writefile(out, file, 0)
}
