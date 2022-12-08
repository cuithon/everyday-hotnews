// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements a typechecker test harness. The packages specified
// in tests are typechecked. Error messages reported by the typechecker are
// compared against the error messages expected in the test files.
//
// Expected errors are indicated in the test files by putting a comment
// of the form /* ERROR "rx" */ immediately following an offending token.
// The harness will verify that an error matching the regular expression
// rx is reported at that source position. Consecutive comments may be
// used to indicate multiple errors for the same token position.
//
// For instance, the following test file indicates that a "not declared"
// error should be reported for the undeclared variable x:
//
//	package p
//	func f() {
//		_ = x /* ERROR "not declared" */ + 1
//	}

package types_test

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/scanner"
	"go/token"
	"internal/testenv"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"

	. "go/types"
)

var (
	haltOnError  = flag.Bool("halt", false, "halt on error")
	verifyErrors = flag.Bool("verify", false, "verify errors (rather than list them) in TestManual")
)

var fset = token.NewFileSet()

func parseFiles(t *testing.T, filenames []string, srcs [][]byte, mode parser.Mode) ([]*ast.File, []error) {
	var files []*ast.File
	var errlist []error
	for i, filename := range filenames {
		file, err := parser.ParseFile(fset, filename, srcs[i], mode)
		if file == nil {
			t.Fatalf("%s: %s", filename, err)
		}
		files = append(files, file)
		if err != nil {
			if list, _ := err.(scanner.ErrorList); len(list) > 0 {
				for _, err := range list {
					errlist = append(errlist, err)
				}
			} else {
				errlist = append(errlist, err)
			}
		}
	}
	return files, errlist
}

func unpackError(fset *token.FileSet, err error) (token.Position, string) {
	switch err := err.(type) {
	case *scanner.Error:
		return err.Pos, err.Msg
	case Error:
		return fset.Position(err.Pos), err.Msg
	}
	panic("unreachable")
}

// delta returns the absolute difference between x and y.
func delta(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}

// parseFlags parses flags from the first line of the given source
// (from src if present, or by reading from the file) if the line
// starts with "//" (line comment) followed by "-" (possibly with
// spaces between). Otherwise the line is ignored.
func parseFlags(filename string, src []byte, flags *flag.FlagSet) error {
	// If there is no src, read from the file.
	const maxLen = 256
	if len(src) == 0 {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}

		var buf [maxLen]byte
		n, err := f.Read(buf[:])
		if err != nil {
			return err
		}
		src = buf[:n]
	}

	// we must have a line comment that starts with a "-"
	const prefix = "//"
	if !bytes.HasPrefix(src, []byte(prefix)) {
		return nil // first line is not a line comment
	}
	src = src[len(prefix):]
	if i := bytes.Index(src, []byte("-")); i < 0 || len(bytes.TrimSpace(src[:i])) != 0 {
		return nil // comment doesn't start with a "-"
	}
	end := bytes.Index(src, []byte("\n"))
	if end < 0 || end > maxLen {
		return fmt.Errorf("flags comment line too long")
	}

	return flags.Parse(strings.Fields(string(src[:end])))
}

