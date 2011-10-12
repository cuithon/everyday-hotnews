// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fmt

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"sync"
	"unicode"
	"utf8"
)

// Some constants in the form of bytes, to avoid string overhead.
// Needlessly fastidious, I suppose.
var (
	commaSpaceBytes = []byte(", ")
	nilAngleBytes   = []byte("<nil>")
	nilParenBytes   = []byte("(nil)")
	nilBytes        = []byte("nil")
	mapBytes        = []byte("map[")
	missingBytes    = []byte("(MISSING)")
	panicBytes      = []byte("(PANIC=")
	extraBytes      = []byte("%!(EXTRA ")
	irparenBytes    = []byte("i)")
	bytesBytes      = []byte("[]byte{")
	widthBytes      = []byte("%!(BADWIDTH)")
	precBytes       = []byte("%!(BADPREC)")
	noVerbBytes     = []byte("%!(NOVERB)")
)

// State represents the printer state passed to custom formatters.
// It provides access to the io.Writer interface plus information about
// the flags and options for the operand's format specifier.
type State interface {
	// Write is the function to call to emit formatted output to be printed.
	Write(b []byte) (ret int, err os.Error)
	// Width returns the value of the width option and whether it has been set.
	Width() (wid int, ok bool)
	// Precision returns the value of the precision option and whether it has been set.
	Precision() (prec int, ok bool)

	// Flag returns whether the flag c, a character, has been set.
	Flag(c int) bool
}

// Formatter is the interface implemented by values with a custom formatter.
// The implementation of Format may call Sprintf or Fprintf(f) etc.
// to generate its output.
type Formatter interface {
	Format(f State, c int)
}

// Stringer is implemented by any value that has a String method,
// which defines the ``native'' format for that value.
// The String method is used to print values passed as an operand
// to a %s or %v format or to an unformatted printer such as Print.
type Stringer interface {
	String() string
}

// GoStringer is implemented by any value that has a GoString method,
// which defines the Go syntax for that value.
// The GoString method is used to print values passed as an operand
// to a %#v format.
type GoStringer interface {
	GoString() string
}

type pp struct {
	n         int
	panicking bool
	buf       bytes.Buffer
	runeBuf   [utf8.UTFMax]byte
	fmt       fmt
}

// A cache holds a set of reusable objects.
// The slice is a stack (LIFO).
// If more are needed, the cache creates them by calling new.
type cache struct {
	mu    sync.Mutex
	saved []interface{}
	new   func() interface{}
}

func (c *cache) put(x interface{}) {
	c.mu.Lock()
	if len(c.saved) < cap(c.saved) {
		c.saved = append(c.saved, x)
	}
	c.mu.Unlock()
}

func (c *cache) get() interface{} {
	c.mu.Lock()
	n := len(c.saved)
	if n == 0 {
		c.mu.Unlock()
		return c.new()
	}
	x := c.saved[n-1]
	c.saved = c.saved[0 : n-1]
	c.mu.Unlock()
	return x
}

func newCache(f func() interface{}) *cache {
	return &cache{saved: make([]interface{}, 0, 100), new: f}
}

var ppFree = newCache(func() interface{} { return new(pp) })

// Allocate a new pp struct or grab a cached one.
func newPrinter() *pp {
	p := ppFree.get().(*pp)
	p.panicking = false
	p.fmt.init(&p.buf)
	return p
}

// Save used pp structs in ppFree; avoids an allocation per invocation.
func (p *pp) free() {
	// Don't hold on to pp structs with large buffers.
	if cap(p.buf.Bytes()) > 1024 {
		return
	}
	p.buf.Reset()
	ppFree.put(p)
}

func (p *pp) Width() (wid int, ok bool) { return p.fmt.wid, p.fmt.widPresent }

func (p *pp) Precision() (prec int, ok bool) { return p.fmt.prec, p.fmt.precPresent }

func (p *pp) Flag(b int) bool {
	switch b {
	case '-':
		return p.fmt.minus
	case '+':
		return p.fmt.plus
	case '#':
		return p.fmt.sharp
	case ' ':
		return p.fmt.space
	case '0':
		return p.fmt.zero
	}
	return false
}

func (p *pp) add(c int) {
	p.buf.WriteRune(c)
}

