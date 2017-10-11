// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssa_test

import (
	"bytes"
	"flag"
	"fmt"
	"internal/testenv"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

var update = flag.Bool("u", false, "update debug_test reference files")
var verbose = flag.Bool("v", false, "print more information about what's happening")
var dryrun = flag.Bool("n", false, "just print the command line and first bits")
var delve = flag.Bool("d", false, "use delve instead of gdb")
var force = flag.Bool("f", false, "force run under not linux-amd64; also do not use tempdir")

var hexRe = regexp.MustCompile("0x[a-zA-Z0-9]+")
var numRe = regexp.MustCompile("-?[0-9]+")
var stringRe = regexp.MustCompile("\"([^\\\"]|(\\.))*\"")

var gdb = "gdb" // Might be "ggdb" on Darwin, because gdb no longer part of XCode

// TestNexting go-builds a file, then uses a debugger (default gdb, optionally delve)
// to next through the generated executable, recording each line landed at, and
// then compares those lines with reference file(s).
// Flag -u updates the reference file(s).
// Flag -d changes the debugger to delve (and uses delve-specific reference files)
// Flag -v is ever-so-slightly verbose.
// Flag -n is for dry-run, and prints the shell and first debug commands.
//
// The file being tested may contain comments of the form
// //DBG-TAG=(v1,v2,v3)
// where DBG = {gdb,dlv} and TAG={dbg,opt}
// each variable may optionally be followed by a / and one or more of S,A,N
// to indicate normalization of Strings, (hex) addresses, and numbers.
// For example:
/*
	if len(os.Args) > 1 { //gdb-dbg=(hist/A,cannedInput/A) //dlv-dbg=(hist/A,cannedInput/A)
*/
// TODO: not implemented for Delve yet, but this is the plan
//
// After a compiler change that causes a difference in the debug behavior, check
// to see if it is sensible or not, and if it is, update the reference files with
// go test debug_test.go -args -u
// (for Delve)
// go test debug_test.go -args -u -d

func TestNexting(t *testing.T) {
	// Skip this test in an ordinary run.bash.  Too many things
	// can cause it to break.
	if testing.Short() {
		t.Skip("skipping in short mode; see issue #22206")
	}

	testenv.MustHaveGoBuild(t)

	if !*delve && !*force && !(runtime.GOOS == "linux" && runtime.GOARCH == "amd64") {
		// Running gdb on OSX/darwin is very flaky.
		// Sometimes it is called ggdb, depending on how it is installed.
		// It also probably requires an admin password typed into a dialog box.
		// Various architectures tend to differ slightly sometimes, and keeping them
		// all in sync is a pain for people who don't have them all at hand,
		// so limit testing to amd64 (for now)

		t.Skip("Skipped unless -d (delve), -f (force), or linux-amd64")
	}

	if *delve {
		_, err := exec.LookPath("dlv")
		if err != nil {
			t.Fatal("dlv specified on command line with -d but no dlv on path")
		}
	} else {
		_, err := exec.LookPath(gdb)
		if err != nil {
			if runtime.GOOS != "darwin" {
				t.Skip("Skipped because gdb not available")
			}
			_, err = exec.LookPath("ggdb")
			if err != nil {
				t.Skip("Skipped because gdb (and also ggdb) not available")
			}
			gdb = "ggdb"
		}
	}

	t.Run("dbg", func(t *testing.T) {
		testNexting(t, "hist", "dbg", "-N -l")
	})
	t.Run("opt", func(t *testing.T) {
		// If this is test is run with a runtime compiled with -N -l, it is very likely to fail.
		// This occurs in the noopt builders (for example).
		if gogcflags := os.Getenv("GO_GCFLAGS"); *force || (!strings.Contains(gogcflags, "-N") && !strings.Contains(gogcflags, "-l")) {
			testNexting(t, "hist", "opt", "")
		} else {
			t.Skip("skipping for unoptimized runtime")
		}
	})
}

func testNexting(t *testing.T, base, tag, gcflags string) {
	// (1) In testdata, build sample.go into sample
	// (2) Run debugger gathering a history
	// (3) Read expected history from testdata/sample.nexts
	// optionally, write out testdata/sample.nexts

	exe := filepath.Join("testdata", base)
	logbase := exe + "-" + tag
	tmpbase := logbase + "-test"

	if !*force {
		tmpdir, err := ioutil.TempDir("", "debug_test")
		if err != nil {
			panic(fmt.Sprintf("Problem creating TempDir, error %v\n", err))
		}
		exe = filepath.Join(tmpdir, base)
		tmpbase = exe + "-" + tag + "-test"
		if *verbose {
			fmt.Printf("Tempdir is %s\n", tmpdir)
		}
		defer os.RemoveAll(tmpdir)
	}

	if gcflags == "" {
		runGo(t, "", "build", "-o", exe, filepath.Join("testdata", base+".go"))
	} else {
		runGo(t, "", "build", "-o", exe, "-gcflags", gcflags, filepath.Join("testdata", base+".go"))
	}
	var h1 *nextHist
	var nextlog, tmplog string
	if *delve {
		h1 = dlvTest(tag, exe, 1000)
		nextlog = logbase + ".delve-nexts"
		tmplog = tmpbase + ".delve-nexts"
	} else {
		h1 = gdbTest(tag, exe, 1000)
		nextlog = logbase + ".gdb-nexts"
		tmplog = tmpbase + ".gdb-nexts"
	}
	if *dryrun {
		fmt.Printf("# Tag for above is %s\n", tag)
		return
	}
	if *update {
		h1.write(nextlog)
	} else {
		h0 := &nextHist{}
		h0.read(nextlog)
		if !h0.equals(h1) {
			// Be very noisy about exactly what's wrong to simplify debugging.
			h1.write(tmplog)
			cmd := exec.Command("diff", "-u", nextlog, tmplog)
			line := asCommandLine("", cmd)
			bytes, err := cmd.CombinedOutput()
			if err != nil && len(bytes) == 0 {
				t.Fatalf("step/next histories differ, diff command %s failed with error=%v", line, err)
			}
			t.Fatalf("step/next histories differ, diff=\n%s", string(bytes))
		}
	}
}

type dbgr interface {
	start()
	do(s string)
	stepnext(s string) bool // step or next, possible with parameter, gets line etc.  returns true for success, false for unsure response
	quit()
	hist() *nextHist
}

func gdbTest(tag, executable string, maxNext int, args ...string) *nextHist {
	dbg := newGdb(tag, executable, args...)
	dbg.start()
	if *dryrun {
		return nil
	}
	for i := 0; i < maxNext; i++ {
		if !dbg.stepnext("n") {
			break
		}
	}
	h := dbg.hist()
	return h
}

func dlvTest(tag, executable string, maxNext int, args ...string) *nextHist {
	dbg := newDelve(tag, executable, args...)
	dbg.start()
	if *dryrun {
		return nil
	}
	for i := 0; i < maxNext; i++ {
		if !dbg.stepnext("n") {
			break
		}
	}
	h := dbg.hist()
	return h
}

func runGo(t *testing.T, dir string, args ...string) string {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(testenv.GoToolPath(t), args...)
	cmd.Dir = dir
	if *dryrun {
		fmt.Printf("%s\n", asCommandLine("", cmd))
		return ""
	}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("error running cmd (%s): %v\nstdout:\n%sstderr:\n%s\n", asCommandLine("", cmd), err, stdout.String(), stderr.String())
	}

	if s := stderr.String(); s != "" {
		t.Fatalf("Stderr = %s\nWant empty", s)
	}

	return stdout.String()
}

