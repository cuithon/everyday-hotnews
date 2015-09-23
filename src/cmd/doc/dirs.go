// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"go/build"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Dirs is a structure for scanning the directory tree.
// Its Next method returns the next Go source directory it finds.
// Although it can be used to scan the tree multiple times, it
// only walks the tree once, caching the data it finds.
type Dirs struct {
	scan   chan string // directories generated by walk.
	paths  []string    // Cache of known paths.
	offset int         // Counter for Next.
}

var dirs Dirs

func init() {
	dirs.paths = make([]string, 0, 1000)
	dirs.scan = make(chan string)
	go dirs.walk()
}

// Reset puts the scan back at the beginning.
func (d *Dirs) Reset() {
	d.offset = 0
}

// Next returns the next directory in the scan. The boolean
// is false when the scan is done.
func (d *Dirs) Next() (string, bool) {
	if d.offset < len(d.paths) {
		path := d.paths[d.offset]
		d.offset++
		return path, true
	}
	path, ok := <-d.scan
	if !ok {
		return "", false
	}
	d.paths = append(d.paths, path)
	d.offset++
	return path, ok
}

// walk walks the trees in GOROOT and GOPATH.
func (d *Dirs) walk() {
	d.walkRoot(build.Default.GOROOT)
	for _, root := range splitGopath() {
		d.walkRoot(root)
	}
	close(d.scan)
}

// walkRoot walks a single directory. Each Go source directory it finds is
// delivered on d.scan.
func (d *Dirs) walkRoot(root string) {
	root = path.Join(root, "src")
	slashDot := string(filepath.Separator) + "."
	// We put a slash on the pkg so can use simple string comparison below
	// yet avoid inadvertent matches, like /foobar matching bar.

	visit := func(pathName string, f os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		// One package per directory. Ignore the files themselves.
		if !f.IsDir() {
			return nil
		}
		// No .git or other dot nonsense please.
		if strings.Contains(pathName, slashDot) {
			return filepath.SkipDir
		}
		// Does the directory contain any Go files? If so, it's a candidate.
		if hasGoFiles(pathName) {
			d.scan <- pathName
			return nil
		}
		return nil
	}

	filepath.Walk(root, visit)
}

// hasGoFiles tests whether the directory contains at least one file with ".go"
// extension
func hasGoFiles(path string) bool {
	dir, err := os.Open(path)
	if err != nil {
		// ignore unreadable directories
		return false
	}
	defer dir.Close()

	names, err := dir.Readdirnames(0)
	if err != nil {
		// ignore unreadable directories
		return false
	}

	for _, name := range names {
		if strings.HasSuffix(name, ".go") {
			return true
		}
	}

	return false
}
