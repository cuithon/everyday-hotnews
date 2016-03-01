// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	flag.Parse()
	if flag.NArg() > 0 {
		f, err := os.Create(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		os.Stdout = f
	}
	fmt.Printf("// generated by wbfat_gen.go; use go generate\n\n")
	fmt.Printf("package runtime\n")
	for i := uint(2); i <= 4; i++ {
		for j := 1; j < 1<<i; j++ {
			fmt.Printf("\n//go:nosplit\n")
			fmt.Printf("func writebarrierfat%0*b(dst *[%d]uintptr, _ uintptr, src [%d]uintptr) {\n", int(i), j, i, i)
			for k := uint(0); k < i; k++ {
				if j&(1<<(i-1-k)) != 0 {
					fmt.Printf("\twritebarrierptr(&dst[%d], src[%d])\n", k, k)
				} else {
					fmt.Printf("\tdst[%d] = src[%d]\n", k, k)
				}
			}
			fmt.Printf("}\n")
		}
	}
}
