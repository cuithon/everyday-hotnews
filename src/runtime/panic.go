// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"internal/goarch"
	"runtime/internal/atomic"
	"runtime/internal/sys"
	"unsafe"
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
func goPanicIndex(x int, y int) {
	panicCheck1(getcallerpc(), "index out of range")
	panic(boundsError{x: int64(x), signed: true, y: y, code: boundsIndex})
}
func goPanicIndexU(x uint, y int) {
	panicCheck1(getcallerpc(), "index out of range")
	panic(boundsError{x: int64(x), signed: false, y: y, code: boundsIndex})
}

// failures in the comparisons for s[:x], 0 <= x <= y (y == len(s) or cap(s))
func goPanicSliceAlen(x int, y int) {
	panicCheck1(getcallerpc(), "slice bounds out of range")
	panic(boundsError{x: int64(x), signed: true, y: y, code: boundsSliceAlen})
}
func goPanicSliceAlenU(x uint, y int) {
	panicCheck1(getcallerpc(), "slice bounds out of range")
	panic(boundsError{x: int64(x), signed: false, y: y, code: boundsSliceAlen})
}
func goPanicSliceAcap(x int, y int) {
	panicCheck1(getcallerpc(), "slice bounds out of range")
	panic(boundsError{x: int64(x), signed: true, y: y, code: boundsSliceAcap})
}
func goPanicSliceAcapU(x uint, y int) {
	panicCheck1(getcallerpc(), "slice bounds out of range")
	panic(boundsError{x: int64(x), signed: false, y: y, code: boundsSliceAcap})
}

// failures in the comparisons for s[x:y], 0 <= x <= y
func goPanicSliceB(x int, y int) {
	panicCheck1(getcallerpc(), "slice bounds out of range")
	panic(boundsError{x: int64(x), signed: true, y: y, code: boundsSliceB})
}
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