func testFiles(t *testing.T, sizes Sizes, filenames []string, srcs [][]byte, manual bool, imp Importer) {
	if len(filenames) == 0 {
		t.Fatal("no source files")
	}

	var conf Config
	conf.Sizes = sizes
	flags := flag.NewFlagSet("", flag.PanicOnError)
	flags.StringVar(&conf.GoVersion, "lang", "", "")
	flags.BoolVar(&conf.FakeImportC, "fakeImportC", false, "")
	flags.BoolVar(addrOldComparableSemantics(&conf), "oldComparableSemantics", false, "")
	if err := parseFlags(filenames[0], srcs[0], flags); err != nil {
		t.Fatal(err)
	}

	files, errlist := parseFiles(t, filenames, srcs, parser.AllErrors)

	pkgName := "<no package>"
	if len(files) > 0 {
		pkgName = files[0].Name.Name
	}

	listErrors := manual && !*verifyErrors
	if listErrors && len(errlist) > 0 {
		t.Errorf("--- %s:", pkgName)
		for _, err := range errlist {
			t.Error(err)
		}
	}

	// typecheck and collect typechecker errors
	if imp == nil {
		imp = importer.Default()
	}
	conf.Importer = imp
	conf.Error = func(err error) {
		if *haltOnError {
			defer panic(err)
		}
		if listErrors {
			t.Error(err)
			return
		}
		// Ignore secondary error messages starting with "\t";
		// they are clarifying messages for a primary error.
		if !strings.Contains(err.Error(), ": \t") {
			errlist = append(errlist, err)
		}
	}
	conf.Check(pkgName, fset, files, nil)

	if listErrors {
		return
	}

	// sort errlist in source order
	sort.Slice(errlist, func(i, j int) bool {
		// TODO(gri) This is not correct as scanner.Errors
		// don't have a correctly set Offset. But we only
		// care about sorting when multiple equal errors
		// appear on the same line, which happens with some
		// type checker errors.
		// For now this works. Will remove need for sorting
		// in a subsequent CL.
		pi, _ := unpackError(fset, errlist[i])
		pj, _ := unpackError(fset, errlist[j])
		return pi.Offset < pj.Offset
	})

	// collect expected errors
	errmap := make(map[string]map[int][]comment)
	for i, filename := range filenames {
		if m := commentMap(srcs[i], regexp.MustCompile("^ ERROR ")); len(m) > 0 {
			errmap[filename] = m
		}
	}

	// match against found errors
	for _, err := range errlist {
		gotPos, gotMsg := unpackError(fset, err)

		// find list of errors for the respective error line
		filename := gotPos.Filename
		filemap := errmap[filename]
		line := gotPos.Line
		var errList []comment
		if filemap != nil {
			errList = filemap[line]
		}

		// one of errors in errList should match the current error
		index := -1 // errList index of matching message, if any
		for i, want := range errList {
			pattern := strings.TrimSpace(want.text[len(" ERROR "):])
			if n := len(pattern); n >= 2 && pattern[0] == '"' && pattern[n-1] == '"' {
				pattern = pattern[1 : n-1]
			}
			rx, err := regexp.Compile(pattern)
			if err != nil {
				t.Errorf("%s:%d:%d: %v", filename, line, want.col, err)
				continue
			}
			if rx.MatchString(gotMsg) {
				index = i
				break
			}
		}
		if index < 0 {
			t.Errorf("%s: no error expected: %q", gotPos, gotMsg)
			continue
		}

		// column position must be within expected colDelta
		const colDelta = 0
		want := errList[index]
		if delta(gotPos.Column, want.col) > colDelta {
			t.Errorf("%s: got col = %d; want %d", gotPos, gotPos.Column, want.col)
		}

		// eliminate from errList
		if n := len(errList) - 1; n > 0 {
			// not the last entry - slide entries down (don't reorder)
			copy(errList[index:], errList[index+1:])
			filemap[line] = errList[:n]
		} else {
			// last entry - remove errList from filemap
			delete(filemap, line)
		}

		// if filemap is empty, eliminate from errmap
		if len(filemap) == 0 {
			delete(errmap, filename)
		}
	}

	// there should be no expected errors left
	if len(errmap) > 0 {
		t.Errorf("--- %s: unreported errors:", pkgName)
		for filename, filemap := range errmap {
			for line, errList := range filemap {
				for _, err := range errList {
					t.Errorf("%s:%d:%d: %s", filename, line, err.col, err.text)
				}
			}
		}
	}
}

func readCode(err Error) int {
	v := reflect.ValueOf(err)
	return int(v.FieldByName("go116code").Int())
}

// addrOldComparableSemantics(conf) returns &conf.oldComparableSemantics (unexported field).
func addrOldComparableSemantics(conf *Config) *bool {
	v := reflect.Indirect(reflect.ValueOf(conf))
	return (*bool)(v.FieldByName("oldComparableSemantics").Addr().UnsafePointer())
}