type tstring struct {
	o string
	e string
}

func (t tstring) String() string {
	return t.o + t.e
}

type pos struct {
	line uint16
	file uint8
}

type nextHist struct {
	f2i   map[string]uint8
	fs    []string
	ps    []pos // TODO: plan to automatically do the minimum distance conversion between a reference and a run for nicer errors.
	texts []string
	vars  [][]string
}

func (h *nextHist) write(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(fmt.Sprintf("Problem opening %s, error %v\n", filename, err))
	}
	defer file.Close()
	var lastfile uint8
	for i, x := range h.texts {
		p := h.ps[i]
		if lastfile != p.file {
			fmt.Fprintf(file, "  %s\n", h.fs[p.file-1])
			lastfile = p.file
		}
		fmt.Fprintf(file, "%d:%s\n", p.line, x)
		// Vars must begin with a dollar-sign.
		// TODO, normalize between gdb and dlv into a common, comparable format.
		for _, y := range h.vars[i] {
			y = strings.TrimSpace(y)
			if y[0] != '$' {
				panic(fmt.Sprintf("Var line '%s' must begin with $, but does not\n", y))
			}
			fmt.Fprintf(file, "%s\n", y)
		}
	}
	file.Close()
}

func (h *nextHist) read(filename string) {
	h.f2i = make(map[string]uint8)
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("Problem reading %s, error %v\n", filename, err))
	}
	var lastfile string
	lines := strings.Split(string(bytes), "\n")
	for i, l := range lines {
		if len(l) > 0 && l[0] != '#' {
			if l[0] == ' ' {
				// file -- first two characters expected to be "  "
				lastfile = strings.TrimSpace(l)
			} else if l[0] == '$' {
				h.addVar(l)
			} else {
				// line number -- <number>:<line>
				colonPos := strings.Index(l, ":")
				if colonPos == -1 {
					panic(fmt.Sprintf("Line %d (%s) in file %s expected to contain '<number>:' but does not.\n", i+1, l, filename))
				}
				h.add(lastfile, l[0:colonPos], l[colonPos+1:])
			}
		}
	}
}

