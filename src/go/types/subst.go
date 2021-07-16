// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements type parameter substitution.

package types

import (
	"bytes"
	"fmt"
	"go/token"
)

// TODO(rFindley) decide error codes for the errors in this file, and check
//                if error spans can be improved

type substMap struct {
	// The targs field is currently needed for *Named type substitution.
	// TODO(gri) rewrite that code, get rid of this field, and make this
	//           struct just the map (proj)
	targs []Type
	proj  map[*TypeParam]Type
}

// makeSubstMap creates a new substitution map mapping tpars[i] to targs[i].
// If targs[i] is nil, tpars[i] is not substituted.
func makeSubstMap(tpars []*TypeName, targs []Type) *substMap {
	assert(len(tpars) == len(targs))
	proj := make(map[*TypeParam]Type, len(tpars))
	for i, tpar := range tpars {
		// We must expand type arguments otherwise *instance
		// types end up as components in composite types.
		// TODO(gri) explain why this causes problems, if it does
		targ := expand(targs[i]) // possibly nil
		targs[i] = targ
		proj[tpar.typ.(*TypeParam)] = targ
	}
	return &substMap{targs, proj}
}

func (m *substMap) String() string {
	return fmt.Sprintf("%s", m.proj)
}

func (m *substMap) empty() bool {
	return len(m.proj) == 0
}

func (m *substMap) lookup(tpar *TypeParam) Type {
	if t := m.proj[tpar]; t != nil {
		return t
	}
	return tpar
}

// subst returns the type typ with its type parameters tpars replaced by
// the corresponding type arguments targs, recursively.
// subst is functional in the sense that it doesn't modify the incoming
// type. If a substitution took place, the result type is different from
// from the incoming type.
func (check *Checker) subst(pos token.Pos, typ Type, smap *substMap) Type {
	if smap.empty() {
		return typ
	}

	// common cases
	switch t := typ.(type) {
	case *Basic:
		return typ // nothing to do
	case *TypeParam:
		return smap.lookup(t)
	}

	// general case
	var subst subster
	subst.pos = pos
	subst.smap = smap
	if check != nil {
		subst.check = check
		subst.typMap = check.typMap
	} else {
		// If we don't have a *Checker and its global type map,
		// use a local version. Besides avoiding duplicate work,
		// the type map prevents infinite recursive substitution
		// for recursive types (example: type T[P any] *T[P]).
		subst.typMap = make(map[string]*Named)
	}
	return subst.typ(typ)
}

type subster struct {
	pos    token.Pos
	smap   *substMap
	check  *Checker // nil if called via Instantiate
	typMap map[string]*Named
}

func (subst *subster) typ(typ Type) Type {
	switch t := typ.(type) {
	case nil:
		// Call typOrNil if it's possible that typ is nil.
		panic("nil typ")

	case *Basic, *top:
		// nothing to do

	case *Array:
		elem := subst.typOrNil(t.elem)
		if elem != t.elem {
			return &Array{len: t.len, elem: elem}
		}

	case *Slice:
		elem := subst.typOrNil(t.elem)
		if elem != t.elem {
			return &Slice{elem: elem}
		}

	case *Struct:
		if fields, copied := subst.varList(t.fields); copied {
			return &Struct{fields: fields, tags: t.tags}
		}

	case *Pointer:
		base := subst.typ(t.base)
		if base != t.base {
			return &Pointer{base: base}
		}

	case *Tuple:
		return subst.tuple(t)

	case *Signature:
		// TODO(gri) rethink the recv situation with respect to methods on parameterized types
		// recv := subst.var_(t.recv) // TODO(gri) this causes a stack overflow - explain
		recv := t.recv
		params := subst.tuple(t.params)
		results := subst.tuple(t.results)
		if recv != t.recv || params != t.params || results != t.results {
			return &Signature{
				rparams: t.rparams,
				// TODO(rFindley) why can't we nil out tparams here, rather than in
				//                instantiate above?
				tparams:  t.tparams,
				scope:    t.scope,
				recv:     recv,
				params:   params,
				results:  results,
				variadic: t.variadic,
			}
		}

	case *Union:
		types, copied := subst.typeList(t.types)
		if copied {
			// TODO(gri) Remove duplicates that may have crept in after substitution
			//           (unlikely but possible). This matters for the Identical
			//           predicate on unions.
			return newUnion(types, t.tilde)
		}

	case *Interface:
		methods, mcopied := subst.funcList(t.methods)
		embeddeds, ecopied := subst.typeList(t.embeddeds)
		if mcopied || ecopied {
			iface := &Interface{methods: methods, embeddeds: embeddeds, complete: t.complete}
			if subst.check == nil {
				panic("internal error: cannot instantiate interfaces yet")
			}
			return iface
		}

	case *Map:
		key := subst.typ(t.key)
		elem := subst.typ(t.elem)
		if key != t.key || elem != t.elem {
			return &Map{key: key, elem: elem}
		}

	case *Chan:
		elem := subst.typ(t.elem)
		if elem != t.elem {
			return &Chan{dir: t.dir, elem: elem}
		}

	case *Named:
		// dump is for debugging
		dump := func(string, ...interface{}) {}
		if subst.check != nil && trace {
			subst.check.indent++
			defer func() {
				subst.check.indent--
			}()
			dump = func(format string, args ...interface{}) {
				subst.check.trace(subst.pos, format, args...)
			}
		}

		if t.TParams() == nil {
			dump(">>> %s is not parameterized", t)
			return t // type is not parameterized
		}

		var newTargs []Type

		if len(t.targs) > 0 {
			// already instantiated
			dump(">>> %s already instantiated", t)
			assert(len(t.targs) == len(t.TParams()))
			// For each (existing) type argument targ, determine if it needs
			// to be substituted; i.e., if it is or contains a type parameter
			// that has a type argument for it.
			for i, targ := range t.targs {
				dump(">>> %d targ = %s", i, targ)
				newTarg := subst.typ(targ)
				if newTarg != targ {
					dump(">>> substituted %d targ %s => %s", i, targ, newTarg)
					if newTargs == nil {
						newTargs = make([]Type, len(t.TParams()))
						copy(newTargs, t.targs)
					}
					newTargs[i] = newTarg
				}
			}

			if newTargs == nil {
				dump(">>> nothing to substitute in %s", t)
				return t // nothing to substitute
			}
		} else {
			// not yet instantiated
			dump(">>> first instantiation of %s", t)
			// TODO(rFindley) can we instead subst the tparam types here?
			newTargs = subst.smap.targs
		}

		// before creating a new named type, check if we have this one already
		h := instantiatedHash(t, newTargs)
		dump(">>> new type hash: %s", h)
		if named, found := subst.typMap[h]; found {
			dump(">>> found %s", named)
			return named
		}

		// create a new named type and populate typMap to avoid endless recursion
		tname := NewTypeName(subst.pos, t.obj.pkg, t.obj.name, nil)
		named := subst.check.newNamed(tname, t, t.Underlying(), t.TParams(), t.methods) // method signatures are updated lazily
		named.targs = newTargs
		subst.typMap[h] = named

		// do the substitution
		dump(">>> subst %s with %s (new: %s)", t.underlying, subst.smap, newTargs)
		named.underlying = subst.typOrNil(t.Underlying())
		named.fromRHS = named.underlying // for cycle detection (Checker.validType)

		return named

	case *TypeParam:
		return subst.smap.lookup(t)

	case *instance:
		// TODO(gri) can we avoid the expansion here and just substitute the type parameters?
		return subst.typ(t.expand())

	default:
		panic("unimplemented")
	}

	return typ
}

