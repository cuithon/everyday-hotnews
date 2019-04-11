// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"runtime/internal/atomic"
	"runtime/internal/sys"
	"unsafe"
)

// Check to make sure we can really generate a panic. If the panic
// was generated from the runtime, or from inside malloc, then convert
// to a throw of msg.
// pc should be the program counter of the compiler-generated code that
// triggered this panic.
func panicCheck1(pc uintptr, msg string) {
	if sys.GoarchWasm == 0 && hasPrefix(funcname(findfunc(pc)), "runtime.") {
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
func panicCheck2(err string) {
	gp := getg()
	if gp != nil && gp.m != nil && gp.m.mallocing != 0 {
		throw(err)
	}
}

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

// Create a new deferred function fn with siz bytes of arguments.
// The compiler turns a defer statement into a call to this.
//go:nosplit
func deferproc(siz int32, fn *funcval) { // arguments of fn follow fn
	if getg().m.curg != getg() {
		// go code on the system stack can't defer
		throw("defer on system stack")
	}

	// the arguments of fn are in a perilous state. The stack map
	// for deferproc does not describe them. So we can't let garbage
	// collection or stack copying trigger until we've copied them out
	// to somewhere safe. The memmove below does that.
	// Until the copy completes, we can only call nosplit routines.
	sp := getcallersp()
	argp := uintptr(unsafe.Pointer(&fn)) + unsafe.Sizeof(fn)
	callerpc := getcallerpc()

	d := newdefer(siz)
	if d._panic != nil {
		throw("deferproc: d.panic != nil after newdefer")
	}
	d.fn = fn
	d.pc = callerpc
	d.sp = sp
	switch siz {
	case 0:
		// Do nothing.
	case sys.PtrSize:
		*(*uintptr)(deferArgs(d)) = *(*uintptr)(unsafe.Pointer(argp))
	default:
		memmove(deferArgs(d), unsafe.Pointer(argp), uintptr(siz))
	}

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
// The defer record must have its siz and fn fields initialized.
// All other fields can contain junk.
// The defer record must be immediately followed in memory by
// the arguments of the defer.
// Nosplit because the arguments on the stack won't be scanned
// until the defer record is spliced into the gp._defer list.
//go:nosplit
func deferprocStack(d *_defer) {
	gp := getg()
	if gp.m.curg != gp {
		// go code on the system stack can't defer
		throw("defer on system stack")
	}
	// siz and fn are already set.
	// The other fields are junk on entry to deferprocStack and
	// are initialized here.
	d.started = false
	d.heap = false
	d.sp = getcallersp()
	d.pc = getcallerpc()
	// The lines below implement:
	//   d.panic = nil
	//   d.link = gp._defer
	//   gp._defer = d
	// But without write barriers. The first two are writes to
	// the stack so they don't need a write barrier, and furthermore
	// are to uninitialized memory, so they must not use a write barrier.
	// The third write does not require a write barrier because we
	// explicitly mark all the defer structures, so we don't need to
	// keep track of pointers to them with a write barrier.
	*(*uintptr)(unsafe.Pointer(&d._panic)) = 0
	*(*uintptr)(unsafe.Pointer(&d.link)) = uintptr(unsafe.Pointer(gp._defer))
	*(*uintptr)(unsafe.Pointer(&gp._defer)) = uintptr(unsafe.Pointer(d))

	return0()
	// No code can go here - the C return register has
	// been set and must not be clobbered.
}

// Small malloc size classes >= 16 are the multiples of 16: 16, 32, 48, 64, 80, 96, 112, 128, 144, ...
// Each P holds a pool for defers with small arg sizes.
// Assign defer allocations to pools by rounding to 16, to match malloc size classes.

const (
	deferHeaderSize = unsafe.Sizeof(_defer{})
	minDeferAlloc   = (deferHeaderSize + 15) &^ 15
	minDeferArgs    = minDeferAlloc - deferHeaderSize
)

// defer size class for arg size sz
//go:nosplit
func deferclass(siz uintptr) uintptr {
	if siz <= minDeferArgs {
		return 0
	}
	return (siz - minDeferArgs + 15) / 16
}

// total size of memory block for defer with arg size sz
func totaldefersize(siz uintptr) uintptr {
	if siz <= minDeferArgs {
		return minDeferAlloc
	}
	return deferHeaderSize + siz
}

// Ensure that defer arg sizes that map to the same defer size class
// also map to the same malloc size class.
func testdefersizes() {
	var m [len(p{}.deferpool)]int32

	for i := range m {
		m[i] = -1
	}
	for i := uintptr(0); ; i++ {
		defersc := deferclass(i)
		if defersc >= uintptr(len(m)) {
			break
		}
		siz := roundupsize(totaldefersize(i))
		if m[defersc] < 0 {
			m[defersc] = int32(siz)
			continue
		}
		if m[defersc] != int32(siz) {
			print("bad defer size class: i=", i, " siz=", siz, " defersc=", defersc, "\n")
			throw("bad defer size class")
		}
	}
}

// The arguments associated with a deferred call are stored
// immediately after the _defer header in memory.
//go:nosplit
func deferArgs(d *_defer) unsafe.Pointer {
	if d.siz == 0 {
		// Avoid pointer past the defer allocation.
		return nil
	}
	return add(unsafe.Pointer(d), unsafe.Sizeof(*d))
}

var deferType *_type // type of _defer struct

func init() {
	var x interface{}
	x = (*_defer)(nil)
	deferType = (*(**ptrtype)(unsafe.Pointer(&x))).elem
}

// Allocate a Defer, usually using per-P pool.
// Each defer must be released with freedefer.
//
// This must not grow the stack because there may be a frame without
// stack map information when this is called.
//
//go:nosplit
func newdefer(siz int32) *_defer {
	var d *_defer
	sc := deferclass(uintptr(siz))
	gp := getg()
	if sc < uintptr(len(p{}.deferpool)) {
		pp := gp.m.p.ptr()
		if len(pp.deferpool[sc]) == 0 && sched.deferpool[sc] != nil {
			// Take the slow path on the system stack so
			// we don't grow newdefer's stack.
			systemstack(func() {
				lock(&sched.deferlock)
				for len(pp.deferpool[sc]) < cap(pp.deferpool[sc])/2 && sched.deferpool[sc] != nil {
					d := sched.deferpool[sc]
					sched.deferpool[sc] = d.link
					d.link = nil
					pp.deferpool[sc] = append(pp.deferpool[sc], d)
				}
				unlock(&sched.deferlock)
			})
		}
		if n := len(pp.deferpool[sc]); n > 0 {
			d = pp.deferpool[sc][n-1]
			pp.deferpool[sc][n-1] = nil
			pp.deferpool[sc] = pp.deferpool[sc][:n-1]
		}
	}
	if d == nil {
		// Allocate new defer+args.
		systemstack(func() {
			total := roundupsize(totaldefersize(uintptr(siz)))
			d = (*_defer)(mallocgc(total, deferType, true))
		})
		if debugCachedWork {
			// Duplicate the tail below so if there's a
			// crash in checkPut we can tell if d was just
			// allocated or came from the pool.
			d.siz = siz
			d.link = gp._defer
			gp._defer = d
			return d
		}
	}
	d.siz = siz
	d.heap = true
	d.link = gp._defer
	gp._defer = d
	return d
}

// Free the given defer.
// The defer cannot be used after this call.
//
// This must not grow the stack because there may be a frame without a
// stack map when this is called.
//
//go:nosplit
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
	sc := deferclass(uintptr(d.siz))
	if sc >= uintptr(len(p{}.deferpool)) {
		return
	}
	pp := getg().m.p.ptr()
	if len(pp.deferpool[sc]) == cap(pp.deferpool[sc]) {
		// Transfer half of local cache to the central cache.
		//
		// Take this slow path on the system stack so
		// we don't grow freedefer's stack.
		systemstack(func() {
			var first, last *_defer
			for len(pp.deferpool[sc]) > cap(pp.deferpool[sc])/2 {
				n := len(pp.deferpool[sc])
				d := pp.deferpool[sc][n-1]
				pp.deferpool[sc][n-1] = nil
				pp.deferpool[sc] = pp.deferpool[sc][:n-1]
				if first == nil {
					first = d
				} else {
					last.link = d
				}
				last = d
			}
			lock(&sched.deferlock)
			last.link = sched.deferpool[sc]
			sched.deferpool[sc] = first
			unlock(&sched.deferlock)
		})
	}

	// These lines used to be simply `*d = _defer{}` but that
	// started causing a nosplit stack overflow via typedmemmove.
	d.siz = 0
	d.started = false
	d.sp = 0
	d.pc = 0
	// d._panic and d.fn must be nil already.
	// If not, we would have called freedeferpanic or freedeferfn above,
	// both of which throw.
	d.link = nil

	pp.deferpool[sc] = append(pp.deferpool[sc], d)
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

// Run a deferred function if there is one.
// The compiler inserts a call to this at the end of any
// function which calls defer.
// If there is a deferred function, this will call runtime·jmpdefer,
// which will jump to the deferred function such that it appears
// to have been called by the caller of deferreturn at the point
// just before deferreturn was called. The effect is that deferreturn
// is called again and again until there are no more deferred functions.
// Cannot split the stack because we reuse the caller's frame to
// call the deferred function.

// The single argument isn't actually used - it just has its address
// taken so it can be matched against pending defers.
//go:nosplit
func deferreturn(arg0 uintptr) {
	gp := getg()
	d := gp._defer
	if d == nil {
		return
	}
	sp := getcallersp()
	if d.sp != sp {
		return
	}

	// Moving arguments around.
	//
	// Everything called after this point must be recursively
	// nosplit because the garbage collector won't know the form
	// of the arguments until the jmpdefer can flip the PC over to
	// fn.
	switch d.siz {
	case 0:
		// Do nothing.
	case sys.PtrSize:
		*(*uintptr)(unsafe.Pointer(&arg0)) = *(*uintptr)(deferArgs(d))
	default:
		memmove(unsafe.Pointer(&arg0), deferArgs(d), uintptr(d.siz))
	}
	fn := d.fn
	d.fn = nil
	gp._defer = d.link
	freedefer(d)
	jmpdefer(fn, uintptr(unsafe.Pointer(&arg0)))
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
			d.fn = nil
			gp._defer = d.link
			freedefer(d)
			continue
		}
		d.started = true
		reflectcall(nil, unsafe.Pointer(d.fn), deferArgs(d), uint32(d.siz), uint32(d.siz))
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
		print("\t")
	}
	print("panic: ")
	printany(p.arg)
	if p.recovered {
		print(" [recovered]")
	}
	print("\n")
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

	for {
		d := gp._defer
		if d == nil {
			break
		}

		// If defer was started by earlier panic or Goexit (and, since we're back here, that triggered a new panic),
		// take defer off list. The earlier panic or Goexit will not continue running.
		if d.started {
			if d._panic != nil {
				d._panic.aborted = true
			}
			d._panic = nil
			d.fn = nil
			gp._defer = d.link
			freedefer(d)
			continue
		}

		// Mark defer as started, but keep on list, so that traceback
		// can find and update the defer's argument frame if stack growth
		// or a garbage collection happens before reflectcall starts executing d.fn.
		d.started = true

		// Record the panic that is running the defer.
		// If there is a new panic during the deferred call, that panic
		// will find d in the list and will mark d._panic (this panic) aborted.
		d._panic = (*_panic)(noescape(unsafe.Pointer(&p)))

		p.argp = unsafe.Pointer(getargp(0))
		reflectcall(nil, unsafe.Pointer(d.fn), deferArgs(d), uint32(d.siz), uint32(d.siz))
		p.argp = nil

		// reflectcall did not panic. Remove d.
		if gp._defer != d {
			throw("bad defer entry in panic")
		}
		d._panic = nil
		d.fn = nil
		gp._defer = d.link

		// trigger shrinkage to test stack copy. See stack_test.go:TestStackPanic
		//GC()

		pc := d.pc
		sp := unsafe.Pointer(d.sp) // must be pointer so it gets adjusted during stack copy
		freedefer(d)
		if p.recovered {
			atomic.Xadd(&runningPanicDefers, -1)

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
func getargp(x int) uintptr {
	// x is an argument mainly so that we can return its address.
	return uintptr(noescape(unsafe.Pointer(&x)))
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
	if p != nil && !p.recovered && argp == uintptr(p.argp) {
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
	// this time returning 1.  The calling function will
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
	return pc == funcPC(abort) || ((GOARCH == "arm" || GOARCH == "arm64") && pc == funcPC(abort)+sys.PCQuantum)
}