func (h *nextHist) add(file, line, text string) {
	fi := h.f2i[file]
	if fi == 0 {
		h.fs = append(h.fs, file)
		fi = uint8(len(h.fs))
		h.f2i[file] = fi
	}

	line = strings.TrimSpace(line)
	var li int
	var err error
	if line != "" {
		li, err = strconv.Atoi(line)
		if err != nil {
			panic(fmt.Sprintf("Non-numeric line: %s, error %v\n", line, err))
		}
	}
	h.ps = append(h.ps, pos{line: uint16(li), file: fi})
	h.texts = append(h.texts, text)
	h.vars = append(h.vars, []string{})
}

func (h *nextHist) addVar(text string) {
	l := len(h.texts)
	h.vars[l-1] = append(h.vars[l-1], text)
}

func invertMapSU8(hf2i map[string]uint8) map[uint8]string {
	hi2f := make(map[uint8]string)
	for hs, i := range hf2i {
		hi2f[i] = hs
	}
	return hi2f
}

func (h *nextHist) equals(k *nextHist) bool {
	if len(h.f2i) != len(k.f2i) {
		return false
	}
	if len(h.ps) != len(k.ps) {
		return false
	}
	hi2f := invertMapSU8(h.f2i)
	ki2f := invertMapSU8(k.f2i)

	for i, hs := range hi2f {
		if hs != ki2f[i] {
			return false
		}
	}

	for i, x := range h.ps {
		if k.ps[i] != x {
			return false
		}
	}
	return true
}

func canonFileName(f string) string {
	i := strings.Index(f, "/src/")
	if i != -1 {
		f = f[i+1:]
	}
	return f
}

/* Delve */

type delveState struct {
	cmd *exec.Cmd
	tag string
	*ioState
	atLineRe         *regexp.Regexp // "\n =>"
	funcFileLinePCre *regexp.Regexp // "^> ([^ ]+) ([^:]+):([0-9]+) .*[(]PC: (0x[a-z0-9]+)"
	line             string
	file             string
	function         string
}

func newDelve(tag, executable string, args ...string) dbgr {
	cmd := exec.Command("dlv", "exec", executable)
	cmd.Env = replaceEnv(cmd.Env, "TERM", "dumb")
	if len(args) > 0 {
		cmd.Args = append(cmd.Args, "--")
		cmd.Args = append(cmd.Args, args...)
	}
	s := &delveState{tag: tag, cmd: cmd}
	// HAHA Delve has control characters embedded to change the color of the => and the line number
	// that would be '(\\x1b\\[[0-9;]+m)?' OR TERM=dumb
	s.atLineRe = regexp.MustCompile("\n=>[[:space:]]+[0-9]+:(.*)")
	s.funcFileLinePCre = regexp.MustCompile("> ([^ ]+) ([^:]+):([0-9]+) .*[(]PC: (0x[a-z0-9]+)[)]\n")
	s.ioState = newIoState(s.cmd)
	return s
}

