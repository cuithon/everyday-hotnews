// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"bytes"
	"context"
	"internal/race"
	"internal/testenv"
	"io"
	"log"
	loginternal "log/internal"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

const timeRE = `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}(Z|[+-]\d{2}:\d{2})`

func TestLogTextHandler(t *testing.T) {
	var buf bytes.Buffer

	l := New(NewTextHandler(&buf))

	check := func(want string) {
		t.Helper()
		if want != "" {
			want = "time=" + timeRE + " " + want
		}
		checkLogOutput(t, buf.String(), want)
		buf.Reset()
	}

	l.Info("msg", "a", 1, "b", 2)
	check(`level=INFO msg=msg a=1 b=2`)

	// By default, debug messages are not printed.
	l.Debug("bg", Int("a", 1), "b", 2)
	check("")

	l.Warn("w", Duration("dur", 3*time.Second))
	check(`level=WARN msg=w dur=3s`)

	l.Error("bad", "a", 1)
	check(`level=ERROR msg=bad a=1`)

	l.Log(nil, LevelWarn+1, "w", Int("a", 1), String("b", "two"))
	check(`level=WARN\+1 msg=w a=1 b=two`)

	l.LogAttrs(nil, LevelInfo+1, "a b c", Int("a", 1), String("b", "two"))
	check(`level=INFO\+1 msg="a b c" a=1 b=two`)

	l.Info("info", "a", []Attr{Int("i", 1)})
	check(`level=INFO msg=info a.i=1`)

	l.Info("info", "a", GroupValue(Int("i", 1)))
	check(`level=INFO msg=info a.i=1`)
}

