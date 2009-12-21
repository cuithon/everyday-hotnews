// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vector


import "testing"


func TestStrZeroLenExp(t *testing.T) {
	a := new(StringVector)
	if a.Len() != 0 {
		t.Errorf("%T: B1) expected 0, got %d", a, a.Len())
	}
	if len(*a) != 0 {
		t.Errorf("%T: B2) expected 0, got %d", a, len(*a))
	}
	var b StringVector
	if b.Len() != 0 {
		t.Errorf("%T: B3) expected 0, got %d", b, b.Len())
	}
	if len(b) != 0 {
		t.Errorf("%T: B4) expected 0, got %d", b, len(b))
	}
}


func TestStrResizeExp(t *testing.T) {
	var a StringVector
	checkSizeExp(t, &a, 0, 0)
	checkSizeExp(t, a.Resize(0, 5), 0, 5)
	checkSizeExp(t, a.Resize(1, 0), 1, 5)
	checkSizeExp(t, a.Resize(10, 0), 10, 10)
	checkSizeExp(t, a.Resize(5, 0), 5, 10)
	checkSizeExp(t, a.Resize(3, 8), 3, 10)
	checkSizeExp(t, a.Resize(0, 100), 0, 100)
	checkSizeExp(t, a.Resize(11, 100), 11, 100)
}


func TestStrResize2Exp(t *testing.T) {
	var a StringVector
	checkSizeExp(t, &a, 0, 0)
	a.Push(int2StrValue(1))
	a.Push(int2StrValue(2))
	a.Push(int2StrValue(3))
	a.Push(int2StrValue(4))
	checkSizeExp(t, &a, 4, 4)
	checkSizeExp(t, a.Resize(10, 0), 10, 10)
	for i := 4; i < a.Len(); i++ {
		if a.At(i) != strzero {
			t.Errorf("%T: expected a.At(%d) == %v; found %v!", a, i, strzero, a.At(i))
		}
	}
	for i := 4; i < len(a); i++ {
		if a[i] != strzero {
			t.Errorf("%T: expected a[%d] == %v; found %v", a, i, strzero, a[i])
		}
	}
}


func checkStrZeroExp(t *testing.T, a *StringVector, i int) {
	for j := 0; j < i; j++ {
		if a.At(j) == strzero {
			t.Errorf("%T: 1 expected a.At(%d) == %d; found %v", a, j, j, a.At(j))
		}
		if (*a)[j] == strzero {
			t.Errorf("%T: 2 expected (*a)[%d] == %d; found %v", a, j, j, (*a)[j])
		}
	}
	for ; i < a.Len(); i++ {
		if a.At(i) != strzero {
			t.Errorf("%T: 3 expected a.At(%d) == %v; found %v", a, i, strzero, a.At(i))
		}
		if (*a)[i] != strzero {
			t.Errorf("%T: 4 expected (*a)[%d] == %v; found %v", a, i, strzero, (*a)[i])
		}
	}
}


func TestStrTrailingElementsExp(t *testing.T) {
	var a StringVector
	for i := 0; i < 10; i++ {
		a.Push(int2StrValue(i + 1))
	}
	checkStrZeroExp(t, &a, 10)
	checkSizeExp(t, &a, 10, 16)
	checkSizeExp(t, a.Resize(5, 0), 5, 16)
	checkSizeExp(t, a.Resize(10, 0), 10, 16)
	checkStrZeroExp(t, &a, 5)
}


func TestStrAccessExp(t *testing.T) {
	const n = 100
	var a StringVector
	a.Resize(n, 0)
	for i := 0; i < n; i++ {
		a.Set(i, int2StrValue(valExp(i)))
	}
	for i := 0; i < n; i++ {
		if elem2StrValue(a.At(i)) != int2StrValue(valExp(i)) {
			t.Error(i)
		}
	}
	var b StringVector
	b.Resize(n, 0)
	for i := 0; i < n; i++ {
		b[i] = int2StrValue(valExp(i))
	}
	for i := 0; i < n; i++ {
		if elem2StrValue(b[i]) != int2StrValue(valExp(i)) {
			t.Error(i)
		}
	}
}


