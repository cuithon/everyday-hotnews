// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file will evolve, since we plan to do a mix of stenciling and passing
// around dictionaries.

package noder

import (
	"bytes"
	"cmd/compile/internal/base"
	"cmd/compile/internal/ir"
	"cmd/compile/internal/objw"
	"cmd/compile/internal/reflectdata"
	"cmd/compile/internal/typecheck"
	"cmd/compile/internal/types"
	"cmd/internal/obj"
	"cmd/internal/src"
	"fmt"
	"go/constant"
	"strconv"
)

func assert(p bool) {
	base.Assert(p)
}

// Temporary - for outputting information on derived types, dictionaries, sub-dictionaries.
// Turn off when running tests.
var infoPrintMode = false

func infoPrint(format string, a ...interface{}) {
	if infoPrintMode {
		fmt.Printf(format, a...)
	}
}

// stencil scans functions for instantiated generic function calls and creates the
// required instantiations for simple generic functions. It also creates
// instantiated methods for all fully-instantiated generic types that have been
// encountered already or new ones that are encountered during the stenciling
// process.
func (g *irgen) stencil() {
	g.instInfoMap = make(map[*types.Sym]*instInfo)
	g.gfInfoMap = make(map[*types.Sym]*gfInfo)

	// Instantiate the methods of instantiated generic types that we have seen so far.
	g.instantiateMethods()

	// Don't use range(g.target.Decls) - we also want to process any new instantiated
	// functions that are created during this loop, in order to handle generic
	// functions calling other generic functions.
	for i := 0; i < len(g.target.Decls); i++ {
		decl := g.target.Decls[i]

		// Look for function instantiations in bodies of non-generic
		// functions or in global assignments (ignore global type and
		// constant declarations).
		switch decl.Op() {
		case ir.ODCLFUNC:
			if decl.Type().HasTParam() {
				// Skip any generic functions
				continue
			}
			// transformCall() below depends on CurFunc being set.
			ir.CurFunc = decl.(*ir.Func)

		case ir.OAS, ir.OAS2, ir.OAS2DOTTYPE, ir.OAS2FUNC, ir.OAS2MAPR, ir.OAS2RECV, ir.OASOP:
			// These are all the various kinds of global assignments,
			// whose right-hand-sides might contain a function
			// instantiation.

		default:
			// The other possible ops at the top level are ODCLCONST
			// and ODCLTYPE, which don't have any function
			// instantiations.
			continue
		}

		// For all non-generic code, search for any function calls using
		// generic function instantiations. Then create the needed
		// instantiated function if it hasn't been created yet, and change
		// to calling that function directly.
		modified := false
		closureRequired := false
		// declInfo will be non-nil exactly if we are scanning an instantiated function
		declInfo := g.instInfoMap[decl.Sym()]

		ir.Visit(decl, func(n ir.Node) {
			if n.Op() == ir.OFUNCINST {
				// generic F, not immediately called
				closureRequired = true
			}
			if n.Op() == ir.OMETHEXPR && len(deref(n.(*ir.SelectorExpr).X.Type()).RParams()) > 0 && !types.IsInterfaceMethod(n.(*ir.SelectorExpr).Selection.Type) {
				// T.M, T a type which is generic, not immediately
				// called. Not necessary if the method selected is
				// actually for an embedded interface field.
				closureRequired = true
			}
			if n.Op() == ir.OCALL && n.(*ir.CallExpr).X.Op() == ir.OFUNCINST {
				// We have found a function call using a generic function
				// instantiation.
				call := n.(*ir.CallExpr)
				inst := call.X.(*ir.InstExpr)
				nameNode, isMeth := g.getInstNameNode(inst)
				targs := typecheck.TypesOf(inst.Targs)
				st := g.getInstantiation(nameNode, targs, isMeth)
				dictValue, usingSubdict := g.getDictOrSubdict(declInfo, n, nameNode, targs, isMeth)
				if infoPrintMode {
					dictkind := "Main dictionary"
					if usingSubdict {
						dictkind = "Sub-dictionary"
					}
					if inst.X.Op() == ir.OMETHVALUE {
						fmt.Printf("%s in %v at generic method call: %v - %v\n", dictkind, decl, inst.X, call)
					} else {
						fmt.Printf("%s in %v at generic function call: %v - %v\n", dictkind, decl, inst.X, call)
					}
				}
				// Replace the OFUNCINST with a direct reference to the
				// new stenciled function
				call.X = st.Nname
				if inst.X.Op() == ir.OMETHVALUE {
					// When we create an instantiation of a method
					// call, we make it a function. So, move the
					// receiver to be the first arg of the function
					// call.
					call.Args.Prepend(inst.X.(*ir.SelectorExpr).X)
				}

				// Add dictionary to argument list.
				call.Args.Prepend(dictValue)
				// Transform the Call now, which changes OCALL
				// to OCALLFUNC and does typecheckaste/assignconvfn.
				transformCall(call)
				modified = true
			}
			if n.Op() == ir.OCALLMETH && n.(*ir.CallExpr).X.Op() == ir.ODOTMETH && len(deref(n.(*ir.CallExpr).X.Type().Recv().Type).RParams()) > 0 {
				// Method call on a generic type, which was instantiated by stenciling.
				// Method calls on explicitly instantiated types will have an OFUNCINST
				// and are handled above.
				call := n.(*ir.CallExpr)
				meth := call.X.(*ir.SelectorExpr)
				targs := deref(meth.Type().Recv().Type).RParams()

				t := meth.X.Type()
				baseSym := deref(t).OrigSym
				baseType := baseSym.Def.(*ir.Name).Type()
				var gf *ir.Name
				for _, m := range baseType.Methods().Slice() {
					if meth.Sel == m.Sym {
						gf = m.Nname.(*ir.Name)
						break
					}
				}

				st := g.getInstantiation(gf, targs, true)
				dictValue, usingSubdict := g.getDictOrSubdict(declInfo, n, gf, targs, true)
				// We have to be using a subdictionary, since this is
				// a generic method call.
				assert(usingSubdict)

				call.SetOp(ir.OCALL)
				call.X = st.Nname
				call.Args.Prepend(dictValue, meth.X)
				// Transform the Call now, which changes OCALL
				// to OCALLFUNC and does typecheckaste/assignconvfn.
				transformCall(call)
				modified = true
			}
		})

		// If we found a reference to a generic instantiation that wasn't an
		// immediate call, then traverse the nodes of decl again (with
		// EditChildren rather than Visit), where we actually change the
		// reference to the instantiation to a closure that captures the
		// dictionary, then does a direct call.
		// EditChildren is more expensive than Visit, so we only do this
		// in the infrequent case of an OFUNCINST without a corresponding
		// call.
		if closureRequired {
			var edit func(ir.Node) ir.Node
			var outer *ir.Func
			if f, ok := decl.(*ir.Func); ok {
				outer = f
			}
			edit = func(x ir.Node) ir.Node {
				ir.EditChildren(x, edit)
				switch {
				case x.Op() == ir.OFUNCINST:
					return g.buildClosure(outer, x)
				case x.Op() == ir.OMETHEXPR && len(deref(x.(*ir.SelectorExpr).X.Type()).RParams()) > 0 &&
					!types.IsInterfaceMethod(x.(*ir.SelectorExpr).Selection.Type): // TODO: test for ptr-to-method case
					return g.buildClosure(outer, x)
				}
				return x
			}
			edit(decl)
		}
		if base.Flag.W > 1 && modified {
			ir.Dump(fmt.Sprintf("\nmodified %v", decl), decl)
		}
		ir.CurFunc = nil
		// We may have seen new fully-instantiated generic types while
		// instantiating any needed functions/methods in the above
		// function. If so, instantiate all the methods of those types
		// (which will then lead to more function/methods to scan in the loop).
		g.instantiateMethods()
	}

}

