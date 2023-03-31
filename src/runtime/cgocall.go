// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Cgo call and callback support.
//
// To call into the C function f from Go, the cgo-generated code calls
// runtime.cgocall(_cgo_Cfunc_f, frame), where _cgo_Cfunc_f is a
// gcc-compiled function written by cgo.
//
// runtime.cgocall (below) calls entersyscall so as not to block
// other goroutines or the garbage collector, and then calls
// runtime.asmcgocall(_cgo_Cfunc_f, frame).
//
// runtime.asmcgocall (in asm_$GOARCH.s) switches to the m->g0 stack
// (assumed to be an operating system-allocated stack, so safe to run
// gcc-compiled code on) and calls _cgo_Cfunc_f(frame).
//
// _cgo_Cfunc_f invokes the actual C function f with arguments
// taken from the frame structure, records the results in the frame,
// and returns to runtime.asmcgocall.
//
// After it regains control, runtime.asmcgocall switches back to the
// original g (m->curg)'s stack and returns to runtime.cgocall.
//
// After it regains control, runtime.cgocall calls exitsyscall, which blocks
// until this m can run Go code without violating the $GOMAXPROCS limit,
// and then unlocks g from m.
//
// The above description skipped over the possibility of the gcc-compiled
// function f calling back into Go. If that happens, we continue down
// the rabbit hole during the execution of f.
//
// To make it possible for gcc-compiled C code to call a Go function p.GoF,
// cgo writes a gcc-compiled function named GoF (not p.GoF, since gcc doesn't
// know about packages).  The gcc-compiled C function f calls GoF.
//
// GoF initializes "frame", a structure containing all of its
// arguments and slots for p.GoF's results. It calls
// crosscall2(_cgoexp_GoF, frame, framesize, ctxt) using the gcc ABI.
//
// crosscall2 (in cgo/asm_$GOARCH.s) is a four-argument adapter from
// the gcc function call ABI to the gc function call ABI. At this
// point we're in the Go runtime, but we're still running on m.g0's
// stack and outside the $GOMAXPROCS limit. crosscall2 calls
// runtime.cgocallback(_cgoexp_GoF, frame, ctxt) using the gc ABI.
// (crosscall2's framesize argument is no longer used, but there's one
// case where SWIG calls crosscall2 directly and expects to pass this
// argument. See _cgo_panic.)
//
// runtime.cgocallback (in asm_$GOARCH.s) switches from m.g0's stack
// to the original g (m.curg)'s stack, on which it calls
// runtime.cgocallbackg(_cgoexp_GoF, frame, ctxt). As part of the
// stack switch, runtime.cgocallback saves the current SP as
// m.g0.sched.sp, so that any use of m.g0's stack during the execution
// of the callback will be done below the existing stack frames.
// Before overwriting m.g0.sched.sp, it pushes the old value on the
// m.g0 stack, so that it can be restored later.
//
// runtime.cgocallbackg (below) is now running on a real goroutine
// stack (not an m.g0 stack).  First it calls runtime.exitsyscall, which will
// block until the $GOMAXPROCS limit allows running this goroutine.
// Once exitsyscall has returned, it is safe to do things like call the memory
// allocator or invoke the Go callback function.  runtime.cgocallbackg
// first defers a function to unwind m.g0.sched.sp, so that if p.GoF
// panics, m.g0.sched.sp will be restored to its old value: the m.g0 stack
// and the m.curg stack will be unwound in lock step.
// Then it calls _cgoexp_GoF(frame).
//
// _cgoexp_GoF, which was generated by cmd/cgo, unpacks the arguments
// from frame, calls p.GoF, writes the results back to frame, and
// returns. Now we start unwinding this whole process.
//
// runtime.cgocallbackg pops but does not execute the deferred
// function to unwind m.g0.sched.sp, calls runtime.entersyscall, and
// returns to runtime.cgocallback.
//
// After it regains control, runtime.cgocallback switches back to
// m.g0's stack (the pointer is still in m.g0.sched.sp), restores the old
// m.g0.sched.sp value from the stack, and returns to crosscall2.
//
// crosscall2 restores the callee-save registers for gcc and returns
// to GoF, which unpacks any result values and returns to f.

package runtime

import (
	"internal/goarch"
	"internal/goexperiment"
	"runtime/internal/sys"
	"unsafe"
)

