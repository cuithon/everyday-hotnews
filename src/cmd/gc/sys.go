// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package PACKAGE

// emitted by compiler, not referred to by go programs

func	mal(int32) *any;
func	throwindex();
func	throwreturn();
func	panicl(int32);

func	printbool(bool);
func	printfloat(float64);
func	printint(int64);
func	printstring(string);
func	printpointer(*any);
func	printinter(any);
func	printarray(any);
func	printnl();
func	printsp();

func	catstring(string, string) string;
func	cmpstring(string, string) int;
func	slicestring(string, int, int) string;
func	indexstring(string, int) byte;
func	intstring(int64) string;
func	byteastring(*byte, int) string;
func	arraystring([]byte) string;

func	ifaceT2I(sigi *byte, sigt *byte, elem any) (ret any);
func	ifaceI2T(sigt *byte, iface any) (ret any);
func	ifaceI2T2(sigt *byte, iface any) (ret any, ok bool);
func	ifaceI2I(sigi *byte, iface any) (ret any);
func	ifaceI2I2(sigi *byte, iface any) (ret any, ok bool);
func	ifaceeq(i1 any, i2 any) (ret bool);

func	newmap(keysize int, valsize int,
			keyalg int, valalg int,
			hint int) (hmap map[any]any);
func	mapaccess1(hmap map[any]any, key any) (val any);
func	mapaccess2(hmap map[any]any, key any) (val any, pres bool);
func	mapassign1(hmap map[any]any, key any, val any);
func	mapassign2(hmap map[any]any, key any, val any, pres bool);
func	mapiterinit(hmap map[any]any, hiter *any);
func	mapiternext(hiter *any);
func	mapiter1(hiter *any) (key any);
func	mapiter2(hiter *any) (key any, val any);

func	newchan(elemsize int, elemalg int, hint int) (hchan chan any);
func	chanrecv1(hchan chan any) (elem any);
func	chanrecv2(hchan chan any) (elem any, pres bool);
func	chanrecv3(hchan chan any, elem *any) (pres bool);
func	chansend1(hchan chan any, elem any);
func	chansend2(hchan chan any, elem any) (pres bool);

func	newselect(size int) (sel *byte);
func	selectsend(sel *byte, hchan chan any, elem any) (selected bool);
func	selectrecv(sel *byte, hchan chan any, elem *any) (selected bool);
func	selectdefault(sel *byte) (selected bool);
func	selectgo(sel *byte);

func	newarray(nel int, cap int, width int) (ary []any);
func	arraysliced(old []any, lb int, hb int, width int) (ary []any);
func	arrayslices(old *any, nel int, lb int, hb int, width int) (ary []any);
func	arrays2d(old *any, nel int) (ary []any);

// used by go programs

export func	Breakpoint();

export func	Reflect(i interface { }) (uint64, string, bool);
export func	Unreflect(uint64, string, bool) (ret interface { });

export var	Args []string;
export var	Envs []string;

export func	Frexp(float64) (float64, int);		// break fp into exp,fract
export func	Ldexp(float64, int) float64;		// make fp from exp,fract
export func	Modf(float64) (float64, float64);	// break fp into double.double
export func	IsInf(float64, int) bool;		// test for infinity
export func	IsNaN(float64) bool;			// test for not-a-number
export func	Inf(int) float64;			// return signed Inf
export func	NaN() float64;				// return a NaN
export func	Float32bits(float32) uint32;		// raw bits
export func	Float64bits(float64) uint64;		// raw bits
export func	Float32frombits(uint32) float32;	// raw bits
export func	Float64frombits(uint64) float64;	// raw bits

export func	Gosched();
export func	Goexit();

export func	BytesToRune(*byte, int, int) (int, int);	// convert bytes to runes
export func	StringToRune(string, int) (int, int);	// convert bytes to runes

export func	Exit(int);

export func	Caller(n int) (pc uint64, file string, line int, ok bool);

export func	SemAcquire(sema *int32);
export func	SemRelease(sema *int32);
