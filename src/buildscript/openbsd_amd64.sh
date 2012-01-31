#!/usr/bin/env bash
# AUTO-GENERATED by buildscript.sh; DO NOT EDIT.
# This script builds the go command (written in Go),
# and then the go command can build the rest of the tree.

export GOOS=openbsd
export GOARCH=amd64
export WORK=$(mktemp -d -t go-build.XXXXXX)
trap "rm -rf $WORK" EXIT SIGINT SIGTERM
set -e



#
# runtime
#

mkdir -p "$WORK"/runtime/_obj/
cd "$GOROOT"/src/pkg/runtime
"$GOROOT"/bin/go-tool/6g -o "$WORK"/runtime/_obj/_go_.6 -p runtime -+ -I "$WORK" ./debug.go ./error.go ./extern.go ./mem.go ./sig.go ./softfloat64.go ./type.go ./zgoarch_amd64.go ./zgoos_openbsd.go ./zruntime_defs_openbsd_amd64.go ./zversion.go
cp "$GOROOT"/src/pkg/runtime/arch_amd64.h "$WORK"/runtime/_obj/arch_GOARCH.h
cp "$GOROOT"/src/pkg/runtime/defs_openbsd_amd64.h "$WORK"/runtime/_obj/defs_GOOS_GOARCH.h
cp "$GOROOT"/src/pkg/runtime/os_openbsd.h "$WORK"/runtime/_obj/os_GOOS.h
cp "$GOROOT"/src/pkg/runtime/signals_openbsd.h "$WORK"/runtime/_obj/signals_GOOS.h
cp "$GOROOT"/src/pkg/runtime/zasm_openbsd_amd64.h "$WORK"/runtime/_obj/zasm_GOOS_GOARCH.h
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/alg.6 -DGOOS_openbsd -DGOARCH_amd64 ./alg.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/atomic_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./atomic_amd64.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/cgocall.6 -DGOOS_openbsd -DGOARCH_amd64 ./cgocall.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/chan.6 -DGOOS_openbsd -DGOARCH_amd64 ./chan.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/closure_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./closure_amd64.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/complex.6 -DGOOS_openbsd -DGOARCH_amd64 ./complex.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/cpuprof.6 -DGOOS_openbsd -DGOARCH_amd64 ./cpuprof.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/float.6 -DGOOS_openbsd -DGOARCH_amd64 ./float.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/hashmap.6 -DGOOS_openbsd -DGOARCH_amd64 ./hashmap.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/iface.6 -DGOOS_openbsd -DGOARCH_amd64 ./iface.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/lock_sema.6 -DGOOS_openbsd -DGOARCH_amd64 ./lock_sema.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/mcache.6 -DGOOS_openbsd -DGOARCH_amd64 ./mcache.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/mcentral.6 -DGOOS_openbsd -DGOARCH_amd64 ./mcentral.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/mem_openbsd.6 -DGOOS_openbsd -DGOARCH_amd64 ./mem_openbsd.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/mfinal.6 -DGOOS_openbsd -DGOARCH_amd64 ./mfinal.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/mfixalloc.6 -DGOOS_openbsd -DGOARCH_amd64 ./mfixalloc.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/mgc0.6 -DGOOS_openbsd -DGOARCH_amd64 ./mgc0.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/mheap.6 -DGOOS_openbsd -DGOARCH_amd64 ./mheap.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/msize.6 -DGOOS_openbsd -DGOARCH_amd64 ./msize.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/print.6 -DGOOS_openbsd -DGOARCH_amd64 ./print.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/proc.6 -DGOOS_openbsd -DGOARCH_amd64 ./proc.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/rune.6 -DGOOS_openbsd -DGOARCH_amd64 ./rune.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/runtime.6 -DGOOS_openbsd -DGOARCH_amd64 ./runtime.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/signal_openbsd_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./signal_openbsd_amd64.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/slice.6 -DGOOS_openbsd -DGOARCH_amd64 ./slice.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/symtab.6 -DGOOS_openbsd -DGOARCH_amd64 ./symtab.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/thread_openbsd.6 -DGOOS_openbsd -DGOARCH_amd64 ./thread_openbsd.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/traceback_x86.6 -DGOOS_openbsd -DGOARCH_amd64 ./traceback_x86.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/zmalloc_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./zmalloc_amd64.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/zmprof_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./zmprof_amd64.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/zruntime1_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./zruntime1_amd64.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/zsema_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./zsema_amd64.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/zsigqueue_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./zsigqueue_amd64.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/zstring_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./zstring_amd64.c
"$GOROOT"/bin/go-tool/6c -FVw -I "$WORK"/runtime/_obj/ -I "$GOROOT"/pkg/openbsd_amd64 -o "$WORK"/runtime/_obj/ztime_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./ztime_amd64.c
"$GOROOT"/bin/go-tool/6a -I "$WORK"/runtime/_obj/ -o "$WORK"/runtime/_obj/asm_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./asm_amd64.s
"$GOROOT"/bin/go-tool/6a -I "$WORK"/runtime/_obj/ -o "$WORK"/runtime/_obj/memmove_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./memmove_amd64.s
"$GOROOT"/bin/go-tool/6a -I "$WORK"/runtime/_obj/ -o "$WORK"/runtime/_obj/rt0_openbsd_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./rt0_openbsd_amd64.s
"$GOROOT"/bin/go-tool/6a -I "$WORK"/runtime/_obj/ -o "$WORK"/runtime/_obj/sys_openbsd_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./sys_openbsd_amd64.s
"$GOROOT"/bin/go-tool/pack grc "$WORK"/runtime.a "$WORK"/runtime/_obj/_go_.6 "$WORK"/runtime/_obj/alg.6 "$WORK"/runtime/_obj/atomic_amd64.6 "$WORK"/runtime/_obj/cgocall.6 "$WORK"/runtime/_obj/chan.6 "$WORK"/runtime/_obj/closure_amd64.6 "$WORK"/runtime/_obj/complex.6 "$WORK"/runtime/_obj/cpuprof.6 "$WORK"/runtime/_obj/float.6 "$WORK"/runtime/_obj/hashmap.6 "$WORK"/runtime/_obj/iface.6 "$WORK"/runtime/_obj/lock_sema.6 "$WORK"/runtime/_obj/mcache.6 "$WORK"/runtime/_obj/mcentral.6 "$WORK"/runtime/_obj/mem_openbsd.6 "$WORK"/runtime/_obj/mfinal.6 "$WORK"/runtime/_obj/mfixalloc.6 "$WORK"/runtime/_obj/mgc0.6 "$WORK"/runtime/_obj/mheap.6 "$WORK"/runtime/_obj/msize.6 "$WORK"/runtime/_obj/print.6 "$WORK"/runtime/_obj/proc.6 "$WORK"/runtime/_obj/rune.6 "$WORK"/runtime/_obj/runtime.6 "$WORK"/runtime/_obj/signal_openbsd_amd64.6 "$WORK"/runtime/_obj/slice.6 "$WORK"/runtime/_obj/symtab.6 "$WORK"/runtime/_obj/thread_openbsd.6 "$WORK"/runtime/_obj/traceback_x86.6 "$WORK"/runtime/_obj/zmalloc_amd64.6 "$WORK"/runtime/_obj/zmprof_amd64.6 "$WORK"/runtime/_obj/zruntime1_amd64.6 "$WORK"/runtime/_obj/zsema_amd64.6 "$WORK"/runtime/_obj/zsigqueue_amd64.6 "$WORK"/runtime/_obj/zstring_amd64.6 "$WORK"/runtime/_obj/ztime_amd64.6 "$WORK"/runtime/_obj/asm_amd64.6 "$WORK"/runtime/_obj/memmove_amd64.6 "$WORK"/runtime/_obj/rt0_openbsd_amd64.6 "$WORK"/runtime/_obj/sys_openbsd_amd64.6
mkdir -p "$GOROOT"/pkg/openbsd_amd64/
cp "$WORK"/runtime.a "$GOROOT"/pkg/openbsd_amd64/runtime.a

