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
	"unicode"
	"utf8"
)

// readRuner is the interface to something that can read runes.  If
// the object provided to Scan does not satisfy this interface, the
// object will be wrapped by a readRune object.
type readRuner interface {
	ReadRune() (rune int, size int, err os.Error)
}

// ScanState represents the scanner state passed to custom scanners.
// Scanners may do rune-at-a-time scanning or ask the ScanState
// to discover the next space-delimited token.
type ScanState interface {
	// GetRune reads the next rune (Unicode code point) from the input.
	GetRune() (rune int, err os.Error)
	// UngetRune causes the next call to Get to return the rune.
	UngetRune(rune int)
	// Token returns the next space-delimited token from the input.
	Token() (token string, err os.Error)
}

// Scanner is implemented by any value that has a Scan method, which scans
// the input for the representation of a value and stores the result in the
// receiver, which must be a pointer to be useful.  The Scan method is called
// for any argument to Scan or Scanln that implements it.
type Scanner interface {
	Scan(ScanState) os.Error
}

// Scan parses text read from standard input, storing successive
// space-separated values into successive arguments.  Newlines count as
// space.  Each argument must be a pointer to a basic type or an
// implementation of the Scanner interface.  It returns the number of items
// successfully parsed.  If that is less than the number of arguments, err
// will report why.
func Scan(a ...interface{}) (n int, err os.Error) {
	return Fscan(os.Stdin, a)
}

// Fscanln parses text read from standard input, storing successive
// space-separated values into successive arguments.  Scanning stops at a
// newline and after the final item there must be a newline or EOF.  Each
// argument must be a pointer to a basic type or an implementation of the
// Scanner interface.  It returns the number of items successfully parsed.
// If that is less than the number of arguments, err will report why.
func Scanln(a ...interface{}) (n int, err os.Error) {
	return Fscanln(os.Stdin, a)
}

// Fscan parses text read from r, storing successive space-separated values
// into successive arguments.  Newlines count as space.  Each argument must
// be a pointer to a basic type or an implementation of the Scanner
// interface.  It returns the number of items successfully parsed.  If that
// is less than the number of arguments, err will report why.
func Fscan(r io.Reader, a ...interface{}) (n int, err os.Error) {
	s := newScanState(r, true)
	n, err = s.doScan(a)
	s.free()
	return
}

// Fscanln parses text read from r, storing successive space-separated values
// into successive arguments.  Scanning stops at a newline and after the
// final item there must be a newline or EOF.  Each argument must be a
// pointer to a basic type or an implementation of the Scanner interface.  It
// returns the number of items successfully parsed.  If that is less than the
// number of arguments, err will report why.
func Fscanln(r io.Reader, a ...interface{}) (n int, err os.Error) {
	s := newScanState(r, false)
	n, err = s.doScan(a)
	s.free()
	return
}

// XXXScanf is incomplete, do not use.
func XXXScanf(format string, a ...interface{}) (n int, err os.Error) {
	return XXXFscanf(os.Stdin, format, a)
}