// Implement Write so we can call Fprintf on a pp (through State), for
// recursive use in custom verbs.
func (p *pp) Write(b []byte) (ret int, err os.Error) {
	return p.buf.Write(b)
}

// These routines end in 'f' and take a format string.

// Fprintf formats according to a format specifier and writes to w.
// It returns the number of bytes written and any write error encountered.
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err os.Error) {
	p := newPrinter()
	p.doPrintf(format, a)
	n64, err := p.buf.WriteTo(w)
	p.free()
	return int(n64), err
}

// Printf formats according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
func Printf(format string, a ...interface{}) (n int, err os.Error) {
	return Fprintf(os.Stdout, format, a...)
}

// Sprintf formats according to a format specifier and returns the resulting string.
func Sprintf(format string, a ...interface{}) string {
	p := newPrinter()
	p.doPrintf(format, a)
	s := p.buf.String()
	p.free()
	return s
}

// Errorf formats according to a format specifier and returns the string 
// as a value that satisfies os.Error.
func Errorf(format string, a ...interface{}) os.Error {
	return os.NewError(Sprintf(format, a...))
}

// These routines do not take a format string

// Fprint formats using the default formats for its operands and writes to w.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func Fprint(w io.Writer, a ...interface{}) (n int, err os.Error) {
	p := newPrinter()
	p.doPrint(a, false, false)
	n64, err := p.buf.WriteTo(w)
	p.free()
	return int(n64), err
}

// Print formats using the default formats for its operands and writes to standard output.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func Print(a ...interface{}) (n int, err os.Error) {
	return Fprint(os.Stdout, a...)
}

// Sprint formats using the default formats for its operands and returns the resulting string.
// Spaces are added between operands when neither is a string.
func Sprint(a ...interface{}) string {
	p := newPrinter()
	p.doPrint(a, false, false)
	s := p.buf.String()
	p.free()
	return s
}

// These routines end in 'ln', do not take a format string,
// always add spaces between operands, and add a newline
// after the last operand.

// Fprintln formats using the default formats for its operands and writes to w.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func Fprintln(w io.Writer, a ...interface{}) (n int, err os.Error) {
	p := newPrinter()
	p.doPrint(a, true, true)
	n64, err := p.buf.WriteTo(w)
	p.free()
	return int(n64), err
}

// Println formats using the default formats for its operands and writes to standard output.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func Println(a ...interface{}) (n int, err os.Error) {
	return Fprintln(os.Stdout, a...)
}

// Sprintln formats using the default formats for its operands and returns the resulting string.
// Spaces are always added between operands and a newline is appended.
func Sprintln(a ...interface{}) string {
	p := newPrinter()
	p.doPrint(a, true, true)
	s := p.buf.String()
	p.free()
	return s
}

// Get the i'th arg of the struct value.
// If the arg itself is an interface, return a value for
// the thing inside the interface, not the interface itself.
func getField(v reflect.Value, i int) reflect.Value {
	val := v.Field(i)
	if i := val; i.Kind() == reflect.Interface {
		if inter := i.Interface(); inter != nil {
			return reflect.ValueOf(inter)
		}
	}
	return val
}

// Convert ASCII to integer.  n is 0 (and got is false) if no number present.
func parsenum(s string, start, end int) (num int, isnum bool, newi int) {
	if start >= end {
		return 0, false, end
	}
	for newi = start; newi < end && '0' <= s[newi] && s[newi] <= '9'; newi++ {
		num = num*10 + int(s[newi]-'0')
		isnum = true
	}
	return
}

func (p *pp) unknownType(v interface{}) {
	if v == nil {
		p.buf.Write(nilAngleBytes)
		return
	}
	p.buf.WriteByte('?')
	p.buf.WriteString(reflect.TypeOf(v).String())
	p.buf.WriteByte('?')
}

func (p *pp) badVerb(verb int, val interface{}) {
	p.add('%')
	p.add('!')
	p.add(verb)
	p.add('(')
	if val == nil {
		p.buf.Write(nilAngleBytes)
	} else {
		p.buf.WriteString(reflect.TypeOf(val).String())
		p.add('=')
		p.printField(val, 'v', false, false, 0)
	}
	p.add(')')
}

