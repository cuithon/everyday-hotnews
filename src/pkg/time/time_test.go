// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
	"os"
	"strings"
	"testing"
	"testing/quick"
	. "time"
)

func init() {
	// Force US Pacific time for daylight-savings
	// tests below (localtests).  Needs to be set
	// before the first call into the time library.
	os.Setenv("TZ", "America/Los_Angeles")
}

type TimeTest struct {
	seconds int64
	golden  Time
}

var utctests = []TimeTest{
	TimeTest{0, Time{1970, 1, 1, 0, 0, 0, Thursday, 0, "UTC"}},
	TimeTest{1221681866, Time{2008, 9, 17, 20, 4, 26, Wednesday, 0, "UTC"}},
	TimeTest{-1221681866, Time{1931, 4, 16, 3, 55, 34, Thursday, 0, "UTC"}},
	TimeTest{-11644473600, Time{1601, 1, 1, 0, 0, 0, Monday, 0, "UTC"}},
	TimeTest{599529660, Time{1988, 12, 31, 0, 1, 0, Saturday, 0, "UTC"}},
	TimeTest{978220860, Time{2000, 12, 31, 0, 1, 0, Sunday, 0, "UTC"}},
	TimeTest{1e18, Time{31688740476, 10, 23, 1, 46, 40, Friday, 0, "UTC"}},
	TimeTest{-1e18, Time{-31688736537, 3, 10, 22, 13, 20, Tuesday, 0, "UTC"}},
	TimeTest{0x7fffffffffffffff, Time{292277026596, 12, 4, 15, 30, 7, Sunday, 0, "UTC"}},
	TimeTest{-0x8000000000000000, Time{-292277022657, 1, 27, 8, 29, 52, Sunday, 0, "UTC"}},
}

var localtests = []TimeTest{
	TimeTest{0, Time{1969, 12, 31, 16, 0, 0, Wednesday, -8 * 60 * 60, "PST"}},
	TimeTest{1221681866, Time{2008, 9, 17, 13, 4, 26, Wednesday, -7 * 60 * 60, "PDT"}},
}

func same(t, u *Time) bool {
	return t.Year == u.Year &&
		t.Month == u.Month &&
		t.Day == u.Day &&
		t.Hour == u.Hour &&
		t.Minute == u.Minute &&
		t.Second == u.Second &&
		t.Weekday == u.Weekday &&
		t.ZoneOffset == u.ZoneOffset &&
		t.Zone == u.Zone
}

func TestSecondsToUTC(t *testing.T) {
	for i := 0; i < len(utctests); i++ {
		sec := utctests[i].seconds
		golden := &utctests[i].golden
		tm := SecondsToUTC(sec)
		newsec := tm.Seconds()
		if newsec != sec {
			t.Errorf("SecondsToUTC(%d).Seconds() = %d", sec, newsec)
		}
		if !same(tm, golden) {
			t.Errorf("SecondsToUTC(%d):", sec)
			t.Errorf("  want=%+v", *golden)
			t.Errorf("  have=%+v", *tm)
		}
	}
}

func TestSecondsToLocalTime(t *testing.T) {
	for i := 0; i < len(localtests); i++ {
		sec := localtests[i].seconds
		golden := &localtests[i].golden
		tm := SecondsToLocalTime(sec)
		newsec := tm.Seconds()
		if newsec != sec {
			t.Errorf("SecondsToLocalTime(%d).Seconds() = %d", sec, newsec)
		}
		if !same(tm, golden) {
			t.Errorf("SecondsToLocalTime(%d):", sec)
			t.Errorf("  want=%+v", *golden)
			t.Errorf("  have=%+v", *tm)
		}
	}
}

func TestSecondsToUTCAndBack(t *testing.T) {
	f := func(sec int64) bool { return SecondsToUTC(sec).Seconds() == sec }
	f32 := func(sec int32) bool { return f(int64(sec)) }
	cfg := &quick.Config{MaxCount: 10000}

	// Try a reasonable date first, then the huge ones.
	if err := quick.Check(f32, cfg); err != nil {
		t.Fatal(err)
	}
	if err := quick.Check(f, cfg); err != nil {
		t.Fatal(err)
	}
}

type TimeFormatTest struct {
	time           Time
	formattedValue string
}

var iso8601Formats = []TimeFormatTest{
	TimeFormatTest{Time{2008, 9, 17, 20, 4, 26, Wednesday, 0, "UTC"}, "2008-09-17T20:04:26Z"},
	TimeFormatTest{Time{1994, 9, 17, 20, 4, 26, Wednesday, -18000, "EST"}, "1994-09-17T20:04:26-0500"},
	TimeFormatTest{Time{2000, 12, 26, 1, 15, 6, Wednesday, 15600, "OTO"}, "2000-12-26T01:15:06+0420"},
}

func TestISO8601Conversion(t *testing.T) {
	for _, f := range iso8601Formats {
		if f.time.Format(ISO8601) != f.formattedValue {
			t.Error("ISO8601:")
			t.Errorf("  want=%+v", f.formattedValue)
			t.Errorf("  have=%+v", f.time.Format(ISO8601))
		}
	}
}

type FormatTest struct {
	name   string
	format string
	result string
}

