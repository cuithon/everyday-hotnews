// Code generated by mklockrank.go; DO NOT EDIT.

package runtime

type lockRank int

// Constants representing the ranks of all non-leaf runtime locks, in rank order.
// Locks with lower rank must be taken before locks with higher rank,
// in addition to satisfying the partial order in lockPartialOrder.
// A few ranks allow self-cycles, which are specified in lockPartialOrder.
const (
	lockRankUnknown lockRank = iota

	lockRankSysmon
	lockRankSweepWaiters
	lockRankAssistQueue
	lockRankCpuprof
	lockRankSweep
	lockRankPollDesc
	lockRankDeadlock
	lockRankItab
	lockRankNotifyList
	lockRankRoot
	lockRankRwmutexW
	lockRankDefer
	lockRankScavenge
	lockRankForcegc
	lockRankSched
	lockRankAllg
	lockRankAllp
	lockRankTimers
	lockRankReflectOffs
	lockRankHchan
	lockRankTraceBuf
	lockRankFin
	lockRankTraceStrings
	lockRankMspanSpecial
	lockRankProfInsert
	lockRankProfBlock
	lockRankProfMemActive
	lockRankGcBitsArenas
	lockRankSpanSetSpine
	lockRankMheapSpecial
	lockRankProfMemFuture
	lockRankTrace
	lockRankTraceStackTab
	lockRankNetpollInit
	lockRankRwmutexR
	lockRankGscan
	lockRankStackpool
	lockRankStackLarge
	lockRankHchanLeaf
	lockRankSudog
	lockRankWbufSpans
	lockRankMheap
	lockRankGlobalAlloc
	lockRankPanic
)

// lockRankLeafRank is the rank of lock that does not have a declared rank,
// and hence is a leaf lock.
const lockRankLeafRank lockRank = 1000

// lockNames gives the names associated with each of the above ranks.
var lockNames = []string{
	lockRankSysmon:        "sysmon",
	lockRankSweepWaiters:  "sweepWaiters",
	lockRankAssistQueue:   "assistQueue",
	lockRankCpuprof:       "cpuprof",
	lockRankSweep:         "sweep",
	lockRankPollDesc:      "pollDesc",
	lockRankDeadlock:      "deadlock",
	lockRankItab:          "itab",
	lockRankNotifyList:    "notifyList",
	lockRankRoot:          "root",
	lockRankRwmutexW:      "rwmutexW",
	lockRankDefer:         "defer",
	lockRankScavenge:      "scavenge",
	lockRankForcegc:       "forcegc",
	lockRankSched:         "sched",
	lockRankAllg:          "allg",
	lockRankAllp:          "allp",
	lockRankTimers:        "timers",
	lockRankReflectOffs:   "reflectOffs",
	lockRankHchan:         "hchan",
	lockRankTraceBuf:      "traceBuf",
	lockRankFin:           "fin",
	lockRankTraceStrings:  "traceStrings",
	lockRankMspanSpecial:  "mspanSpecial",
	lockRankProfInsert:    "profInsert",
	lockRankProfBlock:     "profBlock",
	lockRankProfMemActive: "profMemActive",
	lockRankGcBitsArenas:  "gcBitsArenas",
	lockRankSpanSetSpine:  "spanSetSpine",
	lockRankMheapSpecial:  "mheapSpecial",
	lockRankProfMemFuture: "profMemFuture",
	lockRankTrace:         "trace",
	lockRankTraceStackTab: "traceStackTab",
	lockRankNetpollInit:   "netpollInit",
	lockRankRwmutexR:      "rwmutexR",
	lockRankGscan:         "gscan",
	lockRankStackpool:     "stackpool",
	lockRankStackLarge:    "stackLarge",
	lockRankHchanLeaf:     "hchanLeaf",
	lockRankSudog:         "sudog",
	lockRankWbufSpans:     "wbufSpans",
	lockRankMheap:         "mheap",
	lockRankGlobalAlloc:   "globalAlloc",
	lockRankPanic:         "panic",
}

func (rank lockRank) String() string {
	if rank == 0 {
		return "UNKNOWN"
	}
	if rank == lockRankLeafRank {
		return "LEAF"
	}
	if rank < 0 || int(rank) >= len(lockNames) {
		return "BAD RANK"
	}
	return lockNames[rank]
}

