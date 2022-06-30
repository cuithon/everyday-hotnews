// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

// fastlog2Table contains log2 approximations for 5 binary digits.
// This is used to implement fastlog2, which is used for heap sampling.

package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
)

func main() {
	var buf bytes.Buffer

	fmt.Fprintln(&buf, "// Code generated by mkfastlog2table.go; DO NOT EDIT.")
	fmt.Fprintln(&buf, "// Run go generate from src/runtime to update.")
	fmt.Fprintln(&buf, "// See mkfastlog2table.go for comments.")
	fmt.Fprintln(&buf)
	fmt.Fprintln(&buf, "package runtime")
	fmt.Fprintln(&buf)
	fmt.Fprintln(&buf, "const fastlogNumBits =", fastlogNumBits)
	fmt.Fprintln(&buf)

	fmt.Fprintln(&buf, "var fastlog2Table = [1<<fastlogNumBits + 1]float64{")
	table := computeTable()
	for _, t := range table {
		fmt.Fprintf(&buf, "\t%v,\n", t)
	}
	fmt.Fprintln(&buf, "}")

	if err := os.WriteFile("fastlog2table.go", buf.Bytes(), 0644); err != nil {
		log.Fatalln(err)
	}
}

const fastlogNumBits = 5

func computeTable() []float64 {
	fastlog2Table := make([]float64, 1<<fastlogNumBits+1)
	for i := 0; i <= (1 << fastlogNumBits); i++ {
		fastlog2Table[i] = log2(1.0 + float64(i)/(1<<fastlogNumBits))
	}
	return fastlog2Table
}

// log2 is a local copy of math.Log2 with an explicit float64 conversion
// to disable FMA. This lets us generate the same output on all platforms.
func log2(x float64) float64 {
	frac, exp := math.Frexp(x)
	// Make sure exact powers of two give an exact answer.
	// Don't depend on Log(0.5)*(1/Ln2)+exp being exactly exp-1.
	if frac == 0.5 {
		return float64(exp - 1)
	}
	return float64(nlog(frac)*(1/math.Ln2)) + float64(exp)
}

// nlog is a local copy of math.Log with explicit float64 conversions
// to disable FMA. This lets us generate the same output on all platforms.
func nlog(x float64) float64 {
	const (
		Ln2Hi = 6.93147180369123816490e-01 /* 3fe62e42 fee00000 */
		Ln2Lo = 1.90821492927058770002e-10 /* 3dea39ef 35793c76 */
		L1    = 6.666666666666735130e-01   /* 3FE55555 55555593 */
		L2    = 3.999999999940941908e-01   /* 3FD99999 9997FA04 */
		L3    = 2.857142874366239149e-01   /* 3FD24924 94229359 */
		L4    = 2.222219843214978396e-01   /* 3FCC71C5 1D8E78AF */
		L5    = 1.818357216161805012e-01   /* 3FC74664 96CB03DE */
		L6    = 1.531383769920937332e-01   /* 3FC39A09 D078C69F */
		L7    = 1.479819860511658591e-01   /* 3FC2F112 DF3E5244 */
	)

	// special cases
	switch {
	case math.IsNaN(x) || math.IsInf(x, 1):
		return x
	case x < 0:
		return math.NaN()
	case x == 0:
		return math.Inf(-1)
	}

	// reduce
	f1, ki := math.Frexp(x)
	if f1 < math.Sqrt2/2 {
		f1 *= 2
		ki--
	}
	f := f1 - 1
	k := float64(ki)

	// compute
	s := float64(f / (2 + f))
	s2 := float64(s * s)
	s4 := float64(s2 * s2)
	t1 := s2 * float64(L1+float64(s4*float64(L3+float64(s4*float64(L5+float64(s4*L7))))))
	t2 := s4 * float64(L2+float64(s4*float64(L4+float64(s4*L6))))
	R := float64(t1 + t2)
	hfsq := float64(0.5 * f * f)
	return float64(k*Ln2Hi) - ((hfsq - (float64(s*float64(hfsq+R)) + float64(k*Ln2Lo))) - f)
}
