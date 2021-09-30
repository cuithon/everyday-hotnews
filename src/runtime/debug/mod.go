// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package debug

import (
	"bytes"
	"fmt"
	"strings"
)

// exported from runtime
func modinfo() string

// ReadBuildInfo returns the build information embedded
// in the running binary. The information is available only
// in binaries built with module support.
func ReadBuildInfo() (info *BuildInfo, ok bool) {
	return readBuildInfo(modinfo())
}

// BuildInfo represents the build information read from
// the running binary.
type BuildInfo struct {
	Path string    // The main package path
	Main Module    // The module containing the main package
	Deps []*Module // Module dependencies
}

// Module represents a module.
type Module struct {
	Path    string  // module path
	Version string  // module version
	Sum     string  // checksum
	Replace *Module // replaced by this module
}

func (bi *BuildInfo) MarshalText() ([]byte, error) {
	buf := &bytes.Buffer{}
	if bi.Path != "" {
		fmt.Fprintf(buf, "path\t%s\n", bi.Path)
	}
	var formatMod func(string, Module)
	formatMod = func(word string, m Module) {
		buf.WriteString(word)
		buf.WriteByte('\t')
		buf.WriteString(m.Path)
		mv := m.Version
		if mv == "" {
			mv = "(devel)"
		}
		buf.WriteByte('\t')
		buf.WriteString(mv)
		if m.Replace == nil {
			buf.WriteByte('\t')
			buf.WriteString(m.Sum)
		} else {
			buf.WriteByte('\n')
			formatMod("=>", *m.Replace)
		}
		buf.WriteByte('\n')
	}
	if bi.Main.Path != "" {
		formatMod("mod", bi.Main)
	}
	for _, dep := range bi.Deps {
		formatMod("dep", *dep)
	}

	return buf.Bytes(), nil
}

func readBuildInfo(data string) (*BuildInfo, bool) {
	if len(data) < 32 {
		return nil, false
	}
	data = data[16 : len(data)-16]

	const (
		pathLine = "path\t"
		modLine  = "mod\t"
		depLine  = "dep\t"
		repLine  = "=>\t"
	)

	readEntryFirstLine := func(elem []string) (Module, bool) {
		if len(elem) != 2 && len(elem) != 3 {
			return Module{}, false
		}
		sum := ""
		if len(elem) == 3 {
			sum = elem[2]
		}
		return Module{
			Path:    elem[0],
			Version: elem[1],
			Sum:     sum,
		}, true
	}

	var (
		info = &BuildInfo{}
		last *Module
		line string
		ok   bool
	)
	// Reverse of cmd/go/internal/modload.PackageBuildInfo
	for len(data) > 0 {
		line, data, ok = strings.Cut(data, "\n")
		if !ok {
			break
		}
		switch {
		case strings.HasPrefix(line, pathLine):
			elem := line[len(pathLine):]
			info.Path = elem
		case strings.HasPrefix(line, modLine):
			elem := strings.Split(line[len(modLine):], "\t")
			last = &info.Main
			*last, ok = readEntryFirstLine(elem)
			if !ok {
				return nil, false
			}
		case strings.HasPrefix(line, depLine):
			elem := strings.Split(line[len(depLine):], "\t")
			last = new(Module)
			info.Deps = append(info.Deps, last)
			*last, ok = readEntryFirstLine(elem)
			if !ok {
				return nil, false
			}
		case strings.HasPrefix(line, repLine):
			elem := strings.Split(line[len(repLine):], "\t")
			if len(elem) != 3 {
				return nil, false
			}
			if last == nil {
				return nil, false
			}
			last.Replace = &Module{
				Path:    elem[0],
				Version: elem[1],
				Sum:     elem[2],
			}
			last = nil
		}
	}
	return info, true
}
