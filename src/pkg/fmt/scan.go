// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fmt

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"
	"utf8"
)

// readRuner is the interface to something that can read runes.  If
// the object provided to Scan does not satisfy this interface, the
// object will be wrapped by a readRune object.
type readRuner interface {
	ReadRune() (rune int, size int, err os.Error)
}

// unreadRuner is the interface to something that can unread runes.
// If the object provided to Scan does not satisfy this interface,
// a local buffer will be used to back up the input, but its contents
// will be lost when Scan returns.
type unreadRuner interface {
	UnreadRune() os.Error
}

// ScanState represents the scanner state passed to custom scanners.
// Scanners may do rune-at-a-time scanning or ask the ScanState
// to discover the next space-delimited token.
type ScanState interface {
	// GetRune reads the next rune (Unicode code point) from the input.
	GetRune() (rune int, err os.Error)
	// UngetRune causes the next call to GetRune to return the rune.
	UngetRune()
	// Width returns the value of the width option and whether it has been set.
	// The unit is Unicode code points.
	Width() (wid int, ok bool)
	// Token returns the next space-delimited token from the input. If
	// a width has been specified, the returned token will be no longer
	// than the width.
	Token() (token string, err os.Error)
}

// Scanner is implemented by any value that has a Scan method, which scans
// the input for the representation of a value and stores the result in the
// receiver, which must be a pointer to be useful.  The Scan method is called
// for any argument to Scan or Scanln that implements it.
type Scanner interface {
	Scan(state ScanState, verb int) os.Error
}

// Scan scans text read from standard input, storing successive
// space-separated values into successive arguments.  Newlines count
// as space.  It returns the number of items successfully scanned.
// If that is less than the number of arguments, err will report why.
func Scan(a ...interface{}) (n int, err os.Error) {
	return Fscan(os.Stdin, a)
}

// Scanln is similar to Scan, but stops scanning at a newline and
// after the final item there must be a newline or EOF.
func Scanln(a ...interface{}) (n int, err os.Error) {
	return Fscanln(os.Stdin, a)
}

// Scanf scans text read from standard input, storing successive
// space-separated values into successive arguments as determined by
// the format.  It returns the number of items successfully scanned.
func Scanf(format string, a ...interface{}) (n int, err os.Error) {
	return Fscanf(os.Stdin, format, a)
}

// Sscan scans the argument string, storing successive space-separated
// values into successive arguments.  Newlines count as space.  It
// returns the number of items successfully scanned.  If that is less
// than the number of arguments, err will report why.
func Sscan(str string, a ...interface{}) (n int, err os.Error) {
	return Fscan(strings.NewReader(str), a)
}

// Sscanln is similar to Sscan, but stops scanning at a newline and
// after the final item there must be a newline or EOF.
func Sscanln(str string, a ...interface{}) (n int, err os.Error) {
	return Fscanln(strings.NewReader(str), a)
}

// Sscanf scans the argument string, storing successive space-separated
// values into successive arguments as determined by the format.  It
// returns the number of items successfully parsed.
func Sscanf(str string, format string, a ...interface{}) (n int, err os.Error) {
	return Fscanf(strings.NewReader(str), format, a)
}

// Fscan scans text read from r, storing successive space-separated
// values into successive arguments.  Newlines count as space.  It
// returns the number of items successfully scanned.  If that is less
// than the number of arguments, err will report why.
func Fscan(r io.Reader, a ...interface{}) (n int, err os.Error) {
	s := newScanState(r, true)
	n, err = s.doScan(a)
	s.free()
	return
}

// Fscanln is similar to Fscan, but stops scanning at a newline and
// after the final item there must be a newline or EOF.
func Fscanln(r io.Reader, a ...interface{}) (n int, err os.Error) {
	s := newScanState(r, false)
	n, err = s.doScan(a)
	s.free()
	return
}

// Fscanf scans text read from r, storing successive space-separated
// values into successive arguments as determined by the format.  It
// returns the number of items successfully parsed.
func Fscanf(r io.Reader, format string, a ...interface{}) (n int, err os.Error) {
	s := newScanState(r, false)
	n, err = s.doScanf(format, a)
	s.free()
	return
}

