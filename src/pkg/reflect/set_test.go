// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reflect_test

import (
	"bytes"
	"go/ast"
	"io"
	. "reflect"
	"testing"
	"unsafe"
)

type MyBuffer bytes.Buffer

func TestImplicitMapConversion(t *testing.T) {
	// Test implicit conversions in MapIndex and SetMapIndex.
	{
		// direct
		m := make(map[int]int)
		mv := ValueOf(m)
		mv.SetMapIndex(ValueOf(1), ValueOf(2))
		x, ok := m[1]
		if x != 2 {
			t.Errorf("#1 after SetMapIndex(1,2): %d, %t (map=%v)", x, ok, m)
		}
		if n := mv.MapIndex(ValueOf(1)).Interface().(int); n != 2 {
			t.Errorf("#1 MapIndex(1) = %d", n)
		}
	}
	{
		// convert interface key
		m := make(map[interface{}]int)
		mv := ValueOf(m)
		mv.SetMapIndex(ValueOf(1), ValueOf(2))
		x, ok := m[1]
		if x != 2 {
			t.Errorf("#2 after SetMapIndex(1,2): %d, %t (map=%v)", x, ok, m)
		}
		if n := mv.MapIndex(ValueOf(1)).Interface().(int); n != 2 {
			t.Errorf("#2 MapIndex(1) = %d", n)
		}
	}
	{
		// convert interface value
		m := make(map[int]interface{})
		mv := ValueOf(m)
		mv.SetMapIndex(ValueOf(1), ValueOf(2))
		x, ok := m[1]
		if x != 2 {
			t.Errorf("#3 after SetMapIndex(1,2): %d, %t (map=%v)", x, ok, m)
		}
		if n := mv.MapIndex(ValueOf(1)).Interface().(int); n != 2 {
			t.Errorf("#3 MapIndex(1) = %d", n)
		}
	}
	{
		// convert both interface key and interface value
		m := make(map[interface{}]interface{})
		mv := ValueOf(m)
		mv.SetMapIndex(ValueOf(1), ValueOf(2))
		x, ok := m[1]
		if x != 2 {
			t.Errorf("#4 after SetMapIndex(1,2): %d, %t (map=%v)", x, ok, m)
		}
		if n := mv.MapIndex(ValueOf(1)).Interface().(int); n != 2 {
			t.Errorf("#4 MapIndex(1) = %d", n)
		}
	}
	{
		// convert both, with non-empty interfaces
		m := make(map[io.Reader]io.Writer)
		mv := ValueOf(m)
		b1 := new(bytes.Buffer)
		b2 := new(bytes.Buffer)
		mv.SetMapIndex(ValueOf(b1), ValueOf(b2))
		x, ok := m[b1]
		if x != b2 {
			t.Errorf("#5 after SetMapIndex(b1, b2): %p (!= %p), %t (map=%v)", x, b2, ok, m)
		}
		if p := mv.MapIndex(ValueOf(b1)).Elem().Pointer(); p != uintptr(unsafe.Pointer(b2)) {
			t.Errorf("#5 MapIndex(b1) = %p want %p", p, b2)
		}
	}
	{
		// convert channel direction
		m := make(map[<-chan int]chan int)
		mv := ValueOf(m)
		c1 := make(chan int)
		c2 := make(chan int)
		mv.SetMapIndex(ValueOf(c1), ValueOf(c2))
		x, ok := m[c1]
		if x != c2 {
			t.Errorf("#6 after SetMapIndex(c1, c2): %p (!= %p), %t (map=%v)", x, c2, ok, m)
		}
		if p := mv.MapIndex(ValueOf(c1)).Pointer(); p != ValueOf(c2).Pointer() {
			t.Errorf("#6 MapIndex(c1) = %p want %p", p, c2)
		}
	}
	{
		// convert identical underlying types
		// TODO(rsc): Should be able to define MyBuffer here.
		// 6l prints very strange messages about .this.Bytes etc
		// when we do that though, so MyBuffer is defined
		// at top level.
		m := make(map[*MyBuffer]*bytes.Buffer)
		mv := ValueOf(m)
		b1 := new(MyBuffer)
		b2 := new(bytes.Buffer)
		mv.SetMapIndex(ValueOf(b1), ValueOf(b2))
		x, ok := m[b1]
		if x != b2 {
			t.Errorf("#7 after SetMapIndex(b1, b2): %p (!= %p), %t (map=%v)", x, b2, ok, m)
		}
		if p := mv.MapIndex(ValueOf(b1)).Pointer(); p != uintptr(unsafe.Pointer(b2)) {
			t.Errorf("#7 MapIndex(b1) = %p want %p", p, b2)
		}
	}

}

func TestImplicitSetConversion(t *testing.T) {
	// Assume TestImplicitMapConversion covered the basics.
	// Just make sure conversions are being applied at all.
	var r io.Reader
	b := new(bytes.Buffer)
	rv := ValueOf(&r).Elem()
	rv.Set(ValueOf(b))
	if r != b {
		t.Errorf("after Set: r=%T(%v)", r, r)
	}
}

func TestImplicitSendConversion(t *testing.T) {
	c := make(chan io.Reader, 10)
	b := new(bytes.Buffer)
	ValueOf(c).Send(ValueOf(b))
	if bb := <-c; bb != b {
		t.Errorf("Received %p != %p", bb, b)
	}
}

func TestImplicitCallConversion(t *testing.T) {
	// Arguments must be assignable to parameter types.
	fv := ValueOf(io.WriteString)
	b := new(bytes.Buffer)
	fv.Call([]Value{ValueOf(b), ValueOf("hello world")})
	if b.String() != "hello world" {
		t.Errorf("After call: string=%q want %q", b.String(), "hello world")
	}
}

func TestImplicitAppendConversion(t *testing.T) {
	// Arguments must be assignable to the slice's element type.
	s := []io.Reader{}
	sv := ValueOf(&s).Elem()
	b := new(bytes.Buffer)
	sv.Set(Append(sv, ValueOf(b)))
	if len(s) != 1 || s[0] != b {
		t.Errorf("after append: s=%v want [%p]", s, b)
	}
}

var implementsTests = []struct {
	x interface{}
	t interface{}
	b bool
}{
	{new(*bytes.Buffer), new(io.Reader), true},
	{new(bytes.Buffer), new(io.Reader), false},
	{new(*bytes.Buffer), new(io.ReaderAt), false},
	{new(*ast.Ident), new(ast.Expr), true},
}

func TestImplements(t *testing.T) {
	for _, tt := range implementsTests {
		xv := TypeOf(tt.x).Elem()
		xt := TypeOf(tt.t).Elem()
		if b := xv.Implements(xt); b != tt.b {
			t.Errorf("(%s).Implements(%s) = %v, want %v", xv.String(), xt.String(), b, tt.b)
		}
	}
}

var assignableTests = []struct {
	x interface{}
	t interface{}
	b bool
}{
	{new(chan int), new(<-chan int), true},
	{new(<-chan int), new(chan int), false},
	{new(*int), new(IntPtr), true},
	{new(IntPtr), new(*int), true},
	{new(IntPtr), new(IntPtr1), false},
	// test runs implementsTests too
}

type IntPtr *int
type IntPtr1 *int

func TestAssignableTo(t *testing.T) {
	for _, tt := range append(assignableTests, implementsTests...) {
		xv := TypeOf(tt.x).Elem()
		xt := TypeOf(tt.t).Elem()
		if b := xv.AssignableTo(xt); b != tt.b {
			t.Errorf("(%s).AssignableTo(%s) = %v, want %v", xv.String(), xt.String(), b, tt.b)
		}
	}
}