var instanceHashing = 0

func instantiatedHash(typ *Named, targs []Type) string {
	assert(instanceHashing == 0)
	instanceHashing++
	var buf bytes.Buffer
	writeTypeName(&buf, typ.obj, nil)
	buf.WriteByte('[')
	writeTypeList(&buf, targs, nil, nil)
	buf.WriteByte(']')
	instanceHashing--

	// With respect to the represented type, whether a
	// type is fully expanded or stored as instance
	// does not matter - they are the same types.
	// Remove the instanceMarkers printed for instances.
	res := buf.Bytes()
	i := 0
	for _, b := range res {
		if b != instanceMarker {
			res[i] = b
			i++
		}
	}

	return string(res[:i])
}

func typeListString(list []Type) string {
	var buf bytes.Buffer
	writeTypeList(&buf, list, nil, nil)
	return buf.String()
}

// typOrNil is like typ but if the argument is nil it is replaced with Typ[Invalid].
// A nil type may appear in pathological cases such as type T[P any] []func(_ T([]_))
// where an array/slice element is accessed before it is set up.
func (subst *subster) typOrNil(typ Type) Type {
	if typ == nil {
		return Typ[Invalid]
	}
	return subst.typ(typ)
}

func (subst *subster) var_(v *Var) *Var {
	if v != nil {
		if typ := subst.typ(v.typ); typ != v.typ {
			copy := *v
			copy.typ = typ
			return &copy
		}
	}
	return v
}

func (subst *subster) tuple(t *Tuple) *Tuple {
	if t != nil {
		if vars, copied := subst.varList(t.vars); copied {
			return &Tuple{vars: vars}
		}
	}
	return t
}

func (subst *subster) varList(in []*Var) (out []*Var, copied bool) {
	out = in
	for i, v := range in {
		if w := subst.var_(v); w != v {
			if !copied {
				// first variable that got substituted => allocate new out slice
				// and copy all variables
				new := make([]*Var, len(in))
				copy(new, out)
				out = new
				copied = true
			}
			out[i] = w
		}
	}
	return
}

func (subst *subster) func_(f *Func) *Func {
	if f != nil {
		if typ := subst.typ(f.typ); typ != f.typ {
			copy := *f
			copy.typ = typ
			return &copy
		}
	}
	return f
}

func (subst *subster) funcList(in []*Func) (out []*Func, copied bool) {
	out = in
	for i, f := range in {
		if g := subst.func_(f); g != f {
			if !copied {
				// first function that got substituted => allocate new out slice
				// and copy all functions
				new := make([]*Func, len(in))
				copy(new, out)
				out = new
				copied = true
			}
			out[i] = g
		}
	}
	return
}

func (subst *subster) typeList(in []Type) (out []Type, copied bool) {
	out = in
	for i, t := range in {
		if u := subst.typ(t); u != t {
			if !copied {
				// first function that got substituted => allocate new out slice
				// and copy all functions
				new := make([]Type, len(in))
				copy(new, out)
				out = new
				copied = true
			}
			out[i] = u
		}
	}
	return
}