// scanError represents an error generated by the scanning software.
// It's used as a unique signature to identify such errors when recovering.
type scanError struct {
	err os.Error
}

const EOF = -1

// ss is the internal implementation of ScanState.
type ss struct {
	rr         readRuner    // where to read input
	buf        bytes.Buffer // token accumulator
	nlIsSpace  bool         // whether newline counts as white space
	peekRune   int          // one-rune lookahead
	prevRune   int          // last rune returned by GetRune
	atEOF      bool         // already read EOF
	maxWid     int          // max width of field, in runes
	widPresent bool         // width was specified
	wid        int          // width consumed so far; used in accept()
}

func (s *ss) GetRune() (rune int, err os.Error) {
	if s.peekRune >= 0 {
		rune = s.peekRune
		s.prevRune = rune
		s.peekRune = -1
		return
	}
	rune, _, err = s.rr.ReadRune()
	if err == nil {
		s.prevRune = rune
	}
	return
}

func (s *ss) Width() (wid int, ok bool) {
	return s.maxWid, s.widPresent
}

// The public method returns an error; this private one panics.
// If getRune reaches EOF, the return value is EOF (-1).
func (s *ss) getRune() (rune int) {
	if s.atEOF {
		return EOF
	}
	if s.peekRune >= 0 {
		rune = s.peekRune
		s.prevRune = rune
		s.peekRune = -1
		return
	}
	rune, _, err := s.rr.ReadRune()
	if err == nil {
		s.prevRune = rune
	} else if err != nil {
		if err == os.EOF {
			s.atEOF = true
			return EOF
		}
		s.error(err)
	}
	return
}

// mustGetRune turns os.EOF into a panic(io.ErrUnexpectedEOF).
// It is called in cases such as string scanning where an EOF is a
// syntax error.
func (s *ss) mustGetRune() (rune int) {
	if s.atEOF {
		s.error(io.ErrUnexpectedEOF)
	}
	if s.peekRune >= 0 {
		rune = s.peekRune
		s.peekRune = -1
		return
	}
	rune, _, err := s.rr.ReadRune()
	if err != nil {
		if err == os.EOF {
			err = io.ErrUnexpectedEOF
		}
		s.error(err)
	}
	return
}


func (s *ss) UngetRune() {
	if u, ok := s.rr.(unreadRuner); ok {
		u.UnreadRune()
	} else {
		s.peekRune = s.prevRune
	}
}

func (s *ss) error(err os.Error) {
	panic(scanError{err})
}

func (s *ss) errorString(err string) {
	panic(scanError{os.ErrorString(err)})
}

func (s *ss) Token() (tok string, err os.Error) {
	defer func() {
		if e := recover(); e != nil {
			if se, ok := e.(scanError); ok {
				err = se.err
			} else {
				panic(e)
			}
		}
	}()
	tok = s.token()
	return
}

// readRune is a structure to enable reading UTF-8 encoded code points
// from an io.Reader.  It is used if the Reader given to the scanner does
// not already implement ReadRuner.
type readRune struct {
	reader  io.Reader
	buf     [utf8.UTFMax]byte // used only inside ReadRune
	pending int               // number of bytes in pendBuf; only >0 for bad UTF-8
	pendBuf [utf8.UTFMax]byte // bytes left over
}

// readByte returns the next byte from the input, which may be
// left over from a previous read if the UTF-8 was ill-formed.
func (r *readRune) readByte() (b byte, err os.Error) {
	if r.pending > 0 {
		b = r.pendBuf[0]
		copy(r.pendBuf[0:], r.pendBuf[1:])
		r.pending--
		return
	}
	_, err = r.reader.Read(r.pendBuf[0:1])
	return r.pendBuf[0], err
}

// unread saves the bytes for the next read.
func (r *readRune) unread(buf []byte) {
	copy(r.pendBuf[r.pending:], buf)
	r.pending += len(buf)
}

