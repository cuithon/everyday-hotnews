// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gc

import (
	"cmd/compile/internal/base"
	"cmd/compile/internal/ir"
	"cmd/compile/internal/types"
	"cmd/internal/obj"
	"cmd/internal/src"
	"strconv"
)

// sysfunc looks up Go function name in package runtime. This function
// must follow the internal calling convention.
func sysfunc(name string) *obj.LSym {
	s := ir.Pkgs.Runtime.Lookup(name)
	s.SetFunc(true)
	return s.Linksym()
}

// sysvar looks up a variable (or assembly function) name in package
// runtime. If this is a function, it may have a special calling
// convention.
func sysvar(name string) *obj.LSym {
	return ir.Pkgs.Runtime.Lookup(name).Linksym()
}

// isParamStackCopy reports whether this is the on-stack copy of a
// function parameter that moved to the heap.
func isParamStackCopy(n ir.Node) bool {
	if n.Op() != ir.ONAME {
		return false
	}
	name := n.(*ir.Name)
	return (name.Class_ == ir.PPARAM || name.Class_ == ir.PPARAMOUT) && name.Heapaddr != nil
}

// isParamHeapCopy reports whether this is the on-heap copy of
// a function parameter that moved to the heap.
func isParamHeapCopy(n ir.Node) bool {
	if n.Op() != ir.ONAME {
		return false
	}
	name := n.(*ir.Name)
	return name.Class_ == ir.PAUTOHEAP && name.Name().Stackcopy != nil
}

// autotmpname returns the name for an autotmp variable numbered n.
func autotmpname(n int) string {
	// Give each tmp a different name so that they can be registerized.
	// Add a preceding . to avoid clashing with legal names.
	const prefix = ".autotmp_"
	// Start with a buffer big enough to hold a large n.
	b := []byte(prefix + "      ")[:len(prefix)]
	b = strconv.AppendInt(b, int64(n), 10)
	return types.InternString(b)
}

// make a new Node off the books
func tempAt(pos src.XPos, curfn *ir.Func, t *types.Type) *ir.Name {
	if curfn == nil {
		base.Fatalf("no curfn for tempAt")
	}
	if curfn.Op() == ir.OCLOSURE {
		ir.Dump("tempAt", curfn)
		base.Fatalf("adding tempAt to wrong closure function")
	}
	if t == nil {
		base.Fatalf("tempAt called with nil type")
	}

	s := &types.Sym{
		Name: autotmpname(len(curfn.Dcl)),
		Pkg:  types.LocalPkg,
	}
	n := ir.NewNameAt(pos, s)
	s.Def = n
	n.SetType(t)
	n.Class_ = ir.PAUTO
	n.SetEsc(EscNever)
	n.Curfn = curfn
	n.SetUsed(true)
	n.SetAutoTemp(true)
	curfn.Dcl = append(curfn.Dcl, n)

	dowidth(t)

	return n
}

func temp(t *types.Type) *ir.Name {
	return tempAt(base.Pos, Curfn, t)
}
