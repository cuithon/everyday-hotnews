// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"internal/abi"
	"internal/goarch"
	"runtime/internal/atomic"
	"runtime/internal/sys"
	"unsafe"
)

// throwType indicates the current type of ongoing throw, which affects the
// amount of detail printed to stderr. Higher values include more detail.
type throwType uint32

const (
	// throwTypeNone means that we are not throwing.
	throwTypeNone throwType = iota

	// throwTypeUser is a throw due to a problem with the application.
	//
	// These throws do not include runtime frames, system goroutines, or
	// frame metadata.
	throwTypeUser

	// throwTypeRuntime is a throw due to a problem with Go itself.
	//
	// These throws include as much information as possible to aid in
	// debugging the runtime, including runtime frames, system goroutines,
	// and frame metadata.
	throwTypeRuntime
)

// We have two different ways of doing defers. The older way involves creating a
// defer record at the time that a defer statement is executing and adding it to a
// defer chain. This chain is inspected by the deferreturn call at all function
// exits in order to run the appropriate defer calls. A cheaper way (which we call
// open-coded defers) is used for functions in which no defer statements occur in
// loops. In that case, we simply store the defer function/arg information into
// specific stack slots at the point of each defer statement, as well as setting a
// bit in a bitmask. At each function exit, we add inline code to directly make
// the appropriate defer calls based on the bitmask and fn/arg information stored
// on the stack. During panic/Goexit processing, the appropriate defer calls are
// made using extra funcdata info that indicates the exact stack slots that
// contain the bitmask and defer fn/args.

// Check to make sure we can really generate a panic. If the panic
// was generated from the runtime, or from inside malloc, then convert
// to a throw of msg.
// pc should be the program counter of the compiler-generated code that
// triggered this panic.
func panicCheck1(pc uintptr, msg string) {
	if goarch.IsWasm == 0 && hasPrefix(funcname(findfunc(pc)), "runtime.") {
		// Note: wasm can't tail call, so we can't get the original caller's pc.
		throw(msg)
	}
	// TODO: is this redundant? How could we be in malloc
	// but not in the runtime? runtime/internal/*, maybe?
	gp := getg()
	if gp != nil && gp.m != nil && gp.m.mallocing != 0 {
		throw(msg)
	}
}

// Same as above, but calling from the runtime is allowed.
//
// Using this function is necessary for any panic that may be
// generated by runtime.sigpanic, since those are always called by the
// runtime.
func panicCheck2(err string) {
	// panic allocates, so to avoid recursive malloc, turn panics
	// during malloc into throws.
	gp := getg()
	if gp != nil && gp.m != nil && gp.m.mallocing != 0 {
		throw(err)
	}
}

// Many of the following panic entry-points turn into throws when they
// happen in various runtime contexts. These should never happen in
// the runtime, and if they do, they indicate a serious issue and
// should not be caught by user code.
//
// The panic{Index,Slice,divide,shift} functions are called by
// code generated by the compiler for out of bounds index expressions,
// out of bounds slice expressions, division by zero, and shift by negative.
// The panicdivide (again), panicoverflow, panicfloat, and panicmem
// functions are called by the signal handler when a signal occurs
// indicating the respective problem.
//
// Since panic{Index,Slice,shift} are never called directly, and
// since the runtime package should never have an out of bounds slice
// or array reference or negative shift, if we see those functions called from the
// runtime package we turn the panic into a throw. That will dump the
// entire runtime stack for easier debugging.
//
// The entry points called by the signal handler will be called from
// runtime.sigpanic, so we can't disallow calls from the runtime to
// these (they always look like they're called from the runtime).
// Hence, for these, we just check for clearly bad runtime conditions.
//
// The panic{Index,Slice} functions are implemented in assembly and tail call
// to the goPanic{Index,Slice} functions below. This is done so we can use
// a space-minimal register calling convention.

// failures in the comparisons for s[x], 0 <= x < y (y == len(s))
//
//go:yeswritebarrierrec
func goPanicIndex(x int, y int) {
	panicCheck1(getcallerpc(), "index out of range")
	panic(boundsError{x: int64(x), signed: true, y: y, code: boundsIndex})
}

//go:yeswritebarrierrec
func goPanicIndexU(x uint, y int) {
	panicCheck1(getcallerpc(), "index out of range")
	panic(boundsError{x: int64(x), signed: false, y: y, code: boundsIndex})
}

// failures in the comparisons for s[:x], 0 <= x <= y (y == len(s) or cap(s))
//
//go:yeswritebarrierrec
func goPanicSliceAlen(x int, y int) {
	panicCheck1(getcallerpc(), "slice bounds out of range")
	panic(boundsError{x: int64(x), signed: true, y: y, code: boundsSliceAlen})
}

//go:yeswritebarrierrec
func goPanicSliceAlenU(x uint, y int) {
	panicCheck1(getcallerpc(), "slice bounds out of range")
	panic(boundsError{x: int64(x), signed: false, y: y, code: boundsSliceAlen})
}

