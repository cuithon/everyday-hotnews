// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Package template implements data-driven templates for generating textual
	output such as HTML.

	Templates are executed by applying them to a data structure.
	Annotations in the template refer to elements of the data
	structure (typically a field of a struct or a key in a map)
	to control execution and derive values to be displayed.
	The template walks the structure as it executes and the
	"cursor" @ represents the value at the current location
	in the structure.

	Data items may be values or pointers; the interface hides the
	indirection.

	In the following, 'field' is one of several things, according to the data.

		- The name of a field of a struct (result = data.field),
		- The value stored in a map under that key (result = data[field]), or
		- The result of invoking a niladic single-valued method with that name
		  (result = data.field())

	Major constructs ({} are the default delimiters for template actions;
	[] are the notation in this comment for optional elements):

		{# comment }

	A one-line comment.

		{.section field} XXX [ {.or} YYY ] {.end}

	Set @ to the value of the field.  It may be an explicit @
	to stay at the same point in the data. If the field is nil
	or empty, execute YYY; otherwise execute XXX.

		{.repeated section field} XXX [ {.alternates with} ZZZ ] [ {.or} YYY ] {.end}

	Like .section, but field must be an array or slice.  XXX
	is executed for each element.  If the array is nil or empty,
	YYY is executed instead.  If the {.alternates with} marker
	is present, ZZZ is executed between iterations of XXX.

		{field}
		{field1 field2 ...}
		{field|formatter}
		{field1 field2...|formatter}
		{field|formatter1|formatter2}

	Insert the value of the fields into the output. Each field is
	first looked for in the cursor, as in .section and .repeated.
	If it is not found, the search continues in outer sections
	until the top level is reached.

	If the field value is a pointer, leading asterisks indicate
	that the value to be inserted should be evaluated through the
	pointer.  For example, if x.p is of type *int, {x.p} will
	insert the value of the pointer but {*x.p} will insert the
	value of the underlying integer.  If the value is nil or not a
	pointer, asterisks have no effect.

	If a formatter is specified, it must be named in the formatter
	map passed to the template set up routines or in the default
	set ("html","str","") and is used to process the data for
	output.  The formatter function has signature
		func(wr io.Writer, formatter string, data ...interface{})
	where wr is the destination for output, data holds the field
	values at the instantiation, and formatter is its name at
	the invocation site.  The default formatter just concatenates
	the string representations of the fields.

	Multiple formatters separated by the pipeline character | are
	executed sequentially, with each formatter receiving the bytes
	emitted by the one to its left.

	As well as field names, one may use literals with Go syntax.
	Integer, floating-point, and string literals are supported.
	Raw strings may not span newlines.

	The delimiter strings get their default value, "{" and "}", from
	JSON-template.  They may be set to any non-empty, space-free
	string using the SetDelims method.  Their value can be printed
	in the output using {.meta-left} and {.meta-right}.
*/
package template

import (
	"bytes"
	"container/vector"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"
	"utf8"
)

// Errors returned during parsing and execution.  Users may extract the information and reformat
// if they desire.
type Error struct {
	Line int
	Msg  string
}

func (e *Error) String() string { return fmt.Sprintf("line %d: %s", e.Line, e.Msg) }

// Most of the literals are aces.
var lbrace = []byte{'{'}
var rbrace = []byte{'}'}
var space = []byte{' '}
var tab = []byte{'\t'}

// The various types of "tokens", which are plain text or (usually) brace-delimited descriptors
const (
	tokAlternates = iota
	tokComment
	tokEnd
	tokLiteral
	tokOr
	tokRepeated
	tokSection
	tokText
	tokVariable
)

// FormatterMap is the type describing the mapping from formatter
// names to the functions that implement them.
type FormatterMap map[string]func(io.Writer, string, ...interface{})

// Built-in formatters.
var builtins = FormatterMap{
	"html": HTMLFormatter,
	"str":  StringFormatter,
	"":     StringFormatter,
}

// The parsed state of a template is a vector of xxxElement structs.
// Sections have line numbers so errors can be reported better during execution.

// Plain text.
type textElement struct {
	text []byte
}

// A literal such as .meta-left or .meta-right
type literalElement struct {
	text []byte
}

// A variable invocation to be evaluated
type variableElement struct {
	linenum int
	args    []interface{} // The fields and literals in the invocation.
	fmts    []string      // Names of formatters to apply. len(fmts) > 0
}

// A variableElement arg to be evaluated as a field name
type fieldName string

// A .section block, possibly with a .or
type sectionElement struct {
	linenum int    // of .section itself
	field   string // cursor field for this block
	start   int    // first element
	or      int    // first element of .or block
	end     int    // one beyond last element
}

// A .repeated block, possibly with a .or and a .alternates
type repeatedElement struct {
	sectionElement     // It has the same structure...
	altstart       int // ... except for alternates
	altend         int
}

// Template is the type that represents a template definition.
// It is unchanged after parsing.
type Template struct {
	fmap FormatterMap // formatters for variables
	// Used during parsing:
	ldelim, rdelim []byte // delimiters; default {}
	buf            []byte // input text to process
	p              int    // position in buf
	linenum        int    // position in input
	// Parsed results:
	elems *vector.Vector
}

// Internal state for executing a Template.  As we evaluate the struct,
// the data item descends into the fields associated with sections, etc.
// Parent is used to walk upwards to find variables higher in the tree.
type state struct {
	parent *state          // parent in hierarchy
	data   reflect.Value   // the driver data for this section etc.
	wr     io.Writer       // where to send output
	buf    [2]bytes.Buffer // alternating buffers used when chaining formatters
}

func (parent *state) clone(data reflect.Value) *state {
	return &state{parent: parent, data: data, wr: parent.wr}
}

// New creates a new template with the specified formatter map (which
// may be nil) to define auxiliary functions for formatting variables.
func New(fmap FormatterMap) *Template {
	t := new(Template)
	t.fmap = fmap
	t.ldelim = lbrace
	t.rdelim = rbrace
	t.elems = new(vector.Vector)
	return t
}

// Report error and stop executing.  The line number must be provided explicitly.
func (t *Template) execError(st *state, line int, err string, args ...interface{}) {
	panic(&Error{line, fmt.Sprintf(err, args...)})
}

// Report error, panic to terminate parsing.
// The line number comes from the template state.
func (t *Template) parseError(err string, args ...interface{}) {
	panic(&Error{t.linenum, fmt.Sprintf(err, args...)})
}

// Is this an exported - upper case - name?
func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

// -- Lexical analysis

// Is c a white space character?
func white(c uint8) bool { return c == ' ' || c == '\t' || c == '\r' || c == '\n' }

// Safely, does s[n:n+len(t)] == t?
func equal(s []byte, n int, t []byte) bool {
	b := s[n:]
	if len(t) > len(b) { // not enough space left for a match.
		return false
	}
	for i, c := range t {
		if c != b[i] {
			return false
		}
	}
	return true
}

// isQuote returns true if c is a string- or character-delimiting quote character.
func isQuote(c byte) bool {
	return c == '"' || c == '`' || c == '\''
}

// endQuote returns the end quote index for the quoted string that
// starts at n, or -1 if no matching end quote is found before the end
// of the line.
func endQuote(s []byte, n int) int {
	quote := s[n]
	for n++; n < len(s); n++ {
		switch s[n] {
		case '\\':
			if quote == '"' || quote == '\'' {
				n++
			}
		case '\n':
			return -1
		case quote:
			return n
		}
	}
	return -1
}

// nextItem returns the next item from the input buffer.  If the returned
// item is empty, we are at EOF.  The item will be either a
// delimited string or a non-empty string between delimited
// strings. Tokens stop at (but include, if plain text) a newline.
// Action tokens on a line by themselves drop any space on
// either side, up to and including the newline.
func (t *Template) nextItem() []byte {
	startOfLine := t.p == 0 || t.buf[t.p-1] == '\n'
	start := t.p
	var i int
	newline := func() {
		t.linenum++
		i++
	}
	// Leading white space up to but not including newline
	for i = start; i < len(t.buf); i++ {
		if t.buf[i] == '\n' || !white(t.buf[i]) {
			break
		}
	}
	leadingSpace := i > start
	// What's left is nothing, newline, delimited string, or plain text
	switch {
	case i == len(t.buf):
		// EOF; nothing to do
	case t.buf[i] == '\n':
		newline()
	case equal(t.buf, i, t.ldelim):
		left := i         // Start of left delimiter.
		right := -1       // Will be (immediately after) right delimiter.
		haveText := false // Delimiters contain text.
		i += len(t.ldelim)
		// Find the end of the action.
		for ; i < len(t.buf); i++ {
			if t.buf[i] == '\n' {
				break
			}
			if isQuote(t.buf[i]) {
				i = endQuote(t.buf, i)
				if i == -1 {
					t.parseError("unmatched quote")
					return nil
				}
				continue
			}
			if equal(t.buf, i, t.rdelim) {
				i += len(t.rdelim)
				right = i
				break
			}
			haveText = true
		}
		if right < 0 {
			t.parseError("unmatched opening delimiter")
			return nil
		}
		// Is this a special action (starts with '.' or '#') and the only thing on the line?
		if startOfLine && haveText {
			firstChar := t.buf[left+len(t.ldelim)]
			if firstChar == '.' || firstChar == '#' {
				// It's special and the first thing on the line. Is it the last?
				for j := right; j < len(t.buf) && white(t.buf[j]); j++ {
					if t.buf[j] == '\n' {
						// Yes it is. Drop the surrounding space and return the {.foo}
						t.linenum++
						t.p = j + 1
						return t.buf[left:right]
					}
				}
			}
		}
		// No it's not. If there's leading space, return that.
		if leadingSpace {
			// not trimming space: return leading white space if there is some.
			t.p = left
			return t.buf[start:left]
		}
		// Return the word, leave the trailing space.
		start = left
		break
	default:
		for ; i < len(t.buf); i++ {
			if t.buf[i] == '\n' {
				newline()
				break
			}
			if equal(t.buf, i, t.ldelim) {
				break
			}
		}
	}
	item := t.buf[start:i]
	t.p = i
	return item
}

// Turn a byte array into a white-space-split array of strings,
// taking into account quoted strings.
func words(buf []byte) []string {
	s := make([]string, 0, 5)
	for i := 0; i < len(buf); {
		// One word per loop
		for i < len(buf) && white(buf[i]) {
			i++
		}
		if i == len(buf) {
			break
		}
		// Got a word
		start := i
		if isQuote(buf[i]) {
			i = endQuote(buf, i)
			if i < 0 {
				i = len(buf)
			} else {
				i++
			}
		} else {
			for i < len(buf) && !white(buf[i]) {
				i++
			}
		}
		s = append(s, string(buf[start:i]))
	}
	return s
}

// Analyze an item and return its token type and, if it's an action item, an array of
// its constituent words.
func (t *Template) analyze(item []byte) (tok int, w []string) {
	// item is known to be non-empty
	if !equal(item, 0, t.ldelim) { // doesn't start with left delimiter
		tok = tokText
		return
	}
	if !equal(item, len(item)-len(t.rdelim), t.rdelim) { // doesn't end with right delimiter
		t.parseError("internal error: unmatched opening delimiter") // lexing should prevent this
		return
	}
	if len(item) <= len(t.ldelim)+len(t.rdelim) { // no contents
		t.parseError("empty directive")
		return
	}
	// Comment
	if item[len(t.ldelim)] == '#' {
		tok = tokComment
		return
	}
	// Split into words
	w = words(item[len(t.ldelim) : len(item)-len(t.rdelim)]) // drop final delimiter
	if len(w) == 0 {
		t.parseError("empty directive")
		return
	}
	first := w[0]
	if first[0] != '.' {
		tok = tokVariable
		return
	}
	if len(first) > 1 && first[1] >= '0' && first[1] <= '9' {
		// Must be a float.
		tok = tokVariable
		return
	}
	switch first {
	case ".meta-left", ".meta-right", ".space", ".tab":
		tok = tokLiteral
		return
	case ".or":
		tok = tokOr
		return
	case ".end":
		tok = tokEnd
		return
	case ".section":
		if len(w) != 2 {
			t.parseError("incorrect fields for .section: %s", item)
			return
		}
		tok = tokSection
		return
	case ".repeated":
		if len(w) != 3 || w[1] != "section" {
			t.parseError("incorrect fields for .repeated: %s", item)
			return
		}
		tok = tokRepeated
		return
	case ".alternates":
		if len(w) != 2 || w[1] != "with" {
			t.parseError("incorrect fields for .alternates: %s", item)
			return
		}
		tok = tokAlternates
		return
	}
	t.parseError("bad directive: %s", item)
	return
}

// formatter returns the Formatter with the given name in the Template, or nil if none exists.
func (t *Template) formatter(name string) func(io.Writer, string, ...interface{}) {
	if t.fmap != nil {
		if fn := t.fmap[name]; fn != nil {
			return fn
		}
	}
	return builtins[name]
}

// -- Parsing

// Allocate a new variable-evaluation element.
func (t *Template) newVariable(words []string) *variableElement {
	// After the final space-separated argument, formatters may be specified separated
	// by pipe symbols, for example: {a b c|d|e}

	// Until we learn otherwise, formatters contains a single name: "", the default formatter.
	formatters := []string{""}
	lastWord := words[len(words)-1]
	bar := strings.IndexRune(lastWord, '|')
	if bar >= 0 {
		words[len(words)-1] = lastWord[0:bar]
		formatters = strings.Split(lastWord[bar+1:], "|", -1)
	}

	args := make([]interface{}, len(words))

	// Build argument list, processing any literals
	for i, word := range words {
		var lerr os.Error
		switch word[0] {
		case '"', '`', '\'':
			v, err := strconv.Unquote(word)
			if err == nil && word[0] == '\'' {
				args[i] = []int(v)[0]
			} else {
				args[i], lerr = v, err
			}

		case '.', '+', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			v, err := strconv.Btoi64(word, 0)
			if err == nil {
				args[i] = v
			} else {
				v, err := strconv.Atof64(word)
				args[i], lerr = v, err
			}

		default:
			args[i] = fieldName(word)
		}
		if lerr != nil {
			t.parseError("invalid literal: %q: %s", word, lerr)
		}
	}

	// We could remember the function address here and avoid the lookup later,
	// but it's more dynamic to let the user change the map contents underfoot.
	// We do require the name to be present, though.

	// Is it in user-supplied map?
	for _, f := range formatters {
		if t.formatter(f) == nil {
			t.parseError("unknown formatter: %q", f)
		}
	}

	return &variableElement{t.linenum, args, formatters}
}

// Grab the next item.  If it's simple, just append it to the template.
// Otherwise return its details.
func (t *Template) parseSimple(item []byte) (done bool, tok int, w []string) {
	tok, w = t.analyze(item)
	done = true // assume for simplicity
	switch tok {
	case tokComment:
		return
	case tokText:
		t.elems.Push(&textElement{item})
		return
	case tokLiteral:
		switch w[0] {
		case ".meta-left":
			t.elems.Push(&literalElement{t.ldelim})
		case ".meta-right":
			t.elems.Push(&literalElement{t.rdelim})
		case ".space":
			t.elems.Push(&literalElement{space})
		case ".tab":
			t.elems.Push(&literalElement{tab})
		default:
			t.parseError("internal error: unknown literal: %s", w[0])
		}
		return
	case tokVariable:
		t.elems.Push(t.newVariable(w))
		return
	}
	return false, tok, w
}

// parseRepeated and parseSection are mutually recursive

func (t *Template) parseRepeated(words []string) *repeatedElement {
	r := new(repeatedElement)
	t.elems.Push(r)
	r.linenum = t.linenum
	r.field = words[2]
	// Scan section, collecting true and false (.or) blocks.
	r.start = t.elems.Len()
	r.or = -1
	r.altstart = -1
	r.altend = -1
Loop:
	for {
		item := t.nextItem()
		if len(item) == 0 {
			t.parseError("missing .end for .repeated section")
			break
		}
		done, tok, w := t.parseSimple(item)
		if done {
			continue
		}
		switch tok {
		case tokEnd:
			break Loop
		case tokOr:
			if r.or >= 0 {
				t.parseError("extra .or in .repeated section")
				break Loop
			}
			r.altend = t.elems.Len()
			r.or = t.elems.Len()
		case tokSection:
			t.parseSection(w)
		case tokRepeated:
			t.parseRepeated(w)
		case tokAlternates:
			if r.altstart >= 0 {
				t.parseError("extra .alternates in .repeated section")
				break Loop
			}
			if r.or >= 0 {
				t.parseError(".alternates inside .or block in .repeated section")
				break Loop
			}
			r.altstart = t.elems.Len()
		default:
			t.parseError("internal error: unknown repeated section item: %s", item)
			break Loop
		}
	}
	if r.altend < 0 {
		r.altend = t.elems.Len()
	}
	r.end = t.elems.Len()
	return r
}

func (t *Template) parseSection(words []string) *sectionElement {
	s := new(sectionElement)
	t.elems.Push(s)
	s.linenum = t.linenum
	s.field = words[1]
	// Scan section, collecting true and false (.or) blocks.
	s.start = t.elems.Len()
	s.or = -1
Loop:
	for {
		item := t.nextItem()
		if len(item) == 0 {
			t.parseError("missing .end for .section")
			break
		}
		done, tok, w := t.parseSimple(item)
		if done {
			continue
		}
		switch tok {
		case tokEnd:
			break Loop
		case tokOr:
			if s.or >= 0 {
				t.parseError("extra .or in .section")
				break Loop
			}
			s.or = t.elems.Len()
		case tokSection:
			t.parseSection(w)
		case tokRepeated:
			t.parseRepeated(w)
		case tokAlternates:
			t.parseError(".alternates not in .repeated")
		default:
			t.parseError("internal error: unknown section item: %s", item)
		}
	}
	s.end = t.elems.Len()
	return s
}

func (t *Template) parse() {
	for {
		item := t.nextItem()
		if len(item) == 0 {
			break
		}
		done, tok, w := t.parseSimple(item)
		if done {
			continue
		}
		switch tok {
		case tokOr, tokEnd, tokAlternates:
			t.parseError("unexpected %s", w[0])
		case tokSection:
			t.parseSection(w)
		case tokRepeated:
			t.parseRepeated(w)
		default:
			t.parseError("internal error: bad directive in parse: %s", item)
		}
	}
}

// -- Execution

// Evaluate interfaces and pointers looking for a value that can look up the name, via a
// struct field, method, or map key, and return the result of the lookup.
func (t *Template) lookup(st *state, v reflect.Value, name string) reflect.Value {
	for v.IsValid() {
		typ := v.Type()
		if n := v.Type().NumMethod(); n > 0 {
			for i := 0; i < n; i++ {
				m := typ.Method(i)
				mtyp := m.Type
				if m.Name == name && mtyp.NumIn() == 1 && mtyp.NumOut() == 1 {
					if !isExported(name) {
						t.execError(st, t.linenum, "name not exported: %s in type %s", name, st.data.Type())
					}
					return v.Method(i).Call(nil)[0]
				}
			}
		}
		switch av := v; av.Kind() {
		case reflect.Ptr:
			v = av.Elem()
		case reflect.Interface:
			v = av.Elem()
		case reflect.Struct:
			if !isExported(name) {
				t.execError(st, t.linenum, "name not exported: %s in type %s", name, st.data.Type())
			}
			return av.FieldByName(name)
		case reflect.Map:
			if v := av.MapIndex(reflect.ValueOf(name)); v.IsValid() {
				return v
			}
			return reflect.Zero(typ.Elem())
		default:
			return reflect.Value{}
		}
	}
	return v
}

// indirectPtr returns the item numLevels levels of indirection below the value.
// It is forgiving: if the value is not a pointer, it returns it rather than giving
// an error.  If the pointer is nil, it is returned as is.
func indirectPtr(v reflect.Value, numLevels int) reflect.Value {
	for i := numLevels; v.IsValid() && i > 0; i++ {
		if p := v; p.Kind() == reflect.Ptr {
			if p.IsNil() {
				return v
			}
			v = p.Elem()
		} else {
			break
		}
	}
	return v
}

// Walk v through pointers and interfaces, extracting the elements within.
func indirect(v reflect.Value) reflect.Value {
loop:
	for v.IsValid() {
		switch av := v; av.Kind() {
		case reflect.Ptr:
			v = av.Elem()
		case reflect.Interface:
			v = av.Elem()
		default:
			break loop
		}
	}
	return v
}

// If the data for this template is a struct, find the named variable.
// Names of the form a.b.c are walked down the data tree.
// The special name "@" (the "cursor") denotes the current data.
// The value coming in (st.data) might need indirecting to reach
// a struct while the return value is not indirected - that is,
// it represents the actual named field. Leading stars indicate
// levels of indirection to be applied to the value.
func (t *Template) findVar(st *state, s string) reflect.Value {
	data := st.data
	flattenedName := strings.TrimLeft(s, "*")
	numStars := len(s) - len(flattenedName)
	s = flattenedName
	if s == "@" {
		return indirectPtr(data, numStars)
	}
	for _, elem := range strings.Split(s, ".", -1) {
		// Look up field; data must be a struct or map.
		data = t.lookup(st, data, elem)
		if !data.IsValid() {
			return reflect.Value{}
		}
	}
	return indirectPtr(data, numStars)
}

// Is there no data to look at?
func empty(v reflect.Value) bool {
	v = indirect(v)
	if !v.IsValid() {
		return true
	}
	switch v.Kind() {
	case reflect.Bool:
		return v.Bool() == false
	case reflect.String:
		return v.String() == ""
	case reflect.Struct:
		return false
	case reflect.Map:
		return false
	case reflect.Array:
		return v.Len() == 0
	case reflect.Slice:
		return v.Len() == 0
	}
	return false
}

// Look up a variable or method, up through the parent if necessary.
func (t *Template) varValue(name string, st *state) reflect.Value {
	field := t.findVar(st, name)
	if !field.IsValid() {
		if st.parent == nil {
			t.execError(st, t.linenum, "name not found: %s in type %s", name, st.data.Type())
		}
		return t.varValue(name, st.parent)
	}
	return field
}

func (t *Template) format(wr io.Writer, fmt string, val []interface{}, v *variableElement, st *state) {
	fn := t.formatter(fmt)
	if fn == nil {
		t.execError(st, v.linenum, "missing formatter %s for variable", fmt)
	}
	fn(wr, fmt, val...)
}

// Evaluate a variable, looking up through the parent if necessary.
// If it has a formatter attached ({var|formatter}) run that too.
func (t *Template) writeVariable(v *variableElement, st *state) {
	// Resolve field names
	val := make([]interface{}, len(v.args))
	for i, arg := range v.args {
		if name, ok := arg.(fieldName); ok {
			val[i] = t.varValue(string(name), st).Interface()
		} else {
			val[i] = arg
		}
	}
	for i, fmt := range v.fmts[:len(v.fmts)-1] {
		b := &st.buf[i&1]
		b.Reset()
		t.format(b, fmt, val, v, st)
		val = val[0:1]
		val[0] = b.Bytes()
	}
	t.format(st.wr, v.fmts[len(v.fmts)-1], val, v, st)
}

// Execute element i.  Return next index to execute.
func (t *Template) executeElement(i int, st *state) int {
	switch elem := t.elems.At(i).(type) {
	case *textElement:
		st.wr.Write(elem.text)
		return i + 1
	case *literalElement:
		st.wr.Write(elem.text)
		return i + 1
	case *variableElement:
		t.writeVariable(elem, st)
		return i + 1
	case *sectionElement:
		t.executeSection(elem, st)
		return elem.end
	case *repeatedElement:
		t.executeRepeated(elem, st)
		return elem.end
	}
	e := t.elems.At(i)
	t.execError(st, 0, "internal error: bad directive in execute: %v %T\n", reflect.ValueOf(e).Interface(), e)
	return 0
}

// Execute the template.
func (t *Template) execute(start, end int, st *state) {
	for i := start; i < end; {
		i = t.executeElement(i, st)
	}
}

// Execute a .section
func (t *Template) executeSection(s *sectionElement, st *state) {
	// Find driver data for this section.  It must be in the current struct.
	field := t.varValue(s.field, st)
	if !field.IsValid() {
		t.execError(st, s.linenum, ".section: cannot find field %s in %s", s.field, st.data.Type())
	}
	st = st.clone(field)
	start, end := s.start, s.or
	if !empty(field) {
		// Execute the normal block.
		if end < 0 {
			end = s.end
		}
	} else {
		// Execute the .or block.  If it's missing, do nothing.
		start, end = s.or, s.end
		if start < 0 {
			return
		}
	}
	for i := start; i < end; {
		i = t.executeElement(i, st)
	}
}

// Return the result of calling the Iter method on v, or nil.
func iter(v reflect.Value) reflect.Value {
	for j := 0; j < v.Type().NumMethod(); j++ {
		mth := v.Type().Method(j)
		fv := v.Method(j)
		ft := fv.Type()
		// TODO(rsc): NumIn() should return 0 here, because ft is from a curried FuncValue.
		if mth.Name != "Iter" || ft.NumIn() != 1 || ft.NumOut() != 1 {
			continue
		}
		ct := ft.Out(0)
		if ct.Kind() != reflect.Chan ||
			ct.ChanDir()&reflect.RecvDir == 0 {
			continue
		}
		return fv.Call(nil)[0]
	}
	return reflect.Value{}
}

// Execute a .repeated section
func (t *Template) executeRepeated(r *repeatedElement, st *state) {
	// Find driver data for this section.  It must be in the current struct.
	field := t.varValue(r.field, st)
	if !field.IsValid() {
		t.execError(st, r.linenum, ".repeated: cannot find field %s in %s", r.field, st.data.Type())
	}
	field = indirect(field)

	start, end := r.start, r.or
	if end < 0 {
		end = r.end
	}
	if r.altstart >= 0 {
		end = r.altstart
	}
	first := true

	// Code common to all the loops.
	loopBody := func(newst *state) {
		// .alternates between elements
		if !first && r.altstart >= 0 {
			for i := r.altstart; i < r.altend; {
				i = t.executeElement(i, newst)
			}
		}
		first = false
		for i := start; i < end; {
			i = t.executeElement(i, newst)
		}
	}

	if array := field; array.Kind() == reflect.Array || array.Kind() == reflect.Slice {
		for j := 0; j < array.Len(); j++ {
			loopBody(st.clone(array.Index(j)))
		}
	} else if m := field; m.Kind() == reflect.Map {
		for _, key := range m.MapKeys() {
			loopBody(st.clone(m.MapIndex(key)))
		}
	} else if ch := iter(field); ch.IsValid() {
		for {
			e, ok := ch.Recv()
			if !ok {
				break
			}
			loopBody(st.clone(e))
		}
	} else {
		t.execError(st, r.linenum, ".repeated: cannot repeat %s (type %s)",
			r.field, field.Type())
	}

	if first {
		// Empty. Execute the .or block, once.  If it's missing, do nothing.
		start, end := r.or, r.end
		if start >= 0 {
			newst := st.clone(field)
			for i := start; i < end; {
				i = t.executeElement(i, newst)
			}
		}
		return
	}
}

// A valid delimiter must contain no white space and be non-empty.
func validDelim(d []byte) bool {
	if len(d) == 0 {
		return false
	}
	for _, c := range d {
		if white(c) {
			return false
		}
	}
	return true
}

// checkError is a deferred function to turn a panic with type *Error into a plain error return.
// Other panics are unexpected and so are re-enabled.
func checkError(error *os.Error) {
	if v := recover(); v != nil {
		if e, ok := v.(*Error); ok {
			*error = e
		} else {
			// runtime errors should crash
			panic(v)
		}
	}
}

// -- Public interface

// Parse initializes a Template by parsing its definition.  The string
// s contains the template text.  If any errors occur, Parse returns
// the error.
func (t *Template) Parse(s string) (err os.Error) {
	if t.elems == nil {
		return &Error{1, "template not allocated with New"}
	}
	if !validDelim(t.ldelim) || !validDelim(t.rdelim) {
		return &Error{1, fmt.Sprintf("bad delimiter strings %q %q", t.ldelim, t.rdelim)}
	}
	defer checkError(&err)
	t.buf = []byte(s)
	t.p = 0
	t.linenum = 1
	t.parse()
	return nil
}

// ParseFile is like Parse but reads the template definition from the
// named file.
func (t *Template) ParseFile(filename string) (err os.Error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return t.Parse(string(b))
}

// Execute applies a parsed template to the specified data object,
// generating output to wr.
func (t *Template) Execute(wr io.Writer, data interface{}) (err os.Error) {
	// Extract the driver data.
	val := reflect.ValueOf(data)
	defer checkError(&err)
	t.p = 0
	t.execute(0, t.elems.Len(), &state{parent: nil, data: val, wr: wr})
	return nil
}

// SetDelims sets the left and right delimiters for operations in the
// template.  They are validated during parsing.  They could be
// validated here but it's better to keep the routine simple.  The
// delimiters are very rarely invalid and Parse has the necessary
// error-handling interface already.
func (t *Template) SetDelims(left, right string) {
	t.ldelim = []byte(left)
	t.rdelim = []byte(right)
}

// Parse creates a Template with default parameters (such as {} for
// metacharacters).  The string s contains the template text while
// the formatter map fmap, which may be nil, defines auxiliary functions
// for formatting variables.  The template is returned. If any errors
// occur, err will be non-nil.
func Parse(s string, fmap FormatterMap) (t *Template, err os.Error) {
	t = New(fmap)
	err = t.Parse(s)
	if err != nil {
		t = nil
	}
	return
}

// ParseFile is a wrapper function that creates a Template with default
// parameters (such as {} for metacharacters).  The filename identifies
// a file containing the template text, while the formatter map fmap, which
// may be nil, defines auxiliary functions for formatting variables.
// The template is returned. If any errors occur, err will be non-nil.
func ParseFile(filename string, fmap FormatterMap) (t *Template, err os.Error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return Parse(string(b), fmap)
}

// MustParse is like Parse but panics if the template cannot be parsed.
func MustParse(s string, fmap FormatterMap) *Template {
	t, err := Parse(s, fmap)
	if err != nil {
		panic("template.MustParse error: " + err.String())
	}
	return t
}

// MustParseFile is like ParseFile but panics if the file cannot be read
// or the template cannot be parsed.
func MustParseFile(filename string, fmap FormatterMap) *Template {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("template.MustParseFile error: " + err.String())
	}
	return MustParse(string(b), fmap)
}