#
# errors
#

mkdir -p "$WORK"/errors/_obj/
cd "$GOROOT"/src/pkg/errors
"$GOROOT"/bin/go-tool/6g -o "$WORK"/errors/_obj/_go_.6 -p errors -I "$WORK" ./errors.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/errors.a "$WORK"/errors/_obj/_go_.6
cp "$WORK"/errors.a "$GOROOT"/pkg/openbsd_amd64/errors.a

#
# sync/atomic
#

mkdir -p "$WORK"/sync/atomic/_obj/
cd "$GOROOT"/src/pkg/sync/atomic
"$GOROOT"/bin/go-tool/6g -o "$WORK"/sync/atomic/_obj/_go_.6 -p sync/atomic -I "$WORK" ./doc.go
"$GOROOT"/bin/go-tool/6a -I "$WORK"/sync/atomic/_obj/ -o "$WORK"/sync/atomic/_obj/asm_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./asm_amd64.s
"$GOROOT"/bin/go-tool/pack grc "$WORK"/sync/atomic.a "$WORK"/sync/atomic/_obj/_go_.6 "$WORK"/sync/atomic/_obj/asm_amd64.6
mkdir -p "$GOROOT"/pkg/openbsd_amd64/sync/
cp "$WORK"/sync/atomic.a "$GOROOT"/pkg/openbsd_amd64/sync/atomic.a

