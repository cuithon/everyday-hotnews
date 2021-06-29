// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types2

import (
	"cmd/compile/internal/syntax"
	"sync"
	"sync/atomic"
)

// A Type represents a type of Go.
// All types implement the Type interface.
type Type interface {
	// Underlying returns the underlying type of a type
	// w/o following forwarding chains. Only used by
	// client packages (here for backward-compatibility).
	Underlying() Type

	// String returns a string representation of a type.
	String() string
}

// BasicKind describes the kind of basic type.
type BasicKind int

const (
	Invalid BasicKind = iota // type is invalid

	// predeclared types
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	String
	UnsafePointer

	// types for untyped values
	UntypedBool
	UntypedInt
	UntypedRune
	UntypedFloat
	UntypedComplex
	UntypedString
	UntypedNil

	// aliases
	Byte = Uint8
	Rune = Int32
)

// BasicInfo is a set of flags describing properties of a basic type.
type BasicInfo int

// Properties of basic types.
const (
	IsBoolean BasicInfo = 1 << iota
	IsInteger
	IsUnsigned
	IsFloat
	IsComplex
	IsString
	IsUntyped

	IsOrdered   = IsInteger | IsFloat | IsString
	IsNumeric   = IsInteger | IsFloat | IsComplex
	IsConstType = IsBoolean | IsNumeric | IsString
)

// A Basic represents a basic type.
type Basic struct {
	kind BasicKind
	info BasicInfo
	name string
}

// Kind returns the kind of basic type b.
func (b *Basic) Kind() BasicKind { return b.kind }

// Info returns information about properties of basic type b.
func (b *Basic) Info() BasicInfo { return b.info }

// Name returns the name of basic type b.
func (b *Basic) Name() string { return b.name }

// An Array represents an array type.
type Array struct {
	len  int64
	elem Type
}

// NewArray returns a new array type for the given element type and length.
// A negative length indicates an unknown length.
func NewArray(elem Type, len int64) *Array { return &Array{len: len, elem: elem} }

// Len returns the length of array a.
// A negative result indicates an unknown length.
func (a *Array) Len() int64 { return a.len }

// Elem returns element type of array a.
func (a *Array) Elem() Type { return a.elem }

// A Slice represents a slice type.
type Slice struct {
	elem Type
}

// NewSlice returns a new slice type for the given element type.
func NewSlice(elem Type) *Slice { return &Slice{elem: elem} }

// Elem returns the element type of slice s.
func (s *Slice) Elem() Type { return s.elem }

// A Struct represents a struct type.
type Struct struct {
	fields []*Var
	tags   []string // field tags; nil if there are no tags
}

// NewStruct returns a new struct with the given fields and corresponding field tags.
// If a field with index i has a tag, tags[i] must be that tag, but len(tags) may be
// only as long as required to hold the tag with the largest index i. Consequently,
// if no field has a tag, tags may be nil.
func NewStruct(fields []*Var, tags []string) *Struct {
	var fset objset
	for _, f := range fields {
		if f.name != "_" && fset.insert(f) != nil {
			panic("multiple fields with the same name")
		}
	}
	if len(tags) > len(fields) {
		panic("more tags than fields")
	}
	return &Struct{fields: fields, tags: tags}
}

// NumFields returns the number of fields in the struct (including blank and embedded fields).
func (s *Struct) NumFields() int { return len(s.fields) }

// Field returns the i'th field for 0 <= i < NumFields().
func (s *Struct) Field(i int) *Var { return s.fields[i] }

// Tag returns the i'th field tag for 0 <= i < NumFields().
func (s *Struct) Tag(i int) string {
	if i < len(s.tags) {
		return s.tags[i]
	}
	return ""
}

// A Pointer represents a pointer type.
type Pointer struct {
	base Type // element type
}

// NewPointer returns a new pointer type for the given element (base) type.
func NewPointer(elem Type) *Pointer { return &Pointer{base: elem} }

// Elem returns the element type for the given pointer p.
func (p *Pointer) Elem() Type { return p.base }

