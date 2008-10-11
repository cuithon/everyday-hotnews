// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Regular expression library.

package regexp

import (
	"os";
	"vector";
)


export var ErrUnimplemented = os.NewError("unimplemented");
export var ErrInternal = os.NewError("internal error");
export var ErrUnmatchedLpar = os.NewError("unmatched '('");
export var ErrUnmatchedRpar = os.NewError("unmatched ')'");
export var ErrUnmatchedLbkt = os.NewError("unmatched '['");
export var ErrUnmatchedRbkt = os.NewError("unmatched ']'");
export var ErrBadRange = os.NewError("bad range in character class");
export var ErrExtraneousBackslash = os.NewError("extraneous backslash");
export var ErrEmpty = os.NewError("empty subexpression or alternation");
export var ErrBadClosure = os.NewError("repeated closure (**, ++, etc.)");
export var ErrBareClosure = os.NewError("closure applies to nothing");
export var ErrBadBackslash = os.NewError("illegal backslash escape");

// An instruction executed by the NFA
type Inst interface {
	Type()	int;	// the type of this instruction: CHAR, ANY, etc.
	Next()	Inst;	// the instruction to execute after this one
	SetNext(i Inst);
	Index()	int;
	SetIndex(i int);
	Print();
}

type RE struct {
	expr	string;	// the original expression
	ch	*chan<- *RE;	// reply channel when we're done
	error	*os.Error;	// compile- or run-time error; nil if OK
	inst	*vector.Vector;
	start	Inst;
}

const (
	START	// beginning of program: indexer to start
		= iota;
	END;		// end of program: success
	BOT;		// '^' beginning of text
	EOT;		// '$' end of text
	CHAR;	// 'a' regular character
	CHARCLASS;	// [a-z] character class
	ANY;		// '.' any character
	BRA;		// '(' parenthesized expression
	EBRA;	// ')'; end of '(' parenthesized expression
	ALT;		// '|' alternation
	NOP;		// do nothing; makes it easy to link without patching
)

// --- START start of program
type Start struct {
	next	Inst;
	index	int;
}

func (start *Start) Type() int { return START }
func (start *Start) Next() Inst { return start.next }
func (start *Start) SetNext(i Inst) { start.next = i }
func (start *Start) Index() int { return start.index }
func (start *Start) SetIndex(i int) { start.index = i }
func (start *Start) Print() { print("start") }

// --- END end of program
type End struct {
	next	Inst;
	index	int;
}

func (end *End) Type() int { return END }
func (end *End) Next() Inst { return end.next }
func (end *End) SetNext(i Inst) { end.next = i }
func (end *End) Index() int { return end.index }
func (end *End) SetIndex(i int) { end.index = i }
func (end *End) Print() { print("end") }

// --- BOT beginning of text
type Bot struct {
	next	Inst;
	index	int;
}

func (bot *Bot) Type() int { return BOT }
func (bot *Bot) Next() Inst { return bot.next }
func (bot *Bot) SetNext(i Inst) { bot.next = i }
func (bot *Bot) Index() int { return bot.index }
func (bot *Bot) SetIndex(i int) { bot.index = i }
func (bot *Bot) Print() { print("bot") }

// --- EOT end of text
type Eot struct {
	next	Inst;
	index	int;
}

func (eot *Eot) Type() int { return EOT }
func (eot *Eot) Next() Inst { return eot.next }
func (eot *Eot) SetNext(i Inst) { eot.next = i }
func (eot *Eot) Index() int { return eot.index }
func (eot *Eot) SetIndex(i int) { eot.index = i }
func (eot *Eot) Print() { print("eot") }

// --- CHAR a regular character
type Char struct {
	next	Inst;
	char	int;
	index	int;
}

func (char *Char) Type() int { return CHAR }
func (char *Char) Next() Inst { return char.next }
func (char *Char) SetNext(i Inst) { char.next = i }
func (char *Char) Index() int { return char.index }
func (char *Char) SetIndex(i int) { char.index = i }
func (char *Char) Print() { print("char ", string(char.char)) }

func NewChar(char int) *Char {
	c := new(Char);
	c.char = char;
	return c;
}

// --- CHARCLASS [a-z]

type CClassChar int;	// BUG: Shouldn't be necessary but 6g won't put ints into vectors

