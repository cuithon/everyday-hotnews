// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package template

import (
	"bytes"
	"fmt"
	"http"
	"io"
	"os"
	"reflect"
	"strings"
	"unicode"
	"utf8"
)

// FuncMap is the type of the map defining the mapping from names to functions.
// Each function must have either a single return value, or two return values of
// which the second has type os.Error. If the second argument evaluates to non-nil
// during execution, execution terminates and Execute returns an error.
type FuncMap map[string]interface{}

var builtins = map[string]reflect.Value{
	"and":     reflect.ValueOf(and),
	"html":    reflect.ValueOf(HTMLEscaper),
	"index":   reflect.ValueOf(index),
	"js":      reflect.ValueOf(JSEscaper),
	"not":     reflect.ValueOf(not),
	"or":      reflect.ValueOf(or),
	"print":   reflect.ValueOf(fmt.Sprint),
	"printf":  reflect.ValueOf(fmt.Sprintf),
	"println": reflect.ValueOf(fmt.Sprintln),
	"url":     reflect.ValueOf(URLEscaper),
}

// addFuncs adds to values the functions in funcs, converting them to reflect.Values.
func addFuncs(values map[string]reflect.Value, funcMap FuncMap) {
	for name, fn := range funcMap {
		v := reflect.ValueOf(fn)
		if v.Kind() != reflect.Func {
			panic("value for " + name + " not a function")
		}
		if !goodFunc(v.Type()) {
			panic(fmt.Errorf("can't handle multiple results from method/function %q", name))
		}
		values[name] = v
	}
}

// goodFunc checks that the function or method has the right result signature.
func goodFunc(typ reflect.Type) bool {
	// We allow functions with 1 result or 2 results where the second is an os.Error.
	switch {
	case typ.NumOut() == 1:
		return true
	case typ.NumOut() == 2 && typ.Out(1) == osErrorType:
		return true
	}
	return false
}

// findFunction looks for a function in the template, set, and global map.
func findFunction(name string, tmpl *Template, set *Set) (reflect.Value, bool) {
	if tmpl != nil {
		if fn := tmpl.funcs[name]; fn.IsValid() {
			return fn, true
		}
	}
	if set != nil {
		if fn := set.funcs[name]; fn.IsValid() {
			return fn, true
		}
	}
	if fn := builtins[name]; fn.IsValid() {
		return fn, true
	}
	return reflect.Value{}, false
}

// Indexing.

// index returns the result of indexing its first argument by the following
// arguments.  Thus "index x 1 2 3" is, in Go syntax, x[1][2][3]. Each
// indexed item must be a map, slice, or array.
func index(item interface{}, indices ...interface{}) (interface{}, os.Error) {
	v := reflect.ValueOf(item)
	for _, i := range indices {
		index := reflect.ValueOf(i)
		var isNil bool
		if v, isNil = indirect(v); isNil {
			return nil, fmt.Errorf("index of nil pointer")
		}
		switch v.Kind() {
		case reflect.Array, reflect.Slice:
			var x int64
			switch index.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				x = index.Int()
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				x = int64(index.Uint())
			default:
				return nil, fmt.Errorf("cannot index slice/array with type %s", index.Type())
			}
			if x < 0 || x >= int64(v.Len()) {
				return nil, fmt.Errorf("index out of range: %d", x)
			}
			v = v.Index(int(x))
		case reflect.Map:
			if !index.Type().AssignableTo(v.Type().Key()) {
				return nil, fmt.Errorf("%s is not index type for %s", index.Type(), v.Type())
			}
			if x := v.MapIndex(index); x.IsValid() {
				v = x
			} else {
				v = reflect.Zero(v.Type().Key())
			}
		default:
			return nil, fmt.Errorf("can't index item of type %s", index.Type())
		}
	}
	return v.Interface(), nil
}

// Boolean logic.

func truth(a interface{}) bool {
	t, _ := isTrue(reflect.ValueOf(a))
	return t
}

// and computes the Boolean AND of its arguments, returning
// the first false argument it encounters, or the last argument.
func and(arg0 interface{}, args ...interface{}) interface{} {
	if !truth(arg0) {
		return arg0
	}
	for i := range args {
		arg0 = args[i]
		if !truth(arg0) {
			break
		}
	}
	return arg0
}

// or computes the Boolean OR of its arguments, returning
// the first true argument it encounters, or the last argument.
func or(arg0 interface{}, args ...interface{}) interface{} {
	if truth(arg0) {
		return arg0
	}
	for i := range args {
		arg0 = args[i]
		if truth(arg0) {
			break
		}
	}
	return arg0
}