func TestStrInsertDeleteClearExp(t *testing.T) {
	const n = 100
	var a StringVector

	for i := 0; i < n; i++ {
		if a.Len() != i {
			t.Errorf("T%: A) wrong Len() %d (expected %d)", a, a.Len(), i)
		}
		if len(a) != i {
			t.Errorf("T%: A) wrong len() %d (expected %d)", a, len(a), i)
		}
		a.Insert(0, int2StrValue(valExp(i)))
		if elem2StrValue(a.Last()) != int2StrValue(valExp(0)) {
			t.Error("T%: B", a)
		}
	}
	for i := n - 1; i >= 0; i-- {
		if elem2StrValue(a.Last()) != int2StrValue(valExp(0)) {
			t.Error("T%: C", a)
		}
		if elem2StrValue(a.At(0)) != int2StrValue(valExp(i)) {
			t.Error("T%: D", a)
		}
		if elem2StrValue(a[0]) != int2StrValue(valExp(i)) {
			t.Error("T%: D2", a)
		}
		a.Delete(0)
		if a.Len() != i {
			t.Errorf("T%: E) wrong Len() %d (expected %d)", a, a.Len(), i)
		}
		if len(a) != i {
			t.Errorf("T%: E) wrong len() %d (expected %d)", a, len(a), i)
		}
	}

	if a.Len() != 0 {
		t.Errorf("T%: F) wrong Len() %d (expected 0)", a, a.Len())
	}
	if len(a) != 0 {
		t.Errorf("T%: F) wrong len() %d (expected 0)", a, len(a))
	}
	for i := 0; i < n; i++ {
		a.Push(int2StrValue(valExp(i)))
		if a.Len() != i+1 {
			t.Errorf("T%: G) wrong Len() %d (expected %d)", a, a.Len(), i+1)
		}
		if len(a) != i+1 {
			t.Errorf("T%: G) wrong len() %d (expected %d)", a, len(a), i+1)
		}
		if elem2StrValue(a.Last()) != int2StrValue(valExp(i)) {
			t.Error("T%: H", a)
		}
	}
	a.Resize(0, 0)
	if a.Len() != 0 {
		t.Errorf("T%: I wrong Len() %d (expected 0)", a, a.Len())
	}
	if len(a) != 0 {
		t.Errorf("T%: I wrong len() %d (expected 0)", a, len(a))
	}

	const m = 5
	for j := 0; j < m; j++ {
		a.Push(int2StrValue(j))
		for i := 0; i < n; i++ {
			x := valExp(i)
			a.Push(int2StrValue(x))
			if elem2StrValue(a.Pop()) != int2StrValue(x) {
				t.Error("T%: J", a)
			}
			if a.Len() != j+1 {
				t.Errorf("T%: K) wrong Len() %d (expected %d)", a, a.Len(), j+1)
			}
			if len(a) != j+1 {
				t.Errorf("T%: K) wrong len() %d (expected %d)", a, len(a), j+1)
			}
		}
	}
	if a.Len() != m {
		t.Errorf("T%: L) wrong Len() %d (expected %d)", a, a.Len(), m)
	}
	if len(a) != m {
		t.Errorf("T%: L) wrong len() %d (expected %d)", a, len(a), m)
	}
}


func verify_sliceStrExp(t *testing.T, x *StringVector, elt, i, j int) {
	for k := i; k < j; k++ {
		if elem2StrValue(x.At(k)) != int2StrValue(elt) {
			t.Errorf("T%: M) wrong [%d] element %v (expected %v)", x, k, elem2StrValue(x.At(k)), int2StrValue(elt))
		}
	}

	s := x.Slice(i, j)
	for k, n := 0, j-i; k < n; k++ {
		if elem2StrValue(s.At(k)) != int2StrValue(elt) {
			t.Errorf("T%: N) wrong [%d] element %v (expected %v)", x, k, elem2StrValue(x.At(k)), int2StrValue(elt))
		}
	}
}


func verify_patternStrExp(t *testing.T, x *StringVector, a, b, c int) {
	n := a + b + c
	if x.Len() != n {
		t.Errorf("T%: O) wrong Len() %d (expected %d)", x, x.Len(), n)
	}
	if len(*x) != n {
		t.Errorf("T%: O) wrong len() %d (expected %d)", x, len(*x), n)
	}
	verify_sliceStrExp(t, x, 0, 0, a)
	verify_sliceStrExp(t, x, 1, a, a+b)
	verify_sliceStrExp(t, x, 0, a+b, n)
}


func make_vectorStrExp(elt, len int) *StringVector {
	x := new(StringVector).Resize(len, 0)
	for i := 0; i < len; i++ {
		x.Set(i, int2StrValue(elt))
	}
	return x
}