// buildClosure makes a closure to implement x, a OFUNCINST or OMETHEXPR
// of generic type. outer is the containing function (or nil if closure is
// in a global assignment instead of a function).
func (g *irgen) buildClosure(outer *ir.Func, x ir.Node) ir.Node {
	pos := x.Pos()
	var target *ir.Func   // target instantiated function/method
	var dictValue ir.Node // dictionary to use
	var rcvrValue ir.Node // receiver, if a method value
	typ := x.Type()       // type of the closure
	var outerInfo *instInfo
	if outer != nil {
		outerInfo = g.instInfoMap[outer.Sym()]
	}
	usingSubdict := false
	valueMethod := false
	if x.Op() == ir.OFUNCINST {
		inst := x.(*ir.InstExpr)

		// Type arguments we're instantiating with.
		targs := typecheck.TypesOf(inst.Targs)

		// Find the generic function/method.
		var gf *ir.Name
		if inst.X.Op() == ir.ONAME {
			// Instantiating a generic function call.
			gf = inst.X.(*ir.Name)
		} else if inst.X.Op() == ir.OMETHVALUE {
			// Instantiating a method value x.M.
			se := inst.X.(*ir.SelectorExpr)
			rcvrValue = se.X
			gf = se.Selection.Nname.(*ir.Name)
		} else {
			panic("unhandled")
		}

		// target is the instantiated function we're trying to call.
		// For functions, the target expects a dictionary as its first argument.
		// For method values, the target expects a dictionary and the receiver
		// as its first two arguments.
		// dictValue is the value to use for the dictionary argument.
		target = g.getInstantiation(gf, targs, rcvrValue != nil)
		dictValue, usingSubdict = g.getDictOrSubdict(outerInfo, x, gf, targs, rcvrValue != nil)
		if infoPrintMode {
			dictkind := "Main dictionary"
			if usingSubdict {
				dictkind = "Sub-dictionary"
			}
			if rcvrValue == nil {
				fmt.Printf("%s in %v for generic function value %v\n", dictkind, outer, inst.X)
			} else {
				fmt.Printf("%s in %v for generic method value %v\n", dictkind, outer, inst.X)
			}
		}
	} else { // ir.OMETHEXPR
		// Method expression T.M where T is a generic type.
		se := x.(*ir.SelectorExpr)
		targs := deref(se.X.Type()).RParams()
		if len(targs) == 0 {
			panic("bad")
		}

		// se.X.Type() is the top-level type of the method expression. To
		// correctly handle method expressions involving embedded fields,
		// look up the generic method below using the type of the receiver
		// of se.Selection, since that will be the type that actually has
		// the method.
		recv := deref(se.Selection.Type.Recv().Type)
		baseType := recv.OrigSym.Def.Type()
		var gf *ir.Name
		for _, m := range baseType.Methods().Slice() {
			if se.Sel == m.Sym {
				gf = m.Nname.(*ir.Name)
				break
			}
		}
		if !gf.Type().Recv().Type.IsPtr() {
			// Remember if value method, so we can detect (*T).M case.
			valueMethod = true
		}
		target = g.getInstantiation(gf, targs, true)
		dictValue, usingSubdict = g.getDictOrSubdict(outerInfo, x, gf, targs, true)
		if infoPrintMode {
			dictkind := "Main dictionary"
			if usingSubdict {
				dictkind = "Sub-dictionary"
			}
			fmt.Printf("%s in %v for method expression %v\n", dictkind, outer, x)
		}
	}

	// Build a closure to implement a function instantiation.
	//
	//   func f[T any] (int, int) (int, int) { ...whatever... }
	//
	// Then any reference to f[int] not directly called gets rewritten to
	//
	//   .dictN := ... dictionary to use ...
	//   func(a0, a1 int) (r0, r1 int) {
	//     return .inst.f[int](.dictN, a0, a1)
	//   }
	//
	// Similarly for method expressions,
	//
	//   type g[T any] ....
	//   func (rcvr g[T]) f(a0, a1 int) (r0, r1 int) { ... }
	//
	// Any reference to g[int].f not directly called gets rewritten to
	//
	//   .dictN := ... dictionary to use ...
	//   func(rcvr g[int], a0, a1 int) (r0, r1 int) {
	//     return .inst.g[int].f(.dictN, rcvr, a0, a1)
	//   }
	//
	// Also method values
	//
	//   var x g[int]
	//
	// Any reference to x.f not directly called gets rewritten to
	//
	//   .dictN := ... dictionary to use ...
	//   x2 := x
	//   func(a0, a1 int) (r0, r1 int) {
	//     return .inst.g[int].f(.dictN, x2, a0, a1)
	//   }

	// Make a new internal function.
	fn := ir.NewClosureFunc(pos, outer != nil)
	ir.NameClosure(fn.OClosure, outer)

	// This is the dictionary we want to use.
	// It may be a constant, or it may be a dictionary acquired from the outer function's dictionary.
	// For the latter, dictVar is a variable in the outer function's scope, set to the subdictionary
	// read from the outer function's dictionary.
	var dictVar *ir.Name
	var dictAssign *ir.AssignStmt
	if outer != nil {
		// Note: for now this is a compile-time constant, so we don't really need a closure
		// to capture it (a wrapper function would work just as well). But eventually it
		// will be a read of a subdictionary from the parent dictionary.
		dictVar = ir.NewNameAt(pos, typecheck.LookupNum(".dict", g.dnum))
		g.dnum++
		dictVar.Class = ir.PAUTO
		typed(types.Types[types.TUINTPTR], dictVar)
		dictVar.Curfn = outer
		dictAssign = ir.NewAssignStmt(pos, dictVar, dictValue)
		dictAssign.SetTypecheck(1)
		dictVar.Defn = dictAssign
		outer.Dcl = append(outer.Dcl, dictVar)
	}
	// assign the receiver to a temporary.
	var rcvrVar *ir.Name
	var rcvrAssign ir.Node
	if rcvrValue != nil {
		rcvrVar = ir.NewNameAt(pos, typecheck.LookupNum(".rcvr", g.dnum))
		g.dnum++
		rcvrVar.Class = ir.PAUTO
		typed(rcvrValue.Type(), rcvrVar)
		rcvrVar.Curfn = outer
		rcvrAssign = ir.NewAssignStmt(pos, rcvrVar, rcvrValue)
		rcvrAssign.SetTypecheck(1)
		rcvrVar.Defn = rcvrAssign
		outer.Dcl = append(outer.Dcl, rcvrVar)
	}

	// Build formal argument and return lists.
	var formalParams []*types.Field  // arguments of closure
	var formalResults []*types.Field // returns of closure
	for i := 0; i < typ.NumParams(); i++ {
		t := typ.Params().Field(i).Type
		arg := ir.NewNameAt(pos, typecheck.LookupNum("a", i))
		arg.Class = ir.PPARAM
		typed(t, arg)
		arg.Curfn = fn
		fn.Dcl = append(fn.Dcl, arg)
		f := types.NewField(pos, arg.Sym(), t)
		f.Nname = arg
		formalParams = append(formalParams, f)
	}
	for i := 0; i < typ.NumResults(); i++ {
		t := typ.Results().Field(i).Type
		result := ir.NewNameAt(pos, typecheck.LookupNum("r", i)) // TODO: names not needed?
		result.Class = ir.PPARAMOUT
		typed(t, result)
		result.Curfn = fn
		fn.Dcl = append(fn.Dcl, result)
		f := types.NewField(pos, result.Sym(), t)
		f.Nname = result
		formalResults = append(formalResults, f)
	}

	// Build an internal function with the right signature.
	closureType := types.NewSignature(x.Type().Pkg(), nil, nil, formalParams, formalResults)
	typed(closureType, fn.Nname)
	typed(x.Type(), fn.OClosure)
	fn.SetTypecheck(1)

	// Build body of closure. This involves just calling the wrapped function directly
	// with the additional dictionary argument.

	// First, figure out the dictionary argument.
	var dict2Var ir.Node
	if usingSubdict {
		// Capture sub-dictionary calculated in the outer function
		dict2Var = ir.CaptureName(pos, fn, dictVar)
		typed(types.Types[types.TUINTPTR], dict2Var)
	} else {
		// Static dictionary, so can be used directly in the closure
		dict2Var = dictValue
	}
	// Also capture the receiver variable.
	var rcvr2Var *ir.Name
	if rcvrValue != nil {
		rcvr2Var = ir.CaptureName(pos, fn, rcvrVar)
	}

	// Build arguments to call inside the closure.
	var args []ir.Node

	// First the dictionary argument.
	args = append(args, dict2Var)
	// Then the receiver.
	if rcvrValue != nil {
		args = append(args, rcvr2Var)
	}
	// Then all the other arguments (including receiver for method expressions).
	for i := 0; i < typ.NumParams(); i++ {
		if x.Op() == ir.OMETHEXPR && i == 0 {
			// If we are doing a method expression, we need to
			// explicitly traverse any embedded fields in the receiver
			// argument in order to call the method instantiation.
			arg0 := formalParams[0].Nname.(ir.Node)
			arg0 = typecheck.AddImplicitDots(ir.NewSelectorExpr(base.Pos, ir.OXDOT, arg0, x.(*ir.SelectorExpr).Sel)).X
			if valueMethod && arg0.Type().IsPtr() {
				// For handling the (*T).M case: if we have a pointer
				// receiver after following all the embedded fields,
				// but it's a value method, add a star operator.
				arg0 = ir.NewStarExpr(arg0.Pos(), arg0)
			}
			args = append(args, arg0)
		} else {
			args = append(args, formalParams[i].Nname.(*ir.Name))
		}
	}

	// Build call itself.
	var innerCall ir.Node = ir.NewCallExpr(pos, ir.OCALL, target.Nname, args)
	if len(formalResults) > 0 {
		innerCall = ir.NewReturnStmt(pos, []ir.Node{innerCall})
	}
	// Finish building body of closure.
	ir.CurFunc = fn
	// TODO: set types directly here instead of using typecheck.Stmt
	typecheck.Stmt(innerCall)
	ir.CurFunc = nil
	fn.Body = []ir.Node{innerCall}

	// We're all done with the captured dictionary (and receiver, for method values).
	ir.FinishCaptureNames(pos, outer, fn)

	// Make a closure referencing our new internal function.
	c := ir.UseClosure(fn.OClosure, g.target)
	var init []ir.Node
	if outer != nil {
		init = append(init, dictAssign)
	}
	if rcvrValue != nil {
		init = append(init, rcvrAssign)
	}
	return ir.InitExpr(init, c)
}