//go:yeswritebarrierrec
func goPanicSliceAcap(x int, y int) {
	panicCheck1(getcallerpc(), "slice bounds out of range")
	panic(boundsError{x: int64(x), signed: true, y: y, code: boundsSliceAcap})
}

//go:yeswritebarrierrec
func goPanicSliceAcapU(x uint, y int) {
	panicCheck1(getcallerpc(), "slice bounds out of range")
	panic(boundsError{x: int64(x), signed: false, y: y, code: boundsSliceAcap})
}

// failures in the comparisons for s[x:y], 0 <= x <= y
//
//go:yeswritebarrierrec
func goPanicSliceB(x int, y int) {
	panicCheck1(getcallerpc(), "slice bounds out of range")
	panic(boundsError{x: int64(x), signed: true, y: y, code: boundsSliceB})
}

//go:yeswritebarrierrec
func goPanicSliceBU(x uint, y int) {
	panicCheck1(getcallerpc(), "slice bounds out of range")
	panic(boundsError{x: int64(x), signed: false, y: y, code: boundsSliceB})
}

// failures in the comparisons for s[::x], 0 <= x <= y (y == len(s) or cap(s))
func goPanicSlice3Alen(x int, y int) {
	panicCheck1(getcallerpc(), "slice bounds out of range")
	panic(boundsError{x: int64(x), signed: true, y: y, code: boundsSlice3Alen})
}
func goPanicSlice3AlenU(x uint, y int) {
	panicCheck1(getcallerpc(), "slice bounds out of range")
	panic(boundsError{x: int64(x), signed: false, y: y, code: boundsSlice3Alen})
}
func goPanicSlice3Acap(x int, y int) {
	panicCheck1(getcallerpc(), "slice bounds out of range")
	panic(boundsError{x: int64(x), signed: true, y: y, code: boundsSlice3Acap})
}
func goPanicSlice3AcapU(x uint, y int) {
	panicCheck1(getcallerpc(), "slice bounds out of range")
	panic(boundsError{x: int64(x), signed: false, y: y, code: boundsSlice3Acap})
}

// failures in the comparisons for s[:x:y], 0 <= x <= y
func goPanicSlice3B(x int, y int) {
	panicCheck1(getcallerpc(), "slice bounds out of range")
	panic(boundsError{x: int64(x), signed: true, y: y, code: boundsSlice3B})
}
func goPanicSlice3BU(x uint, y int) {
	panicCheck1(getcallerpc(), "slice bounds out of range")
	panic(boundsError{x: int64(x), signed: false, y: y, code: boundsSlice3B})
}

// failures in the comparisons for s[x:y:], 0 <= x <= y
func goPanicSlice3C(x int, y int) {
	panicCheck1(getcallerpc(), "slice bounds out of range")
	panic(boundsError{x: int64(x), signed: true, y: y, code: boundsSlice3C})
}
func goPanicSlice3CU(x uint, y int) {
	panicCheck1(getcallerpc(), "slice bounds out of range")
	panic(boundsError{x: int64(x), signed: false, y: y, code: boundsSlice3C})
}

// failures in the conversion ([x]T)(s) or (*[x]T)(s), 0 <= x <= y, y == len(s)
func goPanicSliceConvert(x int, y int) {
	panicCheck1(getcallerpc(), "slice length too short to convert to array or pointer to array")
	panic(boundsError{x: int64(x), signed: true, y: y, code: boundsConvert})
}

// Implemented in assembly, as they take arguments in registers.
// Declared here to mark them as ABIInternal.
func panicIndex(x int, y int)
func panicIndexU(x uint, y int)
func panicSliceAlen(x int, y int)
func panicSliceAlenU(x uint, y int)
func panicSliceAcap(x int, y int)
func panicSliceAcapU(x uint, y int)
func panicSliceB(x int, y int)
func panicSliceBU(x uint, y int)
func panicSlice3Alen(x int, y int)
func panicSlice3AlenU(x uint, y int)
func panicSlice3Acap(x int, y int)
func panicSlice3AcapU(x uint, y int)
func panicSlice3B(x int, y int)
func panicSlice3BU(x uint, y int)
func panicSlice3C(x int, y int)
func panicSlice3CU(x uint, y int)
func panicSliceConvert(x int, y int)

var shiftError = error(errorString("negative shift amount"))

//go:yeswritebarrierrec
func panicshift() {
	panicCheck1(getcallerpc(), "negative shift amount")
	panic(shiftError)
}

var divideError = error(errorString("integer divide by zero"))

//go:yeswritebarrierrec
func panicdivide() {
	panicCheck2("integer divide by zero")
	panic(divideError)
}

var overflowError = error(errorString("integer overflow"))

func panicoverflow() {
	panicCheck2("integer overflow")
	panic(overflowError)
}

var floatError = error(errorString("floating point error"))

func panicfloat() {
	panicCheck2("floating point error")
	panic(floatError)
}

var memoryError = error(errorString("invalid memory address or nil pointer dereference"))