// lockPartialOrder is the transitive closure of the lock rank graph.
// An entry for rank X lists all of the ranks that can already be held
// when rank X is acquired.
//
// Lock ranks that allow self-cycles list themselves.
var lockPartialOrder [][]lockRank = [][]lockRank{
	lockRankSysmon:        {},
	lockRankSweepWaiters:  {},
	lockRankAssistQueue:   {},
	lockRankCpuprof:       {},
	lockRankSweep:         {},
	lockRankPollDesc:      {},
	lockRankDeadlock:      {},
	lockRankItab:          {},
	lockRankNotifyList:    {},
	lockRankRoot:          {},
	lockRankRwmutexW:      {},
	lockRankDefer:         {},
	lockRankScavenge:      {lockRankSysmon},
	lockRankForcegc:       {lockRankSysmon},
	lockRankSched:         {lockRankSysmon, lockRankSweepWaiters, lockRankAssistQueue, lockRankCpuprof, lockRankSweep, lockRankPollDesc, lockRankScavenge, lockRankForcegc},
	lockRankAllg:          {lockRankSysmon, lockRankSweepWaiters, lockRankAssistQueue, lockRankCpuprof, lockRankSweep, lockRankPollDesc, lockRankScavenge, lockRankForcegc, lockRankSched},
	lockRankAllp:          {lockRankSysmon, lockRankSweepWaiters, lockRankAssistQueue, lockRankCpuprof, lockRankSweep, lockRankPollDesc, lockRankScavenge, lockRankForcegc, lockRankSched},
	lockRankTimers:        {lockRankSysmon, lockRankSweepWaiters, lockRankAssistQueue, lockRankCpuprof, lockRankSweep, lockRankPollDesc, lockRankScavenge, lockRankForcegc, lockRankSched, lockRankAllp, lockRankTimers},
	lockRankReflectOffs:   {lockRankItab},
	lockRankHchan:         {lockRankSysmon, lockRankSweep, lockRankScavenge, lockRankHchan},
	lockRankTraceBuf:      {lockRankSysmon, lockRankScavenge},
	lockRankFin:           {lockRankSysmon, lockRankSweepWaiters, lockRankAssistQueue, lockRankCpuprof, lockRankSweep, lockRankPollDesc, lockRankItab, lockRankScavenge, lockRankForcegc, lockRankSched, lockRankAllg, lockRankAllp, lockRankTimers, lockRankReflectOffs, lockRankHchan, lockRankTraceBuf},
	lockRankTraceStrings:  {lockRankSysmon, lockRankScavenge, lockRankTraceBuf},
	lockRankMspanSpecial:  {lockRankSysmon, lockRankSweepWaiters, lockRankAssistQueue, lockRankCpuprof, lockRankSweep, lockRankPollDesc, lockRankItab, lockRankNotifyList, lockRankScavenge, lockRankForcegc, lockRankSched, lockRankAllg, lockRankAllp, lockRankTimers, lockRankReflectOffs, lockRankHchan, lockRankTraceBuf, lockRankTraceStrings},
	lockRankProfInsert:    {lockRankSysmon, lockRankSweepWaiters, lockRankAssistQueue, lockRankCpuprof, lockRankSweep, lockRankPollDesc, lockRankItab, lockRankNotifyList, lockRankScavenge, lockRankForcegc, lockRankSched, lockRankAllg, lockRankAllp, lockRankTimers, lockRankReflectOffs, lockRankHchan, lockRankTraceBuf, lockRankTraceStrings},
	lockRankProfBlock:     {lockRankSysmon, lockRankSweepWaiters, lockRankAssistQueue, lockRankCpuprof, lockRankSweep, lockRankPollDesc, lockRankItab, lockRankNotifyList, lockRankScavenge, lockRankForcegc, lockRankSched, lockRankAllg, lockRankAllp, lockRankTimers, lockRankReflectOffs, lockRankHchan, lockRankTraceBuf, lockRankTraceStrings},
	lockRankProfMemActive: {lockRankSysmon, lockRankSweepWaiters, lockRankAssistQueue, lockRankCpuprof, lockRankSweep, lockRankPollDesc, lockRankItab, lockRankNotifyList, lockRankScavenge, lockRankForcegc, lockRankSched, lockRankAllg, lockRankAllp, lockRankTimers, lockRankReflectOffs, lockRankHchan, lockRankTraceBuf, lockRankTraceStrings},
	lockRankGcBitsArenas:  {lockRankSysmon, lockRankSweepWaiters, lockRankAssistQueue, lockRankCpuprof, lockRankSweep, lockRankPollDesc, lockRankItab, lockRankNotifyList, lockRankScavenge, lockRankForcegc, lockRankSched, lockRankAllg, lockRankAllp, lockRankTimers, lockRankReflectOffs, lockRankHchan, lockRankTraceBuf, lockRankTraceStrings},
	lockRankSpanSetSpine:  {lockRankSysmon, lockRankSweepWaiters, lockRankAssistQueue, lockRankCpuprof, lockRankSweep, lockRankPollDesc, lockRankItab, lockRankNotifyList, lockRankScavenge, lockRankForcegc, lockRankSched, lockRankAllg, lockRankAllp, lockRankTimers, lockRankReflectOffs, lockRankHchan, lockRankTraceBuf, lockRankTraceStrings},
	lockRankMheapSpecial:  {lockRankSysmon, lockRankSweepWaiters, lockRankAssistQueue, lockRankCpuprof, lockRankSweep, lockRankPollDesc, lockRankItab, lockRankNotifyList, lockRankScavenge, lockRankForcegc, lockRankSched, lockRankAllg, lockRankAllp, lockRankTimers, lockRankReflectOffs, lockRankHchan, lockRankTraceBuf, lockRankTraceStrings},
	lockRankProfMemFuture: {lockRankSysmon, lockRankSweepWaiters, lockRankAssistQueue, lockRankCpuprof, lockRankSweep, lockRankPollDesc, lockRankItab, lockRankNotifyList, lockRankScavenge, lockRankForcegc, lockRankSched, lockRankAllg, lockRankAllp, lockRankTimers, lockRankReflectOffs, lockRankHchan, lockRankTraceBuf, lockRankTraceStrings, lockRankProfMemActive},
	lockRankTrace:         {lockRankSysmon, lockRankSweepWaiters, lockRankAssistQueue, lockRankCpuprof, lockRankSweep, lockRankPollDesc, lockRankItab, lockRankNotifyList, lockRankRoot, lockRankScavenge, lockRankForcegc, lockRankSched, lockRankAllg, lockRankAllp, lockRankTimers, lockRankReflectOffs, lockRankHchan, lockRankTraceBuf, lockRankFin, lockRankTraceStrings},
	lockRankTraceStackTab: {lockRankSysmon, lockRankSweepWaiters, lockRankAssistQueue, lockRankCpuprof, lockRankSweep, lockRankPollDesc, lockRankItab, lockRankNotifyList, lockRankRoot, lockRankScavenge, lockRankForcegc, lockRankSched, lockRankAllg, lockRankAllp, lockRankTimers, lockRankReflectOffs, lockRankHchan, lockRankTraceBuf, lockRankFin, lockRankTraceStrings, lockRankTrace},
	lockRankNetpollInit:   {lockRankSysmon, lockRankSweepWaiters, lockRankAssistQueue, lockRankCpuprof, lockRankSweep, lockRankPollDesc, lockRankScavenge, lockRankForcegc, lockRankSched, lockRankAllp, lockRankTimers},
	lockRankRwmutexR:      {lockRankSysmon, lockRankRwmutexW},
	lockRankGscan:         {lockRankSysmon, lockRankSweepWaiters, lockRankAssistQueue, lockRankCpuprof, lockRankSweep, lockRankPollDesc, lockRankItab, lockRankNotifyList, lockRankRoot, lockRankScavenge, lockRankForcegc, lockRankSched, lockRankAllg, lockRankAllp, lockRankTimers, lockRankReflectOffs, lockRankHchan, lockRankTraceBuf, lockRankFin, lockRankTraceStrings, lockRankProfInsert, lockRankProfBlock, lockRankProfMemActive, lockRankGcBitsArenas, lockRankSpanSetSpine, lockRankProfMemFuture, lockRankTrace, lockRankTraceStackTab, lockRankNetpollInit},
	lockRankStackpool:     {lockRankSysmon, lockRankSweepWaiters, lockRankAssistQueue, lockRankCpuprof, lockRankSweep, lockRankPollDesc, lockRankItab, lockRankNotifyList, lockRankRoot, lockRankRwmutexW, lockRankScavenge, lockRankForcegc, lockRankSched, lockRankAllg, lockRankAllp, lockRankTimers, lockRankReflectOffs, lockRankHchan, lockRankTraceBuf, lockRankFin, lockRankTraceStrings, lockRankProfInsert, lockRankProfBlock, lockRankProfMemActive, lockRankGcBitsArenas, lockRankSpanSetSpine, lockRankProfMemFuture, lockRankTrace, lockRankTraceStackTab, lockRankNetpollInit, lockRankRwmutexR, lockRankGscan},
	lockRankStackLarge:    {lockRankSysmon, lockRankSweepWaiters, lockRankAssistQueue, lockRankCpuprof, lockRankSweep, lockRankPollDesc, lockRankItab, lockRankNotifyList, lockRankRoot, lockRankScavenge, lockRankForcegc, lockRankSched, lockRankAllg, lockRankAllp, lockRankTimers, lockRankReflectOffs, lockRankHchan, lockRankTraceBuf, lockRankFin, lockRankTraceStrings, lockRankProfInsert, lockRankProfBlock, lockRankProfMemActive, lockRankGcBitsArenas, lockRankSpanSetSpine, lockRankProfMemFuture, lockRankTrace, lockRankTraceStackTab, lockRankNetpollInit, lockRankGscan},
	lockRankHchanLeaf:     {lockRankSysmon, lockRankSweepWaiters, lockRankAssistQueue, lockRankCpuprof, lockRankSweep, lockRankPollDesc, lockRankItab, lockRankNotifyList, lockRankRoot, lockRankScavenge, lockRankForcegc, lockRankSched, lockRankAllg, lockRankAllp, lockRankTimers, lockRankReflectOffs, lockRankHchan, lockRankTraceBuf, lockRankFin, lockRankTraceStrings, lockRankProfInsert, lockRankProfBlock, lockRankProfMemActive, lockRankGcBitsArenas, lockRankSpanSetSpine, lockRankProfMemFuture, lockRankTrace, lockRankTraceStackTab, lockRankNetpollInit, lockRankGscan, lockRankHchanLeaf},
	lockRankSudog:         {lockRankSysmon, lockRankSweep, lockRankNotifyList, lockRankScavenge, lockRankHchan},
	lockRankWbufSpans:     {lockRankSysmon, lockRankSweepWaiters, lockRankAssistQueue, lockRankCpuprof, lockRankSweep, lockRankPollDesc, lockRankItab, lockRankNotifyList, lockRankRoot, lockRankDefer, lockRankScavenge, lockRankForcegc, lockRankSched, lockRankAllg, lockRankAllp, lockRankTimers, lockRankReflectOffs, lockRankHchan, lockRankTraceBuf, lockRankFin, lockRankTraceStrings, lockRankMspanSpecial, lockRankProfInsert, lockRankProfBlock, lockRankProfMemActive, lockRankGcBitsArenas, lockRankSpanSetSpine, lockRankProfMemFuture, lockRankTrace, lockRankTraceStackTab, lockRankNetpollInit, lockRankGscan, lockRankSudog},
	lockRankMheap:         {lockRankSysmon, lockRankSweepWaiters, lockRankAssistQueue, lockRankCpuprof, lockRankSweep, lockRankPollDesc, lockRankItab, lockRankNotifyList, lockRankRoot, lockRankRwmutexW, lockRankDefer, lockRankScavenge, lockRankForcegc, lockRankSched, lockRankAllg, lockRankAllp, lockRankTimers, lockRankReflectOffs, lockRankHchan, lockRankTraceBuf, lockRankFin, lockRankTraceStrings, lockRankMspanSpecial, lockRankProfInsert, lockRankProfBlock, lockRankProfMemActive, lockRankGcBitsArenas, lockRankSpanSetSpine, lockRankProfMemFuture, lockRankTrace, lockRankTraceStackTab, lockRankNetpollInit, lockRankRwmutexR, lockRankGscan, lockRankStackpool, lockRankStackLarge, lockRankSudog, lockRankWbufSpans},
	lockRankGlobalAlloc:   {lockRankSysmon, lockRankSweepWaiters, lockRankAssistQueue, lockRankCpuprof, lockRankSweep, lockRankPollDesc, lockRankItab, lockRankNotifyList, lockRankRoot, lockRankRwmutexW, lockRankDefer, lockRankScavenge, lockRankForcegc, lockRankSched, lockRankAllg, lockRankAllp, lockRankTimers, lockRankReflectOffs, lockRankHchan, lockRankTraceBuf, lockRankFin, lockRankTraceStrings, lockRankMspanSpecial, lockRankProfInsert, lockRankProfBlock, lockRankProfMemActive, lockRankGcBitsArenas, lockRankSpanSetSpine, lockRankMheapSpecial, lockRankProfMemFuture, lockRankTrace, lockRankTraceStackTab, lockRankNetpollInit, lockRankRwmutexR, lockRankGscan, lockRankStackpool, lockRankStackLarge, lockRankSudog, lockRankWbufSpans, lockRankMheap},
	lockRankPanic:         {lockRankDeadlock},
}