#
# sync
#

mkdir -p "$WORK"/sync/_obj/
cd "$GOROOT"/src/pkg/sync
"$GOROOT"/bin/go-tool/6g -o "$WORK"/sync/_obj/_go_.6 -p sync -I "$WORK" ./cond.go ./mutex.go ./once.go ./rwmutex.go ./waitgroup.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/sync.a "$WORK"/sync/_obj/_go_.6
cp "$WORK"/sync.a "$GOROOT"/pkg/openbsd_amd64/sync.a

#
# io
#

mkdir -p "$WORK"/io/_obj/
cd "$GOROOT"/src/pkg/io
"$GOROOT"/bin/go-tool/6g -o "$WORK"/io/_obj/_go_.6 -p io -I "$WORK" ./io.go ./multi.go ./pipe.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/io.a "$WORK"/io/_obj/_go_.6
cp "$WORK"/io.a "$GOROOT"/pkg/openbsd_amd64/io.a

#
# unicode
#

mkdir -p "$WORK"/unicode/_obj/
cd "$GOROOT"/src/pkg/unicode
"$GOROOT"/bin/go-tool/6g -o "$WORK"/unicode/_obj/_go_.6 -p unicode -I "$WORK" ./casetables.go ./digit.go ./graphic.go ./letter.go ./tables.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/unicode.a "$WORK"/unicode/_obj/_go_.6
cp "$WORK"/unicode.a "$GOROOT"/pkg/openbsd_amd64/unicode.a

#
# unicode/utf8
#

mkdir -p "$WORK"/unicode/utf8/_obj/
cd "$GOROOT"/src/pkg/unicode/utf8
"$GOROOT"/bin/go-tool/6g -o "$WORK"/unicode/utf8/_obj/_go_.6 -p unicode/utf8 -I "$WORK" ./utf8.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/unicode/utf8.a "$WORK"/unicode/utf8/_obj/_go_.6
mkdir -p "$GOROOT"/pkg/openbsd_amd64/unicode/
cp "$WORK"/unicode/utf8.a "$GOROOT"/pkg/openbsd_amd64/unicode/utf8.a

#
# bytes
#

mkdir -p "$WORK"/bytes/_obj/
cd "$GOROOT"/src/pkg/bytes
"$GOROOT"/bin/go-tool/6g -o "$WORK"/bytes/_obj/_go_.6 -p bytes -I "$WORK" ./buffer.go ./bytes.go ./bytes_decl.go
"$GOROOT"/bin/go-tool/6a -I "$WORK"/bytes/_obj/ -o "$WORK"/bytes/_obj/asm_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./asm_amd64.s
"$GOROOT"/bin/go-tool/pack grc "$WORK"/bytes.a "$WORK"/bytes/_obj/_go_.6 "$WORK"/bytes/_obj/asm_amd64.6
cp "$WORK"/bytes.a "$GOROOT"/pkg/openbsd_amd64/bytes.a

#
# math
#

