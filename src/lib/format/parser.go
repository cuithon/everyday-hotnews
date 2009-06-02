// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package format

import (
	"container/vector";
	"fmt";
	"format";
	"go/scanner";
	"go/token";
	"io";
	"os";
	"strconv";
	"strings";
)

// ----------------------------------------------------------------------------
// Error handling

// Error describes an individual error. The position Pos, if valid,
// indicates the format source position the error relates to. The
// error is specified with the Msg string.
//
type Error struct {
	Pos token.Position;
	Msg string;
}


func (e *Error) String() string {
	pos := "";
	if e.Pos.IsValid() {
		pos = fmt.Sprintf("%d:%d: ", e.Pos.Line, e.Pos.Column);
	}
	return pos + e.Msg;
}


// An ErrorList is a list of errors encountered during parsing.
type ErrorList []*Error


// ErrorList implements SortInterface and the os.Error interface.

func (p ErrorList) Len() int  { return len(p); }
func (p ErrorList) Swap(i, j int)  { p[i], p[j] = p[j], p[i]; }
func (p ErrorList) Less(i, j int) bool  { return p[i].Pos.Offset < p[j].Pos.Offset; }


func (p ErrorList) String() string {
	switch len(p) {
	case 0: return "unspecified error";
	case 1: return p[0].String();
	}
	return fmt.Sprintf("%s (and %d more errors)", p[0].String(), len(p) - 1);
}


// ----------------------------------------------------------------------------
// Parsing

type parser struct {
	errors vector.Vector;
	scanner scanner.Scanner;
	pos token.Position;  // token position
	tok token.Token;  // one token look-ahead
	lit []byte;  // token literal

	packs map [string] string;  // PackageName -> ImportPath
	rules map [string] expr;  // RuleName -> Expression
}


func (p *parser) next() {
	p.pos, p.tok, p.lit = p.scanner.Scan();
	switch p.tok {
	case token.CHAN, token.FUNC, token.INTERFACE, token.MAP, token.STRUCT:
		// Go keywords for composite types are type names
		// returned by reflect. Accept them as identifiers.
		p.tok = token.IDENT;  // p.lit is already set correctly
	}
}


func (p *parser) init(src []byte) {
	p.errors.Init(0);
	p.scanner.Init(src, p, 0);
	p.next();  // initializes pos, tok, lit
	p.packs = make(map [string] string);
	p.rules = make(map [string] expr);
}


// The parser implements scanner.Error.
func (p *parser) Error(pos token.Position, msg string) {
	// Don't collect errors that are on the same line as the previous error
	// in the hope to reduce the number of spurious errors due to incorrect
	// parser synchronization.
	if p.errors.Len() == 0 || p.errors.Last().(*Error).Pos.Line != pos.Line {
		p.errors.Push(&Error{pos, msg});
	}
}


func (p *parser) errorExpected(pos token.Position, msg string) {
	msg = "expected " + msg;
	if pos.Offset == p.pos.Offset {
		// the error happened at the current position;
		// make the error message more specific
		msg += ", found '" + p.tok.String() + "'";
		if p.tok.IsLiteral() {
			msg += " " + string(p.lit);
		}
	}
	p.Error(pos, msg);
}


func (p *parser) expect(tok token.Token) token.Position {
	pos := p.pos;
	if p.tok != tok {
		p.errorExpected(pos, "'" + tok.String() + "'");
	}
	p.next();  // make progress in any case
	return pos;
}


func (p *parser) parseIdentifier() string {
	name := string(p.lit);
	p.expect(token.IDENT);
	return name;
}


func (p *parser) parseTypeName() (string, bool) {
	pos := p.pos;
	name, isIdent := p.parseIdentifier(), true;
	if p.tok == token.PERIOD {
		// got a package name, lookup package
		if importPath, found := p.packs[name]; found {
			name = importPath;
		} else {
			p.Error(pos, "package not declared: " + name);
		}
		p.next();
		name, isIdent = name + "." + p.parseIdentifier(), false;
	}
	return name, isIdent;
}


// Parses a rule name and returns it. If the rule name is
// a package-qualified type name, the package name is resolved.
// The 2nd result value is true iff the rule name consists of a
// single identifier only (and thus could be a package name).
//
func (p *parser) parseRuleName() (string, bool) {
	name, isIdent := "", false;
	switch p.tok {
	case token.IDENT:
		name, isIdent = p.parseTypeName();
	case token.DEFAULT:
		name = "default";
		p.next();
	case token.QUO:
		name = "/";
		p.next();
	default:
		p.errorExpected(p.pos, "rule name");
		p.next();  // make progress in any case
	}
	return name, isIdent;
}


func (p *parser) parseString() string {
	s := "";
	if p.tok == token.STRING {
		var err os.Error;
		s, err = strconv.Unquote(string(p.lit));
		// Unquote may fail with an error, but only if the scanner found
		// an illegal string in the first place. In this case the error
		// has already been reported.
		p.next();
		return s;
	} else {
		p.expect(token.STRING);
	}
	return s;
}