// ReadRune returns the next UTF-8 encoded code point from the
// io.Reader inside r.
func (r *readRune) ReadRune() (rune int, size int, err os.Error) {
	r.buf[0], err = r.readByte()
	if err != nil {
		return 0, 0, err
	}
	if r.buf[0] < utf8.RuneSelf { // fast check for common ASCII case
		rune = int(r.buf[0])
		return
	}
	var n int
	for n = 1; !utf8.FullRune(r.buf[0:n]); n++ {
		r.buf[n], err = r.readByte()
		if err != nil {
			if err == os.EOF {
				err = nil
				break
			}
			return
		}
	}
	rune, size = utf8.DecodeRune(r.buf[0:n])
	if size < n { // an error
		r.unread(r.buf[size:n])
	}
	return
}


// A leaky bucket of reusable ss structures.
var ssFree = make(chan *ss, 100)

// Allocate a new ss struct.  Probably can grab the previous one from ssFree.
func newScanState(r io.Reader, nlIsSpace bool) *ss {
	s, ok := <-ssFree
	if !ok {
		s = new(ss)
	}
	if rr, ok := r.(readRuner); ok {
		s.rr = rr
	} else {
		s.rr = &readRune{reader: r}
	}
	s.nlIsSpace = nlIsSpace
	s.peekRune = -1
	s.atEOF = false
	s.maxWid = 0
	s.widPresent = false
	return s
}

// Save used ss structs in ssFree; avoid an allocation per invocation.
func (s *ss) free() {
	// Don't hold on to ss structs with large buffers.
	if cap(s.buf.Bytes()) > 1024 {
		return
	}
	s.buf.Reset()
	s.rr = nil
	_ = ssFree <- s
}

// skipSpace skips spaces and maybe newlines.
func (s *ss) skipSpace(stopAtNewline bool) {
	for {
		rune := s.getRune()
		if rune == EOF {
			return
		}
		if rune == '\n' {
			if stopAtNewline {
				break
			}
			if s.nlIsSpace {
				continue
			}
			s.errorString("unexpected newline")
			return
		}
		if !unicode.IsSpace(rune) {
			s.UngetRune()
			break
		}
	}
}

// token returns the next space-delimited string from the input.  It
// skips white space.  For Scanln, it stops at newlines.  For Scan,
// newlines are treated as spaces.
func (s *ss) token() string {
	s.skipSpace(false)
	// read until white space or newline
	for nrunes := 0; !s.widPresent || nrunes < s.maxWid; nrunes++ {
		rune := s.getRune()
		if rune == EOF {
			break
		}
		if unicode.IsSpace(rune) {
			s.UngetRune()
			break
		}
		s.buf.WriteRune(rune)
	}
	return s.buf.String()
}

// typeError indicates that the type of the operand did not match the format
func (s *ss) typeError(field interface{}, expected string) {
	s.errorString("expected field of type pointer to " + expected + "; found " + reflect.Typeof(field).String())
}

var complexError = os.ErrorString("syntax error scanning complex number")
var boolError = os.ErrorString("syntax error scanning boolean")

// accepts checks the next rune in the input.  If it's a byte (sic) in the string, it puts it in the
// buffer and returns true. Otherwise it return false.
func (s *ss) accept(ok string) bool {
	if s.wid >= s.maxWid {
		return false
	}
	rune := s.getRune()
	if rune == EOF {
		return false
	}
	for i := 0; i < len(ok); i++ {
		if int(ok[i]) == rune {
			s.buf.WriteRune(rune)
			s.wid++
			return true
		}
	}
	if rune != EOF {
		s.UngetRune()
	}
	return false
}

// okVerb verifies that the verb is present in the list, setting s.err appropriately if not.
func (s *ss) okVerb(verb int, okVerbs, typ string) bool {
	for _, v := range okVerbs {
		if v == verb {
			return true
		}
	}
	s.errorString("bad verb %" + string(verb) + " for " + typ)
	return false
}