mkdir -p "$WORK"/math/_obj/
cd "$GOROOT"/src/pkg/math
"$GOROOT"/bin/go-tool/6g -o "$WORK"/math/_obj/_go_.6 -p math -I "$WORK" ./abs.go ./acosh.go ./asin.go ./asinh.go ./atan.go ./atan2.go ./atanh.go ./bits.go ./cbrt.go ./const.go ./copysign.go ./dim.go ./erf.go ./exp.go ./expm1.go ./floor.go ./frexp.go ./gamma.go ./hypot.go ./j0.go ./j1.go ./jn.go ./ldexp.go ./lgamma.go ./log.go ./log10.go ./log1p.go ./logb.go ./mod.go ./modf.go ./nextafter.go ./pow.go ./pow10.go ./remainder.go ./signbit.go ./sin.go ./sincos.go ./sinh.go ./sqrt.go ./tan.go ./tanh.go ./unsafe.go
"$GOROOT"/bin/go-tool/6a -I "$WORK"/math/_obj/ -o "$WORK"/math/_obj/abs_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./abs_amd64.s
"$GOROOT"/bin/go-tool/6a -I "$WORK"/math/_obj/ -o "$WORK"/math/_obj/asin_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./asin_amd64.s
"$GOROOT"/bin/go-tool/6a -I "$WORK"/math/_obj/ -o "$WORK"/math/_obj/atan2_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./atan2_amd64.s
"$GOROOT"/bin/go-tool/6a -I "$WORK"/math/_obj/ -o "$WORK"/math/_obj/atan_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./atan_amd64.s
"$GOROOT"/bin/go-tool/6a -I "$WORK"/math/_obj/ -o "$WORK"/math/_obj/dim_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./dim_amd64.s
"$GOROOT"/bin/go-tool/6a -I "$WORK"/math/_obj/ -o "$WORK"/math/_obj/exp2_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./exp2_amd64.s
"$GOROOT"/bin/go-tool/6a -I "$WORK"/math/_obj/ -o "$WORK"/math/_obj/exp_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./exp_amd64.s
"$GOROOT"/bin/go-tool/6a -I "$WORK"/math/_obj/ -o "$WORK"/math/_obj/expm1_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./expm1_amd64.s
"$GOROOT"/bin/go-tool/6a -I "$WORK"/math/_obj/ -o "$WORK"/math/_obj/floor_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./floor_amd64.s
"$GOROOT"/bin/go-tool/6a -I "$WORK"/math/_obj/ -o "$WORK"/math/_obj/fltasm_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./fltasm_amd64.s
"$GOROOT"/bin/go-tool/6a -I "$WORK"/math/_obj/ -o "$WORK"/math/_obj/frexp_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./frexp_amd64.s
"$GOROOT"/bin/go-tool/6a -I "$WORK"/math/_obj/ -o "$WORK"/math/_obj/hypot_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./hypot_amd64.s
"$GOROOT"/bin/go-tool/6a -I "$WORK"/math/_obj/ -o "$WORK"/math/_obj/ldexp_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./ldexp_amd64.s
"$GOROOT"/bin/go-tool/6a -I "$WORK"/math/_obj/ -o "$WORK"/math/_obj/log10_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./log10_amd64.s
"$GOROOT"/bin/go-tool/6a -I "$WORK"/math/_obj/ -o "$WORK"/math/_obj/log1p_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./log1p_amd64.s
"$GOROOT"/bin/go-tool/6a -I "$WORK"/math/_obj/ -o "$WORK"/math/_obj/log_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./log_amd64.s
"$GOROOT"/bin/go-tool/6a -I "$WORK"/math/_obj/ -o "$WORK"/math/_obj/mod_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./mod_amd64.s
"$GOROOT"/bin/go-tool/6a -I "$WORK"/math/_obj/ -o "$WORK"/math/_obj/modf_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./modf_amd64.s
"$GOROOT"/bin/go-tool/6a -I "$WORK"/math/_obj/ -o "$WORK"/math/_obj/remainder_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./remainder_amd64.s
"$GOROOT"/bin/go-tool/6a -I "$WORK"/math/_obj/ -o "$WORK"/math/_obj/sin_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./sin_amd64.s
"$GOROOT"/bin/go-tool/6a -I "$WORK"/math/_obj/ -o "$WORK"/math/_obj/sincos_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./sincos_amd64.s
"$GOROOT"/bin/go-tool/6a -I "$WORK"/math/_obj/ -o "$WORK"/math/_obj/sqrt_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./sqrt_amd64.s
"$GOROOT"/bin/go-tool/6a -I "$WORK"/math/_obj/ -o "$WORK"/math/_obj/tan_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./tan_amd64.s
"$GOROOT"/bin/go-tool/pack grc "$WORK"/math.a "$WORK"/math/_obj/_go_.6 "$WORK"/math/_obj/abs_amd64.6 "$WORK"/math/_obj/asin_amd64.6 "$WORK"/math/_obj/atan2_amd64.6 "$WORK"/math/_obj/atan_amd64.6 "$WORK"/math/_obj/dim_amd64.6 "$WORK"/math/_obj/exp2_amd64.6 "$WORK"/math/_obj/exp_amd64.6 "$WORK"/math/_obj/expm1_amd64.6 "$WORK"/math/_obj/floor_amd64.6 "$WORK"/math/_obj/fltasm_amd64.6 "$WORK"/math/_obj/frexp_amd64.6 "$WORK"/math/_obj/hypot_amd64.6 "$WORK"/math/_obj/ldexp_amd64.6 "$WORK"/math/_obj/log10_amd64.6 "$WORK"/math/_obj/log1p_amd64.6 "$WORK"/math/_obj/log_amd64.6 "$WORK"/math/_obj/mod_amd64.6 "$WORK"/math/_obj/modf_amd64.6 "$WORK"/math/_obj/remainder_amd64.6 "$WORK"/math/_obj/sin_amd64.6 "$WORK"/math/_obj/sincos_amd64.6 "$WORK"/math/_obj/sqrt_amd64.6 "$WORK"/math/_obj/tan_amd64.6
cp "$WORK"/math.a "$GOROOT"/pkg/openbsd_amd64/math.a

