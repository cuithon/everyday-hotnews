// mksyscall_windows.sh winapi.go
// MACHINE GENERATED BY THE COMMAND ABOVE; DO NOT EDIT

package main

import "unsafe"
import "syscall"

var (
	modkernel32 = loadDll("kernel32.dll")
	moduser32   = loadDll("user32.dll")

	procGetModuleHandleW = getSysProcAddr(modkernel32, "GetModuleHandleW")
	procRegisterClassExW = getSysProcAddr(moduser32, "RegisterClassExW")
	procCreateWindowExW  = getSysProcAddr(moduser32, "CreateWindowExW")
	procDefWindowProcW   = getSysProcAddr(moduser32, "DefWindowProcW")
	procDestroyWindow    = getSysProcAddr(moduser32, "DestroyWindow")
	procPostQuitMessage  = getSysProcAddr(moduser32, "PostQuitMessage")
	procShowWindow       = getSysProcAddr(moduser32, "ShowWindow")
	procUpdateWindow     = getSysProcAddr(moduser32, "UpdateWindow")
	procGetMessageW      = getSysProcAddr(moduser32, "GetMessageW")
	procTranslateMessage = getSysProcAddr(moduser32, "TranslateMessage")
	procDispatchMessageW = getSysProcAddr(moduser32, "DispatchMessageW")
	procLoadIconW        = getSysProcAddr(moduser32, "LoadIconW")
	procLoadCursorW      = getSysProcAddr(moduser32, "LoadCursorW")
	procSetCursor        = getSysProcAddr(moduser32, "SetCursor")
	procSendMessageW     = getSysProcAddr(moduser32, "SendMessageW")
	procPostMessageW     = getSysProcAddr(moduser32, "PostMessageW")
)

func GetModuleHandle(modname *uint16) (handle uint32, errno int) {
	r0, _, e1 := syscall.Syscall(procGetModuleHandleW, uintptr(unsafe.Pointer(modname)), 0, 0)
	handle = uint32(r0)
	if handle == 0 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = syscall.EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func RegisterClassEx(wndclass *Wndclassex) (atom uint16, errno int) {
	r0, _, e1 := syscall.Syscall(procRegisterClassExW, uintptr(unsafe.Pointer(wndclass)), 0, 0)
	atom = uint16(r0)
	if atom == 0 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = syscall.EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func CreateWindowEx(exstyle uint32, classname *uint16, windowname *uint16, style uint32, x int32, y int32, width int32, height int32, wndparent uint32, menu uint32, instance uint32, param uintptr) (hwnd uint32, errno int) {
	r0, _, e1 := syscall.Syscall12(procCreateWindowExW, uintptr(exstyle), uintptr(unsafe.Pointer(classname)), uintptr(unsafe.Pointer(windowname)), uintptr(style), uintptr(x), uintptr(y), uintptr(width), uintptr(height), uintptr(wndparent), uintptr(menu), uintptr(instance), uintptr(param))
	hwnd = uint32(r0)
	if hwnd == 0 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = syscall.EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func DefWindowProc(hwnd uint32, msg uint32, wparam int32, lparam int32) (lresult int32) {
	r0, _, _ := syscall.Syscall6(procDefWindowProcW, uintptr(hwnd), uintptr(msg), uintptr(wparam), uintptr(lparam), 0, 0)
	lresult = int32(r0)
	return
}

func DestroyWindow(hwnd uint32) (ok bool, errno int) {
	r0, _, e1 := syscall.Syscall(procDestroyWindow, uintptr(hwnd), 0, 0)
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = syscall.EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func PostQuitMessage(exitcode int32) {
	syscall.Syscall(procPostQuitMessage, uintptr(exitcode), 0, 0)
	return
}

func ShowWindow(hwnd uint32, cmdshow int32) (ok bool) {
	r0, _, _ := syscall.Syscall(procShowWindow, uintptr(hwnd), uintptr(cmdshow), 0)
	ok = bool(r0 != 0)
	return
}

func UpdateWindow(hwnd uint32) (ok bool, errno int) {
	r0, _, e1 := syscall.Syscall(procUpdateWindow, uintptr(hwnd), 0, 0)
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = syscall.EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func GetMessage(msg *Msg, hwnd uint32, MsgFilterMin uint32, MsgFilterMax uint32) (ret int32, errno int) {
	r0, _, e1 := syscall.Syscall6(procGetMessageW, uintptr(unsafe.Pointer(msg)), uintptr(hwnd), uintptr(MsgFilterMin), uintptr(MsgFilterMax), 0, 0)
	ret = int32(r0)
	if ret == -1 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = syscall.EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func TranslateMessage(msg *Msg) (ok bool) {
	r0, _, _ := syscall.Syscall(procTranslateMessage, uintptr(unsafe.Pointer(msg)), 0, 0)
	ok = bool(r0 != 0)
	return
}

func DispatchMessage(msg *Msg) (ret int32) {
	r0, _, _ := syscall.Syscall(procDispatchMessageW, uintptr(unsafe.Pointer(msg)), 0, 0)
	ret = int32(r0)
	return
}

func LoadIcon(instance uint32, iconname *uint16) (icon uint32, errno int) {
	r0, _, e1 := syscall.Syscall(procLoadIconW, uintptr(instance), uintptr(unsafe.Pointer(iconname)), 0)
	icon = uint32(r0)
	if icon == 0 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = syscall.EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func LoadCursor(instance uint32, cursorname *uint16) (cursor uint32, errno int) {
	r0, _, e1 := syscall.Syscall(procLoadCursorW, uintptr(instance), uintptr(unsafe.Pointer(cursorname)), 0)
	cursor = uint32(r0)
	if cursor == 0 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = syscall.EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func SetCursor(cursor uint32) (precursor uint32, errno int) {
	r0, _, e1 := syscall.Syscall(procSetCursor, uintptr(cursor), 0, 0)
	precursor = uint32(r0)
	if precursor == 0 {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = syscall.EINVAL
		}
	} else {
		errno = 0
	}
	return
}

func SendMessage(hwnd uint32, msg uint32, wparam int32, lparam int32) (lresult int32) {
	r0, _, _ := syscall.Syscall6(procSendMessageW, uintptr(hwnd), uintptr(msg), uintptr(wparam), uintptr(lparam), 0, 0)
	lresult = int32(r0)
	return
}

func PostMessage(hwnd uint32, msg uint32, wparam int32, lparam int32) (ok bool, errno int) {
	r0, _, e1 := syscall.Syscall6(procPostMessageW, uintptr(hwnd), uintptr(msg), uintptr(wparam), uintptr(lparam), 0, 0)
	ok = bool(r0 != 0)
	if !ok {
		if e1 != 0 {
			errno = int(e1)
		} else {
			errno = syscall.EINVAL
		}
	} else {
		errno = 0
	}
	return
}
