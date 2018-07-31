// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gc

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"internal/testenv"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TODO: move all these tests elsewhere?
// Perhaps teach test/run.go how to run them with a new action verb.
func runTest(t *testing.T, filename string, flags ...string) {
	t.Parallel()
	doTest(t, filename, "run", flags...)
}
func buildTest(t *testing.T, filename string, flags ...string) {
	t.Parallel()
	doTest(t, filename, "build", flags...)
}
func doTest(t *testing.T, filename string, kind string, flags ...string) {
	testenv.MustHaveGoBuild(t)
	gotool := testenv.GoToolPath(t)

	var stdout, stderr bytes.Buffer
	args := []string{kind}
	if len(flags) == 0 {
		args = append(args, "-gcflags=-d=ssa/check/on")
	} else {
		args = append(args, flags...)
	}
	args = append(args, filepath.Join("testdata", filename))
	cmd := exec.Command(gotool, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed: %v:\nOut: %s\nStderr: %s\n", err, &stdout, &stderr)
	}
	if s := stdout.String(); s != "" {
		t.Errorf("Stdout = %s\nWant empty", s)
	}
	if s := stderr.String(); strings.Contains(s, "SSA unimplemented") {
		t.Errorf("Unimplemented message found in stderr:\n%s", s)
	}
}

// runGenTest runs a test-generator, then runs the generated test.
// Generated test can either fail in compilation or execution.
// The environment variable parameter(s) is passed to the run
// of the generated test.
func runGenTest(t *testing.T, filename, tmpname string, ev ...string) {
	testenv.MustHaveGoRun(t)
	gotool := testenv.GoToolPath(t)
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(gotool, "run", filepath.Join("testdata", filename))
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed: %v:\nOut: %s\nStderr: %s\n", err, &stdout, &stderr)
	}
	// Write stdout into a temporary file
	tmpdir, ok := ioutil.TempDir("", tmpname)
	if ok != nil {
		t.Fatalf("Failed to create temporary directory")
	}
	defer os.RemoveAll(tmpdir)

	rungo := filepath.Join(tmpdir, "run.go")
	ok = ioutil.WriteFile(rungo, stdout.Bytes(), 0600)
	if ok != nil {
		t.Fatalf("Failed to create temporary file " + rungo)
	}

	stdout.Reset()
	stderr.Reset()
	cmd = exec.Command(gotool, "run", "-gcflags=-d=ssa/check/on", rungo)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Env = append(cmd.Env, ev...)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed: %v:\nOut: %s\nStderr: %s\n", err, &stdout, &stderr)
	}
	if s := stderr.String(); s != "" {
		t.Errorf("Stderr = %s\nWant empty", s)
	}
	if s := stdout.String(); s != "" {
		t.Errorf("Stdout = %s\nWant empty", s)
	}
}

func TestGenFlowGraph(t *testing.T) {
	if testing.Short() {
		t.Skip("not run in short mode.")
	}
	runGenTest(t, "flowgraph_generator1.go", "ssa_fg_tmp1")
}

