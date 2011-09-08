// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package regexp

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"utf8"
)

// TestRE2 tests this package's regexp API against test cases
// considered during RE2's exhaustive tests, which run all possible
// regexps over a given set of atoms and operators, up to a given
// complexity, over all possible strings over a given alphabet,
// up to a given size.  Rather than try to link with RE2, we read a
// log file containing the test cases and the expected matches.
// The log file, re2.txt, is generated by running 'make exhaustive-log'
// in the open source RE2 distribution.  http://code.google.com/p/re2/
//
// The test file format is a sequence of stanzas like:
//
//	strings
//	"abc"
//	"123x"
//	regexps
//	"[a-z]+"
//	0-3;0-3
//	-;-
//	"([0-9])([0-9])([0-9])"
//	-;-
//	-;0-3 0-1 1-2 2-3
//
// The stanza begins by defining a set of strings, quoted
// using Go double-quote syntax, one per line.  Then the
// regexps section gives a sequence of regexps to run on
// the strings.  In the block that follows a regexp, each line
// gives the semicolon-separated match results of running
// the regexp on the corresponding string.
// Each match result is either a single -, meaning no match, or a
// space-separated sequence of pairs giving the match and
// submatch indices.  An unmatched subexpression formats
// its pair as a single - (not illustrated above).  For now
// each regexp run produces two match results, one for a
// ``full match'' that restricts the regexp to matching the entire
// string or nothing, and one for a ``partial match'' that gives
// the leftmost first match found in the string.
//
// Lines beginning with # are comments.  Lines beginning with
// a capital letter are test names printed during RE2's test suite
// and are echoed into t but otherwise ignored.
//
// At time of writing, re2.txt is 32 MB but compresses to 760 kB,
// so we store re2.txt.gz in the repository and decompress it on the fly.
//
func TestRE2(t *testing.T) {
	if testing.Short() {
		t.Log("skipping TestRE2 during short test")
		return
	}

	f, err := os.Open("re2.txt.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		t.Fatalf("decompress re2.txt.gz: %v", err)
	}
	defer gz.Close()
	lineno := 0
	r := bufio.NewReader(gz)
	var (
		str       []string
		input     []string
		inStrings bool
		re        *Regexp
		refull    *Regexp
		nfail     int
		ncase     int
	)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == os.EOF {
				break
			}
			t.Fatalf("re2.txt:%d: %v", lineno, err)
		}
		line = line[:len(line)-1] // chop \n
		lineno++
		switch {
		case line == "":
			t.Fatalf("re2.txt:%d: unexpected blank line", lineno)
		case line[0] == '#':
			continue
		case 'A' <= line[0] && line[0] <= 'Z':
			// Test name.
			t.Logf("%s\n", line)
			continue
		case line == "strings":
			str = str[:0]
			inStrings = true
		case line == "regexps":
			inStrings = false
		case line[0] == '"':
			q, err := strconv.Unquote(line)
			if err != nil {
				// Fatal because we'll get out of sync.
				t.Fatalf("re2.txt:%d: unquote %s: %v", lineno, line, err)
			}
			if inStrings {
				str = append(str, q)
				continue
			}
			// Is a regexp.
			if len(input) != 0 {
				t.Fatalf("re2.txt:%d: out of sync: have %d strings left before %#q", lineno, len(input), q)
			}
			re, err = tryCompile(q)
			if err != nil {
				if err.String() == "error parsing regexp: invalid escape sequence: `\\C`" {
					// We don't and likely never will support \C; keep going.
					continue
				}
				t.Errorf("re2.txt:%d: compile %#q: %v", lineno, q, err)
				if nfail++; nfail >= 100 {
					t.Fatalf("stopping after %d errors", nfail)
				}
				continue
			}
			full := `\A(?:` + q + `)\z`
			refull, err = tryCompile(full)
			if err != nil {
				// Fatal because q worked, so this should always work.
				t.Fatalf("re2.txt:%d: compile full %#q: %v", lineno, full, err)
			}
			input = str
		case line[0] == '-' || '0' <= line[0] && line[0] <= '9':
			// A sequence of match results.
			ncase++
			if re == nil {
				// Failed to compile: skip results.
				continue
			}
			if len(input) == 0 {
				t.Fatalf("re2.txt:%d: out of sync: no input remaining", lineno)
			}
			var text string
			text, input = input[0], input[1:]
			if !isSingleBytes(text) && strings.Contains(re.String(), `\B`) {
				// RE2's \B considers every byte position,
				// so it sees 'not word boundary' in the
				// middle of UTF-8 sequences.  This package
				// only considers the positions between runes,
				// so it disagrees.  Skip those cases.
				continue
			}
			res := strings.Split(line, ";")
			if len(res) != len(run) {
				t.Fatalf("re2.txt:%d: have %d test results, want %d", lineno, len(res), len(run))
			}
			for i := range res {
				have, suffix := run[i](re, refull, text)
				want := parseResult(t, lineno, res[i])
				if !same(have, want) {
					t.Errorf("re2.txt:%d: %#q%s.FindSubmatchIndex(%#q) = %v, want %v", lineno, re, suffix, text, have, want)
					if nfail++; nfail >= 100 {
						t.Fatalf("stopping after %d errors", nfail)
					}
					continue
				}
				b, suffix := match[i](re, refull, text)
				if b != (want != nil) {
					t.Errorf("re2.txt:%d: %#q%s.MatchString(%#q) = %v, want %v", lineno, re, suffix, text, b, !b)
					if nfail++; nfail >= 100 {
						t.Fatalf("stopping after %d errors", nfail)
					}
					continue
				}
			}

		default:
			t.Fatalf("re2.txt:%d: out of sync: %s\n", lineno, line)
		}
	}
	if len(input) != 0 {
		t.Fatalf("re2.txt:%d: out of sync: have %d strings left at EOF", lineno, len(input))
	}
	t.Logf("%d cases tested", ncase)
}