// Addresses collected in a cgo backtrace when crashing.
// Length must match arg.Max in x_cgo_callers in runtime/cgo/gcc_traceback.c.
type cgoCallers [32]uintptr

// argset matches runtime/cgo/linux_syscall.c:argset_t
type argset struct {
	args   unsafe.Pointer
	retval uintptr
}

// wrapper for syscall package to call cgocall for libc (cgo) calls.
//
//go:linkname syscall_cgocaller syscall.cgocaller
//go:nosplit
//go:uintptrescapes
func syscall_cgocaller(fn unsafe.Pointer, args ...uintptr) uintptr {
	as := argset{args: unsafe.Pointer(&args[0])}
	cgocall(fn, unsafe.Pointer(&as))
	return as.retval
}

var ncgocall uint64 // number of cgo calls in total for dead m

// Call from Go to C.
//
// This must be nosplit because it's used for syscalls on some
// platforms. Syscalls may have untyped arguments on the stack, so
// it's not safe to grow or scan the stack.
//
//go:nosplit
func cgocall(fn, arg unsafe.Pointer) int32 {
	if !iscgo && GOOS != "solaris" && GOOS != "illumos" && GOOS != "windows" {
		throw("cgocall unavailable")
	}

	if fn == nil {
		throw("cgocall nil")
	}

	if raceenabled {
		racereleasemerge(unsafe.Pointer(&racecgosync))
	}

	mp := getg().m
	mp.ncgocall++
	mp.ncgo++

	// Reset traceback.
	mp.cgoCallers[0] = 0

	// Announce we are entering a system call
	// so that the scheduler knows to create another
	// M to run goroutines while we are in the
	// foreign code.
	//
	// The call to asmcgocall is guaranteed not to
	// grow the stack and does not allocate memory,
	// so it is safe to call while "in a system call", outside
	// the $GOMAXPROCS accounting.
	//
	// fn may call back into Go code, in which case we'll exit the
	// "system call", run the Go code (which may grow the stack),
	// and then re-enter the "system call" reusing the PC and SP
	// saved by entersyscall here.
	entersyscall()

	// Tell asynchronous preemption that we're entering external
	// code. We do this after entersyscall because this may block
	// and cause an async preemption to fail, but at this point a
	// sync preemption will succeed (though this is not a matter
	// of correctness).
	osPreemptExtEnter(mp)

	mp.incgo = true
	errno := asmcgocall(fn, arg)

	// Update accounting before exitsyscall because exitsyscall may
	// reschedule us on to a different M.
	mp.incgo = false
	mp.ncgo--

	osPreemptExtExit(mp)

	exitsyscall()

	// Note that raceacquire must be called only after exitsyscall has
	// wired this M to a P.
	if raceenabled {
		raceacquire(unsafe.Pointer(&racecgosync))
	}

	// From the garbage collector's perspective, time can move
	// backwards in the sequence above. If there's a callback into
	// Go code, GC will see this function at the call to
	// asmcgocall. When the Go call later returns to C, the
	// syscall PC/SP is rolled back and the GC sees this function
	// back at the call to entersyscall. Normally, fn and arg
	// would be live at entersyscall and dead at asmcgocall, so if
	// time moved backwards, GC would see these arguments as dead
	// and then live. Prevent these undead arguments from crashing
	// GC by forcing them to stay live across this time warp.
	KeepAlive(fn)
	KeepAlive(arg)
	KeepAlive(mp)

	return errno
}

// Call from C back to Go. fn must point to an ABIInternal Go entry-point.
//
//go:nosplit
func cgocallbackg(fn, frame unsafe.Pointer, ctxt uintptr) {
	gp := getg()
	if gp != gp.m.curg {
		println("runtime: bad g in cgocallback")
		exit(2)
	}

	// The call from C is on gp.m's g0 stack, so we must ensure
	// that we stay on that M. We have to do this before calling
	// exitsyscall, since it would otherwise be free to move us to
	// a different M. The call to unlockOSThread is in unwindm.
	lockOSThread()

	checkm := gp.m

	// Save current syscall parameters, so m.syscall can be
	// used again if callback decide to make syscall.
	syscall := gp.m.syscall

	// entersyscall saves the caller's SP to allow the GC to trace the Go
	// stack. However, since we're returning to an earlier stack frame and
	// need to pair with the entersyscall() call made by cgocall, we must
	// save syscall* and let reentersyscall restore them.
	savedsp := unsafe.Pointer(gp.syscallsp)
	savedpc := gp.syscallpc
	exitsyscall() // coming out of cgo call
	gp.m.incgo = false
	if gp.m.isextra {
		gp.m.isExtraInC = false
	}

	osPreemptExtExit(gp.m)

	cgocallbackg1(fn, frame, ctxt) // will call unlockOSThread

	// At this point unlockOSThread has been called.
	// The following code must not change to a different m.
	// This is enforced by checking incgo in the schedule function.

	gp.m.incgo = true
	if gp.m.isextra {
		gp.m.isExtraInC = true
	}

	if gp.m != checkm {
		throw("m changed unexpectedly in cgocallbackg")
	}

	osPreemptExtEnter(gp.m)

	// going back to cgo call
	reentersyscall(savedpc, uintptr(savedsp))

	gp.m.syscall = syscall
}

