// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package log implements a simple logging package. It defines a type, Logger,
// with methods for formatting output. It also has a predefined 'standard'
// Logger accessible through helper functions Print[f|ln], Fatal[f|ln], and
// Panic[f|ln], which are easier to use than creating a Logger manually.
// That logger writes to standard error and prints the date and time
// of each logged message.
// Every log message is output on a separate line: if the message being
// printed does not end in a newline, the logger will add one.
// The Fatal functions call os.Exit(1) after writing the log message.
// The Panic functions call panic after writing the log message.
package log

import (
	"fmt"
	"io"
	"log/internal"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// These flags define which text to prefix to each log entry generated by the Logger.
// Bits are or'ed together to control what's printed.
// With the exception of the Lmsgprefix flag, there is no
// control over the order they appear (the order listed here)
// or the format they present (as described in the comments).
// The prefix is followed by a colon only when Llongfile or Lshortfile
// is specified.
// For example, flags Ldate | Ltime (or LstdFlags) produce,
//
//	2009/01/23 01:23:23 message
//
// while flags Ldate | Ltime | Lmicroseconds | Llongfile produce,
//
//	2009/01/23 01:23:23.123123 /a/b/c/d.go:23: message
const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	Lmsgprefix                    // move the "prefix" from the beginning of the line to before the message
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

// A Logger represents an active logging object that generates lines of
// output to an io.Writer. Each logging operation makes a single call to
// the Writer's Write method. A Logger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type Logger struct {
	outMu sync.Mutex
	out   io.Writer // destination for output

	prefix    atomic.Pointer[string] // prefix on each line to identify the logger (but see Lmsgprefix)
	flag      atomic.Int32           // properties
	isDiscard atomic.Bool
}

// New creates a new Logger. The out variable sets the
// destination to which log data will be written.
// The prefix appears at the beginning of each generated log line, or
// after the log header if the Lmsgprefix flag is provided.
// The flag argument defines the logging properties.
func New(out io.Writer, prefix string, flag int) *Logger {
	l := new(Logger)
	l.SetOutput(out)
	l.SetPrefix(prefix)
	l.SetFlags(flag)
	return l
}

// SetOutput sets the output destination for the logger.
func (l *Logger) SetOutput(w io.Writer) {
	l.outMu.Lock()
	defer l.outMu.Unlock()
	l.out = w
	l.isDiscard.Store(w == io.Discard)
}

var std = New(os.Stderr, "", LstdFlags)

// Default returns the standard logger used by the package-level output functions.
func Default() *Logger { return std }

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

// formatHeader writes log header to buf in following order:
//   - l.prefix (if it's not blank and Lmsgprefix is unset),
//   - date and/or time (if corresponding flags are provided),
//   - file and line number (if corresponding flags are provided),
//   - l.prefix (if it's not blank and Lmsgprefix is set).
func formatHeader(buf *[]byte, t time.Time, prefix string, flag int, file string, line int) {
	if flag&Lmsgprefix == 0 {
		*buf = append(*buf, prefix...)
	}
	if flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if flag&LUTC != 0 {
			t = t.UTC()
		}
		if flag&Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if flag&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}
	if flag&(Lshortfile|Llongfile) != 0 {
		if flag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
	if flag&Lmsgprefix != 0 {
		*buf = append(*buf, prefix...)
	}
}

var bufferPool = sync.Pool{New: func() any { return new([]byte) }}

func getBuffer() *[]byte {
	p := bufferPool.Get().(*[]byte)
	*p = (*p)[:0]
	return p
}

func putBuffer(p *[]byte) {
	// Proper usage of a sync.Pool requires each entry to have approximately
	// the same memory cost. To obtain this property when the stored type
	// contains a variably-sized buffer, we add a hard limit on the maximum buffer
	// to place back in the pool.
	//
	// See https://go.dev/issue/23199
	if cap(*p) > 64<<10 {
		*p = nil
	}
	bufferPool.Put(p)
}

// Output writes the output for a logging event. The string s contains
// the text to print after the prefix specified by the flags of the
// Logger. A newline is appended if the last character of s is not
// already a newline. Calldepth is used to recover the PC and is
// provided for generality, although at the moment on all pre-defined
// paths it will be 2.
func (l *Logger) Output(calldepth int, s string) error {
	calldepth++ // +1 for this frame.
	return l.output(0, calldepth, func(b []byte) []byte {
		return append(b, s...)
	})
}

// output can take either a calldepth or a pc to get source line information.
// It uses the pc if it is non-zero.
func (l *Logger) output(pc uintptr, calldepth int, appendOutput func([]byte) []byte) error {
	if l.isDiscard.Load() {
		return nil
	}

	now := time.Now() // get this early.

	// Load prefix and flag once so that their value is consistent within
	// this call regardless of any concurrent changes to their value.
	prefix := l.Prefix()
	flag := l.Flags()

	var file string
	var line int
	if flag&(Lshortfile|Llongfile) != 0 {
		if pc == 0 {
			var ok bool
			_, file, line, ok = runtime.Caller(calldepth)
			if !ok {
				file = "???"
				line = 0
			}
		} else {
			fs := runtime.CallersFrames([]uintptr{pc})
			f, _ := fs.Next()
			file = f.File
			if file == "" {
				file = "???"
			}
			line = f.Line
		}
	}

	buf := getBuffer()
	defer putBuffer(buf)
	formatHeader(buf, now, prefix, flag, file, line)
	*buf = appendOutput(*buf)
	if len(*buf) == 0 || (*buf)[len(*buf)-1] != '\n' {
		*buf = append(*buf, '\n')
	}

	l.outMu.Lock()
	defer l.outMu.Unlock()
	_, err := l.out.Write(*buf)
	return err
}