var run = []func(*Regexp, *Regexp, string) ([]int, string){
	runFull,
	runPartial,
	runFullLongest,
	runPartialLongest,
}

func runFull(re, refull *Regexp, text string) ([]int, string) {
	refull.longest = false
	return refull.FindStringSubmatchIndex(text), "[full]"
}

func runPartial(re, refull *Regexp, text string) ([]int, string) {
	re.longest = false
	return re.FindStringSubmatchIndex(text), ""
}

func runFullLongest(re, refull *Regexp, text string) ([]int, string) {
	refull.longest = true
	return refull.FindStringSubmatchIndex(text), "[full,longest]"
}

func runPartialLongest(re, refull *Regexp, text string) ([]int, string) {
	re.longest = true
	return re.FindStringSubmatchIndex(text), "[longest]"
}

var match = []func(*Regexp, *Regexp, string) (bool, string){
	matchFull,
	matchPartial,
	matchFullLongest,
	matchPartialLongest,
}

func matchFull(re, refull *Regexp, text string) (bool, string) {
	refull.longest = false
	return refull.MatchString(text), "[full]"
}

func matchPartial(re, refull *Regexp, text string) (bool, string) {
	re.longest = false
	return re.MatchString(text), ""
}

func matchFullLongest(re, refull *Regexp, text string) (bool, string) {
	refull.longest = true
	return refull.MatchString(text), "[full,longest]"
}

func matchPartialLongest(re, refull *Regexp, text string) (bool, string) {
	re.longest = true
	return re.MatchString(text), "[longest]"
}

func isSingleBytes(s string) bool {
	for _, c := range s {
		if c >= utf8.RuneSelf {
			return false
		}
	}
	return true
}

func tryCompile(s string) (re *Regexp, err os.Error) {
	// Protect against panic during Compile.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	return Compile(s)
}

func parseResult(t *testing.T, lineno int, res string) []int {
	// A single - indicates no match.
	if res == "-" {
		return nil
	}
	// Otherwise, a space-separated list of pairs.
	n := 1
	for j := 0; j < len(res); j++ {
		if res[j] == ' ' {
			n++
		}
	}
	out := make([]int, 2*n)
	i := 0
	n = 0
	for j := 0; j <= len(res); j++ {
		if j == len(res) || res[j] == ' ' {
			// Process a single pair.  - means no submatch.
			pair := res[i:j]
			if pair == "-" {
				out[n] = -1
				out[n+1] = -1
			} else {
				k := strings.Index(pair, "-")
				if k < 0 {
					t.Fatalf("re2.txt:%d: invalid pair %s", lineno, pair)
				}
				lo, err1 := strconv.Atoi(pair[:k])
				hi, err2 := strconv.Atoi(pair[k+1:])
				if err1 != nil || err2 != nil || lo > hi {
					t.Fatalf("re2.txt:%d: invalid pair %s", lineno, pair)
				}
				out[n] = lo
				out[n+1] = hi
			}
			n += 2
			i = j + 1
		}
	}
	return out
}

func same(x, y []int) bool {
	if len(x) != len(y) {
		return false
	}
	for i, xi := range x {
		if xi != y[i] {
			return false
		}
	}
	return true
}