// instantiateMethods instantiates all the methods (and associated dictionaries) of
// all fully-instantiated generic types that have been added to g.instTypeList.
func (g *irgen) instantiateMethods() {
	for i := 0; i < len(g.instTypeList); i++ {
		typ := g.instTypeList[i]
		if typ.HasShape() {
			// Shape types should not have any methods.
			continue
		}
		// Mark runtime type as needed, since this ensures that the
		// compiler puts out the needed DWARF symbols, when this
		// instantiated type has a different package from the local
		// package.
		typecheck.NeedRuntimeType(typ)
		// Lookup the method on the base generic type, since methods may
		// not be set on imported instantiated types.
		baseSym := typ.OrigSym
		baseType := baseSym.Def.(*ir.Name).Type()
		for j, _ := range typ.Methods().Slice() {
			baseNname := baseType.Methods().Slice()[j].Nname.(*ir.Name)
			// Eagerly generate the instantiations and dictionaries that implement these methods.
			// We don't use the instantiations here, just generate them (and any
			// further instantiations those generate, etc.).
			// Note that we don't set the Func for any methods on instantiated
			// types. Their signatures don't match so that would be confusing.
			// Direct method calls go directly to the instantiations, implemented above.
			// Indirect method calls use wrappers generated in reflectcall. Those wrappers
			// will use these instantiations if they are needed (for interface tables or reflection).
			_ = g.getInstantiation(baseNname, typ.RParams(), true)
			_ = g.getDictionarySym(baseNname, typ.RParams(), true)
		}
	}
	g.instTypeList = nil

}

// getInstNameNode returns the name node for the method or function being instantiated, and a bool which is true if a method is being instantiated.
func (g *irgen) getInstNameNode(inst *ir.InstExpr) (*ir.Name, bool) {
	if meth, ok := inst.X.(*ir.SelectorExpr); ok {
		return meth.Selection.Nname.(*ir.Name), true
	} else {
		return inst.X.(*ir.Name), false
	}
}

// getDictOrSubdict returns, for a method/function call or reference (node n) in an
// instantiation (described by instInfo), a node which is accessing a sub-dictionary
// or main/static dictionary, as needed, and also returns a boolean indicating if a
// sub-dictionary was accessed. nameNode is the particular function or method being
// called/referenced, and targs are the type arguments.
func (g *irgen) getDictOrSubdict(declInfo *instInfo, n ir.Node, nameNode *ir.Name, targs []*types.Type, isMeth bool) (ir.Node, bool) {
	var dict ir.Node
	usingSubdict := false
	if declInfo != nil {
		// Get the dictionary arg via sub-dictionary reference
		entry, ok := declInfo.dictEntryMap[n]
		// If the entry is not found, it may be that this node did not have
		// any type args that depend on type params, so we need a main
		// dictionary, not a sub-dictionary.
		if ok {
			dict = getDictionaryEntry(n.Pos(), declInfo.dictParam, entry, declInfo.dictLen)
			usingSubdict = true
		}
	}
	if !usingSubdict {
		dict = g.getDictionaryValue(nameNode, targs, isMeth)
	}
	return dict, usingSubdict
}

func addGcType(fl []*types.Field, t *types.Type) []*types.Field {
	return append(fl, types.NewField(base.Pos, typecheck.Lookup("F"+strconv.Itoa(len(fl))), t))
}

const INTTYPE = types.TINT64   // XX fix for 32-bit arch
const UINTTYPE = types.TUINT64 // XX fix for 32-bit arch
const INTSTRING = "i8"         // XX fix for 32-bit arch
const UINTSTRING = "u8"        // XX fix for 32-bit arch

// accumGcshape adds fields to fl resulting from the GCshape transformation of
// type t. The string associated with the GCshape transformation of t is added to
// buf. fieldSym is the sym of the field associated with type t, if it is in a
// struct. fieldSym could be used to have special naming for blank fields, etc.
func accumGcshape(fl []*types.Field, buf *bytes.Buffer, t *types.Type, fieldSym *types.Sym) []*types.Field {
	// t.Kind() is already the kind of the underlying type, so no need to
	// reference t.Underlying() to reference the underlying type.
	assert(t.Kind() == t.Underlying().Kind())

	switch t.Kind() {
	case types.TINT8:
		fl = addGcType(fl, types.Types[types.TINT8])
		buf.WriteString("i1")

	case types.TUINT8:
		fl = addGcType(fl, types.Types[types.TUINT8])
		buf.WriteString("u1")

	case types.TINT16:
		fl = addGcType(fl, types.Types[types.TINT16])
		buf.WriteString("i2")

	case types.TUINT16:
		fl = addGcType(fl, types.Types[types.TUINT16])
		buf.WriteString("u2")

	case types.TINT32:
		fl = addGcType(fl, types.Types[types.TINT32])
		buf.WriteString("i4")

	case types.TUINT32:
		fl = addGcType(fl, types.Types[types.TUINT32])
		buf.WriteString("u4")

	case types.TINT64:
		fl = addGcType(fl, types.Types[types.TINT64])
		buf.WriteString("i8")

	case types.TUINT64:
		fl = addGcType(fl, types.Types[types.TUINT64])
		buf.WriteString("u8")

	case types.TINT:
		fl = addGcType(fl, types.Types[INTTYPE])
		buf.WriteString(INTSTRING)

	case types.TUINT, types.TUINTPTR:
		fl = addGcType(fl, types.Types[UINTTYPE])
		buf.WriteString(UINTSTRING)

	case types.TCOMPLEX64:
		fl = addGcType(fl, types.Types[types.TFLOAT32])
		fl = addGcType(fl, types.Types[types.TFLOAT32])
		buf.WriteString("f4")
		buf.WriteString("f4")

	case types.TCOMPLEX128:
		fl = addGcType(fl, types.Types[types.TFLOAT64])
		fl = addGcType(fl, types.Types[types.TFLOAT64])
		buf.WriteString("f8")
		buf.WriteString("f8")

	case types.TFLOAT32:
		fl = addGcType(fl, types.Types[types.TFLOAT32])
		buf.WriteString("f4")

	case types.TFLOAT64:
		fl = addGcType(fl, types.Types[types.TFLOAT64])
		buf.WriteString("f8")

	case types.TBOOL:
		fl = addGcType(fl, types.Types[types.TINT8])
		buf.WriteString("i1")

	case types.TPTR:
		fl = addGcType(fl, types.Types[types.TUNSAFEPTR])
		buf.WriteString("p")

	case types.TFUNC:
		fl = addGcType(fl, types.Types[types.TUNSAFEPTR])
		buf.WriteString("p")

	case types.TSLICE:
		fl = addGcType(fl, types.Types[types.TUNSAFEPTR])
		fl = addGcType(fl, types.Types[INTTYPE])
		fl = addGcType(fl, types.Types[INTTYPE])
		buf.WriteString("p")
		buf.WriteString(INTSTRING)
		buf.WriteString(INTSTRING)

	case types.TARRAY:
		n := t.NumElem()
		if n == 1 {
			fl = accumGcshape(fl, buf, t.Elem(), nil)
		} else if n > 0 {
			// Represent an array with more than one element as its
			// unique type, since it must be treated differently for
			// regabi.
			fl = addGcType(fl, t)
			buf.WriteByte('[')
			buf.WriteString(strconv.Itoa(int(n)))
			buf.WriteString("](")
			var ignore []*types.Field
			// But to determine its gcshape name, we must call
			// accumGcShape() on t.Elem().
			accumGcshape(ignore, buf, t.Elem(), nil)
			buf.WriteByte(')')
		}

	case types.TSTRUCT:
		nfields := t.NumFields()
		for i, f := range t.Fields().Slice() {
			fl = accumGcshape(fl, buf, f.Type, f.Sym)

			// Check if we need to add an alignment field.
			var pad int64
			if i < nfields-1 {
				pad = t.Field(i+1).Offset - f.Offset - f.Type.Width
			} else {
				pad = t.Width - f.Offset - f.Type.Width
			}
			if pad > 0 {
				// There is padding between fields or at end of
				// struct. Add an alignment field.
				fl = addGcType(fl, types.NewArray(types.Types[types.TUINT8], pad))
				buf.WriteString("a")
				buf.WriteString(strconv.Itoa(int(pad)))
			}
		}

	case types.TCHAN:
		fl = addGcType(fl, types.Types[types.TUNSAFEPTR])
		buf.WriteString("p")

	case types.TMAP:
		fl = addGcType(fl, types.Types[types.TUNSAFEPTR])
		buf.WriteString("p")

	case types.TINTER:
		fl = addGcType(fl, types.Types[types.TUNSAFEPTR])
		fl = addGcType(fl, types.Types[types.TUNSAFEPTR])
		buf.WriteString("pp")

	case types.TFORW, types.TANY:
		assert(false)

	case types.TSTRING:
		fl = addGcType(fl, types.Types[types.TUNSAFEPTR])
		fl = addGcType(fl, types.Types[INTTYPE])
		buf.WriteString("p")
		buf.WriteString(INTSTRING)

	case types.TUNSAFEPTR:
		fl = addGcType(fl, types.Types[types.TUNSAFEPTR])
		buf.WriteString("p")

	default: // Everything TTYPEPARAM and below in list of Kinds
		assert(false)
	}

	return fl
}