func (p *parser) parseLiteral() literal {
	s := io.StringBytes(p.parseString());

	// A string literal may contain %-format specifiers. To simplify
	// and speed up printing of the literal, split it into segments
	// that start with "%" possibly followed by a last segment that
	// starts with some other character.
	var list vector.Vector;
	list.Init(0);
	i0 := 0;
	for i := 0; i < len(s); i++ {
		if s[i] == '%' && i+1 < len(s) {
			// the next segment starts with a % format
			if i0 < i {
				// the current segment is not empty, split it off
				list.Push(s[i0 : i]);
				i0 = i;
			}
			i++;  // skip %; let loop skip over char after %
		}
	}
	// the final segment may start with any character
	// (it is empty iff the string is empty)
	list.Push(s[i0 : len(s)]);

	// convert list into a literal
	lit := make(literal, list.Len());
	for i := 0; i < list.Len(); i++ {
		lit[i] = list.At(i).([]byte);
	}

	return lit;
}


func (p *parser) parseField() expr {
	var fname string;
	switch p.tok {
	case token.XOR:
		fname = "^";
		p.next();
	case token.MUL:
		fname = "*";
		p.next();
	case token.IDENT:
		fname = p.parseIdentifier();
	default:
		return nil;
	}

	var ruleName string;
	if p.tok == token.COLON {
		p.next();
		var _ bool;
		ruleName, _ = p.parseRuleName();
	}

	return &field{fname, ruleName};
}


func (p *parser) parseExpression() expr

func (p *parser) parseOperand() (x expr) {
	switch p.tok {
	case token.STRING:
		x = p.parseLiteral();

	case token.LPAREN:
		p.next();
		x = p.parseExpression();
		if p.tok == token.SHR {
			p.next();
			x = &group{x, p.parseExpression()};
		}
		p.expect(token.RPAREN);

	case token.LBRACK:
		p.next();
		x = &option{p.parseExpression()};
		p.expect(token.RBRACK);

	case token.LBRACE:
		p.next();
		x = p.parseExpression();
		var div expr;
		if p.tok == token.QUO {
			p.next();
			div = p.parseExpression();
		}
		x = &repetition{x, div};
		p.expect(token.RBRACE);

	default:
		x = p.parseField();  // may be nil
	}

	return x;
}


func (p *parser) parseSequence() expr {
	var list vector.Vector;
	list.Init(0);

	for x := p.parseOperand(); x != nil; x = p.parseOperand() {
		list.Push(x);
	}

	// no need for a sequence if list.Len() < 2
	switch list.Len() {
	case 0: return nil;
	case 1: return list.At(0).(expr);
	}

	// convert list into a sequence
	seq := make(sequence, list.Len());
	for i := 0; i < list.Len(); i++ {
		seq[i] = list.At(i).(expr);
	}
	return seq;
}


func (p *parser) parseExpression() expr {
	var list vector.Vector;
	list.Init(0);

	for {
		x := p.parseSequence();
		if x != nil {
			list.Push(x);
		}
		if p.tok != token.OR {
			break;
		}
		p.next();
	}

	// no need for an alternatives if list.Len() < 2
	switch list.Len() {
	case 0: return nil;
	case 1: return list.At(0).(expr);
	}

	// convert list into a alternatives
	alt := make(alternatives, list.Len());
	for i := 0; i < list.Len(); i++ {
		alt[i] = list.At(i).(expr);
	}
	return alt;
}


func (p *parser) parseFormat() {
	for p.tok != token.EOF {
		pos := p.pos;

		name, isIdent := p.parseRuleName();
		switch p.tok {
		case token.STRING:
			// package declaration
			importPath := p.parseString();

			// add package declaration
			if !isIdent {
				p.Error(pos, "illegal package name: " + name);
			} else if _, found := p.packs[name]; !found {
				p.packs[name] = importPath;
			} else {
				p.Error(pos, "package already declared: " + name);
			}

		case token.ASSIGN:
			// format rule
			p.next();
			x := p.parseExpression();

			// add rule
			if _, found := p.rules[name]; !found {
				p.rules[name] = x;
			} else {
				p.Error(pos, "format rule already declared: " + name);
			}

		default:
			p.errorExpected(p.pos, "package declaration or format rule");
			p.next();  // make progress in any case
		}

		if p.tok == token.SEMICOLON {
			p.next();
		} else {
			break;
		}
	}
	p.expect(token.EOF);
}


func remap(p *parser, name string) string {
	i := strings.Index(name, ".");
	if i >= 0 {
		packageName := name[0 : i];
		typeName := name[i : len(name)];
		// lookup package
		if importPath, found := p.packs[packageName]; found {
			name = importPath + "." + typeName;
		} else {
			var invalidPos token.Position;
			p.Error(invalidPos, "package not declared: " + packageName);
		}
	}
	return name;
}


// Parse parses a set of format productions from source src. Custom
// formatters may be provided via a map of formatter functions. If
// there are no errors, the result is a Format and the error is nil.
// Otherwise the format is nil and a non-empty ErrorList is returned.
//
func Parse(src []byte, fmap FormatterMap) (Format, os.Error) {
	// parse source
	var p parser;
	p.init(src);
	p.parseFormat();

	// add custom formatters, if any
	for name, form := range fmap {
		name = remap(&p, name);
		if t, found := p.rules[name]; !found {
			p.rules[name] = &custom{name, form};
		} else {
			var invalidPos token.Position;
			p.Error(invalidPos, "formatter already declared: " + name);
		}
	}

	// convert errors list, if any
	if p.errors.Len() > 0 {
		errors := make(ErrorList, p.errors.Len());
		for i := 0; i < p.errors.Len(); i++ {
			errors[i] = p.errors.At(i).(*Error);
		}
		return nil, errors;
	}

	return p.rules, nil;
}