func (p *pp) fmtBool(v bool, verb int, value interface{}) {
	switch verb {
	case 't', 'v':
		p.fmt.fmt_boolean(v)
	default:
		p.badVerb(verb, value)
	}
}

// fmtC formats a rune for the 'c' format.
func (p *pp) fmtC(c int64) {
	rune := int(c) // Check for overflow.
	if int64(rune) != c {
		rune = utf8.RuneError
	}
	w := utf8.EncodeRune(p.runeBuf[0:utf8.UTFMax], rune)
	p.fmt.pad(p.runeBuf[0:w])
}

func (p *pp) fmtInt64(v int64, verb int, value interface{}) {
	switch verb {
	case 'b':
		p.fmt.integer(v, 2, signed, ldigits)
	case 'c':
		p.fmtC(v)
	case 'd', 'v':
		p.fmt.integer(v, 10, signed, ldigits)
	case 'o':
		p.fmt.integer(v, 8, signed, ldigits)
	case 'q':
		if 0 <= v && v <= unicode.MaxRune {
			p.fmt.fmt_qc(v)
		} else {
			p.badVerb(verb, value)
		}
	case 'x':
		p.fmt.integer(v, 16, signed, ldigits)
	case 'U':
		p.fmtUnicode(v)
	case 'X':
		p.fmt.integer(v, 16, signed, udigits)
	default:
		p.badVerb(verb, value)
	}
}

// fmt0x64 formats a uint64 in hexadecimal and prefixes it with 0x or
// not, as requested, by temporarily setting the sharp flag.
func (p *pp) fmt0x64(v uint64, leading0x bool) {
	sharp := p.fmt.sharp
	p.fmt.sharp = leading0x
	p.fmt.integer(int64(v), 16, unsigned, ldigits)
	p.fmt.sharp = sharp
}

// fmtUnicode formats a uint64 in U+1234 form by
// temporarily turning on the unicode flag and tweaking the precision.
func (p *pp) fmtUnicode(v int64) {
	precPresent := p.fmt.precPresent
	sharp := p.fmt.sharp
	p.fmt.sharp = false
	prec := p.fmt.prec
	if !precPresent {
		// If prec is already set, leave it alone; otherwise 4 is minimum.
		p.fmt.prec = 4
		p.fmt.precPresent = true
	}
	p.fmt.unicode = true // turn on U+
	p.fmt.uniQuote = sharp
	p.fmt.integer(int64(v), 16, unsigned, udigits)
	p.fmt.unicode = false
	p.fmt.uniQuote = false
	p.fmt.prec = prec
	p.fmt.precPresent = precPresent
	p.fmt.sharp = sharp
}

func (p *pp) fmtUint64(v uint64, verb int, goSyntax bool, value interface{}) {
	switch verb {
	case 'b':
		p.fmt.integer(int64(v), 2, unsigned, ldigits)
	case 'c':
		p.fmtC(int64(v))
	case 'd':
		p.fmt.integer(int64(v), 10, unsigned, ldigits)
	case 'v':
		if goSyntax {
			p.fmt0x64(v, true)
		} else {
			p.fmt.integer(int64(v), 10, unsigned, ldigits)
		}
	case 'o':
		p.fmt.integer(int64(v), 8, unsigned, ldigits)
	case 'q':
		if 0 <= v && v <= unicode.MaxRune {
			p.fmt.fmt_qc(int64(v))
		} else {
			p.badVerb(verb, value)
		}
	case 'x':
		p.fmt.integer(int64(v), 16, unsigned, ldigits)
	case 'X':
		p.fmt.integer(int64(v), 16, unsigned, udigits)
	case 'U':
		p.fmtUnicode(int64(v))
	default:
		p.badVerb(verb, value)
	}
}

func (p *pp) fmtFloat32(v float32, verb int, value interface{}) {
	switch verb {
	case 'b':
		p.fmt.fmt_fb32(v)
	case 'e':
		p.fmt.fmt_e32(v)
	case 'E':
		p.fmt.fmt_E32(v)
	case 'f':
		p.fmt.fmt_f32(v)
	case 'g', 'v':
		p.fmt.fmt_g32(v)
	case 'G':
		p.fmt.fmt_G32(v)
	default:
		p.badVerb(verb, value)
	}
}

