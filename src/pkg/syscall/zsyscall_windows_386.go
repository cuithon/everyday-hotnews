// mksyscall_windows.pl -l32 syscall_windows.go syscall_windows_386.go
// MACHINE GENERATED BY THE COMMAND ABOVE; DO NOT EDIT

package syscall

import "unsafe"

var (
	modkernel32 = NewLazyDLL("kernel32.dll")
	modadvapi32 = NewLazyDLL("advapi32.dll")
	modshell32  = NewLazyDLL("shell32.dll")
	modmswsock  = NewLazyDLL("mswsock.dll")
	modcrypt32  = NewLazyDLL("crypt32.dll")
	modws2_32   = NewLazyDLL("ws2_32.dll")
	moddnsapi   = NewLazyDLL("dnsapi.dll")
	modiphlpapi = NewLazyDLL("iphlpapi.dll")

	procGetLastError                = modkernel32.NewProc("GetLastError")
	procLoadLibraryW                = modkernel32.NewProc("LoadLibraryW")
	procFreeLibrary                 = modkernel32.NewProc("FreeLibrary")
	procGetProcAddress              = modkernel32.NewProc("GetProcAddress")
	procGetVersion                  = modkernel32.NewProc("GetVersion")
	procFormatMessageW              = modkernel32.NewProc("FormatMessageW")
	procExitProcess                 = modkernel32.NewProc("ExitProcess")
	procCreateFileW                 = modkernel32.NewProc("CreateFileW")
	procReadFile                    = modkernel32.NewProc("ReadFile")
	procWriteFile                   = modkernel32.NewProc("WriteFile")
	procSetFilePointer              = modkernel32.NewProc("SetFilePointer")
	procCloseHandle                 = modkernel32.NewProc("CloseHandle")
	procGetStdHandle                = modkernel32.NewProc("GetStdHandle")
	procFindFirstFileW              = modkernel32.NewProc("FindFirstFileW")
	procFindNextFileW               = modkernel32.NewProc("FindNextFileW")
	procFindClose                   = modkernel32.NewProc("FindClose")
	procGetFileInformationByHandle  = modkernel32.NewProc("GetFileInformationByHandle")
	procGetCurrentDirectoryW        = modkernel32.NewProc("GetCurrentDirectoryW")
	procSetCurrentDirectoryW        = modkernel32.NewProc("SetCurrentDirectoryW")
	procCreateDirectoryW            = modkernel32.NewProc("CreateDirectoryW")
	procRemoveDirectoryW            = modkernel32.NewProc("RemoveDirectoryW")
	procDeleteFileW                 = modkernel32.NewProc("DeleteFileW")
	procMoveFileW                   = modkernel32.NewProc("MoveFileW")
	procGetComputerNameW            = modkernel32.NewProc("GetComputerNameW")
	procSetEndOfFile                = modkernel32.NewProc("SetEndOfFile")
	procGetSystemTimeAsFileTime     = modkernel32.NewProc("GetSystemTimeAsFileTime")
	procSleep                       = modkernel32.NewProc("Sleep")
	procGetTimeZoneInformation      = modkernel32.NewProc("GetTimeZoneInformation")
	procCreateIoCompletionPort      = modkernel32.NewProc("CreateIoCompletionPort")
	procGetQueuedCompletionStatus   = modkernel32.NewProc("GetQueuedCompletionStatus")
	procPostQueuedCompletionStatus  = modkernel32.NewProc("PostQueuedCompletionStatus")
	procCancelIo                    = modkernel32.NewProc("CancelIo")
	procCreateProcessW              = modkernel32.NewProc("CreateProcessW")
	procOpenProcess                 = modkernel32.NewProc("OpenProcess")
	procTerminateProcess            = modkernel32.NewProc("TerminateProcess")
	procGetExitCodeProcess          = modkernel32.NewProc("GetExitCodeProcess")
	procGetStartupInfoW             = modkernel32.NewProc("GetStartupInfoW")
	procGetCurrentProcess           = modkernel32.NewProc("GetCurrentProcess")
	procDuplicateHandle             = modkernel32.NewProc("DuplicateHandle")
	procWaitForSingleObject         = modkernel32.NewProc("WaitForSingleObject")
	procGetTempPathW                = modkernel32.NewProc("GetTempPathW")
	procCreatePipe                  = modkernel32.NewProc("CreatePipe")
	procGetFileType                 = modkernel32.NewProc("GetFileType")
	procCryptAcquireContextW        = modadvapi32.NewProc("CryptAcquireContextW")
	procCryptReleaseContext         = modadvapi32.NewProc("CryptReleaseContext")
	procCryptGenRandom              = modadvapi32.NewProc("CryptGenRandom")
	procGetEnvironmentStringsW      = modkernel32.NewProc("GetEnvironmentStringsW")
	procFreeEnvironmentStringsW     = modkernel32.NewProc("FreeEnvironmentStringsW")
	procGetEnvironmentVariableW     = modkernel32.NewProc("GetEnvironmentVariableW")
	procSetEnvironmentVariableW     = modkernel32.NewProc("SetEnvironmentVariableW")
	procSetFileTime                 = modkernel32.NewProc("SetFileTime")
	procGetFileAttributesW          = modkernel32.NewProc("GetFileAttributesW")
	procSetFileAttributesW          = modkernel32.NewProc("SetFileAttributesW")
	procGetFileAttributesExW        = modkernel32.NewProc("GetFileAttributesExW")
	procGetCommandLineW             = modkernel32.NewProc("GetCommandLineW")
	procCommandLineToArgvW          = modshell32.NewProc("CommandLineToArgvW")
	procLocalFree                   = modkernel32.NewProc("LocalFree")
	procSetHandleInformation        = modkernel32.NewProc("SetHandleInformation")
	procFlushFileBuffers            = modkernel32.NewProc("FlushFileBuffers")
	procGetFullPathNameW            = modkernel32.NewProc("GetFullPathNameW")
	procCreateFileMappingW          = modkernel32.NewProc("CreateFileMappingW")
	procMapViewOfFile               = modkernel32.NewProc("MapViewOfFile")
	procUnmapViewOfFile             = modkernel32.NewProc("UnmapViewOfFile")
	procFlushViewOfFile             = modkernel32.NewProc("FlushViewOfFile")
	procVirtualLock                 = modkernel32.NewProc("VirtualLock")
	procVirtualUnlock               = modkernel32.NewProc("VirtualUnlock")
	procTransmitFile                = modmswsock.NewProc("TransmitFile")
	procReadDirectoryChangesW       = modkernel32.NewProc("ReadDirectoryChangesW")
	procCertOpenSystemStoreW        = modcrypt32.NewProc("CertOpenSystemStoreW")
	procCertEnumCertificatesInStore = modcrypt32.NewProc("CertEnumCertificatesInStore")
	procCertCloseStore              = modcrypt32.NewProc("CertCloseStore")
	procWSAStartup                  = modws2_32.NewProc("WSAStartup")
	procWSACleanup                  = modws2_32.NewProc("WSACleanup")
	procWSAIoctl                    = modws2_32.NewProc("WSAIoctl")
	procsocket                      = modws2_32.NewProc("socket")
	procsetsockopt                  = modws2_32.NewProc("setsockopt")
	procbind                        = modws2_32.NewProc("bind")
	procconnect                     = modws2_32.NewProc("connect")
	procgetsockname                 = modws2_32.NewProc("getsockname")
	procgetpeername                 = modws2_32.NewProc("getpeername")
	proclisten                      = modws2_32.NewProc("listen")
	procshutdown                    = modws2_32.NewProc("shutdown")
	procclosesocket                 = modws2_32.NewProc("closesocket")
	procAcceptEx                    = modmswsock.NewProc("AcceptEx")
	procGetAcceptExSockaddrs        = modmswsock.NewProc("GetAcceptExSockaddrs")
	procWSARecv                     = modws2_32.NewProc("WSARecv")
	procWSASend                     = modws2_32.NewProc("WSASend")
	procWSARecvFrom                 = modws2_32.NewProc("WSARecvFrom")
	procWSASendTo                   = modws2_32.NewProc("WSASendTo")
	procgethostbyname               = modws2_32.NewProc("gethostbyname")
	procgetservbyname               = modws2_32.NewProc("getservbyname")
	procntohs                       = modws2_32.NewProc("ntohs")
	procgetprotobyname              = modws2_32.NewProc("getprotobyname")
	procDnsQuery_W                  = moddnsapi.NewProc("DnsQuery_W")
	procDnsRecordListFree           = moddnsapi.NewProc("DnsRecordListFree")
	procGetIfEntry                  = modiphlpapi.NewProc("GetIfEntry")
	procGetAdaptersInfo             = modiphlpapi.NewProc("GetAdaptersInfo")
)

