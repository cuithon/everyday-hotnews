// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metrics

// Description describes a runtime metric.
type Description struct {
	// Name is the full name of the metric which includes the unit.
	//
	// The format of the metric may be described by the following regular expression.
	//
	// 	^(?P<name>/[^:]+):(?P<unit>[^:*/]+(?:[*/][^:*/]+)*)$
	//
	// The format splits the name into two components, separated by a colon: a path which always
	// starts with a /, and a machine-parseable unit. The name may contain any valid Unicode
	// codepoint in between / characters, but by convention will try to stick to lowercase
	// characters and hyphens. An example of such a path might be "/memory/heap/free".
	//
	// The unit is by convention a series of lowercase English unit names (singular or plural)
	// without prefixes delimited by '*' or '/'. The unit names may contain any valid Unicode
	// codepoint that is not a delimiter.
	// Examples of units might be "seconds", "bytes", "bytes/second", "cpu-seconds",
	// "byte*cpu-seconds", and "bytes/second/second".
	//
	// For histograms, multiple units may apply. For instance, the units of the buckets and
	// the count. By convention, for histograms, the units of the count are always "samples"
	// with the type of sample evident by the metric's name, while the unit in the name
	// specifies the buckets' unit.
	//
	// A complete name might look like "/memory/heap/free:bytes".
	Name string

	// Description is an English language sentence describing the metric.
	Description string

	// Kind is the kind of value for this metric.
	//
	// The purpose of this field is to allow users to filter out metrics whose values are
	// types which their application may not understand.
	Kind ValueKind

	// Cumulative is whether or not the metric is cumulative. If a cumulative metric is just
	// a single number, then it increases monotonically. If the metric is a distribution,
	// then each bucket count increases monotonically.
	//
	// This flag thus indicates whether or not it's useful to compute a rate from this value.
	Cumulative bool
}

