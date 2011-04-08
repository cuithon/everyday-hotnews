// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"exec"
	"flag"
	"fmt"
	"go/parser"
	"go/printer"
	"go/scanner"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	fset     = token.NewFileSet()
	exitCode = 0
)

var allowedRewrites = flag.String("r", "",
	"restrict the rewrites to this comma-separated list")

var allowed map[string]bool

var doDiff = flag.Bool("diff", false, "display diffs instead of rewriting files")

func usage() {
	fmt.Fprintf(os.Stderr, "usage: gofix [-diff] [-r fixname,...] [path ...]\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nAvailable rewrites are:\n")
	for _, f := range fixes {
		fmt.Fprintf(os.Stderr, "\n%s\n", f.name)
		desc := strings.TrimSpace(f.desc)
		desc = strings.Replace(desc, "\n", "\n\t", -1)
		fmt.Fprintf(os.Stderr, "\t%s\n", desc)
	}
	os.Exit(2)
}

func main() {
	sort.Sort(fixes)

	flag.Usage = usage
	flag.Parse()

	if *allowedRewrites != "" {
		allowed = make(map[string]bool)
		for _, f := range strings.Split(*allowedRewrites, ",", -1) {
			allowed[f] = true
		}
	}

	if flag.NArg() == 0 {
		if err := processFile("standard input", true); err != nil {
			report(err)
		}
		os.Exit(exitCode)
	}

	for i := 0; i < flag.NArg(); i++ {
		path := flag.Arg(i)
		switch dir, err := os.Stat(path); {
		case err != nil:
			report(err)
		case dir.IsRegular():
			if err := processFile(path, false); err != nil {
				report(err)
			}
		case dir.IsDirectory():
			walkDir(path)
		}
	}

	os.Exit(exitCode)
}

const (
	tabWidth    = 8
	parserMode  = parser.ParseComments
	printerMode = printer.TabIndent | printer.UseSpaces
)

var printConfig = &printer.Config{
	printerMode,
	tabWidth,
}

func processFile(filename string, useStdin bool) os.Error {
	var f *os.File
	var err os.Error
	var fixlog bytes.Buffer
	var buf bytes.Buffer

	if useStdin {
		f = os.Stdin
	} else {
		f, err = os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
	}

	src, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	file, err := parser.ParseFile(fset, filename, src, parserMode)
	if err != nil {
		return err
	}

	// Apply all fixes to file.
	newFile := file
	fixed := false
	for _, fix := range fixes {
		if allowed != nil && !allowed[fix.desc] {
			continue
		}
		if fix.f(newFile) {
			fixed = true
			fmt.Fprintf(&fixlog, " %s", fix.name)

			// AST changed.
			// Print and parse, to update any missing scoping
			// or position information for subsequent fixers.
			buf.Reset()
			_, err = printConfig.Fprint(&buf, fset, newFile)
			if err != nil {
				return err
			}
			newSrc := buf.Bytes()
			newFile, err = parser.ParseFile(fset, filename, newSrc, parserMode)
			if err != nil {
				return err
			}
		}
	}
	if !fixed {
		return nil
	}
	fmt.Fprintf(os.Stderr, "%s: fixed %s\n", filename, fixlog.String()[1:])

	// Print AST.  We did that after each fix, so this appears
	// redundant, but it is necessary to generate gofmt-compatible
	// source code in a few cases.  The official gofmt style is the
	// output of the printer run on a standard AST generated by the parser,
	// but the source we generated inside the loop above is the
	// output of the printer run on a mangled AST generated by a fixer.
	buf.Reset()
	_, err = printConfig.Fprint(&buf, fset, newFile)
	if err != nil {
		return err
	}
	newSrc := buf.Bytes()

	if *doDiff {
		data, err := diff(src, newSrc)
		if err != nil {
			return fmt.Errorf("computing diff: %s", err)
		}
		fmt.Printf("diff %s fixed/%s\n", filename, filename)
		os.Stdout.Write(data)
		return nil
	}

	if useStdin {
		os.Stdout.Write(newSrc)
		return nil
	}

	return ioutil.WriteFile(f.Name(), newSrc, 0)
}

var gofmtBuf bytes.Buffer

func gofmt(n interface{}) string {
	gofmtBuf.Reset()
	_, err := printConfig.Fprint(&gofmtBuf, fset, n)
	if err != nil {
		return "<" + err.String() + ">"
	}
	return gofmtBuf.String()
}

func report(err os.Error) {
	scanner.PrintError(os.Stderr, err)
	exitCode = 2
}

func walkDir(path string) {
	v := make(fileVisitor)
	go func() {
		filepath.Walk(path, v, v)
		close(v)
	}()
	for err := range v {
		if err != nil {
			report(err)
		}
	}
}

type fileVisitor chan os.Error

func (v fileVisitor) VisitDir(path string, f *os.FileInfo) bool {
	return true
}

func (v fileVisitor) VisitFile(path string, f *os.FileInfo) {
	if isGoFile(f) {
		v <- nil // synchronize error handler
		if err := processFile(path, false); err != nil {
			v <- err
		}
	}
}

func isGoFile(f *os.FileInfo) bool {
	// ignore non-Go files
	return f.IsRegular() && !strings.HasPrefix(f.Name, ".") && strings.HasSuffix(f.Name, ".go")
}

func diff(b1, b2 []byte) (data []byte, err os.Error) {
	f1, err := ioutil.TempFile("", "gofix")
	if err != nil {
		return nil, err
	}
	defer os.Remove(f1.Name())
	defer f1.Close()

	f2, err := ioutil.TempFile("", "gofix")
	if err != nil {
		return nil, err
	}
	defer os.Remove(f2.Name())
	defer f2.Close()

	f1.Write(b1)
	f2.Write(b2)

	diffcmd, err := exec.LookPath("diff")
	if err != nil {
		return nil, err
	}

	c, err := exec.Run(diffcmd, []string{"diff", f1.Name(), f2.Name()}, nil, "",
		exec.DevNull, exec.Pipe, exec.MergeWithStdout)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	return ioutil.ReadAll(c.Stdout)
}