func (p *pp) fmtFloat64(v float64, verb int, value interface{}) {
	switch verb {
	case 'b':
		p.fmt.fmt_fb64(v)
	case 'e':
		p.fmt.fmt_e64(v)
	case 'E':
		p.fmt.fmt_E64(v)
	case 'f':
		p.fmt.fmt_f64(v)
	case 'g', 'v':
		p.fmt.fmt_g64(v)
	case 'G':
		p.fmt.fmt_G64(v)
	default:
		p.badVerb(verb, value)
	}
}

func (p *pp) fmtComplex64(v complex64, verb int, value interface{}) {
	switch verb {
	case 'e', 'E', 'f', 'F', 'g', 'G':
		p.fmt.fmt_c64(v, verb)
	case 'v':
		p.fmt.fmt_c64(v, 'g')
	default:
		p.badVerb(verb, value)
	}
}

func (p *pp) fmtComplex128(v complex128, verb int, value interface{}) {
	switch verb {
	case 'e', 'E', 'f', 'F', 'g', 'G':
		p.fmt.fmt_c128(v, verb)
	case 'v':
		p.fmt.fmt_c128(v, 'g')
	default:
		p.badVerb(verb, value)
	}
}

func (p *pp) fmtString(v string, verb int, goSyntax bool, value interface{}) {
	switch verb {
	case 'v':
		if goSyntax {
			p.fmt.fmt_q(v)
		} else {
			p.fmt.fmt_s(v)
		}
	case 's':
		p.fmt.fmt_s(v)
	case 'x':
		p.fmt.fmt_sx(v)
	case 'X':
		p.fmt.fmt_sX(v)
	case 'q':
		p.fmt.fmt_q(v)
	default:
		p.badVerb(verb, value)
	}
}

func (p *pp) fmtBytes(v []byte, verb int, goSyntax bool, depth int, value interface{}) {
	if verb == 'v' || verb == 'd' {
		if goSyntax {
			p.buf.Write(bytesBytes)
		} else {
			p.buf.WriteByte('[')
		}
		for i, c := range v {
			if i > 0 {
				if goSyntax {
					p.buf.Write(commaSpaceBytes)
				} else {
					p.buf.WriteByte(' ')
				}
			}
			p.printField(c, 'v', p.fmt.plus, goSyntax, depth+1)
		}
		if goSyntax {
			p.buf.WriteByte('}')
		} else {
			p.buf.WriteByte(']')
		}
		return
	}
	s := string(v)
	switch verb {
	case 's':
		p.fmt.fmt_s(s)
	case 'x':
		p.fmt.fmt_sx(s)
	case 'X':
		p.fmt.fmt_sX(s)
	case 'q':
		p.fmt.fmt_q(s)
	default:
		p.badVerb(verb, value)
	}
}

func (p *pp) fmtPointer(field interface{}, value reflect.Value, verb int, goSyntax bool) {
	var u uintptr
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		u = value.Pointer()
	default:
		p.badVerb(verb, field)
		return
	}
	if goSyntax {
		p.add('(')
		p.buf.WriteString(reflect.TypeOf(field).String())
		p.add(')')
		p.add('(')
		if u == 0 {
			p.buf.Write(nilBytes)
		} else {
			p.fmt0x64(uint64(u), true)
		}
		p.add(')')
	} else {
		p.fmt0x64(uint64(u), !p.fmt.sharp)
	}
}

var (
	intBits     = reflect.TypeOf(0).Bits()
	floatBits   = reflect.TypeOf(0.0).Bits()
	complexBits = reflect.TypeOf(1i).Bits()
	uintptrBits = reflect.TypeOf(uintptr(0)).Bits()
)

func (p *pp) catchPanic(val interface{}, verb int) {
	if err := recover(); err != nil {
		// If it's a nil pointer, just say "<nil>". The likeliest causes are a
		// Stringer that fails to guard against nil or a nil pointer for a
		// value receiver, and in either case, "<nil>" is a nice result.
		if v := reflect.ValueOf(val); v.Kind() == reflect.Ptr && v.IsNil() {
			p.buf.Write(nilAngleBytes)
			return
		}
		// Otherwise print a concise panic message. Most of the time the panic
		// value will print itself nicely.
		if p.panicking {
			// Nested panics; the recursion in printField cannot succeed.
			panic(err)
		}
		p.buf.WriteByte('%')
		p.add(verb)
		p.buf.Write(panicBytes)
		p.panicking = true
		p.printField(err, 'v', false, false, 0)
		p.panicking = false
		p.buf.WriteByte(')')
	}
}

