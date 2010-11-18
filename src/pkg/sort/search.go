// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements binary search.

package sort

// Search uses binary search to find and return the smallest index i
// in [0, n) at which f(i) is false, assuming that on the range [0, n), 
// f(i) == false implies f(i+1) == false.  That is, Search requires that
// f is true for some (possibly empty) prefix of the input range [0, n)
// and then false for the (possibly empty) remainder; Search returns
// the first false index.  If there is no such index, Search returns n.
// Search calls f(i) only for i in the range [0, n).
//
// A common use of Search is to find the index i for a value x in
// a sorted, indexable data structure like an array or slice.
// In this case, the argument f, typically a closure, captures the value
// to be searched for, and how the data structure is indexed and
// ordered.
//
// For instance, given a slice data sorted in ascending order,
// the call Search(len(data), func(i int) bool { return data[i] < 23 })
// returns the smallest index i such that data[i] >= 23.  If the caller
// wants to find whether 23 is in the slice, it must test data[i] == 23
// separately.
//
// Searching data sorted in descending order would use the >
// operator instead of the < operator.
//
// To complete the example above, the following code tries to find the value
// x in an integer slice data sorted in ascending order:
//
//	x := 23
//	i := sort.Search(len(data), func(i int) bool { return data[i] < x })
//	if i < len(data) && data[i] == x {
//		// x is present at data[i]
//	} else {
//		// x is not present in data,
//		// but i is the index where it would be inserted.
//	}
//
// As a more whimsical example, this program guesses your number:
//
//	func GuessingGame() {
//		var s string
//		fmt.Printf("Pick an integer from 0 to 100.\n")
//		answer := sort.Search(100, func(i int) bool {
//			fmt.Printf("Is your number > %d? ", i)
//			fmt.Scanf("%s", &s)
//			return s != "" && s[0] == 'y'
//		})
//		fmt.Printf("Your number is %d.\n", answer)
//	}
//
func Search(n int, f func(int) bool) int {
	// Define f(-1) == true and f(n) == false.
	// Invariant: f(i-1) == true, f(j) == false.
	i, j := 0, n
	for i < j {
		h := i + (j-i)/2 // avoid overflow when computing h
		// i ≤ h < j
		if f(h) {
			i = h + 1 // preserves f(i-1) == true
		} else {
			j = h // preserves f(j) == false
		}
	}
	// i == j, f(i-1) == true, and f(j) (= f(i)) == false  =>  answer is i.
	return i
}


// Convenience wrappers for common cases.

// SearchInts searches x in a sorted slice of ints and returns the index
// as specified by Search. The array must be sorted in ascending order.
//
func SearchInts(a []int, x int) int {
	return Search(len(a), func(i int) bool { return a[i] < x })
}


// SearchFloats searches x in a sorted slice of floats and returns the index
// as specified by Search. The array must be sorted in ascending order.
// 
func SearchFloats(a []float, x float) int {
	return Search(len(a), func(i int) bool { return a[i] < x })
}


// SearchStrings searches x in a sorted slice of strings and returns the index
// as specified by Search. The array must be sorted in ascending order.
// 
func SearchStrings(a []string, x string) int {
	return Search(len(a), func(i int) bool { return a[i] < x })
}


// Search returns the result of applying SearchInts to the receiver and x.
func (p IntArray) Search(x int) int { return SearchInts(p, x) }


// Search returns the result of applying SearchFloats to the receiver and x.
func (p FloatArray) Search(x float) int { return SearchFloats(p, x) }


// Search returns the result of applying SearchStrings to the receiver and x.
func (p StringArray) Search(x string) int { return SearchStrings(p, x) }