// TestManual is for manual testing of a package - either provided
// as a list of filenames belonging to the package, or a directory
// name containing the package files - after the test arguments
// (and a separating "--"). For instance, to test the package made
// of the files foo.go and bar.go, use:
//
//	go test -run Manual -- foo.go bar.go
//
// If no source arguments are provided, the file testdata/manual.go
// is used instead.
// Provide the -verify flag to verify errors against ERROR comments
// in the input files rather than having a list of errors reported.
// The accepted Go language version can be controlled with the -lang
// flag.
func TestManual(t *testing.T) {
	testenv.MustHaveGoBuild(t)

	filenames := flag.Args()
	if len(filenames) == 0 {
		filenames = []string{filepath.FromSlash("testdata/manual.go")}
	}

	info, err := os.Stat(filenames[0])
	if err != nil {
		t.Fatalf("TestManual: %v", err)
	}

	DefPredeclaredTestFuncs()
	if info.IsDir() {
		if len(filenames) > 1 {
			t.Fatal("TestManual: must have only one directory argument")
		}
		testDir(t, filenames[0], true)
	} else {
		testPkg(t, filenames, true)
	}
}

func TestLongConstants(t *testing.T) {
	format := "package longconst\n\nconst _ = %s /* ERROR constant overflow */ \nconst _ = %s // ERROR excessively long constant"
	src := fmt.Sprintf(format, strings.Repeat("1", 9999), strings.Repeat("1", 10001))
	testFiles(t, nil, []string{"longconst.go"}, [][]byte{[]byte(src)}, false, nil)
}

// TestIndexRepresentability tests that constant index operands must
// be representable as int even if they already have a type that can
// represent larger values.
func TestIndexRepresentability(t *testing.T) {
	const src = "package index\n\nvar s []byte\nvar _ = s[int64 /* ERROR \"int64\\(1\\) << 40 \\(.*\\) overflows int\" */ (1) << 40]"
	testFiles(t, &StdSizes{4, 4}, []string{"index.go"}, [][]byte{[]byte(src)}, false, nil)
}

func TestIssue47243_TypedRHS(t *testing.T) {
	// The RHS of the shift expression below overflows uint on 32bit platforms,
	// but this is OK as it is explicitly typed.
	const src = "package issue47243\n\nvar a uint64; var _ = a << uint64(4294967296)" // uint64(1<<32)
	testFiles(t, &StdSizes{4, 4}, []string{"p.go"}, [][]byte{[]byte(src)}, false, nil)
}

func TestCheck(t *testing.T) {
	DefPredeclaredTestFuncs()
	testDirFiles(t, "../../internal/types/testdata/check", false)
}
func TestSpec(t *testing.T)      { testDirFiles(t, "../../internal/types/testdata/spec", false) }
func TestExamples(t *testing.T)  { testDirFiles(t, "../../internal/types/testdata/examples", false) }
func TestFixedbugs(t *testing.T) { testDirFiles(t, "../../internal/types/testdata/fixedbugs", false) }
func TestLocal(t *testing.T)     { testDirFiles(t, "testdata/local", false) }

func testDirFiles(t *testing.T, dir string, manual bool) {
	testenv.MustHaveGoBuild(t)
	dir = filepath.FromSlash(dir)

	fis, err := os.ReadDir(dir)
	if err != nil {
		t.Error(err)
		return
	}

	for _, fi := range fis {
		path := filepath.Join(dir, fi.Name())

		// If fi is a directory, its files make up a single package.
		if fi.IsDir() {
			testDir(t, path, manual)
		} else {
			t.Run(filepath.Base(path), func(t *testing.T) {
				testPkg(t, []string{path}, manual)
			})
		}
	}
}

func testDir(t *testing.T, dir string, manual bool) {
	testenv.MustHaveGoBuild(t)

	fis, err := os.ReadDir(dir)
	if err != nil {
		t.Error(err)
		return
	}

	var filenames []string
	for _, fi := range fis {
		filenames = append(filenames, filepath.Join(dir, fi.Name()))
	}

	t.Run(filepath.Base(dir), func(t *testing.T) {
		testPkg(t, filenames, manual)
	})
}

// TODO(rFindley) reconcile the different test setup in go/types with types2.
func testPkg(t *testing.T, filenames []string, manual bool) {
	srcs := make([][]byte, len(filenames))
	for i, filename := range filenames {
		src, err := os.ReadFile(filename)
		if err != nil {
			t.Fatalf("could not read %s: %v", filename, err)
		}
		srcs[i] = src
	}
	testFiles(t, nil, filenames, srcs, manual, nil)
}