type CharClass struct {
	next	Inst;
	index	int;
	char	int;
	negate	bool;	// is character class negated? ([^a-z])
	// Vector of CClassChar, stored pairwise: [a-z] is (a,z); x is (x,x):
	ranges	*vector.Vector;
}

func (cclass *CharClass) Type() int { return CHAR }
func (cclass *CharClass) Next() Inst { return cclass.next }
func (cclass *CharClass) SetNext(i Inst) { cclass.next = i }
func (cclass *CharClass) Index() int { return cclass.index }
func (cclass *CharClass) SetIndex(i int) { cclass.index = i }
func (cclass *CharClass) Print() {
	print("charclass");
	if cclass.negate {
		print(" (negated)");
	}
	for i := 0; i < cclass.ranges.Len(); i += 2 {
		l := cclass.ranges.At(i).(CClassChar);
		r := cclass.ranges.At(i+1).(CClassChar);
		if l == r {
			print(" [", string(l), "]");
		} else {
			print(" [", string(l), "-", string(r), "]");
		}
	}
}

func (cclass *CharClass) AddRange(a, b CClassChar) {
	// range is a through b inclusive
	cclass.ranges.Append(a);
	cclass.ranges.Append(b);
}

func NewCharClass() *CharClass {
	c := new(CharClass);
	c.ranges = vector.New();
	return c;
}

// --- ANY any character
type Any struct {
	next	Inst;
	index	int;
}

func (any *Any) Type() int { return ANY }
func (any *Any) Next() Inst { return any.next }
func (any *Any) SetNext(i Inst) { any.next = i }
func (any *Any) Index() int { return any.index }
func (any *Any) SetIndex(i int) { any.index = i }
func (any *Any) Print() { print("any") }

// --- BRA parenthesized expression
type Bra struct {
	next	Inst;
	index	int;
	n	int;	// subexpression number
}

func (bra *Bra) Type() int { return BRA }
func (bra *Bra) Next() Inst { return bra.next }
func (bra *Bra) SetNext(i Inst) { bra.next = i }
func (bra *Bra) Index() int { return bra.index }
func (bra *Bra) SetIndex(i int) { bra.index = i }
func (bra *Bra) Print() { print("bra"); }

// --- EBRA end of parenthesized expression
type Ebra struct {
	next	Inst;
	index	int;
	n	int;	// subexpression number
}

func (ebra *Ebra) Type() int { return BRA }
func (ebra *Ebra) Next() Inst { return ebra.next }
func (ebra *Ebra) SetNext(i Inst) { ebra.next = i }
func (ebra *Ebra) Index() int { return ebra.index }
func (ebra *Ebra) SetIndex(i int) { ebra.index = i }
func (ebra *Ebra) Print() { print("ebra ", ebra.n); }

// --- ALT alternation
type Alt struct {
	next	Inst;
	index	int;
	left	Inst;	// other branch
}

func (alt *Alt) Type() int { return ALT }
func (alt *Alt) Next() Inst { return alt.next }
func (alt *Alt) SetNext(i Inst) { alt.next = i }
func (alt *Alt) Index() int { return alt.index }
func (alt *Alt) SetIndex(i int) { alt.index = i }
func (alt *Alt) Print() { print("alt(", alt.left.Index(), ")"); }

// --- NOP no operation
type Nop struct {
	next	Inst;
	index	int;
}

func (nop *Nop) Type() int { return NOP }
func (nop *Nop) Next() Inst { return nop.next }
func (nop *Nop) SetNext(i Inst) { nop.next = i }
func (nop *Nop) Index() int { return nop.index }
func (nop *Nop) SetIndex(i int) { nop.index = i }
func (nop *Nop) Print() { print("nop") }

// report error and exit compiling/executing goroutine
func (re *RE) Error(err *os.Error) {
	re.error = err;
	re.ch <- re;
	sys.goexit();
}

func (re *RE) Add(i Inst) Inst {
	i.SetIndex(re.inst.Len());
	re.inst.Append(i);
	return i;
}

type Parser struct {
	re	*RE;
	nbra	int;	// number of brackets in expression, for subexpressions
	nlpar	int;	// number of unclosed lpars
	pos	int;
	ch	int;
}

const EOF = -1

func (p *Parser) c() int {
	return p.ch;
}