// gcshapeType returns the GCshape type and name corresponding to type t.
func gcshapeType(t *types.Type) (*types.Type, string) {
	var fl []*types.Field
	buf := bytes.NewBufferString("")

	// Call CallSize so type sizes and field offsets are available.
	types.CalcSize(t)

	instType := t.Sym() != nil && t.IsFullyInstantiated()
	if instType {
		// We distinguish the gcshape of all top-level instantiated type from
		// normal concrete types, even if they have the exact same underlying
		// "shape", because in a function instantiation, any method call on
		// this type arg will be a generic method call (requiring a
		// dictionary), rather than a direct method call on the underlying
		// type (no dictionary). So, we add the instshape prefix to the
		// normal gcshape name, and will make it a defined type with that
		// name below.
		buf.WriteString("instshape-")
	}
	fl = accumGcshape(fl, buf, t, nil)

	// TODO: Should gcshapes be in a global package, so we don't have to
	// duplicate in each package? Or at least in the specified source package
	// of a function/method instantiation?
	gcshape := types.NewStruct(types.LocalPkg, fl)
	gcname := buf.String()
	if instType {
		// Lookup or create type with name 'gcname' (with instshape prefix).
		newsym := t.Sym().Pkg.Lookup(gcname)
		if newsym.Def != nil {
			gcshape = newsym.Def.Type()
		} else {
			newt := typecheck.NewIncompleteNamedType(t.Pos(), newsym)
			newt.SetUnderlying(gcshape.Underlying())
			gcshape = newt
		}
	}
	assert(gcshape.Size() == t.Size())
	return gcshape, buf.String()
}

// checkFetchBody checks if a generic body can be fetched, but hasn't been loaded
// yet. If so, it imports the body.
func checkFetchBody(nameNode *ir.Name) {
	if nameNode.Func.Body == nil && nameNode.Func.Inl != nil {
		// If there is no body yet but Func.Inl exists, then we can can
		// import the whole generic body.
		assert(nameNode.Func.Inl.Cost == 1 && nameNode.Sym().Pkg != types.LocalPkg)
		typecheck.ImportBody(nameNode.Func)
		assert(nameNode.Func.Inl.Body != nil)
		nameNode.Func.Body = nameNode.Func.Inl.Body
		nameNode.Func.Dcl = nameNode.Func.Inl.Dcl
	}
}

// getInstantiation gets the instantiantion and dictionary of the function or method nameNode
// with the type arguments targs. If the instantiated function is not already
// cached, then it calls genericSubst to create the new instantiation.
func (g *irgen) getInstantiation(nameNode *ir.Name, targs []*types.Type, isMeth bool) *ir.Func {
	checkFetchBody(nameNode)

	// Convert type arguments to their shape, so we can reduce the number
	// of instantiations we have to generate.
	shapes := typecheck.ShapifyList(targs)

	sym := typecheck.MakeInstName(nameNode.Sym(), shapes, isMeth)
	info := g.instInfoMap[sym]
	if info == nil {
		if false {
			// Testing out gcshapeType() and gcshapeName()
			for i, t := range targs {
				gct, gcs := gcshapeType(t)
				fmt.Printf("targ %d: %v %v %v\n", i, gcs, gct, gct.Underlying())
			}
		}
		// If instantiation doesn't exist yet, create it and add
		// to the list of decls.
		gfInfo := g.getGfInfo(nameNode)
		info = &instInfo{
			gf:           nameNode,
			gfInfo:       gfInfo,
			startSubDict: len(targs) + len(gfInfo.derivedTypes),
			dictLen:      len(targs) + len(gfInfo.derivedTypes) + len(gfInfo.subDictCalls),
			dictEntryMap: make(map[ir.Node]int),
		}
		// genericSubst fills in info.dictParam and info.dictEntryMap.
		st := g.genericSubst(sym, nameNode, shapes, targs, isMeth, info)
		info.fun = st
		g.instInfoMap[sym] = info
		// This ensures that the linker drops duplicates of this instantiation.
		// All just works!
		st.SetDupok(true)
		g.target.Decls = append(g.target.Decls, st)
		if base.Flag.W > 1 {
			ir.Dump(fmt.Sprintf("\nstenciled %v", st), st)
		}
	}
	return info.fun
}

// Struct containing info needed for doing the substitution as we create the
// instantiation of a generic function with specified type arguments.
type subster struct {
	g        *irgen
	isMethod bool     // If a method is being instantiated
	newf     *ir.Func // Func node for the new stenciled function
	ts       typecheck.Tsubster
	info     *instInfo // Place to put extra info in the instantiation

	// Which type parameter the shape type came from.
	shape2param map[*types.Type]*types.Type

	// unshapeify maps from shape types to the concrete types they represent.
	// TODO: remove when we no longer need it.
	unshapify  typecheck.Tsubster
	concretify typecheck.Tsubster

	// TODO: some sort of map from <shape type, interface type> to index in the
	// dictionary where a *runtime.itab for the corresponding <concrete type,
	// interface type> pair resides.
}

// genericSubst returns a new function with name newsym. The function is an
// instantiation of a generic function or method specified by namedNode with type
// args targs. For a method with a generic receiver, it returns an instantiated
// function type where the receiver becomes the first parameter. Otherwise the
// instantiated method would still need to be transformed by later compiler
// phases.  genericSubst fills in info.dictParam and info.dictEntryMap.
func (g *irgen) genericSubst(newsym *types.Sym, nameNode *ir.Name, shapes, targs []*types.Type, isMethod bool, info *instInfo) *ir.Func {
	var tparams []*types.Type
	if isMethod {
		// Get the type params from the method receiver (after skipping
		// over any pointer)
		recvType := nameNode.Type().Recv().Type
		recvType = deref(recvType)
		tparams = recvType.RParams()
	} else {
		fields := nameNode.Type().TParams().Fields().Slice()
		tparams = make([]*types.Type, len(fields))
		for i, f := range fields {
			tparams[i] = f.Type
		}
	}
	for i := range targs {
		if targs[i].HasShape() {
			base.Fatalf("generiSubst shape %s %+v %+v\n", newsym.Name, shapes[i], targs[i])
		}
	}
	gf := nameNode.Func
	// Pos of the instantiated function is same as the generic function
	newf := ir.NewFunc(gf.Pos())
	newf.Pragma = gf.Pragma // copy over pragmas from generic function to stenciled implementation.
	newf.Nname = ir.NewNameAt(gf.Pos(), newsym)
	newf.Nname.Func = newf
	newf.Nname.Defn = newf
	newsym.Def = newf.Nname
	savef := ir.CurFunc
	// transformCall/transformReturn (called during stenciling of the body)
	// depend on ir.CurFunc being set.
	ir.CurFunc = newf

	assert(len(tparams) == len(shapes))
	assert(len(tparams) == len(targs))

	subst := &subster{
		g:        g,
		isMethod: isMethod,
		newf:     newf,
		info:     info,
		ts: typecheck.Tsubster{
			Tparams: tparams,
			Targs:   shapes,
			Vars:    make(map[*ir.Name]*ir.Name),
		},
		shape2param: map[*types.Type]*types.Type{},
		unshapify: typecheck.Tsubster{
			Tparams: shapes,
			Targs:   targs,
			Vars:    make(map[*ir.Name]*ir.Name),
		},
		concretify: typecheck.Tsubster{
			Tparams: tparams,
			Targs:   targs,
			Vars:    make(map[*ir.Name]*ir.Name),
		},
	}
	for i := range shapes {
		if !shapes[i].IsShape() {
			panic("must be a shape type")
		}
		subst.shape2param[shapes[i]] = tparams[i]
	}

	newf.Dcl = make([]*ir.Name, 0, len(gf.Dcl)+1)

	// Create the needed dictionary param
	dictionarySym := newsym.Pkg.Lookup(".dict")
	dictionaryType := types.Types[types.TUINTPTR]
	dictionaryName := ir.NewNameAt(gf.Pos(), dictionarySym)
	typed(dictionaryType, dictionaryName)
	dictionaryName.Class = ir.PPARAM
	dictionaryName.Curfn = newf
	newf.Dcl = append(newf.Dcl, dictionaryName)
	for _, n := range gf.Dcl {
		if n.Sym().Name == ".dict" {
			panic("already has dictionary")
		}
		newf.Dcl = append(newf.Dcl, subst.localvar(n))
	}
	dictionaryArg := types.NewField(gf.Pos(), dictionarySym, dictionaryType)
	dictionaryArg.Nname = dictionaryName
	info.dictParam = dictionaryName

	// We add the dictionary as the first parameter in the function signature.
	// We also transform a method type to the corresponding function type
	// (make the receiver be the next parameter after the dictionary).
	oldt := nameNode.Type()
	var args []*types.Field
	args = append(args, dictionaryArg)
	args = append(args, oldt.Recvs().FieldSlice()...)
	args = append(args, oldt.Params().FieldSlice()...)

	// Replace the types in the function signature via subst.fields.
	// Ugly: also, we have to insert the Name nodes of the parameters/results into
	// the function type. The current function type has no Nname fields set,
	// because it came via conversion from the types2 type.
	newt := types.NewSignature(oldt.Pkg(), nil, nil,
		subst.fields(ir.PPARAM, args, newf.Dcl),
		subst.fields(ir.PPARAMOUT, oldt.Results().FieldSlice(), newf.Dcl))

	typed(newt, newf.Nname)
	ir.MarkFunc(newf.Nname)
	newf.SetTypecheck(1)

	// Make sure name/type of newf is set before substituting the body.
	newf.Body = subst.list(gf.Body)

	// Add code to check that the dictionary is correct.
	// TODO: must go away when we move to many->1 shape to concrete mapping.
	newf.Body.Prepend(subst.checkDictionary(dictionaryName, targs)...)

	ir.CurFunc = savef
	// Add any new, fully instantiated types seen during the substitution to
	// g.instTypeList.
	g.instTypeList = append(g.instTypeList, subst.ts.InstTypeList...)
	g.instTypeList = append(g.instTypeList, subst.unshapify.InstTypeList...)
	g.instTypeList = append(g.instTypeList, subst.concretify.InstTypeList...)

	return newf
}