func TestConnections(t *testing.T) {
	var logbuf, slogbuf bytes.Buffer

	// Revert any changes to the default logger. This is important because other
	// tests might change the default logger using SetDefault. Also ensure we
	// restore the default logger at the end of the test.
	currentLogger := Default()
	SetDefault(New(newDefaultHandler(loginternal.DefaultOutput)))
	t.Cleanup(func() {
		SetDefault(currentLogger)
	})

	// The default slog.Logger's handler uses the log package's default output.
	log.SetOutput(&logbuf)
	log.SetFlags(log.Lshortfile &^ log.LstdFlags)
	Info("msg", "a", 1)
	checkLogOutput(t, logbuf.String(), `logger_test.go:\d+: INFO msg a=1`)
	logbuf.Reset()
	Warn("msg", "b", 2)
	checkLogOutput(t, logbuf.String(), `logger_test.go:\d+: WARN msg b=2`)
	logbuf.Reset()
	Error("msg", "err", io.EOF, "c", 3)
	checkLogOutput(t, logbuf.String(), `logger_test.go:\d+: ERROR msg err=EOF c=3`)

	// Levels below Info are not printed.
	logbuf.Reset()
	Debug("msg", "c", 3)
	checkLogOutput(t, logbuf.String(), "")

	t.Run("wrap default handler", func(t *testing.T) {
		// It should be possible to wrap the default handler and get the right output.
		// This works because the default handler uses the pc in the Record
		// to get the source line, rather than a call depth.
		logger := New(wrappingHandler{Default().Handler()})
		logger.Info("msg", "d", 4)
		checkLogOutput(t, logbuf.String(), `logger_test.go:\d+: INFO msg d=4`)
	})

	// Once slog.SetDefault is called, the direction is reversed: the default
	// log.Logger's output goes through the handler.
	SetDefault(New(HandlerOptions{AddSource: true}.NewTextHandler(&slogbuf)))
	log.Print("msg2")
	checkLogOutput(t, slogbuf.String(), "time="+timeRE+` level=INFO source=.*logger_test.go:\d{3} msg=msg2`)

	// The default log.Logger always outputs at Info level.
	slogbuf.Reset()
	SetDefault(New(HandlerOptions{Level: LevelWarn}.NewTextHandler(&slogbuf)))
	log.Print("should not appear")
	if got := slogbuf.String(); got != "" {
		t.Errorf("got %q, want empty", got)
	}

	// Setting log's output again breaks the connection.
	logbuf.Reset()
	slogbuf.Reset()
	log.SetOutput(&logbuf)
	log.SetFlags(log.Lshortfile &^ log.LstdFlags)
	log.Print("msg3")
	checkLogOutput(t, logbuf.String(), `logger_test.go:\d+: msg3`)
	if got := slogbuf.String(); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

type wrappingHandler struct {
	h Handler
}

func (h wrappingHandler) Enabled(ctx context.Context, level Level) bool {
	return h.h.Enabled(ctx, level)
}
func (h wrappingHandler) WithGroup(name string) Handler              { return h.h.WithGroup(name) }
func (h wrappingHandler) WithAttrs(as []Attr) Handler                { return h.h.WithAttrs(as) }
func (h wrappingHandler) Handle(ctx context.Context, r Record) error { return h.h.Handle(ctx, r) }

func TestAttrs(t *testing.T) {
	check := func(got []Attr, want ...Attr) {
		t.Helper()
		if !attrsEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	}

	l1 := New(&captureHandler{}).With("a", 1)
	l2 := New(l1.Handler()).With("b", 2)
	l2.Info("m", "c", 3)
	h := l2.Handler().(*captureHandler)
	check(h.attrs, Int("a", 1), Int("b", 2))
	check(attrsSlice(h.r), Int("c", 3))
}

func sourceLine(r Record) (string, int) {
	f := r.frame()
	return f.File, f.Line
}

func TestCallDepth(t *testing.T) {
	h := &captureHandler{}
	var startLine int

	check := func(count int) {
		t.Helper()
		const wantFile = "logger_test.go"
		wantLine := startLine + count*2
		gotFile, gotLine := sourceLine(h.r)
		gotFile = filepath.Base(gotFile)
		if gotFile != wantFile || gotLine != wantLine {
			t.Errorf("got (%s, %d), want (%s, %d)", gotFile, gotLine, wantFile, wantLine)
		}
	}

	logger := New(h)
	SetDefault(logger)

	// Calls to check must be one line apart.
	// Determine line where calls start.
	f, _ := runtime.CallersFrames([]uintptr{callerPC(2)}).Next()
	startLine = f.Line + 4
	// Do not change the number of lines between here and the call to check(0).

	logger.Log(nil, LevelInfo, "")
	check(0)
	logger.LogAttrs(nil, LevelInfo, "")
	check(1)
	logger.Debug("")
	check(2)
	logger.Info("")
	check(3)
	logger.Warn("")
	check(4)
	logger.Error("")
	check(5)
	Debug("")
	check(6)
	Info("")
	check(7)
	Warn("")
	check(8)
	Error("")
	check(9)
	Log(nil, LevelInfo, "")
	check(10)
	LogAttrs(nil, LevelInfo, "")
	check(11)
}

func TestAlloc(t *testing.T) {
	dl := New(discardHandler{})
	defer func(d *Logger) { SetDefault(d) }(Default())
	SetDefault(dl)

	t.Run("Info", func(t *testing.T) {
		wantAllocs(t, 0, func() { Info("hello") })
	})
	t.Run("Error", func(t *testing.T) {
		wantAllocs(t, 0, func() { Error("hello") })
	})
	t.Run("logger.Info", func(t *testing.T) {
		wantAllocs(t, 0, func() { dl.Info("hello") })
	})
	t.Run("logger.Log", func(t *testing.T) {
		wantAllocs(t, 0, func() { dl.Log(nil, LevelDebug, "hello") })
	})
	t.Run("2 pairs", func(t *testing.T) {
		s := "abc"
		i := 2000
		wantAllocs(t, 2, func() {
			dl.Info("hello",
				"n", i,
				"s", s,
			)
		})
	})
	t.Run("2 pairs disabled inline", func(t *testing.T) {
		l := New(discardHandler{disabled: true})
		s := "abc"
		i := 2000
		wantAllocs(t, 2, func() {
			l.Log(nil, LevelInfo, "hello",
				"n", i,
				"s", s,
			)
		})
	})
	t.Run("2 pairs disabled", func(t *testing.T) {
		l := New(discardHandler{disabled: true})
		s := "abc"
		i := 2000
		wantAllocs(t, 0, func() {
			if l.Enabled(nil, LevelInfo) {
				l.Log(nil, LevelInfo, "hello",
					"n", i,
					"s", s,
				)
			}
		})
	})
	t.Run("9 kvs", func(t *testing.T) {
		s := "abc"
		i := 2000
		d := time.Second
		wantAllocs(t, 11, func() {
			dl.Info("hello",
				"n", i, "s", s, "d", d,
				"n", i, "s", s, "d", d,
				"n", i, "s", s, "d", d)
		})
	})
	t.Run("pairs", func(t *testing.T) {
		wantAllocs(t, 0, func() { dl.Info("", "error", io.EOF) })
	})
	t.Run("attrs1", func(t *testing.T) {
		wantAllocs(t, 0, func() { dl.LogAttrs(nil, LevelInfo, "", Int("a", 1)) })
		wantAllocs(t, 0, func() { dl.LogAttrs(nil, LevelInfo, "", Any("error", io.EOF)) })
	})
	t.Run("attrs3", func(t *testing.T) {
		wantAllocs(t, 0, func() {
			dl.LogAttrs(nil, LevelInfo, "hello", Int("a", 1), String("b", "two"), Duration("c", time.Second))
		})
	})
	t.Run("attrs3 disabled", func(t *testing.T) {
		logger := New(discardHandler{disabled: true})
		wantAllocs(t, 0, func() {
			logger.LogAttrs(nil, LevelInfo, "hello", Int("a", 1), String("b", "two"), Duration("c", time.Second))
		})
	})
	t.Run("attrs6", func(t *testing.T) {
		wantAllocs(t, 1, func() {
			dl.LogAttrs(nil, LevelInfo, "hello",
				Int("a", 1), String("b", "two"), Duration("c", time.Second),
				Int("d", 1), String("e", "two"), Duration("f", time.Second))
		})
	})
	t.Run("attrs9", func(t *testing.T) {
		wantAllocs(t, 1, func() {
			dl.LogAttrs(nil, LevelInfo, "hello",
				Int("a", 1), String("b", "two"), Duration("c", time.Second),
				Int("d", 1), String("e", "two"), Duration("f", time.Second),
				Int("d", 1), String("e", "two"), Duration("f", time.Second))
		})
	})
}

func TestSetAttrs(t *testing.T) {
	for _, test := range []struct {
		args []any
		want []Attr
	}{
		{nil, nil},
		{[]any{"a", 1}, []Attr{Int("a", 1)}},
		{[]any{"a", 1, "b", "two"}, []Attr{Int("a", 1), String("b", "two")}},
		{[]any{"a"}, []Attr{String(badKey, "a")}},
		{[]any{"a", 1, "b"}, []Attr{Int("a", 1), String(badKey, "b")}},
		{[]any{"a", 1, 2, 3}, []Attr{Int("a", 1), Int(badKey, 2), Int(badKey, 3)}},
	} {
		r := NewRecord(time.Time{}, 0, "", 0)
		r.Add(test.args...)
		got := attrsSlice(r)
		if !attrsEqual(got, test.want) {
			t.Errorf("%v:\ngot  %v\nwant %v", test.args, got, test.want)
		}
	}
}

func TestSetDefault(t *testing.T) {
	// Verify that setting the default to itself does not result in deadlock.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	defer func(w io.Writer) { log.SetOutput(w) }(log.Writer())
	log.SetOutput(io.Discard)
	go func() {
		Info("A")
		SetDefault(Default())
		Info("B")
		cancel()
	}()
	<-ctx.Done()
	if err := ctx.Err(); err != context.Canceled {
		t.Errorf("wanted canceled, got %v", err)
	}
}

func TestLoggerError(t *testing.T) {
	var buf bytes.Buffer

	removeTime := func(_ []string, a Attr) Attr {
		if a.Key == TimeKey {
			return Attr{}
		}
		return a
	}
	l := New(HandlerOptions{ReplaceAttr: removeTime}.NewTextHandler(&buf))
	l.Error("msg", "err", io.EOF, "a", 1)
	checkLogOutput(t, buf.String(), `level=ERROR msg=msg err=EOF a=1`)
	buf.Reset()
	l.Error("msg", "err", io.EOF, "a")
	checkLogOutput(t, buf.String(), `level=ERROR msg=msg err=EOF !BADKEY=a`)
}

func TestNewLogLogger(t *testing.T) {
	var buf bytes.Buffer
	h := NewTextHandler(&buf)
	ll := NewLogLogger(h, LevelWarn)
	ll.Print("hello")
	checkLogOutput(t, buf.String(), "time="+timeRE+` level=WARN msg=hello`)
}

func checkLogOutput(t *testing.T, got, wantRegexp string) {
	t.Helper()
	got = clean(got)
	wantRegexp = "^" + wantRegexp + "$"
	matched, err := regexp.MatchString(wantRegexp, got)
	if err != nil {
		t.Fatal(err)
	}
	if !matched {
		t.Errorf("\ngot  %s\nwant %s", got, wantRegexp)
	}
}

// clean prepares log output for comparison.
func clean(s string) string {
	if len(s) > 0 && s[len(s)-1] == '\n' {
		s = s[:len(s)-1]
	}
	return strings.ReplaceAll(s, "\n", "~")
}

type captureHandler struct {
	mu     sync.Mutex
	r      Record
	attrs  []Attr
	groups []string
}

func (h *captureHandler) Handle(ctx context.Context, r Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.r = r
	return nil
}

func (*captureHandler) Enabled(context.Context, Level) bool { return true }

func (c *captureHandler) WithAttrs(as []Attr) Handler {
	c.mu.Lock()
	defer c.mu.Unlock()
	var c2 captureHandler
	c2.r = c.r
	c2.groups = c.groups
	c2.attrs = concat(c.attrs, as)
	return &c2
}

func (c *captureHandler) WithGroup(name string) Handler {
	c.mu.Lock()
	defer c.mu.Unlock()
	var c2 captureHandler
	c2.r = c.r
	c2.attrs = c.attrs
	c2.groups = append(slices.Clip(c.groups), name)
	return &c2
}

type discardHandler struct {
	disabled bool
	attrs    []Attr
}

func (d discardHandler) Enabled(context.Context, Level) bool { return !d.disabled }
func (discardHandler) Handle(context.Context, Record) error  { return nil }
func (d discardHandler) WithAttrs(as []Attr) Handler {
	d.attrs = concat(d.attrs, as)
	return d
}
func (h discardHandler) WithGroup(name string) Handler {
	return h
}

// concat returns a new slice with the elements of s1 followed
// by those of s2. The slice has no additional capacity.
func concat[T any](s1, s2 []T) []T {
	s := make([]T, len(s1)+len(s2))
	copy(s, s1)
	copy(s[len(s1):], s2)
	return s
}

// This is a simple benchmark. See the benchmarks subdirectory for more extensive ones.
func BenchmarkNopLog(b *testing.B) {
	ctx := context.Background()
	l := New(&captureHandler{})
	b.Run("no attrs", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			l.LogAttrs(nil, LevelInfo, "msg")
		}
	})
	b.Run("attrs", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			l.LogAttrs(nil, LevelInfo, "msg", Int("a", 1), String("b", "two"), Bool("c", true))
		}
	})
	b.Run("attrs-parallel", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.LogAttrs(nil, LevelInfo, "msg", Int("a", 1), String("b", "two"), Bool("c", true))
			}
		})
	})
	b.Run("keys-values", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			l.Log(nil, LevelInfo, "msg", "a", 1, "b", "two", "c", true)
		}
	})
	b.Run("WithContext", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			l.LogAttrs(ctx, LevelInfo, "msg2", Int("a", 1), String("b", "two"), Bool("c", true))
		}
	})
	b.Run("WithContext-parallel", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.LogAttrs(ctx, LevelInfo, "msg", Int("a", 1), String("b", "two"), Bool("c", true))
			}
		})
	})
}

// callerPC returns the program counter at the given stack depth.
func callerPC(depth int) uintptr {
	var pcs [1]uintptr
	runtime.Callers(depth, pcs[:])
	return pcs[0]
}

func wantAllocs(t *testing.T, want int, f func()) {
	if race.Enabled {
		t.Skip("skipping test in race mode")
	}
	testenv.SkipIfOptimizationOff(t)
	t.Helper()
	got := int(testing.AllocsPerRun(5, f))
	if got != want {
		t.Errorf("got %d allocs, want %d", got, want)
	}
}