// A Tuple represents an ordered list of variables; a nil *Tuple is a valid (empty) tuple.
// Tuples are used as components of signatures and to represent the type of multiple
// assignments; they are not first class types of Go.
type Tuple struct {
	vars []*Var
}

// NewTuple returns a new tuple for the given variables.
func NewTuple(x ...*Var) *Tuple {
	if len(x) > 0 {
		return &Tuple{vars: x}
	}
	// TODO(gri) Don't represent empty tuples with a (*Tuple)(nil) pointer;
	//           it's too subtle and causes problems.
	return nil
}

// Len returns the number variables of tuple t.
func (t *Tuple) Len() int {
	if t != nil {
		return len(t.vars)
	}
	return 0
}

// At returns the i'th variable of tuple t.
func (t *Tuple) At(i int) *Var { return t.vars[i] }

// A Signature represents a (non-builtin) function or method type.
// The receiver is ignored when comparing signatures for identity.
type Signature struct {
	// We need to keep the scope in Signature (rather than passing it around
	// and store it in the Func Object) because when type-checking a function
	// literal we call the general type checker which returns a general Type.
	// We then unpack the *Signature and use the scope for the literal body.
	rparams  []*TypeName // receiver type parameters from left to right; or nil
	tparams  []*TypeName // type parameters from left to right; or nil
	scope    *Scope      // function scope, present for package-local signatures
	recv     *Var        // nil if not a method
	params   *Tuple      // (incoming) parameters from left to right; or nil
	results  *Tuple      // (outgoing) results from left to right; or nil
	variadic bool        // true if the last parameter's type is of the form ...T (or string, for append built-in only)
}

// NewSignature returns a new function type for the given receiver, parameters,
// and results, either of which may be nil. If variadic is set, the function
// is variadic, it must have at least one parameter, and the last parameter
// must be of unnamed slice type.
func NewSignature(recv *Var, params, results *Tuple, variadic bool) *Signature {
	if variadic {
		n := params.Len()
		if n == 0 {
			panic("types2.NewSignature: variadic function must have at least one parameter")
		}
		if _, ok := params.At(n - 1).typ.(*Slice); !ok {
			panic("types2.NewSignature: variadic parameter must be of unnamed slice type")
		}
	}
	return &Signature{recv: recv, params: params, results: results, variadic: variadic}
}

// Recv returns the receiver of signature s (if a method), or nil if a
// function. It is ignored when comparing signatures for identity.
//
// For an abstract method, Recv returns the enclosing interface either
// as a *Named or an *Interface. Due to embedding, an interface may
// contain methods whose receiver type is a different interface.
func (s *Signature) Recv() *Var { return s.recv }

// TParams returns the type parameters of signature s, or nil.
func (s *Signature) TParams() []*TypeName { return s.tparams }

// RParams returns the receiver type params of signature s, or nil.
func (s *Signature) RParams() []*TypeName { return s.rparams }

// SetTParams sets the type parameters of signature s.
func (s *Signature) SetTParams(tparams []*TypeName) { s.tparams = tparams }

// SetRParams sets the receiver type params of signature s.
func (s *Signature) SetRParams(rparams []*TypeName) { s.rparams = rparams }

// Params returns the parameters of signature s, or nil.
func (s *Signature) Params() *Tuple { return s.params }

// Results returns the results of signature s, or nil.
func (s *Signature) Results() *Tuple { return s.results }

// Variadic reports whether the signature s is variadic.
func (s *Signature) Variadic() bool { return s.variadic }

// An Interface represents an interface type.
type Interface struct {
	obj       Object        // type name object defining this interface; or nil (for better error messages)
	methods   []*Func       // ordered list of explicitly declared methods
	embeddeds []Type        // ordered list of explicitly embedded elements
	embedPos  *[]syntax.Pos // positions of embedded elements; or nil (for error messages) - use pointer to save space
	complete  bool          // indicates that all fields (except for tset) are set up

	tset *TypeSet // type set described by this interface, computed lazily
}

