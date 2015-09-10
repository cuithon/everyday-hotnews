// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package asm

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cmd/asm/internal/lex"
	"cmd/internal/obj"
)

// An end-to-end test for the assembler: Do we print what we parse?
// Output is generated by, in effect, turning on -S and comparing the
// result against a golden file.

func testEndToEnd(t *testing.T, goarch string) {
	lex.InitHist()
	input := filepath.Join("testdata", goarch+".s")
	output := filepath.Join("testdata", goarch+".out")
	architecture, ctxt := setArch(goarch)
	lexer := lex.NewLexer(input, ctxt)
	parser := NewParser(ctxt, architecture, lexer)
	pList := obj.Linknewplist(ctxt)
	var ok bool
	testOut = new(bytes.Buffer) // The assembler writes -S output to this buffer.
	ctxt.Bso = obj.Binitw(os.Stdout)
	defer ctxt.Bso.Flush()
	ctxt.Diag = log.Fatalf
	obj.Binitw(ioutil.Discard)
	pList.Firstpc, ok = parser.Parse()
	if !ok {
		t.Fatalf("asm: %s assembly failed", goarch)
	}
	result := string(testOut.Bytes())
	expect, err := ioutil.ReadFile(output)
	// For Windows.
	result = strings.Replace(result, `testdata\`, `testdata/`, -1)
	if err != nil {
		t.Fatal(err)
	}
	if result != string(expect) {
		if false { // Enable to capture output.
			fmt.Printf("%s", result)
			os.Exit(1)
		}
		t.Errorf("%s failed: output differs", goarch)
		r := strings.Split(result, "\n")
		e := strings.Split(string(expect), "\n")
		if len(r) != len(e) {
			t.Errorf("%s: expected %d lines, got %d", goarch, len(e), len(r))
		}
		n := len(e)
		if n > len(r) {
			n = len(r)
		}
		for i := 0; i < n; i++ {
			if r[i] != e[i] {
				t.Errorf("%s:%d:\nexpected\n\t%s\ngot\n\t%s", output, i, e[i], r[i])
			}
		}
	}
}

func TestPPC64EndToEnd(t *testing.T) {
	testEndToEnd(t, "ppc64")
}

func TestARMEndToEnd(t *testing.T) {
	testEndToEnd(t, "arm")
}

func TestARM64EndToEnd(t *testing.T) {
	testEndToEnd(t, "arm64")
}

func TestAMD64EndToEnd(t *testing.T) {
	testEndToEnd(t, "amd64")
}

func Test386EndToEnd(t *testing.T) {
	testEndToEnd(t, "386")
}

func TestMIPS64EndToEnd(t *testing.T) {
	testEndToEnd(t, "mips64")
}