func (p *Parser) nextc() int {
	if p.pos >= len(p.re.expr) {
		p.ch = EOF
	} else {
		// TODO: stringotorune should take a string*
		c, w := sys.stringtorune(p.re.expr, p.pos);
		p.ch = c;
		p.pos += w;
	}
	return p.ch;
}

func NewParser(re *RE) *Parser {
	parser := new(Parser);
	parser.re = re;
	parser.nextc();	// load p.ch
	return parser;
}

/*

Grammar:
	regexp:
		concatenation { '|' concatenation }
	concatenation:
		{ closure }
	closure:
		term [ '*' | '+' | '?' ]
	term:
		'^'
		'$'
		'.'
		character
		'[' character-ranges ']'
		'(' regexp ')'

*/

func (p *Parser) Regexp() (start, end Inst)

var NULL Inst
type BUGinter interface{}

// same as i == NULL.  TODO: remove when 6g lets me do i == NULL
func isNULL(i Inst) bool {
	return sys.BUG_intereq(i.(BUGinter), NULL.(BUGinter))
}

// same as i == j.  TODO: remove when 6g lets me do i == j
func isEQ(i,j Inst) bool {
	return sys.BUG_intereq(i.(BUGinter), j.(BUGinter))
}

func special(c int) bool {
	s := `\.+*?()|[]`;
	for i := 0; i < len(s); i++ {
		if c == int(s[i]) {
			return true
		}
	}
	return false
}

func specialcclass(c int) bool {
	s := `\-[]`;
	for i := 0; i < len(s); i++ {
		if c == int(s[i]) {
			return true
		}
	}
	return false
}

func (p *Parser) CharClass() Inst {
	cc := NewCharClass();
	p.re.Add(cc);
	if p.c() == '^' {
		cc.negate = true;
		p.nextc();
	}
	left := -1;
	for {
		switch c := p.c(); c {
		case ']', EOF:
			if left >= 0 {
				p.re.Error(ErrBadRange);
			}
			return cc;
		case '-':	// do this before backslash processing
			p.re.Error(ErrBadRange);
		case '\\':
			c = p.nextc();
			switch {
			case c == EOF:
				p.re.Error(ErrExtraneousBackslash);
			case c == 'n':
				c = '\n';
			case specialcclass(c):
				// c is as delivered
			default:
				p.re.Error(ErrBadBackslash);
			}
			fallthrough;
		default:
			p.nextc();
			switch {
			case left < 0:	// first of pair
				if p.c() == '-' {	// range
					p.nextc();
					left = c;
				} else {	// single char
					cc.AddRange(c, c);
				}
			case left <= c:	// second of pair
				cc.AddRange(left, c);
				left = -1;
			default:
				p.re.Error(ErrBadRange);
			}
		}
	}
	return NULL
}

func (p *Parser) Term() (start, end Inst) {
	switch c := p.c(); c {
	case '|', EOF:
		return NULL, NULL;
	case '*', '+', '|':
		p.re.Error(ErrBareClosure);
	case ')':
		if p.nlpar == 0 {
			p.re.Error(ErrUnmatchedRpar);
		}
		return NULL, NULL;
	case ']':
		p.re.Error(ErrUnmatchedRbkt);
	case '^':
		p.nextc();
		start = p.re.Add(new(Bot));
		return start, start;
	case '$':
		p.nextc();
		start = p.re.Add(new(Eot));
		return start, start;
	case '.':
		p.nextc();
		start = p.re.Add(new(Any));
		return start, start;
	case '[':
		p.nextc();
		start = p.CharClass();
		if p.c() != ']' {
			p.re.Error(ErrUnmatchedLbkt);
		}
		p.nextc();
		return start, start;
	case '(':
		p.nextc();
		p.nlpar++;
		start, end = p.Regexp();
		if p.c() != ')' {
			p.re.Error(ErrUnmatchedLpar);
		}
		p.nlpar--;
		p.nextc();
		p.nbra++;
		bra := new(Bra);
		p.re.Add(bra);
		ebra := new(Ebra);
		p.re.Add(ebra);
		bra.n = p.nbra;
		ebra.n = p.nbra;
		if isNULL(start) {
			if !isNULL(end) { p.re.Error(ErrInternal) }
			start = ebra
		} else {
			end.SetNext(ebra);
		}
		bra.SetNext(start);
		return bra, ebra;
	case '\\':
		c = p.nextc();
		switch {
		case c == EOF:
			p.re.Error(ErrExtraneousBackslash);
		case c == 'n':
			c = '\n';
		case special(c):
			// c is as delivered
		default:
			p.re.Error(ErrBadBackslash);
		}
		fallthrough;
	default:
		p.nextc();
		start = NewChar(c);
		p.re.Add(start);
		return start, start
	}
	panic("unreachable");
}

