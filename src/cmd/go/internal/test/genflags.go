// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"flag"
	exec "internal/execabs"
	"log"
	"os"
	"strings"
	"testing"
	"text/template"

	"cmd/go/internal/test/internal/genflags"
)

func main() {
	if err := regenerate(); err != nil {
		log.Fatal(err)
	}
}

func regenerate() error {
	vetAnalyzers, err := genflags.VetAnalyzers()
	if err != nil {
		return err
	}

	t := template.Must(template.New("fileTemplate").Parse(fileTemplate))
	tData := map[string][]string{
		"testFlags":    testFlags(),
		"vetAnalyzers": vetAnalyzers,
	}
	buf := bytes.NewBuffer(nil)
	if err := t.Execute(buf, tData); err != nil {
		return err
	}

	f, err := os.Create("flagdefs.go")
	if err != nil {
		return err
	}

	cmd := exec.Command("gofmt")
	cmd.Stdin = buf
	cmd.Stdout = f
	cmd.Stderr = os.Stderr
	cmdErr := cmd.Run()

	if err := f.Close(); err != nil {
		return err
	}
	if cmdErr != nil {
		os.Remove(f.Name())
		return cmdErr
	}

	return nil
}

func testFlags() []string {
	testing.Init()

	var names []string
	flag.VisitAll(func(f *flag.Flag) {
		if !strings.HasPrefix(f.Name, "test.") {
			return
		}
		name := strings.TrimPrefix(f.Name, "test.")

		switch name {
		case "testlogfile", "paniconexit0", "fuzzcachedir", "fuzzworker":
			// These flags are only for use by cmd/go.
		default:
			names = append(names, name)
		}
	})

	return names
}

const fileTemplate = `// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by genflags.go — DO NOT EDIT.

package test

// passFlagToTest contains the flags that should be forwarded to
// the test binary with the prefix "test.".
var passFlagToTest = map[string]bool {
{{- range .testFlags}}
	"{{.}}": true,
{{- end }}
}

var passAnalyzersToVet = map[string]bool {
{{- range .vetAnalyzers}}
	"{{.}}": true,
{{- end }}
}
`
