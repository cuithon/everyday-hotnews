// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

// This program is run via "go generate" (via a directive in sort.go)
// to generate implementation variants of the underlying sorting algorithm.
// When passed the -generic flag it generates generic variants of sorting;
// otherwise it generates the non-generic variants used by the sort package.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"text/template"
)

type Variant struct {
	// Name is the variant name: should be unique among variants.
	Name string

	// Path is the file path into which the generator will emit the code for this
	// variant.
	Path string

	// Package is the package this code will be emitted into.
	Package string

	// Imports is the imports needed for this package.
	Imports string

	// FuncSuffix is appended to all function names in this variant's code. All
	// suffixes should be unique within a package.
	FuncSuffix string

	// DataType is the type of the data parameter of functions in this variant's
	// code.
	DataType string

	// TypeParam is the optional type parameter for the function.
	TypeParam string

	// ExtraParam is an extra parameter to pass to the function. Should begin with
	// ", " to separate from other params.
	ExtraParam string

	// ExtraArg is an extra argument to pass to calls between functions; typically
	// it invokes ExtraParam. Should begin with ", " to separate from other args.
	ExtraArg string

	// Funcs is a map of functions used from within the template. The following
	// functions are expected to exist:
	//
	//    Less (name, i, j):
	//      emits a comparison expression that checks if the value `name` at
	//      index `i` is smaller than at index `j`.
	//
	//    Swap (name, i, j):
	//      emits a statement that performs a data swap between elements `i` and
	//      `j` of the value `name`.
	Funcs template.FuncMap
}

func main() {
	genGeneric := flag.Bool("generic", false, "generate generic versions")
	flag.Parse()

	if *genGeneric {
		generate(&Variant{
			Name:       "generic_ordered",
			Path:       "zsortordered.go",
			Package:    "slices",
			Imports:    "import \"constraints\"\n",
			FuncSuffix: "Ordered",
			TypeParam:  "[E constraints.Ordered]",
			ExtraParam: "",
			ExtraArg:   "",
			DataType:   "[]E",
			Funcs: template.FuncMap{
				"Less": func(name, i, j string) string {
					return fmt.Sprintf("(%s[%s] < %s[%s])", name, i, name, j)
				},
				"Swap": func(name, i, j string) string {
					return fmt.Sprintf("%s[%s], %s[%s] = %s[%s], %s[%s]", name, i, name, j, name, j, name, i)
				},
			},
		})

		generate(&Variant{
			Name:       "generic_func",
			Path:       "zsortanyfunc.go",
			Package:    "slices",
			FuncSuffix: "LessFunc",
			TypeParam:  "[E any]",
			ExtraParam: ", less func(a, b E) bool",
			ExtraArg:   ", less",
			DataType:   "[]E",
			Funcs: template.FuncMap{
				"Less": func(name, i, j string) string {
					return fmt.Sprintf("less(%s[%s], %s[%s])", name, i, name, j)
				},
				"Swap": func(name, i, j string) string {
					return fmt.Sprintf("%s[%s], %s[%s] = %s[%s], %s[%s]", name, i, name, j, name, j, name, i)
				},
			},
		})
	} else {
		generate(&Variant{
			Name:       "interface",
			Path:       "zsortinterface.go",
			Package:    "sort",
			Imports:    "",
			FuncSuffix: "",
			TypeParam:  "",
			ExtraParam: "",
			ExtraArg:   "",
			DataType:   "Interface",
			Funcs: template.FuncMap{
				"Less": func(name, i, j string) string {
					return fmt.Sprintf("%s.Less(%s, %s)", name, i, j)
				},
				"Swap": func(name, i, j string) string {
					return fmt.Sprintf("%s.Swap(%s, %s)", name, i, j)
				},
			},
		})

		generate(&Variant{
			Name:       "func",
			Path:       "zsortfunc.go",
			Package:    "sort",
			Imports:    "",
			FuncSuffix: "_func",
			TypeParam:  "",
			ExtraParam: "",
			ExtraArg:   "",
			DataType:   "lessSwap",
			Funcs: template.FuncMap{
				"Less": func(name, i, j string) string {
					return fmt.Sprintf("%s.Less(%s, %s)", name, i, j)
				},
				"Swap": func(name, i, j string) string {
					return fmt.Sprintf("%s.Swap(%s, %s)", name, i, j)
				},
			},
		})
	}
}