func (subst *subster) unshapifyTyp(t *types.Type) *types.Type {
	res := subst.unshapify.Typ(t)
	types.CheckSize(res)
	return res
}

// localvar creates a new name node for the specified local variable and enters it
// in subst.vars. It substitutes type arguments for type parameters in the type of
// name as needed.
func (subst *subster) localvar(name *ir.Name) *ir.Name {
	m := ir.NewNameAt(name.Pos(), name.Sym())
	if name.IsClosureVar() {
		m.SetIsClosureVar(true)
	}
	m.SetType(subst.ts.Typ(name.Type()))
	m.BuiltinOp = name.BuiltinOp
	m.Curfn = subst.newf
	m.Class = name.Class
	assert(name.Class != ir.PEXTERN && name.Class != ir.PFUNC)
	m.Func = name.Func
	subst.ts.Vars[name] = m
	m.SetTypecheck(1)
	return m
}

// checkDictionary returns code that does runtime consistency checks
// between the dictionary and the types it should contain.
func (subst *subster) checkDictionary(name *ir.Name, targs []*types.Type) (code []ir.Node) {
	if false {
		return // checking turned off
	}
	// TODO: when moving to GCshape, this test will become harder. Call into
	// runtime to check the expected shape is correct?
	pos := name.Pos()
	// Convert dictionary to *[N]uintptr
	d := ir.NewConvExpr(pos, ir.OCONVNOP, types.Types[types.TUNSAFEPTR], name)
	d.SetTypecheck(1)
	d = ir.NewConvExpr(pos, ir.OCONVNOP, types.NewArray(types.Types[types.TUINTPTR], int64(len(targs))).PtrTo(), d)
	d.SetTypecheck(1)

	// Check that each type entry in the dictionary is correct.
	for i, t := range targs {
		if t.HasShape() {
			// Check the concrete type, not the shape type.
			// TODO: can this happen?
			//t = subst.unshapify.Typ(t)
			base.Fatalf("shape type in dictionary %s %+v\n", name.Sym().Name, t)
			continue
		}
		want := reflectdata.TypePtr(t)
		typed(types.Types[types.TUINTPTR], want)
		deref := ir.NewStarExpr(pos, d)
		typed(d.Type().Elem(), deref)
		idx := ir.NewConstExpr(constant.MakeUint64(uint64(i)), name) // TODO: what to set orig to?
		typed(types.Types[types.TUINTPTR], idx)
		got := ir.NewIndexExpr(pos, deref, idx)
		typed(types.Types[types.TUINTPTR], got)
		cond := ir.NewBinaryExpr(pos, ir.ONE, want, got)
		typed(types.Types[types.TBOOL], cond)
		panicArg := ir.NewNilExpr(pos)
		typed(types.NewInterface(types.LocalPkg, nil), panicArg)
		then := ir.NewUnaryExpr(pos, ir.OPANIC, panicArg)
		then.SetTypecheck(1)
		x := ir.NewIfStmt(pos, cond, []ir.Node{then}, nil)
		x.SetTypecheck(1)
		code = append(code, x)
	}
	return
}

// getDictionaryEntry gets the i'th entry in the dictionary dict.
func getDictionaryEntry(pos src.XPos, dict *ir.Name, i int, size int) ir.Node {
	// Convert dictionary to *[N]uintptr
	// All entries in the dictionary are pointers. They all point to static data, though, so we
	// treat them as uintptrs so the GC doesn't need to keep track of them.
	d := ir.NewConvExpr(pos, ir.OCONVNOP, types.Types[types.TUNSAFEPTR], dict)
	d.SetTypecheck(1)
	d = ir.NewConvExpr(pos, ir.OCONVNOP, types.NewArray(types.Types[types.TUINTPTR], int64(size)).PtrTo(), d)
	d.SetTypecheck(1)

	// Load entry i out of the dictionary.
	deref := ir.NewStarExpr(pos, d)
	typed(d.Type().Elem(), deref)
	idx := ir.NewConstExpr(constant.MakeUint64(uint64(i)), dict) // TODO: what to set orig to?
	typed(types.Types[types.TUINTPTR], idx)
	r := ir.NewIndexExpr(pos, deref, idx)
	typed(types.Types[types.TUINTPTR], r)
	return r
}

// getDictionaryType returns a *runtime._type from the dictionary entry i
// (which refers to a type param or a derived type that uses type params).
func (subst *subster) getDictionaryType(pos src.XPos, i int) ir.Node {
	if i < 0 || i >= subst.info.startSubDict {
		base.Fatalf(fmt.Sprintf("bad dict index %d", i))
	}

	r := getDictionaryEntry(pos, subst.info.dictParam, i, subst.info.startSubDict)
	// change type of retrieved dictionary entry to *byte, which is the
	// standard typing of a *runtime._type in the compiler
	typed(types.Types[types.TUINT8].PtrTo(), r)
	return r
}

