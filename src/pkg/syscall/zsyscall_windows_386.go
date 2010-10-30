// mksyscall_windows.sh -l32 syscall_windows.go syscall_windows_386.go
// MACHINE GENERATED BY THE COMMAND ABOVE; DO NOT EDIT

package syscall

import "unsafe"

var (
	modkernel32 = loadDll("kernel32.dll")
	modadvapi32 = loadDll("advapi32.dll")
	modwsock32  = loadDll("wsock32.dll")
	modws2_32   = loadDll("ws2_32.dll")
	moddnsapi   = loadDll("dnsapi.dll")

	procGetLastError               = getSysProcAddr(modkernel32, "GetLastError")
	procLoadLibraryW               = getSysProcAddr(modkernel32, "LoadLibraryW")
	procFreeLibrary                = getSysProcAddr(modkernel32, "FreeLibrary")
	procGetProcAddress             = getSysProcAddr(modkernel32, "GetProcAddress")
	procGetVersion                 = getSysProcAddr(modkernel32, "GetVersion")
	procFormatMessageW             = getSysProcAddr(modkernel32, "FormatMessageW")
	procExitProcess                = getSysProcAddr(modkernel32, "ExitProcess")
	procCreateFileW                = getSysProcAddr(modkernel32, "CreateFileW")
	procReadFile                   = getSysProcAddr(modkernel32, "ReadFile")
	procWriteFile                  = getSysProcAddr(modkernel32, "WriteFile")
	procSetFilePointer             = getSysProcAddr(modkernel32, "SetFilePointer")
	procCloseHandle                = getSysProcAddr(modkernel32, "CloseHandle")
	procGetStdHandle               = getSysProcAddr(modkernel32, "GetStdHandle")
	procFindFirstFileW             = getSysProcAddr(modkernel32, "FindFirstFileW")
	procFindNextFileW              = getSysProcAddr(modkernel32, "FindNextFileW")
	procFindClose                  = getSysProcAddr(modkernel32, "FindClose")
	procGetFileInformationByHandle = getSysProcAddr(modkernel32, "GetFileInformationByHandle")
	procGetCurrentDirectoryW       = getSysProcAddr(modkernel32, "GetCurrentDirectoryW")
	procSetCurrentDirectoryW       = getSysProcAddr(modkernel32, "SetCurrentDirectoryW")
	procCreateDirectoryW           = getSysProcAddr(modkernel32, "CreateDirectoryW")
	procRemoveDirectoryW           = getSysProcAddr(modkernel32, "RemoveDirectoryW")
	procDeleteFileW                = getSysProcAddr(modkernel32, "DeleteFileW")
	procMoveFileW                  = getSysProcAddr(modkernel32, "MoveFileW")
	procGetComputerNameW           = getSysProcAddr(modkernel32, "GetComputerNameW")
	procSetEndOfFile               = getSysProcAddr(modkernel32, "SetEndOfFile")
	procGetSystemTimeAsFileTime    = getSysProcAddr(modkernel32, "GetSystemTimeAsFileTime")
	procSleep                      = getSysProcAddr(modkernel32, "Sleep")
	procGetTimeZoneInformation     = getSysProcAddr(modkernel32, "GetTimeZoneInformation")
	procCreateIoCompletionPort     = getSysProcAddr(modkernel32, "CreateIoCompletionPort")
	procGetQueuedCompletionStatus  = getSysProcAddr(modkernel32, "GetQueuedCompletionStatus")
	procCreateProcessW             = getSysProcAddr(modkernel32, "CreateProcessW")
	procGetStartupInfoW            = getSysProcAddr(modkernel32, "GetStartupInfoW")
	procGetCurrentProcess          = getSysProcAddr(modkernel32, "GetCurrentProcess")
	procDuplicateHandle            = getSysProcAddr(modkernel32, "DuplicateHandle")
	procWaitForSingleObject        = getSysProcAddr(modkernel32, "WaitForSingleObject")
	procGetTempPathW               = getSysProcAddr(modkernel32, "GetTempPathW")
	procCreatePipe                 = getSysProcAddr(modkernel32, "CreatePipe")
	procGetFileType                = getSysProcAddr(modkernel32, "GetFileType")
	procCryptAcquireContextW       = getSysProcAddr(modadvapi32, "CryptAcquireContextW")
	procCryptReleaseContext        = getSysProcAddr(modadvapi32, "CryptReleaseContext")
	procCryptGenRandom             = getSysProcAddr(modadvapi32, "CryptGenRandom")
	procOpenProcess                = getSysProcAddr(modkernel32, "OpenProcess")
	procGetExitCodeProcess         = getSysProcAddr(modkernel32, "GetExitCodeProcess")
	procGetEnvironmentStringsW     = getSysProcAddr(modkernel32, "GetEnvironmentStringsW")
	procFreeEnvironmentStringsW    = getSysProcAddr(modkernel32, "FreeEnvironmentStringsW")
	procGetEnvironmentVariableW    = getSysProcAddr(modkernel32, "GetEnvironmentVariableW")
	procSetEnvironmentVariableW    = getSysProcAddr(modkernel32, "SetEnvironmentVariableW")
	procSetFileTime                = getSysProcAddr(modkernel32, "SetFileTime")
	procGetFileAttributesW         = getSysProcAddr(modkernel32, "GetFileAttributesW")
	procWSAStartup                 = getSysProcAddr(modwsock32, "WSAStartup")
	procWSACleanup                 = getSysProcAddr(modwsock32, "WSACleanup")
	procsocket                     = getSysProcAddr(modwsock32, "socket")
	procsetsockopt                 = getSysProcAddr(modwsock32, "setsockopt")
	procbind                       = getSysProcAddr(modwsock32, "bind")
	procconnect                    = getSysProcAddr(modwsock32, "connect")
	procgetsockname                = getSysProcAddr(modwsock32, "getsockname")
	procgetpeername                = getSysProcAddr(modwsock32, "getpeername")
	proclisten                     = getSysProcAddr(modwsock32, "listen")
	procshutdown                   = getSysProcAddr(modwsock32, "shutdown")
	procAcceptEx                   = getSysProcAddr(modwsock32, "AcceptEx")
	procGetAcceptExSockaddrs       = getSysProcAddr(modwsock32, "GetAcceptExSockaddrs")
	procWSARecv                    = getSysProcAddr(modws2_32, "WSARecv")
	procWSASend                    = getSysProcAddr(modws2_32, "WSASend")
	procgethostbyname              = getSysProcAddr(modws2_32, "gethostbyname")
	procgetservbyname              = getSysProcAddr(modws2_32, "getservbyname")
	procntohs                      = getSysProcAddr(modws2_32, "ntohs")
	procDnsQuery_W                 = getSysProcAddr(moddnsapi, "DnsQuery_W")
	procDnsRecordListFree          = getSysProcAddr(moddnsapi, "DnsRecordListFree")
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
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func FreeLibrary(handle uint32) (ok bool, errno int) {
	r0, _, e1 := Syscall(procFreeLibrary, uintptr(handle), 0, 0)
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func GetProcAddress(module uint32, procname string) (proc uint32, errno int) {
	r0, _, e1 := Syscall(procGetProcAddress, uintptr(module), uintptr(unsafe.Pointer(StringBytePtr(procname))), 0)
	proc = uint32(r0)
	if proc == 0 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func GetVersion() (ver uint32, errno int) {
	r0, _, e1 := Syscall(procGetVersion, 0, 0, 0)
	ver = uint32(r0)
	if ver == 0 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
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
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
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
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
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
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
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
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func SetFilePointer(handle int32, lowoffset int32, highoffsetptr *int32, whence uint32) (newlowoffset uint32, errno int) {
	r0, _, e1 := Syscall6(procSetFilePointer, uintptr(handle), uintptr(lowoffset), uintptr(unsafe.Pointer(highoffsetptr)), uintptr(whence), 0, 0)
	newlowoffset = uint32(r0)
	if newlowoffset == 0xffffffff {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func CloseHandle(handle int32) (ok bool, errno int) {
	r0, _, e1 := Syscall(procCloseHandle, uintptr(handle), 0, 0)
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func GetStdHandle(stdhandle int32) (handle int32, errno int) {
	r0, _, e1 := Syscall(procGetStdHandle, uintptr(stdhandle), 0, 0)
	handle = int32(r0)
	if handle == -1 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func FindFirstFile(name *uint16, data *Win32finddata) (handle int32, errno int) {
	r0, _, e1 := Syscall(procFindFirstFileW, uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(data)), 0)
	handle = int32(r0)
	if handle == -1 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func FindNextFile(handle int32, data *Win32finddata) (ok bool, errno int) {
	r0, _, e1 := Syscall(procFindNextFileW, uintptr(handle), uintptr(unsafe.Pointer(data)), 0)
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func FindClose(handle int32) (ok bool, errno int) {
	r0, _, e1 := Syscall(procFindClose, uintptr(handle), 0, 0)
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func GetFileInformationByHandle(handle int32, data *ByHandleFileInformation) (ok bool, errno int) {
	r0, _, e1 := Syscall(procGetFileInformationByHandle, uintptr(handle), uintptr(unsafe.Pointer(data)), 0)
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func GetCurrentDirectory(buflen uint32, buf *uint16) (n uint32, errno int) {
	r0, _, e1 := Syscall(procGetCurrentDirectoryW, uintptr(buflen), uintptr(unsafe.Pointer(buf)), 0)
	n = uint32(r0)
	if n == 0 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func SetCurrentDirectory(path *uint16) (ok bool, errno int) {
	r0, _, e1 := Syscall(procSetCurrentDirectoryW, uintptr(unsafe.Pointer(path)), 0, 0)
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func CreateDirectory(path *uint16, sa *byte) (ok bool, errno int) {
	r0, _, e1 := Syscall(procCreateDirectoryW, uintptr(unsafe.Pointer(path)), uintptr(unsafe.Pointer(sa)), 0)
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func RemoveDirectory(path *uint16) (ok bool, errno int) {
	r0, _, e1 := Syscall(procRemoveDirectoryW, uintptr(unsafe.Pointer(path)), 0, 0)
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func DeleteFile(path *uint16) (ok bool, errno int) {
	r0, _, e1 := Syscall(procDeleteFileW, uintptr(unsafe.Pointer(path)), 0, 0)
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func MoveFile(from *uint16, to *uint16) (ok bool, errno int) {
	r0, _, e1 := Syscall(procMoveFileW, uintptr(unsafe.Pointer(from)), uintptr(unsafe.Pointer(to)), 0)
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func GetComputerName(buf *uint16, n *uint32) (ok bool, errno int) {
	r0, _, e1 := Syscall(procGetComputerNameW, uintptr(unsafe.Pointer(buf)), uintptr(unsafe.Pointer(n)), 0)
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func SetEndOfFile(handle int32) (ok bool, errno int) {
	r0, _, e1 := Syscall(procSetEndOfFile, uintptr(handle), 0, 0)
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func GetSystemTimeAsFileTime(time *Filetime) {
	Syscall(procGetSystemTimeAsFileTime, uintptr(unsafe.Pointer(time)), 0, 0)
	return
}

func sleep(msec uint32) {
	Syscall(procSleep, uintptr(msec), 0, 0)
	return
}

func GetTimeZoneInformation(tzi *Timezoneinformation) (rc uint32, errno int) {
	r0, _, e1 := Syscall(procGetTimeZoneInformation, uintptr(unsafe.Pointer(tzi)), 0, 0)
	rc = uint32(r0)
	if rc == 0xffffffff {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func CreateIoCompletionPort(filehandle int32, cphandle int32, key uint32, threadcnt uint32) (handle int32, errno int) {
	r0, _, e1 := Syscall6(procCreateIoCompletionPort, uintptr(filehandle), uintptr(cphandle), uintptr(key), uintptr(threadcnt), 0, 0)
	handle = int32(r0)
	if handle == 0 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func GetQueuedCompletionStatus(cphandle int32, qty *uint32, key *uint32, overlapped **Overlapped, timeout uint32) (ok bool, errno int) {
	r0, _, e1 := Syscall6(procGetQueuedCompletionStatus, uintptr(cphandle), uintptr(unsafe.Pointer(qty)), uintptr(unsafe.Pointer(key)), uintptr(unsafe.Pointer(overlapped)), uintptr(timeout), 0)
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func CreateProcess(appName *int16, commandLine *uint16, procSecurity *int16, threadSecurity *int16, inheritHandles bool, creationFlags uint32, env *uint16, currentDir *uint16, startupInfo *StartupInfo, outProcInfo *ProcessInformation) (ok bool, errno int) {
	var _p0 uint32
	if inheritHandles {
		_p0 = 1
	} else {
		_p0 = 0
	}
	r0, _, e1 := Syscall12(procCreateProcessW, uintptr(unsafe.Pointer(appName)), uintptr(unsafe.Pointer(commandLine)), uintptr(unsafe.Pointer(procSecurity)), uintptr(unsafe.Pointer(threadSecurity)), uintptr(_p0), uintptr(creationFlags), uintptr(unsafe.Pointer(env)), uintptr(unsafe.Pointer(currentDir)), uintptr(unsafe.Pointer(startupInfo)), uintptr(unsafe.Pointer(outProcInfo)), 0, 0)
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func GetStartupInfo(startupInfo *StartupInfo) (ok bool, errno int) {
	r0, _, e1 := Syscall(procGetStartupInfoW, uintptr(unsafe.Pointer(startupInfo)), 0, 0)
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func GetCurrentProcess() (pseudoHandle int32, errno int) {
	r0, _, e1 := Syscall(procGetCurrentProcess, 0, 0, 0)
	pseudoHandle = int32(r0)
	if pseudoHandle == 0 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func DuplicateHandle(hSourceProcessHandle int32, hSourceHandle int32, hTargetProcessHandle int32, lpTargetHandle *int32, dwDesiredAccess uint32, bInheritHandle bool, dwOptions uint32) (ok bool, errno int) {
	var _p0 uint32
	if bInheritHandle {
		_p0 = 1
	} else {
		_p0 = 0
	}
	r0, _, e1 := Syscall9(procDuplicateHandle, uintptr(hSourceProcessHandle), uintptr(hSourceHandle), uintptr(hTargetProcessHandle), uintptr(unsafe.Pointer(lpTargetHandle)), uintptr(dwDesiredAccess), uintptr(_p0), uintptr(dwOptions), 0, 0)
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func WaitForSingleObject(handle int32, waitMilliseconds uint32) (event uint32, errno int) {
	r0, _, e1 := Syscall(procWaitForSingleObject, uintptr(handle), uintptr(waitMilliseconds), 0)
	event = uint32(r0)
	if event == 0xffffffff {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func GetTempPath(buflen uint32, buf *uint16) (n uint32, errno int) {
	r0, _, e1 := Syscall(procGetTempPathW, uintptr(buflen), uintptr(unsafe.Pointer(buf)), 0)
	n = uint32(r0)
	if n == 0 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func CreatePipe(readhandle *uint32, writehandle *uint32, lpsa *byte, size uint32) (ok bool, errno int) {
	r0, _, e1 := Syscall6(procCreatePipe, uintptr(unsafe.Pointer(readhandle)), uintptr(unsafe.Pointer(writehandle)), uintptr(unsafe.Pointer(lpsa)), uintptr(size), 0, 0)
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func GetFileType(filehandle uint32) (n uint32, errno int) {
	r0, _, e1 := Syscall(procGetFileType, uintptr(filehandle), 0, 0)
	n = uint32(r0)
	if n == 0 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func CryptAcquireContext(provhandle *uint32, container *uint16, provider *uint16, provtype uint32, flags uint32) (ok bool, errno int) {
	r0, _, e1 := Syscall6(procCryptAcquireContextW, uintptr(unsafe.Pointer(provhandle)), uintptr(unsafe.Pointer(container)), uintptr(unsafe.Pointer(provider)), uintptr(provtype), uintptr(flags), 0)
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func CryptReleaseContext(provhandle uint32, flags uint32) (ok bool, errno int) {
	r0, _, e1 := Syscall(procCryptReleaseContext, uintptr(provhandle), uintptr(flags), 0)
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func CryptGenRandom(provhandle uint32, buflen uint32, buf *byte) (ok bool, errno int) {
	r0, _, e1 := Syscall(procCryptGenRandom, uintptr(provhandle), uintptr(buflen), uintptr(unsafe.Pointer(buf)))
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func OpenProcess(da uint32, b int, pid uint32) (handle uint32, errno int) {
	r0, _, e1 := Syscall(procOpenProcess, uintptr(da), uintptr(b), uintptr(pid))
	handle = uint32(r0)
	if handle == 0 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func GetExitCodeProcess(h uint32, c *uint32) (ok bool, errno int) {
	r0, _, e1 := Syscall(procGetExitCodeProcess, uintptr(h), uintptr(unsafe.Pointer(c)), 0)
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func GetEnvironmentStrings() (envs *uint16, errno int) {
	r0, _, e1 := Syscall(procGetEnvironmentStringsW, 0, 0, 0)
	envs = (*uint16)(unsafe.Pointer(r0))
	if envs == nil {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func FreeEnvironmentStrings(envs *uint16) (ok bool, errno int) {
	r0, _, e1 := Syscall(procFreeEnvironmentStringsW, uintptr(unsafe.Pointer(envs)), 0, 0)
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func GetEnvironmentVariable(name *uint16, buffer *uint16, size uint32) (n uint32, errno int) {
	r0, _, e1 := Syscall(procGetEnvironmentVariableW, uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(buffer)), uintptr(size))
	n = uint32(r0)
	if n == 0 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func SetEnvironmentVariable(name *uint16, value *uint16) (ok bool, errno int) {
	r0, _, e1 := Syscall(procSetEnvironmentVariableW, uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(value)), 0)
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func SetFileTime(handle int32, ctime *Filetime, atime *Filetime, wtime *Filetime) (ok bool, errno int) {
	r0, _, e1 := Syscall6(procSetFileTime, uintptr(handle), uintptr(unsafe.Pointer(ctime)), uintptr(unsafe.Pointer(atime)), uintptr(unsafe.Pointer(wtime)), 0, 0)
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func GetFileAttributes(name *uint16) (attrs uint32, errno int) {
	r0, _, e1 := Syscall(procGetFileAttributesW, uintptr(unsafe.Pointer(name)), 0, 0)
	attrs = uint32(r0)
	if attrs == INVALID_FILE_ATTRIBUTES {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func WSAStartup(verreq uint32, data *WSAData) (sockerrno int) {
	r0, _, _ := Syscall(procWSAStartup, uintptr(verreq), uintptr(unsafe.Pointer(data)), 0)
	sockerrno = int(r0)
	return
}

func WSACleanup() (errno int) {
	r1, _, e1 := Syscall(procWSACleanup, 0, 0, 0)
	if int(r1) == -1 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func socket(af int32, typ int32, protocol int32) (handle int32, errno int) {
	r0, _, e1 := Syscall(procsocket, uintptr(af), uintptr(typ), uintptr(protocol))
	handle = int32(r0)
	if handle == -1 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func setsockopt(s int32, level int32, optname int32, optval *byte, optlen int32) (errno int) {
	r1, _, e1 := Syscall6(procsetsockopt, uintptr(s), uintptr(level), uintptr(optname), uintptr(unsafe.Pointer(optval)), uintptr(optlen), 0)
	if int(r1) == -1 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func bind(s int32, name uintptr, namelen int32) (errno int) {
	r1, _, e1 := Syscall(procbind, uintptr(s), uintptr(name), uintptr(namelen))
	if int(r1) == -1 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func connect(s int32, name uintptr, namelen int32) (errno int) {
	r1, _, e1 := Syscall(procconnect, uintptr(s), uintptr(name), uintptr(namelen))
	if int(r1) == -1 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func getsockname(s int32, rsa *RawSockaddrAny, addrlen *int32) (errno int) {
	r1, _, e1 := Syscall(procgetsockname, uintptr(s), uintptr(unsafe.Pointer(rsa)), uintptr(unsafe.Pointer(addrlen)))
	if int(r1) == -1 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func getpeername(s int32, rsa *RawSockaddrAny, addrlen *int32) (errno int) {
	r1, _, e1 := Syscall(procgetpeername, uintptr(s), uintptr(unsafe.Pointer(rsa)), uintptr(unsafe.Pointer(addrlen)))
	if int(r1) == -1 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func listen(s int32, backlog int32) (errno int) {
	r1, _, e1 := Syscall(proclisten, uintptr(s), uintptr(backlog), 0)
	if int(r1) == -1 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func shutdown(s int32, how int32) (errno int) {
	r1, _, e1 := Syscall(procshutdown, uintptr(s), uintptr(how), 0)
	if int(r1) == -1 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func AcceptEx(ls uint32, as uint32, buf *byte, rxdatalen uint32, laddrlen uint32, raddrlen uint32, recvd *uint32, overlapped *Overlapped) (ok bool, errno int) {
	r0, _, e1 := Syscall9(procAcceptEx, uintptr(ls), uintptr(as), uintptr(unsafe.Pointer(buf)), uintptr(rxdatalen), uintptr(laddrlen), uintptr(raddrlen), uintptr(unsafe.Pointer(recvd)), uintptr(unsafe.Pointer(overlapped)), 0)
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func GetAcceptExSockaddrs(buf *byte, rxdatalen uint32, laddrlen uint32, raddrlen uint32, lrsa **RawSockaddrAny, lrsalen *int32, rrsa **RawSockaddrAny, rrsalen *int32) {
	Syscall9(procGetAcceptExSockaddrs, uintptr(unsafe.Pointer(buf)), uintptr(rxdatalen), uintptr(laddrlen), uintptr(raddrlen), uintptr(unsafe.Pointer(lrsa)), uintptr(unsafe.Pointer(lrsalen)), uintptr(unsafe.Pointer(rrsa)), uintptr(unsafe.Pointer(rrsalen)), 0)
	return
}

func WSARecv(s uint32, bufs *WSABuf, bufcnt uint32, recvd *uint32, flags *uint32, overlapped *Overlapped, croutine *byte) (errno int) {
	r1, _, e1 := Syscall9(procWSARecv, uintptr(s), uintptr(unsafe.Pointer(bufs)), uintptr(bufcnt), uintptr(unsafe.Pointer(recvd)), uintptr(unsafe.Pointer(flags)), uintptr(unsafe.Pointer(overlapped)), uintptr(unsafe.Pointer(croutine)), 0, 0)
	if int(r1) == -1 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func WSASend(s uint32, bufs *WSABuf, bufcnt uint32, sent *uint32, flags uint32, overlapped *Overlapped, croutine *byte) (errno int) {
	r1, _, e1 := Syscall9(procWSASend, uintptr(s), uintptr(unsafe.Pointer(bufs)), uintptr(bufcnt), uintptr(unsafe.Pointer(sent)), uintptr(flags), uintptr(unsafe.Pointer(overlapped)), uintptr(unsafe.Pointer(croutine)), 0, 0)
	if int(r1) == -1 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func GetHostByName(name string) (h *Hostent, errno int) {
	r0, _, e1 := Syscall(procgethostbyname, uintptr(unsafe.Pointer(StringBytePtr(name))), 0, 0)
	h = (*Hostent)(unsafe.Pointer(r0))
	if h == nil {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func GetServByName(name string, proto string) (s *Servent, errno int) {
	r0, _, e1 := Syscall(procgetservbyname, uintptr(unsafe.Pointer(StringBytePtr(name))), uintptr(unsafe.Pointer(StringBytePtr(proto))), 0)
	s = (*Servent)(unsafe.Pointer(r0))
	if s == nil {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func Ntohs(netshort uint16) (u uint16) {
	r0, _, _ := Syscall(procntohs, uintptr(netshort), 0, 0)
	u = uint16(r0)
	return
}

func DnsQuery(name string, qtype uint16, options uint32, extra *byte, qrs **DNSRecord, pr *byte) (status uint32) {
	r0, _, _ := Syscall6(procDnsQuery_W, uintptr(unsafe.Pointer(StringToUTF16Ptr(name))), uintptr(qtype), uintptr(options), uintptr(unsafe.Pointer(extra)), uintptr(unsafe.Pointer(qrs)), uintptr(unsafe.Pointer(pr)))
	status = uint32(r0)
	return
}

func DnsRecordListFree(rl *DNSRecord, freetype uint32) {
	Syscall(procDnsRecordListFree, uintptr(unsafe.Pointer(rl)), uintptr(freetype), 0)
	return
}
