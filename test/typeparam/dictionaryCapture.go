// run -gcflags=-G=3

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test situations where functions/methods are not
// immediately called and we need to capture the dictionary
// required for later invocation.

// TODO: copy this test file, add -l to gcflags.

package main

func main() {
	functions()
	methodExpressions()
	methodValues()
	interfaceMethods()
}

func g0[T any](x T) {
}
func g1[T any](x T) T {
	return x
}
func g2[T any](x T) (T, T) {
	return x, x
}

func functions() {
	f0 := g0[int]
	f0(7)
	f1 := g1[int]
	is7(f1(7))
	f2 := g2[int]
	is77(f2(7))
}

func is7(x int) {
	if x != 7 {
		println(x)
		panic("assertion failed")
	}
}
func is77(x, y int) {
	if x != 7 || y != 7 {
		println(x,y)
		panic("assertion failed")
	}
}

type s[T any] struct {
	a T
}

func (x s[T]) g0() {
}
func (x s[T]) g1() T {
	return x.a
}
func (x s[T]) g2() (T, T) {
	return x.a, x.a
}

func methodExpressions() {
	x := s[int]{a:7}
	f0 := s[int].g0
	f0(x)
	f1 := s[int].g1
	is7(f1(x))
	f2 := s[int].g2
	is77(f2(x))
}

func methodValues() {
	x := s[int]{a:7}
	f0 := x.g0
	f0()
	f1 := x.g1
	is7(f1())
	f2 := x.g2
	is77(f2())
}

var x interface{
	g0()
	g1()int
	g2()(int,int)
} = s[int]{a:7}
var y interface{} = s[int]{a:7}

func interfaceMethods() {
	x.g0()
	is7(x.g1())
	is77(x.g2())
	y.(interface{g0()}).g0()
	is7(y.(interface{g1()int}).g1())
	is77(y.(interface{g2()(int,int)}).g2())
}
