// $G $D/$F.go && $L $F.$A && ./$A.out

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

func main() {
	i5 := 5;

	switch {  // compiler crash fixable with 'switch true'
	case i5 < 5: dummy := 0;
	case i5 == 5: dummy := 0;
	case i5 > 5: dummy := 0;
	}
}
/*
Segmentation fault
*/