// node is like DeepCopy(), but substitutes ONAME nodes based on subst.ts.vars, and
// also descends into closures. It substitutes type arguments for type parameters
// in all the new nodes.
func (subst *subster) node(n ir.Node) ir.Node {
	// Use closure to capture all state needed by the ir.EditChildren argument.
	var edit func(ir.Node) ir.Node
	edit = func(x ir.Node) ir.Node {
		switch x.Op() {
		case ir.OTYPE:
			return ir.TypeNode(subst.ts.Typ(x.Type()))

		case ir.ONAME:
			if v := subst.ts.Vars[x.(*ir.Name)]; v != nil {
				return v
			}
			return x
		case ir.ONONAME:
			// This handles the identifier in a type switch guard
			fallthrough
		case ir.OLITERAL, ir.ONIL:
			if x.Sym() != nil {
				return x
			}
		}
		m := ir.Copy(x)
		if _, isExpr := m.(ir.Expr); isExpr {
			t := x.Type()
			if t == nil {
				// t can be nil only if this is a call that has no
				// return values, so allow that and otherwise give
				// an error.
				_, isCallExpr := m.(*ir.CallExpr)
				_, isStructKeyExpr := m.(*ir.StructKeyExpr)
				_, isKeyExpr := m.(*ir.KeyExpr)
				if !isCallExpr && !isStructKeyExpr && !isKeyExpr && x.Op() != ir.OPANIC &&
					x.Op() != ir.OCLOSE {
					base.Fatalf(fmt.Sprintf("Nil type for %v", x))
				}
			} else if x.Op() != ir.OCLOSURE {
				m.SetType(subst.ts.Typ(x.Type()))
			}
		}

		for i, de := range subst.info.gfInfo.subDictCalls {
			if de == x {
				// Remember the dictionary entry associated with this
				// node in the instantiated function
				// TODO: make sure this remains correct with respect to the
				// transformations below.
				subst.info.dictEntryMap[m] = subst.info.startSubDict + i
				break
			}
		}

		ir.EditChildren(m, edit)

		m.SetTypecheck(1)
		if typecheck.IsCmp(x.Op()) {
			transformCompare(m.(*ir.BinaryExpr))
		} else {
			switch x.Op() {
			case ir.OSLICE, ir.OSLICE3:
				transformSlice(m.(*ir.SliceExpr))

			case ir.OADD:
				m = transformAdd(m.(*ir.BinaryExpr))

			case ir.OINDEX:
				transformIndex(m.(*ir.IndexExpr))

			case ir.OAS2:
				as2 := m.(*ir.AssignListStmt)
				transformAssign(as2, as2.Lhs, as2.Rhs)

			case ir.OAS:
				as := m.(*ir.AssignStmt)
				if as.Y != nil {
					// transformAssign doesn't handle the case
					// of zeroing assignment of a dcl (rhs[0] is nil).
					lhs, rhs := []ir.Node{as.X}, []ir.Node{as.Y}
					transformAssign(as, lhs, rhs)
				}

			case ir.OASOP:
				as := m.(*ir.AssignOpStmt)
				transformCheckAssign(as, as.X)

			case ir.ORETURN:
				transformReturn(m.(*ir.ReturnStmt))

			case ir.OSEND:
				transformSend(m.(*ir.SendStmt))

			}
		}

		switch x.Op() {
		case ir.OLITERAL:
			t := m.Type()
			if t != x.Type() {
				// types2 will give us a constant with a type T,
				// if an untyped constant is used with another
				// operand of type T (in a provably correct way).
				// When we substitute in the type args during
				// stenciling, we now know the real type of the
				// constant. We may then need to change the
				// BasicLit.val to be the correct type (e.g.
				// convert an int64Val constant to a floatVal
				// constant).
				m.SetType(types.UntypedInt) // use any untyped type for DefaultLit to work
				m = typecheck.DefaultLit(m, t)
			}

		case ir.OXDOT:
			// A method value/call via a type param will have been
			// left as an OXDOT. When we see this during stenciling,
			// finish the transformation, now that we have the
			// instantiated receiver type. We need to do this now,
			// since the access/selection to the method for the real
			// type is very different from the selection for the type
			// param. m will be transformed to an OMETHVALUE node. It
			// will be transformed to an ODOTMETH or ODOTINTER node if
			// we find in the OCALL case below that the method value
			// is actually called.
			mse := m.(*ir.SelectorExpr)
			if src := mse.X.Type(); src.IsShape() {
				// The only dot on a shape type value are methods.
				if mse.X.Op() == ir.OTYPE {
					// Method expression T.M
					// Fall back from shape type to concrete type.
					src = subst.unshapifyTyp(src)
					mse.X = ir.TypeNode(src)
				} else {
					// Implement x.M as a conversion-to-bound-interface
					//  1) convert x to the bound interface
					//  2) call M on that interface
					dst := subst.concretify.Typ(subst.shape2param[src].Bound())
					// Mark that we use the methods of this concrete type.
					// Otherwise the linker deadcode-eliminates them :(
					ix := subst.findDictType(subst.shape2param[src])
					assert(ix >= 0)
					mse.X = subst.convertUsingDictionary(m.Pos(), mse.X, dst, subst.shape2param[src], ix)
				}
			}
			transformDot(mse, false)
			if mse.Op() == ir.OMETHEXPR && mse.X.Type().HasShape() {
				mse.X = ir.TypeNodeAt(mse.X.Pos(), subst.unshapifyTyp(mse.X.Type()))
			}
			m.SetTypecheck(1)

		case ir.OCALL:
			call := m.(*ir.CallExpr)
			convcheck := false
			switch call.X.Op() {
			case ir.OTYPE:
				// Transform the conversion, now that we know the
				// type argument.
				m = transformConvCall(call)
				if m.Op() == ir.OCONVIFACE {
					// Note: srcType uses x.Args[0], not m.X or call.Args[0], because
					// we need the type before the type parameter -> type argument substitution.
					srcType := x.(*ir.CallExpr).Args[0].Type()
					if ix := subst.findDictType(srcType); ix >= 0 {
						c := m.(*ir.ConvExpr)
						m = subst.convertUsingDictionary(c.Pos(), c.X, c.Type(), srcType, ix)
					}
				}

			case ir.OMETHVALUE, ir.OMETHEXPR:
				// Redo the transformation of OXDOT, now that we
				// know the method value is being called. Then
				// transform the call.
				call.X.(*ir.SelectorExpr).SetOp(ir.OXDOT)
				transformDot(call.X.(*ir.SelectorExpr), true)
				call.X.SetType(subst.unshapifyTyp(call.X.Type()))
				transformCall(call)
				convcheck = true

			case ir.ODOT, ir.ODOTPTR:
				// An OXDOT for a generic receiver was resolved to
				// an access to a field which has a function
				// value. Transform the call to that function, now
				// that the OXDOT was resolved.
				transformCall(call)
				convcheck = true

			case ir.ONAME:
				name := call.X.Name()
				if name.BuiltinOp != ir.OXXX {
					switch name.BuiltinOp {
					case ir.OMAKE, ir.OREAL, ir.OIMAG, ir.OLEN, ir.OCAP, ir.OAPPEND:
						// Transform these builtins now that we
						// know the type of the args.
						m = transformBuiltin(call)
					default:
						base.FatalfAt(call.Pos(), "Unexpected builtin op")
					}
					switch m.Op() {
					case ir.OAPPEND:
						// Append needs to pass a concrete type to the runtime.
						// TODO: there's no way to record a dictionary-loaded type for walk to use here
						m.SetType(subst.unshapifyTyp(m.Type()))
					}

				} else {
					// This is the case of a function value that was a
					// type parameter (implied to be a function via a
					// structural constraint) which is now resolved.
					transformCall(call)
					convcheck = true
				}

			case ir.OCLOSURE:
				transformCall(call)
				convcheck = true

			case ir.OFUNCINST:
				// A call with an OFUNCINST will get transformed
				// in stencil() once we have created & attached the
				// instantiation to be called.

			default:
				base.FatalfAt(call.Pos(), fmt.Sprintf("Unexpected op with CALL during stenciling: %v", call.X.Op()))
			}
			if convcheck {
				for i, arg := range x.(*ir.CallExpr).Args {
					if arg.Type().HasTParam() && arg.Op() != ir.OCONVIFACE &&
						call.Args[i].Op() == ir.OCONVIFACE {
						ix := subst.findDictType(arg.Type())
						assert(ix >= 0)
						call.Args[i] = subst.convertUsingDictionary(arg.Pos(), call.Args[i].(*ir.ConvExpr).X, call.Args[i].Type(), arg.Type(), ix)
					}
				}
			}

		case ir.OCLOSURE:
			// We're going to create a new closure from scratch, so clear m
			// to avoid using the ir.Copy by accident until we reassign it.
			m = nil

			x := x.(*ir.ClosureExpr)
			// Need to duplicate x.Func.Nname, x.Func.Dcl, x.Func.ClosureVars, and
			// x.Func.Body.
			oldfn := x.Func
			newfn := ir.NewClosureFunc(oldfn.Pos(), subst.newf != nil)
			ir.NameClosure(newfn.OClosure, subst.newf)

			saveNewf := subst.newf
			ir.CurFunc = newfn
			subst.newf = newfn
			newfn.Dcl = subst.namelist(oldfn.Dcl)

			// Make a closure variable for the dictionary of the
			// containing function.
			cdict := ir.CaptureName(oldfn.Pos(), newfn, subst.info.dictParam)
			typed(types.Types[types.TUINTPTR], cdict)
			ir.FinishCaptureNames(oldfn.Pos(), saveNewf, newfn)
			newfn.ClosureVars = append(newfn.ClosureVars, subst.namelist(oldfn.ClosureVars)...)

			// Create inst info for the instantiated closure. The dict
			// param is the closure variable for the dictionary of the
			// outer function. Since the dictionary is shared, use the
			// same entries for startSubDict, dictLen, dictEntryMap.
			cinfo := &instInfo{
				fun:          newfn,
				dictParam:    cdict,
				startSubDict: subst.info.startSubDict,
				dictLen:      subst.info.dictLen,
				dictEntryMap: subst.info.dictEntryMap,
			}
			subst.g.instInfoMap[newfn.Nname.Sym()] = cinfo

			typed(subst.ts.Typ(oldfn.Nname.Type()), newfn.Nname)
			typed(newfn.Nname.Type(), newfn.OClosure)
			newfn.SetTypecheck(1)

			// Make sure type of closure function is set before doing body.
			newfn.Body = subst.list(oldfn.Body)
			subst.newf = saveNewf
			ir.CurFunc = saveNewf

			m = ir.UseClosure(newfn.OClosure, subst.g.target)
			m.(*ir.ClosureExpr).SetInit(subst.list(x.Init()))

		case ir.OCONVIFACE:
			x := x.(*ir.ConvExpr)
			// Note: x's argument is still typed as a type parameter.
			// m's argument now has an instantiated type.
			t := x.X.Type()
			if ix := subst.findDictType(t); ix >= 0 {
				m = subst.convertUsingDictionary(x.Pos(), m.(*ir.ConvExpr).X, m.Type(), t, ix)
			}
		case ir.OEQ, ir.ONE:
			// Equality between a non-interface and an interface requires the non-interface
			// to be promoted to an interface.
			x := x.(*ir.BinaryExpr)
			m := m.(*ir.BinaryExpr)
			if i := x.Y.Type(); i.IsInterface() {
				if ix := subst.findDictType(x.X.Type()); ix >= 0 {
					m.X = subst.convertUsingDictionary(m.X.Pos(), m.X, i, x.X.Type(), ix)
				}
			}
			if i := x.X.Type(); i.IsInterface() {
				if ix := subst.findDictType(x.Y.Type()); ix >= 0 {
					m.Y = subst.convertUsingDictionary(m.Y.Pos(), m.Y, i, x.X.Type(), ix)
				}
			}

		case ir.ONEW:
			// New needs to pass a concrete type to the runtime.
			// Or maybe it doesn't? We could use a shape type.
			// TODO: need to modify m.X? I don't think any downstream passes use it.
			m.SetType(subst.unshapifyTyp(m.Type()))

		case ir.OPTRLIT:
			m := m.(*ir.AddrExpr)
			// Walk uses the type of the argument of ptrlit. Also could be a shape type?
			m.X.SetType(subst.unshapifyTyp(m.X.Type()))

		case ir.OMETHEXPR:
			se := m.(*ir.SelectorExpr)
			se.X = ir.TypeNodeAt(se.X.Pos(), subst.unshapifyTyp(se.X.Type()))
		case ir.OFUNCINST:
			inst := m.(*ir.InstExpr)
			targs2 := make([]ir.Node, len(inst.Targs))
			for i, n := range inst.Targs {
				targs2[i] = ir.TypeNodeAt(n.Pos(), subst.unshapifyTyp(n.Type()))
				// TODO: need an ir.Name node?
			}
			inst.Targs = targs2
		}
		return m
	}

	return edit(n)
}

