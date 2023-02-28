// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkginit

import (
	"cmd/compile/internal/base"
	"cmd/compile/internal/ir"
	"cmd/compile/internal/noder"
	"cmd/compile/internal/objw"
	"cmd/compile/internal/staticinit"
	"cmd/compile/internal/typecheck"
	"cmd/compile/internal/types"
	"cmd/internal/obj"
	"cmd/internal/objabi"
	"cmd/internal/src"
	"fmt"
	"os"
)

// MakeInit creates a synthetic init function to handle any
// package-scope initialization statements.
//
// TODO(mdempsky): Move into noder, so that the types2-based frontends
// can use Info.InitOrder instead.
func MakeInit() {
	nf := initOrder(typecheck.Target.Decls)
	if len(nf) == 0 {
		return
	}

	// Make a function that contains all the initialization statements.
	base.Pos = nf[0].Pos() // prolog/epilog gets line number of first init stmt
	initializers := typecheck.Lookup("init")
	fn := typecheck.DeclFunc(initializers, nil, nil, nil)
	for _, dcl := range typecheck.InitTodoFunc.Dcl {
		dcl.Curfn = fn
	}
	fn.Dcl = append(fn.Dcl, typecheck.InitTodoFunc.Dcl...)
	typecheck.InitTodoFunc.Dcl = nil
	fn.SetIsPackageInit(true)

	// Outline (if legal/profitable) global map inits.
	newfuncs := []*ir.Func{}
	nf, newfuncs = staticinit.OutlineMapInits(nf)

	// Suppress useless "can inline" diagnostics.
	// Init functions are only called dynamically.
	fn.SetInlinabilityChecked(true)
	for _, nfn := range newfuncs {
		nfn.SetInlinabilityChecked(true)
	}

	fn.Body = nf
	typecheck.FinishFuncBody()

	typecheck.Func(fn)
	ir.WithFunc(fn, func() {
		typecheck.Stmts(nf)
	})
	typecheck.Target.Decls = append(typecheck.Target.Decls, fn)
	if base.Debug.WrapGlobalMapDbg > 1 {
		fmt.Fprintf(os.Stderr, "=-= len(newfuncs) is %d for %v\n",
			len(newfuncs), fn)
	}
	for _, nfn := range newfuncs {
		if base.Debug.WrapGlobalMapDbg > 1 {
			fmt.Fprintf(os.Stderr, "=-= add to target.decls %v\n", nfn)
		}
		typecheck.Target.Decls = append(typecheck.Target.Decls, ir.Node(nfn))
	}

	// Prepend to Inits, so it runs first, before any user-declared init
	// functions.
	typecheck.Target.Inits = append([]*ir.Func{fn}, typecheck.Target.Inits...)

	if typecheck.InitTodoFunc.Dcl != nil {
		// We only generate temps using InitTodoFunc if there
		// are package-scope initialization statements, so
		// something's weird if we get here.
		base.Fatalf("InitTodoFunc still has declarations")
	}
	typecheck.InitTodoFunc = nil
}