// typeSet returns the type set for interface t.
func (t *Interface) typeSet() *TypeSet { return newTypeSet(nil, nopos, t) }

// is reports whether interface t represents types that all satisfy f.
func (t *Interface) is(f func(Type, bool) bool) bool {
	switch t := t.typeSet().types.(type) {
	case nil, *top:
		// TODO(gri) should settle on top or nil to represent this case
		return false // we must have at least one type! (was bug)
	case *Union:
		return t.is(func(typ Type, tilde bool) bool { return f(typ, tilde) })
	default:
		return f(t, false)
	}
}

// emptyInterface represents the empty interface
var emptyInterface = Interface{complete: true, tset: &topTypeSet}

// NewInterface returns a new interface for the given methods and embedded types.
// NewInterface takes ownership of the provided methods and may modify their types
// by setting missing receivers.
//
// Deprecated: Use NewInterfaceType instead which allows arbitrary embedded types.
func NewInterface(methods []*Func, embeddeds []*Named) *Interface {
	tnames := make([]Type, len(embeddeds))
	for i, t := range embeddeds {
		tnames[i] = t
	}
	return NewInterfaceType(methods, tnames)
}

// NewInterfaceType returns a new interface for the given methods and embedded types.
// NewInterfaceType takes ownership of the provided methods and may modify their types
// by setting missing receivers.
func NewInterfaceType(methods []*Func, embeddeds []Type) *Interface {
	if len(methods) == 0 && len(embeddeds) == 0 {
		return &emptyInterface
	}

	// set method receivers if necessary
	typ := new(Interface)
	for _, m := range methods {
		if sig := m.typ.(*Signature); sig.recv == nil {
			sig.recv = NewVar(m.pos, m.pkg, "", typ)
		}
	}

	// sort for API stability
	sortMethods(methods)

	typ.methods = methods
	typ.embeddeds = embeddeds
	typ.complete = true

	return typ
}

// NumExplicitMethods returns the number of explicitly declared methods of interface t.
func (t *Interface) NumExplicitMethods() int { return len(t.methods) }

// ExplicitMethod returns the i'th explicitly declared method of interface t for 0 <= i < t.NumExplicitMethods().
// The methods are ordered by their unique Id.
func (t *Interface) ExplicitMethod(i int) *Func { return t.methods[i] }

// NumEmbeddeds returns the number of embedded types in interface t.
func (t *Interface) NumEmbeddeds() int { return len(t.embeddeds) }

// Embedded returns the i'th embedded defined (*Named) type of interface t for 0 <= i < t.NumEmbeddeds().
// The result is nil if the i'th embedded type is not a defined type.
//
// Deprecated: Use EmbeddedType which is not restricted to defined (*Named) types.
func (t *Interface) Embedded(i int) *Named { tname, _ := t.embeddeds[i].(*Named); return tname }

// EmbeddedType returns the i'th embedded type of interface t for 0 <= i < t.NumEmbeddeds().
func (t *Interface) EmbeddedType(i int) Type { return t.embeddeds[i] }

// NumMethods returns the total number of methods of interface t.
func (t *Interface) NumMethods() int { return t.typeSet().NumMethods() }

// Method returns the i'th method of interface t for 0 <= i < t.NumMethods().
// The methods are ordered by their unique Id.
func (t *Interface) Method(i int) *Func { return t.typeSet().Method(i) }

// Empty reports whether t is the empty interface.
func (t *Interface) Empty() bool { return t.typeSet().IsTop() }

// IsComparable reports whether interface t is or embeds the predeclared interface "comparable".
func (t *Interface) IsComparable() bool { return t.typeSet().IsComparable() }

// IsConstraint reports whether interface t is not just a method set.
func (t *Interface) IsConstraint() bool { return !t.typeSet().IsMethodSet() }

