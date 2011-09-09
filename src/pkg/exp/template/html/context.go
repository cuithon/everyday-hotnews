// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"fmt"
)

// context describes the state an HTML parser must be in when it reaches the
// portion of HTML produced by evaluating a particular template node.
//
// The zero value of type context is the start context for a template that
// produces an HTML fragment as defined at
// http://www.w3.org/TR/html5/the-end.html#parsing-html-fragments
// where the context element is null.
type context struct {
	state   state
	delim   delim
	urlPart urlPart
	jsCtx   jsCtx
	element element
	errLine int
	errStr  string
}

// eq returns whether two contexts are equal.
func (c context) eq(d context) bool {
	return c.state == d.state &&
		c.delim == d.delim &&
		c.urlPart == d.urlPart &&
		c.jsCtx == d.jsCtx &&
		c.element == d.element &&
		c.errLine == d.errLine &&
		c.errStr == d.errStr
}

// state describes a high-level HTML parser state.
//
// It bounds the top of the element stack, and by extension the HTML insertion
// mode, but also contains state that does not correspond to anything in the
// HTML5 parsing algorithm because a single token production in the HTML
// grammar may contain embedded actions in a template. For instance, the quoted
// HTML attribute produced by
//     <div title="Hello {{.World}}">
// is a single token in HTML's grammar but in a template spans several nodes.
type state uint8

const (
	// stateText is parsed character data. An HTML parser is in
	// this state when its parse position is outside an HTML tag,
	// directive, comment, and special element body.
	stateText state = iota
	// stateTag occurs before an HTML attribute or the end of a tag.
	stateTag
	// stateComment occurs inside an <!-- HTML comment -->.
	stateComment
	// stateRCDATA occurs inside an RCDATA element (<textarea> or <title>)
	// as described at http://dev.w3.org/html5/spec/syntax.html#elements-0
	stateRCDATA
	// stateAttr occurs inside an HTML attribute whose content is text.
	stateAttr
	// stateURL occurs inside an HTML attribute whose content is a URL.
	stateURL
	// stateJS occurs inside an event handler or script element.
	stateJS
	// stateJSDqStr occurs inside a JavaScript double quoted string.
	stateJSDqStr
	// stateJSSqStr occurs inside a JavaScript single quoted string.
	stateJSSqStr
	// stateJSRegexp occurs inside a JavaScript regexp literal.
	stateJSRegexp
	// stateJSBlockCmt occurs inside a JavaScript /* block comment */.
	stateJSBlockCmt
	// stateJSLineCmt occurs inside a JavaScript // line comment.
	stateJSLineCmt
	// stateCSS occurs inside a <style> element or style attribute.
	stateCSS
	// stateCSSDqStr occurs inside a CSS double quoted string.
	stateCSSDqStr
	// stateCSSSqStr occurs inside a CSS single quoted string.
	stateCSSSqStr
	// stateCSSDqURL occurs inside a CSS double quoted url("...").
	stateCSSDqURL
	// stateCSSSqURL occurs inside a CSS single quoted url('...').
	stateCSSSqURL
	// stateCSSURL occurs inside a CSS unquoted url(...).
	stateCSSURL
	// stateCSSBlockCmt occurs inside a CSS /* block comment */.
	stateCSSBlockCmt
	// stateCSSLineCmt occurs inside a CSS // line comment.
	stateCSSLineCmt
	// stateError is an infectious error state outside any valid
	// HTML/CSS/JS construct.
	stateError
)

var stateNames = [...]string{
	stateText:        "stateText",
	stateTag:         "stateTag",
	stateComment:     "stateComment",
	stateRCDATA:      "stateRCDATA",
	stateAttr:        "stateAttr",
	stateURL:         "stateURL",
	stateJS:          "stateJS",
	stateJSDqStr:     "stateJSDqStr",
	stateJSSqStr:     "stateJSSqStr",
	stateJSRegexp:    "stateJSRegexp",
	stateJSBlockCmt:  "stateJSBlockCmt",
	stateJSLineCmt:   "stateJSLineCmt",
	stateCSS:         "stateCSS",
	stateCSSDqStr:    "stateCSSDqStr",
	stateCSSSqStr:    "stateCSSSqStr",
	stateCSSDqURL:    "stateCSSDqURL",
	stateCSSSqURL:    "stateCSSSqURL",
	stateCSSURL:      "stateCSSURL",
	stateCSSBlockCmt: "stateCSSBlockCmt",
	stateCSSLineCmt:  "stateCSSLineCmt",
	stateError:       "stateError",
}

