// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reflect

import (
	"io";
	"os";
	"reflect";
	"testing";
	"unsafe";
)

var doprint bool = false

func is_digit(c uint8) bool {
	return '0' <= c && c <= '9'
}

// streq, but '@' in t matches a string of digits
func match(s, t string) bool {
	for i, j := 0, 0; i < len(s) && j < len(t); i, j = i+1, j+1 {
		if s[i] == t[j] {
			continue
		}
		if is_digit(s[i]) && t[j] == '@' {
			for is_digit(s[i+1]) {
				i++
			}
		} else {
			return false
		}
	}
	return true;
}

func assert(s, t string) {
	if doprint {
		println(t)
	}
	if !match(s, t) {
		panicln(s, t)
	}
}

func typedump(s, t string) {
	typ := ParseTypeString("", s);
	assert(typeToString(typ, true), t);
}

func valuedump(s, t string) {
	typ := ParseTypeString("", s);
	v := NewZeroValue(typ);
	if v == nil {
		panicln("valuedump", s);
	}
	switch v.Kind() {
	case IntKind:
		v.(IntValue).Set(132);
	case Int8Kind:
		v.(Int8Value).Set(8);
	case Int16Kind:
		v.(Int16Value).Set(16);
	case Int32Kind:
		v.(Int32Value).Set(32);
	case Int64Kind:
		v.(Int64Value).Set(64);
	case UintKind:
		v.(UintValue).Set(132);
	case Uint8Kind:
		v.(Uint8Value).Set(8);
	case Uint16Kind:
		v.(Uint16Value).Set(16);
	case Uint32Kind:
		v.(Uint32Value).Set(32);
	case Uint64Kind:
		v.(Uint64Value).Set(64);
	case FloatKind:
		v.(FloatValue).Set(3200.0);
	case Float32Kind:
		v.(Float32Value).Set(32.1);
	case Float64Kind:
		v.(Float64Value).Set(64.2);
	case StringKind:
		v.(StringValue).Set("stringy cheese");
	case BoolKind:
		v.(BoolValue).Set(true);
	}
	assert(valueToString(v), t);
}

type T struct { a int; b float64; c string; d *int }