func GetLastError() (lasterrno int) {
	r0, _, _ := Syscall(procGetLastError.Addr(), 0, 0, 0, 0)
	lasterrno = int(r0)
	return
}

func LoadLibrary(libname string) (handle Handle, errno int) {
	r0, _, e1 := Syscall(procLoadLibraryW.Addr(), 1, uintptr(unsafe.Pointer(StringToUTF16Ptr(libname))), 0, 0)
	handle = Handle(r0)
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

func FreeLibrary(handle Handle) (errno int) {
	r1, _, e1 := Syscall(procFreeLibrary.Addr(), 1, uintptr(handle), 0, 0)
	if int(r1) == 0 {
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

func GetProcAddress(module Handle, procname string) (proc uintptr, errno int) {
	r0, _, e1 := Syscall(procGetProcAddress.Addr(), 2, uintptr(module), uintptr(unsafe.Pointer(StringBytePtr(procname))), 0)
	proc = uintptr(r0)
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
	r0, _, e1 := Syscall(procGetVersion.Addr(), 0, 0, 0, 0)
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
	r0, _, e1 := Syscall9(procFormatMessageW.Addr(), 7, uintptr(flags), uintptr(msgsrc), uintptr(msgid), uintptr(langid), uintptr(unsafe.Pointer(_p0)), uintptr(len(buf)), uintptr(unsafe.Pointer(args)), 0, 0)
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
	Syscall(procExitProcess.Addr(), 1, uintptr(exitcode), 0, 0)
	return
}

func CreateFile(name *uint16, access uint32, mode uint32, sa *SecurityAttributes, createmode uint32, attrs uint32, templatefile int32) (handle Handle, errno int) {
	r0, _, e1 := Syscall9(procCreateFileW.Addr(), 7, uintptr(unsafe.Pointer(name)), uintptr(access), uintptr(mode), uintptr(unsafe.Pointer(sa)), uintptr(createmode), uintptr(attrs), uintptr(templatefile), 0, 0)
	handle = Handle(r0)
	if handle == InvalidHandle {
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

func ReadFile(handle Handle, buf []byte, done *uint32, overlapped *Overlapped) (errno int) {
	var _p0 *byte
	if len(buf) > 0 {
		_p0 = &buf[0]
	}
	r1, _, e1 := Syscall6(procReadFile.Addr(), 5, uintptr(handle), uintptr(unsafe.Pointer(_p0)), uintptr(len(buf)), uintptr(unsafe.Pointer(done)), uintptr(unsafe.Pointer(overlapped)), 0)
	if int(r1) == 0 {
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

func WriteFile(handle Handle, buf []byte, done *uint32, overlapped *Overlapped) (errno int) {
	var _p0 *byte
	if len(buf) > 0 {
		_p0 = &buf[0]
	}
	r1, _, e1 := Syscall6(procWriteFile.Addr(), 5, uintptr(handle), uintptr(unsafe.Pointer(_p0)), uintptr(len(buf)), uintptr(unsafe.Pointer(done)), uintptr(unsafe.Pointer(overlapped)), 0)
	if int(r1) == 0 {
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

func SetFilePointer(handle Handle, lowoffset int32, highoffsetptr *int32, whence uint32) (newlowoffset uint32, errno int) {
	r0, _, e1 := Syscall6(procSetFilePointer.Addr(), 4, uintptr(handle), uintptr(lowoffset), uintptr(unsafe.Pointer(highoffsetptr)), uintptr(whence), 0, 0)
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

func CloseHandle(handle Handle) (errno int) {
	r1, _, e1 := Syscall(procCloseHandle.Addr(), 1, uintptr(handle), 0, 0)
	if int(r1) == 0 {
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

func GetStdHandle(stdhandle int) (handle Handle, errno int) {
	r0, _, e1 := Syscall(procGetStdHandle.Addr(), 1, uintptr(stdhandle), 0, 0)
	handle = Handle(r0)
	if handle == InvalidHandle {
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

func FindFirstFile(name *uint16, data *Win32finddata) (handle Handle, errno int) {
	r0, _, e1 := Syscall(procFindFirstFileW.Addr(), 2, uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(data)), 0)
	handle = Handle(r0)
	if handle == InvalidHandle {
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

func FindNextFile(handle Handle, data *Win32finddata) (errno int) {
	r1, _, e1 := Syscall(procFindNextFileW.Addr(), 2, uintptr(handle), uintptr(unsafe.Pointer(data)), 0)
	if int(r1) == 0 {
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

func FindClose(handle Handle) (errno int) {
	r1, _, e1 := Syscall(procFindClose.Addr(), 1, uintptr(handle), 0, 0)
	if int(r1) == 0 {
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

func GetFileInformationByHandle(handle Handle, data *ByHandleFileInformation) (errno int) {
	r1, _, e1 := Syscall(procGetFileInformationByHandle.Addr(), 2, uintptr(handle), uintptr(unsafe.Pointer(data)), 0)
	if int(r1) == 0 {
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
	r0, _, e1 := Syscall(procGetCurrentDirectoryW.Addr(), 2, uintptr(buflen), uintptr(unsafe.Pointer(buf)), 0)
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

func SetCurrentDirectory(path *uint16) (errno int) {
	r1, _, e1 := Syscall(procSetCurrentDirectoryW.Addr(), 1, uintptr(unsafe.Pointer(path)), 0, 0)
	if int(r1) == 0 {
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

func CreateDirectory(path *uint16, sa *SecurityAttributes) (errno int) {
	r1, _, e1 := Syscall(procCreateDirectoryW.Addr(), 2, uintptr(unsafe.Pointer(path)), uintptr(unsafe.Pointer(sa)), 0)
	if int(r1) == 0 {
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

func RemoveDirectory(path *uint16) (errno int) {
	r1, _, e1 := Syscall(procRemoveDirectoryW.Addr(), 1, uintptr(unsafe.Pointer(path)), 0, 0)
	if int(r1) == 0 {
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

func DeleteFile(path *uint16) (errno int) {
	r1, _, e1 := Syscall(procDeleteFileW.Addr(), 1, uintptr(unsafe.Pointer(path)), 0, 0)
	if int(r1) == 0 {
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

func MoveFile(from *uint16, to *uint16) (errno int) {
	r1, _, e1 := Syscall(procMoveFileW.Addr(), 2, uintptr(unsafe.Pointer(from)), uintptr(unsafe.Pointer(to)), 0)
	if int(r1) == 0 {
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

func GetComputerName(buf *uint16, n *uint32) (errno int) {
	r1, _, e1 := Syscall(procGetComputerNameW.Addr(), 2, uintptr(unsafe.Pointer(buf)), uintptr(unsafe.Pointer(n)), 0)
	if int(r1) == 0 {
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

func SetEndOfFile(handle Handle) (errno int) {
	r1, _, e1 := Syscall(procSetEndOfFile.Addr(), 1, uintptr(handle), 0, 0)
	if int(r1) == 0 {
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
	Syscall(procGetSystemTimeAsFileTime.Addr(), 1, uintptr(unsafe.Pointer(time)), 0, 0)
	return
}

func sleep(msec uint32) {
	Syscall(procSleep.Addr(), 1, uintptr(msec), 0, 0)
	return
}

func GetTimeZoneInformation(tzi *Timezoneinformation) (rc uint32, errno int) {
	r0, _, e1 := Syscall(procGetTimeZoneInformation.Addr(), 1, uintptr(unsafe.Pointer(tzi)), 0, 0)
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

func CreateIoCompletionPort(filehandle Handle, cphandle Handle, key uint32, threadcnt uint32) (handle Handle, errno int) {
	r0, _, e1 := Syscall6(procCreateIoCompletionPort.Addr(), 4, uintptr(filehandle), uintptr(cphandle), uintptr(key), uintptr(threadcnt), 0, 0)
	handle = Handle(r0)
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

func GetQueuedCompletionStatus(cphandle Handle, qty *uint32, key *uint32, overlapped **Overlapped, timeout uint32) (errno int) {
	r1, _, e1 := Syscall6(procGetQueuedCompletionStatus.Addr(), 5, uintptr(cphandle), uintptr(unsafe.Pointer(qty)), uintptr(unsafe.Pointer(key)), uintptr(unsafe.Pointer(overlapped)), uintptr(timeout), 0)
	if int(r1) == 0 {
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

func PostQueuedCompletionStatus(cphandle Handle, qty uint32, key uint32, overlapped *Overlapped) (errno int) {
	r1, _, e1 := Syscall6(procPostQueuedCompletionStatus.Addr(), 4, uintptr(cphandle), uintptr(qty), uintptr(key), uintptr(unsafe.Pointer(overlapped)), 0, 0)
	if int(r1) == 0 {
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

func CancelIo(s Handle) (errno int) {
	r1, _, e1 := Syscall(procCancelIo.Addr(), 1, uintptr(s), 0, 0)
	if int(r1) == 0 {
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

func CreateProcess(appName *uint16, commandLine *uint16, procSecurity *SecurityAttributes, threadSecurity *SecurityAttributes, inheritHandles bool, creationFlags uint32, env *uint16, currentDir *uint16, startupInfo *StartupInfo, outProcInfo *ProcessInformation) (errno int) {
	var _p0 uint32
	if inheritHandles {
		_p0 = 1
	} else {
		_p0 = 0
	}
	r1, _, e1 := Syscall12(procCreateProcessW.Addr(), 10, uintptr(unsafe.Pointer(appName)), uintptr(unsafe.Pointer(commandLine)), uintptr(unsafe.Pointer(procSecurity)), uintptr(unsafe.Pointer(threadSecurity)), uintptr(_p0), uintptr(creationFlags), uintptr(unsafe.Pointer(env)), uintptr(unsafe.Pointer(currentDir)), uintptr(unsafe.Pointer(startupInfo)), uintptr(unsafe.Pointer(outProcInfo)), 0, 0)
	if int(r1) == 0 {
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

func OpenProcess(da uint32, inheritHandle bool, pid uint32) (handle Handle, errno int) {
	var _p0 uint32
	if inheritHandle {
		_p0 = 1
	} else {
		_p0 = 0
	}
	r0, _, e1 := Syscall(procOpenProcess.Addr(), 3, uintptr(da), uintptr(_p0), uintptr(pid))
	handle = Handle(r0)
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

func TerminateProcess(handle Handle, exitcode uint32) (errno int) {
	r1, _, e1 := Syscall(procTerminateProcess.Addr(), 2, uintptr(handle), uintptr(exitcode), 0)
	if int(r1) == 0 {
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

func GetExitCodeProcess(handle Handle, exitcode *uint32) (errno int) {
	r1, _, e1 := Syscall(procGetExitCodeProcess.Addr(), 2, uintptr(handle), uintptr(unsafe.Pointer(exitcode)), 0)
	if int(r1) == 0 {
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

func GetStartupInfo(startupInfo *StartupInfo) (errno int) {
	r1, _, e1 := Syscall(procGetStartupInfoW.Addr(), 1, uintptr(unsafe.Pointer(startupInfo)), 0, 0)
	if int(r1) == 0 {
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

func GetCurrentProcess() (pseudoHandle Handle, errno int) {
	r0, _, e1 := Syscall(procGetCurrentProcess.Addr(), 0, 0, 0, 0)
	pseudoHandle = Handle(r0)
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

func DuplicateHandle(hSourceProcessHandle Handle, hSourceHandle Handle, hTargetProcessHandle Handle, lpTargetHandle *Handle, dwDesiredAccess uint32, bInheritHandle bool, dwOptions uint32) (errno int) {
	var _p0 uint32
	if bInheritHandle {
		_p0 = 1
	} else {
		_p0 = 0
	}
	r1, _, e1 := Syscall9(procDuplicateHandle.Addr(), 7, uintptr(hSourceProcessHandle), uintptr(hSourceHandle), uintptr(hTargetProcessHandle), uintptr(unsafe.Pointer(lpTargetHandle)), uintptr(dwDesiredAccess), uintptr(_p0), uintptr(dwOptions), 0, 0)
	if int(r1) == 0 {
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

func WaitForSingleObject(handle Handle, waitMilliseconds uint32) (event uint32, errno int) {
	r0, _, e1 := Syscall(procWaitForSingleObject.Addr(), 2, uintptr(handle), uintptr(waitMilliseconds), 0)
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
	r0, _, e1 := Syscall(procGetTempPathW.Addr(), 2, uintptr(buflen), uintptr(unsafe.Pointer(buf)), 0)
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

func CreatePipe(readhandle *Handle, writehandle *Handle, sa *SecurityAttributes, size uint32) (errno int) {
	r1, _, e1 := Syscall6(procCreatePipe.Addr(), 4, uintptr(unsafe.Pointer(readhandle)), uintptr(unsafe.Pointer(writehandle)), uintptr(unsafe.Pointer(sa)), uintptr(size), 0, 0)
	if int(r1) == 0 {
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

func GetFileType(filehandle Handle) (n uint32, errno int) {
	r0, _, e1 := Syscall(procGetFileType.Addr(), 1, uintptr(filehandle), 0, 0)
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

func CryptAcquireContext(provhandle *Handle, container *uint16, provider *uint16, provtype uint32, flags uint32) (errno int) {
	r1, _, e1 := Syscall6(procCryptAcquireContextW.Addr(), 5, uintptr(unsafe.Pointer(provhandle)), uintptr(unsafe.Pointer(container)), uintptr(unsafe.Pointer(provider)), uintptr(provtype), uintptr(flags), 0)
	if int(r1) == 0 {
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

func CryptReleaseContext(provhandle Handle, flags uint32) (errno int) {
	r1, _, e1 := Syscall(procCryptReleaseContext.Addr(), 2, uintptr(provhandle), uintptr(flags), 0)
	if int(r1) == 0 {
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

func CryptGenRandom(provhandle Handle, buflen uint32, buf *byte) (errno int) {
	r1, _, e1 := Syscall(procCryptGenRandom.Addr(), 3, uintptr(provhandle), uintptr(buflen), uintptr(unsafe.Pointer(buf)))
	if int(r1) == 0 {
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
	r0, _, e1 := Syscall(procGetEnvironmentStringsW.Addr(), 0, 0, 0, 0)
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

func FreeEnvironmentStrings(envs *uint16) (errno int) {
	r1, _, e1 := Syscall(procFreeEnvironmentStringsW.Addr(), 1, uintptr(unsafe.Pointer(envs)), 0, 0)
	if int(r1) == 0 {
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
	r0, _, e1 := Syscall(procGetEnvironmentVariableW.Addr(), 3, uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(buffer)), uintptr(size))
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

func SetEnvironmentVariable(name *uint16, value *uint16) (errno int) {
	r1, _, e1 := Syscall(procSetEnvironmentVariableW.Addr(), 2, uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(value)), 0)
	if int(r1) == 0 {
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

func SetFileTime(handle Handle, ctime *Filetime, atime *Filetime, wtime *Filetime) (errno int) {
	r1, _, e1 := Syscall6(procSetFileTime.Addr(), 4, uintptr(handle), uintptr(unsafe.Pointer(ctime)), uintptr(unsafe.Pointer(atime)), uintptr(unsafe.Pointer(wtime)), 0, 0)
	if int(r1) == 0 {
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
	r0, _, e1 := Syscall(procGetFileAttributesW.Addr(), 1, uintptr(unsafe.Pointer(name)), 0, 0)
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

func SetFileAttributes(name *uint16, attrs uint32) (errno int) {
	r1, _, e1 := Syscall(procSetFileAttributesW.Addr(), 2, uintptr(unsafe.Pointer(name)), uintptr(attrs), 0)
	if int(r1) == 0 {
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

func GetFileAttributesEx(name *uint16, level uint32, info *byte) (errno int) {
	r1, _, e1 := Syscall(procGetFileAttributesExW.Addr(), 3, uintptr(unsafe.Pointer(name)), uintptr(level), uintptr(unsafe.Pointer(info)))
	if int(r1) == 0 {
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

func GetCommandLine() (cmd *uint16) {
	r0, _, _ := Syscall(procGetCommandLineW.Addr(), 0, 0, 0, 0)
	cmd = (*uint16)(unsafe.Pointer(r0))
	return
}

func CommandLineToArgv(cmd *uint16, argc *int32) (argv *[8192]*[8192]uint16, errno int) {
	r0, _, e1 := Syscall(procCommandLineToArgvW.Addr(), 2, uintptr(unsafe.Pointer(cmd)), uintptr(unsafe.Pointer(argc)), 0)
	argv = (*[8192]*[8192]uint16)(unsafe.Pointer(r0))
	if argv == nil {
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

func LocalFree(hmem Handle) (handle Handle, errno int) {
	r0, _, e1 := Syscall(procLocalFree.Addr(), 1, uintptr(hmem), 0, 0)
	handle = Handle(r0)
	if handle != 0 {
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

func SetHandleInformation(handle Handle, mask uint32, flags uint32) (errno int) {
	r1, _, e1 := Syscall(procSetHandleInformation.Addr(), 3, uintptr(handle), uintptr(mask), uintptr(flags))
	if int(r1) == 0 {
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

func FlushFileBuffers(handle Handle) (errno int) {
	r1, _, e1 := Syscall(procFlushFileBuffers.Addr(), 1, uintptr(handle), 0, 0)
	if int(r1) == 0 {
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

func GetFullPathName(path *uint16, buflen uint32, buf *uint16, fname **uint16) (n uint32, errno int) {
	r0, _, e1 := Syscall6(procGetFullPathNameW.Addr(), 4, uintptr(unsafe.Pointer(path)), uintptr(buflen), uintptr(unsafe.Pointer(buf)), uintptr(unsafe.Pointer(fname)), 0, 0)
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

func CreateFileMapping(fhandle Handle, sa *SecurityAttributes, prot uint32, maxSizeHigh uint32, maxSizeLow uint32, name *uint16) (handle Handle, errno int) {
	r0, _, e1 := Syscall6(procCreateFileMappingW.Addr(), 6, uintptr(fhandle), uintptr(unsafe.Pointer(sa)), uintptr(prot), uintptr(maxSizeHigh), uintptr(maxSizeLow), uintptr(unsafe.Pointer(name)))
	handle = Handle(r0)
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

func MapViewOfFile(handle Handle, access uint32, offsetHigh uint32, offsetLow uint32, length uintptr) (addr uintptr, errno int) {
	r0, _, e1 := Syscall6(procMapViewOfFile.Addr(), 5, uintptr(handle), uintptr(access), uintptr(offsetHigh), uintptr(offsetLow), uintptr(length), 0)
	addr = uintptr(r0)
	if addr == 0 {
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

func UnmapViewOfFile(addr uintptr) (errno int) {
	r1, _, e1 := Syscall(procUnmapViewOfFile.Addr(), 1, uintptr(addr), 0, 0)
	if int(r1) == 0 {
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

func FlushViewOfFile(addr uintptr, length uintptr) (errno int) {
	r1, _, e1 := Syscall(procFlushViewOfFile.Addr(), 2, uintptr(addr), uintptr(length), 0)
	if int(r1) == 0 {
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

func VirtualLock(addr uintptr, length uintptr) (errno int) {
	r1, _, e1 := Syscall(procVirtualLock.Addr(), 2, uintptr(addr), uintptr(length), 0)
	if int(r1) == 0 {
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

func VirtualUnlock(addr uintptr, length uintptr) (errno int) {
	r1, _, e1 := Syscall(procVirtualUnlock.Addr(), 2, uintptr(addr), uintptr(length), 0)
	if int(r1) == 0 {
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

func TransmitFile(s Handle, handle Handle, bytesToWrite uint32, bytsPerSend uint32, overlapped *Overlapped, transmitFileBuf *TransmitFileBuffers, flags uint32) (errno int) {
	r1, _, e1 := Syscall9(procTransmitFile.Addr(), 7, uintptr(s), uintptr(handle), uintptr(bytesToWrite), uintptr(bytsPerSend), uintptr(unsafe.Pointer(overlapped)), uintptr(unsafe.Pointer(transmitFileBuf)), uintptr(flags), 0, 0)
	if int(r1) == 0 {
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

func ReadDirectoryChanges(handle Handle, buf *byte, buflen uint32, watchSubTree bool, mask uint32, retlen *uint32, overlapped *Overlapped, completionRoutine uintptr) (errno int) {
	var _p0 uint32
	if watchSubTree {
		_p0 = 1
	} else {
		_p0 = 0
	}
	r1, _, e1 := Syscall9(procReadDirectoryChangesW.Addr(), 8, uintptr(handle), uintptr(unsafe.Pointer(buf)), uintptr(buflen), uintptr(_p0), uintptr(mask), uintptr(unsafe.Pointer(retlen)), uintptr(unsafe.Pointer(overlapped)), uintptr(completionRoutine), 0)
	if int(r1) == 0 {
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

func CertOpenSystemStore(hprov Handle, name *uint16) (store Handle, errno int) {
	r0, _, e1 := Syscall(procCertOpenSystemStoreW.Addr(), 2, uintptr(hprov), uintptr(unsafe.Pointer(name)), 0)
	store = Handle(r0)
	if store == 0 {
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

func CertEnumCertificatesInStore(store Handle, prevContext *CertContext) (context *CertContext) {
	r0, _, _ := Syscall(procCertEnumCertificatesInStore.Addr(), 2, uintptr(store), uintptr(unsafe.Pointer(prevContext)), 0)
	context = (*CertContext)(unsafe.Pointer(r0))
	return
}

func CertCloseStore(store Handle, flags uint32) (errno int) {
	r1, _, e1 := Syscall(procCertCloseStore.Addr(), 2, uintptr(store), uintptr(flags), 0)
	if int(r1) == 0 {
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
	r0, _, _ := Syscall(procWSAStartup.Addr(), 2, uintptr(verreq), uintptr(unsafe.Pointer(data)), 0)
	sockerrno = int(r0)
	return
}

func WSACleanup() (errno int) {
	r1, _, e1 := Syscall(procWSACleanup.Addr(), 0, 0, 0, 0)
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

func WSAIoctl(s Handle, iocc uint32, inbuf *byte, cbif uint32, outbuf *byte, cbob uint32, cbbr *uint32, overlapped *Overlapped, completionRoutine uintptr) (errno int) {
	r1, _, e1 := Syscall9(procWSAIoctl.Addr(), 9, uintptr(s), uintptr(iocc), uintptr(unsafe.Pointer(inbuf)), uintptr(cbif), uintptr(unsafe.Pointer(outbuf)), uintptr(cbob), uintptr(unsafe.Pointer(cbbr)), uintptr(unsafe.Pointer(overlapped)), uintptr(completionRoutine))
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

func socket(af int32, typ int32, protocol int32) (handle Handle, errno int) {
	r0, _, e1 := Syscall(procsocket.Addr(), 3, uintptr(af), uintptr(typ), uintptr(protocol))
	handle = Handle(r0)
	if handle == InvalidHandle {
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

func Setsockopt(s Handle, level int32, optname int32, optval *byte, optlen int32) (errno int) {
	r1, _, e1 := Syscall6(procsetsockopt.Addr(), 5, uintptr(s), uintptr(level), uintptr(optname), uintptr(unsafe.Pointer(optval)), uintptr(optlen), 0)
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

func bind(s Handle, name uintptr, namelen int32) (errno int) {
	r1, _, e1 := Syscall(procbind.Addr(), 3, uintptr(s), uintptr(name), uintptr(namelen))
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

func connect(s Handle, name uintptr, namelen int32) (errno int) {
	r1, _, e1 := Syscall(procconnect.Addr(), 3, uintptr(s), uintptr(name), uintptr(namelen))
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

func getsockname(s Handle, rsa *RawSockaddrAny, addrlen *int32) (errno int) {
	r1, _, e1 := Syscall(procgetsockname.Addr(), 3, uintptr(s), uintptr(unsafe.Pointer(rsa)), uintptr(unsafe.Pointer(addrlen)))
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

func getpeername(s Handle, rsa *RawSockaddrAny, addrlen *int32) (errno int) {
	r1, _, e1 := Syscall(procgetpeername.Addr(), 3, uintptr(s), uintptr(unsafe.Pointer(rsa)), uintptr(unsafe.Pointer(addrlen)))
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

func listen(s Handle, backlog int32) (errno int) {
	r1, _, e1 := Syscall(proclisten.Addr(), 2, uintptr(s), uintptr(backlog), 0)
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

func shutdown(s Handle, how int32) (errno int) {
	r1, _, e1 := Syscall(procshutdown.Addr(), 2, uintptr(s), uintptr(how), 0)
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

func Closesocket(s Handle) (errno int) {
	r1, _, e1 := Syscall(procclosesocket.Addr(), 1, uintptr(s), 0, 0)
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

func AcceptEx(ls Handle, as Handle, buf *byte, rxdatalen uint32, laddrlen uint32, raddrlen uint32, recvd *uint32, overlapped *Overlapped) (errno int) {
	r1, _, e1 := Syscall9(procAcceptEx.Addr(), 8, uintptr(ls), uintptr(as), uintptr(unsafe.Pointer(buf)), uintptr(rxdatalen), uintptr(laddrlen), uintptr(raddrlen), uintptr(unsafe.Pointer(recvd)), uintptr(unsafe.Pointer(overlapped)), 0)
	if int(r1) == 0 {
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
	Syscall9(procGetAcceptExSockaddrs.Addr(), 8, uintptr(unsafe.Pointer(buf)), uintptr(rxdatalen), uintptr(laddrlen), uintptr(raddrlen), uintptr(unsafe.Pointer(lrsa)), uintptr(unsafe.Pointer(lrsalen)), uintptr(unsafe.Pointer(rrsa)), uintptr(unsafe.Pointer(rrsalen)), 0)
	return
}

func WSARecv(s Handle, bufs *WSABuf, bufcnt uint32, recvd *uint32, flags *uint32, overlapped *Overlapped, croutine *byte) (errno int) {
	r1, _, e1 := Syscall9(procWSARecv.Addr(), 7, uintptr(s), uintptr(unsafe.Pointer(bufs)), uintptr(bufcnt), uintptr(unsafe.Pointer(recvd)), uintptr(unsafe.Pointer(flags)), uintptr(unsafe.Pointer(overlapped)), uintptr(unsafe.Pointer(croutine)), 0, 0)
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

func WSASend(s Handle, bufs *WSABuf, bufcnt uint32, sent *uint32, flags uint32, overlapped *Overlapped, croutine *byte) (errno int) {
	r1, _, e1 := Syscall9(procWSASend.Addr(), 7, uintptr(s), uintptr(unsafe.Pointer(bufs)), uintptr(bufcnt), uintptr(unsafe.Pointer(sent)), uintptr(flags), uintptr(unsafe.Pointer(overlapped)), uintptr(unsafe.Pointer(croutine)), 0, 0)
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

func WSARecvFrom(s Handle, bufs *WSABuf, bufcnt uint32, recvd *uint32, flags *uint32, from *RawSockaddrAny, fromlen *int32, overlapped *Overlapped, croutine *byte) (errno int) {
	r1, _, e1 := Syscall9(procWSARecvFrom.Addr(), 9, uintptr(s), uintptr(unsafe.Pointer(bufs)), uintptr(bufcnt), uintptr(unsafe.Pointer(recvd)), uintptr(unsafe.Pointer(flags)), uintptr(unsafe.Pointer(from)), uintptr(unsafe.Pointer(fromlen)), uintptr(unsafe.Pointer(overlapped)), uintptr(unsafe.Pointer(croutine)))
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

func WSASendTo(s Handle, bufs *WSABuf, bufcnt uint32, sent *uint32, flags uint32, to *RawSockaddrAny, tolen int32, overlapped *Overlapped, croutine *byte) (errno int) {
	r1, _, e1 := Syscall9(procWSASendTo.Addr(), 9, uintptr(s), uintptr(unsafe.Pointer(bufs)), uintptr(bufcnt), uintptr(unsafe.Pointer(sent)), uintptr(flags), uintptr(unsafe.Pointer(to)), uintptr(tolen), uintptr(unsafe.Pointer(overlapped)), uintptr(unsafe.Pointer(croutine)))
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
	r0, _, e1 := Syscall(procgethostbyname.Addr(), 1, uintptr(unsafe.Pointer(StringBytePtr(name))), 0, 0)
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
	r0, _, e1 := Syscall(procgetservbyname.Addr(), 2, uintptr(unsafe.Pointer(StringBytePtr(name))), uintptr(unsafe.Pointer(StringBytePtr(proto))), 0)
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
	r0, _, _ := Syscall(procntohs.Addr(), 1, uintptr(netshort), 0, 0)
	u = uint16(r0)
	return
}

func GetProtoByName(name string) (p *Protoent, errno int) {
	r0, _, e1 := Syscall(procgetprotobyname.Addr(), 1, uintptr(unsafe.Pointer(StringBytePtr(name))), 0, 0)
	p = (*Protoent)(unsafe.Pointer(r0))
	if p == nil {
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

func DnsQuery(name string, qtype uint16, options uint32, extra *byte, qrs **DNSRecord, pr *byte) (status uint32) {
	r0, _, _ := Syscall6(procDnsQuery_W.Addr(), 6, uintptr(unsafe.Pointer(StringToUTF16Ptr(name))), uintptr(qtype), uintptr(options), uintptr(unsafe.Pointer(extra)), uintptr(unsafe.Pointer(qrs)), uintptr(unsafe.Pointer(pr)))
	status = uint32(r0)
	return
}

func DnsRecordListFree(rl *DNSRecord, freetype uint32) {
	Syscall(procDnsRecordListFree.Addr(), 2, uintptr(unsafe.Pointer(rl)), uintptr(freetype), 0)
	return
}

func GetIfEntry(pIfRow *MibIfRow) (errcode int) {
	r0, _, _ := Syscall(procGetIfEntry.Addr(), 1, uintptr(unsafe.Pointer(pIfRow)), 0, 0)
	errcode = int(r0)
	return
}

func GetAdaptersInfo(ai *IpAdapterInfo, ol *uint32) (errcode int) {
	r0, _, _ := Syscall(procGetAdaptersInfo.Addr(), 2, uintptr(unsafe.Pointer(ai)), uintptr(unsafe.Pointer(ol)), 0)
	errcode = int(r0)
	return
}
