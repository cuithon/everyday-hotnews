// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io_test

import (
	. "io";
	"os";
	"strings";
	"testing";
)

func checkSize(t *testing.T, path string, size uint64) {
	dir, err := os.Stat(path);
	if err != nil {
		t.Fatalf("Stat %q (looking for size %d): %s", path, size, err)
	}
	if dir.Size != size {
		t.Errorf("Stat %q: size %d want %d", path, dir.Size, size)
	}
}

func TestReadFile(t *testing.T) {
	filename := "rumpelstilzchen";
	contents, err := ReadFile(filename);
	if err == nil {
		t.Fatalf("ReadFile %s: error expected, none found", filename)
	}

	filename = "utils_test.go";
	contents, err = ReadFile(filename);
	if err != nil {
		t.Fatalf("ReadFile %s: %v", filename, err)
	}

	checkSize(t, filename, uint64(len(contents)));
}

func TestWriteFile(t *testing.T) {
	filename := "_obj/rumpelstilzchen";
	data := "Programming today is a race between software engineers striving to "
		"build bigger and better idiot-proof programs, and the Universe trying "
		"to produce bigger and better idiots. So far, the Universe is winning.";

	if err := WriteFile(filename, strings.Bytes(data), 0644); err != nil {
		t.Fatalf("WriteFile %s: %v", filename, err)
	}

	contents, err := ReadFile(filename);
	if err != nil {
		t.Fatalf("ReadFile %s: %v", filename, err)
	}

	if string(contents) != data {
		t.Fatalf("contents = %q\nexpected = %q", string(contents), data)
	}

	// cleanup
	os.Remove(filename);	// ignore error
}


func TestReadDir(t *testing.T) {
	dirname := "rumpelstilzchen";
	_, err := ReadDir(dirname);
	if err == nil {
		t.Fatalf("ReadDir %s: error expected, none found", dirname)
	}

	dirname = ".";
	list, err := ReadDir(dirname);
	if err != nil {
		t.Fatalf("ReadDir %s: %v", dirname, err)
	}

	foundTest := false;
	foundObj := false;
	for _, dir := range list {
		switch {
		case dir.IsRegular() && dir.Name == "utils_test.go":
			foundTest = true
		case dir.IsDirectory() && dir.Name == "_obj":
			foundObj = true
		}
	}
	if !foundTest {
		t.Fatalf("ReadDir %s: test file not found", dirname)
	}
	if !foundObj {
		t.Fatalf("ReadDir %s: _obj directory not found", dirname)
	}
}