// findDictType looks for type t in the typeparams or derived types in the generic
// function info subst.info.gfInfo. This will indicate the dictionary entry with the
// correct concrete type for the associated instantiated function.
func (subst *subster) findDictType(t *types.Type) int {
	for i, dt := range subst.info.gfInfo.tparams {
		if dt == t {
			return i
		}
	}
	for i, dt := range subst.info.gfInfo.derivedTypes {
		if types.Identical(dt, t) {
			return i + len(subst.info.gfInfo.tparams)
		}
	}
	return -1
}

// convertUsingDictionary converts value v from instantiated type src (which is index
// 'ix' in the instantiation's dictionary) to an interface type dst.
func (subst *subster) convertUsingDictionary(pos src.XPos, v ir.Node, dst, src *types.Type, ix int) ir.Node {
	if !dst.IsInterface() {
		base.Fatalf("can only convert type parameters to interfaces %+v -> %+v", src, dst)
	}
	// Load the actual runtime._type of the type parameter from the dictionary.
	rt := subst.getDictionaryType(pos, ix)

	// Convert value to an interface type, so the data field is what we want.
	if !v.Type().IsInterface() {
		v = ir.NewConvExpr(v.Pos(), ir.OCONVIFACE, nil, v)
		typed(types.NewInterface(types.LocalPkg, nil), v)
	}

	// At this point, v is an interface type with a data word we want.
	// But the type word represents a gcshape type, which we don't want.
	// Replace with the instantiated type loaded from the dictionary.
	data := ir.NewUnaryExpr(pos, ir.OIDATA, v)
	typed(types.Types[types.TUNSAFEPTR], data)
	var i ir.Node = ir.NewBinaryExpr(pos, ir.OEFACE, rt, data)
	if !dst.IsEmptyInterface() {
		// We just built an empty interface{}. Type it as such,
		// then assert it to the required non-empty interface.
		typed(types.NewInterface(types.LocalPkg, nil), i)
		i = ir.NewTypeAssertExpr(pos, i, nil)
	}
	typed(dst, i)
	// TODO: we're throwing away the type word of the original version
	// of m here (it would be OITAB(m)), which probably took some
	// work to generate. Can we avoid generating it at all?
	// (The linker will throw them away if not needed, so it would just
	// save toolchain work, not binary size.)
	return i

}

func (subst *subster) namelist(l []*ir.Name) []*ir.Name {
	s := make([]*ir.Name, len(l))
	for i, n := range l {
		s[i] = subst.localvar(n)
		if n.Defn != nil {
			s[i].Defn = subst.node(n.Defn)
		}
		if n.Outer != nil {
			s[i].Outer = subst.node(n.Outer).(*ir.Name)
		}
	}
	return s
}

func (subst *subster) list(l []ir.Node) []ir.Node {
	s := make([]ir.Node, len(l))
	for i, n := range l {
		s[i] = subst.node(n)
	}
	return s
}

// fields sets the Nname field for the Field nodes inside a type signature, based
// on the corresponding in/out parameters in dcl. It depends on the in and out
// parameters being in order in dcl.
func (subst *subster) fields(class ir.Class, oldfields []*types.Field, dcl []*ir.Name) []*types.Field {
	// Find the starting index in dcl of declarations of the class (either
	// PPARAM or PPARAMOUT).
	var i int
	for i = range dcl {
		if dcl[i].Class == class {
			break
		}
	}

	// Create newfields nodes that are copies of the oldfields nodes, but
	// with substitution for any type params, and with Nname set to be the node in
	// Dcl for the corresponding PPARAM or PPARAMOUT.
	newfields := make([]*types.Field, len(oldfields))
	for j := range oldfields {
		newfields[j] = oldfields[j].Copy()
		newfields[j].Type = subst.ts.Typ(oldfields[j].Type)
		// A PPARAM field will be missing from dcl if its name is
		// unspecified or specified as "_". So, we compare the dcl sym
		// with the field sym (or sym of the field's Nname node). (Unnamed
		// results still have a name like ~r2 in their Nname node.) If
		// they don't match, this dcl (if there is one left) must apply to
		// a later field.
		if i < len(dcl) && (dcl[i].Sym() == oldfields[j].Sym ||
			(oldfields[j].Nname != nil && dcl[i].Sym() == oldfields[j].Nname.Sym())) {
			newfields[j].Nname = dcl[i]
			i++
		}
	}
	return newfields
}

// deref does a single deref of type t, if it is a pointer type.
func deref(t *types.Type) *types.Type {
	if t.IsPtr() {
		return t.Elem()
	}
	return t
}

// getDictionarySym returns the dictionary for the named generic function gf, which
// is instantiated with the type arguments targs.
func (g *irgen) getDictionarySym(gf *ir.Name, targs []*types.Type, isMeth bool) *types.Sym {
	if len(targs) == 0 {
		base.Fatalf("%s should have type arguments", gf.Sym().Name)
	}

	// Enforce that only concrete types can make it to here.
	for _, t := range targs {
		if t.IsShape() {
			panic(fmt.Sprintf("shape %+v in dictionary for %s", t, gf.Sym().Name))
		}
	}

	// Get a symbol representing the dictionary.
	sym := typecheck.MakeDictName(gf.Sym(), targs, isMeth)

	// Initialize the dictionary, if we haven't yet already.
	if lsym := sym.Linksym(); len(lsym.P) == 0 {
		info := g.getGfInfo(gf)

		infoPrint("=== Creating dictionary %v\n", sym.Name)
		off := 0
		// Emit an entry for each targ (concrete type or gcshape).
		for _, t := range targs {
			infoPrint(" * %v\n", t)
			s := reflectdata.TypeLinksym(t)
			off = objw.SymPtr(lsym, off, s, 0)
			// Ensure that methods on t don't get deadcode eliminated
			// by the linker.
			// TODO: This is somewhat overkill, we really only need it
			// for types that are put into interfaces.
			reflectdata.MarkTypeUsedInInterface(t, lsym)
		}
		subst := typecheck.Tsubster{
			Tparams: info.tparams,
			Targs:   targs,
		}
		// Emit an entry for each derived type (after substituting targs)
		for _, t := range info.derivedTypes {
			ts := subst.Typ(t)
			infoPrint(" - %v\n", ts)
			s := reflectdata.TypeLinksym(ts)
			off = objw.SymPtr(lsym, off, s, 0)
			reflectdata.MarkTypeUsedInInterface(ts, lsym)
		}
		// Emit an entry for each subdictionary (after substituting targs)
		for _, n := range info.subDictCalls {
			var sym *types.Sym
			switch n.Op() {
			case ir.OCALL:
				call := n.(*ir.CallExpr)
				if call.X.Op() == ir.OXDOT {
					var nameNode *ir.Name
					se := call.X.(*ir.SelectorExpr)
					if types.IsInterfaceMethod(se.Selection.Type) {
						// This is a method call enabled by a type bound.
						tmpse := ir.NewSelectorExpr(base.Pos, ir.OXDOT, se.X, se.Sel)
						tmpse = typecheck.AddImplicitDots(tmpse)
						tparam := tmpse.X.Type()
						assert(tparam.IsTypeParam())
						recvType := targs[tparam.Index()]
						if len(recvType.RParams()) == 0 {
							// No sub-dictionary entry is
							// actually needed, since the
							// typeparam is not an
							// instantiated type that
							// will have generic methods.
							break
						}
						// This is a method call for an
						// instantiated type, so we need a
						// sub-dictionary.
						targs := recvType.RParams()
						genRecvType := recvType.OrigSym.Def.Type()
						nameNode = typecheck.Lookdot1(call.X, se.Sel, genRecvType, genRecvType.Methods(), 1).Nname.(*ir.Name)
						sym = g.getDictionarySym(nameNode, targs, true)
					} else {
						// This is the case of a normal
						// method call on a generic type.
						nameNode = call.X.(*ir.SelectorExpr).Selection.Nname.(*ir.Name)
						subtargs := deref(call.X.(*ir.SelectorExpr).X.Type()).RParams()
						s2targs := make([]*types.Type, len(subtargs))
						for i, t := range subtargs {
							s2targs[i] = subst.Typ(t)
						}
						sym = g.getDictionarySym(nameNode, s2targs, true)
					}
				} else {
					inst := call.X.(*ir.InstExpr)
					var nameNode *ir.Name
					var meth *ir.SelectorExpr
					var isMeth bool
					if meth, isMeth = inst.X.(*ir.SelectorExpr); isMeth {
						nameNode = meth.Selection.Nname.(*ir.Name)
					} else {
						nameNode = inst.X.(*ir.Name)
					}
					subtargs := typecheck.TypesOf(inst.Targs)
					for i, t := range subtargs {
						subtargs[i] = subst.Typ(t)
					}
					sym = g.getDictionarySym(nameNode, subtargs, isMeth)
				}

			case ir.OFUNCINST:
				inst := n.(*ir.InstExpr)
				nameNode := inst.X.(*ir.Name)
				subtargs := typecheck.TypesOf(inst.Targs)
				for i, t := range subtargs {
					subtargs[i] = subst.Typ(t)
				}
				sym = g.getDictionarySym(nameNode, subtargs, false)

			case ir.OXDOT:
				selExpr := n.(*ir.SelectorExpr)
				subtargs := selExpr.X.Type().RParams()
				s2targs := make([]*types.Type, len(subtargs))
				for i, t := range subtargs {
					s2targs[i] = subst.Typ(t)
				}
				nameNode := selExpr.Selection.Nname.(*ir.Name)
				sym = g.getDictionarySym(nameNode, s2targs, true)

			default:
				assert(false)
			}

			if sym == nil {
				// Unused sub-dictionary entry, just emit 0.
				off = objw.Uintptr(lsym, off, 0)
				infoPrint(" - Unused subdict entry\n")
			} else {
				off = objw.SymPtr(lsym, off, sym.Linksym(), 0)
				infoPrint(" - Subdict %v\n", sym.Name)
			}
		}
		objw.Global(lsym, int32(off), obj.DUPOK|obj.RODATA)
		infoPrint("=== Done dictionary\n")

		// Add any new, fully instantiated types seen during the substitution to g.instTypeList.
		g.instTypeList = append(g.instTypeList, subst.InstTypeList...)
	}
	return sym
}
func (g *irgen) getDictionaryValue(gf *ir.Name, targs []*types.Type, isMeth bool) ir.Node {
	sym := g.getDictionarySym(gf, targs, isMeth)

	// Make a node referencing the dictionary symbol.
	n := typecheck.NewName(sym)
	n.SetType(types.Types[types.TUINTPTR]) // should probably be [...]uintptr, but doesn't really matter
	n.SetTypecheck(1)
	n.Class = ir.PEXTERN
	sym.Def = n

	// Return the address of the dictionary.
	np := typecheck.NodAddr(n)
	// Note: treat dictionary pointers as uintptrs, so they aren't pointers
	// with respect to GC. That saves on stack scanning work, write barriers, etc.
	// We can get away with it because dictionaries are global variables.
	// TODO: use a cast, or is typing directly ok?
	np.SetType(types.Types[types.TUINTPTR])
	np.SetTypecheck(1)
	return np
}

