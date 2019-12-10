// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logopt

import (
	"cmd/internal/obj"
	"cmd/internal/objabi"
	"cmd/internal/src"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// This implements (non)optimization logging for -json option to the Go compiler
// The option is -json 0,<destination>.
//
// 0 is the version number; to avoid the need for synchronized updates, if
// new versions of the logging appear, the compiler will support both, for a while,
// and clients will specify what they need.
//
// <destination> is a directory.
// Directories are specified with a leading / or os.PathSeparator,
// or more explicitly with file://directory.  The second form is intended to
// deal with corner cases on Windows, and to allow specification of a relative
// directory path (which is normally a bad idea, because the local directory
// varies a lot in a build, especially with modules and/or vendoring, and may
// not be writeable).
//
// For each package pkg compiled, a url.PathEscape(pkg)-named subdirectory
// is created.  For each source file.go in that package that generates
// diagnostics (no diagnostics means no file),
// a url.PathEscape(file)+".json"-named file is created and contains the
// logged diagnostics.
//
// For example, "cmd%2Finternal%2Fdwarf/%3Cautogenerated%3E.json"
// for "cmd/internal/dwarf" and <autogenerated> (which is not really a file, but the compiler sees it)
//
// If the package string is empty, it is replaced internally with string(0) which encodes to %00.
//
// Each log file begins with a JSON record identifying version,
// platform, and other context, followed by optimization-relevant
// LSP Diagnostic records, one per line (LSP version 3.15, no difference from 3.14 on the subset used here
// see https://microsoft.github.io/language-server-protocol/specifications/specification-3-15/ )
//
// The fields of a Diagnostic are used in the following way:
// Range: the outermost source position, for now begin and end are equal.
// Severity: (always) SeverityInformation (3)
// Source: (always) "go compiler"
// Code: a string describing the missed optimization, e.g., "nilcheck", "cannotInline", "isInBounds", "escape"
// Message: depending on code, additional information, e.g., the reason a function cannot be inlined.
// RelatedInformation: if the missed optimization actually occurred at a function inlined at Range,
//    then the sequence of inlined locations appears here, from (second) outermost to innermost,
//    each with message="inlineLoc".
//
//    In the case of escape analysis explanations, after any outer inlining locations,
//    the lines of the explanation appear, each potentially followed with its own inlining
//    location if the escape flow occurred within an inlined function.
//
// For example <destination>/cmd%2Fcompile%2Finternal%2Fssa/prove.json
// might begin with the following line (wrapped for legibility):
//
// {"version":0,"package":"cmd/compile/internal/ssa","goos":"darwin","goarch":"amd64",
//  "gc_version":"devel +e1b9a57852 Fri Nov 1 15:07:00 2019 -0400",
//  "file":"/Users/drchase/work/go/src/cmd/compile/internal/ssa/prove.go"}
//
// and later contain (also wrapped for legibility):
//
// {"range":{"start":{"line":191,"character":24},"end":{"line":191,"character":24}},
//  "severity":3,"code":"nilcheck","source":"go compiler","message":"",
//  "relatedInformation":[
//    {"location":{"uri":"file:///Users/drchase/work/go/src/cmd/compile/internal/ssa/func.go",
//                 "range":{"start":{"line":153,"character":16},"end":{"line":153,"character":16}}},
//     "message":"inlineLoc"}]}
//
// That is, at prove.go (implicit from context, provided in both filename and header line),
// line 191, column 24, a nilcheck occurred in the generated code.
// The relatedInformation indicates that this code actually came from
// an inlined call to func.go, line 153, character 16.
//
// prove.go:191:
// 	ft.orderS = f.newPoset()
// func.go:152 and 153:
//  func (f *Func) newPoset() *poset {
//	    if len(f.Cache.scrPoset) > 0 {
//
// In the case that the package is empty, the string(0) package name is also used in the header record, for example
//
//  go tool compile -json=0,file://logopt x.go       # no -p option to set the package
//  head -1 logopt/%00/x.json
//  {"version":0,"package":"\u0000","goos":"darwin","goarch":"amd64","gc_version":"devel +86487adf6a Thu Nov 7 19:34:56 2019 -0500","file":"x.go"}

type VersionHeader struct {
	Version   int    `json:"version"`
	Package   string `json:"package"`
	Goos      string `json:"goos"`
	Goarch    string `json:"goarch"`
	GcVersion string `json:"gc_version"`
	File      string `json:"file,omitempty"` // LSP requires an enclosing resource, i.e., a file
}

// DocumentURI, Position, Range, Location, Diagnostic, DiagnosticRelatedInformation all reuse json definitions from gopls.
// See https://github.com/golang/tools/blob/22afafe3322a860fcd3d88448768f9db36f8bc5f/internal/lsp/protocol/tsprotocol.go

type DocumentURI string

type Position struct {
	Line      uint `json:"line"`      // gopls uses float64, but json output is the same for integers
	Character uint `json:"character"` // gopls uses float64, but json output is the same for integers
}

// A Range in a text document expressed as (zero-based) start and end positions.
// A range is comparable to a selection in an editor. Therefore the end position is exclusive.
// If you want to specify a range that contains a line including the line ending character(s)
// then use an end position denoting the start of the next line.
type Range struct {
	/*Start defined:
	 * The range's start position
	 */
	Start Position `json:"start"`

	/*End defined:
	 * The range's end position
	 */
	End Position `json:"end"` // exclusive
}

// A Location represents a location inside a resource, such as a line inside a text file.
type Location struct {
	// URI is
	URI DocumentURI `json:"uri"`

	// Range is
	Range Range `json:"range"`
}

/* DiagnosticRelatedInformation defined:
 * Represents a related message and source code location for a diagnostic. This should be
 * used to point to code locations that cause or related to a diagnostics, e.g when duplicating
 * a symbol in a scope.
 */
type DiagnosticRelatedInformation struct {

	/*Location defined:
	 * The location of this related diagnostic information.
	 */
	Location Location `json:"location"`

	/*Message defined:
	 * The message of this related diagnostic information.
	 */
	Message string `json:"message"`
}

// DiagnosticSeverity defines constants
type DiagnosticSeverity uint

const (
	/*SeverityInformation defined:
	 * Reports an information.
	 */
	SeverityInformation DiagnosticSeverity = 3
)

// DiagnosticTag defines constants
type DiagnosticTag uint

/*Diagnostic defined:
 * Represents a diagnostic, such as a compiler error or warning. Diagnostic objects
 * are only valid in the scope of a resource.
 */
type Diagnostic struct {

	/*Range defined:
	 * The range at which the message applies
	 */
	Range Range `json:"range"`

	/*Severity defined:
	 * The diagnostic's severity. Can be omitted. If omitted it is up to the
	 * client to interpret diagnostics as error, warning, info or hint.
	 */
	Severity DiagnosticSeverity `json:"severity,omitempty"` // always SeverityInformation for optimizer logging.

	/*Code defined:
	 * The diagnostic's code, which usually appear in the user interface.
	 */
	Code string `json:"code,omitempty"` // LSP uses 'number | string' = gopls interface{}, but only string here, e.g. "boundsCheck", "nilcheck", etc.

	/*Source defined:
	 * A human-readable string describing the source of this
	 * diagnostic, e.g. 'typescript' or 'super lint'. It usually
	 * appears in the user interface.
	 */
	Source string `json:"source,omitempty"` // "go compiler"

	/*Message defined:
	 * The diagnostic's message. It usually appears in the user interface
	 */
	Message string `json:"message"` // sometimes used, provides additional information.

	/*Tags defined:
	 * Additional metadata about the diagnostic.
	 */
	Tags []DiagnosticTag `json:"tags,omitempty"` // always empty for logging optimizations.

	/*RelatedInformation defined:
	 * An array of related diagnostic information, e.g. when symbol-names within
	 * a scope collide all definitions can be marked via this property.
	 */
	RelatedInformation []DiagnosticRelatedInformation `json:"relatedInformation,omitempty"`
}

// A LoggedOpt is what the compiler produces and accumulates,
// to be converted to JSON for human or IDE consumption.
type LoggedOpt struct {
	pos    src.XPos      // Source code position at which the event occurred. If it is inlined, outer and all inlined locations will appear in JSON.
	pass   string        // For human/adhoc consumption; does not appear in JSON (yet)
	fname  string        // For human/adhoc consumption; does not appear in JSON (yet)
	what   string        // The (non) optimization; "nilcheck", "boundsCheck", "inline", "noInline"
	target []interface{} // Optional target(s) or parameter(s) of "what" -- what was inlined, why it was not, size of copy, etc. 1st is most important/relevant.
}

type logFormat uint8

const (
	None  logFormat = iota
	Json0           // version 0 for LSP 3.14, 3.15; future versions of LSP may change the format and the compiler may need to support both as clients are updated.
)

var Format = None
var dest string

func LogJsonOption(flagValue string) {
	version, directory := parseLogFlag("json", flagValue)
	if version != 0 {
		log.Fatal("-json version must be 0")
	}
	checkLogPath("json", directory)
	Format = Json0
}

// parseLogFlag checks the flag passed to -json
// for version,destination format and returns the two parts.
func parseLogFlag(flag, value string) (version int, directory string) {
	if Format != None {
		log.Fatal("Cannot repeat -json flag")
	}
	commaAt := strings.Index(value, ",")
	if commaAt <= 0 {
		log.Fatalf("-%s option should be '<version>,<destination>' where <version> is a number", flag)
	}
	v, err := strconv.Atoi(value[:commaAt])
	if err != nil {
		log.Fatalf("-%s option should be '<version>,<destination>' where <version> is a number: err=%v", flag, err)
	}
	version = v
	directory = value[commaAt+1:]
	return
}

// checkLogPath does superficial early checking of the string specifying
// the directory to which optimizer logging is directed, and if
// it passes the test, stores the string in LO_dir
func checkLogPath(flag, destination string) {
	sep := string(os.PathSeparator)
	if strings.HasPrefix(destination, "/") || strings.HasPrefix(destination, sep) {
		err := os.MkdirAll(destination, 0755)
		if err != nil {
			log.Fatalf("optimizer logging destination '<version>,<directory>' but could not create <directory>: err=%v", err)
		}
	} else if strings.HasPrefix(destination, "file://") { // IKWIAD, or Windows C:\foo\bar\baz
		uri, err := url.Parse(destination)
		if err != nil {
			log.Fatalf("optimizer logging destination looked like file:// URI but failed to parse: err=%v", err)
		}
		destination = uri.Host + uri.Path
		err = os.MkdirAll(destination, 0755)
		if err != nil {
			log.Fatalf("optimizer logging destination '<version>,<directory>' but could not create %s: err=%v", destination, err)
		}
	} else {
		log.Fatalf("optimizer logging destination %s was neither %s-prefixed directory nor file://-prefixed file URI", destination, sep)
	}
	dest = destination
}

var loggedOpts []LoggedOpt
var mu = sync.Mutex{} // mu protects loggedOpts.

func LogOpt(pos src.XPos, what, pass, fname string, args ...interface{}) {
	if Format == None {
		return
	}
	pass = strings.Replace(pass, " ", "_", -1)
	mu.Lock()
	defer mu.Unlock()
	// Because of concurrent calls from back end, no telling what the order will be, but is stable-sorted by outer Pos before use.
	loggedOpts = append(loggedOpts, LoggedOpt{pos, pass, fname, what, args})
}

func Enabled() bool {
	switch Format {
	case None:
		return false
	case Json0:
		return true
	}
	panic("Unexpected optimizer-logging level")
}

// byPos sorts diagnostics by source position.
type byPos struct {
	ctxt *obj.Link
	a    []LoggedOpt
}

func (x byPos) Len() int { return len(x.a) }
func (x byPos) Less(i, j int) bool {
	return x.ctxt.OutermostPos(x.a[i].pos).Before(x.ctxt.OutermostPos(x.a[j].pos))
}
func (x byPos) Swap(i, j int) { x.a[i], x.a[j] = x.a[j], x.a[i] }

func writerForLSP(subdirpath, file string) io.WriteCloser {
	basename := file
	lastslash := strings.LastIndexAny(basename, "\\/")
	if lastslash != -1 {
		basename = basename[lastslash+1:]
	}
	lastdot := strings.LastIndex(basename, ".go")
	if lastdot != -1 {
		basename = basename[:lastdot]
	}
	basename = pathEscape(basename)

	// Assume a directory, make a file
	p := filepath.Join(subdirpath, basename+".json")
	w, err := os.Create(p)
	if err != nil {
		log.Fatalf("Could not create file %s for logging optimizer actions, %v", p, err)
	}
	return w
}

func fixSlash(f string) string {
	if os.PathSeparator == '/' {
		return f
	}
	return strings.Replace(f, string(os.PathSeparator), "/", -1)
}

func uriIfy(f string) DocumentURI {
	url := url.URL{
		Scheme: "file",
		Path:   fixSlash(f),
	}
	return DocumentURI(url.String())
}

// Return filename, replacing a first occurrence of $GOROOT with the
// actual value of the GOROOT (because LSP does not speak "$GOROOT").
func uprootedPath(filename string) string {
	if !strings.HasPrefix(filename, "$GOROOT/") {
		return filename
	}
	return objabi.GOROOT + filename[len("$GOROOT"):]
}

// FlushLoggedOpts flushes all the accumulated optimization log entries.
func FlushLoggedOpts(ctxt *obj.Link, slashPkgPath string) {
	if Format == None {
		return
	}

	sort.Stable(byPos{ctxt, loggedOpts}) // Stable is necessary to preserve the per-function order, which is repeatable.
	switch Format {

	case Json0: // LSP 3.15
		var posTmp []src.Pos
		var encoder *json.Encoder
		var w io.WriteCloser

		if slashPkgPath == "" {
			slashPkgPath = string(0)
		}
		subdirpath := filepath.Join(dest, pathEscape(slashPkgPath))
		err := os.MkdirAll(subdirpath, 0755)
		if err != nil {
			log.Fatalf("Could not create directory %s for logging optimizer actions, %v", subdirpath, err)
		}
		diagnostic := Diagnostic{Source: "go compiler", Severity: SeverityInformation}

		// For LSP, make a subdirectory for the package, and for each file foo.go, create foo.json in that subdirectory.
		currentFile := ""
		for _, x := range loggedOpts {
			posTmp = ctxt.AllPos(x.pos, posTmp)
			// Reverse posTmp to put outermost first.
			l := len(posTmp)
			for i := 0; i < l/2; i++ {
				posTmp[i], posTmp[l-i-1] = posTmp[l-i-1], posTmp[i]
			}

			p0 := posTmp[0]
			p0f := uprootedPath(p0.Filename())
			if currentFile != p0f {
				if w != nil {
					w.Close()
				}
				currentFile = p0f
				w = writerForLSP(subdirpath, currentFile)
				encoder = json.NewEncoder(w)
				encoder.Encode(VersionHeader{Version: 0, Package: slashPkgPath, Goos: objabi.GOOS, Goarch: objabi.GOARCH, GcVersion: objabi.Version, File: currentFile})
			}

			// The first "target" is the most important one.
			var target string
			if len(x.target) > 0 {
				target = fmt.Sprint(x.target[0])
			}

			diagnostic.Code = x.what
			diagnostic.Message = target
			diagnostic.Range = Range{Start: Position{p0.Line(), p0.Col()},
				End: Position{p0.Line(), p0.Col()}}
			diagnostic.RelatedInformation = diagnostic.RelatedInformation[:0]

			for i := 1; i < l; i++ {
				p := posTmp[i]
				loc := Location{URI: uriIfy(uprootedPath(p.Filename())),
					Range: Range{Start: Position{p.Line(), p.Col()},
						End: Position{p.Line(), p.Col()}}}
				diagnostic.RelatedInformation = append(diagnostic.RelatedInformation, DiagnosticRelatedInformation{Location: loc, Message: "inlineLoc"})
			}

			encoder.Encode(diagnostic)
		}
		if w != nil {
			w.Close()
		}
	}
}