func panicmem() {
	panicCheck2("invalid memory address or nil pointer dereference")
	panic(memoryError)
}

func panicmemAddr(addr uintptr) {
	panicCheck2("invalid memory address or nil pointer dereference")
	panic(errorAddressString{msg: "invalid memory address or nil pointer dereference", addr: addr})
}

// Create a new deferred function fn, which has no arguments and results.
// The compiler turns a defer statement into a call to this.
func deferproc(fn func()) {
	gp := getg()
	if gp.m.curg != gp {
		// go code on the system stack can't defer
		throw("defer on system stack")
	}

	d := newdefer()
	d.link = gp._defer
	gp._defer = d
	d.fn = fn
	d.pc = getcallerpc()
	// We must not be preempted between calling getcallersp and
	// storing it to d.sp because getcallersp's result is a
	// uintptr stack pointer.
	d.sp = getcallersp()

	// deferproc returns 0 normally.
	// a deferred func that stops a panic
	// makes the deferproc return 1.
	// the code the compiler generates always
	// checks the return value and jumps to the
	// end of the function if deferproc returns != 0.
	return0()
	// No code can go here - the C return register has
	// been set and must not be clobbered.
}

// deferprocStack queues a new deferred function with a defer record on the stack.
// The defer record must have its fn field initialized.
// All other fields can contain junk.
// Nosplit because of the uninitialized pointer fields on the stack.
//
//go:nosplit
func deferprocStack(d *_defer) {
	gp := getg()
	if gp.m.curg != gp {
		// go code on the system stack can't defer
		throw("defer on system stack")
	}
	// fn is already set.
	// The other fields are junk on entry to deferprocStack and
	// are initialized here.
	d.heap = false
	d.sp = getcallersp()
	d.pc = getcallerpc()
	// The lines below implement:
	//   d.panic = nil
	//   d.fd = nil
	//   d.link = gp._defer
	//   gp._defer = d
	// But without write barriers. The first three are writes to
	// the stack so they don't need a write barrier, and furthermore
	// are to uninitialized memory, so they must not use a write barrier.
	// The fourth write does not require a write barrier because we
	// explicitly mark all the defer structures, so we don't need to
	// keep track of pointers to them with a write barrier.
	*(*uintptr)(unsafe.Pointer(&d.link)) = uintptr(unsafe.Pointer(gp._defer))
	*(*uintptr)(unsafe.Pointer(&gp._defer)) = uintptr(unsafe.Pointer(d))

	return0()
	// No code can go here - the C return register has
	// been set and must not be clobbered.
}

// Each P holds a pool for defers.

// Allocate a Defer, usually using per-P pool.
// Each defer must be released with freedefer.  The defer is not
// added to any defer chain yet.
func newdefer() *_defer {
	var d *_defer
	mp := acquirem()
	pp := mp.p.ptr()
	if len(pp.deferpool) == 0 && sched.deferpool != nil {
		lock(&sched.deferlock)
		for len(pp.deferpool) < cap(pp.deferpool)/2 && sched.deferpool != nil {
			d := sched.deferpool
			sched.deferpool = d.link
			d.link = nil
			pp.deferpool = append(pp.deferpool, d)
		}
		unlock(&sched.deferlock)
	}
	if n := len(pp.deferpool); n > 0 {
		d = pp.deferpool[n-1]
		pp.deferpool[n-1] = nil
		pp.deferpool = pp.deferpool[:n-1]
	}
	releasem(mp)
	mp, pp = nil, nil

	if d == nil {
		// Allocate new defer.
		d = new(_defer)
	}
	d.heap = true
	return d
}

// Free the given defer.
// The defer cannot be used after this call.
//
// This is nosplit because the incoming defer is in a perilous state.
// It's not on any defer list, so stack copying won't adjust stack
// pointers in it (namely, d.link). Hence, if we were to copy the
// stack, d could then contain a stale pointer.
//
//go:nosplit
func freedefer(d *_defer) {
	d.link = nil
	// After this point we can copy the stack.

	if d.fn != nil {
		freedeferfn()
	}
	if !d.heap {
		return
	}

	mp := acquirem()
	pp := mp.p.ptr()
	if len(pp.deferpool) == cap(pp.deferpool) {
		// Transfer half of local cache to the central cache.
		var first, last *_defer
		for len(pp.deferpool) > cap(pp.deferpool)/2 {
			n := len(pp.deferpool)
			d := pp.deferpool[n-1]
			pp.deferpool[n-1] = nil
			pp.deferpool = pp.deferpool[:n-1]
			if first == nil {
				first = d
			} else {
				last.link = d
			}
			last = d
		}
		lock(&sched.deferlock)
		last.link = sched.deferpool
		sched.deferpool = first
		unlock(&sched.deferlock)
	}

	*d = _defer{}

	pp.deferpool = append(pp.deferpool, d)

	releasem(mp)
	mp, pp = nil, nil
}

// Separate function so that it can split stack.
// Windows otherwise runs out of stack space.
func freedeferfn() {
	// fn must be cleared before d is unlinked from gp.
	throw("freedefer with d.fn != nil")
}