// isSatisfiedBy reports whether interface t's type list is satisfied by the type typ.
// If the type list is empty (absent), typ trivially satisfies the interface.
// TODO(gri) This is not a great name. Eventually, we should have a more comprehensive
//           "implements" predicate.
func (t *Interface) isSatisfiedBy(typ Type) bool {
	switch t := t.typeSet().types.(type) {
	case nil:
		return true // no type restrictions
	case *Union:
		r, _ := t.intersect(typ, false)
		return r != nil
	default:
		return Identical(t, typ)
	}
}

// Complete computes the interface's type set. It must be called by users of
// NewInterfaceType and NewInterface after the interface's embedded types are
// fully defined and before using the interface type in any way other than to
// form other types. The interface must not contain duplicate methods or a
// panic occurs. Complete returns the receiver.
//
// Deprecated: Type sets are now computed lazily, on demand; this function
//             is only here for backward-compatibility. It does not have to
//             be called explicitly anymore.
func (t *Interface) Complete() *Interface {
	// Some tests are still depending on the state change
	// (string representation of an Interface not containing an
	// /* incomplete */ marker) caused by the explicit Complete
	// call, so we compute the type set eagerly here.
	t.complete = true
	t.typeSet()
	return t
}

// A Map represents a map type.
type Map struct {
	key, elem Type
}

// NewMap returns a new map for the given key and element types.
func NewMap(key, elem Type) *Map {
	return &Map{key: key, elem: elem}
}

// Key returns the key type of map m.
func (m *Map) Key() Type { return m.key }

// Elem returns the element type of map m.
func (m *Map) Elem() Type { return m.elem }

// A Chan represents a channel type.
type Chan struct {
	dir  ChanDir
	elem Type
}

// A ChanDir value indicates a channel direction.
type ChanDir int

// The direction of a channel is indicated by one of these constants.
const (
	SendRecv ChanDir = iota
	SendOnly
	RecvOnly
)

// NewChan returns a new channel type for the given direction and element type.
func NewChan(dir ChanDir, elem Type) *Chan {
	return &Chan{dir: dir, elem: elem}
}

// Dir returns the direction of channel c.
func (c *Chan) Dir() ChanDir { return c.dir }

// Elem returns the element type of channel c.
func (c *Chan) Elem() Type { return c.elem }

// TODO(gri) Clean up Named struct below; specifically the fromRHS field (can we use underlying?).

// A Named represents a named (defined) type.
type Named struct {
	check      *Checker    // for Named.under implementation; nilled once under has been called
	info       typeInfo    // for cycle detection
	obj        *TypeName   // corresponding declared object
	orig       *Named      // original, uninstantiated type
	fromRHS    Type        // type (on RHS of declaration) this *Named type is derived from (for cycle reporting)
	underlying Type        // possibly a *Named during setup; never a *Named once set up completely
	tparams    []*TypeName // type parameters, or nil
	targs      []Type      // type arguments (after instantiation), or nil
	methods    []*Func     // methods declared for this type (not the method set of this type); signatures are type-checked lazily

	resolve func(*Named) ([]*TypeName, Type, []*Func)
	once    sync.Once
}

// NewNamed returns a new named type for the given type name, underlying type, and associated methods.
// If the given type name obj doesn't have a type yet, its type is set to the returned named type.
// The underlying type must not be a *Named.
func NewNamed(obj *TypeName, underlying Type, methods []*Func) *Named {
	if _, ok := underlying.(*Named); ok {
		panic("types2.NewNamed: underlying type must not be *Named")
	}
	return (*Checker)(nil).newNamed(obj, nil, underlying, nil, methods)
}

func (t *Named) expand() *Named {
	if t.resolve == nil {
		return t
	}

	t.once.Do(func() {
		// TODO(mdempsky): Since we're passing t to resolve anyway
		// (necessary because types2 expects the receiver type for methods
		// on defined interface types to be the Named rather than the
		// underlying Interface), maybe it should just handle calling
		// SetTParams, SetUnderlying, and AddMethod instead?  Those
		// methods would need to support reentrant calls though.  It would
		// also make the API more future-proof towards further extensions
		// (like SetTParams).

		tparams, underlying, methods := t.resolve(t)

		switch underlying.(type) {
		case nil, *Named:
			panic("invalid underlying type")
		}

		t.tparams = tparams
		t.underlying = underlying
		t.methods = methods
	})
	return t
}

