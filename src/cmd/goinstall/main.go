// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Experimental Go package installer; see doc.go.

package main

import (
	"bytes"
	"exec"
	"flag"
	"fmt"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func usage() {
	fmt.Fprint(os.Stderr, "usage: goinstall importpath...\n")
	fmt.Fprintf(os.Stderr, "\tgoinstall -a\n")
	flag.PrintDefaults()
	os.Exit(2)
}

var (
	fset          = token.NewFileSet()
	argv0         = os.Args[0]
	errors        = false
	parents       = make(map[string]string)
	visit         = make(map[string]status)
	logfile       = filepath.Join(runtime.GOROOT(), "goinstall.log")
	installedPkgs = make(map[string]bool)

	allpkg            = flag.Bool("a", false, "install all previously installed packages")
	reportToDashboard = flag.Bool("dashboard", true, "report public packages at "+dashboardURL)
	logPkgs           = flag.Bool("log", true, "log installed packages to $GOROOT/goinstall.log for use by -a")
	update            = flag.Bool("u", false, "update already-downloaded packages")
	clean             = flag.Bool("clean", false, "clean the package directory before installing")
	verbose           = flag.Bool("v", false, "verbose")
)

type status int // status for visited map
const (
	unvisited status = iota
	visiting
	done
)

func logf(format string, args ...interface{}) {
	format = "%s: " + format
	args = append([]interface{}{argv0}, args...)
	fmt.Fprintf(os.Stderr, format, args...)
}

func vlogf(format string, args ...interface{}) {
	if *verbose {
		logf(format, args...)
	}
}

func errorf(format string, args ...interface{}) {
	errors = true
	logf(format, args...)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if runtime.GOROOT() == "" {
		fmt.Fprintf(os.Stderr, "%s: no $GOROOT\n", argv0)
		os.Exit(1)
	}

	// special case - "unsafe" is already installed
	visit["unsafe"] = done

	args := flag.Args()
	if *allpkg || *logPkgs {
		readPackageList()
	}
	if *allpkg {
		if len(args) != 0 {
			usage() // -a and package list both provided
		}
		// install all packages that were ever installed
		if len(installedPkgs) == 0 {
			fmt.Fprintf(os.Stderr, "%s: no installed packages\n", argv0)
			os.Exit(1)
		}
		args = make([]string, len(installedPkgs), len(installedPkgs))
		i := 0
		for pkg := range installedPkgs {
			args[i] = pkg
			i++
		}
	}
	if len(args) == 0 {
		usage()
	}
	for _, path := range args {
		if strings.HasPrefix(path, "http://") {
			errorf("'http://' used in remote path, try '%s'\n", path[7:])
			continue
		}

		install(path, "")
	}
	if errors {
		os.Exit(1)
	}
}

// printDeps prints the dependency path that leads to pkg.
func printDeps(pkg string) {
	if pkg == "" {
		return
	}
	if visit[pkg] != done {
		printDeps(parents[pkg])
	}
	fmt.Fprintf(os.Stderr, "\t%s ->\n", pkg)
}

// readPackageList reads the list of installed packages from goinstall.log
func readPackageList() {
	pkglistdata, _ := ioutil.ReadFile(logfile)
	pkglist := strings.Fields(string(pkglistdata))
	for _, pkg := range pkglist {
		installedPkgs[pkg] = true
	}
}

// logPackage logs the named package as installed in goinstall.log, if the package is not found in there
func logPackage(pkg string) {
	if installedPkgs[pkg] {
		return
	}
	fout, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", argv0, err)
		return
	}
	fmt.Fprintf(fout, "%s\n", pkg)
	fout.Close()
}

