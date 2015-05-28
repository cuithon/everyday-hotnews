// mksyscall.pl -l32 -plan9 syscall_plan9.go
// MACHINE GENERATED BY THE COMMAND ABOVE; DO NOT EDIT

// +build 386,plan9

package syscall

import "unsafe"

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func fd2path(fd int, buf []byte) (err error) {
	var _p0 unsafe.Pointer
	if len(buf) > 0 {
		_p0 = unsafe.Pointer(&buf[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall(SYS_FD2PATH, uintptr(fd), uintptr(_p0), uintptr(len(buf)))
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func pipe(p *[2]int32) (err error) {
	r0, _, e1 := Syscall(SYS_PIPE, uintptr(unsafe.Pointer(p)), 0, 0)
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func await(s []byte) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(s) > 0 {
		_p0 = unsafe.Pointer(&s[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall(SYS_AWAIT, uintptr(_p0), uintptr(len(s)), 0)
	n = int(r0)
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func open(path string, mode int) (fd int, err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	r0, _, e1 := Syscall(SYS_OPEN, uintptr(unsafe.Pointer(_p0)), uintptr(mode), 0)
	use(unsafe.Pointer(_p0))
	fd = int(r0)
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func create(path string, mode int, perm uint32) (fd int, err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	r0, _, e1 := Syscall(SYS_CREATE, uintptr(unsafe.Pointer(_p0)), uintptr(mode), uintptr(perm))
	use(unsafe.Pointer(_p0))
	fd = int(r0)
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func remove(path string) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	r0, _, e1 := Syscall(SYS_REMOVE, uintptr(unsafe.Pointer(_p0)), 0, 0)
	use(unsafe.Pointer(_p0))
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func stat(path string, edir []byte) (n int, err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	var _p1 unsafe.Pointer
	if len(edir) > 0 {
		_p1 = unsafe.Pointer(&edir[0])
	} else {
		_p1 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall(SYS_STAT, uintptr(unsafe.Pointer(_p0)), uintptr(_p1), uintptr(len(edir)))
	use(unsafe.Pointer(_p0))
	n = int(r0)
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func bind(name string, old string, flag int) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(name)
	if err != nil {
		return
	}
	var _p1 *byte
	_p1, err = BytePtrFromString(old)
	if err != nil {
		return
	}
	r0, _, e1 := Syscall(SYS_BIND, uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(_p1)), uintptr(flag))
	use(unsafe.Pointer(_p0))
	use(unsafe.Pointer(_p1))
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func mount(fd int, afd int, old string, flag int, aname string) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(old)
	if err != nil {
		return
	}
	var _p1 *byte
	_p1, err = BytePtrFromString(aname)
	if err != nil {
		return
	}
	r0, _, e1 := Syscall6(SYS_MOUNT, uintptr(fd), uintptr(afd), uintptr(unsafe.Pointer(_p0)), uintptr(flag), uintptr(unsafe.Pointer(_p1)), 0)
	use(unsafe.Pointer(_p0))
	use(unsafe.Pointer(_p1))
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func wstat(path string, edir []byte) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	var _p1 unsafe.Pointer
	if len(edir) > 0 {
		_p1 = unsafe.Pointer(&edir[0])
	} else {
		_p1 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall(SYS_WSTAT, uintptr(unsafe.Pointer(_p0)), uintptr(_p1), uintptr(len(edir)))
	use(unsafe.Pointer(_p0))
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func chdir(path string) (err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(path)
	if err != nil {
		return
	}
	r0, _, e1 := Syscall(SYS_CHDIR, uintptr(unsafe.Pointer(_p0)), 0, 0)
	use(unsafe.Pointer(_p0))
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Dup(oldfd int, newfd int) (fd int, err error) {
	r0, _, e1 := Syscall(SYS_DUP, uintptr(oldfd), uintptr(newfd), 0)
	fd = int(r0)
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Pread(fd int, p []byte, offset int64) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(p) > 0 {
		_p0 = unsafe.Pointer(&p[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall6(SYS_PREAD, uintptr(fd), uintptr(_p0), uintptr(len(p)), uintptr(offset), uintptr(offset>>32), 0)
	n = int(r0)
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Pwrite(fd int, p []byte, offset int64) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(p) > 0 {
		_p0 = unsafe.Pointer(&p[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall6(SYS_PWRITE, uintptr(fd), uintptr(_p0), uintptr(len(p)), uintptr(offset), uintptr(offset>>32), 0)
	n = int(r0)
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Close(fd int) (err error) {
	r0, _, e1 := Syscall(SYS_CLOSE, uintptr(fd), 0, 0)
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Fstat(fd int, edir []byte) (n int, err error) {
	var _p0 unsafe.Pointer
	if len(edir) > 0 {
		_p0 = unsafe.Pointer(&edir[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall(SYS_FSTAT, uintptr(fd), uintptr(_p0), uintptr(len(edir)))
	n = int(r0)
	if int32(r0) == -1 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func Fwstat(fd int, edir []byte) (err error) {
	var _p0 unsafe.Pointer
	if len(edir) > 0 {
		_p0 = unsafe.Pointer(&edir[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := Syscall(SYS_FWSTAT, uintptr(fd), uintptr(_p0), uintptr(len(edir)))
	if int32(r0) == -1 {
		err = e1
	}
	return
}