// scanBool returns the value of the boolean represented by the next token.
func (s *ss) scanBool(verb int) bool {
	if !s.okVerb(verb, "tv", "boolean") {
		return false
	}
	// Syntax-checking a boolean is annoying.  We're not fastidious about case.
	switch s.mustGetRune() {
	case '0':
		return false
	case '1':
		return true
	case 't', 'T':
		if s.accept("rR") && (!s.accept("uU") || !s.accept("eE")) {
			s.error(boolError)
		}
		return true
	case 'f', 'F':
		if s.accept("aL") && (!s.accept("lL") || !s.accept("sS") || !s.accept("eE")) {
			s.error(boolError)
		}
		return false
	}
	return false
}

// Numerical elements
const (
	binaryDigits      = "01"
	octalDigits       = "01234567"
	decimalDigits     = "0123456789"
	hexadecimalDigits = "0123456789aAbBcCdDeEfF"
	sign              = "+-"
	period            = "."
	exponent          = "eE"
)

// getBase returns the numeric base represented by the verb and its digit string.
func (s *ss) getBase(verb int) (base int, digits string) {
	s.okVerb(verb, "bdoxXv", "integer") // sets s.err
	base = 10
	digits = decimalDigits
	switch verb {
	case 'b':
		base = 2
		digits = binaryDigits
	case 'o':
		base = 8
		digits = octalDigits
	case 'x', 'X':
		base = 16
		digits = hexadecimalDigits
	}
	return
}

// scanNumber returns the numerical string with specified digits starting here.
func (s *ss) scanNumber(digits string) string {
	if !s.accept(digits) {
		s.errorString("expected integer")
	}
	for s.accept(digits) {
	}
	return s.buf.String()
}

// scanRune returns the next rune value in the input.
func (s *ss) scanRune(bitSize int) int64 {
	rune := int64(s.mustGetRune())
	n := uint(bitSize)
	x := (rune << (64 - n)) >> (64 - n)
	if x != rune {
		s.errorString("overflow on character value " + string(rune))
	}
	return rune
}

// scanInt returns the value of the integer represented by the next
// token, checking for overflow.  Any error is stored in s.err.
func (s *ss) scanInt(verb int, bitSize int) int64 {
	if verb == 'c' {
		return s.scanRune(bitSize)
	}
	base, digits := s.getBase(verb)
	s.skipSpace(false)
	s.accept(sign) // If there's a sign, it will be left in the token buffer.
	tok := s.scanNumber(digits)
	i, err := strconv.Btoi64(tok, base)
	if err != nil {
		s.error(err)
	}
	n := uint(bitSize)
	x := (i << (64 - n)) >> (64 - n)
	if x != i {
		s.errorString("integer overflow on token " + tok)
	}
	return i
}

// scanUint returns the value of the unsigned integer represented
// by the next token, checking for overflow.  Any error is stored in s.err.
func (s *ss) scanUint(verb int, bitSize int) uint64 {
	if verb == 'c' {
		return uint64(s.scanRune(bitSize))
	}
	base, digits := s.getBase(verb)
	s.skipSpace(false)
	tok := s.scanNumber(digits)
	i, err := strconv.Btoui64(tok, base)
	if err != nil {
		s.error(err)
	}
	n := uint(bitSize)
	x := (i << (64 - n)) >> (64 - n)
	if x != i {
		s.errorString("unsigned integer overflow on token " + tok)
	}
	return i
}

// floatToken returns the floating-point number starting here, no longer than swid
// if the width is specified. It's not rigorous about syntax because it doesn't check that
// we have at least some digits, but Atof will do that.
func (s *ss) floatToken() string {
	s.buf.Reset()
	// leading sign?
	s.accept(sign)
	// digits?
	for s.accept(decimalDigits) {
	}
	// decimal point?
	if s.accept(period) {
		// fraction?
		for s.accept(decimalDigits) {
		}
	}
	// exponent?
	if s.accept(exponent) {
		// leading sign?
		s.accept(sign)
		// digits?
		for s.accept(decimalDigits) {
		}
	}
	return s.buf.String()
}