// install installs the package named by path, which is needed by parent.
func install(pkg, parent string) {
	// Make sure we're not already trying to install pkg.
	switch visit[pkg] {
	case done:
		return
	case visiting:
		fmt.Fprintf(os.Stderr, "%s: package dependency cycle\n", argv0)
		printDeps(parent)
		fmt.Fprintf(os.Stderr, "\t%s\n", pkg)
		os.Exit(2)
	}
	visit[pkg] = visiting
	parents[pkg] = parent

	vlogf("%s: visit\n", pkg)

	// Check whether package is local or remote.
	// If remote, download or update it.
	proot, pkg, err := findPackageRoot(pkg)
	// Don't build the standard library.
	if err == nil && proot.goroot && isStandardPath(pkg) {
		if parent == "" {
			errorf("%s: can not goinstall the standard library\n", pkg)
		} else {
			vlogf("%s: skipping standard library\n", pkg)
		}
		visit[pkg] = done
		return
	}
	// Download remote packages if not found or forced with -u flag.
	remote := isRemote(pkg)
	if remote && (err == ErrPackageNotFound || (err == nil && *update)) {
		vlogf("%s: download\n", pkg)
		err = download(pkg, proot.srcDir())
	}
	if err != nil {
		errorf("%s: %v\n", pkg, err)
		visit[pkg] = done
		return
	}
	dir := filepath.Join(proot.srcDir(), pkg)

	// Install prerequisites.
	dirInfo, err := scanDir(dir, parent == "")
	if err != nil {
		errorf("%s: %v\n", pkg, err)
		visit[pkg] = done
		return
	}
	if len(dirInfo.goFiles) == 0 {
		errorf("%s: package has no files\n", pkg)
		visit[pkg] = done
		return
	}
	for _, p := range dirInfo.imports {
		if p != "C" {
			install(p, pkg)
		}
	}

	// Install this package.
	if !errors {
		isCmd := dirInfo.pkgName == "main"
		if err := domake(dir, pkg, proot, isCmd); err != nil {
			errorf("installing: %v\n", err)
		} else if remote && *logPkgs {
			// mark package as installed in $GOROOT/goinstall.log
			logPackage(pkg)
		}
	}
	visit[pkg] = done
}


// Is this a standard package path?  strings container/vector etc.
// Assume that if the first element has a dot, it's a domain name
// and is not the standard package path.
func isStandardPath(s string) bool {
	dot := strings.Index(s, ".")
	slash := strings.Index(s, "/")
	return dot < 0 || 0 < slash && slash < dot
}

// run runs the command cmd in directory dir with standard input stdin.
// If the command fails, run prints the command and output on standard error
// in addition to returning a non-nil os.Error.
func run(dir string, stdin []byte, cmd ...string) os.Error {
	return genRun(dir, stdin, cmd, false)
}

// quietRun is like run but prints nothing on failure unless -v is used.
func quietRun(dir string, stdin []byte, cmd ...string) os.Error {
	return genRun(dir, stdin, cmd, true)
}

// genRun implements run and quietRun.
func genRun(dir string, stdin []byte, cmd []string, quiet bool) os.Error {
	bin, err := exec.LookPath(cmd[0])
	if err != nil {
		return err
	}
	p, err := exec.Run(bin, cmd, os.Environ(), dir, exec.Pipe, exec.Pipe, exec.MergeWithStdout)
	vlogf("%s: %s %s\n", dir, bin, strings.Join(cmd[1:], " "))
	if err != nil {
		return err
	}
	go func() {
		p.Stdin.Write(stdin)
		p.Stdin.Close()
	}()
	var buf bytes.Buffer
	io.Copy(&buf, p.Stdout)
	w, err := p.Wait(0)
	p.Close()
	if err != nil {
		return err
	}
	if !w.Exited() || w.ExitStatus() != 0 {
		if !quiet || *verbose {
			if dir != "" {
				dir = "cd " + dir + "; "
			}
			fmt.Fprintf(os.Stderr, "%s: === %s%s\n", argv0, dir, strings.Join(cmd, " "))
			os.Stderr.Write(buf.Bytes())
			fmt.Fprintf(os.Stderr, "--- %s\n", w)
		}
		return os.ErrorString("running " + cmd[0] + ": " + w.String())
	}
	return nil
}
