// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Fork, exec, wait, etc.

package syscall

import (
	"sync"
	"unsafe"
	"utf16"
)

var ForkLock sync.RWMutex

// EscapeArg rewrites command line argument s as prescribed
// in http://msdn.microsoft.com/en-us/library/ms880421.
// This function returns "" (2 double quotes) if s is empty.
// Alternatively, these transformations are done:
// - every back slash (\) is doubled, but only if immediately
//   followed by double quote (");
// - every double quote (") is escaped by back slash (\);
// - finally, s is wrapped with double quotes (arg -> "arg"),
//   but only if there is space or tab inside s.
func EscapeArg(s string) string {
	if len(s) == 0 {
		return "\"\""
	}
	n := len(s)
	hasSpace := false
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"', '\\':
			n++
		case ' ', '\t':
			hasSpace = true
		}
	}
	if hasSpace {
		n += 2
	}
	if n == len(s) {
		return s
	}

	qs := make([]byte, n)
	j := 0
	if hasSpace {
		qs[j] = '"'
		j++
	}
	slashes := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		default:
			slashes = 0
			qs[j] = s[i]
		case '\\':
			slashes++
			qs[j] = s[i]
		case '"':
			for ; slashes > 0; slashes-- {
				qs[j] = '\\'
				j++
			}
			qs[j] = '\\'
			j++
			qs[j] = s[i]
		}
		j++
	}
	if hasSpace {
		for ; slashes > 0; slashes-- {
			qs[j] = '\\'
			j++
		}
		qs[j] = '"'
		j++
	}
	return string(qs[:j])
}

// makeCmdLine builds a command line out of args by escaping "special"
// characters and joining the arguments with spaces.
func makeCmdLine(args []string) string {
	var s string
	for _, v := range args {
		if s != "" {
			s += " "
		}
		s += EscapeArg(v)
	}
	return s
}

// createEnvBlock converts an array of environment strings into
// the representation required by CreateProcess: a sequence of NUL
// terminated strings followed by a nil.
// Last bytes are two UCS-2 NULs, or four NUL bytes.
func createEnvBlock(envv []string) *uint16 {
	if len(envv) == 0 {
		return &utf16.Encode([]int("\x00\x00"))[0]
	}
	length := 0
	for _, s := range envv {
		length += len(s) + 1
	}
	length += 1

	b := make([]byte, length)
	i := 0
	for _, s := range envv {
		l := len(s)
		copy(b[i:i+l], []byte(s))
		copy(b[i+l:i+l+1], []byte{0})
		i = i + l + 1
	}
	copy(b[i:i+1], []byte{0})

	return &utf16.Encode([]int(string(b)))[0]
}

func CloseOnExec(fd int) {
	SetHandleInformation(int32(fd), HANDLE_FLAG_INHERIT, 0)
}

func SetNonblock(fd int, nonblocking bool) (errno int) {
	return 0
}

// getFullPath retrieves the full path of the specified file.
// Just a wrapper for Windows GetFullPathName api.
func getFullPath(name string) (path string, err int) {
	p := StringToUTF16Ptr(name)
	buf := make([]uint16, 100)
	n, err := GetFullPathName(p, uint32(len(buf)), &buf[0], nil)
	if err != 0 {
		return "", err
	}
	if n > uint32(len(buf)) {
		// Windows is asking for bigger buffer.
		buf = make([]uint16, n)
		n, err = GetFullPathName(p, uint32(len(buf)), &buf[0], nil)
		if err != 0 {
			return "", err
		}
		if n > uint32(len(buf)) {
			return "", EINVAL
		}
	}
	return UTF16ToString(buf[:n]), 0
}

func isSlash(c uint8) bool {
	return c == '\\' || c == '/'
}

func normalizeDir(dir string) (name string, err int) {
	ndir, err := getFullPath(dir)
	if err != 0 {
		return "", err
	}
	if len(ndir) > 2 && isSlash(ndir[0]) && isSlash(ndir[1]) {
		// dir cannot have \\server\share\path form
		return "", EINVAL
	}
	return ndir, 0
}

func volToUpper(ch int) int {
	if 'a' <= ch && ch <= 'z' {
		ch += 'A' - 'a'
	}
	return ch
}