// deferreturn runs deferred functions for the caller's frame.
// The compiler inserts a call to this at the end of any
// function which calls defer.
func deferreturn() {
	var p _panic
	p.deferreturn = true

	p.start(getcallerpc(), unsafe.Pointer(getcallersp()))
	for {
		fn, ok := p.nextDefer()
		if !ok {
			break
		}
		fn()
	}
}

// Goexit terminates the goroutine that calls it. No other goroutine is affected.
// Goexit runs all deferred calls before terminating the goroutine. Because Goexit
// is not a panic, any recover calls in those deferred functions will return nil.
//
// Calling Goexit from the main goroutine terminates that goroutine
// without func main returning. Since func main has not returned,
// the program continues execution of other goroutines.
// If all other goroutines exit, the program crashes.
func Goexit() {
	// Create a panic object for Goexit, so we can recognize when it might be
	// bypassed by a recover().
	var p _panic
	p.goexit = true

	p.start(getcallerpc(), unsafe.Pointer(getcallersp()))
	for {
		fn, ok := p.nextDefer()
		if !ok {
			break
		}
		fn()
	}

	goexit1()
}

// Call all Error and String methods before freezing the world.
// Used when crashing with panicking.
func preprintpanics(p *_panic) {
	defer func() {
		text := "panic while printing panic value"
		switch r := recover().(type) {
		case nil:
			// nothing to do
		case string:
			throw(text + ": " + r)
		default:
			throw(text + ": type " + toRType(efaceOf(&r)._type).string())
		}
	}()
	for p != nil {
		switch v := p.arg.(type) {
		case error:
			p.arg = v.Error()
		case stringer:
			p.arg = v.String()
		}
		p = p.link
	}
}

// Print all currently active panics. Used when crashing.
// Should only be called after preprintpanics.
func printpanics(p *_panic) {
	if p.link != nil {
		printpanics(p.link)
		if !p.link.goexit {
			print("\t")
		}
	}
	if p.goexit {
		return
	}
	print("panic: ")
	printany(p.arg)
	if p.recovered {
		print(" [recovered]")
	}
	print("\n")
}

// readvarintUnsafe reads the uint32 in varint format starting at fd, and returns the
// uint32 and a pointer to the byte following the varint.
//
// There is a similar function runtime.readvarint, which takes a slice of bytes,
// rather than an unsafe pointer. These functions are duplicated, because one of
// the two use cases for the functions would get slower if the functions were
// combined.
func readvarintUnsafe(fd unsafe.Pointer) (uint32, unsafe.Pointer) {
	var r uint32
	var shift int
	for {
		b := *(*uint8)((unsafe.Pointer(fd)))
		fd = add(fd, unsafe.Sizeof(b))
		if b < 128 {
			return r + uint32(b)<<shift, fd
		}
		r += ((uint32(b) &^ 128) << shift)
		shift += 7
		if shift > 28 {
			panic("Bad varint")
		}
	}
}

// A PanicNilError happens when code calls panic(nil).
//
// Before Go 1.21, programs that called panic(nil) observed recover returning nil.
// Starting in Go 1.21, programs that call panic(nil) observe recover returning a *PanicNilError.
// Programs can change back to the old behavior by setting GODEBUG=panicnil=1.
type PanicNilError struct {
	// This field makes PanicNilError structurally different from
	// any other struct in this package, and the _ makes it different
	// from any struct in other packages too.
	// This avoids any accidental conversions being possible
	// between this struct and some other struct sharing the same fields,
	// like happened in go.dev/issue/56603.
	_ [0]*PanicNilError
}

func (*PanicNilError) Error() string { return "panic called with nil argument" }
func (*PanicNilError) RuntimeError() {}

var panicnil = &godebugInc{name: "panicnil"}

// The implementation of the predeclared function panic.
func gopanic(e any) {
	if e == nil {
		if debug.panicnil.Load() != 1 {
			e = new(PanicNilError)
		} else {
			panicnil.IncNonDefault()
		}
	}

	gp := getg()
	if gp.m.curg != gp {
		print("panic: ")
		printany(e)
		print("\n")
		throw("panic on system stack")
	}

	if gp.m.mallocing != 0 {
		print("panic: ")
		printany(e)
		print("\n")
		throw("panic during malloc")
	}
	if gp.m.preemptoff != "" {
		print("panic: ")
		printany(e)
		print("\n")
		print("preempt off reason: ")
		print(gp.m.preemptoff)
		print("\n")
		throw("panic during preemptoff")
	}
	if gp.m.locks != 0 {
		print("panic: ")
		printany(e)
		print("\n")
		throw("panic holding locks")
	}

	var p _panic
	p.arg = e

	runningPanicDefers.Add(1)

	p.start(getcallerpc(), unsafe.Pointer(getcallersp()))
	for {
		fn, ok := p.nextDefer()
		if !ok {
			break
		}
		fn()
	}

	// ran out of deferred calls - old-school panic now
	// Because it is unsafe to call arbitrary user code after freezing
	// the world, we call preprintpanics to invoke all necessary Error
	// and String methods to prepare the panic strings before startpanic.
	preprintpanics(&p)

	fatalpanic(&p)   // should not return
	*(*int)(nil) = 0 // not reached
}

