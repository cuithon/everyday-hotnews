// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
)

var gooses []string

func main() {
	data, err := os.ReadFile("../../go/build/syslist.go")
	if err != nil {
		log.Fatal(err)
	}
	const goosPrefix = `var knownOS = map[string]bool{`
	inGOOS := false
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, goosPrefix) {
			inGOOS = true
		} else if inGOOS && strings.HasPrefix(line, "}") {
			break
		} else if inGOOS {
			goos := strings.Fields(line)[0]
			goos = strings.TrimPrefix(goos, `"`)
			goos = strings.TrimSuffix(goos, `":`)
			gooses = append(gooses, goos)
		}
	}

	for _, target := range gooses {
		if target == "nacl" {
			continue
		}
		var tags []string
		if target == "linux" {
			tags = append(tags, "!android") // must explicitly exclude android for linux
		}
		if target == "solaris" {
			tags = append(tags, "!illumos") // must explicitly exclude illumos for solaris
		}
		if target == "darwin" {
			tags = append(tags, "!ios") // must explicitly exclude ios for darwin
		}
		tags = append(tags, target) // must explicitly include target for bootstrapping purposes
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "// Code generated by gengoos.go using 'go generate'. DO NOT EDIT.\n\n")
		fmt.Fprintf(&buf, "//go:build %s\n\n", strings.Join(tags, " && "))
		fmt.Fprintf(&buf, "package goos\n\n")
		fmt.Fprintf(&buf, "const GOOS = `%s`\n\n", target)
		for _, goos := range gooses {
			value := 0
			if goos == target {
				value = 1
			}
			fmt.Fprintf(&buf, "const Is%s = %d\n", strings.Title(goos), value)
		}
		err := os.WriteFile("zgoos_"+target+".go", buf.Bytes(), 0666)
		if err != nil {
			log.Fatal(err)
		}
	}
}