func joinExeDirAndFName(dir, p string) (name string, err int) {
	if len(p) == 0 {
		return "", EINVAL
	}
	if len(p) > 2 && isSlash(p[0]) && isSlash(p[1]) {
		// \\server\share\path form
		return p, 0
	}
	if len(p) > 1 && p[1] == ':' {
		// has drive letter
		if len(p) == 2 {
			return "", EINVAL
		}
		if isSlash(p[2]) {
			return p, 0
		} else {
			d, err := normalizeDir(dir)
			if err != 0 {
				return "", err
			}
			if volToUpper(int(p[0])) == volToUpper(int(d[0])) {
				return getFullPath(d + "\\" + p[2:])
			} else {
				return getFullPath(p)
			}
		}
	} else {
		// no drive letter
		d, err := normalizeDir(dir)
		if err != 0 {
			return "", err
		}
		if isSlash(p[0]) {
			return getFullPath(d[:2] + p)
		} else {
			return getFullPath(d + "\\" + p)
		}
	}
	// we shouldn't be here
	return "", EINVAL
}

type ProcAttr struct {
	Dir   string
	Env   []string
	Files []int
	Sys   *SysProcAttr
}

type SysProcAttr struct {
	HideWindow bool
	CmdLine    string // used if non-empty, else the windows command line is built by escaping the arguments passed to StartProcess
}

var zeroProcAttr ProcAttr
var zeroSysProcAttr SysProcAttr

func StartProcess(argv0 string, argv []string, attr *ProcAttr) (pid, handle int, err int) {
	if len(argv0) == 0 {
		return 0, 0, EWINDOWS
	}
	if attr == nil {
		attr = &zeroProcAttr
	}
	sys := attr.Sys
	if sys == nil {
		sys = &zeroSysProcAttr
	}

	if len(attr.Files) > 3 {
		return 0, 0, EWINDOWS
	}

	if len(attr.Dir) != 0 {
		// StartProcess assumes that argv0 is relative to attr.Dir,
		// because it implies Chdir(attr.Dir) before executing argv0.
		// Windows CreateProcess assumes the opposite: it looks for
		// argv0 relative to the current directory, and, only once the new
		// process is started, it does Chdir(attr.Dir). We are adjusting
		// for that difference here by making argv0 absolute.
		var err int
		argv0, err = joinExeDirAndFName(attr.Dir, argv0)
		if err != 0 {
			return 0, 0, err
		}
	}
	argv0p := StringToUTF16Ptr(argv0)

	var cmdline string
	// Windows CreateProcess takes the command line as a single string:
	// use attr.CmdLine if set, else build the command line by escaping
	// and joining each argument with spaces
	if sys.CmdLine != "" {
		cmdline = sys.CmdLine
	} else {
		cmdline = makeCmdLine(argv)
	}

	var argvp *uint16
	if len(cmdline) != 0 {
		argvp = StringToUTF16Ptr(cmdline)
	}

	var dirp *uint16
	if len(attr.Dir) != 0 {
		dirp = StringToUTF16Ptr(attr.Dir)
	}

	// Acquire the fork lock so that no other threads
	// create new fds that are not yet close-on-exec
	// before we fork.
	ForkLock.Lock()
	defer ForkLock.Unlock()

	p, _ := GetCurrentProcess()
	fd := make([]int32, len(attr.Files))
	for i := range attr.Files {
		if attr.Files[i] > 0 {
			err := DuplicateHandle(p, int32(attr.Files[i]), p, &fd[i], 0, true, DUPLICATE_SAME_ACCESS)
			if err != 0 {
				return 0, 0, err
			}
			defer CloseHandle(int32(fd[i]))
		}
	}
	si := new(StartupInfo)
	si.Cb = uint32(unsafe.Sizeof(*si))
	si.Flags = STARTF_USESTDHANDLES
	if sys.HideWindow {
		si.Flags |= STARTF_USESHOWWINDOW
		si.ShowWindow = SW_HIDE
	}
	si.StdInput = fd[0]
	si.StdOutput = fd[1]
	si.StdErr = fd[2]

	pi := new(ProcessInformation)

	err = CreateProcess(argv0p, argvp, nil, nil, true, CREATE_UNICODE_ENVIRONMENT, createEnvBlock(attr.Env), dirp, si, pi)
	if err != 0 {
		return 0, 0, err
	}
	defer CloseHandle(pi.Thread)

	return int(pi.ProcessId), int(pi.Process), 0
}

func Exec(argv0 string, argv []string, envv []string) (err int) {
	return EWINDOWS
}