// complexTokens returns the real and imaginary parts of the complex number starting here.
// The number might be parenthesized and has the format (N+Ni) where N is a floating-point
// number and there are no spaces within.
func (s *ss) complexTokens() (real, imag string) {
	// TODO: accept N and Ni independently?
	parens := s.accept("(")
	real = s.floatToken()
	s.buf.Reset()
	// Must now have a sign.
	if !s.accept("+-") {
		s.error(complexError)
	}
	// Sign is now in buffer
	imagSign := s.buf.String()
	imag = s.floatToken()
	if !s.accept("i") {
		s.error(complexError)
	}
	if parens && !s.accept(")") {
		s.error(complexError)
	}
	return real, imagSign + imag
}

// convertFloat converts the string to a float64value.
func (s *ss) convertFloat(str string, n int) float64 {
	f, err := strconv.AtofN(str, n)
	if err != nil {
		s.error(err)
	}
	return f
}

// convertComplex converts the next token to a complex128 value.
// The atof argument is a type-specific reader for the underlying type.
// If we're reading complex64, atof will parse float32s and convert them
// to float64's to avoid reproducing this code for each complex type.
func (s *ss) scanComplex(verb int, n int) complex128 {
	if !s.okVerb(verb, floatVerbs, "complex") {
		return 0
	}
	s.skipSpace(false)
	sreal, simag := s.complexTokens()
	real := s.convertFloat(sreal, n/2)
	imag := s.convertFloat(simag, n/2)
	return cmplx(real, imag)
}

// convertString returns the string represented by the next input characters.
// The format of the input is determined by the verb.
func (s *ss) convertString(verb int) (str string) {
	if !s.okVerb(verb, "svqx", "string") {
		return ""
	}
	s.skipSpace(false)
	switch verb {
	case 'q':
		str = s.quotedString()
	case 'x':
		str = s.hexString()
	default:
		str = s.token() // %s and %v just return the next word
	}
	// Empty strings other than with %q are not OK.
	if len(str) == 0 && verb != 'q' && s.maxWid > 0 {
		s.errorString("Scan: no data for string")
	}
	return
}

// quotedString returns the double- or back-quoted string represented by the next input characters.
func (s *ss) quotedString() string {
	quote := s.mustGetRune()
	switch quote {
	case '`':
		// Back-quoted: Anything goes until EOF or back quote.
		for {
			rune := s.mustGetRune()
			if rune == quote {
				break
			}
			s.buf.WriteRune(rune)
		}
		return s.buf.String()
	case '"':
		// Double-quoted: Include the quotes and let strconv.Unquote do the backslash escapes.
		s.buf.WriteRune(quote)
		for {
			rune := s.mustGetRune()
			s.buf.WriteRune(rune)
			if rune == '\\' {
				// In a legal backslash escape, no matter how long, only the character
				// immediately after the escape can itself be a backslash or quote.
				// Thus we only need to protect the first character after the backslash.
				rune := s.mustGetRune()
				s.buf.WriteRune(rune)
			} else if rune == '"' {
				break
			}
		}
		result, err := strconv.Unquote(s.buf.String())
		if err != nil {
			s.error(err)
		}
		return result
	default:
		s.errorString("expected quoted string")
	}
	return ""
}

// hexDigit returns the value of the hexadecimal digit
func (s *ss) hexDigit(digit int) int {
	switch digit {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return digit - '0'
	case 'a', 'b', 'c', 'd', 'e', 'f':
		return 10 + digit - 'a'
	case 'A', 'B', 'C', 'D', 'E', 'F':
		return 10 + digit - 'A'
	}
	s.errorString("Scan: illegal hex digit")
	return 0
}

// hexByte returns the next hex-encoded (two-character) byte from the input.
// There must be either two hexadecimal digits or a space character in the input.
func (s *ss) hexByte() (b byte, ok bool) {
	rune1 := s.getRune()
	if rune1 == EOF {
		return
	}
	if unicode.IsSpace(rune1) {
		s.UngetRune()
		return
	}
	rune2 := s.mustGetRune()
	return byte(s.hexDigit(rune1)<<4 | s.hexDigit(rune2)), true
}