// start initializes a panic to start unwinding the stack.
//
// If p.goexit is true, then start may return multiple times.
func (p *_panic) start(pc uintptr, sp unsafe.Pointer) {
	gp := getg()

	// Record the caller's PC and SP, so recovery can identify panics
	// that have been recovered. Also, so that if p is from Goexit, we
	// can restart its defer processing loop if a recovered panic tries
	// to jump past it.
	p.startPC = getcallerpc()
	p.startSP = unsafe.Pointer(getcallersp())

	if p.deferreturn {
		p.sp = sp

		if s := (*savedOpenDeferState)(gp.param); s != nil {
			// recovery saved some state for us, so that we can resume
			// calling open-coded defers without unwinding the stack.

			gp.param = nil

			p.retpc = s.retpc
			p.deferBitsPtr = (*byte)(add(sp, s.deferBitsOffset))
			p.slotsPtr = add(sp, s.slotsOffset)
		}

		return
	}

	p.link = gp._panic
	gp._panic = (*_panic)(noescape(unsafe.Pointer(p)))

	// Initialize state machine, and find the first frame with a defer.
	//
	// Note: We could use startPC and startSP here, but callers will
	// never have defer statements themselves. By starting at their
	// caller instead, we avoid needing to unwind through an extra
	// frame. It also somewhat simplifies the terminating condition for
	// deferreturn.
	p.lr, p.fp = pc, sp
	p.nextFrame()
}

// nextDefer returns the next deferred function to invoke, if any.
//
// Note: The "ok bool" result is necessary to correctly handle when
// the deferred function itself was nil (e.g., "defer (func())(nil)").
func (p *_panic) nextDefer() (func(), bool) {
	gp := getg()

	if !p.deferreturn {
		if gp._panic != p {
			throw("bad panic stack")
		}

		if p.recovered {
			mcall(recovery) // does not return
			throw("recovery failed")
		}
	}

	// The assembler adjusts p.argp in wrapper functions that shouldn't
	// be visible to recover(), so we need to restore it each iteration.
	p.argp = add(p.startSP, sys.MinFrameSize)

	for {
		for p.deferBitsPtr != nil {
			bits := *p.deferBitsPtr

			// Check whether any open-coded defers are still pending.
			//
			// Note: We need to check this upfront (rather than after
			// clearing the top bit) because it's possible that Goexit
			// invokes a deferred call, and there were still more pending
			// open-coded defers in the frame; but then the deferred call
			// panic and invoked the remaining defers in the frame, before
			// recovering and restarting the Goexit loop.
			if bits == 0 {
				p.deferBitsPtr = nil
				break
			}

			// Find index of top bit set.
			i := 7 - uintptr(sys.LeadingZeros8(bits))

			// Clear bit and store it back.
			bits &^= 1 << i
			*p.deferBitsPtr = bits

			return *(*func())(add(p.slotsPtr, i*goarch.PtrSize)), true
		}

		if d := gp._defer; d != nil && d.sp == uintptr(p.sp) {
			fn := d.fn
			d.fn = nil

			// TODO(mdempsky): Instead of having each deferproc call have
			// its own "deferreturn(); return" sequence, we should just make
			// them reuse the one we emit for open-coded defers.
			p.retpc = d.pc

			// Unlink and free.
			gp._defer = d.link
			freedefer(d)

			return fn, true
		}

		if !p.nextFrame() {
			return nil, false
		}
	}
}

// nextFrame finds the next frame that contains deferred calls, if any.
func (p *_panic) nextFrame() (ok bool) {
	if p.lr == 0 {
		return false
	}

	gp := getg()
	systemstack(func() {
		var limit uintptr
		if d := gp._defer; d != nil {
			limit = uintptr(d.sp)
		}

		var u unwinder
		u.initAt(p.lr, uintptr(p.fp), 0, gp, 0)
		for {
			if !u.valid() {
				p.lr = 0
				return // ok == false
			}

			// TODO(mdempsky): If we populate u.frame.fn.deferreturn for
			// every frame containing a defer (not just open-coded defers),
			// then we can simply loop until we find the next frame where
			// it's non-zero.

			if u.frame.sp == limit {
				break // found a frame with linked defers
			}

			if p.initOpenCodedDefers(u.frame.fn, unsafe.Pointer(u.frame.varp)) {
				break // found a frame with open-coded defers
			}

			u.next()
		}

		p.lr = u.frame.lr
		p.sp = unsafe.Pointer(u.frame.sp)
		p.fp = unsafe.Pointer(u.frame.fp)

		ok = true
	})

	return
}