func cgocallbackg1(fn, frame unsafe.Pointer, ctxt uintptr) {
	gp := getg()

	// When we return, undo the call to lockOSThread in cgocallbackg.
	// We must still stay on the same m.
	defer unlockOSThread()

	if gp.m.needextram || extraMWaiters.Load() > 0 {
		gp.m.needextram = false
		systemstack(newextram)
	}

	if ctxt != 0 {
		s := append(gp.cgoCtxt, ctxt)

		// Now we need to set gp.cgoCtxt = s, but we could get
		// a SIGPROF signal while manipulating the slice, and
		// the SIGPROF handler could pick up gp.cgoCtxt while
		// tracing up the stack.  We need to ensure that the
		// handler always sees a valid slice, so set the
		// values in an order such that it always does.
		p := (*slice)(unsafe.Pointer(&gp.cgoCtxt))
		atomicstorep(unsafe.Pointer(&p.array), unsafe.Pointer(&s[0]))
		p.cap = cap(s)
		p.len = len(s)

		defer func(gp *g) {
			// Decrease the length of the slice by one, safely.
			p := (*slice)(unsafe.Pointer(&gp.cgoCtxt))
			p.len--
		}(gp)
	}

	if gp.m.ncgo == 0 {
		// The C call to Go came from a thread not currently running
		// any Go. In the case of -buildmode=c-archive or c-shared,
		// this call may be coming in before package initialization
		// is complete. Wait until it is.
		<-main_init_done
	}

	// Check whether the profiler needs to be turned on or off; this route to
	// run Go code does not use runtime.execute, so bypasses the check there.
	hz := sched.profilehz
	if gp.m.profilehz != hz {
		setThreadCPUProfiler(hz)
	}

	// Add entry to defer stack in case of panic.
	restore := true
	defer unwindm(&restore)

	if raceenabled {
		raceacquire(unsafe.Pointer(&racecgosync))
	}

	// Invoke callback. This function is generated by cmd/cgo and
	// will unpack the argument frame and call the Go function.
	var cb func(frame unsafe.Pointer)
	cbFV := funcval{uintptr(fn)}
	*(*unsafe.Pointer)(unsafe.Pointer(&cb)) = noescape(unsafe.Pointer(&cbFV))
	cb(frame)

	if raceenabled {
		racereleasemerge(unsafe.Pointer(&racecgosync))
	}

	// Do not unwind m->g0->sched.sp.
	// Our caller, cgocallback, will do that.
	restore = false
}

func unwindm(restore *bool) {
	if *restore {
		// Restore sp saved by cgocallback during
		// unwind of g's stack (see comment at top of file).
		mp := acquirem()
		sched := &mp.g0.sched
		sched.sp = *(*uintptr)(unsafe.Pointer(sched.sp + alignUp(sys.MinFrameSize, sys.StackAlign)))

		// Do the accounting that cgocall will not have a chance to do
		// during an unwind.
		//
		// In the case where a Go call originates from C, ncgo is 0
		// and there is no matching cgocall to end.
		if mp.ncgo > 0 {
			mp.incgo = false
			mp.ncgo--
			osPreemptExtExit(mp)
		}

		releasem(mp)
	}
}

// called from assembly.
func badcgocallback() {
	throw("misaligned stack in cgocallback")
}

// called from (incomplete) assembly.
func cgounimpl() {
	throw("cgo not implemented")
}

var racecgosync uint64 // represents possible synchronization in C code

