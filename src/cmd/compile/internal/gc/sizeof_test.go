// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !nacl

package gc

import (
	"reflect"
	"testing"
	"unsafe"
)

// Assert that the size of important structures do not change unexpectedly.

func TestSizeof(t *testing.T) {
	const _64bit = unsafe.Sizeof(uintptr(0)) == 8

	var tests = []struct {
		val    interface{} // type as a value
		_32bit uintptr     // size on 32bit platforms
		_64bit uintptr     // size on 64bit platforms
	}{
		{Func{}, 96, 160},
		{Name{}, 36, 56},
		{Param{}, 28, 56},
		{Node{}, 84, 136},
		{Sym{}, 60, 104},
		{Type{}, 52, 88},
		{MapType{}, 20, 40},
		{ForwardType{}, 20, 32},
		{FuncType{}, 28, 48},
		{StructType{}, 12, 24},
		{InterType{}, 4, 8},
		{ChanType{}, 8, 16},
		{ArrayType{}, 12, 16},
		{DDDFieldType{}, 4, 8},
		{FuncArgsType{}, 4, 8},
		{ChanArgsType{}, 4, 8},
		{PtrType{}, 4, 8},
		{SliceType{}, 4, 8},
	}

	for _, tt := range tests {
		want := tt._32bit
		if _64bit {
			want = tt._64bit
		}
		got := reflect.TypeOf(tt.val).Size()
		if want != got {
			t.Errorf("unsafe.Sizeof(%T) = %d, want %d", tt.val, got, want)
		}
	}
}