func (p *_panic) initOpenCodedDefers(fn funcInfo, varp unsafe.Pointer) bool {
	fd := funcdata(fn, abi.FUNCDATA_OpenCodedDeferInfo)
	if fd == nil {
		return false
	}

	if fn.deferreturn == 0 {
		throw("missing deferreturn")
	}

	deferBitsOffset, fd := readvarintUnsafe(fd)
	deferBitsPtr := (*uint8)(add(varp, -uintptr(deferBitsOffset)))
	if *deferBitsPtr == 0 {
		return false // has open-coded defers, but none pending
	}

	slotsOffset, fd := readvarintUnsafe(fd)

	p.retpc = fn.entry() + uintptr(fn.deferreturn)
	p.deferBitsPtr = deferBitsPtr
	p.slotsPtr = add(varp, -uintptr(slotsOffset))

	return true
}

// The implementation of the predeclared function recover.
// Cannot split the stack because it needs to reliably
// find the stack segment of its caller.
//
// TODO(rsc): Once we commit to CopyStackAlways,
// this doesn't need to be nosplit.
//
//go:nosplit
func gorecover(argp uintptr) any {
	// Must be in a function running as part of a deferred call during the panic.
	// Must be called from the topmost function of the call
	// (the function used in the defer statement).
	// p.argp is the argument pointer of that topmost deferred function call.
	// Compare against argp reported by caller.
	// If they match, the caller is the one who can recover.
	gp := getg()
	p := gp._panic
	if p != nil && !p.goexit && !p.recovered && argp == uintptr(p.argp) {
		p.recovered = true
		return p.arg
	}
	return nil
}

//go:linkname sync_throw sync.throw
func sync_throw(s string) {
	throw(s)
}

//go:linkname sync_fatal sync.fatal
func sync_fatal(s string) {
	fatal(s)
}

// throw triggers a fatal error that dumps a stack trace and exits.
//
// throw should be used for runtime-internal fatal errors where Go itself,
// rather than user code, may be at fault for the failure.
//
//go:nosplit
func throw(s string) {
	// Everything throw does should be recursively nosplit so it
	// can be called even when it's unsafe to grow the stack.
	systemstack(func() {
		print("fatal error: ", s, "\n")
	})

	fatalthrow(throwTypeRuntime)
}

// fatal triggers a fatal error that dumps a stack trace and exits.
//
// fatal is equivalent to throw, but is used when user code is expected to be
// at fault for the failure, such as racing map writes.
//
// fatal does not include runtime frames, system goroutines, or frame metadata
// (fp, sp, pc) in the stack trace unless GOTRACEBACK=system or higher.
//
//go:nosplit
func fatal(s string) {
	// Everything fatal does should be recursively nosplit so it
	// can be called even when it's unsafe to grow the stack.
	systemstack(func() {
		print("fatal error: ", s, "\n")
	})

	fatalthrow(throwTypeUser)
}

// runningPanicDefers is non-zero while running deferred functions for panic.
// This is used to try hard to get a panic stack trace out when exiting.
var runningPanicDefers atomic.Uint32

// panicking is non-zero when crashing the program for an unrecovered panic.
var panicking atomic.Uint32

// paniclk is held while printing the panic information and stack trace,
// so that two concurrent panics don't overlap their output.
var paniclk mutex