func (s *delveState) stepnext(ss string) bool {
	x := s.ioState.writeReadExpect(ss+"\n", "[(]dlv[)] ")
	excerpts := s.atLineRe.FindStringSubmatch(x.o)
	locations := s.funcFileLinePCre.FindStringSubmatch(x.o)
	excerpt := ""
	if len(excerpts) > 1 {
		excerpt = excerpts[1]
	}
	if len(locations) > 0 {
		fn := canonFileName(locations[2])
		if *verbose {
			if s.file != fn {
				fmt.Printf("%s\n", locations[2]) // don't canonocalize verbose logging
			}
			fmt.Printf("  %s\n", locations[3])
		}
		s.line = locations[3]
		s.file = fn
		s.function = locations[1]
		s.ioState.history.add(s.file, s.line, excerpt)
		return true
	}
	fmt.Printf("DID NOT MATCH EXPECTED NEXT OUTPUT\nO='%s'\nE='%s'\n", x.o, x.e)
	return false
}

func (s *delveState) start() {
	if *dryrun {
		fmt.Printf("%s\n", asCommandLine("", s.cmd))
		fmt.Printf("b main.main\n")
		fmt.Printf("c\n")
		return
	}
	err := s.cmd.Start()
	if err != nil {
		line := asCommandLine("", s.cmd)
		panic(fmt.Sprintf("There was an error [start] running '%s', %v\n", line, err))
	}
	s.ioState.readExpecting(-1, 5000, "Type 'help' for list of commands.")
	expect("Breakpoint [0-9]+ set at ", s.ioState.writeRead("b main.main\n"))
	s.stepnext("c")
}

func (s *delveState) quit() {
	s.do("q")
}

func (s *delveState) do(ss string) {
	expect("", s.ioState.writeRead(ss+"\n"))
}

/* Gdb */

type gdbState struct {
	cmd  *exec.Cmd
	tag  string
	args []string
	*ioState
	atLineRe         *regexp.Regexp
	funcFileLinePCre *regexp.Regexp
	line             string
	file             string
	function         string
}

func newGdb(tag, executable string, args ...string) dbgr {
	// Turn off shell, necessary for Darwin apparently
	cmd := exec.Command(gdb, "-ex", "set startup-with-shell off", executable)
	cmd.Env = replaceEnv(cmd.Env, "TERM", "dumb")
	s := &gdbState{tag: tag, cmd: cmd, args: args}
	s.atLineRe = regexp.MustCompile("(^|\n)([0-9]+)(.*)")
	s.funcFileLinePCre = regexp.MustCompile(
		"([^ ]+) [(][)][ \\t\\n]+at ([^:]+):([0-9]+)")
	// runtime.main () at /Users/drchase/GoogleDrive/work/go/src/runtime/proc.go:201
	//                                    function              file    line
	// Thread 2 hit Breakpoint 1, main.main () at /Users/drchase/GoogleDrive/work/debug/hist.go:18
	s.ioState = newIoState(s.cmd)
	return s
}

func (s *gdbState) start() {
	run := "run"
	for _, a := range s.args {
		run += " " + a // Can't quote args for gdb, it will pass them through including the quotes
	}
	if *dryrun {
		fmt.Printf("%s\n", asCommandLine("", s.cmd))
		fmt.Printf("b main.main\n")
		fmt.Printf("%s\n", run)
		return
	}
	err := s.cmd.Start()
	if err != nil {
		line := asCommandLine("", s.cmd)
		panic(fmt.Sprintf("There was an error [start] running '%s', %v\n", line, err))
	}
	s.ioState.readExpecting(-1, 5000, "[(]gdb[)] ")
	x := s.ioState.writeReadExpect("b main.main\n", "[(]gdb[)] ")
	expect("Breakpoint [0-9]+ at", x)
	s.stepnext(run)
}

