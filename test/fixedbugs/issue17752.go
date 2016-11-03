// run

// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

func f(m map[string]int) int {
	return m["a"]
}

func g(m map[[8]string]int) int {
	return m[[8]string{"a", "a", "a", "a", "a", "a", "a", "a"}]
}

func main() {
	m := map[[8]string]int{}
	g(m)
}
