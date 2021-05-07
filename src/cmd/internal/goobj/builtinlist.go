// Code generated by mkbuiltin.go. DO NOT EDIT.

package goobj

var builtins = [...]struct {
	name string
	abi  int
}{
	{"runtime.newobject", 1},
	{"runtime.mallocgc", 1},
	{"runtime.panicdivide", 1},
	{"runtime.panicshift", 1},
	{"runtime.panicmakeslicelen", 1},
	{"runtime.panicmakeslicecap", 1},
	{"runtime.throwinit", 1},
	{"runtime.panicwrap", 1},
	{"runtime.gopanic", 1},
	{"runtime.gorecover", 1},
	{"runtime.goschedguarded", 1},
	{"runtime.goPanicIndex", 1},
	{"runtime.goPanicIndexU", 1},
	{"runtime.goPanicSliceAlen", 1},
	{"runtime.goPanicSliceAlenU", 1},
	{"runtime.goPanicSliceAcap", 1},
	{"runtime.goPanicSliceAcapU", 1},
	{"runtime.goPanicSliceB", 1},
	{"runtime.goPanicSliceBU", 1},
	{"runtime.goPanicSlice3Alen", 1},
	{"runtime.goPanicSlice3AlenU", 1},
	{"runtime.goPanicSlice3Acap", 1},
	{"runtime.goPanicSlice3AcapU", 1},
	{"runtime.goPanicSlice3B", 1},
	{"runtime.goPanicSlice3BU", 1},
	{"runtime.goPanicSlice3C", 1},
	{"runtime.goPanicSlice3CU", 1},
	{"runtime.goPanicSliceConvert", 1},
	{"runtime.printbool", 1},
	{"runtime.printfloat", 1},
	{"runtime.printint", 1},
	{"runtime.printhex", 1},
	{"runtime.printuint", 1},
	{"runtime.printcomplex", 1},
	{"runtime.printstring", 1},
	{"runtime.printpointer", 1},
	{"runtime.printuintptr", 1},
	{"runtime.printiface", 1},
	{"runtime.printeface", 1},
	{"runtime.printslice", 1},
	{"runtime.printnl", 1},
	{"runtime.printsp", 1},
	{"runtime.printlock", 1},
	{"runtime.printunlock", 1},
	{"runtime.concatstring2", 1},
	{"runtime.concatstring3", 1},
	{"runtime.concatstring4", 1},
	{"runtime.concatstring5", 1},
	{"runtime.concatstrings", 1},
	{"runtime.cmpstring", 1},
	{"runtime.intstring", 1},
	{"runtime.slicebytetostring", 1},
	{"runtime.slicebytetostringtmp", 1},
	{"runtime.slicerunetostring", 1},
	{"runtime.stringtoslicebyte", 1},
	{"runtime.stringtoslicerune", 1},
	{"runtime.slicecopy", 1},
	{"runtime.decoderune", 1},
	{"runtime.countrunes", 1},
	{"runtime.convI2I", 1},
	{"runtime.convT16", 1},
	{"runtime.convT32", 1},
	{"runtime.convT64", 1},
	{"runtime.convTstring", 1},
	{"runtime.convTslice", 1},
	{"runtime.convT2E", 1},
	{"runtime.convT2Enoptr", 1},
	{"runtime.convT2I", 1},
	{"runtime.convT2Inoptr", 1},
	{"runtime.assertE2I", 1},
	{"runtime.assertE2I2", 1},
	{"runtime.assertI2I", 1},
	{"runtime.assertI2I2", 1},
	{"runtime.panicdottypeE", 1},
	{"runtime.panicdottypeI", 1},
	{"runtime.panicnildottype", 1},
	{"runtime.ifaceeq", 1},
	{"runtime.efaceeq", 1},
	{"runtime.fastrand", 1},
	{"runtime.makemap64", 1},
	{"runtime.makemap", 1},
	{"runtime.makemap_small", 1},
	{"runtime.mapaccess1", 1},
	{"runtime.mapaccess1_fast32", 1},
	{"runtime.mapaccess1_fast64", 1},
	{"runtime.mapaccess1_faststr", 1},
	{"runtime.mapaccess1_fat", 1},
	{"runtime.mapaccess2", 1},
	{"runtime.mapaccess2_fast32", 1},
	{"runtime.mapaccess2_fast64", 1},
	{"runtime.mapaccess2_faststr", 1},
	{"runtime.mapaccess2_fat", 1},
	{"runtime.mapassign", 1},
	{"runtime.mapassign_fast32", 1},
	{"runtime.mapassign_fast32ptr", 1},
	{"runtime.mapassign_fast64", 1},
	{"runtime.mapassign_fast64ptr", 1},
	{"runtime.mapassign_faststr", 1},
	{"runtime.mapiterinit", 1},
	{"runtime.mapdelete", 1},
	{"runtime.mapdelete_fast32", 1},
	{"runtime.mapdelete_fast64", 1},
	{"runtime.mapdelete_faststr", 1},
	{"runtime.mapiternext", 1},
	{"runtime.mapclear", 1},
	{"runtime.makechan64", 1},
	{"runtime.makechan", 1},
	{"runtime.chanrecv1", 1},
	{"runtime.chanrecv2", 1},
	{"runtime.chansend1", 1},
	{"runtime.closechan", 1},
	{"runtime.writeBarrier", 0},
	{"runtime.typedmemmove", 1},
	{"runtime.typedmemclr", 1},
	{"runtime.typedslicecopy", 1},
	{"runtime.selectnbsend", 1},
	{"runtime.selectnbrecv", 1},
	{"runtime.selectsetpc", 1},
	{"runtime.selectgo", 1},
	{"runtime.block", 1},
	{"runtime.makeslice", 1},
	{"runtime.makeslice64", 1},
	{"runtime.makeslicecopy", 1},
	{"runtime.growslice", 1},
	{"runtime.unsafeslice", 1},
	{"runtime.unsafeslice64", 1},
	{"runtime.memmove", 1},
	{"runtime.memclrNoHeapPointers", 1},
	{"runtime.memclrHasPointers", 1},
	{"runtime.memequal", 1},
	{"runtime.memequal0", 1},
	{"runtime.memequal8", 1},
	{"runtime.memequal16", 1},
	{"runtime.memequal32", 1},
	{"runtime.memequal64", 1},
	{"runtime.memequal128", 1},
	{"runtime.f32equal", 1},
	{"runtime.f64equal", 1},
	{"runtime.c64equal", 1},
	{"runtime.c128equal", 1},
	{"runtime.strequal", 1},
	{"runtime.interequal", 1},
	{"runtime.nilinterequal", 1},
	{"runtime.memhash", 1},
	{"runtime.memhash0", 1},
	{"runtime.memhash8", 1},
	{"runtime.memhash16", 1},
	{"runtime.memhash32", 1},
	{"runtime.memhash64", 1},
	{"runtime.memhash128", 1},
	{"runtime.f32hash", 1},
	{"runtime.f64hash", 1},
	{"runtime.c64hash", 1},
	{"runtime.c128hash", 1},
	{"runtime.strhash", 1},
	{"runtime.interhash", 1},
	{"runtime.nilinterhash", 1},
	{"runtime.int64div", 1},
	{"runtime.uint64div", 1},
	{"runtime.int64mod", 1},
	{"runtime.uint64mod", 1},
	{"runtime.float64toint64", 1},
	{"runtime.float64touint64", 1},
	{"runtime.float64touint32", 1},
	{"runtime.int64tofloat64", 1},
	{"runtime.uint64tofloat64", 1},
	{"runtime.uint32tofloat64", 1},
	{"runtime.complex128div", 1},
	{"runtime.getcallerpc", 1},
	{"runtime.getcallersp", 1},
	{"runtime.racefuncenter", 1},
	{"runtime.racefuncexit", 1},
	{"runtime.raceread", 1},
	{"runtime.racewrite", 1},
	{"runtime.racereadrange", 1},
	{"runtime.racewriterange", 1},
	{"runtime.msanread", 1},
	{"runtime.msanwrite", 1},
	{"runtime.msanmove", 1},
	{"runtime.checkptrAlignment", 1},
	{"runtime.checkptrArithmetic", 1},
	{"runtime.libfuzzerTraceCmp1", 1},
	{"runtime.libfuzzerTraceCmp2", 1},
	{"runtime.libfuzzerTraceCmp4", 1},
	{"runtime.libfuzzerTraceCmp8", 1},
	{"runtime.libfuzzerTraceConstCmp1", 1},
	{"runtime.libfuzzerTraceConstCmp2", 1},
	{"runtime.libfuzzerTraceConstCmp4", 1},
	{"runtime.libfuzzerTraceConstCmp8", 1},
	{"runtime.libfuzzerHookStrCmp", 1},
	{"runtime.libfuzzerHookEqualFold", 1},
	{"runtime.x86HasPOPCNT", 0},
	{"runtime.x86HasSSE41", 0},
	{"runtime.x86HasFMA", 0},
	{"runtime.armHasVFPv4", 0},
	{"runtime.arm64HasATOMICS", 0},
	{"runtime.deferproc", 1},
	{"runtime.deferprocStack", 1},
	{"runtime.deferreturn", 1},
	{"runtime.newproc", 1},
	{"runtime.panicoverflow", 1},
	{"runtime.sigpanic", 1},
	{"runtime.gcWriteBarrier", 1},
	{"runtime.duffzero", 1},
	{"runtime.duffcopy", 1},
	{"runtime.morestack", 0},
	{"runtime.morestackc", 0},
	{"runtime.morestack_noctxt", 0},
	{"type:int8", 0},
	{"type:*int8", 0},
	{"type:uint8", 0},
	{"type:*uint8", 0},
	{"type:int16", 0},
	{"type:*int16", 0},
	{"type:uint16", 0},
	{"type:*uint16", 0},
	{"type:int32", 0},
	{"type:*int32", 0},
	{"type:uint32", 0},
	{"type:*uint32", 0},
	{"type:int64", 0},
	{"type:*int64", 0},
	{"type:uint64", 0},
	{"type:*uint64", 0},
	{"type:float32", 0},
	{"type:*float32", 0},
	{"type:float64", 0},
	{"type:*float64", 0},
	{"type:complex64", 0},
	{"type:*complex64", 0},
	{"type:complex128", 0},
	{"type:*complex128", 0},
	{"type:unsafe.Pointer", 0},
	{"type:*unsafe.Pointer", 0},
	{"type:uintptr", 0},
	{"type:*uintptr", 0},
	{"type:bool", 0},
	{"type:*bool", 0},
	{"type:string", 0},
	{"type:*string", 0},
	{"type:error", 0},
	{"type:*error", 0},
	{"type:func(error) string", 0},
	{"type:*func(error) string", 0},
}
