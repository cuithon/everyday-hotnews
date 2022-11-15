// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main_test

import (
	"bytes"
	"cmd/go/internal/script"
	"flag"
	"internal/diff"
	"internal/testenv"
	"os"
	"strings"
	"testing"
	"text/template"
)

var fixReadme = flag.Bool("fixreadme", false, "if true, update ../testdata/script/README")

func checkScriptReadme(t *testing.T, engine *script.Engine, env []string) {
	var args struct {
		Language   string
		Commands   string
		Conditions string
	}

	cmds := new(strings.Builder)
	if err := engine.ListCmds(cmds, true); err != nil {
		t.Fatal(err)
	}
	args.Commands = cmds.String()

	conds := new(strings.Builder)
	if err := engine.ListConds(conds, nil); err != nil {
		t.Fatal(err)
	}
	args.Conditions = conds.String()

	if !testenv.HasExec() {
		t.Skipf("updating script README requires os/exec")
	}

	doc := new(strings.Builder)
	cmd := testenv.Command(t, testGo, "doc", "cmd/go/internal/script")
	cmd.Env = env
	cmd.Stdout = doc
	if err := cmd.Run(); err != nil {
		t.Fatal(cmd, ":", err)
	}
	_, lang, ok := strings.Cut(doc.String(), "# Script Language\n\n")
	if !ok {
		t.Fatalf("%q did not include Script Language section", cmd)
	}
	lang, _, ok = strings.Cut(lang, "\n\nvar ")
	if !ok {
		t.Fatalf("%q did not include vars after Script Language section", cmd)
	}
	args.Language = lang

	tmpl := template.Must(template.New("README").Parse(readmeTmpl[1:]))
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, args); err != nil {
		t.Fatal(err)
	}

	const readmePath = "testdata/script/README"
	old, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatal(err)
	}
	diff := diff.Diff(readmePath, old, "readmeTmpl", buf.Bytes())
	if diff == nil {
		t.Logf("%s is up to date.", readmePath)
		return
	}

	if *fixReadme {
		if err := os.WriteFile(readmePath, buf.Bytes(), 0666); err != nil {
			t.Fatal(err)
		}
		t.Logf("wrote %d bytes to %s", buf.Len(), readmePath)
	} else {
		t.Logf("\n%s", diff)
		t.Errorf("%s is stale. To update, run 'go generate cmd/go'.", readmePath)
	}
}