// TestCode runs all the tests in the testdata directory as subtests.
// These tests are special because we want to run them with different
// compiler flags set (and thus they can't just be _test.go files in
// this directory).
func TestCode(t *testing.T) {
	testenv.MustHaveGoBuild(t)
	gotool := testenv.GoToolPath(t)

	// Make a temporary directory to work in.
	tmpdir, err := ioutil.TempDir("", "TestCode")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpdir)

	// Find all the test functions (and the files containing them).
	var srcs []string // files containing Test functions
	type test struct {
		name      string // TestFoo
		usesFloat bool   // might use float operations
	}
	var tests []test
	files, err := ioutil.ReadDir("testdata")
	if err != nil {
		t.Fatalf("can't read testdata directory: %v", err)
	}
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), "_test.go") {
			continue
		}
		text, err := ioutil.ReadFile(filepath.Join("testdata", f.Name()))
		if err != nil {
			t.Fatalf("can't read testdata/%s: %v", f.Name(), err)
		}
		fset := token.NewFileSet()
		code, err := parser.ParseFile(fset, f.Name(), text, 0)
		if err != nil {
			t.Fatalf("can't parse testdata/%s: %v", f.Name(), err)
		}
		srcs = append(srcs, filepath.Join("testdata", f.Name()))
		foundTest := false
		for _, d := range code.Decls {
			fd, ok := d.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if !strings.HasPrefix(fd.Name.Name, "Test") {
				continue
			}
			if fd.Recv != nil {
				continue
			}
			if fd.Type.Results != nil {
				continue
			}
			if len(fd.Type.Params.List) != 1 {
				continue
			}
			p := fd.Type.Params.List[0]
			if len(p.Names) != 1 {
				continue
			}
			s, ok := p.Type.(*ast.StarExpr)
			if !ok {
				continue
			}
			sel, ok := s.X.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			base, ok := sel.X.(*ast.Ident)
			if !ok {
				continue
			}
			if base.Name != "testing" {
				continue
			}
			if sel.Sel.Name != "T" {
				continue
			}
			// Found a testing function.
			tests = append(tests, test{name: fd.Name.Name, usesFloat: bytes.Contains(text, []byte("float"))})
			foundTest = true
		}
		if !foundTest {
			t.Fatalf("test file testdata/%s has no tests in it", f.Name())
		}
	}

	flags := []string{""}
	if runtime.GOARCH == "arm" || runtime.GOARCH == "mips" || runtime.GOARCH == "mips64" {
		flags = append(flags, ",softfloat")
	}
	for _, flag := range flags {
		args := []string{"test", "-c", "-gcflags=-d=ssa/check/on" + flag, "-o", filepath.Join(tmpdir, "code.test")}
		args = append(args, srcs...)
		out, err := exec.Command(gotool, args...).CombinedOutput()
		if err != nil || len(out) != 0 {
			t.Fatalf("Build failed: %v\n%s\n", err, out)
		}

		// Now we have a test binary. Run it with all the tests as subtests of this one.
		for _, test := range tests {
			test := test
			if flag == ",softfloat" && !test.usesFloat {
				// No point in running the soft float version if the test doesn't use floats.
				continue
			}
			t.Run(fmt.Sprintf("%s%s", test.name[4:], flag), func(t *testing.T) {
				out, err := exec.Command(filepath.Join(tmpdir, "code.test"), "-test.run="+test.name).CombinedOutput()
				if err != nil || string(out) != "PASS\n" {
					t.Errorf("Failed:\n%s\n", out)
				}
			})
		}
	}
}

// TestTypeAssertion tests type assertions.
func TestTypeAssertion(t *testing.T) { runTest(t, "assert.go") }

// TestArithmetic tests that both backends have the same result for arithmetic expressions.
func TestArithmetic(t *testing.T) { runTest(t, "arith.go") }

// TestFP tests that both backends have the same result for floating point expressions.
func TestFP(t *testing.T) { runTest(t, "fp.go") }

func TestFPSoftFloat(t *testing.T) {
	runTest(t, "fp.go", "-gcflags=-d=softfloat,ssa/check/on")
}

// TestArithmeticBoundary tests boundary results for arithmetic operations.
func TestArithmeticBoundary(t *testing.T) { runTest(t, "arithBoundary.go") }

// TestArithmeticConst tests results for arithmetic operations against constants.
func TestArithmeticConst(t *testing.T) { runTest(t, "arithConst.go") }

func TestChan(t *testing.T) { runTest(t, "chan.go") }

// TestComparisonsConst tests results for comparison operations against constants.
func TestComparisonsConst(t *testing.T) { runTest(t, "cmpConst.go") }

func TestCompound(t *testing.T) { runTest(t, "compound.go") }

func TestCtl(t *testing.T) { runTest(t, "ctl.go") }

func TestLoadStore(t *testing.T) { runTest(t, "loadstore.go") }

func TestMap(t *testing.T) { runTest(t, "map.go") }

func TestRegalloc(t *testing.T) { runTest(t, "regalloc.go") }

func TestString(t *testing.T) { runTest(t, "string.go") }

func TestDeferNoReturn(t *testing.T) { buildTest(t, "deferNoReturn.go") }

// TestClosure tests closure related behavior.
func TestClosure(t *testing.T) { runTest(t, "closure.go") }

func TestArray(t *testing.T) { runTest(t, "array.go") }

func TestAppend(t *testing.T) { runTest(t, "append.go") }

func TestZero(t *testing.T) { runTest(t, "zero.go") }

func TestAddressed(t *testing.T) { runTest(t, "addressed.go") }

func TestCopy(t *testing.T) { runTest(t, "copy.go") }

func TestUnsafe(t *testing.T) { runTest(t, "unsafe.go") }

func TestPhi(t *testing.T) { runTest(t, "phi.go") }

func TestSlice(t *testing.T) { runTest(t, "slice.go") }

func TestNamedReturn(t *testing.T) { runTest(t, "namedReturn.go") }

func TestDuplicateLoad(t *testing.T) { runTest(t, "dupLoad.go") }

func TestSqrt(t *testing.T) { runTest(t, "sqrt_const.go") }
