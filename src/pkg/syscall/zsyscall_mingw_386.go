// mksyscall_mingw.sh -l32 syscall_mingw.go syscall_mingw_386.go
// MACHINE GENERATED BY THE COMMAND ABOVE; DO NOT EDIT

package syscall

import "unsafe"

var (
	modKERNEL32        = loadDll("kernel32.dll")
	procGetLastError   = getSysProcAddr(modKERNEL32, "GetLastError")
	procLoadLibraryW   = getSysProcAddr(modKERNEL32, "LoadLibraryW")
	procFreeLibrary    = getSysProcAddr(modKERNEL32, "FreeLibrary")
	procGetProcAddress = getSysProcAddr(modKERNEL32, "GetProcAddress")
	procGetVersion     = getSysProcAddr(modKERNEL32, "GetVersion")
	procFormatMessageW = getSysProcAddr(modKERNEL32, "FormatMessageW")
	procExitProcess    = getSysProcAddr(modKERNEL32, "ExitProcess")
	procCreateFileW    = getSysProcAddr(modKERNEL32, "CreateFileW")
	procReadFile       = getSysProcAddr(modKERNEL32, "ReadFile")
	procWriteFile      = getSysProcAddr(modKERNEL32, "WriteFile")
	procSetFilePointer = getSysProcAddr(modKERNEL32, "SetFilePointer")
	procCloseHandle    = getSysProcAddr(modKERNEL32, "CloseHandle")
	procGetStdHandle   = getSysProcAddr(modKERNEL32, "GetStdHandle")
)

func GetLastError() (lasterrno int) {
	r0, _, _ := Syscall(procGetLastError, 0, 0, 0)
	lasterrno = int(r0)
	return
}

func LoadLibrary(libname string) (handle uint32, errno int) {
	r0, _, e1 := Syscall(procLoadLibraryW, uintptr(unsafe.Pointer(StringToUTF16Ptr(libname))), 0, 0)
	handle = uint32(r0)
	if handle == 0 {
		errno = int(e1)
	} else {
		errno = 0
	}
	return
}

func FreeLibrary(handle uint32) (ok bool, errno int) {
	r0, _, e1 := Syscall(procFreeLibrary, uintptr(handle), 0, 0)
	ok = bool(r0 != 0)
	if !ok {
		errno = int(e1)
	} else {
		errno = 0
	}
	return
}

func GetProcAddress(module uint32, procname string) (proc uint32, errno int) {
	r0, _, e1 := Syscall(procGetProcAddress, uintptr(module), uintptr(unsafe.Pointer(StringBytePtr(procname))), 0)
	proc = uint32(r0)
	if proc == 0 {
		errno = int(e1)
	} else {
		errno = 0
	}
	return
}

func GetVersion() (ver uint32, errno int) {
	r0, _, e1 := Syscall(procGetVersion, 0, 0, 0)
	ver = uint32(r0)
	if ver == 0 {
		errno = int(e1)
	} else {
		errno = 0
	}
	return
}

func FormatMessage(flags uint32, msgsrc uint32, msgid uint32, langid uint32, buf []uint16, args *byte) (n uint32, errno int) {
	var _p0 *uint16
	if len(buf) > 0 {
		_p0 = &buf[0]
	}
	r0, _, e1 := Syscall9(procFormatMessageW, uintptr(flags), uintptr(msgsrc), uintptr(msgid), uintptr(langid), uintptr(unsafe.Pointer(_p0)), uintptr(len(buf)), uintptr(unsafe.Pointer(args)), 0, 0)
	n = uint32(r0)
	if n == 0 {
		errno = int(e1)
	} else {
		errno = 0
	}
	return
}

func ExitProcess(exitcode uint32) {
	Syscall(procExitProcess, uintptr(exitcode), 0, 0)
	return
}

func CreateFile(name *uint16, access uint32, mode uint32, sa *byte, createmode uint32, attrs uint32, templatefile int32) (handle int32, errno int) {
	r0, _, e1 := Syscall9(procCreateFileW, uintptr(unsafe.Pointer(name)), uintptr(access), uintptr(mode), uintptr(unsafe.Pointer(sa)), uintptr(createmode), uintptr(attrs), uintptr(templatefile), 0, 0)
	handle = int32(r0)
	if handle == -1 {
		errno = int(e1)
	} else {
		errno = 0
	}
	return
}

func ReadFile(handle int32, buf []byte, done *uint32, overlapped *Overlapped) (ok bool, errno int) {
	var _p0 *byte
	if len(buf) > 0 {
		_p0 = &buf[0]
	}
	r0, _, e1 := Syscall6(procReadFile, uintptr(handle), uintptr(unsafe.Pointer(_p0)), uintptr(len(buf)), uintptr(unsafe.Pointer(done)), uintptr(unsafe.Pointer(overlapped)), 0)
	ok = bool(r0 != 0)
	if !ok {
		errno = int(e1)
	} else {
		errno = 0
	}
	return
}

func WriteFile(handle int32, buf []byte, done *uint32, overlapped *Overlapped) (ok bool, errno int) {
	var _p0 *byte
	if len(buf) > 0 {
		_p0 = &buf[0]
	}
	r0, _, e1 := Syscall6(procWriteFile, uintptr(handle), uintptr(unsafe.Pointer(_p0)), uintptr(len(buf)), uintptr(unsafe.Pointer(done)), uintptr(unsafe.Pointer(overlapped)), 0)
	ok = bool(r0 != 0)
	if !ok {
		errno = int(e1)
	} else {
		errno = 0
	}
	return
}

func SetFilePointer(handle int32, lowoffset int32, highoffsetptr *int32, whence uint32) (newlowoffset uint32, errno int) {
	r0, _, e1 := Syscall6(procSetFilePointer, uintptr(handle), uintptr(lowoffset), uintptr(unsafe.Pointer(highoffsetptr)), uintptr(whence), 0, 0)
	newlowoffset = uint32(r0)
	if newlowoffset == 0xffffffff {
		errno = int(e1)
	} else {
		errno = 0
	}
	return
}

func CloseHandle(handle int32) (ok bool, errno int) {
	r0, _, e1 := Syscall(procCloseHandle, uintptr(handle), 0, 0)
	ok = bool(r0 != 0)
	if !ok {
		errno = int(e1)
	} else {
		errno = 0
	}
	return
}

func GetStdHandle(stdhandle int32) (handle int32, errno int) {
	r0, _, e1 := Syscall(procGetStdHandle, uintptr(stdhandle), 0, 0)
	handle = int32(r0)
	if handle == -1 {
		errno = int(e1)
	} else {
		errno = 0
	}
	return
}