const readmeTmpl = `
This file is generated by 'go generate cmd/go'. DO NOT EDIT.

This directory holds test scripts *.txt run during 'go test cmd/go'.
To run a specific script foo.txt

	go test cmd/go -run=Script/^foo$

In general script files should have short names: a few words, not whole sentences.
The first word should be the general category of behavior being tested,
often the name of a go subcommand (list, build, test, ...) or concept (vendor, pattern).

Each script is a text archive (go doc internal/txtar).
The script begins with an actual command script to run
followed by the content of zero or more supporting files to
create in the script's temporary file system before it starts executing.

As an example, run_hello.txt says:

	# hello world
	go run hello.go
	stderr 'hello world'
	! stdout .

	-- hello.go --
	package main
	func main() { println("hello world") }

Each script runs in a fresh temporary work directory tree, available to scripts as $WORK.
Scripts also have access to other environment variables, including:

	GOARCH=<target GOARCH>
	GOCACHE=<actual GOCACHE being used outside the test>
	GOEXE=<executable file suffix: .exe on Windows, empty on other systems>
	GOOS=<target GOOS>
	GOPATH=$WORK/gopath
	GOPROXY=<local module proxy serving from cmd/go/testdata/mod>
	GOROOT=<actual GOROOT>
	GOROOT_FINAL=<actual GOROOT_FINAL>
	TESTGO_GOROOT=<GOROOT used to build cmd/go, for use in tests that may change GOROOT>
	HOME=/no-home
	PATH=<actual PATH>
	TMPDIR=$WORK/tmp
	GODEBUG=<actual GODEBUG>
	devnull=<value of os.DevNull>
	goversion=<current Go version; for example, 1.12>

On Plan 9, the variables $path and $home are set instead of $PATH and $HOME.
On Windows, the variables $USERPROFILE and $TMP are set instead of
$HOME and $TMPDIR.

The lines at the top of the script are a sequence of commands to be executed by
a small script engine configured in ../../script_test.go (not the system shell).

The scripts' supporting files are unpacked relative to $GOPATH/src
(aka $WORK/gopath/src) and then the script begins execution in that directory as
well. Thus the example above runs in $WORK/gopath/src with GOPATH=$WORK/gopath
and $WORK/gopath/src/hello.go containing the listed contents.

{{.Language}}

When TestScript runs a script and the script fails, by default TestScript shows
the execution of the most recent phase of the script (since the last # comment)
and only shows the # comments for earlier phases. For example, here is a
multi-phase script with a bug in it:

	# GOPATH with p1 in d2, p2 in d2
	env GOPATH=$WORK${/}d1${:}$WORK${/}d2

	# build & install p1
	env
	go install -i p1
	! stale p1
	! stale p2

	# modify p2 - p1 should appear stale
	cp $WORK/p2x.go $WORK/d2/src/p2/p2.go
	stale p1 p2

	# build & install p1 again
	go install -i p11
	! stale p1
	! stale p2

	-- $WORK/d1/src/p1/p1.go --
	package p1
	import "p2"
	func F() { p2.F() }
	-- $WORK/d2/src/p2/p2.go --
	package p2
	func F() {}
	-- $WORK/p2x.go --
	package p2
	func F() {}
	func G() {}

The bug is that the final phase installs p11 instead of p1. The test failure looks like:

	$ go test -run=Script
	--- FAIL: TestScript (3.75s)
	    --- FAIL: TestScript/install_rebuild_gopath (0.16s)
	        script_test.go:223:
	            # GOPATH with p1 in d2, p2 in d2 (0.000s)
	            # build & install p1 (0.087s)
	            # modify p2 - p1 should appear stale (0.029s)
	            # build & install p1 again (0.022s)
	            > go install -i p11
	            [stderr]
	            can't load package: package p11: cannot find package "p11" in any of:
	            	/Users/rsc/go/src/p11 (from $GOROOT)
	            	$WORK/d1/src/p11 (from $GOPATH)
	            	$WORK/d2/src/p11
	            [exit status 1]
	            FAIL: unexpected go command failure

	        script_test.go:73: failed at testdata/script/install_rebuild_gopath.txt:15 in $WORK/gopath/src

	FAIL
	exit status 1
	FAIL	cmd/go	4.875s
	$

Note that the commands in earlier phases have been hidden, so that the relevant
commands are more easily found, and the elapsed time for a completed phase
is shown next to the phase heading. To see the entire execution, use "go test -v",
which also adds an initial environment dump to the beginning of the log.

Note also that in reported output, the actual name of the per-script temporary directory
has been consistently replaced with the literal string $WORK.

The cmd/go test flag -testwork (which must appear on the "go test" command line after
standard test flags) causes each test to log the name of its $WORK directory and other
environment variable settings and also to leave that directory behind when it exits,
for manual debugging of failing tests:

	$ go test -run=Script -work
	--- FAIL: TestScript (3.75s)
	    --- FAIL: TestScript/install_rebuild_gopath (0.16s)
	        script_test.go:223:
	            WORK=/tmp/cmd-go-test-745953508/script-install_rebuild_gopath
	            GOARCH=
	            GOCACHE=/Users/rsc/Library/Caches/go-build
	            GOOS=
	            GOPATH=$WORK/gopath
	            GOROOT=/Users/rsc/go
	            HOME=/no-home
	            TMPDIR=$WORK/tmp
	            exe=

	            # GOPATH with p1 in d2, p2 in d2 (0.000s)
	            # build & install p1 (0.085s)
	            # modify p2 - p1 should appear stale (0.030s)
	            # build & install p1 again (0.019s)
	            > go install -i p11
	            [stderr]
	            can't load package: package p11: cannot find package "p11" in any of:
	            	/Users/rsc/go/src/p11 (from $GOROOT)
	            	$WORK/d1/src/p11 (from $GOPATH)
	            	$WORK/d2/src/p11
	            [exit status 1]
	            FAIL: unexpected go command failure

	        script_test.go:73: failed at testdata/script/install_rebuild_gopath.txt:15 in $WORK/gopath/src

	FAIL
	exit status 1
	FAIL	cmd/go	4.875s
	$

	$ WORK=/tmp/cmd-go-test-745953508/script-install_rebuild_gopath
	$ cd $WORK/d1/src/p1
	$ cat p1.go
	package p1
	import "p2"
	func F() { p2.F() }
	$

The available commands are:
{{.Commands}}

The available conditions are:
{{.Conditions}}
`