func (s *gdbState) stepnext(ss string) bool {
	x := s.ioState.writeReadExpect(ss+"\n", "[(]gdb[)] ")
	excerpts := s.atLineRe.FindStringSubmatch(x.o)
	locations := s.funcFileLinePCre.FindStringSubmatch(x.o)
	excerpt := ""
	if len(excerpts) == 0 && len(locations) == 0 {
		fmt.Printf("DID NOT MATCH %s", x.o)
		return false
	}
	if len(excerpts) > 0 {
		excerpt = excerpts[3]
	}
	if len(locations) > 0 {
		fn := canonFileName(locations[2])
		if *verbose {
			if s.file != fn {
				fmt.Printf("%s\n", locations[2])
			}
			fmt.Printf("  %s\n", locations[3])
		}
		s.line = locations[3]
		s.file = fn
		s.function = locations[1]
		s.ioState.history.add(s.file, s.line, excerpt)
	}
	if len(excerpts) > 0 {
		if *verbose {
			fmt.Printf("  %s\n", excerpts[2])
		}
		s.line = excerpts[2]
		s.ioState.history.add(s.file, s.line, excerpt)
	}

	// Look for //gdb-<tag>=(v1,v2,v3) and print v1, v2, v3
	vars := varsToPrint(excerpt, "//gdb-"+s.tag+"=(")
	for _, v := range vars {
		slashIndex := strings.Index(v, "/")
		substitutions := ""
		if slashIndex != -1 {
			substitutions = v[slashIndex:]
			v = v[:slashIndex]
		}
		response := s.ioState.writeRead("p " + v + "\n").String()
		// expect something like "$1 = ..."
		dollar := strings.Index(response, "$")
		cr := strings.Index(response, "\n")
		if dollar == -1 {
			if cr == -1 {
				response = strings.TrimSpace(response) // discards trailing newline
				response = strings.Replace(response, "\n", "<BR>", -1)
				s.ioState.history.addVar("$ Malformed response " + response)
				continue
			}
			response = strings.TrimSpace(response[:cr])
			s.ioState.history.addVar("$ " + response)
			continue
		}
		if cr == -1 {
			cr = len(response)
		}
		response = strings.TrimSpace(response[dollar:cr])
		if strings.Contains(substitutions, "A") {
			response = hexRe.ReplaceAllString(response, "<A>")
		}
		if strings.Contains(substitutions, "N") {
			response = numRe.ReplaceAllString(response, "<N>")
		}
		if strings.Contains(substitutions, "S") {
			response = stringRe.ReplaceAllString(response, "<S>")
		}
		s.ioState.history.addVar(response)
	}
	return true
}

func varsToPrint(line, lookfor string) []string {
	var vars []string
	if strings.Contains(line, lookfor) {
		x := line[strings.Index(line, lookfor)+len(lookfor):]
		end := strings.Index(x, ")")
		if end == -1 {
			panic(fmt.Sprintf("Saw variable list begin %s in %s but no closing ')'", lookfor, line))
		}
		vars = strings.Split(x[:end], ",")
		for i, y := range vars {
			vars[i] = strings.TrimSpace(y)
		}
	}
	return vars
}

func (s *gdbState) quit() {
	response := s.ioState.writeRead("q\n")
	if strings.Contains(response.o, "Quit anyway? (y or n)") {
		s.ioState.writeRead("Y\n")
	}
}

func (s *gdbState) do(ss string) {
	expect("", s.ioState.writeRead(ss+"\n"))
}

type ioState struct {
	stdout  io.ReadCloser
	stderr  io.ReadCloser
	stdin   io.WriteCloser
	outChan chan string
	errChan chan string
	last    tstring // Output of previous step
	history *nextHist
}

func newIoState(cmd *exec.Cmd) *ioState {
	var err error
	s := &ioState{}
	s.history = &nextHist{}
	s.history.f2i = make(map[string]uint8)
	s.stdout, err = cmd.StdoutPipe()
	line := asCommandLine("", cmd)
	if err != nil {
		panic(fmt.Sprintf("There was an error [stdoutpipe] running '%s', %v\n", line, err))
	}
	s.stderr, err = cmd.StderrPipe()
	if err != nil {
		panic(fmt.Sprintf("There was an error [stdouterr] running '%s', %v\n", line, err))
	}
	s.stdin, err = cmd.StdinPipe()
	if err != nil {
		panic(fmt.Sprintf("There was an error [stdinpipe] running '%s', %v\n", line, err))
	}

	s.outChan = make(chan string, 1)
	s.errChan = make(chan string, 1)
	go func() {
		buffer := make([]byte, 4096)
		for {
			n, err := s.stdout.Read(buffer)
			if n > 0 {
				s.outChan <- string(buffer[0:n])
			}
			if err == io.EOF || n == 0 {
				break
			}
			if err != nil {
				fmt.Printf("Saw an error forwarding stdout")
				break
			}
		}
		close(s.outChan)
		s.stdout.Close()
	}()

	go func() {
		buffer := make([]byte, 4096)
		for {
			n, err := s.stderr.Read(buffer)
			if n > 0 {
				s.errChan <- string(buffer[0:n])
			}
			if err == io.EOF || n == 0 {
				break
			}
			if err != nil {
				fmt.Printf("Saw an error forwarding stderr")
				break
			}
		}
		close(s.errChan)
		s.stderr.Close()
	}()
	return s
}

