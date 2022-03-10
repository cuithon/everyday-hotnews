// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	exec "internal/execabs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"golang.org/x/mod/semver"
)

// A Dir describes a directory holding code by specifying
// the expected import path and the file system directory.
type Dir struct {
	importPath string // import path for that dir
	dir        string // file system directory
	inModule   bool
}

// Dirs is a structure for scanning the directory tree.
// Its Next method returns the next Go source directory it finds.
// Although it can be used to scan the tree multiple times, it
// only walks the tree once, caching the data it finds.
type Dirs struct {
	scan   chan Dir // Directories generated by walk.
	hist   []Dir    // History of reported Dirs.
	offset int      // Counter for Next.
}

var dirs Dirs

// dirsInit starts the scanning of package directories in GOROOT and GOPATH. Any
// extra paths passed to it are included in the channel.
func dirsInit(extra ...Dir) {
	if buildCtx.GOROOT == "" {
		stdout, err := exec.Command("go", "env", "GOROOT").Output()
		if err != nil {
			if ee, ok := err.(*exec.ExitError); ok && len(ee.Stderr) > 0 {
				log.Fatalf("failed to determine GOROOT: $GOROOT is not set and 'go env GOROOT' failed:\n%s", ee.Stderr)
			}
			log.Fatalf("failed to determine GOROOT: $GOROOT is not set and could not run 'go env GOROOT':\n\t%s", err)
		}
		buildCtx.GOROOT = string(bytes.TrimSpace(stdout))
	}

	dirs.hist = make([]Dir, 0, 1000)
	dirs.hist = append(dirs.hist, extra...)
	dirs.scan = make(chan Dir)
	go dirs.walk(codeRoots())
}

// Reset puts the scan back at the beginning.
func (d *Dirs) Reset() {
	d.offset = 0
}

// Next returns the next directory in the scan. The boolean
// is false when the scan is done.
func (d *Dirs) Next() (Dir, bool) {
	if d.offset < len(d.hist) {
		dir := d.hist[d.offset]
		d.offset++
		return dir, true
	}
	dir, ok := <-d.scan
	if !ok {
		return Dir{}, false
	}
	d.hist = append(d.hist, dir)
	d.offset++
	return dir, ok
}

// walk walks the trees in GOROOT and GOPATH.
func (d *Dirs) walk(roots []Dir) {
	for _, root := range roots {
		d.bfsWalkRoot(root)
	}
	close(d.scan)
}

// bfsWalkRoot walks a single directory hierarchy in breadth-first lexical order.
// Each Go source directory it finds is delivered on d.scan.
func (d *Dirs) bfsWalkRoot(root Dir) {
	root.dir = filepath.Clean(root.dir) // because filepath.Join will do it anyway

	// this is the queue of directories to examine in this pass.
	this := []string{}
	// next is the queue of directories to examine in the next pass.
	next := []string{root.dir}

	for len(next) > 0 {
		this, next = next, this[0:0]
		for _, dir := range this {
			fd, err := os.Open(dir)
			if err != nil {
				log.Print(err)
				continue
			}
			entries, err := fd.Readdir(0)
			fd.Close()
			if err != nil {
				log.Print(err)
				continue
			}
			hasGoFiles := false
			for _, entry := range entries {
				name := entry.Name()
				// For plain files, remember if this directory contains any .go
				// source files, but ignore them otherwise.
				if !entry.IsDir() {
					if !hasGoFiles && strings.HasSuffix(name, ".go") {
						hasGoFiles = true
					}
					continue
				}
				// Entry is a directory.

				// The go tool ignores directories starting with ., _, or named "testdata".
				if name[0] == '.' || name[0] == '_' || name == "testdata" {
					continue
				}
				// When in a module, ignore vendor directories and stop at module boundaries.
				if root.inModule {
					if name == "vendor" {
						continue
					}
					if fi, err := os.Stat(filepath.Join(dir, name, "go.mod")); err == nil && !fi.IsDir() {
						continue
					}
				}
				// Remember this (fully qualified) directory for the next pass.
				next = append(next, filepath.Join(dir, name))
			}
			if hasGoFiles {
				// It's a candidate.
				importPath := root.importPath
				if len(dir) > len(root.dir) {
					if importPath != "" {
						importPath += "/"
					}
					importPath += filepath.ToSlash(dir[len(root.dir)+1:])
				}
				d.scan <- Dir{importPath, dir, root.inModule}
			}
		}

	}
}

var testGOPATH = false // force GOPATH use for testing

// codeRoots returns the code roots to search for packages.
// In GOPATH mode this is GOROOT/src and GOPATH/src, with empty import paths.
// In module mode, this is each module root, with an import path set to its module path.
func codeRoots() []Dir {
	codeRootsCache.once.Do(func() {
		codeRootsCache.roots = findCodeRoots()
	})
	return codeRootsCache.roots
}

var codeRootsCache struct {
	once  sync.Once
	roots []Dir
}

var usingModules bool

