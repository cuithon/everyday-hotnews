// $G $D/$F.go && $L $F.$A && ./$A.out

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main


export type T struct {
	x, y int;
}

func (t *T) m(a int, b float) int {
	return (t.x+a) * (t.y+int(b));
}

func main() {
	var t *T = new(T);
	t.x = 1;
	t.y = 2;
	r10 := t.m(1, 3.0);
}
/*
bug11.go:16: fatal error: walktype: switch 1 unknown op CALLMETH l(16) <int32>INT32
*/