// hexString returns the space-delimited hexpair-encoded string.
func (s *ss) hexString() string {
	for {
		b, ok := s.hexByte()
		if !ok {
			break
		}
		s.buf.WriteByte(b)
	}
	if s.buf.Len() == 0 {
		s.errorString("Scan: no hex data for %x string")
		return ""
	}
	return s.buf.String()
}

const floatVerbs = "eEfFgGv"

// scanOne scans a single value, deriving the scanner from the type of the argument.
func (s *ss) scanOne(verb int, field interface{}) {
	s.buf.Reset()
	var err os.Error
	// If the parameter has its own Scan method, use that.
	if v, ok := field.(Scanner); ok {
		err = v.Scan(s, verb)
		if err != nil {
			s.error(err)
		}
		return
	}
	if !s.widPresent {
		s.maxWid = 1 << 30 // Huge
	}
	s.wid = 0
	switch v := field.(type) {
	case *bool:
		*v = s.scanBool(verb)
	case *complex:
		*v = complex(s.scanComplex(verb, int(complexBits)))
	case *complex64:
		*v = complex64(s.scanComplex(verb, 64))
	case *complex128:
		*v = s.scanComplex(verb, 128)
	case *int:
		*v = int(s.scanInt(verb, intBits))
	case *int8:
		*v = int8(s.scanInt(verb, 8))
	case *int16:
		*v = int16(s.scanInt(verb, 16))
	case *int32:
		*v = int32(s.scanInt(verb, 32))
	case *int64:
		*v = s.scanInt(verb, 64)
	case *uint:
		*v = uint(s.scanUint(verb, intBits))
	case *uint8:
		*v = uint8(s.scanUint(verb, 8))
	case *uint16:
		*v = uint16(s.scanUint(verb, 16))
	case *uint32:
		*v = uint32(s.scanUint(verb, 32))
	case *uint64:
		*v = s.scanUint(verb, 64)
	case *uintptr:
		*v = uintptr(s.scanUint(verb, uintptrBits))
	// Floats are tricky because you want to scan in the precision of the result, not
	// scan in high precision and convert, in order to preserve the correct error condition.
	case *float:
		if s.okVerb(verb, floatVerbs, "float") {
			s.skipSpace(false)
			*v = float(s.convertFloat(s.floatToken(), int(floatBits)))
		}
	case *float32:
		if s.okVerb(verb, floatVerbs, "float32") {
			s.skipSpace(false)
			*v = float32(s.convertFloat(s.floatToken(), 32))
		}
	case *float64:
		if s.okVerb(verb, floatVerbs, "float64") {
			s.skipSpace(false)
			*v = s.convertFloat(s.floatToken(), 64)
		}
	case *string:
		*v = s.convertString(verb)
	case *[]byte:
		// We scan to string and convert so we get a copy of the data.
		// If we scanned to bytes, the slice would point at the buffer.
		*v = []byte(s.convertString(verb))
	default:
		val := reflect.NewValue(v)
		ptr, ok := val.(*reflect.PtrValue)
		if !ok {
			s.errorString("Scan: type not a pointer: " + val.Type().String())
			return
		}
		switch v := ptr.Elem().(type) {
		case *reflect.BoolValue:
			v.Set(s.scanBool(verb))
		case *reflect.IntValue:
			v.Set(s.scanInt(verb, v.Type().Bits()))
		case *reflect.UintValue:
			v.Set(s.scanUint(verb, v.Type().Bits()))
		case *reflect.StringValue:
			v.Set(s.convertString(verb))
		case *reflect.SliceValue:
			// For now, can only handle (renamed) []byte.
			typ := v.Type().(*reflect.SliceType)
			if typ.Elem().Kind() != reflect.Uint8 {
				goto CantHandle
			}
			str := s.convertString(verb)
			v.Set(reflect.MakeSlice(typ, len(str), len(str)))
			for i := 0; i < len(str); i++ {
				v.Elem(i).(*reflect.UintValue).Set(uint64(str[i]))
			}
		case *reflect.FloatValue:
			s.skipSpace(false)
			v.Set(s.convertFloat(s.floatToken(), v.Type().Bits()))
		case *reflect.ComplexValue:
			v.Set(s.scanComplex(verb, v.Type().Bits()))
		default:
		CantHandle:
			s.errorString("Scan: can't handle type: " + val.Type().String())
		}
	}
}

