// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime_test

import (
	"bytes"
	"runtime"
	"testing"
)

// TestGCInfo tests that various objects in heap, data and bss receive correct GC pointer type info.
func TestGCInfo(t *testing.T) {
	verifyGCInfo(t, "bss ScalarPtr", &bssScalarPtr, infoScalarPtr)
	verifyGCInfo(t, "bss PtrScalar", &bssPtrScalar, infoPtrScalar)
	verifyGCInfo(t, "bss BigStruct", &bssBigStruct, infoBigStruct())
	verifyGCInfo(t, "bss string", &bssString, infoString)
	verifyGCInfo(t, "bss slice", &bssSlice, infoSlice)
	verifyGCInfo(t, "bss eface", &bssEface, infoEface)
	verifyGCInfo(t, "bss iface", &bssIface, infoIface)

	verifyGCInfo(t, "data ScalarPtr", &dataScalarPtr, infoScalarPtr)
	verifyGCInfo(t, "data PtrScalar", &dataPtrScalar, infoPtrScalar)
	verifyGCInfo(t, "data BigStruct", &dataBigStruct, infoBigStruct())
	verifyGCInfo(t, "data string", &dataString, infoString)
	verifyGCInfo(t, "data slice", &dataSlice, infoSlice)
	verifyGCInfo(t, "data eface", &dataEface, infoEface)
	verifyGCInfo(t, "data iface", &dataIface, infoIface)

	verifyGCInfo(t, "stack ScalarPtr", new(ScalarPtr), nonStackInfo(infoScalarPtr))
	verifyGCInfo(t, "stack PtrScalar", new(PtrScalar), nonStackInfo(infoPtrScalar))
	verifyGCInfo(t, "stack BigStruct", new(BigStruct), nonStackInfo(infoBigStruct()))
	verifyGCInfo(t, "stack string", new(string), nonStackInfo(infoString))
	verifyGCInfo(t, "stack slice", new([]string), nonStackInfo(infoSlice))
	verifyGCInfo(t, "stack eface", new(interface{}), nonStackInfo(infoEface))
	verifyGCInfo(t, "stack iface", new(Iface), nonStackInfo(infoIface))

	for i := 0; i < 10; i++ {
		verifyGCInfo(t, "heap ScalarPtr", escape(new(ScalarPtr)), infoScalarPtr)
		verifyGCInfo(t, "heap PtrScalar", escape(new(PtrScalar)), infoPtrScalar)
		verifyGCInfo(t, "heap BigStruct", escape(new(BigStruct)), infoBigStruct())
		verifyGCInfo(t, "heap string", escape(new(string)), infoString)
		verifyGCInfo(t, "heap eface", escape(new(interface{})), infoEface)
		verifyGCInfo(t, "heap iface", escape(new(Iface)), infoIface)
	}

}

func verifyGCInfo(t *testing.T, name string, p interface{}, mask0 []byte) {
	mask := runtime.GCMask(p)
	if len(mask) > len(mask0) {
		mask0 = append(mask0, typeDead)
		mask = mask[:len(mask0)]
	}
	if bytes.Compare(mask, mask0) != 0 {
		t.Errorf("bad GC program for %v:\nwant %+v\ngot  %+v", name, mask0, mask)
		return
	}
}

func nonStackInfo(mask []byte) []byte {
	// typeDead is replaced with typeScalar everywhere except stacks.
	mask1 := make([]byte, len(mask))
	for i, v := range mask {
		if v == typeDead {
			v = typeScalar
		}
		mask1[i] = v
	}
	return mask1
}

var gcinfoSink interface{}

func escape(p interface{}) interface{} {
	gcinfoSink = p
	return p
}

const (
	typeDead = iota
	typeScalar
	typePointer
)

const (
	BitsString = iota // unused
	BitsSlice         // unused
	BitsIface
	BitsEface
)

type ScalarPtr struct {
	q int
	w *int
	e int
	r *int
	t int
	y *int
}

var infoScalarPtr = []byte{typeScalar, typePointer, typeScalar, typePointer, typeScalar, typePointer}

type PtrScalar struct {
	q *int
	w int
	e *int
	r int
	t *int
	y int
}

var infoPtrScalar = []byte{typePointer, typeScalar, typePointer, typeScalar, typePointer, typeScalar}

type BigStruct struct {
	q *int
	w byte
	e [17]byte
	r []byte
	t int
	y uint16
	u uint64
	i string
}

func infoBigStruct() []byte {
	switch runtime.GOARCH {
	case "386", "arm":
		return []byte{
			typePointer,                                                // q *int
			typeScalar, typeScalar, typeScalar, typeScalar, typeScalar, // w byte; e [17]byte
			typePointer, typeScalar, typeScalar, // r []byte
			typeScalar, typeScalar, typeScalar, typeScalar, // t int; y uint16; u uint64
			typePointer, typeScalar, // i string
		}
	case "arm64", "amd64", "ppc64", "ppc64le":
		return []byte{
			typePointer,                        // q *int
			typeScalar, typeScalar, typeScalar, // w byte; e [17]byte
			typePointer, typeScalar, typeScalar, // r []byte
			typeScalar, typeScalar, typeScalar, // t int; y uint16; u uint64
			typePointer, typeScalar, // i string
		}
	case "amd64p32":
		return []byte{
			typePointer,                                                // q *int
			typeScalar, typeScalar, typeScalar, typeScalar, typeScalar, // w byte; e [17]byte
			typePointer, typeScalar, typeScalar, // r []byte
			typeScalar, typeScalar, typeScalar, typeScalar, typeScalar, // t int; y uint16; u uint64
			typePointer, typeScalar, // i string
		}
	default:
		panic("unknown arch")
	}
}

type Iface interface {
	f()
}

type IfaceImpl int

func (IfaceImpl) f() {
}

var (
	// BSS
	bssScalarPtr ScalarPtr
	bssPtrScalar PtrScalar
	bssBigStruct BigStruct
	bssString    string
	bssSlice     []string
	bssEface     interface{}
	bssIface     Iface

	// DATA
	dataScalarPtr             = ScalarPtr{q: 1}
	dataPtrScalar             = PtrScalar{w: 1}
	dataBigStruct             = BigStruct{w: 1}
	dataString                = "foo"
	dataSlice                 = []string{"foo"}
	dataEface     interface{} = 42
	dataIface     Iface       = IfaceImpl(42)

	infoString = []byte{typePointer, typeScalar}
	infoSlice  = []byte{typePointer, typeScalar, typeScalar}
	infoEface  = []byte{typePointer, typePointer}
	infoIface  = []byte{typePointer, typePointer}
)