// newNamed is like NewNamed but with a *Checker receiver and additional orig argument.
func (check *Checker) newNamed(obj *TypeName, orig *Named, underlying Type, tparams []*TypeName, methods []*Func) *Named {
	typ := &Named{check: check, obj: obj, orig: orig, fromRHS: underlying, underlying: underlying, tparams: tparams, methods: methods}
	if typ.orig == nil {
		typ.orig = typ
	}
	if obj.typ == nil {
		obj.typ = typ
	}
	// Ensure that typ is always expanded, at which point the check field can be
	// nilled out.
	//
	// Note that currently we cannot nil out check inside typ.under(), because
	// it's possible that typ is expanded multiple times.
	//
	// TODO(gri): clean this up so that under is the only function mutating
	//            named types.
	if check != nil {
		check.later(func() {
			switch typ.under().(type) {
			case *Named, *instance:
				panic("internal error: unexpanded underlying type")
			}
			typ.check = nil
		})
	}
	return typ
}

// Obj returns the type name for the named type t.
func (t *Named) Obj() *TypeName { return t.obj }

// Orig returns the original generic type an instantiated type is derived from.
// If t is not an instantiated type, the result is t.
func (t *Named) Orig() *Named { return t.orig }

// TODO(gri) Come up with a better representation and API to distinguish
//           between parameterized instantiated and non-instantiated types.

// TParams returns the type parameters of the named type t, or nil.
// The result is non-nil for an (originally) parameterized type even if it is instantiated.
func (t *Named) TParams() []*TypeName { return t.expand().tparams }

// SetTParams sets the type parameters of the named type t.
func (t *Named) SetTParams(tparams []*TypeName) { t.expand().tparams = tparams }

// TArgs returns the type arguments after instantiation of the named type t, or nil if not instantiated.
func (t *Named) TArgs() []Type { return t.targs }

// SetTArgs sets the type arguments of the named type t.
func (t *Named) SetTArgs(args []Type) { t.targs = args }

// NumMethods returns the number of explicit methods whose receiver is named type t.
func (t *Named) NumMethods() int { return len(t.expand().methods) }

// Method returns the i'th method of named type t for 0 <= i < t.NumMethods().
func (t *Named) Method(i int) *Func { return t.expand().methods[i] }

// SetUnderlying sets the underlying type and marks t as complete.
func (t *Named) SetUnderlying(underlying Type) {
	if underlying == nil {
		panic("types2.Named.SetUnderlying: underlying type must not be nil")
	}
	if _, ok := underlying.(*Named); ok {
		panic("types2.Named.SetUnderlying: underlying type must not be *Named")
	}
	t.expand().underlying = underlying
}

// AddMethod adds method m unless it is already in the method list.
func (t *Named) AddMethod(m *Func) {
	t.expand()
	if i, _ := lookupMethod(t.methods, m.pkg, m.name); i < 0 {
		t.methods = append(t.methods, m)
	}
}

// Note: This is a uint32 rather than a uint64 because the
// respective 64 bit atomic instructions are not available
// on all platforms.
var lastID uint32

// nextID returns a value increasing monotonically by 1 with
// each call, starting with 1. It may be called concurrently.
func nextID() uint64 { return uint64(atomic.AddUint32(&lastID, 1)) }

// A TypeParam represents a type parameter type.
type TypeParam struct {
	check *Checker  // for lazy type bound completion
	id    uint64    // unique id, for debugging only
	obj   *TypeName // corresponding type name
	index int       // type parameter index in source order, starting at 0
	bound Type      // *Named or *Interface; underlying type is always *Interface
}

// Obj returns the type name for the type parameter t.
func (t *TypeParam) Obj() *TypeName { return t.obj }

