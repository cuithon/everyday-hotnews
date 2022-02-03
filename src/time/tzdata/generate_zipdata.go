// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

// This program generates zipdata.go from $GOROOT/lib/time/zoneinfo.zip.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// header is put at the start of the generated file.
// The string addition avoids this file (generate_zipdata.go) from
// matching the "generated file" regexp.
const header = `// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

` + `// Code generated by generate_zipdata. DO NOT EDIT.

// This file contains an embedded zip archive that contains time zone
// files compiled using the code and data maintained as part of the
// IANA Time Zone Database.
// The IANA asserts that the data is in the public domain.

// For more information, see
// https://www.iana.org/time-zones
// ftp://ftp.iana.org/tz/code/tz-link.html
// https://datatracker.ietf.org/doc/html/rfc6557

package tzdata

const zipdata = `

func main() {
	// We should be run in the $GOROOT/src/time/tzdata directory.
	data, err := os.ReadFile("../../../lib/time/zoneinfo.zip")
	if err != nil {
		die("cannot find zoneinfo.zip file: %v", err)
	}

	of, err := os.Create("zipdata.go")
	if err != nil {
		die("%v", err)
	}

	buf := bufio.NewWriter(of)
	buf.WriteString(header)

	ds := string(data)
	i := 0
	const chunk = 60
	for ; i+chunk < len(data); i += chunk {
		if i > 0 {
			buf.WriteRune('\t')
		}
		fmt.Fprintf(buf, "%s +\n", strconv.Quote(ds[i:i+chunk]))
	}
	fmt.Fprintf(buf, "\t%s\n", strconv.Quote(ds[i:]))

	if err := buf.Flush(); err != nil {
		die("error writing to zipdata.go: %v", err)
	}
	if err := of.Close(); err != nil {
		die("error closing zipdata.go: %v", err)
	}
}

func die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
