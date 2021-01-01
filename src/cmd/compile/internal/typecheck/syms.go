// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typecheck

import (
	"cmd/compile/internal/base"
	"cmd/compile/internal/ir"
	"cmd/compile/internal/types"
	"cmd/internal/obj"
	"cmd/internal/src"
)

func LookupRuntime(name string) *ir.Name {
	s := ir.Pkgs.Runtime.Lookup(name)
	if s == nil || s.Def == nil {
		base.Fatalf("syslook: can't find runtime.%s", name)
	}
	return ir.AsNode(s.Def).(*ir.Name)
}

// SubstArgTypes substitutes the given list of types for
// successive occurrences of the "any" placeholder in the
// type syntax expression n.Type.
// The result of SubstArgTypes MUST be assigned back to old, e.g.
// 	n.Left = SubstArgTypes(n.Left, t1, t2)
func SubstArgTypes(old *ir.Name, types_ ...*types.Type) *ir.Name {
	n := old.CloneName()

	for _, t := range types_ {
		types.CalcSize(t)
	}
	n.SetType(types.SubstAny(n.Type(), &types_))
	if len(types_) > 0 {
		base.Fatalf("substArgTypes: too many argument types")
	}
	return n
}

// AutoLabel generates a new Name node for use with
// an automatically generated label.
// prefix is a short mnemonic (e.g. ".s" for switch)
// to help with debugging.
// It should begin with "." to avoid conflicts with
// user labels.
func AutoLabel(prefix string) *types.Sym {
	if prefix[0] != '.' {
		base.Fatalf("autolabel prefix must start with '.', have %q", prefix)
	}
	fn := ir.CurFunc
	if ir.CurFunc == nil {
		base.Fatalf("autolabel outside function")
	}
	n := fn.Label
	fn.Label++
	return LookupNum(prefix, int(n))
}

func Lookup(name string) *types.Sym {
	return types.LocalPkg.Lookup(name)
}

// InitRuntime loads the definitions for the low-level runtime functions,
// so that the compiler can generate calls to them,
// but does not make them visible to user code.
func InitRuntime() {
	types.Block = 1

	inimport = true
	TypecheckAllowed = true

	typs := runtimeTypes()
	for _, d := range &runtimeDecls {
		sym := ir.Pkgs.Runtime.Lookup(d.name)
		typ := typs[d.typ]
		switch d.tag {
		case funcTag:
			importfunc(ir.Pkgs.Runtime, src.NoXPos, sym, typ)
		case varTag:
			importvar(ir.Pkgs.Runtime, src.NoXPos, sym, typ)
		default:
			base.Fatalf("unhandled declaration tag %v", d.tag)
		}
	}

	TypecheckAllowed = false
	inimport = false
}

// LookupRuntimeFunc looks up Go function name in package runtime. This function
// must follow the internal calling convention.
func LookupRuntimeFunc(name string) *obj.LSym {
	s := ir.Pkgs.Runtime.Lookup(name)
	s.SetFunc(true)
	return s.Linksym()
}

// LookupRuntimeVar looks up a variable (or assembly function) name in package
// runtime. If this is a function, it may have a special calling
// convention.
func LookupRuntimeVar(name string) *obj.LSym {
	return ir.Pkgs.Runtime.Lookup(name).Linksym()
}
