// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"strings"
)

func adjustForeignAttributes(aa []Attribute) {
	for i, a := range aa {
		if a.Key == "" || a.Key[0] != 'x' {
			continue
		}
		switch a.Key {
		case "xlink:actuate", "xlink:arcrole", "xlink:href", "xlink:role", "xlink:show",
			"xlink:title", "xlink:type", "xml:base", "xml:lang", "xml:space", "xmlns:xlink":
			j := strings.Index(a.Key, ":")
			aa[i].Namespace = a.Key[:j]
			aa[i].Key = a.Key[j+1:]
		}
	}
}

// Section 12.2.5.5.
var breakout = map[string]bool{
	"b":          true,
	"big":        true,
	"blockquote": true,
	"body":       true,
	"br":         true,
	"center":     true,
	"code":       true,
	"dd":         true,
	"div":        true,
	"dl":         true,
	"dt":         true,
	"em":         true,
	"embed":      true,
	"font":       true,
	"h1":         true,
	"h2":         true,
	"h3":         true,
	"h4":         true,
	"h5":         true,
	"h6":         true,
	"head":       true,
	"hr":         true,
	"i":          true,
	"img":        true,
	"li":         true,
	"listing":    true,
	"menu":       true,
	"meta":       true,
	"nobr":       true,
	"ol":         true,
	"p":          true,
	"pre":        true,
	"ruby":       true,
	"s":          true,
	"small":      true,
	"span":       true,
	"strong":     true,
	"strike":     true,
	"sub":        true,
	"sup":        true,
	"table":      true,
	"tt":         true,
	"u":          true,
	"ul":         true,
	"var":        true,
}

// TODO: add look-up tables for MathML and SVG adjustments.
