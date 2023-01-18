// Code generated by "go run generator.go"; DO NOT EDIT.

// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types_test

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"testing"
)

var (
	H   = flag.Int("H", 5, "Hilbert matrix size")
	out = flag.String("out", "", "write generated program to out")
)

func TestHilbert(t *testing.T) {
	// generate source
	src := program(*H, *out)
	if *out != "" {
		os.WriteFile(*out, src, 0666)
		return
	}

	mustTypecheck("hilbert.go", string(src), nil, nil)
}

func program(n int, out string) []byte {
	var g gen

	g.p(`// Code generated by: go test -run=Hilbert -H=%d -out=%q. DO NOT EDIT.

// +`+`build ignore

// This program tests arbitrary precision constant arithmetic
// by generating the constant elements of a Hilbert matrix H,
// its inverse I, and the product P = H*I. The product should
// be the identity matrix.
package main

func main() {
	if !ok {
		printProduct()
		return
	}
	println("PASS")
}

`, n, out)
	g.hilbert(n)
	g.inverse(n)
	g.product(n)
	g.verify(n)
	g.printProduct(n)
	g.binomials(2*n - 1)
	g.factorials(2*n - 1)

	return g.Bytes()
}

type gen struct {
	bytes.Buffer
}

func (g *gen) p(format string, args ...interface{}) {
	fmt.Fprintf(&g.Buffer, format, args...)
}

func (g *gen) hilbert(n int) {
	g.p(`// Hilbert matrix, n = %d
const (
`, n)
	for i := 0; i < n; i++ {
		g.p("\t")
		for j := 0; j < n; j++ {
			if j > 0 {
				g.p(", ")
			}
			g.p("h%d_%d", i, j)
		}
		if i == 0 {
			g.p(" = ")
			for j := 0; j < n; j++ {
				if j > 0 {
					g.p(", ")
				}
				g.p("1.0/(iota + %d)", j+1)
			}
		}
		g.p("\n")
	}
	g.p(")\n\n")
}

func (g *gen) inverse(n int) {
	g.p(`// Inverse Hilbert matrix
const (
`)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			s := "+"
			if (i+j)&1 != 0 {
				s = "-"
			}
			g.p("\ti%d_%d = %s%d * b%d_%d * b%d_%d * b%d_%d * b%d_%d\n",
				i, j, s, i+j+1, n+i, n-j-1, n+j, n-i-1, i+j, i, i+j, i)
		}
		g.p("\n")
	}
	g.p(")\n\n")
}

func (g *gen) product(n int) {
	g.p(`// Product matrix
const (
`)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			g.p("\tp%d_%d = ", i, j)
			for k := 0; k < n; k++ {
				if k > 0 {
					g.p(" + ")
				}
				g.p("h%d_%d*i%d_%d", i, k, k, j)
			}
			g.p("\n")
		}
		g.p("\n")
	}
	g.p(")\n\n")
}

func (g *gen) verify(n int) {
	g.p(`// Verify that product is the identity matrix
const ok =
`)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j == 0 {
				g.p("\t")
			} else {
				g.p(" && ")
			}
			v := 0
			if i == j {
				v = 1
			}
			g.p("p%d_%d == %d", i, j, v)
		}
		g.p(" &&\n")
	}
	g.p("\ttrue\n\n")

	// verify ok at type-check time
	if *out == "" {
		g.p("const _ = assert(ok)\n\n")
	}
}

func (g *gen) printProduct(n int) {
	g.p("func printProduct() {\n")
	for i := 0; i < n; i++ {
		g.p("\tprintln(")
		for j := 0; j < n; j++ {
			if j > 0 {
				g.p(", ")
			}
			g.p("p%d_%d", i, j)
		}
		g.p(")\n")
	}
	g.p("}\n\n")
}

func (g *gen) binomials(n int) {
	g.p(`// Binomials
const (
`)
	for j := 0; j <= n; j++ {
		if j > 0 {
			g.p("\n")
		}
		for k := 0; k <= j; k++ {
			g.p("\tb%d_%d = f%d / (f%d*f%d)\n", j, k, j, k, j-k)
		}
	}
	g.p(")\n\n")
}

func (g *gen) factorials(n int) {
	g.p(`// Factorials
const (
	f0 = 1
	f1 = 1
`)
	for i := 2; i <= n; i++ {
		g.p("\tf%d = f%d * %d\n", i, i-1, i)
	}
	g.p(")\n\n")
}