// not returns the Boolean negation of its argument.
func not(arg interface{}) (truth bool) {
	truth, _ = isTrue(reflect.ValueOf(arg))
	return !truth
}

// HTML escaping.

var (
	htmlQuot = []byte("&#34;") // shorter than "&quot;"
	htmlApos = []byte("&#39;") // shorter than "&apos;"
	htmlAmp  = []byte("&amp;")
	htmlLt   = []byte("&lt;")
	htmlGt   = []byte("&gt;")
)

// HTMLEscape writes to w the escaped HTML equivalent of the plain text data b.
func HTMLEscape(w io.Writer, b []byte) {
	last := 0
	for i, c := range b {
		var html []byte
		switch c {
		case '"':
			html = htmlQuot
		case '\'':
			html = htmlApos
		case '&':
			html = htmlAmp
		case '<':
			html = htmlLt
		case '>':
			html = htmlGt
		default:
			continue
		}
		w.Write(b[last:i])
		w.Write(html)
		last = i + 1
	}
	w.Write(b[last:])
}

// HTMLEscapeString returns the escaped HTML equivalent of the plain text data s.
func HTMLEscapeString(s string) string {
	// Avoid allocation if we can.
	if strings.IndexAny(s, `'"&<>`) < 0 {
		return s
	}
	var b bytes.Buffer
	HTMLEscape(&b, []byte(s))
	return b.String()
}

// HTMLEscaper returns the escaped HTML equivalent of the textual
// representation of its arguments.
func HTMLEscaper(args ...interface{}) string {
	ok := false
	var s string
	if len(args) == 1 {
		s, ok = args[0].(string)
	}
	if !ok {
		s = fmt.Sprint(args...)
	}
	return HTMLEscapeString(s)
}

// JavaScript escaping.

var (
	jsLowUni = []byte(`\u00`)
	hex      = []byte("0123456789ABCDEF")

	jsBackslash = []byte(`\\`)
	jsApos      = []byte(`\'`)
	jsQuot      = []byte(`\"`)
	jsLt        = []byte(`\x3C`)
	jsGt        = []byte(`\x3E`)
)

// JSEscape writes to w the escaped JavaScript equivalent of the plain text data b.
func JSEscape(w io.Writer, b []byte) {
	last := 0
	for i := 0; i < len(b); i++ {
		c := b[i]

		if !jsIsSpecial(int(c)) {
			// fast path: nothing to do
			continue
		}
		w.Write(b[last:i])

		if c < utf8.RuneSelf {
			// Quotes, slashes and angle brackets get quoted.
			// Control characters get written as \u00XX.
			switch c {
			case '\\':
				w.Write(jsBackslash)
			case '\'':
				w.Write(jsApos)
			case '"':
				w.Write(jsQuot)
			case '<':
				w.Write(jsLt)
			case '>':
				w.Write(jsGt)
			default:
				w.Write(jsLowUni)
				t, b := c>>4, c&0x0f
				w.Write(hex[t : t+1])
				w.Write(hex[b : b+1])
			}
		} else {
			// Unicode rune.
			rune, size := utf8.DecodeRune(b[i:])
			if unicode.IsPrint(rune) {
				w.Write(b[i : i+size])
			} else {
				// TODO(dsymonds): Do this without fmt?
				fmt.Fprintf(w, "\\u%04X", rune)
			}
			i += size - 1
		}
		last = i + 1
	}
	w.Write(b[last:])
}

// JSEscapeString returns the escaped JavaScript equivalent of the plain text data s.
func JSEscapeString(s string) string {
	// Avoid allocation if we can.
	if strings.IndexFunc(s, jsIsSpecial) < 0 {
		return s
	}
	var b bytes.Buffer
	JSEscape(&b, []byte(s))
	return b.String()
}

func jsIsSpecial(rune int) bool {
	switch rune {
	case '\\', '\'', '"', '<', '>':
		return true
	}
	return rune < ' ' || utf8.RuneSelf <= rune
}

// JSEscaper returns the escaped JavaScript equivalent of the textual
// representation of its arguments.
func JSEscaper(args ...interface{}) string {
	ok := false
	var s string
	if len(args) == 1 {
		s, ok = args[0].(string)
	}
	if !ok {
		s = fmt.Sprint(args...)
	}
	return JSEscapeString(s)
}

// URLEscaper returns the escaped value of the textual representation of its
// arguments in a form suitable for embedding in a URL.
func URLEscaper(args ...interface{}) string {
	s, ok := "", false
	if len(args) == 1 {
		s, ok = args[0].(string)
	}
	if !ok {
		s = fmt.Sprint(args...)
	}
	return http.URLEscape(s)
}