func findCodeRoots() []Dir {
	var list []Dir
	if !testGOPATH {
		// Check for use of modules by 'go env GOMOD',
		// which reports a go.mod file path if modules are enabled.
		stdout, _ := exec.Command("go", "env", "GOMOD").Output()
		gomod := string(bytes.TrimSpace(stdout))

		usingModules = len(gomod) > 0
		if usingModules && buildCtx.GOROOT != "" {
			list = append(list,
				Dir{dir: filepath.Join(buildCtx.GOROOT, "src"), inModule: true},
				Dir{importPath: "cmd", dir: filepath.Join(buildCtx.GOROOT, "src", "cmd"), inModule: true})
		}

		if gomod == os.DevNull {
			// Modules are enabled, but the working directory is outside any module.
			// We can still access std, cmd, and packages specified as source files
			// on the command line, but there are no module roots.
			// Avoid 'go list -m all' below, since it will not work.
			return list
		}
	}

	if !usingModules {
		if buildCtx.GOROOT != "" {
			list = append(list, Dir{dir: filepath.Join(buildCtx.GOROOT, "src")})
		}
		for _, root := range splitGopath() {
			list = append(list, Dir{dir: filepath.Join(root, "src")})
		}
		return list
	}

	// Find module root directories from go list.
	// Eventually we want golang.org/x/tools/go/packages
	// to handle the entire file system search and become go/packages,
	// but for now enumerating the module roots lets us fit modules
	// into the current code with as few changes as possible.
	mainMod, vendorEnabled, err := vendorEnabled()
	if err != nil {
		return list
	}
	if vendorEnabled {
		// Add the vendor directory to the search path ahead of "std".
		// That way, if the main module *is* "std", we will identify the path
		// without the "vendor/" prefix before the one with that prefix.
		list = append([]Dir{{dir: filepath.Join(mainMod.Dir, "vendor"), inModule: false}}, list...)
		if mainMod.Path != "std" {
			list = append(list, Dir{importPath: mainMod.Path, dir: mainMod.Dir, inModule: true})
		}
		return list
	}

	cmd := exec.Command("go", "list", "-m", "-f={{.Path}}\t{{.Dir}}", "all")
	cmd.Stderr = os.Stderr
	out, _ := cmd.Output()
	for _, line := range strings.Split(string(out), "\n") {
		path, dir, _ := strings.Cut(line, "\t")
		if dir != "" {
			list = append(list, Dir{importPath: path, dir: dir, inModule: true})
		}
	}

	return list
}

// The functions below are derived from x/tools/internal/imports at CL 203017.

type moduleJSON struct {
	Path, Dir, GoVersion string
}

var modFlagRegexp = regexp.MustCompile(`-mod[ =](\w+)`)

// vendorEnabled indicates if vendoring is enabled.
// Inspired by setDefaultBuildMod in modload/init.go
func vendorEnabled() (*moduleJSON, bool, error) {
	mainMod, go114, err := getMainModuleAnd114()
	if err != nil {
		return nil, false, err
	}

	stdout, _ := exec.Command("go", "env", "GOFLAGS").Output()
	goflags := string(bytes.TrimSpace(stdout))
	matches := modFlagRegexp.FindStringSubmatch(goflags)
	var modFlag string
	if len(matches) != 0 {
		modFlag = matches[1]
	}
	if modFlag != "" {
		// Don't override an explicit '-mod=' argument.
		return mainMod, modFlag == "vendor", nil
	}
	if mainMod == nil || !go114 {
		return mainMod, false, nil
	}
	// Check 1.14's automatic vendor mode.
	if fi, err := os.Stat(filepath.Join(mainMod.Dir, "vendor")); err == nil && fi.IsDir() {
		if mainMod.GoVersion != "" && semver.Compare("v"+mainMod.GoVersion, "v1.14") >= 0 {
			// The Go version is at least 1.14, and a vendor directory exists.
			// Set -mod=vendor by default.
			return mainMod, true, nil
		}
	}
	return mainMod, false, nil
}

// getMainModuleAnd114 gets the main module's information and whether the
// go command in use is 1.14+. This is the information needed to figure out
// if vendoring should be enabled.
func getMainModuleAnd114() (*moduleJSON, bool, error) {
	const format = `{{.Path}}
{{.Dir}}
{{.GoVersion}}
{{range context.ReleaseTags}}{{if eq . "go1.14"}}{{.}}{{end}}{{end}}
`
	cmd := exec.Command("go", "list", "-m", "-f", format)
	cmd.Stderr = os.Stderr
	stdout, err := cmd.Output()
	if err != nil {
		return nil, false, nil
	}
	lines := strings.Split(string(stdout), "\n")
	if len(lines) < 5 {
		return nil, false, fmt.Errorf("unexpected stdout: %q", stdout)
	}
	mod := &moduleJSON{
		Path:      lines[0],
		Dir:       lines[1],
		GoVersion: lines[2],
	}
	return mod, lines[3] == "go1.14", nil
}