// XXXFscanf is incomplete, do not use.
func XXXFscanf(r io.Reader, format string, a ...interface{}) (n int, err os.Error) {
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

// ss is the internal implementation of ScanState.
type ss struct {
	rr        readRuner    // where to read input
	buf       bytes.Buffer // token accumulator
	nlIsSpace bool         // whether newline counts as white space
	peekRune  int          // one-rune lookahead
}

func (s *ss) GetRune() (rune int, err os.Error) {
	if s.peekRune >= 0 {
		rune = s.peekRune
		s.peekRune = -1
		return
	}
	rune, _, err = s.rr.ReadRune()
	return
}

const EOF = -1

// The public method returns an error; this private one panics.
// If getRune reaches EOF, the return value is EOF (-1).
func (s *ss) getRune() (rune int) {
	if s.peekRune >= 0 {
		rune = s.peekRune
		s.peekRune = -1
		return
	}
	rune, _, err := s.rr.ReadRune()
	if err != nil {
		if err == os.EOF {
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


func (s *ss) UngetRune(rune int) {
	s.peekRune = rune
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
// TODO: readByteRune for things that can read bytes.
type readRune struct {
	reader io.Reader
	buf    [utf8.UTFMax]byte
}

// ReadRune returns the next UTF-8 encoded code point from the
// io.Reader inside r.
func (r readRune) ReadRune() (rune int, size int, err os.Error) {
	_, err = r.reader.Read(r.buf[0:1])
	if err != nil {
		return 0, 0, err
	}
	if r.buf[0] < utf8.RuneSelf { // fast check for common ASCII case
		rune = int(r.buf[0])
		return
	}
	for size := 1; size < utf8.UTFMax; size++ {
		_, err = r.reader.Read(r.buf[size : size+1])
		if err != nil {
			break
		}
		if !utf8.FullRune(r.buf[0:]) {
			continue
		}
		if c, w := utf8.DecodeRune(r.buf[0:size]); w == size {
			rune = c
			return
		}
	}
	return utf8.RuneError, 1, err
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
		s.rr = readRune{reader: r}
	}
	s.nlIsSpace = nlIsSpace
	s.peekRune = -1
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

// skipSpace skips spaces and maybe newlines
func (s *ss) skipSpace() {
	s.buf.Reset()
	for {
		rune := s.getRune()
		if rune == EOF {
			return
		}
		if rune == '\n' {
			if s.nlIsSpace {
				continue
			}
			s.errorString("unexpected newline")
			return
		}
		if !unicode.IsSpace(rune) {
			s.UngetRune(rune)
			break
		}
	}
}

// token returns the next space-delimited string from the input.
// For Scanln, it stops at newlines.  For Scan, newlines are treated as
// spaces.
func (s *ss) token() string {
	s.skipSpace()
	// read until white space or newline
	for {
		rune := s.getRune()
		if rune == EOF {
			break
		}
		if unicode.IsSpace(rune) {
			s.UngetRune(rune)
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

var intBits = uint(reflect.Typeof(int(0)).Size() * 8)
var uintptrBits = uint(reflect.Typeof(int(0)).Size() * 8)
var complexError = os.ErrorString("syntax error scanning complex number")

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
	tok := s.token()
	b, err := strconv.Atob(tok)
	if err != nil {
		s.error(err)
	}
	return b
}

// getBase returns the numeric base represented by the verb.
func (s *ss) getBase(verb int) int {
	s.okVerb(verb, "bdoxXv", "integer") // sets s.err
	base := 10
	switch verb {
	case 'b':
		base = 2
	case 'o':
		base = 8
	case 'x', 'X':
		base = 16
	}
	return base
}

// scanInt returns the value of the integer represented by the next
// token, checking for overflow.  Any error is stored in s.err.
func (s *ss) scanInt(verb int, bitSize uint) int64 {
	base := s.getBase(verb)
	tok := s.token()
	i, err := strconv.Btoi64(tok, base)
	if err != nil {
		s.error(err)
	}
	x := (i << (64 - bitSize)) >> (64 - bitSize)
	if x != i {
		s.errorString("integer overflow on token " + tok)
	}
	return i
}

// scanUint returns the value of the unsigned integer represented
// by the next token, checking for overflow.  Any error is stored in s.err.
func (s *ss) scanUint(verb int, bitSize uint) uint64 {
	base := s.getBase(verb)
	tok := s.token()
	i, err := strconv.Btoui64(tok, base)
	if err != nil {
		s.error(err)
	}
	x := (i << (64 - bitSize)) >> (64 - bitSize)
	if x != i {
		s.errorString("unsigned integer overflow on token " + tok)
	}
	return i
}

// complexParts returns the strings representing the real and imaginary parts of the string.
func (s *ss) complexParts(str string) (real, imag string) {
	if len(str) > 2 && str[0] == '(' && str[len(str)-1] == ')' {
		str = str[1 : len(str)-1]
	}
	real, str = floatPart(str)
	// Must now have a sign.
	if len(str) == 0 || (str[0] != '+' && str[0] != '-') {
		s.error(complexError)
	}
	imag, str = floatPart(str)
	if str != "i" {
		s.error(complexError)
	}
	return real, imag
}

// floatPart returns strings holding the floating point value in the string, followed
// by the remainder of the string.  That is, it splits str into (number,rest-of-string).
func floatPart(str string) (first, last string) {
	i := 0
	// leading sign?
	if len(str) > i && (str[0] == '+' || str[0] == '-') {
		i++
	}
	// digits?
	for len(str) > i && '0' <= str[i] && str[i] <= '9' {
		i++
	}
	// period?
	if str[i] == '.' {
		i++
	}
	// fraction?
	for len(str) > i && '0' <= str[i] && str[i] <= '9' {
		i++
	}
	// exponent?
	if len(str) > i && (str[i] == 'e' || str[i] == 'E') {
		i++
		// leading sign?
		if str[i] == '+' || str[i] == '-' {
			i++
		}
		// digits?
		for len(str) > i && '0' <= str[i] && str[i] <= '9' {
			i++
		}
	}
	return str[0:i], str[i:]
}

// convertFloat converts the string to a float value.
func (s *ss) convertFloat(str string) float64 {
	f, err := strconv.Atof(str)
	if err != nil {
		s.error(err)
	}
	return float64(f)
}

// convertFloat32 converts the string to a float32 value.
func (s *ss) convertFloat32(str string) float64 {
	f, err := strconv.Atof32(str)
	if err != nil {
		s.error(err)
	}
	return float64(f)
}

// convertFloat64 converts the string to a float64 value.
func (s *ss) convertFloat64(str string) float64 {
	f, err := strconv.Atof64(str)
	if err != nil {
		s.error(err)
	}
	return f
}

// convertComplex converts the next token to a complex128 value.
// The atof argument is a type-specific reader for the underlying type.
// If we're reading complex64, atof will parse float32s and convert them
// to float64's to avoid reproducing this code for each complex type.
func (s *ss) scanComplex(verb int, atof func(*ss, string) float64) complex128 {
	if !s.okVerb(verb, floatVerbs, "complex") {
		return 0
	}
	tok := s.token()
	sreal, simag := s.complexParts(tok)
	real := atof(s, sreal)
	imag := atof(s, simag)
	return cmplx(real, imag)
}

// convertString returns the string represented by the next input characters.
// The format of the input is determined by the verb.
func (s *ss) convertString(verb int) string {
	if !s.okVerb(verb, "svqx", "string") {
		return ""
	}
	s.skipSpace()
	switch verb {
	case 'q':
		return s.quotedString()
	case 'x':
		return s.hexString()
	}
	return s.token() // %s and %v just return the next word
}

// quotedString returns the double- or back-quoted string.
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
		s.UngetRune(rune1)
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
	var err os.Error
	// If the parameter has its own Scan method, use that.
	if v, ok := field.(Scanner); ok {
		err = v.Scan(s)
		if err != nil {
			s.error(err)
		}
		return
	}
	switch v := field.(type) {
	case *bool:
		*v = s.scanBool(verb)
	case *complex:
		*v = complex(s.scanComplex(verb, (*ss).convertFloat))
	case *complex64:
		*v = complex64(s.scanComplex(verb, (*ss).convertFloat32))
	case *complex128:
		*v = s.scanComplex(verb, (*ss).convertFloat64)
	case *int:
		*v = int(s.scanInt(verb, intBits))
	case *int8:
		*v = int8(s.scanInt(verb, 8))
	case *int16:
		*v = int16(s.scanInt(verb, 16))
	case *int32:
		*v = int32(s.scanInt(verb, 32))
	case *int64:
		*v = s.scanInt(verb, intBits)
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
			*v = float(s.convertFloat(s.token()))
		}
	case *float32:
		if s.okVerb(verb, floatVerbs, "float32") {
			*v = float32(s.convertFloat32(s.token()))
		}
	case *float64:
		if s.okVerb(verb, floatVerbs, "float64") {
			*v = s.convertFloat64(s.token())
		}
	case *string:
		*v = s.convertString(verb)
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
			v.Set(int(s.scanInt(verb, intBits)))
		case *reflect.Int8Value:
			v.Set(int8(s.scanInt(verb, 8)))
		case *reflect.Int16Value:
			v.Set(int16(s.scanInt(verb, 16)))
		case *reflect.Int32Value:
			v.Set(int32(s.scanInt(verb, 32)))
		case *reflect.Int64Value:
			v.Set(s.scanInt(verb, 64))
		case *reflect.UintValue:
			v.Set(uint(s.scanUint(verb, intBits)))
		case *reflect.Uint8Value:
			v.Set(uint8(s.scanUint(verb, 8)))
		case *reflect.Uint16Value:
			v.Set(uint16(s.scanUint(verb, 16)))
		case *reflect.Uint32Value:
			v.Set(uint32(s.scanUint(verb, 32)))
		case *reflect.Uint64Value:
			v.Set(s.scanUint(verb, 64))
		case *reflect.UintptrValue:
			v.Set(uintptr(s.scanUint(verb, uintptrBits)))
		case *reflect.StringValue:
			v.Set(s.convertString(verb))
		case *reflect.FloatValue:
			v.Set(float(s.convertFloat(s.token())))
		case *reflect.Float32Value:
			v.Set(float32(s.convertFloat(s.token())))
		case *reflect.Float64Value:
			v.Set(s.convertFloat(s.token()))
		case *reflect.ComplexValue:
			v.Set(complex(s.scanComplex(verb, (*ss).convertFloat)))
		case *reflect.Complex64Value:
			v.Set(complex64(s.scanComplex(verb, (*ss).convertFloat32)))
		case *reflect.Complex128Value:
			v.Set(s.scanComplex(verb, (*ss).convertFloat64))
		default:
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

// doScanf does the real work when scanning with a format string.
//  At the moment, it handles only pointers to basic types.
func (s *ss) doScanf(format string, a []interface{}) (numProcessed int, err os.Error) {
	defer errorHandler(&err)
	end := len(format) - 1
	// We process one item per non-trivial format
	for i := 0; i <= end; {
		c, w := utf8.DecodeRuneInString(format[i:])
		if c != '%' || i == end {
			// TODO: WHAT NOW?
			i += w
			continue
		}
		i++
		// TODO: FLAGS
		c, w = utf8.DecodeRuneInString(format[i:])
		i += w
		// percent is special - absorbs no operand
		if c == '%' {
			// TODO: WHAT NOW?
			continue
		}

		if numProcessed >= len(a) { // out of operands
			s.errorString("too few operands for format %" + format[i-w:])
			break
		}
		field := a[numProcessed]

		s.scanOne(c, field)
		numProcessed++
	}
	return
}