// Pointer checking for cgo code.

// We want to detect all cases where a program that does not use
// unsafe makes a cgo call passing a Go pointer to memory that
// contains a Go pointer. Here a Go pointer is defined as a pointer
// to memory allocated by the Go runtime. Programs that use unsafe
// can evade this restriction easily, so we don't try to catch them.
// The cgo program will rewrite all possibly bad pointer arguments to
// call cgoCheckPointer, where we can catch cases of a Go pointer
// pointing to a Go pointer.

// Complicating matters, taking the address of a slice or array
// element permits the C program to access all elements of the slice
// or array. In that case we will see a pointer to a single element,
// but we need to check the entire data structure.

// The cgoCheckPointer call takes additional arguments indicating that
// it was called on an address expression. An additional argument of
// true means that it only needs to check a single element. An
// additional argument of a slice or array means that it needs to
// check the entire slice/array, but nothing else. Otherwise, the
// pointer could be anything, and we check the entire heap object,
// which is conservative but safe.

// When and if we implement a moving garbage collector,
// cgoCheckPointer will pin the pointer for the duration of the cgo
// call.  (This is necessary but not sufficient; the cgo program will
// also have to change to pin Go pointers that cannot point to Go
// pointers.)

// cgoCheckPointer checks if the argument contains a Go pointer that
// points to a Go pointer, and panics if it does.
func cgoCheckPointer(ptr any, arg any) {
	if !goexperiment.CgoCheck2 && debug.cgocheck == 0 {
		return
	}

	ep := efaceOf(&ptr)
	t := ep._type

	top := true
	if arg != nil && (t.kind&kindMask == kindPtr || t.kind&kindMask == kindUnsafePointer) {
		p := ep.data
		if t.kind&kindDirectIface == 0 {
			p = *(*unsafe.Pointer)(p)
		}
		if p == nil || !cgoIsGoPointer(p) {
			return
		}
		aep := efaceOf(&arg)
		switch aep._type.kind & kindMask {
		case kindBool:
			if t.kind&kindMask == kindUnsafePointer {
				// We don't know the type of the element.
				break
			}
			pt := (*ptrtype)(unsafe.Pointer(t))
			cgoCheckArg(pt.elem, p, true, false, cgoCheckPointerFail)
			return
		case kindSlice:
			// Check the slice rather than the pointer.
			ep = aep
			t = ep._type
		case kindArray:
			// Check the array rather than the pointer.
			// Pass top as false since we have a pointer
			// to the array.
			ep = aep
			t = ep._type
			top = false
		default:
			throw("can't happen")
		}
	}

	cgoCheckArg(t, ep.data, t.kind&kindDirectIface == 0, top, cgoCheckPointerFail)
}

const cgoCheckPointerFail = "cgo argument has Go pointer to Go pointer"
const cgoResultFail = "cgo result has Go pointer"

