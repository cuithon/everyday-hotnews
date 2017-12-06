// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// This program generates example_test.go.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/bits"
)

const header = `// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by go run make_examples.go. DO NOT EDIT.

package bits_test

import (
	"fmt"
	"math/bits"
)
`

func main() {
	w := bytes.NewBuffer([]byte(header))

	for _, e := range []struct {
		name string
		in   int
		out  [4]interface{}
		out2 [4]interface{}
	}{
		{
			name: "LeadingZeros",
			in:   1,
			out:  [4]interface{}{bits.LeadingZeros8(1), bits.LeadingZeros16(1), bits.LeadingZeros32(1), bits.LeadingZeros64(1)},
		},
		{
			name: "TrailingZeros",
			in:   14,
			out:  [4]interface{}{bits.TrailingZeros8(14), bits.TrailingZeros16(14), bits.TrailingZeros32(14), bits.TrailingZeros64(14)},
		},
		{
			name: "OnesCount",
			in:   14,
			out:  [4]interface{}{bits.OnesCount8(14), bits.OnesCount16(14), bits.OnesCount32(14), bits.OnesCount64(14)},
		},
		{
			name: "RotateLeft",
			in:   15,
			out:  [4]interface{}{bits.RotateLeft8(15, 2), bits.RotateLeft16(15, 2), bits.RotateLeft32(15, 2), bits.RotateLeft64(15, 2)},
			out2: [4]interface{}{bits.RotateLeft8(15, -2), bits.RotateLeft16(15, -2), bits.RotateLeft32(15, -2), bits.RotateLeft64(15, -2)},
		},
		{
			name: "Reverse",
			in:   19,
			out:  [4]interface{}{bits.Reverse8(19), bits.Reverse16(19), bits.Reverse32(19), bits.Reverse64(19)},
		},
		{
			name: "ReverseBytes",
			in:   15,
			out:  [4]interface{}{nil, bits.ReverseBytes16(15), bits.ReverseBytes32(15), bits.ReverseBytes64(15)},
		},
		{
			name: "Len",
			in:   8,
			out:  [4]interface{}{bits.Len8(8), bits.Len16(8), bits.Len32(8), bits.Len64(8)},
		},
	} {
		for i, size := range []int{8, 16, 32, 64} {
			if e.out[i] == nil {
				continue // function doesn't exist
			}
			f := fmt.Sprintf("%s%d", e.name, size)
			fmt.Fprintf(w, "\nfunc Example%s() {\n", f)
			switch e.name {
			case "RotateLeft", "Reverse", "ReverseBytes":
				fmt.Fprintf(w, "\tfmt.Printf(\"%%0%db\\n\", %d)\n", size, e.in)
				if e.name == "RotateLeft" {
					fmt.Fprintf(w, "\tfmt.Printf(\"%%0%db\\n\", bits.%s(%d, 2))\n", size, f, e.in)
					fmt.Fprintf(w, "\tfmt.Printf(\"%%0%db\\n\", bits.%s(%d, -2))\n", size, f, e.in)
				} else {
					fmt.Fprintf(w, "\tfmt.Printf(\"%%0%db\\n\", bits.%s(%d))\n", size, f, e.in)
				}
				fmt.Fprintf(w, "\t// Output:\n")
				fmt.Fprintf(w, "\t// %0*b\n", size, e.in)
				fmt.Fprintf(w, "\t// %0*b\n", size, e.out[i])
				if e.name == "RotateLeft" && e.out2[i] != nil {
					fmt.Fprintf(w, "\t// %0*b\n", size, e.out2[i])
				}
			default:
				fmt.Fprintf(w, "\tfmt.Printf(\"%s(%%0%db) = %%d\\n\", %d, bits.%s(%d))\n", f, size, e.in, f, e.in)
				fmt.Fprintf(w, "\t// Output:\n")
				fmt.Fprintf(w, "\t// %s(%0*b) = %d\n", f, size, e.in, e.out[i])
			}
			fmt.Fprintf(w, "}\n")
		}
	}

	if err := ioutil.WriteFile("example_test.go", w.Bytes(), 0666); err != nil {
		log.Fatal(err)
	}
}