func (s state) String() string {
	if int(s) < len(stateNames) {
		return stateNames[s]
	}
	return fmt.Sprintf("illegal state %d", s)
}

// delim is the delimiter that will end the current HTML attribute.
type delim uint8

const (
	// delimNone occurs outside any attribute.
	delimNone delim = iota
	// delimDoubleQuote occurs when a double quote (") closes the attribute.
	delimDoubleQuote
	// delimSingleQuote occurs when a single quote (') closes the attribute.
	delimSingleQuote
	// delimSpaceOrTagEnd occurs when a space or right angle bracket (>)
	// closes the attribute.
	delimSpaceOrTagEnd
)

var delimNames = [...]string{
	delimNone:          "delimNone",
	delimDoubleQuote:   "delimDoubleQuote",
	delimSingleQuote:   "delimSingleQuote",
	delimSpaceOrTagEnd: "delimSpaceOrTagEnd",
}

func (d delim) String() string {
	if int(d) < len(delimNames) {
		return delimNames[d]
	}
	return fmt.Sprintf("illegal delim %d", d)
}

// urlPart identifies a part in an RFC 3986 hierarchical URL to allow different
// encoding strategies.
type urlPart uint8

const (
	// urlPartNone occurs when not in a URL, or possibly at the start:
	// ^ in "^http://auth/path?k=v#frag".
	urlPartNone urlPart = iota
	// urlPartPreQuery occurs in the scheme, authority, or path; between the
	// ^s in "h^ttp://auth/path^?k=v#frag".
	urlPartPreQuery
	// urlPartQueryOrFrag occurs in the query portion between the ^s in
	// "http://auth/path?^k=v#frag^".
	urlPartQueryOrFrag
	// urlPartUnknown occurs due to joining of contexts both before and
	// after the query separator.
	urlPartUnknown
)

var urlPartNames = [...]string{
	urlPartNone:        "urlPartNone",
	urlPartPreQuery:    "urlPartPreQuery",
	urlPartQueryOrFrag: "urlPartQueryOrFrag",
	urlPartUnknown:     "urlPartUnknown",
}

func (u urlPart) String() string {
	if int(u) < len(urlPartNames) {
		return urlPartNames[u]
	}
	return fmt.Sprintf("illegal urlPart %d", u)
}

// jsCtx determines whether a '/' starts a regular expression literal or a
// division operator.
type jsCtx uint8

const (
	// jsCtxRegexp occurs where a '/' would start a regexp literal.
	jsCtxRegexp jsCtx = iota
	// jsCtxDivOp occurs where a '/' would start a division operator.
	jsCtxDivOp
)

func (c jsCtx) String() string {
	switch c {
	case jsCtxRegexp:
		return "jsCtxRegexp"
	case jsCtxDivOp:
		return "jsCtxDivOp"
	}
	return fmt.Sprintf("illegal jsCtx %d", c)
}

// element identifies the HTML element when inside a start tag or special body.
// Certain HTML element (for example <script> and <style>) have bodies that are
// treated differently from stateText so the element type is necessary to
// transition into the correct context at the end of a tag and to identify the
// end delimiter for the body.
type element uint8

const (
	// elementNone occurs outside a special tag or special element body.
	elementNone element = iota
	// elementScript corresponds to the raw text <script> element.
	elementScript
	// elementStyle corresponds to the raw text <style> element.
	elementStyle
	// elementTextarea corresponds to the RCDATA <textarea> element.
	elementTextarea
	// elementTitle corresponds to the RCDATA <title> element.
	elementTitle
)

var elementNames = [...]string{
	elementNone:     "elementNone",
	elementScript:   "elementScript",
	elementStyle:    "elementStyle",
	elementTextarea: "elementTextarea",
	elementTitle:    "elementTitle",
}

func (e element) String() string {
	if int(e) < len(elementNames) {
		return elementNames[e]
	}
	return fmt.Sprintf("illegal element %d", e)
}