// cgoCheckArg is the real work of cgoCheckPointer. The argument p
// is either a pointer to the value (of type t), or the value itself,
// depending on indir. The top parameter is whether we are at the top
// level, where Go pointers are allowed.
func cgoCheckArg(t *_type, p unsafe.Pointer, indir, top bool, msg string) {
	if t.ptrdata == 0 || p == nil {
		// If the type has no pointers there is nothing to do.
		return
	}

	switch t.kind & kindMask {
	default:
		throw("can't happen")
	case kindArray:
		at := (*arraytype)(unsafe.Pointer(t))
		if !indir {
			if at.len != 1 {
				throw("can't happen")
			}
			cgoCheckArg(at.elem, p, at.elem.kind&kindDirectIface == 0, top, msg)
			return
		}
		for i := uintptr(0); i < at.len; i++ {
			cgoCheckArg(at.elem, p, true, top, msg)
			p = add(p, at.elem.size)
		}
	case kindChan, kindMap:
		// These types contain internal pointers that will
		// always be allocated in the Go heap. It's never OK
		// to pass them to C.
		panic(errorString(msg))
	case kindFunc:
		if indir {
			p = *(*unsafe.Pointer)(p)
		}
		if !cgoIsGoPointer(p) {
			return
		}
		panic(errorString(msg))
	case kindInterface:
		it := *(**_type)(p)
		if it == nil {
			return
		}
		// A type known at compile time is OK since it's
		// constant. A type not known at compile time will be
		// in the heap and will not be OK.
		if inheap(uintptr(unsafe.Pointer(it))) {
			panic(errorString(msg))
		}
		p = *(*unsafe.Pointer)(add(p, goarch.PtrSize))
		if !cgoIsGoPointer(p) {
			return
		}
		if !top {
			panic(errorString(msg))
		}
		cgoCheckArg(it, p, it.kind&kindDirectIface == 0, false, msg)
	case kindSlice:
		st := (*slicetype)(unsafe.Pointer(t))
		s := (*slice)(p)
		p = s.array
		if p == nil || !cgoIsGoPointer(p) {
			return
		}
		if !top {
			panic(errorString(msg))
		}
		if st.elem.ptrdata == 0 {
			return
		}
		for i := 0; i < s.cap; i++ {
			cgoCheckArg(st.elem, p, true, false, msg)
			p = add(p, st.elem.size)
		}
	case kindString:
		ss := (*stringStruct)(p)
		if !cgoIsGoPointer(ss.str) {
			return
		}
		if !top {
			panic(errorString(msg))
		}
	case kindStruct:
		st := (*structtype)(unsafe.Pointer(t))
		if !indir {
			if len(st.fields) != 1 {
				throw("can't happen")
			}
			cgoCheckArg(st.fields[0].typ, p, st.fields[0].typ.kind&kindDirectIface == 0, top, msg)
			return
		}
		for _, f := range st.fields {
			if f.typ.ptrdata == 0 {
				continue
			}
			cgoCheckArg(f.typ, add(p, f.offset), true, top, msg)
		}
	case kindPtr, kindUnsafePointer:
		if indir {
			p = *(*unsafe.Pointer)(p)
			if p == nil {
				return
			}
		}

		if !cgoIsGoPointer(p) {
			return
		}
		if !top {
			panic(errorString(msg))
		}

		cgoCheckUnknownPointer(p, msg)
	}
}

// cgoCheckUnknownPointer is called for an arbitrary pointer into Go
// memory. It checks whether that Go memory contains any other
// pointer into Go memory. If it does, we panic.
// The return values are unused but useful to see in panic tracebacks.
func cgoCheckUnknownPointer(p unsafe.Pointer, msg string) (base, i uintptr) {
	if inheap(uintptr(p)) {
		b, span, _ := findObject(uintptr(p), 0, 0)
		base = b
		if base == 0 {
			return
		}
		n := span.elemsize
		hbits := heapBitsForAddr(base, n)
		for {
			var addr uintptr
			if hbits, addr = hbits.next(); addr == 0 {
				break
			}
			if cgoIsGoPointer(*(*unsafe.Pointer)(unsafe.Pointer(addr))) {
				panic(errorString(msg))
			}
		}

		return
	}

	for _, datap := range activeModules() {
		if cgoInRange(p, datap.data, datap.edata) || cgoInRange(p, datap.bss, datap.ebss) {
			// We have no way to know the size of the object.
			// We have to assume that it might contain a pointer.
			panic(errorString(msg))
		}
		// In the text or noptr sections, we know that the
		// pointer does not point to a Go pointer.
	}

	return
}

// cgoIsGoPointer reports whether the pointer is a Go pointer--a
// pointer to Go memory. We only care about Go memory that might
// contain pointers.
//
//go:nosplit
//go:nowritebarrierrec
func cgoIsGoPointer(p unsafe.Pointer) bool {
	if p == nil {
		return false
	}

	if inHeapOrStack(uintptr(p)) {
		return true
	}

	for _, datap := range activeModules() {
		if cgoInRange(p, datap.data, datap.edata) || cgoInRange(p, datap.bss, datap.ebss) {
			return true
		}
	}

	return false
}

// cgoInRange reports whether p is between start and end.
//
//go:nosplit
//go:nowritebarrierrec
func cgoInRange(p unsafe.Pointer, start, end uintptr) bool {
	return start <= uintptr(p) && uintptr(p) < end
}

// cgoCheckResult is called to check the result parameter of an
// exported Go function. It panics if the result is or contains a Go
// pointer.
func cgoCheckResult(val any) {
	if !goexperiment.CgoCheck2 && debug.cgocheck == 0 {
		return
	}

	ep := efaceOf(&val)
	t := ep._type
	cgoCheckArg(t, ep.data, t.kind&kindDirectIface == 0, false, cgoResultFail)
}