#
# sort
#

mkdir -p "$WORK"/sort/_obj/
cd "$GOROOT"/src/pkg/sort
"$GOROOT"/bin/go-tool/6g -o "$WORK"/sort/_obj/_go_.6 -p sort -I "$WORK" ./search.go ./sort.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/sort.a "$WORK"/sort/_obj/_go_.6
cp "$WORK"/sort.a "$GOROOT"/pkg/openbsd_amd64/sort.a

#
# container/heap
#

mkdir -p "$WORK"/container/heap/_obj/
cd "$GOROOT"/src/pkg/container/heap
"$GOROOT"/bin/go-tool/6g -o "$WORK"/container/heap/_obj/_go_.6 -p container/heap -I "$WORK" ./heap.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/container/heap.a "$WORK"/container/heap/_obj/_go_.6
mkdir -p "$GOROOT"/pkg/openbsd_amd64/container/
cp "$WORK"/container/heap.a "$GOROOT"/pkg/openbsd_amd64/container/heap.a

#
# strings
#

mkdir -p "$WORK"/strings/_obj/
cd "$GOROOT"/src/pkg/strings
"$GOROOT"/bin/go-tool/6g -o "$WORK"/strings/_obj/_go_.6 -p strings -I "$WORK" ./reader.go ./replace.go ./strings.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/strings.a "$WORK"/strings/_obj/_go_.6
cp "$WORK"/strings.a "$GOROOT"/pkg/openbsd_amd64/strings.a

#
# strconv
#

mkdir -p "$WORK"/strconv/_obj/
cd "$GOROOT"/src/pkg/strconv
"$GOROOT"/bin/go-tool/6g -o "$WORK"/strconv/_obj/_go_.6 -p strconv -I "$WORK" ./atob.go ./atof.go ./atoi.go ./decimal.go ./extfloat.go ./ftoa.go ./itoa.go ./quote.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/strconv.a "$WORK"/strconv/_obj/_go_.6
cp "$WORK"/strconv.a "$GOROOT"/pkg/openbsd_amd64/strconv.a

#
# encoding/base64
#

mkdir -p "$WORK"/encoding/base64/_obj/
cd "$GOROOT"/src/pkg/encoding/base64
"$GOROOT"/bin/go-tool/6g -o "$WORK"/encoding/base64/_obj/_go_.6 -p encoding/base64 -I "$WORK" ./base64.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/encoding/base64.a "$WORK"/encoding/base64/_obj/_go_.6
mkdir -p "$GOROOT"/pkg/openbsd_amd64/encoding/
cp "$WORK"/encoding/base64.a "$GOROOT"/pkg/openbsd_amd64/encoding/base64.a

#
# syscall
#

mkdir -p "$WORK"/syscall/_obj/
cd "$GOROOT"/src/pkg/syscall
"$GOROOT"/bin/go-tool/6g -o "$WORK"/syscall/_obj/_go_.6 -p syscall -I "$WORK" ./bpf_bsd.go ./env_unix.go ./exec_bsd.go ./exec_unix.go ./route_bsd.go ./route_openbsd.go ./sockcmsg_unix.go ./str.go ./syscall.go ./syscall_bsd.go ./syscall_openbsd.go ./syscall_openbsd_amd64.go ./syscall_unix.go ./zerrors_openbsd_amd64.go ./zsyscall_openbsd_amd64.go ./zsysctl_openbsd.go ./zsysnum_openbsd_amd64.go ./ztypes_openbsd_amd64.go
"$GOROOT"/bin/go-tool/6a -I "$WORK"/syscall/_obj/ -o "$WORK"/syscall/_obj/asm_openbsd_amd64.6 -DGOOS_openbsd -DGOARCH_amd64 ./asm_openbsd_amd64.s
"$GOROOT"/bin/go-tool/pack grc "$WORK"/syscall.a "$WORK"/syscall/_obj/_go_.6 "$WORK"/syscall/_obj/asm_openbsd_amd64.6
cp "$WORK"/syscall.a "$GOROOT"/pkg/openbsd_amd64/syscall.a

#
# time
#

mkdir -p "$WORK"/time/_obj/
cd "$GOROOT"/src/pkg/time
"$GOROOT"/bin/go-tool/6g -o "$WORK"/time/_obj/_go_.6 -p time -I "$WORK" ./format.go ./sleep.go ./sys_unix.go ./tick.go ./time.go ./zoneinfo.go ./zoneinfo_unix.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/time.a "$WORK"/time/_obj/_go_.6
cp "$WORK"/time.a "$GOROOT"/pkg/openbsd_amd64/time.a