// failures in the conversion (*[x]T)s, 0 <= x <= y, x == cap(s)
func goPanicSliceConvert(x int, y int) {
	panicCheck1(getcallerpc(), "slice length too short to convert to pointer to array")
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

func panicshift() {
	panicCheck1(getcallerpc(), "negative shift amount")
	panic(shiftError)
}

var divideError = error(errorString("integer divide by zero"))

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
	if d._panic != nil {
		throw("deferproc: d.panic != nil after newdefer")
	}
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
	d.started = false
	d.heap = false
	d.openDefer = false
	d.sp = getcallersp()
	d.pc = getcallerpc()
	d.framepc = 0
	d.varp = 0
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
	*(*uintptr)(unsafe.Pointer(&d._panic)) = 0
	*(*uintptr)(unsafe.Pointer(&d.fd)) = 0
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
func freedefer(d *_defer) {
	if d._panic != nil {
		freedeferpanic()
	}
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

	// These lines used to be simply `*d = _defer{}` but that
	// started causing a nosplit stack overflow via typedmemmove.
	d.started = false
	d.openDefer = false
	d.sp = 0
	d.pc = 0
	d.framepc = 0
	d.varp = 0
	d.fd = nil
	// d._panic and d.fn must be nil already.
	// If not, we would have called freedeferpanic or freedeferfn above,
	// both of which throw.
	d.link = nil

	pp.deferpool = append(pp.deferpool, d)

	releasem(mp)
	mp, pp = nil, nil
}

// Separate function so that it can split stack.
// Windows otherwise runs out of stack space.
func freedeferpanic() {
	// _panic must be cleared before d is unlinked from gp.
	throw("freedefer with d._panic != nil")
}

func freedeferfn() {
	// fn must be cleared before d is unlinked from gp.
	throw("freedefer with d.fn != nil")
}

// deferreturn runs deferred functions for the caller's frame.
// The compiler inserts a call to this at the end of any
// function which calls defer.
func deferreturn() {
	gp := getg()
	for {
		d := gp._defer
		if d == nil {
			return
		}
		sp := getcallersp()
		if d.sp != sp {
			return
		}
		if d.openDefer {
			done := runOpenDeferFrame(gp, d)
			if !done {
				throw("unfinished open-coded defers in deferreturn")
			}
			gp._defer = d.link
			freedefer(d)
			// If this frame uses open defers, then this
			// must be the only defer record for the
			// frame, so we can just return.
			return
		}

		fn := d.fn
		d.fn = nil
		gp._defer = d.link
		freedefer(d)
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
	// Run all deferred functions for the current goroutine.
	// This code is similar to gopanic, see that implementation
	// for detailed comments.
	gp := getg()

	// Create a panic object for Goexit, so we can recognize when it might be
	// bypassed by a recover().
	var p _panic
	p.goexit = true
	p.link = gp._panic
	gp._panic = (*_panic)(noescape(unsafe.Pointer(&p)))

	addOneOpenDeferFrame(gp, getcallerpc(), unsafe.Pointer(getcallersp()))
	for {
		d := gp._defer
		if d == nil {
			break
		}
		if d.started {
			if d._panic != nil {
				d._panic.aborted = true
				d._panic = nil
			}
			if !d.openDefer {
				d.fn = nil
				gp._defer = d.link
				freedefer(d)
				continue
			}
		}
		d.started = true
		d._panic = (*_panic)(noescape(unsafe.Pointer(&p)))
		if d.openDefer {
			done := runOpenDeferFrame(gp, d)
			if !done {
				// We should always run all defers in the frame,
				// since there is no panic associated with this
				// defer that can be recovered.
				throw("unfinished open-coded defers in Goexit")
			}
			if p.aborted {
				// Since our current defer caused a panic and may
				// have been already freed, just restart scanning
				// for open-coded defers from this frame again.
				addOneOpenDeferFrame(gp, getcallerpc(), unsafe.Pointer(getcallersp()))
			} else {
				addOneOpenDeferFrame(gp, 0, nil)
			}
		} else {
			// Save the pc/sp in deferCallSave(), so we can "recover" back to this
			// loop if necessary.
			deferCallSave(&p, d.fn)
		}
		if p.aborted {
			// We had a recursive panic in the defer d we started, and
			// then did a recover in a defer that was further down the
			// defer chain than d. In the case of an outstanding Goexit,
			// we force the recover to return back to this loop. d will
			// have already been freed if completed, so just continue
			// immediately to the next defer on the chain.
			p.aborted = false
			continue
		}
		if gp._defer != d {
			throw("bad defer entry in Goexit")
		}
		d._panic = nil
		d.fn = nil
		gp._defer = d.link
		freedefer(d)
		// Note: we ignore recovers here because Goexit isn't a panic
	}
	goexit1()
}

// Call all Error and String methods before freezing the world.
// Used when crashing with panicking.
func preprintpanics(p *_panic) {
	defer func() {
		if recover() != nil {
			throw("panic while printing panic value")
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

// addOneOpenDeferFrame scans the stack for the first frame (if any) with
// open-coded defers and if it finds one, adds a single record to the defer chain
// for that frame. If sp is non-nil, it starts the stack scan from the frame
// specified by sp. If sp is nil, it uses the sp from the current defer record
// (which has just been finished). Hence, it continues the stack scan from the
// frame of the defer that just finished. It skips any frame that already has an
// open-coded _defer record, which would have been created from a previous
// (unrecovered) panic.
//
// Note: All entries of the defer chain (including this new open-coded entry) have
// their pointers (including sp) adjusted properly if the stack moves while
// running deferred functions. Also, it is safe to pass in the sp arg (which is
// the direct result of calling getcallersp()), because all pointer variables
// (including arguments) are adjusted as needed during stack copies.
func addOneOpenDeferFrame(gp *g, pc uintptr, sp unsafe.Pointer) {
	var prevDefer *_defer
	if sp == nil {
		prevDefer = gp._defer
		pc = prevDefer.framepc
		sp = unsafe.Pointer(prevDefer.sp)
	}
	systemstack(func() {
		gentraceback(pc, uintptr(sp), 0, gp, 0, nil, 0x7fffffff,
			func(frame *stkframe, unused unsafe.Pointer) bool {
				if prevDefer != nil && prevDefer.sp == frame.sp {
					// Skip the frame for the previous defer that
					// we just finished (and was used to set
					// where we restarted the stack scan)
					return true
				}
				f := frame.fn
				fd := funcdata(f, _FUNCDATA_OpenCodedDeferInfo)
				if fd == nil {
					return true
				}
				// Insert the open defer record in the
				// chain, in order sorted by sp.
				d := gp._defer
				var prev *_defer
				for d != nil {
					dsp := d.sp
					if frame.sp < dsp {
						break
					}
					if frame.sp == dsp {
						if !d.openDefer {
							throw("duplicated defer entry")
						}
						return true
					}
					prev = d
					d = d.link
				}
				if frame.fn.deferreturn == 0 {
					throw("missing deferreturn")
				}

				d1 := newdefer()
				d1.openDefer = true
				d1._panic = nil
				// These are the pc/sp to set after we've
				// run a defer in this frame that did a
				// recover. We return to a special
				// deferreturn that runs any remaining
				// defers and then returns from the
				// function.
				d1.pc = frame.fn.entry + uintptr(frame.fn.deferreturn)
				d1.varp = frame.varp
				d1.fd = fd
				// Save the SP/PC associated with current frame,
				// so we can continue stack trace later if needed.
				d1.framepc = frame.pc
				d1.sp = frame.sp
				d1.link = d
				if prev == nil {
					gp._defer = d1
				} else {
					prev.link = d1
				}
				// Stop stack scanning after adding one open defer record
				return false
			},
			nil, 0)
	})
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

// runOpenDeferFrame runs the active open-coded defers in the frame specified by
// d. It normally processes all active defers in the frame, but stops immediately
// if a defer does a successful recover. It returns true if there are no
// remaining defers to run in the frame.
func runOpenDeferFrame(gp *g, d *_defer) bool {
	done := true
	fd := d.fd

	deferBitsOffset, fd := readvarintUnsafe(fd)
	nDefers, fd := readvarintUnsafe(fd)
	deferBits := *(*uint8)(unsafe.Pointer(d.varp - uintptr(deferBitsOffset)))

	for i := int(nDefers) - 1; i >= 0; i-- {
		// read the funcdata info for this defer
		var closureOffset uint32
		closureOffset, fd = readvarintUnsafe(fd)
		if deferBits&(1<<i) == 0 {
			continue
		}
		closure := *(*func())(unsafe.Pointer(d.varp - uintptr(closureOffset)))
		d.fn = closure
		deferBits = deferBits &^ (1 << i)
		*(*uint8)(unsafe.Pointer(d.varp - uintptr(deferBitsOffset))) = deferBits
		p := d._panic
		// Call the defer. Note that this can change d.varp if
		// the stack moves.
		deferCallSave(p, d.fn)
		if p != nil && p.aborted {
			break
		}
		d.fn = nil
		if d._panic != nil && d._panic.recovered {
			done = deferBits == 0
			break
		}
	}

	return done
}

// deferCallSave calls fn() after saving the caller's pc and sp in the
// panic record. This allows the runtime to return to the Goexit defer
// processing loop, in the unusual case where the Goexit may be
// bypassed by a successful recover.
//
// This is marked as a wrapper by the compiler so it doesn't appear in
// tracebacks.
func deferCallSave(p *_panic, fn func()) {
	if p != nil {
		p.argp = unsafe.Pointer(getargp())
		p.pc = getcallerpc()
		p.sp = unsafe.Pointer(getcallersp())
	}
	fn()
	if p != nil {
		p.pc = 0
		p.sp = unsafe.Pointer(nil)
	}
}

// The implementation of the predeclared function panic.
func gopanic(e interface{}) {
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
	p.link = gp._panic
	gp._panic = (*_panic)(noescape(unsafe.Pointer(&p)))

	atomic.Xadd(&runningPanicDefers, 1)

	// By calculating getcallerpc/getcallersp here, we avoid scanning the
	// gopanic frame (stack scanning is slow...)
	addOneOpenDeferFrame(gp, getcallerpc(), unsafe.Pointer(getcallersp()))

	for {
		d := gp._defer
		if d == nil {
			break
		}

		// If defer was started by earlier panic or Goexit (and, since we're back here, that triggered a new panic),
		// take defer off list. An earlier panic will not continue running, but we will make sure below that an
		// earlier Goexit does continue running.
		if d.started {
			if d._panic != nil {
				d._panic.aborted = true
			}
			d._panic = nil
			if !d.openDefer {
				// For open-coded defers, we need to process the
				// defer again, in case there are any other defers
				// to call in the frame (not including the defer
				// call that caused the panic).
				d.fn = nil
				gp._defer = d.link
				freedefer(d)
				continue
			}
		}

		// Mark defer as started, but keep on list, so that traceback
		// can find and update the defer's argument frame if stack growth
		// or a garbage collection happens before executing d.fn.
		d.started = true

		// Record the panic that is running the defer.
		// If there is a new panic during the deferred call, that panic
		// will find d in the list and will mark d._panic (this panic) aborted.
		d._panic = (*_panic)(noescape(unsafe.Pointer(&p)))

		done := true
		if d.openDefer {
			done = runOpenDeferFrame(gp, d)
			if done && !d._panic.recovered {
				addOneOpenDeferFrame(gp, 0, nil)
			}
		} else {
			p.argp = unsafe.Pointer(getargp())
			d.fn()
		}
		p.argp = nil

		// Deferred function did not panic. Remove d.
		if gp._defer != d {
			throw("bad defer entry in panic")
		}
		d._panic = nil

		// trigger shrinkage to test stack copy. See stack_test.go:TestStackPanic
		//GC()

		pc := d.pc
		sp := unsafe.Pointer(d.sp) // must be pointer so it gets adjusted during stack copy
		if done {
			d.fn = nil
			gp._defer = d.link
			freedefer(d)
		}
		if p.recovered {
			gp._panic = p.link
			if gp._panic != nil && gp._panic.goexit && gp._panic.aborted {
				// A normal recover would bypass/abort the Goexit.  Instead,
				// we return to the processing loop of the Goexit.
				gp.sigcode0 = uintptr(gp._panic.sp)
				gp.sigcode1 = uintptr(gp._panic.pc)
				mcall(recovery)
				throw("bypassed recovery failed") // mcall should not return
			}
			atomic.Xadd(&runningPanicDefers, -1)

			// Remove any remaining non-started, open-coded
			// defer entries after a recover, since the
			// corresponding defers will be executed normally
			// (inline). Any such entry will become stale once
			// we run the corresponding defers inline and exit
			// the associated stack frame.
			d := gp._defer
			var prev *_defer
			if !done {
				// Skip our current frame, if not done. It is
				// needed to complete any remaining defers in
				// deferreturn()
				prev = d
				d = d.link
			}
			for d != nil {
				if d.started {
					// This defer is started but we
					// are in the middle of a
					// defer-panic-recover inside of
					// it, so don't remove it or any
					// further defer entries
					break
				}
				if d.openDefer {
					if prev == nil {
						gp._defer = d.link
					} else {
						prev.link = d.link
					}
					newd := d.link
					freedefer(d)
					d = newd
				} else {
					prev = d
					d = d.link
				}
			}

			gp._panic = p.link
			// Aborted panics are marked but remain on the g.panic list.
			// Remove them from the list.
			for gp._panic != nil && gp._panic.aborted {
				gp._panic = gp._panic.link
			}
			if gp._panic == nil { // must be done with signal
				gp.sig = 0
			}
			// Pass information about recovering frame to recovery.
			gp.sigcode0 = uintptr(sp)
			gp.sigcode1 = pc
			mcall(recovery)
			throw("recovery failed") // mcall should not return
		}
	}

	// ran out of deferred calls - old-school panic now
	// Because it is unsafe to call arbitrary user code after freezing
	// the world, we call preprintpanics to invoke all necessary Error
	// and String methods to prepare the panic strings before startpanic.
	preprintpanics(gp._panic)

	fatalpanic(gp._panic) // should not return
	*(*int)(nil) = 0      // not reached
}

// getargp returns the location where the caller
// writes outgoing function call arguments.
//go:nosplit
//go:noinline
func getargp() uintptr {
	return getcallersp() + sys.MinFrameSize
}

// The implementation of the predeclared function recover.
// Cannot split the stack because it needs to reliably
// find the stack segment of its caller.
//
// TODO(rsc): Once we commit to CopyStackAlways,
// this doesn't need to be nosplit.
//go:nosplit
func gorecover(argp uintptr) interface{} {
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

//go:nosplit
func throw(s string) {
	// Everything throw does should be recursively nosplit so it
	// can be called even when it's unsafe to grow the stack.
	systemstack(func() {
		print("fatal error: ", s, "\n")
	})
	gp := getg()
	if gp.m.throwing == 0 {
		gp.m.throwing = 1
	}
	fatalthrow()
	*(*int)(nil) = 0 // not reached
}

// runningPanicDefers is non-zero while running deferred functions for panic.
// runningPanicDefers is incremented and decremented atomically.
// This is used to try hard to get a panic stack trace out when exiting.
var runningPanicDefers uint32

// panicking is non-zero when crashing the program for an unrecovered panic.
// panicking is incremented and decremented atomically.
var panicking uint32

// paniclk is held while printing the panic information and stack trace,
// so that two concurrent panics don't overlap their output.
var paniclk mutex

// Unwind the stack after a deferred function calls recover
// after a panic. Then arrange to continue running as though
// the caller of the deferred function returned normally.
func recovery(gp *g) {
	// Info about defer passed in G struct.
	sp := gp.sigcode0
	pc := gp.sigcode1

	// d's arguments need to be in the stack.
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
func fatalthrow() {
	pc := getcallerpc()
	sp := getcallersp()
	gp := getg()
	// Switch to the system stack to avoid any stack growth, which
	// may make things worse if the runtime is in a bad state.
	systemstack(func() {
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
			atomic.Xadd(&runningPanicDefers, -1)

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
	_g_ := getg()
	if mheap_.cachealloc.size == 0 { // very early
		print("runtime: panic before malloc heap initialized\n")
	}
	// Disallow malloc during an unrecoverable panic. A panic
	// could happen in a signal handler, or in a throw, or inside
	// malloc itself. We want to catch if an allocation ever does
	// happen (even if we're not in one of these situations).
	_g_.m.mallocing++

	// If we're dying because of a bad lock count, set it to a
	// good lock count so we don't recursively panic below.
	if _g_.m.locks < 0 {
		_g_.m.locks = 1
	}

	switch _g_.m.dying {
	case 0:
		// Setting dying >0 has the side-effect of disabling this G's writebuf.
		_g_.m.dying = 1
		atomic.Xadd(&panicking, 1)
		lock(&paniclk)
		if debug.schedtrace > 0 || debug.scheddetail > 0 {
			schedtrace(true)
		}
		freezetheworld()
		return true
	case 1:
		// Something failed while panicking.
		// Just print a stack trace and exit.
		_g_.m.dying = 2
		print("panic during panic\n")
		return false
	case 2:
		// This is a genuine bug in the runtime, we couldn't even
		// print the stack trace successfully.
		_g_.m.dying = 3
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
	_g_ := getg()
	if level > 0 {
		if gp != gp.m.curg {
			all = true
		}
		if gp != gp.m.g0 {
			print("\n")
			goroutineheader(gp)
			traceback(pc, sp, 0, gp)
		} else if level >= 2 || _g_.m.throwing > 0 {
			print("\nruntime stack:\n")
			traceback(pc, sp, 0, gp)
		}
		if !didothers && all {
			didothers = true
			tracebackothers(gp)
		}
	}
	unlock(&paniclk)

	if atomic.Xadd(&panicking, -1) != 0 {
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
func canpanic(gp *g) bool {
	// Note that g is m->gsignal, different from gp.
	// Note also that g->m can change at preemption, so m can go stale
	// if this function ever makes a function call.
	_g_ := getg()
	_m_ := _g_.m

	// Is it okay for gp to panic instead of crashing the program?
	// Yes, as long as it is running Go code, not runtime code,
	// and not stuck in a system call.
	if gp == nil || gp != _m_.curg {
		return false
	}
	if _m_.locks != 0 || _m_.mallocing != 0 || _m_.throwing != 0 || _m_.preemptoff != "" || _m_.dying != 0 {
		return false
	}
	status := readgstatus(gp)
	if status&^_Gscan != _Grunning || gp.syscallsp != 0 {
		return false
	}
	if GOOS == "windows" && _m_.libcallsp != 0 {
		return false
	}
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
	return f.funcID == funcID_abort
}
