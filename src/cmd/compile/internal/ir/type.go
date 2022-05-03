// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ir

import (
	"cmd/compile/internal/base"
	"cmd/compile/internal/types"
	"cmd/internal/src"
	"fmt"
)

// Nodes that represent the syntax of a type before type-checking.
// After type-checking, they serve only as shells around a *types.Type.
// Calling TypeNode converts a *types.Type to a Node shell.

// An Ntype is a Node that syntactically looks like a type.
// It can be the raw syntax for a type before typechecking,
// or it can be an OTYPE with Type() set to a *types.Type.
// Note that syntax doesn't guarantee it's a type: an expression
// like *fmt is an Ntype (we don't know whether names are types yet),
// but at least 1+1 is not an Ntype.
type Ntype interface {
	Node
	CanBeNtype()
}

// A Field is a declared function parameter.
// It is not a Node.
type Field struct {
	Pos   src.XPos
	Sym   *types.Sym
	Type  *types.Type
	IsDDD bool
}

func NewField(pos src.XPos, sym *types.Sym, typ *types.Type) *Field {
	return &Field{Pos: pos, Sym: sym, Type: typ}
}

func (f *Field) String() string {
	if f.Sym != nil {
		return fmt.Sprintf("%v %v", f.Sym, f.Type)
	}
	return fmt.Sprint(f.Type)
}

// A typeNode is a Node wrapper for type t.
type typeNode struct {
	miniNode
	typ *types.Type
}

func newTypeNode(typ *types.Type) *typeNode {
	n := &typeNode{typ: typ}
	n.pos = src.NoXPos
	n.op = OTYPE
	return n
}

func (n *typeNode) Type() *types.Type { return n.typ }
func (n *typeNode) Sym() *types.Sym   { return n.typ.Sym() }
func (n *typeNode) CanBeNtype()       {}

// TypeNode returns the Node representing the type t.
func TypeNode(t *types.Type) Ntype {
	if n := t.Obj(); n != nil {
		if n.Type() != t {
			base.Fatalf("type skew: %v has type %v, but expected %v", n, n.Type(), t)
		}
		return n.(Ntype)
	}
	return newTypeNode(t)
}

// A DynamicType represents the target type in a type switch.
type DynamicType struct {
	miniExpr
	X    Node // a *runtime._type for the targeted type
	ITab Node // for type switches from nonempty interfaces to non-interfaces, this is the itab for that pair.
}

func NewDynamicType(pos src.XPos, x Node) *DynamicType {
	n := &DynamicType{X: x}
	n.pos = pos
	n.op = ODYNAMICTYPE
	return n
}