func TestStrInsertVectorExp(t *testing.T) {
	// 1
	a := make_vectorStrExp(0, 0)
	b := make_vectorStrExp(1, 10)
	a.InsertVector(0, b)
	verify_patternStrExp(t, a, 0, 10, 0)
	// 2
	a = make_vectorStrExp(0, 10)
	b = make_vectorStrExp(1, 0)
	a.InsertVector(5, b)
	verify_patternStrExp(t, a, 5, 0, 5)
	// 3
	a = make_vectorStrExp(0, 10)
	b = make_vectorStrExp(1, 3)
	a.InsertVector(3, b)
	verify_patternStrExp(t, a, 3, 3, 7)
	// 4
	a = make_vectorStrExp(0, 10)
	b = make_vectorStrExp(1, 1000)
	a.InsertVector(8, b)
	verify_patternStrExp(t, a, 8, 1000, 2)
}


func TestStrDoExp(t *testing.T) {
	const n = 25
	const salt = 17
	a := new(StringVector).Resize(n, 0)
	for i := 0; i < n; i++ {
		a.Set(i, int2StrValue(salt*i))
	}
	count := 0
	a.Do(func(e interface{}) {
		i := intf2StrValue(e)
		if i != int2StrValue(count*salt) {
			t.Error(tname(a), "value at", count, "should be", count*salt, "not", i)
		}
		count++
	})
	if count != n {
		t.Error(tname(a), "should visit", n, "values; did visit", count)
	}

	b := new(StringVector).Resize(n, 0)
	for i := 0; i < n; i++ {
		(*b)[i] = int2StrValue(salt * i)
	}
	count = 0
	b.Do(func(e interface{}) {
		i := intf2StrValue(e)
		if i != int2StrValue(count*salt) {
			t.Error(tname(b), "b) value at", count, "should be", count*salt, "not", i)
		}
		count++
	})
	if count != n {
		t.Error(tname(b), "b) should visit", n, "values; did visit", count)
	}

	var c StringVector
	c.Resize(n, 0)
	for i := 0; i < n; i++ {
		c[i] = int2StrValue(salt * i)
	}
	count = 0
	c.Do(func(e interface{}) {
		i := intf2StrValue(e)
		if i != int2StrValue(count*salt) {
			t.Error(tname(c), "c) value at", count, "should be", count*salt, "not", i)
		}
		count++
	})
	if count != n {
		t.Error(tname(c), "c) should visit", n, "values; did visit", count)
	}

}


func TestStrIterExp(t *testing.T) {
	const Len = 100
	x := new(StringVector).Resize(Len, 0)
	for i := 0; i < Len; i++ {
		x.Set(i, int2StrValue(i*i))
	}
	i := 0
	for v := range x.Iter() {
		if elem2StrValue(v) != int2StrValue(i*i) {
			t.Error(tname(x), "Iter expected", i*i, "got", elem2StrValue(v))
		}
		i++
	}
	if i != Len {
		t.Error(tname(x), "Iter stopped at", i, "not", Len)
	}
	y := new(StringVector).Resize(Len, 0)
	for i := 0; i < Len; i++ {
		(*y)[i] = int2StrValue(i * i)
	}
	i = 0
	for v := range y.Iter() {
		if elem2StrValue(v) != int2StrValue(i*i) {
			t.Error(tname(y), "y, Iter expected", i*i, "got", elem2StrValue(v))
		}
		i++
	}
	if i != Len {
		t.Error(tname(y), "y, Iter stopped at", i, "not", Len)
	}
	var z StringVector
	z.Resize(Len, 0)
	for i := 0; i < Len; i++ {
		z[i] = int2StrValue(i * i)
	}
	i = 0
	for v := range z.Iter() {
		if elem2StrValue(v) != int2StrValue(i*i) {
			t.Error(tname(z), "z, Iter expected", i*i, "got", elem2StrValue(v))
		}
		i++
	}
	if i != Len {
		t.Error(tname(z), "z, Iter stopped at", i, "not", Len)
	}
}

func TestStrVectorData(t *testing.T) {
	// verify Data() returns a slice of a copy, not a slice of the original vector
	const Len = 10
	var src StringVector
	for i := 0; i < Len; i++ {
		src.Push(int2StrValue(i * i))
	}
	dest := src.Data()
	for i := 0; i < Len; i++ {
		src[i] = int2StrValue(-1)
		v := elem2StrValue(dest[i])
		if v != int2StrValue(i*i) {
			t.Error(tname(src), "expected", i*i, "got", v)
		}
	}
}
