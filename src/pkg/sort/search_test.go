// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sort

import "testing"


func f(a []int, x int) func(int) bool {
	return func(i int) bool {
		return a[i] <= x
	}
}


var data = []int{0: -10, 1: -5, 2: 0, 3: 1, 4: 2, 5: 3, 6: 5, 7: 7, 8: 11, 9: 100, 10: 100, 11: 100, 12: 1000, 13: 10000}

var tests = []struct {
	name string
	n    int
	f    func(int) bool
	i    int
}{
	{"empty", 0, nil, 0},
	{"1 1", 1, func(i int) bool { return i <= 1 }, 0},
	{"1 true", 1, func(i int) bool { return false }, 0},
	{"1 false", 1, func(i int) bool { return true }, 0},
	{"1e9 991", 1e9, func(i int) bool { return i <= 991 }, 991},
	{"1e9 true", 1e9, func(i int) bool { return false }, 0},
	{"1e9 false", 1e9, func(i int) bool { return true }, 1e9 - 1},
	{"data -20", len(data), f(data, -20), 0},
	{"data -10", len(data), f(data, -10), 0},
	{"data -9", len(data), f(data, -9), 0},
	{"data -6", len(data), f(data, -6), 0},
	{"data -5", len(data), f(data, -5), 1},
	{"data 3", len(data), f(data, 3), 5},
	{"data 99", len(data), f(data, 99), 8},
	{"data 100", len(data), f(data, 100), 11},
	{"data 101", len(data), f(data, 101), 11},
	{"data 10000", len(data), f(data, 10000), 13},
	{"data 10001", len(data), f(data, 10001), 13},
	{"descending a", 7, func(i int) bool { return []int{99, 99, 59, 42, 7, 0, -1, -1}[i] >= 7 }, 4},
	{"descending 7", 1e9, func(i int) bool { return 1e9-i >= 7 }, 1e9 - 7},
}


func TestSearch(t *testing.T) {
	for _, e := range tests {
		i := Search(e.n, e.f)
		if i != e.i {
			t.Errorf("%s: expected index %d; got %d", e.name, e.i, i)
		}
	}
}


// Smoke tests for convenience wrappers - not comprehensive.

var fdata = []float{0: -3.14, 1: 0, 2: 1, 3: 2, 4: 1000.7}
var sdata = []string{0: "f", 1: "foo", 2: "foobar", 3: "x"}

var wrappertests = []struct {
	name   string
	result int
	i      int
}{
	{"SearchInts", SearchInts(data, 11), 8},
	{"SearchFloats", SearchFloats(fdata, 2.1), 3},
	{"SearchStrings", SearchStrings(sdata, ""), 0},
	{"IntArray.Search", IntArray(data).Search(0), 2},
	{"FloatArray.Search", FloatArray(fdata).Search(2.0), 3},
	{"StringArray.Search", StringArray(sdata).Search("x"), 3},
}


func TestSearchWrappers(t *testing.T) {
	for _, e := range wrappertests {
		if e.result != e.i {
			t.Errorf("%s: expected index %d; got %d", e.name, e.i, e.result)
		}
	}
}