// Task makes and returns an initialization record for the package.
// See runtime/proc.go:initTask for its layout.
// The 3 tasks for initialization are:
//  1. Initialize all of the packages the current package depends on.
//  2. Initialize all the variables that have initializers.
//  3. Run any init functions.
func Task() *ir.Name {
	var deps []*obj.LSym // initTask records for packages the current package depends on
	var fns []*obj.LSym  // functions to call for package initialization

	// Find imported packages with init tasks.
	for _, pkg := range typecheck.Target.Imports {
		n, ok := pkg.Lookup(".inittask").Def.(*ir.Name)
		if !ok {
			continue
		}
		if n.Op() != ir.ONAME || n.Class != ir.PEXTERN {
			base.Fatalf("bad inittask: %v", n)
		}
		deps = append(deps, n.Linksym())
	}
	if base.Flag.ASan {
		// Make an initialization function to call runtime.asanregisterglobals to register an
		// array of instrumented global variables when -asan is enabled. An instrumented global
		// variable is described by a structure.
		// See the _asan_global structure declared in src/runtime/asan/asan.go.
		//
		// func init {
		// 		var globals []_asan_global {...}
		// 		asanregisterglobals(&globals[0], len(globals))
		// }
		for _, n := range typecheck.Target.Externs {
			if canInstrumentGlobal(n) {
				name := n.Sym().Name
				InstrumentGlobalsMap[name] = n
				InstrumentGlobalsSlice = append(InstrumentGlobalsSlice, n)
			}
		}
		ni := len(InstrumentGlobalsMap)
		if ni != 0 {
			// Make an init._ function.
			base.Pos = base.AutogeneratedPos
			typecheck.DeclContext = ir.PEXTERN
			name := noder.Renameinit()
			fnInit := typecheck.DeclFunc(name, nil, nil, nil)

			// Get an array of instrumented global variables.
			globals := instrumentGlobals(fnInit)

			// Call runtime.asanregisterglobals function to poison redzones.
			// runtime.asanregisterglobals(unsafe.Pointer(&globals[0]), ni)
			asanf := typecheck.NewName(ir.Pkgs.Runtime.Lookup("asanregisterglobals"))
			ir.MarkFunc(asanf)
			asanf.SetType(types.NewSignature(nil, []*types.Field{
				types.NewField(base.Pos, nil, types.Types[types.TUNSAFEPTR]),
				types.NewField(base.Pos, nil, types.Types[types.TUINTPTR]),
			}, nil))
			asancall := ir.NewCallExpr(base.Pos, ir.OCALL, asanf, nil)
			asancall.Args.Append(typecheck.ConvNop(typecheck.NodAddr(
				ir.NewIndexExpr(base.Pos, globals, ir.NewInt(base.Pos, 0))), types.Types[types.TUNSAFEPTR]))
			asancall.Args.Append(typecheck.ConvNop(ir.NewInt(base.Pos, int64(ni)), types.Types[types.TUINTPTR]))

			fnInit.Body.Append(asancall)
			typecheck.FinishFuncBody()
			typecheck.Func(fnInit)
			ir.CurFunc = fnInit
			typecheck.Stmts(fnInit.Body)
			ir.CurFunc = nil

			typecheck.Target.Decls = append(typecheck.Target.Decls, fnInit)
			typecheck.Target.Inits = append(typecheck.Target.Inits, fnInit)
		}
	}

	// Record user init functions.
	for _, fn := range typecheck.Target.Inits {
		if fn.Sym().Name == "init" {
			// Synthetic init function for initialization of package-scope
			// variables. We can use staticinit to optimize away static
			// assignments.
			s := staticinit.Schedule{
				Plans: make(map[ir.Node]*staticinit.Plan),
				Temps: make(map[ir.Node]*ir.Name),
			}
			for _, n := range fn.Body {
				s.StaticInit(n)
			}
			fn.Body = s.Out
			ir.WithFunc(fn, func() {
				typecheck.Stmts(fn.Body)
			})

			if len(fn.Body) == 0 {
				fn.Body = []ir.Node{ir.NewBlockStmt(src.NoXPos, nil)}
			}
		}

		// Skip init functions with empty bodies.
		if len(fn.Body) == 1 {
			if stmt := fn.Body[0]; stmt.Op() == ir.OBLOCK && len(stmt.(*ir.BlockStmt).List) == 0 {
				continue
			}
		}
		fns = append(fns, fn.Nname.Linksym())
	}

	if len(deps) == 0 && len(fns) == 0 && types.LocalPkg.Path != "main" && types.LocalPkg.Path != "runtime" {
		return nil // nothing to initialize
	}

	// Make an .inittask structure.
	sym := typecheck.Lookup(".inittask")
	task := typecheck.NewName(sym)
	task.SetType(types.Types[types.TUINT8]) // fake type
	task.Class = ir.PEXTERN
	sym.Def = task
	lsym := task.Linksym()
	ot := 0
	ot = objw.Uint32(lsym, ot, 0) // state: not initialized yet
	ot = objw.Uint32(lsym, ot, uint32(len(fns)))
	for _, f := range fns {
		ot = objw.SymPtr(lsym, ot, f, 0)
	}

	// Add relocations which tell the linker all of the packages
	// that this package depends on (and thus, all of the packages
	// that need to be initialized before this one).
	for _, d := range deps {
		r := obj.Addrel(lsym)
		r.Type = objabi.R_INITORDER
		r.Sym = d
	}
	// An initTask has pointers, but none into the Go heap.
	// It's not quite read only, the state field must be modifiable.
	objw.Global(lsym, int32(ot), obj.NOPTR)
	return task
}

// initRequiredForCoverage returns TRUE if we need to force creation
// of an init function for the package so as to insert a coverage
// runtime registration call.
func initRequiredForCoverage(l []ir.Node) bool {
	if base.Flag.Cfg.CoverageInfo == nil {
		return false
	}
	for _, n := range l {
		if n.Op() == ir.ODCLFUNC {
			return true
		}
	}
	return false
}
