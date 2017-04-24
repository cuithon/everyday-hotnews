// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gc

import (
	"cmd/compile/internal/types"
)

// A function named init is a special case.
// It is called by the initialization before main is run.
// To make it unique within a package and also uncallable,
// the name, normally "pkg.init", is altered to "pkg.init.0".
var renameinitgen int

func renameinit() *types.Sym {
	s := lookupN("init.", renameinitgen)
	renameinitgen++
	return s
}

// anyinit reports whether there any interesting init statements.
func anyinit(n []*Node) bool {
	for _, ln := range n {
		switch ln.Op {
		case ODCLFUNC, ODCLCONST, ODCLTYPE, OEMPTY:
		case OAS:
			if !isblank(ln.Left) || !candiscard(ln.Right) {
				return true
			}
		default:
			return true
		}
	}

	// is this main
	if localpkg.Name == "main" {
		return true
	}

	// is there an explicit init function
	if renameinitgen > 0 {
		return true
	}

	// are there any imported init functions
	for _, s := range types.InitSyms {
		if s.Def != nil {
			return true
		}
	}

	// then none
	return false
}

// fninit hand-crafts package initialization code.
//
//      var initdone· uint8                             (1)
//      func init() {                                   (2)
//              if initdone· > 1 {                      (3)
//                      return                          (3a)
//              }
//              if initdone· == 1 {                     (4)
//                      throw()                         (4a)
//              }
//              initdone· = 1                           (5)
//              // over all matching imported symbols
//                      <pkg>.init()                    (6)
//              { <init stmts> }                        (7)
//              init.<n>() // if any                    (8)
//              initdone· = 2                           (9)
//              return                                  (10)
//      }
func fninit(n []*Node) {
	lineno = autogeneratedPos
	nf := initfix(n)
	if !anyinit(nf) {
		return
	}

	var r []*Node

	// (1)
	gatevar := newname(lookup("initdone·"))
	addvar(gatevar, types.Types[TUINT8], PEXTERN)

	// (2)
	initsym := lookup("init")
	fn := dclfunc(initsym, nod(OTFUNC, nil, nil))

	// (3)
	a := nod(OIF, nil, nil)
	a.Left = nod(OGT, gatevar, nodintconst(1))
	a.SetLikely(true)
	r = append(r, a)
	// (3a)
	a.Nbody.Set1(nod(ORETURN, nil, nil))

	// (4)
	b := nod(OIF, nil, nil)
	b.Left = nod(OEQ, gatevar, nodintconst(1))
	// this actually isn't likely, but code layout is better
	// like this: no JMP needed after the call.
	b.SetLikely(true)
	r = append(r, b)
	// (4a)
	b.Nbody.Set1(nod(OCALL, syslook("throwinit"), nil))

	// (5)
	a = nod(OAS, gatevar, nodintconst(1))

	r = append(r, a)

	// (6)
	for _, s := range types.InitSyms {
		if s.Def != nil && s != initsym {
			n := asNode(s.Def)
			n.checkInitFuncSignature()
			a = nod(OCALL, n, nil)
			r = append(r, a)
		}
	}

	// (7)
	r = append(r, nf...)

	// (8)

	// maxInlineInitCalls is the threshold at which we switch
	// from generating calls inline to generating a static array
	// of functions and calling them in a loop.
	// See CL 41500 for more discussion.
	const maxInlineInitCalls = 500

	if renameinitgen < maxInlineInitCalls {
		// Not many init functions. Just call them all directly.
		for i := 0; i < renameinitgen; i++ {
			s := lookupN("init.", i)
			n := asNode(s.Def)
			n.checkInitFuncSignature()
			a = nod(OCALL, n, nil)
			r = append(r, a)
		}
	} else {
		// Lots of init functions.
		// Set up an array of functions and loop to call them.
		// This is faster to compile and similar at runtime.

		// Build type [renameinitgen]func().
		typ := types.NewArray(functype(nil, nil, nil), int64(renameinitgen))

		// Make and fill array.
		fnarr := staticname(typ)
		fnarr.Name.SetReadonly(true)
		for i := 0; i < renameinitgen; i++ {
			s := lookupN("init.", i)
			lhs := nod(OINDEX, fnarr, nodintconst(int64(i)))
			rhs := asNode(s.Def)
			rhs.checkInitFuncSignature()
			as := nod(OAS, lhs, rhs)
			as = typecheck(as, Etop)
			genAsStatic(as)
		}

		// Generate a loop that calls each function in turn.
		// for i := 0; i < renameinitgen; i++ {
		//   fnarr[i]()
		// }
		i := temp(types.Types[TINT])
		fnidx := nod(OINDEX, fnarr, i)
		fnidx.SetBounded(true)

		zero := nod(OAS, i, nodintconst(0))
		cond := nod(OLT, i, nodintconst(int64(renameinitgen)))
		incr := nod(OAS, i, nod(OADD, i, nodintconst(1)))
		body := nod(OCALL, fnidx, nil)

		loop := nod(OFOR, cond, incr)
		loop.Nbody.Set1(body)
		loop.Ninit.Set1(zero)

		loop = typecheck(loop, Etop)
		loop = walkstmt(loop)
		r = append(r, loop)
	}

	// (9)
	a = nod(OAS, gatevar, nodintconst(2))

	r = append(r, a)

	// (10)
	a = nod(ORETURN, nil, nil)

	r = append(r, a)
	exportsym(fn.Func.Nname)

	fn.Nbody.Set(r)
	funcbody(fn)

	Curfn = fn
	fn = typecheck(fn, Etop)
	typecheckslice(r, Etop)
	Curfn = nil
	funccompile(fn)
}

func (n *Node) checkInitFuncSignature() {
	ft := n.Type.FuncType()
	if ft.Receiver.Fields().Len()+ft.Params.Fields().Len()+ft.Results.Fields().Len() > 0 {
		Fatalf("init function cannot have receiver, params, or results: %v (%v)", n, n.Type)
	}
}