#
# os
#

mkdir -p "$WORK"/os/_obj/
cd "$GOROOT"/src/pkg/os
"$GOROOT"/bin/go-tool/6g -o "$WORK"/os/_obj/_go_.6 -p os -I "$WORK" ./dir_unix.go ./doc.go ./env.go ./error.go ./error_posix.go ./exec.go ./exec_posix.go ./exec_unix.go ./file.go ./file_posix.go ./file_unix.go ./getwd.go ./path.go ./path_unix.go ./proc.go ./stat_openbsd.go ./sys_bsd.go ./time.go ./types.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/os.a "$WORK"/os/_obj/_go_.6
cp "$WORK"/os.a "$GOROOT"/pkg/openbsd_amd64/os.a

#
# reflect
#

mkdir -p "$WORK"/reflect/_obj/
cd "$GOROOT"/src/pkg/reflect
"$GOROOT"/bin/go-tool/6g -o "$WORK"/reflect/_obj/_go_.6 -p reflect -I "$WORK" ./deepequal.go ./type.go ./value.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/reflect.a "$WORK"/reflect/_obj/_go_.6
cp "$WORK"/reflect.a "$GOROOT"/pkg/openbsd_amd64/reflect.a

#
# fmt
#

mkdir -p "$WORK"/fmt/_obj/
cd "$GOROOT"/src/pkg/fmt
"$GOROOT"/bin/go-tool/6g -o "$WORK"/fmt/_obj/_go_.6 -p fmt -I "$WORK" ./doc.go ./format.go ./print.go ./scan.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/fmt.a "$WORK"/fmt/_obj/_go_.6
cp "$WORK"/fmt.a "$GOROOT"/pkg/openbsd_amd64/fmt.a

#
# unicode/utf16
#

mkdir -p "$WORK"/unicode/utf16/_obj/
cd "$GOROOT"/src/pkg/unicode/utf16
"$GOROOT"/bin/go-tool/6g -o "$WORK"/unicode/utf16/_obj/_go_.6 -p unicode/utf16 -I "$WORK" ./utf16.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/unicode/utf16.a "$WORK"/unicode/utf16/_obj/_go_.6
cp "$WORK"/unicode/utf16.a "$GOROOT"/pkg/openbsd_amd64/unicode/utf16.a

#
# encoding/json
#

mkdir -p "$WORK"/encoding/json/_obj/
cd "$GOROOT"/src/pkg/encoding/json
"$GOROOT"/bin/go-tool/6g -o "$WORK"/encoding/json/_obj/_go_.6 -p encoding/json -I "$WORK" ./decode.go ./encode.go ./indent.go ./scanner.go ./stream.go ./tags.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/encoding/json.a "$WORK"/encoding/json/_obj/_go_.6
cp "$WORK"/encoding/json.a "$GOROOT"/pkg/openbsd_amd64/encoding/json.a

#
# flag
#

mkdir -p "$WORK"/flag/_obj/
cd "$GOROOT"/src/pkg/flag
"$GOROOT"/bin/go-tool/6g -o "$WORK"/flag/_obj/_go_.6 -p flag -I "$WORK" ./flag.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/flag.a "$WORK"/flag/_obj/_go_.6
cp "$WORK"/flag.a "$GOROOT"/pkg/openbsd_amd64/flag.a

#
# bufio
#

mkdir -p "$WORK"/bufio/_obj/
cd "$GOROOT"/src/pkg/bufio
"$GOROOT"/bin/go-tool/6g -o "$WORK"/bufio/_obj/_go_.6 -p bufio -I "$WORK" ./bufio.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/bufio.a "$WORK"/bufio/_obj/_go_.6
cp "$WORK"/bufio.a "$GOROOT"/pkg/openbsd_amd64/bufio.a

#
# encoding/gob
#

mkdir -p "$WORK"/encoding/gob/_obj/
cd "$GOROOT"/src/pkg/encoding/gob
"$GOROOT"/bin/go-tool/6g -o "$WORK"/encoding/gob/_obj/_go_.6 -p encoding/gob -I "$WORK" ./decode.go ./decoder.go ./doc.go ./encode.go ./encoder.go ./error.go ./type.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/encoding/gob.a "$WORK"/encoding/gob/_obj/_go_.6
cp "$WORK"/encoding/gob.a "$GOROOT"/pkg/openbsd_amd64/encoding/gob.a

#
# go/token
#