func (s *ioState) hist() *nextHist {
	return s.history
}

const (
	interlineDelay = 300
)

func (s *ioState) writeRead(ss string) tstring {
	if *verbose {
		fmt.Printf("=> %s", ss)
	}
	_, err := io.WriteString(s.stdin, ss)
	if err != nil {
		panic(fmt.Sprintf("There was an error writing '%s', %v\n", ss, err))
	}
	return s.readWithDelay(-1, interlineDelay)
}

func (s *ioState) writeReadExpect(ss, expect string) tstring {
	if *verbose {
		fmt.Printf("=> %s", ss)
	}
	_, err := io.WriteString(s.stdin, ss)
	if err != nil {
		panic(fmt.Sprintf("There was an error writing '%s', %v\n", ss, err))
	}
	return s.readExpecting(-1, interlineDelay, expect)
}

func (s *ioState) readWithDelay(millis, interlineTimeout int) tstring {
	return s.readExpecting(millis, interlineTimeout, "")
}

func (s *ioState) readExpecting(millis, interlineTimeout int, expected string) tstring {
	timeout := time.Millisecond * time.Duration(millis)
	interline := time.Millisecond * time.Duration(interlineTimeout)
	s.last = tstring{}
	var re *regexp.Regexp
	if expected != "" {
		re = regexp.MustCompile(expected)
	}
loop:
	for {
		var timer <-chan time.Time
		if timeout > 0 {
			timer = time.After(timeout)
		}
		select {
		case x, ok := <-s.outChan:
			if !ok {
				s.outChan = nil
			}
			s.last.o += x
		case x, ok := <-s.errChan:
			if !ok {
				s.errChan = nil
			}
			s.last.e += x
		case <-timer:
			break loop
		}
		if re != nil {
			if re.MatchString(s.last.o) {
				break
			}
			if re.MatchString(s.last.e) {
				break
			}
		}
		timeout = interline
	}
	if *verbose {
		fmt.Printf("<= %s%s", s.last.o, s.last.e)
	}
	return s.last
}

// replaceEnv returns a new environment derived from env
// by removing any existing definition of ev and adding ev=evv.
func replaceEnv(env []string, ev string, evv string) []string {
	evplus := ev + "="
	var found bool
	for i, v := range env {
		if strings.HasPrefix(v, evplus) {
			found = true
			env[i] = evplus + evv
		}
	}
	if !found {
		env = append(env, evplus+evv)
	}
	return env
}

// asCommandLine renders cmd as something that could be copy-and-pasted into a command line
// If cwd is not empty and different from the command's directory, prepend an approprirate "cd"
func asCommandLine(cwd string, cmd *exec.Cmd) string {
	s := "("
	if cmd.Dir != "" && cmd.Dir != cwd {
		s += "cd" + escape(cmd.Dir) + ";"
	}
	for _, e := range cmd.Env {
		if !strings.HasPrefix(e, "PATH=") &&
			!strings.HasPrefix(e, "HOME=") &&
			!strings.HasPrefix(e, "USER=") &&
			!strings.HasPrefix(e, "SHELL=") {
			s += escape(e)
		}
	}
	for _, a := range cmd.Args {
		s += escape(a)
	}
	s += " )"
	return s
}

// escape inserts escapes appropriate for use in a shell command line
func escape(s string) string {
	s = strings.Replace(s, "\\", "\\\\", -1)
	s = strings.Replace(s, "'", "\\'", -1)
	// Conservative guess at characters that will force quoting
	if strings.ContainsAny(s, "\\ ;#*&$~?!|[]()<>{}`") {
		s = " '" + s + "'"
	} else {
		s = " " + s
	}
	return s
}

func expect(want string, got tstring) {
	if want != "" {
		match, err := regexp.MatchString(want, got.o)
		if err != nil {
			panic(fmt.Sprintf("Error for regexp %s, %v\n", want, err))
		}
		if match {
			return
		}
		match, err = regexp.MatchString(want, got.e)
		if match {
			return
		}
		fmt.Printf("EXPECTED '%s'\n GOT O='%s'\nAND E='%s'\n", want, got.o, got.e)
	}
}
