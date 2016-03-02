// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// Generate Windows callback assembly file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

const maxCallback = 2000

func genasm() {
	var buf bytes.Buffer

	buf.WriteString(`// generated by wincallback.go; run go generate

// runtime·callbackasm is called by external code to
// execute Go implemented callback function. It is not
// called from the start, instead runtime·compilecallback
// always returns address into runtime·callbackasm offset
// appropriately so different callbacks start with different
// CALL instruction in runtime·callbackasm. This determines
// which Go callback function is executed later on.
TEXT runtime·callbackasm(SB),7,$0
`)
	for i := 0; i < maxCallback; i++ {
		buf.WriteString("\tCALL\truntime·callbackasm1(SB)\n")
	}

	err := ioutil.WriteFile("zcallback_windows.s", buf.Bytes(), 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "wincallback: %s\n", err)
		os.Exit(2)
	}
}

func gengo() {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf(`// generated by wincallback.go; run go generate

package runtime

const cb_max = %d // maximum number of windows callbacks allowed
`, maxCallback))
	err := ioutil.WriteFile("zcallback_windows.go", buf.Bytes(), 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "wincallback: %s\n", err)
		os.Exit(2)
	}
}

func main() {
	genasm()
	gengo()
}