mkdir -p "$WORK"/go/token/_obj/
cd "$GOROOT"/src/pkg/go/token
"$GOROOT"/bin/go-tool/6g -o "$WORK"/go/token/_obj/_go_.6 -p go/token -I "$WORK" ./position.go ./serialize.go ./token.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/go/token.a "$WORK"/go/token/_obj/_go_.6
mkdir -p "$GOROOT"/pkg/openbsd_amd64/go/
cp "$WORK"/go/token.a "$GOROOT"/pkg/openbsd_amd64/go/token.a

#
# path/filepath
#

mkdir -p "$WORK"/path/filepath/_obj/
cd "$GOROOT"/src/pkg/path/filepath
"$GOROOT"/bin/go-tool/6g -o "$WORK"/path/filepath/_obj/_go_.6 -p path/filepath -I "$WORK" ./match.go ./path.go ./path_unix.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/path/filepath.a "$WORK"/path/filepath/_obj/_go_.6
mkdir -p "$GOROOT"/pkg/openbsd_amd64/path/
cp "$WORK"/path/filepath.a "$GOROOT"/pkg/openbsd_amd64/path/filepath.a

#
# go/scanner
#

mkdir -p "$WORK"/go/scanner/_obj/
cd "$GOROOT"/src/pkg/go/scanner
"$GOROOT"/bin/go-tool/6g -o "$WORK"/go/scanner/_obj/_go_.6 -p go/scanner -I "$WORK" ./errors.go ./scanner.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/go/scanner.a "$WORK"/go/scanner/_obj/_go_.6
cp "$WORK"/go/scanner.a "$GOROOT"/pkg/openbsd_amd64/go/scanner.a

#
# go/ast
#

mkdir -p "$WORK"/go/ast/_obj/
cd "$GOROOT"/src/pkg/go/ast
"$GOROOT"/bin/go-tool/6g -o "$WORK"/go/ast/_obj/_go_.6 -p go/ast -I "$WORK" ./ast.go ./filter.go ./import.go ./print.go ./resolve.go ./scope.go ./walk.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/go/ast.a "$WORK"/go/ast/_obj/_go_.6
cp "$WORK"/go/ast.a "$GOROOT"/pkg/openbsd_amd64/go/ast.a

#
# io/ioutil
#

mkdir -p "$WORK"/io/ioutil/_obj/
cd "$GOROOT"/src/pkg/io/ioutil
"$GOROOT"/bin/go-tool/6g -o "$WORK"/io/ioutil/_obj/_go_.6 -p io/ioutil -I "$WORK" ./ioutil.go ./tempfile.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/io/ioutil.a "$WORK"/io/ioutil/_obj/_go_.6
mkdir -p "$GOROOT"/pkg/openbsd_amd64/io/
cp "$WORK"/io/ioutil.a "$GOROOT"/pkg/openbsd_amd64/io/ioutil.a

#
# go/parser
#

mkdir -p "$WORK"/go/parser/_obj/
cd "$GOROOT"/src/pkg/go/parser
"$GOROOT"/bin/go-tool/6g -o "$WORK"/go/parser/_obj/_go_.6 -p go/parser -I "$WORK" ./interface.go ./parser.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/go/parser.a "$WORK"/go/parser/_obj/_go_.6
cp "$WORK"/go/parser.a "$GOROOT"/pkg/openbsd_amd64/go/parser.a

#
# log
#

mkdir -p "$WORK"/log/_obj/
cd "$GOROOT"/src/pkg/log
"$GOROOT"/bin/go-tool/6g -o "$WORK"/log/_obj/_go_.6 -p log -I "$WORK" ./log.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/log.a "$WORK"/log/_obj/_go_.6
cp "$WORK"/log.a "$GOROOT"/pkg/openbsd_amd64/log.a

#
# path
#

mkdir -p "$WORK"/path/_obj/
cd "$GOROOT"/src/pkg/path
"$GOROOT"/bin/go-tool/6g -o "$WORK"/path/_obj/_go_.6 -p path -I "$WORK" ./match.go ./path.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/path.a "$WORK"/path/_obj/_go_.6
cp "$WORK"/path.a "$GOROOT"/pkg/openbsd_amd64/path.a

#
# go/build
#

mkdir -p "$WORK"/go/build/_obj/
cd "$GOROOT"/src/pkg/go/build
"$GOROOT"/bin/go-tool/6g -o "$WORK"/go/build/_obj/_go_.6 -p go/build -I "$WORK" ./build.go ./dir.go ./path.go ./syslist.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/go/build.a "$WORK"/go/build/_obj/_go_.6
cp "$WORK"/go/build.a "$GOROOT"/pkg/openbsd_amd64/go/build.a