// Unwind the stack after a deferred function calls recover
// after a panic. Then arrange to continue running as though
// the caller of the deferred function returned normally.
//
// However, if unwinding the stack would skip over a Goexit call, we
// return into the Goexit loop instead, so it can continue processing
// defers instead.
func recovery(gp *g) {
	p := gp._panic
	pc, sp := p.retpc, uintptr(p.sp)
	p0, saveOpenDeferState := p, p.deferBitsPtr != nil && *p.deferBitsPtr != 0

	// Unwind the panic stack.
	for ; p != nil && uintptr(p.startSP) < sp; p = p.link {
		// Don't allow jumping past a pending Goexit.
		// Instead, have its _panic.start() call return again.
		//
		// TODO(mdempsky): In this case, Goexit will resume walking the
		// stack where it left off, which means it will need to rewalk
		// frames that we've already processed.
		//
		// There's a similar issue with nested panics, when the inner
		// panic supercedes the outer panic. Again, we end up needing to
		// walk the same stack frames.
		//
		// These are probably pretty rare occurrences in practice, and
		// they don't seem any worse than the existing logic. But if we
		// move the unwinding state into _panic, we could detect when we
		// run into where the last panic started, and then just pick up
		// where it left off instead.
		//
		// With how subtle defer handling is, this might not actually be
		// worthwhile though.
		if p.goexit {
			pc, sp = p.startPC, uintptr(p.startSP)
			saveOpenDeferState = false // goexit is unwinding the stack anyway
			break
		}

		runningPanicDefers.Add(-1)
	}
	gp._panic = p

	if p == nil { // must be done with signal
		gp.sig = 0
	}

	if gp.param != nil {
		throw("unexpected gp.param")
	}
	if saveOpenDeferState {
		// If we're returning to deferreturn and there are more open-coded
		// defers for it to call, save enough state for it to be able to
		// pick up where p0 left off.
		gp.param = unsafe.Pointer(&savedOpenDeferState{
			retpc: p0.retpc,

			// We need to save deferBitsPtr and slotsPtr too, but those are
			// stack pointers. To avoid issues around heap objects pointing
			// to the stack, save them as offsets from SP.
			deferBitsOffset: uintptr(unsafe.Pointer(p0.deferBitsPtr)) - uintptr(p0.sp),
			slotsOffset:     uintptr(p0.slotsPtr) - uintptr(p0.sp),
		})
	}

	// TODO(mdempsky): Currently, we rely on frames containing "defer"
	// to end with "CALL deferreturn; RET". This allows deferreturn to
	// finish running any pending defers in the frame.
	//
	// But we should be able to tell whether there are still pending
	// defers here. If there aren't, we can just jump directly to the
	// "RET" instruction. And if there are, we don't need an actual
	// "CALL deferreturn" instruction; we can simulate it with something
	// like:
	//
	//	if usesLR {
	//		lr = pc
	//	} else {
	//		sp -= sizeof(pc)
	//		*(*uintptr)(sp) = pc
	//	}
	//	pc = funcPC(deferreturn)
	//
	// So that we effectively tail call into deferreturn, such that it
	// then returns to the simple "RET" epilogue. That would save the
	// overhead of the "deferreturn" call when there aren't actually any
	// pending defers left, and shrink the TEXT size of compiled
	// binaries. (Admittedly, both of these are modest savings.)

	// Ensure we're recovering within the appropriate stack.
	if sp != 0 && (sp < gp.stack.lo || gp.stack.hi < sp) {
		print("recover: ", hex(sp), " not in [", hex(gp.stack.lo), ", ", hex(gp.stack.hi), "]\n")
		throw("bad recovery")
	}

	// Make the deferproc for this d return again,
	// this time returning 1. The calling function will
	// jump to the standard return epilogue.
	gp.sched.sp = sp
	gp.sched.pc = pc
	gp.sched.lr = 0
	gp.sched.ret = 1
	gogo(&gp.sched)
}

// fatalthrow implements an unrecoverable runtime throw. It freezes the
// system, prints stack traces starting from its caller, and terminates the
// process.
//
//go:nosplit
func fatalthrow(t throwType) {
	pc := getcallerpc()
	sp := getcallersp()
	gp := getg()

	if gp.m.throwing == throwTypeNone {
		gp.m.throwing = t
	}

	// Switch to the system stack to avoid any stack growth, which may make
	// things worse if the runtime is in a bad state.
	systemstack(func() {
		if isSecureMode() {
			exit(2)
		}

		startpanic_m()

		if dopanic_m(gp, pc, sp) {
			// crash uses a decent amount of nosplit stack and we're already
			// low on stack in throw, so crash on the system stack (unlike
			// fatalpanic).
			crash()
		}

		exit(2)
	})

	*(*int)(nil) = 0 // not reached
}

// fatalpanic implements an unrecoverable panic. It is like fatalthrow, except
// that if msgs != nil, fatalpanic also prints panic messages and decrements
// runningPanicDefers once main is blocked from exiting.
//
//go:nosplit
func fatalpanic(msgs *_panic) {
	pc := getcallerpc()
	sp := getcallersp()
	gp := getg()
	var docrash bool
	// Switch to the system stack to avoid any stack growth, which
	// may make things worse if the runtime is in a bad state.
	systemstack(func() {
		if startpanic_m() && msgs != nil {
			// There were panic messages and startpanic_m
			// says it's okay to try to print them.

			// startpanic_m set panicking, which will
			// block main from exiting, so now OK to
			// decrement runningPanicDefers.
			runningPanicDefers.Add(-1)

			printpanics(msgs)
		}

		docrash = dopanic_m(gp, pc, sp)
	})

	if docrash {
		// By crashing outside the above systemstack call, debuggers
		// will not be confused when generating a backtrace.
		// Function crash is marked nosplit to avoid stack growth.
		crash()
	}

	systemstack(func() {
		exit(2)
	})

	*(*int)(nil) = 0 // not reached
}

// startpanic_m prepares for an unrecoverable panic.
//
// It returns true if panic messages should be printed, or false if
// the runtime is in bad shape and should just print stacks.
//
// It must not have write barriers even though the write barrier
// explicitly ignores writes once dying > 0. Write barriers still
// assume that g.m.p != nil, and this function may not have P
// in some contexts (e.g. a panic in a signal handler for a signal
// sent to an M with no P).
//
//go:nowritebarrierrec
func startpanic_m() bool {
	gp := getg()
	if mheap_.cachealloc.size == 0 { // very early
		print("runtime: panic before malloc heap initialized\n")
	}
	// Disallow malloc during an unrecoverable panic. A panic
	// could happen in a signal handler, or in a throw, or inside
	// malloc itself. We want to catch if an allocation ever does
	// happen (even if we're not in one of these situations).
	gp.m.mallocing++

	// If we're dying because of a bad lock count, set it to a
	// good lock count so we don't recursively panic below.
	if gp.m.locks < 0 {
		gp.m.locks = 1
	}

	switch gp.m.dying {
	case 0:
		// Setting dying >0 has the side-effect of disabling this G's writebuf.
		gp.m.dying = 1
		panicking.Add(1)
		lock(&paniclk)
		if debug.schedtrace > 0 || debug.scheddetail > 0 {
			schedtrace(true)
		}
		freezetheworld()
		return true
	case 1:
		// Something failed while panicking.
		// Just print a stack trace and exit.
		gp.m.dying = 2
		print("panic during panic\n")
		return false
	case 2:
		// This is a genuine bug in the runtime, we couldn't even
		// print the stack trace successfully.
		gp.m.dying = 3
		print("stack trace unavailable\n")
		exit(4)
		fallthrough
	default:
		// Can't even print! Just exit.
		exit(5)
		return false // Need to return something.
	}
}