// The English language descriptions below must be kept in sync with the
// descriptions of each metric in doc.go.
var allDesc = []Description{
	{
		Name:        "/cgo/go-to-c-calls:calls",
		Description: "Count of calls made from Go to C by the current process.",
		Kind:        KindUint64,
		Cumulative:  true,
	},
	{
		Name: "/cpu/classes/gc/mark/assist:cpu-seconds",
		Description: "Estimated total CPU time goroutines spent performing GC tasks " +
			"to assist the GC and prevent it from falling behind the application. " +
			"This metric is an overestimate, and not directly comparable to " +
			"system CPU time measurements. Compare only with other /cpu/classes " +
			"metrics.",
		Kind:       KindFloat64,
		Cumulative: true,
	},
	{
		Name: "/cpu/classes/gc/mark/dedicated:cpu-seconds",
		Description: "Estimated total CPU time spent performing GC tasks on " +
			"processors (as defined by GOMAXPROCS) dedicated to those tasks. " +
			"This includes time spent with the world stopped due to the GC. " +
			"This metric is an overestimate, and not directly comparable to " +
			"system CPU time measurements. Compare only with other /cpu/classes " +
			"metrics.",
		Kind:       KindFloat64,
		Cumulative: true,
	},
	{
		Name: "/cpu/classes/gc/mark/idle:cpu-seconds",
		Description: "Estimated total CPU time spent performing GC tasks on " +
			"spare CPU resources that the Go scheduler could not otherwise find " +
			"a use for. This should be subtracted from the total GC CPU time to " +
			"obtain a measure of compulsory GC CPU time. " +
			"This metric is an overestimate, and not directly comparable to " +
			"system CPU time measurements. Compare only with other /cpu/classes " +
			"metrics.",
		Kind:       KindFloat64,
		Cumulative: true,
	},
	{
		Name: "/cpu/classes/gc/pause:cpu-seconds",
		Description: "Estimated total CPU time spent with the application paused by " +
			"the GC. Even if only one thread is running during the pause, this is " +
			"computed as GOMAXPROCS times the pause latency because nothing else " +
			"can be executing. This is the exact sum of samples in /gc/pause:seconds " +
			"if each sample is multiplied by GOMAXPROCS at the time it is taken. " +
			"This metric is an overestimate, and not directly comparable to " +
			"system CPU time measurements. Compare only with other /cpu/classes " +
			"metrics.",
		Kind:       KindFloat64,
		Cumulative: true,
	},
	{
		Name: "/cpu/classes/gc/total:cpu-seconds",
		Description: "Estimated total CPU time spent performing GC tasks. " +
			"This metric is an overestimate, and not directly comparable to " +
			"system CPU time measurements. Compare only with other /cpu/classes " +
			"metrics. Sum of all metrics in /cpu/classes/gc.",
		Kind:       KindFloat64,
		Cumulative: true,
	},
	{
		Name: "/cpu/classes/idle:cpu-seconds",
		Description: "Estimated total available CPU time not spent executing any Go or Go runtime code. " +
			"In other words, the part of /cpu/classes/total:cpu-seconds that was unused. " +
			"This metric is an overestimate, and not directly comparable to " +
			"system CPU time measurements. Compare only with other /cpu/classes " +
			"metrics.",
		Kind:       KindFloat64,
		Cumulative: true,
	},
	{
		Name: "/cpu/classes/scavenge/assist:cpu-seconds",
		Description: "Estimated total CPU time spent returning unused memory to the " +
			"underlying platform in response eagerly in response to memory pressure. " +
			"This metric is an overestimate, and not directly comparable to " +
			"system CPU time measurements. Compare only with other /cpu/classes " +
			"metrics.",
		Kind:       KindFloat64,
		Cumulative: true,
	},
	{
		Name: "/cpu/classes/scavenge/background:cpu-seconds",
		Description: "Estimated total CPU time spent performing background tasks " +
			"to return unused memory to the underlying platform. " +
			"This metric is an overestimate, and not directly comparable to " +
			"system CPU time measurements. Compare only with other /cpu/classes " +
			"metrics.",
		Kind:       KindFloat64,
		Cumulative: true,
	},
	{
		Name: "/cpu/classes/scavenge/total:cpu-seconds",
		Description: "Estimated total CPU time spent performing tasks that return " +
			"unused memory to the underlying platform. " +
			"This metric is an overestimate, and not directly comparable to " +
			"system CPU time measurements. Compare only with other /cpu/classes " +
			"metrics. Sum of all metrics in /cpu/classes/scavenge.",
		Kind:       KindFloat64,
		Cumulative: true,
	},
	{
		Name: "/cpu/classes/total:cpu-seconds",
		Description: "Estimated total available CPU time for user Go code " +
			"or the Go runtime, as defined by GOMAXPROCS. In other words, GOMAXPROCS " +
			"integrated over the wall-clock duration this process has been executing for. " +
			"This metric is an overestimate, and not directly comparable to " +
			"system CPU time measurements. Compare only with other /cpu/classes " +
			"metrics. Sum of all metrics in /cpu/classes.",
		Kind:       KindFloat64,
		Cumulative: true,
	},
	{
		Name: "/cpu/classes/user:cpu-seconds",
		Description: "Estimated total CPU time spent running user Go code. This may " +
			"also include some small amount of time spent in the Go runtime. " +
			"This metric is an overestimate, and not directly comparable to " +
			"system CPU time measurements. Compare only with other /cpu/classes " +
			"metrics.",
		Kind:       KindFloat64,
		Cumulative: true,
	},
	{
		Name:        "/gc/cycles/automatic:gc-cycles",
		Description: "Count of completed GC cycles generated by the Go runtime.",
		Kind:        KindUint64,
		Cumulative:  true,
	},
	{
		Name:        "/gc/cycles/forced:gc-cycles",
		Description: "Count of completed GC cycles forced by the application.",
		Kind:        KindUint64,
		Cumulative:  true,
	},
	{
		Name:        "/gc/cycles/total:gc-cycles",
		Description: "Count of all completed GC cycles.",
		Kind:        KindUint64,
		Cumulative:  true,
	},
	{
		Name: "/gc/heap/allocs-by-size:bytes",
		Description: "Distribution of heap allocations by approximate size. " +
			"Note that this does not include tiny objects as defined by " +
			"/gc/heap/tiny/allocs:objects, only tiny blocks.",
		Kind:       KindFloat64Histogram,
		Cumulative: true,
	},
	{
		Name:        "/gc/heap/allocs:bytes",
		Description: "Cumulative sum of memory allocated to the heap by the application.",
		Kind:        KindUint64,
		Cumulative:  true,
	},
	{
		Name: "/gc/heap/allocs:objects",
		Description: "Cumulative count of heap allocations triggered by the application. " +
			"Note that this does not include tiny objects as defined by " +
			"/gc/heap/tiny/allocs:objects, only tiny blocks.",
		Kind:       KindUint64,
		Cumulative: true,
	},
	{
		Name: "/gc/heap/frees-by-size:bytes",
		Description: "Distribution of freed heap allocations by approximate size. " +
			"Note that this does not include tiny objects as defined by " +
			"/gc/heap/tiny/allocs:objects, only tiny blocks.",
		Kind:       KindFloat64Histogram,
		Cumulative: true,
	},
	{
		Name:        "/gc/heap/frees:bytes",
		Description: "Cumulative sum of heap memory freed by the garbage collector.",
		Kind:        KindUint64,
		Cumulative:  true,
	},
	{
		Name: "/gc/heap/frees:objects",
		Description: "Cumulative count of heap allocations whose storage was freed " +
			"by the garbage collector. " +
			"Note that this does not include tiny objects as defined by " +
			"/gc/heap/tiny/allocs:objects, only tiny blocks.",
		Kind:       KindUint64,
		Cumulative: true,
	},
	{
		Name:        "/gc/heap/goal:bytes",
		Description: "Heap size target for the end of the GC cycle.",
		Kind:        KindUint64,
	},
	{
		Name:        "/gc/heap/objects:objects",
		Description: "Number of objects, live or unswept, occupying heap memory.",
		Kind:        KindUint64,
	},
	{
		Name: "/gc/heap/tiny/allocs:objects",
		Description: "Count of small allocations that are packed together into blocks. " +
			"These allocations are counted separately from other allocations " +
			"because each individual allocation is not tracked by the runtime, " +
			"only their block. Each block is already accounted for in " +
			"allocs-by-size and frees-by-size.",
		Kind:       KindUint64,
		Cumulative: true,
	},
	{
		Name: "/gc/limiter/last-enabled:gc-cycle",
		Description: "GC cycle the last time the GC CPU limiter was enabled. " +
			"This metric is useful for diagnosing the root cause of an out-of-memory " +
			"error, because the limiter trades memory for CPU time when the GC's CPU " +
			"time gets too high. This is most likely to occur with use of SetMemoryLimit. " +
			"The first GC cycle is cycle 1, so a value of 0 indicates that it was never enabled.",
		Kind: KindUint64,
	},
	{
		Name:        "/gc/pauses:seconds",
		Description: "Distribution individual GC-related stop-the-world pause latencies.",
		Kind:        KindFloat64Histogram,
		Cumulative:  true,
	},
	{
		Name:        "/gc/stack/starting-size:bytes",
		Description: "The stack size of new goroutines.",
		Kind:        KindUint64,
		Cumulative:  false,
	},
	{
		Name: "/godebug/non-default-behavior/execerrdot:events",
		Description: "The number of non-default behaviors executed by the os/exec package " +
			"due to a non-default GODEBUG=execerrdot=... setting.",
		Kind:       KindUint64,
		Cumulative: true,
	},
	{
		Name: "/godebug/non-default-behavior/http2client:events",
		Description: "The number of non-default behaviors executed by the net/http package " +
			"due to a non-default GODEBUG=http2client=... setting.",
		Kind:       KindUint64,
		Cumulative: true,
	},
	{
		Name: "/godebug/non-default-behavior/http2server:events",
		Description: "The number of non-default behaviors executed by the net/http package " +
			"due to a non-default GODEBUG=http2server=... setting.",
		Kind:       KindUint64,
		Cumulative: true,
	},
	{
		Name: "/godebug/non-default-behavior/installgoroot:events",
		Description: "The number of non-default behaviors executed by the go/build package " +
			"due to a non-default GODEBUG=installgoroot=... setting.",
		Kind:       KindUint64,
		Cumulative: true,
	},
	{
		Name: "/godebug/non-default-behavior/panicnil:events",
		Description: "The number of non-default behaviors executed by the runtime package " +
			"due to a non-default GODEBUG=panicnil=... setting.",
		Kind:       KindUint64,
		Cumulative: true,
	},
	{
		Name: "/godebug/non-default-behavior/randautoseed:events",
		Description: "The number of non-default behaviors executed by the math/rand package " +
			"due to a non-default GODEBUG=randautoseed=... setting.",
		Kind:       KindUint64,
		Cumulative: true,
	},
	{
		Name: "/godebug/non-default-behavior/tarinsecurepath:events",
		Description: "The number of non-default behaviors executed by the archive/tar package " +
			"due to a non-default GODEBUG=tarinsecurepath=... setting.",
		Kind:       KindUint64,
		Cumulative: true,
	},
	{
		Name: "/godebug/non-default-behavior/x509sha1:events",
		Description: "The number of non-default behaviors executed by the crypto/x509 package " +
			"due to a non-default GODEBUG=x509sha1=... setting.",
		Kind:       KindUint64,
		Cumulative: true,
	},
	{
		Name: "/godebug/non-default-behavior/x509usefallbackroots:events",
		Description: "The number of non-default behaviors executed by the crypto/x509 package " +
			"due to a non-default GODEBUG=x509usefallbackroots=... setting.",
		Kind:       KindUint64,
		Cumulative: true,
	},
	{
		Name: "/godebug/non-default-behavior/zipinsecurepath:events",
		Description: "The number of non-default behaviors executed by the archive/zip package " +
			"due to a non-default GODEBUG=zipinsecurepath=... setting.",
		Kind:       KindUint64,
		Cumulative: true,
	},
	{
		Name: "/memory/classes/heap/free:bytes",
		Description: "Memory that is completely free and eligible to be returned to the underlying system, " +
			"but has not been. This metric is the runtime's estimate of free address space that is backed by " +
			"physical memory.",
		Kind: KindUint64,
	},
	{
		Name:        "/memory/classes/heap/objects:bytes",
		Description: "Memory occupied by live objects and dead objects that have not yet been marked free by the garbage collector.",
		Kind:        KindUint64,
	},
	{
		Name: "/memory/classes/heap/released:bytes",
		Description: "Memory that is completely free and has been returned to the underlying system. This " +
			"metric is the runtime's estimate of free address space that is still mapped into the process, " +
			"but is not backed by physical memory.",
		Kind: KindUint64,
	},
	{
		Name:        "/memory/classes/heap/stacks:bytes",
		Description: "Memory allocated from the heap that is reserved for stack space, whether or not it is currently in-use.",
		Kind:        KindUint64,
	},
	{
		Name:        "/memory/classes/heap/unused:bytes",
		Description: "Memory that is reserved for heap objects but is not currently used to hold heap objects.",
		Kind:        KindUint64,
	},
	{
		Name:        "/memory/classes/metadata/mcache/free:bytes",
		Description: "Memory that is reserved for runtime mcache structures, but not in-use.",
		Kind:        KindUint64,
	},
	{
		Name:        "/memory/classes/metadata/mcache/inuse:bytes",
		Description: "Memory that is occupied by runtime mcache structures that are currently being used.",
		Kind:        KindUint64,
	},
	{
		Name:        "/memory/classes/metadata/mspan/free:bytes",
		Description: "Memory that is reserved for runtime mspan structures, but not in-use.",
		Kind:        KindUint64,
	},
	{
		Name:        "/memory/classes/metadata/mspan/inuse:bytes",
		Description: "Memory that is occupied by runtime mspan structures that are currently being used.",
		Kind:        KindUint64,
	},
	{
		Name:        "/memory/classes/metadata/other:bytes",
		Description: "Memory that is reserved for or used to hold runtime metadata.",
		Kind:        KindUint64,
	},
	{
		Name:        "/memory/classes/os-stacks:bytes",
		Description: "Stack memory allocated by the underlying operating system.",
		Kind:        KindUint64,
	},
	{
		Name:        "/memory/classes/other:bytes",
		Description: "Memory used by execution trace buffers, structures for debugging the runtime, finalizer and profiler specials, and more.",
		Kind:        KindUint64,
	},
	{
		Name:        "/memory/classes/profiling/buckets:bytes",
		Description: "Memory that is used by the stack trace hash map used for profiling.",
		Kind:        KindUint64,
	},
	{
		Name:        "/memory/classes/total:bytes",
		Description: "All memory mapped by the Go runtime into the current process as read-write. Note that this does not include memory mapped by code called via cgo or via the syscall package. Sum of all metrics in /memory/classes.",
		Kind:        KindUint64,
	},
	{
		Name:        "/sched/gomaxprocs:threads",
		Description: "The current runtime.GOMAXPROCS setting, or the number of operating system threads that can execute user-level Go code simultaneously.",
		Kind:        KindUint64,
	},
	{
		Name:        "/sched/goroutines:goroutines",
		Description: "Count of live goroutines.",
		Kind:        KindUint64,
	},
	{
		Name:        "/sched/latencies:seconds",
		Description: "Distribution of the time goroutines have spent in the scheduler in a runnable state before actually running.",
		Kind:        KindFloat64Histogram,
	},
	{
		Name:        "/sync/mutex/wait/total:seconds",
		Description: "Approximate cumulative time goroutines have spent blocked on a sync.Mutex or sync.RWMutex. This metric is useful for identifying global changes in lock contention. Collect a mutex or block profile using the runtime/pprof package for more detailed contention data.",
		Kind:        KindFloat64,
		Cumulative:  true,
	},
}

// All returns a slice of containing metric descriptions for all supported metrics.
func All() []Description {
	return allDesc
}
