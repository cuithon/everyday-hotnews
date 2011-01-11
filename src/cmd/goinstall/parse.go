// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Wrappers for Go parser.

package main

import (
	"path"
	"os"
	"log"
	"strings"
	"strconv"
	"go/ast"
	"go/parser"
)


type dirInfo struct {
	goFiles  []string // .go files within dir (including cgoFiles)
	cgoFiles []string // .go files that import "C"
	cFiles   []string // .c files within dir
	imports  []string // All packages imported by goFiles
	pkgName  string   // Name of package within dir
}

// scanDir returns a structure with details about the Go content found
// in the given directory. The list of files will NOT contain the
// following entries:
//
// - Files in package main (unless allowMain is true)
// - Files ending in _test.go
// - Files starting with _ (temporary)
// - Files containing .cgo in their names
//
// The imports map keys are package paths imported by listed Go files,
// and the values are the Go files importing the respective package paths.
func scanDir(dir string, allowMain bool) (info *dirInfo, err os.Error) {
	f, err := os.Open(dir, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	dirs, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}

	goFiles := make([]string, 0, len(dirs))
	cgoFiles := make([]string, 0, len(dirs))
	cFiles := make([]string, 0, len(dirs))
	importsm := make(map[string]bool)
	pkgName := ""
	for i := range dirs {
		d := &dirs[i]
		if strings.HasPrefix(d.Name, "_") || strings.Index(d.Name, ".cgo") != -1 {
			continue
		}
		if strings.HasSuffix(d.Name, ".c") {
			cFiles = append(cFiles, d.Name)
			continue
		}
		if !strings.HasSuffix(d.Name, ".go") || strings.HasSuffix(d.Name, "_test.go") {
			continue
		}
		filename := path.Join(dir, d.Name)
		pf, err := parser.ParseFile(fset, filename, nil, parser.ImportsOnly)
		if err != nil {
			return nil, err
		}
		s := string(pf.Name.Name)
		if s == "main" && !allowMain {
			continue
		}
		if pkgName == "" {
			pkgName = s
		} else if pkgName != s {
			// Only if all files in the directory are in package main
			// do we return pkgName=="main".
			// A mix of main and another package reverts
			// to the original (allowMain=false) behaviour.
			if s == "main" || pkgName == "main" {
				return scanDir(dir, false)
			}
			return nil, os.ErrorString("multiple package names in " + dir)
		}
		goFiles = append(goFiles, d.Name)
		for _, decl := range pf.Decls {
			for _, spec := range decl.(*ast.GenDecl).Specs {
				quoted := string(spec.(*ast.ImportSpec).Path.Value)
				unquoted, err := strconv.Unquote(quoted)
				if err != nil {
					log.Panicf("%s: parser returned invalid quoted string: <%s>", filename, quoted)
				}
				importsm[unquoted] = true
				if unquoted == "C" {
					cgoFiles = append(cgoFiles, d.Name)
				}
			}
		}
	}
	imports := make([]string, len(importsm))
	i := 0
	for p := range importsm {
		imports[i] = p
		i++
	}
	return &dirInfo{goFiles, cgoFiles, cFiles, imports, pkgName}, nil
}