var didothers bool
var deadlock mutex

// gp is the crashing g running on this M, but may be a user G, while getg() is
// always g0.
func dopanic_m(gp *g, pc, sp uintptr) bool {
	if gp.sig != 0 {
		signame := signame(gp.sig)
		if signame != "" {
			print("[signal ", signame)
		} else {
			print("[signal ", hex(gp.sig))
		}
		print(" code=", hex(gp.sigcode0), " addr=", hex(gp.sigcode1), " pc=", hex(gp.sigpc), "]\n")
	}

	level, all, docrash := gotraceback()
	if level > 0 {
		if gp != gp.m.curg {
			all = true
		}
		if gp != gp.m.g0 {
			print("\n")
			goroutineheader(gp)
			traceback(pc, sp, 0, gp)
		} else if level >= 2 || gp.m.throwing >= throwTypeRuntime {
			print("\nruntime stack:\n")
			traceback(pc, sp, 0, gp)
		}
		if !didothers && all {
			didothers = true
			tracebackothers(gp)
		}
	}
	unlock(&paniclk)

	if panicking.Add(-1) != 0 {
		// Some other m is panicking too.
		// Let it print what it needs to print.
		// Wait forever without chewing up cpu.
		// It will exit when it's done.
		lock(&deadlock)
		lock(&deadlock)
	}

	printDebugLog()

	return docrash
}

// canpanic returns false if a signal should throw instead of
// panicking.
//
//go:nosplit
func canpanic() bool {
	gp := getg()
	mp := acquirem()

	// Is it okay for gp to panic instead of crashing the program?
	// Yes, as long as it is running Go code, not runtime code,
	// and not stuck in a system call.
	if gp != mp.curg {
		releasem(mp)
		return false
	}
	// N.B. mp.locks != 1 instead of 0 to account for acquirem.
	if mp.locks != 1 || mp.mallocing != 0 || mp.throwing != throwTypeNone || mp.preemptoff != "" || mp.dying != 0 {
		releasem(mp)
		return false
	}
	status := readgstatus(gp)
	if status&^_Gscan != _Grunning || gp.syscallsp != 0 {
		releasem(mp)
		return false
	}
	if GOOS == "windows" && mp.libcallsp != 0 {
		releasem(mp)
		return false
	}
	releasem(mp)
	return true
}

// shouldPushSigpanic reports whether pc should be used as sigpanic's
// return PC (pushing a frame for the call). Otherwise, it should be
// left alone so that LR is used as sigpanic's return PC, effectively
// replacing the top-most frame with sigpanic. This is used by
// preparePanic.
func shouldPushSigpanic(gp *g, pc, lr uintptr) bool {
	if pc == 0 {
		// Probably a call to a nil func. The old LR is more
		// useful in the stack trace. Not pushing the frame
		// will make the trace look like a call to sigpanic
		// instead. (Otherwise the trace will end at sigpanic
		// and we won't get to see who faulted.)
		return false
	}
	// If we don't recognize the PC as code, but we do recognize
	// the link register as code, then this assumes the panic was
	// caused by a call to non-code. In this case, we want to
	// ignore this call to make unwinding show the context.
	//
	// If we running C code, we're not going to recognize pc as a
	// Go function, so just assume it's good. Otherwise, traceback
	// may try to read a stale LR that looks like a Go code
	// pointer and wander into the woods.
	if gp.m.incgo || findfunc(pc).valid() {
		// This wasn't a bad call, so use PC as sigpanic's
		// return PC.
		return true
	}
	if findfunc(lr).valid() {
		// This was a bad call, but the LR is good, so use the
		// LR as sigpanic's return PC.
		return false
	}
	// Neither the PC or LR is good. Hopefully pushing a frame
	// will work.
	return true
}

// isAbortPC reports whether pc is the program counter at which
// runtime.abort raises a signal.
//
// It is nosplit because it's part of the isgoexception
// implementation.
//
//go:nosplit
func isAbortPC(pc uintptr) bool {
	f := findfunc(pc)
	if !f.valid() {
		return false
	}
	return f.funcID == abi.FuncID_abort
}