// errorHandler turns local panics into error returns.  EOFs are benign.
func errorHandler(errp *os.Error) {
	if e := recover(); e != nil {
		if se, ok := e.(scanError); ok { // catch local error
			if se.err != os.EOF {
				*errp = se.err
			}
		} else {
			panic(e)
		}
	}
}

// doScan does the real work for scanning without a format string.
// At the moment, it handles only pointers to basic types.
func (s *ss) doScan(a []interface{}) (numProcessed int, err os.Error) {
	defer errorHandler(&err)
	for _, field := range a {
		s.scanOne('v', field)
		numProcessed++
	}
	// Check for newline if required.
	if !s.nlIsSpace {
		for {
			rune := s.getRune()
			if rune == '\n' || rune == EOF {
				break
			}
			if !unicode.IsSpace(rune) {
				s.errorString("Scan: expected newline")
				break
			}
		}
	}
	return
}

// advance determines whether the next characters in the input match
// those of the format.  It returns the number of bytes (sic) consumed
// in the format. Newlines included, all runs of space characters in
// either input or format behave as a single space. This routine also
// handles the %% case.  If the return value is zero, either format
// starts with a % (with no following %) or the input is empty.
// If it is negative, the input did not match the string.
func (s *ss) advance(format string) (i int) {
	for i < len(format) {
		fmtc, w := utf8.DecodeRuneInString(format[i:])
		if fmtc == '%' {
			// %% acts like a real percent
			nextc, _ := utf8.DecodeRuneInString(format[i+w:]) // will not match % if string is empty
			if nextc != '%' {
				return
			}
			i += w // skip the first %
		}
		sawSpace := false
		for unicode.IsSpace(fmtc) && i < len(format) {
			sawSpace = true
			i += w
			fmtc, w = utf8.DecodeRuneInString(format[i:])
		}
		if sawSpace {
			// There was space in the format, so there should be space (EOF)
			// in the input.
			inputc := s.getRune()
			if inputc == EOF {
				return
			}
			if !unicode.IsSpace(inputc) {
				// Space in format but not in input: error
				s.errorString("expected space in input to match format")
			}
			s.skipSpace(true)
			continue
		}
		inputc := s.mustGetRune()
		if fmtc != inputc {
			s.UngetRune()
			return -1
		}
		i += w
	}
	return
}

// doScanf does the real work when scanning with a format string.
//  At the moment, it handles only pointers to basic types.
func (s *ss) doScanf(format string, a []interface{}) (numProcessed int, err os.Error) {
	defer errorHandler(&err)
	end := len(format) - 1
	// We process one item per non-trivial format
	for i := 0; i <= end; {
		w := s.advance(format[i:])
		if w > 0 {
			i += w
			continue
		}
		// Either we failed to advance, we have a percent character, or we ran out of input.
		if format[i] != '%' {
			// Can't advance format.  Why not?
			if w < 0 {
				s.errorString("input does not match format")
			}
			// Otherwise at EOF; "too many operands" error handled below
			break
		}
		i++ // % is one byte

		// do we have 20 (width)?
		s.maxWid, s.widPresent, i = parsenum(format, i, end)

		c, w := utf8.DecodeRuneInString(format[i:])
		i += w

		if numProcessed >= len(a) { // out of operands
			s.errorString("too few operands for format %" + format[i-w:])
			break
		}
		field := a[numProcessed]

		s.scanOne(c, field)
		numProcessed++
	}
	if numProcessed < len(a) {
		s.errorString("too many operands")
	}
	return
}