func TestAll(tt *testing.T) {	// TODO(r): wrap up better
	var s string;
	var t Type;

	// Types
	typedump("missing", "$missing$");
	typedump("int", "int");
	typedump("int8", "int8");
	typedump("int16", "int16");
	typedump("int32", "int32");
	typedump("int64", "int64");
	typedump("uint", "uint");
	typedump("uint8", "uint8");
	typedump("uint16", "uint16");
	typedump("uint32", "uint32");
	typedump("uint64", "uint64");
	typedump("float", "float");
	typedump("float32", "float32");
	typedump("float64", "float64");
	typedump("int8", "int8");
	typedump("whoknows.whatsthis", "$missing$");
	typedump("**int8", "**int8");
	typedump("**P.integer", "**P.integer");
	typedump("[32]int32", "[32]int32");
	typedump("[]int8", "[]int8");
	typedump("map[string]int32", "map[string]int32");
	typedump("chan<-string", "chan<-string");
	typedump("struct {c chan *int32; d float32}", "struct{c chan*int32; d float32}");
	typedump("func(a int8, b int32)", "func(a int8, b int32)");
	typedump("struct {c func(? chan *P.integer, ? *int8)}", "struct{c func(chan*P.integer, *int8)}");
	typedump("struct {a int8; b int32}", "struct{a int8; b int32}");
	typedump("struct {a int8; b int8; b int32}", "struct{a int8; b int8; b int32}");
	typedump("struct {a int8; b int8; c int8; b int32}", "struct{a int8; b int8; c int8; b int32}");
	typedump("struct {a int8; b int8; c int8; d int8; b int32}", "struct{a int8; b int8; c int8; d int8; b int32}");
	typedump("struct {a int8; b int8; c int8; d int8; e int8; b int32}", "struct{a int8; b int8; c int8; d int8; e int8; b int32}");
	typedump("struct {a int8 \"hi there\"; }", "struct{a int8 \"hi there\"}");
	typedump("struct {a int8 \"hi \\x00there\\t\\n\\\"\\\\\"; }", "struct{a int8 \"hi \\x00there\\t\\n\\\"\\\\\"}");
	typedump("struct {f func(args ...)}", "struct{f func(args ...)}");
	typedump("interface { a(? func(? func(? int) int) func(? func(? int)) int); b() }", "interface{a (func(func(int)(int))(func(func(int))(int))); b ()}");

	// Values
	valuedump("int8", "8");
	valuedump("int16", "16");
	valuedump("int32", "32");
	valuedump("int64", "64");
	valuedump("uint8", "8");
	valuedump("uint16", "16");
	valuedump("uint32", "32");
	valuedump("uint64", "64");
	valuedump("float32", "32.1");
	valuedump("float64", "64.2");
	valuedump("string", "stringy cheese");
	valuedump("bool", "true");
	valuedump("*int8", "*int8(0)");
	valuedump("**int8", "**int8(0)");
	valuedump("[5]int32", "[5]int32{0, 0, 0, 0, 0}");
	valuedump("**P.integer", "**P.integer(0)");
	valuedump("map[string]int32", "map[string]int32{<can't iterate on maps>}");
	valuedump("chan<-string", "chan<-string");
	valuedump("struct {c chan *int32; d float32}", "struct{c chan*int32; d float32}{chan*int32, 0}");
	valuedump("func(a int8, b int32)", "func(a int8, b int32)(0)");
	valuedump("struct {c func(? chan *P.integer, ? *int8)}", "struct{c func(chan*P.integer, *int8)}{func(chan*P.integer, *int8)(0)}");
	valuedump("struct {a int8; b int32}", "struct{a int8; b int32}{0, 0}");
	valuedump("struct {a int8; b int8; b int32}", "struct{a int8; b int8; b int32}{0, 0, 0}");

	{	var tmp = 123;
		value := NewValue(tmp);
		assert(valueToString(value), "123");
	}
	{	var tmp = 123.4;
		value := NewValue(tmp);
		assert(valueToString(value), "123.4");
	}
	{
		var tmp = byte(123);
		value := NewValue(tmp);
		assert(valueToString(value), "123");
		assert(typeToString(value.Type(), false), "uint8");
	}
	{	var tmp = "abc";
		value := NewValue(tmp);
		assert(valueToString(value), "abc");
	}
	{
		var i int = 7;
		var tmp = &T{123, 456.75, "hello", &i};
		value := NewValue(tmp);
		assert(valueToString(value.(PtrValue).Sub()), "reflect.T{123, 456.75, hello, *int(@)}");
	}
	{
		type C chan *T;	// TODO: should not be necessary
		var tmp = new(C);
		value := NewValue(tmp);
		assert(valueToString(value), "*reflect.C·all_test(@)");
	}
//	{
//		type A [10]int;
//		var tmp A = A{1,2,3,4,5,6,7,8,9,10};
//		value := NewValue(&tmp);
//		assert(valueToString(value.(PtrValue).Sub()), "reflect.A·all_test{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}");
//		value.(PtrValue).Sub().(ArrayValue).Elem(4).(IntValue).Set(123);
//		assert(valueToString(value.(PtrValue).Sub()), "reflect.A·all_test{1, 2, 3, 4, 123, 6, 7, 8, 9, 10}");
//	}
	{
		type AA []int;
		var tmp = AA{1,2,3,4,5,6,7,8,9,10};
		value := NewValue(&tmp);	// TODO: NewValue(tmp) too
		assert(valueToString(value.(PtrValue).Sub()), "reflect.AA·all_test{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}");
		value.(PtrValue).Sub().(ArrayValue).Elem(4).(IntValue).Set(123);
		assert(valueToString(value.(PtrValue).Sub()), "reflect.AA·all_test{1, 2, 3, 4, 123, 6, 7, 8, 9, 10}");
	}

	{
		var ip *int32;
		var i int32 = 1234;
		vip := NewValue(&ip);
		vi := NewValue(i);
		vip.(PtrValue).Sub().(PtrValue).SetSub(vi);
		if *ip != 1234 {
			panicln("SetSub failure", *ip);
		}
	}

	var pt PtrType;
	var st StructType;
	var mt MapType;
	var at ArrayType;
	var ct ChanType;
	var name string;
	var typ Type;
	var tag string;
	var offset int;

	// Type strings
	t = ParseTypeString("", "int8");
	assert(t.String(), "int8");

	t = ParseTypeString("", "*int8");
	assert(t.String(), "*int8");
	pt = t.(PtrType);
	assert(pt.Sub().String(), "int8");

	t = ParseTypeString("", "*struct {c chan *int32; d float32}");
	assert(t.String(), "*struct {c chan *int32; d float32}");
	pt = t.(PtrType);
	assert(pt.Sub().String(), "struct {c chan *int32; d float32}");
	st = pt.Sub().(StructType);
	name, typ, tag, offset = st.Field(0);
	assert(typ.String(), "chan *int32");
	name, typ, tag, offset = st.Field(1);
	assert(typ.String(), "float32");

	t = ParseTypeString("", "interface {a() *int}");
	assert(t.String(), "interface {a() *int}");

	t = ParseTypeString("", "func(a int8, b int32)");
	assert(t.String(), "func(a int8, b int32)");

	t = ParseTypeString("", "func(a int8, b int32) float");
	assert(t.String(), "func(a int8, b int32) float");

	t = ParseTypeString("", "func(a int8, b int32) (a float, b float)");
	assert(t.String(), "func(a int8, b int32) (a float, b float)");

	t = ParseTypeString("", "[32]int32");
	assert(t.String(), "[32]int32");
	at = t.(ArrayType);
	assert(at.Elem().String(), "int32");

	t = ParseTypeString("", "map[string]*int32");
	assert(t.String(), "map[string]*int32");
	mt = t.(MapType);
	assert(mt.Key().String(), "string");
	assert(mt.Elem().String(), "*int32");

	t = ParseTypeString("", "chan<-string");
	assert(t.String(), "chan<-string");
	ct = t.(ChanType);
	assert(ct.Elem().String(), "string");

	// make sure tag strings are not part of element type
	t = ParseTypeString("", "struct{d []uint32 \"TAG\"}");
	st = t.(StructType);
	name, typ, tag, offset = st.Field(0);
	assert(typ.String(), "[]uint32");

	t = ParseTypeString("", "[]int32");
	v := NewSliceValue(t.(ArrayType), 5, 10);
	t1 := ParseTypeString("", "*[]int32");
	v1 := NewZeroValue(t1);
	if v1 == nil { panic("V1 is nil"); }
	v1.(PtrValue).SetSub(v);
	a := v1.Interface().(*[]int32);
	println(&a, len(a), cap(a));
	for i := 0; i < len(a); i++ {
		v.Elem(i).(Int32Value).Set(int32(i));
	}
	for i := 0; i < len(a); i++ {
		println(a[i]);
	}
}