// hasTParamNodes returns true if the type of any node in targs has a typeparam.
func hasTParamNodes(targs []ir.Node) bool {
	for _, n := range targs {
		if n.Type().HasTParam() {
			return true
		}
	}
	return false
}

// hasTParamNodes returns true if any type in targs has a typeparam.
func hasTParamTypes(targs []*types.Type) bool {
	for _, t := range targs {
		if t.HasTParam() {
			return true
		}
	}
	return false
}

// getGfInfo get information for a generic function - type params, derived generic
// types, and subdictionaries.
func (g *irgen) getGfInfo(gn *ir.Name) *gfInfo {
	infop := g.gfInfoMap[gn.Sym()]
	if infop != nil {
		return infop
	}

	checkFetchBody(gn)
	var info gfInfo
	gf := gn.Func
	recv := gf.Type().Recv()
	if recv != nil {
		info.tparams = deref(recv.Type).RParams()
	} else {
		tparams := gn.Type().TParams().FieldSlice()
		info.tparams = make([]*types.Type, len(tparams))
		for i, f := range tparams {
			info.tparams[i] = f.Type
		}
	}
	for _, n := range gf.Dcl {
		addType(&info, n, n.Type())
	}

	if infoPrintMode {
		fmt.Printf(">>> GfInfo for %v\n", gn)
		for _, t := range info.tparams {
			fmt.Printf("  Typeparam %v\n", t)
		}
	}

	var visitFunc func(ir.Node)
	visitFunc = func(n ir.Node) {
		if n.Op() == ir.OFUNCINST && !n.(*ir.InstExpr).Implicit() {
			if hasTParamNodes(n.(*ir.InstExpr).Targs) {
				infoPrint("  Closure&subdictionary required at generic function value %v\n", n.(*ir.InstExpr).X)
				info.subDictCalls = append(info.subDictCalls, n)
			}
		} else if n.Op() == ir.OXDOT && !n.(*ir.SelectorExpr).Implicit() &&
			n.(*ir.SelectorExpr).Selection != nil &&
			len(n.(*ir.SelectorExpr).X.Type().RParams()) > 0 {
			if n.(*ir.SelectorExpr).X.Op() == ir.OTYPE {
				infoPrint("  Closure&subdictionary required at generic meth expr %v\n", n)
			} else {
				infoPrint("  Closure&subdictionary required at generic meth value %v\n", n)
			}
			if hasTParamTypes(deref(n.(*ir.SelectorExpr).X.Type()).RParams()) {
				if n.(*ir.SelectorExpr).X.Op() == ir.OTYPE {
					infoPrint("  Closure&subdictionary required at generic meth expr %v\n", n)
				} else {
					infoPrint("  Closure&subdictionary required at generic meth value %v\n", n)
				}
				info.subDictCalls = append(info.subDictCalls, n)
			}
		}
		if n.Op() == ir.OCALL && n.(*ir.CallExpr).X.Op() == ir.OFUNCINST {
			n.(*ir.CallExpr).X.(*ir.InstExpr).SetImplicit(true)
			if hasTParamNodes(n.(*ir.CallExpr).X.(*ir.InstExpr).Targs) {
				infoPrint("  Subdictionary at generic function/method call: %v - %v\n", n.(*ir.CallExpr).X.(*ir.InstExpr).X, n)
				info.subDictCalls = append(info.subDictCalls, n)
			}
		}
		if n.Op() == ir.OCALL && n.(*ir.CallExpr).X.Op() == ir.OXDOT &&
			n.(*ir.CallExpr).X.(*ir.SelectorExpr).Selection != nil &&
			len(deref(n.(*ir.CallExpr).X.(*ir.SelectorExpr).X.Type()).RParams()) > 0 {
			n.(*ir.CallExpr).X.(*ir.SelectorExpr).SetImplicit(true)
			if hasTParamTypes(deref(n.(*ir.CallExpr).X.(*ir.SelectorExpr).X.Type()).RParams()) {
				infoPrint("  Subdictionary at generic method call: %v\n", n)
				info.subDictCalls = append(info.subDictCalls, n)
			}
		}
		if n.Op() == ir.OCALL && n.(*ir.CallExpr).X.Op() == ir.OXDOT &&
			n.(*ir.CallExpr).X.(*ir.SelectorExpr).Selection != nil &&
			deref(n.(*ir.CallExpr).X.(*ir.SelectorExpr).X.Type()).IsTypeParam() {
			n.(*ir.CallExpr).X.(*ir.SelectorExpr).SetImplicit(true)
			infoPrint("  Optional subdictionary at generic bound call: %v\n", n)
			info.subDictCalls = append(info.subDictCalls, n)
		}
		if n.Op() == ir.OCLOSURE {
			// Visit the closure body and add all relevant entries to the
			// dictionary of the outer function (closure will just use
			// the dictionary of the outer function).
			for _, n1 := range n.(*ir.ClosureExpr).Func.Body {
				ir.Visit(n1, visitFunc)
			}
		}

		addType(&info, n, n.Type())
	}

	for _, stmt := range gf.Body {
		ir.Visit(stmt, visitFunc)
	}
	if infoPrintMode {
		for _, t := range info.derivedTypes {
			fmt.Printf("  Derived type %v\n", t)
		}
		fmt.Printf(">>> Done Gfinfo\n")
	}
	g.gfInfoMap[gn.Sym()] = &info
	return &info
}

// addType adds t to info.derivedTypes if it is parameterized type (which is not
// just a simple type param) that is different from any existing type on
// info.derivedTypes.
func addType(info *gfInfo, n ir.Node, t *types.Type) {
	if t == nil || !t.HasTParam() {
		return
	}
	if t.IsTypeParam() && t.Underlying() == t {
		return
	}
	if t.Kind() == types.TFUNC && n != nil &&
		(n.Op() != ir.ONAME || n.Name().Class == ir.PFUNC) {
		// For now, only record function types that are associate with a
		// local/global variable (a name which is not a named global
		// function).
		return
	}
	if t.Kind() == types.TSTRUCT && t.IsFuncArgStruct() {
		// Multiple return values are not a relevant new type (?).
		return
	}
	// Ignore a derived type we've already added.
	for _, et := range info.derivedTypes {
		if types.Identical(t, et) {
			return
		}
	}
	info.derivedTypes = append(info.derivedTypes, t)
}