func (p *pp) printField(field interface{}, verb int, plus, goSyntax bool, depth int) (wasString bool) {
	if field == nil {
		if verb == 'T' || verb == 'v' {
			p.buf.Write(nilAngleBytes)
		} else {
			p.badVerb(verb, field)
		}
		return false
	}

	// Special processing considerations.
	// %T (the value's type) and %p (its address) are special; we always do them first.
	switch verb {
	case 'T':
		p.printField(reflect.TypeOf(field).String(), 's', false, false, 0)
		return false
	case 'p':
		p.fmtPointer(field, reflect.ValueOf(field), verb, goSyntax)
		return false
	}
	// Is it a Formatter?
	if formatter, ok := field.(Formatter); ok {
		defer p.catchPanic(field, verb)
		formatter.Format(p, verb)
		return false // this value is not a string

	}
	// Must not touch flags before Formatter looks at them.
	if plus {
		p.fmt.plus = false
	}
	// If we're doing Go syntax and the field knows how to supply it, take care of it now.
	if goSyntax {
		p.fmt.sharp = false
		if stringer, ok := field.(GoStringer); ok {
			defer p.catchPanic(field, verb)
			// Print the result of GoString unadorned.
			p.fmtString(stringer.GoString(), 's', false, field)
			return false // this value is not a string
		}
	} else {
		// Is it a Stringer?
		if stringer, ok := field.(Stringer); ok {
			defer p.catchPanic(field, verb)
			p.printField(stringer.String(), verb, plus, false, depth)
			return false // this value is not a string
		}
	}

	// Some types can be done without reflection.
	switch f := field.(type) {
	case bool:
		p.fmtBool(f, verb, field)
		return false
	case float32:
		p.fmtFloat32(f, verb, field)
		return false
	case float64:
		p.fmtFloat64(f, verb, field)
		return false
	case complex64:
		p.fmtComplex64(complex64(f), verb, field)
		return false
	case complex128:
		p.fmtComplex128(f, verb, field)
		return false
	case int:
		p.fmtInt64(int64(f), verb, field)
		return false
	case int8:
		p.fmtInt64(int64(f), verb, field)
		return false
	case int16:
		p.fmtInt64(int64(f), verb, field)
		return false
	case int32:
		p.fmtInt64(int64(f), verb, field)
		return false
	case int64:
		p.fmtInt64(f, verb, field)
		return false
	case uint:
		p.fmtUint64(uint64(f), verb, goSyntax, field)
		return false
	case uint8:
		p.fmtUint64(uint64(f), verb, goSyntax, field)
		return false
	case uint16:
		p.fmtUint64(uint64(f), verb, goSyntax, field)
		return false
	case uint32:
		p.fmtUint64(uint64(f), verb, goSyntax, field)
		return false
	case uint64:
		p.fmtUint64(f, verb, goSyntax, field)
		return false
	case uintptr:
		p.fmtUint64(uint64(f), verb, goSyntax, field)
		return false
	case string:
		p.fmtString(f, verb, goSyntax, field)
		return verb == 's' || verb == 'v'
	case []byte:
		p.fmtBytes(f, verb, goSyntax, depth, field)
		return verb == 's'
	}

	// Need to use reflection
	value := reflect.ValueOf(field)

BigSwitch:
	switch f := value; f.Kind() {
	case reflect.Bool:
		p.fmtBool(f.Bool(), verb, field)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p.fmtInt64(f.Int(), verb, field)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p.fmtUint64(uint64(f.Uint()), verb, goSyntax, field)
	case reflect.Float32, reflect.Float64:
		if f.Type().Size() == 4 {
			p.fmtFloat32(float32(f.Float()), verb, field)
		} else {
			p.fmtFloat64(float64(f.Float()), verb, field)
		}
	case reflect.Complex64, reflect.Complex128:
		if f.Type().Size() == 8 {
			p.fmtComplex64(complex64(f.Complex()), verb, field)
		} else {
			p.fmtComplex128(complex128(f.Complex()), verb, field)
		}
	case reflect.String:
		p.fmtString(f.String(), verb, goSyntax, field)
	case reflect.Map:
		if goSyntax {
			p.buf.WriteString(f.Type().String())
			p.buf.WriteByte('{')
		} else {
			p.buf.Write(mapBytes)
		}
		keys := f.MapKeys()
		for i, key := range keys {
			if i > 0 {
				if goSyntax {
					p.buf.Write(commaSpaceBytes)
				} else {
					p.buf.WriteByte(' ')
				}
			}
			p.printField(key.Interface(), verb, plus, goSyntax, depth+1)
			p.buf.WriteByte(':')
			p.printField(f.MapIndex(key).Interface(), verb, plus, goSyntax, depth+1)
		}
		if goSyntax {
			p.buf.WriteByte('}')
		} else {
			p.buf.WriteByte(']')
		}
	case reflect.Struct:
		if goSyntax {
			p.buf.WriteString(reflect.TypeOf(field).String())
		}
		p.add('{')
		v := f
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				if goSyntax {
					p.buf.Write(commaSpaceBytes)
				} else {
					p.buf.WriteByte(' ')
				}
			}
			if plus || goSyntax {
				if f := t.Field(i); f.Name != "" {
					p.buf.WriteString(f.Name)
					p.buf.WriteByte(':')
				}
			}
			p.printField(getField(v, i).Interface(), verb, plus, goSyntax, depth+1)
		}
		p.buf.WriteByte('}')
	case reflect.Interface:
		value := f.Elem()
		if !value.IsValid() {
			if goSyntax {
				p.buf.WriteString(reflect.TypeOf(field).String())
				p.buf.Write(nilParenBytes)
			} else {
				p.buf.Write(nilAngleBytes)
			}
		} else {
			return p.printField(value.Interface(), verb, plus, goSyntax, depth+1)
		}
	case reflect.Array, reflect.Slice:
		// Byte slices are special.
		if f.Type().Elem().Kind() == reflect.Uint8 {
			// We know it's a slice of bytes, but we also know it does not have static type
			// []byte, or it would have been caught above.  Therefore we cannot convert
			// it directly in the (slightly) obvious way: f.Interface().([]byte); it doesn't have
			// that type, and we can't write an expression of the right type and do a
			// conversion because we don't have a static way to write the right type.
			// So we build a slice by hand.  This is a rare case but it would be nice
			// if reflection could help a little more.
			bytes := make([]byte, f.Len())
			for i := range bytes {
				bytes[i] = byte(f.Index(i).Uint())
			}
			p.fmtBytes(bytes, verb, goSyntax, depth, field)
			return verb == 's'
		}
		if goSyntax {
			p.buf.WriteString(reflect.TypeOf(field).String())
			p.buf.WriteByte('{')
		} else {
			p.buf.WriteByte('[')
		}
		for i := 0; i < f.Len(); i++ {
			if i > 0 {
				if goSyntax {
					p.buf.Write(commaSpaceBytes)
				} else {
					p.buf.WriteByte(' ')
				}
			}
			p.printField(f.Index(i).Interface(), verb, plus, goSyntax, depth+1)
		}
		if goSyntax {
			p.buf.WriteByte('}')
		} else {
			p.buf.WriteByte(']')
		}
	case reflect.Ptr:
		v := f.Pointer()
		// pointer to array or slice or struct?  ok at top level
		// but not embedded (avoid loops)
		if v != 0 && depth == 0 {
			switch a := f.Elem(); a.Kind() {
			case reflect.Array, reflect.Slice:
				p.buf.WriteByte('&')
				p.printField(a.Interface(), verb, plus, goSyntax, depth+1)
				break BigSwitch
			case reflect.Struct:
				p.buf.WriteByte('&')
				p.printField(a.Interface(), verb, plus, goSyntax, depth+1)
				break BigSwitch
			}
		}
		if goSyntax {
			p.buf.WriteByte('(')
			p.buf.WriteString(reflect.TypeOf(field).String())
			p.buf.WriteByte(')')
			p.buf.WriteByte('(')
			if v == 0 {
				p.buf.Write(nilBytes)
			} else {
				p.fmt0x64(uint64(v), true)
			}
			p.buf.WriteByte(')')
			break
		}
		if v == 0 {
			p.buf.Write(nilAngleBytes)
			break
		}
		p.fmt0x64(uint64(v), true)
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		p.fmtPointer(field, value, verb, goSyntax)
	default:
		p.unknownType(f)
	}
	return false
}