func TestInterfaceGet(t *testing.T) {
	var inter struct { e interface{ } };
	inter.e = 123.456;
	v1 := NewValue(&inter);
	v2 := v1.(PtrValue).Sub().(StructValue).Field(0);
	assert(v2.Type().String(), "interface { }");
	i2 := v2.(InterfaceValue).Get();
	v3 := NewValue(i2);
	assert(v3.Type().String(), "float");
}

func TestInterfaceValue(t *testing.T) {
	var inter struct { e interface{ } };
	inter.e = 123.456;
	v1 := NewValue(&inter);
	v2 := v1.(PtrValue).Sub().(StructValue).Field(0);
	assert(v2.Type().String(), "interface { }");
	v3 := v2.(InterfaceValue).Value();
	assert(v3.Type().String(), "float");

	i3 := v2.Interface();
	if f, ok := i3.(float); !ok {
		a, typ, c := unsafe.Reflect(i3);
		t.Error("v2.Interface() did not return float, got ", typ);
	}
}

func TestFunctionValue(t *testing.T) {
	v := NewValue(func() {});
	if v.Interface() != v.Interface() {
		t.Fatalf("TestFunction != itself");
	}
	assert(v.Type().String(), "func()");
}

func TestCopyArray(t *testing.T) {
	a := []int{ 1, 2, 3, 4, 10, 9, 8, 7 };
	b := []int{ 11, 22, 33, 44, 1010, 99, 88, 77, 66, 55, 44 };
	c := []int{ 11, 22, 33, 44, 1010, 99, 88, 77, 66, 55, 44 };
	va := NewValue(&a);
	vb := NewValue(&b);
	for i := 0; i < len(b); i++ {
		if b[i] != c[i] {
			t.Fatalf("b != c before test");
		}
	}
	for tocopy := 1; tocopy <= 7; tocopy++ {
		vb.(PtrValue).Sub().(ArrayValue).CopyFrom(va.(PtrValue).Sub().(ArrayValue), tocopy);
		for i := 0; i < tocopy; i++ {
			if a[i] != b[i] {
				t.Errorf("1 tocopy=%d a[%d]=%d, b[%d]=%d",
					tocopy, i, a[i], i, b[i]);
			}
		}
		for i := tocopy; i < len(b); i++ {
			if b[i] != c[i] {
				if i < len(a) {
					t.Errorf("2 tocopy=%d a[%d]=%d, b[%d]=%d, c[%d]=%d",
						tocopy, i, a[i], i, b[i], i, c[i]);
				} else {
					t.Errorf("3 tocopy=%d b[%d]=%d, c[%d]=%d",
						tocopy, i, b[i], i, c[i]);
				}
			} else {
				t.Logf("tocopy=%d elem %d is okay\n", tocopy, i);
			}
		}
	}
}