// generate generates the code for variant `v` into a file named by `v.Path`.
func generate(v *Variant) {
	// Parse templateCode anew for each variant because Parse requires Funcs to be
	// registered, and it helps type-check the funcs.
	tmpl, err := template.New("gen").Funcs(v.Funcs).Parse(templateCode)
	if err != nil {
		log.Fatal("template Parse:", err)
	}

	var out bytes.Buffer
	err = tmpl.Execute(&out, v)
	if err != nil {
		log.Fatal("template Execute:", err)
	}

	formatted, err := format.Source(out.Bytes())
	if err != nil {
		log.Fatal("format:", err)
	}

	if err := os.WriteFile(v.Path, formatted, 0644); err != nil {
		log.Fatal("WriteFile:", err)
	}
}

var templateCode = `// Code generated by gen_sort_variants.go; DO NOT EDIT.

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package {{.Package}}

{{.Imports}}

// insertionSort{{.FuncSuffix}} sorts data[a:b] using insertion sort.
func insertionSort{{.FuncSuffix}}{{.TypeParam}}(data {{.DataType}}, a, b int {{.ExtraParam}}) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && {{Less "data" "j" "j-1"}}; j-- {
			{{Swap "data" "j" "j-1"}}
		}
	}
}

// siftDown{{.FuncSuffix}} implements the heap property on data[lo:hi].
// first is an offset into the array where the root of the heap lies.
func siftDown{{.FuncSuffix}}{{.TypeParam}}(data {{.DataType}}, lo, hi, first int {{.ExtraParam}}) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && {{Less "data" "first+child" "first+child+1"}} {
			child++
		}
		if !{{Less "data" "first+root" "first+child"}} {
			return
		}
		{{Swap "data" "first+root" "first+child"}}
		root = child
	}
}

func heapSort{{.FuncSuffix}}{{.TypeParam}}(data {{.DataType}}, a, b int {{.ExtraParam}}) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown{{.FuncSuffix}}(data, i, hi, first {{.ExtraArg}})
	}

	// Pop elements, largest first, into end of data.
	for i := hi - 1; i >= 0; i-- {
		{{Swap "data" "first" "first+i"}}
		siftDown{{.FuncSuffix}}(data, lo, i, first {{.ExtraArg}})
	}
}

// Quicksort, loosely following Bentley and McIlroy,
// "Engineering a Sort Function" SP&E November 1993.

// medianOfThree{{.FuncSuffix}} moves the median of the three values data[m0], data[m1], data[m2] into data[m1].
func medianOfThree{{.FuncSuffix}}{{.TypeParam}}(data {{.DataType}}, m1, m0, m2 int {{.ExtraParam}}) {
	// sort 3 elements
	if {{Less "data" "m1" "m0"}} {
		{{Swap "data" "m1" "m0"}}
	}
	// data[m0] <= data[m1]
	if {{Less "data" "m2" "m1"}} {
		{{Swap "data" "m2" "m1"}}
		// data[m0] <= data[m2] && data[m1] < data[m2]
		if {{Less "data" "m1" "m0"}} {
			{{Swap "data" "m1" "m0"}}
		}
	}
	// now data[m0] <= data[m1] <= data[m2]
}

func swapRange{{.FuncSuffix}}{{.TypeParam}}(data {{.DataType}}, a, b, n int {{.ExtraParam}}) {
	for i := 0; i < n; i++ {
		{{Swap "data" "a+i" "b+i"}}
	}
}

func doPivot{{.FuncSuffix}}{{.TypeParam}}(data {{.DataType}}, lo, hi int {{.ExtraParam}}) (midlo, midhi int) {
	m := int(uint(lo+hi) >> 1) // Written like this to avoid integer overflow.
	if hi-lo > 40 {
		// Tukey's "Ninther" median of three medians of three.
		s := (hi - lo) / 8
		medianOfThree{{.FuncSuffix}}(data, lo, lo+s, lo+2*s {{.ExtraArg}})
		medianOfThree{{.FuncSuffix}}(data, m, m-s, m+s {{.ExtraArg}})
		medianOfThree{{.FuncSuffix}}(data, hi-1, hi-1-s, hi-1-2*s {{.ExtraArg}})
	}
	medianOfThree{{.FuncSuffix}}(data, lo, m, hi-1 {{.ExtraArg}})

	// Invariants are:
	//	data[lo] = pivot (set up by ChoosePivot)
	//	data[lo < i < a] < pivot
	//	data[a <= i < b] <= pivot
	//	data[b <= i < c] unexamined
	//	data[c <= i < hi-1] > pivot
	//	data[hi-1] >= pivot
	pivot := lo
	a, c := lo+1, hi-1

	for ; a < c && {{Less "data" "a" "pivot"}}; a++ {
	}
	b := a
	for {
		for ; b < c && !{{Less "data" "pivot" "b"}}; b++ { // data[b] <= pivot
		}
		for ; b < c && {{Less "data" "pivot" "c-1"}}; c-- { // data[c-1] > pivot
		}
		if b >= c {
			break
		}
		// data[b] > pivot; data[c-1] <= pivot
		{{Swap "data" "b" "c-1"}}
		b++
		c--
	}
	// If hi-c<3 then there are duplicates (by property of median of nine).
	// Let's be a bit more conservative, and set border to 5.
	protect := hi-c < 5
	if !protect && hi-c < (hi-lo)/4 {
		// Lets test some points for equality to pivot
		dups := 0
		if !{{Less "data" "pivot" "hi-1"}} { // data[hi-1] = pivot
			{{Swap "data" "c" "hi-1"}}
			c++
			dups++
		}
		if !{{Less "data" "b-1" "pivot"}} { // data[b-1] = pivot
			b--
			dups++
		}
		// m-lo = (hi-lo)/2 > 6
		// b-lo > (hi-lo)*3/4-1 > 8
		// ==> m < b ==> data[m] <= pivot
		if !{{Less "data" "m" "pivot"}} { // data[m] = pivot
			{{Swap "data" "m" "b-1"}}
			b--
			dups++
		}
		// if at least 2 points are equal to pivot, assume skewed distribution
		protect = dups > 1
	}
	if protect {
		// Protect against a lot of duplicates
		// Add invariant:
		//	data[a <= i < b] unexamined
		//	data[b <= i < c] = pivot
		for {
			for ; a < b && !{{Less "data" "b-1" "pivot"}}; b-- { // data[b] == pivot
			}
			for ; a < b && {{Less "data" "a" "pivot"}}; a++ { // data[a] < pivot
			}
			if a >= b {
				break
			}
			// data[a] == pivot; data[b-1] < pivot
			{{Swap "data" "a" "b-1"}}
			a++
			b--
		}
	}
	// Swap pivot into middle
	{{Swap "data" "pivot" "b-1"}}
	return b - 1, c
}

func quickSort{{.FuncSuffix}}{{.TypeParam}}(data {{.DataType}}, a, b, maxDepth int {{.ExtraParam}}) {
	for b-a > 12 { // Use ShellSort for slices <= 12 elements
		if maxDepth == 0 {
			heapSort{{.FuncSuffix}}(data, a, b {{.ExtraArg}})
			return
		}
		maxDepth--
		mlo, mhi := doPivot{{.FuncSuffix}}(data, a, b {{.ExtraArg}})
		// Avoiding recursion on the larger subproblem guarantees
		// a stack depth of at most lg(b-a).
		if mlo-a < b-mhi {
			quickSort{{.FuncSuffix}}(data, a, mlo, maxDepth {{.ExtraArg}})
			a = mhi // i.e., quickSort{{.FuncSuffix}}(data, mhi, b)
		} else {
			quickSort{{.FuncSuffix}}(data, mhi, b, maxDepth {{.ExtraArg}})
			b = mlo // i.e., quickSort{{.FuncSuffix}}(data, a, mlo)
		}
	}
	if b-a > 1 {
		// Do ShellSort pass with gap 6
		// It could be written in this simplified form cause b-a <= 12
		for i := a + 6; i < b; i++ {
			if {{Less "data" "i" "i-6"}} {
				{{Swap "data" "i" "i-6"}}
			}
		}
		insertionSort{{.FuncSuffix}}(data, a, b {{.ExtraArg}})
	}
}

func stable{{.FuncSuffix}}{{.TypeParam}}(data {{.DataType}}, n int {{.ExtraParam}}) {
	blockSize := 20 // must be > 0
	a, b := 0, blockSize
	for b <= n {
		insertionSort{{.FuncSuffix}}(data, a, b {{.ExtraArg}})
		a = b
		b += blockSize
	}
	insertionSort{{.FuncSuffix}}(data, a, n {{.ExtraArg}})

	for blockSize < n {
		a, b = 0, 2*blockSize
		for b <= n {
			symMerge{{.FuncSuffix}}(data, a, a+blockSize, b {{.ExtraArg}})
			a = b
			b += 2 * blockSize
		}
		if m := a + blockSize; m < n {
			symMerge{{.FuncSuffix}}(data, a, m, n {{.ExtraArg}})
		}
		blockSize *= 2
	}
}

// symMerge{{.FuncSuffix}} merges the two sorted subsequences data[a:m] and data[m:b] using
// the SymMerge algorithm from Pok-Son Kim and Arne Kutzner, "Stable Minimum
// Storage Merging by Symmetric Comparisons", in Susanne Albers and Tomasz
// Radzik, editors, Algorithms - ESA 2004, volume 3221 of Lecture Notes in
// Computer Science, pages 714-723. Springer, 2004.
//
// Let M = m-a and N = b-n. Wolog M < N.
// The recursion depth is bound by ceil(log(N+M)).
// The algorithm needs O(M*log(N/M + 1)) calls to data.Less.
// The algorithm needs O((M+N)*log(M)) calls to data.Swap.
//
// The paper gives O((M+N)*log(M)) as the number of assignments assuming a
// rotation algorithm which uses O(M+N+gcd(M+N)) assignments. The argumentation
// in the paper carries through for Swap operations, especially as the block
// swapping rotate uses only O(M+N) Swaps.
//
// symMerge assumes non-degenerate arguments: a < m && m < b.
// Having the caller check this condition eliminates many leaf recursion calls,
// which improves performance.
func symMerge{{.FuncSuffix}}{{.TypeParam}}(data {{.DataType}}, a, m, b int {{.ExtraParam}}) {
	// Avoid unnecessary recursions of symMerge
	// by direct insertion of data[a] into data[m:b]
	// if data[a:m] only contains one element.
	if m-a == 1 {
		// Use binary search to find the lowest index i
		// such that data[i] >= data[a] for m <= i < b.
		// Exit the search loop with i == b in case no such index exists.
		i := m
		j := b
		for i < j {
			h := int(uint(i+j) >> 1)
			if {{Less "data" "h" "a"}} {
				i = h + 1
			} else {
				j = h
			}
		}
		// Swap values until data[a] reaches the position before i.
		for k := a; k < i-1; k++ {
			{{Swap "data" "k" "k+1"}}
		}
		return
	}

	// Avoid unnecessary recursions of symMerge
	// by direct insertion of data[m] into data[a:m]
	// if data[m:b] only contains one element.
	if b-m == 1 {
		// Use binary search to find the lowest index i
		// such that data[i] > data[m] for a <= i < m.
		// Exit the search loop with i == m in case no such index exists.
		i := a
		j := m
		for i < j {
			h := int(uint(i+j) >> 1)
			if !{{Less "data" "m" "h"}} {
				i = h + 1
			} else {
				j = h
			}
		}
		// Swap values until data[m] reaches the position i.
		for k := m; k > i; k-- {
			{{Swap "data" "k" "k-1"}}
		}
		return
	}

	mid := int(uint(a+b) >> 1)
	n := mid + m
	var start, r int
	if m > mid {
		start = n - b
		r = mid
	} else {
		start = a
		r = m
	}
	p := n - 1

	for start < r {
		c := int(uint(start+r) >> 1)
		if !{{Less "data" "p-c" "c"}} {
			start = c + 1
		} else {
			r = c
		}
	}

	end := n - start
	if start < m && m < end {
		rotate{{.FuncSuffix}}(data, start, m, end {{.ExtraArg}})
	}
	if a < start && start < mid {
		symMerge{{.FuncSuffix}}(data, a, start, mid {{.ExtraArg}})
	}
	if mid < end && end < b {
		symMerge{{.FuncSuffix}}(data, mid, end, b {{.ExtraArg}})
	}
}

// rotate{{.FuncSuffix}} rotates two consecutive blocks u = data[a:m] and v = data[m:b] in data:
// Data of the form 'x u v y' is changed to 'x v u y'.
// rotate performs at most b-a many calls to data.Swap,
// and it assumes non-degenerate arguments: a < m && m < b.
func rotate{{.FuncSuffix}}{{.TypeParam}}(data {{.DataType}}, a, m, b int {{.ExtraParam}}) {
	i := m - a
	j := b - m

	for i != j {
		if i > j {
			swapRange{{.FuncSuffix}}(data, m-i, m, j {{.ExtraArg}})
			i -= j
		} else {
			swapRange{{.FuncSuffix}}(data, m-i, m+j-i, i {{.ExtraArg}})
			j -= i
		}
	}
	// i == j
	swapRange{{.FuncSuffix}}(data, m-i, m, i {{.ExtraArg}})
}
`