// NewTypeParam returns a new TypeParam.  bound can be nil (and set later).
func (check *Checker) NewTypeParam(obj *TypeName, index int, bound Type) *TypeParam {
	// Always increment lastID, even if it is not used.
	id := nextID()
	if check != nil {
		check.nextID++
		id = check.nextID
	}
	typ := &TypeParam{check: check, id: id, obj: obj, index: index, bound: bound}
	if obj.typ == nil {
		obj.typ = typ
	}
	return typ
}

// Index returns the index of the type param within its param list.
func (t *TypeParam) Index() int {
	return t.index
}

// SetId sets the unique id of a type param. Should only be used for type params
// in imported generic types.
func (t *TypeParam) SetId(id uint64) {
	t.id = id
}

func (t *TypeParam) Bound() *Interface {
	// we may not have an interface (error reported elsewhere)
	iface, _ := under(t.bound).(*Interface)
	if iface == nil {
		return &emptyInterface
	}
	// use the type bound position if we have one
	pos := nopos
	if n, _ := t.bound.(*Named); n != nil {
		pos = n.obj.pos
	}
	// TODO(gri) switch this to an unexported method on Checker.
	newTypeSet(t.check, pos, iface)
	return iface
}

func (t *TypeParam) SetBound(bound Type) {
	if bound == nil {
		panic("types2.TypeParam.SetBound: bound must not be nil")
	}
	t.bound = bound
}

// optype returns a type's operational type. Except for
// type parameters, the operational type is the same
// as the underlying type (as returned by under). For
// Type parameters, the operational type is determined
// by the corresponding type bound's type list. The
// result may be the bottom or top type, but it is never
// the incoming type parameter.
func optype(typ Type) Type {
	if t := asTypeParam(typ); t != nil {
		// If the optype is typ, return the top type as we have
		// no information. It also prevents infinite recursion
		// via the asTypeParam converter function. This can happen
		// for a type parameter list of the form:
		// (type T interface { type T }).
		// See also issue #39680.
		if a := t.Bound().typeSet().types; a != nil {
			// If we have a union with a single entry, ignore
			// any tilde because under(~t) == under(t).
			if u, _ := a.(*Union); u != nil && u.NumTerms() == 1 {
				a = u.types[0]
			}
			if a != typ {
				// a != typ and a is a type parameter => under(a) != typ, so this is ok
				return under(a)
			}
		}
		return theTop
	}
	return under(typ)
}

// An instance represents an instantiated generic type syntactically
// (without expanding the instantiation). Type instances appear only
// during type-checking and are replaced by their fully instantiated
// (expanded) types before the end of type-checking.
type instance struct {
	check   *Checker     // for lazy instantiation
	pos     syntax.Pos   // position of type instantiation; for error reporting only
	base    *Named       // parameterized type to be instantiated
	targs   []Type       // type arguments
	poslist []syntax.Pos // position of each targ; for error reporting only
	value   Type         // base(targs...) after instantiation or Typ[Invalid]; nil if not yet set
}

// expand returns the instantiated (= expanded) type of t.
// The result is either an instantiated *Named type, or
// Typ[Invalid] if there was an error.
func (t *instance) expand() Type {
	v := t.value
	if v == nil {
		v = t.check.instantiate(t.pos, t.base, t.targs, t.poslist)
		if v == nil {
			v = Typ[Invalid]
		}
		t.value = v
	}
	// After instantiation we must have an invalid or a *Named type.
	if debug && v != Typ[Invalid] {
		_ = v.(*Named)
	}
	return v
}

// expand expands a type instance into its instantiated
// type and leaves all other types alone. expand does
// not recurse.
func expand(typ Type) Type {
	if t, _ := typ.(*instance); t != nil {
		return t.expand()
	}
	return typ
}

// expandf is set to expand.
// Call expandf when calling expand causes compile-time cycle error.
var expandf func(Type) Type

func init() { expandf = expand }

// top represents the top of the type lattice.
// It is the underlying type of a type parameter that
// can be satisfied by any type (ignoring methods),
// because its type constraint contains no restrictions
// besides methods.
type top struct{}

