// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iterable

import (
	"testing";
)

func TestArrayTypes(t *testing.T) {
	// Test that conversion works correctly.
	bytes := ByteArray([]byte{1, 2, 3});
	if x := Data(bytes)[1].(byte); x != 2 {
		t.Error("Data(bytes)[1].(byte) = %v, want 2", x)
	}
	ints := IntArray([]int{1, 2, 3});
	if x := Data(ints)[2].(int); x != 3 {
		t.Error("Data(ints)[2].(int) = %v, want 3", x)
	}
	floats := FloatArray([]float{1, 2, 3});
	if x := Data(floats)[0].(float); x != 1 {
		t.Error("Data(floats)[0].(float) = %v, want 1", x)
	}
	strings := StringArray([]string{"a", "b", "c"});
	if x := Data(strings)[1].(string); x != "b" {
		t.Error(`Data(strings)[1].(string) = %q, want "b"`, x)
	}
}

var oneToFive = IntArray{1, 2, 3, 4, 5}

func isNegative(n interface{}) bool	{ return n.(int) < 0 }
func isPositive(n interface{}) bool	{ return n.(int) > 0 }
func isAbove3(n interface{}) bool	{ return n.(int) > 3 }
func isEven(n interface{}) bool		{ return n.(int)%2 == 0 }
func doubler(n interface{}) interface{}	{ return n.(int) * 2 }
func addOne(n interface{}) interface{}	{ return n.(int) + 1 }
func adder(acc interface{}, n interface{}) interface{} {
	return acc.(int) + n.(int)
}

// A stream of the natural numbers: 0, 1, 2, 3, ...
type integerStream struct{}

func (i integerStream) Iter() <-chan interface{} {
	ch := make(chan interface{});
	go func() {
		for i := 0; ; i++ {
			ch <- i
		}
	}();
	return ch;
}

func TestAll(t *testing.T) {
	if !All(oneToFive, isPositive) {
		t.Error("All(oneToFive, isPositive) == false")
	}
	if All(oneToFive, isAbove3) {
		t.Error("All(oneToFive, isAbove3) == true")
	}
}

func TestAny(t *testing.T) {
	if Any(oneToFive, isNegative) {
		t.Error("Any(oneToFive, isNegative) == true")
	}
	if !Any(oneToFive, isEven) {
		t.Error("Any(oneToFive, isEven) == false")
	}
}

func assertArraysAreEqual(t *testing.T, res []interface{}, expected []int) {
	if len(res) != len(expected) {
		t.Errorf("len(res) = %v, want %v", len(res), len(expected));
		goto missing;
	}
	for i := range res {
		if v := res[i].(int); v != expected[i] {
			t.Errorf("res[%v] = %v, want %v", i, v, expected[i]);
			goto missing;
		}
	}
	return;
missing:
	t.Errorf("res = %v\nwant  %v", res, expected);
}

func TestFilter(t *testing.T) {
	ints := integerStream{};
	moreInts := Filter(ints, isAbove3).Iter();
	res := make([]interface{}, 3);
	for i := 0; i < 3; i++ {
		res[i] = <-moreInts
	}
	assertArraysAreEqual(t, res, []int{4, 5, 6});
}

func TestFind(t *testing.T) {
	ints := integerStream{};
	first := Find(ints, isAbove3);
	if first.(int) != 4 {
		t.Errorf("Find(ints, isAbove3) = %v, want 4", first)
	}
}

func TestInject(t *testing.T) {
	res := Inject(oneToFive, 0, adder);
	if res.(int) != 15 {
		t.Errorf("Inject(oneToFive, 0, adder) = %v, want 15", res)
	}
}

func TestMap(t *testing.T) {
	res := Data(Map(Map(oneToFive, doubler), addOne));
	assertArraysAreEqual(t, res, []int{3, 5, 7, 9, 11});
}

func TestPartition(t *testing.T) {
	ti, fi := Partition(oneToFive, isEven);
	assertArraysAreEqual(t, Data(ti), []int{2, 4});
	assertArraysAreEqual(t, Data(fi), []int{1, 3, 5});
}

func TestTake(t *testing.T) {
	res := Take(oneToFive, 2);
	assertArraysAreEqual(t, Data(res), []int{1, 2});
	assertArraysAreEqual(t, Data(res), []int{1, 2});	// second test to ensure that .Iter() returns a new channel

	// take none
	res = Take(oneToFive, 0);
	assertArraysAreEqual(t, Data(res), []int{});

	// try to take more than available
	res = Take(oneToFive, 20);
	assertArraysAreEqual(t, Data(res), oneToFive);
}

func TestTakeWhile(t *testing.T) {
	// take some
	res := TakeWhile(oneToFive, func(v interface{}) bool { return v.(int) <= 3 });
	assertArraysAreEqual(t, Data(res), []int{1, 2, 3});
	assertArraysAreEqual(t, Data(res), []int{1, 2, 3});	// second test to ensure that .Iter() returns a new channel

	// take none
	res = TakeWhile(oneToFive, func(v interface{}) bool { return v.(int) > 3000 });
	assertArraysAreEqual(t, Data(res), []int{});

	// take all
	res = TakeWhile(oneToFive, func(v interface{}) bool { return v.(int) < 3000 });
	assertArraysAreEqual(t, Data(res), oneToFive);
}

func TestDrop(t *testing.T) {
	// drop none
	res := Drop(oneToFive, 0);
	assertArraysAreEqual(t, Data(res), oneToFive);
	assertArraysAreEqual(t, Data(res), oneToFive);	// second test to ensure that .Iter() returns a new channel

	// drop some
	res = Drop(oneToFive, 2);
	assertArraysAreEqual(t, Data(res), []int{3, 4, 5});
	assertArraysAreEqual(t, Data(res), []int{3, 4, 5});	// second test to ensure that .Iter() returns a new channel

	// drop more than available
	res = Drop(oneToFive, 88);
	assertArraysAreEqual(t, Data(res), []int{});
}

func TestDropWhile(t *testing.T) {
	// drop some
	res := DropWhile(oneToFive, func(v interface{}) bool { return v.(int) < 3 });
	assertArraysAreEqual(t, Data(res), []int{3, 4, 5});
	assertArraysAreEqual(t, Data(res), []int{3, 4, 5});	// second test to ensure that .Iter() returns a new channel

	// test case where all elements are dropped
	res = DropWhile(oneToFive, func(v interface{}) bool { return v.(int) < 100 });
	assertArraysAreEqual(t, Data(res), []int{});

	// test case where none are dropped
	res = DropWhile(oneToFive, func(v interface{}) bool { return v.(int) > 1000 });
	assertArraysAreEqual(t, Data(res), oneToFive);
}