#
# os/exec
#

mkdir -p "$WORK"/os/exec/_obj/
cd "$GOROOT"/src/pkg/os/exec
"$GOROOT"/bin/go-tool/6g -o "$WORK"/os/exec/_obj/_go_.6 -p os/exec -I "$WORK" ./exec.go ./lp_unix.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/os/exec.a "$WORK"/os/exec/_obj/_go_.6
mkdir -p "$GOROOT"/pkg/openbsd_amd64/os/
cp "$WORK"/os/exec.a "$GOROOT"/pkg/openbsd_amd64/os/exec.a

#
# regexp/syntax
#

mkdir -p "$WORK"/regexp/syntax/_obj/
cd "$GOROOT"/src/pkg/regexp/syntax
"$GOROOT"/bin/go-tool/6g -o "$WORK"/regexp/syntax/_obj/_go_.6 -p regexp/syntax -I "$WORK" ./compile.go ./parse.go ./perl_groups.go ./prog.go ./regexp.go ./simplify.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/regexp/syntax.a "$WORK"/regexp/syntax/_obj/_go_.6
mkdir -p "$GOROOT"/pkg/openbsd_amd64/regexp/
cp "$WORK"/regexp/syntax.a "$GOROOT"/pkg/openbsd_amd64/regexp/syntax.a

#
# regexp
#

mkdir -p "$WORK"/regexp/_obj/
cd "$GOROOT"/src/pkg/regexp
"$GOROOT"/bin/go-tool/6g -o "$WORK"/regexp/_obj/_go_.6 -p regexp -I "$WORK" ./exec.go ./regexp.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/regexp.a "$WORK"/regexp/_obj/_go_.6
cp "$WORK"/regexp.a "$GOROOT"/pkg/openbsd_amd64/regexp.a

#
# net/url
#

mkdir -p "$WORK"/net/url/_obj/
cd "$GOROOT"/src/pkg/net/url
"$GOROOT"/bin/go-tool/6g -o "$WORK"/net/url/_obj/_go_.6 -p net/url -I "$WORK" ./url.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/net/url.a "$WORK"/net/url/_obj/_go_.6
mkdir -p "$GOROOT"/pkg/openbsd_amd64/net/
cp "$WORK"/net/url.a "$GOROOT"/pkg/openbsd_amd64/net/url.a

#
# text/template/parse
#

mkdir -p "$WORK"/text/template/parse/_obj/
cd "$GOROOT"/src/pkg/text/template/parse
"$GOROOT"/bin/go-tool/6g -o "$WORK"/text/template/parse/_obj/_go_.6 -p text/template/parse -I "$WORK" ./lex.go ./node.go ./parse.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/text/template/parse.a "$WORK"/text/template/parse/_obj/_go_.6
mkdir -p "$GOROOT"/pkg/openbsd_amd64/text/template/
cp "$WORK"/text/template/parse.a "$GOROOT"/pkg/openbsd_amd64/text/template/parse.a

#
# text/template
#

mkdir -p "$WORK"/text/template/_obj/
cd "$GOROOT"/src/pkg/text/template
"$GOROOT"/bin/go-tool/6g -o "$WORK"/text/template/_obj/_go_.6 -p text/template -I "$WORK" ./doc.go ./exec.go ./funcs.go ./helper.go ./template.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/text/template.a "$WORK"/text/template/_obj/_go_.6
mkdir -p "$GOROOT"/pkg/openbsd_amd64/text/
cp "$WORK"/text/template.a "$GOROOT"/pkg/openbsd_amd64/text/template.a

#
# cmd/go
#

mkdir -p "$WORK"/cmd/go/_obj/
cd "$GOROOT"/src/cmd/go
"$GOROOT"/bin/go-tool/6g -o "$WORK"/cmd/go/_obj/_go_.6 -p cmd/go -I "$WORK" ./bootstrap.go ./build.go ./clean.go ./fix.go ./fmt.go ./get.go ./help.go ./list.go ./main.go ./pkg.go ./run.go ./test.go ./testflag.go ./tool.go ./vcs.go ./version.go ./vet.go
"$GOROOT"/bin/go-tool/pack grc "$WORK"/cmd/go.a "$WORK"/cmd/go/_obj/_go_.6
"$GOROOT"/bin/go-tool/6l -o "$WORK"/cmd/go/_obj/a.out -L "$WORK" "$WORK"/cmd/go.a
mkdir -p "$GOBIN"/
cp "$WORK"/cmd/go/_obj/a.out "$GOBIN"/go_bootstrap