func (p *Parser) Closure() (start, end Inst) {
	start, end = p.Term();
	if isNULL(start) {
		return start, end
	}
	switch p.c() {
	case '*':
		// (start,end)*:
		alt := new(Alt);
		p.re.Add(alt);
		end.SetNext(alt);	// after end, do alt
		alt.left = start;	// alternate brach: return to start
		start = alt;	// alt becomes new (start, end)
		end = alt;
	case '+':
		// (start,end)+:
		alt := new(Alt);
		p.re.Add(alt);
		end.SetNext(alt);	// after end, do alt
		alt.left = start;	// alternate brach: return to start
		end = alt;	// start is unchanged; end is alt
	case '?':
		// (start,end)?:
		alt := new(Alt);
		p.re.Add(alt);
		nop := new(Nop);
		p.re.Add(nop);
		alt.left = start;	// alternate branch is start
		alt.next = nop;	// follow on to nop
		end.SetNext(nop);	// after end, go to nop
		start = alt;	// start is now alt
		end = nop;	// end is nop pointed to by both branches
	default:
		return start, end;
	}
	switch p.nextc() {
	case '*', '+', '?':
		p.re.Error(ErrBadClosure);
	}
	return start, end;
}

func (p *Parser) Concatenation() (start, end Inst) {
	start, end = NULL, NULL;
	for {
		nstart, nend := p.Closure();
		switch {
		case isNULL(nstart):	// end of this concatenation
			if isNULL(start) {	// this is the empty string
				nop := p.re.Add(new(Nop));
				return nop, nop;
			}
			return start, end;
		case isNULL(start):	// this is first element of concatenation
			start, end = nstart, nend;
		default:
			end.SetNext(nstart);
			end = nend;
		}
	}
	panic("unreachable");
}

func (p *Parser) Regexp() (start, end Inst) {
	start, end = p.Concatenation();
	for {
		switch p.c() {
		default:
			return start, end;
		case '|':
			p.nextc();
			nstart, nend := p.Concatenation();
			alt := new(Alt);
			p.re.Add(alt);
			alt.left = start;
			alt.next = nstart;
			nop := new(Nop);
			p.re.Add(nop);
			end.SetNext(nop);
			nend.SetNext(nop);
			start, end = alt, nop;
		}
	}
	panic("unreachable");
}

func UnNop(i Inst) Inst {
	for i.Type() == NOP {
		i = i.Next()
	}
	return i
}

func (re *RE) EliminateNops() {
	for i := 0; i < re.inst.Len(); i++ {
		inst := re.inst.At(i).(Inst);
		if inst.Type() == END {
			continue
		}
		inst.SetNext(UnNop(inst.Next()));
		if inst.Type() == ALT {
			alt := inst.(*Alt);
			alt.left = UnNop(alt.left);
		}
	}
}

func (re *RE) Dump() {
	for i := 0; i < re.inst.Len(); i++ {
		inst := re.inst.At(i).(Inst);
		print(inst.Index(), ": ");
		inst.Print();
		if inst.Type() != END {
			print(" -> ", inst.Next().Index())
		}
		print("\n");
	}
}

func (re *RE) DoParse() {
	parser := NewParser(re);
	start := new(Start);
	re.Add(start);
	s, e := parser.Regexp();
	start.next = s;
	re.start = start;
	e.SetNext(re.Add(new(End)));

	re.Dump();
	println();

	re.EliminateNops();

	re.Dump();
	println();

	re.Error(ErrUnimplemented);
}

func Compiler(str string, ch *chan *RE) {
	re := new(RE);
	re.expr = str;
	re.inst = vector.New();
	re.ch = ch;
	re.DoParse();
	ch <- re;
}

// Public interface has only execute functionality (not yet implemented)
export type Regexp interface {
	// Execute() bool
}

// compile in separate goroutine; wait for result
export func Compile(str string) (regexp Regexp, error *os.Error) {
	ch := new(chan *RE);
	go Compiler(str, ch);
	re := <-ch;
	return re, re.error
}