func TestBigUnnamedStruct(t *testing.T) {
	b := struct{a,b,c,d int64}{1, 2, 3, 4};
	v := NewValue(b);
	b1 := v.Interface().(struct{a,b,c,d int64});
	if b1.a != b.a || b1.b != b.b || b1.c != b.c || b1.d != b.d {
		t.Errorf("NewValue(%v).Interface().(Big) = %v", b, b1);
	}
}

type big struct {
	a, b, c, d, e int64
}
func TestBigStruct(t *testing.T) {
	b := big{1, 2, 3, 4, 5};
	v := NewValue(b);
	b1 := v.Interface().(big);
	if b1.a != b.a || b1.b != b.b || b1.c != b.c || b1.d != b.d || b1.e != b.e {
		t.Errorf("NewValue(%v).Interface().(big) = %v", b, b1);
	}
}

type Basic struct {
	x int;
	y float32
}

type NotBasic Basic

type Recursive struct {
	x int;
	r *Recursive
}

type Complex struct {
	a int;
	b [3]*Complex;
	c *string;
	d map[float]float
}

type DeepEqualTest struct {
	a, b interface{};
	eq bool;
}

var deepEqualTests = []DeepEqualTest {
	// Equalities
	DeepEqualTest{ 1, 1, true },
	DeepEqualTest{ int32(1), int32(1), true },
	DeepEqualTest{ 0.5, 0.5, true },
	DeepEqualTest{ float32(0.5), float32(0.5), true },
	DeepEqualTest{ "hello", "hello", true },
	DeepEqualTest{ make([]int, 10), make([]int, 10), true },
	DeepEqualTest{ &[3]int{ 1, 2, 3 }, &[3]int{ 1, 2, 3 }, true },
	DeepEqualTest{ Basic{ 1, 0.5 }, Basic{ 1, 0.5 }, true },
	// Inequalities
	DeepEqualTest{ 1, 2, false },
	DeepEqualTest{ int32(1), int32(2), false },
	DeepEqualTest{ 0.5, 0.6, false },
	DeepEqualTest{ float32(0.5), float32(0.6), false },
	DeepEqualTest{ "hello", "hey", false },
	DeepEqualTest{ make([]int, 10), make([]int, 11), false },
	DeepEqualTest{ &[3]int{ 1, 2, 3 }, &[3]int{ 1, 2, 4 }, false },
	DeepEqualTest{ Basic{ 1, 0.5 }, Basic{ 1, 0.6 }, false },
	// Mismatched types
	DeepEqualTest{ 1, 1.0, false },
	DeepEqualTest{ int32(1), int64(1), false },
	DeepEqualTest{ 0.5, "hello", false },
	DeepEqualTest{ []int{ 1, 2, 3 }, [3]int{ 1, 2, 3 }, false },
	DeepEqualTest{ &[3]interface{} { 1, 2, 4 }, &[3]interface{} { 1, 2, "s" }, false },
	DeepEqualTest{ Basic{ 1, 0.5 }, NotBasic{ 1, 0.5 }, false },
}

func TestDeepEqual(t *testing.T) {
	for i, test := range deepEqualTests {
		if r := DeepEqual(test.a, test.b); r != test.eq {
			t.Errorf("DeepEqual(%v, %v) = %v, want %v", test.a, test.b, r, test.eq);
		}
	}
}

func TestDeepEqualRecursiveStruct(t *testing.T) {
	a, b := new(Recursive), new(Recursive);
	*a = Recursive{ 12, a };
	*b = Recursive{ 12, b };
	if !DeepEqual(a, b) {
		t.Error("DeepEqual(recursive same) = false, want true");
	}
}

func TestDeepEqualComplexStruct(t *testing.T) {
	m := make(map[float]float);
	stra, strb := "hello", "hello";
	a, b := new(Complex), new(Complex);
	*a = Complex{5, [3]*Complex{a, b, a}, &stra, m};
	*b = Complex{5, [3]*Complex{b, a, a}, &strb, m};
	if !DeepEqual(a, b) {
		t.Error("DeepEqual(complex same) = false, want true");
	}
}