// theTop is the singleton top type.
var theTop = &top{}

// Type-specific implementations of Underlying.
func (t *Basic) Underlying() Type     { return t }
func (t *Array) Underlying() Type     { return t }
func (t *Slice) Underlying() Type     { return t }
func (t *Struct) Underlying() Type    { return t }
func (t *Pointer) Underlying() Type   { return t }
func (t *Tuple) Underlying() Type     { return t }
func (t *Signature) Underlying() Type { return t }
func (t *Interface) Underlying() Type { return t }
func (t *Map) Underlying() Type       { return t }
func (t *Chan) Underlying() Type      { return t }
func (t *Named) Underlying() Type     { return t.expand().underlying }
func (t *TypeParam) Underlying() Type { return t }
func (t *instance) Underlying() Type  { return t }
func (t *top) Underlying() Type       { return t }

// Type-specific implementations of String.
func (t *Basic) String() string     { return TypeString(t, nil) }
func (t *Array) String() string     { return TypeString(t, nil) }
func (t *Slice) String() string     { return TypeString(t, nil) }
func (t *Struct) String() string    { return TypeString(t, nil) }
func (t *Pointer) String() string   { return TypeString(t, nil) }
func (t *Tuple) String() string     { return TypeString(t, nil) }
func (t *Signature) String() string { return TypeString(t, nil) }
func (t *Interface) String() string { return TypeString(t, nil) }
func (t *Map) String() string       { return TypeString(t, nil) }
func (t *Chan) String() string      { return TypeString(t, nil) }
func (t *Named) String() string     { return TypeString(t, nil) }
func (t *TypeParam) String() string { return TypeString(t, nil) }
func (t *instance) String() string  { return TypeString(t, nil) }
func (t *top) String() string       { return TypeString(t, nil) }

// under returns the true expanded underlying type.
// If it doesn't exist, the result is Typ[Invalid].
// under must only be called when a type is known
// to be fully set up.
func under(t Type) Type {
	// TODO(gri) is this correct for *Union?
	if n := asNamed(t); n != nil {
		return n.under()
	}
	return t
}

// Converters
//
// A converter must only be called when a type is
// known to be fully set up. A converter returns
// a type's operational type (see comment for optype)
// or nil if the type argument is not of the
// respective type.

func asBasic(t Type) *Basic {
	op, _ := optype(t).(*Basic)
	return op
}

func asArray(t Type) *Array {
	op, _ := optype(t).(*Array)
	return op
}

func asSlice(t Type) *Slice {
	op, _ := optype(t).(*Slice)
	return op
}

func asStruct(t Type) *Struct {
	op, _ := optype(t).(*Struct)
	return op
}

func asPointer(t Type) *Pointer {
	op, _ := optype(t).(*Pointer)
	return op
}

// asTuple is not needed - not provided

func asSignature(t Type) *Signature {
	op, _ := optype(t).(*Signature)
	return op
}

func asInterface(t Type) *Interface {
	op, _ := optype(t).(*Interface)
	return op
}

func asMap(t Type) *Map {
	op, _ := optype(t).(*Map)
	return op
}

func asChan(t Type) *Chan {
	op, _ := optype(t).(*Chan)
	return op
}

// If the argument to asNamed and asTypeParam is of the respective types
// (possibly after expanding an instance type), these methods return that type.
// Otherwise the result is nil.

func asNamed(t Type) *Named {
	e, _ := expand(t).(*Named)
	return e
}

func asTypeParam(t Type) *TypeParam {
	u, _ := under(t).(*TypeParam)
	return u
}

// Exported for the compiler.

func AsPointer(t Type) *Pointer     { return asPointer(t) }
func AsNamed(t Type) *Named         { return asNamed(t) }
func AsSignature(t Type) *Signature { return asSignature(t) }
func AsInterface(t Type) *Interface { return asInterface(t) }
func AsTypeParam(t Type) *TypeParam { return asTypeParam(t) }