func init() {
	internal.DefaultOutput = func(pc uintptr, data []byte) error {
		return std.output(pc, 0, func(buf []byte) []byte {
			return append(buf, data...)
		})
	}
}

// Print calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Print(v ...any) {
	l.output(0, 2, func(b []byte) []byte {
		return fmt.Append(b, v...)
	})
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Printf(format string, v ...any) {
	l.output(0, 2, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	})
}

// Println calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Println(v ...any) {
	l.output(0, 2, func(b []byte) []byte {
		return fmt.Appendln(b, v...)
	})
}

// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
func (l *Logger) Fatal(v ...any) {
	l.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (l *Logger) Fatalf(format string, v ...any) {
	l.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln is equivalent to l.Println() followed by a call to os.Exit(1).
func (l *Logger) Fatalln(v ...any) {
	l.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}

// Panic is equivalent to l.Print() followed by a call to panic().
func (l *Logger) Panic(v ...any) {
	s := fmt.Sprint(v...)
	l.Output(2, s)
	panic(s)
}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func (l *Logger) Panicf(format string, v ...any) {
	s := fmt.Sprintf(format, v...)
	l.Output(2, s)
	panic(s)
}

// Panicln is equivalent to l.Println() followed by a call to panic().
func (l *Logger) Panicln(v ...any) {
	s := fmt.Sprintln(v...)
	l.Output(2, s)
	panic(s)
}

// Flags returns the output flags for the logger.
// The flag bits are Ldate, Ltime, and so on.
func (l *Logger) Flags() int {
	return int(l.flag.Load())
}

// SetFlags sets the output flags for the logger.
// The flag bits are Ldate, Ltime, and so on.
func (l *Logger) SetFlags(flag int) {
	l.flag.Store(int32(flag))
}

// Prefix returns the output prefix for the logger.
func (l *Logger) Prefix() string {
	if p := l.prefix.Load(); p != nil {
		return *p
	}
	return ""
}

// SetPrefix sets the output prefix for the logger.
func (l *Logger) SetPrefix(prefix string) {
	l.prefix.Store(&prefix)
}

// Writer returns the output destination for the logger.
func (l *Logger) Writer() io.Writer {
	l.outMu.Lock()
	defer l.outMu.Unlock()
	return l.out
}

// SetOutput sets the output destination for the standard logger.
func SetOutput(w io.Writer) {
	std.SetOutput(w)
}

// Flags returns the output flags for the standard logger.
// The flag bits are Ldate, Ltime, and so on.
func Flags() int {
	return std.Flags()
}

// SetFlags sets the output flags for the standard logger.
// The flag bits are Ldate, Ltime, and so on.
func SetFlags(flag int) {
	std.SetFlags(flag)
}

// Prefix returns the output prefix for the standard logger.
func Prefix() string {
	return std.Prefix()
}

// SetPrefix sets the output prefix for the standard logger.
func SetPrefix(prefix string) {
	std.SetPrefix(prefix)
}

// Writer returns the output destination for the standard logger.
func Writer() io.Writer {
	return std.Writer()
}

// These functions write to the standard logger.

// Print calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...any) {
	std.output(0, 2, func(b []byte) []byte {
		return fmt.Append(b, v...)
	})
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...any) {
	std.output(0, 2, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	})
}

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...any) {
	std.output(0, 2, func(b []byte) []byte {
		return fmt.Appendln(b, v...)
	})
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(v ...any) {
	std.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func Fatalf(format string, v ...any) {
	std.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln is equivalent to Println() followed by a call to os.Exit(1).
func Fatalln(v ...any) {
	std.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}

// Panic is equivalent to Print() followed by a call to panic().
func Panic(v ...any) {
	s := fmt.Sprint(v...)
	std.Output(2, s)
	panic(s)
}

// Panicf is equivalent to Printf() followed by a call to panic().
func Panicf(format string, v ...any) {
	s := fmt.Sprintf(format, v...)
	std.Output(2, s)
	panic(s)
}

// Panicln is equivalent to Println() followed by a call to panic().
func Panicln(v ...any) {
	s := fmt.Sprintln(v...)
	std.Output(2, s)
	panic(s)
}

// Output writes the output for a logging event. The string s contains
// the text to print after the prefix specified by the flags of the
// Logger. A newline is appended if the last character of s is not
// already a newline. Calldepth is the count of the number of
// frames to skip when computing the file name and line number
// if Llongfile or Lshortfile is set; a value of 1 will print the details
// for the caller of Output.
func Output(calldepth int, s string) error {
	return std.Output(calldepth+1, s) // +1 for this frame.
}