func TestDeepEqualComplexStructInequality(t *testing.T) {
	m := make(map[float]float);
	stra, strb := "hello", "helloo";  // Difference is here
	a, b := new(Complex), new(Complex);
	*a = Complex{5, [3]*Complex{a, b, a}, &stra, m};
	*b = Complex{5, [3]*Complex{b, a, a}, &strb, m};
	if DeepEqual(a, b) {
		t.Error("DeepEqual(complex different) = true, want false");
	}
}


func check2ndField(x interface{}, offs uintptr, t *testing.T) {
	s := NewValue(x).(StructValue);
	name, ftype, tag, reflect_offset := s.Type().(StructType).Field(1);
	if uintptr(reflect_offset) != offs {
		t.Error("mismatched offsets in structure alignment:", reflect_offset, offs);
	}
}

// Check that structure alignment & offsets viewed through reflect agree with those
// from the compiler itself.
func TestAlignment(t *testing.T) {
	type T1inner struct {
		a int
	}
	type T1 struct {
		T1inner;
		f int;
	}
	type T2inner struct {
		a, b int
	}
	type T2 struct {
		T2inner;
		f int;
	}

	x := T1{T1inner{2}, 17};
	check2ndField(x, uintptr(unsafe.Pointer(&x.f)) - uintptr(unsafe.Pointer(&x)), t);

	x1 := T2{T2inner{2, 3}, 17};
	check2ndField(x1, uintptr(unsafe.Pointer(&x1.f)) - uintptr(unsafe.Pointer(&x1)), t);
}

type Nillable interface {
	IsNil() bool
}

func Nil(a interface{}, t *testing.T) {
	n := NewValue(a).(Nillable);
	if !n.IsNil() {
		t.Errorf("%v should be nil", a)
	}
}

func NotNil(a interface{}, t *testing.T) {
	n := NewValue(a).(Nillable);
	if n.IsNil() {
		t.Errorf("value of type %v should not be nil", NewValue(a).Type().String())
	}
}

func TestIsNil(t *testing.T) {
	// These do not implement IsNil
	doNotNil := []string{"int", "float32", "struct { a int }"};
	// These do implement IsNil
	doNil := []string{"*int", "interface{}", "map[string]int", "func() bool", "chan int", "[]string"};
	for i, ts := range doNotNil {
		ty := ParseTypeString("", ts);
		v := NewZeroValue(ty);
		if nilable, ok := v.(Nillable); ok {
			t.Errorf("%s is nilable; should not be", ts)
		}
	}

	for i, ts := range doNil {
		ty := ParseTypeString("", ts);
		v := NewZeroValue(ty);
		if nilable, ok := v.(Nillable); !ok {
			t.Errorf("%s %T is not nilable; should be", ts, v)
		}
	}
	// Check the implementations
	var pi *int;
	Nil(pi, t);
	pi = new(int);
	NotNil(pi, t);

	var si []int;
	Nil(si, t);
	si = make([]int, 10);
	NotNil(si, t);

	// TODO: map and chan don't work yet

	var ii interface {};
	Nil(ii, t);
	ii = pi;
	NotNil(ii, t);

	var fi func(t *testing.T);
	Nil(fi, t);
	fi = TestIsNil;
	NotNil(fi, t);
}

func TestInterfaceExtraction(t *testing.T) {
	var s struct {
		w io.Writer;
	}

	s.w = os.Stdout;
	v := Indirect(NewValue(&s)).(StructValue).Field(0).Interface();
	if v != s.w.(interface{}) {
		t.Errorf("Interface() on interface: ", v, s.w);
	}
}

func TestInterfaceEditing(t *testing.T) {
	// strings are bigger than one word,
	// so the interface conversion allocates
	// memory to hold a string and puts that
	// pointer in the interface.
	var i interface{} = "hello";

	// if i pass the interface value by value
	// to NewValue, i should get a fresh copy
	// of the value.
	v := NewValue(i);

	// and setting that copy to "bye" should
	// not change the value stored in i.
	v.(StringValue).Set("bye");
	if i.(string) != "hello" {
		t.Errorf(`Set("bye") changed i to %s`, i.(string));
	}

	// the same should be true of smaller items.
	i = 123;
	v = NewValue(i);
	v.(IntValue).Set(234);
	if i.(int) != 123 {
		t.Errorf("Set(234) changed i to %d", i.(int));
	}
}