// intFromArg gets the fieldnumth element of a. On return, isInt reports whether the argument has type int.
func intFromArg(a []interface{}, end, i, fieldnum int) (num int, isInt bool, newi, newfieldnum int) {
	newi, newfieldnum = end, fieldnum
	if i < end && fieldnum < len(a) {
		num, isInt = a[fieldnum].(int)
		newi, newfieldnum = i+1, fieldnum+1
	}
	return
}

func (p *pp) doPrintf(format string, a []interface{}) {
	end := len(format)
	fieldnum := 0 // we process one field per non-trivial format
	for i := 0; i < end; {
		lasti := i
		for i < end && format[i] != '%' {
			i++
		}
		if i > lasti {
			p.buf.WriteString(format[lasti:i])
		}
		if i >= end {
			// done processing format string
			break
		}

		// Process one verb
		i++
		// flags and widths
		p.fmt.clearflags()
	F:
		for ; i < end; i++ {
			switch format[i] {
			case '#':
				p.fmt.sharp = true
			case '0':
				p.fmt.zero = true
			case '+':
				p.fmt.plus = true
			case '-':
				p.fmt.minus = true
			case ' ':
				p.fmt.space = true
			default:
				break F
			}
		}
		// do we have width?
		if i < end && format[i] == '*' {
			p.fmt.wid, p.fmt.widPresent, i, fieldnum = intFromArg(a, end, i, fieldnum)
			if !p.fmt.widPresent {
				p.buf.Write(widthBytes)
			}
		} else {
			p.fmt.wid, p.fmt.widPresent, i = parsenum(format, i, end)
		}
		// do we have precision?
		if i < end && format[i] == '.' {
			if format[i+1] == '*' {
				p.fmt.prec, p.fmt.precPresent, i, fieldnum = intFromArg(a, end, i+1, fieldnum)
				if !p.fmt.precPresent {
					p.buf.Write(precBytes)
				}
			} else {
				p.fmt.prec, p.fmt.precPresent, i = parsenum(format, i+1, end)
				if !p.fmt.precPresent {
					p.fmt.prec = 0
					p.fmt.precPresent = true
				}
			}
		}
		if i >= end {
			p.buf.Write(noVerbBytes)
			continue
		}
		c, w := utf8.DecodeRuneInString(format[i:])
		i += w
		// percent is special - absorbs no operand
		if c == '%' {
			p.buf.WriteByte('%') // We ignore width and prec.
			continue
		}
		if fieldnum >= len(a) { // out of operands
			p.buf.WriteByte('%')
			p.add(c)
			p.buf.Write(missingBytes)
			continue
		}
		field := a[fieldnum]
		fieldnum++

		goSyntax := c == 'v' && p.fmt.sharp
		plus := c == 'v' && p.fmt.plus
		p.printField(field, c, plus, goSyntax, 0)
	}

	if fieldnum < len(a) {
		p.buf.Write(extraBytes)
		for ; fieldnum < len(a); fieldnum++ {
			field := a[fieldnum]
			if field != nil {
				p.buf.WriteString(reflect.TypeOf(field).String())
				p.buf.WriteByte('=')
			}
			p.printField(field, 'v', false, false, 0)
			if fieldnum+1 < len(a) {
				p.buf.Write(commaSpaceBytes)
			}
		}
		p.buf.WriteByte(')')
	}
}

func (p *pp) doPrint(a []interface{}, addspace, addnewline bool) {
	prevString := false
	for fieldnum := 0; fieldnum < len(a); fieldnum++ {
		p.fmt.clearflags()
		// always add spaces if we're doing println
		field := a[fieldnum]
		if fieldnum > 0 {
			isString := field != nil && reflect.TypeOf(field).Kind() == reflect.String
			if addspace || !isString && !prevString {
				p.buf.WriteByte(' ')
			}
		}
		prevString = p.printField(field, 'v', false, false, 0)
	}
	if addnewline {
		p.buf.WriteByte('\n')
	}
}