var formatTests = []FormatTest{
	FormatTest{"ANSIC", ANSIC, "Thu Feb  4 21:00:57 2010"},
	FormatTest{"UnixDate", UnixDate, "Thu Feb  4 21:00:57 PST 2010"},
	FormatTest{"RFC850", RFC850, "Thursday, 04-Feb-10 21:00:57 PST"},
	FormatTest{"RFC1123", RFC1123, "Thu, 04 Feb 2010 21:00:57 PST"},
	FormatTest{"ISO8601", ISO8601, "2010-02-04T21:00:57-0800"},
	FormatTest{"Kitchen", Kitchen, "9:00PM"},
	FormatTest{"am/pm", "3pm", "9pm"},
	FormatTest{"AM/PM", "3PM", "9PM"},
}

func TestFormat(t *testing.T) {
	// The numeric time represents Thu Feb  4 21:00:57 PST 2010
	time := SecondsToLocalTime(1265346057)
	for _, test := range formatTests {
		result := time.Format(test.format)
		if result != test.result {
			t.Errorf("%s expected %q got %q", test.name, test.result, result)
		}
	}
}

type ParseTest struct {
	name   string
	format string
	value  string
	hasTZ  bool // contains a time zone
	hasWD  bool // contains a weekday
}

var parseTests = []ParseTest{
	ParseTest{"ANSIC", ANSIC, "Thu Feb  4 21:00:57 2010", false, true},
	ParseTest{"UnixDate", UnixDate, "Thu Feb  4 21:00:57 PST 2010", true, true},
	ParseTest{"RFC850", RFC850, "Thursday, 04-Feb-10 21:00:57 PST", true, true},
	ParseTest{"RFC1123", RFC1123, "Thu, 04 Feb 2010 21:00:57 PST", true, true},
	ParseTest{"ISO8601", ISO8601, "2010-02-04T21:00:57-0800", true, false},
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		time, err := Parse(test.format, test.value)
		if err != nil {
			t.Errorf("%s error: %v", test.name, err)
		} else {
			checkTime(time, &test, t)
		}
	}
}

func checkTime(time *Time, test *ParseTest, t *testing.T) {
	// The time should be Thu Feb  4 21:00:57 PST 2010
	if time.Year != 2010 {
		t.Errorf("%s: bad year: %d not %d\n", test.name, time.Year, 2010)
	}
	if time.Month != 2 {
		t.Errorf("%s: bad month: %d not %d\n", test.name, time.Month, 2)
	}
	if time.Day != 4 {
		t.Errorf("%s: bad day: %d not %d\n", test.name, time.Day, 4)
	}
	if time.Hour != 21 {
		t.Errorf("%s: bad hour: %d not %d\n", test.name, time.Hour, 21)
	}
	if time.Minute != 0 {
		t.Errorf("%s: bad minute: %d not %d\n", test.name, time.Minute, 0)
	}
	if time.Second != 57 {
		t.Errorf("%s: bad second: %d not %d\n", test.name, time.Second, 57)
	}
	if test.hasTZ && time.ZoneOffset != -28800 {
		t.Errorf("%s: bad tz offset: %d not %d\n", test.name, time.ZoneOffset, -28800)
	}
	if test.hasWD && time.Weekday != 4 {
		t.Errorf("%s: bad weekday: %d not %d\n", test.name, time.Weekday, 4)
	}
}

func TestFormatAndParse(t *testing.T) {
	const fmt = "Mon MST " + ISO8601 // all fields
	f := func(sec int64) bool {
		t1 := SecondsToLocalTime(sec)
		t2, err := Parse(fmt, t1.Format(fmt))
		if err != nil {
			t.Errorf("error: %s", err)
			return false
		}
		if !same(t1, t2) {
			t.Errorf("different: %q %q", t1, t2)
			return false
		}
		return true
	}
	f32 := func(sec int32) bool { return f(int64(sec)) }
	cfg := &quick.Config{MaxCount: 10000}

	// Try a reasonable date first, then the huge ones.
	if err := quick.Check(f32, cfg); err != nil {
		t.Fatal(err)
	}
	if err := quick.Check(f, cfg); err != nil {
		t.Fatal(err)
	}
}

type ParseErrorTest struct {
	format string
	value  string
	expect string // must appear within the error
}

var parseErrorTests = []ParseErrorTest{
	ParseErrorTest{ANSIC, "Feb  4 21:00:60 2010", "parse"}, // cannot parse Feb as Mon
	ParseErrorTest{ANSIC, "Thu Feb  4 21:00:57 @2010", "format"},
	ParseErrorTest{ANSIC, "Thu Feb  4 21:00:60 2010", "second out of range"},
	ParseErrorTest{ANSIC, "Thu Feb  4 21:61:57 2010", "minute out of range"},
	ParseErrorTest{ANSIC, "Thu Feb  4 24:00:60 2010", "hour out of range"},
}

func TestParseErrors(t *testing.T) {
	for _, test := range parseErrorTests {
		_, err := Parse(test.format, test.value)
		if err == nil {
			t.Errorf("expected error for %q %q\n", test.format, test.value)
		} else if strings.Index(err.String(), test.expect) < 0 {
			t.Errorf("expected error with %q for %q %q; got %s\n", test.expect, test.format, test.value, err)
		}
	}
}

func BenchmarkSeconds(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Seconds()
	}
}

func BenchmarkNanoseconds(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Nanoseconds()
	}
}

func BenchmarkFormat(b *testing.B) {
	time := SecondsToLocalTime(1265346057)
	for i := 0; i < b.N; i++ {
		time.Format("Mon Jan  2 15:04:05 2006")
	}
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Parse(ANSIC, "Mon Jan  2 15:04:05 2006")
	}
}
